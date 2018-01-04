package models

import (
    // Standard library
    "log"

    // Own stuff
    "github.com/ohnx/gotodo/database"
)

type (
    // Represent a user
    User struct {
        Id      int     `json:"id"`
        Name    string  `json:"name"`
        // unhashed password
        PwdUH   string  `json:"password"`
    }
)

// Fetch user's ID based on username and password
func (user *User) ReadValues() bool {
    // Check that there is an input Value
    if len(user.Name) == 0 || len(user.PwdUH) == 0 {
        return false
    }

    // Get connection handle
    conn := database.GetConnection()

    // prepare read statement
    stmt, err := conn.Prepare("SELECT id FROM users WHERE name = ? AND password = ?")
    if err != nil {
        log.Printf("Warning: Failed to read database: %s", err)
        return false
    }
    defer stmt.Close()

    // Execute read statement
    res, err := stmt.Query(user.Name, database.Hash(user.PwdUH))
    if err != nil {
        log.Printf("Warning: Failed to read database: %s", err)
        return false
    }
    defer res.Close()

    // Check results
    for res.Next() {
        // Only care about the 1st result
        err = res.Scan(&user.Id)
        if err != nil {
            log.Printf("Warning: Failed to read database: %s", err)
            return false
        }
        return true
    }
    // Incorrect username or password
    return false
}
