package mothership_worker_server

import (
	"github.com/pkg/errors"
	"go.temporal.io/sdk/temporal"
	"net/url"
	"strings"
)

func getObjectPath(uri string) (string, error) {
	// Get object name from URI.
	// Check if object exists.
	// If not, return error.
	parsed, err := url.Parse(uri)
	if err != nil {
		return "", temporal.NewNonRetryableApplicationError(
			"could not parse resource URI",
			"couldNotParseResourceURI",
			errors.Wrap(err, "failed to parse resource URI"),
		)
	}

	// S3 for example must include bucket, while memory:// does not.
	// So memory://test.rpm would be parsed as host=test.rpm, path="".
	// While s3://mship/test.rpm would be parsed as host=mship, path=test.rpm.
	object := strings.TrimPrefix(parsed.Path, "/")
	if object == "" {
		object = parsed.Host
	}

	return object, nil
}
