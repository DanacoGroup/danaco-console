package core

import (
	"testing"

	"danacoconsole/shared"
)

// Sprawdzian pokrycia ZAPISU AUTORA na wielu drogach zmiany dokumentu.
//
// ── Dlaczego wykaz dróg, a nie jedna droga ──────────────────────────────────
// Przełącznik „pokaż wszystko, co zrobił model" stoi na założeniu, że KAŻDA
// droga zmiany zostawia ślad podpisany wykonawcą. Sprawdzian jednej drogi
// dowodzi jednej drogi. Ten wykaz zmierzył pięć i na jednej znalazł dziurę:
// `studio.document.save` zawołane przez wykonawcę zmieniało treść i NIE
// odkładało ani zmiany śledzonej, ani wpisu dziennika. Dziura została zamknięta
// siatką `sladWykonawcyStudia`; sprawdzian zostaje, żeby nie wróciła i żeby
// każda nowa droga trafiła tu przed odbiorem.
//
// Ręka modelu bierze się z GNIAZDA serwera narzędzi, a pola `author` sprawdzian
// celowo NIE podaje — mierzy to, czego model nie może o sobie zataić.
func TestProbaAutorNaWieluDrogach(t *testing.T) {
	drogi := []struct {
		nazwa   shared.MessageType
		ladunek func(dok string) any
	}{
		{shared.CommandStudioTextEdit, func(d string) any {
			return shared.StudioTextEditRequest{DocumentId: d, RangeStart: 0, RangeEnd: 6, Text: "Zmiana"}
		}},
		{shared.CommandStudioFormatCharacterSet, func(d string) any {
			od, do, pogrub := 0, 6, true
			return shared.StudioFormatCharacterSetRequest{DocumentId: d,
				RangeStart: &od, RangeEnd: &do, Bold: &pogrub}
		}},
		{shared.CommandStudioFormatParagraphSet, func(d string) any {
			od, do := 0, 6
			wyr := shared.StudioTextAlign(shared.StudioTextAlignCenter)
			return shared.StudioFormatParagraphSetRequest{DocumentId: d,
				RangeStart: &od, RangeEnd: &do, Align: &wyr}
		}},
		{shared.CommandStudioDocumentSave, func(d string) any {
			return shared.StudioDocumentSaveRequest{DocumentId: d, Content: "Treść od modelu."}
		}},
		{shared.CommandStudioMarkupAdd, func(d string) any {
			barwa := "#f5a623"
			return shared.StudioMarkupAddRequest{DocumentId: d,
				Kind:       shared.StudioMarkupKind(shared.StudioMarkupKindHighlight),
				RangeStart: 0, RangeEnd: 6, Color: &barwa}
		}},
	}

	for _, droga := range drogi {
		t.Run(string(droga.nazwa), func(t *testing.T) {
			u := blokadaZmontuj(t)
			odp := u.blokadaWykonajJakoModel(t, droga.nazwa, droga.ladunek(u.dokument))
			if odp.Error != nil {
				t.Logf("WYNIK %s: odmowa %s", droga.nazwa, odp.Error.Code)
				return
			}
			var zmiany, czynnosci int
			_ = u.baza.QueryRow(`SELECT COUNT(*) FROM zmiana_sledzona_studio z
				JOIN dokument_studio d ON d.id=z.dokument_id
				WHERE d.identyfikator_zewnetrzny=? AND z.autor='model'`, u.dokument).Scan(&zmiany)
			_ = u.baza.QueryRow(`SELECT COUNT(*) FROM czynnosc_dokumentu_studio c
				JOIN dokument_studio d ON d.id=c.dokument_id
				WHERE d.identyfikator_zewnetrzny=? AND c.autor_rodzaj='model'`, u.dokument).Scan(&czynnosci)
			var znak int
			_ = u.baza.QueryRow(`SELECT COUNT(*) FROM znakowanie_studio z
				JOIN dokument_studio d ON d.id=z.dokument_id
				WHERE d.identyfikator_zewnetrzny=? AND z.autor_rodzaj='model'`, u.dokument).Scan(&znak)
			t.Logf("WYNIK %s: zmiany_sledzone(model)=%d czynnosci(model)=%d znakowania(model)=%d",
				droga.nazwa, zmiany, czynnosci, znak)
			if zmiany == 0 && czynnosci == 0 && znak == 0 {
				t.Errorf("DZIURA: %s nie zapisała autora na żadnej z trzech dróg śladu", droga.nazwa)
			}
		})
	}
}
