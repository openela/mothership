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
    name = "memory",
    srcs = ["memory.go"],
    importpath = "github.com/openela/mothership/base/storage/memory",
    visibility = ["//visibility:public"],
    deps = [
        "//base/go/storage",
        "//vendor/github.com/go-git/go-billy/v5:go-billy",
        "//vendor/github.com/pkg/errors",
    ],
)

go_test(
    name = "memory_test",
    size = "small",
    srcs = ["memory_test.go"],
    embed = [":memory"],
    deps = [
        "//base/go/storage",
        "//vendor/github.com/go-git/go-billy/v5/memfs",
        "//vendor/github.com/stretchr/testify/require",
    ],
)
