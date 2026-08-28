// Plik prowadzi przebiegi automatyki: trwałość okna Execution Monitor; przebieg nie jest drugą kolejką, tylko
// zapamiętuje to, czego kolejka nie wie — którą automatykę realizuje, ile etapów miała definicja i dlaczego wykonanie się nie powiodło.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// Przebieg to wiersz tabeli `przebieg_automatyki` niosący stan uruchomienia jednej automatyki w bazie.
type Przebieg struct {
	ID             int64
	Kod            string
	AutomatykaID   int64
	AutomatykaKod  string
	KolejkaID      *int64
	Stan           string
	EtapBiezacy    int
	Etapow         int
	Proba          int
	KomunikatBledu *string
	Rozpoczeto     string
	Zakonczono     *string
}

const (
	kolumnyPrzebiegu = `p.id, p.identyfikator_zewnetrzny, p.automatyka_id, a.identyfikator_zewnetrzny,
	                    p.kolejka_id, p.stan, p.etap_biezacy, p.etapow, p.proba,
	                    p.komunikat_bledu, p.rozpoczeto, p.zakonczono`

	zrodloPrzebiegu = ` FROM przebieg_automatyki p JOIN automatyka a ON a.id = p.automatyka_id`

	zapiszPrzebieg = `INSERT INTO przebieg_automatyki
	                  (identyfikator_zewnetrzny, automatyka_id, kolejka_id, stan,
	                   etap_biezacy, etapow, proba, komunikat_bledu, zakonczono)
	                  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	                  ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                      kolejka_id = excluded.kolejka_id,
	                      stan = excluded.stan,
	                      etap_biezacy = excluded.etap_biezacy,
	                      etapow = excluded.etapow,
	                      proba = excluded.proba,
	                      komunikat_bledu = excluded.komunikat_bledu,
	                      zakonczono = excluded.zakonczono`

	pobierzPrzebieg = `SELECT ` + kolumnyPrzebiegu + zrodloPrzebiegu +
		` WHERE p.identyfikator_zewnetrzny = ?`

	pobierzPrzebiegKolejki = `SELECT ` + kolumnyPrzebiegu + zrodloPrzebiegu +
		` WHERE p.kolejka_id = ?`

	// Zero w pierwszym argumencie znaczy „wszystkie automatyki” — Execution
	// Monitor otwarty bez wskazania automatyki pokazuje przebiegi wszystkich.
	listaPrzebiegow = `SELECT ` + kolumnyPrzebiegu + zrodloPrzebiegu +
		` WHERE (? = 0 OR p.automatyka_id = ?)
		  ORDER BY p.id DESC LIMIT ?`
)

// ZapiszPrzebieg zakłada przebieg albo nadpisuje zastany i zwraca jego pełny stan po zapisie z bazy danych.
func (r *repozytoriumAutomatyk) ZapiszPrzebieg(ctx context.Context, przebieg Przebieg) (Przebieg, error) {
	if przebieg.Kod == "" || przebieg.AutomatykaID == 0 {
		return Przebieg{}, fmt.Errorf("dane: przebieg bez identyfikatora albo bez automatyki")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszPrzebieg)
	if err != nil {
		return Przebieg{}, err
	}
	_, err = polecenie.ExecContext(ctx, przebieg.Kod, przebieg.AutomatykaID,
		liczbaDoKolumny(przebieg.KolejkaID), przebieg.Stan, przebieg.EtapBiezacy,
		przebieg.Etapow, przebieg.Proba, tekstDoKolumny(przebieg.KomunikatBledu),
		tekstDoKolumny(przebieg.Zakonczono))
	if err != nil {
		return Przebieg{}, fmt.Errorf("dane: nie można zapisać przebiegu %q: %w", przebieg.Kod, err)
	}
	return r.Przebieg(ctx, przebieg.Kod)
}

// Przebieg zwraca przebieg automatyki o wskazanym kodzie zewnętrznym wprost z bazy danych repozytorium.
func (r *repozytoriumAutomatyk) Przebieg(ctx context.Context, kod string) (Przebieg, error) {
	return r.jedenPrzebieg(ctx, pobierzPrzebieg, kod, "przebieg "+kod)
}

// PrzebiegKolejki odnajduje przebieg realizowany przez wskazaną kolejkę.
// Kolejka bez przebiegu wraca jako ErrBrakWiersza — nie każda kolejka jest
// wykonaniem automatyki (pętla sesyjna korzysta z tego samego silnika).
func (r *repozytoriumAutomatyk) PrzebiegKolejki(ctx context.Context, kolejkaID int64) (Przebieg, error) {
	return r.jedenPrzebieg(ctx, pobierzPrzebiegKolejki, kolejkaID,
		fmt.Sprintf("przebieg kolejki %d", kolejkaID))
}

// jedenPrzebieg wykonuje odczyt pojedynczego wiersza przebiegu, wspólny obu sposobom jego doboru z bazy.
func (r *repozytoriumAutomatyk) jedenPrzebieg(ctx context.Context, zapytanie string,
	klucz any, opis string) (Przebieg, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return Przebieg{}, err
	}
	przebieg, err := odczytajPrzebieg(polecenie.QueryRowContext(ctx, klucz))
	if errors.Is(err, sql.ErrNoRows) {
		return Przebieg{}, ErrBrakWiersza
	}
	if err != nil {
		return Przebieg{}, fmt.Errorf("dane: nieczytelny %s: %w", opis, err)
	}
	return przebieg, nil
}

// Przebiegi zwraca przebiegi od najnowszego. Automatyka zerowa znaczy wykaz
// wszystkich — Execution Monitor bywa otwierany bez wskazania automatyki.
func (r *repozytoriumAutomatyk) Przebiegi(ctx context.Context,
	automatykaID int64, limit int) ([]Przebieg, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaPrzebiegow)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, automatykaID, automatykaID, granicaWykazu(limit))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać przebiegów automatyki %d: %w", automatykaID, err)
	}
	defer wiersze.Close()

	lista := []Przebieg{}
	for wiersze.Next() {
		przebieg, err := odczytajPrzebieg(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz przebiegu: %w", err)
		}
		lista = append(lista, przebieg)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt przebiegów: %w", err)
	}
	return lista, nil
}

// odczytajPrzebieg składa strukturę przebiegu wprost z jednego wiersza wyniku zapytania SQL do bazy danych.
func odczytajPrzebieg(wiersz skaner) (Przebieg, error) {
	var przebieg Przebieg
	var kolejka sql.NullInt64
	var komunikat, zakonczono sql.NullString
	err := wiersz.Scan(&przebieg.ID, &przebieg.Kod, &przebieg.AutomatykaID, &przebieg.AutomatykaKod,
		&kolejka, &przebieg.Stan, &przebieg.EtapBiezacy, &przebieg.Etapow, &przebieg.Proba,
		&komunikat, &przebieg.Rozpoczeto, &zakonczono)
	if err != nil {
		return Przebieg{}, err
	}
	przebieg.KolejkaID = liczbaZKolumny(kolejka)
	przebieg.KomunikatBledu, przebieg.Zakonczono = tekstZKolumny(komunikat), tekstZKolumny(zakonczono)
	return przebieg, nil
}
