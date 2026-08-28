// Odpowiedzialność pliku: tożsamość wykonawcy przy zmianie śledzonej modułu Studio, czyli kolumny dołożone tabeli zmiana_sledzona_studio.
package dane

import (
	"context"
	"database/sql"
	"fmt"
)

// ZmianaWykonawcyStudia to wiersz `zmiana_sledzona_studio` widziany razem
// z kolumnami tożsamości wykonawcy i postaci sprzed oraz po zmianie.
type ZmianaWykonawcyStudia struct {
	ID               int64
	Kod              string
	DokumentKod      string
	Rodzaj           string
	Autor            string
	AutorAgentKod    *string
	AutorAgentNazwa  *string
	AutorAgentWersja *string
	AutorPodagentKod *string
	ZakresOd         int64
	ZakresDo         int64
	TrescPrzed       *string
	TrescPo          *string
	PostacPrzedJSON  *string
	PostacPoJSON     *string
	CzynnoscKod      *string
	Decyzja          string
	Utworzono        string
}

const (
	kolumnyZmianyWykonawcy = `z.id, z.identyfikator_zewnetrzny, d.identyfikator_zewnetrzny,
	                          z.rodzaj, z.autor, z.autor_agent_kod, z.autor_agent_nazwa,
	                          z.autor_agent_wersja, z.autor_podagent_kod, z.zakres_od,
	                          z.zakres_do, z.tresc_przed, z.tresc_po, z.postac_przed_json,
	                          z.postac_po_json, z.czynnosc_kod, z.decyzja, z.utworzono`

	zrodloZmianyWykonawcy = ` FROM zmiana_sledzona_studio z
	                          JOIN dokument_studio d ON d.id = z.dokument_id`

	// Kolejność wystąpienia w treści, bo podświetlenie zmian wykonawcy jest
	// spisem do przejścia: „następna zmiana modelu" znaczy następna w dokumencie,
	// a nie następna w czasie.
	zmianyWykonawcowLista = `SELECT ` + kolumnyZmianyWykonawcy + zrodloZmianyWykonawcy +
		` WHERE z.dokument_id = ? ORDER BY z.zakres_od, z.id`

	// Tożsamość stempluje się na wierszu już istniejącym, bo zapis zmiany
	// śledzonej należy do innego pliku i nie wolno mu podstawiać drugiej drogi
	// zakładania wiersza.
	zmianaWykonawcyStempluj = `UPDATE zmiana_sledzona_studio SET
	                               autor_agent_kod = ?, autor_agent_nazwa = ?,
	                               autor_agent_wersja = ?, autor_podagent_kod = ?,
	                               postac_przed_json = COALESCE(?, postac_przed_json),
	                               postac_po_json = COALESCE(?, postac_po_json),
	                               czynnosc_kod = COALESCE(?, czynnosc_kod)
	                           WHERE identyfikator_zewnetrzny = ?`
)

// ZmianyWykonawcow oddaje zmiany śledzone dokumentu wraz z tożsamością
// wykonawcy, w kolejności wystąpienia w treści.
func (r *repozytoriumStudia) ZmianyWykonawcow(ctx context.Context,
	dokumentID int64) ([]ZmianaWykonawcyStudia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zmianyWykonawcowLista)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, dokumentID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać zmian śledzonych wykonawców: %w", err)
	}
	defer wiersze.Close()

	lista := []ZmianaWykonawcyStudia{}
	for wiersze.Next() {
		zmiana, err := odczytajZmianeWykonawcyStudia(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz zmiany śledzonej: %w", err)
		}
		lista = append(lista, zmiana)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt zmian śledzonych: %w", err)
	}
	return lista, nil
}

// StemplujTozsamoscZmiany dopisuje do zmiany śledzonej, który wykonawca ją wniósł, wartościami niezerującymi kolumn już wypełnionych.
func (r *repozytoriumStudia) StemplujTozsamoscZmiany(ctx context.Context, kod string,
	agentKod, agentNazwa, agentWersja, podagentKod *string,
	postacPrzed, postacPo, czynnoscKod *string) error {

	if kod == "" {
		return fmt.Errorf("dane: stemplowanie tożsamości bez wskazania zmiany śledzonej")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zmianaWykonawcyStempluj)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, tekstDoKolumny(agentKod), tekstDoKolumny(agentNazwa),
		tekstDoKolumny(agentWersja), tekstDoKolumny(podagentKod), tekstDoKolumny(postacPrzed),
		tekstDoKolumny(postacPo), tekstDoKolumny(czynnoscKod), kod)
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać tożsamości wykonawcy zmiany %q: %w", kod, err)
	}
	return nil
}

func odczytajZmianeWykonawcyStudia(wiersz skaner) (ZmianaWykonawcyStudia, error) {
	var zmiana ZmianaWykonawcyStudia
	var agentKod, agentNazwa, agentWersja, podagent sql.NullString
	var przed, po, postacPrzed, postacPo, czynnosc sql.NullString
	err := wiersz.Scan(&zmiana.ID, &zmiana.Kod, &zmiana.DokumentKod, &zmiana.Rodzaj,
		&zmiana.Autor, &agentKod, &agentNazwa, &agentWersja, &podagent, &zmiana.ZakresOd,
		&zmiana.ZakresDo, &przed, &po, &postacPrzed, &postacPo, &czynnosc,
		&zmiana.Decyzja, &zmiana.Utworzono)
	if err != nil {
		return ZmianaWykonawcyStudia{}, err
	}
	zmiana.AutorAgentKod = tekstZKolumny(agentKod)
	zmiana.AutorAgentNazwa = tekstZKolumny(agentNazwa)
	zmiana.AutorAgentWersja = tekstZKolumny(agentWersja)
	zmiana.AutorPodagentKod = tekstZKolumny(podagent)
	zmiana.TrescPrzed = tekstZKolumny(przed)
	zmiana.TrescPo = tekstZKolumny(po)
	zmiana.PostacPrzedJSON = tekstZKolumny(postacPrzed)
	zmiana.PostacPoJSON = tekstZKolumny(postacPo)
	zmiana.CzynnoscKod = tekstZKolumny(czynnosc)
	return zmiana, nil
}
