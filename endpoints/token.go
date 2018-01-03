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
    // TokenEndpoint represents the controller for operating on the User resource
    TokenEndpoint struct{}
)

func NewTokenEndpoint() *TokenEndpoint {  
    return &TokenEndpoint{}
}

func (te TokenEndpoint) Type(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    // Stub an example user
    u := models.Token{
        Id:     0,
        Type:   0,
        Value:  "0xdeadbeef",
        OwnerId:0,
    }

    // Read in the values
    u.ReadValuesByValue()

    // Marshal provided interface into JSON structure
    uj, _ := json.Marshal(u)

    // Write content-type, statuscode, payload
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(200)
    fmt.Fprintf(w, "%s", uj)
}

func (te TokenEndpoint) New(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    // Stub an example user
    u := models.Token{
        Id:     0,
        Type:   1,
        Value:  "0xdeadbeef",
        OwnerId:1337,
    }

    u.InsertValues()

    // Marshal provided interface into JSON structure
    uj, _ := json.Marshal(u)

    // Write content-type, statuscode, payload
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(200)
    fmt.Fprintf(w, "%s", uj)
}

func (te TokenEndpoint) Invalidate(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    // Stub an example user
    u := models.Token{
        Id:     0,
        Type:   1,
        Value:  "Invalidate",
        OwnerId:0,
    }

    // Marshal provided interface into JSON structure
    uj, _ := json.Marshal(u)

    // Write content-type, statuscode, payload
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(200)
    fmt.Fprintf(w, "%s", uj)
}
