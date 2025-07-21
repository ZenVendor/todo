package main

import "errors"

var (
	ErrNoArguments          = errors.New("Parser: No arguments provided")
	ErrInvalidVerb          = errors.New("Parser: Invalid verb")
	ErrInvalidArgument      = errors.New("Parser: Invalid argument")
	ErrVerbRequiresArgument = errors.New("Parser: An argument is required for verb")
	ErrVerbRequiresValue    = errors.New("Parser: A value is required for verb")
	ErrTooManyArguments     = errors.New("Parser: Too many arguments for the verb")
	ErrInvalidVerbArgument  = errors.New("Parser: Verb does not support argument")
	ErrInvalidDefaultArg    = errors.New("Parser: Verb does not support the default argument")
	ErrInvalidDefaultKwarg  = errors.New("Parser: Verb does not support the default keyword argument")
	ErrInvalidDate          = errors.New("Validator: Invalid date value")
	ErrStringLength         = errors.New("Validator: String value is too long")
    ErrDBVersion            = errors.New("DB Check: Database requires an update!")
    ErrNoConfig             = errors.New("Configuration file not found! Please create the file or reinstall.")
)
