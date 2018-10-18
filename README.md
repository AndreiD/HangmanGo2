# Hangman Go

Simple implementation of Hangman using Go & gRPC.

~~~~
Hangman is a paper and pencil guessing game for two or more players.
One player thinks of a word, phrase or sentence and the other(s) tries to guess it by suggesting letters or numbers

Usage:
  Hangman [command]

Available Commands:
  help        Help about any command
  listGames   Shows your saved games
  newGame     Starts a new game
  resumeGame  Resums a gamewith id.

Flags:
  -h, --help   help for Hangman
  ~~~~

### The "Arhitecture" :)

- gRPC (why: since it doesn't involve any web interface, gRPC is currently the best)
- Server: keeps track of the games
- Client: Cobra CLI
- TSL certificates
- Authentication with username & password

Api is definded in the package "api"


### Build it:

#### You probably want to generate the certificates again:

~~~~
$ openssl genrsa -out server.key 2048
$ openssl req -new -x509 -sha256 -key server.key -out server.crt -days 365
$ openssl req -new -sha256 -key server.key -outserver.csr
$ openssl x509 -req -sha256 -in server.csr -signkey server.key -out server.crt -days 365
~~~~

#### Test it:

~~~~
go test -timeout 30s github.com\AndreiD\HangmanGo2\server -coverprofile=path_to\go-code-cover
~~~~


#### Run it:

~~~~
in the server run:
go build -o server; ./server

in the client
go build -o client; ./server
~~~~

if you're using windows, replace client & server with client.exe & server.exe

### Todo://

- more testing (gomock) & coverage: coverage: 58.3% of statements (server)
- better error handling
- persistance
- benchmarking (maybe user https://github.com/shafreeck/fperf ?)
- integration with CircleCI

### About:

original engine idea & hanging art done by krasi-georgievor miketmoore (not clear)
cobra cli idea by Cityhunteur
custom words are animals, found somewhere on github

### License: MIT
