package nodeman

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/logicmonitor/k8s-asg-lifecycle-manager/pkg/config"
	"github.com/logicmonitor/k8s-asg-lifecycle-manager/pkg/consumer"
	"github.com/logicmonitor/k8s-asg-lifecycle-manager/pkg/queue"
	"github.com/logicmonitor/k8s-asg-lifecycle-manager/pkg/types"
	log "github.com/sirupsen/logrus"
)

// NodeMan struct
type NodeMan struct {
	*types.Base
}

// NewBase instantiates and returns the base structure used throughout nodeman
func NewBase(config *config.Config) (*types.Base, error) {
	// AWS Session
	awsSess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(config.AwsRegion),
	}))

	base := &types.Base{
		AwsSess: awsSess,
		Config:  config,
	}

	return base, nil
}

// NewNodeMan instantiates and returns nodeman.
func NewNodeMan(base *types.Base) (*NodeMan, error) {
	nodeman := &NodeMan{
		Base: base,
	}

	return nodeman, nil
}

// Watch watches the SQS Queue for messages to remove nodes.
func (nm *NodeMan) Watch() {
	initialCtx, cancelCtx := context.WithCancel(context.Background())
	log.Info("Initializing global context")

	handleInterrupt(cancelCtx)

	asgSvc := autoscaling.New(nm.AwsSess)
	ec2Svc := ec2.New(nm.AwsSess)
	sqsSvc := sqs.New(nm.AwsSess)

	queue := &queue.Queue{
		SQS:  *sqsSvc,
		Base: nm.Base,
	}

	// start consumer threads
	for i := 0; i <= nm.Config.ConsumerThreads; i++ {
		consumer := consumer.Consumer{
			ASG:   asgSvc,
			Base:  nm.Base,
			EC2:   ec2Svc,
			Queue: queue,
		}
		log.Infof("Starting thread %v", i)
		err := consumer.Start(initialCtx)
		if err != nil {
			log.Error(err.Error())
			cancelCtx()
		}
	}

	<-initialCtx.Done()
}

func handleInterrupt(cancelCtx context.CancelFunc) {
	// cleanly handle sig int
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		for sig := range c {
			// got sig int
			log.Errorf("Caught sig %v..  exiting...", sig)
			log.Info("Cancelling context...")
			cancelCtx()
			log.Info("Waiting a sec to facilitate cleanup..")
			time.Sleep(time.Second)
			log.Info("Exiting")
			os.Exit(1)
		}
	}()
}
