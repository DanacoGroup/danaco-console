// Moduł Library — metadane osadzone w pliku: rodzaj treści rozpoznany
// z zawartości, wymiary obrazu, liczba stron dokumentu, czas trwania nagrania
// oraz EXIF wraz ze współrzędnymi GPS, IPTC, XMP i ID3.
//
// Odczyt wchodzi wyłącznie na wyraźne żądanie (`includeTechnical`), bo otwiera
// bajty zasobu — przy wykazie kilkuset plików byłby to koszt, którego wykaz nie
// potrzebuje.
//
// Skąd co pochodzi:
//   - rodzaj treści — z zawartości, nie z rozszerzenia (`http.DetectContentType`),
//   - wymiary obrazu — z nagłówka formatu (`image.DecodeConfig`), bez dekodowania
//     całego obrazu do pamięci,
//   - liczba stron dokumentu — z `pdfcpu`, biblioteki wkompilowanej w binarium,
//   - EXIF i GPS — z własnego czytnika niżej: to kilkadziesiąt wierszy pracy na
//     strukturze TIFF, a nie powód, żeby wciągać zależność,
//   - czas trwania nagrania — z `ffprobe`, który stoi na serwerze razem
//     z rdzeniem i jest wołany jedyną dozwoloną drogą (`zewnetrzne.Wolaj`),
//   - IPTC, XMP i ID3 — z `exiftool`, tą samą drogą. Te trzy nie są jedną
//     strukturą jak EXIF, więc czytnik własny nie wchodzi w rachubę; powód
//     stoi przy `dopiszMetadaneOsadzone`.
//
// Współrzędne GPS są jedynym źródłem widoku mapy w Library Explorer — bez nich
// widok nie ma czego nanieść, więc czytnik EXIF ma je wprost, a nie „kiedyś".
package core

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"image"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/session"
	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

// errPustyObrazBiblioteki nazywa obraz bez ani jednego punktu — odcisk nie ma
// wtedy z czego powstać.
var errPustyObrazBiblioteki = errors.New("moduł Library: obraz bez punktów")

// granicaPomiaruNagrania — `ffprobe` czyta nagłówki, nie materiał, więc granica
// jest krótka.
const granicaPomiaruNagrania = 20 * time.Second

// narzedziePomiaruBiblioteki opisuje binarium arsenału używane do pomiaru
// nagrań. Stoi na serwerze razem z rdzeniem; jego brak jest brakiem pomiaru
// czasu trwania, a nie odmową całego odczytu metadanych.
var narzedziePomiaruBiblioteki = zewnetrzne.Narzedzie{
	Nazwa: "ffprobe", Program: "ffprobe", Pakiet: "ffmpeg",
}

// granicaOdczytuMetadanych — program czyta nagłówki i segmenty metadanych,
// nie treść zasobu, więc granica jest krótka. Plik, który każe mu czytać dłużej,
// jest plikiem uszkodzonym, a nie zasobem o bogatym opisie.
const granicaOdczytuMetadanych = 20 * time.Second

// narzedzieMetadanychBiblioteki opisuje binarium czytające IPTC, XMP i ID3.
// Jego brak jest brakiem trzech pól opisu, a nie odmową całego odczytu.
var narzedzieMetadanychBiblioteki = zewnetrzne.Narzedzie{
	Nazwa: "ExifTool", Program: "exiftool", Pakiet: "libimage-exiftool-perl",
}

// metadaneTechniczne czyta metadane osadzone w bajtach zasobu.
//
// Odczyt nie odmawia: zasób bez treści pod odwołaniem, format nieznany czy
// nagłówek uszkodzony oddają mniej pól, a nie błąd całej komendy. Opis zasobu
// jest wtedy niepełny i to jest prawda o pliku, nie usterka rdzenia.
func (a *adapterBiblioteki) metadaneTechniczne(zasob dane.PlikBiblioteki) *shared.LibraryTechnicalMetadata {
	if zasob.TrescOdwolanie == nil || *zasob.TrescOdwolanie == "" {
		return nil
	}
	bajty, err := os.ReadFile(*zasob.TrescOdwolanie)
	if err != nil {
		return nil
	}

	rozmiar := int64(len(bajty))
	skrot := sha256.Sum256(bajty)
	suma := hex.EncodeToString(skrot[:])
	rodzaj := http.DetectContentType(bajty)
	techniczne := &shared.LibraryTechnicalMetadata{
		MimeType: &rodzaj, SizeBytes: &rozmiar, Checksum: &suma,
	}

	if opis, _, err := image.DecodeConfig(bytes.NewReader(bajty)); err == nil {
		szerokosc, wysokosc := opis.Width, opis.Height
		techniczne.Width, techniczne.Height = &szerokosc, &wysokosc
	}
	if bytes.HasPrefix(bajty, []byte("%PDF-")) {
		if strony, err := api.PageCount(bytes.NewReader(bajty), nastawyPdf()); err == nil {
			techniczne.PageCount = &strony
		}
	}
	if wpisy, szerokosc, dlugosc := odczytajExifBiblioteki(bajty); len(wpisy) > 0 {
		if zapis, err := json.Marshal(wpisy); err == nil {
			techniczne.Exif = zapis
		}
		techniczne.GpsLatitude, techniczne.GpsLongitude = szerokosc, dlugosc
	}
	if czas := a.zmierzCzasTrwania(*zasob.TrescOdwolanie, rodzaj); czas != nil {
		techniczne.DurationMs = czas
	}
	a.dopiszMetadaneOsadzone(*zasob.TrescOdwolanie, techniczne)
	return techniczne
}

// dopiszMetadaneOsadzone wypełnia trzy pola kontraktu, których czytnik wyżej nie
// umie przeczytać: IPTC, XMP i ID3.
//
// ── Dlaczego programem, skoro EXIF czyta czytnik własny ─────────────────────
// EXIF jest jedną strukturą TIFF i mieści się w kilkudziesięciu wierszach — to
// jest powód, dla którego stoi wyżej jako kod. Pozostałe trzy nie są jedną
// strukturą: IPTC jest zapisem rekordowym w segmencie APP13, XMP drzewem RDF/XML
// osadzanym inaczej w każdym formacie kontenera, ID3 dwiema niezgodnymi
// rodzinami wersji. Napisanie ich od nowa byłoby przepisaniem cudzej pracy
// wieloletniej, a nie kilkudziesięcioma wierszami — i właśnie takiego przypadku
// dotyczy druga połowa reguły produktu: nie ma biblioteki, program jest
// składnikiem pakietu serwera.
//
// ── Poszerzenie, nie warunek ────────────────────────────────────────────────
// Odczyt opisu nie odmawia z żadnego powodu (patrz `metadaneTechniczne`) i ta
// droga tego nie zmienia: brak programu, plik bez tych metadanych albo
// odpowiedź, której nie da się rozebrać, zostawiają pola puste. Opis zasobu
// jest wtedy węższy, a nie błędny — dokładnie tak, jak przy braku `ffprobe`.
func (a *adapterBiblioteki) dopiszMetadaneOsadzone(sciezka string,
	techniczne *shared.LibraryTechnicalMetadata) {

	if a.uruchamiacz == nil || !zewnetrzne.Stoi(narzedzieMetadanychBiblioteki) {
		return
	}
	ctx, przerwij := context.WithTimeout(context.Background(), granicaOdczytuMetadanych)
	defer przerwij()

	okno := session.Okno{Ustawienia: session.Ustawienia{
		SrodowiskoWykonania: shared.ExecutionEnvCore,
	}}
	zasady := session.Zasady{}
	if a.rozstrzygacz != nil {
		zasady = ZasadyIzolacji(a.rozstrzygacz, konfig.Kontekst{})
	}
	obszar := session.Obszar{}
	if a.katalog != nil {
		obszar = ObszarOkna(a.katalog.Ustal(konfig.Kontekst{}, ""), "")
	}

	// `-g1` grupuje wynik rodziną pierwszą, więc odpowiedź sama mówi, do
	// którego pola kontraktu należy każdy wpis. Bez grupowania trzeba by wołać
	// program trzy razy albo zgadywać przynależność po nazwie znacznika.
	// `-n` wyłącza upiększanie wartości: pole ma nieść to, co stoi w pliku.
	wynik, err := zewnetrzne.Wolaj(ctx, a.uruchamiacz, okno, zasady, obszar,
		narzedzieMetadanychBiblioteki, []string{
			"-json", "-n", "-g1", "-IPTC:all", "-XMP:all", "-ID3:all", sciezka,
		}, "", granicaOdczytuMetadanych)
	if err != nil {
		return
	}

	// Program oddaje tablicę o jednym elemencie na plik — pytamy o jeden plik.
	var odpowiedz []map[string]json.RawMessage
	if err := json.Unmarshal(wynik.Wyjscie, &odpowiedz); err != nil || len(odpowiedz) == 0 {
		return
	}
	iptc, xmp, id3 := map[string]json.RawMessage{}, map[string]json.RawMessage{}, map[string]json.RawMessage{}
	for grupa, wpisy := range odpowiedz[0] {
		// Nazwy grup rodziny pierwszej niosą wariant zapisu w przyrostku
		// (`XMP-dc`, `ID3v2_4`), a kontrakt ma po jednym polu na rodzinę —
		// dlatego rozstrzyga przedrostek, a warianty scalają się w jedno pole.
		switch {
		case strings.HasPrefix(grupa, "IPTC"):
			iptc[grupa] = wpisy
		case strings.HasPrefix(grupa, "XMP"):
			xmp[grupa] = wpisy
		case strings.HasPrefix(grupa, "ID3"):
			id3[grupa] = wpisy
		}
	}
	for _, pole := range []struct {
		wpisy map[string]json.RawMessage
		cel   *json.RawMessage
	}{
		{iptc, &techniczne.Iptc}, {xmp, &techniczne.Xmp}, {id3, &techniczne.Id3},
	} {
		if len(pole.wpisy) == 0 {
			continue
		}
		if zapis, err := json.Marshal(pole.wpisy); err == nil {
			*pole.cel = zapis
		}
	}
}

// zmierzCzasTrwania woła `ffprobe` dla materiału dźwiękowego i filmowego.
//
// Brak binarium nie jest tu odmową: pole czasu trwania po prostu nie wchodzi do
// odpowiedzi. Materiał niebędący nagraniem nie jest w ogóle mierzony — pomiar
// dokumentu byłby uruchomieniem procesu bez powodu.
func (a *adapterBiblioteki) zmierzCzasTrwania(sciezka, rodzaj string) *int {
	if !strings.HasPrefix(rodzaj, "audio/") && !strings.HasPrefix(rodzaj, "video/") {
		return nil
	}
	if a.uruchamiacz == nil {
		return nil
	}
	ctx, przerwij := context.WithTimeout(context.Background(), granicaPomiaruNagrania)
	defer przerwij()

	// Zasięg pomiaru jest platformowy — żądanie opisu zasobu okna nie niesie,
	// bo pyta o plik, nie o okno (ten sam powód co przy rodzinie mediów,
	// `adapter_narzedzia_media.go`).
	okno := session.Okno{Ustawienia: session.Ustawienia{
		SrodowiskoWykonania: shared.ExecutionEnvCore,
	}}
	zasady := session.Zasady{}
	if a.rozstrzygacz != nil {
		zasady = ZasadyIzolacji(a.rozstrzygacz, konfig.Kontekst{})
	}
	obszar := session.Obszar{}
	if a.katalog != nil {
		obszar = ObszarOkna(a.katalog.Ustal(konfig.Kontekst{}, ""), "")
	}

	wynik, err := zewnetrzne.Wolaj(ctx, a.uruchamiacz, okno, zasady, obszar,
		narzedziePomiaruBiblioteki, []string{
			"-v", "error", "-show_entries", "format=duration",
			"-of", "default=noprint_wrappers=1:nokey=1", sciezka,
		}, "", granicaPomiaruNagrania)
	if err != nil {
		return nil
	}
	sekundy, err := strconv.ParseFloat(strings.TrimSpace(string(wynik.Wyjscie)), 64)
	if err != nil {
		return nil
	}
	milisekundy := int(sekundy * 1000)
	return &milisekundy
}

// ── Czytnik EXIF ────────────────────────────────────────────────────────────
//
// EXIF jest strukturą TIFF osadzoną w segmencie APP1 pliku JPEG albo stojącą
// wprost na początku pliku TIFF. Czytnik przechodzi katalog główny, katalog EXIF
// i katalog GPS, biorąc znaczniki, które mówią coś Operatorowi: aparat, czas
// zdjęcia, parametry naświetlenia i położenie.

// znacznikiExifBiblioteki nazywa znaczniki po ludzku — surowy numer znacznika
// w odpowiedzi byłby danymi bez znaczenia.
var znacznikiExifBiblioteki = map[uint16]string{
	0x010F: "wytworca",
	0x0110: "model",
	0x0112: "orientacja",
	0x0132: "czasZmiany",
	0x829A: "czasNaswietlania",
	0x829D: "przyslona",
	0x8827: "czulosc",
	0x9003: "czasZdjecia",
	0x920A: "ogniskowa",
	0xA002: "szerokoscPikseli",
	0xA003: "wysokoscPikseli",
}

// odczytajExifBiblioteki oddaje wpisy EXIF oraz współrzędne GPS, gdy są.
func odczytajExifBiblioteki(bajty []byte) (map[string]string, *float64, *float64) {
	blok := blokTiffBiblioteki(bajty)
	if len(blok) < 8 {
		return nil, nil, nil
	}
	var kolejnosc binary.ByteOrder
	switch {
	case bytes.HasPrefix(blok, []byte("II")):
		kolejnosc = binary.LittleEndian
	case bytes.HasPrefix(blok, []byte("MM")):
		kolejnosc = binary.BigEndian
	default:
		return nil, nil, nil
	}
	poczatek := kolejnosc.Uint32(blok[4:8])
	if int(poczatek) >= len(blok) {
		return nil, nil, nil
	}

	wpisy := map[string]string{}
	gps := map[uint16][]byte{}
	odczytajKatalogExif(blok, kolejnosc, int(poczatek), wpisy, gps, 0)
	szerokosc, dlugosc := wspolrzedneExif(blok, kolejnosc, gps)
	return wpisy, szerokosc, dlugosc
}

// blokTiffBiblioteki wydobywa blok TIFF: z segmentu APP1 pliku JPEG albo
// z początku pliku TIFF.
func blokTiffBiblioteki(bajty []byte) []byte {
	if bytes.HasPrefix(bajty, []byte{0x49, 0x49, 0x2A, 0x00}) ||
		bytes.HasPrefix(bajty, []byte{0x4D, 0x4D, 0x00, 0x2A}) {
		return bajty
	}
	if !bytes.HasPrefix(bajty, []byte{0xFF, 0xD8}) {
		return nil
	}
	pozycja := 2
	for pozycja+4 <= len(bajty) {
		if bajty[pozycja] != 0xFF {
			return nil
		}
		znacznik := bajty[pozycja+1]
		if znacznik == 0xDA || znacznik == 0xD9 {
			return nil
		}
		if pozycja+4 > len(bajty) {
			return nil
		}
		dlugosc := int(binary.BigEndian.Uint16(bajty[pozycja+2 : pozycja+4]))
		poczatek := pozycja + 4
		koniec := poczatek + dlugosc - 2
		if koniec > len(bajty) || koniec < poczatek {
			return nil
		}
		if znacznik == 0xE1 && bytes.HasPrefix(bajty[poczatek:koniec], []byte("Exif\x00\x00")) {
			return bajty[poczatek+6 : koniec]
		}
		pozycja = koniec
	}
	return nil
}

// odczytajKatalogExif przechodzi jeden katalog IFD, schodząc do katalogu EXIF
// i katalogu GPS. Głębokość jest ograniczona, bo plik uszkodzony potrafi
// wskazywać katalog na samego siebie.
func odczytajKatalogExif(blok []byte, kolejnosc binary.ByteOrder, przesuniecie int,
	wpisy map[string]string, gps map[uint16][]byte, glebokosc int) {

	if glebokosc > 3 || przesuniecie+2 > len(blok) {
		return
	}
	liczba := int(kolejnosc.Uint16(blok[przesuniecie : przesuniecie+2]))
	pozycja := przesuniecie + 2
	for indeks := 0; indeks < liczba; indeks++ {
		if pozycja+12 > len(blok) {
			return
		}
		wpis := blok[pozycja : pozycja+12]
		pozycja += 12

		znacznik := kolejnosc.Uint16(wpis[0:2])
		typ := kolejnosc.Uint16(wpis[2:4])
		ile := int(kolejnosc.Uint32(wpis[4:8]))
		tresc := trescWpisuExif(blok, kolejnosc, wpis, typ, ile)

		switch znacznik {
		case 0x8769: // katalog EXIF
			odczytajKatalogExif(blok, kolejnosc, int(kolejnosc.Uint32(wpis[8:12])),
				wpisy, gps, glebokosc+1)
		case 0x8825: // katalog GPS
			odczytajKatalogGps(blok, kolejnosc, int(kolejnosc.Uint32(wpis[8:12])), gps)
		default:
			nazwa, znany := znacznikiExifBiblioteki[znacznik]
			if !znany || tresc == "" {
				continue
			}
			wpisy[nazwa] = tresc
		}
	}
}

// odczytajKatalogGps zbiera surowe wpisy katalogu GPS — przekład na stopnie
// należy do `wspolrzedneExif`.
func odczytajKatalogGps(blok []byte, kolejnosc binary.ByteOrder, przesuniecie int,
	gps map[uint16][]byte) {

	if przesuniecie+2 > len(blok) {
		return
	}
	liczba := int(kolejnosc.Uint16(blok[przesuniecie : przesuniecie+2]))
	pozycja := przesuniecie + 2
	for indeks := 0; indeks < liczba; indeks++ {
		if pozycja+12 > len(blok) {
			return
		}
		wpis := blok[pozycja : pozycja+12]
		pozycja += 12
		gps[kolejnosc.Uint16(wpis[0:2])] = append([]byte{}, wpis...)
	}
}

// trescWpisuExif zamienia wartość wpisu na tekst. Obsługiwane są typy, które
// niosą coś czytelnego: napis, liczba całkowita i ułamek.
func trescWpisuExif(blok []byte, kolejnosc binary.ByteOrder, wpis []byte,
	typ uint16, ile int) string {

	dane := wartoscWpisuExif(blok, kolejnosc, wpis, typ, ile)
	if len(dane) == 0 {
		return ""
	}
	switch typ {
	case 2: // ASCII
		return strings.TrimRight(string(dane), "\x00 ")
	case 3: // SHORT
		if len(dane) < 2 {
			return ""
		}
		return strconv.Itoa(int(kolejnosc.Uint16(dane[0:2])))
	case 4: // LONG
		if len(dane) < 4 {
			return ""
		}
		return strconv.FormatUint(uint64(kolejnosc.Uint32(dane[0:4])), 10)
	case 5, 10: // RATIONAL i SRATIONAL
		if len(dane) < 8 {
			return ""
		}
		licznik := kolejnosc.Uint32(dane[0:4])
		mianownik := kolejnosc.Uint32(dane[4:8])
		if mianownik == 0 {
			return strconv.FormatUint(uint64(licznik), 10)
		}
		return strconv.FormatFloat(float64(licznik)/float64(mianownik), 'f', -1, 64)
	default:
		return ""
	}
}

// wartoscWpisuExif oddaje bajty wartości: wpis mieszczący się w czterech
// bajtach niesie je przy sobie, dłuższy — wskazuje przesunięcie w bloku.
func wartoscWpisuExif(blok []byte, kolejnosc binary.ByteOrder, wpis []byte,
	typ uint16, ile int) []byte {

	rozmiary := map[uint16]int{1: 1, 2: 1, 3: 2, 4: 4, 5: 8, 6: 1, 7: 1, 8: 2, 9: 4, 10: 8}
	rozmiar, znany := rozmiary[typ]
	if !znany || ile <= 0 {
		return nil
	}
	dlugosc := rozmiar * ile
	if dlugosc <= 4 {
		return wpis[8 : 8+dlugosc]
	}
	przesuniecie := int(kolejnosc.Uint32(wpis[8:12]))
	if przesuniecie < 0 || przesuniecie+dlugosc > len(blok) {
		return nil
	}
	return blok[przesuniecie : przesuniecie+dlugosc]
}

// wspolrzedneExif przekłada wpisy GPS na stopnie dziesiętne.
//
// EXIF trzyma położenie jako trójkę stopnie-minuty-sekundy wraz z półkulą
// w osobnym wpisie. Widok mapy potrzebuje liczby, więc przekład jest tutaj, a
// nie w oknie — inaczej każde okno robiłoby go po swojemu.
func wspolrzedneExif(blok []byte, kolejnosc binary.ByteOrder,
	gps map[uint16][]byte) (*float64, *float64) {

	stopnie := func(znacznikWartosci, znacznikPolkuli uint16, ujemna string) *float64 {
		wpis, jest := gps[znacznikWartosci]
		if !jest {
			return nil
		}
		wartosci := wartoscWpisuExif(blok, kolejnosc, wpis, 5, 3)
		if len(wartosci) < 24 {
			return nil
		}
		czlon := func(przesuniecie int) float64 {
			licznik := kolejnosc.Uint32(wartosci[przesuniecie : przesuniecie+4])
			mianownik := kolejnosc.Uint32(wartosci[przesuniecie+4 : przesuniecie+8])
			if mianownik == 0 {
				return 0
			}
			return float64(licznik) / float64(mianownik)
		}
		wynik := czlon(0) + czlon(8)/60 + czlon(16)/3600
		if polkula, jest := gps[znacznikPolkuli]; jest {
			zapis := strings.TrimRight(string(wartoscWpisuExif(blok, kolejnosc, polkula, 2, 2)), "\x00 ")
			if strings.EqualFold(zapis, ujemna) {
				wynik = -wynik
			}
		}
		return &wynik
	}
	return stopnie(0x0002, 0x0001, "S"), stopnie(0x0004, 0x0003, "W")
}
