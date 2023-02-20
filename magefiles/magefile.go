//go:build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	goutils "github.com/l50/goutils"
	"github.com/magefile/mage/mg"

	"github.com/bitfield/script"
)

const (
	debug = false
)

var (
	deploymentDir string
)

func init() {
	os.Setenv("GO111MODULE", "on")
}

// InstallDeps Installs go dependencies
func InstallDeps() error {
	fmt.Println(color.YellowString("Installing dependencies."))

	if err := goutils.Tidy(); err != nil {
		return fmt.Errorf(color.RedString(
			"failed to install dependencies: %v", err))
	}

	if err := goutils.InstallGoPCDeps(); err != nil {
		return fmt.Errorf(color.RedString(
			"failed to install pre-commit dependencies: %v", err))
	}

	if err := goutils.InstallVSCodeModules(); err != nil {
		return fmt.Errorf(color.RedString(
			"failed to install vscode-go modules: %v", err))
	}

	return nil
}

// InstallPreCommitHooks Installs pre-commit hooks locally
func InstallPreCommitHooks() error {
	mg.Deps(InstallDeps)

	fmt.Println(color.YellowString("Installing pre-commit hooks."))
	if err := goutils.InstallPCHooks(); err != nil {
		return err
	}

	return nil
}

// RunPreCommit runs all pre-commit hooks locally
func RunPreCommit() error {
	mg.Deps(InstallDeps)

	fmt.Println(color.YellowString("Updating pre-commit hooks."))
	if err := goutils.UpdatePCHooks(); err != nil {
		return err
	}

	fmt.Println(color.YellowString(
		"Clearing the pre-commit cache to ensure we have a fresh start."))
	if err := goutils.ClearPCCache(); err != nil {
		return err
	}

	fmt.Println(color.YellowString("Running all pre-commit hooks locally."))
	if err := goutils.RunPCHooks(); err != nil {
		return err
	}

	return nil
}

// OnDemand is used to spin up an on-demand deployment
// using an input Docker image.
// Example:
// OnDemand(ubuntu)
// func OnDemand(image string) error {
func OnDemand() error {
	goutils.Cd(fmt.Sprintf(
		filepath.Join(fmt.Sprintf("%s", deploymentDir), "od")))
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
