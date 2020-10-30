package manual_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/srikrsna/nidhi/example/manual"
)

func TestBookMarshal(t *testing.T) {
	cmp := func(tb *manual.Book) error {
		exp, err := json.Marshal(tb)
		if err != nil {
			return err
		}

		stream := jsoniter.ConfigDefault.BorrowStream(nil)
		defer jsoniter.ConfigDefault.ReturnStream(stream)

		if err := tb.MarshalDocument(stream); err != nil {
			return err
		}

		if !bytes.Equal(exp, stream.Buffer()) {
			return fmt.Errorf("json mismatch, act: %s, exp: %s", stream.Buffer(), exp)
		}

		return nil
	}

	tt := []struct {
		Name string
		Book func() *manual.Book
	}{
		{
			"All Set",
			func() *manual.Book {
				tb := getRandomBook()
				return &tb
			},
		},
		{
			"Omit",
			func() *manual.Book {
				tb := getRandomBook()
				tb.Author = nil
				return &tb
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			if err := cmp(tc.Book()); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestBookUnmarshal(t *testing.T) {
	cmp := func(tb *manual.Book) error {
		dat, err := json.Marshal(tb)
		if err != nil {
			return err
		}

		iter := jsoniter.ConfigCompatibleWithStandardLibrary.BorrowIterator(dat)
		defer jsoniter.ConfigCompatibleWithStandardLibrary.ReturnIterator(iter)

		var act manual.Book
		if err := act.UnmarshalDocument(iter); err != nil {
			return err
		}

		if !reflect.DeepEqual(&act, tb) {
			return fmt.Errorf("unmarshal mismatch, act: %v, exp: %v", &act, tb)
		}

		return nil
	}

	tt := []struct {
		Name string
		Book func() *manual.Book
	}{
		{
			"All Set",
			func() *manual.Book {
				tb := getRandomBook()
				return &tb
			},
		},
		{
			"Omit",
			func() *manual.Book {
				tb := getRandomBook()
				tb.Author = nil
				return &tb
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			if err := cmp(tc.Book()); err != nil {
				t.Error(err)
			}
		})
	}

}

func getRandomBook() manual.Book {
	return manual.Book{
		Id:        "identifier",
		PageCount: 12,
		Author: &manual.Author{
			Name: "Name of Author",
			Bio:  "Bio of Author",
		},
		Title: "The Book",
		Pages: []*manual.Page{{
			Number:  12,
			Content: "This is a page",
		}},
	}

}
