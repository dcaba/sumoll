package sumoll

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"time"
)

type (
	// HTTPSourceClient send to Resource HTTP
	HTTPSourceClient struct {
		url       *url.URL
		client    *http.Client
		UserAgent string
		headers   *http.Header
	}
)

// NewHTTPSourceClient create HTTPSourceClient object
func NewHTTPSourceClient(url *url.URL, category, hostname, sourceName *string) *HTTPSourceClient {
	return &HTTPSourceClient{
		url:       url,
		client:    &http.Client{},
		UserAgent: UserAgent(),
		headers:   getHeaders(category, hostname, sourceName),
	}
}

func getHeaders(category *string, hostname *string, sourceName *string) *http.Header {
	if category == nil && hostname == nil && sourceName == nil {
		return nil
	}
	headers := http.Header{}
	if category != nil {
		headers.Add("X-Sumo-Category", *category)
	}
	if category != nil {
		headers.Add("X-Sumo-Host", *hostname)
	}
	if category != nil {
		headers.Add("X-Sumo-Name", *sourceName)
	}
	return &headers
}

// Send object to sumologic
func (h *HTTPSourceClient) Send(body io.Reader) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req, err := h.newRequest(ctx, http.MethodPost, body)
	if err != nil {
		return err
	}

	if h.headers != nil {
		req.Header = *h.headers
	}

	if _, err := h.client.Do(req); err != nil {
		return err
	}

	return nil
}

func (h *HTTPSourceClient) newRequest(ctx context.Context, method string, body io.Reader) (*http.Request, error) {
	u := *h.url

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", h.UserAgent)

	return req, nil
}
