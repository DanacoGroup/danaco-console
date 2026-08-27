// Odpowiedzialność pliku: podstawa odcinka kontroli pracy modułu Studio —
// kontrakt warstwy danych (migracje 363-366 i 370, kolumny wersji z 367),
// rzutowanie repozytorium, rozstrzygnięcie wykonawcy czynności oraz odmowy
// wspólne dla odcinka.
package core

import (
	"context"
	"strconv"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Przedrostki identyfikatorów bytów odcinka. Byt nadany przez rdzeń wychodzi
// kontraktem pod własnym identyfikatorem, nie pod kluczem wiersza.
const (
	przedrostekBlokadyStudia    = "studio-blk-"
	przedrostekCzynnosciStudia  = "studio-czn-"
	przedrostekKopiiStudia      = "studio-kop-"
	przedrostekZnakowaniaStudia = "studio-znk-"
	przedrostekZajeciaStudia    = "studio-zaj-"
	przedrostekSpieciaStudia    = "studio-spc-"
)

// Klucze konfiguracji nastaw pracy wielu wykonawców. Wartości idą zasięgami
// rodziny `config.*` — tą samą tabelą `ustawienie`, którą jedzie cała
// platforma — a nie osobnym magazynem nastaw Studia.
const (
	kluczPetliWykonawczej  = "studio.agents.executionLoopEnabled"
	kluczPracyWieluAgentow = "studio.agents.multiAgentEnabled"
	kluczNajwiecejAgentow  = "studio.agents.maxConcurrentAgents"
	kluczNastawySpiecia    = "studio.agents.conflictPolicy"
	kluczObiegowPetli      = "studio.agents.loopMaxIterations"
	kluczProguBezPostepu   = "studio.agents.loopNoProgressThreshold"
	kluczWymogZajecia      = "studio.agents.requireFragmentClaim"
)

// KontrolaPracyStudia jest kontraktem warstwy danych odcinka kontroli pracy:
// blokady fragmentów, dziennik czynności, kopie zapasowe, nastawy autozapisu,
// znakowanie, zajęcia wykonawców, spięcia i wersje w szeregach.
type KontrolaPracyStudia interface {
	// Blokady fragmentów (migracja 363).
	ZapiszBlokadeFragmentu(ctx context.Context, dokumentID int64,
		blokada dane.BlokadaFragmentuStudia) (dane.BlokadaFragmentuStudia, error)
	BlokadaFragmentu(ctx context.Context, kod string) (dane.BlokadaFragmentuStudia, error)
	BlokadyFragmentow(ctx context.Context, dokumentID int64) ([]dane.BlokadaFragmentuStudia, error)
	UsunBlokadeFragmentu(ctx context.Context, kod string) (bool, error)
	PrzesunBlokadyFragmentow(ctx context.Context, dokumentID, odPozycji, przesuniecie int64) error
	ZapiszBlokadeSzablonu(ctx context.Context, blokada dane.BlokadaSzablonuStudia) error
	BlokadySzablonu(ctx context.Context, szablonKod string) ([]dane.BlokadaSzablonuStudia, error)

	// Odwracalny dziennik czynności (364).
	ZapiszCzynnoscDokumentu(ctx context.Context, dokumentID int64,
		czynnosc dane.CzynnoscDokumentuStudia) (dane.CzynnoscDokumentuStudia, error)
	CzynnoscDokumentu(ctx context.Context, kod string) (dane.CzynnoscDokumentuStudia, error)
	CzynnosciDokumentu(ctx context.Context, dokumentID int64) ([]dane.CzynnoscDokumentuStudia, error)
	ZapiszZaleznoscCzynnosci(ctx context.Context, czynnoscKod, podstawaKod string, powod *string) error
	PrzestawStanCzynnosci(ctx context.Context, kod, stanOczekiwany, stanNowy string) (bool, error)

	// Kopie zapasowe i nastawy autozapisu (366).
	ZapiszKopieZapasowa(ctx context.Context, dokumentID int64,
		kopia dane.KopiaZapasowaStudia) (dane.KopiaZapasowaStudia, error)
	KopiaZapasowa(ctx context.Context, kod string) (dane.KopiaZapasowaStudia, error)
	KopieZapasowe(ctx context.Context, dokumentID int64) ([]dane.KopiaZapasowaStudia, error)
	KopieNiezapisane(ctx context.Context, okno string) ([]dane.KopiaZapasowaStudia, error)
	PrzemiecKopieZapasowe(ctx context.Context, dokumentID, ileZachowac, wygasanieGodzin int64) (int64, error)
	NastawaPracy(ctx context.Context, okno string, dokumentID *int64) (dane.NastawaPracyStudia, error)
	ZapiszNastaweAutozapisu(ctx context.Context, nastawa dane.NastawaPracyStudia) error
	ZapiszSkutekAutozapisu(ctx context.Context, nastawaID int64, chwila *string,
		nieudany bool, powod *string) error

	// Znakowanie i rodzaje znaczników (365).
	ZapiszZnakowanie(ctx context.Context, dokumentID int64,
		znakowanie dane.ZnakowanieStudia) (dane.ZnakowanieStudia, error)
	Znakowanie(ctx context.Context, kod string) (dane.ZnakowanieStudia, error)
	Znakowania(ctx context.Context, dokumentID int64) ([]dane.ZnakowanieStudia, error)
	UsunZnakowanie(ctx context.Context, kod string) (bool, error)
	PrzestawStanZnakowania(ctx context.Context, kod, stan string) (bool, error)
	RodzajeZnacznika(ctx context.Context, dokumentID int64) ([]dane.RodzajZnacznikaStudia, error)
	RodzajZnacznika(ctx context.Context, nazwa string) (dane.RodzajZnacznikaStudia, error)
	ZapiszRodzajZnacznika(ctx context.Context,
		rodzaj dane.RodzajZnacznikaStudia) (dane.RodzajZnacznikaStudia, error)
	UsunRodzajZnacznika(ctx context.Context, nazwa string) (bool, error)

	// Zajęcia fragmentów i spięcia wykonawców (370).
	ZapiszZajecieFragmentu(ctx context.Context, dokumentID int64,
		zajecie dane.ZajecieFragmentuStudia, waznoscSekund int64) (dane.ZajecieFragmentuStudia, error)
	ZajeciaFragmentow(ctx context.Context, dokumentID int64) ([]dane.ZajecieFragmentuStudia, error)
	ZwolnijZajecieFragmentu(ctx context.Context, kod string) (bool, error)
	ZwolnijZajeciaWykonawcy(ctx context.Context, dokumentID int64, agentKod, podagentKod string) (int64, error)
	PrzemiecZajeciaFragmentow(ctx context.Context, dokumentID int64) (int64, error)
	ZapiszSpiecieWykonawcow(ctx context.Context, dokumentID int64,
		spiecie dane.SpiecieWykonawcowStudia) error
	SpieciaWykonawcow(ctx context.Context, dokumentID int64) ([]dane.SpiecieWykonawcowStudia, error)

	// Wersje w szeregu Operatora i w szeregu autozapisu (kolumny z 367).
	ZapiszWersjeSzeregu(ctx context.Context, dokumentID int64,
		wersja dane.WersjaSzereguStudia) (dane.WersjaSzereguStudia, error)
	WersjaSzeregu(ctx context.Context, kod string) (dane.WersjaSzereguStudia, error)
	WersjeSzeregow(ctx context.Context, dokumentID int64) ([]dane.WersjaSzereguStudia, error)
	WersjaZalozycielska(ctx context.Context, dokumentID int64) (dane.WersjaSzereguStudia, error)
	PrzemiecWersjeAutozapisu(ctx context.Context, dokumentID, ileZachowac int64) (int64, error)

	// Nastawy widoku okna pracy — te same wiersze co nastawy autozapisu,
	// inne kolumny (366).
	NastawaWidoku(ctx context.Context, okno string, dokumentID *int64) (dane.NastawaWidokuStudia, error)
	ZapiszNastaweWidoku(ctx context.Context, nastawa dane.NastawaWidokuStudia) error

	// Historia schowka platformy (292) — JEDNA dla Operatora i wykonawców.
	DopiszWpisSchowkaStudia(ctx context.Context, wpis dane.WpisSchowka) (dane.WpisSchowka, bool, error)
	WpisSchowkaStudia(ctx context.Context, kod string) (dane.WpisSchowka, error)
	NajswiezszyWpisSchowkaStudia(ctx context.Context) (dane.WpisSchowka, error)

	// Nastawy wielu wykonawców — zasięgami rodziny `config.*`, tabelą `ustawienie`.
	ZapiszNastaweZasieguStudia(ctx context.Context, ustawienie dane.Ustawienie) error
	NastawaZasieguStudia(ctx context.Context, poziom shared.ConfigScope,
		kluczZasiegu, klucz string) (dane.Ustawienie, bool, error)

	// Zmiany śledzone wraz z tożsamością wykonawcy (kolumny z 369).
	ZmianyWykonawcow(ctx context.Context, dokumentID int64) ([]dane.ZmianaWykonawcyStudia, error)
	StemplujTozsamoscZmiany(ctx context.Context, kod string,
		agentKod, agentNazwa, agentWersja, podagentKod *string,
		postacPrzed, postacPo, czynnoscKod *string) error

	// Pochodzenie fragmentów (368) — wiersz zakłada odcinek wejścia, wykaz
	// czyta ten odcinek.
	ZapiszPochodzenieFragmentu(ctx context.Context,
		pochodzenie dane.PochodzenieFragmentuStudia) (dane.PochodzenieFragmentuStudia, error)
	PochodzenieFragmentow(ctx context.Context,
		dokumentID int64) ([]dane.PochodzenieFragmentuStudia, error)
}

// kontrolaSkladnica zdejmuje kontrakt odcinka z repozytorium Studia. Brak jest
// brakiem montażu rdzenia, nie brakiem funkcji produktu — i tak się nazywa.
func (a *adapterStudia) kontrolaSkladnica() (KontrolaPracyStudia, error) {
	if a == nil || a.repozytorium == nil {
		return nil, kontrolaBladZaplecza("repozytorium Studia nie zostało podane przy montażu rdzenia")
	}
	skladnica, jest := a.repozytorium.(KontrolaPracyStudia)
	if !jest {
		return nil, kontrolaBladZaplecza("repozytorium Studia nie niesie tabel kontroli pracy " +
			"z migracji 363-366 i 370")
	}
	return skladnica, nil
}

// ── Kto woła ────────────────────────────────────────────────────────────────

// kontrolaWykonawca jest rozstrzygnięciem, czyja ręka wykonuje czynność, wraz
// z tożsamością wykonawcy na tyle dokładną, na ile żądanie ją podało.
type kontrolaWykonawca struct {
	// Rodzaj jest grubym rozróżnieniem człowiek-wykonawca; tożsamością wielu
	// agentów jest kod agenta.
	Rodzaj       shared.StudioAuthor
	AgentKod     *string
	AgentNazwa   *string
	AgentWersja  *string
	PodagentKod  *string
	KartaSesjiID *string
	OknoID       *string
	KanalID      *string
}

// czyWykonawca mówi, czy czynność jest czynnością wykonawcy, a nie operatora,
// na podstawie rodzaju zapisanego w rozstrzygnięciu.
func (w kontrolaWykonawca) czyWykonawca() bool {
	return w.Rodzaj == shared.StudioAuthorModel
}

// nazwaWykonawcy składa zdanie, którym odmowa i bilans nazywają rękę. Nazwa
// widoczna stoi przed kodem, bo Operator czyta nazwę, a nie identyfikator;
// wykonawca nienazwany zostaje „modelem" — to prawda węższa, ale prawda.
func (w kontrolaWykonawca) nazwaWykonawcy() string {
	if !w.czyWykonawca() {
		return "Operator"
	}
	if w.AgentNazwa != nil && *w.AgentNazwa != "" {
		if w.AgentKod != nil && *w.AgentKod != "" {
			return *w.AgentNazwa + " (" + *w.AgentKod + ")"
		}
		return *w.AgentNazwa
	}
	if w.AgentKod != nil && *w.AgentKod != "" {
		return "wykonawca " + *w.AgentKod
	}
	if w.PodagentKod != nil && *w.PodagentKod != "" {
		return "podagent " + *w.PodagentKod
	}
	return "model"
}

// jakoAktor przekłada rozstrzygnięcie wykonawcy na strukturę kontraktu, niosącą
// rodzaj, tożsamość agenta oraz identyfikatory sesji, okna i kanału.
func (w kontrolaWykonawca) jakoAktor() shared.StudioActor {
	return shared.StudioActor{
		Kind: w.Rodzaj, AgentId: w.AgentKod, AgentName: w.AgentNazwa,
		AgentVersion: w.AgentWersja, SubagentId: w.PodagentKod,
		SessionId: w.KartaSesjiID, WindowId: w.OknoID, ChannelId: w.KanalID,
	}
}

// kontrolaPodpisZadania niesie pola podpisu wykonawcy wyjęte z ładunku komendy.
// Wszystkie są nieobowiązkowe: żądanie bez podpisu zostaje poprawne.
type kontrolaPodpisZadania struct {
	Author     *shared.StudioAuthor `json:"author,omitempty"`
	AgentId    *string              `json:"agentId,omitempty"`
	AgentName  *string              `json:"agentName,omitempty"`
	SubagentId *string              `json:"subagentId,omitempty"`
}

// kontrolaRozpoznajWykonawce składa rozstrzygnięcie z faktu gniazda i z podpisu
// żądania wedle zasady, że wykonawcą jest ten, kogo wskazuje szerszy z dwóch
// sygnałów; operatorem czynność jest wtedy i tylko wtedy, gdy milczą oba.
func kontrolaRozpoznajWykonawce(ctx context.Context, podpis kontrolaPodpisZadania) kontrolaWykonawca {
	wykonawca := kontrolaWykonawca{Rodzaj: shared.StudioAuthorUzytkownik}

	// Fakt gniazda: serwer narzędzi przedstawia się przy nawiązaniu; okno
	// operatora tak nie wygląda.
	tozsamosc := tozsamoscZKontekstu(ctx)
	if tozsamosc.Narzedzia() {
		wykonawca.Rodzaj = shared.StudioAuthorModel
		if tozsamosc.IdOkna != "" {
			okno := tozsamosc.IdOkna
			wykonawca.OknoID = &okno
		}
	}

	// Podpis żądania podnosi do wykonawcy, nigdy nie obniża do operatora.
	if podpis.Author != nil && *podpis.Author == shared.StudioAuthorModel {
		wykonawca.Rodzaj = shared.StudioAuthorModel
	}
	if kontrolaTekstNiepusty(podpis.AgentId) || kontrolaTekstNiepusty(podpis.AgentName) ||
		kontrolaTekstNiepusty(podpis.SubagentId) {

		wykonawca.Rodzaj = shared.StudioAuthorModel
	}

	// Tożsamość agenta wchodzi wyłącznie z żądania — gniazdo jej nie niesie.
	wykonawca.AgentKod = kontrolaWskaznikNiepusty(podpis.AgentId)
	wykonawca.AgentNazwa = kontrolaWskaznikNiepusty(podpis.AgentName)
	wykonawca.PodagentKod = kontrolaWskaznikNiepusty(podpis.SubagentId)
	return wykonawca
}

// ── Podpis wykonawcy na drodze komendy ──────────────────────────────────────

// kluczWykonawcyStudia jest kluczem, pod którym podpis wykonawcy jedzie
// kontekstem wywołania. Typ własny, nie napis: klucz napisowy dałby się nadpisać
// z zewnątrz pakietu i tożsamość wykonawcy byłaby podrabialna.
type kluczWykonawcyStudia struct{}

// kontrolaZapiszWykonawce wkłada rozpoznanego wykonawcę do kontekstu wywołania,
// skąd czynności odcinka postaci dokumentu czytają go jedną drogą zamiast przez
// parametr każdej sygnatury.
func kontrolaZapiszWykonawce(ctx context.Context, wykonawca kontrolaWykonawca) context.Context {
	return context.WithValue(ctx, kluczWykonawcyStudia{}, wykonawca)
}

// kontrolaWykonawcaZKontekstu czyta podpis wykonawcy z kontekstu wywołania. Brak
// podpisu nie jest usterką: tą drogą idą wywołania spoza rejestru — pętla
// wykonawcza, sprzątanie i sprawdziany; wołający wtedy rozstrzyga sam, kim jest
// wykonawca.
func kontrolaWykonawcaZKontekstu(ctx context.Context) (kontrolaWykonawca, bool) {
	if ctx == nil {
		return kontrolaWykonawca{}, false
	}
	wykonawca, jest := ctx.Value(kluczWykonawcyStudia{}).(kontrolaWykonawca)
	return wykonawca, jest
}

// kontrolaTekstNiepusty mówi, czy wskaźnik niesie napis o treści, odróżniając
// wskaźnik pusty i wskaźnik do napisu pustego od napisu niepustego.
func kontrolaTekstNiepusty(wskazanie *string) bool {
	return wskazanie != nil && *wskazanie != ""
}

// kontrolaWskaznikNiepusty oddaje wskaźnik wyłącznie wtedy, gdy niesie treść. Napis
// pusty ma wyglądać na brak, nie na pustą wartość.
func kontrolaWskaznikNiepusty(wskazanie *string) *string {
	if !kontrolaTekstNiepusty(wskazanie) {
		return nil
	}
	return wskazanie
}

// ── Wspólne przekłady odcinka ───────────────────────────────────────────────

// kontrolaWskaznikLiczby przenosi liczbę bazy typu int64 do pola opcjonalnego
// kontraktu typu int, oddając brak wartości jako wskaźnik pusty.
func kontrolaWskaznikLiczby(wartosc *int64) *int {
	if wartosc == nil {
		return nil
	}
	liczba := int(*wartosc)
	return &liczba
}

// kontrolaWskaznikTekstu przenosi napis do pola opcjonalnego kontraktu, oddając
// napis pusty jako wskaźnik pusty, nie jako pusty napis.
func kontrolaWskaznikTekstu(wartosc string) *string {
	if wartosc == "" {
		return nil
	}
	return &wartosc
}

// kontrolaLiczbaZWskaznika przenosi liczbę opcjonalną z kontraktu typu int do
// kolumny bazy typu int64, zachowując brak wartości jako wskaźnik pusty.
func kontrolaLiczbaZWskaznika(wartosc *int) *int64 {
	if wartosc == nil {
		return nil
	}
	liczba := int64(*wartosc)
	return &liczba
}

// kontrolaZakresWTresci przycina zakres żądania do granic treści liczonych
// w znakach, nie w bajtach, ponieważ kontrakt nazywa początek fragmentu
// w znakach, a cięcie po bajtach rozcinałoby polskie litery dwubajtowe.
func kontrolaZakresWTresci(znaki []rune, od, do int) (int, int, bool) {
	if od < 0 {
		od = 0
	}
	if do > len(znaki) {
		do = len(znaki)
	}
	if od > do {
		return 0, 0, false
	}
	return od, do, true
}

// kontrolaZakresyStykaja mówi, czy dwa zakresy półotwarte mają część wspólną;
// zakres pusty jest punktem wstawienia i styka się z blokadą, w której środku
// leży.
func kontrolaZakresyStykaja(odA, doA, odB, doB int) bool {
	if odA == doA {
		return odA > odB && odA < doB
	}
	if odB == doB {
		return odB > odA && odB < doA
	}
	return odA < doB && odB < doA
}

// ── Odmowy odcinka ──────────────────────────────────────────────────────────

// kontrolaBladZaplecza nazywa brak po stronie montażu rdzenia, na przykład
// repozytorium podane bez tabel kontroli pracy modułu Studio.
func kontrolaBladZaplecza(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"moduł Studio, kontrola pracy: "+powod))
}

// kontrolaBladBlokady jest odmową nazwaną: mówi, który fragment i jaka blokada
// zatrzymały czynność, niosąc nazwę blokady, jej powód i zakres znaków zamiast
// samego zakazu.
func kontrolaBladBlokady(blokada dane.BlokadaFragmentuStudia, wykonawca string) error {
	tresc := "moduł Studio: " + wykonawca + " nie może zmienić fragmentu od znaku " +
		strconv.FormatInt(blokada.ZakresOd, 10) + " do " +
		strconv.FormatInt(blokada.ZakresDo, 10) + " — stoi na nim blokada „" +
		blokada.Nazwa + "”"
	if blokada.Powod != nil && *blokada.Powod != "" {
		tresc += " (" + *blokada.Powod + ")"
	}
	tresc += ". Blokadę zdejmuje wyłącznie Operator; fragment wymagający zmiany " +
		"wskazuje się propozycją na marginesie (studio.markup.add o rodzaju suggestion)."
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodePermissionDenied, tresc))
}

// kontrolaBladZajecia jest odmową nazywającą wykonawcę, który trzyma fragment
// zajęty, wraz z zakresem znaków i chwilą wygaśnięcia zajęcia.
func kontrolaBladZajecia(zajecie dane.ZajecieFragmentuStudia, trzymajacy string) error {
	tresc := "moduł Studio: fragment od znaku " +
		strconv.FormatInt(zajecie.ZakresOd, 10) + " do " +
		strconv.FormatInt(zajecie.ZakresDo, 10) + " trzyma " + trzymajacy
	if zajecie.Wygasa != nil && *zajecie.Wygasa != "" {
		tresc += " do " + *zajecie.Wygasa
	}
	tresc += ". Zajęcie samo wygasa; do tego czasu zmiana tego fragmentu nie wchodzi."
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeConflict, tresc))
}
