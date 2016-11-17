package boolean

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

type item struct {
	typ   itemType
	start int
	val   string
}

func (i item) String() string {
	return fmt.Sprintf("%s %s %d", i.typ, i.val, i.start)
}

type itemType int

const (
	itemError itemType = iota
	itemIdentifier
	itemAnd
	itemOr
	itemNot
	itemLeftParen
	itemRightParen
	itemEOF
)

var key = map[string]itemType{
	"and": itemAnd,
	"or":  itemOr,
	"not": itemNot,
}

const eof = -1

type stateFn func(l *lexer) stateFn

type lexer struct {
	input                         string
	start, pos, width, parenDepth int
	items                         chan item
}

func lex(input string) *lexer {
	l := lexer{
		input: input,
		items: make(chan item),
	}
	go l.run()
	return &l
}

func (l *lexer) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.items)
}

// for non error the parser should call this until itemEOF
func (l *lexer) nextItem() item {
	return <-l.items
}

func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.start, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, width := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = width
	l.pos += width
	return r
}

func (l *lexer) backup() {
	l.pos = l.pos - l.width
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{itemError, l.start, fmt.Sprintf(format, args...)}
	return nil
}

// drain reads all remaining tokens
// TODO: why would we need to continue lexing if all we want to do is stop ?
func (l *lexer) drain() {
	for range l.items {
	}
}

func lexText(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == '(':
			l.emit(itemLeftParen)
			l.parenDepth++
		case r == ')':
			if l.parenDepth == 0 {
				return l.errorf("unexpected right parentheses")
			}
			l.parenDepth--
			l.emit(itemRightParen)
		case isSpace(r):
			l.ignore()
		case isAlphaNumeric(r):
			return lexIdentifier
		case r == eof:
			l.emit(itemEOF)
			return nil
		default:
			return l.errorf("bad character")
		}
	}
	return nil
}

func lexIdentifier(l *lexer) stateFn {
	for {
		r := l.next()
		if !isAlphaNumeric(r) {
			l.backup()
			if v, ok := key[l.input[l.start:l.pos]]; ok {
				l.emit(v)
				return lexText
			}
			l.emit(itemIdentifier)
			return lexText
		}
	}
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}
