package adapter

import (
	"context"
	"database/sql"
	"errors"
	"time"

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

// Open initializes database
func (a *adapter) Open() error {
	if a.db != nil {
		return errors.New("mysql adapter is already connected")
	}

	var err error

	a.dsn = defaultDSN
	a.dbName = defaultDatabase
	a.maxResults = defaultMaxResults
	a.maxMessageResults = defaultMaxMessageResults

	a.db, err = sqlx.Open("mysql", a.dsn)
	if err != nil {
		return err
	}

	// Opening connection.
	err = a.db.Ping()

	if err == nil {
		a.db.SetMaxOpenConns(64)
		a.db.SetMaxIdleConns(64)
		a.db.SetConnMaxLifetime(time.Duration(60) * time.Second)
		a.sqlTimeout = time.Duration(10) * time.Second
		a.txTimeout = time.Duration(float64(10)*txTimeoutMultiplier) * time.Second
	}

	return err
}

func (a *adapter) IsOpen() bool {
	return a.db != nil
}

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
		return nil, nil
	}

	return nil, err
}
