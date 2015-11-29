package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/thoj/go-ircevent"
	"os"
	"sync"
)

func main() {

	// The one and only command line argument should be a path to the config file
	if len(os.Args) != 2 {
		fmt.Printf("Expecting a single argument pointing to a config file")
		os.Exit(1)
	}
	// Parse the config file
	var config Configuration
	if _, err := toml.DecodeFile(os.Args[1], &config); err != nil {
		fmt.Printf("Error: %s", err)
	}

	// Initialize IRC session
	ircSession := irc.IRC(config.Username, config.Username)
	ircSession.Password = config.APIToken
	ircSession.Connect("irc.twitch.tv:6667")
	ircSession.Join(fmt.Sprintf("#%s", config.Username))

	// Create an empty, all white image
	img := freshImage()

	// Initialize the image, name TODO
	f, err := os.Create("image.png")
	if err != nil {
		// this shouldnt happen
	}

	// Initial image write, blank canvas...
	writeImage(f, &img)

	commandChannel := make(chan PaintCommand)
	imageProcessor := ImageProcessor{Chan: commandChannel,
		Image:   img,
		OutFile: f,
	}
	messageProcessor := MessageProcessor{commandChannel}

	ircSession.AddCallback("PRIVMSG", messageProcessor.Handle)

	var wg sync.WaitGroup

	wg.Add(1)
	go imageProcessor.Run()
	wg.Wait()

}
