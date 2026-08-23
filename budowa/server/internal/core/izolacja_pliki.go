// Odpowiedzialność pliku: egzekucja zakresu `odczyt_zapis_plikow` — dostępu do
// plików poza własnym katalogiem roboczym okna.
//
// Obszar dozwolony bierze się z dwóch miejsc: własnego katalogu roboczego okna
// oraz korzeni nadań dostępu tego okna (nadanie żyje per okno). Nadanie niesie
// także tryb, więc korzeń nadany do odczytu nie staje się korzeniem do zapisu.
//
// Straż egzekwuje ścieżki, po które rdzeń sięga w imieniu okna, oraz wykaz korzeni
// podawany procesowi modelu przy uruchomieniu. Nie zabroni uruchomionemu procesowi
// otworzyć pliku samodzielnie — to leży poza zasięgiem rdzenia i wymaga środków
// systemu operacyjnego.
package core

import (
	"strings"

	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// korzenDozwolony to jeden korzeń jawnie udostępniony oknu wraz z trybem.
type korzenDozwolony struct {
	// Sciezka korzenia.
	Sciezka string
	// Zapis mówi, czy nadanie pozwala zmieniać zawartość korzenia.
	Zapis bool
}

// strazPlikow rozstrzyga, czy ścieżka mieści się w obszarze plikowym okna.
type strazPlikow struct {
	zasady   session.Zasady
	wlasny   string
	korzenie []korzenDozwolony
}

// nowaStrazPlikow składa straż z zasad izolacji, własnego katalogu roboczego
// okna i kompletu nadań tego okna. Zakres wyłączony daje straż przepuszczającą
// wszystko — pełny dostęp w ramach uprawnień Operatora jest stanem wyjściowym
// platformy.
func nowaStrazPlikow(zasady session.Zasady, katalogWlasny string,
	nadania []NadanieMostu) *strazPlikow {

	return &strazPlikow{
		zasady:   zasady,
		wlasny:   strings.TrimSpace(katalogWlasny),
		korzenie: korzenieDozwolone(nadania),
	}
}

// korzenieDozwolone wylicza korzenie katalogów jawnie dozwolonych oknu wraz
// z trybem. Zawężenie korzeni nadania do korzeni punktu robi KorzenieNadania —
// reguła zawężenia ma w drzewie jedno miejsce. Nadanie nieczynne
// i punkt nieczynny nie wnoszą nic: brak wiersza znaczy brak uprawnienia, nie
// awarię odczytu.
func korzenieDozwolone(nadania []NadanieMostu) []korzenDozwolony {
	korzenie := make([]korzenDozwolony, 0, len(nadania))
	for _, nadanie := range nadania {
		if !nadanie.Nadanie.Enabled || !nadanie.Punkt.Enabled {
			continue
		}
		if nadanie.Punkt.Kind != shared.AccessPointKindLocalDirectory {
			continue
		}
		for _, korzen := range KorzenieNadania(nadanie) {
			korzenie = append(korzenie, korzenDozwolony{
				Sciezka: korzen,
				Zapis:   nadanie.Nadanie.Mode == shared.AccessModeWrite,
			})
		}
	}
	return korzenie
}

// sprawdzOdczyt odrzuca odczyt ścieżki spoza obszaru okna.
func (s *strazPlikow) sprawdzOdczyt(sciezka string) error {
	return s.sprawdz(sciezka, false)
}

// sprawdzZapis odrzuca zapis ścieżki spoza obszaru okna oraz zapis w korzeniu
// nadanym wyłącznie do odczytu.
func (s *strazPlikow) sprawdzZapis(sciezka string) error {
	return s.sprawdz(sciezka, true)
}

// Katalogi przepuszcza wykaz katalogów podawany procesowi modelu. Katalog spoza
// obszaru okna zatrzymuje uruchomienie — wykaz obcięty po cichu dawałby
// Operatorowi obraz niezgodny z tym, co dostał proces.
func (s *strazPlikow) Katalogi(zadane []string) ([]string, error) {
	if s == nil || !s.zasady.Pliki {
		return zadane, nil
	}
	for _, katalog := range zadane {
		if err := s.sprawdzOdczyt(katalog); err != nil {
			return nil, err
		}
	}
	return zadane, nil
}

// sprawdz rozstrzyga jedną ścieżkę. Straż niewpięta i zakres wyłączony nie
// ograniczają niczego.
func (s *strazPlikow) sprawdz(sciezka string, zapis bool) error {
	if s == nil || !s.zasady.Pliki {
		return nil
	}
	if strings.TrimSpace(sciezka) == "" {
		return session.NoweNaruszenie(konfig.KluczIzolacjaPliki,
			"ścieżka pusta nie mieści się w żadnym korzeniu dozwolonym oknu")
	}
	if session.SciezkaWewnatrz(s.wlasny, sciezka) {
		return nil
	}
	tylkoOdczyt := false
	for _, korzen := range s.korzenie {
		if !session.SciezkaWewnatrz(korzen.Sciezka, sciezka) {
			continue
		}
		if korzen.Zapis || !zapis {
			return nil
		}
		tylkoOdczyt = true
	}
	if tylkoOdczyt {
		return session.NoweNaruszenie(konfig.KluczIzolacjaPliki,
			"ścieżka "+sciezka+" jest nadana oknu wyłącznie do odczytu")
	}
	return session.NoweNaruszenie(konfig.KluczIzolacjaPliki,
		"ścieżka "+sciezka+" leży poza własnym katalogiem okna i poza korzeniami jego nadań")
}
