package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-steputils/stepconf"
)

// Configs ...
type Configs struct {
	Token         string `env:"token,required"`
	App           string `env:"app,required"`
	TestFramework string `env:"framework,opt[appium,calabash,espresso,xcuitest,uitest]"`
	Devices       string `env:"devices,required"`
	Series        string `env:"series,required"`
	Locale        string `env:"locale,required"`
	AppPath       string `env:"app_path,file"`
	DSYMDir       string `env:"dsym_dir"`
	TestDir       string `env:"test_dir,dir"`
}

func uploadTestCommand(apiToken, framework, app, devices, series, local, appPath, dsymDir, testDir string) *command.Model {
	args := []string{"test", "run", string(framework),
		"--token", apiToken,
		"--app", app,
		"--devices", devices,
		"--test-series", series,
		"--locale", local,
		"--async",
		"--app-path", appPath,
	}
	if dsymDir != "" {
		args = append(args, "--dsym-dir", dsymDir)
	}
	if framework == "calabash" {
		args = append(args, "--project-dir", testDir)
	} else {
		args = append(args, "--build-dir", testDir)
	}
	return command.New("appcenter", args...)
}

func mainE() error {
	var cfg Configs
	if err := stepconf.Parse(&cfg); err != nil {
		return fmt.Errorf("Couldn't create config: %s", err)
	}
	stepconf.Print(cfg)

	if _, err := exec.LookPath("appcenter"); err != nil {
		cmd := command.New("npm", "install", "-g", "appcenter-cli")

		log.Infof("\nInstalling appcenter-cli")
		log.Donef("$ %s", cmd.PrintableCommandArgs())

		if out, err := cmd.RunAndReturnTrimmedCombinedOutput(); err != nil {
			return fmt.Errorf("Failed to install appcenter-cli: %s", out)
		}
	}

	cmd := uploadTestCommand(cfg.Token, cfg.TestFramework, cfg.App, cfg.Devices, cfg.Series, cfg.Locale, cfg.AppPath, cfg.DSYMDir, cfg.TestDir).SetStdout(os.Stdout).SetStderr(os.Stderr)

	log.Infof("\nUploading and scheduling tests")
	log.Donef("$ %s", cmd.PrintableCommandArgs())

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Upload failed, error: %s", err)
	}
	return nil
}

func main() {
	if err := mainE(); err != nil {
		log.Errorf("%s", err)
		os.Exit(1)
	}
}
