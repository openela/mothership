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

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "s3",
    srcs = [
        "detector.go",
        "s3.go",
    ],
    importpath = "github.com/openela/mothership/base/storage/s3",
    visibility = ["//visibility:public"],
    deps = [
        "//base/go",
        "//base/go/awsutils",
        "//base/go/storage",
        "//vendor/github.com/aws/aws-sdk-go/aws",
        "//vendor/github.com/aws/aws-sdk-go/aws/awserr",
        "//vendor/github.com/aws/aws-sdk-go/aws/session",
        "//vendor/github.com/aws/aws-sdk-go/service/s3",
        "//vendor/github.com/aws/aws-sdk-go/service/s3/s3manager",
        "//vendor/github.com/pkg/errors",
        "//vendor/github.com/urfave/cli/v2:cli",
    ],
)
