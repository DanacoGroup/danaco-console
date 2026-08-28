// Odpowiedzialność pliku: rodzina `extension.*` — warstwa integracji zewnętrznych: transport i poświadczenie
// (`integracja_rozszerzenia`), webhooki (`webhook_rozszerzenia`) oraz odwzorowania danych (`mapowanie_rozszerzenia`).
package dane

import (
	"context"
	"database/sql"
	"fmt"
)

// IntegracjaRozszerzenia to wiersz tabeli `integracja_rozszerzenia`, niosący transport i poświadczenie integracji.
type IntegracjaRozszerzenia struct {
	RozszerzenieKod  string
	Transport        *string
	Adres            *string
	Polecenie        *string
	SposobLogowania  *string
	OdwolanieSekretu *string
	Zakresy          []string
	Zaktualizowano   int64
}

// WebhookRozszerzenia to wiersz tabeli `webhook_rozszerzenia`, niosący adres i kierunek jednego webhooka.
type WebhookRozszerzenia struct {
	ID               int64
	Kod              string
	RozszerzenieKod  string
	Kierunek         string
	Adres            *string
	AdresNasluchu    *string
	Zdarzenia        []string
	OdwolanieSekretu *string
	Czynny           bool
	Zaktualizowano   int64
}

// MapowanieRozszerzenia to wiersz tabeli `mapowanie_rozszerzenia`, niosący odwzorowanie danych integracji.
type MapowanieRozszerzenia struct {
	ID              int64
	Kod             string
	RozszerzenieKod string
	Nazwa           string
	Reguly          string
	Zaktualizowano  int64
}

const (
	kolumnyIntegracjiRozszerzenia = `rozszerzenie_kod, transport, adres, polecenie,
	                                 sposob_logowania, odwolanie_sekretu, zakresy, zaktualizowano`

	// Pola podane jako brak NIE kasują wartości zastanych: transport i
	// poświadczenie nadaje się osobnymi komendami, a każda zna tylko swoją część.
	zapiszIntegracjeRozszerzenia = `INSERT INTO integracja_rozszerzenia
	                                (rozszerzenie_kod, transport, adres, polecenie,
	                                 sposob_logowania, odwolanie_sekretu, zakresy, zaktualizowano)
	                                VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	                                ON CONFLICT(rozszerzenie_kod) DO UPDATE SET
	                                    transport = IFNULL(excluded.transport,
	                                                       integracja_rozszerzenia.transport),
	                                    adres = IFNULL(excluded.adres, integracja_rozszerzenia.adres),
	                                    polecenie = IFNULL(excluded.polecenie,
	                                                       integracja_rozszerzenia.polecenie),
	                                    sposob_logowania = IFNULL(excluded.sposob_logowania,
	                                                              integracja_rozszerzenia.sposob_logowania),
	                                    odwolanie_sekretu = IFNULL(excluded.odwolanie_sekretu,
	                                                               integracja_rozszerzenia.odwolanie_sekretu),
	                                    zakresy = IFNULL(excluded.zakresy,
	                                                     integracja_rozszerzenia.zakresy),
	                                    zaktualizowano = excluded.zaktualizowano`

	pobierzIntegracjeRozszerzenia = `SELECT ` + kolumnyIntegracjiRozszerzenia + `
	                                 FROM integracja_rozszerzenia WHERE rozszerzenie_kod = ?`

	kolumnyWebhookaRozszerzenia = `id, identyfikator_zewnetrzny, rozszerzenie_kod, kierunek,
	                               adres, adres_nasluchu, zdarzenia, odwolanie_sekretu,
	                               czynny, zaktualizowano`

	zapiszWebhookRozszerzenia = `INSERT INTO webhook_rozszerzenia
	                             (identyfikator_zewnetrzny, rozszerzenie_kod, kierunek, adres,
	                              adres_nasluchu, zdarzenia, odwolanie_sekretu, czynny, zaktualizowano)
	                             VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	                             ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                                 kierunek = excluded.kierunek,
	                                 adres = excluded.adres,
	                                 adres_nasluchu = excluded.adres_nasluchu,
	                                 zdarzenia = excluded.zdarzenia,
	                                 odwolanie_sekretu = excluded.odwolanie_sekretu,
	                                 czynny = excluded.czynny,
	                                 zaktualizowano = excluded.zaktualizowano`

	pobierzWebhookRozszerzenia = `SELECT ` + kolumnyWebhookaRozszerzenia + `
	                              FROM webhook_rozszerzenia WHERE identyfikator_zewnetrzny = ?`

	listaWebhookowRozszerzenia = `SELECT ` + kolumnyWebhookaRozszerzenia + `
	                              FROM webhook_rozszerzenia
	                              WHERE (? = '' OR rozszerzenie_kod = ?) AND (? = '' OR kierunek = ?)
	                              ORDER BY rozszerzenie_kod, kierunek, id`

	kolumnyMapowaniaRozszerzenia = `id, identyfikator_zewnetrzny, rozszerzenie_kod, nazwa,
	                                reguly, zaktualizowano`

	zapiszMapowanieRozszerzenia = `INSERT INTO mapowanie_rozszerzenia
	                               (identyfikator_zewnetrzny, rozszerzenie_kod, nazwa, reguly,
	                                zaktualizowano)
	                               VALUES (?, ?, ?, ?, ?)
	                               ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                                   nazwa = excluded.nazwa,
	                                   reguly = excluded.reguly,
	                                   zaktualizowano = excluded.zaktualizowano`

	pobierzMapowanieRozszerzenia = `SELECT ` + kolumnyMapowaniaRozszerzenia + `
	                                FROM mapowanie_rozszerzenia WHERE identyfikator_zewnetrzny = ?`
)

// ZapiszIntegracjeRozszerzenia zapisuje transport albo poświadczenie integracji rozszerzenia w jednym wierszu tabeli.
func (r *repozytoriumRozszerzen) ZapiszIntegracjeRozszerzenia(ctx context.Context,
	integracja IntegracjaRozszerzenia) (IntegracjaRozszerzenia, error) {

	if integracja.RozszerzenieKod == "" {
		return IntegracjaRozszerzenia{}, fmt.Errorf("dane: integracja bez pozycji katalogu")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszIntegracjeRozszerzenia)
	if err != nil {
		return IntegracjaRozszerzenia{}, err
	}
	_, err = polecenie.ExecContext(ctx, integracja.RozszerzenieKod,
		tekstDoKolumny(integracja.Transport), tekstDoKolumny(integracja.Adres),
		tekstDoKolumny(integracja.Polecenie), tekstDoKolumny(integracja.SposobLogowania),
		tekstDoKolumny(integracja.OdwolanieSekretu), listaDoKolumny(integracja.Zakresy),
		integracja.Zaktualizowano)
	if err != nil {
		return IntegracjaRozszerzenia{}, fmt.Errorf("dane: nie można zapisać integracji %q: %w",
			integracja.RozszerzenieKod, err)
	}
	return r.IntegracjaRozszerzenia(ctx, integracja.RozszerzenieKod)
}

// IntegracjaRozszerzenia zwraca transport i poświadczenie integracji zapisane dla wskazanego rozszerzenia.
func (r *repozytoriumRozszerzen) IntegracjaRozszerzenia(ctx context.Context,
	rozszerzenie string) (IntegracjaRozszerzenia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzIntegracjeRozszerzenia)
	if err != nil {
		return IntegracjaRozszerzenia{}, err
	}
	var integracja IntegracjaRozszerzenia
	var transport, adres, komenda, sposob, sekret, zakresy sql.NullString
	err = polecenie.QueryRowContext(ctx, rozszerzenie).Scan(&integracja.RozszerzenieKod,
		&transport, &adres, &komenda, &sposob, &sekret, &zakresy, &integracja.Zaktualizowano)
	if err == sql.ErrNoRows {
		return IntegracjaRozszerzenia{}, ErrBrakWiersza
	}
	if err != nil {
		return IntegracjaRozszerzenia{}, fmt.Errorf("dane: nieczytelna integracja %q: %w",
			rozszerzenie, err)
	}
	integracja.Transport = tekstZKolumny(transport)
	integracja.Adres = tekstZKolumny(adres)
	integracja.Polecenie = tekstZKolumny(komenda)
	integracja.SposobLogowania = tekstZKolumny(sposob)
	integracja.OdwolanieSekretu = tekstZKolumny(sekret)
	integracja.Zakresy = listaZKolumny(zakresy)
	return integracja, nil
}

// ZapiszWebhookRozszerzenia zapisuje webhook rozszerzenia, nadpisując wiersz istniejący pod tym samym identyfikatorem.
func (r *repozytoriumRozszerzen) ZapiszWebhookRozszerzenia(ctx context.Context,
	webhook WebhookRozszerzenia) (WebhookRozszerzenia, error) {

	if webhook.Kod == "" || webhook.RozszerzenieKod == "" || webhook.Kierunek == "" {
		return WebhookRozszerzenia{}, fmt.Errorf(
			"dane: webhook bez identyfikatora, pozycji albo kierunku")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszWebhookRozszerzenia)
	if err != nil {
		return WebhookRozszerzenia{}, err
	}
	_, err = polecenie.ExecContext(ctx, webhook.Kod, webhook.RozszerzenieKod, webhook.Kierunek,
		tekstDoKolumny(webhook.Adres), tekstDoKolumny(webhook.AdresNasluchu),
		listaDoKolumny(webhook.Zdarzenia), tekstDoKolumny(webhook.OdwolanieSekretu),
		liczbaLogiczna(webhook.Czynny), webhook.Zaktualizowano)
	if err != nil {
		return WebhookRozszerzenia{}, fmt.Errorf("dane: nie można zapisać webhooka %q: %w",
			webhook.Kod, err)
	}
	return r.WebhookRozszerzenia(ctx, webhook.Kod)
}

// WebhookRozszerzenia zwraca jeden webhook rozszerzenia wskazany jego identyfikatorem tekstowym w tabeli.
func (r *repozytoriumRozszerzen) WebhookRozszerzenia(ctx context.Context, kod string) (WebhookRozszerzenia, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWebhookRozszerzenia)
	if err != nil {
		return WebhookRozszerzenia{}, err
	}
	webhook, err := odczytajWebhookRozszerzenia(polecenie.QueryRowContext(ctx, kod))
	if err == sql.ErrNoRows {
		return WebhookRozszerzenia{}, ErrBrakWiersza
	}
	if err != nil {
		return WebhookRozszerzenia{}, fmt.Errorf("dane: nieczytelny webhook %q: %w", kod, err)
	}
	return webhook, nil
}

// WebhookiRozszerzen zwraca wykaz webhooków rozszerzenia, opcjonalnie zawężony pozycją i kierunkiem zdarzenia.
func (r *repozytoriumRozszerzen) WebhookiRozszerzen(ctx context.Context,
	rozszerzenie, kierunek string) ([]WebhookRozszerzenia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaWebhookowRozszerzenia)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, rozszerzenie, rozszerzenie, kierunek, kierunek)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać webhooków: %w", err)
	}
	defer wiersze.Close()

	lista := []WebhookRozszerzenia{}
	for wiersze.Next() {
		webhook, err := odczytajWebhookRozszerzenia(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz webhooka: %w", err)
		}
		lista = append(lista, webhook)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt webhooków: %w", err)
	}
	return lista, nil
}

// ZapiszMapowanieRozszerzenia zapisuje odwzorowanie danych integracji rozszerzenia w jednym wierszu tabeli.
func (r *repozytoriumRozszerzen) ZapiszMapowanieRozszerzenia(ctx context.Context,
	mapowanie MapowanieRozszerzenia) (MapowanieRozszerzenia, error) {

	if mapowanie.Kod == "" || mapowanie.RozszerzenieKod == "" || mapowanie.Nazwa == "" {
		return MapowanieRozszerzenia{}, fmt.Errorf(
			"dane: odwzorowanie bez identyfikatora, pozycji albo nazwy")
	}
	if mapowanie.Reguly == "" {
		mapowanie.Reguly = "{}"
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszMapowanieRozszerzenia)
	if err != nil {
		return MapowanieRozszerzenia{}, err
	}
	_, err = polecenie.ExecContext(ctx, mapowanie.Kod, mapowanie.RozszerzenieKod,
		mapowanie.Nazwa, mapowanie.Reguly, mapowanie.Zaktualizowano)
	if err != nil {
		return MapowanieRozszerzenia{}, fmt.Errorf("dane: nie można zapisać odwzorowania %q: %w",
			mapowanie.Kod, err)
	}
	return r.MapowanieRozszerzenia(ctx, mapowanie.Kod)
}

// MapowanieRozszerzenia zwraca jedno odwzorowanie danych integracji wskazane jego identyfikatorem tekstowym.
func (r *repozytoriumRozszerzen) MapowanieRozszerzenia(ctx context.Context, kod string) (MapowanieRozszerzenia, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzMapowanieRozszerzenia)
	if err != nil {
		return MapowanieRozszerzenia{}, err
	}
	var mapowanie MapowanieRozszerzenia
	err = polecenie.QueryRowContext(ctx, kod).Scan(&mapowanie.ID, &mapowanie.Kod,
		&mapowanie.RozszerzenieKod, &mapowanie.Nazwa, &mapowanie.Reguly, &mapowanie.Zaktualizowano)
	if err == sql.ErrNoRows {
		return MapowanieRozszerzenia{}, ErrBrakWiersza
	}
	if err != nil {
		return MapowanieRozszerzenia{}, fmt.Errorf("dane: nieczytelne odwzorowanie %q: %w", kod, err)
	}
	return mapowanie, nil
}

// odczytajWebhookRozszerzenia składa strukturę webhooka rozszerzenia z jednego wiersza wyniku zapytania.
func odczytajWebhookRozszerzenia(wiersz skaner) (WebhookRozszerzenia, error) {
	var webhook WebhookRozszerzenia
	var adres, nasluch, zdarzenia, sekret sql.NullString
	var czynny int
	err := wiersz.Scan(&webhook.ID, &webhook.Kod, &webhook.RozszerzenieKod, &webhook.Kierunek,
		&adres, &nasluch, &zdarzenia, &sekret, &czynny, &webhook.Zaktualizowano)
	if err != nil {
		return WebhookRozszerzenia{}, err
	}
	webhook.Adres = tekstZKolumny(adres)
	webhook.AdresNasluchu = tekstZKolumny(nasluch)
	webhook.Zdarzenia = listaZKolumny(zdarzenia)
	webhook.OdwolanieSekretu = tekstZKolumny(sekret)
	webhook.Czynny = czynny == 1
	return webhook, nil
}
