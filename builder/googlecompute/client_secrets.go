package googlecompute

import (
	"encoding/json"
	"io/ioutil"
)

// clientSecrets represents the client secrets of a Google Compute Engine
// service account.
type clientSecrets struct {
	Web struct {
		AuthURI     string `json:"auth_uri"`
		ClientEmail string `json:"client_email"`
		ClientId    string `json:"client_id"`
		TokenURI    string `json:"token_uri"`
	}
}

// loadClientSecrets
func loadClientSecrets(path string) (*clientSecrets, error) {
	var cs *clientSecrets
	secretBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(secretBytes, &cs)
	if err != nil {
		return nil, err
	}
	return cs, nil
}
