package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	AgifyURL       = "https://api.agify.io/?name=%s"
	GenderizeURL   = "https://api.genderize.io/?name=%s"
	NationalizeURL = "https://api.nationalize.io/?name=%s"
)

type ExternalClient struct {
	httpClient *http.Client
}

func NewExternalClient() *ExternalClient {
	return &ExternalClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type AgifyResponse struct {
	Age   int    `json:"age"`
	Count int    `json:"count"`
	Name  string `json:"name"`
}

type GenderizeResponse struct {
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
	Count       int     `json:"count"`
	Name        string  `json:"name"`
}

type NationalizeResponse struct {
	Name    string `json:"name"`
	Country []struct {
		CountryID   string  `json:"country_id"`
		Probability float64 `json:"probability"`
	} `json:"country"`
}

func (c *ExternalClient) GetAge(name string) (int, error) {
	url := fmt.Sprintf(AgifyURL, name)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to get age: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to get age: status code %d", resp.StatusCode)
	}

	var agifyResp AgifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&agifyResp); err != nil {
		return 0, fmt.Errorf("failed to decode age response: %w", err)
	}

	return agifyResp.Age, nil
}

func (c *ExternalClient) GetGender(name string) (string, error) {
	url := fmt.Sprintf(GenderizeURL, name)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to get gender: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get gender: status code %d", resp.StatusCode)
	}

	var genderizeResp GenderizeResponse
	if err := json.NewDecoder(resp.Body).Decode(&genderizeResp); err != nil {
		return "", fmt.Errorf("failed to decode gender response: %w", err)
	}

	return genderizeResp.Gender, nil
}

func (c *ExternalClient) GetNationality(name string) (string, error) {
	url := fmt.Sprintf(NationalizeURL, name)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to get nationality: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get nationality: status code %d", resp.StatusCode)
	}

	var nationalizeResp NationalizeResponse
	if err := json.NewDecoder(resp.Body).Decode(&nationalizeResp); err != nil {
		return "", fmt.Errorf("failed to decode nationality response: %w", err)
	}

	if len(nationalizeResp.Country) == 0 {
		return "unknown", nil
	}

	return nationalizeResp.Country[0].CountryID, nil
}

func (c *ExternalClient) EnrichPerson(name string) (int, string, string, error) {
	var (
		age         int
		gender      string
		nationality string
		ageErr      error
		genderErr   error
		nationErr   error
		wg          sync.WaitGroup
	)

	wg.Add(3)

	go func() {
		defer wg.Done()
		age, ageErr = c.GetAge(name)
	}()

	go func() {
		defer wg.Done()
		gender, genderErr = c.GetGender(name)
	}()

	go func() {
		defer wg.Done()
		nationality, nationErr = c.GetNationality(name)
	}()

	wg.Wait()

	if ageErr != nil {
		return 0, "", "", fmt.Errorf("failed to get age: %w", ageErr)
	}
	if genderErr != nil {
		return 0, "", "", fmt.Errorf("failed to get gender: %w", genderErr)
	}
	if nationErr != nil {
		return 0, "", "", fmt.Errorf("failed to get nationality: %w", nationErr)
	}

	return age, gender, nationality, nil
}
