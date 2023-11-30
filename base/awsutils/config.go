package awsutils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"os"
)

// FillOutConfig fills out the config with the default values.
// Supports AWS CLI related environment variables.
func FillOutConfig(cfg *aws.Config) {
	// RESF defaults to region us-east-2
	if cfg.Region == nil || *cfg.Region == "" {
		region := "us-east-2"
		if envRegion := os.Getenv("AWS_REGION"); envRegion != "" {
			region = envRegion
		}
		cfg.Region = aws.String(region)
	}

	// Default to AWS S3 endpoint
	if cfg.Endpoint == nil || *cfg.Endpoint == "" {
		if envEndpoint := os.Getenv("AWS_ENDPOINT"); envEndpoint != "" {
			cfg.Endpoint = aws.String(envEndpoint)
		}
	}

	// By default, only connect securely
	if cfg.DisableSSL == nil {
		// But allow disabling it
		if envDisableSSL := os.Getenv("AWS_DISABLE_SSL"); envDisableSSL != "" {
			cfg.DisableSSL = aws.Bool(envDisableSSL == "true")
		} else {
			cfg.DisableSSL = aws.Bool(false)
		}
	}

	// By default, do not use path style
	if cfg.S3ForcePathStyle == nil {
		// But allow enabling it
		if envS3ForcePathStyle := os.Getenv("AWS_S3_FORCE_PATH_STYLE"); envS3ForcePathStyle != "" {
			cfg.S3ForcePathStyle = aws.Bool(envS3ForcePathStyle == "true")
		} else {
			cfg.S3ForcePathStyle = aws.Bool(false)
		}
	}

	// By default, use dualstack
	if cfg.UseDualStackEndpoint == endpoints.DualStackEndpointStateUnset {
		// But allow disabling it
		if envUseDualStackEndpoint := os.Getenv("AWS_USE_DUALSTACK_ENDPOINT"); envUseDualStackEndpoint != "" {
			if envUseDualStackEndpoint == "false" {
				cfg.UseDualStackEndpoint = endpoints.DualStackEndpointStateDisabled
			} else {
				cfg.UseDualStackEndpoint = endpoints.DualStackEndpointStateEnabled
			}
		} else {
			cfg.UseDualStackEndpoint = endpoints.DualStackEndpointStateEnabled
		}
	}

	// Allow setting access key and secret key with environment variables
	if cfg.Credentials == nil {
		// Make sure both ACCESS_KEY_ID and SECRET_ACCESS_KEY are set
		// Otherwise, do not set credentials
		envAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
		envSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
		if envAccessKeyID != "" && envSecretAccessKey != "" {
			cfg.Credentials = credentials.NewStaticCredentials(envAccessKeyID, envSecretAccessKey, "")
		}
	}
}
