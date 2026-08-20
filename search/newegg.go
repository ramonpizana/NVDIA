package search

import (
	"fmt"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/ramonpizana/NVDIA/internal"
	"github.com/ramonpizana/NVDIA/internal/vcard"
)

const defaultMaxPriceMXN = 80000

// Mexican retailers are used so displayed prices are in MXN and delivery is domestic.
// The monitor excludes foreign stores: duties and courier charges are not reliably
// known until checkout and therefore cannot be promised as "no extra".
var mexicoStores = []store{
	{"DDTech", "https://ddtech.mx/productos/componentes/tarjetas-de-video?orden=menor-mayor&stock=con-existencia", ".product-item", ".product-item-link, .product-item-name", ".price, .product-price", "a"},
	{"Cyberpuerta", "https://www.cyberpuerta.mx/Computo-Hardware/Componentes/Tarjetas-de-Video/Filtro/Procesador-grafico/GeForce-RTX-5090/", ".product__item, .productBox", ".product__name, .productTitle", ".product__price, .price", "a"},
}

type store struct {
	name, url, itemSelector, nameSelector, priceSelector, linkSelector string
}

// SearchNewEgg retains the old exported name for callers. It now searches RTX 5090
// listings at Mexican retailers, so alerts only contain prices quoted in MXN.
func SearchNewEgg(pass string) error {
	cards, err := searchMexicoStores()
	if err != nil {
		return err
	}

	maxPrice := maxPriceMXN()
	filtered := check(cards, maxPrice)
	sort.Slice(filtered, func(i, j int) bool { return filtered[i].Price < filtered[j].Price })
	fmt.Printf("encontré %d RTX 5090 disponibles por hasta $%d MXN\n", len(filtered), maxPrice)
	if len(filtered) == 0 {
		return nil
	}
	return internal.MailInfo(pass, filtered)
}

func searchMexicoStores() ([]*vcard.Vcard, error) {
	var cards []*vcard.Vcard
	var failures []string
	for _, s := range mexicoStores {
		found, err := searchStore(s)
		if err != nil {
			failures = append(failures, fmt.Sprintf("%s: %v", s.name, err))
			continue
		}
		cards = append(cards, found...)
	}
	if len(cards) == 0 && len(failures) > 0 {
		return nil, fmt.Errorf("no se pudo consultar ninguna tienda: %s", strings.Join(failures, "; "))
	}
	if len(failures) > 0 {
		fmt.Printf("advertencia: %s\n", strings.Join(failures, "; "))
	}
	return cards, nil
}

func searchStore(s store) ([]*vcard.Vcard, error) {
	parsed, err := url.Parse(s.url)
	if err != nil {
		return nil, err
	}
	var cards []*vcard.Vcard
	c := colly.NewCollector(colly.AllowedDomains(parsed.Host))
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/120 Safari/537.36")
	})
	c.OnHTML(s.itemSelector, func(h *colly.HTMLElement) {
		name := strings.TrimSpace(h.ChildText(s.nameSelector))
		listingText := strings.ToLower(h.Text)
		if !strings.Contains(strings.ToLower(name), "5090") || !isInStock(listingText) {
			return
		}
		price, conversionErr := PriceConversion(h.ChildText(s.priceSelector))
		if conversionErr != nil {
			return
		}
		link := h.ChildAttr(s.linkSelector, "href")
		if absolute, joinErr := parsed.Parse(link); joinErr == nil {
			link = absolute.String()
		}
		cards = append(cards, vcard.New(name, s.name, price, link, true))
	})
	if err := c.Visit(s.url); err != nil {
		return nil, err
	}
	return cards, nil
}

func isInStock(listingText string) bool {
	for _, unavailable := range []string{"sin stock", "agotado", "out of stock", "no disponible"} {
		if strings.Contains(listingText, unavailable) {
			return false
		}
	}
	return true
}

func maxPriceMXN() int {
	if value, err := strconv.Atoi(os.Getenv("MAX_PRICE_MXN")); err == nil && value > 0 {
		return value
	}
	return defaultMaxPriceMXN
}

func check(cards []*vcard.Vcard, maxPrice int) []*vcard.Vcard {
	var send []*vcard.Vcard
	for _, card := range cards {
		if card.Price <= maxPrice && card.Stock {
			send = append(send, card)
		}
	}
	return send
}
