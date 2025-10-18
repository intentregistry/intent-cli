package httpclient

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/intentregistry/intent-cli/internal/config"
	"github.com/intentregistry/intent-cli/internal/version"
)

type Client struct {
	r       *resty.Client
	baseURL string
	debug   bool
}

func New(cfg config.Config) *Client {
	return NewWithDebug(cfg, false)
}

func NewWithDebug(cfg config.Config, debug bool) *Client {
	r := resty.New().
		SetBaseURL(strings.TrimRight(cfg.APIURL, "/")).
		SetHeader("User-Agent", "intent-cli/"+version.Short()).
		// sane retries
		SetRetryCount(3).
		SetRetryWaitTime(400 * time.Millisecond).
		SetRetryMaxWaitTime(3 * time.Second)

	// Add retry condition with logging
	r.AddRetryCondition(func(resp *resty.Response, err error) bool {
		shouldRetry := err != nil || resp.StatusCode() >= 500
		if shouldRetry && debug {
			if err != nil {
				fmt.Printf("[DEBUG] retry due to error: %v\n", err)
			} else {
				fmt.Printf("[DEBUG] retry due to status %d\n", resp.StatusCode())
			}
		}
		return shouldRetry
	})

	if cfg.Token != "" {
		r.SetAuthToken(cfg.Token)
	}

	// Redact Authorization header in logs
	r.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		if debug || os.Getenv("INTENT_DEBUG") == "1" {
			if auth := req.Header.Get("Authorization"); auth != "" {
				req.Header.Set("Authorization", "Bearer ***redacted***")
			}
		}
		return nil
	})

	if debug || os.Getenv("INTENT_DEBUG") == "1" {
		r.SetDebug(true)
		// Optional: limit debug body size if you expect big payloads
		// r.SetDebugBodyLimit(2048)
	} else {
		// Tame the noisy WARN/ERROR logs from Resty retries
		r.SetLogger(nil)
	}

	return &Client{r: r, baseURL: strings.TrimRight(cfg.APIURL, "/"), debug: debug}
}

func (c *Client) Get(path string, out any) error {
	if c.debug {
		fmt.Printf("[DEBUG] GET %s%s\n", c.baseURL, path)
	}
	resp, err := c.r.R().SetResult(out).Get(path)
	if err != nil {
		// Provide friendlier error messages for common network issues
		if ne, ok := err.(*net.DNSError); ok {
			return fmt.Errorf("cannot resolve API host %q (%v). Set --api-url or INTENT_API_URL", ne.Name, ne)
		}
		return fmt.Errorf("network error: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("GET %s: %s", path, resp.Status())
	}
	return nil
}

func (c *Client) PostMultipart(path string, fields map[string]any, fileField, filePath string, out any) error {
	if c.debug {
		fmt.Printf("[DEBUG] POST %s%s (multipart)\n", c.baseURL, path)
	}

	// Convert to string map for multipart form data
	fd := make(map[string]string, len(fields))
	for k, v := range fields {
		fd[k] = fmt.Sprintf("%v", v)
	}

	req := c.r.R().
		SetMultipartFormData(fd).
		SetFile(fileField, filePath)

	if out != nil {
		req = req.SetResult(out)
	}

	resp, err := req.Post(path)
	if err != nil {
		// Provide friendlier error messages for common network issues
		if ne, ok := err.(*net.DNSError); ok {
			return fmt.Errorf("cannot resolve API host %q (%v). Set --api-url or INTENT_API_URL", ne.Name, ne)
		}
		return fmt.Errorf("network error: %w", err)
	}
	if resp.IsError() {
		// include small body snippet for troubleshooting
		body := resp.String()
		if len(body) > 512 {
			body = body[:512] + "…"
		}
		return fmt.Errorf("POST %s: %s\n%s", path, resp.Status(), body)
	}
	return nil
}