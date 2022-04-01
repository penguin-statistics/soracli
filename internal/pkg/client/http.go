package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

type Penguin struct {
	baseUrl string
	token   string
	client  *http.Client
}

func NewHTTP(baseUrl, token string) *Penguin {
	return &Penguin{
		baseUrl: baseUrl,
		token:   token,
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func NewHTTPFromCliContext(ctx *cli.Context) *Penguin {
	log.Debug().Str("baseUrl", ctx.String("baseUrl")).Str("token", ctx.String("token")).Msg("creating http client")
	return NewHTTP(ctx.String("baseUrl"), ctx.String("token"))
}

func (h *Penguin) NewRequest(method, url string, body any) (*http.Request, error) {
	req, err := http.NewRequest(method, h.baseUrl+url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+h.token)

	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		req.Body = io.NopCloser(bytes.NewReader(b))
	}

	return req, nil
}

func (h *Penguin) GetJSON(url string, v any) error {
	req, err := h.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	log.Debug().Str("method", req.Method).Str("url", req.URL.String()).Msg("making request")

	resp, err := h.client.Do(req)
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return errors.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(v)

	return nil
}

func (h *Penguin) PostJSON(url string, v any) error {
	req, err := h.NewRequest("POST", url, v)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	log.Debug().Str("method", req.Method).Str("url", req.URL.String()).Msg("making request")

	resp, err := h.client.Do(req)
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		// log response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error().Err(err).Msg("failed to read response body")
		} else {
			log.Error().Str("body", string(body)).Msg("request failed")
		}
		return errors.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	return nil
}
