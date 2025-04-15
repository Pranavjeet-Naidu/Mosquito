/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
    "os"

    "github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
    Use:   "mosquito",
    Short: "A steganography tool with MQTT support",
    Long: `Mosquito is a command-line tool for steganography and secure communication.

It allows you to:
- Hide messages and files within images
- Extract hidden content from steganographic images
- Send and receive steganographic images via MQTT

For detailed usage information, use the --help flag with any command.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
    err := rootCmd.Execute()
    if err != nil {
        os.Exit(1)
    }
}

func init() {
    // Remove the toggle flag as it's not useful
    // rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}