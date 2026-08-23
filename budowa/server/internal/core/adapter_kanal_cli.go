package core

import (
	"context"
	"fmt"

	"danacoconsole/server/internal/injection"
	"danacoconsole/server/internal/models"
	"danacoconsole/shared"
)

// kanalGlowny wnosi kanał główny (Claude Code CLI) do rejestru kanałów.
//
// Rejestr modeli zna wyłącznie interfejs models.Kanal i klucz adaptera będący
// wartością danych; pakiet injection zna wyłącznie swój proces, pulę kont
// i strumień. Żaden nie importuje drugiego — łączy je ten adapter,
// mieszkający w warstwie składania.
type kanalGlowny struct {
	definicja models.Definicja
	kanal     *injection.Kanal
	// przejmowanie oddaje proces tury pod uchwyt sesji, dzięki czemu zamknięcie
	// okna kończy także to, co model uruchomił. Puste nie blokuje tury.
	przejmowanie *przejmowanieProcesow
	// zdarzenia odbiera zdarzenia zaczepów ze strumienia: dziennik
	// zdarzeń i diagnostyka. Puste nie blokuje tury — znika ślad.
	zdarzenia *zdarzeniaWykonawcze
	// nazwaKonta tłumaczy wskazanie konta (identyfikator kontraktu albo
	// nazwa) na kod konta puli. Puste nie blokuje tury — wskazanie
	// jedzie wtedy dosłownie.
	nazwaKonta func(string) string
}

// parametrProgramu wskazuje w wierszu rejestru ścieżkę pliku wykonywalnego
// kanału. Jest wartością danych, nie stałą kodu.
const parametrProgramu = "program"

// programDomyslny obowiązuje, gdy wiersz rejestru nie wskazuje ścieżki: nazwa
// programu rozwiązywana ze ścieżki systemowej. Brak wskazania jest wartością
// domyślną, nie odmową pracy.
const programDomyslny = "claude"

// fabrykaKanaluGlownego zwraca fabrykę adaptera kanału głównego dla rejestru
// modeli. Pula kont żyje dłużej niż jedno wywołanie, więc jest wspólna dla
// wszystkich wierszy tego rodzaju — rotacja konta po wyczerpaniu limitu ma sens
// tylko wtedy, gdy pamięć wyczerpania jest jedna.
func fabrykaKanaluGlownego(pula *injection.PulaKont, przejmowanie *przejmowanieProcesow,
	zdarzenia *zdarzeniaWykonawcze, nazwaKonta func(string) string) models.Fabryka {
	return func(d models.Definicja) (models.Kanal, error) {
		return &kanalGlowny{definicja: d, kanal: injection.NowyKanal(pula),
			przejmowanie: przejmowanie, zdarzenia: zdarzenia, nazwaKonta: nazwaKonta}, nil
	}
}

// Kod zwraca kod wiersza rejestru, z którego kanał powstał.
func (k *kanalGlowny) Kod() string {
	return k.definicja.Kod
}

// Definicja zwraca wiersz rejestru.
func (k *kanalGlowny) Definicja() models.Definicja {
	return k.definicja
}

// Wyslij przeprowadza jedną turę kanału głównego i przekazuje jej strumień do
// ujścia. Fragmenty idą w kolejności nadania — pierwszy jest fragment
// prowenancji, który składa kanał.
//
// Błąd tury wraca wynikiem, a nie drugim fragmentem: fragment błędu nadał już
// kanał, a rejestr modeli nie ma powielać tej samej przyczyny.
func (k *kanalGlowny) Wyslij(ctx context.Context, z models.Zapytanie, u models.Ujscie) error {
	// Na progu tury pula bierze bieżący wykaz kont z katalogu, więc konto dodane
	// albo usunięte komendą account.* wchodzi do rotacji bez restartu.
	// Pusta pula i pula bez wpiętego źródła zostają nietknięte.
	k.kanal.Pula().OdswiezZeZrodla()
	zapytanie := injection.Zapytanie{
		// Proces tury trafia pod uchwyt sesji zaraz po starcie, żeby zamknięcie
		// okna kończyło także to, co model uruchomił.
		NaStartProcesu: k.przejmowanie.Haczyk(z.Okno()),
		// Zdarzenia zaczepów jadą do dziennika zdarzeń i diagnostyki.
		NaZdarzenieZaczepu: k.zdarzenia.HaczykZaczepow(z.Okno(), z.Wiadomosc),
		IdOkna:             z.Okno(),
		IdWiadomosci:       z.Wiadomosc,
		Tekst:              z.Tresc,
		Ustawienia:         ustawieniaKanaluGlownego(k.definicja, z),
		Nakladka:           nakladkaKanaluGlownego(z),
	}
	// Konto wskazane w obszarze account konfiguracji sesji albo w wierszu
	// rejestru dojeżdża do procesu; bez tego tor CLI jechałby zawsze rotacją
	// puli. Wskazanie kontraktowe (identyfikator liczbowy) tłumaczy na kod konta
	// puli resolver montażu.
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
		// Pułap kosztu jedzie do --max-budget-usd. Zero znaczy brak przełącznika,
		// więc wiersz argv okna bez nastawy nie zmienia się o bajt.
		PulapKosztuUSD: z.PulapKosztuUSD,
		// Plik ustawień i środowisko z obszarów konfiguracji sesji (tools,
		// permissions, environment, provider). Puste znaczy brak przełącznika
		// i brak zmiennych.
		PlikUstawien: z.PlikUstawien,
		Srodowisko:   z.Srodowisko,
	}
	// Katalog startowy procesu bierze się z ustawienia katalogu sesji. Brak
	// ustalenia schodzi na pierwszy katalog roboczy okna.
	u.KatalogRoboczy = z.KatalogSesji
	if u.KatalogRoboczy == "" && len(z.KatalogiRobocze) > 0 {
		u.KatalogRoboczy = z.KatalogiRobocze[0]
	}
	// Konfiguracja MCP jedzie dwiema drogami, każda osobnym --mcp-config: nadania
	// okna (KonfiguracjaMCP) i wiązania obszaru mcp konfiguracji sesji
	// (DodatkoweMCP). Nie przykrywają się — proces dostaje sumę obu.
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
