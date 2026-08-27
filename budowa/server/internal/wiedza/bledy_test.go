// Sprawdziany odmowy wskaźnika znaczenia.
//
// Mierzą jedną rzecz: że brak biblioteki, brak interpretera i wagi, na których
// nie da się postawić silnika, kończą się ODMOWĄ NAZYWAJĄCĄ BRAK, a nie usterką
// wewnętrzną. Różnica jest cała po stronie czytelnika: usterka wewnętrzna mówi
// „coś się zepsuło" i zachęca do ponowienia, które da to samo, a odmowa mówi,
// czego nie ma, ile to waży i co zrobić, żeby było.
//
// Uruchomienie idzie drogą produkcyjną — `Silnik` → `zewnetrzne.Wolaj` →
// `injection.UruchamiaczOkien` — a nie własnym startem procesu. Sprawdzian
// omijający tę drogę mierzyłby skrypt, a nie zachowanie rdzenia, i przespałby
// każdą zmianę w bramie izolacji albo w rozpoznawaniu braku narzędzia.
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
// tę samą, którą składa adapter rdzenia. Zasady puste znaczą izolację
// wyłączoną, co jest tu właściwe: sprawdzian mierzy odmowę silnika, a nie bramę
// izolacji, która ma własne sprawdziany.
func zasiegSprawdzianu() (session.Okno, session.Zasady, session.Obszar) {
	okno := session.Okno{Ustawienia: session.Ustawienia{
		SrodowiskoWykonania: shared.ExecutionEnvCore,
	}}
	return okno, session.Zasady{}, session.Obszar{}
}

// silnikSprawdzianu składa silnik na prawdziwym uruchamiaczu procesów.
func silnikSprawdzianu(t *testing.T, u Ustawienia) *Silnik {
	t.Helper()
	return NowySilnik(injection.UruchamiaczOkien(), t.TempDir()).ZUstawieniami(u)
}

// brakZOdmowy wyłuskuje typowany brak albo przerywa sprawdzian. Błąd, który
// brakiem nie jest, znaczy dokładnie to, czego te sprawdziany pilnują: rdzeń
// oddał usterkę tam, gdzie miał nazwać brak.
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
//
// Nazwa programu jest tu celowo taka, jakiej nikt nie zainstaluje: sprawdzian
// ma mierzyć brak, a nie to, co akurat stoi na maszynie.
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

// TestWagiBezOpisuKsztaltuDajaBrakNazwany pilnuje, że katalog, w którym leży
// plik wag bez deklaracji kształtu modelu, kończy się brakiem nazwanym.
//
// Ta droga jest tą, którą wskaźnik chodzi na wdrożeniu: nastawa
// `wiedza_katalog_modeli` wskazuje katalog z gotowym modelem, a pomocnik czyta
// przy nim sposób składania tokenów, normalizację i wymiar. Katalog bez tych
// deklaracji jest brakiem, bo zgadnięcie ich dałoby wektory, które są liczbami
// i nie znaczą nic.
//
// Maszyna bez biblioteki osadzeń odpowie na to samo zlecenie brakiem
// biblioteki — wcześniejszym w kolejności i równie nazwanym. Sprawdzian
// przyjmuje oba, bo mierzy KLASĘ odpowiedzi, nie stan maszyny.
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
	// zanim dojdzie do wczytania samych wag. Plik pusty przeszedłby jednak
	// sprawdzenie obecności, nie będąc wagami, więc niesie kilka bajtów.
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

// TestOdmowaNazywaBrakWagePrzyczyneINaprawe przechodzi wszystkie cztery rodzaje
// braku i pilnuje, że każdy z nich składa zdanie o trzech członach.
//
// Sprawdzian jest tabelą, a nie czterema funkcjami, bo mierzy jedną regułę
// w czterech miejscach: rodzaj dopisany bez własnego zdania trafi w gałąź
// domyślną i wyjdzie z niej zdaniem „silnik nie odpowiedział zrozumiale", które
// nie mówi ani czego brak, ani co zrobić.
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
// wypisywane przez pomocnika i rozpoznawane przez rdzeń są tymi samymi napisami.
//
// Rozjazd o jedną literę nie psuje niczego widocznie: brak przechodzi przez
// `odczytaj` jako brak nierozpoznany i wychodzi z niego zdaniem „silnik osadzeń
// nie odpowiedział zrozumiale", czyli odmową gorszą od tej, którą pomocnik już
// napisał. Sprawdzian czyta napisy ze SKRYPTU, a nie powtarza ich za stałymi,
// bo powtórzenie mierzyłoby samo siebie.
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
// stronę. Brak pomocnik umie nazwać sam, więc odpowiedź, której nie da się
// przeczytać, znaczy, że na wyjście pisało coś innego niż on — i to jest
// usterka do zgłoszenia, nie brak do dociągnięcia.
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
// wyjaśnienie, dlaczego rdzeń nie zejdzie po cichu na wyszukiwanie po słowach,
// oraz droga naprawy.
func sprawdzTrzyCzlonyOdmowy(t *testing.T, zdanie string) {
	t.Helper()
	for _, czlon := range []string{"wskaźnik znaczenia: ", "po ZNACZENIU", "naprawa: "} {
		if !strings.Contains(zdanie, czlon) {
			t.Errorf("odmowa bez członu %q: %s", czlon, zdanie)
		}
	}
}
