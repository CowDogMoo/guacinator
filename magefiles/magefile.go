//go:build mage

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/bitfield/script"
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

// InstallDeps Installs go dependencies
func InstallDeps() error {
	fmt.Println("Installing dependencies.")

	if err := mageutils.Tidy(); err != nil {
		return fmt.Errorf("failed to install dependencies: %v", err)
	}

	if err := lint.InstallGoPCDeps(); err != nil {
		return fmt.Errorf("failed to install pre-commit dependencies: %v", err)
	}

	if err := mageutils.InstallVSCodeModules(); err != nil {
		return fmt.Errorf("failed to install vscode-go modules: %v", err)
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

	if err := sys.Cd(filepath.Join(repoRoot, "magefiles")); err != nil {
		return fmt.Errorf("failed to cd to magefiles directory: %v", err)
	}

	repo := docs.Repo{
		Owner: "cowdogmoo",
		Name:  "guacinator",
	}

	templatePath := filepath.Join(repoRoot, "magefiles", "tmpl", "README.md.tmpl")

	if err := docs.CreatePackageDocs(fs, repo, templatePath); err != nil {
		return fmt.Errorf("failed to create package docs: %v", err)
	}

	fmt.Println("Package docs created.")

	return nil
}

// RunPreCommit runs all pre-commit hooks locally
func RunPreCommit() error {
	mg.Deps(InstallDeps)

	fmt.Println("Updating pre-commit hooks.")
	if err := lint.UpdatePCHooks(); err != nil {
		return err
	}

	fmt.Println("Clearing the pre-commit cache to ensure we have a fresh start.")
	if err := lint.ClearPCCache(); err != nil {
		return err
	}

	fmt.Println("Running all pre-commit hooks locally.")
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

// UpdateMirror updates pkg.go.goutils with the release associated with the input tag
func UpdateMirror(tag string) error {
	var err error
	fmt.Printf("Updating pkg.go.goutils with the new tag %s.", tag)

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

	templatePath := "magefiles/tmpl/README.md.tmpl"

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
