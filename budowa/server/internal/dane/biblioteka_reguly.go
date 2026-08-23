// Odpowiedzialność pliku: reguły repozytorium (`regula_biblioteki`, migracja
// 182) oraz wykaz kolekcji w postaci pełnej — z hierarchią, regułą i licznikiem
// zasobów (`library.collection.list`).
//
// Kolekcje mają tu drugie wejście obok `library_kolekcje.go` i to nie jest
// powielenie: tamten plik odpowiada za założenie kolekcji i przypisanie zasobów,
// ten — za odczyt kolekcji jako bytu opisanego (rodzic, reguła, licznik). Podział
// idzie po pytaniu, nie po tabeli.
//
// Reguła nie wykonuje się sama. Warstwa danych przechowuje warunek i wskazuje
// kolekcję docelową; przeliczenie — czyli zamiana warunku na wykaz zasobów —
// należy do rdzenia, bo to on zna znaczenie członów warunku.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// RegulaBiblioteki to wiersz tabeli `regula_biblioteki`.
type RegulaBiblioteki struct {
	ID                   int64
	Kod                  string
	Rodzaj               string
	Nazwa                string
	Warunek              string
	KolekcjaDocelowaKod  *string
	SciezkaObserwowana   *string
	Czynna               bool
	OstatniePrzeliczenie *string
	Utworzono            string
}

const (
	kolumnyRegulyBiblioteki = `id, identyfikator_zewnetrzny, rodzaj, nazwa, warunek, kolekcja_docelowa_kod,
	                 sciezka_obserwowana, czynna, ostatnie_przeliczenie, utworzono`

	zapiszReguleBiblioteki = `INSERT INTO regula_biblioteki
	                          (identyfikator_zewnetrzny, rodzaj, nazwa, warunek,
	                           kolekcja_docelowa_kod, sciezka_obserwowana, czynna,
	                           ostatnie_przeliczenie)
	                          VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	                          ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                              rodzaj = excluded.rodzaj,
	                              nazwa = excluded.nazwa,
	                              warunek = excluded.warunek,
	                              kolekcja_docelowa_kod = excluded.kolekcja_docelowa_kod,
	                              sciezka_obserwowana = excluded.sciezka_obserwowana,
	                              czynna = excluded.czynna,
	                              ostatnie_przeliczenie = excluded.ostatnie_przeliczenie`

	pobierzReguleBiblioteki = `SELECT ` + kolumnyRegulyBiblioteki + ` FROM regula_biblioteki
	                 WHERE identyfikator_zewnetrzny = ?`

	usunReguleBiblioteki = `DELETE FROM regula_biblioteki WHERE identyfikator_zewnetrzny = ?`

	// Wykaz kolekcji niesie licznik zasobów liczony z przypisań oraz kod
	// kolekcji nadrzędnej — kontrakt wskazuje rodzica kodem, nie numerem wiersza.
	wykazKolekcjiBiblioteki = `SELECT k.id, k.identyfikator_zewnetrzny, k.nazwa, k.opis,
	                                  k.utworzono, k.zaktualizowano,
	                                  (SELECT r.identyfikator_zewnetrzny FROM kolekcja_biblioteki r
	                                    WHERE r.id = k.rodzic_id) AS rodzic_kod,
	                                  k.regula_kod,
	                                  (SELECT COUNT(*) FROM przypisanie_kolekcji_biblioteki p
	                                    WHERE p.kolekcja_id = k.id) AS liczba_plikow
	                           FROM kolekcja_biblioteki k`

	wstawPrzypisanieRegulyBiblioteki = `INSERT INTO przypisanie_kolekcji_biblioteki (kolekcja_id, plik_id, zrodlo)
	                          VALUES (?, ?, 'regula')
	                          ON CONFLICT(kolekcja_id, plik_id) DO NOTHING`

	usunPrzypisaniaRegulyBiblioteki = `DELETE FROM przypisanie_kolekcji_biblioteki
	                         WHERE zrodlo = 'regula' AND kolekcja_id =
	                             (SELECT id FROM kolekcja_biblioteki WHERE identyfikator_zewnetrzny = ?)`
)

// ZapiszRegule zakłada regułę albo nadpisuje zastaną.
func (r *repozytoriumBiblioteki) ZapiszRegule(ctx context.Context,
	regula RegulaBiblioteki) (RegulaBiblioteki, error) {

	if strings.TrimSpace(regula.Kod) == "" {
		return RegulaBiblioteki{}, fmt.Errorf("dane: reguła biblioteki bez identyfikatora")
	}
	if strings.TrimSpace(regula.Nazwa) == "" {
		return RegulaBiblioteki{}, fmt.Errorf("dane: reguła biblioteki %q bez nazwy", regula.Kod)
	}
	warunek := regula.Warunek
	if strings.TrimSpace(warunek) == "" {
		warunek = "{}"
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszReguleBiblioteki)
	if err != nil {
		return RegulaBiblioteki{}, err
	}
	_, err = polecenie.ExecContext(ctx, regula.Kod, regula.Rodzaj, regula.Nazwa, warunek,
		tekstDoKolumny(regula.KolekcjaDocelowaKod), tekstDoKolumny(regula.SciezkaObserwowana),
		liczbaLogiczna(regula.Czynna), tekstDoKolumny(regula.OstatniePrzeliczenie))
	if err != nil {
		return RegulaBiblioteki{}, fmt.Errorf("dane: nie można zapisać reguły %q: %w", regula.Kod, err)
	}
	return r.Regula(ctx, regula.Kod)
}

// Regula zwraca regułę o wskazanym kodzie; brak wiersza wraca jako ErrBrakWiersza.
func (r *repozytoriumBiblioteki) Regula(ctx context.Context, kod string) (RegulaBiblioteki, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzReguleBiblioteki)
	if err != nil {
		return RegulaBiblioteki{}, err
	}
	regula, err := odczytajReguleBiblioteki(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return RegulaBiblioteki{}, ErrBrakWiersza
	}
	if err != nil {
		return RegulaBiblioteki{}, fmt.Errorf("dane: nieczytelna reguła biblioteki %q: %w", kod, err)
	}
	return regula, nil
}

// Reguly zwraca reguły uporządkowane po nazwie, zawężone rodzajem i czynnością.
func (r *repozytoriumBiblioteki) Reguly(ctx context.Context, rodzaj *string,
	tylkoCzynne bool) ([]RegulaBiblioteki, error) {

	warunki := []string{"1 = 1"}
	argumenty := []any{}
	if rodzaj != nil && *rodzaj != "" {
		warunki = append(warunki, "rodzaj = ?")
		argumenty = append(argumenty, *rodzaj)
	}
	if tylkoCzynne {
		warunki = append(warunki, "czynna = 1")
	}
	zapytanie := `SELECT ` + kolumnyRegulyBiblioteki + ` FROM regula_biblioteki
	              WHERE ` + strings.Join(warunki, " AND ") + ` ORDER BY nazwa, id`

	wiersze, err := r.db.QueryContext(ctx, zapytanie, argumenty...)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać reguł biblioteki: %w", err)
	}
	defer wiersze.Close()

	lista := []RegulaBiblioteki{}
	for wiersze.Next() {
		regula, err := odczytajReguleBiblioteki(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz reguły biblioteki: %w", err)
		}
		lista = append(lista, regula)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt reguł biblioteki: %w", err)
	}
	return lista, nil
}

// UsunRegule zdejmuje regułę. Przypisania zasobów zostają — zdejmuje je osobno
// `OdepnijPrzypisaniaReguly`, bo kontrakt czyni to wyborem Operatora.
func (r *repozytoriumBiblioteki) UsunRegule(ctx context.Context, kod string) (bool, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, usunReguleBiblioteki)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, kod)
	if err != nil {
		return false, fmt.Errorf("dane: nie można usunąć reguły %q: %w", kod, err)
	}
	zdjete, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nie można policzyć usuniętych reguł: %w", err)
	}
	return zdjete > 0, nil
}

// PrzypiszRegula wpisuje zasoby do kolekcji ze znacznikiem pochodzenia `regula`.
func (r *repozytoriumBiblioteki) PrzypiszRegula(ctx context.Context, kodKolekcji string,
	kodyPlikow []string) (int, error) {

	kolekcja, err := r.kolekcjaPoKodzie(ctx, kodKolekcji)
	if err != nil {
		return 0, err
	}
	poszukiwanie, err := r.zapytania.przygotuj(ctx, idPlikuBibliotekiPoKodzie)
	if err != nil {
		return 0, err
	}
	wstawianie, err := r.zapytania.przygotuj(ctx, wstawPrzypisanieRegulyBiblioteki)
	if err != nil {
		return 0, err
	}
	przypisane := 0
	for _, kodPliku := range kodyPlikow {
		var plikID int64
		err := poszukiwanie.QueryRowContext(ctx, kodPliku).Scan(&plikID)
		if errors.Is(err, sql.ErrNoRows) {
			continue
		}
		if err != nil {
			return 0, fmt.Errorf("dane: nie można odnaleźć pliku %q: %w", kodPliku, err)
		}
		if _, err := wstawianie.ExecContext(ctx, kolekcja.ID, plikID); err != nil {
			return 0, fmt.Errorf("dane: nie można przypisać pliku %q regułą: %w", kodPliku, err)
		}
		przypisane++
	}
	return przypisane, nil
}

// OdepnijPrzypisaniaReguly zdejmuje z kolekcji wyłącznie zasoby wciągnięte
// regułą; przypisanie ręczne Operatora zostaje nietknięte.
func (r *repozytoriumBiblioteki) OdepnijPrzypisaniaReguly(ctx context.Context,
	kodKolekcji string) (int, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, usunPrzypisaniaRegulyBiblioteki)
	if err != nil {
		return 0, err
	}
	wynik, err := polecenie.ExecContext(ctx, kodKolekcji)
	if err != nil {
		return 0, fmt.Errorf("dane: nie można zdjąć przypisań reguły kolekcji %q: %w", kodKolekcji, err)
	}
	zdjete, err := wynik.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("dane: nie można policzyć zdjętych przypisań: %w", err)
	}
	return int(zdjete), nil
}

// KolekcjeWykaz zwraca kolekcje w postaci pełnej — z rodzicem, regułą
// i licznikiem zasobów — zawężone kolekcją nadrzędną i frazą nazwy.
func (r *repozytoriumBiblioteki) KolekcjeWykaz(ctx context.Context, rodzicKod, fraza *string,
	limit int) ([]KolekcjaBiblioteki, int, error) {

	warunki := []string{"1 = 1"}
	argumenty := []any{}
	if rodzicKod != nil && *rodzicKod != "" {
		warunki = append(warunki, `k.rodzic_id = (SELECT id FROM kolekcja_biblioteki
		                                           WHERE identyfikator_zewnetrzny = ?)`)
		argumenty = append(argumenty, *rodzicKod)
	}
	if fraza != nil && strings.TrimSpace(*fraza) != "" {
		warunki = append(warunki, "k.nazwa LIKE ?")
		argumenty = append(argumenty, "%"+strings.TrimSpace(*fraza)+"%")
	}
	warunek := strings.Join(warunki, " AND ")

	zapytanie := wykazKolekcjiBiblioteki + ` WHERE ` + warunek + ` ORDER BY k.nazwa, k.id LIMIT ?`
	wiersze, err := r.db.QueryContext(ctx, zapytanie,
		append(append([]any{}, argumenty...), granicaWykazu(limit))...)
	if err != nil {
		return nil, 0, fmt.Errorf("dane: nie można odczytać wykazu kolekcji: %w", err)
	}
	defer wiersze.Close()

	lista := []KolekcjaBiblioteki{}
	for wiersze.Next() {
		kolekcja, err := odczytajKolekcjePelnaBiblioteki(wiersze)
		if err != nil {
			return nil, 0, fmt.Errorf("dane: nieczytelny wiersz wykazu kolekcji: %w", err)
		}
		lista = append(lista, kolekcja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, 0, fmt.Errorf("dane: przerwany odczyt wykazu kolekcji: %w", err)
	}

	var lacznie int
	zapytanieLiczby := `SELECT COUNT(*) FROM kolekcja_biblioteki k WHERE ` + warunek
	if err := r.db.QueryRowContext(ctx, zapytanieLiczby, argumenty...).Scan(&lacznie); err != nil {
		return nil, 0, fmt.Errorf("dane: nie można policzyć kolekcji: %w", err)
	}
	return lista, lacznie, nil
}

// odczytajReguleBiblioteki składa regułę z jednego wiersza wyniku.
func odczytajReguleBiblioteki(wiersz skaner) (RegulaBiblioteki, error) {
	var regula RegulaBiblioteki
	var kolekcja, sciezka, przeliczenie sql.NullString
	var czynna int
	err := wiersz.Scan(&regula.ID, &regula.Kod, &regula.Rodzaj, &regula.Nazwa, &regula.Warunek,
		&kolekcja, &sciezka, &czynna, &przeliczenie, &regula.Utworzono)
	if err != nil {
		return RegulaBiblioteki{}, err
	}
	regula.KolekcjaDocelowaKod = tekstZKolumny(kolekcja)
	regula.SciezkaObserwowana = tekstZKolumny(sciezka)
	regula.OstatniePrzeliczenie = tekstZKolumny(przeliczenie)
	regula.Czynna = czynna == 1
	return regula, nil
}

// odczytajKolekcjePelnaBiblioteki składa kolekcję wraz z rodzicem, regułą i licznikiem.
func odczytajKolekcjePelnaBiblioteki(wiersz skaner) (KolekcjaBiblioteki, error) {
	var kolekcja KolekcjaBiblioteki
	var opis, rodzic, regula sql.NullString
	err := wiersz.Scan(&kolekcja.ID, &kolekcja.Kod, &kolekcja.Nazwa, &opis,
		&kolekcja.Utworzono, &kolekcja.Zaktualizowano, &rodzic, &regula, &kolekcja.LiczbaPlikow)
	if err != nil {
		return KolekcjaBiblioteki{}, err
	}
	kolekcja.Opis = tekstZKolumny(opis)
	kolekcja.RodzicKod = tekstZKolumny(rodzic)
	kolekcja.RegulaKod = tekstZKolumny(regula)
	return kolekcja, nil
}
