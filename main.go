package main

import (
	"log"
	"os"
	"path/filepath"
	"syscall"

	"github.com/AlecAivazis/survey/v2"
	"github.com/kevinburke/ssh_config"
)

func main() {
	f, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "config"))
	if err != nil {
		log.Fatalf("open config file error %s\n", err)
	}

	cfg, err := ssh_config.Decode(f)
	if err != nil {
		log.Fatalf("parse config file error %s\n", err)
	}

	var options []string
	for _, host := range cfg.Hosts {
		patterns := host.Patterns
		options = append(options, patterns[0].String())
	}

	var qs = []*survey.Question{
		{
			Name: "host",
			Prompt: &survey.Select{
				Message: "Choose a host:",
				Options: options,
				Default: "red",
			},
		},
	}

	answers := struct {
		Host string `survey:"host"` // or you can tag fields to match a specific name
	}{}

	err = survey.Ask(qs, &answers)
	if err != nil {
		log.Fatalf("ask error %s\n", err)
	}

	ssh := "/usr/bin/ssh"
	argv := []string{ssh, answers.Host}
	envv := os.Environ()
	err = syscall.Exec(ssh, argv, envv)
	if err != nil {
		log.Fatalf("run command error: %s", err)
	}
}
