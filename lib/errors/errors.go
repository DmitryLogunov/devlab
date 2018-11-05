package errors

import (
	"fmt"
	"os"
  )

func CheckAndExitIfError(err  error) {
  if err != nil {
		fmt.Println(err) 
		os.Exit(1) 
	}
}

func CheckAndReturnIfError(err  error) bool {
  if err != nil {
		fmt.Println(err) 
		return true
  }
  
  return false
}
