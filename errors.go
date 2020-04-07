package main

import (
	"fmt"
	"os"
)

const (
	ERROR_FAILED_TO_INIT_SDL                  int = 4
	ERROR_FAILED_TO_CREATE_WINDOW             int = 5
	ERROR_FAILED_TO_CREATE_RENDERER           int = 6
	ERROR_FAILED_TO_LOAD_IMAGE                int = 7
	ERROR_FAILED_TO_CREATE_TEXTURE_FROM_IMAGE int = 8
)

func HandleError(message string, err error) {
	fmt.Fprintf(os.Stderr, "%s %s\n", message, err)
}

/*func HandleFatalError(err error) {
	if err != nil {
		//fmt.Fprintf(os.Stderr, err.Error())
		panic(err)
	}
}*/
