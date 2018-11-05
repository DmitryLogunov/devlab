package createDockerCompose

import (
	"devlab/lib/files"
  "devlab/lib/docker-compose-file-builder"
)

func Call(contextName string) {
	config, _ := files.ReadMainConfig()
	dockerComposeData := DockerComposeFileBuilder.CreateDockerComposeObjectExample()
	contextDir := config["contexts-path"] + "/" + contextName
	DockerComposeFileBuilder.Create("./" + contextDir + "/docker-compose.application.yml", dockerComposeData)
}