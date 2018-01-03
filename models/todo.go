package models

type (
    // Represent a todo item
    Todo struct {
        Id      int
        Name    string
        State   int
        TagId   int
        OwnerId int
        Public  bool
        Desc    string
    }
)
