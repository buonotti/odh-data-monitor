package errors

import (
	"os"

	"github.com/joomcode/errorx"

	"github.com/buonotti/apisense/log"
	charmlog "github.com/charmbracelet/log"
)

// fatalTrait is a trait that can be added to any error type to make it fatal and make the program exit with code 1
var fatalTrait = errorx.RegisterTrait("fatal")

// CheckErr is a helper function that handles a given error. If the error is
// not nil it logs the error and exits with code 1. This function should only be
// used in the top level scope of the program or when an error cant be returned
func CheckErr(err error) {
	if _, ok := err.(*errorx.Error); ok {
		handleErrorxError(err.(*errorx.Error))
	} else {
		handleError(err)
	}
}

func handleErrorxError(err *errorx.Error) {
	if err != nil {
		if err.HasTrait(fatalTrait) {
			if charmlog.GetLevel() == charmlog.DebugLevel {
				log.DefaultLogger().Errorf("%+v", err)
			} else {
				log.DefaultLogger().Error(err.Error())
			}
			os.Exit(1)
		} else {
			if charmlog.GetLevel() == charmlog.DebugLevel {
				log.DefaultLogger().Warnf("%+v", err)
			} else {
				log.DefaultLogger().Warn(err.Error())
			}
		}
	}
}

func handleError(err error) {
	if err != nil {
		log.DefaultLogger().Fatal(err.Error())
		os.Exit(1)
	}
}
