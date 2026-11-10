package cla

import (
	"os"
	"strconv"

	m "github.com/vault-thirteen/SWVD/models"
	"github.com/vault-thirteen/auxie/errors"
)

const (
	ArgsCount = 3
)

const (
	UsageHint = `Usage: list.exe <SettingsFilePath> <AuthorName> <PagesCount>`
)

type CommandLineArgs struct {
	SettingsFilePath string
	AuthorName       string
	PagesCount       int
}

func NewCommandLineArgs() (args *CommandLineArgs, err error) {
	if len(os.Args) < ArgsCount+1 {
		return nil, errors.ErrorFromText(m.Err_NotEnoughArguments)
	}

	var pagesCount int
	pagesCount, err = strconv.Atoi(os.Args[3])
	if err != nil {
		return nil, err
	}

	args = &CommandLineArgs{
		SettingsFilePath: os.Args[1],
		AuthorName:       os.Args[2],
		PagesCount:       pagesCount,
	}

	return args, nil
}
