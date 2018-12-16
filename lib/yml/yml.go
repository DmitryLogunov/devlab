package yml

import (
  "github.com/gopkg.in/yaml"
  "devlab/lib/errors"
)

func ParseOneLevelYAML(data string) (parsedData map[string]string, err error) {  
  err = yaml.Unmarshal([]byte(data), &parsedData)
  if( errors.CheckAndReturnIfError(err) ) { return make(map[string]string), err }

  return
}

func ParseTwoLevelYAML(data string) (parsedData map[string]map[string]string, err error) {  
  err = yaml.Unmarshal([]byte(data), &parsedData)
  if( errors.CheckAndReturnIfError(err) ) { return make(map[string]map[string]string), err }

  return
}

func ParseThreeLevelYAML(data string) (parsedData map[string]map[string]map[string]string, err error) {  
  err = yaml.Unmarshal([]byte(data), &parsedData)
  if( errors.CheckAndReturnIfError(err) ) { return make(map[string]map[string]map[string]string), err }

  return
}

