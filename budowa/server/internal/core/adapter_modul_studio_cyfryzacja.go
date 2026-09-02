// Plik obsługuje moduł Studio Ingest/OCR Panel: kolejkę wczytywania po stronie rdzenia,
// rozpoznanie pisma z pełnym sterowaniem, korektę rozpoznania i przyjęcie wyniku jako
// dokumentu roboczego.
package core

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

const (
	przedrostekPozycjiWczytywania = "studio-wcz-"
	// Wartość domyślna języka rozpoznawania. Ta sama, którą bierze rodzina
	// `document.*` — dwa różne domyślne języki w jednym produkcie dawałyby dwa
	// różne odczyty tego samego skanu.
	jezykRozpoznaniaStudia = "pol"
)

// nastawyRozpoznania to nastawy zapamiętane przy pozycji kolejki. Struktura kontraktu
// jedzie tu przez JSON, bo pozycja pamięta je po to, żeby ponowienie po poprawce obrazu
// poszło tymi samymi nastawami, którymi szło pierwsze rozpoznanie.
type nastawyRozpoznania struct {
	Silnik                string   `json:"silnik,omitempty"`
	Jezyki                []string `json:"jezyki,omitempty"`
	MinPewnosc            *float64 `json:"minPewnosc,omitempty"`
	Prostowanie           bool     `json:"prostowanie,omitempty"`
	Odszumianie           bool     `json:"odszumianie,omitempty"`
	Progowanie            bool     `json:"progowanie,omitempty"`
	PrzycinanieMarginesow bool     `json:"przycinanieMarginesow,omitempty"`
	Uklad                 bool     `json:"uklad,omitempty"`
	StronaOd              *int     `json:"stronaOd,omitempty"`
	StronaDo              *int     `json:"stronaDo,omitempty"`
}

// DolozDoKolejki obsługuje komendę studio.ingest.queue.add. Archiwum rozpakowuje się do
// pozycji, a nie zostaje jedną pozycją, bo wsad ma być kolejką materiałów, nie kolejką
// paczek.
func (a *adapterStudia) DolozDoKolejki(ctx context.Context,
	z shared.StudioIngestQueueAddRequest) (shared.StudioIngestQueueAddResponse, error) {

	if z.WindowId == "" {
		return shared.StudioIngestQueueAddResponse{},
			bladWskazaniaStudio("dołożenie do kolejki bez okna")
	}

	nastawy := złóżNastawy(z.Settings)
	nastawyJSON, err := json.Marshal(nastawy)
	if err != nil {
		return shared.StudioIngestQueueAddResponse{}, bladStudio(err)
	}
	zapisaneNastawy := string(nastawyJSON)

	sciezki := append([]string{}, z.SourcePaths...)
	if archiwum := wartoscTekstu(z.ArchivePath); archiwum != "" {
		rozpakowane, err := a.rozpakujWsad(ctx, archiwum)
		if err != nil {
			return shared.StudioIngestQueueAddResponse{}, err
		}
		sciezki = append(sciezki, rozpakowane...)
	}

	if len(sciezki) == 0 && len(z.AssetIds) == 0 {
		return shared.StudioIngestQueueAddResponse{},
			bladWskazaniaStudio("dołożenie do kolejki bez wskazania materiału: " +
				"podaj ścieżkę, zasób magazynu albo archiwum")
	}

	pozycje := make([]shared.StudioIngestItem, 0, len(sciezki)+len(z.AssetIds))
	for _, sciezka := range sciezki {
		wskazanie := sciezka
		zapisana, err := a.repozytorium.ZapiszPozycjeWczytywania(ctx, dane.PozycjaWczytywania{
			Kod:             nowyIdentyfikator(przedrostekPozycjiWczytywania),
			Okno:            z.WindowId,
			SciezkaZrodlowa: &wskazanie,
			Stan:            "oczekuje",
			NastawyJSON:     &zapisaneNastawy,
		})
		if err != nil {
			return shared.StudioIngestQueueAddResponse{}, bladStudio(err)
		}
		pozycje = append(pozycje, złóżPozycjeWczytywania(zapisana))
	}
	for _, zasob := range z.AssetIds {
		wskazanie := zasob
		zapisana, err := a.repozytorium.ZapiszPozycjeWczytywania(ctx, dane.PozycjaWczytywania{
			Kod:         nowyIdentyfikator(przedrostekPozycjiWczytywania),
			Okno:        z.WindowId,
			ZasobID:     &wskazanie,
			Stan:        "oczekuje",
			NastawyJSON: &zapisaneNastawy,
		})
		if err != nil {
			return shared.StudioIngestQueueAddResponse{}, bladStudio(err)
		}
		pozycje = append(pozycje, złóżPozycjeWczytywania(zapisana))
	}

	return shared.StudioIngestQueueAddResponse{Items: pozycje}, nil
}

// KolejkaWczytywania obsługuje komendę studio.ingest.queue.list, oddając pozycje kolejki
// wraz z ich stanem i nastawami.
func (a *adapterStudia) KolejkaWczytywania(ctx context.Context,
	z shared.StudioIngestQueueListRequest) (shared.StudioIngestQueueListResponse, error) {

	if z.WindowId == "" {
		return shared.StudioIngestQueueListResponse{},
			bladWskazaniaStudio("odczyt kolejki bez okna")
	}
	tylkoNieprzetworzone := z.PendingOnly != nil && *z.PendingOnly
	wiersze, err := a.repozytorium.PozycjeWczytywania(ctx, z.WindowId, tylkoNieprzetworzone)
	if err != nil {
		return shared.StudioIngestQueueListResponse{}, bladStudio(err)
	}
	pozycje := make([]shared.StudioIngestItem, 0, len(wiersze))
	for _, wiersz := range wiersze {
		pozycje = append(pozycje, złóżPozycjeWczytywania(wiersz))
	}
	return shared.StudioIngestQueueListResponse{Items: pozycje}, nil
}

// Rozpoznaj obsługuje komendę studio.ingest.recognize. Pozycja przechodzi tu przez trzy
// stany: przetwarzanie na czas wywołania, potem gotowa albo odmowa, a przy pewności
// poniżej progu — ponowienie.
func (a *adapterStudia) Rozpoznaj(ctx context.Context,
	z shared.StudioIngestRecognizeRequest) (shared.StudioIngestRecognizeResponse, error) {

	if z.ItemId == "" {
		return shared.StudioIngestRecognizeResponse{},
			bladWskazaniaStudio("rozpoznanie bez wskazania pozycji kolejki")
	}
	pozycja, err := a.repozytorium.PozycjaWczytywania(ctx, z.ItemId)
	if err != nil {
		if errors.Is(err, dane.ErrBrakWiersza) {
			return shared.StudioIngestRecognizeResponse{},
				bladBrakuStudio("pozycja kolejki nie istnieje: " + z.ItemId)
		}
		return shared.StudioIngestRecognizeResponse{}, bladStudio(err)
	}

	nastawy := nastawyPozycji(pozycja)
	if z.Settings != nil {
		nastawy = złóżNastawy(z.Settings)
	}

	sciezka, err := a.materialPozycji(ctx, pozycja)
	if err != nil {
		return shared.StudioIngestRecognizeResponse{}, err
	}

	pozycja.Stan = "przetwarzanie"
	if _, err := a.repozytorium.ZapiszPozycjeWczytywania(ctx, pozycja); err != nil {
		return shared.StudioIngestRecognizeResponse{}, bladStudio(err)
	}

	odczyt, err := a.rozpoznajMaterial(ctx, sciezka, nastawy)
	if err != nil {
		// Odmowa rozpoznania zostaje przy pozycji, nie tylko w odpowiedzi.

		// Powód jest widoczny przy materiale także po odświeżeniu okna.
		powod := err.Error()
		pozycja.Stan = "odmowa"
		pozycja.PowodOdmowy = &powod
		if _, zapis := a.repozytorium.ZapiszPozycjeWczytywania(ctx, pozycja); zapis != nil {
			return shared.StudioIngestRecognizeResponse{}, bladStudio(zapis)
		}
		return shared.StudioIngestRecognizeResponse{}, err
	}

	pozycja.Tekst = &odczyt.tekst
	pozycja.Stron = &odczyt.stron
	pozycja.UzytoRozpoznania = true
	pozycja.Pewnosc = odczyt.pewnosc
	pozycja.PowodOdmowy = nil
	pozycja.Stan = "gotowa"
	if nastawy.MinPewnosc != nil && odczyt.pewnosc != nil && *odczyt.pewnosc < *nastawy.MinPewnosc {
		pozycja.Stan = "ponowienie"
	}

	slowaJSON, err := json.Marshal(odczyt.slowa)
	if err != nil {
		return shared.StudioIngestRecognizeResponse{}, bladStudio(err)
	}
	zapisSlow := string(slowaJSON)
	pozycja.SlowaJSON = &zapisSlow

	var uklad []shared.StudioLayoutBlock
	if nastawy.Uklad {
		uklad = blokiUkladu(odczyt.slowa)
		ukladJSON, err := json.Marshal(uklad)
		if err != nil {
			return shared.StudioIngestRecognizeResponse{}, bladStudio(err)
		}
		zapisUkladu := string(ukladJSON)
		pozycja.UkladJSON = &zapisUkladu
	}

	zapisana, err := a.repozytorium.ZapiszPozycjeWczytywania(ctx, pozycja)
	if err != nil {
		return shared.StudioIngestRecognizeResponse{}, bladStudio(err)
	}

	return shared.StudioIngestRecognizeResponse{
		Item:   złóżPozycjeWczytywania(zapisana),
		Layout: uklad,
		Words:  odczyt.slowa,
	}, nil
}

// PoprawRozpoznanie obsługuje komendę studio.ingest.correction.set. Poprawka zmienia dwie
// rzeczy naraz: słowo na warstwie tekstowej oraz tekst pozycji, bo zmiana samego słowa
// zostawiłaby tekst w postaci sprzed poprawki.
func (a *adapterStudia) PoprawRozpoznanie(ctx context.Context,
	z shared.StudioIngestCorrectionSetRequest) (shared.StudioIngestCorrectionSetResponse, error) {

	if z.ItemId == "" {
		return shared.StudioIngestCorrectionSetResponse{},
			bladWskazaniaStudio("korekta bez wskazania pozycji kolejki")
	}
	pozycja, err := a.repozytorium.PozycjaWczytywania(ctx, z.ItemId)
	if err != nil {
		if errors.Is(err, dane.ErrBrakWiersza) {
			return shared.StudioIngestCorrectionSetResponse{},
				bladBrakuStudio("pozycja kolejki nie istnieje: " + z.ItemId)
		}
		return shared.StudioIngestCorrectionSetResponse{}, bladStudio(err)
	}
	if pozycja.SlowaJSON == nil || *pozycja.SlowaJSON == "" {
		return shared.StudioIngestCorrectionSetResponse{},
			bladWskazaniaStudio("pozycja nie ma warstwy tekstowej — rozpoznaj ją przed korektą")
	}

	var slowa []shared.StudioRecognizedWord
	if err := json.Unmarshal([]byte(*pozycja.SlowaJSON), &slowa); err != nil {
		return shared.StudioIngestCorrectionSetResponse{}, bladStudio(err)
	}
	miejsce := -1
	for i, slowo := range slowa {
		if slowo.Index == z.WordIndex {
			miejsce = i
			break
		}
	}
	if miejsce < 0 {
		return shared.StudioIngestCorrectionSetResponse{},
			bladBrakuStudio(fmt.Sprintf("warstwa tekstowa nie ma słowa o numerze %d", z.WordIndex))
	}

	slowa[miejsce].Text = z.Text
	slowa[miejsce].Corrected = wskaznikLogiczny(true)

	zapisSlow, err := json.Marshal(slowa)
	if err != nil {
		return shared.StudioIngestCorrectionSetResponse{}, bladStudio(err)
	}
	tekstPoKorekcie := tekstZeSlow(slowa)
	zapis := string(zapisSlow)
	pozycja.SlowaJSON = &zapis
	pozycja.Tekst = &tekstPoKorekcie

	zapisana, err := a.repozytorium.ZapiszPozycjeWczytywania(ctx, pozycja)
	if err != nil {
		return shared.StudioIngestCorrectionSetResponse{}, bladStudio(err)
	}
	return shared.StudioIngestCorrectionSetResponse{Item: złóżPozycjeWczytywania(zapisana)}, nil
}

// PrzyjmijPozycje obsługuje komendę studio.ingest.item.accept. Wiele pozycji składa się w
// jeden dokument w kolejności podania — tak działa skan wielostronicowy rozłożony na
// pliki.
func (a *adapterStudia) PrzyjmijPozycje(ctx context.Context,
	z shared.StudioIngestItemAcceptRequest) (shared.StudioIngestItemAcceptResponse, error) {

	if z.WindowId == "" {
		return shared.StudioIngestItemAcceptResponse{},
			bladWskazaniaStudio("przyjęcie pozycji bez okna")
	}
	if len(z.ItemIds) == 0 {
		return shared.StudioIngestItemAcceptResponse{},
			bladWskazaniaStudio("przyjęcie bez wskazania pozycji kolejki")
	}

	czesci := make([]string, 0, len(z.ItemIds))
	for _, kod := range z.ItemIds {
		pozycja, err := a.repozytorium.PozycjaWczytywania(ctx, kod)
		if err != nil {
			if errors.Is(err, dane.ErrBrakWiersza) {
				return shared.StudioIngestItemAcceptResponse{},
					bladBrakuStudio("pozycja kolejki nie istnieje: " + kod)
			}
			return shared.StudioIngestItemAcceptResponse{}, bladStudio(err)
		}
		if pozycja.Tekst == nil || *pozycja.Tekst == "" {
			return shared.StudioIngestItemAcceptResponse{},
				bladWskazaniaStudio("pozycja " + kod + " nie ma jeszcze tekstu — rozpoznaj ją przed przyjęciem")
		}
		czesci = append(czesci, *pozycja.Tekst)
	}
	tresc := strings.Join(czesci, "\n\n")

	var format shared.StudioDocumentFormat = shared.StudioDocumentFormatTxt
	if z.Format != nil && *z.Format != "" {
		format = *z.Format
	}
	dokument := dane.DokumentStudia{
		Kod:    nowyIdentyfikator(przedrostekDokumentuStudio),
		Okno:   z.WindowId,
		Tytul:  z.Title,
		Format: format,
		Tresc:  &tresc,
	}
	zapisany, err := a.repozytorium.ZapiszDokument(ctx, dokument)
	if err != nil {
		return shared.StudioIngestItemAcceptResponse{}, bladStudio(err)
	}

	autor := string(shared.StudioAuthorModel)
	wersja, err := a.repozytorium.ZapiszWersje(ctx, zapisany.ID, dane.WersjaDokumentu{
		Kod:   nowyIdentyfikator(przedrostekWersjiStudio),
		Tresc: &tresc,
		Autor: &autor,
	})
	if err != nil {
		return shared.StudioIngestItemAcceptResponse{}, bladStudio(err)
	}
	zapisany.WersjaBiezacaKod = &wersja.Kod
	zapisany, err = a.repozytorium.ZapiszDokument(ctx, zapisany)
	if err != nil {
		return shared.StudioIngestItemAcceptResponse{}, bladStudio(err)
	}

	return shared.StudioIngestItemAcceptResponse{
		Document: a.zlozDokument(zapisany),
		Version:  a.zlozWersje(wersja),
	}, nil
}

// UrzadzeniaWejsciowe obsługuje komendę studio.ingest.device.list. Wykaz pochodzi od
// warstwy urządzeń systemu, nie z domysłu rdzenia: brak warstwy jest odpowiedzią „nie ma
// czym szukać", a nie pustym wykazem.
func (a *adapterStudia) UrzadzeniaWejsciowe(ctx context.Context,
	_ shared.StudioIngestDeviceListRequest) (shared.StudioIngestDeviceListResponse, error) {

	urzadzenia, err := a.wykazSkanerow(ctx)
	if err != nil {
		return shared.StudioIngestDeviceListResponse{}, err
	}
	return shared.StudioIngestDeviceListResponse{Devices: urzadzenia}, nil
}

// odczytajUrzadzenia wyciąga urządzenia z wyjścia `scanimage -L`. Wiersz ma
// postać: `device 'net:host:plustek:libusb:001:002' is a Plustek Scanner`.
func odczytajUrzadzenia(wyjscie string) []shared.StudioInputDevice {
	urzadzenia := []shared.StudioInputDevice{}
	for _, wiersz := range strings.Split(wyjscie, "\n") {
		wiersz = strings.TrimSpace(wiersz)
		if !strings.HasPrefix(wiersz, "device ") {
			continue
		}
		poczatek := strings.Index(wiersz, "`")
		if poczatek < 0 {
			poczatek = strings.Index(wiersz, "'")
		}
		koniec := strings.LastIndex(wiersz, "'")
		if poczatek < 0 || koniec <= poczatek {
			continue
		}
		kod := wiersz[poczatek+1 : koniec]
		nazwa := kod
		if opis := strings.TrimSpace(wiersz[koniec+1:]); strings.HasPrefix(opis, "is a ") {
			nazwa = strings.TrimSpace(strings.TrimPrefix(opis, "is a "))
		}
		urzadzenia = append(urzadzenia, shared.StudioInputDevice{
			Id: kod, Name: nazwa, Kind: shared.StudioInputDeviceKindSkaner,
		})
	}
	return urzadzenia
}

// ── Rozpoznanie ─────────────────────────────────────────────────────────────

// odczytMaterialu niesie wynik jednego przebiegu rozpoznania: tekst, liczbę stron i
// pewność, gdy Tesseract ją oddał.
type odczytMaterialu struct {
	tekst   string
	stron   int64
	pewnosc *float64
	slowa   []shared.StudioRecognizedWord
}

// rozpoznajMaterial puszcza materiał przez Tesseracta wyjściem `tsv` i składa
// z niego tekst, słowa i pewność.
func (a *adapterStudia) rozpoznajMaterial(ctx context.Context, sciezka string,
	nastawy nastawyRozpoznania) (odczytMaterialu, error) {

	jezyk := jezykRozpoznaniaStudia
	if len(nastawy.Jezyki) > 0 {
		// Tesseract przyjmuje zestaw języków rozdzielony plusem: to droga wskazania „pol
		// plus eng" naraz.
		jezyk = strings.Join(nastawy.Jezyki, "+")
	}

	material := sciezka
	if czyszczenieZadane(nastawy) {
		oczyszczony, posprzataj, err := a.oczyscMaterial(ctx, sciezka, nastawy)
		if err != nil {
			return odczytMaterialu{}, err
		}
		defer posprzataj()
		material = oczyszczony
	}

	argumenty := []string{material, "stdout", "-l", jezyk, "tsv"}
	wyjscie, err := a.wolajNarzedzie(ctx, narzedzieRozpoznaniaStudia, argumenty)
	if err != nil {
		return odczytMaterialu{}, err
	}

	slowa := odczytajSlowaTsv(string(wyjscie))
	return odczytMaterialu{
		tekst:   tekstZeSlow(slowa),
		stron:   int64(liczbaStron(slowa)),
		pewnosc: sredniaPewnosc(slowa),
		slowa:   slowa,
	}, nil
}

// czyszczenieZadane odpowiada, czy którakolwiek z czterech nastaw obróbki wstępnej jest
// włączona, rozstrzygając, czy unpaper ma w ogóle ruszyć.
func czyszczenieZadane(nastawy nastawyRozpoznania) bool {
	return nastawy.Prostowanie || nastawy.Odszumianie ||
		nastawy.Progowanie || nastawy.PrzycinanieMarginesow
}

// argumentyCzyszczenia składa wiersz wywołania unpapera z nastaw pozycji, wyłączając
// jawnie każdy filtr, o który nikt nie prosił.
func argumentyCzyszczenia(nastawy nastawyRozpoznania, wejscie, wyjscie string) []string {
	argumenty := []string{"--layout", "single", "--no-blackfilter"}
	if !nastawy.Prostowanie {
		argumenty = append(argumenty, "--no-deskew")
	}
	if !nastawy.Odszumianie {
		argumenty = append(argumenty, "--no-noisefilter", "--no-blurfilter", "--no-grayfilter")
	}
	if !nastawy.PrzycinanieMarginesow {
		argumenty = append(argumenty, "--no-border-scan", "--no-border-align")
	}
	if nastawy.Progowanie {
		argumenty = append(argumenty, "--type", "pbm")
	}
	return append(argumenty, wejscie, wyjscie)
}

// oczyscMaterial przepuszcza materiał przez unpapera i oddaje ścieżkę wyniku z
// posprzątaniem katalogu roboczego. unpaper czyta i pisze wyłącznie postać PNM, więc inny
// materiał przechodzi zamianę biblioteką przed wywołaniem.
func (a *adapterStudia) oczyscMaterial(ctx context.Context, sciezka string,
	nastawy nastawyRozpoznania) (string, func(), error) {

	pusto := func() {}
	if !zewnetrzne.Stoi(narzedzieCzyszczeniaSkanu) {
		return "", pusto, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
			"moduł Studio: nastawy pozycji żądają obróbki wstępnej obrazu, a programu "+
				narzedzieCzyszczeniaSkanu.Nazwa+" ("+narzedzieCzyszczeniaSkanu.Program+
				") nie ma na tej maszynie; naprawa: zainstalować pakiet "+
				narzedzieCzyszczeniaSkanu.Pakiet+
				". Droga, która działa bez niego: wyłączyć w nastawach prostowanie, "+
				"odszumianie, progowanie i przycinanie marginesów — rozpoznanie pobiegnie "+
				"na materiale bez obróbki"))
	}

	katalog, err := os.MkdirTemp("", "danaco-skan-")
	if err != nil {
		return "", pusto, bladStudio(err)
	}
	posprzataj := func() { _ = os.RemoveAll(katalog) }

	wejscie, err := materialWPnm(sciezka, katalog)
	if err != nil {
		posprzataj()
		return "", pusto, err
	}
	// Rozszerzenie wyniku idzie za typem, który unpaper zapisze.

	// Przy progowaniu jest to mapa bitowa, poza nim — ten sam typ, co wejście.
	wyjscie := filepath.Join(katalog, "oczyszczony.ppm")
	if nastawy.Progowanie {
		wyjscie = filepath.Join(katalog, "oczyszczony.pbm")
	}

	if _, err := a.wolajNarzedzie(ctx, narzedzieCzyszczeniaSkanu,
		argumentyCzyszczenia(nastawy, wejscie, wyjscie)); err != nil {
		posprzataj()
		return "", pusto, err
	}
	// Program potrafi skończyć się powodzeniem i nie zostawić pliku.

	// Rozpoznanie pobiegłoby wtedy na ścieżce, pod którą nic nie leży.
	if opis, err := os.Stat(wyjscie); err != nil || opis.Size() == 0 {
		posprzataj()
		return "", pusto, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"moduł Studio: unpaper zakończył pracę, ale obrazu po obróbce nie ma pod "+
				wyjscie+" — materiał do rozpoznania nie powstał"))
	}
	return wyjscie, posprzataj, nil
}

// materialWPnm oddaje ścieżkę materiału w postaci, którą unpaper przyjmie, zamieniając go
// najpierw na mapę PPM, gdy nie jest już postacią PNM.
func materialWPnm(sciezka, katalog string) (string, error) {
	bajty, err := os.ReadFile(sciezka)
	if err != nil {
		return "", protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Studio: nie można odczytać materiału pozycji "+sciezka+": "+err.Error()))
	}
	if czyPnm(bajty) {
		return sciezka, nil
	}
	obraz, _, err := image.Decode(bytes.NewReader(bajty))
	if err != nil {
		return "", protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
			"moduł Studio: obróbka wstępna obrazu dotyczy skanu i zdjęcia kartki, a materiału "+
				"pozycji serwer nie umie odczytać jako obrazu ("+err.Error()+
				"); naprawa: wyłączyć w nastawach obróbkę wstępną albo podać materiał "+
				"w formacie obrazu (PNG, JPEG, TIFF, BMP, WEBP, GIF, PNM)"))
	}
	wejscie := filepath.Join(katalog, "wejscie.ppm")
	if err := zapiszPpm(wejscie, obraz); err != nil {
		return "", err
	}
	return wejscie, nil
}

// czyPnm rozpoznaje mapę PNM po jej znaku rozpoznawczym: litera `P` i cyfra od
// 1 do 6. To jest cały nagłówek tego rodzaju plików.
func czyPnm(bajty []byte) bool {
	return len(bajty) >= 2 && bajty[0] == 'P' && bajty[1] >= '1' && bajty[1] <= '6'
}

// zapiszPpm zapisuje obraz jako mapę PPM: trzy bajty na piksel, bez kompresji, na białym
// tle, bo piksel przezroczysty ma składowe pomnożone przez alfę i bez podłożenia bieli
// wyszedłby czarny.
func zapiszPpm(sciezka string, obraz image.Image) error {
	granice := obraz.Bounds()
	naPapierze := image.NewRGBA(granice)
	draw.Draw(naPapierze, granice, image.NewUniform(color.White), image.Point{}, draw.Src)
	draw.Draw(naPapierze, granice, obraz, granice.Min, draw.Over)

	plik, err := os.Create(sciezka)
	if err != nil {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"moduł Studio: nie można założyć pliku materiału do obróbki: "+err.Error()))
	}
	defer func() { _ = plik.Close() }()

	bufor := bufio.NewWriter(plik)
	if _, err := fmt.Fprintf(bufor, "P6\n%d %d\n255\n", granice.Dx(), granice.Dy()); err != nil {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"moduł Studio: nie można zapisać nagłówka materiału do obróbki: "+err.Error()))
	}
	wiersz := make([]byte, 0, granice.Dx()*3)
	for y := granice.Min.Y; y < granice.Max.Y; y++ {
		wiersz = wiersz[:0]
		for x := granice.Min.X; x < granice.Max.X; x++ {
			czerwony, zielony, niebieski, _ := naPapierze.At(x, y).RGBA()
			wiersz = append(wiersz, byte(czerwony>>8), byte(zielony>>8), byte(niebieski>>8))
		}
		if _, err := bufor.Write(wiersz); err != nil {
			return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
				"moduł Studio: nie można zapisać materiału do obróbki: "+err.Error()))
		}
	}
	if err := bufor.Flush(); err != nil {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"moduł Studio: nie można domknąć materiału do obróbki: "+err.Error()))
	}
	return nil
}

// odczytajSlowaTsv czyta wyjście `tesseract … tsv`. Kolumny: level, page_num,
// block_num, par_num, line_num, word_num, left, top, width, height, conf, text.
func odczytajSlowaTsv(wyjscie string) []shared.StudioRecognizedWord {
	slowa := []shared.StudioRecognizedWord{}
	numer := 0
	for _, wiersz := range strings.Split(wyjscie, "\n") {
		pola := strings.Split(wiersz, "\t")
		if len(pola) < 12 || pola[0] == "level" {
			continue
		}
		// Poziom 5 to słowo; poziomy niższe opisują blok, akapit i wiersz.
		if pola[0] != "5" {
			continue
		}
		tekst := strings.TrimSpace(pola[11])
		if tekst == "" {
			continue
		}
		pewnosc, err := strconv.ParseFloat(pola[10], 64)
		if err != nil {
			continue
		}
		slowa = append(slowa, shared.StudioRecognizedWord{
			Index:      numer,
			Text:       tekst,
			Page:       liczbaZTekstu(pola[1]),
			X:          liczbaZTekstu(pola[6]),
			Y:          liczbaZTekstu(pola[7]),
			Width:      liczbaZTekstu(pola[8]),
			Height:     liczbaZTekstu(pola[9]),
			Confidence: pewnosc,
		})
		numer++
	}
	return slowa
}

// tekstZeSlow składa tekst z warstwy słów. Wiersz kończy się tam, gdzie kończy
// go Tesseract — po pionowym skoku ramki — więc tekst przyjęty do edytora ma
// łamania tam, gdzie miał je materiał.
func tekstZeSlow(slowa []shared.StudioRecognizedWord) string {
	var budowa strings.Builder
	poprzedniY, poprzedniaStrona := -1, -1
	for i, slowo := range slowa {
		if i > 0 {
			switch {
			case slowo.Page != poprzedniaStrona:
				budowa.WriteString("\n\n")
			case poprzedniY >= 0 && absolutna(slowo.Y-poprzedniY) > slowo.Height/2:
				budowa.WriteString("\n")
			default:
				budowa.WriteString(" ")
			}
		}
		budowa.WriteString(slowo.Text)
		poprzedniY, poprzedniaStrona = slowo.Y, slowo.Page
	}
	return budowa.String()
}

// blokiUkladu składa bloki układu z warstwy słów, grupując je po stronie
// i wierszu. Rdzeń nie zgaduje rodzaju bloku poza jednym rozstrzygnięciem, które
// da się obronić miarą: wiersz wyraźnie wyższy od mediany jest nagłówkiem.
func blokiUkladu(slowa []shared.StudioRecognizedWord) []shared.StudioLayoutBlock {
	if len(slowa) == 0 {
		return []shared.StudioLayoutBlock{}
	}
	wysokosci := make([]int, 0, len(slowa))
	for _, slowo := range slowa {
		wysokosci = append(wysokosci, slowo.Height)
	}
	sort.Ints(wysokosci)
	mediana := wysokosci[len(wysokosci)/2]

	bloki := []shared.StudioLayoutBlock{}
	biezacy := shared.StudioLayoutBlock{}
	otwarty := false
	var tresc []string
	poprzedniY, poprzedniaStrona := -1, -1

	domknij := func() {
		if !otwarty {
			return
		}
		polaczona := strings.Join(tresc, " ")
		biezacy.Text = &polaczona
		bloki = append(bloki, biezacy)
		otwarty = false
		tresc = nil
	}

	for _, slowo := range slowa {
		nowyBlok := !otwarty || slowo.Page != poprzedniaStrona ||
			(poprzedniY >= 0 && absolutna(slowo.Y-poprzedniY) > slowo.Height)
		if nowyBlok {
			domknij()
			var rodzaj shared.StudioLayoutBlockKind = shared.StudioLayoutBlockKindAkapit
			if slowo.Height > mediana+mediana/3 {
				rodzaj = shared.StudioLayoutBlockKindNaglowek
			}
			biezacy = shared.StudioLayoutBlock{
				Kind: rodzaj, Page: slowo.Page, X: slowo.X, Y: slowo.Y,
				Width: slowo.Width, Height: slowo.Height,
			}
			otwarty = true
		}
		if slowo.X+slowo.Width-biezacy.X > biezacy.Width {
			biezacy.Width = slowo.X + slowo.Width - biezacy.X
		}
		if slowo.Y+slowo.Height-biezacy.Y > biezacy.Height {
			biezacy.Height = slowo.Y + slowo.Height - biezacy.Y
		}
		tresc = append(tresc, slowo.Text)
		poprzedniY, poprzedniaStrona = slowo.Y, slowo.Page
	}
	domknij()
	return bloki
}

// sredniaPewnosc liczy pewność pozycji ze średniej pewności jej słów. Brak słów
// daje brak pewności, a nie zero: zero znaczyłoby „rozpoznano i nic nie pasuje".
func sredniaPewnosc(slowa []shared.StudioRecognizedWord) *float64 {
	if len(slowa) == 0 {
		return nil
	}
	suma := 0.0
	for _, slowo := range slowa {
		suma += slowo.Confidence
	}
	srednia := suma / float64(len(slowa))
	return &srednia
}

func liczbaStron(slowa []shared.StudioRecognizedWord) int {
	strony := map[int]struct{}{}
	for _, slowo := range slowa {
		strony[slowo.Page] = struct{}{}
	}
	if len(strony) == 0 {
		return 0
	}
	return len(strony)
}

// ── Materiał i wsad ─────────────────────────────────────────────────────────

// materialPozycji wskazuje plik, na którym pracuje rozpoznanie, biorąc ścieżkę źródłową
// albo materiał odtworzony z zasobu pozycji.
func (a *adapterStudia) materialPozycji(ctx context.Context,
	pozycja dane.PozycjaWczytywania) (string, error) {

	if pozycja.SciezkaZrodlowa != nil && *pozycja.SciezkaZrodlowa != "" {
		return *pozycja.SciezkaZrodlowa, nil
	}
	if pozycja.ZasobID == nil || *pozycja.ZasobID == "" {
		return "", bladWskazaniaStudio("pozycja kolejki bez wskazania materiału")
	}
	sciezka, err := a.sciezkaZasobu(ctx, *pozycja.ZasobID)
	if err != nil {
		return "", err
	}
	return sciezka, nil
}

// rozpakujWsad rozpakowuje archiwum i oddaje ścieżki jego pozycji, jedną na każdy plik
// materiału znaleziony wewnątrz.
func (a *adapterStudia) rozpakujWsad(ctx context.Context, archiwum string) ([]string, error) {
	katalog, err := a.katalogWsadu(ctx, archiwum)
	if err != nil {
		return nil, err
	}
	if _, err := a.wolajNarzedzie(ctx, narzedzieArchiwumStudia,
		[]string{"x", "-y", "-o" + katalog, archiwum}); err != nil {
		return nil, err
	}
	return plikiKatalogu(katalog)
}

// ── Składanie odpowiedzi ────────────────────────────────────────────────────

func złóżPozycjeWczytywania(wiersz dane.PozycjaWczytywania) shared.StudioIngestItem {
	pozycja := shared.StudioIngestItem{
		Id:            wiersz.Kod,
		WindowId:      wiersz.Okno,
		SourcePath:    wiersz.SciezkaZrodlowa,
		AssetId:       wiersz.ZasobID,
		State:         shared.StudioIngestState(wiersz.Stan),
		Text:          wiersz.Tekst,
		UsedOcr:       wskaznikLogiczny(wiersz.UzytoRozpoznania),
		Confidence:    wiersz.Pewnosc,
		FailureReason: wiersz.PowodOdmowy,
		CreatedAt:     chwilaBazy(wiersz.Utworzono),
	}
	if wiersz.Stron != nil {
		strony := int(*wiersz.Stron)
		pozycja.Pages = &strony
	}
	return pozycja
}

// złóżNastawy sprowadza nastawy kontraktu do postaci zapamiętywanej przy pozycji, oddając
// nastawy domyślne, gdy kontrakt ich nie niesie.
func złóżNastawy(z *shared.StudioRecognitionSettings) nastawyRozpoznania {
	if z == nil {
		return nastawyRozpoznania{}
	}
	nastawy := nastawyRozpoznania{
		Jezyki:     append([]string{}, z.Languages...),
		MinPewnosc: z.MinConfidence,
	}
	if z.Engine != nil {
		nastawy.Silnik = string(*z.Engine)
	}
	nastawy.Prostowanie = wartoscLogicznaZeWskaznika(z.Deskew)
	nastawy.Odszumianie = wartoscLogicznaZeWskaznika(z.Denoise)
	nastawy.Progowanie = wartoscLogicznaZeWskaznika(z.Binarize)
	nastawy.PrzycinanieMarginesow = wartoscLogicznaZeWskaznika(z.TrimMargins)
	nastawy.Uklad = wartoscLogicznaZeWskaznika(z.DetectLayout)
	nastawy.StronaOd = z.PageFrom
	nastawy.StronaDo = z.PageTo
	return nastawy
}

// nastawyPozycji odczytuje nastawy zapamiętane przy pozycji kolejki, wracając do nastaw
// domyślnych, gdy zapis jest pusty albo nieczytelny.
func nastawyPozycji(pozycja dane.PozycjaWczytywania) nastawyRozpoznania {
	if pozycja.NastawyJSON == nil || *pozycja.NastawyJSON == "" {
		return nastawyRozpoznania{}
	}
	var nastawy nastawyRozpoznania
	if err := json.Unmarshal([]byte(*pozycja.NastawyJSON), &nastawy); err != nil {
		// Nastawy nieczytelne nie zatrzymują rozpoznania: wracamy do domyślnych.

		// Materiał jest ważniejszy od zapamiętanej nastawy.
		return nastawyRozpoznania{}
	}
	return nastawy
}

// ── Drobne ──────────────────────────────────────────────────────────────────

func liczbaZTekstu(wartosc string) int {
	liczba, err := strconv.Atoi(strings.TrimSpace(wartosc))
	if err != nil {
		return 0
	}
	return liczba
}

func absolutna(wartosc int) int {
	if wartosc < 0 {
		return -wartosc
	}
	return wartosc
}

func wskaznikLogiczny(wartosc bool) *bool {
	kopia := wartosc
	return &kopia
}

func wartoscLogicznaZeWskaznika(wartosc *bool) bool {
	return wartosc != nil && *wartosc
}

// plikiKatalogu zwraca pliki rozpakowanego wsadu, posortowane nazwą — kolejność
// stron skanu wynika z nazwy pliku, więc kolejność systemu plików byłaby losowa.
func plikiKatalogu(katalog string) ([]string, error) {
	wzorzec := filepath.Join(katalog, "*")
	znalezione, err := filepath.Glob(wzorzec)
	if err != nil {
		return nil, bladStudio(err)
	}
	sort.Strings(znalezione)
	return znalezione, nil
}

// bladBrakuStudio nazywa byt, którego nie ma — tą samą drogą co pozostałe
// odmowy modułu, żeby kod kontraktu składał się w jednym miejscu.
func bladBrakuStudio(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"moduł Studio: "+powod))
}

// narzedzia rodziny cyfryzacji. Nazwy są dla CZŁOWIEKA i wchodzą wprost do
// treści odmowy „nie ma czym", razem z pakietem do dociągnięcia.
var (
	narzedzieRozpoznaniaStudia = zewnetrzne.Narzedzie{
		Nazwa: "Tesseract OCR", Program: "tesseract", Pakiet: "tesseract-ocr tesseract-ocr-pol",
	}
	narzedzieArchiwumStudia = zewnetrzne.Narzedzie{
		Nazwa: "7-Zip", Program: "7z", Pakiet: "p7zip-full",
	}
	narzedzieSkanera = zewnetrzne.Narzedzie{
		Nazwa: "SANE (scanimage)", Program: "scanimage", Pakiet: "sane-utils",
	}
	narzedzieCzyszczeniaSkanu = zewnetrzne.Narzedzie{
		Nazwa: "unpaper", Program: "unpaper", Pakiet: "unpaper",
	}
)
