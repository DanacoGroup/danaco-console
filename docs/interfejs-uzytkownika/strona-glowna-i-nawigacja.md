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
| **Data** | 2026-08-06 |

**Informacje szczegółowe dokumentu:**

| | |
|---|---|
| **Opracowanie** | Danaco Console |
| **Źródło** | Koncepcja platformy; Architektura |

Dokument stanowi szczegółową specyfikację strony głównej platformy Danaco Console (centrum dowodzenia) oraz nawigacji między jej elementami: pionowego układu kolumnowego strony głównej z Chat Window w lewej kolumnie, trzech stref wyboru wraz z formą ich prezentacji, warstw widoczności funkcjonalności, przepływu nawigacji strona główna → środowisko → moduł → okno, mechaniki kart sesji, mechanizmu przełączania i przeładowania przestrzeni roboczej, bocznej nawigacji modułów w środowiskach modułowych oraz panelu orkiestracji jako jej odpowiednika w środowisku MultitaskingAI. Dokument rozwija ustalenia zawarte w Koncepcji platformy (rozdziały 5, 7, 8, 9, 10 i 13) oraz w dokumencie Architektura (rozdziały 5, 7, 12 i 13).

---

## Spis treści

1. [Cel i miejsce dokumentu](#1-cel-i-miejsce-dokumentu)
2. [Strona główna jako centrum dowodzenia](#2-strona-główna-jako-centrum-dowodzenia)
3. [Trzy strefy strony głównej](#3-trzy-strefy-strony-głównej)
3a. [Warstwy widoczności w nawigacji](#3a-warstwy-widoczności-w-nawigacji)
4. [Przepływ nawigacji: strona główna → środowisko → moduł → okno](#4-przepływ-nawigacji-strona-główna--środowisko--moduł--okno)
5. [Boczna nawigacja modułów](#5-boczna-nawigacja-modułów)
6. [Panel orkiestracji — boczna nawigacja środowiska MultitaskingAI](#6-panel-orkiestracji--boczna-nawigacja-środowiska-multitaskingai)
7. [Karty sesji — mechanika](#7-karty-sesji--mechanika)
8. [Przełączanie i przeładowanie przestrzeni roboczej](#8-przełączanie-i-przeładowanie-przestrzeni-roboczej)
9. [Wejście do ustawień i funkcji globalnych](#9-wejście-do-ustawień-i-funkcji-globalnych)
10. [Zgodność z zasadami nadrzędnymi platformy](#10-zgodność-z-zasadami-nadrzędnymi-platformy)
11. [Słowniczek pojęć nawigacyjnych](#11-słowniczek-pojęć-nawigacyjnych)
- [Załącznik A. Macierz dostępności modułów w bocznej nawigacji](#załącznik-a-macierz-dostępności-modułów-w-bocznej-nawigacji)
- [Załącznik B. Panel orkiestracji — zestawienie sekcji](#załącznik-b-panel-orkiestracji--zestawienie-sekcji)
- [Załącznik C. Scenariusz nawigacyjny — pełny przebieg](#załącznik-c-scenariusz-nawigacyjny--pełny-przebieg)

---

## 1. Cel i miejsce dokumentu

Niniejszy dokument rozwija ustalenia Koncepcji platformy do poziomu szczegółowości warstwy nawigacyjnej interfejsu.

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
| Panel orkiestracji środowiska MultitaskingAI | rozdz. 11.8 | Rozdział 6 |

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
| Wejście do środowiska | Prawa kolumna, strefa 1 | Kliknięcie karty środowiska | Otwarcie przestrzeni roboczej jednego z czterech środowisk | 1 | Widoczne bez interakcji |
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
 ═══════════════════════════════════════════════════════════════════════════
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
 ═══════════════════════════════════════════════════════════════════════════
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
| Zasady systemu wizualnego marki | Oszczędne użycie złota jako akcentu; krój nagłówkowy (Cormorant Garamond) zarezerwowany dla tytułów |
| Zgodność formy z rodzajem akcji | Wejście, utworzenie, ustawienie |
| Czytelny rytm strony | Dwie siatki po cztery elementy o różnej skali, dopełnione listwą |

Trzy formy — karty, kafle, listwa — są uporządkowane malejąco pod względem skali i wizualnej masy, co samo w sobie komunikuje hierarchię ważności: użytkownik odczytuje, która strefa jest głównym celem wizyty, zanim przeczyta którykolwiek napis.

### 3.2. Strefa 1 — karty środowisk

Strefa 1 jest punktem wejścia do czterech środowisk platformy: TalkIn, WorkSpace, CodeStudio i MultitaskingAI (Koncepcja platformy, rozdz. 7.1, rozdz. 9). Przyjmuje formę czterech dużych kart wejścia, po jednej na środowisko, stanowiących wizualny środek ciężkości strony głównej.

**Anatomia pojedynczej karty środowiska:**

| Element karty | Zawartość | Uwagi | Warstwa | Sposób wywołania |
|---|---|---|---|---|
| Godło środowiska | Symbol graficzny właściwy danemu środowisku | Zgodny z systemem wizualnym marki | 1 | Widoczne bez interakcji |
| Tytuł | Nazwa środowiska (TalkIn, WorkSpace, CodeStudio, MultitaskingAI) | Złożony krojem nagłówkowym (serif, Cormorant Garamond) | 1 | Widoczny bez interakcji |
| Opis trybu pracy | Jednozdaniowy opis, spójny z definicją środowiska (Koncepcja platformy, rozdz. 9.1–9.4) | Krój tekstowy | 1 | Widoczny bez interakcji |
| Menu karty `⋮` | Otwarcie środowiska w nowej karcie sesji, przypięcie środowiska, wyczyszczenie kart sesji środowiska | Zestaw akcji karty | 3 | Menu kebab `⋮` na karcie |
| Profil startowy środowiska | Wybór modelu, wykonawcy i trybu pracy stosowanego przy wejściu do środowiska | Zwija się po użyciu | 2 | Kliknięcie znacznika kontekstowego w pasku kontekstu |

**Stany karty środowiska:**

| Stan | Wygląd | Znaczenie |
|---|---|---|
| Spoczynkowy | Neutralna powierzchnia karty, bez akcentu złotego | Zgodnie z zasadą oszczędnego użycia złota |
| Najechania (hover) | Akcent złoty na obrysie lub tle karty | Sygnalizuje gotowość do wejścia |
| Aktywny (kliknięcie lub fokus) | Akcent złoty, pierścień fokusu dla nawigacji klawiaturą | Zgodnie z przewodnikiem marki |

**Makieta pojedynczej karty środowiska:**

```
  ┌──────────────────────────────┐
  │                              │
  │            ◈  godło          │
  │                              │
  │          T a l k I n         │  ← tytuł, krój nagłówkowy (Cormorant Garamond)
  │                              │
  │    Wiedza, komunikacja i     │  ← jednozdaniowy opis trybu pracy
  │    praca z treścią           │
  │                              │
  └──────────────────────────────┘
    Spoczynek:  powierzchnia neutralna, bez złota
    Hover / aktywna:  akcent złoty na obrysie + pierścień fokusu
```

Akcent złoty pozostaje zarezerwowany wyłącznie dla stanu aktywnego i najechania, zgodnie z zasadą oszczędnego użycia złota przyjętą w systemie wizualnym marki (przewodnik marki, tokens.json). Afordancja karty odpowiada charakterowi przypisanej jej akcji: „wejdź do przestrzeni roboczej” — kliknięcie karty jest równoznaczne z pierwszym krokiem przepływu nawigacji opisanego w rozdziale 4.1.

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
  │  [ Okno konfiguracji ]    [ Mobile ]    [ Always On Display ] │
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

Makiety w niniejszym dokumencie przedstawiają interfejs w stanie spoczynku: widoczne są wyłącznie elementy warstwy 1 oraz zwinięte wyzwalacze warstw 2–3 — znaczniki kontekstowe, `▼`, `⋮`, `☰`. Elementy warstw 2–4 opisane są w tabelach elementów okna, w kolumnach „Warstwa" i „Sposób wywołania".

---

## 4. Przepływ nawigacji: strona główna → środowisko → moduł → okno

Nawigacja w platformie przebiega w czterech krokach, ustalonych w Koncepcji platformy (rozdz. 8). Niniejszy rozdział rozwija każdy z nich.

**Zestawienie czterech kroków nawigacji:**

| Krok | Z poziomu | Akcja | Prowadzi do | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|
| Krok 1 | Strona główna (strefa 1) | Kliknięcie karty środowiska | Okno środowiska — przestrzeń robocza | 1 | Karta środowiska widoczna bez interakcji; polecenie języka naturalnego w Chat Window |
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
| Chat Window | Okno wspólne obecne na stronie głównej i w każdym z piętnastu modułów, w lewej kolumnie obszaru roboczego; centralny punkt pracy i podstawowy mechanizm sterowania procesami | 1 | Widoczne bez interakcji |
| Execution Loop Window | Okno pętli wykonawczej Koordynator ↔ Wykonawca, otwierane w kolumnie sąsiadującej z Chat Window; dostępne z poziomu strony głównej i z każdego modułu | 2 | Przełącznik `Pętla ▼`; skrót klawiszowy; polecenie języka naturalnego |
| Okna pomocnicze modułu | Panele i monitory otwierane jako rozszerzenia boczne po prawej stronie obszaru roboczego | 3 | Menu `☰` obszaru roboczego; panel popover |

### 4.5. Krok 4 — karty sesji

Aktywne sesje trzymane są w kartach poziomych w oknie środowiska. Karta aktywnej sesji reprezentuje jedną, samodzielną przestrzeń roboczą — z własnym układem okien, historią i kontekstem. Przełączanie kart oznacza przełączanie między sesjami, na wzór kart w przeglądarce internetowej. Pełną mechanikę kart sesji opisuje rozdział 7.

### 4.6. Schemat przepływu nawigacji

```
┌─────────────────────────────────────────────────────────────┐
│ STRONA GŁÓWNA — Centrum dowodzenia                          │
│   Chat Window (lewa kolumna) │ Strefa 1 · Strefa 2 · Strefa 3│
└─────────────────────────────────────────────────────────────┘
        │  KROK 1 — kliknięcie karty środowiska (strefa 1)
        ▼
┌─────────────────────────────────────────────────────────────┐
│ ŚRODOWISKO   TalkIn · WorkSpace · CodeStudio · MultitaskingAI│
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
│  Chat Window │ Execution Loop Window │ okna modułu │ panele  │
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

Pełną macierz dostępności piętnastu modułów w trzech środowiskach modułowych przedstawia Załącznik A (powtórzenie macierzy z rozdziału 10 Koncepcji platformy).

### 5.2. Stan wybrany i wskaźniki

Pozycja modułu aktualnie otwartego w bieżącej karcie sesji jest oznaczona w bocznej nawigacji stanem wybranym (podświetlenie pozycji), tak aby użytkownik zawsze zachowywał orientację co do tego, w jakim module pracuje — zgodnie z zasadą, że zmiana przestrzeni jest zawsze świadomym, widocznym przejściem, a nie ukrytą zmianą stanu (Koncepcja platformy, rozdz. 1).

| Aspekt | Zasada | Warstwa | Sposób wywołania |
|---|---|---|---|
| Oznaczenie modułu otwartego | Stan wybrany — podświetlenie pozycji w bocznej nawigacji | 1 | Widoczne bez interakcji |
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
| Zespoły | Zapisane konfiguracje zespołu — presety ról, powiązań i kolejek; zapis, wczytanie, duplikowanie | rozdz. 11.1 | Konfiguracje ról (rozdz. 11.3) | 1 | Pozycja panelu widoczna bez interakcji |
| Role | Cztery okna robocze: Executor 1, Executor 2, Coordinator, Executor 3 / Validator; przypisanie agentów z modułu Agents do ról; Subagent Network | rozdz. 11.3, 9.15 | Agent Builder, Permissions Center (moduł Agents) | 1 | Pozycja panelu widoczna bez interakcji |
| Kolejki | Definicje kolejek globalnych, lokalnych, modeli, agentów i projektów oraz ich akcje | rozdz. 11.4 | Silnik kolejek, Queue Manager (moduł Automations) | 1 | Pozycja panelu widoczna bez interakcji |
| Orkiestracja | Zależności między modelami, agentami, zadaniami, kolejkami, automatyzacjami i projektami | rozdz. 11.5 | Orchestrator (moduł Automations) | 1 | Pozycja panelu widoczna bez interakcji |
| Harmonogram i automatyki | Harmonogram pracy ciągłej oraz wpięte automatyki z modułu Automations | rozdz. 11.6, 9.3 | Scheduler, Execution Monitor (moduł Automations) | 1 | Pozycja panelu widoczna bez interakcji |
| Monitor procesu | Podgląd przebiegu pętli, statusy przebiegów, hierarchia decyzji; miejsce nadzoru Always On Display | rozdz. 11.2, 7.4, 10.1 | Always On Display, Mobile | 1 | Pozycja panelu widoczna bez interakcji |
| Przypisanie modelu i wykonawcy do roli | Wybór modelu, wykonawcy i poziomu wysiłku dla wybranej roli | rozdz. 11.3 | Agent Builder | 2 | Kliknięcie znacznika kontekstowego roli; selektor zwija się po wyborze |
| Zestaw akcji sekcji | Duplikowanie, eksport, wyczyszczenie i przywrócenie ustawień domyślnych sekcji | rozdz. 11.8 | — | 3 | Menu kebab `⋮` przy sekcji; rozwinięcie `Operacje ▼` |
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
  │  ●  Studio          ✕ │     Research         ✕ │   +   │
  └────────────────────────┴────────────────────────┴───────┘
     ▲ karta aktywna            ▲ karta w tle            ▲ nowa karta
     ●  wskaźnik pracy w tle    ✕  kontrolka zamknięcia
```

**Schemat pojedynczej karty sesji:**

```
  ┌──────────────────────────────┐
  │  ●   Studio               ✕ │
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

Karta sesji jest jednym z siedmiu poziomów zasięgu reguł izolacji (Koncepcja platformy, rozdz. 6.5) — użytkownik może z poziomu okna konfiguracji (rozdz. 9.1) ustanowić jednorazowe współdzielenie historii, pamięci lub kontekstu między konkretnymi kartami, zgodnie z przykładem zastosowania „jednorazowe współdzielenie historii między dwiema otwartymi kartami” z tabeli poziomów zasięgu Koncepcji platformy (rozdz. 6.5).

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

W środowisku MultitaskingAI karta sesji reprezentuje przestrzeń roboczą pojedynczego procesu orkiestracji — obejmującą stan ról, kolejek i orkiestracji widoczny przez panel orkiestracji (rozdz. 6), a nie pojedynczy moduł. Poziom zasięgu „rola” w hierarchii izolacji (Koncepcja platformy, rozdz. 6.5) jest poziomem najbardziej szczegółowym, węższym niż karta sesji.

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
| Narzędzia i workflow dedykowane nowemu modułowi | Pozostałych otwartych kart sesji — każda zachowuje własny, niezależny stan |
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
| Zmiana modułu w bocznej nawigacji (w tej samej karcie) | Zestaw okien operacyjnych, kontekst czatu, narzędzia i workflow | Sama karta sesji, boczna nawigacja, pozostałe karty, funkcje globalne |
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
| Selektor zasięgu | Siedem poziomów: globalny, środowisko, moduł, para modułów, projekt, karta sesji, rola |
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
| Zespoły | Zapisane konfiguracje zespołu — presety ról, powiązań i kolejek | — | 1 | Pozycja panelu widoczna bez interakcji |
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

*Koniec dokumentu. Danaco Console — Strona główna i nawigacja, wersja 2.0.*

---
*Danaco Console — AI Workspace OS · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — [LICENSE](LICENSE). Kontakt: support@danaco-group.pl*
