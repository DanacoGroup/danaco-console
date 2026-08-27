// Plik obsługuje rodzinę `research.source.*` poza samym dodaniem: przegląd
// katalogu, zmianę, usunięcie, scalenie, deduplikację, etykiety, załączniki,
// import bibliografii, przechwycenie strony, transkrypcję nagrania
// i rozstrzyganie identyfikatorów.
package core

import (
	"context"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"errors"
	"os"
	"sort"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// WypiszZrodla obsługuje komendę `research.source.list`: zwraca stronę
// źródeł okna badania zawężoną zapytaniem, rodzajem i pozostałymi filtrami.
func (a *adapterBadan) WypiszZrodla(ctx context.Context,
	z shared.ResearchSourceListRequest) (shared.ResearchSourceListResponse, error) {

	if z.WindowId == "" {
		return shared.ResearchSourceListResponse{}, bladWskazaniaBadan("source.list bez okna badania")
	}
	zrodla, err := a.repozytorium.Zrodla(ctx, z.WindowId)
	if err != nil {
		return shared.ResearchSourceListResponse{}, bladBadan(err)
	}

	dopasowane := make([]dane.ZrodloBadania, 0, len(zrodla))
	for _, zrodlo := range zrodla {
		pasuje, err := a.zrodloPasujeBadania(ctx, zrodlo, z)
		if err != nil {
			return shared.ResearchSourceListResponse{}, err
		}
		if pasuje {
			dopasowane = append(dopasowane, zrodlo)
		}
	}

	// Całość policzona przed wycinkiem strony, nie po nim.
	razem := len(dopasowane)
	wycinek := wycinekListyBadania(razem, z.Offset, z.Limit)
	przelozone := make([]shared.ResearchSource, 0, wycinek.doPozycji-wycinek.odPozycji)
	for _, zrodlo := range dopasowane[wycinek.odPozycji:wycinek.doPozycji] {
		przelozone = append(przelozone, zlozZrodloBadania(zrodlo))
	}
	return shared.ResearchSourceListResponse{Sources: przelozone, Total: razem}, nil
}

// zakresListyBadania opisuje wycinek listy źródeł po zastosowaniu
// przesunięcia i limitu żądania w postaci pary granic indeksu.
type zakresListyBadania struct {
	odPozycji int
	doPozycji int
}

// wycinekListyBadania wylicza granice strony wyniku. Wskazania poza zakresem
// dają stronę pustą, a nie usterkę: „nie ma dalszych pozycji" jest odpowiedzią.
func wycinekListyBadania(razem int, przesuniecie, limit *int) zakresListyBadania {
	od := 0
	if przesuniecie != nil && *przesuniecie > 0 {
		od = *przesuniecie
	}
	if od > razem {
		od = razem
	}
	do := razem
	if limit != nil && *limit > 0 && od+*limit < razem {
		do = od + *limit
	}
	return zakresListyBadania{odPozycji: od, doPozycji: do}
}

// zrodloPasujeBadania rozstrzyga, czy źródło spełnia zawężenia żądania: frazę,
// rodzaj, wiarygodność, stan lektury i wszystkie żądane etykiety.
func (a *adapterBadan) zrodloPasujeBadania(ctx context.Context, zrodlo dane.ZrodloBadania,
	z shared.ResearchSourceListRequest) (bool, error) {

	if z.Query != nil && strings.TrimSpace(*z.Query) != "" {
		szukane := strings.ToLower(strings.TrimSpace(*z.Query))
		wStronie := strings.ToLower(zrodlo.Tytul)
		if zrodlo.Adres != nil {
			wStronie += " " + strings.ToLower(*zrodlo.Adres)
		}
		if !strings.Contains(wStronie, szukane) {
			return false, nil
		}
	}
	if z.Kind != nil && *z.Kind != "" && zrodlo.Rodzaj != *z.Kind {
		return false, nil
	}
	if z.Credibility != nil && *z.Credibility != "" && zrodlo.Wiarygodnosc != *z.Credibility {
		return false, nil
	}
	if z.ReadingState != nil && *z.ReadingState != "" {
		lektura, err := a.repozytorium.LekturaZrodlaBadania(ctx, zrodlo.Kod)
		if err != nil {
			return false, bladBadan(err)
		}
		if lektura.StanLektury != string(*z.ReadingState) {
			return false, nil
		}
	}
	if len(z.Tags) > 0 {
		katalog, err := a.repozytorium.KatalogZrodlaBadania(ctx, zrodlo.Kod)
		if err != nil {
			return false, bladBadan(err)
		}
		for _, etykieta := range z.Tags {
			if !zawieraNapisBadania(katalog.Etykiety, etykieta) {
				return false, nil
			}
		}
	}
	return true, nil
}

// ZmienZrodlo obsługuje `research.source.update`. Pola pominięte zostają bez
// zmiany — żądanie zmiany tytułu nie ma prawa wyczyścić adresu.
func (a *adapterBadan) ZmienZrodlo(ctx context.Context,
	z shared.ResearchSourceUpdateRequest) (shared.ResearchSourceUpdateResponse, error) {

	if z.SourceId == "" {
		return shared.ResearchSourceUpdateResponse{}, bladWskazaniaBadan("source.update bez źródła")
	}
	zastane, err := a.repozytorium.Zrodlo(ctx, z.SourceId)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.ResearchSourceUpdateResponse{}, bladNieznanegoZrodlaBadania(z.SourceId)
	}
	if err != nil {
		return shared.ResearchSourceUpdateResponse{}, bladBadan(err)
	}

	if z.Title != nil {
		zastane.Tytul = *z.Title
	}
	if z.Kind != nil && *z.Kind != "" {
		zastane.Rodzaj = *z.Kind
	}
	if z.Url != nil {
		zastane.Adres = z.Url
	}
	if z.Origin != nil {
		zastane.Pochodzenie = z.Origin
	}
	if z.Credibility != nil && *z.Credibility != "" {
		zastane.Wiarygodnosc = *z.Credibility
	}
	if z.LibraryFileId != nil {
		zastane.PlikBibliotekiID = z.LibraryFileId
	}
	zapisane, err := a.repozytorium.ZapiszZrodlo(ctx, zastane)
	if err != nil {
		return shared.ResearchSourceUpdateResponse{}, bladBadan(err)
	}

	lektura := dane.LekturaZrodlaBadania{}
	if z.ReadingState != nil {
		lektura.StanLektury = string(*z.ReadingState)
	}
	if z.StageIndex != nil {
		etap := int64(*z.StageIndex)
		lektura.EtapIndeks = &etap
	}
	if err := a.repozytorium.UstawLektureZrodla(ctx, z.SourceId, lektura); err != nil {
		return shared.ResearchSourceUpdateResponse{}, bladBadan(err)
	}
	if z.QuestionIds != nil {
		katalog, err := a.repozytorium.KatalogZrodlaBadania(ctx, z.SourceId)
		if err != nil {
			return shared.ResearchSourceUpdateResponse{}, bladBadan(err)
		}
		katalog.Pytania = z.QuestionIds
		if err := a.repozytorium.UstawKatalogZrodla(ctx, z.SourceId, katalog); err != nil {
			return shared.ResearchSourceUpdateResponse{}, bladBadan(err)
		}
	}
	return shared.ResearchSourceUpdateResponse{Source: zlozZrodloBadania(zapisane)}, nil
}

// UsunZrodlo obsługuje `research.source.remove`. Odpowiedź niesie liczbę
// ustaleń, które straciły odwołanie — skutek zdjęcia, nie samo potwierdzenie.
func (a *adapterBadan) UsunZrodlo(ctx context.Context,
	z shared.ResearchSourceRemoveRequest) (shared.ResearchSourceRemoveResponse, error) {

	if z.SourceId == "" {
		return shared.ResearchSourceRemoveResponse{}, bladWskazaniaBadan("source.remove bez źródła")
	}
	odwiazane, err := a.repozytorium.UsunZrodlo(ctx, z.SourceId)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.ResearchSourceRemoveResponse{}, bladNieznanegoZrodlaBadania(z.SourceId)
	}
	if err != nil {
		return shared.ResearchSourceRemoveResponse{}, bladBadan(err)
	}
	return shared.ResearchSourceRemoveResponse{SourceId: z.SourceId, DetachedFindings: odwiazane}, nil
}

// ScalZrodla obsługuje komendę `research.source.merge`: przenosi ustalenia ze
// źródeł scalanych do źródła docelowego i oddaje odświeżone źródło docelowe.
func (a *adapterBadan) ScalZrodla(ctx context.Context,
	z shared.ResearchSourceMergeRequest) (shared.ResearchSourceMergeResponse, error) {

	if z.TargetSourceId == "" {
		return shared.ResearchSourceMergeResponse{}, bladWskazaniaBadan("source.merge bez źródła docelowego")
	}
	if len(z.MergedSourceIds) == 0 {
		return shared.ResearchSourceMergeResponse{},
			bladWskazaniaBadan("source.merge bez źródeł do scalenia — nie ma czego przenieść")
	}
	przeniesione, err := a.repozytorium.PrzeniesUstaleniaZrodel(ctx, z.MergedSourceIds, z.TargetSourceId)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.ResearchSourceMergeResponse{}, bladNieznanegoZrodlaBadania(z.TargetSourceId)
	}
	if err != nil {
		return shared.ResearchSourceMergeResponse{}, bladBadan(err)
	}
	zrodlo, err := a.repozytorium.Zrodlo(ctx, z.TargetSourceId)
	if err != nil {
		return shared.ResearchSourceMergeResponse{}, bladBadan(err)
	}
	return shared.ResearchSourceMergeResponse{
		Source: zlozZrodloBadania(zrodlo), MovedFindings: przeniesione,
	}, nil
}

// SzukajDuplikatowZrodel obsługuje komendę `research.source.duplicates`:
// szuka par źródeł podobnych ponad próg i nazywa podstawę dopasowania.
func (a *adapterBadan) SzukajDuplikatowZrodel(ctx context.Context,
	z shared.ResearchSourceDuplicatesRequest) (shared.ResearchSourceDuplicatesResponse, error) {

	if z.WindowId == "" {
		return shared.ResearchSourceDuplicatesResponse{},
			bladWskazaniaBadan("source.duplicates bez okna badania")
	}
	prog := 70
	if z.Threshold != nil && *z.Threshold > 0 {
		prog = *z.Threshold
	}
	zrodla, err := a.repozytorium.Zrodla(ctx, z.WindowId)
	if err != nil {
		return shared.ResearchSourceDuplicatesResponse{}, bladBadan(err)
	}

	kandydaci := []shared.ResearchDuplicateCandidate{}
	for i := 0; i < len(zrodla); i++ {
		for j := i + 1; j < len(zrodla); j++ {
			podobienstwo, podstawa := podobienstwoZrodelBadania(ctx, a, zrodla[i], zrodla[j])
			if podobienstwo < prog {
				continue
			}
			kandydaci = append(kandydaci, shared.ResearchDuplicateCandidate{
				SourceIds:  []string{zrodla[i].Kod, zrodla[j].Kod},
				Similarity: podobienstwo, Basis: podstawa,
			})
		}
	}
	return shared.ResearchSourceDuplicatesResponse{Candidates: kandydaci}, nil
}

// podobienstwoZrodelBadania ocenia, na ile dwie pozycje są tą samą pracą,
// i oddaje wynik podobieństwa wraz z nazwą podstawy dopasowania.
func podobienstwoZrodelBadania(ctx context.Context, a *adapterBadan,
	pierwsze, drugie dane.ZrodloBadania) (int, string) {

	pierwszaLektura, _ := a.repozytorium.LekturaZrodlaBadania(ctx, pierwsze.Kod)
	drugaLektura, _ := a.repozytorium.LekturaZrodlaBadania(ctx, drugie.Kod)
	if pierwszaLektura.Identyfikat != nil && drugaLektura.Identyfikat != nil &&
		strings.EqualFold(*pierwszaLektura.Identyfikat, *drugaLektura.Identyfikat) &&
		strings.TrimSpace(*pierwszaLektura.Identyfikat) != "" {
		return 100, "identyfikator"
	}
	if pierwsze.Adres != nil && drugie.Adres != nil &&
		strings.EqualFold(strings.TrimSpace(*pierwsze.Adres), strings.TrimSpace(*drugie.Adres)) &&
		strings.TrimSpace(*pierwsze.Adres) != "" {
		return 95, "adres"
	}
	return podobienstwoTekstuBadania(pierwsze.Tytul, drugie.Tytul), "tytuł"
}

// podobienstwoTekstuBadania liczy podobieństwo dwóch tytułów miarą Jaccarda na
// zbiorach słów. Miara na słowach, nie na znakach: „Raport rynkowy 2024"
// i „2024 raport rynkowy" to ta sama praca, a odległość znakowa uznałaby je za
// odległe.
func podobienstwoTekstuBadania(pierwszy, drugi string) int {
	zbiorA := zbiorSlowBadania(pierwszy)
	zbiorB := zbiorSlowBadania(drugi)
	if len(zbiorA) == 0 || len(zbiorB) == 0 {
		return 0
	}
	wspolne := 0
	for slowo := range zbiorA {
		if zbiorB[slowo] {
			wspolne++
		}
	}
	suma := len(zbiorA) + len(zbiorB) - wspolne
	if suma == 0 {
		return 0
	}
	return wspolne * 100 / suma
}

// zbiorSlowBadania rozkłada tekst na zbiór słów sprowadzonych do małych
// liter, pomijając słowa krótsze niż trzy znaki.
func zbiorSlowBadania(tekst string) map[string]bool {
	zbior := map[string]bool{}
	for _, slowo := range strings.FieldsFunc(strings.ToLower(tekst), func(znak rune) bool {
		return !('a' <= znak && znak <= 'z' || '0' <= znak && znak <= '9' || znak > 127)
	}) {
		if len(slowo) > 2 {
			zbior[slowo] = true
		}
	}
	return zbior
}

// OznaczZrodlo obsługuje komendę `research.source.tag`: zapisuje etykiety
// i przynależność do kolekcji wskazanego źródła.
func (a *adapterBadan) OznaczZrodlo(ctx context.Context,
	z shared.ResearchSourceTagRequest) (shared.ResearchSourceTagResponse, error) {

	if z.SourceId == "" {
		return shared.ResearchSourceTagResponse{}, bladWskazaniaBadan("source.tag bez źródła")
	}
	katalog, err := a.repozytorium.KatalogZrodlaBadania(ctx, z.SourceId)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.ResearchSourceTagResponse{}, bladNieznanegoZrodlaBadania(z.SourceId)
	}
	if err != nil {
		return shared.ResearchSourceTagResponse{}, bladBadan(err)
	}
	if z.Tags != nil {
		katalog.Etykiety = z.Tags
	}
	if z.CollectionIds != nil {
		katalog.Kolekcje = z.CollectionIds
	}
	if err := a.repozytorium.UstawKatalogZrodla(ctx, z.SourceId, katalog); err != nil {
		return shared.ResearchSourceTagResponse{}, bladBadan(err)
	}
	zrodlo, err := a.repozytorium.Zrodlo(ctx, z.SourceId)
	if err != nil {
		return shared.ResearchSourceTagResponse{}, bladBadan(err)
	}
	return shared.ResearchSourceTagResponse{Source: zlozZrodloBadania(zrodlo)}, nil
}

// DodajZalacznikZrodla obsługuje `research.source.attachment.add`. Treść
// podana wprost albo ścieżką ląduje w magazynie modułu — załącznik bez bajtów
// byłby wierszem o pliku, którego nie ma.
func (a *adapterBadan) DodajZalacznikZrodla(ctx context.Context,
	z shared.ResearchSourceAttachmentAddRequest) (shared.ResearchSourceAttachmentAddResponse, error) {

	if z.SourceId == "" {
		return shared.ResearchSourceAttachmentAddResponse{},
			bladWskazaniaBadan("source.attachment.add bez źródła")
	}
	if z.Kind == "" {
		return shared.ResearchSourceAttachmentAddResponse{},
			bladWskazaniaBadan("source.attachment.add bez rodzaju załącznika")
	}

	zalacznik := dane.ZalacznikZrodlaBadania{
		Kod: nowyIdentyfikator(przedrostekZalacznikaBadania), ZrodloKod: z.SourceId,
		Rodzaj: string(z.Kind), PlikBibliotekiID: z.LibraryFileId,
	}
	bajty, err := trescZadaniaBadania(z.ContentBase64, z.SourcePath)
	if err != nil {
		return shared.ResearchSourceAttachmentAddResponse{}, err
	}
	if bajty != nil {
		sciezka, rozmiar, err := a.odlozMaterialBadania(bajty)
		if err != nil {
			return shared.ResearchSourceAttachmentAddResponse{}, err
		}
		zalacznik.Sciezka = &sciezka
		zalacznik.RozmiarBajtow = &rozmiar
	} else if z.LibraryFileId == nil {
		return shared.ResearchSourceAttachmentAddResponse{}, bladWskazaniaBadan(
			"source.attachment.add bez treści i bez pliku biblioteki — nie ma czego załączyć")
	}

	zapisany, err := a.repozytorium.ZapiszZalacznik(ctx, zalacznik)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.ResearchSourceAttachmentAddResponse{}, bladNieznanegoZrodlaBadania(z.SourceId)
	}
	if err != nil {
		return shared.ResearchSourceAttachmentAddResponse{}, bladBadan(err)
	}
	return shared.ResearchSourceAttachmentAddResponse{Attachment: zlozZalacznikBadania(zapisany)}, nil
}

// WypiszZalacznikiZrodla obsługuje `research.source.attachment.list`.
// `missingKinds` nazywa braki kompletności — po to menedżer załączników istnieje.
func (a *adapterBadan) WypiszZalacznikiZrodla(ctx context.Context,
	z shared.ResearchSourceAttachmentListRequest) (shared.ResearchSourceAttachmentListResponse, error) {

	if z.SourceId == "" {
		return shared.ResearchSourceAttachmentListResponse{},
			bladWskazaniaBadan("source.attachment.list bez źródła")
	}
	zalaczniki, err := a.repozytorium.Zalaczniki(ctx, z.SourceId)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.ResearchSourceAttachmentListResponse{}, bladNieznanegoZrodlaBadania(z.SourceId)
	}
	if err != nil {
		return shared.ResearchSourceAttachmentListResponse{}, bladBadan(err)
	}

	obecne := map[string]bool{}
	przelozone := make([]shared.ResearchAttachment, 0, len(zalaczniki))
	for _, zalacznik := range zalaczniki {
		obecne[zalacznik.Rodzaj] = true
		przelozone = append(przelozone, zlozZalacznikBadania(zalacznik))
	}
	brakujace := []string{}
	for _, rodzaj := range []string{
		string(shared.ResearchAttachmentKindFulltext),
		string(shared.ResearchAttachmentKindSnapshot),
	} {
		if !obecne[rodzaj] {
			brakujace = append(brakujace, rodzaj)
		}
	}
	return shared.ResearchSourceAttachmentListResponse{
		Attachments: przelozone, MissingKinds: brakujace,
	}, nil
}

// zlozZalacznikBadania przekłada wiersz repozytorium na byt kontraktu wraz
// z rozmiarem i chwilą utworzenia w postaci wymaganej kontraktem.
func zlozZalacznikBadania(z dane.ZalacznikZrodlaBadania) shared.ResearchAttachment {
	return shared.ResearchAttachment{
		Id: z.Kod, SourceId: z.ZrodloKod, Kind: shared.ResearchAttachmentKind(z.Rodzaj),
		LibraryFileId: z.PlikBibliotekiID, SizeBytes: z.RozmiarBajtow,
		CreatedAt: chwilaBazy(z.Utworzono),
	}
}

// WczytajBibliografie obsługuje komendę `research.source.import`: wczytuje
// referencje z pliku wskazanego formatu i zapisuje je jako źródła okna.
func (a *adapterBadan) WczytajBibliografie(ctx context.Context,
	z shared.ResearchSourceImportRequest) (shared.ResearchSourceImportResponse, error) {

	if z.WindowId == "" {
		return shared.ResearchSourceImportResponse{}, bladWskazaniaBadan("source.import bez okna badania")
	}
	bajty, err := trescZadaniaBadania(z.ContentBase64, z.SourcePath)
	if err != nil {
		return shared.ResearchSourceImportResponse{}, err
	}
	if bajty == nil {
		return shared.ResearchSourceImportResponse{},
			bladWskazaniaBadan("source.import bez treści i bez ścieżki — nie ma czego wczytać")
	}

	pozycje, pominieteZParsera, err := pozycjeBibliografiiBadania(z.Format, string(bajty))
	if err != nil {
		return shared.ResearchSourceImportResponse{}, err
	}

	wiarygodnosc := shared.ResearchCredibility("")
	if z.DefaultCredibility != nil {
		wiarygodnosc = *z.DefaultCredibility
	}

	zapisane := []shared.ResearchSource{}
	powody := append([]string{}, pominieteZParsera...)
	for _, pozycja := range pozycje {
		if strings.TrimSpace(pozycja.Tytul) == "" {
			powody = append(powody, "pozycja bez tytułu — nie da się jej pokazać w katalogu")
			continue
		}
		zrodlo := dane.ZrodloBadania{
			Kod: nowyIdentyfikator(przedrostekZrodlaBadania), Okno: z.WindowId,
			Tytul: pozycja.Tytul, Rodzaj: shared.ResearchSourceKindDocument,
			Adres: pozycja.Adres, Wiarygodnosc: wiarygodnosc,
		}
		if pozycja.Pochodzenie != nil {
			zrodlo.Pochodzenie = pozycja.Pochodzenie
		}
		wynik, err := a.repozytorium.ZapiszZrodlo(ctx, zrodlo)
		if err != nil {
			powody = append(powody, "pozycji „"+pozycja.Tytul+"” nie udało się zapisać: "+err.Error())
			continue
		}
		if pozycja.Identyfikator != nil {
			_ = a.repozytorium.UstawLektureZrodla(ctx, wynik.Kod,
				dane.LekturaZrodlaBadania{Identyfikat: pozycja.Identyfikator})
		}
		zapisane = append(zapisane, zlozZrodloBadania(wynik))
	}
	return shared.ResearchSourceImportResponse{
		Sources: zapisane, Imported: len(zapisane),
		Skipped: len(pozycje) + len(pominieteZParsera) - len(zapisane), SkipReasons: powody,
	}, nil
}

// pozycjaBibliografiiBadania jest jedną referencją wczytaną z pliku, sprowadzoną
// do pól, które moduł naprawdę zapisuje.
type pozycjaBibliografiiBadania struct {
	Tytul         string
	Adres         *string
	Identyfikator *string
	Pochodzenie   *string
}

// pozycjeBibliografiiBadania rozkłada plik bibliografii na referencje,
// dobierając parser po formacie wskazanym w żądaniu importu.
func pozycjeBibliografiiBadania(format shared.ResearchImportFormat,
	tresc string) ([]pozycjaBibliografiiBadania, []string, error) {

	switch format {
	case shared.ResearchImportFormatBibtex:
		return pozycjeBibtexBadania(tresc), nil, nil
	case shared.ResearchImportFormatRis:
		return pozycjeRisBadania(tresc), nil, nil
	case shared.ResearchImportFormatCslJson:
		return pozycjeCslBadania(tresc)
	case shared.ResearchImportFormatEndnoteXml:
		return pozycjeEndnoteBadania(tresc)
	case shared.ResearchImportFormatCsv:
		return pozycjeCsvBadania(tresc)
	default:
		return nil, nil, bladWskazaniaBadan("format importu „" + string(format) +
			"” nie jest jednym z: bibtex, ris, cslJson, endnoteXml, csv")
	}
}

// pozycjeBibtexBadania czyta wpisy BibTeX. Parser bierze pola, które moduł
// zapisuje, i nie udaje pełnego czytnika gramatyki BibTeX-a: pole nieznane
// pomija zamiast wywracać cały import przez jeden nawias.
func pozycjeBibtexBadania(tresc string) []pozycjaBibliografiiBadania {
	pozycje := []pozycjaBibliografiiBadania{}
	for _, blok := range strings.Split(tresc, "@")[1:] {
		pola := map[string]string{}
		for _, wiersz := range strings.Split(blok, "\n") {
			czesci := strings.SplitN(wiersz, "=", 2)
			if len(czesci) != 2 {
				continue
			}
			klucz := strings.ToLower(strings.TrimSpace(czesci[0]))
			wartosc := strings.Trim(strings.TrimSpace(czesci[1]), "{}\",")
			wartosc = strings.Trim(wartosc, "{} \t\",")
			if klucz != "" && wartosc != "" {
				pola[klucz] = wartosc
			}
		}
		if len(pola) == 0 {
			continue
		}
		pozycje = append(pozycje, pozycjaZPolBadania(pola["title"], pola["url"], pola["doi"],
			pola["journal"]))
	}
	return pozycje
}

// pozycjeRisBadania czyta wpisy formatu RIS wiersz po wierszu; rekord kończy
// znacznik `ER`, po którym pola zebrane dotąd składają się w referencję.
func pozycjeRisBadania(tresc string) []pozycjaBibliografiiBadania {
	pozycje := []pozycjaBibliografiiBadania{}
	pola := map[string]string{}
	for _, wiersz := range strings.Split(tresc, "\n") {
		wiersz = strings.TrimRight(wiersz, "\r")
		if len(wiersz) < 6 || wiersz[2:6] != "  - " {
			if strings.HasPrefix(strings.TrimSpace(wiersz), "ER") {
				if len(pola) > 0 {
					pozycje = append(pozycje, pozycjaZPolBadania(pola["TI"], pola["UR"],
						pola["DO"], pola["JO"]))
				}
				pola = map[string]string{}
			}
			continue
		}
		pola[wiersz[:2]] = strings.TrimSpace(wiersz[6:])
	}
	if len(pola) > 0 {
		pozycje = append(pozycje, pozycjaZPolBadania(pola["TI"], pola["UR"], pola["DO"], pola["JO"]))
	}
	return pozycje
}

// pozycjeCslBadania czyta listę dokumentów formatu CSL-JSON i przekłada
// wskazane pola każdego wpisu na referencję bibliografii.
func pozycjeCslBadania(tresc string) ([]pozycjaBibliografiiBadania, []string, error) {
	var wpisy []struct {
		Title         string `json:"title"`
		URL           string `json:"URL"`
		DOI           string `json:"DOI"`
		ContainerName string `json:"container-title"`
	}
	if err := json.Unmarshal([]byte(tresc), &wpisy); err != nil {
		return nil, nil, bladWskazaniaBadan("plik CSL-JSON jest nieczytelny: " + err.Error())
	}
	pozycje := make([]pozycjaBibliografiiBadania, 0, len(wpisy))
	for _, wpis := range wpisy {
		pozycje = append(pozycje, pozycjaZPolBadania(wpis.Title, wpis.URL, wpis.DOI, wpis.ContainerName))
	}
	return pozycje, nil, nil
}

// pozycjeEndnoteBadania czyta plik formatu EndNote XML i przekłada rekordy
// jego drzewa znaczników na referencje bibliografii.
func pozycjeEndnoteBadania(tresc string) ([]pozycjaBibliografiiBadania, []string, error) {
	var zbior struct {
		Rekordy []struct {
			Tytuly struct {
				Tytul string `xml:"title"`
			} `xml:"titles"`
			Adresy struct {
				Odnosniki struct {
					Adres string `xml:"url"`
				} `xml:"related-urls"`
			} `xml:"urls"`
			Numery []struct {
				Rodzaj string `xml:"type,attr"`
				Tresc  string `xml:",chardata"`
			} `xml:"electronic-resource-num"`
		} `xml:"records>record"`
	}
	if err := xml.Unmarshal([]byte(tresc), &zbior); err != nil {
		return nil, nil, bladWskazaniaBadan("plik EndNote XML jest nieczytelny: " + err.Error())
	}
	pozycje := make([]pozycjaBibliografiiBadania, 0, len(zbior.Rekordy))
	for _, rekord := range zbior.Rekordy {
		doi := ""
		if len(rekord.Numery) > 0 {
			doi = strings.TrimSpace(rekord.Numery[0].Tresc)
		}
		pozycje = append(pozycje, pozycjaZPolBadania(strings.TrimSpace(rekord.Tytuly.Tytul),
			strings.TrimSpace(rekord.Adresy.Odnosniki.Adres), doi, ""))
	}
	return pozycje, nil, nil
}

// pozycjeCsvBadania czyta arkusz CSV z wierszem nagłówka, dopasowując
// kolumny referencji po nazwie nagłówka w kilku uznanych wariantach.
func pozycjeCsvBadania(tresc string) ([]pozycjaBibliografiiBadania, []string, error) {
	czytnik := csv.NewReader(strings.NewReader(tresc))
	czytnik.FieldsPerRecord = -1
	wiersze, err := czytnik.ReadAll()
	if err != nil {
		return nil, nil, bladWskazaniaBadan("plik CSV jest nieczytelny: " + err.Error())
	}
	if len(wiersze) < 2 {
		return nil, []string{"plik CSV nie ma ani jednego wiersza danych pod nagłówkiem"}, nil
	}
	kolumny := map[string]int{}
	for numer, nazwa := range wiersze[0] {
		kolumny[strings.ToLower(strings.TrimSpace(nazwa))] = numer
	}
	pobierz := func(wiersz []string, nazwy ...string) string {
		for _, nazwa := range nazwy {
			if numer, jest := kolumny[nazwa]; jest && numer < len(wiersz) {
				return strings.TrimSpace(wiersz[numer])
			}
		}
		return ""
	}
	pozycje := make([]pozycjaBibliografiiBadania, 0, len(wiersze)-1)
	for _, wiersz := range wiersze[1:] {
		pozycje = append(pozycje, pozycjaZPolBadania(
			pobierz(wiersz, "title", "tytul", "tytuł"),
			pobierz(wiersz, "url", "adres"),
			pobierz(wiersz, "doi", "identyfikator"),
			pobierz(wiersz, "journal", "source", "zrodlo", "źródło")))
	}
	return pozycje, nil, nil
}

// pozycjaZPolBadania składa referencję z czterech pól wspólnych wszystkim
// obsługiwanym formatom, pomijając pola puste zamiast zapisywać je jako takie.
func pozycjaZPolBadania(tytul, adres, identyfikator, pochodzenie string) pozycjaBibliografiiBadania {
	pozycja := pozycjaBibliografiiBadania{Tytul: strings.TrimSpace(tytul)}
	if wartosc := strings.TrimSpace(adres); wartosc != "" {
		pozycja.Adres = &wartosc
	}
	if wartosc := strings.TrimSpace(identyfikator); wartosc != "" {
		pozycja.Identyfikator = &wartosc
	}
	if wartosc := strings.TrimSpace(pochodzenie); wartosc != "" {
		pozycja.Pochodzenie = &wartosc
	}
	return pozycja
}

// PrzechwycStrone obsługuje `research.source.capture`. Strona wchodzi jako
// źródło z treścią odczytaną w trybie lektury, a przy trybie migawki także
// z bajtami dokumentu odłożonymi w magazynie.
func (a *adapterBadan) PrzechwycStrone(ctx context.Context,
	z shared.ResearchSourceCaptureRequest) (shared.ResearchSourceCaptureResponse, error) {

	if z.WindowId == "" {
		return shared.ResearchSourceCaptureResponse{}, bladWskazaniaBadan("source.capture bez okna badania")
	}
	if strings.TrimSpace(z.Url) == "" {
		return shared.ResearchSourceCaptureResponse{}, bladWskazaniaBadan("source.capture bez adresu")
	}

	adres := z.Url
	if z.UseWayback != nil && *z.UseWayback {
		if archiwalny, err := adresArchiwalnyBadania(ctx, z.Url); err == nil && archiwalny != "" {
			adres = archiwalny
		}
	}
	bajty, _, err := pobierzBadania(ctx, adres)
	if err != nil {
		return shared.ResearchSourceCaptureResponse{}, err
	}
	tytul, tekst := tekstZDokumentuHtmlBadania(string(bajty))
	if strings.TrimSpace(tytul) == "" {
		tytul = z.Url
	}

	zapisany, err := a.repozytorium.ZapiszZrodlo(ctx, dane.ZrodloBadania{
		Kod: nowyIdentyfikator(przedrostekZrodlaBadania), Okno: z.WindowId,
		Tytul: tytul, Rodzaj: shared.ResearchSourceKindWeb, Adres: &z.Url,
	})
	if err != nil {
		return shared.ResearchSourceCaptureResponse{}, bladBadan(err)
	}
	if err := a.repozytorium.UstawTrescZrodla(ctx, zapisany.Kod, dane.TrescZrodlaBadania{
		Tekst: &tekst, WarstwaTekstu: true,
	}); err != nil {
		return shared.ResearchSourceCaptureResponse{}, bladBadan(err)
	}

	odpowiedz := shared.ResearchSourceCaptureResponse{Source: zlozZrodloBadania(zapisany)}
	if z.Mode == shared.ResearchCaptureModeSnapshot || z.Mode == shared.ResearchCaptureModeBoth {
		sciezka, rozmiar, err := a.odlozMaterialBadania(bajty)
		if err != nil {
			return shared.ResearchSourceCaptureResponse{}, err
		}
		zalacznik, err := a.repozytorium.ZapiszZalacznik(ctx, dane.ZalacznikZrodlaBadania{
			Kod: nowyIdentyfikator(przedrostekZalacznikaBadania), ZrodloKod: zapisany.Kod,
			Rodzaj:  string(shared.ResearchAttachmentKindSnapshot),
			Sciezka: &sciezka, RozmiarBajtow: &rozmiar,
		})
		if err != nil {
			return shared.ResearchSourceCaptureResponse{}, bladBadan(err)
		}
		odpowiedz.SnapshotAttachmentId = &zalacznik.Kod
	}
	return odpowiedz, nil
}

// adresArchiwalnyBadania pyta Wayback Machine o wersję archiwalną strony.
// Brak wersji nie jest usterką: przechwycenie idzie wtedy adresem żywym.
func adresArchiwalnyBadania(ctx context.Context, adres string) (string, error) {
	bajty, _, err := pobierzBadania(ctx, "https://archive.org/wayback/available?url="+adres)
	if err != nil {
		return "", err
	}
	var odpowiedz struct {
		ArchivedSnapshots struct {
			Closest struct {
				URL string `json:"url"`
			} `json:"closest"`
		} `json:"archived_snapshots"`
	}
	if err := json.Unmarshal(bajty, &odpowiedz); err != nil {
		return "", err
	}
	return odpowiedz.ArchivedSnapshots.Closest.URL, nil
}

// PrzepiszNagranie obsługuje `research.source.transcribe`. Transkrypcja idzie
// silnikiem rozpoznania mowy rdzenia — tym samym, którym jedzie Assistant.
func (a *adapterBadan) PrzepiszNagranie(ctx context.Context,
	z shared.ResearchSourceTranscribeRequest) (shared.ResearchSourceTranscribeResponse, error) {

	if z.WindowId == "" {
		return shared.ResearchSourceTranscribeResponse{},
			bladWskazaniaBadan("source.transcribe bez okna badania")
	}
	if a.mowa == nil {
		return shared.ResearchSourceTranscribeResponse{},
			protokolBladBadania(shared.ErrorCodeInternalError,
				"silnik rozpoznania mowy nie jest wpięty — naprawa: podpiąć port mowy "+
					"przy składaniu rdzenia")
	}
	odnosnik := ""
	switch {
	case z.SourcePath != nil && strings.TrimSpace(*z.SourcePath) != "":
		odnosnik = *z.SourcePath
	case z.LibraryFileId != nil && strings.TrimSpace(*z.LibraryFileId) != "":
		odnosnik = *z.LibraryFileId
	default:
		return shared.ResearchSourceTranscribeResponse{},
			bladWskazaniaBadan("source.transcribe bez nagrania — wskaż plik biblioteki albo ścieżkę")
	}

	zadanie := shared.SpeechTranscribeRequest{AudioRef: odnosnik}
	if z.Language != nil {
		zadanie.Language = z.Language
	}
	wynik, err := a.mowa.Przepisz(ctx, zadanie)
	if err != nil {
		return shared.ResearchSourceTranscribeResponse{}, err
	}

	tytul := "Transkrypcja: " + odnosnik
	zapisany, err := a.repozytorium.ZapiszZrodlo(ctx, dane.ZrodloBadania{
		Kod: nowyIdentyfikator(przedrostekZrodlaBadania), Okno: z.WindowId,
		Tytul: tytul, Rodzaj: shared.ResearchSourceKindNote, PlikBibliotekiID: z.LibraryFileId,
	})
	if err != nil {
		return shared.ResearchSourceTranscribeResponse{}, bladBadan(err)
	}
	tekst := wynik.Transcript
	if err := a.repozytorium.UstawTrescZrodla(ctx, zapisany.Kod, dane.TrescZrodlaBadania{
		Tekst: &tekst, WarstwaTekstu: true,
	}); err != nil {
		return shared.ResearchSourceTranscribeResponse{}, bladBadan(err)
	}
	// Odcinki liczone z akapitów transkryptu, nie ze znaków.
	odcinki := len(strings.Split(strings.TrimSpace(tekst), "\n"))
	if strings.TrimSpace(tekst) == "" {
		odcinki = 0
	}
	return shared.ResearchSourceTranscribeResponse{
		Source: zlozZrodloBadania(zapisany), SegmentCount: odcinki,
	}, nil
}

// RozstrzygnijIdentyfikator obsługuje komendę `research.source.resolve`:
// pyta serwis właściwy rodzajowi identyfikatora i oddaje jego opis pracy.
func (a *adapterBadan) RozstrzygnijIdentyfikator(ctx context.Context,
	z shared.ResearchSourceResolveRequest) (shared.ResearchSourceResolveResponse, error) {

	identyfikator := strings.TrimSpace(z.Identifier)
	if identyfikator == "" {
		return shared.ResearchSourceResolveResponse{}, bladWskazaniaBadan("source.resolve bez identyfikatora")
	}
	rodzaj := shared.ResearchIdentifierKind(shared.ResearchIdentifierKindDoi)
	if z.Kind != nil && *z.Kind != "" {
		rodzaj = *z.Kind
	} else {
		rodzaj = rodzajIdentyfikatoraBadania(identyfikator)
	}

	switch rodzaj {
	case shared.ResearchIdentifierKindDoi:
		praca, err := pracaCrossref(ctx, identyfikator)
		if err != nil {
			return shared.ResearchSourceResolveResponse{}, err
		}
		return shared.ResearchSourceResolveResponse{Result: praca.jakoWynikBadania()}, nil
	case shared.ResearchIdentifierKindIsbn:
		wynik, err := rozstrzygnijISBN(ctx, identyfikator)
		if err != nil {
			return shared.ResearchSourceResolveResponse{}, err
		}
		return shared.ResearchSourceResolveResponse{Result: wynik}, nil
	case shared.ResearchIdentifierKindPmid:
		wynik, err := rozstrzygnijPMID(ctx, identyfikator)
		if err != nil {
			return shared.ResearchSourceResolveResponse{}, err
		}
		return shared.ResearchSourceResolveResponse{Result: wynik}, nil
	case shared.ResearchIdentifierKindArxiv:
		wyniki, err := szukajArxiv(ctx, "id:"+identyfikator, 1)
		if err != nil {
			return shared.ResearchSourceResolveResponse{}, err
		}
		if len(wyniki) == 0 {
			return shared.ResearchSourceResolveResponse{},
				protokolBladBadania(shared.ErrorCodeNotFound, "arXiv nie zna pozycji "+identyfikator)
		}
		return shared.ResearchSourceResolveResponse{Result: wyniki[0]}, nil
	default:
		return shared.ResearchSourceResolveResponse{}, bladWskazaniaBadan(
			"rodzaj identyfikatora „" + string(rodzaj) + "” nie jest jednym z: doi, isbn, pmid, arxiv")
	}
}

// rodzajIdentyfikatoraBadania rozpoznaje rodzaj po kształcie samego napisu —
// żądanie bez wskazania rodzaju nie musi być odmową, skoro postać identyfikatora
// rozstrzyga sama.
func rodzajIdentyfikatoraBadania(identyfikator string) shared.ResearchIdentifierKind {
	czysty := strings.TrimSpace(strings.ToLower(identyfikator))
	switch {
	case strings.HasPrefix(czysty, "10."):
		return shared.ResearchIdentifierKindDoi
	case strings.HasPrefix(czysty, "arxiv:") || strings.Contains(czysty, "arxiv.org"):
		return shared.ResearchIdentifierKindArxiv
	case len(cyfryBadania(czysty)) >= 10 && len(cyfryBadania(czysty)) <= 13 &&
		len(cyfryBadania(czysty)) == len(strings.ReplaceAll(czysty, "-", "")):
		return shared.ResearchIdentifierKindIsbn
	default:
		if _, err := strconv.Atoi(czysty); err == nil {
			return shared.ResearchIdentifierKindPmid
		}
		return shared.ResearchIdentifierKindDoi
	}
}

// cyfryBadania oddaje same cyfry napisu, pomijając pozostałe znaki, co służy
// do liczenia długości identyfikatora niezależnie od jego separatorów.
func cyfryBadania(tekst string) string {
	var zebrane strings.Builder
	for _, znak := range tekst {
		if '0' <= znak && znak <= '9' {
			zebrane.WriteRune(znak)
		}
	}
	return zebrane.String()
}

// ── Pomocniki wspólne rodziny ──────────────────────────────────────────────

// trescZadaniaBadania oddaje bajty wskazane treścią wprost albo ścieżką.
// Brak obu wskazań daje `nil` bez usterki — wywołujący rozstrzyga, czy brak
// treści jest u niego dopuszczalny.
func trescZadaniaBadania(base64Tresc, sciezka *string) ([]byte, error) {
	if base64Tresc != nil && strings.TrimSpace(*base64Tresc) != "" {
		bajty, err := base64.StdEncoding.DecodeString(strings.TrimSpace(*base64Tresc))
		if err != nil {
			return nil, bladWskazaniaBadan("treść nie jest poprawnym base64: " + err.Error())
		}
		return bajty, nil
	}
	if sciezka != nil && strings.TrimSpace(*sciezka) != "" {
		bajty, err := os.ReadFile(strings.TrimSpace(*sciezka))
		if err != nil {
			return nil, protokolBladBadania(shared.ErrorCodeValidationFailed,
				"nie można odczytać pliku "+*sciezka+": "+err.Error())
		}
		return bajty, nil
	}
	return nil, nil
}

// zawieraNapisBadania mówi, czy lista niesie wskazany napis, porównując
// wartości bez rozróżniania wielkości liter i otaczających odstępów.
func zawieraNapisBadania(lista []string, szukany string) bool {
	for _, wartosc := range lista {
		if strings.EqualFold(strings.TrimSpace(wartosc), strings.TrimSpace(szukany)) {
			return true
		}
	}
	return false
}

// posortowaneNapisyBadania oddaje kopię listy w kolejności ustalonej — odpowiedź
// ma wyglądać tak samo przy dwóch identycznych żądaniach.
func posortowaneNapisyBadania(lista []string) []string {
	kopia := append([]string{}, lista...)
	sort.Strings(kopia)
	return kopia
}

// bladNieznanegoZrodlaBadania odróżnia odpowiedź „źródła nie ma" od usterki
// zapisu, nazywając kod nieznanego źródła w treści odmowy.
func bladNieznanegoZrodlaBadania(kod string) error {
	return protokolBladBadania(shared.ErrorCodeNotFound, "źródło nie istnieje: "+kod)
}
