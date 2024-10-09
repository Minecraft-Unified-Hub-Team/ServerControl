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

### About the ServerControl architecture
The purpose of ServerControl is to manage server in docker container. Besides this, it's gives (grpc) API for 
interaction and change server properties and state. So ServerControl just a "frontend", so it can recieve queries
and break them down into simple ones. After that, it sends "simple" requests to the controlled service units. 
It is important to note that the set of simple queries is quite limited, but at the same time, it is as comprehensive 
as possible. That is, using these simple queries, one can "construct" almost any more complex query.

Service units that can perform "simple" queries:
- Action
- Config
- Health
- Backup
- Log

#### About ActionService
This service unit performs active actions that are directly noticeable on the availability status of the server. 
Such as installation, start, stop, uninstall and other "simple " actions

#### About ConfigService
This service unit performs actions with different configs. For example, creating or changing a config on request
from the API. In the case of user_jvm_args.txt config, it calculates the properties of the container in which 
it is located.

#### About Health
This service unit performs actions that fetch info about server state from environment or other service.
