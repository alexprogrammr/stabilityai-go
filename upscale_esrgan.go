package stabilityai

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type UpscaleESRGANRequest struct {
	Image  io.Reader
	Width  int
	Height int
}

func (c *Client) UpscaleESRGAN(ctx context.Context, req UpscaleESRGANRequest, out io.Writer) error {
	buff := bytes.Buffer{}
	form := multipart.NewWriter(&buff)

	img, err := form.CreateFormField("image")
	if err != nil {
		return fmt.Errorf("failed to create form field image: %w", err)
	}

	if _, err := io.Copy(img, req.Image); err != nil {
		return fmt.Errorf("failed to copy image to form field: %w", err)
	}

	if req.Width > 0 {
		if err := form.WriteField("width", fmt.Sprintf("%d", req.Width)); err != nil {
			return fmt.Errorf("failed to write field width: %w", err)
		}
	}
	if req.Height > 0 {
		if err := form.WriteField("height", fmt.Sprintf("%d", req.Height)); err != nil {
			return fmt.Errorf("failed to write field height: %w", err)
		}
	}

	if err := form.Close(); err != nil {
		return fmt.Errorf("failed to close multipart writer: %w", err)
	}

	url := "https://api.stability.ai/v1/generation/esrgan-v1-x2plus/image-to-image/upscale"
	rq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &buff)
	if err != nil {
		return fmt.Errorf("failed to create request %s: %w", url, err)
	}

	rq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	rq.Header.Set("Content-Type", form.FormDataContentType())
	rq.Header.Set("Accept", "image/png")

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
