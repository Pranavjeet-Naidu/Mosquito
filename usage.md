# Mosquito Usage Guide

## Table of Contents
- [General Usage](#general-usage)
- [Hiding Messages](#hiding-messages)
- [Hiding Images](#hiding-images)
- [Extracting Hidden Data](#extracting-hidden-data)
- [Image Analysis](#image-analysis)
- [MQTT Communication](#mqtt-communication)
- [Example Workflows](#example-workflows)


## General Usage

Mosquito uses a command-based structure:

```bash
mosquito [command] [flags]
```

To see all available commands:

```bash
mosquito --help
```

Each command has its own help information:

```bash
mosquito [command] --help
```

## Hiding Messages

### Basic Text Hiding

```bash
mosquito hideMsg -i cover.png -o stego.png -m "This is a secret message"
```

### Hide Text from a File

```bash
mosquito hideMsg -i cover.png -o stego.png -f secret.txt
```

### Using Different Steganography Modes

Mosquito supports four steganography algorithms with increasing capacity:

```bash
# LSB1 - Uses only the red channel (default, most stealthy)
mosquito hideMsg -i cover.png -o stego.png -m "Secret message" -M 0

# LSB3 - Uses RGB channels for higher capacity
mosquito hideMsg -i cover.png -o stego.png -m "Secret message" -M 1

# LSB4 - Uses 2-bits in R&G channels for even higher capacity
mosquito hideMsg -i cover.png -o stego.png -m "Larger secret message" -M 2

# LSB8 - Uses all channels with 2-bits each for maximum capacity
mosquito hideMsg -i cover.png -o stego.png -f largedatafile.txt -M 3
```

### With Encryption

```bash
mosquito hideMsg -i cover.png -o stego.png -m "Encrypted message" -p "mypassword"
```

## Hiding Images

### Basic Image Hiding

```bash
mosquito hideImg -i cover.png -s secret.jpg -o stego.png
```

### With Different Steganography Modes

```bash
mosquito hideImg -i cover.png -s secret.jpg -o stego.png -M 3
```

### With Encryption

```bash
mosquito hideImg -i cover.png -s secret.jpg -o stego.png -p "mypassword"
```

## Extracting Hidden Data

### Basic Extraction

```bash
mosquito extract -i stego.png -o extracted.bin
```

### Display Text Messages Directly

```bash
mosquito extract -i stego.png -t
```

### Extract with Decryption

```bash
mosquito extract -i stego.png -o extracted.jpg -p "mypassword"
```

### View Steganography Information

```bash
mosquito extract -i stego.png --info
```

## Image Analysis

Use the info command to check an image's capacity for steganography:

```bash
mosquito info -i image.png
```

This will display:
- Image dimensions and pixel count
- Whether the image contains hidden data
- Maximum payload sizes for different steganography modes
- Recommendations for which mode to use

## MQTT Communication

### Send a Steganographic Image

```bash
mosquito mqttSend -b tcp://broker.example.com:1883 -t stego/channel -i stego.png
```

### Receive Steganographic Images

```bash
# Start receiving and save to the 'received' directory
mosquito mqttRecv -b tcp://broker.example.com:1883 -t stego/channel -o ./received
```

This will:
1. Connect to the MQTT broker
2. Subscribe to the specified topic
3. Save any received images to the output directory with timestamps
4. Continue running until interrupted with Ctrl+C

## Example Workflows

### Secure Communication Workflow

```bash
# Sender side:
# 1. Create a steganographic image with hidden text
mosquito hideMsg -i photo.png -o hidden.png -m "Meet me at the park at 5pm" -p "secure123" -M 1

# 2. Send it via MQTT
mosquito mqttSend -b tcp://broker.example.com:1883 -t secret/channel123 -i hidden.png

# Receiver side:
# 1. Receive the image
mosquito mqttRecv -b tcp://broker.example.com:1883 -t secret/channel123 -o ./inbox

# 2. Extract the message
mosquito extract -i ./inbox/received-20250415-123010.png -p "secure123" -t
```

### Image-in-Image Workflow

```bash
# 1. Hide a small image inside a larger one
mosquito hideImg -i largecover.png -s smallsecret.jpg -o combined.png -M 3

# 2. Extract the hidden image
mosquito extract -i combined.png -o recovered.jpg
```

## Tips and Tricks :D


   - Use LSB1 (mode 0) for maximum stealth
   - Use higher modes only when needed for larger payloads


   - Add a password with the `-p` flag for all sensitive content


   - Check the image difference percentage ( displayed after hiding : Lower percentages mean less detectable changes )

