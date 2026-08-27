// Trwałość konsultacji, dziennik tabeli konsultacja_doradcy, stoi osobno od
// pojęcia doradcy i od doboru doradcy: trzy odpowiedzialności, trzy pliki.
package podagenci

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// Stany wpisu dziennika, przepisane z CHECK-u tabeli: doradca odpowiedział albo
// konsultacja się nie odbyła (powód w kolumnie `powod`).
const (
	StanRady       = "rada"
	StanOdmowyRady = "odmowa"
)

const limitWykazuKonsultacji = 50 // obowiązuje przy filtrze bez własnej granicy

// Konsultacja jest wierszem tabeli `konsultacja_doradcy` — pytaniem, radą
// i opisem wywołania w jednym śladzie.
type Konsultacja struct {
	ID                         int64
	Identyfikator              string
	OknoId                     string
	Pytajacy                   string
	DoradcaKanal, DoradcaModel string
	Pytanie, Rada, Skrot       string
	// Prowenancja to fragment strumienia przepisany co do znaku, żeby ślad
	// w bazie się nie rozjechał.
	Prowenancja string
	Stan        string
	Powod       string
	Utworzono   int64
}

// FiltrKonsultacji zawęża wykaz dziennika konsultacji doradcy; limit
// niedodatni znaczy limit domyślny wykazu.
type FiltrKonsultacji struct {
	OknoId string
	Limit  int
}

// Dziennik jest kontraktem trwałości konsultacji doradcy: wołający zna
// interfejs zapisu i wykazu, nie zapytania SQL.
type Dziennik interface {
	Zapisz(ctx context.Context, wpis Konsultacja) (Konsultacja, error)
	Wykaz(ctx context.Context, filtr FiltrKonsultacji) ([]Konsultacja, error)
}

const (
	kolumnyKonsultacji = `k.id, k.identyfikator_zewnetrzny, k.okno_id, k.pytajacy_kanal,
	                      k.doradca_kanal, k.doradca_model, k.pytanie, k.rada, k.skrot,
	                      k.prowenancja, k.stan, k.powod, k.utworzono`

	zrodloKonsultacji  = ` FROM konsultacja_doradcy k`
	pobierzKonsultacje = `SELECT ` + kolumnyKonsultacji + zrodloKonsultacji +
		` WHERE k.identyfikator_zewnetrzny = ?`
	listaKonsultacji = `SELECT ` + kolumnyKonsultacji + zrodloKonsultacji +
		` WHERE (? = '' OR k.okno_id = ?)
		  ORDER BY k.utworzono DESC, k.id DESC
		  LIMIT ?`
	wstawKonsultacje = `INSERT INTO konsultacja_doradcy
	                    (identyfikator_zewnetrzny, okno_id, pytajacy_kanal, doradca_kanal,
	                     doradca_model, pytanie, rada, skrot, prowenancja, stan, powod,
	                     utworzono)
	                    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
)

// dziennikKonsultacji jest jedyną implementacją Dziennika; zgodność sprawdza
// kompilator, a nie pierwsze wywołanie w czasie pracy.
type dziennikKonsultacji struct{ db *sql.DB }

var _ Dziennik = (*dziennikKonsultacji)(nil)

// NowyDziennik zakłada dziennik konsultacji nad otwartą pulą połączeń. Baza
// pusta daje nil — konsultacja odbywa się i tak, traci wyłącznie ślad w bazie;
// ślad w strumieniu zostaje.
func NowyDziennik(db *sql.DB) Dziennik {
	if db == nil {
		return nil
	}
	return &dziennikKonsultacji{db: db}
}

// Zapisz wstawia dany wpis konsultacji do dziennika tej bazy i oddaje go
// odczytany z powrotem po zapisie.
func (d *dziennikKonsultacji) Zapisz(ctx context.Context, wpis Konsultacja) (Konsultacja, error) {
	if err := sprawdzKonsultacje(wpis); err != nil {
		return Konsultacja{}, err
	}
	if _, err := d.db.ExecContext(ctx, wstawKonsultacje, wpis.Identyfikator,
		kolumnaNapisu(wpis.OknoId), wpis.Pytajacy, wpis.DoradcaKanal, wpis.DoradcaModel,
		wpis.Pytanie, kolumnaNapisu(wpis.Rada), kolumnaNapisu(wpis.Skrot),
		kolumnaNapisu(wpis.Prowenancja), wpis.Stan, kolumnaNapisu(wpis.Powod),
		wpis.Utworzono); err != nil {

		return Konsultacja{}, fmt.Errorf("podagenci: nie można zapisać konsultacji %q: %w",
			wpis.Identyfikator, err)
	}
	return d.wpisPoIdentyfikatorze(ctx, wpis.Identyfikator)
}

// Wykaz oddaje wpisy dziennika konsultacji doradcy od najnowszego, zawężone
// filtrem okna i limitem wykazu.
func (d *dziennikKonsultacji) Wykaz(ctx context.Context, filtr FiltrKonsultacji) ([]Konsultacja, error) {
	limit := filtr.Limit
	if limit <= 0 {
		limit = limitWykazuKonsultacji
	}
	wiersze, err := d.db.QueryContext(ctx, listaKonsultacji, filtr.OknoId, filtr.OknoId, limit)
	if err != nil {
		return nil, fmt.Errorf("podagenci: nie można odczytać dziennika konsultacji: %w", err)
	}
	defer wiersze.Close()

	dziennik := make([]Konsultacja, 0, 16)
	for wiersze.Next() {
		wpis, err := odczytajKonsultacje(wiersze)
		if err != nil {
			return nil, fmt.Errorf("podagenci: uszkodzony wiersz konsultacji: %w", err)
		}
		dziennik = append(dziennik, wpis)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("podagenci: przerwany odczyt dziennika konsultacji: %w", err)
	}
	return dziennik, nil
}

// wpisPoIdentyfikatorze odczytuje jeden wpis dziennika konsultacji po jego
// tożsamości zewnętrznej wpisu.
func (d *dziennikKonsultacji) wpisPoIdentyfikatorze(ctx context.Context, identyfikator string) (Konsultacja, error) {
	wiersz := d.db.QueryRowContext(ctx, pobierzKonsultacje, identyfikator)
	wpis, err := odczytajKonsultacje(wiersz)
	if err != nil {
		return Konsultacja{}, fmt.Errorf("podagenci: nie można odczytać konsultacji %q: %w",
			identyfikator, err)
	}
	return wpis, nil
}

// skaner jest bliźniakiem nieeksportowanego `dane.skaner` — ta sama umowa po
// drugiej stronie granicy pakietu.
type skaner interface {
	Scan(cele ...any) error
}

func odczytajKonsultacje(w skaner) (Konsultacja, error) {
	var (
		wpis        Konsultacja
		okno        sql.NullString
		rada        sql.NullString
		skrot       sql.NullString
		prowenancja sql.NullString
		powod       sql.NullString
	)
	if err := w.Scan(&wpis.ID, &wpis.Identyfikator, &okno, &wpis.Pytajacy,
		&wpis.DoradcaKanal, &wpis.DoradcaModel, &wpis.Pytanie, &rada, &skrot,
		&prowenancja, &wpis.Stan, &powod, &wpis.Utworzono); err != nil {
		return Konsultacja{}, err
	}
	wpis.OknoId = okno.String
	wpis.Rada = rada.String
	wpis.Skrot = skrot.String
	wpis.Prowenancja = prowenancja.String
	wpis.Powod = powod.String
	return wpis, nil
}

// sprawdzKonsultacje odrzuca wpis niespójny wcześniej niż CHECK bazy —
// wołający dostaje zdanie po polsku zamiast komunikatu sterownika.
func sprawdzKonsultacje(wpis Konsultacja) error {
	switch {
	case strings.TrimSpace(wpis.Identyfikator) == "":
		return errors.New("podagenci: konsultacja bez identyfikatora zewnętrznego")
	case strings.TrimSpace(wpis.Pytanie) == "":
		return errors.New("podagenci: konsultacja bez treści pytania")
	case strings.TrimSpace(wpis.DoradcaKanal) == "":
		return errors.New("podagenci: konsultacja bez wskazania doradcy")
	case wpis.Stan == StanRady && strings.TrimSpace(wpis.Rada) == "":
		return errors.New("podagenci: konsultacja w stanie rady bez treści rady")
	case wpis.Stan == StanOdmowyRady && strings.TrimSpace(wpis.Powod) == "":
		return errors.New("podagenci: odmowa konsultacji bez powodu")
	case wpis.Stan != StanRady && wpis.Stan != StanOdmowyRady:
		return fmt.Errorf("podagenci: nieznany stan konsultacji %q", wpis.Stan)
	}
	return nil
}

func kolumnaNapisu(wartosc string) any {
	if strings.TrimSpace(wartosc) == "" {
		return nil
	}
	return wartosc
}
