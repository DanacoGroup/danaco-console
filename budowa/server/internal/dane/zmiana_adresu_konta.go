// Odpowiedzialność pliku: trwałość zamówionej zmiany adresu uwierzytelniającego.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// ZmianaAdresuKonta to zamówiona zmiana adresu czekająca na kod z nowego adresu.
// Konto stoi przy adresie poprzednim, dopóki zmiana nie zostanie domknięta.
type ZmianaAdresuKonta struct {
	KontoId         int64
	AdresNowy       string
	AdresPoprzedni  string
	SkrotKodu       string
	SkrotWycofania  string
	WygasaKod       int64
	WygasaWycofanie int64
	Proby           int
	Zamkniete       bool
	Utworzono       int64
	ZrodloIP        string
	Urzadzenie      string
}

// RepozytoriumZmianyAdresu prowadzi zamówione zmiany adresu konta.
type RepozytoriumZmianyAdresu interface {
	// ZamowZmianeAdresu zapisuje zamówienie; poprzednie zamówienia konta zamyka.
	ZamowZmianeAdresu(ctx context.Context, z ZmianaAdresuKonta) error
	// ZmianaCzynnaKonta oddaje niezamkniętą zmianę konta żądania.
	ZmianaCzynnaKonta(ctx context.Context) (ZmianaAdresuKonta, error)
	// ZmianaPoSkrocieWycofania odnajduje zmianę samą drogą z listu, bez sesji.
	ZmianaPoSkrocieWycofania(ctx context.Context, skrot string) (ZmianaAdresuKonta, error)
	// ZamknijZmiane zamyka zamówienie: domknięte kodem albo wycofane drogą.
	ZamknijZmiane(ctx context.Context, kontoId, utworzono int64) error
	// OdnotujPomylke dolicza próbę nietrafioną do zmiany czynnej konta.
	OdnotujPomylke(ctx context.Context, kontoId, utworzono int64) (int, error)
}

const (
	kolumnyZmianyAdresu = `konto_id, adres_nowy, adres_poprzedni, skrot_kodu, skrot_wycofania,
	                       wygasa_kod, wygasa_wycofanie, proby, zamkniete, utworzono,
	                       COALESCE(zrodlo_ip, ), COALESCE(urzadzenie, )`

	wstawZmianeAdresu = `INSERT INTO zmiana_adresu_konta
	                     (konto_id, adres_nowy, adres_poprzedni, skrot_kodu, skrot_wycofania,
	                      wygasa_kod, wygasa_wycofanie, proby, zamkniete, utworzono,
	                      zrodlo_ip, urzadzenie)
	                     VALUES (?, ?, ?, ?, ?, ?, ?, 0, 0, ?, NULLIF(?, ), NULLIF(?, ))`

	zamknijZmianyKonta = `UPDATE zmiana_adresu_konta SET zamkniete = 1
	                      WHERE konto_id = ? AND zamkniete = 0`

	zmianaCzynnaKonta = `SELECT ` + kolumnyZmianyAdresu + `
	                     FROM zmiana_adresu_konta
	                     WHERE konto_id = ? AND zamkniete = 0
	                     ORDER BY utworzono DESC LIMIT 1`

	zmianaPoSkrocieWycofania = `SELECT ` + kolumnyZmianyAdresu + `
	                            FROM zmiana_adresu_konta WHERE skrot_wycofania = ?`

	zamknijZmianeAdresu = `UPDATE zmiana_adresu_konta SET zamkniete = 1
	                       WHERE konto_id = ? AND utworzono = ?`

	dolicznikPomylkiZmiany = `UPDATE zmiana_adresu_konta SET proby = proby + 1
	                          WHERE konto_id = ? AND utworzono = ?
	                          RETURNING proby`
)

type repozytoriumZmianyAdresu struct {
	zapytania *zapytania
}

var _ RepozytoriumZmianyAdresu = (*repozytoriumZmianyAdresu)(nil)

func noweRepozytoriumZmianyAdresu(z *zapytania) *repozytoriumZmianyAdresu {
	return &repozytoriumZmianyAdresu{zapytania: z}
}

// ZamowZmianeAdresu zamyka zamówienia poprzednie tego konta i zapisuje nowe.
// Bez zamknięcia poprzednich w obiegu zostaje tyle kodów i tyle dróg wycofania,
// ile razy Operator poprosił o zmianę.
func (r *repozytoriumZmianyAdresu) ZamowZmianeAdresu(ctx context.Context, z ZmianaAdresuKonta) error {
	zamkniecie, err := r.zapytania.przygotuj(ctx, zamknijZmianyKonta)
	if err != nil {
		return err
	}
	if _, err := zamkniecie.ExecContext(ctx, z.KontoId); err != nil {
		return fmt.Errorf("dane: nie można zamknąć poprzednich zamówień zmiany adresu: %w", err)
	}
	wstaw, err := r.zapytania.przygotuj(ctx, wstawZmianeAdresu)
	if err != nil {
		return err
	}
	if _, err := wstaw.ExecContext(ctx, z.KontoId, z.AdresNowy, z.AdresPoprzedni,
		z.SkrotKodu, z.SkrotWycofania, z.WygasaKod, z.WygasaWycofanie, z.Utworzono,
		z.ZrodloIP, z.Urzadzenie); err != nil {
		return fmt.Errorf("dane: nie można zamówić zmiany adresu: %w", err)
	}
	return nil
}

func (r *repozytoriumZmianyAdresu) ZmianaCzynnaKonta(ctx context.Context) (ZmianaAdresuKonta, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, zmianaCzynnaKonta)
	if err != nil {
		return ZmianaAdresuKonta{}, err
	}
	return odczytajZmiane(polecenie.QueryRowContext(ctx, KontoOperatora(ctx)))
}

func (r *repozytoriumZmianyAdresu) ZmianaPoSkrocieWycofania(ctx context.Context,
	skrot string) (ZmianaAdresuKonta, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zmianaPoSkrocieWycofania)
	if err != nil {
		return ZmianaAdresuKonta{}, err
	}
	return odczytajZmiane(polecenie.QueryRowContext(ctx, skrot))
}

func (r *repozytoriumZmianyAdresu) ZamknijZmiane(ctx context.Context, kontoId, utworzono int64) error {
	polecenie, err := r.zapytania.przygotuj(ctx, zamknijZmianeAdresu)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, kontoId, utworzono); err != nil {
		return fmt.Errorf("dane: nie można zamknąć zamówienia zmiany adresu: %w", err)
	}
	return nil
}

func (r *repozytoriumZmianyAdresu) OdnotujPomylke(ctx context.Context,
	kontoId, utworzono int64) (int, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, dolicznikPomylkiZmiany)
	if err != nil {
		return 0, err
	}
	var proby int
	if err := polecenie.QueryRowContext(ctx, kontoId, utworzono).Scan(&proby); err != nil {
		return 0, fmt.Errorf("dane: nie można doliczyć pomyłki zmiany adresu: %w", err)
	}
	return proby, nil
}

func odczytajZmiane(wiersz *sql.Row) (ZmianaAdresuKonta, error) {
	var z ZmianaAdresuKonta
	var zamkniete int64
	err := wiersz.Scan(&z.KontoId, &z.AdresNowy, &z.AdresPoprzedni, &z.SkrotKodu,
		&z.SkrotWycofania, &z.WygasaKod, &z.WygasaWycofanie, &z.Proby, &zamkniete,
		&z.Utworzono, &z.ZrodloIP, &z.Urzadzenie)
	if errors.Is(err, sql.ErrNoRows) {
		return ZmianaAdresuKonta{}, fmt.Errorf("dane: zamówienia zmiany adresu nie ma: %w", ErrBrakWiersza)
	}
	if err != nil {
		return ZmianaAdresuKonta{}, fmt.Errorf("dane: nie można odczytać zamówienia zmiany adresu: %w", err)
	}
	z.Zamkniete = zamkniete != 0
	return z, nil
}
