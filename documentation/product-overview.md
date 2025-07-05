I am a software engineer who has a commandline application that enables connections and sending SQL commands to relational databases.  I would like to enable usage of this application via agent based development tooling.  To do this, I would like to implement an MCP servr interface that provides access to this CLI, and leverages specific features within the application.

The key feature of this mcp server is to execute the sqlpp application as a child process CLI which is not a visible to logged in users.  The parent application will wait for the results of the sqlpp and stream them back to the CLI client who made the request.  Each tool action enabled will provide inputs which will assist the control of the sqlpp application and its output.

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

Excuting child processes are documented within: https://pkg.go.dev/os/exec
A forum thread also documents this within: https://forum.golangbridge.org/t/running-go-app-without-showing-the-terminal/19361


CDK deployment example: https://github.com/aws-samples/cdk-apprunner-ecr
App Runner and ECR documentation: https://docs.aws.amazon.com/toolkit-for-vscode/latest/userguide/ecr-apprunner.html

information on sqlpp is found in https://github.com/stainedhead/gosqlpp .  https://github.com/stainedhead/gosqlpp/README.md and https://github.com/stainedhead/gosqlpp/documentation/product-overview may be helpful.

### Features 
0. Consider the configuration values which will be needed, create a configuration file for local use
1. Ensure we wrap the @name shortcuts to send the appropriete command and stream the results back to the caller.  Each of these actions will provide a connection parameter, a name-filter parameter and an output parameter, to control the connection we are inteacting with, the filter they want to use on the results and the format of the output returned.  These should all map to the flags available on the command line.  The @name shortcuts to map to are @schema-all, @schema-tables, @schema-views, @schema-procedures, @schema-functions, the tool action name will be the same without the leading @ charactor.  Meaning @schema-tables becomes schema-tables and all have three string parameters of connection, filter and output which map the CLI flags to be used when executing them.
2. Ensure we wrap the --list-connections commandline flag which returns a list of the configured connections.  this will allow the cli client to query on the available connections and to know the correct names to send in the connection parameter to tool requests.
3. Ensure we wrap execution of SQL commands to the database for processing. This action will be called execute_sql_command and take three string parameters connection, command and output.  They will allow the caller to control the database connection to use when sending the command, the SQL command that should be executed, including multiple commands seperated by go statements, and the output format results should be returned in.
4. Ensure we wrap the @drivers shortcut, which returns the list of drivers the sqlpp application provides.  This allows the caller 
4. Ensure we update the README.md to be a high quality and professional file which would help later adopters use, understand the implementation and configuration of the mcp server.
5. Ensure we add documentation to ./documentation/product-generation-results.md of the technical details related to the code that was generated.