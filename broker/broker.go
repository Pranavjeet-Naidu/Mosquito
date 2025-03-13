// Updated broker/broker.go
func handleClient(conn net.Conn) {
	defer conn.Close()
  
	header := make([]byte, 6)
	for {
	  // Read header (2 bytes topic len, 4 bytes payload len)
	  _, err := io.ReadFull(conn, header)
	  if err != nil {
		return
	  }
  
	  topicLen := binary.BigEndian.Uint16(header[:2])
	  payloadLen := binary.BigEndian.Uint32(header[2:6])
  
	  // Read topic
	  topic := make([]byte, topicLen)
	  io.ReadFull(conn, topic)
  
	  // Read payload
	  payload := make([]byte, payloadLen)
	  io.ReadFull(conn, payload)
  
	  // Process message
	  decodedTopic := DecodeFromTopic(string(topic))
	  decodedPayload, _ := DecodeFromImage(payload)
  
	  log.Printf("Hidden in topic: %s", decodedTopic)
	  log.Printf("Hidden in payload: %s", decodedPayload)
	}
  }