package models

type (
    // Represent a todo item
    Todo struct {
        Id      int     `json:"id"`
        Name    string  `json:"name"`
        State   int     `json:"state"`
        TagId   int     `json:"tag_id"`
        OwnerId int     `json:"owner_id"`
        Public  bool    `json:"public"`
        Desc    string  `json:"description"`
    }
)
