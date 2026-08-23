// Odpowiedzialność pliku: warstwa danych pochodzenia fragmentów dokumentu
// Studia — tabela `pochodzenie_fragmentu_studio` z migracji 368.
//
// ── Dlaczego pochodzenie zapisuje odcinek wejścia ───────────────────────────
// Wiersz pochodzenia zakłada ten, kto fragment WNOSI: wniesienie pliku, obrazu,
// fragmentu z Biblioteki albo ze strony sieci. Wykaz pochodzenia (`studio.
// provenance.list`) czyta go inny odcinek — dlatego odczyt jest tu równie pełny
// jak zapis, żeby tamten odcinek nie musiał zakładać drugiego.
//
// Zakres jest liczony w ZNAKACH, tak samo jak zakres blokady. Fragment usunięty
// zostawia wiersz o zerowej długości — świadomie, wedle migracji 368: to, że
// Operator wniósł kiedyś fragment z danego źródła, jest faktem, którego
// usunięcie tekstu nie unieważnia.
package dane

import (
	"context"
	"database/sql"
	"fmt"
)

// PochodzenieFragmentuStudia to wiersz tabeli `pochodzenie_fragmentu_studio`.
type PochodzenieFragmentuStudia struct {
	ID         int64
	Kod        string
	DokumentID int64
	// Rodzaj jest słownikiem tabeli: `web`, `libraryFile`, `research`,
	// `clipboard`, `template`, `importedFile`.
	Rodzaj            string
	ZakresOd          int64
	ZakresDo          int64
	AdresZrodla       *string
	BibliotekaPlikKod *string
	WersjaZrodla      *string
	TytulZrodla       *string
	Siegnieto         *string
	AutorRodzaj       string
	AutorAgentKod     *string
	AutorAgentNazwa   *string
	Utworzono         string
}

// PochodzenieWejsciaStudia jest kontraktem tej warstwy.
type PochodzenieWejsciaStudia interface {
	ZapiszPochodzenieFragmentu(ctx context.Context,
		pochodzenie PochodzenieFragmentuStudia) (PochodzenieFragmentuStudia, error)
	PochodzenieFragmentow(ctx context.Context,
		dokumentID int64) ([]PochodzenieFragmentuStudia, error)
}

const (
	wejscieKolumnyPochodzenia = `id, identyfikator_zewnetrzny, dokument_id, rodzaj, zakres_od,
	                             zakres_do, adres_zrodla, biblioteka_plik_kod, wersja_zrodla,
	                             tytul_zrodla, siegnieto, autor_rodzaj, autor_agent_kod,
	                             autor_agent_nazwa, utworzono`

	wejscieZapiszPochodzenie = `INSERT INTO pochodzenie_fragmentu_studio
	                            (identyfikator_zewnetrzny, dokument_id, rodzaj, zakres_od,
	                             zakres_do, adres_zrodla, biblioteka_plik_kod, wersja_zrodla,
	                             tytul_zrodla, siegnieto, autor_rodzaj, autor_agent_kod,
	                             autor_agent_nazwa)
	                            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	                            ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                                zakres_od = excluded.zakres_od,
	                                zakres_do = excluded.zakres_do,
	                                adres_zrodla = excluded.adres_zrodla,
	                                tytul_zrodla = excluded.tytul_zrodla`

	wejsciePobierzPochodzenie = `SELECT ` + wejscieKolumnyPochodzenia + `
	                             FROM pochodzenie_fragmentu_studio
	                             WHERE identyfikator_zewnetrzny = ?`

	wejscieListaPochodzenia = `SELECT ` + wejscieKolumnyPochodzenia + `
	                           FROM pochodzenie_fragmentu_studio
	                           WHERE dokument_id = ?
	                           ORDER BY zakres_od, id`
)

// ZapiszPochodzenieFragmentu utrwala zapis, skąd fragment dokumentu pochodzi.
func (r *repozytoriumStudia) ZapiszPochodzenieFragmentu(ctx context.Context,
	pochodzenie PochodzenieFragmentuStudia) (PochodzenieFragmentuStudia, error) {

	if pochodzenie.Kod == "" || pochodzenie.DokumentID == 0 {
		return PochodzenieFragmentuStudia{}, fmt.Errorf(
			"dane: pochodzenie fragmentu bez identyfikatora albo bez dokumentu")
	}
	if pochodzenie.AutorRodzaj == "" {
		pochodzenie.AutorRodzaj = "uzytkownik"
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wejscieZapiszPochodzenie)
	if err != nil {
		return PochodzenieFragmentuStudia{}, err
	}
	_, err = polecenie.ExecContext(ctx, pochodzenie.Kod, pochodzenie.DokumentID,
		pochodzenie.Rodzaj, pochodzenie.ZakresOd, pochodzenie.ZakresDo,
		tekstDoKolumny(pochodzenie.AdresZrodla), tekstDoKolumny(pochodzenie.BibliotekaPlikKod),
		tekstDoKolumny(pochodzenie.WersjaZrodla), tekstDoKolumny(pochodzenie.TytulZrodla),
		tekstDoKolumny(pochodzenie.Siegnieto), pochodzenie.AutorRodzaj,
		tekstDoKolumny(pochodzenie.AutorAgentKod), tekstDoKolumny(pochodzenie.AutorAgentNazwa))
	if err != nil {
		return PochodzenieFragmentuStudia{}, fmt.Errorf(
			"dane: nie można zapisać pochodzenia fragmentu %q: %w", pochodzenie.Kod, err)
	}
	polecenie, err = r.zapytania.przygotuj(ctx, wejsciePobierzPochodzenie)
	if err != nil {
		return PochodzenieFragmentuStudia{}, err
	}
	zapisane, err := wejscieOdczytajPochodzenie(polecenie.QueryRowContext(ctx, pochodzenie.Kod))
	if err != nil {
		return PochodzenieFragmentuStudia{}, fmt.Errorf(
			"dane: nieczytelne pochodzenie fragmentu %q: %w", pochodzenie.Kod, err)
	}
	return zapisane, nil
}

// PochodzenieFragmentow oddaje wykaz pochodzenia fragmentów dokumentu.
func (r *repozytoriumStudia) PochodzenieFragmentow(ctx context.Context,
	dokumentID int64) ([]PochodzenieFragmentuStudia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, wejscieListaPochodzenia)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, dokumentID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać pochodzenia fragmentów: %w", err)
	}
	defer wiersze.Close()

	lista := []PochodzenieFragmentuStudia{}
	for wiersze.Next() {
		pochodzenie, err := wejscieOdczytajPochodzenie(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz pochodzenia fragmentu: %w", err)
		}
		lista = append(lista, pochodzenie)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt pochodzenia fragmentów: %w", err)
	}
	return lista, nil
}

func wejscieOdczytajPochodzenie(wiersz skaner) (PochodzenieFragmentuStudia, error) {
	var pochodzenie PochodzenieFragmentuStudia
	var adres, plik, wersja, tytul, siegnieto, agentKod, agentNazwa sql.NullString
	err := wiersz.Scan(&pochodzenie.ID, &pochodzenie.Kod, &pochodzenie.DokumentID,
		&pochodzenie.Rodzaj, &pochodzenie.ZakresOd, &pochodzenie.ZakresDo, &adres, &plik,
		&wersja, &tytul, &siegnieto, &pochodzenie.AutorRodzaj, &agentKod, &agentNazwa,
		&pochodzenie.Utworzono)
	if err != nil {
		return PochodzenieFragmentuStudia{}, err
	}
	pochodzenie.AdresZrodla = tekstZKolumny(adres)
	pochodzenie.BibliotekaPlikKod = tekstZKolumny(plik)
	pochodzenie.WersjaZrodla = tekstZKolumny(wersja)
	pochodzenie.TytulZrodla = tekstZKolumny(tytul)
	pochodzenie.Siegnieto = tekstZKolumny(siegnieto)
	pochodzenie.AutorAgentKod = tekstZKolumny(agentKod)
	pochodzenie.AutorAgentNazwa = tekstZKolumny(agentNazwa)
	return pochodzenie, nil
}
