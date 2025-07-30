GOOS=linux GOARCH=amd64 go build -o ./dist/linux-amd64/loki-shell 
scp dist/linux-amd64/loki-shell docker@10.64.0.140:/app/docker/loki/loki-shell-server
