// Odpowiedzialność pliku: zapis pamięci tłumaczeń modułu Translate — jedyne
// miejsce, w którym para (segment źródłowy, segment przekładu) trafia do tabeli
// `pamiec_tlumaczen`. Metody stoją na wspólnym `*adapterTlumaczenia`, którego
// typ, konstruktor i przedrostki deklaruje `adapter_modul_tlumaczenie.go`.
//
// O tym, kiedy para jest zatwierdzona, rozstrzyga wywołujący. Chwile są dwie,
// obie po treści, która realnie powstała:
//
//  1. `target.add` — model oddał przekład i panel został zapisany. Pary idą do
//     pamięci po zapisie panelu, bo dopiero wtedy istnieje `panel_id`, którego
//     wymaga klucz obcy tabeli.
//  2. `translation.set` — Operator poprawił przekład, więc parę zatwierdził
//     człowiek. Pominięcie jej znaczyłoby, że pamięć pamięta wyłącznie model,
//     a poprawki Operatora zapomina.
//
// Chwilą trzecią nie jest `backtranslation.run`: tłumaczenie zwrotne jest
// kontrolą, a nie przekładem do ponownego użycia, i zanieczyściłoby pamięć
// parami w odwrotną stronę.
//
// Sparowanie segmentów jest częścią trudną. Rdzeń ma jedno narzędzie podziału —
// `podzielNaZdania` (`adapter_modul_tlumaczenie.go`), podział mechaniczny po
// znakach końca zdania, z opisanym tam ograniczeniem. Model tłumaczący nie ma
// obowiązku zachować liczby zdań: scala dwa zdania w jedno, rozbija jedno na
// dwa albo dokłada zdanie wyjaśniające. Stąd rozstrzygnięcie ostrożne:
//
//   - liczba segmentów po obu stronach równa i większa od zera → parowanie
//     kolejne, segment po segmencie. Podział jest wtedy ten sam po obu stronach,
//     a kolejność zdań w przekładzie zachowana;
//   - liczby różne albo któraś strona bez segmentów → jedna para całościowa:
//     cały tekst źródłowy do całego przekładu. Dopasuje się rzadziej, bo
//     `LIKE %fraza%` po długim segmencie trafia tylko przy długim zapytaniu, ale
//     jest prawdziwa. Rozdzielanie n zdań źródła na m zdań przekładu „mniej
//     więcej po kolei" wytworzyłoby pary, w których segment docelowy nie jest
//     przekładem segmentu źródłowego, a pamięć podpowiadałaby je jako gotowe
//     zdania.
//
// Miary podobieństwa, która pozwoliłaby parować przy różnych liczbach (na
// przykład wyrównania długościami albo przez model), rdzeń nie ma i ten plik
// jej nie dorabia.
//
// Błąd zapisu pary jest pomijany, tak samo jak błąd zapisu migawki jakości
// w `DodajPanel`: przekład już powstał i jest w ręku Operatora, więc pamięć jest
// wzbogaceniem, a nie warunkiem przekładu.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/dane"
)

// zapamietajPary dokłada do pamięci tłumaczeń pary segmentów wyprowadzone
// z jednego udanego przekładu. Rozstrzygnięcie o sparowaniu — nagłówek pliku.
//
// Język pary to język panelu, nie język źródłowy: `Podpowiedzi` szuka po
// kolumnie `jezyk` i pyta, co wpisać w języku docelowym dla danego segmentu,
// więc kluczem jest język, na który tłumaczono. Zapisanie tam języka źródłowego
// uczyniłoby całą tabelę nieosiągalną dla odczytu.
//
// Panel bez identyfikatora wewnętrznego (a więc niezapisany) nie ma do czego
// przypiąć wiersza — wychodzimy bez zapisu, zamiast pozwolić kluczowi obcemu
// bazy odrzucić to dopiero przy wstawieniu.
func (a *adapterTlumaczenia) zapamietajPary(ctx context.Context,
	panel dane.PanelTlumaczenia, tekstZrodlowy, przeklad string) {

	if panel.ID == 0 || strings.TrimSpace(panel.Jezyk) == "" {
		return
	}
	for _, para := range sparujSegmenty(tekstZrodlowy, przeklad) {
		wpis := dane.WpisPamieciTlumaczen{
			Kod:             nowyIdentyfikator(przedrostekPamieciTlumaczen),
			Jezyk:           panel.Jezyk,
			SegmentZrodlowy: para.zrodlo,
			SegmentDocelowy: para.cel,
		}
		// Wynik zapisu pomijany świadomie — uzasadnienie w nagłówku pliku.
		_, _ = a.repozytorium.ZapiszPamiec(ctx, panel.ID, wpis)
	}
}

// paraSegmentow to jeden wiersz przyszłej pamięci, zanim dostanie tożsamość.
type paraSegmentow struct {
	zrodlo string
	cel    string
}

// sparujSegmenty wykonuje jedyne dopasowanie, które ten rdzeń potrafi obronić.
// Oddaje pustkę, gdy którakolwiek strona jest pusta — para z pustym segmentem
// źródłowym i tak zostałaby odrzucona przez `ZapiszPamiec`, a para z pustym
// przekładem nie podpowiada niczego.
func sparujSegmenty(tekstZrodlowy, przeklad string) []paraSegmentow {
	zrodlo := strings.TrimSpace(tekstZrodlowy)
	cel := strings.TrimSpace(przeklad)
	if zrodlo == "" || cel == "" {
		return nil
	}

	segmentyZrodla := podzielNaZdania(zrodlo)
	segmentyCelu := podzielNaZdania(cel)

	// Ten sam podział po obu stronach — parowanie kolejne.
	if len(segmentyZrodla) > 0 && len(segmentyZrodla) == len(segmentyCelu) {
		pary := make([]paraSegmentow, 0, len(segmentyZrodla))
		for i := range segmentyZrodla {
			// Segment pusty po którejkolwiek stronie pomijamy zamiast zapisywać —
			// `podzielNaZdania` odsiewa puste odcinki, więc to sytuacja skrajna,
			// ale wiersz z pustym segmentem docelowym byłby podpowiedzią „nic".
			if segmentyZrodla[i] == "" || segmentyCelu[i] == "" {
				continue
			}
			pary = append(pary, paraSegmentow{zrodlo: segmentyZrodla[i], cel: segmentyCelu[i]})
		}
		return pary
	}

	// Liczby się rozjechały, więc dopasowania zdanie-do-zdania nie da się
	// obronić. Jedna para całościowa — prawdziwa, choć grubsza.
	return []paraSegmentow{{zrodlo: zrodlo, cel: cel}}
}
