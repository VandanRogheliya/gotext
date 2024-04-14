package main

import (
	"fmt"
	"math"

	"github.com/nsf/termbox-go"
)

const EDITOR_START_X = 0
const EDITOR_START_Y = 0

type TextEditor struct {
	Text                string
	FileName            string
	CursorXOffset       int
	CursorXOffsetActual int
	CursorYOffset       int
	cellGrid            [][]rune
	scrollY             int
	newLines            map[int]struct{}
}

func InitTextEditor(text string, fileName string) *TextEditor {
	return &TextEditor{
		Text:                text,
		CursorXOffset:       EDITOR_START_X,
		CursorYOffset:       EDITOR_START_Y,
		CursorXOffsetActual: EDITOR_START_X,
		FileName:            fileName,
		scrollY:             EDITOR_START_Y,
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

		totalChars += len(te.cellGrid[i]) - EDITOR_START_X
		if _, ok := te.newLines[i]; ok {
			totalChars++
		}
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

func (te *TextEditor) getMaxVisibleContentHeight() int {
	_, h := termbox.Size()
	return h - 1 // -1 for info strip
}

func (te *TextEditor) getVisibleCellGrid() [][]rune {
	firstVisibleRow := te.scrollY
	lastVisibleRow := te.scrollY + te.getMaxVisibleContentHeight()
	lastVisibleRow = int(math.Min(float64(lastVisibleRow), float64(len(te.cellGrid))))
	return te.cellGrid[firstVisibleRow:lastVisibleRow]

}

func (te *TextEditor) setCursorAndScroll() {
	minY := te.scrollY
	maxY := te.scrollY + te.getMaxVisibleContentHeight()
	if minY > te.CursorYOffset {
		te.scrollY = te.CursorYOffset
	} else if maxY <= te.CursorYOffset {
		te.scrollY += (te.CursorYOffset - maxY + 1)
	}
	termbox.SetCursor(te.CursorXOffset, te.CursorYOffset-te.scrollY)
}

func (te *TextEditor) Draw() {
	te.clearEditor()
	w, h := termbox.Size()
	te.cellGrid = [][]rune{}
	te.newLines = make(map[int]struct{})
	te.cellGrid = append(te.cellGrid, []rune{})

	for _, c := range te.Text {
		lastRowIndex := len(te.cellGrid) - 1
		if c == '\n' || len(te.cellGrid[lastRowIndex]) >= w-1 {
			var row []rune
			te.cellGrid = append(te.cellGrid, row)
			if c == '\n' {
				te.newLines[lastRowIndex] = struct{}{}
				continue
			}
		}
		lastRowIndex = len(te.cellGrid) - 1
		te.cellGrid[lastRowIndex] = append(te.cellGrid[lastRowIndex], c)
	}

	te.setCursorAndScroll()

	visibleCellGrid := te.getVisibleCellGrid()

	for y, row := range visibleCellGrid {
		for x, c := range row {
			te.drawChar(x, y, c, termbox.ColorWhite, termbox.ColorDefault)
		}
	}

	for i := te.height() - te.scrollY; i <= h; i++ {
		te.drawChar(0, i, '~', termbox.ColorCyan, termbox.ColorDefault)
	}
	te.drawBottomInfoStrip()
	err := termbox.Flush()
	if err != nil {
		panic(err)
	}
}

func (te *TextEditor) MoveCursorTo(x, y int) {
	if EDITOR_START_X > x {
		te.MoveCursorUp()
		te.MoveCursorToEndOfTheLine()
		te.MoveCursorLeft()
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
			te.MoveCursorDown()
			te.MoveCursorToBeginningOfTheLine()
			te.MoveCursorRight()
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
	w, _ := termbox.Size()
	if te.CursorXOffset == w-1 {
		te.MoveCursorDown()
		te.MoveCursorToEndOfTheLine()
	} else {
		te.MoveCursorRight()
	}
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
	if te.CursorXOffset == EDITOR_START_X && te.CursorYOffset == EDITOR_START_Y {
		return
	}
	textIndex := te.getTextIndex() - 1
	runeSlice := []rune(te.Text)
	runeSlice = append(runeSlice[:textIndex], runeSlice[textIndex+1:]...)
	if te.CursorXOffset == EDITOR_START_X {
		te.MoveCursorUp()
		te.MoveCursorToEndOfTheLine()
		if _, ok := te.newLines[te.CursorYOffset]; !ok {
			te.MoveCursorLeft()
		}
	} else {
		te.MoveCursorLeft()
	}
	te.Text = string(runeSlice)
	te.Draw()
}
