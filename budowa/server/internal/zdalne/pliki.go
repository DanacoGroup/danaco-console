// Odpowiedzialność pliku: przenosiny plików torem zdalnym. Nośnikiem jest `scp`
// po tym samym SSH i pod tą samą zgodą per host, co tor procesu; każdy wykonany
// ruch bajtów zostawia wiersz prowenancji w `zdalne_przeniesienie`.
//
// Nagrania dźwięku nie jadą tym torem: dźwięk nie opuszcza maszyny Operatora
// (klient pilnuje tego w `dostarczenie-nagrania.ts`), a potrzeby też nie ma —
// silnik mowy jest usługą rdzenia i bierze ścieżkę na maszynie silnika
// (mowa/nagranie.go), więc transkrypcja domyka się przed torem, a do procesu
// zdalnego jedzie wyłącznie tekst. Wykaz rozszerzeń nagrań jest jeden,
// `mowa.FormatyNagran` — drugiej listy ten plik nie zakłada.
//
// Konsument. W rdzeniu funkcję woła spoina katalogów roboczych przy zasięgu
// `remote`.
package zdalne

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"danacoconsole/server/internal/mowa"
)

// Kierunek ruchu bajtów względem maszyny rdzenia — wartości kolumny
// `zdalne_przeniesienie.kierunek` co do znaku.
type Kierunek string

const (
	// Wyslanie — plik z maszyny rdzenia na host zdalny.
	Wyslanie Kierunek = "wyslanie"
	// Pobranie — plik z hosta zdalnego na maszynę rdzenia.
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
	host, err := hostZRejestru(nazwaHosta)
	if err != nil {
		return err
	}
	sciezkaSCP, err := exec.LookPath("scp")
	if err != nil {
		return fmt.Errorf("zdalne: plik %s nie został przeniesiony, bo maszyna rdzenia "+
			"nie ma programu scp w PATH — przenosiny jadą wyłącznie po SSH i wymagają "+
			"jego instalacji: %w", sciezkaLokalna, err)
	}

	zrodlo, cel := sciezkaLokalna, host.adresSCP(sciezkaZdalna)
	if kierunek == Pobranie {
		zrodlo, cel = host.adresSCP(sciezkaZdalna), sciezkaLokalna
	}
	argumenty := append(append([]string{"-q"}, opcjeSSH...), "-P", strconv.Itoa(port(host)), zrodlo, cel)
	if wyjscie, err := exec.CommandContext(kontekst, sciezkaSCP, argumenty...).CombinedOutput(); err != nil {
		return fmt.Errorf("zdalne: przeniesienie %s → %s nie powiodło się: %w (%s)",
			zrodlo, cel, err, strings.TrimSpace(string(wyjscie)))
	}
	return odnotujPrzeniesienie(host, idOkna, kierunek, zrodlo, cel, rozmiar(sciezkaLokalna))
}

// odmowDzwieku egzekwuje regułę o nagraniach — po obu końcach
// drogi, bo pobranie nagrania Z hosta byłoby tym samym ruchem w drugą stronę.
func odmowDzwieku(sciezki ...string) error {
	for _, sciezka := range sciezki {
		if mowa.FormatPrzyjmowany(sciezka) {
			return fmt.Errorf("zdalne: plik %s nie został przeniesiony, bo jest nagraniem "+
				"dźwięku (wykaz formatów mowa.FormatyNagran), a dźwięk nie opuszcza maszyny "+
				"Operatora — rozstrzygnięcie Właściciela; transkrypcję wykonuje silnik mowy "+
				"na maszynie rdzenia i do hosta zdalnego jedzie wyłącznie tekst", sciezka)
		}
	}
	return nil
}

// adresSCP składa zdalny koniec drogi w postaci [użytkownik@]adres:ścieżka.
func (h Host) adresSCP(sciezka string) string {
	przod := h.AdresPolaczenia()
	if uzytkownik := strings.TrimSpace(h.Uzytkownik); uzytkownik != "" {
		przod = uzytkownik + "@" + przod
	}
	return przod + ":" + sciezka
}

// port zwraca port hosta; zero schodzi na 22 — domyślny port SSH.
func port(h Host) int {
	if h.Port == 0 {
		return 22
	}
	return h.Port
}

// rozmiar mierzy plik lokalny; brak pomiaru daje zero, nie odmowę.
func rozmiar(sciezka string) int64 {
	opis, err := os.Stat(sciezka)
	if err != nil {
		return 0
	}
	return opis.Size()
}

// odnotujPrzeniesienie zapisuje wiersz prowenancji przenosin. Zapis następuje
// po ruchu bajtów i niczego nie steruje — to dziennik, nie kolejka.
func odnotujPrzeniesienie(h Host, idOkna string, kierunek Kierunek, zrodlo, cel string, bajty int64) error {
	db := baza()
	if db == nil {
		return odmowaBrakuZasilenia()
	}
	_, err := db.Exec(`INSERT INTO zdalne_przeniesienie
	                   (host_id, okno_id, kierunek, sciezka_zrodlowa, sciezka_docelowa, rozmiar)
	                   VALUES (?, ?, ?, ?, ?, ?)`,
		h.Id, idOkna, string(kierunek), zrodlo, cel, bajty)
	if err != nil {
		return fmt.Errorf("zdalne: plik przeniesiony, ale zapis prowenancji przenosin "+
			"nie powiódł się: %w", err)
	}
	return nil
}
