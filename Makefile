build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -a -installsuffix cgo -o ./out/app cmd/service/main.go
