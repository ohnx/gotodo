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
    TagsEndpointListResponse struct {
        Tags    []models.Tag    `json:"tags"`
    }
)

func NewTagsEndpoint() *TagsEndpoint {  
    return &TagsEndpoint{}
}

func (te TagsEndpoint) List(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    // No input

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
