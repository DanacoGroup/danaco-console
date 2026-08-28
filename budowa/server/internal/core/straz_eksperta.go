// Plik zawiera miejsca, w których rdzeń czyta zakres eksperta i odmawia: nałożenie eksperta na okno modułu spoza jego zakresu, oraz powołanie podagentów przez eksperta z wyłączonym Subagent Network. Straż milczy, gdy zakres nie jest zawężony.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// StrazEksperta rozstrzyga, czy ekspert może działać w danym miejscu. Port stoi po stronie odbiorcy: pytają go adapter okien przy nałożeniu eksperta i adapter podagentów przy powołaniu. Straż niewpięta nie zmienia niczego.
type StrazEksperta interface {
	// SprawdzModulOkna odmawia, gdy zakres eksperta nie obejmuje modułu tego
	// okna. Nil znaczy „wolno".
	SprawdzModulOkna(ctx context.Context, kodEksperta string, modulID int64) error
	// GranicaPodagentowEksperta oddaje górną liczbę jednoczesnych podagentów eksperta i czy je uruchamia.
	GranicaPodagentowEksperta(ctx context.Context, kodEksperta string) (int, bool)
}

// strazEksperta wypełnia port trzema repozytoriami: biblioteką (uprawnienia
// i granica podagentów), zakresem (moduły zastosowania) i katalogiem modułów
// (przekład numeru wiersza okna na kod modułu).
type strazEksperta struct {
	biblioteka dane.RepozytoriumAgentow
	moduly     dane.RepozytoriumModulow
}

var _ StrazEksperta = (*strazEksperta)(nil)

// NowaStrazEksperta wiąże straż z biblioteką ekspertów i katalogiem modułów, tworząc gotowy do użycia port StrazEksperta.
func NowaStrazEksperta(biblioteka dane.RepozytoriumAgentow,
	moduly dane.RepozytoriumModulow) *strazEksperta {

	return &strazEksperta{biblioteka: biblioteka, moduly: moduly}
}

// SprawdzModulOkna odmawia nałożenia eksperta na okno modułu spoza jego zakresu, rozstrzyganego wykazem modułów zastosowania oraz wpisem uprawnienia grupy modules. Ekspert nierozpoznany przechodzi bez odmowy jako brak, nie zawężenie.
func (s *strazEksperta) SprawdzModulOkna(ctx context.Context, kodEksperta string, modulID int64) error {
	if s == nil || s.biblioteka == nil || strings.TrimSpace(kodEksperta) == "" || modulID == 0 {
		return nil
	}
	ekspert, err := s.biblioteka.PoKodzie(ctx, kodEksperta)
	if err != nil {
		return nil
	}
	kodModulu, znany := s.kodModuluOkna(ctx, modulID)
	if !znany {
		return nil
	}
	if len(ekspert.ModulyZastosowania) > 0 && !zawieraKodModulu(ekspert.ModulyZastosowania, kodModulu) {
		return odmowaZakresuEksperta(kodEksperta, kodModulu,
			"moduł nie należy do modułów zastosowania zapisanych w definicji eksperta")
	}
	for _, uprawnienie := range ekspert.Uprawnienia {
		if uprawnienie.Grupa != shared.AgentPermissionGroupModules {
			continue
		}
		if strings.TrimSpace(uprawnienie.Zakres) != kodModulu {
			continue
		}
		if !uprawnienie.Przyznane {
			return odmowaZakresuEksperta(kodEksperta, kodModulu,
				"Permissions Center odebrał temu ekspertowi dostęp do tego modułu")
		}
	}
	return nil
}

// kodModuluOkna przekłada numer wiersza modułu okna na jego kod. Katalog
// nieodczytany daje „nie wiadomo", a nie odmowę: zawężenie ma wynikać z decyzji
// Operatora, nie z nieudanego odczytu.
func (s *strazEksperta) kodModuluOkna(ctx context.Context, modulID int64) (string, bool) {
	if s.moduly == nil {
		return "", false
	}
	wiersze, err := s.moduly.Lista(ctx)
	if err != nil {
		return "", false
	}
	for _, wiersz := range wiersze {
		if wiersz.ID == modulID {
			return wiersz.Kod, true
		}
	}
	return "", false
}

// GranicaPodagentowEksperta oddaje granicę Subagent Network eksperta.
//
// Druga wartość mówi, czy odpowiedź w ogóle dotyczy eksperta: ekspert
// nierozpoznany i straż bez biblioteki dają `false`, a wołający zostaje wtedy
// przy granicy platformy.
func (s *strazEksperta) GranicaPodagentowEksperta(ctx context.Context, kodEksperta string) (int, bool) {
	if s == nil || s.biblioteka == nil || strings.TrimSpace(kodEksperta) == "" {
		return 0, false
	}
	ekspert, err := s.biblioteka.PoKodzie(ctx, kodEksperta)
	if err != nil {
		return 0, false
	}
	return ekspert.LimitPodagentow, true
}

// zawieraKodModulu sprawdza obecność kodu modułu w wykazie modułów zastosowania eksperta, przez proste porównanie.
func zawieraKodModulu(kody []string, szukany string) bool {
	for _, kod := range kody {
		if kod == szukany {
			return true
		}
	}
	return false
}

// odmowaZakresuEksperta składa odmowę nałożenia eksperta wraz z powodem
// i wskazaniem, gdzie Operator zawężenie zdejmie.
func odmowaZakresuEksperta(kodEksperta, kodModulu, powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodePermissionDenied,
		"zakres eksperta: "+kodEksperta+" nie pracuje w module "+kodModulu+" — "+powod+
			"; Operator poszerzy zakres w Permissions Center modułu Agents "+
			"(agent.modules.set albo agent.permission.remove) albo wskaże innego eksperta"))
}

// odmowaPodagentowEksperta składa odmowę powołania podagentów przez eksperta z wyłączonym Subagent Network.
func odmowaPodagentowEksperta(kodEksperta string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodePermissionDenied,
		"zakres eksperta: "+kodEksperta+" ma wyłączony Subagent Network — powołanie podagentów "+
			"nie idzie dalej; Operator włączy go w Permissions Center modułu Agents "+
			"(agent.subagent.set z enabled: true)"))
}
