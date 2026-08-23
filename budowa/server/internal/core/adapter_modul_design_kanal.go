// Odpowiedzialność pliku: wybór kanału obrazowego dla `design.asset.generate` —
// jedno miejsce, w którym rozstrzyga się, którym kanałem pójdzie generowanie,
// i jedno miejsce, w którym nazywa się każdy powód, dla którego nie pójdzie
// żadnym. Metoda stoi na `*adapterDesignu` (`adapter_modul_design.go`).
//
// Kanał obrazowy to nie dowolny kanał modelu. Rejestr rdzenia prowadzi
// w większości do modeli tekstowych oddających fragmenty `text`; kanał obrazowy
// poznaje się po kluczu adaptera „obrazy" (`models.AdapterObrazy`, wiersz
// rodzaju `api` z parametrem `adapter`). Rozróżnienie pada tutaj, przed
// wywołaniem — inaczej moduł wysłałby prompt kanałowi tekstowemu i mówił
// Operatorowi „kanał nie oddał obrazu" zamiast wprost, że wskazał zły kanał.
//
// Odmowy są rozdzielone, bo Operator naprawia je różnymi ruchami: dopisać kanał
// obrazowy, poprawić identyfikator, włączyć kanał, wskazać kanał obrazowy
// zamiast tekstowego, naprawić parametry wiersza. Żadna z nich nie schodzi po
// cichu na kanał inny niż wskazany — obraz wygenerowany innym silnikiem byłby
// obrazem cudzym, a Operator nie miałby jak się o tym dowiedzieć.
//
// Braku poświadczenia ten plik nie rozstrzyga. Wiersz bez `credentialRef` buduje
// się celowo (`models/adapter_obrazy.go`), a odmowa pada w chwili wywołania,
// zdaniem samego kanału: rozróżnia ono „bez odwołania" od „odwołanie bez
// wartości w sejfie", czego ten moduł nie widzi. Powielenie tego sprawdzenia
// tutaj dałoby dwie prawdy o poświadczeniach i uboższą treść odmowy.
package core

import (
	"strings"

	"danacoconsole/server/internal/models"
	"danacoconsole/shared"
)

// kanalObrazowyZadania rozstrzyga, którym kanałem pójdzie generowanie, i oddaje
// jego wiersz (nie sam identyfikator), bo wołający wypisuje identyfikator kanału
// w treściach odmów.
//
// Brak wskazania bierze pierwszy czynny kanał obrazowy — tak opisuje pole
// `channelId` kontrakt. „Pierwszy" znaczy tu: pierwszy w kolejności wykazu rejestru,
// czynny i mający zbudowany adapter; wiersz włączony bez adaptera pomijamy
// milcząco przy szukaniu domyślnego, bo Operator o niego nie prosił i odmówiłby
// przy pierwszej turze.
func (a *adapterDesignu) kanalObrazowyZadania(wskazany *string,
	p shared.DesignPrompt, wariantow int) (models.Definicja, error) {

	if a.kanaly == nil {
		return models.Definicja{}, bladOdmowyGenerowania(shared.ErrorCodeChannelUnavailable,
			"rejestr kanałów modelu nie jest wpięty — nie ma czym wołać silnika obrazów",
			p, wariantow, "wpięty rejestr kanałów modelu")
	}

	kod := ""
	if wskazany != nil {
		kod = strings.TrimSpace(*wskazany)
	}
	if kod == "" {
		return a.pierwszyKanalObrazowy(p, wariantow)
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
		// Rodzaj kanału sprawdza się przed czynnością: kanał tekstowy włączony
		// jest tą samą pomyłką, co kanał tekstowy wyłączony — włączanie go
		// niczego nie naprawi, więc zdanie o czynności wysłałoby Operatora
		// w złą stronę.
		if wiersz.KluczAdaptera() != models.AdapterObrazy {
			return models.Definicja{}, bladOdmowyGenerowania(shared.ErrorCodeValidationFailed,
				"wskazany kanał "+kod+" nie jest kanałem obrazowym (adapter „"+wiersz.KluczAdaptera()+
					"”) — oddaje fragmenty TEKSTU, a fragment tekstowy nie jest obrazem "+
					"i rdzeń nie zrobi z niego zasobu wizualnego; naprawa: wskazać kanał założony "+
					"z config.adapter=\"obrazy\" (channel.list pokazuje konfigurację wierszy)",
				p, wariantow, "kanał obrazowy zamiast wskazanego kanału tekstowego")
		}
		if !wiersz.Aktywny {
			return models.Definicja{}, bladOdmowyGenerowania(shared.ErrorCodeChannelUnavailable,
				"wskazany kanał obrazowy "+kod+" jest nieczynny — generowanie nie zejdzie po cichu "+
					"na inny kanał; naprawa: włączyć kanał (channel.update enabled=true) albo wskazać inny",
				p, wariantow, "czynny kanał obrazowy")
		}
		if _, gotowy := a.kanaly.Kanal(wiersz.Identyfikator()); !gotowy {
			return models.Definicja{}, bladOdmowyGenerowania(shared.ErrorCodeChannelUnavailable,
				"wskazany kanał obrazowy "+kod+" jest włączony, ale nie ma zbudowanego adaptera — "+
					"nie ma czym wołać silnika obrazów; naprawa: sprawdzić parametry kanału, "+
					"zwłaszcza base_url (channel.list pokazuje wiersze pominięte)",
				p, wariantow, "zbudowany adapter wskazanego kanału")
		}
		return wiersz, nil
	}
	return models.Definicja{}, bladOdmowyGenerowania(shared.ErrorCodeNotFound,
		"nie ma kanału o identyfikatorze "+kod+" — generowanie nie zejdzie po cichu na kanał inny "+
			"niż wskazany; naprawa: sprawdzić wykaz kanałów (channel.list) albo pominąć pole channelId, "+
			"żeby jechać pierwszym czynnym kanałem obrazowym",
		p, wariantow, "kanał o wskazanym identyfikatorze")
}

// pierwszyKanalObrazowy szuka domyślnego kanału obrazowego. Brak takiego kanału
// jest brakiem konfiguracji, nie brakiem produktu: adapter obrazowy w rdzeniu
// jest, magazyn jest, pole `channelId` w kontrakcie jest. Dlatego zdanie odmowy
// mówi, jak taki kanał założyć, zamiast opisywać granicę rdzenia.
func (a *adapterDesignu) pierwszyKanalObrazowy(p shared.DesignPrompt,
	wariantow int) (models.Definicja, error) {

	for _, wiersz := range a.kanaly.Wykaz() {
		if !wiersz.Aktywny || wiersz.KluczAdaptera() != models.AdapterObrazy {
			continue
		}
		if _, gotowy := a.kanaly.Kanal(wiersz.Identyfikator()); !gotowy {
			continue
		}
		return wiersz, nil
	}
	return models.Definicja{}, bladOdmowyGenerowania(shared.ErrorCodeChannelUnavailable,
		"w rejestrze nie ma czynnego kanału oddającego BAJTY OBRAZU, a rdzeń nie założy zasobu "+
			"wizualnego bez treści wizualnej i nie podstawi obrazu zastępczego; naprawa: założyć kanał "+
			"komendą channel.add z kind=\"api\" i config.adapter=\"obrazy\", wskazując base_url punktu "+
			"końcowego zgodnego z OpenAI Images oraz credentialRef z kluczem. Prompt strukturalny "+
			"został złożony w całości i jego gotowa treść wraca w `details.polecenie` — nic z pracy "+
			"Prompt Buildera nie ginie razem z tą odmową",
		p, wariantow, "czynny kanał modelu oddający bajty obrazu")
}
