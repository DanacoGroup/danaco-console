// Odpowiedzialność pliku: historia rozmowy okna widziana jako wykaz pozycji —
// odczyt stronicowany kursorem i usunięcie (rodzina `history.*`). Zasada
// przechowywania tego wykazu żyje obok, w `historia_retencja.go`: to nastawa,
// która kasuje sama, i osobna odpowiedzialność.
//
// Repozytorium nie zakłada własnej tabeli historii. Pozycją historii jest
// wiersz `wiadomosc` — tylko on niesie rolę, treść i czas, więc tylko z niego
// złoży się `HistoryEntry`. Bloki są wnętrzem wypowiedzi, nie pozycją.
//
// Osobne repozytorium obok `RepozytoriumWiadomosci` bierze się z drugiego
// pytania o tę samą tabelę. Tamto prowadzi turę: pisze wypowiedź, czyta okno
// po kluczu wewnętrznym i nie usuwa nigdy. Historia pyta po identyfikatorze
// kontraktowym okna, od najnowszej, kursorem czasu — i jako jedyna kasuje;
// usunięcie zabiera też bloki, bo bez klucza obcego kaskada ich nie sprząta.
package dane

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"danacoconsole/shared"
)

// PozycjaHistorii to jedna pozycja wykazu historii okna — wiersz `wiadomosc`
// w kształcie, w jakim pyta o niego rodzina `history.*`.
type PozycjaHistorii struct {
	// Identyfikator kontraktowy: identyfikator zewnętrzny wiadomości, a bez
	// niego klucz wiersza jako napis — bez tożsamości nie da się usunąć.
	Identyfikator string
	OknoKod       string
	// SesjaKod bywa pusty: sesja spoza rdzenia nie ma identyfikatora
	// zewnętrznego. Brak nie jest błędem.
	SesjaKod string
	Rola     shared.MessageRole
	// Tresc jest pełna — skrót wykazu składa warstwa wyższa (zna próg kontraktu).
	Tresc  string
	Chwila int64 // milisekundy epoki

	// klucz i znacznik są wewnętrzne: klucz wiersza i surowy napis kolumny
	// czasu. Kontrakt ich nie zna (kursor niesie same milisekundy), ale
	// domknięcie strony po rodzeństwie o tym samym znaczniku potrzebuje
	// jednego i drugiego — patrz `Pozycje`.
	klucz    int64
	znacznik string
}

// RepozytoriumHistorii jest kontraktem obszaru historii i retencji.
type RepozytoriumHistorii interface {
	// Pozycje zwraca wykaz od najnowszej. Kursor `przed` (milisekundy epoki,
	// 0 znaczy „bez kursora") stronicuje wstecz; limit 0 znaczy całość.
	Pozycje(ctx context.Context, oknoKod string, przed int64, limit int) ([]PozycjaHistorii, error)
	// Policz liczy wszystkie pozycje historii okna — bez kursora i bez limitu
	// (pole `total`). Patrz uzasadnienie przy metodzie.
	Policz(ctx context.Context, oknoKod string) (int, error)
	// Usun kasuje wskazane pozycje; wykaz pusty czyści całą historię okna.
	Usun(ctx context.Context, oknoKod string, identyfikatory []string) (int, error)
	// ZapiszZasade zakłada zasadę zakresu albo nadpisuje istniejącą.
	ZapiszZasade(ctx context.Context, zasada ZasadaPrzechowywania) (ZasadaPrzechowywania, error)
	// ZasadaOkna rozstrzyga zasadę okna: okno przed sesją, sesja przed globalną.
	// Fałsz znaczy „zasady nie ma" — stan poprawny.
	ZasadaOkna(ctx context.Context, oknoKod string) (ZasadaPrzechowywania, bool, error)
	// IstniejeByt mówi, czy okno albo sesja wskazana przez zasadę zakresu jest
	// w bazie; zakres globalny bytu nie wskazuje i zawsze jest prawdziwy.
	IstniejeByt(ctx context.Context, zakres, zakresKod string) (bool, error)
	// OknaZakresu wylicza okna objęte zasadą zakresu — tędy idzie egzekucja.
	OknaZakresu(ctx context.Context, zakres, zakresKod string) ([]string, error)
	// Egzekwuj stosuje zasadę okna i zwraca liczbę pozycji, które odpadły.
	Egzekwuj(ctx context.Context, oknoKod string) (int, error)
}

const (
	// Kursor podany dwukrotnie tym samym parametrem: `?=''` przepuszcza wykaz
	// cały, bo znacznik pusty nie jest datą.
	warunekHistorii = ` FROM wiadomosc w
	                    JOIN okno_komunikacji o ON o.id = w.okno_komunikacji_id
	                    JOIN sesja s ON s.id = o.sesja_id
	                    WHERE o.identyfikator_zewnetrzny = ?
	                      AND (? = '' OR w.utworzono < ?)`
	pozycjeHistorii = `SELECT w.id, COALESCE(w.identyfikator_zewnetrzny, ''),
	                          COALESCE(s.identyfikator_zewnetrzny, ''),
	                          w.rola, COALESCE(w.tresc, ''), w.utworzono` +
		warunekHistorii + `
	                   ORDER BY w.utworzono DESC, w.id DESC
	                   LIMIT CASE WHEN ? > 0 THEN ? ELSE -1 END`
	// Liczba pozycji całego okna — bez warunku kursora. Kursor zaniżałby ją przy
	// każdym „Wczytaj starsze", a panel podaje ją Operatorowi jako „ile zniknie"
	// przy czyszczeniu historii, które kasuje całe okno (`history.delete` bez
	// wskazania pozycji). Liczba zaniżona przy ostrzeżeniu o czynności
	// nieodwracalnej jest gorsza niż jej brak.
	liczbaPozycjiHistorii = `SELECT COUNT(*) FROM wiadomosc w
	                         JOIN okno_komunikacji o ON o.id = w.okno_komunikacji_id
	                         WHERE o.identyfikator_zewnetrzny = ?`

	// Rodzeństwo ostatniej pozycji strony: wiersze o tym samym znaczniku czasu,
	// stojące w porządku wykazu za nią. Domykają stronę, żeby kursor czasu nie
	// przeciął grupy o jednej milisekundzie — patrz `Pozycje`.
	rodzenstwoHistorii = `SELECT w.id, COALESCE(w.identyfikator_zewnetrzny, ''),
	                             COALESCE(s.identyfikator_zewnetrzny, ''),
	                             w.rola, COALESCE(w.tresc, ''), w.utworzono
	                      FROM wiadomosc w
	                      JOIN okno_komunikacji o ON o.id = w.okno_komunikacji_id
	                      JOIN sesja s ON s.id = o.sesja_id
	                      WHERE o.identyfikator_zewnetrzny = ?
	                        AND w.utworzono = ?
	                        AND w.id < ?
	                      ORDER BY w.id DESC`

	// Pozycję wskazuje identyfikator zewnętrzny, a wiersz bez niego — klucz
	// wiersza w postaci napisu. Ta sama tożsamość, którą oddaje odczyt.
	usunPozycjeHistorii = `DELETE FROM wiadomosc
	                       WHERE okno_komunikacji_id = (SELECT id FROM okno_komunikacji
	                                                    WHERE identyfikator_zewnetrzny = ?)
	                         AND (identyfikator_zewnetrzny = ?
	                              OR (identyfikator_zewnetrzny IS NULL AND CAST(id AS TEXT) = ?))`
	usunHistorieOkna = `DELETE FROM wiadomosc
	                    WHERE okno_komunikacji_id = (SELECT id FROM okno_komunikacji
	                                                 WHERE identyfikator_zewnetrzny = ?)`

	// Bloki osierocone — wiszą na identyfikatorze wiadomości, której już nie ma
	// (bez klucza obcego nie ma kaskady).
	usunBlokiOsierocone = `DELETE FROM blok_wiadomosci
	                       WHERE okno_kod = ?
	                         AND wiadomosc_kod NOT IN (
	                             SELECT COALESCE(w.identyfikator_zewnetrzny, '')
	                             FROM wiadomosc w
	                             JOIN okno_komunikacji o ON o.id = w.okno_komunikacji_id
	                             WHERE o.identyfikator_zewnetrzny = ?)`
)

type repozytoriumHistorii struct {
	zapytania *zapytania
	db        *sql.DB
}

func noweRepozytoriumHistorii(z *zapytania, db *sql.DB) *repozytoriumHistorii {
	return &repozytoriumHistorii{zapytania: z, db: db}
}

// Pozycje zwraca wykaz historii okna od najnowszej. Okno nieznane albo puste
// daje wykaz pusty — czytanie nie odmawia z powodu pustki.
//
// Strona nie przecina grupy o jednym znaczniku. Porządek wykazu ma dwa klucze
// (`utworzono DESC, id DESC`), a kursor kontraktu niesie tylko pierwszy z nich
// — milisekundy. Gdy granica strony wypada w środku pozycji o tym samym
// znaczniku (a tak wpada zwykła tura: pytanie i odpowiedź powstają w tej samej
// milisekundzie), strona następna pytana warunkiem ostro mniejszym
// `utworzono < kursor` przeskoczyłaby rodzeństwo bezpowrotnie. Zamiast dokładać
// kontraktowi pola rozstrzygającego remis, domykamy stronę tutaj: po pobraniu
// `limit` wierszy dobieramy jeszcze całe rodzeństwo ostatniego z nich. Strona
// bywa więc odrobinę dłuższa od limitu, za to kursor zawsze pada między grupami
// i nic nie ginie.
func (r *repozytoriumHistorii) Pozycje(ctx context.Context, oknoKod string,
	przed int64, limit int) ([]PozycjaHistorii, error) {

	lista, err := r.stronaPozycji(ctx, oknoKod, przed, limit)
	if err != nil || limit <= 0 || len(lista) < limit {
		return lista, err
	}
	rodzenstwo, err := r.rodzenstwoOstatniej(ctx, oknoKod, lista[len(lista)-1])
	if err != nil {
		return nil, err
	}
	return append(lista, rodzenstwo...), nil
}

// stronaPozycji pobiera samą stronę wykazu — bez domykania rodzeństwa.
func (r *repozytoriumHistorii) stronaPozycji(ctx context.Context, oknoKod string,
	przed int64, limit int) ([]PozycjaHistorii, error) {

	kursor := znacznikKursora(przed)
	polecenie, err := r.zapytania.przygotuj(ctx, pozycjeHistorii)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, oknoKod, kursor, kursor, limit, limit)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać historii okna %q: %w", oknoKod, err)
	}
	defer wiersze.Close()

	lista := []PozycjaHistorii{}
	for wiersze.Next() {
		pozycja, err := odczytajPozycjeHistorii(wiersze, oknoKod)
		if err != nil {
			return nil, err
		}
		lista = append(lista, pozycja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt historii okna %q: %w", oknoKod, err)
	}
	return lista, nil
}

// rodzenstwoOstatniej dobiera pozycje o tym samym znaczniku czasu, co ostatnia
// pozycja strony, stojące za nią w porządku wykazu. Bez nich kursor czasu
// przeciąłby grupę o jednej milisekundzie (patrz `Pozycje`).
func (r *repozytoriumHistorii) rodzenstwoOstatniej(ctx context.Context, oknoKod string,
	ostatnia PozycjaHistorii) ([]PozycjaHistorii, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, rodzenstwoHistorii)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, oknoKod, ostatnia.znacznik, ostatnia.klucz)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można domknąć strony historii okna %q: %w", oknoKod, err)
	}
	defer wiersze.Close()

	lista := []PozycjaHistorii{}
	for wiersze.Next() {
		pozycja, err := odczytajPozycjeHistorii(wiersze, oknoKod)
		if err != nil {
			return nil, err
		}
		lista = append(lista, pozycja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwane domykanie strony historii okna %q: %w", oknoKod, err)
	}
	return lista, nil
}

// Policz liczy wszystkie pozycje historii okna. Kursora nie bierze pod uwagę:
// liczba idzie do pola `total`, a warstwa wyższa podaje ją Operatorowi jako
// rozmiar całej historii okna — także w ostrzeżeniu przed jej wyczyszczeniem,
// które kasuje okno w całości.
func (r *repozytoriumHistorii) Policz(ctx context.Context, oknoKod string) (int, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, liczbaPozycjiHistorii)
	if err != nil {
		return 0, err
	}
	liczba := 0
	if err := polecenie.QueryRowContext(ctx, oknoKod).Scan(&liczba); err != nil {
		return 0, fmt.Errorf("dane: nie można policzyć historii okna %q: %w", oknoKod, err)
	}
	return liczba, nil
}

// Usun kasuje wskazane pozycje, a przy wykazie pustym — całą historię okna.
// Jedna transakcja wraz z blokami: zapis nie zostaje połowiczny.
func (r *repozytoriumHistorii) Usun(ctx context.Context, oknoKod string,
	identyfikatory []string) (int, error) {

	usuniete := 0
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		if len(identyfikatory) == 0 {
			liczba, err := wykonajUsuniecie(ctx, transakcja, usunHistorieOkna, oknoKod)
			if err != nil {
				return err
			}
			usuniete = liczba
			return sprzatnijBloki(ctx, transakcja, oknoKod)
		}
		for _, identyfikator := range identyfikatory {
			liczba, err := wykonajUsuniecie(ctx, transakcja, usunPozycjeHistorii,
				oknoKod, identyfikator, identyfikator)
			if err != nil {
				return err
			}
			usuniete += liczba
		}
		return sprzatnijBloki(ctx, transakcja, oknoKod)
	})
	if err != nil {
		return 0, err
	}
	return usuniete, nil
}

// znacznikCzasuHistorii jest formatem kolumn czasu schematu.
// Porównanie kursora idzie po napisie, więc format musi być dokładnie ten sam,
// co w `strftime` bazy.
const znacznikCzasuHistorii = "2006-01-02T15:04:05.000Z"

// znacznikKursora przekłada kursor kontraktu (milisekundy epoki) na znacznik
// kolumny. Zero znaczy „bez kursora" i daje napis pusty, który warunek
// zapytania przepuszcza — pierwsza strona nie ma czasu odniesienia.
func znacznikKursora(przed int64) string {
	if przed <= 0 {
		return ""
	}
	return time.UnixMilli(przed).UTC().Format(znacznikCzasuHistorii)
}

// odczytajPozycjeHistorii składa pozycję z jednego wiersza wyniku.
func odczytajPozycjeHistorii(wiersz skaner, oknoKod string) (PozycjaHistorii, error) {
	var pozycja PozycjaHistorii
	var klucz int64
	var rola, znacznik string
	err := wiersz.Scan(&klucz, &pozycja.Identyfikator, &pozycja.SesjaKod,
		&rola, &pozycja.Tresc, &znacznik)
	if err != nil {
		return PozycjaHistorii{}, fmt.Errorf("dane: nieczytelny wiersz historii okna: %w", err)
	}
	pozycja.klucz, pozycja.znacznik = klucz, znacznik
	if pozycja.Identyfikator == "" {
		pozycja.Identyfikator = strconv.FormatInt(klucz, 10)
	}
	pozycja.OknoKod = oknoKod
	pozycja.Rola, err = rolaWiadomosciZBazy(rola)
	if err != nil {
		return PozycjaHistorii{}, err
	}
	if chwila, blad := time.Parse(time.RFC3339Nano, znacznik); blad == nil {
		pozycja.Chwila = chwila.UnixMilli()
	}
	return pozycja, nil
}

// wykonajUsuniecie wykonuje polecenie kasujące i oddaje liczbę wierszy.
func wykonajUsuniecie(ctx context.Context, transakcja *sql.Tx,
	polecenie string, argumenty ...any) (int, error) {

	wynik, err := transakcja.ExecContext(ctx, polecenie, argumenty...)
	if err != nil {
		return 0, fmt.Errorf("dane: czyszczenie historii okna: %w", err)
	}
	liczba, err := wynik.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("dane: nieznana liczba usuniętych pozycji historii: %w", err)
	}
	return int(liczba), nil
}

// sprzatnijBloki usuwa bloki wiadomości, których wypowiedzi już nie ma.
func sprzatnijBloki(ctx context.Context, transakcja *sql.Tx, oknoKod string) error {
	if _, err := transakcja.ExecContext(ctx, usunBlokiOsierocone, oknoKod, oknoKod); err != nil {
		return fmt.Errorf("dane: czyszczenie bloków okna %q: %w", oknoKod, err)
	}
	return nil
}
