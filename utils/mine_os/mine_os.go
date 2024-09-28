package mine_os

import (
	"context"
	"os/exec"
	"syscall"
)

// ManagedExecCtx function
// Desc:
//
//	implements upgraded exec command with custom context support
//	and managing child processes (kill them if main process died)
//
// Return:
//
//	(status, error) - if error happens returns error
//	and status:
//	-1 - unexpected error (during start),
//	0 - expected error (during execution),
//	1 - no errors
func ManagedExecCtx(ctx context.Context, command string, args []string) (status int, err error) {
	/* init status variable with value 1 (OK) and err with value nil */
	status, err = 1, nil

	/* init managed process cmd */
	cmd := exec.CommandContext(ctx, command, args...)
	/* init custom pid group for current cmd and its childrens */
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	/* init sync interrupt channel */
	termChan := make(chan struct{}, 1)

	/* start cmd */
	err = cmd.Start()
	if err != nil {
		status = -1 /* error equal to start failed (uncommon error) */
		return
	}

	/* get group processes pid */
	groupID, err := syscall.Getpgid(cmd.Process.Pid)
	if err != nil {
		status = -1                                          /* error equal to start failed (uncommon error) */
		err = syscall.Kill(cmd.Process.Pid, syscall.SIGKILL) /* kill main process */
		return
	}

	/* start goroutine for killing child processes (triggered by termChan) */
	go func() {
		<-termChan
		syscall.Kill(-groupID, syscall.SIGTERM) /* send SIGTERM to process group */
	}()

	/* wait until command finish or fail */
	err = cmd.Wait()
	select {
	case <-ctx.Done():
		status = 1 /* no error, just stopped by context cancelation */
	default:
		status = 0 /* error during execution (common error) */
	}
	/* trigger to stop child processes */
	termChan <- struct{}{}
	return
}

// ManagedExec function
// Desc:
//
//	implements upgraded exec command with background context support
//	and managing child processes (kill them if main process died)
//
// Return:
//
//	(status, error) - if error happens returns error
//	and status:
//	-1 - unexpected error (during start),
//	0 - expected error (during execution),
//	1 - no errors
func ManagedExec(command string, args []string) (int, error) {
	return ManagedExecCtx(context.Background(), command, args)
}
