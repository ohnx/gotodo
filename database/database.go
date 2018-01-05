package database

import (
    // standard library
    "log"
    "os"

    // hashing
    "crypto/sha512"
    "encoding/hex"

    // Database stuff
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

var initStmt = `
CREATE TABLE tokens (
	id integer PRIMARY KEY AUTOINCREMENT,
	type integer,
	value varchar,
	owner_id integer
);

CREATE TABLE users (
	id integer PRIMARY KEY AUTOINCREMENT,
	name varchar,
	password varchar
);

CREATE TABLE todos (
	id integer PRIMARY KEY AUTOINCREMENT,
	state integer,
	tag_id integer,
	owner_id integer,
	public integer,
	name varchar,
	description text
);

CREATE TABLE tags (
	id integer PRIMARY KEY AUTOINCREMENT,
	name varchar
);
`

var (
    // connection handle
    conn *sql.DB
)

func file_exists(filename string) bool {
    _, err := os.Stat(filename)

    if os.IsNotExist(err) {
        return false
    }

    return err == nil
}

func connect(filename string) {
    var err error
    conn, err = sql.Open("sqlite3", filename)
    if err != nil {
        log.Fatalf("Failed to open database: %s", err)
    }
}

func InitializeDatabase(filename string) {
    connect(filename)
    defer Disconnect()
    var err error

    log.Println("Info: Database file not found, creating new database...")

    // Initialize all the tables
    _, err = conn.Exec(initStmt)
    if err != nil {
        log.Fatalf("Failed to initialize database: %s", err)
    }

    // Add admin user
    stmt, err := conn.Prepare("INSERT INTO users(name, password) values(?,?)")
    if err != nil {
        log.Fatalf("Failed to initialize database: %s", err)
    }

    // TODO: configurable
    _, err = stmt.Exec("ohnx", Hash("password"))
    if err != nil {
        log.Fatalf("Failed to initialize database: %s", err)
    }

    // Add initial blank tag
    stmt, err = conn.Prepare("INSERT INTO tags(name) values(?)")
    if err != nil {
        log.Fatalf("Failed to initialize database: %s", err)
    }

    // TODO: configurable
    _, err = stmt.Exec("Unsorted")
    if err != nil {
        log.Fatalf("Failed to initialize database: %s", err)
    }

    // Be responsible!
    stmt.Close()
}

func Hash(str string) string {
    // SHA512 is the password hash being used
    hashalgo := sha512.New()
    hashalgo.Write([]byte(str))
    return hex.EncodeToString(hashalgo.Sum(nil))
}

func Connect(filename string) {
    // Fill the database with blank tables if it doesn't exist
    if !file_exists(filename) {
        InitializeDatabase(filename)
    }

    connect(filename)
}

func GetConnection() *sql.DB {
    return conn
}

func Disconnect() {
    conn.Close()
}
