package main

type Operation struct {
	Operation  OperationType
	Text       string
	StartIndex int
}

type OperationType int

const (
	Insert OperationType = iota
	Remove
	Base
)

type UndoRedoNode struct {
	Operation       OperationType
	Text            string
	StartIndex      int
	Parent          *UndoRedoNode
	UndoOriginChild *UndoRedoNode
}

type UndoRedoTree struct {
	Head *UndoRedoNode
}

func InitUndoRedoTree() *UndoRedoTree {
	initNode := UndoRedoNode{
		Operation: Base,
	}
	return &UndoRedoTree{
		Head: &initNode,
	}
}

func (utr *UndoRedoTree) AddNode(op OperationType, text string, startIndex int) {
	newNode := UndoRedoNode{
		Operation:  op,
		Text:       text,
		StartIndex: startIndex,
		Parent:     utr.Head,
	}

	utr.Head = &newNode
}

func (utr *UndoRedoTree) Undo() *Operation {
	if utr == nil || utr.Head == nil || utr.Head.Operation == Base {
		return nil
	}
	operation := Operation{
		Operation:  utr.Head.Operation,
		Text:       utr.Head.Text,
		StartIndex: utr.Head.StartIndex,
	}

	if utr.Head.Parent != nil {
		utr.Head.Parent.UndoOriginChild = utr.Head
	}

	utr.Head = utr.Head.Parent
	return &operation
}

func (utr *UndoRedoTree) Redo() *Operation {
	if utr.Head.UndoOriginChild == nil {
		return nil
	}

	utr.Head = utr.Head.UndoOriginChild

	operation := Operation{
		Operation:  utr.Head.Operation,
		Text:       utr.Head.Text,
		StartIndex: utr.Head.StartIndex,
	}

	return &operation
}
