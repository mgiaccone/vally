package scanner

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

var (
	eof = rune(0)
)

// Scanner is an implementation of a scanner for lexing an expression
type Scanner struct {
	buf  strings.Builder // buffer for the current token's literal
	pos  int             // current position in the input.
	r    *bufio.Reader   // input reader.
	done bool            // flag to signal done processing
}

// New creates a new scanner for the given input io.Reader.
func New(r io.Reader) *Scanner {
	s := &Scanner{
		pos: 0,
		r:   bufio.NewReader(r),
	}
	return s
}

// Scan moves and returns the next available token
// nolint:gocognit,gocyclo
func (s *Scanner) Scan() Token {
	if s.done {
		return newToken(EOF, "", s.pos, s.pos)
	}

	var (
		ch       rune
		startPos int
		err      error
	)
	for {
		ch = s.read()

		switch ch {
		case eof:
			s.done = true
			return newToken(EOF, "", s.pos, s.pos)
		case '\r', '\n':
			s.done = true
			return newToken(ILLEGAL, "illegal multiple lines", s.pos, s.pos)
		case '(':
			return newToken(LPAREN, "(", s.pos, s.pos)
		case ')':
			return newToken(RPAREN, ")", s.pos, s.pos)
		case ',':
			return newToken(COMMA, ",", s.pos, s.pos)
		case '!':
			return newToken(LNOT, "!", s.pos, s.pos)
		case '&':
			startPos = s.backup()
			if err = scanOperatorLiteral(s, "&&", &s.buf); err != nil {
				s.done = true
				return newToken(ILLEGAL, err.Error(), startPos, s.pos)
			}
			return newToken(LAND, s.buf.String(), startPos, s.pos)
		case '|':
			startPos = s.backup()
			if err = scanOperatorLiteral(s, "||", &s.buf); err != nil {
				s.done = true
				return newToken(ILLEGAL, err.Error(), startPos, s.pos)
			}
			return newToken(LOR, s.buf.String(), startPos, s.pos)
		case '\'':
			startPos = s.backup()
			if err = scanStringLiteral(s, &s.buf); err != nil {
				s.done = true
				return newToken(ILLEGAL, err.Error(), s.pos, s.pos)
			}
			return newToken(STRING, s.buf.String(), startPos, s.pos)
		case '.':
			startPos = s.backup()
			if err = scanFielfRefLiteral(s, &s.buf); err != nil {
				return newToken(ILLEGAL, err.Error(), s.pos, s.pos)
			}
			return newToken(FIELDREF, s.buf.String(), startPos, s.pos)
		}

		if isIdentPrefix(ch) {
			startPos = s.backup()
			if err = scanIdentLiteral(s, &s.buf); err != nil {
				s.done = true
				return newToken(ILLEGAL, err.Error(), s.pos, s.pos)
			}
			return newToken(IDENT, s.buf.String(), startPos, s.pos)
		}

		if isNumberPrefix(ch) {
			startPos = s.backup()
			if err = scanNumberLiteral(s, &s.buf); err != nil {
				s.done = true
				return newToken(ILLEGAL, err.Error(), s.pos, s.pos)
			}
			lit := s.buf.String()
			tt := INT
			if strings.Contains(lit, ".") {
				tt = FLOAT
			}
			return newToken(tt, lit, startPos, s.pos)
		}

		if unicode.IsSpace(ch) {
			continue
		}

		s.done = true
		return newToken(ILLEGAL, fmt.Sprintf("illegal character %q", string(ch)), s.pos, s.pos)
	}
}

func scanFielfRefLiteral(s *Scanner, buf *strings.Builder) error {
	// make sure the literal buffer for the string is empty
	buf.Reset()

	var ch rune
	for {
		ch = s.read()

		switch ch {
		case eof:
			return nil
		case ',', ')', ' ':
			s.backup()
			return nil
		default:
			if isFieldRefLiteral(ch) {
				buf.WriteRune(ch)
				continue
			}
			return fmt.Errorf("illegal field reference character %q", string(ch))
		}
	}
}

// func scanFlagLiteral(s *Scanner, buf *strings.Builder) error {
// 	// make sure the literal buffer for the string is empty
// 	buf.Reset()
//
// 	var ch rune
// 	for {
// 		ch = s.read()
//
// 		switch ch {
// 		case eof:
// 			return nil
// 		case '(', ')', ',', ' ':
// 			s.backup()
// 			return nil
// 		default:
// 			if isIdentLiteral(ch) || ch == '_' {
// 				buf.WriteRune(ch)
// 				continue
// 			}
// 			return fmt.Errorf("illegal flag character %q", string(ch))
// 		}
// 	}
// }

func scanIdentLiteral(s *Scanner, buf *strings.Builder) error {
	// make sure the literal buffer for the string is empty
	buf.Reset()

	var ch rune
	for {
		ch = s.read()

		switch ch {
		case eof:
			return nil
		case '(', ')', ',', ' ':
			s.backup()
			return nil
		default:
			if isIdentLiteral(ch) || ch == '_' {
				buf.WriteRune(ch)
				continue
			}
			return fmt.Errorf("illegal identifier character %q", string(ch))
		}
	}
}

// nolint:gocognit,gocyclo
func scanNumberLiteral(s *Scanner, buf *strings.Builder) error {
	// make sure the literal buffer for the string is empty
	buf.Reset()

	var (
		ds bool // decimal separator flag
		ch rune
	)
	for {
		ch = s.read()

		switch ch {
		case eof:
			return nil
		case ',', ')', ' ':
			s.backup()
			return nil
		case '+', '-':
			// only allowed as first character
			if buf.Len() == 0 {
				buf.WriteRune(ch)
				continue
			}
			buf.WriteRune(ch)
			return fmt.Errorf("illegal numeric value %q", buf.String())
		case '.':
			// only once
			if ds {
				return fmt.Errorf("numbers cannot have multiple decimal separators")
			}

			buf.WriteRune(ch)
			ds = true
		default:
			if unicode.IsDigit(ch) {
				if buf.Len() == 0 && ch == '0' {
					nextCh := s.peek()
					if unicode.IsDigit(nextCh) {
						return fmt.Errorf("integer numbers cannot begin with 0")
					}
				}

				buf.WriteRune(ch)
				continue
			}
			return fmt.Errorf("illegal digit %q", string(ch))
		}
	}
}

func scanOperatorLiteral(s *Scanner, op string, buf *strings.Builder) error {
	// make sure the literal buffer for the string is empty
	buf.Reset()

	var ch rune
	for i := 0; i < len(op); i++ {
		ch = s.read()
		if ch == eof {
			continue
		}
		buf.WriteRune(ch)
	}

	readOp := buf.String()
	if op != readOp {
		return fmt.Errorf("illegal operator %q", readOp)
	}

	return nil
}

// nolint:gocognit,gocyclo
func scanStringLiteral(s *Scanner, buf *strings.Builder) error {
	// make sure the literal buffer for the string is empty
	buf.Reset()

	var (
		quoteOpen bool
		ch        rune
	)
	for {
		ch = s.read()

		switch ch {
		case eof:
			if quoteOpen {
				return fmt.Errorf("unterminated string %q", buf.String())
			}
			return nil
		case '\'':
			buf.WriteRune(ch)

			// beginning of string
			if !quoteOpen && buf.Len() == 1 {
				quoteOpen = true
				continue
			}

			// end of string
			if quoteOpen {
				return nil
			}

			return fmt.Errorf("unexpected state while reading string %q", buf.String())
		default:
			if !unicode.IsControl(ch) && !(ch == '\r' || ch == '\n' || ch == '\t') {
				buf.WriteRune(ch)
				continue
			}
			return fmt.Errorf("illegal string character %q", string(ch))
		}
	}
}

// backup undo reading the last rune from the input.
func (s *Scanner) backup() int {
	if err := s.r.UnreadRune(); err != nil {
		panic(fmt.Errorf("backup input: %w", err))
	}
	s.pos--
	return s.pos
}

// read returns and consume the next rune in the input.
func (s *Scanner) read() rune {
	ch, size, err := s.r.ReadRune()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return eof
		}
		panic(fmt.Errorf("read input: %w", err))
	}
	s.pos += size
	return ch
}

// peek returns the next rune in the input without consuming it.
func (s *Scanner) peek() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return eof
		}
		panic(fmt.Errorf("read rune: %w", err))
	}
	if err = s.r.UnreadRune(); err != nil {
		panic(fmt.Errorf("unread rune: %w", err))
	}

	return ch
}

func isFieldRefLiteral(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '.'
}

// func isFlagLiteral(ch rune) bool {
// 	return unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_'
// }
//
// func isFlagPrefix(ch rune) bool {
// 	return unicode.IsLower(ch)
// }

func isIdentLiteral(ch rune) bool {
	return unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_'
}

func isIdentPrefix(ch rune) bool {
	return unicode.IsLower(ch)
}

func isNumberPrefix(ch rune) bool {
	return unicode.IsDigit(ch) || ch == '+' || ch == '-'
}
