package main

import (
  "fmt"
  "os"
  "devlab/bin/context"
  "devlab/bin/create-docker-compose"
)

func main() {     
  switch os.Args[1] {
  case "context":
    switch os.Args[2] {    
    case "create":
      if err := Context.Create(os.Args[3]); err != nil {
        fmt.Println(err)
        os.Exit(1) 
      }  
      break 
    case "install":
      if err := Context.Install(os.Args[3]); err != nil {
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