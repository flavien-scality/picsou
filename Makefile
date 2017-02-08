BINARIES ?=	picsou
GODIR ?= $GOPATH/src/github.com/scality/picsou

PKG_BASE_DIR ?=	./pkg
CONVEY_PORT ?=	9042
SOURCES :=	$(shell find . -type f -name "*.go")
COMMANDS :=	$(shell go list ./cmd/...)
PACKAGES :=	$(shell go list ./pkg/...)

GO ?=		$(GOENV) go


all:	deps build


.PHONY: build
build:	$(BINARIES)


$(BINARIES):	$(SOURCES)
	$(GO) build -v ./...
	$(GO) build -o $@ ./cmd/$@


.PHONY: lint
lint:
	golint ./...


.PHONY: test
test: deps
	$(GO) test -v $(PACKAGES)


.PHONY: deps
deps:
	$(GO) get -d -t -v ./...


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


profile.txt:	$(SOURCES)
	echo -n "" > coverage.txt
	for d in $(PACKAGES); do \
		go test -v -coverprofile=profile.out $$d || exit 1 ; \
		if [ -f profile.out ]; then \
			cat profile.out >> coverage.txt; \
			rm profile.out \
		; fi \
	done
