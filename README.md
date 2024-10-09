# ServerControl

## About project
An open source project that provides an API and a service to manage a single Minecraft server.

## About project structure
### Directory Structure
The project is organized as follows:

```
├── cmd/
    └── main.go
├── internal/
    ├── api/
        └── api.proto
    ├── serverControl/
        ├── serverControl.go
        ├── serverControlHandler.go
        ├── actionHandler.go
        ├── configHandler.go
        └── healthHandler.go
    ├── install/
        └── installService.go
    ├── action/
        └── actionService.go
    ├── config/
        └── configService.go
    └── health/
        └── healthService.go
└── tests/
    ├── feature_manager.go
    └── features/
        └── start.feature
```

#### cmd/
- Purpose: This directory contains main.go file which is the entry point of the application.
- Files:
  - main.go - The main entry point of the server application.

#### internal/
- Purpose: This directory hosts the internal packages that implement logic and core functionality of the application.

##### internal/api/
- Purpose: Defines the API using Protocol Buffers.
- Files:
  - api.proto - Protocol Buffers definition file for API endpoints.

##### internal/serverControl/
- Purpose: Contains logic to control the server functions, such as handling actions, editing configurations, and health checks.
- Files:
  - serverControl.go - Main server control logic.
  - serverControlHandler.go - Handler for server control requests.
  - actionHandler.go - Handler for action-related requests.
  - configHandler.go - Handler for configuration-related requests.
  - healthHandler.go - Handler for health check requests.

##### internal/action/
- Purpose: Encapsulates the service logic related to actions.
- Files:
  - actionService.go - Service logic to handle actions.

##### internal/config/
- Purpose: Encapsulates the service logic related to configurations.
- Files:
  - configService.go - Service logic to handle configuration settings.

##### internal/health/
- Purpose: Encapsulates the service logic for health checks.
- Files:
  - healthService.go - Service logic to manage health checks of the server.

#### tests/
- Purpose: Contains testing-related files and BDD test scenarios.
- Files:
  - feature_manager.go - Manager that handles feature tests.
- Directories:
  - features/
    - Purpose: Contains feature files for BDD testing.
    - Files:
      - start.feature - A feature file describing test scenarios for starting the server.

### About the ServerControl architecture
The purpose of ServerControl is to efficiently manage a Minecraft server running in a Docker container. It provides gRPC API to interact and change the server properties and states. 
ServerControl receives requests, breaks them down into simple ones, and delegates them to the controlled service units.

Controlled service units include:
- Action
- Config
- Health
- Backup
- Log

#### About ActionService
This service unit performs actions that affect server availability, including installation, starting, stopping, and uninstalling.

#### About ConfigService
This service unit performs actions on configuration files, including creating and modifying them.
In the case of user_jvm_args.txt configuration file, it calculates the properties of the container in which it is located.

#### About Health
This service unit performs actions that retrieve information about the server's state from the environment or other services.
