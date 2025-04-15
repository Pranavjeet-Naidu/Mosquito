package mqtt

import (
    "fmt"
    "os"
    "path/filepath"
    "time"

    MQTT "github.com/eclipse/paho.mqtt.golang"
)

func PublishImage(broker, topic, imgPath string) error {
    opts := MQTT.NewClientOptions().AddBroker(broker)
    client := MQTT.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        return token.Error()
    }

    data, err := os.ReadFile(imgPath)
    if err != nil {
        return err
    }

    token := client.Publish(topic, 0, false, data)
    token.Wait()
    client.Disconnect(250)
    return nil
}

func SubscribeForImages(broker, topic, outputDir string) (MQTT.Client, error) {
    opts := MQTT.NewClientOptions().AddBroker(broker)
    
    // Create a unique client ID
    opts.SetClientID(fmt.Sprintf("mosquito-receiver-%d", time.Now().Unix()))
    
    // Set the message handler
    opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
        // Generate a filename with timestamp
        timestamp := time.Now().Format("20060102-150405")
        filename := filepath.Join(outputDir, fmt.Sprintf("received-%s.png", timestamp))
        
        // Save the received data as an image
        err := os.WriteFile(filename, msg.Payload(), 0644)
        if err != nil {
            fmt.Printf("Error saving received image: %v\n", err)
            return
        }
        
        fmt.Printf("Received image saved to: %s\n", filename)
    })
    
    // Connect to the broker
    client := MQTT.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        return nil, token.Error()
    }
    
    // Subscribe to the topic
    if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
        client.Disconnect(250)
        return nil, token.Error()
    }
    
    return client, nil
}