package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"image-gallery/internal/auth/config"
	"io"
	"net/http"
	"time"
)

type UserTransport struct {
	config config.UserTransport
	logger *zap.SugaredLogger
}

func NewTransport(config config.UserTransport, logger *zap.SugaredLogger) *UserTransport {
	return &UserTransport{
		config: config,
		logger: logger,
	}
}

type GetUserResponse struct {
	Id       int    `json:"Id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (ut *UserTransport) GetUser(ctx context.Context, email string) (*GetUserResponse, error) {
	var response *GetUserResponse

	responseBody, err := ut.makeRequest(ctx, "GET", fmt.Sprintf("/api/v1/user/%s", email), ut.config.Timeout)

	if err != nil {
		ut.logger.Info("Fail in response body")
		return nil, err
	}

	if err := json.Unmarshal(responseBody, &response); err != nil {
		ut.logger.Info("Fail in Unmarshalling")
		return nil, err
	}

	return response, nil
}

func (ut *UserTransport) makeRequest(ctx context.Context, httpMethod string, endpoint string, timeout time.Duration) (b []byte, err error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	requestURL := ut.config.Host + endpoint

	req, err := http.NewRequest(httpMethod, requestURL, nil)

	if err != nil {
		ut.logger.Info("Failed at request in make request method")
		return nil, err
	}
	httpClient := &http.Client{}

	res, err := httpClient.Do(req)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			ut.logger.Info("error in body close")
		}
	}(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
