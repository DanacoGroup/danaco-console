// Pakiet zdalne przenosi pliki torem zdalnym po tym samym SSH i tej samej zgodzie per host, co tor procesu.
package zdalne

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/mowa"
)

// Kierunek ruchu bajtów względem maszyny rdzenia — wartości kolumny
// `zdalne_przeniesienie.kierunek` co do znaku.
type Kierunek string

const (
	// Wyslanie oznacza ruch pliku z maszyny rdzenia na wskazany host zdalny przez dokładnie ten sam tor SSH.
	Wyslanie Kierunek = "wyslanie"
	// Pobranie oznacza ruch pliku z hosta zdalnego z powrotem na maszynę rdzenia przez dokładnie ten sam tor.
	Pobranie Kierunek = "pobranie"
)

// PrzeniesPlik przenosi jeden plik między maszyną rdzenia a hostem zdalnym
// i odnotowuje wykonany ruch w prowenancji przenosin. Zgoda, wykaz hostów
// i odmowy są wspólne z torem procesu (Przeloz) — druga droga zgody nie
// istnieje.
func PrzeniesPlik(kontekst context.Context, nazwaHosta, idOkna string,
	kierunek Kierunek, sciezkaLokalna, sciezkaZdalna string) error {

	if kierunek != Wyslanie && kierunek != Pobranie {
		return fmt.Errorf("zdalne: kierunek przeniesienia %q spoza słownika "+
			"(wyslanie|pobranie)", string(kierunek))
	}
	if err := odmowDzwieku(sciezkaLokalna, sciezkaZdalna); err != nil {
		return err
	}
	host, err := hostZRejestru(kontekst, nazwaHosta)
	if err != nil {
		return err
	}
	sciezkaSCP, err := exec.LookPath("scp")
	if err != nil {
		return fmt.Errorf("zdalne: plik %s nie został przeniesiony, bo maszyna serwera "+
			"nie ma programu scp w PATH — przenosiny jadą wyłącznie po SSH i wymagają "+
			"jego instalacji: %w", sciezkaLokalna, err)
	}

	plikKluczy, err := zapiszZnaneHosty(host)
	if err != nil {
		return err
	}

	zrodlo, cel := sciezkaLokalna, host.adresSCP(sciezkaZdalna)
	if kierunek == Pobranie {
		zrodlo, cel = host.adresSCP(sciezkaZdalna), sciezkaLokalna
	}
	argumenty := append(append([]string{"-q"}, opcjeSSH...),
		"-o", "UserKnownHostsFile="+plikKluczy,
		"-P", strconv.Itoa(port(host)), zrodlo, cel)
	if wyjscie, err := exec.CommandContext(kontekst, sciezkaSCP, argumenty...).CombinedOutput(); err != nil {
		return fmt.Errorf("zdalne: przeniesienie %s → %s nie powiodło się: %w (%s)",
			zrodlo, cel, err, strings.TrimSpace(string(wyjscie)))
	}
	return odnotujPrzeniesienie(kontekst, host, idOkna, kierunek, zrodlo, cel, rozmiar(sciezkaLokalna))
}

// odmowDzwieku egzekwuje regułę o nagraniach — po obu końcach
// drogi, bo pobranie nagrania Z hosta byłoby tym samym ruchem w drugą stronę.
func odmowDzwieku(sciezki ...string) error {
	for _, sciezka := range sciezki {
		if mowa.FormatPrzyjmowany(sciezka) {
			return fmt.Errorf("zdalne: plik %s nie został przeniesiony, bo jest nagraniem "+
				"dźwięku (wykaz formatów mowa.FormatyNagran), a dźwięk nie opuszcza maszyny "+
				"Operatora — rozstrzygnięcie Właściciela; transkrypcję wykonuje silnik mowy "+
				"na maszynie serwera i do hosta zdalnego jedzie wyłącznie tekst", sciezka)
		}
	}
	return nil
}

// Metoda adresSCP składa zdalny koniec drogi w postaci adresu użytkownika,
// hosta oraz jego ścieżki pliku. Ścieżkę czyta powłoka logowania hosta, więc
// idzie zacytowana — inaczej spacja rozbiłaby ją na dwie, a średnik dopisałby
// do przenosin polecenie.
func (h Host) adresSCP(sciezka string) string {
	przod := h.AdresPolaczenia()
	if uzytkownik := strings.TrimSpace(h.Uzytkownik); uzytkownik != "" {
		przod = uzytkownik + "@" + przod
	}
	return przod + ":" + cytuj(sciezka)
}

// Funkcja port zwraca port danego hosta wskazany w jego danych; wartość zero schodzi na domyślny port SSH.
func port(h Host) int {
	if h.Port == 0 {
		return 22
	}
	return h.Port
}

// Funkcja rozmiar mierzy rozmiar wskazanego pliku lokalnego; brak możliwości pomiaru daje zero, nie odmowę.
func rozmiar(sciezka string) int64 {
	opis, err := os.Stat(sciezka)
	if err != nil {
		return 0
	}
	return opis.Size()
}

// odnotujPrzeniesienie zapisuje wiersz prowenancji przenosin. Zapis następuje
// po ruchu bajtów i niczego nie steruje — to dziennik, nie kolejka.
func odnotujPrzeniesienie(kontekst context.Context, h Host, idOkna string, kierunek Kierunek,
	zrodlo, cel string, bajty int64) error {
	db := baza()
	if db == nil {
		return odmowaBrakuZasilenia()
	}
	_, err := db.ExecContext(kontekst, `INSERT INTO zdalne_przeniesienie
	                   (host_id, okno_id, kierunek, sciezka_zrodlowa, sciezka_docelowa, rozmiar, konto_id)
	                   VALUES (?, ?, ?, ?, ?, ?, `+dane.WskazanieKonta+`)`,
		h.Id, idOkna, string(kierunek), zrodlo, cel, bajty, dane.KontoOperatora(kontekst))
	if err != nil {
		return fmt.Errorf("zdalne: plik przeniesiony, ale zapis prowenancji przenosin "+
			"nie powiódł się: %w", err)
	}
	return nil
}
