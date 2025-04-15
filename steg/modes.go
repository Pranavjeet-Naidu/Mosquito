package steg

// StegMode represents different steganography algorithms
type StegMode int

const (
    // LSB1 uses the least significant bit of one channel (R)
    LSB1 StegMode = iota
    // LSB3 uses the least significant bit of all three channels (RGB)
    LSB3
    // LSB4 uses the 2 least significant bits of R and G channels
    LSB4
    // LSB8 uses the least significant bit of all pixel bytes (RGBA, 2-bits each)
    LSB8
)

// ModeNames provides human-readable names for steganography modes
var ModeNames = map[StegMode]string{
    LSB1: "LSB-1 (R channel only)",
    LSB3: "LSB-3 (RGB channels)",
    LSB4: "LSB-4 (2-bits in R & G)",
    LSB8: "LSB-8 (all channels, 2-bits each)",
}

// CapacityFactor returns the number of bits per pixel for each mode
func (m StegMode) CapacityFactor() int {
    switch m {
    case LSB1:
        return 1
    case LSB3:
        return 3
    case LSB4:
        return 4
    case LSB8:
        return 8
    default:
        return 1
    }
}

// GetAvailableModes returns a slice of all available steganography modes
func GetAvailableModes() []StegMode {
    return []StegMode{LSB1, LSB3, LSB4, LSB8}
}