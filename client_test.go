package vaultclient

import (
	"encoding/json"
	"errors"
	"os"
	"testing"
)

var (
	vaultAddr  string
	vaultToken string
)

func set() {
	// set these before running tests, IE: Makefile
	os.Setenv("VAULT_ADDR", os.Getenv("LOCAL_VAULT_ADDR"))
	os.Setenv("VAULT_TOKEN", os.Getenv("LOCAL_VAULT_TOKEN"))

	vaultAddr = os.Getenv("VAULT_ADDR")
	vaultToken = os.Getenv("VAULT_TOKEN")
}

func clear() {
	os.Setenv("VAULT_ADDR", "")
	os.Setenv("VAULT_TOKEN", "")

	vaultAddr = ""
	vaultToken = ""
}

func init() {
	if os.Getenv("LOCAL_VAULT_ADDR") == "" {
		panic(errors.New("LOCAL_VAULT_ADDR is required"))
	}

	if os.Getenv("LOCAL_VAULT_TOKEN") == "" {
		panic(errors.New("LOCAL_VAULT_TOKEN is required"))
	}

	clear()
}

func TestNew(t *testing.T) {
	set()
	v, err := New(true)
	if err != nil {
		t.Errorf("Expected New() to return a nil error, got: " + err.Error())
	}

	addr := v.Client.Address()
	if addr != vaultAddr {
		t.Errorf("Expected New() to return a Vault client with address " + vaultAddr + ", got: " + addr)
	}

	token := v.Client.Token()
	if token != vaultToken {
		t.Errorf("Expected New() to return a Vault client with address " + vaultToken + ", got: " + token)
	}
}

// This test relies on the Makefile
//     docker exec $LOCAL_VAULT_CONTAINER_NAME \
//        vault kv put secret/test test=test >> /dev/null
func TestReadSecret(t *testing.T) {
	set()
	v, err := New(true)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	s, err := v.ReadSecretJSON("secret/test")
	if err != nil {
		t.Errorf("Expected ReadSecretJSON() to return a secret from path 'secret/test', got error: " + err.Error())
	}

	type testType struct {
		test string
	}
	test := testType{}
	if err := json.Unmarshal(s, &test); err != nil {
		t.Errorf("Expected ReadSecretJSON() to return a secret test=test, but got error: " + err.Error())
	}
}

func TestReadSecretBad(t *testing.T) {
	set()
	v, err := New(true)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	s, err := v.ReadSecret("secret/bad/bad/bad/bad")
	if s != nil {
		b, _ := json.Marshal(s)
		t.Errorf("S SHOULD BE NIL WTF " + string(b))
	}
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestWriteSecret(t *testing.T) {
	return
}
