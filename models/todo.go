package models

import (
    // Standard library
    "log"

    // Own stuff
    "github.com/ohnx/gotodo/database"
)

type (
    // Represent a todo item
    Todo struct {
        Id      int     `json:"id"`
        State   int     `json:"state"`
        TagId   int     `json:"tag_id"`
        OwnerId int     `json:"owner_id"`
        Public  bool    `json:"public"`
        Name    string  `json:"name"`
        Desc    string  `json:"description"`
    }
)

// Read in the values of a todo based on id. Returns true if values were read.
func (todo *Todo) ReadValues() bool {
    // Check that there is an input Id
    if todo.Id < 1 {
        return false
    }

    // Get connection handle
    conn := database.GetConnection()

    // prepare read statement
    stmt, err := conn.Prepare("SELECT * FROM todos WHERE id = ?")
    if err != nil {
        log.Printf("Warning: Failed to read database: %s", err)
        return false
    }
    defer stmt.Close()

    // Execute read statement
    res, err := stmt.Query(todo.Id)
    if err != nil {
        log.Printf("Warning: Failed to read database: %s", err)
        return false
    }
    defer res.Close()

    // Check results
    for res.Next() {
        var boolConv int
        // Only care about the 1st result
        err = res.Scan(&todo.Id, &todo.State, &todo.TagId, &todo.OwnerId, &boolConv, &todo.Name, &todo.Desc)
        if err != nil {
            log.Printf("Warning: Failed to read database: %s", err)
            return false
        }
        if boolConv == 1 {
            todo.Public = true
        } else {
            todo.Public = false
        }
        return true
    }

    // No results
    return false
}

// List all viewable todos to the owner_id in the database
func ListAllTodos(owner_id int) []Todo {
    // Get connection handle
    conn := database.GetConnection()

    // prepare read statement
    stmt, err := conn.Prepare("SELECT id,state,tag_id,name FROM todos WHERE public = 1 OR owner_id = ?")
    if err != nil {
        log.Printf("Warning: Failed to read database: %s", err)
        return nil
    }
    defer stmt.Close()

    // Execute read statement
    res, err := stmt.Query(owner_id)
    if err != nil {
        log.Printf("Warning: Failed to read database: %s", err)
        return nil
    }
    defer res.Close()

    // Create new slice for storing output
    var r []Todo
    var todo Todo

    // Check results
    for res.Next() {
        // Read in the values from the database
        var boolConv int
        err = res.Scan(&todo.Id, &todo.State, &todo.TagId, &todo.OwnerId, &boolConv, &todo.Name, &todo.Desc)
        if err != nil {
            log.Printf("Warning: Failed to read database: %s", err)
            return nil
        }

        if boolConv == 1 {
            todo.Public = true
        } else {
            todo.Public = false
        }

        // No errors, append to slice
        r = append(r, todo)
    }

    // Done
    return r
}
