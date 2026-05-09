package util

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func RunCmd(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func RunCmdOutput(ctx context.Context, name string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stderr = os.Stderr
	return cmd.Output()
}

func RunCmdStdin(ctx context.Context, stdin io.Reader, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdin = stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func DryRunOrExec(ctx context.Context, dryRun bool, name string, args ...string) error {
	if dryRun {
		fmt.Printf("      %s %s\n", name, strings.Join(args, " "))
		return nil
	}
	return RunCmd(ctx, name, args...)
}

func DryRunOrExecStdin(ctx context.Context, dryRun bool, stdin io.Reader, name string, args ...string) error {
	if dryRun {
		fmt.Printf("      %s %s\n", name, strings.Join(args, " "))
		return nil
	}
	return RunCmdStdin(ctx, stdin, name, args...)
}
