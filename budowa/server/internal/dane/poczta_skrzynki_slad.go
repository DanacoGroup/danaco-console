// Odpowiedzialność pliku: ślad wysyłki — wiersz tabeli `list_wyslany`
// zakładany po każdym wywołaniu `mail.send`, udanym i nieudanym.
//
// Ślad powstaje bezwarunkowo, bo `mail.send` jest jedyną komendą rdzenia, której
// skutek wychodzi poza maszynę Operatora i której nie da się cofnąć. Zapis jest
// jawny, trwały i niezależny od tego, czy okno rozmowy jeszcze istnieje —
// odtwarza, co wyszło ze skrzynki i kiedy.
//
// Wysyłka nieudana także jest wierszem: `powodzenie = 0` i treść błędu. Reguła
// CHECK schematu pilnuje zgodności obu kolumn — wysyłka udana nie może nieść
// błędu, nieudana musi go nieść. Dziennik z samymi powodzeniami sugerowałby,
// że wszystko doszło.
//
// Treść listu idzie do śladu w całości, bajty załączników nie: leżą w magazynie
// rdzenia pod sumą sha256, a ślad wskazuje je nazwami w treści.
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

// zapiszSladWysylkiSQL wstawia wiersz do `list_wyslany` — drugiej tabeli na
// ślad wysyłki nie ma.
const zapiszSladWysylkiSQL = `INSERT INTO list_wyslany
    (skrzynka_id, adresat, temat, tresc, w_odpowiedzi_na, powodzenie, blad)
    VALUES (?, ?, ?, ?, ?, ?, ?)`

// ZapiszSladWysylki utrwala fakt nadania listu.
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
