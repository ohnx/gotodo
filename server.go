package main

import (
    // standard library
    "net/http"
    "log"
    "os"

    // http router from julienschmidt
    "github.com/julienschmidt/httprouter"

    // own stuff
    "github.com/ohnx/gotodo/endpoints"
    "github.com/ohnx/gotodo/database"
)

func main() {
    // Get database name
    filename := os.Getenv("DB_FILENAME")
    if len(filename) == 0 {
        filename = "data.db"
    }
    log.Printf("Server using database file `%s`", filename)
    database.Connect(filename)
    defer database.Disconnect()

    // Create a new router
    r := httprouter.New()

    // Create endpoints
    tokenEndpoint := endpoints.NewTokenEndpoint()
    //todoEndpoint := endpoints.NewTodoEndpoint()
    todosEndpoint := endpoints.NewTodosEndpoint()
    tagsEndpoint := endpoints.NewTagsEndpoint()

    // Create a handler for endpoints
    r.POST("/api/token/type", tokenEndpoint.Type)
    r.POST("/api/token/new", tokenEndpoint.New)
    r.POST("/api/token/invalidate", tokenEndpoint.Invalidate)
    //
    r.GET("/api/todos/list", todosEndpoint.List)
    r.POST("/api/todos/list", todosEndpoint.List)
    r.POST("/api/tags/list", tagsEndpoint.List)

    // Get the port
    port := os.Getenv("PORT")
    if len(port) == 0 {
        port = "8080"
    }

    // Start server
    log.Printf("Server listening on 0.0.0.0:%s", port)
    err := http.ListenAndServe(":" + port, r)
    if err != nil {
        log.Fatalf("Failed to listen: %s", err)
    }
}
