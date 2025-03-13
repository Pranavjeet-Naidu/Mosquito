// client/publisher.go
package main

import (
	"encoding/binary"
	"net"
	"os"
	"stego"
)

func main() {
	// Read sample image
	imgBytes, _ := os.ReadFile("cover.png")
	
	// Hide secret in image
	secretImg, _ := stego.EncodeInImage(imgBytes, "TOP SECRET")
	
	// Connect to broker
	conn, _ := net.Dial("tcp", "localhost:1883")
	defer conn.Close()
  
	// Create message
	topic := stego.EncodeInTopic("sensors/cam", "hidden123")
	payload := secretImg
  
	// Send framed message
	header := make([]byte, 6)
	binary.BigEndian.PutUint16(header[:2], uint16(len(topic)))
	binary.BigEndian.PutUint32(header[2:6], uint32(len(payload)))
  
	conn.Write(header)
	conn.Write([]byte(topic))
	conn.Write(payload)
  }