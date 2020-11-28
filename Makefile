.DEFAULT_GOAL = help

vendor  := vendor
target  := target
bin     := $(target)/bin
openapi := $(target)/openapi
reports := $(target)/reports

## build: Compile binaries.
go_src := $(shell find * -name *.go -not -path "$(vendor)/*" -not -path "$(target)/*")
go_out := $(patsubst cmd/%/main.go,$(bin)/%,$(wildcard cmd/*/main.go))

.PHONY: build
build: $(go_out)

$(bin)/%: cmd/%/main.go $(go_src) | $(bin)
	go build --trimpath --ldflags='-X "main.version=$(version)"' -o=$@ $<

$(bin):
	mkdir -p $@

## reports: Generate reports.
.PHONY: reports
reports: $(reports)/gosec.xml

$(reports)/gosec.html: $(go_src) | $(reports)
	gosec -fmt=html -out=$@ ./...

$(reports)/gosec.xml: $(go_src) | $(reports)
	gosec -fmt=junit-xml -out=$@ ./...

$(reports)/cover.out $(reports)/tests.xml: go/test

$(reports):
	mkdir -p $@

## coverage: Display test's coverage.
.PHONY: cover coverage
cover coverage: $(reports)/cover.out
	coverage -v -i "$(shell find cmd -name *.go)" target/reports/cover.out

## tests: Run tests and generate coverage report.
.PHONY: test tests
test tests: go/race go/test

.PHONY: go/race
go/race: $(go_src)
	go test -short -race -count=100 ./...

.PHONY: go/test
go/test: $(go_src) | $(reports)
	go test -v -coverpkg ./... -covermode=atomic -coverprofile=$(reports)/cover.out ./... | \
	go-junit-report -set-exit-code > $(reports)/tests.xml

## lint: Run static analysis checks.
.PHONY: lint
lint: go/lint proto/lint proto/break

.PHONY: go/lint
go/lint:
	golangci-lint run

## vendor: Make vendored copy of dependencies.
$(vendor): go.mod go.sum
	go mod vendor -v

## clean: Remove created resources.
.PHONY: clean
clean:
	rm -rf $(vendor) $(target)

## version: Display current version.
git_tag := $(shell git describe --tags 2>/dev/null)
git_sha := $(shell git rev-parse HEAD 2>/dev/null)
version := $(if $(git_tag),$(git_tag),(unknown on $(git_sha)))

.PHONY: version
version:
	@echo "version $(version)"

## help: Display available targets.
.PHONY: help
help: Makefile
	@echo "Usage: make [target]"
	@echo
	@echo "Targets:"
	@sed -n 's/^## //p' $< | awk -F ':' '{printf "  %-16s%s\n",$$1,$$2}'
