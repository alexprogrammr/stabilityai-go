package stabilityai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type GenerateVideoRequest struct {
	Image          io.Reader
	CfgScale       int // [ 0 .. 10 ]
	MotionBucketId int // [ 1 .. 255 ]
}

func (c *Client) GenerateVideo(ctx context.Context, req GenerateVideoRequest, out io.Writer) error {
	id, err := c.generateVideo(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to generate video: %w", err)
	}

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			done, err := c.fetchGenerateVideoResult(ctx, id, out)
			if err != nil {
				return fmt.Errorf("failed to fetch video generation result: %w", err)
			}
			if done {
				return nil
			}
		}
	}
}

func (c *Client) generateVideo(ctx context.Context, req GenerateVideoRequest) (string, error) {
	buff := bytes.Buffer{}
	form := multipart.NewWriter(&buff)

	img, err := form.CreateFormFile("image", "image")
	if err != nil {
		return "", fmt.Errorf("failed to create form field image: %w", err)
	}

	if _, err := io.Copy(img, req.Image); err != nil {
		return "", fmt.Errorf("failed to copy image to form field: %w", err)
	}

	if err := form.WriteField("cfg_scale", fmt.Sprintf("%d", req.CfgScale)); err != nil {
		return "", fmt.Errorf("failed to write field cfg_scale: %w", err)
	}
	if err := form.WriteField("motion_bucket_id", fmt.Sprintf("%d", req.MotionBucketId)); err != nil {
		return "", fmt.Errorf("failed to write field motion_bucket_id: %w", err)
	}

	if err := form.Close(); err != nil {
		return "", fmt.Errorf("failed to close multipart writer: %w", err)
	}

	url := "https://api.stability.ai/v2beta/image-to-video"
	rq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &buff)
	if err != nil {
		return "", fmt.Errorf("failed to create request %s: %w", url, err)
	}

	rq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	rq.Header.Set("Content-Type", form.FormDataContentType())

	resp, err := c.httpClient.Do(rq)
	if err != nil {
		return "", fmt.Errorf("failed to send request %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code %s: %d", url, resp.StatusCode)
	}

	var result struct {
		Id string `json:"id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response body %s: %w", url, err)
	}

	return result.Id, nil
}

func (c *Client) fetchGenerateVideoResult(ctx context.Context, id string, out io.Writer) (bool, error) {
	url := fmt.Sprintf("https://api.stability.ai/v2beta/image-to-video/result/%s", id)
	rq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request %s: %w", url, err)
	}

	rq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	rq.Header.Set("Accept", "video/*")

	resp, err := c.httpClient.Do(rq)
	if err != nil {
		return false, fmt.Errorf("failed to send request %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusAccepted {
		return false, nil
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code %s: %d", url, resp.StatusCode)
	}

	if _, err := io.Copy(out, resp.Body); err != nil {
		return false, fmt.Errorf("failed to copy response body %s: %w", url, err)
	}

	return true, nil
}
