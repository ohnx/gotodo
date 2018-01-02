# gotodo

Todo web app written in golang.

Static files are in `/static` - they are expected to be served by some server
such as nginx. Requests to `/api` are handled by the go server.

By default, the server stores data in an SQLite database in the current directory.
To change this to use MySQL or similar, modifications to the code are necessary.
