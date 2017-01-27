package proto3parser

import (
	"fmt"
	"strconv"
	"strings"
)

type Field struct {
	Name     string
	Type     string
	Repeated bool
	Messages []*Message
	Sequence int
}

// ParseField parsers one field.
func ParseField(f *Field, p *Parser) error {
	for {
		tok, lit := p.scanIgnoreWhitespace()
		switch tok {
		case IDENT:
			// normal type?
			if strings.Contains(TypeTokens, lit) {
				f.Type = lit
				return ParseNormalField(f, p)
			}
			//if tok == ONEOF {}
			//if tok == ONEOFFIELD {}
		case MESSAGE:
			m, err := ParseMessage(p)
			if err != nil {
				return err
			}
			f.Messages = append(f.Messages, m)
		case REPEATED:
			f.Repeated = true
			return ParseField(f, p)
		default:
			goto done
		}
	}
done:
	return nil
}

// ParseNormalField proceeds after reading the type of f.
func ParseNormalField(f *Field, p *Parser) error {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return fmt.Errorf("found %q, expected identifier", lit)
	}
	f.Name = lit
	tok, lit = p.scanIgnoreWhitespace()
	if tok != EQUALS {
		return fmt.Errorf("found %q, expected =", lit)
	}
	_, lit = p.scanIgnoreWhitespace()
	i, err := strconv.Atoi(lit)
	if err != nil {
		return fmt.Errorf("found %q, expected sequence number", lit)
	}
	f.Sequence = i
	return nil
}
