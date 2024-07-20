package stabilityai

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
	Style       Style // Used only in GenerateCore
	Output      Output
}
