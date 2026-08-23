// Odpowiedzialność pliku: przyjęcie polecenia asystenta (głosowego albo
// tekstowego) z okna Actions Monitor — zasila `assistant.voice.command`.
// Plik stoi osobno od `asystent.go` i `asystent_dziennik.go`.
//
// Zlecenie i pierwszy wpis dziennika powstają w jednej transakcji. Polecenie
// przyjęte bez śladu w dzienniku (zlecenie jest, rozmowa nie ma pierwszej
// linii) albo ślad bez zlecenia (wpis wisi na kodzie, którego zlecenie nigdy
// nie powstało) to stan połowiczny — stąd `wTransakcji`, wzorem
// `ZapiszKompozycje` w `design_kompozycje.go` i `ZapiszKroki`
// w `automations_kroki.go`.
//
// Ten plik nie rozpoznaje mowy. Kontrakt `assistant.voice.command` daje
// `AudioRef` (odnośnik do nagranego już pliku) albo `Transcript` (tekst
// poprawiony przez Operatora), a zapis idzie dosłownie: odnośnik do
// `nagranie_odnosnik`, tekst do `tresc`. Brak obu jest błędem żądania, nie
// pustym zapisem — kolumna `tresc` jest `NOT NULL`
// (`migracja_050_asystent.sql`), bo wpis dziennika bez treści nie opisuje
// niczego, co się wydarzyło.
package dane

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// PrzyjmijPolecenie zakłada zlecenie asystenta i zapisuje pierwszy wpis jego
// dziennika (rozmowy) w jednej transakcji. Zwraca oba wiersze po zapisie, ze
// wszystkimi polami uzupełnionymi przez bazę (identyfikator wewnętrzny,
// chwile czasu).
func (r *repozytoriumAsystenta) PrzyjmijPolecenie(ctx context.Context, zlecenie ZlecenieAsystenta,
	wpis WpisDziennikaAsystenta) (ZlecenieAsystenta, WpisDziennikaAsystenta, error) {

	if zlecenie.Kod == "" {
		return ZlecenieAsystenta{}, WpisDziennikaAsystenta{}, fmt.Errorf("dane: polecenie asystenta bez kodu zlecenia")
	}
	if zlecenie.OknoKod == "" {
		return ZlecenieAsystenta{}, WpisDziennikaAsystenta{}, fmt.Errorf("dane: polecenie asystenta %q bez okna", zlecenie.Kod)
	}
	if wpis.Kod == "" {
		return ZlecenieAsystenta{}, WpisDziennikaAsystenta{}, fmt.Errorf("dane: polecenie asystenta %q bez kodu wpisu dziennika", zlecenie.Kod)
	}
	// Ani nagranie, ani tekst — rdzeń nie ma z czego złożyć wpisu, a bez
	// transkrypcji nie ma jak jej dorobić w locie. To błąd żądania,
	// nie sytuacja, w której zapisujemy pusty wiersz.
	brakTresci := wpis.Tresc == ""
	brakNagrania := wpis.NagranieOdnosnik == nil || *wpis.NagranieOdnosnik == ""
	if brakTresci && brakNagrania {
		return ZlecenieAsystenta{}, WpisDziennikaAsystenta{}, fmt.Errorf(
			"dane: polecenie asystenta %q bez treści i bez odnośnika do nagrania", zlecenie.Kod)
	}

	teraz := time.Now().UnixMilli()
	utworzonoZlecenia := zlecenie.Utworzono
	if utworzonoZlecenia == 0 {
		utworzonoZlecenia = teraz
	}
	utworzonoWpisu := wpis.Utworzono
	if utworzonoWpisu == 0 {
		utworzonoWpisu = teraz
	}
	// Wpis dziennika wisi na kodzie zlecenia, które dopiero powstaje w tej
	// samej transakcji — powiązanie jest więc zawsze kodem przyjmowanego
	// zlecenia, niezależnie od tego, co ewentualnie niósł parametr `wpis`.
	kodZlecenia := zlecenie.Kod

	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		zapisZlecenia, err := r.zapytania.wTransakcji(ctx, transakcja, zapiszZlecenieAsystenta)
		if err != nil {
			return err
		}
		if _, err := zapisZlecenia.ExecContext(ctx, zlecenie.Kod, zlecenie.OknoKod,
			tekstDoKolumny(zlecenie.Tytul), zlecenie.Stan, zlecenie.Droga,
			liczbaDoKolumny(zlecenie.EtapBiezacy), liczbaDoKolumny(zlecenie.LiczbaEtapow),
			liczbaDoKolumny(zlecenie.Priorytet), tekstDoKolumny(zlecenie.Wynik),
			tekstDoKolumny(zlecenie.ProfilKod), utworzonoZlecenia, teraz); err != nil {
			return fmt.Errorf("dane: nie można założyć zlecenia asystenta %q: %w", zlecenie.Kod, err)
		}

		zapisWpisu, err := r.zapytania.wTransakcji(ctx, transakcja, wstawWpisDziennikaAsystenta)
		if err != nil {
			return err
		}
		if _, err := zapisWpisu.ExecContext(ctx, wpis.Kod, zlecenie.OknoKod, kodZlecenia,
			wpis.Rodzaj, wpis.Tresc, tekstDoKolumny(wpis.NagranieOdnosnik), utworzonoWpisu); err != nil {
			return fmt.Errorf("dane: nie można zapisać wpisu dziennika %q polecenia asystenta %q: %w",
				wpis.Kod, zlecenie.Kod, err)
		}
		return nil
	})
	if err != nil {
		return ZlecenieAsystenta{}, WpisDziennikaAsystenta{}, err
	}

	zapisaneZlecenie, err := r.Zlecenie(ctx, zlecenie.Kod)
	if err != nil {
		return ZlecenieAsystenta{}, WpisDziennikaAsystenta{}, fmt.Errorf(
			"dane: nie można odczytać zapisanego zlecenia asystenta %q: %w", zlecenie.Kod, err)
	}
	zapisanyWpis := wpis
	zapisanyWpis.OknoKod = zlecenie.OknoKod
	zapisanyWpis.ZlecenieKod = &kodZlecenia
	zapisanyWpis.Utworzono = utworzonoWpisu
	return zapisaneZlecenie, zapisanyWpis, nil
}
