package jul

import (
	"bufio"
	"io"
	"os"
)

type UI interface {
	Write(msg string) error
	Read() (string, error)
}

type DefaultUI struct {
	r *bufio.Reader
	w io.Writer
}

func NewDefaultUI(r io.Reader, w io.Writer) *DefaultUI {
	if r == nil {
		r = os.Stdin
	}
	if w == nil {
		w = os.Stdout
	}
	return &DefaultUI{r: bufio.NewReader(r), w: w}
}

func (ui *DefaultUI) Write(msg string) error {
	_, err := io.WriteString(ui.w, msg)
	return err
}

func (ui *DefaultUI) Read() (string, error) {
	line, err := ui.r.ReadBytes('\n')
	if err != nil {
		return "", err
	}
	return string(line[:len(line)-1]), nil
}
