# Environments can be viewed at https://golang.org/doc/install/source#environment

# $GOOS	$GOARCH
# darwin	386
# darwin	amd64
# darwin	arm
# darwin	arm64
# dragonfly	amd64
# freebsd	386
# freebsd	amd64
# freebsd	arm
# linux	386
# linux	amd64
# linux	arm
# linux	arm64
# linux	ppc64
# linux	ppc64le
# netbsd	386
# netbsd	amd64
# netbsd	arm
# openbsd	386
# openbsd	amd64
# openbsd	arm
# plan9	386
# plan9	amd64
# solaris	amd64
# windows	386
# windows	amd64

VERSION := $(shell cat ./VERSION)

all: darwin-amd64 darwin-386 linux-amd64 linux-386

darwin-amd64:
	@echo "[darwin-amd64]"
	@env GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o build/darwin-amd64/alterant

darwin-386:
	@echo "[darwin-386]"
	@env GOOS=darwin GOARCH=386 go build -ldflags "-X main.version=$(VERSION)" -o build/darwin-386/alterant

linux-amd64:
	@echo "[linux-amd64]"
	@env GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o build/linux-amd64/alterant

linux-386:
	@echo "[linux-386]"
	@env GOOS=linux GOARCH=386 go build -ldflags "-X main.version=$(VERSION)" -o build/linux-386/alterant

release: clean all
	@echo "[Release $(VERSION)]"
	@cd ./build/darwin-386; tar -cJf alterant-darwin-i386-$(VERSION).tar.xz alterant
	@cd ./build/darwin-amd64; tar -cJf alterant-darwin-amd64-$(VERSION).tar.xz alterant
	@cd ./build/linux-386; tar -cJf alterant-linux-i386-$(VERSION).tar.xz alterant
	@cd ./build/linux-amd64; tar -cJf alterant-linux-amd64-$(VERSION).tar.xz alterant

clean:
	@ rm -rf ./build
