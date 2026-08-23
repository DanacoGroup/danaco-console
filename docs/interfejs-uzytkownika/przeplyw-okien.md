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
| **Data** | 2026-08-06 |

**Informacje szczegółowe dokumentu:**

| | |
|---|---|
| **Dokument** | Przepływ okien — specyfikacja nawigacyjna |
| **Odbiorcy** | DESIGNER (co, gdzie, w jakiej formie, do czego) · DEWELOPER (co zbudować) |
| **Opracowanie** | Danaco Console — agent projektowy przepływu okien |
| **Źródło faktów** | Koncepcja platformy (autorytatywna); Strona główna i nawigacja 1.1; Specyfikacja okien operacyjnych 1.1; System wizualny; Model konfiguracji 1.1; Specyfikacja agentów; Specyfikacja MultitaskingAI; Izolacja i zależności; Architektura 2.1; Bezpieczeństwo i uwierzytelnianie; Kontrakty komunikacji; Rozszerzenia |
| **Zakres** | Pełny przepływ okien od uruchomienia aplikacji klienckiej do okna roboczego: okno startowe → okno rejestracji i logowania → strona główna → wybór środowiska → okno środowiska → wybór modułu → okno operacyjne, wraz ze ścieżkami szczegółowymi czterech środowisk oraz wejściami do strefy komponentów własnych i strefy ustawień |
| **Poza zakresem** | Wymiary, tokeny kolorów i typografii, ostateczny układ graficzny (przedmiot systemu wizualnego marki i prac wykonawczych fazy projektowej); szczegółowa zawartość poszczególnych okien operacyjnych modułów poza tym, co potrzebne do zilustrowania przejścia (przedmiot odrębnych opracowań: Specyfikacja modułów, Specyfikacja okien operacyjnych) |

Dokument nie wprowadza nowych funkcji, modułów, okien ani mechanizmów ponad te ustalone w opracowaniach źródłowych. Nazwy własne środowisk, modułów, okien i ról przejęto bez zmian ze źródeł. Tam, gdzie dokument nazywa i opisuje etap techniczny wprost nieujęty pod odrębną nazwą własną w źródłach (okno startowe — faza łączenia z serwerem poprzedzająca okno logowania), fakt ten jest wskazany jawnie w miejscu opisu wraz z podstawą źródłową, zgodnie z zasadą przejrzystości wobec czytelnika.

---

## Spis treści

0. [Jak korzystać z dokumentu](#0-jak-korzystać-z-dokumentu)
1. [Słownik ram przepływu](#1-słownik-ram-przepływu)
2. [Mapa całości — pełny przepływ w jednym diagramie](#2-mapa-całości--pełny-przepływ-w-jednym-diagramie)
3. [Siedem etapów przepływu głównego](#3-siedem-etapów-przepływu-głównego)
   - 3.1 [Etap 1 — Okno startowe](#31-etap-1--okno-startowe)
   - 3.2 [Etap 2 — Okno rejestracji i logowania](#32-etap-2--okno-rejestracji-i-logowania)
   - 3.3 [Etap 3 — Strona główna (Centrum dowodzenia)](#33-etap-3--strona-główna-centrum-dowodzenia)
   - 3.4 [Etap 4 — Wybór środowiska (Strefa 1)](#34-etap-4--wybór-środowiska-strefa-1)
   - 3.5 [Etap 5 — Okno środowiska (przestrzeń robocza)](#35-etap-5--okno-środowiska-przestrzeń-robocza)
   - 3.6 [Etap 6 — Wybór modułu (lub roli / sekcji orkiestracji)](#36-etap-6--wybór-modułu-lub-roli--sekcji-orkiestracji)
   - 3.7 [Etap 7 — Okno operacyjne (okno robocze)](#37-etap-7--okno-operacyjne-okno-robocze)
3a. [Warstwy widoczności w przepływie okien](#3a-warstwy-widoczności-w-przepływie-okien)
4. [Tabela przejść — zestawienie zbiorcze](#4-tabela-przejść--zestawienie-zbiorcze)
5. [Ścieżki szczegółowe](#5-ścieżki-szczegółowe)
   - 5.1 [Ścieżka TalkIn](#51-ścieżka-talkin)
   - 5.2 [Ścieżka WorkSpace](#52-ścieżka-workspace)
   - 5.3 [Ścieżka CodeStudio](#53-ścieżka-codestudio)
   - 5.4 [Ścieżka MultitaskingAI](#54-ścieżka-multitaskingai)
   - 5.5 [Ścieżka Strefy 2 — komponenty własne](#55-ścieżka-strefy-2--komponenty-własne)
   - 5.6 [Ścieżka Strefy 3 — konfiguracja, ustawienia, funkcje globalne](#56-ścieżka-strefy-3--konfiguracja-ustawienia-funkcje-globalne)
6. [Karty sesji i trwałość procesu w tle](#6-karty-sesji-i-trwałość-procesu-w-tle)
6a. [Magistrala kontekstu — przekazywanie artefaktu między modułami](#6a-magistrala-kontekstu--przekazywanie-artefaktu-między-modułami)
7. [Stany okien i elementów interfejsu](#7-stany-okien-i-elementów-interfejsu)
8. [Tabele zbiorcze elementów interfejsu nawigacyjnego](#8-tabele-zbiorcze-elementów-interfejsu-nawigacyjnego)
9. [Zgodność z zasadami nadrzędnymi platformy](#9-zgodność-z-zasadami-nadrzędnymi-platformy)
10. [Słowniczek pojęć](#10-słowniczek-pojęć)
- [Załącznik A. Scenariusz pełny — przejście przez wszystkie warstwy](#załącznik-a-scenariusz-pełny--przejście-przez-wszystkie-warstwy)
- [Załącznik B. Macierz poprzedzania i następstwa ekranów](#załącznik-b-macierz-poprzedzania-i-następstwa-ekranów)

---

## 0. Jak korzystać z dokumentu

Dokument ma dwóch odbiorców czytających go z różnym celem. Poniższa tabela kieruje uwagę każdego z nich do właściwych fragmentów.

| Odbiorca | Pytanie, na które szuka odpowiedzi | Gdzie szukać |
|---|---|---|
| DESIGNER | Co użytkownik widzi na każdym etapie, gdzie leży każdy element, w której kolumnie układu, jaką ma formę, wagę wizualną i warstwę widoczności, jak wygląda w każdym stanie | Makiety tekstowe przy każdym etapie (rozdz. 3), rozdział o warstwach widoczności (rozdz. 3a), tabele elementów interfejsu (rozdz. 8), rozdział o stanach (rozdz. 7) |
| DEWELOPER | Jaki warunek wyzwala przejście, co technicznie się dzieje, jaki jest stan okna przed i po, jak wywoływany jest element danej warstwy widoczności, jak zachowują się procesy sesji w tle | Tabele przejść przy każdym etapie (rozdz. 3) i tabela zbiorcza (rozdz. 4), rozdział o warstwach widoczności (rozdz. 3a), rozdział o trwałości procesów (rozdz. 6) |

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
| Okno rejestracji i logowania | Okno uwierzytelniania Operatora — rejestracja przy pierwszym uruchomieniu, logowanie przy kolejnych; przycisk Pomiń zawsze dostępny i klikalny, wymóg logowania konfigurowalny przez Operatora w oknie konfiguracji | rozdz. 3.2 |
| Strona główna (Centrum dowodzenia) | Przedpokój przed strefą roboczą; trzy strefy: wybór środowiska, komponenty własne, ustawienia | rozdz. 3.3 |
| Środowisko | Najwyższy poziom organizacji pracy: TalkIn, WorkSpace, CodeStudio, MultitaskingAI | rozdz. 3.4 |
| Okno środowiska | Przestrzeń robocza otwierana po wyborze środowiska: pasek marki, kolumna nawigacji (lub panel orkiestracji), Chat Window, pasek kart sesji | rozdz. 3.5 |
| Moduł | Wyspecjalizowany obszar roboczy wewnątrz środowiska; piętnaście modułów, dostępność wg macierzy | rozdz. 3.6 |
| Okno operacyjne | Pojedynczy element zestawu okien modułu — właściwa przestrzeń wykonania zadania; Chat Window jest oknem wspólnym wszystkich modułów | rozdz. 3.7 |
| Chat Window | Główne okno komunikacji Użytkownik ↔ Wykonawca; okno stałe, obecne w każdym module i każdym środowisku, zajmujące lewą kolumnę obszaru roboczego na pełnej wysokości | rozdz. 3.7, 3a |
| Execution Loop Window | Okno pętli wykonawczej Koordynator ↔ Wykonawca; otwierane w kolumnie sąsiadującej z Chat Window, prezentuje dekompozycję zlecenia, kolejkę zadań i sterowanie przebiegiem pętli | rozdz. 3.7, 3a |
| Warstwa widoczności | Przypisanie elementu interfejsu do jednej z czterech warstw ujawniania funkcjonalności; wyznacza, czy element jest widoczny bez interakcji, po użyciu wyzwalacza, w rozwinięciu kontekstowym, czy wyłącznie poleceniem języka naturalnego, skrótem, wyszukiwarką funkcji lub w trybie administracyjnym | rozdz. 3a |
| Karta sesji | Pozioma karta w pasku kart okna środowiska, reprezentująca jedną, samodzielną przestrzeń roboczą | rozdz. 6 |
| Komponent własny | Nazwany wytwór Operatora: automatyka, agent, projekt lub profil asystenta, tworzony ze strefy 2 strony głównej | rozdz. 5.5 |

Cztery pojęcia z powyższej tabeli — Środowisko, Moduł, Okno operacyjne, Karta sesji — powtarzają trójstopniową hierarchię platformy:

```
ŚRODOWISKO  →  MODUŁ  →  OKNO OPERACYJNE
„w jakim trybie      „jakie zadanie        „jakim narzędziem
 pracuję?”            wykonuję?”            realizuję zadanie?”
```

Niniejszy dokument dokłada do tej hierarchii dwa etapy poprzedzające ją (okno startowe, okno logowania) oraz jeden etap pośredni, przez który hierarchia jest zawsze osiągana (strona główna) — pięć elementów łącznie, opisanych w rozdziale 3.

---

## 2. Mapa całości — pełny przepływ w jednym diagramie

Poniższy diagram jest mapą odniesienia dla całego dokumentu — pokazuje wszystkie siedem etapów przepływu głównego wraz z trzema rozgałęzieniami strony głównej (strefa 1, strefa 2, strefa 3) oraz miejscem, w którym każda ścieżka szczegółowa z rozdziału 5 dołącza do całości.

```
╔═══════════════════════════════════════════════════════════════════════════╗
║  0 · URUCHOMIENIE APLIKACJI KLIENCKIEJ                                     ║
║      Operator otwiera natywne okno aplikacji na urządzeniu                 ║
╚════════════════════════════════════════╤════════════════════════════════════╝
                                          ▼
┌───────────────────────────────────────────────────────────────────────────┐
│  1 · OKNO STARTOWE                                                          │
│      nawiązanie kanału WebSocket · uzgodnienie protokołu (connection.hello) │
└───────────────────────────────────────────────────┬─────────────────────────┘
              │                                      │
   token urządzenia NIEWAŻNY                  token urządzenia WAŻNY
   lub brak tokenu                            (urządzenie już uwierzytelnione)
              │                                      │
              ▼                                      ├── zapamiętana aktywna karta ──► [ 7 ] OKNO OPERACYJNE
┌────────────────────────────────────┐               │    (pomija stronę główną — rozdz. 3.1, 6)
│  2 · OKNO REJESTRACJI I LOGOWANIA    │               │
│                                      │               └── brak zapamiętanej karty ───► [ 3 ] STRONA GŁÓWNA
│  FAZA BUDOWY         FAZA PRODUKCYJNA│
│  [Pomiń] klikalny    [Pomiń] klikalny│
│      │                    │          │
│      │              REJESTRACJA  LOGOWANIE
│      │              (konto nie   (konto już
│      │               istnieje)    istnieje)
│      │                    │          │
│      └────────────────────┴──────────┘
│                    │
│           token dostępu wydany urządzeniu
└────────────────────┬────────────────────────────────────────────────────────┘
                      ▼
┌───────────────────────────────────────────────────────────────────────────┐
│  3 · STRONA GŁÓWNA — CENTRUM DOWODZENIA                                     │
│      przedpokój przed strefą roboczą — trzy niezależne strefy               │
│                                                                             │
│   STREFA 1 (waga główna)     STREFA 2 (waga pośrednia)   STREFA 3 (najniższa)│
│   karty środowisk             kafle komponentów własnych   listwa ustawień  │
│   TalkIn·WorkSpace·           Automations·Agents·          Okno konfiguracji│
│   CodeStudio·MultitaskingAI   Workspace·Assistant           Mobile · AOD    │
│         │                            │                            │        │
└─────────┼────────────────────────────┼────────────────────────────┼────────┘
          ▼                            ▼                            ▼
┌──────────────────────┐   ┌─────────────────────────────┐  ┌───────────────────┐
│ 4 · WYBÓR ŚRODOWISKA  │   │  5.5 STREFA 2 — SZCZEGÓŁY    │  │ 5.6 STREFA 3       │
│ kliknięcie karty       │   │  okno budowy komponentu      │  │ okno konfiguracji   │
│ środowiska              │   │  własnego (np. Agent         │  │ (nakładka modal)    │
│                        │   │  Builder) → zapis jako        │  │ albo aktywacja       │
│                        │   │  komponent własny              │  │ Mobile / AOD         │
└───────────┬────────────┘   └───────────────┬────────────────┘  └──────────┬──────────┘
            │                                 ▼                              ▼
            │                 wytwór (agent, automatyka,        aktywacja funkcji
            │                 projekt, profil asystenta)         globalnej — BEZ nowej
            │                 zapisany jako zasób Operatora,     przestrzeni roboczej;
            │                 wybieralny dalej w krokach 6–7 ─┐  bieżący widok trwa
            │                 (wpięty w sesji modułu)          │  niezmieniony pod spodem
            ▼                                                  │
┌───────────────────────────────────────────────────────────┐  │
│  5 · OKNO ŚRODOWISKA                                         │◄─┘
│      pasek marki · kolumna nawigacji modułów                  │
│      (lub panel orkiestracji — MultitaskingAI, rozdz. 5.4)    │
│      · pasek kart sesji                                       │
└─────────────────────────────┬─────────────────────────────────┘
                               ▼
┌───────────────────────────────────────────────────────────────┐
│  6 · WYBÓR MODUŁU (lub roli — rozdz. 5.4)                        │
│      kliknięcie pozycji bocznej nawigacji                        │
│      albo sekcji panelu orkiestracji                             │
└─────────────────────────────┬─────────────────────────────────┘
                               ▼
┌───────────────────────────────────────────────────────────────────────────┐
│  7 · OKNO OPERACYJNE — OKNO ROBOCZE                                          │
│      układ pionowy (podział lewa–prawa):                                     │
│      Chat Window (lewa kolumna, stała) │ Execution Loop Window (kolumna      │
│      sąsiadująca) │ obszar roboczy modułu (kolumna dominująca) │ okna        │
│      pomocnicze jako rozszerzenia boczne po prawej                           │
│      ── tu zaczyna się praca właściwa ──                                    │
│                                                                             │
│      ścieżki szczegółowe:  5.1 TalkIn · 5.2 WorkSpace ·                     │
│                             5.3 CodeStudio · 5.4 MultitaskingAI              │
└───────────────────┬───────────────────────────────────────────────────────┘
                     │  „+” nowa karta sesji (powrót do kroku 6 w nowej karcie)
                     │  powrót do strony głównej (powrót do kroku 3;
                     │     karta bieżąca trwa w tle — rozdz. 6)
                     ▼
              (pętla pracy równoległej — wiele kart, wiele środowisk jednocześnie)
```

Trzy uwagi porządkujące mapę całości, rozwinięte w dalszych rozdziałach:

| Uwaga | Rozwinięcie |
|---|---|
| Etapy 1–2 poprzedzają hierarchię platformy i występują dokładnie raz na sesję połączenia urządzenia (nie przy każdym powrocie do pracy) | rozdz. 3.1–3.2, rozdz. 6 |
| Etap 3 (strona główna) jest punktem, do którego przepływ zawsze może wrócić — nie jest odwiedzany tylko raz | rozdz. 3.3, rozdz. 8.2 dokumentu Strona główna i nawigacja |
| Etapy 4–7 tworzą pętlę: z etapu 7 Operator może otworzyć nową kartę (powrót do etapu 6 w nowej karcie) albo wrócić do strony głównej (powrót do etapu 3), bez zamykania pracy pozostawionej w tle | rozdz. 6 niniejszego dokumentu |

---

## 3. Siedem etapów przepływu głównego

Konwencja opisu każdego etapu jest jednolita: **warunek** (co musi zajść, aby etap się rozpoczął), **co się dzieje** (mechanizm), **stan okna** (jak wygląda i zachowuje się okno w tym etapie), **makieta tekstowa** (rozkład przestrzenny), **tabela elementów interfejsu** (co to jest · do czego służy · forma i waga · stany · zachowanie po interakcji · gdzie występuje).

### 3.1. Etap 1 — Okno startowe

| Aspekt | Treść |
|---|---|
| Warunek wejścia | Operator uruchamia pakiet kliencki na urządzeniu (komputer, telefon, tablet); zdarzenie systemowe, nie decyzja nawigacyjna |
| Co się dzieje | Natywne okno aplikacji (powłoka Tauri) otwiera się i osadza silnik prezentacji; klient otwiera kanał WebSocket pod adresem serwera skonfigurowanym w pakiecie klienckim; klient wysyła komunikat `connection.hello` (wersja protokołu, identyfikator urządzenia, token dostępu — o ile posiadany); serwer odpowiada `connection.hello.ack` z wersją serwera, fazą wdrożenia (budowa / produkcyjna) i wynikiem uwierzytelnienia |
| Stan okna | Przejściowy — trwa od ułamka sekundy do kilku sekund, zależnie od jakości połączenia. Nie jest miejscem żadnej decyzji Operatora; jedyna możliwa interakcja to przerwanie (zamknięcie okna) |
| Wynik | Rozgałęzienie do etapu 2 (okno logowania) albo, gdy urządzenie posiada ważny token, bezpośrednio do etapu 3 lub etapu 7 (rozdz. 6) |

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
| Powrót z ważnym tokenem | Urządzenie było już uwierzytelnione | Jak wyżej, czas trwania zwykle krótszy — bez kroku rejestracji/logowania | Etap 3 albo etap 7 (rozdz. 6) |
| Błąd połączenia | Serwer nieosiągalny, adres nieprawidłowy, przerwanie sieci | Komunikat stanu zastąpiony informacją o braku połączenia oraz kontrolką ponowienia | Ponowienie próby (ten sam etap) — nie jest to blokada dostępu, lecz odzwierciedlenie faktycznego stanu sieci |

**Tabela elementów interfejsu — Okno startowe:**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Gdzie występuje |
|---|---|---|---|---|---|---|---|---|
| Godło marki | Symbol graficzny Danaco Console (`logo-danaco.svg`) wyśrodkowany na powierzchni okna | Identyfikacja marki w czasie oczekiwania | Duży, pojedynczy element graficzny — jedyny akcent wizualny okna | 1 | Widoczne bez interakcji | Statyczny (nie zmienia stanu) | Brak interakcji | Wyłącznie okno startowe |
| Wskaźnik ładowania | Animowany spinner (`.dn-spinner`) z komunikatem tekstowym | Sygnalizacja trwającego połączenia z serwerem | Mały element graficzny + krótki tekst pod godłem | 1 | Widoczny bez interakcji w czasie łączenia | Ładowanie (obrót ciągły, respektuje ograniczenie ruchu) → zastępowany komunikatem błędu przy niepowodzeniu | Brak interakcji w stanie ładowania | Wyłącznie okno startowe |
| Komunikat błędu połączenia | Tekst i ikona ostrzeżenia zastępujące wskaźnik ładowania | Poinformowanie o braku połączenia z serwerem | Mały blok tekstowy z ikoną (`ostrzezenie.svg`) | 1 | Wyświetlany samoczynnie w stanie błędu połączenia | Błąd (widoczny wyłącznie po niepowodzeniu połączenia) | Wyświetlany do chwili automatycznego ponowienia lub ręcznej kontrolki „Spróbuj ponownie” | Wyłącznie okno startowe, stan błędu |
| Kontrolka „Spróbuj ponownie” | Przycisk drugorzędny (`.dn-btn--zarys`) | Ręczne ponowienie próby połączenia | Mały przycisk pod komunikatem błędu | 2 | Widoczna po wystąpieniu błędu połączenia; kliknięcie ponawia etap 1 | Domyślny · hover · aktywny | Ponawia etap 1 od nawiązania kanału WebSocket | Wyłącznie okno startowe, stan błędu |

---

### 3.2. Etap 2 — Okno rejestracji i logowania

| Aspekt | Treść |
|---|---|
| Warunek wejścia | Krok „Uwierzytelnienie” cyklu życia połączenia wskazuje brak ważnego tokenu na urządzeniu (pierwsze połączenie tego urządzenia albo token unieważniony) — **albo** faza budowy, niezależnie od stanu tokenu, gdy okno pozostaje widoczne jako podgląd wyglądu docelowego. Czy okno wymusza skuteczne wypełnienie formularza przed przejściem dalej, ustala pozycja „Wymóg logowania” w oknie konfiguracji (rozdz. 5.6) — ustawienie dostępne Operatorowi w obu fazach wdrożenia, z domyślną wartością nieaktywną w fazie budowy i aktywną w fazie produkcyjnej |
| Co się dzieje | Klient prezentuje okno rejestracji i logowania w pełnym, docelowym wyglądzie niezależnie od fazy wdrożenia; przycisk Pomiń jest obecny i klikalny w obu fazach. Gdy ustawienie „Wymóg logowania” jest nieaktywne, kliknięcie przycisku Pomiń przenosi natychmiast do etapu 3. Gdy ustawienie jest aktywne, kliknięcie Pomiń wyświetla komunikat o obowiązującym wymogu uwierzytelnienia zamiast pominięcia kroku — dalsze przejście wymaga skutecznej rejestracji albo logowania |
| Stan okna | Interaktywny formularz — jedyny etap całego przepływu, w którym Operator wprowadza dane uwierzytelniające. Trzy pod-stany: rejestracja (konto nie istnieje), logowanie (konto istnieje), odzyskiwanie konta (utrata dostępu) |
| Wynik | Wydanie urządzeniu tokenu dostępu i przejście do etapu 3 — albo, gdy ustawienie „Wymóg logowania” jest nieaktywne, przejście do etapu 3 przez przycisk Pomiń, bez wydania tokenu uwierzytelnionego konta |

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
 │              [       Zaloguj się       ]  ← .dn-btn--zloty (CTA)
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
| Zakładki „Zaloguj się / Zarejestruj się” | Przełącznik segmentowy (`.dn-zakladki`) | Wybór między logowaniem a rejestracją | Element średniej wagi, u góry formularza | 1 | Widoczne bez interakcji | Domyślny (Zaloguj się aktywna, gdy konto już istnieje) · wybrana · hover | Przełącza zestaw pól poniżej bez opuszczania okna | Okno rejestracji i logowania |
| Pole „Login” | Pole tekstowe (`.dn-input`) | Wprowadzenie nazwy właściciela konta | Pole średniej szerokości z etykietą | 1 | Widoczne bez interakcji | Domyślny · fokus · wypełnione · błąd (format) | Wartość zapamiętywana do chwili wysłania formularza | Widok logowania i rejestracji |
| Pole „Hasło” | Pole tekstowe maskowane, z przełącznikiem widoczności (ikona `oko.svg`) | Wprowadzenie hasła dostępowego | Pole średniej szerokości z etykietą i ikoną odsłonięcia | 1 | Widoczne bez interakcji | Domyślny · fokus · wypełnione · błąd (niezgodność) | Odsłonięcie/zamaskowanie treści po kliknięciu ikony; wartość przesyłana wyłącznie po zatwierdzeniu formularza | Widok logowania i rejestracji |
| Pole „Adres e-mail uwierzytelniający” | Pole tekstowe (`.dn-input`, walidacja formatu) | Wprowadzenie adresu pełniącego funkcję weryfikacji, logowania i jedynej drogi odzyskania konta | Pole średniej szerokości z etykietą | 1 | Widoczne bez interakcji w widoku rejestracji | Domyślny · fokus · wypełnione · błąd (format) | Po rejestracji inicjuje wysyłkę wiadomości potwierdzającej | Widok rejestracji; widok logowania, gdy metoda e-mail aktywna |
| Wybór metody logowania | Grupa przełączników radiowych (`.dn-check`) | Wskazanie, którą z aktywnych metod Operator się loguje | Mały zestaw pozycji pod polami głównymi | 2 | Zestaw metod aktywnych; wybór zwija się po użyciu | Widoczne wyłącznie metody aktywne (hasło zawsze; PIN i Windows Hello tylko po włączeniu w oknie konfiguracji) | Zmiana metody przełącza układ pól poniżej (np. PIN zamiast hasła) | Widok logowania |
| Przycisk „Zaloguj się” / „Zarejestruj się” | Przycisk główny CTA (`.dn-btn--zloty`) | Zatwierdzenie formularza i wysłanie danych do serwera | Duży, pojedynczy przycisk akcentowany złotem — jedyny CTA widoku | 1 | Widoczny bez interakcji | Domyślny · hover · aktywny (wciśnięcie) · ładowanie (podczas weryfikacji) | Zawsze klikalny; przy niewypełnionych polach wymaganych kliknięcie wyświetla ostrzeżenie przy brakującym polu, bez wysyłki żądania; przy polach wypełnionych wysyła polecenie `auth.login` albo `auth.register` — przy powodzeniu przejście do etapu 3, przy niepowodzeniu komunikat błędu przy właściwym polu, formularz pozostaje wypełniony | Widok logowania i rejestracji |
| Link „Nie pamiętam hasła” | Odnośnik tekstowy | Wejście w funkcję odzyskiwania konta | Mały element tekstowy, waga najniższa formularza | 2 | Kliknięcie odnośnika otwiera widok odzyskiwania konta | Domyślny · hover · fokus | Zamienia formularz na widok odzyskiwania konta (pole adresu e-mail) | Widok logowania |
| Przycisk „Pomiń” | Przycisk drugorzędny bez obrysu (`.dn-btn--duch`) | Wejście do platformy bez wypełniania formularza, gdy wymóg logowania nie jest aktywny | Mały przycisk, wyraźnie oddzielony kreską (`.dn-separator`) od formularza głównego, najniższa waga wizualna okna | 1 | Widoczny bez interakcji | Obecny i klikalny w obu fazach wdrożenia | Gdy ustawienie „Wymóg logowania” (oknie konfiguracji, konfigurowalne przez Operatora niezależnie od fazy) jest nieaktywne — natychmiastowe przejście do etapu 3 bez wydania tokenu, krok „Uwierzytelnienie” pozostaje pominięty; gdy ustawienie jest aktywne — komunikat o obowiązującym wymogu uwierzytelnienia, formularz pozostaje otwarty | Okno logowania, obie fazy |
| Komunikat błędu formularza | Alert liniowy (`.dn-alert--blad`) przy polu, którego weryfikacja się nie powiodła | Wskazanie przyczyny niepowodzenia (dane niepoprawne, adres niepotwierdzony) | Mały blok tekstowy z kreską lewą czerwoną, przy właściwym polu | 1 | Wyświetlany samoczynnie po niepowodzeniu weryfikacji | Widoczny wyłącznie po niepowodzeniu weryfikacji | Znika po poprawnym wypełnieniu pola i ponownym zatwierdzeniu | Widok logowania, rejestracji i odzyskiwania konta |

---

### 3.3. Etap 3 — Strona główna (Centrum dowodzenia)

| Aspekt | Treść |
|---|---|
| Warunek wejścia | Token dostępu wydany (etap 2) i brak zapamiętanej aktywnej karty sesji do przywrócenia — **albo** świadomy powrót Operatora ze środowiska (rozdz. 6) — **albo** zamknięcie okna komponentu własnego strefy 2 lub okna konfiguracji strefy 3 |
| Co się dzieje | Klient wysyła `home.enter` (etap 3 cyklu życia połączenia — „Powiązanie”); serwer wiąże połączenie ze stanem strony głównej. Żadna przestrzeń robocza środowiska nie jest jeszcze otwarta — strona główna nie otwiera jej bezpośrednio, pełni funkcję przedpokoju |
| Stan okna | Stabilny, statyczny układ trzech stref o malejącej wadze wizualnej; brak treści ładowanej na żywo (poza ewentualnym wskaźnikiem aktywności kart sesji trwających w tle) |
| Wynik | Trzy niezależne dalsze ścieżki, wg strefy klikniętej przez Operatora — rozdz. 3.4 (strefa 1), rozdz. 5.5 (strefa 2), rozdz. 5.6 (strefa 3) |

**Makieta tekstowa — Strona główna, trzy strefy:**

```
 ┌───────────────────────────────────────────────────────────────────────┐
 │  ◈ Danaco Console                                    [ 🔍 ]   ● Operator │  ← pasek marki
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

**Tabela elementów interfejsu — Strona główna (poziom strefy, bez rozbicia na pojedyncze karty — te w rozdz. 3.4 i 5.5–5.6):**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Gdzie występuje |
|---|---|---|---|---|---|---|---|---|
| Pasek marki górny | Wąski pasek na całą szerokość okna, tło marki (granat), z godłem, wyszukiwaniem i wskaźnikiem konta | Identyfikacja marki, dostęp do wyszukiwania globalnego, wskazanie zalogowanego Operatora | Pasek pełnej szerokości, waga niska (element ramowy, nie treściowy) | 1 | Widoczny bez interakcji | Statyczny; pole wyszukiwania: domyślny · fokus · z wynikami | Kliknięcie godła wraca zawsze do strony głównej z dowolnego miejsca platformy | Strona główna i okno środowiska (element wspólny) |
| Siatka strefy 1 (cztery karty) | Kontener siatki czterech elementów równej wagi | Grupowanie punktów wejścia do środowisk | Duży blok, środek ciężkości strony | 1 | Widoczna bez interakcji | — (kontener) | — | Wyłącznie strona główna |
| Siatka strefy 2 (cztery kafle) | Kontener siatki czterech elementów lżejszych niż strefa 1 | Grupowanie punktów wejścia do budowy komponentów własnych | Średni blok, pod siatką strefy 1 | 1 | Widoczna bez interakcji | — (kontener) | — | Wyłącznie strona główna |
| Listwa strefy 3 | Zwarty pasek poziomy, trzy pozycje | Grupowanie dostępu do ustawień i funkcji globalnych | Wąski pasek, najniższa waga wizualna całej strony | 1 | Widoczna bez interakcji | — (kontener) | — | Strona główna (dodatkowo Mobile i Always On Display dostępne z każdego środowiska — rozdz. 5.6) |

Szczegółowa anatomia pojedynczej karty środowiska, kafla komponentu i pozycji listwy — wraz z pełną tabelą stanów każdego z nich — znajduje się w rozdziałach 3.4 (strefa 1), 5.5 (strefa 2) i 5.6 (strefa 3), aby uniknąć powtórzenia tej samej treści w dwóch miejscach dokumentu.

---

### 3.4. Etap 4 — Wybór środowiska (Strefa 1)

| Aspekt | Treść |
|---|---|
| Warunek wejścia | Operator znajduje się na stronie głównej (etap 3) |
| Co się dzieje | Kliknięcie jednej z czterech kart środowisk w strefie 1 wysyła polecenie nawigacyjne serwerowi; serwer otwiera lub przywraca kontekst wybranego środowiska dla nowej karty sesji. Jest to **jedyna droga wejścia** do przestrzeni roboczej środowiska — strefy 2 i 3 nie prowadzą do środowisk, lecz do okien konfiguracyjnych i ustawień |
| Stan okna | Strona główna pozostaje widoczna do chwili potwierdzenia przejścia; karta klikana przechodzi przez stan aktywny (wciśnięcie) i najechania, zanim widok zostanie zastąpiony oknem środowiska |
| Wynik | Otwarcie okna środowiska (etap 5) z boczną nawigacją modułów (TalkIn, WorkSpace, CodeStudio) albo panelem orkiestracji (MultitaskingAI), oraz z jedną, nową kartą sesji oczekującą na wybór modułu |

**Cztery karty i ich przeznaczenie:**

| Kolejność | Środowisko | Opis trybu pracy na karcie | Otwiera |
|---|---|---|---|
| 1 | TalkIn | Wiedza, komunikacja i praca z treścią | Okno środowiska z boczną nawigacją dziewięciu modułów |
| 2 | WorkSpace | Produktywność, organizacja i realizacja projektów | Okno środowiska z boczną nawigacją dziewięciu modułów |
| 3 | CodeStudio | Programowanie | Okno środowiska z boczną nawigacją ośmiu modułów |
| 4 | MultitaskingAI | Orkiestracja autonomicznej pracy ciągłej | Okno środowiska z panelem orkiestracji (sześć sekcji, rozdz. 3.6, 5.4) zamiast listy modułów |

**Tabela elementów interfejsu — Karta środowiska (strefa 1):**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Gdzie występuje |
|---|---|---|---|---|---|---|---|---|
| Karta środowiska | Duża karta wejścia (`.dn-karta--interaktywna --akcent`) | Wejście do przestrzeni roboczej jednego z czterech środowisk | Duży panel, jeden z czterech w siatce równej wagi — główny element wizualny strony | 1 | Widoczna bez interakcji | Spoczynek (powierzchnia neutralna, bez złota) · hover (akcent złoty na obrysie, uniesienie) · aktywny/fokus (akcent złoty + pierścień fokusu klawiatury) | Kliknięcie zamyka stronę główną i otwiera okno wybranego środowiska (etap 5) | Wyłącznie strefa 1 strony głównej |
| Godło środowiska | Symbol graficzny właściwy środowisku, u góry karty | Szybka identyfikacja wizualna środowiska niezależnie od tekstu | Mały element graficzny, wewnątrz karty | 1 | Widoczne bez interakcji | Statyczny (dziedziczy kolor `currentColor` z karty) | Brak interakcji własnej — część obszaru klikalnego karty | Karta środowiska |
| Tytuł środowiska | Nazwa środowiska złożona krojem nagłówkowym (Cormorant Garamond) | Nazwanie środowiska | Tekst średniej wagi, największy krój na karcie | 1 | Widoczny bez interakcji | Statyczny | Brak interakcji własnej | Karta środowiska |
| Opis trybu pracy | Jednozdaniowe streszczenie przeznaczenia środowiska | Doprecyzowanie, czym różni się to środowisko od pozostałych trzech | Mały tekst pod tytułem, krój bazowy | 1 | Widoczny bez interakcji | Statyczny | Brak interakcji własnej | Karta środowiska |

---

### 3.5. Etap 5 — Okno środowiska (przestrzeń robocza)

| Aspekt | Treść |
|---|---|
| Warunek wejścia | Kliknięcie karty środowiska (etap 4) — **albo** powrót do środowiska z otwartymi wcześniej kartami sesji (rozdz. 6) |
| Co się dzieje | Strona główna zostaje zastąpiona przestrzenią roboczą środowiska. Otwiera się nowa, pusta karta sesji (jeżeli Operator wchodzi po raz pierwszy) albo przywracana jest karta ostatnio aktywna (jeżeli środowisko było już odwiedzone w tej sesji połączenia). Boczna nawigacja (lub panel orkiestracji dla MultitaskingAI) staje się widoczna i pozostaje widoczna niezależnie od tego, który moduł jest otwarty w bieżącej karcie |
| Stan okna | Układ kolumnowy stały: boczna nawigacja / panel orkiestracji (lewa kolumna nawigacyjna), Chat Window (kolumna stała, pełna wysokość obszaru roboczego — obecne również w karcie bez wybranego modułu), obszar roboczy modułu (kolumna dominująca). Pasek marki i pasek kart sesji stanowią ramę okna. Kolumna dominująca — pusta (oczekiwanie na wybór modułu) albo wypełniony oknem operacyjnym ostatnio wybranego modułu |
| Wynik | Oczekiwanie na wybór modułu (etap 6) w nowej karcie; w karcie przywróconej — natychmiastowe wyświetlenie okna operacyjnego zapamiętanego modułu (etap 7 wprost) |

**Makieta tekstowa — Okno środowiska (przykład: TalkIn, nowa pusta karta):**

```
 ┌───────────────────────────────────────────────────────────────────────┐
 │  ◈ Danaco Console          TalkIn              [ 🔍 ]        ● Operator │  ← pasek marki
 ├───────────────┬───────────────────────────────────────────────────────┤
 │  TalkIn       │  ┌──────────────────┬──────────────────┬───────┐      │  ← pasek kart sesji
 │───────────────│  │ ●  Studio     ✕ │    Research    ✕ │   +   │      │
 │   Studio      │  └──────────────────┴──────────────────┴───────┘      │
 │   Workspace   ├──────────────────────┬────────────────────────────────┤
 │   Browser     │ Chat Window          │                                 │
 │ ▸ (brak)      │ Użytkownik ↔         │   Nowa karta sesji — pusta       │
 │   Library     │ Wykonawca            │   Wybierz moduł z bocznej        │
 │   Translate   │                      │   nawigacji, aby rozpocząć pracę │
 │   Roundtable  │ historia rozmowy      │                                 │
 │   Assistant   │                      │                                 │
 │   Agents      │ [ polecenie…  ] ➤    │                                 │
 │               │                      │                                 │
 │ [TalkIn]      │                      │                                 │
 │ [Fable 5]     │                      │                                 │
 │ Agent ▼       │                      │                                 │
 └───────────────┴──────────────────────┴────────────────────────────────┘
   kolumna           kolumna stała          kolumna dominująca —
   nawigacji         pełnej wysokości       oczekiwanie na wybór modułu
   modułów           (Chat Window)          (etap 6)

 Execution Loop Window otwiera się w kolumnie sąsiadującej z Chat Window
 w chwili przyjęcia pierwszego zlecenia; w karcie bez zlecenia kolumna
 ta jest nieobecna. Regulacji podlega wyłącznie szerokość kolumn.
```

**Tabela elementów interfejsu — Okno środowiska:**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Gdzie występuje |
|---|---|---|---|---|---|---|---|---|
| Boczna nawigacja modułów | Stała, lewa kolumna nawigacyjna z listą pozycji (`.dn-karta--pozycja`) | Wybór modułu dostępnego w bieżącym środowisku | Kolumna nawigacyjna, skrajnie lewa, pełna wysokość okna | 1 | Widoczna bez interakcji | Pozycja: domyślna · hover · wybrana (podświetlenie `▸`) | Kliknięcie przeładowuje kolumnę dominującą bieżącej karty do zestawu okien nowego modułu (rozdz. 3.6) | TalkIn, WorkSpace, CodeStudio |
| Panel orkiestracji | Odpowiednik bocznej nawigacji w środowisku MultitaskingAI — sześć sekcji zamiast listy modułów | Nawigacja po rolach, kolejkach i orkiestracji zamiast po modułach | Kolumna tej samej wagi i położenia co boczna nawigacja | 1 | Widoczny bez interakcji | Jak wyżej | Kliknięcie sekcji „Role” otwiera cztery okna robocze ról; pozostałe sekcje — właściwe okna sterowania (rozdz. 5.4) | Wyłącznie MultitaskingAI |
| Pasek kart sesji | Pasek kart w ramie okna środowiska, z kartami reprezentującymi otwarte przestrzenie robocze | Przełączanie między równolegle prowadzonymi sesjami w tym środowisku | Pasek pełnej szerokości ramy okna, waga średnia | 1 | Widoczny bez interakcji | Karta: aktywna (podświetlona) · w tle · z wskaźnikiem pracy (`●`) | Kliknięcie karty (poza „✕”) przełącza kolumny obszaru roboczego na jej stan | Każde środowisko |
| Kontrolka nowej karty „+” | Mały przycisk na końcu paska kart | Otwarcie nowej, pustej przestrzeni roboczej w bieżącym środowisku | Mała ikonka w obrębie paska kart | 1 | Widoczna bez interakcji | Domyślny · hover · aktywny | Tworzy nową kartę oczekującą na wybór modułu (powrót do etapu 6) | Każde środowisko |
| Kontrolka zamknięcia karty „✕” | Mała ikonka na karcie sesji | Zamknięcie widoku karty | Mała ikonka, widoczna przy najechaniu na kartę | 2 | Ujawniana najechaniem na kartę sesji | Domyślny (ukryta/dyskretna) · hover (widoczna) | Zamyka widok karty w bieżącym oknie; los procesu sesji leżącego u podstaw pozostaje funkcją modelu sesji (rozdz. 6) | Każde środowisko |
| Kolumna dominująca (pusta) | Komunikat pustego stanu (`.dn-pusty-stan`) | Wskazanie, że nowa karta oczekuje na wybór modułu | Duży, centralny blok z ikoną, tytułem i krótkim opisem | 1 | Widoczny bez interakcji | Wyłącznie stan „pusty” — znika po wyborze modułu | Brak interakcji własnej; wybór modułu w bocznej nawigacji zastępuje go oknem operacyjnym | Nowa, jeszcze nieużyta karta sesji |
| Pasek marki (wspólny) | Pasek marki, tożsamy ze stroną główną | Stały dostęp do godła (powrót na stronę główną), wyszukiwania i konta | Pasek pełnej szerokości ramy okna, waga niska | 1 | Widoczny bez interakcji | Statyczny; pole wyszukiwania jak w rozdz. 3.3 | Kliknięcie godła wraca do strony głównej (etap 3); karty bieżącego środowiska trwają w tle (rozdz. 6) | Każde środowisko |

---

### 3.6. Etap 6 — Wybór modułu (lub roli / sekcji orkiestracji)

| Aspekt | Treść |
|---|---|
| Warunek wejścia | Operator znajduje się w oknie środowiska (etap 5), w karcie sesji bez otwartego modułu albo z zamiarem zmiany modułu bieżącej karty |
| Co się dzieje | Kliknięcie pozycji w bocznej nawigacji (TalkIn, WorkSpace, CodeStudio) albo sekcji w panelu orkiestracji (MultitaskingAI). Klient wysyła polecenie otwarcia modułu; serwer przeładowuje kontekst karty sesji do nowego modułu — zestaw okien poprzedniego modułu (jeśli był otwarty) znika, w jego miejsce ładuje się zestaw okien nowego modułu |
| Stan okna | Krótkotrwały stan przejściowy (ładowanie zestawu okien modułu) między kliknięciem a wyświetleniem okna operacyjnego; boczna nawigacja / panel orkiestracji pozostaje widoczna i niezmieniona przez cały czas tego przejścia |
| Wynik | Pozycja modułu przyjmuje stan wybrany (podświetlenie); kolumny obszaru roboczego wypełniają się zestawem okien operacyjnych nowego modułu (etap 7) |

**Zasada wspólna:** przeładowanie obejmuje wyłącznie zestaw okien operacyjnych, kontekst czatu i narzędzia dedykowane modułowi. Nie obejmuje: samej karty sesji (pozostaje tą samą kartą), bocznej nawigacji lub panelu orkiestracji (pozostają widoczne bez zmian), pozostałych otwartych kart (zachowują własny, niezależny stan), funkcji globalnych Mobile i Always On Display (nigdy nie podlegają przeładowaniu).

**Tabela elementów interfejsu — Wybór modułu:**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Gdzie występuje |
|---|---|---|---|---|---|---|---|---|
| Pozycja modułu w bocznej nawigacji | Wiersz listy (`.dn-karta--pozycja`) z nazwą modułu | Wybór jednego z modułów dostępnych w bieżącym środowisku | Mały element listy, jedna linia tekstu | 1 | Widoczna bez interakcji | Domyślny · hover · wybrana (podświetlenie, `▸`) | Przeładowuje kolumnę dominującą bieżącej karty do zestawu okien tego modułu | Boczna nawigacja TalkIn (9 pozycji), WorkSpace (9 pozycji), CodeStudio (8 pozycji) |
| Sekcja panelu orkiestracji | Wiersz listy analogiczny do pozycji modułu, lecz nazwany funkcją sterowania (Zespoły, Role, Kolejki, Orkiestracja, Harmonogram i automatyki, Monitor procesu) | Wybór warstwy sterowania zespołem modeli i agentów | Mały element listy, jedna linia tekstu | 1 | Widoczna bez interakcji | Jak wyżej | Sekcja „Role” otwiera cztery okna robocze; pozostałe sekcje otwierają właściwe im okna sterowania (rozdz. 5.4) | Wyłącznie panel orkiestracji MultitaskingAI (6 pozycji) |
| Wskaźnik ładowania modułu | Krótkotrwały stan pośredni kolumny dominującej | Sygnalizacja trwającego przeładowania zestawu okien | Nakładka lub spinner w kolumnie dominującej | 1 | Wyświetlany samoczynnie na czas przeładowania | Ładowanie (krótkotrwałe) | Automatyczne zniknięcie po wyświetleniu okna operacyjnego | Przejście między modułami w tej samej karcie |

---

### 3.7. Etap 7 — Okno operacyjne (okno robocze)

| Aspekt | Treść |
|---|---|
| Warunek wejścia | Wybór modułu zakończony (etap 6) — **albo** powrót do karty sesji, w której moduł był już otwarty wcześniej |
| Co się dzieje | Obszar roboczy karty sesji wypełnia się zestawem okien operacyjnych właściwym wybranemu modułowi, rozmieszczonym w kolumnach sąsiadujących poziomo. Chat Window — główne okno komunikacji Użytkownik ↔ Wykonawca, wspólne wszystkim piętnastu modułom — zajmuje lewą kolumnę na pełnej wysokości obszaru roboczego i rekonfiguruje się do kontekstu modułu. Execution Loop Window — okno pętli wykonawczej Koordynator ↔ Wykonawca — otwiera się w kolumnie sąsiadującej z Chat Window. Kolumna dominująca po prawej mieści okna właściwe modułowi; okna pomocnicze otwierają się jako rozszerzenia boczne, w kolejnych kolumnach po prawej stronie obszaru roboczego. Tu zaczyna się praca właściwa — dotychczasowe etapy były wyłącznie nawigacją do tego miejsca |
| Stan okna | Pełny, interaktywny zestaw okien w układzie pionowym (podział lewa–prawa): Chat Window w lewej kolumnie stałej, Execution Loop Window w kolumnie sąsiadującej, obszar roboczy modułu (edytor, lista, monitor — zależnie od modułu) w kolumnie dominującej, okna pomocnicze w kolumnach bocznych po prawej. Regulacji podlega wyłącznie szerokość kolumn. Zawartość aktualizowana na żywo kanałem WebSocket tam, gdzie okno monitoruje proces lub odbiera strumień odpowiedzi Wykonawcy. Stan procesu sesji jest trwały po stronie serwera — rozłączenie klienta nie zamyka okna |
| Wynik | Praca w module; dalsze przejścia z tego punktu: zmiana modułu w tej samej karcie (powrót do etapu 6), nowa karta („+”, powrót do etapu 6 w nowej karcie), powrót do strony głównej (powrót do etapu 3, karta trwa w tle), uproszczone menu kontekstowe (szybka zmiana konfiguracji bieżącej sesji bez opuszczania okna) |

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
 ┌─ Uproszczone menu kontekstowe (⋮) ───────────────┐
 │  Zachowanie modeli   ›                            │
 │  Tożsamość modelu    ›                            │
 │  Historia / pamięć   ›                            │
 │  Izolacja tej karty  ›                             │
 │  ──────────────────────────                        │
 │  Otwórz pełne okno konfiguracji →                  │
 └────────────────────────────────────────────────────┘
   zmiana obejmuje wyłącznie warstwę sesji bieżącej karty —
   wartości globalne i pozostałe karty pozostają nienaruszone;
   po zamknięciu menu znika całkowicie z przestrzeni roboczej
```

**Chat Window i Execution Loop Window w przepływie — zachowanie stałe:**

| Zdarzenie przepływu | Chat Window (Użytkownik ↔ Wykonawca) | Execution Loop Window (Koordynator ↔ Wykonawca) |
|---|---|---|
| Otwarcie modułu (etap 6→7) | Otwiera się jako pierwsze okno zestawu, zawsze w lewej kolumnie obszaru roboczego, w tym samym miejscu układu w każdym module i środowisku | Otwiera się w kolumnie sąsiadującej z Chat Window w chwili przyjęcia zlecenia do wykonania; prezentuje dekompozycję zlecenia i kolejkę zadań |
| Zmiana modułu w tej samej karcie | Pozostaje otwarte i zachowuje położenie; kontekst rozmowy rekonfiguruje się do nowego modułu, historia karty pozostaje dostępna | Zamyka się wraz z zestawem okien poprzedniego modułu i otwiera ponownie w kolumnie sąsiadującej dla zlecenia nowego modułu; stan pętli poprzedniego zlecenia trwa po stronie serwera |
| Przełączenie karty sesji | Wyświetla historię rozmowy karty wybranej; historia karty opuszczonej pozostaje nienaruszona | Wyświetla stan pętli wykonawczej karty wybranej; pętla karty opuszczonej biegnie dalej po stronie serwera |
| Powrót do strony głównej (przejście 18) | Znika z widoku wraz z całą kartą; strumień rozmowy trwa po stronie serwera | Znika z widoku; pętla wykonawcza biegnie dalej, wskaźnik pracy `●` pozostaje widoczny na karcie sesji |
| Zamknięcie Execution Loop Window | Bez zmian — pozostaje otwarte, zajmuje zwolnioną szerokość | Znika całkowicie z przestrzeni roboczej; pętla biegnie dalej po stronie serwera, okno przywracane jednym kliknięciem znacznika przebiegu albo poleceniem języka naturalnego w Chat Window |
| Rozłączenie i ponowne połączenie (przejścia 22–23) | Odtwarza pełną historię rozmowy karty w lewej kolumnie | Odtwarza stan pętli: zlecenie, kolejkę zadań, wyniki kontroli jakości i wskaźniki przebiegu |

**Tabela elementów interfejsu — Okno operacyjne:**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Gdzie występuje |
|---|---|---|---|---|---|---|---|---|
| Nagłówek okna | Pasek tytułowy pojedynczego okna operacyjnego | Identyfikacja okna i jego przynależności do modułu | Wąski pasek nagłówkowy kolumny okna | 1 | Widoczny bez interakcji | Statyczny | Brak interakcji własnej poza kontrolką rozwinięcia kontekstowego | Każde okno operacyjne |
| Chat Window — pole poleceń | Pole tekstowe wieloliniowe z przyciskiem wysłania | Wprowadzenie polecenia języka naturalnego dla Wykonawcy w kontekście bieżącego modułu | Pole szerokości lewej kolumny, waga średnia — stały punkt interakcji | 1 | Widoczne bez interakcji | Domyślny · fokus · wypełnione · ładowanie (podczas generowania odpowiedzi) | Wysłanie polecenia dodaje wpis do historii i otwiera strumień odpowiedzi na żywo | Każdy z piętnastu modułów i każde środowisko (okno stałe) |
| Chat Window — historia rozmowy | Przewijana lista wymienionych wiadomości | Prezentacja przebiegu rozmowy Użytkownik ↔ Wykonawca w bieżącym kontekście | Panel wypełniający lewą kolumnę, stała, pełna wysokość obszaru roboczego | 1 | Widoczna bez interakcji | Pusta (nowa karta) · z treścią · aktualizacja na żywo (strumień) | Przewijanie, zaznaczenie fragmentu jako przedmiotu dalszej operacji | Każdy z piętnastu modułów i każde środowisko |
| Przycisk „Wyślij” | Mały przycisk ikonowy (`.dn-btn-ikona`) przy polu poleceń | Zatwierdzenie i wysłanie polecenia | Mała ikonka (36×36 px) | 1 | Widoczny bez interakcji | Domyślny · hover · aktywny · ładowanie | Zawsze klikalny; przy pustym polu poleceń kliknięcie wyświetla krótki komunikat zamiast wysyłki; przy polu wypełnionym inicjuje wysyłkę polecenia i otwarcie strumienia odpowiedzi | Chat Window, każdy moduł |
| Execution Loop Window | Okno pętli wykonawczej Koordynator ↔ Wykonawca — zlecenie i jego dekompozycja na zadania, kolejka i stan zadań, wymiana komunikatów sterujących, wyniki kontroli jakości i decyzje o ponowieniu, wskaźniki przebiegu pętli | Prowadzenie pętli wykonawczej, koordynacja zadań, nadzór nad realizacją i kontrola przebiegu procesów | Kolumna sąsiadująca z Chat Window, pełna wysokość obszaru roboczego | 1 | Otwierane przyjęciem zlecenia do wykonania; przywracane jednym kliknięciem znacznika przebiegu albo poleceniem języka naturalnego w Chat Window | Bez zlecenia (kolumna nieobecna) · zlecenie w toku · wstrzymane · zakończone · z niezgodnością zgłoszoną przez kontrolę jakości | Kliknięcie zadania w kolejce otwiera jego szczegóły w tej samej kolumnie; zamknięcie okna usuwa kolumnę całkowicie, pętla biegnie dalej po stronie serwera | Każdy moduł i każde środowisko, w którym realizowane jest zlecenie |
| Sterowanie przebiegiem pętli | Zestaw akcji Execution Loop Window: wstrzymanie, wznowienie, przerwanie, korekta zlecenia | Ingerencja Użytkownika w przebieg pętli wykonawczej | Jeden element zbiorczy `Przebieg ▼` w nagłówku kolumny | 3 | Rozwinięcie elementu zbiorczego `Przebieg ▼`; równoważnie polecenie języka naturalnego w Chat Window | Zwinięty (spoczynek) · rozwinięty | Wybór akcji zmienia stan pętli i zwija listę; rozwinięcie znika całkowicie po użyciu | Execution Loop Window |
| Kontrolka rozwinięcia kontekstowego „⋮” | Mała ikonka w nagłówku okna operacyjnego | Otwarcie uproszczonego menu szybkiej zmiany konfiguracji bieżącej sesji | Mała ikonka | 3 | Kliknięcie `⋮` w nagłówku okna | Domyślny · hover · otwarte (menu widoczne) | Otwiera menu z podzbiorem ustawień warstwy sesji; po zamknięciu menu znika całkowicie z przestrzeni roboczej | Każde okno operacyjne |
| Znaczniki kontekstowe | Lekkie znaczniki środowiska, repozytorium, projektu, modelu i wykonawcy, na przykład `[Danaco Console] [Ubuntu] [Worktree] [Fable 5] [Ultra]` | Wskazanie kontekstu pracy bieżącej karty i dostęp do właściwych selektorów | Zwarty zestaw małych znaczników w kolumnie nawigacji | 1 (znacznik) / 2 (selektor) | Znacznik widoczny bez interakcji; kliknięcie znacznika otwiera odpowiedni selektor, który zwija się samoczynnie po wyborze | Domyślny · hover · selektor otwarty | Wybór wartości aktualizuje znacznik i zwija selektor | Każde okno operacyjne |
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
| Menu progresywne | Zbiór jednorodnych wyborów prezentowany jako jeden element zwinięty (`Agent ▼`, `Przebieg ▼`); lista pozycji rozwija się po kliknięciu | rozdz. 3.7, 5.4 |
| Panele wysuwane | Okna pomocnicze modułu otwierane jako rozszerzenia boczne w kolumnach po prawej stronie obszaru roboczego; po zamknięciu kolumna znika całkowicie | rozdz. 3.7, 5.1–5.3 |
| Grupowanie logiczne akcji | Zamiast zestawu przycisków jeden element zbiorczy (`Operacje ▼`), którego rozwinięcie zawiera pełną listę akcji | rozdz. 3.7, 5.5 |
| Znaczniki kontekstowe | Środowisko, repozytorium, projekt, model i wykonawca jako lekkie znaczniki `[Danaco Console] [Ubuntu] [Worktree] [Fable 5] [Ultra]`; kliknięcie znacznika otwiera odpowiedni selektor | rozdz. 3.5, 3.7 |

**Makiety w stanie spoczynku.** Wszystkie makiety tekstowe tego dokumentu rysowane są w stanie spoczynku interfejsu: widoczne są wyłącznie elementy warstwy 1 oraz zwinięte wyzwalacze warstw 2–3 (znaczniki, `▼`, `⋮`, `☰`). Elementy warstw 2–4 opisane są wyłącznie w tabelach elementów okna, z podaniem warstwy i sposobu wywołania.

**Warstwy a etapy przepływu:**

| Etap | Elementy warstwy 1 | Elementy warstw 2–4 obecne na etapie |
|---|---|---|
| 1 — Okno startowe | Godło marki, wskaźnik stanu połączenia | Kontrolka ponowienia (warstwa 2, widoczna wyłącznie w stanie błędu) |
| 2 — Okno rejestracji i logowania | Zakładki, pola formularza, przycisk zatwierdzenia, przycisk „Pomiń” | Wybór metody logowania (warstwa 2), odzyskiwanie konta (warstwa 2, wywoływane odnośnikiem) |
| 3 — Strona główna | Trzy strefy, pasek marki | Wyszukiwarka funkcji (warstwa 4, skrót klawiszowy i pole wyszukiwania paska marki) |
| 4 — Wybór środowiska | Cztery karty środowisk | — |
| 5 — Okno środowiska | Boczna nawigacja lub panel orkiestracji, pasek kart sesji, znaczniki kontekstowe | Selektory otwierane znacznikiem (warstwa 2), zestaw akcji karty sesji (warstwa 3) |
| 6 — Wybór modułu | Pozycja modułu lub sekcja panelu orkiestracji | — |
| 7 — Okno operacyjne | Chat Window, Execution Loop Window, okno wiodące modułu | Okna pomocnicze (warstwa 2), uproszczone menu kontekstowe i sterowanie przebiegiem pętli (warstwa 3), administracja procesami i diagnostyka niskiego poziomu (warstwa 4) |

---

## 4. Tabela przejść — zestawienie zbiorcze

Poniższa tabela zestawia wszystkie przejścia opisane w rozdziałach 2, 3 i 6, w kolejności zbliżonej do przebiegu głównego, z rozgałęzieniami zgrupowanymi bezpośrednio pod przejściem macierzystym.

| Nr | Grupa | Przejście (Z → Do) | Warunek | Co się dzieje | Stan okna docelowego |
|---|---|---|---|---|---|
| 1 | Rdzeń | Uruchomienie aplikacji → Okno startowe | Operator otwiera pakiet kliencki | Otwarcie natywnego okna; nawiązanie kanału WebSocket; `connection.hello` | Ładowanie (łączenie z serwerem) |
| 2 | Rdzeń | Okno startowe → Okno rejestracji i logowania | Brak ważnego tokenu urządzenia LUB faza budowy | Serwer potwierdza fazę wdrożenia; klient prezentuje formularz | Interaktywny formularz (logowanie domyślnie) |
| 2a | Wariant | Okno startowe → Strona główna (z pominięciem logowania) | Token urządzenia ważny, brak zapamiętanej aktywnej karty | Krok „Uwierzytelnienie” spełniony ważnym tokenem; `home.enter` | Stabilny, trzy strefy |
| 2b | Wariant | Okno startowe → Okno operacyjne (z pominięciem logowania i strony głównej) | Token ważny i klient zgłasza powrót do zapamiętanej aktywnej karty | `session.bind`; serwer przywraca pełny stan karty (układ, historia, kontekst) | Pełny, jak przed rozłączeniem |
| 3 | Rdzeń | Okno logowania → Strona główna (rejestracja) | Konto nie istnieje; dane wymagane podane i e-mail potwierdzony | Utworzenie jedynej encji Konto; wydanie tokenu urządzeniu | Stabilny, trzy strefy |
| 4 | Rdzeń | Okno logowania → Strona główna (logowanie) | Konto istnieje; login+hasło lub metoda dodatkowa poprawne | Weryfikacja przez serwer; wydanie tokenu urządzeniu | Stabilny, trzy strefy |
| 5 | Wariant | Okno logowania → Strona główna (Pomiń) | Ustawienie „Wymóg logowania” nieaktywne (domyślnie w fazie budowy, konfigurowalne przez Operatora w każdej fazie) | Wejście bez uwierzytelnienia; krok „Uwierzytelnienie” pomijany | Stabilny, trzy strefy |
| 6 | Wariant | Okno logowania → Odzyskiwanie konta → Okno logowania | Operator wskazuje utratę dostępu i podaje e-mail uwierzytelniający | Wysyłka drogi potwierdzenia; ustawienie nowego hasła; unieważnienie dawnych tokenów | Formularz logowania z nowym hasłem |
| 7 | Wariant | Okno logowania (dane niepoprawne) → Okno logowania | Weryfikacja negatywna | Komunikat błędu przy polu; formularz pozostaje wypełniony | Formularz, stan błędu |
| 8 | Rdzeń | Strona główna → Okno środowiska | Kliknięcie karty środowiska (strefa 1) | Zamknięcie widoku strony głównej; otwarcie nowej lub przywróconej karty sesji środowiska | Układ kolumnowy: kolumna nawigacji, Chat Window, kolumna dominująca; pasek marki i pasek kart sesji jako rama okna |
| 9 | Rozgałęzienie | Strona główna → Okno budowy komponentu własnego | Kliknięcie kafla komponentu (strefa 2) | Otwarcie okna/nakładki budowy komponentu (Agent Builder i analogiczne — rozdz. 5.5) | Formularz budowy komponentu |
| 10 | Rozgałęzienie | Strona główna → Okno konfiguracji | Kliknięcie pozycji „Okno konfiguracji” (strefa 3) | Otwarcie nakładki modalnej nad bieżącym widokiem | Nakładka modalna, trzynaście zakresów w nawigacji |
| 11 | Rozgałęzienie | Strona główna → Aktywacja Mobile | Kliknięcie pozycji „Mobile” (strefa 3) | Aktywacja funkcji globalnej — bez nowej przestrzeni roboczej | Bieżący widok bez zmian + nakładka/panel Mobile |
| 12 | Rozgałęzienie | Strona główna → Aktywacja Always On Display | Kliknięcie pozycji „Always On Display” (strefa 3) | Aktywacja globalnego agenta towarzyszącego | Bieżący widok bez zmian + pływający avatar AOD |
| 13 | Rdzeń | Okno środowiska → Wybór modułu | Kliknięcie pozycji bocznej nawigacji (lub sekcji panelu orkiestracji) | Przeładowanie kolumny dominującej bieżącej karty | Krótkotrwałe ładowanie zestawu okien modułu |
| 14 | Rdzeń | Wybór modułu → Okno operacyjne | Zestaw okien modułu załadowany | Chat Window w lewej kolumnie, Execution Loop Window w kolumnie sąsiadującej i okna właściwe modułowi wypełniają obszar roboczy | Pełny, interaktywny zestaw okien |
| 15 | Pętla | Okno operacyjne → Wybór modułu (ta sama karta) | Kliknięcie innej pozycji bocznej nawigacji | Zamknięcie okien poprzedniego modułu; otwarcie zestawu nowego modułu | Jak poz. 13–14 |
| 16 | Pętla | Okno środowiska → nowa karta sesji → Wybór modułu | Kliknięcie kontrolki „+” paska kart | Nowa, pusta przestrzeń robocza; poprzednie karty niezmienione w tle | Kolumna dominująca pusta — oczekiwanie na moduł |
| 17 | Pętla | Karta sesji → zamknięcie karty | Kliknięcie kontrolki „✕” | Zamknięcie widoku karty w oknie bieżącym | Karta znika z paska; pozostałe karty niezmienione |
| 18 | Pętla | Okno środowiska / okno operacyjne → Strona główna | Kliknięcie godła w pasku marki | Powrót do centrum dowodzenia; karty bieżącego środowiska trwają w tle jako procesy serwera | Stabilny, trzy strefy |
| 19 | Pętla | Strona główna → Okno środowiska (powrót) | Ponowne kliknięcie karty środowiska odwiedzonego wcześniej w tej sesji połączenia | Przywrócenie pełnego stanu wszystkich otwartych tam kart (układ, historia, kontekst) | Jak przed opuszczeniem |
| 20 | Wariant | Okno operacyjne → rozwinięcie kontekstowe → Okno operacyjne | Kliknięcie kontrolki „⋮” w nagłówku okna | Szybka zmiana ustawień warstwy sesji bez opuszczania okna | Okno operacyjne niezmienione poza zaktualizowaną wartością ustawienia |
| 21 | Rozgałęzienie | Okno budowy komponentu własnego → zapis → Strona główna lub kontynuacja | Zatwierdzenie definicji komponentu | Zapis jako komponent własny w pamięci aplikacji; wytwór staje się zasobem wybieralnym | Strona główna albo dalsza edycja w tym samym oknie |
| 22 | Wariant | Rozłączenie klienta (dowolny moment) | Zamknięcie okna aplikacji, utrata sieci | Kanał WebSocket zamknięty; proces sesji i wszystkie procesy podrzędne trwają dalej po stronie serwera | Aplikacja zamknięta po stronie klienta; brak zmiany stanu po stronie serwera |
| 23 | Wariant | Ponowne połączenie → Okno startowe → stan zastany | Ponowne uruchomienie aplikacji na tym samym lub innym urządzeniu | Etapy 1–3 cyklu życia połączenia powtórzone; token urządzenia przedstawiony | Odtworzenie stanu sprzed rozłączenia (poz. 2b) albo strona główna (poz. 2a) |

---

## 5. Ścieżki szczegółowe

Poniższe sześć podrozdziałów rozwija etapy 4–7 (rozdz. 3.4–3.7) dla każdego z czterech środowisk z osobna oraz dla dwóch pozostałych stref strony głównej. Każda ścieżka pokazuje pełny przebieg od karty środowiska (lub kafla/pozycji listwy) do konkretnego okna operacyjnego, na jednym reprezentatywnym przykładzie modułu.

### 5.1. Ścieżka TalkIn

Środowisko wiedzy, komunikacji i pracy z treścią. Boczna nawigacja udostępnia dziewięć modułów: Studio, Workspace, Browser, Research, Library, Translate, Roundtable, Assistant, Agents.

```
 STRONA GŁÓWNA
     │  Strefa 1 → karta „TalkIn”
     ▼
 OKNO ŚRODOWISKA TalkIn
     │  boczna nawigacja (9 pozycji widocznych)
     │  ┌─────────────┬─────────────┬──────────┬────────────┬──────────┐
     │  │ Studio      │ Workspace   │ Browser  │ Research   │ Library  │
     │  ├─────────────┼─────────────┼──────────┼────────────┼──────────┤
     │  │ Translate   │ Roundtable  │ Assistant│ Agents     │          │
     │  └─────────────┴─────────────┴──────────┴────────────┴──────────┘
     ▼  (przykład: wybór „Studio”)
 OKNO OPERACYJNE — moduł Studio (stan spoczynku, układ kolumnowy)
 ═══════════════════════════════════════════════════════════════════════
  Chat Window          │ Studio Editor          │ Tools Panel
  Użytkownik ↔         │ treść dokumentu:       │ (rozszerzenie boczne)
  Wykonawca            │ PDF·DOCX·TXT·Markdown  │ operacje kontekstowe
                       │                        │ Wykonawcy wobec
  [ polecenie…  ] ➤    │                    [⋮] │ dokumentu
  ──────────────────   │                        │
  Execution Loop       │                        │
  Koordynator ↔        │                        │
  Wykonawca            │                        │
 ═══════════════════════════════════════════════════════════════════════
  kolumna stała          kolumna dominująca       kolumny boczne
  pełnej wysokości                                (Tools Panel, Diff/Grep
                                                   Panel, Preview Window,
                                                   Session Repository —
                                                   każde jako rozszerzenie
                                                   boczne po prawej)
     │
     │  „+” nowa karta w tym samym środowisku → wybór „Research”
     ▼
 DRUGA KARTA SESJI — moduł Research (równolegle do karty Studio)
 ═══════════════════════════════════════════════════════════════════════
  Chat Window          │ Research Workspace     │ Sources Manager
  Użytkownik ↔         │ zakres · przebieg      │ źródła z metadanymi
  Wykonawca            │ badania                │ (rozszerzenie boczne)
  ──────────────────   │                    [⋮] │
  Execution Loop       │                        │ Findings Panel ·
  Koordynator ↔        │                        │ Report Builder ·
  Wykonawca            │                        │ Export Panel —
                       │                        │ kolejne kolumny boczne
 ═══════════════════════════════════════════════════════════════════════

 Obie karty (Studio, Research) działają jednocześnie — pasek kart:
 [ ● Studio  ✕ ] [ Research  ✕ ] [ + ]
```

**Tabela przejść — ścieżka TalkIn:**

| Krok | Warunek | Co się dzieje | Stan okna |
|---|---|---|---|
| Strona główna → TalkIn | Kliknięcie karty „TalkIn” | Otwarcie okna środowiska, nowa pusta karta | Boczna nawigacja aktywna, Chat Window w lewej kolumnie, kolumna dominująca pusta |
| TalkIn → moduł Studio | Kliknięcie pozycji „Studio” | Załadowanie zestawu: Studio Editor, Tools Panel, Diff/Grep Panel, Session Repository, Preview Window, Chat Window | Pełny zestaw okien, edytor pusty do chwili wczytania dokumentu |
| Karta Studio → nowa karta Research | Kliknięcie „+”, następnie „Research” w bocznej nawigacji nowej karty | Druga, niezależna karta z własnym kontekstem; karta Studio niezmieniona w tle | Dwie karty aktywne jednocześnie w pasku kart |

### 5.2. Ścieżka WorkSpace

Środowisko produktywności, organizacji i realizacji projektów. Boczna nawigacja udostępnia dziewięć modułów: Studio, Workspace, Browser, Research, Library, Roundtable, Design, Apps, Agents.

```
 STRONA GŁÓWNA
     │  Strefa 1 → karta „WorkSpace”
     ▼
 OKNO ŚRODOWISKA WorkSpace
     │  boczna nawigacja (9 pozycji widocznych)
     │  ┌─────────────┬─────────────┬──────────┬────────────┬──────────┐
     │  │ Studio      │ Workspace   │ Browser  │ Research   │ Library  │
     │  ├─────────────┼─────────────┼──────────┼────────────┼──────────┤
     │  │ Roundtable  │ Design      │ Apps     │ Agents     │          │
     │  └─────────────┴─────────────┴──────────┴────────────┴──────────┘
     ▼  (przykład: wybór modułu „Apps” — budowa produktu cyfrowego)
 OKNO OPERACYJNE — moduł Apps (stan spoczynku, układ kolumnowy)
 ═══════════════════════════════════════════════════════════════════════
  Chat Window          │ Product Builder        │ Architecture Designer
  Użytkownik ↔         │ punkt wejścia budowy   │ komponenty rozwiązania
  Wykonawca            │ produktu               │ i zależności
                       │                    [⋮] │ (rozszerzenie boczne)
  [ polecenie…  ] ➤    │                        │
  ──────────────────   │                        │
  Execution Loop       │                        │
  Koordynator ↔        │                        │
  Wykonawca            │                        │
  zlecenie · zadania   │                        │
 ═══════════════════════════════════════════════════════════════════════
  kolumna stała          kolumna dominująca       kolumny boczne
  pełnej wysokości                                (Architecture Designer,
                                                   Frontend Workspace,
                                                   Backend Workspace,
                                                   Deployment Panel —
                                                   każde jako rozszerzenie
                                                   boczne po prawej)
```

**Tabela przejść — ścieżka WorkSpace:**

| Krok | Warunek | Co się dzieje | Stan okna |
|---|---|---|---|
| Strona główna → WorkSpace | Kliknięcie karty „WorkSpace” | Otwarcie okna środowiska, nowa pusta karta | Boczna nawigacja aktywna, Chat Window w lewej kolumnie, kolumna dominująca pusta |
| WorkSpace → moduł Apps | Kliknięcie pozycji „Apps” | Załadowanie: Product Builder, Architecture Designer, Frontend Workspace, Backend Workspace, Deployment Panel, Chat Window | Product Builder jako punkt wejścia; Deployment Panel dostępny od razu — kliknięcie przed zakończeniem prac frontend/backend wyświetla komunikat o niezakończonych etapach |
| Apps → moduł Workspace (projekt) | Operator otwiera dodatkowo moduł Workspace w nowej karcie, aby zarządzać projektem nadrzędnym | Project Dashboard, Instructions Panel, Context Memory, Project Library, Agent Manager | Druga karta niezależna; Agent Manager pozwala przypisać agentów z modułu Agents do tego projektu |

### 5.3. Ścieżka CodeStudio

Środowisko programistyczne. Boczna nawigacja udostępnia osiem modułów: Workspace, Roundtable, Design, Terminal, Developer, Diagnostics, Apps, Agents.

```
 STRONA GŁÓWNA
     │  Strefa 1 → karta „CodeStudio”
     ▼
 OKNO ŚRODOWISKA CodeStudio
     │  boczna nawigacja (8 pozycji widocznych)
     │  ┌─────────────┬─────────────┬──────────┬────────────┐
     │  │ Workspace   │ Roundtable  │ Design   │ Terminal   │
     │  ├─────────────┼─────────────┼──────────┼────────────┤
     │  │ Developer   │ Diagnostics │ Apps     │ Agents     │
     │  └─────────────┴─────────────┴──────────┴────────────┘
     ▼  (przykład: wybór modułu „Developer”)
 OKNO OPERACYJNE — moduł Developer (stan spoczynku, układ kolumnowy)
 ═══════════════════════════════════════════════════════════════════════
  Chat Window       │ Project │ Code Editor        │ Git Panel
  Użytkownik ↔      │ Tree    │ treść plików       │ status zmian
  Wykonawca         │ struk-  │ źródłowych         │ historia rewizji
                    │ tura    │ (wsparcie          │ (rozszerzenie
  [ polecenie…] ➤   │ katalo- │  Wykonawcy)    [⋮] │  boczne)
  ───────────────   │ gów     │                    │
  Execution Loop    │         │                    │
  Koordynator ↔     │         │                    │
  Wykonawca         │         │                    │
  zlecenie · zadania│         │                    │
 ═══════════════════════════════════════════════════════════════════════
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
| Strona główna → CodeStudio | Kliknięcie karty „CodeStudio” | Otwarcie okna środowiska, nowa pusta karta | Boczna nawigacja aktywna, Chat Window w lewej kolumnie, kolumna dominująca pusta |
| CodeStudio → moduł Developer | Kliknięcie pozycji „Developer” | Załadowanie: Project Tree, Code Editor, Git Panel, Build Output, Chat Window | Project Tree odzwierciedla strukturę repozytorium powiązanego z sesją |
| Developer → Terminal (powiązanie jawne) | Operator, świadomą decyzją, otwiera dodatkowo moduł Terminal (nowa karta lub powiązanie skonfigurowane) | Terminal Tabs, Output Console, Process Monitor | Powiązanie nie jest domyślne — wymaga jawnej konfiguracji, zgodnie z zasadą pełnej konfigurowalności |

### 5.4. Ścieżka MultitaskingAI

Środowisko orkiestracji autonomicznej pracy ciągłej. **Różni się strukturalnie od pozostałych trzech** — nie ma bocznej nawigacji modułów; w jej miejscu panel orkiestracji o sześciu sekcjach, nawigujący po rolach, kolejkach i orkiestracji zamiast po modułach.

```
 STRONA GŁÓWNA
     │  Strefa 1 → karta „MultitaskingAI”
     ▼
 OKNO ŚRODOWISKA MultitaskingAI
     │  panel orkiestracji (6 sekcji, zamiast listy modułów)
     │  ┌────────────────────────────┐
     │  │  Zespoły                    │  presety ról, powiązań, kolejek
     │  │  Role                       │  → 4 okna robocze (poniżej)
     │  │  Kolejki                    │  → Queue Manager (Automations)
     │  │  Orkiestracja                │  → Orchestrator (Automations)
     │  │  Harmonogram i automatyki    │  → Scheduler, Execution Monitor
     │  │  Monitor procesu             │  → Always On Display, Mobile
     │  └────────────────────────────┘
     ▼  (wybór sekcji „Role”)
 CZTERY OKNA ROBOCZE RÓL — jednoczesny widok zespołu w kolumnach sąsiadujących
 ═════════════════════════════════════════════════════════════════════════
  Chat Window      │ Executor Chat   │ Executor Chat   │ Results Analyzer
  Użytkownik ↔     │ (Executor 1)    │ (Executor 2)    │ (Executor 3 /
  Wykonawca        │ wykonanie zadań │ drugi, równo-   │  Validator)
                   │ z kolejki       │ legły tor pracy │ ocena wyników ·
  [ polecenie…] ➤  │ Subagent        │ Subagent        │ zgłaszanie
  ───────────────  │ Network (do 15) │ Network (do 15) │ niezgodności ·
  Execution Loop   │ tryb współpracy │ niezależna ·    │ rozstrzyganie
  Koordynator ↔    │ z Executorem 2  │ przekazywanie   │ konfliktów
  Wykonawca        │ wg ustalenia    │ wyników ·       │
  (Coordinator     │ Koordynatora    │ naprzemienna ·  │ wcielenie
   Chat)           │                 │ iteracyjna      │ konfigurowalne:
  plan etapów ·    │                 │                 │ Validator ·
  podział pracy ·  │                 │                 │ Reviewer ·
  budowa promptów ·│                 │                 │ Security Auditor ·
  stan kolejki     │                 │                 │ Architect ·
  Przebieg ▼       │             [⋮] │             [⋮] │ Product Owner ·
                   │                 │                 │ QA Lead ·
                   │                 │                 │ Arbitrator
 ═════════════════════════════════════════════════════════════════════════
  kolumna stała      kolumna           kolumna           kolumna boczna
  pełnej wysokości   dominująca        sąsiadująca       (rozszerzenie)

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
| Strona główna → MultitaskingAI | Kliknięcie karty „MultitaskingAI” | Otwarcie okna środowiska z panelem orkiestracji zamiast bocznej nawigacji | Panel orkiestracji aktywny, Chat Window w lewej kolumnie, kolumna dominująca pusta |
| MultitaskingAI → sekcja „Role” | Kliknięcie sekcji „Role” | Otwarcie czterech okien roboczych jednocześnie: Coordinator Chat, dwa Executor Chat, Results Analyzer | Cztery okna widoczne równolegle — nie sekwencyjnie |
| Rola → przypisanie wykonawcy | Operator przypisuje agenta z modułu Agents (strefa 2) albo model bazowy do roli | Agent lub model staje się wykonawcą roli; wybór kanału modelu (API/CLI/SSH/HTTP) per rola | Okno roli aktywne, gotowe do sterowania przez Coordinatora |
| Role → sekcja „Harmonogram i automatyki” | Operator konfiguruje powiązanie z modułem Automations (świadoma decyzja, nie domyślna) | Scheduler i Execution Monitor modułu Automations stają się operacyjnym zapleczem harmonogramu | Praca ciągła 24/7/365 możliwa po tej konfiguracji |
| Rola → sekcja „Monitor procesu” | Kliknięcie sekcji „Monitor procesu” | Otwarcie widoku nadzoru; Always On Display może pełnić funkcję obserwatora lub operatora | Widok na żywo statusów przebiegów i hierarchii decyzji |

### 5.5. Ścieżka Strefy 2 — komponenty własne

Strefa 2 strony głównej udostępnia cztery kafle: Automations, Agents, Workspace, Assistant. Każdy otwiera okno budowy właściwego rodzaju komponentu własnego — **nie** przestrzeń roboczą środowiska. Wytwór, po zapisie, trafia do pamięci aplikacji jako zasób Operatora i staje się wybieralny operacyjnie w środowiskach (etapy 6–7).

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
              (etap 7)            Automations)          MultitaskingAI)
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
| Kafel komponentu własnego | Karta (`.dn-karta`, bez wariantu `--akcent`) | Wejście do okna budowy jednego rodzaju komponentu własnego | Średni panel, lżejszy wizualnie niż karta środowiska — brak godła w skali karty środowiska, brak kroju nagłówkowego | 1 | Widoczny bez interakcji | Spoczynek · hover (tło subtelne, bez akcentu złotego na spoczynku) · aktywny/fokus | Kliknięcie otwiera modal albo dedykowane okno budowy komponentu (nie: przestrzeń roboczą środowiska) | Wyłącznie strefa 2 strony głównej |
| Ikona komponentu | Mały symbol graficzny właściwy rodzajowi komponentu | Szybka identyfikacja rodzaju (automatyka, agent, projekt, profil głosowy) | Mała ikonka, mniejsza skala niż godło środowiska | 1 | Widoczna bez interakcji | Statyczny | Brak interakcji własnej | Kafel komponentu własnego |
| Nazwa komponentu | Automations / Agents / Workspace / Assistant | Nazwanie rodzaju komponentu | Tekst średniej wagi, krój bazowy (nie nagłówkowy) | 1 | Widoczna bez interakcji | Statyczny | Brak interakcji własnej | Kafel komponentu własnego |
| Etykieta działania | Krótkie sformułowanie zorientowane na tworzenie („Zbuduj automatykę”, „Skonfiguruj agenta”, „Załóż projekt”, „Ustaw profil asystenta”) | Odróżnienie charakteru akcji od karty środowiska („wejdź” vs. „zbuduj”) | Mały tekst pod nazwą komponentu | 1 | Widoczna bez interakcji | Statyczny | Brak interakcji własnej | Kafel komponentu własnego |
| Agent Builder — okno nadrzędne | Formularz spinający cztery okna konfiguracyjne w jeden zapisywalny wytwór | Skompletowanie pełnej definicji agenta | Duże okno / dedykowany widok budowy | 1 | Kliknięcie kafla „Agents” albo polecenie języka naturalnego w Chat Window | Pusty (nowy agent) · wypełniany · gotowy do zapisu · zapisany | Zapis tworzy komponent własny „agent”, wybieralny dalej w środowiskach i rolach | Strefa 2 (kafel „Agents”) oraz okno modułowe Agents w środowiskach |
| Przycisk „Zapisz agenta” | Przycisk główny CTA w Agent Builder | Zatwierdzenie i zapisanie definicji agenta jako komponentu własnego | Duży, pojedynczy przycisk akcentowany | 1 | Widoczny bez interakcji | Domyślny · hover · aktywny · zapisano (potwierdzenie) | Zawsze klikalny; przy niekompletnej definicji kliknięcie zapisuje agenta i wyświetla ostrzeżenie wskazujące brakujące elementy — zapis nie jest blokowany; toast potwierdzający; agent staje się wybieralny | Agent Builder |

### 5.6. Ścieżka Strefy 3 — konfiguracja, ustawienia, funkcje globalne

Strefa 3 strony głównej udostępnia listwę ustawień: pozycję „Okno konfiguracji” oraz dwie funkcje globalne, Mobile i Always On Display. W odróżnieniu od stref 1 i 2, wszystkie trzy pozycje są dostępne **z każdego miejsca platformy**, nie tylko ze strony głównej — obecność w listwie jest jednym z kilku równoważnych punktów dostępu.

```
 STRONA GŁÓWNA — Strefa 3                    (dostępne też z KAŻDEGO środowiska
     │                                         i KAŻDEGO okna operacyjnego)
     ├── „Okno konfiguracji” ──► OKNO KONFIGURACJI
     │                            nakładka modalna nad bieżącym widokiem
     │                            ┌─────────────────────┬──────────────────┐
     │                            │ NAWIGACJA (13 zakr.) │ ZAWARTOŚĆ ZAKRESU │
     │                            │ Aplikacja             │ pozycje ustawień  │
     │                            │ Procesy               │ z wartością,      │
     │                            │ Akcje                 │ warstwą (globalna │
     │                            │ Zachowanie modeli     │ / środowisko /    │
     │                            │ Tożsamość modeli      │ projekt / sesja)  │
     │                            │ Prompty systemowe      │ i objaśnieniem[?] │
     │                            │ Rozszerzenia           │                   │
     │                            │ Integracje             │                   │
     │                            │ Historia               │                   │
     │                            │ Pamięć                 │                   │
     │                            │ Izolacja  ────────────►│ (rozwinięcie niżej)│
     │                            │ Karty sesji             │                   │
     │                            │ Komponenty własne       │                   │
     │                            └─────────────────────┴──────────────────┘
     │                                    │  pozycja „Izolacja”
     │                                    ▼
     │                     OKNO KONFIGURACJI PUNKTÓW IZOLACJI (trzy panele)
     │                     ┌────────────────┬──────────────────┬───────────────┐
     │                     │ SELEKTOR        │ MACIERZ IZOLACJI  │ PROFIL I       │
     │                     │ ZASIĘGU (7 poz.)│ (3 kontekst. +     │ PODGLĄD        │
     │                     │ globalny·środo- │  8 technicznych,   │ zapisz/wczytaj/│
     │                     │ wisko·moduł·    │  każdy: przełącznik│ przypisz profil│
     │                     │ para modułów·   │  z objaśnieniem[?])│ warstwa: domyśl│
     │                     │ projekt·karta   │                    │ na / sesji     │
     │                     │ sesji·rola      │                    │ polityka efekt.│
     │                     └────────────────┴──────────────────┴───────────────┘
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

Zakres „Aplikacja” obejmuje m.in. ustawienie „Wymóg logowania”, którym Operator określa — niezależnie od fazy wdrożenia — czy okno rejestracji i logowania (rozdz. 3.2) wymusza skuteczne uwierzytelnienie przed wejściem na platformę; domyślnie nieaktywne w fazie budowy, aktywne w fazie produkcyjnej, zmienialne w każdej chwili.

**Tabela elementów interfejsu — listwa ustawień (strefa 3):**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Gdzie występuje |
|---|---|---|---|---|---|---|---|---|
| Pozycja „Okno konfiguracji” | Element listwy (`.dn-btn-ikona` + etykieta) | Otwarcie pełnego zakresu ustawień platformy | Mały element paska, waga najniższa | 1 | Widoczna bez interakcji; równoważnie skrót klawiszowy i wyszukiwarka funkcji | Domyślny · hover · fokus | Otwiera nakładkę modalną z trzynastoma zakresami ustawień | Listwa strony głównej i każde środowisko/okno operacyjne |
| Pozycja „Mobile” | Element listwy analogiczny | Aktywacja globalnego trybu dostępu mobilnego | Mały element paska | 1 | Widoczna bez interakcji | Jak wyżej | Aktywuje funkcję — nie tworzy nowej przestrzeni roboczej, nie przeładowuje niczego | Listwa strony głównej i każde środowisko/okno operacyjne |
| Pozycja „Always On Display” | Element listwy analogiczny | Aktywacja globalnego agenta towarzyszącego | Mały element paska | 1 | Widoczna bez interakcji | Jak wyżej | Aktywuje pływający avatar obecny nad całą platformą | Listwa strony głównej i każde środowisko/okno operacyjne |
| Nakładka modalna okna konfiguracji | Duży panel (`.dn-modal`) nad przyciemnionym tłem (`.dn-nakladka`) | Konfiguracja dowolnej zależności lub zachowania platformy z jednego miejsca | Duży panel, dwudzielny układ (nawigacja + zawartość) | 2 | Kliknięcie pozycji listwy, rozwinięcie kontekstowe `⋮` okna operacyjnego, skrót klawiszowy albo wyszukiwarka funkcji; po zamknięciu znika całkowicie z przestrzeni roboczej | Otwarta · zamknięta (powrót do widoku pod spodem, niezmienionego) | Zamknięcie przywraca dokładnie stan widoku sprzed otwarcia | Wywoływana z listwy strony głównej albo uproszczonego menu kontekstowego okna operacyjnego |
| Przełącznik macierzy izolacji | Suwak (`.dn-suwak`) | Ustawienie „współdzielone/odrębne” (kontekst) albo „włączony/wyłączony” (technika) dla jednej pozycji izolacji | Mały element w wierszu tabeli/listy | 2 | Widoczny po otwarciu okna konfiguracji punktów izolacji | Domyślny (wyłączony/odrębny) · włączony/współdzielony · hover · fokus | Zmiana natychmiast widoczna w podglądzie polityki efektywnej panelu prawego | Okno konfiguracji punktów izolacji |
| Objaśnienie kontekstowe `[?]` | Ikona z dymkiem (`.dn-tooltip`) | Wyjaśnienie działania ustawienia i jego wpływu na aplikację | Bardzo mała ikonka przy każdej pozycji ustawienia | 3 | Najechanie lub fokus na ikonie `[?]`; dymek znika po odsunięciu wskaźnika | Domyślny (ukryty tekst) · najechanie/fokus (dymek widoczny) | Wyświetla treść objaśnienia, nie zmienia wartości ustawienia | Każda pozycja konfiguracji w oknie konfiguracji |

---

## 6. Karty sesji i trwałość procesu w tle

Mechanizm kart sesji przenika etapy 5–7 i tłumaczy, dlaczego „stan okna” w tym dokumencie tak często rozróżnia „nową kartę” od „karty przywróconej”. Każda sesja — niezależnie od tego, czy jej karta jest obecnie widoczna — działa jako odrębny proces po stronie serwera.

```
 UTWORZENIE KARTY SESJI  (etap 6, kliknięcie „+” albo pierwszy wybór środowiska)
        │  serwer uruchamia odrębny proces (własny katalog roboczy, środowisko procesu)
        ▼
 ┌──────────────────┐   rozłączenie klienta      ┌────────────────────────────┐
 │  AKTYWNA          │ ─────────────────────────►│  AKTYWNA BEZ KLIENTA        │
 │  karta widoczna,  │   (zamknięcie aplikacji,   │  proces trwa po stronie     │
 │  klient podłączony │    utrata sieci, powrót    │  serwera; karta niewidoczna │
 │                    │    do strony głównej)      │  w oknie klienta            │
 └─────────┬──────────┘ ◄─────────────────────────└──────────────┬─────────────┘
           │              ponowne połączenie /                    │
           │              powrót do środowiska                    │
           │              (odtworzenie pełnego stanu:              │
           │               układ, historia, kontekst)               │
           └───────────────────────◄──────────────────────────────┘
```

| Zdarzenie | Wpływ na kartę bieżącą | Wpływ na pozostałe karty |
|---|---|---|
| Powrót do strony głównej (etap 3, przejście 18) | Znika z widoku; proces trwa dalej na serwerze | Bez zmian — trwają w tle na równi z kartą opuszczoną |
| Wybór innego środowiska na stronie głównej | Poprzednie środowisko i wszystkie jego karty pozostają aktywne w tle | Odtwarzają pełny stan przy powrocie do tego środowiska |
| Zamknięcie okna aplikacji (rozłączenie, przejście 22) | Kanał WebSocket zamknięty; stan zapisany trwale (SQLite + pliki treści) | Wszystkie procesy sesji tego Operatora trwają dalej po stronie serwera |
| Ponowne uruchomienie aplikacji (przejście 23) | Etapy 1–3 cyklu życia połączenia powtórzone | Karty odtwarzają się dokładnie w stanie sprzed rozłączenia |
| Zmiana modułu w tej samej karcie (etap 6→7 w pętli) | Przeładowanie zestawu okien tej jednej karty; Chat Window pozostaje otwarte w lewej kolumnie i rekonfiguruje kontekst | Bez wpływu na pozostałe karty |
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

Magistrala przechowuje referencje do artefaktów, nie ich kopie. Przekazanie nie tworzy nowego artefaktu ani nowej wersji — moduł docelowy pracuje na tym samym wpisie encji `artefakt` (`architektura/model-danych.md`, rozdz. 16.1), a każde utrwalenie zmiany powstaje jako kolejna `wersja_artefaktu` (rozdz. 16.2) tego samego artefaktu.

**Waga wizualna i miejsce w układzie.** Kolumna boczna, otwierana jako rozszerzenie boczne — kolejna kolumna po prawej stronie obszaru roboczego, o pełnej wysokości obszaru roboczego. Po zamknięciu kolumna znika całkowicie z przestrzeni roboczej. W stanie spoczynku mechanizm reprezentuje wyłącznie licznik pozycji magistrali w pasku kontekstu (warstwa 1).

**Warstwa widoczności i sposób wywołania.**

| Element | Warstwa | Sposób wywołania |
|---|---|---|
| Licznik pozycji magistrali | 1 | Widoczny bez interakcji w pasku kontekstu |
| Kolumna magistrali kontekstu | 3 | Kliknięcie licznika, pozycja menu hamburger (☰) środowiska albo skrót klawiszowy |
| Odłożenie i pobranie artefaktu | 3 | Menu kontekstowe artefaktu, przeciągnięcie |
| Reguły kierowania artefaktów | 4 | Okno konfiguracji, polecenie języka naturalnego w Chat Window, paleta poleceń `Ctrl/Cmd + K` |

**Zasięg.** Zasięg magistrali jest wartością ustawienia konfiguracyjnego i przyjmuje jeden z poziomów zasięgu opisanych w `architektura/izolacja-i-zaleznosci.md` (rozdz. 6): sesja, projekt albo środowisko. Przy zasięgu „sesja” magistrala obejmuje moduły jednej karty sesji; przy zasięgu „projekt” — wszystkie karty grupy projektu; przy zasięgu „środowisko” — wszystkie karty środowiska. Zasięg globalny nie występuje: artefakt nie przechodzi między środowiskami magistralą, lecz przez moduł Library.

**Relacja do modułu Library.** Magistrala jest przepływem roboczym o ograniczonej trwałości, moduł Library (`moduly/library.md`) jest repozytorium trwałym. Artefakt odłożony do magistrali pozostaje w niej przez okres wskazany ustawieniem retencji; przeniesienie do Library jest jawnym działaniem Użytkownika i jedyną drogą trwałego zachowania artefaktu oraz jedyną drogą jego przeniesienia między środowiskami.

**Relacja do punktów izolacji.** Magistrala respektuje reguły izolacji kontekstu obowiązujące parę modułów i zasięg bieżącej karty (`architektura/izolacja-i-zaleznosci.md`, rozdz. 4, 6, 9). Reguła izolacji kontekstu ustawiona dla pary modułów rozstrzyga, czy artefakt odłożony w module źródłowym jest widoczny w module docelowym; przekazanie niezgodne z regułą nie następuje, a magistrala prezentuje przy pozycji opis obowiązującej reguły wraz z odnośnikiem do okna konfiguracji punktów izolacji (rozdz. 10 tamże).

**Relacja do encji artefaktu w modelu danych.** Pozycja magistrali odwołuje się do wiersza `artefakt` (`architektura/model-danych.md`, rozdz. 16.1) przez jego identyfikator, zachowując moduł pochodzenia i kartę sesji odłożenia. Kolekcje i etykiety artefaktu (rozdz. 16.3–16.5) pozostają niezmienione — magistrala ich nie tworzy i nie modyfikuje.

---

## 7. Stany okien i elementów interfejsu

Poniższa tabela ustala wspólny, czterostanowy model interakcji (Domyślny, Hover, Active, Fokus) obowiązujący dla każdego interaktywnego elementu wymienionego w rozdziałach 3 i 5 — zgodnie z regułami stanów i fokusu systemu wizualnego marki — uzupełniony dwoma stanami procesu (Ładowanie, Błąd), z których żaden nie ogranicza klikalności kontrolki. Żaden interaktywny element nie przyjmuje stanu wyłączonego (disabled) jako blokady wykonania. Kolumna „Zastosowanie w przepływie” wskazuje, gdzie stan wybija się na pierwszy plan tego dokumentu.

| Stan | Reguła wizualna | Zastosowanie w przepływie |
|---|---|---|
| Domyślny (spoczynek) | Powierzchnia neutralna, bez akcentu złotego | Karty środowisk i kafle komponentów na stronie głównej przed interakcją (rozdz. 3.3–3.4, 5.5) |
| Hover (najechanie) | Subtelna zmiana tła albo uniesienie z akcentem złotym dla elementów akcentowanych | Karta środowiska, kafel komponentu, pozycja bocznej nawigacji, karta sesji (rozdz. 3.4–3.6) |
| Active / wybrany | Wciśnięcie (`translateY(1px)`) chwilowo; stan trwały — akcent złoty (podkreślenie, obrys, tło) | Pozycja modułu podświetlona w bocznej nawigacji (rozdz. 3.6); karta sesji aktywna w pasku kart (rozdz. 3.5) |
| Fokus (nawigacja klawiaturą) | Pierścień złoty 2 px + odsunięcie 2 px, wyłącznie `:focus-visible`, obowiązkowy dla każdej kontrolki interaktywnej | Wszystkie elementy formularza logowania (rozdz. 3.2); wszystkie przełączniki macierzy izolacji (rozdz. 5.6) |
| Ładowanie | Wskaźnik `.dn-spinner`, obrót ciągły, respektuje ograniczenie ruchu systemowego | Okno startowe (rozdz. 3.1); przycisk „Zaloguj się” w trakcie weryfikacji (rozdz. 3.2); przeładowanie modułu (rozdz. 3.6) |
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

### 8.1. Elementy poziomu „okno” (jednostki nawigacji pełnoekranowej)

| Okno / widok | Forma i waga | Warstwa | Sposób wywołania | Zamknięcie / opuszczenie | Trwałość stanu |
|---|---|---|---|---|---|
| Okno startowe | Pełnoekranowe, jednolite, bez nawigacji | 1 | Uruchomienie aplikacji | Automatyczne po uzgodnieniu połączenia | Brak — stan przejściowy |
| Okno rejestracji i logowania | Pełnoekranowy formularz | 1 | Brak ważnego tokenu lub faza budowy | Zatwierdzenie formularza albo „Pomiń” (skuteczny, gdy „Wymóg logowania” nieaktywne) | Wartości pól zachowane do zatwierdzenia; nie trwałe między uruchomieniami |
| Strona główna | Pełnoekranowa, trzy strefy | 1 | Token wydany; powrót z dowolnego miejsca platformy | Wybór strefy 1/2/3 | Nie przechowuje własnego stanu — jest zawsze tym samym widokiem |
| Okno środowiska | Pełnoekranowe, układ kolumnowy (nawigacja · Chat Window · kolumna dominująca) | 1 | Wybór karty środowiska (strefa 1) | Powrót do strony głównej | Karty sesji trwają w tle po opuszczeniu (rozdz. 6) |
| Chat Window | Lewa kolumna, stała, pełna wysokość obszaru roboczego | 1 | Obecne bez wywołania w każdym module i każdym środowisku | Nie podlega zamknięciu — okno stałe układu | Historia rozmowy trwała po stronie serwera niezależnie od stanu klienta (rozdz. 6) |
| Execution Loop Window | Kolumna sąsiadująca z Chat Window, pełna wysokość obszaru roboczego | 1 | Przyjęcie zlecenia do wykonania; przywrócenie jednym kliknięciem znacznika przebiegu albo poleceniem języka naturalnego w Chat Window | Zamknięcie kolumny — znika całkowicie z przestrzeni roboczej, pętla biegnie dalej po stronie serwera | Stan pętli (zlecenie, kolejka zadań, wyniki kontroli jakości) trwały po stronie serwera |
| Okno pomocnicze modułu | Kolumna boczna, otwierana jako rozszerzenie boczne | 2 | Kliknięcie pozycji zestawu okien modułu, skrót klawiszowy albo polecenie języka naturalnego w Chat Window | Zamknięcie — kolumna znika całkowicie, zwolniona szerokość przypada kolumnom pozostałym | Stan okna trwały po stronie serwera w obrębie karty sesji |
| Okno operacyjne | Kolumny obszaru roboczego karty sesji: Chat Window (lewa, stała), Execution Loop Window (kolumna sąsiadująca), okna modułu (kolumna dominująca i kolumny boczne) | 1 | Wybór modułu w bocznej nawigacji / panelu orkiestracji | Zmiana modułu, zamknięcie karty, powrót do strony głównej | Trwały po stronie serwera niezależnie od stanu klienta (rozdz. 6) |
| Okno budowy komponentu własnego | Kolumna dominująca obszaru roboczego albo nakładka modalna | 2 | Wybór kafla strefy 2, skrót klawiszowy albo polecenie języka naturalnego w Chat Window | Zapis komponentu albo zamknięcie bez zapisu | Wersja robocza — zależnie od implementacji; zapis trwały po zatwierdzeniu |
| Okno konfiguracji | Nakładka modalna nad bieżącym widokiem | 2 | Pozycja „Okno konfiguracji” (strefa 3 lub rozwinięcie kontekstowe `⋮`), skrót klawiszowy albo wyszukiwarka funkcji | Zamknięcie nakładki | Zmiany zapisywane natychmiast na warstwie wskazanej (globalna / sesji) |

### 8.2. Elementy poziomu „nawigacja” (stałe elementy ramy interfejsu)

| Element | Obecność | Forma i waga | Warstwa | Sposób wywołania | Zawartość |
|---|---|---|---|---|---|
| Chat Window | Każdy moduł i każde środowisko | Lewa kolumna, stała, pełna wysokość obszaru roboczego | 1 | Obecne bez wywołania | Historia rozmowy Użytkownik ↔ Wykonawca, pole poleceń języka naturalnego, strumień odpowiedzi i wyników, zatwierdzanie i przerywanie działań |
| Execution Loop Window | Każdy moduł i każde środowisko realizujące zlecenie | Kolumna sąsiadująca z Chat Window, pełna wysokość obszaru roboczego | 1 | Przyjęcie zlecenia do wykonania | Zlecenie i jego dekompozycja na zadania, kolejka i stan zadań, komunikaty sterujące Koordynator ↔ Wykonawca, wyniki kontroli jakości, wskaźniki przebiegu pętli, sterowanie przebiegiem |
| Pasek marki górny | Strona główna i każde środowisko | Wąski pasek pełnej szerokości | 1 | Widoczny bez interakcji | Godło (powrót do strony głównej), wyszukiwanie, wskaźnik konta |
| Boczna nawigacja modułów | TalkIn, WorkSpace, CodeStudio | Lewa kolumna nawigacyjna, stała | 1 | Widoczna bez interakcji | Lista modułów wg macierzy dostępności (8–9 pozycji) |
| Panel orkiestracji | Wyłącznie MultitaskingAI | Lewa kolumna nawigacyjna, stała, ta sama pozycja co boczna nawigacja | 1 | Widoczny bez interakcji | Sześć sekcji sterowania zespołem |
| Pasek kart sesji | Każde środowisko | Pasek kart w ramie okna środowiska | 1 | Widoczny bez interakcji | Karty otwartych przestrzeni roboczych + kontrolka „+” |
| Listwa ustawień | Strona główna (jawnie) + dostępna z każdego miejsca | Zwarta listwa, waga najniższa | 1 | Widoczna bez interakcji; równoważnie skrót klawiszowy i wyszukiwarka funkcji | Okno konfiguracji, Mobile, Always On Display |

### 8.3. Elementy poziomu „kontrolka” (pojedyncze akcje)

| Kontrolka | Forma i waga | Warstwa | Sposób wywołania | Stany istotne dla przepływu | Skutek |
|---|---|---|---|---|---|
| Karta środowiska | Duży panel interaktywny | 1 | Widoczna bez interakcji | spoczynek / hover / aktywny+fokus | Otwiera okno środowiska (nieodwracalnie zmienia widok — jedyna droga wejścia) |
| Kafel komponentu własnego | Średni panel interaktywny | 1 | Widoczny bez interakcji | spoczynek / hover / aktywny+fokus | Otwiera okno budowy komponentu (nie zmienia środowiska) |
| Pozycja listwy ustawień | Mały element paska | 1 | Widoczna bez interakcji | domyślny / hover / fokus | Otwiera nakładkę konfiguracji albo aktywuje funkcję globalną |
| Pozycja bocznej nawigacji / sekcji orkiestracji | Mały element listy | 1 | Widoczna bez interakcji | domyślny / hover / wybrana | Przeładowuje kolumnę dominującą bieżącej karty |
| Karta sesji | Element paska kart | 1 | Widoczna bez interakcji | aktywna / w tle / ze wskaźnikiem pracy | Przełącza widoczną zawartość kolumn obszaru roboczego |
| Pole poleceń Chat Window | Pole tekstowe wieloliniowe w lewej kolumnie | 1 | Widoczne bez interakcji | domyślny / fokus / wypełnione / ładowanie | Kieruje polecenie języka naturalnego do Wykonawcy; jest zarazem drogą wywołania funkcji warstwy 4 |
| Sterowanie przebiegiem pętli `Przebieg ▼` | Element zbiorczy w nagłówku Execution Loop Window | 3 | Rozwinięcie elementu zbiorczego albo polecenie języka naturalnego w Chat Window | zwinięty / rozwinięty | Wstrzymuje, wznawia, przerywa albo koryguje zlecenie pętli wykonawczej |
| Znacznik kontekstowy | Lekki znacznik w kolumnie nawigacji (`[Ubuntu]`, `[Fable 5]`) | 1 (znacznik) / 2 (selektor) | Kliknięcie znacznika otwiera selektor, który zwija się po wyborze | domyślny / hover / selektor otwarty | Zmienia środowisko, repozytorium, projekt, model albo wykonawcę bieżącej karty |
| Kontrolka „+” (nowa karta) | Mała ikonka | 1 | Widoczna bez interakcji | domyślny / hover / aktywny | Tworzy nową, pustą kartę sesji |
| Kontrolka „✕” (zamknij kartę) | Mała ikonka | 2 | Ujawniana najechaniem na kartę sesji | ukryta do najechania / hover | Zamyka widok karty |
| Przycisk „Pomiń” | Mały przycisk bez obrysu | 1 | Widoczny bez interakcji | obecny i klikalny w obu fazach wdrożenia | Wejście bez uwierzytelnienia, gdy „Wymóg logowania” nieaktywne; w przeciwnym razie komunikat o obowiązującym wymogu logowania |
| Kontrolka „⋮” (rozwinięcie kontekstowe) | Mała ikonka w nagłówku okna | 3 | Kliknięcie `⋮` w nagłówku okna | domyślny / hover / otwarte | Otwiera podzbiór ustawień warstwy sesji; po zamknięciu rozwinięcie znika całkowicie |
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
| Jawność i konfigurowalność zależności | Każde przełączenie środowiska lub modułu jest świadomym, widocznym przejściem (etapy 4, 6), nie ukrytą zmianą stanu; powiązania międzymodułowe pokazane wprost w ścieżkach 5.1–5.3 są możliwością, nie regułą wbudowaną | rozdz. 3.4, 3.6, 5.1–5.3 |
| Stopniowe ujawnianie funkcjonalności | Każdy element przepływu należy do jednej z czterech warstw widoczności; makiety rysowane są w stanie spoczynku, okna i panele warstw 2–3 znikają całkowicie z przestrzeni roboczej po zamknięciu, funkcje warstwy 4 uruchamiane są poleceniem języka naturalnego, skrótem klawiszowym, wyszukiwarką funkcji albo w trybie administracyjnym; każda ukryta funkcja pozostaje osiągalna jednym kliknięciem | rozdz. 3a |
| Dwa kanały komunikacji operacyjnej | Chat Window (Użytkownik ↔ Wykonawca) jest oknem stałym lewej kolumny w każdym module i każdym środowisku; Execution Loop Window (Koordynator ↔ Wykonawca) otwiera się w kolumnie sąsiadującej i prowadzi pętlę wykonawczą — oba są elementami pierwszoplanowymi każdego etapu przepływu | rozdz. 3.7, 3a, 6 |
| Układ pionowy (podział lewa–prawa) | Wszystkie okna robocze, okna komunikacji i okna pomocnicze rozmieszczone są w kolumnach sąsiadujących poziomo; okna pomocnicze otwierają się jako rozszerzenia boczne po prawej stronie obszaru roboczego, a regulacji podlega wyłącznie szerokość kolumn | rozdz. 3.5, 3.7, 5.1–5.4 |
| Pełna kompozycyjność | Liczba jednocześnie otwartych kart sesji i kombinacje modułów w poszczególnych kartach nie są ograniczone przez platformę (rozdz. 5.1: dwie karty TalkIn równolegle) | rozdz. 5.1, 6 |
| Orkiestracja pracy zespołu | Ścieżka MultitaskingAI (rozdz. 5.4) realizuje tę zasadę w warstwie nawigacyjnej — cztery okna robocze ról w kolumnach sąsiadujących, widoczne jednocześnie, wraz z pętlą pracy ciągłej 24/7/365 prowadzoną przez Koordynatora | rozdz. 5.4 |

---

## 10. Słowniczek pojęć

| Pojęcie | Definicja |
|---|---|
| Okno startowe | Pierwsze okno natywne prezentowane po uruchomieniu aplikacji klienckiej, obejmujące nawiązanie kanału WebSocket i wstępne uzgodnienie protokołu z serwerem, poprzedzające decyzję o wyświetleniu okna logowania albo bezpośrednim wejściu do platformy |
| Okno rejestracji i logowania | Okno uwierzytelniania Operatora; przycisk Pomiń obecny i klikalny w obu fazach wdrożenia — czy okno wymusza skuteczną rejestrację lub logowanie, zależy od ustawienia „Wymóg logowania” w oknie konfiguracji, konfigurowalnego przez Operatora niezależnie od fazy |
| Faza budowy | Faza wdrożenia, w której rdzeń działa lokalnie na komputerze Operatora; ustawienie „Wymóg logowania” domyślnie nieaktywne, zmienialne przez Operatora w oknie konfiguracji |
| Faza produkcyjna | Faza wdrożenia po migracji rdzenia na serwer wirtualny i otwarciu portu kanału WebSocket; ustawienie „Wymóg logowania” domyślnie aktywne, zmienialne przez Operatora w oknie konfiguracji |
| Token dostępu | Poświadczenie wydawane urządzeniu po uwierzytelnieniu, wykorzystywane do nawiązania lub potwierdzenia kolejnych połączeń WebSocket bez ponownego logowania |
| Strona główna (Centrum dowodzenia) | Przedpokój przed strefą roboczą platformy; trzy strefy — wybór środowiska, komponenty własne, ustawienia |
| Strefa 1 | Strefa strony głównej z czterema kartami wejścia do środowisk; jedyna droga wejścia do przestrzeni roboczej środowiska |
| Strefa 2 | Strefa strony głównej z czterema kaflami budowy komponentów własnych: Automations, Agents, Workspace, Assistant |
| Strefa 3 | Strefa strony głównej — listwa ustawień: okno konfiguracji, Mobile, Always On Display |
| Środowisko | Najwyższy poziom organizacji pracy: TalkIn, WorkSpace, CodeStudio, MultitaskingAI |
| Okno środowiska | Przestrzeń robocza otwierana po wyborze środowiska, w układzie pionowym (podział lewa–prawa): kolumna nawigacji modułów lub panelu orkiestracji, Chat Window, kolumna dominująca obszaru roboczego; pasek marki i pasek kart sesji stanowią ramę okna |
| Boczna nawigacja modułów | Stała, lewa kolumna nawigacyjna środowisk TalkIn, WorkSpace, CodeStudio, wyświetlająca listę dostępnych modułów |
| Panel orkiestracji | Odpowiednik bocznej nawigacji w środowisku MultitaskingAI, zajmujący tę samą lewą kolumnę nawigacyjną — sześć sekcji: Zespoły, Role, Kolejki, Orkiestracja, Harmonogram i automatyki, Monitor procesu |
| Moduł | Wyspecjalizowany obszar roboczy wewnątrz środowiska; piętnaście modułów, dostępność wg macierzy |
| Karta sesji | Pozioma karta w pasku kart okna środowiska; jedna, samodzielna przestrzeń robocza z własnym układem, historią i kontekstem |
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
| 3 | Trafia na stronę główną | Trzy strefy widoczne; brak zapamiętanej aktywnej karty (pierwsze uruchomienie) | Strona główna — rozdz. 3.3 |
| 4 | W strefie 2 klika kafel „Agents” | Otwiera się Agent Builder; konfiguruje model bazowy, tożsamość, skille, uprawnienia; zapisuje agenta | Strefa 2 — rozdz. 5.5 |
| 5 | Wraca na stronę główną, klika kartę „TalkIn” | Otwiera się okno środowiska TalkIn z nową, pustą kartą | Wybór środowiska — rozdz. 3.4 |
| 6 | W bocznej nawigacji wybiera „Studio” | Karta ładuje zestaw: Studio Editor, Tools Panel, Diff/Grep Panel, Session Repository, Preview Window, Chat Window | Okno operacyjne — rozdz. 3.7, 5.1 |
| 7 | Otwiera nową kartę („+”), wybiera „Research” | Druga karta z odrębnym kontekstem; karta Studio niezmieniona w tle | Karty równoległe — rozdz. 5.1, 6 |
| 8 | Wraca do strony głównej, klika kartę „CodeStudio” | Obie karty TalkIn trwają w tle; otwiera się nowe okno środowiska CodeStudio | Przełączanie środowiska — rozdz. 6 |
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
| Okno rejestracji i logowania | Okno startowe (brak ważnego tokenu lub faza budowy) | Strona główna |
| Strona główna | Okno startowe (token ważny) **albo** Okno rejestracji i logowania **albo** powrót z okna środowiska/operacyjnego **albo** zamknięcie okna strefy 2/3 | Okno środowiska (strefa 1) **albo** okno budowy komponentu (strefa 2) **albo** okno konfiguracji / aktywacja funkcji globalnej (strefa 3) |
| Okno środowiska | Strona główna (wybór karty środowiska) | Wybór modułu → okno operacyjne |
| Okno operacyjne | Wybór modułu w oknie środowiska | Zmiana modułu (ta sama karta) **albo** nowa karta **albo** powrót do strony głównej **albo** zamknięcie aplikacji (proces trwa na serwerze) |
| Okno budowy komponentu własnego (strefa 2) | Strona główna (kafel strefy 2) **albo** okno modułowe środowiska (dla Agents i Workspace) | Zapis → powrót do miejsca wejścia; wytwór wybieralny w oknie operacyjnym |
| Okno konfiguracji (strefa 3) | Strona główna **albo** dowolne okno środowiska/operacyjne (uproszczone menu kontekstowe) | Zamknięcie → powrót do dokładnie tego samego widoku sprzed otwarcia |

---

*Koniec dokumentu. Danaco Console — Przepływ okien: od uruchomienia aplikacji do okna roboczego, wersja 2.0.*

---
*Danaco Console — AI Workspace OS · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — [LICENSE](LICENSE). Kontakt: support@danaco-group.pl*
