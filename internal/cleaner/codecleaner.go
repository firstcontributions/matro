package cleaner

import "os"

var pathsToBeCleaned = []string{
	"internal/grpc",
	"assets",
	"api",
}

// Clean will clean existing code under given basePath all code should not be
// cleaned as there can be some custom code added by the user.
// For now tracking a  list of paths to be cleaned. later some configs files
// (.matroignore) like .gitignore can be intriduced
func Clean(basePath string) error {
	for _, path := range pathsToBeCleaned {
		if err := os.RemoveAll(basePath + "/" + path); err != nil {
			return err
		}
	}
	return nil
}
