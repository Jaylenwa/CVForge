package seed

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"openresume/internal/infra/database"
	"openresume/internal/pkg/logger"

	"go.uber.org/zap"
)

func RunCLI(args []string) int {
	if len(args) == 0 {
		_, _ = fmt.Fprintln(os.Stderr, "usage: openresume seed <import-default|import> [--file path]")
		return 2
	}

	switch args[0] {
	case "import-default":
		s, err := LoadDefaultSeed()
		if err != nil {
			logger.L().Error("load default seed failed", zap.Error(err))
			return 1
		}
		counts, err := Import(context.Background(), database.DB, s)
		if err != nil {
			logger.L().Error("seed import failed", zap.Error(err))
			return 1
		}
		_ = json.NewEncoder(os.Stdout).Encode(map[string]any{"success": true, "counts": counts})
		return 0

	case "import":
		fs := flag.NewFlagSet("openresume seed import", flag.ContinueOnError)
		fs.SetOutput(os.Stderr)
		filePath := fs.String("file", "", "seed json file path")
		if err := fs.Parse(args[1:]); err != nil {
			return 2
		}
		if *filePath == "" {
			_, _ = fmt.Fprintln(os.Stderr, "usage: openresume seed import --file path/to/seed.json")
			return 2
		}
		s, err := LoadFromFilePath(*filePath)
		if err != nil {
			logger.L().Error("load seed file failed", zap.Error(err))
			return 1
		}
		counts, err := Import(context.Background(), database.DB, s)
		if err != nil {
			logger.L().Error("seed import failed", zap.Error(err))
			return 1
		}
		_ = json.NewEncoder(os.Stdout).Encode(map[string]any{"success": true, "counts": counts})
		return 0

	default:
		_, _ = fmt.Fprintln(os.Stderr, "usage: openresume seed <import-default|import> [--file path]")
		return 2
	}
}
