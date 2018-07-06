package tableprinter

import (
	"bytes"
	"fmt"
	"testing"
)

type (
	book struct {
		Title       string    `header:"Title"`
		Description string    `header:"Desc"`
		Sales       int       `header:"Sales"`
		Publisher   publisher `header:"inline"`
	}

	publisher struct {
		Name    string  `header:"Publisher Name"`
		Country country `header:"Publisher Country"`
	}

	country struct {
		Name string
		Code string
	}
)

func (c country) String() string {
	return c.Name
}

func buildBooks(n int) []book {
	books := make([]book, n)
	var b book

	for i := 1; i <= n; i++ {
		b.Title = fmt.Sprintf("Title for Book %d", i)
		b.Description = fmt.Sprintf("Description for Book %d", i)
		b.Sales = i * 12000
		b.Publisher = publisher{
			fmt.Sprintf("Publisher Name %d", i),
			country{fmt.Sprintf("Country Name for Publisher %d", i), "Code doesn't matter"},
		}

		books[i-1] = b
	}

	return books
}

// quite fast: using one table per printer & using cache for types that already scanned for headers.
// goos: linux
// goarch: amd64
// pkg: github.com/landoop/tableprinter
// BenchmarkPrint-8          100000             22545 ns/op            3973 B/op        181 allocs/op
// PASS
// ok      github.com/landoop/tableprinter 2.919s
func BenchmarkPrint(b *testing.B) {
	var (
		w       = new(bytes.Buffer)
		printer = New(w)
		books   = buildBooks(1)
	)

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		printer.Print(books)
	}

	b.StopTimer()
	w.Reset()
}
