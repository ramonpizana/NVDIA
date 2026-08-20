package vcard_test

import (
	"reflect"
	"testing"

	"github.com/ramonpizana/NVDIA/internal/vcard"
)

func TestVcard(t *testing.T) {
	t.Run("add a vcard", func(t *testing.T) {
		got := vcard.New("RTX 5090", "Tienda", 1500, "some link", false)

		want := &vcard.Vcard{
			Name:  "RTX 5090",
			Store: "Tienda",
			Price: 1500,
			Link:  "some link",
			Stock: false,
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("wanted %v, got %v", want, got)
		}

	})
	t.Run("add 2 cards", func(t *testing.T) {
		got := vcard.New("RTX 5090", "Tienda", 1500, "some link", false)
		got2 := vcard.New("RTX 5090 OC", "Otra tienda", 1700, "another link", true)

		want := &vcard.Vcard{
			Name:  "RTX 5090",
			Store: "Tienda",
			Price: 1500,
			Link:  "some link",
			Stock: false,
		}
		want2 := &vcard.Vcard{
			Name:  "RTX 5090 OC",
			Store: "Otra tienda",
			Price: 1700,
			Link:  "another link",
			Stock: true,
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("wanted %v, got %v", want, got)
		}

		if !reflect.DeepEqual(want2, got2) {
			t.Errorf("wanted %v, got %v", want2, got2)
		}

	})
}
