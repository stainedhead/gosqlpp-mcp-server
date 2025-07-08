I am a software engineer who has a commandline application that enables connections and sending SQL commands to relational databases.  I would like to enable usage of this application via agent based development tooling.  To do this, I would like to implement an MCP servr interface that provides access to this CLI, and leverages specific features within the application.

The key feature of this mcp server is to execute the sqlpp application as a child process CLI which is not a visible to logged in users.  The parent application will wait for the results of the sqlpp and stream them back to the CLI client who made the request.  Each tool action enabled will provide inputs which will assist the control of the sqlpp application and its output.

Can you assist me in the development of this application.  I would like you to review the feature I will list below, formulate a plan for development, show me that plan to confirm it is acceptable or needs to be updated and when ready have you generate the code.

Your job is to assist me in the development of this application.  I would like you to review the feature I will list below, formulate a plan for development, show me that plan to confirm it is acceptable or needs to be updated and when ready have you generate the code.

 

### Technical Details
This application will be written in golang, and be containerized for deployement
The application will be an MCP server that provides tool actions to the caller.
The application will be deployed to AWS App Runner in test and production
AWS CDK will be leveraged to deploy any infrastructure
GitHub Actions will be used to build, test and deploy the application with CDK to assist
AWS has been configured to use OIDC based GitHub Actions
us-east-1 is our default region, should be configurable
Deployment commands from our dev environment should be provided

### Examples and Context
The core package that will enable the MCP interaction is: https://github.com/modelcontextprotocol/go-sdk
Importable packages are available at: github.com/modelcontextprotocol/go-sdk/mcp
and github.com/modelcontextprotocol/go-sdk/jsonschema

Other packages should be selected to support commandline flags and ENV variables and configuration file settings.  Select core framework packages if available, then select the most popular packages within the golang community if required.  This would also include logging.

The MCP protocol is documented within: https://docs.anthropic.com/en/docs/mcp
documentation on MCP tools can be found: https://modelcontextprotocol.io/docs/concepts/tools
a simple MCP example can be found: https://github.com/modelcontextprotocol/go-sdk/blob/main/examples/hello/main.go


Excuting child processes are documented within: https://pkg.go.dev/os/exec
A forum thread also documents this within: https://forum.golangbridge.org/t/running-go-app-without-showing-the-terminal/19361
CDK deployment example: https://github.com/aws-samples/cdk-apprunner-ecr
App Runner and ECR documentation: https://docs.aws.amazon.com/toolkit-for-vscode/latest/userguide/ecr-apprunner.html

information on sqlpp is found in https://github.com/stainedhead/gosqlpp .  https://github.com/stainedhead/gosqlpp/README.md and https://github.com/stainedhead/gosqlpp/documentation/product-overview may be helpful.

### Features 
0. Consider the configuration values which will be needed, create a configuration file for local use
1. Provide the standard MCP Server protocol to allow MCP Clients to interact with this server.
2. Ensure you support STDIO and HTTP+SSE to provide flexibility on usage and testing
3. Implement unit tests that ensure the system is stable and working as expected
4. Ensure we wrap the @name shortcuts to send the appropriete command and stream the results back to the caller.  Each of these actions will provide a connection parameter, a name-filter parameter and an output parameter, to control the connection we are inteacting with, the filter they want to use on the results and the format of the output returned.  These should all map to the flags available on the command line.  The @name shortcuts to map to are @schema-all, @schema-tables, @schema-views, @schema-procedures, @schema-functions, the tool action name will follow MCP standard {verb}_{noun} format with list_ prefix.  Meaning @schema-tables becomes list_schema_tables and all have three string parameters of connection, filter and output which map the CLI flags to be used when executing them.
5. Ensure we wrap the --list-connections commandline flag which returns a list of the configured connections.  this will allow the cli client to query on the available connections and to know the correct names to send in the connection parameter to tool requests. The tool will be named list_connections.
6. Ensure we wrap execution of SQL commands to the database for processing. This action will be called execute_sql_command and take three string parameters connection, command and output.  They will allow the caller to control the database connection to use when sending the command, the SQL command that should be executed, including multiple commands seperated by go statements, and the output format results should be returned in.
7. Ensure we wrap the @drivers shortcut, which returns the list of drivers the sqlpp application provides. This will be named list_drivers following MCP standards. 
8. Ensure we update the README.md to be a high quality and professional file which would help later adopters use, understand the implementation and configuration of the mcp server.
9. Ensure we add documentation to ./documentation/product-generation-results.md of the technical details related to the code that was generated.

