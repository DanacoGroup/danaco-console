// Odpowiedzialność pliku: most modułu Translate do modelu w jednym miejscu —
// wybór kanału, złożenie polecenia przekładu, wywołanie modelu i rozgłoszenie
// zmiany panelu. Metody stoją na wspólnym `*adapterTlumaczenia`, którego typ,
// konstruktor i przedrostki deklaruje `adapter_modul_tlumaczenie.go`.
//
// Droga modelu (przekład panelu, rozpoznanie języka) to jedna odpowiedzialność
// wołana z trzech komend (`target.add`, `backtranslation.run`,
// `source.detect`), więc leży w jednym pliku i ma jeden opis granicy: brak
// kanału jest odmową wprost, nigdy pustym napisem udającym przekład.
//
// Kanał wskazuje Operator polem `channelId`, a rozstrzyga to `kanalZadania`:
//
//   - wskazania nie ma → kanał domyślny czynny (brak wskazania jest wskazaniem
//     na wartość domyślną, nie odmową);
//   - wskazanie wskazuje kanał czynny i gotowy → przekład idzie tym kanałem;
//   - wskazanie wskazuje kanał nieznany albo nieczynny → odmowa nazwana. Ciche
//     zejście na kanał domyślny wypełniłoby panel przekładem modelu, którego
//     Operator nie wybrał, i nic by o tym nie powiedziało.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// ZWyjsciem wpina szynę zdarzeń rdzenia — most, którym moduł rozgłasza
// `translate.translation.changed` po zmianie treści panelu. Wpięcie robi
// `zarejestrujTlumaczenie` (ten sam emiter, który dostają pozostałe moduły
// w `kompozycja.go`), więc montaż portów nie musi znać tej zależności.
// Nadajnik niepodłączony nie jest błędem: rdzeń tłumaczy także wtedy, gdy nikt
// nie słucha zdarzeń — `emiter.wyslij` sam odsiewa pusty nadajnik.
func (a *adapterTlumaczenia) ZWyjsciem(wyjscie *emiter) *adapterTlumaczenia {
	a.wyjscie = wyjscie
	return a
}

// kanalZadania rozstrzyga, którym kanałem pójdzie operacja modelu — jedno
// miejsce dla wszystkich komend modelowych modułu, żeby wskazanie Operatora
// znaczyło wszędzie to samo.
//
// Odmowy są trzy i każda mówi o czym innym: kanał nieznany, kanał znany, lecz
// wyłączony, i kanał włączony, lecz bez zbudowanego adaptera. Naprawia się je
// trzema różnymi ruchami (poprawić identyfikator, włączyć kanał, poprawić jego
// konfigurację), więc każda ma własne zdanie. Żadna nie schodzi po cichu na
// kanał domyślny, bo przekład wykonany innym modelem niż wskazany byłby
// przekładem, o którym Operator nie miałby skąd wiedzieć.
//
// Kanał „gotowy" znaczy tu to samo, co w `Rejestr.Kontrakt(true)`: wiersz
// czynny, który ma zbudowany adapter. Wiersz z `aktywny = 1` bez
// adaptera (rodzaj bez fabryki, fabryka odmówiła budowy) odmówi przy pierwszej
// turze, więc lepiej powiedzieć to teraz niż w połowie przekładu.
func (a *adapterTlumaczenia) kanalZadania(wskazany *string) (string, error) {
	if a.kanaly == nil {
		return "", bladBrakuKanalowTlumaczenia()
	}
	kod := ""
	if wskazany != nil {
		kod = strings.TrimSpace(*wskazany)
	}
	// Brak wskazania jest wskazaniem na kanał domyślny — droga sprzed pola
	// `channelId`, zachowana bez zmiany.
	if kod == "" {
		domyslny, ok := a.domyslnyKanalModelu()
		if !ok {
			return "", bladBrakuKanalowTlumaczenia()
		}
		return domyslny, nil
	}

	// Wykaz, a nie `Kontrakt(true)`: wykaz niesie także wiersze nieczynne, więc
	// da się odróżnić „nie ma takiego kanału" od „jest, ale wyłączony". Gdyby
	// szukać wyłącznie wśród czynnych, obie sytuacje zlałyby się w jedną odmowę
	// i Operator nie wiedziałby, czy pomylił identyfikator, czy zapomniał
	// włączyć kanał.
	for _, wiersz := range a.kanaly.Wykaz() {
		if wiersz.Identyfikator() != kod && strings.TrimSpace(wiersz.Kod) != kod {
			continue
		}
		if !wiersz.Aktywny {
			return "", protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
				"moduł Translate: wskazany kanał "+kod+" jest nieczynny — przekład nie zejdzie po cichu "+
					"na kanał domyślny; naprawa: włączyć kanał (channel.update enabled=true) "+
					"albo wskazać inny"))
		}
		if _, gotowy := a.kanaly.Kanal(wiersz.Identyfikator()); !gotowy {
			return "", protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
				"moduł Translate: wskazany kanał "+kod+" jest włączony, ale nie ma zbudowanego adaptera "+
					"— nie ma czym wołać modelu; naprawa: sprawdzić rodzaj i parametry kanału "+
					"(channel.list pokazuje wiersze pominięte)"))
		}
		return wiersz.Identyfikator(), nil
	}
	return "", protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"moduł Translate: nie ma kanału modelu o identyfikatorze "+kod+
			" — przekład nie zejdzie po cichu na kanał domyślny; naprawa: sprawdzić wykaz kanałów "+
			"(channel.list) albo pominąć pole channelId, żeby jechać kanałem domyślnym"))
}

// domyslnyKanalModelu ustala domyślny czynny kanał modelu — pierwszy wiersz
// rejestru czynny i gotowy do pracy (wzór `adapter_kolejki.go`:
// `rozwiazKanalPozycji`). Droga dla żądań bez `channelId`. Brak rejestru albo
// brak czynnego kanału znaczy „nie ma czym wołać modelu": druga wartość false,
// a wywołujący odmawia wprost, zamiast oddać wynik udający przekład.
func (a *adapterTlumaczenia) domyslnyKanalModelu() (string, bool) {
	if a.kanaly == nil {
		return "", false
	}
	czynne := a.kanaly.Kontrakt(true)
	if len(czynne) == 0 {
		return "", false
	}
	return czynne[0].Id, true
}

// przetlumaczModelem woła model, żeby przełożyć tekst źródłowy na język
// docelowy panelu, i oddaje sam przekład. Kanał rozstrzyga `kanalZadania` ze
// wskazania Operatora; jego brak albo wskazanie niedobre kończy się odmową
// wprost. Pusta odpowiedź modelu również jest odmową, bo pusty panel „po
// przekładzie" byłby atrapą.
//
// Słownik Operatora wchodzi dwa razy: raz do treści polecenia
// (`poleceniePrzekladu` dostaje `wiazania`), raz po odpowiedzi modelu jako
// mechaniczna podmiana terminów zostawionych w brzmieniu źródłowym
// (`zastosujTerminySlownika`, ta sama funkcja co w `glossary.apply`).
// Uzasadnienie obu dróg niesie nagłówek
// `adapter_modul_tlumaczenie_polecenia.go`. Nieudany odczyt słownika odmawia
// całego przekładu, bo przekład pomijający terminologię Operatora oddawałby po
// cichu co innego, niż zamówiono.
//
// Ton panelu też idzie do polecenia — inaczej kolumna `ton` (`panel.tone.set`)
// byłaby zapisem bez skutku.
func (a *adapterTlumaczenia) przetlumaczModelem(ctx context.Context,
	oknoKod, jezykDocelowy string, ton *string, tekstZrodlowy string,
	wskazanyKanal *string) (string, error) {

	kanal, err := a.kanalZadania(wskazanyKanal)
	if err != nil {
		return "", err
	}
	wiazania, err := a.wiazaniaSlownikaDlaJezyka(ctx, jezykDocelowy)
	if err != nil {
		return "", err
	}
	przeklad, err := a.zapytajModel(ctx, oknoKod, kanal,
		poleceniePrzekladu(jezykDocelowy, ton, wiazania, tekstZrodlowy))
	if err != nil {
		return "", err
	}
	przeklad = strings.TrimSpace(przeklad)
	if przeklad == "" {
		return "", bladTlumaczenia(errPustyPrzeklad)
	}
	// Siatka bezpieczeństwa słownika. Podmiana obejmuje wyłącznie odpowiedniki
	// języka docelowego — terminy nietykalne mają w wyniku zostać w brzmieniu
	// źródłowym, więc podmienianie ich byłoby złamaniem zakazu Operatora, a nie
	// jego pilnowaniem. Liczba podmian nie idzie nigdzie dalej: kontrakt
	// `target.add` nie ma pola na taką liczbę, a `ChangedCount` należy do
	// `glossary.apply`, nie do przekładu.
	przeklad, _ = zastosujTerminySlownika(przeklad, wiazania.odpowiedniki)
	return przeklad, nil
}

// przetlumaczZwrotnieModelem woła model o przekład panelu z powrotem na język
// źródłowy okna — obsługa `backtranslation.run`. Słownik tu nie wchodzi (powód
// przy `polecenieTlumaczeniaZwrotnego`), bo kontrola wierności, która sama
// naprawia terminologię, niczego nie kontroluje.
func (a *adapterTlumaczenia) przetlumaczZwrotnieModelem(ctx context.Context,
	oknoKod, jezykZrodlowy, tekstPanelu string, wskazanyKanal *string) (string, error) {

	kanal, err := a.kanalZadania(wskazanyKanal)
	if err != nil {
		return "", err
	}
	zwrotne, err := a.zapytajModel(ctx, oknoKod, kanal,
		polecenieTlumaczeniaZwrotnego(jezykZrodlowy, tekstPanelu))
	if err != nil {
		return "", err
	}
	zwrotne = strings.TrimSpace(zwrotne)
	if zwrotne == "" {
		return "", bladTlumaczenia(errPusteTlumaczenieZwrotne)
	}
	return zwrotne, nil
}

// rozpoznajJezykModelem pyta model o język tekstu i oddaje surową, przyciętą
// odpowiedź. Kanał jak wyżej — domyślny czynny, brak to odmowa wprost. Pusta
// odpowiedź modelu jest odmową, nie „nierozpoznanym językiem": rdzeń nie
// zgaduje ani nie udaje rozpoznania.
func (a *adapterTlumaczenia) rozpoznajJezykModelem(ctx context.Context, tekst string) (string, error) {
	kanal, ok := a.domyslnyKanalModelu()
	if !ok {
		return "", bladBrakuKanalowTlumaczenia()
	}
	// Rozpoznanie języka nie należy do żadnego okna — Zasięg okna zostaje pusty
	// (żądanie `source.detect` nie niesie okna, tylko sam tekst).
	odpowiedz, err := a.zapytajModel(ctx, "", kanal, polecenieRozpoznaniaJezyka(tekst))
	if err != nil {
		return "", err
	}
	odpowiedz = strings.TrimSpace(odpowiedz)
	if odpowiedz == "" {
		return "", bladTlumaczenia(errPusteRozpoznanie)
	}
	return odpowiedz, nil
}

// rozglosZmianePanelu rozgłasza `translate.translation.changed` po zmianie
// treści panelu. Zdarzenie idzie bez identyfikatora sesji: schemat okna
// tłumaczenia nie wiąże sesji (patrz `dane.OknoTlumaczenia`), więc zmiana
// rozgłasza się bez niej, zamiast nie rozgłaszać się wcale. Nadajnik
// niepodłączony jest odsiewany w `emiter.wyslij`, a nil-emiter — w jego
// odbiorniku nil, więc wywołanie jest bezpieczne bez wpiętej szyny.
func (a *adapterTlumaczenia) rozglosZmianePanelu(zmiana shared.ChangeKind, panel dane.PanelTlumaczenia) {
	a.wyjscie.wyslij(shared.EventTranslateTranslationChanged, "",
		shared.TranslateTranslationChangedEvent{Change: zmiana, Panel: zlozPanelTlumaczenia(panel)})
}

// polecenieRozpoznaniaJezyka składa treść wywołania modelu dla rozpoznania
// języka. Prompt zawęża odpowiedź do samej nazwy języka, żeby pole `Language`
// kontraktu nie niosło całego zdania modelu.
func polecenieRozpoznaniaJezyka(tekst string) string {
	return "Rozpoznaj język poniższego tekstu. Odpowiedz wyłącznie nazwą języka " +
		"(albo kodem ISO 639-1), bez zdania, bez komentarza.\n\nTekst:\n" + tekst
}
