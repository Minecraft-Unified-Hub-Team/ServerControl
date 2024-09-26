# ServerControl

## About project
An open source project that provides an API and service for managing a single minecraft server.

## About project structure
### Directory Structure
The project directory is organized as follows:

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
- Purpose: This directory contains the main.go file which is the entry point of the application.
- Files:
  - main.go - The main entry point of the server application.

#### internal/
- Purpose: This directory hosts the internal packages that implement the business logic and core functionality of the application, safely hidden from external access.

##### internal/api/
- Purpose: Defines the API using Protocol Buffers.
- Files:
  - api.proto - Protocol Buffers definition file for API endpoints.

##### internal/serverControl/
- Purpose: Contains logic for server control functionalities such as handling actions, editing configurations, and health checks.
- Files:
  - serverControl.go - Main server control logic.
  - serverControlHandler.go - Handler for server control requests.
  - actionHandler.go - Handler for action-related requests.
  - configHandler.go - Handler for configuration-related requests.
  - healthHandler.go - Handler for health check requests.

##### internal/action/
- Purpose: Encapsulates the service logic related to actions.
- Files:
  - actionService.go - Service logic for handling actions.

##### internal/config/
- Purpose: Encapsulates the service logic related to configurations.
- Files:
  - configService.go - Service logic for handling configuration settings.

##### internal/health/
- Purpose: Encapsulates the service logic for health checks.
- Files:
  - healthService.go - Service logic for managing health checks of the server.

#### tests/
- Purpose: Contains testing-related files and BDD test scenarios.
- Files:
  - feature_manager.go - Manager to handle feature tests.
- Directories:
  - features/
    - Purpose: Contains feature files for BDD testing.
    - Files:
      - start.feature - A feature file describing test scenarios for starting the server.

## About ServerControl
This structure provides an API for managing the server, potentially able to handle complex scenarios by breaking them down into simpler ones. These simpler scenarios are implemented using "Services".
### About ActionService
The Action Service handles simple actions that the server can perform explicitly. For instance, it can start and stop server.
About functionality:
- StartCtx - starts the server with a non-blocking call, if the request is not processed during the lifetime of the context, it returns an error
- StopCtx - stops the server with a non-blocking call, if the request is not processed during the lifetime of the context, it returns an error
### About ConfigService
The Config Service handles scenarios where configurations need to be changed, affecting server settings. It performs actions such as modifying the usual server settings, managing mods, and plugins, among others.
- GetSettingsCtx - returns the list of server settings.
- UpdateSettingsCtx - updates the values in the server configuration.
### About HealthService
The Health Service can provide information about the current state of a server. This includes whether the server is running, has stopped, or dead due to an error.
- GetStateCtx - returns the current state of the server. For example, running, stopped, or dead.
### About LogService
TODO
### About BackupService
TODO
