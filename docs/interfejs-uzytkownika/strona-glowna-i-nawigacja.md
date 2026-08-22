# Danaco Console — Strona główna i nawigacja

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
| **Tytuł** | Strona główna i nawigacja |
| **Klasa dokumentu** | Specyfikacja docelowa |
| **Odbiorcy** | projektant · deweloper · redaktor treści interfejsu |
| **Przeznaczenie** | Ustala szczegółową specyfikację strony głównej platformy (centrum dowodzenia) oraz nawigacji między jej elementami — pionowego układu kolumnowego, trzech stref wyboru wraz z formą ich prezentacji, warstw widoczności, przepływu nawigacji strona główna → środowisko → moduł → okno, mechaniki kart sesji, bocznej nawigacji modułów i panelu orkiestracji — tak aby deweloper zbudował je bez rozstrzygania czegokolwiek samodzielnie |
| **Zakres** | układ trzech stref i forma ich prezentacji, warstwy widoczności w nawigacji, czterokrokowy przepływ nawigacji, mechanika kart sesji, przełączanie i przeładowanie przestrzeni roboczej, boczna nawigacja modułów, panel orkiestracji środowiska MultitaskingAI, wejście do ustawień i funkcji globalnych |
| **Poza zakresem** | szczegółowy opis samych środowisk i modułów jako bytów (Koncepcja platformy); rama okna aplikacji — [Rama okna aplikacji](rama-okna.md); okno konfiguracji punktów izolacji jako takie — poza niniejszym katalogiem |
| **Dokument nadrzędny** | [Elementy okien](elementy-okien.md) |
| **Dokumenty powiązane** | [Rama okna aplikacji](rama-okna.md) · [System wizualny](system-wizualny.md) · [Katalog komponentów](katalog-komponentow.md) · [Okno nowego projektu](okno-nowego-projektu.md) · [Okno historii sesji](okno-historii-sesji.md) |
| **Prototypy odniesienia** | `design/05-okna/` — prototypy platformowe niosące stronę główną i nawigację |
| **Źródła normatywne** | `design/zasoby/css/komponenty.css` (rodziny „Karta”, „Boczna nawigacja modułów”) · `design/zasoby/zetony/zetony.css` · `design/zasoby/ikony/manifest.json` (cztery ikony środowisk i ich motta) |
| **Zasada nadrzędna** | Strona główna jest centrum dowodzenia o wyraźnej hierarchii wagi, nie jednorodną siatką równoważnych elementów — forma prezentacji każdej z trzech stref (karty, kafle, listwa) odpowiada wprost jej wadze i rodzajowi akcji. |

Dokument stanowi szczegółową specyfikację strony głównej platformy Danaco Console (centrum dowodzenia) oraz nawigacji między jej elementami: pionowego układu kolumnowego strony głównej z Chat Window w lewej kolumnie, trzech stref wyboru wraz z formą ich prezentacji, warstw widoczności funkcjonalności, przepływu nawigacji strona główna → środowisko → moduł → okno, mechaniki kart sesji, mechanizmu przełączania i przeładowania przestrzeni roboczej, bocznej nawigacji modułów w środowiskach modułowych oraz panelu orkiestracji jako jej odpowiednika w środowisku MultitaskingAI. Dokument rozwija ustalenia zawarte w Koncepcji platformy (rozdziały 5, 7, 8, 9, 10 i 13) oraz w dokumencie Architektura (rozdziały 5, 7, 12 i 13).

---

## Spis treści

1. [Cel i miejsce dokumentu](#1-cel-i-miejsce-dokumentu)
2. [Strona główna jako centrum dowodzenia](#2-strona-główna-jako-centrum-dowodzenia)
   - [2.1 Chat Window na stronie głównej](#21-chat-window-na-stronie-głównej)
   - [2.2 Execution Loop Window](#22-execution-loop-window)
   - [2.3 Ścieżki działania z centrum dowodzenia](#23-ścieżki-działania-z-centrum-dowodzenia)
3. [Trzy strefy strony głównej](#3-trzy-strefy-strony-głównej)
   - [3.1 Przegląd i hierarchia wizualna](#31-przegląd-i-hierarchia-wizualna)
   - [3.2 Strefa 1 — karty środowisk](#32-strefa-1--karty-środowisk)
   - [3.3 Strefa 2 — kafle komponentów własnych](#33-strefa-2--kafle-komponentów-własnych)
   - [3.4 Strefa 3 — listwa ustawień](#34-strefa-3--listwa-ustawień)
   - [3.5 Schemat całościowy strony głównej](#35-schemat-całościowy-strony-głównej)
4. [3a. Warstwy widoczności w nawigacji](#3a-warstwy-widoczności-w-nawigacji)
   - [3 a.1. Zasada nadrzędna](#3a1-zasada-nadrzędna)
   - [3 a.2. Cztery warstwy widoczności](#3a2-cztery-warstwy-widoczności)
   - [3 a.3. Mechanizmy ukrywania funkcjonalności](#3a3-mechanizmy-ukrywania-funkcjonalności)
   - [3 a.4. Zasada jednego kliknięcia](#3a4-zasada-jednego-kliknięcia)
   - [3 a.5. Zasada rysowania makiet](#3a5-zasada-rysowania-makiet)
5. [Przepływ nawigacji: strona główna → środowisko → moduł → okno](#4-przepływ-nawigacji-strona-główna--środowisko--moduł--okno)
   - [4.1 Krok 1 — strona główna → środowisko](#41-krok-1--strona-główna--środowisko)
   - [4.2 Krok 2 — boczna nawigacja → moduł (środowiska modułowe)](#42-krok-2--boczna-nawigacja--moduł-środowiska-modułowe)
   - [4.3 Krok 2 (wariant) — panel orkiestracji (środowisko MultitaskingAI)](#43-krok-2-wariant--panel-orkiestracji-środowisko-multitaskingai)
   - [4.4 Krok 3 — moduł → okno kontekstowe](#44-krok-3--moduł--okno-kontekstowe)
   - [4.5 Krok 4 — karty sesji](#45-krok-4--karty-sesji)
   - [4.6 Schemat przepływu nawigacji](#46-schemat-przepływu-nawigacji)
6. [Boczna nawigacja modułów](#5-boczna-nawigacja-modułów)
   - [5.1 Mechanika i widoczność według macierzy dostępności modułów](#51-mechanika-i-widoczność-według-macierzy-dostępności-modułów)
   - [5.2 Stan wybrany i wskaźniki](#52-stan-wybrany-i-wskaźniki)
   - [5.3 Przykład: nawigacja w środowisku TalkIn](#53-przykład-nawigacja-w-środowisku-talkin)
7. [Panel orkiestracji — boczna nawigacja środowiska MultitaskingAI](#6-panel-orkiestracji--boczna-nawigacja-środowiska-multitaskingai)
   - [6.1 Różnica względem bocznej nawigacji modułowej](#61-różnica-względem-bocznej-nawigacji-modułowej)
   - [6.2 Sześć sekcji panelu](#62-sześć-sekcji-panelu)
   - [6.3 Konfigurowalność kolejności i widoczności](#63-konfigurowalność-kolejności-i-widoczności)
8. [Karty sesji — mechanika](#7-karty-sesji--mechanika)
   - [7.1 Anatomia karty sesji](#71-anatomia-karty-sesji)
   - [7.2 Otwieranie, zamykanie, przełączanie](#72-otwieranie-zamykanie-przełączanie)
   - [7.3 Domyślna izolacja karty sesji i możliwość współdzielenia](#73-domyślna-izolacja-karty-sesji-i-możliwość-współdzielenia)
   - [7.4 Karta sesji a rola w środowisku MultitaskingAI](#74-karta-sesji-a-rola-w-środowisku-multitaskingai)
9. [Przełączanie i przeładowanie przestrzeni roboczej](#8-przełączanie-i-przeładowanie-przestrzeni-roboczej)
   - [8.1 Przeładowanie w obrębie karty — zmiana modułu](#81-przeładowanie-w-obrębie-karty--zmiana-modułu)
   - [8.2 Przełączanie środowiska — powrót do centrum dowodzenia](#82-przełączanie-środowiska--powrót-do-centrum-dowodzenia)
   - [8.3 Trwałość sesji w tle](#83-trwałość-sesji-w-tle)
   - [8.4 Zestawienie: co się zmienia, a co jest zachowywane](#84-zestawienie-co-się-zmienia-a-co-jest-zachowywane)
10. [Wejście do ustawień i funkcji globalnych](#9-wejście-do-ustawień-i-funkcji-globalnych)
   - [9.1 Okno konfiguracji](#91-okno-konfiguracji)
   - [9.2 Mobile i Always On Display](#92-mobile-i-always-on-display)
11. [Zgodność z zasadami nadrzędnymi platformy](#10-zgodność-z-zasadami-nadrzędnymi-platformy)
12. [Słowniczek pojęć nawigacyjnych](#11-słowniczek-pojęć-nawigacyjnych)
13. [Wykazy normatywne i kryteria odbioru](#12-wykazy-normatywne-i-kryteria-odbioru)
   - [12.1 Etykiety interfejsu](#121-etykiety-interfejsu)
   - [12.2 Komunikaty](#122-komunikaty)
   - [12.3 Żetony projektowe](#123-żetony-projektowe)
   - [12.4 Komponenty interfejsu](#124-komponenty-interfejsu)
   - [12.5 Skróty klawiszowe](#125-skróty-klawiszowe)
   - [12.6 Stany kontrolek](#126-stany-kontrolek)
   - [12.7 Punkty łamania](#127-punkty-łamania)
   - [12.8 Kryteria odbioru](#128-kryteria-odbioru)
14. [Załącznik A. Macierz dostępności modułów w bocznej nawigacji](#załącznik-a-macierz-dostępności-modułów-w-bocznej-nawigacji)
15. [Załącznik B. Panel orkiestracji — zestawienie sekcji](#załącznik-b-panel-orkiestracji--zestawienie-sekcji)
16. [Załącznik C. Scenariusz nawigacyjny — pełny przebieg](#załącznik-c-scenariusz-nawigacyjny--pełny-przebieg)
17. [Załącznik D. Pełny wykaz komend kontraktu strony głównej i nawigacji](#załącznik-d-pełny-wykaz-komend-kontraktu-strony-głównej-i-nawigacji)
   - [Obszar `home` — 1 komenda](#obszar-home--1-komenda)
   - [Obszar `environment` — 2 komendy](#obszar-environment--2-komendy)
   - [Obszar `session` — 19 komend (rozdz. 7, 8)](#obszar-session--19-komend-rozdz-7-8)
   - [Obszar `window` — 7 komend (rozdz. 4, 8)](#obszar-window--7-komend-rozdz-4-8)

---

## 1. Cel i miejsce dokumentu

Niniejszy dokument rozwija ustalenia opracowania [Koncepcja platformy](../architektura/koncepcja-platformy.md)
do poziomu szczegółowości warstwy nawigacyjnej interfejsu.

**Konwencja przywołań obowiązująca w całym dokumencie.** Trzy opracowania są przywoływane
w dalszej części skróconą nazwą, bez powtarzania odsyłacza przy każdym wystąpieniu. Nazwa
skrócona wskazuje zawsze ten sam plik:

| Nazwa użyta w treści | Przywoływane opracowanie |
|---|---|
| Koncepcja platformy | [Koncepcja platformy](../architektura/koncepcja-platformy.md) |
| Architektura | [Architektura techniczna](../architektura/architektura.md) |
| System wizualny | [System wizualny](system-wizualny.md) |

**Ustalenia źródłowe rozwijane w tym dokumencie:**

| Ustalenie źródłowe | Miejsce w Koncepcji platformy | Rozwinięcie w niniejszym dokumencie |
|---|---|---|
| Podział strony głównej na trzy strefy | rozdz. 5 | Rozdział 3 — opis każdej strefy wraz z formą prezentacji |
| Forma prezentacji trzech stref | rozdz. 5.4 | Rozdział 3 — karty środowisk, kafle komponentów, listwa ustawień |
| Chat Window jako centralny punkt pracy użytkownika | rozdz. 2.4 | Rozdział 2 — lewa kolumna strony głównej i każdego modułu |
| Execution Loop Window — pętla wykonawcza Koordynator ↔ Wykonawca | rozdz. 2.4, 11.2 | Rozdział 2 — kolumna sąsiadująca, otwierana ze strony głównej i z każdego modułu |
| Warstwy widoczności funkcjonalności | rozdz. 12 | Rozdział 3a — cztery warstwy, mechanizmy ukrywania, zasada jednego kliknięcia |
| Czterokrokowy przepływ nawigacji | rozdz. 6 | Rozdział 4 — każdy z czterech kroków z osobna |
| Mechanika kart sesji | rozdz. 6 (krok 4) | Rozdział 7 — anatomia, cykl życia i schemat karty |
| Boczna nawigacja modułów | rozdz. 8 | Rozdział 5 |
| Panel orkiestracji środowiska MultitaskingAI | rozdz. 13.8 | Rozdział 6 |

**Rozgraniczenie zakresu — co dokument opisuje, a co pozostaje w źródle:**

| W zakresie niniejszego dokumentu | Poza zakresem — opisane w całości w źródle |
|---|---|
| Układ trzech stref i forma ich prezentacji | Szczegółowy opis środowisk (Koncepcja platformy, rozdz. 9) |
| Przepływ nawigacji na czterech krokach | Szczegółowy opis modułów (Koncepcja platformy, rozdz. 11) |
| Mechanika kart sesji | Funkcje globalne (Koncepcja platformy, rozdz. 12) |
| Przeładowanie i przełączanie przestrzeni roboczej | Szczegółowy interfejs okna konfiguracji punktów izolacji (Koncepcja platformy, rozdz. 6.4–6.6) |
| Boczna nawigacja modułów i panel orkiestracji | — |

Elementy z prawej kolumny są w niniejszym dokumencie przywoływane wyłącznie w zakresie potrzebnym do opisania nawigacji do nich prowadzącej.

---

## 2. Strona główna jako centrum dowodzenia

Strona główna jest centrum dowodzenia platformy. Jej układ jest wyłącznie pionowy — podział lewa–prawa. Lewa kolumna, stała, o pełnej wysokości obszaru roboczego, zawiera **Chat Window** — główne okno komunikacji Użytkownik ↔ Wykonawca. Kolumna dominująca po prawej stronie zawiera trzy strefy wyboru: karty środowisk, kafle komponentów własnych i listwę ustawień.

### 2.1. Chat Window na stronie głównej

Chat Window jest obecne już na stronie głównej, w tym samym miejscu układu, w którym występuje w każdym środowisku i w każdym module — w lewej kolumnie obszaru roboczego. Stanowi centralny punkt pracy użytkownika i podstawowy mechanizm sterowania wszystkimi procesami realizowanymi przez platformę.

| Cecha Chat Window na stronie głównej | Treść |
|---|---|
| Kanał | Użytkownik ↔ Wykonawca (AI / agent / system wykonawczy) |
| Położenie | Lewa kolumna, stała, pełna wysokość obszaru roboczego |
| Rola | Centralny punkt pracy użytkownika; podstawowy mechanizm sterowania procesami platformy |
| Zakres działania | Przyjmowanie poleceń w języku naturalnym, strumień odpowiedzi i wyników, zatwierdzanie i przerywanie działań, wyjaśnianie wyniku i kontekstu |
| Sterowanie nawigacją | Polecenie języka naturalnego otwiera środowisko, moduł, okno konfiguracji i funkcję globalną bez korzystania ze stref wyboru |
| Ciągłość | Przejście ze strony głównej do środowiska nie zmienia położenia okna — zmienia się wyłącznie zawartość kolumny prawej |
| Warstwa widoczności | Warstwa 1 — zawsze widoczna |

### 2.2. Execution Loop Window

**Execution Loop Window** jest oknem pętli wykonawczej prezentującym komunikację Koordynator ↔ Wykonawca. Otwiera się w kolumnie sąsiadującej z Chat Window i jest dostępne z poziomu strony głównej oraz z poziomu każdego modułu.

| Cecha Execution Loop Window | Treść |
|---|---|
| Kanał | Koordynator ↔ Wykonawca |
| Położenie | Kolumna sąsiadująca z Chat Window |
| Dostępność | Strona główna oraz każdy moduł każdego środowiska |
| Zawartość | Bieżące zlecenie i jego dekompozycja na zadania, kolejka i stan zadań, wymiana komunikatów sterujących, wyniki kontroli jakości i decyzje o ponowieniu, wskaźniki przebiegu pętli |
| Sterowanie przebiegiem | Wstrzymanie, wznowienie, przerwanie, korekta zlecenia |
| Warstwa widoczności | Warstwa 2 — otwierana przełącznikiem `Pętla ▼` w pasku znaczników kontekstowych; kolumna znika po zamknięciu |

Role rozdzielone są jednoznacznie: **Użytkownik** zleca i zatwierdza, **Koordynator** dekomponuje zlecenie oraz przydziela i nadzoruje zadania, **Wykonawca** realizuje zadania.

### 2.3. Ścieżki działania z centrum dowodzenia

```
Użytkownik
    │
    ▼
STRONA GŁÓWNA — Centrum dowodzenia
    │
    ├─▶ Chat Window (lewa kolumna) ─────▶ sterowanie procesami poleceniem języka naturalnego
    ├─▶ Execution Loop Window ──────────▶ pętla wykonawcza Koordynator ↔ Wykonawca
    ├─▶ Strefa 1 · wybór środowiska ────▶ przestrzeń robocza środowiska (rozdz. 9 Koncepcji)
    ├─▶ Strefa 2 · komponent własny ────▶ okno konfiguracji komponentu
    └─▶ Strefa 3 · ustawienia ──────────▶ okno konfiguracji lub funkcja globalna
```

| Ścieżka | Miejsce układu | Akcja | Rezultat | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|
| Sterowanie procesami | Lewa kolumna | Polecenie w Chat Window | Uruchomienie, zatwierdzenie albo przerwanie procesu platformy | 1 | Widoczne bez interakcji |
| Nadzór nad pętlą wykonawczą | Kolumna sąsiadująca | Otwarcie Execution Loop Window | Podgląd i sterowanie przebiegiem pętli Koordynator ↔ Wykonawca | 2 | Przełącznik `Pętla ▼`; skrót klawiszowy |
| Wejście do środowiska | Prawa kolumna, strefa 1 | Kliknięcie karty środowiska | Otwarcie przedsionka jednego z czterech środowisk — widoku wejściowego z kaflami modułów. Przestrzeń robocza otwiera się dopiero po wyborze modułu | 1 | Widoczne bez interakcji |
| Utworzenie lub skonfigurowanie komponentu własnego | Prawa kolumna, strefa 2 | Kliknięcie kafla komponentu | Otwarcie okna konfiguracji komponentu własnego | 1 | Widoczne bez interakcji |
| Ustawienia lub funkcja globalna | Prawa kolumna, strefa 3 | Kliknięcie pozycji listwy | Okno konfiguracji platformy albo aktywacja jednej z dwóch funkcji globalnych | 1 | Widoczne bez interakcji |

Centrum dowodzenia nie jest jednorodną siatką równoważnych elementów, lecz przestrzenią o wyraźnej hierarchii wagi (Koncepcja platformy, rozdz. 7.4).

| Element | Funkcja w centrum dowodzenia | Waga wizualna | Warstwa | Sposób wywołania |
|---|---|---|---|---|
| Chat Window | Centralny punkt pracy i sterowania procesami | Lewa kolumna, stała, pełna wysokość obszaru roboczego | 1 | Widoczne bez interakcji |
| Execution Loop Window | Nadzór nad pętlą wykonawczą | Kolumna sąsiadująca, otwierana | 2 | Przełącznik `Pętla ▼`; skrót klawiszowy; polecenie języka naturalnego |
| Strefa 1 | Główny punkt wejścia do pracy | Główna | 1 | Widoczne bez interakcji |
| Strefa 2 | Zaplecze twórcze do budowy komponentów własnych | Pośrednia | 1 | Widoczne bez interakcji |
| Strefa 3 | Warstwa narzędziowa ustawień | Najniższa | 1 | Widoczne bez interakcji |

**Makieta strony głównej w stanie spoczynku:**

```
 ══════════════════════════════════════════════════════════════════════════════════════
  Chat Window                │  DANACO CONSOLE — CENTRUM DOWODZENIA
  Użytkownik ↔ Wykonawca     │
                             │  STREFA 1 · KARTY ŚRODOWISK
  ─────────────────────────  │  ┌────────┐ ┌──────────┐ ┌───────────┐ ┌───────────────┐
  │ Polecenie…            │  │  │ TalkIn │ │WorkSpace │ │CodeStudio │ │MultitaskingAI │
  ─────────────────────────  │  └────────┘ └──────────┘ └───────────┘ └───────────────┘
                             │
  [Danaco Console] [Ubuntu]  │  STREFA 2 · KAFLE KOMPONENTÓW WŁASNYCH
  [Worktree] [Fable 5]       │  ┌───────────┐ ┌──────┐ ┌─────────┐ ┌─────────┐
  [Ultra]  Pętla ▼      ⋮    │  │Automations│ │Agents│ │Workspace│ │Assistant│
                             │  └───────────┘ └──────┘ └─────────┘ └─────────┘
                             │
                             │  STREFA 3 · LISTWA USTAWIEŃ
                             │  [ Okno konfiguracji ] [ Mobile ] [ Always On Display ]
 ══════════════════════════════════════════════════════════════════════════════════════
   Widoczna wyłącznie warstwa 1 oraz zwinięte wyzwalacze warstw 2–3:
   znaczniki kontekstowe, `Pętla ▼`, `⋮`
```

Otwarcie Execution Loop Window wstawia jego kolumnę bezpośrednio po prawej stronie Chat Window; strefy wyboru przesuwają się w prawo, zachowując pionowy podział lewa–prawa. Regulacji podlega wyłącznie szerokość kolumn.

Hierarchia wagi determinuje formę prezentacji każdej ze stref, opisaną szczegółowo w rozdziale 3. Strona główna jest jednocześnie punktem powrotu: użytkownik znajdujący się w dowolnym środowisku wraca do centrum dowodzenia, aby wybrać inne środowisko, skonfigurować kolejny komponent własny lub otworzyć ustawienia. Mechanizm powrotu i jego skutki dla otwartych kart sesji opisuje rozdział 8.2.

---

## 3. Trzy strefy strony głównej

### 3.1. Przegląd i hierarchia wizualna

Trzy strefy zajmują prawą, dominującą kolumnę strony głównej. Każda ma własną formę prezentacji, dobraną do swojej wagi i funkcji.

| Strefa | Zawartość | Forma prezentacji | Waga wizualna | Rodzaj akcji | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|
| Strefa 1 — wybór środowiska | TalkIn, WorkSpace, CodeStudio, MultitaskingAI | Karty środowisk — cztery duże karty wejścia | Główna (środek ciężkości) | Wejście do przestrzeni roboczej | 1 | Widoczna bez interakcji |
| Strefa 2 — komponenty własne | Automations, Agents, Workspace, Assistant | Kafle komponentów własnych — siatka czterech kafli | Pośrednia | Utworzenie komponentu | 1 | Widoczna bez interakcji |
| Strefa 3 — ustawienia | Okno konfiguracji, Mobile, Always On Display | Listwa ustawień — zwarty pasek narzędziowy | Najniższa | Ustawienie | 1 | Widoczna bez interakcji |

Forma prezentacji opiera się na czterech przesłankach:

| Przesłanka | Treść |
|---|---|
| Metafora centrum dowodzenia | Z natury ma hierarchię, a nie jednorodną siatkę |
| Zasady [systemu wizualnego](system-wizualny.md) | Jeden błękit sygnałowy jako jedyny akcent chromatyczny poza szarością (`--dn-sygnal-*`, [System wizualny](system-wizualny.md) rozdz. 2.2); krój nagłówkowy `--dn-ff-naglowek` (Space Grotesk) zarezerwowany dla tytułów |
| Zgodność formy z rodzajem akcji | Wejście, utworzenie, ustawienie |
| Czytelny rytm strony | Dwie siatki po cztery elementy o różnej skali, dopełnione listwą |

Trzy formy — karty, kafle, listwa — są uporządkowane malejąco pod względem skali i wizualnej masy, co samo w sobie komunikuje hierarchię ważności: użytkownik odczytuje, która strefa jest głównym celem wizyty, zanim przeczyta którykolwiek napis.

### 3.2. Strefa 1 — karty środowisk

Strefa 1 jest punktem wejścia do czterech środowisk platformy: TalkIn, WorkSpace, CodeStudio i MultitaskingAI (Koncepcja platformy, rozdz. 7.1, rozdz. 9). Przyjmuje formę czterech dużych kart wejścia, po jednej na środowisko, stanowiących wizualny środek ciężkości strony głównej.

**Anatomia pojedynczej karty środowiska**, zweryfikowana wobec klasy `.dn-karta-srodowiska` i jej
czterech elementów pochodnych w [Systemie wizualnym](system-wizualny.md) rozdz. 14.11:

| Element karty | Klasa | Zawartość | Uwagi | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|
| Godło środowiska | `.dn-karta-srodowiska-godlo` | Symbol graficzny właściwy danemu środowisku — jedna z czterech ikon własnych katalogu (`srodowisko-talkin`, `srodowisko-workspace`, `srodowisko-codestudio`, `srodowisko-multitaskingai`, [System wizualny](system-wizualny.md) Załącznik C) | wymiar zapisany w regule `.dn-karta-srodowiska-godlo` arkusza `design/zasoby/css/komponenty.css`, kolor `--dn-tekst`, odstęp dolny `--dn-od-4` | 1 | Widoczne bez interakcji |
| Tytuł | `.dn-karta-srodowiska-tytul` | Nazwa środowiska (TalkIn, WorkSpace, CodeStudio, MultitaskingAI) | krój `--dn-ff-naglowek`, stopień `--dn-fs-3xl`, waga `--dn-fw-polgruba` | 1 | Widoczny bez interakcji |
| Motto | `.dn-karta-srodowiska-motto` | Krótkie hasło środowiska — trzy czasowniki rozkazujące (tabela mott w dalszej części rozdz. 3.2) | krój `--dn-ff-mono`, stopień `--dn-fs-sm`, kolor `--dn-tekst-3` | 1 | Widoczne bez interakcji |
| Opis trybu pracy | `.dn-karta-srodowiska-opis` | Jednozdaniowy opis, spójny z definicją środowiska (Koncepcja platformy, rozdz. 9.1–9.4) | krój bazowy `--dn-ff-bazowa`, stopień `--dn-fs-base`, kolor `--dn-tekst-2`, długość opisu do czterdziestu czterech znaków | 1 | Widoczny bez interakcji |
| Menu karty `⋮` | — | Otwarcie środowiska w nowej karcie sesji, przypięcie środowiska, wyczyszczenie kart sesji środowiska | Zestaw akcji karty | 3 | Menu kebab `⋮` na karcie |
| Profil startowy środowiska | — | Wybór modelu, wykonawcy i trybu pracy stosowanego przy wejściu do środowiska | Zwija się po użyciu | 2 | Kliknięcie znacznika kontekstowego w pasku kontekstu |

Cztery motta, zweryfikowane w [Systemie wizualnym](system-wizualny.md) Załącznik C (pole
`zastosowanie` czterech ikon środowisk) — każde z dokładnie trzech czasowników rozkazujących:

| Środowisko | Motto (`.dn-karta-srodowiska-motto`) |
|---|---|
| TalkIn | „Myśl. Analizuj. Rozumiej.” |
| WorkSpace | „Planuj. Organizuj. Realizuj.” |
| CodeStudio | „Projektuj. Buduj. Rozwijaj.” |
| MultitaskingAI | „Deleguj. Koordynuj. Nadzoruj.” |

**Stany karty środowiska**, zweryfikowane wobec reguł CSS `.dn-karta-srodowiska` — komentarz
źródłowy nazywa wprost: „Spoczynek bez sygnału; najechanie/aktywność: wstęga górna 2 px”:

| Stan | Selektor | Wygląd | Znaczenie |
|---|---|---|---|
| Spoczynkowy | `.dn-karta-srodowiska` | obrys `--dn-obrys`, tło `--dn-powierzchnia`, brak wstęgi (`::before` z `opacity: 0`) | karta neutralna, bez akcentu sygnałowego |
| Wskazanie kursorem / aktywny | `.dn-karta-srodowiska:hover`, `.dn-karta-srodowiska[aria-current='true']` | obrys → `--dn-obrys-mocny`, cień → `--dn-cien-2`, `transform: translateY(-2px)` (uniesienie 2px), wstęga górna grubości `--dn-wym-wstega` (2px) w kolorze `--dn-kropka` pojawia się (`opacity: 1`) | sygnalizuje gotowość do wejścia; ten sam wzorzec wstęgi co pozycja bieżąca w [Systemie wizualnym](system-wizualny.md) rozdz. 14.11 |
| Ognisko klawiatury | pierścień `--dn-fokus` wspólny wszystkim kontrolkom interaktywnym | — | nawigacja klawiaturą |

**Makieta pojedynczej karty środowiska:**

```
  ┌──────────────────────────────┐
  │╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍│  ← wstęga górna 2px, `--dn-kropka` (wyłącznie wskazanie kursorem/aktywny)
  │            ◈  godło          │  ← .dn-karta-srodowiska-godlo, 40×40px
  │                              │
  │          T a l k I n         │  ← .dn-karta-srodowiska-tytul, Space Grotesk 30px
  │      Myśl. Analizuj.         │  ← .dn-karta-srodowiska-motto, IBM Plex Mono
  │        Rozumiej.             │
  │    Wiedza, komunikacja i     │  ← .dn-karta-srodowiska-opis, do 44 znaków
  │    praca z treścią           │
  │                              │
  └──────────────────────────────┘
    Spoczynek:  powierzchnia neutralna, brak wstęgi
    Wskazanie kursorem / aktywna:  wstęga sygnałowa 2px + uniesienie 2px + pierścień fokusu
```

Wstęga sygnałowa pozostaje zarezerwowana wyłącznie dla stanu aktywnego i wskazania kursorem, zgodnie z
zasadą „stan nigdy samym kolorem” ([System wizualny](system-wizualny.md) rozdz. 1.2) — karta niesie
zarówno zmianę koloru wstęgi, jak i zmianę kształtu (uniesienie, cień), nie samą barwę. Afordancja
karty odpowiada charakterowi przypisanej jej akcji: „wejdź do przestrzeni roboczej” — kliknięcie
karty jest równoznaczne z pierwszym krokiem przepływu nawigacji opisanego w rozdziale 4.1.

**Kolejność kart i opis trybu pracy na karcie** (zgodnie z Koncepcją platformy, rozdz. 5.1, rozdz. 9 — od środowiska pracy z treścią i komunikacją, przez produktywność projektową i programowanie, do orkiestracji autonomicznej pracy ciągłej):

| Kolejność | Środowisko | Opis trybu pracy na karcie (Koncepcja platformy, rozdz. 9.1–9.4) |
|---|---|---|
| 1 | TalkIn | Wiedza, komunikacja i praca z treścią |
| 2 | WorkSpace | Produktywność, organizacja i realizacja projektów |
| 3 | CodeStudio | Programowanie |
| 4 | MultitaskingAI | Orkiestracja autonomicznej pracy ciągłej |

### 3.3. Strefa 2 — kafle komponentów własnych

Strefa 2 jest punktem wejścia do konfiguracji czterech komponentów własnych: Automations, Agents, Workspace i Assistant (Koncepcja platformy, rozdz. 7.2, rozdz. 2.3). Przyjmuje formę siatki czterech kafli, wizualnie lżejszych od kart strefy 1 — mniejsza skala, brak godła na poziomie karty środowiska, brak zarezerwowanego dla tytułów kroju nagłówkowego.

**Anatomia pojedynczego kafla komponentu własnego:**

| Element kafla | Zawartość | Uwagi | Warstwa | Sposób wywołania |
|---|---|---|---|---|
| Ikona komponentu | Symbol graficzny właściwy rodzajowi komponentu (automatyka, agent, projekt, profil asystenta) | Mniejsza skala niż godło środowiska | 1 | Widoczna bez interakcji |
| Nazwa komponentu | Automations, Agents, Workspace, Assistant | Krój tekstowy, nie nagłówkowy | 1 | Widoczna bez interakcji |
| Etykieta działania | Zorientowana na tworzenie | Odróżnia kafel od karty wejścia do środowiska | 1 | Widoczna bez interakcji |
| Lista komponentów zapisanych | Wcześniej utworzone automatyki, agenci, projekty i profile asystenta | Panel wysuwany, znikający po zamknięciu | 3 | Panel popover kafla |
| Menu `Operacje ▼` | Duplikowanie, eksport, import i usunięcie komponentu własnego | Grupowanie logiczne akcji | 3 | Rozwinięcie `Operacje ▼` w panelu kafla |
| Edycja definicji komponentu w formacie źródłowym | Bezpośrednia praca na zapisie konfiguracji komponentu | Funkcja ekspercka | 4 | Polecenie języka naturalnego w Chat Window; wyszukiwarka funkcji; tryb administracyjny |

**Cztery kafle — nazwa, etykieta działania, powstający wytwór** (kolejność wg Koncepcji platformy, rozdz. 7.2):

| Kolejność | Komponent własny | Etykieta działania | Powstający wytwór |
|---|---|---|---|
| 1 | Automations | „Zbuduj automatykę” | Automatyka |
| 2 | Agents | „Skonfiguruj agenta” | Agent |
| 3 | Workspace | „Załóż projekt” | Projekt |
| 4 | Assistant | „Ustaw profil asystenta” | Profil asystenta |

**Makieta pojedynczego kafla komponentu własnego:**

```
  ┌─────────────────┐
  │     ▧  ikona    │
  │   Automations   │  ← krój tekstowy (nie nagłówkowy)
  │  „Zbuduj        │  ← etykieta zorientowana na tworzenie
  │   automatykę”   │
  └─────────────────┘
```

Odmienna forma wobec strefy 1 oddziela w odbiorze „wejście do środowiska” od „zbudowania komponentu” — kliknięcie kafla nie przenosi użytkownika do przestrzeni roboczej środowiska, lecz otwiera okno konfiguracji właściwe danemu rodzajowi komponentu, w którym powstaje nazwany wytwór użytkownika (Koncepcja platformy, rozdz. 2.3, rozdz. 7.2). Utworzony komponent trafia do pamięci aplikacji jako zasób własny i pozostaje wybieralny później, w module, w którym ma zastosowanie — zgodnie z Załącznikiem D.1 Koncepcji platformy (szablon komponentu własnego).

**Weryfikacja formy wobec klasy `.dn-kafel`.** [System wizualny](system-wizualny.md) rozdz. 14
potwierdza kontrast opisany powyżej wprost w wartościach żetonów: kafel komponentu własnego
(`.dn-kafel`) niesie tło `--dn-panel`, krój bazowy `--dn-ff-bazowa`, zaokrąglenie `--dn-r-lg` —
podczas gdy karta środowiska (`.dn-karta-srodowiska`, rozdz. 3.2) niesie tło `--dn-powierzchnia`,
krój nagłówkowy `--dn-ff-naglowek` (Space Grotesk) w stopniu `--dn-fs-3xl` i zaokrąglenie większe
`--dn-r-xl`. Jedyny stan interaktywny kafla — wskazanie kursorem (`.dn-kafel:hover`) — zmienia obrys
na `--dn-obrys-mocny` i tło na `--dn-powierzchnia`, bez uniesienia (`transform`) ani wstęgi górnej
właściwej karcie środowiska (rozdz. 3.2) — kafel jest komponentem wizualnie spokojniejszym,
zgodnie z jego wagą pośrednią w hierarchii trzech stref (rozdz. 3.1).

### 3.4. Strefa 3 — listwa ustawień

Strefa 3 udostępnia okno konfiguracji i ustawień platformy oraz dwie funkcje globalne: Mobile i Always On Display (Koncepcja platformy, rozdz. 7.3, rozdz. 12). Przyjmuje formę zwartego paska narzędziowego — listwy — o najniższej wadze wizualnej spośród trzech stref, sygnalizującej warstwę stale dostępną, a nie punkt docelowy odwiedzin.

**Zawartość listwy ustawień:**

| Pozycja | Prowadzi do | Charakter | Warstwa | Sposób wywołania |
|---|---|---|---|---|
| Okno konfiguracji | Pełny zakres ustawień platformy, w tym okno konfiguracji punktów izolacji (Koncepcja platformy, rozdz. 6.4–6.6; Architektura, rozdz. 13) | Ustawienie | 1 | Widoczna bez interakcji |
| Mobile | Globalny tryb dostępu do platformy z urządzeń mobilnych (Koncepcja platformy, rozdz. 12.1; `funkcje-globalne/mobile.md`) | Funkcja globalna | 1 | Widoczna bez interakcji |
| Always On Display | Globalny agent towarzyszący (Koncepcja platformy, rozdz. 12.2) | Funkcja globalna | 1 | Widoczna bez interakcji |
| Ustawienia szybkie listwy | Wygląd interfejsu, język, gęstość układu | Ustawienie | 3 | Menu kebab `⋮` listwy |
| Narzędzia diagnostyczne platformy | Dzienniki procesów, podgląd stanu sesji, diagnostyka połączeń modeli | Funkcja ekspercka | 4 | Skrót klawiszowy; wyszukiwarka funkcji; tryb administracyjny |

**Makieta listwy ustawień:**

```
  ┌──────────────────────────────────────────────────────────────┐
  │  [ Okno konfiguracji ]    [ Mobile ]    [ Always On Display ]│
  └──────────────────────────────────────────────────────────────┘
    zwarty pasek narzędziowy · najniższa waga wizualna · stale dostępny
```

Pozycje listwy nie mają formy kart ani kafli — są zwartymi elementami paska, spójnymi wizualnie z resztą interfejsu narzędziowego platformy, bez własnego godła ani rozbudowanego opisu. Kliknięcie pozycji „Okno konfiguracji” otwiera okno opisane w rozdziale 9.1. Kliknięcie pozycji Mobile lub Always On Display aktywuje odpowiednią funkcję globalną — obie pozostają dostępne również z poziomu każdego środowiska i modułu, ponieważ funkcje globalne nie tworzą własnej przestrzeni roboczej i nie podlegają mechanizmowi przełączania i przeładowania (Koncepcja platformy, rozdz. 5.3, rozdz. 12). Ich obecność w listwie strony głównej jest zatem jednym z kilku równoważnych punktów dostępu, a nie jedynym.

### 3.5. Schemat całościowy strony głównej

```
 ═══════════════════════════════════════════════════════════════════════════
  Chat Window            │ Execution Loop     │ DANACO CONSOLE — CENTRUM
  Użytkownik ↔ Wykonawca │ Koordynator ↔      │ DOWODZENIA
                         │ Wykonawca          │
                         │ (kolumna otwarta)  │ STREFA 1 · KARTY ŚRODOWISK
                         │                    │ ┌──────┐┌────────┐┌────────┐
  ───────────────────    │ Zlecenie           │ │TalkIn││WorkSpa.││CodeSt. │
  │ Polecenie…      │    │ Zadania · stan     │ └──────┘└────────┘└────────┘
  ───────────────────    │ Kontrola jakości   │ ┌──────────────┐
                         │ Sterowanie pętlą   │ │MultitaskingAI│
  [Danaco Console]       │                    │ └──────────────┘
  [Ubuntu] [Worktree]    │                    │
  [Fable 5] [Ultra]      │                    │ STREFA 2 · KAFLE
  Pętla ▼            ⋮   │                    │ ┌─────────┐┌──────┐┌───────┐
                         │                    │ │Automat. ││Agents││Worksp.│
                         │                    │ └─────────┘└──────┘└───────┘
                         │                    │ ┌─────────┐
                         │                    │ │Assistant│
                         │                    │ └─────────┘
                         │                    │
                         │                    │ STREFA 3 · LISTWA USTAWIEŃ
                         │                    │ [Okno konfiguracji] [Mobile]
                         │                    │ [Always On Display]        ⋮
 ═══════════════════════════════════════════════════════════════════════════
```

Powyższy schemat porządkuje układ opisany w rozdziałach 2 i 3.2–3.4: lewa kolumna z Chat Window, sąsiadująca kolumna Execution Loop Window oraz prawa kolumna dominująca z dwiema siatkami po cztery elementy o różnej skali (strefa 1, strefa 2), dopełnionymi listwą (strefa 3). W stanie spoczynku kolumna Execution Loop Window jest zwinięta, a jej wyzwalaczem pozostaje przełącznik `Pętla ▼`. Wymiary, odstępy oraz tokeny kolorów i typografii określa system projektowy marki — przewodnik marki i pliki tokenów (Koncepcja platformy, rozdz. 7.4).

---

## 3a. Warstwy widoczności w nawigacji

### 3a.1. Zasada nadrzędna

Interfejs Danaco Console realizuje zasadę: **jeżeli funkcja nie jest potrzebna do realizacji aktualnego zadania, nie jest widoczna**. Interfejs ujawnia możliwości systemu stopniowo — zależnie od kontekstu, roli użytkownika i wykonywanej czynności — zachowując maksymalną moc funkcjonalną przy minimalnej złożoności wizualnej. Złożoność platformy istnieje w architekturze i pozostaje niewidoczna w interfejsie do chwili wystąpienia potrzeby użycia danej funkcji. Liczba modułów, agentów, przepływów pracy, komponentów, narzędzi, paneli, ustawień i funkcji administracyjnych nie wpływa na postrzeganą prostotę interfejsu.

Strona główna i nawigacja prezentują wyłącznie warstwę 1: główne okno komunikacji, aktywne okno wiodące, kontekst pracy, podstawową nawigację oraz wskaźniki stanu wykonania. Warstwa ta zajmuje ponad 80% powierzchni interfejsu.

### 3a.2. Cztery warstwy widoczności

| Warstwa | Nazwa | Zawartość w warstwie nawigacyjnej | Sposób wywołania |
|---|---|---|---|
| 1 | Zawsze widoczna | Chat Window, aktywne okno wiodące, karty środowisk, kafle komponentów, listwa ustawień, boczna nawigacja modułów, panel orkiestracji, pasek kart sesji, wskaźniki stanu wykonania | Widoczna bez interakcji |
| 2 | Widoczna na żądanie | Wybór modułu, środowiska, modelu, wykonawcy i trybu pracy, poziom wysiłku, parametry przepływu pracy, otwarcie Execution Loop Window | Znacznik kontekstowy, ikona, przycisk, przełącznik; selektor zwija się samoczynnie po użyciu |
| 3 | Rozwinięcia kontekstowe | Zestawy akcji karty środowiska, kafla, karty sesji i pozycji bocznej nawigacji; ustawienia szybkie; warianty operacji | Menu kebab `⋮`, menu hamburger `☰`, menu kontekstowe, panel popover, lista rozwijana |
| 4 | Funkcje eksperckie | Konfiguracja punktów izolacji na poziomie roli, edycja definicji komponentów w formacie źródłowym, narzędzia diagnostyczne niskiego poziomu, tryby administracyjne | Polecenie języka naturalnego w Chat Window, skrót klawiszowy, wyszukiwarka funkcji, tryb administracyjny, konfiguracja roli; użytkownik podstawowy nie widzi tych elementów |

Wybór modułu, środowiska, modelu, wykonawcy i trybu pracy realizowany jest przez znaczniki kontekstowe i selektory zwijane po użyciu — należy do warstwy 2. Zestawy akcji realizowane są przez menu `⋮` i `☰` oraz panele popover — należą do warstwy 3. Funkcje eksperckie i administracyjne realizowane są wyłącznie przez polecenie języka naturalnego, skrót klawiszowy, wyszukiwarkę funkcji albo tryb administracyjny — należą do warstwy 4.

### 3a.3. Mechanizmy ukrywania funkcjonalności

| Mechanizm | Działanie w warstwie nawigacyjnej | Przykład |
|---|---|---|
| Menu progresywne | Zbiór jednorodnych wyborów prezentowany jest jako jeden element zwinięty; lista pozycji rozwija się po kliknięciu | `Agent ▼` w pasku kontekstu karty sesji |
| Panele wysuwane | Funkcjonalność umieszczona jest w panelach bocznych, wysuwanych, oknach popover i panelach kontekstowych; po zamknięciu panel znika całkowicie z przestrzeni roboczej | Panel zapisanych komponentów własnych kafla strefy 2 |
| Grupowanie logiczne akcji | Zamiast zestawu przycisków prezentowany jest jeden element zbiorczy, którego rozwinięcie zawiera pełną listę akcji | `Operacje ▼` w panelu kafla i w panelu orkiestracji |
| Znaczniki kontekstowe | Środowisko, repozytorium, projekt, model i wykonawca występują jako lekkie znaczniki w pasku kontekstu; kliknięcie znacznika otwiera odpowiedni selektor | `[Danaco Console] [Ubuntu] [Worktree] [Fable 5] [Ultra]` |

**Pasek znaczników kontekstowych — stan spoczynku i stan rozwinięty:**

```
  Spoczynek:
  [Danaco Console] [Ubuntu] [Worktree] [Fable 5] [Ultra]   Pętla ▼   ⋮

  Po kliknięciu znacznika [Fable 5] — selektor modelu:
  ┌────────────────────┐
  │ Fable 5      ✓     │   selektor zwija się samoczynnie
  │ Fable 5 Mini       │   natychmiast po dokonaniu wyboru
  │ Model zewnętrzny…  │
  └────────────────────┘
```

### 3a.4. Zasada jednego kliknięcia

Każda ukryta funkcja warstwy nawigacyjnej jest osiągalna jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego wydanym w Chat Window. Zagnieżdżanie funkcji głęboko w hierarchii menu jest zabronione — ukrycie zmniejsza chaos wizualny, nie utrudnia dostępu.

### 3a.5. Zasada rysowania makiet

Makiety w niniejszym dokumencie przedstawiają interfejs w stanie spoczynku: widoczne są wyłącznie elementy warstwy 1 oraz zwinięte wyzwalacze warstw 2–3 — znaczniki kontekstowe, `▼`, `⋮`, `☰`. Elementy warstw 2–4 opisane są w tabelach elementów okna, w kolumnach „Warstwa” i „Sposób wywołania”.

---

## 4. Przepływ nawigacji: strona główna → środowisko → moduł → okno

Nawigacja w platformie przebiega w czterech krokach, ustalonych w Koncepcji platformy (rozdz. 8). Niniejszy rozdział rozwija każdy z nich.

**Zestawienie czterech kroków nawigacji:**

| Krok | Z poziomu | Akcja | Prowadzi do | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|
| Krok 1 | Strona główna (strefa 1) | Kliknięcie karty środowiska | Przedsionek środowiska — kafle modułów, szyna sesji, listwa działań | 1 | Karta środowiska widoczna bez interakcji |
| Krok 2 | Okno środowiska | Wybór w bocznej nawigacji modułów | Moduł (TalkIn, WorkSpace, CodeStudio) | 1 | Pozycja bocznej nawigacji; znacznik kontekstowy modułu |
| Krok 2 (wariant) | Okno środowiska | Wybór sekcji w panelu orkiestracji | Rola lub sekcja sterowania (MultitaskingAI, rozdz. 6) | 1 | Sekcja panelu orkiestracji |
| Krok 3 | Moduł | Otwarcie okien właściwych modułowi | Okna operacyjne: Chat Window w lewej kolumnie, okna modułu w kolumnie dominującej | 1 | Otwierane automatycznie wraz z modułem |
| Krok 3 (kanał drugi) | Moduł | Otwarcie Execution Loop Window | Okno pętli wykonawczej Koordynator ↔ Wykonawca w kolumnie sąsiadującej | 2 | Przełącznik `Pętla ▼`; skrót klawiszowy; polecenie języka naturalnego |
| Krok 4 | Okno środowiska | Praca w karcie sesji; przełączanie kart | Równoległe sesje (rozdz. 7) | 1 | Pasek kart sesji |

### 4.1. Krok 1 — strona główna → środowisko

Użytkownik klika kartę środowiska w strefie 1 (rozdz. 3.2). Kliknięcie karty jest jedyną drogą wejścia do przestrzeni roboczej środowiska — strefy 2 i 3 prowadzą do okien konfiguracyjnych i ustawień, nie do środowisk. Wybór środowiska jest pierwszą i najbardziej ogólną decyzją nawigacyjną.

| Wybrana karta | Tryb pracy | Co zastępuje stronę główną |
|---|---|---|
| TalkIn | Wiedza i komunikacja | Okno środowiska z boczną nawigacją modułów oraz paskiem kart sesji |
| WorkSpace | Produktywność projektowa | Okno środowiska z boczną nawigacją modułów oraz paskiem kart sesji |
| CodeStudio | Programowanie | Okno środowiska z boczną nawigacją modułów oraz paskiem kart sesji |
| MultitaskingAI | Orkiestracja autonomicznej pracy ciągłej | Okno środowiska z panelem orkiestracji (rozdz. 6) oraz paskiem kart sesji |

### 4.2. Krok 2 — boczna nawigacja → moduł (środowiska modułowe)

W środowiskach TalkIn, WorkSpace i CodeStudio boczna nawigacja wyświetla stałą listę modułów dostępnych w danym środowisku, zgodnie z macierzą dostępności modułów (Koncepcja platformy, rozdz. 10; Załącznik A). Lista jest widoczna niezależnie od tego, który moduł jest aktualnie otwarty w bieżącej karcie sesji. Kliknięcie pozycji modułu na liście przeładowuje przestrzeń roboczą bieżącej karty do układu właściwego temu modułowi — mechanizm przeładowania opisano w rozdziale 8.1. Pełny opis bocznej nawigacji zawiera rozdział 5.

### 4.3. Krok 2 (wariant) — panel orkiestracji (środowisko MultitaskingAI)

Środowisko MultitaskingAI nie ma bocznej nawigacji modułów — jego odpowiednikiem jest panel orkiestracji, opisany szczegółowo w rozdziale 6. Różnica wynika z charakteru środowiska: MultitaskingAI organizuje nie zadania modułowe, lecz zespół modeli i agentów realizujących wspólny proces (Koncepcja platformy, rozdz. 13.1), dlatego jego boczna nawigacja prowadzi po rolach, kolejkach i orkiestracji, a nie po modułach.

### 4.4. Krok 3 — moduł → okno kontekstowe

Wybrany moduł otwiera właściwy mu zestaw okien operacyjnych — dopiero na tym poziomie zaczyna się praca właściwa.

| Cecha | Zasada | Warstwa | Sposób wywołania |
|---|---|---|---|
| Zestaw okien modułu | Każdy moduł definiuje własny zestaw okien, rozmieszczonych w kolumnach na prawo od okien komunikacji (Koncepcja platformy, rozdz. 11; Załącznik A Koncepcji platformy) | 1 | Otwierane wraz z modułem |
| Chat Window | Okno wspólne obecne na stronie głównej i w każdym z piętnastu modułów wyliczonych w Koncepcji platformy, rozdz. 11.1–11.15, w lewej kolumnie obszaru roboczego; centralny punkt pracy i podstawowy mechanizm sterowania procesami | 1 | Widoczne bez interakcji |
| Execution Loop Window | Okno pętli wykonawczej Koordynator ↔ Wykonawca, otwierane w kolumnie sąsiadującej z Chat Window; dostępne z poziomu strony głównej i z każdego modułu | 2 | Przełącznik `Pętla ▼`; skrót klawiszowy; polecenie języka naturalnego |
| Okna pomocnicze modułu | Panele i monitory otwierane jako rozszerzenia boczne po prawej stronie obszaru roboczego | 3 | Menu `☰` obszaru roboczego; panel popover |

### 4.5. Krok 4 — karty sesji

Aktywne sesje trzymane są w kartach poziomych w oknie środowiska. Karta aktywnej sesji reprezentuje jedną, samodzielną przestrzeń roboczą — z własnym układem okien, historią i kontekstem. Przełączanie kart oznacza przełączanie między sesjami, na wzór kart w przeglądarce internetowej. Pełną mechanikę kart sesji opisuje rozdział 7.

### 4.6. Schemat przepływu nawigacji

```
┌─────────────────────────────────────────────────────────────┐
│ STRONA GŁÓWNA — Centrum dowodzenia                          │
│  Chat Window (lewa kolumna) │ Strefa 1 · Strefa 2 · Strefa 3│
└─────────────────────────────────────────────────────────────┘
        │  KROK 1 — kliknięcie karty środowiska (strefa 1)
        ▼
┌─────────────────────────────────────────────────────────────┐
│ ŚRODOWISKO  TalkIn · WorkSpace · CodeStudio · MultitaskingAI│
└─────────────────────────────────────────────────────────────┘
        │  KROK 2 — boczna nawigacja modułów  (TalkIn, WorkSpace, CodeStudio)
        │           albo panel orkiestracji    (MultitaskingAI — rozdz. 6)
        ▼
┌─────────────────────────────────────────────────────────────┐
│ MODUŁ  (np. Studio, Research, Developer)                    │
│   albo  ROLA / SEKCJA panelu orkiestracji                   │
└─────────────────────────────────────────────────────────────┘
        │  KROK 3 — otwarcie okien właściwych modułowi
        ▼
┌─────────────────────────────────────────────────────────────┐
│ OKNA OPERACYJNE — układ kolumnowy lewa–prawa                │
│  Chat Window │ Execution Loop Window │ okna modułu │ panele │
└─────────────────────────────────────────────────────────────┘
        │  KROK 4 — praca odbywa się w karcie sesji bieżącej
        ▼
  [ Karta sesji A · Studio ]  [ Karta sesji B · Research ]  [ + nowa karta ]
```

---

## 5. Boczna nawigacja modułów

### 5.1. Mechanika i widoczność według macierzy dostępności modułów

Boczna nawigacja jest stałym, pionowym elementem interfejsu środowiska. Środowisko jest profilem widoczności modułów w bocznej nawigacji, a nie pojemnikiem, do którego moduł należy wyłącznie (Koncepcja platformy, rozdz. 10).

| Cecha bocznej nawigacji | Charakterystyka | Warstwa | Sposób wywołania |
|---|---|---|---|
| Charakter elementu | Skrajna lewa kolumna interfejsu środowiska, stała, pełna wysokość | 1 | Widoczna bez interakcji |
| Widoczność | Widoczna niezależnie od tego, który moduł jest otwarty w bieżącej karcie sesji (Koncepcja platformy, rozdz. 8, krok 2) | 1 | Widoczna bez interakcji |
| Zawartość | Lista modułów dostępnych w danym środowisku, zgodnie z macierzą dostępności (Załącznik A) | 1 | Widoczna bez interakcji |
| Zależność od środowiska | Środowisko jest profilem widoczności modułów; moduł jest dostępny w tych środowiskach, dla których stanowi naturalne narzędzie pracy | 1 | Widoczna bez interakcji |
| Wspólność dla kart sesji | Wspólna dla wszystkich kart sesji otwartych w danym środowisku — nie jest odrębna per karta | 1 | Widoczna bez interakcji |
| Moduł Automations | Nie występuje w żadnej z trzech list — jest konfigurowany wyłącznie ze strony głównej (strefa 2) i trafia do sesji jako gotowa automatyka, komponent własny (Koncepcja platformy, rozdz. 10, rozdz. 11.3) | 1 | Kafel strefy 2 strony głównej |
| Zestaw akcji pozycji modułu | Otwarcie modułu w nowej karcie sesji, przypięcie modułu, ustawienia domyślne modułu | 3 | Menu kebab `⋮` przy pozycji modułu |
| Zwinięcie bocznej nawigacji | Zwinięcie listy do pasa ikon i jej ponowne rozwinięcie | 3 | Menu hamburger `☰` nad listą |
| Wybór modelu i wykonawcy dla modułu | Zmiana modelu, wykonawcy oraz trybu pracy obowiązującego w otwartym module | 2 | Kliknięcie znacznika kontekstowego w pasku kontekstu |
| Konfiguracja widoczności pozycji nawigacji dla roli | Ustalenie zestawu modułów widocznych dla danej roli użytkownika | 4 | Tryb administracyjny; konfiguracja roli; polecenie języka naturalnego w Chat Window |

Pełną macierz dostępności piętnastu modułów w trzech środowiskach modułowych (TalkIn, WorkSpace, CodeStudio) przedstawia Załącznik A (powtórzenie macierzy z rozdziału 10 Koncepcji platformy).

### 5.2. Stan wybrany i wskaźniki

Pozycja modułu aktualnie otwartego w bieżącej karcie sesji jest oznaczona w bocznej nawigacji stanem wybranym (podświetlenie pozycji), tak aby użytkownik zawsze zachowywał orientację co do tego, w jakim module pracuje — zgodnie z zasadą, że zmiana przestrzeni jest zawsze świadomym, widocznym przejściem, a nie ukrytą zmianą stanu (Koncepcja platformy, rozdz. 1).

| Aspekt | Zasada | Warstwa | Sposób wywołania |
|---|---|---|---|
| Oznaczenie modułu otwartego | Stan wybrany — podświetlenie pozycji w bocznej nawigacji, zweryfikowane wobec `.dn-boczna-pozycja[aria-current='page']` ([System wizualny](system-wizualny.md) rozdz. 14.12): tło `--dn-sygnal-tlo` oraz znacznik kreskowy dodatkowy — pseudo-element szerokości `--dn-wym-wstega` (2px), wysokości 16px, tło `--dn-kropka`, przy lewej krawędzi pozycji, nazwany w komentarzu źródłowym „krótka kreska sygnału zamiast pełnej wstęgi” | 1 | Widoczne bez interakcji |
| Wskaźnik stanu wykonania | Sygnalizacja trwającej pracy modułu przy pozycji nawigacji | 1 | Widoczny bez interakcji |
| Skutek kliknięcia innej pozycji | Przeładowanie przestrzeni roboczej bieżącej karty do układu nowo wybranego modułu (rozdz. 8.1) | 1 | Kliknięcie pozycji |
| Który moduł jest podświetlony | Zależy od modułu otwartego w aktualnie aktywnej (frontowej) karcie sesji | 1 | Widoczne bez interakcji |
| Zmiana podświetlenia | Następuje wraz z przełączaniem kart sesji (rozdz. 7.2) | 1 | Przełączenie karty sesji |

### 5.3. Przykład: nawigacja w środowisku TalkIn

Zgodnie z macierzą dostępności modułów (Załącznik A) boczna nawigacja środowiska TalkIn wyświetla dziewięć pozycji.

| Widoczne w bocznej nawigacji TalkIn (9 pozycji) | Niewidoczne — poza macierzą TalkIn |
|---|---|
| Studio, Workspace, Browser, Research, Library, Translate, Roundtable, Assistant, Agents | Design, Terminal, Developer, Diagnostics, Apps; ponadto Automations (konfigurowany wyłącznie ze strony głównej) |

**Makieta środowiska TalkIn w stanie spoczynku (układ kolumnowy lewa–prawa):**

```
 ═══════════════════════════════════════════════════════════════════════════
  Boczna       │ Chat Window            │ Obszar roboczy modułu Research
  nawigacja  ☰ │ Użytkownik ↔ Wykonawca │
  ─────────────│                        │ Research Workspace · Sources
    Studio     │                        │ Manager · Findings Panel
    Workspace  │                        │
    Browser    │ ──────────────────     │
  ▸ Research ⋮ │ │ Polecenie…     │     │
    Library    │ ──────────────────     │
    Translate  │                        │
    Roundtable │ [Danaco Console]       │
    Assistant  │ [Ubuntu] [Worktree]    │
    Agents     │ [Fable 5] [Ultra]      │
               │ Pętla ▼            ⋮   │
 ═══════════════════════════════════════════════════════════════════════════
   ▸ stan wybrany (podświetlenie) · widoczna wyłącznie warstwa 1
   oraz zwinięte wyzwalacze warstw 2–3: znaczniki, `Pętla ▼`, `⋮`, `☰`
```

Otwarcie Execution Loop Window wstawia kolumnę pętli wykonawczej między Chat Window a obszar roboczy modułu.

**Mini-przepływ — wybór modułu Research w środowisku TalkIn:**

1. Użytkownik klika pozycję Research w bocznej nawigacji TalkIn.
2. Bieżąca karta sesji przeładowuje się do zestawu okien modułu Research: Chat Window, Research Workspace, Sources Manager, Findings Panel, Report Builder, Export Panel (Koncepcja platformy, rozdz. 11.5).
3. Pozycja Research przyjmuje w bocznej nawigacji stan wybrany.

---

## 6. Panel orkiestracji — boczna nawigacja środowiska MultitaskingAI

### 6.1. Różnica względem bocznej nawigacji modułowej

Środowisko MultitaskingAI organizuje nie zadania modułowe, lecz zespół modeli i agentów realizujących wspólny proces (Koncepcja platformy, rozdz. 13.1), dlatego jego boczna nawigacja nie zawiera listy modułów, lecz panel orkiestracji.

| Cecha | Boczna nawigacja modułów (TalkIn, WorkSpace, CodeStudio) | Panel orkiestracji (MultitaskingAI) |
|---|---|---|
| Co organizuje | Zadania modułowe | Zespół modeli i agentów realizujących wspólny proces |
| Zawartość nawigacji | Lista modułów wg macierzy dostępności | Sześć sekcji sterowania |
| Po czym nawiguje | Po modułach | Po rolach, kolejkach i orkiestracji |
| Podstawa | Koncepcja platformy, rozdz. 10 | Koncepcja platformy, rozdz. 13.2–13.6, rozdz. 13.8 |

Panel orkiestracji jest funkcjonalnym odpowiednikiem bocznej nawigacji modułów: nawiguje po tym, czym w tym środowisku faktycznie się steruje (Koncepcja platformy, rozdz. 13.8).

### 6.2. Sześć sekcji panelu

Panel orkiestracji zawiera sześć sekcji, ustalonych w Koncepcji platformy (rozdz. 13.8):

| Sekcja | Zawartość | Podstawa (Koncepcja platformy) | Powiązane mechanizmy i okna | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|
| Zespoły | Zapisane konfiguracje zespołu — szablony ról, powiązań i kolejek; zapis, wczytanie, duplikowanie | rozdz. 13.1 | Konfiguracje ról (rozdz. 13.3) | 1 | Pozycja panelu widoczna bez interakcji |
| Role | Cztery okna robocze: Executor 1, Executor 2, Coordinator, Executor 3 / Validator; przypisanie agentów z modułu Agents do ról; Subagent Network | rozdz. 13.3, 11.15 | Agent Builder, Permissions Center (moduł Agents) | 1 | Pozycja panelu widoczna bez interakcji |
| Kolejki | Definicje kolejek globalnych, lokalnych, modeli, agentów i projektów oraz ich akcje | rozdz. 13.4 | Silnik kolejek, Queue Manager (moduł Automations) | 1 | Pozycja panelu widoczna bez interakcji |
| Orkiestracja | Zależności między modelami, agentami, zadaniami, kolejkami, automatyzacjami i projektami | rozdz. 13.5 | Orchestrator (moduł Automations) | 1 | Pozycja panelu widoczna bez interakcji |
| Harmonogram i automatyki | Harmonogram pracy ciągłej oraz wpięte automatyki z modułu Automations | rozdz. 13.6, 11.3 | Scheduler, Execution Monitor (moduł Automations) | 1 | Pozycja panelu widoczna bez interakcji |
| Monitor procesu | Podgląd przebiegu pętli, statusy przebiegów, hierarchia decyzji; miejsce nadzoru Always On Display | rozdz. 13.2, 9.4, 12.2 | Always On Display, Mobile | 1 | Pozycja panelu widoczna bez interakcji |
| Przypisanie modelu i wykonawcy do roli | Wybór modelu, wykonawcy i poziomu wysiłku dla wybranej roli | rozdz. 13.3 | Agent Builder | 2 | Kliknięcie znacznika kontekstowego roli; selektor zwija się po wyborze |
| Zestaw akcji sekcji | Duplikowanie, eksport, wyczyszczenie i przywrócenie ustawień domyślnych sekcji | rozdz. 13.8 | — | 3 | Menu kebab `⋮` przy sekcji; rozwinięcie `Operacje ▼` |
| Kolejność i widoczność sekcji | Ustawienie zestawu i porządku sekcji panelu | rozdz. 6.3 niniejszego dokumentu | Okno konfiguracji (rozdz. 9.1) | 4 | Tryb administracyjny; wyszukiwarka funkcji; polecenie języka naturalnego w Chat Window |

**Makieta środowiska MultitaskingAI w stanie spoczynku (układ kolumnowy lewa–prawa):**

```
 ═══════════════════════════════════════════════════════════════════════════
  Panel               │ Chat Window            │ Obszar roboczy procesu
  orkiestracji      ☰ │ Użytkownik ↔ Wykonawca │
  ────────────────────│                        │ Executor 1 · Executor 2
    Zespoły           │                        │ Coordinator
  ▸ Role            ⋮ │ ──────────────────     │ Executor 3 / Validator
    Kolejki           │ │ Polecenie…     │     │
    Orkiestracja      │ ──────────────────     │
    Harmonogram       │                        │
      i automatyki    │ [Danaco Console]       │
    Monitor procesu   │ [Ubuntu] [Fable 5]     │
                      │ [Ultra]  Pętla ▼   ⋮   │
 ═══════════════════════════════════════════════════════════════════════════
   ▸ stan wybrany · sekcja Monitor procesu jest miejscem nadzoru
   Always On Display · widoczna wyłącznie warstwa 1
```

Kliknięcie sekcji „Role” otwiera cztery okna robocze odpowiadające czterem rolom środowiska (Executor 1, Executor 2, Coordinator, Executor 3 / Validator), w których użytkownik przypisuje agentów z modułu Agents lub modele bazowe, zgodnie z szablonem roli (Koncepcja platformy, Załącznik D.2). Kliknięcie sekcji „Monitor procesu” otwiera widok nadzoru, w którym Always On Display może pełnić funkcję obserwatora lub operatora procesu (Koncepcja platformy, rozdz. 13.2).

### 6.3. Konfigurowalność kolejności i widoczności

Zgodnie z zasadą pełnej konfigurowalności (Koncepcja platformy, rozdz. 6, rozdz. 14, zasada 2) kolejność i widoczność sekcji panelu orkiestracji podlegają konfiguracji z poziomu okna konfiguracji. Przedstawiony w rozdziale 6.2 zestaw sześciu sekcji jest zestawem domyślnym, obowiązującym przy braku odmiennego ustawienia użytkownika — zgodnie z zasadą nadrzędną „brak ustawienia = wartość domyślna”. Panel nigdy nie wymusza ukrycia żadnej sekcji ani nie blokuje możliwości przywrócenia zestawu domyślnego.

**Szablon konfiguracji panelu orkiestracji (redakcyjny):**

```
# Panel orkiestracji — kolejność i widoczność sekcji
# Zasada nadrzędna: brak ustawienia = zestaw i kolejność domyślna
panel_orkiestracji:
  sekcje:                                    # kolejność domyślna wg rozdz. 6.2
    - nazwa: "Zespoły"                    widoczna: tak   kolejnosc: 1
    - nazwa: "Role"                       widoczna: tak   kolejnosc: 2
    - nazwa: "Kolejki"                    widoczna: tak   kolejnosc: 3
    - nazwa: "Orkiestracja"               widoczna: tak   kolejnosc: 4
    - nazwa: "Harmonogram i automatyki"   widoczna: tak   kolejnosc: 5
    - nazwa: "Monitor procesu"            widoczna: tak   kolejnosc: 6
  zrodlo_konfiguracji: "okno konfiguracji"   # rozdz. 6.3, rozdz. 9.1
  przywrocenie_domyslnych: dostepne          # panel nigdy nie wymusza ukrycia sekcji
```

---

## 7. Karty sesji — mechanika

### 7.1. Anatomia karty sesji

Karta sesji jest poziomym elementem paska kart w oknie środowiska, reprezentującym jedną, samodzielną przestrzeń roboczą. Mechanika jest analogiczna do kart w przeglądarce internetowej (Koncepcja platformy, rozdz. 8, krok 4).

| Element karty sesji | Zawartość | Uwagi | Warstwa | Sposób wywołania |
|---|---|---|---|---|
| Tytuł karty | Nazwa modułu otwartego w karcie (np. „Studio”, „Developer”) | Aktualizuje się przy przeładowaniu modułu (rozdz. 8.1) | 1 | Widoczny bez interakcji |
| Wskaźnik stanu wykonania | Sygnalizacja aktywnej pracy w tle (trwający proces automatyki, aktywna rola MultitaskingAI) | Powiązany z funkcją Mobile — monitoring procesów (Koncepcja platformy, rozdz. 12.1) | 1 | Widoczny bez interakcji |
| Kontrolka zamknięcia | Zamyka kartę sesji | Analogicznie do kart przeglądarki | 1 | Widoczna bez interakcji |
| Kontrolka nowej karty („+”) | Otwiera nową, pustą kartę sesji w bieżącym środowisku | Oczekuje na wybór modułu z bocznej nawigacji | 1 | Widoczna bez interakcji |
| Obszar przełączania | Kliknięcie dowolnego miejsca karty (poza kontrolką zamknięcia) czyni ją kartą aktywną | Przełącza widok obszaru roboczego | 1 | Kliknięcie karty |
| Znaczniki kontekstu karty | Środowisko, repozytorium, projekt, model i wykonawca przypisane karcie | Selektor zwija się po wyborze | 2 | Kliknięcie znacznika w pasku kontekstu |
| Zestaw akcji karty | Duplikowanie karty, przypięcie, zamknięcie pozostałych kart, zmiana nazwy | Grupowanie logiczne akcji | 3 | Menu kontekstowe karty; menu kebab `⋮` |
| Profil izolacji karty na poziomie roli | Przypisanie profilu izolacji do konkretnej roli w obrębie karty (rozdz. 7.4) | Funkcja ekspercka | 4 | Konfiguracja roli; tryb administracyjny; polecenie języka naturalnego w Chat Window |

**Schemat paska kart sesji:**

```
  ┌────────────────────────┬────────────────────────┬───────┐
  │  ●  Studio           ✕ │     Research         ✕ │   +   │
  └────────────────────────┴────────────────────────┴───────┘
     ▲ karta aktywna            ▲ karta w tle            ▲ nowa karta
     ●  wskaźnik pracy w tle    ✕  kontrolka zamknięcia
```

**Schemat pojedynczej karty sesji:**

```
  ┌──────────────────────────────┐
  │  ●   Studio               ✕  │
  │  │      │                  │ │
  │  │      │                  └── kontrolka zamknięcia karty
  │  │      └── tytuł = nazwa otwartego modułu
  │  │          (aktualizuje się przy przeładowaniu, rozdz. 8.1)
  │  └── wskaźnik stanu = aktywna praca w tle
  │      (powiązany z funkcją Mobile, Koncepcja platformy, rozdz. 12.1)
  └──────────────────────────────┘
   Kliknięcie w dowolnym miejscu karty (poza „✕”) = przełączenie na tę sesję
```

### 7.2. Otwieranie, zamykanie, przełączanie

| Operacja | Wyzwalacz | Skutek |
|---|---|---|
| Otwarcie nowej karty | Kontrolka „+” | Nowa, pusta przestrzeń robocza w bieżącym środowisku; boczna nawigacja (lub panel orkiestracji) aktywna; oczekiwanie na wybór modułu — analogicznie do nowej karty przeglądarki bez załadowanej strony |
| Przełączenie karty | Kliknięcie karty (poza „✕”) | Główny obszar roboczy przełącza się na układ okien, historię i kontekst zapisane w tej karcie; boczna nawigacja lub panel orkiestracji podświetla właściwy moduł lub sekcję (rozdz. 5.2, rozdz. 6.1) |
| Zamknięcie karty | Kontrolka „✕” | Zakończenie przestrzeni roboczej karty w widoku bieżącym; los sesji leżącej u podstaw karty jest funkcją modelu sesji i procesów (Architektura, rozdz. 7) i pozostaje poza zakresem niniejszego dokumentu nawigacyjnego |

Równoległa praca w wielu modułach lub w wielu miejscach tego samego środowiska odbywa się przez otwarcie kilku kart sesji obok siebie — to mechanizm, który rozstrzyga pozorną sprzeczność między zarządzaniem przez platformę wieloma równoległymi przestrzeniami roboczymi (Koncepcja platformy, rozdz. 1) a przeładowaniem przestrzeni roboczej przy zmianie modułu w obrębie jednej karty (Koncepcja platformy, rozdz. 2.2; rozdz. 10.1): w obrębie jednej karty widoczny jest w danej chwili interfejs jednego modułu, a równoległość realizują wielokrotne, jednocześnie otwarte karty.

### 7.3. Domyślna izolacja karty sesji i możliwość współdzielenia

Każda nowa karta sesji otrzymuje domyślnie odrębny układ okien, odrębną historię i odrębny kontekst — jest to stan wyjściowy okna konfiguracji punktów izolacji (Koncepcja platformy, rozdz. 6.1, rozdz. 6.6), a nie ograniczenie możliwości platformy.

| Aspekt karty sesji | Stan domyślny | Ustawienie konfiguracyjne | Warstwa | Sposób wywołania |
|---|---|---|---|---|
| Układ okien | Odrębny dla każdej nowej karty | Wynika z wybranego modułu | 1 | Widoczny bez interakcji |
| Historia | Odrębna | Jednorazowe współdzielenie między konkretnymi kartami (okno konfiguracji, rozdz. 9.1) | 4 | Okno konfiguracji punktów izolacji; wyszukiwarka funkcji; polecenie języka naturalnego w Chat Window |
| Pamięć | Odrębna | Jak wyżej | 4 | Jak wyżej |
| Kontekst | Odrębny | Jak wyżej | 4 | Jak wyżej |

Karta sesji jest jednym z ośmiu poziomów zasięgu reguł izolacji (Koncepcja platformy, rozdz. 6.5; [Izolacja i zależności](../architektura/izolacja-i-zaleznosci.md) rozdz. 6) — użytkownik może z poziomu okna konfiguracji (rozdz. 9.1) ustanowić jednorazowe współdzielenie historii, pamięci lub kontekstu między konkretnymi kartami, zgodnie z przykładem zastosowania „jednorazowe współdzielenie historii między dwiema otwartymi kartami” z tabeli poziomów zasięgu Koncepcji platformy (rozdz. 6.5).

**Szablon konfiguracji karty sesji (redakcyjny):**

```
# Karta sesji — izolacja domyślna i współdzielenie
# Stan wyjściowy = domyślny, uporządkowany podział pracy (rozdz. 7.3; Koncepcja platformy, rozdz. 6.1, 6.6)
karta_sesji:
  izolacja_kontekstu:                        # poziom zasięgu: karta sesji (rozdz. 6.5 Koncepcji)
    historia: odrebna                        # domyślnie odrębna dla każdej nowej karty
    pamiec:   odrebna
    kontekst: odrebny
  wspoldzielenie_miedzy_kartami:             # jednorazowe, świadoma decyzja użytkownika
    dozwolone: tak                           # ustanawiane z okna konfiguracji (rozdz. 9.1)
    zakres:    [historia, pamiec, kontekst]
  pierwszenstwo_zasiegu:                      # rola węższa niż karta sesji (rozdz. 7.4)
    rola_ma_pierwszenstwo_przed_karta: tak    # np. profil roli Executor 1 (MultitaskingAI)
```

### 7.4. Karta sesji a rola w środowisku MultitaskingAI

W środowisku MultitaskingAI karta sesji reprezentuje przestrzeń roboczą pojedynczego procesu orkiestracji — obejmującą stan ról, kolejek i orkiestracji widoczny przez panel orkiestracji (rozdz. 6), a nie pojedynczy moduł. Poziom zasięgu „rola” w hierarchii izolacji (Koncepcja platformy, rozdz. 6.5) jest węższy niż karta sesji; poziomem najbardziej szczegółowym jest okno komunikacji (`window`), węższe od roli — dwa okna jednej karty sesji mogą mieć odrębną pamięć i odrębne uprawnienia ([Izolacja i zależności](../architektura/izolacja-i-zaleznosci.md) rozdz. 6).

```
Zasięg reguł izolacji (pierwszeństwo rośnie w prawo):
  … ─▶ Karta sesji ─▶ Rola (np. Executor 1)
                       ▲ pierwszeństwo — profil roli wygrywa z ustawieniem karty sesji
                         (zasada pierwszeństwa zasięgu najbardziej szczegółowego)
```

Profil izolacji przypisany konkretnej roli (na przykład Executor 1) ma pierwszeństwo przed ustawieniem izolacji ustalonym na poziomie karty sesji, zgodnie z zasadą pierwszeństwa zasięgu najbardziej szczegółowego.

---

## 8. Przełączanie i przeładowanie przestrzeni roboczej

### 8.1. Przeładowanie w obrębie karty — zmiana modułu

Przełączenie modułu w bocznej nawigacji (lub sekcji w panelu orkiestracji) w obrębie tej samej karty sesji przeładowuje przestrzeń roboczą tej karty: znikają okna właściwe poprzedniemu modułowi, a w ich miejsce pojawia się zestaw okien operacyjnych dopasowany do nowo wybranego modułu (Koncepcja platformy, rozdz. 2.2). Ten mechanizm przeładowania utrzymuje czytelność pojedynczej karty sesji — w obrębie jednej karty widoczny jest zawsze interfejs dokładnie jednego modułu.

| Przeładowanie OBEJMUJE | Przeładowanie NIE OBEJMUJE |
|---|---|
| Zestaw okien operacyjnych (zgodnie z katalogiem modułu, Koncepcja platformy, rozdz. 11) | Samej karty sesji jako takiej — karta pozostaje tą samą kartą, zmienia się jedynie jej zawartość |
| Kontekst czatu (rekonfigurowany do nowego modułu, Koncepcja platformy, rozdz. 2.4) | Bocznej nawigacji lub panelu orkiestracji — pozostają widoczne w niezmienionej formie |
| Narzędzia i przebieg pracy dedykowane nowemu modułowi | Pozostałych otwartych kart sesji — każda zachowuje własny, niezależny stan |
| — | Funkcji globalnych Mobile i Always On Display, które nie podlegają przeładowaniu w żadnym wypadku (Koncepcja platformy, rozdz. 5.3) |

### 8.2. Przełączanie środowiska — powrót do centrum dowodzenia

Zmiana środowiska następuje przez powrót do strony głównej i wybór innej karty środowiska w strefie 1 (rozdz. 4.1). W odróżnieniu od przeładowania modułu w obrębie karty, przełączenie środowiska jest decyzją bardziej ogólną: zmienia całą przestrzeń roboczą — dostępne narzędzia, model nawigacji (boczna nawigacja modułów albo panel orkiestracji) oraz historię sesji przypisaną do tego trybu pracy (Koncepcja platformy, rozdz. 2.1). Zgodnie z zasadą, że zmiana przestrzeni jest zawsze świadomym, widocznym przejściem, a nie ukrytą zmianą stanu (Koncepcja platformy, rozdz. 1), przejście przez stronę główną — a nie ukryty przełącznik wewnątrz okna środowiska — jest tym widocznym, świadomym krokiem.

### 8.3. Trwałość sesji w tle

Powrót do strony głównej w celu przełączenia środowiska nie kończy kart sesji otwartych w środowisku, które użytkownik opuszcza. Zasada wynika z modelu sesji i procesów (Architektura, rozdz. 7).

```
Sesja jako odrębny proces po stronie serwera (Architektura, rozdz. 7):
  rozłączenie klienta ─▶ proces sesji pracuje dalej ─▶ ponowne połączenie podłącza klienta
                                                        (odtworzony pełny stan: układ, historia, kontekst)
```

Karty sesji otwarte w środowisku opuszczonym przez użytkownika pozostają aktywne w tle i odtwarzają swój pełny stan (układ, historię, kontekst) w chwili powrotu do tego środowiska — zgodnie z charakterem platformy jako systemu operacyjnego zarządzającego wieloma równoległymi przestrzeniami roboczymi (Koncepcja platformy, rozdz. 1). Ta trwałość jest tym, co pozwala funkcji Mobile monitorować i nadzorować procesy uruchomione w tle niezależnie od tego, które środowisko jest aktualnie otwarte na ekranie (Koncepcja platformy, rozdz. 12.1).

### 8.4. Zestawienie: co się zmienia, a co jest zachowywane

| Zdarzenie | Co się przeładowuje / zmienia | Co pozostaje zachowane |
|---|---|---|
| Zmiana modułu w bocznej nawigacji (w tej samej karcie) | Zestaw okien operacyjnych, kontekst czatu, narzędzia i przebieg pracy | Sama karta sesji, boczna nawigacja, pozostałe karty, funkcje globalne |
| Otwarcie nowej karty sesji | — (nowa, pusta przestrzeń robocza) | Wszystkie dotychczas otwarte karty, boczna nawigacja / panel orkiestracji, funkcje globalne |
| Zamknięcie karty sesji | Znika widok karty zamkniętej | Pozostałe karty, stan środowiska, funkcje globalne |
| Powrót do strony głównej i wybór innego środowiska | Cała przestrzeń robocza: narzędzia, model nawigacji, historia sesji przypisana do trybu pracy | Karty sesji poprzedniego środowiska (trwają w tle, rozdz. 8.3), funkcje globalne |
| Aktywacja Mobile lub Always On Display | Nic — funkcje globalne nie tworzą nowej przestrzeni roboczej | Cała bieżąca przestrzeń robocza, wszystkie karty sesji |

---

## 9. Wejście do ustawień i funkcji globalnych

### 9.1. Okno konfiguracji

Pozycja „Okno konfiguracji” w listwie ustawień strony głównej (rozdz. 3.4) otwiera okno udostępniające pełny zakres ustawień wpływających na aplikację. To samo okno jest dostępne również z poziomu każdego środowiska — nie jest unikalne dla strony głównej, lecz stanowi warstwę narzędziową obecną w całej platformie, zgodnie z zasadą pełnej konfigurowalności (Koncepcja platformy, rozdz. 14, zasada 2). Każdy element konfiguracji zawiera objaśnienie kontekstowe (`[?]`), opisujące jego działanie i wpływ na aplikację (Architektura, rozdz. 13).

**Zakres ustawień okna konfiguracji** (Architektura, rozdz. 13):

| Obszar ustawień | Zakres |
|---|---|
| Aplikacja i procesy | Zachowanie aplikacji, procesy, akcje |
| Modele | Zachowanie modeli, tożsamość modeli, prompty systemowe |
| Rozszerzenia | Rozszerzenia, integracje |
| Kontekst | Historia, pamięć, izolacja |

Wśród ustawień dostępnych w tym oknie znajduje się okno konfiguracji punktów izolacji, łączące izolację kontekstu (historia, pamięć, kontekst) i izolację techniczną procesu sesji (osiem zakresów zdefiniowanych w rozdziale 12 Architektury) w trzech panelach:

| Panel okna konfiguracji punktów izolacji | Zawartość |
|---|---|
| Selektor zasięgu | Osiem poziomów: globalny, środowisko, moduł, para modułów, projekt, karta sesji, rola, okno komunikacji ([Izolacja i zależności](../architektura/izolacja-i-zaleznosci.md) rozdz. 6) |
| Macierz izolacji | Przełączniki dla obu rodzajów izolacji: izolacji kontekstu (historia, pamięć, kontekst) oraz izolacji technicznej (osiem zakresów, Architektura, rozdz. 12) |
| Profil i podgląd polityki efektywnej | Zapis, wczytanie i przypisanie profilu izolacji; wybór warstwy konfiguracji (domyślna lub sesji); podgląd polityki efektywnej |

Pełna specyfikacja tego okna — układ paneli, reguła pierwszeństwa zasięgu najbardziej szczegółowego, profile izolacji i warstwy konfiguracji — znajduje się w Koncepcji platformy (rozdz. 6.4–6.6), do której niniejszy dokument odsyła w zakresie mechanizmów izolacji przywołanych przy kartach sesji (rozdz. 7.3) i przy rolach środowiska MultitaskingAI (rozdz. 7.4).

### 9.2. Mobile i Always On Display

Pozycje „Mobile” i „Always On Display” w listwie ustawień strony głównej są jednym z kilku równoważnych punktów dostępu do tych dwóch funkcji globalnych.

| Funkcja globalna | Punkty dostępu | Cecha wspólna |
|---|---|---|
| Mobile | Listwa ustawień strony głównej; każde środowisko; każdy moduł | Nie tworzy własnej przestrzeni roboczej, nie przeładowuje środowiska, nie tworzy nowych sesji (Koncepcja platformy, rozdz. 5.3, rozdz. 12); pełne opracowanie funkcji: `funkcje-globalne/mobile.md` |
| Always On Display | Listwa ustawień strony głównej; każde środowisko; każdy moduł | Jak wyżej |

W kontekście strony głównej ich obecność w strefie 3 sygnalizuje, że są one częścią warstwy narzędziowej platformy, dostępnej niezależnie od tego, czy użytkownik znajduje się na stronie głównej, czy w dowolnym środowisku.

---

## 10. Zgodność z zasadami nadrzędnymi platformy

Nawigacja opisana w niniejszym dokumencie pozostaje w pełnej zgodności z czterema zasadami nadrzędnymi platformy (Koncepcja platformy, rozdz. 14).

| Zasada nadrzędna | Zastosowanie w warstwie nawigacyjnej | Miejsce w dokumencie |
|---|---|---|
| Pełna kompozycyjność | Kolejność kart sesji, liczba jednocześnie otwartych kart oraz kombinacje modułów otwartych w poszczególnych kartach nie są ograniczone przez platformę; użytkownik komponuje własny układ pracy z dostępnych środowisk, modułów i kart sesji | rozdz. 7 |
| Pełna konfigurowalność (zasada centralna) | Zakres współdzielenia historii, pamięci i kontekstu między kartami sesji, kolejność i widoczność sekcji panelu orkiestracji oraz wszystkie punkty izolacji nawigacyjnej konfiguruje się z poziomu okna konfiguracji; brak ustawienia oznacza wartość domyślną, nigdy blokadę | rozdz. 6.3, 7.3, 9.1 |
| Jawność i konfigurowalność zależności | Każde przełączenie środowiska lub modułu jest świadomym, widocznym przejściem, a nie ukrytą zmianą stanu; żadna zależność nawigacyjna między kartami, modułami czy środowiskami nie jest zaszyta na stałe | rozdz. 8.2 |
| Rozszerzenie orkiestracji | Panel orkiestracji środowiska MultitaskingAI jest obecnym, najpełniejszym wcieleniem tego kierunku w warstwie nawigacyjnej platformy | rozdz. 6 |

Warstwa nawigacyjna realizuje ponadto regułę warstw widoczności (rozdz. 3a): strona główna i nawigacja prezentują wyłącznie warstwę 1, a funkcje warstw 2–4 ujawniane są przez znaczniki kontekstowe, menu `⋮` i `☰`, panele popover, skróty klawiszowe oraz polecenia języka naturalnego wydawane w Chat Window. Liczba modułów, agentów, przepływów pracy i ustawień nie wpływa na postrzeganą prostotę interfejsu.

Zgodnie z dyrektywą braku twardych blokad żaden z opisanych mechanizmów — przeładowanie przestrzeni roboczej, domyślna izolacja karty sesji, domyślny zestaw i kolejność sekcji panelu orkiestracji — nie jest nieodwracalny ani niekonfigurowalny: każdy z nich jest stanem wyjściowym, możliwym do zmiany z poziomu okna konfiguracji (rozdz. 9.1).

---

## 11. Słowniczek pojęć nawigacyjnych

| Pojęcie | Definicja | Rozdział |
|---|---|---|
| Centrum dowodzenia | Inna nazwa strony głównej aplikacji; układ kolumnowy lewa–prawa złożony z Chat Window w lewej kolumnie oraz trzech stref w kolumnie dominującej: wyboru środowiska, konfiguracji komponentów własnych i ustawień platformy | rozdz. 2 |
| Chat Window | Główne okno komunikacji Użytkownik ↔ Wykonawca; centralny punkt pracy użytkownika i podstawowy mechanizm sterowania procesami platformy; lewa kolumna, stała, pełna wysokość obszaru roboczego, obecna na stronie głównej i w każdym module | rozdz. 2.1 |
| Execution Loop Window | Okno pętli wykonawczej Koordynator ↔ Wykonawca, otwierane w kolumnie sąsiadującej z Chat Window; dostępne z poziomu strony głównej i z każdego modułu | rozdz. 2.2 |
| Warstwa widoczności | Jedna z czterech warstw ujawniania funkcjonalności: zawsze widoczna, widoczna na żądanie, rozwinięcia kontekstowe, funkcje eksperckie | rozdz. 3a |
| Znacznik kontekstowy | Lekki element paska kontekstu wskazujący środowisko, repozytorium, projekt, model albo wykonawcę; kliknięcie otwiera odpowiedni selektor, który zwija się po użyciu | rozdz. 3a.3 |
| Karta środowiska | Element strefy 1 strony głównej; jedna z czterech dużych kart wejścia do środowisk platformy, z tytułem złożonym krojem nagłówkowym, opisem trybu pracy i godłem środowiska | rozdz. 3.2 |
| Kafel komponentu własnego | Element strefy 2 strony głównej; jeden z czterech lżejszych kafli prowadzących do okna konfiguracji komponentu własnego, z etykietą zorientowaną na działanie twórcze | rozdz. 3.3 |
| Listwa ustawień | Element strefy 3 strony głównej; zwarty pasek narzędziowy prowadzący do okna konfiguracji oraz do funkcji globalnych Mobile i Always On Display | rozdz. 3.4 |
| Boczna nawigacja modułów | Stały, pionowy element interfejsu środowisk TalkIn, WorkSpace i CodeStudio, wyświetlający listę modułów dostępnych w danym środowisku zgodnie z macierzą dostępności modułów | rozdz. 5 |
| Panel orkiestracji | Boczna nawigacja środowiska MultitaskingAI, zastępująca listę modułów; złożona z sześciu sekcji: Zespoły, Role, Kolejki, Orkiestracja, Harmonogram i automatyki, Monitor procesu | rozdz. 6 |
| Karta sesji | Poziomy element paska kart w oknie środowiska, reprezentujący jedną, samodzielną przestrzeń roboczą z własnym układem, historią i kontekstem; mechanika analogiczna do kart w przeglądarce internetowej | rozdz. 7 |
| Przeładowanie przestrzeni roboczej | Mechanizm wymiany zestawu okien operacyjnych, kontekstu czatu oraz narzędzi w obrębie jednej karty sesji, następujący przy zmianie modułu | rozdz. 8.1 |
| Przełączanie środowiska | Zmiana całej przestrzeni roboczej (narzędzia, model nawigacji, historia sesji) przez powrót do strony głównej i wybór innej karty środowiska w strefie 1 | rozdz. 8.2 |

---

## 12. Wykazy normatywne i kryteria odbioru

Rozdział zamyka korpus dokumentu kompletem wykazów obowiązkowych wskazanych w [Standardzie redakcyjnym i językowym](../STANDARD-REDAKCYJNY-I-JEZYKOWY.md), rozdz. 4.2 — etykiety interfejsu, komunikaty, żetony, komponenty, skróty klawiszowe, stany kontrolek i punkty łamania — a kończy go zbiór kryteriów odbioru. Dosłowne brzmienia pochodzą z prototypu `design/05-okna/przeplyw/centrum-dowodzenia.html`, wartości wymiarów i barw z arkuszy `design/zasoby/zetony/zetony.css` i `design/zasoby/css/komponenty.css`, a nazwy komend z `budowa/shared/contract.json`. Wykaz komend kontraktu prowadzony jest osobno, w Załączniku D, i nie jest tu powtarzany.

### 12.1. Etykiety interfejsu

Brzmienia przejęte z prototypu; etykiety wymienione już w rozdziałach 3.2–3.4 nie są powtarzane.

| Element | Dosłowne brzmienie | Miejsce wystąpienia |
|---|---|---|
| Tytuł widoku | „Danaco Console — Strona główna · Centrum dowodzenia” | Belka tytułowa (rozdz. 2) |
| Ścieżka widoku | „Danaco Console › Centrum dowodzenia” | Pasek kontekstu (rozdz. 2) |
| Zdanie prowadzące strony głównej | „Wybierz środowisko poniżej — resztę podpowie samouczek po prawej.” | Nagłówek kolumny dominującej (rozdz. 3.1) |
| Motta kart środowisk | „9 modułów · rozmowa i wiedza” · „9 modułów · produktywność” · „8 modułów · programowanie” · „panel orkiestracji · 6 sekcji” | Karty środowisk, strefa 1 (rozdz. 3.2) |
| Opisy trybu pracy kart środowisk | „Wiedza, komunikacja i praca z treścią.” · „Produktywność, organizacja i realizacja projektów.” · „Programowanie — edytor, terminal, kontrola wersji.” · „Orkiestracja autonomicznej pracy ciągłej.” | Karty środowisk, strefa 1 (rozdz. 3.2) |
| Nazwy kafli komponentów własnych | „Automations” · „Agents” · „Workspace” · „Assistant” | Kafle, strefa 2 (rozdz. 3.3) |
| Opisy kafli komponentów własnych | „Automatyka wykonawcza” · „Agent autonomiczny” · „Projekt roboczy” · „Profil asystenta” | Kafle, strefa 2 (rozdz. 3.3) |
| Pozycje listwy ustawień | „Ustawienia i konfiguracja” · „Mobile” · „Always On Display” | Listwa, strefa 3 (rozdz. 3.4) |
| Podpowiedź pola wyszukiwania | „Szukaj w sesji, projekcie i środowisku…  Ctrl+K” | Pole wyszukiwania ramy okna (rozdz. 2.3) |
| Podpowiedź wyszukiwania w widoku | „Znajdź w bieżącym widoku…” | Pole „Znajdź w widoku” (rozdz. 2.3) |
| Nagłówek wyboru środowiska dla nowej sesji | „Nowa sesja — w którym środowisku” | Strefa pracy (rozdz. 2.3, 4.1) |
| Nagłówek wyboru środowiska dla nowego projektu | „Nowy projekt — w którym środowisku” | Strefa pracy (rozdz. 2.3) |
| Objaśnienie wyboru środowiska | „Otwarcie nowej sesji — najpierw wskazanie środowiska, w którym ma powstać.” | Strefa pracy (rozdz. 4.1) |
| Pozycje strefy pracy | „Nowa sesja” · „Nowy projekt” · „Otwórz projekt istniejący” · „Powtórz ostatnią sesję” · „Historia sesji” | Strefa pracy (rozdz. 2.3) |
| Opis pełnej historii sesji | „Pełny rejestr sesji konta: czynne, zakończone i archiwalne.” | Nakładka „Historia sesji” (rozdz. 9.1) |
| Licznik sesji czynnych | „4 czynne” · „Sesje w tle:” | Pasek stanu i szyna sesji (rozdz. 8.3) |
| Anatomia karty sesji | „Nazwa sesji i moduł” · „Dwa wiersze: nazwa sesji oraz moduł, w którym sesja pracuje.” · „Kropka sygnałowa przy sesji, w której trwa zadanie.” | Rozdz. 7.1 |
| Pozycje menu karty sesji | „Konfiguracja karty sesji” · „Kopiuj sesję” · „Eksportuj sesję…” · „Archiwizuj” | Menu `⋮` karty sesji (rozdz. 7.2) |
| Sekcje ustawień szybkich | „Motyw i gęstość widoku” · „Sesje i karty” · „Ustawienia aplikacji, środowiska, sesji i modeli.” | Listwa i okno konfiguracji (rozdz. 3.4, 9.1) |

### 12.2. Komunikaty

| Sytuacja | Dosłowna treść albo wzorzec | Rodzaj |
|---|---|---|
| Brak projektów na koncie | „Brak projektów. Projekt zakłada się w przestrzeni Workspace środowiska.” | Stan pusty |
| Pierwsze wejście do centrum dowodzenia | „Wybierz środowisko poniżej — resztę podpowie samouczek po prawej.” | Objaśnienie kontekstowe |
| Brak otwartych kart sesji | Pasmo kart puste, bez komunikatu — strona główna sama jest stanem wyjściowym (rozdz. 7.2) | Stan pusty |
| Brak sesji w historii | Stan pusty nakładki „Historia sesji” z podpowiedzią założenia nowej sesji (rozdz. 9.1) | Stan pusty |
| Moduł niedostępny w danym środowisku | Pozycja nie jest wyświetlana w bocznej nawigacji zgodnie z macierzą dostępności (Załącznik A) — brak komunikatu odmowy, brak pozycji wygaszonej | Stan (pominięcie) |
| Przeładowanie przestrzeni roboczej przy zmianie modułu | Wymiana zestawu okien w obrębie karty sesji zapowiedziana wprost, bez okna potwierdzenia; historia karty pozostaje (rozdz. 8.1) | Potwierdzenie |
| Zamknięcie karty sesji z zadaniem w toku | Dymek powiadomienia (`.dn-toast`) z opcją „Cofnij”; zadanie serwerowe trwa dalej jako sesja w tle (rozdz. 8.3) | Potwierdzenie |
| Współdzielenie kontekstu między kartami sesji | Komunikat kontekstowy o zakresie współdzielenia przy zmianie ustawienia; stan wyjściowy to izolacja karty (rozdz. 7.3) | Ostrzeżenie |
| Utrata połączenia z rdzeniem | Wskaźnik stanu w pasku stanu; strona główna pozostaje w pełni nawigowalna, karty sesji nie są zamykane (rozdz. 8.3) | Ostrzeżenie |
| Niepowodzenie wznowienia sesji z historii | Komunikat wskazujący przyczynę przy pozycji listy, bez zamykania nakładki (kod błędu `not_found` albo `conflict`, Załącznik D) | Błąd |

### 12.3. Żetony projektowe

| Miejsce zastosowania | Żeton `--dn-*` | Czego dotyczy |
|---|---|---|
| Tło strony głównej i tła stref | `--dn-tlo`, `--dn-panel`, `--dn-powierzchnia`, `--dn-powierzchnia-2` | Tło kolumny dominującej, kafli strefy 2 (`--dn-panel`) i kart strefy 1 (`--dn-powierzchnia`) — rozdz. 3.3 |
| Obrysy i uniesienia trzech stref | `--dn-obrys`, `--dn-obrys-mocny`, `--dn-cien-1`, `--dn-cien-2` | Spoczynek i wskazanie kursorem karty środowiska oraz kafla (rozdz. 3.2, 3.3) |
| Wstęga karty środowiska | `--dn-wym-wstega`, `--dn-kropka` | Wstęga górna 2px pojawiająca się przy wskazaniu kursorem i na karcie bieżącej (rozdz. 3.2) |
| Akcent sygnałowy | `--dn-sygnal`, `--dn-sygnal-500`, `--dn-sygnal-tlo`, `--dn-sygnal-wypelnienie` | Jedyny akcent chromatyczny platformy — pozycja bieżąca nawigacji, akcja główna strefy pracy (rozdz. 5.2) |
| Ognisko kontrolek | `--dn-fokus`, `--dn-fokus-cien`, `--dn-wym-fokus`, `--dn-wym-fokus-odsuniecie` | Pierścień ogniska kart, kafli, pozycji listwy i pozycji bocznej nawigacji |
| Typografia kart środowisk | `--dn-ff-naglowek`, `--dn-fs-3xl`, `--dn-fw-polgruba`, `--dn-ff-mono`, `--dn-fs-sm`, `--dn-tekst-3` | Tytuł, motto i opis karty środowiska (rozdz. 3.2) |
| Typografia kafli i listwy | `--dn-ff-bazowa`, `--dn-fs-base`, `--dn-tekst-2`, `--dn-lh-bazowy` | Nazwa i opis kafla, etykiety pozycji listwy (rozdz. 3.3, 3.4) |
| Zaokrąglenia | `--dn-r-xl`, `--dn-r-lg`, `--dn-r-md`, `--dn-r-sm`, `--dn-r-pill` | Karta środowiska (`--dn-r-xl`), kafel (`--dn-r-lg`), karta sesji (`--dn-r-md`), pozycja listwy (`--dn-r-sm`), plakietka licznika (`--dn-r-pill`) |
| Odstępy siatki trzech stref | `--dn-od-2`, `--dn-od-4`, `--dn-od-6`, `--dn-od-8`, `--dn-siatka-przerwa`, `--dn-odstep-sekcji` | Przerwy między kartami, kaflami i sekcjami stref (rozdz. 3.1) |
| Siatka i szerokość treści | `--dn-siatka-kolumny`, `--dn-tresc-max` | Dwanaście kolumn siatki i maksymalna szerokość kolumny dominującej (rozdz. 3.5) |
| Wymiary powierzchni nawigacyjnych | `--dn-wym-boczna`, `--dn-wym-pas-kart`, `--dn-wym-pas-komunikacji`, `--dn-wym-belka`, `--dn-wym-stan` | Boczna nawigacja i panel orkiestracji, pasmo kart sesji, kolumna Chat Window, belka tytułowa, pasek stanu (rozdz. 2, 5, 6, 7) |
| Wymiary kontrolek nawigacji | `--dn-wym-ikonowy`, `--dn-wym-ikona`, `--dn-wym-ikona-szyna`, `--dn-wym-kropka`, `--dn-wym-awatar` | Pozycje listwy, ikony bocznej nawigacji, kropka sygnałowa karty sesji, awatar profilu |
| Stany wskazania i wciśnięcia | `--dn-hover`, `--dn-wcisniecie` | Pozycje bocznej nawigacji, kart sesji i listwy ustawień |
| Czasy i krzywa przejść | `--dn-czas-1`, `--dn-czas-2`, `--dn-czas-3`, `--dn-ease`, `--dn-czas-tetno` | Mikroreakcje kart i kafli, przeładowanie przestrzeni roboczej, tętno kropki sesji pracującej |
| Warstwy widoczności elementów | `--dn-z-podloga`, `--dn-z-boczna`, `--dn-z-nakladka`, `--dn-z-modal`, `--dn-z-tooltip`, `--dn-z-powiadomienie`, `--dn-z-centrum-polecen`, `--dn-z-aod` | Kolejność nakładania: strefy, boczna nawigacja, panele popover, nakładki, paleta poleceń i warstwa Always On Display (rozdz. 3a) |
| Stany semantyczne liczników i plakietek | `--dn-sukces-tekst`, `--dn-ostrzezenie-tekst`, `--dn-blad-tekst`, `--dn-informacja-tekst` wraz z odpowiadającymi im żetonami `-tlo` i `-obrys` | Plakietka sesji czynnych, ostrzeżenie o współdzieleniu kontekstu, błąd wznowienia sesji |

### 12.4. Komponenty interfejsu

Klasy zweryfikowane wobec `design/zasoby/css/komponenty.css`.

| Klasa `.dn-*` | Rola | Modyfikatory | Stany |
|---|---|---|---|
| `.dn-karta-srodowiska` | Karta wejścia do środowiska, strefa 1 | elementy pochodne `-godlo`, `-tytul`, `-motto`, `-opis` | spoczynek · wskazanie kursorem · bieżąca (`aria-current`) · ognisko |
| `.dn-kafel` | Kafel komponentu własnego, strefa 2 | elementy pochodne `-ikona`, `-etykieta`, `-opis` | spoczynek · wskazanie kursorem · ognisko |
| `.dn-listwa` | Listwa ustawień, strefa 3 | pozycja `-pozycja` | spoczynek · wskazanie kursorem · ognisko |
| `.dn-boczna` | Boczna nawigacja modułów oraz panel orkiestracji | nagłówek `-naglowek`, pozycja `-pozycja` | rozwinięta · ikonowa |
| `.dn-boczna-pozycja` | Pozycja modułu albo sekcji panelu | — | spoczynek · wskazanie kursorem · bieżąca (`aria-current='page'`) · ognisko |
| `.dn-karty-sesji` | Pasmo kart sesji | — | wypełnione · puste · przewijane · przycięte |
| `.dn-karta-sesji` | Pojedyncza karta sesji | `-zamknij` (znacznik zamknięcia) | aktywna · w tle · z pracą w toku · ognisko |
| `.dn-kropka` | Kropka sygnałowa sesji pracującej | `--tetno` · `--sygnal` · `--sukces` · `--ostrzezenie` · `--blad` · `--neutralna` | aktywna · bezczynna |
| `.dn-plakietka` | Plakietka licznika sesji i statusu | `--sygnal` · `--sukces` · `--ostrzezenie` · `--blad` · `--informacja` | spoczynek |
| `.dn-szukaj` | Pole wyszukiwania w środowiskach | — | pusty · wypełniony · ognisko |
| `.dn-pasek` | Pasek kontekstu strony głównej | `-godlo` · `-logotyp` · `-prawa` · `-szukaj` | spoczynek |
| `.dn-modal` | Nakładka „Historia sesji”, okno konfiguracji | `--szeroki` | otwarta · zamykana |
| `.dn-toast` | Dymek powiadomienia z opcją „Cofnij” | `--sukces` · `--ostrzezenie` · `--blad` · `--informacja` | widoczny · wygaszany |
| `.dn-pusty-stan` | Stan pusty listy projektów i historii sesji | `-tytul` · `-opis` | spoczynek |
| `.dn-tooltip` | Etykieta pozycji zwiniętej bocznej nawigacji | — | ukryta · widoczna |
| `.dn-prompt` | Pole polecenia Chat Window na stronie głównej | `-obszar` · `-grot` | pusty · wypełniony · ognisko · wysyłanie |
| `.dn-wpis` | Wpis w Chat Window i w Execution Loop Window | `--czlowiek` · `--inteligencja` · `--system` · `--pracuje` | zakończony · strumieniowany |
| `.dn-aod` | Awatar i powierzchnia Always On Display | `-rdzen` · `-tresc` | zwinięta · rozwinięta |
| `.dn-spinner` | Wskaźnik ładowania przy przeładowaniu przestrzeni roboczej | — | aktywny |

### 12.5. Skróty klawiszowe

Wykaz podaje **wartości domyślne**. Zgodnie z rozstrzygnięciem platformowym żaden skrót nie jest ustalony na sztywno — wszystkie nadaje Operator samodzielnie w oknie konfiguracji i ustawień (rozdz. 9.1), a „Przywróć domyślne” przywraca stan z tej tabeli. Wszystkie kombinacje poniżej należą do ramy okna i są wspólne dla całej platformy — [Rama okna aplikacji](rama-okna.md), rozdz. 16 i 22.4; strona główna ich nie przedefiniowuje. Prototypy `design/05-okna/` są wzorcem wyglądu i układu, nie źródłem normatywnym przypisań klawiszowych.

| Kombinacja | Działanie | Zasięg | Kolizje |
|---|---|---|---|
| `Ctrl+Home` | Powrót do centrum dowodzenia — strony głównej (rozdz. 8.2) | globalny | — |
| `Ctrl+E` | Przedsionek środowiska bieżącego (rozdz. 4.1) | globalny | — |
| `Ctrl+N` | Nowa sesja — z wyborem środowiska (rozdz. 2.3, 4.1) | globalny | — |
| `Ctrl+K` | Paleta poleceń i wyszukiwanie w środowiskach (rozdz. 2.3) | globalny | — |
| `Ctrl+F` | Znajdź w bieżącym widoku | widok | — |
| `Ctrl+Tab` | Następna karta sesji (rozdz. 7.2) | powłoka | — |
| `Ctrl+Shift+Tab` | Poprzednia karta sesji (rozdz. 7.2) | powłoka | — |
| `Ctrl+B` | Szyna sesji — przełącznik widoczności | globalny | — |
| `Ctrl+Shift+B` | Szyna nawigacji — przełącznik widoczności (rozdz. 5) | globalny | — |
| `Ctrl+Shift+F` | Tryb skupienia | globalny | — |
| `Ctrl+,` | Okno konfiguracji i ustawień (rozdz. 9.1) | globalny | — |
| `Ctrl+/` | Ściągawka skrótów aktywnych | globalny | — |
| `Ctrl+O` | Otwórz projekt istniejący (rozdz. 2.3) | globalny | dzieli kombinację z pozycją „Otwórz projekt…” menu aplikacji — [Rama okna aplikacji](rama-okna.md), rozdz. 22.4 |

Czynności nawigacyjne swoiste dla strony głównej — wejście do wskazanej karty środowiska, otwarcie kafla komponentu własnego, przejście do pozycji listwy ustawień — **nie mają wartości fabrycznej**; kombinację nadaje im Operator w edytorze skrótów, a pole przypisania pozostaje puste do pierwszego nadania. Edytor odmawia zapisania kombinacji już zajętej i wskazuje czynność, która ją zajmuje.

### 12.6. Stany kontrolek

Dziewięć kolumn za standardem rozdz. 4.2. Kolumna „nieaktywny” niemal wszędzie brzmi „nie dotyczy”, ponieważ platforma nie stawia twardych blokad (rozdz. 10) — moduł niedostępny w danym środowisku nie jest wygaszany, lecz pomijany w wykazie pozycji (rozdz. 5.1).

| Kontrolka | Spoczynek | Wskazanie kursorem | Wciśnięcie | Ognisko | Nieaktywny | Ładowanie | Pusty | Błąd |
|---|---|---|---|---|---|---|---|---|
| Karta środowiska (`.dn-karta-srodowiska`) | Obrys `--dn-obrys`, tło `--dn-powierzchnia`, wstęga ukryta | Obrys `--dn-obrys-mocny`, cień `--dn-cien-2`, uniesienie 2px, wstęga górna `--dn-kropka` | Wejście do środowiska następuje na puszczeniu — bez odrębnego stanu wizualnego | Pierścień `--dn-fokus` | nie dotyczy — cztery środowiska są zawsze dostępne | nie dotyczy — przejście do przedsionka jest lokalne | nie dotyczy — karta niesie treść stałą | nie dotyczy — błąd wejścia zgłasza dymek powiadomienia, nie karta |
| Kafel komponentu własnego (`.dn-kafel`) | Tło `--dn-panel`, obrys `--dn-obrys` | Obrys `--dn-obrys-mocny`, tło `--dn-powierzchnia`, bez uniesienia i bez wstęgi | Otwarcie okna konfiguracji na puszczeniu | Pierścień `--dn-fokus` | nie dotyczy — cztery rodzaje komponentów zawsze dostępne | Spinner w panelu kafla przy wczytywaniu listy komponentów zapisanych | „Brak projektów. Projekt zakłada się w przestrzeni Workspace środowiska.” — panel kafla bez zapisanych pozycji | Komunikat w panelu kafla przy niepowodzeniu wczytania listy |
| Pozycja listwy ustawień (`.dn-listwa-pozycja`) | Tekst `--dn-tekst-2` | Tło `--dn-hover`, tekst `--dn-tekst` | Tło `--dn-wcisniecie` | Pierścień `--dn-fokus` | nie dotyczy — trzy pozycje stale dostępne | nie dotyczy | nie dotyczy | nie dotyczy — błąd zgłasza otwierane okno, nie pozycja |
| Pozycja bocznej nawigacji (`.dn-boczna-pozycja`) | Tekst `--dn-tekst-2`, tło przezroczyste | Tło `--dn-hover` | Przeładowanie przestrzeni roboczej karty | Pierścień `--dn-fokus` | nie dotyczy — moduł niedostępny jest pomijany, nie wygaszany (rozdz. 5.1) | Spinner przy pozycji w trakcie przeładowania przestrzeni roboczej (rozdz. 8.1) | nie dotyczy — wykaz pozycji nigdy nie jest pusty | Komunikat kontekstowy przy pozycji, gdy moduł nie odpowiada |
| Karta sesji (`.dn-karta-sesji`) | Tło `--dn-panel`, nazwa sesji i moduł w dwóch wierszach | Ujawnienie znacznika zamknięcia (`.dn-karta-sesji-zamknij`) | Przełączenie ogniska na kartę | Pierścień `--dn-fokus` | nie dotyczy | Kropka sygnałowa (`.dn-kropka--tetno`) przy sesji, w której trwa zadanie | Pasmo bez kart — puste, bez komunikatu (rozdz. 7.2) | Plakietka `--dn-blad-tlo` przy karcie sesji zakończonej niepowodzeniem |
| Pole wyszukiwania w środowiskach (`.dn-szukaj`) | Obrys `--dn-obrys` | Obrys `--dn-sygnal-obrys` | nie dotyczy — pole nie jest przyciskiem | Pierścień `--dn-fokus` | nie dotyczy | Wskaźnik pracy indeksu w polu | Podpowiedź „Szukaj w sesji, projekcie i środowisku…  Ctrl+K” | Obrys `--dn-blad-obrys` z komunikatem pod polem |
| Pozycja nakładki „Historia sesji” (`.dn-modal`) | Wiersz sesji z nazwą, środowiskiem i czasem | Tło `--dn-hover` | Wznowienie sesji | Pierścień `--dn-fokus` | nie dotyczy — sesja archiwalna pozostaje wznawialna | Spinner w wierszu w trakcie wznawiania | Stan pusty rejestru sesji (`.dn-pusty-stan`) | Komunikat przy wierszu — `not_found` albo `conflict` (Załącznik D) |
| Pole polecenia Chat Window (`.dn-prompt`) | Obrys `--dn-obrys` | Obrys `--dn-sygnal-obrys` | nie dotyczy | Pierścień `--dn-fokus` | nie dotyczy — zasada zero blokad | Ikona wysłania ustępuje wskaźnikowi strumieniowania | Podpowiedź pola polecenia | Obrys `--dn-blad-obrys`, komunikat pod polem |
| Znacznik kontekstowy paska (rozdz. 3a.3) | Zwarty napis z wartością bieżącą | Tło `--dn-hover` | Rozwinięcie selektora | Pierścień `--dn-fokus` | nie dotyczy | nie dotyczy | Brak wartości ustalonej — napis podpowiadający wybór | nie dotyczy |

### 12.7. Punkty łamania

Progi z `design/zasoby/zetony/zetony.css`, rodzina `--dn-bp-*`.

| Próg szerokości | Zachowanie układu |
|---|---|
| `--dn-bp-w1` — 640px | Widok mobilny: Chat Window i kolumna dominująca układają się jedna pod drugą, karty środowisk w jednej kolumnie, listwa ustawień zwija się do menu; dostęp opisuje [Mobile](../funkcje-globalne/mobile.md) |
| `--dn-bp-w2` — 960px | Boczna nawigacja modułów i panel orkiestracji zwijają się do samych ikon z etykietą ujawnianą wskazaniem kursorem (rozdz. 5.1, 6.3); karty środowisk w dwóch kolumnach |
| `--dn-bp-w3` — 1280px | Pełny układ strony głównej: Chat Window w lewej kolumnie i trzy strefy w kolumnie dominującej widoczne jednocześnie, karty środowisk w jednym rzędzie (rozdz. 3.5) |
| `--dn-bp-w4` — 1600px | Szerokie biurko: Chat Window i Execution Loop Window rozwinięte równocześnie obok kolumny dominującej, bez zwężania stref (rozdz. 2.2) |

Szerokość kolumny dominującej ograniczona jest żetonem `--dn-tresc-max` (1200px), a jej podział opiera się na siatce `--dn-siatka-kolumny` (dwanaście kolumn) z przerwą `--dn-siatka-przerwa`.

### 12.8. Kryteria odbioru

| Warunek sprawdzalny | Sposób sprawdzenia |
|---|---|
| Strona główna prezentuje trzy strefy o malejącej masie wizualnej, w formach karty–kafel–listwa | Przegląd strony głównej wobec rozdz. 3.1–3.4; potwierdzenie odmiennego kroju, tła i zaokrąglenia karty i kafla |
| Chat Window jest obecne na stronie głównej w lewej kolumnie o pełnej wysokości | Otwarcie strony głównej; pomiar wysokości kolumny (rozdz. 2.1) |
| Cztery karty środowisk noszą motta z dokładnie trzech czasowników rozkazujących | Odczyt czterech mott wobec tabeli rozdz. 3.2 |
| Wstęga górna karty środowiska pojawia się wyłącznie przy wskazaniu kursorem i na karcie bieżącej | Wskazanie kursorem karty oraz odczyt karty z `aria-current='true'` (rozdz. 3.2) |
| Kliknięcie kafla otwiera okno konfiguracji komponentu, nie przestrzeń roboczą środowiska | Kliknięcie każdego z czterech kafli (rozdz. 3.3) |
| Boczna nawigacja pokazuje wyłącznie moduły dostępne w danym środowisku, bez pozycji wygaszonych | Przegląd bocznej nawigacji czterech środowisk wobec macierzy Załącznika A |
| Panel orkiestracji zastępuje listę modułów sześcioma sekcjami w środowisku MultitaskingAI | Wejście do środowiska MultitaskingAI; przegląd panelu wobec rozdz. 6.2 i Załącznika B |
| Interfejs prezentuje wyłącznie warstwę 1, a warstwy 2–4 są ujawniane wyzwalaczami | Przegląd strony głównej wobec rozdz. 3a.2; potwierdzenie obecności znaczników kontekstowych, `⋮` i `☰` |
| Każda funkcja warstw 2–4 jest osiągalna co najwyżej jednym kliknięciem od wyzwalacza | Przegląd wyzwalaczy wobec zasady rozdz. 3a.4 |
| Zmiana modułu przeładowuje przestrzeń roboczą karty, zachowując historię i kontekst karty | Zmiana modułu w karcie sesji; porównanie z zestawieniem rozdz. 8.4 |
| Przełączenie środowiska prowadzi przez powrót do centrum dowodzenia i nie zamyka sesji w tle | Przełączenie środowiska; przegląd szyny sesji (rozdz. 8.2, 8.3) |
| Karta sesji pozostaje domyślnie odizolowana, a współdzielenie jest świadomym ustawieniem | Otwarcie dwóch kart sesji; potwierdzenie odrębności historii i kontekstu (rozdz. 7.3) |
| Żadna kontrolka nawigacyjna nie jest trwale nieaktywna | Przegląd kontrolek z rozdz. 12.4 pod kątem atrybutu `disabled`; zero wystąpień poza stanem „ładowanie” |
| Każdy skrót jest zmienialny, a czynność bez wartości fabrycznej czeka na nadanie | Otwarcie edytora skrótów w oknie konfiguracji; próba przypisania kombinacji zajętej (rozdz. 12.5) |
| Układ odpowiada czterem progom łamania | Zmiana szerokości okna kolejno przez 640px, 960px, 1280px i 1600px; obserwacja zachowań z rozdz. 12.7 |
| Każdy żeton i każda klasa przywołane w dokumencie istnieją w arkuszach źródłowych | Zestawienie nazw z rozdz. 12.3 i 12.4 z `design/zasoby/zetony/zetony.css` i `design/zasoby/css/komponenty.css` |
| Każda komenda przywołana w korpusie występuje w wykazie Załącznika D | Zestawienie nazw komend z korpusu z Załącznikiem D i z `budowa/shared/contract.json` |

---

## Załącznik A. Macierz dostępności modułów w bocznej nawigacji

Zestawienie powtarza macierz dostępności modułów z Koncepcji platformy (rozdz. 10) w kontekście bocznej nawigacji: określa, w której bocznej nawigacji środowiska modułowego dany moduł jest widoczny.

| Moduł | TalkIn | WorkSpace | CodeStudio |
|---|---|---|---|
| Studio | TAK | TAK | NIE |
| Workspace | TAK | TAK | TAK |
| Automations | — | — | — |
| Browser | TAK | TAK | NIE |
| Research | TAK | TAK | NIE |
| Library | TAK | TAK | NIE |
| Translate | TAK | NIE | NIE |
| Roundtable | TAK | TAK | TAK |
| Design | NIE | TAK | TAK |
| Assistant | TAK | NIE | NIE |
| Terminal | NIE | NIE | TAK |
| Developer | NIE | NIE | TAK |
| Diagnostics | NIE | NIE | TAK |
| Apps | NIE | TAK | TAK |
| Agents | TAK | TAK | TAK |
| **Liczba pozycji widocznych** | **9** | **9** | **8** |

Automations nie ma pozycji w żadnej z trzech bocznych nawigacji — moduł jest konfigurowany wyłącznie ze strony głównej (strefa 2), a gotowe automatyki trafiają do sesji jako komponenty własne. Środowisko MultitaskingAI nie ma bocznej nawigacji modułów, dlatego nie występuje jako kolumna macierzy — jego odpowiednik opisuje Załącznik B.

---

## Załącznik B. Panel orkiestracji — zestawienie sekcji

Powtórzenie zestawienia sześciu sekcji panelu orkiestracji (Koncepcja platformy, rozdz. 13.8) w kontekście bocznej nawigacji środowiska MultitaskingAI.

| Sekcja | Zawartość | Powiązane okna | Warstwa | Sposób wywołania |
|---|---|---|---|---|
| Zespoły | Zapisane konfiguracje zespołu — szablony ról, powiązań i kolejek | — | 1 | Pozycja panelu widoczna bez interakcji |
| Role | Executor 1, Executor 2, Coordinator, Executor 3 / Validator; przypisanie agentów; Subagent Network | Agent Builder, Permissions Center | 1 | Pozycja panelu widoczna bez interakcji |
| Kolejki | Kolejki globalne, lokalne, modeli, agentów, projektów | Queue Manager | 1 | Pozycja panelu widoczna bez interakcji |
| Orkiestracja | Zależności między modelami, agentami, zadaniami, kolejkami, automatyzacjami, projektami | Orchestrator | 1 | Pozycja panelu widoczna bez interakcji |
| Harmonogram i automatyki | Harmonogram pracy ciągłej; wpięte automatyki | Scheduler, Execution Monitor | 1 | Pozycja panelu widoczna bez interakcji |
| Monitor procesu | Podgląd przebiegu pętli, statusy, hierarchia decyzji | Always On Display, Mobile | 1 | Pozycja panelu widoczna bez interakcji |

---

## Załącznik C. Scenariusz nawigacyjny — pełny przebieg

Poniższy scenariusz ilustruje pełny przebieg nawigacji opisany w rozdziałach 4, 7 i 8, na przykładzie użytkownika prowadzącego równolegle pracę redakcyjną i pracę badawczą, a następnie sprawdzającego kod aplikacji rozwijanej równolegle.

| Krok | Akcja użytkownika | Co się dzieje | Karty aktywne | Rozdział |
|---|---|---|---|---|
| 1 | Otwiera centrum dowodzenia i klika kartę środowiska TalkIn (strefa 1) | Strona główna zostaje zastąpiona oknem środowiska TalkIn | — | rozdz. 4.1 |
| 2 | W bocznej nawigacji TalkIn wybiera moduł Studio | Pierwsza karta sesji ładuje zestaw okien: Chat Window, Studio Editor, Tools Panel, Diff/Grep Panel, Session Repository, Preview Window | TalkIn: Studio | rozdz. 4.2, 5.3 |
| 3 | Otwiera nową kartę sesji („+”) i w jej bocznej nawigacji wybiera moduł Research | Druga karta ładuje się z własnym, odrębnym zestawem okien i odrębnym kontekstem; pierwsza karta (Studio) pozostaje niezmieniona w tle | TalkIn: Studio, Research | rozdz. 7.2, 7.3 |
| 4 | Przełącza się między kartami, klikając ich tytuły | Za każdym razem widoczny jest inny, dedykowany układ pracy | TalkIn: Studio, Research | rozdz. 7.2 |
| 5 | Wraca do strony głównej i klika kartę środowiska CodeStudio | Obie karty sesji otwarte w TalkIn (Studio, Research) pozostają aktywne w tle i odtworzą pełny stan przy powrocie | TalkIn w tle: Studio, Research | rozdz. 8.2, 8.3 |
| 6 | W środowisku CodeStudio wybiera moduł Developer i pracuje nad kodem w nowej karcie sesji | Praca nad kodem w karcie właściwej środowisku CodeStudio | CodeStudio: Developer; TalkIn w tle: Studio, Research | rozdz. 8.3 |
| 7 | Po zakończeniu pracy wraca do strony głównej i ponownie wybiera kartę TalkIn | Obie wcześniej otwarte karty (Studio, Research) są widoczne dokładnie w takim stanie, w jakim je pozostawił | TalkIn: Studio, Research | rozdz. 8.3 |


---

## Załącznik D. Pełny wykaz komend kontraktu strony głównej i nawigacji

Wykaz wygenerowany wprost z `budowa/shared/contract.json` — jedynego źródła nazw komend, pól ładunków i zdarzeń. Każda pozycja jest wiążąca dla warstwy klienckiej i serwerowej; deweloper nie dodaje ani nie zmienia nazw poza tym wykazem.

Cztery obszary tego wykazu — `home`, `environment`, `session`, `window` — nie wyczerpują wszystkich
sześćdziesięciu ośmiu obszarów kontraktu (`budowa/shared/contract.json`, pole `obszary`); są
wyłącznie obszarami, na których operuje strona główna i mechanizm nawigacji między środowiskiem,
modułem i oknem opisane w rozdziałach 2–8 niniejszego dokumentu. Obszary właściwe treści modułów
poszczególnych (np. `workspace`, `automation`, `browser`) należą do dokumentów tych modułów, poza
zakresem niniejszego opracowania.

### Obszar `home` — 1 komenda

| Komenda | Przeznaczenie | Pola żądania | Pola wyniku |
|---|---|---|---|
| `home.enter` | Wejście na stronę główną; zwraca karty środowisk i sesje czynne konta | `clientId:string` (wym) | `environments:Environment[]` (wym)<br>`sessions:Session[]` (wym)<br>`focusedSessionId:string` (opc)<br>`presence:SessionPresence[]` (opc) |

### Obszar `environment` — 2 komendy

| Komenda | Przeznaczenie | Pola żądania | Pola wyniku |
|---|---|---|---|
| `environment.list` | Zwraca środowiska platformy | `includeModules:bool` (opc) | `environments:Environment[]` (wym) |
| `environment.enter` | Wchodzi do środowiska; zwraca jego nawigację i karty sesji | `environmentId:string` (wym)<br>`clientId:string` (wym)<br>`sessionId:string` (opc) | `environment:Environment` (wym)<br>`modules:Module[]` (wym)<br>`sessions:Session[]` (wym)<br>`focusedSessionId:string` (opc) |

### Obszar `session` — 19 komend (rozdz. 7, 8)

| Komenda | Przeznaczenie | Pola żądania | Pola wyniku |
|---|---|---|---|
| `session.create` | Zakłada sesję | `title:string` (opc)<br>`projectId:string` (opc)<br>`metadata:json` (opc) | `session:Session` (wym) |
| `session.list` | Zwraca listę sesji | `status:SessionStatus` (opc)<br>`limit:int` (opc)<br>`offset:int` (opc)<br>`includePresence:bool` (opc) | `sessions:Session[]` (wym)<br>`total:int` (wym)<br>`presence:SessionPresence[]` (opc) |
| `session.focus` | Przenosi ognisko na wskazaną kartę sesji i opcjonalnie na okno w jej wnętrzu. Asystent przestawia ognisko Operatorowi polem `targetClientId` — inaczej jego posunięcia nie byłyby widoczne na ekranie | `sessionId:string` (wym)<br>`clientId:string` (wym)<br>`windowId:string` (opc)<br>`targetClientId:string` (opc) | `sessionId:string` (wym)<br>`windowId:string` (opc)<br>`previousSessionId:string` (opc)<br>`focusedAt:int64` (wym) |
| `session.bind` | Wiąże bieżące połączenie z sesją trwającą na rdzeniu; odtwarza jej okna | `sessionId:string` (wym)<br>`clientId:string` (wym)<br>`windowIds:string[]` (opc) | `session:Session` (wym)<br>`windows:Window[]` (wym)<br>`bound:bool` (wym)<br>`resumed:bool` (wym) |
| `session.open` | Otwiera sesję wraz z jej oknami | `sessionId:string` (wym) | `session:Session` (wym)<br>`windows:Window[]` (wym) |
| `session.close` | Zamyka sesję | `sessionId:string` (wym) | `session:Session` (wym) |
| `session.delete` | Usuwa trwale wskazane sesje wraz z całym ich zapisem; jedyna droga utraty danych sesji | `sessionIds:string[]` (wym)<br>`confirm:bool` (wym) | `deletedIds:string[]` (wym)<br>`deletedCount:int` (wym)<br>`missingIds:string[]` (opc) |
| `session.rename` | Zmienia nazwę sesji w historii | `sessionId:string` (wym)<br>`title:string` (wym) | `session:Session` (wym) |
| `session.copy` | Kopiuje sesję wraz z jej zapisem; kopia jest osobnym bytem historii | `sessionId:string` (wym)<br>`title:string` (opc) | `session:Session` (wym)<br>`copiedMessages:int` (wym) |
| `session.project.set` | Przenosi sesję do projektu; projekt wskazany albo zakładany nowy | `sessionIds:string[]` (wym)<br>`projectId:string` (opc)<br>`projectName:string` (opc) | `projectId:string` (wym)<br>`movedIds:string[]` (wym) |
| `session.project.clear` | Wyjmuje sesję z projektu; sesja zostaje w historii bez projektu | `sessionIds:string[]` (wym) | `clearedIds:string[]` (wym) |
| `session.archive` | Przenosi sesje do archiwum; zapis zostaje w całości, sesja znika z historii bieżącej | `sessionIds:string[]` (wym) | `archivedIds:string[]` (wym) |
| `session.restore` | Przywraca sesje z archiwum do historii bieżącej | `sessionIds:string[]` (wym) | `restoredIds:string[]` (wym) |
| `session.archive.list` | Zwraca sesje archiwum; wgląd z okna ustawień | `offset:int` (opc)<br>`limit:int` (opc) | `sessions:Session[]` (wym)<br>`total:int` (wym) |
| `session.resume` | Wznawia zamkniętą sesję wraz z jej oknami | `sessionId:string` (wym) | `session:Session` (wym)<br>`windows:Window[]` (wym) |
| `session.stop` | Zatrzymuje tury biegnące we wszystkich oknach sesji | `sessionId:string` (wym) | `sessionId:string` (wym)<br>`stoppedWindowIds:string[]` (wym) |
| `session.tool.attach` | Dokłada narzędzie albo umiejętność do sesji na czas jej trwania. Definicji eksperta nie zmienia: dołożenie żyje w stanie sesji, przeżywa rozłączenie klienta i kończy się wraz z sesją albo z usunięciem rozmowy | `sessionId:string` (wym)<br>`toolName:string` (wym)<br>`source:SessionToolSource` (opc) | `tool:SessionTool` (wym)<br>`tools:SessionTool[]` (wym)<br>`alreadyAttached:bool` (wym) |
| `session.tool.detach` | Zdejmuje dołożenie z sesji. Zestaw wraca do podstawy z definicji eksperta; sama definicja pozostaje niezmieniona | `sessionId:string` (wym)<br>`toolName:string` (wym) | `detached:bool` (wym)<br>`tools:SessionTool[]` (wym) |
| `session.tool.list` | Zwraca narzędzia dołożone do sesji. Zestaw narzędzi tury to definicja eksperta wraz z tymi dołożeniami — bez tego odczytu druga połowa zestawu byłaby niewidoczna | `sessionId:string` (wym) | `tools:SessionTool[]` (wym)<br>`total:int` (wym) |

**Zdarzenia obszaru `session` — 4:**

| Zdarzenie | Przeznaczenie | Pola |
|---|---|---|
| `session.changed` | Zmiana sesji | — |
| `session.focus.changed` | Zmiana ogniska karty sesji; nośnik synchronizacji ogniska między urządzeniami konta | — |
| `session.tool.attached` | Narzędzie dołożone do sesji. Skoro model dostaje narzędzie, którego nie miał, Operator musi to zobaczyć — każde posunięcie widać na ekranie | — |
| `session.tool.detached` | Dołożenie zdjęte z sesji. Widoczność obowiązuje w obie strony: zestaw, który urósł na oczach Operatora, nie może skurczyć się po cichu — bez tego zdarzenia sąsiednie okno pokazywałoby narzędzie, którego model już nie ma | — |

### Obszar `window` — 7 komend (rozdz. 4, 8)

| Komenda | Przeznaczenie | Pola żądania | Pola wyniku |
|---|---|---|---|
| `window.create` | Zakłada okno komunikacji w sesji | `sessionId:string` (wym)<br>`moduleId:string` (wym)<br>`modelChannelId:string` (wym)<br>`workingDirs:string[]` (wym)<br>`executionEnv:ExecutionEnv` (wym)<br>`permissionMode:PermissionMode` (wym)<br>`windowRole:WindowRole` (wym)<br>`coordinatorWindowId:string` (opc)<br>`title:string` (opc) | `window:Window` (wym) |
| `window.list` | Zwraca okna komunikacji | `sessionId:string` (opc)<br>`status:WindowStatus` (opc) | `windows:Window[]` (wym) |
| `window.state.get` | Zwraca stan okna komunikacji: parametry wykonania, stan procesu i konfigurację efektywną | `windowId:string` (wym)<br>`includeConfig:bool` (opc)<br>`includeAccess:bool` (opc) | `window:Window` (wym)<br>`processStatus:ProgressStatus` (wym)<br>`messageCount:int` (wym)<br>`lastMessageId:string` (opc)<br>`streaming:bool` (wym)<br>`config:ConfigEntry[]` (opc)<br>`accessGrants:AccessGrant[]` (opc)<br>`loop:LoopState` (opc) |
| `window.update` | Zmienia parametry okna komunikacji | `windowId:string` (wym)<br>`moduleId:string` (opc)<br>`modelChannelId:string` (opc)<br>`agentId:string` (opc)<br>`workingDirs:string[]` (opc)<br>`executionEnv:ExecutionEnv` (opc)<br>`permissionMode:PermissionMode` (opc)<br>`windowRole:WindowRole` (opc)<br>`coordinatorWindowId:string` (opc)<br>`title:string` (opc) | `window:Window` (wym) |
| `window.close` | Zamyka okno komunikacji | `windowId:string` (wym) | `window:Window` (wym) |
| `window.action` | Wykonuje akcję panelu akcji okna operacyjnego; wspólną każdemu oknu | `windowId:string` (wym)<br>`actionId:string` (wym)<br>`parameters:json` (opc) | `windowId:string` (wym)<br>`actionId:string` (wym)<br>`result:json` (opc) |
| `window.handoff` | Przekazuje zlecenie z okna koordynatora do okna wykonawcy; zakłada pozycję kolejki | `sessionId:string` (wym)<br>`fromWindowId:string` (wym)<br>`toWindowId:string` (wym)<br>`instruction:string` (wym)<br>`bundle:ContextBundle` (opc) | `window:Window` (wym)<br>`queueItemId:string` (wym) |

**Zdarzenia obszaru `window` — 2:**

| Zdarzenie | Przeznaczenie | Pola |
|---|---|---|
| `window.changed` | Zmiana okna komunikacji | — |
| `window.state.changed` | Zmiana stanu okna operacyjnego; wspólna każdemu oknu | — |


Razem w wykazie: **29 komend** (1 + 2 + 19 + 7) z 4 obszarów kontraktu, plus sześć zdarzeń zasilających widok (4 obszaru `session`, 2 obszaru `window`).

---

*Koniec dokumentu. Strona główna i nawigacja — Specyfikacja docelowa, wersja 2.0, 2026-08-20.*

---
*Danaco Console — Platforma AI Workspace OS · v2.0 · status Deweloperski*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz.*
*Warunki korzystania: [Licencja produktu](../LICENSE.md). Kontakt: support@danaco-group.pl*
