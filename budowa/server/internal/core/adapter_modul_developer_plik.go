// Odpowiedzialność pliku: developer.file.open i developer.file.save — warstwa
// plikowa okna Code Editor. Plik binarny wraca bez treści, a nie jako błąd:
// DeveloperFile.content jest w kontrakcie polem opcjonalnym.
package core

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// granicaTresciEdytora jest największym plikiem wpuszczanym do Code Editor.
// Powyżej niej odczyt przestaje być otwarciem pliku, a staje się przesłaniem
// paczki danych przez gniazdo zdarzeń.
const granicaTresciEdytora = 4 << 20

// OtworzPlik obsługuje developer.file.open — wczytuje wskazany plik do Code Editor tego okna Operatora.
func (a *adapterDevelopera) OtworzPlik(ctx context.Context,
	z shared.DeveloperFileOpenRequest) (shared.DeveloperFileOpenResponse, error) {

	okno, err := a.oknoDevelopera(z.WindowId)
	if err != nil {
		return shared.DeveloperFileOpenResponse{}, err
	}
	sciezka, err := sciezkaWObszarze(a.korzenieOkna(okno), z.Path, czyIstnieje)
	if err != nil {
		return shared.DeveloperFileOpenResponse{}, err
	}

	opis, err := os.Stat(sciezka)
	if errors.Is(err, os.ErrNotExist) {
		return shared.DeveloperFileOpenResponse{}, bladZasobuDevelopera("plik " + sciezka + " nie istnieje")
	}
	if err != nil {
		return shared.DeveloperFileOpenResponse{}, bladWykonaniaDevelopera(
			"nie można odczytać opisu pliku " + sciezka + ": " + err.Error())
	}
	if opis.IsDir() {
		return shared.DeveloperFileOpenResponse{}, bladZadaniaDevelopera(
			sciezka + " jest katalogiem — Code Editor otwiera plik, drzewo katalogów daje developer.tree.get")
	}

	plik := opisPliku(sciezka, opis)
	if opis.Size() <= granicaTresciEdytora {
		bajty, err := os.ReadFile(sciezka)
		if err != nil {
			return shared.DeveloperFileOpenResponse{}, bladWykonaniaDevelopera(
				"nie można odczytać pliku " + sciezka + ": " + err.Error())
		}
		if !czyBinarna(bajty) {
			tresc := string(bajty)
			plik.Content = &tresc
		} else {
			plik.Language = wskaznikTekstu(jezykBinarny)
		}
	}
	plik.VersionId = a.ostatniaWersja(ctx, okno.Id, sciezka)
	return shared.DeveloperFileOpenResponse{File: plik}, nil
}

// ZapiszPlik obsługuje developer.file.save — zapisuje treść pliku po edycji w Code Editor tego okna Operatora.
func (a *adapterDevelopera) ZapiszPlik(ctx context.Context,
	z shared.DeveloperFileSaveRequest) (shared.DeveloperFileSaveResponse, error) {

	okno, err := a.oknoDevelopera(z.WindowId)
	if err != nil {
		return shared.DeveloperFileSaveResponse{}, err
	}
	if err := sprawdzZmianeSystemu(okno.TrybUprawnien, "zapis pliku"); err != nil {
		return shared.DeveloperFileSaveResponse{}, err
	}
	sciezka, err := sciezkaWObszarze(a.korzenieOkna(okno), z.Path, czyIstnieje)
	if err != nil {
		return shared.DeveloperFileSaveResponse{}, err
	}
	if opis, err := os.Stat(sciezka); err == nil && opis.IsDir() {
		return shared.DeveloperFileSaveResponse{}, bladZadaniaDevelopera(
			sciezka + " jest katalogiem — zapis pliku pod tą ścieżką zniszczyłby jego zawartość")
	}
	// Katalog nadrzędny musi istnieć, inaczej zapis zamieniłby literówkę w ścieżce na nowy katalog.
	katalog := filepath.Dir(sciezka)
	if opis, err := os.Stat(katalog); err != nil || !opis.IsDir() {
		return shared.DeveloperFileSaveResponse{}, bladZasobuDevelopera(
			"katalog " + katalog + " nie istnieje — zapis nie zakłada katalogów")
	}

	wersja, err := a.zalozWersje(ctx, okno.Id, sciezka, z.CreateVersion)
	if err != nil {
		return shared.DeveloperFileSaveResponse{}, err
	}
	if err := os.WriteFile(sciezka, []byte(z.Content), 0o644); err != nil {
		return shared.DeveloperFileSaveResponse{}, bladWykonaniaDevelopera(
			"nie można zapisać pliku " + sciezka + ": " + err.Error())
	}

	opis, err := os.Stat(sciezka)
	if err != nil {
		return shared.DeveloperFileSaveResponse{}, bladWykonaniaDevelopera(
			"plik " + sciezka + " zapisany, lecz nie można odczytać jego opisu: " + err.Error())
	}
	plik := opisPliku(sciezka, opis)
	tresc := z.Content
	plik.Content = &tresc
	if wersja != "" {
		plik.VersionId = &wersja
	} else {
		plik.VersionId = a.ostatniaWersja(ctx, okno.Id, sciezka)
	}
	return shared.DeveloperFileSaveResponse{File: plik}, nil
}

// zalozWersje odkłada migawkę treści sprzed zapisu i zwraca jej identyfikator.
// Żądanie wersji przy braku dziennika kończy się odmową, nie cichym pominięciem.
func (a *adapterDevelopera) zalozWersje(ctx context.Context, oknoKod, sciezka string,
	zadana *bool) (string, error) {

	if zadana == nil || !*zadana {
		return "", nil
	}
	if a.repozytorium == nil {
		return "", bladWykonaniaDevelopera(
			"rdzeń nie ma dziennika wersji, więc nie założy punktu powrotu do treści sprzed zapisu")
	}
	poprzednia, err := os.ReadFile(sciezka)
	if errors.Is(err, os.ErrNotExist) {
		// Plik zakładany od zera nie ma treści sprzed zapisu — wersji nie ma z czego zrobić.
		return "", nil
	}
	if err != nil {
		return "", bladWykonaniaDevelopera(
			"nie można odczytać treści sprzed zapisu pliku " + sciezka + ": " + err.Error())
	}
	kod := nowyIdentyfikator(przedrostekWersjiPliku)
	if err := a.repozytorium.ZapiszWersje(ctx, dane.WersjaPliku{
		Kod:     kod,
		OknoKod: oknoKod,
		Sciezka: sciezka,
		Tresc:   string(poprzednia),
		Rozmiar: int64(len(poprzednia)),
	}); err != nil {
		return "", bladWykonaniaDevelopera("nie można założyć wersji pliku " + sciezka + ": " + err.Error())
	}
	return kod, nil
}

// ostatniaWersja zwraca identyfikator najnowszej migawki pliku. Brak dziennika
// i brak migawki dają ten sam wynik — pusty; plik bez wersji jest stanem
// zwykłym, a nie usterką odczytu.
func (a *adapterDevelopera) ostatniaWersja(ctx context.Context, oknoKod, sciezka string) *string {
	if a.repozytorium == nil {
		return nil
	}
	wersja, err := a.repozytorium.OstatniaWersja(ctx, oknoKod, sciezka)
	if err != nil || wersja.Kod == "" {
		return nil
	}
	return &wersja.Kod
}

// czyIstnieje mówi, czy pod ścieżką coś leży. Rozstrzyga wybór katalogu
// roboczego dla wskazania względnego.
func czyIstnieje(sciezka string) bool {
	_, err := os.Stat(sciezka)
	return err == nil
}

// czyBinarna rozpoznaje treść nienadającą się do edytora tekstu. Bajt zerowy
// w początkowej porcji jest tym samym rozpoznaniem, którego używa git —
// wystarczająco pewnym, a nieporównanie tańszym niż analiza kodowania.
func czyBinarna(bajty []byte) bool {
	if len(bajty) > 8000 {
		bajty = bajty[:8000]
	}
	return strings.IndexByte(string(bajty), 0) >= 0
}
