//go:build mage
// +build mage

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/bitfield/script"
	"github.com/fatih/color"
	"github.com/l50/goutils/v2/dev/lint"
	mageutils "github.com/l50/goutils/v2/dev/mage"
	"github.com/l50/goutils/v2/docs"
	"github.com/l50/goutils/v2/git"
	"github.com/l50/goutils/v2/sys"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/spf13/afero"
)

var (
	deploymentDir string
)

func init() {
	os.Setenv("GO111MODULE", "on")
}

// InstallDeps installs the Go dependencies necessary for developing
// on the project.
//
// Example usage:
//
// ```go
// mage installdeps
// ```
//
// **Returns:**
//
// error: An error if any issue occurs while trying to
// install the dependencies.
func InstallDeps() error {
	fmt.Println(color.YellowString("Installing dependencies."))
	if err := lint.InstallGoPCDeps(); err != nil {
		return fmt.Errorf("failed to install pre-commit dependencies: %v", err)
	}

	if err := mageutils.InstallVSCodeModules(); err != nil {
		return fmt.Errorf(color.RedString(
			"failed to install vscode-go modules: %v", err))
	}

	return nil
}

// FindExportedFuncsWithoutTests finds exported functions without tests
func FindExportedFuncsWithoutTests(pkg string) ([]string, error) {
	funcs, err := mageutils.FindExportedFuncsWithoutTests(os.Args[1])

	if err != nil {
		log.Fatalf("failed to find exported functions without tests: %v", err)
	}

	for _, funcName := range funcs {
		fmt.Println(funcName)
	}

	return funcs, nil

}

// GeneratePackageDocs generates package documentation
// for packages in the current directory and its subdirectories.
func GeneratePackageDocs() error {
	fs := afero.NewOsFs()

	repoRoot, err := git.RepoRoot()
	if err != nil {
		return fmt.Errorf("failed to get repo root: %v", err)
	}

	if err := sys.Cd(repoRoot); err != nil {
		return fmt.Errorf("failed to change directory to repo root: %v", err)
	}

	repo := docs.Repo{
		Owner: "cowdogmoo",
		Name:  "guacinator",
	}

	templatePath := filepath.Join(repoRoot, "templates", "README.md.tmpl")
	// Set the packages to exclude (optional)
	excludedPkgs := []string{"main"}
	if err := docs.CreatePackageDocs(fs, repo, templatePath, excludedPkgs...); err != nil {
		return fmt.Errorf("failed to create package docs: %v", err)
	}

	return nil
}

// RunPreCommit updates, clears, and executes all pre-commit hooks
// locally. The function follows a three-step process:
//
// First, it updates the pre-commit hooks.
// Next, it clears the pre-commit cache to ensure a clean environment.
// Lastly, it executes all pre-commit hooks locally.
//
// Example usage:
//
// ```go
// mage runprecommit
// ```
//
// **Returns:**
//
// error: An error if any issue occurs at any of the three stages
// of the process.
func RunPreCommit() error {
	if !sys.CmdExists("pre-commit") {
		return fmt.Errorf("pre-commit is not installed, please install it " +
			"with the following command: `python3 -m pip install pre-commit`")
	}

	fmt.Println(color.YellowString("Updating pre-commit hooks."))
	if err := lint.UpdatePCHooks(); err != nil {
		return err
	}

	fmt.Println(color.YellowString("Clearing the pre-commit cache to ensure we have a fresh start."))
	if err := lint.ClearPCCache(); err != nil {
		return err
	}

	fmt.Println(color.YellowString("Running all pre-commit hooks locally."))
	if err := lint.RunPCHooks(); err != nil {
		return err
	}

	return nil
}

// RunTests runs all of the unit tests
func RunTests() error {
	mg.Deps(InstallDeps)

	fmt.Println("Running unit tests.")
	if err := sh.RunV(filepath.Join(".hooks", "go-unit-tests.sh"), "all"); err != nil {
		return fmt.Errorf("failed to run unit tests: %v", err)
	}

	return nil
}

// UpdateMirror updates pkg.go.dev with the release associated with the
// input tag
//
// Example usage:
//
// ```go
// mage updatemirror v2.0.1
// ```
//
// **Parameters:**
//
// tag: the tag to update pkg.go.dev with
//
// **Returns:**
//
// error: An error if any issue occurs while updating pkg.go.dev
func UpdateMirror(tag string) error {
	var err error
	fmt.Printf("Updating pkg.go.dev with the new tag %s.", tag)

	err = sh.RunV("curl", "--silent", fmt.Sprintf(
		"https://sum.golang.org/lookup/github.com/cowdogmoo/guacinator/@%s",
		tag))
	if err != nil {
		return fmt.Errorf("failed to update proxy.golang.org: %w", err)
	}

	err = sh.RunV("curl", "--silent", fmt.Sprintf(
		"https://proxy.golang.org/github.com/cowdogmoo/guacinator/@v/%s.info",
		tag))
	if err != nil {
		return fmt.Errorf("failed to update pkg.go.dev: %w", err)
	}

	return nil
}

// UpdateDocs updates the package documentation
// for packages in the current directory and its subdirectories.
func UpdateDocs() error {
	repo := docs.Repo{
		Owner: "cowdogmoo",
		Name:  "guacinator",
	}

	fs := afero.NewOsFs()

	templatePath := "templates/README.md.tmpl"

	fmt.Println("Updating docs.")
	if err := docs.CreatePackageDocs(fs, repo, templatePath); err != nil {
		return fmt.Errorf("failed to update docs: %v", err)
	}

	return nil
}

// OnDemand is used to spin up an on-demand deployment
// using an input Docker image.
// Example:
// OnDemand(ubuntu)
// func OnDemand(image string) error {
func OnDemand() error {
	if err := sys.Cd(fmt.Sprintf(
		filepath.Join(fmt.Sprintf("%s", deploymentDir), "od"))); err != nil {
		return err
	}

	cmds := []string{
		"kubectl apply -f ubuntu-deployment.yaml",
	}

	for _, cmd := range cmds {
		if _, err := script.Exec(cmd).Stdout(); err != nil {
			return err
		}
	}

	fmt.Println("Run the following command to get the IP address to access your new system:")
	fmt.Println("kubectl get service -o wide | grep vnc | awk '{print $3}'")

	return nil
}
