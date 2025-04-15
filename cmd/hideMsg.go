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
    hideMsgInputImage  string
    hideMsgOutputImage string
    hideMsgText        string
    hideMsgFile        string
    hideMsgPassword    string
    hideMsgMode        int
)

// hideMsgCmd represents the hideMsg command
var hideMsgCmd = &cobra.Command{
    Use:   "hideMsg",
    Short: "Hide a text message inside an image",
    Long: `Hide a text message inside an image using steganography.
    
Example:
  mosquito hideMsg -i input.png -o output.png -m "Secret message"
  mosquito hideMsg -i input.png -o output.png -f message.txt
  mosquito hideMsg -i input.png -o output.png -m "Secret message" -p mypassword -M 3`,
    Run: func(cmd *cobra.Command, args []string) {
        if hideMsgInputImage == "" || hideMsgOutputImage == "" {
            fmt.Println("Error: Input and output image paths are required")
            cmd.Help()
            return
        }

        if hideMsgText == "" && hideMsgFile == "" {
            fmt.Println("Error: Either a message or message file must be provided")
            cmd.Help()
            return
        }

        var message []byte
        var err error

        if hideMsgFile != "" {
            message, err = os.ReadFile(hideMsgFile)
            if err != nil {
                fmt.Printf("Error reading message file: %v\n", err)
                return
            }
        } else {
            message = []byte(hideMsgText)
        }

        // Load the input image
        img, err := steg.LoadImage(hideMsgInputImage)
        if err != nil {
            fmt.Printf("Error loading image: %v\n", err)
            return
        }

        // Verify the mode is valid
        modes := steg.GetAvailableModes()
        if hideMsgMode < 0 || hideMsgMode >= len(modes) {
            fmt.Printf("Error: Invalid mode specified. Valid modes are 0-%d\n", len(modes)-1)
            for i, mode := range modes {
                fmt.Printf("  %d: %s\n", i, steg.ModeNames[mode])
            }
            return
        }
        
        selectedMode := modes[hideMsgMode]
        
        // Check if the image has enough capacity
        hasCapacity, available, required := steg.Capacity(img, len(message), selectedMode)
        if !hasCapacity {
            fmt.Printf("Error: Image too small to encode message\n")
            fmt.Printf("  Required: %d bytes, Available: %d bytes\n", required, available)
            fmt.Printf("  Try using a different mode with higher capacity (current: %s)\n", steg.ModeNames[selectedMode])
            return
        }

        // Encode the message
        var encoded image.Image
        if hideMsgPassword != "" {
            encoded, err = steg.EncodeMessageWithPassword(img, message, hideMsgPassword, selectedMode, false)
            if err != nil {
                fmt.Printf("Error encoding message: %v\n", err)
                return
            }
            fmt.Println("Message encrypted with provided password")
        } else {
            encoded, err = steg.EncodeMessage(img, message, selectedMode)
            if err != nil {
                fmt.Printf("Error encoding message: %v\n", err)
                return
            }
        }

        // Save the output image
        err = steg.SaveImage(encoded, hideMsgOutputImage)
        if err != nil {
            fmt.Printf("Error saving image: %v\n", err)
            return
        }

        fmt.Printf("Message successfully hidden in %s using %s\n", hideMsgOutputImage, steg.ModeNames[selectedMode])
        
        // Calculate detection metrics
        diff := steg.MeasureImageDifference(img, encoded)
        fmt.Printf("Image difference: %.2f%% (lower is better)\n", diff*100)
    },
}

func init() {
    rootCmd.AddCommand(hideMsgCmd)

    // Add flags
    hideMsgCmd.Flags().StringVarP(&hideMsgInputImage, "input", "i", "", "Input image path (required)")
    hideMsgCmd.Flags().StringVarP(&hideMsgOutputImage, "output", "o", "", "Output image path (required)")
    hideMsgCmd.Flags().StringVarP(&hideMsgText, "message", "m", "", "Text message to hide")
    hideMsgCmd.Flags().StringVarP(&hideMsgFile, "file", "f", "", "File containing message to hide")
    hideMsgCmd.Flags().StringVarP(&hideMsgPassword, "password", "p", "", "Password for encrypting the message")
    hideMsgCmd.Flags().IntVarP(&hideMsgMode, "mode", "M", 0, "Steganography mode (0=LSB1, 1=LSB3, 2=LSB4, 3=LSB8)")

    // Mark required flags
    hideMsgCmd.MarkFlagRequired("input")
    hideMsgCmd.MarkFlagRequired("output")
}