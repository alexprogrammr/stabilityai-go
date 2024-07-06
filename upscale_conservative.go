package stabilityai

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type UpscaleConservativeRequest struct {
	Image          io.Reader
	Prompt         string
	NegativePrompt string
	Seed           int
	Output         Output
	Creativity     float64
}

func (c *Client) UpscaleConservative(ctx context.Context, req UpscaleConservativeRequest, out io.Writer) error {
	buff := bytes.Buffer{}
	form := multipart.NewWriter(&buff)

	img, err := form.CreateFormFile("image", "image")
	if err != nil {
		return fmt.Errorf("failed to create form field image: %w", err)
	}

	if _, err := io.Copy(img, req.Image); err != nil {
		return fmt.Errorf("failed to copy image to form field: %w", err)
	}

	if err := form.WriteField("prompt", req.Prompt); err != nil {
		return fmt.Errorf("failed to write field prompt: %w", err)
	}

	if req.NegativePrompt != "" {
		if err := form.WriteField("negative_prompt", req.NegativePrompt); err != nil {
			return fmt.Errorf("failed to write field negative_prompt: %w", err)
		}
	}
	if req.Seed > 0 {
		if err := form.WriteField("seed", fmt.Sprintf("%d", req.Seed)); err != nil {
			return fmt.Errorf("failed to write field seed: %w", err)
		}
	}
	if req.Output != "" {
		if err := form.WriteField("output_format", string(req.Output)); err != nil {
			return fmt.Errorf("failed to write field output_format: %w", err)
		}
	}
	if req.Creativity > 0 {
		if err := form.WriteField("creativity", fmt.Sprintf("%f", req.Creativity)); err != nil {
			return fmt.Errorf("failed to write field creativity: %w", err)
		}
	}

	if err := form.Close(); err != nil {
		return fmt.Errorf("failed to close multipart writer: %w", err)
	}

	url := "https://api.stability.ai/v2beta/stable-image/upscale/conservative"
	rq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &buff)
	if err != nil {
		return fmt.Errorf("failed to create request %s: %w", url, err)
	}

	rq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	rq.Header.Set("Content-Type", form.FormDataContentType())
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
