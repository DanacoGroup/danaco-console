// Odpowiedzialność pliku: dostęp do kolejek (tabela `kolejka`). Jeden silnik
// obsługuje pętlę sesyjną i MultitaskingAI — bez drugiego kompletu tabel.
// Zmiana stanu dotyka stanu i dziennika akcji, więc idzie w transakcji.
// Zlecenia kolejki obsługuje `pozycje_kolejki.go`.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"danacoconsole/shared"
)

// Kolejka to wiersz tabeli `kolejka`.
type Kolejka struct {
	ID                 int64
	Nazwa              string
	Rodzaj             string
	SesjaID            *int64
	OknoKoordynatoraID *int64
	Stan               shared.QueueStatus
	Utworzono          string
	Zaktualizowano     string
}

// RepozytoriumKolejek jest kontraktem obszaru kolejek.
type RepozytoriumKolejek interface {
	UtworzKolejke(ctx context.Context, kolejka Kolejka) (int64, error)
	PobierzKolejke(ctx context.Context, id int64) (Kolejka, error)
	ZmienStanKolejki(ctx context.Context, id int64, stan shared.QueueStatus, akcja shared.QueueAction) error
	DodajPozycje(ctx context.Context, pozycja Pozycja) (int64, error)
	ZmienStanPozycji(ctx context.Context, id int64, stan string, werdykt *string) error
	ZwiekszObieg(ctx context.Context, pozycjaID int64) (int, error)
	ListaPozycji(ctx context.Context, kolejkaID int64) ([]Pozycja, error)
	Dziennik(ctx context.Context, kolejkaID int64, limit int) ([]WpisDziennika, error)
	// LiczbaCzynnych liczy kolejki poza stanem końcowym — miara stanu platformy
	// dla mobilnego centrum dowodzenia. Rachunek stoi tutaj, bo tabelę `kolejka`
	// prowadzi to repozytorium i drugiego czytelnika mieć nie będzie;
	// ciało metody leży w `mobile.go`.
	LiczbaCzynnych(ctx context.Context) (int, error)

	// Zlecenia kolejki i polityka kolejki (`zlecenia_kolejki.go`, migracje
	// 268-269) — byty rodziny `queue.item.*`, `queue.policy.set`,
	// `queue.dead.list` i `queue.depth.get`.
	DodajZlecenie(ctx context.Context, zlecenie Zlecenie) (Zlecenie, error)
	Zlecenie(ctx context.Context, kod string) (Zlecenie, error)
	ZlecenieKluczem(ctx context.Context, kolejkaID int64, klucz string) (Zlecenie, error)
	ZleceniaKolejki(ctx context.Context, kolejkaID int64, stan string, limit int) ([]Zlecenie, int, error)
	ZleceniaMartwe(ctx context.Context, kolejkaID int64, limit int) ([]Zlecenie, error)
	ZmienStanZlecenia(ctx context.Context, zlecenieID int64, stan string) error
	OdlozZlecenie(ctx context.Context, zlecenieID int64, termin string) error
	UstawWarunekZlecenia(ctx context.Context, zlecenieID int64, warunek *string) error
	SkierujZlecenie(ctx context.Context, zlecenieID, kolejkaID int64, ekspert *string) error
	GlebokoscKolejki(ctx context.Context, kolejkaID int64, odcinekSekundy int,
		od string) ([]PunktGlebokosci, error)
	PolitykaKolejki(ctx context.Context, kolejkaID int64) (PolitykaKolejki, error)
	ZapiszPolitykeKolejki(ctx context.Context, kolejkaID int64, polityka PolitykaKolejki) error
}

// domyslnyRodzajKolejki odpowiada wartości domyślnej kolumny w schemacie.
const domyslnyRodzajKolejki = "sesyjna"

const (
	kolumnyKolejki = `id, nazwa, rodzaj, sesja_id, okno_koordynatora_id, stan,
	                  utworzono, zaktualizowano`

	wstawKolejke = `INSERT INTO kolejka (nazwa, rodzaj, sesja_id, okno_koordynatora_id, stan)
	                VALUES (?, ?, ?, ?, ?)`

	pobierzKolejke = `SELECT ` + kolumnyKolejki + ` FROM kolejka WHERE id = ?`

	zmienStanKolejki = `UPDATE kolejka
	                    SET stan = ?, zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                    WHERE id = ?`

	odczytStanuKolejki = `SELECT stan FROM kolejka WHERE id = ?`
)

type repozytoriumKolejek struct {
	zapytania *zapytania
	db        *sql.DB
}

func noweRepozytoriumKolejek(z *zapytania, db *sql.DB) *repozytoriumKolejek {
	return &repozytoriumKolejek{zapytania: z, db: db}
}

// UtworzKolejke zakłada kolejkę i odnotowuje założenie w dzienniku.
func (r *repozytoriumKolejek) UtworzKolejke(ctx context.Context, kolejka Kolejka) (int64, error) {
	stan, err := stanKolejkiNaBaze(kolejka.Stan)
	if err != nil {
		return 0, err
	}
	rodzaj := kolejka.Rodzaj
	if rodzaj == "" {
		rodzaj = domyslnyRodzajKolejki
	}
	var id int64
	err = wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawKolejke)
		if err != nil {
			return err
		}
		wynik, err := polecenie.ExecContext(ctx, kolejka.Nazwa, rodzaj,
			liczbaDoKolumny(kolejka.SesjaID), liczbaDoKolumny(kolejka.OknoKoordynatoraID), stan)
		if err != nil {
			return fmt.Errorf("dane: nie można założyć kolejki %q: %w", kolejka.Nazwa, err)
		}
		if id, err = wynik.LastInsertId(); err != nil {
			return fmt.Errorf("dane: nieznany identyfikator założonej kolejki: %w", err)
		}
		return dopiszWpisDziennika(ctx, r.zapytania, transakcja, WpisDziennika{
			KolejkaID: id, Akcja: "utworzenie", StanPo: &stan,
		})
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

// PobierzKolejke zwraca kolejkę o wskazanym identyfikatorze.
func (r *repozytoriumKolejek) PobierzKolejke(ctx context.Context, id int64) (Kolejka, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzKolejke)
	if err != nil {
		return Kolejka{}, err
	}
	var kolejka Kolejka
	var sesja, okno sql.NullInt64
	var stan string
	err = polecenie.QueryRowContext(ctx, id).Scan(&kolejka.ID, &kolejka.Nazwa, &kolejka.Rodzaj,
		&sesja, &okno, &stan, &kolejka.Utworzono, &kolejka.Zaktualizowano)
	if errors.Is(err, sql.ErrNoRows) {
		return Kolejka{}, fmt.Errorf("dane: kolejka %d nie istnieje", id)
	}
	if err != nil {
		return Kolejka{}, fmt.Errorf("dane: nieczytelny wiersz kolejki %d: %w", id, err)
	}
	kolejka.SesjaID = liczbaZKolumny(sesja)
	kolejka.OknoKoordynatoraID = liczbaZKolumny(okno)
	if kolejka.Stan, err = stanKolejkiZBazy(stan); err != nil {
		return Kolejka{}, err
	}
	return kolejka, nil
}

// ZmienStanKolejki zapisuje nowy stan wraz z wywołaną akcją. Każda zmiana
// zostawia wpis w dzienniku, razem ze stanem sprzed przejścia.
func (r *repozytoriumKolejek) ZmienStanKolejki(ctx context.Context, id int64,
	stan shared.QueueStatus, akcja shared.QueueAction) error {

	kolumna, err := stanKolejkiNaBaze(stan)
	if err != nil {
		return err
	}
	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		poprzedni, err := stanKolejkiSprzed(ctx, r.zapytania, transakcja, id)
		if err != nil {
			return err
		}
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, zmienStanKolejki)
		if err != nil {
			return err
		}
		wynik, err := polecenie.ExecContext(ctx, kolumna, id)
		if err != nil {
			return fmt.Errorf("dane: nie można zmienić stanu kolejki %d: %w", id, err)
		}
		if err := sprawdzTrafienie(wynik, "kolejka", id); err != nil {
			return err
		}
		return dopiszWpisDziennika(ctx, r.zapytania, transakcja, WpisDziennika{
			KolejkaID: id, Akcja: string(akcja), StanPrzed: poprzedni, StanPo: &kolumna,
		})
	})
}

// stanKolejkiSprzed odczytuje stan sprzed zmiany, żeby dziennik pokazywał
// przejście, a nie samą wartość docelową.
func stanKolejkiSprzed(ctx context.Context, z *zapytania, transakcja *sql.Tx, id int64) (*string, error) {
	polecenie, err := z.wTransakcji(ctx, transakcja, odczytStanuKolejki)
	if err != nil {
		return nil, err
	}
	var stan string
	err = polecenie.QueryRowContext(ctx, id).Scan(&stan)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("dane: kolejka %d nie istnieje", id)
	}
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać stanu kolejki %d: %w", id, err)
	}
	return &stan, nil
}
