# Mosquito

Mosquito is a powerful steganography tool with MQTT support, allowing you to hide and transmit secret data within images.

## What it can do 

- **Steganography**
  - Hide text messages in images
  - Hide one image inside another image
  - Multiple encoding algorithms (LSB1, LSB3, LSB4, LSB8)
  - Verification of image capacity before encoding
  - Support for multiple image formats (PNG, JPEG, GIF, BMP, TIFF, WebP)

- **Security**
  - AES-256-GCM encryption for protected content
  - Password-based protection
  - Image difference analysis to assess stealth

- **MQTT **
  - Send steganographic images via MQTT
  - Receive and automatically save incoming steganographic images
  - Secure communication channels


## Installation

### Prerequisites

- Go 1.18 or later

### From Source

```bash
# Clone the repository
git clone https://github.com/Pranavjeet-Naidu/Mosquito.git
cd Mosquito

# Build the application
go build

# Install globally (optional)
go install
```

### Using Go

```bash
go install github.com/Pranavjeet-Naidu/Mosquito@latest
```

## Quick Start

```bash
# Hide a text message in an image
./Mosquito hideMsg -i cover.png -o hidden.png -m "This is a secret message"

# Extract a hidden message
./Mosquito extract -i hidden.png -t

# Check an image's steganography capacity
./Mosquito info -i image.png

# Send a steganographic image via MQTT
./Mosquito mqttSend -b tcp://broker.example.com:1883 -t stego/channel -i hidden.png
```



## License

This project is licensed under the MIT License - see the LICENSE file for details.

