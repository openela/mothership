# Copyright 2024 The Mothership Authors
# SPDX-License-Identifier: Apache-2.0

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "gitlab",
    srcs = ["gitlab.go"],
    importpath = "github.com/openela/mothership/base/forge/gitlab",
    visibility = ["//visibility:public"],
    deps = [
        "//base/go/forge",
        "//vendor/github.com/go-git/go-git/v5/plumbing/transport/http",
    ],
)
