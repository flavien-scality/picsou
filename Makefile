BINARIES ?=	picsou
GODIR ?= $(shell echo $$GOPATH)/src/github.com/scality/picsou

GOFLAGS ?= GOOS=linux GOARCH=amd64 CGO_ENABLED=0
PKG_BASE_DIR ?=	./pkg
CONVEY_PORT ?=	9042
SOURCES :=	$(shell find . -type f -name "*.go")
COMMANDS :=	$(shell go list ./cmd/...)
PACKAGES :=	$(shell go list ./pkg/...)

GO ?=		$(GOENV) go


all:	build


.PHONY: build
build:	$(BINARIES)


$(BINARIES):	$(SOURCES)
	$(GOFLAGS) $(GO) get ./...
	$(GOFLAGS) $(GO) build -v ./...
	$(GOFLAGS) $(GO) build -o $@ ./cmd/$@


.PHONY: lint
lint:
	golint ./...


.PHONY: test
test: deps
	$(GO) test -v $(PACKAGES)


.PHONY: deps
deps:
	pip install -r requirements.txt
	$(GO) get -t ./...


.PHONY: install
install:
	$(GO) install $(COMMANDS) 


.PHONY: clean
clean:
	rm -f $(BINARIES)
	rm -f coverage.*


.PHONY: re
re:	clean all


.PHONY: convey
convey:
	goconvey -cover -port=$(CONVEY_PORT) -workDir="$(realpath $(PKG_BASE_DIR))" -depth=1


.PHONY:	cover
cover:	profile.txt
	$(GO) tool cover -html=coverage.txt -o coverage.html

.PHONY: docker-build
docker-build:
	docker run --rm -v "$(shell echo $$HOME)/.ssh:/root/.ssh" -v "$(shell echo $$PWD)":/go/src/github.com/scality/picsou -w /go/src/github.com/scality/picsou golang:1.8.0-onbuild bash -c make

.PHONY: docker-build
docker-build:
	docker run --rm -v "$(shell echo $$HOME)/.ssh:/root/.ssh" -v "$(shell echo $$PWD)":/go/src/github.com/scality/picsou -w /go/src/github.com/scality/picsou golang:1.8.0-onbuild bash -c make


profile.txt:	$(SOURCES)
	echo -n "" > coverage.txt
	for d in $(PACKAGES); do \
		go test -v -coverprofile=profile.out $$d || exit 1 ; \
		if [ -f profile.out ]; then \
			cat profile.out >> coverage.txt; \
			rm profile.out \
		; fi \
	done
