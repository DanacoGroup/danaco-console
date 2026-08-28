// Plik przekłada wartości wyliczeniowe tożsamości modelu między kontraktem
// a kolumnami tabel kategoria_tozsamosci i dokument_tozsamosci, sprawdzając
// przynależność do zbioru kontraktu i uzupełniając wartość domyślną.
package dane

import (
	"fmt"

	"danacoconsole/shared"
)

// warstwyTozsamosci zestawia trzy warstwy nakładki tożsamości modelu w
// kolejności według krytyczności: konstytucja, profil, kompetencja.
var warstwyTozsamosci = map[string]shared.IdentityLayer{
	shared.IdentityLayerConstitution: shared.IdentityLayerConstitution,
	shared.IdentityLayerProfile:      shared.IdentityLayerProfile,
	shared.IdentityLayerExpertise:    shared.IdentityLayerExpertise,
}

// trybyTozsamosci zestawia dwa tryby podania tożsamości modelu. Tryb ZASTAP
// zastępuje prompt fabryczny w całości, tryb DOLACZ dokłada treść do niego.
var trybyTozsamosci = map[string]shared.IdentityMode{
	shared.IdentityModeZASTAP: shared.IdentityModeZASTAP,
	shared.IdentityModeDOLACZ: shared.IdentityModeDOLACZ,
}

// osieTozsamosci zestawia trzy osie rozstrzygania treści tożsamości modelu:
// platforma, wskazany model i wskazane konto.
var osieTozsamosci = map[string]shared.ConfigAxis{
	shared.ConfigAxisPlatform: shared.ConfigAxisPlatform,
	shared.ConfigAxisModel:    shared.ConfigAxisModel,
	shared.ConfigAxisAccount:  shared.ConfigAxisAccount,
}

func warstwaTozsamosciZBazy(kolumna string) (shared.IdentityLayer, error) {
	return zBazy(warstwyTozsamosci, kolumna, "kategoria_tozsamosci.warstwa")
}

func trybTozsamosciZBazy(kolumna string) (shared.IdentityMode, error) {
	return zBazy(trybyTozsamosci, kolumna, "dokument_tozsamosci.tryb")
}

// trybTozsamosciNaBaze przekłada tryb na wartość kolumny. Brak trybu znaczy
// ZASTAP — tożsamość jest zamieniana, nie dołączana (wartość domyślna klucza
// `tozsamosc.tryb_domyslny`).
func trybTozsamosciNaBaze(tryb shared.IdentityMode) (string, error) {
	return naBaze(slownikBazyTozsamosci(trybyTozsamosci), tryb, shared.IdentityModeZASTAP,
		"dokument_tozsamosci.tryb")
}

func osTozsamosciZBazy(kolumna string) (shared.ConfigAxis, error) {
	return zBazy(osieTozsamosci, kolumna, "dokument_tozsamosci.os")
}

// osTozsamosciNaBaze przekłada oś na wartość kolumny. Brak osi znaczy platform
// (kontrakt: „oś pominięta znaczy platform”).
func osTozsamosciNaBaze(os shared.ConfigAxis) (string, error) {
	return naBaze(slownikBazyTozsamosci(osieTozsamosci), os, shared.ConfigAxisPlatform,
		"dokument_tozsamosci.os")
}

// slownikBazyTozsamosci odwraca słownik rozpoznania na słownik zapisu. Kolumna trzyma
// wartość kontraktu wprost, więc odwrócenie jest tożsamościowe — powstaje po to,
// by zapis szedł tą samą drogą sprawdzenia co odczyt, a nie obok niej.
func slownikBazyTozsamosci[T ~string](slownik map[string]T) map[T]string {
	odwrocony := make(map[T]string, len(slownik))
	for kolumna, wartosc := range slownik {
		odwrocony[wartosc] = kolumna
	}
	return odwrocony
}

// sprawdzBytOsiTozsamosci pilnuje zgody osi z jej bytem: oś platform nie ma bytu, oś
// model i konto muszą go wskazać. Ten sam warunek stoi w schemacie bazy
// — tu wraca komunikatem zrozumiałym dla okna konfiguracji, zamiast surowym
// naruszeniem CHECK.
func sprawdzBytOsiTozsamosci(os shared.ConfigAxis, byt string) error {
	if os == shared.ConfigAxisPlatform && byt != "" {
		return fmt.Errorf("dane: oś %q nie ma bytu, a podano %q", string(os), byt)
	}
	if os != shared.ConfigAxisPlatform && byt == "" {
		return fmt.Errorf("dane: oś %q wymaga wskazania bytu osi", string(os))
	}
	return nil
}
