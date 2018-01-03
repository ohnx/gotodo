package models

import (
    // Standard library
    "log"

    // Own stuff
    "github.com/ohnx/gotodo/database"
)

type (
    // Represent a token
    Token struct {
        Id      int
        Type    int
        Value   string
        OwnerId int
    }
)

func (token *Token) InsertValues() int {
    if token.Id > 0 {
        return -2
    }

    // Get connection handle
    conn := database.GetConnection()

    // prepare read statement
    stmt, err := conn.Prepare("INSERT INTO tokens(type, value, owner_id) values(?,?,?)")

    if err != nil {
        log.Printf("Warning: Failed to read database: %s", err)
        return -1
    }
    defer stmt.Close()

    // Execute read statement
    _, err = stmt.Exec(token.Type, token.Value, token.OwnerId)
    if err != nil {
        log.Printf("Warning: Failed to read database: %s", err)
        return -1
    }

    // No error
    return 0
}

func (token *Token) ReadValuesByValue() {
    // Default values
    token.Id = -1
    token.Type = 9
    token.OwnerId = -1

    // Check that there is an input Value
    if len(token.Value) == 0 {
        return
    }

    // Get connection handle
    conn := database.GetConnection()

    // prepare read statement
    stmt, err := conn.Prepare("SELECT * FROM tokens WHERE value = ?")
    if err != nil {
        log.Printf("Warning: Failed to read database: %s", err)
        return
    }
    defer stmt.Close()

    // Execute read statement
    res, err := stmt.Query(token.Value)
    if err != nil {
        log.Printf("Warning: Failed to read database: %s", err)
        return
    }
    defer res.Close()

    // Check results
    for res.Next() {
        err = res.Scan(&token.Id, &token.Type, &token.Value, &token.OwnerId)
        log.Printf("Read database values for token Id %d", token.Id)
        if err != nil {
            log.Printf("Warning: Failed to read database: %s", err)
            return
        }
        // Only read 1st result... there shouldn't be any more...
        return
    }
    log.Printf("No results found...")
}
