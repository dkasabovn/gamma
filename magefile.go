//go:build mage

package main

import (
	"errors"
	"gamma/app/datastore/event"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type (
	Run      mg.Namespace
	Generate mg.Namespace
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
	eventdb := event.EventDB()
	if eventdb == nil {
		return errors.New("Couldn't connect to database")
	}

	err := sh.RunV("xo", "schema", "pgsql://docker:nhPldb98Rt@localhost:5432/eventsvcdb?sslmode=disable", "-o", "./app/datastore/event/models")
	eventdb.Close()
	return err
}
