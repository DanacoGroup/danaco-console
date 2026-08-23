// Moduł Library — higiena repozytorium: `library.duplicate.scan`,
// `library.duplicate.merge`, `library.fixity.check`, `library.name.normalize`,
// `library.stats.get`.
//
// Wszystkie pięć czynności są czynnościami RDZENIA, nie okna, i z jednego
// powodu: sięgają po to, czego okno nie ma. Rozpoznanie przybliżone żąda
// porównania każdego zasobu z każdym, weryfikacja integralności — przeliczenia
// sumy z bajtów leżących na dysku, pulpit stanu — przejścia po całym zbiorze.
// Okno liczące te rzeczy z odczytanej strony wykazu orzeka o próbce.
//
// Rozpoznanie dokładne idzie po sumie kontrolnej i jest pewne. Rozpoznania
// przybliżone niosą trafność, bo to wnioski: podobieństwo treści liczy się
// odciskiem słów (shingling), podobieństwo obrazu — odciskiem percepcyjnym
// zbudowanym z obrazu zeskalowanego do siatki. Oba liczy wkompilowany kod Go;
// żaden nie woła programu z zewnątrz.
package core

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// granicaZbioruHigieny jest górną granicą zbioru branego do czynności
// obchodzących całe repozytorium. Bez niej repozytorium o setkach tysięcy
// zasobów wciągnęłoby cały wykaz do pamięci przy jednym żądaniu.
const granicaZbioruHigieny = 5000

// progPodobienstwaDomyslny — trafność, od której dwa zasoby uznaje się za
// niemal-duplikaty. Wartość w setnych, zgodnie z kontraktem.
const progPodobienstwaDomyslny = 85

// SkanujDuplikaty obsługuje `library.duplicate.scan`.
func (a *adapterBiblioteki) SkanujDuplikaty(ctx context.Context,
	z shared.LibraryDuplicateScanRequest) (shared.LibraryDuplicateScanResponse, error) {

	zasoby, err := a.zasobyZbioru(ctx, nil, z.CollectionId)
	if err != nil {
		return shared.LibraryDuplicateScanResponse{}, err
	}
	rodzaje := z.Kinds
	if len(rodzaje) == 0 {
		rodzaje = []shared.LibraryDuplicateKind{shared.LibraryDuplicateKindExact}
	}
	prog := progPodobienstwaDomyslny
	if z.Threshold != nil && *z.Threshold > 0 {
		prog = *z.Threshold
	}

	grupy := []shared.LibraryDuplicateGroup{}
	pominiete := 0
	for _, rodzaj := range rodzaje {
		switch rodzaj {
		case shared.LibraryDuplicateKindExact:
			dokladne, bezSumy := grupyDokladneBiblioteki(zasoby)
			grupy = append(grupy, dokladne...)
			pominiete = bezSumy
		case shared.LibraryDuplicateKindNearText:
			grupy = append(grupy, a.grupyPrzyblizone(zasoby, prog, false)...)
		case shared.LibraryDuplicateKindNearImage:
			grupy = append(grupy, a.grupyPrzyblizone(zasoby, prog, true)...)
		default:
			return shared.LibraryDuplicateScanResponse{}, bladWskazaniaBiblioteki(
				"rodzaj rozpoznania " + string(rodzaj) + " nie ma odwzorowania w rdzeniu")
		}
	}
	sort.SliceStable(grupy, func(pierwsza, druga int) bool {
		return len(grupy[pierwsza].FileIds) > len(grupy[druga].FileIds)
	})
	lacznie := len(grupy)
	if z.Limit != nil && *z.Limit > 0 && len(grupy) > *z.Limit {
		grupy = grupy[:*z.Limit]
	}
	a.odnotuj(ctx, shared.LibraryAuditActionAccess, nil,
		"skanowanie duplikatów: grup "+strconv.Itoa(lacznie))
	return shared.LibraryDuplicateScanResponse{
		Groups: grupy, Total: lacznie, SkippedWithoutChecksum: pominiete,
	}, nil
}

// grupyDokladneBiblioteki składa grupy zasobów o wspólnej sumie kontrolnej
// i liczy zasoby, których sumy nie ma — brak sumy znaczy niewiadomą, nie brak
// duplikatu.
func grupyDokladneBiblioteki(zasoby []dane.PlikBiblioteki) ([]shared.LibraryDuplicateGroup, int) {
	wedlugSumy := map[string][]string{}
	bezSumy := 0
	for _, zasob := range zasoby {
		if zasob.SumaKontrolna == nil || *zasob.SumaKontrolna == "" {
			bezSumy++
			continue
		}
		wedlugSumy[*zasob.SumaKontrolna] = append(wedlugSumy[*zasob.SumaKontrolna], zasob.Kod)
	}
	grupy := []shared.LibraryDuplicateGroup{}
	sumy := make([]string, 0, len(wedlugSumy))
	for suma := range wedlugSumy {
		sumy = append(sumy, suma)
	}
	sort.Strings(sumy)
	for _, suma := range sumy {
		kody := wedlugSumy[suma]
		if len(kody) < 2 {
			continue
		}
		wspolna := suma
		grupy = append(grupy, shared.LibraryDuplicateGroup{
			Kind: shared.LibraryDuplicateKindExact, Checksum: &wspolna, FileIds: kody,
		})
	}
	return grupy, bezSumy
}

// grupyPrzyblizone składa grupy zasobów podobnych treścią albo obrazem.
//
// Zasoby o identycznej sumie kontrolnej są pomijane: należą do rozpoznania
// dokładnego, a wystawione tu po raz drugi kazałyby Operatorowi rozstrzygać tę
// samą parę dwa razy.
func (a *adapterBiblioteki) grupyPrzyblizone(zasoby []dane.PlikBiblioteki, prog int,
	obrazy bool) []shared.LibraryDuplicateGroup {

	type odcisk struct {
		kod   string
		slowa map[string]struct{}
		obraz uint64
	}
	odciski := make([]odcisk, 0, len(zasoby))
	for _, zasob := range zasoby {
		if zasob.TrescOdwolanie == nil || *zasob.TrescOdwolanie == "" {
			continue
		}
		if obrazy != rodzajPodgladuJestObrazem(zasob) {
			continue
		}
		bajty, err := os.ReadFile(*zasob.TrescOdwolanie)
		if err != nil {
			continue
		}
		if obrazy {
			skrot, err := odciskObrazuBiblioteki(bajty)
			if err != nil {
				continue
			}
			odciski = append(odciski, odcisk{kod: zasob.Kod, obraz: skrot})
			continue
		}
		odciski = append(odciski, odcisk{kod: zasob.Kod, slowa: odciskTekstuBiblioteki(bajty)})
	}

	grupy := []shared.LibraryDuplicateGroup{}
	wgrupowane := map[string]bool{}
	for pierwszy := range odciski {
		if wgrupowane[odciski[pierwszy].kod] {
			continue
		}
		grupa := []string{odciski[pierwszy].kod}
		najlepsza := 0
		for drugi := pierwszy + 1; drugi < len(odciski); drugi++ {
			if wgrupowane[odciski[drugi].kod] {
				continue
			}
			var trafnosc int
			if obrazy {
				trafnosc = podobienstwoOdciskowObrazuBiblioteki(odciski[pierwszy].obraz, odciski[drugi].obraz)
			} else {
				trafnosc = podobienstwoZbiorowSlowBiblioteki(odciski[pierwszy].slowa, odciski[drugi].slowa)
			}
			if trafnosc < prog || trafnosc >= 100 {
				// Trafność stuprocentowa znaczy treść identyczną — to jest
				// rozpoznanie dokładne, nie przybliżone.
				continue
			}
			grupa = append(grupa, odciski[drugi].kod)
			wgrupowane[odciski[drugi].kod] = true
			if trafnosc > najlepsza {
				najlepsza = trafnosc
			}
		}
		if len(grupa) < 2 {
			continue
		}
		wgrupowane[odciski[pierwszy].kod] = true
		trafnosc := najlepsza
		rodzaj := shared.LibraryDuplicateKind(shared.LibraryDuplicateKindNearText)
		if obrazy {
			rodzaj = shared.LibraryDuplicateKindNearImage
		}
		grupy = append(grupy, shared.LibraryDuplicateGroup{
			Kind: rodzaj, Score: &trafnosc, FileIds: grupa,
		})
	}
	return grupy
}

// PolaczDuplikaty obsługuje `library.duplicate.merge`.
//
// Zasoby wchłaniane trafiają do archiwum, nie znikają: scalenie bywa pomyłką,
// a kosz repozytorium jest odwracalny — usunięcie trwałe ma własną komendę
// i własne potwierdzenie.
func (a *adapterBiblioteki) PolaczDuplikaty(ctx context.Context,
	z shared.LibraryDuplicateMergeRequest) (shared.LibraryDuplicateMergeResponse, error) {

	docelowy, err := a.plik(ctx, z.TargetFileId)
	if err != nil {
		return shared.LibraryDuplicateMergeResponse{}, err
	}
	if len(z.SourceFileIds) == 0 {
		return shared.LibraryDuplicateMergeResponse{}, bladWskazaniaBiblioteki(
			"scalanie bez wskazania zasobów wchłanianych")
	}

	etykietyDocelowe, err := a.repozytorium.Etykiety(ctx, docelowy.ID)
	if err != nil {
		return shared.LibraryDuplicateMergeResponse{}, bladBiblioteki(err)
	}
	kolekcjeDocelowe, err := a.repozytorium.KolekcjePliku(ctx, docelowy.ID)
	if err != nil {
		return shared.LibraryDuplicateMergeResponse{}, bladBiblioteki(err)
	}
	przeniesione := map[string]bool{}
	for _, etykieta := range etykietyDocelowe {
		przeniesione[etykieta] = false
	}

	przeniesioneEtykiety := []string{}
	wersjePrzeniesione := 0
	wchloniete := 0
	przenosWersje := z.KeepVersions != nil && *z.KeepVersions

	for _, kodZrodla := range z.SourceFileIds {
		if strings.TrimSpace(kodZrodla) == "" || kodZrodla == docelowy.Kod {
			continue
		}
		zrodlo, err := a.plik(ctx, kodZrodla)
		if err != nil {
			return shared.LibraryDuplicateMergeResponse{}, err
		}
		etykiety, err := a.repozytorium.Etykiety(ctx, zrodlo.ID)
		if err != nil {
			return shared.LibraryDuplicateMergeResponse{}, bladBiblioteki(err)
		}
		for _, etykieta := range etykiety {
			if _, jest := przeniesione[etykieta]; jest {
				continue
			}
			przeniesione[etykieta] = true
			etykietyDocelowe = append(etykietyDocelowe, etykieta)
			przeniesioneEtykiety = append(przeniesioneEtykiety, etykieta)
		}
		kolekcje, err := a.repozytorium.KolekcjePliku(ctx, zrodlo.ID)
		if err != nil {
			return shared.LibraryDuplicateMergeResponse{}, bladBiblioteki(err)
		}
		kolekcjeDocelowe = polaczWykazyBiblioteki(kolekcjeDocelowe, kolekcje)

		if przenosWersje {
			wersje, err := a.repozytorium.Wersje(ctx, zrodlo.ID)
			if err != nil {
				return shared.LibraryDuplicateMergeResponse{}, bladBiblioteki(err)
			}
			for _, wersja := range wersje {
				// Wersja wchłonięta wchodzi do historii docelowego jako nowy
				// wpis pod własnym identyfikatorem: kod wersji jest unikalny
				// w całym repozytorium, więc przeniesienie go wprost byłoby
				// zderzeniem kluczy.
				_, err := a.repozytorium.ZapiszWersje(ctx, docelowy.ID, dane.WersjaPlikuBiblioteki{
					Kod: nowyIdentyfikator(przedrostekWersjiBiblioteki), PlikID: docelowy.ID,
					Etykieta: wskazanieBiblioteki("z zasobu " + zrodlo.Nazwa),
					Autor:    wersja.Autor, RozmiarBajtow: wersja.RozmiarBajtow,
					SumaKontrolna: wersja.SumaKontrolna, TrescOdwolanie: wersja.TrescOdwolanie,
				})
				if err != nil {
					return shared.LibraryDuplicateMergeResponse{}, bladBiblioteki(err)
				}
				wersjePrzeniesione++
			}
		}

		if _, err := a.repozytorium.UstawStanPlikow(ctx,
			[]string{zrodlo.Kod}, dane.StanZasobuZarchiwizowany); err != nil {
			return shared.LibraryDuplicateMergeResponse{}, bladBiblioteki(err)
		}
		a.odnotuj(ctx, shared.LibraryAuditActionArchive, &zrodlo.Kod,
			"wchłonięty przy scaleniu duplikatów do "+docelowy.Kod)
		wchloniete++
	}

	if len(przeniesioneEtykiety) > 0 {
		if _, err := a.repozytorium.UstawEtykiety(ctx, docelowy.Kod, etykietyDocelowe); err != nil {
			return shared.LibraryDuplicateMergeResponse{}, bladBiblioteki(err)
		}
	}
	if len(kolekcjeDocelowe) > 0 {
		if _, err := a.repozytorium.UstawKolekcjePliku(ctx, docelowy.Kod, kolekcjeDocelowe); err != nil {
			return shared.LibraryDuplicateMergeResponse{}, bladBiblioteki(err)
		}
	}
	wiersz, err := a.repozytorium.Plik(ctx, docelowy.Kod)
	if err != nil {
		return shared.LibraryDuplicateMergeResponse{}, bladNieznanegoPlikuBiblioteki(docelowy.Kod, err)
	}
	kontrakt, err := a.zloz(ctx, wiersz)
	if err != nil {
		return shared.LibraryDuplicateMergeResponse{}, err
	}
	a.odnotuj(ctx, shared.LibraryAuditActionChange, &docelowy.Kod,
		"scalenie duplikatów: wchłonięto "+strconv.Itoa(wchloniete))
	a.zglosNasluchom(shared.LibraryWebhookEventFileChanged, docelowy.Kod)
	return shared.LibraryDuplicateMergeResponse{
		File: kontrakt, MergedCount: wchloniete,
		CarriedTags: przeniesioneEtykiety, CarriedVersions: wersjePrzeniesione,
	}, nil
}

// SprawdzIntegralnosc obsługuje `library.fixity.check`.
//
// Suma liczy się z bajtów leżących pod odwołaniem i porównuje z zapisaną przy
// zasobie. Brak treści pod odwołaniem to trzeci wynik, nie odmiana niezgodności:
// sum nie ma czego porównać, a zasób i tak jest uszkodzony.
func (a *adapterBiblioteki) SprawdzIntegralnosc(ctx context.Context,
	z shared.LibraryFixityCheckRequest) (shared.LibraryFixityCheckResponse, error) {

	zasoby, err := a.zasobyZbioru(ctx, z.FileIds, z.CollectionId)
	if err != nil {
		return shared.LibraryFixityCheckResponse{}, err
	}
	tylkoNiezgodne := z.MismatchedOnly != nil && *z.MismatchedOnly
	teraz := time.Now().UnixMilli()

	wyniki := []shared.LibraryFixityResult{}
	niezgodne, brakujace := 0, 0
	for _, zasob := range zasoby {
		wynik := shared.LibraryFixityResult{
			FileId: zasob.Kod, ExpectedChecksum: zasob.SumaKontrolna, CheckedAt: teraz,
		}
		if zasob.TrescOdwolanie == nil || *zasob.TrescOdwolanie == "" {
			wynik.Missing = true
			brakujace++
		} else if bajty, err := os.ReadFile(*zasob.TrescOdwolanie); err != nil {
			wynik.Missing = true
			brakujace++
		} else {
			skrot := sha256.Sum256(bajty)
			wyliczona := hex.EncodeToString(skrot[:])
			wynik.ActualChecksum = &wyliczona
			wynik.Matched = zasob.SumaKontrolna != nil && *zasob.SumaKontrolna == wyliczona
			if !wynik.Matched {
				niezgodne++
			}
		}
		if tylkoNiezgodne && wynik.Matched {
			continue
		}
		wyniki = append(wyniki, wynik)
	}
	a.odnotuj(ctx, shared.LibraryAuditActionAccess, nil,
		"weryfikacja integralności: sprawdzono "+strconv.Itoa(len(zasoby))+
			", niezgodnych "+strconv.Itoa(niezgodne))
	return shared.LibraryFixityCheckResponse{
		Results: wyniki, CheckedCount: len(zasoby),
		MismatchedCount: niezgodne, MissingCount: brakujace,
	}, nil
}

// NormalizujNazwy obsługuje `library.name.normalize`.
//
// Przebieg próbny pokazuje wynik bez zapisu — nazwa jest tym, po czym Operator
// odnajduje zasób, więc masowa zmiana bez podglądu byłaby zmianą w ciemno.
func (a *adapterBiblioteki) NormalizujNazwy(ctx context.Context,
	z shared.LibraryNameNormalizeRequest) (shared.LibraryNameNormalizeResponse, error) {

	zasoby, err := a.zasobyZbioru(ctx, z.FileIds, z.CollectionId)
	if err != nil {
		return shared.LibraryNameNormalizeResponse{}, err
	}
	probny := z.DryRun != nil && *z.DryRun
	transliteracja := z.Transliterate == nil || *z.Transliterate
	schemat := strings.TrimSpace(wartoscTekstu(z.Pattern))

	wyniki := make([]shared.LibraryNormalizeResult, 0, len(zasoby))
	zmienione := 0
	for indeks, zasob := range zasoby {
		nowa := nazwaZnormalizowanaBiblioteki(zasob, schemat, transliteracja, indeks+1)
		wynik := shared.LibraryNormalizeResult{
			FileId: zasob.Kod, PreviousName: zasob.Nazwa, NewName: nowa,
		}
		if nowa != zasob.Nazwa {
			zmienione++
			if !probny {
				if _, err := a.repozytorium.PrzemianujPlik(ctx, zasob.Kod, nowa); err != nil {
					return shared.LibraryNameNormalizeResponse{}, bladBiblioteki(err)
				}
				wynik.Applied = true
				a.zglosNasluchom(shared.LibraryWebhookEventFileChanged, zasob.Kod)
			}
		}
		wyniki = append(wyniki, wynik)
	}
	if !probny && zmienione > 0 {
		a.odnotuj(ctx, shared.LibraryAuditActionChange, nil,
			"normalizacja nazw: zmieniono "+strconv.Itoa(zmienione))
	}
	return shared.LibraryNameNormalizeResponse{Results: wyniki, ChangedCount: zmienione}, nil
}

// PulpitStanu obsługuje `library.stats.get`.
func (a *adapterBiblioteki) PulpitStanu(ctx context.Context,
	z shared.LibraryStatsGetRequest) (shared.LibraryStatsGetResponse, error) {

	topN := 0
	if z.TopN != nil {
		topN = *z.TopN
	}
	stat, err := a.repozytorium.Statystyki(ctx, z.CollectionId, z.ProjectId, topN)
	if err != nil {
		return shared.LibraryStatsGetResponse{}, bladBiblioteki(err)
	}
	return shared.LibraryStatsGetResponse{Stats: shared.LibraryStats{
		FileCount: stat.LiczbaZasobow, ArchivedCount: stat.LiczbaArchiwalnych,
		TotalBytes: stat.LacznyRozmiar, OrphanCount: stat.LiczbaOsieroconych,
		DuplicateCount: stat.LiczbaDuplikatow, MissingChecksumCount: stat.BezSumyKontrolnej,
		ByMimeType:      rozkladKontraktuBiblioteki(stat.WedlugRodzajuTresci),
		BySourceModule:  rozkladKontraktuBiblioteki(stat.WedlugModuluZrodla),
		TagUsage:        rozkladKontraktuBiblioteki(stat.UzycieEtykiet),
		CollectionUsage: rozkladKontraktuBiblioteki(stat.UzycieKolekcji),
	}}, nil
}

// zasobyZbioru składa zbiór zasobów wskazany żądaniem: wykaz kodów, kolekcja
// albo — gdy nie ma ani jednego wskazania — całe repozytorium czynne.
func (a *adapterBiblioteki) zasobyZbioru(ctx context.Context, kody []string,
	kolekcjaKod *string) ([]dane.PlikBiblioteki, error) {

	if len(kody) > 0 {
		zasoby := make([]dane.PlikBiblioteki, 0, len(kody))
		for _, kod := range kody {
			zasob, err := a.plik(ctx, kod)
			if err != nil {
				return nil, err
			}
			zasoby = append(zasoby, zasob)
		}
		return zasoby, nil
	}
	stan := dane.StanZasobuCzynny
	wiersze, _, err := a.repozytorium.Pliki(ctx, dane.FiltrPlikow{
		KolekcjaKod: kolekcjaKod, Stan: &stan, Limit: granicaZbioruHigieny,
	})
	if err != nil {
		return nil, bladBiblioteki(err)
	}
	return wiersze, nil
}

// nazwaZnormalizowanaBiblioteki składa nazwę zasobu według reguły nazewnictwa.
//
// Schemat rozumie trzy wzorce: `{nazwa}`, `{rozszerzenie}` i `{numer}`. Bez
// schematu nazwa zostaje ta sama, a zmianie podlega wyłącznie jej zapis:
// znaki diakrytyczne schodzą do postaci podstawowej, znaki spoza zakresu
// bezpiecznego ustępują myślnikowi, a wielkość liter idzie w dół.
func nazwaZnormalizowanaBiblioteki(zasob dane.PlikBiblioteki, schemat string, transliteracja bool,
	numer int) string {

	trzon, rozszerzenie := trzonIRozszerzenieBiblioteki(zasob.Nazwa)
	if schemat != "" {
		zamiana := strings.NewReplacer(
			"{nazwa}", trzon,
			"{rozszerzenie}", rozszerzenie,
			"{numer}", strconv.Itoa(numer),
		)
		zlozona := zamiana.Replace(schemat)
		trzon, rozszerzenie = trzonIRozszerzenieBiblioteki(zlozona)
	}
	trzon = zapisBezpiecznyNazwyBiblioteki(trzon, transliteracja)
	if trzon == "" {
		trzon = "zasob"
	}
	if rozszerzenie == "" {
		return trzon
	}
	return trzon + "." + zapisBezpiecznyNazwyBiblioteki(rozszerzenie, transliteracja)
}

// trzonIRozszerzenieBiblioteki rozdziela nazwę na człon główny i rozszerzenie.
func trzonIRozszerzenieBiblioteki(nazwa string) (string, string) {
	kropka := strings.LastIndex(nazwa, ".")
	if kropka <= 0 || kropka == len(nazwa)-1 {
		return nazwa, ""
	}
	return nazwa[:kropka], nazwa[kropka+1:]
}

// zapisBezpiecznyNazwyBiblioteki sprowadza nazwę do postaci ujednoliconej.
func zapisBezpiecznyNazwyBiblioteki(tekst string, transliteracja bool) string {
	if transliteracja {
		// Rozkład kanoniczny, zdjęcie znaków łączących, złożenie z powrotem —
		// „ą" staje się „a", a „ł" zostaje, bo nie jest literą ze znakiem
		// łączącym. Zamiana liter osobnych idzie niżej, wprost.
		bezZnakow := transform.Chain(norm.NFD,
			runes.Remove(runes.In(unicode.Mn)), norm.NFC)
		if wynik, _, err := transform.String(bezZnakow, tekst); err == nil {
			tekst = wynik
		}
		tekst = strings.NewReplacer("ł", "l", "Ł", "L").Replace(tekst)
	}
	tekst = strings.ToLower(strings.TrimSpace(tekst))

	var zapis strings.Builder
	poprzedniMyslnik := false
	for _, znak := range tekst {
		switch {
		case znak >= 'a' && znak <= 'z', znak >= '0' && znak <= '9':
			zapis.WriteRune(znak)
			poprzedniMyslnik = false
		default:
			if !poprzedniMyslnik && zapis.Len() > 0 {
				zapis.WriteRune('-')
				poprzedniMyslnik = true
			}
		}
	}
	return strings.Trim(zapis.String(), "-")
}

// polaczWykazyBiblioteki scala dwa wykazy kodów bez powtórzeń.
func polaczWykazyBiblioteki(pierwszy, drugi []string) []string {
	obecne := map[string]bool{}
	wynik := make([]string, 0, len(pierwszy)+len(drugi))
	for _, wykaz := range [][]string{pierwszy, drugi} {
		for _, kod := range wykaz {
			if kod == "" || obecne[kod] {
				continue
			}
			obecne[kod] = true
			wynik = append(wynik, kod)
		}
	}
	return wynik
}

// rozkladKontraktuBiblioteki przenosi rozkład warstwy danych na kontrakt.
func rozkladKontraktuBiblioteki(pozycje []dane.LiczbaWedlugKlucza) []shared.LibraryCount {
	wynik := make([]shared.LibraryCount, 0, len(pozycje))
	for _, pozycja := range pozycje {
		wynik = append(wynik, shared.LibraryCount{Key: pozycja.Klucz, Count: pozycja.Liczba})
	}
	return wynik
}

// terazZnacznikBiblioteki oddaje chwilę bieżącą w zapisie znacznika schematu.
func terazZnacznikBiblioteki() string {
	return time.Now().UTC().Format(formatZnacznikaBazy)
}
