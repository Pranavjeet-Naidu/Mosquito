package steg

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/sha256"
    "image"
    "image/color"
    "io"
)

// Capacity checks if the image can store the payload using the given mode
func Capacity(img image.Image, dataSize int, mode StegMode) (bool, int, int) {
    bounds := img.Bounds()
    pixelCount := bounds.Dx() * bounds.Dy()
    
    // Calculate bits per pixel for the mode
    bitsPerPixel := mode.CapacityFactor()
    
    // Total capacity in bits
    totalBits := pixelCount * bitsPerPixel
    
    // Header size in bits (8 bytes * 8 bits)
    headerBits := 8 * 8
    
    // Available bits for payload
    availableBits := totalBits - headerBits
    
    // Required bits for payload (dataSize in bytes * 8 bits)
    requiredBits := dataSize * 8
    
    return availableBits >= requiredBits, availableBits / 8, requiredBits / 8
}

// EncodeMessage embeds a message into an image
func EncodeMessage(img image.Image, msg []byte, mode StegMode) (image.Image, error) {
    return EncodeMessageWithPassword(img, msg, "", mode, false)
}

func IsStegImage(img image.Image) bool {
    // Try all modes for header detection
    for _, mode := range []StegMode{LSB1, LSB3, LSB4, LSB8} {
        var headerData []byte
        switch mode {
        case LSB1:
            headerData = decodeLSB1(img, 8, 0)
        case LSB3:
            headerData = decodeLSB3(img, 8, 0)
        case LSB4:
            headerData = decodeLSB4(img, 8, 0)
        case LSB8:
            headerData = decodeLSB8(img, 8, 0)
        }
        
        if len(headerData) >= 8 && headerData[0] == MagicByte {
            return true
        }
    }
    return false
}

// EncodeMessageWithPassword embeds an encrypted message into an image
func EncodeMessageWithPassword(img image.Image, msg []byte, password string, mode StegMode, isImage bool) (image.Image, error) {
    bounds := img.Bounds()
    
    // Check if the image has enough capacity
    hasCapacity, _, _ := Capacity(img, len(msg), mode)
    if !hasCapacity {
        return nil, ErrImageTooSmall
    }
    
    // Create the output image
    out := image.NewRGBA(bounds)
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            r, g, b, a := img.At(x, y).RGBA()
            out.Set(x, y, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)})
        }
    }
    
    // Set up the header
    header := Header{
        Magic:     MagicByte,
        Version:   Version,
        Mode:      mode,
        Flags:     0,
        PayloadLen: uint32(len(msg)),
    }
    
    // Set image flag if payload is an image
    if isImage {
        header.Flags |= FlagImage
    }
    
    // Encrypt the payload if a password is provided
    finalMsg := msg
    if password != "" {
        encryptedMsg, err := encrypt(msg, password)
        if err != nil {
            return nil, err
        }
        finalMsg = encryptedMsg
        header.Flags |= FlagEncrypted
        header.PayloadLen = uint32(len(encryptedMsg))
    }
    
    // Encode the header and payload
    headerData := MarshalHeader(header)
    data := append(headerData, finalMsg...)
    
    // Encode the data using the specified mode
    switch mode {
    case LSB1:
        encodeLSB1(out, data)
    case LSB3:
        encodeLSB3(out, data)
    case LSB4:
        encodeLSB4(out, data)
    case LSB8:
        encodeLSB8(out, data)
    default:
        return nil, ErrUnsupportedMode
    }
    
    return out, nil
}



// DecodeMessage extracts a message from an image
func DecodeMessage(img image.Image) ([]byte, error) {
    return DecodeMessageWithPassword(img, "")
}

// DecodeMessageWithPassword extracts and decrypts a message from an image
// DecodeMessageWithPassword extracts and decrypts a message from an image
func DecodeMessageWithPassword(img image.Image, password string) ([]byte, error) {
    // Try to extract the header using all possible modes
    var header Header
    // var err error
    var headerFound bool
    
    for _, mode := range []StegMode{LSB1, LSB3, LSB4, LSB8} {
        var headerData []byte
        switch mode {
        case LSB1:
            headerData = decodeLSB1(img, 8, 0)
        case LSB3:
            headerData = decodeLSB3(img, 8, 0)
        case LSB4:
            headerData = decodeLSB4(img, 8, 0)
        case LSB8:
            headerData = decodeLSB8(img, 8, 0)
        }
        
        if len(headerData) >= 8 && headerData[0] == MagicByte {
            tmpHeader, tmpErr := UnmarshalHeader(headerData)
            if tmpErr == nil {
                header = tmpHeader
                headerFound = true
                break
            }
        }
    }
    
    if !headerFound {
        return nil, ErrInvalidHeader
    }
    
    // Extract data based on the mode in the header
    var data []byte
    switch header.Mode {
    case LSB1:
        data = decodeLSB1(img, int(header.PayloadLen), 8) // Skip 8 bytes of header
    case LSB3:
        data = decodeLSB3(img, int(header.PayloadLen), 8)
    case LSB4:
        data = decodeLSB4(img, int(header.PayloadLen), 8)
    case LSB8:
        data = decodeLSB8(img, int(header.PayloadLen), 8)
    default:
        return nil, ErrUnsupportedMode
    }
    
    // Decrypt the data if it's encrypted
    if header.IsEncrypted() {
        if password == "" {
            return nil, ErrDecryptionFailed
        }
        
        decrypted, err := decrypt(data, password)
        if err != nil {
            return nil, err
        }
        return decrypted, nil
    }
    
    return data, nil
}
// GetImageInfo extracts information about a steganographic image
func GetImageInfo(img image.Image) (Header, error) {
    // Try all modes for header detection
    for _, mode := range []StegMode{LSB1, LSB3, LSB4, LSB8} {
        var headerData []byte
        switch mode {
        case LSB1:
            headerData = decodeLSB1(img, 8, 0)
        case LSB3:
            headerData = decodeLSB3(img, 8, 0)
        case LSB4:
            headerData = decodeLSB4(img, 8, 0)
        case LSB8:
            headerData = decodeLSB8(img, 8, 0)
        }
        
        if len(headerData) >= 8 && headerData[0] == MagicByte {
            header, err := UnmarshalHeader(headerData)
            if err == nil {
                return header, nil
            }
        }
    }
    
    return Header{}, ErrInvalidHeader
}
// ========================= LSB Encoding Functions =========================

// encodeLSB1 embeds data using LSB of the red channel only
func encodeLSB1(img *image.RGBA, data []byte) {
    bounds := img.Bounds()
    bitIndex := 0
    totalBits := len(data) * 8
    
    for y := bounds.Min.Y; y < bounds.Max.Y && bitIndex < totalBits; y++ {
        for x := bounds.Min.X; x < bounds.Max.X && bitIndex < totalBits; x++ {
            byteIndex := bitIndex / 8
            bitPos := 7 - (bitIndex % 8)
            bit := (data[byteIndex] >> bitPos) & 1
            
            idx := (y-bounds.Min.Y)*img.Stride + (x-bounds.Min.X)*4
            
            // Modify only the least significant bit of the red channel
            img.Pix[idx] = (img.Pix[idx] & 0xFE) | bit
            
            bitIndex++
        }
    }
}

// encodeLSB3 embeds data using LSB of RGB channels
func encodeLSB3(img *image.RGBA, data []byte) {
    bounds := img.Bounds()
    bitIndex := 0
    totalBits := len(data) * 8
    
    for y := bounds.Min.Y; y < bounds.Max.Y && bitIndex < totalBits; y++ {
        for x := bounds.Min.X; x < bounds.Max.X && bitIndex < totalBits; x++ {
            idx := (y-bounds.Min.Y)*img.Stride + (x-bounds.Min.X)*4
            
            // Modify R, G, B channels (3 bits per pixel)
            for c := 0; c < 3 && bitIndex < totalBits; c++ {
                byteIndex := bitIndex / 8
                bitPos := 7 - (bitIndex % 8)
                bit := (data[byteIndex] >> bitPos) & 1
                
                img.Pix[idx+c] = (img.Pix[idx+c] & 0xFE) | bit
                
                bitIndex++
            }
        }
    }
}

// encodeLSB4 embeds data using 2 LSBs of R and G channels
func encodeLSB4(img *image.RGBA, data []byte) {
    bounds := img.Bounds()
    bitIndex := 0
    totalBits := len(data) * 8
    
    for y := bounds.Min.Y; y < bounds.Max.Y && bitIndex < totalBits; y++ {
        for x := bounds.Min.X; x < bounds.Max.X && bitIndex < totalBits; x++ {
            idx := (y-bounds.Min.Y)*img.Stride + (x-bounds.Min.X)*4
            
            // Modify 2 LSB of R (2 bits)
            for i := 0; i < 2 && bitIndex < totalBits; i++ {
                byteIndex := bitIndex / 8
                bitPos := 7 - (bitIndex % 8)
                bit := (data[byteIndex] >> bitPos) & 1
                
                img.Pix[idx] = (img.Pix[idx] & (0xFF - (1 << i))) | (bit << i)
                
                bitIndex++
            }
            
            // Modify 2 LSB of G (2 bits)
            for i := 0; i < 2 && bitIndex < totalBits; i++ {
                byteIndex := bitIndex / 8
                bitPos := 7 - (bitIndex % 8)
                bit := (data[byteIndex] >> bitPos) & 1
                
                img.Pix[idx+1] = (img.Pix[idx+1] & (0xFF - (1 << i))) | (bit << i)
                
                bitIndex++
            }
        }
    }
}

// encodeLSB8 embeds data using 2 LSBs of all RGBA channels
func encodeLSB8(img *image.RGBA, data []byte) {
    bounds := img.Bounds()
    bitIndex := 0
    totalBits := len(data) * 8
    
    for y := bounds.Min.Y; y < bounds.Max.Y && bitIndex < totalBits; y++ {
        for x := bounds.Min.X; x < bounds.Max.X && bitIndex < totalBits; x++ {
            idx := (y-bounds.Min.Y)*img.Stride + (x-bounds.Min.X)*4
            
            // Modify 2 LSB of R, G, B, A (8 bits per pixel)
            for c := 0; c < 4 && bitIndex < totalBits; c++ {
                for i := 0; i < 2 && bitIndex < totalBits; i++ {
                    byteIndex := bitIndex / 8
                    bitPos := 7 - (bitIndex % 8)
                    bit := (data[byteIndex] >> bitPos) & 1
                    
                    img.Pix[idx+c] = (img.Pix[idx+c] & (0xFF - (1 << i))) | (bit << i)
                    
                    bitIndex++
                }
            }
        }
    }
}

// ========================= LSB Decoding Functions =========================

// decodeLSB1 extracts data from LSB of the red channel
func decodeLSB1(img image.Image, dataSize int, offset int) []byte {
    bounds := img.Bounds()
    output := make([]byte, dataSize)
    bitIndex := 0
    totalBits := dataSize * 8
    
    // Skip offset bits
    skipBits := offset * 8
    
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            if bitIndex < skipBits {
                bitIndex++
                continue
            }
            
            if bitIndex-skipBits >= totalBits {
                return output
            }
            
            r, _, _, _ := img.At(x, y).RGBA()
            // Fix: add shift to extract the correct bit
            bit := byte((r >> 8) & 1)
            
            outIndex := (bitIndex - skipBits) / 8
            bitPos := 7 - ((bitIndex - skipBits) % 8)
            
            output[outIndex] |= bit << bitPos
            
            bitIndex++
        }
    }
    
    return output
}

// decodeLSB3 extracts data from LSB of RGB channels
func decodeLSB3(img image.Image, dataSize int, offset int) []byte {
    bounds := img.Bounds()
    output := make([]byte, dataSize)
    bitIndex := 0
    totalBits := dataSize * 8
    
    // Skip offset bits
    skipBits := offset * 8
    totalSkipAndData := skipBits + totalBits
    
    for y := bounds.Min.Y; y < bounds.Max.Y && bitIndex < totalSkipAndData; y++ {
        for x := bounds.Min.X; x < bounds.Max.X && bitIndex < totalSkipAndData; x++ {
            r, g, b, _ := img.At(x, y).RGBA()
            channels := []uint32{r, g, b}
            
            for c := 0; c < 3 && bitIndex < totalSkipAndData; c++ {
                // Skip bits according to offset
                if bitIndex < skipBits {
                    bitIndex++
                    continue
                }
                
                // Extract bit from the channel's LSB
                bit := byte((channels[c] >> 8) & 1)
                
                // Place the bit in the output byte
                outIndex := (bitIndex - skipBits) / 8
                bitPos := 7 - ((bitIndex - skipBits) % 8)
                
                output[outIndex] |= bit << bitPos
                
                bitIndex++
            }
        }
    }
    
    return output
}

// decodeLSB4 extracts data from 2 LSBs of R and G channels
func decodeLSB4(img image.Image, dataSize int, offset int) []byte {
    bounds := img.Bounds()
    output := make([]byte, dataSize)
    bitIndex := 0
    totalBits := dataSize * 8
    
    // Skip offset bits
    skipBits := offset * 8
    
    for y := bounds.Min.Y; y < bounds.Max.Y && (bitIndex-skipBits) < totalBits; y++ {
        for x := bounds.Min.X; x < bounds.Max.X && (bitIndex-skipBits) < totalBits; x++ {
            r, g, _, _ := img.At(x, y).RGBA()
            
            // Extract 2 bits from R
            for i := 0; i < 2 && (bitIndex-skipBits) < totalBits; i++ {
                if bitIndex < skipBits {
                    bitIndex++
                    continue
                }
                
                // Fix: Add shift for proper bit extraction
                bit := byte((r >> (8 + i)) & 1)
                
                outIndex := (bitIndex - skipBits) / 8
                bitPos := 7 - ((bitIndex - skipBits) % 8)
                
                output[outIndex] |= bit << bitPos
                
                bitIndex++
            }
            
            // Extract 2 bits from G
            for i := 0; i < 2 && (bitIndex-skipBits) < totalBits; i++ {
                if bitIndex < skipBits {
                    bitIndex++
                    continue
                }
                
                // Fix: Add shift for proper bit extraction
                bit := byte((g >> (8 + i)) & 1)
                
                outIndex := (bitIndex - skipBits) / 8
                bitPos := 7 - ((bitIndex - skipBits) % 8)
                
                output[outIndex] |= bit << bitPos
                
                bitIndex++
            }
        }
    }
    
    return output
}

// decodeLSB8 extracts data from 2 LSBs of all RGBA channels
func decodeLSB8(img image.Image, dataSize int, offset int) []byte {
    bounds := img.Bounds()
    output := make([]byte, dataSize)
    bitIndex := 0
    totalBits := dataSize * 8
    
    // Skip offset bits
    skipBits := offset * 8
    
    for y := bounds.Min.Y; y < bounds.Max.Y && (bitIndex-skipBits) < totalBits; y++ {
        for x := bounds.Min.X; x < bounds.Max.X && (bitIndex-skipBits) < totalBits; x++ {
            r, g, b, a := img.At(x, y).RGBA()
            channels := []uint32{r, g, b, a}
            
            for c := 0; c < 4 && (bitIndex-skipBits) < totalBits; c++ {
                for i := 0; i < 2 && (bitIndex-skipBits) < totalBits; i++ {
                    if bitIndex < skipBits {
                        bitIndex++
                        continue
                    }
                    
                    // Fix: Add shift for proper bit extraction
                    bit := byte((channels[c] >> (8 + i)) & 1)
                    
                    outIndex := (bitIndex - skipBits) / 8
                    bitPos := 7 - ((bitIndex - skipBits) % 8)
                    
                    output[outIndex] |= bit << bitPos
                    
                    bitIndex++
                }
            }
        }
    }
    
    return output
}

// ========================= Encryption Functions =========================

// encrypt data with AES-256-GCM
func encrypt(data []byte, password string) ([]byte, error) {
    // Create a key from the password
    key := sha256.Sum256([]byte(password))
    
    block, err := aes.NewCipher(key[:])
    if err != nil {
        return nil, err
    }
    
    // Create GCM mode
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    // Create a nonce
    nonce := make([]byte, gcm.NonceSize())
    if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }
    
    // Encrypt
    ciphertext := gcm.Seal(nonce, nonce, data, nil)
    
    return ciphertext, nil
}

// decrypt data with AES-256-GCM
func decrypt(data []byte, password string) ([]byte, error) {
    // Create a key from the password
    key := sha256.Sum256([]byte(password))
    
    block, err := aes.NewCipher(key[:])
    if err != nil {
        return nil, err
    }
    
    // Create GCM mode
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    // Extract nonce
    nonceSize := gcm.NonceSize()
    if len(data) < nonceSize {
        return nil, ErrDecryptionFailed
    }
    
    nonce, ciphertext := data[:nonceSize], data[nonceSize:]
    
    // Decrypt
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return nil, ErrDecryptionFailed
    }
    
    return plaintext, nil
}

// Public versions of the encoding/decoding functions

// EncodeLSB1 is the public version of encodeLSB1
func EncodeLSB1(img *image.RGBA, data []byte) {
    encodeLSB1(img, data)
}

// EncodeLSB3 is the public version of encodeLSB3
func EncodeLSB3(img *image.RGBA, data []byte) {
    encodeLSB3(img, data)
}

// EncodeLSB4 is the public version of encodeLSB4
func EncodeLSB4(img *image.RGBA, data []byte) {
    encodeLSB4(img, data)
}

// EncodeLSB8 is the public version of encodeLSB8
func EncodeLSB8(img *image.RGBA, data []byte) {
    encodeLSB8(img, data)
}

// DecodeLSB1 is the public version of decodeLSB1
func DecodeLSB1(img image.Image, dataSize int, offset int) []byte {
    return decodeLSB1(img, dataSize, offset)
}

// DecodeLSB3 is the public version of decodeLSB3
func DecodeLSB3(img image.Image, dataSize int, offset int) []byte {
    return decodeLSB3(img, dataSize, offset)
}

// DecodeLSB4 is the public version of decodeLSB4
func DecodeLSB4(img image.Image, dataSize int, offset int) []byte {
    return decodeLSB4(img, dataSize, offset)
}

// DecodeLSB8 is the public version of decodeLSB8
func DecodeLSB8(img image.Image, dataSize int, offset int) []byte {
    return decodeLSB8(img, dataSize, offset)
}