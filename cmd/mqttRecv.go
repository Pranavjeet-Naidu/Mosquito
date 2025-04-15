/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"

    "github.com/Pranavjeet-Naidu/Mosquito/mqtt"
    "github.com/spf13/cobra"
)

var (
    mqttRecvBroker    string
    mqttRecvTopic     string
    mqttRecvOutputDir string
)

// mqttRecvCmd represents the mqttRecv command
var mqttRecvCmd = &cobra.Command{
    Use:   "mqttRecv",
    Short: "Receive steganographic images via MQTT",
    Long: `Subscribe to an MQTT topic and receive steganographic images.
Images will be saved to the specified output directory.
    
Example:
  mosquito mqttRecv -b tcp://broker.example.com:1883 -t stego/images -o ./received`,
    Run: func(cmd *cobra.Command, args []string) {
        if mqttRecvBroker == "" || mqttRecvTopic == "" || mqttRecvOutputDir == "" {
            fmt.Println("Error: Broker URL, topic, and output directory are required")
            cmd.Help()
            return
        }

        // Ensure output directory exists
        if err := os.MkdirAll(mqttRecvOutputDir, 0755); err != nil {
            fmt.Printf("Error creating output directory: %v\n", err)
            return
        }

        // Setup signal handling for graceful shutdown
        sigChan := make(chan os.Signal, 1)
        signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

        // Start receiving messages
        client, err := mqtt.SubscribeForImages(mqttRecvBroker, mqttRecvTopic, mqttRecvOutputDir)
        if err != nil {
            fmt.Printf("Error subscribing: %v\n", err)
            return
        }

        fmt.Printf("Subscribed to %s on topic %s\n", mqttRecvBroker, mqttRecvTopic)
        fmt.Println("Waiting for images... (Press Ctrl+C to stop)")

        // Wait for termination signal
        <-sigChan
        fmt.Println("\nShutting down...")
        client.Disconnect(250)
    },
}

func init() {
    rootCmd.AddCommand(mqttRecvCmd)

    // Add flags
    mqttRecvCmd.Flags().StringVarP(&mqttRecvBroker, "broker", "b", "", "MQTT broker URL (required)")
    mqttRecvCmd.Flags().StringVarP(&mqttRecvTopic, "topic", "t", "", "MQTT topic to subscribe to (required)")
    mqttRecvCmd.Flags().StringVarP(&mqttRecvOutputDir, "output", "o", "", "Directory to save received images (required)")

    // Mark required flags
    mqttRecvCmd.MarkFlagRequired("broker")
    mqttRecvCmd.MarkFlagRequired("topic")
    mqttRecvCmd.MarkFlagRequired("output")
}