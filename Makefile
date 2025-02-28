CWD=$(shell pwd)
GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

# The HTTP URI to the database being used. Adjust to taste.
VFS_URI=http://localhost:8081/ca-search.db

server:
	go run 	-mod $(GOMOD) \
		-tags "icu json1 fts5" \
		cmd/httpd/main.go \
		-server-uri http://localhost:8080 \
		-spelunker-uri 'sql://sqlite3?dsn=$(DSN)'

lambda:
	@make lambda-server

# Things to note:
# 1) This requires 'zig' to build
# 2) This either needs to be run with SPELUNKER_SERVER_URI set to lambda:// or using
#    the https://github.com/awslabs/aws-lambda-web-adapter "layer"

lambda-server:
	if test -f bootstrap; then rm -f bootstrap; fi
	if test -f server.zip; then rm -f server.zip; fi
	GOOS=linux \
		GOARCH=arm64 \
		CGO_ENABLED=1 \
		CGO_CFLAGS=-D_LARGEFILE64_SOURCE \
		CC='zig cc -target aarch64-linux-musl' \
		CXX='zig cc -target aarch64-linux-musl' \
		go build -mod $(GOMOD) \
		-ldflags="$(LDFLAGS) -linkmode external -extldflags -static" \
		-tags "netgo,sqlite_omit_load_extension,fts5,lambda.norpc" \
		-o bootstrap \
		cmd/httpd/main.go
	zip server.zip bootstrap
	rm -f bootstrap	


