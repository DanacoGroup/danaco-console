package injection

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// zapasNaZamkniecie daje procesowi chwilę na domknięcie potoków po odwołaniu
// kontekstu; bez tego oczekiwanie potrafi wisieć na niedomkniętym wyjściu.
const zapasNaZamkniecie = 5 * time.Second

// limitBledow ogranicza zapamiętane wyjście diagnostyczne procesu. Treść służy
// rozpoznaniu przyczyny, nie archiwum.
const limitBledow = 64 << 10

// Proces jest jednym uruchomieniem programu `claude` w trybie stream-json,
// udostępniającym identyfikator procesu i honorującym odwołanie kontekstu.
type Proces struct {
	start *Start
}

// Uruchom startuje program z podanym wierszem argumentów i katalogiem
// konfiguracji konta. CLAUDE_CONFIG_DIR jest jedynym nośnikiem tożsamości
// konta; ścieżka nie jest sekretem, sekrety zostają w profilu.
func Uruchom(kontekst context.Context, u Ustawienia, argv []string, katalogKonfiguracji string) (*Proces, error) {
	start, err := Wystartuj(kontekst, Rozruch{
		Program:           u.Program,
		Argumenty:         argv,
		Katalog:           u.KatalogRoboczy,
		Srodowisko:        srodowisko(u.Srodowisko, katalogKonfiguracji),
		ZapasNaZamkniecie: zapasNaZamkniecie,
	})
	if err != nil {
		return nil, err
	}
	return &Proces{start: start}, nil
}

// Pid zwraca identyfikator procesu — warstwa sesji podpina go pod obiekt
// zadania systemu i tylko ona ubija drzewo.
func (p *Proces) Pid() int {
	return p.start.Pid()
}

// Wyslij podaje wypowiedź użytkownika na wejście procesu w kształcie
// JSON-lines wymaganym przez --input-format stream-json.
func (p *Proces) Wyslij(tekst string) error {
	linia, err := json.Marshal(wiadomoscWejsciowa{
		Type:    "user",
		Message: trescWejsciowa{Role: "user", Content: tekst},
	})
	if err != nil {
		return fmt.Errorf("injection: kodowanie wypowiedzi: %w", err)
	}
	if _, err := p.start.Wejscie().Write(append(linia, '\n')); err != nil {
		return fmt.Errorf("injection: zapis na wejście procesu: %w", err)
	}
	return nil
}

// ZamknijWejscie domyka strumień wejściowy procesu — program kończy turę
// dopiero po zobaczeniu końca wejścia.
func (p *Proces) ZamknijWejscie() error {
	return p.start.Wejscie().Close()
}

// Wyjscie daje strumień wyjściowy procesu do odczytu linia po linii,
// w kształcie strumienia JSON-lines.
func (p *Proces) Wyjscie() io.Reader {
	return p.start.Wyjscie()
}

// Bledy zwraca zapamiętane wyjście diagnostyczne procesu, ograniczone
// limitem pojemności zapisanego bufora.
func (p *Proces) Bledy() string {
	return p.start.Bledy()
}

// Czekaj czeka na zakończenie procesu. Kod wyjścia różny od zera jest błędem
// bieżącego wywołania, nie awarią sesji ani konta.
func (p *Proces) Czekaj() error {
	return p.start.Czekaj()
}

// srodowisko składa środowisko procesu: dziedziczone, dołożone z ustawień,
// a na końcu katalog konfiguracji konta, który wygrywa z każdym wcześniejszym.
func srodowisko(dodatkowe map[string]string, katalogKonfiguracji string) []string {
	zmienne := make(map[string]string, len(dodatkowe)+8)
	kolejnosc := make([]string, 0, len(dodatkowe)+8)
	dodaj := func(klucz, wartosc string) {
		if _, jest := zmienne[klucz]; !jest {
			kolejnosc = append(kolejnosc, klucz)
		}
		zmienne[klucz] = wartosc
	}
	for _, wpis := range os.Environ() {
		if klucz, wartosc, jest := strings.Cut(wpis, "="); jest {
			dodaj(klucz, wartosc)
		}
	}
	for klucz, wartosc := range dodatkowe {
		dodaj(klucz, wartosc)
	}
	if katalogKonfiguracji != "" {
		dodaj("CLAUDE_CONFIG_DIR", katalogKonfiguracji)
	}
	wynik := make([]string, 0, len(kolejnosc))
	for _, klucz := range kolejnosc {
		wynik = append(wynik, klucz+"="+zmienne[klucz])
	}
	return wynik
}

// buforBledow zbiera wyjście diagnostyczne procesu z górnym ograniczeniem
// pojemności zapamiętanego bufora.
type buforBledow struct {
	mu    sync.Mutex
	tresc []byte
}

func (b *buforBledow) Write(dane []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if wolne := limitBledow - len(b.tresc); wolne > 0 {
		if len(dane) < wolne {
			wolne = len(dane)
		}
		b.tresc = append(b.tresc, dane[:wolne]...)
	}
	return len(dane), nil
}

func (b *buforBledow) Tresc() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return string(b.tresc)
}

type wiadomoscWejsciowa struct {
	Type    string         `json:"type"`
	Message trescWejsciowa `json:"message"`
}

type trescWejsciowa struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
