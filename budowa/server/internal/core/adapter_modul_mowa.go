// Wypełnienie portu Mowa dwiema komendami obszaru `speech.*`: przełożenie żądań
// kontraktu na zlecenia silnika `mowa` i przełożenie jego typowanych odmów na
// kody kontraktu.
//
// ── Skąd biorą się okno, zasady i obszar ────────────────────────────────────
// `mowa.Silnik` wymaga `session.Okno`, `session.Zasady` i `session.Obszar`, bo
// pomocnik transkrypcji startuje tym samym uruchamiaczem i przez tę samą bramę
// izolacji, co każdy inny proces drzewa. Terminal i Developer biorą tę trójkę
// z okna żądania: rejestr okien daje `session.Okno`, `ZasadyIzolacji` nad
// rozstrzygaczem daje zasady dla `konfig.Kontekst{Okno: …}`, a `ObszarOkna` nad
// ustalaczem katalogu roboczego daje obszar.
//
// Żądania `speech.*` okna nie niosą: `SpeechAvailabilityGetRequest` nie ma pól,
// a `SpeechTranscribeRequest` niesie wyłącznie `audioRef`, `language` i `model`.
// Rodzina `speech.*` jest zdolnością platformy, nie okna — pyta o to, czy dana
// maszyna umie rozpoznać mowę. Trójka składa się więc w zasięgu platformy, czyli
// dla pustego `konfig.Kontekst{}`, tą samą drogą co dla okna:
//
//	zasady := ZasadyIzolacji(rozstrzygacz, konfig.Kontekst{})
//	obszar := ObszarOkna(katalog.Ustal(konfig.Kontekst{}, ""), "")
//
// Pusty kontekst zasięgu jest poprawnym adresem najszerszego poziomu:
// rozstrzygacz oddaje wtedy politykę platformy, a ustalacz — katalog roboczy
// platformy. Wpisane z ręki `session.Zasady{}` i `session.Obszar{}` znaczyłyby
// „izolacja wyłączona" niezależnie od ustawień.
//
// Okno jest jedyną wartością, którą adapter wypełnia sam:
// `session.Okno{SrodowiskoWykonania: shared.ExecutionEnvCore}`. `audioRef`
// wskazuje ścieżkę na maszynie silnika, więc pomocnik musi ruszyć na hoście
// rdzenia; puste pole środowiska dałoby ten sam rozruch gałęzią `case ""`
// w `injection/uruchamiacz_okna.go`, ale bez zapisanego wskazania.
// Identyfikatora okna nie ma skąd wziąć, więc wpis dziennika transkrypcji nie
// dostaje odnośnika okna.
//
// ── Silnik powstaje na każde wywołanie ──────────────────────────────────────
// `mowa.Silnik.ZUstawieniami` mutuje byt, więc jedna instancja współdzielona
// przez równoległe żądania oznaczałaby wyścig o nastawy: jedno żądanie ustawia
// model, drugie go podmienia, pierwsze rozpoznaje cudzym. Przy okazji nastawy są
// świeże — zmiana `mowa_model` komendą `config.set` obowiązuje od następnej
// transkrypcji, bez restartu rdzenia.
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

// adapterMowy wypełnia port Mowa.
type adapterMowy struct {
	// uruchamiacz jest portem warstwy kanału — jedyną drogą startu procesu
	// w drzewie. Bez niego moduł nie ruszy pomocnika.
	uruchamiacz session.Uruchamiacz
	// dziennik daje trwały ślad transkrypcji. Zależność opcjonalna: silnik bez
	// dziennika rozpoznaje mowę tak samo, traci wyłącznie ślad.
	dziennik mowa.Dziennik
	// rozstrzygacz i katalog składają zasady izolacji i obszar zasięgu platformy
	// — te same dwa źródła, z których korzystają Terminal i Developer.
	rozstrzygacz *konfig.Rozstrzygacz
	katalog      *KatalogRoboczy
	// nagrania, katalogDanych i konfiguracja obsługują cztery komendy dobudowane
	// obok transkrypcji: przyjęcie i oddanie bajtów nagrania oraz nastawę
	// wybudzania. Zależności są opcjonalne — bez nich te komendy odmawiają,
	// nazywając brak, a rozpoznawanie mowy pracuje bez zmian.
	nagrania      dane.RepozytoriumNagranMowy
	katalogDanych string
	konfiguracja  dane.RepozytoriumKonfiguracji
	// nasluchy trzyma nasłuchy ciągłe okien, a nadajnik ogłasza to, co rdzeń
	// w nich usłyszał.
	nasluchy *rejestrNasluchow
	nadajnik *emiter
}

// nowyAdapterMowy wiąże port z uruchamiaczem procesów.
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

// ZWyjsciem wpina nadajnik zdarzeń nasłuchu ciągłego.
func (a *adapterMowy) ZWyjsciem(e *emiter) *adapterMowy {
	a.nadajnik = e
	return a
}

// ZDziennikiem podpina trwały ślad transkrypcji nad bazą rdzenia.
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

// Gotowosc obsługuje `speech.availability.get`.
//
// Brak silnika jest odpowiedzią, nie odmową — tak stanowi kontrakt tej komendy.
// Silnik rozstrzyga to po swojej stronie: brak interpretera i brak biblioteki
// wracają jako `Gotowy` fałszywe z powodem.
//
// Błędem zostaje u niego wyłącznie nieczytelna odpowiedź pomocnika i tę jedną
// adapter również sprowadza do `available=false` z powodem, zamiast oddać
// odmowę: komenda zadaje jedno pytanie — czy rysować mikrofon — a odmowa
// zostawiłaby klienta bez odpowiedzi na nie. Powód niesie trójczęściowy
// komunikat pakietu wraz ze wskazaniem naprawy.
func (a *adapterMowy) Gotowosc(ctx context.Context,
	_ shared.SpeechAvailabilityGetRequest) (shared.SpeechAvailabilityGetResponse, error) {

	ustawienia := a.ustawienia(ctx)
	okno, zasady, obszar := a.zasiegPlatformy()

	// Odsłuch liczony jest RAZ i osobno od dyktowania, bo jedzie innym łańcuchem:
	// piper z głosem .onnx albo espeak-ng, nie python z faster-whisper. Nawet gdy
	// dyktowanie odmawia (brak Pythona), odsłuch bywa gotowy — dlatego wynik idzie
	// do obu gałęzi, a nie tylko do udanej. Bez tego komenda meldowałaby o odsłuchu
	// wyłącznie wtedy, gdy działa dyktowanie, choć to dwie niezależne zdolności.
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
	// Trzy dobudowane zdolności rodziny meldują się osobno, bo osobno znikają.
	// Przyjęcie nagrania (`speech.audio.upload`) zależy WYŁĄCZNIE od magazynu
	// rdzenia i działa nawet bez silnika mowy — bajty da się odłożyć i odsłuchać
	// bez rozpoznawania czegokolwiek. Wybudzenie i nasłuch ciągły rozpoznają
	// każdy odcinek, więc znikają razem z silnikiem.
	przyjmowanie := a.nagrania != nil && strings.TrimSpace(a.katalogDanych) != ""
	return shared.SpeechAvailabilityGetResponse{
		Available:         dostepnosc.Gotowy,
		UploadAvailable:   wskaznikPrawdy(przyjmowanie),
		WakeWordAvailable: wskaznikPrawdy(dostepnosc.Gotowy),
		ListenAvailable:   wskaznikPrawdy(dostepnosc.Gotowy && przyjmowanie),
		Python:            wskaznikTekstu(dostepnosc.Python),
		Engine:            wskaznikTekstu(dostepnosc.Silnik),
		// Model zastany na dysku, a przy jego braku — model ustawiony. Kontrakt
		// pyta o ten drugi, pomocnik odpowiada tym pierwszym; oddanie pustki, gdy
		// wag jeszcze nie pobrano, gubiłoby nastawę widoczną w oknie konfiguracji.
		Model:              wskaznikTekstu(pierwszyNiepustyTekst(dostepnosc.Model, ustawienia.Model)),
		Reason:             wskaznikTekstu(dostepnosc.Powod),
		SynthesisAvailable: wskaznikPrawdy(odsluchGotowy),
		SynthesisReason:    wskaznikTekstu(odsluchPowod),
	}, nil
}

// gotowoscOdsluchu mierzy, czy synteza mowy ruszy tu i teraz.
//
// Używa `dobierzSyntezator` — tej samej drogi, którą idzie faktyczny odsłuch —
// więc pomiar nie może rozejść się z wykonaniem: jeżeli dobór silnika kończy się
// odmową, odsłuch odmówi tak samo, a jego powód jest tym, co Operator zobaczy.
// Zgodność mierzonego z wykonywanym jest tu warunkiem sensu: komenda, która
// mówi „gotowy" o rzeczy, która za chwilę odmówi, jest gorsza niż jej brak.
//
// Język jest polski, bo produkt jest polskojęzyczny, a odsłuch czyta panele
// właśnie po polsku; głos innego języka i tak nie przeczytałby polskiego panelu
// naturalnie. Pusty powód znaczy: odsłuch gotowy.
func gotowoscOdsluchu() (bool, string) {
	if _, err := dobierzSyntezator("pl"); err != nil {
		return false, err.Error()
	}
	return true, ""
}

// Przepisz obsługuje `speech.transcribe`.
//
// Cisza jest wynikiem pomiaru, nie błędem. Pole `processed` rozdziela trzy stany
// i tylko trzeci jest odmową: przetworzono z tekstem, przetworzono bez mowy
// (`mowa.StanBezMowy` → `processed=true`, `transcript` pusty) oraz nie
// przetworzono.
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
		// Przetworzono, bo silnik nie oddał odmowy. Wartość nie wynika z długości
		// tekstu: nagranie bez mowy też jest przetworzone.
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

// bladSilnikaMowy przekłada typowane odmowy pakietu `mowa` na kody kontraktu —
// cztery przypadki, trzy kody:
//
//   - `mowa.BrakNagrania` → `validation_failed`. Jedyna z odmów wywołana daną
//     przysłaną przez klienta: `audioRef` nie wskazuje pliku, który da się
//     przepisać (nie ma go, jest pusty albo ma format spoza wykazu). Kod
//     nieponawialny — to samo żądanie powtórzone da to samo.
//
//   - `mowa.BrakPomocnika` → `channel_unavailable`. Nie ma interpretera albo nie
//     ma skryptu pomocnika. Żądanie było poprawne, więc nie `validation_failed`;
//     `not_found` mówiłby o nagraniu, nie o instalacji; `internal_error`
//     znaczyłby wadę produktu, a produkt mówi wprost, czego dołożyć. Kod
//     ponawialny: po naprawie z komunikatu to samo żądanie się powiedzie.
//
//   - `mowa.BrakInterpretera` i `mowa.BrakSilnika` → `channel_unavailable`, ten
//     sam kod. Katalog `ErrorCode` nie ma pozycji odróżniającej „nie ma
//     interpretera" od „interpreter jest, brakuje w nim biblioteki".
//     Rozróżnienie, którego wymaga naprawa, niesie TREŚĆ odmowy — i musi je
//     nieść, bo trzy ogniwa łańcucha naprawia się trzema różnymi czynnościami:
//     dołożeniem katalogu pomocników, instalacją Pythona 3 i instalacją
//     `faster-whisper`. Odmowa przypisująca brak niewłaściwemu ogniwu prowadzi
//     Operatora do naprawy bezskutecznej.
//
//   - naruszenie izolacji (`session.ErrIzolacja`) → `permission_denied`. Punkt
//     izolacji zatrzymał uruchomienie; tak samo znakuje je Terminal.
//
// Odmowa nierozpoznana schodzi na `internal_error`: to przypadek, którego rdzeń
// nie przewidział, i ma się zgłosić jako taki, a nie udawać znany.
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

// bladZadaniaMowy znakuje wadę żądania kodem kontraktu.
func bladZadaniaMowy(powod string) error {
	return odmowaMowy(shared.ErrorCodeValidationFailed, powod)
}

// odmowaMowy składa odmowę obszaru z kodem kontraktu i trójczęściową treścią.
// Przedrostek nazywa obszar, żeby czytelnik odmowy wiedział, kto odmówił, zanim
// przeczyta dlaczego.
func odmowaMowy(kod shared.ErrorCode, powod string) error {
	return protocol.JakoError(protocol.NowyBlad(kod, "silnik mowy: "+powod))
}

// pierwszyNiepustyTekst oddaje pierwszą wartość, która coś niesie.
func pierwszyNiepustyTekst(wartosci ...string) string {
	for _, wartosc := range wartosci {
		if strings.TrimSpace(wartosc) != "" {
			return wartosc
		}
	}
	return ""
}
