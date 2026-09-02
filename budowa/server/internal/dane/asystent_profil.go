// Plik czyta profil asystenta: profil wskazany kodem oraz profil domyślny.
// Zlecenia, dziennik i przyjęcie polecenia leżą w innych plikach tego samego
// repozytorium.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// ProfilAsystenta to wiersz profilu niosący warstwę promptu zlecenia oraz
// nastawy tury: kanał, głos odczytu, zasięg urządzenia i zasięg pracy.
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
	                          WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	// Domyślny jest co najwyżej jeden w koncie — indeks częściowy migracji 485
	// obejmuje wskazanie konta. Bez warunku konta odczyt oddawałby profil konta
	// najstarszego; LIMIT 1 zostaje na wypadek zmiany schematu.
	pobierzProfilDomyslnyAsystenta = `SELECT ` + kolumnyProfiluAsystenta + ` FROM profil_asystenta
	                                  WHERE domyslny = 1 AND ` + WarunekKonta + `
	                                  ORDER BY id LIMIT 1`
)

// Profil zwraca profil o wskazanym kodzie. Kod nieznany wraca jako
// ErrBrakWiersza — warstwa wyższa odróżnia „takiego profilu nie ma” od „odczyt
// się nie powiódł” i tylko pierwsze z nich wolno jej przemilczeć.
func (r *repozytoriumAsystenta) Profil(ctx context.Context, kod string) (ProfilAsystenta, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzProfilAsystenta)
	if err != nil {
		return ProfilAsystenta{}, err
	}
	profil, err := odczytajProfilAsystenta(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
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
	profil, err := odczytajProfilAsystenta(polecenie.QueryRowContext(ctx, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return ProfilAsystenta{}, ErrBrakWiersza
	}
	if err != nil {
		return ProfilAsystenta{}, fmt.Errorf("dane: nieczytelny wiersz domyślnego profilu asystenta: %w", err)
	}
	return profil, nil
}

// odczytajProfilAsystenta składa strukturę profilu asystenta z jednego
// zwróconego wiersza wyniku bazy.
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
