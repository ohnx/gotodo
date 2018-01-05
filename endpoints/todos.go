package endpoints

import (
    // stdlib
    "fmt"
    "encoding/json"
    "net/http"

    // HTTP router
    "github.com/julienschmidt/httprouter"

    // own stuff
    "github.com/ohnx/gotodo/models"
)

type (
    // TodosEndpoint represents the controller for operating on the Todos resource
    TodosEndpoint struct {}

    // List endpoint
    TodosEndpointListRequest struct {
        Auth    string          `json:"authority"`
    }
    TodosEndpointListResponse struct {
        Todos   []models.Todo   `json:"todos"`
    }
)

func NewTodosEndpoint() *TodosEndpoint {  
    return &TodosEndpoint{}
}

func (te TodosEndpoint) List(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    // Input type
    var telr TodosEndpointListRequest

    // Create a decoder
    decoder := json.NewDecoder(r.Body)
    // Decode into the input type
    decoder.Decode(&telr)

    // The OwnerId of this request
    var ownerId int = -1

    if len(telr.Auth) != 0 {
        // Stub an example token
        token := models.Token{
            Value:  telr.Auth,
        }

        // Read in the values
        token.ReadValues()

        // Check type
        if token.Type == 1 {
            // This token is authorized to view private todos from this owner
            ownerId = token.OwnerId
        }
    }

    // Read all the tags
    resp := TodosEndpointListResponse{
        Todos: models.ListAllTodos(ownerId),
    }

    // Create JSON response
    jresp, _ := json.Marshal(resp)

    // Write OK + payload
    w.WriteHeader(200)
    fmt.Fprintf(w, "%s", jresp)
}
