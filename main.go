package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nsf/termbox-go"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a file path")
		os.Exit(1)
	}

	path := os.Args[1]
	var debug string = ""
	if len(os.Args) > 2 {
		debug = os.Args[2]
	}
	err := termbox.Init()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer termbox.Close()
	if debug != "" {
		f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(fmt.Errorf("error opening file: %v", err))
		}
		defer f.Close()
		log.SetOutput(f)
	}

	termbox.SetInputMode(termbox.InputEsc)

	sections := strings.Split(path, "/")
	fileName := sections[len(sections)-1]

	content := ReadFile(path)

	textEditor := InitTextEditor(content, fileName)
	textEditor.Draw()

mainLoop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {

		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break mainLoop
			case termbox.KeyCtrlS:
				WriteFile(path, textEditor.Text)
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
			case termbox.KeyCtrlR:
				textEditor.Redo()
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
