package main

import (
	"errors"
	"fmt"
	"github.com/thoj/go-ircevent"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

/***
 * This file contains methods for MessageProcessor, called  when an IRC message is received.
 * These MessageProcessor methods parses the message to extract
 * an (X,Y) coordinate and a color code. Converts the message (if formatted
 * correctly) to a PaintCommand object and passes it on to the painting system
 ***/

// Handle takes an irc message, parses it, and sends valid commands to the
// painting subsystem
func (mp *MessageProcessor) Handle(event *irc.Event) {
	message := event.Message()
	message = mp.StripInput(message)

	paintCommand, err := mp.Parse(message)
	if err == nil {
		mp.Chan <- *paintCommand
	} else {
		fmt.Printf("ERROR: command %q encountered the following error: %q", message, err)
	}

}

// Parse will return a PaintCommand if the provided message can be parsed into
// one, else nil and an error are returned
func (mp *MessageProcessor) Parse(message string) (*PaintCommand, error) {

	// Regular expression to match a valid command
	re, _ := regexp.Compile("^\\(([0-9]?[0-9]?[0-9]?),([0-9]?[0-9]?[0-9]?),([a-fA-F0-9]{2})([a-fA-F0-9]{2})([a-fA-F0-9]{2})\\)$")
	commandFields := re.FindStringSubmatch(message)

	// If the regex doesnt match, return immediately
	if commandFields == nil {
		return nil, errors.New("Message does not match command regular expression")
	}

	// Parse out every field
	// Skip commandFields[0] because that will be the entire match
	x, xErr := strconv.ParseUint(commandFields[1], 10, 16)
	y, yErr := strconv.ParseUint(commandFields[2], 10, 16)
	r, rErr := strconv.ParseUint(commandFields[3], 16, 8)
	g, gErr := strconv.ParseUint(commandFields[4], 16, 8)
	b, bErr := strconv.ParseUint(commandFields[5], 16, 8)

	// See if any errors occurred during field parsing
	if xErr != nil || yErr != nil || rErr != nil || gErr != nil || bErr != nil {
		return nil, errors.New("Error converting field string to integer")
	}

	// Check if the parsed X and Y are out of the 800x600 boundaries
	if !(x < 800) {
		return nil, errors.New("X coordinate must be in [0,799]")
	} else if !(y < 600) {
		return nil, errors.New("Y coordinate must be in [0,599]")

	}

	// Create a PaintCommand and return it
	return &PaintCommand{
		X:     uint16(x),
		Y:     uint16(y),
		Red:   uint8(r),
		Blue:  uint8(b),
		Green: uint8(g),
	}, nil
}

// StripInput returns the string with only characters relevant to constructing
// a PaintCommand (parentheses, commas, and hex digits)
func (mp *MessageProcessor) StripInput(input string) string {

	validCharacters := [25]rune{
		'1', '2', '3', '4', '5', '6', '7', '8', '9', '0',
		'a', 'b', 'c', 'd', 'e', 'f', 'A', 'B', 'C', 'D', 'E', 'F',
		'(', ')', ','}
	return strings.Map(func(c rune) rune {
		if unicode.IsSpace(c) {
			return -1
		}
		for _, validChar := range validCharacters {
			if validChar == c {
				return c
			}
		}
		return -1
	}, input)
}
