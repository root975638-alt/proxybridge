package aws

import (
	"os"
	"testing"
)

func TestAWSBedrockProvider_Name(t *testing.T) {
	p := NewAWSBedrockProvider()
	if p.Name() != "AWS Bedrock" {
		t.Errorf("expected name 'AWS Bedrock', got '%s'", p.Name())
	}
}

func TestAWSBedrockProvider_ID(t *testing.T) {
	p := NewAWSBedrockProvider()
	if p.ID() != "aws" {
		t.Errorf("expected ID 'aws', got '%s'", p.ID())
	}
}

func TestAWSBedrockProvider_GetDefaultModel(t *testing.T) {
	p := NewAWSBedrockProvider()
	expected := "anthropic.claude-3-5-sonnet-20241022-v2:0"
	if p.GetDefaultModel() != expected {
		t.Errorf("expected default model '%s', got '%s'", expected, p.GetDefaultModel())
	}
}

func TestAWSBedrockProvider_GetConfig(t *testing.T) {
	p := NewAWSBedrockProvider()
	config := p.GetConfig()
	if config["provider"] != "aws" {
		t.Errorf("expected provider 'aws', got '%s'", config["provider"])
	}
	if config["type"] != "bedrock" {
		t.Errorf("expected type 'bedrock', got '%s'", config["type"])
	}
}

func TestAWSBedrockProvider_GetEnvironment(t *testing.T) {
	p := NewAWSBedrockProvider()
	env := p.GetEnvironment()
	if env["AWS_ACCESS_KEY_ID"] != "YOUR_AWS_ACCESS_KEY_ID" {
		t.Errorf("expected AWS_ACCESS_KEY_ID placeholder")
	}
	if env["AWS_SECRET_ACCESS_KEY"] != "YOUR_AWS_SECRET_ACCESS_KEY" {
		t.Errorf("expected AWS_SECRET_ACCESS_KEY placeholder")
	}
}

func TestAWSBedrockProvider_IsInstalled(t *testing.T) {
	p := NewAWSBedrockProvider()
	if !p.IsInstalled() {
		t.Error("AWS Bedrock provider should be considered installed")
	}
}

func TestAWSBedrockProvider_Install(t *testing.T) {
	p := NewAWSBedrockProvider()
	if err := p.Install(); err != nil {
		t.Errorf("Install returned error: %v", err)
	}
}

func TestAWSBedrockProvider_Validate(t *testing.T) {
	// Set up temporary environment variables
	os.Setenv("AWS_ACCESS_KEY_ID", "test-key")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "test-secret")
	os.Setenv("AWS_REGION", "us-east-1")
	defer func() {
		os.Unsetenv("AWS_ACCESS_KEY_ID")
		os.Unsetenv("AWS_SECRET_ACCESS_KEY")
		os.Unsetenv("AWS_REGION")
	}()

	p := NewAWSBedrockProvider()
	if err := p.Validate(); err != nil {
		t.Errorf("Validate returned error: %v", err)
	}
}

func TestAWSBedrockProvider_IsValidAWSRegion(t *testing.T) {
	p := NewAWSBedrockProvider()

	validRegions := []string{"us-east-1", "us-west-2", "eu-west-1", "ap-northeast-1"}
	for _, region := range validRegions {
		if !p.isValidAWSRegion(region) {
			t.Errorf("expected region '%s' to be valid", region)
		}
	}

	invalidRegions := []string{"invalid-region", "us-fake-1", ""}
	for _, region := range invalidRegions {
		if p.isValidAWSRegion(region) {
			t.Errorf("expected region '%s' to be invalid", region)
		}
	}
}
