package cswizard

import (
	"encoding/csv"
	"errors"
)

type Writer interface {
	AddHeader(string) (id uint64)
	LockHeaders() error
	CreateRow() []string
	CommitRow([]string) error
}

var (
	errHeadersLocked    = errors.New("cswizard: trying to add a header after header lock")
	errHeadersNotLocked = errors.New("cswizard: trying to add a row before header lock")
)

type BasicWriter struct {
	headers       []string
	headersLocked bool

	w   *csv.Writer
	buf []string
}

func New(w *csv.Writer) Writer {
	return &BasicWriter{
		headers:       make([]string, 0, 10),
		headersLocked: false,

		w:   w,
		buf: nil,
	}
}

func (this *BasicWriter) AddHeader(str string) (ret uint64) {
	if this.headersLocked {
		panic(errHeadersLocked)
	}

	ret = uint64(len(this.headers))
	this.headers = append(this.headers, str)
	return
}

func (this *BasicWriter) LockHeaders() error {
	this.headersLocked = true
	l := len(this.headers)
	this.buf = make([]string, l, l)
	return this.w.Write(this.headers)
}

func (this *BasicWriter) CreateRow() []string {
	if !this.headersLocked {
		panic(errHeadersNotLocked)
	}

	return this.buf
}

func (this *BasicWriter) CommitRow(row []string) (err error) {
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
