//go:build mage

package main

import (
	"errors"
	"fmt"
	"gamma/app/datastore/events"
	"gamma/build/db"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type (
	Run      mg.Namespace
	Generate mg.Namespace
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

func (Generate) EventModels() error {
	eventdb := events.EventDB()
	if eventdb == nil {
		return errors.New("Couldn't connect to database")
	}
	eventdb.Close()
	err := sh.RunV("xo", "schema", fmt.Sprintf("pgsql://%s:%s@localhost:%d/%s?sslmode=disable", POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_PORT, POSTGRES_DB), "-o", "./app/datastore/event/models")
	return err
}

func (Generate) EventModelQueries() error {
	queries := db.FindXOQueries("event")

	for _, query := range queries {
		err := sh.RunV(
			"xo",
			"query",
			fmt.Sprintf("pg://%s:%s@localhost:%d/%s?sslmode=disable", POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_PORT, POSTGRES_DB),
			"-o",
			"./app/datastore/event/models",
			"-M",
			"-B",
			"-2",
			"-T",
			query.TypeName,
			"-Q",
			fmt.Sprintf("%s", query.Query),
		)
		if err != nil {
			return err
		}
	}
	return nil
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
