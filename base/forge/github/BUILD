# Copyright 2023 Peridot Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "github",
    srcs = ["github.go"],
    importpath = "github.com/openela/mothership/base/forge/github",
    visibility = ["//visibility:public"],
    deps = [
        "//base/go/forge",
        "//vendor/github.com/go-git/go-git/v5/plumbing/transport/http",
        "//vendor/github.com/golang-jwt/jwt/v5:jwt",
    ],
)

go_test(
    name = "github_test",
    size = "small",
    srcs = ["github_test.go"],
    embed = [":github"],
    deps = [
        "//vendor/github.com/go-git/go-git/v5/plumbing/transport/http",
        "//vendor/github.com/jarcoal/httpmock",
        "//vendor/github.com/stretchr/testify/require",
    ],
)
