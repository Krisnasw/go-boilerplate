//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/magefile/mage/mg"
)

// Gen namespace for proto generation tasks
type Gen mg.Namespace

// Proto generates protobuf files for a specified service or all app.
func (Gen) Proto(service string) error {
	if service == "all" {
		return generateAllProtoFiles()
	}

	return generateProtoFilesForService(service)
}

// generateProtoFilesForService generates protobuf files for a specific service.
func generateProtoFilesForService(service string) error {
	basePath := fmt.Sprintf("../contracts/%s", service)

	// Check if the directory exists
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		return fmt.Errorf("contracts directory for service '%s' not found: %s", service, basePath)
	}

	// Find all .proto files directly in the contracts directory
	protoFiles, err := filepath.Glob(fmt.Sprintf("%s/*.proto", basePath))
	if err != nil || len(protoFiles) == 0 {
		return fmt.Errorf("no .proto files found in %s", basePath)
	}

	// Loop through and generate protobuf files for each .proto file
	for _, protoFile := range protoFiles {
		fmt.Printf("Generating protobuf for: %s\n", protoFile)

		cmd := exec.Command(
			"protoc",
			"-I", basePath,
			"--go_out", basePath,
			"--go_opt", "paths=source_relative",
			"--go-grpc_out", basePath,
			"--go-grpc_opt", "paths=source_relative",
			"--grpc-gateway_out", basePath,
			"--grpc-gateway_opt", "paths=source_relative",
			"--openapiv2_out", basePath,
			"--openapiv2_opt", "use_go_templates=true",
			protoFile,
		)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to generate protobuf files for %s: %w", protoFile, err)
		}
	}

	fmt.Printf("Protobuf files generated successfully for service '%s'\n", service)
	return nil
}

// generateAllProtoFiles generates protobuf files for all services directly under the contracts folder.
func generateAllProtoFiles() error {
	contractsDir := "../contracts"

	// Check if the directory exists
	if _, err := os.Stat(contractsDir); os.IsNotExist(err) {
		return fmt.Errorf("contracts directory not found: %s", contractsDir)
	}

	// Find all .proto files directly in the contracts directory
	protoFiles, err := filepath.Glob(fmt.Sprintf("%s/*.proto", contractsDir))
	if err != nil || len(protoFiles) == 0 {
		return fmt.Errorf("no .proto files found in %s", contractsDir)
	}

	// Loop through and generate protobuf files for each .proto file
	for _, protoFile := range protoFiles {
		fmt.Printf("Generating protobuf for: %s\n", protoFile)

		basePath := filepath.Dir(protoFile)
		cmd := exec.Command(
			"protoc",
			"-I", contractsDir,
			"--go_out", basePath,
			"--go_opt", "paths=source_relative",
			"--go-grpc_out", basePath,
			"--go-grpc_opt", "paths=source_relative",
			"--grpc-gateway_out", basePath,
			"--grpc-gateway_opt", "paths=source_relative",
			"--openapiv2_out", basePath,
			"--openapiv2_opt", "use_go_templates=true",
			protoFile,
		)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to generate protobuf files for %s: %w", protoFile, err)
		}
	}

	fmt.Println("Protobuf generation completed for all contracts!")
	return nil
}
