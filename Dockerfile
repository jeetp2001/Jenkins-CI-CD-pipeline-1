from golang:1.18-bullseye
workdir /app
copy go.mod .
run go mod download
copy &.go
run go build -t .
expose 3000
run ls -l
