package main

import (
	"devlab/cmd/context"
	contextErrors "devlab/cmd/context/common/errors"
	createDockerCompose "devlab/cmd/create-docker-compose"
	"fmt"
	"os"
)

func main() {
	switch os.Args[1] {
	case "context":
		switch os.Args[2] {
		case "create":
			if err := context.Create(os.Args[3]); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			break
		case "install":
			if err := context.Install(os.Args[3]); err != nil {
				if err != contextErrors.ErrContextIsNotCreated {
					fmt.Println(err)
				}
				os.Exit(1)
			}
		}
		break
}
