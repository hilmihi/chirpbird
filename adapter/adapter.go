package adapter

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	ms "github.com/go-sql-driver/mysql"
	guuid "github.com/google/uuid"
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
	State      int       `json:"int"`
	Lastseen   time.Time `json:"lastseen"`
}

type Room struct {
	ID         int64     `json:"id"`
	CreateDate time.Time `json:"createdate"`
	UpdateDate time.Time `json:"updatedate"`
	Name       string    `json:"name"`
	Owner      int64     `json:"owner"`
	Seqid      int64     `json:"seqid"`
	Public     interface{}
	Tags       interface{}
}

type Subscription struct {
	ID         int64      `json:"id"`
	CreateDate time.Time  `json:"createdate"`
	UpdateDate time.Time  `json:"updatedate"`
	DeleteDate *time.Time `json:",omitempty"`
	UserID     int64      `json:"userid"`
	Room       string     `json:"room"`
	Seqid      int64      `json:"seqid"`
	Comment    string     `json:"comment"`
	Public     interface{}
}

type Message struct {
	ID         int64      `json:"id"`
	CreateDate time.Time  `json:"createdate"`
	UpdateDate time.Time  `json:"updatedate"`
	DeleteDate *time.Time `json:",omitempty"`
	Seqid      int64      `json:"seqid"`
	Room       string     `json:"room"`
	From       int64      `json:"from"`
	Content    interface{}
	Username   string `json:"username"`
}

type QueryOpt struct {
	// Subscription query
	User  int64
	Topic string
}

var a adapter

func (a *adapter) getContext() (context.Context, context.CancelFunc) {
	if a.sqlTimeout > 0 {
		return context.WithTimeout(context.Background(), a.sqlTimeout)
	}
	return context.Background(), nil
}

func (a *adapter) getContextForTx() (context.Context, context.CancelFunc) {
	if a.txTimeout > 0 {
		return context.WithTimeout(context.Background(), a.txTimeout)
	}
	return context.Background(), nil
}

// Open initializes database
func Open(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) error {
	if a.db != nil {
		return errors.New("mysql is already connected")
	}

	var err error

	a.dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
	a.dbName = DbName
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

func IsOpen() bool {
	return a.db != nil
}

func Close() error {
	var err error
	if a.db != nil {
		err = a.db.Close()
		a.db = nil
		a.version = -1
	}
	return err
}

func CreateDb(reset bool) error {
	var err error
	var tx *sql.Tx

	// Can't use an existing connection because it's configured with a database name which may not exist.
	// Don't care if it does not close cleanly.
	a.db.Close()

	// This DSN has been parsed before and produced no error, not checking for errors here.
	cfg, _ := ms.ParseDSN(a.dsn)
	// Clear database name
	cfg.DBName = ""

	a.db, err = sqlx.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return err
	}

	if tx, err = a.db.Begin(); err != nil {
		return err
	}

	defer func() {
		if err != nil {
			// FIXME: This is useless: MySQL auto-commits on every CREATE TABLE.
			// Maybe DROP DATABASE instead.
			tx.Rollback()
		}
	}()

	if reset {
		if _, err = tx.Exec("DROP DATABASE IF EXISTS " + a.dbName); err != nil {
			return err
		}
	}

	if _, err = tx.Exec("CREATE DATABASE " + a.dbName + " CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci"); err != nil {
		return err
	}

	if _, err = tx.Exec("USE " + a.dbName); err != nil {
		return err
	}

	if _, err = tx.Exec(
		`CREATE TABLE users(
			id        INT NOT NULL AUTO_INCREMENT,
			uname   VARCHAR(32) NOT NULL,
			createdat DATETIME(3) NOT NULL,
			updatedat DATETIME(3) NOT NULL,
			state     SMALLINT NOT NULL DEFAULT 0,
			lastseen  DATETIME,
			PRIMARY KEY(id),
			UNIQUE INDEX auth_uname(uname)
		)`); err != nil {
		return err
	}

	// Rooms
	if _, err = tx.Exec(
		`CREATE TABLE rooms(
			id        INT NOT NULL AUTO_INCREMENT,
			createdat DATETIME(3) NOT NULL,
			updatedat DATETIME(3) NOT NULL,
			name      CHAR(25) NOT NULL,
			owner     INT NOT NULL DEFAULT 0,
			seqid     INT NOT NULL DEFAULT 0,
			public    JSON,
			tags      JSON,
			PRIMARY KEY(id),
			UNIQUE INDEX room_name(name),
			INDEX room_owner(owner)
		)`); err != nil {
		return err
	}

	// Subscriptions
	if _, err = tx.Exec(
		`CREATE TABLE subscriptions(
			id        INT NOT NULL AUTO_INCREMENT,
			createdat DATETIME(3) NOT NULL,
			updatedat DATETIME(3) NOT NULL,
			deletedat DATETIME(3),
			userid    INT NOT NULL,
			room      CHAR(25) NOT NULL,
			seqid     INT DEFAULT 0,
			comment   JSON,
			PRIMARY KEY(id),
			FOREIGN KEY(userid) REFERENCES users(id),
			UNIQUE INDEX subscriptions_room_userid(room, userid),
			INDEX subscriptions_room(room)
		)`); err != nil {
		return err
	}

	// Messages
	if _, err = tx.Exec(
		`CREATE TABLE messages(
			id        INT NOT NULL AUTO_INCREMENT,
			createdat DATETIME(3) NOT NULL,
			updatedat DATETIME(3) NOT NULL,
			deletedat DATETIME(3),
			seqid     INT NOT NULL,
			room      CHAR(25) NOT NULL,` +
			"`from`   BIGINT NOT NULL," +
			`content   JSON,
			PRIMARY KEY(id),
			FOREIGN KEY(room) REFERENCES rooms(name),
			UNIQUE INDEX messages_room_seqid(room, seqid)
		);`); err != nil {
		return err
	}

	return tx.Commit()
}

func UserGet(username string) (*User, error) {
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

func UserGetByID(id int64) (*User, error) {
	ctx, cancel := a.getContext()
	if cancel != nil {
		defer cancel()
	}
	// Fetch room by name
	var usr = new(User)
	err := a.db.GetContext(ctx, usr,
		"SELECT * FROM users WHERE id=?",
		id)

	if err != nil {
		if err == sql.ErrNoRows {
			// Nothing found - clear the error
			err = nil
		}
		return nil, err
	}

	return usr, err
}

func UserCreate(user *User) (int64, error) {
	ctx, cancel := a.getContextForTx()
	if cancel != nil {
		defer cancel()
	}
	tx, err := a.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	id, err := a.userCreate(tx, user)
	if err != nil {
		return 0, err
	}
	return id, tx.Commit()
}

func (a *adapter) userCreate(tx *sqlx.Tx, user *User) (int64, error) {
	result, err := tx.Exec("INSERT INTO users(username,createdate,updatedate,state,lastseen) "+
		"VALUES(?,?,?,?,?)",
		user.Username, user.CreateDate, user.UpdateDate, user.State, user.Lastseen)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil
}

func RoomCreate(room *Room) (int64, error) {
	ctx, cancel := a.getContextForTx()
	if cancel != nil {
		defer cancel()
	}
	tx, err := a.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	id, err := a.roomCreate(tx, room)
	if err != nil {
		return 0, err
	}
	return id, tx.Commit()
}

func (a *adapter) roomCreate(tx *sqlx.Tx, room *Room) (int64, error) {
	result, err := tx.Exec("INSERT INTO rooms(createdate,updatedate,name,owner,public,tags) "+
		"VALUES(?,?,?,?,?,?)",
		room.CreateDate, room.UpdateDate, room.Name, room.Owner, toJSON(room.Public), toJSON(room.Tags))
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil
}

func createSubscription(tx *sqlx.Tx, sub *Subscription) error {
	_, err := tx.Exec(
		"INSERT INTO subscriptions(createdate,updatedate,deletedate,userid,room,comment) "+
			"VALUES(?,?,NULL,?,?,?)",
		sub.CreateDate, sub.UpdateDate, sub.UserID, sub.Room, toJSON(sub.Comment))

	if err != nil && isDuplicate(err) {
		_, err = tx.Exec("UPDATE subscriptions SET createdate=?,updatedate=?,deletedate=NULL,"+
			"comment=? WHERE room=? AND userid=?",
			sub.CreateDate, sub.UpdateDate, sub.Comment, sub.Room, sub.UserID)
	}

	return err
}

func RoomGetByID(id int64) (*Room, error) {
	ctx, cancel := a.getContext()
	if cancel != nil {
		defer cancel()
	}
	// Fetch room by name
	var rm = new(Room)
	err := a.db.GetContext(ctx, rm,
		"SELECT id, createdate,updatedate,name,owner,seqid,public,tags "+
			"FROM rooms WHERE id=?",
		id)

	if err != nil {
		if err == sql.ErrNoRows {
			// Nothing found - clear the error
			err = nil
		}
		return nil, err
	}

	return rm, err
}

func RoomCreateP2P(initiator *Subscription, invited []User) error {
	ctx, cancel := a.getContextForTx()
	if cancel != nil {
		defer cancel()
	}
	tx, err := a.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = createSubscription(tx, initiator)
	if err != nil {
		return err
	}

	for _, m := range invited {
		invt := &Subscription{
			CreateDate: time.Now().UTC().Round(time.Millisecond),
			UpdateDate: time.Now().UTC().Round(time.Millisecond),
			UserID:     m.ID,
			Room:       initiator.Room,
		}
		err = createSubscription(tx, invt)
	}

	if err != nil {
		return err
	}

	return tx.Commit()
}

func encodeUidString(uUID guuid.UUID) {
	panic("unimplemented")
}

// Get Room by name.
func RoomGet(room string) (*Room, error) {
	ctx, cancel := a.getContext()
	if cancel != nil {
		defer cancel()
	}
	// Fetch room by name
	var rm = new(Room)
	err := a.db.GetContext(ctx, rm,
		"SELECT createdate,updatedate,name AS id,owner,seqid,public,tags "+
			"FROM rooms WHERE name=?",
		room)

	if err != nil {
		if err == sql.ErrNoRows {
			// Nothing found - clear the error
			err = nil
		}
		return nil, err
	}

	return rm, nil
}

func UsersSub(topic string, opts *QueryOpt) ([]Subscription, error) {
	// Fetch all subscribed users. The number of users is not large
	q := `SELECT s.createdate,s.updatedate,s.deletedate,s.userid,s.room
		FROM subscriptions AS s JOIN users AS u ON s.userid=u.id 
		WHERE s.topic=? 
		LIMIT ?`
	args := []interface{}{topic, a.maxResults}

	ctx, cancel := a.getContext()
	if cancel != nil {
		defer cancel()
	}
	rows, err := a.db.QueryxContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}

	// Fetch subscriptions
	var sub Subscription
	var subs []Subscription
	for rows.Next() {
		if err = rows.Scan(
			&sub.CreateDate, &sub.UpdateDate, &sub.DeleteDate,
			&sub.UserID, &sub.Room); err != nil {
			break
		}

		subs = append(subs, sub)
	}
	if err == nil {
		err = rows.Err()
	}
	rows.Close()

	return subs, err
}

func SubsByUser(userid int64) ([]Subscription, error) {
	q := `SELECT r.id, s.createdate,s.updatedate,s.deletedate,s.userid,s.room,r.public
	FROM subscriptions AS s JOIN rooms AS r ON s.room=r.name
	WHERE s.userid=? AND deletedate IS NULL`
	args := []interface{}{userid}

	ctx, cancel := a.getContext()
	if cancel != nil {
		defer cancel()
	}
	rows, err := a.db.QueryxContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}

	var subs []Subscription
	var sub Subscription
	for rows.Next() {
		if err = rows.StructScan(&sub); err != nil {
			break
		}
		sub.Public = fromJSON(sub.Public)
		subs = append(subs, sub)
	}
	if err == nil {
		err = rows.Err()
	}
	rows.Close()

	return subs, err
}

func SubsDelete(topic string, userid int64) error {
	tx, err := a.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	ctx, cancel := a.getContext()
	if cancel != nil {
		defer cancel()
	}

	now := time.Now().UTC().Round(time.Millisecond)
	res, err := tx.ExecContext(ctx,
		"UPDATE subscriptions SET updatedate=?,deletedate=? WHERE topic=? AND userid=? AND deletedate IS NULL",
		now, now, topic, userid)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err == nil && affected == 0 {
		return errors.New("Not Found")
	}

	return tx.Commit()
}

func FindUsers(word string, userid int64) ([]User, error) {

	query := `
		SELECT u.id,u.username,u.createdate,u.updatedate,u.state,u.lastseen
		FROM users AS u 
		WHERE u.state=? AND u.username LIKE ? AND u.id <> ?`
	args := []interface{}{1, "%" + word + "%", userid}

	ctx, cancel := a.getContext()
	if cancel != nil {
		defer cancel()
	}
	// Get users matched by tags, sort by number of matches from high to low.
	rows, err := a.db.QueryxContext(ctx, query, args...)

	if err != nil {
		return nil, err
	}

	var users []User
	var user User
	for rows.Next() {
		if err = rows.StructScan(&user); err != nil {
			break
		}
		users = append(users, user)
	}
	if err == nil {
		err = rows.Err()
	}
	rows.Close()

	return users, err
}

func SaveMessage(msg *Message) error {
	ctx, cancel := a.getContext()
	if cancel != nil {
		defer cancel()
	}
	// store assignes message ID, but we don't use it. Message IDs are not used anywhere.
	// Using a sequential ID provided by the database.
	_, err := a.db.ExecContext(
		ctx,
		"INSERT INTO messages(createdate,updatedate,deletedate,seqid,room,`from`,content) VALUES(?,?,NULL,?,?,?,?)",
		msg.CreateDate, msg.UpdateDate, msg.Seqid, msg.Room, msg.From, toJSON(msg.Content))

	return err
}

func MessageByRoom(room string) ([]Message, error) {
	var limit = a.maxMessageResults

	ctx, cancel := a.getContext()
	if cancel != nil {
		defer cancel()
	}
	rows, err := a.db.QueryxContext(
		ctx,
		"SELECT m.id, m.createdate,m.updatedate,m.deletedate,m.seqid,m.room,m.`from`,m.content, u.username"+
			" FROM messages AS m JOIN users AS u ON m.from=u.id"+
			" WHERE m.room=?"+
			" ORDER BY m.seqid ASC LIMIT ?",
		room, limit)

	if err != nil {
		return nil, err
	}

	msgs := make([]Message, 0, limit)
	for rows.Next() {
		var msg Message
		if err = rows.StructScan(&msg); err != nil {
			break
		}
		msg.Content = fromJSON(msg.Content)
		msgs = append(msgs, msg)
	}
	if err == nil {
		err = rows.Err()
	}
	rows.Close()
	return msgs, err
}

// Check if MySQL error is a Error Code: 1062. Duplicate entry
func isDuplicate(err error) bool {
	if err == nil {
		return false
	}

	myerr, ok := err.(*ms.MySQLError)
	return ok && myerr.Number == 1062
}

func fromJSON(src interface{}) interface{} {
	if src == nil {
		return nil
	}
	if bb, ok := src.([]byte); ok {
		var out interface{}
		json.Unmarshal(bb, &out)
		return out
	}
	return nil
}

func toJSON(src interface{}) []byte {
	if src == nil {
		return nil
	}

	jval, _ := json.Marshal(src)
	return jval
}
