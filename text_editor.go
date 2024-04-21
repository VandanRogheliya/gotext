package main

import (
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/nsf/termbox-go"
)

const EDITOR_START_X = 0
const EDITOR_START_Y = 0

type selection struct {
	StartX int
	StartY int
}

type TextEditor struct {
	Text                string
	FileName            string
	CursorXOffset       int
	CursorXOffsetActual int
	CursorYOffset       int
	cellGrid            [][]rune
	scrollY             int
	newLines            map[int]struct{}
	Selection           selection
	isSelection         bool
	undoRedoTree        UndoRedoTree
}

func InitTextEditor(text string, fileName string) *TextEditor {
	return &TextEditor{
		Text:                text,
		CursorXOffset:       EDITOR_START_X,
		CursorYOffset:       EDITOR_START_Y,
		CursorXOffsetActual: EDITOR_START_X,
		FileName:            fileName,
		scrollY:             EDITOR_START_Y,
		isSelection:         false,
		undoRedoTree:        *InitUndoRedoTree(),
	}
}

func (te *TextEditor) getSelectionRange() (startX int, startY int, endX int, endY int) {
	startY = int(math.Min(float64(te.CursorYOffset), float64(te.Selection.StartY)))
	endY = int(math.Max(float64(te.CursorYOffset), float64(te.Selection.StartY)))

	if startY == endY {
		startX = int(math.Min(float64(te.CursorXOffset), float64(te.Selection.StartX)))
		endX = int(math.Max(float64(te.CursorXOffset), float64(te.Selection.StartX)))
	} else {
		if startY == te.Selection.StartY {
			startX = te.Selection.StartX
			endX = te.CursorXOffset
		} else {
			startX = te.CursorXOffset
			endX = te.Selection.StartX

		}
	}

	return startX, startY, endX, endY
}

func (te *TextEditor) checkIsCharSelected(x, y int) bool {
	if !te.isSelection {
		return false
	}

	startX, startY, endX, endY := te.getSelectionRange()

	if y < startY || y > endY {
		return false
	}

	if startY == endY {
		return startX <= x && x <= endX
	} else if y == startY {
		return startX <= x
	} else if y == endY {
		return x <= endX
	}

	return true
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

func (te *TextEditor) getTextIndex(x, y int) int {
	totalChars := 0
	for i := 0; i < y; i++ {
		totalChars += len(te.cellGrid[i]) - EDITOR_START_X
		if _, ok := te.newLines[i]; ok {
			totalChars++
		}
	}
	totalChars += x - EDITOR_START_X
	textIndex := totalChars
	return textIndex
}

func (te *TextEditor) drawBottomInfoStrip() {
	infoText := fmt.Sprintf("%d,%d | %s | Selection (Ctrl+Q): %t", te.CursorXOffset, te.CursorYOffset, te.FileName, te.isSelection)
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

func (te *TextEditor) getSelectedTextIndexRange() (startIndex int, endIndex int) {
	startX, startY, endX, endY := te.getSelectionRange()

	startIndex = te.getTextIndex(startX, startY)
	endIndex = te.getTextIndex(endX, endY)
	return startIndex, endIndex
}

func (te *TextEditor) getSelectedText() string {
	startIndex, endIndex := te.getSelectedTextIndexRange()
	return te.Text[startIndex:endIndex]
}

func (te *TextEditor) Draw() {
	te.clearEditor()
	w, h := termbox.Size()
	te.cellGrid = [][]rune{}
	te.newLines = make(map[int]struct{})
	te.cellGrid = append(te.cellGrid, []rune{})

	for _, c := range te.Text {
		lastRowIndex := len(te.cellGrid) - 1
		if c == '\n' || len(te.cellGrid[lastRowIndex]) > w {
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
			isSelected := te.checkIsCharSelected(x, y)
			fg := termbox.ColorWhite
			bg := termbox.ColorDefault
			if isSelected {
				fg = termbox.ColorBlack
				bg = termbox.ColorWhite
			}
			te.drawChar(x, y, c, fg, bg)
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
		if y == EDITOR_START_Y {
			return
		}
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
			if y == te.height()-1 {
				return
			}
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

func (te *TextEditor) getLinePosFromTextIndex(textIndex int) int {
	w, _ := termbox.Size()
	lines := strings.Split(te.Text, "\n")
	editorLines := [][]string{}
	for _, line := range lines {
		if len(line) > w {
			runeSlice := []rune(line)
			newLineCnt := len(line) / w
			startIndex := 0
			endIndex := w
			newLines := []string{}
			for i := 0; i <= newLineCnt; i++ {
				endIndex = int(math.Min(float64(endIndex), float64(len(runeSlice))))
				newRuneSlice := runeSlice[startIndex:endIndex]
				newLine := string(newRuneSlice)
				newLines = append(newLines, newLine)
				startIndex += w
				endIndex += w
			}
			editorLines = append(editorLines, newLines)

		} else {
			editorLines = append(editorLines, []string{line})
		}
	}

	y := EDITOR_START_Y - 1
	charsCnt := 0
	for _, lines := range editorLines {
		for _, line := range lines {
			y++
			charsCnt += len(line)
			if charsCnt >= textIndex {
				return y
			}
		}
		charsCnt++
	}
	height := te.height()
	return height
}

func (te *TextEditor) InsertString(text string) {
	if te.isSelection {
		te.RemoveSelection()
	}
	textIndex := te.getTextIndex(te.CursorXOffset, te.CursorYOffset)
	runeSlice := []rune(te.Text)
	runeSlice = append(runeSlice[:textIndex], append([]rune(text), runeSlice[textIndex:]...)...)
	te.Text = string(runeSlice)
	te.Draw()
	w, _ := termbox.Size()
	if te.CursorXOffset == w-1 {
		te.MoveCursorDown()
		te.MoveCursorToEndOfTheLine()
	} else {
		te.MoveCursorRight()
	}
	te.undoRedoTree.AddNode(Insert, text, textIndex)
}

func (te *TextEditor) AddNewLine() {
	te.RemoveSelection()
	textIndex := te.getTextIndex(te.CursorXOffset, te.CursorYOffset)
	runeSlice := []rune(te.Text)
	runeSlice = append(runeSlice[:textIndex], append([]rune{'\n'}, runeSlice[textIndex:]...)...)
	te.Text = string(runeSlice)
	te.MoveCursorToBeginningOfTheLine()
	te.MoveCursorDown()
	te.Draw()
}

func (te *TextEditor) RemoveChar() {
	if te.isSelection {
		te.RemoveSelection()
		return
	}
	if te.CursorXOffset == EDITOR_START_X && te.CursorYOffset == EDITOR_START_Y {
		return
	}
	textIndex := te.getTextIndex(te.CursorXOffset, te.CursorYOffset) - 1
	removedChar := te.Text[textIndex]
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
	te.undoRedoTree.AddNode(Remove, string(removedChar), textIndex)
}

func (te *TextEditor) RemoveSelection() {
	if !te.isSelection {
		return
	}
	startIndex, endIndex := te.getSelectedTextIndexRange()
	selectedText := te.Text[startIndex:endIndex]
	runeSlice := []rune(te.Text)
	runeSlice = append(runeSlice[:startIndex], runeSlice[endIndex:]...)
	startX, startY, _, _ := te.getSelectionRange()
	te.MoveCursorTo(startX, startY)
	te.Text = string(runeSlice)
	te.ToggleSelectionMode()
	te.Draw()
	te.undoRedoTree.AddNode(Remove, selectedText, startIndex)
}

func (te *TextEditor) ToggleSelectionMode() {
	te.isSelection = !te.isSelection
	if te.isSelection {
		te.Selection.StartX = te.CursorXOffset
		te.Selection.StartY = te.CursorYOffset
	}
	te.Draw()
}

func (te *TextEditor) CopySelection() {
	if !te.isSelection {
		return
	}
	selectedText := te.getSelectedText()
	err := clipboard.WriteAll(selectedText)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func (te *TextEditor) CutSelection() {
	if !te.isSelection {
		return
	}
	selectedText := te.getSelectedText()
	err := clipboard.WriteAll(selectedText)
	if err != nil {
		log.Fatalf(err.Error())
	}
	te.RemoveSelection()
}

func (te *TextEditor) Paste() {
	text, err := clipboard.ReadAll()
	if err != nil {
		log.Fatalf(err.Error())
	}
	te.InsertString(text)
}

func (te *TextEditor) Undo() {
	op := te.undoRedoTree.Undo()
	if op == nil {
		return
	}
	if op.Operation == Insert {
		runeSlice := []rune(te.Text)
		startIndex := op.StartIndex
		endIndex := startIndex + len(op.Text)
		linePos := te.getLinePosFromTextIndex(startIndex)
		te.MoveCursorTo(EDITOR_START_X, linePos)
		runeSlice = append(runeSlice[:startIndex], runeSlice[endIndex:]...)
		te.Text = string(runeSlice)
		te.Draw()
	} else if op.Operation == Remove {
		textIndex := op.StartIndex
		linePos := te.getLinePosFromTextIndex(textIndex)
		te.MoveCursorTo(EDITOR_START_X, linePos)
		runeSlice := []rune(te.Text)
		runeSlice = append(runeSlice[:textIndex], append([]rune(op.Text), runeSlice[textIndex:]...)...)
		te.Text = string(runeSlice)
		te.Draw()
	}
}

func (te *TextEditor) Redo() {
	op := te.undoRedoTree.Redo()
	if op == nil {
		return
	}

	if op.Operation == Remove {
		runeSlice := []rune(te.Text)
		startIndex := op.StartIndex
		endIndex := startIndex + len(op.Text)
		linePos := te.getLinePosFromTextIndex(startIndex)
		te.MoveCursorTo(EDITOR_START_X, linePos)
		runeSlice = append(runeSlice[:startIndex], runeSlice[endIndex:]...)
		te.Text = string(runeSlice)
		te.Draw()
	} else if op.Operation == Insert {
		textIndex := op.StartIndex
		linePos := te.getLinePosFromTextIndex(textIndex)
		te.MoveCursorTo(EDITOR_START_X, linePos)
		runeSlice := []rune(te.Text)
		runeSlice = append(runeSlice[:textIndex], append([]rune(op.Text), runeSlice[textIndex:]...)...)
		te.Text = string(runeSlice)
		te.Draw()
	}
}
