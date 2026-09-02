// Sprawdziany trwałości wyników osi obrazu: drugie zapytanie czyta składnicę, zmiana modelu albo obrazu unieważnia zapisany wynik.
package wiedza

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"danacoconsole/server/internal/session"
	"danacoconsole/server/internal/store"
)

type atrapaPomocnikaObrazu struct {
	odpowiedzi []string
	wolania    int
}

func (a *atrapaPomocnikaObrazu) UruchomProces(context.Context, session.Okno,
	session.Polecenie) (session.UchwytProcesu, error) {

	if a.wolania >= len(a.odpowiedzi) {
		return nil, errors.New("atrapa pomocnika osi obrazu: wywołanie ponad przygotowaną kolejkę odpowiedzi")
	}
	odpowiedz := a.odpowiedzi[a.wolania]
	a.wolania++
	return &uchwytAtrapyObrazu{wyjscie: strings.NewReader(odpowiedz)}, nil
}

// PID stały i niedodatni: drzewo procesu sesji nie przejmuje ani nie ubija.
type uchwytAtrapyObrazu struct {
	wyjscie io.Reader
}

func (u *uchwytAtrapyObrazu) Pid() int                { return 1 }
func (u *uchwytAtrapyObrazu) Wejscie() io.WriteCloser { return zapisNigdzie{} }
func (u *uchwytAtrapyObrazu) Wyjscie() io.Reader      { return u.wyjscie }
func (u *uchwytAtrapyObrazu) Diagnostyka() io.Reader  { return strings.NewReader("") }
func (u *uchwytAtrapyObrazu) Czekaj() error           { return nil }
func (u *uchwytAtrapyObrazu) Ubij() error             { return nil }

type zapisNigdzie struct{}

func (zapisNigdzie) Write(p []byte) (int, error) { return len(p), nil }
func (zapisNigdzie) Close() error                { return nil }

func skladnicaSprawdzianuObrazu(t *testing.T) *Skladnica {
	t.Helper()
	baza, err := store.Otworz(filepath.Join(t.TempDir(), "dane.sqlite"))
	if err != nil {
		t.Fatalf("nie można otworzyć bazy sprawdzianu: %v", err)
	}
	t.Cleanup(func() { _ = baza.Zamknij() })
	return NowaSkladnica(baza.DB)
}

func silnikObrazuSprawdzianu(t *testing.T, model string,
	odpowiedzi ...string) (*SilnikObrazu, *atrapaPomocnikaObrazu) {

	t.Helper()
	atrapa := &atrapaPomocnikaObrazu{odpowiedzi: odpowiedzi}
	silnik := NowySilnikObrazu(atrapa, t.TempDir()).
		ZUstawieniami(Ustawienia{Program: "/bin/sh", ModelObrazu: model}).
		ZeSkladnica(skladnicaSprawdzianuObrazu(t))
	return silnik, atrapa
}

func plikSprawdzianu(t *testing.T, nazwa, tresc string) string {
	t.Helper()
	sciezka := filepath.Join(t.TempDir(), nazwa)
	if err := os.WriteFile(sciezka, []byte(tresc), 0o600); err != nil {
		t.Fatalf("nie można złożyć pliku sprawdzianu: %v", err)
	}
	return sciezka
}

// Przed naprawą sprawdzian pada: Dopasuj wołało pomocnika przy każdym wywołaniu.
func TestDopasujDrugieZapytanieOTenSamObrazNieWolaPomocnika(t *testing.T) {
	silnik, atrapa := silnikObrazuSprawdzianu(t, "model-sprawdzianu",
		`{"ok":true,"model":"model-sprawdzianu","oceny":[0.42],"pominiete":[]}`)
	sciezka := plikSprawdzianu(t, "obraz.png", "treść obrazu")
	okno, zasady, obszar := zasiegSprawdzianu()
	ctx := context.Background()

	oceny1, _, err := silnik.Dopasuj(ctx, okno, zasady, obszar, "czerwone koło",
		[]string{sciezka}, time.Minute)
	if err != nil {
		t.Fatalf("pierwsze zapytanie odmówiło: %v", err)
	}
	if len(oceny1) != 1 || oceny1[0] != 0.42 {
		t.Fatalf("pierwsze zapytanie oddało %v zamiast [0.42]", oceny1)
	}
	if atrapa.wolania != 1 {
		t.Fatalf("pierwsze zapytanie wywołało pomocnika %d razy zamiast raz", atrapa.wolania)
	}

	oceny2, _, err := silnik.Dopasuj(ctx, okno, zasady, obszar, "czerwone koło",
		[]string{sciezka}, time.Minute)
	if err != nil {
		t.Fatalf("drugie zapytanie odmówiło: %v", err)
	}
	if len(oceny2) != 1 || oceny2[0] != 0.42 {
		t.Fatalf("drugie zapytanie oddało %v zamiast [0.42] ze składnicy", oceny2)
	}
	if atrapa.wolania != 1 {
		t.Fatalf("drugie zapytanie o ten sam obraz wywołało pomocnika zewnętrznego — "+
			"licznik wywołań wynosi %d zamiast pozostać przy jednym", atrapa.wolania)
	}
}

// Kryterium 2: wynik innego modelu opisuje inną przestrzeń, więc zmiana wiedza_model_obrazu unieważnia zapis.
func TestZmianaModeluObrazuWymuszaPonowneLiczenie(t *testing.T) {
	silnik, atrapa := silnikObrazuSprawdzianu(t, "pierwszy-model",
		`{"ok":true,"model":"pierwszy-model","oceny":[0.3],"pominiete":[]}`,
		`{"ok":true,"model":"drugi-model","oceny":[0.9],"pominiete":[]}`)
	sciezka := plikSprawdzianu(t, "obraz.png", "treść obrazu")
	okno, zasady, obszar := zasiegSprawdzianu()
	ctx := context.Background()

	if _, _, err := silnik.Dopasuj(ctx, okno, zasady, obszar, "pytanie",
		[]string{sciezka}, time.Minute); err != nil {
		t.Fatalf("pierwsze zapytanie odmówiło: %v", err)
	}

	silnik.ZUstawieniami(Ustawienia{Program: "/bin/sh", ModelObrazu: "drugi-model"})
	oceny, _, err := silnik.Dopasuj(ctx, okno, zasady, obszar, "pytanie",
		[]string{sciezka}, time.Minute)
	if err != nil {
		t.Fatalf("zapytanie po zmianie modelu odmówiło: %v", err)
	}
	if len(oceny) != 1 || oceny[0] != 0.9 {
		t.Fatalf("zapytanie po zmianie modelu oddało %v zamiast wyniku nowego modelu [0.9] — "+
			"zmieszało wyniki dwóch przestrzeni", oceny)
	}
	if atrapa.wolania != 2 {
		t.Fatalf("zmiana modelu nie wymusiła ponownego liczenia: pomocnik wywołany %d razy "+
			"zamiast dwóch", atrapa.wolania)
	}
}

func TestZmianaObrazuWymuszaPonowneLiczenie(t *testing.T) {
	silnik, atrapa := silnikObrazuSprawdzianu(t, "model-sprawdzianu",
		`{"ok":true,"model":"model-sprawdzianu","oceny":[0.1],"pominiete":[]}`,
		`{"ok":true,"model":"model-sprawdzianu","oceny":[0.8],"pominiete":[]}`)
	sciezka := plikSprawdzianu(t, "obraz.png", "treść pierwsza")
	okno, zasady, obszar := zasiegSprawdzianu()
	ctx := context.Background()

	if _, _, err := silnik.Dopasuj(ctx, okno, zasady, obszar, "pytanie",
		[]string{sciezka}, time.Minute); err != nil {
		t.Fatalf("pierwsze zapytanie odmówiło: %v", err)
	}

	if err := os.WriteFile(sciezka, []byte("treść druga, dłuższa niż pierwsza"), 0o600); err != nil {
		t.Fatalf("nie można nadpisać pliku sprawdzianu: %v", err)
	}

	oceny, _, err := silnik.Dopasuj(ctx, okno, zasady, obszar, "pytanie",
		[]string{sciezka}, time.Minute)
	if err != nil {
		t.Fatalf("zapytanie po zmianie obrazu odmówiło: %v", err)
	}
	if len(oceny) != 1 || oceny[0] != 0.8 {
		t.Fatalf("zapytanie po zmianie obrazu oddało %v zamiast świeżego wyniku [0.8] — "+
			"użyło wyniku policzonego dla obrazu sprzed zmiany", oceny)
	}
	if atrapa.wolania != 2 {
		t.Fatalf("zmiana obrazu nie wymusiła ponownego liczenia: pomocnik wywołany %d razy "+
			"zamiast dwóch", atrapa.wolania)
	}
}
