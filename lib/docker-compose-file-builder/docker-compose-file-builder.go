package DockerComposeFileBuilder

import (
  "devlab/lib/files"
)

type Service struct {
  image string
  env_files []string
  volumes []string
  ports []string
  restart string    
}

type DockerComposeFile struct {
  version string
  services map[string]Service
  networks map[string]map[string]map[string]string
}

func CreateDockerComposeObjectExample() *DockerComposeFile {
	dockerComposeExample := new(DockerComposeFile)
  dockerComposeExample.version = "2"
  dockerComposeExample.services = make(map[string]Service)
  dockerComposeExample.services["dlp-service-config"] = Service{
    image: "${IMAGES_PREFIX}dlp-service-config",
    env_files: []string{"${BUILD_DIR}/dlp-service-config/.env"},
    volumes: []string{"${DEVENV_ROOT_DIR}/dlp-service-config:/usr/src/app"},
    ports: []string{"4004"},
    restart: "always"}

	dockerComposeExample.networks = map[string]map[string]map[string]string{ "default": {"external": { "name": "bedrock" } } }
	
	return dockerComposeExample
}


func Create(dockerComposeFilePath string, dockerComposeData *DockerComposeFile) {     
  files.WriteAppendFileWithIndent(dockerComposeFilePath, "version: " + dockerComposeData.version, 0)
  files.WriteAppendFileWithIndent(dockerComposeFilePath, "services: ", 0)

  for serviceName, serviceData := range dockerComposeData.services {
    files.WriteAppendFileWithIndent(dockerComposeFilePath, serviceName + ":", 2)
    
    files.WriteAppendFileWithIndent(dockerComposeFilePath, "image: " + serviceData.image, 4)
    
    if len(serviceData.env_files) > 0 {
      files.WriteAppendFileWithIndent(dockerComposeFilePath, "env_files: ", 4)
      for _, envFile := range serviceData.env_files {
        files.WriteAppendFileWithIndent(dockerComposeFilePath, "- " + envFile, 6)
      }
    }

    if len(serviceData.volumes) > 0 {
      files.WriteAppendFileWithIndent(dockerComposeFilePath, "volumes: ", 4)
      for _, volume := range serviceData.volumes {
        files.WriteAppendFileWithIndent(dockerComposeFilePath, "- " + volume, 6)
      }
    }

    if len(serviceData.ports) > 0 {
      files.WriteAppendFileWithIndent(dockerComposeFilePath, "ports: ", 4)
      for _, port := range serviceData.ports {
        files.WriteAppendFileWithIndent(dockerComposeFilePath, "- " + port, 6)
      }
    }        
     
    if serviceData.restart != "" {
      files.WriteAppendFileWithIndent(dockerComposeFilePath, "restart: " + serviceData.restart, 4)
    } 
  }

  files.WriteAppendFileWithIndent(dockerComposeFilePath, "networks: ", 0)
  files.WriteAppendFileWithIndent(dockerComposeFilePath, "default: ", 2)
  files.WriteAppendFileWithIndent(dockerComposeFilePath, "external: ", 4)
  files.WriteAppendFileWithIndent(dockerComposeFilePath, "name: " + dockerComposeData.networks["default"]["external"]["name"] , 6)
}