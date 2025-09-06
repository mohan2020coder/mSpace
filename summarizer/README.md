📝 How to use:

Setup PostgreSQL and run scripts/init_db.sql.

Configure config.yaml with your DB DSN and Ollama endpoint.

go build -o summarizer cmd/main.go

go build -o summarizer.exe cmd/main.go



./summarizer --input path/to/judgment.pdf --query "summary"

./summarizer --input path/to/judgment.pdf --query "What was the court's ruling on liability?"


summarizer.exe --input OP_1_2021.pdf --query "summary"


summarizer.exe --input OP_1_2021.pdf --query ""