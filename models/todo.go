package models

import (
    // Standard library
    "log"
    "time"

    // Own stuff
    "github.com/ohnx/gotodo/database"
)

type (
    // Represent a todo item
    Todo struct {
        Id      int         `json:"id"`
        State   int         `json:"state"`
        TagId   int         `json:"tag_id"`
        OwnerId int         `json:"owner_id"`
        Public  bool        `json:"public"`
        Name    string      `json:"name"`
        DueDate time.Time   `json:"due_date"`
        Desc    string      `json:"description"`
    }
)

// Inserts a new todo. Returns true on success, false on error.
func (todo *Todo) InsertValues() bool {
    // Check that there is no input Id
    if todo.Id > 0 || len(todo.Name) == 0{
        return false
    }

    // Get connection handle
    conn := database.GetConnection()

    // prepare insert statement
    stmt, err := conn.Prepare("INSERT INTO todos(state, tag_id, owner_id, public, name, duedate, description) values(?,?,?,?,?,?,?)")

    if err != nil {
        log.Printf("Warning: Failed to write to database: %s", err)
        return false
    }
    defer stmt.Close()

    // Execute insert statement
    _, err = stmt.Exec(todo.State, todo.TagId, todo.OwnerId, todo.Public, todo.Name, todo.DueDate, todo.Desc)
    if err != nil {
        log.Printf("Warning: Failed to write to database: %s", err)
        return false
    }

    // No error
    return true
}

// Updates an existing todo. Returns true on success, false on error.
func (todo *Todo) WriteValues() bool {
    // Check that there is an input Id
    if todo.Id <= 0 || len(todo.Name) == 0 {
        return false
    }

    // Get connection handle
    conn := database.GetConnection()

    // prepare insert statement
    stmt, err := conn.Prepare("UPDATE todos SET state = ?, tag_id = ?, public = ?, name = ?, duedate = ?, description = ? WHERE id = ?")

    if err != nil {
        log.Printf("Warning: Failed to write to database: %s", err)
        return false
    }
    defer stmt.Close()

    // Execute insert statement
    _, err = stmt.Exec(todo.State, todo.TagId, todo.Public, todo.Name, todo.DueDate, todo.Desc, todo.Id)
    if err != nil {
        log.Printf("Warning: Failed to write to database: %s", err)
        return false
    }

    // No error
    return true
}

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
        err = res.Scan(&todo.Id, &todo.State, &todo.TagId, &todo.OwnerId, &boolConv, &todo.Name, &todo.DueDate, &todo.Desc)
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

// Read in the only the owner and publicness of a todo id. Returns true if values were read.
func (todo *Todo) ReadPermissions() bool {
    // Check that there is an input Id
    if todo.Id < 1 {
        return false
    }

    // Get connection handle
    conn := database.GetConnection()

    // prepare read statement
    stmt, err := conn.Prepare("SELECT owner_id, public FROM todos WHERE id = ?")
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
        // Only care about the 1st result
        var boolConv int
        err = res.Scan(&todo.OwnerId, &boolConv)
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

// Remove a todo from the database based on Id. Returns true on successful removal.
func (todo *Todo) Remove() bool {
    // Check that there is an input Id
    if todo.Id <= 0 {
        return false
    }

    // Get connection handle
    conn := database.GetConnection()

    // prepare delete statement
    stmt, err := conn.Prepare("DELETE FROM todos WHERE id = ?")
    if err != nil {
        log.Printf("Warning: Failed to write to database: %s", err)
        return false
    }
    defer stmt.Close()

    // Execute delete statement
    _, err = stmt.Exec(todo.Id)
    if err != nil {
        log.Printf("Warning: Failed to write to database: %s", err)
        return false
    }

    return true
}

// List all viewable todos to the owner_id in the database
func ListAllTodos(owner_id int) []Todo {
    // Get connection handle
    conn := database.GetConnection()

    // prepare read statement
    stmt, err := conn.Prepare("SELECT id,state,tag_id,name,duedate FROM todos WHERE public = 1 OR owner_id = ?")
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
        err = res.Scan(&todo.Id, &todo.State, &todo.TagId, &todo.Name, &todo.DueDate)
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
