// Odpowiedzialność pliku: przekład wartości wyliczeniowych między kontraktem
// a kolumnami bazy. Słowniki pochodzą wyłącznie z pakietu `shared` —
// warstwa trwałości nie powiela literałów ani nie prowadzi własnego przekładu.
package dane

import (
	"fmt"

	"danacoconsole/shared"
)

// naBaze przekłada wartość kontraktu na wartość kolumny. Pusta wartość oznacza
// „nie ustawiono” i daje wartość domyślną kolumny.
func naBaze[T ~string](slownik map[T]string, wartosc T, domyslna T, pole string) (string, error) {
	if wartosc == "" {
		wartosc = domyslna
	}
	kolumna, znana := slownik[wartosc]
	if !znana {
		return "", fmt.Errorf("dane: wartość %q nie należy do słownika kontraktu pola %s", string(wartosc), pole)
	}
	return kolumna, nil
}

// zBazy przekłada wartość kolumny na wartość kontraktu.
func zBazy[T ~string](slownik map[string]T, kolumna string, pole string) (T, error) {
	wartosc, znana := slownik[kolumna]
	if !znana {
		var pusta T
		return pusta, fmt.Errorf("dane: wartość kolumny %q pola %s nie ma odpowiednika w kontrakcie", kolumna, pole)
	}
	return wartosc, nil
}

// stanSesjiNaBaze i pozostałe funkcje nazwane wprost wiążą pole struktury
// z jego słownikiem, dzięki czemu repozytorium nie wybiera słownika samo.
func stanSesjiNaBaze(stan shared.SessionStatus) (string, error) {
	return naBaze(shared.WartosciBazySessionStatus, stan, shared.SessionStatusActive, "sesja.stan")
}

func stanSesjiZBazy(kolumna string) (shared.SessionStatus, error) {
	return zBazy(shared.WartosciKontraktuSessionStatus, kolumna, "sesja.stan")
}

func stanOknaNaBaze(stan shared.WindowStatus) (string, error) {
	return naBaze(shared.WartosciBazyWindowStatus, stan, shared.WindowStatusOpen, "okno_komunikacji.stan")
}

func stanOknaZBazy(kolumna string) (shared.WindowStatus, error) {
	return zBazy(shared.WartosciKontraktuWindowStatus, kolumna, "okno_komunikacji.stan")
}

func rolaOknaNaBaze(rola shared.WindowRole) (string, error) {
	return naBaze(shared.WartosciBazyWindowRole, rola, shared.WindowRoleStandalone, "okno_komunikacji.rola_okna")
}

func rolaOknaZBazy(kolumna string) (shared.WindowRole, error) {
	return zBazy(shared.WartosciKontraktuWindowRole, kolumna, "okno_komunikacji.rola_okna")
}

func srodowiskoWykonaniaNaBaze(srodowisko shared.ExecutionEnv) (string, error) {
	return naBaze(shared.WartosciBazyExecutionEnv, srodowisko, shared.ExecutionEnvLocal,
		"okno_komunikacji.srodowisko_wykonania")
}

func srodowiskoWykonaniaZBazy(kolumna string) (shared.ExecutionEnv, error) {
	return zBazy(shared.WartosciKontraktuExecutionEnv, kolumna, "okno_komunikacji.srodowisko_wykonania")
}

func trybUprawnienNaBaze(tryb shared.PermissionMode) (string, error) {
	return naBaze(shared.WartosciBazyPermissionMode, tryb, shared.PermissionModeManual,
		"okno_komunikacji.tryb_uprawnien")
}

func trybUprawnienZBazy(kolumna string) (shared.PermissionMode, error) {
	return zBazy(shared.WartosciKontraktuPermissionMode, kolumna, "okno_komunikacji.tryb_uprawnien")
}

func rolaWiadomosciNaBaze(rola shared.MessageRole) (string, error) {
	return naBaze(shared.WartosciBazyMessageRole, rola, shared.MessageRoleUser, "wiadomosc.rola")
}

func rolaWiadomosciZBazy(kolumna string) (shared.MessageRole, error) {
	return zBazy(shared.WartosciKontraktuMessageRole, kolumna, "wiadomosc.rola")
}

func stanWiadomosciNaBaze(stan shared.MessageStatus) (string, error) {
	return naBaze(shared.WartosciBazyMessageStatus, stan, shared.MessageStatusPending, "wiadomosc.stan")
}

func stanWiadomosciZBazy(kolumna string) (shared.MessageStatus, error) {
	return zBazy(shared.WartosciKontraktuMessageStatus, kolumna, "wiadomosc.stan")
}

func rodzajTresciNaBaze(rodzaj shared.ChunkKind) (string, error) {
	return naBaze(shared.WartosciBazyChunkKind, rodzaj, shared.ChunkKindText, "wiadomosc.rodzaj_tresci")
}

func rodzajTresciZBazy(kolumna string) (shared.ChunkKind, error) {
	return zBazy(shared.WartosciKontraktuChunkKind, kolumna, "wiadomosc.rodzaj_tresci")
}

func stanKolejkiNaBaze(stan shared.QueueStatus) (string, error) {
	return naBaze(shared.WartosciBazyQueueStatus, stan, shared.QueueStatusIdle, "kolejka.stan")
}

func stanKolejkiZBazy(kolumna string) (shared.QueueStatus, error) {
	return zBazy(shared.WartosciKontraktuQueueStatus, kolumna, "kolejka.stan")
}

func poziomZasieguNaBaze(poziom shared.ConfigScope) (string, error) {
	return naBaze(shared.WartosciBazyConfigScope, poziom, shared.ConfigScopeGlobal, "poziom_zasiegu.kod")
}

func poziomZasieguZBazy(kolumna string) (shared.ConfigScope, error) {
	return zBazy(shared.WartosciKontraktuConfigScope, kolumna, "poziom_zasiegu.kod")
}

// Oś rozstrzygania nie ma w kontrakcie słownika przekładu, bo kolumny
// `os_zasiegu.kod` i `ustawienie.os` niosą wprost wartości wyliczenia
// ConfigAxis. Mapy poniżej powstają ze stałych pakietu shared, więc nazwy osi
// nie ma tu jako literału.
var wartosciBazyConfigAxis = map[shared.ConfigAxis]string{
	shared.ConfigAxisPlatform: shared.ConfigAxisPlatform,
	shared.ConfigAxisModel:    shared.ConfigAxisModel,
	shared.ConfigAxisAccount:  shared.ConfigAxisAccount,
}

var wartosciKontraktuConfigAxis = map[string]shared.ConfigAxis{
	shared.ConfigAxisPlatform: shared.ConfigAxisPlatform,
	shared.ConfigAxisModel:    shared.ConfigAxisModel,
	shared.ConfigAxisAccount:  shared.ConfigAxisAccount,
}

func osZasieguNaBaze(os shared.ConfigAxis) (string, error) {
	return naBaze(wartosciBazyConfigAxis, os, shared.ConfigAxisPlatform, "ustawienie.os")
}

func osZasieguZBazy(kolumna string) (shared.ConfigAxis, error) {
	return zBazy(wartosciKontraktuConfigAxis, kolumna, "ustawienie.os")
}
