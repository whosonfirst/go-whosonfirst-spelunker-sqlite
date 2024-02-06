CWD=$(shell pwd)

GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")

# SPELUNKER_URI=sql://sqlite3?vfs=file:///usr/local/data/ca-search.db
SPELUNKER_URI=sql://sqlite3?vfs=http://localhost:8081/ca-search.db

server:
	go run -mod $(GOMOD) -tags "icu json1 fts5" cmd/httpd/main.go \
		-server-uri http://localhost:8080 \
		-spelunker-uri $(SPELUNKER_URI)
