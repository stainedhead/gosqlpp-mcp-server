from aws_cdk import (
    Stack,
    aws_apprunner as apprunner,
    aws_ecr as ecr,
    aws_iam as iam,
    aws_logs as logs,
    CfnOutput,
    RemovalPolicy,
)
from constructs import Construct


class GosqlppMcpStack(Stack):
    def __init__(self, scope: Construct, construct_id: str, environment: str, **kwargs) -> None:
        super().__init__(scope, construct_id, **kwargs)

        self.environment = environment

        # Create ECR repository
        self.ecr_repository = self._create_ecr_repository()

        # Create IAM roles
        self.instance_role = self._create_instance_role()
        self.access_role = self._create_access_role()

        # Create CloudWatch log group
        self.log_group = self._create_log_group()

        # Create App Runner service
        self.app_runner_service = self._create_app_runner_service()

        # Create outputs
        self._create_outputs()

    def _create_ecr_repository(self) -> ecr.Repository:
        """Create ECR repository for container images"""
        repository = ecr.Repository(
            self,
            "GosqlppMcpRepository",
            repository_name=f"gosqlpp-mcp-server-{self.environment}",
            image_scan_on_push=True,
            lifecycle_rules=[
                ecr.LifecycleRule(
                    description="Keep last 10 images",
                    max_image_count=10,
                    rule_priority=1,
                )
            ],
            removal_policy=RemovalPolicy.DESTROY if self.environment == "development" else RemovalPolicy.RETAIN,
        )

        CfnOutput(
            self,
            "ECRRepositoryURI",
            value=repository.repository_uri,
            description="ECR Repository URI",
        )

        return repository

    def _create_instance_role(self) -> iam.Role:
        """Create IAM role for App Runner instance"""
        role = iam.Role(
            self,
            "GosqlppMcpInstanceRole",
            assumed_by=iam.ServicePrincipal("tasks.apprunner.amazonaws.com"),
            description="IAM role for gosqlpp MCP server App Runner instance",
        )

        # Add CloudWatch logs permissions
        role.add_to_policy(
            iam.PolicyStatement(
                effect=iam.Effect.ALLOW,
                actions=[
                    "logs:CreateLogGroup",
                    "logs:CreateLogStream",
                    "logs:PutLogEvents",
                    "logs:DescribeLogStreams",
                ],
                resources=[f"arn:aws:logs:{self.region}:{self.account}:log-group:/aws/apprunner/*"],
            )
        )

        return role

    def _create_access_role(self) -> iam.Role:
        """Create IAM role for App Runner to access ECR"""
        role = iam.Role(
            self,
            "GosqlppMcpAccessRole",
            assumed_by=iam.ServicePrincipal("build.apprunner.amazonaws.com"),
            description="IAM role for App Runner to access ECR",
        )

        # Add ECR permissions
        role.add_to_policy(
            iam.PolicyStatement(
                effect=iam.Effect.ALLOW,
                actions=[
                    "ecr:GetAuthorizationToken",
                    "ecr:BatchCheckLayerAvailability",
                    "ecr:GetDownloadUrlForLayer",
                    "ecr:BatchGetImage",
                ],
                resources=["*"],
            )
        )

        return role

    def _create_log_group(self) -> logs.LogGroup:
        """Create CloudWatch log group"""
        log_group = logs.LogGroup(
            self,
            "GosqlppMcpLogGroup",
            log_group_name=f"/aws/apprunner/gosqlpp-mcp-server-{self.environment}",
            retention=logs.RetentionDays.ONE_WEEK if self.environment == "development" else logs.RetentionDays.ONE_MONTH,
            removal_policy=RemovalPolicy.DESTROY if self.environment == "development" else RemovalPolicy.RETAIN,
        )

        return log_group

    def _create_app_runner_service(self) -> apprunner.CfnService:
        """Create App Runner service"""
        
        # Environment variables for the container
        environment_variables = [
            {
                "name": "GOSQLPP_MCP_SERVER_TRANSPORT",
                "value": "http"
            },
            {
                "name": "GOSQLPP_MCP_SERVER_HOST",
                "value": "0.0.0.0"
            },
            {
                "name": "GOSQLPP_MCP_SERVER_PORT",
                "value": "8080"
            },
            {
                "name": "GOSQLPP_MCP_LOG_LEVEL",
                "value": "debug" if self.environment == "development" else "info"
            },
            {
                "name": "GOSQLPP_MCP_LOG_FORMAT",
                "value": "json"
            },
            {
                "name": "GOSQLPP_MCP_AWS_REGION",
                "value": self.region
            },
            {
                "name": "GOSQLPP_MCP_AWS_ENVIRONMENT",
                "value": self.environment
            }
        ]

        service = apprunner.CfnService(
            self,
            "GosqlppMcpService",
            service_name=f"gosqlpp-mcp-server-{self.environment}",
            source_configuration=apprunner.CfnService.SourceConfigurationProperty(
                auto_deployments_enabled=True,
                image_repository=apprunner.CfnService.ImageRepositoryProperty(
                    image_identifier=f"{self.ecr_repository.repository_uri}:latest",
                    image_configuration=apprunner.CfnService.ImageConfigurationProperty(
                        port="8080",
                        runtime_environment_variables=environment_variables,
                        start_command="./gosqlpp-mcp-server --transport http --host 0.0.0.0 --port 8080",
                    ),
                    image_repository_type="ECR",
                ),
                access_role_arn=self.access_role.role_arn,
            ),
            instance_configuration=apprunner.CfnService.InstanceConfigurationProperty(
                cpu="0.25 vCPU",
                memory="0.5 GB",
                instance_role_arn=self.instance_role.role_arn,
            ),
            health_check_configuration=apprunner.CfnService.HealthCheckConfigurationProperty(
                protocol="HTTP",
                path="/health",
                interval=30,
                timeout=10,
                healthy_threshold=2,
                unhealthy_threshold=3,
            ),
            network_configuration=apprunner.CfnService.NetworkConfigurationProperty(
                egress_configuration=apprunner.CfnService.EgressConfigurationProperty(
                    egress_type="DEFAULT"
                )
            ),
        )

        return service

    def _create_outputs(self):
        """Create CloudFormation outputs"""
        CfnOutput(
            self,
            "AppRunnerServiceURL",
            value=f"https://{self.app_runner_service.attr_service_url}",
            description="App Runner service URL",
        )

        CfnOutput(
            self,
            "AppRunnerServiceArn",
            value=self.app_runner_service.attr_service_arn,
            description="App Runner service ARN",
        )

        CfnOutput(
            self,
            "LogGroupName",
            value=self.log_group.log_group_name,
            description="CloudWatch log group name",
        )
