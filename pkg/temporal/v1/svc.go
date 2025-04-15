package temporalv1

import (
	"context"
	"github.com/ladmakhi81/learnup/pkg/env"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"time"
)

type TemporalSvc struct {
	client client.Client
	config *env.EnvConfig
}

func NewTemporalSvc(config *env.EnvConfig) *TemporalSvc {
	return &TemporalSvc{
		config: config,
	}
}

func (svc *TemporalSvc) Init() error {
	temporalClient, temporalClientErr := client.Dial(client.Options{
		HostPort: svc.config.Temporal.Endpoint,
	})
	if temporalClientErr != nil {
		return temporalClientErr
	}
	svc.client = temporalClient
	return nil
}

func (svc *TemporalSvc) AddWorker(queueName string, workflowFn any, activitiesFn ...any) error {
	w := worker.New(svc.client, queueName, worker.Options{})
	w.RegisterWorkflow(workflowFn)
	for _, activity := range activitiesFn {
		w.RegisterActivity(activity)
	}
	if err := w.Start(); err != nil {
		return err
	}
	return nil
}

func (svc *TemporalSvc) ExecuteWorker(ctx context.Context, queueName string, workflowFn any, data any) error {
	_, err := svc.client.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{TaskQueue: queueName},
		workflowFn,
		data,
	)
	if err != nil {
		return err
	}
	return nil
}

func (svc *TemporalSvc) ExecuteTask(ctx workflow.Context, activityFn any, data any, result any) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour * 24,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	return workflow.ExecuteActivity(ctx, activityFn, data).Get(ctx, result)
}
