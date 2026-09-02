// Wypełnienie portu Mowa dwiema komendami obszaru speech: przełożenie żądań
// kontraktu na zlecenia silnika mowy oraz przełożenie jego typowanych odmów na
// kody kontraktu, ze złożeniem trójki okno, zasady i obszar zasięgu platformy
// dla każdego wywołania.
package core

import (
	"context"
	"errors"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/mowa"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

type adapterMowy struct {
	uruchamiacz   session.Uruchamiacz
	dziennik      mowa.Dziennik
	rozstrzygacz  *konfig.Rozstrzygacz
	katalog       *KatalogRoboczy
	nagrania      dane.RepozytoriumNagranMowy
	katalogDanych string
	konfiguracja  dane.RepozytoriumKonfiguracji
	nasluchy      *rejestrNasluchow
	nadajnik      *emiter
}

func nowyAdapterMowy(uruchamiacz session.Uruchamiacz) *adapterMowy {
	return &adapterMowy{uruchamiacz: uruchamiacz, nasluchy: nowyRejestrNasluchow()}
}

func (a *adapterMowy) ZMagazynemNagran(nagrania dane.RepozytoriumNagranMowy,
	katalogDanych string, konfiguracja dane.RepozytoriumKonfiguracji) *adapterMowy {

	a.nagrania = nagrania
	a.katalogDanych = katalogDanych
	a.konfiguracja = konfiguracja
	return a
}

func (a *adapterMowy) ZWyjsciem(e *emiter) *adapterMowy {
	a.nadajnik = e
	return a
}

func (a *adapterMowy) ZDziennikiem(d mowa.Dziennik) *adapterMowy {
	a.dziennik = d
	return a
}

func (a *adapterMowy) ZIzolacja(rozstrzygacz *konfig.Rozstrzygacz, katalog *KatalogRoboczy) *adapterMowy {
	a.rozstrzygacz, a.katalog = rozstrzygacz, katalog
	return a
}

func (a *adapterMowy) Gotowosc(ctx context.Context,
	_ shared.SpeechAvailabilityGetRequest) (shared.SpeechAvailabilityGetResponse, error) {

	ustawienia := a.ustawienia(ctx)
	okno, zasady, obszar := a.zasiegPlatformy(ctx)

	odsluchGotowy, odsluchPowod := gotowoscOdsluchu()

	dostepnosc, err := a.silnik(ustawienia).Dostepnosc(ctx, okno, zasady, obszar)
	if err != nil {
		return shared.SpeechAvailabilityGetResponse{
			Available:          false,
			Model:              wskaznikTekstu(ustawienia.Model),
			Reason:             wskaznikTekstu(err.Error()),
			SynthesisAvailable: wskaznikPrawdy(odsluchGotowy),
			SynthesisReason:    wskaznikTekstu(odsluchPowod),
		}, nil
	}
	przyjmowanie := a.nagrania != nil && strings.TrimSpace(a.katalogDanych) != ""
	return shared.SpeechAvailabilityGetResponse{
		Available:          dostepnosc.Gotowy,
		UploadAvailable:    wskaznikPrawdy(przyjmowanie),
		WakeWordAvailable:  wskaznikPrawdy(dostepnosc.Gotowy),
		ListenAvailable:    wskaznikPrawdy(dostepnosc.Gotowy && przyjmowanie),
		Python:             wskaznikTekstu(dostepnosc.Python),
		Engine:             wskaznikTekstu(dostepnosc.Silnik),
		Model:              wskaznikTekstu(pierwszyNiepustyTekst(dostepnosc.Model, ustawienia.Model)),
		Reason:             wskaznikTekstu(dostepnosc.Powod),
		SynthesisAvailable: wskaznikPrawdy(odsluchGotowy),
		SynthesisReason:    wskaznikTekstu(odsluchPowod),
	}, nil
}

func gotowoscOdsluchu() (bool, string) {
	if _, err := dobierzSyntezator("pl"); err != nil {
		return false, err.Error()
	}
	return true, ""
}

func (a *adapterMowy) Przepisz(ctx context.Context,
	z shared.SpeechTranscribeRequest) (shared.SpeechTranscribeResponse, error) {

	odnosnik := strings.TrimSpace(z.AudioRef)
	if odnosnik == "" {
		return shared.SpeechTranscribeResponse{}, bladZadaniaMowy(
			"transkrypcja bez odnośnika nagrania nie ma czego przepisać")
	}
	// Odnośnik pod katalogiem nagrań przyjętych wydaje wyłącznie rejestr
	// zawężony do konta żądania. Ścieżka spoza tego katalogu to plik Operatora
	// (`research.source.transcribe`, dyktowanie z dysku) i idzie do silnika
	// bez zmiany.
	if a.wKataloguNagranPrzyjetych(odnosnik) {
		if a.nagrania == nil {
			return shared.SpeechTranscribeResponse{}, bladZapleczaNagran(
				"serwer nie ma wpiętego rejestru nagrań")
		}
		wiersz, err := a.nagrania.NagranieMowyPoSciezce(ctx, odnosnik)
		if err != nil {
			return shared.SpeechTranscribeResponse{}, bladObcegoNagrania(odnosnik)
		}
		odnosnik = wiersz.Sciezka
	}

	ustawienia := a.ustawienia(ctx)
	okno, zasady, obszar := a.zasiegPlatformy(ctx)

	transkrypcja, err := a.silnik(ustawienia).Transkrybuj(ctx, mowa.Zlecenie{
		Odnosnik: odnosnik,
		Jezyk:    wartoscTekstu(z.Language),
		Model:    wartoscTekstu(z.Model),
		Okno:     okno,
		Zasady:   zasady,
		Obszar:   obszar,
	})
	if err != nil {
		return shared.SpeechTranscribeResponse{}, bladSilnikaMowy(err)
	}

	return shared.SpeechTranscribeResponse{
		Processed:  true,
		Transcript: transkrypcja.Tekst,
		Characters: transkrypcja.Znakow,
		DurationMs: int(transkrypcja.TrwanieMs),
		Model:      transkrypcja.Model,
		Language:   wskaznikTekstu(transkrypcja.Jezyk),
	}, nil
}

func (a *adapterMowy) silnik(ustawienia mowa.Ustawienia) *mowa.Silnik {
	return mowa.NowySilnik(a.uruchamiacz).ZUstawieniami(ustawienia).ZDziennikiem(a.dziennik)
}

func (a *adapterMowy) ustawienia(ctx context.Context) mowa.Ustawienia {
	komplet := mowa.UstawieniaDomyslne()
	if a.rozstrzygacz == nil {
		return komplet
	}
	for _, klucz := range []string{mowa.KluczProgram, mowa.KluczModel,
		mowa.KluczJezyk, mowa.KluczKatalogModeli} {

		wynik := a.rozstrzygacz.Rozstrzygnij(konfig.Kontekst{KontoOperatora: dane.KontoOperatora(ctx)}, klucz)
		if wynik.Pochodzenie == konfig.PochodzenieNieznane {
			continue
		}
		komplet = mowa.Nanies(komplet, klucz, wynik.Wartosc)
	}
	return komplet
}

func (a *adapterMowy) zasiegPlatformy(ctx context.Context) (session.Okno, session.Zasady, session.Obszar) {
	okno := session.Okno{Ustawienia: session.Ustawienia{
		SrodowiskoWykonania: shared.ExecutionEnvCore,
	}}
	kontekst := konfig.Kontekst{KontoOperatora: dane.KontoOperatora(ctx)}
	zasady := session.Zasady{}
	if a.rozstrzygacz != nil {
		zasady = ZasadyIzolacji(a.rozstrzygacz, kontekst)
	}
	obszar := session.Obszar{}
	if a.katalog != nil {
		obszar = ObszarOkna(a.katalog.Ustal(kontekst, ""), "")
	}
	return okno, zasady, obszar
}

func bladSilnikaMowy(err error) error {
	var brakNagrania *mowa.BrakNagrania
	if errors.As(err, &brakNagrania) {
		return bladZadaniaMowy(brakNagrania.Error())
	}

	var brakPomocnika *mowa.BrakPomocnika
	var brakInterpretera *mowa.BrakInterpretera
	var brakSilnika *mowa.BrakSilnika
	if errors.As(err, &brakPomocnika) ||
		errors.As(err, &brakInterpretera) ||
		errors.As(err, &brakSilnika) {
		return odmowaMowy(shared.ErrorCodeChannelUnavailable, err.Error())
	}

	if errors.Is(err, session.ErrIzolacja) {
		return odmowaMowy(shared.ErrorCodePermissionDenied, err.Error())
	}
	return odmowaMowy(shared.ErrorCodeInternalError, err.Error())
}

func bladZadaniaMowy(powod string) error {
	return odmowaMowy(shared.ErrorCodeValidationFailed, powod)
}

func odmowaMowy(kod shared.ErrorCode, powod string) error {
	return protocol.JakoError(protocol.NowyBlad(kod, "silnik mowy: "+powod))
}

func pierwszyNiepustyTekst(wartosci ...string) string {
	for _, wartosc := range wartosci {
		if strings.TrimSpace(wartosc) != "" {
			return wartosc
		}
	}
	return ""
}
