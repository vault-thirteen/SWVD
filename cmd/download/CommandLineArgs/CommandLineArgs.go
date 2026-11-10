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
	UsageHint = `Usage: download.exe <SettingsFilePath> <AuthorName> <Quality>`
)

type CommandLineArgs struct {
	SettingsFilePath string
	AuthorName       string
	Quality          m.QualityInt
}

func NewCommandLineArgs() (args *CommandLineArgs, err error) {
	if len(os.Args) < ArgsCount+1 {
		return nil, errors.ErrorFromText(m.Err_NotEnoughArguments)
	}

	var q int
	q, err = strconv.Atoi(os.Args[3])
	if err != nil {
		return nil, err
	}

	args = &CommandLineArgs{
		SettingsFilePath: os.Args[1],
		AuthorName:       os.Args[2],
		Quality:          m.QualityInt(q),
	}

	return args, nil
}
