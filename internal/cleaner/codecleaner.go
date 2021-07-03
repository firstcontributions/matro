package cleaner

import "os"

var pathsToBeCleaned = []string{
	"internal/grpc",
	"assets",
	"api",
}

func Clean(basePath string) error {
	for _, path := range pathsToBeCleaned {
		if err := os.RemoveAll(basePath + "/" + path); err != nil {
			return err
		}
	}
	return nil
}
