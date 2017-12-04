> **Note:** ASG Lifecycle Manager is a community driven project. LogicMonitor support will not assist in any issues related to ASG Lifecycle Manager.

## ASG Lifecycle Manager is a tool for managing AWS auto scaling group events
in your Kubernetes cluster.
-   **Automated Node Draining:** Leverages AWS auto scaling group lifecycle
hooks to automatically drain pods from nodes scheduled for termination.

## ASG Lifecycle Manager Overview
The ASG Lifecycle Manager leverages [Auto Scaling Lifecycle Hooks](https://docs.aws.amazon.com/autoscaling/latest/userguide/lifecycle-hooks.html#sqs-notifications) and [SQS](https://aws.amazon.com/documentation/sqs/) to
automatically drain pods from nodes that are scheduled for termination by an
auto scaling group.

Lifecycle events are sent to an SQS queue which is consumed by the Lifecycle
Manager. When the Lifecycle Manager receives a termination event, it will
identify the targeted node(s) and instruct Kubernetes to drain all pods from
the affected node(s). Once the pods are drained, the Lifecycle Manager will
notify that auto scaling group that it may now proceed with the termination.

## ASG Lifecycle Manager Configuration Requirements
- For this to have any value, you should be using an ** [EC2 Auto Scaling Group](https://docs.aws.amazon.com/autoscaling/latest/userguide/AutoScalingGroup.html) **
to manage at least some of your worker nodes
- The auto scaling group must be configured to send [Auto Scaling Lifecycle Hooks](https://docs.aws.amazon.com/autoscaling/latest/userguide/lifecycle-hooks.html#sqs-notifications)
to an [SQS](https://aws.amazon.com/documentation/sqs/).
- The Lifecycle Manager must run in the cluster using a service account with
appropriate RBAC permissions for deleting pods.
**Note:** this documentation should be improved to accurately enumerate specific
permissions.
- It is *not* required but currently recommended to run the Lifecycle Manager
on the master node to protect against the scenario in which the worker node
running the Lifecycle Manager gets scheduled for termination

## ASG Lifecycle Manager Configuration Options
| Name                                   | Type   | Required | Default | Description                                                        |
|----------------------------------------|--------|----------|---------|--------------------------------------------------------------------|
| NODEMAN_AWS_REGION                     | string | yes      |         | AWS region containing the source SQS queue                         |
| NODEMAN_AWS_SQS_QUEUE_URL              | string | yes      |         | URL of the source SQS queue                                        |
| NODEMAN_CONSUMER_THREADS               | int    | no       | 5       | Number of Lifecycle Manager consumer threads                       |
| NODEMAN_DEBUG                          | bool   | no       | false   | Enable debug logging                                               |
| NODEMAN_DEFAULT_VISIBILITY_TIMEOUT_SEC | int    | no       | 300     | SQS message visibility timeout for processing hooks                |
| NODEMAN_ERROR_VISIBILITY_TIMEOUT_SEC   | int    | no       | 60      | SQS message visibility timeout for messages with processing errors |
| NODEMAN_QUEUE_WAIT_TIME_SEC            | int    | no       | 5       | Delay between attempted SQS message retrievals                     |

**Note:** The Lifecycle Manager will let the AWS SDK transparently decide the
most appropriate way to authenticate with AWS. For example, you may choose to
leverage EC2 roles for the node running the Lifecycle Manager, mount a shared
credentials file into the Lifecycle Manager container, or create the container
with AWS access tokens set as environment variables. See the section on
(Configuring Credentials)[https://github.com/aws/aws-sdk-go] for more information.

## ASG Lifecycle Manager Setup Examples
The examples folder contains an example Helm chart and Terraform module.
These examples are meant to illustrate the various dependencies and requirements
of the Lifecycle Manager. This should not be considered production-ready code.
Copy/paste at your own risk.

- [Helm Chart](examples/chart)
- [Terraform module](examples/terraform)

### License
[![license](https://img.shields.io/github/license/logicmonitor/k8s-argus.svg?style=flat-square)](https://github.com/logicmonitor/k8s-argus/blob/master/LICENSE)


See the [documentation](https://docs.aws.amazon.com/autoscaling/latest/userguide/lifecycle-hooks.html) to discover more about AWS auto scaling lifecycle hooks.
