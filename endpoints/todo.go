package endpoints

import (
    // stdlib
    "fmt"
    //"encoding/json"
    "net/http"

    // HTTP router
    "github.com/julienschmidt/httprouter"

    // own stuff
    //"github.com/ohnx/gotodo/models"
)

type (
    // TodoEndpoint represents the controller for operating on the Todo resource
    TodoEndpoint struct {}

    // List endpoint
)

func NewTodoEndpoint() *TodoEndpoint {  
    return &TodoEndpoint{}
}

func (te TodoEndpoint) List(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    // Write OK + payload
    w.WriteHeader(200)
    fmt.Fprintf(w, "hi")
}
