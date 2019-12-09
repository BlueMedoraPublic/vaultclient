package vaultclient

import (
	"io/ioutil"
	"strings"
	"encoding/json"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

// VERSION is the library version
const VERSION = "0.1.0"
const tokenError = "failed to retrieve token from environment or ~/.vault-token"

// Vault represents the Vault configuration
type Vault struct {
	Client *api.Client

	config *api.Config
}

// List represents the structured data returned by the list command
type List struct {
	Keys []string
}

// New returns a new config
func New(init bool) (v Vault, err error) {
	if init {
		err = v.Init()
	}
	return v, err
}

// Init initilizes the vault client
func (v *Vault) Init() (err error) {
	v.config = &api.Config{}
	v.config = api.DefaultConfig()
	if v.config.Error != nil {
		return v.config.Error
	}
	return v.initClient()
}

func (v *Vault) initClient() (err error) {
	v.Client, err = api.NewClient(v.config)
	if err != nil {
		return err
	}

	if len(v.Client.Token()) == 0 {
		t, err := readTokenFromFile()
		if err != nil {
			return errors.Wrap(err, tokenError)
		}
		if len(t) == 0 {
			return errors.New(tokenError)
		}
		v.Client.SetToken(t)
	}

	return nil
}

// ReadSecret returns a secret
func (v Vault) ReadSecret(path string) (*api.Secret, error) {
	s, err := v.Client.Logical().Read(path)
	if err != nil {
		return s, err
	}

	if s == nil {
		return s, errors.New("secret at path " + path + " does not exist")
	}
	return s, nil
}

// ReadSecretJSON returns a secret endcoded in json
func (v Vault) ReadSecretJSON(path string) ([]byte, error) {
	s, err := v.ReadSecret(path)
	if err != nil {
		return nil, err
	}
	return json.Marshal(s.Data)
}

// ListSecret lists secrets at a given path
func (v Vault) ListSecret(path string) (*api.Secret, error) {
	l, err := v.Client.Logical().List(path)
	if err != nil {
		return l, err
	}

	if l == nil {
		return l, errors.New("path " + path + " does not exist")
	}
	return l, nil
}

// ListSecretJSON returns a secret endcoded in json
func (v Vault) ListSecretJSON(path string) ([]byte, error) {
	s, err := v.ListSecret(path)
	if err != nil {
		return nil, err
	}
	return json.Marshal(s.Data)
}

// ListSecretStruct returns a list secret as a List object
func (v Vault) ListSecretStruct(path string) (List, error) {
	list := List{}

	s, err := v.ListSecretJSON(path)
	if err != nil {
		return list, err
	}

	err = json.Unmarshal(s, &list)
	return list, err
}

func readTokenFromFile() (string, error) {
	path, err := vaultTokenFilePath()
	if err != nil {
		return "", err
	}

	f, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(f), "\n")
	return lines[0], nil
}

func vaultTokenFilePath() (string, error) {
	h, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return h + "/.vault-token", nil
}
