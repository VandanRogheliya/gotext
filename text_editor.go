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
	color := termbox.ColorWhite
	if x < EDITOR_START_X {
		color = termbox.ColorGreen
	}
	if c == '~' {
		color = termbox.ColorCyan
	}
	termbox.SetCell(x, y, c, color, termbox.ColorDefault)
}

func (te *TextEditor) height() int {
	return len(te.cellGrid)
}

func (te *TextEditor) width(row int) int {
	return len(te.cellGrid[row])
}

func (te *TextEditor) Draw() {
	_, h := termbox.Size()
	te.cellGrid = [][]rune{}
	te.cellGrid = append(te.cellGrid, []rune{'1', ' '})

	for _, c := range te.Text {
		lastRowIndex := len(te.cellGrid) - 1
		if c == '\n' {
			var row []rune
			row = append(row, rune(strconv.Itoa(len(te.cellGrid) + 1)[0]))
			row = append(row, ' ')
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
			te.drawChar(x, y, c)
		}
	}

	for i := te.height(); i <= h; i++ {
		te.drawChar(0, i, '~')
	}

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
