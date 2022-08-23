package models

import (
    // Standard library
    "log"
    "math/rand"
    "time"

    // Own stuff
    "github.com/ohnx/gotodo/database"
)

type (
    // Represent a token
    Token struct {
        Id      int     `json:"id"`
        // TODO: I have no clue if this is correct or not
        // Token type 1: Master token (create tokens?)
        // Token type 2: Read, write, edit todos
        // Token type 3: Create-only
        // Token type 4: [WIP] Read-only
        Type    int     `json:"type"`
        Value   string  `json:"value"`
        OwnerId int     `json:"owner_id"`
    }
)

// Inserts a new token. Returns true on success, false on error.
func (token *Token) InsertValues() bool {
    // Check that there is no input Id
    if token.Id > 0 {
        return false
    }

    // Get connection handle
    conn := database.GetConnection()

    // prepare insert statement
    stmt, err := conn.Prepare("INSERT INTO tokens(type, value, owner_id) values(?,?,?)")

    if err != nil {
        log.Printf("Warning: Failed to write to database: %s", err)
        return false
    }
    defer stmt.Close()

    // Execute insert statement
    _, err = stmt.Exec(token.Type, token.Value, token.OwnerId)
    if err != nil {
        log.Printf("Warning: Failed to write to database: %s", err)
        return false
    }

    // No error
    log.Printf("Info: Created new token %s type %d", token.Value, token.Type)
    return true
}

// Read in the values of a token based on value. Returns true if values were read
func (token *Token) ReadValues() bool {
    // Check that there is an input Value
    if len(token.Value) == 0 {
        return false
    }

    // Get connection handle
    conn := database.GetConnection()

    // prepare read statement
    stmt, err := conn.Prepare("SELECT * FROM tokens WHERE value = ?")
    if err != nil {
        log.Printf("Warning: Failed to read database: %s", err)
        return false
    }
    defer stmt.Close()

    // Execute read statement
    res, err := stmt.Query(token.Value)
    if err != nil {
        log.Printf("Warning: Failed to read database: %s", err)
        return false
    }
    defer res.Close()

    // Check results
    for res.Next() {
        err = res.Scan(&token.Id, &token.Type, &token.Value, &token.OwnerId)
        if err != nil {
            log.Printf("Warning: Failed to read database: %s", err)
            return false
        }
        // Only read 1st result... there shouldn't be any more...
        return true
    }

    // Invalid token
    token.Type = 9
    return false
}

// Remove a token from the database based on Id. Returns true on successful removal.
func (token *Token) Remove() bool {
    // Check that there is an input Id
    if token.Id <= 0 {
        return false
    }

    // Get connection handle
    conn := database.GetConnection()

    // prepare delete statement
    stmt, err := conn.Prepare("DELETE FROM tokens WHERE id = ?")
    if err != nil {
        log.Printf("Warning: Failed to write to database: %s", err)
        return false
    }
    defer stmt.Close()

    // Execute delete statement
    _, err = stmt.Exec(token.Id)
    if err != nil {
        log.Printf("Warning: Failed to write to database: %s", err)
        return false
    }

    log.Printf("Info: Removed token #%d = %s from database", token.Id, token.Value)
    return true
}

// Adapted from StackOverflow answer https://stackoverflow.com/questions/22892120
const tLen = 32
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
    letterIdxBits = 6                    // 6 bits to represent a letter index
    letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
    letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)
var src = rand.NewSource(time.Now().UnixNano())
func (token *Token) GenValue() {
    b := make([]byte, tLen)
    // A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
    for i, cache, remain := tLen-1, src.Int63(), letterIdxMax; i >= 0; {
        if remain == 0 {
            cache, remain = src.Int63(), letterIdxMax
        }
        if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
            b[i] = letterBytes[idx]
            i--
        }
        cache >>= letterIdxBits
        remain--
    }

    token.Value = string(b)
}
