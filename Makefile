# Generate build files run this after adding dependencies or adding build targets like new binary or tests
.PHONY: generate
generate:
	bazel run //:gazelle

.PHONY: build
build:
	bazel build //:project

