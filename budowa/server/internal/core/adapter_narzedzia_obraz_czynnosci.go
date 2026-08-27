// Czynności narzędzi obrazu image.inspect, image.transform, image.adjust
// i image.convert. Każda liczy się w procesie biblioteką wkompilowaną,
// a przy formacie bez rachunku Go korzysta z programu pakietu serwera jako
// drogi zapasowej.
package core

import (
	"context"
	"errors"
	"image"
	"strconv"
	"strings"

	"danacoconsole/shared"
)

// domyslneSilyPoprawki zbiera wartości `amount` brane przy braku wskazania.
// Wartości są procentami w rozumieniu tej rodziny narzędzi — przekład na
// parametry rachunku robią `poprawObrazArsenalu` i `argumentyPoprawki`.
var domyslneSilyPoprawki = map[shared.ImageAdjustKind]int{
	shared.ImageAdjustKindBrightness: 10,
	shared.ImageAdjustKindContrast:   10,
	shared.ImageAdjustKindSaturation: 20,
	shared.ImageAdjustKindSharpen:    50,
	shared.ImageAdjustKindBlur:       50,
	shared.ImageAdjustKindDenoise:    50,
}

// formatyDocelowe to zamknięty zbiór formatów, które `image.convert` przyjmuje.
// Nazwa formatu trafia w argument programu drogi zapasowej, więc dowolny tekst
// dawałby modelowi wpływ szerszy niż zamiana formatu.
var formatyDocelowe = map[string]struct{}{
	"png": {}, "jpeg": {}, "jpg": {}, "webp": {}, "avif": {}, "tiff": {}, "gif": {},
}

// czynnoscObrazu opisuje jedno żądanie zmiany w obu drogach naraz: `rachunek`
// liczy je w procesie, `argumenty` składają wiersz poleceń programu pakietu
// serwera dla formatów bez rachunku Go. Obie muszą znaczyć to samo.
type czynnoscObrazu struct {
	rachunek    func(image.Image) (image.Image, error)
	argumenty   []string
	format      string
	jakosc      *int
	bezstratnie *bool
	nazwa       string
	// kompresuj włącza dogniecenie zapisu programem; niesie je tylko
	// konwersja.
	kompresuj bool
}

// Zbadaj obsługuje `image.inspect` — czynność czytającą, którą model wywołuje
// przed każdą zmianą, żeby znać format i wymiary. Nagłówek pliku czyta
// biblioteka wkompilowana, bo pomiar to jeden odczyt kilkuset bajtów.
func (a *adapterNarzedziObrazu) Zbadaj(ctx context.Context,
	z shared.ImageInspectRequest) (shared.ImageInspectResponse, error) {

	zrodlo, err := a.rozwiazZrodlo(ctx, z.AssetId, z.SourcePath)
	if err != nil {
		return shared.ImageInspectResponse{}, err
	}

	opis, err := zbadajObrazArsenalu(zrodlo.sciezka)
	if err != nil {
		if !errors.Is(err, errBrakRachunkuGoObrazu) {
			return shared.ImageInspectResponse{}, bladPrzetwarzaniaObrazu(err.Error())
		}
		return a.zbadajProgramem(ctx, zrodlo)
	}

	odpowiedz := shared.ImageInspectResponse{
		Format:    opis.format,
		Width:     opis.szerokosc,
		Height:    opis.wysokosc,
		SizeBytes: zrodlo.rozmiar,
	}
	// Pola opcjonalne obsadza się tylko wtedy, gdy pomiar coś dał — pusty tekst
	// udawałby metadane.
	if opis.przestrzen != "" {
		przestrzen := opis.przestrzen
		odpowiedz.ColorSpace = &przestrzen
	}
	if opis.metadane != "" {
		metadane := opis.metadane
		odpowiedz.Metadata = &metadane
	}
	return odpowiedz, nil
}

// zbadajProgramem mierzy programem pakietu serwera plik bez dekodera w Go,
// na przykład AVIF. Indeks `[0]` bierze pierwszą klatkę, inaczej GIF animowany
// oddałby wymiary swojej ostatniej klatki.
func (a *adapterNarzedziObrazu) zbadajProgramem(ctx context.Context,
	zrodlo zrodloObrazu) (shared.ImageInspectResponse, error) {

	argumenty := []string{
		"-format", "%m\n%w\n%h\n%[colorspace]\n%[EXIF:*]",
		zrodlo.sciezka + "[0]",
	}
	wyjscie, err := a.wolajImageMagick(ctx, "identify", argumenty, "identify")
	if err != nil {
		return shared.ImageInspectResponse{}, err
	}

	wiersze := strings.Split(string(wyjscie), "\n")
	if len(wiersze) < 4 {
		return shared.ImageInspectResponse{}, bladPrzetwarzaniaObrazu(
			"program pakietu serwera nie opisał obrazu w oczekiwanym kształcie: " +
				strings.TrimSpace(string(wyjscie)))
	}
	odpowiedz := shared.ImageInspectResponse{
		Format:    strings.ToLower(strings.TrimSpace(wiersze[0])),
		Width:     liczbaZOpisu(wiersze[1]),
		Height:    liczbaZOpisu(wiersze[2]),
		SizeBytes: zrodlo.rozmiar,
	}
	if przestrzen := strings.TrimSpace(wiersze[3]); przestrzen != "" {
		odpowiedz.ColorSpace = &przestrzen
	}
	if metadane := strings.TrimSpace(strings.Join(wiersze[4:], "\n")); metadane != "" {
		odpowiedz.Metadata = &metadane
	}
	return odpowiedz, nil
}

// Przeksztalc obsługuje `image.transform` — geometrię obrazu. Składa argumenty
// dla obu dróg i deleguje rachunek do `przetworz`, wspólnego dla wszystkich
// czynności zmieniających obraz.
func (a *adapterNarzedziObrazu) Przeksztalc(ctx context.Context,
	z shared.ImageTransformRequest) (shared.ImageTransformResponse, error) {

	zrodlo, err := a.rozwiazZrodlo(ctx, z.AssetId, z.SourcePath)
	if err != nil {
		return shared.ImageTransformResponse{}, err
	}
	argumenty, err := argumentyPrzeksztalcenia(z)
	if err != nil {
		return shared.ImageTransformResponse{}, err
	}
	zasob, _, err := a.przetworz(ctx, zrodlo, z.WindowId, czynnoscObrazu{
		rachunek: func(obraz image.Image) (image.Image, error) {
			return przeksztalcObrazArsenalu(obraz, z)
		},
		argumenty: argumenty,
		nazwa:     string(z.Operation),
	})
	if err != nil {
		return shared.ImageTransformResponse{}, err
	}
	return shared.ImageTransformResponse{Asset: zasob}, nil
}

// Popraw obsługuje `image.adjust` — retusz barw i ostrości obrazu. Składa
// argumenty dla obu dróg i deleguje rachunek do `przetworz`, tak samo
// jak Przeksztalc.
func (a *adapterNarzedziObrazu) Popraw(ctx context.Context,
	z shared.ImageAdjustRequest) (shared.ImageAdjustResponse, error) {

	zrodlo, err := a.rozwiazZrodlo(ctx, z.AssetId, z.SourcePath)
	if err != nil {
		return shared.ImageAdjustResponse{}, err
	}
	argumenty, err := argumentyPoprawki(z)
	if err != nil {
		return shared.ImageAdjustResponse{}, err
	}
	zasob, _, err := a.przetworz(ctx, zrodlo, z.WindowId, czynnoscObrazu{
		rachunek: func(obraz image.Image) (image.Image, error) {
			return poprawObrazArsenalu(obraz, z)
		},
		argumenty: argumenty,
		nazwa:     string(z.Operation),
	})
	if err != nil {
		return shared.ImageAdjustResponse{}, err
	}
	return shared.ImageAdjustResponse{Asset: zasob}, nil
}

// Przekonwertuj obsługuje `image.convert` — zmianę formatu i kompresji. Pole
// `savedBytes` wypełnia się tylko przy realnej oszczędności, bo wartość ujemna
// nazywałaby stratę oszczędnością; `sizeBytes` obok zawsze podaje rozmiar
// wyniku.
func (a *adapterNarzedziObrazu) Przekonwertuj(ctx context.Context,
	z shared.ImageConvertRequest) (shared.ImageConvertResponse, error) {

	zrodlo, err := a.rozwiazZrodlo(ctx, z.AssetId, z.SourcePath)
	if err != nil {
		return shared.ImageConvertResponse{}, err
	}
	format, argumenty, err := argumentyKonwersji(z)
	if err != nil {
		return shared.ImageConvertResponse{}, err
	}
	zasob, rozmiar, err := a.przetworz(ctx, zrodlo, z.WindowId, czynnoscObrazu{
		// Konwersja nie rusza pikseli — rachunkiem jest przepisanie obrazu
		// do kodera formatu docelowego.
		rachunek:    func(obraz image.Image) (image.Image, error) { return obraz, nil },
		argumenty:   argumenty,
		format:      format,
		jakosc:      z.Quality,
		bezstratnie: z.Lossless,
		nazwa:       "convert:" + format,
		kompresuj:   true,
	})
	if err != nil {
		return shared.ImageConvertResponse{}, err
	}
	odpowiedz := shared.ImageConvertResponse{Asset: zasob, SizeBytes: rozmiar}
	if oszczednosc := zrodlo.rozmiar - rozmiar; oszczednosc > 0 {
		odpowiedz.SavedBytes = &oszczednosc
	}
	return odpowiedz, nil
}

// przetworz jest wspólną drogą trzech czynności zmieniających obraz: liczy go
// w procesie, a przy formacie bez rachunku Go korzysta z programu pakietu
// serwera, po czym odkłada wynik jako nowy zasób.
func (a *adapterNarzedziObrazu) przetworz(ctx context.Context, zrodlo zrodloObrazu,
	oknoZadane *string, czynnosc czynnoscObrazu) (shared.DesignAsset, int, error) {

	format := strings.TrimSpace(czynnosc.format)
	if format == "" {
		format = formatZrodla(zrodlo.sciezka)
	}

	bajty, err := a.policzWkompilowanym(zrodlo, format, czynnosc)
	if err != nil {
		if !errors.Is(err, errBrakRachunkuGoObrazu) {
			return shared.DesignAsset{}, 0, err
		}
		bajty, err = a.policzProgramem(ctx, zrodlo, format, czynnosc)
		if err != nil {
			return shared.DesignAsset{}, 0, err
		}
	}
	if czynnosc.kompresuj {
		// Krok bez odmowy: przy braku programu dogniatającego oddaje te same
		// bajty bez zmiany.
		bajty = a.dogniecZapisObrazu(ctx, bajty, format, czynnosc.bezstratnie)
	}
	return a.odlozZasob(ctx, zrodlo, oknoZadane, bajty, format, czynnosc.nazwa)
}

// policzWkompilowanym przeprowadza czynność biblioteką w procesie. Odmowa
// rachunku wraca wprost i kończy czynność, ponieważ droga zapasowa nie potrafi
// jej naprawić.
func (a *adapterNarzedziObrazu) policzWkompilowanym(zrodlo zrodloObrazu, format string,
	czynnosc czynnoscObrazu) ([]byte, error) {

	obraz, err := odczytajObrazArsenalu(zrodlo.sciezka)
	if err != nil {
		return nil, err
	}
	wynik, err := czynnosc.rachunek(obraz)
	if err != nil {
		return nil, err
	}
	return zakodujObrazArsenalu(wynik, format, czynnosc.jakosc, czynnosc.bezstratnie)
}

// policzProgramem przeprowadza czynność programem pakietu serwera — drogą dla
// AVIF-a i WEBP-a stratnego. Kolejność argumentów jest wiążąca: ścieżka źródła
// stoi przed operatorami, a cel `format:-` na końcu.
func (a *adapterNarzedziObrazu) policzProgramem(ctx context.Context, zrodlo zrodloObrazu,
	format string, czynnosc czynnoscObrazu) ([]byte, error) {

	argumenty := append([]string{zrodlo.sciezka}, czynnosc.argumenty...)
	argumenty = append(argumenty, format+":-")
	return a.wolajImageMagick(ctx, "", argumenty, "convert")
}

// argumentyPrzeksztalcenia składa operatory geometryczne dla drogi zapasowej.
// Każda gałąź sprawdza własne wskazania, bo kadr bez wymiarów albo obrót bez
// kąta to żądanie niewykonalne, a nie żądanie o skutku pustym.
func argumentyPrzeksztalcenia(z shared.ImageTransformRequest) ([]string, error) {
	switch z.Operation {
	case shared.ImageTransformKindResize:
		miara, err := miaraSkalowania(z.Width, z.Height, z.KeepAspect)
		if err != nil {
			return nil, err
		}
		return []string{"-resize", miara}, nil

	case shared.ImageTransformKindThumbnail:
		// `-thumbnail` zdejmuje profil i metadane, inaczej niż `-resize`.
		// Brak wymiarów bierze bok 256.
		szerokosc, wysokosc := z.Width, z.Height
		if szerokosc == nil && wysokosc == nil {
			bok := 256
			szerokosc = &bok
		}
		miara, err := miaraSkalowania(szerokosc, wysokosc, z.KeepAspect)
		if err != nil {
			return nil, err
		}
		return []string{"-thumbnail", miara}, nil

	case shared.ImageTransformKindCrop:
		if !dodatnia(z.Width) || !dodatnia(z.Height) {
			return nil, bladWskazaniaObrazu(
				"kadrowanie wymaga pól width i height dodatnich — kadr bez wymiarów nie istnieje")
		}
		miara := strconv.Itoa(*z.Width) + "x" + strconv.Itoa(*z.Height) +
			przesuniecie(z.X) + przesuniecie(z.Y)
		// `+repage` kasuje ślad po pierwotnym płótnie, inaczej obrót po kadrze
		// wyszedłby przesunięty.
		return []string{"-crop", miara, "+repage"}, nil

	case shared.ImageTransformKindRotate:
		if z.Degrees == nil {
			return nil, bladWskazaniaObrazu("obrót wymaga pola degrees — obrót o nieznany kąt nie istnieje")
		}
		return []string{"-rotate", strconv.Itoa(*z.Degrees)}, nil

	case shared.ImageTransformKindFlipHorizontal:
		// `-flop` odbija w poziomie, `-flip` w pionie; nazwy łatwo pomylić.
		return []string{"-flop"}, nil

	case shared.ImageTransformKindFlipVertical:
		return []string{"-flip"}, nil
	}
	return nil, bladWskazaniaObrazu("nieznane przekształcenie " + string(z.Operation) +
		"; kontrakt zna: resize, crop, rotate, flipHorizontal, flipVertical, thumbnail")
}

// argumentyPoprawki składa operatory retuszu dla drogi zapasowej wraz
// z przekładem siły na parametr programu.
func argumentyPoprawki(z shared.ImageAdjustRequest) ([]string, error) {
	sila := domyslneSilyPoprawki[z.Operation]
	if z.Amount != nil {
		sila = *z.Amount
	}

	switch z.Operation {
	case shared.ImageAdjustKindBrightness:
		return []string{"-brightness-contrast", strconv.Itoa(sila) + "x0"}, nil
	case shared.ImageAdjustKindContrast:
		return []string{"-brightness-contrast", "0x" + strconv.Itoa(sila)}, nil
	case shared.ImageAdjustKindSaturation:
		// `-modulate` liczy w procentach, gdzie 100 znaczy bez zmian — siła
		// jest przyrostem, nie wartością.
		return []string{"-modulate", "100," + strconv.Itoa(100+sila) + ",100"}, nil
	case shared.ImageAdjustKindSharpen:
		return []string{"-sharpen", "0x" + promienZSily(sila)}, nil
	case shared.ImageAdjustKindBlur:
		return []string{"-blur", "0x" + promienZSily(sila)}, nil
	case shared.ImageAdjustKindDenoise:
		// `-median` bierze promień sąsiedztwa — ten sam filtr, co rachunek
		// wkompilowany.
		operatory := make([]string, 0, 2*przebiegiOdszumianiaArsenalu(sila))
		for i := 0; i < przebiegiOdszumianiaArsenalu(sila); i++ {
			operatory = append(operatory, "-median", "1")
		}
		return operatory, nil
	case shared.ImageAdjustKindGrayscale:
		return []string{"-colorspace", "Gray"}, nil
	case shared.ImageAdjustKindAutoLevels:
		return []string{"-auto-level"}, nil
	}
	return nil, bladWskazaniaObrazu("nieznana poprawka " + string(z.Operation) +
		"; kontrakt zna: brightness, contrast, saturation, sharpen, blur, denoise, " +
		"grayscale, autoLevels")
}

// argumentyKonwersji rozstrzyga format docelowy i parametry kompresji.
// Sprawdza, że format należy do zbioru `formatyDocelowe`, i odrzuca żądanie
// bezstratności, którego wskazany format nie potrafi spełnić.
func argumentyKonwersji(z shared.ImageConvertRequest) (string, []string, error) {
	format := strings.ToLower(strings.TrimSpace(z.Format))
	if format == "" {
		return "", nil, bladWskazaniaObrazu("konwersja bez pola format — nie wiadomo, na co zamieniać")
	}
	if _, znany := formatyDocelowe[format]; !znany {
		return "", nil, bladWskazaniaObrazu("format " + format +
			" nie jest formatem docelowym tej czynności; " +
			"kontrakt zna: png, jpeg, webp, avif, tiff, gif")
	}

	operatory := []string{}
	if z.Quality != nil {
		if *z.Quality < 1 || *z.Quality > 100 {
			return "", nil, bladWskazaniaObrazu("jakość " + strconv.Itoa(*z.Quality) +
				" jest poza zakresem 1-100")
		}
		operatory = append(operatory, "-quality", strconv.Itoa(*z.Quality))
	}
	if z.Lossless != nil && *z.Lossless {
		// Bezstratność ma sens tylko tam, gdzie format ją zna — na JPEG-u
		// byłaby niespełnialną obietnicą.
		switch format {
		case "webp", "avif":
			operatory = append(operatory, "-define", format+":lossless=true")
		case "png", "tiff", "gif":
			// Te formaty są bezstratne z natury — wskazanie jest spełnione bez
			// jednego argumentu więcej.
		default:
			return "", nil, bladWskazaniaObrazu("format " + format +
				" nie zna kompresji bezstratnej; naprawa: wybrać webp, avif, png albo tiff")
		}
	}
	return format, operatory, nil
}

// miaraSkalowania składa miarę skalowania programu z pary wymiarów, w tym
// samym rozumieniu, co `przeskalujObrazArsenalu`. Proporcje zachowuje
// domyślnie; wyłączenie ich wymaga obu wymiarów, bo rozciągnięcie do jednego
// boku nie ma czego rozciągnąć.
func miaraSkalowania(szerokosc, wysokosc *int, zachowajProporcje *bool) (string, error) {
	if !dodatnia(szerokosc) && !dodatnia(wysokosc) {
		return "", bladWskazaniaObrazu(
			"skalowanie wymaga pola width albo height dodatniego — rozmiar docelowy nie jest znany")
	}
	proporcje := zachowajProporcje == nil || *zachowajProporcje
	if !proporcje {
		if !dodatnia(szerokosc) || !dodatnia(wysokosc) {
			return "", bladWskazaniaObrazu(
				"skalowanie z keepAspect=false wymaga OBU wymiarów — bez nich nie ma czego rozciągnąć")
		}
		return strconv.Itoa(*szerokosc) + "x" + strconv.Itoa(*wysokosc) + "!", nil
	}
	if dodatnia(szerokosc) && dodatnia(wysokosc) {
		return strconv.Itoa(*szerokosc) + "x" + strconv.Itoa(*wysokosc), nil
	}
	if dodatnia(szerokosc) {
		return strconv.Itoa(*szerokosc), nil
	}
	return "x" + strconv.Itoa(*wysokosc), nil
}

// przesuniecie zapisuje odsunięcie kadru w składni programu: znak jest zawsze
// obecny, brak wskazania znaczy zero (kadr od lewego górnego rogu).
func przesuniecie(wartosc *int) string {
	liczba := 0
	if wartosc != nil {
		liczba = *wartosc
	}
	if liczba < 0 {
		return strconv.Itoa(liczba)
	}
	return "+" + strconv.Itoa(liczba)
}

// promienZSily przekłada siłę w procentach na sigmę rozmycia albo wyostrzenia —
// tym samym rachunkiem, co `sigmaZSilyArsenalu`.
func promienZSily(sila int) string {
	return strconv.FormatFloat(sigmaZSilyArsenalu(sila), 'f', 2, 64)
}

// dodatnia orzeka, czy wskazany wymiar jest wymiarem. Zero i liczby ujemne są
// tu brakiem, nie wartością: obraz o szerokości zero nie istnieje.
func dodatnia(wartosc *int) bool {
	return wartosc != nil && *wartosc > 0
}

// liczbaZOpisu czyta liczbę z wiersza wypisanego przez program. Wiersz
// nieczytelny daje zero — a zero w polu `width` odpowiedzi jest widoczne
// natychmiast i nie udaje zmierzonego wymiaru.
func liczbaZOpisu(wiersz string) int {
	liczba, err := strconv.Atoi(strings.TrimSpace(wiersz))
	if err != nil {
		return 0
	}
	return liczba
}

// formatZrodla rozpoznaje format pliku z jego nagłówka, tą samą drogą, którą
// robi to wniesienie zasobu (`rozpoznajObrazZasobu`). Format nierozpoznany
// oddaje `png`, tak jak opisuje to `przetworz`.
func formatZrodla(sciezka string) string {
	format, _, _ := rozpoznajObrazZasobu(sciezka)
	if format != nil && strings.TrimSpace(*format) != "" {
		return strings.ToLower(strings.TrimSpace(*format))
	}
	return "png"
}
