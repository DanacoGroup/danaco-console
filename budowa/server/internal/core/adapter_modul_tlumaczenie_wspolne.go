// Plik obsługuje drobne narzędzia wspólne dobudowanym rodzinom modułu Translate: odczyt pól nieobowiązkowych żądania, odmowa o pliku, miara podobieństwa segmentów i trwały podział okna, potrzebne w ośmiu plikach modułu.
package core

import (
	"context"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// napisZeWskaznika oddaje treść pola nieobowiązkowego albo pusty napis, gdy wskaźnik żądania jest pusty.
func napisZeWskaznika(wartosc *string) string {
	if wartosc == nil {
		return ""
	}
	return *wartosc
}

// liczbaCalkowitaZeWskaznika oddaje wartość pola nieobowiązkowego albo zero, gdy wskaźnik żądania jest pusty.
func liczbaCalkowitaZeWskaznika(wartosc *int) int {
	if wartosc == nil {
		return 0
	}
	return *wartosc
}

// wskaznikNapisu oddaje wskaźnik na treść albo nic, gdy treść jest pusta —
// pusty napis w polu nieobowiązkowym kontraktu udawałby wartość podaną.
func wskaznikNapisu(wartosc string) *string {
	if strings.TrimSpace(wartosc) == "" {
		return nil
	}
	return &wartosc
}

// bladPlikuTlumaczenia nazywa niepowodzenie pracy na pliku wskazanym przez Operatora. Plik nieobecny albo niedostępny jest pomyłką wskazania, nie usterką rdzenia, więc kod odmowy nie zachęca do pętli ponowień.
func bladPlikuTlumaczenia(sciezka string, err error) error {
	if os.IsNotExist(err) {
		return bladWskazaniaTlumaczenia("pliku " + sciezka + " nie ma pod wskazaną ścieżką")
	}
	return bladWskazaniaTlumaczenia("nie można sięgnąć po plik " + sciezka + ": " + err.Error())
}

// podobienstwoSegmentow mierzy bliskość dwóch segmentów w procentach: sto znaczy identyczne, zero — rozbieżne całkowicie. Miarą jest odległość edycyjna (Levenshtein) odniesiona do długości segmentu dłuższego.
func podobienstwoSegmentow(pierwszy, drugi string) int {
	a := []rune(strings.ToLower(strings.Join(strings.Fields(pierwszy), " ")))
	b := []rune(strings.ToLower(strings.Join(strings.Fields(drugi), " ")))
	if len(a) == 0 && len(b) == 0 {
		return 100
	}
	if len(a) == 0 || len(b) == 0 {
		return 0
	}
	// Wiersz poprzedni i bieżący wystarczą, bo nikt tu nie odtwarza drogi edycji.
	poprzedni := make([]int, len(b)+1)
	biezacy := make([]int, len(b)+1)
	for j := range poprzedni {
		poprzedni[j] = j
	}
	for i := 1; i <= len(a); i++ {
		biezacy[0] = i
		for j := 1; j <= len(b); j++ {
			koszt := 1
			if a[i-1] == b[j-1] {
				koszt = 0
			}
			biezacy[j] = min(min(biezacy[j-1]+1, poprzedni[j]+1), poprzedni[j-1]+koszt)
		}
		poprzedni, biezacy = biezacy, poprzedni
	}
	odleglosc := poprzedni[len(b)]
	dluzszy := max(len(a), len(b))
	return 100 * (dluzszy - odleglosc) / dluzszy
}

// segmentyOkna oddaje podział okna: trwały, gdy Operator scalał albo dzielił segmenty, a przy jego braku — podział wyprowadzony z tekstu źródłowego. Branie podziału mechanicznego mimo istnienia trwałego cofałoby każde scalenie przy kolejnym wywołaniu.
func (a *adapterTlumaczenia) segmentyOkna(ctx context.Context,
	okno dane.OknoTlumaczenia) ([]string, error) {

	trwale, err := a.repozytorium.SegmentyOkna(ctx, okno.ID)
	if err != nil {
		return nil, bladTlumaczenia(err)
	}
	if len(trwale) > 0 {
		return trwale, nil
	}
	if okno.TekstZrodlowy == nil {
		return nil, nil
	}
	return podzielNaZdania(*okno.TekstZrodlowy), nil
}

// trescPanelu oddaje treść panelu albo odmowę, gdy panel jest pusty. Czynność
// na pustym panelu (korekta, napisy, dubbing, wydanie) nie ma materiału, a wynik
// pusty nie jest wynikiem.
func trescPanelu(panel dane.PanelTlumaczenia) (string, error) {
	if panel.Tresc == nil || strings.TrimSpace(*panel.Tresc) == "" {
		return "", bladWskazaniaTlumaczenia("panel " + panel.Kod +
			" nie ma treści przekładu — nie ma na czym pracować")
	}
	return *panel.Tresc, nil
}

// liczbaZnakow liczy znaki, nie bajty — kontrola długości linii napisów
// i czytelności ma mierzyć to, co widzi czytelnik.
func liczbaZnakow(tekst string) int {
	return utf8.RuneCountInString(tekst)
}

// liczbaZeWskazania czyta numer podany napisem. Kontrakt niesie numery
// segmentów wykazem napisów, więc odczyt jest tu, a nie w każdej komendzie
// z osobna.
func liczbaZeWskazania(wskazanie string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(wskazanie))
}

// wagaOdmowyKontroli sprowadza wagę zapisaną w bazie do wartości kontraktu, odczytaną z tabeli kontroli jakości.
func wagaOdmowyKontroli(waga string) shared.ProofreadSeverity {
	return shared.ProofreadSeverity(waga)
}
