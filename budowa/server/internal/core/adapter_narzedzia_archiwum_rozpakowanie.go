// Odpowiedzialność pliku: komenda `archive.unpack` — rozpakowanie archiwum do
// katalogu roboczego okna albo do magazynu zasobów, wraz z całą obroną przed
// wyjściem poza katalog docelowy. Wyrok o zawartości archiwum wydaje
// `adapter_narzedzia_archiwum_spis.go`, zaplecze rodziny leży
// w `adapter_narzedzia_archiwum.go`.
//
// Warstwy obrony są trzy, bo żadna pojedyncza nie daje pewności:
//  1. Spis przed zapisem. Zawartość archiwum oglądamy `7z l -slt`, zanim poleci
//     pierwszy bajt, i odmawiamy pozycjom bezwzględnym, pozycjom z członem `..`
//     i dowiązaniom (uzasadnienie: plik spisu).
//  2. Kwarantanna. Rozpakowujemy do katalogu świeżo założonego i pustego, a nie
//     wprost do katalogu Operatora: w katalogu, w którym nie ma ani jednego
//     pliku, nie ma czego nadpisać. Kwarantanna stoi na tym samym nośniku co
//     cel, więc przeniesienie gotowej treści jest przemianowaniem, a nie drugim
//     kopiowaniem.
//  3. Przejście po wyniku. Po rozpakowaniu obchodzimy kwarantannę i patrzymy,
//     co powstało: `filepath.WalkDir` nie idzie za dowiązaniami i rozpoznaje je
//     po typie wpisu, więc dowiązanie, które przeszłoby przez spis, zatrzymuje
//     się tutaj. Ta warstwa mierzy też prawdziwy rozmiar wyniku — deklaracja
//     w nagłówku archiwum jest cudzą obietnicą.
//
// Dopiero po trzeciej warstwie treść wchodzi do katalogu Operatora. Odmowa na
// każdej z nich zostawia jego drzewo nietknięte, bo do tej chwili nic w nim nie
// powstało.
//
// Miejsca docelowe są dwa, bo zamiary są dwa. Kontrakt opisuje `targetPath`
// słowami „katalog docelowy w katalogu roboczym okna; brak kładzie zawartość
// w magazynie": ze ścieżką model rozpakowuje po to, żeby na tych plikach
// pracować, bez ścieżki — żeby zawartość przechować pod sumami kontrolnymi, nie
// zaśmiecając katalogu Operatora.
//
// Plik istniejący nie jest nadpisywany: przy przenoszeniu z kwarantanny nazwa
// zajęta w celu jest odmową. Nadpisanie cudzej pracy zawartością archiwum jest
// tą samą szkodą, przed którą stoi reszta tego pliku.
package core

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// Rozpakuj obsługuje `archive.unpack`.
func (a *adapterNarzedziArchiwum) Rozpakuj(ctx context.Context,
	z shared.ArchiveUnpackRequest) (shared.ArchiveUnpackResponse, error) {

	katalogOkna, err := a.katalogRoboczyOkna()
	if err != nil {
		return shared.ArchiveUnpackResponse{}, err
	}
	zrodlo, err := a.zrodloRozpakowania(ctx, z, katalogOkna)
	if err != nil {
		return shared.ArchiveUnpackResponse{}, err
	}

	// Cel rozstrzygamy przed rozpakowaniem, bo od niego zależy, na którym
	// nośniku ma stanąć kwarantanna — a przeniesienie między nośnikami nie jest
	// przemianowaniem.
	var cel string
	doMagazynu := bezWartosci(z.TargetPath)
	podstawaKwarantanny := a.katalogDanych
	if !doMagazynu {
		cel, err = sciezkaWzgledemKatalogu(katalogOkna, *z.TargetPath)
		if err != nil {
			return shared.ArchiveUnpackResponse{}, err
		}
		podstawaKwarantanny = katalogOkna
	}

	kwarantanna, err := katalogRoboczyTymczasowy(podstawaKwarantanny, ".danaco-rozpakowanie-*")
	if err != nil {
		return shared.ArchiveUnpackResponse{}, err
	}
	defer sprzatnij(kwarantanna)

	tresc, err := a.rozpakujDoKwarantanny(ctx, zrodlo, kwarantanna)
	if err != nil {
		return shared.ArchiveUnpackResponse{}, err
	}
	sciezki, err := zweryfikujKwarantanne(tresc)
	if err != nil {
		return shared.ArchiveUnpackResponse{}, err
	}

	if doMagazynu {
		liczba, err := a.wniesDoMagazynu(ctx, tresc, sciezki)
		if err != nil {
			return shared.ArchiveUnpackResponse{}, err
		}
		// Pole `paths` zostaje puste. Zawartość wylądowała w magazynie
		// rdzenia, a tam nie ma ścieżek — są zasoby pod identyfikatorami.
		// Wypisanie tu nazw z wnętrza archiwum dałoby napisy wyglądające jak
		// ścieżki na dysku Operatora, pod którymi nie ma nic.
		// Pole jest w kontrakcie nieobowiązkowe właśnie dlatego, że jedna
		// z dwóch dróg tej komendy ścieżek nie wytwarza.
		return shared.ArchiveUnpackResponse{Entries: liczba}, nil
	}

	wydane, err := przeniesDoCelu(tresc, cel, sciezki, katalogOkna)
	if err != nil {
		return shared.ArchiveUnpackResponse{}, err
	}
	return shared.ArchiveUnpackResponse{Entries: len(wydane), Paths: wydane}, nil
}

// zrodloRozpakowania przekłada parę `assetId?|sourcePath?` na plik archiwum.
//
// Zasób ma pierwszeństwo przed ścieżką — treść zasobu leży pod sumą kontrolną
// i nie zmieni się między wskazaniem a odczytem. Podmiana pliku między
// obejrzeniem spisu a rozpakowaniem byłaby zresztą sposobem obejścia wyroku,
// więc droga przez magazyn jest tu nie tylko wygodniejsza, ale i pewniejsza.
func (a *adapterNarzedziArchiwum) zrodloRozpakowania(ctx context.Context,
	z shared.ArchiveUnpackRequest, katalogOkna string) (string, error) {

	if !bezWartosci(z.AssetId) {
		_, blob, err := a.zasobArchiwum(ctx, *z.AssetId)
		return blob, err
	}
	if !bezWartosci(z.SourcePath) {
		pelna, err := sciezkaWzgledemKatalogu(katalogOkna, *z.SourcePath)
		if err != nil {
			return "", err
		}
		opis, err := os.Stat(pelna)
		if err != nil {
			return "", bladWskazaniaArchiwum("nie można odczytać archiwum " +
				strings.TrimSpace(*z.SourcePath) + ": " + err.Error())
		}
		if opis.IsDir() {
			return "", bladWskazaniaArchiwum(strings.TrimSpace(*z.SourcePath) +
				" jest katalogiem, a nie archiwum")
		}
		return pelna, nil
	}
	return "", bladWskazaniaArchiwum(
		"żądanie bez wskazania archiwum: brakuje pola assetId albo sourcePath — " +
			"narzędzie archiwum nie zgaduje, co ma otworzyć")
}

// rozpakujDoKwarantanny odwija archiwum i oddaje katalog z jego zawartością.
//
// `tar.gz` wymaga dwóch przebiegów: `7z l` na `.tar.gz` widzi jedną pozycję —
// zawinięty `.tar` — więc wyrok wydany na tym spisie nie oglądałby zawartości
// w ogóle. Najpierw
// zdejmujemy warstwę `gzip`, potem oglądamy i rozpakowujemy warstwę `tar`.
// Zdjęcie pierwszej warstwy jest bezpieczne bez wyroku, bo `gzip` nie niesie
// ścieżek: to jeden strumień, który ląduje jednym plikiem w kwarantannie.
func (a *adapterNarzedziArchiwum) rozpakujDoKwarantanny(ctx context.Context,
	zrodlo, kwarantanna string) (string, error) {

	format, err := rozpoznajFormatArchiwum(zrodlo)
	if err != nil {
		return "", err
	}

	doOtwarcia := zrodlo
	if format == formatTarGz {
		warstwa := filepath.Join(kwarantanna, "warstwa")
		if err := a.wolajRozpakowanie(ctx, zrodlo, warstwa); err != nil {
			return "", err
		}
		doOtwarcia, err = jedynyPlik(warstwa)
		if err != nil {
			return "", err
		}
	}

	spis, err := a.spisArchiwum(ctx, doOtwarcia)
	if err != nil {
		return "", err
	}
	if err := sprawdzSpisArchiwum(spis); err != nil {
		return "", err
	}

	tresc := filepath.Join(kwarantanna, "tresc")
	if err := a.wolajRozpakowanie(ctx, doOtwarcia, tresc); err != nil {
		return "", err
	}
	return tresc, nil
}

// wolajRozpakowanie składa jedno wywołanie `7z x`.
//
// Tryb `x`, a nie `e`: `x` zachowuje strukturę katalogów, `e` spłaszcza
// wszystko do jednego poziomu. Spłaszczenie wyglądałoby na obronę przed
// ucieczką ze ścieżki, ale jest gorsze niż odmowa: dwa pliki o tej samej nazwie
// z różnych katalogów nadpisałyby się nawzajem, a Operator dostałby drzewo inne
// niż to, które przysłał, i nie dowiedziałby się o tym.
func (a *adapterNarzedziArchiwum) wolajRozpakowanie(ctx context.Context, archiwum, cel string) error {
	if err := os.MkdirAll(cel, 0o755); err != nil {
		return bladZapleczaArchiwum("nie można założyć katalogu kwarantanny: " + err.Error())
	}
	_, err := a.wolaj7z(ctx, []string{"x", "-bd", "-y", "-o" + cel, "--", archiwum},
		filepath.Dir(archiwum))
	return err
}

// jedynyPlik oddaje jedyny plik katalogu. Warstwa `gzip` niesie dokładnie jeden
// strumień, więc cokolwiek innego niż jeden plik znaczy, że rozpakowanie poszło
// inaczej, niż zakładamy — a wtedy lepiej odmówić niż zgadywać, który wziąć.
func jedynyPlik(katalog string) (string, error) {
	wpisy, err := os.ReadDir(katalog)
	if err != nil {
		return "", bladZapleczaArchiwum("nie można odczytać kwarantanny: " + err.Error())
	}
	var znalezione []string
	for _, wpis := range wpisy {
		if !wpis.IsDir() {
			znalezione = append(znalezione, filepath.Join(katalog, wpis.Name()))
		}
	}
	if len(znalezione) != 1 {
		return "", bladPrzetwarzaniaArchiwum(
			"zdjęcie warstwy gzip nie dało jednego pliku — archiwum nie ma spodziewanego kształtu tar.gz")
	}
	return znalezione[0], nil
}

// zweryfikujKwarantanne obchodzi rozpakowaną treść i oddaje ścieżki plików
// względem korzenia kwarantanny.
//
// Rozmiar liczymy z tego, co legło na nośniku, a nie z nagłówka archiwum:
// nagłówek jest obietnicą jego twórcy, a bomba dekompresyjna obiecuje, co
// zechce. Dowiązania i pliki nie-zwykłe (gniazda, urządzenia, potoki nazwane)
// odmawiamy tutaj po typie wpisu, bo `WalkDir` czyta go bez wchodzenia w to, na
// co wpis wskazuje.
func zweryfikujKwarantanne(korzen string) ([]string, error) {
	var sciezki []string
	var suma int64

	err := filepath.WalkDir(korzen, func(sciezka string, wpis fs.DirEntry, err error) error {
		if err != nil {
			return bladZapleczaArchiwum("nie można obejść rozpakowanej treści: " + err.Error())
		}
		wzgledna, blad := filepath.Rel(korzen, sciezka)
		if blad != nil || wzgledna == ".." || strings.HasPrefix(wzgledna, ".."+string(filepath.Separator)) {
			return bladSciezkiPozaKatalogiem(sciezka, "wylądowała poza katalogiem kwarantanny")
		}
		if wzgledna == "." || wpis.IsDir() {
			return nil
		}
		if wpis.Type()&os.ModeSymlink != 0 {
			return bladWskazaniaArchiwum("rozpakowana pozycja " + wzgledna +
				" okazała się dowiązaniem — dowiązanie wyprowadza zapis poza katalog docelowy; " +
				"naprawa: przysłać archiwum z samymi plikami i katalogami")
		}
		if !wpis.Type().IsRegular() {
			return bladWskazaniaArchiwum("rozpakowana pozycja " + wzgledna +
				" nie jest zwykłym plikiem — rdzeń wydaje Operatorowi wyłącznie pliki i katalogi")
		}
		opis, blad := wpis.Info()
		if blad != nil {
			return bladZapleczaArchiwum("nie można zmierzyć " + wzgledna + ": " + blad.Error())
		}
		suma += opis.Size()
		if suma > granicaRozpakowaniaBajty {
			return bladPrzekroczonejGranicyRozmiaru(suma)
		}
		if len(sciezki) >= granicaRozpakowaniaPozycji {
			return bladWskazaniaArchiwum("rozpakowana treść przekroczyła granicę " +
				"liczby pozycji — rdzeń rozpakowuje wyłącznie do tej granicy")
		}
		sciezki = append(sciezki, wzgledna)
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(sciezki) == 0 {
		return nil, bladWskazaniaArchiwum("archiwum nie zawierało ani jednego pliku")
	}
	return sciezki, nil
}

// przeniesDoCelu przenosi zawartość kwarantanny do katalogu docelowego i oddaje
// ścieżki wydane Operatorowi — liczone względem KATALOGU ROBOCZEGO OKNA, bo to
// jest układ odniesienia, w którym model i Operator się umawiają. Ścieżka
// bezwzględna wypisywałaby w odpowiedzi układ katalogów maszyny, o który nikt
// nie pytał.
//
// KOLEJNOŚĆ JEST ZAMIERZONA: najpierw sprawdzamy WSZYSTKIE kolizje, potem
// przenosimy. Przenoszenie z jednoczesnym sprawdzaniem zostawiłoby przy
// kolizji połowę treści w celu i połowę w kwarantannie — czyli rozpakowanie
// częściowe podane jako odmowa.
func przeniesDoCelu(tresc, cel string, sciezki []string, katalogOkna string) ([]string, error) {
	if err := os.MkdirAll(cel, 0o755); err != nil {
		return nil, bladZapleczaArchiwum("nie można założyć katalogu docelowego: " + err.Error())
	}

	wpisy, err := os.ReadDir(tresc)
	if err != nil {
		return nil, bladZapleczaArchiwum("nie można odczytać rozpakowanej treści: " + err.Error())
	}
	for _, wpis := range wpisy {
		if _, err := os.Lstat(filepath.Join(cel, wpis.Name())); err == nil {
			return nil, bladWskazaniaArchiwum("w katalogu docelowym leży już " + wpis.Name() +
				" — rozpakowanie nadpisałoby pracę, której to archiwum nie dotyczy; " +
				"naprawa: wskazać pusty katalog docelowy albo usunąć to, co tam leży")
		}
	}
	for _, wpis := range wpisy {
		if err := os.Rename(filepath.Join(tresc, wpis.Name()), filepath.Join(cel, wpis.Name())); err != nil {
			return nil, bladZapleczaArchiwum("nie można wydać " + wpis.Name() +
				" do katalogu docelowego: " + err.Error())
		}
	}

	wydane := make([]string, 0, len(sciezki))
	for _, wzgledna := range sciezki {
		odOkna, err := filepath.Rel(katalogOkna, filepath.Join(cel, wzgledna))
		if err != nil {
			odOkna = wzgledna
		}
		wydane = append(wydane, filepath.ToSlash(odOkna))
	}
	return wydane, nil
}

// wniesDoMagazynu odkłada każdy rozpakowany plik jako osobny zasób — droga dla
// żądania BEZ `targetPath`.
//
// PO CO TA DROGA W OGÓLE ISTNIEJE: model, który dostał archiwum i chce się
// dowiedzieć, co w nim jest, nie musi zaśmiecać katalogu roboczego Operatora.
// Zasoby ma pod identyfikatorami i sięga po nie pozostałymi narzędziami —
// obrazu, dokumentu, mediów.
func (a *adapterNarzedziArchiwum) wniesDoMagazynu(ctx context.Context,
	tresc string, sciezki []string) (int, error) {

	if a.magazyn == nil || a.repozytorium == nil {
		return 0, bladZapleczaArchiwum(
			"magazyn zasobów nie jest wpięty — nie ma gdzie odłożyć rozpakowanej treści")
	}

	for _, wzgledna := range sciezki {
		odwolanie, _, _, err := a.magazyn.ZapiszZePliku(filepath.Join(tresc, wzgledna))
		if err != nil {
			return 0, bladZapleczaArchiwum("nie można utrwalić " + wzgledna + ": " + err.Error())
		}
		// Nazwą zasobu jest ścieżka WEWNĄTRZ archiwum, a nie sama nazwa pliku:
		// archiwum z dziesięcioma plikami `index.html` w różnych katalogach
		// dałoby inaczej dziesięć kafelków nie do odróżnienia.
		zasob := dane.ZasobDesignu{
			Kod:    nowyIdentyfikator(przedrostekZasobuDesign),
			Okno:   oknoZasobowNarzedziArchiwum,
			Nazwa:  wskaznikTekstu(filepath.ToSlash(wzgledna)),
			Rodzaj: shared.DesignAssetKindImage,
			URI:    &odwolanie,
		}
		if rozszerzenie := strings.TrimPrefix(filepath.Ext(wzgledna), "."); rozszerzenie != "" {
			zasob.Format = wskaznikTekstu(strings.ToLower(rozszerzenie))
		}
		if _, err := a.repozytorium.ZapiszZasob(ctx, zasob); err != nil {
			return 0, bladZapleczaArchiwum("nie można założyć wiersza zasobu dla " +
				wzgledna + ": " + err.Error())
		}
	}
	return len(sciezki), nil
}
