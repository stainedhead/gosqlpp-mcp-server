#!/usr/bin/env python3
import os
import aws_cdk as cdk
from stacks.gosqlpp_mcp_stack import GosqlppMcpStack

app = cdk.App()

# Get environment variables
account = os.environ.get('CDK_DEFAULT_ACCOUNT')
region = os.environ.get('CDK_DEFAULT_REGION', 'us-east-1')
environment = os.environ.get('ENVIRONMENT', 'development')

# Create stack
GosqlppMcpStack(
    app, 
    f"GosqlppMcpStack-{environment}",
    env=cdk.Environment(account=account, region=region),
    environment=environment
)

app.synth()
