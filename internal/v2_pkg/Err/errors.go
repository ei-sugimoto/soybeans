package Err

import "errors"

var ManyArgs = errors.New("too many arguments")
var CreateContainerDir = errors.New("failed to create container directory")

var SaveContainerState = errors.New("failed to save container state")
