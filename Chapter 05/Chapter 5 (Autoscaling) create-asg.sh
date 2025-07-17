#!/bin/bash

# Fail fast
set -e

# Configurable Variables - Replace these with your actual values
TEMPLATE_NAME="my-web-app-template"
TEMPLATE_VERSION="1"
AMI_ID="ami-0abcdef1234567890"
INSTANCE_TYPE="t3.micro"
SECURITY_GROUP_ID="sg-0abcd1234efgh5678"
KEY_NAME="my-keypair"
ASG_NAME="my-asg"
MIN_SIZE=2
MAX_SIZE=10
DESIRED_CAPACITY=2
SUBNETS="subnet-abc123,subnet-def456"
POLICY_NAME="cpu-scale-out"
ALARM_NAME="HighCPU"
REGION="us-east-1"  # Change as per your region
ACCOUNT_ID="$(aws sts get-caller-identity --query Account --output text)"

# Create Launch Template
echo "Creating Launch Template..."
aws ec2 create-launch-template \
  --launch-template-name $TEMPLATE_NAME \
  --version-description "v1" \
  --launch-template-data "{
    \"ImageId\": \"$AMI_ID\",
    \"InstanceType\": \"$INSTANCE_TYPE\",
    \"SecurityGroupIds\": [\"$SECURITY_GROUP_ID\"],
    \"KeyName\": \"$KEY_NAME\"
  }" \
  --region $REGION || echo "Launch template may already exist."

# Create Auto Scaling Group
echo "Creating Auto Scaling Group..."
aws autoscaling create-auto-scaling-group \
  --auto-scaling-group-name $ASG_NAME \
  --launch-template LaunchTemplateName=$TEMPLATE_NAME,Version=$TEMPLATE_VERSION \
  --min-size $MIN_SIZE \
  --max-size $MAX_SIZE \
  --desired-capacity $DESIRED_CAPACITY \
  --vpc-zone-identifier "$SUBNETS" \
  --region $REGION

# Create Scaling Policy
echo "Attaching Scaling Policy..."
POLICY_ARN=$(aws autoscaling put-scaling-policy \
  --auto-scaling-group-name $ASG_NAME \
  --policy-name $POLICY_NAME \
  --scaling-adjustment 1 \
  --adjustment-type ChangeInCapacity \
  --region $REGION \
  --query PolicyARN --output text)

echo "Scaling Policy ARN: $POLICY_ARN"

# Create CloudWatch Alarm
echo "Linking Policy to CloudWatch Alarm..."
aws cloudwatch put-metric-alarm \
  --alarm-name $ALARM_NAME \
  --metric-name CPUUtilization \
  --namespace AWS/EC2 \
  --statistic Average \
  --period 300 \
  --threshold 70 \
  --comparison-operator GreaterThanThreshold \
  --dimensions Name=AutoScalingGroupName,Value=$ASG_NAME \
  --evaluation-periods 2 \
  --alarm-actions $POLICY_ARN \
  --region $REGION

echo "Auto Scaling Group setup complete!"
