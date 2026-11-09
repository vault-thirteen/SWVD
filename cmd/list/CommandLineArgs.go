package main

import (
	"os"
	"strconv"

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
	pagesCount      int
	cookiesFilePath string
	dataFolder      string
}

func NewCommandLineArgs() (args *CommandLineArgs, err error) {
	if len(os.Args) < ArgsCount+1 {
		return nil, errors.ErrorFromText(Err_NotEnoughArguments)
	}

	var pagesCount int
	pagesCount, err = strconv.Atoi(os.Args[2])
	if err != nil {
		return nil, err
	}

	args = &CommandLineArgs{
		authorName:      os.Args[1],
		pagesCount:      pagesCount,
		cookiesFilePath: os.Args[3],
		dataFolder:      os.Args[4],
	}

	return args, nil
}
