// Odpowiedzialność pliku: podagenci MultitaskingAI (tabela `podagent`) —
// trwałość panelu Subagent Network i panelu zadań w tle.
//
// Podagent to zadanie w tle, którego trwałą tożsamością jest pozycja kolejki,
// a proces modelu jest wyłącznie sposobem jej wykonania. Stąd kształt tego
// repozytorium: wiersz zna swoją pozycję kolejki (`pozycja_kolejki_id`), a cyklu
// życia zlecenia nie prowadzi — prowadzi go silnik kolejek.
//
// Nie ma tu odczytu tabeli `pozycja_kolejki`: jej czytelnikiem jest repozytorium
// kolejek. Wiersz podagenta niesie wyłącznie to, czego pozycja nie wie — pod kim
// biegnie, jak się nazywa i ile kosztował.
//
// Repozytorium wchodzi metodą zestawu, nie polem — ten sam wzorzec co
// `Zestaw.Rozszerzenia()` i `Zestaw.RoleOkien()`: rejestr nie trzyma stanu poza
// wskaźnikiem na wspólną pamięć zapytań, więc złożenie go na żądanie kosztuje
// tyle, co odczyt pola.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// Podagent to wiersz tabeli `podagent` wraz z identyfikatorami zewnętrznymi
// okna i sesji — kontrakt niesie je, a nie klucze wierszy.
type Podagent struct {
	ID               int64
	Kod              string
	OknoID           int64
	OknoKod          *string
	SesjaKod         *string
	BiegID           *int64
	PozycjaKolejkiID *int64
	Nazwa            *string
	Zadanie          string
	Stan             string
	Wynik            *string
	Zetonow          int
	Narzedzi         int
	Rozpoczeto       *string
	Zakonczono       *string
	Utworzono        string
}

// FiltrPodagentow zawęża wykaz. Pole puste znaczy „bez zawężenia" — pusty filtr
// oddaje komplet, a pusty wykaz jest poprawną odpowiedzią, nie odmową.
type FiltrPodagentow struct {
	OknoKod  string
	SesjaKod string
	Stan     string
}

// RepozytoriumPodagentow jest kontraktem obszaru podagentów. Rdzeń
// bierze je metodą `Zestaw.Podagenci()`.
type RepozytoriumPodagentow interface {
	ZalozPodagentow(ctx context.Context, podagenci []Podagent) ([]Podagent, error)
	Podagenci(ctx context.Context, filtr FiltrPodagentow) ([]Podagent, error)
	PodagenciPoKodach(ctx context.Context, kody []string) ([]Podagent, error)
	PrzypiszPozycje(ctx context.Context, kod string, pozycjaID int64) error
	UstawStan(ctx context.Context, kod, stan string, wynik *string) error
	// ZapiszWynikPozycji utrwala zebraną treść tury na wierszu podagenta
	// wskazanym pozycją kolejki — tak pyta ujście wyniku silnika, które zna
	// pozycję, a podagenta nie. Pozycja bez podagenta jest zapisem donikąd
	// i nie jest błędem: silnik jest jeden, więc
	// tą drogą przechodzą też pozycje pętli sesyjnej i Automations.
	ZapiszWynikPozycji(ctx context.Context, pozycjaID int64, wynik string) error
	// ZywotnoscPodagentow dokłada odpowiedzi na pytanie „co się z podagentem
	// dzieje po awarii rdzenia i po restarcie". Idzie osobnym
	// kontraktem, bo osobna jest odpowiedzialność: powyżej stoi trwałość
	// powołania i wyniku, poniżej — przynależność pracy do uruchomienia.
	ZywotnoscPodagentow
}

// Stany podagenta — słownik zamknięty więzem CHECK, zgodny
// z wyliczeniem `SubagentStatus` kontraktu.
const (
	StanPodagentaOczekuje   = "pending"
	StanPodagentaWBiegu     = "running"
	StanPodagentaUkonczony  = "done"
	StanPodagentaBledny     = "failed"
	StanPodagentaZatrzymany = "stopped"
)

const (
	kolumnyPodagenta = `p.id, p.identyfikator_zewnetrzny, p.okno_wykonawcy_id,
	                    o.identyfikator_zewnetrzny, s.identyfikator_zewnetrzny,
	                    p.bieg_id, p.pozycja_kolejki_id, p.nazwa, p.zadanie, p.stan,
	                    p.wynik, p.zetonow, p.narzedzi, p.rozpoczeto, p.zakonczono,
	                    p.utworzono`

	// Sesja wchodzi przez okno, bo `podagent` kolumny sesji nie ma i mieć nie
	// powinien: sesja podagenta jest sesją okna, które go powołało.
	zrodloPodagenta = ` FROM podagent p
	                    JOIN okno_komunikacji o ON o.id = p.okno_wykonawcy_id
	                    JOIN sesja s ON s.id = o.sesja_id`

	wstawPodagenta = `INSERT INTO podagent
	                  (identyfikator_zewnetrzny, okno_wykonawcy_id, bieg_id,
	                   pozycja_kolejki_id, nazwa, zadanie, stan)
	                  VALUES (?, ?, ?, ?, ?, ?, ?)`

	pobierzPodagentaPoID = `SELECT ` + kolumnyPodagenta + zrodloPodagenta + ` WHERE p.id = ?`

	// Zawężenia idą warunkiem „puste znaczy wszystko", więc jedno zapytanie
	// obsługuje cztery kształty pytania panelu i nie ma czterech zapytań, które
	// mogłyby się rozjechać.
	listaPodagentow = `SELECT ` + kolumnyPodagenta + zrodloPodagenta +
		` WHERE (? = '' OR o.identyfikator_zewnetrzny = ?)
		    AND (? = '' OR s.identyfikator_zewnetrzny = ?)
		    AND (? = '' OR p.stan = ?)
		  ORDER BY p.id`

	przypiszPozycjePodagenta = `UPDATE podagent SET pozycja_kolejki_id = ?
	                            WHERE identyfikator_zewnetrzny = ?`

	// Wynik pisze się po pozycji, bo w chwili zapisu znana jest pozycja, nie
	// podagent. Indeks częściowy `idx_podagent_pozycja` czyni
	// odczyt punktowym, a UNIQUE gwarantuje najwyżej jeden wiersz.
	zapiszWynikPoPozycji = `UPDATE podagent SET wynik = ?
	                        WHERE pozycja_kolejki_id = ?`

	// Znaczniki czasu stawia baza, nie rdzeń: chwila rozpoczęcia zapisuje się
	// raz, przy pierwszym wejściu w bieg, a chwila zakończenia raz, przy
	// pierwszym stanie końcowym. Dzięki temu powtórzony zapis stanu nie
	// przesuwa historii podagenta.
	// Powód zakończenia i oznaka życia dokładają się tym samym poleceniem:
	// stan mówi co, powód mówi dlaczego, a oznaka — kiedy rdzeń
	// ostatni raz tego wiersza dotknął. Powód dla stanu niekońcowego zostaje
	// zastany, bo przejście 'pending'→'running' niczego nie kończy.
	ustawStanPodagenta = `UPDATE podagent
	                         SET stan = ?,
	                             wynik = COALESCE(?, wynik),
	                             oznaka_zycia = strftime('%Y-%m-%dT%H:%M:%fZ','now'),
	                             powod_zakonczenia = CASE ?
	                                 WHEN 'done'    THEN 'ukonczony'
	                                 WHEN 'failed'  THEN 'blad'
	                                 WHEN 'stopped' THEN 'zatrzymany'
	                                 ELSE powod_zakonczenia END,
	                             rozpoczeto = CASE
	                                 WHEN ? = 'running' AND rozpoczeto IS NULL
	                                 THEN strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                                 ELSE rozpoczeto END,
	                             zakonczono = CASE
	                                 WHEN ? IN ('done','failed','stopped')
	                                 THEN COALESCE(zakonczono, strftime('%Y-%m-%dT%H:%M:%fZ','now'))
	                                 ELSE zakonczono END
	                       WHERE identyfikator_zewnetrzny = ?`
)

type repozytoriumPodagentow struct {
	zapytania *zapytania
	db        *sql.DB
}

// Zgodność implementacji z kontraktem sprawdza kompilator, a nie dopiero montaż.
var _ RepozytoriumPodagentow = (*repozytoriumPodagentow)(nil)

// Podagenci oddaje repozytorium podagentów nad pamięcią zapytań zestawu.
func (z *Zestaw) Podagenci() RepozytoriumPodagentow {
	if z == nil || z.zapytania == nil {
		return nil
	}
	return &repozytoriumPodagentow{zapytania: z.zapytania, db: z.zapytania.db}
}

// ZalozPodagentow zakłada komplet podagentów jednego powołania i oddaje je
// w stanie po zapisie.
//
// Jedna transakcja na całe powołanie. `subagent.spawn` powołuje od jednego do
// piętnastu podagentów jednym żądaniem; zapis wierszami osobnymi zostawiałby po
// awarii połowę powołania, czyli podagentów bez reszty ich pracy.
//
// Wykaz pusty nie jest błędem — zwraca wykaz pusty.
func (r *repozytoriumPodagentow) ZalozPodagentow(ctx context.Context,
	podagenci []Podagent) ([]Podagent, error) {

	numery := make([]int64, 0, len(podagenci))
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawPodagenta)
		if err != nil {
			return err
		}
		for _, podagent := range podagenci {
			if podagent.Kod == "" || podagent.OknoID == 0 || podagent.Zadanie == "" {
				return fmt.Errorf("dane: podagent bez identyfikatora, okna albo zadania")
			}
			stan := podagent.Stan
			if stan == "" {
				stan = StanPodagentaOczekuje
			}
			wynik, err := polecenie.ExecContext(ctx, podagent.Kod, podagent.OknoID,
				liczbaDoKolumny(podagent.BiegID), liczbaDoKolumny(podagent.PozycjaKolejkiID),
				tekstDoKolumny(podagent.Nazwa), podagent.Zadanie, stan)
			if err != nil {
				return fmt.Errorf("dane: nie można założyć podagenta %q: %w", podagent.Kod, err)
			}
			id, err := wynik.LastInsertId()
			if err != nil {
				return fmt.Errorf("dane: nieznany numer założonego podagenta: %w", err)
			}
			numery = append(numery, id)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	zalozeni := make([]Podagent, 0, len(numery))
	for _, id := range numery {
		podagent, err := r.podagentPoID(ctx, id)
		if err != nil {
			return nil, err
		}
		zalozeni = append(zalozeni, podagent)
	}
	return zalozeni, nil
}

// Podagenci zwraca wykaz zawężony filtrem, w kolejności powołania.
func (r *repozytoriumPodagentow) Podagenci(ctx context.Context,
	filtr FiltrPodagentow) ([]Podagent, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaPodagentow)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, filtr.OknoKod, filtr.OknoKod,
		filtr.SesjaKod, filtr.SesjaKod, filtr.Stan, filtr.Stan)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wykazu podagentów: %w", err)
	}
	return zbierzPodagentow(wiersze)
}

// PodagenciPoKodach zwraca podagentów wskazanych wprost — tak pyta
// `subagent.result.collect`, gdy Operator zbiera wyniki wybranych.
//
// Liczba miejsc w zapytaniu rośnie z wykazem, a pamięć poleceń trzyma po jednym
// wariancie na długość. Wariantów jest najwyżej tyle, ilu podagentów da się
// powołać, więc pamięć nie puchnie.
func (r *repozytoriumPodagentow) PodagenciPoKodach(ctx context.Context,
	kody []string) ([]Podagent, error) {

	if len(kody) == 0 {
		return []Podagent{}, nil
	}
	zapytanie := `SELECT ` + kolumnyPodagenta + zrodloPodagenta +
		` WHERE p.identyfikator_zewnetrzny IN (` +
		strings.TrimSuffix(strings.Repeat("?,", len(kody)), ",") + `) ORDER BY p.id`
	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return nil, err
	}
	klucze := make([]any, 0, len(kody))
	for _, kod := range kody {
		klucze = append(klucze, kod)
	}
	wiersze, err := polecenie.QueryContext(ctx, klucze...)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wskazanych podagentów: %w", err)
	}
	return zbierzPodagentow(wiersze)
}

// PrzypiszPozycje wiąże podagenta z pozycją kolejki, która niesie jego pracę.
// Wiązanie powstaje po powołaniu, bo kolejka zakłada się dopiero wtedy — kolumna
// jest pusta wyłącznie w tym oknie czasu.
func (r *repozytoriumPodagentow) PrzypiszPozycje(ctx context.Context,
	kod string, pozycjaID int64) error {

	if kod == "" || pozycjaID == 0 {
		return fmt.Errorf("dane: wiązanie podagenta z pozycją bez podagenta albo bez pozycji")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, przypiszPozycjePodagenta)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, pozycjaID, kod); err != nil {
		return fmt.Errorf("dane: nie można związać podagenta %q z pozycją %d: %w",
			kod, pozycjaID, err)
	}
	return nil
}

// UstawStan zapisuje stan podagenta, a przy stanie końcowym także jego wynik.
// Wynik pusty zostawia zastany: przejście przez stany nie ma prawa wymazać
// pracy, którą podagent już oddał.
func (r *repozytoriumPodagentow) UstawStan(ctx context.Context,
	kod, stan string, wynik *string) error {

	if kod == "" || stan == "" {
		return fmt.Errorf("dane: zapis stanu podagenta bez podagenta albo bez stanu")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, ustawStanPodagenta)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, stan, tekstDoKolumny(wynik), stan, stan, stan, kod); err != nil {
		return fmt.Errorf("dane: nie można zapisać stanu podagenta %q: %w", kod, err)
	}
	return nil
}

// ZapiszWynikPozycji utrwala treść tury na wierszu podagenta związanym
// z pozycją. Zapis bez trafienia (pozycja spoza podagentów) przechodzi bez
// błędu — patrz kontrakt interfejsu wyżej.
func (r *repozytoriumPodagentow) ZapiszWynikPozycji(ctx context.Context,
	pozycjaID int64, wynik string) error {

	if pozycjaID == 0 {
		return fmt.Errorf("dane: zapis wyniku bez wskazania pozycji kolejki")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszWynikPoPozycji)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, wynik, pozycjaID); err != nil {
		return fmt.Errorf("dane: nie można zapisać wyniku pozycji %d: %w", pozycjaID, err)
	}
	return nil
}

// podagentPoID odczytuje pojedynczy wiersz po kluczu głównym.
func (r *repozytoriumPodagentow) podagentPoID(ctx context.Context, id int64) (Podagent, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzPodagentaPoID)
	if err != nil {
		return Podagent{}, err
	}
	podagent, err := odczytajPodagenta(polecenie.QueryRowContext(ctx, id))
	if errors.Is(err, sql.ErrNoRows) {
		return Podagent{}, ErrBrakWiersza
	}
	if err != nil {
		return Podagent{}, fmt.Errorf("dane: nieczytelny podagent %d: %w", id, err)
	}
	return podagent, nil
}

// zbierzPodagentow składa wykaz z otwartego wyniku zapytania.
func zbierzPodagentow(wiersze *sql.Rows) ([]Podagent, error) {
	defer wiersze.Close()

	lista := []Podagent{}
	for wiersze.Next() {
		podagent, err := odczytajPodagenta(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz podagenta: %w", err)
		}
		lista = append(lista, podagent)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt podagentów: %w", err)
	}
	return lista, nil
}

// odczytajPodagenta składa strukturę z jednego wiersza wyniku.
func odczytajPodagenta(wiersz skaner) (Podagent, error) {
	var p Podagent
	var oknoKod, sesjaKod, nazwa, wynik, rozpoczeto, zakonczono sql.NullString
	var bieg, pozycja sql.NullInt64
	err := wiersz.Scan(&p.ID, &p.Kod, &p.OknoID, &oknoKod, &sesjaKod, &bieg, &pozycja,
		&nazwa, &p.Zadanie, &p.Stan, &wynik, &p.Zetonow, &p.Narzedzi,
		&rozpoczeto, &zakonczono, &p.Utworzono)
	if err != nil {
		return Podagent{}, err
	}
	p.OknoKod, p.SesjaKod = tekstZKolumny(oknoKod), tekstZKolumny(sesjaKod)
	p.BiegID, p.PozycjaKolejkiID = liczbaZKolumny(bieg), liczbaZKolumny(pozycja)
	p.Nazwa, p.Wynik = tekstZKolumny(nazwa), tekstZKolumny(wynik)
	p.Rozpoczeto, p.Zakonczono = tekstZKolumny(rozpoczeto), tekstZKolumny(zakonczono)
	return p, nil
}
