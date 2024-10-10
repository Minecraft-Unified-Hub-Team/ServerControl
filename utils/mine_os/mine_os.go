package mine_os

import (
	"context"
	"fmt"
	"os/exec"
	"syscall"
)

const (
	SYSTEM_ERROR = iota
	COMMON_ERROR
	NO_ERROR
)

func ExecCtx(ctx context.Context, command string, args []string) error {
	var err error = nil
	var errorFormat string = fmt.Sprintf("mine_os.ExecCtx(ctx, %s, %v):", command, args) + ": %w"

	cmd := exec.CommandContext(ctx, command, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	return err
}

func ManagedExecCtx(ctx context.Context, command string, args []string) (int, error) {
	var status int = NO_ERROR
	var err error = nil
	var errorFormat string = fmt.Sprintf("mine_os.ManagedExecCtx(ctx, %s, %v)", command, args) + ": %w"

	/* init managed process cmd */
	cmd := exec.CommandContext(ctx, command, args...)
	/* init custom pid group for current cmd and its childrens */
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	/* init sync interrupt channel */
	termChan := make(chan struct{}, 1)

	err = cmd.Start()
	if err != nil {
		status = SYSTEM_ERROR /* error equal to start failed (uncommon error) */
		return status, fmt.Errorf(errorFormat, err)
	}

	/* get group processes pid */
	groupID, err := syscall.Getpgid(cmd.Process.Pid)
	if err != nil {
		status = SYSTEM_ERROR                                /* error equal to start failed (uncommon error) */
		err = syscall.Kill(cmd.Process.Pid, syscall.SIGKILL) /* kill main process */
		return status, fmt.Errorf(errorFormat, err)
	}

	/* start goroutine for killing child processes (can be triggered by termChan) */
	go func() {
		<-termChan
		syscall.Kill(-groupID, syscall.SIGTERM) /* send SIGTERM to process group */
	}()

	/* wait until command finish or fail */
	err = cmd.Wait()
	select {
	case <-ctx.Done():
		status = NO_ERROR /* no error, just stopped by context cancelation */
	default:
		status = COMMON_ERROR /* error during execution (common error) */
	}

	/* trigger to stop child processes */
	termChan <- struct{}{}
	if status == NO_ERROR {
		return status, nil
	} else {
		return status, fmt.Errorf(errorFormat, err)
	}
}
