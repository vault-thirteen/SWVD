package main

import (
	"os"
	"strconv"

	m "github.com/vault-thirteen/SWVD/models"
	"github.com/vault-thirteen/auxie/errors"
)

const (
	Err_NotEnoughArguments = "not enough arguments"
)

const (
	ArgsCount = 4
)

type CommandLineArgs struct {
	authorName      string
	cookiesFilePath string
	dataFolder      string
	quality         m.QualityInt
}

func NewCommandLineArgs() (args *CommandLineArgs, err error) {
	if len(os.Args) < ArgsCount+1 {
		return nil, errors.ErrorFromText(Err_NotEnoughArguments)
	}

	var q int
	q, err = strconv.Atoi(os.Args[4])
	if err != nil {
		return nil, err
	}

	args = &CommandLineArgs{
		authorName:      os.Args[1],
		cookiesFilePath: os.Args[2],
		dataFolder:      os.Args[3],
		quality:         m.QualityInt(q),
	}

	return args, nil
}
