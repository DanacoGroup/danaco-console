// Plik obsługuje czynności channel.check i channel.credential.status: sprawdzenie kanału modelu prawdziwym wywołaniem oraz odczyt stanu poświadczenia bez ujawniania jego treści.
package core

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

const (
	// granicaSprawdzeniaKanalu domyka jedno sprawdzenie po piętnastu sekundach, bo dłuższe czekanie nie jest odpowiedzią przy przycisku.
	granicaSprawdzeniaKanalu = 15 * time.Second

	// trescSprawdzeniaKanalu jest krótkim zapytaniem sprawdzającym: pyta, czy kanał odpowiada, zamiast prosić o treść.
	trescSprawdzeniaKanalu = "ping"
)

// sejfPoswiadczen jest tą częścią sejfu, której potrzebuje stan poświadczenia: samym pytaniem, czy pod odwołaniem coś leży.
type sejfPoswiadczen interface {
	Odczytaj(ctx context.Context, byt string) (string, bool)
}

// ZSejfem wpina sejf poświadczeń jako źródło stanu poświadczenia, nie jego treści, do adaptera kanałów rejestru modeli.
func (a *adapterKanalow) ZSejfem(sejf sejfPoswiadczen) *adapterKanalow {
	a.sejf = sejf
	return a
}

// Sprawdz obsługuje channel.check: wysyła krótkie zapytanie sprawdzające tym samym rejestrem kanałów, którym jedzie okno rozmowy.
func (a *adapterKanalow) Sprawdz(ctx context.Context,
	z shared.ChannelCheckRequest) (shared.ChannelCheckResponse, error) {

	kod := strings.TrimSpace(z.ChannelId)
	if kod == "" {
		return shared.ChannelCheckResponse{}, bladWskazaniaKanalu(
			"sprawdzenie bez wskazania kanału")
	}
	chwila := time.Now()
	if a.rejestr == nil {
		szczegol := "serwer nie ma wpiętego rejestru kanałów — nie ma czym wykonać sprawdzenia"
		return shared.ChannelCheckResponse{
			Reachable: false, CheckedAt: chwila.UnixMilli(), Detail: &szczegol,
		}, nil
	}
	if _, jest := a.rejestr.Kanal(kod); !jest {
		// Kanał, którego nie ma w rejestrze, oddaje odpowiedź nie odpowiada, a nie odmowę sprawdzenia.
		szczegol := "kanału nie ma w rejestrze kanałów serwera albo jest wyłączony"
		return shared.ChannelCheckResponse{
			Reachable: false, CheckedAt: chwila.UnixMilli(), Detail: &szczegol,
		}, nil
	}

	sprawdzenie, odwolaj := context.WithTimeout(ctx, granicaSprawdzeniaKanalu)
	defer odwolaj()

	odpowiedzial := false
	ujscie := models.UjscieFunkcji(func(_ context.Context, f models.Fragment) error {
		if f.Kind == shared.ChunkKindText {
			odpowiedzial = true
		}
		return nil
	})
	zapytanie := models.Zapytanie{
		Wiadomosc: nowyIdentyfikator(przedrostekSprawdzeniaKanalu),
		Tresc:     trescSprawdzeniaKanalu,
		Kanal:     kod,
	}

	poczatek := time.Now()
	err := a.rejestr.Wyslij(sprawdzenie, zapytanie, ujscie)
	czas := int(time.Since(poczatek).Milliseconds())
	sprawdzono := time.Now().UnixMilli()

	if err != nil {
		szczegol := err.Error()
		return shared.ChannelCheckResponse{
			Reachable: false, CheckedAt: sprawdzono, Detail: &szczegol,
		}, nil
	}
	if !odpowiedzial {
		szczegol := "kanał zakończył wywołanie bez błędu, ale nie oddał ani jednego fragmentu treści"
		return shared.ChannelCheckResponse{
			Reachable: false, CheckedAt: sprawdzono, LatencyMs: &czas, Detail: &szczegol,
		}, nil
	}
	return shared.ChannelCheckResponse{
		Reachable: true, CheckedAt: sprawdzono, LatencyMs: &czas,
	}, nil
}

// StanPoswiadczenia obsługuje channel.credential.status i oddaje wyłącznie stan poświadczenia kanału, nigdy jego treść.
func (a *adapterKanalow) StanPoswiadczenia(ctx context.Context,
	z shared.ChannelCredentialStatusRequest) (shared.ChannelCredentialStatusResponse, error) {

	kod := strings.TrimSpace(z.ChannelId)
	if kod == "" {
		return shared.ChannelCredentialStatusResponse{}, bladWskazaniaKanalu(
			"pytanie o poświadczenie bez wskazania kanału")
	}
	if a.repozytorium == nil {
		return shared.ChannelCredentialStatusResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeInternalError,
			"kanały: serwer nie ma wpiętego rejestru kanałów"))
	}
	kanal, err := a.repozytorium.PobierzPoKodzie(ctx, kod)
	if err != nil {
		return shared.ChannelCredentialStatusResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeNotFound, "kanały: nie ma kanału o identyfikatorze "+kod))
	}

	stan := shared.ChannelCredentialStatus{ChannelId: kod}
	odwolanie := strings.TrimSpace(wartoscTekstu(kanal.PoswiadczenieOdwolanie))
	if odwolanie == "" {
		// Kanał bez odwołania nie wymaga uwierzytelnienia, na przykład model lokalny; to odpowiedź, nie brak.
		return shared.ChannelCredentialStatusResponse{Status: stan}, nil
	}

	rodzaj := rodzajPoswiadczeniaKanalu(kanal)
	stan.Kind = &rodzaj
	zarzadca := zarzadcaPoswiadczeniaKanalu(kanal)
	stan.ManagedBy = &zarzadca

	if a.sejf != nil {
		if wartosc, jest := a.sejf.Odczytaj(ctx, odwolanie); jest {
			// Wartość służy wyłącznie do rozstrzygnięcia, czy jest niepusta, i nigdzie indziej nie wychodzi.
			stan.Present = strings.TrimSpace(wartosc) != ""
		}
	}
	if !stan.Present {
		// Poświadczenie bywa też zmienną środowiskową maszyny rdzenia; odpowiedź podaje, gdzie rdzeń szukał.
		zarzadca = zarzadcaPoswiadczeniaKanalu(kanal) + " (odwołanie: " + odwolanie + ")"
		stan.ManagedBy = &zarzadca
	}
	return shared.ChannelCredentialStatusResponse{Status: stan}, nil
}

// przedrostekSprawdzeniaKanalu znakuje identyfikator wywołania sprawdzającego w rejestrze kanałów modeli.
const przedrostekSprawdzeniaKanalu = "sprawdzenie-kanalu-"

// rodzajPoswiadczeniaKanalu nazywa rodzaj poświadczenia na podstawie parametrów
// wiersza. Kontrakt nie ma osobnego pola rodzaju przy zapisie kanału, więc
// rodzaj bierze się z tego, co wiersz naprawdę niesie.
func rodzajPoswiadczeniaKanalu(kanal dane.Kanal) string {
	var parametry map[string]any
	if len(kanal.ParametryJSON) > 0 {
		_ = json.Unmarshal([]byte(kanal.ParametryJSON), &parametry)
	}
	if rodzaj, jest := parametry["credentialKind"].(string); jest && strings.TrimSpace(rodzaj) != "" {
		return rodzaj
	}
	if kanal.KontoID != nil {
		return "konto rejestru kont"
	}
	return "klucz kanału API"
}

// zarzadcaPoswiadczeniaKanalu nazywa miejsce, w którym poświadczenie kanału mieszka: rejestr kont albo sejf rdzenia.
func zarzadcaPoswiadczeniaKanalu(kanal dane.Kanal) string {
	if kanal.KontoID != nil {
		return "rejestr kont platformy"
	}
	return "sejf poświadczeń serwera"
}

// bladWskazaniaKanalu nazywa niepoprawne żądanie czynności rejestru kanałów —
// błąd Operatora, nie rdzenia.
func bladWskazaniaKanalu(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"kanały: "+powod))
}
