// Moduł Translate — typ adaptera, konstruktor i przedrostki identyfikatorów
// bytów modułu (okno, panel, termin, ślady wymiany, synteza, eksport). Tu leżą
// też trzy komendy przekładu: ustalenie tekstu źródłowego, segmentacja i próba
// rozpoznania języka.
//
// Pozostałe pliki modułu dopisują metody do tego samego `*adapterTlumaczenia`;
// port `Tlumaczenie` i `zarejestrujTlumaczenie` deklaruje
// `adapter_modul_tlumaczenie_uchwyty.go`.
//
// Granice modułu:
//  1. `source.detect` rozpoznaje język modelem. Rdzeń nie ma własnego silnika
//     rozpoznania, ma za to most do rejestru kanałów: pyta domyślny czynny
//     kanał, jaki to język. Bez wpiętego rejestru albo bez czynnego kanału
//     odmawia wprost, zamiast zgadywać po znakach diakrytycznych
//     (`RozpoznajJezyk` niżej).
//  2. `source.set` nie wytwarza treści tłumaczenia — zapisuje sam tekst
//     źródłowy i liczbę segmentów. Treść panelu powstaje przy `target.add`,
//     gdzie znany jest język docelowy (`DodajPanel`,
//     `adapter_modul_tlumaczenie_panele.go`). Bez czynnego kanału `target.add`
//     odmawia zamiast zakładać panel pusty.
//  3. Segmentacja (`source.segment`) idzie podziałem własnym
//     (`podzielNaZdania`), bez biblioteki zewnętrznej.
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
	// woła model. Wpięty wzorem Roundtable (`ZKanalami`). Bez niego operacja
	// modelowa odmawia wprost zamiast oddać pusty wynik udający tłumaczenie.
	kanaly *models.Rejestr
	// wyjscie jest szyną zdarzeń rdzenia — mostem, którym moduł rozgłasza
	// `translate.translation.changed` po zmianie treści panelu. Wpinany przez
	// `ZWyjsciem` (`adapter_modul_tlumaczenie_model.go`) w chwili rejestracji
	// komend. Niepodłączony nie jest błędem.
	wyjscie *emiter
	// Cztery pola niżej składają silnik syntezy mowy (`speech.synthesize`),
	// wpinany przez `ZSynteza` (`adapter_modul_tlumaczenie_mowa_silnik.go`). Ta
	// sama czwórka co przy rozpoznawaniu mowy: port startu procesu, dwa źródła
	// izolacji i katalog danych rdzenia jako miejsce na nagrania. Niewpięte nie
	// psują pozostałych komend modułu — sama synteza odmawia wtedy, nazywając
	// brak.
	// biblioteka i magazynWytworow składają drogę wytworu na zewnątrz modułu
	// (`translate.artifact.publish`): magazyn odkłada bajty pod sumą kontrolną,
	// repozytorium biblioteki zakłada wiersz pliku, po którym reszta platformy
	// wytwór widzi. Wpina je `ZWytworami`; ich brak nie psuje pozostałych komend
	// — samo wydanie wytworu odmawia wtedy, nazywając brak.
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
	"moduł Translate: rdzeń nie ma wpiętego magazynu wytworów ani repozytorium biblioteki — " +
		"wytwór nie miałby gdzie leżeć; naprawa: podpiąć ZWytworami przy składaniu rdzenia")

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

// nowyAdapterTlumaczenia wiąże adapter z repozytorium modułu.
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

// zapytajModel woła wskazany kanał modelu i zbiera całą odpowiedź tekstową.
// Jest wspólnym mostem tego modułu do modelu: przekład, kontrola jakości czy
// tłumaczenie zwrotne wołane modelem idą tędy, a nie każde własną kopią pętli
// strumienia. Kanał musi wskazać wywołujący — rdzeń nie zgaduje, na
// którym kanale okno pracuje; gdy rejestr nie jest wpięty albo kanał pusty,
// most odmawia wprost, nie oddaje pustego napisu udającego przekład.
//
// Kanał wskazuje wywołujący, nie ten most: kontrakt niesie pole `channelId`
// (`target.add`, `backtranslation.run`), a przekład wskazania na kanał składa
// `kanalZadania` (`adapter_modul_tlumaczenie_model.go`). Tutaj zostaje sama
// zasada: pusty `kanal` to odmowa wprost, bo most nie dobiera kanału za
// wywołującego.
func (a *adapterTlumaczenia) zapytajModel(ctx context.Context, okno, kanal, tresc string) (string, error) {
	if a.kanaly == nil {
		return "", bladBrakuKanalowTlumaczenia()
	}
	if strings.TrimSpace(kanal) == "" {
		return "", bladWskazaniaTlumaczenia("żądanie bez wskazania kanału modelu — rdzeń nie zgaduje kanału tłumaczenia")
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

// UstawZrodlo obsługuje `translate.source.set`. Zapisuje tekst źródłowy okna
// i — gdy okno zakłada się po raz pierwszy albo Operator zażądał ponownego
// podziału — liczbę segmentów policzoną przez `podzielNaZdania`. Nie wywołuje
// żadnego rozpoznania języka: pole `SourceLanguage` odpowiedzi niesie to, co
// podał Operator, albo puste, gdy nie podał.
//
// Panele wychodzą stąd bez treści — ich treść powstaje w `DodajPanel`
// (`adapter_modul_tlumaczenie_panele.go`). Adapter nie zna jeszcze paneli okna
// w chwili pierwszego zapisu źródła, więc odpowiedź niesie wtedy pustą listę.
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

	// Segmentacja liczy się od nowa, gdy Operator o to poprosił, albo gdy
	// okno jeszcze nie istnieje (liczba wtedy jest niepoznana skądinąd).
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

// liczbaSegmentowOdpowiedzi przekłada wskaźnik wiersza na wskaźnik kontraktu.
func liczbaSegmentowOdpowiedzi(liczba *int64) *int {
	if liczba == nil {
		return nil
	}
	n := int(*liczba)
	return &n
}

// PodzielNaSegmenty obsługuje `translate.source.segment`. Dzieli wskazany
// tekst — albo, gdy żądanie go nie niesie, aktualny tekst źródłowy okna —
// na zdania metodą `podzielNaZdania`. Komenda nie zapisuje nic do bazy:
// kontrakt oddaje z niej samą listę napisów bez identyfikatorów, a segmenty nie
// mają własnej tabeli, więc powtórne wywołanie na tym samym tekście daje ten
// sam wynik bez efektu ubocznego.
//
// Żądanie bez wskazania tekstu i bez okna, którego tekst źródłowy dałoby się
// wziąć, jest błędem wskazania — nie ma z czego dzielić.
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

// podzielNaZdania dzieli tekst na zdania po znakach końca zdania (kropka,
// wykrzyknik, pytajnik), które są kolejno spacją albo końcem tekstu — jedyny
// sygnał mechaniczny, jaki mamy bez wiedzy językowej o skrótach, inicjałach
// czy cudzysłowach zagnieżdżających kropkę. Ograniczenie: skróty w rodzaju
// „np." albo „ul." rozłamią zdanie tam, gdzie językoznawczo zdanie się nie
// kończy — podział semantyczny wymagałby słownika skrótów albo modelu, których
// rdzeń nie ma. Puste odcinki (wielokrotne białe znaki) są pomijane.
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

// RozpoznajJezyk obsługuje `translate.source.detect`. Rozpoznaje język modelem:
// pyta domyślny czynny kanał modelu, jaki to język (`rozpoznajJezykModelem`
// w `adapter_modul_tlumaczenie_model.go`). Gdy rejestr kanałów nie jest wpięty
// albo nie ma czynnego kanału, komenda odmawia wprost
// (`bladBrakuKanalowTlumaczenia`) zamiast zgadywać po znakach diakrytycznych.
//
// Żądanie niesie sam tekst (kontrakt nie ma pola okna), więc puste `text` jest
// błędem wskazania — nie ma z czego rozpoznać języka. `Confidence` zostaje
// puste: model oddaje nazwę języka, nie miarę pewności, a rdzeń nie dorabia
// liczby, której nikt uczciwie nie wypełni.
func (a *adapterTlumaczenia) RozpoznajJezyk(ctx context.Context,
	z shared.TranslateSourceDetectRequest) (shared.TranslateSourceDetectResponse, error) {

	tekst := ""
	if z.Text != nil {
		tekst = strings.TrimSpace(*z.Text)
	}
	if tekst == "" {
		return shared.TranslateSourceDetectResponse{}, bladWskazaniaTlumaczenia(
			"żądanie translate.source.detect bez tekstu — kontrakt nie niesie okna, więc rdzeń nie ma czego rozpoznać")
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
		// Kod zewnętrzny okna, doczytany złączeniem w warstwie danych; numer
		// wiersza nie zaadresowałby po stronie klienta żadnego okna.
		WindowId:  p.OknoKod,
		Language:  p.Jezyk,
		Text:      p.Tresc,
		Status:    shared.TranslationStatus(p.Stan),
		Tone:      p.Ton,
		UpdatedAt: p.Zaktualizowano,
		// Migawka obiegu zatwierdzeń (`translate.approval.set`). Panel, którego
		// nikt nie przeprowadził przez obieg, wychodzi z pustym etapem — nie
		// z etapem `translation` udającym, że praca ruszyła.
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
// powód, a nie samo „nie udało się". Błąd, któremu kod już nadano, przechodzi
// bez zmiany; dopiero usterka bez kodu staje się usterką wewnętrzną rdzenia.
// Wspólny dla całego modułu.
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
