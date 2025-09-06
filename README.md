go run cmd/api/main.go   # start Gin API
go run cmd/cli/main.go   # start CLI


docker-compose up -d


curl -T OP_1_2025_1.pdf -H "Accept: text/plain" http://localhost:9998/tika


curl -X POST http://localhost:8080/api/items/1/file   -F "file=@OP_1_2025_1.pdf"
