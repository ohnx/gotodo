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

// Ref: https://willnorris.com/2014/05/go-rest-apis-and-pointers
type (
    // TokenEndpoint represents the controller for operating on the Token resource
    TokenEndpoint struct {}

    // Type endpoint
    TokenEndpointTypeRequest struct {
        Token   string      `json:"token"`
    }
    TokenEndpointTypeResponse struct {
        Type    int         `json:"type"`
    }

    // New endpoint
    TokenEndpointNewRequest struct {
        Type    int         `json:"type"`
        UName   *string     `json:"username"`
        UPwdUH  *string     `json:"password"`
        Auth    *string     `json:"authority"`
    }
    TokenEndpointNewResponse struct {
        Error   string      `json:"error,omitempty"`
        Token   string      `json:"token,omitempty"`
    }

    // Invalidate endpoint
    TokenEndpointInvalidateRequest struct {
        Token   string      `json:"token"`
        UName   *string     `json:"username"`
        UPwdUH  *string     `json:"password"`
        Auth    *string     `json:"authority"`
    }
    TokenEndpointInvalidateResponse struct {
        Error   string      `json:"error,omitempty"`
    }
)

func NewTokenEndpoint() *TokenEndpoint {  
    return &TokenEndpoint{}
}

func (te TokenEndpoint) Type(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    // Input type
    var tetr TokenEndpointTypeRequest

    // Create a decoder
    decoder := json.NewDecoder(r.Body)
    // Decode into the input type
    err := decoder.Decode(&tetr)

    // Check for errors
    if err != nil || len(tetr.Token) == 0 {
        // Failed to parse user input... call it a user error
        w.WriteHeader(400)
        return
    }

    // Stub an example token
    token := models.Token{
        Value:  tetr.Token,
    }

    // Read in the values
    token.ReadValues()

    // Create response
    resp := TokenEndpointTypeResponse{
        Type:   token.Type,
    }

    // Create JSON response
    jresp, _ := json.Marshal(resp)

    // Write OK + payload
    w.WriteHeader(200)
    fmt.Fprintf(w, "%s", jresp)
}

func (te TokenEndpoint) New(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    // Input type
    var tenr TokenEndpointNewRequest

    // Create a decoder
    decoder := json.NewDecoder(r.Body)
    // Decode into the input type
    err := decoder.Decode(&tenr)

    // Check for errors
    if err != nil {
        // Failed to parse user input... call it a user error
        w.WriteHeader(400)
        return
    }

    // Check that the user provided sufficient authorization
    // either username + password OR auth
    if !((tenr.UName != nil && tenr.UPwdUH != nil) || tenr.Auth != nil) {
        // insufficient auth provided
        resp := TokenEndpointInvalidateResponse{
            Error: "Insufficient authorization provided",
        }
        jresp, _ := json.Marshal(resp)

        // Write error + payload
        w.WriteHeader(400)
        fmt.Fprintf(w, "%s", jresp)
        return
    }

    // 2 code paths - user or token specified
    if tenr.UName != nil && tenr.UPwdUH != nil {
        // Username and password specified, first check if they are valid
        user := models.User{
            Name:   *tenr.UName,
            PwdUH:  *tenr.UPwdUH,
        }

        if !user.ReadValues() {
            // User not authorized
            resp := TokenEndpointNewResponse{
                Error: "Invalid username and password combination",
            }
            jresp, _ := json.Marshal(resp)

            // Write error + payload
            w.WriteHeader(403)
            fmt.Fprintf(w, "%s", jresp)
            return
        }

        // Check for valid token type
        if tenr.Type < 1 || tenr.Type > 3 {
            // Invalid kind of token
            resp := TokenEndpointNewResponse{
                Error: "Invalid requested token type",
            }
            jresp, _ := json.Marshal(resp)

            // Write error + payload
            w.WriteHeader(400)
            fmt.Fprintf(w, "%s", jresp)
            return
        }

        // User is authorized, create a new token!
        token := models.Token{
            Type:   tenr.Type,
            OwnerId:user.Id,
        }
        token.GenValue()

        // Write to database
        if !token.InsertValues() {
            // Database error
            resp := TokenEndpointNewResponse{
                Error: "Database error",
            }
            jresp, _ := json.Marshal(resp)

            // Write error + payload
            w.WriteHeader(500)
            fmt.Fprintf(w, "%s", jresp)
            return
        }

        // Done!
        resp := TokenEndpointNewResponse{
            Token: token.Value,
        }
        jresp, _ := json.Marshal(resp)

        // Write error + payload
        w.WriteHeader(200)
        fmt.Fprintf(w, "%s", jresp)
        return
    } else if tenr.Auth != nil {
        // Token specified, check if it is valid
        auth := models.Token{
            Value:  *tenr.Auth,
        }

        // Check the privileges on the auth token
        if !auth.ReadValues() || auth.Type != 1 {
            // User not authorized
            resp := TokenEndpointNewResponse{
                Error: "Authorization token lacks creation privilege",
            }
            jresp, _ := json.Marshal(resp)

            // Write error + payload
            w.WriteHeader(403)
            fmt.Fprintf(w, "%s", jresp)
            return
        }

        // Check for valid token type
        if tenr.Type < 2 || tenr.Type > 3 {
            // Unrecognized token type or invalid type
            resp := TokenEndpointNewResponse{
                Error: "Invalid requested token type",
            }
            jresp, _ := json.Marshal(resp)

            // Write error + payload
            w.WriteHeader(400)
            fmt.Fprintf(w, "%s", jresp)
            return
        }

        // Token is authorized, create a new token!
        token := models.Token{
            Type:   tenr.Type,
            OwnerId:auth.OwnerId,
        }
        token.GenValue()

        // Write to database
        if !token.InsertValues() {
            // Database error
            resp := TokenEndpointNewResponse{
                Error: "Database error",
            }
            jresp, _ := json.Marshal(resp)

            // Write error + payload
            w.WriteHeader(500)
            fmt.Fprintf(w, "%s", jresp)
            return
        }

        // Done!
        resp := TokenEndpointNewResponse{
            Token: token.Value,
        }
        jresp, _ := json.Marshal(resp)

        // Write error + payload
        w.WriteHeader(200)
        fmt.Fprintf(w, "%s", jresp)
        return
    }
}

func (te TokenEndpoint) Invalidate(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    // Input type
    var teir TokenEndpointInvalidateRequest

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

    // Check that the user provided sufficient authorization
    // either username + password OR auth
    if !((teir.UName != nil && teir.UPwdUH != nil) || teir.Auth != nil) {
        // insufficient auth provided
        resp := TokenEndpointInvalidateResponse{
            Error: "Insufficient authorization provided",
        }
        jresp, _ := json.Marshal(resp)

        // Write error + payload
        w.WriteHeader(400)
        fmt.Fprintf(w, "%s", jresp)
        return
    }

    // First fetch this token's Id
    token := models.Token{
        Value: teir.Token,
    }

    // Try reading in the token's values
    if !token.ReadValues() {
        // Invalid token
        resp := TokenEndpointInvalidateResponse{
            Error: "Invalid field token",
        }
        jresp, _ := json.Marshal(resp)

        // Write error + payload
        w.WriteHeader(400)
        fmt.Fprintf(w, "%s", jresp)
        return
    }

    if teir.UName != nil && teir.UPwdUH != nil {
        // Username and password specified, first check if they are valid
        user := models.User{
            Name:   *teir.UName,
            PwdUH:  *teir.UPwdUH,
        }

        if !user.ReadValues() {
            // User not authorized
            resp := TokenEndpointInvalidateResponse{
                Error: "Invalid username and password combination",
            }
            jresp, _ := json.Marshal(resp)

            // Write error + payload
            w.WriteHeader(403)
            fmt.Fprintf(w, "%s", jresp)
            return
        }

        if user.Id != token.OwnerId {
            // Owner Id doesn't match
            resp := TokenEndpointInvalidateResponse{
                Error: "User does not own token",
            }
            jresp, _ := json.Marshal(resp)

            // Write error + payload
            w.WriteHeader(400)
            fmt.Fprintf(w, "%s", jresp)
            return
        }

        // No errors! ready to delete
        if !token.Remove() {
            // Database error
            resp := TokenEndpointInvalidateResponse{
                Error: "Database error",
            }
            jresp, _ := json.Marshal(resp)

            // Write error + payload
            w.WriteHeader(500)
            fmt.Fprintf(w, "%s", jresp)
            return
        }
        resp := TokenEndpointInvalidateResponse{}
        jresp, _ := json.Marshal(resp)

        // Write error + payload
        w.WriteHeader(200)
        fmt.Fprintf(w, "%s", jresp)
        return
    } else if teir.Auth != nil {
        // Now check if this token can be invalidated using the given auth
        if token.Type == 1 {
            // Needs to be either itself or a username/password combination
            if token.Value == *teir.Auth {
                // Authorized delete, the token is itself
                if !token.Remove() {
                    // Database error
                    resp := TokenEndpointInvalidateResponse{
                        Error: "Database error",
                    }
                    jresp, _ := json.Marshal(resp)

                    // Write error + payload
                    w.WriteHeader(500)
                    fmt.Fprintf(w, "%s", jresp)
                    return
                }
                resp := TokenEndpointInvalidateResponse{}
                jresp, _ := json.Marshal(resp)

                // Write error + payload
                w.WriteHeader(200)
                fmt.Fprintf(w, "%s", jresp)
                return
            } else {
                // Insufficient authorization
                resp := TokenEndpointInvalidateResponse{
                    Error: "Insufficient authorization provided",
                }
                jresp, _ := json.Marshal(resp)

                // Write error + payload
                w.WriteHeader(403)
                fmt.Fprintf(w, "%s", jresp)
                return
            }
        } else {
            // Type 2 or 3 - can be invalidated by any master token, so we check if auth is a master token
            auth := models.Token{
                Value: *teir.Auth,
            }

            // Check the privileges on the auth token
            if !auth.ReadValues() || auth.Type != 1 {
                // User not authorized
                resp := TokenEndpointInvalidateResponse{
                    Error: "Authorization token lacks removal privilege",
                }
                jresp, _ := json.Marshal(resp)

                // Write error + payload
                w.WriteHeader(403)
                fmt.Fprintf(w, "%s", jresp)
                return
            }

            if auth.OwnerId != token.OwnerId {
                // Owner Id doesn't match
                resp := TokenEndpointInvalidateResponse{
                    Error: "User does not own token",
                }
                jresp, _ := json.Marshal(resp)

                // Write error + payload
                w.WriteHeader(400)
                fmt.Fprintf(w, "%s", jresp)
                return
            }

            // Token is valid, time to invalidate!
            if !token.Remove() {
                // Database error
                resp := TokenEndpointInvalidateResponse{
                    Error: "Database error",
                }
                jresp, _ := json.Marshal(resp)

                // Write error + payload
                w.WriteHeader(500)
                fmt.Fprintf(w, "%s", jresp)
                return
            }
            resp := TokenEndpointInvalidateResponse{}
            jresp, _ := json.Marshal(resp)

            // Write error + payload
            w.WriteHeader(200)
            fmt.Fprintf(w, "%s", jresp)
                return
        }
    }
}
