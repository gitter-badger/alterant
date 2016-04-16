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

all: linux-amd64

darwin-amd64:
	@echo "[darwin-amd64]"
	@env GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o build/darwin-amd64/alterant

darwin-386:
	@echo "[darwin-386]"
	@env GOOS=darwin GOARCH=386 go build -ldflags "-X main.version=$(VERSION)" -o build/darwin-386/alterant

linux-amd64:
	@echo "[linux-amd64]"
	@env CGO_LDFLAGS="-L/vagrant/go/src/github.com/autonomy/alterant/install/lib -lgit2 -lssh2 -lrt -lssl -lcrypto -ldl -lz" PKG_CONFIG_PATH="$(PWD)/install/lib/pkgconfig:$(PWD)/install/lib64/pkgconfig" GOOS=linux GOARCH=amd64 go build -x -ldflags '-X main.version=$(VERSION)'  -o build/linux-amd64/alterant

linux-386:
	@echo "[linux-386]"
	@env GOOS=linux GOARCH=386 go build -ldflags "-X main.version=$(VERSION)" -o build/linux-386/alterant

build-libgit2:
	cd ./submodules/git2go; git checkout next; git submodule update --init
	patch submodules/git2go/script/build-libgit2-static.sh < patch_build_libbgit2_static.patch
	./submodules/git2go/script/build-libgit2-static.sh

build-openssl:
	./script/build-openssl-static.sh

build-libssh2:
	./script/build-libssh2-static.sh

# -L is being interpreted relative to the location of the package being built
#  so we execute the script from within git2go
deps: build-openssl build-libssh2 build-libgit2
	gvt restore
	cd submodules/git2go; ../../script/with-static-all.sh go install ./...

release: clean all
	@echo "[Release $(VERSION)]"
	@cd ./build/darwin-386; tar -cJf alterant-darwin-i386-$(VERSION).tar.xz alterant
	@cd ./build/darwin-amd64; tar -cJf alterant-darwin-amd64-$(VERSION).tar.xz alterant
	@cd ./build/linux-386; tar -cJf alterant-linux-i386-$(VERSION).tar.xz alterant
	@cd ./build/linux-amd64; tar -cJf alterant-linux-amd64-$(VERSION).tar.xz alterant

deps-clean:
	@rm -rf install
	@rm -rf submodules/*

clean:
	@ rm -rf ./build
