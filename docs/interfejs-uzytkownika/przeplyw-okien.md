# Danaco Console — Przepływ okien: od uruchomienia aplikacji do okna roboczego

| | |
|---|---|
| **Produkt** | Danaco Console |
| **Rodzaj** | Platforma AI Workspace OS |
| **Opis** | Platforma jest wielośrodowiskowym systemem operacyjnym dla sztucznej inteligencji, integrującym komunikację, zarządzanie wiedzą, tworzenie treści, projektowanie, automatyzacje procesów oraz rozwój oprogramowania. |
| **Producent** | Danaco Holding Group Sp. z o.o. |
| **Twórca** | Dariusz Naharnowicz |
| **Wersja** | v2.0 |
| **Status** | Deweloperski |
| **Data** | 2026-08-20 |

**Informacje szczegółowe dokumentu:**

| | |
|---|---|
| **Tytuł** | Przepływ okien: od uruchomienia aplikacji do okna roboczego |
| **Klasa dokumentu** | Specyfikacja docelowa |
| **Odbiorcy** | projektant · deweloper |
| **Przeznaczenie** | Ustala pełny przepływ nawigacyjny okien platformy od uruchomienia aplikacji klienckiej do okna roboczego, wraz ze stanami przejściowymi, warstwami widoczności i trwałością procesów w tle |
| **Zakres** | okno startowe → okno rejestracji i logowania → przygotowanie środowiska pracy → strona główna → wybór środowiska → przedsionek środowiska → wybór modułu → okno operacyjne, wraz ze ścieżkami szczegółowymi czterech środowisk, wejściami do strefy komponentów własnych i strefy ustawień, kartami sesji i trwałością procesu w tle, magistralą kontekstu, stanami okien i elementów interfejsu |
| **Poza zakresem** | wymiary, tokeny kolorów i typografii, ostateczny układ graficzny — [System wizualny](system-wizualny.md); szczegółowa zawartość poszczególnych okien operacyjnych modułów poza tym, co potrzebne do zilustrowania przejścia — opracowania katalogu `moduly/` oraz [Specyfikacja okien operacyjnych](../specyfikacje/specyfikacja-okien-operacyjnych.md) |
| **Dokument nadrzędny** | [Elementy okien](elementy-okien.md) |
| **Dokumenty powiązane** | [Elementy okien](elementy-okien.md) · [Katalog komponentów](katalog-komponentow.md) · [Okno Konfiguracji](konfiguracja.md) · [Okno Ustawień](ustawienia.md) · [Strona główna i nawigacja](strona-glowna-i-nawigacja.md) · [Koncepcja platformy](../architektura/koncepcja-platformy.md) · [Specyfikacja okien operacyjnych](../specyfikacje/specyfikacja-okien-operacyjnych.md) · [System wizualny](system-wizualny.md) · [Model konfiguracji](../architektura/model-konfiguracji.md) · [Specyfikacja agentów](../specyfikacje/specyfikacja-agentow.md) · [Izolacja i konfigurowalność zależności](../architektura/izolacja-i-zaleznosci.md) · [Architektura techniczna](../architektura/architektura.md) · [Bezpieczeństwo i uwierzytelnianie](../architektura/bezpieczenstwo-i-uwierzytelnianie.md) · [Rozszerzenia](../architektura/rozszerzenia.md) |
| **Prototypy odniesienia** | `design/05-okna/przeplyw/okno-startowe.html` · `design/05-okna/przeplyw/rejestracja-i-logowanie.html` · `design/05-okna/przeplyw/przygotowanie-srodowiska.html` · `design/05-okna/przeplyw/centrum-dowodzenia.html` · `design/05-okna/srodowiska/talkin-przedsionek.html` · `design/05-okna/srodowiska/workspace-przedsionek.html` · `design/05-okna/srodowiska/codestudio-przedsionek.html` · `design/05-okna/srodowiska/multitaskingai-przedsionek.html` |
| **Źródła normatywne** | `budowa/shared/contract.json` (obszary `window`, `session`, `home`, `environment`, `connection`, `auth`) · `design/zasoby/zetony/zetony.css` · `design/zasoby/css/komponenty.css`, `rama.css`, `prototyp.css` |
| **Zasada nadrzędna** | Dokument nie wprowadza nowych funkcji, modułów, okien ani mechanizmów ponad te ustalone w opracowaniach źródłowych; żaden interaktywny element nawigacyjny nie przyjmuje stanu wyłączonego jako blokady wykonania |

Dokument nie wprowadza nowych funkcji, modułów, okien ani mechanizmów ponad te ustalone w opracowaniach źródłowych. Nazwy własne środowisk, modułów, okien i ról przejęto bez zmian ze źródeł. Tam, gdzie dokument nazywa i opisuje etap techniczny wprost nieujęty pod odrębną nazwą własną w źródłach (okno startowe — faza łączenia z serwerem poprzedzająca okno logowania), fakt ten jest wskazany jawnie w miejscu opisu wraz z podstawą źródłową, zgodnie z zasadą przejrzystości wobec czytelnika.

---

## Spis treści

1. [Jak korzystać z dokumentu](#0-jak-korzystać-z-dokumentu)
2. [Słownik ram przepływu](#1-słownik-ram-przepływu)
3. [Mapa całości — pełny przepływ w jednym diagramie](#2-mapa-całości--pełny-przepływ-w-jednym-diagramie)
4. [Osiem etapów przepływu głównego](#3-osiem-etapów-przepływu-głównego)
   - [3.1 Etap 1 — Okno startowe](#31-etap-1--okno-startowe)
   - [3.2 Etap 2 — Okno rejestracji i logowania](#32-etap-2--okno-rejestracji-i-logowania)
   - [3.3 Etap 3 — Przygotowanie środowiska pracy](#33-etap-3--przygotowanie-środowiska-pracy)
   - [3.4 Etap 4 — Strona główna (Centrum dowodzenia)](#34-etap-4--strona-główna-centrum-dowodzenia)
   - [3.5 Etap 5 — Wybór środowiska (Strefa 1)](#35-etap-5--wybór-środowiska-strefa-1)
   - [3.6 Etap 6 — Przedsionek środowiska](#36-etap-6--przedsionek-środowiska)
   - [3.7 Etap 7 — Wybór modułu (lub roli / sekcji orkiestracji)](#37-etap-7--wybór-modułu-lub-roli--sekcji-orkiestracji)
   - [3.8 Etap 8 — Okno operacyjne (okno robocze)](#38-etap-8--okno-operacyjne-okno-robocze)
5. [3a. Warstwy widoczności w przepływie okien](#3a-warstwy-widoczności-w-przepływie-okien)
6. [Tabela przejść — zestawienie zbiorcze](#4-tabela-przejść--zestawienie-zbiorcze)
7. [Ścieżki szczegółowe](#5-ścieżki-szczegółowe)
   - [5.1 Ścieżka TalkIn](#51-ścieżka-talkin)
   - [5.2 Ścieżka WorkSpace](#52-ścieżka-workspace)
   - [5.3 Ścieżka CodeStudio](#53-ścieżka-codestudio)
   - [5.4 Ścieżka MultitaskingAI](#54-ścieżka-multitaskingai)
   - [5.5 Ścieżka Strefy 2 — komponenty własne](#55-ścieżka-strefy-2--komponenty-własne)
   - [5.6 Ścieżka Strefy 3 — konfiguracja, ustawienia, funkcje globalne](#56-ścieżka-strefy-3--konfiguracja-ustawienia-funkcje-globalne)
   - [5.7 Ścieżki strefy pracy szyny nawigacji](#57-ścieżki-strefy-pracy-szyny-nawigacji)
8. [Karty sesji i trwałość procesu w tle](#6-karty-sesji-i-trwałość-procesu-w-tle)
9. [6a. Magistrala kontekstu — przekazywanie artefaktu między modułami](#6a-magistrala-kontekstu--przekazywanie-artefaktu-między-modułami)
10. [Stany okien i elementów interfejsu](#7-stany-okien-i-elementów-interfejsu)
11. [Tabele zbiorcze elementów interfejsu nawigacyjnego](#8-tabele-zbiorcze-elementów-interfejsu-nawigacyjnego)
   - [8.1 Elementy poziomu „okno” (jednostki nawigacji pełnoekranowej)](#81-elementy-poziomu-okno-jednostki-nawigacji-pełnoekranowej)
   - [8.2 Elementy poziomu „nawigacja” (stałe elementy ramy interfejsu)](#82-elementy-poziomu-nawigacja-stałe-elementy-ramy-interfejsu)
   - [8.3 Elementy poziomu „kontrolka” (pojedyncze akcje)](#83-elementy-poziomu-kontrolka-pojedyncze-akcje)
12. [Zgodność z zasadami nadrzędnymi platformy](#9-zgodność-z-zasadami-nadrzędnymi-platformy)
13. [Słowniczek pojęć](#10-słowniczek-pojęć)
14. [Załącznik A. Scenariusz pełny — przejście przez wszystkie warstwy](#załącznik-a-scenariusz-pełny--przejście-przez-wszystkie-warstwy)
15. [Załącznik B. Macierz poprzedzania i następstwa ekranów](#załącznik-b-macierz-poprzedzania-i-następstwa-ekranów)

---

## 0. Jak korzystać z dokumentu

Dokument ma dwóch odbiorców czytających go z różnym celem. Poniższa tabela kieruje uwagę każdego z nich do właściwych fragmentów.

| Odbiorca | Pytanie, na które szuka odpowiedzi | Gdzie szukać |
|---|---|---|
| Projektant | Co użytkownik widzi na każdym etapie, gdzie leży każdy element, w której kolumnie układu, jaką ma formę, wagę wizualną i warstwę widoczności, jak wygląda w każdym stanie | Makiety tekstowe przy każdym etapie (rozdz. 3), rozdział o warstwach widoczności (rozdz. 3a), tabele elementów interfejsu (rozdz. 8), rozdział o stanach (rozdz. 7) |
| Deweloper | Jaki warunek wyzwala przejście, co technicznie się dzieje, jaki jest stan okna przed i po, jak wywoływany jest element danej warstwy widoczności, jak zachowują się procesy sesji w tle | Tabele przejść przy każdym etapie (rozdz. 3) i tabela zbiorcza (rozdz. 4), rozdział o warstwach widoczności (rozdz. 3a), rozdział o trwałości procesów (rozdz. 6) |

Dokument czyta się w dwóch trybach. **Tryb liniowy** — od rozdziału 1 do 10, wraz z rozdziałem 3a, dla pełnego obrazu przepływu od uruchomienia aplikacji do pracy w oknie operacyjnym. **Tryb wyszukiwania** — wejście od razu do rozdziału 5 dla jednej, konkretnej ścieżki (np. „co się dzieje, gdy Operator wybiera CodeStudio i moduł Developer”) bez czytania całości.

Konwencja symboli używana w diagramach ASCII całego dokumentu:

| Symbol | Znaczenie |
|---|---|
| `▼` `▲` `►` `◄` | Kierunek przejścia między stanami lub oknami |
| `├──` `└──` | Rozgałęzienie na warianty tego samego przejścia |
| `┌─┬─┐` `│` `└─┴─┘` | Obrys okna lub panelu w makiecie tekstowej (rozkład przestrzenny) |
| `[?]` | Objaśnienie kontekstowe towarzyszące elementowi konfiguracji |
| `●` | Wskaźnik aktywności lub procesu w tle |
| `✕` | Kontrolka zamknięcia |
| `+` | Kontrolka utworzenia (nowa karta, nowy komponent) |
| `▸` | Stan wybrany / podświetlony pozycji nawigacji |
| `⋮` `☰` | Zwinięty wyzwalacz rozwinięcia kontekstowego (warstwa 3) |
| `Nazwa ▼` | Zwinięte menu progresywne — jeden element zbiorczy zamiast listy pozycji (warstwa 2) |
| `[Znacznik]` | Lekki znacznik kontekstowy w pasku kontekstu; kliknięcie otwiera właściwy selektor |

Makiety tekstowe całego dokumentu rysowane są w stanie spoczynku interfejsu: widoczne są wyłącznie elementy warstwy 1 oraz zwinięte wyzwalacze warstw 2–3. Elementy warstw 2–4 opisane są w tabelach elementów okna, w kolumnach „Warstwa” i „Sposób wywołania” (rozdz. 3a).

---

## 1. Słownik ram przepływu

Poniższe pojęcia organizują cały dokument. Pełny słowniczek pojęć platformy zawiera rozdział 10; tabela poniżej ogranicza się do pojęć niezbędnych do czytania diagramów przepływu.

| Pojęcie | Definicja skrócona | Rozwinięcie |
|---|---|---|
| Okno startowe | Pierwsze okno natywne prezentowane po uruchomieniu aplikacji klienckiej; obejmuje nawiązanie połączenia z serwerem i wstępne uzgodnienie protokołu | rozdz. 3.1 |
| Przygotowanie środowiska pracy | Widok przejściowy między uwierzytelnieniem a stroną główną; odtwarza stan pracy zapamiętany przy ostatnim zamknięciu aplikacji | rozdz. 3.3 |
| Okno rejestracji i logowania | Okno uwierzytelniania Operatora — rejestracja przy pierwszym uruchomieniu, logowanie przy kolejnych; przycisk Pomiń zawsze dostępny i klikalny, wymóg logowania konfigurowalny przez Operatora w oknie konfiguracji | rozdz. 3.2 |
| Strona główna (Centrum dowodzenia) | Przedpokój przed strefą roboczą; trzy strefy: wybór środowiska, komponenty własne, ustawienia | rozdz. 3.4 |
| Środowisko | Najwyższy poziom organizacji pracy: TalkIn, WorkSpace, CodeStudio, MultitaskingAI | rozdz. 3.5 |
| Przedsionek środowiska | Widok wejściowy środowiska przed wyborem modułu: rama okna (szyna nawigacji, belka tytułowa, wstążka), szyna sesji, kafle modułów, listwa działań środowiska | rozdz. 3.6 |
| Przestrzeń robocza środowiska | Widok otwierany po wyborze modułu: rama okna (szyna nawigacji, belka tytułowa, wstążka z pasmem kart sesji, pasek stanu), boczna nawigacja modułów (lub panel orkiestracji), Chat Window, okna operacyjne | rozdz. 3.7–3.8 |
| Moduł | Wyspecjalizowany obszar roboczy wewnątrz środowiska; piętnaście modułów, dostępność wg macierzy | rozdz. 3.7 |
| Okno operacyjne | Pojedynczy element zestawu okien modułu — właściwa przestrzeń wykonania zadania; Chat Window jest oknem wspólnym wszystkich modułów | rozdz. 3.8 |
| Chat Window | Główne okno komunikacji Użytkownik ↔ Wykonawca; okno stałe, obecne w każdym module i każdym środowisku, zajmujące lewą kolumnę obszaru roboczego na pełnej wysokości | rozdz. 3.8, 3a |
| Execution Loop Window | Okno pętli wykonawczej Koordynator ↔ Wykonawca; otwierane w kolumnie sąsiadującej z Chat Window, prezentuje dekompozycję zlecenia, kolejkę zadań i sterowanie przebiegiem pętli | rozdz. 3.8, 3a |
| Warstwa widoczności | Przypisanie elementu interfejsu do jednej z czterech warstw ujawniania funkcjonalności; wyznacza, czy element jest widoczny bez interakcji, po użyciu wyzwalacza, w rozwinięciu kontekstowym, czy wyłącznie poleceniem języka naturalnego, skrótem, wyszukiwarką funkcji lub w trybie administracyjnym | rozdz. 3a |
| Karta sesji | Pozioma karta w pasie kart przestrzeni roboczej, reprezentująca jedną, samodzielną przestrzeń roboczą | rozdz. 6 |
| Komponent własny | Nazwany wytwór Operatora: automatyka, agent, projekt lub profil asystenta, tworzony ze strefy 2 strony głównej | rozdz. 5.5 |

Trzy pojęcia z powyższej tabeli — Środowisko, Moduł, Okno operacyjne — powtarzają trójstopniową hierarchię platformy:

```
ŚRODOWISKO  →  MODUŁ  →  OKNO OPERACYJNE
„w jakim trybie      „jakie zadanie        „jakim narzędziem
 pracuję?”            wykonuję?”            realizuję zadanie?”
```

Niniejszy dokument dokłada do tej hierarchii trzy etapy poprzedzające ją (okno startowe, okno rejestracji i logowania, przygotowanie środowiska pracy) oraz jeden etap pośredni, przez który hierarchia jest zawsze osiągana (strona główna). Cztery etapy dołożone i cztery etapy hierarchii — wybór środowiska, przedsionek środowiska, wybór modułu, okno operacyjne — dają osiem etapów opisanych w rozdziale 3.

---

## 2. Mapa całości — pełny przepływ w jednym diagramie

Poniższy diagram jest mapą odniesienia dla całego dokumentu — pokazuje wszystkie osiem etapów przepływu głównego wraz z trzema rozgałęzieniami strony głównej (strefa 1, strefa 2, strefa 3) oraz miejscem, w którym każda ścieżka szczegółowa z rozdziału 5 dołącza do całości.

```
╔═════════════════════════════════════════════════════════════════════════════╗
║  0 · URUCHOMIENIE APLIKACJI KLIENCKIEJ                                      ║
║      Operator otwiera natywne okno aplikacji na urządzeniu                  ║
╚═════════════════════════════════════════╤═══════════════════════════════════╝
                                          ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│  1 · OKNO STARTOWE                                                          │
│      nawiązanie kanału WebSocket · uzgodnienie protokołu (connection.hello) │
└─────────────┬──────────────────────────────────────┬────────────────────────┘
              │                                      │
   token urządzenia NIEWAŻNY                  token urządzenia WAŻNY
   lub brak tokenu                            (urządzenie już uwierzytelnione)
              │                                      │
              ▼                                      ├── zapamiętana aktywna karta
┌──────────────────────────────────────┐             │      ──► [ 8 ] OKNO OPERACYJNE
│  2 · OKNO REJESTRACJI I LOGOWANIA    │             │      (pomija stronę główną — rozdz. 3.1, 6)
│                                      │             │
│  FAZA BUDOWY         FAZA PRODUKCYJNA│             └── brak zapamiętanej karty
│  [Pomiń] klikalny    [Pomiń] klikalny│                    ──► [ 4 ] STRONA GŁÓWNA
│                                      │
│  REJESTRACJA         LOGOWANIE       │
│  (konto nie          (konto już      │
│   istnieje)           istnieje)      │
│                                      │
│  token dostępu wydany urządzeniu     │
└─────────────────────┬────────────────┘
                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│  3 · PRZYGOTOWANIE ŚRODOWISKA PRACY                                         │
│      widok przejściowy — profil, uprawnienia, przywracanie sesji,           │
│      kanały modeli, magistrala kontekstu                                    │
│      bryła warstw platformy w obrocie = jedyny ruch ciągły widoku           │
└─────────────────────┬───────────────────────────────────────────────────────┘
                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│  4 · STRONA GŁÓWNA — CENTRUM DOWODZENIA                                     │
│      przedpokój przed strefą roboczą — trzy niezależne strefy               │
│                                                                             │
│   STREFA 1 (waga główna)     STREFA 2 (waga pośrednia)  STREFA 3 (najniższa)│
│   karty środowisk             kafle komponentów własnych   listwa ustawień  │
│   TalkIn·WorkSpace·           Automations·Agents·          Okno konfiguracji│
│   CodeStudio·MultitaskingAI   Workspace·Assistant           Mobile · AOD    │
│         │                            │                            │         │
└─────────┼────────────────────────────┼────────────────────────────┼─────────┘
          ▼                            ▼                            ▼
┌────────────────────────┐   ┌──────────────────────────────┐   ┌──────────────────────┐
│ 5 · WYBÓR ŚRODOWISKA   │   │ 5.5 STREFA 2 — SZCZEGÓŁY     │   │ 5.6 STREFA 3         │
│ kliknięcie karty       │   │ okno budowy komponentu       │   │ okno konfiguracji    │
│ środowiska             │   │ własnego (np. Agent          │   │ (nakładka modalna)   │
│                        │   │ Builder) → zapis jako        │   │ albo aktywacja       │
│                        │   │ komponent własny             │   │ Mobile / AOD         │
└───────────┬────────────┘   └────────────────┬─────────────┘   └────────────┬─────────┘
            │                                 ▼                              ▼
            │                 wytwór (agent, automatyka,        aktywacja funkcji
            │                 projekt, profil asystenta)         globalnej — BEZ nowej
            │                 zapisany jako zasób Operatora,     przestrzeni roboczej;
            │                 wybieralny dalej w krokach 7–8 ─┐  bieżący widok trwa
            │                 (wpięty w sesji modułu)         │  niezmieniony pod spodem
            ▼                                                 │
┌──────────────────────────────────────────────────────────┐  │
│  6 · PRZEDSIONEK ŚRODOWISKA                              │◄─┘
│      rama okna · szyna sesji („Twoje sesje”)             │
│      · kafle modułów · listwa działań środowiska         │
│      ── żaden moduł nie jest jeszcze otwarty ──          │
└──────────────────────────────┬───────────────────────────┘
                               ▼
┌──────────────────────────────────────────────────────────┐
│  7 · WYBÓR MODUŁU (lub sekcji orkiestracji — rozdz. 5.4) │
│      kliknięcie kafla w przedsionku — otwiera przestrzeń │
│      roboczą; później: pozycja bocznej nawigacji         │
└──────────────────────────────┬───────────────────────────┘
                               ▼
┌──────────────────────────────────────────────────────────────────────────────┐
│  8 · OKNO OPERACYJNE — OKNO ROBOCZE                                          │
│      układ pionowy (podział lewa–prawa):                                     │
│      Chat Window (lewa kolumna, stała) │ Execution Loop Window (kolumna      │
│      sąsiadująca) │ obszar roboczy modułu (kolumna dominująca) │ okna        │
│      pomocnicze jako rozszerzenia boczne po prawej                           │
│      ── tu zaczyna się praca właściwa ──                                     │
│                                                                              │
│      ścieżki szczegółowe:  5.1 TalkIn · 5.2 WorkSpace ·                      │
│                             5.3 CodeStudio · 5.4 MultitaskingAI              │
└────────────────────┬─────────────────────────────────────────────────────────┘
                     │  „+” nowa karta sesji (powrót do etapu 7 w nowej karcie)
                     │  powrót do strony głównej (powrót do etapu 4;
                     │     karta bieżąca trwa w tle — rozdz. 6)
                     ▼
              (pętla pracy równoległej — wiele kart, wiele środowisk jednocześnie)
```

Trzy uwagi porządkujące mapę całości, rozwinięte w dalszych rozdziałach:

| Uwaga | Rozwinięcie |
|---|---|
| Etapy 1–3 poprzedzają hierarchię platformy i występują dokładnie raz na sesję połączenia urządzenia (nie przy każdym powrocie do pracy) | rozdz. 3.1–3.3, rozdz. 6 |
| Etap 4 (strona główna) jest punktem, do którego przepływ zawsze może wrócić — nie jest odwiedzany tylko raz | rozdz. 3.4, rozdz. 8.2 dokumentu Strona główna i nawigacja |
| Etapy 5–8 tworzą pętlę: z etapu 8 Operator może otworzyć nową kartę (powrót do etapu 7 w nowej karcie) albo wrócić do strony głównej (powrót do etapu 4), bez zamykania pracy pozostawionej w tle | rozdz. 6 niniejszego dokumentu |

---

## 3. Osiem etapów przepływu głównego

Konwencja opisu każdego etapu jest jednolita: **warunek** (co musi zajść, aby etap się rozpoczął), **co się dzieje** (mechanizm), **stan okna** (jak wygląda i zachowuje się okno w tym etapie), **makieta tekstowa** (rozkład przestrzenny), **tabela elementów interfejsu** (co to jest · do czego służy · forma i waga · stany · zachowanie po interakcji · gdzie występuje).

### 3.1 Etap 1 — Okno startowe

| Aspekt | Treść |
|---|---|
| Warunek wejścia | Operator uruchamia pakiet kliencki na urządzeniu (komputer, telefon, tablet); zdarzenie systemowe, nie decyzja nawigacyjna |
| Co się dzieje | Natywne okno aplikacji (powłoka Tauri) otwiera się i osadza silnik prezentacji; klient otwiera kanał WebSocket pod adresem serwera skonfigurowanym w pakiecie klienckim; klient wysyła komendę `connection.hello` (wersja protokołu, identyfikator klienta, wersja klienta, token dostępu — o ile posiadany); serwer odpowiada synchronicznie polami wyniku tej samej komendy: wersja serwera, wersja protokołu rdzenia, `authenticated` (czy połączenie jest związane z ważną sesją bramki), `gatewayConfigured` (czy sekret bramki jest już ustawiony) i `loginRequired` (czy na tym nasłuchu obowiązuje wymóg logowania) |
| Stan okna | Przejściowy — trwa od ułamka sekundy do kilku sekund, zależnie od jakości połączenia. Nie jest miejscem żadnej decyzji Operatora; jedyna możliwa interakcja to przerwanie (zamknięcie okna) |
| Wynik | Rozgałęzienie do etapu 2 (okno rejestracji i logowania) albo, gdy urządzenie posiada ważny token, wprost do etapu 3 (przygotowanie środowiska pracy) |

**Uzasadnienie nazwy.** Źródła nie nadają odrębnej nazwy własnej fazie łączenia poprzedzającej okno logowania, lecz opisują ją jako etapy 1–2 cyklu życia połączenia (nawiązanie kanału WebSocket, uzgodnienie i uwierzytelnienie). Niniejszy dokument nazywa tę fazę „oknem startowym”, ponieważ w interfejsie musi mieć widoczną, choćby krótkotrwałą reprezentację — Operator nie może patrzeć na puste okno w czasie trwania uzgodnienia. Nazwa porządkuje istniejący mechanizm; nie wprowadza nowego.

```
 Makieta — Okno startowe
 ┌─────────────────────────────────────────────┐
 │                                               │
 │                                               │
 │                  ◈  Danaco Console              │  ← godło marki, wyśrodkowane
 │                                               │
 │              ○ ○ ●  łączenie z serwerem…      │  ← wskaźnik stanu (spinner)
 │                                               │
 │                                               │
 └─────────────────────────────────────────────┘
   Bez paska nawigacji, bez bocznych paneli — okno jednolite, jeden komunikat stanu
```

**Warianty stanu okna startowego:**

| Wariant | Wyzwalacz | Co widzi Operator | Dalsze przejście |
|---|---|---|---|
| Łączenie (domyślny) | Otwarcie aplikacji | Godło marki i wskaźnik ładowania z komunikatem „łączenie z serwerem” | Automatyczne, po odpowiedzi serwera |
| Powrót z ważnym tokenem | Urządzenie było już uwierzytelnione | Jak wyżej, czas trwania zwykle krótszy — bez kroku rejestracji/logowania | Etap 4 albo etap 8 (rozdz. 6) |
| Błąd połączenia | Serwer nieosiągalny, adres nieprawidłowy, przerwanie sieci | Komunikat stanu zastąpiony informacją o braku połączenia oraz kontrolką ponowienia | Ponowienie próby (ten sam etap) — nie jest to blokada dostępu, lecz odzwierciedlenie faktycznego stanu sieci |

**Tabela elementów interfejsu — Okno startowe:**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Gdzie występuje |
|---|---|---|---|---|---|---|---|---|
| Godło marki | Sygnet „Delegacja” (`design/zasoby/marka/logo/sygnet.svg`), wklejony inline w dużej skali, z animowaną kropką sygnału, wyśrodkowany na powierzchni okna | Identyfikacja marki w czasie oczekiwania | Duży, pojedynczy element graficzny — jedyny akcent wizualny okna | 1 | Widoczne bez interakcji | Statyczny (nie zmienia stanu) | Brak interakcji | Wyłącznie okno startowe |
| Wskaźnik ładowania | Animowany spinner (`.dn-spinner`) z komunikatem tekstowym | Sygnalizacja trwającego połączenia z serwerem | Mały element graficzny + krótki tekst pod godłem | 1 | Widoczny bez interakcji w czasie łączenia | Ładowanie (obrót ciągły, respektuje ograniczenie ruchu) → zastępowany komunikatem błędu przy niepowodzeniu | Brak interakcji w stanie ładowania | Wyłącznie okno startowe |
| Komunikat błędu połączenia | Tekst i ikona ostrzeżenia zastępujące wskaźnik ładowania | Poinformowanie o braku połączenia z serwerem | Mały blok tekstowy z ikoną `ostrzezenie` | 1 | Wyświetlany samoczynnie w stanie błędu połączenia | Błąd (widoczny wyłącznie po niepowodzeniu połączenia) | Wyświetlany do chwili automatycznego ponowienia lub ręcznej kontrolki „Spróbuj ponownie” | Wyłącznie okno startowe, stan błędu |
| Kontrolka „Spróbuj ponownie” | Przycisk drugorzędny (`.dn-btn--zarys`) | Ręczne ponowienie próby połączenia | Mały przycisk pod komunikatem błędu | 2 | Widoczna po wystąpieniu błędu połączenia; kliknięcie ponawia etap 1 | Domyślny · wskazanie kursorem · aktywny | Ponawia etap 1 od nawiązania kanału WebSocket | Wyłącznie okno startowe, stan błędu |

---

### 3.2 Etap 2 — Okno rejestracji i logowania

| Aspekt | Treść |
|---|---|
| Warunek wejścia | Krok „Uwierzytelnienie” cyklu życia połączenia wskazuje brak ważnego tokenu na urządzeniu (pierwsze połączenie tego urządzenia albo token unieważniony) — **albo** faza budowy, niezależnie od stanu tokenu, gdy okno pozostaje widoczne jako podgląd wyglądu docelowego. Czy okno wymusza skuteczne wypełnienie formularza przed przejściem dalej, ustala pozycja „Wymóg logowania” w oknie konfiguracji (rozdz. 5.6) — ustawienie dostępne Operatorowi w obu fazach wdrożenia, z domyślną wartością nieaktywną w fazie budowy i aktywną w fazie produkcyjnej |
| Co się dzieje | Klient prezentuje okno rejestracji i logowania w pełnym, docelowym wyglądzie niezależnie od fazy wdrożenia; przycisk Pomiń jest obecny i klikalny w obu fazach. Gdy ustawienie „Wymóg logowania” jest nieaktywne, kliknięcie przycisku Pomiń przenosi natychmiast do etapu 4. Gdy ustawienie jest aktywne, kliknięcie Pomiń wyświetla komunikat o obowiązującym wymogu uwierzytelnienia zamiast pominięcia kroku — dalsze przejście wymaga skutecznej rejestracji albo logowania |
| Stan okna | Interaktywny formularz — jedyny etap całego przepływu, w którym Operator wprowadza dane uwierzytelniające. Trzy pod-stany: rejestracja (konto nie istnieje), logowanie (konto istnieje), odzyskiwanie konta (utrata dostępu) |
| Limit prób | Logowanie przyjmuje **pięć nieudanych prób**. Każdy komunikat błędu nazywa liczbę prób pozostałych, a po wykorzystaniu wszystkich pięciu logowanie zostaje **zablokowane na godzinę**; przez ten czas dostęp przywraca wyłącznie odzyskiwanie konta adresem e-mail |
| Wynik | Wydanie urządzeniu tokenu dostępu i przejście do etapu 3 (przygotowanie środowiska pracy) — albo, gdy ustawienie „Wymóg logowania” jest nieaktywne, przejście do etapu 3 przez przycisk Pomiń, bez wydania tokenu uwierzytelnionego konta |

Rozstrzygające dla wyglądu okna jest fazowanie wdrożenia; rozstrzygające dla wymogu logowania jest niezależne ustawienie Operatora w oknie konfiguracji — wspólny kod obsługi w obu fazach:

```
 FAZA BUDOWY                              FAZA PRODUKCYJNA
 rdzeń lokalnie na komputerze    ────►     rdzeń zmigrowany na serwer docelowy
 Operatora                                  + otwarty port kanału WebSocket

  okno obecne, wygląd docelowy               okno obecne, wygląd docelowy
  [ Pomiń ]  zawsze klikalny                 [ Pomiń ]  zawsze klikalny
  „Wymóg logowania”: domyślnie NIEAKTYWNY    „Wymóg logowania”: domyślnie AKTYWNY
  (Operator zmienia w oknie konfiguracji —    (Operator zmienia w oknie konfiguracji —
   zakres „Aplikacja”, w każdej z faz)         zakres „Aplikacja”, w każdej z faz)

        ── ten sam, wspólny kod obsługi rejestracji, logowania,
           metod dodatkowych i odzyskiwania konta w obu fazach ──
           przejście między fazami nie wymaga nowego kodu —
           „Wymóg logowania” jest ustawieniem okna konfiguracji,
           dostępnym Operatorowi w obu fazach, z domyślną wartością
           zależną od fazy i zmienialną w każdej chwili
```

**Trzy pod-stany okna, w zależności od sytuacji Operatora:**

```
Okno rejestracji i logowania
        │
        ├── Konto właściciela NIE ISTNIEJE (pierwsze uruchomienie platformy)
        │       │
        │       ▼
        │   REJESTRACJA
        │   pola: login · adres e-mail uwierzytelniający · hasło
        │       │
        │       ▼
        │   weryfikacja adresu e-mail (potwierdzenie)
        │       │
        │       ▼
        │   konto Operatora utworzone — metody aktywne: hasło + e-mail
        │   uwierzytelniający; metody nieaktywne: PIN, Windows Hello
        │
        ├── Konto właściciela ISTNIEJE (kolejne uruchomienie / kolejne urządzenie)
        │       │
        │       ▼
        │   LOGOWANIE
        │   login + hasło  albo  aktywna metoda dodatkowa:
        │   PIN · Windows Hello · e-mail uwierzytelniający
        │       │
        │       ├── weryfikacja pozytywna ──► token dostępu wydany urządzeniu
        │       └── weryfikacja negatywna ──► ponowna próba albo:
        │
        └── ODZYSKIWANIE KONTA (utrata dostępu do dotychczasowego materiału)
                │
                ▼
            podanie adresu e-mail uwierzytelniającego
                │
                ▼
            potwierdzenie tożsamości drogą e-mail
                │
                ▼
            ustawienie nowego hasła — zastępuje dotychczasowy skrót
                │
                ▼
            unieważnienie tokenów wydanych wcześniej
            (ponowne logowanie wymagane na wszystkich powiązanych urządzeniach)
                │
                ▼
            powrót do LOGOWANIA z nowym hasłem
```

**Makieta tekstowa — Okno rejestracji i logowania (widok logowania):**

```
 ┌───────────────────────────────────────────────────────┐
 │                     ◈  Danaco Console                     │
 │                                                          │
 │   ┌───────────────┬───────────────┐                     │
 │   │ ▸ Zaloguj się │  Zarejestruj  │  ← zakładki (.dn-zakladki)
 │   └───────────────┴───────────────┘                     │
 │                                                          │
 │   Login             [_________________________]         │
 │   Hasło             [_________________________]  [👁]   │
 │                                                          │
 │   Metoda logowania:  ( ) hasło  ( ) PIN  ( ) Windows     │
 │                       ( ) e-mail uwierzytelniający       │
 │                       — widoczne tylko metody aktywne —  │
 │                                                          │
 │              [       Zaloguj się       ]  ← .dn-btn--sygnal (CTA)
 │                                                          │
 │   Nie pamiętam hasła →                                   │
 │   ─────────────────────────────────────────────         │
 │              [          Pomiń          ]  ← .dn-btn--duch, zawsze klikalny (obie fazy)
 │              wejście bez zakładania konta                │
 └───────────────────────────────────────────────────────┘
```

**Tabela elementów interfejsu — Okno rejestracji i logowania:**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Gdzie występuje |
|---|---|---|---|---|---|---|---|---|
| Zakładki „Zaloguj się / Zarejestruj się” | Przełącznik segmentowy (`.dn-zakladki`) | Wybór między logowaniem a rejestracją | Element średniej wagi, u góry formularza | 1 | Widoczne bez interakcji | Domyślny (Zaloguj się aktywna, gdy konto już istnieje) · wybrana · wskazanie kursorem | Przełącza zestaw pól poniżej bez opuszczania okna | Okno rejestracji i logowania |
| Pole „Login” | Pole tekstowe (`.dn-pole-kontrolka`) | Wprowadzenie nazwy właściciela konta | Pole średniej szerokości z etykietą | 1 | Widoczne bez interakcji | Domyślny · fokus · wypełnione · błąd (format) | Wartość zapamiętywana do chwili wysłania formularza | Widok logowania i rejestracji |
| Pole „Hasło” | Pole tekstowe maskowane, z przełącznikiem widoczności (ikona `oko`) | Wprowadzenie hasła dostępowego | Pole średniej szerokości z etykietą i ikoną odsłonięcia | 1 | Widoczne bez interakcji | Domyślny · fokus · wypełnione · błąd (niezgodność) | Odsłonięcie/zamaskowanie treści po kliknięciu ikony; wartość przesyłana wyłącznie po zatwierdzeniu formularza | Widok logowania i rejestracji |
| Pole „Adres e-mail uwierzytelniający” | Pole tekstowe (`.dn-pole-kontrolka`, walidacja formatu) | Wprowadzenie adresu pełniącego funkcję weryfikacji, logowania i jedynej drogi odzyskania konta | Pole średniej szerokości z etykietą | 1 | Widoczne bez interakcji w widoku rejestracji | Domyślny · fokus · wypełnione · błąd (format) | Po rejestracji inicjuje wysyłkę wiadomości potwierdzającej | Widok rejestracji; widok logowania, gdy metoda e-mail aktywna |
| Wybór metody logowania | Grupa przełączników radiowych (`.dn-check`) | Wskazanie, którą z aktywnych metod Operator się loguje | Mały zestaw pozycji pod polami głównymi | 2 | Zestaw metod aktywnych; wybór zwija się po użyciu | Widoczne wyłącznie metody aktywne (hasło zawsze; PIN i Windows Hello tylko po włączeniu w oknie konfiguracji) | Zmiana metody przełącza układ pól poniżej (np. PIN zamiast hasła) | Widok logowania |
| Przycisk „Zaloguj się” / „Zarejestruj się” | Przycisk główny CTA (`.dn-btn--sygnal`) | Zatwierdzenie formularza i wysłanie danych do serwera | Duży, pojedynczy przycisk akcentowany sygnałem — jedyny CTA widoku | 1 | Widoczny bez interakcji | Domyślny · wskazanie kursorem · aktywny (wciśnięcie) · ładowanie (podczas weryfikacji) | Zawsze klikalny; przy niewypełnionych polach wymaganych kliknięcie wyświetla ostrzeżenie przy brakującym polu, bez wysyłki żądania; przy polach wypełnionych wysyła polecenie `auth.login` albo `auth.register` — przy powodzeniu przejście do etapu 4, przy niepowodzeniu komunikat błędu przy właściwym polu, formularz pozostaje wypełniony | Widok logowania i rejestracji |
| Link „Nie pamiętam hasła” | Odnośnik tekstowy | Wejście w funkcję odzyskiwania konta | Mały element tekstowy, waga najniższa formularza | 2 | Kliknięcie odnośnika otwiera widok odzyskiwania konta | Domyślny · wskazanie kursorem · fokus | Zamienia formularz na widok odzyskiwania konta (pole adresu e-mail) | Widok logowania |
| Przycisk „Pomiń” | Przycisk drugorzędny bez obrysu (`.dn-btn--duch`) | Wejście do platformy bez wypełniania formularza, gdy wymóg logowania nie jest aktywny | Mały przycisk, wyraźnie oddzielony kreską rozdzielającą (bez odrębnej klasy w arkuszu) od formularza głównego, najniższa waga wizualna okna | 1 | Widoczny bez interakcji | Obecny i klikalny w obu fazach wdrożenia | Gdy ustawienie „Wymóg logowania” (oknie konfiguracji, konfigurowalne przez Operatora niezależnie od fazy) jest nieaktywne — natychmiastowe przejście do etapu 4 bez wydania tokenu, krok „Uwierzytelnienie” pozostaje pominięty; gdy ustawienie jest aktywne — komunikat o obowiązującym wymogu uwierzytelnienia, formularz pozostaje otwarty | Okno logowania, obie fazy |
| Komunikat błędu formularza | Alert liniowy (`.dn-toast--blad`) przy polu, którego weryfikacja się nie powiodła | Wskazanie przyczyny niepowodzenia (dane niepoprawne, adres niepotwierdzony) | Mały blok tekstowy z kreską lewą czerwoną, przy właściwym polu | 1 | Wyświetlany samoczynnie po niepowodzeniu weryfikacji | Widoczny wyłącznie po niepowodzeniu weryfikacji | Znika po poprawnym wypełnieniu pola i ponownym zatwierdzeniu | Widok logowania, rejestracji i odzyskiwania konta |

---

### 3.3 Etap 3 — Przygotowanie środowiska pracy

| Aspekt | Treść |
|---|---|
| Warunek wejścia | Uwierzytelnienie zakończone powodzeniem — token wydany urządzeniu — albo wejście przyciskiem „Pomiń” przy nieaktywnym wymogu logowania |
| Co się dzieje | Klient odtwarza stan pracy zapamiętany przy ostatnim zamknięciu aplikacji: profil Operatora i uprawnienia, karty sesji trwające po stronie serwera, kanały modeli i konektory, magistralę kontekstu i pamięć projektów. Kroki wykonują się kolejno, a widok pokazuje, który z nich trwa |
| Stan okna | Widok przejściowy, jednolita powierzchnia. Po lewej bryła przestrzenna czterech warstw platformy — rdzeń, środowiska, moduły, stanowisko — w powolnym obrocie; po prawej wykaz kroków przygotowania ze stanem każdego z nich oraz pasek postępu. Okno nie jest miejscem żadnej decyzji Operatora |
| Wynik | Otwarcie strony głównej (etap 4). Gdy przywracanie sesji zostanie pominięte, strona główna otwiera się bez odtworzonych kart — sesje pozostają dostępne w szynie przedsionka każdego środowiska |

**Dlaczego etap jest osobny.** Uwierzytelnienie kończy się w chwili wydania
tokenu, a strona główna wymaga stanu, którego w tej chwili jeszcze nie ma:
listy sesji trwających w tle, dostępnych kanałów modeli i pamięci projektów.
Wyświetlenie strony głównej przed skompletowaniem tego stanu pokazałoby
Operatorowi obraz nieprawdziwy — środowiska bez sesji i karty bez treści.

**Ruch w widoku.** Bryła warstw jest jedynym ruchem ciągłym tego widoku
i zastępuje wskaźnik nieoznaczonego postępu. Przy `prefers-reduced-motion`
obrót ustaje, a bryła pozostaje w położeniu nieruchomym — informację o trwaniu
pracy niesie wówczas wykaz kroków i pasek postępu.

**Makieta tekstowa — Przygotowanie środowiska pracy:**

```
 ┌───────────────────────────────────────────────────────────────────────────┐
 │ ◈ Danaco Console        Przygotowanie środowiska pracy      ─  □  ✕       │
 ├───────────────────────────────────────────────────────────────────────────┤
 │ ☰ ⇤ │ 🔍 ⛶ 📋 │ ← →     Przepływ główny › Etap 3     [ 🔍 ]  ⌂ ⟳ ⚙ ◐ ◉  │
 ├───────────────────────────────────────────────────────────────────────────┤
 │                                                                            │
 │        ╱▔▔▔▔▔▔▔╲              DANACO CONSOLE                              │
 │       ╱ stanowisko╲            Przygotowuję środowisko pracy                │
 │      ╱▔▔▔▔▔▔▔╲                                                            │
 │     ╱  moduły   ╲             ✓ Uwierzytelnienie      token wydany         │
 │    ╱▔▔▔▔▔▔▔╲                  ✓ Profil i uprawnienia  18 pozycji           │
 │   ╱ środowiska  ╲             ● Przywracanie sesji    3 z 7 kart           │
 │  ╱▔▔▔▔▔▔▔╲                    4 Kanały modeli         oczekuje            │
 │ ╱   rdzeń       ╲             5 Magistrala kontekstu  oczekuje            │
 │ ╲▁▁▁▁▁▁▁╱                                                                 │
 │    bryła w obrocie             Przywracanie sesji              62%         │
 │    (jedyny ruch ciągły)        ▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░░░░░                      │
 │                                                                            │
 │                                [Pomiń przywracanie] [Przerwij i wyloguj]   │
 └───────────────────────────────────────────────────────────────────────────┘
```

**Tabela elementów interfejsu — Przygotowanie środowiska pracy:**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Gdzie występuje |
|---|---|---|---|---|---|---|---|---|
| Bryła warstw platformy | Cztery płaszczyzny w rzucie przestrzennym, powolny obrót wokół osi pionowej | Wskazanie, że praca trwa, oraz przypomnienie budowy platformy: rdzeń → środowiska → moduły → stanowisko | Element graficzny wagi głównej, lewa połowa widoku | 1 | Widoczna bez interakcji | Obrót ciągły · nieruchoma przy `prefers-reduced-motion` | Brak interakcji własnej | Wyłącznie ten widok |
| Wykaz kroków przygotowania | Lista pięciu kroków ze znacznikiem stanu i miarą postępu każdego | Pokazanie, co dokładnie się dzieje i na czym praca stoi | Lista wagi średniej, prawa połowa widoku | 1 | Widoczny bez interakcji | Krok: oczekuje · pracuje (kropka tętna) · gotowy (znak potwierdzenia) | Brak interakcji własnej; stany zmieniają się samoczynnie | Wyłącznie ten widok |
| Pasek postępu | Tor z wartością (`.dn-postep`) wraz z podpisem procentowym | Miara postępu kroku bieżącego | Tor 4 px, waga niska | 1 | Widoczny bez interakcji | Wartość od 0 do 100 | Brak interakcji własnej | Wyłącznie ten widok |
| Kontrolka „Pomiń przywracanie sesji” | Przycisk `--duch`, mały | Przejście do strony głównej bez odtwarzania kart sesji | Przycisk tekstowy, waga najniższa | 1 | Widoczna bez interakcji | Domyślny · wskazanie kursorem · fokus | Przerywa przywracanie i otwiera stronę główną; sesje pozostają dostępne w przedsionkach środowisk | Wyłącznie ten widok |
| Kontrolka „Przerwij i wyloguj” | Przycisk `--duch`, mały | Wyjście z przygotowania i powrót do okna rejestracji i logowania | Przycisk tekstowy, waga najniższa | 1 | Widoczna bez interakcji | Domyślny · wskazanie kursorem · fokus | Unieważnia token wydany urządzeniu i wraca do etapu 2 | Wyłącznie ten widok |

---

### 3.4 Etap 4 — Strona główna (Centrum dowodzenia)

| Aspekt | Treść |
|---|---|
| Warunek wejścia | Token dostępu wydany (etap 2) i brak zapamiętanej aktywnej karty sesji do przywrócenia — **albo** świadomy powrót Operatora ze środowiska (rozdz. 6) — **albo** zamknięcie okna komponentu własnego strefy 2 lub okna konfiguracji strefy 3 |
| Co się dzieje | Klient wysyła `home.enter` (etap 4 cyklu życia połączenia — „Powiązanie”); serwer wiąże połączenie ze stanem strony głównej. Żadna przestrzeń robocza środowiska nie jest jeszcze otwarta — strona główna nie otwiera jej bezpośrednio, pełni funkcję przedpokoju |
| Stan okna | Stabilny, statyczny układ trzech stref o malejącej wadze wizualnej; brak treści ładowanej na żywo (poza ewentualnym wskaźnikiem aktywności kart sesji trwających w tle) |
| Wynik | Trzy niezależne dalsze ścieżki, wg strefy klikniętej przez Operatora — rozdz. 3.5 (strefa 1), rozdz. 5.5 (strefa 2), rozdz. 5.6 (strefa 3) |

**Makieta tekstowa — Strona główna, trzy strefy:**

```
 ┌───────────────────────────────────────────────────────────────────────┐
 │  ◈ Danaco Console                                    [ 🔍 ]   ● Operator │  ← belka + wstążka
 ├───────────────────────────────────────────────────────────────────────┤
 │                                                                         │
 │  STREFA 1 · WYBÓR ŚRODOWISKA                          (waga główna)     │
 │  ┌───────────┐  ┌───────────┐  ┌────────────┐  ┌─────────────────┐    │
 │  │  ◈ godło   │  │  ◈ godło  │  │  ◈ godło    │  │  ◈ godło          │    │
 │  │  TalkIn    │  │ WorkSpace │  │ CodeStudio  │  │ MultitaskingAI    │    │
 │  │  Wiedza,   │  │ Produkty- │  │ Programo-   │  │ Orkiestracja       │    │
 │  │  komunika- │  │ wność,    │  │ wanie       │  │ autonomicznej      │    │
 │  │  cja i     │  │ organiza- │  │             │  │ pracy ciągłej      │    │
 │  │  treść     │  │ cja       │  │             │  │                    │    │
 │  └───────────┘  └───────────┘  └────────────┘  └─────────────────┘    │
 │                                                                         │
 │  STREFA 2 · KOMPONENTY WŁASNE                       (waga pośrednia)   │
 │  ┌────────────┐  ┌─────────┐  ┌───────────┐  ┌───────────┐            │
 │  │ ▧ Automa-  │  │ ▧ Agents│  │ ▧ Work-   │  │ ▧ Assistant│            │
 │  │  tions     │  │         │  │  space    │  │            │            │
 │  │ „Zbuduj    │  │ „Skonfi-│  │ „Załóż    │  │ „Ustaw     │            │
 │  │  automaty- │  │  guruj  │  │  projekt” │  │  profil    │            │
 │  │  kę”       │  │  agenta”│  │           │  │  asystenta”│            │
 │  └────────────┘  └─────────┘  └───────────┘  └───────────┘            │
 │                                                                         │
 │  STREFA 3 · USTAWIENIA                              (waga najniższa)   │
 │  [ ⚙ Okno konfiguracji ]   [ 📱 Mobile ]   [ ● Always On Display ]     │
 │                                                                         │
 └───────────────────────────────────────────────────────────────────────┘
```

**Tabela elementów interfejsu — Strona główna (poziom strefy, bez rozbicia na pojedyncze karty — te w rozdz. 3.5 i 5.5–5.6):**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Gdzie występuje |
|---|---|---|---|---|---|---|---|---|
| Belka tytułowa | Pas atramentowy na prawo od szyny nawigacji, ze znakiem i nazwą Danaco Console oraz tytułem bieżącego widoku | Identyfikacja marki i rozpoznanie okna; sterowanie oknem systemowym | Pas ramy, waga niska (element ramowy, nie treściowy) | 1 | Widoczny bez interakcji | Statyczny; tytuł zmienny wraz z widokiem | Znak wraca na stronę główną; wyszukiwanie i konto niesie wstążka oraz pasek stanu | Strona główna, przedsionek i przestrzeń robocza (element wspólny) |
| Siatka strefy 1 (cztery karty) | Kontener siatki czterech elementów równej wagi | Grupowanie punktów wejścia do środowisk | Duży blok, środek ciężkości strony | 1 | Widoczna bez interakcji | — (kontener) | — | Wyłącznie strona główna |
| Siatka strefy 2 (cztery kafle) | Kontener siatki czterech elementów lżejszych niż strefa 1 | Grupowanie punktów wejścia do budowy komponentów własnych | Średni blok, pod siatką strefy 1 | 1 | Widoczna bez interakcji | — (kontener) | — | Wyłącznie strona główna |
| Listwa strefy 3 | Zwarty pasek poziomy, trzy pozycje | Grupowanie dostępu do ustawień i funkcji globalnych | Wąski pasek, najniższa waga wizualna całej strony | 1 | Widoczna bez interakcji | — (kontener) | — | Strona główna (dodatkowo Mobile i Always On Display dostępne z każdego środowiska — rozdz. 5.6) |

Szczegółowa anatomia pojedynczej karty środowiska, kafla komponentu i pozycji listwy — wraz z pełną tabelą stanów każdego z nich — znajduje się w rozdziałach 3.5 (strefa 1), 5.5 (strefa 2) i 5.6 (strefa 3), aby uniknąć powtórzenia tej samej treści w dwóch miejscach dokumentu.

---

### 3.5 Etap 5 — Wybór środowiska (Strefa 1)

| Aspekt | Treść |
|---|---|
| Warunek wejścia | Operator znajduje się na stronie głównej (etap 4) |
| Co się dzieje | Kliknięcie jednej z czterech kart środowisk w strefie 1 wysyła polecenie nawigacyjne serwerowi; serwer otwiera lub przywraca kontekst wybranego środowiska dla nowej karty sesji. Jest to **jedyna droga wejścia** do przestrzeni roboczej środowiska — strefy 2 i 3 nie prowadzą do środowisk, lecz do okien konfiguracyjnych i ustawień |
| Stan okna | Strona główna pozostaje widoczna do chwili potwierdzenia przejścia; karta klikana przechodzi przez stan aktywny (wciśnięcie) i stan wskazania kursorem, zanim widok zostanie zastąpiony oknem środowiska |
| Wynik | Otwarcie **przedsionka środowiska** (etap 6): szyna sesji, kafle modułów i listwa działań środowiska. Przestrzeń robocza — boczna nawigacja, pas kart sesji, Chat Window — nie otwiera się przed wyborem modułu |

**Cztery karty i ich przeznaczenie:**

| Kolejność | Środowisko | Opis trybu pracy na karcie | Otwiera |
|---|---|---|---|
| 1 | TalkIn | Wiedza, komunikacja i praca z treścią | Przedsionek z dziewięcioma kaflami modułów |
| 2 | WorkSpace | Produktywność, organizacja i realizacja projektów | Przedsionek z dziewięcioma kaflami modułów |
| 3 | CodeStudio | Programowanie | Przedsionek z ośmioma kaflami modułów |
| 4 | MultitaskingAI | Orkiestracja autonomicznej pracy ciągłej | Przedsionek z sześcioma kaflami sekcji orkiestracji (rozdz. 3.6, 5.4) zamiast kafli modułów |

**Tabela elementów interfejsu — Karta środowiska (strefa 1):**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Gdzie występuje |
|---|---|---|---|---|---|---|---|---|
| Karta środowiska | Duża karta wejścia (`.dn-karta-srodowiska`) | Wejście do przestrzeni roboczej jednego z czterech środowisk | Duży panel, jeden z czterech w siatce równej wagi — główny element wizualny strony | 1 | Widoczna bez interakcji | Spoczynek (powierzchnia neutralna, bez wstęgi) · wskazanie kursorem (wstęga sygnałowa na obrysie, uniesienie) · aktywny/fokus (wstęga sygnałowa + pierścień fokusu klawiatury) | Kliknięcie zamyka stronę główną i otwiera przedsionek wybranego środowiska (etap 6) | Wyłącznie strefa 1 strony głównej |
| Godło środowiska | Symbol graficzny właściwy środowisku, u góry karty | Szybka identyfikacja wizualna środowiska niezależnie od tekstu | Mały element graficzny, wewnątrz karty | 1 | Widoczne bez interakcji | Statyczny (dziedziczy kolor `currentColor` z karty) | Brak interakcji własnej — część obszaru klikalnego karty | Karta środowiska |
| Tytuł środowiska | Nazwa środowiska złożona krojem nagłówkowym (Cormorant Garamond) | Nazwanie środowiska | Tekst średniej wagi, największy krój na karcie | 1 | Widoczny bez interakcji | Statyczny | Brak interakcji własnej | Karta środowiska |
| Opis trybu pracy | Jednozdaniowe streszczenie przeznaczenia środowiska | Doprecyzowanie, czym różni się to środowisko od pozostałych trzech | Mały tekst pod tytułem, krój bazowy | 1 | Widoczny bez interakcji | Statyczny | Brak interakcji własnej | Karta środowiska |

---

### 3.6 Etap 6 — Przedsionek środowiska

| Aspekt | Treść |
|---|---|
| Warunek wejścia | Kliknięcie karty środowiska na stronie głównej (etap 5), przełączenie środowiska kontrolką trybów na pasku albo powrót do środowiska odwiedzonego wcześniej w tej sesji połączenia |
| Co się dzieje | Strona główna zostaje zastąpiona **przedsionkiem środowiska** — widokiem wejściowym, który przedstawia zawartość środowiska i pozwala wybrać moduł. Przestrzeń robocza nie otwiera się w tym etapie: **żaden moduł nie jest wskazany, więc żaden nie zostaje otwarty**. Chat Window, Execution Loop Window i boczna nawigacja modułów pojawiają się dopiero po wyborze (etap 7) |
| Stan okna | Rama okna (szyna nawigacji, belka tytułowa, wstążka) tożsama z przestrzenią roboczą. Poniżej dwa rejony: **szyna sesji** przy lewej krawędzi (wykaz sesji Operatora w tym środowisku) oraz **płótno przedsionka** — nagłówek środowiska, kafle modułów i listwa działań środowiska. Pas kart sesji jest nieobecny: karty należą do przestrzeni roboczej, nie do przedsionka |
| Wynik | Oczekiwanie na wybór modułu (etap 7), wznowienie sesji z szyny (przejście wprost do etapu 8) albo działanie środowiskowe z listwy — nowy projekt, konfiguracja środowiska, ustawienia |

**Zasada architektoniczna.** Wejście do środowiska nie jest wejściem do modułu.
Środowisko jest zbiorem modułów, nie jednym z nich; otwarcie modułu pierwszego
w kolejności byłoby decyzją podjętą za Operatora i zafałszowałoby hierarchię
platformy. Przedsionek środowiska pełni wobec modułów tę samą rolę, którą
Centrum dowodzenia pełni wobec środowisk: pokazuje zbiór, nie wybiera z niego.
Reguła obowiązuje wszystkie cztery środowiska bez wyjątku.

**Makieta tekstowa — przedsionek środowiska (przykład: TalkIn):**

```
 ┌───────────────────────────────────────────────────────────────────────────┐
 │ ☰  ◈ Danaco Console  [TalkIn][WorkSpace][CodeStudio][MultitaskingAI]      │ ← szyna + belka
 │    Danaco Console › TalkIn            [ 🔍 szukaj ]      ● 1 w tle  ◐  ◉  │
 ├─────────────────────┬─────────────────────────────────────────────────────┤
 │ TWOJE SESJE   3 · 7 │  ŚRODOWISKO PRACY                                    │
 │ [Czynne][Wszystkie] │  TalkIn                                              │
 │ ─────────────────── │  Myśl. Analizuj. Rozumiej.                           │
 │ Raport kwartalny    │  Wiedza, komunikacja i praca z treścią.              │
 │  ● Redakcja rozdz.  │                                                       │
 │    Studio · pracuje │  MODUŁY ŚRODOWISKA                                    │
 │    Streszczenie     │  ┌────────┬────────┬────────┬────────┬────────┐      │
 │    Studio · 12 min  │  │ Studio │Library │Browser │Research│Translate│     │
 │ Rejestr umów        │  ├────────┼────────┼────────┼────────┼────────┘      │
 │    Klasyfikacja     │  │Round-  │Assist- │Work-   │Agents  │               │
 │    Kary umowne      │  │table   │ant     │space   │        │               │
 │ Materiały reg.      │  └────────┴────────┴────────┴────────┘               │
 │    Przegląd zmian   │                                                       │
 │    Przekład noty    │  ┌───────────────────────────────────────────────┐   │
 │ ─────────────────── │  │ [+ Nowy projekt] │ Konfiguracja │ Ustawienia  │   │
 │ [Pełna historia]    │  └───────────────────────────────────────────────┘   │
 └─────────────────────┴─────────────────────────────────────────────────────┘
   szyna sesji             płótno przedsionka — trzy linie o malejącej masie:
   (lewa krawędź)          nagłówek → kafle modułów → listwa działań
```

**Trzy linie płótna.** Układ powtarza hierarchię strony głównej: masa maleje
z góry na dół, a forma prezentacji oddziela rodzaje działania. Nagłówek
środowiska niesie tożsamość; kafle modułów są punktami wejścia w pracę; listwa
gromadzi działania dotyczące środowiska jako całości, nie pojedynczego modułu.

**Tabela elementów interfejsu — przedsionek środowiska:**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Gdzie występuje |
|---|---|---|---|---|---|---|---|---|
| Szyna sesji | Stała kolumna przy lewej krawędzi, zatytułowana „Twoje sesje”, z sesjami zgrupowanymi projektami | Wznowienie pracy przerwanej oraz orientacja w tym, co dzieje się w środowisku | Kolumna nawigacyjna pełnej wysokości, waga średnia | 1 | Widoczna bez interakcji | Pozycja: domyślna · wskazanie kursorem · wybrana · ze wskaźnikiem pracy w tle | Kliknięcie pozycji otwiera przestrzeń roboczą tej sesji wraz z jej modułem — przejście wprost do etapu 8 | Przedsionek każdego środowiska |
| Filtr zakresu sesji | Grupa trzech kontrolek: czynne · wszystkie · zakończone | Zawężenie wykazu sesji do stanu, który Operator chce zobaczyć | Trzy małe kontrolki w nagłówku szyny | 1 | Widoczny bez interakcji | Domyślny · wskazanie kursorem · wybrany | Przeładowuje wykaz szyny bez zmiany płótna | Przedsionek każdego środowiska |
| Nagłówek środowiska | Godło środowiska, nazwa, motto i zdanie o przedmiocie pracy | Potwierdzenie, w którym środowisku znalazł się Operator, i przypomnienie jego zakresu | Blok o wadze głównej, krój nagłówkowy w największym stopniu | 1 | Widoczny bez interakcji | Statyczny | Brak interakcji własnej | Przedsionek każdego środowiska |
| Kafel modułu | Kafel z ikoną modułu, nazwą, zdaniem o przeznaczeniu i miarą bieżącego obłożenia | Wybór modułu, w którym rozpocznie się praca | Kafel w siatce, waga pośrednia — lżejszy od karty środowiska na stronie głównej | 1 | Widoczny bez interakcji | Domyślny · wskazanie kursorem · fokus | Otwiera przestrzeń roboczą środowiska z tym modułem jako wiodącym (etap 7 → 8) | Przedsionek każdego środowiska; liczba kafli równa liczbie modułów środowiska |
| Kafel modułu Workspace | Kafel na równych prawach z pozostałymi | Wejście w moduł organizacji projektów środowiska | Jak wyżej — **kafel, nie otwarte okno**; przedsionek nie osadza w sobie żadnego okna operacyjnego | 1 | Widoczny bez interakcji | Jak wyżej | Jak wyżej | Przedsionek TalkIn, WorkSpace, CodeStudio, MultitaskingAI |
| Listwa działań środowiska | Poziomy pas małych kontrolek pod kaflami | Działania dotyczące środowiska jako całości: założenie nowego projektu, konfiguracja środowiska, ustawienia | Listwa jednorzędowa, waga najniższa | 1 | Widoczna bez interakcji | Domyślny · wskazanie kursorem · fokus | „Nowy projekt” otwiera okno zakładania projektu; pozostałe pozycje otwierają Okno konfiguracji albo Okno ustawień jako nakładkę nad przedsionkiem | Przedsionek każdego środowiska |
| Kontrolka pełnej historii sesji | Przycisk w stopce szyny | Otwarcie pełnego wykazu sesji środowiska, także zakończonych i zarchiwizowanych | Przycisk zarysowany, pełna szerokość szyny | 1 | Widoczna bez interakcji | Domyślny · wskazanie kursorem · fokus | Otwiera menedżer sesji środowiska jako nakładkę | Przedsionek każdego środowiska |
| Rama okna (wspólna) | Szyna nawigacji, belka tytułowa, wstążka i pasek stanu tożsame z przestrzenią roboczą | Stały dostęp do godła (belka), środowisk i menu (szyna), wyszukiwania i kart (wstążka), konta i stanu (pasek stanu) | Cztery pasy ramy okna, waga niska | 1 | Widoczny bez interakcji | Statyczny | Znak w belce wraca do Centrum dowodzenia; szyna nawigacji przenosi do przedsionka innego środowiska | Przedsionek i przestrzeń robocza każdego środowiska |

**Czego przedsionek nie zawiera.** Chat Window, Execution Loop Window, bocznej
nawigacji modułów, pasa kart sesji ani żadnego okna operacyjnego. Wszystkie
te elementy należą do przestrzeni roboczej i pojawiają się dopiero po wyborze
modułu. Przedsionek nie jest przestrzenią roboczą o pustej treści — jest
odrębnym widokiem o własnym przeznaczeniu.

---

### 3.7 Etap 7 — Wybór modułu (lub roli / sekcji orkiestracji)

| Aspekt | Treść |
|---|---|
| Warunek wejścia | Operator znajduje się w przedsionku środowiska (etap 6) albo w przestrzeni roboczej z zamiarem zmiany modułu bieżącej karty |
| Co się dzieje | **Z przedsionka:** kliknięcie kafla modułu. Przedsionek ustępuje przestrzeni roboczej środowiska — pojawiają się boczna nawigacja modułów (albo panel orkiestracji), pas kart sesji, Chat Window i zestaw okien wybranego modułu; otwiera się pierwsza karta sesji. **Z przestrzeni roboczej:** kliknięcie pozycji w bocznej nawigacji albo sekcji w panelu orkiestracji. Klient wysyła polecenie otwarcia modułu; serwer przeładowuje kontekst karty sesji do nowego modułu — zestaw okien poprzedniego modułu (jeśli był otwarty) znika, w jego miejsce ładuje się zestaw okien nowego modułu |
| Stan okna | Krótkotrwały stan przejściowy (ładowanie zestawu okien modułu) między kliknięciem a wyświetleniem okna operacyjnego; boczna nawigacja / panel orkiestracji pozostaje widoczna i niezmieniona przez cały czas tego przejścia |
| Wynik | Pozycja modułu przyjmuje stan wybrany (podświetlenie); kolumny obszaru roboczego wypełniają się zestawem okien operacyjnych nowego modułu (etap 8) |

**Zasada wspólna:** przeładowanie obejmuje wyłącznie zestaw okien operacyjnych, kontekst czatu i narzędzia dedykowane modułowi. Nie obejmuje: samej karty sesji (pozostaje tą samą kartą), bocznej nawigacji lub panelu orkiestracji (pozostają widoczne bez zmian), pozostałych otwartych kart (zachowują własny, niezależny stan), funkcji globalnych Mobile i Always On Display (nigdy nie podlegają przeładowaniu).

**Tabela elementów interfejsu — Wybór modułu:**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Gdzie występuje |
|---|---|---|---|---|---|---|---|---|
| Kafel modułu w przedsionku | Kafel z ikoną, nazwą modułu i zdaniem o przeznaczeniu | Pierwszy wybór modułu po wejściu do środowiska | Kafel w siatce, waga pośrednia | 1 | Widoczny bez interakcji | Domyślny · wskazanie kursorem · fokus | Zamyka przedsionek i otwiera przestrzeń roboczą z tym modułem jako wiodącym | Przedsionek każdego środowiska |
| Pozycja modułu w bocznej nawigacji | Wiersz listy (`.dn-karta--klikalna`) z nazwą modułu | Wybór jednego z modułów dostępnych w bieżącym środowisku | Mały element listy, jedna linia tekstu | 1 | Widoczna bez interakcji | Domyślny · wskazanie kursorem · wybrana (podświetlenie, `▸`) | Przeładowuje kolumnę dominującą bieżącej karty do zestawu okien tego modułu | Boczna nawigacja TalkIn (9 pozycji), WorkSpace (9 pozycji), CodeStudio (8 pozycji) |
| Sekcja panelu orkiestracji | Wiersz listy analogiczny do pozycji modułu, lecz nazwany funkcją sterowania (Zespoły, Role, Kolejki, Orkiestracja, Harmonogram i automatyki, Monitor procesu) | Wybór warstwy sterowania zespołem modeli i agentów | Mały element listy, jedna linia tekstu | 1 | Widoczna bez interakcji | Jak wyżej | Sekcja „Role” otwiera cztery okna robocze; pozostałe sekcje otwierają właściwe im okna sterowania (rozdz. 5.4) | Wyłącznie panel orkiestracji MultitaskingAI (6 pozycji) |
| Wskaźnik ładowania modułu | Krótkotrwały stan pośredni kolumny dominującej | Sygnalizacja trwającego przeładowania zestawu okien | Nakładka lub spinner w kolumnie dominującej | 1 | Wyświetlany samoczynnie na czas przeładowania | Ładowanie (krótkotrwałe) | Automatyczne zniknięcie po wyświetleniu okna operacyjnego | Przejście między modułami w tej samej karcie |

---

### 3.8 Etap 8 — Okno operacyjne (okno robocze)

| Aspekt | Treść |
|---|---|
| Warunek wejścia | Wybór modułu zakończony (etap 7) — **albo** powrót do karty sesji, w której moduł był już otwarty wcześniej |
| Co się dzieje | Obszar roboczy karty sesji wypełnia się zestawem okien operacyjnych właściwym wybranemu modułowi, rozmieszczonym w kolumnach sąsiadujących poziomo. Chat Window — główne okno komunikacji Użytkownik ↔ Wykonawca, wspólne wszystkim piętnastu modułom — zajmuje lewą kolumnę na pełnej wysokości obszaru roboczego i rekonfiguruje się do kontekstu modułu. Execution Loop Window — okno pętli wykonawczej Koordynator ↔ Wykonawca — otwiera się w kolumnie sąsiadującej z Chat Window. Kolumna dominująca po prawej mieści okna właściwe modułowi; okna pomocnicze otwierają się jako rozszerzenia boczne, w kolejnych kolumnach po prawej stronie obszaru roboczego. Tu zaczyna się praca właściwa — dotychczasowe etapy były wyłącznie nawigacją do tego miejsca |
| Stan okna | Pełny, interaktywny zestaw okien w układzie pionowym (podział lewa–prawa): Chat Window w lewej kolumnie stałej, Execution Loop Window w kolumnie sąsiadującej, obszar roboczy modułu (edytor, lista, monitor — zależnie od modułu) w kolumnie dominującej, okna pomocnicze w kolumnach bocznych po prawej. Regulacji podlega wyłącznie szerokość kolumn. Zawartość aktualizowana na żywo kanałem WebSocket tam, gdzie okno monitoruje proces lub odbiera strumień odpowiedzi Wykonawcy. Stan procesu sesji jest trwały po stronie serwera — rozłączenie klienta nie zamyka okna |
| Wynik | Praca w module; dalsze przejścia z tego punktu: zmiana modułu w tej samej karcie (powrót do etapu 7), nowa karta („+”, powrót do etapu 7 w nowej karcie), powrót do strony głównej (powrót do etapu 4, karta trwa w tle), uproszczone menu kontekstowe (szybka zmiana konfiguracji bieżącej sesji bez opuszczania okna) |

**Makieta tekstowa — Okno operacyjne (schemat ogólny, niezależny od modułu, stan spoczynku):**

```
 ═══════════════════════════════════════════════════════════════════════════
  Boczna    │ Chat Window          │ Obszar roboczy modułu    │ Panel
  nawigacja │ Użytkownik ↔         │ nazwa okna · moduł   [⋮] │ pomocniczy
  modułów   │ Wykonawca            │                          │ (rozszerzenie
            │                      │ treść właściwa modułowi  │  boczne)
  [Danaco   │ historia rozmowy     │ (edytor / lista /        │
   Console] │                      │  monitor / podgląd)      │ zawartość
  [Ubuntu]  │ [ polecenie…    ] ➤  │                          │ właściwa
  [Fable 5] │                      │ aktualizacja na żywo     │ oknu
  [Ultra]   │ ──────────────────   │ (WebSocket) tam, gdzie   │ pomocniczemu
            │ Execution Loop       │ okno monitoruje proces   │
  Agent ▼   │ Koordynator ↔        │                          │
            │ Wykonawca            │                          │
            │ zlecenie · zadania · │                          │
            │ kontrola · przebieg  │                          │
 ═══════════════════════════════════════════════════════════════════════════
  kolumna     kolumna stała          kolumna dominująca         kolumna
  nawigacji   pełnej wysokości                                  boczna
```

Kolumny sąsiadują poziomo; regulacji podlega wyłącznie ich szerokość. Kolejne okna pomocnicze modułu otwierają się jako dalsze rozszerzenia boczne, po prawej stronie obszaru roboczego, i po zamknięciu znikają całkowicie z przestrzeni roboczej (rozdz. 3a).

**Uproszczone menu kontekstowe** — rozwinięcie kontekstowe warstwy 3, wywoływane kontrolką `⋮` w nagłówku dowolnego okna operacyjnego, bez opuszczania trwającej sesji:

```
 ┌─ Uproszczone menu kontekstowe (⋮) ────────────────┐
 │  Zachowanie modeli   ›                            │
 │  Tożsamość modelu    ›                            │
 │  Historia / pamięć   ›                            │
 │  Izolacja tej karty  ›                            │
 │  ──────────────────────────                       │
 │  Otwórz pełne okno konfiguracji →                 │
 └───────────────────────────────────────────────────┘
   zmiana obejmuje wyłącznie warstwę sesji bieżącej karty —
   wartości globalne i pozostałe karty pozostają nienaruszone;
   po zamknięciu menu znika całkowicie z przestrzeni roboczej
```

**Chat Window i Execution Loop Window w przepływie — zachowanie stałe:**

| Zdarzenie przepływu | Chat Window (Użytkownik ↔ Wykonawca) | Execution Loop Window (Koordynator ↔ Wykonawca) |
|---|---|---|
| Otwarcie modułu (etap 7→8) | Otwiera się jako pierwsze okno zestawu, zawsze w lewej kolumnie obszaru roboczego, w tym samym miejscu układu w każdym module i środowisku | Otwiera się w kolumnie sąsiadującej z Chat Window w chwili przyjęcia zlecenia do wykonania; prezentuje dekompozycję zlecenia i kolejkę zadań |
| Zmiana modułu w tej samej karcie | Pozostaje otwarte i zachowuje położenie; kontekst rozmowy rekonfiguruje się do nowego modułu, historia karty pozostaje dostępna | Zamyka się wraz z zestawem okien poprzedniego modułu i otwiera ponownie w kolumnie sąsiadującej dla zlecenia nowego modułu; stan pętli poprzedniego zlecenia trwa po stronie serwera |
| Przełączenie karty sesji | Wyświetla historię rozmowy karty wybranej; historia karty opuszczonej pozostaje nienaruszona | Wyświetla stan pętli wykonawczej karty wybranej; pętla karty opuszczonej biegnie dalej po stronie serwera |
| Powrót do strony głównej (przejście 18) | Znika z widoku wraz z całą kartą; strumień rozmowy trwa po stronie serwera | Znika z widoku; pętla wykonawcza biegnie dalej, wskaźnik pracy `●` pozostaje widoczny na karcie sesji |
| Zamknięcie Execution Loop Window | Bez zmian — pozostaje otwarte, zajmuje zwolnioną szerokość | Znika całkowicie z przestrzeni roboczej; pętla biegnie dalej po stronie serwera, okno przywracane jednym kliknięciem znacznika przebiegu albo poleceniem języka naturalnego w Chat Window |
| Rozłączenie i ponowne połączenie (przejścia 22–23) | Odtwarza pełną historię rozmowy karty w lewej kolumnie | Odtwarza stan pętli: zlecenie, kolejkę zadań, wyniki kontroli jakości i wskaźniki przebiegu |

**Stan okna komunikacji na poziomie kontraktu.** Chat Window i Execution Loop Window są, każde z osobna, encją „okno komunikacji” obszaru `window` kontraktu (`budowa/shared/contract.json`), niosącą pole `status` wyliczenia `WindowStatus` — dwuwartościowe: `open` (okno otwarte) i `closed` (okno zamknięte) — oraz pole `role` wyliczenia `WindowRole`: `executor` (okno wykonawcy — Chat Window; zakończenie tury wybudza koordynatora), `coordinator` (okno koordynatora — Execution Loop Window; przekazuje zlecenie przez narzędzie) albo `standalone` (okno samodzielne, poza pętlą). Stan `closed` opisany w kolumnach powyżej („znika z widoku”, „znika całkowicie z przestrzeni roboczej”) jest zawsze stanem **widoku klienta** — proces sesji leżący u podstaw okna trwa po stronie serwera niezależnie od stanu `open`/`closed` obserwowanego przez dowolne podłączone urządzenie (rozdz. 6).

```
Diagram przejść stanu — WindowStatus (obszar `window` kontraktu)

        window.create
   (zaklada okno komunikacji w sesji)
              │
              ▼
        ┌───────────┐   window.close     ┌────────────┐
        │   open    │ ──────────────────►│   closed   │
        │ (otwarte) │                    │ (zamknięte)│
        └─────┬─────┘ ◄───────────────── └────────────┘
              │             window.create
              │             (ponowne założenie —
              │              np. przywrócenie z ●
              │              znacznika przebiegu)
              │
   window.update / window.action / window.state.get / window.handoff
   (parametry, akcje panelu, odczyt stanu, przekazanie zlecenia —
    nie zmieniają wartości status; okno pozostaje `open`)
```

Zamknięcie widoku (`window.close`) nigdy nie jest wywoływane jako skutek rozłączenia klienta ani przełączenia karty sesji — oba te zdarzenia pozostawiają okno w stanie `open` po stronie serwera (kolumny „Rozłączenie i ponowne połączenie”, „Przełączenie karty sesji” powyżej); `window.close` odpowiada wyłącznie świadomemu zamknięciu kolumny przez Operatora (rozdz. 3a, „Znikanie po zamknięciu”) albo zakończeniu zlecenia pętli wykonawczej. Ponowne otwarcie okna uprzednio zamkniętego (np. przywrócenie Execution Loop Window znacznikiem przebiegu) zakłada nową encję `window.create`, nie przywraca poprzedniej encji `closed`.

**Tabela elementów interfejsu — Okno operacyjne:**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Gdzie występuje |
|---|---|---|---|---|---|---|---|---|
| Nagłówek okna | Pasek tytułowy pojedynczego okna operacyjnego | Identyfikacja okna i jego przynależności do modułu | Wąski pasek nagłówkowy kolumny okna | 1 | Widoczny bez interakcji | Statyczny | Brak interakcji własnej poza kontrolką rozwinięcia kontekstowego | Każde okno operacyjne |
| Chat Window — pole poleceń | Pole tekstowe wieloliniowe z przyciskiem wysłania | Wprowadzenie polecenia języka naturalnego dla Wykonawcy w kontekście bieżącego modułu | Pole szerokości lewej kolumny, waga średnia — stały punkt interakcji | 1 | Widoczne bez interakcji | Domyślny · fokus · wypełnione · ładowanie (podczas generowania odpowiedzi) | Wysłanie polecenia dodaje wpis do historii i otwiera strumień odpowiedzi na żywo | Każdy z piętnastu modułów i każde środowisko (okno stałe) |
| Chat Window — historia rozmowy | Przewijana lista wymienionych wiadomości | Prezentacja przebiegu rozmowy Użytkownik ↔ Wykonawca w bieżącym kontekście | Panel wypełniający lewą kolumnę, stała, pełna wysokość obszaru roboczego | 1 | Widoczna bez interakcji | Pusta (nowa karta) · z treścią · aktualizacja na żywo (strumień) | Przewijanie, zaznaczenie fragmentu jako przedmiotu dalszej operacji | Każdy z piętnastu modułów i każde środowisko |
| Przycisk „Wyślij” | Mały przycisk ikonowy (`.dn-btn-ikona`) przy polu poleceń | Zatwierdzenie i wysłanie polecenia | Mała ikonka (36×36 px) | 1 | Widoczny bez interakcji | Domyślny · wskazanie kursorem · aktywny · ładowanie | Zawsze klikalny; przy pustym polu poleceń kliknięcie wyświetla krótki komunikat zamiast wysyłki; przy polu wypełnionym inicjuje wysyłkę polecenia i otwarcie strumienia odpowiedzi | Chat Window, każdy moduł |
| Execution Loop Window | Okno pętli wykonawczej Koordynator ↔ Wykonawca — zlecenie i jego dekompozycja na zadania, kolejka i stan zadań, wymiana komunikatów sterujących, wyniki kontroli jakości i decyzje o ponowieniu, wskaźniki przebiegu pętli | Prowadzenie pętli wykonawczej, koordynacja zadań, nadzór nad realizacją i kontrola przebiegu procesów | Kolumna sąsiadująca z Chat Window, pełna wysokość obszaru roboczego | 1 | Otwierane przyjęciem zlecenia do wykonania; przywracane jednym kliknięciem znacznika przebiegu albo poleceniem języka naturalnego w Chat Window | Bez zlecenia (kolumna nieobecna) · zlecenie w toku · wstrzymane · zakończone · z niezgodnością zgłoszoną przez kontrolę jakości | Kliknięcie zadania w kolejce otwiera jego szczegóły w tej samej kolumnie; zamknięcie okna usuwa kolumnę całkowicie, pętla biegnie dalej po stronie serwera | Każdy moduł i każde środowisko, w którym realizowane jest zlecenie |
| Sterowanie przebiegiem pętli | Zestaw akcji Execution Loop Window: wstrzymanie, wznowienie, przerwanie, korekta zlecenia | Ingerencja Użytkownika w przebieg pętli wykonawczej | Jeden element zbiorczy `Przebieg ▼` w nagłówku kolumny | 3 | Rozwinięcie elementu zbiorczego `Przebieg ▼`; równoważnie polecenie języka naturalnego w Chat Window | Zwinięty (spoczynek) · rozwinięty | Wybór akcji zmienia stan pętli i zwija listę; rozwinięcie znika całkowicie po użyciu | Execution Loop Window |
| Kontrolka rozwinięcia kontekstowego „⋮” | Mała ikonka w nagłówku okna operacyjnego | Otwarcie uproszczonego menu szybkiej zmiany konfiguracji bieżącej sesji | Mała ikonka | 3 | Kliknięcie `⋮` w nagłówku okna | Domyślny · wskazanie kursorem · otwarte (menu widoczne) | Otwiera menu z podzbiorem ustawień warstwy sesji; po zamknięciu menu znika całkowicie z przestrzeni roboczej | Każde okno operacyjne |
| Znaczniki kontekstowe | Lekkie znaczniki środowiska, repozytorium, projektu, modelu i wykonawcy, na przykład `[Danaco Console] [Ubuntu] [Worktree] [Fable 5] [Ultra]` | Wskazanie kontekstu pracy bieżącej karty i dostęp do właściwych selektorów | Zwarty zestaw małych znaczników w kolumnie nawigacji | 1 (znacznik) / 2 (selektor) | Znacznik widoczny bez interakcji; kliknięcie znacznika otwiera odpowiedni selektor, który zwija się samoczynnie po wyborze | Domyślny · wskazanie kursorem · selektor otwarty | Wybór wartości aktualizuje znacznik i zwija selektor | Każde okno operacyjne |
| Okno pomocnicze modułu | Okno operacyjne dedykowane funkcji modułu (np. Studio Editor, Code Editor, Project Dashboard) | Realizacja konkretnego zadania właściwego modułowi | Kolumna boczna, otwierana jako rozszerzenie boczne po prawej stronie obszaru roboczego; regulacji podlega wyłącznie szerokość | 2 | Kliknięcie pozycji zestawu okien modułu, skrót klawiszowy albo polecenie języka naturalnego w Chat Window | Zestaw stanów właściwy typowi okna (edycyjne, monitorujące, budujące — zob. rozdz. 5) | Zależne od typu okna — zob. katalog szczegółowy w Specyfikacji okien operacyjnych; zamknięcie usuwa kolumnę całkowicie z przestrzeni roboczej | Zależnie od modułu — 1 do 5 okien pomocniczych na moduł |
| Narzędzia diagnostyczne niskiego poziomu i tryb administracyjny | Operacje na stanie procesu sesji, podglądzie kanału WebSocket i dziennikach wykonania | Diagnostyka przebiegu i administracja platformą | Bez reprezentacji w stanie spoczynku interfejsu | 4 | Polecenie języka naturalnego w Chat Window, skrót klawiszowy, wyszukiwarka funkcji albo tryb administracyjny | Wywołane · zamknięte | Wynik prezentowany w kolumnie bocznej, znikającej całkowicie po zamknięciu | Każde okno operacyjne; niewidoczne dla użytkownika podstawowego |

---

## 3a. Warstwy widoczności w przepływie okien

Zasadą nadrzędną interfejsu jest stopniowe ujawnianie funkcjonalności: jeżeli funkcja nie jest potrzebna do realizacji aktualnego zadania, nie jest widoczna. Przepływ okien realizuje tę zasadę na każdym etapie — liczba modułów, okien pomocniczych, paneli, ustawień i funkcji administracyjnych dostępnych na danym etapie nie wpływa na postrzeganą prostotę widoku. Złożoność platformy istnieje w architekturze i pozostaje niewidoczna w interfejsie do chwili wystąpienia potrzeby użycia danej funkcji.

Każdy element interfejsu opisany w tym dokumencie należy do dokładnie jednej z czterech warstw widoczności. Warstwę i sposób wywołania podają kolumny „Warstwa” i „Sposób wywołania” tabel elementów interfejsu w rozdziałach 3 i 5.

| Warstwa | Nazwa | Zawartość w przepływie okien | Sposób dostępu |
|---|---|---|---|
| 1 | Zawsze widoczna | Chat Window w lewej kolumnie, Execution Loop Window w kolumnie sąsiadującej, aktywne okno wiodące modułu, kontekst pracy, boczna nawigacja modułów lub panel orkiestracji, pasek kart sesji ze wskaźnikami stanu wykonania. Zajmuje ponad 80% powierzchni interfejsu | Widoczna bez interakcji |
| 2 | Widoczna na żądanie | Wybór modelu, wybór wykonawcy, wybór środowiska pracy modułu, wybór trybu pracy, poziom wysiłku, parametry przepływu pracy, okna pomocnicze modułu otwierane jako rozszerzenia boczne | Ikona, przycisk, przełącznik, znacznik kontekstowy; po użyciu element zwija się samoczynnie |
| 3 | Rozwinięcia kontekstowe | Uproszczone menu kontekstowe okna operacyjnego (`⋮`), sterowanie przebiegiem pętli wykonawczej, zestawy akcji karty sesji, ustawienia szybkie, warianty operacji | Menu kebab (`⋮`), menu hamburger (`☰`), menu kontekstowe, panel popover, lista rozwijana |
| 4 | Funkcje eksperckie | Najbardziej zaawansowane operacje przepływu: administracja procesami sesji, podgląd kanału WebSocket i dzienników wykonania, narzędzia diagnostyczne niskiego poziomu, tryby administracyjne | Polecenie języka naturalnego w Chat Window, skrót klawiszowy, wyszukiwarka funkcji, tryb administracyjny, konfiguracja roli. Użytkownik podstawowy nie widzi tych elementów |

**Znikanie po zamknięciu.** Okna i panele warstw 2 i 3 — okna pomocnicze otwierane jako rozszerzenia boczne, panele wysuwane, okna popover, panele kontekstowe, rozwinięcia menu — po zamknięciu znikają całkowicie z przestrzeni roboczej. Zamknięte okno pomocnicze nie pozostawia w układzie kolumny szczątkowej, zakładki ani uchwytu: zwolniona szerokość przypada kolumnom pozostałym. Zamknięcie widoku nie kończy procesu — proces sesji i pętla wykonawcza trwają po stronie serwera (rozdz. 6).

**Wywołanie funkcji warstwy 4.** Funkcje warstwy 4 uruchamiane są poleceniem języka naturalnego w Chat Window, skrótem klawiszowym, wyszukiwarką funkcji albo w trybie administracyjnym. Nie mają reprezentacji w stanie spoczynku interfejsu i nie zajmują miejsca w żadnej kolumnie do chwili wywołania; wynik ich działania prezentowany jest w kolumnie bocznej, znikającej całkowicie po zamknięciu.

**Zasada jednego kliknięcia.** Każda ukryta funkcja przepływu jest osiągalna jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego. Zagnieżdżanie funkcji głęboko w hierarchii menu jest zabronione — ukrycie zmniejsza chaos wizualny, nie utrudnia dostępu.

**Mechanizmy ukrywania stosowane w przepływie okien:**

| Mechanizm | Zastosowanie w przepływie | Rozdział |
|---|---|---|
| Menu progresywne | Zbiór jednorodnych wyborów prezentowany jako jeden element zwinięty (`Agent ▼`, `Przebieg ▼`); lista pozycji rozwija się po kliknięciu | rozdz. 3.8, 5.4 |
| Panele wysuwane | Okna pomocnicze modułu otwierane jako rozszerzenia boczne w kolumnach po prawej stronie obszaru roboczego; po zamknięciu kolumna znika całkowicie | rozdz. 3.8, 5.1–5.3 |
| Grupowanie logiczne akcji | Zamiast zestawu przycisków jeden element zbiorczy (`Operacje ▼`), którego rozwinięcie zawiera pełną listę akcji | rozdz. 3.8, 5.5 |
| Znaczniki kontekstowe | Środowisko, repozytorium, projekt, model i wykonawca jako lekkie znaczniki `[Danaco Console] [Ubuntu] [Worktree] [Fable 5] [Ultra]`; kliknięcie znacznika otwiera odpowiedni selektor | rozdz. 3.6, 3.8 |

**Makiety w stanie spoczynku.** Wszystkie makiety tekstowe tego dokumentu rysowane są w stanie spoczynku interfejsu: widoczne są wyłącznie elementy warstwy 1 oraz zwinięte wyzwalacze warstw 2–3 (znaczniki, `▼`, `⋮`, `☰`). Elementy warstw 2–4 opisane są wyłącznie w tabelach elementów okna, z podaniem warstwy i sposobu wywołania.

**Warstwy a etapy przepływu:**

| Etap | Elementy warstwy 1 | Elementy warstw 2–4 obecne na etapie |
|---|---|---|
| 1 — Okno startowe | Godło marki, wskaźnik stanu połączenia | Kontrolka ponowienia (warstwa 2, widoczna wyłącznie w stanie błędu) |
| 2 — Okno rejestracji i logowania | Zakładki, pola formularza, przycisk zatwierdzenia, przycisk „Pomiń” | Wybór metody logowania (warstwa 2), odzyskiwanie konta (warstwa 2, wywoływane odnośnikiem) |
| 3 — Przygotowanie środowiska pracy | Bryła warstw platformy, wykaz kroków przygotowania, pasek postępu, kontrolki „Pomiń przywracanie sesji” i „Przerwij i wyloguj” | — |
| 4 — Strona główna | Trzy strefy, rama okna | Wyszukiwarka funkcji (warstwa 4, skrót klawiszowy i pole wyszukiwania wstążki) |
| 5 — Wybór środowiska | Cztery karty środowisk | — |
| 6 — Przedsionek środowiska | Szyna sesji, kafle modułów, listwa działań środowiska | Menedżer sesji i okno zakładania projektu (warstwa 2) |
| 7 — Wybór modułu | Pozycja modułu lub sekcja panelu orkiestracji | — |
| 8 — Okno operacyjne | Chat Window, Execution Loop Window, okno wiodące modułu | Okna pomocnicze (warstwa 2), uproszczone menu kontekstowe i sterowanie przebiegiem pętli (warstwa 3), administracja procesami i diagnostyka niskiego poziomu (warstwa 4) |

---

## 4. Tabela przejść — zestawienie zbiorcze

Poniższa tabela zestawia wszystkie przejścia opisane w rozdziałach 2, 3 i 6, w kolejności zbliżonej do przebiegu głównego, z rozgałęzieniami zgrupowanymi bezpośrednio pod przejściem macierzystym.

| Nr | Grupa | Przejście (Z → Do) | Warunek | Co się dzieje | Stan okna docelowego |
|---|---|---|---|---|---|
| 1 | Rdzeń | Uruchomienie aplikacji → Okno startowe | Operator otwiera pakiet kliencki | Otwarcie natywnego okna; nawiązanie kanału WebSocket; `connection.hello` | Ładowanie (łączenie z serwerem) |
| 2 | Rdzeń | Okno startowe → Okno rejestracji i logowania | Brak ważnego tokenu urządzenia LUB faza budowy | Serwer potwierdza fazę wdrożenia; klient prezentuje formularz | Interaktywny formularz (logowanie domyślnie) |
| 2a | Wariant | Okno startowe → Przygotowanie środowiska → Strona główna (z pominięciem logowania) | Token urządzenia ważny, brak zapamiętanej aktywnej karty | Krok „Uwierzytelnienie” spełniony ważnym tokenem; `home.enter` | Stabilny, trzy strefy |
| 2b | Wariant | Okno startowe → Przygotowanie środowiska → Okno operacyjne (z pominięciem logowania i strony głównej) | Token ważny i klient zgłasza powrót do zapamiętanej aktywnej karty | `session.bind`; serwer przywraca pełny stan karty (układ, historia, kontekst) | Pełny, jak przed rozłączeniem |
| 3 | Rdzeń | Okno logowania → Strona główna (rejestracja) | Konto nie istnieje; dane wymagane podane i e-mail potwierdzony | Utworzenie jedynej encji Konto; wydanie tokenu urządzeniu | Stabilny, trzy strefy |
| 4 | Rdzeń | Okno logowania → Strona główna (logowanie) | Konto istnieje; login+hasło lub metoda dodatkowa poprawne | Weryfikacja przez serwer; wydanie tokenu urządzeniu | Stabilny, trzy strefy |
| 5 | Wariant | Okno logowania → Strona główna (Pomiń) | Ustawienie „Wymóg logowania” nieaktywne (domyślnie w fazie budowy, konfigurowalne przez Operatora w każdej fazie) | Wejście bez uwierzytelnienia; krok „Uwierzytelnienie” pomijany | Stabilny, trzy strefy |
| 6 | Wariant | Okno logowania → Odzyskiwanie konta → Okno logowania | Operator wskazuje utratę dostępu i podaje e-mail uwierzytelniający | Wysyłka drogi potwierdzenia; ustawienie nowego hasła; unieważnienie dawnych tokenów | Formularz logowania z nowym hasłem |
| 7 | Wariant | Okno logowania (dane niepoprawne) → Okno logowania | Weryfikacja negatywna | Komunikat błędu przy polu; formularz pozostaje wypełniony | Formularz, stan błędu |
| 7a | Rdzeń | Okno logowania / Pomiń → Przygotowanie środowiska pracy | Token wydany urządzeniu albo wejście bez konta | Klient odtwarza profil, uprawnienia, karty sesji, kanały modeli i magistralę kontekstu; widok pokazuje kroki i postęp | Widok przejściowy, bryła warstw w obrocie |
| 7b | Rdzeń | Przygotowanie środowiska pracy → Strona główna | Wszystkie kroki zakończone albo naciśnięcie „Pomiń przywracanie sesji” | Stan skompletowany; przedpokój otwiera się z aktualną liczbą sesji w tle | Stabilny, trzy strefy |
| 8 | Rdzeń | Strona główna → Przedsionek środowiska | Kliknięcie karty środowiska (strefa 1) | Zamknięcie widoku strony głównej; otwarcie przedsionka środowiska — szyny sesji, kafli modułów i listwy działań. Przestrzeń robocza nie otwiera się | Układ kolumnowy: boczna nawigacja modułów, Chat Window, kolumna dominująca; szyna nawigacji, belka tytułowa i pasmo kart sesji na wstążce jako rama okna |
| 9 | Rozgałęzienie | Strona główna → Okno budowy komponentu własnego | Kliknięcie kafla komponentu (strefa 2) | Otwarcie okna/nakładki budowy komponentu (Agent Builder i analogiczne — rozdz. 5.5) | Formularz budowy komponentu |
| 10 | Rozgałęzienie | Strona główna → Okno konfiguracji | Kliknięcie pozycji „Okno konfiguracji” (strefa 3) | Otwarcie nakładki modalnej nad bieżącym widokiem | Nakładka modalna, trzynaście zakresów w nawigacji |
| 11 | Rozgałęzienie | Strona główna → Aktywacja Mobile | Kliknięcie pozycji „Mobile” (strefa 3) | Aktywacja funkcji globalnej — bez nowej przestrzeni roboczej | Bieżący widok bez zmian + nakładka/panel Mobile |
| 12 | Rozgałęzienie | Strona główna → Aktywacja Always On Display | Kliknięcie pozycji „Always On Display” (strefa 3) | Aktywacja globalnego agenta towarzyszącego | Bieżący widok bez zmian + pływający avatar AOD |
| 13 | Rdzeń | Przedsionek środowiska → Wybór modułu | Kliknięcie kafla modułu w przedsionku (albo pozycji bocznej nawigacji, gdy przestrzeń robocza jest już otwarta) | Przedsionek ustępuje przestrzeni roboczej; otwiera się pierwsza karta sesji z zestawem okien wybranego modułu | Krótkotrwałe ładowanie zestawu okien modułu |
| 14 | Rdzeń | Wybór modułu → Okno operacyjne | Zestaw okien modułu załadowany | Chat Window w lewej kolumnie, Execution Loop Window w kolumnie sąsiadującej i okna właściwe modułowi wypełniają obszar roboczy | Pełny, interaktywny zestaw okien |
| 15 | Pętla | Okno operacyjne → Wybór modułu (ta sama karta) | Kliknięcie innej pozycji bocznej nawigacji | Zamknięcie okien poprzedniego modułu; otwarcie zestawu nowego modułu | Jak poz. 13–14 |
| 16 | Pętla | Przestrzeń robocza → nowa karta sesji → Wybór modułu | Kliknięcie kontrolki „+” paska kart | Nowa, pusta karta; poprzednie karty niezmienione w tle | Kolumna dominująca pusta — oczekiwanie na moduł |
| 17 | Pętla | Karta sesji → zamknięcie karty | Kliknięcie kontrolki „✕” | Zamknięcie widoku karty w oknie bieżącym | Karta znika z paska; pozostałe karty niezmienione |
| 18 | Pętla | Przedsionek / okno operacyjne → Strona główna | Kliknięcie godła w belce tytułowej | Powrót do centrum dowodzenia; karty bieżącego środowiska trwają w tle jako procesy serwera | Stabilny, trzy strefy |
| 19 | Pętla | Strona główna → Przedsionek środowiska (powrót) | Ponowne kliknięcie karty środowiska odwiedzonego wcześniej w tej sesji połączenia | Przedsionek z szyną sesji wskazującą karty trwające w tle; wznowienie wybranej sesji przywraca jej pełny stan | Przedsionek, szyna z sesjami czynnymi |
| 20 | Wariant | Okno operacyjne → rozwinięcie kontekstowe → Okno operacyjne | Kliknięcie kontrolki „⋮” w nagłówku okna | Szybka zmiana ustawień warstwy sesji bez opuszczania okna | Okno operacyjne niezmienione poza zaktualizowaną wartością ustawienia |
| 21 | Rozgałęzienie | Okno budowy komponentu własnego → zapis → Strona główna lub kontynuacja | Zatwierdzenie definicji komponentu | Zapis jako komponent własny w pamięci aplikacji; wytwór staje się zasobem wybieralnym | Strona główna albo dalsza edycja w tym samym oknie |
| 22 | Wariant | Rozłączenie klienta (dowolny moment) | Zamknięcie okna aplikacji, utrata sieci | Kanał WebSocket zamknięty; proces sesji i wszystkie procesy podrzędne trwają dalej po stronie serwera | Aplikacja zamknięta po stronie klienta; brak zmiany stanu po stronie serwera |
| 23 | Wariant | Ponowne połączenie → Okno startowe → stan zastany | Ponowne uruchomienie aplikacji na tym samym lub innym urządzeniu | Etapy 1–3 cyklu życia połączenia powtórzone; token urządzenia przedstawiony | Odtworzenie stanu sprzed rozłączenia (poz. 2b) albo strona główna (poz. 2a) |

---

## 5. Ścieżki szczegółowe

Poniższe sześć podrozdziałów rozwija etapy 5–8 (rozdz. 3.5–3.8) dla każdego z czterech środowisk z osobna oraz dla dwóch pozostałych stref strony głównej. Każda ścieżka pokazuje pełny przebieg od karty środowiska (lub kafla/pozycji listwy) do konkretnego okna operacyjnego, na jednym reprezentatywnym przykładzie modułu.

### 5.1 Ścieżka TalkIn

Środowisko wiedzy, komunikacji i pracy z treścią. Boczna nawigacja udostępnia dziewięć modułów: Studio, Workspace, Browser, Research, Library, Translate, Roundtable, Assistant, Agents.

```
 STRONA GŁÓWNA
     │  Strefa 1 → karta „TalkIn”
     ▼
 PRZEDSIONEK ŚRODOWISKA TalkIn
     │  szyna sesji + siatka dziewięciu kafli modułów
     │  ┌─────────────┬─────────────┬──────────┬────────────┬──────────┐
     │  │ Studio      │ Workspace   │ Browser  │ Research   │ Library  │
     │  ├─────────────┼─────────────┼──────────┼────────────┼──────────┤
     │  │ Translate   │ Roundtable  │ Assistant│ Agents     │          │
     │  └─────────────┴─────────────┴──────────┴────────────┴──────────┘
     ▼  (przykład: wybór „Studio”)
 OKNO OPERACYJNE — moduł Studio (stan spoczynku, układ kolumnowy)
┌─────────────────────┬────────────────────────┬───────────────────────┐
│ Chat Window         │ Studio Editor          │ Tools Panel           │
│ Użytkownik ↔        │ treść dokumentu:       │ (rozszerzenie boczne) │
│ Wykonawca           │ PDF·DOCX·TXT·Markdown  │ operacje kontekstowe  │
│                     │                        │ Wykonawcy wobec       │
│ [ polecenie…  ] ➤   │                     [⋮]│ dokumentu             │
│ ──────────────────  │                        │                       │
│ Execution Loop      │                        │                       │
│ Koordynator ↔       │                        │                       │
│ Wykonawca           │                        │                       │
└─────────────────────┴────────────────────────┴───────────────────────┘
 kolumna stała         kolumna dominująca       kolumny boczne
 pełnej wysokości                               (Tools Panel, Diff/Grep
                                                 Panel, Preview Window,
                                                 Session Repository —
                                                 każde jako rozszerzenie
                                                 boczne po prawej)
     │
     │  „+” nowa karta w tym samym środowisku → wybór „Research”
     ▼
 DRUGA KARTA SESJI — moduł Research (równolegle do karty Studio)
┌─────────────────────┬────────────────────────┬───────────────────────┐
│ Chat Window         │ Research Workspace     │ Sources Manager       │
│ Użytkownik ↔        │ zakres · przebieg      │ źródła z metadanymi   │
│ Wykonawca           │ badania                │ (rozszerzenie boczne) │
│ ──────────────────  │                     [⋮]│                       │
│ Execution Loop      │                        │ Findings Panel ·      │
│ Koordynator ↔       │                        │ Report Builder ·      │
│ Wykonawca           │                        │ Export Panel —        │
│                     │                        │ kolejne kolumny boczne│
└─────────────────────┴────────────────────────┴───────────────────────┘

 Obie karty (Studio, Research) działają jednocześnie — pasek kart:
 [ ● Studio  ✕ ] [ Research  ✕ ] [ + ]
```

**Tabela przejść — ścieżka TalkIn:**

| Krok | Warunek | Co się dzieje | Stan okna |
|---|---|---|---|
| Strona główna → TalkIn | Kliknięcie karty „TalkIn” | Otwarcie przedsionka środowiska | Szyna sesji przy lewej krawędzi, kafle modułów na płótnie, listwa działań pod kaflami |
| TalkIn → moduł Studio | Kliknięcie pozycji „Studio” | Załadowanie zestawu: Studio Editor, Tools Panel, Diff/Grep Panel, Session Repository, Preview Window, Chat Window | Pełny zestaw okien, edytor pusty do chwili wczytania dokumentu |
| Karta Studio → nowa karta Research | Kliknięcie „+”, następnie „Research” w bocznej nawigacji nowej karty | Druga, niezależna karta z własnym kontekstem; karta Studio niezmieniona w tle | Dwie karty aktywne jednocześnie w pasku kart |

### 5.2 Ścieżka WorkSpace

Środowisko produktywności, organizacji i realizacji projektów. Boczna nawigacja udostępnia dziewięć modułów: Studio, Workspace, Browser, Research, Library, Roundtable, Design, Apps, Agents.

```
 STRONA GŁÓWNA
     │  Strefa 1 → karta „WorkSpace”
     ▼
 PRZEDSIONEK ŚRODOWISKA WorkSpace
     │  szyna sesji + siatka dziewięciu kafli modułów
     │  ┌─────────────┬─────────────┬──────────┬────────────┬──────────┐
     │  │ Studio      │ Workspace   │ Browser  │ Research   │ Library  │
     │  ├─────────────┼─────────────┼──────────┼────────────┼──────────┤
     │  │ Roundtable  │ Design      │ Apps     │ Agents     │          │
     │  └─────────────┴─────────────┴──────────┴────────────┴──────────┘
     ▼  (przykład: wybór modułu „Apps” — budowa produktu cyfrowego)
 OKNO OPERACYJNE — moduł Apps (stan spoczynku, układ kolumnowy)
┌─────────────────────┬────────────────────────┬───────────────────────┐
│ Chat Window         │ Product Builder        │ Architecture Designer │
│ Użytkownik ↔        │ punkt wejścia budowy   │ komponenty rozwiązania│
│ Wykonawca           │ produktu               │ i zależności          │
│                     │                     [⋮]│ (rozszerzenie boczne) │
│ [ polecenie…  ] ➤   │                        │                       │
│ ──────────────────  │                        │                       │
│ Execution Loop      │                        │                       │
│ Koordynator ↔       │                        │                       │
│ Wykonawca           │                        │                       │
│ zlecenie · zadania  │                        │                       │
└─────────────────────┴────────────────────────┴───────────────────────┘
 kolumna stała         kolumna dominująca       kolumny boczne
 pełnej wysokości                               (Architecture Designer,
                                                 Frontend Workspace,
                                                 Backend Workspace,
                                                 Deployment Panel —
                                                 każde jako rozszerzenie
                                                 boczne po prawej)
```

**Tabela przejść — ścieżka WorkSpace:**

| Krok | Warunek | Co się dzieje | Stan okna |
|---|---|---|---|
| Strona główna → WorkSpace | Kliknięcie karty „WorkSpace” | Otwarcie przedsionka środowiska | Szyna sesji przy lewej krawędzi, kafle modułów na płótnie, listwa działań pod kaflami |
| WorkSpace → moduł Apps | Kliknięcie pozycji „Apps” | Załadowanie: Product Builder, Architecture Designer, Frontend Workspace, Backend Workspace, Deployment Panel, Chat Window | Product Builder jako punkt wejścia; Deployment Panel dostępny od razu — kliknięcie przed zakończeniem prac frontend/backend wyświetla komunikat o niezakończonych etapach |
| Apps → moduł Workspace (projekt) | Operator otwiera dodatkowo moduł Workspace w nowej karcie, aby zarządzać projektem nadrzędnym | Project Dashboard, Instructions Panel, Context Memory, Project Library, Agent Manager | Druga karta niezależna; Agent Manager pozwala przypisać agentów z modułu Agents do tego projektu |

### 5.3 Ścieżka CodeStudio

Środowisko programistyczne. Boczna nawigacja udostępnia osiem modułów: Workspace, Roundtable, Design, Terminal, Developer, Diagnostics, Apps, Agents.

```
 STRONA GŁÓWNA
     │  Strefa 1 → karta „CodeStudio”
     ▼
 PRZEDSIONEK ŚRODOWISKA CodeStudio
     │  szyna sesji + siatka ośmiu kafli modułów
     │  ┌─────────────┬─────────────┬──────────┬────────────┐
     │  │ Workspace   │ Roundtable  │ Design   │ Terminal   │
     │  ├─────────────┼─────────────┼──────────┼────────────┤
     │  │ Developer   │ Diagnostics │ Apps     │ Agents     │
     │  └─────────────┴─────────────┴──────────┴────────────┘
     ▼  (przykład: wybór modułu „Developer”)
 OKNO OPERACYJNE — moduł Developer (stan spoczynku, układ kolumnowy)
┌────────────────────┬─────────┬────────────────────┬──────────────────┐
│ Chat Window        │ Project │ Code Editor        │ Git Panel        │
│ Użytkownik ↔       │ Tree    │ treść plików       │ status zmian     │
│ Wykonawca          │ struk-  │ źródłowych         │ historia rewizji │
│                    │ tura    │ (wsparcie          │ (rozszerzenie    │
│ [ polecenie…] ➤    │ katalo- │ Wykonawcy)      [⋮]│  boczne)         │
│ ───────────────    │ gów     │                    │                  │
│ Execution Loop     │         │                    │                  │
│ Koordynator ↔      │         │                    │                  │
│ Wykonawca          │         │                    │                  │
│ zlecenie · zadania │         │                    │                  │
└────────────────────┴─────────┴────────────────────┴──────────────────┘
 kolumna stała       kolumna   kolumna dominująca   kolumny boczne
 pełnej wysokości    boczna                         (Git Panel,
                                                     Build Output —
                                                     log budowania
                                                     i wynik testów
                                                     na żywo)
     │
     │  powiązanie międzymodułowe jawne (konfigurowalne)
     ├──► Terminal Tabs (moduł Terminal) — warstwa wykonawcza poleceń
     └──► Diagnostics Center (moduł Diagnostics) — analiza błędów
```

**Tabela przejść — ścieżka CodeStudio:**

| Krok | Warunek | Co się dzieje | Stan okna |
|---|---|---|---|
| Strona główna → CodeStudio | Kliknięcie karty „CodeStudio” | Otwarcie przedsionka środowiska | Szyna sesji przy lewej krawędzi, kafle modułów na płótnie, listwa działań pod kaflami |
| CodeStudio → moduł Developer | Kliknięcie pozycji „Developer” | Załadowanie: Project Tree, Code Editor, Git Panel, Build Output, Chat Window | Project Tree odzwierciedla strukturę repozytorium powiązanego z sesją |
| Developer → Terminal (powiązanie jawne) | Operator, świadomą decyzją, otwiera dodatkowo moduł Terminal (nowa karta lub powiązanie skonfigurowane) | Terminal Tabs, Output Console, Process Monitor | Powiązanie nie jest domyślne — wymaga jawnej konfiguracji, zgodnie z zasadą pełnej konfigurowalności |

### 5.4 Ścieżka MultitaskingAI

Środowisko orkiestracji autonomicznej pracy ciągłej. **Różni się strukturalnie od pozostałych trzech** — nie ma bocznej nawigacji modułów; w jej miejscu panel orkiestracji o sześciu sekcjach, nawigujący po rolach, kolejkach i orkiestracji zamiast po modułach.

```
 STRONA GŁÓWNA
     │  Strefa 1 → karta „MultitaskingAI”
     ▼
 PRZEDSIONEK ŚRODOWISKA MultitaskingAI
     │  szyna sesji + siatka sześciu kafli sekcji orkiestracji
     │  ┌─────────────────────────────┐
     │  │  Zespoły                    │  presety ról, powiązań, kolejek
     │  │  Role                       │  → 4 okna robocze (poniżej)
     │  │  Kolejki                    │  → Queue Manager (Automations)
     │  │  Orkiestracja               │  → Orchestrator (Automations)
     │  │  Harmonogram i automatyki   │  → Scheduler, Execution Monitor
     │  │  Monitor procesu            │  → Always On Display, Mobile
     │  └─────────────────────────────┘
     ▼  (wybór sekcji „Role”)
 CZTERY OKNA ROBOCZE RÓL — jednoczesny widok zespołu w kolumnach sąsiadujących
┌──────────────────┬────────────────┬────────────────┬───────────────────┐
│ Chat Window      │ Executor Chat  │ Executor Chat  │ Results Analyzer  │
│ Użytkownik ↔     │ (Executor 1)   │ (Executor 2)   │ (Executor 3 /     │
│ Wykonawca        │ wykonanie zadań│ drugi, równo-  │  Validator)       │
│                  │ z kolejki      │ legły tor pracy│ ocena wyników ·   │
│ [ polecenie…] ➤  │ Subagent       │ Subagent       │ zgłaszanie        │
│ ───────────────  │ Network (do 15)│ Network (do 15)│ niezgodności ·    │
│ Execution Loop   │ tryb współpracy│ niezależna ·   │ rozstrzyganie     │
│ Koordynator ↔    │ z Executorem 2 │ przekazywanie  │ konfliktów        │
│ Wykonawca        │ wg ustalenia   │ wyników ·      │                   │
│ (Coordinator     │ Koordynatora   │ naprzemienna · │ wcielenie         │
│  Chat)           │                │ iteracyjna     │ konfigurowalne:   │
│ plan etapów ·    │                │                │ Validator ·       │
│ podział pracy ·  │                │                │ Reviewer ·        │
│ budowa promptów ·│                │                │ Security Auditor ·│
│ stan kolejki     │                │                │ Architect ·       │
│ Przebieg ▼       │             [⋮]│             [⋮]│ Product Owner ·   │
│                  │                │                │ QA Lead ·         │
│                  │                │                │ Arbitrator        │
└──────────────────┴────────────────┴────────────────┴───────────────────┘
 kolumna stała      kolumna          kolumna          kolumna boczna
 pełnej wysokości   dominująca       sąsiadująca      (rozszerzenie)

 Sterowanie pętlą (start · stop · pauza · wznowienie · przekazanie ·
 powtórzenie · walidacja) rozwija się z elementu zbiorczego `Przebieg ▼`
 w nagłówku Execution Loop Window (warstwa 3) albo jest wywoływane
 poleceniem języka naturalnego w Chat Window.
     │
     │  po skonfigurowaniu integracji z modułem Automations —
     │  pętla pracy ciągłej 24/7/365 prowadzona przez Koordynatora
     ▼
 Silnik kolejek → Orkiestracja → Scheduler/Queue Manager/Orchestrator/
 Execution Monitor (moduł Automations) → kolejny przebieg (enqueue) ─┐
                                                                       │
        ◄──────────────────────────────────────────────────────────┘
```

**Tabela przejść — ścieżka MultitaskingAI:**

| Krok | Warunek | Co się dzieje | Stan okna |
|---|---|---|---|
| Strona główna → MultitaskingAI | Kliknięcie karty „MultitaskingAI” | Otwarcie przedsionka środowiska | Szyna sesji przy lewej krawędzi, kafle sześciu sekcji orkiestracji na płótnie, listwa działań pod kaflami |
| MultitaskingAI → sekcja „Role” | Kliknięcie sekcji „Role” | Otwarcie czterech okien roboczych jednocześnie: Coordinator Chat, dwa Executor Chat, Results Analyzer | Cztery okna widoczne równolegle — nie sekwencyjnie |
| Rola → przypisanie wykonawcy | Operator przypisuje agenta z modułu Agents (strefa 2) albo model bazowy do roli | Agent lub model staje się wykonawcą roli; wybór kanału modelu (API/CLI/SSH/HTTP) per rola | Okno roli aktywne, gotowe do sterowania przez Coordinatora |
| Role → sekcja „Harmonogram i automatyki” | Operator konfiguruje powiązanie z modułem Automations (świadoma decyzja, nie domyślna) | Scheduler i Execution Monitor modułu Automations stają się operacyjnym zapleczem harmonogramu | Praca ciągła 24/7/365 możliwa po tej konfiguracji |
| Rola → sekcja „Monitor procesu” | Kliknięcie sekcji „Monitor procesu” | Otwarcie widoku nadzoru; Always On Display może pełnić funkcję obserwatora lub operatora | Widok na żywo statusów przebiegów i hierarchii decyzji |

### 5.5 Ścieżka Strefy 2 — komponenty własne

Strefa 2 strony głównej udostępnia cztery kafle: Automations, Agents, Workspace, Assistant. Każdy otwiera okno budowy właściwego rodzaju komponentu własnego — **nie** przestrzeń roboczą środowiska. Wytwór, po zapisie, trafia do pamięci aplikacji jako zasób Operatora i staje się wybieralny operacyjnie w środowiskach (etapy 7–8).

**Schemat wspólny czterech ścieżek strefy 2:**

```
 STRONA GŁÓWNA — Strefa 2
     │
     ├── kafel „Automations” ──► Workflow Builder (+ Scheduler, Queue Manager,
     │                            Orchestrator, Execution Monitor, Chat Window)
     ├── kafel „Agents”      ──► Agent Builder (rozwinięcie poniżej — najbardziej
     │                            rozbudowana ścieżka strefy 2)
     ├── kafel „Workspace”   ──► Project Dashboard (+ Instructions Panel,
     │                            Context Memory, Project Library, Agent Manager)
     └── kafel „Assistant”   ──► Voice Console (+ Actions Monitor, Activity Feed)
                                        │
                                        ▼
                          ZAPIS JAKO KOMPONENT WŁASNY
                    (zasób Operatora w pamięci aplikacji)
                                        │
                    ┌───────────────────┼───────────────────┐
                    ▼                   ▼                   ▼
              wybieralny w        wpięty w sesji       (dla Agents:
              sesji modułu        modułu Developer     dodatkowo rola w
              właściwego          i innych (dla        panelu orkiestracji
              (etap 8)            Automations)          MultitaskingAI)
```

**Rozwinięcie — ścieżka Agents (przykład flagowy, wskazany wprost w zakresie dokumentu):**

Moduł Agents ma dwa punkty wejścia prowadzące do tego samego okna Agent Builder i tej samej definicji agenta: strefa 2 strony głównej (tworzenie i zarządzanie agentami jako takimi) oraz okno modułowe wewnątrz środowiska (tworzenie lub edycja agenta w toku bieżącej pracy nad zadaniem, dostępne w TalkIn, WorkSpace i CodeStudio).

```
 STRONA GŁÓWNA — Strefa 2, kafel „Agents” („Skonfiguruj agenta”)
     │
     ▼
 AGENT BUILDER  (okno nadrzędne — scala definicję w komponent własny)
     │
     ├──────────────┬────────────────────────┬──────────────┐
     ▼              ▼                        ▼              ▼
 Model            Skills                 Connectors      Permissions
 Configuration    Manager                Manager         Center
 model bazowy,    dobór umiejętności     pluginy,        zakres uprawnień
 kanał modelu                            konektory       agenta
     │              │                        │              │
     └──────────────┴───────────┬────────────┴──────────────┘
                                 ▼
                  + tożsamość, instrukcje systemowe,
                    konfiguracja pamięci (Agent Builder)
                                 ▼
                  ZAPIS JAKO KOMPONENT WŁASNY
                  (zasób Operatora w pamięci aplikacji)
                                 │
        ┌────────────────────────┼────────────────────────┐
        ▼                        ▼                        ▼
  Wykonawca zadań          Agent Manager             Rola w panelu
  w TalkIn / WorkSpace /   modułu Workspace           orkiestracji
  CodeStudio (wybór        — przypisanie do           środowiska
  komponentu w sesji)      projektu                   MultitaskingAI
                                                        (sekcja Role, 5.4)

  Chat Window (Użytkownik ↔ Wykonawca) — lewa kolumna, stała, pełna
  wysokość obszaru roboczego; wsparcie konwersacyjne budowy i testowania
  agenta, obecne przez cały proces w Agent Builder.
  Execution Loop Window (Koordynator ↔ Wykonawca) — kolumna sąsiadująca,
  otwierana na czas przebiegu testowego definicji agenta
```

**Tabela elementów interfejsu — Kafel komponentu własnego (strefa 2, wzorzec wspólny czterech kafli):**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Gdzie występuje |
|---|---|---|---|---|---|---|---|---|
| Kafel komponentu własnego | Kafel (`.dn-kafel`) | Wejście do okna budowy jednego rodzaju komponentu własnego | Średni panel, lżejszy wizualnie niż karta środowiska — brak godła w skali karty środowiska, brak kroju nagłówkowego | 1 | Widoczny bez interakcji | Spoczynek · wskazanie kursorem (tło subtelne, bez wstęgi sygnałowej na spoczynku) · aktywny/fokus | Kliknięcie otwiera okno nakładkowe albo dedykowane okno budowy komponentu (nie: przestrzeń roboczą środowiska) | Wyłącznie strefa 2 strony głównej |
| Ikona komponentu | Mały symbol graficzny właściwy rodzajowi komponentu | Szybka identyfikacja rodzaju (automatyka, agent, projekt, profil głosowy) | Mała ikonka, mniejsza skala niż godło środowiska | 1 | Widoczna bez interakcji | Statyczny | Brak interakcji własnej | Kafel komponentu własnego |
| Nazwa komponentu | Automations / Agents / Workspace / Assistant | Nazwanie rodzaju komponentu | Tekst średniej wagi, krój bazowy (nie nagłówkowy) | 1 | Widoczna bez interakcji | Statyczny | Brak interakcji własnej | Kafel komponentu własnego |
| Etykieta działania | Krótkie sformułowanie zorientowane na tworzenie („Zbuduj automatykę”, „Skonfiguruj agenta”, „Załóż projekt”, „Ustaw profil asystenta”) | Odróżnienie charakteru akcji od karty środowiska („wejdź” vs. „zbuduj”) | Mały tekst pod nazwą komponentu | 1 | Widoczna bez interakcji | Statyczny | Brak interakcji własnej | Kafel komponentu własnego |
| Agent Builder — okno nadrzędne | Formularz spinający cztery okna konfiguracyjne w jeden zapisywalny wytwór | Skompletowanie pełnej definicji agenta | Duże okno / dedykowany widok budowy | 1 | Kliknięcie kafla „Agents” albo polecenie języka naturalnego w Chat Window | Pusty (nowy agent) · wypełniany · gotowy do zapisu · zapisany | Zapis tworzy komponent własny „agent”, wybieralny dalej w środowiskach i rolach | Strefa 2 (kafel „Agents”) oraz okno modułowe Agents w środowiskach |
| Przycisk „Zapisz agenta” | Przycisk główny CTA w Agent Builder | Zatwierdzenie i zapisanie definicji agenta jako komponentu własnego | Duży, pojedynczy przycisk akcentowany | 1 | Widoczny bez interakcji | Domyślny · wskazanie kursorem · aktywny · zapisano (potwierdzenie) | Zawsze klikalny; przy niekompletnej definicji kliknięcie zapisuje agenta i wyświetla ostrzeżenie wskazujące brakujące elementy — zapis nie jest blokowany; dymek powiadomienia potwierdzający; agent staje się wybieralny | Agent Builder |

### 5.6 Ścieżka Strefy 3 — konfiguracja, ustawienia, funkcje globalne

Strefa 3 strony głównej udostępnia listwę ustawień: pozycję „Okno konfiguracji” oraz dwie funkcje globalne, Mobile i Always On Display. W odróżnieniu od stref 1 i 2, wszystkie trzy pozycje są dostępne **z każdego miejsca platformy**, nie tylko ze strony głównej — obecność w listwie jest jednym z kilku równoważnych punktów dostępu.

```
 STRONA GŁÓWNA — Strefa 3                    (dostępne też z KAŻDEGO środowiska
     │                                         i KAŻDEGO okna operacyjnego)
     ├── „Okno konfiguracji” ──► OKNO KONFIGURACJI
     │                            nakładka modalna nad bieżącym widokiem
     │                            ┌─────────────────────┬────────────────────┐
     │                            │ NAWIGACJA (13 zakr.)│ ZAWARTOŚĆ ZAKRESU  │
     │                            │ Aplikacja           │ pozycje ustawień   │
     │                            │ Procesy             │ z wartością,       │
     │                            │ Akcje               │ warstwą (globalna  │
     │                            │ Zachowanie modeli   │ / środowisko /     │
     │                            │ Tożsamość modeli    │ projekt / sesja)   │
     │                            │ Prompty systemowe   │ i objaśnieniem[?]  │
     │                            │ Rozszerzenia        │                    │
     │                            │ Integracje          │                    │
     │                            │ Historia            │                    │
     │                            │ Pamięć              │                    │
     │                            │ Izolacja  ─────────►│ (rozwinięcie niżej)│
     │                            │ Karty sesji         │                    │
     │                            │ Komponenty własne   │                    │
     │                            └─────────────────────┴────────────────────┘
     │                                    │  pozycja „Izolacja”
     │                                    ▼
     │                     OKNO KONFIGURACJI PUNKTÓW IZOLACJI (trzy panele)
     │                     ┌─────────────────┬────────────────────┬────────────────┐
     │                     │ SELEKTOR        │ MACIERZ IZOLACJI   │ PROFIL I       │
     │                     │ ZASIĘGU (7 poz.)│ (3 kontekst. +     │ PODGLĄD        │
     │                     │ globalny·środo- │  8 technicznych,   │ zapisz/wczytaj/│
     │                     │ wisko·moduł·    │  każdy: przełącznik│ przypisz profil│
     │                     │ para modułów·   │  z objaśnieniem[?])│ warstwa: domyśl│
     │                     │ projekt·karta   │                    │ na / sesji     │
     │                     │ sesji·rola      │                    │ polityka efekt.│
     │                     └─────────────────┴────────────────────┴────────────────┘
     │
     ├── „Mobile” ──► AKTYWACJA funkcji globalnej Mobile
     │                 bez nowej przestrzeni roboczej; dostęp mobilny,
     │                 monitoring procesów, zarządzanie zadaniami z dowolnego miejsca
     │
     └── „Always On Display” ──► AKTYWACJA globalnego agenta towarzyszącego
                                   pływający avatar obecny jednocześnie wszędzie;
                                   proaktywne doradztwo, pomoc kontekstowa,
                                   komunikacja głosowa i tekstowa
```

Wśród pozycji zakresu „Aplikacja” stoi ustawienie „Wymóg logowania”, którym Operator określa — niezależnie od fazy wdrożenia — czy okno rejestracji i logowania (rozdz. 3.2) wymusza skuteczne uwierzytelnienie przed wejściem na platformę; domyślnie nieaktywne w fazie budowy, aktywne w fazie produkcyjnej, zmienialne w każdej chwili.

**Tabela elementów interfejsu — listwa ustawień (strefa 3):**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Gdzie występuje |
|---|---|---|---|---|---|---|---|---|
| Pozycja „Okno konfiguracji” | Element listwy (`.dn-btn-ikona` + etykieta) | Otwarcie pełnego zakresu ustawień platformy | Mały element paska, waga najniższa | 1 | Widoczna bez interakcji; równoważnie skrót klawiszowy i wyszukiwarka funkcji | Domyślny · wskazanie kursorem · fokus | Otwiera nakładkę modalną z trzynastoma zakresami ustawień | Listwa strony głównej i każde środowisko/okno operacyjne |
| Pozycja „Mobile” | Element listwy analogiczny | Aktywacja globalnego trybu dostępu mobilnego | Mały element paska | 1 | Widoczna bez interakcji | Jak wyżej | Aktywuje funkcję — nie tworzy nowej przestrzeni roboczej, nie przeładowuje niczego | Listwa strony głównej i każde środowisko/okno operacyjne |
| Pozycja „Always On Display” | Element listwy analogiczny | Aktywacja globalnego agenta towarzyszącego | Mały element paska | 1 | Widoczna bez interakcji | Jak wyżej | Aktywuje pływający avatar obecny nad całą platformą | Listwa strony głównej i każde środowisko/okno operacyjne |
| Nakładka modalna okna konfiguracji | Duży panel (`.dn-modal`) nad przyciemnionym tłem (token `--dn-nakladka` na `::backdrop`) | Konfiguracja dowolnej zależności lub zachowania platformy z jednego miejsca | Duży panel, dwudzielny układ (nawigacja + zawartość) | 2 | Kliknięcie pozycji listwy, rozwinięcie kontekstowe `⋮` okna operacyjnego, skrót klawiszowy albo wyszukiwarka funkcji; po zamknięciu znika całkowicie z przestrzeni roboczej | Otwarta · zamknięta (powrót do widoku pod spodem, niezmienionego) | Zamknięcie przywraca dokładnie stan widoku sprzed otwarcia | Wywoływana z listwy strony głównej albo uproszczonego menu kontekstowego okna operacyjnego |
| Przełącznik macierzy izolacji | Suwak (`.dn-suwak`) | Ustawienie „współdzielone/odrębne” (kontekst) albo „włączony/wyłączony” (technika) dla jednej pozycji izolacji | Mały element w wierszu tabeli/listy | 2 | Widoczny po otwarciu okna konfiguracji punktów izolacji | Domyślny (wyłączony/odrębny) · włączony/współdzielony · wskazanie kursorem · fokus | Zmiana natychmiast widoczna w podglądzie polityki efektywnej panelu prawego | Okno konfiguracji punktów izolacji |
| Objaśnienie kontekstowe `[?]` | Ikona z dymkiem (`.dn-tooltip`) | Wyjaśnienie działania ustawienia i jego wpływu na aplikację | Bardzo mała ikonka przy każdej pozycji ustawienia | 3 | Wskazanie kursorem lub fokus na ikonie `[?]`; dymek znika po odsunięciu wskaźnika | Domyślny (ukryty tekst) · wskazanie kursorem albo fokus (dymek widoczny) | Wyświetla treść objaśnienia, nie zmienia wartości ustawienia | Każda pozycja konfiguracji w oknie konfiguracji |

---

### 5.7 Ścieżki strefy pracy szyny nawigacji

Trzy pozycje strefy pracy prowadzą ścieżki niezależne od tego, gdzie Operator
właśnie jest. Wszystkie trzy są dostępne z każdego okna platformy.

| Wywołanie | Warunek | Co się dzieje | Okno docelowe |
|---|---|---|---|
| Nowa sesja — będąc w przestrzeni roboczej środowiska | Środowisko rozpoznane z kontekstu | Sesja powstaje bez pytania, jako kolejna karta tego środowiska | Przestrzeń robocza, nowa karta sesji |
| Nowa sesja — będąc poza środowiskiem | Środowisko nierozpoznane | Menu żąda wskazania środowiska; modułu nie pyta | Przedsionek wskazanego środowiska |
| Historia sesji | Zawsze dostępna | Otwiera nakładkę okna ustawień z wybraną sekcją „Historia sesji”; widok pod spodem pozostaje nietknięty | Nakładka rejestru sesji |
| Historia sesji → wejście w sesję czynną | Sesja pracuje po stronie serwera | Powrót do sesji w stanie, w jakim została pozostawiona | Przestrzeń robocza właściwego środowiska |
| Historia sesji → wejście w sesję zakończoną | Zapis sesji istnieje | Otwiera się nowa sesja z odtworzonym kontekstem poprzedniej | Przestrzeń robocza właściwego środowiska |
| Historia sesji → wejście w sesję archiwalną | Sesja w archiwum konta | Wymaga przywrócenia; dopiero potem otwiera się jak sesja zakończona | Nakładka rejestru (po przywróceniu) |
| Nowy projekt | Zawsze dostępny | Menu żąda wskazania środowiska; **przedsionek jest pomijany** | Okno zakładania projektu w Workspace wskazanego środowiska |
| Zakładanie projektu → „Załóż projekt i otwórz sesję” | Formularz wypełniony | Projekt powstaje, pierwsza sesja otwiera się w przedsionku środowiska | Przedsionek środowiska |
| Zakładanie projektu → „Załóż bez otwierania sesji” | Formularz wypełniony | Projekt czeka w Workspace środowiska; w historii sesji pojawia się dopiero po pierwszym otwarciu | Workspace środowiska |

**Zależności niewidoczne w oknach.** Katalog projektu jest jedynym miejscem
zapisu wytworów — poza nim sesja nie zapisze niczego bez osobnej zgody
Operatora. Instrukcja stała projektu obowiązuje każdą jego sesję i każdego
agenta; sesja może ją uzupełnić, nie może jej znieść. Profil izolacji kontekstu
rozstrzyga zasięg sięgania po dane (katalog projektu · biblioteka środowiska ·
sieć) i jest tym samym ustawieniem, które niesie kontrolka „Izolacja kontekstu”
na wstążce. Zakończenie sesji nie usuwa jej wytworów: dokumenty, wersje
i wpisy pamięci pozostają w projekcie.

---

## 6. Karty sesji i trwałość procesu w tle

Mechanizm kart sesji przenika etapy 6–8 i tłumaczy, dlaczego „stan okna” w tym dokumencie tak często rozróżnia „nową kartę” od „karty przywróconej”. Każda sesja — niezależnie od tego, czy jej karta jest obecnie widoczna — działa jako odrębny proces po stronie serwera.

```
 UTWORZENIE KARTY SESJI  (etap 7, kliknięcie „+” albo pierwszy wybór środowiska)
        │  serwer uruchamia odrębny proces (własny katalog roboczy, środowisko procesu)
        ▼
 ┌────────────────────┐   rozłączenie klienta     ┌──────────────────────────────┐
 │  AKTYWNA           │ ─────────────────────────►│  AKTYWNA BEZ KLIENTA         │
 │  karta widoczna,   │   (zamknięcie aplikacji,  │  proces trwa po stronie      │
 │  klient podłączony │    utrata sieci, powrót   │  serwera; karta niewidoczna  │
 │                    │    do strony głównej)     │  w oknie klienta             │
 └─────────┬──────────┘ ◄─────────────────────────└───────────────┬──────────────┘
           │              ponowne połączenie /                    │
           │              powrót do środowiska                    │
           │              (odtworzenie pełnego stanu:             │
           │               układ, historia, kontekst)             │
           └───────────────────────◄──────────────────────────────┘
```

| Zdarzenie | Wpływ na kartę bieżącą | Wpływ na pozostałe karty |
|---|---|---|
| Powrót do strony głównej (etap 4, przejście 18) | Znika z widoku; proces trwa dalej na serwerze | Bez zmian — trwają w tle na równi z kartą opuszczoną |
| Wybór innego środowiska na stronie głównej | Poprzednie środowisko i wszystkie jego karty pozostają aktywne w tle | Odtwarzają pełny stan przy powrocie do tego środowiska |
| Zamknięcie okna aplikacji (rozłączenie, przejście 22) | Kanał WebSocket zamknięty; stan zapisany trwale (SQLite + pliki treści) | Wszystkie procesy sesji tego Operatora trwają dalej po stronie serwera |
| Ponowne uruchomienie aplikacji (przejście 23) | Etapy 1–3 cyklu życia połączenia powtórzone | Karty odtwarzają się dokładnie w stanie sprzed rozłączenia |
| Zmiana modułu w tej samej karcie (etap 7→8 w pętli) | Przeładowanie zestawu okien tej jednej karty; Chat Window pozostaje otwarte w lewej kolumnie i rekonfiguruje kontekst | Bez wpływu na pozostałe karty |
| Zamknięcie Execution Loop Window albo okna pomocniczego | Kolumna znika całkowicie z przestrzeni roboczej, zwolniona szerokość przypada kolumnom pozostałym | Bez wpływu — proces pętli wykonawczej biegnie dalej po stronie serwera |

Ta trwałość jest bezpośrednią podstawą funkcji globalnej Mobile — pozwala monitorować i nadzorować procesy uruchomione w tle niezależnie od tego, które okno jest aktualnie otwarte na ekranie danego urządzenia, oraz podstawą pracy ciągłej 24/7/365 środowiska MultitaskingAI opisanej w rozdziale 5.4: role zespołu wykonują pracę również wtedy, gdy żadne urządzenie Operatora nie jest w danej chwili podłączone. Zamknięcie widoku Execution Loop Window nie przerywa pętli wykonawczej — Koordynator prowadzi ją dalej, a okno przywracane jest jednym kliknięciem znacznika przebiegu albo poleceniem języka naturalnego w Chat Window.

## 6a. Magistrala kontekstu — przekazywanie artefaktu między modułami

Magistrala kontekstu jest mechanizmem przekazywania artefaktu między modułami platformy. Występuje pod tą jedną nazwą w każdym środowisku i w każdym module; dokumenty środowiskowe opisują wyłącznie specyfikę właściwą danemu środowisku i odsyłają do niniejszego rozdziału.

**Przeznaczenie.** Magistrala przenosi artefakt powstały w jednym module do innego modułu, do Chat Window albo do Execution Loop Window, bez eksportu, zapisu pośredniego i bez opuszczania karty sesji. Przenosi wyłącznie kontekst — nie wykonuje pracy modułów: decyzję o działaniu na przekazanym artefakcie podejmuje Użytkownik w Chat Window albo Koordynator w pętli wykonawczej. Przedmiotem przekazania jest plik, zaznaczenie, wynik, błąd, fragment treści, źródło badawcze albo odnośnik.

**Sposób przekazania artefaktu.**

| Krok | Działanie | Warstwa | Sposób wywołania |
|---|---|---|---|
| Odłożenie | Artefakt zostaje odłożony do magistrali z okna operacyjnego modułu źródłowego | 3 | Pozycja „Do magistrali” w menu kontekstowym artefaktu albo przeciągnięcie na krawędź obszaru roboczego |
| Przegląd | Kolumna magistrali prezentuje odłożone pozycje z modułem pochodzenia, typem i znacznikiem czasu | 3 | Kliknięcie licznika magistrali w pasku kontekstu |
| Pobranie | Artefakt zostaje osadzony w module docelowym, w polu polecenia Chat Window albo jako wejście zadania pętli wykonawczej | 3 | Pozycja „Otwórz w…”, „Dołącz do rozmowy” albo przeciągnięcie pozycji z kolumny magistrali |
| Skierowanie wprost | Artefakt zostaje przekazany do wskazanego modułu jednym krokiem, z pominięciem przeglądu | 3 | Pozycja „Przekaż do modułu” w menu kontekstowym artefaktu |
| Zwolnienie | Pozycja znika z magistrali; artefakt pozostaje w module pochodzenia i w module docelowym | 3 | Pozycja „Usuń z magistrali” w menu kontekstowym pozycji |

Magistrala przechowuje referencje do artefaktów, nie ich kopie. Przekazanie nie tworzy nowego artefaktu ani nowej wersji — moduł docelowy pracuje na tym samym wpisie encji `artefakt` ([Model danych](../architektura/model-danych.md), rozdz. 16.1), a każde utrwalenie zmiany powstaje jako kolejna `wersja_artefaktu` (rozdz. 16.2) tego samego artefaktu.

**Waga wizualna i miejsce w układzie.** Kolumna boczna, otwierana jako rozszerzenie boczne — kolejna kolumna po prawej stronie obszaru roboczego, o pełnej wysokości obszaru roboczego. Po zamknięciu kolumna znika całkowicie z przestrzeni roboczej. W stanie spoczynku mechanizm reprezentuje wyłącznie licznik pozycji magistrali w pasku kontekstu (warstwa 1).

**Warstwa widoczności i sposób wywołania.**

| Element | Warstwa | Sposób wywołania |
|---|---|---|
| Licznik pozycji magistrali | 1 | Widoczny bez interakcji w pasku kontekstu |
| Kolumna magistrali kontekstu | 3 | Kliknięcie licznika, pozycja menu hamburger (☰) środowiska albo skrót klawiszowy |
| Odłożenie i pobranie artefaktu | 3 | Menu kontekstowe artefaktu, przeciągnięcie |
| Reguły kierowania artefaktów | 4 | Okno konfiguracji, polecenie języka naturalnego w Chat Window, paleta poleceń `Ctrl/Cmd + K` |

**Zasięg.** Zasięg magistrali jest wartością ustawienia konfiguracyjnego i przyjmuje jeden z poziomów zasięgu opisanych w opracowaniu [Izolacja i konfigurowalność zależności](../architektura/izolacja-i-zaleznosci.md) (rozdz. 6): sesja, projekt albo środowisko. Przy zasięgu „sesja” magistrala obejmuje moduły jednej karty sesji; przy zasięgu „projekt” — wszystkie karty grupy projektu; przy zasięgu „środowisko” — wszystkie karty środowiska. Zasięg globalny nie występuje: artefakt nie przechodzi między środowiskami magistralą, lecz przez moduł Library.

**Relacja do modułu Library.** Magistrala jest przepływem roboczym o ograniczonej trwałości, moduł Library ([Moduł Library](../moduly/library.md)) jest repozytorium trwałym. Artefakt odłożony do magistrali pozostaje w niej przez okres wskazany ustawieniem retencji; przeniesienie do Library jest jawnym działaniem Użytkownika i jedyną drogą trwałego zachowania artefaktu oraz jedyną drogą jego przeniesienia między środowiskami.

**Relacja do punktów izolacji.** Magistrala respektuje reguły izolacji kontekstu obowiązujące parę modułów i zasięg bieżącej karty ([Izolacja i konfigurowalność zależności](../architektura/izolacja-i-zaleznosci.md), rozdz. 4, 6, 9). Reguła izolacji kontekstu ustawiona dla pary modułów rozstrzyga, czy artefakt odłożony w module źródłowym jest widoczny w module docelowym; przekazanie niezgodne z regułą nie następuje, a magistrala prezentuje przy pozycji opis obowiązującej reguły wraz z odnośnikiem do okna konfiguracji punktów izolacji (rozdz. 10 tamże).

**Relacja do encji artefaktu w modelu danych.** Pozycja magistrali odwołuje się do wiersza `artefakt` ([Model danych](../architektura/model-danych.md), rozdz. 16.1) przez jego identyfikator, zachowując moduł pochodzenia i kartę sesji odłożenia. Kolekcje i etykiety artefaktu (rozdz. 16.3–16.5) pozostają niezmienione — magistrala ich nie tworzy i nie modyfikuje.

---

## 7. Stany okien i elementów interfejsu

Poniższa tabela ustala wspólny, czterostanowy model interakcji (Domyślny, Wskazanie kursorem, Active, Fokus) obowiązujący dla każdego interaktywnego elementu wymienionego w rozdziałach 3 i 5 — zgodnie z regułami stanów i fokusu systemu wizualnego marki — uzupełniony dwoma stanami procesu (Ładowanie, Błąd), z których żaden nie ogranicza klikalności kontrolki. Żaden interaktywny element nie przyjmuje stanu wyłączonego (disabled) jako blokady wykonania. Kolumna „Zastosowanie w przepływie” wskazuje, gdzie stan wybija się na pierwszy plan tego dokumentu.

| Stan | Reguła wizualna | Zastosowanie w przepływie |
|---|---|---|
| Domyślny (spoczynek) | Powierzchnia neutralna, bez akcentu sygnałowego | Karty środowisk i kafle komponentów na stronie głównej przed interakcją (rozdz. 3.4–3.5, 5.5) |
| Wskazanie kursorem | Subtelna zmiana tła albo uniesienie z akcentem sygnałowym dla elementów akcentowanych | Karta środowiska, kafel komponentu, pozycja bocznej nawigacji, karta sesji (rozdz. 3.5–3.7) |
| Active / wybrany | Wciśnięcie (`translateY(1px)`) chwilowo; stan trwały — akcent sygnałowy (podkreślenie, obrys, tło) | Pozycja modułu podświetlona w bocznej nawigacji (rozdz. 3.7); karta sesji aktywna w pasku kart (rozdz. 3.6) |
| Fokus (nawigacja klawiaturą) | Pierścień sygnałowy 2 px + odsunięcie 2 px, wyłącznie `:focus-visible`, obowiązkowy dla każdej kontrolki interaktywnej | Wszystkie elementy formularza logowania (rozdz. 3.2); wszystkie przełączniki macierzy izolacji (rozdz. 5.6) |
| Ładowanie | Wskaźnik `.dn-spinner`, obrót ciągły, respektuje ograniczenie ruchu systemowego | Okno startowe (rozdz. 3.1); przycisk „Zaloguj się” w trakcie weryfikacji (rozdz. 3.2); przeładowanie modułu (rozdz. 3.7) |
| Błąd | Kreska/obrys w kolorze statusu błędu, komunikat tekstowy towarzyszący polu lub widokowi | Pole formularza logowania po niepowodzeniu weryfikacji (rozdz. 3.2); okno startowe przy braku połączenia z serwerem (rozdz. 3.1) |

**Zasada „zero blokad” w tym dokumencie.** Żaden interaktywny element nawigacyjny opisany w tym dokumencie nie przyjmuje stanu wyłączonego (disabled) jako blokady wykonania — kontrolki pozostają zawsze klikalne. Chwilowa niekompletność danych wejściowych (np. pusty formularz) albo niezakończony poprzedni etap procesu (np. panel wdrożenia oczekujący na zakończenie budowy) są sygnalizowane **po kliknięciu** — komunikatem lub ostrzeżeniem przy właściwym polu bądź elemencie — albo opisowo obok kontrolki, nigdy przez odebranie interakcji. Wyjątkiem sankcjonowanym pozostaje wyłącznie uwierzytelnianie (rozdz. 3.2): próg dostępu do platformy jest jedyną granicą bezpieczeństwa w tym przepływie, lecz jego egzekwowanie jest ustawieniem Operatora w oknie konfiguracji („Wymóg logowania”), nie sztywną blokadą poza jego zasięgiem.

**Reguła obecności i sygnalizacji niegotowości — rozróżnienie ważne dla dewelopera:**

| Sytuacja | Reguła | Przykład |
|---|---|---|
| Element zależy od danych wprowadzonych przez Operatora w tym samym widoku | Element **zawsze klikalny**; brak lub błąd danych sygnalizowany komunikatem po kliknięciu, przy właściwym polu — akcja nie jest blokowana | Przycisk „Zaloguj się” przy niewypełnionym formularzu (rozdz. 3.2) |
| Element zależy od zakończenia poprzedniego etapu tego samego procesu | Element **zawsze klikalny**; kliknięcie przed zakończeniem poprzedniego etapu wyświetla komunikat o brakującym warunku, wskazanym opisowo obok elementu | Deployment Panel modułu Apps przed zakończeniem Frontend/Backend Workspace (rozdz. 5.2) |
| Element zależy od świadomej konfiguracji Operatora (izolacja, powiązanie międzymodułowe) | Element **obecny i w pełni dostępny do skonfigurowania** — brak konfiguracji nie ukrywa ani nie ogranicza elementu, oznacza wyłącznie wartość domyślną | Każdy przełącznik macierzy izolacji (rozdz. 5.6); każde powiązanie międzymodułowe (rozdz. 5.1–5.3) |
| Element dotyczy progu uwierzytelnienia (wyjątek sankcjonowany jako jedyna granica bezpieczeństwa) | Przycisk „Pomiń” **zawsze klikalny** w obu fazach; skuteczność zależy od ustawienia „Wymóg logowania” Operatora w oknie konfiguracji — gdy aktywne, kliknięcie wyświetla komunikat zamiast pominięcia | Przycisk „Pomiń” w oknie rejestracji i logowania (rozdz. 3.2) |

---

## 8. Tabele zbiorcze elementów interfejsu nawigacyjnego

Rozdziały 3 i 5 opisują każdy element w miejscu jego pierwszego wystąpienia. Poniższe tabele porządkują je dodatkowo w przekroju poziomym — wszystkie elementy jednego typu w jednym miejscu — na potrzeby szybkiego odniesienia projektowego.

### 8.1 Elementy poziomu „okno” (jednostki nawigacji pełnoekranowej)

| Okno / widok | Forma i waga | Warstwa | Sposób wywołania | Zamknięcie / opuszczenie | Trwałość stanu |
|---|---|---|---|---|---|
| Okno startowe | Pełnoekranowe, jednolite, bez nawigacji | 1 | Uruchomienie aplikacji | Automatyczne po uzgodnieniu połączenia | Brak — stan przejściowy |
| Przygotowanie środowiska pracy | Pełnoekranowe, jednolite, bez nawigacji | 1 | Zakończone uwierzytelnienie | Automatyczne po skompletowaniu stanu | Brak — stan przejściowy |
| Okno rejestracji i logowania | Pełnoekranowy formularz | 1 | Brak ważnego tokenu lub faza budowy | Zatwierdzenie formularza albo „Pomiń” (skuteczny, gdy „Wymóg logowania” nieaktywne) | Wartości pól zachowane do zatwierdzenia; nie trwałe między uruchomieniami |
| Strona główna | Pełnoekranowa, trzy strefy | 1 | Token wydany; powrót z dowolnego miejsca platformy | Wybór strefy 1/2/3 | Nie przechowuje własnego stanu — jest zawsze tym samym widokiem |
| Przedsionek środowiska | Pełnoekranowy, dwa rejony (szyna sesji · płótno z kaflami) | 1 | Wybór karty środowiska (strefa 1) | Powrót do strony głównej albo wybór modułu | Sesje trwają w tle po opuszczeniu (rozdz. 6) |
| Przestrzeń robocza środowiska | Pełnoekranowa, układ kolumnowy (nawigacja · Chat Window · kolumna dominująca) | 1 | Wybór modułu w przedsionku | Powrót do przedsionka albo do strony głównej | Karty sesji trwają w tle po opuszczeniu (rozdz. 6) |
| Chat Window | Lewa kolumna, stała, pełna wysokość obszaru roboczego | 1 | Obecne bez wywołania w każdym module i każdym środowisku | Nie podlega zamknięciu — okno stałe układu | Historia rozmowy trwała po stronie serwera niezależnie od stanu klienta (rozdz. 6) |
| Execution Loop Window | Kolumna sąsiadująca z Chat Window, pełna wysokość obszaru roboczego | 1 | Przyjęcie zlecenia do wykonania; przywrócenie jednym kliknięciem znacznika przebiegu albo poleceniem języka naturalnego w Chat Window | Zamknięcie kolumny — znika całkowicie z przestrzeni roboczej, pętla biegnie dalej po stronie serwera | Stan pętli (zlecenie, kolejka zadań, wyniki kontroli jakości) trwały po stronie serwera |
| Okno pomocnicze modułu | Kolumna boczna, otwierana jako rozszerzenie boczne | 2 | Kliknięcie pozycji zestawu okien modułu, skrót klawiszowy albo polecenie języka naturalnego w Chat Window | Zamknięcie — kolumna znika całkowicie, zwolniona szerokość przypada kolumnom pozostałym | Stan okna trwały po stronie serwera w obrębie karty sesji |
| Okno operacyjne | Kolumny obszaru roboczego karty sesji: Chat Window (lewa, stała), Execution Loop Window (kolumna sąsiadująca), okna modułu (kolumna dominująca i kolumny boczne) | 1 | Wybór modułu w bocznej nawigacji / panelu orkiestracji | Zmiana modułu, zamknięcie karty, powrót do strony głównej | Trwały po stronie serwera niezależnie od stanu klienta (rozdz. 6) |
| Okno budowy komponentu własnego | Kolumna dominująca obszaru roboczego albo nakładka modalna | 2 | Wybór kafla strefy 2, skrót klawiszowy albo polecenie języka naturalnego w Chat Window | Zapis komponentu albo zamknięcie bez zapisu | Wersja robocza — zależnie od implementacji; zapis trwały po zatwierdzeniu |
| Historia sesji | **Nakładka modalna** nad bieżącym widokiem — sekcja okna ustawień i konfiguracji; trzy zbiory sesji w jednym rejestrze | 2 | Pozycja „Historia sesji” strefy pracy szyny nawigacji, przycisk „Pełna historia sesji” w przedsionku albo sekcja „Historia sesji” w oknie ustawień | Zamknięcie nakładki (widok pod spodem pozostaje nietknięty) albo wejście w wybraną sesję | Rejestr jest widokiem stanu serwera — nie przechowuje stanu własnego |
| Zakładanie projektu (Workspace środowiska) | Pełnoekranowe okno Workspace wskazanego środowiska; pięć grup ustaleń oraz kolumna podsumowania | 1 | Pozycja „Nowy projekt” strefy pracy szyny nawigacji, po wskazaniu środowiska | Założenie projektu (z sesją albo bez), zapis szablonu, opuszczenie bez zapisu | Wersja robocza formularza trwała w obrębie sesji klienta; projekt trwały po założeniu |
| Okno konfiguracji | Nakładka modalna nad bieżącym widokiem | 2 | Pozycja „Okno konfiguracji” (strefa 3 lub rozwinięcie kontekstowe `⋮`), skrót klawiszowy albo wyszukiwarka funkcji | Zamknięcie nakładki | Zmiany zapisywane natychmiast na warstwie wskazanej (globalna / sesji) |

### 8.2 Elementy poziomu „nawigacja” (stałe elementy ramy interfejsu)

| Element | Obecność | Forma i waga | Warstwa | Sposób wywołania | Zawartość |
|---|---|---|---|---|---|
| Chat Window | Każdy moduł i każde środowisko | Lewa kolumna, stała, pełna wysokość obszaru roboczego | 1 | Obecne bez wywołania | Historia rozmowy Użytkownik ↔ Wykonawca, pole poleceń języka naturalnego, strumień odpowiedzi i wyników, zatwierdzanie i przerywanie działań |
| Execution Loop Window | Każdy moduł i każde środowisko realizujące zlecenie | Kolumna sąsiadująca z Chat Window, pełna wysokość obszaru roboczego | 1 | Przyjęcie zlecenia do wykonania | Zlecenie i jego dekompozycja na zadania, kolejka i stan zadań, komunikaty sterujące Koordynator ↔ Wykonawca, wyniki kontroli jakości, wskaźniki przebiegu pętli, sterowanie przebiegiem |
| Belka tytułowa | Strona główna i każde środowisko | Pas ramy na prawo od szyny nawigacji | 1 | Widoczny bez interakcji | Godło (powrót do strony głównej), tytuł widoku, sterowanie oknem |
| Boczna nawigacja modułów | TalkIn, WorkSpace, CodeStudio | Lewa kolumna nawigacyjna, stała | 1 | Widoczna bez interakcji | Lista modułów wg macierzy dostępności (8–9 pozycji) |
| Panel orkiestracji | Wyłącznie MultitaskingAI | Lewa kolumna nawigacyjna, stała, ta sama pozycja co boczna nawigacja | 1 | Widoczny bez interakcji | Sześć sekcji sterowania zespołem |
| Pas kart sesji | Każde środowisko | Pas kart w ramie przestrzeni roboczej; nieobecny w przedsionku | 1 | Widoczny bez interakcji | Karty otwartych przestrzeni roboczych + kontrolka „+” |
| Listwa ustawień | Strona główna (jawnie) + dostępna z każdego miejsca | Zwarta listwa, waga najniższa | 1 | Widoczna bez interakcji; równoważnie skrót klawiszowy i wyszukiwarka funkcji | Okno konfiguracji, Mobile, Always On Display |
| Szyna nawigacji — strefa pracy | Każde okno platformy | Trzy pozycje przy górnej krawędzi szyny, nad środowiskami | 1 | Widoczna bez interakcji | Nowa sesja (menu wyboru środowiska), Historia sesji (okno rejestru konta), Nowy projekt (menu wyboru środowiska → Workspace środowiska) |

### 8.3 Elementy poziomu „kontrolka” (pojedyncze akcje)

| Kontrolka | Forma i waga | Warstwa | Sposób wywołania | Stany istotne dla przepływu | Skutek |
|---|---|---|---|---|---|
| Karta środowiska | Duży panel interaktywny | 1 | Widoczna bez interakcji | spoczynek / wskazanie kursorem / aktywny+fokus | Otwiera przedsionek środowiska (jedyna droga wejścia do środowiska) |
| Kafel komponentu własnego | Średni panel interaktywny | 1 | Widoczny bez interakcji | spoczynek / wskazanie kursorem / aktywny+fokus | Otwiera okno budowy komponentu (nie zmienia środowiska) |
| Pozycja listwy ustawień | Mały element paska | 1 | Widoczna bez interakcji | domyślny / wskazanie kursorem / fokus | Otwiera nakładkę konfiguracji albo aktywuje funkcję globalną |
| Pozycja bocznej nawigacji / sekcji orkiestracji | Mały element listy | 1 | Widoczna bez interakcji | domyślny / wskazanie kursorem / wybrana | Przeładowuje kolumnę dominującą bieżącej karty |
| Karta sesji | Element paska kart | 1 | Widoczna bez interakcji | aktywna / w tle / ze wskaźnikiem pracy | Przełącza widoczną zawartość kolumn obszaru roboczego |
| Pole poleceń Chat Window | Pole tekstowe wieloliniowe w lewej kolumnie | 1 | Widoczne bez interakcji | domyślny / fokus / wypełnione / ładowanie | Kieruje polecenie języka naturalnego do Wykonawcy; jest zarazem drogą wywołania funkcji warstwy 4 |
| Sterowanie przebiegiem pętli `Przebieg ▼` | Element zbiorczy w nagłówku Execution Loop Window | 3 | Rozwinięcie elementu zbiorczego albo polecenie języka naturalnego w Chat Window | zwinięty / rozwinięty | Wstrzymuje, wznawia, przerywa albo koryguje zlecenie pętli wykonawczej |
| Znacznik kontekstowy | Lekki znacznik w kolumnie nawigacji (`[Ubuntu]`, `[Fable 5]`) | 1 (znacznik) / 2 (selektor) | Kliknięcie znacznika otwiera selektor, który zwija się po wyborze | domyślny / wskazanie kursorem / selektor otwarty | Zmienia środowisko, repozytorium, projekt, model albo wykonawcę bieżącej karty |
| Kontrolka „+” (nowa karta) | Mała ikonka | 1 | Widoczna bez interakcji | domyślny / wskazanie kursorem / aktywny | Tworzy nową, pustą kartę sesji |
| Kontrolka „✕” (zamknij kartę) | Mała ikonka | 2 | Ujawniana wskazaniem kursorem na karcie sesji | ukryta do chwili wskazania kursorem / wskazanie kursorem | Zamyka widok karty |
| Przycisk „Pomiń” | Mały przycisk bez obrysu | 1 | Widoczny bez interakcji | obecny i klikalny w obu fazach wdrożenia | Wejście bez uwierzytelnienia, gdy „Wymóg logowania” nieaktywne; w przeciwnym razie komunikat o obowiązującym wymogu logowania |
| Kontrolka „⋮” (rozwinięcie kontekstowe) | Mała ikonka w nagłówku okna | 3 | Kliknięcie `⋮` w nagłówku okna | domyślny / wskazanie kursorem / otwarte | Otwiera podzbiór ustawień warstwy sesji; po zamknięciu rozwinięcie znika całkowicie |
| Przełącznik macierzy izolacji | Suwak | 2 | Widoczny po otwarciu okna konfiguracji punktów izolacji | wyłączony/odrębny (domyślnie) / włączony/współdzielony | Zmienia zakres współdzielenia lub izolacji technicznej |
| Wywołanie funkcji eksperckiej | Bez reprezentacji w stanie spoczynku interfejsu | 4 | Polecenie języka naturalnego w Chat Window, skrót klawiszowy, wyszukiwarka funkcji, tryb administracyjny albo konfiguracja roli | wywołane / zamknięte | Uruchamia administrację procesami sesji, podgląd kanału WebSocket, dzienniki wykonania i narzędzia diagnostyczne niskiego poziomu |

---

## 9. Zgodność z zasadami nadrzędnymi platformy

Przepływ opisany w niniejszym dokumencie pozostaje w pełnej zgodności z zasadami nadrzędnymi platformy. Poniższa tabela wskazuje, gdzie każda zasada odciska się na konkretnym elemencie przepływu.

| Zasada nadrzędna | Zastosowanie w przepływie okien | Rozdział |
|---|---|---|
| Domyślne zachowanie systemu = wykonanie; brak twardych blokad | Żaden interaktywny element nawigacyjny nie przyjmuje stanu wyłączonego jako blokady wykonania — kontrolki pozostają klikalne, niekompletność danych albo niezakończony poprzedni etap procesu sygnalizowane komunikatem po kliknięciu; izolacja i wymóg logowania są zachowaniami skonfigurowanymi przez Operatora dla własnego, jedynego konta, nie barierą narzuconą z zewnątrz | rozdz. 7 |
| Faza budowy nie wymaga dodatkowej czynności od dewelopera/Operatora | Przycisk „Pomiń” pozwala wejść do platformy bez zakładania konta, gdy Operator nie aktywował ustawienia „Wymóg logowania” — tym samym, wspólnym kodem uwierzytelniania co faza produkcyjna | rozdz. 3.2 |
| Pełna konfigurowalność (zasada centralna) | Zakres współdzielenia historii, pamięci i kontekstu między kartami, moduły powiązane jawnie, kolejność sekcji panelu orkiestracji — wszystko konfigurowalne z jednego okna konfiguracji, dostępnego z każdego miejsca przepływu | rozdz. 5.6, 7 |
| Jawność i konfigurowalność zależności | Każde przełączenie środowiska lub modułu jest świadomym, widocznym przejściem (etapy 5, 7), nie ukrytą zmianą stanu; powiązania międzymodułowe pokazane wprost w ścieżkach 5.1–5.3 są możliwością, nie regułą wbudowaną | rozdz. 3.5, 3.7, 5.1–5.3 |
| Stopniowe ujawnianie funkcjonalności | Każdy element przepływu należy do jednej z czterech warstw widoczności; makiety rysowane są w stanie spoczynku, okna i panele warstw 2–3 znikają całkowicie z przestrzeni roboczej po zamknięciu, funkcje warstwy 4 uruchamiane są poleceniem języka naturalnego, skrótem klawiszowym, wyszukiwarką funkcji albo w trybie administracyjnym; każda ukryta funkcja pozostaje osiągalna jednym kliknięciem | rozdz. 3a |
| Dwa kanały komunikacji operacyjnej | Chat Window (Użytkownik ↔ Wykonawca) jest oknem stałym lewej kolumny w każdym module i każdym środowisku; Execution Loop Window (Koordynator ↔ Wykonawca) otwiera się w kolumnie sąsiadującej i prowadzi pętlę wykonawczą — oba są elementami pierwszoplanowymi każdego etapu przepływu | rozdz. 3.8, 3a, 6 |
| Układ pionowy (podział lewa–prawa) | Wszystkie okna robocze, okna komunikacji i okna pomocnicze rozmieszczone są w kolumnach sąsiadujących poziomo; okna pomocnicze otwierają się jako rozszerzenia boczne po prawej stronie obszaru roboczego, a regulacji podlega wyłącznie szerokość kolumn | rozdz. 3.6, 3.8, 5.1–5.4 |
| Pełna kompozycyjność | Liczba jednocześnie otwartych kart sesji i kombinacje modułów w poszczególnych kartach nie są ograniczone przez platformę (rozdz. 5.1: dwie karty TalkIn równolegle) | rozdz. 5.1, 6 |
| Orkiestracja pracy zespołu | Ścieżka MultitaskingAI (rozdz. 5.4) realizuje tę zasadę w warstwie nawigacyjnej — cztery okna robocze ról w kolumnach sąsiadujących, widoczne jednocześnie, wraz z pętlą pracy ciągłej 24/7/365 prowadzoną przez Koordynatora | rozdz. 5.4 |

---

## 10. Słowniczek pojęć

| Pojęcie | Definicja |
|---|---|
| Okno startowe | Pierwsze okno natywne prezentowane po uruchomieniu aplikacji klienckiej, obejmujące nawiązanie kanału WebSocket i wstępne uzgodnienie protokołu z serwerem, poprzedzające decyzję o wyświetleniu okna logowania albo bezpośrednim wejściu do platformy |
| Przygotowanie środowiska pracy | Widok przejściowy po uwierzytelnieniu; odtwarza profil, uprawnienia, karty sesji, kanały modeli i magistralę kontekstu, zanim otworzy się strona główna |
| Okno rejestracji i logowania | Okno uwierzytelniania Operatora; przycisk Pomiń obecny i klikalny w obu fazach wdrożenia — czy okno wymusza skuteczną rejestrację lub logowanie, zależy od ustawienia „Wymóg logowania” w oknie konfiguracji, konfigurowalnego przez Operatora niezależnie od fazy |
| Faza budowy | Faza wdrożenia, w której rdzeń działa lokalnie na komputerze Operatora; ustawienie „Wymóg logowania” domyślnie nieaktywne, zmienialne przez Operatora w oknie konfiguracji |
| Faza produkcyjna | Faza wdrożenia po migracji rdzenia na serwer wirtualny i otwarciu portu kanału WebSocket; ustawienie „Wymóg logowania” domyślnie aktywne, zmienialne przez Operatora w oknie konfiguracji |
| Token dostępu | Poświadczenie wydawane urządzeniu po uwierzytelnieniu, wykorzystywane do nawiązania lub potwierdzenia kolejnych połączeń WebSocket bez ponownego logowania |
| Strona główna (Centrum dowodzenia) | Przedpokój przed strefą roboczą platformy; trzy strefy — wybór środowiska, komponenty własne, ustawienia |
| Strefa 1 | Strefa strony głównej z czterema kartami wejścia do środowisk; jedyna droga wejścia do przestrzeni roboczej środowiska |
| Strefa 2 | Strefa strony głównej z czterema kaflami budowy komponentów własnych: Automations, Agents, Workspace, Assistant |
| Strefa 3 | Strefa strony głównej — listwa ustawień: okno konfiguracji, Mobile, Always On Display |
| Środowisko | Najwyższy poziom organizacji pracy: TalkIn, WorkSpace, CodeStudio, MultitaskingAI |
| Przedsionek środowiska | Widok wejściowy środowiska: szyna sesji przy lewej krawędzi oraz płótno z nagłówkiem środowiska, kaflami modułów i listwą działań środowiska. Nie zawiera Chat Window, bocznej nawigacji ani pasa kart sesji |
| Przestrzeń robocza środowiska | Widok otwierany po wyborze modułu, w układzie pionowym (podział lewa–prawa): boczna nawigacja modułów lub panel orkiestracji, Chat Window, kolumna dominująca obszaru roboczego; szyna nawigacji, belka tytułowa, wstążka z pasmem kart sesji i pasek stanu stanowią ramę okna |
| Boczna nawigacja modułów | Stała, lewa kolumna nawigacyjna środowisk TalkIn, WorkSpace, CodeStudio, wyświetlająca listę dostępnych modułów |
| Panel orkiestracji | Odpowiednik bocznej nawigacji w środowisku MultitaskingAI, zajmujący tę samą lewą kolumnę nawigacyjną — sześć sekcji: Zespoły, Role, Kolejki, Orkiestracja, Harmonogram i automatyki, Monitor procesu |
| Moduł | Wyspecjalizowany obszar roboczy wewnątrz środowiska; piętnaście modułów, dostępność wg macierzy |
| Karta sesji | Pozioma karta w pasie kart przestrzeni roboczej; jedna, samodzielna przestrzeń robocza z własnym układem, historią i kontekstem |
| Okno operacyjne | Pojedynczy element zestawu okien modułu; Chat Window jest oknem wspólnym wszystkich piętnastu modułów |
| Magistrala kontekstu | Mechanizm przekazywania artefaktu między modułami — kolumna boczna z referencjami do odłożonych artefaktów, o zasięgu sesji, projektu albo środowiska (rozdz. 6a) |
| Chat Window | Główne okno komunikacji Użytkownik ↔ Wykonawca; okno stałe, obecne w każdym module i każdym środowisku, zajmujące lewą kolumnę obszaru roboczego na pełnej wysokości. Przyjmuje polecenia w języku naturalnym, prezentuje strumień odpowiedzi i wyników, umożliwia zatwierdzanie oraz przerywanie działań i wyjaśnianie wyniku; jest zarazem drogą wywołania funkcji warstwy 4 |
| Execution Loop Window | Okno pętli wykonawczej Koordynator ↔ Wykonawca, otwierane w kolumnie sąsiadującej z Chat Window; zawiera zlecenie i jego dekompozycję na zadania, kolejkę i stan zadań, wymianę komunikatów sterujących, wyniki kontroli jakości i decyzje o ponowieniu, wskaźniki przebiegu pętli oraz sterowanie przebiegiem (wstrzymanie, wznowienie, przerwanie, korekta zlecenia) |
| Użytkownik | Rola zlecająca i zatwierdzająca pracę platformy |
| Koordynator | Komponent orkiestrujący platformy: dekomponuje zlecenie, przydziela i nadzoruje zadania |
| Wykonawca | AI, agent lub system wykonawczy realizujący zadania |
| Warstwa widoczności | Jedna z czterech warstw ujawniania funkcjonalności (zawsze widoczna, widoczna na żądanie, rozwinięcie kontekstowe, funkcja ekspercka) przypisana każdemu elementowi interfejsu (rozdz. 3a) |
| Rozszerzenie boczne | Kolumna boczna otwierana po prawej stronie obszaru roboczego, mieszcząca okno pomocnicze; po zamknięciu znika całkowicie z przestrzeni roboczej |
| Przeładowanie przestrzeni roboczej | Wymiana zestawu okien operacyjnych, kontekstu czatu i narzędzi w obrębie jednej karty sesji przy zmianie modułu |
| Przełączanie środowiska | Zmiana całej przestrzeni roboczej przez powrót do strony głównej i wybór innej karty środowiska |
| Komponent własny | Nazwany wytwór Operatora — automatyka, agent, projekt lub profil asystenta — tworzony w oknie budowy właściwym strefie 2 i wybieralny operacyjnie po zapisie |
| Uproszczone menu kontekstowe | Rozwinięcie kontekstowe warstwy 3, wywoływane kontrolką `⋮` w nagłówku okna operacyjnego, udostępniające podzbiór ustawień ograniczony do warstwy sesji bieżącej; po zamknięciu znika całkowicie z przestrzeni roboczej |
| Okno konfiguracji | Nakładka modalna z pełnym zakresem trzynastu obszarów ustawień platformy, dostępna z listwy strony głównej i z każdego okna operacyjnego |
| Okno konfiguracji punktów izolacji | Trzypanelowe okno (selektor zasięgu, macierz izolacji, profil i podgląd) łączące izolację kontekstu i izolację techniczną |

---

## Załącznik A. Scenariusz pełny — przejście przez wszystkie warstwy

Scenariusz ilustruje pełny przebieg opisany w rozdziałach 2–6 na przykładzie dewelopera pracującego w fazie budowy, który zakłada agenta, prowadzi równolegle pracę redakcyjną i badawczą, po czym sprawdza kod aplikacji rozwijanej równolegle.

| Krok | Akcja | Co się dzieje | Etap / rozdział |
|---|---|---|---|
| 1 | Uruchamia aplikację kliencką na komputerze | Natywne okno się otwiera; kanał WebSocket nawiązany; `connection.hello` zwraca fazę „budowa” | Okno startowe — rozdz. 3.1 |
| 2 | Widzi okno logowania, klika „Pomiń” | Wejście do platformy bez zakładania konta; krok „Uwierzytelnienie” pominięty | Okno logowania — rozdz. 3.2 |
| 3 | Trafia na stronę główną | Trzy strefy widoczne; brak zapamiętanej aktywnej karty (pierwsze uruchomienie) | Strona główna — rozdz. 3.4 |
| 4 | W strefie 2 klika kafel „Agents” | Otwiera się Agent Builder; konfiguruje model bazowy, tożsamość, skille, uprawnienia; zapisuje agenta | Strefa 2 — rozdz. 5.5 |
| 5 | Wraca na stronę główną, klika kartę „TalkIn” | Otwiera się przedsionek środowiska TalkIn — kafle dziewięciu modułów, szyna sesji | Wybór środowiska — rozdz. 3.5 |
| 6 | W bocznej nawigacji wybiera „Studio” | Karta ładuje zestaw: Studio Editor, Tools Panel, Diff/Grep Panel, Session Repository, Preview Window, Chat Window | Okno operacyjne — rozdz. 3.8, 5.1 |
| 7 | Otwiera nową kartę („+”), wybiera „Research” | Druga karta z odrębnym kontekstem; karta Studio niezmieniona w tle | Karty równoległe — rozdz. 5.1, 6 |
| 8 | Wraca do strony głównej, klika kartę „CodeStudio” | Obie karty TalkIn trwają w tle; otwiera się przedsionek środowiska CodeStudio | Przełączanie środowiska — rozdz. 6 |
| 9 | Wybiera moduł „Developer”, przypisuje uprzednio utworzonego agenta jako wykonawcę zadania | Code Editor, Project Tree, Git Panel, Build Output się ładują; agent z kroku 4 dostępny jako komponent własny do wyboru w sesji | Okno operacyjne — rozdz. 5.3, 5.5 |
| 10 | Zamyka aplikację na koniec dnia | Kanał WebSocket zamknięty; wszystkie trzy karty (Studio, Research, Developer) trwają jako procesy serwera | Rozłączenie — rozdz. 6 |
| 11 | Następnego dnia uruchamia aplikację ponownie | Okno startowe; token urządzenia ważny; klient zgłasza powrót do ostatnio aktywnej karty | Ponowne połączenie — rozdz. 3.1, 6 |
| 12 | Trafia bezpośrednio do karty „Developer” w stanie sprzed zamknięcia | Pełny stan odtworzony: układ, historia, kontekst — bez przechodzenia przez stronę główną | Wariant 2b — rozdz. 4, poz. 2b |

---

## Załącznik B. Macierz poprzedzania i następstwa ekranów

Zestawienie porządkowe: dla każdego ekranu/okna głównego przepływu — co go bezpośrednio poprzedza i co bezpośrednio po nim następuje, z pominięciem rozgałęzień pobocznych (ujętych w rozdz. 4).

| Ekran / okno | Bezpośrednio poprzedza go | Bezpośrednio następuje po nim |
|---|---|---|
| Okno startowe | Uruchomienie aplikacji (zdarzenie systemowe) | Okno rejestracji i logowania **albo** Strona główna **albo** Okno operacyjne (rozdz. 4, poz. 2/2a/2b) |
| Przygotowanie środowiska pracy | Okno rejestracji i logowania (token wydany) **albo** Okno startowe (token ważny) | Strona główna |
| Okno rejestracji i logowania | Okno startowe (brak ważnego tokenu lub faza budowy) | Strona główna |
| Strona główna | Okno startowe (token ważny) **albo** Okno rejestracji i logowania **albo** powrót z przedsionka lub okna operacyjnego **albo** zamknięcie okna strefy 2/3 | Przedsionek środowiska (strefa 1) **albo** okno budowy komponentu (strefa 2) **albo** okno konfiguracji / aktywacja funkcji globalnej (strefa 3) |
| Przedsionek środowiska | Strona główna (wybór karty środowiska) | Wybór modułu → przestrzeń robocza z oknem operacyjnym |
| Okno operacyjne | Wybór modułu w oknie środowiska | Zmiana modułu (ta sama karta) **albo** nowa karta **albo** powrót do strony głównej **albo** zamknięcie aplikacji (proces trwa na serwerze) |
| Okno budowy komponentu własnego (strefa 2) | Strona główna (kafel strefy 2) **albo** okno modułowe środowiska (dla Agents i Workspace) | Zapis → powrót do miejsca wejścia; wytwór wybieralny w oknie operacyjnym |
| Okno konfiguracji (strefa 3) | Strona główna **albo** dowolne okno środowiska/operacyjne (uproszczone menu kontekstowe) | Zamknięcie → powrót do dokładnie tego samego widoku sprzed otwarcia |

---

*Koniec dokumentu. Przepływ okien: od uruchomienia aplikacji do okna roboczego — Specyfikacja docelowa, wersja 2.0, 2026-08-20.*

---
*Danaco Console — Platforma AI Workspace OS · v2.0 · status Deweloperski*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz.*
*Warunki korzystania: [Licencja produktu](../LICENSE.md). Kontakt: support@danaco-group.pl*
