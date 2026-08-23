package core

import (
	"log"

	"danacoconsole/shared"
)

// Porty jest kompletem zależności rdzenia. Każde pole jest interfejsem z
// porty.go, więc rdzeń nie wie, kto je wypełnia: pakiet sesji, pakiet
// danych, pakiet modeli albo — w teście — atrapa portu.
//
// Pole niewypełnione nie zatrzymuje startu. Domena bez portu nie rejestruje
// obsługiwaczy, a jej komendy odpowiadają `*.unknown`; pozostałe domeny pracują
// bez zmian. Rdzeń niczego nie wymaga i niczym nie warunkuje startu.
type Porty struct {
	Sesje         Sesje
	Okna          Okna
	Rozmowa       Rozmowa
	Ustawienia    Ustawienia
	Kanaly        Kanaly
	Kolejki       Kolejki
	Przenoszenie  Przenoszenie
	Nawigacja     Nawigacja
	WiazanieSesji WiazanieSesji
	Akcje         Akcje
	// CentrumPowiadomien jest rejestrem zdarzeń centrum powiadomień. Zgłasza
	// do niego sam rdzeń; kontrakt niesie odczyt i zmianę stanu.
	CentrumPowiadomien CentrumPowiadomien
	// Pięć portów okna konfiguracji: katalog pozycji, punkty dostępu i nadania
	// per okno rozmowy, rejestr kont oraz zasady i tożsamość modelu.
	KatalogUstawien KatalogUstawien
	PunktyDostepu   PunktyDostepu
	NadaniaDostepu  NadaniaDostepu
	Konta           Konta
	Tozsamosc       Tozsamosc
	// Agenci obsługuje moduł Agents: bibliotekę ekspertów wraz z ich modelem
	// bazowym, umiejętnościami, konektorami i uprawnieniami.
	Agenci Agenci
	// PrzestrzenRobocza obsługuje moduł Workspace: projekt, jego instrukcje
	// warstwowe, pamięć, bibliotekę i przypisanych ekspertów.
	PrzestrzenRobocza PrzestrzenRobocza
	// Terminal obsługuje moduł Terminal: karty powłok, rejestr procesów rdzenia
	// i strumień wyjścia.
	Terminal Terminal
	// Developer obsługuje moduł Developer: plik repozytorium, drzewo projektu,
	// czynności repozytorium i przebiegi budowania. Procesy startuje tym samym
	// uruchamiaczem co Terminal.
	Developer Developer
	// Automatyki obsługuje moduł Automations: definicje, harmonogramy, układ
	// zależności i przebiegi. Kroki wykonuje silnik kolejek — ten sam, którym
	// pracuje pętla sesyjna i MultitaskingAI.
	Automatyki Automatyki
	// Biblioteka obsługuje moduł Library: pliki repozytorium wiedzy, ich wersje,
	// etykiety i kolekcje.
	Biblioteka Biblioteka
	// Studio obsługuje moduł Studio: dokumenty edytora, ich wersje i propozycje
	// zmian. Wykaz wersji jest repozytorium dokumentu.
	Studio Studio
	// Przegladarka obsługuje moduł Browser: migawki stron, zebrane źródła
	// i notatki. Rdzeń pobiera stronę sam (HTTP GET z granicami protokołu,
	// czasu, rozmiaru i przekierowań); nie renderuje jej JavaScriptu.
	Przegladarka Przegladarka
	// Design obsługuje moduł Design: prompty strukturalne, zasoby wizualne
	// i kompozycje.
	Design Design
	// Badania obsługuje moduł Research: źródła, ustalenia, raport i przestrzeń.
	Badania Badania
	// Asystent obsługuje moduł Assistant: polecenia, stan zleceń i dziennik
	// czynności.
	Asystent Asystent
	// Aplikacje obsługuje moduł Apps: architektura, warsztat i wdrożenia.
	Aplikacje Aplikacje
	// Komponenty obsługuje rodzinę `component.*` — kafle Strefy 2 Strony głównej.
	// Kafel wskazuje byt magazynu modułowego; prawda o projekcie, ekspercie
	// i automatyce zostaje w ich własnych tabelach i komendach.
	Komponenty Komponenty
	// PrzekazanieOkna obsługuje obszar window.*: przekazanie zlecenia między
	// oknami i akcje panelu okna. `action.list` tu nie należy — obsługuje ją
	// `zarejestrujAkcje` niżej.
	PrzekazanieOkna PrzekazanieOkna
	// Tlumaczenie obsługuje cały moduł Translate — szesnaście komend trzech
	// obszarów (przekład, słownik, jakość i mowa) na jednej powierzchni.
	// Trzy porty nad jednym repozytorium byłyby trzema prawdami o jednym
	// module.
	Tlumaczenie Tlumaczenie
	// Diagnostyka obsługuje moduł Diagnostics: dziennik rdzenia, wykaz błędów,
	// analizy i rekomendacje. Jest też odbiorcą odmów wykonania komend — dzięki
	// temu Errors Panel pokazuje odmowy jako fakty, a nie jako ciszę.
	Diagnostyka Diagnostyka
	// Prowenancja i Zuzycie są portami przekrojowymi, nie modułowymi: ślad
	// wywołania modelu czyta Provenance Explorer, ale też Execution Loop Window,
	// a rozliczenie zużycia — Always On Display. Jeden magazyn, dwa porty.
	Prowenancja Prowenancja
	Zuzycie     Zuzycie
	// Kondycja i Alerty opisują STAN PRODUKTU: sondy wraz z serią pomiarów oraz
	// reguły wyzwalania wraz z rejestrem wyzwoleń. Porty przekrojowe, nie
	// modułowe — kondycję czyta Health & Uptime Panel, ale też Always On Display.
	Kondycja Kondycja
	Alerty   Alerty
	// Schowek, SkrotyTekstowe i KontekstyPamieci obsługują tekst poza jednym
	// oknem: historię schowka, słownik skrótów rozwijanych we wszystkich polach
	// platformy oraz nazwane zestawy pamięci wraz z zasadami retencji.
	Schowek          Schowek
	SkrotyTekstowe   SkrotyTekstowe
	KontekstyPamieci KontekstyPamieci
	// Wywolywacz obsługuje rodzinę `launcher.*` — skrót globalny otwierający
	// wywoływacz poleceń. Nastawę trzyma rdzeń, skrót przechwytuje powłoka.
	Wywolywacz Wywolywacz
	// ZajetoscKontekstu obsługuje `context.usage.get` — pomiar zajętości okna
	// kontekstu tokenizatorem wkompilowanym w rdzeń.
	ZajetoscKontekstu ZajetoscKontekstu
	// zdolnosciKlienta jest mapą deklaracji spisanych w powitaniu. Pole jest
	// nieeksportowane, bo nie jest portem: wypełnia je montaż rdzenia, a czyta
	// wyłącznie `launcher.hotkey.*`. Pole puste znaczy, że rdzeń nie zna
	// deklaracji żadnego klienta — wtedy skrót globalny melduje się jako
	// niewspierany, co jest prawdą, a nie zgadywaniem.
	zdolnosciKlienta *zdolnosciKlientow
	// WarsztatPdf jest częścią modułu Studio o innym komplecie zależności:
	// materiał wchodzi i wychodzi zasobem magazynu, a czynność pracuje
	// biblioteką wkompilowaną w rdzeń. Port osobny, bo zależności osobne.
	WarsztatPdf WarsztatPdf
	// BezpieczenstwoDokumentu stoi obok warsztatu i na tych samych zależnościach,
	// ale port jest osobny, bo rodzina jest osobna: czynność, która dokument
	// otwiera, i czynność, która go zamyka, nie mają wspólnego powodu do awarii.
	BezpieczenstwoDokumentu BezpieczenstwoDokumentu
	// Debata obsługuje moduł Roundtable: skład debaty, tury, wypowiedzi
	// i stanowisko końcowe. Jest jedynym portem, którego jedno okno rozmawia
	// z wieloma kanałami naraz — po tym samym rejestrze kanałów, którym jedzie
	// okno rozmowy.
	Debata Debata
	// Izolacja obsługuje rodzinę `isolation.*`: profile izolacji, warstwy,
	// zakresy techniczne i podgląd polityki. Profil rozstrzyga, co okno widzi
	// z historii, pamięci i kontekstu sąsiada — dlatego port bierze rozstrzygacz
	// polityki, a nie samo repozytorium.
	Izolacja Izolacja
	// Uwierzytelnianie obsługuje rodzinę `auth.*`: bramka, metody i sesja
	// Operatora (`dane/auth.go`).
	Uwierzytelnianie Uwierzytelnianie
	// Urzadzenia obsługuje rodzinę `device.*` — wykaz maszyn powiązanych
	// z kontem i unieważnienie ich tokenów. Stoi na tym samym repozytorium co
	// uwierzytelnianie, bo urządzeniem konta jest to, które weszło przez bramkę.
	Urzadzenia Urzadzenia
	// NakladkaAod obsługuje rodzinę `aod.*` — nakładkę zawsze-na-wierzchu.
	// Stoi na nadzorcy sesji, telemetrii postępu i rejestrze obecności, a nie
	// na własnym repozytorium: nakładka pokazuje stan biegnącej pracy, więc
	// druga tabela byłaby drugą prawdą o tym samym.
	NakladkaAod NakladkaAod
	// Rozszerzenia obsługuje rodzinę `extension.*` — katalog rozszerzeń.
	// `Zestaw.Rozszerzenia()` stoi metodą, nie polem, bo rejestr nie trzyma
	// stanu poza wskaźnikiem na wspólną pamięć zapytań.
	Rozszerzenia Rozszerzenia
	// Role obsługuje rodzinę `role.*` — rolę okna, więź koordynatora, wykaz
	// nadań i zdjęcie roli. `Zestaw.RoleOkien()` stoi metodą, bo rola okna nie
	// jest osobnym obszarem danych, tylko widokiem na dwie kolumny obszaru okien.
	// Port jeden, nie dwa: nadanie roli to para pól okna, więc wykaz jest
	// widokiem na ten sam adapter, a nie drugim repozytorium.
	Role RoleWykaz
	// Panele obsługuje `panel.sections.*` — układ sekcji panelu okna: kolejność,
	// zwinięcie, zdjęcie z widoku.
	Panele SekcjePaneli
	// Historia obsługuje rodzinę `history.*` i `retention.set` — wykaz pozycji
	// rozmowy okna oraz zasadę ich przechowywania. Port osobny, nie rozszerzenie
	// portu Rozmowa: tamten prowadzi turę (wysyła, zatrzymuje, wykazuje), a ten
	// prowadzi zapis po turze — czyta go wstecz kursorem czasu i jako jedyny
	// kasuje.
	Historia Historia

	// Podagenci obsługuje rodzinę `subagent.*` — powołanie podagentów okna
	// wykonawcy, ich wykaz i zbieranie wyników. Port osobny, nie rozszerzenie
	// portu Automatyki: podagent wisi na oknie wykonawcy, a nie na zapisanej
	// definicji automatyki. Pracę wykonuje ten sam silnik kolejek — drugiego
	// silnika nie ma.
	Podagenci Podagenci

	// WersjeEksperta obsługuje historię tożsamości i archiwum eksperta.
	// Port osobny, nie rozszerzenie portu Agenci: historia i archiwum siedzą
	// w innych tabelach i mają własne repozytoria, a wtopienie ich w Agenci
	// rozdęłoby port biblioteki o pięć czynności niezwiązanych z biblioteką.
	WersjeEksperta WersjeEksperta

	// ZakresEksperta obsługuje zakres działania eksperta: Permissions Center
	// w całości oraz cztery dopełnienia okien Skills, Connectors i historii
	// wersji. Port osobny, nie rozszerzenie portu Agenci: tamten opisuje
	// bibliotekę, ten — co ekspertowi wolno zrobić w systemie.
	ZakresEksperta ZakresEksperta

	// ZakresyNarzedzi obsługuje rodzinę `tools.scope.*` — zakresy uprawnień
	// i limity wywołań pozycji katalogu dla profilu asystenta. Ten sam adapter
	// jest strażą, którą rdzeń pyta przed skierowaniem komendy ręki modelu.
	ZakresyNarzedzi ZakresyNarzedzi

	// Zespoly obsługuje rodzinę `team.*` — zapisane składy ekspertów. Port
	// osobny, nie rozszerzenie portu Agenci: zespół jest bytem Operatora nad
	// biblioteką, a nie kolejną własnością eksperta.
	Zespoly Zespoly
	// Mowa obsługuje rodzinę `speech.*` — sprawdzenie gotowości silnika
	// rozpoznawania i zamianę nagrania na tekst (pakiet `mowa`).
	//
	// Synteza mowy zostaje przy porcie `Tlumaczenie`, bo idzie w drugą stronę
	// (tekst na dźwięk) i dotyczy panelu tłumaczenia, nie zdolności platformy.
	// Rdzeń syntezuje syntezatorem lokalnym `espeak-ng`, wpiętym tym samym
	// uruchamiaczem i przez tę samą bramę izolacji, co pomocnik rozpoznawania
	// (`adapter_modul_tlumaczenie_mowa_silnik.go`).
	Mowa Mowa
	// Wiedza obsługuje rodzinę `knowledge.*` — wskaźnik znaczenia wiedzy
	// Operatora i wyszukiwanie po sensie, a nie po literach (pakiet `wiedza`).
	// Uzupełnia to, czego nie umie indeks słów, i tego braku nie ukrywa: bez
	// silnika osadzeń obie komendy odmawiają, nazywając brak, zamiast zejść po
	// cichu na `library.file.search`.
	Wiedza Wiedza
	// NarzedziaObrazu obsługuje rodzinę `image.*` — cztery narzędzia modelu do
	// pracy na obrazie (zbadanie, geometria, retusz, format). Port osobny, nie
	// rozszerzenie portu Design: moduł Design jest oknem Operatora nad zasobami,
	// a to są czynności modelu wołane z rozmowy. Wynik i tak ląduje w magazynie
	// zasobów Designu — drugiego magazynu bajtów nie ma.
	NarzedziaObrazu NarzedziaObrazu
	// NarzedziaObrazuModelu obsługuje `image.upscale` i `image.background.remove`
	// — dwie czynności, które nie przerabiają pikseli, tylko puszczają nad nimi
	// sieć neuronową (Real-ESRGAN, U²-Net). Port osobny od NarzedziaObrazu
	// wyżej, bo silniki są osobne i bywają nieobecne osobno: maszyna
	// z ImageMagickiem i bez Real-ESRGAN-a ma stracić powiększanie, a nie retusz.
	NarzedziaObrazuModelu NarzedziaObrazuModelu
	// Doradcy obsługuje `advisor.consult` — konsultację modelu u doradcy
	// w trakcie tury. Port osobny, nie rozszerzenie portu Rozmowa: tamten
	// prowadzi turę okna, a ten pyta o radę poza nią i niczego w stanie rdzenia
	// nie zmienia (konsultacja, nie delegacja).
	Doradcy Doradcy
	// NarzedziaMedia obsługuje rodzinę `media.*` — pomiar i przetworzenie
	// dźwięku oraz filmu binariami `ffprobe`/`ffmpeg`, wołanymi jedyną drogą
	// arsenału (`server/internal/zewnetrzne`). Port osobny z tego samego powodu
	// co NarzedziaObrazu wyżej: to są czynności modelu wołane z rozmowy, a nie
	// okno Operatora — wspólny jest wyłącznie magazyn zasobów.
	NarzedziaMedia NarzedziaMedia
	// Dokumenty obsługuje rodzinę `document.*` — dwa narzędzia modelu do pracy
	// na dokumencie: zamianę formatu i odczyt treści (PDF, skan, zdjęcie
	// kartki). Port osobny z tego samego powodu co dwa wyżej: to są czynności
	// modelu wołane z rozmowy, a nie okno Operatora.
	Dokumenty Dokumenty
	// NarzedziaArchiwum obsługuje rodzinę `archive.*` — złożenie pracy w jeden
	// plik i otwarcie archiwum przysłanego przez Operatora. Port osobny z tego
	// samego powodu co trzy wyżej: to są czynności modelu wołane z rozmowy,
	// a nie okno Operatora. Wytworzone archiwum ląduje w magazynie zasobów
	// Designu — drugiego magazynu bajtów nie ma.
	NarzedziaArchiwum NarzedziaArchiwum
	// Poczta obsługuje rodzinę `mail.*` — dziesięć komend, którymi rdzeń jest
	// klientem skrzynki Operatora i nikim więcej; własnego serwera poczty
	// platforma nie ma. Port osobny, nie rozszerzenie portu Asystent: poczta jest
	// zasobem urządzenia, po który sięga model w dowolnym oknie, tak samo jak po
	// narzędzia obrazu czy dokumentów. Załączniki listu lądują w magazynie
	// zasobów Designu — drugiego magazynu bajtów nie ma.
	Poczta Poczta
	// NarzedziaSesji obsługuje rodzinę `session.tool.*` oraz `tools.catalog.list`
	// — doraźne dołożenie narzędzia na czas sesji, zdjęcie go, wykaz dołożeń
	// i wykaz pozycji po ukośniku, z którego się je wybiera. Port osobny, nie
	// rozszerzenie portu Sesje: tamten prowadzi sesję, a ten prowadzi zestaw
	// narzędzi modelu w jej obrębie. Ten sam adapter wnosi dołożenia do zestawu
	// narzędzi tury (montaz_rozmowa.go) — drugi byłby drugą prawdą o tym, czym
	// model w tej sesji dysponuje.
	NarzedziaSesji NarzedziaSesji
	Nadajnik       Nadajnik
	Nasluch        Nasluch
	Dziennik       *log.Logger
}

// Zloz składa rdzeń: zakłada rejestr, wpina obsługiwacze wszystkich domen
// kontraktu i buduje rozpoznanie nazw z tego, co rzeczywiście zostało wpięte.
//
// To jedyne miejsce, w którym powstaje mapa komend rdzenia. Kolejność wpinania
// nie ma znaczenia — rejestr jest mapą nazwa→obsługiwacz, a nie łańcuchem
// warunków.
func Zloz(p Porty) *Rdzen {
	rejestr := NowyRejestr()
	nadawca := nowyEmiter(p.Nadajnik)
	// Więź połączenie → sesja bramki. Jedna na rdzeń: powitanie ją zakłada,
	// bramka ją nadpisuje przy wejściu, zmiana hasła ją czyta (`wiez_polaczenia.go`).
	wiez := nowaWiezBramki()
	// Deklaracje zdolności klientów spisane w powitaniu. Jedna mapa na rdzeń:
	// czyta ją `launcher.hotkey.*`, żeby wiedzieć, czy skrót globalny ma kto
	// przechwycić (`adapter_wywolywacz.go`). Mapę zakłada montaż i podaje ją
	// tutaj wraz z portami; rdzeń złożony bez montażu (na przykład w teście)
	// zakłada własną, żeby powitanie nie sięgało po nic.
	zdolnosci := p.zdolnosciKlienta
	if zdolnosci == nil {
		zdolnosci = noweZdolnosciKlientow()
	}

	// Nastawy poziomu `aplikacja` czyta ten sam adapter, co całą konfigurację —
	// rozszerzenie portu, nie drugi port. Adapter nieznający
	// rozszerzenia daje port pusty i powitanie milczy o wymogu.
	nastawy, _ := p.Ustawienia.(NastawyAplikacji)
	zarejestrujPolaczenie(rejestr, p.Uwierzytelnianie, wiez, nastawy, zdolnosci)
	zarejestrujSesje(rejestr, p.Sesje, nadawca)
	zarejestrujOkna(rejestr, p.Okna, nadawca)
	zarejestrujKanalModelu(rejestr, p.Okna.(KanalModelu), nadawca)
	zarejestrujRozmowe(rejestr, p.Rozmowa, nadawca)
	zarejestrujUstawienia(rejestr, p.Ustawienia, nadawca)
	zarejestrujProwenancjeKonfiguracji(rejestr, p.Ustawienia.(ProwenancjaKonfiguracji))
	zarejestrujKanaly(rejestr, p.Kanaly)
	zarejestrujKolejki(rejestr, p.Kolejki, nadawca)
	zarejestrujWiazaniaKolejek(rejestr, p.Kolejki, nadawca)
	zarejestrujZleceniaKolejek(rejestr, p.Kolejki, nadawca)
	zarejestrujPrzenoszenie(rejestr, p.Przenoszenie, nadawca)
	zarejestrujNawigacje(rejestr, p.Nawigacja, nadawca)
	zarejestrujMonitor(rejestr, p.Nawigacja)
	zarejestrujWiazanieSesji(rejestr, p.WiazanieSesji, nadawca)
	zarejestrujAkcje(rejestr, p.Akcje, shared.CommandActionList)
	zarejestrujCentrumPowiadomien(rejestr, p.CentrumPowiadomien)
	zarejestrujKatalogUstawien(rejestr, p.KatalogUstawien)
	zarejestrujPunktyDostepu(rejestr, p.PunktyDostepu, nadawca)
	zarejestrujNadaniaDostepu(rejestr, p.NadaniaDostepu, nadawca)
	zarejestrujKonta(rejestr, p.Konta, nadawca)
	zarejestrujTozsamosc(rejestr, p.Tozsamosc, nadawca)
	zarejestrujAgentow(rejestr, p.Agenci, nadawca)
	zarejestrujWersjeEksperta(rejestr, p.WersjeEksperta, nadawca)
	zarejestrujZakresEksperta(rejestr, p.ZakresEksperta, nadawca)
	zarejestrujZakresyNarzedzi(rejestr, p.ZakresyNarzedzi)
	zarejestrujZespoly(rejestr, p.Zespoly, nadawca)
	zarejestrujPrzestrzenRobocza(rejestr, p.PrzestrzenRobocza, nadawca)
	// Moduł Workspace wchodzi trzema portami ponad `PrzestrzenRobocza`: pamięć
	// projektu, hub planowania i obszar wiedzy. Wszystkie wypełnia ten sam
	// adapter, ale rodzina liczy czterdzieści komend i jeden port byłby wykazem
	// wszystkiego, co moduł umie, zamiast wykazem bytu.
	//
	// Rzutowanie idzie w postaci DWUWARTOŚCIOWEJ z rozmysłem. Rzutowanie bez
	// sprawdzenia zamieniało brak jednej metody w panikę przy składaniu rdzenia,
	// czyli kładło cały produkt z powodu jednej rodziny komend. Rdzeń złożony
	// bez portu ma pracować dalej w pozostałych czynnościach, a brak widać
	// w wykazie komend powitania — tak samo jak przy porcie zerowym, który
	// każda funkcja `zarejestruj*` przyjmuje bez awarii.
	if pamiec, jest := p.PrzestrzenRobocza.(PamiecPrzestrzeni); jest {
		zarejestrujPamiec(rejestr, pamiec, nadawca)
	}
	if planowanie, jest := p.PrzestrzenRobocza.(PlanowanieProjektu); jest {
		zarejestrujPlanowanieProjektu(rejestr, planowanie, nadawca)
	}
	if wiedza, jest := p.PrzestrzenRobocza.(WiedzaProjektu); jest {
		zarejestrujWiedzeProjektu(rejestr, wiedza, nadawca)
	}
	zarejestrujAutomatyki(rejestr, p.Automatyki, nadawca)
	zarejestrujDobudoweAutomatyk(rejestr, p.Automatyki, nadawca)
	zarejestrujHarmonogramy(rejestr, p.Automatyki.(Harmonogramy), nadawca)
	zarejestrujOrkiestracje(rejestr, p.Automatyki.(Orkiestracja), nadawca)
	zarejestrujTerminal(rejestr, p.Terminal, nadawca)
	zarejestrujWyjscieTerminala(rejestr, p.Terminal.(WyjscieTerminala))
	// Wyposażenie modułu — karty, pliki, hosty, skrypty, klucze, tunele
	// i obserwacje — wchodzi rzutowaniem DWUWARTOŚCIOWYM z tego samego powodu co
	// porty Workspace'u wyżej: rdzeń złożony z portem niepełnym ma pracować dalej
	// w pozostałych czynnościach, a brak widać w wykazie komend powitania.
	if wyposazenie, jest := p.Terminal.(WyposazenieTerminala); jest {
		zarejestrujWyposazenieTerminala(rejestr, wyposazenie)
	}
	zarejestrujDevelopera(rejestr, p.Developer, nadawca)
	zarejestrujDiagnostyke(rejestr, p.Diagnostyka, nadawca)
	zarejestrujProwenancje(rejestr, p.Prowenancja)
	zarejestrujKondycje(rejestr, p.Kondycja)
	zarejestrujAlerty(rejestr, p.Alerty)
	zarejestrujSchowek(rejestr, p.Schowek)
	zarejestrujSkrotyTekstowe(rejestr, p.SkrotyTekstowe)
	zarejestrujKontekstyPamieci(rejestr, p.KontekstyPamieci)
	zarejestrujWywolywacz(rejestr, p.Wywolywacz)
	zarejestrujZajetoscKontekstu(rejestr, p.ZajetoscKontekstu)
	zarejestrujZuzycie(rejestr, p.Zuzycie)
	zarejestrujWarsztatPdf(rejestr, p.WarsztatPdf)
	zarejestrujBezpieczenstwoDokumentu(rejestr, p.BezpieczenstwoDokumentu)
	zarejestrujBiblioteke(rejestr, p.Biblioteka, nadawca)
	zarejestrujStudio(rejestr, p.Studio, nadawca)
	// Pętla wykonawcza Studia (`studio.plan.*`) stoi na tym samym porcie, ale
	// ma stan własny — magazyn rozkładów — więc wchodzi osobnym wpięciem.
	// Rejestr komend jedzie jej argumentem, bo nastawy pętli czyta ona
	// obsługiwaczem `studio.agents.settings.get`, a nie własnym odczytem
	// konfiguracji: drugiego magazynu nastaw w tym module nie ma.
	zarejestrujPetleWykonawczaStudia(rejestr, p.Studio, nadawca)
	// Zapora blokad fragmentu owija to, co w rejestrze JUŻ stoi, więc jedzie po
	// obu wpięciach Studia — nie przed nimi. Owinięcia nie da się założyć na
	// komendę, której w rejestrze jeszcze nie ma, a bez owinięcia zapora stałaby
	// napisana i nieaktywna: blokada fragmentu byłaby wtedy zapisem bez skutku.
	// Port jest interfejsem, a zapora potrzebuje adaptera, bo sprawdza blokadę
	// jego drogą do warstwy danych; wpięcie zachodzi tylko wtedy, gdy port jest
	// tym właśnie adapterem.
	if adapter, jest := p.Studio.(*adapterStudia); jest {
		zaporaBlokadStudia(rejestr, adapter)
		// Siatka śladu autora idzie PO zaporze, więc owija ją z zewnątrz. Ta
		// kolejność jest zamierzona: czynność zatrzymana blokadą nie ma zostawić
		// śladu w dzienniku, bo nic nie zrobiła — a gdyby siatka stała na
		// zewnątrz, mierzyłaby treść sprzed i po odmowie i słusznie nie znalazła
		// różnicy. Odwrotna kolejność dawałaby ten sam skutek dopiero przez
		// przypadek; tu wynika z układu.
		sladWykonawcyStudia(rejestr, adapter)
	}
	zarejestrujPrzegladarke(rejestr, p.Przegladarka, nadawca)
	zarejestrujDesign(rejestr, p.Design, nadawca)
	zarejestrujBadania(rejestr, p.Badania, nadawca)
	zarejestrujAsystenta(rejestr, p.Asystent, nadawca)
	zarejestrujAplikacje(rejestr, p.Aplikacje, nadawca)
	zarejestrujKomponenty(rejestr, p.Komponenty, nadawca)
	zarejestrujPrzekazanieOkna(rejestr, p.PrzekazanieOkna, nadawca)
	zarejestrujTlumaczenie(rejestr, p.Tlumaczenie, nadawca)
	zarejestrujDebate(rejestr, p.Debata, nadawca)
	zarejestrujIzolacje(rejestr, p.Izolacja)
	zarejestrujUwierzytelnianie(rejestr, p.Uwierzytelnianie, nadawca, wiez)
	zarejestrujUrzadzenia(rejestr, p.Urzadzenia, nadawca, wiez)
	zarejestrujNakladkeAod(rejestr, p.NakladkaAod, nadawca)
	// Dwa wpięcia bez nowego portu: obie funkcje stoją na portach, które rdzeń
	// już niesie.
	zarejestrujHistorieSesji(rejestr, p.Sesje, nadawca)
	zarejestrujWarstweMobilna(rejestr, p.Nawigacja)
	zarejestrujRozszerzenia(rejestr, p.Rozszerzenia)
	zarejestrujRole(rejestr, p.Role, nadawca)
	zarejestrujWykazRol(rejestr, p.Role, nadawca)
	zarejestrujSekcjePaneli(rejestr, p.Panele)
	zarejestrujHistorie(rejestr, p.Historia, nadawca)
	zarejestrujMowe(rejestr, p.Mowa)
	zarejestrujWiedze(rejestr, p.Wiedza)
	zarejestrujOrkiestracjePodagentow(rejestr, p.Podagenci)
	zarejestrujDoradcow(rejestr, p.Doradcy, nadawca)
	zarejestrujDokumenty(rejestr, p.Dokumenty)
	zarejestrujNarzedziaMediow(rejestr, p.NarzedziaMedia)
	zarejestrujNarzedziaObrazu(rejestr, p.NarzedziaObrazu)
	zarejestrujNarzedziaObrazuModelu(rejestr, p.NarzedziaObrazuModelu)
	zarejestrujNarzedziaArchiwum(rejestr, p.NarzedziaArchiwum)
	zarejestrujPoczte(rejestr, p.Poczta, nadawca)
	zarejestrujNarzedziaSesji(rejestr, p.NarzedziaSesji, nadawca)

	return &Rdzen{
		rejestr:  rejestr,
		komendy:  rejestr.RejestrKomend(),
		nasluch:  p.Nasluch,
		dziennik: p.Dziennik,
		wiez:     wiez,
	}
}
