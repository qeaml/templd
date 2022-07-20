package templd

import (
	"errors"
	"fmt"
	"io"
)

var parse = &parser{}

func render(src io.Reader, out io.Writer, vars Vars) (err error) {
	parse.source(src, vars)
	var line line
	inCode := false
	for {
		line, err = parse.readLine()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			var perr parseError
			if errors.As(err, &perr) {
				return fmt.Errorf("error on line %d: %s", perr.line, perr)
			}
			return err
		}
		if inCode && line.effect != '`' {
			_, err = out.Write([]byte("</code>"))
			if err != nil {
				return
			}
			inCode = false
		}
		switch line.effect {
		case ' ':
			_, err = out.Write(append([]byte(line.content), ' '))
			if err != nil {
				return
			}
		case '\n', '\r':
			_, err = out.Write([]byte("<br/><br/>"))
			if err != nil {
				return
			}
		case '=':
			_, err = out.Write([]byte("<h1>"))
			if err != nil {
				return
			}
			_, err = out.Write([]byte(line.content))
			if err != nil {
				return
			}
			_, err = out.Write([]byte("</h1>"))
			if err != nil {
				return
			}
		case '-':
			_, err = out.Write([]byte("<h2>"))
			if err != nil {
				return
			}
			_, err = out.Write([]byte(line.content))
			if err != nil {
				return
			}
			_, err = out.Write([]byte("</h2>"))
			if err != nil {
				return
			}
		case '_':
			_, err = out.Write([]byte("<h3>"))
			if err != nil {
				return
			}
			_, err = out.Write([]byte(line.content))
			if err != nil {
				return
			}
			_, err = out.Write([]byte("</h3>"))
			if err != nil {
				return
			}
		case '`':
			if !inCode {
				_, err = out.Write([]byte("<code>"))
				if err != nil {
					return
				}
				inCode = true
			}
			_, err = out.Write([]byte(line.content))
			if err != nil {
				return
			}
		}
	}
}
