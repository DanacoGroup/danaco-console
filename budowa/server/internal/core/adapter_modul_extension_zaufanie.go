// Rodzina `extension.*` — Permissions & Trust Center: uprawnienia deklarowane
// i nadane, weryfikacja podpisu, skaner manifestu, rejestr referencji
// sekretów. Ostrzeżenie ani brak podpisu nie blokują instalacji: kontrola
// idzie przez stan i zakres.
package core

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"sort"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// Kody spostrzeżeń skanera manifestu, wystawiane jako sygnał w Permissions
// & Trust Center, nigdy jako blokada.
const (
	kodSkaneraBezPodpisu     = "extensionWithoutSignature"
	kodSkaneraPodpisNiepewny = "extensionSignatureUnverified"
	kodSkaneraSzerokiZakres  = "permissionWithoutTarget"
	kodSkaneraNadmiarowe     = "permissionGrantedButNotDeclared"
	kodSkaneraBezWydawcy     = "extensionWithoutPublisher"
	kodSkaneraProcesy        = "permissionProcessSpawn"
)

// WypiszUprawnienia obsługuje `extension.permission.list`: oddaje zakresy
// nadane wraz z uprawnieniami nadmiarowymi.
func (a *adapterRozszerzen) WypiszUprawnienia(ctx context.Context,
	z shared.ExtensionPermissionListRequest) (shared.ExtensionPermissionListResponse, error) {

	wiersz, err := a.pozycjaRozszerzeniaZadania(ctx, z.ExtensionId)
	if err != nil {
		return shared.ExtensionPermissionListResponse{}, err
	}
	wiersze, err := a.rejestr.UprawnieniaRozszerzenia(ctx, wiersz.Identyfikator)
	if err != nil {
		return shared.ExtensionPermissionListResponse{}, bladRozszerzenia(err)
	}

	deklarowane := []shared.ExtensionPermission{}
	nadane := []shared.ExtensionPermission{}
	zakresyDeklarowane := map[string]struct{}{}
	zakresyNadane := map[string]struct{}{}
	for _, uprawnienie := range wiersze {
		if uprawnienie.Nadane {
			nadane = append(nadane, uprawnienieKontraktuRozszerzen(uprawnienie))
			zakresyNadane[uprawnienie.Zakres] = struct{}{}
			continue
		}
		deklarowane = append(deklarowane, uprawnienieKontraktuRozszerzen(uprawnienie))
		zakresyDeklarowane[uprawnienie.Zakres] = struct{}{}
	}

	nadmiarowe := []string{}
	for zakres := range zakresyNadane {
		if _, jest := zakresyDeklarowane[zakres]; jest {
			continue
		}
		nadmiarowe = append(nadmiarowe, zakres)
	}
	sort.Strings(nadmiarowe)

	return shared.ExtensionPermissionListResponse{
		Declared: deklarowane, Granted: nadane, Excessive: nadmiarowe,
	}, nil
}

// NadajUprawnienia obsługuje `extension.permission.grant`. Nadanie wymienia
// komplet uprawnień nadanych, nie dokłada do nich: żądanie z krótszym
// wykazem cofa te, których w nim nie ma.
func (a *adapterRozszerzen) NadajUprawnienia(ctx context.Context,
	z shared.ExtensionPermissionGrantRequest) (shared.ExtensionPermissionGrantResponse, error) {

	wiersz, err := a.pozycjaRozszerzeniaZadania(ctx, z.ExtensionId)
	if err != nil {
		return shared.ExtensionPermissionGrantResponse{}, err
	}
	if err := a.sprawdzEksperta(ctx, wartoscTekstu(z.AgentId)); err != nil {
		return shared.ExtensionPermissionGrantResponse{}, err
	}
	teraz := a.teraz()
	wiersze := make([]dane.UprawnienieRozszerzenia, 0, len(z.Permissions))
	for _, uprawnienie := range z.Permissions {
		if err := sprawdzZakresUprawnieniaRozszerzenia(uprawnienie.Scope); err != nil {
			return shared.ExtensionPermissionGrantResponse{}, err
		}
		wpis := dane.UprawnienieRozszerzenia{
			Zakres: string(uprawnienie.Scope), Byt: uprawnienie.Target,
			Objasnienie: uprawnienie.Explanation, AgentKod: z.AgentId, Nadano: &teraz,
		}
		if uprawnienie.Mode != nil {
			tryb := string(*uprawnienie.Mode)
			wpis.Tryb = &tryb
		}
		wiersze = append(wiersze, wpis)
	}
	if err := a.rejestr.ZapiszUprawnieniaRozszerzenia(ctx, wiersz.Identyfikator, true, wiersze); err != nil {
		return shared.ExtensionPermissionGrantResponse{}, bladRozszerzenia(err)
	}

	zapisane, err := a.rejestr.UprawnieniaRozszerzenia(ctx, wiersz.Identyfikator)
	if err != nil {
		return shared.ExtensionPermissionGrantResponse{}, bladRozszerzenia(err)
	}
	nadane := []shared.ExtensionPermission{}
	for _, uprawnienie := range zapisane {
		if !uprawnienie.Nadane {
			continue
		}
		nadane = append(nadane, uprawnienieKontraktuRozszerzen(uprawnienie))
	}
	a.odnotujCyklZycia(ctx, wiersz, shared.ExtensionLifecycleActionConfigured, nil, nil,
		"nadano zakres uprawnień: "+strconv.Itoa(len(nadane)))
	return shared.ExtensionPermissionGrantResponse{Granted: nadane}, nil
}

// ZweryfikujPodpis obsługuje `extension.signature.verify`: liczy podpis od
// nowa kluczem publicznym pozycji.
func (a *adapterRozszerzen) ZweryfikujPodpis(ctx context.Context,
	z shared.ExtensionSignatureVerifyRequest) (shared.ExtensionSignatureVerifyResponse, error) {

	wiersz, err := a.pozycjaRozszerzeniaZadania(ctx, z.ExtensionId)
	if err != nil {
		return shared.ExtensionSignatureVerifyResponse{}, err
	}
	teraz := a.teraz()

	podpis, err := a.rejestr.PodpisRozszerzenia(ctx, wiersz.Identyfikator)
	if err != nil {
		if !isBrakWierszaApp(err) {
			return shared.ExtensionSignatureVerifyResponse{}, bladRozszerzenia(err)
		}
		// Pozycja bez podpisu nie jest odmową: brak podpisu jest faktem do
		// pokazania, nie usterką odczytu.
		sygnatura := shared.ExtensionSignature{
			Signed: false, Verified: false,
			TrustLevel: poziomZaufaniaBezPodpisu(wiersz),
			Detail: wskaznikNapisuApp(
				"pozycja nie ma zapisanego podpisu — pochodzenie potwierdza wyłącznie źródło"),
		}
		return shared.ExtensionSignatureVerifyResponse{Signature: sygnatura, CheckedAt: teraz}, nil
	}

	sygnatura := sygnaturaKontraktu(podpis, wiersz)
	sygnatura.Verified = false

	switch {
	case podpis.PodpisBase64 == nil || podpis.KluczBase64 == nil ||
		podpis.SumaKontrolna == nil:
		sygnatura.Detail = wskaznikNapisuApp(
			"podpis nie niesie kompletu: potrzebne są bajty podpisu, klucz publiczny i suma pakietu")
	default:
		bajtyPodpisu, bladPodpisu := base64.StdEncoding.DecodeString(*podpis.PodpisBase64)
		bajtyKlucza, bladKlucza := base64.StdEncoding.DecodeString(*podpis.KluczBase64)
		suma, bladSumy := hex.DecodeString(*podpis.SumaKontrolna)
		switch {
		case bladPodpisu != nil || bladKlucza != nil || bladSumy != nil:
			sygnatura.Detail = wskaznikNapisuApp("materiał podpisu jest nieczytelny")
		case len(bajtyKlucza) != ed25519.PublicKeySize:
			sygnatura.Detail = wskaznikNapisuApp("klucz publiczny ma nieprawidłową długość")
		default:
			sygnatura.Verified = ed25519.Verify(ed25519.PublicKey(bajtyKlucza), suma, bajtyPodpisu)
			if !sygnatura.Verified {
				sygnatura.Detail = wskaznikNapisuApp(
					"podpis nie zgadza się z sumą kontrolną pakietu ani z kluczem wydawcy")
			}
		}
	}

	// Poziom zaufania wynika z weryfikacji: bez potwierdzenia podpisem
	// wydawca jest niezweryfikowany.
	if !sygnatura.Verified &&
		sygnatura.TrustLevel == shared.ExtensionTrustLevelVerifiedPublisher {
		sygnatura.TrustLevel = shared.ExtensionTrustLevel(shared.ExtensionTrustLevelUnverifiedPersonal)
	}
	if sygnatura.Verified {
		sygnatura.TrustLevel = shared.ExtensionTrustLevel(shared.ExtensionTrustLevelVerifiedPublisher)
		if wiersz.ZrodloPochodzenia == shared.ExtensionOriginDanaco {
			sygnatura.TrustLevel = shared.ExtensionTrustLevel(shared.ExtensionTrustLevelDanacoPlugin)
		}
	}
	return shared.ExtensionSignatureVerifyResponse{Signature: sygnatura, CheckedAt: teraz}, nil
}

// poziomZaufaniaBezPodpisu nazywa zaufanie pozycji, której nikt nie podpisał
// — nie ma czym go potwierdzić.
func poziomZaufaniaBezPodpisu(wiersz dane.Rozszerzenie) shared.ExtensionTrustLevel {
	if wiersz.ZrodloPochodzenia == shared.ExtensionOriginDanaco {
		return shared.ExtensionTrustLevel(shared.ExtensionTrustLevelDanacoPlugin)
	}
	return shared.ExtensionTrustLevel(shared.ExtensionTrustLevelUnverifiedPersonal)
}

// SkanujManifest obsługuje `extension.manifest.scan`: wystawia spostrzeżenia,
// nigdy blokadę. Spostrzeżenia liczą się z tego, co pozycja naprawdę
// deklaruje i co ma nadane.
func (a *adapterRozszerzen) SkanujManifest(ctx context.Context,
	z shared.ExtensionManifestScanRequest) (shared.ExtensionManifestScanResponse, error) {

	wiersz, err := a.pozycjaRozszerzeniaZadania(ctx, z.ExtensionId)
	if err != nil {
		return shared.ExtensionManifestScanResponse{}, err
	}
	uprawnienia, err := a.rejestr.UprawnieniaRozszerzenia(ctx, wiersz.Identyfikator)
	if err != nil {
		return shared.ExtensionManifestScanResponse{}, bladRozszerzenia(err)
	}

	spostrzezenia := []shared.ExtensionScanFinding{}
	zakresyDeklarowane := map[string]struct{}{}
	for _, uprawnienie := range uprawnienia {
		if !uprawnienie.Nadane {
			zakresyDeklarowane[uprawnienie.Zakres] = struct{}{}
		}
	}

	for _, uprawnienie := range uprawnienia {
		zakres := shared.ExtensionPermissionScope(uprawnienie.Zakres)
		// Uprawnienie bez wskazania bytu obejmuje wszystko: sieć bez domeny,
		// zapis bez korzenia katalogu.
		if uprawnienie.Byt == nil || strings.TrimSpace(*uprawnienie.Byt) == "" {
			spostrzezenia = append(spostrzezenia, shared.ExtensionScanFinding{
				Severity:        shared.AppValidationSeverityWarning,
				Code:            kodSkaneraSzerokiZakres,
				Message:         "uprawnienie " + uprawnienie.Zakres + " nie wskazuje bytu, którego dotyczy",
				PermissionScope: &zakres,
			})
		}
		// Uruchamianie procesów jest zakresem najszerszym: pozycja z nim
		// zrobi wszystko na maszynie Operatora.
		if uprawnienie.Zakres == shared.ExtensionPermissionScopeProcessSpawn {
			spostrzezenia = append(spostrzezenia, shared.ExtensionScanFinding{
				Severity:        shared.AppValidationSeverityWarning,
				Code:            kodSkaneraProcesy,
				Message:         "pozycja żąda uruchamiania procesów — to zakres najszerszy z możliwych",
				PermissionScope: &zakres,
			})
		}
		if !uprawnienie.Nadane {
			continue
		}
		if _, jest := zakresyDeklarowane[uprawnienie.Zakres]; jest {
			continue
		}
		spostrzezenia = append(spostrzezenia, shared.ExtensionScanFinding{
			Severity:        shared.AppValidationSeverityError,
			Code:            kodSkaneraNadmiarowe,
			Message:         "zakres " + uprawnienie.Zakres + " jest nadany, choć manifest go nie deklaruje",
			PermissionScope: &zakres,
		})
	}

	podpis, err := a.rejestr.PodpisRozszerzenia(ctx, wiersz.Identyfikator)
	switch {
	case err != nil && isBrakWierszaApp(err):
		spostrzezenia = append(spostrzezenia, shared.ExtensionScanFinding{
			Severity: shared.AppValidationSeverityWarning,
			Code:     kodSkaneraBezPodpisu,
			Message:  "pozycja nie ma podpisu — pochodzenia nie da się potwierdzić kryptograficznie",
		})
	case err != nil:
		return shared.ExtensionManifestScanResponse{}, bladRozszerzenia(err)
	default:
		if podpis.Wydawca == nil || strings.TrimSpace(*podpis.Wydawca) == "" {
			spostrzezenia = append(spostrzezenia, shared.ExtensionScanFinding{
				Severity: shared.AppValidationSeverityInfo,
				Code:     kodSkaneraBezWydawcy,
				Message:  "podpis nie wskazuje wydawcy",
			})
		}
		if podpis.PodpisBase64 == nil || podpis.KluczBase64 == nil {
			spostrzezenia = append(spostrzezenia, shared.ExtensionScanFinding{
				Severity: shared.AppValidationSeverityWarning,
				Code:     kodSkaneraPodpisNiepewny,
				Message:  "podpis nie niesie materiału do weryfikacji — nie da się go sprawdzić",
			})
		}
	}

	return shared.ExtensionManifestScanResponse{
		Findings: spostrzezenia, ScannedAt: a.teraz(),
	}, nil
}

// WypiszSekrety obsługuje `extension.secret.list` — rejestr referencji wraz
// z przypomnieniem o wygasających poświadczeniach.
func (a *adapterRozszerzen) WypiszSekrety(ctx context.Context,
	z shared.ExtensionSecretListRequest) (shared.ExtensionSecretListResponse, error) {

	if err := a.sprawdzKatalog(); err != nil {
		return shared.ExtensionSecretListResponse{}, err
	}
	kod := strings.TrimSpace(wartoscTekstu(z.ExtensionId))
	if kod != "" {
		if _, err := a.pozycjaRozszerzeniaZadania(ctx, kod); err != nil {
			return shared.ExtensionSecretListResponse{}, err
		}
	}
	// Zawężenie po terminie: „wygasające w N dni” to górna granica czasu.
	// Zero znaczy brak zawężenia.
	granicaCzasu := int64(0)
	if z.ExpiringWithinDays != nil && *z.ExpiringWithinDays > 0 {
		granicaCzasu = a.teraz() + int64(*z.ExpiringWithinDays)*24*60*60*1000
	}

	wiersze, err := a.rejestr.SekretyRozszerzen(ctx, granicaCzasu)
	if err != nil {
		return shared.ExtensionSecretListResponse{}, bladRozszerzenia(err)
	}
	sekrety := make([]shared.SecretRef, 0, len(wiersze))
	for _, wiersz := range wiersze {
		// Zawężenie po pozycji: referencja należy do niej, gdy jest
		// współdzielona albo to ona ją powiązała.
		if kod != "" && !zawieraNapis(wiersz.KodyRozszerzen, kod) {
			integracja, err := a.rejestr.IntegracjaRozszerzenia(ctx, kod)
			if err != nil && !isBrakWierszaApp(err) {
				return shared.ExtensionSecretListResponse{}, bladRozszerzenia(err)
			}
			if integracja.OdwolanieSekretu == nil || *integracja.OdwolanieSekretu != wiersz.Odwolanie {
				continue
			}
		}
		sekrety = append(sekrety, sekretKontraktu(wiersz))
	}
	return shared.ExtensionSecretListResponse{Secrets: sekrety, Total: len(sekrety)}, nil
}

// UdostepnijSekret obsługuje `extension.secret.share`: ustala, które pozycje
// i role mają dostęp do referencji. Treści poświadczenia ta droga nie dotyka —
// zmienia wyłącznie zakres jego współdzielenia.
func (a *adapterRozszerzen) UdostepnijSekret(ctx context.Context,
	z shared.ExtensionSecretShareRequest) (shared.ExtensionSecretShareResponse, error) {

	if err := a.sprawdzKatalog(); err != nil {
		return shared.ExtensionSecretShareResponse{}, err
	}
	odwolanie := strings.TrimSpace(z.SecretRef)
	if odwolanie == "" {
		return shared.ExtensionSecretShareResponse{}, bladWskazaniaRozszerzenia(
			"udostępnienie bez odwołania do sekretu")
	}
	// Pozycje, którym udostępnia się referencję, muszą istnieć: wskazanie
	// nieznanej to dostęp nikomu.
	for _, kod := range z.ExtensionIds {
		if _, err := a.pozycjaRozszerzeniaZadania(ctx, kod); err != nil {
			return shared.ExtensionSecretShareResponse{}, err
		}
	}

	zapisany, err := a.rejestr.ZapiszSekretRozszerzenia(ctx, dane.SekretRozszerzenia{
		Odwolanie: odwolanie, KodyRozszerzen: z.ExtensionIds, KodyRol: z.RoleIds,
		Zaktualizowano: a.teraz(),
	}, true)
	if err != nil {
		return shared.ExtensionSecretShareResponse{}, bladRozszerzenia(err)
	}
	return shared.ExtensionSecretShareResponse{Secret: sekretKontraktu(zapisany)}, nil
}

// sekretKontraktu przekłada wiersz referencji na kształt kontraktu. Treści
// poświadczenia nie ma tu i nie ma jak się pojawić — wiersz jej nie zna.
func sekretKontraktu(wiersz dane.SekretRozszerzenia) shared.SecretRef {
	odwolanie := shared.SecretRef{
		Ref: wiersz.Odwolanie, Label: wiersz.Etykieta, ExpiresAt: wiersz.Wygasa,
		SharedWithExtensionIds: wiersz.KodyRozszerzen, SharedWithRoleIds: wiersz.KodyRol,
	}
	if wiersz.SposobLogowania != nil {
		sposob := shared.ExtensionAuthKind(*wiersz.SposobLogowania)
		odwolanie.AuthKind = &sposob
	}
	return odwolanie
}

// zawieraNapis rozstrzyga obecność napisu w wykazie porównaniem dokładnym,
// bez normalizacji wielkości liter.
func zawieraNapis(wykaz []string, szukany string) bool {
	for _, wpis := range wykaz {
		if wpis == szukany {
			return true
		}
	}
	return false
}

// sprawdzZakresUprawnieniaRozszerzenia dopuszcza wyłącznie zakresy uprawnień
// znane kontraktowi rozszerzeń.
func sprawdzZakresUprawnieniaRozszerzenia(zakres shared.ExtensionPermissionScope) error {
	switch zakres {
	case shared.ExtensionPermissionScopeNetwork, shared.ExtensionPermissionScopeFileRead,
		shared.ExtensionPermissionScopeFileWrite, shared.ExtensionPermissionScopeProcessSpawn,
		shared.ExtensionPermissionScopeSecretRead, shared.ExtensionPermissionScopeModelCall:
		return nil
	}
	return bladWskazaniaRozszerzenia("nieznany zakres uprawnienia " +
		strconv.Quote(string(zakres)) +
		" — dopuszczalne: network, fileRead, fileWrite, processSpawn, secretRead, modelCall")
}
