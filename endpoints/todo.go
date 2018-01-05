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
    // TodoEndpoint represents the controller for operating on the Todo resource
    TodoEndpoint struct {}

    // Update endpoint
    TodoEndpointUpdateRequest struct {
        Todo    models.Todo     `json:"todo"`
        Auth    string          `json:"authority"`
    }
    TodoEndpointUpdateResponse struct {
        Error   string          `json:"error,omitempty"`
    }

    // Remove endpoint
    TodoEndpointRemoveRequest struct {
        Todo    models.Todo     `json:"todo"`
        Auth    string          `json:"authority"`
    }
    TodoEndpointRemoveResponse struct {
        Error   string          `json:"error,omitempty"`
    }

    // Info endpoint
    TodoEndpointInfoRequest struct {
        Todo    models.Todo     `json:"todo"`
        Auth    string          `json:"authority"`
    }
    TodoEndpointInfoResponse struct {
        Error   string          `json:"error,omitempty"`
        Todo    models.Todo     `json:"todo,omitempty"`
    }
)

func NewTodoEndpoint() *TodoEndpoint {  
    return &TodoEndpoint{}
}

func (te TodoEndpoint) Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    // Input type
    var teur TodoEndpointUpdateRequest

    // Create a decoder
    decoder := json.NewDecoder(r.Body)
    // Decode into the input type
    err := decoder.Decode(&teur)

    // Check for errors
    if err != nil {
        // Failed to parse user input... call it a user error
        w.WriteHeader(400)
        return
    }

    // Two possiblities - new or update existing
    if teur.Todo.Id < 0 {
        // New, check if auth token <= 3
        // Token specified, check if it is valid
        auth := models.Token{
            Value:  teur.Auth,
        }

        // Check the privileges on the auth token
        if !auth.ReadValues() || auth.Type > 3 {
            // User not authorized
            resp := TodoEndpointUpdateResponse{
                Error: "Authorization token lacks creation privilege",
            }
            jresp, _ := json.Marshal(resp)

            // Write error + payload
            w.WriteHeader(403)
            fmt.Fprintf(w, "%s", jresp)
            return
        }

        // Authorized to create todo, so write it!
        teur.Todo.OwnerId = auth.OwnerId
        if !teur.Todo.InsertValues() {
            // Database error
            resp := TodoEndpointUpdateResponse{
                Error: "Database error",
            }
            jresp, _ := json.Marshal(resp)

            // Write error + payload
            w.WriteHeader(500)
            fmt.Fprintf(w, "%s", jresp)
            return
        }

        // Everything is good!
        resp := TodoEndpointUpdateResponse{}
        jresp, _ := json.Marshal(resp)

        // Write error + payload
        w.WriteHeader(200)
        fmt.Fprintf(w, "%s", jresp)
        return
    } else {
        // Existing todo, check if auth token <= 2
        // Token specified, check if it is valid
        auth := models.Token{
            Value:  teur.Auth,
        }

        // Check the privileges on the auth token
        if !auth.ReadValues() || auth.Type > 2 {
            // User not authorized
            resp := TodoEndpointUpdateResponse{
                Error: "Authorization token lacks modification privilege",
            }
            jresp, _ := json.Marshal(resp)

            // Write error + payload
            w.WriteHeader(403)
            fmt.Fprintf(w, "%s", jresp)
            return
        }

        // Valid key, but does it belong to the right owner?
        todo := models.Todo{
            Id:     teur.Todo.Id,
        }
        if !todo.ReadPermissions() {
            // Database error
            resp := TodoEndpointUpdateResponse{
                Error: "Todo not found in database",
            }
            jresp, _ := json.Marshal(resp)

            // Write error + payload
            w.WriteHeader(400)
            fmt.Fprintf(w, "%s", jresp)
            return
        }
        if todo.OwnerId != auth.OwnerId {
            // Todo doesn't belong to the right owner
            resp := TodoEndpointUpdateResponse{
                Error: "User does not own todo",
            }
            jresp, _ := json.Marshal(resp)

            // Write error + payload
            w.WriteHeader(403)
            fmt.Fprintf(w, "%s", jresp)
            return
        }

        // Everything looks good! Time to update the todo
        if !teur.Todo.WriteValues() {
            // Database error
            resp := TodoEndpointUpdateResponse{
                Error: "Database error",
            }
            jresp, _ := json.Marshal(resp)

            // Write error + payload
            w.WriteHeader(500)
            fmt.Fprintf(w, "%s", jresp)
            return
        }

        // Everything is good!
        resp := TodoEndpointUpdateResponse{}
        jresp, _ := json.Marshal(resp)

        // Write error + payload
        w.WriteHeader(200)
        fmt.Fprintf(w, "%s", jresp)
        return
    }
}

func (te TodoEndpoint) Remove(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    // Input type
    var terr TodoEndpointRemoveRequest

    // Create a decoder
    decoder := json.NewDecoder(r.Body)
    // Decode into the input type
    err := decoder.Decode(&terr)

    // Check for errors
    if err != nil || len(terr.Auth) == 0 {
        // Failed to parse user input... call it a user error
        w.WriteHeader(400)
        return
    }

    // Stub an example token
    auth := models.Token{
        Value:  terr.Auth,
    }

    // Check the privileges on the auth token
    if !auth.ReadValues() || auth.Type != 1 {
        // User not authorized
        resp := TodoEndpointRemoveResponse{
            Error: "Authorization token lacks removal privilege",
        }
        jresp, _ := json.Marshal(resp)

        // Write error + payload
        w.WriteHeader(403)
        fmt.Fprintf(w, "%s", jresp)
        return
    }

    // Check if the todo is owned by the correct person
    if !terr.Todo.ReadPermissions() {
        // Database error
        resp := TodoEndpointRemoveResponse{
            Error: "Todo not found in database",
        }
        jresp, _ := json.Marshal(resp)

        // Write error + payload
        w.WriteHeader(400)
        fmt.Fprintf(w, "%s", jresp)
        return
    }
    if terr.Todo.OwnerId != auth.OwnerId {
        // Todo doesn't belong to the right owner
        resp := TodoEndpointRemoveResponse{
            Error: "User does not own todo",
        }
        jresp, _ := json.Marshal(resp)

        // Write error + payload
        w.WriteHeader(403)
        fmt.Fprintf(w, "%s", jresp)
        return
    }

    // Ok, looks like we can remove the todo now.
    if !terr.Todo.Remove() {
        // Database error
        resp := TodoEndpointRemoveResponse{
            Error: "Database error",
        }
        jresp, _ := json.Marshal(resp)

        // Write error + payload
        w.WriteHeader(500)
        fmt.Fprintf(w, "%s", jresp)
        return
    }

    // Create response
    resp := TodoEndpointRemoveResponse{}

    // Create JSON response
    jresp, _ := json.Marshal(resp)

    // Write OK + payload
    w.WriteHeader(200)
    fmt.Fprintf(w, "%s", jresp)
}

func (te TodoEndpoint) Info(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    // Input type
    var teir TodoEndpointInfoRequest

    // Create a decoder
    decoder := json.NewDecoder(r.Body)
    // Decode into the input type
    err := decoder.Decode(&teir)

    // Check for errors
    if err != nil {
        // Failed to parse user input... call it a user error
        w.WriteHeader(400)
        return
    }

    // Read the permissions first to decide what to do
    if !teir.Todo.ReadPermissions() {
        // Database error
        resp := TodoEndpointUpdateResponse{
            Error: "Todo not found in database",
        }
        jresp, _ := json.Marshal(resp)

        // Write error + payload
        w.WriteHeader(400)
        fmt.Fprintf(w, "%s", jresp)
        return
    }

    if !teir.Todo.Public {
        // Todo is not public, need to check a token
        if len(teir.Auth) == 0 {
            // Token not specified, fake a not known error
            resp := TodoEndpointUpdateResponse{
                Error: "Todo not found in database",
            }
            jresp, _ := json.Marshal(resp)

            // Write error + payload
            w.WriteHeader(400)
            fmt.Fprintf(w, "%s", jresp)
            return
        }

        auth := models.Token{
            Value:  teir.Auth,
        }

        // Check the privileges on the auth token
        if !auth.ReadValues() || auth.Type > 2 {
            // User not authorized
            resp := TodoEndpointInfoResponse{
                Error: "Authorization token lacks information privilege",
            }
            jresp, _ := json.Marshal(resp)

            // Write error + payload
            w.WriteHeader(403)
            fmt.Fprintf(w, "%s", jresp)
            return
        }

        // Check if the todo is owned by the correct person
        if teir.Todo.OwnerId != auth.OwnerId  {
            // Todo doesn't belong to the right owner
            resp := TodoEndpointUpdateResponse{
                Error: "User does not own todo",
            }
            jresp, _ := json.Marshal(resp)

            // Write error + payload
            w.WriteHeader(403)
            fmt.Fprintf(w, "%s", jresp)
            return
        }
    }

    // All good, ready to send information now
    if !teir.Todo.ReadValues() {
        // Database error
        resp := TodoEndpointUpdateResponse{
            Error: "Database error",
        }
        jresp, _ := json.Marshal(resp)

        // Write error + payload
        w.WriteHeader(500)
        fmt.Fprintf(w, "%s", jresp)
        return
    }

    // Create response
    resp := TodoEndpointInfoResponse{
        Todo:   teir.Todo,
    }

    // Create JSON response
    jresp, _ := json.Marshal(resp)

    // Write OK + payload
    w.WriteHeader(200)
    fmt.Fprintf(w, "%s", jresp)
}
