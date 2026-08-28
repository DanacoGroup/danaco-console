// Pakiet dane zawiera repozytoria nad pakietem `internal/store`: jedno na
// obszar. Warstwy wyższe sięgają po dane wyłącznie przez interfejsy tego
// pakietu, nigdy przez SQL ani przez wnętrze `store`.
package dane

import (
	"context"
	"database/sql"
	"fmt"

	"danacoconsole/server/internal/store"
)

// Zestaw to komplet repozytoriów jednej otwartej bazy: po jednym polu na
// każdy obszar danych aplikacji.
type Zestaw struct {
	Sesje        RepozytoriumSesji
	Okna         RepozytoriumOkien
	Wiadomosci   RepozytoriumWiadomosci
	Kanaly       RepozytoriumKanalow
	Konta        RepozytoriumKont
	Konfiguracja RepozytoriumKonfiguracji
	// KonfiguracjaOsi to Konfiguracja widziana z osią rozstrzygania: jedna
	// implementacja, dwa widoki.
	KonfiguracjaOsi RepozytoriumKonfiguracjiOsi
	KatalogUstawien RepozytoriumKatalogUstawien
	Kolejki         RepozytoriumKolejek
	Pamiec          RepozytoriumPamieci
	Srodowiska      RepozytoriumSrodowisk
	Moduly          RepozytoriumModulow
	// Macierz czyta wiersz tabeli `srodowisko_modul` w całości; `Moduly` zna ją
	// wyłącznie jako filtr.
	Macierz        RepozytoriumMacierzy
	KartySesji     RepozytoriumKartSesji
	OknaOperacyjne RepozytoriumOkienOperacyjnych
	Akcje          RepozytoriumAkcji
	PunktyDostepu  RepozytoriumPunktowDostepu
	Nadania        RepozytoriumNadan
	Tozsamosc      RepozytoriumTozsamosci
	// Urzadzenia to katalog maszyn (tabela `urzadzenie`), wskazywany przez
	// punkt dostępu `localDirectory`.
	Urzadzenia RepozytoriumUrzadzen
	// CentrumPowiadomien to rejestr zdarzeń centrum powiadomień; to nie jest
	// kolejka doręczeń Mobile.
	CentrumPowiadomien RepozytoriumCentrumPowiadomien
	// Agenci to biblioteka ekspertów modułu Agents; ekspert jest komponentem
	// własnym, żyje obok sesji.
	Agenci RepozytoriumAgentow
	// WarstwyAgenta to tożsamość własna eksperta: warstwy promptu i wtyczki,
	// katalog rozszerzeń powłoki.
	WarstwyAgenta RepozytoriumWarstwAgenta
	// WersjeAgenta i ArchiwumAgentow to jedna implementacja widziana dwoma
	// widokami.
	WersjeAgenta    RepozytoriumWersjiAgenta
	ArchiwumAgentow RepozytoriumArchiwumAgentow
	// ZakresAgenta to zakres działania eksperta: moduły, izolacja,
	// uprawnienia i konektory.
	ZakresAgenta RepozytoriumZakresuAgenta
	// ZakresyNarzedzi to zakresy uprawnień i limity wywołań narzędzi dla
	// profilu asystenta.
	ZakresyNarzedzi RepozytoriumZakresowNarzedzi
	// UkladOrkiestracji to bramka, grupa kroków, krok wycofujący i spięcie
	// kolejek automatyki.
	UkladOrkiestracji RepozytoriumUkladuOrkiestracji
	// Zespoly to nazwane składy biblioteki ekspertów; wskazują ekspertów
	// kodem, nie numerem wiersza.
	Zespoly RepozytoriumZespolow
	// PrzestrzenRobocza to projekt, jego pamięć i przypisania ekspertów —
	// trwałość modułu Workspace.
	PrzestrzenRobocza RepozytoriumPrzestrzeniRoboczej
	// Terminal to karty powłok i dziennik procesów; proces czynny prowadzi
	// rdzeń.
	Terminal RepozytoriumTerminala
	// Developer to wersje plików edytora i dziennik budowania; treść pliku
	// roboczego zostaje na dysku.
	Developer RepozytoriumDevelopera
	// Automatyki to definicje, harmonogramy i przebiegi Automations; wykonuje
	// je silnik kolejek.
	Automatyki RepozytoriumAutomatyk
	// Biblioteka to pliki repozytorium wiedzy, ich wersje, etykiety i
	// kolekcje; treść zostaje na dysku.
	Biblioteka RepozytoriumBiblioteki
	// Prowenancja to ślad wywołań kanału modelu wraz z drzewem odcinków,
	// źródło dla Explorer i rozliczeń.
	Prowenancja RepozytoriumProwenancji
	// Studio to dokumenty edytora, ich wersje i propozycje zmian modułu
	// Studio.
	Studio RepozytoriumStudia
	// Przegladanie to migawki stron, zebrane źródła i notatki modułu Browser.
	Przegladanie RepozytoriumPrzegladania
	// Design to prompty strukturalne, zasoby wizualne i kompozycje; rdzeń
	// zapisuje, nie generuje zasobów.
	Design RepozytoriumDesignu
	// Badania to źródła, ustalenia, raporty i przestrzeń modułu Research;
	// rdzeń tylko składa raport.
	Badania RepozytoriumBadan
	// Asystent to zlecenia i dziennik czynności modułu Assistant; zlecenie
	// ma własny automat stanu.
	Asystent RepozytoriumAsystenta
	// Aplikacje to architektura, pliki warsztatu i przebiegi wdrożeń Apps;
	// rdzeń niczego nie wdraża.
	Aplikacje RepozytoriumAplikacji
	// Przekazania to zlecenia przekazania między oknami i dziennik akcji
	// okna.
	Przekazania RepozytoriumPrzekazan
	// Tlumaczenia to cały moduł Translate: okno źródłowe, panele, słownik,
	// pamięć, kontrola jakości.
	Tlumaczenia RepozytoriumTlumaczen
	// Diagnostyka to dziennik, błędy, analizy i rekomendacje modułu
	// Diagnostics.
	Diagnostyka RepozytoriumDiagnostyki
	// Roundtable to skład debaty, jej tury, wypowiedzi i stanowisko końcowe.
	Roundtable RepozytoriumRoundtable
	// Komponenty to rejestr kafli strefy komponentów Strony głównej,
	// wskazujący byt kolumną byt_docelowy.
	Komponenty RepozytoriumKomponentow
	// ZdarzeniaWykonawcze to dziennik zdarzeń zaczepów i zamknięcia tur.
	ZdarzeniaWykonawcze RepozytoriumZdarzenWykonawczych
	// Bloki to nietekstowe fragmenty strumienia odpowiedzi z tury:
	// rozumowanie, narzędzia, prowenancja.
	Bloki RepozytoriumBlokow
	// Historia to wykaz pozycji rozmowy okna z zasadą przechowywania; pyta
	// o wiadomosc inaczej niż inne.
	Historia RepozytoriumHistorii
	// SzukanieRozmow to odczyt indeksu pełnotekstowego (FTS5) nad
	// `wiadomosc.tresc`.
	SzukanieRozmow RepozytoriumSzukaniaRozmow
	// KoszSesji to odwrotna strona tabeli `sesja`: widzi wyłącznie wiersze
	// ze znacznikiem `usunieto_o`.
	KoszSesji RepozytoriumKoszaSesji
	// Uwierzytelnienie to bramka Operatora: metody wejścia i sesje; sekret
	// nie leży tu w żadnej postaci.
	Uwierzytelnienie RepozytoriumUwierzytelnienia
	// KontoWlasciciela mówi, czyja jest bramka: login, e-mail, stan
	// potwierdzenia i drogi odzyskiwania.
	KontoWlasciciela RepozytoriumKontaWlasciciela
	// Kondycja i Alerty opisują stan produktu: sondy z pomiarami oraz reguły
	// wyzwalania z wyzwoleniami.
	Kondycja RepozytoriumKondycji
	Alerty   RepozytoriumAlertow
	// Schowek i SkrotyTekstowe należą do rdzenia, by historia schowka
	// i skrót nie ginęły razem z kartą.
	Schowek        RepozytoriumSchowka
	SkrotyTekstowe RepozytoriumSkrotow
	// KontekstyPamieci dopełnia rodzinę `memory.*`: nazwane zestawy wskazań
	// i zasady retencji.
	KontekstyPamieci RepozytoriumKontekstowPamieci
	// WylaczeniaPamieci trzyma wyłączenia pamięci w zasięgu; to nie jest
	// wpis pamięci ani jego zmiana.
	WylaczeniaPamieci RepozytoriumWylaczenPamieci
	// WyciszeniaNakladki trzyma wyciszenia nakładki Always On Display wraz
	// z sygnałami wyzwalającymi.
	WyciszeniaNakladki RepozytoriumWyciszenNakladki
	// NagraniaMowy jest rejestrem bajtów przyjętych od okna, odkładanych
	// przez mikrofon karty.
	NagraniaMowy RepozytoriumNagranMowy

	zapytania *zapytania
	// baza jest połączeniem potrzebnym repozytoriom składanym na żądanie do
	// transakcji.
	baza *sql.DB
}

// Otworz składa repozytoria nad otwartą bazą. Baza pochodzi z pakietu `store`
// — warstwa dostępu do danych nie otwiera pliku ani nie prowadzi migracji.
func Otworz(ctx context.Context, baza *store.Baza) (*Zestaw, error) {
	if baza == nil || baza.DB == nil {
		return nil, fmt.Errorf("dane: baza nie jest otwarta")
	}
	zapytania := noweZapytania(baza.DB)
	// Jedna instancja: historia tożsamości i archiwum eksperta stoją na tej
	// samej tabeli.
	wersjeAgenta := noweRepozytoriumWersjiAgenta(zapytania, baza.DB)
	konfiguracja := noweRepozytoriumKonfiguracji(zapytania)
	return &Zestaw{
		Sesje:               noweRepozytoriumSesji(zapytania),
		Okna:                noweRepozytoriumOkien(zapytania, baza.DB),
		Wiadomosci:          noweRepozytoriumWiadomosci(zapytania),
		Kanaly:              noweRepozytoriumKanalow(zapytania),
		Konta:               noweRepozytoriumKont(zapytania, baza.DB),
		Konfiguracja:        konfiguracja,
		KonfiguracjaOsi:     konfiguracja,
		KatalogUstawien:     noweRepozytoriumKatalogUstawien(zapytania),
		Kolejki:             noweRepozytoriumKolejek(zapytania, baza.DB),
		Pamiec:              noweRepozytoriumPamieci(zapytania),
		Srodowiska:          noweRepozytoriumSrodowisk(zapytania),
		Moduly:              noweRepozytoriumModulow(zapytania),
		Macierz:             noweRepozytoriumMacierzy(zapytania),
		KartySesji:          noweRepozytoriumKartSesji(zapytania),
		OknaOperacyjne:      noweRepozytoriumOkienOperacyjnych(zapytania),
		Akcje:               noweRepozytoriumAkcji(zapytania),
		PunktyDostepu:       noweRepozytoriumPunktowDostepu(zapytania, baza.DB),
		Nadania:             noweRepozytoriumNadan(zapytania, baza.DB),
		Tozsamosc:           noweRepozytoriumTozsamosci(zapytania),
		Urzadzenia:          noweRepozytoriumUrzadzen(zapytania, baza.DB),
		CentrumPowiadomien:  noweRepozytoriumCentrumPowiadomien(zapytania),
		PrzestrzenRobocza:   noweRepozytoriumPrzestrzeniRoboczej(zapytania, baza.DB),
		Agenci:              noweRepozytoriumAgentow(zapytania, baza.DB),
		WarstwyAgenta:       noweRepozytoriumWarstwAgenta(zapytania, baza.DB),
		WersjeAgenta:        wersjeAgenta,
		ArchiwumAgentow:     wersjeAgenta,
		ZakresAgenta:        noweRepozytoriumZakresuAgenta(zapytania, baza.DB),
		ZakresyNarzedzi:     noweRepozytoriumZakresowNarzedzi(zapytania, baza.DB),
		UkladOrkiestracji:   noweRepozytoriumUkladuOrkiestracji(zapytania, baza.DB),
		Zespoly:             noweRepozytoriumZespolow(zapytania, baza.DB),
		Terminal:            noweRepozytoriumTerminala(zapytania, baza.DB),
		Developer:           noweRepozytoriumDevelopera(zapytania),
		Automatyki:          noweRepozytoriumAutomatyk(zapytania, baza.DB),
		Roundtable:          noweRepozytoriumRoundtable(zapytania, baza.DB),
		Diagnostyka:         noweRepozytoriumDiagnostyki(zapytania, baza.DB),
		Biblioteka:          noweRepozytoriumBiblioteki(zapytania, baza.DB),
		Prowenancja:         noweRepozytoriumProwenancji(zapytania, baza.DB),
		Studio:              noweRepozytoriumStudia(zapytania, baza.DB),
		Przegladanie:        noweRepozytoriumPrzegladania(zapytania, baza.DB),
		Design:              noweRepozytoriumDesignu(zapytania, baza.DB),
		Badania:             noweRepozytoriumBadan(zapytania, baza.DB),
		Asystent:            noweRepozytoriumAsystenta(zapytania, baza.DB),
		Aplikacje:           noweRepozytoriumAplikacji(zapytania, baza.DB),
		Przekazania:         noweRepozytoriumPrzekazan(zapytania, baza.DB),
		Tlumaczenia:         noweRepozytoriumTlumaczen(zapytania, baza.DB),
		Komponenty:          noweRepozytoriumKomponentow(zapytania),
		Uwierzytelnienie:    noweRepozytoriumUwierzytelnienia(zapytania),
		KontoWlasciciela:    noweRepozytoriumKontaWlasciciela(zapytania),
		ZdarzeniaWykonawcze: noweRepozytoriumZdarzenWykonawczych(zapytania),
		Bloki:               noweRepozytoriumBlokow(zapytania),
		Historia:            noweRepozytoriumHistorii(zapytania, baza.DB),
		SzukanieRozmow:      noweRepozytoriumSzukaniaRozmow(zapytania),
		KoszSesji:           noweRepozytoriumKoszaSesji(zapytania, baza.DB),
		Kondycja:            noweRepozytoriumKondycji(zapytania, baza.DB),
		Alerty:              noweRepozytoriumAlertow(zapytania, baza.DB),
		Schowek:             noweRepozytoriumSchowka(zapytania, baza.DB),
		SkrotyTekstowe:      noweRepozytoriumSkrotow(zapytania, baza.DB),
		KontekstyPamieci:    noweRepozytoriumKontekstowPamieci(zapytania, baza.DB),
		WylaczeniaPamieci:   noweRepozytoriumWylaczenPamieci(zapytania),
		WyciszeniaNakladki:  noweRepozytoriumWyciszenNakladki(zapytania),
		NagraniaMowy:        noweRepozytoriumNagranMowy(zapytania, baza.DB),
		zapytania:           zapytania,
		baza:                baza.DB,
	}, nil
}

// Zamknij zwalnia przygotowane zapytania zestawu; bazy nie zamyka, bo zamyka
// ją zawsze ten, kto ją otworzył.
func (z *Zestaw) Zamknij() error {
	if z == nil || z.zapytania == nil {
		return nil
	}
	return z.zapytania.zamknij()
}
