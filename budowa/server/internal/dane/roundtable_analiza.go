// Odpowiedzialność pliku: ustalenia analizy debaty i rejestr dowodów, zapisywane trwale, bo analiza kosztuje wywołanie kanału modelu.
package dane

import (
	"context"
	"database/sql"
	"fmt"
)

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
	                          uczestnik, tresc, pewnosc, konto_id)
	                         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ` + WskazanieKonta + `)`

	usunUstaleniaDebaty = `DELETE FROM debata_ustalenie
	                       WHERE okno = ? AND rodzaj = ? AND (? = '' OR tura = ?)
	                         AND ` + WarunekKonta

	pobierzUstaleniaDebaty = `SELECT ` + kolumnyUstaleniaDebaty + `
	                          FROM debata_ustalenie
	                          WHERE okno = ? AND (? = '' OR rodzaj = ?) AND (? = '' OR tura = ?)
	                            AND ` + WarunekKonta + `
	                          ORDER BY id ASC`

	kolumnyDowoduDebaty = `identyfikator_zewnetrzny, okno, tura, wypowiedz, twierdzenie,
	                       zrodlo, adres, poparte, utworzono`

	// debata_dowod nie ma konto_id (migracja 193, krok 484 ją pominął): granica idzie przez
	// wypowiedź do tury, której konto_id wiąże niekwalifikowaną nazwę warunku.
	granicaWypowiedziDowodu = `EXISTS (SELECT 1 FROM debata_wypowiedz w
	                                    JOIN debata_tura t ON t.id = w.tura_id
	                                    WHERE w.identyfikator_zewnetrzny = `

	zapiszDowodDebaty = `INSERT INTO debata_dowod
	                     (identyfikator_zewnetrzny, okno, tura, wypowiedz, twierdzenie,
	                      zrodlo, adres, poparte)
	                     SELECT ?, ?, ?, ?, ?, ?, ?, ?
	                     WHERE ` + granicaWypowiedziDowodu + `? AND ` + WarunekKonta + `)`

	usunDowodyDebaty = `DELETE FROM debata_dowod
	                    WHERE okno = ? AND (? = '' OR tura = ?)
	                      AND ` + granicaWypowiedziDowodu + `debata_dowod.wypowiedz
	                                                          AND ` + WarunekKonta + `)`

	pobierzDowodyDebaty = `SELECT ` + kolumnyDowoduDebaty + `
	                       FROM debata_dowod
	                       WHERE okno = ? AND (? = '' OR tura = ?) AND (? = 0 OR poparte = 0)
	                         AND ` + granicaWypowiedziDowodu + `debata_dowod.wypowiedz
	                                                             AND ` + WarunekKonta + `)
	                       ORDER BY id ASC`
)

// Ponowny przebieg tej samej analizy na tym samym zapisie jest jednym ustaleniem powtórzonym, nie dwoma.
func (r *repozytoriumRoundtable) ZastapUstaleniaDebaty(ctx context.Context,
	okno, rodzaj, tura string, ustalenia []UstalenieDebaty) error {

	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		wyczysc, err := r.zapytania.wTransakcji(ctx, transakcja, usunUstaleniaDebaty)
		if err != nil {
			return err
		}
		if _, err := wyczysc.ExecContext(ctx, okno, rodzaj, tura, tura,
			KontoOperatora(ctx)); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić ustaleń analizy okna %q: %w", okno, err)
		}
		wstaw, err := r.zapytania.wTransakcji(ctx, transakcja, zapiszUstalenieDebaty)
		if err != nil {
			return err
		}
		for _, ustalenie := range ustalenia {
			if _, err := wstaw.ExecContext(ctx, ustalenie.Kod, okno, rodzaj, ustalenie.Tura,
				ustalenie.Wypowiedz, ustalenie.Wezel, ustalenie.Uczestnik, ustalenie.Tresc,
				ustalenie.Pewnosc, KontoOperatora(ctx)); err != nil {
				return fmt.Errorf("dane: nie można zapisać ustalenia analizy %q: %w",
					ustalenie.Kod, err)
			}
		}
		return nil
	})
}

func (r *repozytoriumRoundtable) UstaleniaDebaty(ctx context.Context,
	okno, rodzaj, tura string) ([]UstalenieDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzUstaleniaDebaty)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, rodzaj, rodzaj, tura, tura,
		KontoOperatora(ctx))
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

// ZastapDowodyDebaty wymienia rejestr dowodów okna albo jednej tury; pozycja bez wypowiedzi
// konta zamawiającego kończy się ErrKolizjaWiersza.
func (r *repozytoriumRoundtable) ZastapDowodyDebaty(ctx context.Context,
	okno, tura string, dowody []DowodDebaty) error {

	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		wyczysc, err := r.zapytania.wTransakcji(ctx, transakcja, usunDowodyDebaty)
		if err != nil {
			return err
		}
		if _, err := wyczysc.ExecContext(ctx, okno, tura, tura, KontoOperatora(ctx)); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić rejestru dowodów okna %q: %w", okno, err)
		}
		wstaw, err := r.zapytania.wTransakcji(ctx, transakcja, zapiszDowodDebaty)
		if err != nil {
			return err
		}
		for _, dowod := range dowody {
			wynik, err := wstaw.ExecContext(ctx, dowod.Kod, okno, dowod.Tura, dowod.Wypowiedz,
				dowod.Twierdzenie, dowod.Zrodlo, dowod.Adres, dowod.Poparte,
				dowod.Wypowiedz, KontoOperatora(ctx))
			if err != nil {
				return fmt.Errorf("dane: nie można zapisać pozycji rejestru dowodów %q: %w",
					dowod.Kod, err)
			}
			if err := sprawdzTrafienieZapisu(wynik, "wypowiedź dowodu", dowod.Wypowiedz); err != nil {
				return err
			}
		}
		return nil
	})
}

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
	wiersze, err := polecenie.QueryContext(ctx, okno, tura, tura, zawezenie, KontoOperatora(ctx))
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
