package search

import (
	"net/url"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestParseStoreDocument(t *testing.T) {
	tests := []struct {
		name       string
		store      store
		html       string
		wantMatch  int
		wantCards  int
		wantPrice  int
		wantInStock int
	}{
		{
			name: "DDTech current markup",
			store: mexicoStores[0],
			html: `<div class="product">
				<h3 class="name"><a href="/producto/rtx-5090?id=1">Tarjeta NVIDIA GeForce RTX 5090</a></h3>
				<div class="product-price"><span class="price">$79,999.00</span></div>
				<div class="stock"><span class="with-stock">CON EXISTENCIA</span></div>
			</div>`,
			wantMatch: 1, wantCards: 1, wantPrice: 79999, wantInStock: 1,
		},
		{
			name: "Cyberpuerta excludes unavailable card",
			store: mexicoStores[1],
			html: `<div class="cpd-product-card-catalog-list">
				<a href="/gpu-disponible"><div class="cpd-product-card-catalog-list__product-name">RTX 5090 disponible</div></a>
				<div class="cpd-product-card-catalog-list__price"><span class="cp-text--price-total">$83,459.00</span></div>
				<div>26 pzas. disponibles</div>
			</div>
			<div class="cpd-product-card-catalog-list">
				<a href="/gpu-agotada"><div class="cpd-product-card-catalog-list__product-name">RTX 5090 agotada</div></a>
				<div class="cpd-product-card-catalog-list__price"><span class="cp-text--price-total">$69,159.00</span></div>
				<div>Sin piezas disponibles. Crear alerta. Agotado.</div>
			</div>`,
			wantMatch: 2, wantCards: 1, wantPrice: 83459, wantInStock: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(test.html))
			if err != nil {
				t.Fatal(err)
			}
			base, _ := url.Parse(test.store.url)
			result := parseStoreDocument(doc, base, test.store)
			if result.matched != test.wantMatch || result.available != test.wantInStock || len(result.cards) != test.wantCards {
				t.Fatalf("got matched=%d available=%d cards=%d", result.matched, result.available, len(result.cards))
			}
			if result.cards[0].Price != test.wantPrice {
				t.Fatalf("got price %d, want %d", result.cards[0].Price, test.wantPrice)
			}
		})
	}
}

func TestIsInStock(t *testing.T) {
	if isInStock("Sin piezas disponibles. Crear alerta") {
		t.Fatal("unavailable listing was treated as in stock")
	}
	if !isInStock("Solo 2 pzas. disponibles") {
		t.Fatal("available listing was treated as out of stock")
	}
}
