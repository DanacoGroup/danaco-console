package core

import (
	"context"
	"fmt"

	"danacoconsole/server/internal/injection"
	"danacoconsole/server/internal/models"
	"danacoconsole/shared"
)

// kanalGlowny wnosi kanał główny Claude Code CLI do rejestru kanałów, łącząc pakiet injection oraz rejestr modeli w warstwie składania.
type kanalGlowny struct {
	definicja models.Definicja
	kanal     *injection.Kanal
	// przejmowanie oddaje proces tury pod uchwyt sesji, by zamknięcie okna kończyło też pracę modelu.
	przejmowanie *przejmowanieProcesow
	// zdarzenia odbiera zdarzenia zaczepów ze strumienia do dziennika zdarzeń i diagnostyki.
	zdarzenia *zdarzeniaWykonawcze
	// nazwaKonta tłumaczy wskazanie konta, identyfikator kontraktu albo nazwę, na kod konta puli.
	nazwaKonta func(string) string
}

// parametrProgramu wskazuje w wierszu rejestru ścieżkę pliku wykonywalnego
// kanału. Jest wartością danych, nie stałą kodu.
const parametrProgramu = "program"

// programDomyslny obowiązuje, gdy wiersz rejestru nie wskazuje ścieżki: nazwa
// programu rozwiązywana ze ścieżki systemowej. Brak wskazania jest wartością
// domyślną, nie odmową pracy.
const programDomyslny = "claude"

// fabrykaKanaluGlownego zwraca fabrykę adaptera kanału głównego dla rejestru modeli, dzieloną przez wszystkie wiersze tego rodzaju.
func fabrykaKanaluGlownego(pula *injection.PulaKont, przejmowanie *przejmowanieProcesow,
	zdarzenia *zdarzeniaWykonawcze, nazwaKonta func(string) string) models.Fabryka {
	return func(d models.Definicja) (models.Kanal, error) {
		return &kanalGlowny{definicja: d, kanal: injection.NowyKanal(pula),
			przejmowanie: przejmowanie, zdarzenia: zdarzenia, nazwaKonta: nazwaKonta}, nil
	}
}

// Kod zwraca kod wiersza rejestru modeli, z którego adapter kanału głównego powstał podczas budowy fabryki.
func (k *kanalGlowny) Kod() string {
	return k.definicja.Kod
}

// Definicja zwraca pełny wiersz rejestru modeli, na podstawie którego adapter kanału głównego został utworzony.
func (k *kanalGlowny) Definicja() models.Definicja {
	return k.definicja
}

// Wyslij przeprowadza jedną turę kanału głównego i przekazuje jej strumień do ujścia, w kolejności nadania fragmentów.
func (k *kanalGlowny) Wyslij(ctx context.Context, z models.Zapytanie, u models.Ujscie) error {
	// Pula bierze na progu tury bieżący wykaz kont z katalogu, więc zmiana konta wchodzi bez restartu.
	k.kanal.Pula().OdswiezZeZrodla()
	zapytanie := injection.Zapytanie{
		// Proces trafia pod uchwyt sesji zaraz po starcie, żeby zamknięcie okna kończyło też pracę modelu.
		NaStartProcesu: k.przejmowanie.Haczyk(z.Okno()),
		// Zdarzenia zaczepów jadą do dziennika zdarzeń i diagnostyki.
		NaZdarzenieZaczepu: k.zdarzenia.HaczykZaczepow(z.Okno(), z.Wiadomosc),
		IdOkna:             z.Okno(),
		IdWiadomosci:       z.Wiadomosc,
		Tekst:              z.Tresc,
		Ustawienia:         ustawieniaKanaluGlownego(k.definicja, z),
		Nakladka:           nakladkaKanaluGlownego(z),
	}
	// Konto wskazane w konfiguracji sesji albo w wierszu rejestru dojeżdża do procesu zamiast rotacji.
	zapytanie.Ustawienia.Konto = k.przelozoneKonto(z)
	var przyczyna error
	for fragment := range k.kanal.Rozmowa(ctx, zapytanie) {
		if fragment.Chunk.Kind == shared.ChunkKindError {
			przyczyna = fmt.Errorf("kanał %s: %s", k.definicja.Kod, models.TrescFragmentu(fragment.Chunk))
		}
		if err := u.Fragment(ctx, fragment.Chunk); err != nil {
			return err
		}
	}
	return przyczyna
}

// przelozoneKonto rozstrzyga wskazanie konta wywołania i tłumaczy je na kod
// konta puli. Brak wskazania i brak resolvera zostawiają dawne zachowanie:
// pusty napis znaczy rotację puli.
func (k *kanalGlowny) przelozoneKonto(z models.Zapytanie) string {
	wskazanie := z.WybraneKonto(k.definicja)
	if wskazanie == "" || k.nazwaKonta == nil {
		return wskazanie
	}
	return k.nazwaKonta(wskazanie)
}

// ustawieniaKanaluGlownego składa parametry wywołania z wiersza rejestru
// i z okna komunikacji. Okno jest poziomem najwęższym, więc jego wskazanie
// wygrywa z wierszem rejestru.
func ustawieniaKanaluGlownego(d models.Definicja, z models.Zapytanie) injection.Ustawienia {
	u := injection.Ustawienia{
		Program:       d.ParametrLub(parametrProgramu, programDomyslny),
		Model:         z.WybranyModel(d),
		ModelZapasowy: z.ModelZapasowy,
		Naklad:        z.NakladRozumowania,
		TrybUprawnien: z.TrybUprawnien,
		Katalogi:      z.KatalogiRobocze,
		Wznowienie:    z.Wznowienie,
		// Pułap kosztu jedzie do przełącznika budżetu; zero znaczy brak przełącznika w wierszu argv.
		PulapKosztuUSD: z.PulapKosztuUSD,
		// Plik ustawień i środowisko pochodzą z obszarów konfiguracji sesji; puste znaczy brak przełącznika.
		PlikUstawien: z.PlikUstawien,
		Srodowisko:   z.Srodowisko,
	}
	// Katalog startowy procesu bierze się z katalogu sesji, a w braku ustalenia z pierwszego katalogu.
	u.KatalogRoboczy = z.KatalogSesji
	if u.KatalogRoboczy == "" && len(z.KatalogiRobocze) > 0 {
		u.KatalogRoboczy = z.KatalogiRobocze[0]
	}
	// Konfiguracja MCP jedzie dwiema drogami, nadaniem okna i wiązaniem sesji; proces dostaje sumę obu.
	if z.KonfiguracjaMCP != "" {
		u.KonfiguracjaMCP = append(u.KonfiguracjaMCP, z.KonfiguracjaMCP)
	}
	u.KonfiguracjaMCP = append(u.KonfiguracjaMCP, z.DodatkoweMCP...)
	return u
}

// nakladkaKanaluGlownego przenosi warstwy nakładki: konstytucja → profil/rola →
// ekspertyza zadaniowa. Treść pochodzi z konfiguracji; rdzeń nie zna
// ani jednego zdania promptu.
func nakladkaKanaluGlownego(z models.Zapytanie) injection.Nakladka {
	return injection.NowaNakladka(z.Nakladka.Tryb,
		injection.Warstwa{Nazwa: injection.WarstwaKonstytucja, Tresc: z.Nakladka.Konstytucja},
		injection.Warstwa{Nazwa: injection.WarstwaProfil, Tresc: z.Nakladka.ProfilRoli},
		injection.Warstwa{Nazwa: injection.WarstwaEkspertyza, Tresc: z.Nakladka.Ekspertyza},
	)
}
