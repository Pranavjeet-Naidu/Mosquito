package steg

import (
    "image"
    "image/color"
    "fmt"
    "os"
    "path/filepath"
    "strings"
    
    // Import all image formats
    "image/gif"
    "image/jpeg"
    "image/png"
    _ "golang.org/x/image/bmp"
    _ "golang.org/x/image/tiff"
    _ "golang.org/x/image/webp"
)

// LoadImage loads an image from a file
func LoadImage(path string) (image.Image, error) {
    // Check if the file exists first
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return nil, fmt.Errorf("image file not found: %s", path)
    }
    
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    
    img, _, err := image.Decode(f)
    if err != nil {
        return nil, err
    }
    
    return img, nil
}

// SaveImage saves an image to a file with appropriate format based on extension
func SaveImage(img image.Image, path string) error {
    f, err := os.Create(path)
    if err != nil {
        return err
    }
    defer f.Close()
    
    ext := strings.ToLower(filepath.Ext(path))
    
    switch ext {
    case ".jpg", ".jpeg":
        return jpeg.Encode(f, img, &jpeg.Options{Quality: 95})
    case ".png":
        return png.Encode(f, img)
    case ".gif":
        return gif.Encode(f, img, &gif.Options{NumColors: 256})
    default:
        // Default to PNG if extension not recognized
        return png.Encode(f, img)
    }
}

// ImageInfo returns information about an image
func ImageInfo(img image.Image) (width, height, maxBits int) {
    bounds := img.Bounds()
    width = bounds.Dx()
    height = bounds.Dy()
    maxBits = width * height * 3 // RGB channels
    return
}

// CalculateMaxPayloadSize calculates the maximum payload size in bytes for an image
func CalculateMaxPayloadSize(img image.Image, mode StegMode) int {
    width, height, _ := ImageInfo(img)
    pixelCount := width * height
    
    // Calculate bits per pixel for the mode
    bitsPerPixel := mode.CapacityFactor()
    
    // Total capacity in bits
    totalBits := pixelCount * bitsPerPixel
    
    // Header size in bits (8 bytes * 8 bits)
    headerBits := 8 * 8
    
    // Available bits for payload
    availableBits := totalBits - headerBits
    
    // Convert to bytes (integer division)
    return availableBits / 8
}

// ConvertToRGBA converts any image to RGBA format
func ConvertToRGBA(img image.Image) *image.RGBA {
    if rgba, ok := img.(*image.RGBA); ok {
        return rgba
    }
    
    bounds := img.Bounds()
    rgba := image.NewRGBA(bounds)
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            rgba.Set(x, y, img.At(x, y))
        }
    }
    
    return rgba
}

// MeasureImageDifference calculates the average pixel difference between two images
func MeasureImageDifference(img1, img2 image.Image) float64 {
    bounds1 := img1.Bounds()
    bounds2 := img2.Bounds()
    
    if bounds1.Dx() != bounds2.Dx() || bounds1.Dy() != bounds2.Dy() {
        return 1.0 // Maximum difference
    }
    
    totalDiff := 0.0
    pixelCount := bounds1.Dx() * bounds1.Dy()
    
    for y := bounds1.Min.Y; y < bounds1.Max.Y; y++ {
        for x := bounds1.Min.X; x < bounds1.Max.X; x++ {
            r1, g1, b1, _ := img1.At(x, y).RGBA()
            r2, g2, b2, _ := img2.At(x, y).RGBA()
            
            // Calculate color channel differences (normalized to 0-255)
            rDiff := float64(abs(int(r1>>8) - int(r2>>8)))
            gDiff := float64(abs(int(g1>>8) - int(g2>>8)))
            bDiff := float64(abs(int(b1>>8) - int(b2>>8)))
            
            // Average channel difference for this pixel (0-255)
            pixelDiff := (rDiff + gDiff + bDiff) / 3.0
            
            // Add to total (normalized to 0-1)
            totalDiff += pixelDiff / 255.0
        }
    }
    
    // Return average difference (0-1)
    return totalDiff / float64(pixelCount)
}

// abs returns the absolute value of x
func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}

// IsGrayscale checks if an image is grayscale
func IsGrayscale(img image.Image) bool {
    bounds := img.Bounds()
    
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            r, g, b, _ := img.At(x, y).RGBA()
            
            // If R, G, B values are not all equal, the image is not grayscale
            if r != g || g != b {
                return false
            }
        }
    }
    
    return true
}

// CreateColorGrid creates a test image with a grid of colors
func CreateColorGrid(width, height int) *image.RGBA {
    img := image.NewRGBA(image.Rect(0, 0, width, height))
    
    // Define colors for the grid
    colors := []color.RGBA{
        {255, 0, 0, 255},    // Red
        {0, 255, 0, 255},    // Green
        {0, 0, 255, 255},    // Blue
        {255, 255, 0, 255},  // Yellow
        {0, 255, 255, 255},  // Cyan
        {255, 0, 255, 255},  // Magenta
        {255, 255, 255, 255},// White
        {0, 0, 0, 255},      // Black
    }
    
    blockWidth := width / 4
    blockHeight := height / 2
    
    for y := 0; y < 2; y++ {
        for x := 0; x < 4; x++ {
            colorIndex := y*4 + x
            if colorIndex < len(colors) {
                fillRect(img, x*blockWidth, y*blockHeight, (x+1)*blockWidth, (y+1)*blockHeight, colors[colorIndex])
            }
        }
    }
    
    return img
}

// fillRect fills a rectangle with a color
func fillRect(img *image.RGBA, x1, y1, x2, y2 int, col color.RGBA) {
    for y := y1; y < y2; y++ {
        for x := x1; x < x2; x++ {
            img.Set(x, y, col)
        }
    }
}