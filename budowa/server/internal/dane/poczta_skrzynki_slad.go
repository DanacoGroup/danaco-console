// Plik utrwala ślad każdej wysyłki listu w tabeli list_wyslany, zapisywany
// bezwarunkowo po każdym wywołaniu wysyłki, zarówno udanym, jak i nieudanym.
package dane

import (
	"context"
	"fmt"
)

// SladWysylki to jeden zapis faktu nadania. `Blad` jest niepusty wyłącznie
// przy `Powodzenie == false` — odwrotność łamie regułę CHECK schematu i zapis
// odmówi, zamiast utrwalić ślad wewnętrznie sprzeczny.
type SladWysylki struct {
	SkrzynkaID    int64
	Adresat       string
	Temat         string
	Tresc         string
	WOdpowiedziNa *string
	Powodzenie    bool
	Blad          *string
}

// zapiszSladWysylkiSQL wstawia jeden wiersz do tabeli list_wyslany, jedynej
// tabeli przechowującej ślad wysyłki listu w tym module.
const zapiszSladWysylkiSQL = `INSERT INTO list_wyslany
    (skrzynka_id, adresat, temat, tresc, w_odpowiedzi_na, powodzenie, blad)
    VALUES (?, ?, ?, ?, ?, ?, ?)`

// ZapiszSladWysylki utrwala fakt nadania listu, zapisując adresata, treść
// i wynik wysyłki niezależnie od tego, czy okno rozmowy nadal istnieje.
func (r *repozytoriumSkrzynek) ZapiszSladWysylki(ctx context.Context, s SladWysylki) error {
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszSladWysylkiSQL)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, s.SkrzynkaID, s.Adresat, s.Temat, s.Tresc,
		tekstDoKolumny(s.WOdpowiedziNa), liczbaLogiczna(s.Powodzenie), tekstDoKolumny(s.Blad))
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać śladu wysyłki do %q: %w", s.Adresat, err)
	}
	return nil
}
