package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type EnrichmentClient struct {
	client *http.Client
	logger *zap.Logger
}

func NewEnrichmentClient(logger *zap.Logger) *EnrichmentClient {
	return &EnrichmentClient{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		logger: logger.With(zap.String("component", "enrichment_client")),
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
	url := fmt.Sprintf("https://api.agify.io/?name=%s", name)
	c.logger.Debug("Requesting age",
		zap.String("name", name),
		zap.String("url", url),
	)

	resp, err := c.client.Get(url)
	if err != nil {
		c.logger.Error("Age API request failed",
			zap.String("name", name),
			zap.Error(err),
		)
		return 0, fmt.Errorf("agify request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		c.logger.Error("Unexpected status from age API",
			zap.Int("status_code", resp.StatusCode),
			zap.String("body", string(body)),
			zap.String("name", name),
		)
		return 0, fmt.Errorf("agify returned status: %d", resp.StatusCode)
	}

	var ageResp AgeResponse
	if err := json.Unmarshal(body, &ageResp); err != nil {
		c.logger.Error("Failed to parse age response",
			zap.String("body", string(body)),
			zap.String("name", name),
			zap.Error(err),
		)
		return 0, fmt.Errorf("failed to parse agify response: %w", err)
	}

	c.logger.Debug("Age received",
		zap.String("name", name),
		zap.Int("age", ageResp.Age),
	)
	return ageResp.Age, nil
}

func (c *EnrichmentClient) GetGender(name string) (string, error) {
	url := fmt.Sprintf("https://api.genderize.io/?name=%s", name)
	c.logger.Debug("Requesting gender",
		zap.String("name", name),
		zap.String("url", url),
	)

	resp, err := c.client.Get(url)
	if err != nil {
		c.logger.Error("Gender API request failed",
			zap.String("name", name),
			zap.Error(err),
		)
		return "", fmt.Errorf("genderize request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		c.logger.Error("Unexpected status from gender API",
			zap.Int("status_code", resp.StatusCode),
			zap.String("body", string(body)),
			zap.String("name", name),
		)
		return "", fmt.Errorf("genderize returned status: %d", resp.StatusCode)
	}

	var genderResp GenderResponse
	if err := json.Unmarshal(body, &genderResp); err != nil {
		c.logger.Error("Failed to parse gender response",
			zap.String("body", string(body)),
			zap.String("name", name),
			zap.Error(err),
		)
		return "", fmt.Errorf("failed to parse genderize response: %w", err)
	}

	c.logger.Debug("Gender received",
		zap.String("name", name),
		zap.String("gender", genderResp.Gender),
	)
	return genderResp.Gender, nil
}
func (c *EnrichmentClient) GetNationality(name string) (string, error) {
	url := fmt.Sprintf("https://api.nationalize.io/?name=%s", name)
	c.logger.Debug("Requesting nationality",
		zap.String("name", name),
		zap.String("url", url),
	)

	resp, err := c.client.Get(url)
	if err != nil {
		c.logger.Error("Nationality API request failed",
			zap.String("name", name),
			zap.Error(err),
		)
		return "", fmt.Errorf("nationalize request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("Failed to read nationality response body",
			zap.String("name", name),
			zap.Error(err),
		)
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if len(body) == 0 {
		c.logger.Warn("Empty response from nationality API",
			zap.String("name", name),
		)
		return "", fmt.Errorf("empty response body from nationalize.io")
	}

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("Unexpected status from nationality API",
			zap.Int("status_code", resp.StatusCode),
			zap.String("body", string(body)),
			zap.String("name", name),
		)
		return "", fmt.Errorf("nationalize returned status: %d", resp.StatusCode)
	}

	var nationalityResp NationalityResponse
	if err := json.Unmarshal(body, &nationalityResp); err != nil {
		c.logger.Error("Failed to parse nationality response",
			zap.String("body", string(body)),
			zap.String("name", name),
			zap.Error(err),
		)
		return "", fmt.Errorf("failed to parse nationalize response: %w", err)
	}

	if len(nationalityResp.Country) == 0 {
		c.logger.Warn("No countries in nationality response",
			zap.String("name", name),
		)
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

	c.logger.Debug("Nationality received",
		zap.String("name", name),
		zap.String("country", result),
		zap.Float64("probability", maxProb),
	)
	return result, nil
}
