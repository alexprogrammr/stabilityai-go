package stabilityai

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type AspectRatio string

const (
	AspectRatio1x1  AspectRatio = "1:1"
	AspectRatio3x2  AspectRatio = "3:2"
	AspectRatio2x3  AspectRatio = "2:3"
	AspectRatio5x4  AspectRatio = "5:4"
	AspectRatio4x5  AspectRatio = "4:5"
	AspectRatio16x9 AspectRatio = "16:9"
	AspectRatio9x16 AspectRatio = "9:16"
	AspectRatio21x9 AspectRatio = "21:9"
	AspectRatio9x21 AspectRatio = "9:21"
)

type Style string

const (
	Style3dModel      Style = "3d-model"
	StyleAnalogFilm   Style = "analog-film"
	StyleAnime        Style = "anime"
	StyleCinematic    Style = "cinematic"
	StyleComicBook    Style = "comic-book"
	StyleDigitalArt   Style = "digital-art"
	StyleEnhance      Style = "enhance"
	StyleFantasyArt   Style = "fantasy-art"
	StyleIsometric    Style = "isometric"
	StyleLineArt      Style = "line-art"
	StyleLowPoly      Style = "low-poly"
	StyleModeling     Style = "modeling-compound"
	StyleNeonPunk     Style = "neon-punk"
	StyleOrigami      Style = "origami"
	StylePhotographic Style = "photographic"
	StylePixelArt     Style = "pixel-art"
	StyleTileTexture  Style = "tile-texture"
)

type Output string

const (
	OutputPNG  Output = "png"
	OutputJPEG Output = "jpeg"
	OutputWEBP Output = "webp"
)

type GenerateRequest struct {
	Prompt      string
	AspectRatio AspectRatio
	Style       Style
	Output      Output
}

func (c *Client) Generate(ctx context.Context, req GenerateRequest, out io.Writer) error {
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
	if req.Style != "" {
		if err := w.WriteField("style_preset", string(req.Style)); err != nil {
			return fmt.Errorf("failed to write field style_preset: %w", err)
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

	url := "https://api.stability.ai/v2beta/stable-image/generate/core"
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
