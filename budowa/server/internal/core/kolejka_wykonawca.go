// Odpowiedzialność pliku: most między pozycją kolejki a realnym wykonaniem.
// Silnik kolejek (`kolejka_silnik.go`) prowadzi pozycje przez stany; wykonawca
// robi to, czego stan sam nie robi — uruchamia pracę pozycji, która weszła
// w stan `wykonywana`, i zwraca jej wynik.
//
// Droga wysyłki modelu jest jedna. Wykonawca nie buduje drugiego silnika ani
// drugiej drogi do modelu: sięga po ten sam rejestr kanałów i tę samą metodę
// `Wyslij`, którą jedzie tura okna (`adapter_rozmowa.go`) i głos debaty
// (`adapter_modul_roundtable_glos.go`). Strumień odpowiedzi idzie wspólnym
// nadajnikiem, więc Process Monitor widzi pracę pozycji tak samo jak turę okna.
//
// Czego pozycja nie niesie: treść zlecenia ma (`tresc_zlecenia`), lecz kanału
// modelu nie — ani schemat `pozycja_kolejki`, ani kontrakt nie mają pola
// wskazującego kanał, którym pozycję wykonać. Kanał
// dostarcza więc rozwiązywacz wpięty przy montażu: zna okno wykonawcy pozycji
// albo koordynatora kolejki i z niego bierze kanał. Pozycja bez treści albo bez
// kanału to nie cichy sukces — to realny błąd wykonania, który silnik
// zamienia na stan `bledna`.
package core

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/models"
	"danacoconsole/shared"
)

// wykonawcaKroku uruchamia realną pracę pozycji kolejki. Zwrócony błąd znaczy
// niepowodzenie wykonania — silnik przekłada je na stan `bledna`; brak błędu
// znaczy pracę zakończoną, którą silnik przesuwa do `do_weryfikacji`.
//
// Treść wyniku wraca do wołającego, bo bez niej `Subagent.result` zostawałby
// pusty, choć praca się odbyła i jej treść przepłynęła strumieniem. Silnik
// oddaje treść ujściu wyniku (`kolejka_silnik.go`), a schemat pozycji zostaje
// nietknięty: `pozycja_kolejki` nie ma kolumny wyniku i ten podpis jej nie
// dorabia — trwałość wyniku należy do wiersza, który go pokazuje (podagent).
type wykonawcaKroku interface {
	Wykonaj(ctx context.Context, pozycja dane.Pozycja) (string, error)
}

// Odmowy wykonania nazwane raz. Obie znaczą „nie ma czym wykonać kroku", więc
// silnik zamyka pozycję stanem `bledna`, a nie udanym `do_weryfikacji`.
var (
	errPozycjaBezTresci = errors.New(
		"pozycja kolejki bez treści zlecenia — nie ma czego wykonać")

	errPozycjaBezKanalu = errors.New(
		"pozycja kolejki bez kanału modelu — nie ma czym wykonać kroku")

	errWykonawcaNiegotowy = errors.New(
		"wykonawca kroku bez rejestru kanałów albo rozwiązywacza kanału")
)

// kanalPozycji rozwiązuje kanał modelu i zasięgi (sesja, okno) dla pozycji.
// Zwraca `false`, gdy pozycji nie da się przypisać kanału — wykonanie jest
// wtedy niemożliwe, a pozycja idzie w stan błędu. Rozwiązywacz wpina montaż,
// bo mapowanie okna wykonawcy na kanał żyje w pakiecie sesji, nie tutaj.
type kanalPozycji func(ctx context.Context, pozycja dane.Pozycja) (kanal string, zasiegi models.Zasiegi, ok bool)

// wykonawcaModelu wykonuje pozycję turą kanału modelu. Treść zlecenia pozycji
// staje się treścią zapytania; kanał podaje rozwiązywacz; strumień odpowiedzi
// idzie wspólnym nadajnikiem zdarzeń.
type wykonawcaModelu struct {
	kanaly   *models.Rejestr
	nadajnik Nadajnik
	kanal    kanalPozycji
}

// nowyWykonawcaModelu składa wykonawcę z rejestru kanałów, nadajnika strumienia
// i rozwiązywacza kanału pozycji.
func nowyWykonawcaModelu(kanaly *models.Rejestr, nadajnik Nadajnik, kanal kanalPozycji) wykonawcaModelu {
	return wykonawcaModelu{kanaly: kanaly, nadajnik: nadajnik, kanal: kanal}
}

// Wykonaj uruchamia turę kanału dla pozycji i zwraca zebraną treść odpowiedzi
// oraz wynik. Powodzenie tury znaczy pracę wykonaną; błąd kanału albo brak
// drogi wykonania znaczy niepowodzenie, które silnik pokaże jako stan `bledna`.
// Treść wraca także przy błędzie — to, co model zdążył oddać, jest częścią
// prawdy o nieudanej turze, nie odpadem.
func (w wykonawcaModelu) Wykonaj(ctx context.Context, pozycja dane.Pozycja) (string, error) {
	tresc := strings.TrimSpace(wartoscTekstu(pozycja.TrescZlecenia))
	if tresc == "" {
		return "", errPozycjaBezTresci
	}
	if w.kanaly == nil || w.kanal == nil {
		return "", errWykonawcaNiegotowy
	}
	kanal, zasiegi, ok := w.kanal(ctx, pozycja)
	if !ok || strings.TrimSpace(kanal) == "" {
		return "", errPozycjaBezKanalu
	}

	idPozycji := strconv.FormatInt(pozycja.ID, 10)
	strumien := nowyNadawcaStrumienia(w.nadajnik, idPozycji, zasiegi.Sesja)
	var zebrana strings.Builder
	ujscie := models.UjscieFunkcji(func(c context.Context, f models.Fragment) error {
		if f.Kind == shared.ChunkKindText {
			zebrana.WriteString(models.TrescFragmentu(f))
		}
		return strumien.Fragment(c, f)
	})

	blad := w.kanaly.Wyslij(ctx, models.Zapytanie{
		Zasiegi: zasiegi, Wiadomosc: idPozycji, Tresc: tresc, Kanal: kanal,
	}, ujscie)
	strumien.Zakoncz(zasiegi.Okno, idPozycji, zebrana.String(), blad)
	return zebrana.String(), blad
}
