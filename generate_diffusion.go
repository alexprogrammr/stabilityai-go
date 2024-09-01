package stabilityai

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

func (c *Client) GenerateDiffusion(ctx context.Context, req GenerateRequest, out io.Writer) error {
	body := bytes.Buffer{}
	w := multipart.NewWriter(&body)

	if err := w.WriteField("prompt", req.Prompt); err != nil {
		return fmt.Errorf("failed to write field prompt: %w", err)
	}

	if req.AspectRatio != "" {
		if err := w.WriteField("aspect_ratio", string(req.AspectRatio)); err != nil {
			return fmt.Errorf("failed to write field aspect_ratio: %w", err)
		}
	}
	if req.Output != "" {
		if err := w.WriteField("output_format", string(req.Output)); err != nil {
			return fmt.Errorf("failed to write field output_format: %w", err)
		}
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("failed to close multipart writer: %w", err)
	}

	url := "https://api.stability.ai/v2beta/stable-image/generate/sd3"
	rq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &body)
	if err != nil {
		return fmt.Errorf("failed to create request %s: %w", url, err)
	}

	rq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	rq.Header.Set("Content-Type", w.FormDataContentType())
	rq.Header.Set("Accept", "image/*")

	resp, err := c.httpClient.Do(rq)
	if err != nil {
		return fmt.Errorf("failed to send request %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %s: %d", url, resp.StatusCode)
	}

	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("failed to copy response body %s: %w", url, err)
	}

	return nil
}
