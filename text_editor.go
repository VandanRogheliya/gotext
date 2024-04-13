package main

import (
	"fmt"
	"log"

	"github.com/nsf/termbox-go"
)

const EDITOR_START_X = 0

type TextEditor struct {
	Text                string
	FileName            string
	CursorXOffset       int
	CursorXOffsetActual int
	CursorYOffset       int
	cellGrid            [][]rune
	scrollX             int
	scrollY             int
}

func InitTextEditor(text string, fileName string) *TextEditor {
	return &TextEditor{
		Text:                text,
		CursorXOffset:       EDITOR_START_X,
		CursorYOffset:       0,
		CursorXOffsetActual: EDITOR_START_X,
		FileName:            fileName,
	}
}

func (te *TextEditor) drawChar(x, y int, c rune, fgColor termbox.Attribute, bgColor termbox.Attribute) {
	termbox.SetCell(x, y, c, fgColor, bgColor)
}

func (te *TextEditor) clearEditor() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}

func (te *TextEditor) height() int {
	return len(te.cellGrid)
}

func (te *TextEditor) width(row int) int {
	return len(te.cellGrid[row])
}

func (te *TextEditor) getTextIndex() int {
	totalChars := 0
	for i := 0; i < te.CursorYOffset; i++ {
		totalChars += len(te.cellGrid[i]) - EDITOR_START_X + 1 // +1 for '\n' char
	}
	totalChars += te.CursorXOffset - EDITOR_START_X
	textIndex := totalChars
	return textIndex
}

func (te *TextEditor) drawBottomInfoStrip() {
	infoText := fmt.Sprintf("%d,%d | %s", te.CursorXOffset, te.CursorYOffset, te.FileName)
	termboxWidth, termboxHeight := termbox.Size()

	for i := 0; i < termboxWidth; i++ {
		if i >= len(infoText) {
			te.drawChar(i, termboxHeight-1, ' ', termbox.ColorWhite, termbox.ColorWhite)
			continue
		}
		te.drawChar(i, termboxHeight-1, rune(infoText[i]), termbox.ColorDefault, termbox.ColorWhite)
	}
}

func (te *TextEditor) Draw() {
	te.clearEditor()
	_, h := termbox.Size()
	te.cellGrid = [][]rune{}
	te.cellGrid = append(te.cellGrid, []rune{})

	for _, c := range te.Text {
		lastRowIndex := len(te.cellGrid) - 1
		if c == '\n' {
			var row []rune
			te.cellGrid = append(te.cellGrid, row)
			if c == '\n' {
				continue
			}
		}
		lastRowIndex = len(te.cellGrid) - 1
		te.cellGrid[lastRowIndex] = append(te.cellGrid[lastRowIndex], c)
	}

	if te.height() > h {
		panic("Text will overflow. Height of terminal isnt enough")
	}

	for y, row := range te.cellGrid {
		for x, c := range row {
			te.drawChar(x, y, c, termbox.ColorWhite, termbox.ColorDefault)
		}
	}

	for i := te.height(); i <= h; i++ {
		te.drawChar(0, i, '~', termbox.ColorCyan, termbox.ColorDefault)
	}
	te.drawBottomInfoStrip()
	termbox.SetCursor(te.CursorXOffset, te.CursorYOffset)
	err := termbox.Flush()
	if err != nil {
		panic(err)
	}
}

func (te *TextEditor) MoveCursorTo(x, y int) {
	// TODO: test limits: edge of editor
	if EDITOR_START_X > x {
		return
	}

	if te.height() <= y {
		y = te.height() - 1
		x = te.width(y)
		te.CursorXOffset, te.CursorYOffset = x, y
		return
	}

	if 0 > y {
		y = 0
		x = EDITOR_START_X
		te.CursorXOffset, te.CursorYOffset = x, y
		return
	}

	if y != te.CursorYOffset {
		x = te.CursorXOffsetActual
		if te.width(y) < x {
			x = te.width(y)
		}
	} else if x != te.CursorXOffset {
		if te.width(y) < x {
			return
		}
		te.CursorXOffsetActual = x
	}

	te.CursorXOffset, te.CursorYOffset = x, y
}

func (te *TextEditor) MoveCursorRight() {
	te.MoveCursorTo(te.CursorXOffset+1, te.CursorYOffset)
	te.Draw()
}

func (te *TextEditor) MoveCursorLeft() {
	te.MoveCursorTo(te.CursorXOffset-1, te.CursorYOffset)
	te.Draw()
}

func (te *TextEditor) MoveCursorUp() {
	te.MoveCursorTo(te.CursorXOffset, te.CursorYOffset-1)
	te.Draw()
}

func (te *TextEditor) MoveCursorDown() {
	te.MoveCursorTo(te.CursorXOffset, te.CursorYOffset+1)
	te.Draw()
}

func (te *TextEditor) MoveCursorToEndOfTheLine() {
	x := te.width(te.CursorYOffset)
	te.MoveCursorTo(x, te.CursorYOffset)
	te.Draw()
}

func (te *TextEditor) MoveCursorToBeginningOfTheLine() {
	te.MoveCursorTo(EDITOR_START_X, te.CursorYOffset)
	te.Draw()
}

func (te *TextEditor) InsertChar(c rune) {
	textIndex := te.getTextIndex()
	runeSlice := []rune(te.Text)
	runeSlice = append(runeSlice[:textIndex], append([]rune{c}, runeSlice[textIndex:]...)...)
	te.Text = string(runeSlice)
	te.Draw()
	te.MoveCursorRight()
}

func (te *TextEditor) AddNewLine() {
	textIndex := te.getTextIndex()
	runeSlice := []rune(te.Text)
	runeSlice = append(runeSlice[:textIndex], append([]rune{'\n'}, runeSlice[textIndex:]...)...)
	te.Text = string(runeSlice)
	te.MoveCursorToBeginningOfTheLine()
	te.MoveCursorDown()
	te.Draw()
}

func (te *TextEditor) RemoveChar() {
	if te.CursorXOffset == EDITOR_START_X && te.CursorYOffset == 0 {
		return
	}
	textIndex := te.getTextIndex() - 1
	log.Printf("Deleting: %c", (te.Text[textIndex]))
	runeSlice := []rune(te.Text)
	runeSlice = append(runeSlice[:textIndex], runeSlice[textIndex+1:]...)
	if te.CursorXOffset == EDITOR_START_X {
		te.MoveCursorUp()
		te.MoveCursorToEndOfTheLine()
	} else {
		te.MoveCursorLeft()
	}
	te.Text = string(runeSlice)
	te.Draw()
}
