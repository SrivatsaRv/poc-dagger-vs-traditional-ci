// A generated module for PocDaggerVsTraditionalCi functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"dagger/poc-dagger-vs-traditional-ci/internal/dagger"
)

type PocDaggerVsTraditionalCi struct{}

// Returns a container that echoes whatever string argument is provided
func (m *PocDaggerVsTraditionalCi) ContainerEcho(stringArg string) *dagger.Container {
	return dag.Container().From("alpine:latest").WithExec([]string{"echo", stringArg})
}

// Returns lines that match a pattern in the files of the provided Directory
func (m *PocDaggerVsTraditionalCi) GrepDir(ctx context.Context, directoryArg *dagger.Directory, pattern string) (string, error) {
	return dag.Container().
		From("alpine:latest").
		WithMountedDirectory("/mnt", directoryArg).
		WithWorkdir("/mnt").
		WithExec([]string{"grep", "-R", pattern, "."}).
		Stdout(ctx)
}

// Runs golangci-lint inside a container and lints the current directory
func (m *PocDaggerVsTraditionalCi) Lint(ctx context.Context, directoryArg *dagger.Directory) (string, error) {
	return dag.Container().
		From("golangci/golangci-lint:latest").
		WithMountedDirectory("/app", directoryArg.Directory(("todolist-app"))).
		WithWorkdir("/app").
		WithExec([]string{"golangci-lint", "run", "./..."}).
		Stdout(ctx)
}

// Runs Go tests inside a container and returns the test output
func (m *PocDaggerVsTraditionalCi) Test(ctx context.Context, directoryArg *dagger.Directory) (string, error) {
	return dag.Container().
		From("golang:latest").                                                // Use the official Go image
		WithMountedDirectory("/app", directoryArg.Directory("todolist-app")). // Mount the todolist-app directory
		WithWorkdir("/app").                                                  // Set /app as the working directory
		WithExec([]string{"go", "test", "-v", "./..."}).                      // Run the Go test command
		Stdout(ctx)                                                           // Capture and return the output
}

// // Runs the application by executing `docker compose up --build -d` inside a container     --------------------------->>>> THIS WAS CHALLENGING
// func (m *PocDaggerVsTraditionalCi) AppUp(ctx context.Context, directoryArg *dagger.Directory) error {
// 	dockerSocketPath := "/var/run/docker.sock" // Path to the Docker socket on the host

// 	_, err := dag.Container().
// 		From("docker:latest").                                                     // Use Docker's official image
// 		WithUnixSocket(dockerSocketPath, dag.Host().UnixSocket(dockerSocketPath)). // Mount Docker socket into container
// 		WithMountedDirectory("/app", directoryArg.Directory("todolist-app")).      // Mount the todolist-app directory
// 		WithWorkdir("/app").                                                       // Set working directory
// 		WithExec([]string{"docker-compose", "up", "--build", "-d"}).               // Execute docker-compose command
// 		Stdout(ctx)                                                                // Capture and log output
// 	return err
// }

// Builds a Go binary for the application and returns the container for further actions
func (m *PocDaggerVsTraditionalCi) Build(ctx context.Context, source *dagger.Directory) *dagger.Container {
	return dag.Container().
		From("golang:latest").                                   // Use Go's official image for a reliable build environment
		WithDirectory("/src", source.Directory("todolist-app")). // Mount the todolist-app directory
		WithWorkdir("/src/cmd/server").                          // Navigate to the correct working directory
		WithExec([]string{"go", "build", "-o", "server", "."})   // Build the Go binary named 'server'
}

// Calls a Build Function by Reference - Uses It's Cached Output to Publish Image
func (m *PocDaggerVsTraditionalCi) PublishImage(ctx context.Context, source *dagger.Directory) (string, error) {
	return m.Build(ctx, source).
		WithEntrypoint([]string{"/src/todolist-app"}).
		Publish(ctx, "ttl.sh/todolist-app-one2n")
}
