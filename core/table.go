package core

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"strings"
)

func emitStr(s tcell.Screen, x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)
		x += w
	}
}

type TextAlign rune

var (
	TextAlignNone   TextAlign = ' '
	TextAlignLeft   TextAlign = 'L'
	TextAlignRight  TextAlign = 'R'
	TextAlignCenter TextAlign = 'C'
)

type Pos struct {
	x, y int
}

func NewPos(x, y int) *Pos {
	return &Pos{x, y}
}

func (p *Pos) SetX(x int) *Pos {
	p.x = x
	return p
}
func (p *Pos) SetY(y int) *Pos {
	p.y = y
	return p
}

func (p *Pos) GetX() int {
	return p.x
}

func (p *Pos) GetY() int {
	return p.y
}

func (p *Pos) GetXY() (x, y int) {
	return p.x, p.y
}

func (p *Pos) GetPos() Pos {
	return Pos{p.x, p.y}
}

type Row struct {
	screen       tcell.Screen
	style        tcell.Style
	pos          *Pos      // 起始坐标
	contextAlign TextAlign // L、 C、 R

	content string

	// column
	cols []*Col

	// position X
	colX int

	// position Y
	rowY int
}

func NewRow(s tcell.Screen, pos *Pos) *Row {
	return &Row{
		screen:       s,
		style:        tcell.StyleDefault,
		contextAlign: 'L',
		pos:          pos,
		colX:         pos.GetX(),
		rowY:         pos.GetY(),
		cols:         make([]*Col, 0),
	}
}

func (r *Row) TextAlign(align TextAlign) *Row {
	r.contextAlign = align

	return r
}

func (r *Row) SetContent(content string) *Row {
	r.content = content

	return r
}

func (r *Row) SetStyle(style tcell.Style) *Row {
	r.style = style

	return r
}

func (r *Row) Print() {
	// print all columns content
	if len(r.cols) > 0 {
		str := strings.Repeat(" ", r.colX-r.pos.x-10)
		emitStr(r.screen, r.pos.x, r.pos.y, r.style, str)

		for _, col := range r.cols {
			col.Print()
		}

		return
	}

	text := r.content
	// print row's content
	w, _ := r.screen.Size()

	switch r.contextAlign {
	case 'L': // left
		r.pos.SetX(0)
	case 'C': // cente
		length := len(text)
		r.pos.SetX(w/2 - length)
	case 'R': // right
		length := len(text)
		r.pos.SetX(w - length)
	}

	// x := r.pos.GetX()
	// y := r.getRowY()

	x, y := r.pos.GetXY()

	emitStr(r.screen, x, y, r.style, text)
}

func (r *Row) getColX(width int) int {
	x := r.colX
	// next td postion X
	r.colX += width
	return x
}

func (r *Row) AddCol(col *Col) *Row {
	width := col.GetWidth()
	col.SetPos(NewPos(r.getColX(width), r.pos.GetY()))
	col.SetStyle(r.style)

	r.cols = append(r.cols, col)
	return r
}

func (r *Row) AddColWithWidth(text string, width int) *Row {
	col := NewCol(r.screen, text)
	col.SetWidth(width).
		SetPos(NewPos(r.getColX(width), r.pos.GetY())).
		SetStyle(r.style)

	r.cols = append(r.cols, col)
	return r
}

func (r *Row) AddColl(text string, opts ...Option) *Row {
	col := NewCol(r.screen, text, opts...)
	r.cols = append(r.cols, col)
	return r
}

type Option func(*Col)

func WithWidth(n int) Option {
	return func(opt *Col) {
		opt.width = n
	}
}
func WithContent(str string) Option {
	return func(opt *Col) {
		opt.content = str
	}
}

func WithPos(pos *Pos) Option {
	return func(opt *Col) {
		opt.pos = pos
	}
}

func WithStyle(style tcell.Style) Option {
	return func(opt *Col) {
		opt.style = style
	}
}

type Col struct {
	screen  tcell.Screen
	style   tcell.Style
	pos     *Pos
	content string
	width   int
}

var defaultColWidth = 10

func NewCol(s tcell.Screen, text string, opts ...Option) *Col {
	col := &Col{
		screen:  s,
		style:   tcell.StyleDefault,
		content: text,
		width:   defaultColWidth,
	}

	for _, apply := range opts {
		apply(col)
	}

	return col
}

func (c *Col) SetPos(pos *Pos) *Col {
	c.pos = pos
	return c
}

func (c *Col) SetContent(content string) *Col {
	c.content = content
	return c
}
func (c *Col) SetStyle(style tcell.Style) *Col {
	c.style = style

	return c
}

func (c *Col) SetWidth(width int) *Col {
	c.width = width
	return c
}

func (c *Col) GetWidth() int {
	return c.width
}

func (c *Col) Print() {
	text := c.content
	if len(text) > c.width {
		text = text[:c.width]
	}

	x, y := c.pos.GetXY()

	emitStr(c.screen, x, y, c.style, text)
}

type Table struct {
	screen tcell.Screen
	style  tcell.Style
	pos    *Pos
	rows   []*Row
}

func NewTable(s tcell.Screen, pos *Pos) *Table {
	t := &Table{
		screen: s,
		style:  tcell.StyleDefault,
		pos:    pos,
		rows:   make([]*Row, 0),
	}

	return t
}

func (t *Table) print() {
	for _, r := range t.rows {
		for _, td := range r.cols {
			fmt.Printf("%#v\n", td)
		}
	}
}

func (t *Table) AddRow() *Row {
	row := NewRow(t.screen, NewPos(t.pos.x, t.pos.y))
	t.pos.y++
	// move Y position

	return row
}

func (t *Table) GetPos() Pos {
	return Pos{
		x: t.pos.x,
		y: t.pos.y,
	}
}
