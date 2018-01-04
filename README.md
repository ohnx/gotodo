# gotodo

Todo web app written in golang.

Static files are in `/static` - they are expected to be served by some server
such as nginx. Requests to `/api` are handled by the go server.

By default, the server stores data in an SQLite database in the current directory.
To change this to use MySQL or similar, modifications to the code are necessary.

## Note

This was my first project using Go! I am open to any criticism of this code; I'd
rather correct my mistakes and C-isms now rather than later :)
