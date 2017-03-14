// Package cswizard is a CSV writer that doesn't stand in your way as your
// system evolves. Using `encoding/csv` directly is fine when you want to just
// do the thing and forget it. Long-living projects, however, are rarely
// done this way: every day business demands to add new columns,
// remove and reorder them, and under this conditions `encoding/csv`
// becomes too fragile due to nature of its API: you just can't change one thing
// and be sure that everything else would keep working. With CSWizard, you can.
//
// So, in a nutshell, that's a small wrapper around `encoding/csv` for reports
// that change often in various ways.
//
// A typical usage would be something like that:
//
//	cw := csv.NewWriter(os.Stdout)
//	w := cswizard.New(cw)
//
//	colName := w.AddHeader("Client name")
//	colAge := w.AddHeader("Client age")
//	colHeight := w.AddHeader("Client height (predicted)")
//	w.LockHeaders()
//
//	for _, c := range clients {
//		row := w.CreateRow()
//		row[colName] = c.Name
//		row[colAge] = strconv.FormatUint(c.Age, 10)
//		row[colHeight] = strconv.FormatUint(c.Height, 10)
//
//		err := w.CommitRow(row)
//		if err != nil {
//			return
//		}
//	}
//	cw.Flush()
package cswizard

import (
	"encoding/csv"
	"errors"
)

// Writer provides methods for composing a CSV.
type Writer interface {
	// AddHeader adds a new header and returns its index which you can
	// use when you'll be filling a row.
	AddHeader(string) (id uint64)
	// LockHeaders locks header list. After that you can start adding the rows.
	LockHeaders() error
	// CreateRow creates a new row and returns a buffer which you can populate
	// using indexes you've got from AddHeader calls, and later commit with
	// CommitRow.
	//
	// Internally, this buffer is reused between calls. This may change
	// in the future.
	CreateRow() []string
	// CommitRow accepts a buffer of columns (which you normally obtain
	// from CreateRow) and adds it into CSV.
	CommitRow([]string) error
}

var (
	errHeadersLocked    = errors.New("cswizard: trying to add a header after header lock")
	errHeadersNotLocked = errors.New("cswizard: trying to add a row before header lock")
)

type basicWriter struct {
	headers       []string
	headersLocked bool

	w   *csv.Writer
	buf []string
}

// New wraps an existing csv.Writer into cswizard.Writer.
func New(w *csv.Writer) Writer {
	return &basicWriter{
		headers:       make([]string, 0, 10),
		headersLocked: false,

		w:   w,
		buf: nil,
	}
}

func (this *basicWriter) AddHeader(str string) (ret uint64) {
	if this.headersLocked {
		panic(errHeadersLocked)
	}

	ret = uint64(len(this.headers))
	this.headers = append(this.headers, str)
	return
}

func (this *basicWriter) LockHeaders() error {
	this.headersLocked = true
	l := len(this.headers)
	this.buf = make([]string, l, l)
	return this.w.Write(this.headers)
}

func (this *basicWriter) CreateRow() []string {
	if !this.headersLocked {
		panic(errHeadersNotLocked)
	}

	return this.buf
}

func (this *basicWriter) CommitRow(row []string) (err error) {
	if !this.headersLocked {
		panic(errHeadersNotLocked)
	}

	err = this.w.Write(row)
	if err != nil {
		return
	}

	for i := 0; i < len(this.buf); i++ {
		this.buf[i] = ""
	}
	return
}
