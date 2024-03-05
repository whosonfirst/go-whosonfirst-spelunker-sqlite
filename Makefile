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

# This does not work when run under MacOS...
#
# https://github.com/mattn/go-sqlite3
# https://github.com/psanford/sqlite3vfs
#
# brew install filosottile/musl-cross/musl-cross --build-from-source --without-x86_64 --without-aarch64 --with-arm-hf --with-arm-linux-gnueabihf
# brew install arm-linux-gnueabihf-binutils
#
# fails with:
# cgo: C compiler "arm-linux-gnueabihf-gcc" not found: exec: "arm-linux-gnueabihf-gcc": executable file not found in $PATH
# There are no gcc or g++ tools in /usr/local/Cellar/arm-linux-gnueabihf-binutils/2.41_1/bin/
#
# If I try to use the (AArch32 bare-metal target (arm-none-eabi)) tools from:
# https://developer.arm.com/downloads/-/arm-gnu-toolchain-downloads
#
# At least it's a different error...
#
# arm-none-eabi-gcc: error: unrecognized command-line option '-pthread'; did you mean '-fpthread'?
#
# Same same if I run from  brew install --cask gcc-arm-embedded
# ARMGXX=/Applications/ArmGNUToolchain/13.2.Rel1/arm-none-eabi/bin/arm-none-eabi-g++
# ARMGCC=/Applications/ArmGNUToolchain/13.2.Rel1/arm-none-eabi/bin/arm-none-eabi-gcc
#
# Maybe this doesn't work with "arm-none-eabi" ?

ARMBIN=$(CWD)/work/arm-gnu-toolchain-13.2.Rel1-darwin-arm64-arm-none-eabi/bin
ARMGCC=$(ARMBIN)/arm-none-eabi-gcc
ARMGXX=$(ARMBIN)/arm-none-eabi-g++

lambda-server:
	if test -f bootstrap; then rm -f bootstrap; fi
	if test -f server.zip; then rm -f server.zip; fi
	# GOARCH=arm64 GOOS=linux go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -tags "lambda.norpc" -o bootstrap cmd/httpd/main.go
	CC=$(ARMGCC) \
		CXX=$(ARMGXX) \
		CGO_ENABLED=1 \
		GOARCH=arm64 \
		GOOS=linux \
		go build -mod $(GOMOD) \
		-ldflags="$(LDFLAGS) -linkmode external -extldflags -static" \
		-tags "lambda.norpc icu json1 fts5" \
		-o bootstrap \
		cmd/httpd/main.go
	zip server.zip bootstrap
	rm -f bootstrap	
