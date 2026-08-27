// Sprawdziany odmowy wskaźnika znaczenia; mierzą, że brak biblioteki, brak
// interpretera i wagi kończą się ODMOWĄ NAZYWAJĄCĄ BRAK, a nie usterką
// wewnętrzną.
package wiedza

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"danacoconsole/server/internal/injection"
	"danacoconsole/server/internal/session"
	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

// zasiegSprawdzianu oddaje trójkę uruchomienia procesu w zasięgu platformy —
// tę samą, którą składa adapter rdzenia; zasady puste znaczą izolację
// wyłączoną.
func zasiegSprawdzianu() (session.Okno, session.Zasady, session.Obszar) {
	okno := session.Okno{Ustawienia: session.Ustawienia{
		SrodowiskoWykonania: shared.ExecutionEnvCore,
	}}
	return okno, session.Zasady{}, session.Obszar{}
}

// silnikSprawdzianu składa silnik na prawdziwym uruchamiaczu procesów,
// tym samym, który mierzy sprawdziany odmowy w tym pakiecie.
func silnikSprawdzianu(t *testing.T, u Ustawienia) *Silnik {
	t.Helper()
	return NowySilnik(injection.UruchamiaczOkien(), t.TempDir()).ZUstawieniami(u)
}

// brakZOdmowy wyłuskuje typowany brak albo przerywa sprawdzian; błąd, który
// brakiem nie jest, znaczy, że rdzeń oddał usterkę zamiast nazwać brak.
func brakZOdmowy(t *testing.T, err error) *BrakSilnika {
	t.Helper()
	if err == nil {
		t.Fatal("silnik zameldował gotowość tam, gdzie nie ma czym liczyć")
	}
	var nazwany *BrakSilnika
	if !errors.As(err, &nazwany) {
		t.Fatalf("odmowa nie jest nazwanym brakiem, tylko usterką wewnętrzną: %v", err)
	}
	return nazwany
}

// TestBrakInterpreteraNazywaBrakINaprawe pilnuje, że nieobecny interpreter
// kończy się brakiem nazwanym, a nie komunikatem systemu o nieznanym pliku.
func TestBrakInterpreteraNazywaBrakINaprawe(t *testing.T) {
	ustawienia := UstawieniaDomyslne()
	ustawienia.Program = "danaco-interpreter-ktorego-nie-ma"

	okno, zasady, obszar := zasiegSprawdzianu()
	err := silnikSprawdzianu(t, ustawienia).Gotowy(context.Background(),
		okno, zasady, obszar, LimitZapytania)

	nazwany := brakZOdmowy(t, err)
	if nazwany.Rodzaj != BrakInterpretera {
		t.Fatalf("brak rozpoznany jako %q, a nie jako brak interpretera", nazwany.Rodzaj)
	}
	sprawdzTrzyCzlonyOdmowy(t, nazwany.Error())
}

// TestWagiBezOpisuKsztaltuDajaBrakNazwany pilnuje, że katalog, w którym
// leży plik wag bez deklaracji kształtu modelu, kończy się brakiem
// nazwanym.
func TestWagiBezOpisuKsztaltuDajaBrakNazwany(t *testing.T) {
	if _, err := exec.LookPath(interpreterPreferowany); err != nil {
		t.Skip("nie ma interpretera " + interpreterPreferowany +
			", więc pomocnika osadzeń nie ma czym uruchomić")
	}

	katalogWag := t.TempDir()
	if err := os.MkdirAll(filepath.Join(katalogWag, "onnx"), 0o700); err != nil {
		t.Fatalf("nie można przygotować katalogu wag: %v", err)
	}
	// Treść pliku nie ma znaczenia: pomocnik odmawia po odczycie deklaracji,
	// przed wczytaniem wag.
	if err := os.WriteFile(filepath.Join(katalogWag, "onnx", "model.onnx"),
		[]byte("nie-sa-to-wagi"), 0o600); err != nil {
		t.Fatalf("nie można położyć pliku wag: %v", err)
	}

	ustawienia := UstawieniaDomyslne()
	ustawienia.KatalogModeli = katalogWag

	okno, zasady, obszar := zasiegSprawdzianu()
	err := silnikSprawdzianu(t, ustawienia).Gotowy(context.Background(),
		okno, zasady, obszar, LimitZapytania)

	nazwany := brakZOdmowy(t, err)
	if nazwany.Rodzaj != brakWagStojacych && nazwany.Rodzaj != brakBiblioteki {
		t.Fatalf("brak rozpoznany jako %q; katalog z plikiem wag bez deklaracji kształtu "+
			"ma dać brak wag stojących, a maszyna bez biblioteki — brak biblioteki",
			nazwany.Rodzaj)
	}
	sprawdzTrzyCzlonyOdmowy(t, nazwany.Error())
}

// TestOdmowaNazywaBrakWagePrzyczyneINaprawe przechodzi wszystkie cztery
// rodzaje braku i pilnuje, że każdy z nich składa zdanie o trzech
// członach.
func TestOdmowaNazywaBrakWagePrzyczyneINaprawe(t *testing.T) {
	const model = "sentence-transformers/paraphrase-multilingual-mpnet-base-v2"

	przypadki := []struct {
		rodzaj    string
		wCzesci   string
		wNaprawie string
	}{
		{brakBiblioteki, "fastembed", "pip install fastembed"},
		{brakModelu, "nie ma na dysku", KluczKatalogModeli},
		{brakWagStojacych, KluczKatalogModeli, "modules.json"},
		{BrakInterpretera, "interpreter", "pip install fastembed"},
	}
	for _, przypadek := range przypadki {
		t.Run(przypadek.rodzaj, func(t *testing.T) {
			zdanie := (&BrakSilnika{Rodzaj: przypadek.rodzaj, Model: model,
				WagaMb: WagaModeluMb, Powod: "powód od pomocnika"}).Error()

			if !strings.Contains(zdanie, przypadek.wCzesci) {
				t.Errorf("odmowa nie nazywa braku — brakuje %q w zdaniu: %s",
					przypadek.wCzesci, zdanie)
			}
			if !strings.Contains(zdanie, "naprawa: ") {
				t.Errorf("odmowa nie mówi, co zrobić: %s", zdanie)
			}
			if !strings.Contains(zdanie, przypadek.wNaprawie) {
				t.Errorf("droga naprawy nie wskazuje %q: %s", przypadek.wNaprawie, zdanie)
			}
			if !strings.Contains(zdanie, "powód od pomocnika") {
				t.Errorf("odmowa gubi to, co powiedział o sobie pomocnik: %s", zdanie)
			}
		})
	}
}

// TestOdpowiedzPomocnikaZamieniaSieWNazwanyBrak pilnuje, że rodzaje braku
// wypisywane przez pomocnika i rozpoznawane przez rdzeń są tymi samymi
// napisami.
func TestOdpowiedzPomocnikaZamieniaSieWNazwanyBrak(t *testing.T) {
	for _, rodzaj := range []string{brakBiblioteki, brakModelu, brakWagStojacych} {
		if !strings.Contains(skryptPomocnika, `"`+rodzaj+`"`) {
			t.Errorf("rdzeń rozpoznaje brak %q, którego pomocnik nigdy nie wypisze — "+
				"odmowa wyszłaby jako brak nierozpoznany", rodzaj)
		}
	}

	silnik := NowySilnik(nil, t.TempDir())
	odmowa := zewnetrzne.Wynik{Wyjscie: []byte(`{"ok":false,"brak":"` + brakWagStojacych +
		`","powod":"wagi leżą, ale nie ma przy nich wykazu warstw","wagaMb":2200}`)}

	_, err := silnik.odczytaj(odmowa)
	nazwany := brakZOdmowy(t, err)
	if nazwany.Rodzaj != brakWagStojacych {
		t.Fatalf("rodzaj braku zgubiony po drodze: %q", nazwany.Rodzaj)
	}
	if nazwany.WagaMb != 2200 {
		t.Fatalf("waga braku zgubiona po drodze: %d", nazwany.WagaMb)
	}
	sprawdzTrzyCzlonyOdmowy(t, nazwany.Error())
}

// TestOdpowiedzNieczytelnaJestUsterkaANieBrakiem pilnuje granicy w drugą
// stronę: brak pomocnik umie nazwać sam, a odpowiedź nieczytelna to
// usterka.
func TestOdpowiedzNieczytelnaJestUsterkaANieBrakiem(t *testing.T) {
	silnik := NowySilnik(nil, t.TempDir())
	_, err := silnik.odczytaj(zewnetrzne.Wynik{
		Wyjscie:     []byte("Traceback (most recent call last):"),
		Diagnostyka: "ModuleNotFoundError",
	})
	if err == nil {
		t.Fatal("nieczytelna odpowiedź przeszła jako wynik")
	}
	var nazwany *BrakSilnika
	if errors.As(err, &nazwany) {
		t.Fatalf("nieczytelna odpowiedź uznana za brak %q — Operator dostałby "+
			"podpowiedź instalacyjną tam, gdzie trzeba zgłosić usterkę", nazwany.Rodzaj)
	}
	if !strings.Contains(err.Error(), "naprawa: ") {
		t.Fatalf("usterka nie mówi, co zrobić: %v", err)
	}
}

// sprawdzTrzyCzlonyOdmowy pilnuje kształtu zdania odmowy: nazwa wskaźnika,
// wyjaśnienie braku i droga naprawy, w kolejności czytelnej dla Operatora.
func sprawdzTrzyCzlonyOdmowy(t *testing.T, zdanie string) {
	t.Helper()
	for _, czlon := range []string{"wskaźnik znaczenia: ", "po ZNACZENIU", "naprawa: "} {
		if !strings.Contains(zdanie, czlon) {
			t.Errorf("odmowa bez członu %q: %s", czlon, zdanie)
		}
	}
}
