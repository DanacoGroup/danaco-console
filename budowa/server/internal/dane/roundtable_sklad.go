// Odpowiedzialność pliku: skład debaty poza dopisaniem uczestnika: zmiana tożsamości, usunięcie ze składu, zespoły do ponownego użycia i role.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// UczestnikZespoluDebaty to jeden wiersz zapamiętanego składu. Nie ma okna ani
// wyciszenia: jedno i drugie należy do debaty, nie do zespołu.
type UczestnikZespoluDebaty struct {
	KanalModelu     string
	NazwaTozsamosci *string
	PromptSystemowy *string
	Rola            string
	Waga            float64
	Awatar          *string
	OpisRoli        *string
	Kolejnosc       int
}

// ZespolDebaty to nazwany skład wraz z formatem debaty, jeśli ten format zapisano razem z tym składem.
type ZespolDebaty struct {
	Kod        string
	Nazwa      string
	Format     string
	Utworzono  string
	Uczestnicy []UczestnikZespoluDebaty
}

// RolaDebaty to pozycja biblioteki ról dostępnych do przypisania uczestnikom każdej prowadzonej debaty.
type RolaDebaty struct {
	Kod             string
	Nazwa           string
	PromptSystemowy string
	Opis            *string
	Fabryczna       bool
}

// RepozytoriumDebatySkladu jest częścią kontraktu obszaru odpowiadającą za
// skład: uczestnika po zmianie, zespoły i role.
type RepozytoriumDebatySkladu interface {
	ZmienUczestnika(ctx context.Context, uczestnik UczestnikDebaty) error
	UsunUczestnika(ctx context.Context, kod string) error

	ZapiszZespol(ctx context.Context, zespol ZespolDebaty) (ZespolDebaty, error)
	Zespol(ctx context.Context, kod string) (ZespolDebaty, error)
	Zespoly(ctx context.Context, fraza string, limit int) ([]ZespolDebaty, error)

	Role(ctx context.Context, fraza string) ([]RolaDebaty, error)
}

const (
	// Zmiana tożsamości nie dotyka wyciszenia ani kolejności: obie są
	// czynnościami moderatora, a `roundtable.model.update` opisuje uczestnika,
	// nie przebieg tury.
	zmienUczestnikaDebaty = `UPDATE debata_uczestnik
	                         SET nazwa_tozsamosci = ?, prompt_systemowy = ?, kluczowy = ?,
	                             waga = ?, rola = ?, agent = ?, awatar = ?, opis_roli = ?,
	                             liczba_probek = ?,
	                             zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                         WHERE identyfikator_zewnetrzny = ?`

	usunUczestnikaDebaty = `DELETE FROM debata_uczestnik WHERE identyfikator_zewnetrzny = ?`

	zalozZespolDebaty = `INSERT INTO debata_zespol (identyfikator_zewnetrzny, nazwa, format)
	                     VALUES (?, ?, ?)`

	zalozUczestnikaZespoluDebaty = `INSERT INTO debata_zespol_uczestnik
	                          (zespol_id, kanal_modelu, nazwa_tozsamosci, prompt_systemowy,
	                           rola, waga, awatar, opis_roli, kolejnosc)
	                          SELECT z.id, ?, ?, ?, ?, ?, ?, ?, ?
	                            FROM debata_zespol z
	                           WHERE z.identyfikator_zewnetrzny = ?`

	pobierzZespolDebaty = `SELECT identyfikator_zewnetrzny, nazwa, format, utworzono
	                 FROM debata_zespol WHERE identyfikator_zewnetrzny = ?`

	// Fraza pusta przepuszcza wszystko: warunek porównuje z wzorcem procentowym, któremu odpowiada każda nazwa uczestnika.
	pobierzZespolyDebaty = `SELECT identyfikator_zewnetrzny, nazwa, format, utworzono
	                  FROM debata_zespol
	                  WHERE nazwa LIKE '%' || ? || '%'
	                  ORDER BY id DESC LIMIT (CASE WHEN ? > 0 THEN ? ELSE -1 END)`

	pobierzUczestnikowZespoluDebaty = `SELECT u.kanal_modelu, u.nazwa_tozsamosci, u.prompt_systemowy,
	                                    u.rola, u.waga, u.awatar, u.opis_roli, u.kolejnosc
	                               FROM debata_zespol_uczestnik u
	                               JOIN debata_zespol z ON z.id = u.zespol_id
	                              WHERE z.identyfikator_zewnetrzny = ?
	                              ORDER BY u.kolejnosc ASC, u.id ASC`

	pobierzRoleDebaty = `SELECT identyfikator_zewnetrzny, nazwa, prompt_systemowy, opis, fabryczna
	               FROM debata_rola
	               WHERE nazwa LIKE '%' || ? || '%'
	               ORDER BY fabryczna DESC, nazwa ASC`
)

// ZmienUczestnika zapisuje tożsamość, rolę i wagę głosu wskazanego uczestnika tej samej prowadzonej debaty.
func (r *repozytoriumRoundtable) ZmienUczestnika(ctx context.Context, uczestnik UczestnikDebaty) error {
	polecenie, err := r.zapytania.przygotuj(ctx, zmienUczestnikaDebaty)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, uczestnik.NazwaTozsamosci, uczestnik.PromptSystemowy,
		uczestnik.Kluczowy, uczestnik.Waga, uczestnik.Rola, uczestnik.Agent, uczestnik.Awatar,
		uczestnik.OpisRoli, uczestnik.LiczbaProbek, uczestnik.Kod)
	if err != nil {
		return fmt.Errorf("dane: nie można zmienić uczestnika debaty %q: %w", uczestnik.Kod, err)
	}
	return trafienieDebaty(wynik)
}

// UsunUczestnika zdejmuje uczestnika ze składu. Jego wypowiedzi zostają: zapis
// tury nie ma prawa się zmienić dlatego, że mówca wypadł ze składu.
func (r *repozytoriumRoundtable) UsunUczestnika(ctx context.Context, kod string) error {
	polecenie, err := r.zapytania.przygotuj(ctx, usunUczestnikaDebaty)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, kod)
	if err != nil {
		return fmt.Errorf("dane: nie można usunąć uczestnika debaty %q: %w", kod, err)
	}
	return trafienieDebaty(wynik)
}

// ZapiszZespol utrwala zespół razem ze składem w jednej transakcji. Zespół
// z połową składu byłby układem, którego nikt nie zapisywał.
func (r *repozytoriumRoundtable) ZapiszZespol(ctx context.Context,
	zespol ZespolDebaty) (ZespolDebaty, error) {

	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		naglowek, err := r.zapytania.wTransakcji(ctx, transakcja, zalozZespolDebaty)
		if err != nil {
			return err
		}
		if _, err := naglowek.ExecContext(ctx, zespol.Kod, zespol.Nazwa, zespol.Format); err != nil {
			return fmt.Errorf("dane: nie można założyć zespołu debaty %q: %w", zespol.Kod, err)
		}
		wiersz, err := r.zapytania.wTransakcji(ctx, transakcja, zalozUczestnikaZespoluDebaty)
		if err != nil {
			return err
		}
		for pozycja, uczestnik := range zespol.Uczestnicy {
			if _, err := wiersz.ExecContext(ctx, uczestnik.KanalModelu, uczestnik.NazwaTozsamosci,
				uczestnik.PromptSystemowy, uczestnik.Rola, uczestnik.Waga, uczestnik.Awatar,
				uczestnik.OpisRoli, pozycja+1, zespol.Kod); err != nil {

				return fmt.Errorf("dane: nie można zapisać uczestnika zespołu %q: %w", zespol.Kod, err)
			}
		}
		return nil
	})
	if err != nil {
		return ZespolDebaty{}, err
	}
	return r.Zespol(ctx, zespol.Kod)
}

// Zespol zwraca zespół wraz z jego zapamiętanym pełnym składem uczestników tej samej debaty operacyjnej.
func (r *repozytoriumRoundtable) Zespol(ctx context.Context, kod string) (ZespolDebaty, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzZespolDebaty)
	if err != nil {
		return ZespolDebaty{}, err
	}
	var zespol ZespolDebaty
	err = polecenie.QueryRowContext(ctx, kod).Scan(&zespol.Kod, &zespol.Nazwa,
		&zespol.Format, &zespol.Utworzono)
	if errors.Is(err, sql.ErrNoRows) {
		return ZespolDebaty{}, ErrBrakWiersza
	}
	if err != nil {
		return ZespolDebaty{}, fmt.Errorf("dane: nieczytelny wiersz zespołu debaty %q: %w", kod, err)
	}
	zespol.Uczestnicy, err = r.uczestnicyZespolu(ctx, kod)
	if err != nil {
		return ZespolDebaty{}, err
	}
	return zespol, nil
}

// Zespoly zwraca zespoły tego okna od najnowszego, wraz ze składem uczestników każdego z tych zespołów.
func (r *repozytoriumRoundtable) Zespoly(ctx context.Context,
	fraza string, limit int) ([]ZespolDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzZespolyDebaty)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, strings.TrimSpace(fraza), limit, limit)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać zespołów debaty: %w", err)
	}
	defer wiersze.Close()

	zespoly := make([]ZespolDebaty, 0, 8)
	for wiersze.Next() {
		var zespol ZespolDebaty
		if err := wiersze.Scan(&zespol.Kod, &zespol.Nazwa, &zespol.Format,
			&zespol.Utworzono); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz zespołu debaty: %w", err)
		}
		zespoly = append(zespoly, zespol)
	}
	if err := wiersze.Err(); err != nil {
		return nil, err
	}
	for i := range zespoly {
		skladu, err := r.uczestnicyZespolu(ctx, zespoly[i].Kod)
		if err != nil {
			return nil, err
		}
		zespoly[i].Uczestnicy = skladu
	}
	return zespoly, nil
}

// uczestnicyZespolu czyta zapamiętany skład jednego zespołu tej samej debaty, w kolejności jego zapisania.
func (r *repozytoriumRoundtable) uczestnicyZespolu(ctx context.Context,
	kod string) ([]UczestnikZespoluDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzUczestnikowZespoluDebaty)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, kod)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać składu zespołu %q: %w", kod, err)
	}
	defer wiersze.Close()

	skladu := make([]UczestnikZespoluDebaty, 0, 8)
	for wiersze.Next() {
		var uczestnik UczestnikZespoluDebaty
		if err := wiersze.Scan(&uczestnik.KanalModelu, &uczestnik.NazwaTozsamosci,
			&uczestnik.PromptSystemowy, &uczestnik.Rola, &uczestnik.Waga, &uczestnik.Awatar,
			&uczestnik.OpisRoli, &uczestnik.Kolejnosc); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz składu zespołu: %w", err)
		}
		skladu = append(skladu, uczestnik)
	}
	return skladu, wiersze.Err()
}

// Role zwraca całą bibliotekę ról tej debaty operacyjnej; role fabryczne idą przed rolami własnymi Operatora.
func (r *repozytoriumRoundtable) Role(ctx context.Context, fraza string) ([]RolaDebaty, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzRoleDebaty)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, strings.TrimSpace(fraza))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać ról debaty: %w", err)
	}
	defer wiersze.Close()

	role := make([]RolaDebaty, 0, 8)
	for wiersze.Next() {
		var rola RolaDebaty
		if err := wiersze.Scan(&rola.Kod, &rola.Nazwa, &rola.PromptSystemowy, &rola.Opis,
			&rola.Fabryczna); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz roli debaty: %w", err)
		}
		role = append(role, rola)
	}
	return role, wiersze.Err()
}
