/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
    "fmt"

    "github.com/Pranavjeet-Naidu/Mosquito/mqtt"
    "github.com/spf13/cobra"
)

var (
    mqttSendBroker string
    mqttSendTopic  string
    mqttSendImage  string
)

// mqttSendCmd represents the mqttSend command
var mqttSendCmd = &cobra.Command{
    Use:   "mqttSend",
    Short: "Send a steganographic image via MQTT",
    Long: `Send a steganographic image to an MQTT broker.
    
Example:
  mosquito mqttSend -b tcp://broker.example.com:1883 -t stego/images -i stego.png`,
    Run: func(cmd *cobra.Command, args []string) {
        if mqttSendBroker == "" || mqttSendTopic == "" || mqttSendImage == "" {
            fmt.Println("Error: Broker URL, topic, and image path are required")
            cmd.Help()
            return
        }

        err := mqtt.PublishImage(mqttSendBroker, mqttSendTopic, mqttSendImage)
        if err != nil {
            fmt.Printf("Error sending image: %v\n", err)
            return
        }

        fmt.Printf("Image successfully sent to %s on topic %s\n", mqttSendBroker, mqttSendTopic)
    },
}

func init() {
    rootCmd.AddCommand(mqttSendCmd)

    // Add flags
    mqttSendCmd.Flags().StringVarP(&mqttSendBroker, "broker", "b", "", "MQTT broker URL (required)")
    mqttSendCmd.Flags().StringVarP(&mqttSendTopic, "topic", "t", "", "MQTT topic (required)")
    mqttSendCmd.Flags().StringVarP(&mqttSendImage, "image", "i", "", "Image path to send (required)")

    // Mark required flags
    mqttSendCmd.MarkFlagRequired("broker")
    mqttSendCmd.MarkFlagRequired("topic")
    mqttSendCmd.MarkFlagRequired("image")
}