# Danaco Console — Instalacja i konfiguracja

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
| **Tytuł** | Instalacja i konfiguracja Danaco Console |
| **Klasa dokumentu** | Stan wdrożenia |
| **Odbiorcy** | Operator — osoba instalująca, konfigurująca i utrzymująca Danaco Console |
| **Przeznaczenie** | Prowadzi Operatora od pustej stacji roboczej do produktu zdolnego wykonać pierwszą turę modelu: wymagania, budowa ze źródeł, złożenie katalogu produktu, konfiguracja rdzenia, bazy, kanałów modelu, kont i narzędzi oraz diagnostyka. Na jego podstawie powstaje działająca instalacja. |
| **Zakres** | wymagania systemowe, środowiska uruchomieniowe i zależności, przygotowanie stacji Windows, instalacja, struktura katalogów, pierwsze uruchomienie, konfiguracja rdzenia i klienta, bazy danych, kanałów modelu, kont i tożsamości, zmienne środowiskowe, katalogi robocze i nadania, środowiska lokalne i zdalne, mosty MCP, dostawcy modeli, sesje i okna, weryfikacja instalacji, procedura doprowadzenia świeżej instalacji do pracy, diagnostyka, aktualizacja i deinstalacja, załączniki |
| **Poza zakresem** | budowa i architektura produktu — [Opis produktu](README.md); codzienna praca z aplikacją — [Instrukcja użytkowania](INSTRUKCJA-UZYTKOWANIA.md); warunki korzystania — [Licencja produktu](LICENSE.md) |
| **Dokument nadrzędny** | [Spis opracowań](SPIS-OPRACOWAN.md) |
| **Dokumenty powiązane** | [Opis produktu](README.md) · [Instrukcja użytkowania](INSTRUKCJA-UZYTKOWANIA.md) · [Licencja produktu](LICENSE.md) · [Standard redakcyjny i językowy](STANDARD-REDAKCYJNY-I-JEZYKOWY.md) |
| **Prototypy odniesienia** | nie dotyczy — dokument klasy Stan wdrożenia opisuje instalację produktu zbudowanego, nie projekt okna |
| **Źródła normatywne** | `budowa/scripts/wydanie.sh` · `budowa/scripts/pakowanie.sh` · `budowa/server/internal/konfiguracja/` · `budowa/server/internal/store/` · `budowa/desktop/src-tauri/` · `budowa/shared/contract.json` |
| **Zasada nadrzędna** | Dokument nie obiecuje zachowań, których nie da się wskazać w kodzie. |

**Zastrzeżenie o statusie wydania.** Status „Deweloperski” nie jest ozdobnikiem. Danaco Console w bieżącej
wersji jest kompletnym fundamentem technicznym (kontrakt, rdzeń, warstwa danych, transport, powłoka natywna),
lecz część funkcji zapowiedzianych koncepcją nie ma jeszcze implementacji wykonawczej. Niniejszy dokument
konsekwentnie odróżnia to, co **działa**, od tego, co jest **przewidziane koncepcją, a w v2.0
niezaimplementowane**. Każde takie miejsce jest oznaczone wprost. Dokument nie obiecuje zachowań, których nie
da się wskazać w kodzie.

---

## Spis treści

1. [Wprowadzenie, zakres i sposób czytania dokumentu](#1-wprowadzenie-zakres-i-sposób-czytania-dokumentu)
   - [1.1 Czym jest Danaco Console z punktu widzenia instalacji](#11-czym-jest-danaco-console-z-punktu-widzenia-instalacji)
   - [1.2 Trzy części produktu](#12-trzy-części-produktu)
   - [1.3 Konwencje dokumentu](#13-konwencje-dokumentu)
2. [Wymagania systemowe](#2-wymagania-systemowe)
   - [2.1 System operacyjny](#21-system-operacyjny)
   - [2.2 Sprzęt](#22-sprzęt)
   - [2.3 Sieć i porty](#23-sieć-i-porty)
   - [2.4 Uprawnienia i konto systemowe](#24-uprawnienia-i-konto-systemowe)
3. [Wymagane środowiska uruchomieniowe i zależności](#3-wymagane-środowiska-uruchomieniowe-i-zależności)
   - [3.1 Zestawienie zbiorcze](#31-zestawienie-zbiorcze)
   - [3.2 Go](#32-go)
   - [3.3 Node.js i npm](#33-nodejs-i-npm)
   - [3.4 Rust i Tauri CLI](#34-rust-i-tauri-cli)
   - [3.5 Program `claude` — Claude Code CLI](#35-program-claude--claude-code-cli)
   - [3.6 Powłoka Bash](#36-powłoka-bash)
   - [3.7 Zależności opcjonalne](#37-zależności-opcjonalne)
   - [3.8 Zależności biblioteczne rdzenia](#38-zależności-biblioteczne-rdzenia)
4. [Przygotowanie stacji roboczej Windows](#4-przygotowanie-stacji-roboczej-windows)
   - [4.1 Kolejność przygotowania](#41-kolejność-przygotowania)
   - [4.2 Weryfikacja narzędzi](#42-weryfikacja-narzędzi)
   - [4.3 Kodowanie znaków i długie ścieżki](#43-kodowanie-znaków-i-długie-ścieżki)
   - [4.4 Zapora systemowa i oprogramowanie ochronne](#44-zapora-systemowa-i-oprogramowanie-ochronne)
5. [Instalacja Danaco Console](#5-instalacja-danaco-console)
   - [5.1 Stan faktyczny dróg instalacji](#51-stan-faktyczny-dróg-instalacji)
   - [5.2 Droga podstawowa — skrypt wydania](#52-droga-podstawowa--skrypt-wydania)
   - [5.3 Co dokładnie robi skrypt wydania](#53-co-dokładnie-robi-skrypt-wydania)
   - [5.4 Wariant rdzeń + klient (bez powłoki)](#54-wariant-rdzeń--klient-bez-powłoki)
   - [5.5 Droga ręczna — budowa części osobno](#55-droga-ręczna--budowa-części-osobno)
   - [5.6 Instalator NSIS](#56-instalator-nsis)
6. [Struktura instalacji i model katalogów](#6-struktura-instalacji-i-model-katalogów)
   - [6.1 Rozdział źródeł od wyjścia](#61-rozdział-źródeł-od-wyjścia)
   - [6.2 Zawartość katalogu produktu](#62-zawartość-katalogu-produktu)
   - [6.3 Katalog danych rdzenia](#63-katalog-danych-rdzenia)
   - [6.4 Katalog profili kanału głównego](#64-katalog-profili-kanału-głównego)
   - [6.5 Reguły higieny katalogu produktu](#65-reguły-higieny-katalogu-produktu)
7. [Pierwsze uruchomienie](#7-pierwsze-uruchomienie)
   - [7.1 Uruchomienie rdzenia](#71-uruchomienie-rdzenia)
   - [7.2 Interfejs z pakietu serwowanego przez rdzeń](#72-interfejs-z-pakietu-serwowanego-przez-rdzeń)
   - [7.3 Interfejs z serwera deweloperskiego Vite](#73-interfejs-z-serwera-deweloperskiego-vite)
   - [7.4 Powłoka natywna Tauri](#74-powłoka-natywna-tauri)
   - [7.5 Czego oczekiwać po pierwszym starcie](#75-czego-oczekiwać-po-pierwszym-starcie)
8. [Konfiguracja rdzenia i klienta](#8-konfiguracja-rdzenia-i-klienta)
   - [8.1 Warstwy konfiguracji](#81-warstwy-konfiguracji)
   - [8.2 Przełączniki wiersza poleceń](#82-przełączniki-wiersza-poleceń)
   - [8.3 Rola procesu](#83-rola-procesu)
   - [8.4 Port i adres gniazda](#84-port-i-adres-gniazda)
   - [8.5 Ustalanie adresu rdzenia po stronie klienta](#85-ustalanie-adresu-rdzenia-po-stronie-klienta)
9. [Konfiguracja bazy danych](#9-konfiguracja-bazy-danych)
   - [9.1 Plik bazy i jego położenie](#91-plik-bazy-i-jego-położenie)
   - [9.2 Migracje](#92-migracje)
   - [9.3 Pragmy i tryb dziennika](#93-pragmy-i-tryb-dziennika)
   - [9.4 Kopia zapasowa i przeniesienie](#94-kopia-zapasowa-i-przeniesienie)
10. [Konfiguracja kanałów modelu](#10-konfiguracja-kanałów-modelu)
   - [10.1 Rejestr kanałów sterowany danymi](#101-rejestr-kanałów-sterowany-danymi)
   - [10.2 Komendy kontraktu do zarządzania rejestrem](#102-komendy-kontraktu-do-zarządzania-rejestrem)
   - [10.3 Kanał `cli` — Claude Code](#103-kanał-cli--claude-code)
   - [10.4 Kanał `api`](#104-kanał-api)
   - [10.5 Kanał `echo`](#105-kanał-echo)
   - [10.6 Kanał SSH — stan faktyczny](#106-kanał-ssh--stan-faktyczny)
   - [10.7 Kanał HTTP — stan faktyczny](#107-kanał-http--stan-faktyczny)
11. [Konta, pula rotacji i profile tożsamości](#11-konta-pula-rotacji-i-profile-tożsamości)
   - [11.1 Tabela `konto`](#111-tabela-konto)
   - [11.2 `CLAUDE_CONFIG_DIR` per konto](#112-claude_config_dir-per-konto)
   - [11.3 Pula rotacji i wyczerpanie limitu](#113-pula-rotacji-i-wyczerpanie-limitu)
   - [11.4 Katalog profili a katalog konfiguracji konta](#114-katalog-profili-a-katalog-konfiguracji-konta)
   - [11.5 Nakładka tożsamości](#115-nakładka-tożsamości)
12. [Zmienne środowiskowe](#12-zmienne-środowiskowe)
   - [12.1 Wykaz zupełny](#121-wykaz-zupełny)
   - [12.2 Opis pojedynczych zmiennych](#122-opis-pojedynczych-zmiennych)
   - [12.3 Zmienne świadomie nieobecne](#123-zmienne-świadomie-nieobecne)
   - [12.4 Zmienne przekazywane procesowi kanału](#124-zmienne-przekazywane-procesowi-kanału)
   - [12.5 Plik `.env` i zmienne budowania interfejsu](#125-plik-env-i-zmienne-budowania-interfejsu)
   - [12.6 Zmienne wdrożenia serwerowego i składania wydania](#126-zmienne-wdrożenia-serwerowego-i-składania-wydania)
13. [Konfiguracja katalogów roboczych i nadań](#13-konfiguracja-katalogów-roboczych-i-nadań)
   - [13.1 Katalog roboczy sesji](#131-katalog-roboczy-sesji)
   - [13.2 Katalogi dodatkowe okna](#132-katalogi-dodatkowe-okna)
   - [13.3 Punkty dostępu i nadania](#133-punkty-dostępu-i-nadania)
14. [Środowiska lokalne i zdalne](#14-środowiska-lokalne-i-zdalne)
   - [14.1 Rola `hub` i rola `agent`](#141-rola-hub-i-rola-agent)
   - [14.2 Zasięg wykonania — stan faktyczny](#142-zasięg-wykonania--stan-faktyczny)
   - [14.3 Praca zdalna dostępna dzisiaj](#143-praca-zdalna-dostępna-dzisiaj)
15. [Konfiguracja narzędzi — mosty MCP](#15-konfiguracja-narzędzi--mosty-mcp)
   - [15.1 Model mostu](#151-model-mostu)
   - [15.2 Punkt dostępu rodzaju `mcpBridge`](#152-punkt-dostępu-rodzaju-mcpbridge)
   - [15.3 Przekazanie mostu do procesu kanału](#153-przekazanie-mostu-do-procesu-kanału)
   - [15.4 Stałe mostu i ich konsekwencje](#154-stałe-mostu-i-ich-konsekwencje)
16. [Konfiguracja dostawców modeli](#16-konfiguracja-dostawców-modeli)
   - [16.1 Wiersz rejestru kanału jako opis dostawcy](#161-wiersz-rejestru-kanału-jako-opis-dostawcy)
   - [16.2 Kształt `parametry_json`](#162-kształt-parametry_json)
   - [16.3 Poświadczenia dostawcy](#163-poświadczenia-dostawcy)
17. [Konfiguracja sesji i okien operacyjnych](#17-konfiguracja-sesji-i-okien-operacyjnych)
   - [17.1 Ustawienia wpływające na wywołanie](#171-ustawienia-wpływające-na-wywołanie)
   - [17.2 Ustawienia zapisywane, lecz nieaplikowane](#172-ustawienia-zapisywane-lecz-nieaplikowane)
   - [17.3 Ciągłość rozmowy](#173-ciągłość-rozmowy)
18. [Konfiguracja komponentów dostępnych w bieżącej wersji](#18-konfiguracja-komponentów-dostępnych-w-bieżącej-wersji)
   - [18.1 Komponenty konfigurowalne](#181-komponenty-konfigurowalne)
   - [18.2 Komponenty przewidziane koncepcją, niezaimplementowane](#182-komponenty-przewidziane-koncepcją-niezaimplementowane)
19. [Weryfikacja poprawności instalacji](#19-weryfikacja-poprawności-instalacji)
   - [19.1 Kontrole źródeł](#191-kontrole-źródeł)
   - [19.2 Brama jakości](#192-brama-jakości)
   - [19.3 Kontrola produktu złożonego](#193-kontrola-produktu-złożonego)
   - [19.4 Kontrola połączenia i kontraktu](#194-kontrola-połączenia-i-kontraktu)
20. [Procedura doprowadzenia świeżej instalacji do stanu gotowego do pracy](#20-procedura-doprowadzenia-świeżej-instalacji-do-stanu-gotowego-do-pracy)
   - [20.1 Dlaczego procedura jest potrzebna](#201-dlaczego-procedura-jest-potrzebna)
   - [20.2 Krok 1 — profil konta kanału głównego](#202-krok-1--profil-konta-kanału-głównego)
   - [20.3 Krok 2 — start rdzenia](#203-krok-2--start-rdzenia)
   - [20.4 Krok 3 — zasianie wiersza kanału](#204-krok-3--zasianie-wiersza-kanału)
   - [20.5 Krok 4 — zasianie konta](#205-krok-4--zasianie-konta)
   - [20.6 Krok 5 — kontrola zasiania i pierwsza tura](#206-krok-5--kontrola-zasiania-i-pierwsza-tura)
   - [20.7 Wariant awaryjny — zapis wprost do bazy](#207-wariant-awaryjny--zapis-wprost-do-bazy)
   - [20.8 Ograniczenie procedury](#208-ograniczenie-procedury)
21. [Diagnostyka typowych problemów](#21-diagnostyka-typowych-problemów)
   - [21.1 Gdzie szukać śladów](#211-gdzie-szukać-śladów)
   - [21.2 Pusty rejestr kanałów](#212-pusty-rejestr-kanałów)
   - [21.3 Pusta pula kont — komunikat mylący](#213-pusta-pula-kont--komunikat-mylący)
   - [21.4 Brak programu `claude`](#214-brak-programu-claude)
   - [21.5 Port zajęty](#215-port-zajęty)
   - [21.6 Powłoka bez pakietu interfejsu — „pusta strona”](#216-powłoka-bez-pakietu-interfejsu--pusta-strona)
   - [21.7 Baza zablokowana albo niezgodna](#217-baza-zablokowana-albo-niezgodna)
   - [21.8 Historia znika po odświeżeniu strony](#218-historia-znika-po-odświeżeniu-strony)
22. [Aktualizacja, ponowne uruchomienie i deinstalacja](#22-aktualizacja-ponowne-uruchomienie-i-deinstalacja)
   - [22.1 Co przeżywa ponowne uruchomienie rdzenia](#221-co-przeżywa-ponowne-uruchomienie-rdzenia)
   - [22.2 Ponowne wydanie i aktualizacja](#222-ponowne-wydanie-i-aktualizacja)
   - [22.3 Ponowne uruchomienie](#223-ponowne-uruchomienie)
   - [22.4 Deinstalacja](#224-deinstalacja)
23. [Załączniki](#23-załączniki)
   - [23.1 Wykaz plików konfiguracyjnych i ich rola](#231-wykaz-plików-konfiguracyjnych-i-ich-rola)
   - [23.2 Komendy kontraktu przywołane w dokumencie](#232-komendy-kontraktu-przywołane-w-dokumencie)
   - [23.3 Tabele bazy istotne dla instalacji](#233-tabele-bazy-istotne-dla-instalacji)
   - [23.4 Wykaz kwestii pozostawionych do decyzji Operatora](#234-wykaz-kwestii-pozostawionych-do-decyzji-operatora)
   - [23.5 Słownik pojęć użytych w dokumencie](#235-słownik-pojęć-użytych-w-dokumencie)

---

## 1. Wprowadzenie, zakres i sposób czytania dokumentu

### 1.1 Czym jest Danaco Console z punktu widzenia instalacji

Danaco Console to środowisko operacyjne pracy z modelami językowymi, zbudowane jako aplikacja instalowana
na stacji Operatora. Z perspektywy wdrożeniowej nie jest to pojedynczy plik wykonywalny ani usługa sieciowa
kupowana w chmurze, lecz zestaw trzech współpracujących części, budowanych ze źródeł i składanych w jeden
katalog produktu. Instalacja u Operatora idzie przez kreator instalacji sześciu kroków
(`Danaco Console — Instalator_2.0.0_x64-setup.exe`, rozdział 5.1), a rdzeń stoi na serwerze wdrożenia
z pakietu `.deb` (rozdział 12.6); budowa ze źródeł opisana niżej jest drogą stacji budującej.

Konsekwencja jest praktyczna: stacja, na której Danaco Console ma powstać, musi mieć zainstalowane
środowiska uruchomieniowe kompilatorów (Go, Node.js, opcjonalnie Rust). Stacja, na której produkt ma jedynie
pracować, potrzebuje wyłącznie skopiowanego katalogu produktu oraz zewnętrznego programu `claude`
(Claude Code CLI), bez którego kanał główny modelu nie ma czego uruchomić.

### 1.2 Trzy części produktu

| Część | Nazwa w dokumentacji | Technologia | Artefakt |
| --- | --- | --- | --- |
| Serwer aplikacyjny | **rdzeń** | Go (moduł `danacoconsole`) | `danaco-console.exe` |
| Interfejs użytkownika | **klient** / **interfejs** | TypeScript + Vite | katalog `client/dist` |
| Okno natywne systemu | **powłoka** | Rust + Tauri 2 | `Danaco Console.exe` |

Do tego dochodzi czwarty, niewykonywalny składnik: **kontrakt** (`budowa/shared`) — jedno źródło prawdy
nazw komend, zdarzeń i narzędzi, z którego generowane są `contract.ts` (dla klienta) i `contract.go`
(dla rdzenia). Kontrakt nie jest osobno instalowany, lecz wkompilowuje się w obie strony; jego rola
w instalacji sprowadza się do tego, że wszystkie nazwy komend przywoływane w niniejszym dokumencie
pochodzą właśnie stamtąd i nie wolno ich zapisywać z pamięci.

Rdzeń jest jedynym procesem, który dotyka bazy danych, uruchamia procesy kanału modelu i nasłuchuje
transportu WebSocket. Klient jest kompletnie bezstanowy wobec trwałości — nie zapisuje niczego poza
przeglądarką. Powłoka nie zawiera logiki produktu: stawia rdzeń w tle, otwiera okno i ładuje do niego
pakiet interfejsu.

### 1.3 Konwencje dokumentu

Dokument opisuje **rzeczywisty stan produktu** w bieżącym stanie repozytorium — po zamknięciu prac naprawczych
uruchomionych na podstawie audytu technicznego. Tam, gdzie prace zmieniły stan opisany pierwotnie
przez audyt, dokument przywołuje obie rewizje wprost, żeby Operator widział, co było, a co jest.
Stosowane są następujące oznaczenia stanu:

- **[DZIAŁA]** — funkcja potwierdzona w kodzie i w wyniku audytu; można na niej polegać przy wdrożeniu.
- **[DZIAŁA CZĘŚCIOWO]** — droga wykonawcza istnieje, lecz jest niekompletna albo wymaga obejścia;
  dokument opisuje dokładnie granicę działania.
- **[NIEZINTEGROWANE]** — kod istnieje po jednej lub obu stronach, lecz nie tworzy zamkniętej drogi;
  Operator nie osiągnie efektu bez ingerencji ręcznej.
- **[ATRAPA]** — element widoczny w interfejsie (pozycja menu, kafel, przycisk), za którym nie stoi
  realizacja.
- **[BRAK]** — przewidziane koncepcją, w bieżącej wersji niezaimplementowane: nazwa, menu albo definicja
  istnieje, brak wykonawcy. Takiej pozycji nie wolno traktować jako funkcji.
- **[DO DECYZJI OPERATORA]** — kwestia otwarta, której dokumentacja nie rozstrzyga, bo rozstrzygnięcie
  należy do Operatora, nie do autora dokumentu.

Zestaw sześcioelementowy jest wspólny dla wszystkich opracowań klasy Stan wdrożenia — wiąże go
[Standard redakcyjny i językowy](STANDARD-REDAKCYJNY-I-JEZYKOWY.md) rozdz. 2.6, a posługują się nim
także [Opis produktu](README.md) i [Instrukcja użytkowania](INSTRUKCJA-UZYTKOWANIA.md).

Ścieżki źródeł podawane są względem katalogu `budowa/` w repozytorium. Ścieżki produktu podawane są
bezwzględnie, w postaci Windows albo w postaci z ukośnikami zwykłymi tam, gdzie trafiają do powłoki Bash.
Polecenia oznaczone jako niezweryfikowane uruchomieniem są opatrzone wyraźną adnotacją — dokument nie
składa kategorycznych obietnic co do zachowań, których nie dało się wykonać w toku prac.

**Podstawa faktograficzna.** Opis opiera się na kodzie repozytorium w stanie bieżącym —
**po zamknięciu prac naprawczych** uruchomionych na podstawie audytu technicznego. Podstawa
ta obejmuje `budowa/scripts/wydanie.sh` i `pakowanie.sh`, bramę mierzącą wywołania i dym
pionu, komplet stylów CSS v2.0, pulpit i stronę główną zasilane z rdzenia oraz porządek
dokumentacyjny.

**Zakres platformowy.** Dokument opisuje stację Operatora pod Windows 11. Rdzeń kompiluje
się także na innych systemach, lecz opisana jest wyłącznie ścieżka Windows.

**Język produktu i dokumentacji.** Produkt i jego dokumentacja są w całości polskojęzyczne.

---

## 2. Wymagania systemowe

### 2.1 System operacyjny

Docelową platformą Operatora jest **Windows 11** (weryfikowano na Windows 11 Pro, kompilacja 10.0.26200).
Produkt nie zawiera żadnego kodu warunkowanego wersją systemu, natomiast trzy elementy wiążą go
z Windows w sposób praktyczny:

1. Domyślny katalog danych rdzenia rozwiązywany jest do `%LOCALAPPDATA%\DanacoConsole` — droga
   charakterystyczna dla Windows (rozdział 6.3).
2. Skrypty budowy i wydania (`budowa/scripts/*.sh`) są skryptami powłoki Bash; na Windows wymagają
   Git Bash albo równoważnej powłoki POSIX (rozdział 3.6).
3. Artefakty produktu noszą rozszerzenie `.exe`, a skrypt pakowania składa je pod nazwami
   `danaco-console.exe` i `Danaco Console.exe`.

Rdzeń napisany jest w Go bez zależności wymagających kompilatora C (sterownik SQLite `modernc.org/sqlite`
jest napisany w całości w języku Go), więc kompiluje się także na Linux i macOS. Niniejszy dokument opisuje
jednak wyłącznie ścieżkę Windows i nie orzeka o zachowaniu na innych systemach.

### 2.2 Sprzęt

Produkt nie deklaruje w kodzie żadnych progów sprzętowych. Poniższe wartości wynikają z charakteru
budowanych składników, nie z kontroli wbudowanej w program:

| Zasób | Do samej pracy produktu | Do budowy ze źródeł |
| --- | --- | --- |
| Procesor | dowolny x86-64 | x86-64, wielordzeniowy skraca kompilację Rust |
| Pamięć operacyjna | 4 GB (rdzeń jest lekki, ciężar biorą na siebie procesy kanału) | 8 GB — kompilacja Tauri/Rust jest najbardziej pamięciożerna |
| Dysk | ok. 100 MB na katalog produktu i bazę | kilka GB: `node_modules`, pamięć podręczna Cargo, moduły Go |

Rzeczywiste zapotrzebowanie w pracy zależy przede wszystkim od liczby równolegle otwartych okien
operacyjnych, bo każda tura kanału `cli` uruchamia osobny proces zewnętrznego programu `claude`.

### 2.3 Sieć i porty

Rdzeń nasłuchuje transportu WebSocket na **porcie 17870**. Wartość jest domyślną w
`server/internal/konfiguracja` i tę samą liczbę zna klient w `client/src/polaczenie/adres-rdzenia.ts`.
Port wybrano świadomie poza pulą zajętą przez inne systemy Danaco (8090, 8091, 8093, 8095, 8760, 8766,
8771, 8772, 8790, 9980, 18765).

Ruch produktu jest w całości ruchem pętli zwrotnej: klient łączy się z rdzeniem na tej samej maszynie.
Otwieranie portu 17870 na zewnątrz nie jest wymagane do pracy i nie jest zalecane; wdrożenie serwerowe
wystawia rdzeń przez Caddy z jawnym wymogiem logowania (rozdział 2.4, rozstrzygnięcie 29).

Dostęp do sieci publicznej jest natomiast potrzebny:

- programowi `claude` — do rozmowy z usługą modelu,
- kanałom rodzaju `api` — do punktu końcowego dostawcy,
- procesowi budowy — do pobrania modułów Go, pakietów npm i skrzyń Rust.

### 2.4 Uprawnienia i konto systemowe

Praca produktu nie wymaga uprawnień administratora. Rdzeń zapisuje wyłącznie do katalogu danych
w profilu użytkownika i do katalogów roboczych wskazanych przez Operatora. Uprawnienia podwyższone mogą
być potrzebne jednorazowo przy instalowaniu środowisk uruchomieniowych (Go, Node.js, Rust) oraz przy
budowie instalatora NSIS.

Na stacji deweloperskiej rdzeń nasłuchuje na pętli zwrotnej i tam zdejmuje wymóg logowania sam z siebie:
każdy proces tej maszyny jest Operatorem. Nasłuch na interfejsie sieciowym rozstrzyga rozstrzygnięcie 29
(`prowadzenie/decyzje.md`): rdzeń wdrożenia stoi na `127.0.0.1:17870` z jawnym
`DANACO_WYMOG_LOGOWANIA=true` w jednostce systemd, a na świat wystawia go odwrotne proxy Caddy jako
`console.danaco-group.pl:443` z TLS z ACME. Port 17870 nie jest otwierany na zaporze;
`DANACO_WSZYSTKIE_INTERFEJSY=1` nie wchodzi do wdrożenia. Bez jawnego wymogu Caddy wystawiłoby konsolę
publicznie bez bramki — rdzeń na pętli zwrotnej nie odróżnia proxy od procesu lokalnego.

---

## 3. Wymagane środowiska uruchomieniowe i zależności

### 3.1 Zestawienie zbiorcze

| Składnik | Wersja / źródło ustalenia | Potrzebny do budowy | Potrzebny do pracy |
| --- | --- | --- | --- |
| Go | `go 1.26` — deklaracja w `budowa/go.mod` | tak | nie |
| Node.js + npm | wymuszona przez `vite@6.3.5` i `typescript@5.8.3` z `client/package.json` | tak | nie |
| Rust (toolchain) | `rust-version = "1.77"` — `desktop/src-tauri/Cargo.toml` | tylko dla powłoki | nie |
| Tauri CLI (`cargo tauri`) | Tauri 2 — `tauri = { version = "2" }`, `"$schema": ".../config/2"` | tylko dla powłoki | nie |
| Program `claude` (Claude Code CLI) | zewnętrzny; ścieżka z wiersza rejestru kanału albo stała `claude` | nie | **tak**, dla kanału `cli` |
| Powłoka Bash | Git Bash lub równoważna; skrypty `budowa/scripts/*.sh` | tak | nie |
| Klient `ssh` | wyłącznie jako program mostu MCP, nie jako kanał modelu | nie | opcjonalnie |
| NSIS | `"targets": ["nsis"]` w `tauri.conf.json` | tylko dla instalatora | nie |

### 3.2 Go

Moduł rdzenia deklaruje `go 1.26`. Nazwa modułu to `danacoconsole`, a jego korzeń obejmuje jednocześnie
`server/` i `shared/` — to celowe, bo bez wspólnego modułu rdzeń nie mógłby zaimportować wygenerowanego
kontraktu `shared/contract.go`. Ścieżki importu układają się następująco: `danacoconsole/shared` dla
kontraktu, `danacoconsole/server/internal/...` dla pakietów rdzenia oraz
`danacoconsole/server/cmd/danaco-console` dla punktu wejścia.

Instalacja Go na Windows odbywa się z pakietu instalacyjnego producenta języka albo menedżerem pakietów
systemu. Po instalacji `go` musi być osiągalne w `PATH` — skrypt wydania wywołuje `go build` bez ścieżki
bezwzględnej.

### 3.3 Node.js i npm

Klient jest projektem Vite z TypeScriptem. `client/package.json` deklaruje jedną zależność wykonawczą
(`@tauri-apps/api` w wersji `2.5.0`) i cztery narzędziowe: `jsdom@^29.1.1`, `typescript@5.8.3`,
`vite@6.3.5`, `vitest@3.2.7`. Pozycja `jsdom` jest środowiskiem DOM, w którym `vitest` wykonuje
sprawdziany interfejsu — bez niej testy dotykające drzewa dokumentu nie wystartują. Zdefiniowane
polecenia to:

| Polecenie npm | Działanie |
| --- | --- |
| `npm run dev` | serwer deweloperski Vite (domyślnie port 5173) |
| `npm run build` | `tsc --noEmit && vite build` — kontrola typów, a po niej pakiet do `client/dist` |
| `npm run preview` | podgląd zbudowanego pakietu |
| `npm run typy` | sama kontrola typów (`tsc --noEmit`) |
| `npm run testy` | zestaw testów `vitest run` |

Wersję Node.js narzuca w praktyce Vite 6 — projekt nie zawiera pliku `.nvmrc` ani pola `engines`,
więc dokument nie podaje konkretnego numeru jako wymogu produktu. Należy użyć wydania Node.js zgodnego
z wymaganiami Vite 6.3, w linii aktywnego wsparcia długoterminowego. **[DO DECYZJI OPERATORA]** —
ustalenie i zapisanie w repozytorium wiążącej wersji Node.js (`engines` w `package.json` albo `.nvmrc`),
by budowa była powtarzalna na każdej stacji.

Skrypt wydania instaluje zależności poleceniem `npm ci --no-audit --no-fund`, ale **wyłącznie wtedy,
gdy katalog `client/node_modules` nie istnieje**. Wymuszenie ponownej instalacji wymaga uprzedniego
usunięcia tego katalogu.

### 3.4 Rust i Tauri CLI

Powłoka natywna jest skrzynią Rust o nazwie `danaco-console-powloka`, wersji `1.0.0` i wymaganiu
`rust-version = "1.77"`, w edycji 2021. Zależności to `tauri` 2 z cechami `tray-icon`, `image-png`
i `image-ico`, `tauri-plugin-dialog` 2 oraz `serde` 1 z cechą `derive`. Profil wydania jest zestrojony
pod rozmiar artefaktu: `opt-level = "s"`, `lto = true`, `codegen-units = 1`, `panic = "abort"`,
`strip = true`.

Do budowy powłoki potrzebne są trzy rzeczy: łańcuch narzędzi Rust (na Windows — wariant MSVC wraz
z narzędziami budowania Visual Studio, wymaganymi przez konsolidator), Tauri CLI zainstalowane
poleceniem `cargo install tauri-cli` (tak podpowiada sam skrypt wydania) oraz środowisko WebView2,
obecne standardowo w Windows 11.

Powłoka jest **składnikiem opcjonalnym**. Skrypt wydania pomija jej budowę, gdy brakuje `cargo` albo
`cargo tauri`, i wypisuje wtedy komunikat o wariancie rdzeń+klient — wydanie kończy się powodzeniem,
tylko bez pliku `Danaco Console.exe`.

### 3.5 Program `claude` — Claude Code CLI

Kanał główny modelu (rodzaj `cli`) nie zawiera własnego klienta usługi modelu. Rdzeń uruchamia
**zewnętrzny program `claude`** — narzędzie Claude Code CLI — w trybie nieinteraktywnym, ze strumieniowym
wyjściem JSON, i czyta jego strumień wyjściowy. Bez tego programu kanał `cli` nie ma czego wykonać.

Ścieżkę programu rdzeń bierze z wiersza rejestru kanału (pole parametrów kanału); gdy pole jest puste,
używa gołej nazwy `claude`, co oznacza wyszukanie w `PATH`. Istotna uwaga wynikająca z audytu:
w katalogu ustawień istnieje klucz `harness.program_claude`, lecz **nie jest on odczytywany przez
wykonanie** — ustawienie go w oknie konfiguracji nie zmieni ścieżki programu.
Jedyne dwie skuteczne drogi to wpis w wierszu rejestru kanału albo obecność `claude` w `PATH`.

Program `claude` wymaga własnej konfiguracji poświadczeń. Rdzeń nie przechowuje poświadczeń w bazie:
przekazuje procesowi zmienną `CLAUDE_CONFIG_DIR` wskazującą katalog konfiguracji przypisany do konta
(rozdział 11.2). Wartości uwierzytelniające pozostają na dysku, poza repozytorium i poza bazą produktu.

### 3.6 Powłoka Bash

Skrypty budowy, wydania i bram jakości są skryptami `bash` z nagłówkiem `#!/usr/bin/env bash` i trybem
`set -euo pipefail`. Na Windows należy uruchamiać je z **Git Bash** albo równoważnej powłoki POSIX.
Uruchomienie ich z `cmd.exe` albo z PowerShell bez pośrednictwa `bash` nie zadziała; wywołanie ma postać
`bash budowa/scripts/wydanie.sh`.

Skrypt pakowania korzysta z `find`, `cp -R`, `wc` i `mkdir -p`, więc powłoka musi dostarczać komplet
podstawowych narzędzi POSIX — Git Bash je zawiera.

### 3.7 Zależności opcjonalne

- **Klient `ssh`.** Produkt **nie posiada kanału modelu rodzaju SSH** (rozdział 10.6). `ssh` pojawia się
  wyłącznie jako program uruchamiany przez most MCP, gdy Operator skonfiguruje punkt dostępu wskazujący
  na taki most. Jest więc zależnością warunkową, wynikającą z konfiguracji, a nie z produktu.
- **NSIS.** Potrzebny wyłącznie do wytworzenia instalatora z konfiguracji `bundle` w `tauri.conf.json`
  (`"targets": ["nsis"]`). Sam produkt działa bez instalatora.
- **`sqlite3` (wiersz poleceń).** Nie jest wymagany — rdzeń niesie własny sterownik SQLite. Bywa jednak
  wygodny przy wariancie awaryjnym z rozdziału 20.7 i przy diagnostyce zawartości bazy.

### 3.8 Zależności biblioteczne rdzenia

Rdzeń jest celowo ubogi w zależności zewnętrzne. `budowa/go.mod` deklaruje trzy zależności bezpośrednie:

| Moduł | Wersja | Rola |
| --- | --- | --- |
| `github.com/coder/websocket` | `v1.8.15` | transport WebSocket rdzenia |
| `golang.org/x/sys` | `v0.47.0` | dostęp do funkcji systemowych |
| `modernc.org/sqlite` | `v1.56.0` | sterownik SQLite bez zależności od C |

Pozostałe pozycje w `go.mod` są zależnościami pośrednimi sterownika SQLite. Wybór `modernc.org/sqlite`
ma bezpośredni skutek wdrożeniowy: budowa rdzenia **nie wymaga kompilatora C ani `CGO_ENABLED=1`**,
co znacząco upraszcza przygotowanie stacji Windows.

---

## 4. Przygotowanie stacji roboczej Windows

### 4.1 Kolejność przygotowania

Zalecana kolejność czynności na czystej stacji Windows 11:

1. **Git wraz z Git Bash** — daje jednocześnie klienta repozytorium i powłokę POSIX wymaganą przez
   skrypty (rozdział 3.6).
2. **Go** w wydaniu zgodnym z `go 1.26`; kontrola: `go version`.
3. **Node.js z npm** w wydaniu zgodnym z Vite 6.3; kontrola: `node --version` i `npm --version`.
4. **Program `claude`** (Claude Code CLI) wraz z jego uwierzytelnieniem; kontrola: `claude --version`.
5. **Rust** wraz z narzędziami budowania MSVC — krok wymagany wyłącznie przy budowie powłoki natywnej — a następnie
   `cargo install tauri-cli`; kontrola: `cargo --version` i `cargo tauri --version`.
6. **NSIS** — krok wymagany wyłącznie przy składaniu instalatora.
7. Sklonowanie repozytorium do `C:\DanacoConsole_Git` oraz utworzenie katalogu produktu
   `C:\DanacoConsole_App` — model katalogów opisuje rozdział 6.

Kroki 5 i 6 można pominąć całkowicie. Wynikiem jest wtedy wariant rdzeń+klient: pełna funkcjonalność
produktu bez okna natywnego, z interfejsem otwieranym w przeglądarce.

### 4.2 Weryfikacja narzędzi

Przed pierwszą budową warto potwierdzić, że wszystkie narzędzia są osiągalne w `PATH` z tej samej
powłoki, z której będzie uruchamiany skrypt wydania. Kontrola w Git Bash:

```bash
go version
node --version
npm --version
claude --version
cargo --version        # tylko dla powłoki natywnej
cargo tauri --version  # tylko dla powłoki natywnej
```

Ostatnie dwa polecenia mogą zakończyć się niepowodzeniem bez szkody dla wydania — skrypt sam wykrywa
ich brak i przechodzi do wariantu rdzeń+klient (rozdział 5.4).

### 4.3 Kodowanie znaków i długie ścieżki

Produkt jest w całości polskojęzyczny: nazwy plików źródłowych, komunikaty, dokumentacja i komentarze
używają znaków diakrytycznych. Wpływa to na dwa ustawienia stacji:

- **Kodowanie konsoli.** Komunikaty rdzenia i skryptów zawierają polskie znaki oraz znaki ramek.
  Git Bash obsługuje UTF-8 domyślnie. W PowerShell warto ustawić `[Console]::OutputEncoding` na UTF-8,
  by komunikaty nie wyświetlały się jako sekwencje zastępcze. Jest to kwestia czytelności, nie poprawności
  działania.
- **Długie ścieżki.** Drzewo źródeł zawiera głębokie ścieżki (`budowa/server/internal/...`,
  `client/node_modules/...`). Na stacji z wyłączoną obsługą długich ścieżek zalecane jest jej włączenie
  w systemie oraz `git config --system core.longpaths true`. Bez tego operacje na `node_modules` i na
  katalogu `target` Rust mogą zawodzić.

### 4.4 Zapora systemowa i oprogramowanie ochronne

Przy pierwszym uruchomieniu `danaco-console.exe` zapora Windows zapyta o zezwolenie na nasłuch.
Ponieważ ruch produktu jest ruchem pętli zwrotnej (rozdział 2.3), wystarczy zezwolenie dla sieci
prywatnej; udostępnianie w sieciach publicznych nie jest potrzebne i nie jest zalecane.

Oprogramowanie ochronne bywa źródłem dwóch utrudnień: spowalnia budowę przez skanowanie katalogów
`node_modules` i `target` w czasie rzeczywistym oraz może kwarantannować świeżo zbudowane, niepodpisane
pliki wykonywalne. Podpis kodu rozstrzyga rozstrzygnięcie 27 (`prowadzenie/decyzje.md`): certyfikat
OV od Certum, klucz w usłudze SimplySign, podpis przez klienta SimplySign (PKCS#11, `osslsigncode`),
znacznik czasu `http://time.certum.pl`. Do zakupu certyfikatu (dokumenty spółki, czynność Właściciela)
wydania idą bez podpisu z jawnym `DANACO_PODPIS=pomijany`, wykaz niesie `podpisany: false`, a strona
„Pobierz” mówi o tym przy każdej pozycji Windows. Wyłączenia katalogów budowy z monitorowania czasu
rzeczywistego pozostają nastawą stacji, nie produktu.

---

## 5. Instalacja Danaco Console

### 5.1 Stan faktyczny dróg instalacji

Operator pobiera ze strony „Pobierz” kreator instalacji sześciu kroków — pakiet NSIS
`Danaco Console — Instalator_2.0.0_x64-setup.exe` składany skryptem
`budowa/scripts/instalka-kreatora-win-x64.sh`. Kreator rozpoznaje architekturę maszyny, czyta wykaz
wydań z kanału `https://pobierz.danaco-group.pl/` (kopia wkompilowana jest zapasem) i w kroku 5 sam
ściąga powłokę programu (`Danaco Console_2.0.0_hybryda_x64-setup.exe`, składana
`instalka-hybryda-win-x64.sh`). Rdzenia w żadnym z tych plików nie ma — stoi na serwerze wdrożenia
(rozstrzygnięcie 8, model hybrydowy). Pozycja ARM64 stoi w wykazie jako „w przygotowaniu” bez pliku
(rozstrzygnięcie 28); kreator na maszynie ARM odmawia nazwanym powodem. Poniższy opis budowy ze źródeł
dotyczy stacji budującej, nie stacji Operatora.

Do niedawna nie istniała nawet zautomatyzowana droga tego złożenia: audyt w stanie sprzed prac naprawczych odnotował
wprost, że **żaden skrypt nie wytwarzał katalogu `C:\DanacoConsole_App`**, a Operator musiał budować trzy
części osobno i kopiować je ręcznie. Stan ten zmieniła **prace naprawcze po audycie**, która wprowadziła
dwa skrypty:

| Skrypt | Odpowiedzialność |
| --- | --- |
| `budowa/scripts/wydanie.sh` | buduje wszystkie trzy części produktu i wywołuje pakowanie |
| `budowa/scripts/pakowanie.sh` | wyłącznie składa gotowe artefakty w katalogu wyjściowym; niczego nie buduje |

Niniejszy dokument opisuje drogę przez `wydanie.sh` jako **drogę podstawową**, jednocześnie oznaczając
jednoznacznie jej pochodzenie: jest to wynik prac naprawczych, nie stan zastany w chwili audytu.
Droga ręczna (rozdział 5.5) pozostaje udokumentowana, bo pozwala budować pojedynczą część bez przechodzenia
przez pełne wydanie.

### 5.2 Droga podstawowa — skrypt wydania

Wydanie uruchamia się z powłoki Bash jednym poleceniem:

```bash
bash budowa/scripts/wydanie.sh
```

Bez argumentu katalogiem wyjściowym jest `C:/DanacoConsole_App`. Inny katalog wskazuje się pierwszym
argumentem pozycyjnym:

```bash
bash budowa/scripts/wydanie.sh D:/Danaco/Produkt
```

Zachowanie skryptu modyfikuje jedna zmienna środowiskowa: `DANACO_BEZ_POWLOKI=1` pomija budowę powłoki
natywnej i kończy wydanie w wariancie rdzeń+klient.

```bash
DANACO_BEZ_POWLOKI=1 bash budowa/scripts/wydanie.sh
```

Skrypt pracuje w trybie `set -euo pipefail`: pierwsze niepowodzenie któregokolwiek kroku przerywa całość
i pozostawia katalog wyjściowy w stanie sprzed nieudanego kroku.

### 5.3 Co dokładnie robi skrypt wydania

Skrypt przechodzi pięć ponumerowanych kroków, wypisując nagłówek każdego z nich:

1. **Zależności klienta.** Wchodzi do `budowa/client`. Jeżeli katalog `node_modules` nie istnieje —
   wykonuje `npm ci --no-audit --no-fund`. Jeżeli istnieje — instalację pomija i informuje o tym
   komunikatem; wymuszenie ponownej instalacji wymaga usunięcia `node_modules`.
2. **Pakiet klienta.** Wykonuje `npm run build`, czyli `tsc --noEmit && vite build`. Kontrola typów jest
   częścią budowy: błąd typów **zatrzymuje wydanie**. Wynik trafia do `budowa/client/dist`.
3. **Rdzeń.** Wraca do `budowa/` i wykonuje `go build -o "$STAGING/danaco-console.exe"
./server/cmd/danaco-console`, gdzie `$STAGING` to katalog przejściowy `danaco-wydanie` w katalogu
   tymczasowym systemu (`$TEMP`, a przy jego braku `/tmp`).
4. **Powłoka natywna.** Krok warunkowy. Skrypt pomija go, gdy ustawiono `DANACO_BEZ_POWLOKI=1` albo gdy
   `cargo` bądź `cargo tauri --version` nie są dostępne — w obu wypadkach wypisuje powód i podpowiada
   `cargo install tauri-cli`. W przeciwnym razie wykonuje w `budowa/desktop/src-tauri` polecenie
   `cargo tauri build --no-bundle` i oczekuje pliku
   `desktop/src-tauri/target/release/danaco-console-powloka.exe`. Brak tego pliku po zgłoszonym powodzeniu
   budowy jest traktowany jako błąd i przerywa wydanie.
5. **Pakowanie.** Wywołuje `pakowanie.sh` z katalogiem wyjściowym, binarką rdzenia ze `$STAGING`,
   katalogiem `client/dist` oraz — o ile powstała — ścieżką powłoki.

Skrypt pakowania realizuje trzy operacje i jedną regułę ochronną:

- kopiuje binarkę rdzenia jako `danaco-console.exe` w katalogu wyjściowym;
- odświeża `client/dist`: **usuwa z podkatalogu `client/` wszystkie pliki poza `*.md`** (nazwy plików
  pakietu zawierają sumy kontrolne, więc stare pliki muszą ustąpić), usuwa opróżnione katalogi, po czym
  kopiuje nowy pakiet;
- kopiuje powłokę jako `Danaco Console.exe`, o ile została wskazana.

Reguła ochronna brzmi: **skrypt nigdy nie usuwa ani nie nadpisuje plików `*.md` w katalogu wyjściowym**.
Katalog produktu zawiera równolegle dokumentację produktu (w tym niniejszy dokument), a ponowne wydanie
nie może jej skasować. Przed kopiowaniem skrypt sprawdza obecność binarki rdzenia oraz pliku
`index.html` w katalogu pakietu klienta i przerywa pracę z komunikatem, gdy któregoś brakuje.

Po zamknięciu wydania skrypt wypisuje podsumowanie zależne od wariantu: przy wariancie rdzeń+klient —
wskazanie na `danaco-console.exe` i adres `http://127.0.0.1:17870/`; przy wariancie pełnym — wskazanie na
`Danaco Console.exe` z adnotacją, że powłoka stawia rdzeń w tle.

### 5.4 Wariant rdzeń + klient (bez powłoki)

Wariant ten powstaje w dwóch sytuacjach: na żądanie (`DANACO_BEZ_POWLOKI=1`) albo samoczynnie, gdy stacja
nie ma łańcucha narzędzi Rust. Katalog produktu zawiera wtedy `danaco-console.exe` i `client/dist`,
lecz **nie** `Danaco Console.exe`.

Praktycznie oznacza to jedną różnicę: interfejs otwiera się w przeglądarce zamiast w oknie natywnym.
Rdzeń serwuje pakiet interfejsu obok gniazda transportu, więc adres `http://127.0.0.1:17870/` prowadzi
do kompletnego interfejsu. Cała funkcjonalność produktu jest po stronie rdzenia i klienta — powłoka
nie wnosi żadnej funkcji dziedzinowej, a jedynie okno systemowe, ikonę w zasobniku i automatyczne
podniesienie rdzenia.

Wariant rdzeń+klient jest zalecany na stacjach, gdzie powłoka nie jest potrzebna, oraz jako droga
najkrótsza do sprawdzenia, czy produkt w ogóle wstaje.

### 5.5 Droga ręczna — budowa części osobno

Poszczególne części można budować niezależnie — jest to przydatne w pracy nad jedną warstwą oraz przy
diagnostyce niepowodzeń wydania.

```bash
# Rdzeń — z katalogu budowa/
go build -o danaco-console.exe ./server/cmd/danaco-console

# Klient — z katalogu budowa/client/
npm ci --no-audit --no-fund      # tylko przy pierwszym uruchomieniu
npm run build                    # tsc --noEmit && vite build → client/dist

# Powłoka — z katalogu budowa/desktop/src-tauri/
cargo tauri build --no-bundle
```

Złożenie artefaktów w katalogu produktu można wykonać wprost skryptem pakowania, bez powtarzania budowy:

```bash
bash budowa/scripts/pakowanie.sh C:/DanacoConsole_App \
  budowa/danaco-console.exe \
  budowa/client/dist \
  budowa/desktop/src-tauri/target/release/danaco-console-powloka.exe
```

Ostatni argument jest opcjonalny; jego pominięcie daje wariant rdzeń+klient. Skrypt wymaga co najmniej
trzech argumentów i przy mniejszej liczbie wypisuje wzorzec użycia.

### 5.6 Instalator NSIS

Konfiguracja pakowania w `tauri.conf.json` deklaruje `"bundle": { "active": true, "targets": ["nsis"] }`
wraz z metadanymi producenta (`publisher`: Danaco Holding Group Sp. z o.o., `copyright`, `category`:
DeveloperTool) oraz zestawem ikon. Oznacza to, że projekt jest **przygotowany** do wytworzenia
instalatora Windows.

Instalatory NSIS składają na Linuksie dwa skrypty wydania: `budowa/scripts/instalka-hybryda-win-x64.sh`
(powłoka programu, ze wskazaniem rdzenia wdrożenia z rozstrzygnięcia 29 wpisanym przy kompilacji —
rozdział 12.6) i `budowa/scripts/instalka-kreatora-win-x64.sh` (kreator sześciu kroków). Oba mierzą
katalog Security gotowego pliku i przyjmują `DANACO_PODPIS=wymagany|pomijany`.

Instalator NSIS niesie **wyłącznie powłokę** z osadzonym interfejsem — i tak ma być: rdzeń stoi na
serwerze wdrożenia (rozstrzygnięcie 8), a instalator pełny obejmujący rdzeń nie powstaje. Podpis
Authenticode rozstrzyga rozstrzygnięcie 27 (rozdział 4.4): do zakupu certyfikatu OV od Certum oba
skrypty składają z jawnym `DANACO_PODPIS=pomijany`; po zakupie polecenie podpisujące wchodzi do
`bundle.windows` obu profili i składanie idzie z `DANACO_PODPIS=wymagany`.

---

## 6. Struktura instalacji i model katalogów

### 6.1 Rozdział źródeł od wyjścia

Model katalogów produktu opiera się na twardym rozdziale dwóch drzew:

| Ścieżka | Rola | Zawartość |
| --- | --- | --- |
| `C:\DanacoConsole_Git` | **źródła** | repozytorium: `budowa/`, `opracowania/`, dokumenty projektu |
| `C:\DanacoConsole_App` | **wyjście** | produkt gotowy do uruchomienia oraz dokumentacja produktu |

Rozdział ten nie jest konwencją redakcyjną, lecz obowiązuje w narzędziach: skrypt wydania domyślnie
składa produkt do `C:/DanacoConsole_App`, a skrypt pakowania jest jedynym miejscem, które do tego katalogu
zapisuje. Odwrotnie — proces budowy nie umieszcza artefaktów w katalogu źródeł poza katalogami roboczymi
kompilatorów (`client/node_modules`, `client/dist`, `desktop/src-tauri/target`).

Praktyczne następstwa dla Operatora:

- Katalog produktu można skasować w całości bez utraty czegokolwiek poza dokumentacją — pod warunkiem
  jej uprzedniego zabezpieczenia. Wydanie odtworzy pliki wykonywalne, ale nie odtworzy dokumentów.
- **Katalog produktu nie jest katalogiem danych.** Baza SQLite i cała trwałość leżą gdzie indziej
  (rozdział 6.3), więc usunięcie katalogu produktu nie kasuje historii rozmów, kont ani konfiguracji.
- Katalogu produktu nie należy trzymać wewnątrz repozytorium ani w katalogu synchronizowanym z chmurą —
  odświeżanie `client/dist` przy każdym wydaniu generuje intensywny ruch plikowy.

### 6.2 Zawartość katalogu produktu

Po pełnym wydaniu katalog `C:\DanacoConsole_App` ma następującą postać:

```
C:\DanacoConsole_App\
├── danaco-console.exe        ← rdzeń: rola all, nasłuch 17870, serwuje interfejs
├── Danaco Console.exe        ← powłoka natywna (tylko w wariancie pełnym)
├── client\
│   └── dist\                 ← pakiet interfejsu (index.html + zasoby z sumami w nazwach)
└── *.md                      ← dokumentacja produktu, nietykana przez wydanie
```

Układ ten nie jest dowolny — odpowiada dwóm miejscom wyszukiwania w kodzie powłoki:
`desktop/src-tauri/src/rdzen/lokalizacja.rs` szuka `danaco-console.exe` obok pliku powłoki, a
`desktop/src-tauri/src/rdzen/pakiet_klienta.rs` szuka katalogu `client/dist` obok binarki rdzenia.
Przeniesienie któregokolwiek pliku w inne miejsce zrywa to wyszukiwanie.

Domyślna wartość katalogu pakietu interfejsu po stronie rdzenia to ścieżka **względna** `client\dist`
(stałe `KatalogKlientaZrodlo = "client"` i `KatalogKlientaWydanie = "dist"`). Względność jest zamierzona:
katalog pakietu podąża za katalogiem roboczym procesu. Skutek praktyczny — rdzeń uruchomiony z innego
katalogu roboczego niż katalog produktu nie znajdzie interfejsu, dopóki nie wskaże się go przełącznikiem
`--klient` albo zmienną `DANACO_KATALOG_KLIENTA`. Brak pakietu **nie wstrzymuje startu**: gniazdo
transportu pracuje bez plików statycznych.

### 6.3 Katalog danych rdzenia

Katalog danych jest jedynym miejscem trwałości produktu. Rozstrzyganie jego ścieżki przebiega
następująco (`server/internal/konfiguracja/katalog_danych.go`):

1. Jeżeli ustawiona jest zmienna systemowa `LOCALAPPDATA` — katalogiem danych jest
   `%LOCALAPPDATA%\DanacoConsole`. Stała nazwy katalogu to `DanacoConsole`.
2. W przeciwnym razie — podkatalog `DanacoConsole` w katalogu konfiguracyjnym użytkownika zwracanym
   przez system.
3. W ostateczności — `.\DanacoConsole` względem katalogu roboczego procesu.

Trzystopniowość istnieje po to, by start rdzenia nigdy nie zależał od obecności jednej zmiennej.
Na typowej stacji Windows obowiązuje krok pierwszy, a katalogiem danych jest:

```
C:\Users\<Operator>\AppData\Local\DanacoConsole
```

Katalog zakładany jest przy starcie wraz z brakującymi katalogami nadrzędnymi; katalog już istniejący
nie jest modyfikowany. Ustawienie własnej ścieżki odbywa się przełącznikiem `--dane` albo zmienną
`DANACO_KATALOG_DANYCH`.

Wewnątrz katalogu danych leży **jeden plik bazy**:

```
C:\Users\<Operator>\AppData\Local\DanacoConsole\danaco-console.db
```

Nazwa `danaco-console.db` jest stałą w kodzie (`NazwaPlikuBazy`), a ścieżka powstaje przez złożenie
katalogu danych z tą nazwą. Nie istnieje osobny przełącznik wskazujący plik bazy — **zmiana katalogu
danych przenosi całą trwałość naraz**. Jest to świadoma decyzja projektowa, upraszczająca przenoszenie
instalacji i wykonywanie kopii zapasowych.

Sam rdzeń innych plików w tym katalogu nie zakłada (poza plikami towarzyszącymi bazie `-wal` i `-shm`,
opisanymi w rozdziale 9.2). **Przy pracy z powłoką natywną** w tym samym katalogu powstaje jednak
dodatkowo **plik dziennika rdzenia** `rdzen-powloki.log` — powłoka kieruje do niego wyjście
uruchamianego przez siebie procesu rdzenia oraz dopisuje własne wiersze startowe. Zawartość pliku,
sposób jego czytania i skutki braku rotacji opisuje rozdział 21.1.

### 6.4 Katalog profili kanału głównego

Katalog profili to miejsce, w którym leżą katalogi konfiguracyjne kont kanału głównego — czyli katalogi
przekazywane procesowi `claude` przez zmienną `CLAUDE_CONFIG_DIR`. Wskazuje się go przełącznikiem
`--profile` albo zmienną `DANACO_KATALOG_PROFILI`.

Zasada rozdziału jest tu istotna dla bezpieczeństwa: **rdzeń trzyma wyłącznie odwołanie do katalogu
profilu, nigdy samych poświadczeń**. Wartości uwierzytelniające pozostają na dysku, poza repozytorium
i poza bazą produktu. Baza zna ścieżkę, nie zna sekretu.

Wartość domyślna katalogu profili jest **pusta**, co jest stanem dopuszczalnym: pula kont startuje wtedy
pusta, a brak konta ujawnia się dopiero przy próbie rozmowy. Ma to bezpośredni związek z procedurą
z rozdziału 20 — świeża instalacja bez zasianego konta odmówi wykonania tury modelu.

### 6.5 Reguły higieny katalogu produktu

1. **Dokumentacja produktu (`*.md`) mieszka w katalogu produktu i jest chroniona przez skrypt
   pakowania.** Nie należy jej przenosić do `client/`, gdzie ochrona obejmuje wprawdzie `*.md`, ale
   katalog bywa czyszczony z pozostałych plików.
2. **Nie należy ręcznie modyfikować zawartości `client/dist`.** Każde wydanie usuwa i odtwarza ten
   katalog; zmiany zostaną utracone bez ostrzeżenia.
3. **Nie należy umieszczać katalogu danych wewnątrz katalogu produktu.** Nic tego nie blokuje
   technicznie, lecz odbiera własność opisaną w rozdziale 6.1 — bezpieczne kasowanie katalogu produktu.
4. **Katalog produktu nie jest miejscem na katalogi robocze sesji.** Katalogi robocze wskazuje Operator
   w konfiguracji okna i powinny leżeć tam, gdzie leżą właściwe dane pracy.

---

## 7. Pierwsze uruchomienie

### 7.1 Uruchomienie rdzenia

Rdzeń uruchamia się przez wykonanie `danaco-console.exe`. Kluczowa jest przy tym rola katalogu roboczego:
domyślna ścieżka pakietu interfejsu jest względna, więc rdzeń należy uruchamiać **z katalogu produktu**.

```powershell
Set-Location C:\DanacoConsole_App
.\danaco-console.exe
```

Przy starcie rdzeń wykonuje kolejno: ustalenie konfiguracji z trzech warstw (rozdział 8.1), utworzenie
katalogu danych, otwarcie pliku bazy wraz z zastosowaniem migracji (rozdział 9.2), odtworzenie sesji
i okien z bazy oraz otwarcie nasłuchu transportu na porcie 17870. Rdzeń wypisuje jednowierszowy zapis
konfiguracji w postaci `rola=<rola> port=<port> dane=<katalog>`.

Katalog przełączników i wykaz zmiennych środowiskowych można wypisać wbudowaną pomocą:

```powershell
.\danaco-console.exe --help
```

Pomoc drukuje nagłówek `danaco-console — rdzeń Danaco Console`, katalog przełączników wraz z wartościami
domyślnymi oraz wiersz z nazwami zmiennych środowiskowych — wykaz pochodzi z tej samej funkcji, która
zasila odczyt konfiguracji, więc nie może rozjechać się z implementacją.

### 7.2 Interfejs z pakietu serwowanego przez rdzeń

Rdzeń serwuje pakiet interfejsu obok gniazda transportu. Po starcie interfejs jest osiągalny pod:

```
http://127.0.0.1:17870/
```

Jest to droga zalecana w wariancie rdzeń+klient, ponieważ nie wymaga uruchamiania niczego poza rdzeniem
ani obecności Node.js na stacji pracy. Warunkiem jest obecność katalogu `client/dist` w katalogu
roboczym procesu — przy uruchomieniu z katalogu produktu warunek jest spełniony z automatu.

Gdy pakietu brakuje, rdzeń **nie przerywa pracy**: gniazdo transportu działa normalnie, a pod adresem
HTTP nie ma czego wyświetlić. Jest to zachowanie zamierzone (brak elementu nie blokuje startu),
lecz dla Operatora wygląda jak „pusta strona” — rozpoznanie tego przypadku opisuje rozdział 21.

### 7.3 Interfejs z serwera deweloperskiego Vite

W pracy nad interfejsem używa się serwera deweloperskiego Vite, uruchamianego z katalogu `budowa/client`:

```bash
npm run dev
```

Serwer nasłuchuje domyślnie pod `http://localhost:5173` — tę wartość potwierdza `devUrl`
w `tauri.conf.json`. Rdzeń musi wtedy pracować równolegle, we własnym procesie: Vite serwuje wyłącznie
interfejs, a cała logika pozostaje po stronie rdzenia na porcie 17870.

Droga ta jest przeznaczona do pracy deweloperskiej, nie do eksploatacji produktu. Do eksploatacji służy
pakiet serwowany przez rdzeń (rozdział 7.2) albo powłoka natywna (rozdział 7.4).

### 7.4 Powłoka natywna Tauri

Powłokę uruchamia się przez `Danaco Console.exe` w katalogu produktu. Powłoka podnosi rdzeń w tle
i otwiera okno interfejsu; dostarcza także ikonę w zasobniku systemowym (cecha `tray-icon` w `Cargo.toml`)
oraz okna dialogowe systemu (`tauri-plugin-dialog`).

**Stan faktyczny i luka `frontendDist`.** Konfiguracja Tauri wskazuje `"frontendDist": "../../client/dist"`
— ścieżkę względem `desktop/src-tauri`, prowadzącą do katalogu `budowa/client/dist`. Katalog ten
**nie istnieje w świeżo sklonowanym repozytorium**: powstaje dopiero po wykonaniu `npm run build`
w kliencie. Audyt odnotował ten stan jako lukę uniemożliwiającą wytworzenie pakietu produkcyjnego:
konfiguracja wskazywała nieistniejący katalog, a `beforeBuildCommand` nie było zdefiniowane, więc budowa
powłoki nie wywoływała budowy klienta.

W bieżącym stanie drzewa roboczego, po zamknięciu prac naprawczych, `tauri.conf.json` zawiera:

```json
"beforeBuildCommand": { "script": "npm run build", "cwd": "../../client" },
"beforeDevCommand":   { "script": "npm run dev",   "cwd": "../../client", "wait": false },
"devUrl": "http://localhost:5173",
"frontendDist": "../../client/dist"
```

Oznacza to, że `cargo tauri build` sam wywoła budowę pakietu klienta przed budową powłoki, a
`cargo tauri dev` sam podniesie serwer Vite. Skrypt wydania świadomie dubluje ten krok — komentarz w
`wydanie.sh` nazywa powtórkę zamierzoną i tanią. Pozostaje jednak ograniczenie wynikające z samej
konstrukcji ścieżki: **droga budowy powłoki jest związana z drzewem źródeł**. Zbudowanie powłoki
wymaga kompletnego repozytorium z działającym klientem, a nie jedynie katalogu produktu.

Druga, niezależna obserwacja audytu dotyczy adresu transportu: klient ustala adres rdzenia na podstawie
`location.hostname`. Rozwiązanie działa poprawnie, gdy interfejs jest serwowany przez rdzeń albo przez
Vite po HTTP. Przy interfejsie ładowanym z zasobów osadzonych w pakiecie powłoki (bez pochodzenia HTTP)
zgadywanie adresu nie ma z czego wywieść wartości. Szczegóły i skutki opisuje rozdział 8.5.
**[DO DECYZJI OPERATORA]** — czy wprowadzić jawne ustawienie adresu rdzenia po stronie klienta
— wstrzykiwane przez powłokę — zamiast wywodzenia go z adresu strony.

### 7.5 Czego oczekiwać po pierwszym starcie

Uczciwy opis stanu, który zastanie Operator po pierwszym uruchomieniu świeżej instalacji:

- Rdzeń wstaje, zakłada katalog danych, tworzy bazę i stosuje migracje. **[DZIAŁA]**
- Interfejs ładuje się i łączy z rdzeniem po WebSocket; praca przeżywa rozłączenie klienta. **[DZIAŁA]**
- Okna operacyjne, karty sesji, motyw wizualny i nawigacja rysują się poprawnie. **[DZIAŁA]**
- **Wysłanie wiadomości do modelu zakończy się niepowodzeniem.** Rejestr kanałów modelu jest pusty
  (żadna migracja nie zasiewa wiersza kanału), a pula kont jest pusta, więc kanał główny odmawia
  wykonania. **[NIEZINTEGROWANE]** Doprowadzenie instalacji do stanu roboczego opisuje **rozdział 20**
  i jest to krok **obowiązkowy**, nie opcjonalny.
- Pulpit Mission Control **rysuje się ze stanu rdzenia** (`session.list`, `window.list`, `channel.list`
  oraz zdarzenia `session.changed`, `window.changed`, `queue.changed`, `progress.changed`) — plik danych
  przykładowych, na którym pulpit rysował się przed pracami naprawczymi, został usunięty. **[DZIAŁA CZĘŚCIOWO]**
  Kolumna kolejek pozostaje pusta do pierwszego zdarzenia, bo kontrakt nie ma odczytu `queue.list`, a
  same kolejki nadal nie mają silnika wykonania (rozdział 18.2) — pulpit pokazuje to uczciwym stanem
  pustym, nie danymi zmyślonymi.
- Strona główna niesie strefę **„Sesje w tle”**, zasilaną tym samym `session.list` i zdarzeniami
  `session.changed`/`window.changed`; powrót do sesji trwającej na rdzeniu idzie przez `session.bind`,
  a następnie `session.focus`. Trzy postacie uczciwe (oczekiwanie, pusto, błąd) zastępują dane
  miejscowe. **[DZIAŁA]**
- Pozycje nawigacji odpowiadające modułom i środowiskom nie przełączają przestrzeni roboczej —
  interfejs pozostaje na jednym ekranie sceny okien. **[PRZEWIDZIANE KONCEPCJĄ, W v2.0
  NIEZAIMPLEMENTOWANE]**

---

## 8. Konfiguracja rdzenia i klienta

### 8.1 Warstwy konfiguracji

Rdzeń ustala komplet swoich ustawień jednorazowo, w chwili startu procesu, i nie zmienia ich później.
Funkcja `Wczytaj` (`server/internal/konfiguracja/wczytanie.go`) nakłada trzy warstwy w ściśle ustalonej
kolejności:

| Kolejność | Warstwa | Źródło |
| --- | --- | --- |
| 1 | wartości domyślne | `Domyslna` w `konfiguracja/ustawienia.go` |
| 2 | zmienne środowiska | `zastosujSrodowisko` w `konfiguracja/srodowisko.go` |
| 3 | argumenty wywołania | `zastosujArgumenty` w `konfiguracja/argumenty.go` |

Warstwa późniejsza wygrywa z wcześniejszą: **argument wiersza poleceń przebija zmienną środowiska,
a zmienna środowiska przebija wartość domyślną**. Mechanizm jest zrealizowany elegancko — wartości
domyślne zestawu przełączników pochodzą z konfiguracji już zmodyfikowanej przez środowisko, więc
przełącznik niepodany nie kasuje ustawienia ze zmiennej.

Dwie reguły uzupełniające, obie potwierdzone w kodzie:

- **Zmienna nieustawiona albo pusta nie zmienia niczego.** Warunek brzmi `if tekst := odczyt(...);
  tekst != ""`, więc ustawienie zmiennej na pusty łańcuch jest równoważne jej nieustawieniu. Nie da się
  „wyzerować” wartości domyślnej pustą zmienną.
- **Po nałożeniu wszystkich warstw następuje sprawdzenie spójności** (`konfiguracja/sprawdzenie.go`):
  rola musi należeć do katalogu ról, port musi mieścić się w zakresie 1–65535, katalog danych nie może
  być pusty. Naruszenie któregokolwiek warunku kończy proces komunikatem `konfiguracja: <przyczyna>` —
  rdzeń nie wstaje z niepoprawnym ustawieniem.

Danaco Console **nie posiada pliku konfiguracyjnego rdzenia**. Nie istnieje żaden `config.yaml`,
`ustawienia.json` ani równoważny plik czytany przy starcie: całość konfiguracji startowej opisują
pięć przełączników i pięć zmiennych środowiska. Konfiguracja dziedzinowa (kategorie, definicje, opcje
okna konfiguracji) mieszka w bazie, nie w plikach — opisuje ją rozdział 17.

### 8.2 Przełączniki wiersza poleceń

Katalog przełączników jest zamknięty i liczy pięć pozycji:

| Przełącznik | Wartość domyślna | Znaczenie |
| --- | --- | --- |
| `--role` | `all` | rola procesu: `hub`, `agent` albo `all` |
| `--port` | `17870` | port nasłuchu transportu WebSocket |
| `--dane` | `%LOCALAPPDATA%\DanacoConsole` (rozdz. 6.3) | katalog danych rdzenia |
| `--klient` | `client\dist` (ścieżka względna) | katalog pakietu interfejsu |
| `--profile` | wartość pusta | katalog profili kanału głównego |

Przełączniki obsługuje standardowy pakiet `flag` języka Go, więc dopuszczalne są obie postacie zapisu
(`--port 17870` i `--port=17870`, a także wariant z jednym myślnikiem). Przełącznik `--help` drukuje
pomoc i **kończy proces bez błędu** — `main.go` rozpoznaje `flag.ErrHelp` i wraca normalnie.

Przykład uruchomienia z jawnie wskazanym katalogiem danych i profili:

```powershell
.\danaco-console.exe --port 17870 `
  --dane "C:\Users\Operator\AppData\Local\DanacoConsole" `
  --profile "C:\Users\Operator\AppData\Local\DanacoConsole\profile"
```

### 8.3 Rola procesu

Rola rozstrzyga, **które tory pracy** uruchamia proces (`server/cmd/danaco-console/uruchomienie/
rozgalezienie.go`). Ta sama binarka obsługuje wszystkie trzy warianty:

| Rola | Tor interfejsu (gniazdo WebSocket) | Tor wykonawczy (wiersze na stdin/stdout) | Zajmuje port |
| --- | --- | --- | --- |
| `hub` | tak | nie | tak |
| `agent` | nie | tak | **nie** |
| `all` | tak — tor wiodący | tak — w tle | tak |

Tor wykonawczy przyjmuje żądania kontraktu **wierszami strumienia wejścia** i odsyła odpowiedzi
wierszami strumienia wyjścia; jeden wiersz to jedna koperta kontraktu. Z tego powodu dziennik rdzenia
idzie na wyjście diagnostyczne (`stderr`), a na wyjściu standardowym pojawia się wyłącznie kontrakt.

Trzy własności istotne przy wdrożeniu:

1. **Rola nierozpoznana nie zatrzymuje procesu.** Wartość spoza katalogu ról podana na etapie
   rozgałęzienia schodzi na zachowanie roli `all` z komunikatem w dzienniku. Uwaga: dotyczy to wyłącznie
   rozgałęzienia — sprawdzenie konfiguracji z rozdziału 8.1 odrzuci taką wartość wcześniej, więc
   w praktyce proces nie wstanie. Zachowanie zapasowe jest zabezpieczeniem konstrukcyjnym, nie drogą
   użytkową.
2. **Rola `all` uruchomiona bez podłączonego wejścia pracuje normalnie.** Wyczerpanie strumienia wejścia
   kończy wyłącznie tor wykonawczy; tor interfejsu pracuje dalej. To jest sytuacja typowa: rdzeń
   uruchomiony podwójnym kliknięciem albo podniesiony przez powłokę nie ma sensownego `stdin`.
3. **Rola nie zmienia niczego poza wyborem torów.** Nie różnicuje zakresu komend, nie ogranicza tabel
   ani nie zmienia miejsca uruchamiania procesów kanału. W szczególności rola `agent` **nie jest**
   mechanizmem pracy zdalnej — o tym mówi wprost rozdział 14.

Dla wdrożenia na stacji Operatora właściwą i domyślną wartością jest `all`. Ról `hub` i `agent`
dokument nie zaleca w bieżącej wersji: rozdzielenie procesów ma sens dopiero wraz z drogą wykonawczą
pracy zdalnej, której v2.0 nie posiada.

### 8.4 Port i adres gniazda

Portem domyślnym jest **17870** — stała `PortDomyslny` w `konfiguracja/ustawienia.go`. Zmienia się go
przełącznikiem `--port` albo zmienną `DANACO_PORT`.

**Ostrzeżenie o rozspojeniu.** Ta sama liczba jest zapisana po stronie klienta jako stała
`PORT_RDZENIA_LOKALNEGO` w `client/src/polaczenie/adres-rdzenia.ts`. Klient **nie pyta rdzenia o port** —
składa adres gniazda z własnej stałej. Zmiana portu rdzenia bez ponownego zbudowania klienta z
odpowiadającym ustawieniem powoduje, że interfejs nie trafi w rdzeń: strona się załaduje, lecz połączenie
WebSocket nie powstanie. Rozjazd tych dwóch wartości był odnotowany jako usterka.

Reguła praktyczna jest zatem jednoznaczna: **portu 17870 nie należy zmieniać**, dopóki nie zachodzi
konflikt z inną usługą. Jeżeli zmiana jest konieczna, trzeba wykonać obie czynności naraz:

1. uruchomić rdzeń z nowym portem (`--port` albo `DANACO_PORT`),
2. zbudować klienta ze zmienną budowania `VITE_ADRES_RDZENIA` wskazującą pełny adres gniazda:
   `VITE_ADRES_RDZENIA=ws://127.0.0.1:17999/ws npm run build`.

Ścieżka gniazda to `/ws` (stała `SCIEZKA_GNIAZDA`). Pełny adres domyślny ma postać
`ws://127.0.0.1:17870/ws`, a interfejs serwowany przez rdzeń — `http://127.0.0.1:17870/`.

**[DO DECYZJI OPERATORA]** — czy wprowadzić przekazywanie portu z rdzenia do klienta jedną z dwóch dróg:
punktem końcowym HTTP z konfiguracją albo wstrzyknięciem do dokumentu, by usunąć konieczność podwójnej zmiany.
Dokument opisuje stan zastany, nie postuluje rozwiązania.

### 8.5 Ustalanie adresu rdzenia po stronie klienta

Klient ustala adres gniazda w dwóch krokach (`adres-rdzenia.ts`):

1. **Ustawienie jawne.** Jeżeli przy budowie podano zmienną `VITE_ADRES_RDZENIA` i jest ona niepustym
   łańcuchem, klient używa jej dosłownie. Jest to jedyna droga wskazania rdzenia pod innym adresem niż
   strona serwująca interfejs.
2. **Wywiedzenie z lokalizacji dokumentu.** W przeciwnym razie klient bierze `location.protocol`
   i `location.hostname`, wybiera schemat `wss:` dla `https:` oraz `ws:` w pozostałych wypadkach, po czym
   składa adres z własną stałą portu i ścieżką `/ws`.

Zachowanie zapasowe jest zdefiniowane: gdy `location` nie istnieje albo protokół to `file:`, klient
wraca do `ws://127.0.0.1:17870/ws`. Podobnie pusty `hostname` schodzi na `127.0.0.1`.

Konsekwencje wdrożeniowe:

- Interfejs serwowany przez rdzeń pod `http://127.0.0.1:17870/` — adres wywodzi się poprawnie. **[DZIAŁA]**
- Interfejs z serwera deweloperskiego Vite pod `http://localhost:5173` — `hostname` to `localhost`,
  port podstawia stała, więc adres wychodzi `ws://localhost:17870/ws` i trafia w rdzeń. **[DZIAŁA]**
- Interfejs ładowany z zasobów osadzonych w pakiecie powłoki, bez pochodzenia HTTP — wywodzenie nie ma
  z czego wziąć wartości i sprawa opiera się o zachowanie zapasowe. Audyt odnotował to jako punkt kruchy
  drogi produkcyjnej powłoki. **[DZIAŁA CZĘŚCIOWO]** Zalecenie: budując powłokę do wydania, podawać
  `VITE_ADRES_RDZENIA` jawnie zamiast polegać na wywiedzeniu.

Poza adresem gniazda klient **nie posiada własnej konfiguracji**. Nie czyta żadnego pliku ustawień,
nie ma panelu ustawień połączenia i nie zapisuje trwale niczego poza tym, co przechowuje przeglądarka.
Cała trwałość produktu leży po stronie rdzenia — w bazie opisanej w rozdziale 9.

---

## 9. Konfiguracja bazy danych

### 9.1 Plik bazy i jego położenie

Trwałość Danaco Console mieści się w **jednym pliku SQLite**. Nie ma serwera bazy danych, nie ma
osobnej usługi do zainstalowania i nie ma poświadczeń do bazy — sterownik `modernc.org/sqlite` jest
wkompilowany w rdzeń, a plik otwierany jest bezpośrednio.

Ścieżka pliku powstaje przez złożenie katalogu danych ze stałą nazwą `danaco-console.db`
(`konfiguracja/plik_bazy.go`). Na typowej stacji Windows jest to:

```
C:\Users\<Operator>\AppData\Local\DanacoConsole\danaco-console.db
```

**Nie istnieje przełącznik ani zmienna wskazująca sam plik bazy.** Jedyną drogą przeniesienia trwałości
jest zmiana katalogu danych (`--dane` albo `DANACO_KATALOG_DANYCH`). Decyzja jest świadoma: jeden
katalog niesie całą trwałość, więc kopia zapasowa i przenosiny sprowadzają się do operacji na katalogu.

Otwarcie bazy (`store/baza.go`, funkcja `Otworz`) przebiega tak:

1. Pusta ścieżka jest odrzucana błędem `store: pusta ścieżka pliku bazy`.
2. Katalog nadrzędny pliku jest zakładany wraz z brakującymi katalogami pośrednimi. Baza powstaje więc
   sama — Operator nie tworzy jej ręcznie i nie uruchamia żadnego skryptu inicjującego.
3. Otwierana jest pula połączeń z ciągiem połączenia zawierającym pragmy (rozdział 9.3), po czym
   wykonywany jest `Ping`. Brak odpowiedzi kończy start błędem `store: baza %q nie odpowiada`.
4. Stosowane są wszystkie niezastosowane migracje (rozdział 9.2). Niepowodzenie migracji zamyka bazę
   i **przerywa start rdzenia**.

Plik bazy nie leży w katalogu produktu. Skasowanie `C:\DanacoConsole_App` nie kasuje danych; skasowanie
katalogu danych kasuje **wszystko**: sesje, wiadomości, okna, konta, kanały, punkty dostępu, ustawienia
i dokumenty tożsamości.

### 9.2 Migracje

Schemat bazy jest doprowadzany do bieżącej postaci **automatycznie przy każdym starcie rdzenia**.
Operator nie uruchamia migracji osobnym poleceniem — nie ma takiego polecenia i nie jest potrzebne.

Cechy mechanizmu, wszystkie potwierdzone w `store/migracje.go` i `store/zrodlo_migracji.go`:

- **Kroki są wkompilowane w binarkę** dyrektywą `go:embed migracja_*.sql`. Wdrożenie nie wymaga
  kopiowania plików `.sql` obok programu — katalog produktu nie zawiera i nie musi zawierać schematu.
- **Nazewnictwo kroku** ma postać `migracja_NNN_nazwa.sql`. Numer wersji i nazwa wyprowadzane są
  z nazwy pliku; plik o innej postaci nazwy zatrzymuje start błędem.
- **Rejestr zastosowanych kroków** to tabela `migracja` (kolumny `id`, `wersja`, `nazwa`,
  `suma_kontrolna`, `zastosowano`), z warunkiem `UNIQUE` na kolumnie `wersja`. Jest to jedyna tabela
  zakładana z kodu, a nie z pliku `.sql` — musi istnieć, zanim ruszy pierwszy krok.
- **Każdy krok wykonuje się w jednej transakcji razem z wpisem do rejestru.** Schemat nigdy nie zostaje
  zastosowany połowicznie: albo krok przechodzi w całości i jest odnotowany, albo transakcja jest
  wycofywana.
- **Suma kontrolna SHA-256 treści kroku** jest zapisywana w rejestrze i sprawdzana przy każdym
  kolejnym starcie. Zmiana treści migracji już zastosowanej daje błąd
  `store: migracja NNN (nazwa) zmieniła treść po zastosowaniu` i **nie pozwala rdzeniowi wstać**.
  Reguła praktyczna: migracji wydanej nie wolno edytować — poprawka to nowy krok o nowym numerze.
- **Kroki stosowane są rosnąco po numerze wersji.** Krok już odnotowany jest pomijany.

Bieżący wykaz obejmuje **trzynaście kroków** o numerach 001–009 oraz 012–015:

| Krok | Plik | Zakres |
| --- | --- | --- |
| 001 | `migracja_001_fundament.sql` | fundament schematu |
| 002 | `migracja_002_okna.sql` | okna |
| 003 | `migracja_003_kolejki.sql` | kolejki |
| 004 | `migracja_004_pamiec.sql` | pamięć |
| 005 | `migracja_005_okna_operacyjne.sql` | okna operacyjne |
| 006 | `migracja_006_identyfikatory_zewnetrzne.sql` | identyfikatory zewnętrzne |
| 007 | `migracja_007_zaczyn_slownikow.sql` | zaczyn słowników |
| 008 | `migracja_008_katalog_akcji.sql` | katalog akcji |
| 009 | `migracja_009_zaczyn_akcji.sql` | zaczyn akcji |
| 012 | `migracja_012_katalog_ustawien.sql` | katalog ustawień okna konfiguracji |
| 013 | `migracja_013_punkty_dostepu.sql` | punkty dostępu i nadania |
| 014 | `migracja_014_katalog_kont.sql` | katalog kont |
| 015 | `migracja_015_tozsamosc_modelu.sql` | dokumenty tożsamości |

**Luka numeracji 010–011 nie jest usterką.** Numery 010 i 011 nie istnieją w wykazie i nigdy nie
powstały: zostały zarezerwowane przy pracy równoległej dla wykonawców, którzy ostatecznie ich nie użyli.
Mechanizm stosuje kroki rosnąco po numerze, więc wolny numer w środku ciągu nie zaburza niczego.
Osobny sprawdzian (`store/numeracja_migracji_test.go`) pilnuje **jednoznaczności i porządku**
numerów, świadomie nie wymagając ich ciągłości. Operator, który zobaczy w tabeli `migracja` wersje
1–9 oraz 12–15, widzi stan poprawny — nie brakujące migracje.

Bieżącą wersję schematu zwraca `SELECT MAX(wersja) FROM migracja`; dla świeżej, poprawnie zmigrowanej
bazy jest to **15**.

### 9.3 Pragmy i tryb dziennika

Pragmy SQLite trafiają do **ciągu połączenia**, nie do zapytania wykonywanego po otwarciu. Jest to
rozstrzygnięcie istotne: pula `database/sql` otwiera połączenia wielokrotnie i w różnym czasie, więc
pragma ustawiona pojedynczym zapytaniem obowiązywałaby tylko jedno z nich. Wykaz z `store/baza.go`:

| Pragma | Wartość | Skutek |
| --- | --- | --- |
| `foreign_keys` | `1` | więzy kluczy obcych są egzekwowane — łańcuch sesja → karta → okno → wiadomość jest pilnowany przez bazę |
| `busy_timeout` | `5000` (ms) | zajęta baza jest ponawiana przez 5 sekund, zanim padnie błąd blokady |
| `journal_mode` | `WAL` | dziennik z wyprzedzeniem zapisu: odczyt nie blokuje zapisu ani odwrotnie |
| `synchronous` | `NORMAL` | równowaga między trwałością zapisu a wydajnością, typowa dla trybu WAL |

Tryb **WAL** ma bezpośrednie następstwo dla kopii zapasowych i przenoszenia instalacji: obok pliku
`danaco-console.db` pojawiają się pliki towarzyszące `danaco-console.db-wal` (dziennik) oraz
`danaco-console.db-shm` (pamięć wspólna). Kopiowanie samego pliku głównego przy pracującym rdzeniu może
dać kopię niekompletną.

Rdzeń udostępnia własną kontrolę spójności (`store/spojnosc.go`): `PRAGMA integrity_check` z oczekiwanym
wynikiem `ok` oraz `PRAGMA foreign_key_check` zwracającą listę naruszeń. Kontrola jest wykorzystywana
w sprawdzianach i dostępna diagnostycznie; **nie jest** uruchamiana samoczynnie przy każdym starcie.

### 9.4 Kopia zapasowa i przeniesienie

Procedura kopii zapasowej jest krótka, bo trwałość jest skupiona w jednym miejscu:

1. **Zatrzymać rdzeń** (zamknąć powłokę albo przerwać proces `danaco-console.exe`). Krok jest istotny
   właśnie z powodu trybu WAL — zamknięcie bazy scala dziennik z plikiem głównym.
2. Skopiować **cały katalog danych**, nie sam plik `.db`. Kopiowanie przy zatrzymanym rdzeniu obejmuje
   komplet: plik główny, — jeśli pozostały — pliki `-wal` i `-shm`, a przy pracy z powłoką natywną
   także plik dziennika rdzenia `rdzen-powloki.log` (rozdział 21.1). Dziennik nie jest daną
   produkcyjną i do odtworzenia pracy nie jest potrzebny, lecz kopia katalogu zabiera go razem
   z bazą — bywa to przydatne przy odtwarzaniu okoliczności usterki z dnia wykonania kopii.
3. Kopię przechowywać poza katalogiem produktu, bo katalog produktu jest odtwarzalny wydaniem, a dane
   nie są.

Przeniesienie instalacji na inną stację sprowadza się do: zbudowania produktu na stacji docelowej
(rozdział 5), skopiowania katalogu danych w to samo miejsce w profilu użytkownika albo wskazania go
przełącznikiem `--dane`, oraz **osobnego przeniesienia katalogów `CLAUDE_CONFIG_DIR` kont** — baza
przechowuje wyłącznie ścieżki do nich, nigdy samych poświadczeń (rozdział 11.2).

Przywrócenie kopii do rdzenia nowszej wersji jest bezpieczne w jedną stronę: migracje niezastosowane
zostaną dołożone przy starcie. **Droga powrotna nie istnieje** — mechanizm nie posiada migracji
odwrotnych, więc bazy zmigrowanej do wyższej wersji schematu nie da się użyć ze starszym rdzeniem.
**[DO DECYZJI OPERATORA]** — czy ustanowić cykliczną, automatyczną kopię katalogu danych (zadanie
harmonogramu Windows), oraz jaki okres przechowywania kopii przyjąć. Produkt nie posiada wbudowanego
mechanizmu kopii zapasowej.

---

## 10. Konfiguracja kanałów modelu

### 10.1 Rejestr kanałów sterowany danymi

**Kanał modelu** to droga, którą rdzeń rozmawia z modelem. W Danaco Console kanał nie jest typem
w kodzie, lecz **wierszem w tabeli `kanal_modelu`**. Nowy kanał powstaje przez dopisanie wiersza,
nie przez zmianę programu — jest to przyjęte rozstrzygnięcie architektoniczne, konsekwentnie zrealizowane.

Kształt wiersza (`migracja_001_fundament.sql`, odczyt w `models/zrodlo_bazy.go`):

| Kolumna | Typ | Uwagi |
| --- | --- | --- |
| `id` | INTEGER | klucz główny, nadawany automatycznie |
| `kod` | TEXT NOT NULL UNIQUE | kod kanału; **drugi klucz wyszukiwania** (rozdz. 10.1.1) |
| `nazwa` | TEXT NOT NULL | nazwa widoczna dla Operatora |
| `dostawca` | TEXT NOT NULL | nazwa dostawcy; przy zakładaniu przez kontrakt powtarza rodzaj |
| `identyfikator_modelu` | TEXT NOT NULL | identyfikator modelu przekazywany kanałowi |
| `rodzaj_kanalu` | TEXT NOT NULL | **`CHECK(rodzaj_kanalu IN ('cli','api','sdk','lokalny'))`** |
| `konto_id` | INTEGER | odwołanie do `konto(id)`, `ON DELETE SET NULL` |
| `poswiadczenie_odwolanie` | TEXT | **odwołanie** do danych dostępowych, nigdy sam sekret |
| `parametry_json` | TEXT NOT NULL DEFAULT `'{}'` | parametry kanału jako obiekt JSON |
| `multimodalny` | INTEGER 0/1 | znacznik informacyjny |
| `aktywny` | INTEGER 0/1, domyślnie 1 | wiersz nieczynny nie trafia do rejestru czynnych |
| `kolejnosc` | INTEGER, domyślnie 0 | porządek prezentacji |
| `utworzono` | TEXT | znacznik czasu ISO 8601 UTC, nadawany automatycznie |

Rdzeń buduje z tych wierszy **rejestr kanałów** (`models/rejestr.go`) przy starcie i po każdej zmianie
wywołanej komendą `channel.*`. Trzy własności rejestru są istotne wdrożeniowo:

- **Wiersz nieczynny (`aktywny = 0`) nie daje kanału** — jest w wykazie, lecz nie da się na niego
  wysłać zapytania.
- **Wiersz, dla którego nie ma fabryki adaptera, nie wywraca rejestru.** Trafia na listę pominiętych
  wraz z powodem (`brak fabryki adaptera <klucz>`), a pozostałe kanały pracują normalnie.
- **Wiersz niezmieniony zachowuje swoją instancję adaptera** przy odświeżeniu, więc dopisanie kanału
  nie przerywa strumienia biegnącego na innym kanale.

#### 10.1.1 Rozstrzyganie adaptera i kluczy wyszukiwania

Dwa mechanizmy rozstrzygania łatwo pomylić, a różnica jest praktyczna:

1. **Klucz adaptera** (`Definicja.KluczAdaptera`) wskazuje fabrykę, która zbuduje wykonawcę kanału.
   Rozstrzyga go parametr `"adapter"` z `parametry_json`, a **dopiero w jego braku** kolumna
   `rodzaj_kanalu`. Wbudowane fabryki stoją pod kluczami `echo` oraz `api`; klucz `cli` dokłada pakiet
   `internal/injection` przy montażu rdzenia.
2. **Klucze wyszukiwania kanału** (`kluczeKanalu`) — kanał jest osiągalny **pod dwiema wartościami
   naraz**: pod identyfikatorem wiersza (tekstowa postać kolumny `id`) oraz pod kolumną `kod`.
   Okno komunikacji przechowuje jedną z nich w polu `modelChannelId`.

Z połączenia obu reguł wynika **krytyczne ograniczenie kolumny `rodzaj_kanalu`**: warunek `CHECK`
dopuszcza wyłącznie wartości `cli`, `api`, `sdk` i `lokalny`. Wartość `echo` **nie jest dopuszczalna
jako rodzaj kanału** — próba zapisu `rodzaj_kanalu = 'echo'` zostanie odrzucona przez bazę. Kanał
echo zakłada się więc jako wiersz rodzaju `lokalny` (albo innego dopuszczalnego) z parametrem
`{"adapter":"echo"}`. Wzorzec ten jest utrwalony w sprawdzianach rdzenia i jest jedyną poprawną drogą.

Zauważmy też, że rodzaj `sdk` przechodzi warunek `CHECK`, ale **nie ma fabryki adaptera** — wiersz
takiego rodzaju bez parametru `"adapter"` wyląduje na liście pominiętych. Rodzaj `sdk` jest dziś
wartością słownikową bez drogi wykonawczej. **[BRAK]**

### 10.2 Komendy kontraktu do zarządzania rejestrem

Rejestrem zarządzają cztery komendy kontraktu, obsługiwane przez rdzeń po WebSocket
(`core/adapter_kanaly.go`):

| Komenda | Treść żądania | Treść wyniku |
| --- | --- | --- |
| `channel.add` | `name` (wymagane), `kind` (wymagane), `model?`, `enabled?`, `config?` | `channel` |
| `channel.update` | `channelId` (wymagane), `name?`, `model?`, `enabled?`, `config?` | `channel` |
| `channel.remove` | `channelId` | `channelId` |
| `channel.list` | `enabledOnly?` | `channels` |

Zachowania warte odnotowania, wszystkie sprawdzone w kodzie:

- **`channel.add` nie przyjmuje kodu kanału.** Kod nadaje rdzeń automatycznie w postaci
  `kanal-<licznik>-<losowe>` (przedrostek `kanal-`). Operator nie ma przez kontrakt wpływu na wartość
  kolumny `kod`.
- **`channel.add` nie przyjmuje pola dostawcy.** Kolumna `dostawca` bierze wartość parametru
  `"provider"` z `config`, a w jego braku **powtarza `kind`**.
- **`channel.update` i `channel.remove` odnajdują wiersz po kolumnie `kod`**, nie po numerze wiersza
  (`PobierzPoKodzie`). Pole `channelId` w tych komendach oznacza więc kod.
- **`channel.list` zwraca w polu `id` tekstową postać kolumny `id`**, czyli numer wiersza — nie kod.
  Wynik `channel.add` zwraca natomiast w tym samym polu **kod**. **[DZIAŁA CZĘŚCIOWO]** Rozbieżność
  jest realna i ma następstwo praktyczne: identyfikatora wziętego z `channel.list` nie da się podać
  do `channel.update` ani `channel.remove`. Kod kanału odczytuje się wtedy wprost z bazy
  (`SELECT id, kod FROM kanal_modelu`). Dla samego **wysłania wiadomości rozbieżność jest nieszkodliwa**,
  bo rejestr rozpoznaje obie wartości (rozdział 10.1.1).
- Po każdej zmianie rejestr jest odświeżany, więc **dopisanie kanału działa natychmiast, bez restartu
  rdzenia**. Jest to istotny kontrast wobec puli kont, która restartu wymaga (rozdział 11.3).

**Brak zaczynu i brak interfejsu — stan faktyczny.** Żadna migracja nie zasiewa wiersza kanału:
**świeża instalacja ma rejestr kanałów pusty**. Jednocześnie interfejs **nie posiada widoku zakładania
kanałów** — klient wywołuje wyłącznie `channel.list` (moduły `sterowanie/rejestr-kanalow.ts`,
`modele/zrodlo-modeli.ts`, `mission-control/zrodlo-pulpitu.ts`). Jedyne miejsce wywołujące `channel.add`
to moduł `rozmowa/zapewnienie-kanalu.ts`, podpięty do stanowiska podglądu rozmowy, a nie do
produkcyjnej ścieżki okna komunikacji. **[NIEZINTEGROWANE]** Doprowadzenie instalacji do stanu
roboczego opisuje rozdział 20 i wymaga działania Operatora.

### 10.3 Kanał `cli` — Claude Code

Kanał główny. Wykonawcę dostarcza pakiet `internal/injection`, wnoszony do rejestru pod kluczem
danych `cli`. Adapter uruchamia zewnętrzny program `claude` w trybie nieinteraktywnym ze strumieniowym
wyjściem JSON i przekłada jego strumień na fragmenty odpowiedzi.

Do pracy kanał wymaga **łącznie trzech rzeczy**:

1. wiersza w `kanal_modelu` o rodzaju `cli` (albo z parametrem `{"adapter":"cli"}`),
2. dostępnego programu `claude` (rozdział 3.5),
3. **co najmniej jednego konta rodzaju `cli` w puli kont** (rozdział 11).

Brak któregokolwiek z tych elementów kończy turę błędem, przy czym komunikat o pustej puli kont bywa
mylący — opisuje to rozdział 21.

Kanał ten jest tym, do którego domyślnie kieruje interfejs: pole `kanalModelu` opisu okna przyjmuje
wartość `KnownChannelKinds[0]`, czyli literał **`'cli'`** (`client/src/okno-komunikacji/opis-okna.ts`).
Jest to **rodzaj kanału, nie identyfikator wiersza**. Rejestr rozpozna tę wartość tylko wtedy, gdy
istnieje wiersz o `kod = 'cli'`. Fakt ten jest podstawą procedury z rozdziału 20 i wprost tłumaczy,
dlaczego wiersz zakładany komendą `channel.add` (z kodem `kanal-…`) nie wystarcza domyślnemu ustawieniu
okna. Wartość domyślną da się nadpisać przy budowie klienta zmienną `VITE_KANAL_MODELU`.

### 10.4 Kanał `api`

Generyczny kanał sieciowy (`models/adapter_api.go`), sparametryzowany wyłącznie wierszem rejestru —
**bez zaszytego dostawcy**. Szczegóły parametrów opisuje rozdział 16.

Ograniczenie zasadnicze: adapter `api` **nie jest zintegrowany z pulą kont**. Konta rodzaju `api`
i `sdk` można założyć w bazie i widać je w interfejsie, lecz nie prowadzi od nich żadna droga
wykonawcza do kanału `api` — poświadczenie kanał bierze z własnych parametrów wiersza, nie z konta.
**[NIEZINTEGROWANE]**

### 10.5 Kanał `echo`

Kanał próbny bez sieci i bez procesu zewnętrznego (`models/adapter_echo.go`). Odsyła treść zapytania
porcjami, poprzedzoną prowenancją wywołania — dokładnie tak, jak zrobiłby to kanał rzeczywisty.
Jest to **pełna, zamknięta droga**: rejestr → zapytanie → prowenancja → strumień → ujście → zapis
do bazy, bez zależności od dostawcy, konta i łącza. **[DZIAŁA]**

Dwa parametry wiersza sterują jego zachowaniem:

| Parametr | Znaczenie | Wartość domyślna |
| --- | --- | --- |
| `przedrostek` | tekst doklejany przed odpowiedzią | brak (pusty) |
| `porcja` | liczba znaków w jednym fragmencie | `24` |

Wartość `porcja` nieczytelna albo niedodatnia schodzi na wartość domyślną, zamiast zatrzymywać kanał.

Kanał echo jest **zalecanym narzędziem kontroli instalacji**: pozwala potwierdzić, że transport,
rejestr, okno, strumień i zapis wiadomości działają, zanim do gry wejdzie program `claude`
i poświadczenia. Wykorzystuje to krok kontrolny z rozdziału 20.6. Przypomnienie z rozdziału 10.1.1:
wiersz zakłada się z rodzajem `lokalny` i parametrem `{"adapter":"echo"}`.

### 10.6 Kanał SSH — stan faktyczny

**Kanał modelu rodzaju SSH nie istnieje.** Nie ma go w warunku `CHECK` kolumny `rodzaj_kanalu`, nie ma
fabryki adaptera pod takim kluczem, nie ma pakietu realizującego zdalne wykonanie tury. **[BRAK]**

Program `ssh` występuje w produkcie w **jednym miejscu i w innej roli**: jako program uruchamiany przez
**most MCP** (rozdział 15), gdy Operator skonfiguruje punkt dostępu wskazujący na most zdalny. Jest to
narzędzie udostępniane modelowi, a nie droga wykonania tury modelu. Mylenie tych dwóch rzeczy prowadzi
do fałszywego wniosku, że produkt potrafi wykonać turę na maszynie zdalnej — nie potrafi
(rozdział 14.2).

### 10.7 Kanał HTTP — stan faktyczny

Nie istnieje osobny kanał rodzaju „HTTP”. Rolę kanału sieciowego pełni w całości adapter `api`
(rozdziały 10.4 i 16), który komunikuje się po HTTP z punktem końcowym wskazanym w parametrach wiersza.
Dokument nie wymienia „kanału HTTP” jako odrębnego bytu, bo takiego bytu w produkcie nie ma.

---

## 11. Konta, pula rotacji i profile tożsamości

### 11.1 Tabela `konto`

**Konto** jest profilem uwierzytelnienia kanału. Tabela `konto` powstała w migracji 001 i została
rozszerzona migracją 014 do pełnego zakresu kontraktu. Kształt wiersza po obu migracjach:

| Kolumna | Pochodzenie | Uwagi |
| --- | --- | --- |
| `id` | 001 | klucz główny |
| `nazwa` | 001 | `NOT NULL UNIQUE`; **to ona jest kodem konta w puli rotacji** |
| `rodzaj` | 001 | `CHECK(rodzaj IN ('cli','api','sdk'))` |
| `identyfikator_zewnetrzny` | 001 | identyfikator u dostawcy |
| `poswiadczenie_odwolanie` | 001 | **odwołanie** do wpisu w magazynie sekretów albo ścieżki profilu |
| `katalog_konfiguracji` | 001 | katalog przekazywany procesowi jako `CLAUDE_CONFIG_DIR` |
| `aktywne` | 001 | 0/1, domyślnie 1 — przełącznik Operatora „konto w ogóle istnieje w pracy” |
| `kolejnosc` | 001 | porządek rotacji |
| `utworzono`, `zaktualizowano` | 001 | znaczniki czasu |
| `dostawca` | 014 | `NOT NULL DEFAULT ''` |
| `model_domyslny` | 014 | identyfikator modelu domyślnego konta |
| `adres_bazowy` | 014 | adres punktu końcowego dostawcy |
| `domyslne` | 014 | 0/1; indeks częściowy `idx_konto_domyslne_rodzaj` wymusza **dokładnie jedno konto domyślne na rodzaj** |
| `stan` | 014 | `CHECK(stan IN ('aktywne','wyczerpane','zawieszone'))`, domyślnie `aktywne` |
| `wyczerpane_do` | 014 | znacznik odnowienia limitu |

Rozdział `aktywne` od `stan` jest zamierzony: pierwsza kolumna odpowiada na pytanie Operatora
(„czy konto jest w użyciu”), druga na pytanie rotacji („czy konto jest teraz zdatne”). Kolumny `stan`
i `wyczerpane_do` są **dziennikiem tego, co pula już rozpoznała**, a nie drugą implementacją rotacji.

**Do pracy kanału głównego liczy się jeden rodzaj: `cli`.** Konta rodzaju `api` i `sdk` można założyć
i widać je w oknie modeli, lecz **nie prowadzi od nich żadna droga wykonawcza** — adapter `api` bierze
poświadczenie z parametrów wiersza kanału, nie z konta. **[NIEZINTEGROWANE]**

Wiązanie konta z kanałem prowadzi wyłącznie kolumna `kanal_modelu.konto_id` z regułą
`ON DELETE SET NULL`: usunięcie konta **odłącza** kanały, lecz ich nie kasuje. Kontrakt oddaje to
polem `detachedChannelIds` w wyniku `account.remove`.

Konta obsługuje rodzina komend `account.*` (pełny odczyt i zapis, z widokiem w oknie modeli/kont).
**[DZIAŁA]** — z zastrzeżeniami rozdziału 11.3.

### 11.2 `CLAUDE_CONFIG_DIR` per konto

Rdzeń **nie przechowuje poświadczeń**. Cała tożsamość uwierzytelnienia kanału głównego przenoszona jest
jedną zmienną środowiskową procesu potomnego: `CLAUDE_CONFIG_DIR` wskazuje katalog konfiguracji konta
na dysku. Baza zna ścieżkę; sekret leży w katalogu.

Sposób składania środowiska procesu (`injection/proces.go`, funkcja `srodowisko`) ma trzy warstwy
o rosnącym pierwszeństwie:

1. środowisko dziedziczone po procesie rdzenia (`os.Environ`),
2. zmienne dołożone z ustawień wywołania,
3. **`CLAUDE_CONFIG_DIR` z konta — nadpisuje każdą wcześniejszą wartość.**

Wniosek praktyczny: ustawienie `CLAUDE_CONFIG_DIR` w środowisku systemowym stacji **nie zmieni** konta,
którym pracuje kanał, gdy pula podała konto z katalogiem. Wartość z konta zawsze wygrywa. Konto bez
wypełnionej kolumny `katalog_konfiguracji` **w ogóle nie wchodzi do puli** — funkcja `naPuleKont`
pomija takie wiersze.

Katalog konfiguracji konta zakłada i uwierzytelnia Operator, posługując się programem `claude`; produkt
nie ma do tego własnego mechanizmu i nie podejmuje próby uwierzytelnienia w imieniu Operatora.
**[DO DECYZJI OPERATORA]** — gdzie umieścić katalogi kont (proponowana konwencja: podkatalogi wspólnego
katalogu profili wskazanego przełącznikiem `--profile`) oraz jak zabezpieczyć je uprawnieniami systemu
plików.

### 11.3 Pula rotacji i wyczerpanie limitu

Pula kont (`injection/pula_kont.go`) trzyma kolejność kont i pamięć o wyczerpanych limitach. Zasada
działania jest prosta: wyczerpanie limitu jednego konta nie kończy rozmowy — pula podaje następne,
a zmiana konta jest zdarzeniem kanału, nie końcem sesji.

**Skąd biorą się konta w puli** (`core/montaz_zrodla.go`, funkcja `pulaKont`):

1. **Źródło pierwsze — tabela `konto`.** Zapytanie `KontaRotacji` wybiera wiersze warunkiem
   `rodzaj = 'cli' AND aktywne = 1 AND stan <> 'zawieszone'`, porządkiem
   `ORDER BY domyslne DESC, kolejnosc, id`. Konto domyślne idzie pierwsze. Do puli trafiają tylko te
   wiersze, które mają wypełniony `katalog_konfiguracji`; **kodem konta w puli jest kolumna `nazwa`**.
2. **Źródło zapasowe — katalog profili z dysku.** Używane wtedy i tylko wtedy, gdy krok pierwszy nie
   dał ani jednego konta, a przełącznik `--profile` (albo `DANACO_KATALOG_PROFILI`) wskazuje katalog.
   Każdy **podkatalog** staje się kontem: kodem jest nazwa podkatalogu, katalogiem konfiguracji —
   jego ścieżka. Podkatalogi zaczynające się od kropki albo podkreślenia są pomijane jako zaplecze
   narzędzi. Konta sortowane są alfabetycznie po nazwie.
3. Gdy oba źródła są puste, pula jest pusta i rdzeń rusza normalnie. Brak konta ujawnia się dopiero
   przy próbie rozmowy.

**Rozpoznanie wyczerpania limitu** (`injection/wyczerpanie.go`) ma trzy stopnie wiarygodności: pole
stanu limitu podane przez program `claude` (stan `rejected` wraz z chwilą odnowienia), treść
podsumowania tury, wreszcie wyjście diagnostyczne procesu. Dwa ostatnie rozpoznawane są wzorcami
tekstowymi: `usage limit reached`, `rate limit exceeded`, `rate_limit_error`, `quota exceeded`, `429`.
Wyczerpanie bez podanej chwili odnowienia oznacza konto jako niedostępne **na godzinę**.

**Dwa ograniczenia, które trzeba znać przy wdrożeniu:**

- **Pula budowana jest raz, przy starcie rdzenia.** Nie ma odświeżenia puli po zmianie tabeli `konto`
— kontrast wobec rejestru kanałów, który odświeża się po każdej komendzie `channel.*`. Dodanie,
  usunięcie, zawieszenie albo zmiana kolejności konta **wymaga ponownego uruchomienia rdzenia**, żeby
  odniosło skutek w wykonaniu. **[DZIAŁA CZĘŚCIOWO]** Punkt ten jest wpleciony w procedurę
  z rozdziału 20 i jest najczęstszym źródłem zdziwienia („po dodaniu konta system nadal odmawia”).
- **Stan wyczerpania nie jest utrwalany.** Pamięć o wyczerpanych limitach żyje w pamięci procesu;
  restart rdzenia ją kasuje i pula zaczyna od konta pierwszego. Kolumny `stan` i `wyczerpane_do`
  istnieją w tabeli, lecz pula ich nie zapisuje. **[DZIAŁA CZĘŚCIOWO]**

Osobna, ważna uwaga diagnostyczna: **przy zerowej liczbie kont komunikat o wyczerpanym limicie jest
mylący.** Pula pusta i pula w całości wyczerpana kończą się tym samym brakiem konta do podania.
Rozróżnienie opisuje rozdział 21.

### 11.4 Katalog profili a katalog konfiguracji konta

Dwa pojęcia o podobnych nazwach, o zupełnie różnej roli:

| Pojęcie | Skąd | Rola |
| --- | --- | --- |
| **katalog profili** | `--profile` / `DANACO_KATALOG_PROFILI` | katalog nadrzędny; jego **podkatalogi** stają się kontami w źródle zapasowym puli |
| **katalog konfiguracji konta** | kolumna `konto.katalog_konfiguracji` | katalog jednego konta, przekazywany procesowi jako `CLAUDE_CONFIG_DIR` |

Oba mogą współistnieć: typową i zalecaną konwencją jest umieszczenie katalogów kont jako podkatalogów
katalogu profili i wpisanie ich pełnych ścieżek do kolumny `katalog_konfiguracji`. Tabela `konto`
pozostaje wtedy źródłem pierwszym, a katalog profili — czytelnym miejscem na dysku.

Wartość domyślna katalogu profili jest **pusta** i jest to stan dopuszczalny (rozdział 6.4): przy pustej
tabeli `konto` i pustym katalogu profili pula startuje pusta.

### 11.5 Nakładka tożsamości

**Nakładka tożsamości** to prompt systemowy składany z warstw i podawany programowi `claude`
(`injection/nakladka.go`). Mechanizm działa od końca do końca. **[DZIAŁA]**

Trzy warstwy, w stałej kolejności krytyczności:

| Kolejność | Warstwa | Wartość w kontrakcie |
| --- | --- | --- |
| 1 | konstytucja | `constitution` |
| 2 | profil | `profile` |
| 3 | ekspertyza zadaniowa | `expertise` |

Warstwy o nazwach spoza katalogu układają się za nimi, w kolejności podania; sortowanie jest stabilne.
Warstwa o pustej treści jest pomijana. Warstwy sklejane są pustym wierszem.

Dwa tryby podania rozstrzygają, **czy nakładka zastępuje prompt programu, czy się do niego dokłada**:

| Tryb (kontrakt) | Stała silnika | Przełącznik programu | Skutek |
| --- | --- | --- | --- |
| `ZASTAP` | `TrybZastap` | `--system-prompt` | prompt programu jest podmieniany w całości |
| `DOLACZ` | `TrybDopisz` | `--append-system-prompt` | nakładka dokłada się; narzędzia i zachowanie powłoki zostają |

Tryb pusty znaczy dołączenie. Nakładka pusta nie daje żadnego przełącznika — kanał rusza z samą
powłoką programu. Skrót SHA-256 złożonego promptu trafia do prowenancji wywołania i służy **wyłącznie**
diagnostyce: niczego nie dopuszcza i niczego nie blokuje.

Treść warstw pochodzi w całości z konfiguracji — pakiet nie zawiera ani jednego zdania promptu.
Nośnikiem treści są tabele założone migracją 015 i rodzina komend kontraktu:

| Komenda | Rola |
| --- | --- |
| `identity.category.list` | katalog kategorii zasad |
| `identity.document.get` | odczyt dokumentu tożsamości |
| `identity.document.set` | zapis dokumentu tożsamości |
| `identity.document.remove` | usunięcie dokumentu |
| `identity.effective.get` | nakładka obowiązująca: tryb i warstwy po złożeniu |

Dokumenty tożsamości zapisywane są **per oś**: `platform` (dla całej platformy), `model` (dla
identyfikatora modelu) albo `account` (dla `konto.id`). Oś jest prostopadła do poziomu zasięgu
konfiguracji. Trybem domyślnym jest **`ZASTAP`** — wartość kolumny `tryb_domyslny` odpowiada wartości
domyślnej klucza ustawienia `tozsamosc.tryb_domyslny`.

Kategoria zasad jest **wierszem katalogu**, nie stałą w kodzie: dołożenie kategorii to dopisanie
wiersza. Zmiana kategorii i dokumentów **nie wymaga restartu rdzenia** — inaczej niż zmiana puli kont.

---

## 12. Zmienne środowiskowe

### 12.1 Wykaz zupełny

Rdzeń czyta **dokładnie pięć** zmiennych środowiskowych. Wykaz jest zamknięty i prowadzony w jednym
miejscu — funkcja `ZmienneSrodowiska` w `server/internal/konfiguracja/srodowisko.go`. Ta sama funkcja
zasila pomoc wiersza poleceń i sprawdzian zgodności ze wzorcem `budowa/.env.example`, więc trzy rzeczy
— kod, pomoc i wzorzec — nie mogą się rozjechać. Osobny sprawdzian zabrania odczytu zmiennych po nazwie
dosłownej poza tym pakietem (zabezpieczenie przed nawrotem tej usterki).

Wszystkie nazwy noszą przedrostek `DANACO_`:

| Zmienna | Odpowiednik w wierszu poleceń | Wartość domyślna |
| --- | --- | --- |
| `DANACO_ROLA` | `--role` | `all` |
| `DANACO_PORT` | `--port` | `17870` |
| `DANACO_KATALOG_DANYCH` | `--dane` | `%LOCALAPPDATA%\DanacoConsole` |
| `DANACO_KATALOG_KLIENTA` | `--klient` | `client\dist` (ścieżka względna) |
| `DANACO_KATALOG_PROFILI` | `--profile` | wartość pusta |

Powłoka natywna czyta dodatkowo pięć zmiennych **na własny użytek**
(`desktop/src-tauri/src/ustawienia.rs`, `rdzen/pakiet_klienta.rs` oraz `rdzen/dziennik.rs`) — nie są to
zmienne rdzenia, lecz zmienne procesu powłoki. Dwie z nich (`DANACO_PORT` i `DANACO_KATALOG_DANYCH`)
noszą nazwy wspólne z rdzeniem: powłoka je **wyłącznie czyta**, aby trafić w to samo miejsce, w które
trafia rdzeń, i nigdy nie rozstrzyga ich za rdzeń.

| Zmienna | Rola w powłoce | Zachowanie przy braku |
| --- | --- | --- |
| `DANACO_PORT` | port, pod którym powłoka szuka rdzenia (`http://127.0.0.1:<port>`) | `17870`; wartość niebędąca liczbą schodzi na domyślną, nie przerywa startu |
| `DANACO_RDZEN` | jawnie wskazana ścieżka binarki rdzenia | wyszukanie w miejscach znanych powłoce |
| `DANACO_ADRES_INTERFEJSU` | jawnie wskazany adres interfejsu (serwer deweloperski albo rdzeń) | rozstrzygnięcie wewnętrzne powłoki |
| `DANACO_KATALOG_KLIENTA` | katalog pakietu przekazywany rdzeniowi uruchamianemu w tle | wyszukanie `client/dist` obok binarki rdzenia |
| `DANACO_KATALOG_DANYCH` | katalog, w którym powłoka zakłada plik dziennika rdzenia `rdzen-powloki.log` (rozdział 21.1) | `%LOCALAPPDATA%\DanacoConsole`; poza Windows `$HOME/DanacoConsole`; gdy brak i tej zmiennej — `./DanacoConsole` względem katalogu roboczego |

Powłoka traktuje wartość pustą (także złożoną z samych spacji) jak brak ustawienia.

### 12.2 Opis pojedynczych zmiennych

**`DANACO_ROLA`** — rola procesu, jedna z wartości `hub`, `agent`, `all`. Rozstrzyga, które tory pracy
uruchamia proces (rozdział 8.3). Wartość spoza katalogu ról **zatrzymuje start** komunikatem
`DANACO_ROLA: nieznana rola "…"; dopuszczalne: hub|agent|all`. Wartością właściwą dla stacji Operatora
jest `all`.

**`DANACO_PORT`** — port nasłuchu transportu WebSocket. Wartość musi być liczbą; tekst nieliczbowy
zatrzymuje start komunikatem `DANACO_PORT: wartość "…" nie jest liczbą`. Wartość poza zakresem 1–65535
zatrzymuje start na sprawdzeniu spójności. **Ostrzeżenie z rozdziału 8.4 obowiązuje: zmiana tej
zmiennej bez odpowiadającej zmiany po stronie klienta rozspaja interfejs z rdzeniem.**

**`DANACO_KATALOG_DANYCH`** — katalog danych rdzenia, w którym leży plik bazy `danaco-console.db`
(rozdział 9.1), a przy pracy z powłoką natywną także plik dziennika rdzenia `rdzen-powloki.log`
(rozdział 21.1) — powłoka czyta tę samą zmienną, aby założyć dziennik w tym samym miejscu, w którym
rdzeń trzyma bazę. Wartość pusta oznacza wartość domyślną rozstrzyganą trzystopniowo (rozdział 6.3).
Katalog wskazany tą zmienną jest zakładany przy starcie wraz z katalogami nadrzędnymi. **Zmiana tej
zmiennej przenosi całą trwałość produktu** — sesje, wiadomości, konta, kanały i ustawienia.

**`DANACO_KATALOG_KLIENTA`** — katalog pakietu interfejsu serwowanego obok gniazda. Wartość pusta
oznacza ścieżkę **względną** `client\dist`, liczoną od katalogu roboczego procesu. Katalog nieistniejący
**nie wstrzymuje nasłuchu** — gniazdo pracuje bez plików statycznych, a pod adresem HTTP nie ma czego
wyświetlić. Zmienna jest właściwym narzędziem wtedy, gdy rdzeń uruchamiany jest z katalogu roboczego
innego niż katalog produktu.

**`DANACO_KATALOG_PROFILI`** — katalog profili kanału głównego. Jego **podkatalogi** stają się kontami
w zapasowym źródle puli rotacji (rozdział 11.3). Zmienna trzyma katalogi profili, **nie poświadczenia**:
rdzeń zna wyłącznie odwołanie, a wartości uwierzytelniające zostają na dysku, poza repozytorium i poza
bazą. Wartość pusta jest stanem dopuszczalnym — pula kont startuje wtedy pusta, a brak konta zgłasza
dopiero próba rozmowy.

Reguła wspólna dla wszystkich pięciu: **zmienna nieustawiona albo pusta pozostawia wartość
dotychczasową** (rozdział 8.1). Nie istnieje wartość „wyłącz ustawienie”.

### 12.3 Zmienne świadomie nieobecne

Wzorzec `budowa/.env.example` prowadzi listę ustawień, których **celowo nie ma**, wraz z powodem.
Reguła brzmi: deklarowana jest wyłącznie zmienna, którą rdzeń rzeczywiście czyta — zmienna
zadeklarowana bez odczytu byłaby atrapą (obowiązuje zakaz atrap). Wpis wraca do wzorca razem z kodem, który go
odczyta, nie wcześniej.

| Ustawienie nieobecne | Powód |
| --- | --- |
| adres rdzenia dla agenta | tor wykonawczy roli `agent` pracuje na strumieniach procesu, więc adresu nie czyta; wpis wróci po powstaniu transportu sieciowego agenta |
| **ścieżka programu `claude`** | nie ma jej we wzorcu; ścieżkę bierze wiersz rejestru kanału albo `PATH` (rozdział 3.5) |
| klucze kanałów API | z założenia nie tutaj: rejestr kanałów buduje się z tabeli `kanal_modelu` w czasie działania, a odwołania do kluczy wpisuje Operator |
| faza wdrożenia | po powstaniu uwierzytelniania |
| poziom dziennika | po powstaniu dziennika sterowanego ustawieniem |

Wniosek dla Operatora: **nie ma zmiennej środowiskowej wskazującej program `claude`** i nie należy jej
szukać ani wymyślać. Skuteczne są dokładnie dwie drogi: wpis w parametrach wiersza rejestru kanału albo
obecność `claude` w `PATH`. Klucz katalogu ustawień `harness.program_claude` istnieje, lecz **nie jest
odczytywany przez wykonanie** (rozdział 17.2). **[NIEZINTEGROWANE]**

### 12.4 Zmienne przekazywane procesowi kanału

Proces programu `claude` dostaje środowisko składane z trzech warstw o rosnącym pierwszeństwie
(`injection/proces.go`):

1. **środowisko dziedziczone** po procesie rdzenia — cały `os.Environ`,
2. **zmienne dołożone z ustawień wywołania** (pole `Srodowisko` ustawień kanału),
3. **`CLAUDE_CONFIG_DIR` z konta puli** — wartość nadpisująca wszystko powyżej.

Skutki praktyczne, warte zapamiętania przy diagnostyce:

- Zmienne systemowe stacji (w tym `PATH`, potrzebny do odnalezienia programu `claude`) **są dziedziczone**
  przez proces kanału. Zmiana `PATH` po starcie rdzenia nie odniesie skutku — rdzeń dziedziczy
  środowisko z chwili własnego startu.
- **`CLAUDE_CONFIG_DIR` ustawiony ręcznie w systemie zostanie nadpisany** wartością z konta. Jest to
  zachowanie zamierzone: to konto, a nie środowisko stacji, rozstrzyga tożsamość uwierzytelnienia.
- Konto bez katalogu konfiguracji nie wchodzi do puli, więc nie istnieje przypadek „konto podane,
  a `CLAUDE_CONFIG_DIR` niepodany”.

### 12.5 Plik `.env` i zmienne budowania interfejsu

**Rdzeń nie czyta pliku `.env`.** W kodzie Go nie ma żadnego mechanizmu wczytywania takiego pliku —
odczyt zmiennych idzie wyłącznie przez `os.Getenv`. Plik `budowa/.env.example` jest **wzorcem
dokumentacyjnym**: opisuje zmienne, ich znaczenie i wartości domyślne. Utworzony z niego
`budowa/.env` pozostaje poza kontrolą wersji i **nie jest wczytywany samoczynnie** — Operator musi
ustawić zmienne w środowisku procesu (sesja PowerShell, skrót uruchomieniowy, definicja usługi albo
przedrostek polecenia w Git Bash).

```powershell
# Ustawienie na czas jednej sesji PowerShell
$env:DANACO_KATALOG_PROFILI = "C:\Users\Operator\AppData\Local\DanacoConsole\profile"
.\danaco-console.exe
```

```bash
# Ustawienie na czas jednego wywołania w Git Bash
DANACO_PORT=17870 ./danaco-console.exe
```

**[DO DECYZJI OPERATORA]** — czy wprowadzić do rdzenia wczytywanie pliku `.env`, czy pozostać przy
ustawianiu zmiennych w środowisku procesu i przekazywać je skrótem uruchomieniowym. Dokument opisuje
stan zastany: plik jest wzorcem, nie źródłem konfiguracji.

**Zmienne budowania interfejsu (`VITE_*`) to osobna rodzina.** Nie należą do `budowa/.env.example`
(Vite czyta je z katalogu `client/`), nie są czytane w czasie pracy i **wkompilowują się w pakiet
w chwili budowy**. Zmiana którejkolwiek z nich wymaga ponownego `npm run build`. Pełny wykaz zmiennych
odczytywanych przez klienta:

| Zmienna budowania | Rola | Wartość domyślna |
| --- | --- | --- |
| `VITE_ADRES_RDZENIA` | pełny adres gniazda WebSocket rdzenia | wywiedzenie z adresu strony (rozdz. 8.5) |
| `VITE_PROJEKT` | nazwa projektu w opisie okna | `Danaco Console` |
| `VITE_TYTUL_OKNA` | tytuł okna komunikacji | `Okno komunikacji` |
| `VITE_MODUL` | moduł początkowy okna | `talkin` (pierwsza pozycja `KnownModuleIds`) |
| `VITE_KANAL_MODELU` | kanał modelu okna | `cli` (pierwsza pozycja `KnownChannelKinds`) |
| `VITE_KATALOGI_ROBOCZE` | katalogi robocze okna, rozdzielone średnikiem | lista pusta |

Zmienna `VITE_KANAL_MODELU` zasługuje na osobną uwagę: jest jedyną drogą zmiany domyślnego wskazania
kanału po stronie klienta bez ingerencji w bazę. Wartość domyślna `cli` jest **rodzajem kanału**, nie
identyfikatorem wiersza — konsekwencje opisują rozdziały 10.3 i 20.

### 12.6 Zmienne wdrożenia serwerowego i składania wydania

Rozdziały 12.1–12.2 opisują pięć zmiennych stacji Operatora. Wdrożenie serwerowe — pakiet `.deb`
z `packaging/drzewo/`, jednostka `danaco-console.service`, rozstrzygnięcie 29 rejestru
`prowadzenie/decyzje.md` — dokłada zmienne czytane przez rdzeń poza tą piątką oraz zmienne czytane
wyłącznie przy składaniu instalki. Pełny wykaz nazw prowadzi nadal `ZmienneSrodowiska` i wzorzec
`budowa/.env.example`; poniżej te, bez których wdrożenie nie stanie albo stanie niezabezpieczone.

**`DANACO_KLUCZ_SEJFU`** — ścieżka pliku klucza pieczętującego sejf poświadczeń (hasła IMAP/SMTP,
klucze API; `server/internal/dane/sejf_poswiadczen.go`). Plik niesie 32 bajty albo 64 znaki
szesnastkowe; poza Windows rdzeń przyjmuje wyłącznie prawa `0400` i przy innych odmawia otwarcia sejfu
komunikatem `dane: sejf: klucz … ma prawa …, wymagane 0400`. Jednostka systemd wskazuje
`/etc/danaco-console/sejf.klucz`; skrypt `postinst` pakietu zakłada ten plik, gdy go nie ma (64 znaki
szesnastkowe z `/dev/urandom`, właściciel `danaco-console`, prawa `0400`), a plik istniejący zostawia —
nowy klucz odciąłby sejf zapieczętowany poprzednim. Wartość pusta = klucz własny `sejf.klucz` zakładany
w katalogu danych przy pierwszym użyciu sejfu (rozstrzygnięcia 26 i 30). Na serwerze plik ma leżeć poza
katalogiem danych: kopia bazy ani kopia katalogu danych nie otwiera wtedy sejfu. Utrata pliku jest
utratą dostępu do zapieczętowanych poświadczeń, więc kopia zapasowa (rozdział 9.4) obejmuje go osobno.

**`DANACO_SEKRET_NAWIAZANIA`** — sekret, którym klient przedstawia się przed uaktualnieniem połączenia
do WebSocket: parametrem `sekret` w adresie gniazda (klient składa go w `adresNawiazania`,
`klient/src/polaczenie/adres-rdzenia.ts`, wartość bierze z powłoki) albo nagłówkiem `X-Danaco-Sekret`.
Rdzeń czyta zmienną wprost ze środowiska procesu (warstwa transportu, nie nastawy bazy) i porównuje
czasem stałym; niezgodność odrzuca nawiązanie przed otwarciem gniazda. Puste = bez sprawdzenia.
Miejsce wpisu na serwerze: `/etc/danaco-console/srodowisko` (plik `EnvironmentFile` jednostki, prawa
`0640`, `root:danaco-console`). Tę samą wartość musi nieść powłoka na urządzeniu Operatora
(`ZMIENNA_SEKRET_NAWIAZANIA` w `desktop/src-tauri/src/ustawienia.rs`) — w instalkę sekretu się nie
wpisuje, bo napis wkompilowany w plik wykonywalny czyta `strings`. Bez odpowiednika w wierszu poleceń.

**Zmienne składania wydania** — czytane wyłącznie przy kompilacji powłoki (`option_env!`
w `desktop/src-tauri/src/ustawienia.rs`) i wpisywane w binarium instalki jako napis; po budowie nie ma
ich jak dopisać, a kreator instalacji o adres rdzenia nie pyta (rozstrzygnięcie 8). Skrypty
`scripts/instalka-hybryda-win-x64.sh` i `scripts/instalka-hybryda-win-arm.sh` odmawiają budowy bez
którejkolwiek z nich i po budowie sprawdzają, że host wszedł do binarium. Wartości wdrożenia
z rozstrzygnięcia 29:

| Zmienna | Wartość wdrożenia | Znaczenie |
| --- | --- | --- |
| `DANACO_HOST_WDROZENIA` | `console.danaco-group.pl` | nazwa rdzenia, do którego łączy się powłoka |
| `DANACO_PORT_WDROZENIA` | `443` | port na świat (Caddy, TLS z ACME); rdzeń za nim nasłuchuje na `127.0.0.1:17870` |
| `DANACO_SCHEMAT_WDROZENIA` | `https` | schemat łącza; klient wywodzi z niego `wss:`. Przyjmowane wyłącznie `http` albo `https`; bez `https` hasło Operatora szłoby otwartym tekstem |

```bash
DANACO_HOST_WDROZENIA=console.danaco-group.pl \
DANACO_PORT_WDROZENIA=443 \
DANACO_SCHEMAT_WDROZENIA=https \
  scripts/instalka-hybryda-win-x64.sh
```

TLS kończy się na Caddy, więc rdzeń wdrożenia nie potrzebuje `DANACO_TLS_CERTYFIKAT` ani
`DANACO_TLS_KLUCZ`. Adres podglądu budowy `51.75.62.180:80` nigdy nie jest celem wydania. Na
urządzeniu Operatora wskazanie z instalki nadpisują zmienne powłoki `DANACO_HOST_RDZENIA`
i `DANACO_SCHEMAT_RDZENIA`.

---

## 13. Konfiguracja katalogów roboczych i nadań

Rozdział opisuje trzy różne byty, których nie wolno mylić:

| Byt | Odpowiada na pytanie | Nośnik |
| --- | --- | --- |
| **katalog roboczy sesji** | gdzie model zostawia własne pliki | ustawienia `katalog.roboczy.*` |
| **katalogi dodatkowe okna** | do jakich katalogów model ma sięgać w czasie tury | lista katalogów okna → `--add-dir` |
| **punkty dostępu i nadania** | do jakich maszyn i zasobów model ma wgląd | tabele `punkt_dostepu`, `nadanie_dostepu` |

Rozdzielenie jest rozstrzygnięciem projektowym, zapisanym wprost w kodzie: moduł katalogu roboczego
**nie zna punktów dostępu i nigdy o nie nie pyta**. Katalog roboczy nie jest dostępem.

### 13.1 Katalog roboczy sesji

Katalog roboczy sesji jest miejscem, w którym powstają katalogi sesyjne i pliki robocze modelu.
Nie ma własnej tabeli — jest **pozycją katalogu ustawień**, rozstrzyganą przez rezolwer ośmiu poziomów
zasięgu, dwoma kluczami:

| Klucz ustawienia | Znaczenie | Wartość domyślna |
| --- | --- | --- |
| `katalog.roboczy.podstawa` | katalog, w którym powstają katalogi sesyjne | pusty — znaczy **miejsce instalacji aplikacji głównej** |
| `katalog.roboczy.wzorzec_sesji` | wzorzec nazwy katalogu jednej sesji, liczony względem podstawy | `sesje/<identyfikator>` |

**Miejsce instalacji** ustalane jest w chwili startu procesu (`KatalogInstalacji`): jest to katalog pliku
wykonywalnego rdzenia, po rozwiązaniu dowiązań symbolicznych. Gdy ścieżki pliku wykonywalnego nie da
się ustalić, przyjmowany jest katalog bieżący procesu, a w ostateczności zapis względny `.`. Żadna
z tych ścieżek nie kończy się błędem.

Dla instalacji zbudowanej wg rozdziału 6 domyślnym katalogiem sesji jest zatem:

```
C:\DanacoConsole_App\sesje\<identyfikator-sesji>\
```

**Składanie ścieżki** (funkcja `SciezkaSesji`) rządzi się czterema regułami:

1. Znacznik `<identyfikator>` we wzorcu jest miejscem podstawienia identyfikatora sesji.
2. Wzorzec **bez znacznika** dostaje identyfikator na końcu: wzorzec `pliki` daje
   `<podstawa>/pliki/<identyfikator>`.
3. Wzorzec, który wyprowadzałby **poza podstawę** — ścieżka bezwzględna albo wiodące `..` — jest
   odrzucany na rzecz wzorca domyślnego. Ustawienie Operatora steruje układem katalogów **wewnątrz**
   podstawy, nie omija samej podstawy.
4. Identyfikator sesji jest oczyszczany do nazwy bezpiecznej dla systemu plików: litery, cyfry, kreska,
   podkreślenie i kropka zostają, wszystko pozostałe staje się kreską. Identyfikator, z którego nie
   zostaje ani jeden znak, dostaje nazwę zastępczą `sesja`.

**Degradacja przy braku prawa zapisu.** Katalog żądany, którego nie da się przygotować do zapisu, nie
zatrzymuje startu sesji: sesja dostaje katalog zastępczy, a Operator — wiadomość o tym, że pracuje
gdzie indziej, niż ustawił. Struktura degradacji niesie ścieżkę żądaną, ścieżkę faktycznie użytą,
przyczynę oraz znacznik skuteczności. Wartość fałszywa tego ostatniego znaczy, że **nie ma gdzie
pisać** — sesja i tak startuje, a zapis pliku roboczego zawiedzie dopiero przy próbie.

Katalog roboczy sesji trafia do wywołania programu `claude` jako **katalog startowy procesu**
(`Ustawienia.KatalogRoboczy`). Gdy ustalenia brak, katalogiem startowym zostaje **pierwszy katalog
roboczy okna**, a gdy i tego nie ma — katalog bieżący procesu rdzenia. **[DZIAŁA]**

### 13.2 Katalogi dodatkowe okna

Okno komunikacji niesie **listę katalogów roboczych** (`Ustawienia.Katalogi`). Każdy z nich jedzie do
programu `claude` **osobnym przełącznikiem `--add-dir`**. Lista pusta oznacza brak przełączników,
nie błąd.

Wskazuje się je w konfiguracji okna; przy budowie klienta wartość początkową można podać zmienną
`VITE_KATALOGI_ROBOCZE` jako listę rozdzieloną średnikiem (rozdział 12.5).

Zasada higieny z rozdziału 6.5 obowiązuje: katalogi robocze **nie powinny leżeć w katalogu produktu**.
Właściwym miejscem jest tam, gdzie leżą rzeczywiste dane pracy. Katalog produktu jest odtwarzany przy
każdym wydaniu.

Przekazanie katalogów do wywołania jest drogą zamkniętą i sprawdzoną. **[DZIAŁA]**

### 13.3 Punkty dostępu i nadania

**Punkt dostępu** mówi, do czego model ma wgląd. **Nadanie** wiąże punkt z konkretnym oknem
komunikacji. Oba byty założyła migracja 013.

Tabela `punkt_dostepu` — kolumny istotne przy konfiguracji:

| Kolumna | Uwagi |
| --- | --- |
| `kod` | `NOT NULL UNIQUE`, nadawany przez rdzeń (przedrostek punktu dostępu) |
| `rodzaj` | `CHECK(rodzaj IN ('mcpBridge','localDirectory'))` |
| `urzadzenie_id` | odwołanie do `urzadzenie(id)`, `ON DELETE CASCADE` |
| `host`, `port`, `uzytkownik`, `sciezka_klucza` | adres mostu; `port` domyślnie `22` |
| `polecenie_startu`, `nazwa_mostu` | parametry mostu MCP (rozdział 15) |
| `poswiadczenie_odwolanie` | **odwołanie** do wpisu w magazynie sekretów, nigdy treść |
| `tryb_domyslny` | `CHECK(tryb_domyslny IN ('read','write'))`, domyślnie `read` |
| `stan` | `CHECK(stan IN ('unknown','reachable','unreachable'))`, domyślnie `unknown` |
| `aktywny`, `kolejnosc`, `utworzono`, `zaktualizowano` | jak w pozostałych katalogach |

Dwa warunki tabelowe pilnują sensowności wiersza:

```sql
CHECK(rodzaj <> 'localDirectory' OR urzadzenie_id IS NOT NULL)
CHECK(rodzaj <> 'mcpBridge'      OR host <> '')
```

Katalogi udostępniane przez punkt zapisywane są w osobnej tabeli `korzen_punktu_dostepu`; nadania mają
własną tabelę korzeni `korzen_nadania`, a argumenty trybów mostu — `argument_trybu_mostu`.

Tabela `nadanie_dostepu` wiąże okno z punktem: `okno_komunikacji_id`, `punkt_dostepu_id`, `tryb`
(`read`/`write`, domyślnie `read`), `kolejnosc`, `glowne`, `aktywne`. Obie relacje mają regułę
`ON DELETE CASCADE`, więc usunięcie okna albo punktu kasuje nadanie.

Rodzina komend kontraktu:

| Komenda | Rola |
| --- | --- |
| `access.point.add` / `.update` / `.remove` | zarządzanie punktami dostępu |
| `access.point.list` | wykaz punktów |
| `access.point.check` | sprawdzenie osiągalności punktu |
| `access.grant.add` / `.update` / `.remove` | zarządzanie nadaniami okna |
| `access.grant.list` | wykaz nadań |

Zapis tych rejestrów jest **zastrzeżony Operatorowi**: model ma do nich dostęp wyłącznie odczytowy
(`access.point.list`, `access.point.check`, `access.grant.list`). Model nie rozszerza własnego dostępu.

**Ograniczenie punktu rodzaju `localDirectory` — stan faktyczny.** Wiersz tego rodzaju wymaga
niepustego `urzadzenie_id`, a więc istniejącego wiersza w tabeli `urzadzenie`. Tabela istnieje od
migracji 001, lecz **produkt nie ma dla niej repozytorium ani komend kontraktu** — nie istnieje żadna
komenda zakładająca urządzenie. Punktu dostępu rodzaju `localDirectory` **nie da się więc założyć drogą
kontraktową**; próba zapisu bez wskazania urządzenia zostanie odrzucona przez warunek tabelowy.
**[NIEZINTEGROWANE]** Punktem zdatnym do konfiguracji w bieżącej wersji jest wyłącznie rodzaj
`mcpBridge` (rozdział 15).

**[DO DECYZJI OPERATORA]** — czy udostępnianie katalogów lokalnych prowadzić przez most MCP
o odpowiednich korzeniach, czy poczekać na drogę kontraktową dla urządzeń. Do czasu rozstrzygnięcia
katalogi lokalne wskazuje się drogą z rozdziału 13.2 (`--add-dir`), która działa i jest niezależna
od punktów dostępu.

---

## 14. Środowiska lokalne i zdalne

Rozdział wymaga szczególnej ostrożności redakcyjnej, bo nazewnictwo produktu sugeruje możliwości,
których wersja v2.0 nie posiada. Poniższy opis trzyma się wyłącznie tego, co da się wskazać w kodzie.

### 14.1 Rola `hub` i rola `agent`

Rola procesu (rozdział 8.3) jest jedynym mechanizmem rozdzielenia rdzenia na dwa procesy. Rola `hub`
prowadzi wyłącznie tor interfejsu — nasłuch transportu na porcie rdzenia. Rola `agent` prowadzi
wyłącznie tor wykonawczy — żądania kontraktu wierszami wejścia standardowego, odpowiedzi wierszami
wyjścia; **żaden port nie jest zajmowany**.

Konstrukcja jest przygotowaniem pod pracę rozproszoną, lecz **przygotowaniem niedokończonym**: agent
komunikuje się przez strumienie własnego procesu, a nie przez sieć. Wzorzec `.env.example` mówi to
wprost, wymieniając „adres rdzenia dla agenta” wśród ustawień **świadomie nieobecnych**, z powodem:
tor wykonawczy roli `agent` pracuje na strumieniach procesu, więc adresu nie czyta. Wpis wróci po
powstaniu transportu sieciowego agenta.

Wniosek wdrożeniowy: **rola `agent` nie jest dziś drogą pracy na maszynie zdalnej**. Uruchomienie
rdzenia z rolą `agent` na innej maszynie da proces, z którym nie ma jak się połączyć siecią.
**[NIEZINTEGROWANE]**

### 14.2 Zasięg wykonania — stan faktyczny

Okno komunikacji niesie pole **zasięgu wykonania** (`executionEnv` w kontrakcie, kolumna
`okno_komunikacji.srodowisko_wykonania` w bazie) o trzech wartościach:

| Wartość kontraktu | Wartość w bazie | Zamysł koncepcyjny |
| --- | --- | --- |
| `local` | `lokalne` | wykonanie na urządzeniu Operatora |
| `core` | `rdzen` | wykonanie na maszynie rdzenia |
| `remote` | `zdalny` | wykonanie na wskazanym hoście zdalnym |

**Stan faktyczny: zasięg wykonania jest dziś wyłącznie etykietą.** Wartość jest przyjmowana przez
kontrakt, zapisywana w bazie, przekładana w obie strony i **przenoszona do prowenancji wywołania** jako
pole informacyjne `executionEnv`. Na tym jej rola się kończy. **Nie istnieje ani jedna gałąź kodu,
która na podstawie tej wartości uruchomiłaby proces w innym miejscu.** Proces programu `claude` startuje
zawsze na maszynie rdzenia, niezależnie od ustawienia okna. **[NIEZINTEGROWANE]**

Trzy następstwa dla Operatora:

1. Ustawienie oknu zasięgu `remote` **nie przeniesie wykonania** nigdzie. Zmieni wyłącznie etykietę
   widoczną w prowenancji.
2. Kanał modelu rodzaju SSH **nie istnieje** (rozdział 10.6), więc nie ma nawet drogi transportowej,
   którą wykonanie mogłoby wyjść poza maszynę rdzenia.
3. **Host wykonania nie jest podłączony.** Nazwa konkretnego hosta nie jest wartością wyliczenia —
   miała być ustawieniem poziomu okna. Klient zapisuje ją komendą `config.set` pod kluczem
   **`window.executionHost`** (`client/src/sterowanie/klucze-ustawien.ts`). Serwer **nie zna tego
   klucza**: nie ma go w katalogu ustawień założonym migracją 012 ani w żadnej stałej rdzenia, a jego
   nazwa nie odpowiada kolumnie `host_wykonania` z opisu koncepcyjnego. Zapis trafia więc do tabeli
   ustawień i tam pozostaje, bez odbiorcy. **[NIEZINTEGROWANE]**

Ten sam rozjazd dotyczy dwóch pokrewnych kluczy zapisywanych przez klienta na poziomie okna:
`window.fallbackChannelId` (kanał zapasowy) oraz `window.reasoningEffort` (nakład rozumowania).
Pełny wykaz rozjazdów zawiera rozdział 17.2.

> **Zapis kluczy nastaw.** Nazwy w postaci `obszar.grupa.nastawa` użyte w tym rozdziale są
> kluczami konfiguracji, nie komendami kontraktu — konwencję zapisu wiąże
> [Standard redakcyjny i językowy](STANDARD-REDAKCYJNY-I-JEZYKOWY.md) rozdz. 7.3.

### 14.3 Praca zdalna dostępna dzisiaj

Aby zamknąć rozdział uczciwie: **jedyną formą sięgania poza maszynę rdzenia, która w v2.0 działa, jest
most MCP** (rozdział 15). Most jest narzędziem udostępnianym modelowi — pozwala modelowi wykonywać
czynności na maszynie zdalnej za pośrednictwem programu `ssh` uruchamianego przez punkt dostępu.
Nie jest to wykonanie tury modelu na maszynie zdalnej: sama tura, czyli proces programu `claude`,
biegnie przy rdzeniu.

Różnica ma znaczenie praktyczne przy planowaniu wdrożenia:

| Potrzeba | Droga dostępna w v2.0 |
| --- | --- |
| model ma czytać i zmieniać pliki na maszynie zdalnej | most MCP (rozdział 15) — **[DZIAŁA CZĘŚCIOWO]**, z ograniczeniami stałych |
| tura modelu ma biec na maszynie zdalnej (odciążenie stacji) | **brak drogi** — **[BRAK]** |
| rdzeń ma pracować na serwerze, a interfejs na stacji | droga wdrożenia z rozstrzygnięcia 29: rdzeń na `127.0.0.1:17870` z wymogiem logowania za Caddy — patrz rozdz. 2.4 i 12.6 |

**[DO DECYZJI OPERATORA]** — czy i kiedy wprowadzić transport sieciowy agenta wraz z uwierzytelnianiem,
bo dopiero to zamyka pion pracy rozproszonej. Do tego czasu dokument zaleca wdrożenie jednomaszynowe
z rolą `all`.

---

## 15. Konfiguracja narzędzi — mosty MCP

### 15.1 Model mostu

**Most MCP** jest narzędziem udostępnianym modelowi: pozwala mu sięgnąć do zasobów maszyny wskazanej
punktem dostępu. Danaco Console nie realizuje protokołu MCP samodzielnie — składa **konfigurację
`mcpServers`** i przekazuje ją programowi `claude` przełącznikiem `--mcp-config`. Właściwa rozmowa
protokołu toczy się między programem `claude` a mostem po drugiej stronie.

Rozdział ról jest ścisły i warto go zapamiętać:

- **punkt dostępu** mówi, **gdzie** stoi maszyna (adres, konto, port, odwołanie do klucza),
- **nadanie okna** mówi, **w jakim trybie** i **do których korzeni** sięga to jedno okno.

Generator wpisów (`core/most_mcp.go`) jest w całości zbiorem funkcji czystych: dane wchodzą, tekst
JSON wychodzi. Nie ma tam ani jednego połączenia sieciowego, ani jednego wywołania `ssh` — kształt
konfiguracji da się sprawdzić bez maszyny po drugiej stronie.

### 15.2 Punkt dostępu rodzaju `mcpBridge`

Most konfiguruje się jako punkt dostępu rodzaju **`mcpBridge`** (rozdział 13.3). Warunek tabelowy
wymaga dla tego rodzaju **niepustej kolumny `host`**. Pola wykorzystywane przez generator:

| Pole punktu (kontrakt) | Kolumna bazy | Rola |
| --- | --- | --- |
| `endpoint` | `host`/`port`/`uzytkownik` po rozłożeniu | adres w zapisie `[ssh://][konto@]host[:port]` — **ma pierwszeństwo** |
| `host` | `host` | nazwa maszyny, gdy `endpoint` jej nie podał |
| `credentialRef` | `poswiadczenie_odwolanie` | **odwołanie** do klucza; wchodzi jako `-i <klucz>` |
| `bridgeName` | `nazwa_mostu` | nazwa mostu |
| `defaultMode` | `tryb_domyslny` | `read` albo `write` |

Rozkład adresu (`core/most_mcp_adres.go`) obsługuje przedrostek `ssh://`, część `konto@` oraz `:port`,
przy czym za port uznawany jest wyłącznie zapis będący liczbą — dzięki temu adres IPv6 albo ścieżka
nie zostaną pomylone z portem. Wartości brakujące schodzą na domyślne mostu, nie na odmowę złożenia
wpisu.

Tryb nadania nie jest wpisywany w kod. Słowo, jakim dany most nazywa tryb — „odczyt”, „zapis” —
mieszka w tabeli **`argument_trybu_mostu`** (migracja 013) i dojeżdża do generatora jako argument
trybu. Most o innym słownictwie wymaga więc **nowego wiersza, nie zmiany rdzenia**. Argument pusty
oznacza uruchomienie samego skryptu — most `mcp-danaco-pulpit-console` wchodzi wtedy w tryb odczytu,
czyli w wariant bezpieczniejszy.

Do konfiguracji wchodzą wyłącznie **nadania i punkty czynne** oraz **punkty rodzaju `mcpBridge`**;
katalog lokalny nie jest mostem i nie ma czego uruchamiać przez `ssh`. Brak nadań daje pustą, lecz
poprawną mapę — brak dostępów nie jest błędem konfiguracji.

### 15.3 Przekazanie mostu do procesu kanału

Droga jest zamknięta i przebiega tak (`core/adapter_rozmowa_srodowisko.go`, `injection/argumenty.go`):

1. Dla okna rozmowy składany jest tekst konfiguracji `mcpServers` z jego nadań.
2. Tekst trafia do pola `KonfiguracjaMCP` zapytania kanału.
3. Kanał główny przekłada je na przełącznik **`--mcp-config`** — każda pozycja osobnym przełącznikiem.

Klucze mapy `mcpServers` są uporządkowane, więc **ten sam zbiór nadań zawsze daje ten sam tekst**
konfiguracji; klucz wpisu składa się z przedrostka `mcp-danaco-pulpit-console-` i identyfikatora
maszyny. Punkt, z którego identyfikatora nie zostaje ani jeden znak dopuszczalny, dostaje nazwę
zastępczą `maszyna`.

Brak nadań pozostawia proces **bez przełącznika `--mcp-config`** i rozmowa toczy się normalnie
(zasada „brak nie blokuje startu”). **[DZIAŁA]** — w zakresie składania i przekazania konfiguracji.

Kształt pojedynczego wpisu odpowiada dokładnie temu, czego oczekuje strona serwerowa mostu:

```json
{
  "mcpServers": {
    "mcp-danaco-pulpit-console-<maszyna>": {
      "type": "stdio",
      "command": "ssh",
      "args": ["-o", "BatchMode=yes", "-o", "StrictHostKeyChecking=accept-new",
               "-p", "22", "-i", "<odwołanie do klucza>",
               "<konto>@<host>", "/opt/danaco/most-konsoli/uruchom-stdio.sh <tryb>"],
      "env": { "DANACO_MOST_KORZENIE": "<korzeń1>:<korzeń2>" }
    }
  }
}
```

Argument `-i` pojawia się **wyłącznie wtedy, gdy punkt niesie odwołanie do klucza**; brak odwołania
zostawia rozstrzygnięcie konfiguracji `ssh`, zamiast wstawiać ścieżkę zmyśloną. Podobnie pole `env`
występuje tylko przy niepustej liście korzeni.

### 15.4 Stałe mostu i ich konsekwencje

Generator posługuje się zestawem **stałych wpisanych w kod** (`core/most_mcp.go`). Ich znajomość jest
warunkiem poprawnego wdrożenia, bo wyznaczają wymagania wobec maszyny po drugiej stronie:

| Stała | Wartość | Konsekwencja |
| --- | --- | --- |
| ścieżka uruchomienia mostu | `/opt/danaco/most-konsoli/uruchom-stdio.sh` | maszyna zdalna **musi** mieć ten skrypt pod tą dokładnie ścieżką |
| program połączenia | `ssh` | wymagany klient `ssh` w `PATH` procesu rdzenia |
| rodzaj wpisu | `stdio` | most rozmawia przez strumienie, nie po sieci HTTP |
| zmienna korzeni | `DANACO_MOST_KORZENIE` | korzenie przekazywane środowiskiem |
| rozdzielnik korzeni | `:` (dwukropek) | **korzeń zawierający dwukropek, jak ścieżka Windows `C:\…`, rozjedzie wykaz** |
| użytkownik domyślny | `ubuntu` | adres bez konta zakłada system rodziny Linux |
| port domyślny | `22` | standardowy port SSH |
| przedrostek klucza wpisu | `mcp-danaco-pulpit-console-` | nazwa mostu jest po stronie serwera stała |

Cztery wnioski wdrożeniowe wynikają z tego wprost:

1. **Most jest przewidziany dla maszyny linuksowej.** Ścieżka `/opt/danaco/...`, użytkownik domyślny
   `ubuntu` i skrypt powłoki nie mają odpowiednika na Windows. Most wskazujący maszynę Windows nie
   zadziała bez zmian po stronie mostu.
2. **Ścieżki skryptu nie da się zmienić konfiguracją.** Nie ma dla niej ani kolumny, ani parametru
   punktu — jest stałą kodu. Maszyna zdalna musi się do niej dostosować. **[DZIAŁA CZĘŚCIOWO]**
3. **Korzenie muszą być ścieżkami bez dwukropka**, bo rozdzielnikiem jest dwukropek. Ścieżki
   Windows w korzeniach mostu są więc wykluczone.
4. Klient `ssh` musi być osiągalny w `PATH` **procesu rdzenia** (a nie tylko w powłoce Operatora),
   ponieważ to proces kanału go uruchamia, dziedzicząc środowisko po rdzeniu (rozdział 12.4).

Uwaga bezpieczeństwa: wpis niesie `StrictHostKeyChecking=accept-new`, co oznacza **automatyczne
przyjęcie klucza hosta przy pierwszym połączeniu**. Jest to wybór wygody nad rygorem.
**[DO DECYZJI OPERATORA]** — czy wcześniej rozprowadzić klucze hostów maszyn mostowych do pliku
`known_hosts` stacji, by pierwsze połączenie nie było przyjmowane w ciemno.

---

## 16. Konfiguracja dostawców modeli

### 16.1 Wiersz rejestru kanału jako opis dostawcy

Danaco Console **nie zna żadnego dostawcy modeli z nazwy**. W kodzie nie ma ani jednej gałęzi
warunkowanej nazwą dostawcy, nie ma zaszytych adresów punktów końcowych i nie ma wbudowanych kształtów
żądań konkretnych usług. Adapter `api` (`models/adapter_api.go`) jest **generyczny**: wszystko, co
odróżnia jednego dostawcę od drugiego, pochodzi z **wiersza rejestru kanałów** — z kolumny
`parametry_json` i z kolumny `poswiadczenie_odwolanie`.

Konsekwencja jest wprost taka, jak zapowiada przyjęte rozstrzygnięcie: **nowy dostawca to nowy wiersz danych, nie nowy
typ w kodzie**. Konfiguracja dostawcy jest więc czynnością operacyjną, nie programistyczną.

Kanał `api` odmawia budowy w **jednym jedynym przypadku**: gdy wiersz nie niesie parametru `base_url`.
Bez adresu nie ma dokąd wysłać żądania. Wiersz taki trafia na listę pominiętych z powodem
`models: kanał "…" bez parametru base_url`, a pozostałe kanały pracują normalnie.

### 16.2 Kształt `parametry_json`

Pełny wykaz parametrów odczytywanych przez adapter `api`:

| Parametr | Rola | Wartość domyślna |
| --- | --- | --- |
| `base_url` | pełny adres punktu końcowego | **wymagany** |
| `strumien` | czy żądać odpowiedzi strumieniem | `true` (wyłącznie wartość `false` wyłącza) |
| `sciezka_tekstu` | ścieżka do porcji tekstu w zdarzeniu strumienia | `choices.0.delta.content` |
| `sciezka_odpowiedzi` | ścieżka do treści w odpowiedzi bez strumienia | `choices.0.message.content` |
| `sciezka_bledu` | ścieżka do komunikatu błędu w odpowiedzi | brak — użyta zostaje surowa treść |
| `naglowek_klucza` | nazwa nagłówka niosącego dane dostępowe | `Authorization` |
| `przedrostek_klucza` | przedrostek wartości tego nagłówka | `Bearer ` (ze spacją) |
| `naglowki` | obiekt dodatkowych nagłówków (wartości tekstowe) | brak |
| `cialo_dodatkowe` | obiekt scalany z ciałem żądania | brak |
| `limit_sekund` | czas oczekiwania na odpowiedź | `600` (dziesięć minut) |
| `adapter` | wskazanie fabryki adaptera (rozdział 10.1.1) | kolumna `rodzaj_kanalu` |

Uwagi do wartości domyślnych:

- **Parametr `przedrostek_klucza` odróżnia brak od pustki.** Ustawienie go na pusty łańcuch jest
  świadomym ustawieniem („nagłówek bez przedrostka”) i wygrywa z wartością domyślną `Bearer `.
  Służy temu osobna funkcja odczytu rozróżniająca oba przypadki.
- **Limit czasu nieczytelny albo niedodatni** schodzi na dziesięć minut. Wartość jest dobrana pod
  długie odpowiedzi modeli rozumujących.

Żądanie budowane jest zawsze metodą **POST** z nagłówkiem `Content-Type: application/json`, a przy
odpowiedzi strumieniowej dodatkowo `Accept: text/event-stream`. Ciało domyślne ma postać:

```json
{
  "model":    "<identyfikator_modelu wiersza albo wskazanie okna>",
  "messages": [{"role": "system", "content": "<nakładka tożsamości>"},
               {"role": "user",   "content": "<treść zapytania>"}],
  "stream":   true
}
```

Wiadomość systemowa pojawia się **wyłącznie wtedy, gdy nakładka tożsamości jest niepusta**. Obiekt
z parametru `cialo_dodatkowe` jest scalany z powyższym i może nadpisać każdy z jego kluczy — to nim
dostraja się kształt żądania do wymagań konkretnego dostawcy, bez zmiany kodu.

Przykładowy wiersz parametrów dla punktu końcowego zgodnego z powszechnym kształtem rozmowy:

```json
{
  "base_url": "https://api.przyklad.pl/v1/chat/completions",
  "strumien": true,
  "naglowek_klucza": "Authorization",
  "przedrostek_klucza": "Bearer ",
  "sciezka_tekstu": "choices.0.delta.content",
  "sciezka_bledu": "error.message",
  "limit_sekund": 600
}
```

### 16.3 Poświadczenia dostawcy

Rozstrzygnięcie jest jednoznaczne i konsekwentne z zasadą trzymania sekretów poza repozytorium: **klucz dostępowy nie mieszka ani
w bazie, ani w repozytorium**. Kolumna `poswiadczenie_odwolanie` wiersza kanału niesie **nazwę
zmiennej środowiskowej**, a adapter odczytuje jej wartość funkcją `os.Getenv` w chwili budowania
żądania.

Zachowania warte zapamiętania:

- **Odwołanie puste** oznacza punkt końcowy bez uwierzytelnienia — żądanie idzie bez nagłówka klucza.
- **Odwołanie wskazujące zmienną nieustawioną albo pustą** kończy wywołanie błędem
  `models: kanał "…": brak wartości pod odwołaniem <NAZWA>`. Jest to błąd jednego wywołania, nie awaria
  rdzenia.
- **Do prowenancji wywołania trafia wyłącznie nazwa odwołania**, nigdy wartość klucza. Prowenancja jest
  bezpieczna do zapisania i pokazania.
- Zmienna musi istnieć w środowisku **procesu rdzenia** — to on wykonuje żądanie HTTP. Ustawienie jej
  po starcie rdzenia nie odniesie skutku (rozdział 12.4).

Przykład: wiersz z `poswiadczenie_odwolanie = 'DOSTAWCA_KLUCZ_API'` wymaga uruchomienia rdzenia
w środowisku, w którym zmienna `DOSTAWCA_KLUCZ_API` niesie klucz.

```powershell
$env:DOSTAWCA_KLUCZ_API = "<klucz dostawcy>"
.\danaco-console.exe
```

**Ograniczenie zasadnicze — konta a kanał `api`.** Konta rodzaju `api` i `sdk` są w produkcie
kompletne po stronie katalogu: mają tabelę, komendy `account.*`, widok w oknie modeli, kolumny
dostawcy, modelu domyślnego i adresu bazowego. **Nie prowadzi jednak od nich żadna droga wykonawcza.**
Adapter `api` bierze adres i poświadczenie **wyłącznie z wiersza kanału**, a pula rotacji obsługuje
wyłącznie rodzaj `cli`. Kolumna `kanal_modelu.konto_id` jest przechowywana, lecz adapter `api` jej nie
czyta. **[NIEZINTEGROWANE]**

Praktycznie znaczy to tyle: **konfigurując dostawcę API, cała konfiguracja musi znaleźć się w wierszu
kanału.** Zakładanie konta rodzaju `api` jest dziś czynnością porządkową (ewidencja), nie
konfiguracyjną. **[DO DECYZJI OPERATORA]** — czy prowadzić ewidencję kont `api` mimo braku drogi
wykonawczej, czy odłożyć ją do czasu spięcia kont z adapterem.

---

## 17. Konfiguracja sesji i okien operacyjnych

Okno komunikacji jest **najwęższym poziomem zasięgu** konfiguracji i to ono rozstrzyga parametry
wykonania tury. Rozdział opisuje, które z ustawień okna **rzeczywiście docierają do wywołania**,
a które są zapisywane i nie mają odbiorcy. Rozdzielenie jest tu najważniejszą treścią całego dokumentu
po rozdziale 20 — pozwala nie tracić czasu na strojenie ustawień, które nic nie robią.

### 17.1 Ustawienia wpływające na wywołanie

Wiersz argumentów programu `claude` składa funkcja `Argumenty` (`injection/argumenty.go`). Szkielet
każdego wywołania jest stały:

```
claude -p --input-format stream-json --output-format stream-json --verbose …
```

Bez `--verbose` program nie wypuszcza zdarzeń pośrednich, więc strumień przestałby być strumieniem.
Do szkieletu doklejane są przełączniki **wyłącznie wtedy, gdy odpowiadające im ustawienie jest
wypełnione** — brak ustawienia nie jest błędem, tylko brakiem przełącznika.

Zestawienie tego, co **dociera** do wywołania:

| Ustawienie | Skąd pochodzi | Przełącznik | Stan |
| --- | --- | --- | --- |
| ścieżka programu | parametr `program` wiersza kanału; brak → `claude` | (nazwa programu) | **[DZIAŁA]** |
| identyfikator modelu | kolumna `identyfikator_modelu` **wiersza kanału** | `--model` | **[DZIAŁA CZĘŚCIOWO]** — patrz niżej |
| tryb uprawnień | pole okna komunikacji | `--permission-mode` | **[DZIAŁA]** |
| katalogi robocze okna | lista katalogów okna | `--add-dir` (powtarzany) | **[DZIAŁA]** |
| katalog roboczy sesji | ustawienia `katalog.roboczy.*` | katalog startowy procesu | **[DZIAŁA]** |
| konfiguracja mostów MCP | nadania okna | `--mcp-config` | **[DZIAŁA]** |
| nakładka tożsamości | dokumenty tożsamości i tryb | `--system-prompt` albo `--append-system-prompt` | **[DZIAŁA]** |
| `CLAUDE_CONFIG_DIR` | konto z puli rotacji | zmienna środowiska procesu | **[DZIAŁA]** |

Zastrzeżenie do modelu: wywołanie bierze model z **wiersza rejestru kanału**, a nie z ustawienia okna.
Rozstrzyganie (`WybranyModel`) najpierw sprawdza wskazanie okna, a dopiero potem kolumnę wiersza —
lecz **wskazanie okna nigdy nie jest wypełniane** (rozdział 17.2). Skutek praktyczny: **modelem steruje
się dziś przez kolumnę `identyfikator_modelu` wiersza kanału**, komendą `channel.update` z polem
`model` albo zapisem do bazy. Jest to droga skuteczna i sprawdzona.

### 17.2 Ustawienia zapisywane, lecz nieaplikowane

Poniższe ustawienia mają **kompletny zapis**: pole w kontrakcie albo klucz katalogu ustawień, obsługę
komendy, zapis do bazy i sterowanie w interfejsie. Nie mają **odbiorcy w wykonaniu** — wartość zostaje
w bazie i nic z niej nie wynika. Wykaz pochodzi z audytu i został potwierdzony w kodzie składania
zapytania kanału (`core/adapter_rozmowa.go`, funkcja `zapytanieKanalu`), które nie wypełnia
odpowiadających pól.

| Ustawienie | Gdzie jest zapisywane | Dlaczego nie działa | Stan |
| --- | --- | --- | --- |
| **model per okno** | klucze sterowania okna | pole `Model` zapytania nigdy nie wypełniane; wygrywa kolumna wiersza kanału | **[NIEZINTEGROWANE]** |
| **model zapasowy** | `config.set`, klucz `window.fallbackChannelId` | serwer nie zna tego klucza; pole `ModelZapasowy` nigdy nie wypełniane, `--fallback-model` nigdy nie podawany | **[NIEZINTEGROWANE]** |
| **nakład rozumowania** | `config.set`, klucz `window.reasoningEffort` | jak wyżej; `--effort` nigdy nie podawany | **[NIEZINTEGROWANE]** |
| **host wykonania** | `config.set`, klucz `window.executionHost` | serwer nie zna klucza; brak drogi wykonania zdalnego (rozdz. 14.2) | **[NIEZINTEGROWANE]** |
| **zasięg wykonania** (`local`/`core`/`remote`) | pole okna, kolumna `srodowisko_wykonania` | wyłącznie etykieta w prowenancji | **[NIEZINTEGROWANE]** |
| **`harness.program_claude`** | katalog ustawień (migracja 012) | wykonanie czyta parametr `program` wiersza kanału, nie ten klucz | **[NIEZINTEGROWANE]** |
| **`harness.plik_ustawien`** | katalog ustawień | pole `PlikUstawien` nigdy nie wypełniane; `--settings` nigdy nie podawany | **[NIEZINTEGROWANE]** |
| **`harness.konfiguracja_mcp`** | katalog ustawień | `--mcp-config` pochodzi z nadań okna, nie z tego klucza | **[NIEZINTEGROWANE]** |
| **ustawienia izolacji** (11 pozycji × 8 poziomów) | katalog ustawień | zdefiniowane i rozstrzygane, **zero egzekutorów** | **[NIEZINTEGROWANE]** |

Rozjazd nazw jest tu zjawiskiem systemowym, nie pojedynczą pomyłką: klient zapisuje trzy klucze okna
(`window.executionHost`, `window.fallbackChannelId`, `window.reasoningEffort`) wymyślone po swojej
stronie, ponieważ — jak mówi komentarz w `client/src/sterowanie/klucze-ustawien.ts` — „klucze ustawień
nie mają w kontrakcie swojego słownika”. Serwer prowadzi własny katalog kluczy (migracja 012) i tych
trzech nie zawiera. Zapis się udaje, odczyt przez klienta się udaje, a wykonanie ich nie widzi.

**Wniosek dla Operatora: nie należy oczekiwać, że zmiana modelu, nakładu rozumowania czy hosta
w panelu sterowania okna zmieni cokolwiek w zachowaniu modelu.** Skuteczne dźwignie konfiguracyjne to
te z rozdziału 17.1.

**[DO DECYZJI OPERATORA]** — czy uzgodnić słownik kluczy ustawień okna po obu stronach (dopisanie
wierszy do katalogu ustawień oraz odczyt w składaniu zapytania), czy usunąć sterowania z interfejsu
do czasu spięcia. Dokumentacja nie rozstrzyga tego wyboru.

### 17.3 Ciągłość rozmowy

Rozstrzygnięcie tego podrozdziału ma najpoważniejsze następstwa użytkowe w całym produkcie i wymaga
postawienia sprawy wprost.

**Model nie pamięta poprzednich wiadomości okna.** Każda tura uruchamia **nowy proces programu
`claude` bez historii**. Przełącznik `--resume` jest obsłużony w budowie argumentów (pole `Wznowienie`
ustawień), lecz **pole to nigdy nie jest wypełniane** przy składaniu zapytania. Identyfikator rozmowy
programu (`IdSesjiCLI`) jest odczytywany ze strumienia poprzedniej tury i porzucany.
**[NIEZINTEGROWANE]**

Praktycznie znaczy to, że rozmowa w oknie jest ciągiem niezależnych zapytań: model odpowiadający na
drugą wiadomość nie zna pierwszej. Kontekst trzeba nieść w treści zapytania.

Co natomiast **działa** i nie należy tego mylić z powyższym:

- **Historia jest utrwalana.** Wiadomości zapisywane są do tabeli `wiadomosc` w łańcuchu
  sesja → karta → okno → wiadomość, pilnowanym więzami kluczy obcych. **[DZIAŁA]**
- **Sesje i okna są odtwarzane z bazy po restarcie rdzenia.** **[DZIAŁA]**
- **Praca przeżywa rozłączenie klienta** — tura biegnie dalej, a jej wynik trafia do bazy. **[DZIAŁA]**
- **Zatrzymanie tury** komendą `message.stop` działa. **[DZIAŁA]**

Czego **nie ma** wokół historii:

- **Interfejs nigdy nie wywołuje `message.list`.** Historia jest w bazie, lecz po odświeżeniu strony
  znika z widoku. **[NIEZINTEGROWANE]**
- **`session.create` i `window.create` nie zapisują wiersza** aż do pierwszej wiadomości.
- **`window.update` nie zapisuje do bazy**, a sesje trafiają zawsze do pierwszego środowiska
  (`talkin`, „Karta domyślna”). **[NIEZINTEGROWANE]**

**[DO DECYZJI OPERATORA]** — czy do czasu spięcia `--resume` prowadzić pracę w konwencji „jedna tura =
jedno kompletne polecenie z kontekstem”, czy wstrzymać wdrożenie produkcyjne do czasu zamknięcia
ciągłości rozmowy. Jest to rozstrzygnięcie o charakterze eksploatacyjnym, nie technicznym.

---

## 18. Konfiguracja komponentów dostępnych w bieżącej wersji

### 18.1 Komponenty konfigurowalne

Poniższe elementy interfejsu mają **kompletną drogę**: widok, komendę kontraktu, obsługę w rdzeniu
i zapis do bazy. Są to jedyne miejsca, w których konfiguracja produktu prowadzona z interfejsu daje
trwały skutek.

| Komponent | Zakres konfiguracji | Komendy kontraktu | Stan |
| --- | --- | --- | --- |
| **Okno komunikacji** | prowadzenie rozmowy, zatrzymanie tury, strumień odpowiedzi | `message.send`, `message.stop` | **[DZIAŁA]** |
| **Okno konfiguracji** | katalog kategorii, definicji i opcji sterowany danymi z bazy (migracja 012); osiem poziomów zasięgu i trzy osie | `config.get`, `config.set`, `config.reset` | **[DZIAŁA]** — zapis; **odbiorcy części kluczy brak** (rozdz. 17.2) |
| **Okno dostępów** | punkty dostępu rodzaju `mcpBridge` oraz nadania per okno | `access.point.*`, `access.grant.*` | **[DZIAŁA CZĘŚCIOWO]** — bez rodzaju `localDirectory` (rozdz. 13.3) |
| **Okno modeli, kont i tożsamości** | katalog kont, dokumenty tożsamości per oś | `account.*`, `identity.*` | **[DZIAŁA CZĘŚCIOWO]** — pula wymaga restartu, konta `api`/`sdk` bez drogi wykonawczej |
| **Motyw wizualny** | dwa motywy (jasny, ciemny), gęstość zwarta, komplet ikon wg manifestu | — (warstwa prezentacji) | **[DZIAŁA]** |
| **Scena okien operacyjnych** | układ jedno-, dwu- i trójokienny, role okien | `window.*` | **[DZIAŁA CZĘŚCIOWO]** — `window.update` nie zapisuje do bazy |

Poza konfiguracją w ścisłym sensie (zapis trwały) dwa kolejne widoki mają po pracach nad zasilaniem widoków z rdzenia **kompletną
drogę odczytu i nawigacji** ze stanu rdzenia, choć niczego nie konfigurują: pulpit Mission Control
(rozdział 7.5, rozdział 18.2) i strefa „Sesje w tle” strony głównej (rozdział 7.5) — obie opisane
w miejscu właściwym dla ich charakteru, nie w tabeli powyżej.

Uzupełnienie do okna konfiguracji: katalog ustawień jest **sterowany danymi**, więc nowa pozycja
okna konfiguracji to nowy wiersz migracji, nie zmiana kodu. Rozstrzyganie prowadzi rezolwer ośmiu
poziomów zasięgu (od globalnego po okno) i trzech osi (`platform`, `model`, `account`). Zapis pod
adresem złożonym (poziom + byt poziomu + oś + byt osi) działa niezawodnie. Ograniczeniem jest wyłącznie
liczba kluczy, które mają odbiorcę w wykonaniu.

### 18.2 Komponenty przewidziane koncepcją, niezaimplementowane

Poniższe pozycje **istniały w nawigacji, nazewnictwie albo definicjach kontraktu bez wykonawcy**
w chwili audytu. Dokument wymienia je po to, by Operator nie próbował ich konfigurować — nie dlatego,
że są funkcjami produktu. **Trzy pozycje częściowo wyszły z tej kategorii wraz z pakietem D fali
naprawczej** — oznaczone poniżej wprost jako **[DZIAŁA CZĘŚCIOWO]** wraz z granicą działania; reszta
wykazu pozostaje bez wykonawcy tak, jak zastał ją audyt.

| Pozycja | Stan faktyczny |
| --- | --- |
| **Moduły operacyjne** (Studio, Workspace, Browser, Research, Library, Translate, Roundtable, Design, Assistant, Terminal, Developer, Diagnostics, Apps, Agents, Automations) | pozycje menu bez widoków; **64 z 79 wierszy katalogu okien operacyjnych w bazie nie mają odpowiadającego widoku** ([INSTRUKCJA-UZYTKOWANIA.md](INSTRUKCJA-UZYTKOWANIA.md) rozdz. 7); wybór modułu w nawigacji **nie przeładowuje przestrzeni roboczej** — 26 pozycji nawigacji modułowej trzech środowisk modułowych, a wraz z sześcioma sekcjami panelu orkiestracji MultitaskingAI wszystkie 32 pozycje nawigacji czterech środowisk, prowadzi do jednego ekranu sceny okien ([INSTRUKCJA-UZYTKOWANIA.md](INSTRUKCJA-UZYTKOWANIA.md) rozdz. 4.2). **[BRAK]** |
| **Środowiska** TalkIn, WorkSpace, CodeStudio, MultitaskingAI | nawigacja i motta istnieją, realizacja funkcjonalna **brak**; MultitaskingAI to sześć etykiet — bez panelu orkiestracji, ról, kolejek ról, sieci podagentów, harmonogramu i monitora. **[BRAK]** |
| **Pulpit Mission Control** | po pracach nad zasilaniem widoków z rdzenia rysuje się **ze stanu rdzenia** (`session.list`/`window.list`/`channel.list` i zdarzenia pochodne), nie z danych przykładowych — opis pełny w rozdz. 7.5. Kolumna kolejek pozostaje uczciwie pusta, bo silnika kolejek nadal brak. **[DZIAŁA CZĘŚCIOWO]** |
| **Kolejki** | `queue.create` i `queue.action` przestawiają wyłącznie kolumnę stanu; pozycje kolejki nigdy nie powstają, silnika wykonania brak. Pulpit nasłuchuje `queue.changed`, ale samo zdarzenie nie niesie więcej niż zmianę kolumny stanu. **[NIEZINTEGROWANE]** |
| **Przekazanie zlecenia koordynator → wykonawca w interfejsie** | wyłącznie animacja lokalna, bez komendy kontraktu. Sama pętla koordynator–wykonawca **istnieje w rdzeniu** i startuje turą, lecz jej postęp nie ma odbiorcy w interfejsie poza pulpitem (patrz telemetria postępu niżej). **[NIEZINTEGROWANE]** |
| **Narzędzia modelu z kontraktu** (39 deklaracji) | wygenerowane, **zero konsumentów** — sterowanie platformą przez model nie działa. **[NIEZINTEGROWANE]** |
| **Nawigacja platformy przez kontrakt** (`home.enter`, `environment.list`/`enter`, `module.list`, `workspace.enter`, `window.state.get`) | zaimplementowana po obu stronach, **niewywoływana przez żaden widok**. **[NIEZINTEGROWANE]** **Wyjątek po pracach nad zasilaniem widoków z rdzenia:** `session.bind` i `session.focus` **są już wywoływane** — przez czynność powrotu strefy „Sesje w tle” strony głównej (rozdz. 7.5); reszta wykazu pozostaje martwa. |
| **Pamięć wielopoziomowa** | repozytorium istnieje, zero komend kontraktu, zero konsumentów. **[NIEZINTEGROWANE]** |
| **Rozszerzenia i marketplace** | brak w kontrakcie wykonawczym. **[BRAK]** |
| **Telemetria postępu w interfejsie** | po pracach nad zasilaniem widoków z rdzenia `progress.changed` **ma konsumenta** — kolumnę procesów pulpitu; `session.changed`/`window.changed` mają konsumentów w pulpicie i w strefie „Sesje w tle” strony głównej. Poza tymi dwoma widokami zdarzenia nadal nie mają subskrybenta, a `window.state.get` zawsze zwraca stan oczekiwania. **[DZIAŁA CZĘŚCIOWO]** |
| **Elementy paska górnego** (dzwonek powiadomień, awatar, przełącznik AOD, pole wyszukiwania) | atrapy prezentacyjne. **[ATRAPA]** |
| **Treści obszerne jako pliki** | kolumna odwołania do treści nigdy nie jest wypełniana. **[BRAK]** |
| **Uwierzytelnianie i rozdział ról użytkowników** | nie istnieje; kto ma dostęp do gniazda, ma pełnię możliwości (rozdz. 2.4). **[BRAK]** |

Uczciwe podsumowanie rozdziału: **Danaco Console w wersji v2.0 jest fundamentem technicznym
z jedną kompletną drogą użytkową — oknem komunikacji na kanale modelu — oraz kompletem narzędzi
konfiguracyjnych wokół niego.** Warstwa modułów i środowisk jest zapowiedzią, nie funkcją. Wdrożenie
należy planować wyłącznie w oparciu o rozdział 18.1.

---

## 19. Weryfikacja poprawności instalacji

Rozdział opisuje kontrole wykonywane **na drzewie źródeł** (rozdziały 19.1–19.2) oraz kontrole
**produktu złożonego i pracującego** (19.3–19.4). Pierwsze wymagają kompletnego środowiska budowy;
drugie wystarczy uruchomić na stacji pracy.

### 19.1 Kontrole źródeł

Trzy kontrole podstawowe wykonuje się z katalogu `budowa/`:

```bash
gofmt -l server shared     # wypisuje pliki niesformatowane; wynik pusty = poprawnie
go vet ./...               # analiza statyczna rdzenia
go test ./...              # zestaw sprawdzianów rdzenia
```

Kontrola `gofmt -l` **wypisuje nazwy plików wymagających formatowania** — poprawny wynik to brak
wyjścia. Zestaw sprawdzianów rdzenia obejmuje ponad sto plików testowych i pokrywa sześć obszarów: migracje na
czystej bazie, jednoznaczność numeracji migracji, rejestr kanałów, rotację kont, budowę argumentów
programu oraz przekrój wszystkich komend kontraktu.

Po stronie klienta, z katalogu `budowa/client/`:

```bash
npm run typy      # tsc --noEmit — sama kontrola typów
npm run testy     # vitest run
npm run build     # tsc --noEmit && vite build — kontrola typów i pakiet
```

### 19.2 Brama jakości

Skrypt `budowa/scripts/brama.sh` spina wszystkie kontrole w jeden automat. Uruchamia się go z korzenia
repozytorium:

```bash
bash budowa/scripts/brama.sh            # tryb pełny
bash budowa/scripts/brama.sh --szybka   # tryb szybki
```

Brama **nie modyfikuje drzewa roboczego** — sprawdzian mutacyjny kontraktu pracuje na kopii w katalogu
tymczasowym poza repozytorium. Kod wyjścia `0` oznacza przejście, `1` — naruszenia.

Zakres kontroli:

| Sekcja bramy | Co sprawdza | Tryb szybki |
| --- | --- | --- |
| Rdzeń Go | `go build ./...`, `go vet ./...`, `go test ./...`, `gofmt` | tak |
| Klient TypeScript | `tsc --noEmit`, `npm run testy`, a w trybie pełnym `npm run build` | częściowo |
| Osiągalność drzewa interfejsu | czy każdy plik `client/src` jest osiągalny importem z `main.ts` | tak |
| Powłoka natywna Rust | formatowanie; w trybie pełnym budowa, testy i analiza | częściowo |
| Modularność | progi rozmiaru plików | tak |
| Atrapy | znaczniki niedokończonej roboty, pliki bez treści logicznej, pakiety Go bez importu, eksporty TypeScriptu bez odbiorcy, nazwy-wysypiska (`utils`, `helpers`, `common`, `misc`) | tak |
| Żywotność kontraktu | komenda kontraktu bez wywołania w kliencie; eksport pakietu Go bez odbiorcy | tak |
| Sprawdzian obustronny kontraktu | czy zmiana **identyfikatora** komendy przerywa kompilację po obu stronach | nie |
| **Dymny sprawdzian pionu** | pełny pion na kanale echo, uruchomieniowo | nie |

Tryb szybki pomija budowę klienta, budowę powłoki, sprawdzian obustronny kontraktu i dymny sprawdzian
pionu — jest przeznaczony do haka przed zatwierdzeniem zmian, gdzie liczy się czas.

**Dymny sprawdzian pionu** (`budowa/scripts/dym-pionu.mjs`) zasługuje na wyróżnienie, bo jest jedyną
kontrolą **uruchomieniową**: buduje rdzeń, uruchamia go w katalogu tymczasowym na porcie **17999**
(port 17870 bywa zajęty przez rdzeń roboczy) i przechodzi cały pion na kanale echo — bez sieci, bez
konta i bez programu `claude`. Można go uruchomić samodzielnie:

```bash
node budowa/scripts/dym-pionu.mjs        # port domyślny 17999
node budowa/scripts/dym-pionu.mjs 17998  # port wskazany
```

Kolejność kroków sprawdzianu jest jednocześnie **wzorcem poprawnego pionu**:

```
connection.hello → channel.add(echo) → session.create → window.create
→ message.send(stream) → stream.chunk… → message.changed(complete)
```

Po przejściu proces rdzenia jest zatrzymywany, a pliki tymczasowe usuwane. Sprawdzian nie dotyka
katalogu danych instalacji roboczej.

### 19.3 Kontrola produktu złożonego

Na stacji, gdzie produkt ma pracować, warto potwierdzić cztery rzeczy w tej kolejności:

1. **Obecność artefaktów.** W katalogu produktu leżą `danaco-console.exe`, katalog `client\dist`
   z plikiem `index.html`, a w wariancie pełnym także `Danaco Console.exe` (rozdział 6.2).
2. **Start rdzenia.** Uruchomienie z katalogu produktu wypisuje wiersz startowy postaci
   `start: rola=all port=17870 dane=<katalog> baza=<ścieżka>`. Wiersz ten potwierdza naraz cztery
   rzeczy: rolę, port, katalog danych i otwarcie bazy.
3. **Powstanie bazy.** W katalogu danych pojawił się plik `danaco-console.db`.
4. **Serwowanie interfejsu.** Adres `http://127.0.0.1:17870/` zwraca stronę interfejsu, a nie pustkę.

Kontrola wersji schematu bazy (przy zainstalowanym `sqlite3`, opcjonalnie):

```bash
sqlite3 "%LOCALAPPDATA%\DanacoConsole\danaco-console.db" \
  "SELECT MAX(wersja) FROM migracja;"
```

Poprawny wynik dla bieżącego wydania to **15**. Pełny wykaz zastosowanych kroków:

```sql
SELECT wersja, nazwa, zastosowano FROM migracja ORDER BY wersja;
```

Oczekiwany zbiór wersji to 1–9 oraz 12–15 (luka 010–011 jest zamierzona, rozdział 9.2).

### 19.4 Kontrola połączenia i kontraktu

Kontrola najgłębsza sprawdza, czy klient i rdzeń mówią tym samym kontraktem. Wykonuje się ją po
połączeniu z gniazdem `ws://127.0.0.1:17870/ws` i wysłaniu **koperty** kontraktu w kształcie
`{ type, id, payload, timestamp }`.

**Krok 1 — powitanie.** Komenda `connection.hello` z treścią żądania:

```json
{ "clientId": "kontrola", "clientVersion": "1.0", "protocolVersion": "1.0" }
```

Wynik niesie `serverVersion` (wersja rdzenia), `protocolVersion` (wersja protokołu rdzenia) oraz —
opcjonalnie — `capabilities` i `commands` (wykaz komend obsługiwanych przez rdzeń). **Zgodność pola
`protocolVersion` po obu stronach jest właściwą kontrolą kontraktu.** Rozbieżność oznacza, że klient
i rdzeń pochodzą z różnych wydań i nie wolno ich łączyć.

**Krok 2 — rejestr kanałów.** Komenda `channel.list` z treścią `{}` (albo `{"enabledOnly": true}`).
Wynik niesie pole `channels`. Na **świeżej instalacji lista jest pusta** i jest to stan oczekiwany,
nie usterka — rejestr kanałów nie ma zaczynu (rozdział 10.2). Pusta lista jest sygnałem, że należy
wykonać procedurę z rozdziału 20.

Odpowiedź o statusie `error` na `channel.list` oznacza natomiast rzeczywisty problem: najczęściej
niedostępną bazę albo niezgodność kontraktu.

Uwaga metodyczna: dokument nie podaje gotowego polecenia klienta WebSocket dla Windows, bo produkt
takiego narzędzia nie dostarcza i nie należy zakładać obecności zewnętrznego. Najprostszą dostępną
drogą kontroli pionu jest **uruchomienie `dym-pionu.mjs`** (rozdział 19.2), które wykonuje oba
powyższe kroki i cały pion, korzystając z klienta WebSocket wbudowanego w Node.js.

---

## 20. Procedura doprowadzenia świeżej instalacji do stanu gotowego do pracy

### 20.1 Dlaczego procedura jest potrzebna

Świeżo zbudowana i uruchomiona instalacja Danaco Console **nie wykona ani jednej tury modelu**. Nie
jest to usterka konfiguracji stacji ani błąd budowy — jest to **stan zastany produktu**, stwierdzony
audytem i potwierdzony w kodzie w trzech niezależnych miejscach:

1. **Rejestr kanałów jest pusty.** Żadna z trzynastu migracji nie zasiewa wiersza w tabeli
   `kanal_modelu` (rozdział 10.2). Bez wiersza kanału okno komunikacji nie ma czym rozmawiać.
2. **Interfejs nie ma widoku zakładania kanałów.** Klient wywołuje wyłącznie `channel.list`; jedyne
   miejsce wywołujące `channel.add` to moduł podpięty do stanowiska podglądu rozmowy, nie do
   produkcyjnej ścieżki okna.
3. **Pula kont jest pusta**, a kanał główny bez konta odmawia wykonania (rozdział 11.3).

Procedura poniżej jest **obejściem brakującego zaczynu**, wykonywanym raz na instalację. Jest to
działanie ustalone jako wynik audytu i tak też należy je traktować: jako czynność wdrożeniową
konieczną w wersji v2.0, nie jako docelowy sposób konfiguracji produktu.

Wykonać należy wszystkie pięć kroków; pominięcie któregokolwiek zostawia instalację niezdolną
do pracy.

### 20.2 Krok 1 — profil konta kanału głównego

Krok wykonywany **poza produktem**, przy użyciu programu `claude`.

1. Utworzyć katalog konfiguracji konta:

   ```powershell
   New-Item -ItemType Directory -Force `
     "C:\Users\Operator\AppData\Local\DanacoConsole\profile\konto-glowne"
   ```

2. Uwierzytelnić program `claude` w tym katalogu, ustawiając mu `CLAUDE_CONFIG_DIR` na czas
   uwierzytelnienia. Sposób uwierzytelnienia należy do programu `claude`; Danaco Console nie
   uczestniczy w tej czynności i nie przechowuje poświadczeń (rozdział 11.2).
3. Potwierdzić, że program działa z tym katalogiem — to znaczy, że w katalogu powstały pliki
   konfiguracyjne programu.

Zalecana konwencja: wszystkie katalogi kont jako podkatalogi jednego katalogu profili
(rozdział 11.4). Ułatwia to kopie zapasowe i przenoszenie instalacji.

**[DO DECYZJI OPERATORA]** — liczba kont w puli rotacji. Jedno konto wystarcza do pracy; wiele kont
ma sens dopiero przy realnym ryzyku wyczerpania limitu.

### 20.3 Krok 2 — start rdzenia

Rdzeń musi pracować, bo dwa kolejne kroki wykonywane są komendami kontraktu po WebSocket.

```powershell
Set-Location C:\DanacoConsole_App
.\danaco-console.exe --profile "C:\Users\Operator\AppData\Local\DanacoConsole\profile"
```

Wskazanie `--profile` nie jest tu obowiązkowe — jest zabezpieczeniem: gdyby krok 4 (zasianie konta
w bazie) został pominięty, pula sięgnie po źródło zapasowe i weźmie konta z podkatalogów katalogu
profili (rozdział 11.3).

Potwierdzeniem startu jest wiersz `start: rola=all port=17870 dane=… baza=…` oraz osiągalny adres
`http://127.0.0.1:17870/`.

### 20.4 Krok 3 — zasianie wiersza kanału

Krok wykonywany komendą kontraktu **`channel.add`** po WebSocket (`ws://127.0.0.1:17870/ws`), w
kopercie `{ type, id, payload, timestamp }`.

**Wariant A — kanał główny (Claude Code CLI).** Treść żądania:

```json
{
  "name":    "Kanał główny — Claude Code CLI",
  "kind":    "cli",
  "model":   "<identyfikator modelu>",
  "enabled": true,
  "config":  {}
}
```

Pole `model` wypełnia kolumnę `identyfikator_modelu` i jest **jedyną skuteczną drogą wskazania
modelu** (rozdział 17.1). Aby wskazać ścieżkę programu `claude` inną niż z `PATH`, należy podać
w `config` parametr `program`:

```json
{ "program": "C:\\Program Files\\Claude\\claude.exe" }
```

**Wariant B — kanał próbny echo.** Rodzaj musi należeć do zbioru dopuszczonego warunkiem `CHECK`
(`cli`, `api`, `sdk`, `lokalny`), a adapter wybiera parametr `adapter`:

```json
{
  "name":    "Echo kontrolne",
  "kind":    "lokalny",
  "enabled": true,
  "config":  { "adapter": "echo", "porcja": 16, "przedrostek": "echo: " }
}
```

Kształt ten jest **zweryfikowany uruchomieniowo** — posługuje się nim dymny sprawdzian pionu
(rozdział 19.2).

Wynik `channel.add` niesie obiekt `channel`, którego pole `id` zawiera **kod kanału** postaci
`kanal-<licznik>-<losowe>`. **Wartość tę należy zapamiętać** — to ona idzie do pola `modelChannelId`
przy zakładaniu okna oraz do `channelId` komend `channel.update` i `channel.remove`.

**Uwaga o wartości domyślnej klienta.** Interfejs kieruje okno domyślnie na wartość `cli` — rodzaj
kanału, nie kod wiersza (rozdział 10.3). Wiersz założony komendą `channel.add` dostaje kod
`kanal-…`, więc **domyślne wskazanie okna go nie odnajdzie**. Są dwie drogi rozwiązania:

- wskazać kod kanału jawnie przy zakładaniu okna (`modelChannelId` w `window.create`) — droga
  zalecana, bo nie wymaga ingerencji w bazę;
- albo zasiać wiersz z kodem dosłownie równym `cli` — droga z rozdziału 20.7, po której domyślne
  wskazanie klienta trafia w cel bez żadnej dalszej konfiguracji.

### 20.5 Krok 4 — zasianie konta

Krok wykonywany komendą kontraktu **`account.add`**. Treść żądania:

```json
{
  "name":        "konto-glowne",
  "kind":        "cli",
  "provider":    "anthropic",
  "configDir":   "C:\\Users\\Operator\\AppData\\Local\\DanacoConsole\\profile\\konto-glowne",
  "makeDefault": true,
  "enabled":     true
}
```

Cztery uwagi o polach, wszystkie potwierdzone w kodzie:

- **`name`** trafia do kolumny `nazwa` (`NOT NULL UNIQUE`) i **jest kodem konta w puli rotacji**.
  Pojawia się w prowenancji wywołania, więc warto nadać nazwę czytelną.
- **`kind` musi być `cli`** — pula rotacji kanału głównego bierze wyłącznie ten rodzaj.
- **`configDir` jest obowiązkowe w praktyce**: konto bez katalogu konfiguracji **nie wchodzi do puli**
  i jest z punktu widzenia wykonania niewidoczne.
- **`provider`** jest polem wymaganym kontraktu i wartością danych, nie typem kodu.

Pole `credential` istnieje w kontrakcie, lecz **nie należy go używać do przekazania sekretu w tej
procedurze**: poświadczenia kanału głównego mieszkają w katalogu konfiguracji, do którego prowadzi
`configDir`. Odpowiedź `account.add` nigdy nie niesie poświadczenia.

### 20.6 Krok 5 — kontrola zasiania i pierwsza tura

**Krok 5a — ponowne uruchomienie rdzenia. Obowiązkowe.**

Pula kont budowana jest **raz, przy starcie rdzenia** (rozdział 11.3). Konto założone komendą
`account.add` **nie trafia do puli pracującego procesu**. Bez restartu kanał główny nadal będzie
odmawiał wykonania, mimo poprawnie założonego konta — jest to najczęstsza pułapka całej procedury.

Rejestr kanałów zachowuje się odwrotnie: odświeża się po każdej komendzie `channel.*` i restartu nie
wymaga. Restart jest potrzebny **wyłącznie z powodu puli kont**.

```powershell
# zatrzymać proces rdzenia, po czym uruchomić ponownie
Set-Location C:\DanacoConsole_App
.\danaco-console.exe --profile "C:\Users\Operator\AppData\Local\DanacoConsole\profile"
```

**Krok 5b — kontrola rejestru.** Komenda `channel.list` z treścią `{"enabledOnly": true}` powinna
zwrócić co najmniej jeden kanał. Komenda `account.list` z treścią `{"kind": "cli"}` powinna zwrócić
założone konto wraz z polem `total`.

**Krok 5c — pierwsza tura kontrolna na kanale echo.** Zalecane przed próbą rozmowy z modelem
rzeczywistym, bo oddziela sprawdzenie produktu od sprawdzenia programu `claude` i poświadczeń.
Kolejność komend jest ta sama, którą przechodzi dymny sprawdzian pionu:

```
session.create  → { "title": "Kontrola instalacji" }              → sesja.id
window.create   → { "sessionId": "<sesja.id>",
                    "moduleId": "talkin.rozmowa",
                    "modelChannelId": "<kod kanału echo>",
                    "workingDirs": [], "executionEnv": "core",
                    "permissionMode": "manual",
                    "windowRole": "standalone",
                    "title": "Okno kontrolne" }                   → okno.id
message.send    → { "windowId": "<okno.id>",
                    "content": "kontrola instalacji",
                    "stream": true }
```

Oczekiwany przebieg: strumień zdarzeń `stream.chunk` z fragmentami tekstu, zakończony zdarzeniem
`message.changed` z wiadomością roli `assistant` w stanie `complete`. Treść odpowiedzi kanału echo
to przedrostek doklejony do treści zapytania.

Przejście tego kroku potwierdza **cały pion produktu**: transport, kontrakt, rejestr kanałów, sesję,
okno, strumień i zapis wiadomości do bazy — bez udziału sieci, konta i programu `claude`.

**Krok 5d — pierwsza tura na kanale głównym.** Ta sama sekwencja z kodem kanału `cli` zamiast echa.
Dopiero jej niepowodzenie wskazuje na program `claude`, konto albo poświadczenia — rozpoznanie
przyczyn opisuje rozdział 21.

Najkrótszą drogą sprawdzenia samego produktu, bez ręcznego wykonywania komend, pozostaje uruchomienie
dymnego sprawdzianu pionu na osobnym porcie i osobnym katalogu danych (rozdział 19.2). Nie zastępuje
on jednak procedury: sprawdzian pracuje na własnej bazie tymczasowej, a instalacja robocza pozostaje
niezasiana.

### 20.7 Wariant awaryjny — zapis wprost do bazy

Wariant przeznaczony na wypadek braku możliwości wykonania komend po WebSocket oraz — co ważniejsze —
na wypadek, gdy potrzebny jest **wiersz kanału o kodzie dosłownie równym `cli`**, w który trafia
domyślne wskazanie klienta. Kodu nie da się nadać komendą `channel.add`: rdzeń generuje go
automatycznie.

Warunki wykonania: **rdzeń musi być zatrzymany** (tryb WAL i blokady zapisu), a baza dostępna pod
ścieżką z rozdziału 9.1.

```sql
-- Kanał główny o kodzie 'cli' — trafia w domyślne wskazanie okna po stronie klienta.
INSERT INTO kanal_modelu
  (kod, nazwa, dostawca, identyfikator_modelu, rodzaj_kanalu, parametry_json, aktywny, kolejnosc)
VALUES
  ('cli', 'Kanał główny — Claude Code CLI', 'anthropic', '<identyfikator modelu>',
   'cli', '{}', 1, 1);

-- Konto rotacji kanału głównego.
INSERT INTO konto
  (nazwa, rodzaj, dostawca, katalog_konfiguracji, aktywne, kolejnosc, domyslne, stan)
VALUES
  ('konto-glowne', 'cli', 'anthropic',
   'C:\Users\Operator\AppData\Local\DanacoConsole\profile\konto-glowne',
   1, 1, 1, 'aktywne');
```

Wiersz kanału echo w tym samym wariancie:

```sql
INSERT INTO kanal_modelu
  (kod, nazwa, dostawca, identyfikator_modelu, rodzaj_kanalu, parametry_json, aktywny, kolejnosc)
VALUES
  ('echo', 'Echo kontrolne', 'lokalny', 'echo-1', 'lokalny',
   '{"adapter":"echo","porcja":16,"przedrostek":"echo: "}', 1, 2);
```

Reguły, których zapis wprost do bazy nie może naruszyć:

- `kanal_modelu.rodzaj_kanalu` **musi** należeć do zbioru `('cli','api','sdk','lokalny')`; wartość
  `echo` zostanie odrzucona (adapter wskazuje się parametrem `adapter`).
- `konto.rodzaj` **musi** należeć do zbioru `('cli','api','sdk')`.
- `konto.stan` **musi** należeć do zbioru `('aktywne','wyczerpane','zawieszone')`.
- Indeks częściowy dopuszcza **dokładnie jedno konto domyślne na rodzaj** — `domyslne = 1` dla
  drugiego konta rodzaju `cli` zostanie odrzucone.
- Kolumny `kod` (kanału) i `nazwa` (konta) są **unikalne**.

Po zapisie należy uruchomić rdzeń; rejestr kanałów i pula kont zbudują się z nowych wierszy przy
starcie.

Program `sqlite3` nie jest częścią produktu (rozdział 3.7) — wariant awaryjny wymaga zainstalowania go
osobno albo użycia innego narzędzia obsługującego SQLite.

### 20.8 Ograniczenie procedury

Procedura doprowadza instalację do stanu, w którym **okno komunikacji wykonuje tury modelu**. Nie
usuwa i nie może usunąć ograniczeń opisanych w rozdziałach 17 i 18. Po jej wykonaniu nadal
obowiązuje:

- **model nie pamięta poprzednich wiadomości okna** — każda tura startuje bez historii
  (rozdział 17.3);
- **historia znika z widoku po odświeżeniu strony**, mimo że jest zapisana w bazie;
- **zmiana konta wymaga ponownego uruchomienia rdzenia** (rozdział 11.3);
- **ustawienia modelu, nakładu i hosta w panelu okna nie mają odbiorcy** (rozdział 17.2);
- **moduły i środowiska pozostają atrapami nawigacyjnymi** (rozdział 18.2).

Procedura jest więc warunkiem koniecznym pracy, nie warunkiem wystarczającym do wdrożenia
produkcyjnego. **[DO DECYZJI OPERATORA]** — czy uruchamiać wdrożenie pilotażowe w obecnym zakresie
(jedno okno komunikacji, tury bez ciągłości), czy poczekać na zamknięcie ciągłości rozmowy
i odczytu historii.

---

## 21. Diagnostyka typowych problemów

### 21.1 Gdzie szukać śladów

Dostępność śladów **zależy od drogi uruchomienia** i to rozróżnienie jest pierwszą rzeczą, którą należy
ustalić przed jakąkolwiek diagnostyką:

- **rdzeń uruchomiony samodzielnie z konsoli** (`danaco-console.exe` wywołany z wiersza poleceń) —
  wyjście standardowe i diagnostyczne idzie do okna konsoli i **nie jest nigdzie utrwalane**. Po
  zamknięciu okna ślad przepada; sam rdzeń pliku dziennika nie zakłada;
- **rdzeń podniesiony przez powłokę natywną** — powłoka pracuje w podsystemie okienkowym, więc nie ma
  konsoli, do której mogłaby oddać wyjście procesu potomnego. Zamiast je porzucić, **zakłada plik
  dziennika i kieruje do niego oba strumienie rdzenia** — wyjście standardowe i diagnostyczne.

Plik nosi nazwę `rdzen-powloki.log` (stała `NAZWA_PLIKU` w `desktop/src-tauri/src/rdzen/dziennik.rs`)
i leży **w katalogu danych**, czyli obok pliku bazy `danaco-console.db`:

```
C:\Users\<Operator>\AppData\Local\DanacoConsole\rdzen-powloki.log
```

Funkcja `strumienie` otwiera ten plik w trybie dopisywania i zwraca go jako `Stdio` dla wyjścia
i błędu procesu rdzenia (użycie w `rdzen/uruchomienie.rs`, funkcja `postaw`). Ta sama powłoka dokłada
do pliku **własne wiersze** z przedrostkiem `powłoka:` przez funkcję `dopisz` — wołaną przy montażu
(`montaz.rs`) trzykrotnie: opis uchwytu rdzenia, opis źródła interfejsu oraz potwierdzenie otwarcia
okna i ustawienia ikony w zasobniku. Dziennik jest więc **wspólny dla powłoki i rdzenia**, co pozwala
czytać kolejność zdarzeń startu w jednym miejscu.

Pozostałe stwierdzenia o dzienniku pozostają w mocy: **nie ma rotacji** — plik rośnie bez ograniczenia
i bez podziału na pliki historyczne, a jego przycięcie należy do Operatora; **nie ma ustawienia poziomu
szczegółowości** — wzorzec `.env.example` wymienia „poziom dziennika” wśród ustawień świadomie
nieobecnych, bo dziennik sterowany ustawieniem jeszcze nie powstał (rozdział 12.3). Niepowodzenie
otwarcia pliku **nie wstrzymuje startu**: strumienie schodzą wtedy na urządzenie puste, a wyjście
rdzenia przepada — brak dziennika nie jest bramą.

Ślady rozłożone są w pięciu miejscach:

| Źródło | Co niesie | Jak sięgnąć |
| --- | --- | --- |
| **plik `<katalog danych>\rdzen-powloki.log`** | pełne wyjście standardowe i diagnostyczne rdzenia podniesionego przez powłokę oraz wiersze samej powłoki z przedrostkiem `powłoka:` | odczyt pliku w katalogu danych (domyślnie `%LOCALAPPDATA%\DanacoConsole`) dowolnym czytnikiem tekstu |
| **wyjście diagnostyczne rdzenia** (`stderr`) | wiersz startowy, komunikaty transportu, ostrzeżenia rejestrów, ślad pętli koordynator–wykonawca | uruchomić rdzeń z konsoli i czytać na bieżąco |
| **fragmenty błędu w strumieniu** | przyczyna niepowodzenia tury — trafia do okna komunikacji | okno komunikacji, wiadomość w stanie błędu |
| **tabela `wiadomosc`** | stan wiadomości (`complete`, `error`, `stopped`) i treść | odczyt bazy |
| **wyjście diagnostyczne procesu `claude`** | przyczyny po stronie programu zewnętrznego | trafia do bufora rdzenia i jest cytowane w rozpoznaniu wyczerpania limitu |

Ścieżki dziennika **nie trzeba zgadywać**. Powłoka wylicza ją funkcją `dziennik::sciezka` i podaje
w opisie uchwytu rdzenia — polu `sciezka_dziennika` struktury opisu (`rdzen/uruchomienie.rs`). Jest to
droga pewna także wtedy, gdy katalog danych został przesunięty zmienną `DANACO_KATALOG_DANYCH`
(rozdział 12.1): opis uchwytu wskazuje miejsce faktyczne, nie domyślne.

Wniosek praktyczny jest odwrotny do intuicji: **przy pracy z powłoką natywną należy czytać plik
`rdzen-powloki.log`, a przy uruchomieniu ręcznym — okno konsoli**. To start z powłoki utrwala wyjście,
a start z konsoli go nie zachowuje. Jeżeli więc usterka wystąpiła w codziennej pracy z powłoką, ślad
najpewniej już istnieje na dysku i wystarczy go otworzyć; jeżeli rdzeń był uruchamiany ręcznie,
odtworzenie usterki wymaga powtórzenia jej przy otwartej konsoli.

**[DO DECYZJI OPERATORA]** — czy przy uruchomieniu ręcznym przekierowywać `stderr` rdzenia do pliku
skrótem uruchomieniowym, oraz jaki przyjąć tryb przycinania `rdzen-powloki.log` (produkt nie prowadzi
rotacji, więc plik rośnie bez ograniczenia).

### 21.2 Pusty rejestr kanałów

**Objaw.** Wysłanie wiadomości kończy się fragmentem błędu; w treści pojawia się
`models: kanał "…" nie istnieje w rejestrze albo jest nieczynny` albo
`models: okno … bez wskazanego kanału modelu`. Odczyt `channel.list` zwraca pustą listę.

**Przyczyny, w kolejności prawdopodobieństwa:**

1. **Procedura z rozdziału 20 nie została wykonana.** Świeża instalacja ma rejestr pusty — to stan
   domyślny, nie usterka.
2. **Wiersz istnieje, lecz okno wskazuje inną wartość.** Domyślne wskazanie klienta to literał `cli`
   (rodzaj), a wiersz założony komendą `channel.add` ma kod `kanal-…`. Rozwiązanie: wskazać kod jawnie
   przy zakładaniu okna albo zasiać wiersz o kodzie `cli` (rozdział 20.7).
3. **Wiersz istnieje i jest nieczynny** (`aktywny = 0`) — do rejestru czynnych nie trafia.
4. **Wiersz istnieje, lecz nie ma dla niego fabryki adaptera.** Dotyczy rodzaju `sdk` bez parametru
   `adapter` oraz parametru `adapter` o wartości nieznanej. Wiersz trafia wtedy na listę pominiętych
   z powodem `brak fabryki adaptera <klucz>`.

**Rozpoznanie.** Odczyt wprost z bazy rozstrzyga wszystkie cztery przypadki naraz:

```sql
SELECT id, kod, rodzaj_kanalu, aktywny, parametry_json FROM kanal_modelu ORDER BY kolejnosc, id;
```

### 21.3 Pusta pula kont — komunikat mylący

**Objaw.** Tura na kanale głównym kończy się fragmentem błędu o treści:

```
wszystkie konta puli mają wyczerpany limit; kolejne próby pozostają otwarte
```

**Ostrzeżenie o mylącym komunikacie.** Ten sam komunikat pojawia się w **dwóch zupełnie różnych
sytuacjach**, ponieważ obie kończą się brakiem konta do podania przez pulę:

- pula zawiera konta i **wszystkie mają czynne wyczerpanie limitu** — komunikat jest wtedy prawdziwy;
- **pula jest pusta** (zero kont) — komunikat jest wtedy **fałszywy**: żaden limit nie został
  wyczerpany, po prostu nie ma czego rotować.

Na świeżej instalacji obowiązuje niemal zawsze przypadek drugi. Komunikat nie powinien być brany
dosłownie i **nie oznacza problemu z limitami u dostawcy**.

**Rozpoznanie.** Odczyt puli z bazy — dokładnie tym warunkiem, którym posługuje się rdzeń:

```sql
SELECT id, nazwa, rodzaj, aktywne, stan, katalog_konfiguracji, domyslne, kolejnosc
  FROM konto
 WHERE rodzaj = 'cli' AND aktywne = 1 AND stan <> 'zawieszone'
 ORDER BY domyslne DESC, kolejnosc, id;
```

Pusty wynik oznacza pustą pulę. Wynik niepusty, w którym **wszystkie wiersze mają pustą kolumnę
`katalog_konfiguracji`**, znaczy to samo: takie konta do puli nie wchodzą (rozdział 11.3).

**Rozwiązanie:** wykonać kroki 4 i 5a procedury z rozdziału 20 — założyć konto **i ponownie uruchomić
rdzeń**. Konto założone przy pracującym rdzeniu nie trafia do puli.

### 21.4 Brak programu `claude`

**Objaw.** Tura kończy się fragmentem błędu z treścią zaczynającą się od
`injection: uruchomienie claude:` i komunikatem systemu o nieodnalezieniu pliku.

**Przyczyny i rozwiązania:**

| Przyczyna | Rozpoznanie | Rozwiązanie |
| --- | --- | --- |
| program nie jest zainstalowany | `claude --version` w konsoli | zainstalować Claude Code CLI |
| program jest, lecz poza `PATH` procesu rdzenia | `claude --version` działa w konsoli, a tura nadal pada | podać ścieżkę parametrem `program` wiersza kanału |
| `PATH` zmieniono po starcie rdzenia | — | **ponownie uruchomić rdzeń** — środowisko dziedziczone jest z chwili startu (rozdz. 12.4) |
| ustawiono `harness.program_claude` licząc, że zadziała | wartość jest w bazie, wykonanie jej nie czyta | użyć parametru `program` wiersza kanału (rozdz. 12.3) |

Zmiana ścieżki programu w wierszu kanału:

```json
{ "channelId": "<kod kanału>", "config": { "program": "C:\\ścieżka\\do\\claude.exe" } }
```

Komenda `channel.update` **nadpisuje całe `parametry_json`**, gdy pole `config` jest podane — należy
więc przesłać komplet parametrów, nie sam zmieniany klucz.

### 21.5 Port zajęty

**Objaw.** Rdzeń przerywa start komunikatem `transport: nasłuch <adres>:` wraz z przyczyną systemową
o zajętym adresie.

**Najczęstsza przyczyna:** działa już inny proces `danaco-console.exe` — podniesiony
wcześniej przez powłokę natywną w tle. Powłoka stawia rdzeń samodzielnie, więc uruchomienie rdzenia
ręcznie **przy otwartej powłoce** kończy się konfliktem.

**Rozpoznanie na Windows:**

```powershell
Get-NetTCPConnection -LocalPort 17870 -State Listen |
  Select-Object -Property LocalAddress, LocalPort, OwningProcess
Get-Process -Id <OwningProcess>
```

**Rozwiązania:** zatrzymać zbędny proces rdzenia albo — jeżeli konflikt pochodzi od zupełnie innej
usługi — zmienić port, pamiętając o **podwójnej zmianie** opisanej w rozdziale 8.4 (rdzeń **i** klient).

**Uwaga o zakresie nasłuchu.** Pusty adres nasłuchu znaczy pętlę zwrotną (`127.0.0.1`,
`transport/ustawienia.go`); nasłuch na wszystkich interfejsach wymaga jawnego
`DANACO_WSZYSTKIE_INTERFEJSY=1` i włącza wymóg logowania sam z siebie. Wdrożenie z rozstrzygnięcia 29
nasłuchu nie poszerza: rdzeń stoi na `127.0.0.1:17870` z jawnym wymogiem logowania, a na świat wystawia
go Caddy (rozdział 2.4). Reguła zapory dla portu 17870 nie jest potrzebna.

### 21.6 Powłoka bez pakietu interfejsu — „pusta strona”

**Objaw.** Rdzeń pracuje, adres `http://127.0.0.1:17870/` odpowiada, lecz nie ma czego wyświetlić;
albo okno powłoki natywnej pozostaje puste.

**Przyczyna.** Brak katalogu pakietu interfejsu. Rdzeń **nie przerywa pracy** przy jego braku —
gniazdo transportu działa normalnie, brakuje wyłącznie plików statycznych. Domyślna ścieżka
`client\dist` jest **względna**, więc wystarczy uruchomić rdzeń z innego katalogu roboczego, by
przestał ją znajdować.

**Rozpoznanie i rozwiązanie:**

1. Sprawdzić, czy w katalogu produktu istnieje `client\dist\index.html`. Brak — wykonać wydanie
   (rozdział 5.2); skrypt pakowania odmawia pracy bez tego pliku.
2. Sprawdzić katalog roboczy procesu. Uruchamiać rdzeń **z katalogu produktu** (rozdział 7.1) albo
   wskazać pakiet jawnie: `--klient "C:\DanacoConsole_App\client\dist"` bądź zmienną
   `DANACO_KATALOG_KLIENTA`.
3. Przy powłoce natywnej — sprawdzić, czy `danaco-console.exe` leży **obok** `Danaco Console.exe`,
   a katalog `client\dist` obok binarki rdzenia (rozdział 6.2). Powłoka szuka obu w tych miejscach.

Osobny przypadek: **strona ładuje się, lecz interfejs nie łączy się z rdzeniem**. To nie jest brak
pakietu, lecz rozjazd adresu gniazda — rozdziały 8.4 i 8.5.

### 21.7 Baza zablokowana albo niezgodna

**Objaw A — blokada.** Operacja na bazie kończy się błędem o zajętej bazie mimo ustawionego
`busy_timeout` (5 sekund).

**Przyczyna.** Do jednego pliku bazy sięga **więcej niż jeden proces rdzenia** — najczęściej rdzeń
uruchomiony ręcznie i rdzeń podniesiony przez powłokę. Tryb WAL pozwala czytać podczas zapisu, lecz
nie znosi konfliktu dwóch pisarzy.

**Rozwiązanie.** Pozostawić jeden proces rdzenia. Rozpoznanie jak w rozdziale 21.5.

**Objaw B — zmieniona migracja.** Start przerywa komunikat
`store: migracja NNN (nazwa) zmieniła treść po zastosowaniu`.

**Przyczyna.** Suma kontrolna kroku migracji w rejestrze nie zgadza się z sumą kroku w binarce.
Zdarza się przy sięgnięciu bazą po rdzeń z innego wydania albo po edycji wydanej migracji.

**Rozwiązanie.** Użyć rdzenia z tego samego wydania, z którego pochodzi baza. Migracji **nie da się
cofnąć** — mechanizm nie posiada migracji odwrotnych (rozdział 9.4).

**Objaw C — niezgodność kontraktu.** Interfejs łączy się, lecz komendy kończą się błędami.
Kontrola: pole `protocolVersion` w wyniku `connection.hello` po obu stronach (rozdział 19.4).
Rozbieżność oznacza klienta i rdzeń z różnych wydań — należy zbudować obie części z tego samego
drzewa.

### 21.8 Historia znika po odświeżeniu strony

**Objaw.** Rozmowa jest widoczna do czasu odświeżenia interfejsu; po przeładowaniu okno jest puste,
mimo że wiadomości były wysłane i doczekały odpowiedzi.

**Przyczyna.** **Nie jest to utrata danych.** Wiadomości są utrwalone w tabeli `wiadomosc`, a sesje
i okna odtwarzają się z bazy po restarcie rdzenia. Brakuje wyłącznie odczytu po stronie interfejsu:
klient **nigdy nie wywołuje `message.list`**. **[NIEZINTEGROWANE]**

**Potwierdzenie, że dane są w bazie:**

```sql
SELECT COUNT(*) FROM wiadomosc;
```

**Rozwiązanie w bieżącej wersji: brak.** Do czasu spięcia odczytu historii jedyną drogą do treści
poprzednich rozmów jest odczyt bazy. Nie należy tego mylić z brakiem ciągłości rozmowy po stronie
modelu, który jest osobnym ograniczeniem (rozdział 17.3).

---

## 22. Aktualizacja, ponowne uruchomienie i deinstalacja

### 22.1 Co przeżywa ponowne uruchomienie rdzenia

Rozgraniczenie tego, co trwałe, od tego, co ulotne, jest niezbędne przy planowaniu aktualizacji.

**Przeżywa restart** (leży w bazie, katalog danych):

- sesje, karty sesji i okna komunikacji — rdzeń odtwarza je z bazy przy starcie;
- wiadomości wraz ze stanem i treścią (tabela `wiadomosc`);
- rejestr kanałów modelu (`kanal_modelu`);
- konta (`konto`) wraz z kolejnością i wskazaniem konta domyślnego;
- punkty dostępu, nadania, korzenie i argumenty trybów mostu;
- ustawienia okna konfiguracji (`ustawienie`) na wszystkich poziomach i osiach;
- dokumenty tożsamości i katalog kategorii;
- rejestr zastosowanych migracji.

**Nie przeżywa restartu** (żyje w pamięci procesu):

- **stan wyczerpania kont w puli rotacji** — po restarcie pula zaczyna od konta pierwszego
  (rozdział 11.3);
- **bieżący wskaźnik rotacji** — który konto jest „bieżące”;
- **biegnące tury** — proces programu `claude` kończy się razem z rdzeniem;
- **stan pętli koordynator–wykonawca** (licznik obiegów, powody zatrzymania);
- **połączenia WebSocket** — klienci muszą połączyć się ponownie.

**Nie przeżywa odświeżenia strony po stronie interfejsu:**

- **widok historii rozmowy** — dane pozostają w bazie, lecz klient ich nie odczytuje
  (rozdziały 17.3, 21.8);
- stan sceny okien nietrafiony do bazy — `window.update` nie zapisuje.

Osobno warto pamiętać, że **zmiany w tabeli `konto` odnoszą skutek dopiero po restarcie rdzenia**,
podczas gdy zmiany w tabeli `kanal_modelu` działają natychmiast. Jest to jedyne miejsce, w którym
restart jest częścią normalnej procedury konfiguracyjnej, a nie reakcją na awarię.

### 22.2 Ponowne wydanie i aktualizacja

Aktualizacja produktu sprowadza się do ponownego wydania ze zaktualizowanych źródeł:

```bash
# w katalogu repozytorium
git pull
bash budowa/scripts/wydanie.sh
```

Co skrypt zrobi z katalogiem produktu (rozdział 5.3):

- **nadpisze** `danaco-console.exe` i — w wariancie pełnym — `Danaco Console.exe`;
- **wyczyści i odtworzy** `client\dist`: usunie z podkatalogu `client\` wszystkie pliki poza `*.md`,
  usunie opróżnione katalogi i skopiuje nowy pakiet;
- **nie tknie plików `*.md`** w katalogu produktu — dokumentacja przeżywa każde wydanie;
- **nie tknie katalogu danych** — jest poza katalogiem produktu.

Kolejność czynności przy aktualizacji instalacji roboczej:

1. **Zatrzymać rdzeń i powłokę.** Nadpisanie pracującego pliku wykonywalnego nie powiedzie się na
   Windows.
2. **Wykonać kopię katalogu danych** (rozdział 9.4). Migracje aplikują się samoczynnie i **nie da się
   ich cofnąć** — kopia jest jedyną drogą powrotu do poprzedniego wydania.
3. Uruchomić wydanie.
4. Uruchomić rdzeń i sprawdzić wiersz startowy oraz wersję schematu (rozdział 19.3).
5. Sprawdzić zgodność `protocolVersion` (rozdział 19.4) — po aktualizacji obie części pochodzą
   z tego samego drzewa, więc rozbieżność wskazywałaby na niekompletne wydanie.

Wymuszenie ponownej instalacji zależności klienta wymaga usunięcia `budowa/client/node_modules` przed
wydaniem — skrypt instaluje je **wyłącznie wtedy, gdy katalogu nie ma** (rozdział 3.3).

**Uwaga o zgodności wstecznej.** Baza zmigrowana do wyższej wersji schematu **nie zadziała ze starszym
rdzeniem**: mechanizm migracji nie ma kroków odwrotnych, a niezgodność sumy kontrolnej zatrzyma start
(rozdział 9.2). Powrót do poprzedniego wydania wymaga odtworzenia katalogu danych z kopii.

### 22.3 Ponowne uruchomienie

Rdzeń zatrzymuje się na sygnał przerwania (Ctrl+C w konsoli) albo na sygnał zakończenia procesu.
Zamknięcie jest uporządkowane: rdzeń zamyka kanały, zwalnia zasoby adapterów i **zamyka bazę**, co
scala dziennik WAL z plikiem głównym. Wypisuje przy tym wiersz `zatrzymanie: rdzeń zamknięty`.

Zabicie procesu bez sygnału, poleceniem `Stop-Process -Force`, pozostawia pliki `-wal` i `-shm`.
Nie jest to uszkodzenie bazy — SQLite odtworzy stan przy następnym otwarciu — lecz kopia zapasowa
wykonana w takiej chwili wymaga skopiowania kompletu plików.

Przy pracy z powłoką natywną rdzeń jest procesem potomnym powłoki: zamknięcie okna powłoki zamyka
też rdzeń. Uruchamianie rdzenia ręcznie przy otwartej powłoce prowadzi do konfliktu portu
(rozdział 21.5).

### 22.4 Deinstalacja

Danaco Console nie posiada deinstalatora, bo nie posiada instalatora (rozdział 5.1). Usunięcie
sprowadza się do skasowania katalogów, przy czym **kolejność i świadomość zawartości mają znaczenie**:

| Krok | Co usunąć | Skutek |
| --- | --- | --- |
| 1 | zatrzymać rdzeń i powłokę | zwolnienie plików i bazy |
| 2 | `C:\DanacoConsole_App` | usunięcie produktu **wraz z dokumentacją** — dokumenty `*.md` należy wcześniej zabezpieczyć |
| 3 | `%LOCALAPPDATA%\DanacoConsole` | **nieodwracalne usunięcie wszystkich danych**: sesji, wiadomości, kont, kanałów, dostępów, ustawień i tożsamości, a wraz z nimi pliku dziennika rdzenia `rdzen-powloki.log` (rozdział 21.1) |
| 4 | katalogi profili kont (`CLAUDE_CONFIG_DIR`) | usunięcie poświadczeń programu `claude` — leżą poza produktem |
| 5 | `C:\DanacoConsole_Git` | usunięcie źródeł, jeżeli nie są dalej potrzebne |
| 6 | reguła zapory dla `danaco-console.exe` | porządek w konfiguracji systemu |

Krok 2 bez kroku 3 daje **odinstalowanie z zachowaniem danych** — stan przydatny przy przenoszeniu
instalacji na nową stację albo przy czasowym usunięciu produktu. Kroki 3 i 4 są nieodwracalne
i wymagają świadomej decyzji.

Produkt **nie zapisuje niczego w rejestrze systemu Windows**, nie zakłada usług systemowych i nie
umieszcza plików poza wymienionymi katalogami. Wyjątkiem jest ewentualna reguła zapory utworzona przy
pierwszym uruchomieniu (rozdział 4.4) oraz — w razie budowy instalatora NSIS — wpisy pochodzące
z niego, których niniejszy dokument nie opisuje, bo droga NSIS nie jest zalecana (rozdział 5.6).

---

## 23. Załączniki

### 23.1 Wykaz plików konfiguracyjnych i ich rola

Danaco Console **nie posiada pliku konfiguracyjnego rdzenia**. Poniższy wykaz obejmuje wszystkie pliki
mające jakikolwiek wpływ na konfigurację produktu — wraz z jednoznacznym wskazaniem, czy są czytane
w czasie pracy, czy wyłącznie w czasie budowy.

| Plik | Kiedy czytany | Rola |
| --- | --- | --- |
| `budowa/.env.example` | **nigdy przez rdzeń** | wzorzec dokumentacyjny zmiennych środowiskowych (rozdz. 12.5) |
| `budowa/.env` | **nigdy przez rdzeń** | plik Operatora poza kontrolą wersji; zmienne trzeba ustawić w środowisku procesu |
| `budowa/go.mod` | budowa rdzenia | wersja Go i zależności rdzenia |
| `budowa/client/package.json` | budowa klienta | zależności i polecenia klienta |
| `budowa/client/.env*` (jeśli utworzony) | **budowa klienta** | zmienne `VITE_*` wkompilowywane w pakiet (rozdz. 12.5) |
| `budowa/desktop/src-tauri/tauri.conf.json` | budowa powłoki | `frontendDist`, `devUrl`, `beforeBuildCommand`, pakowanie NSIS |
| `budowa/desktop/src-tauri/Cargo.toml` | budowa powłoki | wersja Rust, cechy Tauri, profil wydania |
| `budowa/shared/contract.json` | generowanie kontraktu | jedno źródło nazw komend, zdarzeń i narzędzi |
| `%LOCALAPPDATA%\DanacoConsole\danaco-console.db` | **każdy start i praca rdzenia** | **jedyne miejsce konfiguracji dziedzinowej**: kanały, konta, dostępy, ustawienia, tożsamość |
| `/etc/danaco-console/sejf.klucz` (serwer) | pierwsze użycie sejfu po starcie | klucz sejfu poświadczeń wskazany `DANACO_KLUCZ_SEJFU`; `postinst` zakłada go z prawami `0400`, gdy go nie ma (rozdz. 12.6) |
| `/etc/danaco-console/srodowisko` (serwer) | start jednostki systemd | nastawy Operatora czytane `EnvironmentFile`: TLS, `DANACO_SEKRET_NAWIAZANIA`, konto nadawcze (rozdz. 12.6) |
| katalogi `CLAUDE_CONFIG_DIR` kont | każde wywołanie kanału `cli` | poświadczenia programu `claude`; poza produktem i poza bazą |

Wniosek: **konfiguracja pracującego produktu mieszka w bazie, nie w plikach.** Pliki wpływają na
konfigurację wyłącznie w chwili budowy albo jako wzorce dokumentacyjne.

### 23.2 Komendy kontraktu przywołane w dokumencie

Wszystkie nazwy pochodzą z kontraktu (`budowa/shared`). Tabela obejmuje wyłącznie komendy i zdarzenia
przywołane w niniejszym dokumencie — nie jest pełnym katalogiem kontraktu.

| Nazwa | Rola | Rozdział | Stan wg dokumentu |
| --- | --- | --- | --- |
| `connection.hello` | powitanie klienta; ustala wersję protokołu | 19.4 | **[DZIAŁA]** |
| `channel.add` | dopisuje wiersz rejestru kanałów | 10.2, 20.4 | **[DZIAŁA]** — brak UI |
| `channel.update` | zmienia wiersz rejestru (po **kodzie**) | 10.2, 21.4 | **[DZIAŁA]** |
| `channel.remove` | usuwa wiersz rejestru (po **kodzie**) | 10.2 | **[DZIAŁA]** |
| `channel.list` | zwraca rejestr kanałów (pole `id` = **numer wiersza**) | 10.2, 19.4 | **[DZIAŁA]** |
| `account.add` | zakłada konto | 20.5 | **[DZIAŁA]** — pula wymaga restartu |
| `account.list` | wykaz kont; poświadczeń nie niesie | 20.6 | **[DZIAŁA]** |
| `account.update`, `account.remove`, `account.default.set` | prowadzenie katalogu kont | 11.1 | **[DZIAŁA]** |
| `access.point.add/update/remove/list/check` | punkty dostępu | 13.3, 15.2 | **[DZIAŁA CZĘŚCIOWO]** — bez `localDirectory` |
| `access.grant.add/update/remove/list` | nadania okna | 13.3, 15.3 | **[DZIAŁA]** |
| `identity.category.list` | katalog kategorii tożsamości | 11.5 | **[DZIAŁA]** |
| `identity.document.get/set/remove` | dokumenty tożsamości per oś | 11.5 | **[DZIAŁA]** |
| `identity.effective.get` | nakładka obowiązująca | 11.5 | **[DZIAŁA]** |
| `config.get`, `config.set`, `config.reset` | ustawienia na ośmiu poziomach i trzech osiach | 17, 18.1 | **[DZIAŁA]** — część kluczy bez odbiorcy |
| `session.create` | zakłada sesję | 20.6 | **[DZIAŁA CZĘŚCIOWO]** — wiersz do pierwszej wiadomości |
| `window.create` | zakłada okno komunikacji | 20.6 | **[DZIAŁA CZĘŚCIOWO]** |
| `window.update` | zmienia okno | 17.3 | **[NIEZINTEGROWANE]** — nie zapisuje do bazy |
| `window.state.get` | stan okna | 18.2 | **[NIEZINTEGROWANE]** — zawsze stan oczekiwania |
| `message.send` | wysyła wiadomość i uruchamia turę | 20.6 | **[DZIAŁA]** |
| `message.stop` | zatrzymuje turę | 17.3 | **[DZIAŁA]** |
| `message.list` | odczyt historii | 17.3, 21.8 | **[NIEZINTEGROWANE]** — klient nie wywołuje |
| `queue.create`, `queue.action` | kolejki | 18.2 | **[NIEZINTEGROWANE]** |
| `session.focus`, `session.bind` | powrót do sesji trwającej na rdzeniu | 7.5, 18.2 | **[DZIAŁA]** — wołane przez stronę główną (strefa „Sesje w tle”) |
| `home.enter`, `environment.list/enter`, `module.list`, `workspace.enter` | nawigacja platformy | 18.2 | **[NIEZINTEGROWANE]** — niewywoływane przez widoki |
| `stream.chunk` (zdarzenie) | fragmenty strumienia odpowiedzi | 19.2, 20.6 | **[DZIAŁA]** |
| `message.changed` (zdarzenie) | zmiana stanu wiadomości | 19.2, 20.6 | **[DZIAŁA]** |
| `session.list` | wykaz sesji (z żywym stanem) | 7.5, 19.4 | **[DZIAŁA]** — pulpit i strona główna |
| `progress.changed`, `session.changed`, `window.changed` (zdarzenia) | telemetria postępu i stanu | 7.5, 18.2 | **[DZIAŁA CZĘŚCIOWO]** — konsument w pulpicie i w stronie głównej; poza nimi bez subskrybenta |

Koperta transportu ma niezmiennie kształt `{ type, id, payload, timestamp }`; odpowiedź niesie ten sam
`type` i `id` oraz pole `status` o wartości `ok` albo `error`.

### 23.3 Tabele bazy istotne dla instalacji

Wykaz obejmuje tabele, których Operator dotyka przy instalacji, konfiguracji albo diagnostyce.
Nie jest pełnym modelem danych produktu.

| Tabela | Migracja | Rola przy instalacji | Rozdział |
| --- | --- | --- | --- |
| `migracja` | zakładana z kodu | rejestr zastosowanych kroków schematu wraz z sumami kontrolnymi; źródło wersji schematu | 9.2, 19.3 |
| `kanal_modelu` | 001 | **rejestr kanałów modelu — wymaga zasiania** | 10, 20.4 |
| `konto` | 001 + 014 | **katalog kont i pula rotacji — wymaga zasiania** | 11, 20.5 |
| `urzadzenie` | 001 | wymagana dla punktu `localDirectory`; **brak repozytorium i komend** | 13.3 |
| `okno_komunikacji` | 002 / 005 | okna komunikacji wraz z zasięgiem wykonania i trybem uprawnień | 17 |
| `wiadomosc` | 001 | historia rozmowy; łańcuch sesja → karta → okno → wiadomość | 17.3, 21.8 |
| `ustawienie` | 001 + 012 | wartości ustawień pod adresem złożonym (poziom + oś) | 17, 18.1 |
| `definicja_ustawienia` | 012 | katalog definicji okna konfiguracji (klucze, rodzaje, wartości domyślne) | 17.2, 18.1 |
| `poziom_zasiegu` | 001 | osiem poziomów zasięgu konfiguracji | 18.1 |
| `punkt_dostepu` | 013 | punkty dostępu (`mcpBridge`, `localDirectory`) | 13.3, 15.2 |
| `korzen_punktu_dostepu` | 013 | korzenie udostępniane przez punkt | 13.3 |
| `nadanie_dostepu` | 013 | nadania punktu dla okna wraz z trybem | 13.3, 15.2 |
| `korzen_nadania` | 013 | korzenie nadania | 13.3 |
| `argument_trybu_mostu` | 013 | słownictwo trybu konkretnego mostu MCP | 15.2 |
| `akcja` | 008 + 009 | katalog akcji sterowany danymi | 18.1 |

Warunki `CHECK` istotne przy zapisie ręcznym (rozdział 20.7):

| Kolumna | Dopuszczalne wartości |
| --- | --- |
| `kanal_modelu.rodzaj_kanalu` | `cli`, `api`, `sdk`, `lokalny` |
| `konto.rodzaj` | `cli`, `api`, `sdk` |
| `konto.stan` | `aktywne`, `wyczerpane`, `zawieszone` |
| `punkt_dostepu.rodzaj` | `mcpBridge`, `localDirectory` |
| `punkt_dostepu.tryb_domyslny`, `nadanie_dostepu.tryb` | `read`, `write` |
| `punkt_dostepu.stan` | `unknown`, `reachable`, `unreachable` |
| `okno_komunikacji.srodowisko_wykonania` | `lokalne`, `rdzen`, `zdalny` |
| `poziom_zasiegu.kod` | `globalny`, `srodowisko`, `modul`, `para_modulow`, `projekt`, `karta_sesji`, `rola`, `okno` |

### 23.4 Wykaz kwestii pozostawionych do decyzji Operatora

Zestawienie wszystkich miejsc oznaczonych w dokumencie znacznikiem **[DO DECYZJI OPERATORA]**.
Żadna z tych kwestii nie blokuje instalacji; wszystkie wymagają rozstrzygnięcia przed wdrożeniem
produkcyjnym.

| Nr | Rozdział | Kwestia |
| --- | --- | --- |
| 1 | 2.4 | rozstrzygnięte: rozstrzygnięcie 29 — nasłuch `127.0.0.1:17870` z wymogiem logowania, na świat Caddy `console.danaco-group.pl:443` |
| 2 | 3.3 | Ustalenie wiążącej wersji Node.js (`engines` albo `.nvmrc`) dla powtarzalnej budowy |
| 3 | 4.4 | rozstrzygnięte: rozstrzygnięcie 27 — OV od Certum, klucz w SimplySign; do zakupu `DANACO_PODPIS=pomijany`. Otwarte: wyłączenia z monitorowania na stacji budującej |
| 4 | 5.6 | rozstrzygnięte: rozstrzygnięcia 8 i 27 — instalator pełny nie powstaje (rdzeń na serwerze), podpis jak w pozycji 3 |
| 5 | 7.4 | Czy wprowadzić jawne ustawienie adresu rdzenia po stronie klienta |
| 6 | 8.4 | Czy przekazywać port z rdzenia do klienta, by usunąć konieczność podwójnej zmiany |
| 7 | 9.4 | Czy ustanowić cykliczną kopię katalogu danych i jaki okres przechowywania przyjąć |
| 8 | 11.2 | Gdzie umieścić katalogi konfiguracji kont i jak zabezpieczyć je uprawnieniami |
| 9 | 12.5 | Czy wprowadzić wczytywanie pliku `.env`, czy pozostać przy zmiennych środowiska procesu |
| 10 | 13.3 | Czy udostępniać katalogi lokalne przez most MCP, czy czekać na drogę kontraktową urządzeń |
| 11 | 14.3 | Czy i kiedy wprowadzić transport sieciowy agenta wraz z uwierzytelnianiem |
| 12 | 15.4 | Czy rozprowadzić klucze hostów mostów do `known_hosts` przed pierwszym połączeniem |
| 13 | 16.3 | Czy prowadzić ewidencję kont `api` mimo braku drogi wykonawczej |
| 14 | 17.2 | Czy uzgodnić słownik kluczy ustawień okna po obu stronach, czy usunąć sterowania z interfejsu |
| 15 | 17.3 | Czy pracować w konwencji „jedna tura = jedno kompletne polecenie”, czy czekać na ciągłość rozmowy |
| 16 | 20.2 | Liczba kont w puli rotacji |
| 17 | 20.8 | Czy uruchamiać wdrożenie pilotażowe w obecnym zakresie |
| 18 | 21.1 | Czy przekierowywać wyjście diagnostyczne rdzenia do pliku |
| 19 | 21.5 | rozstrzygnięte: rozstrzygnięcie 29 — nasłuch stoi na pętli zwrotnej z jednostki, reguła zapory zbędna |

### 23.5 Słownik pojęć użytych w dokumencie

**Danaco Console (Platforma AI Workspace OS)** — produkt opisywany dokumentem; środowisko operacyjne
pracy z modelami językowymi instalowane na stacji Operatora.

**Operator** — osoba instalująca, konfigurująca i utrzymująca produkt; adresat dokumentu.

**rdzeń** — serwer aplikacyjny w języku Go (`danaco-console.exe`); jedyny proces dotykający bazy,
uruchamiający procesy kanału i nasłuchujący transportu.

**klient / interfejs** — warstwa interfejsu użytkownika w TypeScripcie, budowana do katalogu
`client/dist`; bezstanowa wobec trwałości.

**powłoka** — okno natywne systemu zbudowane w Rust/Tauri (`Danaco Console.exe`); nie zawiera logiki
dziedzinowej, stawia rdzeń w tle i ładuje interfejs.

**kontrakt** — jedno źródło nazw komend, zdarzeń i narzędzi (`budowa/shared`), z którego generowane są
`contract.ts` i `contract.go`; nazwy przywoływane w dokumencie pochodzą stamtąd.

**koperta** — jednostka transportu o kształcie `{ type, id, payload, timestamp }`.

**kanał modelu** — droga rozmowy z modelem; **wiersz** tabeli `kanal_modelu`, nie typ w kodzie.
Rodzaje dopuszczone przez bazę: `cli`, `api`, `sdk`, `lokalny`.

**rejestr kanałów** — zbiór adapterów zbudowanych w czasie działania z wierszy tabeli `kanal_modelu`;
odświeżany po każdej komendzie `channel.*`.

**adapter** — wykonawca kanału; wskazywany parametrem `adapter` wiersza albo kolumną `rodzaj_kanalu`.
Wbudowane: `echo`, `api`; kanał główny `cli` wnosi pakiet wstrzyknięcia.

**konto** — profil uwierzytelnienia kanału; wiersz tabeli `konto`. Kodem konta w puli jest kolumna
`nazwa`.

**pula kont** — uporządkowany zbiór kont rodzaju `cli` używany do rotacji; **budowany raz, przy starcie
rdzenia**.

**rotacja** — przejście na kolejne konto po rozpoznaniu wyczerpania limitu; nie przerywa sesji.

**nakładka tożsamości** — prompt systemowy złożony z warstw konstytucja → profil → ekspertyza,
podawany programowi w trybie `ZASTAP` (`--system-prompt`) albo `DOLACZ` (`--append-system-prompt`).

**prowenancja** — opis warunków wywołania (kanał, adapter, model, konto, adres, argumenty, zasięg
wykonania, skrót nakładki) nadawany jako pierwszy fragment strumienia, przed treścią odpowiedzi.

**okno komunikacji** — byt pośredni między sesją a wiadomością; niesie kanał, katalogi robocze, tryb
uprawnień, zasięg wykonania i rolę. **Najwęższy poziom zasięgu konfiguracji.**

**okno operacyjne** — pojęcie koncepcyjne modułów; katalog okien operacyjnych w bazie liczy siedemdziesiąt dziewięć wierszy, a w v2.0 istnieje **jeden widok** — okno komunikacji, któremu odpowiada piętnaście wierszy katalogu, po jednym na moduł (rozdział 18.2 oraz [INSTRUKCJA-UZYTKOWANIA.md](INSTRUKCJA-UZYTKOWANIA.md) rozdz. 7).

**sesja / karta sesji** — nadrzędne byty grupujące okna i wiadomości.

**moduł** — jednostka funkcjonalna platformy; w v2.0 pozycje nawigacji bez widoków.

**środowiska** (TalkIn, WorkSpace, CodeStudio, MultitaskingAI) — cztery przestrzenie produktu;
w v2.0 nawigacja bez realizacji funkcjonalnej.

**punkt dostępu** — wpis mówiący, **do czego** model ma wgląd (`mcpBridge` albo `localDirectory`).

**nadanie** — wiązanie punktu dostępu z jednym oknem, wraz z trybem i korzeniami.

**most MCP** — narzędzie udostępniane modelowi; konfiguracja `mcpServers` przekazywana programowi
przełącznikiem `--mcp-config`.

**katalog danych** — katalog trwałości produktu (`%LOCALAPPDATA%\DanacoConsole`), zawierający jedyny
plik bazy oraz — przy pracy z powłoką natywną — plik dziennika `rdzen-powloki.log` (rozdział 21.1).

**katalog profili** — katalog nadrzędny, którego podkatalogi stają się kontami w zapasowym źródle puli.

**katalog roboczy sesji** — miejsce, w którym model zostawia własne pliki; **nie jest dostępem**.

**zasięg wykonania** (`local`, `core`, `remote`) — pole okna; w v2.0 **wyłącznie etykieta prowenancji**.

**rola procesu** (`hub`, `agent`, `all`) — wybór torów pracy rdzenia.

**brama jakości** — automat `budowa/scripts/brama.sh` spinający wszystkie kontrole przed zamknięciem
fali prac.

**Decyzje architektoniczne** — rozstrzygnięcia przyjęte dla budowy platformy, opisane
w dokumentacji i obowiązujące w kodzie.

**Wykaz znanych luk** — zbiór odnotowanych usterek i ich nawrotów; w dokumencie przywołane
dwie: rozjazd portu oraz rozjazd wzorca zmiennych z kodem. Zbiór ten nie ma dziś własnego
nośnika ani w zbiorze dokumentacji, ani wśród źródeł normatywnych — każda luka opisana jest
w miejscu swojego wystąpienia, a zestawienie ograniczeń całego produktu niesie
[Opis produktu](README.md) rozdz. 18. **[DO DECYZJI OPERATORA]** — czy wykaz ma otrzymać
własne opracowanie zbioru dokumentacji, czy pojęcie ma zniknąć ze słownika.

---

*Koniec dokumentu. Instalacja i konfiguracja — Stan wdrożenia, wersja 2.0, 2026-08-20.*

---
*Danaco Console — Platforma AI Workspace OS · v2.0 · status Deweloperski*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz.*
*Warunki korzystania: [Licencja produktu](LICENSE.md). Kontakt: support@danaco-group.pl*
