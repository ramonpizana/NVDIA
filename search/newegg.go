package search

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ramonpizana/NVDIA/internal"
	"github.com/ramonpizana/NVDIA/internal/vcard"
)

const defaultMaxPriceMXN = 80000

// These endpoints display prices in MXN and sell domestically in Mexico. DDTech's
// main host blocks GitHub runners, while its NVIDIA catalog host serves the same
// public product catalog without returning HTTP 403.
var mexicoStores = []store{
	{
		name:          "DDTech",
		url:           "https://api.nvidia.ddtech.mx/productos/componentes/tarjetas-de-video?geforce-rtx-serie-50%5B%5D=geforce-rtx-5090&orden=primero-existencia",
		itemSelector:  ".product",
		nameSelector:  "h3.name a",
		priceSelector: ".product-price .price",
		linkSelector:  "h3.name a",
	},
	{
		name:          "Cyberpuerta",
		url:           "https://www.cyberpuerta.mx/Computo-Hardware/Componentes/Tarjetas-de-Video/Filtro/Procesador-grafico/GeForce-RTX-5090/",
		itemSelector:  ".cpd-product-card-catalog-list",
		nameSelector:  ".cpd-product-card-catalog-list__product-name",
		priceSelector: ".cpd-product-card-catalog-list__price .cp-text--price-total",
		linkSelector:  "a",
	},
}

type store struct {
	name, url, itemSelector, nameSelector, priceSelector, linkSelector string
}

type storeResult struct {
	cards     []*vcard.Vcard
	matched   int
	available int
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
	client := &http.Client{Timeout: 30 * time.Second}
	var cards []*vcard.Vcard
	var failures []string
	successfulStores := 0

	for _, s := range mexicoStores {
		result, err := searchStore(client, s)
		if err != nil {
			failures = append(failures, fmt.Sprintf("%s: %v", s.name, err))
			continue
		}
		successfulStores++
		cards = append(cards, result.cards...)
		fmt.Printf("%s: %d disponibles de %d RTX 5090 detectadas\n", s.name, result.available, result.matched)
	}

	if successfulStores == 0 {
		return nil, fmt.Errorf("no se pudo consultar ninguna tienda: %s", strings.Join(failures, "; "))
	}
	if len(failures) > 0 {
		fmt.Printf("advertencia: %s\n", strings.Join(failures, "; "))
	}
	return cards, nil
}

func searchStore(client *http.Client, s store) (storeResult, error) {
	parsed, err := url.Parse(s.url)
	if err != nil {
		return storeResult{}, err
	}
	req, err := http.NewRequest(http.MethodGet, s.url, nil)
	if err != nil {
		return storeResult{}, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/120 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml")
	req.Header.Set("Accept-Language", "es-MX,es;q=0.9")

	response, err := client.Do(req)
	if err != nil {
		return storeResult{}, err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return storeResult{}, fmt.Errorf("HTTP %d %s", response.StatusCode, response.Status)
	}

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return storeResult{}, err
	}
	result := parseStoreDocument(document, parsed, s)
	if result.matched == 0 {
		return storeResult{}, fmt.Errorf("la página respondió, pero no se detectaron productos RTX 5090")
	}
	if result.available > 0 && len(result.cards) == 0 {
		return storeResult{}, fmt.Errorf("se detectaron %d disponibles, pero no se pudo leer su precio", result.available)
	}
	return result, nil
}

func parseStoreDocument(document *goquery.Document, baseURL *url.URL, s store) storeResult {
	result := storeResult{}
	document.Find(s.itemSelector).Each(func(_ int, item *goquery.Selection) {
		name := strings.TrimSpace(item.Find(s.nameSelector).First().Text())
		if !strings.Contains(strings.ToLower(name), "5090") {
			return
		}
		result.matched++
		if !isInStock(strings.ToLower(item.Text())) {
			return
		}
		result.available++

		price, err := PriceConversion(item.Find(s.priceSelector).First().Text())
		if err != nil {
			return
		}
		link, exists := item.Find(s.linkSelector).First().Attr("href")
		if !exists || strings.TrimSpace(link) == "" {
			return
		}
		if absolute, joinErr := baseURL.Parse(link); joinErr == nil {
			link = absolute.String()
		}
		result.cards = append(result.cards, vcard.New(name, s.name, price, link, true))
	})
	return result
}

func isInStock(listingText string) bool {
	listingText = strings.ToLower(listingText)
	for _, unavailable := range []string{
		"sin stock", "sin existencia", "sin piezas", "agotado", "out of stock", "no disponible", "crear alerta",
	} {
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
