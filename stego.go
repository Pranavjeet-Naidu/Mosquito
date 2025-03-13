package stego

import (
  "bytes"
  "image"
  "image/png"
  "encoding/base64"
  "strings"
)

import (
  
)

// Hide data in an image using LSB steganography
func EncodeInImage(imgBytes []byte, secret string) ([]byte, error) {
  img, _, err := image.Decode(bytes.NewReader(imgBytes))
  if err != nil {
    return nil, err
  }

  bounds := img.Bounds()
  rgbaImg := image.NewNRGBA(bounds)
  for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
    for x := bounds.Min.X; x < bounds.Max.X; x++ {
      rgbaImg.Set(x, y, img.At(x, y))
    }
  }

  secretBits := strToBits(secret + "\x00") // Null terminator
  bitIdx := 0

  // Embed in LSB of red channel
  for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
    for x := bounds.Min.X; x < bounds.Max.X; x++ {
      if bitIdx >= len(secretBits) {
        break
      }

      c := rgbaImg.NRGBAAt(x, y)
      c.R = embedByte(c.R, secretBits[bitIdx:bitIdx+1])
      bitIdx++
      rgbaImg.SetNRGBA(x, y, c)
    }
  }

  buf := new(bytes.Buffer)
  if err := png.Encode(buf, rgbaImg); err != nil {
    return nil, err
  }
  return buf.Bytes(), nil
}

// Extract data from image
func DecodeFromImage(imgBytes []byte) (string, error) {
  img, _, err := image.Decode(bytes.NewReader(imgBytes))
  if err != nil {
    return "", err
  }

  var bits []byte
  bounds := img.Bounds()

  // Extract from LSB of red channel
  for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
    for x := bounds.Min.X; x < bounds.Max.X; x++ {
      r, _, _, _ := img.At(x, y).RGBA()
      bits = append(bits, byte(r&1))
    }
  }

  return bitsToStr(bits), nil
}

// Helper functions
func strToBits(s string) []byte {
  bits := make([]byte, 0, len(s)*8)
  for _, c := range []byte(s) {
    for i := 7; i >= 0; i-- {
      bits = append(bits, (c>>i)&1)
    }
  }
  return bits
}

func bitsToStr(bits []byte) string {
  var result []byte
  for i := 0; i < len(bits); i += 8 {
    if i+8 > len(bits) {
      break
    }
    var char byte
    for j := 0; j < 8; j++ {
      char |= bits[i+j] << (7 - j)
    }
    if char == 0 {
      break
    }
    result = append(result, char)
  }
  return string(result)
}

func embedByte(original, bits byte) byte {
  return (original &^ 1) | (bits & 1)
}


// EncodeInTopic hides a secret message in a topic string
func EncodeInTopic(topic, secret string) string {
    // Simple implementation: Add the secret as a base64-encoded suffix
    encodedSecret := base64.StdEncoding.EncodeToString([]byte(secret))
    return topic + "/." + encodedSecret
}

// DecodeFromTopic extracts a hidden message from a topic string
func DecodeFromTopic(topic string) (string, string) {
    parts := strings.Split(topic, "/.")
    if len(parts) <= 1 {
        return topic, "" // No hidden message
    }
    
    // Extract original topic and encoded part
    originalTopic := parts[0]
    encoded := parts[len(parts)-1]
    
    decoded, err := base64.StdEncoding.DecodeString(encoded)
    if err != nil {
        return originalTopic, "" // Invalid encoding
    }
    
    return originalTopic, string(decoded)
}