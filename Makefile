CWD=$(shell pwd)
GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

VFS_URI=http://localhost:8081/ca-search.db

server:
	go run -mod $(GOMOD) -tags "icu json1 fts5" cmd/httpd/main.go \
		-server-uri http://localhost:8080 \
		-spelunker-uri 'sql://sqlite3?vfs=$(VFS_URI)'

lambda:
	@make lambda-server

# Note: This requires 'zig'

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


