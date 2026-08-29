// Plik deklaruje typ adaptera modułu Translate, konstruktor, przedrostki
// identyfikatorów bytów modułu oraz trzy komendy przekładu: zapis tekstu
// źródłowego, segmentacja i rozpoznanie języka.
package core

import (
	"context"
	"errors"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// Przedrostki identyfikatorów bytów całego modułu Translate. Zadeklarowane tu
// w całości — łącznie z tymi, których ten plik nie używa — żeby pozostałe pliki
// modułu nie deklarowały ich po raz drugi.
const (
	przedrostekOknaTlumaczenia   = "okt-"
	przedrostekPaneluTlumaczenia = "pan-"
	przedrostekTerminuSlownika   = "trm-"
	przedrostekImportuSlownika   = "imp-"
	przedrostekEksportuSlownika  = "eks-"
	przedrostekPamieciTlumaczen  = "pam-"
	przedrostekSyntezyMowy       = "syn-"
	przedrostekEksportuPanelu    = "epa-"
)

// adapterTlumaczenia wypełnia część portu Tlumaczenie. Zależność jest jedna:
// wspólne repozytorium modułu, rozłożone po stronie danych na kilka plików
// wedle odpowiedzialności, ale niosące jeden typ.
type adapterTlumaczenia struct {
	repozytorium dane.RepozytoriumTlumaczen
	// kanaly jest rejestrem kanałów modelu rdzenia — jedyną drogą, którą moduł
	// woła model.
	kanaly *models.Rejestr
	// wyjscie jest szyną zdarzeń rdzenia; moduł rozgłasza nią zmianę treści
	// panelu, gdy jest podłączona.
	wyjscie *emiter
	// Pola niżej składają wydanie wytworu poza moduł i syntezę mowy lokalnym
	// silnikiem.
	biblioteka       dane.RepozytoriumBiblioteki
	magazynWytworow  *magazynTresciBiblioteki
	uruchamiacz      session.Uruchamiacz
	rozstrzygaczMowy *konfig.Rozstrzygacz
	katalogIzolacji  *KatalogRoboczy
	katalogDanych    string
}

// errPustyPrzeklad i errPusteRozpoznanie znakują pustą odpowiedź modelu na
// operacji, która miała dać treść. Pusty wynik nie jest tłumaczeniem ani
// rozpoznaniem — most modelu zamienia go w odmowę, nie w atrapę.
var (
	errPustyPrzeklad    = errors.New("model oddał pusty przekład — panel nie dostał treści")
	errPusteRozpoznanie = errors.New("model oddał pustą odpowiedź — język nierozpoznany")
	// Ten sam powód co wyżej, ale o kontroli wierności: pusty wynik tłumaczenia
	// zwrotnego nie znaczy „przekład wierny", tylko „nie ma czym sprawdzić".
	errPusteTlumaczenieZwrotne = errors.New("model oddał puste tłumaczenie zwrotne — kontrola wierności nie ma wyniku")
)

// errBrakMagazynuWytworow znakuje odmowę wydania wytworu w rdzeniu złożonym
// bez magazynu treści i repozytorium biblioteki. Wytwór bez miejsca, w którym
// leży, i bez wiersza, który o nim wie, nie jest wytworem.
var errBrakMagazynuWytworow = errors.New(
	"moduł Translate: serwer nie ma wpiętego magazynu wytworów ani repozytorium biblioteki — " +
		"wytwór nie miałby gdzie leżeć; naprawa: podpiąć ZWytworami przy składaniu serwera")

// ZWytworami wpina drogę wydania wytworu: repozytorium biblioteki i katalog
// danych rdzenia, pod którym stoi magazyn treści.
func (a *adapterTlumaczenia) ZWytworami(biblioteka dane.RepozytoriumBiblioteki,
	katalogDanych string) *adapterTlumaczenia {
	a.biblioteka = biblioteka
	if strings.TrimSpace(katalogDanych) != "" {
		a.magazynWytworow = nowyMagazynTresciBiblioteki(katalogDanych)
	}
	return a
}

// nowyAdapterTlumaczenia wiąże adapter z repozytorium modułu, pozostawiając
// pozostałe zależności do wpięcia osobno.
func nowyAdapterTlumaczenia(repozytorium dane.RepozytoriumTlumaczen) *adapterTlumaczenia {
	return &adapterTlumaczenia{repozytorium: repozytorium}
}

// ZKanalami wpina rejestr kanałów modelu — most, którym tłumaczenie woła model
// (wzór: `adapter_modul_roundtable.go`). Zwraca adapter, żeby montaż wiązał
// zależność w łańcuchu, tak jak robią to pozostałe moduły korzystające z modelu.
func (a *adapterTlumaczenia) ZKanalami(kanaly *models.Rejestr) *adapterTlumaczenia {
	a.kanaly = kanaly
	return a
}

// zapytajModel woła wskazany kanał modelu i zbiera całą odpowiedź tekstową —
// wspólny most modułu do modelu dla operacji słownych. Pusty kanał albo brak
// rejestru kończy się odmową wprost.
func (a *adapterTlumaczenia) zapytajModel(ctx context.Context, okno, kanal, tresc string) (string, error) {
	if a.kanaly == nil {
		return "", bladBrakuKanalowTlumaczenia()
	}
	if strings.TrimSpace(kanal) == "" {
		return "", bladWskazaniaTlumaczenia("żądanie bez wskazania kanału modelu — serwer nie zgaduje kanału tłumaczenia")
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
		return "", bladTlumaczenia(err)
	}
	return zebrane.String(), nil
}

// bladBrakuKanalowTlumaczenia znakuje odmowę, gdy moduł ma wykonać operację
// modelu, a rejestr kanałów nie został wpięty — kanał niedostępny, nie usterka
// wewnętrzna (wzór `bladBrakuKanalow` z Roundtable).
func bladBrakuKanalowTlumaczenia() error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
		"moduł Translate: rejestr kanałów modelu nie jest wpięty — tłumaczenie nie ma czym wołać modelu"))
}

// UstawZrodlo obsługuje `translate.source.set`: zapisuje tekst źródłowy okna
// i, gdy trzeba, liczbę segmentów, nie wywołując rozpoznania języka.
func (a *adapterTlumaczenia) UstawZrodlo(ctx context.Context,
	z shared.TranslateSourceSetRequest) (shared.TranslateSourceSetResponse, error) {

	if z.WindowId == "" {
		return shared.TranslateSourceSetResponse{}, bladWskazaniaTlumaczenia("żądanie bez okna tłumaczenia")
	}
	if z.Text == "" {
		return shared.TranslateSourceSetResponse{}, bladWskazaniaTlumaczenia("żądanie bez tekstu źródłowego")
	}

	tekst := z.Text
	okno := dane.OknoTlumaczenia{
		Kod:           z.WindowId,
		TekstZrodlowy: &tekst,
		JezykZrodlowy: z.SourceLanguage,
	}

	// Segmentacja liczy się od nowa, gdy operator o to poprosił albo gdy
	// okno jeszcze nie istnieje.
	if z.Resegment == nil || *z.Resegment {
		liczba := int64(len(podzielNaZdania(z.Text)))
		okno.LiczbaSegmentow = &liczba
	}

	zapisane, err := a.repozytorium.ZapiszOkno(ctx, okno)
	if err != nil {
		return shared.TranslateSourceSetResponse{}, bladTlumaczenia(err)
	}

	panele, err := a.repozytorium.Panele(ctx, zapisane.ID)
	if err != nil {
		return shared.TranslateSourceSetResponse{}, bladTlumaczenia(err)
	}

	return shared.TranslateSourceSetResponse{
		SourceLanguage: jezykZrodlowyZadania(zapisane.JezykZrodlowy),
		SegmentCount:   liczbaSegmentowOdpowiedzi(zapisane.LiczbaSegmentow),
		Panels:         zlozPaneleTlumaczenia(panele),
	}, nil
}

// jezykZrodlowyZadania przekłada pole opcjonalne wiersza. Brak rozpoznania
// zostaje pustym napisem — kontrakt każe oddać string, nie wskaźnik, ale
// pusty string nie udaje rozpoznanego języka.
func jezykZrodlowyZadania(jezyk *string) string {
	if jezyk != nil {
		return *jezyk
	}
	return ""
}

// liczbaSegmentowOdpowiedzi przekłada wskaźnik wiersza repozytorium na
// wskaźnik pola kontraktu, zachowując brak wartości.
func liczbaSegmentowOdpowiedzi(liczba *int64) *int {
	if liczba == nil {
		return nil
	}
	n := int(*liczba)
	return &n
}

// PodzielNaSegmenty obsługuje `translate.source.segment`: dzieli wskazany
// tekst, albo tekst źródłowy okna, na zdania metodą `podzielNaZdania`, bez
// zapisu do bazy.
func (a *adapterTlumaczenia) PodzielNaSegmenty(ctx context.Context,
	z shared.TranslateSourceSegmentRequest) (shared.TranslateSourceSegmentResponse, error) {

	tekst := ""
	if z.Text != nil {
		tekst = *z.Text
	}
	if tekst == "" {
		return shared.TranslateSourceSegmentResponse{}, bladWskazaniaTlumaczenia(
			"żądanie bez tekstu do podziału — translate.source.segment nie ma dostępu do okna bez tekstu wskazanego wprost")
	}

	return shared.TranslateSourceSegmentResponse{Segments: podzielNaZdania(tekst)}, nil
}

// podzielNaZdania dzieli tekst na zdania po znakach końca zdania, kropce,
// wykrzykniku albo pytajniku, po których stoi biały znak albo koniec tekstu.
func podzielNaZdania(tekst string) []string {
	var zdania []string
	poczatek := 0
	for i, r := range tekst {
		if r != '.' && r != '!' && r != '?' {
			continue
		}
		koniecZnaku := i + len(string(r))
		czyGraniczne := koniecZnaku >= len(tekst) || tekst[koniecZnaku] == ' ' ||
			tekst[koniecZnaku] == '\n' || tekst[koniecZnaku] == '\t'
		if !czyGraniczne {
			continue
		}
		kandydat := strings.TrimSpace(tekst[poczatek:koniecZnaku])
		if kandydat != "" {
			zdania = append(zdania, kandydat)
		}
		poczatek = koniecZnaku
	}
	ogon := strings.TrimSpace(tekst[poczatek:])
	if ogon != "" {
		zdania = append(zdania, ogon)
	}
	return zdania
}

// RozpoznajJezyk obsługuje `translate.source.detect`: pyta domyślny czynny
// kanał modelu, jaki to język, i odmawia wprost bez czynnego kanału.
func (a *adapterTlumaczenia) RozpoznajJezyk(ctx context.Context,
	z shared.TranslateSourceDetectRequest) (shared.TranslateSourceDetectResponse, error) {

	tekst := ""
	if z.Text != nil {
		tekst = strings.TrimSpace(*z.Text)
	}
	if tekst == "" {
		return shared.TranslateSourceDetectResponse{}, bladWskazaniaTlumaczenia(
			"żądanie translate.source.detect bez tekstu — kontrakt nie niesie okna, więc serwer nie ma czego rozpoznać")
	}

	jezyk, err := a.rozpoznajJezykModelem(ctx, tekst)
	if err != nil {
		return shared.TranslateSourceDetectResponse{}, err
	}
	return shared.TranslateSourceDetectResponse{Language: jezyk}, nil
}

// zlozPaneleTlumaczenia przekłada wiersze repozytorium na byty kontraktu.
// Niezgodności panelu (`Issues`) zostają tu puste — wypełnia je odczyt obszaru
// jakości, tak jak ustalenia badania nie są polem `ZrodloBadania`
// (wzór z `adapter_modul_badania.go`).
func zlozPaneleTlumaczenia(panele []dane.PanelTlumaczenia) []shared.TranslationPanel {
	wynik := make([]shared.TranslationPanel, 0, len(panele))
	for _, p := range panele {
		wynik = append(wynik, zlozPanelTlumaczenia(p))
	}
	return wynik
}

// zlozPanelTlumaczenia przekłada jeden wiersz panelu na byt kontraktu. Wspólna
// dla całego modułu: każdy adapter, który odczyta panel z repozytorium, składa
// go tą samą funkcją.
func zlozPanelTlumaczenia(p dane.PanelTlumaczenia) shared.TranslationPanel {
	return shared.TranslationPanel{
		Id: p.Kod,
		// Kod zewnętrzny okna, doczytany złączeniem w warstwie danych, nie
		// numer wiersza wewnętrznego.
		WindowId:  p.OknoKod,
		Language:  p.Jezyk,
		Text:      p.Tresc,
		Status:    shared.TranslationStatus(p.Stan),
		Tone:      p.Ton,
		UpdatedAt: p.Zaktualizowano,
		// Migawka obiegu zatwierdzeń; panel nieprzeprowadzony przez obieg
		// wychodzi z pustym etapem.
		ApprovalStage: etapZatwierdzeniaPanelu(p.EtapZatwierdzenia),
		ApprovedBy:    p.Zatwierdzil,
		ApprovedAt:    p.Zatwierdzono,
	}
}

// etapZatwierdzeniaPanelu przekłada kolumnę migawki na pole kontraktu. Pusty
// napis nie jest etapem, więc zostaje wskaźnikiem pustym.
func etapZatwierdzeniaPanelu(etap *string) *shared.ApprovalStage {
	if etap == nil || strings.TrimSpace(*etap) == "" {
		return nil
	}
	wartosc := shared.ApprovalStage(*etap)
	return &wartosc
}

// bladTlumaczenia znakuje usterkę kodem kontraktu, żeby okno modułu pokazało
// powód. Błąd, któremu kod już nadano, przechodzi bez zmiany.
func bladTlumaczenia(err error) error {
	if err == nil {
		return nil
	}
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError, err))
}

// bladWskazaniaTlumaczenia nazywa brak danych w żądaniu — błąd żądania, nie
// rdzenia. Wspólny dla całego modułu.
func bladWskazaniaTlumaczenia(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed, "moduł Translate: "+powod))
}

// bladNieznanegoOkna odróżnia brak okna tłumaczenia od usterki wewnętrznej —
// ten sam kształt co `bladNieznanegoPanelu`, ale o innym bycie. Dwa byty, dwa
// komunikaty: Operator ma wiedzieć, czego nie ma.
func bladNieznanegoOkna(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Translate: nie ma okna tłumaczenia o identyfikatorze "+kod))
	}
	return bladTlumaczenia(err)
}

// bladNieznanegoPanelu odróżnia „panelu nie ma" (odmowa wprost, kod niezgodny
// wskazania) od usterki wewnętrznej odczytu — wzór `bladNieznanegoPrzebiegu`
// (`adapter_modul_automations_przebiegi.go`).
func bladNieznanegoPanelu(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Translate: nie ma panelu tłumaczenia o identyfikatorze "+kod))
	}
	return bladTlumaczenia(err)
}
