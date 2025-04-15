package temporal

import (
	"context"
	"go.temporal.io/sdk/workflow"
)

type Temporal interface {
	Init() error
	AddWorker(queueName string, workflowFn any, activitiesFn ...any) error
	ExecuteWorker(ctx context.Context, queueName string, workflowFn any, data any) error
	ExecuteTask(ctx workflow.Context, activityFn any, data any, result any) error
}
