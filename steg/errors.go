package steg

import "errors"

// Error types for steganography operations
var (
    ErrInvalidHeader     = errors.New("invalid steganography header")
    ErrInvalidMagic      = errors.New("invalid magic byte, not a Mosquito steganographic image")
    ErrUnsupportedMode   = errors.New("unsupported steganography mode")
    ErrImageTooSmall     = errors.New("image too small to encode payload")
    ErrMessageCorrupted  = errors.New("message data corrupted or truncated")
    ErrEncryptionFailed  = errors.New("encryption failed")
    ErrDecryptionFailed  = errors.New("decryption failed, invalid key or corrupted data")
    ErrInvalidImage      = errors.New("invalid or unsupported image format")
    ErrInvalidKey        = errors.New("invalid encryption key")
)