package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type VoIPProvisioner struct {
	Username string
	Password string
	Endpoint string
}

type VoIPSettings struct {
	Codec         string
	Quality       string
	CallForwarding bool
	Voicemail      bool
}

type VoIPAccount struct {
	ID          string
	PhoneNumber string
	Settings    VoIPSettings
}

type VoIPProvider struct {
	Endpoint string
}

func NewVoIPProvisioner(username, password, endpoint string) *VoIPProvisioner {
	return &VoIPProvisioner{
		Username: username,
		Password: password,
		Endpoint: endpoint,
	}
}

func NewVoIPAccount(id, phoneNumber string, settings VoIPSettings) *VoIPAccount {
	return &VoIPAccount{
		ID:          id,
		PhoneNumber: phoneNumber,
		Settings:    settings,
	}
}

func NewVoIPProvider(endpoint string) *VoIPProvider {
	return &VoIPProvider{
		Endpoint: endpoint,
	}
}

func (p *VoIPProvisioner) ProvisionAccount(account *VoIPAccount) error {
	accountData := map[string]interface{}{
		"id":          account.ID,
		"phone_number": account.PhoneNumber,
		"username":     p.Username,
		"password":     p.Password,
		"settings":     account.Settings,
	}

	response, err := p.SendProvisionRequest(accountData)
	if err != nil {
		return err
	}

	if response["success"] == true {
		fmt.Println("Account provisioned successfully.")
	} else {
		fmt.Printf("Account provisioning failed. Error: %v\n", response["error"])
	}

	return nil
}

func (p *VoIPProvisioner) SendProvisionRequest(data map[string]interface{}) (map[string]interface{}, error) {
	endpointURL := p.Endpoint + "/provision"
	values, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(endpointURL, "application/json", bytes.NewBuffer(values))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func main() {
	username := "your_username"
	password := "your_password"
	endpoint := "https://voip-provider-api.com"

	provisioner := NewVoIPProvisioner(username, password, endpoint)

	accountID := "12345"
	phoneNumber := "+1234567890"

	settings := VoIPSettings{
		Codec:         "G.729",
		Quality:       "Standard",
		CallForwarding: true,
		Voicemail:      false,
	}

	account := NewVoIPAccount(accountID, phoneNumber, settings)

	if err := provisioner.ProvisionAccount(account); err != nil {
		fmt.Printf("Error provisioning account: %v\n", err)
	}
}
