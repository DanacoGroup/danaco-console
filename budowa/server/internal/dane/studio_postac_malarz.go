// Odpowiedzialność pliku: warstwa danych MALARZA FORMATÓW modułu Studio —
// tabela `postac_malarza_studio` z migracji 368.
//
// ── Dlaczego malarz ma wiersz, a nie pamięć procesu ─────────────────────────
// Migracja 368 zapisała powód wprost i ten plik go wykonuje: malarz kopiuje
// POSTAĆ, nie treść, i nanosi ją w innym miejscu, więc między pobraniem
// a naniesieniem stoją DWIE osobne komendy (`studio.format.painter.copy`
// i `.apply`). Postać trzymana w pamięci procesu przepadała przy przeładowaniu
// rdzenia — Operator pobierał postać, rdzeń wstawał od nowa, a naniesienie
// odmawiało „takiej postaci nie znam". Przy pracy modelu przepadała jeszcze
// łatwiej: model pobiera postać jednym narzędziem i nanosi drugim, być może po
// kilku innych czynnościach.
//
// ── Dlaczego wiersz wygasa ──────────────────────────────────────────────────
// Malarz jest narzędziem JEDNEJ czynności. Postać pobrana wczoraj i naniesiona
// dziś byłaby zaskoczeniem, nie pomocą — stąd kolumna `wygasa` i odczyt, który
// wpisu wygasłego nie oddaje. Wygasły wiersz nie jest przy tym kasowany
// w odczycie: sprzątanie idzie osobnym wywołaniem, bo odczyt, który po cichu
// usuwa wiersze, jest odczytem zmieniającym stan.
//
// ── Dlaczego kluczem jest OKNO, a nie dokument ──────────────────────────────
// Malarz przenosi postać MIĘDZY dokumentami — to jest jego zwykłe użycie
// w pakiecie biurowym. Kluczem jest więc okno, w którym Operator pracuje;
// dokument, z którego postać zabrano, stoi obok jako wiedza, a nie jako warunek.
//
// ── Przedrostek nazw pomocniczych ───────────────────────────────────────────
// Nazwy pomocnicze tego pliku niosą przedrostek `malarz` — przestrzeń nazw
// pakietu `dane` jest dzielona z innymi wykonawcami.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// PostacMalarzaStudia to wiersz tabeli `postac_malarza_studio` — postać zabrana
// malarzem formatów.
type PostacMalarzaStudia struct {
	ID          int64
	Kod         string
	Okno        string
	DokumentKod *string
	// PostacZnakuJSON i PostacAkapituJSON niosą postać kontraktu w zapisie JSON.
	// Kolumn na pojedyncze cechy nie ma z zamysłu: postać znaku ma siedemnaście
	// cech, a postać akapitu dwadzieścia, i rosną razem z kontraktem. Kolumna na
	// każdą z nich znaczyłaby migrację przy każdym dołożonym polu.
	PostacZnakuJSON   *string
	PostacAkapituJSON *string
	StylNazwany       *string
	Utworzono         string
	// Wygasa jest chwilą, po której wpisu nie wolno nanieść. Puste znaczy „nie
	// wygasa" — droga zostawiona świadomie, bo postać zabraną na potrzeby
	// szablonu bywa potrzebna dłużej niż jedną czynność.
	Wygasa *string
}

const (
	malarzKolumny = `id, identyfikator_zewnetrzny, okno, dokument_kod, postac_znaku_json,
	                 postac_akapitu_json, styl_nazwany, utworzono, wygasa`

	malarzZapisz = `INSERT INTO postac_malarza_studio
	                (identyfikator_zewnetrzny, okno, dokument_kod, postac_znaku_json,
	                 postac_akapitu_json, styl_nazwany, wygasa)
	                VALUES (?, ?, ?, ?, ?, ?, ?)`

	// Odczyt po uchwycie NIE pyta o okno: uchwyt jest identyfikatorem nadanym
	// przez rdzeń i dowodzi sam z siebie. Warunek okna kazałby naniesieniu
	// podawać okno drugi raz, a model, który pobrał postać jednym narzędziem,
	// niesie do drugiego uchwyt, nie okno.
	malarzPobierz = `SELECT ` + malarzKolumny + ` FROM postac_malarza_studio
	                 WHERE identyfikator_zewnetrzny = ?`

	// Najświeższa postać okna — droga naniesienia bez podanego uchwytu
	// („nanieś to, co ostatnio zabrałem"). Wygasłych nie oddaje.
	malarzNajswiezsza = `SELECT ` + malarzKolumny + ` FROM postac_malarza_studio
	                     WHERE okno = ?
	                       AND (wygasa IS NULL OR wygasa > strftime('%Y-%m-%dT%H:%M:%fZ','now'))
	                     ORDER BY utworzono DESC, id DESC LIMIT 1`

	// Sprzątanie zdejmuje wpisy wygasłe. Idzie osobnym wywołaniem, nie w odczycie:
	// odczyt kasujący wiersze byłby odczytem zmieniającym stan, a przy dwóch
	// oknach pytających naraz jedno zabierałoby postać drugiemu.
	malarzSprzataj = `DELETE FROM postac_malarza_studio
	                  WHERE wygasa IS NOT NULL
	                    AND wygasa <= strftime('%Y-%m-%dT%H:%M:%fZ','now')`
)

// MalarzFormatowStudia jest kontraktem tej warstwy.
type MalarzFormatowStudia interface {
	ZapiszPostacMalarza(ctx context.Context,
		postac PostacMalarzaStudia) (PostacMalarzaStudia, error)
	PostacMalarza(ctx context.Context, kod string) (PostacMalarzaStudia, error)
	NajswiezszaPostacMalarza(ctx context.Context, okno string) (PostacMalarzaStudia, error)
	SprzatnijPostacieMalarza(ctx context.Context) (int, error)
}

// ZapiszPostacMalarza odkłada postać zabraną malarzem i oddaje ją zapisaną.
func (r *repozytoriumStudia) ZapiszPostacMalarza(ctx context.Context,
	postac PostacMalarzaStudia) (PostacMalarzaStudia, error) {

	kod := strings.TrimSpace(postac.Kod)
	if kod == "" {
		return PostacMalarzaStudia{}, fmt.Errorf("dane: postać malarza bez identyfikatora")
	}
	// Postać bez ani jednej z trzech treści nie jest postacią — naniesienie
	// takiego wpisu nie zmieniłoby niczego i oddałoby „naniesiono".
	if postac.PostacZnakuJSON == nil && postac.PostacAkapituJSON == nil &&
		postac.StylNazwany == nil {

		return PostacMalarzaStudia{}, fmt.Errorf(
			"dane: postać malarza %q nie niesie ani postaci znaku, ani postaci akapitu, "+
				"ani stylu nazwanego", kod)
	}
	polecenie, err := r.zapytania.przygotuj(ctx, malarzZapisz)
	if err != nil {
		return PostacMalarzaStudia{}, err
	}
	if _, err := polecenie.ExecContext(ctx, kod, strings.TrimSpace(postac.Okno),
		tekstDoKolumny(postac.DokumentKod), tekstDoKolumny(postac.PostacZnakuJSON),
		tekstDoKolumny(postac.PostacAkapituJSON), tekstDoKolumny(postac.StylNazwany),
		tekstDoKolumny(postac.Wygasa)); err != nil {

		return PostacMalarzaStudia{}, fmt.Errorf(
			"dane: nie można zapisać postaci malarza %q: %w", kod, err)
	}
	return r.PostacMalarza(ctx, kod)
}

// PostacMalarza oddaje postać o wskazanym uchwycie.
//
// Wpis wygasły jest tu BRAKIEM WIERSZA, nie wierszem z datą w przeszłości:
// rdzeń ma powiedzieć „ta postać już nie obowiązuje", a nie nanieść postać,
// której Operator nie pamięta.
func (r *repozytoriumStudia) PostacMalarza(ctx context.Context,
	kod string) (PostacMalarzaStudia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, malarzPobierz)
	if err != nil {
		return PostacMalarzaStudia{}, err
	}
	postac, err := malarzOdczytaj(polecenie.QueryRowContext(ctx, strings.TrimSpace(kod)))
	if errors.Is(err, sql.ErrNoRows) {
		return PostacMalarzaStudia{}, ErrBrakWiersza
	}
	if err != nil {
		return PostacMalarzaStudia{}, fmt.Errorf(
			"dane: nieczytelny wiersz postaci malarza %q: %w", kod, err)
	}
	return postac, nil
}

// NajswiezszaPostacMalarza oddaje ostatnią niewygasłą postać zabraną w oknie.
func (r *repozytoriumStudia) NajswiezszaPostacMalarza(ctx context.Context,
	okno string) (PostacMalarzaStudia, error) {

	if strings.TrimSpace(okno) == "" {
		return PostacMalarzaStudia{}, ErrBrakWiersza
	}
	polecenie, err := r.zapytania.przygotuj(ctx, malarzNajswiezsza)
	if err != nil {
		return PostacMalarzaStudia{}, err
	}
	postac, err := malarzOdczytaj(polecenie.QueryRowContext(ctx, strings.TrimSpace(okno)))
	if errors.Is(err, sql.ErrNoRows) {
		return PostacMalarzaStudia{}, ErrBrakWiersza
	}
	if err != nil {
		return PostacMalarzaStudia{}, fmt.Errorf(
			"dane: nieczytelny wiersz najświeższej postaci malarza okna %q: %w", okno, err)
	}
	return postac, nil
}

// SprzatnijPostacieMalarza zdejmuje wpisy wygasłe i oddaje, ile ich było.
func (r *repozytoriumStudia) SprzatnijPostacieMalarza(ctx context.Context) (int, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, malarzSprzataj)
	if err != nil {
		return 0, err
	}
	wynik, err := polecenie.ExecContext(ctx)
	if err != nil {
		return 0, fmt.Errorf("dane: nie można zdjąć wygasłych postaci malarza: %w", err)
	}
	zdjete, err := wynik.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("dane: nieznany skutek zdjęcia wygasłych postaci malarza: %w", err)
	}
	return int(zdjete), nil
}

// malarzOdczytaj składa wiersz postaci malarza.
func malarzOdczytaj(wiersz interface{ Scan(...any) error }) (PostacMalarzaStudia, error) {
	var postac PostacMalarzaStudia
	var dokument, znak, akapit, styl, wygasa sql.NullString
	if err := wiersz.Scan(&postac.ID, &postac.Kod, &postac.Okno, &dokument, &znak,
		&akapit, &styl, &postac.Utworzono, &wygasa); err != nil {

		return PostacMalarzaStudia{}, err
	}
	postac.DokumentKod = tekstZKolumny(dokument)
	postac.PostacZnakuJSON = tekstZKolumny(znak)
	postac.PostacAkapituJSON = tekstZKolumny(akapit)
	postac.StylNazwany = tekstZKolumny(styl)
	postac.Wygasa = tekstZKolumny(wygasa)
	return postac, nil
}
