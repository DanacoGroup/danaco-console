// Odpowiedzialność pliku: obszar modułu Developer — okno komunikacji, katalogi robocze,
// sprawdzenie ścieżki wobec obszaru, brama trybu uprawnień i kody odmów.
package core

import (
	"errors"
	"path/filepath"
	"runtime"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

func (a *adapterDevelopera) oknoDevelopera(oknoKod string) (session.Okno, error) {
	kod := strings.TrimSpace(oknoKod)
	if kod == "" {
		return session.Okno{}, bladZadaniaDevelopera("czynność wymaga wskazania okna modułu")
	}
	if a.okna == nil {
		return session.Okno{}, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"moduł Developer: serwer nie ma rejestru okien, więc nie zna obszaru pracy okna"))
	}
	okno, err := a.okna.Okno(kod)
	if err != nil {
		return session.Okno{}, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Developer: okno "+kod+" nie jest otwarte"))
	}
	return okno, nil
}

func (a *adapterDevelopera) korzenieOkna(okno session.Okno) []string {
	korzenie := make([]string, 0, len(okno.KatalogiRobocze)+1)
	for _, katalog := range okno.KatalogiRobocze {
		if czysty := oczyscKorzen(katalog); czysty != "" {
			korzenie = append(korzenie, czysty)
		}
	}
	if len(korzenie) > 0 {
		return korzenie
	}
	if a.katalog == nil {
		return nil
	}
	zastepczy := oczyscKorzen(a.katalog.Ustal(konfig.Kontekst{Okno: okno.Id}, okno.IdSesji).Sciezka)
	if zastepczy == "" {
		return nil
	}
	return []string{zastepczy}
}

func (a *adapterDevelopera) obszarDevelopera(okno session.Okno) session.Obszar {
	if a.katalog == nil {
		return session.Obszar{IdOkna: okno.Id}
	}
	return ObszarOkna(a.katalog.Ustal(konfig.Kontekst{Okno: okno.Id}, okno.IdSesji), okno.Id)
}

func sciezkaWObszarze(korzenie []string, wskazanie string, istnieje func(string) bool) (string, error) {
	tresc := strings.TrimSpace(wskazanie)
	if tresc == "" {
		return "", bladZadaniaDevelopera("czynność wymaga wskazania ścieżki")
	}
	if len(korzenie) == 0 {
		return "", protocol.JakoError(protocol.NowyBlad(shared.ErrorCodePermissionDenied,
			"moduł Developer: okno nie ma ani jednego katalogu roboczego, "+
				"więc nie ma obszaru, w którym wolno mu czytać i pisać"))
	}

	if filepath.IsAbs(tresc) {
		return wObszarze(korzenie, filepath.Clean(tresc))
	}
	for _, korzen := range korzenie {
		kandydat := filepath.Join(korzen, tresc)
		if istnieje != nil && istnieje(kandydat) {
			return wObszarze(korzenie, kandydat)
		}
	}
	return wObszarze(korzenie, filepath.Join(korzenie[0], tresc))
}

func wObszarze(korzenie []string, sciezka string) (string, error) {
	for _, korzen := range korzenie {
		if wewnatrzKatalogu(korzen, sciezka) {
			return sciezka, nil
		}
	}
	return "", protocol.JakoError(protocol.NowyBlad(shared.ErrorCodePermissionDenied,
		"moduł Developer: ścieżka "+sciezka+" leży poza katalogami roboczymi okna ("+
			strings.Join(korzenie, ", ")+")"))
}

func wewnatrzKatalogu(korzen, sciezka string) bool {
	k, s := znormalizuj(korzen), znormalizuj(sciezka)
	if k == "" || s == "" {
		return false
	}
	if k == s {
		return true
	}
	return strings.HasPrefix(s, k+string(filepath.Separator))
}

func znormalizuj(sciezka string) string {
	czysta := filepath.Clean(strings.TrimSpace(sciezka))
	if czysta == "." {
		return ""
	}
	if runtime.GOOS == "windows" {
		return strings.ToLower(czysta)
	}
	return czysta
}

// Katalog względny odpada: obszar liczony od katalogu procesu znaczyłby co innego po każdym uruchomieniu.
func oczyscKorzen(katalog string) string {
	tresc := strings.TrimSpace(katalog)
	if tresc == "" || !filepath.IsAbs(tresc) {
		return ""
	}
	return filepath.Clean(tresc)
}

// Tryb planistyczny (plan) wyklucza zapis, repozytorium i budowanie; pozostałe tryby je przepuszczają.
func sprawdzZmianeSystemu(tryb shared.PermissionMode, czynnosc string) error {
	if tryb != shared.PermissionModePlan {
		return nil
	}
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodePermissionDenied,
		"moduł Developer: "+czynnosc+" nie ruszy — okno pracuje w trybie planistycznym (plan), "+
			"który wyklucza zmiany w systemie"))
}

func bladZadaniaDevelopera(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed, "moduł Developer: "+powod))
}

func bladZasobuDevelopera(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound, "moduł Developer: "+powod))
}

func bladDostepuDevelopera(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodePermissionDenied, "moduł Developer: "+powod))
}

func bladWykonaniaDevelopera(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError, "moduł Developer: "+powod))
}

// Kolizja wiersza jest odmową kontraktu (conflict); internal_error klient ponawia (KodyPonawialne).
func bladZapisuDevelopera(err error, powod string) error {
	if errors.Is(err, dane.ErrKolizjaWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeConflict, "moduł Developer: "+powod))
	}
	return bladWykonaniaDevelopera(powod)
}
