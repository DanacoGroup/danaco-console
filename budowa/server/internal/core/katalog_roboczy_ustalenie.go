package core

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"danacoconsole/server/internal/konfig"
)

// DegradacjaKatalogu opisuje zejście na lokalizację zastępczą. Katalog roboczy
// niedostępny do zapisu nie przerywa startu sesji: sesja dostaje katalog
// zastępczy, a Operator dostaje wiadomość o tym, że pracuje gdzie indziej,
// niż ustawił.
type DegradacjaKatalogu struct {
	// Zadana to ścieżka wynikająca z ustawień Operatora.
	Zadana string
	// Zastepcza to ścieżka faktycznie użyta.
	Zastepcza string
	// Powod niesie przyczynę zejścia — treść błędu systemu plików.
	Powod string
	// Skuteczna mówi, czy lokalizacja zastępcza sama dała się przygotować; sesja startuje mimo to.
	Skuteczna bool
}

// UstalenieKatalogu to rozstrzygnięty katalog roboczy jednej sesji wraz
// ze wskazaniem, skąd wzięła się podstawa i czy doszło do degradacji.
type UstalenieKatalogu struct {
	// Podstawa to katalog, w którym powstają katalogi sesyjne.
	Podstawa string
	// Wzorzec to wzorzec nazwy katalogu sesji liczony względem podstawy.
	Wzorzec string
	// Sciezka to katalog tej sesji.
	Sciezka string
	// PochodzeniePodstawy mówi, czy podstawa pochodzi z zapisu Operatora, z wartości domyślnej.
	PochodzeniePodstawy konfig.Pochodzenie
	// PoziomPodstawy wskazuje poziom zasięgu, z którego wzięła się podstawa.
	PoziomPodstawy konfig.Poziom
	// Degradacja jest niepusta wyłącznie wtedy, gdy katalog żądany okazał się niezdatny do zapisu.
	Degradacja *DegradacjaKatalogu
}

// obserwatorKatalogu przyjmuje wiadomość o degradacji katalogu roboczego.
// Rdzeń nie zakłada, kto słucha: dziennik, telemetria albo nadajnik zdarzeń
// transportu. Brak obserwatora nie zmienia zachowania ustalenia.
type obserwatorKatalogu interface {
	KatalogZdegradowany(DegradacjaKatalogu)
}

// ObserwatorKataloguFunkcja podpina funkcję jako obserwatora degradacji tego katalogu roboczego sesji.
type ObserwatorKataloguFunkcja func(DegradacjaKatalogu)

// KatalogZdegradowany wykonuje funkcję obserwatora zarejestrowaną dla tego katalogu roboczego tej sesji.
func (f ObserwatorKataloguFunkcja) KatalogZdegradowany(d DegradacjaKatalogu) { f(d) }

// KatalogRoboczy ustala katalog roboczy sesji z ustawień Operatora, a przy ich
// braku — z miejsca instalacji aplikacji głównej.
type KatalogRoboczy struct {
	rozstrzygacz *konfig.Rozstrzygacz
	instalacja   string
	zastepcza    string
	przygotuj    func(string) error
	obserwator   obserwatorKatalogu
}

// opcjaKatalogu zmienia jedną nastawę katalogu roboczego przy jego składaniu w tym module rdzenia platformy.
type opcjaKatalogu func(*KatalogRoboczy)

// ZObserwatoremKatalogu podpina odbiorcę wiadomości o degradacji katalogu roboczego tej sesji rdzenia.
func ZObserwatoremKatalogu(obserwator obserwatorKatalogu) opcjaKatalogu {
	return func(k *KatalogRoboczy) { k.obserwator = obserwator }
}

// NowyKatalogRoboczy składa ustalacz katalogu roboczego. Brak rozstrzygacza nie
// jest błędem — ustalenie schodzi wtedy na wartości domyślne.
func NowyKatalogRoboczy(rozstrzygacz *konfig.Rozstrzygacz, opcje ...opcjaKatalogu) *KatalogRoboczy {
	katalog := &KatalogRoboczy{
		rozstrzygacz: rozstrzygacz,
		instalacja:   KatalogInstalacji(),
		zastepcza:    lokalizacjaZastepczaDomyslna(),
		przygotuj:    przygotujKatalogZapisywalny,
	}
	for _, opcja := range opcje {
		if opcja != nil {
			opcja(katalog)
		}
	}
	if katalog.przygotuj == nil {
		katalog.przygotuj = przygotujKatalogZapisywalny
	}
	return katalog
}

// lokalizacjaZastepczaDomyslna zwraca katalog zastępczy platformy: podkatalog
// katalogu plików tymczasowych systemu. Jest to ostatnie miejsce, w którym
// proces prawie zawsze może pisać.
func lokalizacjaZastepczaDomyslna() string {
	return filepath.Join(os.TempDir(), "danaco-console")
}

// Ustal zwraca katalog roboczy sesji obowiązujący w kontekście zasięgu.
// Metoda nie zwraca błędu z zamysłem: każda ścieżka kończy się katalogiem,
// a niepowodzenie przygotowania zapisuje się jako degradacja.
func (k *KatalogRoboczy) Ustal(kontekst konfig.Kontekst, identyfikatorSesji string) UstalenieKatalogu {
	if k == nil {
		return UstalenieKatalogu{}
	}
	podstawaWynik := k.rozstrzygnij(kontekst, KluczKatalogRoboczyPodstawa)
	wzorzecWynik := k.rozstrzygnij(kontekst, KluczKatalogRoboczyWzorzecSesji)

	ustalenie := UstalenieKatalogu{
		Podstawa:            PodstawaLubInstalacja(podstawaWynik.Wartosc, k.instalacja),
		Wzorzec:             WzorzecLubDomyslny(wzorzecWynik.Wartosc),
		PochodzeniePodstawy: podstawaWynik.Pochodzenie,
		PoziomPodstawy:      podstawaWynik.Poziom,
	}
	ustalenie.Sciezka = SciezkaSesji(ustalenie.Podstawa, ustalenie.Wzorzec, identyfikatorSesji)

	blad := k.przygotuj(ustalenie.Sciezka)
	if blad == nil {
		return ustalenie
	}
	return k.zdegraduj(ustalenie, identyfikatorSesji, blad)
}

// zdegraduj schodzi na lokalizację zastępczą i zgłasza to zdarzenie zarejestrowanemu obserwatorowi katalogu.
func (k *KatalogRoboczy) zdegraduj(ustalenie UstalenieKatalogu, identyfikatorSesji string, przyczyna error) UstalenieKatalogu {
	zadana := ustalenie.Sciezka
	podstawa := PodstawaLubInstalacja(k.zastepcza, lokalizacjaZastepczaDomyslna())
	zastepcza := SciezkaSesji(podstawa, ustalenie.Wzorzec, identyfikatorSesji)

	degradacja := DegradacjaKatalogu{
		Zadana:    zadana,
		Zastepcza: zastepcza,
		Powod:     przyczyna.Error(),
		Skuteczna: k.przygotuj(zastepcza) == nil,
	}
	ustalenie.Podstawa = podstawa
	ustalenie.Sciezka = zastepcza
	ustalenie.Degradacja = &degradacja
	if k.obserwator != nil {
		k.obserwator.KatalogZdegradowany(degradacja)
	}
	return ustalenie
}

// rozstrzygnij pyta rezolwer ośmiu poziomów zasięgu o jedno ustawienie.
// Brak rezolwera daje wynik pusty, a nie odmowę ustalenia.
func (k *KatalogRoboczy) rozstrzygnij(kontekst konfig.Kontekst, klucz string) konfig.Wynik {
	if k.rozstrzygacz == nil {
		return konfig.Wynik{Klucz: klucz, Pochodzenie: konfig.PochodzenieNieznane}
	}
	return k.rozstrzygacz.Rozstrzygnij(kontekst, klucz)
}

// przygotujKatalogZapisywalny tworzy katalog i sprawdza, czy da się w nim pisać.
// Sprawdzenie jest zapisem próbnym, bo bity uprawnień nie odpowiadają na to
// pytanie w systemach z listami kontroli dostępu.
func przygotujKatalogZapisywalny(sciezka string) error {
	if strings.TrimSpace(sciezka) == "" {
		return errors.New("pusta ścieżka katalogu roboczego")
	}
	if err := os.MkdirAll(sciezka, 0o750); err != nil {
		return err
	}
	proba, err := os.CreateTemp(sciezka, ".danaco-zapis-*")
	if err != nil {
		return err
	}
	nazwa := proba.Name()
	if err := proba.Close(); err != nil {
		os.Remove(nazwa)
		return err
	}
	return os.Remove(nazwa)
}
