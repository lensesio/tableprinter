package main

import (
	"fmt"
	"os"

	"github.com/landoop/tableprinter"
)

type (
	book struct {
		Title       string    `header:"title"`
		Description string    `header:"desc"`
		Sales       int       `header:"sales"`
		Publisher   publisher `header:"inline"`
	}

	publisher struct {
		Name    string  `header:"publisher name"`
		Country country `header:"publisher country"`
	}

	country struct {
		Name string
		Code string
	}
)

func (c country) String() string {
	return c.Name
}

func main() {
	n := 5
	books := make([]book, n, n)
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

	/*
	  TITLE (5)          DESC                     SALES   PUBLISHER NAME
	 ------------------ ------------------------ ------- ------------------
	  Title for Book 1   Description for Book 1   12.0K   Publisher Name 1
	  Title for Book 2   Description for Book 2   24.0K   Publisher Name 2
	  Title for Book 3   Description for Book 3   36.0K   Publisher Name 3
	  Title for Book 4   Description for Book 4   48.0K   Publisher Name 4
	  Title for Book 5   Description for Book 5   60.0K   Publisher Name 5
	*/
	tableprinter.Print(os.Stdout, books)
}
