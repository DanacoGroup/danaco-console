// Odpowiedzialność pliku: magazyn artefaktów debaty — miejsce, w którym leżą
// bajty transkryptu, grafu i nagrania, oraz jedna droga ich wydania. Magazyn
// jest własny, a nie wspólny z biblioteką, bo artefakt debaty nie jest
// plikiem Operatora.
package core

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfiguracja"
)

const (
	// podkatalogDebaty oddziela magazyn modułu od reszty katalogu danych,
	// swoim własnym podkatalogiem, osobnym dla każdego modułu platformy.
	podkatalogDebaty = "roundtable"
	// podkatalogArtefaktowDebaty mieści same bajty wydanych artefaktów,
	// wewnątrz podkatalogu modułu Roundtable.
	podkatalogArtefaktowDebaty = "artefakty"
	// korzenArtefaktowDebaty jest przedrostkiem odwołania wychodzącego z rdzenia,
	// ścieżką względną liczoną od katalogu danych.
	korzenArtefaktowDebaty = podkatalogDebaty + "/" + podkatalogArtefaktowDebaty

	// Rodzaje artefaktu — wartości kolumny `rodzaj` z migracji 199, transkrypt,
	// graf debaty i odsłuch nagrania mowy uczestników.
	rodzajArtefaktuTranskryptu = "transkrypt"
	rodzajArtefaktuGrafu       = "graf"
	rodzajArtefaktuOdsluchu    = "odsluch"

	przedrostekArtefaktu = "artdeb-"

	// prawaKataloguArtefaktow: artefakt debaty należy do Operatora i do nikogo
	// więcej — katalog nie jest czytelny dla innych kont na maszynie.
	prawaKataloguArtefaktow = 0o700
)

// ZKatalogiemDanych podpina katalog, pod którym leży magazyn artefaktów.
// Katalog pusty spada na domyślny, żeby adapter nigdy nie stanął bez magazynu
// i nie odmawiał z powodu własnego montażu zamiast z powodu żądania.
func (a *adapterDebaty) ZKatalogiemDanych(katalogDanych string) *adapterDebaty {
	if strings.TrimSpace(katalogDanych) == "" {
		katalogDanych = konfiguracja.KatalogDanychDomyslny()
	}
	a.katalogArtefaktow = filepath.Join(katalogDanych, podkatalogDebaty,
		podkatalogArtefaktowDebaty)
	return a
}

// wydajArtefaktDebaty kładzie bajty w magazynie i odnotowuje artefakt w bazie.
// Nazwą pliku jest suma kontrolna jego treści, więc dwa wydania tego samego
// transkryptu zajmują jedno miejsce.
func (a *adapterDebaty) wydajArtefaktDebaty(ctx context.Context, okno, rodzaj, format string,
	bajty []byte, dlugoscMs int) (dane.ArtefaktDebaty, error) {

	if a.katalogArtefaktow == "" {
		return dane.ArtefaktDebaty{}, odmowaBrakuMagazynuDebaty()
	}
	if len(bajty) == 0 {
		return dane.ArtefaktDebaty{}, odmowaPustegoArtefaktu(rodzaj)
	}

	suma := sha256.Sum256(bajty)
	nazwa := hex.EncodeToString(suma[:])
	katalogBloku := filepath.Join(a.katalogArtefaktow, nazwa[:2])
	if err := os.MkdirAll(katalogBloku, prawaKataloguArtefaktow); err != nil {
		return dane.ArtefaktDebaty{}, bladDebaty(err)
	}
	docelowa := filepath.Join(katalogBloku, nazwa+"."+format)

	if stan, err := os.Stat(docelowa); err != nil || stan.IsDir() || stan.Size() == 0 {
		// Plik tymczasowy leży w katalogu docelowym: przemianowanie ma być niepodzielne.
		tymczasowy, err := os.CreateTemp(katalogBloku, "artefakt-*.czesciowy")
		if err != nil {
			return dane.ArtefaktDebaty{}, bladDebaty(err)
		}
		nazwaTymczasowa := tymczasowy.Name()
		if _, err := tymczasowy.Write(bajty); err != nil {
			tymczasowy.Close()
			_ = os.Remove(nazwaTymczasowa)
			return dane.ArtefaktDebaty{}, bladDebaty(err)
		}
		if err := tymczasowy.Sync(); err != nil {
			tymczasowy.Close()
			_ = os.Remove(nazwaTymczasowa)
			return dane.ArtefaktDebaty{}, bladDebaty(err)
		}
		if err := tymczasowy.Close(); err != nil {
			_ = os.Remove(nazwaTymczasowa)
			return dane.ArtefaktDebaty{}, bladDebaty(err)
		}
		if err := os.Rename(nazwaTymczasowa, docelowa); err != nil {
			_ = os.Remove(nazwaTymczasowa)
			return dane.ArtefaktDebaty{}, bladDebaty(err)
		}
	}

	artefakt := dane.ArtefaktDebaty{
		Kod: nowyIdentyfikator(przedrostekArtefaktu), Okno: okno, Rodzaj: rodzaj,
		Format: format, Odwolanie: odwolanieMagazynu(docelowa, korzenArtefaktowDebaty),
		Rozmiar: int64(len(bajty)), SumaKontrolna: nazwa, DlugoscMs: dlugoscMs,
	}
	if err := a.repozytorium.ZapiszArtefaktDebaty(ctx, artefakt); err != nil {
		return dane.ArtefaktDebaty{}, bladDebaty(err)
	}
	zapisany, err := a.repozytorium.ArtefaktDebatyPoKodzie(ctx, artefakt.Kod)
	if err != nil {
		return dane.ArtefaktDebaty{}, bladDebaty(err)
	}
	return zapisany, nil
}

// odwolanieArtefaktu oddaje odwołanie w postaci, w jakiej niesie je kontrakt.
// Odwołanie puste nie wychodzi wcale: pole `uri` jest niewymagane, a napis
// pusty klient odczytałby jako ścieżkę, pod którą nic nie leży.
func odwolanieArtefaktu(artefakt dane.ArtefaktDebaty) *string {
	if strings.TrimSpace(artefakt.Odwolanie) == "" {
		return nil
	}
	odwolanie := artefakt.Odwolanie
	return &odwolanie
}
