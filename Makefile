.PHONY: all
all: kc

.PHONY: kc
kc:
	bazel build //cmd/kc

.PHONY: gofmt
gofmt:
	gofmt -w -s cmd/ pkg/

.PHONY: gazelle
gazelle:
	bazel run //:gazelle

.PHONY: dep
dep:
	dep ensure
	find vendor -name "BUILD" -delete
	find vendor -name "BUILD.bazel" -delete
	bazel run //:gazelle -- -proto disable

.PHONY: goimports
goimports:
	goimports -w cmd/ pkg/
