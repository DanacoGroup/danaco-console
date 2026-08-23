// Odpowiedzialność pliku: jedyna droga modułu Studio do programu zewnętrznego
// oraz do miejsca na dysku, w którym taki program wolno puścić.
//
// ── Dlaczego to stoi osobno ─────────────────────────────────────────────────
// Rodzina cyfryzacji woła Tesseracta, wsad woła 7-Zipa. Obie potrzebują tego
// samego: uruchamiacza, zasad izolacji obowiązujących okno i obszaru, w którym
// proces ma pracować. Gdyby każda rodzina składała ten komplet u siebie, jedna
// z nich prędzej czy później pominęłaby zasady izolacji — a to jest dokładnie
// ta pomyłka, której punkt izolacji ma nie dopuścić.
//
// Warsztat dokumentu tędy NIE idzie i iść nie będzie: `studio.pdf.*` oraz
// `studio.security.*` pracują biblioteką wkompilowaną w rdzeń, bez ani jednego
// procesu potomnego. Pilnuje tego zapora `zapora_warsztatu_pdf_test.go`.
//
// ── Brak narzędzia to odmowa nazwana, nie panika ────────────────────────────
// Rdzeń złożony bez warstwy kanału ma powiedzieć, czego mu brakuje, i pracować
// dalej w pozostałych czynnościach. Odmowa niesie NAZWĘ programu i pakiet do
// dociągnięcia, bo Operator ma przeczytać, co zainstalować, a nie szukać sam.
package core

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

// Granice czasu rodzin modułu. Rozpoznanie pisma jest najdłuższe, bo pracuje
// na obrazie strona po stronie; rozpakowanie wsadu bywa równie długie przy
// archiwum wielostronicowego skanu.
const (
	granicaRozpoznaniaStudia = 3 * time.Minute
	granicaArsenalStudia     = 2 * time.Minute
)

// wolajNarzedzie uruchamia jeden program arsenału i oddaje jego wyjście.
func (a *adapterStudia) wolajNarzedzie(ctx context.Context, n zewnetrzne.Narzedzie,
	argumenty []string) ([]byte, error) {

	if a.uruchamiacz == nil {
		return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"moduł Studio: rdzeń nie ma uruchamiacza procesów, więc narzędzie "+n.Nazwa+
				" nie ma czym wystartować; naprawa: podpiąć warstwę kanału przy składaniu rdzenia"))
	}
	granica := granicaArsenalStudia
	if n.Program == narzedzieRozpoznaniaStudia.Program {
		granica = granicaRozpoznaniaStudia
	}
	okno, zasady, obszar := a.zasiegStudia()
	wynik, err := zewnetrzne.Wolaj(ctx, a.uruchamiacz, okno, zasady, obszar, n, argumenty, "", granica)
	if err != nil {
		return wynik.Wyjscie, bladArsenaluStudia(err)
	}
	return wynik.Wyjscie, nil
}

// zasiegStudia składa okno, zasady izolacji i obszar roboczy — ten sam komplet,
// który bierze adapter narzędzi dokumentu.
func (a *adapterStudia) zasiegStudia() (session.Okno, session.Zasady, session.Obszar) {
	okno := session.Okno{Ustawienia: session.Ustawienia{
		SrodowiskoWykonania: shared.ExecutionEnvCore,
	}}
	zasady := session.Zasady{}
	if a.rozstrzygacz != nil {
		zasady = ZasadyIzolacji(a.rozstrzygacz, konfig.Kontekst{})
	}
	obszar := session.Obszar{}
	if a.katalog != nil {
		obszar = ObszarOkna(a.katalog.Ustal(konfig.Kontekst{}, ""), "")
	}
	return okno, zasady, obszar
}

// bladArsenaluStudia przekłada odmowy pakietu `zewnetrzne` na kody kontraktu.
// Rozstrzygnięcie jest to samo, co w rodzinie narzędzi dokumentu: brak
// binarium jest zapleczem niedostępnym i kodem PONAWIALNYM, bo po instalacji to
// samo żądanie przejdzie; naruszenie izolacji jest odmową uprawnienia; reszta
// to usterka wewnętrzna wraz z diagnostyką programu.
func bladArsenaluStudia(err error) error {
	var brak *zewnetrzne.BrakNarzedzia
	if errors.As(err, &brak) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
			"moduł Studio: "+brak.Error()))
	}
	if errors.Is(err, session.ErrIzolacja) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodePermissionDenied,
			"moduł Studio: "+err.Error()))
	}
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError, err))
}

// katalogWsadu zakłada katalog na rozpakowany wsad.
//
// Katalog powstaje POD obszarem roboczym okna, a nie w katalogu tymczasowym
// systemu: rozpakowany skan jest materiałem Operatora i ma podlegać tym samym
// zasadom izolacji co reszta jego pracy. Nazwa bierze się z sumy nazwy
// archiwum, więc rozpakowanie tego samego archiwum dwa razy trafia w to samo
// miejsce, zamiast mnożyć katalogi.
func (a *adapterStudia) katalogWsadu(archiwum string) (string, error) {
	korzen := os.TempDir()
	if a.katalog != nil {
		if ustalenie := a.katalog.Ustal(konfig.Kontekst{}, ""); ustalenie.Sciezka != "" {
			korzen = ustalenie.Sciezka
		}
	}
	katalog := filepath.Join(korzen, "studio-wsad", nazwaBezpieczna(filepath.Base(archiwum)))
	if err := os.MkdirAll(katalog, 0o755); err != nil {
		return "", protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError, err))
	}
	return katalog, nil
}

// nazwaBezpieczna sprowadza nazwę pliku do postaci, która nie wyjdzie
// z katalogu: bez separatorów ścieżki i bez odwołań do katalogu nadrzędnego.
func nazwaBezpieczna(nazwa string) string {
	oczyszczona := make([]rune, 0, len(nazwa))
	for _, znak := range nazwa {
		switch znak {
		case '/', '\\', ':', '.':
			oczyszczona = append(oczyszczona, '-')
		default:
			oczyszczona = append(oczyszczona, znak)
		}
	}
	if len(oczyszczona) == 0 {
		return "wsad"
	}
	return string(oczyszczona)
}

// sciezkaZasobu wskazuje plik zasobu magazynu rdzenia.
//
// Droga jest jedna dla całego modułu: wiersz zasobu niesie odwołanie do bajtów,
// a bajty leżą w magazynie zasobów pod sumą kontrolną treści. Rdzeń złożony bez
// repozytorium zasobów odmawia zdaniem nazywającym brak i podaje drogę, która
// działa bez niego — wskazanie materiału ścieżką widzianą przez rdzeń. Milczące
// zejście na pustą ścieżkę dałoby odczyt pliku, którego nie ma, i wynik
// wyglądający na prawdziwy.
func (a *adapterStudia) sciezkaZasobu(ctx context.Context, kodZasobu string) (string, error) {
	kod := strings.TrimSpace(kodZasobu)
	if a.zasoby == nil {
		return "", protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"moduł Studio: wskazanie materiału zasobem magazynu ("+kod+") nie ma drogi — "+
				"rdzeń złożony bez repozytorium zasobów. Wskaż materiał ścieżką widzianą "+
				"przez rdzeń; ta droga działa w całości."))
	}
	wiersz, err := a.zasoby.Zasob(ctx, kod)
	if err != nil {
		return "", protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Studio: zasobu "+kod+" nie ma w magazynie rdzenia"))
	}
	if wiersz.URI == nil || strings.TrimSpace(*wiersz.URI) == "" {
		return "", protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Studio: zasób "+kod+" nie ma odwołania do bajtów — nie ma czego odczytać"))
	}
	return *wiersz.URI, nil
}
