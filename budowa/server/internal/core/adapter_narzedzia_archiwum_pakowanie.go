// Odpowiedzialność pliku: komenda `archive.pack` — złożenie wskazanych zasobów
// albo katalogu roboczego w jedno archiwum i odłożenie go jako zasobu.
// Zaplecze rodziny leży w `adapter_narzedzia_archiwum.go`, czytanie spisu
// w `adapter_narzedzia_archiwum_spis.go`.
//
// ── sourcePath rozstrzygamy względem katalogu okna ─────────────────────────
// Pakowanie wygląda na czynność czytającą, więc ścieżka bezwzględna zdawałaby
// się tu nieszkodliwa. Nie jest. Wynikiem `archive.pack` jest zasób w magazynie
// rdzenia, a zasób model potrafi odczytać i rozpakować. Przyjęcie ścieżki
// bezwzględnej dałoby więc drogę: spakuj `~/.ssh`, odłóż jako zasób, rozpakuj
// u siebie — czyli wyniesienie dowolnego pliku maszyny Operatora do materiału,
// którym model dysponuje. Katalog roboczy okna jest jedynym miejscem, o którym
// produkt umówił się z Operatorem, że model tam sięga, i pakowanie tej umowy
// nie łamie.
//
// ── zasoby idą przez rusztowanie, a nie wprost do `7z` ─────────────────────
// Blob magazynu nazywa się swoją sumą kontrolną. Podanie blobów `7z` wprost
// dałoby archiwum, w którym pliki nazywają się
// `3f786850e387550fdab836ed7e6dc881de23001b` — Operator dostałby dwadzieścia
// nazw, z których żadna nic nie znaczy. Rusztowanie jest katalogiem
// tymczasowym, w którym każdy blob dostaje swoją nazwę czytelną z wiersza
// zasobu; dopiero ono trafia do `7z`.
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

// Spakuj obsługuje `archive.pack`.
//
// Porządek jest zamierzony: najpierw rusztowanie, potem archiwum, potem spis
// gotowego archiwum, a dopiero na końcu magazyn. `entries` liczymy z tego, co
// w archiwum naprawdę jest, a nie z tego, ile pozycji kazaliśmy spakować —
// liczba wzięta z żądania byłaby w odpowiedzi nieodróżnialna od policzonej,
// a rozjechałaby się z prawdą przy każdym pliku, którego `7z` nie wziął.
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

	// Praca pośrednia stoi przy katalogu danych rdzenia, a nie w katalogu
	// roboczym okna: archiwum w budowie nie jest pracą Operatora i nie ma
	// zaśmiecać katalogu, który on ogląda.
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

// zrodloPakowania rozstrzyga, co pakujemy, i oddaje katalog, z którego woła się
// `7z`, wraz z listą pozycji do wzięcia oraz nazwą domyślną archiwum.
//
// Zasoby mają pierwszeństwo przed ścieżką, tak samo i z tego samego powodu co
// w rodzinie obrazu: treść zasobu leży już pod sumą kontrolną i nie zmieni się
// między wskazaniem a odczytem, a plik na dysku jest treścią żywą.
//
// Żądanie bez jednego i drugiego jest odmową. Podstawienie „katalogu roboczego,
// skoro nic nie podał" byłoby zgadywaniem przedmiotu czynności — a spakowanie
// całego katalogu roboczego zamiast dwóch plików jest pomyłką kosztowną.
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
		// Wołamy `7z` z katalogu nadrzędnego i podajemy samą nazwę. Dzięki temu
		// archiwum niesie `projekt/plik.txt`, a nie ścieżkę od korzenia maszyny —
		// a przy rozpakowaniu powstaje jeden katalog, nie rozsypane pliki.
		return filepath.Dir(pelna), []string{filepath.Base(pelna)}, filepath.Base(pelna), nil
	}
	return "", nil, "", bladWskazaniaArchiwum(
		"żądanie bez wskazania, co spakować: brakuje pola assetIds albo sourcePath — " +
			"narzędzie archiwum nie zgaduje przedmiotu pakowania")
}

// rusztowanieZasobow stawia katalog tymczasowy, w którym każdy wskazany zasób
// leży pod swoją nazwą czytelną, i oddaje listę tych nazw.
//
// Kopiujemy, a nie dowiązujemy. Dowiązanie do bloba oszczędziłoby bajty, ale
// `7z` zapisałby wtedy do archiwum dowiązanie albo poszedł za nim zależnie od
// przełącznika — a przy rozpakowaniu ta sama rodzina dowiązania odmawia.
// Archiwum, którego własny produkt nie umie otworzyć, byłoby rozjazdem
// wewnątrz jednej rodziny.
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

// nazwaZasobuWArchiwum składa nazwę, pod którą zasób wejdzie do archiwum.
//
// Nazwa musi spełnić trzy warunki naraz. Musi być czytelna (więc bierzemy
// `nazwa` wiersza, a `kod` dopiero w zastępstwie), musi być samą nazwą (więc
// obcinamy katalogi i człony `..` — nazwa zasobu pochodzi od Operatora
// i mogłaby nieść ścieżkę, a archiwum z pozycją `../coś` odmówiłaby własna
// rodzina przy rozpakowaniu) i musi być jedyna (dwa zasoby o tej samej nazwie
// nadpisałyby się nawzajem w rusztowaniu i jeden zniknąłby bez słowa).
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
	// Rozszerzenie z pola `format`, gdy nazwa go nie niesie — Operator
	// rozpakowujący archiwum ma dostać plik, który jego system otworzy.
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

// nazwaArchiwum rozstrzyga nazwę pliku archiwum — bez rozszerzenia, bo to
// dokłada format.
//
// Nazwa od modelu jest obcinana do samej nazwy z tych samych powodów, co nazwa
// zasobu wyżej: `name` postaci `../../wydanie` kazałoby zapisać archiwum poza
// katalogiem pracy.
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
	// Rozszerzenie dokłada format, więc zdejmujemy to, które model dopisał sam —
	// inaczej powstałoby `wydanie.zip.zip`.
	for _, koncowka := range []string{".zip", ".7z", ".tar.gz", ".tgz", ".tar"} {
		podstawa = strings.TrimSuffix(podstawa, koncowka)
	}
	if podstawa == "" {
		return "archiwum"
	}
	return podstawa
}

// zbudujArchiwum woła `7z` i oddaje ścieżkę gotowego archiwum oraz ścieżkę
// pliku, z którego czyta się spis jego zawartości.
//
// Dla `tar.gz` są to dwa różne pliki i dwa wywołania. `tar.gz` jest archiwum
// dwuwarstwowym: `tar` niesie pozycje, `gzip` opakowuje całość jednym
// strumieniem, więc `7z l` na gotowym `.tar.gz` pokazuje jedną pozycję —
// zawinięty plik `.tar` — a nie zawartość. Spis czytamy zatem z warstwy `tar`,
// zanim ją zawiniemy, bo tam pozycje naprawdę są. Dwa wywołania zamiast potoku,
// bo `zewnetrzne.Wolaj` prowadzi jeden proces — i dobrze, bo potok dwóch
// programów miałby dwa kody wyjścia i jedną odpowiedź o powodzeniu.
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
		// Drugie wywołanie bierze gotowy `.tar` jako jedyną pozycję i woła się
		// z jego katalogu, żeby do środka nie weszła ścieżka bezwzględna.
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

// wolajPakowanie składa jedno wywołanie `7z a`.
//
// Przełączniki są wybrane, nie przepisane: `-bd` zdejmuje pasek postępu (w
// strumieniu bez terminala byłby śmieciem w diagnostyce), `-y` odpowiada „tak"
// na pytania programu (proces bez terminala nie ma komu odpowiadać i czekałby
// do granicy czasu), a `--` zamyka listę przełączników — bez niego plik
// nazwany `-sdel` zostałby wzięty za polecenie, a nie za nazwę.
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
