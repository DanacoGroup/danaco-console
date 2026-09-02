// Plik niesie udostępnienia odnośnikiem i nasłuchy zewnętrzne wraz z ich
// zdarzeniami. Odwołanie udostępnienia zapisuje znacznik czasu zamiast kasować
// wiersz. Zapis nasłuchu podmienia komplet jego zdarzeń, nie tylko różnicę.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// UdostepnienieBiblioteki to wiersz udostępnienia odnośnikiem: token dostępu,
// zasięg, cel, termin wygaśnięcia oraz znacznik ewentualnego odwołania.
type UdostepnienieBiblioteki struct {
	ID        int64
	Kod       string
	Zasieg    string
	CelKod    string
	Token     string
	Wygasa    *string
	Odwolano  *string
	Utworzono string
}

// WebhookBiblioteki to wiersz nasłuchu zewnętrznego wraz z jego zdarzeniami:
// adres docelowy, sekret podpisu, znacznik czynności i ostatnie zgłoszenie.
type WebhookBiblioteki struct {
	ID                 int64
	Kod                string
	Adres              string
	Sekret             *string
	Czynny             bool
	Zdarzenia          []string
	OstatnieZgloszenie *string
	Utworzono          string
}

const (
	kolumnyUdostepnieniaBiblioteki = `id, identyfikator_zewnetrzny, zasieg, cel_kod, token,
	                                  wygasa, odwolano, utworzono`

	zapiszUdostepnienieBiblioteki = `INSERT INTO udostepnienie_biblioteki
	                                 (identyfikator_zewnetrzny, zasieg, cel_kod, token, wygasa,
	                                  konto_id)
	                                 VALUES (?, ?, ?, ?, ?, ` + WskazanieKonta + `)`

	pobierzUdostepnienieBiblioteki = `SELECT ` + kolumnyUdostepnieniaBiblioteki + `
	                                  FROM udostepnienie_biblioteki
	                                  WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	odwolajUdostepnienieBiblioteki = `UPDATE udostepnienie_biblioteki
	                                  SET odwolano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                                  WHERE identyfikator_zewnetrzny = ? AND odwolano IS NULL
	                                    AND ` + WarunekKonta

	kolumnyWebhookaBiblioteki = `id, identyfikator_zewnetrzny, adres, sekret, czynny,
	                             ostatnie_zgloszenie, utworzono`

	zapiszWebhookBiblioteki = `INSERT INTO webhook_biblioteki
	                           (identyfikator_zewnetrzny, adres, sekret, czynny, ostatnie_zgloszenie,
	                            konto_id)
	                           VALUES (?, ?, ?, ?, ?, ` + WskazanieKonta + `)
	                           ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                               adres = excluded.adres,
	                               sekret = excluded.sekret,
	                               czynny = excluded.czynny,
	                               ostatnie_zgloszenie = excluded.ostatnie_zgloszenie
	                           WHERE ` + WarunekKonta

	pobierzWebhookBiblioteki = `SELECT ` + kolumnyWebhookaBiblioteki + ` FROM webhook_biblioteki
	                            WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	usunWebhookBiblioteki = `DELETE FROM webhook_biblioteki
	                         WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	usunZdarzeniaWebhookaBiblioteki = `DELETE FROM zdarzenie_webhooka_biblioteki WHERE webhook_id = ?`

	wstawZdarzenieWebhookaBiblioteki = `INSERT INTO zdarzenie_webhooka_biblioteki (webhook_id, zdarzenie)
	                                    VALUES (?, ?)
	                                    ON CONFLICT(webhook_id, zdarzenie) DO NOTHING`

	listaZdarzenWebhookaBiblioteki = `SELECT zdarzenie FROM zdarzenie_webhooka_biblioteki
	                                  WHERE webhook_id = ? ORDER BY zdarzenie`
)

// ZapiszUdostepnienie wystawia odnośnik wraz z tokenem i oddaje zapisany wiersz.
// Kod albo token puste są odrzucane jako błąd.
func (r *repozytoriumBiblioteki) ZapiszUdostepnienie(ctx context.Context,
	udostepnienie UdostepnienieBiblioteki) (UdostepnienieBiblioteki, error) {

	if strings.TrimSpace(udostepnienie.Kod) == "" || strings.TrimSpace(udostepnienie.Token) == "" {
		return UdostepnienieBiblioteki{}, fmt.Errorf("dane: udostępnienie bez identyfikatora albo tokenu")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszUdostepnienieBiblioteki)
	if err != nil {
		return UdostepnienieBiblioteki{}, err
	}
	_, err = polecenie.ExecContext(ctx, udostepnienie.Kod, udostepnienie.Zasieg,
		udostepnienie.CelKod, udostepnienie.Token, tekstDoKolumny(udostepnienie.Wygasa),
		KontoOperatora(ctx))
	if err != nil {
		return UdostepnienieBiblioteki{}, fmt.Errorf("dane: nie można zapisać udostępnienia %q: %w",
			udostepnienie.Kod, err)
	}
	odczyt, err := r.zapytania.przygotuj(ctx, pobierzUdostepnienieBiblioteki)
	if err != nil {
		return UdostepnienieBiblioteki{}, err
	}
	zapisane, err := odczytajUdostepnienieBiblioteki(
		odczyt.QueryRowContext(ctx, udostepnienie.Kod, KontoOperatora(ctx)))
	if err != nil {
		return UdostepnienieBiblioteki{}, fmt.Errorf("dane: nieczytelne udostępnienie %q: %w",
			udostepnienie.Kod, err)
	}
	return zapisane, nil
}

// Udostepnienia zwraca udostępnienia od najnowszego, opcjonalnie zawężone
// do jednego celu i do wierszy wciąż czynnych.
func (r *repozytoriumBiblioteki) Udostepnienia(ctx context.Context, celKod *string,
	tylkoCzynne bool) ([]UdostepnienieBiblioteki, error) {

	warunki := []string{WarunekKonta}
	argumenty := []any{KontoOperatora(ctx)}
	if celKod != nil && *celKod != "" {
		warunki = append(warunki, "cel_kod = ?")
		argumenty = append(argumenty, *celKod)
	}
	if tylkoCzynne {
		// Czynne znaczy nieodwołane i nieprzeterminowane.
		warunki = append(warunki, `odwolano IS NULL AND
		                           (wygasa IS NULL OR wygasa > strftime('%Y-%m-%dT%H:%M:%fZ','now'))`)
	}
	zapytanie := `SELECT ` + kolumnyUdostepnieniaBiblioteki + ` FROM udostepnienie_biblioteki
	              WHERE ` + strings.Join(warunki, " AND ") + ` ORDER BY utworzono DESC, id DESC`

	wiersze, err := r.db.QueryContext(ctx, zapytanie, argumenty...)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać udostępnień: %w", err)
	}
	defer wiersze.Close()

	lista := []UdostepnienieBiblioteki{}
	for wiersze.Next() {
		udostepnienie, err := odczytajUdostepnienieBiblioteki(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz udostępnienia: %w", err)
		}
		lista = append(lista, udostepnienie)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt udostępnień: %w", err)
	}
	return lista, nil
}

// OdwolajUdostepnienie znakuje udostępnienie jako odwołane. Powtórne odwołanie
// tego samego wiersza oddaje fałsz — nie było czego odwoływać.
func (r *repozytoriumBiblioteki) OdwolajUdostepnienie(ctx context.Context, kod string) (bool, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, odwolajUdostepnienieBiblioteki)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, kod, KontoOperatora(ctx))
	if err != nil {
		return false, fmt.Errorf("dane: nie można odwołać udostępnienia %q: %w", kod, err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nie można policzyć odwołanych udostępnień: %w", err)
	}
	return zmienione > 0, nil
}

// ZapiszWebhook zakłada nasłuch albo zmienia zastany wraz z kompletem zdarzeń;
// kod bez adresu jest odrzucany jako błąd.
func (r *repozytoriumBiblioteki) ZapiszWebhook(ctx context.Context,
	webhook WebhookBiblioteki) (WebhookBiblioteki, error) {

	if strings.TrimSpace(webhook.Kod) == "" {
		return WebhookBiblioteki{}, fmt.Errorf("dane: nasłuch biblioteki bez identyfikatora")
	}
	if strings.TrimSpace(webhook.Adres) == "" {
		return WebhookBiblioteki{}, fmt.Errorf("dane: nasłuch biblioteki %q bez adresu", webhook.Kod)
	}
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		zapis, err := r.zapytania.wTransakcji(ctx, transakcja, zapiszWebhookBiblioteki)
		if err != nil {
			return err
		}
		wynik, err := zapis.ExecContext(ctx, webhook.Kod, webhook.Adres,
			tekstDoKolumny(webhook.Sekret), liczbaLogiczna(webhook.Czynny),
			tekstDoKolumny(webhook.OstatnieZgloszenie), KontoOperatora(ctx),
			KontoOperatora(ctx))
		if err != nil {
			return fmt.Errorf("dane: nie można zapisać nasłuchu %q: %w", webhook.Kod, err)
		}
		zmienione, err := wynik.RowsAffected()
		if err != nil {
			return fmt.Errorf("dane: nieznana liczba zapisanych nasłuchów biblioteki: %w", err)
		}
		if zmienione == 0 {
			return fmt.Errorf("dane: nasłuch %q należy do innego konta: %w",
				webhook.Kod, ErrKolizjaWiersza)
		}
		odczyt, err := r.zapytania.wTransakcji(ctx, transakcja, pobierzWebhookBiblioteki)
		if err != nil {
			return err
		}
		zapisany, err := odczytajWebhookBiblioteki(
			odczyt.QueryRowContext(ctx, webhook.Kod, KontoOperatora(ctx)))
		if err != nil {
			return fmt.Errorf("dane: nieczytelny nasłuch %q: %w", webhook.Kod, err)
		}
		czyszczenie, err := r.zapytania.wTransakcji(ctx, transakcja, usunZdarzeniaWebhookaBiblioteki)
		if err != nil {
			return err
		}
		if _, err := czyszczenie.ExecContext(ctx, zapisany.ID); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić zdarzeń nasłuchu %q: %w", webhook.Kod, err)
		}
		wstawienie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawZdarzenieWebhookaBiblioteki)
		if err != nil {
			return err
		}
		for _, zdarzenie := range webhook.Zdarzenia {
			if zdarzenie == "" {
				continue
			}
			if _, err := wstawienie.ExecContext(ctx, zapisany.ID, zdarzenie); err != nil {
				return fmt.Errorf("dane: nie można zapisać zdarzenia %q nasłuchu %q: %w",
					zdarzenie, webhook.Kod, err)
			}
		}
		return nil
	})
	if err != nil {
		return WebhookBiblioteki{}, err
	}
	return r.webhookPoKodzie(ctx, webhook.Kod)
}

// Webhooki zwraca nasłuchy od najnowszego wraz z ich zdarzeniami, opcjonalnie
// zawężone do wierszy czynnych.
func (r *repozytoriumBiblioteki) Webhooki(ctx context.Context, tylkoCzynne bool) ([]WebhookBiblioteki, error) {
	warunek := WarunekKonta
	if tylkoCzynne {
		warunek += " AND czynny = 1"
	}
	zapytanie := `SELECT ` + kolumnyWebhookaBiblioteki + ` FROM webhook_biblioteki
	              WHERE ` + warunek + ` ORDER BY utworzono DESC, id DESC`

	wiersze, err := r.db.QueryContext(ctx, zapytanie, KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać nasłuchów biblioteki: %w", err)
	}
	defer wiersze.Close()

	lista := []WebhookBiblioteki{}
	for wiersze.Next() {
		webhook, err := odczytajWebhookBiblioteki(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz nasłuchu: %w", err)
		}
		lista = append(lista, webhook)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt nasłuchów: %w", err)
	}
	for indeks := range lista {
		zdarzenia, err := r.zdarzeniaWebhooka(ctx, lista[indeks].ID)
		if err != nil {
			return nil, err
		}
		lista[indeks].Zdarzenia = zdarzenia
	}
	return lista, nil
}

// UsunWebhook zdejmuje nasłuch wraz z jego zdarzeniami przez kaskadę schematu
// i oddaje fałsz, gdy nasłuch o podanym kodzie nie istniał.
func (r *repozytoriumBiblioteki) UsunWebhook(ctx context.Context, kod string) (bool, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, usunWebhookBiblioteki)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, kod, KontoOperatora(ctx))
	if err != nil {
		return false, fmt.Errorf("dane: nie można usunąć nasłuchu %q: %w", kod, err)
	}
	zdjete, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nie można policzyć usuniętych nasłuchów: %w", err)
	}
	return zdjete > 0, nil
}

// webhookPoKodzie odczytuje nasłuch wraz ze zdarzeniami po kodzie zewnętrznym,
// oddając ErrBrakWiersza, gdy nasłuch nie istnieje.
func (r *repozytoriumBiblioteki) webhookPoKodzie(ctx context.Context, kod string) (WebhookBiblioteki, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWebhookBiblioteki)
	if err != nil {
		return WebhookBiblioteki{}, err
	}
	webhook, err := odczytajWebhookBiblioteki(
		polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return WebhookBiblioteki{}, ErrBrakWiersza
	}
	if err != nil {
		return WebhookBiblioteki{}, fmt.Errorf("dane: nieczytelny nasłuch %q: %w", kod, err)
	}
	zdarzenia, err := r.zdarzeniaWebhooka(ctx, webhook.ID)
	if err != nil {
		return WebhookBiblioteki{}, err
	}
	webhook.Zdarzenia = zdarzenia
	return webhook, nil
}

// zdarzeniaWebhooka odczytuje zdarzenia jednego nasłuchu, uporządkowane
// alfabetycznie po nazwie zdarzenia.
func (r *repozytoriumBiblioteki) zdarzeniaWebhooka(ctx context.Context, webhookID int64) ([]string, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaZdarzenWebhookaBiblioteki)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, webhookID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać zdarzeń nasłuchu %d: %w", webhookID, err)
	}
	defer wiersze.Close()

	lista := []string{}
	for wiersze.Next() {
		var zdarzenie string
		if err := wiersze.Scan(&zdarzenie); err != nil {
			return nil, fmt.Errorf("dane: nieczytelne zdarzenie nasłuchu %d: %w", webhookID, err)
		}
		lista = append(lista, zdarzenie)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt zdarzeń nasłuchu %d: %w", webhookID, err)
	}
	return lista, nil
}

// odczytajUdostepnienieBiblioteki składa udostępnienie z jednego wiersza
// wyniku, w kolejności kolumn kolumnyUdostepnieniaBiblioteki.
func odczytajUdostepnienieBiblioteki(wiersz skaner) (UdostepnienieBiblioteki, error) {
	var udostepnienie UdostepnienieBiblioteki
	var wygasa, odwolano sql.NullString
	err := wiersz.Scan(&udostepnienie.ID, &udostepnienie.Kod, &udostepnienie.Zasieg,
		&udostepnienie.CelKod, &udostepnienie.Token, &wygasa, &odwolano, &udostepnienie.Utworzono)
	if err != nil {
		return UdostepnienieBiblioteki{}, err
	}
	udostepnienie.Wygasa, udostepnienie.Odwolano = tekstZKolumny(wygasa), tekstZKolumny(odwolano)
	return udostepnienie, nil
}

// odczytajWebhookBiblioteki składa nasłuch z jednego wiersza wyniku, w kolejności
// kolumn kolumnyWebhookaBiblioteki; zdarzenia dokłada strona wołająca.
func odczytajWebhookBiblioteki(wiersz skaner) (WebhookBiblioteki, error) {
	var webhook WebhookBiblioteki
	var sekret, zgloszenie sql.NullString
	var czynny int
	err := wiersz.Scan(&webhook.ID, &webhook.Kod, &webhook.Adres, &sekret, &czynny,
		&zgloszenie, &webhook.Utworzono)
	if err != nil {
		return WebhookBiblioteki{}, err
	}
	webhook.Sekret, webhook.OstatnieZgloszenie = tekstZKolumny(sekret), tekstZKolumny(zgloszenie)
	webhook.Czynny = czynny == 1
	webhook.Zdarzenia = []string{}
	return webhook, nil
}
