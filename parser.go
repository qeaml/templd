package templd

import (
	"fmt"
	"io"
)

type Vars map[string]string

type parseError struct {
	line uint
	msg  string
}

func (e parseError) Error() string {
	return e.msg
}

type line struct {
	no      uint
	effect  byte
	content string
}

type parser struct {
	lineNo uint
	src    io.Reader
	vars   Vars
}

func (p *parser) error(format string, args ...any) error {
	return parseError{p.lineNo, fmt.Sprintf(format, args...)}
}

func (p *parser) readByte() (b byte, err error) {
	buf := make([]byte, 1)
	_, err = p.src.Read(buf)
	b = buf[0]
	return
}

func (p *parser) readLine() (line line, err error) {
	line.no = p.lineNo
	p.lineNo++
	b, err := p.readByte()
	if err != nil {
		return
	}
	line.effect = b
	for {
		b, err = p.readByte()
		if err != nil {
			return
		}
		switch b {
		case '%':
			ident := ""
			for {
				b, err = p.readByte()
				if err != nil {
					return
				}
				if b == '%' {
					break
				}
				ident += string(b)
			}
			if val, ok := p.vars[ident]; ok {
				line.content += val
			} else {
				err = p.error("unknown variable `%s`", ident)
				return
			}
		case '\n':
			return
		}
		line.content += string(b)
	}
}

func (p *parser) source(newsrc io.Reader, newvars Vars) {
	p.lineNo = 0
	p.src = newsrc
	p.vars = newvars
}
