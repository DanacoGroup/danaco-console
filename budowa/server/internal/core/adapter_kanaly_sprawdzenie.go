// Odpowiedzialność pliku: `channel.check` i `channel.credential.status` — dwie
// czynności pytające o kanał modelu, dołożone do rejestru z `adapter_kanaly.go`.
//
// ── Sprawdzenie jest narzędziem pomocniczym, nie bramką ─────────────────────
// Wynik `channel.check` NICZEGO nie warunkuje: nie wstrzymuje zapisu eksperta,
// nie blokuje wysłania tury, nie wyłącza kanału. Operator pyta „czy to
// odpowiada", dostaje odpowiedź i sam decyduje, co z nią zrobić. Kanał, który
// nie odpowiedział minutę temu, bywa sprawny teraz — i odwrotnie.
//
// Sprawdzenie jest PRAWDZIWYM wywołaniem: rdzeń wysyła krótkie zapytanie tym
// samym rejestrem kanałów, którym jedzie okno rozmowy, i mierzy czas
// odpowiedzi. Sprawdzenie, które ogląda wyłącznie wiersz rejestru, mówiłoby
// „skonfigurowany", a Operator czyta z niego „działa".
//
// ── Stan poświadczenia to STAN, nigdy treść ─────────────────────────────────
// `channel.credential.status` oddaje: czy poświadczenie jest ustawione, jakiego
// jest rodzaju i kiedy zostało zmienione. TREŚCI NIE ODDAJE I ODDAWAĆ NIE MOŻE.
// W tym pliku nie ma ani jednej ścieżki, którą sekret wychodzi do klienta —
// wartość z sejfu służy wyłącznie do rozstrzygnięcia, czy jest niepusta.
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
	// granicaSprawdzeniaKanalu domyka jedno sprawdzenie. Piętnaście sekund
	// starcza każdemu kanałowi na odpowiedź na jedno słowo, a dłuższe czekanie
	// i tak nie jest odpowiedzią, na którą ktokolwiek czeka przy przycisku.
	granicaSprawdzeniaKanalu = 15 * time.Second

	// trescSprawdzeniaKanalu jest zapytaniem sprawdzającym. Krótkie z zamysłu:
	// pytamy, czy kanał odpowiada, a nie prosimy o treść, za którą Operator
	// zapłaci przy każdym kliknięciu w „sprawdź".
	trescSprawdzeniaKanalu = "ping"
)

// sejfPoswiadczen jest tą częścią sejfu, której potrzebuje stan poświadczenia:
// samo pytanie, czy pod odwołaniem coś leży.
//
// Interfejs zamiast typu, bo adapter ma widzieć wyłącznie odczyt — węższy
// widok jest tu granicą, a nie ozdobą: adapter, który nie umie zapisać, nie
// nadpisze cudzego klucza przy pomyłce.
type sejfPoswiadczen interface {
	Odczytaj(ctx context.Context, byt string) (string, bool)
}

// ZSejfem wpina sejf poświadczeń — źródło stanu, nie treści.
func (a *adapterKanalow) ZSejfem(sejf sejfPoswiadczen) *adapterKanalow {
	a.sejf = sejf
	return a
}

// Sprawdz obsługuje `channel.check`.
func (a *adapterKanalow) Sprawdz(ctx context.Context,
	z shared.ChannelCheckRequest) (shared.ChannelCheckResponse, error) {

	kod := strings.TrimSpace(z.ChannelId)
	if kod == "" {
		return shared.ChannelCheckResponse{}, bladWskazaniaKanalu(
			"sprawdzenie bez wskazania kanału")
	}
	chwila := time.Now()
	if a.rejestr == nil {
		szczegol := "rdzeń nie ma wpiętego rejestru kanałów — nie ma czym wykonać sprawdzenia"
		return shared.ChannelCheckResponse{
			Reachable: false, CheckedAt: chwila.UnixMilli(), Detail: &szczegol,
		}, nil
	}
	if _, jest := a.rejestr.Kanal(kod); !jest {
		// Kanał, którego nie ma w rejestrze, nie jest odmową sprawdzenia: jest
		// odpowiedzią „nie odpowiada, bo go nie ma". Odmowa kazałaby oknu
		// pokazać błąd tam, gdzie Operator ma wiersz wyłączony albo źle wpisany.
		szczegol := "kanału nie ma w rejestrze kanałów rdzenia albo jest wyłączony"
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

// StanPoswiadczenia obsługuje `channel.credential.status`.
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
			"kanały: rdzeń nie ma wpiętego rejestru kanałów"))
	}
	kanal, err := a.repozytorium.PobierzPoKodzie(ctx, kod)
	if err != nil {
		return shared.ChannelCredentialStatusResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeNotFound, "kanały: nie ma kanału o identyfikatorze "+kod))
	}

	stan := shared.ChannelCredentialStatus{ChannelId: kod}
	odwolanie := strings.TrimSpace(wartoscTekstu(kanal.PoswiadczenieOdwolanie))
	if odwolanie == "" {
		// Kanał bez odwołania to kanał bez uwierzytelnienia — na przykład model
		// lokalny. To odpowiedź, nie brak: `present: false` bez rodzaju.
		return shared.ChannelCredentialStatusResponse{Status: stan}, nil
	}

	rodzaj := rodzajPoswiadczeniaKanalu(kanal)
	stan.Kind = &rodzaj
	zarzadca := zarzadcaPoswiadczeniaKanalu(kanal)
	stan.ManagedBy = &zarzadca

	if a.sejf != nil {
		if wartosc, jest := a.sejf.Odczytaj(ctx, odwolanie); jest {
			// Wartość służy WYŁĄCZNIE do rozstrzygnięcia, czy jest niepusta.
			// Poza tym warunkiem nie jest nigdzie użyta i nigdzie nie wychodzi.
			stan.Present = strings.TrimSpace(wartosc) != ""
		}
	}
	if !stan.Present {
		// Sejf nie zna odwołania — ale poświadczenie bywa też zmienną
		// środowiskową maszyny rdzenia, bo `credentialRef` jest właśnie jej
		// nazwą. Odpowiedź musi rozróżnić „nie ustawiono" od „ustawiono gdzie
		// indziej", więc mówi, gdzie rdzeń szukał.
		zarzadca = zarzadcaPoswiadczeniaKanalu(kanal) + " (odwołanie: " + odwolanie + ")"
		stan.ManagedBy = &zarzadca
	}
	return shared.ChannelCredentialStatusResponse{Status: stan}, nil
}

// przedrostekSprawdzeniaKanalu znakuje identyfikator wywołania sprawdzającego.
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

// zarzadcaPoswiadczeniaKanalu nazywa miejsce, w którym poświadczenie mieszka.
func zarzadcaPoswiadczeniaKanalu(kanal dane.Kanal) string {
	if kanal.KontoID != nil {
		return "rejestr kont platformy"
	}
	return "sejf poświadczeń rdzenia"
}

// bladWskazaniaKanalu nazywa niepoprawne żądanie czynności rejestru kanałów —
// błąd Operatora, nie rdzenia.
func bladWskazaniaKanalu(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"kanały: "+powod))
}
