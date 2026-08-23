// Odpowiedzialność pliku: źródła zebrane w toku przeglądania (tabela
// `zrodlo_przegladania`) — trwałość szuflady Sources.
//
// Tabela `zrodlo_przegladania` nie jest tabelą `zrodlo_badania` modułu Research.
// `BrowserSource` to odcisk strony zebrany w toku przeglądania — ma `snapshotId`
// i `key`, jest zawsze powiązany z oknem operacyjnym. `ResearchSource` ocenia
// wiarygodność zasobu badawczego i niesie inny kształt (`kind`, `credibility`,
// `libraryFileId`). Różne kształty, różne cykle życia — osobna tabela, nie
// współdzielenie.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// ZrodloPrzegladania to wiersz tabeli `zrodlo_przegladania`.
type ZrodloPrzegladania struct {
	ID                  int64
	Kod                 string
	Okno                string
	Url                 string
	Tytul               *string
	MigawkaZewnetrznaID *string
	Kluczowe            bool
	// Grupa niesie zestaw tematyczny, do którego źródło należy
	// (`browser.source.group.set`, migracja 179). Pusty wskaźnik znaczy „poza
	// zestawami", nie „zestaw bez nazwy".
	Grupa     *string
	Utworzono string
}

const (
	kolumnyZrodlaPrzegladania = `id, identyfikator_zewnetrzny, okno, url, tytul,
	                             migawka_zewnetrzna_id, kluczowe, grupa, utworzono`

	zapiszZrodloPrzegladania = `INSERT INTO zrodlo_przegladania
	                            (identyfikator_zewnetrzny, okno, url, tytul,
	                             migawka_zewnetrzna_id, kluczowe, grupa)
	                            VALUES (?, ?, ?, ?, ?, ?, ?)
	                            ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                                okno = excluded.okno,
	                                url = excluded.url,
	                                tytul = excluded.tytul,
	                                migawka_zewnetrzna_id = excluded.migawka_zewnetrzna_id,
	                                kluczowe = excluded.kluczowe,
	                                grupa = excluded.grupa`

	pobierzZrodloPrzegladania = `SELECT ` + kolumnyZrodlaPrzegladania + `
	                             FROM zrodlo_przegladania WHERE identyfikator_zewnetrzny = ?`

	// Jedno zapytanie na wszystkie zawężenia zamiast sklejania SQL-a w locie:
	// zerowy `tylkoKluczowe` wyłącza warunek kluczowości, więc plan i wpis
	// w podręcznej pamięci `przygotuj` są jedne, niezależnie od filtru.
	// Kolejność zgodna z indeksem idx_zrodlo_przegladania_okno (okno, utworzono DESC, id).
	listaZrodelPrzegladania = `SELECT ` + kolumnyZrodlaPrzegladania + `
	                           FROM zrodlo_przegladania
	                           WHERE okno = ? AND (? = 0 OR kluczowe = 1)
	                             AND (? = '' OR grupa = ?)
	                             AND (? = '' OR url LIKE ? OR IFNULL(tytul,'') LIKE ?)
	                           ORDER BY utworzono DESC, id DESC LIMIT ?`
)

// FiltrZrodelPrzegladania zawęża wykaz źródeł okna — obsługuje pola żądania
// `browser.source.list`. Filtr zamiast trzech argumentów, bo kontrakt może
// dołożyć kolejne zawężenie, a wtedy zmieniałby się każdy wołający.
type FiltrZrodelPrzegladania struct {
	// Okno operacyjne, którego wykaz dotyczy — pole wymagane kontraktem.
	Okno string
	// TylkoKluczowe odsiewa źródła nieoznaczone jako kluczowe.
	TylkoKluczowe bool
	// Zestaw zawęża wykaz do jednego zestawu tematycznego źródeł.
	Zestaw string
	// Szukaj przegląda adres i tytuł naraz — Operator pamięta jedno z dwojga.
	Szukaj string
	// Limit 0 lub ujemny znaczy wykaz pełny, nie wykaz pusty.
	Limit int
}

// ZapiszZrodlo zakłada źródło albo nadpisuje zastane (dopasowane po kodzie
// zewnętrznym) i zwraca stan po zapisie — `browser.source.add` jest
// idempotentne wobec ponownego wywołania z tym samym `Id`.
func (r *repozytoriumPrzegladania) ZapiszZrodlo(ctx context.Context, zrodlo ZrodloPrzegladania) (ZrodloPrzegladania, error) {
	if zrodlo.Kod == "" || zrodlo.Okno == "" || zrodlo.Url == "" {
		return ZrodloPrzegladania{}, fmt.Errorf("dane: źródło przeglądania bez identyfikatora, okna albo adresu")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszZrodloPrzegladania)
	if err != nil {
		return ZrodloPrzegladania{}, err
	}
	_, err = polecenie.ExecContext(ctx, zrodlo.Kod, zrodlo.Okno, zrodlo.Url,
		tekstDoKolumny(zrodlo.Tytul), tekstDoKolumny(zrodlo.MigawkaZewnetrznaID),
		liczbaLogiczna(zrodlo.Kluczowe), tekstDoKolumny(zrodlo.Grupa))
	if err != nil {
		return ZrodloPrzegladania{}, fmt.Errorf("dane: nie można zapisać źródła przeglądania %q: %w", zrodlo.Kod, err)
	}
	return r.jednoZrodlo(ctx, zrodlo.Kod)
}

// jednoZrodlo odczytuje pojedynczy wiersz źródła po kodzie zewnętrznym.
func (r *repozytoriumPrzegladania) jednoZrodlo(ctx context.Context, kod string) (ZrodloPrzegladania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzZrodloPrzegladania)
	if err != nil {
		return ZrodloPrzegladania{}, err
	}
	zrodlo, err := odczytajZrodloPrzegladania(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return ZrodloPrzegladania{}, ErrBrakWiersza
	}
	if err != nil {
		return ZrodloPrzegladania{}, fmt.Errorf("dane: nieczytelne źródło przeglądania %q: %w", kod, err)
	}
	return zrodlo, nil
}

// Zrodla zwraca źródła okna od najnowszego — szuflada Sources pokazuje to, co
// zebrano w danym oknie operacyjnym, ewentualnie zawężone filtrem. Wykaz pusty
// jest prawidłowym wynikiem, nie brakiem wiersza: okno bez źródeł nie jest
// oknem nieznanym (rozróżnienie robi `OknoZnane`).
func (r *repozytoriumPrzegladania) Zrodla(ctx context.Context,
	filtr FiltrZrodelPrzegladania) ([]ZrodloPrzegladania, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaZrodelPrzegladania)
	if err != nil {
		return nil, err
	}
	okno := filtr.Okno
	wzorzec := ""
	if filtr.Szukaj != "" {
		wzorzec = "%" + filtr.Szukaj + "%"
	}
	wiersze, err := polecenie.QueryContext(ctx, okno,
		liczbaLogiczna(filtr.TylkoKluczowe), filtr.Zestaw, filtr.Zestaw,
		filtr.Szukaj, wzorzec, wzorzec, granicaWykazu(filtr.Limit))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać źródeł przeglądania okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	lista := []ZrodloPrzegladania{}
	for wiersze.Next() {
		zrodlo, err := odczytajZrodloPrzegladania(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz źródła przeglądania: %w", err)
		}
		lista = append(lista, zrodlo)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt źródeł przeglądania: %w", err)
	}
	return lista, nil
}

// odczytajZrodloPrzegladania składa strukturę z jednego wiersza wyniku.
func odczytajZrodloPrzegladania(wiersz skaner) (ZrodloPrzegladania, error) {
	var zrodlo ZrodloPrzegladania
	var tytul, migawkaZewnetrznaID, grupa sql.NullString
	var kluczowe int
	err := wiersz.Scan(&zrodlo.ID, &zrodlo.Kod, &zrodlo.Okno, &zrodlo.Url, &tytul,
		&migawkaZewnetrznaID, &kluczowe, &grupa, &zrodlo.Utworzono)
	if err != nil {
		return ZrodloPrzegladania{}, err
	}
	zrodlo.Tytul = tekstZKolumny(tytul)
	zrodlo.MigawkaZewnetrznaID = tekstZKolumny(migawkaZewnetrznaID)
	zrodlo.Kluczowe = kluczowe != 0
	zrodlo.Grupa = tekstZKolumny(grupa)
	return zrodlo, nil
}
