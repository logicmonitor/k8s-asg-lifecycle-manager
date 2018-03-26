package node

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/logicmonitor/k8s-asg-lifecycle-manager/pkg/kubectl"
	log "github.com/sirupsen/logrus"
)

// Node an EC2 node that needs to be drained
type Node struct {
	EC2           *ec2.EC2
	EC2InstanceID string
	name          string
	instance      *ec2.Instance
}

// New instantiates and returns a Node.
func New(ec2 *ec2.EC2, id string, short bool) (*Node, error) {
	i, err := instance(ec2, id)
	if err != nil {
		return nil, err
	}
	name, err := name(i, short)
	if err != nil {
		return nil, err
	}

	n := &Node{
		EC2:           ec2,
		EC2InstanceID: id,
		name:          name,
		instance:      i,
	}

	return n, nil
}

// Drain all pods from the node using its aws private hostname
func (n *Node) Drain() error {
	log.Infof("Draining node %s", n.name)
	k := &kubectl.Kubectl{}
	args := []string{
		"drain",
		n.name,
		"--delete-local-data",
		"--ignore-daemonsets",
	}
	err := k.Exec(args)
	if err != nil {
		return fmt.Errorf("drain node: %s", err.Error())
	}
	return nil
}

// Delete the node from the cluster
func (n *Node) Delete() error {
	log.Infof("Deleting node %s", n.name)
	k := &kubectl.Kubectl{}
	args := []string{
		"delete",
		"node",
		n.name,
		"--force",
	}
	err := k.Exec(args)
	if err != nil {
		return fmt.Errorf("delete node: %s", err.Error())
	}
	return nil
}

// State of the EC2 instance
func (n *Node) State() (string, error) {
	state := *n.instance.State.Name
	return state, nil
}

func instance(ec2client *ec2.EC2, id string) (*ec2.Instance, error) {
	log.Infof("Retrieving EC2 information for Instance ID %s", id)
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("instance-id"),
				Values: []*string{
					aws.String(id),
				},
			},
		},
	}
	res, err := ec2client.DescribeInstances(params)
	if err != nil {
		return nil, err
	}

	switch {
	case len(res.Reservations[0].Instances) < 1:
		log.Errorf("No instances found with ID %s", id)
		return nil, fmt.Errorf("No instances found with ID %s", id)
	case len(res.Reservations[0].Instances) > 1:
		log.Errorf("Too many instances found with ID %s", id)
		return nil, fmt.Errorf("Too many instances found with ID %s", id)
	default:
		log.Infof("Found instance with ID %s", id)
		return res.Reservations[0].Instances[0], nil
	}
}

func name(i *ec2.Instance, short bool) (string, error) {
	var name string
	if i.PrivateDnsName == nil {
		return "", fmt.Errorf("cannot determine node name, EC2 private DNS name is nil")
	}

	if short {
		parts := strings.Split(*i.PrivateDnsName, ".")
		name = parts[0]
	} else {
		name = *i.PrivateDnsName
	}
	log.Infof("Using node name %s", name)

	return name, nil
}
