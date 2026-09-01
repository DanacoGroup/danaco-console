//go:build !windows

// Odpowiedzialność pliku: postać klucza własnego sejfu poza Windows — jawny zapis szesnastkowy, chroniony prawami pliku 0400.
package dane

import "encoding/hex"

// zabezpieczKluczWlasny oddaje klucz w zapisie szesnastkowym; poza Windows nie ma
// magazynu systemowego związanego z kontem, więc chronią go prawa pliku.
func zabezpieczKluczWlasny(surowy []byte) ([]byte, error) {
	return []byte(hex.EncodeToString(surowy) + "\n"), nil
}

// odbezpieczKluczWlasny mówi, że plik nie jest zapieczętowany — czyta go rozbierzKlucz.
func odbezpieczKluczWlasny(_ []byte) ([]byte, bool, error) {
	return nil, false, nil
}
