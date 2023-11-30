package mothershipadmin_rpc

import (
	"context"
	"github.com/openela/mothership/base"
	mshipadminpb "github.com/openela/mothership/proto/admin/v1"
	mothershippb "github.com/openela/mothership/proto/v1"
	"github.com/openela/mothership/third_party/googleapis/google/longrunning"
	v11 "go.temporal.io/api/enums/v1"
	"go.temporal.io/api/serviceerror"
	"go.temporal.io/api/workflowservice/v1"
	rpccode "google.golang.org/genproto/googleapis/rpc/code"
	rpcstatus "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) describeWorkflowToOperation(ctx context.Context, res *workflowservice.DescribeWorkflowExecutionResponse) (*longrunning.Operation, error) {
	if res.WorkflowExecutionInfo == nil {
		return nil, status.Error(codes.NotFound, "workflow not found")
	}
	if res.WorkflowExecutionInfo.Execution == nil {
		return nil, status.Error(codes.NotFound, "workflow not found")
	}

	op := &longrunning.Operation{
		Name: res.WorkflowExecutionInfo.Execution.WorkflowId,
	}

	// If the workflow is not running, we can mark the operation as done
	if res.WorkflowExecutionInfo.Status != v11.WORKFLOW_EXECUTION_STATUS_RUNNING {
		op.Done = true
	}

	// Add metadata
	rpmMetadata := &mshipadminpb.RetractEntryMetadata{
		StartTime: nil,
		EndTime:   nil,
	}
	st := res.WorkflowExecutionInfo.GetStartTime()
	if st != nil {
		rpmMetadata.StartTime = timestamppb.New(*st)
	}

	et := res.WorkflowExecutionInfo.GetCloseTime()
	if et != nil {
		rpmMetadata.EndTime = timestamppb.New(*et)
	}

	rpmMetadataAny, err := anypb.New(rpmMetadata)
	if err != nil {
		return op, nil
	}
	op.Metadata = rpmMetadataAny

	// If completed, add result
	// If failed, add error
	if res.WorkflowExecutionInfo.Status == v11.WORKFLOW_EXECUTION_STATUS_COMPLETED {
		// Complete, we need to get the result using GetWorkflow
		run := s.temporal.GetWorkflow(ctx, op.Name, "")

		var res mothershippb.ProcessRPMResponse
		if err := run.Get(ctx, &res); err != nil {
			return nil, err
		}

		resAny, err := anypb.New(&res)
		if err != nil {
			return nil, err
		}
		op.Result = &longrunning.Operation_Response{Response: resAny}
	} else if res.WorkflowExecutionInfo.Status == v11.WORKFLOW_EXECUTION_STATUS_FAILED {
		// Failed, we need to get the error using GetWorkflow
		run := s.temporal.GetWorkflow(ctx, op.Name, "")
		err := run.Get(ctx, nil)
		// No error so return with a generic error
		if err == nil {
			op.Result = &longrunning.Operation_Error{
				Error: &rpcstatus.Status{
					Code:    int32(rpccode.Code_INTERNAL),
					Message: "workflow failed",
				},
			}
			return op, nil
		}

		// Error, so return with the error
		op.Result = &longrunning.Operation_Error{
			Error: &rpcstatus.Status{
				Code:    int32(rpccode.Code_FAILED_PRECONDITION),
				Message: err.Error(),
			},
		}
	} else if res.WorkflowExecutionInfo.Status == v11.WORKFLOW_EXECUTION_STATUS_CANCELED {
		// Error, so return with the error
		op.Result = &longrunning.Operation_Error{
			Error: &rpcstatus.Status{
				Code:    int32(rpccode.Code_CANCELLED),
				Message: "workflow canceled",
			},
		}
	}

	return op, nil
}

func (s *Server) getOperation(ctx context.Context, name string) (*longrunning.Operation, error) {
	res, err := s.temporal.DescribeWorkflowExecution(ctx, name, "")
	if err != nil {
		if _, ok := err.(*serviceerror.NotFound); ok {
			return nil, status.Error(codes.NotFound, "workflow not found")
		}

		// Log error, but user doesn't need to know about it
		base.LogErrorf("failed to describe workflow: %v", err)
		return &longrunning.Operation{
			Name: name,
		}, nil
	}

	return s.describeWorkflowToOperation(ctx, res)
}

func (s *Server) GetOperation(ctx context.Context, req *longrunning.GetOperationRequest) (*longrunning.Operation, error) {
	// Get from Temporal. We don't care about long term storage, so we don't
	// need to store the operation in the database.
	return s.getOperation(ctx, req.Name)
}
