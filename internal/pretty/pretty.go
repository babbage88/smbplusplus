package pretty

import "fmt"

type PrettyPrintOption func(p *prettyPrinter)

type PrettyPrinter interface {
	Print(s string)
	PrintWarning(s string)
	PrintError(s string)
}

type prettyPrinter struct {
	InfoColor int32 `json:"infoColor"`
	WarnColor int32 `json:"warnColor"`
	ErrColor  int32 `json:"errorColor"`
}

func WithInfoColor(c int32) PrettyPrintOption {
	return func(p *prettyPrinter) {
		p.InfoColor = c
	}
}

func WithWarnColor(c int32) PrettyPrintOption {
	return func(p *prettyPrinter) {
		p.WarnColor = c
	}
}

func WithErrColor(c int32) PrettyPrintOption {
	return func(p *prettyPrinter) {
		p.ErrColor = c
	}
}

func NewPrettyPrinter(opts ...PrettyPrintOption) *prettyPrinter {
	const (
		infoColor = int32(92)
		warnColor = int32(93)
		errColor  = int32(91)
	)

	printer := &prettyPrinter{
		InfoColor: infoColor,
		WarnColor: warnColor,
		ErrColor:  errColor,
	}
	for _, opt := range opts {
		opt(printer)
	}
	return printer
}

func (p *prettyPrinter) Print(s string) {
	fmt.Printf("\x1b[1;%dm%s\x1b[0m\n", p.InfoColor, s)
}

func (p *prettyPrinter) PrintWarning(s string) {
	fmt.Printf("\x1b[1;%dm%s\x1b[0m\n", p.WarnColor, s)
}

func (p *prettyPrinter) PrintError(s string) {
	fmt.Printf("\x1b[1;%dm%s\x1b[0m\n", p.ErrColor, s)
}

func Print(s string) {
	printer := NewPrettyPrinter()
	printer.Print(s)
}

func PrintWarning(s string) {
	printer := NewPrettyPrinter()
	printer.PrintWarning(s)
}

func PrintError(s string) {
	printer := NewPrettyPrinter()
	printer.PrintError(s)
}
