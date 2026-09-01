// Odpowiedzialność pliku: sejf poświadczeń, magazyn sekretów bytów rdzenia trzymany poza bazą produktu, w osobnym pliku katalogu danych.
package dane

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// nazwaPlikuSejfu to plik sekretów w katalogu danych rdzenia, zapisywany z prawami tylko dla właściciela.
const nazwaPlikuSejfu = "poswiadczenia.sejf"

// przedrostekOdwolania znakuje odwołanie oddawane bazie, żeby było widać, że
// wskazuje wpis sejfu, a nie zmienną środowiskową ani ścieżkę profilu.
const przedrostekOdwolania = "sejf:"

// zmiennaKlucza wskazuje plik klucza pieczętującego sejf. Na serwerze wdrożenia
// klucz leży poza katalogiem danych: kopia bazy ani kopia katalogu danych nie
// wystarczą do odczytania haseł IMAP/SMTP i kluczy API.
const zmiennaKlucza = "DANACO_KLUCZ_SEJFU"

// nazwaPlikuKlucza to klucz własny rdzenia, zakładany w katalogu danych przy
// pierwszym użyciu sejfu, gdy zmienna nie wskazuje klucza. Instalacja na
// urządzeniu Operatora nie ma skąd wziąć zmiennej, a bez klucza nie da się
// założyć pierwszego konta (rozstrzygnięcie 26 w prowadzenie/decyzje.md).
const nazwaPlikuKlucza = "sejf.klucz"

// prawaKlucza to jedyne prawa dopuszczane dla pliku klucza — odczyt wyłącznie
// dla właściciela procesu rdzenia, bez prawa zapisu. Na Windows bity praw nie
// niosą tej treści (system zwraca 0444 albo 0666 wedle atrybutu „tylko do
// odczytu”), dlatego tam sprawdzenie nie obowiązuje.
const prawaKlucza = 0o400

// dlugoscKlucza wymusza AES-256: klucz krótszy zmieniałby siłę pieczęci bez śladu.
const dlugoscKlucza = 32

// znacznikPostaci otwiera plik zapieczętowany. Plik bez tego znacznika pochodzi
// sprzed pieczętowania i jest czytany jako czysty JSON — pierwszy zapis pieczętuje go.
const znacznikPostaci = "danaco-sejf-aes256gcm-v1\n"

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
// sekret i oddaje to samo odwołanie. Bez klucza zapis kończy się odmową; sekret
// nie ląduje na dysku otwartym tekstem nawet wtedy, gdy odmowa zatrzyma czynność
// Operatora.
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

// Odczytaj zwraca poświadczenie bytu i znacznik, czy wpis w sejfie istnieje; brak pliku i brak wpisu dają ten sam wynik bez błędu.
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

// Usun kasuje poświadczenie bytu z sejfu wraz z jego odwołaniem w bazie danych; brak wpisu nie jest błędem.
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
// błąd — pierwszy zapis dopiero go założy. Plik zapieczętowany wymaga klucza;
// plik bez znacznika postaci pochodzi sprzed pieczętowania i idzie jako czysty
// JSON, bo odmowa odczytu nie usunęłaby go z dysku, a odcięłaby dostęp do
// poświadczeń już zapisanych. Wołane pod zamkiem.
func (s *SejfPlikowy) wczytaj() (map[string]string, error) {
	surowe, err := os.ReadFile(s.sciezka)
	if os.IsNotExist(err) {
		return map[string]string{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("dane: sejf: odczyt %s: %w", s.sciezka, err)
	}
	wpisy := map[string]string{}
	if len(bytes.TrimSpace(surowe)) == 0 {
		return wpisy, nil
	}
	if bytes.HasPrefix(surowe, []byte(znacznikPostaci)) {
		klucz, err := s.wczytajKlucz()
		if err != nil {
			return nil, err
		}
		surowe, err = odpieczetuj(klucz, surowe)
		if err != nil {
			return nil, fmt.Errorf("dane: sejf: %s nie otwiera się tym kluczem: %w", s.sciezka, err)
		}
	}
	if err := json.Unmarshal(surowe, &wpisy); err != nil {
		return nil, fmt.Errorf("dane: sejf: treść %s nieczytelna: %w", s.sciezka, err)
	}
	return wpisy, nil
}

// zapisz utrwala mapę sekretów zapieczętowaną kluczem spoza bazy, z prawami
// tylko dla właściciela. Zapis idzie przez plik tymczasowy i przemianowanie, żeby
// awaria w połowie nie zostawiła pliku obciętego. Wołane pod zamkiem.
func (s *SejfPlikowy) zapisz(wpisy map[string]string) error {
	if err := os.MkdirAll(filepath.Dir(s.sciezka), 0o700); err != nil {
		return fmt.Errorf("dane: sejf: katalog %s: %w", filepath.Dir(s.sciezka), err)
	}
	klucz, err := s.wczytajKlucz()
	if err != nil {
		return err
	}
	jawne, err := json.Marshal(wpisy)
	if err != nil {
		return fmt.Errorf("dane: sejf: kodowanie: %w", err)
	}
	zapieczetowane, err := zapieczetuj(klucz, jawne)
	if err != nil {
		return err
	}
	tymczasowy := s.sciezka + ".tmp"
	if err := os.WriteFile(tymczasowy, zapieczetowane, 0o600); err != nil {
		return fmt.Errorf("dane: sejf: zapis %s: %w", tymczasowy, err)
	}
	if err := os.Rename(tymczasowy, s.sciezka); err != nil {
		return fmt.Errorf("dane: sejf: przemianowanie %s: %w", s.sciezka, err)
	}
	return nil
}

// wczytajKlucz podaje klucz pieczęci: ze zmiennej środowiska, a bez niej z pliku
// klucza własnego w katalogu danych, założonego przy pierwszym użyciu. Prawa
// szersze niż odczyt właściciela i zła długość kończą się odmową nazwaną —
// sekret nie ma innej drogi na dysk niż przez ten klucz.
func (s *SejfPlikowy) wczytajKlucz() ([]byte, error) {
	sciezka := strings.TrimSpace(os.Getenv(zmiennaKlucza))
	if sciezka == "" {
		sciezka = s.sciezkaKluczaWlasnego()
		if err := zalozKluczWlasny(sciezka); err != nil {
			return nil, err
		}
	}
	stan, err := os.Stat(sciezka)
	if err != nil {
		return nil, fmt.Errorf("dane: sejf: klucz %s niedostępny: %w", sciezka, err)
	}
	if stan.IsDir() {
		return nil, fmt.Errorf("dane: sejf: klucz %s jest katalogiem, nie plikiem", sciezka)
	}
	if runtime.GOOS != "windows" && stan.Mode().Perm() != prawaKlucza {
		return nil, fmt.Errorf("dane: sejf: klucz %s ma prawa %04o, wymagane %04o",
			sciezka, stan.Mode().Perm(), prawaKlucza)
	}
	tresc, err := os.ReadFile(sciezka)
	if err != nil {
		return nil, fmt.Errorf("dane: sejf: odczyt klucza %s: %w", sciezka, err)
	}
	return rozbierzKlucz(sciezka, tresc)
}

// sciezkaKluczaWlasnego wskazuje klucz zakładany przez rdzeń obok pliku sejfu.
func (s *SejfPlikowy) sciezkaKluczaWlasnego() string {
	return filepath.Join(filepath.Dir(s.sciezka), nazwaPlikuKlucza)
}

// zalozKluczWlasny losuje klucz i zapisuje go szesnastkowo z prawami tylko do
// odczytu dla właściciela. Plik istniejący zostaje nietknięty: nadpisanie
// odcięłoby dostęp do wszystkiego, co sejf już zapieczętował.
func zalozKluczWlasny(sciezka string) error {
	if _, err := os.Stat(sciezka); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("dane: sejf: klucz %s: %w", sciezka, err)
	}
	surowy := make([]byte, dlugoscKlucza)
	if _, err := io.ReadFull(rand.Reader, surowy); err != nil {
		return fmt.Errorf("dane: sejf: nie można wylosować klucza: %w", err)
	}
	tymczasowy := sciezka + ".tmp"
	if err := os.WriteFile(tymczasowy, []byte(hex.EncodeToString(surowy)+"\n"), prawaKlucza); err != nil {
		return fmt.Errorf("dane: sejf: zapis klucza %s: %w", tymczasowy, err)
	}
	if err := os.Rename(tymczasowy, sciezka); err != nil {
		return fmt.Errorf("dane: sejf: przemianowanie klucza %s: %w", sciezka, err)
	}
	return nil
}

// OpisKlucza mówi do dziennika startu, skąd sejf bierze klucz. Nie czyta ani
// nie zakłada klucza — to robi pierwsze użycie sejfu.
func (s *SejfPlikowy) OpisKlucza() string {
	if sciezka := strings.TrimSpace(os.Getenv(zmiennaKlucza)); sciezka != "" {
		return "klucz ze zmiennej " + zmiennaKlucza + ": " + sciezka
	}
	return "klucz własny rdzenia: " + s.sciezkaKluczaWlasnego()
}

// rozbierzKlucz przyjmuje klucz zapisany szesnastkowo albo trzydziestoma dwoma
// bajtami wprost. Skracania ani rozciągania tu nie ma: klucz o innej długości
// zmieniałby siłę pieczęci bez śladu w pliku.
func rozbierzKlucz(sciezka string, tresc []byte) ([]byte, error) {
	oczyszczona := bytes.TrimSpace(tresc)
	if len(oczyszczona) == 2*dlugoscKlucza {
		klucz, err := hex.DecodeString(string(oczyszczona))
		if err != nil {
			return nil, fmt.Errorf("dane: sejf: klucz %s nie jest zapisem szesnastkowym: %w", sciezka, err)
		}
		return klucz, nil
	}
	if len(oczyszczona) == dlugoscKlucza {
		return oczyszczona, nil
	}
	return nil, fmt.Errorf("dane: sejf: klucz %s ma %d bajtów — wymagane %d bajtów albo %d znaków szesnastkowych",
		sciezka, len(oczyszczona), dlugoscKlucza, 2*dlugoscKlucza)
}

// zapieczetuj składa plik sejfu: znacznik postaci, a po nim jednorazowa wartość
// i szyfrogram AES-256-GCM zapisane base64. Jednorazowa wartość idzie przed
// szyfrogramem, żeby odczyt nie potrzebował niczego poza samym plikiem i kluczem.
func zapieczetuj(klucz, jawne []byte) ([]byte, error) {
	pieczec, err := pieczecGCM(klucz)
	if err != nil {
		return nil, err
	}
	jednorazowa := make([]byte, pieczec.NonceSize())
	if _, err := io.ReadFull(rand.Reader, jednorazowa); err != nil {
		return nil, fmt.Errorf("dane: sejf: nie można wylosować wartości jednorazowej: %w", err)
	}
	szyfrogram := pieczec.Seal(jednorazowa, jednorazowa, jawne, nil)
	plik := make([]byte, 0, len(znacznikPostaci)+base64.StdEncoding.EncodedLen(len(szyfrogram))+1)
	plik = append(plik, znacznikPostaci...)
	plik = append(plik, base64.StdEncoding.EncodeToString(szyfrogram)...)
	plik = append(plik, '\n')
	return plik, nil
}

// odpieczetuj odczytuje plik złożony przez zapieczetuj. Błąd rozpieczętowania
// znaczy zły klucz albo naruszoną treść — GCM nie rozróżnia tych dwóch przypadków.
func odpieczetuj(klucz, plik []byte) ([]byte, error) {
	pieczec, err := pieczecGCM(klucz)
	if err != nil {
		return nil, err
	}
	tresc := bytes.TrimSpace(plik[len(znacznikPostaci):])
	szyfrogram, err := base64.StdEncoding.DecodeString(string(tresc))
	if err != nil {
		return nil, fmt.Errorf("zapis base64 nieczytelny: %w", err)
	}
	if len(szyfrogram) < pieczec.NonceSize() {
		return nil, fmt.Errorf("treść krótsza od wartości jednorazowej")
	}
	jednorazowa := szyfrogram[:pieczec.NonceSize()]
	return pieczec.Open(nil, jednorazowa, szyfrogram[pieczec.NonceSize():], nil)
}

// pieczecGCM składa szyfr uwierzytelniający. AES-GCM, a nie sam AES, bo plik
// sejfu leży poza bazą i nikt poza tą pieczęcią nie wykaże, że nie został podmieniony.
func pieczecGCM(klucz []byte) (cipher.AEAD, error) {
	blok, err := aes.NewCipher(klucz)
	if err != nil {
		return nil, fmt.Errorf("dane: sejf: klucz odrzucony przez szyfr: %w", err)
	}
	pieczec, err := cipher.NewGCM(blok)
	if err != nil {
		return nil, fmt.Errorf("dane: sejf: nie można złożyć pieczęci: %w", err)
	}
	return pieczec, nil
}
