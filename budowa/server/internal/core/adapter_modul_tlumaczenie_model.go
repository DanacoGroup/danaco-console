// Odpowiedzialność pliku: most modułu Translate do modelu — wybór kanału, złożenie polecenia przekładu, wywołanie modelu i rozgłoszenie zmiany panelu, dla przekładu, przekładu zwrotnego i rozpoznania języka.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// ZWyjsciem wpina szynę zdarzeń rdzenia — most, którym moduł rozgłasza translate.translation.changed po zmianie treści panelu.
func (a *adapterTlumaczenia) ZWyjsciem(wyjscie *emiter) *adapterTlumaczenia {
	a.wyjscie = wyjscie
	return a
}

// kanalZadania rozstrzyga, którym kanałem pójdzie operacja modelu, na podstawie wskazania Operatora, z odmową wprost przy kanale nieznanym, wyłączonym albo bez zbudowanego adaptera.
func (a *adapterTlumaczenia) kanalZadania(wskazany *string) (string, error) {
	if a.kanaly == nil {
		return "", bladBrakuKanalowTlumaczenia()
	}
	kod := ""
	if wskazany != nil {
		kod = strings.TrimSpace(*wskazany)
	}
	// Brak wskazania jest wskazaniem na kanał domyślny, drogą sprzed pola channelId.
	if kod == "" {
		domyslny, ok := a.domyslnyKanalModelu()
		if !ok {
			return "", bladBrakuKanalowTlumaczenia()
		}
		return domyslny, nil
	}

	// Wykaz niesie też wiersze nieczynne, więc odróżnia brak kanału od kanału wyłączonego.
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

// domyslnyKanalModelu ustala domyślny czynny kanał modelu: pierwszy wiersz rejestru czynny i gotowy do pracy, droga dla żądań bez channelId.
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

// przetlumaczModelem woła model, żeby przełożyć tekst źródłowy na język docelowy panelu, stosując słownik Operatora w poleceniu i po odpowiedzi, oraz oddaje sam przekład.
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
	// Podmiana słownika obejmuje odpowiedniki celu; terminy nietykalne zostają w brzmieniu źródłowym.
	przeklad, _ = zastosujTerminySlownika(przeklad, wiazania.odpowiedniki)
	return przeklad, nil
}

// przetlumaczZwrotnieModelem woła model o przekład panelu z powrotem na język źródłowy okna, obsługując backtranslation.run bez udziału słownika.
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

// rozpoznajJezykModelem pyta model o język tekstu i oddaje surową, przyciętą odpowiedź; pusta odpowiedź modelu jest odmową.
func (a *adapterTlumaczenia) rozpoznajJezykModelem(ctx context.Context, tekst string) (string, error) {
	kanal, ok := a.domyslnyKanalModelu()
	if !ok {
		return "", bladBrakuKanalowTlumaczenia()
	}
	// Rozpoznanie języka nie należy do żadnego okna: żądanie source.detect niesie sam tekst.
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

// rozglosZmianePanelu rozgłasza translate.translation.changed po zmianie treści panelu, bez identyfikatora sesji, bo schemat okna tłumaczenia jej nie wiąże.
func (a *adapterTlumaczenia) rozglosZmianePanelu(ctx context.Context, zmiana shared.ChangeKind, panel dane.PanelTlumaczenia) {
	a.wyjscie.wyslijDoKonta(ctx, shared.EventTranslateTranslationChanged, "",
		shared.TranslateTranslationChangedEvent{Change: zmiana, Panel: zlozPanelTlumaczenia(panel)})
}

// polecenieRozpoznaniaJezyka składa treść wywołania modelu dla rozpoznania
// języka. Prompt zawęża odpowiedź do samej nazwy języka, żeby pole `Language`
// kontraktu nie niosło całego zdania modelu.
func polecenieRozpoznaniaJezyka(tekst string) string {
	return "Rozpoznaj język poniższego tekstu. Odpowiedz wyłącznie nazwą języka " +
		"(albo kodem ISO 639-1), bez zdania, bez komentarza.\n\nTekst:\n" + tekst
}
