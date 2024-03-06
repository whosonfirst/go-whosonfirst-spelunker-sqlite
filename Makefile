CWD=$(shell pwd)
GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

# SPELUNKER_URI=sql://sqlite3?vfs=file:///usr/local/data/ca-search.db
SPELUNKER_URI=sql://sqlite3?vfs=http://localhost:8081/ca-search.db

server:
	go run -mod $(GOMOD) -tags "icu json1 fts5" cmd/httpd/main.go \
		-server-uri http://localhost:8080 \
		-spelunker-uri $(SPELUNKER_URI)

lambda:
	@make lambda-server

lambda-server:
	if test -f bootstrap; then rm -f bootstrap; fi
	if test -f server.zip; then rm -f server.zip; fi
	GOOS=linux \
		GOARCH=arm64 \
		CGO_ENABLED=1 \
		CC='zig cc -target aarch64-linux-musl' \
		CXX='zig cc -target aarch64-linux-musl' \
		CGO_CFLAGS=-D_LARGEFILE64_SOURCE \
		go build -mod $(GOMOD) \
		-ldflags="$(LDFLAGS) -linkmode external -extldflags -static" \
		-tags "netgo,sqlite_omit_load_extension,fts5,lambda.norpc" \
		-o bootstrap \
		cmd/httpd/main.go
	zip server.zip bootstrap
	rm -f bootstrap	


