package dane

import (
	"context"
	"fmt"
)

// KopiujZapisSesji powiela okna i wiadomości sesji źródłowej do sesji docelowej, kopiując wiersze bazy danych.
func (z *Zestaw) KopiujZapisSesji(ctx context.Context,
	oknaZrodlowe []Okno, wgOkna map[int64]int64) (int, error) {

	skopiowanych := 0
	for _, zrodlowe := range oknaZrodlowe {
		docelowe, jest := wgOkna[zrodlowe.ID]
		if !jest {
			continue
		}
		wiadomosci, err := z.Wiadomosci.ListaOkna(ctx, zrodlowe.ID, 0)
		if err != nil {
			return skopiowanych, fmt.Errorf("dane: odczyt wiadomości okna %d: %w", zrodlowe.ID, err)
		}
		for _, wiadomosc := range wiadomosci {
			kopia := wiadomosc
			kopia.ID = 0
			kopia.OknoID = docelowe
			kopia.IdentyfikatorZewnetrzny = nil
			kopia.OknoZrodloweID = nil
			if _, err := z.Wiadomosci.Dopisz(ctx, kopia); err != nil {
				return skopiowanych, fmt.Errorf("dane: zapis kopii wiadomości: %w", err)
			}
			skopiowanych++
		}
	}
	return skopiowanych, nil
}
