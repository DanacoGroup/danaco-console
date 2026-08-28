// Odpowiedzialność pliku: słownik skrótów tekstowych, rodzina komend snippet.*; skrót rozwija się we wszystkich polach tekstowych platformy.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// SkrotTekstowy to wiersz tabeli skrot_tekstowy, niosący kod, treść rozwinięcia i profil tego samego skrótu.
type SkrotTekstowy struct {
	Kod            string
	ProfilKod      string
	Skrot          string
	Tresc          string
	Opis           *string
	PolaJSON       string
	Czynny         bool
	Utworzono      int64
	Zaktualizowano int64
}

// RepozytoriumSkrotow jest kontraktem słownika skrótów tekstowych dla wszystkich warstw wyższych rdzenia.
type RepozytoriumSkrotow interface {
	ZapiszSkrotTekstowy(ctx context.Context, skrot SkrotTekstowy) (SkrotTekstowy, error)
	SkrotTekstowyPoKodzie(ctx context.Context, kod string) (SkrotTekstowy, error)
	SkrotyTekstowe(ctx context.Context, fraza, profil string, granica int) ([]SkrotTekstowy, int, error)
	UsunSkrotTekstowy(ctx context.Context, kod string) (bool, error)
}

const (
	kolumnySkrotuTekstowego = `identyfikator_zewnetrzny, profil_kod, skrot, tresc, opis,
	                           pola_json, czynny, utworzono, zaktualizowano`

	zapiszSkrotTekstowy = `INSERT INTO skrot_tekstowy (` + kolumnySkrotuTekstowego + `)
	                       VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	                       ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                           profil_kod = excluded.profil_kod,
	                           skrot = excluded.skrot,
	                           tresc = excluded.tresc,
	                           opis = excluded.opis,
	                           pola_json = excluded.pola_json,
	                           czynny = excluded.czynny,
	                           zaktualizowano = excluded.zaktualizowano`

	pobierzSkrotTekstowy = `SELECT ` + kolumnySkrotuTekstowego +
		` FROM skrot_tekstowy WHERE identyfikator_zewnetrzny = ?`

	usunSkrotTekstowy = `DELETE FROM skrot_tekstowy WHERE identyfikator_zewnetrzny = ?`
)

type repozytoriumSkrotow struct {
	zapytania *zapytania
	db        *sql.DB
}

// noweRepozytoriumSkrotow zakłada słownik skrótów nad wspólną bazą tego samego zestawu, gotowy do użycia.
func noweRepozytoriumSkrotow(z *zapytania, db *sql.DB) *repozytoriumSkrotow {
	return &repozytoriumSkrotow{zapytania: z, db: db}
}

// ZapiszSkrotTekstowy zakłada skrót albo nadpisuje zastany po kodzie, zwracając jego pełny stan po zapisie.
func (r *repozytoriumSkrotow) ZapiszSkrotTekstowy(ctx context.Context,
	skrot SkrotTekstowy) (SkrotTekstowy, error) {

	if strings.TrimSpace(skrot.Kod) == "" {
		return SkrotTekstowy{}, fmt.Errorf("dane: skrót tekstowy bez identyfikatora")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszSkrotTekstowy)
	if err != nil {
		return SkrotTekstowy{}, err
	}
	_, err = polecenie.ExecContext(ctx, skrot.Kod, skrot.ProfilKod, skrot.Skrot, skrot.Tresc,
		tekstDoKolumny(skrot.Opis), skrot.PolaJSON, liczbaLogiczna(skrot.Czynny),
		skrot.Utworzono, skrot.Zaktualizowano)
	if err != nil {
		return SkrotTekstowy{}, fmt.Errorf("dane: nie można zapisać skrótu tekstowego %q: %w",
			skrot.Kod, err)
	}
	return r.SkrotTekstowyPoKodzie(ctx, skrot.Kod)
}

// SkrotTekstowyPoKodzie zwraca jedną pozycję słownika po jej kodzie zewnętrznym, wraz z jej pełną treścią.
func (r *repozytoriumSkrotow) SkrotTekstowyPoKodzie(ctx context.Context,
	kod string) (SkrotTekstowy, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzSkrotTekstowy)
	if err != nil {
		return SkrotTekstowy{}, err
	}
	skrot, err := odczytajSkrotTekstowy(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return SkrotTekstowy{}, fmt.Errorf("dane: skrót tekstowy %q nie istnieje: %w",
			kod, ErrBrakWiersza)
	}
	return skrot, err
}

// SkrotyTekstowe zwraca słownik w kolejności alfabetycznej wraz z liczbą
// pozycji spełniających zawężenie.
func (r *repozytoriumSkrotow) SkrotyTekstowe(ctx context.Context, fraza, profil string,
	granica int) ([]SkrotTekstowy, int, error) {

	warunki := []string{"1 = 1"}
	argumenty := []any{}
	if szukane := strings.TrimSpace(fraza); szukane != "" {
		warunki = append(warunki, "(skrot LIKE ? OR tresc LIKE ?)")
		argumenty = append(argumenty, "%"+szukane+"%", "%"+szukane+"%")
	}
	if strings.TrimSpace(profil) != "" {
		warunki = append(warunki, "profil_kod = ?")
		argumenty = append(argumenty, profil)
	}
	gdzie := strings.Join(warunki, " AND ")

	wszystkich := 0
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM skrot_tekstowy WHERE `+gdzie,
		argumenty...).Scan(&wszystkich); err != nil {
		return nil, 0, fmt.Errorf("dane: nie można policzyć skrótów tekstowych: %w", err)
	}

	tekst := `SELECT ` + kolumnySkrotuTekstowego + ` FROM skrot_tekstowy WHERE ` + gdzie +
		` ORDER BY skrot, profil_kod`
	if granica > 0 {
		tekst += fmt.Sprintf(" LIMIT %d", granica)
	}
	wiersze, err := r.db.QueryContext(ctx, tekst, argumenty...)
	if err != nil {
		return nil, 0, fmt.Errorf("dane: nie można odczytać skrótów tekstowych: %w", err)
	}
	defer wiersze.Close()

	lista := []SkrotTekstowy{}
	for wiersze.Next() {
		skrot, err := odczytajSkrotTekstowy(wiersze)
		if err != nil {
			return nil, 0, err
		}
		lista = append(lista, skrot)
	}
	if err := wiersze.Err(); err != nil {
		return nil, 0, fmt.Errorf("dane: przerwany odczyt skrótów tekstowych: %w", err)
	}
	return lista, wszystkich, nil
}

// UsunSkrotTekstowy kasuje pozycję słownika. Brak wiersza nie jest awarią —
// oddaje fałsz, a kontrakt niesie to polem `deleted`.
func (r *repozytoriumSkrotow) UsunSkrotTekstowy(ctx context.Context, kod string) (bool, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, usunSkrotTekstowy)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, kod)
	if err != nil {
		return false, fmt.Errorf("dane: nie można usunąć skrótu tekstowego %q: %w", kod, err)
	}
	usuniete, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nie można policzyć usuniętych skrótów: %w", err)
	}
	return usuniete > 0, nil
}

// odczytajSkrotTekstowy przekłada wiersz tabeli na pozycję słownika, kolumna po kolumnie tego zapytania.
func odczytajSkrotTekstowy(s skaner) (SkrotTekstowy, error) {
	var skrot SkrotTekstowy
	var opis sql.NullString
	var czynny int
	err := s.Scan(&skrot.Kod, &skrot.ProfilKod, &skrot.Skrot, &skrot.Tresc, &opis,
		&skrot.PolaJSON, &czynny, &skrot.Utworzono, &skrot.Zaktualizowano)
	if err != nil {
		return SkrotTekstowy{}, err
	}
	skrot.Opis = tekstZKolumny(opis)
	skrot.Czynny = czynny == 1
	return skrot, nil
}
