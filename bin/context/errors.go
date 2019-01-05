package context

import (
	"errors"
)

var (
	ErrNotDefinedConfigurationPath     = errors.New("ERROR: context => configuration path is not defined")
	ErrCouldntReadConfiguration        = errors.New("ERROR: context => couldn't read configuration")
	ErrCouldntParseConfiguration       = errors.New("ERROR: context => couldn't parse configurations")
	ErrCouldntReadConfig               = errors.New("ERROR: context => couldn't read config")
	ErrNotDefinedContextName           = errors.New("ERROR: context => context name is not defined")
	ErrNotDefinedContextsPath          = errors.New("ERROR: context => contexts path is not defined")
	ErrNotDefinedTemplatesContextsPath = errors.New("ERROR: context => templtes contexts path is not defined")
	ErrContextIsNotCreated             = errors.New("ERROR: context => contexts is not created")
)
