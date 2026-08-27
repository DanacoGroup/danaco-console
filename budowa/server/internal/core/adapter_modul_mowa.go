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

// adapterMowy wypełnia port Mowa: łączy uruchamiacz procesu pomocnika, trwały
// dziennik transkrypcji, źródła zasad izolacji i obszaru roboczego oraz
// zależności czterech komend dobudowanych obok rozpoznawania mowy.
type adapterMowy struct {
	// uruchamiacz jest portem warstwy kanału, jedyną drogą startu procesu
	// pomocnika w drzewie.
	uruchamiacz session.Uruchamiacz
	// dziennik daje trwały ślad transkrypcji; brak zależności nie zatrzymuje
	// rozpoznawania mowy.
	dziennik mowa.Dziennik
	// rozstrzygacz i katalog składają zasady izolacji i obszar zasięgu
	// platformy dla uruchomienia.
	rozstrzygacz *konfig.Rozstrzygacz
	katalog      *KatalogRoboczy
	// nagrania, katalogDanych i konfiguracja obsługują cztery komendy
	// dobudowane obok transkrypcji.
	nagrania      dane.RepozytoriumNagranMowy
	katalogDanych string
	konfiguracja  dane.RepozytoriumKonfiguracji
	// nasluchy trzyma nasłuchy ciągłe okien, a nadajnik ogłasza to, co rdzeń
	// w nich usłyszał.
	nasluchy *rejestrNasluchow
	nadajnik *emiter
}

// nowyAdapterMowy wiąże port z uruchamiaczem procesów i zakłada pusty rejestr
// nasłuchów ciągłych okien, gotowy do wpięcia pozostałych zależności metodami
// budującymi.
func nowyAdapterMowy(uruchamiacz session.Uruchamiacz) *adapterMowy {
	return &adapterMowy{uruchamiacz: uruchamiacz, nasluchy: nowyRejestrNasluchow()}
}

// ZMagazynemNagran wpina rejestr nagrań, katalog danych rdzenia i magazyn
// konfiguracji — trzy zależności czterech komend dobudowanych obok transkrypcji
// (`adapter_mowa_nagrania.go`, `adapter_mowa_wybudzenie.go`).
func (a *adapterMowy) ZMagazynemNagran(nagrania dane.RepozytoriumNagranMowy,
	katalogDanych string, konfiguracja dane.RepozytoriumKonfiguracji) *adapterMowy {

	a.nagrania = nagrania
	a.katalogDanych = katalogDanych
	a.konfiguracja = konfiguracja
	return a
}

// ZWyjsciem wpina nadajnik zdarzeń, przez który nasłuch ciągły okna ogłasza to,
// co rdzeń usłyszał, kiedy proces mowy działa w tle.
func (a *adapterMowy) ZWyjsciem(e *emiter) *adapterMowy {
	a.nadajnik = e
	return a
}

// ZDziennikiem podpina trwały ślad transkrypcji nad bazą rdzenia, z którego
// korzysta silnik rozpoznawania mowy przy każdym wywołaniu komendy.
func (a *adapterMowy) ZDziennikiem(d mowa.Dziennik) *adapterMowy {
	a.dziennik = d
	return a
}

// ZIzolacja podpina rozstrzygacz zasięgu i ustalacz katalogu roboczego — dwa
// źródła, z których powstają zasady i obszar egzekwowane przy uruchomieniu.
func (a *adapterMowy) ZIzolacja(rozstrzygacz *konfig.Rozstrzygacz, katalog *KatalogRoboczy) *adapterMowy {
	a.rozstrzygacz, a.katalog = rozstrzygacz, katalog
	return a
}

// Gotowosc obsługuje speech.availability.get. Brak silnika jest odpowiedzią,
// nie odmową, bo komenda zadaje jedno pytanie: czy rysować mikrofon. Powód
// niesie trójczęściowy komunikat pakietu wraz ze wskazaniem naprawy.
func (a *adapterMowy) Gotowosc(ctx context.Context,
	_ shared.SpeechAvailabilityGetRequest) (shared.SpeechAvailabilityGetResponse, error) {

	ustawienia := a.ustawienia(ctx)
	okno, zasady, obszar := a.zasiegPlatformy()

	// Odsłuch mierzy się osobno od dyktowania, bo jedzie innym łańcuchem
	// silnika mowy.
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
	// Przyjęcie nagrania zależy wyłącznie od magazynu rdzenia i działa nawet
	// bez silnika mowy.
	przyjmowanie := a.nagrania != nil && strings.TrimSpace(a.katalogDanych) != ""
	return shared.SpeechAvailabilityGetResponse{
		Available:         dostepnosc.Gotowy,
		UploadAvailable:   wskaznikPrawdy(przyjmowanie),
		WakeWordAvailable: wskaznikPrawdy(dostepnosc.Gotowy),
		ListenAvailable:   wskaznikPrawdy(dostepnosc.Gotowy && przyjmowanie),
		Python:            wskaznikTekstu(dostepnosc.Python),
		Engine:            wskaznikTekstu(dostepnosc.Silnik),
		// Model zastany na dysku, a przy jego braku model ustawiony w
		// konfiguracji rdzenia.
		Model:              wskaznikTekstu(pierwszyNiepustyTekst(dostepnosc.Model, ustawienia.Model)),
		Reason:             wskaznikTekstu(dostepnosc.Powod),
		SynthesisAvailable: wskaznikPrawdy(odsluchGotowy),
		SynthesisReason:    wskaznikTekstu(odsluchPowod),
	}, nil
}

// gotowoscOdsluchu mierzy, czy synteza mowy ruszy tu i teraz, drogą
// dobierzSyntezator, żeby pomiar nie rozjechał się z wykonaniem. Język jest
// polski, bo produkt jest polskojęzyczny. Pusty powód znaczy: odsłuch gotowy.
func gotowoscOdsluchu() (bool, string) {
	if _, err := dobierzSyntezator("pl"); err != nil {
		return false, err.Error()
	}
	return true, ""
}

// Przepisz obsługuje speech.transcribe. Cisza jest wynikiem pomiaru, nie
// błędem: pole processed rozdziela przetworzenie z tekstem, przetworzenie bez
// mowy i nieudane przetworzenie.
func (a *adapterMowy) Przepisz(ctx context.Context,
	z shared.SpeechTranscribeRequest) (shared.SpeechTranscribeResponse, error) {

	if strings.TrimSpace(z.AudioRef) == "" {
		return shared.SpeechTranscribeResponse{}, bladZadaniaMowy(
			"transkrypcja bez odnośnika nagrania nie ma czego przepisać")
	}

	ustawienia := a.ustawienia(ctx)
	okno, zasady, obszar := a.zasiegPlatformy()

	transkrypcja, err := a.silnik(ustawienia).Transkrybuj(ctx, mowa.Zlecenie{
		Odnosnik: z.AudioRef,
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
		// Przetworzono, bo silnik nie oddał odmowy, niezależnie od długości
		// tekstu.
		Processed:  true,
		Transcript: transkrypcja.Tekst,
		Characters: transkrypcja.Znakow,
		DurationMs: int(transkrypcja.TrwanieMs),
		Model:      transkrypcja.Model,
		Language:   wskaznikTekstu(transkrypcja.Jezyk),
	}, nil
}

// silnik składa silnik mowy na nastawach obowiązujących dla tego wywołania.
// Powód budowy per wywołanie stoi w nagłówku pliku.
func (a *adapterMowy) silnik(ustawienia mowa.Ustawienia) *mowa.Silnik {
	return mowa.NowySilnik(a.uruchamiacz).ZUstawieniami(ustawienia).ZDziennikiem(a.dziennik)
}

// ustawienia czyta cztery nastawy silnika z katalogu ustawień w zasięgu
// platformy. Brak rozstrzygacza daje komplet domyślny, a nie odmowę: brak
// ustawienia jest wskazaniem na wartość domyślną.
func (a *adapterMowy) ustawienia(context.Context) mowa.Ustawienia {
	komplet := mowa.UstawieniaDomyslne()
	if a.rozstrzygacz == nil {
		return komplet
	}
	for _, klucz := range []string{mowa.KluczProgram, mowa.KluczModel,
		mowa.KluczJezyk, mowa.KluczKatalogModeli} {

		wynik := a.rozstrzygacz.Rozstrzygnij(konfig.Kontekst{}, klucz)
		if wynik.Pochodzenie == konfig.PochodzenieNieznane {
			continue
		}
		komplet = mowa.Nanies(komplet, klucz, wynik.Wartosc)
	}
	return komplet
}

// zasiegPlatformy składa trójkę okno–zasady–obszar dla zasięgu platformy.
// Rozstrzygnięcie i jego uzasadnienie stoją w nagłówku pliku.
func (a *adapterMowy) zasiegPlatformy() (session.Okno, session.Zasady, session.Obszar) {
	okno := session.Okno{Ustawienia: session.Ustawienia{
		SrodowiskoWykonania: shared.ExecutionEnvCore,
	}}
	zasady := session.Zasady{}
	if a.rozstrzygacz != nil {
		zasady = ZasadyIzolacji(a.rozstrzygacz, konfig.Kontekst{})
	}
	obszar := session.Obszar{}
	if a.katalog != nil {
		obszar = ObszarOkna(a.katalog.Ustal(konfig.Kontekst{}, ""), "")
	}
	return okno, zasady, obszar
}

// bladSilnikaMowy przekłada typowane odmowy pakietu mowa na trzy kody
// kontraktu, dobierane osobno dla braku danych, braku pomocnika i naruszenia
// izolacji.
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

// bladZadaniaMowy znakuje wadę żądania kodem kontraktu validation_failed, gdy
// przysłane dane nie pozwalają wykonać transkrypcji.
func bladZadaniaMowy(powod string) error {
	return odmowaMowy(shared.ErrorCodeValidationFailed, powod)
}

// odmowaMowy składa odmowę obszaru z kodem kontraktu i trójczęściową treścią.
// Przedrostek nazywa obszar, żeby czytelnik odmowy wiedział, kto odmówił, zanim
// przeczyta dlaczego.
func odmowaMowy(kod shared.ErrorCode, powod string) error {
	return protocol.JakoError(protocol.NowyBlad(kod, "silnik mowy: "+powod))
}

// pierwszyNiepustyTekst oddaje pierwszą wartość spośród podanych, która niesie
// niepusty tekst po przycięciu białych znaków.
func pierwszyNiepustyTekst(wartosci ...string) string {
	for _, wartosc := range wartosci {
		if strings.TrimSpace(wartosc) != "" {
			return wartosc
		}
	}
	return ""
}
