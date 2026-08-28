// Plik składa prompt systemowy z wybranych treści i rozstrzyga tryb podania nakładki. Składanie jest deterministyczne: ta sama konfiguracja daje bajtowo ten sam prompt. Tryb ZASTAP wygrywa z DOLACZ, jeżeli choć jedna kategoria żąda zastąpienia.
package core

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"danacoconsole/server/internal/injection"
	"danacoconsole/shared"
)

// spoinaWarstw łączy warstwy nakładki. Ta sama spoina co w silniku nakładki —
// prompt złożony tutaj i prompt złożony przez injection.Nakladka.Prompt() mają
// być tym samym ciągiem bajtów, co potwierdza test składania.
const spoinaWarstw = "\n\n"

// porzadekWarstwy zwraca krytyczność warstwy: konstytucja, profil, ekspertyza — kolejność wyliczenia kontraktu IdentityLayer wiążącego swoje wartości ze stałymi silnika nakładki. Warstwa spoza wyliczenia ląduje na końcu.
func porzadekWarstwy(warstwa shared.IdentityLayer) int {
	switch warstwa {
	case shared.IdentityLayerConstitution:
		return 1
	case shared.IdentityLayerProfile:
		return 2
	case shared.IdentityLayerExpertise:
		return 3
	default:
		return 4
	}
}

// warstwaSilnika przekłada warstwę kontraktu na nazwę warstwy silnika nakładki.
// Przekład jest jedynym miejscem styku obu nazewnictw.
func warstwaSilnika(warstwa shared.IdentityLayer) string {
	switch warstwa {
	case shared.IdentityLayerConstitution:
		return injection.WarstwaKonstytucja
	case shared.IdentityLayerProfile:
		return injection.WarstwaProfil
	case shared.IdentityLayerExpertise:
		return injection.WarstwaEkspertyza
	default:
		return string(warstwa)
	}
}

// zlozNakladke buduje odpowiedź kontraktu z katalogu kategorii i zebranych treści, ustalając tryb i odcisk promptu.
func zlozNakladke(kategorie []shared.IdentityCategory, dokumenty []shared.IdentityDocument,
	trybDomyslny shared.IdentityMode) shared.IdentityEffectiveGetResponse {

	wybrane := wybierzTresci(kategorie, dokumenty)
	warstwy := make([]shared.IdentityLayerContent, 0, len(wybrane))
	czesci := make([]string, 0, len(wybrane))
	wnoszace := make([]shared.IdentityDocument, 0, len(wybrane))
	braki := []string{}

	for _, wybor := range wybrane {
		if !wybor.Wnosi() {
			if wybor.Kategoria.Required {
				braki = append(braki, wybor.Kategoria.Id)
			}
			continue
		}
		tresc := wybor.Tresc()
		warstwy = append(warstwy, shared.IdentityLayerContent{
			Layer:      wybor.Kategoria.Layer,
			CategoryId: wybor.Kategoria.Id,
			Name:       wybor.Kategoria.Name,
			Content:    tresc,
			Order:      len(warstwy) + 1,
			Axis:       wybor.Dokument.Axis,
			AxisId:     wybor.Dokument.AxisId,
		})
		czesci = append(czesci, tresc)
		wnoszace = append(wnoszace, wybor.Dokument)
	}

	prompt := strings.Join(czesci, spoinaWarstw)
	return shared.IdentityEffectiveGetResponse{
		Mode:                       trybObowiazujacy(wnoszace, trybDomyslny),
		Layers:                     warstwy,
		Prompt:                     prompt,
		PromptHash:                 odciskPromptu(prompt),
		MissingRequiredCategoryIds: braki,
	}
}

// trybObowiazujacy rozstrzyga tryb podania całej nakładki. Bez treści
// obowiązuje tryb domyślny konfiguracji — Operator ma widzieć, co obowiązywałoby
// po wpisaniu treści.
func trybObowiazujacy(dokumenty []shared.IdentityDocument, trybDomyslny shared.IdentityMode) shared.IdentityMode {
	domyslny := trybLubDomyslny(trybDomyslny, shared.IdentityModeZASTAP)
	if len(dokumenty) == 0 {
		return domyslny
	}
	for _, dokument := range dokumenty {
		if trybLubDomyslny(dokument.Mode, domyslny) == shared.IdentityModeZASTAP {
			return shared.IdentityModeZASTAP
		}
	}
	return shared.IdentityModeDOLACZ
}

// trybLubDomyslny sprowadza tryb do wartości kontraktu; wartość nieznana albo
// pusta znaczy tryb domyślny, nigdy odmowę składania.
func trybLubDomyslny(tryb, domyslny shared.IdentityMode) shared.IdentityMode {
	switch tryb {
	case shared.IdentityModeZASTAP, shared.IdentityModeDOLACZ:
		return tryb
	default:
		return domyslny
	}
}

// odciskPromptu liczy skrót SHA-256 złożonego promptu. Służy wyłącznie
// diagnostyce prowenancji i rozpoznaniu zmiany nakładki: niczego nie dopuszcza
// i niczego nie blokuje.
func odciskPromptu(prompt string) string {
	if prompt == "" {
		return ""
	}
	suma := sha256.Sum256([]byte(prompt))
	return hex.EncodeToString(suma[:])
}

// NakladkaSilnika przekłada nakładkę obowiązującą na nakładkę silnika wraz
// z trybem. Silnik zostaje nietknięty, a rozstrzygnięcie „zastąpić czy dopisać"
// przychodzi z danych katalogu, nie z kodu kanału.
func NakladkaSilnika(wynik shared.IdentityEffectiveGetResponse) injection.Nakladka {
	warstwy := make([]injection.Warstwa, 0, len(wynik.Layers))
	for _, warstwa := range wynik.Layers {
		warstwy = append(warstwy, injection.Warstwa{
			Nazwa: warstwaSilnika(warstwa.Layer),
			Tresc: warstwa.Content,
		})
	}
	return injection.NowaNakladka(trybSilnika(wynik.Mode), warstwy...)
}

// trybSilnika przekłada tryb kontraktu na tryb silnika nakładki: ZASTAP na
// TrybZastap (`--system-prompt`), DOLACZ na TrybDopisz
// (`--append-system-prompt`).
func trybSilnika(tryb shared.IdentityMode) string {
	if trybLubDomyslny(tryb, shared.IdentityModeZASTAP) == shared.IdentityModeDOLACZ {
		return injection.TrybDopisz
	}
	return injection.TrybZastap
}
