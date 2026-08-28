// Plik pomiaru materiału obsługuje media.inspect i przekłada odpowiedź ffprobe w zapisie strukturalnym na pola kontraktu.
package core

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"danacoconsole/shared"
)

// opisMediow jest kopertą odpowiedzi `ffprobe` — dokładnie tymi polami, których
// używamy. Reszty odpowiedzi nie odczytujemy: pola, którego nikt nie czyta, nie
// trzeba utrzymywać przy zmianie wydania programu.
type opisMediow struct {
	Format kontenerMediow `json:"format"`
	// Strumienie zostają surowe, bo kontrakt każe oddać ich opis w zapisie tekstowym w całości.
	Strumienie json.RawMessage `json:"streams"`
}

// kontenerMediow niesie pola sekcji `format`. Czas trwania i rozmiar przychodzą
// z `ffprobe` jako tekst (sekundy z ułamkiem, bajty) i tak je bierzemy, zamiast
// wymuszać typ, którego program nie obiecuje.
type kontenerMediow struct {
	Nazwa    string `json:"format_name"`
	Trwanie  string `json:"duration"`
	Rozmiar  string `json:"size"`
	Strumien int    `json:"nb_streams"`
}

// strumienMediow niesie te cechy pojedynczego strumienia, po których rodzina
// rozstrzyga: rodzaj (obraz czy dźwięk), kodek (kontener docelowy przy
// wyodrębnianiu dźwięku) i wymiary (zasób wynikowy).
type strumienMediow struct {
	Rodzaj    string `json:"codec_type"`
	Kodek     string `json:"codec_name"`
	Szerokosc int    `json:"width"`
	Wysokosc  int    `json:"height"`
}

// Zbadaj obsługuje media.inspect: mierzy materiał tą samą drogą źródła co przy przetwarzaniu, żeby opis dotyczył treści zapisanej w magazynie.
func (a *adapterNarzedziMediow) Zbadaj(ctx context.Context,
	z shared.MediaInspectRequest) (shared.MediaInspectResponse, error) {

	// Komenda media.inspect niczego nie wytwarza, więc okno źródła zostaje tu porzucone.
	sciezka, _, err := a.zrodloMediow(ctx, "media.inspect", z.AssetId, z.SourcePath)
	if err != nil {
		return shared.MediaInspectResponse{}, err
	}

	opis, surowe, err := a.zmierzMediow(ctx, sciezka)
	if err != nil {
		return shared.MediaInspectResponse{}, err
	}

	// Rozmiar bierze się z odczytu systemowego pliku; pole rozmiaru odpowiedzi jest zapasem awaryjnym.
	rozmiar := rozmiarPlikuMediow(sciezka)
	if rozmiar == 0 {
		rozmiar = liczbaCalkowitaMediow(opis.Format.Rozmiar)
	}

	return shared.MediaInspectResponse{
		DurationMs: trwanieMediowMs(opis.Format.Trwanie),
		Streams:    surowe,
		Format:     opis.Format.Nazwa,
		SizeBytes:  rozmiar,
	}, nil
}

// zmierzMediow woła ffprobe i rozbiera jego odpowiedź, oddając przy okazji surowy opis strumieni gotowy do pola streams kontraktu.
func (a *adapterNarzedziMediow) zmierzMediow(ctx context.Context,
	sciezka string) (opisMediow, string, error) {

	wynik, err := a.wolajMediow(ctx, narzedzieFfprobe, []string{
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		sciezka,
	}, granicaPomiaruMediow)
	if err != nil {
		return opisMediow{}, "", bladNarzedziMediow(err)
	}

	var opis opisMediow
	if err := json.Unmarshal(wynik.Wyjscie, &opis); err != nil {
		return opisMediow{}, "", odmowaNarzedziMediow(shared.ErrorCodeInternalError,
			"ffprobe oddał odpowiedź, której nie da się odczytać: "+err.Error()+
				"; naprawa: sprawdzić, czy wskazany plik jest materiałem dźwiękowym "+
				"albo filmowym")
	}

	surowe := strings.TrimSpace(string(opis.Strumienie))
	if surowe == "" || surowe == "null" {
		// Plik bez strumienia nie jest materiałem: pusty opis wyglądałby jak udany pomiar.
		return opisMediow{}, "", odmowaNarzedziMediow(shared.ErrorCodeValidationFailed,
			"wskazany plik nie niesie żadnego strumienia dźwięku ani obrazu — "+
				"nie ma czego zmierzyć")
	}
	return opis, surowe, nil
}

// strumienieMediow rozbiera surowy opis strumieni na cechy, po których rodzina rozstrzyga rodzaj, kodek i wymiary materiału.
func strumienieMediow(surowe string) []strumienMediow {
	var strumienie []strumienMediow
	if err := json.Unmarshal([]byte(surowe), &strumienie); err != nil {
		return nil
	}
	return strumienie
}

// wymiaryMediow wyjmuje wymiary pierwszego strumienia, który je podaje.
//
// Brak takiego strumienia (materiał czysto dźwiękowy) oddaje dwie pustki, bo
// zasób bez wymiarów ma mieć `null` w kolumnie; zero znaczyłoby zmierzone zero
// pikseli.
func wymiaryMediow(strumienie []strumienMediow) (*int, *int) {
	for _, strumien := range strumienie {
		if strumien.Szerokosc > 0 && strumien.Wysokosc > 0 {
			szerokosc, wysokosc := strumien.Szerokosc, strumien.Wysokosc
			return &szerokosc, &wysokosc
		}
	}
	return nil, nil
}

// kodekDzwiekuMediow oddaje nazwę kodeka pierwszego strumienia dźwięku.
// Pustka znaczy „materiał nie niesie dźwięku" i wołający robi z niej odmowę.
func kodekDzwiekuMediow(strumienie []strumienMediow) string {
	for _, strumien := range strumienie {
		if strumien.Rodzaj == "audio" {
			return strings.TrimSpace(strumien.Kodek)
		}
	}
	return ""
}

// trwanieMediowMs przekłada czas trwania z sekund z ułamkiem na milisekundy
// kontraktu. Wartość nieczytelna albo niedodatnia daje zero, czyli brak wartości.
func trwanieMediowMs(sekundy string) int {
	wartosc, err := strconv.ParseFloat(strings.TrimSpace(sekundy), 64)
	if err != nil || wartosc <= 0 {
		return 0
	}
	return int(wartosc * 1000)
}

// liczbaCalkowitaMediow czyta liczbę zapisaną tekstem. Zero przy każdym
// niepowodzeniu, tak samo jak przy czasie trwania.
func liczbaCalkowitaMediow(tekst string) int {
	wartosc, err := strconv.Atoi(strings.TrimSpace(tekst))
	if err != nil || wartosc < 0 {
		return 0
	}
	return wartosc
}
