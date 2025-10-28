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
	r          *resty.Client
	baseURL    string
	debug      bool
	telemetry  bool
}

func New(cfg config.Config) *Client {
	return NewWithOptions(cfg, false, cfg.Telemetry)
}

func NewWithDebug(cfg config.Config, debug bool) *Client {
	return NewWithOptions(cfg, debug, cfg.Telemetry)
}

func NewWithOptions(cfg config.Config, debug bool, telemetry bool) *Client {
	r := resty.New().
		SetBaseURL(strings.TrimRight(cfg.APIURL, "/")).
		SetHeader("User-Agent", "intent-cli/"+version.Short()).
		// Exponential backoff retries
		SetRetryCount(3).
		SetRetryWaitTime(500 * time.Millisecond).
		SetRetryMaxWaitTime(5 * time.Second).
		SetRetryAfter(func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
			// Exponential backoff: 500ms, 1s, 2s, 4s (capped at 5s)
			attempt := resp.Request.Attempt
			waitTime := time.Duration(500) * time.Millisecond * time.Duration(1<<attempt)
			if waitTime > 5*time.Second {
				waitTime = 5 * time.Second
			}
			return waitTime, nil
		})

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

	// Add telemetry header if enabled
	if telemetry {
		r.SetHeader("X-Telemetry-Enabled", "true")
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

	return &Client{r: r, baseURL: strings.TrimRight(cfg.APIURL, "/"), debug: debug, telemetry: telemetry}
}

func (c *Client) Get(path string, out any) error {
	if c.debug {
		fmt.Printf("[DEBUG] GET %s%s\n", c.baseURL, path)
	}
	resp, err := c.r.R().SetResult(out).Get(path)
	if err != nil {
		// Provide friendlier error messages for common network issues
		if ne, ok := err.(*net.DNSError); ok {
			return fmt.Errorf("cannot resolve API host %q (%v). Try: --api-url or check connectivity", ne.Name, ne)
		}
		if opErr, ok := err.(*net.OpError); ok {
			if opErr.Timeout() {
				return fmt.Errorf("request timeout to %s. Try: --api-url or check connectivity", c.baseURL)
			}
			return fmt.Errorf("network error connecting to %s: %w. Try: --api-url or check connectivity", c.baseURL, err)
		}
		return fmt.Errorf("network error: %w. Try: --api-url or check connectivity", err)
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
			return fmt.Errorf("cannot resolve API host %q (%v). Try: --api-url or check connectivity", ne.Name, ne)
		}
		if opErr, ok := err.(*net.OpError); ok {
			if opErr.Timeout() {
				return fmt.Errorf("request timeout to %s. Try: --api-url or check connectivity", c.baseURL)
			}
			return fmt.Errorf("network error connecting to %s: %w. Try: --api-url or check connectivity", c.baseURL, err)
		}
		return fmt.Errorf("network error: %w. Try: --api-url or check connectivity", err)
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