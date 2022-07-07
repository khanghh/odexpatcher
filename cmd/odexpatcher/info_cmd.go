package main

import (
	"errors"
	"io/ioutil"
	"log"
	v124 "odexpatcher/inernal/oat/v124"
)

type InfoCmd struct {
}

func (cmd *InfoCmd) Execute(args []string) error {
	fileData, err := ioutil.ReadFile(args[0])
	if err != nil {
		log.Fatalln("Could not read input file:", err)
	}

	dex, err := v124.ParseDex(fileData)
	if err == nil {
		dex.PrintInfo()
		return nil
	}

	odex, err := v124.ParseOdex(fileData)
	if err == nil {
		odex.PrintInfo()
		return nil
	}

	vdex, err := v124.ParseVdex(fileData)
	if err == nil {
		vdex.PrintInfo()
		return nil
	}

	return errors.New("Could not parse input file.")
}
