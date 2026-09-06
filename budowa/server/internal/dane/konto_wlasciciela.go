// Odpowiedzialność pliku: konto właściciela (tabela `konto_wlasciciela`) oraz jednorazowe drogi potwierdzenia
// tożsamości (tabela `potwierdzenie_tozsamosci`) — trwałość rejestracji, weryfikacji adresu i odzyskiwania konta.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

const (
	CelWeryfikacja = "weryfikacja"
	CelOdzyskanie  = "odzyskanie"
)

// Jednoznaczności pilnują wskaźniki na loginie i adresie, nie identyfikator wiersza (migracja 406, `decyzje.md` poz. 22).
type KontoWlasciciela struct {
	Id           int64
	Login        string
	Email        string
	Potwierdzone bool
	Utworzono    int64
}

type PotwierdzenieTozsamosci struct {
	Skrot     string
	Cel       string
	Wygasa    int64
	Uzyte     bool
	Utworzono int64
	// Bez wskazania konta potwierdzenie adresu przy dwóch kontach naraz otwierałoby konto najstarsze.
	KontoId int64
}

// Byt odrębny od katalogu metod wejścia: katalog mówi, CZYM otworzyć bramkę, konto — CZYJA ona jest.
type RepozytoriumKontaWlasciciela interface {
	Konto(ctx context.Context) (KontoWlasciciela, error)
	KontoPoTozsamosci(ctx context.Context, wskazanie string) (KontoWlasciciela, error)
	KontoPoId(ctx context.Context, kontoId int64) (KontoWlasciciela, error)
	// Praca procesu bez zamawiającego idzie po kontach po kolei (decyzja 34).
	Konta(ctx context.Context) ([]KontoWlasciciela, error)
	ZalozKonto(ctx context.Context, konto KontoWlasciciela) (int64, error)
	PotwierdzKonto(ctx context.Context, kontoId int64) error
	UsunKonto(ctx context.Context, kontoId int64) error
	// UstawAdresKonta przenosi konto na adres potwierdzony kodem ze zmiany adresu.
	UstawAdresKonta(ctx context.Context, kontoId int64, adres string) error

	ZalozPotwierdzenie(ctx context.Context, p PotwierdzenieTozsamosci) error
	PotwierdzeniePoSkrocie(ctx context.Context, skrot string) (PotwierdzenieTozsamosci, error)
	ZuzyjPotwierdzenie(ctx context.Context, skrot string, teraz int64) (bool, error)
}

const (
	kolumnyKontaWlasciciela = `id, login, email, potwierdzone, utworzono`

	// Kolejność po identyfikatorze: znacznik czasu dwóch kont z tej samej milisekundy nie rozstrzyga.
	kontoWlascicielaWiersz = `SELECT ` + kolumnyKontaWlasciciela +
		` FROM konto_wlasciciela ORDER BY id LIMIT 1`

	kontoWlascicielaPoId = `SELECT ` + kolumnyKontaWlasciciela +
		` FROM konto_wlasciciela WHERE id = ?`

	kontaWlasciciela = `SELECT ` + kolumnyKontaWlasciciela + ` FROM konto_wlasciciela ORDER BY id`

	// Okno logowania przyjmuje login i adres w tym samym polu i nie rozstrzyga, które podano.
	kontoWlascicielaPoTozsamosci = `SELECT ` + kolumnyKontaWlasciciela +
		` FROM konto_wlasciciela
		  WHERE login = ? COLLATE NOCASE OR email = ? COLLATE NOCASE
		  ORDER BY id LIMIT 1`

	wstawKontoWlasciciela = `INSERT INTO konto_wlasciciela
	                         (login, email, potwierdzone, utworzono)
	                         VALUES (?, ?, ?, ?)`

	potwierdzKontoWlasciciela = `UPDATE konto_wlasciciela SET potwierdzone = 1 WHERE id = ?`

	usunKontoWlasciciela = `DELETE FROM konto_wlasciciela WHERE id = ?`

	ustawAdresKontaWlasciciela = `UPDATE konto_wlasciciela SET email = ? WHERE id = ?`

	wstawPotwierdzenieTozsamosci = `INSERT INTO potwierdzenie_tozsamosci
	                                (skrot, cel, wygasa, uzyte, utworzono, konto_id)
	                                VALUES (?, ?, ?, 0, ?, NULLIF(?, 0))`

	potwierdzenieTozsamosciPoSkrocie = `SELECT skrot, cel, wygasa, uzyte, utworzono,
	                                           COALESCE(konto_id, 0)
	                                    FROM potwierdzenie_tozsamosci WHERE skrot = ?`

	// Warunek `uzyte = 0` w WHERE: dwa równoległe żądania dają jedno zamknięcie; niepodzielność stoi w bazie.
	zuzyjPotwierdzenieTozsamosci = `UPDATE potwierdzenie_tozsamosci
	                                SET uzyte = 1
	                                WHERE skrot = ? AND uzyte = 0 AND wygasa > ?`
)

type repozytoriumKontaWlasciciela struct {
	zapytania *zapytania
}

var _ RepozytoriumKontaWlasciciela = (*repozytoriumKontaWlasciciela)(nil)

func noweRepozytoriumKontaWlasciciela(z *zapytania) *repozytoriumKontaWlasciciela {
	return &repozytoriumKontaWlasciciela{zapytania: z}
}

func odczytajKontoWlasciciela(wiersz *sql.Row) (KontoWlasciciela, error) {
	var konto KontoWlasciciela
	var potwierdzone int64
	err := wiersz.Scan(&konto.Id, &konto.Login, &konto.Email, &potwierdzone, &konto.Utworzono)
	if errors.Is(err, sql.ErrNoRows) {
		return KontoWlasciciela{}, fmt.Errorf("dane: konta nie ma: %w", ErrBrakWiersza)
	}
	if err != nil {
		return KontoWlasciciela{}, fmt.Errorf("dane: nie można odczytać konta: %w", err)
	}
	konto.Potwierdzone = potwierdzone != 0
	return konto, nil
}

// Służy pytaniu „czy platforma ma w ogóle konto”, nie wskazaniu, czyja jest bieżąca sesja.
func (r *repozytoriumKontaWlasciciela) Konto(ctx context.Context) (KontoWlasciciela, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, kontoWlascicielaWiersz)
	if err != nil {
		return KontoWlasciciela{}, err
	}
	return odczytajKontoWlasciciela(polecenie.QueryRowContext(ctx))
}

func (r *repozytoriumKontaWlasciciela) KontoPoTozsamosci(ctx context.Context,
	wskazanie string) (KontoWlasciciela, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, kontoWlascicielaPoTozsamosci)
	if err != nil {
		return KontoWlasciciela{}, err
	}
	return odczytajKontoWlasciciela(polecenie.QueryRowContext(ctx, wskazanie, wskazanie))
}

func (r *repozytoriumKontaWlasciciela) KontoPoId(ctx context.Context,
	kontoId int64) (KontoWlasciciela, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, kontoWlascicielaPoId)
	if err != nil {
		return KontoWlasciciela{}, err
	}
	return odczytajKontoWlasciciela(polecenie.QueryRowContext(ctx, kontoId))
}

func (r *repozytoriumKontaWlasciciela) Konta(ctx context.Context) ([]KontoWlasciciela, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, kontaWlasciciela)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać kont: %w", err)
	}
	defer wiersze.Close()
	konta := []KontoWlasciciela{}
	for wiersze.Next() {
		var konto KontoWlasciciela
		var potwierdzone int64
		if err := wiersze.Scan(&konto.Id, &konto.Login, &konto.Email, &potwierdzone,
			&konto.Utworzono); err != nil {
			return nil, fmt.Errorf("dane: nie można odczytać konta: %w", err)
		}
		konto.Potwierdzone = potwierdzone != 0
		konta = append(konta, konto)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt kont: %w", err)
	}
	return konta, nil
}

func (r *repozytoriumKontaWlasciciela) ZalozKonto(ctx context.Context,
	konto KontoWlasciciela) (int64, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, wstawKontoWlasciciela)
	if err != nil {
		return 0, err
	}
	wynik, err := polecenie.ExecContext(ctx, konto.Login, konto.Email,
		liczbaLogiczna(konto.Potwierdzone), konto.Utworzono)
	if czyKolizja(err) {
		return 0, fmt.Errorf("dane: login albo adres jest już zajęty: %w", ErrKolizjaWiersza)
	}
	if err != nil {
		return 0, fmt.Errorf("dane: nie można założyć konta: %w", err)
	}
	id, err := wynik.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("dane: nie można odczytać identyfikatora założonego konta: %w", err)
	}
	return id, nil
}

func (r *repozytoriumKontaWlasciciela) PotwierdzKonto(ctx context.Context, kontoId int64) error {
	polecenie, err := r.zapytania.przygotuj(ctx, potwierdzKontoWlasciciela)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, kontoId); err != nil {
		return fmt.Errorf("dane: nie można potwierdzić konta: %w", err)
	}
	return nil
}

func (r *repozytoriumKontaWlasciciela) UsunKonto(ctx context.Context, kontoId int64) error {
	polecenie, err := r.zapytania.przygotuj(ctx, usunKontoWlasciciela)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, kontoId); err != nil {
		return fmt.Errorf("dane: nie można usunąć konta: %w", err)
	}
	return nil
}

func (r *repozytoriumKontaWlasciciela) UstawAdresKonta(ctx context.Context,
	kontoId int64, adres string) error {

	polecenie, err := r.zapytania.przygotuj(ctx, ustawAdresKontaWlasciciela)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, adres, kontoId); err != nil {
		return fmt.Errorf("dane: nie można zmienić adresu konta: %w", err)
	}
	return nil
}

func (r *repozytoriumKontaWlasciciela) ZalozPotwierdzenie(ctx context.Context,
	p PotwierdzenieTozsamosci) error {

	polecenie, err := r.zapytania.przygotuj(ctx, wstawPotwierdzenieTozsamosci)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, p.Skrot, p.Cel, p.Wygasa, p.Utworzono, p.KontoId); err != nil {
		return fmt.Errorf("dane: nie można założyć drogi potwierdzenia: %w", err)
	}
	return nil
}

func (r *repozytoriumKontaWlasciciela) PotwierdzeniePoSkrocie(ctx context.Context,
	skrot string) (PotwierdzenieTozsamosci, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, potwierdzenieTozsamosciPoSkrocie)
	if err != nil {
		return PotwierdzenieTozsamosci{}, err
	}
	var p PotwierdzenieTozsamosci
	var uzyte int64
	err = polecenie.QueryRowContext(ctx, skrot).Scan(
		&p.Skrot, &p.Cel, &p.Wygasa, &uzyte, &p.Utworzono, &p.KontoId)
	if errors.Is(err, sql.ErrNoRows) {
		return PotwierdzenieTozsamosci{}, fmt.Errorf(
			"dane: drogi potwierdzenia nie ma: %w", ErrBrakWiersza)
	}
	if err != nil {
		return PotwierdzenieTozsamosci{}, fmt.Errorf(
			"dane: nie można odczytać drogi potwierdzenia: %w", err)
	}
	p.Uzyte = uzyte != 0
	return p, nil
}

func (r *repozytoriumKontaWlasciciela) ZuzyjPotwierdzenie(ctx context.Context,
	skrot string, teraz int64) (bool, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zuzyjPotwierdzenieTozsamosci)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, skrot, teraz)
	if err != nil {
		return false, fmt.Errorf("dane: nie można zamknąć drogi potwierdzenia: %w", err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nie można odczytać skutku zamknięcia drogi: %w", err)
	}
	return zmienione > 0, nil
}
