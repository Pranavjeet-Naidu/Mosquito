/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/Pranavjeet-Naidu/Mosquito/steg"
    "github.com/spf13/cobra"
)

var (
    extractInputImage string
    extractOutputFile string
    extractShowText   bool
    extractPassword   string
    extractInfo       bool
)

// extractCmd represents the extract command
var extractCmd = &cobra.Command{
    Use:   "extract",
    Short: "Extract hidden data from an image",
    Long: `Extract hidden data from a steganographic image.
    
Example:
  mosquito extract -i stego.png -o extracted_data.bin
  mosquito extract -i stego.png -t                     # Display text message
  mosquito extract -i stego.png -o secret.jpg -p pass  # Extract with password
  mosquito extract -i stego.png --info                 # Show steganography info`,
    Run: func(cmd *cobra.Command, args []string) {
        if extractInputImage == "" {
            fmt.Println("Error: Input image path is required")
            cmd.Help()
            return
        }

        // Load the steganographic image
        img, err := steg.LoadImage(extractInputImage)
        if err != nil {
            fmt.Printf("Error loading image: %v\n", err)
            return
        }

        // Check if this is a steganographic image
        if !steg.IsStegImage(img) {
            fmt.Println("Error: The image does not appear to contain hidden data")
            return
        }

        // Just show info about the steganographic image if requested
        if extractInfo {
            header, err := steg.GetImageInfo(img)
            if err != nil {
                fmt.Printf("Error reading steganography header: %v\n", err)
                return
            }
            
            fmt.Println("Steganographic Image Information:")
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
            fmt.Printf("  Compression: %s\n", func() string {
                if header.IsCompressed() {
                    return "Compressed"
                }
                return "Not compressed"
            }())
            return
        }

        // Extract the hidden data
        var data []byte
        if extractPassword != "" {
            data, err = steg.DecodeMessageWithPassword(img, extractPassword)
            if err != nil {
                fmt.Printf("Error extracting data: %v\n", err)
                return
            }
        } else {
            data, err = steg.DecodeMessage(img)
            if err != nil {
                fmt.Printf("Error extracting data: %v\n", err)
                return
            }
        }

        // Get header to check if this is an image
        header, _ := steg.GetImageInfo(img)
        isImage := header.IsImage()

        if extractShowText && !isImage {
            // Display the extracted data as text
            fmt.Println("Extracted message:")
            fmt.Println(string(data))
        } else if extractOutputFile != "" {
            // Save the extracted data to a file
            err = os.WriteFile(extractOutputFile, data, 0644)
            if err != nil {
                fmt.Printf("Error writing output file: %v\n", err)
                return
            }
            fmt.Printf("Data successfully extracted to %s\n", extractOutputFile)
            
            // If extracted data is an image, try to determine format
            if isImage {
                fmt.Println("Extracted data appears to be an image")
                
                // Check file extension
                ext := filepath.Ext(extractOutputFile)
                if ext == "" || ext == ".bin" || ext == ".dat" {
                    fmt.Println("Note: You may need to rename the file with an appropriate image extension (.png, .jpg, etc.)")
                }
            }
        } else {
            // If no output file specified and text display not requested,
            // or if it's an image and text display was requested
            if isImage {
                fmt.Println("Extracted data is an image. Please specify an output file with -o to save it.")
            } else {
                // Default to showing the data as text
                fmt.Println("Extracted message:")
                fmt.Println(string(data))
            }
        }
    },
}

func init() {
    rootCmd.AddCommand(extractCmd)

    // Add flags
    extractCmd.Flags().StringVarP(&extractInputImage, "input", "i", "", "Steganographic image path (required)")
    extractCmd.Flags().StringVarP(&extractOutputFile, "output", "o", "", "Output file for extracted data")
    extractCmd.Flags().BoolVarP(&extractShowText, "text", "t", false, "Display extracted data as text")
    extractCmd.Flags().StringVarP(&extractPassword, "password", "p", "", "Password for decrypting the data")
    extractCmd.Flags().BoolVar(&extractInfo, "info", false, "Show information about the steganographic image")

    // Mark required flags
    extractCmd.MarkFlagRequired("input")
}