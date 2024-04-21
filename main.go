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

	textEditor := InitTextEditor("Voluptate nostrud aliqua cupidatat amet Lorem nulla laborum id dolore reprehenderit eu consectetur tempor aliquip. Aliquip non anim commodo nisi. Ut culpa aute ex deserunt id consequat aliqua amet labore proident ex. Ad ex aliqua in ullamco. Voluptate ullamco deserunt cupidatat minim veniam aute magna adipisicing occaecat duis. Ullamco adipisicing nisi occaecat qui.", "my-file.txt")
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
			case termbox.KeyCtrlQ:
				textEditor.ToggleSelectionMode()
			case termbox.KeyCtrlC:
				textEditor.CopySelection()
			case termbox.KeyCtrlX:
				textEditor.CutSelection()
			case termbox.KeyCtrlV:
				textEditor.Paste()
			case termbox.KeyCtrlZ:
				textEditor.Undo()
			case termbox.KeySpace:
				textEditor.InsertString(" ")
			case termbox.KeyEnter:
				textEditor.AddNewLine()
			case termbox.KeyBackspace2:
				textEditor.RemoveChar()
			default:
				if ev.Ch != 0 {
					textEditor.InsertString(string(ev.Ch))
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}

}
