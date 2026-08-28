// Plik niesie komendę archive.pack: złożenie wskazanych zasobów albo katalogu roboczego w jedno archiwum i odłożenie go jako zasobu. Zaplecze rodziny leży w adapter_narzedzia_archiwum.go, czytanie spisu w adapter_narzedzia_archiwum_spis.go.
package core

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"danacoconsole/shared"
)

// Spakuj obsługuje archive.pack. Porządek jest zamierzony: najpierw rusztowanie, potem archiwum, potem spis gotowego archiwum, a dopiero na końcu magazyn.
func (a *adapterNarzedziArchiwum) Spakuj(ctx context.Context,
	z shared.ArchivePackRequest) (shared.ArchivePackResponse, error) {

	format, err := rozstrzygnijFormatDocelowy(z.Format)
	if err != nil {
		return shared.ArchivePackResponse{}, err
	}
	katalogOkna, err := a.katalogRoboczyOkna()
	if err != nil {
		return shared.ArchivePackResponse{}, err
	}

	// Praca pośrednia stoi przy katalogu danych rdzenia, nie w katalogu roboczym, który ogląda Operator.
	praca, err := katalogRoboczyTymczasowy(a.katalogDanych, "pakowanie-*")
	if err != nil {
		return shared.ArchivePackResponse{}, err
	}
	defer sprzatnij(praca)

	katalogZrodla, pozycje, nazwaDomyslna, err := a.zrodloPakowania(ctx, z, katalogOkna, praca)
	if err != nil {
		return shared.ArchivePackResponse{}, err
	}

	nazwa := nazwaArchiwum(z.Name, nazwaDomyslna)
	sciezkaArchiwum, sciezkaSpisu, err := a.zbudujArchiwum(ctx, katalogZrodla, pozycje,
		filepath.Join(praca, nazwa), format)
	if err != nil {
		return shared.ArchivePackResponse{}, err
	}

	spis, err := a.spisArchiwum(ctx, sciezkaSpisu)
	if err != nil {
		return shared.ArchivePackResponse{}, err
	}

	zasob, rozmiar, err := a.odlozArchiwumJakoZasob(ctx, sciezkaArchiwum,
		nazwa+"."+string(format), format)
	if err != nil {
		return shared.ArchivePackResponse{}, err
	}

	return shared.ArchivePackResponse{
		Asset:     zasob,
		Entries:   policzPliki(spis),
		SizeBytes: int(rozmiar),
	}, nil
}

// zrodloPakowania rozstrzyga, co pakujemy, i oddaje katalog, z którego woła się 7z, wraz z listą pozycji do wzięcia oraz nazwą domyślną archiwum. Zasoby mają pierwszeństwo przed ścieżką.
func (a *adapterNarzedziArchiwum) zrodloPakowania(ctx context.Context,
	z shared.ArchivePackRequest, katalogOkna, praca string) (string, []string, string, error) {

	if len(z.AssetIds) > 0 {
		katalog, pozycje, err := a.rusztowanieZasobow(ctx, z.AssetIds, praca)
		return katalog, pozycje, "zasoby", err
	}
	if !bezWartosci(z.SourcePath) {
		pelna, err := sciezkaWzgledemKatalogu(katalogOkna, *z.SourcePath)
		if err != nil {
			return "", nil, "", err
		}
		if _, err := os.Stat(pelna); err != nil {
			return "", nil, "", bladWskazaniaArchiwum("nie można odczytać " +
				strings.TrimSpace(*z.SourcePath) + ": " + err.Error())
		}
		// Woła się 7z z katalogu nadrzędnego i podaje się samą nazwę, więc archiwum niesie ścieżkę względną.
		return filepath.Dir(pelna), []string{filepath.Base(pelna)}, filepath.Base(pelna), nil
	}
	return "", nil, "", bladWskazaniaArchiwum(
		"żądanie bez wskazania, co spakować: brakuje pola assetIds albo sourcePath — " +
			"narzędzie archiwum nie zgaduje przedmiotu pakowania")
}

// rusztowanieZasobow stawia katalog tymczasowy, w którym każdy wskazany zasób leży pod swoją nazwą czytelną, i oddaje listę tych nazw. Kopiuje się, a nie dowiazuje się, bo 7z traktuje dowiązania niejednoznacznie.
func (a *adapterNarzedziArchiwum) rusztowanieZasobow(ctx context.Context,
	kody []string, praca string) (string, []string, error) {

	rusztowanie := filepath.Join(praca, "rusztowanie")
	if err := os.MkdirAll(rusztowanie, 0o755); err != nil {
		return "", nil, bladZapleczaArchiwum("nie można założyć katalogu rusztowania: " + err.Error())
	}

	zajete := map[string]bool{}
	var nazwy []string
	for _, kod := range kody {
		zasob, blob, err := a.zasobArchiwum(ctx, kod)
		if err != nil {
			return "", nil, err
		}
		nazwa := nazwaZasobuWArchiwum(zasob.Nazwa, zasob.Kod, zasob.Format, zajete)
		if err := skopiujPlik(blob, filepath.Join(rusztowanie, nazwa)); err != nil {
			return "", nil, err
		}
		nazwy = append(nazwy, nazwa)
	}
	return rusztowanie, nazwy, nil
}

// nazwaZasobuWArchiwum składa nazwę, pod którą zasób wejdzie do archiwum. Nazwa musi być czytelna, samą nazwą bez katalogów i człony `..`, oraz jedyna w obrębie archiwum.
func nazwaZasobuWArchiwum(nazwa *string, kod string, format *string, zajete map[string]bool) string {
	podstawa := ""
	if nazwa != nil {
		podstawa = strings.TrimSpace(*nazwa)
	}
	podstawa = nazwaBezKatalogow(podstawa)
	podstawa = strings.ReplaceAll(podstawa, "..", "")
	podstawa = strings.TrimSpace(podstawa)
	if podstawa == "" || podstawa == "." {
		podstawa = kod
	}
	// Rozszerzenie z pola format wchodzi, gdy nazwa własnego rozszerzenia nie niesie.
	if format != nil && !strings.Contains(podstawa, ".") {
		if rozszerzenie := strings.TrimSpace(*format); rozszerzenie != "" {
			podstawa += "." + rozszerzenie
		}
	}

	nazwaKoncowa := podstawa
	for licznik := 2; zajete[nazwaKoncowa]; licznik++ {
		rozszerzenie := filepath.Ext(podstawa)
		rdzen := strings.TrimSuffix(podstawa, rozszerzenie)
		nazwaKoncowa = rdzen + "-" + strconv.Itoa(licznik) + rozszerzenie
	}
	zajete[nazwaKoncowa] = true
	return nazwaKoncowa
}

// nazwaArchiwum rozstrzyga nazwę pliku archiwum, bez rozszerzenia, bo to dokłada format; nazwa od modelu jest obcinana do samej nazwy.
func nazwaArchiwum(zadana *string, domyslna string) string {
	podstawa := ""
	if !bezWartosci(zadana) {
		podstawa = strings.TrimSpace(*zadana)
	}
	podstawa = nazwaBezKatalogow(podstawa)
	podstawa = strings.ReplaceAll(podstawa, "..", "")
	podstawa = strings.TrimSpace(podstawa)
	if podstawa == "" {
		podstawa = strings.TrimSpace(domyslna)
	}
	if podstawa == "" || podstawa == "." {
		return "archiwum"
	}
	// Rozszerzenie dokłada format, więc zdejmuje się to, które model dopisał sam w nazwie.
	for _, koncowka := range []string{".zip", ".7z", ".tar.gz", ".tgz", ".tar"} {
		podstawa = strings.TrimSuffix(podstawa, koncowka)
	}
	if podstawa == "" {
		return "archiwum"
	}
	return podstawa
}

// zbudujArchiwum woła 7z i oddaje ścieżkę gotowego archiwum oraz ścieżkę pliku, z którego czyta się spis jego zawartości. Dla tar.gz są to dwa różne pliki i dwa wywołania.
func (a *adapterNarzedziArchiwum) zbudujArchiwum(ctx context.Context,
	katalogZrodla string, pozycje []string, podstawaWyjscia string,
	format formatArchiwum) (string, string, error) {

	if len(pozycje) == 0 {
		return "", "", bladWskazaniaArchiwum("nie wskazano ani jednej pozycji do spakowania")
	}

	switch format {
	case formatTarGz:
		sciezkaTar := podstawaWyjscia + ".tar"
		if err := a.wolajPakowanie(ctx, katalogZrodla, "-ttar", sciezkaTar, pozycje); err != nil {
			return "", "", err
		}
		sciezkaGz := sciezkaTar + ".gz"
		// Drugie wywołanie bierze gotowy plik tar jako jedyną pozycję do spakowania warstwą gzip.
		if err := a.wolajPakowanie(ctx, filepath.Dir(sciezkaTar), "-tgzip", sciezkaGz,
			[]string{filepath.Base(sciezkaTar)}); err != nil {
			return "", "", err
		}
		return sciezkaGz, sciezkaTar, nil
	case format7z:
		sciezka := podstawaWyjscia + ".7z"
		err := a.wolajPakowanie(ctx, katalogZrodla, "-t7z", sciezka, pozycje)
		return sciezka, sciezka, err
	default:
		sciezka := podstawaWyjscia + ".zip"
		err := a.wolajPakowanie(ctx, katalogZrodla, "-tzip", sciezka, pozycje)
		return sciezka, sciezka, err
	}
}

// wolajPakowanie składa jedno wywołanie 7z a, z przełącznikami dobranymi pod pracę bez terminala i bez dwuznaczności nazw pozycji.
func (a *adapterNarzedziArchiwum) wolajPakowanie(ctx context.Context,
	katalog, typ, wyjscie string, pozycje []string) error {

	argumenty := append([]string{"a", typ, "-bd", "-y", "--", wyjscie}, pozycje...)
	if _, err := a.wolaj7z(ctx, argumenty, katalog); err != nil {
		return err
	}
	if _, err := os.Stat(wyjscie); err != nil {
		return bladPrzetwarzaniaArchiwum(
			"7-Zip zakończył się powodzeniem, ale archiwum nie powstało: " + err.Error())
	}
	return nil
}

// skopiujPlik przepisuje blob do rusztowania strumieniem — plik dowolnej
// wielkości przechodzi tak samo i nie ląduje w całości w pamięci rdzenia.
func skopiujPlik(zrodlo, cel string) error {
	wejscie, err := os.Open(zrodlo)
	if err != nil {
		return bladZapleczaArchiwum("nie można odczytać treści zasobu: " + err.Error())
	}
	defer wejscie.Close()

	wyjscie, err := os.Create(cel)
	if err != nil {
		return bladZapleczaArchiwum("nie można założyć pliku w rusztowaniu: " + err.Error())
	}
	if _, err := io.Copy(wyjscie, wejscie); err != nil {
		wyjscie.Close()
		return bladZapleczaArchiwum("nie można przepisać treści zasobu: " + err.Error())
	}
	if err := wyjscie.Close(); err != nil {
		return bladZapleczaArchiwum("nie można domknąć pliku w rusztowaniu: " + err.Error())
	}
	return nil
}
