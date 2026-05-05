package helm

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Run(ctx context.Context, dryRun bool, args ...string) error {
	if dryRun {
		fmt.Printf("      helm %s\n", strings.Join(args, " "))
		return nil
	}
	cmd := exec.CommandContext(ctx, "helm", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
