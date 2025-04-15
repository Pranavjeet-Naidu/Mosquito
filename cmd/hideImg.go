/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
    "fmt"
    "os"
	"image"
    "github.com/Pranavjeet-Naidu/Mosquito/steg"
    "github.com/spf13/cobra"
)

var (
    hideImgInputImage  string
    hideImgOutputImage string
    hideImgSecretImage string
    hideImgPassword    string
    hideImgMode        int
)

// hideImgCmd represents the hideImg command
var hideImgCmd = &cobra.Command{
    Use:   "hideImg",
    Short: "Hide an image inside another image",
    Long: `Hide an image inside another image using steganography.
    
Example:
  mosquito hideImg -i cover.png -s secret.png -o output.png
  mosquito hideImg -i cover.png -s secret.png -o output.png -p mypassword -M 3`,
    Run: func(cmd *cobra.Command, args []string) {
        if hideImgInputImage == "" || hideImgOutputImage == "" || hideImgSecretImage == "" {
            fmt.Println("Error: Input, output, and secret image paths are required")
            cmd.Help()
            return
        }

        // Load the cover image
        coverImg, err := steg.LoadImage(hideImgInputImage)
        if err != nil {
            fmt.Printf("Error loading cover image: %v\n", err)
            return
        }

        // Read the secret image as binary data
        secretData, err := os.ReadFile(hideImgSecretImage)
        if err != nil {
            fmt.Printf("Error reading secret image: %v\n", err)
            return
        }

        // Verify the mode is valid
        modes := steg.GetAvailableModes()
        if hideImgMode < 0 || hideImgMode >= len(modes) {
            fmt.Printf("Error: Invalid mode specified. Valid modes are 0-%d\n", len(modes)-1)
            for i, mode := range modes {
                fmt.Printf("  %d: %s\n", i, steg.ModeNames[mode])
            }
            return
        }
        
        selectedMode := modes[hideImgMode]
        
        // Check if the image has enough capacity
        hasCapacity, available, required := steg.Capacity(coverImg, len(secretData), selectedMode)
        if !hasCapacity {
            fmt.Printf("Error: Cover image too small to encode secret image\n")
            fmt.Printf("  Required: %d bytes, Available: %d bytes\n", required, available)
            fmt.Printf("  Try using a different mode with higher capacity (current: %s)\n", steg.ModeNames[selectedMode])
            return
        }

        // Encode the image data
        var encoded image.Image
        if hideImgPassword != "" {
            encoded, err = steg.EncodeMessageWithPassword(coverImg, secretData, hideImgPassword, selectedMode, true)
            if err != nil {
                fmt.Printf("Error encoding image: %v\n", err)
                return
            }
            fmt.Println("Image encrypted with provided password")
        } else {
            // Set isImage flag to true
            header := steg.Header{
                Magic:     steg.MagicByte,
                Version:   steg.Version,
                Mode:      selectedMode,
                Flags:     steg.FlagImage,
                PayloadLen: uint32(len(secretData)),
            }
            
            // Manually create the encoded image with the header
            headerData := steg.MarshalHeader(header)
            data := append(headerData, secretData...)
            
            // Create the output image
           
            out := steg.ConvertToRGBA(coverImg)

			// Encode the data - remove the unused bounds variable
            switch selectedMode {
            case steg.LSB1:
                steg.EncodeLSB1(out, data)
            case steg.LSB3:
                steg.EncodeLSB3(out, data)
            case steg.LSB4:
                steg.EncodeLSB4(out, data)
            case steg.LSB8:
                steg.EncodeLSB8(out, data)
            }
            
            encoded = out
        }

        // Save the output image
        err = steg.SaveImage(encoded, hideImgOutputImage)
        if err != nil {
            fmt.Printf("Error saving image: %v\n", err)
            return
        }

        fmt.Printf("Image successfully hidden in %s using %s\n", hideImgOutputImage, steg.ModeNames[selectedMode])
        
        // Calculate detection metrics
        diff := steg.MeasureImageDifference(coverImg, encoded)
        fmt.Printf("Image difference: %.2f%% (lower is better)\n", diff*100)
    },
}

func init() {
    rootCmd.AddCommand(hideImgCmd)

    // Add flags
    hideImgCmd.Flags().StringVarP(&hideImgInputImage, "input", "i", "", "Cover image path (required)")
    hideImgCmd.Flags().StringVarP(&hideImgSecretImage, "secret", "s", "", "Secret image to hide (required)")
    hideImgCmd.Flags().StringVarP(&hideImgOutputImage, "output", "o", "", "Output image path (required)")
    hideImgCmd.Flags().StringVarP(&hideImgPassword, "password", "p", "", "Password for encrypting the image")
    hideImgCmd.Flags().IntVarP(&hideImgMode, "mode", "M", 0, "Steganography mode (0=LSB1, 1=LSB3, 2=LSB4, 3=LSB8)")

    // Mark required flags
    hideImgCmd.MarkFlagRequired("input")
    hideImgCmd.MarkFlagRequired("secret")
    hideImgCmd.MarkFlagRequired("output")
}