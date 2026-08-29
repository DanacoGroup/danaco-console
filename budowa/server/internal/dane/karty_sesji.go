// Odpowiedzialność pliku: dostęp do kart sesji (tabela `karta_sesji`). Karta jest
// kontenerem sesji w obrębie środowiska — sesja bez karty nie istnieje, bo więz klucza obcego jest obowiązkowy.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// KartaSesji to wiersz tabeli `karta_sesji`, pełniący rolę kontenera sesji w obrębie środowiska pracy.
type KartaSesji struct {
	ID             int64
	SrodowiskoID   int64
	Nazwa          string
	Opis           *string
	Przypieta      bool
	Kolejnosc      int
	Utworzono      string
	Zaktualizowano string
	// KontoId wskazuje konto, do którego karta należy. Zero znaczy kartę zastaną,
	// sprzed migracji 407, gdy praca była wspólna.
	KontoId int64
}

// RepozytoriumKartSesji jest kontraktem obszaru kart sesji, określającym operacje dostępne na wykazie kart.
// Każda czynność wskazuje konto: karty jednego konta nie są widoczne dla
// drugiego, więc konto jest częścią pytania, nie doprecyzowaniem odpowiedzi.
type RepozytoriumKartSesji interface {
	Zapewnij(ctx context.Context, srodowiskoID, kontoId int64, nazwa string) (int64, error)
	Lista(ctx context.Context, srodowiskoID, kontoId int64) ([]KartaSesji, error)
}

const (
	kolumnyKartySesji = `id, srodowisko_id, nazwa, opis, przypieta, kolejnosc,
	                     utworzono, zaktualizowano, COALESCE(konto_id, 0)`

	// Karta zastana, bez wskazania konta, należy do konta najstarszego — inaczej
	// instalacja sprzed rozdzielenia straciłaby całą dotychczasową pracę z oczu.
	warunekKonta = `(konto_id = ?
	                 OR (konto_id IS NULL
	                     AND ? = (SELECT id FROM konto_wlasciciela ORDER BY id LIMIT 1)))`

	kartaSesjiPoNazwie = `SELECT ` + kolumnyKartySesji + ` FROM karta_sesji
	                      WHERE srodowisko_id = ? AND nazwa = ? AND ` + warunekKonta + `
	                      ORDER BY id LIMIT 1`

	// `NULLIF` zapisuje pustą wartość zamiast zera: karta założona bez
	// rozpoznanego konta ma należeć do tej samej rodziny co karty zastane,
	// a nie do konta o numerze zero, którego nie ma.
	wstawKarteSesji = `INSERT INTO karta_sesji (srodowisko_id, nazwa, konto_id, kolejnosc)
	                   VALUES (?, ?, NULLIF(?, 0), (SELECT COALESCE(MAX(kolejnosc) + 1, 0)
	                                                FROM karta_sesji WHERE srodowisko_id = ?))`

	listaKartSesji = `SELECT ` + kolumnyKartySesji + ` FROM karta_sesji
	                  WHERE srodowisko_id = ? AND ` + warunekKonta + `
	                  ORDER BY kolejnosc, id`
)

type repozytoriumKartSesji struct {
	zapytania *zapytania
}

func noweRepozytoriumKartSesji(z *zapytania) *repozytoriumKartSesji {
	return &repozytoriumKartSesji{zapytania: z}
}

// Zapewnij zwraca identyfikator karty o wskazanej nazwie, zakładając ją, gdy
// jeszcze nie istnieje. Wywołanie powtórne nie mnoży wierszy.
func (r *repozytoriumKartSesji) Zapewnij(ctx context.Context, srodowiskoID, kontoId int64, nazwa string) (int64, error) {
	odczyt, err := r.zapytania.przygotuj(ctx, kartaSesjiPoNazwie)
	if err != nil {
		return 0, err
	}
	karta, err := odczytajKarteSesji(odczyt.QueryRowContext(ctx, srodowiskoID, nazwa, kontoId, kontoId))
	if err == nil {
		return karta.ID, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}
	zapis, err := r.zapytania.przygotuj(ctx, wstawKarteSesji)
	if err != nil {
		return 0, err
	}
	wynik, err := zapis.ExecContext(ctx, srodowiskoID, nazwa, kontoId, srodowiskoID)
	if err != nil {
		return 0, fmt.Errorf("dane: nie można zapisać karty sesji %q: %w", nazwa, err)
	}
	return wynik.LastInsertId()
}

// Lista zwraca karty sesji środowiska w kolejności ustalonej do wyświetlania w interfejsie użytkownika.
func (r *repozytoriumKartSesji) Lista(ctx context.Context, srodowiskoID, kontoId int64) ([]KartaSesji, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaKartSesji)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, srodowiskoID, kontoId, kontoId)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać kart środowiska %d: %w", srodowiskoID, err)
	}
	defer wiersze.Close()

	lista := []KartaSesji{}
	for wiersze.Next() {
		karta, err := odczytajKarteSesji(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, karta)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt kart środowiska %d: %w", srodowiskoID, err)
	}
	return lista, nil
}

// odczytajKarteSesji składa pełną strukturę karty sesji z jednego wiersza wyniku zapytania do bazy danych.
func odczytajKarteSesji(wiersz skaner) (KartaSesji, error) {
	var karta KartaSesji
	var opis sql.NullString
	var przypieta int
	err := wiersz.Scan(&karta.ID, &karta.SrodowiskoID, &karta.Nazwa, &opis, &przypieta,
		&karta.Kolejnosc, &karta.Utworzono, &karta.Zaktualizowano, &karta.KontoId)
	if err != nil {
		return KartaSesji{}, err
	}
	karta.Opis = tekstZKolumny(opis)
	karta.Przypieta = przypieta != 0
	return karta, nil
}
