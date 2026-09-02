// Plik obsługuje czynności channel.check i channel.credential.status: sprawdzenie kanału modelu prawdziwym wywołaniem oraz odczyt stanu poświadczenia bez ujawniania jego treści.
package core

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

const (
	// Dłuższe czekanie niż piętnaście sekund nie jest odpowiedzią przy przycisku.
	granicaSprawdzeniaKanalu = 15 * time.Second

	trescSprawdzeniaKanalu = "ping"
)

type sejfPoswiadczen interface {
	Odczytaj(ctx context.Context, byt string) (string, bool)
}

func (a *adapterKanalow) ZSejfem(sejf sejfPoswiadczen) *adapterKanalow {
	a.sejf = sejf
	return a
}

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
	if _, jest := kanalKonta(ctx, a.repozytorium, a.rejestr, kod); !jest {
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
		// Kanał bez odwołania nie wymaga uwierzytelnienia (model lokalny); to odpowiedź, nie brak.
		return shared.ChannelCredentialStatusResponse{Status: stan}, nil
	}

	rodzaj := rodzajPoswiadczeniaKanalu(kanal)
	stan.Kind = &rodzaj
	zarzadca := zarzadcaPoswiadczeniaKanalu(kanal)
	stan.ManagedBy = &zarzadca

	if a.sejf != nil {
		if wartosc, jest := a.sejf.Odczytaj(ctx, odwolanie); jest {
			// Wartość służy wyłącznie rozstrzygnięciu, czy jest niepusta, i nigdzie nie wychodzi.
			stan.Present = strings.TrimSpace(wartosc) != ""
		}
	}
	if !stan.Present {
		zarzadca = zarzadcaPoswiadczeniaKanalu(kanal) + " (odwołanie: " + odwolanie + ")"
		stan.ManagedBy = &zarzadca
	}
	return shared.ChannelCredentialStatusResponse{Status: stan}, nil
}

const przedrostekSprawdzeniaKanalu = "sprawdzenie-kanalu-"

func rodzajPoswiadczeniaKanalu(kanal dane.Kanal) string {
	var parametry map[string]any
	if len(kanal.ParametryJSON) > 0 {
		_ = json.Unmarshal([]byte(kanal.ParametryJSON), &parametry)
	}
	if rodzaj, jest := parametry["credentialKind"].(string); jest && strings.TrimSpace(rodzaj) != "" {
		return rodzaj
	}
	if kanal.KontoDostawcyID != nil {
		return "konto rejestru kont"
	}
	return "klucz kanału API"
}

func zarzadcaPoswiadczeniaKanalu(kanal dane.Kanal) string {
	if kanal.KontoDostawcyID != nil {
		return "rejestr kont platformy"
	}
	return "sejf poświadczeń serwera"
}

func bladWskazaniaKanalu(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"kanały: "+powod))
}

// Więz klucza obcego znaczy okna albo zlecenia na kanale — stan do naprawy przez
// Operatora, nie usterka rdzenia; odmowa nazywa kanał kodem wołającego.
func bladUsunieciaKanalu(err error, kod string) error {
	if err != nil && strings.Contains(err.Error(), "FOREIGN KEY") {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeConflict,
			"kanały: kanał "+strconv.Quote(kod)+
				" jest w użyciu — wiąże go okno albo zlecenie; zdejmij wiązania przed usunięciem"))
	}
	return err
}

// Wykaz rodzajów kanału stoi w więzie kolumny `kanal_modelu.rodzaj_kanalu`;
// rodzaj spoza wykazu wraca surową treścią więzu i wymaga przełożenia na odmowę.
func bladZapisuKanalu(err error, rodzaj string) error {
	if err != nil && strings.Contains(err.Error(), "rodzaj_kanalu") {
		return bladWskazaniaKanalu("rodzaj kanału spoza wykazu rejestru: " +
			strconv.Quote(rodzaj) + "; znane rodzaje: " +
			strings.Join(shared.KnownChannelKinds(), ", "))
	}
	return err
}
