## General 
This code implements the Action interface and provides a service that can perform the following tasks:

1. Install the server
2. Start the server
3. Stop the server
4. Uninstall the server
5. Get the current state of the server

## Code documentation:

### ActionService struct

The ActionService struct is responsible for managing the lifecycle and state synchronization of a minecraft server. It encapsulates the context management and state synchronization mechanisms utilized within the server's binary execution.

#### Fields

- **aliveCtx**:  
  Type: **context.Context**
  Description: This is the context that remains active while the server is running. It provides the ability to manage deadlines, cancellations, and other resource-scoped operations. It terminates when the server is intentionally stopped or encounters a failure.
  
- **stopCtx**:  
  Type: **context.CancelFunc** 
  Description: This function is used to stop the server binary execution safely. It enables controlled shutdowns, ensuring that all resources are released appropriately.
  
- syncedState:  
  Type: **mine_state.SyncedState**
  Description: A reference to the SyncedState channel object that holds the state of the server. It ensures that the server state is consistent and provides synchronization across different server activities.

### NewActionService function

The NewActionService function is designed to instantiate a new ActionService object. It initializes the service's state by obtaining a SyncedState object from the mine_state package, which tracks the operational status of the server. The current implementation sets the state to mine_state.Stopped.

#### Parameters

- None

#### Return Values

- **ActionService**:  
  A pointer to the newly created ActionService instance, initialized with a SyncedState set to mine_state.Stopped.

- **error**:  
  Returns nil in the current implementation since no errors are expected during initialization. However, it can potentially return an error if modifications introduce failure points during initialization.

### dowloadJar function

The downloadJar method is a member of the ActionService struct and is responsible for downloading a specified version of a Forge jar file from a given URL. It uses a context to manage the command execution, allowing the operation to be canceled or timeout.

#### Parameters

- ctx **context.Context**:  
  The context to control command execution. It allows the download operation to be canceled or to respect deadlines, which is essential for managing long-running processes or in network-limited environments.

- version **string**:  
  The version of the Forge jar file to download. This string is used to construct the proper URL and filename based on the template strings baseURL and installerName.

#### Return Value

- **error**:  
  Returns an error if the download operation fails. The error message is formatted to indicate the failed function and version. If the operation is successful, it returns nil.

### installJar function

The installJar method is a member of the ActionService struct, responsible for executing the installation of a Forge server jar file using the Java runtime. It leverages command-line execution to automate server setup for a specified version.

#### Parameters

- ctx **context.Context**:  
  This context controls the command execution. It provides capabilities for timeout management and cancellation, essential for handling long-running or potentially blocking operations.

- version **string**:  
  Specifies the version of the Forge jar file to be installed. This version is used to construct appropriate file paths for the installation process.

#### Return Value

- **error**:  
  Returns an error if the installation operation encounters an issue. The error message is formatted to include the invoking function and the version of the jar file. If successful, the method returns nil.

### removeJar function

The removeJar method is a member of the ActionService struct, responsible for removing a specific version of a server jar file from the filesystem. This method leverages command-line operations to delete files and manages execution through context.

#### Parameters

- ctx **context.Context**:
  The context for managing the lifecycle of the command execution. This context is used to support cancellations and time-based constraints, ensuring the operation is controlled effectively.

- version **string**:
  The version of the jar file to be removed. This parameter helps construct the appropriate path for the file targeted for deletion.

#### Return Value

- **error**:
  Returns an error if the removal operation fails. The error message is formatted to include the method name and version of the jar file for clarity. If the operation completes successfully, it returns nil.

### Install function

The Install method is a comprehensive function of the ActionService struct that manages the end-to-end process of downloading, installing, and cleaning up a Minecraft server with specific version. It orchestrates several subordinate methods to ensure the server is set up correctly.

#### Parameters

- ctx **context.Context**:  
  Context used to manage command execution lifecycles, supporting timeout, cancellation, and other resource management strategies across the installation processes.

- version **string**:  
  Represents the specific version of the Minecraft server to be installed. This version is used in constructing URLs and file paths necessary for downloading and installing the server software.

#### Return Value

- **error**:  
  Returns an error if any stage of the installation process fails. The error is wrapped with context indicating the stage (download, install, remove) where the failure occurred. If all operations are successful, it returns nil.

### Start function

The Start method of the ActionService struct is designed to initiate the execution of a server process. It ensures that the server is not already running before starting it and manages the server process in a separate goroutine for asynchronous operation.

#### Parameters

- ctx **context.Context**:  
  This is the parent context used for initial checks and setup. It is a good practice to have it available for any additional functional extensions or checks that might be necessary in the implementation.

#### Return Value

- **error**:  
  Returns an error if the server is already running. Otherwise, it returns nil indicating that the server has been successfully initiated.

### Stop function

The Stop method of the ActionService struct is responsible for halting the currently running server process. It achieves this by canceling the context controlling the server's execution.

#### Parameters

- ctx **context.Context**:  
  While this method includes the context parameter for consistency with other methods, it is primarily used as a formality in this implementation. The stopping action is achieved through the stopCtx function, which is set up internally in the ActionService.

#### Return Value

- **error**:  
  Returns nil as this method's operation is straightforward and does not currently account for error scenarios besides setup failures, which are generally managed elsewhere in the service lifecycle.

### GetState function

The GetState method in the ActionService struct is designed to retrieve the current state of the server. It provides an encapsulated way to check the server's status using the internal state management system.

#### Parameters

- ctx **context.Context**:  
  This method includes the context parameter for consistency with other methods, it is primarily used as a formality in this implementation. The stopping action is achieved through the stopCtx function, which is set up internally in the ActionService.

#### Return Value

- **mine_state.State**:  
  Returns the current state of the server, represented by the mine_state.State type. The state reflects whether the server is running, stopped, or in another defined condition, based on the internal state management mechanism.


