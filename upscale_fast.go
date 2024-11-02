package stabilityai

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type UpscaleFastRequest struct {
	Image  io.Reader
	Output Output
}

func (c *Client) UpscaleFast(ctx context.Context, req UpscaleFastRequest, out io.Writer) error {
	buff := bytes.Buffer{}
	form := multipart.NewWriter(&buff)

	img, err := form.CreateFormFile("image", "image")
	if err != nil {
		return fmt.Errorf("failed to create form field image: %w", err)
	}

	if _, err := io.Copy(img, req.Image); err != nil {
		return fmt.Errorf("failed to copy image to form field: %w", err)
	}

	if req.Output != "" {
		if err := form.WriteField("output_format", string(req.Output)); err != nil {
			return fmt.Errorf("failed to write field output_format: %w", err)
		}
	}

	if err := form.Close(); err != nil {
		return fmt.Errorf("failed to close multipart writer: %w", err)
	}

	url := "https://api.stability.ai/v2beta/stable-image/upscale/fast"
	rq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &buff)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	rq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	rq.Header.Set("Content-Type", form.FormDataContentType())
	rq.Header.Set("Accept", "image/*")

	resp, err := c.httpClient.Do(rq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("failed to copy response body: %w", err)
	}

	return nil
}
