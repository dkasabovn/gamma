//go:build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type (
	Run      mg.Namespace
	Generate mg.Namespace
	Test     mg.Namespace
)

const (
	POSTGRES_PASSWORD = "nhPldb98Rt"
	POSTGRES_PORT     = 5432
	POSTGRES_USER     = "postgres"
	POSTGRES_DB       = "postgres"
)

func (Run) EventDB() error {
	// build the database
	err := sh.RunV("docker", "build", "-t", "gamma/eventdb", "-f", "./db/event/Dockerfile", "./db/event/")
	if err != nil {
		return err
	}

	sh.RunV("docker", "container", "stop", "gamma_db")
	sh.RunV("docker", "container", "rm", "gamma_db")

	err = sh.RunV("docker", "run", "-p", "5432:5432", "-d", "--name=gamma_db", "gamma/eventdb")

	return err
}

func (Test) All() error {
	return sh.RunV("ginkgo", "./app/...")
}

func (Generate) PrivatePublicKeys() error {
	err := sh.RunV(
		"openssl",
		"ecparam",
		"-name",
		"prime256v1",
		"-genkey",
		"-noout",
		"-out",
		"private-key.pem",
	)

	if err != nil {
		return err
	}

	return sh.RunV(
		"openssl",
		"ec",
		"-in",
		"private-key.pem",
		"-pubout",
		"-out",
		"public-key.pem",
	)
}
