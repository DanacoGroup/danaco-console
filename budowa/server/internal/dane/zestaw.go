// Pakiet dane zawiera repozytoria nad pakietem `internal/store`: jedno na
// obszar, każde w osobnym pliku. Warstwy wyższe (rdzeń, sesje, modele,
// konfiguracja warstwowa) sięgają po dane wyłącznie przez interfejsy tego
// pakietu, nie przez SQL ani przez wnętrze `store`.
//
// Odpowiedzialność tego pliku: złożenie repozytoriów w jeden zestaw i zwolnienie
// zasobów. Zero SQL, zero logiki obszarowej.
package dane

import (
	"context"
	"database/sql"
	"fmt"

	"danacoconsole/server/internal/store"
)

// Zestaw to komplet repozytoriów jednej otwartej bazy.
type Zestaw struct {
	Sesje        RepozytoriumSesji
	Okna         RepozytoriumOkien
	Wiadomosci   RepozytoriumWiadomosci
	Kanaly       RepozytoriumKanalow
	Konta        RepozytoriumKont
	Konfiguracja RepozytoriumKonfiguracji
	// KonfiguracjaOsi to ten sam byt co Konfiguracja, widziany razem z osią
	// rozstrzygania. Jedna implementacja, dwa widoki.
	KonfiguracjaOsi RepozytoriumKonfiguracjiOsi
	KatalogUstawien RepozytoriumKatalogUstawien
	Kolejki         RepozytoriumKolejek
	Pamiec          RepozytoriumPamieci
	Srodowiska      RepozytoriumSrodowisk
	Moduly          RepozytoriumModulow
	// Macierz to ta sama tabela `srodowisko_modul`, ktora `Moduly` zna wylacznie
	// jako filtr wykazu. Tu bytem jest wiersz macierzy — para kodow, kolejnosc
	// i widocznosc — czytany w calosci jednym zapytaniem. Jedna tabela, dwa
	// pytania; drugiej prawdy o module nie ma.
	Macierz        RepozytoriumMacierzy
	KartySesji     RepozytoriumKartSesji
	OknaOperacyjne RepozytoriumOkienOperacyjnych
	Akcje          RepozytoriumAkcji
	PunktyDostepu  RepozytoriumPunktowDostepu
	Nadania        RepozytoriumNadan
	Tozsamosc      RepozytoriumTozsamosci
	// Urzadzenia to katalog maszyn (tabela `urzadzenie`). Punkt dostępu rodzaju
	// `localDirectory` wskazuje tu urządzenie przez `urzadzenie_id`,
	// a rozpoznanie startowe zakłada wiersz maszyny bieżącej.
	Urzadzenia RepozytoriumUrzadzen
	// CentrumPowiadomien to trwały rejestr zdarzeń centrum powiadomień
	// (tabela `powiadomienie_centrum`). To NIE jest kolejka doręczeń funkcji
	// Mobile — tamta mieszka w tabeli `powiadomienie` i ma własny cykl życia.
	CentrumPowiadomien RepozytoriumCentrumPowiadomien
	// Agenci to biblioteka ekspertów modułu Agents. Ekspert jest komponentem
	// własnym, więc żyje obok sesji, nie w niej.
	Agenci RepozytoriumAgentow
	// WarstwyAgenta to tożsamość własna eksperta: warstwy jego promptu i jego
	// wtyczki. Ta sama biblioteka ekspertów, inny byt — warstwa
	// jest stanem promptu, a wtyczka katalogiem rozszerzeń powłoki, odrębnym
	// od konektora.
	WarstwyAgenta RepozytoriumWarstwAgenta
	// WersjeAgenta i ArchiwumAgentow to jedna implementacja widziana dwoma
	// widokami — historia tożsamości i archiwum eksperta.
	WersjeAgenta    RepozytoriumWersjiAgenta
	ArchiwumAgentow RepozytoriumArchiwumAgentow
	// ZakresAgenta to zakres działania eksperta: moduły zastosowania, osiem
	// zakresów izolacji technicznej, granica Subagent Network, zdjęcie wpisów
	// uprawnień oraz odczyt konektorów i przypisań od strony eksperta.
	// Repozytorium osobne od biblioteki: odpowiada nie na pytanie „jaki jest ten
	// ekspert", lecz „co temu ekspertowi wolno zrobić w systemie".
	ZakresAgenta RepozytoriumZakresuAgenta
	// ZakresyNarzedzi to zakresy uprawnień i limity wywołań pozycji katalogu
	// narzędzi dla profilu asystenta.
	ZakresyNarzedzi RepozytoriumZakresowNarzedzi
	// UkladOrkiestracji to trzy dopełnienia układu zależności — bramka
	// dołączenia, grupa kroków i krok wycofujący — oraz spięcie kolejek
	// automatyki z silnikiem kolejek środowiska MultitaskingAI.
	UkladOrkiestracji RepozytoriumUkladuOrkiestracji
	// Zespoly to nazwane składy biblioteki ekspertów. Zespół wskazuje ekspertów
	// kodem, nie numerem wiersza, i nie powiela ich tożsamości — prawdą
	// o ekspercie zostaje tabela `agent`.
	Zespoly RepozytoriumZespolow
	// PrzestrzenRobocza to projekt, jego pamięć i przypisania ekspertów —
	// trwałość modułu Workspace.
	PrzestrzenRobocza RepozytoriumPrzestrzeniRoboczej
	// Terminal to karty powłok i dziennik procesów modułu Terminal. Proces
	// czynny prowadzi rdzeń; tu leży ślad po nim.
	Terminal RepozytoriumTerminala
	// Developer to wersje plików edytora i dziennik przebiegów budowania.
	// Treść pliku roboczego zostaje na dysku — tutaj leży wyłącznie migawka
	// założona na żądanie i ślad po budowaniu.
	Developer RepozytoriumDevelopera
	// Automatyki to definicje, harmonogramy i przebiegi modułu Automations.
	// Wykonaniem kroków zajmuje się silnik kolejek — to repozytorium opisuje
	// wyłącznie definicję i zapis przebiegu.
	Automatyki RepozytoriumAutomatyk
	// Biblioteka to pliki repozytorium wiedzy, ich wersje, etykiety i kolekcje
	// modułu Library. Treść pliku zostaje na dysku — baza trzyma odwołanie.
	Biblioteka RepozytoriumBiblioteki
	// Prowenancja to ślad wywołań kanału modelu wraz z drzewem odcinków. Jedno
	// źródło dla Provenance Explorer i dla rozliczenia zużycia — obie
	// odpowiedzi powstają z tych samych wierszy.
	Prowenancja RepozytoriumProwenancji
	// Studio to dokumenty edytora, ich wersje i propozycje zmian modułu Studio.
	// Wykaz wersji jest repozytorium dokumentu — Repository Panel nie ma
	// osobnej struktury.
	Studio RepozytoriumStudia
	// Przegladanie to migawki stron, zebrane źródła i notatki modułu Browser.
	// Historia nawigacji to kolejne migawki, nie osobna tabela.
	Przegladanie RepozytoriumPrzegladania
	// Design to prompty strukturalne, zasoby wizualne i kompozycje modułu
	// Design. Rdzeń zasobów nie generuje — zapisuje to, co o nich wie.
	Design RepozytoriumDesignu
	// Badania to źródła, ustalenia, raporty i przestrzeń modułu Research.
	// Rdzeń badań nie prowadzi — składa raport z tego, co Operator zebrał.
	Badania RepozytoriumBadan
	// Asystent to zlecenia i dziennik czynności modułu Assistant. Zlecenie wisi
	// na oknie i ma własny automat stanu — to inny byt niż pozycja kolejki
	// sesyjnej i niż statyczny katalog akcji.
	Asystent RepozytoriumAsystenta
	// Aplikacje to architektura, pliki warsztatu i przebiegi wdrożeń modułu
	// Apps. Rdzeń niczego nie wdraża — trzyma ślad zlecenia i stan.
	Aplikacje RepozytoriumAplikacji
	// Przekazania to zlecenia przekazania między oknami i dziennik akcji okna.
	// Sama więź koordynator–wykonawca mieszka w kolumnie
	// `okno_komunikacji.okno_koordynatora_id` — tu jej nie ma drugi raz.
	Przekazania RepozytoriumPrzekazan
	// Tlumaczenia to cały moduł Translate: okno źródłowe, panele języków
	// docelowych, słownik, pamięć tłumaczeń, kontrola jakości, ślady syntezy
	// mowy i eksportów. Jedno repozytorium, nie trzy.
	// Rdzeń nie tłumaczy i nie rozpoznaje języka; trzyma to, co dostał.
	Tlumaczenia RepozytoriumTlumaczen
	// Diagnostyka to dziennik, błędy, analizy i rekomendacje modułu
	// Diagnostics. Przenosi wyłącznie fakty zgłoszone przez rdzeń — nie liczy
	// stanu systemu i nie wytwarza rekomendacji.
	Diagnostyka RepozytoriumDiagnostyki
	// Roundtable to skład debaty, jej tury, wypowiedzi i stanowisko końcowe.
	// Wywołanie kanałów uczestników prowadzi rdzeń przez rejestr kanałów —
	// tu leży wyłącznie zapis debaty.
	Roundtable RepozytoriumRoundtable
	// Komponenty to rejestr kafli strefy komponentów Strony głównej (rodzina
	// `component.*`). Kafel wskazuje byt magazynu modułowego kolumną
	// `byt_docelowy` i nie powiela go — prawdą o projekcie, ekspercie
	// i automatyce pozostają ich własne tabele i własne komendy.
	Komponenty RepozytoriumKomponentow
	// ZdarzeniaWykonawcze to dziennik zdarzeń zaczepów i zamknięcia tur. Ślad
	// po tym, co zaszło w turze — zapis następuje po zdarzeniu i niczego
	// nie steruje.
	ZdarzeniaWykonawcze RepozytoriumZdarzenWykonawczych
	// Bloki to nietekstowe fragmenty strumienia odpowiedzi zapisane w trakcie
	// tury: rozumowanie, narzędzia, prowenancja. Tekst wypowiedzi mieszka
	// w `wiadomosc.tresc` i tu go nie ma.
	Bloki RepozytoriumBlokow
	// Historia to wykaz pozycji rozmowy okna (rodzina `history.*`) wraz z zasadą
	// przechowywania. Osobnej tabeli historii nie ma: pozycją jest wiersz
	// `wiadomosc`, a to repozytorium pyta o niego z drugiej
	// strony niż `Wiadomosci` — po identyfikatorze kontraktowym okna, od
	// najnowszej, kursorem czasu — i jako jedyne kasuje. Jedna tabela, dwa
	// pytania, wzorem par Moduly/Macierz i Sesje/KoszSesji.
	Historia RepozytoriumHistorii
	// SzukanieRozmow to odczyt indeksu pełnotekstowego nad `wiadomosc.tresc`
	// (FTS5). Indeks jest zewnętrzny (content=) — treść ma jedną
	// prawdę w tabeli wiadomości, repozytorium wyłącznie pyta.
	SzukanieRozmow RepozytoriumSzukaniaRozmow
	// KoszSesji to odwrotna strona tabeli `sesja`: widzi
	// wyłącznie wiersze ze znacznikiem `usunieto_o`, których wykaz sesji
	// żywych nie widzi wcale. Jedna tabela, dwa pytania.
	KoszSesji RepozytoriumKoszaSesji
	// Uwierzytelnienie to bramka Operatora: metody wejścia i sesje bramki
	// (rodzina `auth.*`). Mówi, CZYM otworzyć bramkę. Sekret nie leży tu
	// w żadnej postaci: baza zna wyłącznie odwołanie do sejfu poświadczeń.
	Uwierzytelnienie RepozytoriumUwierzytelnienia
	// KontoWlasciciela mówi, CZYJA jest bramka: login i adres e-mail
	// uwierzytelniający jedynego właściciela, jego stan potwierdzenia oraz
	// jednorazowe drogi potwierdzenia tożsamości wysyłane listem przy
	// rejestracji i przy odzyskiwaniu konta.
	KontoWlasciciela RepozytoriumKontaWlasciciela
	// Kondycja i Alerty to dwie rodziny przekrojowe opisujące STAN PRODUKTU:
	// sondy wraz z serią pomiarów oraz reguły wyzwalania wraz z rejestrem
	// wyzwoleń. Magazyny są osobne, bo osobne są pytania: „czy to działa"
	// i „kiedy mam zawołać".
	Kondycja RepozytoriumKondycji
	Alerty   RepozytoriumAlertow
	// Schowek i SkrotyTekstowe należą do rdzenia, choć obsługują pola tekstowe
	// klienta: historia schowka ginęłaby razem z kartą, a skrót rozwijany
	// w jednym oknie rozwijałby się inaczej w drugim.
	Schowek        RepozytoriumSchowka
	SkrotyTekstowe RepozytoriumSkrotow
	// KontekstyPamieci dopełnia rodzinę `memory.*`: nazwane zestawy wskazań
	// i zasady retencji. Wpisy pamięci zostają tam, gdzie były
	// (`PrzestrzenRobocza`) — kontekst jest wskazaniem, nie właścicielem treści.
	KontekstyPamieci RepozytoriumKontekstowPamieci
	// WylaczeniaPamieci trzyma wyłączenia pamięci w zasięgu: całkiem,
	// w środowisku, w projekcie, w module, w parze modułów, w karcie sesji.
	// Osobno od `Pamiec` i `PrzestrzenRobocza`, bo wyłączenie NIE JEST wpisem
	// pamięci ani jego zmianą — nie dotyka treści i znosi się jednym ruchem.
	WylaczeniaPamieci RepozytoriumWylaczenPamieci
	// WyciszeniaNakladki trzyma wyciszenia nakładki Always On Display wraz
	// z sygnałami klas zdarzeń wyzwalających. Wyciszenie jest bytem rdzenia, nie
	// stanem jednego okna: bez wiersza Operator wyciszał w jednej powłoce,
	// a w drugiej sugestie wchodziły dalej.
	WyciszeniaNakladki RepozytoriumWyciszenNakladki
	// NagraniaMowy jest rejestrem bajtów przyjętych od okna: mikrofon karty ma
	// dokąd odłożyć nagranie, a silnik mowy dostaje ścieżkę, jakiej oczekuje.
	NagraniaMowy RepozytoriumNagranMowy

	zapytania *zapytania
	// baza jest połączeniem, którego potrzebują repozytoria składane na żądanie
	// (`Rozszerzenia`, `NarzedziaSesji`) do transakcji — te, które powstają
	// w `Otworz`, dostają je wprost w konstruktorze.
	baza *sql.DB
}

// Otworz składa repozytoria nad otwartą bazą. Baza pochodzi z pakietu `store`
// — warstwa dostępu do danych nie otwiera pliku ani nie prowadzi migracji.
func Otworz(ctx context.Context, baza *store.Baza) (*Zestaw, error) {
	if baza == nil || baza.DB == nil {
		return nil, fmt.Errorf("dane: baza nie jest otwarta")
	}
	zapytania := noweZapytania(baza.DB)
	// Jedna instancja, dwa widoki: historia tożsamości i archiwum eksperta stoją
	// na tej samej tabeli, więc drugie repozytorium byłoby drugim źródłem
	// prawdy o tym samym wierszu.
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

// Zamknij zwalnia przygotowane zapytania. Bazy nie zamyka — zamyka ją ten, kto
// ją otworzył.
func (z *Zestaw) Zamknij() error {
	if z == nil || z.zapytania == nil {
		return nil
	}
	return z.zapytania.zamknij()
}
