package core

import (
	"strings"

	"danacoconsole/server/internal/injection"
	"danacoconsole/server/internal/models"
	"danacoconsole/shared"
)

// Scalenie warstw eksperta z nakładką osi — jeden składacz, nie drugi.
//
// To jest miejsce, w którym tożsamość eksperta wchodzi do promptu. Nie ma tu
// własnej drogi do procesu i nie może jej być: warstwy jadą dalej dokładnie tą
// samą trasą co dotąd — `nakladkaOkna` → `models.Nakladka` →
// `nakladkaKanaluGlownego` → `injection.Nakladka` → przełącznik CLI. Druga
// droga do promptu znaczyłaby dwie prawdy o tym, co model dostał, a przy
// sporze nie dałoby się ustalić, która zadziałała.
//
// Dowodem, że droga jest jedna, jest komenda `config.explain.get`: liczy
// prowenancję tymi samymi funkcjami, którymi jedzie tura. Gdyby ekspert wchodził
// obok, podgląd pokazywałby wiersz, którego tura nigdy nie wykona.

// nakladkaZAgentem dopisuje tożsamość eksperta do nakładki wyliczonej z osi.
//
// Ekspert dopisuje, nigdy nie zastępuje: paczka wbudowana (system prompt)
// obowiązuje zawsze, a warstwy eksperta dokładają się do niej jako zakres
// typu user.
//
// Skutek jest zupełny i bez wyjątku: treść osi zostaje nietknięta w każdej
// warstwie, tryb silnika idzie ku dopisaniu, a ekspert nie ma żadnej drogi,
// którą mógłby prompt globalny zdjąć.
//
// Warstwa zerowa — `agent.instrukcje_systemowe`.
// Instrukcje systemowe eksperta idą pierwsze spośród tego, co ekspert wnosi,
// i lądują w warstwie najbardziej krytycznej. „Zerowa" znaczy: przed
// konstytucją eksperta, ale po treści osi — bo przed treścią osi nie ma niczego,
// co ekspert mógłby postawić.
//
// Kolejność w warstwie konstytucji jest więc taka:
//
//	[oś: platforma · model · konto]  →  [ekspert: instrukcje systemowe]  →  [ekspert: konstytucja]
//
// Warstwa wyłączona (`aktywna = 0`) i warstwa o pustej treści są pomijane
// jednakowo: ekspert, który nie ma nic do powiedzenia w danej warstwie, po
// prostu nic w niej nie mówi.
func nakladkaZAgentem(nakladka models.Nakladka, tozsamosc TozsamoscAgenta) models.Nakladka {
	if tozsamosc.Pusta() || !wnosiTresc(tozsamosc) {
		return nakladka
	}
	// Warstwa zerowa przed warstwami — kolejność dopisywania jest kolejnością
	// w prompcie, bo `zlozWarstwe` dokleja na końcu.
	if instrukcje := strings.TrimSpace(tozsamosc.InstrukcjeSystemowe); instrukcje != "" {
		nakladka.Konstytucja = zlozWarstwe(nakladka.Konstytucja, instrukcje)
	}
	for _, warstwa := range tozsamosc.Warstwy {
		if !warstwa.Aktywna {
			continue
		}
		tresc := strings.TrimSpace(warstwa.Tresc)
		if tresc == "" {
			continue
		}
		cel := celWarstwy(&nakladka, warstwa.Warstwa)
		if cel == nil {
			continue
		}
		*cel = zlozWarstwe(*cel, tresc)
	}
	// Prompt wbudowany zostaje w mocy. Ekspert dopisuje się do niego, więc
	// przełącznik musi być dopisujący. To jedyne miejsce, w którym ekspert
	// dotyka trybu silnika, i dotyka go w jedną stronę.
	//
	// Oś ustawiona na ZASTAP traci przy ekspercie swoje żądanie zastąpienia,
	// ale nie traci ani zdania treści: jej warstwy jadą dalej, pierwsze, a przed
	// nimi zostaje prompt własny programu `claude`. Okno z ekspertem dostaje
	// o jedną paczkę więcej, nie o jedną mniej. Okno bez eksperta zachowuje
	// rozstrzygnięcie osi bez zmiany.
	nakladka.Tryb = injection.TrybDopisz
	return nakladka
}

// wnosiTresc mówi, czy ekspert ma w ogóle co dopisać.
//
// Liczy się warstwa zerowa tak samo jak trzy pozostałe: ekspert z samymi
// instrukcjami systemowymi, bez ani jednej warstwy, wnosi treść i przestawia
// tryb silnika. Pominięcie go tutaj znaczyłoby, że jego instrukcje dojeżdżają
// do promptu, a paczka wbudowana mimo to znika.
//
// Ekspert bez treści nie rusza trybu silnika — obowiązuje wtedy rozstrzygnięcie
// osi. Brak zdania nie jest zdaniem.
func wnosiTresc(tozsamosc TozsamoscAgenta) bool {
	if strings.TrimSpace(tozsamosc.InstrukcjeSystemowe) != "" {
		return true
	}
	for _, warstwa := range tozsamosc.Warstwy {
		if warstwa.Aktywna && strings.TrimSpace(warstwa.Tresc) != "" {
			return true
		}
	}
	return false
}

// celWarstwy wskazuje pole nakładki odpowiadające nazwie warstwy z bazy.
//
// Rozjazd nazw jest tu domykany, nie obchodzony. Tabela `agent_warstwa`
// przechowuje wartości kontraktu (`constitution` · `profile` · `expertise`),
// a `models.Nakladka` nazywa te same warstwy po polsku. Przekład stoi w jednym
// miejscu; nazwa spoza trójki daje `nil`, czyli warstwę pominiętą — nowy rodzaj
// warstwy nie ma prawa po cichu wylądować w cudzym polu.
func celWarstwy(n *models.Nakladka, nazwa string) *string {
	switch shared.IdentityLayer(nazwa) {
	case shared.IdentityLayerConstitution:
		return &n.Konstytucja
	case shared.IdentityLayerProfile:
		return &n.ProfilRoli
	case shared.IdentityLayerExpertise:
		return &n.Ekspertyza
	}
	return nil
}

// zlozWarstwe dopisuje treść eksperta po treści osi.
//
// Trybu tu nie ma. Kolumna `tryb` w `agent_warstwa` (migracja 073) i pole
// `AgentLayer.mode` w kontrakcie nadal istnieją i nadal są zapisywane przez
// `agent.layer.set` — ale na złożenie promptu nie mają wpływu, bo ekspert
// dopisuje się do paczki wbudowanej zawsze.
// Honorowanie tu `ZASTAP` znaczyłoby, że Operator jednym polem formularza kasuje
// prompt systemowy platformy, a ten ma obowiązywać zawsze.
//
// Rozbieżność między tym, co formularz pozwala zapisać, a tym, co zmienia
// wynik, jest znana: zdjęcie kolumny albo nadanie jej innego znaczenia to
// osobna zmiana. Kolumnę zakłada `migracja_073_agent_warstwy.sql`.
func zlozWarstwe(osi, eksperta string) string {
	if strings.TrimSpace(osi) == "" {
		return eksperta
	}
	return osi + spoinaWarstw + eksperta
}
