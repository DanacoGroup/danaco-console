// Odpowiedzialność pliku: moduł Translate — obsługa `speech.synthesize`
// i `panel.export` na `*adapterTlumaczenia`. Typ, konstruktor i przedrostki
// deklaruje `adapter_modul_tlumaczenie.go`; port `Tlumaczenie` i rejestracja
// komend stoją osobno — ten plik dokłada tylko metody.
//
// Mowę syntezuje `espeak-ng`, syntezator lokalny uruchamiany portem
// `session.Uruchamiacz` — tą samą drogą, którą chodzi silnik rozpoznawania mowy
// (pakiet `server/internal/mowa`). Wybór syntezatora, jego cena (głos brzydki)
// i powód odrzucenia `pipera` — nagłówek
// `adapter_modul_tlumaczenie_mowa_silnik.go`.
//
// Kolumna `nagranie_odnosnik` (`migracja_055_jakosc_i_mowa.sql`, tabela
// `panel_tlumaczenia_synteza_mowy`) niesie ścieżkę pliku, który naprawdę
// powstał, i tę samą ścieżkę oddaje pole `Path` odpowiedzi. Ślad zapisujemy
// dopiero po syntezie — wiersz z odnośnikiem do nagrania, którego nie ma, mówiłby
// nieprawdę.
//
// Brak syntezatora jest odmową, nie atrapą. Gdy programu nie ma na maszynie, gdy
// rdzeń nie ma uruchamiacza albo gdy syntezator nie zna głosu dla języka panelu
// — komenda odmawia, nazywając brak i wskazując naprawę, zamiast oddać pustą
// ścieżkę udającą nagranie.
//
// Rdzeń nie ma magazynu blobów (ten sam brak, co w Library, Research i przy
// słowniku tego samego modułu). `panel.export` zapisuje ślad — format i czas
// w tabeli `panel_tlumaczenia_eksport` — ale nie wytwarza pliku na dysku.
// Kolumna `plik_odnosnik` wraca NULL, a pole `Path` odpowiedzi puste. `Format`
// bierze wartość kontraktu wprost (`pdf`, `docx`, `markdown`, `html`, `txt`),
// bez tłumaczenia wartości.
package core

import (
	"context"
	"errors"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// SyntezujMowe obsługuje `speech.synthesize`. Czyta treść panelu, oddaje ją
// syntezatorowi i zwraca ścieżkę pliku dźwiękowego, który powstał.
//
// Trzy odmowy wprost, każda o czym innym:
//  1. żądanie bez panelu — nie ma czego odsłuchać;
//  2. panel bez treści — odsłuch pustki dałby nagranie ciszy udające przeczytany
//     przekład; brak treści jest tu wiadomością, a nie plikiem do wytworzenia;
//  3. panel bez języka — syntezator dostaje głos z pola `jezyk` panelu i rdzeń
//     nie podstawia za nie własnego domyślnego (nagłówek silnika: „głosu nie
//     zgadujemy”).
//
// Odmowy samego silnika (brak programu, brak głosu dla języka, izolacja)
// przychodzą z `zsyntezujDoPliku` już oznakowane kodem kontraktu.
func (a *adapterTlumaczenia) SyntezujMowe(ctx context.Context,
	z shared.TranslateSpeechSynthesizeRequest) (shared.TranslateSpeechSynthesizeResponse, error) {

	if z.PanelId == "" {
		return shared.TranslateSpeechSynthesizeResponse{}, bladWskazaniaTlumaczenia("speech.synthesize bez panelu")
	}

	panel, err := a.repozytorium.Panel(ctx, z.PanelId)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.TranslateSpeechSynthesizeResponse{}, bladNieznanegoPanelu(z.PanelId, err)
	}
	if err != nil {
		return shared.TranslateSpeechSynthesizeResponse{}, bladTlumaczenia(err)
	}

	tresc := ""
	if panel.Tresc != nil {
		tresc = strings.TrimSpace(*panel.Tresc)
	}
	if tresc == "" {
		return shared.TranslateSpeechSynthesizeResponse{}, bladWskazaniaTlumaczenia(
			"panel " + z.PanelId + " nie ma treści — nie ma czego przeczytać na głos")
	}
	jezyk := strings.TrimSpace(panel.Jezyk)
	if jezyk == "" {
		return shared.TranslateSpeechSynthesizeResponse{}, bladWskazaniaTlumaczenia(
			"panel " + z.PanelId + " nie ma języka — syntezator nie wie, którym głosem czytać")
	}

	sciezka, err := a.zsyntezujDoPliku(ctx, panel.Kod, jezyk, tresc)
	if err != nil {
		return shared.TranslateSpeechSynthesizeResponse{}, err
	}

	// Ślad zapisujemy po syntezie, z odnośnikiem do pliku, który istnieje.
	// Nieudany zapis śladu nie przewraca komendy: nagranie już powstało, a jego
	// ścieżka jest dla Operatora wartościowsza niż wiersz historii, który można
	// powtórzyć. Ta sama zasada, co przy migawce jakości w `DodajPanel`.
	_, _ = a.repozytorium.ZapiszSyntezeMowy(ctx, dane.SyntezaMowy{PanelID: panel.ID, NagranieOdnosnik: &sciezka})

	return shared.TranslateSpeechSynthesizeResponse{PanelId: z.PanelId, Path: sciezka}, nil
}

// EksportujPanel obsługuje `panel.export`. Zapisuje ślad eksportu panelu w
// formacie żądanym przez Operatora — `Path` odpowiedzi zostaje pusty, bo
// rdzeń nie ma magazynu blobów (patrz nagłówek pliku).
func (a *adapterTlumaczenia) EksportujPanel(ctx context.Context,
	z shared.TranslatePanelExportRequest) (shared.TranslatePanelExportResponse, error) {

	if z.PanelId == "" {
		return shared.TranslatePanelExportResponse{}, bladWskazaniaTlumaczenia("panel.export bez panelu")
	}

	panel, err := a.repozytorium.Panel(ctx, z.PanelId)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.TranslatePanelExportResponse{}, bladNieznanegoPanelu(z.PanelId, err)
	}
	if err != nil {
		return shared.TranslatePanelExportResponse{}, bladTlumaczenia(err)
	}

	eksport := dane.EksportPanelu{PanelID: panel.ID, Format: string(z.Format)}
	if _, err := a.repozytorium.ZapiszEksportPanelu(ctx, eksport); err != nil {
		return shared.TranslatePanelExportResponse{}, bladTlumaczenia(err)
	}

	return shared.TranslatePanelExportResponse{PanelId: z.PanelId}, nil
}
