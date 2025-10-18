package httpclient

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/intentregistry/intent-cli/internal/config"
)

type Client struct {
	r *resty.Client
}

func New(cfg config.Config) *Client {
	r := resty.New().
		SetBaseURL(cfg.APIURL).
		SetRetryCount(3)
	if cfg.Token != "" {
		r.SetAuthToken(cfg.Token)
	}
	// opcional: logs de depuración si INTENT_DEBUG=1
	if os.Getenv("INTENT_DEBUG") == "1" {
		r.SetDebug(true)
	}
	return &Client{r: r}
}

func (c *Client) Get(path string, out any) error {
	resp, err := c.r.R().SetResult(out).Get(path)
	if err != nil { return err }
	if resp.IsError() { return fmt.Errorf("GET %s: %s", path, resp.Status()) }
	return nil
}

func (c *Client) PostMultipart(path string, fields map[string]any, fileField, filePath string, out any) error {
	req := c.r.R()
	for k, v := range fields {
		req.SetFormData(map[string]string{k: fmt.Sprintf("%v", v)})
	}
	req.SetFile(fileField, filePath)
	if out != nil { req.SetResult(out) }
	resp, err := req.Post(path)
	if err != nil { return err }
	if resp.IsError() { return fmt.Errorf("POST %s: %s", path, resp.Status()) }
	return nil
}