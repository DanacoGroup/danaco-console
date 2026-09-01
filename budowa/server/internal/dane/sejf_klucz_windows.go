//go:build windows

// Odpowiedzialność pliku: postać klucza własnego sejfu na Windows — zapieczętowana DPAPI w zakresie użytkownika, żeby kopia katalogu danych na inne konto nie otwierała sejfu (rozstrzygnięcie 30).
package dane

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

// znacznikKluczaDPAPI otwiera plik klucza zapieczętowany DPAPI. Plik bez
// znacznika jest czytany jak poza Windows — zapis szesnastkowy.
const znacznikKluczaDPAPI = "danaco-klucz-dpapi-v1\n"

// zabezpieczKluczWlasny pieczętuje klucz DPAPI w zakresie użytkownika, bez okna
// dialogowego — rdzeń pracuje jako proces poboczny powłoki, bez pulpitu.
func zabezpieczKluczWlasny(surowy []byte) ([]byte, error) {
	wejscie := windows.DataBlob{Size: uint32(len(surowy)), Data: &surowy[0]}
	var wyjscie windows.DataBlob
	if err := windows.CryptProtectData(&wejscie, nil, nil, 0, nil,
		windows.CRYPTPROTECT_UI_FORBIDDEN, &wyjscie); err != nil {
		return nil, fmt.Errorf("DPAPI nie zapieczętowało klucza: %w", err)
	}
	defer windows.LocalFree(windows.Handle(unsafe.Pointer(wyjscie.Data)))
	zapieczetowany := unsafe.Slice(wyjscie.Data, wyjscie.Size)
	plik := make([]byte, 0, len(znacznikKluczaDPAPI)+base64.StdEncoding.EncodedLen(len(zapieczetowany))+1)
	plik = append(plik, znacznikKluczaDPAPI...)
	plik = append(plik, base64.StdEncoding.EncodeToString(zapieczetowany)...)
	plik = append(plik, '\n')
	return plik, nil
}

// odbezpieczKluczWlasny otwiera plik ze znacznikiem DPAPI. Plik bez znacznika
// oddaje false — czyta go rozbierzKlucz.
func odbezpieczKluczWlasny(tresc []byte) ([]byte, bool, error) {
	if !bytes.HasPrefix(tresc, []byte(znacznikKluczaDPAPI)) {
		return nil, false, nil
	}
	zapieczetowany, err := base64.StdEncoding.DecodeString(
		string(bytes.TrimSpace(tresc[len(znacznikKluczaDPAPI):])))
	if err != nil {
		return nil, true, fmt.Errorf("zapis base64 nieczytelny: %w", err)
	}
	if len(zapieczetowany) == 0 {
		return nil, true, fmt.Errorf("plik klucza pusty")
	}
	wejscie := windows.DataBlob{Size: uint32(len(zapieczetowany)), Data: &zapieczetowany[0]}
	var wyjscie windows.DataBlob
	if err := windows.CryptUnprotectData(&wejscie, nil, nil, 0, nil,
		windows.CRYPTPROTECT_UI_FORBIDDEN, &wyjscie); err != nil {
		return nil, true, fmt.Errorf("DPAPI nie otwiera klucza na tym koncie: %w", err)
	}
	defer windows.LocalFree(windows.Handle(unsafe.Pointer(wyjscie.Data)))
	klucz := make([]byte, wyjscie.Size)
	copy(klucz, unsafe.Slice(wyjscie.Data, wyjscie.Size))
	if len(klucz) != dlugoscKlucza {
		return nil, true, fmt.Errorf("klucz ma %d bajtów, wymagane %d", len(klucz), dlugoscKlucza)
	}
	return klucz, true, nil
}
