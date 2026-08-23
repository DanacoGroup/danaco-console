// Odpowiedzialność pliku: moduł Research — typ adaptera, konstruktor oraz
// katalogowanie źródeł i zapis ustaleń (Sources Manager, Findings Tracker).
// Raport, eksport i przestrzeń badania stoją w
// `adapter_modul_badania_raport.go` i `adapter_modul_badania_uchwyty.go`, gdzie
// deklarowany jest port `Badania` i `zarejestrujBadania`.
//
// Typ adaptera i konstruktor stoją tutaj, bo metody pisane w trzech plikach
// wiszą na tym samym `*adapterBadan` (jak `dane/badania*.go` na
// `*repozytoriumBadan`) — deklaracja typu w dwóch miejscach byłaby dwiema
// prawdami o jednym bycie.
//
// Rdzeń nie ocenia wiarygodności sam. `research.source.add` zapisuje ocenę,
// którą podał Operator w żądaniu — nie wylicza jej z adresu ani z treści, i
// nie sięga do sieci po źródło. Brak oceny w żądaniu zostaje pusty; warstwa
// danych nadaje mu wtedy wartość domyślną `unverified` (brak oceny
// nie udaje oceny wyliczonej).
package core

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"

	"danacoconsole/server/internal/konfiguracja"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Przedrostki identyfikatorów bytów modułu. Zadeklarowane tu w całości —
// łącznie z tymi, których ten plik nie używa — żeby pozostałe pliki modułu nie
// deklarowały ich po raz drugi.
const (
	przedrostekZrodlaBadania    = "zrod-"
	przedrostekUstaleniaBadania = "ust-"
	przedrostekRaportuBadania   = "rap-"
	przedrostekSekcjiRaportu    = "sekc-"
	przedrostekEksportuRaportu  = "eksp-"

	// Przedrostki bytów dobudowanych migracją 150.
	przedrostekZalacznikaBadania    = "zal-"
	przedrostekAdnotacjiBadania     = "adn-"
	przedrostekTabeliBadania        = "tab-"
	przedrostekKoduBadania          = "kod-"
	przedrostekSprzecznosciBadania  = "sprz-"
	przedrostekWatkuBadania         = "wat-"
	przedrostekPytaniaBadania       = "pyt-"
	przedrostekMonitoraBadania      = "mon-"
	przedrostekWersjiRaportuBadania = "wers-"
	przedrostekKomentarzaBadania    = "kom-"
	przedrostekBlokuBadania         = "blok-"
	przedrostekSzablonuBadania      = "szab-"
	przedrostekStyluBadania         = "styl-"
	przedrostekPartiiBadania        = "part-"
	przedrostekUdostepnieniaBadania = "udos-"
)

// adapterBadan wypełnia część portu Badania. Zależność jest jedna: wspólne
// repozytorium modułu, rozłożone po stronie danych na trzy pliki wedle
// odpowiedzialności, ale niosące jeden typ.
type adapterBadan struct {
	repozytorium dane.RepozytoriumBadan
	// magazyn trzyma BAJTY modułu: migawki stron, załączniki źródeł i pliki
	// eksportu. Ten sam magazyn treści, którym jadą Design i Library — moduł nie
	// zakłada trzeciego miejsca na te same bajty, tylko własny podkatalog.
	magazyn *magazynTresciBiblioteki
	// dokumenty jest portem arsenału dokumentowego (Pandoc, poppler, Tesseract).
	// Lektura źródła, rozpoznanie pisma i zamiana formatu eksportu idą przez
	// niego — moduł nie uruchamia ani jednego procesu sam.
	dokumenty Dokumenty
	// mowa jest silnikiem rozpoznania mowy; przez niego idzie transkrypcja
	// nagrania na źródło cytowalne. Silnik jest jeden na rdzeń, wspólny
	// z modułem Assistant.
	mowa Mowa
	// katalogDanych jest korzeniem magazynów rdzenia. Trzymany osobno, bo
	// odwołania zasobów wychodzące z innych portów są ścieżkami WZGLĘDNYMI tego
	// korzenia (`design/zasoby/...`) — bez niego eksport nie miałby jak sięgnąć
	// po bajty, które wytworzył mu arsenał.
	katalogDanych string
	// kanaly jest rejestrem kanałów modelu rdzenia, wpiętym wzorem Roundtable
	// (`ZKanalami`). Operacje badania, które wymagają modelu — streszczanie
	// źródła, porównanie źródeł, redakcja sekcji raportu ze zebranych ustaleń —
	// idą przez ten rejestr. Bez niego taka operacja odmawia wprost.
	kanaly *models.Rejestr
}

// nowyAdapterBadan wiąże adapter z repozytorium modułu.
func nowyAdapterBadan(repozytorium dane.RepozytoriumBadan) *adapterBadan {
	return &adapterBadan{
		repozytorium:  repozytorium,
		magazyn:       magazynZasobowBadania(konfiguracja.KatalogDanychDomyslny()),
		katalogDanych: konfiguracja.KatalogDanychDomyslny(),
	}
}

// ZKatalogiemDanych przestawia magazyn bajtów modułu na katalog wskazany
// konfiguracją — ten sam, którym jadą magazyn biblioteki i zasoby Designu.
func (a *adapterBadan) ZKatalogiemDanych(katalog string) *adapterBadan {
	if strings.TrimSpace(katalog) != "" {
		a.magazyn = magazynZasobowBadania(katalog)
		a.katalogDanych = katalog
	}
	return a
}

// ZDokumentami wpina port arsenału dokumentowego: lektura źródła, rozpoznanie
// pisma i zamiana formatu eksportu. Bez niego te czynności odmawiają wprost.
func (a *adapterBadan) ZDokumentami(d Dokumenty) *adapterBadan {
	a.dokumenty = d
	return a
}

// ZMowa wpina silnik rozpoznania mowy — transkrypcja nagrania na źródło.
func (a *adapterBadan) ZMowa(m Mowa) *adapterBadan {
	a.mowa = m
	return a
}

// magazynZasobowBadania składa magazyn bajtów modułu nad katalogiem danych.
// Osobny podkatalog, bo moduły nie dzielą stanu: sprzątnięcie materiałów
// badania nie ma prawa ruszyć treści biblioteki ani zasobów Designu.
func magazynZasobowBadania(katalogDanych string) *magazynTresciBiblioteki {
	if strings.TrimSpace(katalogDanych) == "" {
		return nil
	}
	return &magazynTresciBiblioteki{
		katalog: filepath.Join(katalogDanych, podkatalogBadania, podkatalogMaterialowBadania),
	}
}

// Podkatalogi magazynu modułu.
const (
	podkatalogBadania           = "badania"
	podkatalogMaterialowBadania = "materialy"
)

// odczytajPlikBadania oddaje bajty spod ścieżki, sprowadzając wskazanie względne
// do korzenia magazynów rdzenia. Odwołania wychodzące z portów są względne
// (granica kontraktu: rdzeń nie wypuszcza ścieżek z dysku Operatora), więc
// otwarcie ich wprost trafiałoby w nieistniejący plik.
func (a *adapterBadan) odczytajPlikBadania(sciezka string) ([]byte, error) {
	pelna := sciezka
	if !filepath.IsAbs(pelna) {
		pelna = filepath.Join(a.katalogDanych, sciezka)
	}
	bajty, err := os.ReadFile(pelna)
	if err != nil {
		return nil, protokolBladBadania(shared.ErrorCodeInternalError,
			"nie można odczytać pliku "+sciezka+": "+err.Error())
	}
	return bajty, nil
}

// protokolBladBadania buduje odmowę modułu o wskazanym kodzie kontraktu.
// Jedno miejsce składania zdania odmowy: przedrostek „moduł Research" jest ten
// sam we wszystkich rodzinach, więc Operator poznaje nadawcę po pierwszym słowie.
func protokolBladBadania(kod shared.ErrorCode, powod string) error {
	return protocol.JakoError(protocol.NowyBlad(kod, "moduł Research: "+powod))
}

// odlozMaterialBadania utrwala bajty w magazynie modułu pod sumą ich treści
// i oddaje ścieżkę oraz rozmiar. Nazwa jest sumą zawartości, więc ta sama
// migawka wgrana dwa razy leży w magazynie raz.
func (a *adapterBadan) odlozMaterialBadania(bajty []byte) (string, int64, error) {
	if a.magazyn == nil {
		return "", 0, protokolBladBadania(shared.ErrorCodeInternalError,
			"magazyn materiałów badania nie jest wpięty — naprawa: podać katalog danych "+
				"przy składaniu rdzenia")
	}
	suma := sha256.Sum256(bajty)
	sciezka, err := a.magazyn.Zapisz(bajty, hex.EncodeToString(suma[:]))
	if err != nil {
		return "", 0, bladBadan(err)
	}
	return sciezka, int64(len(bajty)), nil
}

// ZKanalami wpina rejestr kanałów modelu (wzór: `adapter_modul_roundtable.go`).
// Zwraca adapter, żeby montaż wiązał zależność w łańcuchu. Katalogowanie źródeł
// i zapis ustaleń modelu nie potrzebują — most służy operacjom słownym, które
// dopiero kontrakt wskaże kanałem.
func (a *adapterBadan) ZKanalami(kanaly *models.Rejestr) *adapterBadan {
	a.kanaly = kanaly
	return a
}

// zapytajModel woła wskazany kanał modelu i zbiera całą odpowiedź tekstową —
// wspólny most modułu do modelu dla operacji słownych (streszczenie, porównanie
// źródeł, redakcja sekcji). Kanał wskazuje wywołujący; brak rejestru albo pusty
// kanał to odmowa wprost, nie streszczenie zmyślone bez modelu.
//
// Granica kontraktu: żądania modułu Research (`research.report.build` i inne)
// nie niosą `channelId`, więc wywołujący nie ma czym wskazać kanału z żądania.
// Dopóki kontrakt tego pola nie ma, operacja bierze domyślny czynny kanał
// rejestru przez `domyslnyKanalBadania`.
func (a *adapterBadan) zapytajModel(ctx context.Context, okno, kanal, tresc string) (string, error) {
	if a.kanaly == nil {
		return "", bladBrakuKanalowBadan()
	}
	if strings.TrimSpace(kanal) == "" {
		return "", bladWskazaniaBadan("żądanie bez wskazania kanału modelu — rdzeń nie zgaduje kanału badania")
	}
	var zebrane strings.Builder
	ujscie := models.UjscieFunkcji(func(_ context.Context, f models.Fragment) error {
		if f.Kind == shared.ChunkKindText {
			zebrane.WriteString(models.TrescFragmentu(f))
		}
		return nil
	})
	zapytanie := models.Zapytanie{
		Zasiegi:   models.Zasiegi{Okno: okno},
		Wiadomosc: okno,
		Tresc:     tresc,
		Kanal:     kanal,
	}
	if err := a.kanaly.Wyslij(ctx, zapytanie, ujscie); err != nil {
		return "", bladBadan(err)
	}
	return zebrane.String(), nil
}

// domyslnyKanalBadania wskazuje kanał operacji słownej badania. Żądania modułu
// Research nie niosą `channelId`, więc most bierze pierwszy czynny kanał
// rejestru — ten sam, który `channel.list` pokazuje jako
// gotowy (`Kontrakt(true)`). Brak wpiętego rejestru albo brak czynnego kanału
// to odmowa wprost, nie streszczenie zmyślone bez modelu.
func (a *adapterBadan) domyslnyKanalBadania() (string, error) {
	if a.kanaly == nil {
		return "", bladBrakuKanalowBadan()
	}
	czynne := a.kanaly.Kontrakt(true)
	if len(czynne) == 0 {
		return "", bladBrakuCzynnegoKanaluBadan()
	}
	return czynne[0].Id, nil
}

// bladBrakuKanalowBadan znakuje odmowę, gdy operacja badania ma iść modelem,
// a rejestr kanałów nie został wpięty — kanał niedostępny, nie usterka rdzenia.
func bladBrakuKanalowBadan() error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
		"moduł Research: rejestr kanałów modelu nie jest wpięty — operacja badania nie ma czym wołać modelu"))
}

// bladBrakuCzynnegoKanaluBadan znakuje odmowę, gdy rejestr jest wpięty, lecz nie
// ma ani jednego czynnego kanału — redakcja sekcji nie ma czym wołać modelu,
// więc odmawia wprost zamiast składać streszczenie bez modelu.
func bladBrakuCzynnegoKanaluBadan() error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
		"moduł Research: rejestr kanałów nie ma czynnego kanału — redakcja sekcji raportu nie ma czym wołać modelu"))
}

// DodajZrodlo obsługuje `research.source.add`. Katalog źródła i nic więcej —
// patrz nagłówek pliku o uczciwości oceny wiarygodności.
func (a *adapterBadan) DodajZrodlo(ctx context.Context,
	z shared.ResearchSourceAddRequest) (shared.ResearchSourceAddResponse, error) {

	if z.WindowId == "" {
		return shared.ResearchSourceAddResponse{}, bladWskazaniaBadan("żądanie bez okna badania")
	}
	if z.Title == "" {
		return shared.ResearchSourceAddResponse{}, bladWskazaniaBadan("żądanie bez tytułu źródła")
	}

	zrodlo := dane.ZrodloBadania{
		Kod:              nowyIdentyfikator(przedrostekZrodlaBadania),
		Okno:             z.WindowId,
		Tytul:            z.Title,
		Rodzaj:           z.Kind,
		Adres:            z.Url,
		Pochodzenie:      z.Origin,
		Wiarygodnosc:     wiarygodnoscZadania(z.Credibility),
		PlikBibliotekiID: z.LibraryFileId,
	}
	zapisane, err := a.repozytorium.ZapiszZrodlo(ctx, zrodlo)
	if err != nil {
		return shared.ResearchSourceAddResponse{}, bladBadan(err)
	}

	// Katalogowanie i cechy lektury idą tym samym żądaniem, więc zapisują się
	// razem ze źródłem. Rozbicie ich na drugą komendę zostawiłoby okno w stanie,
	// w którym Operator wpisał etykiety, a katalog ich nie ma — i nie wiedziałby,
	// czy pominął krok, czy rdzeń je zgubił.
	if z.Tags != nil || z.CollectionIds != nil || z.QuestionIds != nil {
		if err := a.repozytorium.UstawKatalogZrodla(ctx, zapisane.Kod, dane.KatalogZrodlaBadania{
			Etykiety: z.Tags, Kolekcje: z.CollectionIds, Pytania: z.QuestionIds,
		}); err != nil {
			return shared.ResearchSourceAddResponse{}, bladBadan(err)
		}
	}
	lektura := dane.LekturaZrodlaBadania{Identyfikat: z.Identifier, CslJson: z.CslJson}
	if z.ReadingState != nil {
		lektura.StanLektury = string(*z.ReadingState)
	}
	if z.StageIndex != nil {
		etap := int64(*z.StageIndex)
		lektura.EtapIndeks = &etap
	}
	if err := a.repozytorium.UstawLektureZrodla(ctx, zapisane.Kod, lektura); err != nil {
		return shared.ResearchSourceAddResponse{}, bladBadan(err)
	}
	return shared.ResearchSourceAddResponse{Source: zlozZrodloBadania(zapisane)}, nil
}

// wiarygodnoscZadania przekłada pole opcjonalne żądania. Brak wskazania
// zostaje pustą wartością, którą warstwa danych zamienia na `unverified` —
// adapter nie wylicza oceny sam (nagłówek pliku).
func wiarygodnoscZadania(wskazana *shared.ResearchCredibility) shared.ResearchCredibility {
	if wskazana != nil {
		return *wskazana
	}
	return ""
}

// zlozZrodloBadania przekłada wiersz repozytorium na byt kontraktu.
func zlozZrodloBadania(z dane.ZrodloBadania) shared.ResearchSource {
	return shared.ResearchSource{
		Id: z.Kod, WindowId: z.Okno, Kind: z.Rodzaj, Url: z.Adres, Title: z.Tytul,
		Origin: z.Pochodzenie, Credibility: z.Wiarygodnosc, LibraryFileId: z.PlikBibliotekiID,
		AcquiredAt: chwilaBazy(z.PozyskanoO),
	}
}

// DodajUstalenie obsługuje `research.finding.add`. Wiąże ustalenie ze
// źródłami wskazanymi w żądaniu, ale odpowiedź niesie komplet, który
// warstwa danych naprawdę powiązała — kod źródła nieistniejącego jest
// pomijany po stronie `dane/badania_ustalenia.go`, nie wywraca zapisu, a
// odpowiedź mówi prawdę o skutku, nie powtarza żądania.
func (a *adapterBadan) DodajUstalenie(ctx context.Context,
	z shared.ResearchFindingAddRequest) (shared.ResearchFindingAddResponse, error) {

	if z.WindowId == "" {
		return shared.ResearchFindingAddResponse{}, bladWskazaniaBadan("żądanie bez okna badania")
	}
	if z.Content == "" {
		return shared.ResearchFindingAddResponse{}, bladWskazaniaBadan("żądanie bez treści ustalenia")
	}

	tresc := z.Content
	ustalenie := dane.UstalenieBadania{
		Kod:   identyfikatorUstaleniaBadania(z.FindingId),
		Okno:  z.WindowId,
		Tresc: &tresc,
		Stan:  stanUstaleniaZadania(z.Status),
	}
	zapisane, err := a.repozytorium.ZapiszUstalenie(ctx, ustalenie, z.SourceIds)
	if err != nil {
		return shared.ResearchFindingAddResponse{}, bladBadan(err)
	}
	// Klasyfikacja, waga, kotwica i wymóg potwierdzenia przychodzą tym samym
	// żądaniem — zapisują się razem z ustaleniem, a nie drugą komendą.
	szczegoly := dane.SzczegolyUstaleniaBadania{AdnotacjaKod: z.AnnotationId}
	if z.Kind != nil {
		szczegoly.Rodzaj = string(*z.Kind)
	}
	if z.Weight != nil {
		szczegoly.Waga = string(*z.Weight)
	}
	if z.NeedsConfirmation != nil {
		szczegoly.WymagaPotwierdzeni = *z.NeedsConfirmation
	}
	if z.Anchor != nil {
		szczegoly.Kotwica = kotwicaDoBazyBadania(*z.Anchor)
	}
	if err := a.repozytorium.UstawSzczegolyUstalenia(ctx, zapisane.Kod, szczegoly); err != nil {
		return shared.ResearchFindingAddResponse{}, bladBadan(err)
	}
	if len(z.CodeIds) > 0 {
		if err := a.repozytorium.UstawKodyUstalenia(ctx, zapisane.Kod, z.CodeIds); err != nil {
			return shared.ResearchFindingAddResponse{}, bladBadan(err)
		}
	}
	// Ślad prowenancji zaczyna się przy powstaniu ustalenia, a nie przy jego
	// pierwszej zmianie — inaczej historia zaczynałaby się od drugiego zdarzenia.
	wpis := dane.WpisProwenancjiBadania{
		UstalenieKod: zapisane.Kod, Aktor: string(shared.ActorKindOperator),
		Czynnosc: "zapis ustalenia",
	}
	if len(z.SourceIds) > 0 {
		kod := z.SourceIds[0]
		wpis.ZrodloKod = &kod
	}
	_ = a.repozytorium.ZapiszProwenancje(ctx, wpis)

	zrodla, err := a.repozytorium.ZrodlaUstalenia(ctx, zapisane.ID)
	if err != nil {
		return shared.ResearchFindingAddResponse{}, bladBadan(err)
	}
	return shared.ResearchFindingAddResponse{Finding: zlozUstalenieBadania(zapisane, zrodla)}, nil
}

// identyfikatorUstaleniaBadania zwraca identyfikator ustalenia zmienianego
// albo nadaje nowy — brak wskazania w żądaniu znaczy ustalenie nowe.
func identyfikatorUstaleniaBadania(wskazany *string) string {
	if wskazany != nil && *wskazany != "" {
		return *wskazany
	}
	return nowyIdentyfikator(przedrostekUstaleniaBadania)
}

// stanUstaleniaZadania przekłada pole opcjonalne żądania. Brak wskazania
// zostaje pustą wartością, którą warstwa danych zamienia na `open`.
func stanUstaleniaZadania(wskazany *shared.ResearchFindingStatus) shared.ResearchFindingStatus {
	if wskazany != nil {
		return *wskazany
	}
	return ""
}

// zlozUstalenieBadania przekłada wiersz repozytorium i jego źródła powiązane
// na byt kontraktu.
func zlozUstalenieBadania(u dane.UstalenieBadania, zrodla []dane.ZrodloBadania) shared.ResearchFinding {
	tresc := ""
	if u.Tresc != nil {
		tresc = *u.Tresc
	}
	kodyZrodel := make([]string, 0, len(zrodla))
	for _, z := range zrodla {
		kodyZrodel = append(kodyZrodel, z.Kod)
	}
	return shared.ResearchFinding{
		Id: u.Kod, WindowId: u.Okno, Content: tresc, SourceIds: kodyZrodel, Status: u.Stan,
		CreatedAt: chwilaBazy(u.Utworzono), UpdatedAt: chwilaBazy(u.Zaktualizowano),
	}
}

// bladBadan znakuje usterkę kodem kontraktu, żeby okno modułu pokazało
// powód, a nie samo „nie udało się". Błąd, któremu kod już nadano, przechodzi
// bez zmiany; dopiero usterka bez kodu staje się usterką wewnętrzną rdzenia.
func bladBadan(err error) error {
	if err == nil {
		return nil
	}
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError, err))
}

// bladWskazaniaBadan nazywa brak danych w żądaniu — błąd Operatora, nie rdzenia.
func bladWskazaniaBadan(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed, "moduł Research: "+powod))
}
