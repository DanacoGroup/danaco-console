// Egzekucja zakresu `dostep_sieciowy` — połączeń wychodzących zestawianych
// i konfigurowanych przez rdzeń w imieniu okna.
//
// Zakres włączony znaczy: wolno sięgać wyłącznie tam, gdzie Operator dał oknu
// jawne nadanie. Wykaz maszyn bierze się więc z nadań dostępu okna, nie
// z osobnego ustawienia — drugiego miejsca, w którym Operator wskazywałby
// maszyny, w produkcie nie ma.
//
// Granica egzekucji: straż rozstrzyga adresy, po które sięga rdzeń, oraz wykaz
// mostów podawany procesowi modelu w konfiguracji MCP. Gniazda otwierane przez
// sam proces modelu pozostają poza jej zasięgiem — na to potrzeba zapory albo
// przestrzeni nazw sieci, czyli środka systemu, nie rdzenia.
package core

import (
	"strings"

	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// strazSieci rozstrzyga, czy adres mieści się w obszarze sieciowym okna.
type strazSieci struct {
	zasady session.Zasady
	hosty  map[string]struct{}
}

// nowaStrazSieci składa straż z zasad izolacji i kompletu nadań okna. Zakres
// wyłączony daje straż przepuszczającą wszystko — wspólny dostęp sieciowy
// serwera jest stanem wyjściowym platformy.
func nowaStrazSieci(zasady session.Zasady, nadania []NadanieMostu) *strazSieci {
	return &strazSieci{zasady: zasady, hosty: hostyNadan(nadania)}
}

// hostyNadan wylicza maszyny jawnie nadane oknu. Odsiew nadań nieczynnych
// i punktów niebędących maszyną robi uporzadkowaneNadania, a nazwę maszyny
// wyprowadza AdresMostu — ten sam kod, który składa wpis mostu w konfiguracji
// MCP.
func hostyNadan(nadania []NadanieMostu) map[string]struct{} {
	hosty := make(map[string]struct{}, len(nadania))
	for _, nadanie := range uporzadkowaneNadania(nadania) {
		if host := strings.ToLower(strings.TrimSpace(AdresMostu(nadanie.Punkt).Host)); host != "" {
			hosty[host] = struct{}{}
		}
	}
	return hosty
}

// sprawdzAdres odrzuca połączenie wychodzące do maszyny, której Operator nie
// nadał oknu. Adres bez nazwy maszyny również odpada: przy izolacji włączonej
// nie ma czego porównać z wykazem nadań.
func (s *strazSieci) sprawdzAdres(adres string) error {
	if s == nil || !s.zasady.DostepSieciowy {
		return nil
	}
	host := hostAdresu(adres)
	if host == "" {
		return session.NoweNaruszenie(konfig.KluczIzolacjaDostepSieciowy,
			"adres "+strings.TrimSpace(adres)+" nie wskazuje maszyny, więc nie mieści się w nadaniach okna")
	}
	if _, jest := s.hosty[host]; !jest {
		return session.NoweNaruszenie(konfig.KluczIzolacjaDostepSieciowy,
			"maszyna "+host+" nie jest nadana temu oknu")
	}
	return nil
}

// mosty przepuszcza wykaz punktów mostowych podawany procesowi modelu. Most
// spoza nadań okna zatrzymuje złożenie konfiguracji — wykaz obcięty po cichu
// dawałby Operatorowi obraz niezgodny z tym, co dostał proces.
func (s *strazSieci) mosty(punkty []shared.AccessPoint) ([]shared.AccessPoint, error) {
	if s == nil || !s.zasady.DostepSieciowy {
		return punkty, nil
	}
	for _, punkt := range punkty {
		if punkt.Kind != shared.AccessPointKindMcpBridge {
			continue
		}
		if err := s.sprawdzAdres(AdresMostu(punkt).Host); err != nil {
			return nil, err
		}
	}
	return punkty, nil
}

// hostAdresu wyłuskuje nazwę maszyny z adresu połączenia wychodzącego: adresu
// URL, zapisu host:port albo samej nazwy. Zapis IPv6 w nawiasach kwadratowych
// wraca bez nawiasów, żeby porównanie z wykazem nadań szło po tej samej postaci.
func hostAdresu(adres string) string {
	tekst := strings.TrimSpace(adres)
	if podzial := strings.Index(tekst, "://"); podzial >= 0 {
		tekst = tekst[podzial+3:]
	}
	if podzial := strings.IndexAny(tekst, "/?#"); podzial >= 0 {
		tekst = tekst[:podzial]
	}
	if podzial := strings.LastIndex(tekst, "@"); podzial >= 0 {
		tekst = tekst[podzial+1:]
	}
	if strings.HasPrefix(tekst, "[") {
		if koniec := strings.Index(tekst, "]"); koniec > 0 {
			return strings.ToLower(tekst[1:koniec])
		}
	}
	if podzial := strings.LastIndex(tekst, ":"); podzial >= 0 && portPoprawny(tekst[podzial+1:]) {
		tekst = tekst[:podzial]
	}
	return strings.ToLower(strings.TrimSpace(tekst))
}
