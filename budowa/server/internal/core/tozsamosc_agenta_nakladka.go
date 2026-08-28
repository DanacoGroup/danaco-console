package core

import (
	"strings"

	"danacoconsole/server/internal/injection"
	"danacoconsole/server/internal/models"
	"danacoconsole/shared"
)

// Scalenie warstw eksperta z nakładką osi idzie jednym składaczem, jedną drogą do procesu.

// nakladkaZAgentem dopisuje tożsamość eksperta do nakładki wyliczonej z osi. Ekspert dopisuje, nigdy nie zastępuje: warstwy eksperta dokładają się do paczki wbudowanej. Instrukcje systemowe eksperta lądują przed konstytucją eksperta, po treści osi.
func nakladkaZAgentem(nakladka models.Nakladka, tozsamosc TozsamoscAgenta) models.Nakladka {
	if tozsamosc.Pusta() || !wnosiTresc(tozsamosc) {
		return nakladka
	}
	// Warstwa zerowa przed warstwami — kolejność dopisywania jest kolejnością w prompcie.
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
	// Prompt wbudowany zostaje w mocy; ekspert dopisuje się do niego, więc przełącznik jest dopisujący.
	nakladka.Tryb = injection.TrybDopisz
	return nakladka
}

// wnosiTresc mówi, czy ekspert ma w ogóle co dopisać, licząc warstwę zerową tak samo jak trzy pozostałe. Ekspert bez treści nie rusza trybu silnika — obowiązuje wtedy rozstrzygnięcie osi, bo brak zdania nie jest zdaniem.
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

// celWarstwy wskazuje pole nakładki odpowiadające nazwie warstwy z bazy, domykając rozjazd nazw między tabelą agent_warstwa a strukturą models.Nakladka. Nazwa spoza trójki daje nil, czyli warstwę pominiętą.
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

// zlozWarstwe dopisuje treść eksperta po treści osi. Trybu tu nie ma: kolumna tryb w agent_warstwa i pole AgentLayer.mode nadal istnieją, ale na złożenie promptu nie mają wpływu, bo ekspert dopisuje się do paczki wbudowanej zawsze.
func zlozWarstwe(osi, eksperta string) string {
	if strings.TrimSpace(osi) == "" {
		return eksperta
	}
	return osi + spoinaWarstw + eksperta
}
