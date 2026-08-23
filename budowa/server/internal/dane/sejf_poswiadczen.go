// Odpowiedzialność pliku: sejf poświadczeń — magazyn sekretów bytów rdzenia
// (kont i punktów dostępu) trzymany poza bazą produktu.
//
// Baza zna wyłącznie odwołanie — nazwę wpisu w sejfie. Sam
// sekret leży w osobnym pliku katalogu danych, z prawami tylko dla właściciela.
// Dzięki temu odczyt katalogu kont nigdy nie może wynieść sekretu (kolumny na
// niego nie ma), a mimo to poświadczenie jest trwałe i `hasCredential` mówi
// prawdę: skoro odwołanie zapisano, sekret istnieje.
//
// Sejf jest celowo prosty: jeden plik JSON kluczowany bytem, pod zamkiem. To nie
// jest magazyn klasy KMS — jest to trwały schowek na sekret, którego rdzeń nie
// wpuszcza do bazy ani do odpowiedzi. Wymianę na zewnętrzny magazyn (Windows
// Credential Manager, plik profilu) domyka ten sam interfejs Zapisz/Usun.
package dane

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// nazwaPlikuSejfu to plik sekretów w katalogu danych rdzenia.
const nazwaPlikuSejfu = "poswiadczenia.sejf"

// przedrostekOdwolania znakuje odwołanie oddawane bazie, żeby było widać, że
// wskazuje wpis sejfu, a nie zmienną środowiskową ani ścieżkę profilu.
const przedrostekOdwolania = "sejf:"

// SejfPlikowy jest sejfem poświadczeń opartym o plik katalogu danych. Spełnia
// port core.SejfPoswiadczen strukturalnie (metody Zapisz/Usun) — pakiet dane
// nie importuje pakietu core.
type SejfPlikowy struct {
	mu      sync.Mutex
	sciezka string
}

// NowySejfPlikowy składa sejf nad plikiem w podanym katalogu danych. Katalog nie
// musi jeszcze istnieć — powstanie przy pierwszym zapisie.
func NowySejfPlikowy(katalogDanych string) *SejfPlikowy {
	return &SejfPlikowy{sciezka: filepath.Join(katalogDanych, nazwaPlikuSejfu)}
}

// Zapisz umieszcza poświadczenie bytu w sejfie i zwraca odwołanie do niego.
// Odwołanie jest stabilne dla bytu — ponowny zapis tego samego bytu nadpisuje
// sekret i oddaje to samo odwołanie.
func (s *SejfPlikowy) Zapisz(_ context.Context, byt, poswiadczenie string) (string, error) {
	byt = strings.TrimSpace(byt)
	if byt == "" {
		return "", fmt.Errorf("dane: sejf: zapis bez wskazania bytu")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	wpisy, err := s.wczytaj()
	if err != nil {
		return "", err
	}
	wpisy[byt] = poswiadczenie
	if err := s.zapisz(wpisy); err != nil {
		return "", err
	}
	return przedrostekOdwolania + byt, nil
}

// Odczytaj zwraca poświadczenie bytu i znacznik, czy wpis w sejfie istnieje.
// Byt jest tu surowym kluczem wpisu — bez przedrostka odwołania: rozbiera go
// wołający (kanał API), a sejf kluczuje bytem dokładnie tak, jak zapisał
// w Zapisz. Brak pliku, brak wpisu i nieczytelna treść dają ten sam wynik
// „nie ma" bez błędu: warstwa wyżej odmawia wtedy wywołania tak samo
// jak przy pustej zmiennej środowiskowej, zamiast wywracać się na braku sekretu.
func (s *SejfPlikowy) Odczytaj(_ context.Context, byt string) (string, bool) {
	byt = strings.TrimSpace(byt)
	if byt == "" {
		return "", false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	wpisy, err := s.wczytaj()
	if err != nil {
		return "", false
	}
	poswiadczenie, jest := wpisy[byt]
	return poswiadczenie, jest
}

// Usun kasuje poświadczenie bytu. Brak wpisu nie jest błędem.
func (s *SejfPlikowy) Usun(_ context.Context, byt string) error {
	byt = strings.TrimSpace(byt)
	if byt == "" {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	wpisy, err := s.wczytaj()
	if err != nil {
		return err
	}
	if _, jest := wpisy[byt]; !jest {
		return nil
	}
	delete(wpisy, byt)
	return s.zapisz(wpisy)
}

// wczytaj odczytuje mapę sekretów z pliku. Brak pliku znaczy sejf pusty, nie
// błąd — pierwszy zapis dopiero go założy. Wołane pod zamkiem.
func (s *SejfPlikowy) wczytaj() (map[string]string, error) {
	surowe, err := os.ReadFile(s.sciezka)
	if os.IsNotExist(err) {
		return map[string]string{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("dane: sejf: odczyt %s: %w", s.sciezka, err)
	}
	wpisy := map[string]string{}
	if len(strings.TrimSpace(string(surowe))) == 0 {
		return wpisy, nil
	}
	if err := json.Unmarshal(surowe, &wpisy); err != nil {
		return nil, fmt.Errorf("dane: sejf: treść %s nieczytelna: %w", s.sciezka, err)
	}
	return wpisy, nil
}

// zapisz utrwala mapę sekretów z prawami tylko dla właściciela. Zapis idzie przez
// plik tymczasowy i przemianowanie, żeby awaria w połowie nie zostawiła pliku
// obciętego. Wołane pod zamkiem.
func (s *SejfPlikowy) zapisz(wpisy map[string]string) error {
	if err := os.MkdirAll(filepath.Dir(s.sciezka), 0o700); err != nil {
		return fmt.Errorf("dane: sejf: katalog %s: %w", filepath.Dir(s.sciezka), err)
	}
	surowe, err := json.Marshal(wpisy)
	if err != nil {
		return fmt.Errorf("dane: sejf: kodowanie: %w", err)
	}
	tymczasowy := s.sciezka + ".tmp"
	if err := os.WriteFile(tymczasowy, surowe, 0o600); err != nil {
		return fmt.Errorf("dane: sejf: zapis %s: %w", tymczasowy, err)
	}
	if err := os.Rename(tymczasowy, s.sciezka); err != nil {
		return fmt.Errorf("dane: sejf: przemianowanie %s: %w", s.sciezka, err)
	}
	return nil
}
