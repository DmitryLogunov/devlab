package errors

import (
	"errors"
)

var (
	ErrNotDefinedConfigurationPath     = errors.New("context: configuration path is not defined")
	ErrCouldntReadConfiguration        = errors.New("context: couldn't read configuration")
	ErrCouldntParseConfiguration       = errors.New("context: couldn't parse configurations")
	ErrCouldntReadConfig               = errors.New("context: couldn't read config")
	ErrNotDefinedContextName           = errors.New("context: context name is not defined")
	ErrNotDefinedContextsPath          = errors.New("context: contexts path is not defined")
	ErrNotDefinedTemplatesContextsPath = errors.New("context: templtes contexts path is not defined")
	ErrContextIsNotCreated             = errors.New("context: contexts is not created")
)
