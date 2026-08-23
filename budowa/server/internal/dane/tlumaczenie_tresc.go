// Odpowiedzialność pliku: zmiany treści panelu tłumaczenia — tłumaczenie
// (korekta Operatora), tłumaczenie zwrotne, ton i stan. Cztery metody typu
// `*repozytoriumTlumaczen`, którego typ i konstruktor deklaruje
// `dane/tlumaczenie.go`; ten plik nie wnosi własnego typu ani interfejsu.
//
// Ton jest kolumną panelu, nie osobnym bytem: `panel.tone.set` pisze przez
// metodę tego pliku, bo `ton` mieszka w `panel_tlumaczenia`
// (`migracja_053_tlumaczenie.sql`).
//
// Każda z czterech metod przy nieznanym kodzie panelu wraca ErrBrakWiersza.
// Cicha zgoda na zmianę bytu, którego nie ma, byłaby potwierdzeniem czynności,
// która się nie odbyła, więc każda metoda sprawdza `RowsAffected` wzorem
// `UstawStanZlecenia` w `dane/asystent.go`.
//
// Czas jest liczbą milisekund epoki, wzorem `dane/asystent.go` — cały moduł
// Translate idzie tym samym sposobem.
package dane

import (
	"context"
	"fmt"
	"time"
)

const (
	ustawTrescPanelu = `UPDATE panel_tlumaczenia
	                    SET tresc = ?, tresc_odwolanie = ?, zaktualizowano = ?
	                    WHERE identyfikator_zewnetrzny = ?`

	ustawTrescZwrotnaPanelu = `UPDATE panel_tlumaczenia
	                           SET tresc_zwrotna = ?, zaktualizowano = ?
	                           WHERE identyfikator_zewnetrzny = ?`

	ustawTonPanelu = `UPDATE panel_tlumaczenia
	                  SET ton = ?, zaktualizowano = ?
	                  WHERE identyfikator_zewnetrzny = ?`

	ustawStanPanelu = `UPDATE panel_tlumaczenia
	                   SET stan = ?, zaktualizowano = ?
	                   WHERE identyfikator_zewnetrzny = ?`
)

// UstawTlumaczenie zapisuje korektę Operatora — treść krótką wprost w `tresc`
// i (dla treści obszernej) odwołanie do pliku w `tresc_odwolanie`.
// Oba pola przychodzą jako wskaźniki: `nil` zostawia kolumnę pustą, bo
// `translation.set` może nieść samą treść krótką albo samo odwołanie, zależnie
// od rozmiaru tekstu — rdzeń nie zgaduje, które pole miało zostać wyczyszczone.
// Kod panelu nieznany wraca jako ErrBrakWiersza.
func (r *repozytoriumTlumaczen) UstawTlumaczenie(ctx context.Context, kodPanelu string,
	tresc, odwolanie *string) (PanelTlumaczenia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, ustawTrescPanelu)
	if err != nil {
		return PanelTlumaczenia{}, err
	}
	wynik, err := polecenie.ExecContext(ctx, tekstDoKolumny(tresc), tekstDoKolumny(odwolanie),
		time.Now().UnixMilli(), kodPanelu)
	if err != nil {
		return PanelTlumaczenia{}, fmt.Errorf("dane: nie można ustawić tłumaczenia panelu %q: %w", kodPanelu, err)
	}
	dotknietych, err := wynik.RowsAffected()
	if err != nil {
		return PanelTlumaczenia{}, fmt.Errorf("dane: nieczytelny wynik zmiany tłumaczenia panelu %q: %w", kodPanelu, err)
	}
	if dotknietych == 0 {
		return PanelTlumaczenia{}, ErrBrakWiersza
	}
	return r.Panel(ctx, kodPanelu)
}

// UstawTlumaczenieZwrotne zapisuje wynik `backtranslation.run` w kolumnie
// `tresc_zwrotna`. Kontrakt oddaje jeden tekst na panel bez historii przebiegów
// (`migracja_053_tlumaczenie.sql`) — zapis nadpisuje, nie dokłada wiersza. Kod panelu nieznany
// wraca jako ErrBrakWiersza.
func (r *repozytoriumTlumaczen) UstawTlumaczenieZwrotne(ctx context.Context,
	kodPanelu, tresc string) (PanelTlumaczenia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, ustawTrescZwrotnaPanelu)
	if err != nil {
		return PanelTlumaczenia{}, err
	}
	wynik, err := polecenie.ExecContext(ctx, tresc, time.Now().UnixMilli(), kodPanelu)
	if err != nil {
		return PanelTlumaczenia{}, fmt.Errorf("dane: nie można ustawić tłumaczenia zwrotnego panelu %q: %w", kodPanelu, err)
	}
	dotknietych, err := wynik.RowsAffected()
	if err != nil {
		return PanelTlumaczenia{}, fmt.Errorf("dane: nieczytelny wynik zmiany tłumaczenia zwrotnego panelu %q: %w", kodPanelu, err)
	}
	if dotknietych == 0 {
		return PanelTlumaczenia{}, ErrBrakWiersza
	}
	return r.Panel(ctx, kodPanelu)
}

// UstawTon zapisuje ton panelu (kolumna `ton`) na żądanie `panel.tone.set`,
// które pisze przez tę metodę, bo `ton` jest polem panelu, a nie osobnym bytem.
// Kod panelu nieznany wraca jako ErrBrakWiersza.
func (r *repozytoriumTlumaczen) UstawTon(ctx context.Context, kodPanelu, ton string) (PanelTlumaczenia, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, ustawTonPanelu)
	if err != nil {
		return PanelTlumaczenia{}, err
	}
	wynik, err := polecenie.ExecContext(ctx, ton, time.Now().UnixMilli(), kodPanelu)
	if err != nil {
		return PanelTlumaczenia{}, fmt.Errorf("dane: nie można ustawić tonu panelu %q: %w", kodPanelu, err)
	}
	dotknietych, err := wynik.RowsAffected()
	if err != nil {
		return PanelTlumaczenia{}, fmt.Errorf("dane: nieczytelny wynik zmiany tonu panelu %q: %w", kodPanelu, err)
	}
	if dotknietych == 0 {
		return PanelTlumaczenia{}, ErrBrakWiersza
	}
	return r.Panel(ctx, kodPanelu)
}

// UstawStanPanelu zmienia stan panelu (pending/translating/ready/error,
// CHECK w `migracja_053_tlumaczenie.sql`, wartości kontraktu wprost). Kod panelu nieznany wraca
// jako ErrBrakWiersza — z tego samego powodu co pozostałe trzy metody pliku.
func (r *repozytoriumTlumaczen) UstawStanPanelu(ctx context.Context, kodPanelu, stan string) (PanelTlumaczenia, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, ustawStanPanelu)
	if err != nil {
		return PanelTlumaczenia{}, err
	}
	wynik, err := polecenie.ExecContext(ctx, stan, time.Now().UnixMilli(), kodPanelu)
	if err != nil {
		return PanelTlumaczenia{}, fmt.Errorf("dane: nie można ustawić stanu panelu %q: %w", kodPanelu, err)
	}
	dotknietych, err := wynik.RowsAffected()
	if err != nil {
		return PanelTlumaczenia{}, fmt.Errorf("dane: nieczytelny wynik zmiany stanu panelu %q: %w", kodPanelu, err)
	}
	if dotknietych == 0 {
		return PanelTlumaczenia{}, ErrBrakWiersza
	}
	return r.Panel(ctx, kodPanelu)
}
