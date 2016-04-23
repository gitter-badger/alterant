VERSION := $(shell cat ./VERSION)
OBJS := install/lib/libcrypto.a install/lib/libssl.a install/lib/libssh2.a install/lib/libgit2.a
DARWIN_OUTPUTS := ./build/darwin-amd64/alterant
LINUX_OUTPUTS := ./build/linux-amd64/alterant

.PHONY: vendor

vendor:
	gvt restore
	git clone https://github.com/libgit2/git2go.git ./vendor/github.com/libgit2/git2go
	cd ./vendor/github.com/libgit2/git2go; git checkout next; git submodule update --init
	git submodule update --init

# -L is being interpreted relative to the location of the package being built
#  so we execute the script from within git2go
#  NOTE: the order of dependent targets is important, ensuring that system libraries are not linked
deps: vendor $(OBJS)
	cd vendor/github.com/libgit2/git2go; ../../../../script/with-static-all.sh go install ./...

install/lib/libssl.a:
	@echo [Building libssl]
	./script/build-openssl-static.sh

install/lib/libcrypto.a: install/lib/libssl.a

install/lib/libssh2.a:
	@echo [Building libssh2]
	./script/build-libssh2-static.sh

install/lib/libgit2.a:
	@echo [Building libgit2]
	chmod +x vendor/github.com/libgit2/git2go/script/*.sh
	patch vendor/github.com/libgit2/git2go/script/build-libgit2-static.sh < script/patch_build_libbgit2_static.patch
	./vendor/github.com/libgit2/git2go/script/build-libgit2-static.sh

build/darwin-amd64/alterant:
	env CGO_LDFLAGS="-L$(PWD)/install/lib -lgit2 -lssh2 -lssl -lcrypto -lcurl -liconv -ldl -lz -framework CoreFoundation -framework Security" PKG_CONFIG_PATH="$(PWD)/install/lib/pkgconfig:$(PWD)/install/lib64/pkgconfig" GOOS=darwin GOARCH=amd64 go build -x -ldflags "-X main.version=$(VERSION)" -o build/darwin-amd64/alterant

build/linux-amd64/alterant:
	@echo [Building alterant for linux amd64]
	env CGO_LDFLAGS="-L$(PWD)/install/lib -lgit2 -lssh2 -lrt -lssl -lcrypto -ldl -lz" PKG_CONFIG_PATH="$(PWD)/install/lib/pkgconfig:$(PWD)/install/lib64/pkgconfig" GOOS=linux GOARCH=amd64 go build -x -ldflags '-X main.version=$(VERSION)'  -o build/linux-amd64/alterant

darwin: build/darwin-amd64/alterant

linux: build/linux-amd64/alterant

release: clean $(TARGET)
	@echo [Building alterant release for $(TARGET) amd64]
ifeq ($(TARGET),linux)
	cd ./build/linux-amd64; tar -cJf alterant-linux-amd64-$(VERSION).tar.xz alterant
endif
ifeq ($(TARGET),darwin)
	cd ./build/darwin-amd64; tar -cJf alterant-darwin-amd64-$(VERSION).tar.xz alterant
endif

deps-clean:
	rm -rf install
	rm -rf submodules/*
	find ./vendor/* -type d -maxdepth 0 -exec rm -rfv {} \;

clean:
	rm -rf ./build
