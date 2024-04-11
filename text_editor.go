package main

import (
	"strconv"

	"github.com/nsf/termbox-go"
)

const EDITOR_START_X = 2

type TextEditor struct {
	Text                string
	CursorXOffset       int
	CursorXOffsetActual int
	CursorYOffset       int
	cellGrid            [][]rune
}

func InitTextEditor(text string) *TextEditor {
	return &TextEditor{Text: text, CursorXOffset: 2, CursorYOffset: 0, CursorXOffsetActual: 2}
}

func (te *TextEditor) drawChar(x, y int, c rune) {
	termbox.SetCell(x, y, c, termbox.ColorRed, termbox.ColorDefault)
}

func (te *TextEditor) height() int {
	return len(te.cellGrid)
}

func (te *TextEditor) width(row int) int {
	return len(te.cellGrid[row])
}

func (te *TextEditor) Draw() {
	w, h := termbox.Size()

	var newCellGrid [][]rune

	newCellGrid = append(newCellGrid, []rune{'1', ' '})

	for _, c := range te.Text {
		lastRowIndex := len(newCellGrid) - 1
		if len(newCellGrid[lastRowIndex]) == w || c == '\n' {
			var row []rune
			row = append(row, rune(strconv.Itoa(len(newCellGrid) + 1)[0]))
			row = append(row, ' ')
			newCellGrid = append(newCellGrid, row)
			if c == '\n' {
				continue
			}
		}
		lastRowIndex = len(newCellGrid) - 1
		newCellGrid[lastRowIndex] = append(newCellGrid[lastRowIndex], c)
	}

	if len(newCellGrid) > h {
		panic("Text will overflow. Height of terminal isnt enough")
	}

	for y, row := range newCellGrid {
		for x, c := range row {
			te.drawChar(x, y, c)
		}
	}

	te.cellGrid = newCellGrid

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
}

func (te *TextEditor) MoveCursorLeft() {
	te.MoveCursorTo(te.CursorXOffset-1, te.CursorYOffset)
}

func (te *TextEditor) MoveCursorUp() {
	te.MoveCursorTo(te.CursorXOffset, te.CursorYOffset-1)
}

func (te *TextEditor) MoveCursorDown() {
	te.MoveCursorTo(te.CursorXOffset, te.CursorYOffset+1)
}

func (te *TextEditor) MoveCursorToEndOfTheLine() {
	x := te.width(te.CursorYOffset)
	te.MoveCursorTo(x, te.CursorYOffset)
}

func (te *TextEditor) MoveCursorToBeginningOfTheLine() {
	te.MoveCursorTo(EDITOR_START_X, te.CursorYOffset)
}
