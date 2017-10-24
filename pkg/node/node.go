package node

import (
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/logicmonitor/k8s-asg-lifecycle-manager/pkg/kubectl"
	log "github.com/sirupsen/logrus"
)

// Node an EC2 node that needs to be drained
type Node struct {
	EC2           *ec2.EC2
	EC2InstanceID string
	hostname      string
}

// Drain all pods from the node using its aws private hostname
func (n Node) Drain() error {
	hostname, err := n.PrivateHostname()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	log.Infof("Draining node %s", hostname)
	k := &kubectl.Kubectl{}
	err = k.Exec([]string{"drain", hostname,
		"--delete-local-data",
		"--ignore-daemonsets"})
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

// Delete the node from the cluster
func (n Node) Delete() error {
	hostname, err := n.PrivateHostname()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	log.Infof("Deleting node %s", hostname)
	k := &kubectl.Kubectl{}
	err = k.Exec([]string{"delete", "node", hostname,
		"--force",
	})
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

// State of the EC2 instance
func (n Node) State() (string, error) {
	i, err := n.instance()
	if err != nil {
		log.Error(err.Error())
		return "", err
	}

	state := *i.State.Name
	return state, nil
}

func (n Node) instance() (*ec2.Instance, error) {
	log.Infof("Retrieving EC2 information for Instance ID %s", n.EC2InstanceID)
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("instance-id"),
				Values: []*string{
					aws.String(n.EC2InstanceID),
				},
			},
		},
	}
	res, err := n.EC2.DescribeInstances(params)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	switch {
	case len(res.Reservations[0].Instances) < 1:
		log.Errorf("No instances found with ID %s", n.EC2InstanceID)
		return nil, fmt.Errorf("No instances found with ID %s", n.EC2InstanceID)
	case len(res.Reservations[0].Instances) > 1:
		log.Errorf("Too many instances found with ID %s", n.EC2InstanceID)
		return nil, fmt.Errorf("Too many instances found with ID %s", n.EC2InstanceID)
	default:
		log.Infof("Found instance with ID %s", n.EC2InstanceID)
		return res.Reservations[0].Instances[0], nil
	}
}

// PrivateHostname the EC2 instance's private hostname
func (n Node) PrivateHostname() (string, error) {
	if n.hostname != "" {
		return n.hostname, nil
	}

	i, err := n.instance()
	if err != nil {
		log.Error(err.Error())
		return "", err
	}

	n.hostname = n.formatHostname(*i.PrivateDnsName)
	log.Infof("Found instance private hostname %s", n.hostname)
	return n.hostname, nil
}

func (n Node) formatHostname(hostname string) string {
	// turn ip-10-35-120-96.us-west-1.ec2.internal
	// into ip-10-35-120-96.us-west-1.compute.internal
	var re = regexp.MustCompile(`(ec2)`)
	return re.ReplaceAllString(hostname, `compute`)
}
