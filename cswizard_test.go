package cswizard_test

import (
	"github.com/m1kc/go-cswizard"

	"bytes"
	"encoding/csv"
	"strings"

	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCSWizard(t *testing.T) {
	Convey("test suite", t, func() {
		Convey("Normal test", func() {
			b := new(bytes.Buffer)
			wr := csv.NewWriter(b)
			wr.Comma = ';'

			// make the CSV
			w := cswizard.New(wr)

			x := w.AddHeader("x")
			y := w.AddHeader("y")
			w.LockHeaders()

			row1 := w.CreateRow()
			row1[x] = "123"
			row1[y] = "456"
			err := w.CommitRow(row1)
			So(err, ShouldBeNil)

			row2 := w.CreateRow()
			row2[y] = "777"
			err = w.CommitRow(row2)
			So(err, ShouldBeNil)

			// flush and check
			wr.Flush()
			So(strings.TrimSpace(b.String()), ShouldEqual, "x;y\n123;456\n;777")
		})

		Convey("Errors", func() {
			Convey("Adding headers after lock", func() {
				b := new(bytes.Buffer)
				wr := csv.NewWriter(b)
				w := cswizard.New(wr)

				_ = w.AddHeader("x")
				w.LockHeaders()
				So(func() {
					_ = w.AddHeader("y")
				}, ShouldPanic)

			})

			Convey("Adding rows before lock", func() {
				b := new(bytes.Buffer)
				wr := csv.NewWriter(b)
				w := cswizard.New(wr)

				_ = w.AddHeader("x")
				_ = w.AddHeader("y")

				So(func() {
					_ = w.CreateRow()
				}, ShouldPanic)
				So(func() {
					_ = w.CommitRow(make([]string, 0, 0))
				}, ShouldPanic)
			})
		})
	})
}
