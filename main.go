package main

import (
    "Mosquito/broker"
    "Mosquito/config"
    "log"
)

func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }
    
    broker.StartBroker(cfg)
}