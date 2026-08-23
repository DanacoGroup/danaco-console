// Odpowiedzialność pliku: postać dokumentu modułu Studio — drzewo postaci
// (tabela `postac_dokumentu_studio`), arkusz stylów nazwanych
// (`styl_nazwany_studio`) i sekcje o własnych nastawach
// (`sekcja_dokumentu_studio`).
//
// Ten plik deklaruje kontrakt obszaru postaci — drzewa postaci, arkusza stylów,
// sekcji oraz obiektów, aparatu i pól, które leżą w pliku sąsiednim
// (`studio_postac_obiekty.go`). Jeden kontrakt w jednym miejscu, wzorem
// `studio.go`.
//
// Interfejs `RepozytoriumPostaciStudia` wchodzi do `RepozytoriumStudia` przez
// zagnieżdżenie — Studio ma jedno repozytorium, nie dwa, więc adapter modułu
// dostaje postać tą samą zależnością, którą dostaje dokument.
//
// ── Dlaczego drzewo postaci idzie jednym zapisem, a style i sekcje nie ───────
// Powód stoi w migracji 361 i nie powtarzam go tu w całości: po drzewie się nie
// pyta, drzewo się czyta i zapisuje całe; po stylu i po sekcji się PYTA („ile
// miejsc używa tego stylu", „która sekcja obejmuje ten znak"), więc mają wiersze.
//
// ── Przedrostek nazw pomocniczych ────────────────────────────────────────────
// Wszystkie nazwy pomocnicze tego odcinka niosą przedrostek `postac` —
// przestrzeń nazw pakietu `dane` jest dzielona z innymi wykonawcami.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// PostacDokumentuStudia to wiersz tabeli `postac_dokumentu_studio`.
//
// `PostacJSON` niesie drzewo postaci w kształcie kontraktowego
// `StudioDocumentForm` — bloki, fragmenty o jednolitej postaci znaku, tabele,
// listy. Warstwa danych go NIE rozbiera: przekład na kontrakt należy do rdzenia,
// a baza jest tu magazynem, nie rachunkiem.
type PostacDokumentuStudia struct {
	ID                int64
	DokumentID        int64
	PostacJSON        string
	NastawyStronyJSON *string
	WersjaPostaci     int64
	Utworzono         string
	Zaktualizowano    string
}

// StylNazwanyStudia to wiersz tabeli `styl_nazwany_studio`. Nazwa jest jedynym
// identyfikatorem stylu w obrębie dokumentu — tak samo jak w pakiecie biurowym.
type StylNazwanyStudia struct {
	ID                int64
	DokumentID        int64
	Nazwa             string
	NazwaWidoczna     *string
	Rodzaj            string
	StylNadrzedny     *string
	StylNastepny      *string
	PostacZnakuJSON   *string
	PostacAkapituJSON *string
	Fabryczny         bool
	Utworzono         string
	Zaktualizowano    string
}

// SekcjaDokumentuStudia to wiersz tabeli `sekcja_dokumentu_studio`.
type SekcjaDokumentuStudia struct {
	ID                int64
	Kod               string
	DokumentID        int64
	Kolejnosc         int64
	Tytul             *string
	ZakresOd          int64
	ZakresDo          int64
	Rozpoczecie       string
	NastawyStronyJSON *string
	NaglowkiJSON      *string
	NumeracjaJSON     *string
	ZnakWodnyJSON     *string
	Utworzono         string
	Zaktualizowano    string
}

// RepozytoriumPostaciStudia jest kontraktem obszaru postaci dokumentu.
type RepozytoriumPostaciStudia interface {
	// --- drzewo postaci (ten plik) ---
	ZapiszPostacDokumentu(ctx context.Context, postac PostacDokumentuStudia) (PostacDokumentuStudia, error)
	PostacDokumentu(ctx context.Context, dokumentID int64) (PostacDokumentuStudia, error)

	// --- arkusz stylów nazwanych (ten plik) ---
	ZapiszStylNazwany(ctx context.Context, styl StylNazwanyStudia) (StylNazwanyStudia, error)
	StylNazwany(ctx context.Context, dokumentID int64, nazwa string) (StylNazwanyStudia, error)
	StyleNazwane(ctx context.Context, dokumentID int64, rodzaj string) ([]StylNazwanyStudia, error)
	StyleDziedziczace(ctx context.Context, dokumentID int64, nazwaNadrzednego string) ([]StylNazwanyStudia, error)
	UsunStylNazwany(ctx context.Context, dokumentID int64, nazwa string) (bool, error)

	// --- sekcje (ten plik) ---
	ZapiszSekcje(ctx context.Context, sekcja SekcjaDokumentuStudia) (SekcjaDokumentuStudia, error)
	Sekcja(ctx context.Context, kod string) (SekcjaDokumentuStudia, error)
	Sekcje(ctx context.Context, dokumentID int64) ([]SekcjaDokumentuStudia, error)
	UsunSekcje(ctx context.Context, kod string) (bool, error)

	// --- obiekty, aparat i pola (`studio_postac_obiekty.go`) ---
	ZapiszObiektDokumentu(ctx context.Context, obiekt ObiektDokumentuStudia) (ObiektDokumentuStudia, error)
	ObiektDokumentu(ctx context.Context, kod string) (ObiektDokumentuStudia, error)
	ObiektyDokumentu(ctx context.Context, dokumentID int64, rodzaj string) ([]ObiektDokumentuStudia, error)
	UsunObiektDokumentu(ctx context.Context, kod string) (bool, error)
	ZapiszElementAparatu(ctx context.Context, element ElementAparatuStudia) (ElementAparatuStudia, error)
	ElementAparatu(ctx context.Context, kod string) (ElementAparatuStudia, error)
	ElementyAparatu(ctx context.Context, dokumentID int64, rodzaj string) ([]ElementAparatuStudia, error)
	UsunElementAparatu(ctx context.Context, kod string) (bool, error)
	ZapiszPoleDokumentu(ctx context.Context, pole PoleDokumentuStudia) (PoleDokumentuStudia, error)
	PolaDokumentu(ctx context.Context, dokumentID int64, rodzaj string) ([]PoleDokumentuStudia, error)
	UsunPoleDokumentu(ctx context.Context, kod string) (bool, error)

	// ── Czego w tym interfejsie NIE MA i dlaczego ────────────────────────────
	// Blokady fragmentów, dziennik czynności, znakowanie, kopie zapasowe,
	// nastawy pracy, zajęcia fragmentów i spięcia wykonawców stoją nad tymi
	// samymi tabelami (migracje 363-370), ale ich metody napisał inny wykonawca
	// w plikach `studio_kontrola_pracy.go` i `studio_znakowanie_wykonawcy.go`.
	// Nie dopisuję ich tutaj i nie zakładam drugich: dwa zestawy metod nad jedną
	// tabelą to dwie prawdy o tym samym wierszu. Kontrakt tamtego obszaru
	// należy do tamtych plików — wymaga ogłoszenia interfejsem tak samo jak ten
	// i jest to wypisane w sprawozdaniu jako rzecz do domknięcia.
}

const (
	postacKolumnyPostaci = `id, dokument_id, postac_json, nastawy_strony_json,
	                        wersja_postaci, utworzono, zaktualizowano`

	// Zapis zakłada postać albo nadpisuje zastaną. Numer porządkowy postaci
	// przy nadpisaniu ROŚNIE po stronie bazy, nie wywołującego: gdyby liczył go
	// rdzeń, dwa zapisy z tym samym numerem byłyby możliwe i okno nie miałoby po
	// czym poznać, że trzyma stan przestarzały.
	postacZapiszPostac = `INSERT INTO postac_dokumentu_studio
	                      (dokument_id, postac_json, nastawy_strony_json)
	                      VALUES (?, ?, ?)
	                      ON CONFLICT(dokument_id) DO UPDATE SET
	                          postac_json = excluded.postac_json,
	                          nastawy_strony_json = excluded.nastawy_strony_json,
	                          wersja_postaci = postac_dokumentu_studio.wersja_postaci + 1,
	                          zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	postacPobierzPostac = `SELECT ` + postacKolumnyPostaci + ` FROM postac_dokumentu_studio
	                       WHERE dokument_id = ?`

	postacKolumnyStylu = `id, dokument_id, nazwa, nazwa_widoczna, rodzaj, styl_nadrzedny,
	                      styl_nastepny, postac_znaku_json, postac_akapitu_json, fabryczny,
	                      utworzono, zaktualizowano`

	// Zmiana stylu nazwanego jest nadpisaniem wiersza, nie założeniem drugiego:
	// nazwa jest tożsamością stylu, a dwa wiersze o tej samej nazwie znaczyłyby
	// dwie prawdy o tym, jak wygląda „Nagłówek 1”.
	postacZapiszStyl = `INSERT INTO styl_nazwany_studio
	                    (dokument_id, nazwa, nazwa_widoczna, rodzaj, styl_nadrzedny,
	                     styl_nastepny, postac_znaku_json, postac_akapitu_json, fabryczny)
	                    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	                    ON CONFLICT(dokument_id, nazwa) DO UPDATE SET
	                        nazwa_widoczna = excluded.nazwa_widoczna,
	                        rodzaj = excluded.rodzaj,
	                        styl_nadrzedny = excluded.styl_nadrzedny,
	                        styl_nastepny = excluded.styl_nastepny,
	                        postac_znaku_json = excluded.postac_znaku_json,
	                        postac_akapitu_json = excluded.postac_akapitu_json,
	                        zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	postacPobierzStyl = `SELECT ` + postacKolumnyStylu + ` FROM styl_nazwany_studio
	                     WHERE dokument_id = ? AND nazwa = ?`

	postacListaStylow = `SELECT ` + postacKolumnyStylu + ` FROM styl_nazwany_studio
	                     WHERE dokument_id = ? ORDER BY rodzaj, nazwa`

	postacListaStylowRodzaju = `SELECT ` + postacKolumnyStylu + ` FROM styl_nazwany_studio
	                            WHERE dokument_id = ? AND rodzaj = ? ORDER BY nazwa`

	postacListaStylowDziedziczacych = `SELECT ` + postacKolumnyStylu + ` FROM styl_nazwany_studio
	                                   WHERE dokument_id = ? AND styl_nadrzedny = ? ORDER BY nazwa`

	// Stylu fabrycznego nie usuwa warstwa danych — o odmowie rozstrzyga rdzeń,
	// bo to on umie nazwać powód słowami Operatora. Warunek na `fabryczny` stoi
	// tu jako druga zapora: gdyby rdzeń kiedyś zapomniał sprawdzić, styl
	// fabryczny nie zniknie po cichu.
	postacUsunStyl = `DELETE FROM styl_nazwany_studio
	                  WHERE dokument_id = ? AND nazwa = ? AND fabryczny = 0`

	postacKolumnySekcji = `id, identyfikator_zewnetrzny, dokument_id, kolejnosc, tytul,
	                       zakres_od, zakres_do, rozpoczecie, nastawy_strony_json,
	                       naglowki_json, numeracja_json, znak_wodny_json,
	                       utworzono, zaktualizowano`

	postacZapiszSekcje = `INSERT INTO sekcja_dokumentu_studio
	                      (identyfikator_zewnetrzny, dokument_id, kolejnosc, tytul,
	                       zakres_od, zakres_do, rozpoczecie, nastawy_strony_json,
	                       naglowki_json, numeracja_json, znak_wodny_json)
	                      VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	                      ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                          kolejnosc = excluded.kolejnosc,
	                          tytul = excluded.tytul,
	                          zakres_od = excluded.zakres_od,
	                          zakres_do = excluded.zakres_do,
	                          rozpoczecie = excluded.rozpoczecie,
	                          nastawy_strony_json = excluded.nastawy_strony_json,
	                          naglowki_json = excluded.naglowki_json,
	                          numeracja_json = excluded.numeracja_json,
	                          znak_wodny_json = excluded.znak_wodny_json,
	                          zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	postacPobierzSekcje = `SELECT ` + postacKolumnySekcji + ` FROM sekcja_dokumentu_studio
	                       WHERE identyfikator_zewnetrzny = ?`

	postacListaSekcji = `SELECT ` + postacKolumnySekcji + ` FROM sekcja_dokumentu_studio
	                     WHERE dokument_id = ? ORDER BY kolejnosc, id`

	postacUsunSekcje = `DELETE FROM sekcja_dokumentu_studio WHERE identyfikator_zewnetrzny = ?`
)

// ZapiszPostacDokumentu zakłada postać dokumentu albo nadpisuje zastaną
// i oddaje stan po zapisie wraz z podniesionym numerem porządkowym.
func (r *repozytoriumStudia) ZapiszPostacDokumentu(ctx context.Context,
	postac PostacDokumentuStudia) (PostacDokumentuStudia, error) {

	if postac.DokumentID == 0 {
		return PostacDokumentuStudia{}, fmt.Errorf("dane: postać dokumentu studio bez dokumentu")
	}
	if postac.PostacJSON == "" {
		postac.PostacJSON = "{}"
	}
	polecenie, err := r.zapytania.przygotuj(ctx, postacZapiszPostac)
	if err != nil {
		return PostacDokumentuStudia{}, err
	}
	_, err = polecenie.ExecContext(ctx, postac.DokumentID, postac.PostacJSON,
		tekstDoKolumny(postac.NastawyStronyJSON))
	if err != nil {
		return PostacDokumentuStudia{}, fmt.Errorf(
			"dane: nie można zapisać postaci dokumentu studio %d: %w", postac.DokumentID, err)
	}
	return r.PostacDokumentu(ctx, postac.DokumentID)
}

// PostacDokumentu zwraca postać dokumentu. Brak wiersza wraca jako
// ErrBrakWiersza — dokument bez zapisanej postaci jest normalnym stanem
// (dokumenty sprzed dobudowy postaci go mają), a rdzeń podstawia wtedy postać
// domyślną. Zamiana braku na pustą postać tutaj odebrałaby rdzeniowi możliwość
// odróżnienia „nie ma jeszcze” od „jest i jest puste”.
func (r *repozytoriumStudia) PostacDokumentu(ctx context.Context,
	dokumentID int64) (PostacDokumentuStudia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, postacPobierzPostac)
	if err != nil {
		return PostacDokumentuStudia{}, err
	}
	postac, err := postacOdczytajPostac(polecenie.QueryRowContext(ctx, dokumentID))
	if errors.Is(err, sql.ErrNoRows) {
		return PostacDokumentuStudia{}, ErrBrakWiersza
	}
	if err != nil {
		return PostacDokumentuStudia{}, fmt.Errorf(
			"dane: nieczytelny wiersz postaci dokumentu studio %d: %w", dokumentID, err)
	}
	return postac, nil
}

// ZapiszStylNazwany zakłada styl nazwany albo nadpisuje zastany.
func (r *repozytoriumStudia) ZapiszStylNazwany(ctx context.Context,
	styl StylNazwanyStudia) (StylNazwanyStudia, error) {

	if styl.DokumentID == 0 || styl.Nazwa == "" {
		return StylNazwanyStudia{}, fmt.Errorf("dane: styl nazwany studio bez dokumentu albo bez nazwy")
	}
	if styl.Rodzaj == "" {
		styl.Rodzaj = "paragraph"
	}
	polecenie, err := r.zapytania.przygotuj(ctx, postacZapiszStyl)
	if err != nil {
		return StylNazwanyStudia{}, err
	}
	_, err = polecenie.ExecContext(ctx, styl.DokumentID, styl.Nazwa,
		tekstDoKolumny(styl.NazwaWidoczna), styl.Rodzaj,
		tekstDoKolumny(styl.StylNadrzedny), tekstDoKolumny(styl.StylNastepny),
		tekstDoKolumny(styl.PostacZnakuJSON), tekstDoKolumny(styl.PostacAkapituJSON),
		liczbaLogiczna(styl.Fabryczny))
	if err != nil {
		return StylNazwanyStudia{}, fmt.Errorf("dane: nie można zapisać stylu %q: %w", styl.Nazwa, err)
	}
	return r.StylNazwany(ctx, styl.DokumentID, styl.Nazwa)
}

// StylNazwany zwraca styl o wskazanej nazwie.
func (r *repozytoriumStudia) StylNazwany(ctx context.Context, dokumentID int64,
	nazwa string) (StylNazwanyStudia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, postacPobierzStyl)
	if err != nil {
		return StylNazwanyStudia{}, err
	}
	styl, err := postacOdczytajStyl(polecenie.QueryRowContext(ctx, dokumentID, nazwa))
	if errors.Is(err, sql.ErrNoRows) {
		return StylNazwanyStudia{}, ErrBrakWiersza
	}
	if err != nil {
		return StylNazwanyStudia{}, fmt.Errorf("dane: nieczytelny wiersz stylu %q: %w", nazwa, err)
	}
	return styl, nil
}

// StyleNazwane zwraca arkusz stylów dokumentu; puste `rodzaj` znaczy wszystkie.
func (r *repozytoriumStudia) StyleNazwane(ctx context.Context, dokumentID int64,
	rodzaj string) ([]StylNazwanyStudia, error) {

	zapytanie, argumenty := postacListaStylow, []any{dokumentID}
	if rodzaj != "" {
		zapytanie, argumenty = postacListaStylowRodzaju, []any{dokumentID, rodzaj}
	}
	return r.postacWykazStylow(ctx, zapytanie, argumenty)
}

// StyleDziedziczace zwraca style, które dziedziczą po wskazanym stylu.
//
// Służy jednej rzeczy, po której poznaje się styl nazwany: zmiana stylu
// nadrzędnego ma przestawić WSZYSTKIE miejsca, które go używają — także te,
// które używają go pośrednio, przez styl potomny. Bez tego zapytania rdzeń
// musiałby czytać cały arkusz i składać drzewo dziedziczenia przy każdej zmianie.
func (r *repozytoriumStudia) StyleDziedziczace(ctx context.Context, dokumentID int64,
	nazwaNadrzednego string) ([]StylNazwanyStudia, error) {

	return r.postacWykazStylow(ctx, postacListaStylowDziedziczacych,
		[]any{dokumentID, nazwaNadrzednego})
}

// UsunStylNazwany usuwa styl własny. Styl fabryczny zostaje i metoda oddaje
// fałsz — powód odmowy nazywa rdzeń.
func (r *repozytoriumStudia) UsunStylNazwany(ctx context.Context, dokumentID int64,
	nazwa string) (bool, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, postacUsunStyl)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, dokumentID, nazwa)
	if err != nil {
		return false, fmt.Errorf("dane: nie można usunąć stylu %q: %w", nazwa, err)
	}
	usuniete, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nieznany skutek usunięcia stylu %q: %w", nazwa, err)
	}
	return usuniete > 0, nil
}

// ZapiszSekcje zakłada sekcję albo nadpisuje zastaną.
func (r *repozytoriumStudia) ZapiszSekcje(ctx context.Context,
	sekcja SekcjaDokumentuStudia) (SekcjaDokumentuStudia, error) {

	if sekcja.Kod == "" || sekcja.DokumentID == 0 {
		return SekcjaDokumentuStudia{}, fmt.Errorf("dane: sekcja studio bez identyfikatora albo bez dokumentu")
	}
	if sekcja.Rozpoczecie == "" {
		sekcja.Rozpoczecie = "continuous"
	}
	polecenie, err := r.zapytania.przygotuj(ctx, postacZapiszSekcje)
	if err != nil {
		return SekcjaDokumentuStudia{}, err
	}
	_, err = polecenie.ExecContext(ctx, sekcja.Kod, sekcja.DokumentID, sekcja.Kolejnosc,
		tekstDoKolumny(sekcja.Tytul), sekcja.ZakresOd, sekcja.ZakresDo, sekcja.Rozpoczecie,
		tekstDoKolumny(sekcja.NastawyStronyJSON), tekstDoKolumny(sekcja.NaglowkiJSON),
		tekstDoKolumny(sekcja.NumeracjaJSON), tekstDoKolumny(sekcja.ZnakWodnyJSON))
	if err != nil {
		return SekcjaDokumentuStudia{}, fmt.Errorf("dane: nie można zapisać sekcji %q: %w", sekcja.Kod, err)
	}
	return r.Sekcja(ctx, sekcja.Kod)
}

// Sekcja zwraca sekcję o wskazanym kodzie.
func (r *repozytoriumStudia) Sekcja(ctx context.Context, kod string) (SekcjaDokumentuStudia, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, postacPobierzSekcje)
	if err != nil {
		return SekcjaDokumentuStudia{}, err
	}
	sekcja, err := postacOdczytajSekcje(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return SekcjaDokumentuStudia{}, ErrBrakWiersza
	}
	if err != nil {
		return SekcjaDokumentuStudia{}, fmt.Errorf("dane: nieczytelny wiersz sekcji %q: %w", kod, err)
	}
	return sekcja, nil
}

// Sekcje zwraca sekcje dokumentu w kolejności czytania.
func (r *repozytoriumStudia) Sekcje(ctx context.Context, dokumentID int64) ([]SekcjaDokumentuStudia, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, postacListaSekcji)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, dokumentID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać sekcji dokumentu %d: %w", dokumentID, err)
	}
	defer wiersze.Close()

	lista := []SekcjaDokumentuStudia{}
	for wiersze.Next() {
		sekcja, err := postacOdczytajSekcje(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz sekcji: %w", err)
		}
		lista = append(lista, sekcja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt sekcji dokumentu %d: %w", dokumentID, err)
	}
	return lista, nil
}

// UsunSekcje usuwa sekcję. Treści sekcji nie rusza — przeniesienie jej do sekcji
// poprzedniej należy do rdzenia, bo to on wie, co znaczy „poprzednia”.
func (r *repozytoriumStudia) UsunSekcje(ctx context.Context, kod string) (bool, error) {
	return r.postacUsunWiersz(ctx, postacUsunSekcje, kod, "sekcji")
}

// postacWykazStylow wykonuje zapytanie wykazu stylów i składa wynik.
func (r *repozytoriumStudia) postacWykazStylow(ctx context.Context, zapytanie string,
	argumenty []any) ([]StylNazwanyStudia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, argumenty...)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać stylów nazwanych: %w", err)
	}
	defer wiersze.Close()

	lista := []StylNazwanyStudia{}
	for wiersze.Next() {
		styl, err := postacOdczytajStyl(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz stylu nazwanego: %w", err)
		}
		lista = append(lista, styl)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt stylów nazwanych: %w", err)
	}
	return lista, nil
}

// postacUsunWiersz wykonuje usunięcie po identyfikatorze zewnętrznym i oddaje,
// czy wiersz istniał. Jedno miejsce dla wszystkich bytów postaci — usunięcie
// obiektu, elementu aparatu i pola różni się wyłącznie zapytaniem.
func (r *repozytoriumStudia) postacUsunWiersz(ctx context.Context, zapytanie, kod,
	nazwaBytu string) (bool, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, kod)
	if err != nil {
		return false, fmt.Errorf("dane: nie można usunąć %s %q: %w", nazwaBytu, kod, err)
	}
	usuniete, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nieznany skutek usunięcia %s %q: %w", nazwaBytu, kod, err)
	}
	return usuniete > 0, nil
}

// postacOdczytajPostac składa postać dokumentu z jednego wiersza wyniku.
func postacOdczytajPostac(wiersz skaner) (PostacDokumentuStudia, error) {
	var postac PostacDokumentuStudia
	var nastawy sql.NullString
	err := wiersz.Scan(&postac.ID, &postac.DokumentID, &postac.PostacJSON, &nastawy,
		&postac.WersjaPostaci, &postac.Utworzono, &postac.Zaktualizowano)
	if err != nil {
		return PostacDokumentuStudia{}, err
	}
	postac.NastawyStronyJSON = tekstZKolumny(nastawy)
	return postac, nil
}

// postacOdczytajStyl składa styl nazwany z jednego wiersza wyniku.
func postacOdczytajStyl(wiersz skaner) (StylNazwanyStudia, error) {
	var styl StylNazwanyStudia
	var widoczna, nadrzedny, nastepny, znak, akapit sql.NullString
	var fabryczny int64
	err := wiersz.Scan(&styl.ID, &styl.DokumentID, &styl.Nazwa, &widoczna, &styl.Rodzaj,
		&nadrzedny, &nastepny, &znak, &akapit, &fabryczny,
		&styl.Utworzono, &styl.Zaktualizowano)
	if err != nil {
		return StylNazwanyStudia{}, err
	}
	styl.NazwaWidoczna = tekstZKolumny(widoczna)
	styl.StylNadrzedny = tekstZKolumny(nadrzedny)
	styl.StylNastepny = tekstZKolumny(nastepny)
	styl.PostacZnakuJSON = tekstZKolumny(znak)
	styl.PostacAkapituJSON = tekstZKolumny(akapit)
	styl.Fabryczny = fabryczny != 0
	return styl, nil
}

// postacOdczytajSekcje składa sekcję z jednego wiersza wyniku.
func postacOdczytajSekcje(wiersz skaner) (SekcjaDokumentuStudia, error) {
	var sekcja SekcjaDokumentuStudia
	var tytul, nastawy, naglowki, numeracja, znakWodny sql.NullString
	err := wiersz.Scan(&sekcja.ID, &sekcja.Kod, &sekcja.DokumentID, &sekcja.Kolejnosc,
		&tytul, &sekcja.ZakresOd, &sekcja.ZakresDo, &sekcja.Rozpoczecie,
		&nastawy, &naglowki, &numeracja, &znakWodny,
		&sekcja.Utworzono, &sekcja.Zaktualizowano)
	if err != nil {
		return SekcjaDokumentuStudia{}, err
	}
	sekcja.Tytul = tekstZKolumny(tytul)
	sekcja.NastawyStronyJSON = tekstZKolumny(nastawy)
	sekcja.NaglowkiJSON = tekstZKolumny(naglowki)
	sekcja.NumeracjaJSON = tekstZKolumny(numeracja)
	sekcja.ZnakWodnyJSON = tekstZKolumny(znakWodny)
	return sekcja, nil
}
