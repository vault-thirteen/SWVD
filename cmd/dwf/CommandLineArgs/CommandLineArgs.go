package cla

import (
	"os"

	m "github.com/vault-thirteen/SWVD/models"
	"github.com/vault-thirteen/auxie/boolean"
	"github.com/vault-thirteen/auxie/errors"
)

const (
	ArgsCount = 2
)

const (
	UsageHint = `Usage: dwf.exe <SettingsFilePath> <DeleteReportAfterSuccess>`
)

type CommandLineArgs struct {
	SettingsFilePath         string
	DeleteReportAfterSuccess bool
}

func NewCommandLineArgs() (args *CommandLineArgs, err error) {
	if len(os.Args) < ArgsCount+1 {
		return nil, errors.ErrorFromText(m.Err_NotEnoughArguments)
	}

	args = &CommandLineArgs{
		SettingsFilePath:         os.Args[1],
		DeleteReportAfterSuccess: false,
	}

	args.DeleteReportAfterSuccess, err = boolean.FromString(os.Args[2])
	if err != nil {
		return nil, err
	}

	return args, nil
}
