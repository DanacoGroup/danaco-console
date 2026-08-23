package core

import (
	"context"

	"danacoconsole/server/internal/dane"
)

// CiagloscRozmowy utrwala identyfikator rozmowy nadany przez program `claude`,
// dzięki któremu kolejna tura wznawia rozmowę zamiast zaczynać od zera.
//
// Port jest celowo wąski — dwie czynności na jednej kolumnie. Rozmowa należy do
// okna, nie do sesji: dwa okna jednej sesji prowadzą dwie niezależne rozmowy
// z modelem i muszą mieć osobne wznowienia.
type CiagloscRozmowy interface {
	// Zapamietaj utrwala identyfikator rozmowy przy oknie.
	Zapamietaj(kontekst context.Context, idOkna, idRozmowy string) error
	// Przypomnij zwraca identyfikator rozmowy okna albo pusty napis, gdy okno
	// nie rozmawiało jeszcze z modelem. Pusty napis jest stanem poprawnym i
	// znaczy „zacznij nową rozmowę".
	Przypomnij(kontekst context.Context, idOkna string) string
}

// ciagloscNadOknami realizuje port nad repozytorium okien.
type ciagloscNadOknami struct {
	okna dane.RepozytoriumOkien
}

// nowaCiagloscRozmowy składa port nad repozytorium okien. Repozytorium puste
// daje port pusty — rozmowa nadal działa, tylko bez ciągłości.
func nowaCiagloscRozmowy(okna dane.RepozytoriumOkien) CiagloscRozmowy {
	if okna == nil {
		return nil
	}
	return &ciagloscNadOknami{okna: okna}
}

// ciagloscRozmowy wyjmuje repozytorium okien z zestawu i składa port. Zestaw
// pusty daje port pusty — rdzeń wstaje bez bazy i prowadzi rozmowę bez
// ciągłości, zamiast odmówić startu.
func ciagloscRozmowy(repozytoria *dane.Zestaw) CiagloscRozmowy {
	if repozytoria == nil || repozytoria.Okna == nil {
		return nil
	}
	return nowaCiagloscRozmowy(repozytoria.Okna)
}

func (c *ciagloscNadOknami) Zapamietaj(kontekst context.Context, idOkna, idRozmowy string) error {
	wiersz, err := c.okna.PoIdentyfikatorze(kontekst, idOkna)
	if err != nil {
		return err
	}
	return c.okna.ZapiszRozmoweCLI(kontekst, wiersz.ID, idRozmowy)
}

func (c *ciagloscNadOknami) Przypomnij(kontekst context.Context, idOkna string) string {
	wiersz, err := c.okna.PoIdentyfikatorze(kontekst, idOkna)
	if err != nil {
		return ""
	}
	idRozmowy, err := c.okna.RozmowaCLI(kontekst, wiersz.ID)
	if err != nil {
		return ""
	}
	return idRozmowy
}
