package helm

import (
	"context"
	"path/filepath"
	"strings"

	"k8s-manager/internal/logger"
	"k8s-manager/internal/util"
)

func Run(ctx context.Context, dryRun bool, args ...string) error {
	normalized := make([]string, len(args))
	for i, a := range args {
		normalized[i] = filepath.ToSlash(a)
	}
	if dryRun {
		logger.Step("      helm %s", strings.Join(normalized, " "))
		return nil
	}
	return util.RunCmd(ctx, "helm", normalized...)
}
