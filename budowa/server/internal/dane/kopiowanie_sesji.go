package dane

import (
	"context"
	"fmt"
)

// KopiujZapisSesji powiela okna i wiadomości sesji źródłowej do sesji docelowej.
//
// Kopiowanie idzie po wierszach bazy, nie po rejestrze żywym: kopia ma być
// wiernym odbiciem zapisu, a rejestr żywy niesie wyłącznie sesje otwarte.
// Sesja sprzed restartu ma dać się skopiować tak samo jak ta z bieżącej pracy.
//
// Odwzorowanie okien źródłowych na docelowe podaje wywołujący — to on nadaje
// nowym oknom identyfikatory rdzenia i zna ich wiersze. Warstwa danych nie
// zakłada okien sama, bo cykl życia okna należy do pakietu sesji.
//
// Wiadomość kopiowana traci powiązanie z oknem źródłowym (`okno_zrodlowe_id`):
// w kopii nie ma bytu, na który mogłoby wskazywać, a wskazanie na okno oryginału
// byłoby więzią między dwiema niezależnymi sesjami.
// Zwraca liczbę skopiowanych wiadomości — jedyną miarę, po którą sięga
// wywołujący. Struktura wyniku byłaby typem bez odbiorcy.
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
