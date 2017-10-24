package constants

var (
	// Version is the NodeMan version and is set at build time.
	Version string
)

const (
	// AsgActionContinue ASG continue action
	AsgActionContinue = "CONTINUE"
	// InstanceTerminating is the string indicating an AWS ASG terminate action
	InstanceTerminating = "autoscaling:EC2_INSTANCE_TERMINATING"
	// InstanceTerminatedState is the string indicating an AWS instance has terminated
	InstanceTerminatedState = "terminated"
	// UserAgentBase is the base string for the User-Agent HTTP header.
	UserAgentBase = "LogicMonitor NodeMan/"
)
