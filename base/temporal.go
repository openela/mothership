// Copyright 2024 The Mothership Authors
// SPDX-License-Identifier: Apache-2.0

package base

import (
	"context"
	"crypto/tls"
	"github.com/pkg/errors"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log/slog"
	"strings"
	"time"
)

type temporalTQInterceptor struct {
	interceptor.ClientInterceptorBase
	interceptor.ClientOutboundInterceptorBase

	taskQueue string
}

func (tqi *temporalTQInterceptor) InterceptClient(next interceptor.ClientOutboundInterceptor) interceptor.ClientOutboundInterceptor {
	return &temporalTQInterceptor{
		ClientOutboundInterceptorBase: interceptor.ClientOutboundInterceptorBase{
			Next: next,
		},
		taskQueue: tqi.taskQueue,
	}
}

func (tqi *temporalTQInterceptor) ExecuteWorkflow(ctx context.Context, in *interceptor.ClientExecuteWorkflowInput) (client.WorkflowRun, error) {
	in.Options.TaskQueue = tqi.taskQueue
	return tqi.Next.ExecuteWorkflow(ctx, in)
}

func NewTemporalClient(host string, namespace string, taskQueue string, opts client.Options) (client.Client, error) {
	// If host contains :443, then use TLS
	if strings.Contains(host, ":443") {
		opts.ConnectionOptions = client.ConnectionOptions{
			DialOptions: []grpc.DialOption{
				grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
			},
		}
	}

	opts.HostPort = host

	slog.Info("Using Temporal namespace", "namespace", namespace)

	// Register namespace (ignore error if already exists)
	nscl, err := client.NewNamespaceClient(opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create namespace client")
	}

	// Set default retention period to 5 days
	dur := 5 * 24 * time.Hour
	err = nscl.Register(context.TODO(), &workflowservice.RegisterNamespaceRequest{
		Namespace:                        namespace,
		WorkflowExecutionRetentionPeriod: &dur,
	})
	if err != nil && !strings.Contains(err.Error(), "Namespace already exists") {
		return nil, errors.Wrap(err, "failed to register namespace")
	}

	// Set namespace in opts
	opts.Namespace = namespace

	// Set interceptor to set task queue
	if opts.Interceptors == nil {
		opts.Interceptors = []interceptor.ClientInterceptor{}
	}
	opts.Interceptors = append(opts.Interceptors, &temporalTQInterceptor{
		taskQueue: taskQueue,
	})

	slog.Info("Using Temporal task queue", "taskQueue", taskQueue)

	// Dial Temporal
	cl, err := client.Dial(opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to dial temporal")
	}

	return cl, nil
}
