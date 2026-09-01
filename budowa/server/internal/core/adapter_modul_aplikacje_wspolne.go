// Moduł Apps — drobiazgi wspólne wszystkim obszarom modułu: sprawdzenie okna w żądaniu, odmowa nazywająca brak bytu, wskaźnik na napis, rodzaj zmiany przy zapisie oraz rozgłoszenie etapu budowy.
package core

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"path/filepath"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// oknoAplikacji przycina i sprawdza okno żądania. Nazwa komendy wchodzi do treści odmowy, bo bez niej każda odmowa brzmiałaby tym samym zdaniem.
func oknoAplikacji(okno, komenda string) (string, error) {
	przyciete := strings.TrimSpace(okno)
	if przyciete == "" {
		return "", bladWskazaniaAplikacji(komenda + " wymaga okna")
	}
	return przyciete, nil
}

// bladNieznanegoBytuApp odróżnia „bytu nie ma" od „odczyt się nie powiódł". Okno pokazuje wtedy inny komunikat i inaczej podpowiada Operatorowi — tak samo jak `bladNieznanejArchitektury` przy architekturze.
func bladNieznanegoBytuApp(nazwaBytu, kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Apps: "+nazwaBytu+" nie istnieje: "+kod))
	}
	return bladAplikacji(err)
}

// wskaznikNapisuApp oddaje wskaźnik na napis; pusty napis daje brak wartości, bo pole opcjonalne kontraktu niosące pusty łańcuch mówiłoby „jest, ale nic".
func wskaznikNapisuApp(wartosc string) *string {
	if wartosc == "" {
		return nil
	}
	kopia := wartosc
	return &kopia
}

// sumaTresciApp liczy sumę kontrolną bajtów wytworu — nazwą pliku w magazynie jest właśnie ta suma, a `AppArtifact.checksumSha256` niesie ją do kontraktu.
func sumaTresciApp(bajty []byte) string {
	suma := sha256.Sum256(bajty)
	return hex.EncodeToString(suma[:])
}

// isBrakWierszaApp nazywa „wiersza nie ma" jednym sprawdzeniem. Powtarzane `errors.Is(err, dane.ErrBrakWiersza)` w siedmiu miejscach obszaru różniłoby się tylko literówką, którą kompilator by przepuścił.
func isBrakWierszaApp(err error) bool {
	return errors.Is(err, dane.ErrBrakWiersza)
}

// Magazyn wytworów modułu leży w katalogu danych rdzenia, obok magazynu biblioteki i zasobów Designu. Osobny podkatalog, bo moduły nie dzielą stanu: skasowanie wytworów Apps nie ma prawa ruszyć treści biblioteki.
const (
	podkatalogAplikacji = "aplikacje"
	podkatalogWytworow  = "wytwory"
)

// korzenWytworowApp jest korzeniem magazynu liczonym od katalogu danych — ta sama para podkatalogów, którą składa `magazynWytworowApp`, więc postać odwołania i miejsce zapisu nie mają jak się rozjechać.
const korzenWytworowApp = podkatalogAplikacji + "/" + podkatalogWytworow

// magazynWytworowApp składa magazyn bajtów wytworów modułu nad katalogiem danych rdzenia. Katalog nie musi istnieć — powstaje przy pierwszym zapisie.
func magazynWytworowApp(katalogDanych string) *magazynTresciBiblioteki {
	if strings.TrimSpace(katalogDanych) == "" {
		return nil
	}
	return &magazynTresciBiblioteki{
		katalog: filepath.Join(katalogDanych, podkatalogAplikacji, podkatalogWytworow),
	}
}

// odwolanieWytworuApp oddaje odwołanie w postaci, którą wolno wypuścić z rdzenia: ścieżkę względną magazynu, nie ścieżkę na dysku Operatora. Powód i cena stoją w nagłówku `odwolanieMagazynu` — to jedna decyzja dla całego produktu.
func odwolanieWytworuApp(sciezka string) string {
	return odwolanieMagazynu(sciezka, korzenWytworowApp)
}

// bladBrakuMagazynuApp nazywa brak katalogu danych. Moduł ma wtedy odmówić z powodem, a nie zameldować wytwór, za którym nie ma ani jednego bajtu — to jest dokładnie ten wzorzec szkody, który w tym produkcie już wystąpił.
func bladBrakuMagazynuApp(czynnosc string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"moduł Apps: "+czynnosc+" wymaga magazynu treści, a serwer zmontowano bez katalogu danych"))
}

// zmianaZalozenia nazywa rodzaj zmiany w rozgłoszeniu. Rozstrzyga się przed zapisem, bo zapis jest UPSERT-em i po fakcie nie widać, czy wiersz powstał.
func zmianaZalozenia(nowy bool) shared.ChangeKind {
	if nowy {
		return shared.ChangeKindCreated
	}
	return shared.ChangeKindUpdated
}

// rozglosEtap oddaje etap obsługiwaczowi, który rozsyła apps.build.changed. Zdarzenie łączy etap budowy z wdrożeniem w jednym kształcie; przy zmianie samego etapu pole Deployment zostaje puste.
func (a *adapterAplikacji) rozglosEtap(ctx context.Context, zmiana shared.ChangeKind, etap shared.AppStage) {
	if a.przyrostEtapu == nil {
		return
	}
	a.przyrostEtapu(ctx, zmiana, etap)
}
