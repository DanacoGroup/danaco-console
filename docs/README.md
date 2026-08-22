# Danaco Console — Platforma AI Workspace OS

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
| **Tytuł** | Platforma AI Workspace OS — główna karta techniczno-produktowa |
| **Klasa dokumentu** | Stan wdrożenia |
| **Odbiorcy** | Operator · deweloper · projektant |
| **Przeznaczenie** | Opisuje produkt rzeczywisty — taki, jaki został zbudowany — w podziale na warstwy architektury, moduły, środowiska, kanały wykonawcze, konfigurację, dane i warstwę wizualną, z jawnym oznaczeniem stanu wdrożenia każdej funkcji. Na jego podstawie ustala się, co w bieżącym wydaniu wolno przedstawiać jako działające. |
| **Zakres** | wielośrodowiskowy system operacyjny dla sztucznej inteligencji, umożliwiający budowę i nadzorowanie cyfrowej organizacji złożonej z modeli, agentów, procesów i automatyzacji realizujących rzeczywistą pracę: architektura rdzenia, powłoki i klienta, komponenty, moduły, środowiska, modele i kanały wykonawcze, konfiguracja, sesje i okna komunikacji, dane, rozszerzenia, integracje, struktura katalogów, warstwa wizualna, scenariusze zastosowania, wymagania, uruchomienie, znane ograniczenia |
| **Poza zakresem** | procedura instalacji i pełny wykaz kluczy konfiguracji — [Instalacja i konfiguracja](INSTALACJA-I-KONFIGURACJA.md); obsługa produktu przez Operatora krok po kroku — [Instrukcja użytkowania](INSTRUKCJA-UZYTKOWANIA.md); warunki korzystania — [Licencja produktu](LICENSE.md) |
| **Dokument nadrzędny** | [Spis opracowań](SPIS-OPRACOWAN.md) |
| **Dokumenty powiązane** | [Instrukcja użytkowania](INSTRUKCJA-UZYTKOWANIA.md) · [Instalacja i konfiguracja](INSTALACJA-I-KONFIGURACJA.md) · [Licencja produktu](LICENSE.md) · [Standard redakcyjny i językowy](STANDARD-REDAKCYJNY-I-JEZYKOWY.md) |
| **Prototypy odniesienia** | nie dotyczy — dokument klasy Stan wdrożenia opisuje kod zbudowany, nie projekt okna; prototypy przywołują opracowania klasy Specyfikacja docelowa |
| **Źródła normatywne** | `budowa/shared/contract.json` · `budowa/core/` · `budowa/client/src/` · `budowa/rdzen/dane/migracje/` · `design/zasoby/zetony/zetony.css` |
| **Zasada nadrzędna** | Żadna funkcja o stanie ATRAPA, NIEZINTEGROWANE albo BRAK nie jest przedstawiana jako działająca. |

**Nawigacja zbioru:** [Spis opracowań](SPIS-OPRACOWAN.md) ·
[Instrukcja użytkowania](INSTRUKCJA-UZYTKOWANIA.md) ·
[Instalacja i konfiguracja](INSTALACJA-I-KONFIGURACJA.md) ·
[Licencja produktu](LICENSE.md)

Niniejszy dokument jest **główną kartą techniczno-produktową** Danaco Console.
Opisuje produkt rzeczywisty — taki, jaki został zbudowany — a nie zamierzony.
Wszystkie ścieżki, nazwy komend, przełączniki i klucze konfiguracji zostały
zweryfikowane w kodzie źródłowym.

**Podstawa opisu i jej dwie warstwy.** Zakres funkcjonalny i statystyki ustaleń
pochodzą z audytu technicznego przeprowadzonego w stanie sprzed prac naprawczych. Po tej
rewizji wykonano i **zatwierdzono** prace naprawcze w trzech rewizjach —
trzy kolejne etapy prac naprawczych (HEAD, stan na dzień odcięcia) — które
unieważniły sześć ustaleń audytu: konfiguracja powłoki (`tauri.conf.json`)
wiąże budowanie interfejsu z budowaniem pakietu; repozytorium zawiera skrypty
wydania i pakowania produktu; pulpit Mission Control czerpie dane z rdzenia
zamiast z pliku przykładowego; komendy `session.focus` i `session.bind` mają
konsumenta w widoku strony głównej; stos dymków powiadomień i pozostałe trzynaście
klas ze słownika sprzed design v2.0 mają dziś pełne pokrycie w arkuszach
(`komponenty/aliasy-zgodnosci.css`); most MCP po SSH otrzymał zapis
decyzyjny o statusie `Zaproponowana`. Rozdziały [2.5](#25-decyzje-architektoniczne-jako-źródło-prawdy),
[3.3](#33-powłoka-natywna--tauri), [3.8](#38-decyzje-architektoniczne),
[6](#6-środowiska), [12.2](#122-mosty-mcp-po-ssh),
[13.1](#131-model-czterech-katalogów), [14](#14-warstwa-wizualna),
[17.3](#173-wydanie-produktu-do-cdanacoconsole_app),
[18](#18-znane-ograniczenia-wersji-v20) i [20](#20-kwestie-do-decyzji-operatora)
opisują stan bieżący w bieżącym stanie repozytorium, nie stan rewizji odniesienia audytu.
Pomiary liczbowe podane w dokumencie są zliczeniami własnymi wykonanymi na
drzewie roboczym HEAD w dniu odcięcia (drzewo czyste, bez zmian niezatwierdzonych);
przy każdym z nich wskazano metodę i miejsce pomiaru.

**Podstawa dokumentu.** Opis powstał na bieżącym stanie repozytorium, przy czystym
drzewie roboczym — po zamknięciu prac naprawczych domykających ustalenia audytu
technicznego. Prace objęły ścieżkę wydania i bramy mierzące wywołania, domknięcie
warstwy wizualnej v2.0, zasilenie pulpitu i strony głównej danymi z rdzenia,
sprawdziany widoków oraz porządek dokumentacyjny; łącznie względem stanu sprzed
prac naprawczych: 53 pliki zmienione, 28 pozycji nowych, 1 usunięta.

**Data odcięcia opisu.** Dniem pomiaru drzewa roboczego jest 2026-08-12. Pole
**Data** metryki produktowej niesie datę wydania dokumentu, nie datę pomiaru;
obie daty są odrębnymi informacjami i nie podlegają uzgodnieniu.

Dokument leży w katalogu `C:\DanacoConsole_App` obok plików wykonywalnych
złożonych z tego samego drzewa i opisuje ten sam stan kodu.

---

## Oznaczenia stanu

Produkt jest w wersji **v2.0 o statusie Deweloperskim**. Część warstw jest
domknięta i sprawdzona, część istnieje jako mechanizm bez wpięcia w przepływ,
część pozostaje wyłącznie pozycją nawigacji. Dokument rozdziela te stany
jednoznacznie, przy każdej sekcji funkcjonalnej:

| Oznaczenie | Znaczenie |
|---|---|
| **[DZIAŁA]** | Funkcja domknięta od interfejsu (albo od komendy kontraktu) do skutku; potwierdzona audytem. |
| **[DZIAŁA CZĘŚCIOWO]** | Mechanizm działa, ale w ograniczonym zakresie albo z udokumentowanym zastrzeżeniem. |
| **[NIEZINTEGROWANE]** | Kod istnieje i jest poprawny po jednej albo obu stronach, lecz nie ma konsumenta — funkcja nie jest osiągalna dla Operatora. |
| **[ATRAPA]** | Element widoczny w interfejsie (pozycja menu, kafel, przycisk), za którym nie stoi realizacja. |
| **[BRAK]** | Przewidziane koncepcją, w bieżącej wersji niezaimplementowane. |
| **[DO DECYZJI OPERATORA]** | Kwestia nierozstrzygnięta; dokument jej nie przesądza. |

Zasada nadrzędna dokumentu: **żadna funkcja o stanie ATRAPA, NIEZINTEGROWANE
albo BRAK nie jest przedstawiana jako działająca.** Tam, gdzie mechanizm istnieje,
lecz nie jest dostępny z interfejsu, jest to napisane wprost.

---

## Spis treści

1. [Oznaczenia stanu](#oznaczenia-stanu)
2. [Czym jest Danaco Console](#1-czym-jest-danaco-console)
   - [1.1 Definicja produktu](#11-definicja-produktu)
   - [1.2 Czym Danaco Console nie jest](#12-czym-danaco-console-nie-jest)
   - [1.3 Stan wersji v2.0 — obraz uczciwy](#13-stan-wersji-v20--obraz-uczciwy)
3. [Filozofia i założenia produktu](#2-filozofia-i-założenia-produktu)
   - [2.1 Współpraca człowieka i sztucznej inteligencji](#21-współpraca-człowieka-i-sztucznej-inteligencji)
   - [2.2 Cyfrowa organizacja](#22-cyfrowa-organizacja)
   - [2.3 Operator](#23-operator)
   - [2.4 Zasady konstrukcyjne](#24-zasady-konstrukcyjne)
   - [2.5 Decyzje architektoniczne jako źródło prawdy](#25-decyzje-architektoniczne-jako-źródło-prawdy)
4. [Architektura](#3-architektura)
   - [3.1 Obraz ogólny](#31-obraz-ogólny)
   - [3.2 Rdzeń — serwer Go](#32-rdzeń--serwer-go)
   - [3.3 Powłoka natywna — Tauri](#33-powłoka-natywna--tauri)
   - [3.4 Klient — interfejs TypeScript](#34-klient--interfejs-typescript)
   - [3.5 Kontrakt `shared/`](#35-kontrakt-shared)
   - [3.6 Transport: WebSocket i JSON](#36-transport-websocket-i-json)
   - [3.7 Trwałość: SQLite](#37-trwałość-sqlite)
   - [3.8 Decyzje architektoniczne](#38-decyzje-architektoniczne)
   - [3.9 Reguły modularności i brama weryfikacyjna](#39-reguły-modularności-i-brama-weryfikacyjna)
5. [Główne komponenty](#4-główne-komponenty)
   - [4.1 Pakiety rdzenia](#41-pakiety-rdzenia)
   - [4.2 Obszary klienta](#42-obszary-klienta)
   - [4.3 Moduły powłoki](#43-moduły-powłoki)
6. [Moduły](#5-moduły)
7. [Środowiska](#6-środowiska)
8. [Modele i kanały wykonawcze](#7-modele-i-kanały-wykonawcze)
   - [7.1 Rejestr kanałów sterowany danymi](#71-rejestr-kanałów-sterowany-danymi)
   - [7.2 Kanał główny — Claude Code CLI](#72-kanał-główny--claude-code-cli)
   - [7.3 Kanał `api`](#73-kanał-api)
   - [7.4 Kanał `echo`](#74-kanał-echo)
   - [7.5 Konta i pula rotacji](#75-konta-i-pula-rotacji)
   - [7.6 Nakładka tożsamości](#76-nakładka-tożsamości)
   - [7.7 Prowenancja wywołania](#77-prowenancja-wywołania)
   - [7.8 Czego warstwa kanałów w tej wersji nie robi](#78-czego-warstwa-kanałów-w-tej-wersji-nie-robi)
9. [Konfiguracja](#8-konfiguracja)
   - [8.1 Katalog ustawień sterowany danymi](#81-katalog-ustawień-sterowany-danymi)
   - [8.2 Osiem poziomów zasięgu](#82-osiem-poziomów-zasięgu)
   - [8.3 Osie rozstrzygania](#83-osie-rozstrzygania)
   - [8.4 Które klucze realnie sterują wykonaniem](#84-które-klucze-realnie-sterują-wykonaniem)
   - [8.5 Izolacja zasięgu](#85-izolacja-zasięgu)
   - [8.6 Konfiguracja startowa procesu rdzenia](#86-konfiguracja-startowa-procesu-rdzenia)
10. [Sesje i okna komunikacji](#9-sesje-i-okna-komunikacji)
   - [9.1 Model bytów](#91-model-bytów)
   - [9.2 Cykl życia](#92-cykl-życia)
   - [9.3 Okna równoległe i role](#93-okna-równoległe-i-role)
   - [9.4 Pętla koordynator–wykonawca](#94-pętla-koordynatorwykonawca)
   - [9.5 Kolejki](#95-kolejki)
   - [9.6 Anulowanie i odporność](#96-anulowanie-i-odporność)
11. [Dane](#10-dane)
   - [10.1 Schemat i migracje](#101-schemat-i-migracje)
   - [10.2 Trwałość rozmowy](#102-trwałość-rozmowy)
   - [10.3 Odtworzenie stanu](#103-odtworzenie-stanu)
   - [10.4 Pamięć wielopoziomowa](#104-pamięć-wielopoziomowa)
12. [Rozszerzenia](#11-rozszerzenia)
13. [Integracje](#12-integracje)
   - [12.1 Punkty dostępu i nadania](#121-punkty-dostępu-i-nadania)
   - [12.2 Mosty MCP po SSH](#122-mosty-mcp-po-ssh)
   - [12.3 Claude Code CLI](#123-claude-code-cli)
14. [Struktura aplikacji i katalogów](#13-struktura-aplikacji-i-katalogów)
   - [13.1 Model czterech katalogów](#131-model-czterech-katalogów)
   - [13.2 Drzewo repozytorium budowy](#132-drzewo-repozytorium-budowy)
   - [13.3 Katalog danych czasu pracy](#133-katalog-danych-czasu-pracy)
15. [Warstwa wizualna](#14-warstwa-wizualna)
16. [Scenariusze zastosowania](#15-scenariusze-zastosowania)
   - [15.1 Scenariusze wykonalne w v2.0](#151-scenariusze-wykonalne-w-v20)
   - [15.2 Scenariusze przewidziane koncepcją, w v2.0 niewykonalne](#152-scenariusze-przewidziane-koncepcją-w-v20-niewykonalne)
17. [Wymagania](#16-wymagania)
   - [16.1 Środowisko uruchomieniowe](#161-środowisko-uruchomieniowe)
   - [16.2 Środowisko budowy](#162-środowisko-budowy)
18. [Podstawowe uruchomienie](#17-podstawowe-uruchomienie)
   - [17.1 Zarys uruchomienia w drzewie deweloperskim](#171-zarys-uruchomienia-w-drzewie-deweloperskim)
   - [17.2 Doprowadzenie świeżej instalacji do pracy — czynność obowiązkowa](#172-doprowadzenie-świeżej-instalacji-do-pracy--czynność-obowiązkowa)
   - [17.3 Wydanie produktu do `C:\DanacoConsole_App`](#173-wydanie-produktu-do-cdanacoconsole_app)
19. [Znane ograniczenia wersji v2.0](#18-znane-ograniczenia-wersji-v20)
20. [Pozostała dokumentacja](#19-pozostała-dokumentacja)
21. [Kwestie do decyzji Operatora](#20-kwestie-do-decyzji-operatora)

---

## 1. Czym jest Danaco Console

### 1.1 Definicja produktu

**Danaco Console jest operacyjnym systemem współpracy człowieka i sztucznej
inteligencji.** Nazwa kategorii — „Platforma AI Workspace OS” — nie jest ozdobnikiem
marketingowym, lecz opisem konstrukcji: produkt nie dostarcza jednego okna rozmowy
z jednym modelem, tylko środowisko, w którym modele, konta, kanały wykonawcze,
katalogi robocze, tożsamości, uprawnienia, sesje i procesy są **bytami pierwszej
klasy**, konfigurowalnymi i podlegającymi nadzorowi.

Produkt składa się z czterech wyspecjalizowanych środowisk pracy — wiedzy
(**TalkIn**), pracy (**WorkSpace**), technologii (**CodeStudio**) oraz orkiestracji
inteligencji (**MultitaskingAI**) — połączonych wspólną warstwą komunikacji,
konfiguracji, dostępów, pamięci i trwałości.

Celem produktu nie jest generowanie odpowiedzi. Celem jest **umożliwienie
Operatorowi budowy i zarządzania cyfrową organizacją**, która realizuje rzeczywistą
pracę: pracuje na plikach, wykonuje polecenia na maszynach, prowadzi projekty
i utrzymuje ciągłość zadania także wtedy, gdy Operator patrzy w inne miejsce.

Architektura wynika wprost z tego celu. Rdzeń jest procesem serwerowym, który
przeżywa rozłączenie interfejsu; sesja i okno komunikacji są bytami rdzenia, nie
stanami przeglądarki; kanał modelu jest wierszem w rejestrze, a nie gałęzią
w kodzie; wywołanie modelu niesie pełną prowenancję tego, co do modelu poszło.

### 1.2 Czym Danaco Console nie jest

Rozgraniczenie jest istotne dla poprawnego czytania dalszych rozdziałów:

- **Nie jest komunikatorem ani nakładką na czat.** Okno komunikacji jest jednym
  z bytów systemu, a nie jego całością; niesie moduł, kanał modelu, listę katalogów
  roboczych, zasięg wykonania, tryb uprawnień i rolę w pętli koordynator–wykonawca.
- **Nie jest bramką do jednego dostawcy modeli.** Rejestr kanałów jest sterowany
  danymi: nowy kanał to nowy wiersz tabeli, nie nowy typ w kodzie.
- **Nie jest systemem kontroli dostępu.** Zgodnie z przyjętą zasadą jedyną kontrolą dostępu
  jest uwierzytelnianie, a jego egzekwowanie jest konfigurowalne przez Operatora;
  w fazie budowy pozostaje wyłączone.
- **Nie jest usługą chmurową.** Rdzeń pracuje lokalnie na maszynie Operatora;
  docelowe przeniesienie na serwer `danaco-system` jest przewidziane, ale zasięg
  wykonania modelu jest od umiejscowienia rdzenia **niezależny** i ustalany per okno
  komunikacji.

### 1.3 Stan wersji v2.0 — obraz uczciwy

Wersja v2.0 ma status **Deweloperski** i taki właśnie jest jej stan faktyczny.
Fundament techniczny — kontrakt komunikacyjny, rdzeń Go, warstwa danych, warstwa
konfiguracji, warstwa tożsamości, warstwa dostępów oraz higiena repozytorium — jest
zbudowany rzetelnie i zgodnie z decyzjami architektonicznymi. Warstwa funkcjonalna produktu
jest natomiast wykonana w zakresie wąskim: z 79 wierszy katalogu okien operacyjnych
w bazie ([INSTRUKCJA-UZYTKOWANIA.md](INSTRUKCJA-UZYTKOWANIA.md) rozdz. 7) istnieje
**jeden widok** (Chat Window), a nawigacja modułowa czterech środowisk nie
przeładowuje przestrzeni roboczej.

Trzy stwierdzenia, których dokument nie ukrywa, bo determinują odbiór produktu:

1. **Świeża instalacja nie wykona ani jednej tury modelu bez ręcznego
   doprowadzenia.** Rejestr kanałów modelu nie ma zaczynu w migracjach, pula kont
   przy braku wpisu odmawia wywołania, a interfejs w `window.create` podaje **rodzaj**
   kanału (`cli`) w miejscu przeznaczonym na identyfikator wiersza rejestru.
   Szczegóły i konsekwencje: rozdział [7](#7-modele-i-kanały-wykonawcze)
   oraz [17](#17-podstawowe-uruchomienie).
2. **Rozmowa nie ma ciągłości między turami.** Każda tura kanału głównego startuje
   nowy proces bez historii; przełącznik `--resume` nie jest podawany, mimo że
   identyfikator rozmowy programu jest odczytywany ze strumienia. Model nie pamięta
   poprzednich wypowiedzi okna.
3. **Historia rozmowy nie wraca na ekran.** Rdzeń utrwala każdą wiadomość
   i odpowiedź oraz oddaje historię okna komendą `message.list`, ale interfejs nie
   wywołuje jej ani razu. Widoczna historia jest buforem bieżącego połączenia:
   odświeżenie widoku albo ponowne połączenie czyści ekran, choć dane w bazie
   pozostają nienaruszone. Szczegóły: [9.2](#92-cykl-życia).

Ścieżka wydania produktu **istnieje** i jest jednym poleceniem
(`bash budowa/scripts/wydanie.sh`) — zbudowanie i rozmieszczenie produktu
w `C:\DanacoConsole_App` opisuje rozdział [17.3](#173-wydanie-produktu-do-cdanacoconsole_app).
Odłożonym krokiem pozostaje wyłącznie instalator NSIS.

Statystyka audytu zamykającego etap, przeprowadzonego na rewizji odniesienia
stanu sprzed prac naprawczych: **169 ustaleń** w podziale na wagę — krytyczna 16,
wysoka 67, średnia 60, niska 26 — oraz w podziale na stany:
`NIEZINTEGROWANE` 54, `DZIAŁA CZĘŚCIOWO` 58, `BRAK` 22, `NIEZGODNE Z KONCEPCJĄ` 19,
`ATRAPA/STUB` 9, `REGRESJA` 2, `DZIAŁA` 5. Sześć z tych ustaleń zostało po audycie
zamknięte pracami opisanymi w nagłówku dokumentu; przy każdym takim miejscu jest
to odnotowane.

Wnioskiem nie jest, że produkt nie istnieje. Wnioskiem jest, że **istnieje jego
fundament i pierwszy przekrój pionowy**, a warstwa modułów pozostaje przed
wykonaniem.

---

## 2. Filozofia i założenia produktu

### 2.1 Współpraca człowieka i sztucznej inteligencji

Danaco Console przyjmuje, że model nie jest narzędziem odpowiadającym na pytania,
lecz **wykonawcą pracy o określonej tożsamości, uprawnieniach i miejscu pracy**.
Konsekwencje tego założenia są wpisane w architekturę:

- **Model ma tożsamość, nie prompt.** Nakładka tożsamości składa się z trzech warstw
  uporządkowanych według krytyczności: konstytucja → profil (rola) → ekspertyza
  zadaniowa. Tożsamość jest **zamieniana**, nie dołączana do promptu fabrycznego —
  wartością domyślną trybu jest `ZASTAP`.
- **Model ma miejsce pracy.** Okno komunikacji niesie listę katalogów roboczych,
  katalog roboczy sesji oraz zbiór nadań dostępu do maszyn.
- **Model ma zakres zgody.** Tryb uprawnień okna przekłada się bezpośrednio na
  przełącznik uruchomienia kanału głównego.
- **Praca jest przejrzysta, nie bramkowana.** Zamiast blokad — prowenancja: pierwszy
  fragment strumienia odpowiedzi niesie opis tego, co poszło do modelu.

### 2.2 Cyfrowa organizacja

Pojęcie „cyfrowej organizacji” ma w produkcie znaczenie techniczne, a nie
metaforyczne. Wyraża się trzema mechanizmami:

1. **Role okien.** Okno komunikacji ma rolę: `wykonawca`, `koordynator` albo
   `samodzielne`. Okno samodzielne pozostaje poza pętlą.
2. **Pętla koordynator–wykonawca.** Zakończenie tury wykonawcy wybudza koordynatora;
   jeden silnik obsługuje pętlę sesyjną i orkiestrację MultitaskingAI — dwie
   równoległe implementacje są zakazane.
3. **Bieg naprawczy bez limitu.** Krok zakończony błędem powtarza się dowolną liczbę
   razy; przerwanie należy do Operatora. Zamiast bramy — przejrzystość: licznik
   obiegów, narastające zużycie, wpis w dzienniku i zawsze aktywny przycisk
   zatrzymania.

Silnik pętli istnieje w rdzeniu i startuje turą **[DZIAŁA CZĘŚCIOWO]**; warstwa
interfejsu orkiestracji — panel MultitaskingAI, kolejki ról, harmonogram, monitor —
w bieżącej wersji **[BRAK]**. Szczegóły w rozdziałach [6](#6-środowiska)
i [9](#9-sesje-i-okna-komunikacji).

### 2.3 Operator

**Operator** to rola człowieka pracującego z produktem — nie „użytkownik” w sensie
konsumenta interfejsu. Operator:

- zakłada i konfiguruje kanały modelu, konta i punkty dostępu,
- ustala tożsamość modeli na dowolnym z ośmiu poziomów zasięgu i trzech osiach,
- decyduje o zasięgu wykonania, trybie uprawnień i katalogach roboczych okna,
- nadzoruje pracę biegnącą w tle i przerywa ją, gdy uzna to za właściwe,
- podejmuje decyzje architektoniczne i produktowe, których wykonawca (człowiek albo
  agent) podjąć nie może.

Ostatni punkt jest regułą wiążącą także dla niniejszej dokumentacji: kwestie
nierozstrzygnięte są w tekście oznaczone **[DO DECYZJI OPERATORA]** i zebrane
w rozdziale [20](#20-kwestie-do-decyzji-operatora).

### 2.4 Zasady konstrukcyjne

Sześć zasad przenika cały produkt i wyjaśnia większość rozwiązań technicznych
opisanych dalej:

| Zasada | Treść |
|---|---|
| **Fail-open** | Żaden element opcjonalny nie blokuje uruchomienia. Brak ustawienia znaczy „wartość domyślna”, nie „odmowa”. Nieznana komenda kontraktu wraca odpowiedzią `*.unknown`, a nie zerwaniem połączenia. |
| **Zero blokad w interfejsie** | Brak wyłączanych przycisków jako formy blokady; niegotowość sygnalizuje się opisowo. Okna nakładkowe i potwierdzenia są opcją, nie bramą. |
| **Prowenancja zamiast bramki** | To, co poszło do modelu (argumenty, prompt systemowy, plik ustawień, suma kontrolna nakładki), jest raportowane dla przejrzystości — nigdy nie służy dopuszczaniu ani blokowaniu wywołania. |
| **Bez wartości wpisanych na sztywno** | Prompt systemowy, powłoka wykonawcza, hooki i model pochodzą wyłącznie z konfiguracji. Kod nie zna ani jednego zdania promptu. |
| **Konfiguracja zamiast logiki** | Zachowanie wynika z ustawień, nie z gałęzi w kodzie; brak ustawienia = wartość domyślna. |
| **Sterowanie danymi** | Rejestry — kanałów, akcji, środowisk, modułów, ustawień, tożsamości — są tabelami. Nowa pozycja to nowy wiersz migracji, nie nowa gałąź rdzenia. |

Błąd techniczny jest błędem **wyłącznie bieżącego wywołania** — nie blokuje trwale
sesji, konta ani kolejnych prób.

### 2.5 Decyzje architektoniczne jako źródło prawdy

Decyzje architektoniczne produktu są rozstrzygnięciami **obowiązującymi**, a nie
propozycjami do ponownego otwierania przy każdej fali budowy. Każda z nich jest
notowana w dokumentacji i opisana statusem: `Przyjęta` (obowiązuje),
`Zaproponowana` (przed przyjęciem), `Zastąpiona` (wyparta przez nowszą),
`Odrzucona` (rozstrzygnięta odmownie). Część
rozstrzygnięć została przeniesiona z projektu poprzedzającego, pozostałe powstały
w toku budowy Danaco Console; najnowszym z nich jest most MCP po SSH, zapisany
w pracach porządkowych nad dokumentacją.

Tryb jest bezwzględny: nowa decyzja architektoniczna trafia **najpierw do
dokumentacji, dopiero potem do kodu**. Dokumentacja odnotowuje jawnie, że w kilku
przypadkach kod wyprzedził zapis decyzji, i traktuje to jako nieprawidłowość
odnotowaną, nie zaakceptowaną.

Znane luki produktu są odnotowywane odrębnie, poza treścią niniejszego dokumentu.

---

## 3. Architektura

### 3.1 Obraz ogólny

Produkt składa się z czterech warstw o rozłącznych odpowiedzialnościach:

```
┌──────────────────────────────────────────────────────────────┐
│  POWŁOKA NATYWNA (Tauri / Rust)                              │
│  okno systemowe · ikona zasobnika · rdzeń w tle              │
│  natywny wybór katalogu roboczego                            │
└───────────────┬──────────────────────────────────────────────┘
                │ webview platformy
┌───────────────▼──────────────────────────────────────────────┐
│  KLIENT (TypeScript + Vite)                                  │
│  strona główna · powłoka środowiska · scena okien            │
│  okno komunikacji · konfiguracja · dostępy · modele          │
└───────────────┬──────────────────────────────────────────────┘
                │ WebSocket, JSON, koperta {type, sessionId, payload}
┌───────────────▼──────────────────────────────────────────────┐
│  RDZEŃ (Go)                                                  │
│  transport · protokół · rejestr komend · sesje i okna        │
│  konfiguracja warstwowa · rejestr kanałów · injection        │
└───────────────┬──────────────────────────────────────────────┘
                │ database/sql
┌───────────────▼──────────────────────────────────────────────┐
│  TRWAŁOŚĆ (SQLite, jeden plik, WAL)                          │
└──────────────────────────────────────────────────────────────┘

        ↑ wszystkie warstwy nazywają byty tak samo:
          KONTRAKT shared/contract.json → contract.go + contract.ts
```

Poza tym łańcuchem rdzeń uruchamia **proces modelu** (kanał wykonawczy) i przekazuje
mu konfigurację mostów MCP — opis w rozdziałach [7](#7-modele-i-kanały-wykonawcze)
i [12](#12-integracje).

### 3.2 Rdzeń — serwer Go

**Decyzja.** Rdzeń jest jednym programem w języku Go (moduł `danacoconsole`,
Go 1.26), którego punktem wejścia jest `server/cmd/danaco-console/main.go` — plik
72-wierszowy, wyłącznie kompozycyjny, bez logiki i bez definicji typów.

Rdzeń pracuje w jednej z trzech ról, wybieranej przełącznikiem `--role` albo zmienną
`DANACO_ROLA`:

| Rola | Zakres |
|---|---|
| `hub` | tor interfejsu: nasłuch transportu na porcie rdzenia, bez toru wykonawczego na strumieniach procesu |
| `agent` | wyłącznie tor wykonawczy: żądania kontraktu wierszami wejścia standardowego, odpowiedzi wierszami wyjścia; żaden port nie jest zajmowany |
| `all` | oba tory w jednym procesie (wartość domyślna — rdzeń lokalny) |

Rola `agent` istnieje, ponieważ przyjęta topologia rdzenia rozdziela **umiejscowienie rdzenia** od
**zasięgu wykonania modelu**: po przeniesieniu rdzenia na serwer instalacja na
urządzeniu Operatora zawiera klienta oraz agenta lokalnego, bez czego model utraciłby
dostęp do plików urządzenia.

Właściwości rdzenia potwierdzone audytem **[DZIAŁA]**:

- **Kompozycja w jednym miejscu.** Cała mapa komend powstaje w `internal/core`
  (`kompozycja.go`, `montaz.go`); port niewypełniony nie wstrzymuje startu.
- **Odwrócenie zależności transport ↔ rdzeń.** Pakiet `transport` zna wyłącznie
  interfejsy; nie importuje `core`.
- **Każde żądanie w osobnym biegu, na kontekście życia rdzenia, nie połączenia.**
  Dzięki temu `message.stop` dociera w trakcie trwającej tury, a zerwanie połączenia
  interfejsu nie przerywa pracy modelu.
- **Odporność na panikę obsługiwacza.** Panika kończy wyłącznie jedno wywołanie
  i wraca kodem błędu wewnętrznego.
- **Pełne pokrycie kontraktu.** Wszystkie 53 komendy kontraktu mają zarejestrowany
  obsługiwacz; pilnuje tego test przekroju, który przerywa budowę przy dopisaniu
  komendy bez obsługiwacza.
- **Jakość.** `gofmt -l .` i `go vet ./...` przechodzą bez uwag; 114 plików testów
  w 11 pakietach; żaden plik kodu nie przekracza twardego progu 400 wierszy.

### 3.3 Powłoka natywna — Tauri

**Decyzja.** Powłoka jest programem w języku Rust zbudowanym na Tauri 2.
Punkt wejścia `desktop/src-tauri/src/main.rs` ma 40 wierszy i wyłącznie deklaruje
moduły oraz składa budowniczego aplikacji.

Odpowiedzialności powłoki:

- **Uruchomienie rdzenia w tle** wraz z odnalezieniem binarki. Kolejność miejsc
  poszukiwania: ścieżka wskazana zmienną `DANACO_RDZEN` → katalog obok pliku
  wykonywalnego powłoki (oraz jego podkatalog `rdzen`) → miejsca w drzewie budowy
  (`budowa/`, `budowa/bin/`, `budowa/server/`) → `C:\DanacoConsole_App`. Brak binarki
  **nie wstrzymuje** otwarcia okna, a przyczyna trafia do opisu stanu
  rdzenia.
- **Rozstrzygnięcie źródła interfejsu**: adres wskazany zmienną
  `DANACO_ADRES_INTERFEJSU` → serwer rozwojowy w trybie `cargo tauri dev` → pakiet
  `client/dist` serwowany przez nasłuchujący rdzeń → pakiet osadzony w pliku
  wykonywalnym powłoki.
- **Okno główne, menu zasobnika systemowego, obsługa zamknięcia.**
- **Natywny wybór katalogu roboczego** — jedyne polecenie powłoki wywoływane dziś
  przez interfejs **[DZIAŁA]**.

**Łańcuch budowania pakietu interfejsu przez powłokę [DZIAŁA].** Plik
`desktop/src-tauri/tauri.conf.json` wiąże obie drogi uruchomienia z katalogiem
klienta, więc powłoka nigdy nie sięga po pakiet, którego nie zbudowała:

| Klucz `build` | Wartość | Skutek |
|---|---|---|
| `beforeDevCommand` | `npm run dev` w katalogu `../../client`, `wait: false` | `cargo tauri dev` podnosi serwer rozwojowy Vite w tle |
| `devUrl` | `http://localhost:5173` | okno trybu deweloperskiego bierze interfejs z serwera Vite |
| `beforeBuildCommand` | `npm run build` w katalogu `../../client` | `cargo tauri build` przebudowuje `client/dist` przed osadzeniem |
| `frontendDist` | `../../client/dist` | katalog osadzany w pliku wykonywalnym powłoki |

Skrypt `npm run build` klienta wykonuje `tsc --noEmit && vite build`, więc pakiet
nie powstanie przy błędzie typów. Pakiet wynikowy trafia trzema drogami: do
zasobu osadzonego w powłoce, do katalogu serwowanego przez rdzeń oraz — przez
skrypt wydania — do `C:\DanacoConsole_App\client\dist` (rozdział
[17.3](#173-wydanie-produktu-do-cdanacoconsole_app)). Bundle powłoki jest
skonfigurowany na cel `nsis`, ale samo złożenie instalatora wymaga narzędzia NSIS
i pozostaje krokiem odłożonym **[BRAK]**.

Ograniczenia powłoki w wersji v2.0:

- Polecenia `stan_rdzenia`, `adres_rdzenia` i `zatrzymaj_rdzen` są zarejestrowane,
  lecz interfejs ich nie wywołuje **[NIEZINTEGROWANE]**.
- Interfejs wylicza adres gniazda z `location.hostname`. Strona podana przez rdzeń
  ma `hostname` równy `127.0.0.1` i trafia w gniazdo; strona z zasobu osadzonego
  w powłoce ma `hostname` równy `tauri.localhost`, więc połączenie w pakiecie
  osadzonym nie zostaje zestawione **[DZIAŁA CZĘŚCIOWO]** — luka odnotowana.

### 3.4 Klient — interfejs TypeScript

**Decyzja.** Interfejs jest napisany w TypeScripcie i renderowany przez
webview platformy; budowany narzędziem Vite, sprawdzany `tsc --noEmit` i `vitest`.
Punkt wejścia `client/src/main.ts` ma 24 wiersze i wyłącznie komponuje trzy
czynności, z jawnie uzasadnioną kolejnością „router → transport”.

Cechy warstwy interfejsu **[DZIAŁA]**:

- **Zero literałów nazw komend i zdarzeń** w kodzie produkcyjnym — wszystkie nazwy
  pochodzą ze stałych kontraktu. Pomiar na drzewie roboczym w dniu odcięcia:
  **129 plików `.ts`** w `client/src` niesie faktyczną instrukcję importu
  z `shared/contract` (z tego 121 plików produkcyjnych i 8 plików sprawdzianów);
  kolejne dwa pliki (`rozmowa/indeks.ts`, `strona-glowna/pozycje-komponentow.ts`)
  wspominają o kontrakcie wyłącznie w komentarzu, co daje 131 plików odwołujących
  się do niego w jakiejkolwiek formie. Liczba rośnie z każdym widokiem i jest
  miarą pokrycia, nie progiem.
- **Router bez biblioteki zewnętrznej**; widok opuszczony jest porządkowany.
- **Osiągalność drzewa** — kontrola obchodu importów od `main.ts`
  (`scripts/osiagalnosc.mjs`) melduje na 467 plikach `client/src`: **435 osiągalnych,
  0 nieosiągalnych, zero zerwanych odwołań**, przy 29 zwolnieniach jawnych
  (stanowiska podglądu, sprawdziany widoków dodane w falach prac naprawczych, manifest ikon,
  licencje krojów). Osobno wykazanych jest 17 plików importowanych wyłącznie
  przez rusztowania podglądowe — jawna miara „zbudowane, niepodłączone”.
- **Motyw v2.0** — dwa równoprawne motywy (jasny i ciemny) sterowane atrybutem
  `data-theme` z uwzględnieniem `prefers-color-scheme`; **82 ikony zgodne
  z manifestem** pakietu design (`client/src/ikony/manifest.json`, pole
  `liczba-ikon: 82` i wykaz `ikony` o 82 pozycjach — obie wartości zgodne).
  W katalogu `client/src/ikony/svg/` leżą 83 pliki: 82 ikony zestawu oraz plik
  nadmiarowy wobec manifestu `logo-danaco.svg` — sygnet marki w odmianie
  pełnokolorowej, który nie jest ikoną obrysową zestawu (kod wskazuje go osobną
  stałą `IKONA_PELNOKOLOROWA`) i dlatego celowo nie podlega zasadom manifestu
  (siatka 24×24, obrys 1.75, barwa `currentColor`).

Zastrzeżenie o zasięgu: warstwa interfejsu realizuje dziś **jedno okno operacyjne**
(Chat Window) oraz okna narzędziowe konfiguracji, dostępów i modeli. Nawigacja
modułowa czterech środowisk nie przeładowuje przestrzeni roboczej — szczegóły
w rozdziałach [5](#5-moduły) i [6](#6-środowiska).

### 3.5 Kontrakt `shared/`

**Decyzja.** Katalog `budowa/shared/` jest **jedynym źródłem prawdy** nazw
komend, zdarzeń, kształtów ładunków i koperty. Źródłem jest plik `contract.json`;
generator w Node (`shared/gen/generate.mjs`) emituje z niego dwa artefakty:

- `shared/contract.go` — importowany przez rdzeń (moduł Go obejmuje `server/`
  i `shared/` jednym korzeniem, żeby import był fizycznie możliwy),
- `shared/contract.ts` — importowany przez klienta.

Zakres kontraktu wykonawczego w wersji v2.0:

| Element | Liczba |
|---|---|
| Komendy | **53** |
| Zdarzenia | **31** (12 zdarzeń dziedzinowych + 19 zdarzeń `*.unknown`) |
| Deklaracje narzędzi modelu | **39** |

Powyższe liczby opisują **zakres wykonawczy** wersji v2.0 — komendy, zdarzenia
i deklaracje obsłużone po obu stronach kontraktu. Nie są to liczby całej deklaracji
pliku `budowa/shared/contract.json`, która niesie 1119 komend, 68 obszarów,
552 struktury, 377 wyliczeń, 78 zdarzeń i 8 kodów błędów; pełną deklarację przywołuje
rozdział [19](#19-pozostała-dokumentacja) i opracowanie
[Kontrakty komunikacji](architektura/kontrakty-komunikacji.md).

Stan potwierdzony audytem **[DZIAŁA]**: statyczne odtworzenie generatora daje pliki
`contract.ts` i `contract.go` identyczne co do bajta z tymi na dysku; maszynowe
zestawienie kontraktu z opracowaniem
[Kontrakty komunikacji](architektura/kontrakty-komunikacji.md) zwraca zero rozjazdów;
w kodzie produkcyjnym po obu stronach nie ma ani jednego powielonego literału nazwy.
Sprawdzian w bramie weryfikacyjnej jest testem mutacyjnym na identyfikatorze, nie na
wartości napisu — zmiana nazwy komendy w `contract.json` musi przerwać kompilację po
obu stronach.

Rodziny komend kontraktu:

| Rodzina | Komendy |
|---|---|
| Połączenie | `connection.hello` |
| Nawigacja platformy | `home.enter` · `environment.list` · `environment.enter` · `module.list` · `workspace.enter` |
| Sesje | `session.create` · `session.list` · `session.focus` · `session.bind` · `session.open` · `session.close` · `session.delete` |
| Okna komunikacji | `window.create` · `window.list` · `window.state.get` · `window.update` · `window.close` |
| Rozmowa | `message.send` · `message.stop` · `message.list` |
| Konfiguracja | `config.get` · `config.set` · `config.reset` · `settings.category.list` · `settings.definition.list` |
| Kanały modelu | `channel.add` · `channel.update` · `channel.remove` · `channel.list` |
| Kolejki | `queue.create` · `queue.action` |
| Kontekst i akcje | `context.transfer` · `action.list` |
| Dostępy | `access.point.add/list/update/remove/check` · `access.grant.add/list/update/remove` |
| Konta | `account.add` · `account.list` · `account.update` · `account.remove` · `account.default.set` |
| Tożsamość | `identity.category.list` · `identity.document.get/set/remove` · `identity.effective.get` |

Zdarzenia dziedzinowe: `session.changed`, `session.focus.changed`, `window.changed`,
`message.changed`, `config.changed`, `queue.changed`, `stream.chunk`,
`progress.changed`, `access.point.changed`, `access.grant.changed`, `account.changed`,
`identity.changed`.

**Narzędzia modelu — [NIEZINTEGROWANE].** Generator emituje 39
deklaracji narzędzi (nazwa mechaniczna `danaco_<komenda>`, schemat wejścia wprost
z pól żądania). Poza plikami generowanymi nie mają one w bieżącej wersji ani jednego
konsumenta: sterowanie platformą przez model **nie działa**.

### 3.6 Transport: WebSocket i JSON

**Decyzja.** Klient i rdzeń rozmawiają po WebSocket, formatem JSON.
Koperta ma trzy pola: `{type, sessionId, payload}`; po stronie Go jest aliasem typu
kontraktu, nie jego kopią.

Właściwości transportu **[DZIAŁA]**:

- Port nasłuchu domyślnie **17870** (wartość dobrana świadomie poza portami zajętymi
  na serwerze `danaco-system`).
- Jedna pętla zapisu i jedna pętla odczytu na połączenie — wysyłka z wielu miejsc
  rdzenia nie miesza ramek.
- Każde żądanie obsługiwane w osobnym biegu; `message.stop` dociera w trakcie tury.
- Fail-open `*.unknown` symetryczny po obu stronach: nieznana komenda wraca
  odpowiedzią rodzaju `<obszar>.unknown` z ładunkiem opisującym, czego nie
  rozpoznano; klient prowadzi dziennik nieznanych ramek.
- Praca modelu przeżywa rozłączenie interfejsu.
- Ponawianie połączenia po stronie klienta jest wykładnicze, z pułapem
  i rozproszeniem, a ramki wychodzące są kolejkowane.

Ograniczenia transportu w wersji v2.0:

- **Ponowne połączenie nie odtwarza uzgodnienia** — po odzyskaniu łącza klient nie
  wysyła ponownie `connection.hello` ani `session.bind`, więc nowe połączenie jest
  po stronie rdzenia anonimowe **[DZIAŁA CZĘŚCIOWO]**.
- **Fragment strumienia odrzucony przy pełnej kolejce wyjściowej nie jest ponawiany**;
  tekst wypowiedzi pozostaje trwale ucięty, ponieważ interfejs nie wywołuje
  `message.list` **[DZIAŁA CZĘŚCIOWO]**.
- **Rozgłoszenie jest kierowane do wszystkich połączeń**; mechanizmy adresowania
  pojedynczego urządzenia i zawężania rozgłoszeń do konta są zadeklarowane
  w kontrakcie transportu, ale bez realizacji **[NIEZINTEGROWANE]**.
- **Brak kontroli żywotności połączenia** (ping/pong) po obu stronach **[BRAK]**.
- Wersja protokołu jest uzgadniana w `connection.hello`, lecz klient jej nie odczytuje
  **[NIEZINTEGROWANE]**.

### 3.7 Trwałość: SQLite

**Decyzja.** Trwałość opiera się na **jednym pliku SQLite**
(`danaco-console.db`) w katalogu danych rdzenia. Domyślny katalog danych to
`%LOCALAPPDATA%\DanacoConsole`; poza Windows rdzeń sięga po katalog konfiguracyjny
użytkownika, aby start nigdy nie zależał od obecności jednej zmiennej środowiska.

Właściwości warstwy danych **[DZIAŁA]**:

- **Migracje wkompilowane w binarium** (`go:embed`), stosowane rosnąco, każdy krok
  w jednej transakcji, z sumą kontrolną wykrywającą zmianę treści już zastosowanego
  kroku i z rejestrem wersji schematu.
- **Pragmy przez DSN** — `foreign_keys(1)`, `journal_mode(WAL)`, `busy_timeout` —
  obowiązują każde połączenie z puli.
- **Przekład wyliczeń wyłącznie słownikami kontraktu** — warstwa danych nie powiela
  ani jednego literału wyliczenia.
- **Odtworzenie sesji i okien z bazy po restarcie rdzenia** pod tymi samymi
  identyfikatorami.

Druga część decyzji o trwałości — „treści obszerne jako pliki, baza trzyma ścieżkę” — jest
w bieżącej wersji **[BRAK]**: kolumny `tresc_odwolanie` mają w repozytoriach pełną
ścieżkę zapisu i odczytu, ale nigdy nie są wypełniane; nie ma zapisu treści do pliku.

### 3.8 Decyzje architektoniczne

Zestawienie decyzji, które bezpośrednio kształtują opisywane dalej mechanizmy:

| Decyzja | Status |
|---|---|
| Rdzeń serwera w języku Go | Przyjęta |
| Klient jako natywne okno Tauri (warstwa Rust) | Przyjęta |
| Interfejs w TypeScript renderowany przez webview platformy | Przyjęta |
| Trwałość na SQLite (jeden plik, transakcje); treści obszerne jako pliki | Przyjęta |
| Komunikacja klient–serwer po WebSocket, format JSON | Przyjęta |
| Model hybrydowy: serwer wykonawczy + cienki klient urządzenia | Przyjęta |
| Interface-first — rozdzielenie przez kontrakty | Przyjęta |
| Jedna odpowiedzialność = jeden plik; cienki punkt wejścia | Przyjęta |
| Bez wartości wpisanych na sztywno w nakładce: prompt, harness, hooki, model wyłącznie z konfiguracji | Przyjęta |
| Konfiguracja zamiast logiki; brak ustawienia = wartość domyślna | Przyjęta |
| Kanał główny Code CLI (`CLAUDE_CONFIG_DIR`, rotacja kont) + API/Agent SDK | Przyjęta |
| Warstwy nakładki wg krytyczności: konstytucja → profil → ekspertyza | Przyjęta |
| Fail-open — element opcjonalny nie blokuje uruchomienia | Przyjęta |
| Suma kontrolna konstytucji wyłącznie do diagnostyki prowenancji | Przyjęta |
| Błąd techniczny = błąd tylko bieżącego wywołania | Przyjęta |
| Prowenancja „co poszło do modelu” — transparentność, nie bramka | Przyjęta |
| Zero blokad w interfejsie | Przyjęta |
| Uwierzytelnianie jedyną kontrolą dostępu; w fazie budowy wyłączone | Przyjęta |
| Sekrety poza repozytorium (`.env`, repozytorium trzyma `.env.example`) | Przyjęta |
| Wersjonowanie: v2.0 do pierwszej publikacji, status Deweloperski | Przyjęta |
| `shared/` jedynym źródłem prawdy kontraktu; koperta; fail-open `*.unknown` | Przyjęta |
| Jednolity kontrakt kanałów; wspólny strumień fragmentów, `Chunk` niesie `Kind` | Przyjęta |
| Otwarty rejestr kanałów sterowany danymi — nowy kanał = nowy wiersz | Przyjęta |
| Topologia rdzenia; zasięg wykonania niezależny od umiejscowienia rdzenia | Przyjęta |
| Okno komunikacji jako byt pośredni między sesją a wiadomością | Przyjęta |
| Sterowanie platformą przez narzędzia modelu generowane z kontraktu | Przyjęta |
| Pętla koordynator–wykonawca na jednym silniku kolejek | Przyjęta |
| Bieg naprawczy bez limitu; przerwanie należy do Operatora | Przyjęta |
| Koordynator widzi pełny strumień wykonawcy | Przyjęta |
| Telemetria postępu w kontrakcie — jedno zdarzenie dla dwóch monitorów | Przyjęta |
| Przenoszenie kontekstu jako cecha systemowa (jedna komenda, nie ścieżka per moduł) | Przyjęta |
| Rejestr akcji sterowany danymi; Panel akcji budowany dynamicznie | Przyjęta |
| Progi modularności (plik ≤ 200/400, funkcja ≤ 50/100, wejście ≤ 100, arkusz ≤ 300) | Przyjęta |
| Tryb budowy: fale subagentów o rozłącznych zestawach plików, zamykane bramą | Przyjęta |
| Katalog rozszerzeń w zakresie; marketplace wydawniczy odłożony | Przyjęta |
| Ósmy poziom zasięgu konfiguracji: okno komunikacji | Przyjęta |
| Zakaz atrap — kod niepodłączony do punktu wejścia nie liczy się jako wykonany | Przyjęta |
| Zasięg wykonania jest rodzajem (`local` · `core` · `remote`), nie nazwą hosta | **Zaproponowana** |
| Pakiet design v2.0 zastępuje w całości dotychczasową warstwę wizualną | Przyjęta |
| Makieta określa wygląd, dokumentacja określa zakres funkcji | Przyjęta |
| `model.channel.set` zastąpiona rodziną `channel.*` | **Zaproponowana** |
| `window.changed` nazwą kanoniczną zdarzenia zmiany okna | **Zaproponowana** |
| Zakres pierwszego przekroju pionowego | **Zaproponowana** |
| Konwencja narzędzi modelu generowanych z kontraktu | **Zaproponowana** |
| Środowisko ≠ punkt dostępu ≠ katalog roboczy — trzy rozłączne byty | Przyjęta |
| Dostęp jest **zbiorem** nadań per okno rozmowy | Przyjęta |
| Dosłownie wszystko konfigurowalne; katalog ustawień sterowany danymi; osie | Przyjęta |
| Tożsamość modelu jest **zamieniana**, nie dołączana (domyślnie `ZASTAP`) | Przyjęta |
| Start zawsze na stronie głównej: `connection.hello` → `home.enter` | Przyjęta |
| Oś ustawienia jest prostopadła do poziomu zasięgu | **Zaproponowana** |
| Poświadczenie wchodzi żądaniem i nie wychodzi nigdy | **Zaproponowana** |
| Sterowanie silnikiem nakładki z danych katalogu; `ZASTAP` wygrywa z `DOLACZ` | **Zaproponowana** |
| Most MCP po SSH jako droga dostępu modelu do maszyn | **Zaproponowana** |

Zestawienie podaje brzmienie skrócone; pełne uzasadnienia i zastrzeżenia każdej
z decyzji opisują rozdziały poświęcone odpowiednim mechanizmom.

### 3.9 Reguły modularności i brama weryfikacyjna

Produkt jest budowany falami subagentów o rozłącznych zestawach plików,
a każdy etap zamyka **brama weryfikacyjna pięciu kontroli**: zgodność
z koncepcją, kryterium etapu na dowodzie, kontrakt (obie strony importują
z `shared/`), kompletność (zakaz atrap i kodu niepodłączonego), metryka wielkości.

Progi modularności:

| Byt | Cel | Próg twardy |
|---|---|---|
| Plik kodu (`.go`, `.ts`, `.rs`) | ≤ 200 wierszy | **400 wierszy** |
| Funkcja / metoda | ≤ 50 wierszy | 100 wierszy |
| Punkt wejścia | wyłącznie kompozycja | **100 wierszy, zero logiki** |
| Arkusz stylów | ≤ 300 wierszy | podział per komponent |

Wyłączone jawnie: pliki generowane (`contract.ts`, `contract.go`), migracje SQL,
pliki zasobów. Zakazane są nazwy `utils`, `helpers`, `common`, `misc`.

Bramę realizuje **piętnaście skryptów** w `budowa/scripts/` oraz podkatalog
`haki/`. Pełny wykaz, w podziale na role:

| Skrypt | Rola |
|---|---|
| `brama.sh` | przebieg zbiorczy bramy; tryb `--szybka` i tryb pełny (Go + Node + Rust) |
| `metryka-plikow.sh` | progi modularności: wielkość pliku, funkcji, punktu wejścia, arkusza stylów |
| `atrapy.sh` | kod bez odbiorcy i znaczniki niedokończenia po stronie repozytorium |
| `atrapy-eksporty.mjs` | eksporty TypeScript bez konsumenta |
| `eksporty-go.mjs` | eksporty Go bez konsumenta |
| `eksporty-go-zwolnienia.mjs` | jawny wykaz zwolnień dla kontroli eksportów Go |
| `osiagalnosc.mjs` | osiągalność drzewa interfejsu obchodem importów od `main.ts` |
| `osiagalnosc-odwolania.mjs` | zerwane odwołania w drzewie interfejsu |
| `wywolania-kontraktu.mjs` | które komendy kontraktu mają realne wywołanie po stronie klienta |
| `dym-pionu.mjs` | dymny sprawdzian całego pionu na kanale `echo` — od `connection.hello` po strumień odpowiedzi |
| `kontrola-rust.sh` | warstwa powłoki: formatowanie, `clippy`, budowa |
| `wydanie.sh` | wydanie produktu: klient → rdzeń → powłoka → pakowanie |
| `pakowanie.sh` | złożenie katalogu produktu z gotowych artefaktów |
| `pokaz.sh` | podgląd deweloperski stanowisk |
| `kopia.sh` | kopia robocza drzewa |

Trzy pozycje wykazu domykają szczeliny wytknięte audytem. `eksporty-go.mjs`
wraz z wykazem zwolnień wykrywa eksporty Go bez odbiorcy; `wywolania-kontraktu.mjs`
bada **wywołanie**, a nie sam import, więc funkcja zbudowana i nigdy niewywołana
przestaje przechodzić niezauważona; `dym-pionu.mjs` przechodzi cały pion na kanale
`echo` bez sieci i bez konta, w środowisku Node bez zależności zewnętrznych.

**Brama jako proces [DZIAŁA CZĘŚCIOWO].** Hak `scripts/haki/pre-commit` uruchamia
bramę w trybie szybkim przed każdym zatwierdzeniem i jest w tym drzewie
zainstalowany (`git config core.hooksPath budowa/scripts/haki`); jego pominięcie
(`git commit --no-verify`) pozostaje decyzją Operatora, nie stanem domyślnym.
Przepływ ciągły `.github/workflows/brama.yml` uruchamia bramę pełną na `push`
i `pull_request` — jest przygotowany, lecz **bezczynny**, ponieważ projekt nie ma
dziś zdalnego repozytorium GitHub; zadziała z chwilą pierwszego wypchnięcia.

Stan bramy zmierzony na drzewie roboczym w dniu odcięcia: zero naruszeń progu
twardego (największy plik kodu liczy 254 wiersze i jest plikiem sprawdzianu;
największy plik produkcyjny — 233 wiersze), zero znaczników `TODO`/`FIXME`, zero
plików-wysypisk, zgodność `.env.example` z kodem co do joty. Pokrycie sprawdzianami:
114 plików testów Go i 10 plików `*.test.ts` po stronie klienta.

---

## 4. Główne komponenty

### 4.1 Pakiety rdzenia

Rdzeń dzieli się na dziesięć pakietów wewnętrznych o rozłącznych
odpowiedzialnościach (`budowa/server/internal/`):

| Pakiet | Odpowiedzialność |
|---|---|
| `transport` | nasłuch WebSocket, nawiązanie połączenia, pętle odbioru i wysyłki, rozgłoszenie, serwowanie pakietu interfejsu |
| `protocol` | koperta, odkodowanie żądania, rejestr nazw komend, fail-open `*.unknown` |
| `core` | kompozycja rdzenia, rejestr obsługiwaczy, adaptery domen (rozmowa, sesje, okna, kanały, konfiguracja, dostępy, konta, tożsamość, kolejki, nawigacja), dziennik rozmowy, telemetria, most MCP |
| `session` | nadzorca sesji i okien, rejestr bytów w pamięci procesu, pętla koordynator–wykonawca, podsystem procesów okien |
| `models` | rejestr kanałów modelu, wspólny interfejs kanału, adaptery `cli`/`api`/`echo`, zapytanie kanału, prowenancja, ujście strumienia |
| `injection` | budowa argumentów procesu modelu, nakładka tożsamości, pula kont i rotacja, uruchomienie procesu, parser strumienia `stream-json`, rozpoznanie wyczerpania limitu |
| `konfig` | ośmiopoziomowy rozstrzygacz konfiguracji z osiami, definicje wykonania i izolacji, rejestr definicji z katalogu bazy |
| `konfiguracja` | konfiguracja **startowa procesu** rdzenia: rola, port, katalogi; warstwy domyślne → środowisko → argumenty |
| `dane` | repozytoria warstwy trwałości, przekład wyliczeń, łańcuch więzów rozmowy |
| `store` | otwarcie bazy, pragmy, migracje wkompilowane, kontrola spójności |

Rozróżnienie `konfig` ↔ `konfiguracja` jest celowe i nie stanowi duplikacji:
pierwszy pakiet rozstrzyga ustawienia Operatora w czasie pracy, drugi ustala
parametry startu procesu.

### 4.2 Obszary klienta

Interfejs dzieli się na obszary odpowiadające widokom i warstwom
(`budowa/client/src/`):

| Obszar | Zawartość |
|---|---|
| `aplikacja` | router, sceny, widoki trasy, komunikaty (dymki), arkusze powłoki |
| `powloka` | pasek górny, karty sesji, boczna nawigacja modułów, obszar roboczy, katalog środowisk |
| `strona-glowna` | Centrum dowodzenia: kafle środowisk i komponentów |
| `rozmowa`, `okno-komunikacji` | okno komunikacji: nadawca, historia, przepływ komunikatów, metadane konta |
| `okna-rownolegle` | scena 1/2/3 okien, gniazda okien, pas relacji, przekazanie zlecenia |
| `sterowanie` | panel sterowania okna: rola, kanał, zasięg i host wykonania, nakład rozumowania |
| `konfiguracja` | okno konfiguracji sterowane katalogiem z rdzenia |
| `dostepy`, `modele` | punkty dostępu i nadania; kanały, konta i tożsamość |
| `mission-control` | pulpit operacyjny: źródło danych z rdzenia, złożenie kompletu, sekcje pulpitu |
| `protokol`, `polaczenie` | koperta, korelacja żądań, komendy platformy, transport, dziennik nieznanych ramek |
| `motyw`, `ikony`, `komponenty` | żetony motywu, zestaw 82 ikon wraz z sygnetem marki, biblioteka komponentów v2.0 |

### 4.3 Moduły powłoki

Powłoka Rust liczy 17 plików w podziale: `main.rs` (kompozycja), `montaz.rs`,
`okno.rs`, `zasobnik.rs`, `menu_zasobnika.rs`, `zamkniecie.rs`, `polecenia.rs`,
`ustawienia.rs`, `dialog_katalogu.rs`, `zrodlo_interfejsu.rs` oraz podkatalog
`rdzen/` (`mod.rs` — deklaracja modułu, `lokalizacja.rs`, `uruchomienie.rs`,
`nasluch.rs`, `uchwyt.rs`, `dziennik.rs`, `pakiet_klienta.rs`).

---

## 5. Moduły

**Moduł** jest wyspecjalizowaną przestrzenią roboczą. Katalog modułów jest sterowany
danymi: piętnaście wierszy tabeli `modul` zasianych migracją 007, wraz z macierzą
widoczności `srodowisko_modul` (9 pozycji w TalkIn, 9 w WorkSpace, 8 w CodeStudio).
Rdzeń wystawia ten katalog komendami `module.list` i `environment.enter`
**[DZIAŁA]**.

**Stan realizacji warstwy interfejsu jest jednak jednoznaczny: żaden z piętnastu
modułów nie ma własnego widoku [ATRAPA].** Boczna nawigacja wystawia komplet
pozycji, ale wybór pozycji zmienia wyłącznie tytuł karty i podpis w pasku
kontekstu — obszar roboczy jest trwale zajęty przez scenę okien komunikacji.
Wszystkie 32 pozycje nawigacji czterech środowisk prowadzą do tego samego ekranu:
26 pozycji modułowych trzech środowisk modułowych (TalkIn 9, WorkSpace 9,
CodeStudio 8) oraz sześć sekcji panelu orkiestracji środowiska MultitaskingAI
([INSTRUKCJA-UZYTKOWANIA.md](INSTRUKCJA-UZYTKOWANIA.md) rozdz. 4.2).
Z piętnastu prototypów okien modułowych w `design/05-okna/moduly/` — jednego na moduł,
część pełnego zestawu **38 prototypów** interfejsu przechowywanych w `design/05-okna/`
(zliczenie: `find design/05-okna -name '*.html' | wc -l`) — wdrożone jest **jedno**:
Chat Window, wspólne wszystkim modułom.

| Moduł | Przeznaczenie wg makiety | Stan w v2.0 |
|---|---|---|
| **Studio** | praca z dokumentem: edycja, operacje kontekstowe, wersjonowanie sesji | **[ATRAPA]** — brak Studio Editor, Tools Panel, Diff/Grep, Session Repository |
| **Workspace** | projekt: zadania, zasoby, instrukcje systemowe, pamięć kontekstu | **[ATRAPA]** — brak Project Dashboard, Instructions Panel, Context Memory, Project Library |
| **Automations** | przebieg pracy, harmonogram, kolejki, orkiestracja | **[BRAK]** — zgodnie z koncepcją bez okna modułowego, ale bez poprawnego wejścia ze strony głównej |
| **Browser** | przeglądanie stron z gromadzeniem źródeł i notatek | **[ATRAPA]** — brak Browser Window, Sources Panel, Notes Panel |
| **Research** | badanie: źródła, ustalenia, raport | **[ATRAPA]** — brak Research Workspace, Sources Manager, Findings Panel, Report Builder |
| **Library** | repozytorium plików, wersji, etykiet i kolekcji | **[ATRAPA]** — brak Library Explorer, Tags & Collections, File Preview, Versioning Panel |
| **Translate** | tłumaczenie równoległe z glosariuszem | **[ATRAPA]** — brak Source Panel, Translation Panels, Glossary Manager |
| **Roundtable** | debata wielu modeli z moderacją | **[ATRAPA]** — brak Model Panels, Debate Panel, Moderator Panel, Consensus Panel |
| **Design** | praca wizualna: kompozycja, zasoby, generowanie | **[ATRAPA]** — brak Design Board, Assets Panel, Prompt Builder, Preview Window |
| **Assistant** | asystent głosowy z historią działań | **[ATRAPA]** — brak Voice Console, Actions Monitor, Activity Feed; brak toru przetwarzania: rozpoznawanie mowy → model → synteza mowy |
| **Terminal** | powłoki, polecenia i procesy urządzenia | **[ATRAPA]** — brak Terminal Tabs, Output Console, Process Monitor |
| **Developer** | kod, budowanie, repozytorium Git | **[ATRAPA]** — brak okien modułowych i komend `developer.*` |
| **Diagnostics** | stan systemu, logi i analiza błędów | **[ATRAPA]** — brak okien modułowych |
| **Apps** | budowa produktu od architektury po wdrożenie | **[ATRAPA]** — brak okien modułowych |
| **Agents** | agenci: tożsamość, model, uprawnienia, konektory | **[ATRAPA]** — brak Agent Builder, Skills Manager, Connectors Manager, Permissions Center |

Dodatkowe zastrzeżenia wobec warstwy modułowej:

- **Rejestr akcji [NIEZINTEGROWANE].** Rdzeń ma pełny, sterowany danymi
  katalog akcji (migracje 008 i 009) oraz komendę `action.list`; interfejs nie
  wywołuje jej ani razu i nie ma Panelu akcji. Zaczyn obejmuje 29 komend platformy,
  a nie ponad 260 akcji modułowych z inwentarza.
- **Pojęcie modułu w kontrakcie [DO DECYZJI OPERATORA].** Lista `KnownModuleIds`
  kontraktu niesie kody **środowisk** (`talkin`, `workspace`, `codestudio`,
  `multitaskingai`), a jest jednocześnie źródłem domyślnego modułu okna i listy
  wyboru w panelu sterowania. Skutkiem jest zapisywanie okien zawsze z modułem
  `workspace`, niezależnie od nawigacji. Rozstrzygnięcie — przemianowanie listy
  w kontrakcie albo dodanie osobnej listy kodów modułów — należy do Operatora.
- **Kliknięcie modułu bez widoku nie zwraca komunikatu**, co jest odstępstwem od
  wzorca zera blokad (element klikalny, kliknięcie zwraca informację zamiast prowadzić
  donikąd).

---

## 6. Środowiska

**Środowisko** jest — zgodnie z przyjętym rozdziałem pojęć — **profilem widoczności modułów**
w bocznej nawigacji. Nie jest pojemnikiem na maszyny ani na katalogi; opis maszyn
należy do punktów dostępu, a miejsce pracy modelu do katalogów roboczych.

Cztery środowiska są zasiane migracją 007 i wystawiane komendami
`environment.list` / `environment.enter` **[DZIAŁA po stronie rdzenia]**:

| Środowisko | Motto | Przeznaczenie | Moduły |
|---|---|---|---|
| **TalkIn** | Myśl. Analizuj. Rozumiej. | wiedza, komunikacja i praca z treścią | 9 pozycji |
| **WorkSpace** | Planuj. Organizuj. Realizuj. | produktywność, organizacja i realizacja projektów | 9 pozycji |
| **CodeStudio** | Projektuj. Buduj. Rozwijaj. | programowanie | 8 pozycji |
| **MultitaskingAI** | Deleguj. Koordynuj. Nadzoruj. | orkiestracja autonomicznej pracy ciągłej | panel orkiestracji (6 sekcji) |

Stan realizacji funkcjonalnej:

- **TalkIn / WorkSpace / CodeStudio — [ATRAPA].** Nawigacja i motta istnieją; żadna
  pozycja nie otwiera przestrzeni roboczej modułu. Zmiana środowiska prowadzi przez
  stronę główną, zgodnie z koncepcją, ale karty opuszczonego środowiska nie trwają
  jako byty rdzenia — pas kart wewnątrz środowiska jest dziś pasem okien sceny,
  najwyżej trzech, i nie wywołuje `session.close`. Sesje trwające na rdzeniu są
  natomiast widoczne i osiągalne ze **strony głównej**, która odpytuje `session.list`,
  nasłuchuje `session.changed` i wykonuje powrót ścieżką `session.bind` →
  `session.focus` **[DZIAŁA]**.
- **MultitaskingAI — [ATRAPA].** Środowisko ma w interfejsie sześć etykiet bocznej
  listy (Zespoły, Role, Kolejki, Orkiestracja, Harmonogram, Monitor) obsługiwanych
  tym samym kodem co lista modułów. Panelu orkiestracji, czterech okien ról,
  Subagent Network, hierarchii ról, harmonogramu i monitora procesu nie ma ani
  w interfejsie, ani w kontrakcie wykonawczym.

**Mission Control (pulpit operacyjny) — [DZIAŁA CZĘŚCIOWO].** Ustalenie audytu
(„rysuje się wyłącznie z danych przykładowych”) **przestało obowiązywać**. Plik
`mission-control/dane-przykladowe.ts` — 2143 wiersze liczb wpisanych w kod — został
z drzewa usunięty, a jego miejsce zajęły cztery moduły: `zrodlo-pulpitu.ts`
(odczyty kontraktu i subskrypcje zdarzeń), `zlozenie-danych.ts` (czyste przełożenie
stanu na komplet widoku), `stan-zrodla.ts` (stan surowy) oraz `stan-pusty.ts` (stan
sprzed pierwszego odczytu). Pulpit pozostaje trzecią trasą aplikacji, obok Centrum
dowodzenia. Stan bieżący, zweryfikowany w kodzie:

| Element pulpitu | Źródło | Stan |
|---|---|---|
| Matryca sesji, kolumna zespołu, kafle aktywności | `session.list` (z obecnością), `window.list` | **[DZIAŁA]** |
| Sekcja operacji (kanały modelu) | `channel.list` | **[DZIAŁA]** |
| Kolumna procesów | zdarzenie `progress.changed` | **[DZIAŁA]** — pulpit jest jedynym konsumentem tej telemetrii |
| Kolumna kolejek | zdarzenie `queue.changed` | **[DZIAŁA CZĘŚCIOWO]** — kontrakt nie ma odczytu `queue.list`, więc stan buduje się wyłącznie ze zdarzeń |
| Wejście do sesji z matrycy | `session.open` | **[DZIAŁA]** |
| Pasek transportu kolejki (uruchom, wstrzymaj, zatrzymaj, ponów, wznów, wyczyść) | `queue.action` | **[DZIAŁA CZĘŚCIOWO]** — komenda dochodzi, lecz przestawia wyłącznie kolumnę stanu (patrz [9.5](#95-kolejki)) |
| Przycisk „Przekaż” | `context.transfer` | **[NIEZINTEGROWANE]** — komenda wymaga kompletu kontekstu, którego pas kolejek nie niesie |
| Przycisk „Podnieś priorytet” | — | **[BRAK]** — czynność bez odpowiednika w kontrakcie |
| Przyciski sekcji „Utwórz”, pas decyzji | — | **[BRAK]** — kreatory poza zakresem bieżącej fali |

Pulpit nie niesie ani jednej liczby wymyślonej: miara bez źródła w kontrakcie
zostaje pusta, a nagłówek pulpitu wprost oznacza, czy komplet pochodzi już
z rdzenia, czy dopiero z oczekiwania na pierwszy odczyt. Wartości pozostające
trwale puste, bo kontrakt nie ma dla nich producenta, to: modele analizujące,
walidatory oczekujące, wstrzymane decyzje, pas relacji i liczba podagentów.
Trzy czynności bez drogi przez kontrakt nie są udawane — naciśnięcie przycisku
zwraca komunikat mówiący wprost, czego brakuje.

Pozostałe zastrzeżenia warstwy środowisk:

- **Ścieżka nawigacji platformy — [NIEZINTEGROWANE w zakresie sześciu komend].**
  Twierdzenie audytu o komplecie ośmiu komend bez konsumenta wymaga zawężenia.
  Dwie z nich mają dziś konsumenta w warstwie widoku: `client/src/strona-glowna/wpiecie-sesji.ts`
  wykonuje powrót do sesji trwającej na rdzeniu ścieżką `session.bind` →
  `session.focus` (moduły protokołu `powiazanie-sesji` i `ognisko-sesji`), a wykaz
  sesji strony głównej odpytuje `session.list` i nasłuchuje `session.changed`
  **[DZIAŁA]**. Bez konsumenta pozostaje **sześć komend**: `home.enter`,
  `environment.list`, `environment.enter`, `module.list`, `workspace.enter`
  i `window.state.get`. Są zaimplementowane po obu stronach i złożone w jeden byt
  warstwy protokołu (`protokol/komendy-platformy.ts`, `protokol/nawigacja-platformy.ts`,
  `protokol/stan-okna.ts`), lecz żaden widok nie wywołuje ani jednej z ich metod.
  Ruch otwierający interfejsu biegnie `connection.hello` → `session.create` →
  `window.create`, wbrew przyjętej regule startu, która przewiduje
  `connection.hello` → `home.enter`.
- **Katalog środowisk powielony w interfejsie — [DZIAŁA CZĘŚCIOWO].** Powłoka
  interfejsu trzyma własny, zaszyty w kodzie rejestr czterech środowisk wraz z pełną
  macierzą widoczności, równolegle do słownika w bazie.

Nazwy, motta i macierz widoczności odpowiadają dokumentacji — **zredukowana została
realizacja, nie koncepcja**.

---

## 7. Modele i kanały wykonawcze

### 7.1 Rejestr kanałów sterowany danymi

**Decyzja.** Kanał modelu jest **wierszem tabeli `kanal_modelu`**,
nie typem w kodzie. Wiersz niesie: kod, nazwę, dostawcę, identyfikator modelu,
rodzaj kanału, dowiązanie do konta, odwołanie do poświadczenia, parametry
w formacie JSON, znacznik multimodalności, aktywność i kolejność.

**Rodzaj kanału — dwa rozłączne pojęcia pod jedną nazwą.** Jest to miejsce, w którym
dokumentacja produktu bywała niespójna, więc dokument rozdziela je wprost:

| Pojęcie | Zbiór wartości | Miejsce | Rola |
|---|---|---|---|
| **Wartość kolumny schematu** `kanal_modelu.rodzaj_kanalu` | `cli` · `api` · `sdk` · `lokalny` | więz `CHECK` migracji 001; ta sama czwórka w kontrakcie jako `KnownChannelKinds` | klasyfikacja wiersza rejestru — czym kanał **jest** |
| **Rodzaj kanału wykonawczego** (klucz adaptera) | `cli` · `api` · `echo` | fabryki adapterów w rdzeniu (`models.FabrykiWbudowane` + adapter kanału głównego) | rozstrzygnięcie, **którym kodem** tura zostanie wykonana |

Klucz adaptera rozstrzyga metoda `Definicja.KluczAdaptera`: **wygrywa parametr
`adapter` wiersza, a dopiero w jego braku kolumna `rodzaj_kanalu`**. Obie wartości
są danymi, więc kanał na nowym adapterze powstaje wpisem do tabeli, nie zmianą kodu.
Praktyczna konsekwencja tego rozdziału: kanał `echo` **nie jest**
dopuszczalną wartością kolumny schematu — zakłada się go wierszem o rodzaju
`lokalny` z parametrem `{"adapter": "echo"}`. Dokładnie tak postępuje dymny
sprawdzian pionu `scripts/dym-pionu.mjs`.

Terminologia obowiązująca w pozostałej dokumentacji produktu
([INSTRUKCJA-UZYTKOWANIA.md](INSTRUKCJA-UZYTKOWANIA.md) rozdz. 1.4 i 6.1) posługuje się
trójką **cli / api / echo** — i jest to trójka kluczy adaptera, czyli rodzajów **kanału
wykonawczego**, zgodna z powyższą tabelą. Czwórkę `cli · api · sdk · lokalny` opisuje
[LICENSE.md](LICENSE.md) rozdz. 2.3 jako dopuszczalne wartości kolumny schematu,
zastrzegając wprost, że `echo` nie jest samodzielną wartością tej kolumny.
**[DO DECYZJI OPERATORA]:** czy oba zbiory
mają zostać rozdzielone także nazewniczo w kontrakcie (`KnownChannelKinds` wobec
osobnego wykazu kluczy adapterów), czy różnica pozostaje wyjaśniona wyłącznie
dokumentacją. Rodzaje `sdk` i `lokalny` nie mają dziś własnego adaptera
wykonawczego **[BRAK]**.

Rejestr kanałów w rdzeniu jest budowany z wierszy tabeli i **odświeżany po każdej
zmianie**; instancje kanałów niezmienionych są zachowywane. Wszystkie kanały mówią
jednym strumieniem fragmentów, a fragment niesie swój rodzaj.

Cykl życia rejestru obsługują cztery komendy: `channel.add`, `channel.update`,
`channel.remove`, `channel.list` **[DZIAŁA po stronie rdzenia]**.

**Zastrzeżenie [BRAK]:** interfejs nie ma okna zakładania ani edycji kanałów —
w kodzie produkcyjnym wywoływana jest wyłącznie `channel.list` (podpowiedź modeli
i lista wyboru w panelu sterowania okna). Operator nie ma w produkcie drogi, żeby
założyć kanał modelu; pozostaje komenda kontraktu wysłana po WebSocket.

**Zastrzeżenie [DO DECYZJI OPERATORA]:** identyfikator kanału w kontrakcie jest
niespójny — `channel.add` i `channel.update` posługują się kodem wiersza, natomiast
`channel.list` oddaje wykaz z rejestru modeli, w którym identyfikatorem jest numer
wiersza. Wymaga to jednego rozstrzygnięcia w `shared/`.

### 7.2 Kanał główny — Claude Code CLI

**Decyzja.** Kanałem głównym produktu jest **Claude Code CLI** — zewnętrzny
program `claude` uruchamiany przez rdzeń jako proces potomny. Kanał jest zbudowany
kompletnie i starannie **[DZIAŁA CZĘŚCIOWO — patrz zastrzeżenia]**.

**Szkielet każdego wywołania** (stałe argumenty):

```
-p  --input-format stream-json  --output-format stream-json  --verbose
```

Tryb bezinterakcyjny, wejście i wyjście strumieniem JSON-lines oraz pełne wyjście
zdarzeń — bez `--verbose` program nie wypuszcza zdarzeń pośrednich, więc strumień
przestałby być strumieniem.

**Argumenty warunkowe.** Przełącznik pojawia się wyłącznie wtedy, gdy odpowiadające
mu ustawienie jest wypełnione; brak ustawienia nie jest błędem, tylko brakiem
przełącznika:

| Przełącznik | Źródło wartości | Stan |
|---|---|---|
| `--permission-mode` | tryb uprawnień okna komunikacji | **[DZIAŁA]** |
| `--add-dir` (powtarzalny) | lista katalogów roboczych okna | **[DZIAŁA]** |
| `--mcp-config` (powtarzalny) | konfiguracja mostów MCP złożona ze zbioru nadań okna | **[DZIAŁA]** |
| `--system-prompt` / `--append-system-prompt` | nakładka tożsamości, tryb `ZASTAP` / `DOLACZ` | **[DZIAŁA]** |
| `--model` | wskazanie okna, a w jego braku identyfikator modelu z wiersza rejestru kanału | **[DZIAŁA CZĘŚCIOWO]** — wartość z wiersza rejestru kanału; wskazanie per okno nie dociera |
| `--fallback-model` | pole zapytania kanału | **[NIEZINTEGROWANE]** — pole bez producenta |
| `--effort` | nakład rozumowania | **[NIEZINTEGROWANE]** — pole bez producenta |
| `--settings` | plik ustawień powłoki wykonawczej | **[NIEZINTEGROWANE]** — pole bez producenta |
| `--resume` | identyfikator rozmowy programu | **[NIEZINTEGROWANE]** — przełącznik zbudowany, pole zapytania bez producenta |

**Uściślenie do wiersza `--model`.** Łańcuch jest domknięty po stronie rejestru
i przerwany po stronie okna, a rozstrzyga to jedna metoda:
`server/internal/models/zapytanie.go` — `WybranyModel(d)` zwraca wskazanie okna,
a w jego braku model z wiersza rejestru kanałów (`d.Model`). Wynik trafia przez
`server/internal/core/adapter_kanal_cli.go` do pola `Model` ustawień wstrzyknięcia,
a `server/internal/injection/argumenty.go` dopisuje z niego przełącznik `--model`.
**Przełącznik jest zatem podawany zawsze, ilekroć wiersz `kanal_modelu` niesie
identyfikator modelu** — a niesie go zawsze, bo kolumna `identyfikator_modelu`
jest w schemacie obowiązkowa. Nie dociera natomiast **wskazanie per okno**:
funkcja `zapytanieKanalu` w `server/internal/core/adapter_rozmowa.go` składa
zapytanie z parametrów okna i pola `Model` nie wypełnia. Model wywołania jest więc
własnością **kanału**, nie okna; zmiana modelu wymaga wskazania oknu innego wiersza
rejestru, a nie zapisu ustawienia. Ta sama luka producenta dotyczy pól
`ModelZapasowy`, `NakladRozumowania` i `Wznowienie` — z tą różnicą, że dla nich
rejestr nie podstawia żadnej wartości zastępczej, więc odpowiadające im przełączniki
nie pojawiają się nigdy.

**Katalog startowy procesu** pochodzi z katalogu roboczego sesji **[DZIAŁA]**.

**Uwierzytelnianie i profil konta.** Kanał ustawia procesowi zmienną
`CLAUDE_CONFIG_DIR` wskazującą katalog profilu wybranego konta. Katalog konfiguracji
jest w tym produkcie wyłącznie nośnikiem tożsamości; poświadczenia zostają na dysku,
poza repozytorium i poza bazą.

**Parser strumienia** czyta wyjście `stream-json` wiersz po wierszu, jest odporny na
linie nieczytelne (zlicza je i wykazuje w podsumowaniu tury), rozpoznaje wyczerpanie
limitu tempa z trzech niezależnych źródeł i wychwytuje identyfikator rozmowy programu
ze zdarzenia wyniku.

**Ścieżka programu.** Rdzeń używa ścieżki z wiersza rejestru kanału albo stałej
`claude` (program musi być osiągalny w `PATH`). Klucz katalogu ustawień
`harness.program_claude` istnieje i jest opisany jako obowiązujący, gdy wiersz
rejestru ścieżki nie niesie — ale **nie jest odczytywany** **[NIEZINTEGROWANE]**.

### 7.3 Kanał `api`

Kanał generyczny HTTP (`models/adapter_api.go`, `models/zadanie_api.go`) jest
zbudowany poprawnie i **w pełni sparametryzowany wierszem rejestru** — nie zna
żadnego zaszytego dostawcy. Czytnik strumienia obsługuje zarówno linie SSE, jak
i gołe linie JSON; wyjmowanie treści odbywa się ścieżką po dokumencie JSON zapisaną
w parametrach wiersza.

**Stan: [NIEZINTEGROWANE].** Kanał nie ma drogi wykonawczej przez katalog kont:
konta rodzajów `api` i `sdk` mają w schemacie dostawcę, model domyślny i adres bazowy,
ale do wykonania dociera wyłącznie rodzaj `cli`. Kanał `api` ignoruje też tryb
nakładki — prompt jest zawsze dopisywany jako wiadomość systemowa.

### 7.4 Kanał `echo`

Kanał testowy bez sieci, realizujący pełną drogę: prowenancja → konto → porcje
odpowiedzi → domknięcie strumienia. Służy do sprawdzenia toru komunikacyjnego bez
zewnętrznego programu i bez dostawcy **[DZIAŁA]**.

### 7.5 Konta i pula rotacji

**Konto** jest wierszem tabeli `konto` o jednym z trzech rodzajów: `cli`, `api`,
`sdk`. Pełny cykl życia kont obsługują komendy `account.add`, `account.list`,
`account.update`, `account.remove`, `account.default.set`, a interfejs ma dla nich
okno (obszar `modele/`) **[DZIAŁA]**.

**Pula rotacji** kanału głównego jest zbudowana z kont rodzaju `cli`. Mechanizm:
kanał bierze konto bieżące, uruchamia proces z jego `CLAUDE_CONFIG_DIR`, a po
rozpoznaniu wyczerpania limitu przechodzi do kolejnego konta puli, zapamiętując
chwilę odnowienia limitu **[DZIAŁA CZĘŚCIOWO]**.

Zastrzeżenia wobec puli:

- **Pula jest budowana jednorazowo przy montażu rdzenia.** Komendy `account.add`,
  `account.remove` i `account.default.set` zapisują wiersze, ale nie odświeżają puli —
  **zmiana katalogu kont wymaga restartu rdzenia** (inaczej niż rejestr kanałów,
  który jest przebudowywany po każdej zmianie).
- **Stan wyczerpania nie jest utrwalany.** Kolumny `konto.stan` i `wyczerpane_do`
  istnieją wraz z metodą zapisu, ale pula ich nie wywołuje **[NIEZINTEGROWANE]**.
- **Rotacja jest niewidoczna dla Operatora.** Kontrakt ma rodzaj fragmentu `account`,
  rdzeń ma strukturę metadanych konta, interfejs ma gotowy moduł jej odczytu — ale
  kanał główny takiego fragmentu nie nadaje **[NIEZINTEGROWANE]**.
- **Pusta pula odmawia wywołania.** Przy zerowej liczbie kont kanał główny kończy
  turę fragmentem błędu o wyczerpaniu limitów, zamiast uruchomić program z profilem
  domyślnym. Komunikat myli brak konfiguracji z wyczerpaniem limitu.

Poświadczenia podlegają regule: wchodzą żądaniem (`account.add`,
`account.update`, `access.point.add`, `access.point.update`) i **nie wychodzą nigdy** —
odpowiedzi niosą wyłącznie informację o obecności poświadczenia albo jego odwołanie.

### 7.6 Nakładka tożsamości

**Decyzje.** Tożsamość modelu składa się z trzech warstw
uporządkowanych wg krytyczności:

1. **konstytucja** — warstwa najwyższa,
2. **profil** (rola),
3. **ekspertyza** zadaniowa.

Tryb podania nakładki ma dwie wartości: **`ZASTAP`** (przełącznik `--system-prompt`,
podmienia prompt fabryczny w całości) oraz **`DOLACZ`** (przełącznik
`--append-system-prompt`, dokłada nakładkę do promptu programu). Wartością domyślną
jest `ZASTAP`. Reguła rozstrzygnięcia trybu całej nakładki: **`ZASTAP` wygrywa
z `DOLACZ`** — jeżeli choć jedna kategoria wnosząca treść żąda zastąpienia, cała
nakładka idzie przełącznikiem `--system-prompt`.

Silnik nakładki nie zna ani jednego zdania promptu; treść pochodzi
wyłącznie z dokumentów tożsamości przechowywanych w bazie i rozstrzyganych per oś.
Suma kontrolna konstytucji służy wyłącznie diagnostyce prowenancji i nigdy nie
dopuszcza ani nie blokuje wywołania.

Warstwa tożsamości jest obsługiwana komendami `identity.category.list`,
`identity.document.get/set/remove`, `identity.effective.get` i ma własne okno
w interfejsie. **Cały mechanizm działa od interfejsu do skutku [DZIAŁA]** — jest to, obok
katalogu ustawień i warstwy dostępów, najlepiej domknięty obszar produktu.

### 7.7 Prowenancja wywołania

**Decyzja.** Pierwszym fragmentem strumienia odpowiedzi jest **prowenancja**
— opis tego, co poszło do modelu: program, tryb nakładki, warstwy nakładki, suma
kontrolna nakładki, nakład, tryb uprawnień, katalogi, katalog roboczy, plik ustawień,
katalog konfiguracji konta. Fragment jest nadawany **przed uruchomieniem procesu**,
więc odbiorca zna warunki wywołania niezależnie od jego powodzenia **[DZIAŁA]**.

Zastrzeżenie **[DO DECYZJI OPERATORA]**: ten sam rodzaj fragmentu jest nadawany
w dwóch niezgodnych kształtach JSON — jednym przez kanał główny, drugim przez warstwę
modeli. Scalenie tych dwóch kształtów pozostaje nierozstrzygnięte i należy je wykonać
przed rozbudową narzędzi diagnostycznych opartych na prowenancji.

Odnotowana wcześniej luka **jest zamknięta pracami porządkowymi nad dokumentacją**
(bieżący stan repozytorium). Do audytu na
stanu sprzed prac naprawczych rodzaje fragmentu `provenance` i `account` rzeczywiście nie
występowały w wyliczeniu rodzajów fragmentu kontraktu i żyły w trzech definicjach
lokalnych. Stan dzisiejszy jest inny: wyliczenie `ChunkKind` w `shared/contract.go`
i `shared/contract.ts` niesie dziewięć wartości — `text`, `thinking`, `tool_use`,
`tool_result`, `image`, `audio`, `error`, `provenance`, `account` — a trzy dawne
definicje lokalne (`server/internal/models/rodzaj_fragmentu.go`,
`server/internal/injection/fragment.go`, `client/src/rozmowa/rodzaje-fragmentow.ts`)
nie powielają już literału: odwołują się wyłącznie do stałych kontraktu
(`shared.ChunkKindProvenance`, `shared.ChunkKindAccount` po stronie rdzenia;
`ChunkKind.Provenance`, `ChunkKind.Account` po stronie klienta).

### 7.8 Czego warstwa kanałów w tej wersji nie robi

Zestawienie zbiorcze, żeby nie było wątpliwości co do zakresu:

| Funkcja | Stan | Uzasadnienie |
|---|---|---|
| Ciągłość rozmowy między turami | **[BRAK]** | przełącznik `--resume` jest zbudowany w warstwie argumentów, ale pole `Wznowienie` zapytania nie ma producenta, więc przełącznik nie trafia do żadnego wywołania; identyfikator rozmowy programu jest odczytywany ze strumienia i porzucany. **Model nie pamięta poprzednich wiadomości okna.** |
| Wybór modelu per okno | **[NIEZINTEGROWANE]** | pole zapytania bez producenta; klucz interfejsu nieznany rdzeniowi. Samo wywołanie **nie zostaje bez modelu** — podstawia się identyfikator z wiersza rejestru kanału (patrz [7.2](#72-kanał-główny--claude-code-cli)) |
| Model zapasowy per okno | **[NIEZINTEGROWANE]** | jak wyżej |
| Nakład rozumowania per okno | **[NIEZINTEGROWANE]** | jak wyżej; dodatkowo interfejs zapisuje liczbę 1–5 przy katalogu przewidującym wyliczenie |
| Zasięg wykonania `local`/`core`/`remote` | **[ATRAPA]** | wartość wyłącznie kopiowana do prowenancji; proces zawsze startuje na maszynie rdzenia |
| Host wykonania — nazwa hosta, jak `danaco-system` | **[NIEZINTEGROWANE]** | ustawienie istnieje, konsumenta wykonawczego brak |
| Kanał wykonawczy SSH | **[BRAK]** | `ssh` występuje wyłącznie jako polecenie mostu MCP |
| Konta rodzaju `api` / `sdk` w wykonaniu | **[NIEZINTEGROWANE]** | droga wykonawcza obejmuje wyłącznie rodzaj `cli` |
| Narzędzia modelu (sterowanie platformą) | **[NIEZINTEGROWANE]** | 39 deklaracji bez konsumenta |

---

## 8. Konfiguracja

### 8.1 Katalog ustawień sterowany danymi

**Decyzja.** Okno konfiguracji nie ma w kodzie ani jednej listy pól.
Kategorie, definicje ustawień, ich rodzaje, wartości domyślne, podpowiedzi, opcje
wyliczeń, dopuszczalne poziomy zasięgu i osie pochodzą z **katalogu w bazie**
(migracja 012) i docierają do interfejsu komendami `settings.category.list`
i `settings.definition.list`. Nowa pozycja konfiguracji to nowy wiersz migracji,
nie nowa gałąź w rdzeniu **[DZIAŁA]**.

Osiem kategorii katalogu:

| Kod | Nazwa | Zakres |
|---|---|---|
| `modele` | Ustawienia modeli | kanał modelu, kanał zapasowy, nakład rozumowania |
| `harness` | Harness | powłoka wykonawcza modelu: program, plik ustawień, mosty MCP, próg biegu naprawczego |
| `tozsamosc` | Tożsamość modelu | tryb podania tożsamości: zastąpienie promptu fabrycznego albo dopisanie |
| `bezpieczenstwo` | Security | zakres zgody wydanej modelowi, egzekwowanie uwierzytelniania |
| `katalog_roboczy` | Katalog roboczy | miejsce powstawania katalogów sesyjnych i plików roboczych modelu |
| `wykonanie` | Wykonanie okna | rola okna w pętli, zasięg i host wykonania |
| `izolacja` | Izolacja zasięgu | jedenaście punktów izolacji |
| `personalizacja` | Personalizacja | wygląd interfejsu Operatora |

Kolumna „Nazwa” podaje **dosłowne brzmienie** zapisane w katalogu ustawień
(`budowa/server/internal/store/migracja_012_katalog_ustawien.sql`), nie tłumaczenie
redakcyjne. Dwie z ośmiu nazw są w tym zaczynie angielskie — „Harness” i „Security” —
i tak też wyświetla je interfejs **[DO DECYZJI OPERATORA]**: czy zaczyn katalogu ma
zostać ujednolicony do polszczyzny pozostałych sześciu kategorii.

Zapis i odczyt obsługują komendy `config.get`, `config.set` i `config.reset`;
zmiana rozgłaszana jest zdarzeniem `config.changed`. **Zapis działa** — wartość
trafia do tabeli `ustawienie` pod adresem złożonym z klucza, poziomu, bytu poziomu,
osi i bytu osi.

### 8.2 Osiem poziomów zasięgu

**Decyzja.** Ustawienie może być zapisane na jednym z ośmiu poziomów;
poziom najwęższy wygrywa z szerszymi:

| Kolejność | Poziom (kontrakt) | Kod w bazie | Znaczenie |
|---|---|---|---|
| 1 (najszerszy) | `global` | `globalny` | cała platforma |
| 2 | `environment` | `srodowisko` | jedno ze środowisk |
| 3 | `module` | `modul` | jeden moduł |
| 4 | `modulePair` | `para_modulow` | para modułów |
| 5 | `project` | `projekt` | projekt |
| 6 | `session` | `karta_sesji` | karta sesji |
| 7 | `role` | `rola` | rola okna |
| 8 (najwęższy) | `window` | `okno` | okno komunikacji |

Rozstrzygacz jest zbudowany poprawnie i pokryty testami: reguła pierwszeństwa
najwęższego poziomu, adresy generowane dwuwymiarowo (poziom × oś), pominięcie
poziomu, którego pole kontekstu jest puste **[DZIAŁA jako mechanizm]**.

**Zastrzeżenie [NIEZINTEGROWANE]:** w rdzeniu kontekst rozstrzygania nigdy nie niesie
więcej niż pole okna (albo osie modelu i konta). W praktyce rozstrzyganie odbywa się
dziś na **dwóch poziomach z ośmiu** — okno i globalny. Ustawienie zapisane na
poziomie środowiska, modułu, pary modułów, projektu, karty sesji albo roli nigdy nie
wygra, ponieważ poziom nie wchodzi do rachunku. Domknięcie wymaga doprowadzenia do
rdzenia odwzorowania okno → karta sesji → projekt → moduł → środowisko; dane te są
w bazie, a łańcuch więzów istnieje.

### 8.3 Osie rozstrzygania

**Decyzja.** Oś jest **prostopadła** do poziomu zasięgu: poziom mówi „jak
wąsko”, oś mówi „dla czego”. Trzy osie: `platform` · `model` · `account`. Klucz
rozstrzygania jest złożony: `klucz + poziom + byt poziomu + oś + byt osi`. Oś wchodzi
jako pola żądania do istniejącej rodziny `config.*`, a nie jako druga rodzina komend.
Brak osi znaczy `platform`.

Mechanizm jest zbudowany i zapis działa **[DZIAŁA]**.

**Zastrzeżenie [DZIAŁA CZĘŚCIOWO]:** okno konfiguracji odczytuje wartości komendą
`config.get` bez wskazania poziomu, licząc na komplet zapisów ze wszystkich poziomów
(cały rachunek dziedziczenia wykonuje po swojej stronie), natomiast rdzeń w tym
przypadku zwraca politykę efektywną liczoną w kontekście zawsze odwzorowującym
identyfikator zakresu na poziom okna. Podgląd łańcucha dziedziczenia nie może przez
to działać. **[DO DECYZJI OPERATORA]**: czy `config.get` bez poziomu ma oddawać
politykę efektywną, czy komplet zapisów.

### 8.4 Które klucze realnie sterują wykonaniem

Rozdział kluczowy dla poprawnego korzystania z produktu. Katalog ustawień jest
kompletny, ale **nie wszystkie zdefiniowane klucze mają dziś konsumenta
wykonawczego**. Poniższe zestawienie rozdziela jedno od drugiego.

**Sterują wykonaniem [DZIAŁA]:**

| Ustawienie / mechanizm | Skutek w wywołaniu modelu |
|---|---|
| Tożsamość: konstytucja → profil → ekspertyza, tryb `ZASTAP`/`DOLACZ` | `--system-prompt` albo `--append-system-prompt` |
| Tryb uprawnień okna | `--permission-mode` |
| Katalogi robocze okna | powtarzalny `--add-dir` |
| Katalog roboczy sesji (klucze `katalog.roboczy.*`) | katalog startowy procesu modelu |
| Nadania dostępu okna (mosty MCP) | powtarzalny `--mcp-config` |
| Kanał modelu okna | wybór wiersza rejestru, a przez to programu i parametrów |

**Zdefiniowane, lecz nieaplikowane [NIEZINTEGROWANE]:**

| Klucz katalogu | Powód |
|---|---|
| `kanal_modelu` (jako ustawienie zasięgu) | rozstrzygacz nie jest źródłem parametrów wywołania |
| `kanal_modelu_zapasowy` | pole zapytania bez producenta |
| `naklad_rozumowania` | pole zapytania bez producenta |
| `srodowisko_wykonania` | wartość wyłącznie w prowenancji |
| `host_wykonania` | brak konsumenta wykonawczego |
| `rola_okna` (jako ustawienie zasięgu) | rola żyje w wierszu okna, nie w rozstrzygaczu |
| `harness.program_claude`, `harness.plik_ustawien`, `harness.konfiguracja_mcp` | cała kategoria bez konsumenta |
| `petla.prog_braku_postepu` | próg biegu naprawczego nieodczytywany |
| jedenaście kluczy `izolacja_*` | brak egzekutora — patrz [8.5](#85-izolacja-zasięgu) |
| egzekwowanie uwierzytelniania, motyw | bez konsumenta |

> **Zapis kluczy nastaw.** Nazwy w postaci `obszar.grupa.nastawa` użyte w tym rozdziale są
> kluczami konfiguracji, nie komendami kontraktu — konwencję zapisu wiąże
> [Standard redakcyjny i językowy](STANDARD-REDAKCYJNY-I-JEZYKOWY.md) rozdz. 7.3.

**Rozjazd słowników [DO DECYZJI OPERATORA].** Panel sterowania okna zapisuje trzy
ustawienia pod własnymi kluczami — `window.executionHost`, `window.fallbackChannelId`,
`window.reasoningEffort` — podczas gdy rdzeń i katalog ustawień znają te same pojęcia
jako `host_wykonania`, `kanal_modelu_zapasowy`, `naklad_rozumowania`. Nie ma między
tymi zbiorami żadnego przekładu, więc **wartości zapisane z panelu sterowania są
martwe**. Rozstrzygnięcia wymaga miejsce, w którym żyje słownik kluczy ustawień:
kontrakt `shared/` czy katalog w bazie.

### 8.5 Izolacja zasięgu

Opracowanie [Izolacja i zależności](architektura/izolacja-i-zaleznosci.md)
opisuje mechanizm rdzenny: **jedenaście punktów
izolacji** — trzy wymiary kontekstu i osiem zakresów technicznych — rozstrzyganych
osobno dla każdego zakresu przez osiem poziomów zasięgu.

| Grupa | Klucze |
|---|---|
| Kontekst | `izolacja_historia` · `izolacja_pamiec` · `izolacja_kontekst` (wartości: `odrebna` / `wspoldzielona`) |
| Zakresy techniczne | `izolacja_katalog_roboczy_sesji` · `izolacja_srodowisko_procesu` · `izolacja_katalog_danych_modelu` · `izolacja_dostep_sieciowy` · `izolacja_odczyt_zapis_plikow` · `izolacja_konto_i_token` · `izolacja_model_procesu` · `izolacja_serwer_wykonania` (wartości: `wlaczony` / `wylaczony`) |

Stan wyjściowy platformy: kontekst odrębny, żaden zakres techniczny niewłączony.
Wszelka dalsza izolacja jest decyzją Operatora, nigdy ustawieniem narzuconym.

**Stan w v2.0: [NIEZINTEGROWANE].** Klucze są zdefiniowane, zasiane do katalogu,
rozstrzygane i wystawiane do odczytu, ale **żaden z nich nie wpływa na uruchomienie
procesu modelu** — nie ma egzekutora. Izolacja per okno i per sesja, kluczowa dla
pracy okien równoległych, w tej wersji nie działa. Kontrakt nie ma dla izolacji ani
jednej komendy, ani jednego zdarzenia, a interfejs nie ma dla niej widoku.

### 8.6 Konfiguracja startowa procesu rdzenia

Konfiguracja startu procesu jest warstwowa: **wartość domyślna → zmienna środowiska
→ argument wywołania** (argument wygrywa ze zmienną, zmienna z wartością domyślną).
Wykaz nazw zmiennych prowadzi jedna funkcja w kodzie, a `budowa/.env.example` jest
z nią zgodny co do joty — deklarowana jest wyłącznie zmienna, którą rdzeń
rzeczywiście czyta.

| Zmienna | Przełącznik | Znaczenie | Wartość domyślna |
|---|---|---|---|
| `DANACO_ROLA` | `--role` | rola procesu: `hub` \| `agent` \| `all` | `all` |
| `DANACO_PORT` | `--port` | port nasłuchu rdzenia | `17870` |
| `DANACO_KATALOG_DANYCH` | `--dane` | katalog danych rdzenia (baza SQLite i pozostała trwałość) | `%LOCALAPPDATA%\DanacoConsole` |
| `DANACO_KATALOG_KLIENTA` | `--klient` | katalog pakietu interfejsu serwowanego obok gniazda | `client\dist` względem katalogu roboczego procesu |
| `DANACO_KATALOG_PROFILI` | `--profile` | katalog profili kanału głównego (konta rotacji) | puste — pula kont startuje pusta |

Zmienne czytane przez **powłokę** (nie przez rdzeń):

| Zmienna | Znaczenie |
|---|---|
| `DANACO_RDZEN` | ścieżka binarki rdzenia, gdy leży poza miejscami znanymi powłoce |
| `DANACO_ADRES_INTERFEJSU` | adres interfejsu (serwer rozwojowy albo rdzeń serwujący pakiet) |
| `DANACO_PORT` | port rdzenia — powłoka odczytuje wartość należącą do rdzenia i przekazuje procesowi potomnemu |

Katalog nieistniejący nie wstrzymuje nasłuchu, wartość nieliczbowa portu nie
przerywa startu, brak katalogu profili nie blokuje pracy — wszystko zgodnie
z regułą fail-open.

**Zastrzeżenie [DZIAŁA CZĘŚCIOWO]:** wartość portu `17870` jest zapisana w trzech
miejscach (rdzeń, powłoka, klient), spięta wyłącznie komentarzem. Zmiana
`DANACO_PORT` bez odpowiedniej zmiany w kliencie rozspaja interfejs z rdzeniem.

Sekrety nigdy nie trafiają do repozytorium: żyją w pliku `budowa\.env` poza kontrolą
wersji, a repozytorium zawiera wyłącznie wzorzec `.env.example`. Baza
przechowuje **odwołania** do danych dostępowych, nigdy ich treści.

---

## 9. Sesje i okna komunikacji

### 9.1 Model bytów

**Decyzja.** Między sesją a wiadomością stoi **okno komunikacji** — byt
pośredni, który niesie parametry pracy modelu:

```
środowisko → karta sesji → sesja → okno komunikacji → wiadomość
```

| Byt | Zawartość |
|---|---|
| **Środowisko** | profil widoczności modułów |
| **Karta sesji** | pojemnik sesji w środowisku |
| **Sesja** | wspólna dla plików, pamięci, projektu i agentów |
| **Okno komunikacji** | moduł · kanał modelu · lista katalogów roboczych · zasięg wykonania · tryb uprawnień · rola w pętli |
| **Wiadomość** | rola, treść, stan, prowenancja |

Identyfikator modułu przeszedł wraz z wprowadzeniem okna komunikacji z sesji
na okno, a relacja procesu
sesji wobec sesji zmieniła się z 1:1 na 1:N.

### 9.2 Cykl życia

Komendy cyklu życia: `session.create`, `session.list`, `session.open`,
`session.close`, `session.delete`, `session.focus`, `session.bind` oraz
`window.create`, `window.list`, `window.update`, `window.close`, `window.state.get`.
Wszystkie mają obsługiwacz w rdzeniu **[DZIAŁA po stronie rdzenia]**.

Zastrzeżenia istotne dla trwałości pracy:

- **`session.create` i `window.create` nie zapisują wiersza do bazy
  [DZIAŁA CZĘŚCIOWO].** Byty powstają wyłącznie w rejestrze pamięciowym nadzorcy;
  wiersz `sesja` i `okno_komunikacji` powstaje dopiero przy pierwszej utrwalonej
  wiadomości. Sesja założona i niewykorzystana nie przeżywa restartu rdzenia.
- **`window.update` nie zapisuje niczego [NIEZINTEGROWANE].** Adapter zmienia byt
  w pamięci nadzorcy; metoda zapisu repozytorium okien — wraz z transakcyjnym
  przepisaniem katalogów okna — nie ma wywołania.
- **Sesje trafiają zawsze do pierwszego środowiska** (`talkin`) i pod stałą nazwę
  „Karta domyślna”; przypisanie sesji do WorkSpace, CodeStudio i MultitaskingAI nie
  jest zapisywane **[NIEZGODNE Z KONCEPCJĄ]**.
- **`workspace.enter` zakłada okno bez kanału modelu [DZIAŁA CZĘŚCIOWO].** Wejście
  do przestrzeni roboczej uzupełnia środowisko wykonania, tryb uprawnień i rolę, ale
  nie uzupełnia kanału modelu — pierwsze `message.send` w takim oknie kończy się
  błędem o oknie bez wskazanego kanału.
- **Historia rozmowy nie jest odtwarzana w interfejsie [BRAK].** Rdzeń utrwala każdą
  wiadomość i odpowiedź oraz oddaje historię okna komendą `message.list`, ale
  interfejs nie wywołuje jej ani razu. Historia okna jest wyłącznie buforem
  interfejsu zapełnianym zdarzeniami bieżącego połączenia — po odświeżeniu widoku
  albo ponownym połączeniu **znika z ekranu** (dane w bazie pozostają).

### 9.3 Okna równoległe i role

Scena sesji mieści **do trzech okien komunikacji** jednocześnie, z przełącznikiem
liczby okien 1/2/3. Każde okno jest zamawiane w rdzeniu komendą `window.create`
w chwili wejścia gniazda na scenę **[DZIAŁA CZĘŚCIOWO]**.

Rola okna w pętli koordynator–wykonawca ma trzy wartości: `executor` (wykonawca),
`coordinator` (koordynator), `standalone` (samodzielne). Okno samodzielne pozostaje
poza pętlą.

Zastrzeżenia:

- **Przekazanie zlecenia koordynator → wykonawca jest w interfejsie wyłącznie
  animacją [NIEZINTEGROWANE].** Wpisy systemowe w historii obu okien, mrugnięcie
  nagłówków, żeton po pasie relacji i przestawienie stanu po 900 ms wykonują się
  lokalnie, **bez jednej komendy kontraktu**. Kontrakt nie ma dziś komendy
  przekazania zlecenia — luka odnotowana.
- **Zmiana roli okna nie wraca na scenę** — figura koordynator/wykonawca może
  rozjechać się ze stanem rdzenia **[DZIAŁA CZĘŚCIOWO]**.
- **Scena mieści trzy okna, a MultitaskingAI wymaga czterech okien ról**
  **[DO DECYZJI OPERATORA]**.

### 9.4 Pętla koordynator–wykonawca

**Decyzje.** W rdzeniu istnieje **jeden silnik pętli**
obsługujący pętlę sesyjną i orkiestrację, z licznikiem obiegów, progiem braku
postępu i jawnym, testowalnym bytem stanu. Obieg startuje turą rozmowy — jedyną
drogą wejścia **[DZIAŁA CZĘŚCIOWO]**.

Telemetria postępu ma producenta: rdzeń emituje zdarzenie
`progress.changed` z identyfikatorem procesu, etapem, liczbą etapów, stopniem
ukończenia i stanem pętli. **Ma też — od prac po rewizji odniesienia — dokładnie
jednego konsumenta [DZIAŁA CZĘŚCIOWO]:** źródło pulpitu Mission Control
(`client/src/mission-control/zrodlo-pulpitu.ts`) subskrybuje je i zasila nim
kolumnę procesów oraz kafle aktywności. Ustalenie audytu o zdarzeniu „idącym po
łączu do nikogo” przestało zatem obowiązywać — z zastrzeżeniem, że konsumentem
jest **jeden widok**: okno komunikacji ani scena okien równoległych nadal nie
pokazują postępu tury. To samo dotyczy zdarzenia `session.changed`, które czytają
dziś dwa miejsca: pulpit oraz wykaz sesji na stronie głównej.

Komenda `window.state.get` zwraca zawsze stan „oczekujący”, ponieważ podsystem
procesów okien jest w montażu rdzenia pozbawiony źródła poleceń i pozostaje
nieosiągalny z punktu wejścia **[NIEZINTEGROWANE]** — luka odnotowana.
**[DO DECYZJI OPERATORA]:** czy podsystem procesów okien zostaje wpięty jako droga
uruchamiania procesu modelu, czy zostaje usunięty na rzecz drogi przez rejestr
kanałów.

### 9.5 Kolejki

Kontrakt ma dla kolejek dwie komendy: `queue.create` i `queue.action` (start,
wznowienie, ponowienie, wstrzymanie, zatrzymanie), oraz zdarzenie `queue.changed`.
Zapis nazwy, stanu i dowiązania sesji działa.

**Stan: [ATRAPA].** Komenda `queue.action` ma od prac po rewizji odniesienia realnego
nadawcę — pasek transportu kolejki na pulpicie Mission Control wysyła ją kanałem
kontraktu, z sześcioma wartościami wyliczenia `QueueAction` i z komunikatem błędu
przy odmowie rdzenia. Nie zmienia to jednak oceny warstwy: **`queue.action`
przestawia wyłącznie kolumnę stanu** — nic nie jest uruchamiane, nic nie jest
kolejkowane, żadne okno nie dostaje pracy. **Pozycje kolejki nigdy nie powstają**:
metody dodania pozycji, zmiany jej stanu, zwiększenia obiegu, wykazu pozycji
i dziennika akcji nie mają w czasie pracy ani jednego wywołania. Silnika wykonania
kolejek nie ma. Blokuje to środowisko MultitaskingAI. Kontrakt nie ma też odczytu
`queue.list`, więc pulpit buduje stan kolejek wyłącznie ze zdarzeń `queue.changed`
i do pierwszego zdarzenia pokazuje stan pusty.

Odnotowana rozbieżność **[DO DECYZJI OPERATORA]**: pulpit dowodzenia przewiduje
podniesienie priorytetu pozycji w kolejce, ale kontrakt nie ma dla tego ani wartości
akcji, ani komendy — luka odnotowana. Interfejs nie udaje tej czynności —
naciśnięcie przycisku zwraca komunikat mówiący wprost, czego w kontrakcie brakuje.
Rozstrzygnięcie: albo akcja kolejki otrzymuje wartość podnoszącą
priorytet, albo funkcja znika z pulpitu.

### 9.6 Anulowanie i odporność

- **`message.stop` działa [DZIAŁA]** — kontekst tury pochodzi z życia rdzenia, nie
  z połączenia, więc żądanie zatrzymania dociera w trakcie trwającej tury i przerywa
  proces modelu.
- **Zerwanie połączenia interfejsu nie kończy pracy [DZIAŁA]** — rdzeń pracuje dalej,
  zgodnie z regułą błędu bieżącego wywołania.
- **Tura bez ani jednego fragmentu i tak zostaje domknięta [DZIAŁA]**.
- Zastrzeżenia: przy błędzie kanału interfejs otrzymuje **trzy** fragmenty błędu
  zamiast jednego; wyścig w rejestrze biegnących tur może unieważnić zatrzymanie
  nowej tury; proces tury kanału głównego startuje bez atrybutów drzewa procesów,
  więc zatrzymanie nie ubija potomstwa procesu modelu **[DZIAŁA CZĘŚCIOWO]**.

---

## 10. Dane

### 10.1 Schemat i migracje

Cała trwałość mieści się w jednym pliku SQLite. Schemat powstaje z **trzynastu
kroków migracji** wkompilowanych w binarium rdzenia i stosowanych automatycznie przy
starcie:

| Krok | Zakres |
|---|---|
| 001 | fundament: środowiska, moduły, karty sesji, sesje, kanały modelu, konta, poziomy zasięgu, ustawienia |
| 002 | okna komunikacji, katalogi okna, procesy sesji |
| 003 | kolejki, pozycje kolejki, dziennik akcji kolejki |
| 004 | pamięć wielopoziomowa: zasoby pamięci, konfiguracja pamięci sesji |
| 005 | katalog okien operacyjnych |
| 006 | identyfikatory zewnętrzne bytów |
| 007 | zaczyn słowników: 4 środowiska, 15 modułów, macierz widoczności, katalog okien operacyjnych |
| 008 | katalog akcji |
| 009 | zaczyn katalogu akcji |
| 012 | katalog ustawień: kategorie, definicje, opcje, zasięgi i osie |
| 013 | punkty dostępu i nadania, argumenty trybu mostu |
| 014 | katalog kont: dostawca, model domyślny, adres bazowy, stan i chwila odnowienia limitu |
| 015 | tożsamość modelu: kategorie i dokumenty tożsamości per oś |

Numeracja ma **lukę na krokach 010 i 011** — miejsca zarezerwowane, bez śladu
w dokumentacji modelu danych. Sama luka nie jest usterką warstwy migracji:
[INSTALACJA-I-KONFIGURACJA.md](INSTALACJA-I-KONFIGURACJA.md) rozdz. 9.2 rozstrzyga,
że numery 010 i 011 nigdy nie istniały, a zbiór wersji 1–9 oraz 12–15 w tabeli
`migracja` jest stanem poprawnym. Usterką pozostaje wyłącznie brak odnotowania tej
luki w dokumentacji modelu danych **[DZIAŁA CZĘŚCIOWO]**.

Właściwości warstwy migracji **[DZIAŁA]**: każdy krok w jednej transakcji, suma
kontrolna wykrywająca zmianę treści kroku już zastosowanego, rejestr wersji,
idempotentne wstawienia zaczynu (`ON CONFLICT DO NOTHING`), test stosujący cały
łańcuch migracji na czystej bazie.

**Pokrycie modelu danych [DZIAŁA CZĘŚCIOWO]:** schemat zakłada **35 tabel** wobec
**68 tabel** docelowego modelu opisanego w [Model danych](architektura/model-danych.md).

**Źródło pomiaru.** Liczba pochodzi ze zliczenia poleceń `CREATE TABLE`
w `budowa/server/internal/store/*.sql`: migracje 001–015 niosą **36 takich poleceń**,
z czego jedno tworzy tabelę przejściową `ustawienie_z_osia` — migracja 012 usuwa
starą tabelę `ustawienie` i przemianowuje na nią tabelę przejściową, więc do stanu
końcowego schematu wchodzi **35 tabel odrębnych**. Docelowy model danych opisany
w [Model danych](architektura/model-danych.md) obejmuje **68 tabel**, z rozbiciem
na grupy tematyczne (Załącznik B tego opracowania) — różnica wobec 35 tabel wdrożonych
wynika z encji przewidzianych koncepcją, którym schemat SQLite nie nadał jeszcze
trwałości. Podana wcześniej liczba 34 była błędem zliczenia.

Trzy całe grupy tematyczne
— tabele izolacji (`profil_izolacji` i pokrewne), załączniki wiadomości oraz
artefakty — nie mają jeszcze trwałości; izolacja żyje wyłącznie jako klucze w tabeli
`ustawienie`.

Trzy tabele pierwszego przekroju pionowego — `urzadzenie`, `polaczenie`,
`proces_sesji` — nie mają repozytorium ani ani jednego zapisu **[BRAK]**. Skutkiem
ubocznym jest niezapisywalność punktów dostępu rodzaju `localDirectory`, które
wymagają dowiązania do urządzenia.

### 10.2 Trwałość rozmowy

Historia rozmowy jest zapisywana do tabeli `wiadomosc` przez łańcuch więzów
`środowisko → karta sesji → sesja → okno → wiadomość`. Baza jest źródłem prawdy;
bufor pamięci obowiązuje wyłącznie jako jawna degradacja per okno, gdy zapis nie
jest możliwy **[DZIAŁA]**.

Zastrzeżenia:

- Kolumny `persona`, `okno_zrodlowe_id`, `tokeny_wejscia`, `tokeny_wyjscia` nigdy nie
  są zapisywane **[NIEZINTEGROWANE]**.
- Kolumna `okno_koordynatora_id` nigdy nie jest zapisywana, więc pętla
  koordynator–wykonawca nie przeżywa restartu rdzenia **[NIEZINTEGROWANE]**.
- Odwołanie do treści obszernej (`tresc_odwolanie`) nie jest wypełniane — druga
  część decyzji o trwałości pozostaje niewykonana **[BRAK]**.

### 10.3 Odtworzenie stanu

Po restarcie rdzeń odtwarza sesje i okna z bazy pod tymi samymi identyfikatorami
**[DZIAŁA]**. Odtworzenie obejmuje byty, które zdążyły zostać utrwalone — czyli te,
w których padła co najmniej jedna wiadomość (patrz [9.2](#92-cykl-życia)).

Kontrola spójności bazy istnieje w kodzie, ale jest osiągalna wyłącznie z testów
**[NIEZINTEGROWANE]**.

### 10.4 Pamięć wielopoziomowa

Repozytorium pamięci (tabele `zasob_pamieci`, `konfiguracja_pamieci_sesji`) jest
zbudowane i złożone w zestawie repozytoriów, ale **kontrakt nie ma dla obszaru
pamięci ani jednej komendy**, a montaż rdzenia nie wpina repozytorium do żadnego
obsługiwacza **[NIEZINTEGROWANE]**.

---

## 11. Rozszerzenia

**Decyzja.** Katalog rozszerzeń należy do zakresu
produktu: przeglądanie, instalacja, zarządzanie i kategorie, oparte na tabeli
`rozszerzenie` i rodzinie komend `extension.*`. Marketplace jako model wydawniczy
(publikowanie i dystrybucja do innych Operatorów) został świadomie odłożony do
osobnego opracowania.

**Stan w v2.0: [BRAK].** Kontrakt wykonawczy nie zawiera ani jednej komendy
`extension.*`, schemat nie zawiera tabeli `rozszerzenie`, a interfejs nie ma widoku
katalogu rozszerzeń. Decyzja jest przyjęta, implementacji nie ma.

---

## 12. Integracje

### 12.1 Punkty dostępu i nadania

**Decyzje.** Dostęp modelu do maszyn i katalogów opisują dwa byty:

- **Punkt dostępu** (`access.point.*`) — mówi, do jakich maszyn i katalogów model ma
  wgląd. Dwa rodzaje: `mcpBridge` (most MCP do maszyny) i `localDirectory` (katalog
  lokalny urządzenia).
- **Nadanie** (`access.grant.*`) — wiąże punkt z **oknem rozmowy**. Okno ma **zbiór**
  nadań, nie jedno; kolejność i oznaczenie nadania głównego niosą znaczenie. Każde
  nadanie ma tryb (`read` / `write`) i podzbiór korzeni punktu.

Warstwa jest domknięta od interfejsu do bazy dla rodzaju `mcpBridge`: pełny cykl
życia punktów i nadań, zdarzenia `access.point.changed` i `access.grant.changed`,
komenda sprawdzenia punktu, okno dostępów w interfejsie **[DZIAŁA]**.

Rodzaj `localDirectory` jest **[BRAK]** — wymaga dowiązania do urządzenia, a tabela
`urzadzenie` nie ma repozytorium ani żadnej ścieżki zapisu.

### 12.2 Mosty MCP po SSH

Most MCP jest mechanizmem, którym model uzyskuje dostęp do maszyn. Rdzeń składa
z nadań okna wpis konfiguracji `mcpServers` i przekazuje go procesowi modelu
przełącznikiem `--mcp-config` **[DZIAŁA]**.

Kształt wpisu odpowiada dosłownie stronie serwerowej mostu:

| Element | Wartość |
|---|---|
| Przedrostek klucza wpisu | `mcp-danaco-pulpit-console-` + identyfikator maszyny |
| Rodzaj połączenia | `stdio` |
| Program zestawiający połączenie | `ssh` |
| Skrypt uruchamiający most po stronie maszyny | `/opt/danaco/most-konsoli/uruchom-stdio.sh` |
| Zmienna przenosząca korzenie dostępu | `DANACO_MOST_KORZENIE` (korzenie rozdzielone dwukropkiem) |
| Domyślny użytkownik | `ubuntu` (gdy adres punktu nie niesie nazwy konta) |
| Domyślny port | `22` |

Tryb uprawnień mostu pochodzi **z nadania**, nie z ustawienia globalnego — nadanie
żyje per okno rozmowy. Słownictwo trybu — „odczyt”, „zapis” — nie jest zaszyte
w kodzie: należy do konkretnego mostu i mieszka w tabeli `argument_trybu_mostu`
(migracja 013).

Cały generator konfiguracji mostu to funkcje czyste — dane wchodzą, tekst wychodzi —
dzięki czemu kształt wpisu daje się sprawdzić bez maszyny po drugiej stronie.

Zastrzeżenia:

- **Mechanizm ma zapis decyzyjny o statusie `Zaproponowana`
  [DO DECYZJI OPERATORA].** Zapis powstał w pracach porządkowych nad dokumentacją
  i odnotowuje jawnie, że most MCP po SSH wszedł do kodu przed decyzją — to samo
  odwrócenie trybu pracy, które dokumentacja odnotowała już wcześniej przy kilku
  innych rozstrzygnięciach; odstępstwo od trybu jest zapisane, nie przemilczane.
  Ratyfikacja wpisu należy do Operatora ([20](#20-kwestie-do-decyzji-operatora),
  pkt 12).
- Ścieżka skryptu, program połączenia, użytkownik i port są **stałymi w kodzie**,
  mimo istnienia odpowiadających im kolumn w bazie **[DZIAŁA CZĘŚCIOWO]**.
  Dokumentacja rekomenduje przy ratyfikacji przeniesienie tych wartości do wiersza
  punktu dostępu, zgodnie z zasadą konfiguracji zamiast logiki wpisanej na
  stałe — rozstrzygnięcie pozostaje otwarte.

### 12.3 Claude Code CLI

Integracja z zewnętrznym programem `claude` (Claude Code CLI) jest opisana
w rozdziale [7.2](#72-kanał-główny--claude-code-cli). Z punktu widzenia wdrożenia
istotne są trzy fakty:

1. Program musi być osiągalny — ze ścieżki zapisanej w wierszu rejestru kanału albo
   z `PATH` pod nazwą `claude`.
2. Wartości trybu uprawnień kontraktu (`manual`, `acceptEdits`, `plan`, `auto`,
   `dontAsk`, `bypassPermissions`) trafiają dosłownie do przełącznika
   `--permission-mode`; zostały zweryfikowane wobec pomocy programu.
3. Uwierzytelnianie odbywa się przez katalog profilu wskazywany zmienną
   `CLAUDE_CONFIG_DIR` — produkt nie przechowuje poświadczeń.

---

## 13. Struktura aplikacji i katalogów

### 13.1 Model czterech katalogów

Projekt rozdziela opracowania, kod, próby techniczne i produkt. Materiały jednego
obszaru nie są mieszane z innym:

| Katalog | Przeznaczenie | Kontrola wersji |
|---|---|---|
| `C:\DanacoConsole_Git\opracowania\` | źródło prawdy dokumentacyjnej: koncepcja, decyzje, ustalenia Operatora. **Może zawierać dane wrażliwe.** | poza kontrolą wersji |
| `C:\DanacoConsole_Git\budowa\` | repozytorium kodu i procesu deweloperskiego | objęte kontrolą wersji |
| `C:\DanacoConsole_Git\spike\` | próby techniczne rozbrajające ryzyka; wyrzucane po zamknięciu etapu | poza kontrolą wersji |
| `C:\DanacoConsole_App\` | **wyjście budowania** — gotowy produkt do uruchomienia; nigdy nie zawiera źródeł | poza kontrolą wersji |

**Reguła katalogu produktu wraz z zastrzeżeniem stanu v2.0.** Zasadą jest, że
zawartość wykonywalna tego katalogu powstaje **wyłącznie z budowania**, a nie
z ręcznego kopiowania: wytwarza ją skrypt `budowa/scripts/wydanie.sh`, który jako
krok ostatni wywołuje `budowa/scripts/pakowanie.sh` (rozdział
[17.3](#173-wydanie-produktu-do-cdanacoconsole_app)). Reguła ma w bieżącej wersji
**trzy jawne wyjątki**, więc nie należy jej czytać jako bezwzględnego zakazu
dotykania katalogu:

1. **Dokumentacja produktu (`*.md`) trafia tu poza budowaniem.** Skrypt pakowania
   ma to zapisane jako regułę nadrzędną: nigdy nie usuwa ani nie nadpisuje plików
   `*.md`, a czyszczenie ogranicza do podkatalogu `client/` również z wyłączeniem
   `*.md`. Niniejszy dokument, [LICENSE.md](LICENSE.md), [INSTRUKCJA-UZYTKOWANIA.md](INSTRUKCJA-UZYTKOWANIA.md)
   i [INSTALACJA-I-KONFIGURACJA.md](INSTALACJA-I-KONFIGURACJA.md) są w tym katalogu umieszczane niezależnie od
   wydania.
2. **Wariant bez powłoki.** Gdy w środowisku nie ma `cargo` albo `tauri-cli`, albo
   gdy Operator poda `DANACO_BEZ_POWLOKI=1`, wydanie składa wariant rdzeń+klient,
   a pliku `Danaco Console.exe` w katalogu nie ma.
3. **Ręczne wypełnienie jako droga zastępcza.** Dokument
   [INSTRUKCJA-UZYTKOWANIA.md](INSTRUKCJA-UZYTKOWANIA.md) (rozdz. 2.5) opisuje wypełnienie tego katalogu
   ręcznie — jako miejsca, w którym powłoka odnajdzie plik wykonywalny rdzenia.
   Ta droga pozostaje **dopuszczalna** i bywa jedyną możliwą, gdy łańcuch narzędzi
   budowania nie jest w środowisku dostępny. Reguła niniejszego rozdziału jej nie
   zakazuje.

**[DO DECYZJI OPERATORA]:** rozdz. 2.5 [INSTRUKCJA-UZYTKOWANIA.md](INSTRUKCJA-UZYTKOWANIA.md) oraz rozdz. 2.2
tego samego dokumentu stwierdzają dziś, że katalogu produktu nie wytwarza żaden
skrypt repozytorium. Twierdzenie to było prawdziwe na rewizji odniesienia
stanu sprzed prac naprawczych, a przestało być nim po dopisaniu `wydanie.sh`
i `pakowanie.sh`. Wymaga
to jednego przejścia po dokumentacji produktu, żeby wszystkie cztery dokumenty
opisywały ten sam stan; niniejszy dokument opisuje stan bieżący.

### 13.2 Drzewo repozytorium budowy

```
C:\DanacoConsole_Git\budowa\
├── .env.example                 wzorzec konfiguracji środowiska rdzenia
├── go.mod, go.sum               moduł Go `danacoconsole` (obejmuje server/ i shared/)
├── shared\                      KONTRAKT — jedyne źródło prawdy nazw
│   ├── contract.json            źródło
│   ├── contract.go              generat dla rdzenia
│   ├── contract.ts              generat dla interfejsu
│   └── gen\                     generator i sprawdzian zgodności z dokumentem
├── server\
│   ├── cmd\danaco-console\      punkt wejścia (kompozycja) + uruchomienie\
│   └── internal\                transport · protocol · core · session · models
│                                injection · konfig · konfiguracja · dane · store
├── client\                      interfejs TypeScript (Vite, vitest)
│   ├── package.json             skrypty: dev · build · preview · typy · testy
│   └── src\                     main.ts + obszary widoków i warstw
├── desktop\src-tauri\           powłoka natywna (Rust, Tauri 2)
│   ├── tauri.conf.json          konfiguracja pakietu (bundle NSIS)
│   └── src\                     main.rs + moduły powłoki + rdzen\
├── (dokumentacja produktu)      nie leży w budowa\ — mieszka w korzeniu repozytorium,
│                                w katalogu docs\: architektura\ (kontrakty komunikacji,
│                                model danych, izolacja i zależności) · moduly\ ·
│                                interfejs-uzytkownika\ · srodowiska\ · specyfikacje\
└── scripts\                     piętnaście skryptów bramy, wydania i podglądu
    ├── brama.sh                 przebieg zbiorczy bramy (tryb szybki i pełny)
    ├── metryka-plikow.sh        progi modularności
    ├── atrapy.sh                kod bez odbiorcy i znaczniki niedokończenia
    ├── atrapy-eksporty.mjs      eksporty TypeScript bez konsumenta
    ├── eksporty-go.mjs          eksporty Go bez konsumenta
    ├── eksporty-go-zwolnienia.mjs  jawny wykaz zwolnień kontroli eksportów Go
    ├── osiagalnosc.mjs          osiągalność drzewa interfejsu od main.ts
    ├── osiagalnosc-odwolania.mjs   zerwane odwołania w drzewie interfejsu
    ├── wywolania-kontraktu.mjs  komendy kontraktu z realnym wywołaniem u klienta
    ├── dym-pionu.mjs            dymny sprawdzian całego pionu na kanale echo
    ├── kontrola-rust.sh         warstwa powłoki: format, clippy, budowa
    ├── wydanie.sh               wydanie produktu do C:\DanacoConsole_App
    ├── pakowanie.sh             złożenie katalogu produktu z artefaktów
    ├── pokaz.sh                 podgląd deweloperski stanowisk
    ├── kopia.sh                 kopia robocza drzewa
    └── haki\pre-commit          brama w trybie szybkim przed zatwierdzeniem
```

Poza drzewem `budowa\` leży jeszcze `.github\workflows\brama.yml` — opis przepływu
ciągłego uruchamiającego bramę pełną (stan: przygotowany, bezczynny — patrz
[3.9](#39-reguły-modularności-i-brama-weryfikacyjna)).

Odstępstwo od mapy dokumentu nadrzędnego: zapowiedziany katalog `budowa\tests\`
nie istnieje — sprawdziany Go żyją obok kodu, a sprawdziany interfejsu w drzewie
`client/src` **[BRAK]**.

### 13.3 Katalog danych czasu pracy

```
%LOCALAPPDATA%\DanacoConsole\
└── danaco-console.db            jedyny plik trwałości (SQLite, tryb WAL)
```

Zmiana katalogu danych (`DANACO_KATALOG_DANYCH` albo `--dane`) przenosi całą trwałość,
bez osobnego przełącznika. Katalogi profili kanału głównego (`CLAUDE_CONFIG_DIR`)
leżą poza tym drzewem i są wskazywane zmienną `DANACO_KATALOG_PROFILI` albo wpisami
kont.

---

## 14. Warstwa wizualna

**Decyzja.** Obowiązującym kierunkiem wizualnym jest pakiet **design v2.0** —
„monochromatyczna precyzja”:

- odcienie bieli w motywie jasnym, odcienie czerni w ciemnym,
- **jeden błękit sygnałowy** `#3B6FE0`,
- zakaz czystej czerni jako tła i tekstu (skala kończy się na `#0A0A0A`, tekst
  `#181818`),
- działanie główne to **inwersja atramentu**, nie sygnał; sygnał wskazuje, nie działa,
- kroje: **Space Grotesk** · **IBM Plex Sans** · **IBM Plex Mono** (licencja OFL),
- gęstość **zwarta** domyślnie,
- pasek górny **zawsze atramentowy w obu motywach** — rama kokpitu,
- element sygnaturowy: **kropka sygnału**.

Wcześniejsza warstwa zastępcza (granat + złoto, tarcza z wagą, Cormorant Garamond,
Inter, JetBrains Mono) **nie obowiązuje**.

Stan wdrożenia **[DZIAŁA CZĘŚCIOWO]**:

- Żetony motywu zostały przeniesione z pakietu design **dosłownie**; oba motywy są
  równoprawne; **zestaw 82 ikon odpowiada manifestowi** co do pozycji i co do
  zadeklarowanej liczby (`liczba-ikon: 82`) **[DZIAŁA]**. Katalog `ikony/svg/`
  zawiera 83 pliki — 83. z nich, `logo-danaco.svg`, jest sygnetem marki
  w odmianie pełnokolorowej, a nie ikoną obrysową zestawu; nie podlega zasadom
  manifestu i jest osobno wskazany w kodzie (por. [3.4](#34-klient--interfejs-typescript)).
- Biblioteka komponentów została wymieniona na wersję v2.0; po wymianie **warstwa
  widoków fali 3 nadal adresowała część klas ze słownika sprzed wymiany** —
  czternaście nazw klas nie miało pokrycia w arkuszach, ze stosem dymków powiadomień
  bez reguł stylu jako skutkiem najdotkliwszym, bo to jedyny nośnik odpowiedzi na
  kliknięcie w reżimie zera blokad. prace nad arkuszami stylów domknął tę
  lukę arkuszem `komponenty/aliasy-zgodnosci.css`, który pokrywa regułami
  wartości słownika v2.0 wszystkie klasy, po które sięgają widoki — stos dymków
  ma dziś pełne style w obu motywach (klasy `dn-toast`, `dn-toast--blad`,
  `dn-toast--informacja`, `dn-toast--sukces`, `dn-toast-tresc`, `dn-toast-tytul`,
  `dn-toasty` zdefiniowane w `komponenty/powiadomienie.css` i
  `komponenty/aliasy-zgodnosci.css`), co sprawdza wprost
  `komponenty/stos-dymkow.test.ts` (dziewięć sprawdzianów, każda klasa wariantu
  użyta w DOM ma regułę w arkuszu) **[DZIAŁA]**. Jeden przypadek zostaje poza tym
  zamknięciem: klasa `.dn-karta-sesji`, nadawana w
  `powloka/karty-sesji.ts` pasowi kart sesji, nie ma dziś żadnej reguły CSS
  w drzewie klienta — nawet komentarz we własnym arkuszu `powloka/karty-sesji.css`
  przywołuje tę klasę z nazwy, nie definiując jej; dług nieopisany dotąd w żadnym
  rejestrze **[BRAK]**.
- **Stany interfejsu wymagane opracowaniem [Elementy okien](interfejs-uzytkownika/elementy-okien.md)
  w większości nie mają w widokach żadnego
  wystąpienia [DZIAŁA CZĘŚCIOWO]** — biblioteka dostarcza wskaźnik ładowania,
  komunikat blokowy, dymek objaśnienia `[?]` i ekran „łączenie z serwerem…”.
  prace nad arkuszami stylów podpiął pierwsze użycia w widoku sterowania: dymki
  `[?]` objaśnień (`widok-sterowania/dymek-objasnienia.ts`) i wskaźnik ładowania
  przy pozycjach podsumowania ustawień w trakcie odczytu wykazu
  (`widok-sterowania/podsumowanie-ustawien.ts`, sprawdzone testem
  `podsumowanie-ustawien.test.ts`). Poza widokiem sterowania wymóg
  z [Elementy okien](interfejs-uzytkownika/elementy-okien.md), by każde okno miało stan pusty, ładowania i błędu, nadal
  nie jest spełniony.
- **Pasek górny zawiera atrapy**: dzwonek powiadomień, awatar, przełącznik trybu
  zawsze widocznego i pole wyszukiwania nie mają realizacji **[ATRAPA]**.

Zgodnie z przyjętą zasadą makieta określa wygląd, a dokumentacja — zakres funkcji;
odwzorowanie makiety nie jest wykonaniem funkcji.

---

## 15. Scenariusze zastosowania

Rozdział rozdziela scenariusze **wykonalne w wersji v2.0** od scenariuszy
**przewidzianych koncepcją, dziś niewykonalnych**. Rozdział nie opisuje zamiarów jako
możliwości.

### 15.1 Scenariusze wykonalne w v2.0

**Warunek wstępny wspólny dla wszystkich sześciu scenariuszy.** Każdy z nich zakłada
wykonanie co najmniej jednej tury modelu, a **scenariusze zakładają instalację
doprowadzoną do pracy wg [17.2](#172-doprowadzenie-świeżej-instalacji-do-pracy--czynność-obowiązkowa)**
— to jest po ręcznym zasianiu wiersza kanału głównego w tabeli `kanal_modelu`, po
założeniu co najmniej jednego konta rodzaju `cli` i po wskazaniu oknu komunikacji
identyfikatora wiersza rejestru zamiast rodzaju kanału. Na świeżej instalacji żaden
z poniższych scenariuszy nie dojdzie do skutku, ponieważ pierwsze `message.send`
zakończy się odmową. Scenariusz sprawdzający sam tor komunikacyjny, bez zewnętrznego
programu i bez konta, prowadzi kanał `echo` (por. [7.4](#74-kanał-echo)
i skrypt `scripts/dym-pionu.mjs`).

**Scenariusz pierwszy. Ustanowienie tożsamości modelu i sprawdzenie, co realnie poszło do modelu.**
Operator otwiera okno tożsamości, zapisuje dokument konstytucji na poziomie
globalnym, profil roli na poziomie okna, a ekspertyzę zadaniową na osi modelu.
Po wysłaniu wiadomości pierwszy fragment strumienia niesie prowenancję: tryb
nakładki, jej warstwy, sumę kontrolną, tryb uprawnień, katalogi i katalog roboczy.
Scenariusz sprawdza najlepiej domkniętą część produktu i jest właściwym testem
poprawności instalacji.

**Scenariusz drugi. Nadanie modelowi dostępu do maszyny przez most MCP.** Operator zakłada punkt
dostępu rodzaju `mcpBridge` (adres maszyny, korzenie), nadaje go oknu rozmowy
w trybie odczytu albo zapisu i porządkuje kolejność nadań. Rdzeń składa konfigurację
`mcpServers` i przekazuje ją procesowi modelu. Praca modelu na maszynie odbywa się
w granicach korzeni nadania.

**Scenariusz trzeci. Praca na katalogach lokalnych.** Operator wskazuje katalogi robocze okna
(natywnym oknem wyboru powłoki) i katalog roboczy sesji. Katalogi trafiają do
wywołania jako powtarzalny `--add-dir`, a katalog roboczy sesji jako katalog startowy
procesu modelu.

**Scenariusz czwarty. Praca w oknach równoległych.** Operator otwiera dwa albo trzy okna
komunikacji obok siebie, nadaje im role i prowadzi w nich niezależne tury. Każde
okno **jest w rdzeniu odrębnym bytem** i niesie własny kanał, własne katalogi
robocze oraz własny tryb uprawnień; jedno okno prowadzi jedną turę naraz, a kolejne
okna sesji biegną równolegle. Zakres tej odrębności wymaga jednak trzech zastrzeżeń,
bez których zdanie powyżej byłoby nieprawdziwe:

- **Parametry okna żyją w pamięci nadzorcy, nie w bazie.** Komenda `window.update`
  zmienia byt w rejestrze pamięciowym, ale nie zapisuje wiersza
  (por. [9.2](#92-cykl-życia)), więc zmiana kanału, katalogów albo trybu uprawnień
  okna **nie przeżywa restartu rdzenia**.
- **Wskazanie kanału z interfejsu nie trafia w wiersz rejestru.** Interfejs podaje
  w `window.create` **rodzaj** kanału (`cli`) w miejscu przeznaczonym na identyfikator
  wiersza (por. [17.2](#172-doprowadzenie-świeżej-instalacji-do-pracy--czynność-obowiązkowa)),
  więc „własny kanał okna” jest w praktyce kanałem wskazanym ręcznie przy
  doprowadzaniu instalacji do pracy, a nie wybranym w oknie.
- **Odrębność nie jest izolacją.** Przekazanie zlecenia między oknami jest dziś
  animacją interfejsu, a nie komendą kontraktu; jedenaście punktów izolacji nie ma
  egzekutora (por. [8.5](#85-izolacja-zasięgu)), więc okna dzielą środowisko procesu,
  katalog danych modelu i dostęp sieciowy maszyny rdzenia.

**Scenariusz piąty. Zatrzymanie pracy i praca w tle.** Operator przerywa trwającą turę
przyciskiem zatrzymania; przy zerwaniu połączenia interfejsu rdzeń kontynuuje pracę
i utrwala wynik w bazie. Zastrzeżenie: po ponownym połączeniu historia nie wraca na
ekran, ponieważ interfejs nie odczytuje historii z rdzenia.

**Scenariusz szósty. Konfiguracja warstwowa z podglądem prowenancji ustawienia.** Operator zapisuje
ustawienie na wybranym poziomie i osi, a okno konfiguracji pokazuje wartość
obowiązującą. Zastrzeżenie: realny wpływ na wywołanie mają wyłącznie klucze
wymienione w [8.4](#84-które-klucze-realnie-sterują-wykonaniem).

### 15.2 Scenariusze przewidziane koncepcją, w v2.0 niewykonalne

| Scenariusz | Przeszkoda |
|---|---|
| Wielotygodniowa rozmowa z zachowaniem kontekstu | brak `--resume`: każda tura startuje bez historii |
| Praca w module (Studio, Research, Terminal, Developer…) | brak okien operacyjnych; nawigacja modułowa nie przeładowuje przestrzeni roboczej |
| Orkiestracja zespołu modeli w MultitaskingAI | brak panelu orkiestracji, ról, kolejek ról i silnika wykonania kolejek |
| Sterowanie platformą poleceniem w oknie rozmowy | narzędzia modelu wygenerowane, bez konsumenta |
| Wykonanie pracy na hoście zdalnym | zasięg wykonania jest etykietą; kanał SSH nie istnieje |
| Odizolowanie okna od reszty platformy | jedenaście punktów izolacji bez egzekutora |
| Instalacja rozszerzenia z katalogu | brak komend `extension.*` i tabeli `rozszerzenie` |
| Praca z modelem innego dostawcy przez konto | konta `api`/`sdk` bez drogi wykonawczej |
| Nadzór postępu kolejek ról na pulpicie | pulpit czyta stan kolejek wyłącznie ze zdarzeń — kontrakt nie ma odczytu `queue.list`, a `queue.action` nie uruchamia pracy |

---

## 16. Wymagania

### 16.1 Środowisko uruchomieniowe

| Element | Wymaganie |
|---|---|
| System | Windows 11 (środowisko referencyjne budowy i uruchomienia) |
| Rdzeń | pojedyncza binarka `danaco-console.exe`; brak zależności zewnętrznych poza systemem operacyjnym |
| Baza | brak instalacji serwera — SQLite wbudowany w rdzeń (sterownik napisany w całości w języku Go) |
| Port | wolny port TCP **17870** (albo inny wskazany `DANACO_PORT`, z równoległą zmianą po stronie interfejsu) |
| Kanał główny | program `claude` (Claude Code CLI) osiągalny w `PATH` albo ze ścieżki w wierszu rejestru kanału |
| Konta kanału głównego | katalogi profili wskazywane przez `CLAUDE_CONFIG_DIR` |
| Mosty MCP | program `ssh` po stronie rdzenia oraz skrypt `/opt/danaco/most-konsoli/uruchom-stdio.sh` po stronie maszyny |

### 16.2 Środowisko budowy

| Warstwa | Narzędzia |
|---|---|
| Rdzeń | Go 1.26 (`go build ./server/cmd/danaco-console`, `go test ./...`, `go vet ./...`, `gofmt`) — polecenia wydawane z katalogu `budowa`, w którym leży `go.mod` |
| Interfejs | Node.js z `npm`; TypeScript 5.8, Vite 6.3, Vitest 3.2 (`npm run build`, `npm run typy`, `npm run testy`) |
| Powłoka | łańcuch narzędzi Rust wraz z `cargo` i `cargo tauri` (Tauri 2, pakiet NSIS) |
| Kontrakt | Node do uruchomienia generatora `shared/gen/generate.mjs` |
| Brama | powłoka zgodna z POSIX do uruchomienia skryptów `budowa/scripts/*.sh` |

Instalacja narzędzi i bibliotek, które aplikacja ma mieć wbudowane, jest objęta
odrębnym upoważnieniem; nie obejmuje ono zmian w środowiskach produkcyjnych ani
operacji nieodwracalnych.

---

## 17. Podstawowe uruchomienie

Rozdział podaje wyłącznie **zarys**. Pełna procedura wraz z doprowadzeniem świeżej
instalacji do stanu zdolnego do pracy znajduje się w **[INSTALACJA-I-KONFIGURACJA.md](INSTALACJA-I-KONFIGURACJA.md)**;
opis pracy z gotową instalacją — w **[INSTRUKCJA-UZYTKOWANIA.md](INSTRUKCJA-UZYTKOWANIA.md)**.

### 17.1 Zarys uruchomienia w drzewie deweloperskim

1. **Rdzeń.** Zbudować binarkę **z katalogu `budowa`** (tam leży `go.mod` modułu
   `danacoconsole`, obejmującego `server/` i `shared/`):
   `go build ./server/cmd/danaco-console`, a następnie uruchomić ją — bez
   przełączników rdzeń przyjmuje rolę `all`, port `17870` i katalog danych
   `%LOCALAPPDATA%\DanacoConsole`. Migracje stosują się same przy starcie.
   Ścieżka punktu wejścia jest **względna wobec katalogu `budowa`**, a nie wobec
   `budowa/server` — odpowiada drzewu z [13.2](#132-drzewo-repozytorium-budowy),
   krokowi 3 skryptu wydania oraz rozdz. 2.2 [INSTRUKCJA-UZYTKOWANIA.md](INSTRUKCJA-UZYTKOWANIA.md).
2. **Interfejs.** Z katalogu `budowa/client`: `npm install`, a następnie
   `npm run dev` (serwer rozwojowy) albo `npm run build` (pakiet `client/dist`
   serwowany przez rdzeń).
3. **Powłoka.** Z katalogu `budowa/desktop/src-tauri`: `cargo tauri dev`. Powłoka
   sama odnajduje i uruchamia rdzeń w tle oraz otwiera okno interfejsu.

### 17.2 Doprowadzenie świeżej instalacji do pracy — czynność obowiązkowa

**Świeża instalacja nie wykona ani jednej tury modelu bez ręcznego zasiania danych.**
Wynika to z trzech niezależnych przyczyn opisanych w [1.3](#13-stan-wersji-v20--obraz-uczciwy)
i [7](#7-modele-i-kanały-wykonawcze). Do pracy potrzebne są:

1. **Wiersz kanału głównego** w tabeli `kanal_modelu` (rodzaj `cli`) — zakładany
   komendą kontraktu `channel.add` wysłaną po WebSocket; **interfejs nie ma dla tego
   okna**.
2. **Co najmniej jedno konto rodzaju `cli`** w tabeli `konto`, wskazujące katalog
   profilu (`CLAUDE_CONFIG_DIR`) — zakładane komendą `account.add`; po zmianie kont
   **konieczny jest restart rdzenia**, ponieważ pula rotacji jest budowana raz przy
   montażu.
3. **Wskazanie oknu komunikacji identyfikatora wiersza kanału**, a nie rodzaju
   kanału. Interfejs w bieżącej wersji podaje w `window.create` wartość `cli`
   (rodzaj), co nie odpowiada żadnemu wierszowi rejestru.

W drzewie źródeł istnieje moduł zapewnienia kanału głównego, który wykonuje dokładnie
tę czynność (sprawdza `channel.list`, a przy braku kanału zakłada go komendą
`channel.add`), lecz jest podłączony wyłącznie do stanowiska podglądu i nie wchodzi
do pakietu produkcyjnego **[NIEZINTEGROWANE]**.

### 17.3 Wydanie produktu do `C:\DanacoConsole_App`

**Stan: [DZIAŁA].** Wydanie produktu jest **jednym poleceniem**, wydawanym z korzenia
repozytorium:

```
bash budowa/scripts/wydanie.sh [katalog-wyjścia]
```

Argument jest opcjonalny; jego brak daje `C:/DanacoConsole_App`. Skrypt prowadzi
pięć kroków, po czym oddaje złożenie katalogu skryptowi `pakowanie.sh`:

| Krok | Czynność | Uwaga |
|---|---|---|
| 1 | zależności klienta — `npm ci --no-audit --no-fund` | wykonywane wyłącznie przy braku `node_modules`; wymuszenie = usunięcie katalogu |
| 2 | pakiet interfejsu — `npm run build` w `budowa/client` (`tsc --noEmit && vite build` → `client/dist`) | błąd typów przerywa wydanie |
| 3 | rdzeń — `go build ./server/cmd/danaco-console` z katalogu `budowa` | wynik odkładany do katalogu roboczego wydania |
| 4 | powłoka natywna — `cargo tauri build --no-bundle` w `desktop/src-tauri` | krok **pomijalny**, patrz niżej |
| 5 | pakowanie — `pakowanie.sh` przenosi artefakty do katalogu wyjściowego | niczego nie buduje |

**Krok czwarty jest pomijany trzema drogami**, zgodnie z regułą fail-open:
na żądanie zmienną `DANACO_BEZ_POWLOKI=1`, przy braku `cargo` oraz przy braku
`tauri-cli`. Wynikiem jest wtedy **wariant rdzeń+klient**, uruchamiany plikiem
`danaco-console.exe` z interfejsem pod adresem `http://127.0.0.1:17870/`. Wydanie
z powłoką uruchamia się plikiem `Danaco Console.exe`, który stawia rdzeń w tle.
Przełącznik `--no-bundle` daje sam plik wykonywalny; `beforeBuildCommand`
z `tauri.conf.json` przebudowuje przy okazji pakiet klienta — powtórka kroku 2 jest
zamierzona i tania.

**Struktura katalogu wyjściowego** odpowiada dokładnie kolejności poszukiwań
zaszytej w powłoce (`rdzen/lokalizacja.rs` szuka rdzenia obok pliku powłoki,
`rdzen/pakiet_klienta.rs` szuka `client/dist` obok binarki rdzenia):

```
C:\DanacoConsole_App\
├── danaco-console.exe       rdzeń (rola all, nasłuch 17870); serwuje interfejs
├── Danaco Console.exe       powłoka natywna — obecna, gdy krok 4 się powiódł
├── client\dist\             pakiet interfejsu odnajdywany obok rdzenia
└── *.md                     dokumentacja produktu — NIETYKANA przez wydanie
```

Reguła nadrzędna skryptu pakowania: **pliki `*.md` nigdy nie są usuwane ani
nadpisywane**. Czyszczenie obejmuje wyłącznie podkatalog `client/`, i również tam
z wyłączeniem `*.md`. Dzięki temu dokumentacja produktu i pakiet interfejsu żyją
w jednym katalogu bez konfliktu, a kolejne wydania nie kasują dokumentów.
Skrypt pakowania odmawia pracy, gdy nie widzi binarki rdzenia albo pliku
`index.html` w pakiecie klienta — brak artefaktu jest błędem wydania, nie cichym
pominięciem.

**Krok odłożony [BRAK]: instalator NSIS.** Sekcja `bundle` w `tauri.conf.json`
deklaruje cel `nsis`, ale złożenie instalatora wymaga narzędzia NSIS i jest osobnym
krokiem (`cargo tauri build` bez `--no-bundle`, wykonywane w `desktop/src-tauri`).
W bieżącej wersji produkt rozmieszcza się kopią katalogu, a nie instalacją.
**[DO DECYZJI OPERATORA]:** czy instalator NSIS wchodzi do zakresu wydania v2.0,
oraz czy rdzeń ma w nim być osadzony jako zasób powłoki, czy — jak dziś — leżeć
obok niej jako osobny plik wykonywalny.

---

## 18. Znane ograniczenia wersji v2.0

Zestawienie zbiorcze, uporządkowane wg wagi. Każda pozycja jest rozwinięta we
wskazanym rozdziale.

**Blokujące podstawowy przepływ pracy:**

1. Świeża instalacja nie wykona tury modelu — brak zaczynu kanału, pusta pula kont
   odmawia, interfejs podaje rodzaj kanału zamiast identyfikatora wiersza ([17.2](#172-doprowadzenie-świeżej-instalacji-do-pracy--czynność-obowiązkowa)).
2. Brak ciągłości rozmowy — model nie pamięta poprzednich tur okna ([7.8](#78-czego-warstwa-kanałów-w-tej-wersji-nie-robi)).
3. 64 z 79 wierszy katalogu okien operacyjnych w bazie nie mają odpowiadającego
   widoku ([INSTRUKCJA-UZYTKOWANIA.md](INSTRUKCJA-UZYTKOWANIA.md) rozdz. 7);
   pozostałe piętnaście wierszy to jeden wiersz Chat Window na każdy z piętnastu
   modułów, którym odpowiada **jeden** widok produktu — stąd zdanie o jednym widoku
   w rozdziale [1.3](#13-stan-wersji-v20--obraz-uczciwy); wybór modułu nie
   przeładowuje przestrzeni roboczej ([5](#5-moduły)).

**Ograniczające zakres funkcjonalny:**

4. Ruch otwierający interfejsu niezgodny z przyjętą regułą startu; sześć komend nawigacji
   platformy (`home.enter`, `environment.list`, `environment.enter`, `module.list`,
   `workspace.enter`, `window.state.get`) bez konsumenta ([6](#6-środowiska)).
5. Historia rozmowy nie jest odtwarzana w interfejsie ([9.2](#92-cykl-życia)).
6. Model zapasowy i nakład rozumowania nie docierają do wywołania, a wskazanie
   modelu per okno nie ma producenta — model wywołania pochodzi z wiersza rejestru
   kanału; do tego rozjazd słowników kluczy ustawień
   ([7.2](#72-kanał-główny--claude-code-cli), [8.4](#84-które-klucze-realnie-sterują-wykonaniem)).
7. Sześć z ośmiu poziomów zasięgu jest nieosiągalnych ([8.2](#82-osiem-poziomów-zasięgu)).
8. Izolacja bez egzekutora ([8.5](#85-izolacja-zasięgu)).
9. Kolejki bez silnika wykonania; MultitaskingAI bez panelu orkiestracji
   ([9.5](#95-kolejki), [6](#6-środowiska)).
10. Narzędzia modelu bez konsumenta — sterowanie platformą przez model nie działa
    ([3.5](#35-kontrakt-shared)).
11. Pamięć wielopoziomowa bez komend kontraktu ([10.4](#104-pamięć-wielopoziomowa)).
12. Rozszerzenia bez implementacji ([11](#11-rozszerzenia)).
13. Pulpit Mission Control czerpie dane z rdzenia, ale część jego sekcji nie ma
    w kontrakcie źródła ani czynności: brak odczytu `queue.list`, brak wartości
    akcji podnoszącej priorytet, `context.transfer` bez kompletu kontekstu po
    stronie pulpitu ([6](#6-środowiska)).

**Ograniczające jakość pracy:**

14. `window.state.get` zawsze zwraca stan „oczekujący” — podsystem procesów okien
    pozostaje w montażu bez źródła poleceń ([9.4](#94-pętla-koordynatorwykonawca)).
15. Stany interfejsu (ładowanie, komunikat, dymek, ekran łączenia) bez wystąpień
    poza widokiem sterowania, gdzie prace nad arkuszami stylów podpiął pierwsze
    użycia; poza tym widokiem wymóg stanu pustego/ładowania/błędu na okno nadal
    niespełniony. Pokrycie klas CSS domknięte arkuszem
    `komponenty/aliasy-zgodnosci.css` z jednym wyjątkiem resztkowym —
    `.dn-karta-sesji` bez reguły ([14](#14-warstwa-wizualna)).
16. Zmiana kont wymaga restartu rdzenia; stan wyczerpania kont nieutrwalany
    ([7.5](#75-konta-i-pula-rotacji)).
17. `window.update` nie zapisuje do bazy; sesje trafiają zawsze do pierwszego
    środowiska ([9.2](#92-cykl-życia)).
18. Adres gniazda zgadywany z nazwy hosta — pakiet osadzony w powłoce nie łączy się
    z rdzeniem ([3.3](#33-powłoka-natywna--tauri)).
19. Druga część decyzji o trwałości (treści obszerne jako pliki) niezaimplementowana
    ([10.2](#102-trwałość-rozmowy)).
20. Instalator NSIS nie jest składany — produkt rozmieszcza się kopią katalogu
    wydania, nie instalacją ([17.3](#173-wydanie-produktu-do-cdanacoconsole_app)).

---

## 19. Pozostała dokumentacja

Niniejszy dokument jest kartą techniczno-produktową. Pozostałe dokumenty produktu:

| Dokument | Zawartość |
|---|---|
| **[LICENSE.md](LICENSE.md)** | warunki licencyjne Danaco Console; prawa producenta i twórcy; zakres dozwolonego użycia |
| **[INSTRUKCJA-UZYTKOWANIA.md](INSTRUKCJA-UZYTKOWANIA.md)** | praca z gotową instalacją: strona główna, środowiska, okno komunikacji, okna równoległe, okno konfiguracji, dostępy, konta i tożsamość — wraz z zakresem realnie dostępnym w v2.0 |
| **[INSTALACJA-I-KONFIGURACJA.md](INSTALACJA-I-KONFIGURACJA.md)** | pełna procedura instalacji, budowania warstw, doprowadzenia świeżej instalacji do pracy, wykaz zmiennych środowiska i przełączników, rozwiązywanie problemów |

**Wspólna podstawa faktograficzna.** Wszystkie cztery dokumenty produktu opisują
stan bieżącego stanu repozytorium, zamykającej prace naprawcze — każdy deklaruje tę
podstawę we własnej metryce nagłówkowej. Cztery zagadnienia, przy których fala
naprawcze zmieniły stan wobec rewizji odniesienia stanu sprzed prac naprawczych, są w całej
dokumentacji opisane zgodnie:

| Zagadnienie | Stan przed pracami naprawczymi (stanu sprzed prac naprawczych) | Stan opisywany zgodnie przez wszystkie dokumenty (bieżącego stanu repozytorium) |
|---|---|---|
| Wytworzenie katalogu `C:\DanacoConsole_App` | nie wytwarzał go żaden skrypt; wypełnienie ręczne | wytwarzają je `wydanie.sh` + `pakowanie.sh`; wypełnienie ręczne pozostaje drogą zastępczą ([17.3](#173-wydanie-produktu-do-cdanacoconsole_app), [13.1](#131-model-czterech-katalogów)) |
| Budowanie pakietu interfejsu przez powłokę | brak `beforeBuildCommand` | oba klucze obecne, łańcuch domknięty ([3.3](#33-powłoka-natywna--tauri)) |
| Pulpit Mission Control | dane przykładowe | dane z rdzenia, część sekcji bez źródła w kontrakcie ([6](#6-środowiska)) |
| Rodzaje kanału modelu | trójka `cli` / `api` / `echo` | trójka rodzajów **wykonawczych** obok czwórki wartości kolumny schematu — pojęcia rozdzielone ([7.1](#71-rejestr-kanałów-sterowany-danymi)) |

Punkt wejścia budowania rdzenia (`./server/cmd/danaco-console` z katalogu `budowa`)
jest w niniejszym dokumencie **wyrównany** do brzmienia z [INSTRUKCJA-UZYTKOWANIA.md](INSTRUKCJA-UZYTKOWANIA.md)
rozdz. 2.2, zgodnego z drzewem repozytorium i ze skryptem wydania.

Dokumentacja techniczna towarzysząca kodowi repozytorium budowy, w katalogu `docs\`
korzenia repozytorium (nie w `budowa\` — rozdz. 13.2 wyżej):

| Dokument | Zawartość |
|---|---|
| [Kontrakty komunikacji](architektura/kontrakty-komunikacji.md) | pełny opis kontraktu komunikacyjnego: komendy, zdarzenia, ładunki, narzędzia modelu — źródło normatywne `budowa/shared/contract.json` (1119 komend, 68 obszarów, 552 struktury, 377 wyliczeń, 78 zdarzeń, 8 kodów błędów) |
| [Model danych](architektura/model-danych.md) | model danych: 68 tabel docelowych (Załącznik B) wobec 35 tabel wdrożonych w bieżącym schemacie |
| [Izolacja i zależności](architektura/izolacja-i-zaleznosci.md) | mechanizm izolacji: jedenaście punktów, osiem poziomów |
| [Opracowania modułów](SPIS-OPRACOWAN.md#15-moduly) | katalog `moduly/` — piętnaście opracowań, jedno na moduł platformy |
| [Opracowania środowisk](SPIS-OPRACOWAN.md#16-srodowiska) | katalog `srodowiska/` — cztery opracowania, jedno na środowisko platformy |
| [Elementy okien](interfejs-uzytkownika/elementy-okien.md) | katalog stanów okna (pusty, ładowania, błędu) i wymogów nawigacji |

Streszczenia decyzji architektonicznych są jawne w rozdz. 3.8. Materiały
koncepcyjne i decyzyjne spoza repozytorium mogą pozostawać poza kontrolą wersji
i zawierać dane wrażliwe:

| Dokument | Zawartość | W repozytorium |
|---|---|---|
| `design/` | pakiet design v2.0: kierunek, marka, żetony, arkusze, ikony, makiety, księga marki | tak |

---

## 20. Kwestie do decyzji Operatora

Poniższe kwestie są w produkcie **nierozstrzygnięte**. Dokument ich nie przesądza
i nie proponuje rozwiązania w miejsce Operatora.

1. **Sposób doprowadzenia świeżej instalacji do pracy.** Czy kanał główny i konto
   mają powstawać zaczynem migracji, kreatorem pierwszego uruchomienia w interfejsie,
   czy przez wpięcie istniejącego modułu zapewnienia kanału do ścieżki produkcyjnej
   ([17.2](#172-doprowadzenie-świeżej-instalacji-do-pracy--czynność-obowiązkowa)).
2. **Miejsce słownika kluczy ustawień** — kontrakt `shared/` czy katalog w bazie —
   oraz kierunek uzgodnienia trzech kluczy panelu sterowania z kluczami rdzenia
   ([8.4](#84-które-klucze-realnie-sterują-wykonaniem)).
3. **Semantyka `config.get` bez wskazania poziomu**: polityka efektywna czy komplet
   zapisów ze wszystkich poziomów ([8.3](#83-osie-rozstrzygania)).
4. **Identyfikator kanału w kontrakcie** — kod wiersza czy numer wiersza; wymaga
   jednego rozstrzygnięcia i wyrównania obu ścieżek ([7.1](#71-rejestr-kanałów-sterowany-danymi)).
5. **Pojęcie modułu w kontrakcie** — przemianowanie listy niosącej kody środowisk
   albo dodanie osobnej listy kodów modułów ([5](#5-moduły)).
6. **Podsystem procesów okien** — wpięcie jako droga uruchamiania procesu modelu albo
   usunięcie na rzecz drogi przez rejestr kanałów ([9.4](#94-pętla-koordynatorwykonawca)).
7. **Kształt prowenancji** — scalenie dwóch niezgodnych struktur JSON, jednej
   nadawanej przez kanał główny, drugiej przez warstwę modeli. Dopisanie rodzajów
   fragmentu `provenance` i `account` do wyliczenia `ChunkKind` kontraktu — druga
   część tej samej kwestii na rewizji audytu stanu sprzed prac naprawczych — jest już wykonane
   (luka zamknięta pracami porządkowymi nad dokumentacją) i nie jest przedmiotem tej decyzji
   ([7.7](#77-prowenancja-wywołania)).
8. **Podniesienie priorytetu w kolejce** — wartość akcji w kontrakcie albo usunięcie
   funkcji z pulpitu ([9.5](#95-kolejki)).
9. **Liczba okien na scenie** — trzy okna wobec czterech okien ról wymaganych przez
   MultitaskingAI ([9.3](#93-okna-równoległe-i-role)).
10. **Umiejscowienie stałych mostu MCP po SSH.** Zapis decyzyjny istnieje (status
    `Zaproponowana`); nierozstrzygnięte pozostaje, czy przy ratyfikacji ścieżka
    skryptu, program połączenia, użytkownik i port przechodzą ze stałych kodu
    do wiersza `punkt_dostepu` (rekomendacja dokumentacji, zgodna z zasadami
    konfiguracji zamiast logiki i sterowania danymi), czy pozostają stałymi produktu ([12.2](#122-mosty-mcp-po-ssh)).
11. **Instalator NSIS** — jedyny krok wydania faktycznie odłożony. Wydanie produktu
    do `C:\DanacoConsole_App` jest domknięte skryptem `wydanie.sh`; nierozstrzygnięte
    pozostaje, czy do zakresu v2.0 wchodzi złożenie instalatora (cel `nsis`
    zadeklarowany w `tauri.conf.json`, wymaga narzędzia NSIS) oraz czy w instalatorze
    rdzeń ma zostać osadzony jako zasób powłoki, czy — jak w wydaniu bieżącym —
    leżeć obok niej osobnym plikiem wykonywalnym
    ([17.3](#173-wydanie-produktu-do-cdanacoconsole_app)).
12. **Ratyfikacja dziewięciu decyzji o statusie `Zaproponowana`** — część z nich jest
    już realizowana przez kod, co dokumentacja odnotowuje jako nieprawidłowość
    trybu. Rekomendacja dla wszystkich dziewięciu rozstrzygnięć została
    przedłożona w pracach porządkowych nad dokumentacją
    ([3.8](#38-decyzje-architektoniczne)).
13. **Sposób adresowania rozgłoszeń** — czy zdarzenia i fragmenty mają docierać do
    wszystkich połączeń, czy podlegać zawężeniu do konta i urządzenia
    ([3.6](#36-transport-websocket-i-json)).
14. **Rozdział pojęć „rodzaj kanału”** — czy wartości kolumny schematu
    (`cli` · `api` · `sdk` · `lokalny`, dziś powielone w kontrakcie jako
    `KnownChannelKinds`) mają zostać nazewniczo oddzielone od kluczy adapterów
    wykonawczych (`cli` · `api` · `echo`), czy różnica pozostaje wyjaśniona wyłącznie
    dokumentacją; wiąże się z tym pytanie o adaptery dla rodzajów `sdk` i `lokalny`
    ([7.1](#71-rejestr-kanałów-sterowany-danymi)).
15. **Tryb utrzymania czterech dokumentów produktu wobec kolejnych fal** — komplet
    dokumentacji opisuje dziś zgodnie stan bieżącego stanu repozytorium, każdy dokument
    deklaruje tę podstawę we własnej metryce. Rozstrzygnięcia wymaga tryb dalszy:
    czy dokumentacja jest przepisywana po każdej fali budowy, czy zamrożona na
    rewizji wydania i odświeżana wyłącznie przy publikacji. Przy trwających falach
    modułowych rozjazd opisu ze stanem kodu narasta z każdą godziną
    ([13.1](#131-model-czterech-katalogów), [19](#19-pozostała-dokumentacja)).

---

*Koniec dokumentu. Opis produktu — Stan wdrożenia, wersja 2.0, 2026-08-20.*

---
*Danaco Console — Platforma AI Workspace OS · v2.0 · status Deweloperski*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz.*
*Warunki korzystania: [Licencja produktu](LICENSE.md). Kontakt: support@danaco-group.pl*
