// Odpowiedzialność pliku: konto właściciela (tabela `konto_wlasciciela`) oraz jednorazowe drogi potwierdzenia
// tożsamości (tabela `potwierdzenie_tozsamosci`) — trwałość rejestracji, weryfikacji adresu i odzyskiwania konta.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// Cele drogi potwierdzenia. Są dwa i nie rosną: adres potwierdza się przy
// rejestracji, a konto odzyskuje po utracie hasła.
const (
	CelWeryfikacja = "weryfikacja"
	CelOdzyskanie  = "odzyskanie"
)

// KontoWlasciciela to wiersz tabeli `konto_wlasciciela` — jedyne konto platformy, ograniczone warunkiem schematu.
type KontoWlasciciela struct {
	Login        string
	Email        string
	Potwierdzone bool
	Utworzono    int64
}

// PotwierdzenieTozsamosci to jednorazowa droga wysłana listem: skrót materiału,
// cel, czas wygaśnięcia i znacznik użycia.
type PotwierdzenieTozsamosci struct {
	Skrot     string
	Cel       string
	Wygasa    int64
	Uzyte     bool
	Utworzono int64
}

// RepozytoriumKontaWlasciciela jest kontraktem trwałości rejestracji.
//
// Byt jest odrębny od katalogu metod wejścia, bo odpowiada na inne pytanie:
// katalog metod mówi, CZYM otworzyć bramkę, konto mówi, CZYJA ona jest.
type RepozytoriumKontaWlasciciela interface {
	// Konto zwraca konto właściciela; brak wiersza znaczy platformę przed rejestracją.
	Konto(ctx context.Context) (KontoWlasciciela, error)
	// ZalozKonto zapisuje jedyne konto platformy; drugie założenie odbija się o warunek schematu.
	ZalozKonto(ctx context.Context, konto KontoWlasciciela) error
	// PotwierdzKonto przenosi konto ze stanu niepotwierdzonego do potwierdzonego.
	PotwierdzKonto(ctx context.Context) error
	// UsunKonto kasuje konto właściciela — istnieje wyłącznie po to, by cofnąć nieudaną rejestrację.
	UsunKonto(ctx context.Context) error

	// ZalozPotwierdzenie zapisuje skrót drogi potwierdzenia wraz z celem
	// i czasem wygaśnięcia.
	ZalozPotwierdzenie(ctx context.Context, p PotwierdzenieTozsamosci) error
	// PotwierdzeniePoSkrocie zwraca drogę rozpoznaną skrótem; brak wiersza znaczy drogę nieznaną.
	PotwierdzeniePoSkrocie(ctx context.Context, skrot string) (PotwierdzenieTozsamosci, error)
	// ZuzyjPotwierdzenie zamyka drogę po użyciu; drugi wynik mówi, czy wiersz dało się zamknąć.
	ZuzyjPotwierdzenie(ctx context.Context, skrot string, teraz int64) (bool, error)
}

const (
	kolumnyKontaWlasciciela = `login, email, potwierdzone, utworzono`

	kontoWlascicielaWiersz = `SELECT ` + kolumnyKontaWlasciciela +
		` FROM konto_wlasciciela WHERE id = 1`

	wstawKontoWlasciciela = `INSERT INTO konto_wlasciciela
	                         (id, login, email, potwierdzone, utworzono)
	                         VALUES (1, ?, ?, ?, ?)`

	potwierdzKontoWlasciciela = `UPDATE konto_wlasciciela SET potwierdzone = 1 WHERE id = 1`

	usunKontoWlasciciela = `DELETE FROM konto_wlasciciela WHERE id = 1`

	wstawPotwierdzenieTozsamosci = `INSERT INTO potwierdzenie_tozsamosci
	                                (skrot, cel, wygasa, uzyte, utworzono)
	                                VALUES (?, ?, ?, 0, ?)`

	potwierdzenieTozsamosciPoSkrocie = `SELECT skrot, cel, wygasa, uzyte, utworzono
	                                    FROM potwierdzenie_tozsamosci WHERE skrot = ?`

	// Zamknięcie drogi jest warunkowe: `uzyte = 0` w klauzuli WHERE sprawia, że
	// dwa równoległe żądania z tym samym materiałem dają jedno zamknięcie i jedno
	// zero zmienionych wierszy. Niepodzielność stoi w bazie, nie w kodzie.
	zuzyjPotwierdzenieTozsamosci = `UPDATE potwierdzenie_tozsamosci
	                                SET uzyte = 1
	                                WHERE skrot = ? AND uzyte = 0 AND wygasa > ?`
)

type repozytoriumKontaWlasciciela struct {
	zapytania *zapytania
}

// Zgodność implementacji repozytorium konta właściciela z kontraktem jest sprawdzana przy kompilacji pakietu.
var _ RepozytoriumKontaWlasciciela = (*repozytoriumKontaWlasciciela)(nil)

func noweRepozytoriumKontaWlasciciela(z *zapytania) *repozytoriumKontaWlasciciela {
	return &repozytoriumKontaWlasciciela{zapytania: z}
}

// Konto zwraca jedyne konto platformy, zapisane w tabeli konta właściciela w bazie danych rdzenia systemu.
func (r *repozytoriumKontaWlasciciela) Konto(ctx context.Context) (KontoWlasciciela, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, kontoWlascicielaWiersz)
	if err != nil {
		return KontoWlasciciela{}, err
	}
	var konto KontoWlasciciela
	var potwierdzone int64
	err = polecenie.QueryRowContext(ctx).Scan(
		&konto.Login, &konto.Email, &potwierdzone, &konto.Utworzono)
	if errors.Is(err, sql.ErrNoRows) {
		return KontoWlasciciela{}, fmt.Errorf(
			"dane: konta właściciela nie ma: %w", ErrBrakWiersza)
	}
	if err != nil {
		return KontoWlasciciela{}, fmt.Errorf("dane: nie można odczytać konta właściciela: %w", err)
	}
	konto.Potwierdzone = potwierdzone != 0
	return konto, nil
}

// ZalozKonto zapisuje jedyne konto platformy w tabeli konta właściciela w bazie danych systemu rdzenia.
func (r *repozytoriumKontaWlasciciela) ZalozKonto(ctx context.Context, konto KontoWlasciciela) error {
	polecenie, err := r.zapytania.przygotuj(ctx, wstawKontoWlasciciela)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, konto.Login, konto.Email,
		liczbaLogiczna(konto.Potwierdzone), konto.Utworzono)
	if czyKolizja(err) {
		return fmt.Errorf("dane: konto właściciela już istnieje: %w", ErrKolizjaWiersza)
	}
	if err != nil {
		return fmt.Errorf("dane: nie można założyć konta właściciela: %w", err)
	}
	return nil
}

// PotwierdzKonto przenosi konto właściciela do stanu potwierdzonego po weryfikacji adresu rejestracji.
func (r *repozytoriumKontaWlasciciela) PotwierdzKonto(ctx context.Context) error {
	polecenie, err := r.zapytania.przygotuj(ctx, potwierdzKontoWlasciciela)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx); err != nil {
		return fmt.Errorf("dane: nie można potwierdzić konta właściciela: %w", err)
	}
	return nil
}

// UsunKonto kasuje konto właściciela — droga cofnięcia nieudanej rejestracji platformy bez potwierdzenia.
func (r *repozytoriumKontaWlasciciela) UsunKonto(ctx context.Context) error {
	polecenie, err := r.zapytania.przygotuj(ctx, usunKontoWlasciciela)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx); err != nil {
		return fmt.Errorf("dane: nie można usunąć konta właściciela: %w", err)
	}
	return nil
}

// ZalozPotwierdzenie zapisuje skrót drogi potwierdzenia tożsamości wraz z terminem jego ważności czasowej.
func (r *repozytoriumKontaWlasciciela) ZalozPotwierdzenie(ctx context.Context,
	p PotwierdzenieTozsamosci) error {

	polecenie, err := r.zapytania.przygotuj(ctx, wstawPotwierdzenieTozsamosci)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, p.Skrot, p.Cel, p.Wygasa, p.Utworzono); err != nil {
		return fmt.Errorf("dane: nie można założyć drogi potwierdzenia: %w", err)
	}
	return nil
}

// PotwierdzeniePoSkrocie zwraca drogę potwierdzenia tożsamości rozpoznaną jej skrótem zapisanym w bazie.
func (r *repozytoriumKontaWlasciciela) PotwierdzeniePoSkrocie(ctx context.Context,
	skrot string) (PotwierdzenieTozsamosci, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, potwierdzenieTozsamosciPoSkrocie)
	if err != nil {
		return PotwierdzenieTozsamosci{}, err
	}
	var p PotwierdzenieTozsamosci
	var uzyte int64
	err = polecenie.QueryRowContext(ctx, skrot).Scan(
		&p.Skrot, &p.Cel, &p.Wygasa, &uzyte, &p.Utworzono)
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

// ZuzyjPotwierdzenie zamyka drogę potwierdzenia tożsamości natychmiast po jej pierwszym wykorzystaniu.
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
