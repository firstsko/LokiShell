#/bin/bash
GOOS=linux GOARCH=amd64 go build  -o ./dist/linux-amd64/lokishell 
GOOS=windows GOARCH=amd64 go build  -o ./dist/windows-amd64/lokishell.exe
GOOS=darwin GOARCH=amd64 go build  -o ./dist/darwin-amd64/lokishell

