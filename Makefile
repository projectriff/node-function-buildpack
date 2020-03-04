.PHONY: clean build test acceptance all
GO_SOURCES = $(shell find . -type f -name '*.go')

PACK=go run github.com/buildpacks/pack/cmd/pack

all: test build acceptance

build: artifactory/io/projectriff/node/io.projectriff.node

test:
	go test -v ./...

acceptance: acceptance/testdata/builder.toml
	$(PACK) create-builder -b acceptance/testdata/builder.toml projectriff/builder
	docker pull cloudfoundry/build:base-cnb
	docker pull cloudfoundry/run:base-cnb
	GO111MODULE=on go test -v -tags=acceptance ./acceptance

artifactory/io/projectriff/node/io.projectriff.node: buildpack.toml $(GO_SOURCES)
	rm -fR $@ 							&& \
	./ci/package.sh						&& \
	mkdir $@/latest 					&& \
	tar -C $@/latest -xzf $@/*/*.tgz

acceptance/testdata/builder.toml: acceptance/testdata/builder.toml.tpl go.mod
	./ci/apply-template.sh acceptance/testdata/builder.toml.tpl > acceptance/testdata/builder.toml

clean:
	rm -fR artifactory/
	rm -fR dependency-cache/
	rm acceptance/testdata/builder.toml
