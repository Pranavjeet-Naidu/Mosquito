// client/subscriber.go
func main() {
	conn, _ := net.Dial("tcp", "localhost:1883")
	defer conn.Close()
  
	for {
	  header := make([]byte, 6)
	  io.ReadFull(conn, header)
	  
	  topicLen := binary.BigEndian.Uint16(header[:2])
	  payloadLen := binary.BigEndian.Uint32(header[2:6])
  
	  topic := make([]byte, topicLen)
	  io.ReadFull(conn, topic)
  
	  payload := make([]byte, payloadLen)
	  io.ReadFull(conn, payload)
  
	  // Extract hidden data
	  secret, _ := stego.DecodeFromImage(payload)
	  fmt.Printf("Extracted secret: %s\n", secret)
	}
  }