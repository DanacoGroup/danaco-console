// Odpowiedzialność pliku: wypełnienie portu KatalogTozsamosci repozytorium
// warstwy danych. Port mówi typami kontraktu, repozytorium wierszami tabel —
// tutaj i tylko tutaj jedno przechodzi w drugie.
package core

import (
	"context"
	"strconv"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// zrodloTozsamosci wypełnia port danymi z tabel `kategoria_tozsamosci`
// i `dokument_tozsamosci`.
type zrodloTozsamosci struct {
	repozytorium dane.RepozytoriumTozsamosci
}

// NoweZrodloTozsamosci wiąże katalog tożsamości z repozytorium warstwy danych.
// Zwraca port, bo składacz ma widzieć wyłącznie kontrakt, nigdy wiersz tabeli.
func NoweZrodloTozsamosci(repozytorium dane.RepozytoriumTozsamosci) KatalogTozsamosci {
	return &zrodloTozsamosci{repozytorium: repozytorium}
}

// Kategorie zwraca katalog kategorii przełożony na struktury kontraktu.
func (z *zrodloTozsamosci) Kategorie(ctx context.Context, tylkoAktywne bool) ([]shared.IdentityCategory, error) {
	if z == nil || z.repozytorium == nil {
		return []shared.IdentityCategory{}, nil
	}
	wiersze, err := z.repozytorium.Kategorie(ctx, tylkoAktywne)
	if err != nil {
		return nil, err
	}
	kategorie := make([]shared.IdentityCategory, 0, len(wiersze))
	for _, wiersz := range wiersze {
		kategorie = append(kategorie, kategoriaKontraktu(wiersz))
	}
	return kategorie, nil
}

// Dokumenty zwraca treści jednej osi. Odczyt obejmuje wyłącznie zapisy czynne —
// zapis nieczynny przepuszcza oś szerszą, tak samo jak brak wiersza.
func (z *zrodloTozsamosci) Dokumenty(ctx context.Context, os shared.ConfigAxis,
	bytOsi string) ([]shared.IdentityDocument, error) {

	if z == nil || z.repozytorium == nil {
		return []shared.IdentityDocument{}, nil
	}
	wiersze, err := z.repozytorium.Dokumenty(ctx, dane.FiltrTozsamosci{
		Os: os, OsByt: bytOsi, TylkoAktywne: true,
	})
	if err != nil {
		return nil, err
	}
	dokumenty := make([]shared.IdentityDocument, 0, len(wiersze))
	for _, wiersz := range wiersze {
		dokumenty = append(dokumenty, dokumentKontraktu(wiersz))
	}
	return dokumenty, nil
}

// kategoriaKontraktu przekłada wiersz katalogu na strukturę kontraktu.
// Identyfikatorem kategorii jest jej kod — trwały i czytelny dla okna
// konfiguracji, w przeciwieństwie do numeru wiersza.
func kategoriaKontraktu(wiersz dane.KategoriaTozsamosci) shared.IdentityCategory {
	kategoria := shared.IdentityCategory{
		Id:          wiersz.Kod,
		Name:        wiersz.Nazwa,
		Layer:       wiersz.Warstwa,
		Order:       wiersz.Kolejnosc,
		Required:    wiersz.Obowiazkowa,
		DefaultMode: wiersz.TrybDomyslny,
		Enabled:     wiersz.Aktywna,
	}
	if wiersz.Opis != "" {
		opis := wiersz.Opis
		kategoria.Description = &opis
	}
	return kategoria
}

// dokumentKontraktu przekłada wiersz treści na strukturę kontraktu.
func dokumentKontraktu(wiersz dane.DokumentTozsamosci) shared.IdentityDocument {
	dokument := shared.IdentityDocument{
		Id:         strconv.FormatInt(wiersz.ID, 10),
		CategoryId: wiersz.KodKategorii,
		Axis:       wiersz.Os,
		Mode:       wiersz.Tryb,
		Content:    wiersz.Tresc,
		Enabled:    wiersz.Aktywny,
		UpdatedAt:  chwilaBazy(wiersz.Zaktualizowano),
	}
	if wiersz.OsByt != "" {
		byt := wiersz.OsByt
		dokument.AxisId = &byt
	}
	if wiersz.OdciskTresci != "" {
		odcisk := wiersz.OdciskTresci
		dokument.ContentHash = &odcisk
	}
	return dokument
}
