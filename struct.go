package main

import (
	"image"
	"os"
)

// MessageProcessor handles messages and sends them into the painting system
type MessageProcessor struct {
	Chan chan PaintCommand
}

// ImageProcessor handles writing PaintCommands to an image file
type ImageProcessor struct {
	Chan    chan PaintCommand
	Image   image.RGBA
	OutFile *os.File
}

// PaintCommand encodes a chat-entered command to color a pixel
type PaintCommand struct {
	X     uint16
	Y     uint16
	Red   uint8
	Blue  uint8
	Green uint8
}

// Configuration holds configuration parsed when the program starts up
type Configuration struct {
	Username string `toml:"username"`
	APIToken string `toml:"apitoken"`
}
