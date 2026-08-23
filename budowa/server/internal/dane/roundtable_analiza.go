// Odpowiedzialność pliku: ustalenia analizy debaty i rejestr dowodów
// (`store/migracja_193_roundtable_analiza.sql`).
//
// Ustalenie zostaje w zapisie, bo analiza kosztuje wywołanie kanału modelu.
// Odczyt po fakcie — na przykład przy wydaniu transkryptu z ustaleniami — nie
// ma prawa wołać modelu po raz drugi po to samo.
package dane

import (
	"context"
	"database/sql"
	"fmt"
)

// UstalenieDebaty to jedno ustalenie analizy zapisu debaty.
type UstalenieDebaty struct {
	Kod       string
	Okno      string
	Rodzaj    string
	Tura      string
	Wypowiedz string
	Wezel     string
	Uczestnik string
	Tresc     string
	Pewnosc   float64
	Utworzono string
}

// DowodDebaty to pozycja rejestru dowodów: twierdzenie wraz z tym, czym je
// poparto — albo z odnotowanym brakiem poparcia.
type DowodDebaty struct {
	Kod         string
	Okno        string
	Tura        string
	Wypowiedz   string
	Twierdzenie string
	Zrodlo      *string
	Adres       *string
	Poparte     bool
	Utworzono   string
}

// RepozytoriumDebatyAnalizy jest częścią kontraktu obszaru odpowiadającą za
// wyniki analizy zapisu debaty.
type RepozytoriumDebatyAnalizy interface {
	ZastapUstaleniaDebaty(ctx context.Context, okno, rodzaj, tura string,
		ustalenia []UstalenieDebaty) error
	UstaleniaDebaty(ctx context.Context, okno, rodzaj, tura string) ([]UstalenieDebaty, error)

	ZastapDowodyDebaty(ctx context.Context, okno, tura string, dowody []DowodDebaty) error
	DowodyDebaty(ctx context.Context, okno, tura string, tylkoNiepoparte bool) ([]DowodDebaty, error)
}

const (
	kolumnyUstaleniaDebaty = `identyfikator_zewnetrzny, okno, rodzaj, tura, wypowiedz, wezel,
	                          uczestnik, tresc, pewnosc, utworzono`

	zapiszUstalenieDebaty = `INSERT INTO debata_ustalenie
	                         (identyfikator_zewnetrzny, okno, rodzaj, tura, wypowiedz, wezel,
	                          uczestnik, tresc, pewnosc)
	                         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	usunUstaleniaDebaty = `DELETE FROM debata_ustalenie
	                       WHERE okno = ? AND rodzaj = ? AND (? = '' OR tura = ?)`

	pobierzUstaleniaDebaty = `SELECT ` + kolumnyUstaleniaDebaty + `
	                          FROM debata_ustalenie
	                          WHERE okno = ? AND (? = '' OR rodzaj = ?) AND (? = '' OR tura = ?)
	                          ORDER BY id ASC`

	kolumnyDowoduDebaty = `identyfikator_zewnetrzny, okno, tura, wypowiedz, twierdzenie,
	                       zrodlo, adres, poparte, utworzono`

	zapiszDowodDebaty = `INSERT INTO debata_dowod
	                     (identyfikator_zewnetrzny, okno, tura, wypowiedz, twierdzenie,
	                      zrodlo, adres, poparte)
	                     VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	usunDowodyDebaty = `DELETE FROM debata_dowod WHERE okno = ? AND (? = '' OR tura = ?)`

	pobierzDowodyDebaty = `SELECT ` + kolumnyDowoduDebaty + `
	                       FROM debata_dowod
	                       WHERE okno = ? AND (? = '' OR tura = ?) AND (? = 0 OR poparte = 0)
	                       ORDER BY id ASC`
)

// ZastapUstaleniaDebaty wymienia ustalenia jednego rodzaju analizy. Ponowny
// przebieg zastępuje poprzedni wynik, bo dwa przebiegi tej samej analizy na tym
// samym zapisie nie są dwoma ustaleniami, tylko jednym powtórzonym.
func (r *repozytoriumRoundtable) ZastapUstaleniaDebaty(ctx context.Context,
	okno, rodzaj, tura string, ustalenia []UstalenieDebaty) error {

	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		wyczysc, err := r.zapytania.wTransakcji(ctx, transakcja, usunUstaleniaDebaty)
		if err != nil {
			return err
		}
		if _, err := wyczysc.ExecContext(ctx, okno, rodzaj, tura, tura); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić ustaleń analizy okna %q: %w", okno, err)
		}
		wstaw, err := r.zapytania.wTransakcji(ctx, transakcja, zapiszUstalenieDebaty)
		if err != nil {
			return err
		}
		for _, ustalenie := range ustalenia {
			if _, err := wstaw.ExecContext(ctx, ustalenie.Kod, okno, rodzaj, ustalenie.Tura,
				ustalenie.Wypowiedz, ustalenie.Wezel, ustalenie.Uczestnik, ustalenie.Tresc,
				ustalenie.Pewnosc); err != nil {
				return fmt.Errorf("dane: nie można zapisać ustalenia analizy %q: %w",
					ustalenie.Kod, err)
			}
		}
		return nil
	})
}

// UstaleniaDebaty zwraca ustalenia okna; rodzaj i tura puste nie zawężają.
func (r *repozytoriumRoundtable) UstaleniaDebaty(ctx context.Context,
	okno, rodzaj, tura string) ([]UstalenieDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzUstaleniaDebaty)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, rodzaj, rodzaj, tura, tura)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać ustaleń analizy okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	ustalenia := make([]UstalenieDebaty, 0, 8)
	for wiersze.Next() {
		var ustalenie UstalenieDebaty
		if err := wiersze.Scan(&ustalenie.Kod, &ustalenie.Okno, &ustalenie.Rodzaj, &ustalenie.Tura,
			&ustalenie.Wypowiedz, &ustalenie.Wezel, &ustalenie.Uczestnik, &ustalenie.Tresc,
			&ustalenie.Pewnosc, &ustalenie.Utworzono); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz ustalenia analizy: %w", err)
		}
		ustalenia = append(ustalenia, ustalenie)
	}
	return ustalenia, wiersze.Err()
}

// ZastapDowodyDebaty wymienia rejestr dowodów okna albo jednej tury.
func (r *repozytoriumRoundtable) ZastapDowodyDebaty(ctx context.Context,
	okno, tura string, dowody []DowodDebaty) error {

	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		wyczysc, err := r.zapytania.wTransakcji(ctx, transakcja, usunDowodyDebaty)
		if err != nil {
			return err
		}
		if _, err := wyczysc.ExecContext(ctx, okno, tura, tura); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić rejestru dowodów okna %q: %w", okno, err)
		}
		wstaw, err := r.zapytania.wTransakcji(ctx, transakcja, zapiszDowodDebaty)
		if err != nil {
			return err
		}
		for _, dowod := range dowody {
			if _, err := wstaw.ExecContext(ctx, dowod.Kod, okno, dowod.Tura, dowod.Wypowiedz,
				dowod.Twierdzenie, dowod.Zrodlo, dowod.Adres, dowod.Poparte); err != nil {
				return fmt.Errorf("dane: nie można zapisać pozycji rejestru dowodów %q: %w",
					dowod.Kod, err)
			}
		}
		return nil
	})
}

// DowodyDebaty zwraca rejestr dowodów; zawężenie do niepopartych wyciąga na
// wierzch twierdzenia bez źródła.
func (r *repozytoriumRoundtable) DowodyDebaty(ctx context.Context, okno, tura string,
	tylkoNiepoparte bool) ([]DowodDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzDowodyDebaty)
	if err != nil {
		return nil, err
	}
	zawezenie := 0
	if tylkoNiepoparte {
		zawezenie = 1
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, tura, tura, zawezenie)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać rejestru dowodów okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	dowody := make([]DowodDebaty, 0, 8)
	for wiersze.Next() {
		var dowod DowodDebaty
		if err := wiersze.Scan(&dowod.Kod, &dowod.Okno, &dowod.Tura, &dowod.Wypowiedz,
			&dowod.Twierdzenie, &dowod.Zrodlo, &dowod.Adres, &dowod.Poparte,
			&dowod.Utworzono); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz rejestru dowodów: %w", err)
		}
		dowody = append(dowody, dowod)
	}
	return dowody, wiersze.Err()
}
