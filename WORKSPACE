# Go rules
http_archive(
    name = "io_bazel_rules_go",
    url = "https://github.com/bazelbuild/rules_go/releases/download/0.10.0/rules_go-0.10.0.tar.gz",
    sha256 = "53c8222c6eab05dd49c40184c361493705d4234e60c42c4cd13ab4898da4c6be",
)
load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains", "go_repositories")
# Gazelle Go build file generator
http_archive(
    name = "bazel_gazelle",
    url = "https://github.com/bazelbuild/bazel-gazelle/releases/download/0.10.0/bazel-gazelle-0.10.0.tar.gz",
    sha256 = "6228d9618ab9536892aa69082c063207c91e777e51bd3c5544c9c060cafe1bd8",
)
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

# Docker rules
git_repository(
    name = "io_bazel_rules_docker",
    remote = "https://github.com/bazelbuild/rules_docker.git",
    tag = "v0.4.0",
)
load(
    "@io_bazel_rules_docker//go:image.bzl",
    _go_image_repos = "repositories",
)

go_rules_dependencies()
go_register_toolchains()
gazelle_dependencies()
_go_image_repos()


