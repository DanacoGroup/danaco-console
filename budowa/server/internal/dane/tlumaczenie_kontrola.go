// Faseta kontroli modułu Translate: wykaz terminów zawężony, profile
// kontroli jakości, obieg zatwierdzeń panelu i ustalenia korekty językowej —
// trzy byty odpowiadające na pytanie, czy przekład wolno wypuścić.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// FiltrTerminow zawęża wykaz terminów słownika po języku, dziedzinie,
// stanie i frazie; pola puste nie zawężają wyniku.
type FiltrTerminow struct {
	Jezyk     string
	Dziedzina string
	Stan      string
	Fraza     string
	Limit     int
	Offset    int
}

// TerminyZawezone oddaje wykaz terminów słownika dopasowanych filtrem wraz
// z liczbą wszystkich pasujących wierszy.
func (r *repozytoriumTlumaczen) TerminyZawezone(ctx context.Context,
	filtr FiltrTerminow) ([]TerminSlownika, int, error) {

	warunki := []string{}
	argumenty := []any{}
	if strings.TrimSpace(filtr.Jezyk) != "" {
		warunki = append(warunki, "jezyk = ?")
		argumenty = append(argumenty, filtr.Jezyk)
	}
	if strings.TrimSpace(filtr.Dziedzina) != "" {
		warunki = append(warunki, "dziedzina = ?")
		argumenty = append(argumenty, filtr.Dziedzina)
	}
	if strings.TrimSpace(filtr.Stan) != "" {
		warunki = append(warunki, "stan = ?")
		argumenty = append(argumenty, filtr.Stan)
	}
	if fraza := strings.TrimSpace(filtr.Fraza); fraza != "" {
		warunki = append(warunki, "(zrodlo LIKE ? OR cel LIKE ?)")
		argumenty = append(argumenty, "%"+fraza+"%", "%"+fraza+"%")
	}
	warunek := ""
	if len(warunki) > 0 {
		warunek = " WHERE " + strings.Join(warunki, " AND ")
	}

	var razem int
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM termin_slownika`+warunek, argumenty...).Scan(&razem); err != nil {
		return nil, 0, fmt.Errorf("dane: nie można policzyć terminów słownika: %w", err)
	}

	zapytanie := `SELECT ` + kolumnyTerminuSlownika + ` FROM termin_slownika` + warunek +
		` ORDER BY jezyk, zrodlo`
	if filtr.Limit > 0 {
		zapytanie += " LIMIT ?"
		argumenty = append(argumenty, filtr.Limit)
		if filtr.Offset > 0 {
			zapytanie += " OFFSET ?"
			argumenty = append(argumenty, filtr.Offset)
		}
	}
	wiersze, err := r.db.QueryContext(ctx, zapytanie, argumenty...)
	if err != nil {
		return nil, 0, fmt.Errorf("dane: nie można odczytać terminów słownika: %w", err)
	}
	defer wiersze.Close()

	lista := []TerminSlownika{}
	for wiersze.Next() {
		termin, err := odczytajTerminSlownika(wiersze)
		if err != nil {
			return nil, 0, fmt.Errorf("dane: nieczytelny wiersz terminu słownika: %w", err)
		}
		lista = append(lista, termin)
	}
	return lista, razem, wiersze.Err()
}

// ProfilQa to wiersz `profil_qa` wraz z kontrolami — profil bez kontroli nie
// jest profilem, więc oba byty chodzą razem.
type ProfilQa struct {
	ID             int64
	Kod            string
	Nazwa          string
	Opis           *string
	Zasieg         string
	ZasiegID       *string
	Kontrole       []KontrolaProfiluQa
	Zaktualizowano int64
}

// KontrolaProfiluQa to wiersz tabeli `profil_qa_kontrola`: pojedyncza
// kontrola jakości wchodząca w skład profilu, wraz z wagą i stanem włączenia.
type KontrolaProfiluQa struct {
	Rodzaj   string
	Waga     string
	Wlaczona bool
}

// ProfileQa oddaje profile kontroli jakości zawężone zasięgiem wraz z ich
// kontrolami; pusty zasięg oddaje komplet profili.
func (r *repozytoriumTlumaczen) ProfileQa(ctx context.Context,
	zasieg, zasiegID string) ([]ProfilQa, error) {

	warunki := []string{}
	argumenty := []any{}
	if strings.TrimSpace(zasieg) != "" {
		warunki = append(warunki, "zasieg = ?")
		argumenty = append(argumenty, zasieg)
	}
	if strings.TrimSpace(zasiegID) != "" {
		warunki = append(warunki, "zasieg_id = ?")
		argumenty = append(argumenty, zasiegID)
	}
	zapytanie := `SELECT id, identyfikator_zewnetrzny, nazwa, opis, zasieg, zasieg_id, zaktualizowano
	                FROM profil_qa`
	if len(warunki) > 0 {
		zapytanie += " WHERE " + strings.Join(warunki, " AND ")
	}
	zapytanie += " ORDER BY nazwa"

	wiersze, err := r.db.QueryContext(ctx, zapytanie, argumenty...)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać profili kontroli jakości: %w", err)
	}
	defer wiersze.Close()

	profile := []ProfilQa{}
	for wiersze.Next() {
		var profil ProfilQa
		var opis, zasiegKolumna sql.NullString
		if err := wiersze.Scan(&profil.ID, &profil.Kod, &profil.Nazwa, &opis,
			&profil.Zasieg, &zasiegKolumna, &profil.Zaktualizowano); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny profil kontroli jakości: %w", err)
		}
		profil.Opis = tekstZKolumny(opis)
		profil.ZasiegID = tekstZKolumny(zasiegKolumna)
		profile = append(profile, profil)
	}
	if err := wiersze.Err(); err != nil {
		return nil, err
	}
	for i := range profile {
		kontrole, err := r.kontroleProfiluQa(ctx, profile[i].ID)
		if err != nil {
			return nil, err
		}
		profile[i].Kontrole = kontrole
	}
	return profile, nil
}

// kontroleProfiluQa doczytuje kontrole jednego profilu kontroli jakości,
// uporządkowane według rodzaju kontroli.
func (r *repozytoriumTlumaczen) kontroleProfiluQa(ctx context.Context,
	profilID int64) ([]KontrolaProfiluQa, error) {

	wiersze, err := r.db.QueryContext(ctx,
		`SELECT rodzaj, waga, wlaczona FROM profil_qa_kontrola WHERE profil_id = ? ORDER BY rodzaj`,
		profilID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać kontroli profilu: %w", err)
	}
	defer wiersze.Close()

	kontrole := []KontrolaProfiluQa{}
	for wiersze.Next() {
		var kontrola KontrolaProfiluQa
		var wlaczona int64
		if err := wiersze.Scan(&kontrola.Rodzaj, &kontrola.Waga, &wlaczona); err != nil {
			return nil, fmt.Errorf("dane: nieczytelna kontrola profilu: %w", err)
		}
		kontrola.Wlaczona = wlaczona == 1
		kontrole = append(kontrole, kontrola)
	}
	return kontrole, wiersze.Err()
}

// ProfilQaPoKodzie oddaje jeden profil kontroli jakości po kodzie
// zewnętrznym, wraz z pełnym wykazem jego kontroli.
func (r *repozytoriumTlumaczen) ProfilQaPoKodzie(ctx context.Context, kod string) (ProfilQa, error) {
	var profil ProfilQa
	var opis, zasiegID sql.NullString
	err := r.db.QueryRowContext(ctx,
		`SELECT id, identyfikator_zewnetrzny, nazwa, opis, zasieg, zasieg_id, zaktualizowano
		   FROM profil_qa WHERE identyfikator_zewnetrzny = ?`, kod).
		Scan(&profil.ID, &profil.Kod, &profil.Nazwa, &opis, &profil.Zasieg, &zasiegID,
			&profil.Zaktualizowano)
	if errors.Is(err, sql.ErrNoRows) {
		return ProfilQa{}, ErrBrakWiersza
	}
	if err != nil {
		return ProfilQa{}, fmt.Errorf("dane: nieczytelny profil kontroli jakości %q: %w", kod, err)
	}
	profil.Opis = tekstZKolumny(opis)
	profil.ZasiegID = tekstZKolumny(zasiegID)
	kontrole, err := r.kontroleProfiluQa(ctx, profil.ID)
	if err != nil {
		return ProfilQa{}, err
	}
	profil.Kontrole = kontrole
	return profil, nil
}

// ZapiszProfilQa zakłada profil albo nadpisuje zastany i wymienia jego kontrole
// w całości — kontrakt nadsyła wykaz kontroli kompletem.
func (r *repozytoriumTlumaczen) ZapiszProfilQa(ctx context.Context, profil ProfilQa) (ProfilQa, error) {
	if strings.TrimSpace(profil.Kod) == "" {
		return ProfilQa{}, fmt.Errorf("dane: profil kontroli jakości bez identyfikatora")
	}
	teraz := time.Now().UnixMilli()
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		if _, err := transakcja.ExecContext(ctx, `INSERT INTO profil_qa
			(identyfikator_zewnetrzny, nazwa, opis, zasieg, zasieg_id, zaktualizowano)
			VALUES (?, ?, ?, ?, ?, ?)
			ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
				nazwa = excluded.nazwa, opis = excluded.opis, zasieg = excluded.zasieg,
				zasieg_id = excluded.zasieg_id, zaktualizowano = excluded.zaktualizowano`,
			profil.Kod, profil.Nazwa, tekstDoKolumny(profil.Opis), profil.Zasieg,
			tekstDoKolumny(profil.ZasiegID), teraz); err != nil {
			return fmt.Errorf("dane: nie można zapisać profilu kontroli jakości %q: %w", profil.Kod, err)
		}
		var profilID int64
		if err := transakcja.QueryRowContext(ctx,
			`SELECT id FROM profil_qa WHERE identyfikator_zewnetrzny = ?`, profil.Kod).
			Scan(&profilID); err != nil {
			return fmt.Errorf("dane: nie można odczytać profilu %q po zapisie: %w", profil.Kod, err)
		}
		if _, err := transakcja.ExecContext(ctx,
			`DELETE FROM profil_qa_kontrola WHERE profil_id = ?`, profilID); err != nil {
			return fmt.Errorf("dane: nie można zdjąć kontroli profilu %q: %w", profil.Kod, err)
		}
		for _, kontrola := range profil.Kontrole {
			if _, err := transakcja.ExecContext(ctx, `INSERT INTO profil_qa_kontrola
				(profil_id, rodzaj, waga, wlaczona) VALUES (?, ?, ?, ?)`,
				profilID, kontrola.Rodzaj, kontrola.Waga,
				wartoscLogicznaDoKolumny(kontrola.Wlaczona)); err != nil {
				return fmt.Errorf("dane: nie można zapisać kontroli profilu %q: %w", profil.Kod, err)
			}
		}
		return nil
	})
	if err != nil {
		return ProfilQa{}, err
	}
	return r.ProfilQaPoKodzie(ctx, profil.Kod)
}

// UsunProfilQa kasuje profil kontroli jakości wraz z jego kontrolami; klucz
// obcy kaskadowy zdejmuje wiersze zależne.
func (r *repozytoriumTlumaczen) UsunProfilQa(ctx context.Context, kod string) (bool, error) {
	wynik, err := r.db.ExecContext(ctx,
		`DELETE FROM profil_qa WHERE identyfikator_zewnetrzny = ?`, kod)
	if err != nil {
		return false, fmt.Errorf("dane: nie można usunąć profilu kontroli jakości %q: %w", kod, err)
	}
	zeszlo, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nie można ustalić skutku usunięcia profilu: %w", err)
	}
	return zeszlo > 0, nil
}

// ZatwierdzeniePanelu to wiersz tabeli `zatwierdzenie_panelu`: jeden krok
// obiegu zatwierdzeń panelu tłumaczenia.
type ZatwierdzeniePanelu struct {
	Kod       string
	PanelID   int64
	PanelKod  string
	Etap      string
	Autor     string
	Uwaga     *string
	Utworzono int64
}

// ZapiszZatwierdzenie dokłada krok obiegu i przestawia migawkę panelu w jednej
// transakcji — inaczej wiersz obiegu i stan panelu rozjechałyby się przy awarii
// między dwoma zapisami.
func (r *repozytoriumTlumaczen) ZapiszZatwierdzenie(ctx context.Context,
	zapis ZatwierdzeniePanelu) (ZatwierdzeniePanelu, error) {

	teraz := time.Now().UnixMilli()
	if zapis.Utworzono == 0 {
		zapis.Utworzono = teraz
	}
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		if _, err := transakcja.ExecContext(ctx, `INSERT INTO zatwierdzenie_panelu
			(identyfikator_zewnetrzny, panel_id, etap, autor, uwaga, utworzono)
			VALUES (?, ?, ?, ?, ?, ?)`,
			zapis.Kod, zapis.PanelID, zapis.Etap, zapis.Autor,
			tekstDoKolumny(zapis.Uwaga), zapis.Utworzono); err != nil {
			return fmt.Errorf("dane: nie można zapisać zatwierdzenia panelu: %w", err)
		}
		if _, err := transakcja.ExecContext(ctx, `UPDATE panel_tlumaczenia
			SET etap_zatwierdzenia = ?, zatwierdzil = ?, zatwierdzono = ?, zaktualizowano = ?
			WHERE id = ?`, zapis.Etap, zapis.Autor, zapis.Utworzono, teraz, zapis.PanelID); err != nil {
			return fmt.Errorf("dane: nie można przestawić etapu panelu: %w", err)
		}
		return nil
	})
	if err != nil {
		return ZatwierdzeniePanelu{}, err
	}
	return zapis, nil
}

// Zatwierdzenia oddaje obieg zatwierdzeń panelu albo, gdy `panelID` jest
// zerem, obieg zatwierdzeń wszystkich paneli okna.
func (r *repozytoriumTlumaczen) Zatwierdzenia(ctx context.Context,
	panelID, oknoID int64) ([]ZatwierdzeniePanelu, error) {

	zapytanie := `SELECT z.identyfikator_zewnetrzny, z.panel_id, p.identyfikator_zewnetrzny,
	                     z.etap, z.autor, z.uwaga, z.utworzono
	                FROM zatwierdzenie_panelu z
	                JOIN panel_tlumaczenia p ON p.id = z.panel_id`
	argumenty := []any{}
	switch {
	case panelID > 0:
		zapytanie += " WHERE z.panel_id = ?"
		argumenty = append(argumenty, panelID)
	case oknoID > 0:
		zapytanie += " WHERE p.okno_id = ?"
		argumenty = append(argumenty, oknoID)
	}
	zapytanie += " ORDER BY z.utworzono DESC, z.id DESC"

	wiersze, err := r.db.QueryContext(ctx, zapytanie, argumenty...)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać obiegu zatwierdzeń: %w", err)
	}
	defer wiersze.Close()

	zapisy := []ZatwierdzeniePanelu{}
	for wiersze.Next() {
		var zapis ZatwierdzeniePanelu
		var uwaga sql.NullString
		if err := wiersze.Scan(&zapis.Kod, &zapis.PanelID, &zapis.PanelKod, &zapis.Etap,
			&zapis.Autor, &uwaga, &zapis.Utworzono); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz obiegu zatwierdzeń: %w", err)
		}
		zapis.Uwaga = tekstZKolumny(uwaga)
		zapisy = append(zapisy, zapis)
	}
	return zapisy, wiersze.Err()
}

// UstalenieKorekty to wiersz `ustalenie_korekty` — jedno zastrzeżenie korekty
// językowej wobec treści panelu.
type UstalenieKorekty struct {
	Kod         string
	PanelID     int64
	Rodzaj      string
	Waga        string
	Segment     *string
	Szczegol    string
	Propozycja  *string
	Zastosowano *int64
	Odrzucono   *int64
	Utworzono   int64
}

// ZapiszUstaleniaKorekty wymienia otwarte ustalenia panelu na nadesłane.
// Ustalenia rozstrzygnięte (zastosowane albo odrzucone) zostają — kolejny
// przebieg nie ma prawa skasować odpowiedzi Operatora.
func (r *repozytoriumTlumaczen) ZapiszUstaleniaKorekty(ctx context.Context,
	panelID int64, ustalenia []UstalenieKorekty) error {

	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		if _, err := transakcja.ExecContext(ctx,
			`DELETE FROM ustalenie_korekty
			  WHERE panel_id = ? AND zastosowano IS NULL AND odrzucono IS NULL`, panelID); err != nil {
			return fmt.Errorf("dane: nie można zdjąć otwartych ustaleń korekty: %w", err)
		}
		for _, ustalenie := range ustalenia {
			if _, err := transakcja.ExecContext(ctx, `INSERT INTO ustalenie_korekty
				(identyfikator_zewnetrzny, panel_id, rodzaj, waga, segment, szczegol,
				 propozycja, utworzono)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
				ustalenie.Kod, panelID, ustalenie.Rodzaj, ustalenie.Waga,
				tekstDoKolumny(ustalenie.Segment), ustalenie.Szczegol,
				tekstDoKolumny(ustalenie.Propozycja), ustalenie.Utworzono); err != nil {
				return fmt.Errorf("dane: nie można zapisać ustalenia korekty: %w", err)
			}
		}
		return nil
	})
}

// UstaleniaKorekty oddaje wszystkie ustalenia korekty panelu, rozstrzygnięte
// i otwarte, w kolejności od najnowszego.
func (r *repozytoriumTlumaczen) UstaleniaKorekty(ctx context.Context,
	panelID int64) ([]UstalenieKorekty, error) {

	wiersze, err := r.db.QueryContext(ctx,
		`SELECT identyfikator_zewnetrzny, panel_id, rodzaj, waga, segment, szczegol,
		        propozycja, zastosowano, odrzucono, utworzono
		   FROM ustalenie_korekty WHERE panel_id = ? ORDER BY utworzono DESC, id DESC`, panelID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać ustaleń korekty: %w", err)
	}
	defer wiersze.Close()

	ustalenia := []UstalenieKorekty{}
	for wiersze.Next() {
		ustalenie, err := odczytajUstalenieKorekty(wiersze)
		if err != nil {
			return nil, err
		}
		ustalenia = append(ustalenia, ustalenie)
	}
	return ustalenia, wiersze.Err()
}

// UstalenieKorektyPoKodzie oddaje jedno ustalenie korekty po kodzie
// zewnętrznym, jednoznacznie identyfikującym wiersz.
func (r *repozytoriumTlumaczen) UstalenieKorektyPoKodzie(ctx context.Context,
	kod string) (UstalenieKorekty, error) {

	wiersz := r.db.QueryRowContext(ctx,
		`SELECT identyfikator_zewnetrzny, panel_id, rodzaj, waga, segment, szczegol,
		        propozycja, zastosowano, odrzucono, utworzono
		   FROM ustalenie_korekty WHERE identyfikator_zewnetrzny = ?`, kod)
	ustalenie, err := odczytajUstalenieKorekty(wiersz)
	if errors.Is(err, sql.ErrNoRows) {
		return UstalenieKorekty{}, ErrBrakWiersza
	}
	return ustalenie, err
}

// odczytajUstalenieKorekty odczytuje pojedynczą strukturę UstalenieKorekty
// z jednego wiersza wyniku zapytania SQL.
func odczytajUstalenieKorekty(wiersz skaner) (UstalenieKorekty, error) {
	var ustalenie UstalenieKorekty
	var segment, propozycja sql.NullString
	var zastosowano, odrzucono sql.NullInt64
	err := wiersz.Scan(&ustalenie.Kod, &ustalenie.PanelID, &ustalenie.Rodzaj, &ustalenie.Waga,
		&segment, &ustalenie.Szczegol, &propozycja, &zastosowano, &odrzucono, &ustalenie.Utworzono)
	if err != nil {
		return UstalenieKorekty{}, err
	}
	ustalenie.Segment = tekstZKolumny(segment)
	ustalenie.Propozycja = tekstZKolumny(propozycja)
	ustalenie.Zastosowano = liczbaZKolumny(zastosowano)
	ustalenie.Odrzucono = liczbaZKolumny(odrzucono)
	return ustalenie, nil
}

// RozstrzygnijUstalenieKorekty znakuje ustalenie jako zastosowane albo
// odrzucone. Rozstrzygnięcie jest chwilą, nie wartością logiczną: „kiedy" niesie
// więcej niż „czy", a „czy" da się z „kiedy" odczytać.
func (r *repozytoriumTlumaczen) RozstrzygnijUstalenieKorekty(ctx context.Context,
	kod string, odrzucone bool) error {

	kolumna := "zastosowano"
	if odrzucone {
		kolumna = "odrzucono"
	}
	_, err := r.db.ExecContext(ctx,
		`UPDATE ustalenie_korekty SET `+kolumna+` = ? WHERE identyfikator_zewnetrzny = ?`,
		time.Now().UnixMilli(), kod)
	if err != nil {
		return fmt.Errorf("dane: nie można rozstrzygnąć ustalenia korekty %q: %w", kod, err)
	}
	return nil
}
