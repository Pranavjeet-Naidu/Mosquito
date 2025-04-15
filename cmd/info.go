/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
    "fmt"

    "github.com/Pranavjeet-Naidu/Mosquito/steg"
    "github.com/spf13/cobra"
)

var (
    infoImagePath string
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
    Use:   "info",
    Short: "Show information about an image",
    Long: `Show information about an image and its steganography capacity.
    
Example:
  mosquito info -i image.png`,
    Run: func(cmd *cobra.Command, args []string) {
        if infoImagePath == "" {
            fmt.Println("Error: Image path is required")
            cmd.Help()
            return
        }

        // Load the image
        img, err := steg.LoadImage(infoImagePath)
        if err != nil {
            fmt.Printf("Error loading image: %v\n", err)
            return
        }

        // Get basic image info
        width, height, _ := steg.ImageInfo(img)
        
        fmt.Println("Image Information:")
        fmt.Printf("  File: %s\n", infoImagePath)
        fmt.Printf("  Dimensions: %d x %d pixels\n", width, height)
        fmt.Printf("  Total pixels: %d\n", width*height)
        fmt.Printf("  Type: %s\n", func() string {
            if steg.IsGrayscale(img) {
                return "Grayscale"
            }
            return "Color"
        }())
        
        // Check if it's a steganography image
        if steg.IsStegImage(img) {
            fmt.Println("\nSteganography Information:")
            header, err := steg.GetImageInfo(img)
            if err != nil {
                fmt.Printf("  Error reading steganography header: %v\n", err)
            } else {
                fmt.Printf("  Mode: %s\n", steg.ModeNames[header.Mode])
                fmt.Printf("  Payload size: %d bytes\n", header.PayloadLen)
                fmt.Printf("  Contains: %s\n", func() string {
                    if header.IsImage() {
                        return "Image data"
                    }
                    return "Text/binary data"
                }())
                fmt.Printf("  Encryption: %s\n", func() string {
                    if header.IsEncrypted() {
                        return "Encrypted (password required)"
                    }
                    return "Not encrypted"
                }())
            }
        } else {
            fmt.Println("\nSteganography Capacity:")
            
            for _, mode := range steg.GetAvailableModes() {
                maxBytes := steg.CalculateMaxPayloadSize(img, mode)
                
                fmt.Printf("  %s: %d bytes (%.1f KB)\n", 
                    steg.ModeNames[mode], 
                    maxBytes, 
                    float64(maxBytes)/1024.0,
                )
            }
            
            fmt.Println("\nRecommended Mode:")
            if width*height < 1000 {
                fmt.Println("  This image is very small. Use LSB8 for maximum capacity.")
            } else if width*height < 10000 {
                fmt.Println("  This image is small. LSB3 or LSB4 recommended for balance.")
            } else {
                fmt.Println("  This image is large enough for LSB1 to hide most messages securely.")
            }
        }
    },
}

func init() {
    rootCmd.AddCommand(infoCmd)

    // Add flags
    infoCmd.Flags().StringVarP(&infoImagePath, "image", "i", "", "Image path (required)")

    // Mark required flags
    infoCmd.MarkFlagRequired("image")
}