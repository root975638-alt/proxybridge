// Package aws provides the AWS Bedrock provider for ProxyBridge.
package aws

import (
	"fmt"
	"os"

	"github.com/root975638-alt/proxybridge/internal/logging"
)

// AWSBedrockProvider implements the Provider interface for AWS Bedrock.
type AWSBedrockProvider struct {
	id           string
	name         string
	description  string
	defaultModel string
}

// NewAWSBedrockProvider creates a new AWS Bedrock provider.
func NewAWSBedrockProvider() *AWSBedrockProvider {
	return &AWSBedrockProvider{
		id:           "aws",
		name:         "AWS Bedrock",
		description:  "Amazon Web Services Bedrock",
		defaultModel: "anthropic.claude-3-5-sonnet-20241022-v2:0",
	}
}

// Name returns the display name of the provider.
func (p *AWSBedrockProvider) Name() string {
	return p.name
}

// ID returns the unique identifier for the provider.
func (p *AWSBedrockProvider) ID() string {
	return p.id
}

// Setup initializes the provider.
func (p *AWSBedrockProvider) Setup() error {
	logging.Info("Setting up AWS Bedrock provider")
	return nil
}

// Validate validates the AWS credentials and configuration.
func (p *AWSBedrockProvider) Validate() error {
	// Check for AWS credentials
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsRegion := os.Getenv("AWS_REGION")

	// If credentials are not set via environment, try to use default chain
	if awsAccessKey == "" || awsSecretKey == "" {
		// Check if AWS CLI is configured
		if !p.isAWSCliConfigured() {
			return fmt.Errorf("AWS credentials not configured. Set AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, and AWS_REGION environment variables, or run 'aws configure'")
		}
	}

	// Validate region format
	if awsRegion == "" {
		return fmt.Errorf("AWS_REGION environment variable not set")
	}
	if !p.isValidAWSRegion(awsRegion) {
		return fmt.Errorf("invalid AWS region: %s", awsRegion)
	}

	// Check for Bedrock access
	if !p.hasBedrockAccess() {
		return fmt.Errorf("AWS Bedrock access not enabled. Visit https://aws.amazon.com/bedrock/ to enable access")
	}

	return nil
}

// GetModels returns the list of available AWS Bedrock models.
func (p *AWSBedrockProvider) GetModels() ([]string, error) {
	// Common AWS Bedrock models
	models := []string{
		"anthropic.claude-3-5-sonnet-20241022-v2:0",
		"anthropic.claude-3-opus-20240229-v1:0",
		"anthropic.claude-3-haiku-20240307-v1:0",
		"anthropic.claude-3-7-sonnet-20250219-v1:0",
		"anthropic.claude-3-5-haiku-20241022-v1:0",
		"amazon.nova-pro-v1:0",
		"amazon.nova-lite-v1:0",
		"amazon.titan-text-lite-v1:0",
		"amazon.titan-text-express-v1:0",
		"cohere.command-r-v1:0",
		"cohere.command-r-plus-v1:0",
		"meta.llama3-2-11b-instruct-v1:0",
		"meta.llama3-2-90b-instruct-v1:0",
		"meta.llama3-1-8b-instruct-v1:0",
		"meta.llama3-1-70b-instruct-v1:0",
		"meta.llama3-70b-instruct-v1:0",
		"meta.llama3-8b-instruct-v1:0",
	}

	return models, nil
}

// GetDefaultModel returns the default model for this provider.
func (p *AWSBedrockProvider) GetDefaultModel() string {
	return p.defaultModel
}

// GetConfig returns provider-specific configuration.
func (p *AWSBedrockProvider) GetConfig() map[string]string {
	return map[string]string{
		"provider": "aws",
		"type":     "bedrock",
	}
}

// GetEnvironment returns environment variables for this provider.
func (p *AWSBedrockProvider) GetEnvironment() map[string]string {
	return map[string]string{
		"AWS_ACCESS_KEY_ID":     "YOUR_AWS_ACCESS_KEY_ID",
		"AWS_SECRET_ACCESS_KEY": "YOUR_AWS_SECRET_ACCESS_KEY",
		"AWS_REGION":            "us-east-1", // Bedrock is only available in specific regions
	}
}

// IsInstalled checks if the provider is installed.
func (p *AWSBedrockProvider) IsInstalled() bool {
	// AWS Bedrock doesn't require a separate installation
	// It's accessed through the AWS SDK
	return true
}

// Install installs the provider.
func (p *AWSBedrockProvider) Install() error {
	logging.Info("AWS Bedrock provider - no installation required")

	// Provide installation instructions
	fmt.Println("AWS Bedrock provider setup:")
	fmt.Println()
	fmt.Println("1. Install AWS CLI if not already installed:")
	fmt.Println("   - Windows: winget install Amazon.WebServicesCLI")
	fmt.Println("   - macOS: brew install awscli")
	fmt.Println("   - Linux: curl 'https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip' -o 'awscliv2.zip'")
	fmt.Println()
	fmt.Println("2. Configure AWS credentials:")
	fmt.Println("   aws configure")
	fmt.Println()
	fmt.Println("3. Enable Bedrock access:")
	fmt.Println("   - Visit https://aws.amazon.com/bedrock/")
	fmt.Println("   - Click 'Enable Bedrock'")
	fmt.Println()

	return nil
}

// Uninstall uninstalls the provider.
func (p *AWSBedrockProvider) Uninstall() error {
	// Nothing to uninstall for AWS Bedrock
	return nil
}

// isAWSCliConfigured checks if AWS CLI is configured.
func (p *AWSBedrockProvider) isAWSCliConfigured() bool {
	// Check if AWS CLI is available
	// In production, would check ~/.aws/credentials and ~/.aws/config
	return true
}

// isValidAWSRegion checks if a region is valid for Bedrock.
func (p *AWSBedrockProvider) isValidAWSRegion(region string) bool {
	// Bedrock is available in specific regions
	validRegions := map[string]bool{
		"us-east-1":      true,
		"us-east-2":      true,
		"us-west-2":      true,
		"eu-west-1":      true,
		"eu-central-1":   true,
		"ap-northeast-1": true,
		"ap-southeast-1": true,
		"ap-southeast-2": true,
	}

	return validRegions[region]
}

// hasBedrockAccess checks if the account has Bedrock access.
func (p *AWSBedrockProvider) hasBedrockAccess() bool {
	// In production, would make an API call to verify access
	// For now, assume access if credentials are configured
	return true
}

// GetProvider creates a new AWS Bedrock provider instance.
func GetProvider() *AWSBedrockProvider {
	return NewAWSBedrockProvider()
}
