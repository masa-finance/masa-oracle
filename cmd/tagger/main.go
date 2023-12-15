/*
Package main - tagger

The tagger is a utility tool designed to automate the process of versioning, building, and tagging Go projects. It streamlines the workflow by encapsulating several steps into a single command, reducing the potential for human error and increasing efficiency.

The tagger performs the following operations:

1. Updates the project: It navigates to the project directory and pulls the latest changes from the specified branch. It then tidies up the dependencies using 'go mod tidy'.

2. Pushes changes: It stages and commits changes with a user-provided commit message, then pushes the commit to the specified branch.

3. Tags the project: It tags the current commit with a user-provided version tag and pushes the tag to the remote repository.

4. Builds the project: It builds the project for multiple platforms (Linux, Darwin, and Windows) and outputs the binaries to the 'bin' directory in the project directory.

The tagger is invoked from the command line with four flags: 'projectDir' (the project directory), 'version' (the version to tag), 'branch' (the branch to use), and 'commitMessage' (the commit message). All flags must be provided.

Example usage:
go run tag_project.go -projectDir=/path/to/project -version=1.0.0 -branch=main -commitMessage="Update go.mod"

This utility tool is intended for use by developers who need to regularly version, build, and tag the project. It assumes that the user has a working Go environment and is familiar with basic Git operations.
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

func updateProject(projectDir, branch string) error {
	// Change to the project directory
	if err := os.Chdir(projectDir); err != nil {
		return err
	}

	// Execute git pull
	if err := runCommand("git", "pull", "origin", branch); err != nil {
		return err
	}
	// Execute go mod tidy
	if err := runCommand("go", "mod", "tidy"); err != nil {
		return fmt.Errorf("failed to tidy dependencies: %w", err)
	}
	return nil
}

func pushChanges(projectDir, commitMessage, branch string) (err error) {
	// Change to the project directory
	if err := os.Chdir(projectDir); err != nil {
		return err
	}

	// Execute git add .
	if err := runCommand("git", "add", "."); err != nil {
		return err
	}

	// Execute git commit -m "update go.mod"
	if err := runCommand("git", "commit", "-m", commitMessage); err != nil {
		return err
	}

	// Execute git push origin main
	if err := runCommand("git", "push", "origin", branch); err != nil {
		return err
	}
	return nil
}

func tagProject(projectDir, version string) error {
	// Change to the project directory
	if err := os.Chdir(projectDir); err != nil {
		return err
	}
	if err := runCommand("git", "tag", "-a", version, "-m", fmt.Sprintf("Version %s", version)); err != nil {
		return err
	}

	// Execute git push --tags
	if err := runCommand("git", "push", "--tags"); err != nil {
		return err
	}
	return nil
}

func buildProject(projectDir string) error {
	// Define the target platforms
	platforms := []struct {
		goos, goarch, ext string
	}{
		{"linux", "amd64", ""},
		{"darwin", "amd64", ""},
		{"windows", "amd64", ".exe"},
	}

	// Build the project for each platform
	for _, platform := range platforms {
		// Set the target platform
		err := os.Setenv("GOOS", platform.goos)
		if err != nil {
			return err
		}
		err = os.Setenv("GOARCH", platform.goarch)
		if err != nil {
			return err
		}

		// Define the output path
		outputPath := filepath.Join(projectDir, "bin", "masa-oracle-"+platform.goos+platform.ext)

		// Execute go build
		if err := runCommand("go", "build", "-o", outputPath, "."); err != nil {
			return err
		}
	}

	return nil
}

func runCommand(command string, args ...string) error {
	fmt.Printf("Running command: %s %v\n", command, args)
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func TagAndBuild(projectDir, version, branch, commitMessage string) (err error) {
	if projectDir == "" {
		usr, err := user.Current()
		if err != nil {
			log.Fatal("could not find user.home directory")
		}
		projectDir = filepath.Join(usr.HomeDir, "github", "masa-finance")
	}
	err = updateProject(projectDir, branch)
	if err != nil {
		return err
	}
	err = pushChanges(projectDir, commitMessage, branch)
	if err != nil {
		return err
	}
	err = tagProject(projectDir, version)
	if err != nil {
		return err
	}
	err = buildProject(projectDir)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	// Define flags
	projectDir := flag.String("projectDir", "", "The project directory")
	version := flag.String("version", "", "The version to tag")
	branch := flag.String("branch", "", "The branch to use")
	commitMessage := flag.String("commitMessage", "", "The commit message")

	// Parse flags
	flag.Parse()

	// Check that all flags have been provided
	if *projectDir == "" || *version == "" || *branch == "" || *commitMessage == "" {
		fmt.Println("Error: All flags must be provided")
		fmt.Println("Usage: tag_project -projectDir=<projectDir> -version=<version> -branch=<branch> -commitMessage=<commitMessage>")
		os.Exit(1)
	}

	// Call TagAndBuild with the parsed arguments
	err := TagAndBuild(*projectDir, *version, *branch, *commitMessage)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Update and tagging completed successfully.")
	}
}
