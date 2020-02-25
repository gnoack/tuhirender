package imgwriter

import "io"

type W interface {
	Step()
	Write(io.Writer) error
}
