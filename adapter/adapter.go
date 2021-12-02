package adapter

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	ms "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	defaultDSN               = "root:@tcp(localhost:3306)/chirpbird?parseTime=true"
	defaultDatabase          = "chirpbird"
	adpVersion               = 111
	adapterName              = "mysql"
	defaultMaxResults        = 1024
	defaultMaxMessageResults = 100
	txTimeoutMultiplier      = 1.5
)

type adapter struct {
	db                *sqlx.DB
	dsn               string
	dbName            string
	maxResults        int
	maxMessageResults int
	version           int
	sqlTimeout        time.Duration
	txTimeout         time.Duration
}

type configType struct {
	ms.Config
	DSN             string `json:"dsn,omitempty"`
	Database        string `json:"database,omitempty"`
	MaxOpenConns    int    `json:"max_open_conns,omitempty"`
	MaxIdleConns    int    `json:"max_idle_conns,omitempty"`
	ConnMaxLifetime int    `json:"conn_max_lifetime,omitempty"`
	SqlTimeout      int    `json:"sql_timeout,omitempty"`
}

type User struct {
	ID         int64     `json:"id"`
	Username   string    `json:"username"`
	CreateDate time.Time `json:"createdate"`
	UpdateDate time.Time `json:"updatedate"`
	State      string    `json:"state"`
	Lastseen   time.Time `json:"lastseen"`
	Public     string    `json:"public"`
	Tags       string    `json:"tags"`
}

func (a *adapter) getContext() (context.Context, context.CancelFunc) {
	if a.sqlTimeout > 0 {
		return context.WithTimeout(context.Background(), a.sqlTimeout)
	}
	return context.Background(), nil
}

// Open initializes database session
func (a *adapter) Open(jsonconfig json.RawMessage) error {
	if a.db != nil {
		return errors.New("mysql adapter is already connected")
	}

	if len(jsonconfig) < 2 {
		return errors.New("adapter mysql missing config")
	}

	var err error
	defaultCfg := ms.NewConfig()
	config := configType{Config: *defaultCfg}
	if err = json.Unmarshal(jsonconfig, &config); err != nil {
		return errors.New("mysql adapter failed to parse config: " + err.Error())
	}

	if dsn := config.FormatDSN(); dsn != defaultCfg.FormatDSN() {
		// MySql config is specified. Use it.
		a.dbName = config.DBName
		a.dsn = dsn
		if config.DSN != "" || config.Database != "" {
			return errors.New("mysql config: `dsn` and `database` fields are deprecated. Please, specify individual connection settings via mysql.Config: https://pkg.go.dev/github.com/go-sql-driver/mysql#Config")
		}
	} else {
		// Otherwise, use DSN and Database to configure database connection.
		// Note: this method is deprecated.
		if config.DSN != "" {
			a.dsn = config.DSN
		} else {
			a.dsn = defaultDSN
		}
		a.dbName = config.Database
	}

	if a.dbName == "" {
		a.dbName = defaultDatabase
	}

	if a.maxResults <= 0 {
		a.maxResults = defaultMaxResults
	}

	if a.maxMessageResults <= 0 {
		a.maxMessageResults = defaultMaxMessageResults
	}

	// This just initializes the driver but does not open the network connection.
	a.db, err = sqlx.Open("mysql", a.dsn)
	if err != nil {
		return err
	}

	// Opening connection.
	err = a.db.Ping()

	if err == nil {
		if config.MaxOpenConns > 0 {
			a.db.SetMaxOpenConns(config.MaxOpenConns)
		}
		if config.MaxIdleConns > 0 {
			a.db.SetMaxIdleConns(config.MaxIdleConns)
		}
		if config.ConnMaxLifetime > 0 {
			a.db.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Second)
		}
		if config.SqlTimeout > 0 {
			a.sqlTimeout = time.Duration(config.SqlTimeout) * time.Second
			// We allocate txTimeoutMultiplier times sqlTimeout for transactions.
			a.txTimeout = time.Duration(float64(config.SqlTimeout)*txTimeoutMultiplier) * time.Second
		}
	}
	return err
}

// Close closes the underlying database connection
func (a *adapter) Close() error {
	var err error
	if a.db != nil {
		err = a.db.Close()
		a.db = nil
		a.version = -1
	}
	return err
}

func (a *adapter) UserGet(username string) (*User, error) {
	ctx, cancel := a.getContext()
	if cancel != nil {
		defer cancel()
	}
	var user User
	err := a.db.GetContext(ctx, &user, "SELECT * FROM users WHERE username=?", username)
	if err == nil {
		return &user, nil
	}

	if err == sql.ErrNoRows {
		// Clear the error if user does not exist or marked as soft-deleted.
		return nil, nil
	}

	return nil, err
}
