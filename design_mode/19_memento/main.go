package main

import "fmt"

type Originator interface {
	Save() Memento
	Restore(m Memento)
}

type Memento interface {
	GetState() string
}

type TextEditor struct {
	state string
}

func (t *TextEditor) Save() Memento {
	return &textMemento{state: t.state}
}

func (t *TextEditor) Restore(m Memento) {
	t.state = m.GetState()
}

func (t *TextEditor) SetState(state string) {
	t.state = state
}

func (t *TextEditor) GetState() string {
	return t.state
}

type textMemento struct {
	state string
}

func (t *textMemento) GetState() string {
	return t.state
}

type Caretaker struct {
	mementos     []Memento
	currentIndex int
}

func (c *Caretaker) AddMemento(m Memento) {
	c.mementos = append(c.mementos, m)
	c.currentIndex = len(c.mementos) - 1
}

func (c *Caretaker) Undo(t *TextEditor) {
	if c.currentIndex > 0 {
		c.currentIndex--
		m := c.mementos[c.currentIndex]
		t.Restore(m)
	}
}

func (c *Caretaker) Redo(t *TextEditor) {
	if c.currentIndex < len(c.mementos)-1 {
		c.currentIndex++
		m := c.mementos[c.currentIndex]
		t.Restore(m)
	}
}

func main() {
	editor := &TextEditor{}
	caretaker := &Caretaker{}

	editor.SetState("State #1")
	caretaker.AddMemento(editor.Save())

	editor.SetState("State #2")
	caretaker.AddMemento(editor.Save())

	editor.SetState("State #3")
	caretaker.AddMemento(editor.Save())

	caretaker.Undo(editor)
	caretaker.Undo(editor)
	fmt.Println(editor.GetState())

	caretaker.Redo(editor)
	fmt.Println(editor.GetState())
}
