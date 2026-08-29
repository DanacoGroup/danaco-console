// Odpowiedzialność pliku: jedyna droga modułu Studio do programu zewnętrznego
// oraz do miejsca na dysku, w którym taki program wolno puścić; brak
// narzędzia to odmowa nazwana, nie panika.
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

// wolajNarzedzie uruchamia jeden program arsenału i oddaje jego wyjście,
// dobierając uruchamiacz, zasady izolacji i obszar roboczy okna.
func (a *adapterStudia) wolajNarzedzie(ctx context.Context, n zewnetrzne.Narzedzie,
	argumenty []string) ([]byte, error) {

	if a.uruchamiacz == nil {
		return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"moduł Studio: serwer nie ma uruchamiacza procesów, więc narzędzie "+n.Nazwa+
				" nie ma czym wystartować; naprawa: podpiąć warstwę kanału przy składaniu serwera"))
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

// bladArsenaluStudia przekłada odmowy pakietu zewnetrzne na kody kontraktu;
// brak binarium jest zapleczem niedostępnym i kodem PONAWIALNYM, reszta to
// usterka wewnętrzna z diagnostyką.
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

// katalogWsadu zakłada katalog na rozpakowany wsad POD obszarem roboczym
// okna, a nie w katalogu tymczasowym systemu, żeby podlegał tym samym
// zasadom izolacji.
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

// sciezkaZasobu wskazuje plik zasobu magazynu rdzenia; droga jest jedna dla
// całego modułu, a rdzeń złożony bez repozytorium zasobów odmawia zdaniem
// nazywającym brak.
func (a *adapterStudia) sciezkaZasobu(ctx context.Context, kodZasobu string) (string, error) {
	kod := strings.TrimSpace(kodZasobu)
	if a.zasoby == nil {
		return "", protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"moduł Studio: wskazanie materiału zasobem magazynu ("+kod+") nie ma drogi — "+
				"serwer złożony bez repozytorium zasobów. Wskaż materiał ścieżką widzianą "+
				"przez serwer; ta droga działa w całości."))
	}
	wiersz, err := a.zasoby.Zasob(ctx, kod)
	if err != nil {
		return "", protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Studio: zasobu "+kod+" nie ma w magazynie serwera"))
	}
	if wiersz.URI == nil || strings.TrimSpace(*wiersz.URI) == "" {
		return "", protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Studio: zasób "+kod+" nie ma odwołania do bajtów — nie ma czego odczytać"))
	}
	return *wiersz.URI, nil
}
