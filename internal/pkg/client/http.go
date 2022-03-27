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

type HTTP struct {
	baseUrl string
	token   string
	client  *http.Client
}

func NewHTTP(baseUrl, token string) *HTTP {
	return &HTTP{
		baseUrl: baseUrl,
		token:   token,
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func NewHTTPFromCliContext(ctx *cli.Context) *HTTP {
	log.Debug().Str("baseUrl", ctx.String("baseUrl")).Str("token", ctx.String("token")).Msg("creating http client")
	return NewHTTP(ctx.String("baseUrl"), ctx.String("token"))
}

func (h *HTTP) NewRequest(method, url string, body any) (*http.Request, error) {
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

func (h *HTTP) GetJSON(url string, v any) error {
	req, err := h.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	log.Debug().Interface("headers", req.Header).Msg("request")

	resp, err := h.client.Do(req)
	if err != nil {
		return err
	} else if resp.StatusCode != 200 {
		return errors.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(v)

	return nil
}
