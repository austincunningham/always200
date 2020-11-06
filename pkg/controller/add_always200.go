package controller

import (
	"always200/pkg/controller/always200"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, always200.Add)
}
