package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/jroimartin/gocui"
)

var (
	conn       net.Conn
	g          *gocui.Gui
	history    []string
	historyMux sync.Mutex
)

func main() {
	var err error
	conn, err = net.Dial("tcp", "localhost:8989")
	if err != nil {
		log.Fatalf("Failed to connect to server: %s", err.Error())
	}
	defer conn.Close()

	// Initialize gocui
	g, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatalf("Failed to initialize gocui: %v", err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)
	if err := keybindings(g); err != nil {
		log.Fatalf("Failed to set keybindings: %v", err)
	}

	// Start receiving messages in a goroutine
	go receiveMessages(conn)

	// Main loop for the UI
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatalf("Error in main loop: %v", err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// Chat view
	if v, err := g.SetView("chat", 0, 0, maxX-1, maxY-5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Chat"
		v.Autoscroll = true
		v.Wrap = true
		historyMux.Lock()
		for _, msg := range history {
			fmt.Fprintln(v, msg)
		}
		historyMux.Unlock()
	}

	// Input view
	if v, err := g.SetView("input", 0, maxY-4, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Input"
		v.Editable = true
		v.Wrap = true
		if _, err := g.SetCurrentView("input"); err != nil {
			return err
		}
	}

	return nil
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, sendMessage); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	return nil
}

func sendMessage(g *gocui.Gui, v *gocui.View) error {
	inputView, err := g.View("input")
	if err != nil {
		return err
	}
	msg := strings.TrimSpace(inputView.Buffer())
	if msg == "" {
		return nil // Ignore empty messages
	}

	_, err = fmt.Fprintf(conn, "%s\n", msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}

	inputView.Clear()
	inputView.SetCursor(0, 0)
	return nil
}

func receiveMessages(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()
		historyMux.Lock()
		history = append(history, msg)
		historyMux.Unlock()

		// Update UI
		g.Update(func(g *gocui.Gui) error {
			chatView, err := g.View("chat")
			if err != nil {
				return err
			}
			chatView.Clear()
			historyMux.Lock()
			for _, msg := range history {
				fmt.Fprintln(chatView, msg)
			}
			historyMux.Unlock()
			return nil
		})
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from server: %v", err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
