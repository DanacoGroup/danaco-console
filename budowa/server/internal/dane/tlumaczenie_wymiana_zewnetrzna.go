// Plik utrzymuje nastawy silników i wymianę zewnętrzną modułu Translate:
// profile silników wraz z kanałami, politykę tłumaczenia pivotowego wraz
// z parami języków, pakiety przekazania wykonawcy, przebieg pakietowy i most do dokumentu.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// ProfilSilnika to wiersz tabeli profil_silnika_tlumaczenia wraz z kanałami przypisanymi mu w osobnej tabeli.
type ProfilSilnika struct {
	ID             int64
	Kod            string
	Nazwa          string
	Dziedzina      *string
	Zasieg         string
	ZasiegID       *string
	Kanaly         []string
	Adaptacyjny    bool
	ZasiegPamieci  *string
	Temperatura    *float64
	Zaktualizowano int64
}

// ProfileSilnikow oddaje profile silników zawężone zasięgiem i identyfikatorem zasięgu, w kolejności nazwy.
func (r *repozytoriumTlumaczen) ProfileSilnikow(ctx context.Context,
	zasieg, zasiegID string) ([]ProfilSilnika, error) {

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
	zapytanie := `SELECT id, identyfikator_zewnetrzny, nazwa, dziedzina, zasieg, zasieg_id,
	                     adaptacyjny, zasieg_pamieci, temperatura, zaktualizowano
	                FROM profil_silnika_tlumaczenia`
	if len(warunki) > 0 {
		zapytanie += " WHERE " + strings.Join(warunki, " AND ")
	}
	zapytanie += " ORDER BY nazwa"

	wiersze, err := r.db.QueryContext(ctx, zapytanie, argumenty...)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać profili silników: %w", err)
	}
	defer wiersze.Close()

	profile := []ProfilSilnika{}
	for wiersze.Next() {
		profil, err := odczytajProfilSilnika(wiersze)
		if err != nil {
			return nil, err
		}
		profile = append(profile, profil)
	}
	if err := wiersze.Err(); err != nil {
		return nil, err
	}
	for i := range profile {
		kanaly, err := r.kanalyProfiluSilnika(ctx, profile[i].ID)
		if err != nil {
			return nil, err
		}
		profile[i].Kanaly = kanaly
	}
	return profile, nil
}

// odczytajProfilSilnika składa strukturę profilu silnika tłumaczenia z jednego wiersza wyniku zapytania.
func odczytajProfilSilnika(wiersz skaner) (ProfilSilnika, error) {
	var profil ProfilSilnika
	var dziedzina, zasiegID, zasiegPamieci sql.NullString
	var adaptacyjny int64
	var temperatura sql.NullFloat64
	err := wiersz.Scan(&profil.ID, &profil.Kod, &profil.Nazwa, &dziedzina, &profil.Zasieg,
		&zasiegID, &adaptacyjny, &zasiegPamieci, &temperatura, &profil.Zaktualizowano)
	if err != nil {
		return ProfilSilnika{}, fmt.Errorf("dane: nieczytelny profil silnika: %w", err)
	}
	profil.Dziedzina = tekstZKolumny(dziedzina)
	profil.ZasiegID = tekstZKolumny(zasiegID)
	profil.ZasiegPamieci = tekstZKolumny(zasiegPamieci)
	profil.Adaptacyjny = adaptacyjny == 1
	if temperatura.Valid {
		wartosc := temperatura.Float64
		profil.Temperatura = &wartosc
	}
	return profil, nil
}

// kanalyProfiluSilnika doczytuje kanały danego profilu silnika w kolejności wskazania przy jego zapisie.
func (r *repozytoriumTlumaczen) kanalyProfiluSilnika(ctx context.Context,
	profilID int64) ([]string, error) {

	wiersze, err := r.db.QueryContext(ctx,
		`SELECT kanal_kod FROM kanal_profilu_silnika WHERE profil_id = ? ORDER BY kolejnosc`, profilID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać kanałów profilu silnika: %w", err)
	}
	defer wiersze.Close()

	kanaly := []string{}
	for wiersze.Next() {
		var kanal string
		if err := wiersze.Scan(&kanal); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny kanał profilu silnika: %w", err)
		}
		kanaly = append(kanaly, kanal)
	}
	return kanaly, wiersze.Err()
}

// ProfilSilnikaPoKodzie oddaje jeden profil silnika wraz z jego kanałami po kodzie zewnętrznym profilu.
func (r *repozytoriumTlumaczen) ProfilSilnikaPoKodzie(ctx context.Context,
	kod string) (ProfilSilnika, error) {

	wiersz := r.db.QueryRowContext(ctx,
		`SELECT id, identyfikator_zewnetrzny, nazwa, dziedzina, zasieg, zasieg_id,
		        adaptacyjny, zasieg_pamieci, temperatura, zaktualizowano
		   FROM profil_silnika_tlumaczenia WHERE identyfikator_zewnetrzny = ?`, kod)
	profil, err := odczytajProfilSilnika(wiersz)
	if errors.Is(err, sql.ErrNoRows) {
		return ProfilSilnika{}, ErrBrakWiersza
	}
	if err != nil {
		return ProfilSilnika{}, err
	}
	kanaly, err := r.kanalyProfiluSilnika(ctx, profil.ID)
	if err != nil {
		return ProfilSilnika{}, err
	}
	profil.Kanaly = kanaly
	return profil, nil
}

// ZapiszProfilSilnika zakłada profil silnika tłumaczenia albo nadpisuje zastany i wymienia jego kanały.
func (r *repozytoriumTlumaczen) ZapiszProfilSilnika(ctx context.Context,
	profil ProfilSilnika) (ProfilSilnika, error) {

	if strings.TrimSpace(profil.Kod) == "" {
		return ProfilSilnika{}, fmt.Errorf("dane: profil silnika bez identyfikatora")
	}
	teraz := time.Now().UnixMilli()
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		var temperatura any
		if profil.Temperatura != nil {
			temperatura = *profil.Temperatura
		}
		if _, err := transakcja.ExecContext(ctx, `INSERT INTO profil_silnika_tlumaczenia
			(identyfikator_zewnetrzny, nazwa, dziedzina, zasieg, zasieg_id, adaptacyjny,
			 zasieg_pamieci, temperatura, zaktualizowano)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
				nazwa = excluded.nazwa, dziedzina = excluded.dziedzina,
				zasieg = excluded.zasieg, zasieg_id = excluded.zasieg_id,
				adaptacyjny = excluded.adaptacyjny, zasieg_pamieci = excluded.zasieg_pamieci,
				temperatura = excluded.temperatura, zaktualizowano = excluded.zaktualizowano`,
			profil.Kod, profil.Nazwa, tekstDoKolumny(profil.Dziedzina), profil.Zasieg,
			tekstDoKolumny(profil.ZasiegID), wartoscLogicznaDoKolumny(profil.Adaptacyjny),
			tekstDoKolumny(profil.ZasiegPamieci), temperatura, teraz); err != nil {
			return fmt.Errorf("dane: nie można zapisać profilu silnika %q: %w", profil.Kod, err)
		}
		var profilID int64
		if err := transakcja.QueryRowContext(ctx,
			`SELECT id FROM profil_silnika_tlumaczenia WHERE identyfikator_zewnetrzny = ?`,
			profil.Kod).Scan(&profilID); err != nil {
			return fmt.Errorf("dane: nie można odczytać profilu silnika %q: %w", profil.Kod, err)
		}
		if _, err := transakcja.ExecContext(ctx,
			`DELETE FROM kanal_profilu_silnika WHERE profil_id = ?`, profilID); err != nil {
			return fmt.Errorf("dane: nie można zdjąć kanałów profilu %q: %w", profil.Kod, err)
		}
		for numer, kanal := range profil.Kanaly {
			if _, err := transakcja.ExecContext(ctx,
				`INSERT INTO kanal_profilu_silnika (profil_id, kanal_kod, kolejnosc) VALUES (?, ?, ?)`,
				profilID, kanal, numer); err != nil {
				return fmt.Errorf("dane: nie można zapisać kanału profilu %q: %w", profil.Kod, err)
			}
		}
		return nil
	})
	if err != nil {
		return ProfilSilnika{}, err
	}
	return r.ProfilSilnikaPoKodzie(ctx, profil.Kod)
}

// PolitykaPivota to wiersz tabeli polityka_pivota wraz z parami języków, przez które idzie tłumaczenie pośrednie.
type PolitykaPivota struct {
	ID             int64
	Zasieg         string
	ZasiegID       string
	JezykDomyslny  *string
	Pary           []ParaPivota
	Zaktualizowano int64
}

// ParaPivota to wiersz tabeli para_pivota: wskazuje, przez jaki język pośredni idzie dana para źródło-cel.
type ParaPivota struct {
	JezykZrodla string
	JezykCelu   string
	JezykPivota string
}

// PolitykaPivotaZasiegu oddaje politykę zasięgu wraz z parami. Brak wiersza
// wraca jako ErrBrakWiersza — polityki nieustawionej nie udajemy pustą.
func (r *repozytoriumTlumaczen) PolitykaPivotaZasiegu(ctx context.Context,
	zasieg, zasiegID string) (PolitykaPivota, error) {

	var polityka PolitykaPivota
	var domyslny sql.NullString
	err := r.db.QueryRowContext(ctx,
		`SELECT id, zasieg, zasieg_id, jezyk_domyslny, zaktualizowano
		   FROM polityka_pivota WHERE zasieg = ? AND zasieg_id = ?`, zasieg, zasiegID).
		Scan(&polityka.ID, &polityka.Zasieg, &polityka.ZasiegID, &domyslny, &polityka.Zaktualizowano)
	if errors.Is(err, sql.ErrNoRows) {
		return PolitykaPivota{}, ErrBrakWiersza
	}
	if err != nil {
		return PolitykaPivota{}, fmt.Errorf("dane: nieczytelna polityka pivota: %w", err)
	}
	polityka.JezykDomyslny = tekstZKolumny(domyslny)

	wiersze, err := r.db.QueryContext(ctx,
		`SELECT jezyk_zrodla, jezyk_celu, jezyk_pivota FROM para_pivota
		  WHERE polityka_id = ? ORDER BY jezyk_zrodla, jezyk_celu`, polityka.ID)
	if err != nil {
		return PolitykaPivota{}, fmt.Errorf("dane: nie można odczytać par pivota: %w", err)
	}
	defer wiersze.Close()
	polityka.Pary = []ParaPivota{}
	for wiersze.Next() {
		var para ParaPivota
		if err := wiersze.Scan(&para.JezykZrodla, &para.JezykCelu, &para.JezykPivota); err != nil {
			return PolitykaPivota{}, fmt.Errorf("dane: nieczytelna para pivota: %w", err)
		}
		polityka.Pary = append(polityka.Pary, para)
	}
	return polityka, wiersze.Err()
}

// ZapiszPolitykePivota zakłada politykę pivota zasięgu albo nadpisuje zastaną i wymienia jej pary w całości.
func (r *repozytoriumTlumaczen) ZapiszPolitykePivota(ctx context.Context,
	polityka PolitykaPivota) (PolitykaPivota, error) {

	teraz := time.Now().UnixMilli()
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		if _, err := transakcja.ExecContext(ctx, `INSERT INTO polityka_pivota
			(zasieg, zasieg_id, jezyk_domyslny, zaktualizowano) VALUES (?, ?, ?, ?)
			ON CONFLICT(zasieg, zasieg_id) DO UPDATE SET
				jezyk_domyslny = excluded.jezyk_domyslny,
				zaktualizowano = excluded.zaktualizowano`,
			polityka.Zasieg, polityka.ZasiegID, tekstDoKolumny(polityka.JezykDomyslny),
			teraz); err != nil {
			return fmt.Errorf("dane: nie można zapisać polityki pivota: %w", err)
		}
		var politykaID int64
		if err := transakcja.QueryRowContext(ctx,
			`SELECT id FROM polityka_pivota WHERE zasieg = ? AND zasieg_id = ?`,
			polityka.Zasieg, polityka.ZasiegID).Scan(&politykaID); err != nil {
			return fmt.Errorf("dane: nie można odczytać polityki pivota po zapisie: %w", err)
		}
		if _, err := transakcja.ExecContext(ctx,
			`DELETE FROM para_pivota WHERE polityka_id = ?`, politykaID); err != nil {
			return fmt.Errorf("dane: nie można zdjąć par pivota: %w", err)
		}
		for _, para := range polityka.Pary {
			if _, err := transakcja.ExecContext(ctx, `INSERT INTO para_pivota
				(polityka_id, jezyk_zrodla, jezyk_celu, jezyk_pivota) VALUES (?, ?, ?, ?)`,
				politykaID, para.JezykZrodla, para.JezykCelu, para.JezykPivota); err != nil {
				return fmt.Errorf("dane: nie można zapisać pary pivota: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return PolitykaPivota{}, err
	}
	return r.PolitykaPivotaZasiegu(ctx, polityka.Zasieg, polityka.ZasiegID)
}

// PakietPrzekazania to wiersz tabeli pakiet_przekazania, jednego pakietu materiału przekazanego wykonawcy.
type PakietPrzekazania struct {
	Kod            string
	OknoID         int64
	OknoKod        string
	Zawartosci     []string
	Panele         []string
	Instrukcje     *string
	Sciezka        string
	Stan           string
	Utworzono      int64
	Zaktualizowano int64
}

// ZapiszPakietPrzekazania zakłada pakiet albo nadpisuje zastany po kodzie —
// zwrot materiału (`handoff.receive`) przestawia ten sam wiersz na stan
// `returned`, zamiast zakładać drugi.
func (r *repozytoriumTlumaczen) ZapiszPakietPrzekazania(ctx context.Context,
	pakiet PakietPrzekazania) (PakietPrzekazania, error) {

	if strings.TrimSpace(pakiet.Kod) == "" {
		return PakietPrzekazania{}, fmt.Errorf("dane: pakiet przekazania bez identyfikatora")
	}
	teraz := time.Now().UnixMilli()
	utworzono := pakiet.Utworzono
	if utworzono == 0 {
		utworzono = teraz
	}
	_, err := r.db.ExecContext(ctx, `INSERT INTO pakiet_przekazania
		(identyfikator_zewnetrzny, okno_id, zawartosci, panele, instrukcje, sciezka,
		 stan, utworzono, zaktualizowano)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
			zawartosci = excluded.zawartosci, panele = excluded.panele,
			instrukcje = excluded.instrukcje, sciezka = excluded.sciezka,
			stan = excluded.stan, zaktualizowano = excluded.zaktualizowano`,
		pakiet.Kod, pakiet.OknoID, strings.Join(pakiet.Zawartosci, "\n"),
		strings.Join(pakiet.Panele, "\n"), tekstDoKolumny(pakiet.Instrukcje),
		pakiet.Sciezka, pakiet.Stan, utworzono, teraz)
	if err != nil {
		return PakietPrzekazania{}, fmt.Errorf("dane: nie można zapisać pakietu przekazania %q: %w",
			pakiet.Kod, err)
	}
	return r.PakietPrzekazania(ctx, pakiet.Kod)
}

// PakietPrzekazania oddaje pakiet przekazania po jego kodzie zewnętrznym wraz z kodem okna źródłowego.
func (r *repozytoriumTlumaczen) PakietPrzekazania(ctx context.Context,
	kod string) (PakietPrzekazania, error) {

	var pakiet PakietPrzekazania
	var zawartosci, panele string
	var instrukcje sql.NullString
	err := r.db.QueryRowContext(ctx,
		`SELECT p.identyfikator_zewnetrzny, p.okno_id, o.identyfikator_zewnetrzny,
		        p.zawartosci, p.panele, p.instrukcje, p.sciezka, p.stan, p.utworzono, p.zaktualizowano
		   FROM pakiet_przekazania p
		   JOIN okno_tlumaczenia o ON o.id = p.okno_id
		  WHERE p.identyfikator_zewnetrzny = ?`, kod).
		Scan(&pakiet.Kod, &pakiet.OknoID, &pakiet.OknoKod, &zawartosci, &panele,
			&instrukcje, &pakiet.Sciezka, &pakiet.Stan, &pakiet.Utworzono, &pakiet.Zaktualizowano)
	if errors.Is(err, sql.ErrNoRows) {
		return PakietPrzekazania{}, ErrBrakWiersza
	}
	if err != nil {
		return PakietPrzekazania{}, fmt.Errorf("dane: nieczytelny pakiet przekazania %q: %w", kod, err)
	}
	pakiet.Instrukcje = tekstZKolumny(instrukcje)
	if zawartosci != "" {
		pakiet.Zawartosci = strings.Split(zawartosci, "\n")
	}
	if panele != "" {
		pakiet.Panele = strings.Split(panele, "\n")
	}
	return pakiet, nil
}

// PozycjaPakietuTlumaczenia to wiersz tabeli pozycja_pakietu_tlumaczenia, jednego panelu w przebiegu pakietowym.
type PozycjaPakietuTlumaczenia struct {
	PanelKod string
	Operacja string
	Stan     string
	Szczegol *string
}

// ZalozZleceniePakietu zakłada przebieg pakietowy wraz z pozycjami w jednej
// transakcji i oddaje kod zlecenia.
func (r *repozytoriumTlumaczen) ZalozZleceniePakietu(ctx context.Context,
	kod string, oknoID int64, pozycje []PozycjaPakietuTlumaczenia) error {

	teraz := time.Now().UnixMilli()
	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		wynik, err := transakcja.ExecContext(ctx, `INSERT INTO zlecenie_pakietu_tlumaczenia
			(identyfikator_zewnetrzny, okno_id, utworzono) VALUES (?, ?, ?)`, kod, oknoID, teraz)
		if err != nil {
			return fmt.Errorf("dane: nie można założyć zlecenia pakietowego %q: %w", kod, err)
		}
		zlecenieID, err := wynik.LastInsertId()
		if err != nil {
			return fmt.Errorf("dane: nie można ustalić klucza zlecenia pakietowego: %w", err)
		}
		for _, pozycja := range pozycje {
			if _, err := transakcja.ExecContext(ctx, `INSERT INTO pozycja_pakietu_tlumaczenia
				(zlecenie_id, panel_kod, operacja, stan, szczegol, zaktualizowano)
				VALUES (?, ?, ?, ?, ?, ?)`,
				zlecenieID, pozycja.PanelKod, pozycja.Operacja, pozycja.Stan,
				tekstDoKolumny(pozycja.Szczegol), teraz); err != nil {
				return fmt.Errorf("dane: nie można zapisać pozycji pakietu: %w", err)
			}
		}
		return nil
	})
}

// PozycjePakietu oddaje pozycje przebiegu pakietowego po kodzie zlecenia, w kolejności ich zapisu do tabeli.
func (r *repozytoriumTlumaczen) PozycjePakietu(ctx context.Context,
	kod string) ([]PozycjaPakietuTlumaczenia, error) {

	wiersze, err := r.db.QueryContext(ctx,
		`SELECT p.panel_kod, p.operacja, p.stan, p.szczegol
		   FROM pozycja_pakietu_tlumaczenia p
		   JOIN zlecenie_pakietu_tlumaczenia z ON z.id = p.zlecenie_id
		  WHERE z.identyfikator_zewnetrzny = ? ORDER BY p.id`, kod)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać pozycji pakietu %q: %w", kod, err)
	}
	defer wiersze.Close()

	pozycje := []PozycjaPakietuTlumaczenia{}
	for wiersze.Next() {
		var pozycja PozycjaPakietuTlumaczenia
		var szczegol sql.NullString
		if err := wiersze.Scan(&pozycja.PanelKod, &pozycja.Operacja, &pozycja.Stan, &szczegol); err != nil {
			return nil, fmt.Errorf("dane: nieczytelna pozycja pakietu: %w", err)
		}
		pozycja.Szczegol = tekstZKolumny(szczegol)
		pozycje = append(pozycje, pozycja)
	}
	return pozycje, wiersze.Err()
}

// MostTlumaczenia to wiersz `most_tlumaczenia` — wiązanie okna z dokumentem,
// z którego okno wzięło materiał.
type MostTlumaczenia struct {
	OknoID        int64
	DokumentKod   string
	Zakotwiczenie *string
}

// ZapiszMost zakłada wiązanie okna tłumaczenia z dokumentem albo przestawia zastane wiązanie na nowy dokument.
func (r *repozytoriumTlumaczen) ZapiszMost(ctx context.Context, most MostTlumaczenia) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO most_tlumaczenia
		(okno_id, dokument_kod, zakotwiczenie, zaktualizowano) VALUES (?, ?, ?, ?)
		ON CONFLICT(okno_id) DO UPDATE SET
			dokument_kod = excluded.dokument_kod, zakotwiczenie = excluded.zakotwiczenie,
			zaktualizowano = excluded.zaktualizowano`,
		most.OknoID, most.DokumentKod, tekstDoKolumny(most.Zakotwiczenie),
		time.Now().UnixMilli())
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać mostu okna %d: %w", most.OknoID, err)
	}
	return nil
}

// Most oddaje wiązanie okna z dokumentem tłumaczenia po identyfikatorze okna. Brak wiersza wraca jako ErrBrakWiersza.
func (r *repozytoriumTlumaczen) Most(ctx context.Context, oknoID int64) (MostTlumaczenia, error) {
	var most MostTlumaczenia
	var zakotwiczenie sql.NullString
	err := r.db.QueryRowContext(ctx,
		`SELECT okno_id, dokument_kod, zakotwiczenie FROM most_tlumaczenia WHERE okno_id = ?`,
		oknoID).Scan(&most.OknoID, &most.DokumentKod, &zakotwiczenie)
	if errors.Is(err, sql.ErrNoRows) {
		return MostTlumaczenia{}, ErrBrakWiersza
	}
	if err != nil {
		return MostTlumaczenia{}, fmt.Errorf("dane: nieczytelny most okna %d: %w", oknoID, err)
	}
	most.Zakotwiczenie = tekstZKolumny(zakotwiczenie)
	return most, nil
}
