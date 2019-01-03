package main

import (
  "os"
  "devlab/bin/context"
  "devlab/bin/create-docker-compose"
)

func main() {     
  switch os.Args[1] {
  case "context":
    switch os.Args[2] {
//    case "create":
//      Context.Create(os.Args[3])         
    case "set":
      Context.Set(os.Args[3])   
    }
    break 
  case "create-docker-compose":
    createDockerCompose.Call(os.Args[2])
    break
  }
}