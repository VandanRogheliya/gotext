package main

import (
	"fmt"
	"log"
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
	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	termbox.SetInputMode(termbox.InputEsc)

	textEditor := InitTextEditor("Voluptate nostrud aliqua cupidatat amet\nLorem nulla laborum id dolore\nreprehenderit eu consectetur tempor aliquip.", "my-file.txt")
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
			case termbox.KeyArrowLeft:
				textEditor.MoveCursorLeft()
			case termbox.KeyArrowUp:
				textEditor.MoveCursorUp()
			case termbox.KeyArrowDown:
				textEditor.MoveCursorDown()
			case termbox.KeyCtrlL:
				textEditor.MoveCursorToEndOfTheLine()
			case termbox.KeyCtrlH:
				textEditor.MoveCursorToBeginningOfTheLine()
			case termbox.KeySpace:
				textEditor.InsertChar(' ')
			case termbox.KeyEnter:
				textEditor.AddNewLine()
			case termbox.KeyBackspace2:
				textEditor.RemoveChar()
			default:
				if ev.Ch != 0 {
					textEditor.InsertChar(ev.Ch)
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}

}
