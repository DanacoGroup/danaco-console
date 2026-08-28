// Odpowiedzialność pliku: zapis pamięci tłumaczeń modułu Translate — jedyne
// miejsce, w którym para segmentów źródła i przekładu trafia do tabeli
// `pamiec_tlumaczen`. O zatwierdzeniu pary rozstrzyga wywołujący, po treści,
// która realnie powstała.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/dane"
)

// zapamietajPary dokłada do pamięci tłumaczeń pary segmentów wyprowadzone z
// jednego udanego przekładu. Język pary to język panelu, na który tłumaczono,
// nie język źródłowy. Panel bez identyfikatora wewnętrznego nie ma do czego
// przypiąć wiersza.
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
		// Wynik zapisu jest pomijany świadomie: przekład już powstał, pamięć jest
		// wzbogaceniem, nie warunkiem.
		_, _ = a.repozytorium.ZapiszPamiec(ctx, panel.ID, wpis)
	}
}

// paraSegmentow to jeden wiersz przyszłej pamięci tłumaczeń, zanim otrzyma
// własną tożsamość w bazie danych.
type paraSegmentow struct {
	zrodlo string
	cel    string
}

// sparujSegmenty wykonuje jedyne dopasowanie, które ten rdzeń potrafi obronić.
// Oddaje pustkę, gdy którakolwiek strona jest pusta.
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
			// Segment pusty po którejkolwiek stronie zostaje pominięty zamiast
			// zapisany.
			if segmentyZrodla[i] == "" || segmentyCelu[i] == "" {
				continue
			}
			pary = append(pary, paraSegmentow{zrodlo: segmentyZrodla[i], cel: segmentyCelu[i]})
		}
		return pary
	}

	// Liczby się rozjechały, dopasowania zdanie-do-zdania nie da się obronić.
	return []paraSegmentow{{zrodlo: zrodlo, cel: cel}}
}
