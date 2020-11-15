module github.com/BlueMedoraPublic/vaultclient

go 1.15

require (
	// we need this commit https://github.com/hashicorp/vault/pull/7611
	// therefore we import the entire vault package because the importing only
	// api will result in too old of a version for some reason (thanks go mod..)
	github.com/hashicorp/vault v1.6.0
	github.com/hashicorp/vault/api v1.0.5-0.20201001211907-38d91b749c77
	
	github.com/mitchellh/go-homedir v1.1.0
	github.com/pkg/errors v0.9.1
)
