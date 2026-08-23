// Odpowiedzialność pliku: rachunek WKOMPILOWANY czterech czynności obrazu
// modelu — `image.inspect`, `image.transform`, `image.adjust`, `image.convert`.
// Czynności i ich odmowy stoją w `adapter_narzedzia_obraz_czynnosci.go`, wspólne
// zaplecze (źródło, magazyn wyniku) w `adapter_narzedzia_obraz.go`.
//
// ── Dlaczego ten plik powstał ───────────────────────────────────────────────
// Te cztery czynności liczył wcześniej program zewnętrzny, choć każdą z nich
// wykonuje w całości biblioteka Go wkompilowana w binarium. Program zewnętrzny
// wołany tam, gdzie biblioteka wystarcza, jest regresem: kosztuje uruchomienie
// procesu, wiąże funkcję z wersją cudzego wydania i przy niekompletnym serwerze
// zamienia retusz w odmowę. Rachunek stoi więc tutaj i idzie w procesie:
// `disintegration/imaging` (skalowanie Lanczosem, kadr, obrót, odbicia,
// korekcje barwne, rozmycie, wyostrzenie), `golang.org/x/image` (dekodery WEBP,
// TIFF, BMP oraz kodery TIFF i BMP), `HugoSmits86/nativewebp` (zapis WEBP
// bezstratnego), `rwcarlsen/goexif` (odczyt metadanych EXIF) i rachunek własny
// na odszumianie medianą oraz rozciągnięcie poziomów.
//
// ── Gdzie rachunku Go NIE MA — i co się wtedy dzieje ────────────────────────
// Dwa wyjścia nie mają w Go kodera i nie da się ich tu policzyć: **AVIF**
// (kodera czysto-Go nie ma wcale) oraz **WEBP stratny** (`nativewebp` zapisuje
// wyłącznie bezstratny VP8L). Tak samo AVIF nie ma dekodera, więc obraz w tym
// formacie nie wchodzi. Te przypadki oddają `errBrakRachunkuGoObrazu`, a
// czynność sięga wtedy po program pakietu serwera — jedyna droga, jaka zostaje,
// i lepsza od odmowy, bo odmowa jest brakiem funkcji. Zapora
// `zapora_narzedzi_obrazu_test.go` pilnuje, żeby ta droga została wyjątkiem
// nazwanym, a nie wróciła jako droga podstawowa.
package core

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"os"
	"sort"
	"strings"

	"github.com/HugoSmits86/nativewebp"
	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"

	"danacoconsole/shared"
)

// errBrakRachunkuGoObrazu znaczy „tego nie policzy tu żadna biblioteka
// wkompilowana" — a nie „czynność się nie udała". Wołający rozpoznaje ten błąd
// i przechodzi na program pakietu serwera; każdy inny błąd jest odmową wprost,
// bo obraz uszkodzony ma zostać nazwany, a nie oddany drugiej drodze.
var errBrakRachunkuGoObrazu = errors.New("rachunek wkompilowany nie zna tego formatu")

// tloObrotuArsenalu wypełnia narożniki powstałe przy obrocie o kąt niebędący
// wielokrotnością prostego. Biel, nie przezroczystość: obrót zapisany potem
// w JPEG-u nie ma czym nieść kanału alfa i narożniki wyszłyby czarne.
var tloObrotuArsenalu = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}

// opisObrazuArsenalu jest pomiarem pliku dla `image.inspect`: format, wymiary,
// przestrzeń barw i metadane. Puste pole znaczy „nie zmierzono", nigdy „zero".
type opisObrazuArsenalu struct {
	format     string
	szerokosc  int
	wysokosc   int
	przestrzen string
	metadane   string
}

// zbadajObrazArsenalu czyta sam nagłówek pliku — bez rozkodowania wszystkich
// pikseli, bo do formatu i wymiarów nie są potrzebne.
//
// Metadane EXIF czytamy osobno i miękko: ich brak jest zwykłym stanem pliku
// (PNG z ekranu nie ma EXIF-u), więc nie przerywa pomiaru.
func zbadajObrazArsenalu(sciezka string) (opisObrazuArsenalu, error) {
	plik, err := os.Open(sciezka)
	if err != nil {
		return opisObrazuArsenalu{}, err
	}
	defer plik.Close()

	nastawy, format, err := image.DecodeConfig(plik)
	if err != nil {
		// Format bez dekodera w drzewie — pomiar zostaje dla programu pakietu
		// serwera, zamiast oddać wymiary zgadnięte.
		return opisObrazuArsenalu{}, errBrakRachunkuGoObrazu
	}

	opis := opisObrazuArsenalu{
		format:     strings.ToLower(strings.TrimSpace(format)),
		szerokosc:  nastawy.Width,
		wysokosc:   nastawy.Height,
		przestrzen: przestrzenBarwArsenalu(nastawy.ColorModel),
	}
	opis.metadane = metadaneExifArsenalu(sciezka)
	return opis, nil
}

// przestrzenBarwArsenalu nazywa przestrzeń barw modelem koloru, który oddał
// dekoder. Nazwy zostają te, które produkt wypisywał do tej pory (`sRGB`,
// `Gray`, `CMYK`), bo czyta je model w treści odpowiedzi — zmiana słownika
// byłaby zmianą kontraktu przy okazji zmiany rachunku.
//
// YCbCr jest zapisem JPEG-owym barw sRGB, nie osobną przestrzenią widzianą przez
// Operatora, więc wraca jako `sRGB`.
func przestrzenBarwArsenalu(model color.Model) string {
	switch model {
	case color.GrayModel, color.Gray16Model:
		return "Gray"
	case color.CMYKModel:
		return "CMYK"
	case color.YCbCrModel, color.NYCbCrAModel:
		return "sRGB"
	case color.RGBAModel, color.RGBA64Model, color.NRGBAModel, color.NRGBA64Model,
		color.AlphaModel, color.Alpha16Model:
		return "sRGB"
	}
	// Model nierozpoznany (paleta GIF-a, model własny dekodera) zostaje
	// nienazwany: puste pole mówi „nie wiem", a wpisana nazwa kłamałaby.
	return ""
}

// metadaneExifArsenalu składa metadane w wiersze „nazwa: wartość".
//
// Pusty tekst przy braku EXIF-u jest tu wynikiem poprawnym, nie usterką: pole
// `metadata` odpowiedzi zostaje wtedy nieobsadzone, tak jak przy pliku bez
// metadanych.
func metadaneExifArsenalu(sciezka string) string {
	plik, err := os.Open(sciezka)
	if err != nil {
		return ""
	}
	defer plik.Close()

	dane, err := exif.Decode(plik)
	if err != nil {
		return ""
	}
	zbieracz := &zbieraczExifArsenalu{}
	if err := dane.Walk(zbieracz); err != nil {
		return ""
	}
	// Porządek wypisania jest kolejnością obchodzenia znaczników — a ta idzie po
	// mapie, więc różni się między wywołaniami. Sortowanie daje odpowiedź
	// powtarzalną dla tego samego zdjęcia, a od niej zależy sprawdzian.
	sort.Strings(zbieracz.wiersze)
	return strings.Join(zbieracz.wiersze, "\n")
}

// zbieraczExifArsenalu wypełnia interfejs obchodzenia znaczników EXIF-u:
// biblioteka nie ma odpowiednika funkcyjnego, więc odbiorcą jest typ.
type zbieraczExifArsenalu struct {
	wiersze []string
}

// Walk zapisuje jeden znacznik jako wiersz „nazwa: wartość". Błędu nie oddaje
// nigdy: znacznik nieczytelny ma nie przerywać odczytu pozostałych, bo metadane
// są dodatkiem do pomiaru, a nie jego treścią.
func (z *zbieraczExifArsenalu) Walk(nazwa exif.FieldName, znacznik *tiff.Tag) error {
	z.wiersze = append(z.wiersze, string(nazwa)+": "+znacznik.String())
	return nil
}

// przeksztalcObrazArsenalu liczy geometrię `image.transform`.
//
// Filtr Lanczosa dla każdego skalowania: to on jest wyborem domyślnym warsztatu
// fotografii tego produktu, a dwa różne filtry w dwóch miejscach dawałyby dwa
// różne wyniki tej samej prośby modelu.
func przeksztalcObrazArsenalu(obraz image.Image,
	z shared.ImageTransformRequest) (image.Image, error) {

	switch z.Operation {
	case shared.ImageTransformKindResize:
		return przeskalujObrazArsenalu(obraz, z.Width, z.Height, z.KeepAspect)

	case shared.ImageTransformKindThumbnail:
		// Brak wymiarów bierze bok 256 — miniatura bez rozmiaru domyślnego nie
		// byłaby czynnością osobną od skalowania. Metadanych miniatura nie
		// niesie z natury tej drogi: rachunek składa nowy obraz z pikseli, więc
		// EXIF źródła nie ma czym przejść.
		szerokosc, wysokosc := z.Width, z.Height
		if szerokosc == nil && wysokosc == nil {
			bok := 256
			szerokosc = &bok
		}
		return przeskalujObrazArsenalu(obraz, szerokosc, wysokosc, z.KeepAspect)

	case shared.ImageTransformKindCrop:
		if !dodatnia(z.Width) || !dodatnia(z.Height) {
			return nil, bladWskazaniaObrazu(
				"kadrowanie wymaga pól width i height dodatnich — kadr bez wymiarów nie istnieje")
		}
		odsuniecieX, odsuniecieY := 0, 0
		if z.X != nil {
			odsuniecieX = *z.X
		}
		if z.Y != nil {
			odsuniecieY = *z.Y
		}
		granice := obraz.Bounds()
		kadr := image.Rect(
			granice.Min.X+odsuniecieX,
			granice.Min.Y+odsuniecieY,
			granice.Min.X+odsuniecieX+*z.Width,
			granice.Min.Y+odsuniecieY+*z.Height,
		)
		// Kadr poza obrazem jest odmową, nie obrazem pustym: prostokąt rozminięty
		// z powierzchnią oddałby zero pikseli opisanych jako skutek kadrowania.
		if kadr.Intersect(granice).Empty() {
			return nil, bladWskazaniaObrazu(
				"kadr nie ma części wspólnej z obrazem — poza jego powierzchnią nie ma czego wyciąć")
		}
		return imaging.Crop(obraz, kadr), nil

	case shared.ImageTransformKindRotate:
		if z.Degrees == nil {
			return nil, bladWskazaniaObrazu("obrót wymaga pola degrees — obrót o nieznany kąt nie istnieje")
		}
		// Znak przeciwny: kontrakt liczy kąt zgodnie z ruchem wskazówek zegara,
		// biblioteka — przeciwnie. Bez tej zamiany obrót o 90 stopni położyłby
		// zdjęcie na drugą stronę.
		return imaging.Rotate(obraz, -float64(*z.Degrees), tloObrotuArsenalu), nil

	case shared.ImageTransformKindFlipHorizontal:
		return imaging.FlipH(obraz), nil

	case shared.ImageTransformKindFlipVertical:
		return imaging.FlipV(obraz), nil
	}
	return nil, bladWskazaniaObrazu("nieznane przekształcenie " + string(z.Operation) +
		"; kontrakt zna: resize, crop, rotate, flipHorizontal, flipVertical, thumbnail")
}

// przeskalujObrazArsenalu rozstrzyga skalowanie z pary wymiarów.
//
// Proporcje zachowujemy domyślnie, tak jak mówi kontrakt: model prosi zwykle
// o „szerokość 800", a nie o rozciągnięcie zdjęcia. Obie miary podane przy
// zachowanych proporcjach znaczą „zmieść się w tej ramce" (`Fit`), a nie
// „rozciągnij do niej" — rozciągnięcie wymaga wskazania go wprost przez
// `keepAspect=false`.
func przeskalujObrazArsenalu(obraz image.Image, szerokosc, wysokosc *int,
	zachowajProporcje *bool) (image.Image, error) {

	if !dodatnia(szerokosc) && !dodatnia(wysokosc) {
		return nil, bladWskazaniaObrazu(
			"skalowanie wymaga pola width albo height dodatniego — rozmiar docelowy nie jest znany")
	}
	proporcje := zachowajProporcje == nil || *zachowajProporcje
	if !proporcje {
		if !dodatnia(szerokosc) || !dodatnia(wysokosc) {
			return nil, bladWskazaniaObrazu(
				"skalowanie z keepAspect=false wymaga OBU wymiarów — bez nich nie ma czego rozciągnąć")
		}
		return imaging.Resize(obraz, *szerokosc, *wysokosc, imaging.Lanczos), nil
	}
	if dodatnia(szerokosc) && dodatnia(wysokosc) {
		return imaging.Fit(obraz, *szerokosc, *wysokosc, imaging.Lanczos), nil
	}
	if dodatnia(szerokosc) {
		return imaging.Resize(obraz, *szerokosc, 0, imaging.Lanczos), nil
	}
	return imaging.Resize(obraz, 0, *wysokosc, imaging.Lanczos), nil
}

// poprawObrazArsenalu liczy retusz `image.adjust`.
//
// Siła jest procentem w rozumieniu tej rodziny narzędzi. Brak pola `amount`
// bierze wartość domyślną operacji, nie zero: zero byłoby poprawką bez skutku,
// a model prosząc „rozjaśnij" bez liczby dostałby obraz nieodróżnialny od
// źródła.
func poprawObrazArsenalu(obraz image.Image, z shared.ImageAdjustRequest) (image.Image, error) {
	sila := domyslneSilyPoprawki[z.Operation]
	if z.Amount != nil {
		sila = *z.Amount
	}

	switch z.Operation {
	case shared.ImageAdjustKindBrightness:
		return imaging.AdjustBrightness(obraz, float64(sila)), nil
	case shared.ImageAdjustKindContrast:
		return imaging.AdjustContrast(obraz, float64(sila)), nil
	case shared.ImageAdjustKindSaturation:
		return imaging.AdjustSaturation(obraz, float64(sila)), nil
	case shared.ImageAdjustKindSharpen:
		return imaging.Sharpen(obraz, sigmaZSilyArsenalu(sila)), nil
	case shared.ImageAdjustKindBlur:
		return imaging.Blur(obraz, sigmaZSilyArsenalu(sila)), nil
	case shared.ImageAdjustKindDenoise:
		return odszumMedianaArsenalu(obraz, przebiegiOdszumianiaArsenalu(sila)), nil
	case shared.ImageAdjustKindGrayscale:
		// Skala szarości nie ma stopni pośrednich w tej czynności — `amount` jest
		// tu bez znaczenia i nie udajemy, że coś z nim robimy.
		return imaging.Grayscale(obraz), nil
	case shared.ImageAdjustKindAutoLevels:
		return rozciagnijPoziomyArsenalu(obraz), nil
	}
	return nil, bladWskazaniaObrazu("nieznana poprawka " + string(z.Operation) +
		"; kontrakt zna: brightness, contrast, saturation, sharpen, blur, denoise, " +
		"grayscale, autoLevels")
}

// sigmaZSilyArsenalu przekłada siłę w procentach na sigmę rozmycia albo
// wyostrzenia. Pięćdziesiąt procent daje sigmę 1.0 — wartość, przy której skutek
// jest widoczny na ekranie, a obraz nie wygląda na uszkodzony. Zero i wartości
// ujemne dają najmniejszą sigmę o skutku widocznym, bo sigma zerowa nie
// zrobiłaby nic, a czynność ma robić to, o co poproszono.
func sigmaZSilyArsenalu(sila int) float64 {
	if sila <= 0 {
		return 0.1
	}
	return float64(sila) / 50.0
}

// przebiegiOdszumianiaArsenalu przekłada siłę na liczbę przebiegów mediany.
// Granica trzech jest ta sama, co przy poprzednim rachunku: czwarty przebieg
// zaciera szczegóły, które odszumienie ma zachować.
func przebiegiOdszumianiaArsenalu(sila int) int {
	przebiegi := sila / 50
	if przebiegi < 1 {
		return 1
	}
	if przebiegi > 3 {
		return 3
	}
	return przebiegi
}

// odszumMedianaArsenalu odszumia filtrem medianowym 3×3.
//
// Mediana, nie rozmycie: szum pojedynczych punktów (sól i pieprz z matrycy przy
// wysokiej czułości) jest wartością odstającą, a mediana odstających nie bierze,
// zamiast rozmazywać je po sąsiedztwie. Krawędzie zostają, bo po obu ich
// stronach mediana wskazuje wartość strony liczniejszej — czego rozmycie Gaussa
// nie robi.
//
// Kanał alfa liczy się tą samą medianą co barwy, osobno: mieszanie go z barwami
// zmieniłoby przezroczystość na krawędziach wycięcia.
func odszumMedianaArsenalu(obraz image.Image, przebiegi int) image.Image {
	wynik := imaging.Clone(obraz)
	for i := 0; i < przebiegi; i++ {
		wynik = medianaJednymPrzebiegiemArsenalu(wynik)
	}
	return wynik
}

// medianaJednymPrzebiegiemArsenalu liczy jeden przebieg filtru medianowego.
// Brzeg bierze tylko istniejących sąsiadów — dopisanie brakującym wartości
// czarnej wciągnęłoby ciemną obwódkę wokół obrazu.
func medianaJednymPrzebiegiemArsenalu(zrodlo *image.NRGBA) *image.NRGBA {
	granice := zrodlo.Bounds()
	wynik := image.NewNRGBA(granice)

	var czerwone, zielone, niebieskie, krycie [9]uint8
	for y := granice.Min.Y; y < granice.Max.Y; y++ {
		for x := granice.Min.X; x < granice.Max.X; x++ {
			liczba := 0
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					sasiadX, sasiadY := x+dx, y+dy
					if sasiadX < granice.Min.X || sasiadX >= granice.Max.X ||
						sasiadY < granice.Min.Y || sasiadY >= granice.Max.Y {
						continue
					}
					przesuniecie := zrodlo.PixOffset(sasiadX, sasiadY)
					czerwone[liczba] = zrodlo.Pix[przesuniecie]
					zielone[liczba] = zrodlo.Pix[przesuniecie+1]
					niebieskie[liczba] = zrodlo.Pix[przesuniecie+2]
					krycie[liczba] = zrodlo.Pix[przesuniecie+3]
					liczba++
				}
			}
			przesuniecie := wynik.PixOffset(x, y)
			wynik.Pix[przesuniecie] = medianaSkladowejArsenalu(czerwone[:liczba])
			wynik.Pix[przesuniecie+1] = medianaSkladowejArsenalu(zielone[:liczba])
			wynik.Pix[przesuniecie+2] = medianaSkladowejArsenalu(niebieskie[:liczba])
			wynik.Pix[przesuniecie+3] = medianaSkladowejArsenalu(krycie[:liczba])
		}
	}
	return wynik
}

// medianaSkladowejArsenalu wybiera wartość środkową z sąsiedztwa. Sortujemy
// kopię krótką (najwyżej dziewięć bajtów) — porządek w tablicy wołającego nie ma
// znaczenia, bo następny punkt nadpisuje ją w całości.
func medianaSkladowejArsenalu(wartosci []uint8) uint8 {
	sort.Slice(wartosci, func(i, j int) bool { return wartosci[i] < wartosci[j] })
	return wartosci[len(wartosci)/2]
}

// rozciagnijPoziomyArsenalu rozciąga histogram każdej składowej barwnej na pełny
// zakres — to jest „autoLevels" kontraktu.
//
// Skrajne wartości bierzemy z całego obrazu, ale rozciągamy tylko wtedy, gdy
// zakres jest węższy niż pełny: obraz już rozciągnięty przeszedłby przez
// mnożenie bez zmiany, a dzielenie przez zero przy obrazie jednobarwnym
// wywróciłoby rachunek. Kanał alfa zostaje nietknięty — przezroczystość nie jest
// jasnością.
func rozciagnijPoziomyArsenalu(obraz image.Image) image.Image {
	zrodlo := imaging.Clone(obraz)
	granice := zrodlo.Bounds()

	najmniejsze := [3]uint8{255, 255, 255}
	najwieksze := [3]uint8{0, 0, 0}
	for y := granice.Min.Y; y < granice.Max.Y; y++ {
		for x := granice.Min.X; x < granice.Max.X; x++ {
			przesuniecie := zrodlo.PixOffset(x, y)
			for skladowa := 0; skladowa < 3; skladowa++ {
				wartosc := zrodlo.Pix[przesuniecie+skladowa]
				if wartosc < najmniejsze[skladowa] {
					najmniejsze[skladowa] = wartosc
				}
				if wartosc > najwieksze[skladowa] {
					najwieksze[skladowa] = wartosc
				}
			}
		}
	}

	var tablice [3][256]uint8
	for skladowa := 0; skladowa < 3; skladowa++ {
		dol := najmniejsze[skladowa]
		gora := najwieksze[skladowa]
		for wejscie := 0; wejscie < 256; wejscie++ {
			if gora <= dol {
				// Składowa jednobarwna — zostaje, jaka jest. Rozciąganie zakresu
				// zerowej szerokości nie ma czego rozciągnąć.
				tablice[skladowa][wejscie] = uint8(wejscie)
				continue
			}
			rozciagniete := (float64(wejscie) - float64(dol)) * 255.0 /
				(float64(gora) - float64(dol))
			switch {
			case rozciagniete < 0:
				tablice[skladowa][wejscie] = 0
			case rozciagniete > 255:
				tablice[skladowa][wejscie] = 255
			default:
				tablice[skladowa][wejscie] = uint8(rozciagniete + 0.5)
			}
		}
	}

	for y := granice.Min.Y; y < granice.Max.Y; y++ {
		for x := granice.Min.X; x < granice.Max.X; x++ {
			przesuniecie := zrodlo.PixOffset(x, y)
			for skladowa := 0; skladowa < 3; skladowa++ {
				zrodlo.Pix[przesuniecie+skladowa] =
					tablice[skladowa][zrodlo.Pix[przesuniecie+skladowa]]
			}
		}
	}
	return zrodlo
}

// odczytajObrazArsenalu dekoduje plik źródłowy do obrazu.
//
// Odmowa dekodera znaczy „nie ma czym tego przeczytać w procesie", a nie „plik
// jest zepsuty": pod tą samą odmową kryje się AVIF, którego dekodera w Go nie
// ma. Rozstrzyga to wołający, przechodząc na program pakietu serwera — gdyby
// plik był naprawdę uszkodzony, tamta droga powie to wprost.
func odczytajObrazArsenalu(sciezka string) (image.Image, error) {
	obraz, err := odczytajObrazPliku(sciezka)
	if err != nil {
		return nil, errBrakRachunkuGoObrazu
	}
	return obraz, nil
}

// zakodujObrazArsenalu zapisuje obraz w formacie docelowym.
//
// Jakość dotyczy wyłącznie zapisów stratnych. Wskazana przy formacie
// bezstratnym nie jest błędem żądania — model prosi zwykle o jedno i drugie
// naraz — i po prostu nie ma na co wpłynąć.
func zakodujObrazArsenalu(obraz image.Image, format string, jakosc *int,
	bezstratnie *bool) ([]byte, error) {

	format = strings.ToLower(strings.TrimSpace(format))
	var bufor bytes.Buffer

	switch format {
	case "png":
		if err := imaging.Encode(&bufor, obraz, imaging.PNG); err != nil {
			return nil, err
		}
	case "jpeg", "jpg":
		opcje := []imaging.EncodeOption{}
		if jakosc != nil {
			opcje = append(opcje, imaging.JPEGQuality(*jakosc))
		}
		if err := imaging.Encode(&bufor, obraz, imaging.JPEG, opcje...); err != nil {
			return nil, err
		}
	case "gif":
		if err := imaging.Encode(&bufor, obraz, imaging.GIF); err != nil {
			return nil, err
		}
	case "tiff", "tif":
		if err := imaging.Encode(&bufor, obraz, imaging.TIFF); err != nil {
			return nil, err
		}
	case "bmp":
		if err := imaging.Encode(&bufor, obraz, imaging.BMP); err != nil {
			return nil, err
		}
	case "webp":
		// Koder czysto-Go zapisuje wyłącznie WEBP bezstratny. Prośba o zapis
		// stratny (jakość podana, bezstratność niewskazana) nie ma tu rachunku
		// i idzie do programu pakietu serwera — zapis bezstratny podany jako
		// spełnienie prośby o kompresję byłby plikiem większym, niż model prosił.
		if jakosc != nil && (bezstratnie == nil || !*bezstratnie) {
			return nil, errBrakRachunkuGoObrazu
		}
		if err := nativewebp.Encode(&bufor, obraz, nil); err != nil {
			return nil, err
		}
	default:
		// AVIF i wszystko, czego biblioteka nie zapisuje.
		return nil, errBrakRachunkuGoObrazu
	}
	return bufor.Bytes(), nil
}
