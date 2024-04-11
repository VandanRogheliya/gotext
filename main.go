package main

import (
	"fmt"
	"os"

	"github.com/nsf/termbox-go"
)

func main() {
	err := termbox.Init()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc)

	textEditor := InitTextEditor("Hello world! Hello world! Hello world! Hello world! Hello world! Hello world! Hello world!\nHello world!\n Hello world! Hello world! Hello world! Hello world! Hello world! Hello world! Hello world! Hello world! Hello world! Hello world! Hello world! Hello world!")
	textEditor.Draw()

mainLoop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {

		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break mainLoop
			case termbox.KeyArrowRight:
				textEditor.MoveCursorRight()
				textEditor.Draw()
			case termbox.KeyArrowLeft:
				textEditor.MoveCursorLeft()
				textEditor.Draw()
			case termbox.KeyArrowUp:
				textEditor.MoveCursorUp()
				textEditor.Draw()
			case termbox.KeyArrowDown:
				textEditor.MoveCursorDown()
				textEditor.Draw()
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}

	// time.Sleep(time.Second)
}
