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

// KontoWlasciciela to wiersz tabeli `konto_wlasciciela`. Kont jest tyle, ilu
// użytkowników: jednoznaczność pilnują wskaźniki na loginie i adresie, nie
// identyfikator wiersza (migracja 406, `decyzje.md` poz. 22).
type KontoWlasciciela struct {
	// Id wiersza konta. Zero znaczy konto jeszcze niezapisane.
	Id           int64
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
	// KontoId wskazuje konto, do którego droga prowadzi. Bez tego wskazania
	// potwierdzenie adresu przy dwóch kontach naraz otwierałoby konto najstarsze,
	// nie to, na które poszedł list.
	KontoId int64
}

// RepozytoriumKontaWlasciciela jest kontraktem trwałości rejestracji.
//
// Byt jest odrębny od katalogu metod wejścia, bo odpowiada na inne pytanie:
// katalog metod mówi, CZYM otworzyć bramkę, konto mówi, CZYJA ona jest.
type RepozytoriumKontaWlasciciela interface {
	// Konto zwraca konto założone najwcześniej; brak wiersza znaczy platformę przed rejestracją.
	Konto(ctx context.Context) (KontoWlasciciela, error)
	// KontoPoTozsamosci odnajduje konto po loginie albo adresie, bez względu na
	// wielkość liter. Brak wiersza znaczy tożsamość nieznaną, nie platformę pustą.
	KontoPoTozsamosci(ctx context.Context, wskazanie string) (KontoWlasciciela, error)
	// KontoPoId odnajduje konto wskazane identyfikatorem. Droga potwierdzenia
	// niesie identyfikator konta, do którego prowadzi, i tylko po nim wolno
	// sięgnąć po konto — inaczej list otwierałby konto najstarsze.
	KontoPoId(ctx context.Context, kontoId int64) (KontoWlasciciela, error)
	// ZalozKonto zapisuje konto i zwraca jego identyfikator; zajęty login albo adres daje ErrKolizjaWiersza.
	ZalozKonto(ctx context.Context, konto KontoWlasciciela) (int64, error)
	// PotwierdzKonto przenosi wskazane konto ze stanu niepotwierdzonego do potwierdzonego.
	PotwierdzKonto(ctx context.Context, kontoId int64) error
	// UsunKonto kasuje wskazane konto — istnieje wyłącznie po to, by cofnąć nieudaną rejestrację.
	UsunKonto(ctx context.Context, kontoId int64) error

	// ZalozPotwierdzenie zapisuje skrót drogi potwierdzenia wraz z celem
	// i czasem wygaśnięcia.
	ZalozPotwierdzenie(ctx context.Context, p PotwierdzenieTozsamosci) error
	// PotwierdzeniePoSkrocie zwraca drogę rozpoznaną skrótem; brak wiersza znaczy drogę nieznaną.
	PotwierdzeniePoSkrocie(ctx context.Context, skrot string) (PotwierdzenieTozsamosci, error)
	// ZuzyjPotwierdzenie zamyka drogę po użyciu; drugi wynik mówi, czy wiersz dało się zamknąć.
	ZuzyjPotwierdzenie(ctx context.Context, skrot string, teraz int64) (bool, error)
}

const (
	kolumnyKontaWlasciciela = `id, login, email, potwierdzone, utworzono`

	// Konto najstarsze, gdy pyta się bez wskazania tożsamości — kolejność po
	// identyfikatorze, bo znacznik czasu dwóch kont założonych w tej samej
	// milisekundzie nie rozstrzyga.
	kontoWlascicielaWiersz = `SELECT ` + kolumnyKontaWlasciciela +
		` FROM konto_wlasciciela ORDER BY id LIMIT 1`

	kontoWlascicielaPoId = `SELECT ` + kolumnyKontaWlasciciela +
		` FROM konto_wlasciciela WHERE id = ?`

	// Jedno zapytanie na login i na adres: okno logowania przyjmuje oba w tym
	// samym polu i nie rozstrzyga, które podano.
	kontoWlascicielaPoTozsamosci = `SELECT ` + kolumnyKontaWlasciciela +
		` FROM konto_wlasciciela
		  WHERE login = ? COLLATE NOCASE OR email = ? COLLATE NOCASE
		  ORDER BY id LIMIT 1`

	wstawKontoWlasciciela = `INSERT INTO konto_wlasciciela
	                         (login, email, potwierdzone, utworzono)
	                         VALUES (?, ?, ?, ?)`

	potwierdzKontoWlasciciela = `UPDATE konto_wlasciciela SET potwierdzone = 1 WHERE id = ?`

	usunKontoWlasciciela = `DELETE FROM konto_wlasciciela WHERE id = ?`

	wstawPotwierdzenieTozsamosci = `INSERT INTO potwierdzenie_tozsamosci
	                                (skrot, cel, wygasa, uzyte, utworzono, konto_id)
	                                VALUES (?, ?, ?, 0, ?, NULLIF(?, 0))`

	potwierdzenieTozsamosciPoSkrocie = `SELECT skrot, cel, wygasa, uzyte, utworzono,
	                                           COALESCE(konto_id, 0)
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

// odczytajKontoWlasciciela rozbiera wiersz konta z gotowego zapytania; wspólne dla odczytu
// bez wskazania tożsamości i z jej wskazaniem, żeby dwie ścieżki nie rozjechały
// się w kolejności kolumn.
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

// Konto zwraca konto założone najwcześniej. Służy pytaniu „czy platforma ma
// w ogóle konto", nie wskazaniu, czyja jest bieżąca sesja.
func (r *repozytoriumKontaWlasciciela) Konto(ctx context.Context) (KontoWlasciciela, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, kontoWlascicielaWiersz)
	if err != nil {
		return KontoWlasciciela{}, err
	}
	return odczytajKontoWlasciciela(polecenie.QueryRowContext(ctx))
}

// KontoPoTozsamosci odnajduje konto po loginie albo adresie e-mail.
func (r *repozytoriumKontaWlasciciela) KontoPoTozsamosci(ctx context.Context,
	wskazanie string) (KontoWlasciciela, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, kontoWlascicielaPoTozsamosci)
	if err != nil {
		return KontoWlasciciela{}, err
	}
	return odczytajKontoWlasciciela(polecenie.QueryRowContext(ctx, wskazanie, wskazanie))
}

// KontoPoId odnajduje konto wskazane identyfikatorem drogi potwierdzenia.
func (r *repozytoriumKontaWlasciciela) KontoPoId(ctx context.Context,
	kontoId int64) (KontoWlasciciela, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, kontoWlascicielaPoId)
	if err != nil {
		return KontoWlasciciela{}, err
	}
	return odczytajKontoWlasciciela(polecenie.QueryRowContext(ctx, kontoId))
}

// ZalozKonto zapisuje konto i zwraca jego identyfikator. Kolizja znaczy zajęty
// login albo zajęty adres — pilnują tego wskaźniki jednoznaczności, nie kod.
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

// PotwierdzKonto przenosi konto właściciela do stanu potwierdzonego po weryfikacji adresu rejestracji.
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

// UsunKonto kasuje konto właściciela — droga cofnięcia nieudanej rejestracji platformy bez potwierdzenia.
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

// ZalozPotwierdzenie zapisuje skrót drogi potwierdzenia tożsamości wraz z terminem jego ważności czasowej.
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
