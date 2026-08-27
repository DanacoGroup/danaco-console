// Plik prowadzi obszar automatyk modułu Automations: definicję automatyki wraz z kontraktem obszaru; kroki, harmonogram,
// przebiegi i kolejka leżą w osobnych plikach jednego repozytorium, a silnika wykonania kroków tu nie ma.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// Automatyka to wiersz tabeli `automatyka`. Kod jest identyfikatorem, którym
// automatyka wychodzi kontraktem (`AutomationWorkflow.id`).
type Automatyka struct {
	ID             int64
	Kod            string
	Nazwa          string
	Opis           *string
	Czynna         bool
	Wersja         int
	Utworzono      string
	Zaktualizowano string
	// WersjaOpublikowana pusta znaczy brak rozdziału wersji; wykonywana produkcyjnie jest wersja bieżąca.
	WersjaOpublikowana *int
	Udostepniona       bool
	BudzetPrzebiegu    int
	BudzetKroku        int
	RegulaBudzetu      *string
}

// RepozytoriumAutomatyk jest kontraktem obszaru Automations: definicje, kroki, harmonogram i przebiegi automatyk.
type RepozytoriumAutomatyk interface {
	ZapiszAutomatyke(ctx context.Context, automatyka Automatyka) (Automatyka, error)
	Automatyka(ctx context.Context, kod string) (Automatyka, error)
	AutomatykaPoID(ctx context.Context, id int64) (Automatyka, error)
	Automatyki(ctx context.Context, tylkoCzynne bool, limit int) ([]Automatyka, error)

	ZapiszKroki(ctx context.Context, automatykaID int64, kroki []KrokAutomatyki) error
	Kroki(ctx context.Context, automatykaID int64) ([]KrokAutomatyki, error)
	ZapiszZaleznosci(ctx context.Context, automatykaID int64, zaleznosci []ZaleznoscKroku) error
	Zaleznosci(ctx context.Context, automatykaID int64) ([]ZaleznoscKroku, error)

	ZapiszHarmonogram(ctx context.Context, harmonogram Harmonogram,
		wyzwalacze []WyzwalaczAutomatyki) (Harmonogram, error)
	Harmonogram(ctx context.Context, automatykaID int64) (Harmonogram, error)
	HarmonogramyNalezne(ctx context.Context, teraz string) ([]Harmonogram, error)
	UstawNastepneUruchomienie(ctx context.Context, automatykaID int64, nastepne *string) error
	Wyzwalacze(ctx context.Context, harmonogramID int64) ([]WyzwalaczAutomatyki, error)

	ZapiszPrzebieg(ctx context.Context, przebieg Przebieg) (Przebieg, error)
	Przebieg(ctx context.Context, kod string) (Przebieg, error)
	PrzebiegKolejki(ctx context.Context, kolejkaID int64) (Przebieg, error)
	Przebiegi(ctx context.Context, automatykaID int64, limit int) ([]Przebieg, error)

	UstawKolejnoscPozycji(ctx context.Context, pozycjaID int64, kolejnosc int) error
	PrzeniesPozycje(ctx context.Context, pozycjaID, kolejkaDocelowaID int64) error

	// Historia wersji definicji, publikacja, udostępnienie i budżety czasu automatyki, w osobnym pliku.
	ZapiszWersjeAutomatyki(ctx context.Context, wersja WersjaAutomatyki) error
	WersjeAutomatyki(ctx context.Context, automatykaID int64, limit int) ([]WersjaAutomatyki, error)
	WersjaAutomatykiNumer(ctx context.Context, automatykaID int64, numer int) (WersjaAutomatyki, error)
	OpublikujWersjeAutomatyki(ctx context.Context, automatykaID int64, numer int) error
	UstawUdostepnienieAutomatyki(ctx context.Context, automatykaID int64, udostepniona bool) error
	UstawBudzetyAutomatyki(ctx context.Context, automatykaID int64,
		budzetPrzebiegu, budzetKroku int, regula *string) error

	// Etykiety, zmienne przepływu, mapowania i adnotacje kanwy
	// (`automations_kanwa.go`, migracje 261-264).
	UstawEtykietyAutomatyki(ctx context.Context, automatykaID int64, etykiety []string) error
	EtykietyAutomatyki(ctx context.Context, automatykaID int64) ([]string, error)
	ZapiszZmienneAutomatyki(ctx context.Context, automatykaID int64, zmienne []ZmiennaAutomatyki) error
	ZmienneAutomatyki(ctx context.Context, automatykaID int64) ([]ZmiennaAutomatyki, error)
	ZapiszMapowaniaAutomatyki(ctx context.Context, automatykaID int64, mapowania []MapowanieDanych) error
	MapowaniaAutomatyki(ctx context.Context, automatykaID int64) ([]MapowanieDanych, error)
	UstawNotatkeKroku(ctx context.Context, automatykaID int64, krokKod string, notatka *string) error
	ZapiszPolozeniaKrokow(ctx context.Context, automatykaID int64, polozenia []AdnotacjaKroku) error
	AdnotacjeKrokow(ctx context.Context, automatykaID int64) ([]AdnotacjaKroku, error)

	// Biblioteka szablonów przepływów (`automations_szablony.go`, migracje 265-266).
	ZapiszSzablonAutomatyki(ctx context.Context, szablon SzablonAutomatyki,
		parametry []ParametrSzablonu) (SzablonAutomatyki, error)
	SzablonAutomatyki(ctx context.Context, kod string) (SzablonAutomatyki, error)
	SzablonyAutomatyki(ctx context.Context, limit int) ([]SzablonAutomatyki, error)
	ParametrySzablonuAutomatyki(ctx context.Context, szablonID int64) ([]ParametrSzablonu, error)

	// Trwały zapis przebiegu automatyki: log, kroki wraz z ładunkami, punkty wznowienia, w osobnym pliku.
	DopiszLogPrzebiegu(ctx context.Context, wpis WpisLoguPrzebiegu) error
	LogPrzebiegu(ctx context.Context, przebiegID int64, poziom, krokKod string,
		limit int) ([]WpisLoguPrzebiegu, error)
	ZapiszKrokPrzebiegu(ctx context.Context, krok KrokPrzebiegu) error
	KrokiPrzebiegu(ctx context.Context, przebiegID int64) ([]KrokPrzebiegu, error)
	KrokPrzebieguPoKodzie(ctx context.Context, przebiegID int64, krokKod string) (KrokPrzebiegu, error)
	ZapiszPunktWznowienia(ctx context.Context, punkt PunktWznowienia) error
	PunktyWznowienia(ctx context.Context, przebiegID int64) ([]PunktWznowienia, error)
	PunktWznowieniaPoKodzie(ctx context.Context, kod string) (PunktWznowienia, error)

	// Reguły alarmowania, skarbiec referencji i dziennik audytu automatyki, w osobnym pliku repozytorium.
	ZapiszRegulealarmowania(ctx context.Context, regula RegulaAlarmowania) (RegulaAlarmowania, error)
	RegulyAlarmowania(ctx context.Context, automatykaID int64) ([]RegulaAlarmowania, error)
	ZapiszPoswiadczenieAutomatyki(ctx context.Context, poswiadczenie PoswiadczenieAutomatyki) error
	PoswiadczeniaAutomatyki(ctx context.Context, zasieg, zasiegID string) ([]PoswiadczenieAutomatyki, error)
	PoswiadczenieAutomatykiPoOdwolaniu(ctx context.Context, odwolanie string) (PoswiadczenieAutomatyki, error)
	UsunPoswiadczenieAutomatyki(ctx context.Context, odwolanie string) (bool, error)
	DopiszAudytAutomatyki(ctx context.Context, wpis WpisAudytuAutomatyki) error
	AudytAutomatyki(ctx context.Context, automatykaID int64, od, do string,
		limit int) ([]WpisAudytuAutomatyki, error)

	// Harmonogram po własnym identyfikatorze, okna wykonania, nadzór obecności i historia wyzwoleń.
	HarmonogramPoKodzie(ctx context.Context, kod string) (Harmonogram, error)
	Harmonogramy(ctx context.Context, tylkoCzynne bool) ([]Harmonogram, error)
	ZapiszOknaWykonania(ctx context.Context, harmonogramID int64, okna []OknoWykonania) error
	OknaWykonania(ctx context.Context, harmonogramID int64) ([]OknoWykonania, error)
	UstawNadzorHarmonogramu(ctx context.Context, harmonogramID int64,
		tolerancjaSekundy int, regula *string) error
	UstawPodpisHarmonogramu(ctx context.Context, harmonogramID int64, odwolanie string) error
	DopiszWyzwolenie(ctx context.Context, automatykaID int64, harmonogramID *int64,
		przyczyna string, przebiegID *int64) error
	Wyzwolenia(ctx context.Context, automatykaID, harmonogramID int64,
		limit int) ([]WyzwolenieAutomatyki, error)
}

const (
	kolumnyAutomatyki = `id, identyfikator_zewnetrzny, nazwa, opis, czynna, wersja,
	                     utworzono, zaktualizowano, wersja_opublikowana, udostepniona,
	                     budzet_przebiegu_sekundy, budzet_kroku_sekundy, regula_budzetu`

	// Zapis podnosi wersję definicji przy każdej zmianie — panel akcji Workflow
	// Buildera ma pozycję „Wersje” i musi mieć co pokazać.
	zapiszAutomatyke = `INSERT INTO automatyka
	                    (identyfikator_zewnetrzny, nazwa, opis, czynna)
	                    VALUES (?, ?, ?, ?)
	                    ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                        nazwa = excluded.nazwa,
	                        opis = excluded.opis,
	                        czynna = excluded.czynna,
	                        wersja = automatyka.wersja + 1,
	                        zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	pobierzAutomatyke = `SELECT ` + kolumnyAutomatyki + ` FROM automatyka
	                     WHERE identyfikator_zewnetrzny = ?`

	pobierzAutomatykePoID = `SELECT ` + kolumnyAutomatyki + ` FROM automatyka WHERE id = ?`

	listaAutomatyk = `SELECT ` + kolumnyAutomatyki + ` FROM automatyka
	                  WHERE (? = 0 OR czynna = 1)
	                  ORDER BY nazwa, id LIMIT ?`
)

type repozytoriumAutomatyk struct {
	zapytania *zapytania
	db        *sql.DB
}

func noweRepozytoriumAutomatyk(z *zapytania, db *sql.DB) *repozytoriumAutomatyk {
	return &repozytoriumAutomatyk{zapytania: z, db: db}
}

// ZapiszAutomatyke zakłada automatykę albo nadpisuje zastaną i zwraca stan po
// zapisie wraz z podniesioną wersją definicji.
func (r *repozytoriumAutomatyk) ZapiszAutomatyke(ctx context.Context,
	automatyka Automatyka) (Automatyka, error) {

	if automatyka.Kod == "" {
		return Automatyka{}, fmt.Errorf("dane: automatyka bez identyfikatora")
	}
	if automatyka.Nazwa == "" {
		automatyka.Nazwa = automatyka.Kod
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszAutomatyke)
	if err != nil {
		return Automatyka{}, err
	}
	_, err = polecenie.ExecContext(ctx, automatyka.Kod, automatyka.Nazwa,
		tekstDoKolumny(automatyka.Opis), liczbaLogiczna(automatyka.Czynna))
	if err != nil {
		return Automatyka{}, fmt.Errorf("dane: nie można zapisać automatyki %q: %w", automatyka.Kod, err)
	}
	return r.Automatyka(ctx, automatyka.Kod)
}

// Automatyka zwraca definicję o wskazanym kodzie. Brak wiersza wraca jako
// ErrBrakWiersza — warstwa wyższa odróżnia „nie ma” od „odczyt się nie powiódł”.
func (r *repozytoriumAutomatyk) Automatyka(ctx context.Context, kod string) (Automatyka, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzAutomatyke)
	if err != nil {
		return Automatyka{}, err
	}
	automatyka, err := odczytajAutomatyke(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return Automatyka{}, ErrBrakWiersza
	}
	if err != nil {
		return Automatyka{}, fmt.Errorf("dane: nieczytelny wiersz automatyki %q: %w", kod, err)
	}
	return automatyka, nil
}

// AutomatykaPoID zwraca definicję po kluczu wiersza. Służy budzikowi
// harmonogramu, który zna automatykę przez klucz obcy harmonogramu
// (`harmonogram_automatyki.automatyka_id`), a nie przez identyfikator zewnętrzny.
func (r *repozytoriumAutomatyk) AutomatykaPoID(ctx context.Context, id int64) (Automatyka, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzAutomatykePoID)
	if err != nil {
		return Automatyka{}, err
	}
	automatyka, err := odczytajAutomatyke(polecenie.QueryRowContext(ctx, id))
	if errors.Is(err, sql.ErrNoRows) {
		return Automatyka{}, ErrBrakWiersza
	}
	if err != nil {
		return Automatyka{}, fmt.Errorf("dane: nieczytelny wiersz automatyki %d: %w", id, err)
	}
	return automatyka, nil
}

// Automatyki zwraca definicje uporządkowane po nazwie — w takiej kolejności
// pokazuje je wykaz Workflow Buildera.
func (r *repozytoriumAutomatyk) Automatyki(ctx context.Context,
	tylkoCzynne bool, limit int) ([]Automatyka, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaAutomatyk)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, liczbaLogiczna(tylkoCzynne), granicaWykazu(limit))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać automatyk: %w", err)
	}
	defer wiersze.Close()

	lista := []Automatyka{}
	for wiersze.Next() {
		automatyka, err := odczytajAutomatyke(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz automatyki: %w", err)
		}
		lista = append(lista, automatyka)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt automatyk: %w", err)
	}
	return lista, nil
}

// odczytajAutomatyke składa strukturę automatyki wprost z jednego wiersza wyniku zapytania do bazy SQL.
func odczytajAutomatyke(wiersz skaner) (Automatyka, error) {
	var automatyka Automatyka
	var opis, regula sql.NullString
	var opublikowana sql.NullInt64
	var czynna, udostepniona int
	err := wiersz.Scan(&automatyka.ID, &automatyka.Kod, &automatyka.Nazwa, &opis,
		&czynna, &automatyka.Wersja, &automatyka.Utworzono, &automatyka.Zaktualizowano,
		&opublikowana, &udostepniona, &automatyka.BudzetPrzebiegu, &automatyka.BudzetKroku,
		&regula)
	if err != nil {
		return Automatyka{}, err
	}
	automatyka.Opis = tekstZKolumny(opis)
	automatyka.Czynna = czynna == 1
	automatyka.Udostepniona = udostepniona == 1
	automatyka.RegulaBudzetu = tekstZKolumny(regula)
	if opublikowana.Valid {
		numer := int(opublikowana.Int64)
		automatyka.WersjaOpublikowana = &numer
	}
	return automatyka, nil
}
