package models

import (
    // Standard library
    "log"

    // Own stuff
    "github.com/ohnx/gotodo/database"
)

type (
    // Represent a tag
    Tag struct {
        Id      int     `json:"id"`
        Name    string  `json:"name"`
    }
)

// List all tags in the database
func ListAllTags() []Tag {
    // Get connection handle
    conn := database.GetConnection()

    // prepare read statement
    stmt, err := conn.Prepare("SELECT * FROM tags")
    if err != nil {
        log.Printf("Warning: Failed to read database: %s", err)
        return nil
    }
    defer stmt.Close()

    // Execute read statement
    res, err := stmt.Query()
    if err != nil {
        log.Printf("Warning: Failed to read database: %s", err)
        return nil
    }
    defer res.Close()

    // Create new slice for storing output
    var r []Tag
    var a Tag

    // Check results
    for res.Next() {
        // Read in the values from the database
        err = res.Scan(&a.Id, &a.Name)
        // Check for errors
        if err != nil {
            log.Printf("Warning: Failed to read database: %s", err)
            return r
        }
        // No errors, append to slice
        r = append(r, a)
    }

    // Done
    return r
}
