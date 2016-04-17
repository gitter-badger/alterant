# Environments can be viewed at https://golang.org/doc/install/source#environment

VERSION := $(shell cat ./VERSION)
OBJS := install/lib/libcrypto.a install/lib/libssl.a install/lib/libssh2.a install/lib/libgit2.a
DARWIN_OUTPUTS := ./build/darwin-386/alterant ./build/darwin-amd64/alterant
LINUX_OUTPUTS := ./build/linux-386/alterant ./build/linux-amd64/alterant

.PHONY: vendor

vendor:
	find ./vendor/* -type d -maxdepth 0 -exec rm -rfv {} \;
	gvt restore

install/lib/libssl.a:
	./script/build-openssl-static.sh

install/lib/libcrypto.a: install/lib/libssl.a

install/lib/libssh2.a:
	./script/build-libssh2-static.sh

install/lib/libgit2.a:
	cd ./submodules/git2go; git checkout next; git submodule update --init
	patch submodules/git2go/script/build-libgit2-static.sh < script/patch_build_libbgit2_static.patch
	./submodules/git2go/script/build-libgit2-static.sh

# -L is being interpreted relative to the location of the package being built
#  so we execute the script from within git2go
#  NOTE: the order of dependent targets is important, ensuring that system libraries are not linked
deps: vendor $(OBJS)
	cd submodules/git2go; ../../script/with-static-all.sh go install ./...

build/darwin-amd64/alterant:
	env CGO_LDFLAGS="-L$(PWD)/install/lib -lgit2 -lssh2 -lssl -lcrypto -lcurl -liconv -ldl -lz -framework CoreFoundation -framework Security" PKG_CONFIG_PATH="$(PWD)/install/lib/pkgconfig:$(PWD)/install/lib64/pkgconfig" GOOS=darwin GOARCH=amd64 go build -x -ldflags "-X main.version=$(VERSION)" -o build/darwin-amd64/alterant

build/darwin-386/alterant:
	env GOOS=darwin GOARCH=386 go build -ldflags "-X main.version=$(VERSION)" -o build/darwin-386/alterant

build/linux-amd64/alterant:
	env CGO_LDFLAGS="-L$(PWD)/install/lib -lgit2 -lssh2 -lrt -lssl -lcrypto -ldl -lz" PKG_CONFIG_PATH="$(PWD)/install/lib/pkgconfig:$(PWD)/install/lib64/pkgconfig" GOOS=linux GOARCH=amd64 go build -x -ldflags '-X main.version=$(VERSION)'  -o build/linux-amd64/alterant

build/linux-386/alterant:
	env GOOS=linux GOARCH=386 go build -ldflags "-X main.version=$(VERSION)" -o build/linux-386/alterant

darwin: build/darwin-amd64/alterant #build/darwin-386/alterant 

linux: build/linux-amd64/alterant #build/linux-386/alterant

release: clean $(TARGET)
ifeq ($(TARGET),linux)
	# cd ./build/linux-386; tar -cJf alterant-linux-i386-$(VERSION).tar.xz alterant
	cd ./build/linux-amd64; tar -cJf alterant-linux-amd64-$(VERSION).tar.xz alterant
endif
ifeq ($(TARGET),darwin)
	# cd ./build/darwin-386; tar -cJf alterant-darwin-i386-$(VERSION).tar.xz alterant
	cd ./build/darwin-amd64; tar -cJf alterant-darwin-amd64-$(VERSION).tar.xz alterant
endif

deps-clean:
	rm -rf install
	rm -rf submodules/*

clean:
	rm -rf ./build
