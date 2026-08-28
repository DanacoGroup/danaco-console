# Danaco Console — Instrukcja użytkowania

| | |
|---|---|
| **Produkt** | Danaco Console |
| **Rodzaj** | Platforma AI Workspace OS |
| **Producent** | Danaco Holding Group Sp. z o.o. |
| **Twórca** | Dariusz Naharnowicz |
| **Wersja** | v2.0 |
| **Status** | Deweloperski |
| **Data** | 2026-08-20 |

**Informacje szczegółowe dokumentu:**

| | |
|---|---|
| **Tytuł** | Instrukcja użytkowania Danaco Console |
| **Klasa dokumentu** | Stan wdrożenia |
| **Odbiorcy** | Operator — jedyna rola, którą produkt dziś zna (rozdz. 1.1) |
| **Przeznaczenie** | Prowadzi Operatora przez codzienną pracę z aplikacją w stanie faktycznym bieżącego wydania: uruchomienie, nawigację, sesje i okna komunikacji, kanały modelu, konfigurację okna oraz rozwiązywanie trudności. Na jego podstawie Operator wie, co w produkcie działa, a co jest wyłącznie pozycją nawigacji. |
| **Zakres** | uruchamianie aplikacji, struktura interfejsu, nawigacja, sesje i okna komunikacji, modele, konta, dostawcy i kanały, konfiguracja sesji i okna, środowiska, praca z modułami, obsługa narzędzi, praca z plikami i artefaktami, historia i dane, ustawienia, typowe scenariusze pracy, komunikaty i diagnostyka, zakończenie i wznowienie pracy, zbiorcze zestawienie stanu funkcji |
| **Poza zakresem** | budowa i architektura produktu — [Opis produktu](README.md); instalacja, wdrożenie i pełny wykaz kluczy konfiguracji — [Instalacja i konfiguracja](INSTALACJA-I-KONFIGURACJA.md); warunki korzystania — [Licencja produktu](LICENSE.md) |
| **Dokument nadrzędny** | [Spis opracowań](SPIS-OPRACOWAN.md) |
| **Dokumenty powiązane** | [Opis produktu](README.md) · [Instalacja i konfiguracja](INSTALACJA-I-KONFIGURACJA.md) · [Licencja produktu](LICENSE.md) · [Standard redakcyjny i językowy](STANDARD-REDAKCYJNY-I-JEZYKOWY.md) |
| **Prototypy odniesienia** | nie dotyczy — dokument klasy Stan wdrożenia opisuje interfejs zbudowany, nie projektowany; prototypy przywołują opracowania klasy Specyfikacja docelowa |
| **Źródła normatywne** | `budowa/shared/contract.json` · `budowa/client/src/` · `budowa/server/internal/store/` · `budowa/scripts/` |
| **Zasada nadrzędna** | Instrukcja opisuje wyłącznie to, co Operator może dziś wykonać; funkcja nieosiągalna z interfejsu jest oznaczona, nie przemilczana. |

---

## Spis treści

1. [Wprowadzenie](#1-wprowadzenie)
   - [1.1 Do kogo skierowany jest ten dokument](#11-do-kogo-skierowany-jest-ten-dokument)
   - [1.2 Czym jest Danaco Console](#12-czym-jest-danaco-console)
   - [1.3 Co oznacza status „Deweloperski”](#13-co-oznacza-status-deweloperski)
   - [1.4 Terminologia obowiązująca](#14-terminologia-obowiązująca)
   - [1.5 Oznaczenia stanu używane w tym dokumencie](#15-oznaczenia-stanu-używane-w-tym-dokumencie)
   - [1.6 Miejsce tego dokumentu wśród dokumentów produktu](#16-miejsce-tego-dokumentu-wśród-dokumentów-produktu)
   - [1.7 Podstawa opisu — stan kodu, na którym oparto instrukcję](#17-podstawa-opisu--stan-kodu-na-którym-oparto-instrukcję)
2. [Uruchamianie aplikacji](#2-uruchamianie-aplikacji)
   - [2.1 Warunki wstępne](#21-warunki-wstępne)
   - [2.2 Uruchomienie przez powłokę Tauri w drzewie deweloperskim](#22-uruchomienie-przez-powłokę-tauri-w-drzewie-deweloperskim)
   - [2.3 Uruchomienie samego rdzenia](#23-uruchomienie-samego-rdzenia)
   - [2.4 Uruchomienie interfejsu z serwera rozwojowego](#24-uruchomienie-interfejsu-z-serwera-rozwojowego)
   - [2.5 Wydanie produktu do `C:\DanacoConsole_App`](#25-wydanie-produktu-do-cdanacoconsole_app)
   - [2.6 Zmienne środowiska i przełączniki wiersza poleceń](#26-zmienne-środowiska-i-przełączniki-wiersza-poleceń)
   - [2.7 Doprowadzenie świeżej instalacji do stanu zdolnego do rozmowy](#27-doprowadzenie-świeżej-instalacji-do-stanu-zdolnego-do-rozmowy)
3. [Struktura interfejsu](#3-struktura-interfejsu)
   - [3.1 Trzy widoki najwyższego rzędu](#31-trzy-widoki-najwyższego-rzędu)
   - [3.2 Centrum dowodzenia — strona główna](#32-centrum-dowodzenia--strona-główna)
   - [3.3 Powłoka środowiska — cztery pasy](#33-powłoka-środowiska--cztery-pasy)
   - [3.4 Mission Control](#34-mission-control)
   - [3.5 Motyw jasny i ciemny](#35-motyw-jasny-i-ciemny)
   - [3.6 Zasada zera blokad i dymki komunikatów](#36-zasada-zera-blokad-i-dymki-komunikatów)
4. [Nawigacja](#4-nawigacja)
   - [4.1 Wejście do środowiska](#41-wejście-do-środowiska)
   - [4.2 Wybór modułu — co robi, a czego nie robi](#42-wybór-modułu--co-robi-a-czego-nie-robi)
   - [4.3 Powrót i przełączanie widoków](#43-powrót-i-przełączanie-widoków)
   - [4.4 Adres dokumentu i odświeżenie](#44-adres-dokumentu-i-odświeżenie)
5. [Sesje i okna komunikacji](#5-sesje-i-okna-komunikacji)
   - [5.1 Pojęcia: sesja, karta sesji, okno komunikacji](#51-pojęcia-sesja-karta-sesji-okno-komunikacji)
   - [5.2 Powstanie pierwszego okna](#52-powstanie-pierwszego-okna)
   - [5.3 Scena jednego, dwóch i trzech okien](#53-scena-jednego-dwóch-i-trzech-okien)
   - [5.4 Role okien](#54-role-okien)
   - [5.5 Więź koordynator–wykonawca i przekazanie zlecenia](#55-więź-koordynatorwykonawca-i-przekazanie-zlecenia)
   - [5.6 Prowadzenie rozmowy](#56-prowadzenie-rozmowy)
   - [5.7 Co pokazuje wpis rozmowy](#57-co-pokazuje-wpis-rozmowy)
   - [5.8 Brak pamięci między turami — ograniczenie krytyczne](#58-brak-pamięci-między-turami--ograniczenie-krytyczne)
6. [Modele, konta, dostawcy i kanały](#6-modele-konta-dostawcy-i-kanały)
   - [6.1 Kanał modelu — pojęcie i rodzaje](#61-kanał-modelu--pojęcie-i-rodzaje)
   - [6.2 Rejestr kanałów i brak interfejsu do ich zakładania](#62-rejestr-kanałów-i-brak-interfejsu-do-ich-zakładania)
   - [6.3 Konta i pula rotacji](#63-konta-i-pula-rotacji)
   - [6.4 Co realnie steruje wywołaniem modelu](#64-co-realnie-steruje-wywołaniem-modelu)
7. [Środowiska](#7-środowiska)
8. [Konfiguracja sesji i okna — panel sterowania](#8-konfiguracja-sesji-i-okna--panel-sterowania)
   - [8.1 Umiejscowienie i obsługa szuflady](#81-umiejscowienie-i-obsługa-szuflady)
   - [8.2 Dziewięć kontrolek](#82-dziewięć-kontrolek)
   - [8.3 Zestawienie skuteczności ustawień okna](#83-zestawienie-skuteczności-ustawień-okna)
9. [Praca z modułami](#9-praca-z-modułami)
10. [Obsługa narzędzi — mosty MCP, katalogi, nadania](#10-obsługa-narzędzi--mosty-mcp-katalogi-nadania)
   - [10.1 Okno „Dostępy i katalog roboczy”](#101-okno-dostępy-i-katalog-roboczy)
   - [10.2 Punkty dostępu](#102-punkty-dostępu)
   - [10.3 Nadania dostępu per okno](#103-nadania-dostępu-per-okno)
   - [10.4 Katalog roboczy modelu](#104-katalog-roboczy-modelu)
   - [10.5 Narzędzia platformy dla modelu](#105-narzędzia-platformy-dla-modelu)
11. [Praca z plikami i artefaktami](#11-praca-z-plikami-i-artefaktami)
   - [11.1 Jak model uzyskuje dostęp do plików](#111-jak-model-uzyskuje-dostęp-do-plików)
   - [11.2 Gdzie powstają pliki modelu](#112-gdzie-powstają-pliki-modelu)
   - [11.3 Artefakty i treści obszerne](#113-artefakty-i-treści-obszerne)
12. [Historia i dane](#12-historia-i-dane)
   - [12.1 Co trafia do bazy](#121-co-trafia-do-bazy)
   - [12.2 Czego klient nie odtwarza](#122-czego-klient-nie-odtwarza)
   - [12.3 Co przeżywa restart rdzenia](#123-co-przeżywa-restart-rdzenia)
13. [Ustawienia](#13-ustawienia)
   - [13.1 Okno konfiguracji](#131-okno-konfiguracji)
   - [13.2 Okno „Modele, konta i tożsamość”](#132-okno-modele-konta-i-tożsamość)
   - [13.3 Okno dostępów](#133-okno-dostępów)
   - [13.4 Pozycje listwy bez własnego ekranu](#134-pozycje-listwy-bez-własnego-ekranu)
14. [Typowe scenariusze pracy](#14-typowe-scenariusze-pracy)
   - [14.1 Od uruchomienia do pierwszej rozmowy](#141-od-uruchomienia-do-pierwszej-rozmowy)
   - [14.2 Praca z tożsamością modelu](#142-praca-z-tożsamością-modelu)
   - [14.3 Praca z wieloma oknami](#143-praca-z-wieloma-oknami)
   - [14.4 Nadanie modelowi dostępu do maszyny](#144-nadanie-modelowi-dostępu-do-maszyny)
15. [Komunikaty i diagnostyka](#15-komunikaty-i-diagnostyka)
   - [15.1 Najczęstsze komunikaty i ich przyczyny](#151-najczęstsze-komunikaty-i-ich-przyczyny)
   - [15.2 Gdzie szukać przyczyn](#152-gdzie-szukać-przyczyn)
   - [15.3 Prowenancja jako narzędzie diagnostyczne](#153-prowenancja-jako-narzędzie-diagnostyczne)
16. [Zakończenie i wznowienie pracy](#16-zakończenie-i-wznowienie-pracy)
   - [16.1 Zamknięcie okna i zakończenie powłoki](#161-zamknięcie-okna-i-zakończenie-powłoki)
   - [16.2 Zatrzymanie rdzenia](#162-zatrzymanie-rdzenia)
   - [16.3 Wznowienie pracy](#163-wznowienie-pracy)
   - [16.4 Powrót do sesji trwającej na rdzeniu](#164-powrót-do-sesji-trwającej-na-rdzeniu)
17. [Pozostałe funkcje wymagające objaśnienia](#17-pozostałe-funkcje-wymagające-objaśnienia)
   - [17.1 Obsługa klawiaturą](#171-obsługa-klawiaturą)
   - [17.2 Odporność na nazwy nieznane](#172-odporność-na-nazwy-nieznane)
   - [17.3 Łączność i kolejka wychodząca](#173-łączność-i-kolejka-wychodząca)
   - [17.4 Adresowanie rozgłoszeń](#174-adresowanie-rozgłoszeń)
   - [17.5 Wersjonowanie kontraktu](#175-wersjonowanie-kontraktu)
18. [Zbiorcze zestawienie stanu funkcji](#18-zbiorcze-zestawienie-stanu-funkcji)
19. [Kwestie pozostawione do decyzji Operatora](#19-kwestie-pozostawione-do-decyzji-operatora)

---

## 1. Wprowadzenie

### 1.1 Do kogo skierowany jest ten dokument

Instrukcja opisuje obsługę Danaco Console przez **Operatora** — jedyną rolę człowieka, jaką zna terminologia obowiązująca produktu.

**Operator** jest osobą fizyczną faktycznie obsługującą produkt: konfigurującą go **oraz** prowadzącą sesje pracy z modelem, odpowiadającą za polecenia wydawane modelowi i za wykorzystanie wyników jego pracy. Definicja ta jest wspólna dla całej dokumentacji produktu — wynika z [LICENSE.md](LICENSE.md) (rozdział 2.1, „Podmioty”) i z [README.md](README.md) (rozdział 2.3), gdzie zastrzeżono wprost, że Operator **nie jest „użytkownikiem” w sensie konsumenta interfejsu**. Operator zakłada i konfiguruje kanały modelu, konta oraz punkty dostępu, ustala tożsamość modeli, decyduje o trybie uprawnień i katalogach roboczych okna, nadzoruje pracę biegnącą w tle, a także podejmuje decyzje architektoniczne i produktowe, których wykonawca — człowiek albo model — podjąć nie może.

Produkt nie rozdziela dziś praw ani widoków między dwie odrębne role: nie ma logowania, kont użytkowników platformy ani żadnego mechanizmu, który odbierałby komukolwiek dostęp do konfiguracji. Wszystko, co opisuje ta instrukcja, dostępne jest z jednego interfejsu i dla jednej roli. Dlatego dokument mówi wyłącznie o Operatorze.

Tam, gdzie w tekście pojawia się zwrot **„praca w oknie komunikacji”**, chodzi o tę część czynności Operatora, która nie wymaga wcześniejszej konfiguracji instalacji — prowadzenie rozmowy z modelem w gotowym środowisku. Nie jest to odrębna rola, lecz odrębny **rodzaj czynności**. Rozróżnienie ma znaczenie praktyczne: rozdziały [2](#2-uruchamianie-aplikacji), [6](#6-modele-konta-dostawcy-i-kanały) i [13](#13-ustawienia) opisują czynności konfiguracyjne, a rozdziały [5](#5-sesje-i-okna-komunikacji), [11](#11-praca-z-plikami-i-artefaktami) i [14](#14-typowe-scenariusze-pracy) — czynności prowadzenia pracy.

Dokument nie wymaga znajomości kodu źródłowego. Każda opisana czynność interfejsu została zweryfikowana w kodzie klienta; opisane są wyłącznie elementy istniejące i obsłużone. Tam, gdzie funkcja jest zapowiedziana koncepcją, lecz nie została zaimplementowana, dokument mówi to wprost — nie w przypisie, lecz w miejscu, w którym Operator jej szuka.

### 1.2 Czym jest Danaco Console

Danaco Console jest operacyjnym środowiskiem współpracy człowieka i modeli językowych. Koncepcja produktu przewiduje cztery wyspecjalizowane środowiska pracy — **TalkIn** (wiedza), **WorkSpace** (praca), **CodeStudio** (technologia) oraz **MultitaskingAI** (orkiestracja) — połączone wspólną warstwą komunikacji, konfiguracji, dostępów i tożsamości modelu.

Produkt składa się z trzech warstw wykonawczych:

| Warstwa | Nazwa w dokumentacji | Technologia | Rola |
|---|---|---|---|
| Serwer | **rdzeń** | Go | trwałość, rejestr kanałów, uruchamianie procesów modelu, transport |
| Aplikacja okienkowa | **powłoka** | Tauri (Rust) | okno natywne, ikona w zasobniku, uruchomienie rdzenia w tle |
| Warstwa widoku | **klient** / **interfejs** | TypeScript | wszystko, co Operator widzi i czym steruje |

Klient i rdzeń rozmawiają jednym kanałem WebSocket, w formacie JSON, kopertą `{type, sessionId, payload}`. Nazwy komend i zdarzeń pochodzą z jednego kontraktu (katalog `shared/`), wspólnego dla obu stron.

Zwrot `{type, sessionId, payload}` jest skrótem terminologicznym nazywającym trzy pola
najczęściej przywoływane w tym dokumencie; pełny opis koperty w kontrakcie jest szerszy
i warto go znać, ponieważ rozdział [15](#15-komunikaty-i-diagnostyka) odwołuje się do
zachowań strumienia. Kontrakt wymaga obowiązkowo trzech pól: **`type`** (typ komunikatu
w notacji `obszar.zasób.akcja`), **`id`** (identyfikator zadania, powtarzany w odpowiedzi
i w każdym fragmencie strumienia) oraz **`timestamp`** (czas nadania w milisekundach
epoki). Pola **`sessionId`** i **`payload`** są opcjonalne — `sessionId` pozostaje pusty
w powitaniu połączenia, a `payload` niosą wyłącznie komunikaty mające treść własną.
Odpowiedź dokłada dwa pola: **`status`** oraz **`error`** (to drugie wyłącznie przy
statusie błędu). Fragment strumienia dokłada natomiast **`seq`** — numer fragmentu
liczony od jedynki — oraz **`done`**, prawdziwe w fragmencie ostatnim. Dzięki polu `id`
odpowiedź i wszystkie fragmenty da się jednoznacznie przypisać do zadania, które je
wywołało, a `seq` wraz z `done` pozwalają rozpoznać, czy strumień domknął się poprawnie,
czy urwał.

### 1.3 Co oznacza status „Deweloperski”

Wersja v2.0 ma status **Deweloperski** i tak należy ją traktować w praktyce eksploatacyjnej. Konkretnie oznacza to trzy rzeczy, które Operator odczuje od pierwszej minuty:

1. **Nie istnieje instalator ani procedura aktualizacji.** Repozytorium ma wprawdzie skrypt wydania, który składa gotowy do uruchomienia produkt w katalogu `C:\DanacoConsole_App` (rozdział [2.5](#25-wydanie-produktu-do-cdanacoconsole_app)) — a niniejszy dokument leży w wyniku takiego wydania, obok wytworzonych nim plików wykonywalnych. Wydanie **nie jest jednak instalacją**: nie tworzy wpisu w rejestrze systemu, nie zakłada skrótów, nie sprawdza zależności zewnętrznych i nie ma drogi odinstalowania. Zbudowanie instalatora NSIS pozostaje osobnym, ręcznym krokiem, wymagającym narzędzia NSIS.
2. **Świeża instalacja nie wykona ani jednej tury modelu bez ręcznego doprowadzenia.** Rejestr kanałów modelu startuje pusty, katalog kont startuje pusty, a interfejs nie zawiera ekranu do założenia kanału. Procedurę doprowadzenia opisuje rozdział [2.7](#27-doprowadzenie-świeżej-instalacji-do-stanu-zdolnego-do-rozmowy).
3. **Model nie pamięta poprzednich tur okna.** Każda wypowiedź jest osobnym wywołaniem bez historii. Jest to najważniejsze ograniczenie tej wersji i opisuje je rozdział [5.8](#58-brak-pamięci-między-turami--ograniczenie-krytyczne).

Poza tymi trzema punktami działa jednak pełna, sprawna droga wykonawcza: wypowiedź Operatora trafia do procesu modelu, odpowiedź wraca strumieniem, rozmowa zapisuje się w bazie, a tożsamość modelu, tryb uprawnień, katalogi robocze i mosty MCP realnie zmieniają sposób uruchomienia procesu.

### 1.4 Terminologia obowiązująca

| Pojęcie | Znaczenie w produkcie |
|---|---|
| **Operator** | osoba fizyczna faktycznie obsługująca produkt: konfigurująca go **oraz** prowadząca sesje pracy z modelem; odpowiada za polecenia wydawane modelowi i za wykorzystanie wyników jego pracy (zgodnie z [LICENSE.md](LICENSE.md) rozdz. 2.1 i [README.md](README.md) rozdz. 2.3). Nie jest „użytkownikiem” w sensie konsumenta interfejsu; produkt nie zna innej roli człowieka |
| **rdzeń** | proces serwera (Go), właściciel bazy i procesów modelu |
| **powłoka** | aplikacja natywna Tauri: okno, zasobnik, uruchomienie rdzenia |
| **klient**, **interfejs** | warstwa widoku w TypeScript |
| **środowisko** | jedno z czterech: TalkIn, WorkSpace, CodeStudio, MultitaskingAI |
| **moduł** | pozycja nawigacji bocznej środowiska (Studio, Workspace, Browser, …) |
| **okno operacyjne** | ekran roboczy modułu; w v2.0 istnieje wyłącznie okno komunikacji |
| **okno komunikacji** | okno rozmowy z modelem — jednostka, do której należy kanał, tryb uprawnień, katalogi i rola |
| **sesja** | byt nadrzędny wobec okien; powstaje przy uzgodnieniu z rdzeniem |
| **karta sesji** | zakładka w drugim pasie powłoki środowiska |
| **kanał modelu** | wiersz rejestru wskazujący, przez co uruchamiany jest model; rodzaje kontraktowe: `cli`, `api`, `sdk`, `lokalny` |
| **adapter kanału** | kod budujący kanał z wiersza rejestru; wskazywany parametrem `adapter` wiersza, a w jego braku rodzajem kanału. Adaptery wbudowane: `echo`, `api`, `cli` |
| **konto** | nośnik tożsamości technicznej wobec dostawcy; rodzaje: `cli`, `api`, `sdk` |
| **pula kont** | zbiór kont rodzaju `cli` podlegających rotacji po wyczerpaniu limitu |
| **nakładka tożsamości** | trójwarstwowy prompt systemowy: konstytucja → profil roli → ekspertyza |
| **prowenancja** | pełny opis tego, co poszło do modelu: program, argumenty, prompt, konto, katalogi |
| **kontrakt** | katalog `shared/` — jedyne źródło nazw komend i zdarzeń |
| **koperta** | struktura ramki komunikatu: `{type, sessionId, payload}`; pełny zestaw pól kontraktu opisuje rozdział [1.2](#12-czym-jest-danaco-console) |
| **Decyzje architektoniczne** | rozstrzygnięcia przyjęte dla budowy platformy, obowiązujące w kodzie i w tej instrukcji |

### 1.5 Oznaczenia stanu używane w tym dokumencie

Aby uniknąć wrażenia, że opisywana jest funkcja niedziałająca, każdy element o niepełnym stanie jest oznaczony jednym z sześciu zwrotów — zestaw ujednolicony z tym, którym posługuje się [README.md](README.md):

- **DZIAŁA** — element wykonuje to, co zapowiada jego nazwa, na całej drodze od interfejsu do skutku.
- **DZIAŁA CZĘŚCIOWO** — element wykonuje część zapowiadanego skutku; zakres podany jest w miejscu opisu.
- **NIEZINTEGROWANE** — kod istnieje i jest poprawny po jednej albo obu stronach, lecz nie ma konsumenta — funkcja nie jest osiągalna dla Operatora; należy tu również przypadek, w którym interfejs przyjmuje wartość i zapisuje ją, lecz wykonanie jej nie odczytuje.
- **ATRAPA** — element widoczny w interfejsie (pozycja menu, kafel, przycisk), za którym nie stoi realizacja.
- **BRAK** — przewidziane koncepcją, w bieżącej wersji niezaimplementowane; nazwa występuje wyłącznie w nawigacji albo w dokumentacji koncepcyjnej.

Osobno oznaczane są kwestie wymagające rozstrzygnięcia poza dokumentacją: **[DO DECYZJI OPERATORA]**.

### 1.6 Miejsce tego dokumentu wśród dokumentów produktu

Katalog produktu zawiera cztery dokumenty. Ich zakresy zachodzą na siebie w jednym miejscu i Operator musi wiedzieć, który z nich rozstrzyga.

| Dokument | Zakres deklarowany w [README.md](README.md) (rozdz. 17 i 19) |
|---|---|
| **[README.md](README.md)** | karta techniczno-produktowa: koncepcja, model katalogów, stos technologiczny, decyzje, zestawienie znanych ograniczeń |
| **[LICENSE.md](LICENSE.md)** | warunki licencyjne, definicje podmiotów, zakres dozwolonego użycia |
| **[INSTRUKCJA-UZYTKOWANIA.md](INSTRUKCJA-UZYTKOWANIA.md)** (ten dokument) | praca z gotową instalacją: strona główna, środowiska, okno komunikacji, okna równoległe, okno konfiguracji, dostępy, konta i tożsamość |
| **[INSTALACJA-I-KONFIGURACJA.md](INSTALACJA-I-KONFIGURACJA.md)** | pełna procedura instalacji, budowania warstw, doprowadzenia świeżej instalacji do pracy, wykaz zmiennych środowiska i przełączników, rozwiązywanie problemów |

**Rozejście się zakresów — stan rzeczywisty.** Rozdział [2](#2-uruchamianie-aplikacji) niniejszej instrukcji opisuje budowanie rdzenia, klienta i powłoki (2.2), procedurę wydania (2.5), przełączniki wiersza poleceń (2.3), wykaz zmiennych środowiska (2.6) oraz doprowadzenie świeżej instalacji do stanu zdolnego do rozmowy (2.7). Jest to dokładnie ten zakres, który [README.md](README.md) przypisuje dokumentowi instalacyjnemu. Dwie procedury uruchomienia opisane niezależnie w dwóch dokumentach mogą się z czasem rozejść, a Operator nie ma wtedy sposobu rozstrzygnięcia, która obowiązuje.

**Zasada obowiązująca do czasu rozstrzygnięcia.** W razie rozbieżności między niniejszą instrukcją a [INSTALACJA-I-KONFIGURACJA.md](INSTALACJA-I-KONFIGURACJA.md) w sprawach instalacyjnych — wymagań systemowych, kolejności budowania warstw, wartości zmiennych środowiska, procedury pierwszego uruchomienia i procedury doprowadzenia świeżej instalacji do pracy — **rozstrzyga dokument instalacyjny**. Rozdział 2 tej instrukcji należy czytać jako skrót operacyjny, wystarczający do uruchomienia produktu w typowym przypadku, lecz nie jako źródło prawdy o instalacji.

**[DO DECYZJI OPERATORA]** Podział zakresów między [INSTRUKCJA-UZYTKOWANIA.md](INSTRUKCJA-UZYTKOWANIA.md) a [INSTALACJA-I-KONFIGURACJA.md](INSTALACJA-I-KONFIGURACJA.md): czy rozdział 2 niniejszego dokumentu ma zostać skrócony do odesłania, czy przeciwnie — deklaracja podziału w [README.md](README.md) (rozdz. 17 i 19) ma zostać poprawiona tak, by uwzględnić skrót uruchomieniowy w instrukcji użytkowania. Do czasu rozstrzygnięcia obowiązuje zasada podana wyżej.

### 1.7 Podstawa opisu — stan kodu, na którym oparto instrukcję

Instrukcja opisuje zapis repozytorium bieżącego stanu repozytorium na — to znaczy HEAD gałęzi w chwili wydania tego dokumentu. Praca nad instrukcją zaczęła się w stanie sprzed prac naprawczych — stanie kodu, na którym oparł się audyt stojący u źródeł oznaczeń stanu z rozdziału 1.5 i zbiorczego zestawienia z rozdziału 18. Zanim dokument został domknięty, do gałęzi doszła **prace naprawcze**, złożona z sześciu pakietów (A–F), zamykająca część ustaleń audytu trzema kolejnymi zapisami:

| Obszar prac | Zawartość |
|---|---|
| Ścieżka wydania | złożenie produktu skryptami wydania oraz bramy jakości mierzące wywołania, nie tylko istnienie plików |
| Warstwa widoków | domknięcie warstwy wizualnej v2.0, zasilenie pulpitu danymi z rdzenia, pierwsze sprawdziany widoków |
| Dokumentacja | prace porządkowe nad dokumentacją: cztery pliki dokumentacji repozytorium ([README.md](README.md), `kontrakt.md`, `model-danych.md`, `shared/README.md` — ścieżki historyczne sprzed reorganizacji dokumentacji, zob. wiersz „Dokumentacja repozytorium” w tabeli niżej) wyrównane do stanu kodu po pracach naprawczych; zapis decyzji architektonicznych i znanych luk zaktualizowany deklaratywnie tym samym zapisem, zob. zastrzeżenie niżej |

Rozróżnienie nie jest formalnością: między stanem sprzed prac naprawczych a stanem bieżącym leży kilkadziesiąt zmienionych, dodanych i usuniętych plików, a część z nich dotyczy wprost zachowań opisanych w tej instrukcji — dokument musiał zostać przepisany w tych miejscach, nie tylko opatrzony nową datą.

Zmiany istotne dla treści dokumentu, wprowadzone pracami naprawczymi względem zapisu stanu sprzed prac naprawczych:

| Obszar | Zmiana | Rozdział opisujący stan po zmianie |
|---|---|---|
| Wydanie produktu | doszły skrypty `budowa/scripts/wydanie.sh` i `budowa/scripts/pakowanie.sh`; jednym poleceniem powstaje pakiet w `C:\DanacoConsole_App` — pakiet ten został zbudowany i uruchomiony (odpowiedź HTTP 200 na stronie głównej); instalatora NSIS nadal brak wobec braku narzędzia NSIS na maszynie budującej | [2.5](#25-wydanie-produktu-do-cdanacoconsole_app) |
| Powłoka | `tauri.conf.json` dostał `beforeBuildCommand` i `beforeDevCommand`; z `Cargo.toml` usunięto nieużywaną zależność `serde_json` | [2.2](#22-uruchomienie-przez-powłokę-tauri-w-drzewie-deweloperskim) |
| Bramy jakości | brama mierzy dziś **wywołania**, nie tylko obecność plików: doszły kontrole `budowa/scripts/wywolania-kontraktu.mjs` (komendy kontraktu bez wywołania) i `budowa/scripts/eksporty-go.mjs` (eksporty Go bez odbiorcy w torze produkcyjnym); doszły próg wielkości funkcji, próg wielkości arkusza CSS, hak pre-commit i przebieg CI, oraz `budowa/scripts/dym-pionu.mjs` — uruchomieniowy sprawdzian pełnego pionu na kanale echo | rozdz. 18 (tabela zbiorcza) |
| Warstwa wizualna | liczba klas używanych w widokach bez definicji CSS spadła z 36 do **1** (nie do zera) — pozostaje `.dn-karta-sesji` w `client/src/powloka/karty-sesji.ts:56`, zweryfikowane sprawdzianem pokrycia klas zastosowanym do kodu w stanie bieżącym; skrypt tego sprawdzianu, `budowa/scripts/klasy-css.mjs`, jeszcze nie należy do zapisu bieżącego stanu repozytorium — jego dopisanie do bramy jest osobnym, nieukończonym zadaniem; stos dymków komunikatów dostał pełne reguły stylu w obu motywach; widok sterowania dostał dymki `[?]` objaśnień i stan ładowania | [3.5](#35-motyw-jasny-i-ciemny), [8](#8-konfiguracja-sesji-i-okna--panel-sterowania) |
| Mission Control | plik `mission-control/dane-przykladowe.ts` **usunięty**; pulpit bierze stan wyłącznie z rdzenia | [3.4](#34-mission-control) |
| Strona główna | doszła strefa „Sesje w tle” zasilana odczytem `session.list` oraz zdarzeniami `session.changed` i `window.changed`, z uczciwymi stanami pustymi i powrotem do sesji przez `session.bind`/`session.focus`; zdarzenie `progress.changed` dostało pierwszego konsumenta (wykaz procesów w tle na pulpicie) | [3.2.3](#323-strefa-trzecia--sesje-w-tle) |
| Dokumentacja repozytorium | `shared/README.md` doprowadzony do 53 komend; `model-danych.md` dostał sekcję pokrycia; [README.md](README.md) i rozdziały 1–3 `kontrakt.md` wyrównane do kontraktu wykonawczego — te cztery pliki ([README.md](README.md), `budowa/docs/kontrakt.md`, `budowa/docs/model-danych.md`, `budowa/shared/README.md`) są jedynymi zmienionymi w bieżącym stanie repozytorium (zweryfikowano zestawieniem zmian w drzewie roboczym); ścieżki `budowa/docs/kontrakt.md`, `budowa/docs/model-danych.md` i `budowa/shared/README.md` opisują drzewo repozytorium budowy sprzed reorganizacji dokumentacji — nie istnieją w bieżącym drzewie, a ich rolę informacyjną przejęły [Kontrakty komunikacji](architektura/kontrakty-komunikacji.md) i [Model danych](architektura/model-danych.md) w katalogu `docs/architektura/` korzenia repozytorium. Komunikat tego zapisu deklaruje ponadto aktualizację zapisu decyzji architektonicznych i znanych luk — dopisanie decyzji o moście MCP oraz zamknięcie trzech luk i aktualizację czwartej — lecz materiały te nie są obecne w tym repozytorium. Operator czytający wyłącznie kod repozytorium nie ma jak zweryfikować tej części zmiany z samego drzewa źródłowego — musi albo zaufać komunikatowi zapisu, albo sięgnąć po katalog `opracowania/` osobno, poza tym pobraniem drzewa roboczego | rozdz. 6, 10, 18 (w miejscach cytujących te rejestry) |
| Sprawdziany | doszły sprawdziany widoków w jsdom (stos dymków, motyw, pulpit, strefa sesji, podsumowanie ustawień) oraz testy logiki startu powłoki w Ruście; kontrole żywotności obecne w bramie w bieżącym stanie repozytorium — osiągalność drzewa interfejsu (`budowa/scripts/osiagalnosc.mjs`: 435/435 osiągalnych, 0 naruszeń), wywołania komend kontraktu i eksporty Go bez odbiorcy (`budowa/scripts/wywolania-kontraktu.mjs`, `budowa/scripts/eksporty-go.mjs`) — kończą się bez naruszeń niezastrzeżonych długiem jawnym, zweryfikowane ponownym uruchomieniem na czystej kopii tego zapisu; kontrola pokrycia klas CSS (`budowa/scripts/klasy-css.mjs`) do bramy na tym zapisie jeszcze nie należy — jej dopisanie jest osobnym, nieukończonym zadaniem — i przy zastosowaniu do kodu z tego zapisu wykazuje jedną klasę bez definicji — zob. wiersz „Warstwa wizualna” wyżej; zdanie „brama jest pełna: 0 naruszeń” obejmuje więc kontrole, które na tym zapisie faktycznie istnieją, nie całość kontroli dopisanych później | rozdz. 18 (tabela zbiorcza) |

**Co prace naprawcze NIE zmieniły** — pozostaje w mocy wszystko, co ten dokument opisuje jako ATRAPĘ, NIEZINTEGROWANE albo BRAK: zaczyn kanału głównego i konta na świeżej instalacji, ciągłość rozmowy między turami (`--resume` nadal nigdy nie podawane), dotarcie modelu/modelu zapasowego/nakładu rozumowania/hosta wykonania do wywołania, egzekutory izolacji, 64 z 79 wierszy katalogu okien operacyjnych w bazie (rozdz. 7), silnik kolejek, 39 narzędzi platformy dla modelu bez konsumenta, większość nawigacji platformy (wyjątek opisany niżej), historia rozmowy po odświeżeniu (`message.list` nadal bez wywołania), zasięg wykonania jako etykieta, kanał SSH. Jedyny wyjątek w nawigacji platformy: `session.list`/`session.bind`/`session.focus` są dziś wywoływane — przez nową strefę „Sesje w tle” na stronie głównej.

**Powód, dla którego sama adnotacja o zapisie nie wystarcza.** Instrukcja jest dostarczana **wewnątrz** pakietu produktu, obok plików wykonywalnych zbudowanych z tego samego drzewa roboczego. Gdyby dokument opisywał stan starszy niż zbudowany produkt, Operator czytałby opis funkcji innych niż te, które ma przed sobą — a dokument, którego zadaniem jest umożliwienie samodzielnej obsługi, stałby się źródłem błędu. Dlatego wszystkie twierdzenia o zachowaniu klienta, powłoki i skryptów zostały zweryfikowane ponownie względem zapisu bieżącego stanu repozytorium, a nie względem zapisu stanu sprzed prac naprawczych, na którym praca się zaczęła.

**Stan opisywany.** Dokument opisuje bieżący stan repozytorium po zamknięciu prac
naprawczych domykających ustalenia audytu technicznego; drzewo robocze zostało
zmierzone 12 sierpnia 2026. Pole **Data** metryki produktowej niesie datę wydania
dokumentu, nie datę tego pomiaru — obie daty są odrębnymi informacjami.

**Zgodność z pakietem.** Dokument opisuje ten sam stan kodu, z którego zbudowano
pliki wykonywalne leżące obok niego w `C:\DanacoConsole_App`.

---

## 2. Uruchamianie aplikacji

> **Zakres tego rozdziału.** Rozdział podaje procedurę wystarczającą do uruchomienia produktu i doprowadzenia go do stanu zdolnego do rozmowy. Pełna procedura instalacyjna — wymagania systemowe, przygotowanie stacji roboczej, konfiguracja bazy, dostawców i narzędzi, weryfikacja poprawności instalacji oraz aktualizacja i deinstalacja — należy do dokumentu **[INSTALACJA-I-KONFIGURACJA.md](INSTALACJA-I-KONFIGURACJA.md)** i w razie rozbieżności to on rozstrzyga (rozdział [1.6](#16-miejsce-tego-dokumentu-wśród-dokumentów-produktu)).

### 2.1 Warunki wstępne

Do uruchomienia produktu w drzewie deweloperskim potrzebne są:

| Składnik | Przeznaczenie | Uwaga |
|---|---|---|
| **Go** | zbudowanie rdzenia (`danaco-console`) | moduł `danacoconsole`, plik `go.mod` w katalogu `budowa` |
| **Node.js wraz z npm** | zbudowanie i uruchomienie interfejsu | klient korzysta z Vite 6 i TypeScriptu 5.8 |
| **Rust wraz z cargo i Tauri CLI** | zbudowanie powłoki natywnej | wymagane wyłącznie dla wariantu z oknem natywnym |
| **program `claude`** (Claude Code CLI) | wykonanie tury modelu kanałem głównym | musi być osiągalny w `PATH` albo wskazany w wierszu rejestru kanału |

Program `claude` jest zależnością zewnętrzną kanału głównego. Rdzeń wywołuje go w trybie strumieniowym i czyta strumień JSON-lines. Jego brak nie wstrzymuje startu rdzenia ani otwarcia okna — objawia się dopiero przy pierwszej wypowiedzi, jako błąd kanału.

> **Uwaga.** Katalog ustawień rdzenia zawiera klucz `harness.program_claude`, opisany jako wskazanie ścieżki programu. Klucz ten **nie jest dziś odczytywany** przez kod wykonawczy — ścieżka pochodzi wyłącznie z wiersza rejestru kanału, a w jego braku ze stałej `claude` rozstrzyganej przez `PATH`. Wpisanie wartości pod ten klucz jest ZAPISEM BEZ SKUTKU.

### 2.2 Uruchomienie przez powłokę Tauri w drzewie deweloperskim

Powłoka jest wariantem zalecanym, ponieważ sama stawia rdzeń w tle i sama rozstrzyga, skąd wziąć interfejs.

**Kolejność czynności:**

1. **Zbuduj rdzeń.** Z katalogu `budowa` wykonaj budowanie punktu wejścia `./server/cmd/danaco-console`. Plik wykonywalny nazywa się `danaco-console.exe` (Windows) albo `danaco-console`.
2. **Umieść plik wykonywalny w miejscu, w którym powłoka go znajdzie.** Powłoka przeszukuje miejsca w ustalonej kolejności:
   1. ścieżka wskazana zmienną `DANACO_RDZEN`,
   2. katalog pliku wykonywalnego powłoki oraz jego podkatalog `rdzen`,
   3. katalog `budowa`, `budowa\bin`, `budowa\server`,
   4. `C:\DanacoConsole_App`.
3. **Zainstaluj zależności interfejsu.** Z katalogu `budowa\client` wykonaj `npm install`. Ręczne zbudowanie pakietu (`npm run build`) **nie jest już konieczne** przy uruchamianiu przez powłokę: konfiguracja powłoki zawiera dziś polecenia budowania klienta i sama je wywoła — `cargo tauri build` uruchamia przed budowaniem `npm run build` w katalogu `..\..\client`, a `cargo tauri dev` uruchamia `npm run dev` w tym samym katalogu, nie czekając na jego zakończenie, i pobiera interfejs spod adresu `http://localhost:5173`. Ręczne zbudowanie pakietu pozostaje potrzebne wyłącznie wtedy, gdy interfejs ma być serwowany przez sam rdzeń, bez powłoki (rozdział [2.3](#23-uruchomienie-samego-rdzenia)).
4. **Uruchom powłokę** poleceniem `cargo tauri dev` z katalogu `budowa\desktop\src-tauri` albo — po zbudowaniu — uruchamiając wytworzony plik wykonywalny.

**Co robi powłoka przy starcie:**

- odczytuje ustawienia ze zmiennych środowiska (`DANACO_PORT`, `DANACO_RDZEN`, `DANACO_ADRES_INTERFEJSU`);
- sprawdza, czy pod portem rdzenia ktoś już nasłuchuje; **jeżeli tak — nie stawia drugiego procesu**, tylko dołącza do zastanego;
- w przeciwnym razie stawia rdzeń jako proces potomny z przełącznikami `--role all --port <port>`, bez okna konsoli i we własnej grupie procesów, po czym czeka na nasłuch do dziesięciu sekund;
- kieruje wyjście rdzenia do pliku `rdzen-powloki.log` w katalogu danych;
- rozstrzyga źródło interfejsu w kolejności: adres wskazany zmienną → tryb deweloperski → pakiet serwowany przez rdzeń pod `http://127.0.0.1:<port>` → pakiet osadzony w powłoce;
- otwiera okno o rozmiarze początkowym 1440 × 900 (minimalny 960 × 640) i stawia ikonę w zasobniku systemowym.

Brak binarki rdzenia, odmowa startu albo minięcie czasu oczekiwania **nie wstrzymują otwarcia okna**. Okno pojawia się mimo to, a przyczyna trafia do opisu stanu rdzenia dostępnego z menu zasobnika.

**Menu ikony w zasobniku** zawiera pięć zawsze czynnych pozycji:

| Pozycja | Skutek |
|---|---|
| Pokaż okno | przywraca i ogniskuje okno główne |
| Wskaż katalog roboczy… | otwiera natywne okno wyboru katalogu i rozgłasza wybór do interfejsu |
| Stan rdzenia | pokazuje opis: czy rdzeń pracuje, jego identyfikator procesu, adres i ścieżkę dziennika |
| Zatrzymaj rdzeń | jawne zatrzymanie procesu rdzenia wraz z komunikatem o wyniku |
| Zakończ powłokę | kończy powłokę; rdzeń i procesy sesji pracują dalej |

Podwójne kliknięcie ikony zasobnika przywraca okno.

### 2.3 Uruchomienie samego rdzenia

Rdzeń da się uruchomić bez powłoki — jest to wariant właściwy dla pracy z interfejsem w przeglądarce oraz dla diagnostyki.

Plik wykonywalny przyjmuje pięć przełączników:

| Przełącznik | Znaczenie | Wartość domyślna |
|---|---|---|
| `--role` | rola procesu: `hub`, `agent`, `all` | `all` |
| `--port` | port nasłuchu | `17870` |
| `--dane` | katalog danych rdzenia | `%LOCALAPPDATA%\DanacoConsole` |
| `--klient` | katalog pakietu interfejsu serwowanego obok gniazda | `client\dist` względem katalogu roboczego |
| `--profile` | katalog profili kanału głównego (konta rotacji) | pusty |

Znaczenie ról:

- **`hub`** — wyłącznie tor interfejsu: nasłuch WebSocket na wskazanym porcie;
- **`agent`** — wyłącznie tor wykonawczy: żądania kontraktu przyjmowane wierszami wejścia standardowego, odpowiedzi wierszami wyjścia; żaden port nie jest zajmowany;
- **`all`** — oba tory w jednym procesie; jest to rola, w której powłoka stawia rdzeń.

Przy pierwszym uruchomieniu rdzeń zakłada katalog danych i stosuje komplet migracji bazy. Migracje są wkompilowane w plik wykonywalny, stosowane w jednej transakcji na krok i opatrzone sumą kontrolną, więc uruchomienie na istniejącej bazie nie wymaga żadnej czynności ręcznej. Baza jest jednym plikiem `danaco-console.db` wewnątrz katalogu danych.

Po uruchomieniu rdzenia interfejs otwiera się pod adresem `http://127.0.0.1:17870/`, o ile katalog `client\dist` istnieje i został wskazany (bezpośrednio przełącznikiem `--klient`, zmienną `DANACO_KATALOG_KLIENTA` albo pośrednio przez powłokę). Nieistniejący katalog pakietu **nie wstrzymuje nasłuchu** — gniazdo pracuje bez plików statycznych, lecz przeglądarka nie ma wtedy czego wyświetlić.

### 2.4 Uruchomienie interfejsu z serwera rozwojowego

Wariant najszybszy w pracy nad interfejsem: rdzeń i klient uruchamiane osobno.

1. W katalogu `budowa\client` uruchom `npm run dev`. Serwer rozwojowy Vite podnosi się na porcie 5173.
2. Uruchom rdzeń (rozdział 2.3) na porcie 17870.
3. Otwórz `http://localhost:5173/`.

Klient wylicza adres gniazda WebSocket z lokalizacji dokumentu: schemat (`ws:` albo `wss:`), nazwa hosta strony oraz **stały port 17870** i ścieżka `/ws`. Adres da się nadpisać zmienną budowania `VITE_ADRES_RDZENIA`.

> **Ostrzeżenie eksploatacyjne.** Numer portu jest w kliencie wartością stałą. Zmiana `DANACO_PORT` po stronie rdzenia **rozspaja klienta z rdzeniem**, o ile nie zostanie równocześnie podana zmienna budowania `VITE_ADRES_RDZENIA`.

> **Ograniczenie znane — pakiet osadzony w powłoce nie zestawia połączenia.** Klient wylicza nazwę hosta gniazda z lokalizacji dokumentu. Strona podana przez rdzeń ma nazwę hosta `127.0.0.1` i trafia w gniazdo poprawnie. Strona pobrana z zasobu **osadzonego w pliku wykonywalnym powłoki** ma natomiast nazwę hosta `tauri.localhost`, więc wyliczony adres gniazda nie wskazuje rdzenia i **połączenie nie zostaje zestawione**. Stan ten jest **luką produktu**, a nie rozwiązaniem zamierzonym: [README.md](README.md) (rozdz. 3.3 oraz zestawienie ograniczeń w rozdz. 18) klasyfikuje go jako **DZIAŁA CZĘŚCIOWO** i nazywa **luką nazwy hosta zasobu osadzonego**. Pierwszeństwo, jakie powłoka daje pakietowi serwowanemu przez rdzeń, jest **obejściem tej luki**, nie decyzją projektową — dzięki niemu okno działa mimo wady, dopóki rdzeń nasłuchuje i serwuje pakiet. Praktyczny wniosek: pakiet osadzony w powłoce nie jest dziś drogą sprawną i nie należy na nim opierać wdrożenia.

Repozytorium zawiera również skrypt `budowa/scripts/pokaz.sh`, który wykonuje powyższe czynności jednym poleceniem: sprawdza kontrolę typów, buduje pakiet klienta, podnosi serwer rozwojowy i rdzeń, po czym otwiera przeglądarkę. Zatrzymanie: `bash budowa/scripts/pokaz.sh --stop`. Skrypt jest narzędziem podglądu stanu prac, nie procedurą wdrożeniową.

### 2.5 Wydanie produktu do `C:\DanacoConsole_App`

Model katalogów produktu przewiduje `C:\DanacoConsole_App` jako wyjście budowania — gotowy do uruchomienia produkt, nigdy niezawierający źródeł. Katalog ten jest dziś **wytwarzany skryptem wydania** repozytorium. Dowodem materialnym jest sam ten dokument: leży w `C:\DanacoConsole_App` obok plików `danaco-console.exe`, `Danaco Console.exe` oraz katalogu `client\dist`, wytworzonych opisaną niżej procedurą.

#### 2.5.1 Skrypty wydania i pakowania

Za wydanie odpowiadają dwa skrypty o rozdzielonych zadaniach:

| Skrypt | Zadanie |
|---|---|
| `budowa/scripts/wydanie.sh` | **buduje** trzy części produktu i przekazuje je skryptowi pakowania |
| `budowa/scripts/pakowanie.sh` | **składa** katalog wyjściowy z gotowych artefaktów; niczego nie buduje |

Podział jest celowy: pakowanie da się wywołać samodzielnie, gdy artefakty już istnieją, a wydanie pozostaje jedynym miejscem, w którym uruchamiane są narzędzia budowania.

**Wywołanie:**

```
bash budowa/scripts/wydanie.sh [katalog-wyjścia]
```

Katalog wyjścia jest argumentem opcjonalnym; jego wartością domyślną jest `C:/DanacoConsole_App`.

#### 2.5.2 Przebieg wydania

Skrypt wykonuje pięć kroków w ustalonej kolejności:

| Krok | Czynność | Uwaga |
|---|---|---|
| 1 | **Zależności klienta** — `npm ci` w katalogu `budowa/client` | pomijane, gdy katalog `node_modules` już istnieje; wymuszenie instalacji = usunięcie tego katalogu |
| 2 | **Pakiet klienta** — `npm run build` (kontrola typów `tsc`, następnie `vite build`) | wynik trafia do `client/dist` |
| 3 | **Rdzeń** — `go build ./server/cmd/danaco-console` | wynik trafia do katalogu przejściowego, nie od razu do wyjścia |
| 4 | **Powłoka natywna** — `cargo tauri build --no-bundle` w katalogu `desktop/src-tauri` | krok pomijalny, patrz niżej |
| 5 | **Pakowanie** — złożenie katalogu wyjściowego | wykonywane skryptem `pakowanie.sh` |

**Struktura katalogu wyjściowego po wydaniu:**

| Element | Zawartość |
|---|---|
| `danaco-console.exe` | rdzeń; uruchomiony bez przełączników przyjmuje rolę `all` i port `17870`, serwuje też interfejs |
| `Danaco Console.exe` | powłoka natywna — obecna wyłącznie wtedy, gdy jej budowanie się powiodło |
| `client\dist\` | pakiet interfejsu; leży obok rdzenia — trafia do niego powłoka albo sam rdzeń, na warunkach opisanych niżej |
| `*.md` | dokumentacja produktu |

Układ ten odpowiada kolejności poszukiwań stosowanej przez powłokę (rozdział 2.2): powłoka szuka pliku `danaco-console.exe` obok siebie, a katalog `client\dist` wskazuje rdzeniowi zmienną `DANACO_KATALOG_KLIENTA`. Rdzeń **sam** rozstrzyga tę ścieżkę wyłącznie względem katalogu roboczego swojego procesu (`filepath.Join("client","dist")`, rozdział [2.3](#23-uruchomienie-samego-rdzenia)) — nie sięga po własną ścieżkę pliku wykonywalnego. Obok binarki rdzenia szuka **powłoka** (najpierw `<katalog rdzenia>\client\dist`, potem `<katalog rdzenia>\dist`, potem `client\dist` z drzewa repozytorium) i przekazuje wynik rdzeniowi tą samą zmienną `DANACO_KATALOG_KLIENTA`; powłoka dodatkowo stawia katalog roboczy procesu rdzenia na katalogu jego binarki, więc w praktyce obie ścieżki poszukiwań prowadzą do tego samego miejsca. Wydanie nie wymaga zatem żadnej dodatkowej konfiguracji ścieżek — pod warunkiem uruchamiania rdzenia przez powłokę albo z katalogu jego własnej binarki jako katalogu roboczego; uruchomienie rdzenia z innego katalogu roboczego, bez powłoki, wymaga jawnego przełącznika `--klient` albo zmiennej `DANACO_KATALOG_KLIENTA`.

> **Zasada ochrony dokumentacji.** Skrypt pakowania **nigdy nie usuwa ani nie nadpisuje plików `*.md`** w katalogu wyjściowym. Czyszczenie obejmuje wyłącznie podkatalog `client\`, i to z wyłączeniem plików `*.md`. Powtórne wydanie do katalogu, w którym leży dokumentacja produktu, jest więc bezpieczne — dokumenty przetrwają wymianę plików wykonywalnych i pakietu interfejsu.

#### 2.5.3 Warianty wydania

**Wariant pełny** — rdzeń, pakiet interfejsu i powłoka natywna. Wymaga zainstalowanych narzędzi Go, Node.js wraz z npm, Rust wraz z cargo oraz `tauri-cli`.

**Wariant „rdzeń + klient”** — bez powłoki natywnej. Powstaje w dwóch przypadkach:

1. **na żądanie**, przez ustawienie zmiennej `DANACO_BEZ_POWLOKI=1` przed wywołaniem skryptu;
2. **samoczynnie**, gdy w systemie nie ma `cargo` albo `tauri-cli`; skrypt zgłasza wtedy pominięcie kroku wraz z poleceniem instalacji narzędzia i **kontynuuje pracę**, zamiast przerywać wydanie.

Zachowanie z punktu drugiego jest zastosowaniem zasady fail-open: brak narzędzia opcjonalnego nie unieważnia wydania, lecz zawęża jego wynik. Wynik wariantu „rdzeń + klient” uruchamia się plikiem `danaco-console.exe`, a interfejs otwiera się w przeglądarce pod adresem `http://127.0.0.1:17870/`. Skrypt wypisuje na końcu, który wariant powstał i jak go uruchomić.

**Powtórka budowania pakietu klienta jest zamierzona.** Krok 2 buduje pakiet klienta, a krok 4 uruchamia `cargo tauri build`, które przez `beforeBuildCommand` buduje ten pakiet ponownie. Powtórzenie jest tanie i celowe: dzięki niemu wariant bez powłoki również dostaje świeży pakiet interfejsu.

#### 2.5.4 Instalator NSIS — osobny krok

Wydanie wywołuje `cargo tauri build` z przełącznikiem `--no-bundle`, więc jego wynikiem jest **sam plik wykonywalny powłoki**, bez instalatora. Konfiguracja powłoki przewiduje wprawdzie cel `nsis` wraz z metryką wydawcy, prawami autorskimi i kompletem ikon, lecz zbudowanie instalatora wymaga zainstalowanego narzędzia NSIS i pozostaje **osobnym, ręcznym krokiem**: `cargo tauri build` (bez `--no-bundle`) uruchomione w katalogu `desktop/src-tauri`.

#### 2.5.5 Czego wydanie nie robi

Aby uniknąć nieporozumienia co do dojrzałości procedury, poniżej wprost to, czego wydanie **nie** obejmuje:

- **nie jest instalacją** — nie zakłada wpisu w rejestrze systemu, nie tworzy skrótów, nie rejestruje usługi ani zadania startowego;
- **nie sprawdza zależności zewnętrznych** — nie weryfikuje obecności programu `claude`, od którego zależy kanał główny;
- **nie zasiewa danych** — katalog danych, rejestr kanałów i katalog kont pozostają w stanie wyjściowym; doprowadzenie instalacji do stanu zdolnego do rozmowy opisuje rozdział [2.7](#27-doprowadzenie-świeżej-instalacji-do-stanu-zdolnego-do-rozmowy);
- **nie ma drogi odinstalowania ani aktualizacji przyrostowej** — kolejne wydanie wymienia pliki wykonywalne i pakiet interfejsu w miejscu;
- **nie podpisuje plików wykonywalnych** żadnym certyfikatem.

**[DO DECYZJI OPERATORA]** Sposób pakowania produktu: czy rdzeń ma być zasobem osadzonym w powłoce, czy — jak dziś — odrębnym plikiem wykonywalnym obok niej; oraz czy pakiet interfejsu ma być serwowany przez rdzeń, czy osadzany w powłoce. Rozstrzygnięcie ma dziś skutek praktyczny, ponieważ droga pakietu osadzonego nie zestawia połączenia (luka nazwy hosta zasobu osadzonego, rozdział 2.4). Osobno wymaga rozstrzygnięcia, czy instalator NSIS ma zostać włączony do skryptu wydania jako krok warunkowy, czy pozostać czynnością odrębną.

### 2.6 Zmienne środowiska i przełączniki wiersza poleceń

Rdzeń odczytuje dokładnie pięć zmiennych środowiska. Warstwy ustalania wartości: **wartość domyślna → zmienna środowiska → argument wywołania**; argument wygrywa ze zmienną, zmienna z wartością domyślną. Zmienna pusta lub nieustawiona pozostawia wartość dotychczasową.

| Zmienna | Odpowiednik w wierszu poleceń | Znaczenie |
|---|---|---|
| `DANACO_ROLA` | `--role` | rola procesu: `hub`, `agent`, `all` |
| `DANACO_PORT` | `--port` | port nasłuchu rdzenia |
| `DANACO_KATALOG_DANYCH` | `--dane` | katalog danych; puste = `%LOCALAPPDATA%\DanacoConsole` |
| `DANACO_KATALOG_KLIENTA` | `--klient` | katalog pakietu interfejsu |
| `DANACO_KATALOG_PROFILI` | `--profile` | katalog profili kanału głównego |

Powłoka odczytuje trzy zmienne:

| Zmienna | Znaczenie |
|---|---|
| `DANACO_PORT` | port, pod którym powłoka szuka i stawia rdzeń |
| `DANACO_RDZEN` | jawna ścieżka pliku wykonywalnego rdzenia |
| `DANACO_ADRES_INTERFEJSU` | jawny adres, spod którego okno pobiera interfejs |

Klient przyjmuje wartości w chwili budowania (przedrostek `VITE_`): `VITE_ADRES_RDZENIA`, `VITE_PROJEKT`, `VITE_TYTUL_OKNA`, `VITE_MODUL`, `VITE_KANAL_MODELU` oraz `VITE_KATALOGI_ROBOCZE` (lista rozdzielona średnikiem). Zmienne te należą do katalogu `client`, nie do pliku `.env` rdzenia.

Wzorzec `budowa/.env.example` zawiera wyłącznie zmienne rzeczywiście odczytywane przez kod. Nie ma w nim ani ścieżki programu `claude`, ani kluczy dostawców API — te ostatnie z założenia nie żyją w środowisku, lecz w rejestrze kanałów i w katalogu kont.

Osobno stoi jedna zmienna, która nie należy ani do rdzenia, ani do powłoki, ani do klienta, lecz do **procedury wydania**: `DANACO_BEZ_POWLOKI`. Ustawiona na `1` przed wywołaniem skryptu wydania pomija budowanie powłoki natywnej i daje wariant „rdzeń + klient” (rozdział [2.5](#25-wydanie-produktu-do-cdanacoconsole_app)). Nie ma ona żadnego wpływu na działanie już zbudowanego produktu.

Pełny, rozstrzygający wykaz zmiennych środowiska wraz z ich zastosowaniem w scenariuszach wdrożeniowych zawiera dokument [INSTALACJA-I-KONFIGURACJA.md](INSTALACJA-I-KONFIGURACJA.md) (rozdział 12).

### 2.7 Doprowadzenie świeżej instalacji do stanu zdolnego do rozmowy

To jest najważniejszy rozdział procedury uruchomieniowej. **Świeżo utworzona baza nie pozwala na wykonanie ani jednej tury modelu.** Przyczyny są trzy i występują równocześnie:

1. **Rejestr kanałów modelu jest pusty.** Żadna migracja nie wstawia wiersza do tabeli kanałów, więc rejestr startuje bez wpisów.
2. **Interfejs wskazuje kanał wartością `cli`, czyli RODZAJEM kanału, a nie identyfikatorem jego wiersza.** Nawet gdyby wiersz istniał, domyślne wskazanie okna go nie trafi, o ile jego kod nie brzmi dokładnie `cli`.
3. **Pula kont kanału głównego jest pusta**, a pusta pula odmawia wywołania komunikatem o wyczerpaniu limitów.

Rdzeń obsługuje komendy zakładania kanału i konta (`channel.add`, `account.add`) w pełnym zakresie. **Interfejs nie zawiera jednak ekranu do zakładania kanału modelu** — wywołuje wyłącznie odczyt wykazu kanałów. Konta natomiast zakłada się z interfejsu, w oknie „Modele, konta i tożsamość”.

**Kolejność doprowadzenia:**

| Krok | Czynność | Droga |
|---|---|---|
| 1 | Założenie wiersza kanału głównego rodzaju `cli` | komenda kontraktu `channel.add` przez WebSocket — **brak drogi z głównej aplikacji** (uściślenie poniżej: moduł automatycznego zapewnienia kanału istnieje, lecz jest wpięty tylko do stanowiska podglądu) |
| 2 | Założenie co najmniej jednego konta rodzaju `cli` ze wskazaniem katalogu konfiguracji programu | interfejs: Centrum dowodzenia → „Modele, konta i tożsamość” → zakładka „Konta modeli i code CLI” |
| 3 | Restart rdzenia | konieczny, ponieważ pula rotacji kont budowana jest jednorazowo przy montażu rdzenia |
| 4 | Wskazanie kanału w oknie komunikacji | panel sterowania okna → kontrolka „Model” |

Krok 3 nie jest zaleceniem ostrożnościowym, lecz wymogiem: zmiany katalogu kont **nie odświeżają puli rotacji w działającym procesie**. Inaczej zachowuje się rejestr kanałów — ten jest przebudowywany natychmiast po każdej zmianie wiersza.

Krok 1 wymaga wysłania komendy kontraktu poza interfejsem — klientem WebSocket skierowanym na `ws://127.0.0.1:17870/ws`, z kopertą żądania `channel.add`. Kształt ładunku opisuje kontrakt w katalogu `shared/`.

#### Moduł zapewnienia kanału — napisany, lecz niepodpięty do aplikacji

Stan opisany wyżej wymaga jednego uściślenia, bez którego obraz byłby niepełny. W kodzie klienta **istnieje gotowy moduł automatycznego zapewnienia kanału głównego** (`client/src/rozmowa/zapewnienie-kanalu.ts`). Realizuje on pełną procedurę:

1. odpytuje rejestr komendą odczytu wykazu kanałów, ograniczoną do kanałów czynnych;
2. szuka w odpowiedzi kanału rodzaju pierwszego z wykazu kontraktu, czyli `cli`, oznaczonego jako czynny;
3. gdy taki wiersz zastanie — zwraca jego identyfikator;
4. gdy nie zastanie — zakłada wiersz komendą `channel.add`, pod nazwą „Kanał główny — Claude Code CLI”, rodzaju `cli`, oznaczony jako czynny, i zwraca identyfikator nowo powstałego wiersza;
5. przy niepowodzeniu nie przerywa niczego poza własnym wywołaniem — oddaje opis przeszkody, a okno pozostaje czynne.

**Moduł jest jednak wpięty wyłącznie do stanowiska podglądu rozmowy** (`client/src/rozmowa/podglad-rozmowy.html`) — osobnej strony służącej pracy nad warstwą rozmowy. Główne wejście aplikacji (`index.html` → `main.ts`) go **nie wywołuje**. Z punktu widzenia Operatora pracującego z produktem droga automatycznego zapewnienia kanału **nie istnieje**, mimo że kod ją realizujący jest napisany, kompletny i podpięty do kanału kontraktu.

Fakt ten ma znaczenie dla rozstrzygnięcia poniżej: jeden z trzech wymienionych wariantów nie wymaga napisania nowego kodu, lecz wyłącznie wpięcia istniejącego modułu do ścieżki produkcyjnej wraz z rozstrzygnięciem, w którym momencie uzgodnienia ma on działać.

**[DO DECYZJI OPERATORA]** Sposób trwałego rozwiązania stanu wyjściowego instalacji: (a) zaczyn migracji wstawiający wiersz kanału głównego; (b) ekran zakładania kanałów w interfejsie; (c) wpięcie istniejącego modułu zapewnienia kanału do głównego wejścia aplikacji, z automatycznym zapewnieniem kanału przy uzgodnieniu. Wszystkie trzy warianty są technicznie osiągalne, przy czym wariant (c) jest już napisany i wymaga wyłącznie podpięcia. Wybór należy do Operatora — instrukcja go nie przesądza, ponieważ automatyczne zakładanie wiersza rejestru przy każdym uzgodnieniu jest decyzją produktową o skutkach wykraczających poza wygodę pierwszego uruchomienia.

---

## 3. Struktura interfejsu

### 3.1 Trzy widoki najwyższego rzędu

Interfejs ma trzy trasy — trzy widoki najwyższego rzędu, między którymi przełącza się grupą przycisków na pasku:

| Trasa | Nazwa widoczna | Rola |
|---|---|---|
| `strona-glowna` | **Centrum dowodzenia** | wejście do produktu, wybór środowiska, listwa ustawień |
| `srodowisko` | **Środowisko** | powłoka środowiska ze sceną okien komunikacji |
| `pulpit` | **Mission Control** | pulpit orkiestracji |

Trasą początkową jest Centrum dowodzenia — uruchomienie pokazuje przedpokój pracy, nie okno komunikacji wyrwane z kontekstu. Nazwa trasy zapisuje się w adresie dokumentu, więc odświeżenie strony wraca do tego samego widoku.

Widoki powstają leniwie, przy pierwszym wejściu, i nie są porzucane przy wyjściu. Karty opuszczonego środowiska trwają w tle i wracają z pełnym stanem. Wyjątkiem jest widok środowiska: przygotowuje się od razu przy uruchomieniu, ponieważ to on podpina przepływ komunikatów — uzgodnienie z rdzeniem zdąży się dokonać, zanim Operator wybierze środowisko.

Przełącznik tras stoi w grupie akcji paska górnego; trasa bieżąca jest oznaczona atrybutem dostępności, a nie odebraniem klikalności. Ikony pozycji: dom (Centrum dowodzenia), folder (Środowisko), zegar (Mission Control).

### 3.2 Centrum dowodzenia — strona główna

Centrum dowodzenia dzieli się na **cztery strefy** o malejącej wadze wizualnej: karty czterech środowisk, kafle komponentów własnych, wykaz sesji trwających w tle oraz listwę ustawień.

#### 3.2.1 Strefa pierwsza — karty czterech środowisk

Cztery karty, każda z godłem, nazwą kanoniczną, mottem i jednozdaniowym opisem trybu pracy:

| Środowisko | Motto | Opis |
|---|---|---|
| **TalkIn** | Myśl. Analizuj. Rozumiej. | Wiedza, treść, dokumenty, badania i tłumaczenia. |
| **WorkSpace** | Planuj. Organizuj. Realizuj. | Projekty, procesy, automatyzacje i produkty. |
| **CodeStudio** | Projektuj. Buduj. Rozwijaj. | Programowanie, terminale i architektura systemów. |
| **MultitaskingAI** | Deleguj. Koordynuj. Nadzoruj. | Orkiestracja zespołów modeli i procesy autonomiczne. |

Wybór karty przenosi do widoku środowiska i przestawia powłokę na wskazane środowisko. Środowisko, w którym Operator pracuje, dostaje na stronie głównej oznaczenie czynnego.

#### 3.2.2 Strefa druga — kafle komponentów własnych

Cztery kafle prowadzące do zbudowania rzeczy, która potem pracuje w środowisku: **Automations** („Zbuduj automatykę”), **Agents** („Skonfiguruj agenta”), **Workspace** („Załóż projekt”), **Assistant** („Ustaw profil asystenta”).

**Stan faktyczny: DZIAŁA CZĘŚCIOWO.** Naciśnięcie kafla nie otwiera kreatora komponentu — przenosi do środowiska i wskazuje pozycję nawigacji odpowiadającą komponentowi (Automations prowadzi do środowiska WorkSpace, ponieważ nie występuje w żadnym wykazie nawigacji). Ponieważ wybór pozycji nawigacji nie otwiera dziś okna operacyjnego (rozdział 4.2), kafel doprowadza Operatora do sceny okien komunikacji, a nie do konfiguracji komponentu.

#### 3.2.3 Strefa trzecia — sesje w tle

**Stan faktyczny: DZIAŁA.** Trzecia strefa Centrum dowodzenia nosi nagłówek **„Sesje w tle”** i pokazuje sesje trwające na rdzeniu **poza bieżącym połączeniem**. Jest to jedyne miejsce w produkcie, w którym Operator widzi pracę toczącą się bez otwartego przed nim okna.

**Skąd biorą się dane.** Strefa zasila się odczytem wykazu sesji, wysyłanym z żądaniem żywego odpisu obecności, i odświeża się na dwóch zdarzeniach rdzenia: zmiany sesji oraz zmiany okna. Zdarzenia okien są scalane — kilka zmian następujących po sobie daje jedno ponowne odpytanie, po krótkiej zwłoce. Odpowiedź przedawniona nie nadpisuje świeższej, więc wykaz nie „skacze” przy szybkich zmianach.

**Co pokazuje wiersz sesji.** Każdy napis pochodzi z rdzenia, żadna liczba nie jest wyliczana miejscowo:

| Element wiersza | Źródło |
|---|---|
| tytuł sesji | pole tytułu sesji; przy jego braku — identyfikator sesji, nigdy nazwa wymyślona |
| środowisko i moduł | żywy odpis obecności; odpis nieobecny pomija cały fragment zamiast pokazać wartość zastępczą |
| liczba okien i liczba okien w strumieniu | żywy odpis obecności, w postaci „2 okna · 1 w strumieniu” |
| chwila ostatniej czynności | odpis obecności, przedstawiony względnie („przed chwilą”, „14 min temu”), a powyżej doby — datą |
| plakietka stanu | stan sesji z kontraktu: czynna, wstrzymana, zakończona, archiwalna |
| kropka sygnału | tętno przy sesji ze strumieniem, wypełnienie pełne przy sesji żywej na rdzeniu |

**Trzy stany uczciwe.** Strefa nigdy nie udaje danych, których nie ma:

| Stan | Kiedy | Co pokazuje |
|---|---|---|
| **oczekiwanie** | rdzeń nie odpowiedział jeszcze na odczyt wykazu | „Oczekiwanie na rdzeń” wraz z wyjaśnieniem, że wykaz pojawi się po odpowiedzi |
| **pusto** | rdzeń odpowiedział, a sesji w tle nie ma | „Brak sesji w tle” wraz z informacją, że sesja rozłączona pojawi się tu sama |
| **błąd** | rdzeń odmówił albo odpowiedź miała zły kształt | „Rdzeń odmówił wykazu sesji” wraz z treścią odmowy — stan pusty nie udaje wtedy braku sesji |

**Powrót do sesji.** Wiersz sesji ma przycisk „Wróć do sesji”, wykonujący dwie komendy kontraktu po kolei: **powiązanie połączenia z sesją**, a następnie **przeniesienie ogniska** na okno, które w tej sesji było ogniskowane. Powiązanie odtwarza okna sesji i kieruje ich strumienie na bieżące połączenie; od tej chwili koperty wychodzące niosą identyfikator sesji powiązanej. Odmowa przeniesienia ogniska **nie cofa powiązania** — sesja jest już związana, więc przerwanie powrotu byłoby stratą. Po udanym powrocie widok przechodzi do środowiska wskazanego przez odpis obecności, o ile kod tego środowiska należy do wykazu znanego interfejsowi.

Przycisk powrotu pojawia się **wyłącznie wtedy**, gdy połączenie ma tożsamość klienta nadaną w powitaniu. Bez niej wykaz pozostaje czysto informacyjny — przycisk, którego naciśnięcie nie miałoby skutku, byłby atrapą, a tych produkt nie dopuszcza.

**Znaczenie tej strefy.** Rozłączenie klienta nie kończy sesji ani jej procesów — rdzeń pracuje dalej. Bez tej strefy Operator, który zamknął przeglądarkę albo stracił łączność, nie miałby żadnego sposobu dowiedzenia się, że praca trwa, ani drogi powrotu do niej. Strefa sesji w tle jest tą drogą i jedyną w produkcie.

> **Zastrzeżenie.** Powrót do sesji przywraca **powiązanie z sesją i jej oknami**, a nie treść rozmowy. Historia widoczna w oknie po powrocie zaczyna się od chwili powiązania, ponieważ interfejs nie wczytuje wiadomości z bazy (rozdział [12.2](#122-czego-klient-nie-odtwarza)).

#### 3.2.4 Strefa czwarta — listwa ustawień

Pięć pozycji, wszystkie klikalne zawsze:

| Pozycja | Skutek | Stan |
|---|---|---|
| **Okno konfiguracji** | otwiera okno konfiguracji platformy | DZIAŁA |
| **Dostępy i katalog roboczy** | otwiera okno punktów dostępu, nadań i katalogu roboczego | DZIAŁA |
| **Modele, konta i tożsamość** | otwiera okno kont, ustawień osi, tożsamości i podglądu promptu | DZIAŁA |
| **Mobile** | dymek z informacją, że widok czeka na własny ekran | BRAK |
| **Always On Display** | jak wyżej | BRAK |

Dwie ostatnie pozycje nie są wyszarzone. Naciśnięcie zwraca komunikat mówiący wprost, że widok nie wchodzi w skład tej fali — zamiast otwierać ekran, którego nie ma.

### 3.3 Powłoka środowiska — cztery pasy

Widok środowiska zbudowany jest z czterech pasów ułożonych od góry.

#### 3.3.1 Pas 1 — pasek górny

Pasek na tle marki, o wysokości 56 pikseli, zawiera od lewej:

- **godło i nazwę produktu**;
- **kontekst pracy** — para „środowisko · moduł”, odświeżana przy każdym wyborze pozycji nawigacji;
- **pole poleceń** z podpowiedzią „Polecenie, moduł albo wyszukanie” oraz oznaczeniem skrótu `Ctrl K`;
- **grupę akcji** przy prawej krawędzi.

Grupa akcji zawiera, w kolejności: wskaźnik łączności z rdzeniem, przełącznik trzech widoków, uchwyt szuflady sterowania okna, dzwonek powiadomień z licznikiem, przełącznik Always On Display, przełącznik motywu i awatar Operatora z inicjałami.

**Stan faktyczny elementów paska:**

| Element | Stan |
|---|---|
| godło, nazwa produktu, kontekst „środowisko · moduł” | DZIAŁA |
| wskaźnik łączności (Rozłączony / Łączenie / Połączony / Ponawianie, z liczbą ramek w kolejce) | DZIAŁA |
| przełącznik tras | DZIAŁA |
| uchwyt szuflady sterowania | DZIAŁA |
| przełącznik motywu | DZIAŁA |
| **pole poleceń wraz ze skrótem `Ctrl K`** | ATRAPA — pole ma wygląd i ognisko, wykonanie polecenia należy do niezaimplementowanego okna „Command Center” — nazwa własna okna produktu, nie opis funkcji; skrót nie jest obsłużony |
| **dzwonek powiadomień** | ATRAPA — licznik jest sterowany metodą, którą powłoka wywołuje z wartością zero; nic go nie zasila |
| **przełącznik Always On Display** | ATRAPA — przełącza własny stan wciśnięcia i nic poza tym |
| **awatar Operatora** | ATRAPA — pokazuje inicjały napisu „Operator”; naciśnięcie nie otwiera menu |

#### 3.3.2 Pas 2 — karty sesji

Poziomy pas zakładek o mechanice ARIA. Powłoka otwiera się z jedną kartą.

**Obsługa:**

- **＋** przy prawej krawędzi pasa — zakłada kartę; tytuł nowej karty bierze się z nazwy modułu wybranego w nawigacji bocznej;
- kliknięcie karty — czyni ją czynną;
- **strzałki** w lewo/prawo/górę/dół — przenoszą wybór między kartami;
- **Home** i **End** — skaczą na krańce pasa;
- **Enter** oraz **spacja** — wybierają kartę pod ogniskiem;
- **Delete** — zamyka kartę pod ogniskiem;
- ikona zamknięcia na karcie — zamyka kartę.

Po zamknięciu karty czynnej wybór przechodzi na sąsiada z prawej, a przy jego braku z lewej.

**Stan faktyczny: DZIAŁA CZĘŚCIOWO — pas kart nie jest pasem sesji rdzenia.** Karta ma identyfikator wyłącznie miejscowy i nie odpowiada bytowi sesji po stronie rdzenia. Naciśnięcie ＋ nie zakłada nowej sesji, lecz wprowadza na scenę kolejne **okno komunikacji** — a ponieważ scena mieści najwyżej trzy okna, przy próbie założenia czwartej karty pojawia się komunikat „Scena jest pełna”. Zamknięcie karty zdejmuje okno ze sceny, ale nie zamyka okna w rdzeniu. Tytuł karty czynnej jest przemianowywany na nazwę modułu wybranego w nawigacji.

#### 3.3.3 Pas 3 — nawigacja boczna

Pionowa, stała kolumna zawierająca nagłówek środowiska (nazwa krojem szeryfowym, pod nią motto i liczba pozycji w odmianie polskiej) oraz wykaz pozycji. Każda pozycja ma ikonę, nazwę i dymek z opisem roli. Pozycja wybrana dostaje oznaczenie `aria-current` oraz pasek błękitu sygnałowego przy lewej krawędzi — barwa pochodzi z żetonu błękitu sygnałowego pakietu design v2.0 ([README.md](README.md) rozdz. 14), nie z nieobowiązującej warstwy zastępczej.

Wykazy pozycji per środowisko:

| Środowisko | Liczba pozycji | Pozycje |
|---|---|---|
| **TalkIn** | 9 modułów | Studio, Workspace, Browser, Research, Library, Translate, Roundtable, Assistant, Agents |
| **WorkSpace** | 9 modułów | Studio, Workspace, Browser, Research, Library, Roundtable, Design, Apps, Agents |
| **CodeStudio** | 8 modułów | Workspace, Roundtable, Design, Terminal, Developer, Diagnostics, Apps, Agents |
| **MultitaskingAI** | 6 sekcji | Zespoły, Role, Kolejki, Orkiestracja, Harmonogram i automatyki, Monitor procesu |

MultitaskingAI nie organizuje pracy Operatora, lecz zespół wykonawców, dlatego zamiast wykazu modułów niesie sekcje panelu orkiestracji. Nagłówek kolumny opisuje to poprawnie („Panel orkiestracji” zamiast „Moduły środowiska”).

Żadna pozycja nie traci klikalności — wykaz wynika ze środowiska, a nie z gotowości modułu. Skutek wyboru opisuje rozdział 4.2.

#### 3.3.4 Pas 4 — obszar roboczy

Kontener na widoki modułu. Dopóki nic nie jest w nim osadzone, pokazuje zapowiedź: ikonę modułu, nagłówek „Środowisko · Moduł” i zdanie „…Okna operacyjne modułu wejdą w ten obszar.”

**Stan faktyczny.** W widoku środowiska obszar roboczy zostaje **jednorazowo, przy montażu widoku, wypełniony sceną okien komunikacji** i pozostaje nią zajęty. Zapowiedź modułu jest wprawdzie odświeżana przy każdym wyborze pozycji nawigacji, lecz trafia do elementu, który został wtedy odłączony od dokumentu — Operator jej nie zobaczy.

### 3.4 Mission Control

Trzecia trasa, dostępna z przełącznika widoków. Pulpit orkiestracji pokazujący sesje, okna, kanały, kolejki ról, procesy w tle oraz przepływy oczekujące na decyzję.

**Stan faktyczny: DZIAŁA CZĘŚCIOWO — dane pochodzą z rdzenia, lecz kontrakt nie pokrywa wszystkich sekcji pulpitu.** Pulpit stoi **obok** Centrum dowodzenia, a nie zamiast niego: Centrum dowodzenia jest wejściem dla Operatora zaczynającego pracę, pulpit — dla wracającego do niej.

#### 3.4.1 Skąd pulpit bierze dane

Pulpit ma jedno źródło danych, prowadzące wyłącznie kanałem kontraktu. Odczyty wysyłane przy wejściu na trasę:

| Odczyt | Co przynosi |
|---|---|
| wykaz sesji, z żądaniem żywego odpisu obecności | sesje wraz ze środowiskiem, modułem, liczbą okien i chwilą ostatniej czynności |
| wykaz okien | okna komunikacji istniejące w rdzeniu |
| wykaz kanałów | rejestr kanałów modelu |

Subskrypcje utrzymywane przez cały czas obecności na trasie:

| Zdarzenie | Skutek |
|---|---|
| zmiana sesji | dopisanie, zmiana albo usunięcie sesji w stanie pulpitu |
| zmiana okna | to samo dla okna |
| zmiana kolejki | to samo dla kolejki |
| **zmiana postępu** | dopisanie albo odświeżenie procesu w wykazie procesów w tle |

Transport kolejkuje ramki do chwili zestawienia połączenia, więc odczyt wysłany przed otwarciem gniazda dochodzi po nim — pulpit nie musi nasłuchiwać stanu łącza, żeby zadziałać poprawnie po uruchomieniu.

#### 3.4.2 Stany uczciwe zamiast danych zastępczych

Sekcja bez danych pokazuje **stan pusty** z tytułem i zdaniem mówiącym wprost, czego brakuje — odczytu, który jeszcze nie nadszedł, albo źródła, którego kontrakt nie ma. Pulpit nie pokazuje ani jednej liczby, której nie oddał rdzeń.

Rozróżnienie, które warto znać przy czytaniu pulpitu:

- **przed pierwszym odczytem** sekcje pokazują stan oczekiwania na rdzeń;
- **po odczycie, bez danych** sekcje mówią, że rdzeń nie ma czego pokazać;
- **przy braku źródła w kontrakcie** sekcja mówi to wprost i nazywa brakującą komendę.

#### 3.4.3 Ograniczenie realne — brak odczytu wykazu kolejek

Kontrakt **nie zawiera odczytu wykazu kolejek**. Wykaz kolejek na pulpicie buduje się więc wyłącznie ze zdarzeń zmiany kolejki, a to znaczy, że **przed nadejściem pierwszego takiego zdarzenia sekcja pozostaje pusta**, nawet gdyby kolejki w rdzeniu istniały. Sekcja mówi to wprost, zamiast pokazywać zero udające pomiar. To samo ograniczenie dotyczy wykazu procesów w tle: nie ma odczytu bieżących procesów, więc wykaz zaczyna się od pierwszego zdarzenia postępu, jakie nadejdzie w czasie obecności na trasie.

Ograniczenie to jest **właściwością kontraktu**, nie usterką interfejsu — pulpit nie ma czego odpytać.

Osobno pozostaje w mocy ograniczenie opisane w rozdziale [9](#9-praca-z-modułami): **silnika wykonania kolejek nie ma**, więc pozycje kolejki nigdy nie powstają. Nawet komplet odczytów nie pokazałby zatem pracy kolejek, ponieważ takiej pracy się dziś nie prowadzi.

#### 3.4.4 Zamiary zgłaszane z pulpitu

Zamiary zgłaszane z pulpitu zachowują się różnie i jest to zachowanie zamierzone:

- **wejście do sesji** wysyła realną komendę otwarcia sesji i przenosi do środowiska;
- **działanie na kolejce** wysyła realną komendę działania kolejki — która jednak, po stronie rdzenia, wyłącznie przestawia kolumnę stanu kolejki; **silnika wykonania kolejek nie ma**, pozycje kolejki nigdy nie powstają;
- **podniesienie priorytetu** zwraca komunikat mówiący wprost, że czynność nie ma dziś odpowiednika w kontrakcie — ani w `QueueAction`, ani wśród komend (rozbieżność zgłoszona do rozstrzygnięcia);
- **przekazanie kontekstu** różni się od powyższego tym, że komenda **istnieje w kontrakcie** (`context.transfer`) — lecz pas kolejek pulpitu nie niesie kompletu `ContextBundle`, którego ta komenda wymaga. Pulpit zwraca więc komunikat, że przekazanie idzie z okna koordynatora, nie z pulpitu, i nie wysyła wywołania. Żadna z dwóch czynności nie jest udawana.

### 3.5 Motyw jasny i ciemny

Produkt ma dwa równoprawne motywy. Żaden nie jest wartością „domyślną” — przy braku wyboru Operatora rozstrzyga preferencja systemu operacyjnego, a jej zmiana działa na żywo, bez ponownego uruchomienia.

**Obsługa.** Przycisk motywu w grupie akcji paska górnego pokazuje motyw, w który przejdzie po naciśnięciu: słońce przy motywie ciemnym, księżyc przy jasnym. Wybór zapisuje się w pamięci trwałej przeglądarki pod kluczem `danaco-console.motyw` i przeżywa odświeżenie strony. Niedostępność pamięci trwałej nie odbiera możliwości przełączenia motywu w bieżącej sesji.

Warstwa wizualna oparta jest na żetonach systemu wizualnego v2.0, w gęstości zwartej, z zestawem **82 ikon** o jednolitym polu rysunku (`client/src/ikony/manifest.json`, pole `liczba-ikon`; godło marki — jedyny znak o barwach własnych — leży poza tym zestawem, jako plik odrębny). Stan elementu nigdy nie jest sygnalizowany samą barwą — plakietki niosą ikonę i etykietę, wskazanie pozycji nawigacji dublowane jest atrybutem dostępności.

> **Uwaga o stanie wdrożenia warstwy wizualnej.** Wdrożenie systemu v2.0 zostało w bibliotece komponentów niemal domknięte: prace naprawcze sprowadziły liczbę klas używanych w widokach bez definicji CSS z 36 do **1** — nie do zera. Pozostaje `.dn-karta-sesji`, użyta w `client/src/powloka/karty-sesji.ts:56` i niezdefiniowana w żadnym arkuszu; ustalone sprawdzianem pokrycia klas zastosowanym do kodu w stanie bieżącym (skrypt sprawdzianu, `budowa/scripts/klasy-css.mjs`, na tym zapisie jeszcze nie istnieje jako plik repozytorium — dopisuje go osobne, nieukończone zadanie; brakująca reguła jest więc luką w arkuszach stylu, niezależną od tego, kiedy powstał skrypt, który ją wykrywa). Poza tym jednym wyjątkiem domknięcie miało dwie drogi jednocześnie. Tam, gdzie widok adresował nazwę klasy ze słownika sprzed wymiany, powstał arkusz aliasów zgodności (`client/src/komponenty/aliasy-zgodnosci.css`), wczytywany jako **ostatni** w punkcie zbiorczym biblioteki, więc przy równej szczegółowości wygrywający z regułami bazowymi; jego mapowania są dosłownymi kopiami odpowiedników v2.0 — nie wprowadza ani jednej nowej barwy, ani jednego nowego wymiaru. Tam, gdzie brakowało reguł stylu samego komponentu — jak w przypadku stosu dymków komunikatów — pokrycie doszło wprost do słownika v2.0, a nie przez alias: **stos dymków ma dziś pełne reguły stylu w `client/src/komponenty/powiadomienie.css`** (klasa `dn-toasty` dla stosu, `dn-toast dn-toast--{informacja|sukces|ostrzezenie|blad}` dla dymka — nazwy klas bez polskich znaków diakrytycznych, dosłownie jak w arkuszu) — położenie stałe przy prawej dolnej krawędzi okna, własna warstwa stosu, układ kolumnowy z odstępem oraz obramowanie boczne i barwa ikony właściwe wadze dymka. Klasy `dn-toast-stos`, `dn-toast--ostrz`, `dn-toast--info` oraz reguła barwy ikony wewnątrz `.dn-toast` — nazwana w arkuszu aliasów zgodności klasą ikony, której słownik v2.0 nie deklaruje — czyli wszystko to, co arkusz aliasów zgodności również zawiera ponad wymienione wyżej reguły stosu, nie ma w repozytorium ani jednego użytkownika — są pozostałością bez zastosowania, nie ścieżką, którą stos faktycznie korzysta.
>
> **Stan tego rozwiązania: w większości alias zgodności do wycofania, nie brak pokrycia — z jednym wyjątkiem, który jest brakiem pokrycia.** Arkusz aliasów jest rozwiązaniem przejściowym i jest tak w kodzie opisany; dla klas, które przez niego przechodzą, dla Operatora skutek jest żaden — warstwa wizualna wygląda i zachowuje się zgodnie z systemem v2.0. Klasa `.dn-karta-sesji` (`client/src/powloka/karty-sesji.ts:56`) do tej kategorii nie należy: nie ma reguły ani w słowniku v2.0, ani w arkuszu aliasów — jest to więc rzeczywisty, pojedynczy brak pokrycia, nie tylko dług przejściowy. Docelowo widoki mają zostać przepięte na słownik v2.0, a arkusz skreślony — wraz z nim martwe reguły dymka powiadomienia, które nigdy nie miały odbiorcy. Do tego czasu obowiązuje reguła utrzymaniowa: zmiana pierwowzoru w bibliotece wymaga naniesienia tej samej zmiany w arkuszu aliasów — albo, lepiej, przepięcia widoku i usunięcia aliasu.

### 3.6 Zasada zera blokad i dymki komunikatów

Produkt konsekwentnie stosuje regułę: **żaden element nie jest wyszarzany ani pozbawiany klikalności**. Niegotowość opisuje się słowem, nie odebraniem możliwości działania. Reguła ta ma dwa praktyczne następstwa, które warto znać:

1. **Pozycja klikalna nie oznacza funkcji gotowej.** Wykaz modułów jest kompletny niezależnie od tego, czy moduł ma widok.
2. **Każde naciśnięcie musi dać odpowiedź.** Gdy czynność nie ma odpowiednika w kontrakcie albo rdzeń odmawia, pojawia się dymek mówiący wprost, co się stało i dlaczego. Dymek nigdy nie udaje wykonanej czynności.

Dymki mają cztery wagi: informacja, sukces, ostrzeżenie, błąd. Każda niesie własną ikonę, więc rozróżnienie nie zależy od koloru. Dymek błędu jest ogłaszany technologiom wspomagającym jako alert, pozostałe jako status. Czas życia dymka wynosi sześć sekund.

Podobną rolę pełni **pasek komunikatów w panelu sterowania okna**: potwierdza zmianę ustawienia albo zgłasza jej niepowodzenie wraz z kodem i treścią błędu z rdzenia.

---

## 4. Nawigacja

### 4.1 Wejście do środowiska

Zmiana środowiska prowadzi przez Centrum dowodzenia. Kolejność czynności przy wejściu: router pokazuje widok środowiska, następnie powłoka przestawia się na wskazane środowisko wraz z jego wykazem pozycji, a Centrum dowodzenia oznacza to środowisko jako czynne.

Do środowiska można wejść także wprost, przełącznikiem tras — środowisko domyślne (TalkIn) istnieje od pierwszej chwili i nie powstaje dopiero po wyborze na stronie głównej.

Przestawienie powłoki na inne środowisko wymienia wykaz pozycji nawigacji i **wybiera jego pierwszą pozycję**, co odświeża tytuł karty czynnej i kontekst na pasku górnym.

### 4.2 Wybór modułu — co robi, a czego nie robi

To jest miejsce, w którym stan bieżącej wersji odbiega od koncepcji najmocniej, dlatego opis jest jednoznaczny.

**Co wybór pozycji nawigacji robi:**

- oznacza pozycję jako bieżącą (pasek błękitu sygnałowego wg żetonu przy krawędzi oraz atrybut `aria-current`);
- przestawia atrybut środowiska na elemencie powłoki;
- wypisuje na pasku górnym parę „środowisko · moduł”;
- przemianowuje tytuł karty sesji czynnej na nazwę modułu.

**Czego wybór pozycji nawigacji NIE robi:**

- **nie otwiera żadnego okna operacyjnego** — obszar roboczy pozostaje zajęty przez scenę okien komunikacji;
- nie przeładowuje przestrzeni roboczej;
- nie zmienia kanału modelu, katalogów roboczych ani żadnego ustawienia okna;
- nie zgłasza rdzeniowi, w jakim module pracuje Operator.

W efekcie **wszystkie 32 pozycje nawigacji czterech środowisk prowadzą do tego samego ekranu**. Liczba ta rozkłada się następująco: 26 pozycji modułowych trzech środowisk modułowych (TalkIn 9, WorkSpace 9, CodeStudio 8) oraz 6 sekcji panelu orkiestracji środowiska MultitaskingAI — zgodnie z wykazami podanymi w tabeli rozdziału [3.3.3](#333-pas-3--nawigacja-boczna). Zmienia się wyłącznie tytuł karty i podpis kontekstu na pasku. Z siedemdziesięciu dziewięciu wierszy katalogu okien operacyjnych w bazie — liczby obowiązującej w całej dokumentacji produktu, ustanowionej w rozdziale [7](#7-środowiska) — w bieżącej wersji istnieje jeden widok: okno komunikacji, któremu odpowiada piętnaście wierszy katalogu, po jednym na moduł.

Rdzeń jest w tym zakresie gotowy: obsługuje komplet komend nawigacji platformy (wejście na stronę główną, wykaz i wejście do środowiska, wykaz modułów, wejście do przestrzeni roboczej, ognisko i powiązanie sesji, stan okna), a warstwa protokołu klienta je opakowuje. **Żaden widok ich jednak nie wywołuje** — nawigacja między środowiskami dzieje się wyłącznie po stronie przeglądarki.

### 4.3 Powrót i przełączanie widoków

Powrót do Centrum dowodzenia następuje przyciskiem „Centrum dowodzenia” w przełączniku tras. Widok środowiska nie jest przy tym niszczony: karty, okna, treść rozmów i stan sterowania trwają w tle i wracają w komplecie przy ponownym wejściu.

Ta sama zasada dotyczy zakładek wewnątrz okien ustawień: przełączenie zakładki nie porzuca obszaru, z którego Operator wychodzi — tekst wpisany w edytorze tożsamości przeżywa zajrzenie do kont.

### 4.4 Adres dokumentu i odświeżenie

Nazwa trasy zapisuje się w adresie dokumentu. Odświeżenie strony wraca do tego samego widoku, a nie na początek. Adres z nazwą nierozpoznaną nie zatrzymuje uruchomienia — router sprowadza go do trasy początkowej.

> **Ostrzeżenie eksploatacyjne.** Odświeżenie strony przywraca widok, lecz **nie przywraca treści rozmowy**. Historia okna budowana jest wyłącznie z tego, co przypłynęło w bieżącym połączeniu. Zagadnienie opisuje rozdział [12](#12-historia-i-dane).

---

## 5. Sesje i okna komunikacji

### 5.1 Pojęcia: sesja, karta sesji, okno komunikacji

Trzy pojęcia bywają mylone, a różnią się istotnie:

| Byt | Kto go tworzy | Do czego należy |
|---|---|---|
| **sesja** | rdzeń, przy uzgodnieniu z klientem | nadrzędny kontener; niesie identyfikator używany w każdej kopercie komunikatu |
| **karta sesji** | interfejs, w drugim pasie powłoki | zakładka wizualna; w bieżącej wersji odpowiada oknu na scenie, nie sesji rdzenia |
| **okno komunikacji** | rdzeń, na żądanie interfejsu | jednostka pracy: własny kanał modelu, moduł, katalogi robocze, środowisko wykonania, tryb uprawnień i rola |

Okno komunikacji jest bytem, do którego przypisane są **wszystkie** ustawienia wykonawcze. Dwa okna jednej sesji mogą pracować na dwóch różnych kanałach modelu, z dwoma różnymi zestawami katalogów i w dwóch różnych trybach uprawnień równocześnie.

### 5.2 Powstanie pierwszego okna

Po nawiązaniu połączenia klient przeprowadza uzgodnienie z rdzeniem: przedstawia się, zakłada sesję, po czym zakłada pierwsze okno komunikacji. Pierwsze gniazdo sceny dostaje to właśnie okno.

Domyślne parametry okna zakładanego przy uzgodnieniu:

| Parametr | Wartość domyślna | Źródło |
|---|---|---|
| projekt | „Danaco Console” | `VITE_PROJEKT` |
| tytuł | „Okno komunikacji” | `VITE_TYTUL_OKNA` |
| moduł | pierwsza pozycja wykazu kontraktu | `VITE_MODUL` |
| kanał modelu | pierwszy rodzaj kanału z wykazu kontraktu, czyli `cli` | `VITE_KANAL_MODELU` |
| katalogi robocze | lista pusta | `VITE_KATALOGI_ROBOCZE` |
| środowisko wykonania | urządzenie Operatora (`local`) | stała |
| tryb uprawnień | ręczny (pytanie przed każdą zmianą) | stała |
| rola | samodzielne | stała |

> **Uwaga krytyczna.** Domyślna wartość kanału modelu to **rodzaj** kanału, a nie identyfikator wiersza rejestru. Jest to jedna z trzech przyczyn, dla których świeża instalacja nie wykona tury (rozdział 2.7). Po założeniu wiersza kanału należy wskazać go w panelu sterowania okna, kontrolką „Model”.

### 5.3 Scena jednego, dwóch i trzech okien

Scena mieści od jednego do trzech okien obok siebie. Liczbę ustawia się na dwa sposoby:

1. **przełącznikiem „Okna komunikacji: 1 2 3”** na listwie nad sceną — grupa trzech zawsze czynnych pozycji jednokrotnego wyboru;
2. **przyciskiem ＋ w pasie kart sesji** — zakłada kartę i równocześnie wprowadza kolejne okno.

Wartość spoza zakresu jest przycinana do najbliższej dopuszczalnej, a nie odrzucana błędem. Próba wprowadzenia czwartego okna kończy się dymkiem „Scena jest pełna. Obok siebie mieszczą się najwyżej trzy okna komunikacji. Zamknij jedno, aby otworzyć kolejne.”

**Kiedy powstaje okno w rdzeniu.** Drugie i trzecie okno zamawiane jest komendą kontraktu dokładnie w chwili, gdy wchodzi na scenę — nie wcześniej. Zamówienie złożone przed zakończeniem uzgodnienia nie jest odrzucane: czeka i rusza zaraz po założeniu sesji.

Zejście okna ze sceny (zmniejszenie liczby albo zamknięcie karty) zdejmuje je z widoku. **Okno nie jest przy tym zamykane w rdzeniu** — interfejs nie wywołuje komendy zamknięcia okna ani sesji.

### 5.4 Role okien

Każde okno ma rolę w pętli koordynator–wykonawca. Rola widoczna jest jako plakietka w nagłówku gniazda, niosąca ikonę i etykietę wersalikową:

| Rola | Ikona | Znaczenie | Wygląd plakietki |
|---|---|---|---|
| **Koordynator** | waga | rozdziela i ocenia pracę | plakietka w błękicie sygnałowym wg żetonu |
| **Wykonawca** | strzałka uruchomienia | wykonuje pracę | plakietka informacyjna |
| **Samodzielne** | sylwetka | pracuje wprost z Operatorem, poza pętlą | plakietka neutralna |

**Role domyślne przy zmianie liczby okien:**

- jedno okno → samodzielne (pętla nie ma z kim biec);
- dwa okna → koordynator i wykonawca;
- trzy okna → koordynator i dwaj wykonawcy.

Rola nadana ręcznie — kontrolką „Rola okna” w panelu sterowania — **nie zostaje odebrana** przy kolejnej zmianie liczby okien. Domyślne są wartością wyjściową, nie bramą.

> **Uwaga.** Zmiana roli dokonana w panelu sterowania idzie komendą zmiany okna do rdzenia, lecz **nie wraca na scenę**: plakietka roli w nagłówku gniazda nie jest przez nią przestawiana. Figura koordynator–wykonawca widoczna na scenie może się więc rozejść ze stanem zapisanym w rdzeniu. Aby zobaczyć figurę zgodną, należy zmieniać liczbę okien przełącznikiem, który przelicza role całej sceny.

### 5.5 Więź koordynator–wykonawca i przekazanie zlecenia

Pod torem okien biegnie pas relacji pokazujący więź między koordynatorem a wykonawcą oraz stan pętli. Nagłówki okien objętych więzią niosą opis kierunku, a plakietka stanu pokazuje fazę pracy pary.

**Przekazanie zlecenia z okna koordynatora do okna wykonawcy — stan faktyczny: NIEZINTEGROWANE.** Czynność ta wykonuje wyłącznie skutki miejscowe: dopisuje wpisy systemowe do historii obu okien, mruga nagłówkami, puszcza żeton po pasie relacji i po 900 milisekundach przestawia stan na „wykonawca pracuje”. **Nie jest przy tym wysyłana żadna komenda kontraktu.** Rdzeń nie dowiaduje się o przekazaniu, wykonawca nie dostaje pracy, a zmiana stanu jest animacją.

Sama pętla koordynator–wykonawca **istnieje w rdzeniu** i jest uruchamiana turą rozmowy: rdzeń prowadzi licznik obiegów i emituje zdarzenie postępu.

**Odbiorca zdarzenia postępu w interfejsie istnieje — lecz nie w oknie komunikacji.** Zdarzenie postępu jest subskrybowane przez źródło danych pulpitu Mission Control (rozdział [3.4](#34-mission-control)) i zasila tam wykaz procesów w tle: każde nadejście zdarzenia dopisuje albo odświeża pozycję procesu, a sekcja aktywności liczy z nich procesy w biegu. **W oknie komunikacji odbiorcy nie ma** — ani nagłówek gniazda, ani pas relacji, ani plakietka stanu pary nie są tym zdarzeniem ożywiane. Praktyczny skutek dla Operatora prowadzącego rozmowę: licznik obiegów pętli i telemetria postępu **nie docierają do sceny okien**; aby je zobaczyć, trzeba przejść na trasę pulpitu. Ograniczeniem pozostaje przy tym brak odczytu bieżących procesów w kontrakcie — pulpit pokazuje wyłącznie te procesy, o których zdarzenie nadeszło w czasie obecności na trasie.

Widoczny w nagłówku gniazda stan pary jest ożywiany zdarzeniem zmiany kolejki, nie zdarzeniem postępu.

### 5.6 Prowadzenie rozmowy

Okno komunikacji składa się z historii wpisów oraz pola wysyłki u dołu.

**Pole wysyłki** zawiera obszar wpisywania (trzy wiersze), wskaźnik stanu, przycisk **Przerwij** i przycisk **Wyślij**.

**Klawiatura:**

| Klawisz | Skutek |
|---|---|
| `Enter` | wysyła wypowiedź |
| `Shift` + `Enter` | dokłada wiersz bez wysyłania |

Jest to zwyczaj konsoli, nie edytora. Treść pusta albo złożona z samych znaków odstępu nie jest wysyłana.

**Przebieg tury:**

1. Wypowiedź Operatora trafia natychmiast do historii, a pod nią powstaje wpis oczekujący na odpowiedź modelu.
2. Klient wysyła komendę wysłania wiadomości z żądaniem strumieniowania.
3. Rdzeń rozstrzyga kanał modelu okna, składa nakładkę tożsamości, buduje wiersz wywołania i uruchamia proces.
4. **Pierwszym fragmentem strumienia jest prowenancja** — nadawana zanim proces w ogóle wystartuje, więc dociera także wtedy, gdy uruchomienie się nie powiedzie.
5. Kolejne fragmenty narastają w treści wpisu; strumień domyka fragment ostatni.
6. Zdarzenie zmiany wiadomości potwierdza stan końcowy: zakończona, przerwana albo błędna.

**Przerwanie tury.** Przycisk **Przerwij** jest czynny zawsze, także gdy nic nie biegnie — rdzeń odpowiada wtedy informacją, że nic nie zatrzymano, a nie błędem. Przerwanie anuluje kontekst tury, co kończy proces modelu.

> **Ograniczenie znane.** Proces tury kanału głównego stoi poza nadzorem procesów sesji, więc zatrzymanie kończy proces główny, ale **nie ubija jego potomstwa**. Odczyt stanu procesu okna zwraca stale wartość „oczekujący” — podsystem procesów okien nie jest podłączony.

**Zerwanie połączenia w trakcie strumienia.** Rdzeń pracuje dalej — zerwanie połączenia klienta nie przerywa tury. Fragmenty wysłane w czasie rozłączenia przepadają jednak bezpowrotnie, a klient po ponownym połączeniu **nie odtwarza uzgodnienia ani powiązania sesji**. Praktyczny skutek: tekst wypowiedzi modelu może zostać trwale ucięty w miejscu zerwania.

**Kolejka wychodząca.** Treść wpisana przy braku łączności nie ginie — czeka w kolejce wychodzącej, a jej rozmiar widać przy wskaźniku łączności („Ponawianie · w kolejce 3”). Ponawianie połączenia jest wykładnicze, z pułapem piętnastu sekund i rozproszeniem losowym, i nie ma ograniczenia liczby prób.

### 5.7 Co pokazuje wpis rozmowy

Strumień odpowiedzi niesie dziewięć rodzajów fragmentów. Wpis rozmowy rozróżnia je i pokazuje w odrębnych blokach:

| Rodzaj fragmentu | Co pokazuje |
|---|---|
| **tekst** | treść odpowiedzi; narasta na żywo |
| **rozumowanie** | tok rozumowania modelu, w bloku zwijanym |
| **wywołanie narzędzia** | narzędzie wywołane przez model wraz z argumentami |
| **wynik narzędzia** | to, co narzędzie zwróciło modelowi |
| **obraz**, **dźwięk** | treść nietekstowa |
| **błąd** | błąd w trakcie strumienia; kończy wywołanie, nie sesję |
| **prowenancja** | pełny opis wywołania — blok diagnostyczny |
| **konto** | metadane konta użytego przez kanał |

**Blok prowenancji** jest najważniejszym narzędziem diagnostycznym produktu i odpowiada na pytanie „co dokładnie poszło do modelu”. Zawiera: plik wykonywalny kanału, pełny wiersz wywołania procesu, złożony prompt systemowy, tryb nakładki i nazwy jej warstw wraz ze skrótem, model główny i zapasowy, nakład rozumowania, tryb uprawnień, katalogi udostępnione, katalog roboczy procesu, plik ustawień, wskazane konfiguracje MCP, identyfikator wznawianej rozmowy, kod konta i jego katalog konfiguracji, numer próby i jej powód oraz chwilę złożenia.

Prowenancja jest przejrzystością, nie bramą — niczego nie dopuszcza i niczego nie wstrzymuje. Ładunek nieczytelny nie przerywa strumienia: fragment zostaje pokazany jako sam fakt wywołania, bez szczegółów.

> **Uwaga o fragmencie konta.** Rodzaj fragmentu „konto” istnieje w kontrakcie, rdzeń ma dla niego strukturę, a klient gotowy moduł odczytu — lecz **kanał główny go dziś nie nadaje**. Rotacja konta po wyczerpaniu limitu zachodzi więc milcząco; jej ślad widać wyłącznie w polach numeru i powodu próby w prowenancji.

### 5.8 Brak pamięci między turami — ograniczenie krytyczne

**To jest najważniejsze ograniczenie bieżącej wersji i należy je znać przed rozpoczęciem jakiejkolwiek pracy.**

**Model nie pamięta poprzednich wypowiedzi okna.** Każda tura uruchamia nowy proces programu `claude` z samą bieżącą wypowiedzią Operatora. Rozmowa nie ma ciągłości: pytanie „a co z tym drugim punktem?” nie ma dla modelu żadnego odniesienia, ponieważ model nie widział ani punktu pierwszego, ani drugiego.

Przyczyna techniczna: rdzeń odczytuje ze strumienia identyfikator rozmowy programu i pokazuje go w podsumowaniu tury, lecz **nigdzie go nie zapamiętuje i nigdy nie podaje w kolejnej turze** — przełącznik wznowienia nie pojawia się w wierszu wywołania. Pole wznowienia w zapytaniu kanału istnieje i ma odbiorcę, ale nie ma nadawcy.

**Skutki praktyczne, z którymi trzeba pracować:**

1. Każda wypowiedź musi być **samowystarczalna**. Kontekst, który ma być znany modelowi, trzeba powtórzyć w treści wypowiedzi.
2. Historia widoczna w oknie jest historią **dla Operatora**, nie dla modelu. To, co widać na ekranie, nie jest tym, co widzi model.
3. Praca wieloetapowa wymaga ręcznego przenoszenia ustaleń między turami — przez wklejenie streszczenia poprzedniego kroku do wypowiedzi otwierającej turę następną.
4. Trwałym nośnikiem kontekstu, który **działa**, jest nakładka tożsamości (konstytucja, profil roli, ekspertyza) — dociera do każdej tury, ponieważ jest składana od nowa przy każdym wywołaniu. Wiedzę mającą obowiązywać w całym cyklu pracy warto umieścić właśnie tam.

**[DO DECYZJI OPERATORA]** Miejsce trwałego zapamiętania identyfikatora rozmowy kanału: kolumna przy oknie komunikacji, ustawienie poziomu okna albo wpis w dzienniku rozmowy. Rozstrzygnięcie warunkuje domknięcie ciągłości rozmowy.

---

## 6. Modele, konta, dostawcy i kanały

### 6.1 Kanał modelu — pojęcie i rodzaje

**Kanał modelu** jest wierszem rejestru odpowiadającym na pytanie: przez co uruchamiany jest model. Kanał, a nie „model”, jest jednostką wyboru w interfejsie — kontrolka nosi wprawdzie nazwę „Model”, lecz wybiera się nią wiersz rejestru kanałów.

Rejestr jest sterowany danymi: nowy model to nowy wiersz rejestru, nie zmiana w kodzie. Rdzeń buduje rejestr z wierszy tabeli w czasie działania i przebudowuje go natychmiast po każdej zmianie wiersza, zachowując przy tym instancje kanałów niezmienionych.

#### Cztery rodzaje kanału

**Rodzaj kanału** jest wartością kolumny `rodzaj_kanalu` wiersza rejestru. Schemat bazy dopuszcza dokładnie **cztery** wartości, a wykaz kontraktu (`KnownChannelKinds`) niesie te same cztery. Wartość spoza tego zbioru **zostanie odrzucona przez więz bazy**, więc komenda zakładania kanału z rodzajem nieznanym po prostu się nie powiedzie.

| Rodzaj | Rola | Stan |
|---|---|---|
| **`cli`** — kanał główny | uruchamia zewnętrzny program `claude` (Claude Code CLI) w trybie strumieniowym i czyta strumień JSON-lines; obsługuje katalog konfiguracji per konto, rotację kont i rozpoznanie wyczerpania limitu | DZIAŁA |
| **`api`** — kanał sieciowy | generyczny kanał HTTP, w pełni sparametryzowany wierszem rejestru, bez zaszytego dostawcy; ścieżka wydobycia treści z odpowiedzi jest parametrem wiersza | NIEZINTEGROWANE z katalogiem kont — do wykonania dociera wyłącznie rodzaj `cli` |
| **`sdk`** — kanał zestawu narzędzi dostawcy | rodzaj przewidziany kontraktem i schematem bazy | BRAK — nie ma adaptera o tym kluczu, więc wiersz tego rodzaju bez wskazanego adaptera zostanie w rejestrze pominięty |
| **`lokalny`** — kanał bez sieci | rodzaj przeznaczony dla kanałów pracujących bez łącza i bez procesu zewnętrznego; to pod nim mieszka kanał testowy | DZIAŁA — pod warunkiem wskazania adaptera parametrem wiersza |

#### `echo` nie jest rodzajem kanału, lecz adapterem

Jest to rozróżnienie, którego pomylenie kończy się niewykonalnym poleceniem, dlatego opis jest szczegółowy.

**Rodzaj kanału** mówi, do jakiej klasy wiersz należy. **Adapter** mówi, który kod zbuduje z tego wiersza działający kanał. Rdzeń rozstrzyga adapter w dwóch krokach: najpierw czyta parametr `adapter` z pola `parametry_json` wiersza, a **dopiero w jego braku** sięga po kolumnę rodzaju. Obie wartości są danymi, nie nazwami typów w kodzie — dzięki temu kanał na nowym adapterze powstaje wpisem do tabeli, a nie zmianą w kodzie.

Adaptery wbudowane w rdzeń mają klucze: **`echo`**, **`api`** oraz **`cli`** (ten ostatni dokłada osobny pakiet wstrzykiwania procesu). Klucz `echo` **nie występuje** wśród rodzajów kanału — nie da się więc założyć kanału rodzaju `echo`.

**Kanał testowy powstaje jako wiersz rodzaju `lokalny` ze wskazanym adapterem `echo`:**

| Kolumna wiersza | Wartość |
|---|---|
| `rodzaj_kanalu` | `lokalny` |
| `parametry_json` | `{"adapter":"echo"}` |

Parametry rozpoznawane przez adapter `echo`, oba opcjonalne:

| Parametr | Znaczenie | Wartość domyślna |
|---|---|---|
| `adapter` | klucz adaptera; dla kanału testowego `echo` | — (bez niego rozstrzyga rodzaj kanału) |
| `porcja` | liczba znaków niesionych przez jeden fragment tekstu | 24 |
| `przedrostek` | tekst doklejany przed odpowiedzią | pusty |

Pełny przykład parametrów kanału testowego: `{"adapter":"echo","porcja":8,"przedrostek":"echo: "}`.

Adapter `echo` odsyła treść zapytania porcjami, zachowując pełną kolejność fragmentów kontraktu: najpierw prowenancja wywołania, następnie — o ile wiersz niesie odwołanie do poświadczenia albo zapytanie wskazuje konto — metadane konta, a dopiero po nich tekst odpowiedzi. Dzięki temu kanał testowy jest praktycznym narzędziem diagnostycznym: pozwala sprawdzić całą drogę od pola wysyłki po historię rozmowy, wraz z blokiem prowenancji, bez zależności od programu zewnętrznego i bez zużywania limitu konta.

> **Skutek operacyjny.** Polecenie „załóż kanał rodzaju `echo`” jest **niewykonalne** — komenda zakładania kanału z taką wartością rodzaju odpadnie na więzie bazy. Poprawne polecenie brzmi: „załóż kanał rodzaju `lokalny` z parametrem `adapter` równym `echo`”.

> **Skutek dla wiersza rejestru bez adaptera.** Wiersz rodzaju, dla którego nie ma ani wskazanego adaptera, ani adaptera wbudowanego o kluczu równym rodzajowi, zostaje w rejestrze **pominięty wraz z podaniem powodu** — pozostałe kanały pracują dalej. Rejestr nie gaśnie z powodu jednego wiersza nie do zbudowania.

Kanał zdalny po SSH **nie istnieje**. Program `ssh` występuje w produkcie wyłącznie jako polecenie zestawiające most MCP, a nie jako kanał wykonawczy modelu.

### 6.2 Rejestr kanałów i brak interfejsu do ich zakładania

**Stan faktyczny: BRAK UI ZAKŁADANIA KANAŁÓW.** Kontrakt zawiera komplet czterech komend zarządzania rejestrem — dodanie, zmiana, usunięcie i odczyt wykazu — a rdzeń obsługuje wszystkie cztery wraz z natychmiastowym odświeżeniem rejestru. **Główna aplikacja wywołuje wyłącznie odczyt wykazu.** Nie ma w produkcie ekranu, w którym Operator założyłby kanał modelu, zmienił jego parametry albo go usunął.

Jedno uściślenie: komenda zakładania kanału **jest** wywoływana w kodzie klienta — przez moduł automatycznego zapewnienia kanału głównego, opisany w rozdziale [2.7](#27-doprowadzenie-świeżej-instalacji-do-stanu-zdolnego-do-rozmowy). Moduł ten jest jednak wpięty wyłącznie do stanowiska podglądu rozmowy, więc dla Operatora pracującego z produktem droga ta pozostaje niedostępna. Zakładanie kanału pozostaje czynnością wykonywaną komendą kontraktu poza interfejsem.

Wykaz kanałów wykorzystywany jest w dwóch miejscach interfejsu:

1. kontrolka „Model” w panelu sterowania okna — lista wyboru kanału obsługującego okno;
2. podpowiedzi identyfikatorów modeli w oknie „Modele, konta i tożsamość”.

Kanał oznaczony jako nieczynny **zostaje na liście i pozostaje wybieralny** — o jego stanie mówi nazwa, a nie wyszarzenie.

> **Uwaga o identyfikatorze kanału.** Odczyt wykazu oddaje kanały z rejestru, w którym identyfikatorem jest numer wiersza, natomiast komendy zmiany i usunięcia przyjmują kod wiersza. Rozjazd ten należy znać przy pracy komendami kontraktu poza interfejsem.

### 6.3 Konta i pula rotacji

**Konto** jest nośnikiem tożsamości technicznej wobec dostawcy. Kontrakt zna trzy rodzaje: `api` (konto API modelu), `cli` (konto programu code CLI), `sdk` (konto SDK). Rodzaju nie zmienia się po założeniu konta.

Konta zakłada się i edytuje z interfejsu, w oknie „Modele, konta i tożsamość”, zakładka „Konta modeli i code CLI”. Panel dzieli się na wykaz kont po lewej i formularz jednego konta po prawej; pasek nad nimi zawiera ograniczenie wykazu do rodzaju, przycisk „Nowe konto” i przycisk „Odczytaj rejestr”.

**Pola formularza konta:**

| Pole | Znaczenie |
|---|---|
| Nazwa konta | nazwa widoczna w wykazie i w kanale modelu |
| Rodzaj konta | `api`, `cli` albo `sdk`; po założeniu pokazywany jako wydruk, nie kontrolka |
| Dostawca | wartość danych, nie typ kodu — wpisywana tekstem |
| Identyfikator konta u dostawcy | pole opcjonalne |
| Model domyślny konta | identyfikator modelu używany, gdy okno nie wskaże innego |
| Adres punktu końcowego | adres bazowy dostawcy |
| Katalog konfiguracji programu code CLI | katalog, w którym program trzyma własną konfigurację; ma sens wyłącznie dla rodzaju `cli` |
| Poświadczenie konta | pole tajne; wpisanie zapisuje nowe poświadczenie |
| Konto czynne | konto nieczynne zostaje w rejestrze i nie wchodzi do rotacji |
| Konto domyślne swojego rodzaju | domyślne jest dokładnie jedno konto na rodzaj; poprzednie traci oznaczenie |

**Zasada bezpieczeństwa.** Żadna komenda kontraktu nie zwraca poświadczenia. Pole poświadczenia jest zawsze puste, także dla konta, które poświadczenie ma. Baza przechowuje odwołania do danych dostępowych, nigdy ich treści; katalogi profili programu `claude` żyją na dysku, poza repozytorium i poza bazą.

**Pula rotacji kont** obejmuje konta rodzaju `cli`. Rdzeń przełącza konto po rozpoznaniu wyczerpania limitu i zapamiętuje chwilę odnowienia. Rozpoznanie ma trzy niezależne źródła sprawdzane w kolejności.

**Ograniczenia puli, które trzeba znać:**

| Ograniczenie | Skutek praktyczny |
|---|---|
| **Pula budowana jest raz, przy montażu rdzenia** | zmiany katalogu kont (dodanie, usunięcie, zmiana konta domyślnego) **wymagają restartu rdzenia**, aby weszły do rotacji |
| **Stan wyczerpania nie jest utrwalany** | restart rdzenia zeruje wiedzę o wyczerpanych kontach; kolumny stanu w bazie pozostają niewypełnione |
| **Pusta pula odmawia wywołania** | komunikat mówi o wyczerpaniu limitów, choć rzeczywistą przyczyną jest brak konta |
| **Konta rodzaju `api` i `sdk` nie mają drogi wykonawczej** | do wykonania dociera wyłącznie rodzaj `cli`, i to jako katalog konfiguracji na dysku |
| **Rotacja jest niewidoczna** | kanał nie nadaje fragmentu konta; przełączenie widać wyłącznie w prowenancji |
| **Konto preferowane kanału jest ignorowane** | wybór konta następuje wyłącznie przez rotację puli |

### 6.4 Co realnie steruje wywołaniem modelu

To zestawienie jest kluczowe dla zrozumienia, dlaczego część ustawień daje efekt, a część nie.

**Parametry, które DZIAŁAJĄ — docierają do procesu modelu:**

| Parametr | Droga do wykonania |
|---|---|
| **kanał modelu okna** | wybór wiersza rejestru → rozstrzygnięcie kanału → uruchomienie właściwego adaptera |
| **nakładka tożsamości** | konstytucja → profil roli → ekspertyza, składane wg krytyczności, z trybem ZASTĄP albo DOŁĄCZ → przełącznik promptu systemowego |
| **tryb uprawnień okna** | wartość okna → przełącznik trybu uprawnień w wierszu wywołania |
| **katalogi robocze okna** | lista okna → wielokrotny przełącznik dodania katalogu |
| **katalog roboczy sesji** | ustawienia `katalog.roboczy.*` → katalog startowy procesu |
| **mosty MCP z nadań okna** | punkty dostępu i nadania → wygenerowana konfiguracja MCP → przełącznik konfiguracji MCP |
| **konto rodzaju `cli`** | katalog konfiguracji konta → zmienna środowiska procesu modelu |

**Parametry, które są ZAPISYWANE BEZ SKUTKU:**

| Parametr | Dlaczego nie działa |
|---|---|
| **wybór modelu głównego** (poza wyborem kanału) | pole modelu w zapytaniu kanału nie ma nadawcy; przełącznik modelu w wierszu wywołania pochodzi wyłącznie z parametrów wiersza rejestru |
| **model zapasowy** | interfejs zapisuje pod kluczem, którego rdzeń nie zna; pole zapytania bez nadawcy |
| **nakład rozumowania** | j.w.; dodatkowo interfejs zapisuje liczbę 1–5, a katalog ustawień przewiduje wyliczenie słowne |
| **host wykonania** | j.w.; wartość nie ma konsumenta wykonawczego |
| **środowisko wykonania** (`local` / `core` / `remote`) | wartość jest wyłącznie etykietą w prowenancji; proces zawsze startuje na maszynie rdzenia |
| **plik ustawień powłoki** | pole bez nadawcy |
| **zmienne środowiska procesu modelu** (poza katalogiem konfiguracji konta) | pole bez nadawcy |
| **jedenaście punktów izolacji** | klucze zdefiniowane i rozstrzygane, bez ani jednego egzekutora w czasie wykonania |
| **klucze kategorii `harness`** | katalog je zna, kod ich nie odczytuje |

Przyczyną większości powyższych jest jeden fakt techniczny: **rozstrzygacz ośmiu poziomów zasięgu nie jest źródłem parametrów wywołania modelu**. Zapytanie kanału składane jest z pól wiersza okna, a nie z rozstrzygniętej konfiguracji. Wyjątkiem są dwa miejsca, w których rozstrzygacz jest realnie wywoływany: katalog roboczy sesji oraz tryb nakładki tożsamości — i one działają.

Dodatkowo dwa z ośmiu poziomów zasięgu są w praktyce osiągalne (okno i globalny). Kontekst rozstrzygania nigdy nie niesie środowiska, modułu, pary modułów, projektu, karty sesji ani roli, więc ustawienie zapisane na tych poziomach nie wygra — poziom nie wchodzi do rachunku.

---

## 7. Środowiska

Cztery środowiska są w bieżącej wersji **profilami widoczności nawigacji**, a nie odrębnymi przestrzeniami pracy o odmiennych funkcjach.

**Co działa:**

- karty czterech środowisk na stronie głównej wraz z godłami, mottami i opisami;
- przestawienie powłoki na wskazane środowisko;
- wykaz pozycji nawigacji właściwy środowisku (9 / 9 / 8 modułów oraz 6 sekcji orkiestracji);
- nagłówek kolumny nawigacji: nazwa, motto, liczba pozycji w poprawnej odmianie polskiej;
- kontekst „środowisko · moduł” na pasku górnym.

**Czego nie ma:**

- **realizacja funkcjonalna środowisk** — żadne z czterech nie ma własnych okien operacyjnych ani własnego zestawu czynności;
- **MultitaskingAI** jest sześcioma etykietami listy; nie ma panelu orkiestracji, okien ról, sieci podagentów, hierarchii zespołów, harmonogramu ani monitora procesu — wszystkie te pozycje mają stan BRAK;
- **przypisanie sesji do środowiska nie jest utrwalane.** Każda utrwalona sesja trafia do pierwszego środowiska (TalkIn), pod stałą nazwą karty. Praca prowadzona w WorkSpace, CodeStudio albo MultitaskingAI zapisze się jako praca w TalkIn i takim wróci po restarcie rdzenia;
- **rdzeń nie wie, w jakim środowisku pracuje Operator** — komendy nawigacji platformy nie są wywoływane.

Rdzeń ma po swojej stronie komplet: cztery środowiska, piętnaście modułów, macierz widoczności modułów w środowiskach oraz katalog okien operacyjnych — wszystko zasiane migracją i wystawione komendami. **Interfejs tego katalogu nie czyta** i prowadzi własny, zaszyty w kodzie wykaz środowisk i pozycji. Oba wykazy są dziś zgodne co do treści, lecz są dwoma źródłami tej samej wiedzy.

> **Uściślenie liczbowe — inwentarz a katalog w bazie.** W dokumentacji produktu występują dwie różne liczby okien operacyjnych i nie są one ze sobą sprzeczne, lecz mierzą dwie różne rzeczy.
>
> **Prototypy interfejsu.** Katalog `design/05-okna/` niesie **38 prototypów** okien produktu (zliczenie: `find design/05-okna -name '*.html' | wc -l`), z czego piętnaście w podkatalogu `design/05-okna/moduly/` — jeden na moduł. Dawny inwentarz repozytorium budowy, który wyliczał okna operacyjne w rozbiciu na jednostki unikatowe i wystąpienia w module, nie należy do bieżącego drzewa repozytorium; liczbą obowiązującą w całej dokumentacji produktu, w tym w [README.md](README.md), jest odtąd liczba wierszy katalogu w bazie, podana niżej. Do niej odnoszą się rozdziały [4.2](#42-wybór-modułu--co-robi-a-czego-nie-robi), [9](#9-praca-z-modułami) i [18](#18-zbiorcze-zestawienie-stanu-funkcji) niniejszej instrukcji.
>
> **Katalog w bazie** zasiewany jest natomiast **wierszami**, nie oknami unikatowymi: migracja wstawia jeden wiersz Chat Window na każdy z piętnastu modułów oraz 64 wiersze pozostałych okien, co daje **79 wierszy** tabeli okien operacyjnych. Liczba wierszy nie jest liczbą okien produktu i nie należy jej tak czytać.
>
> Niezależnie od obu liczb obowiązuje stan faktyczny: **w bieżącej wersji istnieje jedno okno operacyjne — okno komunikacji.** Wszystkie pozostałe pozycje katalogu są zapisami w bazie bez odpowiadającego im widoku.

> **Rozjazd pojęciowy, który warto znać.** Wykaz „znanych identyfikatorów modułów” w kontrakcie niesie w rzeczywistości **kody środowisk**. Z tego wykazu bierze się domyślny moduł okna oraz zawartość kontrolki „Moduł okna” w panelu sterowania — dlatego kontrolka ta oferuje cztery środowiska, a nie piętnaście modułów. Rdzeń dodatkowo zamienia kod nierozpoznany na moduł `workspace`, przez co okno zapisuje się zawsze z tym samym modułem, niezależnie od nawigacji.

---

## 8. Konfiguracja sesji i okna — panel sterowania

### 8.1 Umiejscowienie i obsługa szuflady

Komplet sterowania powstaje **osobno dla każdego okna komunikacji** i domyka się na jego identyfikatorze. Dwa okna obok siebie mają dwa niezależne komplety, bez jednej wspólnej zmiennej. Wspólny pozostaje wyłącznie wykaz kanałów modelu — jest katalogiem wyboru, nie ustawieniem okna.

Kolumna sterowania montuje się obok sceny okna. Zawiera nagłówek z identyfikatorem okna, szufladę oraz pasek komunikatów.

**Szuflada** ma dwie warstwy: podsumowanie ośmiu wartości (widoczne po zwinięciu) i komplet kontrolek (widoczny po rozwinięciu). Szuflada **startuje rozwinięta** — Operator ma zobaczyć ustawienia od razu, bez szukania, co nacisnąć.

**Obsługa szuflady:**

- uchwyt szuflady w kolumnie sterowania — przełącza stan; napis zmienia się między „Ustawienia okna” a „Rozwiń ustawienia okna”;
- uchwyt w grupie akcji paska górnego — ten sam stan, ogłaszany magistralą, więc oba uchwyty pokazują tę samą prawdę;
- klawisz `Escape` wewnątrz szuflady — zwija ją i wraca ogniskiem na uchwyt;
- naciśnięcie wartości w podsumowaniu — rozwija szufladę i prowadzi ognisko do kompletu kontrolek.

Zwinięcie nie jest blokadą: żaden element nie traci klikalności, zmienia się wyłącznie to, która warstwa zajmuje miejsce.

**Podsumowanie ośmiu wartości.** Warstwa zwinięta pokazuje osiem wierszy, każdy z ikoną, nazwą ustawienia i jego bieżącą wartością:

| Wiersz | Wartość pokazywana | Wartość zastępcza przy braku |
|---|---|---|
| Środowisko wykonania | środowisko, a po kropce host wykonania | — |
| Katalogi robocze | wykaz katalogów okna | informacja o braku katalogu |
| Model | nazwa kanału modelu z rejestru | „Bez wskazania” |
| Model zapasowy | nazwa kanału z rejestru | „Bez modelu zapasowego” |
| Nakład rozumowania | nazwa stopnia | — |
| Tryb uprawnień | nazwa trybu | — |
| Rola okna | nazwa roli | — |
| Moduł | nazwa modułu okna | — |

Dwa wiersze — Model i Model zapasowy — mają **wskaźnik oczekiwania**: dopóki wykaz kanałów nie przypłynął z rdzenia, a okno ma wskazany kanał, wiersz pokazuje, że nazwa jest w drodze, zamiast orzekać, że kanału nie ma. Rozróżnienie „odpowiedź jeszcze nie nadeszła” od „rejestr jest pusty” nie jest dziś możliwe po stronie interfejsu, więc wskaźnik jest jedynym uczciwym rozwiązaniem.

**Objaśnienia kontekstowe [?].** Przy nazwie każdego z ośmiu wierszy stoi znak zapytania niosący jednozdaniowe objaśnienie ustawienia — czym ono steruje i co znaczy jego wartość. Objaśnienie pokazuje się przy najechaniu wskaźnikiem **oraz przy ognisku klawiatury**, bez konieczności kliknięcia i bez osobnego zamykania. Znak jest przyciskiem, nie ozdobą: naciśnięcie prowadzi do niego ognisko, a ognisko pokazuje treść — każde naciśnięcie daje więc odpowiedź. Treść objaśnienia czytają także technologie wspomagające, ponieważ jest ona etykietą dostępności samego znaku.

Naciśnięcie **nazwy** wiersza (nie znaku zapytania) rozwija szufladę i prowadzi ognisko do odpowiadającej kontrolki w komplecie sterowań — podsumowanie jest więc równocześnie spisem treści panelu.

### 8.2 Dziewięć kontrolek

Komplet zawiera dziewięć sterowań w jednej siatce:

| Kontrolka | Rodzaj | Droga zapisu |
|---|---|---|
| **Środowisko wykonania** | lista wyboru | komenda zmiany okna |
| **Host wykonania** | pole tekstowe | zapis ustawienia na poziomie okna |
| **Moduł okna** | lista wyboru (cztery pozycje z kontraktu) | komenda zmiany okna |
| **Model** (kanał modelu) | lista wyboru z rejestru kanałów | komenda zmiany okna |
| **Model zapasowy** | lista wyboru z rejestru kanałów | zapis ustawienia na poziomie okna |
| **Nakład rozumowania** | suwak 1–5 | zapis ustawienia na poziomie okna |
| **Tryb uprawnień** | lista wyboru | komenda zmiany okna |
| **Rola okna** | lista wyboru | komenda zmiany okna |
| **Katalogi robocze** | lista z dodawaniem i usuwaniem | komenda zmiany okna |

**Nazwy stopni nakładu rozumowania:** 1 — „Najszybciej — odpowiedź bez namysłu”, 2 — „Szybciej”, 3 — „Równowaga”, 4 — „Mądrzej”, 5 — „Najmądrzej — pełny namysł”.

**Tryby uprawnień** dostępne w kontrolce: ręczny (pytanie przed każdą zmianą), akceptuj zmiany plików, plan (bez zmian w systemie), automatyczny (decyduje model), bez pytania z zachowaniem ograniczeń, pominięcie kontroli uprawnień.

**Katalogi robocze** dodaje się i usuwa po jednym; każda operacja wysyła pełną listę po zmianie, ponieważ kontrakt niesie katalogi jako całość. Lista pusta nie jest błędem — okno pracuje wtedy bez wskazanego katalogu.

**Potwierdzanie zmian.** Zmiana idąca komendą zmiany okna jest potwierdzana dwiema drogami: odpowiedzią na komendę oraz zdarzeniem zmiany okna rozgłaszanym do wszystkich urządzeń. Zdarzenie dotyczące innego okna nie zmienia tego kompletu — filtr po identyfikatorze jest adresowaniem, nie bramą. Niepowodzenie zmiany **nie blokuje sterowania**: widok wraca do stanu potwierdzonego, a pasek komunikatów podaje treść i kod błędu.

### 8.3 Zestawienie skuteczności ustawień okna

Poniższa tabela jest praktycznym podsumowaniem rozdziału 6.4, odniesionym do konkretnych kontrolek.

| Kontrolka | Zapis | Wpływ na wywołanie modelu |
|---|---|---|
| Model (kanał modelu) | DZIAŁA | **DZIAŁA** — rozstrzyga, czym uruchomiony zostanie model |
| Tryb uprawnień | DZIAŁA | **DZIAŁA** — trafia do wiersza wywołania |
| Katalogi robocze | DZIAŁA | **DZIAŁA** — trafiają do wiersza wywołania |
| Rola okna | DZIAŁA (rdzeń) | **DZIAŁA** dla pętli w rdzeniu; nie wraca na scenę |
| Moduł okna | DZIAŁA CZĘŚCIOWO | brak wpływu na wywołanie; rdzeń podmienia kod nierozpoznany |
| Środowisko wykonania | DZIAŁA | **BEZ SKUTKU** — wyłącznie etykieta w prowenancji |
| Host wykonania | NIEZINTEGROWANE — zapis bez skutku | brak — klucz nieznany rdzeniowi |
| Model zapasowy | NIEZINTEGROWANE — zapis bez skutku | brak — klucz nieznany rdzeniowi |
| Nakład rozumowania | NIEZINTEGROWANE — zapis bez skutku | brak — klucz nieznany rdzeniowi, typ wartości niezgodny |

Wartość ustawienia zapisanego bez skutku **da się odczytać w prowenancji**: pola modelu, modelu zapasowego i nakładu w bloku prowenancji pozostają puste, co jest bezpośrednim potwierdzeniem, że parametr nie dojechał do wywołania.

**[DO DECYZJI OPERATORA]** Miejsce, w którym mieszka słownik kluczy ustawień okna: kontrakt (`shared/`) czy katalog ustawień w bazie. Rozstrzygnięcie jest warunkiem uzgodnienia trzech kluczy zapisywanych dziś przez interfejs pod nazwami nieznanymi rdzeniowi.

---

## 9. Praca z modułami

Rozdział ten opisuje stan, a nie procedurę — ponieważ procedury pracy z modułem w bieżącej wersji nie ma.

**Stan faktyczny.** Wszystkie piętnaście modułów produktu — Studio, Workspace, Browser, Research, Library, Translate, Roundtable, Design, Assistant, Terminal, Developer, Diagnostics, Apps, Agents, Automations — występuje w produkcie **wyłącznie jako pozycje nawigacji**. Żaden nie ma okna operacyjnego, żaden nie ma własnych czynności i żaden nie zmienia zawartości obszaru roboczego.

Z siedemdziesięciu dziewięciu wierszy katalogu okien operacyjnych w bazie (rozdział [7](#7-środowiska)) istnieje jeden widok: **okno komunikacji**, wspólne wszystkim modułom, któremu odpowiada piętnaście wierszy katalogu. Nie istnieją w szczególności: edytor Studia, panel narzędzi, porównanie różnic, eksplorator biblioteki, zakładki terminala, panel repozytorium, tablica projektowa, panel debaty, menedżer kolejek ani kreator agentów.

Praktyczny wniosek dla Operatora: **całość pracy w bieżącej wersji odbywa się w oknie komunikacji**. Zróżnicowanie sposobu pracy między modułami osiąga się nie przez wybór modułu w nawigacji, lecz przez:

1. **nakładkę tożsamości** — profil roli i ekspertyza zapisane dla wskazanego modelu albo konta;
2. **katalogi robocze okna** — zakres plików, w które model ma wgląd;
3. **nadania dostępu okna** — maszyny osiągalne mostem MCP;
4. **tryb uprawnień okna** — swoboda modelu w dokonywaniu zmian.

Moduł **Automations** nie występuje w żadnym wykazie nawigacji — zgodnie z koncepcją jego wejście miało prowadzić ze strony głównej. Kafel Automations na stronie głównej prowadzi dziś do środowiska WorkSpace i pozycji Workspace, więc **modułu Automations nie da się osiągnąć poprawnym wejściem**.

**[DO DECYZJI OPERATORA]** Droga otwarcia okien modułu Automations oraz sposób osadzania okien operacyjnych w obszarze roboczym: rozszerzenie kontraktu o komendy modułowe albo jedna komenda czynności okna. Rozstrzygnięcie warunkuje realizację całej warstwy modułów.

---

## 10. Obsługa narzędzi — mosty MCP, katalogi, nadania

### 10.1 Okno „Dostępy i katalog roboczy”

Okno otwiera się z listwy ustawień Centrum dowodzenia. Jest oknem nakładkowym opartym na natywnym mechanizmie przeglądarki, więc zamyka się klawiszem `Escape`, a warstwa tła i stos okien pochodzą z przeglądarki.

Okno rozróżnia **trzy byty, których nie wolno mieszać** i mówi o tym wprost w nagłówku:

| Byt | Odpowiada na pytanie | Gdzie żyje |
|---|---|---|
| **punkt dostępu** wraz z **nadaniem** | do jakich maszyn i katalogów model sięga | punkt — platforma; nadanie — okno rozmowy |
| **katalog roboczy** | gdzie model zostawia swoje pliki | ustawienie warstwowe |
| **środowisko** | profil widoczności modułów w nawigacji | poza tym oknem — sekcja go nie dotyka |

Układ okna: po lewej wykaz punktów dostępu, po prawej zbiór nadań okna, pod nimi — wyraźnie oddzielony — obszar katalogu roboczego. Przycisk „Odczytaj ponownie” w nagłówku odświeża wszystkie trzy obszary.

### 10.2 Punkty dostępu

Punkt dostępu jest bytem platformy — opisuje maszynę albo katalog raz, dla wszystkich okien. Kontrakt zna dwa rodzaje:

| Rodzaj | Nazwa w interfejsie | Stan |
|---|---|---|
| **most MCP** | „Maszyna (most MCP)” | DZIAŁA |
| **katalog lokalny** | „Katalog lokalny” | **NIEZAPISYWALNY** — schemat wymaga wskazania urządzenia, a tabela urządzeń nie ma żadnej drogi zapisu |

Karta punktu w wykazie zawiera: ikonę rodzaju, nazwę, plakietkę rodzaju, plakietkę stanu ostatniego sprawdzenia, oznaczenie „nadany temu oknu” (o ile nadanie istnieje), adres, chwilę ostatniego sprawdzenia, przycisk **Sprawdź**, przycisk usunięcia oraz formularz nadania.

**Stany punktu:** „odpowiada”, „nie odpowiada”, „niesprawdzony”. Stan nierozpoznany dostaje plakietkę neutralną — brak wiedzy nie jest błędem.

**Przycisk Sprawdź** wysyła zapytanie o osiągalność punktu; wynik trafia do plakietki stanu i do zdania obok.

**Dodanie katalogu z „Mój komputer”.** W powłoce natywnej dostępne jest natywne okno wyboru katalogu — zwraca ścieżkę istniejącą, rozwiniętą i zapisaną w konwencji systemu. Interfejs otwarty w przeglądarce nie ma powłoki i nie ma natywnego okna: wskazanie wraca wtedy pustą wartością, dokładnie tak samo jak rezygnacja Operatora, a widok zostawia drogę wpisania ścieżki ręcznie.

### 10.3 Nadania dostępu per okno

**Nadanie jest bytem okna, nie platformy.** Ta sama maszyna może być nadana jednemu oknu w trybie odczytu, drugiemu w trybie zapisu, a trzeciemu wcale.

Formularz nadania na karcie punktu zawiera:

- **przełącznik trybu** — odczyt albo zapis (kolejność rosnącego uprawnienia);
- **wybór korzeni** — katalogi, poza które nadanie nie wychodzi;
- **przycisk „Nadaj oknu dostęp”**.

**Jak nadanie dociera do modelu.** Rdzeń składa z nadań okna konfigurację serwerów MCP i przekazuje ją procesowi modelu przełącznikiem konfiguracji MCP. Wpis mostu ma postać połączenia w trybie standardowego wejścia i wyjścia, zestawianego programem `ssh` do skryptu uruchomieniowego mostu po stronie maszyny. Tryb nadania wchodzi argumentem pozycyjnym, a korzenie — zmienną środowiska rozdzieloną dwukropkiem.

Cechy istotne w eksploatacji:

- **tryb bierze się z nadania, nie z ustawienia globalnego** — nadanie żyje per okno rozmowy;
- **słownictwo trybu należy do konkretnego mostu** i pochodzi z wiersza w bazie, nie z kodu; tryb pusty uruchamia sam skrypt, co po stronie mostu oznacza wariant bezpieczniejszy, czyli odczyt;
- **nadania i punkty wygaszone są pomijane**, tak samo jak punkty rodzaju innego niż maszyna;
- **brak nadań nie jest błędem** — daje pustą, lecz poprawną konfigurację.

Wymagania po stronie maszyny: dostępny program `ssh` po stronie rdzenia oraz skrypt uruchomieniowy mostu w ustalonej ścieżce po stronie maszyny. Wartości domyślne połączenia (nazwa konta, port) obowiązują, gdy adres punktu ich nie niesie.

### 10.4 Katalog roboczy modelu

Katalog roboczy jest **zwykłym ustawieniem** — idzie tymi samymi komendami i przez ten sam rozstrzygacz co reszta konfiguracji. Nie ma dla niego osobnej komendy i mieć nie powinien.

Obszar w oknie dostępów pokazuje dwa klucze:

| Klucz | Znaczenie | Wartość domyślna |
|---|---|---|
| `katalog.roboczy.podstawa` | katalog, w którym powstają katalogi sesyjne modelu | ustalana przez rdzeń przy pierwszym uruchomieniu, w miejscu instalacji aplikacji |
| `katalog.roboczy.wzorzec_sesji` | wzorzec nazwy katalogu jednej sesji wewnątrz podstawy | `sesje/<identyfikator>` |

Znacznik `<identyfikator>` zastępowany jest identyfikatorem sesji. Pole podstawy ma przycisk wyboru katalogu (w powłoce natywnej); pole wzorca jest polem tekstowym.

Interfejs **nie zgaduje ścieżki domyślnej**: mówi wprost, że katalog tworzy rdzeń, i prosi o wpisanie ścieżki tylko wtedy, gdy Operator chce to nadpisać. Pełny wybór ośmiu poziomów zasięgu i trzech osi dla tych kluczy pozostaje w oknie konfiguracji — obszar w oknie dostępów jest skrótem, nie drugą implementacją.

**Rozróżnienie, którego nie wolno zatrzeć:** dostęp mówi, **do czego** model sięga; katalog roboczy mówi, **gdzie** model zostawia swoje pliki. Model może mieć wgląd w cudzy katalog i nic w nim nie zapisywać, a swoje katalogi sesyjne trzymać w miejscu instalacji. Zlanie tych ustawień w jedno kazałoby otworzyć zapis wszędzie tam, gdzie model ma tylko czytać.

### 10.5 Narzędzia platformy dla modelu

Koncepcja produktu przewiduje sterowanie platformą przez model: komendy kontraktu udostępniane kanałowi modelu jako narzędzia, bez osobnego parsera intencji — polecenie wpisane w oknie rozmowy wykonuje model wywołaniem narzędzia.

**Stan faktyczny: BRAK.** Generator wytwarza dziś trzydzieści dziewięć deklaracji narzędzi po obu stronach kontraktu. **Żaden plik poza wygenerowanym kontraktem ich nie odczytuje** — ani adapter kanału modelu, ani warstwa uruchamiania procesu, ani klient. Nie ma również ścieżki zwrotnej przyjmującej wywołanie narzędzia od modelu.

Praktyczny wniosek: **model nie może dziś sterować platformą**. Nie założy okna, nie zmieni ustawienia, nie otworzy sesji i nie odczyta wykazu kanałów. Narzędzia, które model realnie ma do dyspozycji, pochodzą wyłącznie z programu `claude` oraz z mostów MCP nadanych oknu.

Powiązany brak: **panel czynności** (rejestr akcji sterowany danymi) istnieje po stronie rdzenia wraz z zaczynem katalogu i komendą odczytu, lecz interfejs nie wywołuje go ani razu i nie ma panelu czynności w oknie komunikacji.

---

## 11. Praca z plikami i artefaktami

### 11.1 Jak model uzyskuje dostęp do plików

W bieżącej wersji istnieje **jedna działająca droga** udostępnienia modelowi plików lokalnych: **katalogi robocze okna**. Lista katalogów okna trafia do wiersza wywołania procesu jako wielokrotny przełącznik dodania katalogu, więc model widzi dokładnie te katalogi, które w niej stoją.

Kolejność czynności:

1. Otwórz szufladę sterowania okna.
2. W kontrolce **Katalogi robocze** wpisz ścieżkę katalogu i zatwierdź przyciskiem dodania.
3. Sprawdź w bloku prowenancji następnej tury, czy katalog pojawił się w polu katalogów udostępnionych i w wierszu wywołania.

Katalog usuwa się przyciskiem przy jego pozycji. Każda operacja wysyła pełną listę po zmianie.

Drugą drogą — działającą, lecz dotyczącą maszyn zdalnych, nie plików lokalnych — są **mosty MCP nadane oknu** (rozdział 10.3). Trzecia, katalog lokalny jako punkt dostępu, **nie jest dziś zapisywalna** (rozdział 10.2).

### 11.2 Gdzie powstają pliki modelu

Katalog roboczy procesu ustala rdzeń na podstawie ustawień `katalog.roboczy.podstawa` i `katalog.roboczy.wzorzec_sesji` (rozdział 10.4). Jest to jedno z niewielu ustawień, które rzeczywiście przechodzi przez rozstrzygacz poziomów zasięgu i dociera do procesu jako katalog startowy.

Katalog roboczy widoczny jest w bloku prowenancji jako osobne pole, obok listy katalogów udostępnionych. Rozróżnienie tych dwóch pól jest istotne przy diagnozie: model może mieć wgląd w katalogi, a pracować w zupełnie innym miejscu.

### 11.3 Artefakty i treści obszerne

Decyzja architektoniczna produktu przewiduje, że baza przechowuje ścieżkę, a treść obszerna trafia do pliku w katalogu danych. **Część plikowa tej decyzji nie jest zaimplementowana.**

Stan faktyczny:

- kolumny odwołania do treści istnieją w schemacie (przy wiadomości, przy zasobie pamięci, przy pozycji kolejki) i mają w repozytoriach pełną drogę zapisu i odczytu, lecz **nigdy nie są wypełniane**;
- **nie ma ani jednego zapisu treści do pliku**;
- nie istnieją tabele załączników wiadomości ani artefaktów;
- w interfejsie nie ma zatem przeglądarki artefaktów, wersjonowania ani eksportu treści do pliku.

Praktyczny wniosek: **cała treść rozmowy przechowywana jest w bazie jako tekst kolumny wiadomości**. Produkt nie zarządza artefaktami jako osobnym bytem. Pliki wytworzone przez model powstają wyłącznie na dysku, w katalogu roboczym procesu, i pozostają poza wiedzą platformy — Danaco Console ich nie indeksuje, nie wersjonuje i nie pokazuje.

**[DO DECYZJI OPERATORA]** Miejsce składowania treści obszernych w katalogu danych oraz kształt tabel załącznika wiadomości i artefaktu. Rozstrzygnięcie warunkuje domknięcie plikowej części decyzji o trwałości.

---

## 12. Historia i dane

### 12.1 Co trafia do bazy

Rdzeń utrwala każdą wypowiedź Operatora i każdą odpowiedź modelu w tabeli wiadomości, w komplecie więzów: środowisko → karta sesji → sesja → okno komunikacji → wiadomość. Baza jest źródłem prawdy; bufor pamięci procesu wchodzi wyłącznie jako jawna degradacja per okno, gdy utrwalenie się nie powiedzie.

Trwałość działa również dla:

- ustawień warstwowych (poziom zasięgu wraz z osią);
- katalogu kont wraz z ich parametrami;
- punktów dostępu i nadań;
- kategorii i dokumentów tożsamości modelu;
- rejestru kanałów modelu.

### 12.2 Czego klient nie odtwarza

**Ograniczenie kluczowe: interfejs nigdy nie wczytuje historii rozmowy.** Kontrakt zawiera komendę odczytu wiadomości okna, rdzeń ją obsługuje i oddaje historię z bazy — **klient nie wywołuje jej ani razu**.

Historia widoczna w oknie budowana jest wyłącznie z tego, co przypłynęło w bieżącym połączeniu: z fragmentów strumienia i ze zdarzeń zmiany wiadomości. Praktyczne skutki:

| Zdarzenie | Skutek dla widocznej historii |
|---|---|
| **odświeżenie strony** | historia znika z widoku, choć pozostaje w bazie |
| **ponowne połączenie po zerwaniu** | historia nie wraca; klient nie powtarza uzgodnienia ani powiązania sesji |
| **powrót do środowiska po wyjściu do Centrum dowodzenia** | historia **zostaje** — widok nie jest niszczony |
| **restart rdzenia** | historia znika z widoku; sesje i okna wracają w rdzeniu, wiadomości pozostają w bazie |

**Czego interfejs nie wywołuje — stan ścisły.** W całym kodzie klienta **nie ma ani jednego odwołania do komendy odczytu wiadomości okna**; jest to potwierdzone przeglądem katalogu klienta. Nie są też wywoływane komendy zamknięcia sesji ani zamknięcia okna — zejście okna ze sceny zdejmuje je z widoku, lecz zostawia w rdzeniu (rozdział [5.3](#53-scena-jednego-dwóch-i-trzech-okien)).

**Co interfejs wywołuje.** Odczyt **wykazu sesji** jest natomiast wywoływany, i to w dwóch miejscach, oba z żądaniem żywego odpisu obecności:

| Miejsce | Do czego służy |
|---|---|
| strefa „Sesje w tle” Centrum dowodzenia | wykaz sesji trwających na rdzeniu poza bieżącym połączeniem wraz z drogą powrotu do nich (rozdział [3.2.3](#323-strefa-trzecia--sesje-w-tle)) |
| źródło danych pulpitu Mission Control | sekcje sesji, matrycy i aktywności pulpitu (rozdział [3.4](#34-mission-control)) |

Rozróżnienie jest istotne praktycznie: Operator **widzi**, że sesja trwa, i **może do niej wrócić**, lecz po powrocie nie zobaczy tego, co w niej wcześniej powiedziano. Wykaz sesji odtwarza dostęp do pracy, nie jej treść. Drogi obejrzenia historii rozmowy poza oknem — przeglądarki archiwum, wyszukiwania w wiadomościach, eksportu — w produkcie nie ma.

### 12.3 Co przeżywa restart rdzenia

| Byt | Przeżywa restart | Uwaga |
|---|---|---|
| **wiadomości rozmowy** | TAK | zapisane w tabeli wiadomości |
| **sesja i okno komunikacji** | TAK, **warunkowo** | wiersz sesji i okna powstaje dopiero przy pierwszej utrwalonej wiadomości; sesja bez wypowiedzi nie zostawia śladu |
| **ustawienia warstwowe** | TAK | |
| **konta, punkty dostępu, nadania, tożsamość, kanały** | TAK | |
| **zmiany ustawień okna dokonane komendą zmiany okna** | NIE | komenda zmienia byt w pamięci nadzorcy; droga zapisu do bazy istnieje, lecz nie jest wywoływana |
| **przypisanie sesji do środowiska** | NIE | każda utrwalona sesja trafia do pierwszego środowiska pod stałą nazwą karty |
| **stan wyczerpania kont w puli rotacji** | NIE | nieutrwalany |
| **więź koordynator–wykonawca** | NIE | pole koordynatora okna nigdy nie jest zapisywane |
| **stan procesów okien** | NIE | podsystem procesów okien nie jest podłączony |

Rdzeń po restarcie odtwarza z bazy sesje i okna pod tymi samymi identyfikatorami, więc odwołania do nich pozostają ważne.

---

## 13. Ustawienia

Produkt ma trzy okna ustawień, wszystkie otwierane z listwy Centrum dowodzenia. Każde jest oknem nakładkowym przeglądarki: zamyka się klawiszem `Escape`, ma własny nagłówek z ikoną i przyciskiem zamknięcia oraz stopkę z przyciskiem odczytu ponownego i zamknięcia. Każde otwiera się natychmiast, **przed odpowiedzią rdzenia** — dane dojeżdżają do niego odpowiedzią, a brak danych zostawia komunikat, nie pusty prostokąt.

### 13.1 Okno konfiguracji

Okno konfiguracji jest głównym narzędziem konfiguracyjnym platformy. **Nie zawiera ani jednej listy pól zaszytej w kodzie** — kategorie, definicje pól, rodzaje kontrolek, opcje i dopuszczalne poziomy pochodzą z katalogu w bazie rdzenia.

**Układ okna:**

1. **pasek punktu widzenia** u góry — ustala miejsce, względem którego liczone jest dziedziczenie każdego pola;
2. **kolumna kategorii** po lewej;
3. **formularz kategorii czynnej** po prawej;
4. **stopka** z przyciskiem „Odczytaj katalog ponownie” i „Zamknij”.

**Dwa prostopadłe wymiary konfiguracji.**

**Poziom zasięgu** mówi, **jak wąsko** obowiązuje wartość. Osiem poziomów, w kolejności od najwęższego (wygrywającego) do najszerszego:

| Kolejność | Poziom | Byt wskazywany przy poziomie |
|---|---|---|
| 1 | okno komunikacji | identyfikator okna |
| 2 | rola | nazwa roli |
| 3 | karta sesji | identyfikator sesji |
| 4 | projekt | identyfikator projektu |
| 5 | para modułów | para kodów modułów |
| 6 | moduł | kod modułu |
| 7 | środowisko | kod środowiska |
| 8 | globalny | — |

**Oś** mówi, **dla czego** wartość obowiązuje. Trzy osie, od najwęższej: konto (identyfikator konta), model (identyfikator modelu), platforma (bez bytu). Oś pominięta znaczy platformę.

Byt poziomu i byt osi **znikają tam, gdzie nie mają sensu**: poziom globalny bytu nie ma, oś platformy również. Znikające pole nie jest blokadą — jest usunięciem kontrolki bez znaczenia.

**Jak zapisać wartość.** Przy każdym polu stoi własny wybór adresu zapisu — poziom, byt poziomu, oś, byt osi — zawężony do poziomów i osi dopuszczonych przez katalog dla tego pola. Pasek u góry ustawia natomiast punkt widzenia dla całego formularza.

**Stan faktyczny okna konfiguracji:**

| Element | Stan |
|---|---|
| katalog kategorii i pól z rdzenia | DZIAŁA |
| zapis wartości pod adresem poziom + oś | DZIAŁA |
| wskaźnik dziedziczenia i pochodzenia wartości | DZIAŁA CZĘŚCIOWO — odczyt bez poziomu liczy politykę zawsze w kontekście poziomu okna, więc podgląd łańcucha dziedziczenia nie odzwierciedla pełnej drabiny |
| skuteczność zapisanych wartości w wykonaniu | patrz rozdział 6.4 — **większość kluczy katalogu nie ma dziś konsumenta wykonawczego** |
| konfiguracja per środowisko i per moduł | DZIAŁA CZĘŚCIOWO — kod bytu wpisuje się ręcznie, bez podpowiedzi z katalogu |
| okno konfiguracji punktów izolacji | BRAK — brak widoku i brak komend w kontrakcie |

Warto zapamiętać, które kategorie ustawień mają realny skutek: **katalog roboczy** oraz **tryb domyślny nakładki tożsamości**. Pozostałe zapisują się poprawnie, lecz wykonanie ich dziś nie odczytuje.

### 13.2 Okno „Modele, konta i tożsamość”

Okno zbiera cztery obszary związane jedną osią. Pasek osi stoi nad zakładkami, ponieważ wszystkie obszary mówią o tym samym bycie: o wskazanym modelu albo o wskazanym koncie.

| Zakładka | Odpowiada na pytanie | Stan |
|---|---|---|
| **Konta modeli i code CLI** | czym się łączymy | DZIAŁA |
| **Ustawienia osi** | jak model ma działać | zapis DZIAŁA; skutek wykonawczy — patrz 6.4 |
| **Zasady i tożsamość** | kim model ma być | DZIAŁA |
| **Podgląd złożonego promptu** | co z tego naprawdę pojedzie | DZIAŁA |

**Pasek osi** nie dotyczy zakładki kont — rejestr kont jest wspólny dla osi, więc pasek na niej znika. Byt osi jest polem otwartym: podpowiedzi pochodzą z rejestru kanałów i z kont, lecz identyfikator modelu jeszcze nieużywanego wolno wpisać wprost.

**Zakładka „Ustawienia osi”** korzysta w całości z mechanizmów okna konfiguracji — tego samego katalogu kategorii, tych samych kontrolek i tego samego zapisu. Dokłada jedną rzecz: punkt widzenia zawężony do wskazanej osi. Zmiana osi przebudowuje formularz, a nie tylko odświeża wartości — formularz zbudowany dla modelu A zapisywałby inaczej pod adres modelu B.

**Zakładka „Zasady i tożsamość”** dzieli się na wykaz kategorii po lewej i edytor jednej kategorii po prawej. Kategorie pochodzą z katalogu rdzenia. Treść zapisana w kategorii składa się na nakładkę systemową wskazanej osi.

**Nakładka tożsamości — trzy warstwy w kolejności krytyczności:**

| Warstwa | Nazwa w interfejsie | Rola |
|---|---|---|
| 1 | **konstytucja** | zasady najwyższej wagi, obowiązujące bezwarunkowo |
| 2 | **profil roli** | kim model jest w tej pracy |
| 3 | **ekspertyza** | wiedza dziedzinowa |

**Dwa tryby podania nakładki:**

| Tryb | Skutek | Ostrzeżenie pokazywane w interfejsie |
|---|---|---|
| **ZASTĄP** (domyślny) | podmienia prompt fabryczny w całości | „ZASTĄPIENIE ZDEJMUJE PROMPT FABRYCZNY W CAŁOŚCI. Model dostanie wyłącznie treść złożoną z kategorii tej sekcji…” |
| **DOŁĄCZ** | dokłada treść do promptu fabrycznego | „DOŁĄCZENIE ZOSTAWIA PROMPT FABRYCZNY…” |

Rozróżnienie nie jest kosmetyczne — decyduje o tym, którym przełącznikiem nakładka pojedzie do procesu modelu. Jest to jeden z **niewielu mechanizmów konfiguracyjnych produktu działających od interfejsu do wykonania**.

**Zakładka „Podgląd złożonego promptu”** pokazuje nakładkę obowiązującą — tę, którą naprawdę dostanie model. Zawiera plakietkę trybu, odcisk (skrót) nakładki, wykaz warstw w kolejności złożenia, imienne wyliczenie kategorii wymaganych pozostawionych bez treści oraz przycisk skopiowania złożonego promptu. Podgląd **nie skleja treści po stronie interfejsu** — pobiera nakładkę obowiązującą z rdzenia, więc pokazuje to samo, co pojedzie do procesu.

Bez tego podglądu konfiguracja tożsamości byłaby ślepa: Operator zapisuje treść w kilkunastu kategoriach na trzech osiach, a pytanie „co z tego złożyło się w jeden prompt i czy prompt fabryczny został zastąpiony” musi mieć odpowiedź przed uruchomieniem okna.

### 13.3 Okno dostępów

Opisane w rozdziale [10](#10-obsługa-narzędzi--mosty-mcp-katalogi-nadania). Okno to jest równocześnie miejscem, w którym ustawia się katalog roboczy modelu.

### 13.4 Pozycje listwy bez własnego ekranu

Pozycje **Mobile** i **Always On Display** listwy ustawień nie mają widoków. Naciśnięcie zwraca dymek z wyjaśnieniem pozycji i informacją, że widok czeka na własny ekran. Jest to zachowanie zamierzone: reguła zera blokad zakazuje wyszarzania, a zakaz atrap zabrania otwierania pustego ekranu udającego funkcję.

Poza tymi dwiema pozycjami, koncepcja przewiduje w oknie ustawień sześć sekcji; zrealizowane są trzy (konfiguracja, dostępy, modele/konta/tożsamość).

---

## 14. Typowe scenariusze pracy

### 14.1 Od uruchomienia do pierwszej rozmowy

Scenariusz zakłada świeżą instalację i uwzględnia obejścia opisane w rozdziale 2.7.

| Krok | Czynność | Miejsce |
|---|---|---|
| 1 | Zbuduj rdzeń i pakiet interfejsu | wiersz poleceń, katalogi `budowa` i `budowa\client` |
| 2 | Uruchom powłokę albo sam rdzeń | rozdziały 2.2 / 2.3 |
| 3 | Sprawdź, że wskaźnik łączności na pasku pokazuje **Połączony** | pasek górny widoku środowiska |
| 4 | Załóż wiersz kanału głównego rodzaju `cli` | komenda kontraktu poza interfejsem — **brak drogi z głównej aplikacji** (rozdział [2.7](#27-doprowadzenie-świeżej-instalacji-do-stanu-zdolnego-do-rozmowy)) |
| 5 | Załóż konto rodzaju `cli` wskazujące katalog konfiguracji programu | Centrum dowodzenia → „Modele, konta i tożsamość” → „Konta modeli i code CLI” → „Nowe konto” |
| 6 | Zrestartuj rdzeń, aby konto weszło do puli rotacji | menu zasobnika: „Zatrzymaj rdzeń”, następnie ponowne uruchomienie |
| 7 | Wejdź do środowiska — TalkIn, WorkSpace, CodeStudio albo MultitaskingAI | karta na stronie głównej |
| 8 | Rozwiń szufladę sterowania i w kontrolce **Model** wskaż założony kanał | kolumna sterowania obok okna |
| 9 | Ustaw tryb uprawnień odpowiedni do zamierzonej pracy | kontrolka „Tryb uprawnień” |
| 10 | Dodaj katalogi robocze, w które model ma mieć wgląd | kontrolka „Katalogi robocze” |
| 11 | Wpisz wypowiedź i naciśnij `Enter` | pole wysyłki okna |
| 12 | Rozwiń blok prowenancji pierwszej odpowiedzi i sprawdź wiersz wywołania | historia okna |

Krok 12 jest krokiem kontrolnym, nie ozdobnym: prowenancja jest jedynym miejscem, w którym widać, co naprawdę poszło do modelu — z jakim promptem systemowym, z jakim kontem, z jakimi katalogami i w jakim katalogu roboczym.

**Jeżeli tura kończy się błędem**, przejdź do rozdziału [15](#15-komunikaty-i-diagnostyka).

**Wariant sprawdzenia bez programu zewnętrznego.** Jeżeli celem jest wyłącznie potwierdzenie, że droga strumienia działa, załóż **kanał testowy** i wskaż go w kontrolce „Model”. Kanał testowy odsyła wypowiedź porcjami, bez procesu zewnętrznego i bez zużycia limitu konta, zachowując przy tym pełną kolejność fragmentów — z blokiem prowenancji na czele.

Kanał testowy zakłada się jako wiersz **rodzaju `lokalny`** ze wskazanym **adapterem `echo`**:

| Pole wiersza rejestru | Wartość |
|---|---|
| rodzaj kanału | `lokalny` |
| parametry (`parametry_json`) | `{"adapter":"echo"}` — opcjonalnie z `porcja` i `przedrostek` |

> **Uwaga.** Nie należy zakładać kanału „rodzaju `echo`” — takiego rodzaju nie ma. Komenda zakładania kanału z tą wartością **odpadnie na więzie bazy**, ponieważ schemat dopuszcza wyłącznie rodzaje `cli`, `api`, `sdk` i `lokalny`. Pojęcia rodzaju kanału i adaptera rozróżnia rozdział [6.1](#61-kanał-modelu--pojęcie-i-rodzaje).

Krok 6 procedury (restart rdzenia) jest przy tym wariancie **zbędny**: kanał testowy nie korzysta z puli rotacji kont, więc pusta pula go nie zatrzyma. Wariant ten jest zatem najkrótszą drogą potwierdzenia, że rdzeń, transport, interfejs i zapis rozmowy do bazy współpracują poprawnie — jeszcze zanim Operator zajmie się kontami i programem zewnętrznym.

### 14.2 Praca z tożsamością modelu

Nakładka tożsamości jest **najskuteczniejszym narzędziem sterowania zachowaniem modelu** w bieżącej wersji — działa na całej drodze od interfejsu do procesu, i to przy każdej turze.

| Krok | Czynność |
|---|---|
| 1 | Otwórz „Modele, konta i tożsamość” z listwy Centrum dowodzenia |
| 2 | Na pasku osi wybierz oś: **platforma** (warstwa wspólna), **model** albo **konto**, i wskaż byt osi |
| 3 | Przejdź na zakładkę „Zasady i tożsamość” |
| 4 | Wybierz kategorię w wykazie po lewej i wpisz treść w edytorze |
| 5 | Ustal tryb podania: **ZASTĄP** albo **DOŁĄCZ** — przeczytaj ostrzeżenie pod przełącznikiem |
| 6 | Przejdź na zakładkę „Podgląd złożonego promptu” i sprawdź: plakietkę trybu, kolejność warstw, wykaz kategorii wymaganych bez treści |
| 7 | Wróć do okna komunikacji i wykonaj turę |
| 8 | Rozwiń blok prowenancji i porównaj prompt systemowy oraz nazwy warstw ze złożonym podglądem |

**Zasady praktyczne:**

- **Oś platformy jest tłem.** Zapis dla modelu albo konta leży na niej; odjęcie osi platformy odebrałoby jedyne miejsce, w którym widać i zmienia się warstwę wspólną.
- **ZASTĄP zdejmuje prompt fabryczny w całości.** Wszystko, co ma obowiązywać, musi wtedy stać w kategoriach tożsamości. Jest to tryb domyślny.
- **Kategoria wymagana bez treści nie wstrzymuje uruchomienia**, lecz jest imiennie wyliczana w podglądzie — nie należy tego przeoczyć.
- **Nakładka jest jedynym trwałym nośnikiem kontekstu**, ponieważ model nie pamięta poprzednich tur (rozdział 5.8). Wiedzę mającą obowiązywać w całym cyklu pracy warto umieścić właśnie tutaj, a nie w pierwszej wypowiedzi rozmowy.

### 14.3 Praca z wieloma oknami

| Krok | Czynność |
|---|---|
| 1 | Na listwie nad sceną ustaw przełącznikiem liczbę okien: **2** albo **3** |
| 2 | Sprawdź plakietki ról: pierwsze okno zostaje koordynatorem, pozostałe wykonawcami |
| 3 | Dla każdego okna rozwiń jego własną kolumnę sterowania i ustaw kanał modelu, tryb uprawnień oraz katalogi robocze |
| 4 | Prowadź rozmowy w oknach niezależnie — każde ma własną historię, własny strumień i własne przerwanie |

**Co warto wiedzieć:**

- **Każde okno jest niezależne.** Dwa okna jednej sesji mogą pracować na dwóch różnych kanałach modelu, w dwóch trybach uprawnień i z dwoma zestawami katalogów równocześnie.
- **Nie ma automatycznego przekazania pracy między oknami.** Przekazanie zlecenia z okna koordynatora do wykonawcy jest w bieżącej wersji animacją bez komendy kontraktu (rozdział 5.5). Współpracę okien prowadzi się ręcznie: przenosząc treść między nimi.
- **Zmiana roli w panelu sterowania nie wraca na scenę.** Aby figura koordynator–wykonawca była widocznie spójna, zmieniaj liczbę okien przełącznikiem, który przelicza role całej sceny.
- **Scena mieści najwyżej trzy okna.** Środowisko MultitaskingAI przewiduje w koncepcji cztery okna ról — takiego układu nie ma.
- **Zejście okna ze sceny nie zamyka go w rdzeniu.**

### 14.4 Nadanie modelowi dostępu do maszyny

| Krok | Czynność |
|---|---|
| 1 | Otwórz „Dostępy i katalog roboczy” z listwy Centrum dowodzenia |
| 2 | W wykazie punktów odszukaj maszynę albo dodaj punkt rodzaju „Maszyna (most MCP)” |
| 3 | Naciśnij **Sprawdź** i potwierdź, że plakietka stanu pokazuje „odpowiada” |
| 4 | W formularzu nadania ustaw tryb: **odczyt** albo **zapis** |
| 5 | Wskaż korzenie, poza które nadanie nie wychodzi |
| 6 | Naciśnij **Nadaj oknu dostęp** — nadanie dotyczy okna, z którym związana jest sekcja |
| 7 | Wróć do okna komunikacji i wykonaj turę |
| 8 | W bloku prowenancji sprawdź pole wskazanych konfiguracji MCP |

**Zasady praktyczne:**

- **Tryb należy do nadania, nie do punktu.** Ta sama maszyna może być nadana jednemu oknu w trybie odczytu, a drugiemu w trybie zapisu.
- **Korzenie są ograniczeniem zasięgu.** Most nie wyjdzie poza wskazane katalogi.
- **Brak nadań nie jest błędem** — okno pracuje wtedy bez mostów.
- **Punkt rodzaju „Katalog lokalny” nie jest dziś zapisywalny.** Do udostępnienia katalogów lokalnych służą katalogi robocze okna (rozdział 11.1).

---

## 15. Komunikaty i diagnostyka

### 15.1 Najczęstsze komunikaty i ich przyczyny

| Komunikat / objaw | Przyczyna | Postępowanie |
|---|---|---|
| **„kanał %s nie istnieje w rejestrze albo jest nieczynny”** | okno wskazuje kanał, którego nie ma w rejestrze — najczęściej dlatego, że wskazanie domyślne niesie rodzaj `cli` zamiast identyfikatora wiersza, albo rejestr jest pusty | załóż wiersz kanału (rozdział 2.7), następnie wskaż go w kontrolce „Model” |
| **„okno %s bez wskazanego kanału modelu”** | pole kanału okna jest puste; okno założone drogą wejścia do przestrzeni roboczej nie dostaje kanału | wskaż kanał w kontrolce „Model” panelu sterowania |
| **„wszystkie konta puli mają wyczerpany limit; kolejne próby pozostają otwarte”** | pula kont rodzaju `cli` jest pusta **albo** wszystkie konta rozpoznano jako wyczerpane; komunikat nie rozróżnia tych przypadków | sprawdź, czy istnieje choć jedno czynne konto rodzaju `cli`; po dodaniu konta **zrestartuj rdzeń** |
| **odpowiedź `*.unknown`** | rdzeń albo klient otrzymał nazwę komendy lub zdarzenia spoza kontraktu; nazwa nierozpoznana **nie zrywa kanału** — zwracana jest odpowiedź o nieznanej komendzie i praca trwa dalej | sprawdź zgodność wersji klienta i rdzenia; nieznane nazwy trafiają po stronie klienta do dziennika nieznanych |
| **„Scena jest pełna”** | próba wprowadzenia czwartego okna komunikacji | zamknij jedno okno albo zmniejsz liczbę przełącznikiem |
| **wskaźnik łączności „Ponawianie · w kolejce N”** | rdzeń nie odpowiada; ramki czekają w kolejce wychodzącej | sprawdź, czy rdzeń nasłuchuje na porcie 17870 i czy klient nie został zbudowany z innym adresem |
| **okno interfejsu puste, rdzeń pracuje** | brak katalogu pakietu interfejsu; nieistniejący katalog nie wstrzymuje nasłuchu | zbuduj pakiet klienta i wskaż go przełącznikiem `--klient` albo zmienną katalogu klienta |
| **komunikat o niepowodzeniu zmiany ustawienia w pasku panelu sterowania** | rdzeń odrzucił komendę zmiany okna | odczytaj kod i treść błędu z paska; widok wrócił już do stanu potwierdzonego, można zmieniać dalej |
| **odpowiedź modelu urwana w połowie** | zerwanie połączenia w trakcie strumienia albo przepełnienie kolejki wyjściowej; fragment przepada bezpowrotnie, a klient nie dociąga treści | powtórz turę; treść pozostaje kompletna w bazie, lecz interfejs jej nie odtwarza |
| **model „nie pamięta” poprzedniej wypowiedzi** | zachowanie zamierzone w tej wersji — każda tura startuje bez historii | patrz rozdział 5.8 |

### 15.2 Gdzie szukać przyczyn

| Źródło | Zawartość | Ścieżka |
|---|---|---|
| **blok prowenancji w oknie rozmowy** | pełny wiersz wywołania, prompt systemowy, konto, katalogi, katalog roboczy, konfiguracje MCP, numer i powód próby | pierwszy fragment każdej odpowiedzi |
| **dziennik rdzenia prowadzonego przez powłokę** | wyjście standardowe i błędów rdzenia oraz wpisy powłoki o starcie i źródle interfejsu | `<katalog danych>\rdzen-powloki.log` |
| **opis stanu rdzenia** | czy rdzeń pracuje, identyfikator procesu, adres, ścieżka dziennika, przyczyna niepowodzenia | menu zasobnika → „Stan rdzenia” |
| **wskaźnik łączności** | stan transportu i liczba ramek oczekujących | pasek górny widoku środowiska |
| **pasek komunikatów panelu sterowania** | potwierdzenia i niepowodzenia zmian ustawień okna | kolumna sterowania |
| **dymki komunikatów** | skutki czynności bez odpowiednika w kontrakcie oraz odmowy rdzenia | prawa strona ekranu |
| **baza danych** | wiadomości, ustawienia, konta, punkty dostępu, tożsamość | `<katalog danych>\danaco-console.db` |

Katalog danych to `%LOCALAPPDATA%\DanacoConsole`, chyba że wskazano inny zmienną katalogu danych albo przełącznikiem `--dane`.

### 15.3 Prowenancja jako narzędzie diagnostyczne

Prowenancja odpowiada na większość pytań diagnostycznych bez sięgania do dziennika. Nadawana jest **przed uruchomieniem procesu**, więc dociera także wtedy, gdy proces w ogóle nie wystartuje.

| Pytanie | Pole prowenancji |
|---|---|
| Czy w ogóle uruchomiono właściwy program? | plik wykonywalny kanału, pełny wiersz wywołania |
| Czy prompt tożsamości dojechał? | prompt systemowy, tryb nakładki, nazwy warstw, skrót nakładki |
| Czy model dostał wskazane katalogi? | katalogi udostępnione |
| Gdzie model pracuje? | katalog roboczy |
| Czy mosty MCP zostały przekazane? | wskazane konfiguracje MCP |
| Którego konta użyto? | kod konta, katalog konfiguracji konta |
| Czy nastąpiła rotacja konta? | numer próby (druga i dalsze), powód próby |
| Czy ustawienia okna dojechały? | pola modelu, modelu zapasowego, nakładu — **puste oznacza, że parametr nie dotarł** |
| Czy rozmowa jest wznawiana? | pole wznowienia — **w tej wersji zawsze puste** |

Zasada nadrzędna: prowenancja jest przejrzystością, nie bramą — niczego nie dopuszcza i niczego nie wstrzymuje.

---

## 16. Zakończenie i wznowienie pracy

### 16.1 Zamknięcie okna i zakończenie powłoki

**Zamknięcie okna nie kończy pracy.** Naciśnięcie przycisku zamknięcia okna głównego chowa je do zasobnika: rdzeń pracuje dalej, procesy sesji biegną, a ponowne otwarcie wraca do tej samej pracy. Okno jest widokiem, nie właścicielem pracy.

Nie jest to blokada — przycisk działa natychmiast i bez pytania, zmienia się wyłącznie skutek.

Okno przywraca się z zasobnika na dwa sposoby: podwójnym kliknięciem ikony albo pozycją „Pokaż okno” w menu.

**Zakończenie powłoki** („Zakończ powłokę” w menu zasobnika) kończy aplikację okienkową i **pozostawia rdzeń oraz procesy sesji przy pracy**. Zatrzymanie rdzenia jest osobnym, jawnym poleceniem.

### 16.2 Zatrzymanie rdzenia

Pozycja „Zatrzymaj rdzeń” w menu zasobnika zatrzymuje proces rdzenia i melduje wynik komunikatem. Rdzeń pracuje we własnej grupie procesów i bez okna konsoli, więc nie zostanie przerwany przypadkiem — przerwanie powłoki go nie dotyka.

Restart rdzenia jest konieczny w jednym przypadku eksploatacyjnym: **po każdej zmianie katalogu kont**, aby pula rotacji zobaczyła nowy stan.

### 16.3 Wznowienie pracy

| Po czym wznawiamy | Co wraca | Co nie wraca |
|---|---|---|
| **zminimalizowanie okna do zasobnika** | wszystko: karty, okna, historia widoczna, stan sterowania | — |
| **wyjście do Centrum dowodzenia i powrót** | wszystko, jak wyżej | — |
| **odświeżenie strony interfejsu** | widok trasy zapisany w adresie, motyw, konfiguracja | **historia rozmowy**; sesja i okna zakładane są od nowa |
| **zakończenie powłoki i ponowne uruchomienie** (rdzeń pracował) | powłoka dołącza do zastanego rdzenia bez stawiania drugiego procesu | historia rozmowy w widoku |
| **restart rdzenia** | sesje i okna odtwarzane z bazy pod tymi samymi identyfikatorami; ustawienia, konta, dostępy, tożsamość, kanały | historia w widoku; przypisanie sesji do środowiska; stan wyczerpania kont; więź koordynator–wykonawca |

**Zasada praktyczna:** wznowienie pracy w bieżącej wersji oznacza wznowienie **konfiguracji**, nie wznowienie **rozmowy**. Konfiguracja jest trwała i wraca w komplecie; treść rozmowy jest trwała w bazie, lecz interfejs jej nie odtwarza, a model i tak jej nie widzi.

### 16.4 Powrót do sesji trwającej na rdzeniu

Rozłączenie klienta **nie kończy sesji ani jej procesów**. Sesja, którą Operator zostawił — zamykając przeglądarkę, tracąc łączność albo odświeżając stronę — trwa na rdzeniu wraz ze swoimi oknami. Drogą powrotu do niej jest strefa **„Sesje w tle”** Centrum dowodzenia (rozdział [3.2.3](#323-strefa-trzecia--sesje-w-tle)).

**Przebieg powrotu:**

| Krok | Czynność | Miejsce |
|---|---|---|
| 1 | Wejdź na Centrum dowodzenia | przełącznik tras na pasku górnym |
| 2 | Odszukaj sesję w strefie „Sesje w tle” — po tytule, środowisku, liczbie okien albo chwili ostatniej czynności | trzecia strefa strony głównej |
| 3 | Naciśnij **Wróć do sesji** | wiersz sesji |
| 4 | Sprawdź, że widok przeszedł do środowiska sesji | pasek górny, kontekst „środowisko · moduł” |

**Co powrót przywraca:** powiązanie połączenia z sesją, odtworzenie jej okien w rdzeniu wraz ze skierowaniem ich strumieni na bieżące połączenie oraz przeniesienie ogniska na okno, które w tej sesji było ogniskowane.

**Czego powrót nie przywraca:** treści rozmowy prowadzonej w tej sesji wcześniej. Historia widoczna po powrocie zaczyna się od chwili powiązania. Powrót odzyskuje **dostęp do pracy**, nie jej **zapis**.

**Kiedy przycisk powrotu się nie pojawi:** gdy połączenie nie ma tożsamości klienta nadanej w powitaniu. Wykaz pozostaje wtedy czysto informacyjny — pokazuje, że sesja trwa, lecz nie oferuje drogi do niej. Jest to zachowanie zamierzone: przycisk bez skutku byłby atrapą.

**Gdy odmowa przyjdzie w połowie:** odmowa przeniesienia ogniska **nie cofa powiązania**. Sesja pozostaje wtedy powiązana i pracuje, a nieustawione ognisko oznacza wyłącznie, że karta czynna nie została przełączona. Powtórne naciśnięcie przycisku jest bezpieczne.

---

## 17. Pozostałe funkcje wymagające objaśnienia

### 17.1 Obsługa klawiaturą

| Kontekst | Klawisz | Skutek |
|---|---|---|
| pole wysyłki | `Enter` | wysyła wypowiedź |
| pole wysyłki | `Shift` + `Enter` | dokłada wiersz |
| pas kart sesji | strzałki | przenoszą wybór między kartami |
| pas kart sesji | `Home` / `End` | skok na krańce pasa |
| pas kart sesji | `Enter` / spacja | wybiera kartę pod ogniskiem |
| pas kart sesji | `Delete` | zamyka kartę pod ogniskiem |
| szuflada sterowania | `Escape` | zwija szufladę i wraca ogniskiem na uchwyt |
| okna ustawień (nakładkowe) | `Escape` | zamyka okno |
| przełącznik liczby okien | `Tab`, strzałki | ognisko wędruje z wyborem; grupa ma jedno wejście z `Tab` |

Skrót `Ctrl K` widoczny przy polu poleceń **nie jest obsłużony** — jest oznaczeniem zapowiadanym, nie działającym.

### 17.2 Odporność na nazwy nieznane

Produkt stosuje regułę „nieznana nazwa nie zrywa kanału” po obu stronach łącza:

- **rdzeń** odpowiada na nierozpoznaną komendę odpowiedzią typu `*.unknown` z ładunkiem opisującym nieznaną komendę, zamiast zamykać połączenie;
- **klient** przyjmuje nierozpoznane zdarzenie i odnotowuje je w dzienniku nieznanych, zamiast przerywać pracę;
- **wartość spoza wyliczenia** nie gaśnie w interfejsie — pokazywana jest jako własny kod, żeby Operator zobaczył, co przysłał rdzeń;
- **panika obsługiwacza** po stronie rdzenia kończy wyłącznie jedno wywołanie i wraca kodem błędu wewnętrznego.

Ta sama reguła obowiązuje przy uruchamianiu: brak binarki rdzenia, brak pakietu interfejsu, brak katalogu profili, nieczytelna wartość portu — żadne z nich nie wstrzymuje startu.

### 17.3 Łączność i kolejka wychodząca

Transport klienta ponawia połączenie bez ograniczenia liczby prób, z odstępem rosnącym wykładniczo do pułapu piętnastu sekund, z rozproszeniem losowym do 25 procent (co zapobiega zbieganiu się prób wielu okien naraz). Kolejka wychodząca nie jest ograniczona rozmiarem.

**Ograniczenie:** po odzyskaniu połączenia klient **nie powtarza uzgodnienia ani powiązania sesji**. Nowe połączenie jest po stronie rdzenia anonimowe, a uzgodnienie startuje wyłącznie wtedy, gdy okno jeszcze nie powstało.

**Brak kontroli żywotności połączenia.** Ani rdzeń, ani klient nie prowadzą wymiany kontrolnej potwierdzającej, że druga strona żyje.

### 17.4 Adresowanie rozgłoszeń

Rdzeń rozgłasza zdarzenia i fragmenty strumienia **do wszystkich otwartych połączeń**, bez zawężenia do konta połączenia. Mechanizm adresowania pojedynczego urządzenia istnieje w warstwie transportu, lecz nie jest wykorzystywany. Praktyczny skutek: dwa równocześnie otwarte interfejsy zobaczą te same zdarzenia.

### 17.5 Wersjonowanie kontraktu

Przy uzgodnieniu rdzeń odsyła wersję protokołu oraz wykaz obsługiwanych komend. **Klient nie odczytuje żadnej z tych wartości** — nie ma więc dziś sprawdzenia zgodności wersji między interfejsem a rdzeniem.

---

## 18. Zbiorcze zestawienie stanu funkcji

Zestawienie służy szybkiej orientacji przed rozpoczęciem pracy. Kolumna „Stan” używa oznaczeń z rozdziału 1.5.

| Obszar | Funkcja | Stan |
|---|---|---|
| Uruchamianie | powłoka Tauri stawia rdzeń w tle i otwiera okno | DZIAŁA CZĘŚCIOWO (drzewo deweloperskie) |
| Uruchamianie | rdzeń samodzielnie, z migracjami bazy | DZIAŁA |
| Uruchamianie | skrypt wydania składający rdzeń, klienta i powłokę w `C:\DanacoConsole_App` | DZIAŁA (`wydanie.sh` + `pakowanie.sh`; wariant bez powłoki przez `DANACO_BEZ_POWLOKI=1`) |
| Uruchamianie | budowanie pakietu klienta wywoływane przez powłokę | DZIAŁA (`beforeBuildCommand` i `beforeDevCommand` w konfiguracji powłoki) |
| Uruchamianie | instalator NSIS | BRAK — cel `nsis` jest zadeklarowany w konfiguracji powłoki, lecz wydanie wywołuje `cargo tauri build --no-bundle` i instalator nie powstaje; zgodnie z rozdz. 1.3 pkt 1, [README.md](README.md) rozdz. 17.3 i [LICENSE.md](LICENSE.md) rozdz. 10.3 |
| Uruchamianie | instalacja, aktualizacja przyrostowa, deinstalacja, podpis plików wykonywalnych | BRAK |
| Uruchamianie | pakiet interfejsu osadzony w powłoce | DZIAŁA CZĘŚCIOWO — nazwa hosta zasobu osadzonego rozspaja gniazdo (luka nazwy hosta zasobu osadzonego, rozdz. 2.4) |
| Interfejs | trzy trasy, powłoka środowiska, cztery pasy | DZIAŁA |
| Interfejs | motyw jasny i ciemny, żetony v2.0, 82 ikony (stos i wagi dymków ze słownika v2.0 — `komponenty/powiadomienie.css`) | DZIAŁA |
| Interfejs | pole poleceń, dzwonek, AOD, awatar na pasku | ATRAPA |
| Interfejs | dymki komunikatów wraz ze stosem | DZIAŁA (stos spozycjonowany arkuszem aliasów zgodności) |
| Interfejs | nazwy klas ze słownika sprzed design v2.0 | DZIAŁA CZĘŚCIOWO — pokryte arkuszem `aliasy-zgodnosci.css` poza jednym wyjątkiem (`.dn-karta-sesji`, `powloka/karty-sesji.ts:56`, bez definicji w żadnym arkuszu); alias przejściowy, do wycofania po przepięciu widoków |
| Interfejs | strefa „Sesje w tle” na stronie głównej wraz z powrotem do sesji | DZIAŁA |
| Interfejs | objaśnienia kontekstowe [?] przy ustawieniach okna | DZIAŁA (osiem wierszy podsumowania panelu sterowania) |
| Nawigacja | wejście do środowiska i wykazy pozycji | DZIAŁA |
| Nawigacja | wybór modułu otwiera okno operacyjne | BRAK (64 z 79 wierszy katalogu okien w bazie nie ma widoku — rozdz. 7) |
| Nawigacja | komendy nawigacji platformy | NIEZINTEGROWANE (rdzeń gotowy, klient nie wywołuje) |
| Sesje | scena 1/2/3 okien, role, więź | DZIAŁA CZĘŚCIOWO |
| Sesje | przekazanie zlecenia koordynator → wykonawca | NIEZINTEGROWANE (animacja lokalna) |
| Rozmowa | wysyłka, strumień, przerwanie, prowenancja | DZIAŁA |
| Rozmowa | ciągłość kontekstu między turami | BRAK |
| Rozmowa | odtworzenie historii po odświeżeniu | BRAK (odczyt wiadomości okna niewywoływany ani razu w kliencie) |
| Sesje | odczyt wykazu sesji z żywym odpisem obecności | DZIAŁA (strona główna i pulpit) |
| Sesje | powrót do sesji trwającej na rdzeniu (powiązanie + ognisko) | DZIAŁA |
| Sesje | zamknięcie sesji i okna z interfejsu | BRAK (komendy niewywoływane) |
| Modele | rejestr kanałów sterowany danymi | DZIAŁA (bez UI zakładania) |
| Modele | kanał główny rodzaju `cli` | DZIAŁA |
| Modele | kanał testowy — rodzaj `lokalny`, adapter `echo` | DZIAŁA |
| Modele | kanał rodzaju `api` (adapter `api`) | NIEZINTEGROWANE z kontami |
| Modele | kanał rodzaju `sdk` | BRAK (brak adaptera o tym kluczu) |
| Modele | kanał SSH / wykonanie zdalne | BRAK |
| Konta | pełne zarządzanie kontami z interfejsu | DZIAŁA |
| Konta | rotacja po wyczerpaniu limitu | DZIAŁA CZĘŚCIOWO (wymaga restartu po zmianach, stan nieutrwalany, rotacja niewidoczna) |
| Tożsamość | nakładka konstytucja → profil → ekspertyza | DZIAŁA |
| Tożsamość | podgląd złożonego promptu | DZIAŁA |
| Konfiguracja | katalog kategorii i pól z rdzenia, zapis poziom + oś | DZIAŁA |
| Konfiguracja | wpływ zapisanych kluczy na wykonanie | DZIAŁA CZĘŚCIOWO (patrz 6.4) |
| Konfiguracja | punkty izolacji | NIEZINTEGROWANE (brak egzekutorów, brak UI) |
| Dostępy | punkty rodzaju most MCP, nadania per okno | DZIAŁA |
| Dostępy | punkty rodzaju katalog lokalny | BRAK (niezapisywalne) |
| Pliki | katalogi robocze okna | DZIAŁA |
| Pliki | treści obszerne jako pliki, artefakty | BRAK |
| Dane | trwałość rozmowy, ustawień, kont, dostępów, tożsamości | DZIAŁA |
| Dane | zapis zmian okna do bazy | NIEZINTEGROWANE |
| Dane | przypisanie sesji do środowiska | DZIAŁA CZĘŚCIOWO — sesja trafia zawsze do pierwszego środowiska, niezależnie od wyboru Operatora |
| Orkiestracja | pętla koordynator–wykonawca w rdzeniu | DZIAŁA |
| Orkiestracja | telemetria postępu w interfejsie | DZIAŁA CZĘŚCIOWO — odbiorca istnieje na pulpicie Mission Control (wykaz procesów w tle); w oknie komunikacji odbiorcy brak |
| Orkiestracja | silnik kolejek | ATRAPA (pozycje kolejki nigdy nie powstają) |
| Orkiestracja | Mission Control — dane z rdzenia | DZIAŁA CZĘŚCIOWO — odczyty wykazu sesji, okien i kanałów oraz cztery subskrypcje zdarzeń; uczciwe stany pusty / oczekiwanie / błąd |
| Orkiestracja | Mission Control — wykaz kolejek i procesów | DZIAŁA CZĘŚCIOWO — kontrakt nie ma odczytu wykazu kolejek ani bieżących procesów; wykazy budują się wyłącznie ze zdarzeń |
| Narzędzia | narzędzia platformy dla modelu (39 deklaracji) | NIEZINTEGROWANE (zero konsumentów) |
| Narzędzia | panel czynności | NIEZINTEGROWANE |
| Pamięć | pamięć wielopoziomowa | NIEZINTEGROWANE (repozytorium bez komend) |
| Rozszerzenia | rozszerzenia i marketplace | BRAK |

---

## 19. Kwestie pozostawione do decyzji Operatora

Poniższe zagadnienia wykraczają poza zakres dokumentacji i wymagają rozstrzygnięcia produktowego albo architektonicznego. Instrukcja odnotowuje je, nie rozstrzyga.

1. **Sposób pakowania produktu** — czy rdzeń ma być zasobem osadzonym w powłoce, czy — jak dziś — odrębnym plikiem wykonywalnym obok niej; oraz czy pakiet interfejsu ma być serwowany przez rdzeń, czy osadzany w powłoce. Rozstrzygnięcie ma skutek praktyczny, ponieważ droga pakietu osadzonego nie zestawia połączenia (luka nazwy hosta zasobu osadzonego). Osobno: czy instalator NSIS ma zostać włączony do skryptu wydania jako krok warunkowy, czy pozostać czynnością odrębną. (rozdziały 2.4 i 2.5)
2. **Trwałe rozwiązanie stanu wyjściowego instalacji** — zaczyn migracji z wierszem kanału głównego, ekran zakładania kanałów w interfejsie, albo wpięcie istniejącego modułu zapewnienia kanału do głównego wejścia aplikacji. Wariant trzeci jest już napisany i wymaga wyłącznie podpięcia oraz rozstrzygnięcia, w którym momencie uzgodnienia ma działać. (rozdział 2.7)
3. **Miejsce trwałego zapamiętania identyfikatora rozmowy kanału** — warunek domknięcia ciągłości rozmowy między turami. (rozdział 5.8)
4. **Miejsce, w którym mieszka słownik kluczy ustawień okna** — kontrakt czy katalog ustawień w bazie; warunek uzgodnienia trzech kluczy zapisywanych dziś pod nazwami nieznanymi rdzeniowi. (rozdział 8.3)
5. **Droga otwarcia okien modułu Automations oraz sposób osadzania okien operacyjnych w obszarze roboczym** — rozszerzenie kontraktu o komendy modułowe albo jedna komenda czynności okna. (rozdział 9)
6. **Miejsce składowania treści obszernych w katalogu danych oraz kształt tabel załącznika wiadomości i artefaktu** — warunek domknięcia plikowej części decyzji o trwałości. (rozdział 11.3)
7. **Zachowanie odczytu konfiguracji bez wskazania poziomu** — czy ma oddawać politykę efektywną, czy komplet zapisów; od tego zależy poprawność podglądu dziedziczenia w oknie konfiguracji. (rozdział 13.1)
8. **Zakres i sposób odświeżania puli rotacji kont bez restartu rdzenia** oraz sposób utrwalania stanu wyczerpania kont. (rozdział 6.3)
9. **Jeden identyfikator kanału w kontrakcie** — numer wiersza czy kod wiersza; dziś odczyt wykazu i komendy zmiany używają różnych. (rozdział 6.2)
10. **Sposób doprowadzenia kont rodzaju `api` i `sdk` do drogi wykonawczej** — dziś do wykonania dociera wyłącznie rodzaj `cli`. (rozdział 6.3)
11. **Podział zakresów między [INSTRUKCJA-UZYTKOWANIA.md](INSTRUKCJA-UZYTKOWANIA.md) a [INSTALACJA-I-KONFIGURACJA.md](INSTALACJA-I-KONFIGURACJA.md)** — czy rozdział 2 niniejszego dokumentu ma zostać skrócony do odesłania, czy deklaracja podziału w [README.md](README.md) (rozdz. 17 i 19) ma zostać poprawiona tak, by uwzględnić skrót uruchomieniowy w instrukcji użytkowania. Do czasu rozstrzygnięcia w sprawach instalacyjnych rozstrzyga dokument instalacyjny. (rozdział 1.6)
12. **Termin i sposób wycofania arkusza aliasów zgodności** — przepięcie widoków na słownik design v2.0 i skreślenie arkusza, albo utrzymywanie aliasów jako warstwy stałej. Do czasu rozstrzygnięcia obowiązuje reguła: zmiana pierwowzoru w bibliotece wymaga naniesienia jej także w arkuszu aliasów. (rozdział 3.5)
13. **Rodzaj kanału `sdk`** — czy ma otrzymać własny adapter, czy zostać usunięty ze zbioru dopuszczalnych wartości kolumny rodzaju i z wykazu kontraktu. Dziś jest wartością dopuszczalną, dla której nie ma kodu budującego kanał. (rozdział 6.1)
14. **Uzupełnienie kontraktu o odczyty wykazu kolejek i bieżących procesów** — bez nich pulpit Mission Control buduje te dwa wykazy wyłącznie ze zdarzeń, więc przed pierwszym zdarzeniem pozostają puste niezależnie od stanu rdzenia. (rozdział 3.4.3)
15. **Odbiorca telemetrii postępu w oknie komunikacji** — czy licznik obiegów pętli koordynator–wykonawca ma być pokazywany na scenie okien, czy pozostać wyłącznie na pulpicie. (rozdział 5.5)

---

*Koniec dokumentu. Instrukcja użytkowania — Stan wdrożenia, wersja 2.0, 2026-08-20.*

---
*Danaco Console — Platforma AI Workspace OS · v2.0 · status Deweloperski*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz.*
*Warunki korzystania: [Licencja produktu](LICENSE.md). Kontakt: support@danaco-group.pl*
