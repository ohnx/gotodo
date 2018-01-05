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
    // TagsEndpoint represents the controller for operating on the Tags resource
    TagsEndpoint struct {}

    // List endpoint
    TagsEndpointListRequest struct {
        Auth    string          `json:"authority"`
    }
    TagsEndpointListResponse struct {
        Error   string          `json:"error,omitempty"`
        Tags    []models.Tag    `json:"tags,omitempty"`
    }
)

func NewTagsEndpoint() *TagsEndpoint {  
    return &TagsEndpoint{}
}

func (te TagsEndpoint) List(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    // Input type
    var telr TagsEndpointListRequest

    // Create a decoder
    decoder := json.NewDecoder(r.Body)
    // Decode into the input type
    err := decoder.Decode(&telr)

    // Check for errors
    if err != nil {
        // Failed to parse user input... call it a user error
        w.WriteHeader(400)
        return
    }

    // Stub an example token
    token := models.Token{
        Value:  telr.Auth,
    }

    // Read in the values
    token.ReadValues()

    // Check type
    if token.Type != 1 {
        // Not enough permissions to list tags
        resp := TagsEndpointListResponse{
            Error: "Insufficient authorization provided",
        }
        jresp, _ := json.Marshal(resp)

        // Write error + payload
        w.WriteHeader(400)
        fmt.Fprintf(w, "%s", jresp)
        return
    }

    // Read all the tags
    resp := TagsEndpointListResponse{
        Tags: models.ListAllTags(),
    }

    // Create JSON response
    jresp, _ := json.Marshal(resp)

    // Write OK + payload
    w.WriteHeader(200)
    fmt.Fprintf(w, "%s", jresp)
}
