package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type EnrichmentClient struct {
	client *http.Client
}

func NewEnrichmentClient() *EnrichmentClient {
	return &EnrichmentClient{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type AgeResponse struct {
	Age int `json:"age"`
}

type GenderResponse struct {
	Gender string `json:"gender"`
}

type NationalityResponse struct {
	Country []struct {
		CountryID   string  `json:"country_id"`
		Probability float64 `json:"probability"`
	} `json:"country"`
}

func (c *EnrichmentClient) GetAge(name string) (int, error) {
	resp, err := c.client.Get(fmt.Sprintf("https://api.agify.io/?name=%s", name))
	if err != nil {
		return 0, fmt.Errorf("agify request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("agify returned status: %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var ageResp AgeResponse
	if err := json.Unmarshal(body, &ageResp); err != nil {
		return 0, fmt.Errorf("failed to parse agify response: %w", err)
	}

	return ageResp.Age, nil
}

func (c *EnrichmentClient) GetGender(name string) (string, error) {
	resp, err := c.client.Get(fmt.Sprintf("https://api.genderize.io/?name=%s", name))
	if err != nil {
		return "", fmt.Errorf("genderize request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("genderize returned status: %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var genderResp GenderResponse
	if err := json.Unmarshal(body, &genderResp); err != nil {
		return "", fmt.Errorf("failed to parse genderize response: %w", err)
	}

	return genderResp.Gender, nil
}

func (c *EnrichmentClient) GetNationality(name string) (string, error) {
	resp, err := c.client.Get(fmt.Sprintf("https://api.nationalize.io/?name=%s", name))
	if err != nil {
		return "", fmt.Errorf("nationalize request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("nationalize returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if len(body) == 0 {
		return "", fmt.Errorf("empty response body from nationalize.io")
	}

	var nationalityResp NationalityResponse
	if err := json.Unmarshal(body, &nationalityResp); err != nil {
		return "", fmt.Errorf("failed to parse nationalize response: %w", err)
	}

	if len(nationalityResp.Country) == 0 {
		return "", fmt.Errorf("no nationality found")
	}

	var maxProb float64 = 0.0
	result := ""
	for _, country := range nationalityResp.Country {
		if country.Probability > maxProb {
			maxProb = country.Probability
			result = country.CountryID
		}
	}

	return result, nil
}
