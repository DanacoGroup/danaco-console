// Odpowiedzialność pliku: trwałość powiadomień — rejestracja urządzenia do
// wołania, kolejka powiadomień i ślad doręczenia. Trzy tabele obsługuje jedno
// repozytorium, bo opisują trzy strony jednej sprawy: kogo wołać, czym i czy
// doszło.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// Kanały wołania. Słownik jest zamknięty na to, co rdzeń potrafi doręczyć, i ma
// pokrywać się z warunkiem CHECK na kolumnie `kanal`. Dopisanie tu kanału usługi
// zewnętrznej bez kodu, który go obsłuży, przepuści rejestrację, której żaden
// takt nie doręczy.
const (
	// KanalPolaczenie doręcza po gnieździe, które Operator już otworzył; to
	// jedyny kanał ujęty w słowniku, więc rejestracja bez wskazania kanału
	// przyjmuje go jako wartość domyślną.
	KanalPolaczenie = "polaczenie"
)

// Stany powiadomienia. Wartości muszą być te same, co w warunku CHECK na
// kolumnie `stan` — inaczej zapis odbije się o schemat.
const (
	StanPowiadomieniaOczekuje    = "oczekuje"
	StanPowiadomieniaDostarczone = "dostarczone"
	StanPowiadomieniaOdwolane    = "odwolane"
	StanPowiadomieniaWygasle     = "wygasle"
	StanPowiadomieniaPorzucone   = "porzucone"
)

// Priorytety powiadomienia. Kolejność doręczeń w zapytaniu należnych stawia
// powiadomienia pilne przed zwykłymi, niezależnie od kolejności wpisu do
// kolejki.
const (
	PriorytetZwykly = "zwykly"
	PriorytetPilny  = "pilny"
)

// RejestracjaPowiadomien to wiersz tabeli `urzadzenie_powiadomien`: zgoda
// jednego urządzenia na wołanie jedną drogą. `KluczKanalu` dla kanału
// `polaczenie` jest adresem gniazda, a nie odrębną tożsamością urządzenia.
type RejestracjaPowiadomien struct {
	ID                  int64
	UrzadzenieID        int64
	Kanal               string
	KluczKanalu         string
	Etykieta            *string
	Aktywne             bool
	Zarejestrowano      string
	Wyrejestrowano      *string
	OstatnioDostarczono *string
}

// Powiadomienie to wiersz kolejki.
//
// `BytRodzaj` i `BytID` mówią, czego powiadomienie dotyczy, i to po nich idzie
// odwołanie. Kotwica jest miękka, bez klucza obcego, bo jedno powiadomienie ma
// z założenia móc dotyczyć bytów z różnych tabel.
type Powiadomienie struct {
	ID             int64
	Tytul          string
	Tresc          string
	Priorytet      string
	Powod          string
	BytRodzaj      *string
	BytID          *string
	Stan           string
	ProbIle        int
	NastepnaProba  string
	Wygasa         string
	Utworzono      string
	Dostarczono    *string
	Odwolano       *string
	PowodOdwolania *string
}

// PlanTaktu jest tym, co silnik podaje repozytorium na jeden takt jednego
// powiadomienia. Wysyłka jest domknięciem, a nie wynikiem, żeby sprawdzenie
// stanu wiersza i sama wysyłka zmieściły się w jednej transakcji zapisu.
type PlanTaktu struct {
	// Teraz jest chwilą taktu w formacie znacznika bazy.
	Teraz string
	// Wyslij dostaje wiersz pod zamkiem i zwraca klucze rejestracji, które
	// kopertę przyjęły.
	Wyslij func(ctx context.Context, p Powiadomienie) ([]int64, error)
	// NastepnaProba wylicza termin kolejnego podejścia z liczby prób już
	// odbytych.
	NastepnaProba func(probIle int) string
}

// WynikTaktu opisuje, co takt naprawdę zrobił z jednym powiadomieniem: czy je
// pominięto, w jakim stanie zostawił wiersz i ilu odbiorcom go doręczono.
type WynikTaktu struct {
	// Pominieto znaczy, że wiersz nie był już oczekuje w chwili zamknięcia
	// zamka; to poprawny wynik.
	Pominieto bool
	// Stan wiersza po takcie.
	Stan string
	// Odbiorcow to liczba urządzeń, które kopertę przyjęły.
	Odbiorcow int
}

// RepozytoriumPowiadomien jest kontraktem trwałości powiadomień: obejmuje
// rejestrację urządzeń do wołania, kolejkę powiadomień i ślad ich doręczenia.
type RepozytoriumPowiadomien interface {
	// Zarejestruj zapisuje zgodę urządzenia na wołanie wskazaną drogą;
	// powtórzona odświeża wiersz.
	Zarejestruj(ctx context.Context, r RejestracjaPowiadomien) (int64, error)
	// Wyrejestruj cofa zgodę, zostawiając wiersz, żeby ślad wołania pozostał
	// czytelny.
	Wyrejestruj(ctx context.Context, kanal, kluczKanalu string, teraz string) error
	// AktywneRejestracje zwraca komplet czynnych zgód.
	AktywneRejestracje(ctx context.Context) ([]RejestracjaPowiadomien, error)

	// Wstaw wnosi powiadomienie do kolejki.
	Wstaw(ctx context.Context, p Powiadomienie) (int64, error)
	// Pobierz zwraca wiersz kolejki.
	Pobierz(ctx context.Context, id int64) (Powiadomienie, error)
	// Nalezne zwraca powiadomienia, których termin podejścia już minął.
	Nalezne(ctx context.Context, teraz string) ([]Powiadomienie, error)
	// Takt wykonuje jedno podejście pod zamkiem zapisu.
	Takt(ctx context.Context, id int64, plan PlanTaktu) (WynikTaktu, error)
	// Wygas zamyka powiadomienia, którym minął termin ważności.
	Wygas(ctx context.Context, teraz string) (int, error)
	// Odwolaj gasi oczekujące powiadomienia o wskazanym bycie, gdy decyzja w
	// jego sprawie już zapadła.
	Odwolaj(ctx context.Context, bytRodzaj, bytID, powod, teraz string) (int, error)
	// Potwierdz zapisuje, że Operator widział powiadomienie na wskazanym
	// urządzeniu.
	Potwierdz(ctx context.Context, powiadomienieID int64, kanal, kluczKanalu, teraz string) error
	// Dostarczenia zwraca klucze rejestracji, którym powiadomienie doręczono.
	Dostarczenia(ctx context.Context, powiadomienieID int64) ([]int64, error)
}

const kolumnyPowiadomienia = `id, tytul, tresc, priorytet, powod, byt_rodzaj, byt_id, stan,
                              prob_ile, nastepna_proba, wygasa, utworzono, dostarczono,
                              odwolano, powod_odwolania`

const kolumnyRejestracjiPowiadomien = `id, urzadzenie_id, kanal, klucz_kanalu, etykieta,
                                       aktywne, zarejestrowano, wyrejestrowano,
                                       ostatnio_dostarczono`

type repozytoriumPowiadomien struct {
	zapytania *zapytania
	db        *sql.DB
}

// Zgodność implementacji z kontraktem sprawdzana jest przy kompilacji:
// przypisanie repozytoriumPowiadomien do zmiennej typu RepozytoriumPowiadomien
// nie skompiluje się, gdy metody przestaną się pokrywać.
var _ RepozytoriumPowiadomien = (*repozytoriumPowiadomien)(nil)

// NowePowiadomienia zakłada repozytorium powiadomień nad otwartą pulą połączeń.
// Konstruktor bierze uchwyt puli, a nie zestaw repozytoriów, żeby silnik
// powiadomień mógł złożyć repozytorium z uchwytu, który jego pakiet już ma.
func NowePowiadomienia(db *sql.DB) RepozytoriumPowiadomien {
	if db == nil {
		return nil
	}
	return &repozytoriumPowiadomien{zapytania: noweZapytania(db), db: db}
}

// ── REJESTRACJA ─────────────────────────────────────────────────────────────

const zapiszRejestracjePowiadomien = `
	INSERT INTO urzadzenie_powiadomien (urzadzenie_id, kanal, klucz_kanalu, etykieta, aktywne, konto_id)
	VALUES (?, ?, ?, ?, 1, ` + WskazanieKonta + `)
	ON CONFLICT(urzadzenie_id, kanal, klucz_kanalu) DO UPDATE SET
	    etykieta       = COALESCE(excluded.etykieta, urzadzenie_powiadomien.etykieta),
	    aktywne        = 1,
	    wyrejestrowano = NULL
	WHERE ` + WarunekKonta + `
	RETURNING id`

func (r *repozytoriumPowiadomien) Zarejestruj(ctx context.Context, rej RejestracjaPowiadomien) (int64, error) {
	kanal := rej.Kanal
	if kanal == "" {
		kanal = KanalPolaczenie
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszRejestracjePowiadomien)
	if err != nil {
		return 0, err
	}
	var id int64
	konto := KontoOperatora(ctx)
	err = polecenie.QueryRowContext(ctx, rej.UrzadzenieID, kanal, rej.KluczKanalu,
		rej.Etykieta, konto, konto).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("dane: nie można zapisać rejestracji powiadomień urządzenia %d "+
			"na kanale %q: %w", rej.UrzadzenieID, kanal, err)
	}
	return id, nil
}

const wyrejestrujPowiadomienia = `
	UPDATE urzadzenie_powiadomien
	   SET aktywne = 0, wyrejestrowano = ?
	 WHERE kanal = ? AND klucz_kanalu = ? AND aktywne = 1 AND ` + WarunekKonta

// Wyrejestruj cofa zgodę. Brak wiersza nie jest błędem: cofnięcie zgody, której
// nie było, kończy się tym samym stanem, o który wołający prosił.
func (r *repozytoriumPowiadomien) Wyrejestruj(ctx context.Context, kanal, kluczKanalu, teraz string) error {
	polecenie, err := r.zapytania.przygotuj(ctx, wyrejestrujPowiadomienia)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, teraz, kanal, kluczKanalu,
		KontoOperatora(ctx)); err != nil {
		return fmt.Errorf("dane: nie można wyrejestrować powiadomień dla %q/%q: %w",
			kanal, kluczKanalu, err)
	}
	return nil
}

const listaAktywnychRejestracji = `SELECT ` + kolumnyRejestracjiPowiadomien + `
	FROM urzadzenie_powiadomien WHERE aktywne = 1 ORDER BY id`

func (r *repozytoriumPowiadomien) AktywneRejestracje(ctx context.Context) ([]RejestracjaPowiadomien, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaAktywnychRejestracji)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać czynnych rejestracji powiadomień: %w", err)
	}
	defer wiersze.Close()

	lista := []RejestracjaPowiadomien{}
	for wiersze.Next() {
		rej, err := odczytajRejestracjePowiadomien(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, rej)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt rejestracji powiadomień: %w", err)
	}
	return lista, nil
}

func odczytajRejestracjePowiadomien(wiersz skaner) (RejestracjaPowiadomien, error) {
	var rej RejestracjaPowiadomien
	var etykieta, wyrejestrowano, dostarczono sql.NullString
	var aktywne int
	err := wiersz.Scan(&rej.ID, &rej.UrzadzenieID, &rej.Kanal, &rej.KluczKanalu, &etykieta,
		&aktywne, &rej.Zarejestrowano, &wyrejestrowano, &dostarczono)
	if err != nil {
		return RejestracjaPowiadomien{}, fmt.Errorf("dane: nieczytelny wiersz rejestracji "+
			"powiadomień: %w", err)
	}
	rej.Etykieta = tekstZKolumny(etykieta)
	rej.Aktywne = aktywne != 0
	rej.Wyrejestrowano = tekstZKolumny(wyrejestrowano)
	rej.OstatnioDostarczono = tekstZKolumny(dostarczono)
	return rej, nil
}

// ── KOLEJKA ─────────────────────────────────────────────────────────────────

const wstawPowiadomienie = `
	INSERT INTO powiadomienie (tytul, tresc, priorytet, powod, byt_rodzaj, byt_id,
	                           nastepna_proba, wygasa, konto_id)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ` + WskazanieKonta + `)
	RETURNING id`

func (r *repozytoriumPowiadomien) Wstaw(ctx context.Context, p Powiadomienie) (int64, error) {
	priorytet := p.Priorytet
	if priorytet == "" {
		priorytet = PriorytetZwykly
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawPowiadomienie)
	if err != nil {
		return 0, err
	}
	var id int64
	err = polecenie.QueryRowContext(ctx, p.Tytul, p.Tresc, priorytet, p.Powod,
		p.BytRodzaj, p.BytID, p.NastepnaProba, p.Wygasa, KontoOperatora(ctx)).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("dane: nie można wnieść powiadomienia %q do kolejki: %w", p.Tytul, err)
	}
	return id, nil
}

const pobierzPowiadomienie = `SELECT ` + kolumnyPowiadomienia + ` FROM powiadomienie
	WHERE id = ? AND ` + WarunekKonta

func (r *repozytoriumPowiadomien) Pobierz(ctx context.Context, id int64) (Powiadomienie, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzPowiadomienie)
	if err != nil {
		return Powiadomienie{}, err
	}
	p, err := odczytajPowiadomienie(polecenie.QueryRowContext(ctx, id, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return Powiadomienie{}, fmt.Errorf("dane: powiadomienie %d nie istnieje: %w", id, ErrBrakWiersza)
	}
	if err != nil {
		return Powiadomienie{}, err
	}
	return p, nil
}

// Nalezne czyta powiadomienia, którym termin podejścia już minął. Porządek jest
// rozmyślny: pilne przed zwykłymi, a w obrębie stopnia — te, które czekają
// najdłużej.
const nalezniePowiadomienia = `SELECT ` + kolumnyPowiadomienia + `
	  FROM powiadomienie
	 WHERE stan = 'oczekuje' AND nastepna_proba <= ?
	 ORDER BY CASE priorytet WHEN 'pilny' THEN 0 ELSE 1 END, nastepna_proba, id`

func (r *repozytoriumPowiadomien) Nalezne(ctx context.Context, teraz string) ([]Powiadomienie, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, nalezniePowiadomienia)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, teraz)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać należnych powiadomień: %w", err)
	}
	defer wiersze.Close()

	lista := []Powiadomienie{}
	for wiersze.Next() {
		p, err := odczytajPowiadomienie(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, p)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt należnych powiadomień: %w", err)
	}
	return lista, nil
}

// zajmijPowiadomienie w jednym poleceniu bierze zamek zapisu na wierszu,
// sprawdza, że wiersz nadal czeka, i zlicza podejście. Rozbicie tego na SELECT
// i osobny UPDATE otwiera szczelinę, w której wiersz zmienia stan między
// sprawdzeniem a zajęciem.
const zajmijPowiadomienie = `
	UPDATE powiadomienie SET prob_ile = prob_ile + 1
	 WHERE id = ? AND stan = 'oczekuje'`

const stanPowiadomieniaPodZamkiem = `SELECT ` + kolumnyPowiadomienia + `
	FROM powiadomienie WHERE id = ?`

const zamknijDostarczeniem = `
	UPDATE powiadomienie SET stan = 'dostarczone', dostarczono = ? WHERE id = ?`

const zamknijWygasnieciem = `
	UPDATE powiadomienie SET stan = 'wygasle' WHERE id = ? AND stan = 'oczekuje'`

const przelozProbe = `
	UPDATE powiadomienie SET nastepna_proba = ? WHERE id = ? AND stan = 'oczekuje'`

const zapiszDostarczenie = `
	INSERT INTO powiadomienie_dostarczenie (powiadomienie_id, urzadzenie_powiadomien_id, dostarczono)
	VALUES (?, ?, ?)
	ON CONFLICT(powiadomienie_id, urzadzenie_powiadomien_id) DO NOTHING`

const odnotujDostarczenieUrzadzeniu = `
	UPDATE urzadzenie_powiadomien SET ostatnio_dostarczono = ? WHERE id = ?`

// Takt wykonuje jedno podejście do jednego powiadomienia pod zamkiem zapisu.
// Kolejność zajęcia wiersza, sprawdzenia terminu ważności i wysyłki wewnątrz
// jednej transakcji jest rozstrzygająca i nie wolno jej przestawić.
func (r *repozytoriumPowiadomien) Takt(ctx context.Context, id int64, plan PlanTaktu) (WynikTaktu, error) {
	if plan.Wyslij == nil || plan.NastepnaProba == nil {
		return WynikTaktu{}, fmt.Errorf("dane: takt powiadomienia %d bez nadajnika albo bez "+
			"reguły ponowienia — takt bez nich niczego nie robi i nie wolno mu udawać, "+
			"że zrobił", id)
	}

	wynik := WynikTaktu{}
	err := wTransakcji(ctx, r.db, func(tx *sql.Tx) error {
		zajecie, err := tx.ExecContext(ctx, zajmijPowiadomienie, id)
		if err != nil {
			return fmt.Errorf("dane: nie można zająć powiadomienia %d: %w", id, err)
		}
		zajete, err := zajecie.RowsAffected()
		if err != nil {
			return fmt.Errorf("dane: nieczytelny wynik zajęcia powiadomienia %d: %w", id, err)
		}
		if zajete == 0 {
			// Wiersz nie czeka już na wysyłkę; odczyt stanu służy dziennikowi
			// nazwanym powodem pominięcia.
			wynik.Pominieto = true
			biezacy, err := odczytajPowiadomienie(tx.QueryRowContext(ctx, stanPowiadomieniaPodZamkiem, id))
			if errors.Is(err, sql.ErrNoRows) {
				wynik.Stan = ""
				return nil
			}
			if err != nil {
				return err
			}
			wynik.Stan = biezacy.Stan
			return nil
		}

		p, err := odczytajPowiadomienie(tx.QueryRowContext(ctx, stanPowiadomieniaPodZamkiem, id))
		if err != nil {
			return err
		}

		// Termin ważności rozstrzyga przed wysyłką — powiadomienie
		// przeterminowane nie ma dolecieć.
		if p.Wygasa <= plan.Teraz {
			if _, err := tx.ExecContext(ctx, zamknijWygasnieciem, id); err != nil {
				return fmt.Errorf("dane: nie można zamknąć wygasłego powiadomienia %d: %w", id, err)
			}
			wynik.Stan = StanPowiadomieniaWygasle
			return nil
		}

		przyjeli, err := plan.Wyslij(ctx, p)
		if err != nil {
			return fmt.Errorf("dane: wysyłka powiadomienia %d nie powiodła się: %w", id, err)
		}

		if len(przyjeli) == 0 {
			// Nie było komu: wiersz zostaje w kolejce z podbitym licznikiem
			// prób i nowym terminem.
			if _, err := tx.ExecContext(ctx, przelozProbe, plan.NastepnaProba(p.ProbIle+1), id); err != nil {
				return fmt.Errorf("dane: nie można przełożyć próby powiadomienia %d: %w", id, err)
			}
			wynik.Stan = StanPowiadomieniaOczekuje
			return nil
		}

		for _, rejestracja := range przyjeli {
			if _, err := tx.ExecContext(ctx, zapiszDostarczenie, id, rejestracja, plan.Teraz); err != nil {
				return fmt.Errorf("dane: nie można zapisać doręczenia powiadomienia %d "+
					"na rejestracji %d: %w", id, rejestracja, err)
			}
			if _, err := tx.ExecContext(ctx, odnotujDostarczenieUrzadzeniu, plan.Teraz, rejestracja); err != nil {
				return fmt.Errorf("dane: nie można odnotować doręczenia na rejestracji %d: %w",
					rejestracja, err)
			}
		}
		if _, err := tx.ExecContext(ctx, zamknijDostarczeniem, plan.Teraz, id); err != nil {
			return fmt.Errorf("dane: nie można zamknąć doręczonego powiadomienia %d: %w", id, err)
		}
		wynik.Stan = StanPowiadomieniaDostarczone
		wynik.Odbiorcow = len(przyjeli)
		return nil
	})
	if err != nil {
		return WynikTaktu{}, err
	}
	return wynik, nil
}

const wygasPowiadomienia = `
	UPDATE powiadomienie SET stan = 'wygasle' WHERE stan = 'oczekuje' AND wygasa <= ?`

func (r *repozytoriumPowiadomien) Wygas(ctx context.Context, teraz string) (int, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, wygasPowiadomienia)
	if err != nil {
		return 0, err
	}
	wynik, err := polecenie.ExecContext(ctx, teraz)
	if err != nil {
		return 0, fmt.Errorf("dane: nie można wygasić przeterminowanych powiadomień: %w", err)
	}
	ile, err := wynik.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("dane: nieczytelny wynik wygaszania powiadomień: %w", err)
	}
	return int(ile), nil
}

// Odwolaj gasi wyłącznie powiadomienia oczekujące. Wiersz już doręczony zostaje
// doręczony — cofnąć budzika, który zabrzmiał, i tak się nie da, a przepisanie
// go na `odwolane` byłoby zmyśleniem historii.
const odwolajPowiadomienia = `
	UPDATE powiadomienie
	   SET stan = 'odwolane', odwolano = ?, powod_odwolania = ?
	 WHERE stan = 'oczekuje' AND byt_rodzaj = ? AND byt_id = ? AND ` + WarunekKonta

func (r *repozytoriumPowiadomien) Odwolaj(ctx context.Context, bytRodzaj, bytID, powod, teraz string) (int, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, odwolajPowiadomienia)
	if err != nil {
		return 0, err
	}
	wynik, err := polecenie.ExecContext(ctx, teraz, powod, bytRodzaj, bytID, KontoOperatora(ctx))
	if err != nil {
		return 0, fmt.Errorf("dane: nie można odwołać powiadomień o bycie %s/%s: %w",
			bytRodzaj, bytID, err)
	}
	ile, err := wynik.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("dane: nieczytelny wynik odwołania powiadomień: %w", err)
	}
	return int(ile), nil
}

const potwierdzDostarczenie = `
	UPDATE powiadomienie_dostarczenie
	   SET potwierdzono = ?
	 WHERE powiadomienie_id = ?
	   AND urzadzenie_powiadomien_id = (SELECT id FROM urzadzenie_powiadomien
	                                     WHERE kanal = ? AND klucz_kanalu = ?
	                                       AND ` + WarunekKonta + `)`

// Potwierdz zapisuje, że Operator widział powiadomienie na tym urządzeniu.
// Brak wiersza doręczenia nie jest błędem: potwierdzenie z urządzenia, któremu
// nic nie doręczono, jest pomyłką wołającego, a nie awarią rdzenia.
func (r *repozytoriumPowiadomien) Potwierdz(ctx context.Context, powiadomienieID int64,
	kanal, kluczKanalu, teraz string) error {

	polecenie, err := r.zapytania.przygotuj(ctx, potwierdzDostarczenie)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, teraz, powiadomienieID, kanal, kluczKanalu,
		KontoOperatora(ctx)); err != nil {
		return fmt.Errorf("dane: nie można potwierdzić powiadomienia %d na %q/%q: %w",
			powiadomienieID, kanal, kluczKanalu, err)
	}
	return nil
}

const dostarczeniaPowiadomienia = `
	SELECT urzadzenie_powiadomien_id FROM powiadomienie_dostarczenie
	 WHERE powiadomienie_id = ? ORDER BY urzadzenie_powiadomien_id`

func (r *repozytoriumPowiadomien) Dostarczenia(ctx context.Context, powiadomienieID int64) ([]int64, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, dostarczeniaPowiadomienia)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, powiadomienieID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać doręczeń powiadomienia %d: %w",
			powiadomienieID, err)
	}
	defer wiersze.Close()

	lista := []int64{}
	for wiersze.Next() {
		var id int64
		if err := wiersze.Scan(&id); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz doręczenia: %w", err)
		}
		lista = append(lista, id)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt doręczeń: %w", err)
	}
	return lista, nil
}

func odczytajPowiadomienie(wiersz skaner) (Powiadomienie, error) {
	var p Powiadomienie
	var bytRodzaj, bytID, dostarczono, odwolano, powodOdwolania sql.NullString
	err := wiersz.Scan(&p.ID, &p.Tytul, &p.Tresc, &p.Priorytet, &p.Powod, &bytRodzaj, &bytID,
		&p.Stan, &p.ProbIle, &p.NastepnaProba, &p.Wygasa, &p.Utworzono, &dostarczono,
		&odwolano, &powodOdwolania)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Powiadomienie{}, err
		}
		return Powiadomienie{}, fmt.Errorf("dane: nieczytelny wiersz powiadomienia: %w", err)
	}
	p.BytRodzaj = tekstZKolumny(bytRodzaj)
	p.BytID = tekstZKolumny(bytID)
	p.Dostarczono = tekstZKolumny(dostarczono)
	p.Odwolano = tekstZKolumny(odwolano)
	p.PowodOdwolania = tekstZKolumny(powodOdwolania)
	return p, nil
}
