package types

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/logicmonitor/k8s-asg-lifecycle-manager/pkg/config"
)

// Base is a struct for embedding
type Base struct {
	AwsSess *session.Session
	Config  *config.Config
}
