// Odpowiedzialność pliku: profil asystenta (tabela `profil_asystenta`,
// migracja 117) — odczyt profilu wskazanego kodem oraz profilu domyślnego.
// Zlecenia leżą w `asystent.go`, dziennik w `asystent_dziennik.go`, przyjęcie
// polecenia w `asystent_polecenia.go`: jedno repozytorium, cztery pliki wedle
// odpowiedzialności.
//
// Profil niesie warstwę promptu zlecenia — zdanie, które mówi modelowi, że jest
// klawiaturą Operatora, a nie autorem odpowiedzi. To ono rozstrzyga, czy model
// sięgnie po narzędzia platformy, czy odpisze tekstem. Reszta kolumn to nastawy
// tury (kanał, głos odczytu, zasięg urządzenia, zasięg pracy), które w innym
// razie bierze wiersz okna.
//
// Plik ma sam odczyt, bez zapisu. Zakładanie profilu, wykaz profili i wskazanie
// domyślnego to trzy czynności Operatora, a kontrakt nie ma dla nich ani jednej
// komendy (`assistant.profile.*` nie istnieje). Metoda zapisu bez wołającego
// byłaby drogą, której nikt nie przechodzi, więc jej tu nie ma.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// ProfilAsystenta to wiersz tabeli `profil_asystenta`.
//
// `SrodowiskoWykonania` i `TrybUprawnien` niosą wartości kontraktu wprost
// (shared.ExecutionEnv, shared.PermissionMode) — kolumna ma na nich warunek
// CHECK, a rdzeń wkłada je do zapytania kanału bez przekładu. Pakiet `dane` nie
// zależy od `shared`, więc typem jest tu napis; jedynym miejscem, w którym te
// napisy stają się typami kontraktu, jest adapter modułu.
//
// Wskaźniki przy `WarstwaPromptu`, `KanalModelu` i `GlosSyntezy` są rozmyślne:
// NULL znaczy „profil nie ma w tej sprawie zdania" i wtedy obowiązuje nastawa
// okna. Pusty napis znaczyłby „profil kasuje nastawę okna" — a profil ma zasięg
// pracy poszerzać, nie zabierać.
type ProfilAsystenta struct {
	ID                  int64
	Kod                 string
	Nazwa               string
	WarstwaPromptu      *string
	KanalModelu         *string
	GlosSyntezy         *string
	SrodowiskoWykonania string
	TrybUprawnien       string
	Domyslny            bool
	Utworzono           int64
	Zaktualizowano      int64
}

const (
	kolumnyProfiluAsystenta = `id, identyfikator_zewnetrzny, nazwa, warstwa_promptu, kanal_modelu,
	                           glos_syntezy, srodowisko_wykonania, tryb_uprawnien, domyslny,
	                           utworzono, zaktualizowano`

	pobierzProfilAsystenta = `SELECT ` + kolumnyProfiluAsystenta + ` FROM profil_asystenta
	                          WHERE identyfikator_zewnetrzny = ?`

	// Domyślny jest co najwyżej jeden — pilnuje tego indeks częściowy migracji
	// 117. LIMIT 1 stoi tu mimo to, żeby odczyt nie zależał od tego, czy indeks
	// przetrwał każdą przyszłą zmianę schematu.
	pobierzProfilDomyslnyAsystenta = `SELECT ` + kolumnyProfiluAsystenta + ` FROM profil_asystenta
	                                  WHERE domyslny = 1 ORDER BY id LIMIT 1`
)

// Profil zwraca profil o wskazanym kodzie. Kod nieznany wraca jako
// ErrBrakWiersza — warstwa wyższa odróżnia „takiego profilu nie ma” od „odczyt
// się nie powiódł” i tylko pierwsze z nich wolno jej przemilczeć.
func (r *repozytoriumAsystenta) Profil(ctx context.Context, kod string) (ProfilAsystenta, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzProfilAsystenta)
	if err != nil {
		return ProfilAsystenta{}, err
	}
	profil, err := odczytajProfilAsystenta(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return ProfilAsystenta{}, ErrBrakWiersza
	}
	if err != nil {
		return ProfilAsystenta{}, fmt.Errorf("dane: nieczytelny wiersz profilu asystenta %q: %w", kod, err)
	}
	return profil, nil
}

// ProfilDomyslny zwraca profil oznaczony jako domyślny. Brak takiego profilu to
// ErrBrakWiersza, a nie usterka: instalacja, w której Operator nie wskazał
// jeszcze żadnego profilu, ma pracować dalej — bez warstwy promptu.
func (r *repozytoriumAsystenta) ProfilDomyslny(ctx context.Context) (ProfilAsystenta, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzProfilDomyslnyAsystenta)
	if err != nil {
		return ProfilAsystenta{}, err
	}
	profil, err := odczytajProfilAsystenta(polecenie.QueryRowContext(ctx))
	if errors.Is(err, sql.ErrNoRows) {
		return ProfilAsystenta{}, ErrBrakWiersza
	}
	if err != nil {
		return ProfilAsystenta{}, fmt.Errorf("dane: nieczytelny wiersz domyślnego profilu asystenta: %w", err)
	}
	return profil, nil
}

// odczytajProfilAsystenta składa strukturę z jednego wiersza wyniku.
func odczytajProfilAsystenta(wiersz skaner) (ProfilAsystenta, error) {
	var profil ProfilAsystenta
	var warstwa, kanal, glos sql.NullString
	var domyslny int64
	err := wiersz.Scan(&profil.ID, &profil.Kod, &profil.Nazwa, &warstwa, &kanal, &glos,
		&profil.SrodowiskoWykonania, &profil.TrybUprawnien, &domyslny,
		&profil.Utworzono, &profil.Zaktualizowano)
	if err != nil {
		return ProfilAsystenta{}, err
	}
	profil.WarstwaPromptu = tekstZKolumny(warstwa)
	profil.KanalModelu = tekstZKolumny(kanal)
	profil.GlosSyntezy = tekstZKolumny(glos)
	profil.Domyslny = domyslny != 0
	return profil, nil
}
