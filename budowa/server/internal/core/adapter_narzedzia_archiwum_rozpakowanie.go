// Adapter obsługuje `archive.unpack`: rozpakowuje archiwum do katalogu roboczego okna albo
// do magazynu zasobów, trzema warstwami obrony przed wyjściem poza katalog docelowy — spisem
// przed zapisem, kwarantanną i przejściem po wyniku.
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

// Rozpakuj obsługuje komendę `archive.unpack`: rozstrzyga cel, prowadzi treść przez
// kwarantannę i wydaje ją do katalogu docelowego albo do magazynu zasobów.
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

	// Cel rozstrzyga się przed rozpakowaniem, bo od niego zależy nośnik kwarantanny.
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
		// Pole paths zostaje puste: zawartość leży w magazynie pod identyfikatorami, nie pod ścieżkami.
		return shared.ArchiveUnpackResponse{Entries: liczba}, nil
	}

	wydane, err := przeniesDoCelu(tresc, cel, sciezki, katalogOkna)
	if err != nil {
		return shared.ArchiveUnpackResponse{}, err
	}
	return shared.ArchiveUnpackResponse{Entries: len(wydane), Paths: wydane}, nil
}

// zrodloRozpakowania przekłada parę `assetId?|sourcePath?` na plik archiwum. Zasób ma
// pierwszeństwo przed ścieżką, bo jego treść leży pod sumą kontrolną i nie zmieni się między
// wskazaniem a odczytem.
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

// rozpakujDoKwarantanny odwija archiwum i oddaje katalog z jego zawartością; archiwum tar.gz
// przechodzi dwa przebiegi, bo spis ogląda jedną zawiniętą pozycję zamiast treści warstwy tar.
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

// wolajRozpakowanie składa jedno wywołanie 7z w trybie x, który zachowuje strukturę
// katalogów; tryb e spłaszczałby wynik i nadpisywał pliki o tej samej nazwie z różnych
// katalogów.
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
// inaczej, niż zaklada się — a wtedy lepiej odmówić niż zgadywać, który wziąć.
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

// zweryfikujKwarantanne obchodzi rozpakowaną treść i oddaje ścieżki plików względem korzenia
// kwarantanny; rozmiar liczy z tego, co legło na nośniku, a dowiązania i pliki nie-zwykłe
// odmawia po typie wpisu.
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

// przeniesDoCelu przenosi zawartość kwarantanny do katalogu docelowego i oddaje ścieżki
// wydane Operatorowi względem katalogu roboczego okna; najpierw sprawdza wszystkie kolizje,
// potem przenosi, aby uniknąć rozpakowania częściowego.
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

// wniesDoMagazynu odkłada każdy rozpakowany plik jako osobny zasób pod identyfikatorem —
// droga dla żądania bez pola targetPath, żeby model mógł obejrzeć zawartość archiwum bez
// zaśmiecania katalogu roboczego Operatora.
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
		// Nazwą zasobu jest ścieżka wewnątrz archiwum, nie nazwa pliku, by uniknąć nierozróżnialnych kafelków.
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
