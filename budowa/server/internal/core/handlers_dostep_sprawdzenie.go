// Plik usuwa punkt dostępu i sprawdza, czy punkt odpowiada oraz jakie
// korzenie potwierdza, korzystając z próby właściwej rodzajowi punktu.
package core

import (
	"context"
	"os"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// ProbaPunktu bada osiągalność punktu dostępu. Zwraca stan, korzenie
// potwierdzone oraz szczegół niepowodzenia.
type ProbaPunktu func(ctx context.Context, punkt dane.PunktDostepu) (shared.AccessPointStatus, []string, string)

// Usun kasuje punkt wraz z nadaniami, które się na niego powoływały — kasuje je
// więz klucza obcego, nie kod rdzenia. Punkt nieznany nie jest błędem: wynik
// mówi wtedy, że nic nie usunięto. Brak wpiętego katalogu jest odmową.
func (a *adapterPunktowDostepu) Usun(ctx context.Context,
	z shared.AccessPointRemoveRequest) (shared.AccessPointRemoveResponse, error) {

	if a == nil || a.repozytorium == nil {
		return shared.AccessPointRemoveResponse{}, bladBrakuKatalogu("punktów dostępu")
	}
	punkt, err := a.repozytorium.PoKodzie(ctx, z.AccessPointId)
	if brakWiersza(err) {
		return shared.AccessPointRemoveResponse{Removed: false}, nil
	}
	if err != nil {
		return shared.AccessPointRemoveResponse{}, err
	}
	if err := a.repozytorium.Usun(ctx, punkt.ID); err != nil {
		if brakWiersza(err) {
			return shared.AccessPointRemoveResponse{Removed: false}, nil
		}
		return shared.AccessPointRemoveResponse{}, err
	}
	usunPoswiadczenie(ctx, a.sejf, punkt.Kod)
	return shared.AccessPointRemoveResponse{Removed: true}, nil
}

// Sprawdz bada punkt i utrwala wynik badania w katalogu, żeby okno konfiguracji
// pokazywało stan także po ponownym otwarciu.
func (a *adapterPunktowDostepu) Sprawdz(ctx context.Context,
	z shared.AccessPointCheckRequest) (shared.AccessPointCheckResponse, error) {

	chwila := time.Now().UTC()
	if a == nil || a.repozytorium == nil {
		return wynikSprawdzenia(shared.AccessPointStatusUnknown, chwila, nil,
			"katalog punktów dostępu nie jest wpięty"), nil
	}
	punkt, err := a.repozytorium.PoKodzie(ctx, z.AccessPointId)
	if brakWiersza(err) {
		return wynikSprawdzenia(shared.AccessPointStatusUnknown, chwila, nil,
			"punkt dostępu "+z.AccessPointId+" nie istnieje"), nil
	}
	if err != nil {
		return shared.AccessPointCheckResponse{}, err
	}
	stan, korzenie, szczegol := a.proba()(ctx, punkt)
	if err := a.repozytorium.ZapiszWynikSprawdzenia(ctx, punkt.ID, stan,
		chwila.Format(formatZnacznikaBazy)); err != nil {
		return shared.AccessPointCheckResponse{}, err
	}
	return wynikSprawdzenia(stan, chwila, korzenie, szczegol), nil
}

// proba zwraca badanie wpięte w adapter albo badanie wbudowane, bo brak
// wpięcia nie wyłącza komendy sprawdzenia.
func (a *adapterPunktowDostepu) proba() ProbaPunktu {
	if a.sprawdzenie != nil {
		return a.sprawdzenie
	}
	return probaKorzeniLokalnych
}

// probaKorzeniLokalnych sprawdza korzenie katalogu na systemie plików
// widzianym przez rdzeń, zwracając stan nierozpoznany zamiast fałszywej
// niedostępności.
func probaKorzeniLokalnych(_ context.Context, punkt dane.PunktDostepu) (shared.AccessPointStatus, []string, string) {
	if punkt.Rodzaj != shared.AccessPointKindLocalDirectory {
		return shared.AccessPointStatusUnknown, nil,
			"punkt mostowy sprawdza klient modelu; serwer nie odpytuje mostu MCP"
	}
	potwierdzone := make([]string, 0, len(punkt.Korzenie))
	brakujace := ""
	for _, korzen := range punkt.Korzenie {
		stan, err := os.Stat(korzen)
		if err != nil || !stan.IsDir() {
			if brakujace == "" {
				brakujace = korzen
			}
			continue
		}
		potwierdzone = append(potwierdzone, korzen)
	}
	if brakujace != "" {
		return shared.AccessPointStatusUnknown, potwierdzone,
			"korzenia " + brakujace + " nie widać z maszyny serwera; katalog jest zapisany dla urządzenia"
	}
	if len(potwierdzone) == 0 {
		return shared.AccessPointStatusUnknown, potwierdzone,
			"punkt nie ma korzeni, więc nie ma czego potwierdzić"
	}
	return shared.AccessPointStatusReachable, potwierdzone, ""
}

// wynikSprawdzenia składa odpowiedź kontraktu na sprawdzenie punktu, niosącą
// stan, korzenie potwierdzone i szczegół.
func wynikSprawdzenia(stan shared.AccessPointStatus, chwila time.Time,
	korzenie []string, szczegol string) shared.AccessPointCheckResponse {

	return shared.AccessPointCheckResponse{
		Status:    stan,
		CheckedAt: chwila.UnixMilli(),
		Roots:     korzenie,
		Detail:    tekstOpcjonalny(szczegol),
	}
}
