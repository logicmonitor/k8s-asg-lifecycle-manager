# the kubernetes worker asg
resource "aws_autoscaling_group" "worker_asg" {
  name                      = "workers"
  max_size                  = "${var.worker_max_size}"
  min_size                  = "${var.worker_min_size}"
  health_check_grace_period = 300
  health_check_type         = "EC2"
  force_delete              = false
  launch_configuration      = "${aws_launch_configuration.worker.name}"
  vpc_zone_identifier       = ["${var.subnet_ids}"]
}

# launch configuration and instance profile for asg workers
resource "aws_launch_configuration" "worker" {
  name_prefix   = "worker-"
  image_id      = "${var.ami_id}"
  instance_type = "${var.worker_instance_type}"

  key_name             = "${var.key_pair_id}"
  iam_instance_profile = "${aws_iam_instance_profile.worker.id}"
  security_groups      = ["${aws_security_group.worker.id}"]
}

resource "aws_iam_role" "worker" {
    name               = "worker"
    path               = "/"
    assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": { "Service": "ec2.amazonaws.com"},
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_instance_profile" "worker" {
  name = "worker"
  role = "${aws_iam_role.worker.id}"
}

resource "aws_security_group" "worker" {
  name   = "worker"
  vpc_id = "${var.vpc_id}"
}

# asg lifecycle hook configuration for sending termination event hooks
resource "aws_autoscaling_lifecycle_hook" "termination_hook" {
  autoscaling_group_name  = "${aws_autoscaling_group.worker_asg.name}"
  default_result          = "CONTINUE"
  heartbeat_timeout       = 300
  lifecycle_transition    = "autoscaling:EC2_INSTANCE_TERMINATING"
  name                    = "termination-hook-worker"
  notification_target_arn = "${aws_sqs_queue.asg_lifecycle_queue.arn}"
  role_arn                = "${aws_iam_role.asg_lifecycle_role.arn}"
}

# sqs queue to receive termination hooks
resource "aws_sqs_queue" "asg_lifecycle_queue" {
  message_retention_seconds   = 600
  name                        = "asg-lifecycle-queue"
  receive_wait_time_seconds   = 10
  visibility_timeout_seconds  = 300
}

# policy, role, and attachments to allow asg to submit messages to sqs queue
data "aws_iam_policy_document" "asg_lifecycle_policy" {
  statement {
    actions = [
      "sqs:GetQueueUrl",
      "sqs:SendMessage",
    ]

    resources = [
      "${aws_sqs_queue.asg_lifecycle_queue.arn}",
    ]
  }
}


resource "aws_iam_role_policy" "asg_lifecycle_policy" {
  name    = "asg-lifecycle-role-policy"
  policy  = "${data.aws_iam_policy_document.asg_lifecycle_policy.json}"
  role    = "${aws_iam_role.asg_lifecycle_role.id}"
}

resource "aws_iam_role" "asg_lifecycle_role" {
  name = "asg-lifecycle-role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "autoscaling.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

# policy, role, and attachments to allow lifecycle manager to RUD hook messages
data "aws_iam_policy_document" "asg_lifecycle_consumer_policy" {
  statement {
    actions = [
      "sqs:ChangeMessageVisibility",
      "sqs:ChangeMessageVisibilityBatch",
      "sqs:DeleteMessage",
      "sqs:DeleteMessageBatch",
      "sqs:ReceiveMessage",
    ]

    resources = [
      "${aws_sqs_queue.asg_lifecycle_queue.arn}",
    ]
  }
}

resource "aws_iam_role_policy" "asg_lifecycle_policy_consumer" {
  name    = "asg-lifecycle-consumer-role-policy"
  policy  = "${data.aws_iam_policy_document.asg_lifecycle_consumer_policy.json}"
  role    = "${aws_iam_role.worker.id}"
}
