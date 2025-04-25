package file

import "fmt"

type Kind int

const (
	UnknownKind = Kind(iota)
	Axs
	Axi
	Axb
	Lib
)

func (k Kind) String() string {
	switch k {
	case Axs:
		return "Axs"
	case Axi:
		return "Axi"
	case Axb:
		return "Axb"
	case Lib:
		return "Lib"
	default:
		return fmt.Sprintf("internal error: unknown file kind %d", k)
	}
}
