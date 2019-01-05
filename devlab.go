package main

import (
	"devlab/bin/context"
	createDockerCompose "devlab/bin/create-docker-compose"
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
				fmt.Println(err)
				os.Exit(1)
			}
		}
		break
	case "create-docker-compose":
		createDockerCompose.Call(os.Args[2])
		break
	}
}
