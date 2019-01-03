package Context

type SettingsParametersDescription struct {
  Maintainer string
  Created string
}

type SettingsParametersGit struct {
  BaseBranch string
}

type SettingsParametersDocker struct {
  ImagesPrefix string
  Network string
}

type SettingsParametersPaths struct {
  Templates string
}

type SettingsParametersBuilding struct {
  Template string
}

type SettingsParametersDeploying struct {
  Template string
}

type SettingsParametersSystemServices struct {
  Template string
}

type SettingsParametersApplicatonServices struct {
  Template string
}

type SettingsParameters struct {
  Name string
  Description SettingsParametersDescription
  Git SettingsParametersGit
  Docker SettingsParametersDocker
  Paths SettingsParametersPaths
  Building SettingsParametersBuilding
  Deploying SettingsParametersDeploying
  SystemServices SettingsParametersSystemServices
  ApplicatonServices SettingsParametersApplicatonServices
}