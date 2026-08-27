// Plik jest mostem między pozycją kolejki a realnym wykonaniem: wykonawca uruchamia pracę pozycji,
// która weszła w stan wykonywana, i zwraca jej wynik, sięgając po ten sam rejestr kanałów co tura okna.
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

// wykonawcaKroku uruchamia realną pracę pozycji kolejki. Zwrócony błąd znaczy niepowodzenie
// wykonania, brak błędu znaczy pracę zakończoną, którą silnik przesuwa do weryfikacji.
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

// kanalPozycji rozwiązuje kanał modelu i zasięgi (sesja, okno) dla pozycji, zwracając fałsz, gdy
// pozycji nie da się przypisać kanału, bo mapowanie okna wykonawcy na kanał żyje w pakiecie sesji.
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

// Wykonaj uruchamia turę kanału dla pozycji i zwraca zebraną treść odpowiedzi oraz wynik. Treść
// wraca także przy błędzie, bo to, co model zdążył oddać, jest częścią prawdy o nieudanej turze.
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
