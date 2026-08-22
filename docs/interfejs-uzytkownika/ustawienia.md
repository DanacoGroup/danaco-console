# Danaco Console — Okno Ustawień (poziom aplikacji)

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
| **Tytuł** | Okno Ustawień (poziom aplikacji) |
| **Klasa dokumentu** | Specyfikacja docelowa |
| **Odbiorcy** | projektant · deweloper |
| **Przeznaczenie** | Ustala projekt okna Ustawień — okna poziomu aplikacji, odrębnego od okna Konfiguracji, gromadzącego wszystko, co dotyczy Operatora i jego relacji z platformą jako całością |
| **Zakres** | konto Operatora, uwierzytelnianie, wygląd (motyw i język), urządzenia (lista, parowanie Mobile, synchronizacja), powiadomienia, zarządzanie kontami modeli, ustawienia Chat Window i Execution Loop Window, warstwy widoczności według roli użytkownika, Always On Display, makieta, katalog elementów, stany, komunikacja klient–serwer |
| **Poza zakresem** | pozostałe dwanaście zakresów Modelu konfiguracji (Procesy, Akcje, Zachowanie modeli, Tożsamość modeli, Prompty systemowe, Rozszerzenia, Integracje, Historia, Pamięć, Izolacja, Karty sesji, Komponenty własne) — [Okno Konfiguracji](konfiguracja.md); pełny protokół parowania i pełny opis funkcji Always On Display — [Mobile](../funkcje-globalne/mobile.md) oraz [Always On Display](../funkcje-globalne/always-on-display.md) |
| **Dokument nadrzędny** | [Elementy okien](elementy-okien.md) |
| **Dokumenty powiązane** | [Elementy okien](elementy-okien.md) · [Katalog komponentów](katalog-komponentow.md) · [Przepływ okien](przeplyw-okien.md) · [Okno Konfiguracji](konfiguracja.md) · [Koncepcja platformy](../architektura/koncepcja-platformy.md) · [Model konfiguracji](../architektura/model-konfiguracji.md) · [Bezpieczeństwo i uwierzytelnianie](../architektura/bezpieczenstwo-i-uwierzytelnianie.md) · [Strona główna i nawigacja](strona-glowna-i-nawigacja.md) · [System wizualny](system-wizualny.md) · [Integracja modeli](../architektura/integracja-modeli.md) · [Mobile](../funkcje-globalne/mobile.md) · [Always On Display](../funkcje-globalne/always-on-display.md) |
| **Prototypy odniesienia** | `design/05-okna/platformowe/ustawienia.html` |
| **Źródła normatywne** | `budowa/shared/contract.json` (obszary `auth`, `account`, `device`, `config`, `panel`, `aod`) · `design/zasoby/zetony/zetony.css` · `design/zasoby/css/komponenty.css`, `rama.css`, `prototyp.css` |
| **Zasada nadrzędna** | Żadne ustawienie w tym oknie nie jest wymuszone — każda metoda logowania, każda klasa powiadomień i każde konto modelu jest ustawieniem konfiguracyjnym Operatora, nie wymogiem |

Dokument jest projektem okna Ustawień — okna poziomu aplikacji, odrębnego od okna Konfiguracji. Ustawienia gromadzą wszystko, co dotyczy Operatora i jego relacji z platformą jako całością: konto, uwierzytelnianie, wygląd, urządzenia, powiadomienia oraz zarządzanie kontami modeli wykorzystywanymi przez platformę. Konfiguracja pozostaje odrębnym oknem odpowiedzialnym za pozostałe zakresy ustawień platformy — zachowanie procesów, akcji, modeli, promptów, rozszerzeń, integracji, historii, pamięci, izolacji, kart sesji i komponentów własnych. Rozdział 1 określa dokładny podział między oboma oknami. Dokument zawiera pełną anatomię okna, makiety tekstowe w układzie pionowym (podział lewa–prawa), katalog elementów interfejsu ze wszystkimi wymaganymi atrybutami wraz z warstwą widoczności i sposobem wywołania oraz zestawienie stanów — zgodnie ze standardem przyjętym dla całego opracowania dokumentacji projektowej interfejsu. Okno Ustawień pracuje w tym samym kontekście pracy co pozostałe okna platformy: w lewej kolumnie obszaru roboczego stale obecne jest Chat Window — główne okno komunikacji Użytkownik ↔ Wykonawca, a w kolumnie z nim sąsiadującej Execution Loop Window — okno pętli wykonawczej Koordynator ↔ Wykonawca (rozdz. 9).

---

## Spis treści

1. [Wprowadzenie i mapa źródeł](#wprowadzenie-i-mapa-źródeł)
2. [Miejsce Ustawień w architekturze platformy](#1-miejsce-ustawień-w-architekturze-platformy)
   - [1.1 Ustawienia a Konfiguracja — rozgraniczenie](#11-ustawienia-a-konfiguracja--rozgraniczenie)
   - [1.2 Klasyfikacja okna](#12-klasyfikacja-okna)
   - [1.3 Zasada zerowych blokad zastosowana do Ustawień](#13-zasada-zerowych-blokad-zastosowana-do-ustawień)
3. [Struktura okna, nawigacja wewnętrzna i warstwy widoczności](#2-struktura-okna-nawigacja-wewnętrzna-i-warstwy-widoczności)
   - [2.1 Układ kolumnowy okna](#21-układ-kolumnowy-okna)
   - [2.2 Dziewięć sekcji selektora](#22-dziewięć-sekcji-selektora)
   - [2.3 Mechanizm objaśnień kontekstowych](#23-mechanizm-objaśnień-kontekstowych)
   - [2.4 Nawigacja z poziomu środowiska](#24-nawigacja-z-poziomu-środowiska)
   - [2.5 Warstwy widoczności w oknie](#25-warstwy-widoczności-w-oknie)
4. [Konto Operatora](#3-konto-operatora)
   - [3.1 Model jednego konta](#31-model-jednego-konta)
   - [3.2 Zawartość sekcji](#32-zawartość-sekcji)
   - [3.3 Makieta sekcji Konto](#33-makieta-sekcji-konto)
   - [3.4 Zachowanie: zmiana hasła](#34-zachowanie-zmiana-hasła)
5. [Uwierzytelnianie](#4-uwierzytelnianie)
   - [4.1 Zasada: wszystkie metody logowania jako ustawienie konfiguracyjne, hasło jako domyślna](#41-zasada-wszystkie-metody-logowania-jako-ustawienie-konfiguracyjne-hasło-jako-domyślna)
   - [4.2 Makieta sekcji Uwierzytelnianie](#42-makieta-sekcji-uwierzytelnianie)
   - [4.3 Wymóg logowania — ustawienie Konfiguracji, dostępne Operatorowi](#43-wymóg-logowania--ustawienie-konfiguracji-dostępne-operatorowi)
   - [4.4 Urządzenia połączone — odesłanie](#44-urządzenia-połączone--odesłanie)
6. [Wygląd: motyw i język](#5-wygląd-motyw-i-język)
   - [5.1 Motyw wizualny](#51-motyw-wizualny)
   - [5.2 Język interfejsu](#52-język-interfejsu)
   - [5.3 Makieta sekcji Wygląd i język](#53-makieta-sekcji-wygląd-i-język)
   - [5.4 Zachowanie](#54-zachowanie)
7. [Urządzenia: lista, parowanie Mobile, synchronizacja](#6-urządzenia-lista-parowanie-mobile-synchronizacja)
   - [6.1 Wykaz urządzeń połączonych](#61-wykaz-urządzeń-połączonych)
   - [6.2 Parowanie urządzenia mobilnego](#62-parowanie-urządzenia-mobilnego)
   - [6.2 a Zarządzanie sparowanymi urządzeniami](#62a-zarządzanie-sparowanymi-urządzeniami)
   - [6.3 Synchronizacja](#63-synchronizacja)
   - [6.4 Makieta sekcji Urządzenia](#64-makieta-sekcji-urządzenia)
8. [Powiadomienia](#7-powiadomienia)
   - [7.1 Zakres ustawienia źródłowego](#71-zakres-ustawienia-źródłowego)
   - [7.2 Klasy zdarzeń](#72-klasy-zdarzeń)
   - [7.3 Kanał dostarczenia](#73-kanał-dostarczenia)
   - [7.4 Makieta sekcji Powiadomienia](#74-makieta-sekcji-powiadomienia)
   - [7.5 Zachowanie](#75-zachowanie)
9. [Zarządzanie kontami modeli](#8-zarządzanie-kontami-modeli)
   - [8.1 Miejsce w mechanizmie kanału modelu](#81-miejsce-w-mechanizmie-kanału-modelu)
   - [8.2 Zawartość sekcji](#82-zawartość-sekcji)
   - [8.3 Dodanie konta — okno nakładkowe](#83-dodanie-konta--okno-nakładkowe)
   - [8.4 Przełączanie aktywnego konta — zasada jawności](#84-przełączanie-aktywnego-konta--zasada-jawności)
   - [8.5 Makieta sekcji Konta modeli](#85-makieta-sekcji-konta-modeli)
   - [8.6 Stan pustej listy](#86-stan-pustej-listy)
10. [Okna komunikacji operacyjnej i ich ustawienia](#9-okna-komunikacji-operacyjnej-i-ich-ustawienia)
   - [9.1 Chat Window — kanał Użytkownik ↔ Wykonawca](#91-chat-window--kanał-użytkownik--wykonawca)
   - [9.2 Execution Loop Window — kanał Koordynator ↔ Wykonawca](#92-execution-loop-window--kanał-koordynator--wykonawca)
   - [9.3 Sekcja „Okna komunikacji” — ustawienia obu kanałów](#93-sekcja-okna-komunikacji--ustawienia-obu-kanałów)
   - [9.4 Makieta sekcji Okna komunikacji](#94-makieta-sekcji-okna-komunikacji)
11. [Warstwy widoczności — ustawienia według roli użytkownika](#10-warstwy-widoczności--ustawienia-według-roli-użytkownika)
   - [10.1 Grupa ustawień „Warstwy widoczności według roli”](#101-grupa-ustawień-warstwy-widoczności-według-roli)
   - [10.2 Makieta sekcji Warstwy widoczności](#102-makieta-sekcji-warstwy-widoczności)
12. [Always On Display](#11-always-on-display)
   - [11.1 Zakres sekcji](#111-zakres-sekcji)
   - [11.2 Zawartość sekcji](#112-zawartość-sekcji)
   - [11.3 Makieta sekcji Always On Display](#113-makieta-sekcji-always-on-display)
   - [11.4 Zachowanie](#114-zachowanie)
13. [Makieta całościowa okna Ustawień](#12-makieta-całościowa-okna-ustawień)
14. [Katalog elementów interfejsu](#13-katalog-elementów-interfejsu)
15. [Stany — zestawienie zbiorcze](#14-stany--zestawienie-zbiorcze)
   - [14.1 Stany okna jako całości](#141-stany-okna-jako-całości)
   - [14.2 Stany sekcji z listami i tabelami](#142-stany-sekcji-z-listami-i-tabelami)
   - [14.3 Stany formularzy i okien nakładkowych](#143-stany-formularzy-i-okien-nakładkowych)
   - [14.4 Stany przełączników (wzorzec wspólny)](#144-stany-przełączników-wzorzec-wspólny)
16. [Komunikacja — polecenia i zdarzenia](#15-komunikacja--polecenia-i-zdarzenia)
   - [15.1 Polecenia i zdarzenia przejęte wprost](#151-polecenia-i-zdarzenia-przejęte-wprost)
   - [15.2 Polecenia i zdarzenia zakresu Konta Operatora, kont modeli, okien komunikacji i warstw widoczności](#152-polecenia-i-zdarzenia-zakresu-konta-operatora-kont-modeli-okien-komunikacji-i-warstw-widoczności)
   - [15.3 Schemat przepływu — zapis pojedynczego ustawienia](#153-schemat-przepływu--zapis-pojedynczego-ustawienia)
17. [Zgodność z zasadami nadrzędnymi platformy](#16-zgodność-z-zasadami-nadrzędnymi-platformy)
18. [Słowniczek pojęć](#17-słowniczek-pojęć)
19. [Załącznik A. Scenariusze użycia](#załącznik-a-scenariusze-użycia)
20. [Załącznik B. Szablony konfiguracji](#załącznik-b-szablony-konfiguracji)
21. [Załącznik C. Mapa zgodności z dokumentami źródłowymi](#załącznik-c-mapa-zgodności-z-dokumentami-źródłowymi)

---

## Wprowadzenie i mapa źródeł

Koncepcja platformy (rozdz. 7.3) umieszcza w Strefie 3 strony głównej „okno konfiguracji i ustawień platformy”. Model konfiguracji (rozdz. 5.1) grupuje pod nazwą zakresu „Aplikacja” pięć ustawień: motyw wizualny, język interfejsu, metody uwierzytelniania, urządzenia połączone i powiadomienia. Zakres ten jest samodzielnym, pełnoprawnym oknem poziomu aplikacji — Ustawieniami — dostępnym wprost z listwy ustawień strony głównej, obok (nie: wewnątrz) okna Konfiguracji. Zakres „Aplikacja” obejmuje ponadto jawne zarządzanie profilem Konta Operatora (login, e-mail, hasło), zarządzanie kontami modeli wykorzystywanymi przez platformę w kanale wiersza poleceń, ustawienia obu okien komunikacji operacyjnej (rozdz. 9) oraz ustawienia warstw widoczności według roli użytkownika (rozdz. 10). Poniższa mapa wiąże każdy element Ustawień z jego źródłem i miejscem w dokumencie.

| Element Ustawień | Źródło ustalenia | Rozdział źródła | Miejsce w dokumencie |
|---|---|---|---|
| Wydzielenie „Aplikacja” jako samodzielnego okna | Koncepcja platformy; Model konfiguracji | rozdz. 5.3–5.4; rozdz. 3.1, 5.1 | Rozdz. 1 |
| Konto Operatora — model jednego konta i wielu urządzeń | Bezpieczeństwo i uwierzytelnianie | rozdz. 3, 11.1 | Przejęte wprost |
| Konto Operatora — edycja loginu, e-maila, hasła w trakcie sesji | — | — | Rozdz. 3, rozdz. 15 |
| Metody uwierzytelniania (hasło, PIN, Windows Hello, e-mail) | Bezpieczeństwo i uwierzytelnianie | rozdz. 6, 13 | Przejęte wprost |
| Wymóg logowania i przycisk Pomiń w oknie logowania | Bezpieczeństwo i uwierzytelnianie | rozdz. 8 | Wymóg logowania jest ustawieniem okna Konfiguracji, dostępnym Operatorowi w każdej chwili (rozdz. 4.3) — nie ustawieniem Ustawień |
| Urządzenia połączone (wykaz, unieważnienie tokenu) | Bezpieczeństwo i uwierzytelnianie | rozdz. 3, 11.2, 13 | Przejęte wprost |
| Parowanie urządzenia mobilnego | [Mobile](../funkcje-globalne/mobile.md) | rozdz. 9 | Przejęte przez odesłanie — pełny protokół pozostaje w dokumencie funkcji globalnej; tu opisany interfejs parowania i zarządzania sparowanymi urządzeniami |
| Synchronizacja na żywo | Architektura | rozdz. 13 | Przejęte wprost, zaprezentowane jako wskaźnik stanu |
| Motyw wizualny, mechanizm przełączania | System wizualny; Model konfiguracji | rozdz. 9; rozdz. 5.1 | Przejęte wprost |
| Język interfejsu | Model konfiguracji | rozdz. 5.1 | Przejęte wprost |
| Powiadomienia — zakres i kanał | Model konfiguracji; System wizualny; Katalog komponentów | rozdz. 5.1; rozdz. 8.10; rozdz. 11.6 | Przejęte i rozwinięte do poziomu elementów interfejsu |
| Kanał modelu — encja i dane dostępowe poza bazą | Integracja modeli; Bezpieczeństwo i uwierzytelnianie | rozdz. 5, 8, 10; rozdz. 10 | Podstawa mechanizmu; zastosowana do nowego zakresu |
| Zarządzanie kontami modeli (rotacja kont, katalogi profili) | — | — | Rozdz. 8, na mechanizmie kanału modelu |
| Chat Window i Execution Loop Window — ustawienia obu kanałów komunikacji operacyjnej | Koncepcja platformy | rozdz. 1, 11.8 | Rozdz. 9 |
| Warstwy widoczności interfejsu według roli użytkownika | Koncepcja platformy; System wizualny | rozdz. 12; rozdz. 10 | Rozdz. 2.5, rozdz. 10 |
| Always On Display — obecność, tryb, reguły wyzwalania, wyciszanie, tor głosowy | funkcje-globalne/always-on-display.md; Model konfiguracji | rozdz. 3, 5, 10; rozdz. 5.17 | Przejęte przez odesłanie — pełny opis funkcji pozostaje w dokumencie funkcji globalnej; tu opisana sekcja ustawień (rozdz. 11) |

---

## 1. Miejsce Ustawień w architekturze platformy

### 1.1 Ustawienia a Konfiguracja — rozgraniczenie

Model konfiguracji (rozdz. 5) wylicza siedemnaście zakresów ustawień platformy. Pięć z nich — „Aplikacja” (5.1), „Chat Window” (5.14), „Execution Loop Window” (5.15), „Warstwy widoczności funkcji” (5.16) i „Always On Display” (5.17) — ma status osobnego okna, dostępnego bezpośrednio z listwy ustawień strony głównej. Pozostałe dwanaście zakresów pozostaje w oknie Konfiguracji, opisanym odrębnym opracowaniem. Rozgraniczenie nie jest podziałem technicznym, lecz podziałem według rodzaju decyzji: Ustawienia gromadzą to, co dotyczy Operatora i jego tożsamości w platformie; Konfiguracja gromadzi to, co dotyczy zachowania samej platformy podczas pracy.

| Cecha | Ustawienia (niniejszy dokument) | Konfiguracja (dokument odrębny) |
|---|---|---|
| Pytanie, na które odpowiada | „Kim jestem w platformie i jak z niej korzystam na moich urządzeniach?” | „Jak platforma ma się zachowywać podczas pracy?” |
| Zakres z Modelu konfiguracji | 5.1 Aplikacja | 5.2–5.13 — Procesy, Akcje, Zachowanie modeli, Tożsamość modeli, Prompty systemowe, Rozszerzenia, Integracje, Historia, Pamięć, Izolacja, Karty sesji, Komponenty własne |
| Warstwowość | Głównie warstwa globalna — dotyczy Operatora i jego urządzeń jako takich | Cztery warstwy: globalna, środowisko, projekt, sesja (Model konfiguracji, rozdz. 4) |
| Częstość zmian w typowym użyciu | Rzadka — ustawiane raz, zmieniane okazjonalnie | Częsta — dostosowywana do bieżącej pracy, sesji i roli |
| Punkt wejścia | Listwa ustawień strony głównej, pozycja „Ustawienia” | Listwa ustawień strony głównej, pozycja „Konfiguracja”; uproszczone menu kontekstowe okna operacyjnego |
| Klasyfikacja okna (rozdz. 1.2) | Okno globalne | Okno globalne |

Schemat 1 osadza obie ścieżki w strefie 3 strony głównej, zgodnie z Koncepcją platformy (rozdz. 7.3) i Stroną główną i nawigacją (rozdz. 3.4).

```
Schemat 1 — Strefa 3 strony głównej po wydzieleniu Ustawień

STRONA GŁÓWNA — Centrum dowodzenia
├─ Strefa 1 · wybór środowiska        (waga główna)
├─ Strefa 2 · komponenty własne       (waga pośrednia)
└─ Strefa 3 · listwa ustawień         (waga najniższa)
      │
      ├──[ Ustawienia ]──────► NINIEJSZY DOKUMENT
      │                         Konto · Uwierzytelnianie · Wygląd i język ·
      │                         Urządzenia · Powiadomienia · Konta modeli
      │
      ├──[ Konfiguracja ]────► dokument odrębny
      │                         Procesy · Akcje · Zachowanie modeli ·
      │                         Tożsamość modeli · Prompty systemowe ·
      │                         Rozszerzenia · Integracje · Historia ·
      │                         Pamięć · Izolacja · Karty sesji ·
      │                         Komponenty własne
      │
      ├──[ Mobile ]──────────► funkcja globalna (Koncepcja platformy, rozdz. 12.1;
      │                          funkcje-globalne/mobile.md)
      │                         parowanie opisane w rozdz. 6.2 niniejszego dokumentu
      │
      └──[ Always On Display ]► funkcja globalna (funkcje-globalne/always-on-display.md)
                                 nadzór nad Mobile i Monitorem procesu — poza zakresem
```

### 1.2 Klasyfikacja okna

Okno Ustawień, tak jak okno konfiguracji punktów izolacji (Specyfikacja okien operacyjnych, rozdz. 8.1), nie jest oknem modułowym przypisanym do jednego z piętnastu modułów — jest oknem o zasięgu globalnym.

| Cecha | Wartość |
|---|---|
| Kategoria okna | Okno globalne (nieprzypisane do modułu, nieprzypisane do środowiska) |
| Miejsce otwarcia | Listwa ustawień, strefa 3 strony głównej; dostępne również z poziomu każdego środowiska analogicznie do Konfiguracji (Model konfiguracji, rozdz. 3.1) |
| Zasięg oddziaływania | Cała platforma, wszystkie urządzenia Operatora jednocześnie (rozdz. 3, rozdz. 6) |
| Przeładowanie przy zmianie modułu | Nie dotyczy — Ustawienia nie należą do zestawu okien żadnego modułu (Specyfikacja okien operacyjnych, rozdz. 3) |
| Trwałość stanu | Zmiany zapisywane natychmiast po stronie serwera; widoczne na wszystkich urządzeniach przez zdarzenia rozgłoszeniowe (rozdz. 15.2) |
| Podstawa | Koncepcja platformy, rozdz. 7.3; Model konfiguracji, rozdz. 5.1 |

### 1.3 Zasada zerowych blokad zastosowana do Ustawień

Zgodnie z zasadą nadrzędną pełnej konfigurowalności (Koncepcja platformy, rozdz. 14, zasada 2) żadne ustawienie w tym oknie nie jest wymuszone. Tabela poniżej stosuje tę zasadę wprost do każdej sekcji Ustawień.

| Sekcja | Ustawienie konfiguracyjne, nie wymóg | Warstwa | Sposób wywołania |
|---|---|---|---|
| Uwierzytelnianie (rozdz. 4) | Wszystkie metody logowania — hasło, PIN, Windows Hello, e-mail — włącza i wyłącza Operator niezależnie; hasło pozostaje metodą domyślną, aktywną od rejestracji, ale nie jest przymusem: wyłączenie ostatniej aktywnej metody logowania sekcja sygnalizuje ostrzeżeniem, nie blokadą | 2 | Pozycja „Uwierzytelnianie” w selektorze sekcji `☰` |
| Wygląd (rozdz. 5) | Motyw i język mają wartość domyślną (motyw systemowy, język polski); brak ustawienia nie ogranicza pracy | 2 | Pozycja „Wygląd i język” w selektorze sekcji `☰` |
| Urządzenia (rozdz. 6) | Brak ograniczenia liczby sparowanych urządzeń; synchronizacja jest właściwością platformy, nie ustawieniem do włączenia | 2 | Pozycja „Urządzenia” w selektorze sekcji `☰` |
| Powiadomienia (rozdz. 7) | Każda klasa zdarzeń wyłączalna niezależnie; wyłączenie wszystkich nie ogranicza dostępu do żadnej funkcji platformy | 2 | Pozycja „Powiadomienia” w selektorze sekcji `☰` |
| Konta modeli (rozdz. 8) | Wielość kont jest ustawieniem Operatora, nie wymogiem; jedno skonfigurowane konto jest w pełni wystarczające do pracy | 2 | Pozycja „Konta modeli” w selektorze sekcji `☰` |
| Okna komunikacji (rozdz. 9) | Sposób prezentacji Chat Window i Execution Loop Window jest sterowany ustawieniami; oba okna pozostają obecne w każdym kontekście pracy | 2 | Pozycja „Okna komunikacji” w selektorze sekcji `☰` |
| Warstwy widoczności (rozdz. 10) | Przypisanie warstw do roli użytkownika jest ustawieniem Operatora; wartości domyślne wystarczają do pracy bez interwencji | 4 | Pozycja „Warstwy widoczności” w selektorze sekcji `☰`, tryb administracyjny, polecenie języka naturalnego w Chat Window |

---

## 2. Struktura okna, nawigacja wewnętrzna i warstwy widoczności

### 2.1 Układ kolumnowy okna

Okno Ustawień przyjmuje układ pionowy (podział lewa–prawa) obowiązujący w całej platformie. Kolumny sąsiadują ze sobą poziomo i mają pełną wysokość obszaru roboczego; regulacji podlega wyłącznie ich szerokość.

| Kolumna | Zawartość | Waga wizualna |
|---|---|---|
| Lewa, stała, pełna wysokość | **Chat Window** — główne okno komunikacji Użytkownik ↔ Wykonawca; centralny punkt pracy i podstawowy mechanizm sterowania procesami platformy, w tym każdą pozycją Ustawień (rozdz. 9.1) | Lewa kolumna, stała, pełna wysokość obszaru roboczego |
| Kolumna sąsiadująca (otwierana) | **Execution Loop Window** — okno pętli wykonawczej Koordynator ↔ Wykonawca; drugi kanał komunikacji operacyjnej (rozdz. 9.2) | Kolumna sąsiadująca, otwierana |
| Kolumna selektora sekcji | Selektor sekcji Ustawień `☰` — wybór jednej z dziewięciu sekcji | Kolumna wąska, waga niska |
| Prawa, dominująca | Panel zawartości sekcji — pozycje ustawień właściwe wybranej sekcji | Kolumna dominująca obszaru roboczego |
| Kolejne kolumny boczne | Okna pomocnicze i panele — otwierane jako rozszerzenia boczne, po prawej stronie obszaru roboczego | Kolumna boczna, otwierana jako rozszerzenie boczne |

```
Schemat 2 — Układ kolumnowy okna Ustawień (podział lewa–prawa)

 ═════════════════════════════════════════════════════════════════════════════════
  Pasek kontekstu:  [Danaco Console] [Ubuntu] [Worktree] [Fable 5] [Ultra]
 ═════════════════════════════════════════════════════════════════════════════════
  Chat Window      │ Execution Loop ▼│ Sekcje ☰            │ ZAWARTOŚĆ SEKCJI                      │
  Użytkownik ↔     │ Koordynator ↔   │                     │                                       │
  Wykonawca        │ Wykonawca       │ ◈ Konto             │ [ nazwa ustawienia ]            [?] ⋮ │
                   │                 │   Uwierzytelnianie  │     wartość bieżąca                   │
  [ polecenie… ]   │ [ pętla: 0 zad.]│   Wygląd i język    │                                       │
                   │                 │   Urządzenia        │ [ nazwa ustawienia ]            [?] ⋮ │
                   │                 │   Powiadomienia     │     wartość bieżąca                   │
                   │                 │   Konta modeli      │                                       │
                   │                 │   Okna komunikacji  │ …                                     │
                   │                 │  Warstwy widoczności│                                       │
                   │                 │   Always On Display │                                       │
 ═════════════════════════════════════════════════════════════════════════════════

Stan spoczynku: widoczne wyłącznie elementy warstwy 1 oraz zwinięte
wyzwalacze warstw 2–3 (`▼`, `⋮`, `☰`) i znaczniki paska kontekstu.
```

### 2.2 Dziewięć sekcji selektora

| Sekcja | Zawartość skrótowa | Ikona (System wizualny, Załącznik C) | Warstwa | Sposób wywołania | Rozdział niniejszego dokumentu |
|---|---|---|---|---|---|
| Konto | Login, e-mail uwierzytelniający, hasło, awatar, data utworzenia konta | `uzytkownik` | 1 | Sekcja otwierana domyślnie przy wejściu do okna | 3 |
| Uwierzytelnianie | Metody logowania (hasło i metody dodatkowe), stan wymogu logowania | `klodka` | 2 | Pozycja selektora sekcji `☰` | 4 |
| Wygląd i język | Motyw jasny/ciemny, język interfejsu | `slonce` / `ksiezyc` | 2 | Pozycja selektora sekcji `☰` | 5 |
| Urządzenia | Wykaz urządzeń połączonych, parowanie Mobile, stan synchronizacji | `telefon` (urządzenie mobilne) | 2 | Pozycja selektora sekcji `☰` | 6 |
| Powiadomienia | Przełącznik główny, klasy zdarzeń, kanał dostarczenia | `dzwonek` | 2 | Pozycja selektora sekcji `☰` | 7 |
| Konta modeli | Wykaz kont modeli w kanale CLI, w tym kont narzędzia Code CLI, katalogi profili | `kod` | 2 | Pozycja selektora sekcji `☰` | 8 |
| Okna komunikacji | Ustawienia Chat Window i Execution Loop Window — dwóch kanałów komunikacji operacyjnej | `rozmowa` | 2 | Pozycja selektora sekcji `☰`; polecenie języka naturalnego w Chat Window | 9 |
| Warstwy widoczności | Przypisanie warstw widoczności funkcji do roli użytkownika | `warstwy` | 4 | Tryb administracyjny, wyszukiwarka funkcji, skrót klawiszowy, polecenie języka naturalnego w Chat Window | 10 |
| Always On Display | Obecność i tryb funkcji globalnej Always On Display, reguły wyzwalania sugestii, wyciszanie, tor głosowy | `dzwonek` | 2 | Pozycja selektora sekcji `☰`; polecenie języka naturalnego w Chat Window | 11 |

### 2.3 Mechanizm objaśnień kontekstowych

Każda pozycja ustawienia w oknie Ustawień jest opatrzona objaśnieniem kontekstowym `[?]`, zgodnie z zasadą jednolitą dla całego okna Konfiguracji i przyjętą tu wprost (Model konfiguracji, rozdz. 3.3; Architektura, rozdz. 13). Nośnikiem jest komponent `.dn-tooltip` (System wizualny, rozdz. 8.11), wyzwalany wskazaniem kursorem lub fokusem klawiatury.

| Ustawienie | Przykładowa treść objaśnienia `[?]` |
|---|---|
| Metoda PIN | „Krótki kod numeryczny do szybkiego potwierdzenia tożsamości na tym urządzeniu. Nie zastępuje hasła — obie metody działają niezależnie od siebie; hasło pozostaje metodą domyślną, chyba że Operator wyłączy je osobno w tej samej sekcji.” |
| Motyw wizualny | „Wybór wariantu jasnego lub ciemnego interfejsu. Bez wyboru platforma stosuje motyw zgodny z ustawieniem systemu operacyjnego urządzenia.” |
| Konto modelu — katalog profilu | „Lokalny katalog przechowujący dane logowania i pliki konfiguracyjne tego konta. Każde konto ma odrębny katalog, dzięki czemu konta nie współdzielą sesji ani pamięci podręcznej narzędzia.” |

### 2.4 Nawigacja z poziomu środowiska

Analogicznie do okna Konfiguracji (Model konfiguracji, rozdz. 3.1), Ustawienia są dostępne nie tylko ze strony głównej, lecz z poziomu każdego z czterech środowisk — jako pozycja stale obecna, niezależna od aktywnego modułu i karty sesji, ponieważ dotyczą Operatora jako całości, a nie pojedynczej sesji. Otwarcie Ustawień nie zamyka ani nie przesłania Chat Window i Execution Loop Window: oba kanały komunikacji operacyjnej pozostają w swoich kolumnach, a Operator zmienia dowolne ustawienie zarówno przez element interfejsu w panelu zawartości, jak i poleceniem języka naturalnego skierowanym do Wykonawcy w Chat Window (rozdz. 9).

### 2.5 Warstwy widoczności w oknie

Okno Ustawień stosuje regułę stopniowego ujawniania funkcjonalności obowiązującą w całej platformie: jeżeli funkcja nie jest potrzebna do realizacji aktualnego zadania, nie jest widoczna. Każdy element okna należy do dokładnie jednej z czterech warstw widoczności; przynależność elementu i sposób jego wywołania podaje katalog elementów (rozdz. 13) oraz tabele ustawień w rozdziałach 3–11.

| Warstwa | Zawartość w oknie Ustawień | Sposób wywołania |
|---|---|---|
| 1 — zawsze widoczna | Chat Window w lewej kolumnie, pasek kontekstu, nagłówek okna, panel zawartości sekcji otwartej (domyślnie „Konto”), pola i działania tej sekcji, wskaźnik stanu synchronizacji, plakietka stanu wymogu logowania | Widoczne bez interakcji |
| 2 — widoczna na żądanie | Execution Loop Window, selektor sekcji `☰` i pozostałe sekcje Ustawień, przełączniki metod logowania, przełącznik motywu, lista języka, tabela urządzeń, tabela klas zdarzeń powiadomień, tabela kont modeli, przełącznik obecności awatara Always On Display, przełącznik trybu obecności i przełącznik toru głosowego funkcji globalnej, objaśnienia kontekstowe `[?]`, selektory znaczników paska kontekstu | Wyzwalacz `▼`, wyzwalacz `☰`, znacznik kontekstowy; po użyciu element zwija się samoczynnie |
| 3 — rozwinięcia kontekstowe | Menu akcji wiersza `⋮` tabel urządzeń i kont modeli, okno nakładkowe zmiany hasła, okno nakładkowe dodania konta modelu, formularz „Ustaw PIN”, działanie „Testuj połączenie”, menu wyciszania sugestii Always On Display `⋮`, odnośnik do przełącznika wymogu logowania w oknie Konfiguracji | Menu kebab `⋮`, menu hamburger `☰`, menu kontekstowe, panel popover, lista rozwijana |
| 4 — funkcje eksperckie | Sekcja „Warstwy widoczności” (rozdz. 10), tabela reguł wyzwalania sugestii Always On Display (rozdz. 11.2), podgląd i podmiana danych dostępowych kont modeli, narzędzia diagnostyczne kanałów komunikacji, tryb administracyjny okna | Polecenie języka naturalnego w Chat Window, skrót klawiszowy, wyszukiwarka funkcji, tryb administracyjny, konfiguracja roli; użytkownik podstawowy nie widzi tych elementów |

Zawsze widoczne pozostają: Chat Window jako centralny punkt pracy i podstawowy mechanizm sterowania procesami platformy, pasek kontekstu ze znacznikami środowiska, repozytorium, projektu, modelu i wykonawcy oraz zawartość sekcji aktualnie otwartej. Pod selektorem `☰` znajdują się wszystkie pozostałe sekcje Ustawień oraz Execution Loop Window, otwierane wyzwalaczem `▼` w kolumnie sąsiadującej. W menu kontekstowym `⋮` znajdują się akcje wierszy tabel i warianty operacji. Wyłącznie poleceniem języka naturalnego w Chat Window, skrótem klawiszowym, wyszukiwarką funkcji albo w trybie administracyjnym dostępne są: konfiguracja warstw widoczności według roli, podgląd danych dostępowych kont modeli, diagnostyka kanałów komunikacji oraz zbiorcze przywrócenie wartości domyślnych całego okna.

Obowiązuje zasada jednego kliknięcia: każda ukryta funkcja okna jest osiągalna jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego skierowanym do Wykonawcy w Chat Window. Ukrycie zmniejsza chaos wizualny i nie wydłuża drogi dostępu.

Mechanizmy ukrywania funkcjonalności zastosowane w oknie: menu progresywne (`Sekcje ☰`, `Operacje ▼`), panele wysuwane (okna nakładkowe, panele pomocnicze jako rozszerzenia boczne), grupowanie logiczne akcji (menu wiersza `⋮`) oraz znaczniki kontekstowe paska kontekstu, których kliknięcie otwiera odpowiedni selektor.


---

## 3. Konto Operatora

### 3.1 Model jednego konta

Platforma nie jest systemem wielokontowym (Bezpieczeństwo i uwierzytelnianie, rozdz. 3). Sekcja Konto zarządza jedyną encją Konto właściwą Operatorowi oraz wykazem urządzeń, na których się uwierzytelnił.

```
Schemat 3 — Jedno konto, wiele urządzeń (Bezpieczeństwo i uwierzytelnianie, rozdz. 3)

              ┌───────────────────────────────┐
              │  KONTO — jedyne (Operator)    │
              │  zarządzane w sekcji „Konto”  │
              └────────────────┬──────────────┘
                               │  jedno konto, wiele urządzeń
        ┌────────────┬─────────┴─────────┬────────────┐
        ▼            ▼                   ▼            ▼
   Urządzenie    Urządzenie          Urządzenie      …
   (komputer)     (telefon)           (tablet)   bez limitu
   — token A       — token B           — token C
        wykaz i zarządzanie w sekcji „Urządzenia” (rozdz. 6)
```

### 3.2 Zawartość sekcji

| Pole | Opis | Źródło wartości | Edytowalność |
|---|---|---|---|
| Awatar | Symbol graficzny identyfikujący Operatora w oknach wielomodelowych i w panelu ról MultitaskingAI | Inicjał loginu domyślnie; obraz wgrywalny przez Operatora | Edytowalny (wgranie / usunięcie obrazu) |
| Login | Nazwa właściciela konta wykorzystywana przy logowaniu (Bezpieczeństwo i uwierzytelnianie, rozdz. 4.2) | Ustalony przy rejestracji | Edytowalny — potwierdzenie bieżącym hasłem jest ustawieniem konfiguracyjnym włączanym w sekcji Uwierzytelnianie (rozdz. 4.1); domyślnie zmiana zapisuje się bez tego kroku |
| Adres e-mail uwierzytelniający | Pełni funkcję weryfikacji, logowania i odzyskiwania konta (Bezpieczeństwo i uwierzytelnianie, rozdz. 4.3, 6.3, 7) | Ustalony przy rejestracji | Edytowalny — zmiana uruchamia ponowną weryfikację adresu (ten sam mechanizm co rozdz. 4.3) |
| Stan weryfikacji e-mail | Potwierdzony / niepotwierdzony | Wynik ostatniej weryfikacji | Tylko do odczytu; działanie: „Wyślij ponownie potwierdzenie” |
| Hasło | Metoda bazowa uwierzytelniania, obecna zawsze (rozdz. 4) | Ustalone przy rejestracji | Nie wyświetlane; działanie: „Zmień hasło” otwiera formularz osobny |
| Data utworzenia konta | Znacznik czasu rejestracji (Bezpieczeństwo i uwierzytelnianie, rozdz. 11.1) | Zapisana przy rejestracji | Tylko do odczytu |

### 3.3 Makieta sekcji Konto

```
Makieta 1 — Sekcja „Konto”

 ═════════════════════════════════════════════════════════════════════════════════
  Pasek kontekstu:  [Danaco Console] [Ubuntu] [Worktree] [Fable 5] [Ultra]
 ═════════════════════════════════════════════════════════════════════════════════
  Chat Window      │ Execution Loop ▼│ Sekcje ☰            │ Konto Operatora                       │
  Użytkownik ↔     │ Koordynator ↔   │                     │ ┌───────┐                             │
  Wykonawca        │ Wykonawca       │ ◈ Konto             │ │ awatar│  Zmień awatar               │
                   │                 │   Uwierzytelnianie  │ └───────┘                             │
  [ polecenie… ]   │ [ pętla: 0 zad.]│   Wygląd i język    │                                       │
                   │                 │   Urządzenia        │ Login                       [?] ⋮     │
                   │                 │   Powiadomienia     │ [ dnaharnowicz             ]          │
                   │                 │   Konta modeli      │                                       │
                   │                 │   Okna komunikacji  │ Adres e-mail uwierzytelniający  [?]   │
                   │                 │  Warstwy widoczności│ [ dnaharnowicz@danacogroup ] ✓ potw.  │
                   │                 │                     │                                       │
                   │                 │                     │ Hasło                                 │
                   │                 │                     │ ••••••••     [ Zmień hasło ]          │
                   │                 │                     │                                       │
                   │                 │                     │ Konto utworzone: 12.03.2026           │
                   │                 │                     │                                       │
                   │                 │                     │         [ Zapisz zmiany ]             │
 ═════════════════════════════════════════════════════════════════════════════════

Stan spoczynku: widoczne wyłącznie elementy warstwy 1 oraz zwinięte
wyzwalacze warstw 2–3 (`▼`, `⋮`, `☰`) i znaczniki paska kontekstu.
```

### 3.4 Zachowanie: zmiana hasła

Formularz „Zmień hasło” otwierany jest jako okno nakładkowe `.dn-modal` (System wizualny, rozdz. 8.9) i przyjmuje nowe hasło podane dwukrotnie. Pole hasła bieżącego pojawia się w formularzu, gdy Operator włączył w sekcji Uwierzytelnianie (rozdz. 4.1) ustawienie potwierdzania zmian hasłem bieżącym — zabezpieczenie dodatkowe, nie wymóg blokujący zapis. Nowe hasło zastępuje dotychczasowy skrót hasła przechowywany poza bazą danych, zgodnie z zasadą przyjętą dla wszystkich danych dostępowych (Bezpieczeństwo i uwierzytelnianie, rozdz. 10). W odróżnieniu od odzyskiwania konta (Bezpieczeństwo i uwierzytelnianie, rozdz. 7), które działa dla Operatora niemającego dostępu do żadnego uwierzytelnionego urządzenia, zmiana hasła z Ustawień jest czynnością Operatora już zalogowanego — polecenie komunikacji jest przez to odrębne (rozdz. 15.1).

```
Diagram 1 — Zmiana hasła z poziomu Ustawień

Operator zalogowany, otwiera sekcję „Konto” → „Zmień hasło”
        │
        ▼
Okno nakładkowe: podanie nowego hasła (dwukrotnie); pole hasła bieżącego widoczne,
gdy ustawienie „Wymagaj hasła bieżącego” jest włączone w Uwierzytelnianiu (rozdz. 4.1)
        │
        ▼
Kliknięcie „Zmień hasło” — przycisk aktywny niezależnie od stanu wypełnienia pól
        │
        ▼
Serwer sprawdza dane: przy włączonym ustawieniu weryfikuje hasło bieżące wobec skrótu
poza bazą (rozdz. 10 Bezpieczeństwa); zawsze sprawdza zgodność nowego hasła z powtórzeniem
        │
        ├── dane poprawne ──► zapis nowego skrótu hasła ──► dymek powiadomienia potwierdzenia
        │                                                    (zdarzenie auth.changed, rozdz. 15.2)
        └── ostrzeżenie (hasło bieżące niezgodne / nowe hasła się różnią / pole puste)
              ──► komunikat w oknie nakładkowym, formularz pozostaje otwarty do poprawy i ponowienia
```

---

## 4. Uwierzytelnianie

### 4.1 Zasada: wszystkie metody logowania jako ustawienie konfiguracyjne, hasło jako domyślna

Sekcja realizuje wprost mechanizm ustalony w Bezpieczeństwie i uwierzytelnianiu (rozdz. 6): Operator konfiguruje z tej sekcji cztery metody logowania, w tym hasło. Żadna z nich, hasło włącznie, nie jest aktywna z konieczności technicznej; każda jest włączana i wyłączana świadomą decyzją Operatora (zasada pełnej konfigurowalności, Koncepcja platformy, rozdz. 14, zasada 2). Hasło pozostaje metodą domyślną — aktywną od chwili rejestracji konta — lecz podobnie jak pozostałe metody może zostać przez Operatora wyłączone; przy próbie wyłączenia ostatniej aktywnej metody logowania sekcja wyświetla ostrzeżenie, nie blokadę (rozdz. 1.3).

| Metoda | Charakter | Przełącznik w Ustawieniach | Stan domyślny | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|
| Hasło | Bazowa, domyślna | `.dn-suwak` — konfigurowalny jak pozostałe metody | Aktywna | 2 | Sekcja „Uwierzytelnianie” w selektorze sekcji `☰` |
| Adres e-mail uwierzytelniający jako metoda logowania | Dodatkowa | `.dn-suwak` | Aktywna | 2 | Sekcja „Uwierzytelnianie” w selektorze sekcji `☰` |
| PIN | Dodatkowa | `.dn-suwak` | Nieaktywna | 2 | Sekcja „Uwierzytelnianie” w selektorze sekcji `☰` |
| Windows Hello | Dodatkowa, wiązana z konkretnym urządzeniem | `.dn-suwak` | Nieaktywna | 2 | Sekcja „Uwierzytelnianie” w selektorze sekcji `☰` |

Wyłączenie dowolnej metody dodatkowej nie ogranicza dostępu do platformy, dopóki pozostaje aktywna metoda bazowa albo inna metoda dodatkowa (Bezpieczeństwo i uwierzytelnianie, rozdz. 6.4). Wyłączenie hasła pozostawia dostęp przez metody dodatkowe wciąż aktywne; próba wyłączenia ostatniej aktywnej metody logowania nie jest blokowana, lecz sygnalizowana ostrzeżeniem (rozdz. 1.3). Ustawienie ma warstwę globalną i obowiązuje jednolicie na wszystkich urządzeniach Operatora.

### 4.2 Makieta sekcji Uwierzytelnianie

```
Makieta 2 — Sekcja „Uwierzytelnianie”

 ═════════════════════════════════════════════════════════════════════════════════
  Pasek kontekstu:  [Danaco Console] [Ubuntu] [Worktree] [Fable 5] [Ultra]
 ═════════════════════════════════════════════════════════════════════════════════
  Chat Window      │ Execution Loop ▼│ Sekcje ☰            │ Uwierzytelnianie                      │
  Użytkownik ↔     │ Koordynator ↔   │                     │                                       │
  Wykonawca        │ Wykonawca       │   Konto             │ Hasło                       [?] ⋮     │
                   │                 │ ◈ Uwierzytelnianie  │ metoda bazowa, domyślna   ( ● wł.)    │
  [ polecenie… ]   │ [ pętla: 0 zad.]│   Wygląd i język    │                                       │
                   │                 │   Urządzenia        │ Adres e-mail jako metoda    [?]       │
                   │                 │   Powiadomienia     │ logowania                 ( ● wł.)    │
                   │                 │   Konta modeli      │                                       │
                   │                 │   Okna komunikacji  │ PIN                         [?]       │
                   │                 │  Warstwy widoczności│ nieaktywny                ( ○ wył.)   │
                   │                 │                     │                                       │
                   │                 │                     │ Windows Hello               [?]       │
                   │                 │                     │ wiąże metodę z urządzeniem ( ○ wył.)  │
                   │                 │                     │                                       │
                   │                 │                     │ ┌──────────────────────────────────┐  │
                   │                 │                     │ │ Wymóg logowania: wyłączony ·      │ │
                   │                 │                     │ │ zmień w Konfiguracji (rozdz. 4.3) │ │
                   │                 │                     │ └──────────────────────────────────┘  │
 ═════════════════════════════════════════════════════════════════════════════════

Stan spoczynku: widoczne wyłącznie elementy warstwy 1 oraz zwinięte
wyzwalacze warstw 2–3 (`▼`, `⋮`, `☰`) i znaczniki paska kontekstu.
```

### 4.3 Wymóg logowania — ustawienie Konfiguracji, dostępne Operatorowi

Bezpieczeństwo i uwierzytelnianie (rozdz. 8) opisuje kontrolę dostępu do platformy. Zgodnie z zasadą zerowych blokad (rozdz. 1.3) egzekwowanie wymogu logowania w oknie logowania nie jest twardą bramką poza zasięgiem Operatora — jest ustawieniem, które Operator włącza i wyłącza w oknie Konfiguracji. Poniższa tabela streszcza mechanizm w zakresie istotnym dla tej sekcji.

| Cecha | Zachowanie | Warstwa | Sposób wywołania |
|---|---|---|---|
| Okno logowania | Obecne, pełny wygląd systemu wizualnego; wymusza uwierzytelnienie wtedy i tylko wtedy, gdy wymóg logowania jest włączony | 1 | Widoczne bez interakcji przy uruchomieniu platformy |
| Przycisk Pomiń w oknie logowania | Zawsze widoczny i klikalny; przy włączonym wymogu logowania jego użycie wyświetla komunikat i zawraca do okna logowania zamiast wpuszczać bez uwierzytelnienia | 1 | Widoczny bez interakcji w oknie logowania |
| Sekcja Uwierzytelnianie w Ustawieniach | W pełni funkcjonalna niezależnie od stanu wymogu logowania (Bezpieczeństwo i uwierzytelnianie, rozdz. 8.3) | 2 | Pozycja selektora sekcji `☰` |
| Przełącznik wymogu logowania | Ustawienie Operatora w oknie Konfiguracji, zmienialne w każdej chwili | 3 | Menu kontekstowe `⋮` sekcji Uwierzytelnianie — odnośnik do pozycji w oknie Konfiguracji; polecenie języka naturalnego w Chat Window |
| Plakietka stanu wymogu logowania | Informacja, czy okno logowania wymusza uwierzytelnienie | 1 | Widoczna bez interakcji w sekcji Uwierzytelnianie |

Przełącznik wymogu logowania **nie jest** pozycją tej sekcji ani żadnej innej sekcji Ustawień — jest ustawieniem okna Konfiguracji, zgodnie z rozgraniczeniem przyjętym w rozdz. 1.1 (Ustawienia: tożsamość Operatora; Konfiguracja: zachowanie platformy). Pozostaje jednak w pełni dostępny i zmienialny przez Operatora (zasada zerowych blokad, rozdz. 1.3) — nie jest parametrem poza jego zasięgiem. Sekcja Uwierzytelnianie prezentuje informacyjną plakietkę stanu (`.dn-plakietka--informacja`), tak aby Operator wiedział, czy okno logowania obecnie wymusza uwierzytelnienie, wraz z odnośnikiem do pozycji w Konfiguracji, gdzie może tę wartość zmienić. Operator zmienia tę wartość również poleceniem języka naturalnego skierowanym do Wykonawcy w Chat Window (rozdz. 9.1).

### 4.4 Urządzenia połączone — odesłanie

Wykaz urządzeń i unieważnianie ich tokenów dostępu jest zawartością sekcji Urządzenia (rozdz. 6), nie sekcji Uwierzytelnianie — mimo że oba mechanizmy opisuje wspólnie Bezpieczeństwo i uwierzytelnianie (rozdz. 3, rozdz. 13). Rozdzielenie w Ustawieniach odpowiada różnicy pytań: Uwierzytelnianie odpowiada na „jakimi metodami mogę się zalogować”, Urządzenia — na „które moje urządzenia mają dziś dostęp”.

---

## 5. Wygląd: motyw i język

### 5.1 Motyw wizualny

Motyw jest sterowany atrybutem `data-theme="light|dark"` na elemencie głównym dokumentu (System wizualny, rozdz. 9); każdy token semantyczny ma zdefiniowaną wartość dla obu wariantów, więc przełączenie motywu zmienia wyłącznie wartości zmiennych, nie reguły komponentów.

| Ustawienie | Opis | Warstwa konfiguracji | Wartość domyślna | Warstwa widoczności | Sposób wywołania |
|---|---|---|---|---|---|
| Motyw wizualny | Wariant jasny albo ciemny systemu wizualnego marki | globalna | brak jawnego ustawienia → wartość z preferencji systemowej urządzenia (`prefers-color-scheme`) | 2 | Sekcja „Wygląd i język” w selektorze sekcji `☰`; skrót klawiszowy przełączenia motywu |

### 5.2 Język interfejsu

| Ustawienie | Opis | Warstwa konfiguracji | Wartość domyślna | Warstwa widoczności | Sposób wywołania |
|---|---|---|---|---|---|
| Język interfejsu | Język wyświetlania interfejsu platformy | globalna | polski | 2 | Sekcja „Wygląd i język” w selektorze sekcji `☰` |

Polski jest językiem wbudowanym i domyślnym. Pole jest listą wyboru zasilaną katalogiem języków platformy, nie parą stałych wartości — dodanie języka do katalogu nie wymaga zmiany tego elementu interfejsu (Architektura, rozdz. 14).

### 5.3 Makieta sekcji Wygląd i język

```
Makieta 3 — Sekcja „Wygląd i język”

 ═════════════════════════════════════════════════════════════════════════════════
  Pasek kontekstu:  [Danaco Console] [Ubuntu] [Worktree] [Fable 5] [Ultra]
 ═════════════════════════════════════════════════════════════════════════════════
  Chat Window      │ Execution Loop ▼│ Sekcje ☰            │ Wygląd i język                        │
  Użytkownik ↔     │ Koordynator ↔   │                     │                                       │
  Wykonawca        │ Wykonawca       │   Konto             │ Motyw wizualny              [?] ⋮     │
                   │                 │   Uwierzytelnianie  │ ┌────────┬─────────┬────────┐         │
  [ polecenie… ]   │ [ pętla: 0 zad.]│ ◈ Wygląd i język    │ │ ☀ Jasny│ 🌙 Ciemny│ ⚙ System│        │
                   │                 │   Urządzenia        │ └────────┴─────────┴────────┘         │
                   │                 │   Powiadomienia     │   .dn-zakladki--pigulki, wybór        │
                   │                 │   Konta modeli      │   zaznaczony                          │
                   │                 │   Okna komunikacji  │                                       │
                   │                 │  Warstwy widoczności│ Język interfejsu            [?]       │
                   │                 │                     │ [ Polski                  ▾ ]         │
                   │                 │                     │   .dn-pole-kontrolka                  │
 ═════════════════════════════════════════════════════════════════════════════════

Stan spoczynku: widoczne wyłącznie elementy warstwy 1 oraz zwinięte
wyzwalacze warstw 2–3 (`▼`, `⋮`, `☰`) i znaczniki paska kontekstu.
```

### 5.4 Zachowanie

Zmiana motywu i języka jest natychmiastowa i widoczna na bieżąco we wszystkich otwartych oknach operacyjnych — bez przeładowania karty sesji, ponieważ obie pozycje należą do zakresu Aplikacja, nie do zakresu żadnej sesji (Model konfiguracji, rozdz. 4.1). Zmiana motywu wykorzystuje przejście animowane w czasie tokenu `--dn-czas-2`, z poszanowaniem `prefers-reduced-motion` (System wizualny, rozdz. 9).

---

## 6. Urządzenia: lista, parowanie Mobile, synchronizacja

### 6.1 Wykaz urządzeń połączonych

Sekcja prezentuje wykaz encji Urządzenie powiązanych z jedynym Kontem (Bezpieczeństwo i uwierzytelnianie, rozdz. 11.2), w formie tabeli `.dn-tabela` z kolumną akcji — wzorzec okien zarządców (System wizualny, rozdz. 10.7).

| Ustawienie | Opis | Warstwa konfiguracji | Wartość domyślna | Warstwa widoczności | Sposób wywołania |
|---|---|---|---|---|---|
| Urządzenia połączone | Wykaz urządzeń powiązanych z Kontem oraz zarządzanie ich dostępem | globalna | brak ograniczeń liczby urządzeń | 2 | Sekcja „Urządzenia” w selektorze sekcji `☰` |
| Odłączenie urządzenia | Unieważnienie tokenu dostępu wskazanego urządzenia | globalna | brak | 3 | Menu akcji wiersza `⋮` w tabeli urządzeń |
| Parowanie urządzenia mobilnego | Wydanie poświadczenia nowemu urządzeniu | globalna | brak sparowanych urządzeń mobilnych | 3 | Działanie „Sparuj urządzenie” w sekcji Urządzenia; polecenie języka naturalnego w Chat Window |

| Kolumna tabeli | Zawartość |
|---|---|
| Urządzenie | Nazwa urządzenia + ikona typu (komputer, telefon, tablet) |
| Identyfikator | Skrócony identyfikator techniczny urządzenia |
| Ostatnie połączenie | Data i godzina ostatniego połączenia z platformą |
| Stan | Plakietka: „to urządzenie” dla urządzenia bieżącego; „aktywne” dla pozostałych |
| Akcje | `.dn-btn-ikona` — „Odłącz” (unieważnienie tokenu) |

### 6.2 Parowanie urządzenia mobilnego

Przycisk „Sparuj urządzenie” uruchamia protokół parowania opisany w całości w dokumencie funkcji globalnej ([Mobile](../funkcje-globalne/mobile.md), rozdz. 9). Parowanie jest drogą wejścia urządzenia przenośnego do funkcji Mobile, zakładane w obecności Operatora przy maszynie już uwierzytelnionej, i wydaje nowemu urządzeniu poświadczenie odrębne od pozostałych — odwołanie jednego urządzenia nie odcina pozostałych. Niniejszy rozdział dokumentuje interfejs parowania po stronie okna Ustawień: okno nakładkowe kodu, potwierdzenie zgłoszenia i miejsce obu kroków w sekcji Urządzenia.

```
Diagram 2 — Parowanie z sekcji „Urządzenia”

Sekcja „Urządzenia” → [ Sparuj urządzenie ]        (device.pair.start)
        │
        ▼
Okno nakładkowe kodu parowania: ciąg znaków + kod graficzny + licznik ważności
        │
        ▼
Operator odczytuje kod aplikacją na urządzeniu przenośnym  (device.pair.confirm)
        │
        ▼
Okno nakładkowe potwierdzenia zgłoszenia: typ, nazwa wstępna, czas zgłoszenia
        │
        ├── [ Potwierdź ] ──► wydanie poświadczenia (device.pair.approve)
        │                      nowy wiersz w tabeli urządzeń
        │                      dymek powiadomienia „Sparowano: <nazwa urządzenia>”
        │                      (zdarzenie device.changed, rozdz. 14.2)
        │
        ├── [ Odrzuć ] ──────► kod unieważniony (device.pair.reject)
        │
        └── czas ważności upłynął ──► dymek powiadomienia błędu, kod unieważniony,
                                      działanie „Wydaj nowy kod”
```

| Element interfejsu parowania | Co to jest | Do czego służy | Stany | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|
| Okno nakładkowe kodu parowania | Okno nakładkowe `.dn-modal` z kodem w dwóch postaciach | Przekazanie kodu urządzeniu przenośnemu | Generowanie · kod aktywny · kod wygasły · kod unieważniony | 3 | Przycisk „Sparuj urządzenie” |
| Kod w postaci znakowej | Ciąg znaków do wpisania ręcznego | Parowanie urządzenia bez odczytu kamerą | Domyślny · skopiowany | 3 | Widoczny w oknie nakładkowym |
| Kod graficzny | Kod do odczytu kamerą urządzenia przenośnego | Parowanie bez wpisywania znaków | Domyślny · wygasły | 3 | Widoczny w oknie nakładkowym |
| Licznik ważności kodu | Wskaźnik pozostałego czasu ważności | Sygnał, ile czasu pozostaje na odczyt kodu | Odliczanie · wyzerowany | 3 | Widoczny w oknie nakładkowym |
| Działanie „Unieważnij kod” | Przycisk pomocniczy okna nakładkowego | Zamknięcie parowania przed jego użyciem | Domyślny · ładowanie | 3 | Widoczne w oknie nakładkowym |
| Działanie „Wydaj nowy kod” | Przycisk pomocniczy okna nakładkowego | Powtórzenie kroku wydania kodu po jego wygaśnięciu | Domyślny · ładowanie | 3 | Widoczne po wygaśnięciu kodu |
| Okno nakładkowe potwierdzenia zgłoszenia | Okno nakładkowe z danymi urządzenia zgłaszającego | Rozstrzygnięcie, czy urządzenie wchodzi do rejestru | Oczekuje · potwierdzone · odrzucone | 3 | Zdarzenie zgłoszenia urządzenia |
| Pole nazwy urządzenia w potwierdzeniu | Pole `.dn-pole-kontrolka` z nazwą wstępną | Nadanie nazwy rozpoznawalnej w rejestrze | Domyślne · zmienione | 3 | Widoczne w oknie nakładkowym potwierdzenia |
| Komunikat kontekstowy błędu parowania | Wiersz komunikatu w oknie nakładkowym | Wskazanie przyczyny: kod niepoprawny, wygasły, użyty, zgłoszenie odrzucone | Ukryty · widoczny | 3 | Pojawia się przy warunku błędu ([Mobile](../funkcje-globalne/mobile.md), rozdz. 9.11) |

### 6.2a Zarządzanie sparowanymi urządzeniami

Tabela urządzeń sekcji (rozdz. 6.1) jest zarazem rejestrem sparowanych urządzeń. Poza polami wykazu prezentuje stan poświadczenia każdego urządzenia i udostępnia operacje zarządzania dostępem.

| Kolumna rejestru | Zawartość | Źródło |
|---|---|---|
| Tryb Mobile | Plakietka „Mobile” przy urządzeniu korzystającym z funkcji globalnej | `urzadzenie.tryb_mobile_aktywny` |
| Data parowania | Data wydania poświadczenia | `poswiadczenie_urzadzenia.data_wydania` (Model danych, rozdz. 3.4) |
| Ważność poświadczenia | Data wygaśnięcia poświadczenia | `poswiadczenie_urzadzenia.data_waznosci` |
| Stan poświadczenia | Plakietka: aktywne · zbliża się wygaśnięcie · wygasłe · odwołane | `poswiadczenie_urzadzenia.stan` |

| Operacja | Opis | Warstwa | Sposób wywołania |
|---|---|---|---|
| Zmiana nazwy urządzenia | Nadanie nazwy rozpoznawalnej w rejestrze (`device.rename`) | 3 | Menu akcji wiersza `⋮` |
| Odwołanie dostępu | Unieważnienie poświadczenia wskazanego urządzenia i zamknięcie jego połączeń — wylogowanie zdalne (`device.revoke`) | 3 | Menu akcji wiersza `⋮`; działanie „Odłącz” w kolumnie akcji |
| Odwołanie dostępu wszystkich urządzeń | Unieważnienie wszystkich poświadczeń wydanych przed rozstrzygnięciem | 3 | Akcja zbiorcza nagłówka tabeli |
| Odświeżenie poświadczenia | Wydanie poświadczenia o nowym terminie ważności — brak dedykowanej komendy kontraktu; `[DO DECYZJI OPERATORA]`, czy działanie uruchamia nowe parowanie (`device.pair.start`) czy zyskuje własną komendę | 4 | Polecenie języka naturalnego w Chat Window; wyszukiwarka funkcji |

Postępowanie przy utracie urządzenia — rozpoznanie wiersza po nazwie, typie i dacie ostatniego połączenia, odwołanie poświadczenia i zapis w dzienniku audytu — opisuje opracowanie [Mobile](../funkcje-globalne/mobile.md), rozdz. 9.9.

### 6.3 Synchronizacja

Synchronizacja odbywa się na żywo kanałem WebSocket — zmiana dokonana na jednym urządzeniu jest przekazywana przez serwer do pozostałych (Architektura, rozdz. 15). Nie jest to ustawienie do włączenia, lecz stała właściwość platformy; sekcja prezentuje wyłącznie jej stan bieżący jako informację, nie jako przełącznik.

| Element | Forma | Treść |
|---|---|---|
| Wskaźnik stanu synchronizacji | `.dn-plakietka--sygnal` z `.dn-kropka` | „Synchronizacja aktywna” (kropka `--sukces`) / „Brak połączenia — zmiany zostaną zsynchronizowane po powrocie łączności” (kropka `--ostrz`) |

### 6.4 Makieta sekcji Urządzenia

```
Makieta 4 — Sekcja „Urządzenia”

 ═════════════════════════════════════════════════════════════════════════════════
  Pasek kontekstu:  [Danaco Console] [Ubuntu] [Worktree] [Fable 5] [Ultra]
 ═════════════════════════════════════════════════════════════════════════════════
  Chat Window      │ Execution Loop ▼│ Sekcje ☰            │ Urządzenia      ● Synchronizacja      │
  Użytkownik ↔     │ Koordynator ↔   │                     │                   aktywna             │
  Wykonawca        │ Wykonawca       │   Konto             │                                       │
                   │                 │   Uwierzytelnianie  │ ┌────────────────────────────────┐    │
  [ polecenie… ]   │ [ pętla: 0 zad.]│   Wygląd i język    │ │ Urządzenie   │Ostatnie poł.│ ⋮ │    │
                   │                 │ ◈ Urządzenia        │ ├────────────────────────────────┤    │
                   │                 │   Powiadomienia     │ │ 💻 DP-Biuro  │teraz        │ ⋮ │     │
                   │                 │   Konta modeli      │ │ 📱 iPhone    │dziś 08:12   │ ⋮ │     │
                   │                 │   Okna komunikacji  │ │ 📱 Pixel     │3 dni temu   │ ⋮ │     │
                   │                 │  Warstwy widoczności│ └────────────────────────────────┘    │
                   │                 │                     │   .dn-tabela, akcje w menu ⋮          │
                   │                 │                     │                                       │
                   │                 │                     │       [ Sparuj urządzenie ]           │
 ═════════════════════════════════════════════════════════════════════════════════

Stan spoczynku: widoczne wyłącznie elementy warstwy 1 oraz zwinięte
wyzwalacze warstw 2–3 (`▼`, `⋮`, `☰`) i znaczniki paska kontekstu.
```

---

## 7. Powiadomienia

### 7.1 Zakres ustawienia źródłowego

Sekcja steruje jednym mechanizmem powiadamiania obowiązującym w całej platformie — centrum powiadomień, którego taksonomię klas zdarzeń, postać wizualną, umiejscowienie w układzie pionowym, warstwę widoczności, sposób wywołania, stany i działania opisuje karta komponentu w [Katalogu komponentów](katalog-komponentow.md) (rozdz. 11.6). Encję zdarzenia opisuje [Model danych](../architektura/model-danych.md) (rozdz. 18.4).

| Ustawienie | Opis | Warstwa konfiguracji | Wartość domyślna | Warstwa widoczności | Sposób wywołania |
|---|---|---|---|---|---|
| Powiadomienia | Zakres klas zdarzeń i kanał dostarczenia powiadomień centrum powiadomień, w tym zdarzeń nadzorowanych przez funkcję Mobile (Koncepcja platformy, rozdz. 12.1; [Mobile](../funkcje-globalne/mobile.md), rozdz. 8) | globalna, sesja | powiadomienia aktywne | 2 | Sekcja „Powiadomienia” w selektorze sekcji `☰` |

Sekcja rozwija tę jedną pozycję Modelu konfiguracji do poziomu elementów interfejsu: przełącznika głównego, klas zdarzeń i kanału dostarczenia.

### 7.2 Klasy zdarzeń

Klasy zdarzeń odpowiadają taksonomii centrum powiadomień ([Katalog komponentów](katalog-komponentow.md), rozdz. 11.6) i zdarzeniom zdefiniowanym w kontraktach komunikacji poszczególnych zakresów platformy — sekcja nie wprowadza nowych zdarzeń, lecz udostępnia Operatorowi wybór, które z nich mają skutkować powiadomieniem w centrum powiadomień i na urządzeniu mobilnym.

| Klasa zdarzenia | Przykładowe zdarzenie źródłowe | Domyślnie | Warstwa | Sposób wywołania |
|---|---|---|---|---|
| Zakończenie | Zakończenie przebiegu pętli wykonawczej, zakończenie procesu automatyki — Execution Monitor, moduł Automations (Specyfikacja okien operacyjnych, rozdz. 6.3) | aktywne | 2 | Tabela klas zdarzeń w sekcji „Powiadomienia” |
| Decyzja | Krok procesu oczekujący na zatwierdzenie — Monitor procesu, panel orkiestracji (Koncepcja platformy, rozdz. 13.8) | aktywne | 2 | Tabela klas zdarzeń w sekcji „Powiadomienia” |
| Błąd | Niepowodzenie zadania, naruszenie zależności orkiestracji, powtarzające się niepowodzenie walidacji | aktywne | 2 | Tabela klas zdarzeń w sekcji „Powiadomienia” |
| Wzmianka | Odwołanie do Operatora w treści, komentarz przypisany do artefaktu | aktywne | 2 | Tabela klas zdarzeń w sekcji „Powiadomienia” |
| Termin | Termin zadania projektu, zbliżający się cykl harmonogramu | aktywne | 2 | Tabela klas zdarzeń w sekcji „Powiadomienia” |
| Automatyka | Wpięcie automatyki, wynik cyklu harmonogramu | aktywne | 2 | Tabela klas zdarzeń w sekcji „Powiadomienia” |
| System | `auth.changed`, `device.changed` (Bezpieczeństwo i uwierzytelnianie, rozdz. 12.2); `config.changed`, `memory.changed`, `history.changed` (Model konfiguracji, rozdz. 8); wiadomości i zdolności telefonu ([Mobile](../funkcje-globalne/mobile.md), rozdz. 8.1) | aktywne | 2 | Tabela klas zdarzeń w sekcji „Powiadomienia” |

### 7.3 Kanał dostarczenia

Centrum powiadomień jest kanałem podstawowym każdej klasy zdarzeń — każde zdarzenie objęte ustawieniem trafia do rejestru centrum niezależnie od pozostałych kanałów. Kolumna kanału określa dodatkowe drogi dostarczenia tego samego zdarzenia.

| Kanał | Opis | Warstwa | Sposób wywołania |
|---|---|---|---|
| Centrum powiadomień | Pozycja w kolumnie centrum powiadomień oraz Dymek powiadomienia (`.dn-toast`, System wizualny, rozdz. 8.10) w chwili wystąpienia zdarzenia | 2 | Kolumna kanału w tabeli klas zdarzeń |
| Mobile — powiadomienie wypychane | Dostarczane dwoma kanałami o różnej poufności: kanałem własnym przy aplikacji na pierwszym planie oraz kanałem pośrednika poza nią — pełny mechanizm w opracowaniu [Mobile](../funkcje-globalne/mobile.md), rozdz. 8.5 | 2 | Kolumna kanału w tabeli klas zdarzeń |
| E-mail | Na adres e-mail uwierzytelniający (rozdz. 3.2) | 2 | Kolumna kanału w tabeli klas zdarzeń |

### 7.4 Makieta sekcji Powiadomienia

```
Makieta 5 — Sekcja „Powiadomienia”

 ═════════════════════════════════════════════════════════════════════════════════
  Pasek kontekstu:  [Danaco Console] [Ubuntu] [Worktree] [Fable 5] [Ultra]
 ═════════════════════════════════════════════════════════════════════════════════
  Chat Window      │ Execution Loop ▼│ Sekcje ☰            │ Powiadomienia                         │
  Użytkownik ↔     │ Koordynator ↔   │                     │                                       │
  Wykonawca        │ Wykonawca       │   Konto             │ Powiadomienia aktywne       [?] ⋮     │
                   │                 │   Uwierzytelnianie  │                           ( ● wł.)    │
  [ polecenie… ]   │ [ pętla: 0 zad.]│   Wygląd i język    │                                       │
                   │                 │   Urządzenia        │ Klasy zdarzeń        Kanał            │
                   │                 │ ◈ Powiadomienia     │ ──────────────────────────────────    │
                   │                 │   Konta modeli      │ Zakończenie          [●] centrum      │
                   │                 │   Okna komunikacji  │                      [●] Mobile       │
                   │                 │  Warstwy widoczności│ Decyzja              [●] centrum      │
                   │                 │                     │                      [●] Mobile       │
                   │                 │                     │ Błąd                 [●] centrum      │
                   │                 │                     │                      [●] Mobile       │
                   │                 │                     │ Wzmianka             [●] centrum      │
                   │                 │                     │ Termin               [●] centrum      │
                   │                 │                     │ Automatyka           [●] centrum      │
                   │                 │                     │ System               [●] centrum      │
                   │                 │                     │                      [●] e-mail       │
                   │                 │                     │   .dn-tabela, .dn-check w komórce     │
 ═════════════════════════════════════════════════════════════════════════════════

Stan spoczynku: widoczne wyłącznie elementy warstwy 1 oraz zwinięte
wyzwalacze warstw 2–3 (`▼`, `⋮`, `☰`) i znaczniki paska kontekstu.
```

### 7.5 Zachowanie

Wyłączenie przełącznika głównego wygasza wszystkie klasy zdarzeń jednocześnie, zachowując ich indywidualne ustawienia w pamięci — ponowne włączenie przywraca poprzedni stan każdej klasy, a nie stan domyślny. Zmiana pojedynczej klasy zdarzenia lub kanału jest zapisywana natychmiast, bez przycisku zbiorczego zapisu, zgodnie z konwencją przyjętą dla macierzy przełączników w oknie konfiguracji punktów izolacji (Koncepcja platformy, rozdz. 6.4).

---

## 8. Zarządzanie kontami modeli

### 8.1 Miejsce w mechanizmie kanału modelu

Integracja modeli (rozdz. 5) definiuje kanał modelu jako encję reprezentującą skonfigurowane połączenie z konkretnym modelem przez jeden z czterech kanałów integracji: API, CLI, SSH, HTTP. Dane dostępowe kanału są przechowywane poza bazą danych, a encja kanału modelu przechowuje wyłącznie odwołanie do nich (Integracja modeli, rozdz. 5; Bezpieczeństwo i uwierzytelnianie, rozdz. 10) — dokładnie ten sam wzorzec, jaki niniejszy dokument stosuje w rozdz. 3 do hasła Operatora.

Sekcja „Konta modeli” jest katalogiem Operatora dla jednego szczególnego przypadku kanału CLI — narzędzia wiersza poleceń Code CLI, którym platforma może się posługiwać jako jednym z czterech kanałów integracji modeli. Podczas gdy przypisanie kanału modelu do sesji lub roli pozostaje czynnością wykonywaną w oknie Konfiguracji, zakres „Zachowanie modeli” (Model konfiguracji, rozdz. 5.4) — utrzymywanie samego wykazu dostępnych kont, ich danych dostępowych i katalogów profili jest czynnością właściwą Operatorowi jako takiemu, a nie pojedynczej sesji, i dlatego należy do Ustawień.

```
Schemat 4 — Relacja sekcji „Konta modeli” do kanału modelu (Integracja modeli, rozdz. 5, 7, 8)

USTAWIENIA → Konta modeli                        KONFIGURACJA → Zachowanie modeli
┌──────────────────────────────┐                  ┌──────────────────────────────┐
│ Katalog kont Code CLI        │   konto wybrane  │ Kanał modelu przypisany do   │
│ (niniejszy rozdział)         │  jako dane dostę-│ sesji lub roli               │
│                              │  powe kanału CLI │ (Model konfiguracji, rozdz.  │
│  • nazwa konta               │ ───────────────► │  5.4)                        │
│  • katalog profilu           │                  │                              │
│  • dane dostępowe (poza bazą)│                  │  typ kanału: CLI             │
│  • konto aktywne             │                  │ odwołanie → konto z tej listy│
└──────────────────────────────┘                  └──────────────────────────────┘
   „kim się logujemy do CLI”                          „które zadanie tego używa”
```

### 8.2 Zawartość sekcji

| Ustawienie | Opis | Wartość domyślna | Warstwa | Sposób wywołania |
|---|---|---|---|---|
| Wykaz kont modeli | Nazwane, zapisane konta wiersza poleceń (w tym konta narzędzia Code CLI) wraz z katalogiem profilu i stanem | pusty do chwili dodania pierwszego konta | 2 | Sekcja „Konta modeli” w selektorze sekcji `☰` |
| Konto aktywne | Konto stosowane domyślnie dla nowo tworzonych kanałów CLI, dopóki sesja lub rola nie wskaże innego wprost | pierwsze dodane konto | 3 | Menu akcji wiersza `⋮` → „Ustaw jako aktywne”; polecenie języka naturalnego w Chat Window |
| Dane dostępowe konta | Token narzędzia wiersza poleceń przechowywany poza bazą danych | wartość podana przy dodaniu konta | 4 | Tryb administracyjny, wyszukiwarka funkcji, polecenie języka naturalnego w Chat Window |

| Kolumna tabeli kont | Zawartość |
|---|---|
| Nazwa konta | Nazwa własna nadana przez Operatora (np. „Danaco — konto główne”, „Zespół prawny — konto zapasowe”) |
| Katalog profilu | Ścieżka lokalnego katalogu przechowującego dane logowania i pliki konfiguracyjne tego konta, odrębna dla każdego konta |
| Stan | Plakietka: „aktywne” / „skonfigurowane” / „błąd uwierzytelnienia” |
| Ostatnio użyte | Data ostatniego wywołania modelu przez to konto |
| Akcje | `.dn-btn-ikona` — „Ustaw jako aktywne”, „Testuj połączenie”, „Edytuj”, „Usuń” |

### 8.3 Dodanie konta — okno nakładkowe

Dodanie konta otwiera okno nakładkowe `.dn-modal` z formularzem czterech pól: nazwa konta, katalog profilu, dane dostępowe (token narzędzia wiersza poleceń) oraz działanie testu połączenia wykonywane przed zapisem. Po zapisie dane dostępowe trafiają do magazynu danych dostępowych poza bazą (Bezpieczeństwo i uwierzytelnianie, rozdz. 10); pole danych dostępowych, tak w formularzu jak i przy późniejszym podglądzie, wyświetla wartość domyślnie w postaci jawnej — maskowanie jest ustawieniem konfiguracyjnym włączanym przez Operatora, z zawsze dostępnym działaniem „Pokaż” / „Ukryj” (zasada zerowych blokad, rozdz. 1.3).

```
Makieta 6 — Okno nakładkowe „Dodaj konto modelu”

┌──────────────────────────────────────────────────┐
│  Dodaj konto modelu                      ✕       │
├──────────────────────────────────────────────────┤
│  Nazwa konta                             [?]     │
│  [ np. Danaco — konto główne            ]        │
│                                                  │
│  Katalog profilu                         [?]     │
│  [ ~/.danaco/profile-danaco-glowne      ]        │
│                                                  │
│  Dane dostępowe (token CLI)              [?]     │
│  [ <token w postaci jawnej>             ] 👁      │
│  jawne domyślnie · przechowywane poza bazą       │
│  danych (rozdz. 8.1) · maskowanie sterowane      │
│  ustawieniem konfiguracyjnym                     │
│                                                  │
│  [ Testuj połączenie ]     stan: — nie testowano │
│  warstwa 3 — rozwinięcie kontekstowe             │
├──────────────────────────────────────────────────┤
│                    [ Anuluj ]  [ Zapisz konto ]  │
│                                   .dn-btn--sygnal│
└──────────────────────────────────────────────────┘
```

### 8.4 Przełączanie aktywnego konta — zasada jawności

Zgodnie z zasadą platformy, że zmiana kontekstu jest zawsze świadomym, widocznym przejściem, a nie ukrytą zmianą stanu (Koncepcja platformy, rozdz. 1), przełączenie konta aktywnego jest wyłącznie jawną, ręczną akcją Operatora — kliknięciem „Ustaw jako aktywne” przy wybranym wierszu tabeli. Sekcja nie zawiera mechanizmu automatycznego, niewidocznego przełączania kont między wywołaniami modelu.

```
Diagram 3 — Przełączenie konta aktywnego

Operator klika „Ustaw jako aktywne” przy koncie B
        │
        ▼
Poprzednie konto aktywne (A) traci plakietkę „aktywne” → „skonfigurowane”
        │
        ▼
Konto B otrzymuje plakietkę „aktywne”
        │
        ▼
Zdarzenie account.changed rozgłoszone do wszystkich urządzeń (rozdz. 15.2)
        │
        ▼
Nowe sesje i role wskazujące kanał CLI bez własnego, jawnego przypisania konta
(Model konfiguracji, rozdz. 5.4) korzystają odtąd z konta B; sesje już
uruchomione zachowują konto, z którym zostały uruchomione
```

### 8.5 Makieta sekcji Konta modeli

```
Makieta 7 — Sekcja „Konta modeli”

 ═════════════════════════════════════════════════════════════════════════════════
  Pasek kontekstu:  [Danaco Console] [Ubuntu] [Worktree] [Fable 5] [Ultra]
 ═════════════════════════════════════════════════════════════════════════════════
  Chat Window      │ Execution Loop ▼│ Sekcje ☰            │ Konta modeli                          │
  Użytkownik ↔     │ Koordynator ↔   │                     │                                       │
  Wykonawca        │ Wykonawca       │   Konto             │ ┌────────────────────────────────┐    │
                   │                 │   Uwierzytelnianie  │ │Nazwa         │Stan     │  ⋮   │     │
  [ polecenie… ]   │ [ pętla: 0 zad.]│   Wygląd i język    │ ├────────────────────────────────┤    │
                   │                 │   Urządzenia        │ │Danaco — główne│●aktywne │  ⋮   │    │
                   │                 │   Powiadomienia     │ │Zespół prawny │skonfig. │  ⋮   │     │
                   │                 │ ◈ Konta modeli      │ │Zapasowe      │błąd     │  ⋮   │     │
                   │                 │   Okna komunikacji  │ └────────────────────────────────┘    │
                   │                 │  Warstwy widoczności│   .dn-tabela, akcje w menu ⋮          │
                   │                 │                     │                                       │
                   │                 │                     │       [ Dodaj konto modelu ]          │
                   │                 │                     │                                       │
                   │                 │                     │ Przypisanie kanałów do sesji i ról:   │
                   │                 │                     │ Konfiguracja → Zachowanie modeli →    │
 ═════════════════════════════════════════════════════════════════════════════════

Stan spoczynku: widoczne wyłącznie elementy warstwy 1 oraz zwinięte
wyzwalacze warstw 2–3 (`▼`, `⋮`, `☰`) i znaczniki paska kontekstu.
```

### 8.6 Stan pustej listy

Przed dodaniem pierwszego konta sekcja prezentuje `.dn-pusty-stan` (System wizualny, rozdz. 8.12): ikona `kod`, tytuł „Brak skonfigurowanych kont modeli”, opis „Dodaj konto narzędzia Code CLI lub innego narzędzia wiersza poleceń, aby platforma mogła się nim posługiwać jako kanałem CLI”, przycisk „Dodaj konto modelu”.

---

## 9. Okna komunikacji operacyjnej i ich ustawienia

Platforma prowadzi pracę dwoma kanałami komunikacji operacyjnej obecnymi w każdym kontekście pracy użytkownika, w tym przy otwartym oknie Ustawień. Oba kanały są elementami pierwszoplanowymi architektury interfejsu, nie funkcjami dodatkowymi. Sekcja „Okna komunikacji” gromadzi ustawienia sterujące oboma oknami.

### 9.1 Chat Window — kanał Użytkownik ↔ Wykonawca

Chat Window jest głównym oknem komunikacji między Użytkownikiem a Wykonawcą (AI, agent lub system wykonawczy). Stanowi centralny punkt pracy użytkownika i podstawowy mechanizm sterowania wszystkimi procesami realizowanymi przez platformę: przyjmuje polecenia w języku naturalnym, prezentuje strumień odpowiedzi i wyników, umożliwia zatwierdzanie oraz przerywanie działań, a także wyjaśnianie wyniku i kontekstu. Okno zajmuje lewą kolumnę obszaru roboczego — stałą, o pełnej wysokości — i występuje w tym samym miejscu układu we wszystkich modułach i środowiskach. Każde ustawienie opisane w rozdziałach 3–10 Operator zmienia zarówno elementem interfejsu w panelu zawartości, jak i poleceniem języka naturalnego skierowanym do Wykonawcy w tym oknie.

### 9.2 Execution Loop Window — kanał Koordynator ↔ Wykonawca

Execution Loop Window jest drugim kanałem komunikacji operacyjnej i prezentuje wymianę między Koordynatorem a Wykonawcą. Odpowiada za prowadzenie pętli wykonawczej, koordynację zadań, nadzór nad realizacją, orkiestrację działań oraz kontrolę realizacji procesów. Zawiera bieżące zlecenie i jego dekompozycję na zadania, kolejkę i stan zadań, wymianę komunikatów sterujących, wyniki kontroli jakości i decyzje o ponowieniu, wskaźniki przebiegu pętli oraz sterowanie przebiegiem: wstrzymanie, wznowienie, przerwanie i korektę zlecenia. Okno otwierane jest jako kolumna sąsiadująca z Chat Window.

Role: **Użytkownik** zleca i zatwierdza; **Koordynator** jest komponentem orkiestrującym platformy — dekomponuje zlecenie, przydziela i nadzoruje zadania; **Wykonawca** — AI, agent lub system wykonawczy — realizuje zadania.

### 9.3 Sekcja „Okna komunikacji” — ustawienia obu kanałów

| Ustawienie | Opis | Wartość domyślna | Warstwa | Sposób wywołania |
|---|---|---|---|---|
| Zapamiętane położenie paska Chat Window | Szerokość lewej kolumny obszaru roboczego nie jest nastawą wpisywaną suwakiem — Operator ustawia ją przesunięciem uchwytu rozdzielającego kolumny, a okno zapamiętuje ustawione położenie paska między sesjami; położenie kolumny pozostaje stałe | domyślne położenie startowe systemu wizualnego | 2 | Uchwyt rozdzielający kolumny obszaru roboczego |
| Gęstość strumienia Chat Window | Zwarta albo rozstrzelona prezentacja wiadomości w strumieniu | zwarta | 2 | Sekcja „Okna komunikacji” w selektorze sekcji `☰` |
| Potwierdzanie działań w Chat Window | Zakres działań Wykonawcy wymagających zatwierdzenia Użytkownika przed wykonaniem | działania zmieniające stan platformy | 2 | Sekcja „Okna komunikacji” w selektorze sekcji `☰` |
| Sterowanie ustawieniami z Chat Window | Przyjmowanie poleceń języka naturalnego zmieniających pozycje Ustawień | włączone | 2 | Sekcja „Okna komunikacji” w selektorze sekcji `☰` |
| Otwarcie Execution Loop Window | Stan kolumny sąsiadującej: otwarta albo zwinięta do wyzwalacza `▼` | zwinięta, otwierana samoczynnie przy rozpoczęciu pętli | 2 | Wyzwalacz `▼` „Execution Loop”; polecenie języka naturalnego w Chat Window |
| Zapamiętane położenie paska Execution Loop Window | Szerokość kolumny sąsiadującej ustawiana wyłącznie przesunięciem uchwytu rozdzielającego; okno zapamiętuje ustawione położenie paska między sesjami, bez pola ani suwaka nastawy | domyślne położenie startowe systemu wizualnego | 2 | Uchwyt rozdzielający kolumny obszaru roboczego |
| Zakres komunikatów sterujących | Komunikaty pętli prezentowane w oknie: zlecenia i zadania albo pełna wymiana Koordynator ↔ Wykonawca wraz z wynikami kontroli jakości | zlecenia i zadania | 3 | Menu `⋮` nagłówka Execution Loop Window |
| Próg automatycznego ponowienia zadania | Liczba ponowień zadania podejmowanych przez Koordynatora przed zatrzymaniem pętli i zapytaniem Użytkownika | 2 | 3 | Menu `⋮` nagłówka Execution Loop Window |
| Powiadomienie o zatrzymaniu pętli | Kategoria powiadomień wiązana ze zdarzeniami pętli wykonawczej (rozdz. 7.2) | aktywne | 2 | Sekcja „Powiadomienia” w selektorze sekcji `☰` |
| Diagnostyka kanałów komunikacji | Podgląd surowej wymiany komunikatów obu kanałów wraz z kopertami komunikatów | wyłączona | 4 | Tryb administracyjny; wyszukiwarka funkcji; polecenie języka naturalnego w Chat Window; skrót klawiszowy |

### 9.4 Makieta sekcji Okna komunikacji

```
Makieta 8 — Sekcja „Okna komunikacji”

 ═════════════════════════════════════════════════════════════════════════════════
  Pasek kontekstu:  [Danaco Console] [Ubuntu] [Worktree] [Fable 5] [Ultra]
 ═════════════════════════════════════════════════════════════════════════════════
  Chat Window      │ Execution Loop ▼│ Sekcje ☰            │ Okna komunikacji                      │
  Użytkownik ↔     │ Koordynator ↔   │                     │                                       │
  Wykonawca        │ Wykonawca       │   Konto             │ Chat Window                  [?] ⋮    │
                   │                 │   Uwierzytelnianie  │ Użytkownik ↔ Wykonawca                │
  [ polecenie… ]   │ [ pętla: 0 zad.]│   Wygląd i język    │ gęstość strumienia [ zwarta ▾ ]       │
                   │                 │   Urządzenia        │ potwierdzanie działań ( ● wł.)        │
                   │                 │   Powiadomienia     │ sterowanie ustawieniami ( ● wł.)      │
                   │                 │   Konta modeli      │                                       │
                   │                 │ ◈ Okna komunikacji  │                                       │
                   │                 │  Warstwy widoczności│ Execution Loop Window        [?] ⋮    │
                   │                 │                     │ Koordynator ↔ Wykonawca               │
                   │                 │                     │ otwarcie kolumny   ( ○ zwinięta )     │
                   │                 │                     │ zakres komunikatów [ zlecenia ▾ ]     │
                   │                 │                     │ próg ponowienia    [ 2 ]              │
 ═════════════════════════════════════════════════════════════════════════════════

Stan spoczynku: widoczne wyłącznie elementy warstwy 1 oraz zwinięte
wyzwalacze warstw 2–3 (`▼`, `⋮`, `☰`) i znaczniki paska kontekstu.
```

---

## 10. Warstwy widoczności — ustawienia według roli użytkownika

Sekcja przypisuje warstwy widoczności (rozdz. 2.5) rolom użytkownika platformy. Przypisanie rozstrzyga, które elementy interfejsu są dla danej roli widoczne bez interakcji, które pod wyzwalaczem, które w menu kontekstowym, a które dostępne wyłącznie poleceniem języka naturalnego, skrótem klawiszowym, wyszukiwarką funkcji albo w trybie administracyjnym. Sekcja należy do warstwy 4 — użytkownik podstawowy jej nie widzi.

### 10.1 Grupa ustawień „Warstwy widoczności według roli”

| Ustawienie | Opis | Wartość domyślna | Warstwa | Sposób wywołania |
|---|---|---|---|---|
| Rola użytkownika | Rola, dla której konfigurowane jest przypisanie warstw: Operator, użytkownik podstawowy, użytkownik zaawansowany, administrator | Operator | 4 | Tryb administracyjny; wyszukiwarka funkcji; polecenie języka naturalnego w Chat Window |
| Najwyższa warstwa widoczna dla roli | Warstwa, do której włącznie elementy są dostępne bez trybu administracyjnego | 3 dla Operatora, 2 dla użytkownika podstawowego, 4 dla administratora | 4 | Tryb administracyjny; skrót klawiszowy |
| Elementy warstwy 1 dla roli | Zestaw elementów widocznych bez interakcji: Chat Window, pasek kontekstu, panel zawartości sekcji otwartej | Chat Window, pasek kontekstu, panel zawartości | 4 | Tryb administracyjny |
| Wyzwalacze warstwy 2 dla roli | Zestaw wyzwalaczy widocznych w stanie spoczynku: `▼` Execution Loop Window, `☰` selektor sekcji, znaczniki kontekstowe | wszystkie wymienione | 4 | Tryb administracyjny |
| Zakres menu kontekstowych warstwy 3 dla roli | Akcje udostępniane w menu `⋮` wierszy tabel i nagłówków okien | pełny zakres akcji sekcji | 4 | Tryb administracyjny |
| Dostęp do funkcji warstwy 4 dla roli | Udostępnienie funkcji eksperckich: diagnostyka kanałów, podgląd danych dostępowych, konfiguracja warstw | administrator | 4 | Tryb administracyjny; konfiguracja roli |
| Kanały wywołania funkcji ukrytych | Kanały, którymi rola sięga po funkcje ukryte: polecenie języka naturalnego, skrót klawiszowy, wyszukiwarka funkcji | wszystkie trzy kanały | 4 | Tryb administracyjny; polecenie języka naturalnego w Chat Window |

Zmiana przypisania obowiązuje natychmiast na wszystkich urządzeniach Operatora (mechanizm zapisu — rozdz. 15.2, `[DO DECYZJI OPERATORA]`). Przypisanie nie usuwa funkcji z platformy — zmienia wyłącznie sposób jej wywołania; każda funkcja pozostaje osiągalna jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego w Chat Window (zasada jednego kliknięcia, rozdz. 2.5).

### 10.2 Makieta sekcji Warstwy widoczności

```
Makieta 9 — Sekcja „Warstwy widoczności” (tryb administracyjny)

 ═════════════════════════════════════════════════════════════════════════════════
  Pasek kontekstu:  [Danaco Console] [Ubuntu] [Worktree] [Fable 5] [Ultra]
 ═════════════════════════════════════════════════════════════════════════════════
  Chat Window      │ Execution Loop ▼│ Sekcje ☰            │ Warstwy widoczności          [?] ⋮    │
  Użytkownik ↔     │ Koordynator ↔   │                     │ warstwa 4 — tryb administracyjny      │
  Wykonawca        │ Wykonawca       │   Konto             │                                       │
                   │                 │   Uwierzytelnianie  │ Rola  [ Operator            ▾ ]       │
  [ polecenie… ]   │ [ pętla: 0 zad.]│   Wygląd i język    │                                       │
                   │                 │   Urządzenia        │ Najwyższa warstwa widoczna  [ 3 ▾ ]   │
                   │                 │   Powiadomienia     │                                       │
                   │                 │   Konta modeli      │ Warstwa 1  [●] Chat Window            │
                   │                 │   Okna komunikacji  │            [●] pasek kontekstu        │
                   │                 │ ◈ Warstwy widoczn.  │            [●] panel zawartości       │
                   │                 │                     │ Warstwa 2  [●] ▼ Execution Loop       │
                   │                 │                     │            [●] ☰ selektor sekcji      │
                   │                 │                     │ Warstwa 3  [●] menu ⋮ wierszy         │
                   │                 │                     │ Warstwa 4  [ ] diagnostyka kanałów    │
                   │                 │                     │                                       │
                   │                 │                     │ Kanały funkcji ukrytych:              │
                   │                 │                     │ [●] język naturalny [●] skrót         │
                   │                 │                     │ [●] szukaj                            │
 ═════════════════════════════════════════════════════════════════════════════════

Stan spoczynku: widoczne wyłącznie elementy warstwy 1 oraz zwinięte
wyzwalacze warstw 2–3 (`▼`, `⋮`, `☰`) i znaczniki paska kontekstu.
```


---

## 11. Always On Display

### 11.1 Zakres sekcji

Sekcja udostępnia ustawienia funkcji globalnej Always On Display — agenta towarzyszącego obecnego ponad wszystkimi środowiskami, modułami, projektami i sesjami. Pełną dokumentację funkcji, w tym reguły wyzwalania sugestii, katalog ich rodzajów oraz rozgraniczenie toru głosowego wobec modułu Assistant, zawiera opracowanie [Always On Display](../funkcje-globalne/always-on-display.md). Pełny zakres ustawień funkcji obejmuje okno Konfiguracji (Model konfiguracji, rozdz. 5.17); niniejsza sekcja udostępnia pozycje o zastosowaniu codziennym.

| Ustawienie | Opis | Warstwa konfiguracji | Wartość domyślna | Warstwa widoczności | Sposób wywołania |
|---|---|---|---|---|---|
| Always On Display | Obecność, tryb, reguły wyzwalania sugestii, wyciszanie i tor głosowy funkcji globalnej (Koncepcja platformy, rozdz. 12.2) | globalna, środowisko, sesja | funkcja czynna, tryb pełny | 2 | Sekcja „Always On Display” w selektorze sekcji `☰` |

### 11.2 Zawartość sekcji

| Pozycja | Zawartość | Wartość domyślna | Warstwa | Sposób wywołania |
|---|---|---|---|---|
| Obecność awatara | Przełącznik widoczności pływającego awatara; ukrycie nie przerywa działania funkcji w tle — dostęp pozostaje przez szybki wybór szyny nawigacji, skrót klawiszowy i polecenie języka naturalnego | awatar widoczny | 2 | Sekcja „Always On Display” w selektorze sekcji `☰` |
| Tryb obecności | Wybór trybu: pełny, cichy, ukryty | pełny | 2 | Sekcja „Always On Display” w selektorze sekcji `☰` |
| Reguły wyzwalania | Tabela klas zdarzeń wyzwalających sugestie (stan pętli wykonawczej, stan kolejki zadań, wynik kontroli jakości, zdarzenia modułów, harmonogram, kontekst pracy) wraz z progami i częstotliwością ujawniania | wszystkie klasy czynne, wartości progów domyślne | 4 | Tryb administracyjny, wyszukiwarka funkcji, skrót klawiszowy, polecenie języka naturalnego w Chat Window |
| Wyciszanie | Wybór czasu wyciszenia (15 minut, 1 godzina, do końca dnia) i jego zakresu (cała platforma, bieżący moduł, bieżąca karta sesji, wskazana klasa zdarzeń) oraz podgląd wyciszeń czynnych | brak wyciszenia | 3 | Menu akcji `⋮` w nagłówku sekcji |
| Tor głosowy | Przełącznik globalnego toru głosowego funkcji, globalna fraza wybudzająca, synteza mowy odpowiedzi | tor czynny, synteza czynna | 2 | Sekcja „Always On Display” w selektorze sekcji `☰` |

Tor głosowy tej sekcji dotyczy wyłącznie krótkich poleceń ponadkontekstowych obsługiwanych przez funkcję globalną. Tor sesyjny — rozmowa wieloturowa, dyktowanie, makra i profil głosu — należy do okna Voice Console modułu Assistant i konfigurowany jest w profilu asystenta.

### 11.3 Makieta sekcji Always On Display

```
Makieta 10 — Sekcja „Always On Display”

 ═════════════════════════════════════════════════════════════════════════════════
  Pasek kontekstu:  [Danaco Console] [Ubuntu] [Worktree] [Fable 5] [Ultra]
 ═════════════════════════════════════════════════════════════════════════════════
  Chat Window      │ Execution Loop ▼│ Sekcje ☰            │ Always On Display              ⋮      │
  Użytkownik ↔     │ Koordynator ↔   │                     │                                       │
  Wykonawca        │ Wykonawca       │   Konto             │ Obecność awatara            [?]       │
                   │                 │   Uwierzytelnianie  │                           ( ● wł.)    │
  [ polecenie… ]   │ [ pętla: 0 zad.]│   Wygląd i język    │                                       │
                   │                 │   Urządzenia        │ Tryb obecności              [?]       │
                   │                 │   Powiadomienia     │ [ pełny ][ cichy ][ ukryty ]          │
                   │                 │   Konta modeli      │                                       │
                   │                 │   Okna komunikacji  │ Tor głosowy                 [?]       │
                   │                 │  Warstwy widoczności│                           ( ● wł.)    │
                   │                 │ ◈ Always On Display │ Fraza wybudzająca                     │
                   │                 │                     │ [ Danaco                   ]          │
                   │                 │                     │                                       │
                   │                 │                     │ Wyciszanie i reguły wyzwalania        │
                   │                 │                     │   — menu akcji ⋮ nagłówka sekcji      │
 ═════════════════════════════════════════════════════════════════════════════════

Stan spoczynku: widoczne wyłącznie elementy warstwy 1 oraz zwinięte
wyzwalacze warstw 2–3 (`▼`, `⋮`, `☰`) i znaczniki paska kontekstu.
```

### 11.4 Zachowanie

Zmiana każdej pozycji zapisywana jest natychmiast i rozgłaszana zdarzeniem `config.changed` na wszystkie urządzenia Operatora. Ukrycie awatara i tryb cichy nie zatrzymują wykrywania zdarzeń — sugestie gromadzą się w liście oczekujących funkcji i pozostają dostępne po otwarciu jej powierzchni interakcji. Punkt decyzyjny pętli wykonawczej wstrzymujący proces ujawnia się mimo wyciszenia, w postaci plakietki bez dymka.

---

## 12. Makieta całościowa okna Ustawień

Poniższa makieta zbiera układ całego okna w jednym schemacie kolumnowym, w stanie spoczynku interfejsu: Chat Window i Execution Loop Window po lewej, selektor dziewięciu sekcji `☰`, panel zawartości z sekcją „Konto” otwieraną domyślnie przy wejściu do okna oraz panel pomocniczy jako rozszerzenie boczne.

```
Makieta 11 — Okno Ustawień, widok całościowy (stan spoczynku, sekcja domyślna: Konto)

 ═════════════════════════════════════════════════════════════════════════════════
  Pasek kontekstu:  [Danaco Console] [Ubuntu] [Worktree] [Fable 5] [Ultra]
 ═════════════════════════════════════════════════════════════════════════════════
  Chat Window      │ Execution Loop ▼│ Sekcje ☰            │ USTAWIENIA                        ✕   │
  Użytkownik ↔     │ Koordynator ↔   │                     │                                       │
  Wykonawca        │ Wykonawca       │ ◈ Konto             │ Konto Operatora                       │
                   │                 │   Uwierzytelnianie  │ ┌───────┐                             │
  [ polecenie… ]   │ [ pętla: 0 zad.]│   Wygląd i język    │ │ awatar│  Zmień awatar               │
                   │                 │   Urządzenia        │ └───────┘                             │
                   │                 │   Powiadomienia     │                                       │
                   │                 │   Konta modeli      │ Login                       [?] ⋮     │
                   │                 │   Okna komunikacji  │ [ dnaharnowicz             ]          │
                   │                 │  Warstwy widoczności│                                       │
                   │                 │   Always On Display │ Adres e-mail uwierzytelniający  [?]   │
                   │                 │                     │ [ dnaharnowicz@danacogroup ] ✓ potw.  │
                   │                 │                     │                                       │
                   │                 │                     │ Hasło                                 │
                   │                 │                     │ ••••••••     [ Zmień hasło ]          │
                   │                 │                     │                                       │
                   │                 │                     │ Konto utworzone: 12.03.2026           │
                   │                 │                     │                                       │
                   │                 │                     │         [ Zapisz zmiany ]             │
 ═════════════════════════════════════════════════════════════════════════════════

Stan spoczynku: widoczne wyłącznie elementy warstwy 1 oraz zwinięte
wyzwalacze warstw 2–3 (`▼`, `⋮`, `☰`) i znaczniki paska kontekstu.
```

Okno otwiera się jako pełnoekranowy widok analogiczny do Konfiguracji — nie jako okno nakładkowe `.dn-modal` — ponieważ liczba i objętość sekcji (w szczególności tabele Urządzeń i Kont modeli) przekracza zakres wygodny dla nakładki modalnej; zasada zgodna z klasyfikacją okna Konfiguracji punktów izolacji jako okna globalnego, nie okna nakładkowego (Specyfikacja okien operacyjnych, rozdz. 8).

---

## 13. Katalog elementów interfejsu

Poniższa tabela katalogowa obejmuje każdy element interfejsu występujący w oknie Ustawień, zgodnie ze standardem dokumentu: co to jest, do czego służy, forma i waga wizualna, stany, zachowanie po interakcji, miejsce występowania, warstwa widoczności (rozdz. 2.5) i sposób wywołania.

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Selektor sekcji (panel rozwinięty) | Pionowa lista dziewięciu pozycji, komponent `.dn-karta--klikalna` | Wybór sekcji ustawień do wyświetlenia w panelu zawartości | Panel boczny, szerokość stała, waga niska (tło neutralne) | domyślny · wybrany (podświetlenie `--dn-sygnal-tlo`) · wskazanie kursorem | Kliknięcie pozycji zamienia zawartość panelu prawego bez przeładowania okna | Całe okno (rozdz. 2.1) | 1 | Widoczny bez interakcji; selektor `☰` w stanie zwiniętym na wąskich szerokościach |
| Awatar Operatora | Obraz koła, komponent `.dn-awatar` | Identyfikacja wizualna Operatora w interfejsie i oknach wielomodelowych | Mała ikona okrągła (rozmiar `--sm`/`--lg`, ok. 40–48 px) | domyślny (inicjał) · z obrazem własnym · ładowanie (podczas wgrywania) | Kliknięcie otwiera wybór pliku obrazu; zapis natychmiastowy | Sekcja Konto (3.2) | 1 | Widoczny bez interakcji w sekcji Konto |
| Pole „Login” | Pole tekstowe jednowierszowe, `.dn-pole-kontrolka` w `.dn-pole` | Wprowadzenie i edycja nazwy logowania | Pole formularza, waga średnia, pełna szerokość kolumny | domyślny · fokus · wypełnione · ostrzeżenie (login zajęty) · zapisywanie | Edycja oznacza zmiany niezapisane przy przycisku „Zapisz zmiany” (przycisk pozostaje aktywny cały czas); zapis wymaga potwierdzenia hasłem bieżącym, gdy Operator włączył to ustawienie w Uwierzytelnianiu (rozdz. 4.1) | Sekcja Konto (3.2) | 1 | Widoczne bez interakcji w sekcji Konto |
| Pole „Adres e-mail uwierzytelniający” | Pole tekstowe jednowierszowe z plakietką stanu, `.dn-pole-kontrolka` + `.dn-plakietka` | Wprowadzenie i edycja adresu e-mail pełniącego funkcję weryfikacji, logowania i odzyskiwania | Pole formularza, waga średnia | domyślny · fokus · potwierdzony (`.dn-plakietka--sukces`) · niepotwierdzony (`.dn-plakietka--ostrzezenie`) · ostrzeżenie (format nieprawidłowy) | Zmiana adresu uruchamia ponowną weryfikację e-mail (link potwierdzający); zapis nie jest blokowany formatem — ostrzeżenie towarzyszy polu do czasu poprawy | Sekcja Konto (3.2) | 1 | Widoczne bez interakcji w sekcji Konto |
| Przycisk „Zmień hasło” | Przycisk drugorzędny, `.dn-btn--zarys` | Otwarcie formularza zmiany hasła | Mały przycisk tekstowy | domyślny · wskazanie kursorem · ładowanie (zapis innej zmiany w toku — przycisk pozostaje klikalny) | Otwiera okno nakładkowe `.dn-modal` z formularzem: nowe hasło i jego powtórzenie zawsze; pole hasła bieżącego widoczne, gdy ustawienie jest włączone w Uwierzytelnianiu (rozdz. 4.1) | Sekcja Konto (3.2, 3.4) | 2 | Działanie w sekcji Konto; polecenie języka naturalnego w Chat Window |
| Okno nakładkowe zmiany hasła | Okno nakładkowe, `.dn-modal` (token `--dn-nakladka` na `::backdrop`) | Przeprowadzenie zmiany hasła, z weryfikacją hasła bieżącego sterowaną ustawieniem (rozdz. 4.1) | Duży panel nakładkowy, wyśrodkowany | ładowanie (zapis) · ostrzeżenie (hasło bieżące niezgodne, gdy ustawienie jest włączone / nowe hasła niezgodne ze sobą / pole puste) · sukces (zamknięcie + dymek powiadomienia) | Przycisk zatwierdzenia jest zawsze aktywny; kliknięcie z niekompletnymi lub niezgodnymi danymi wyświetla ostrzeżenie przy właściwym polu, formularz pozostaje otwarty do poprawy | Sekcja Konto (3.4) | 3 | Rozwinięcie kontekstowe po użyciu działania „Zmień hasło” |
| Przycisk „Zapisz zmiany” | Przycisk główny akcji, `.dn-btn--sygnal` | Zapisanie zmian pól Konta | Duży przycisk, jedyny CTA widoku | domyślny (aktywny niezależnie od obecności zmian) · ładowanie · ostrzeżenie (pole wymagane puste lub reguła niespełniona — zapis pozostaje do ponowienia) · błąd (zapis nie powiódł się po stronie serwera) | Przycisk klikalny zawsze; zapis przez `account.update` (rozdz. 15.1) — brak zmian niezapisanych sprawia, że kliknięcie nic nie wysyła, a brak w polu wymaganym sygnalizowany jest ostrzeżeniem przy polu, nie blokadą przycisku | Sekcja Konto (3.2) | 1 | Widoczny bez interakcji w sekcji Konto |
| Przełącznik metody logowania | Przełącznik, `.dn-suwak` | Włączenie/wyłączenie metody logowania (hasło, e-mail, PIN, Windows Hello) | Mały suwak inline, tor w gradiencie sygnałowym gdy włączony | wyłączony · włączony · ładowanie (podczas zapisu) · ostrzeżenie (np. Windows Hello niedostępne na tym urządzeniu — przełącznik pozostaje klikalny, komunikat po kliknięciu) · błąd (zapis się nie powiódł, przełącznik wraca do poprzedniego stanu) | Przełączenie zapisuje stan natychmiast (`auth.method.add` przy włączeniu, `auth.method.remove` przy wyłączeniu); przy PIN i Windows Hello dodatkowo otwiera krok konfiguracji; przy próbie wyłączenia ostatniej aktywnej metody logowania — ostrzeżenie zamiast blokady (rozdz. 1.3) | Sekcja Uwierzytelnianie (4.1, 4.2) | 2 | Sekcja „Uwierzytelnianie” w selektorze sekcji `☰` |
| Przycisk „Ustaw PIN” | Przycisk drugorzędny, `.dn-btn--zarys` | Uruchomienie kroku ustanowienia kodu PIN | Mały przycisk tekstowy, widoczny wyłącznie po włączeniu przełącznika PIN | ukryty (przełącznik wyłączony) · domyślny · ładowanie | Otwiera krótki formularz dwukrotnego podania kodu numerycznego | Sekcja Uwierzytelnianie (4.2) | 3 | Rozwinięcie kontekstowe po włączeniu przełącznika PIN |
| Plakietka stanu wymogu logowania | Plakietka informacyjna, `.dn-plakietka--informacja` | Poinformowanie Operatora, czy okno logowania aktualnie wymusza uwierzytelnienie | Mały blok informacyjny z kreską lewą, tylko do odczytu | wymóg logowania włączony · wymóg logowania wyłączony | Brak — element wyłącznie informacyjny, bez interakcji | Sekcja Uwierzytelnianie (4.2, 4.3) | 1 | Widoczna bez interakcji w sekcji Uwierzytelnianie |
| Przełącznik motywu | Zakładki segmentowe, `.dn-zakladki--pigulki`, trzy warianty | Wybór wariantu jasnego, ciemnego albo zgodnego z systemem | Mały segmentowany przełącznik, ikony `slonce`/`ksiezyc` | domyślny (System zaznaczony) · Jasny zaznaczony · Ciemny zaznaczony | Wybór natychmiast przełącza motyw całej platformy z animacją przejścia | Sekcja Wygląd i język (5.1, 5.3) | 2 | Sekcja „Wygląd i język” w selektorze sekcji `☰`; skrót klawiszowy |
| Lista wyboru języka | Pole wyboru, `.dn-pole-kontrolka` | Wybór języka interfejsu | Pole formularza z rozwijaną listą, waga mała | domyślny (Polski) · rozwinięty · wybrany inny język | Wybór natychmiast przełącza język całego interfejsu | Sekcja Wygląd i język (5.2, 5.3) | 2 | Sekcja „Wygląd i język” w selektorze sekcji `☰` |
| Tabela urządzeń połączonych | Tabela danych, `.dn-tabela` | Prezentacja wszystkich urządzeń powiązanych z Kontem | Duża tabela pełnej szerokości panelu | ładowanie · wypełniona · pusta (teoretyczna — zawsze min. jedno urządzenie bieżące) | Wiersz „to urządzenie” bez akcji odłączenia; pozostałe wiersze z akcją „Odłącz” | Sekcja Urządzenia (6.1, 6.4) | 2 | Sekcja „Urządzenia” w selektorze sekcji `☰` |
| Przycisk „Odłącz” (wiersz urządzenia) | Przycisk ikonowy w wierszu tabeli, `.dn-btn-ikona` (36×36 px) | Unieważnienie tokenu dostępu wskazanego urządzenia | Mała ikona w kolumnie akcji | domyślny · wskazanie kursorem · ładowanie | Wykonuje `device.revoke` od razu; dymek powiadomienia „Odłączono: <nazwa urządzenia>” z dostępną przez kilka sekund akcją „Cofnij”; okno nakładkowe potwierdzenia przed odłączeniem jest ustawieniem konfiguracyjnym sekcji Uwierzytelnianie (rozdz. 4), domyślnie wyłączonym | Sekcja Urządzenia (6.1) | 3 | Menu akcji wiersza `⋮` w tabeli urządzeń |
| Wskaźnik stanu synchronizacji | Plakietka ze wskaźnikiem, `.dn-plakietka--sygnal` + `.dn-kropka` | Informowanie o stanie połączenia synchronizacji na żywo | Mała plakietka w nagłówku sekcji | aktywna (`--sukces`) · przerwana (`--ostrz`) | Brak interakcji — wyłącznie odczyt stanu | Sekcja Urządzenia (6.3, 6.4) | 1 | Widoczny bez interakcji w nagłówku sekcji Urządzenia |
| Przycisk „Sparuj urządzenie” | Przycisk główny akcji, `.dn-btn--sygnal` | Uruchomienie procedury parowania nowego urządzenia mobilnego | Duży przycisk, jedyny CTA sekcji | domyślny · ładowanie (generowanie kodu) · aktywny kod wyświetlony | Otwiera okno nakładkowe z kodem w postaci znakowej i kodem graficznym (protokół: [Mobile](../funkcje-globalne/mobile.md), rozdz. 9) | Sekcja Urządzenia (6.2, 6.4) | 2 | Sekcja „Urządzenia”; polecenie języka naturalnego w Chat Window |
| Przełącznik główny powiadomień | Przełącznik, `.dn-suwak` | Włączenie/wyłączenie wszystkich powiadomień jednocześnie | Mały suwak inline, umieszczony na szczycie sekcji | włączony · wyłączony | Wyłączenie wygasza (nie kasuje) ustawienia poszczególnych kategorii poniżej | Sekcja Powiadomienia (7.4, 7.5) | 2 | Sekcja „Powiadomienia” w selektorze sekcji `☰` |
| Tabela klas zdarzeń | Tabela z polami wyboru, `.dn-tabela` + `.dn-check` w komórkach | Wybór, które klasy zdarzeń i którym kanałem mają być dostarczane | Duża tabela, kolumny: klasa zdarzenia, centrum powiadomień, Mobile, e-mail | wypełniona · przygaszona wizualnie (gdy przełącznik główny nieaktywny), w pełni klikalna | Zaznaczenie i odznaczenie pola wyboru zapisuje natychmiast, bez zbiorczego zapisu — zmiany dokonane przy wyłączonym przełączniku głównym obowiązują po jego ponownym włączeniu | Sekcja Powiadomienia (7.2, 7.4) | 2 | Sekcja „Powiadomienia” w selektorze sekcji `☰` |
| Przełącznik obecności awatara Always On Display | Przełącznik, `.dn-suwak` | Włączenie i ukrycie pływającego awatara funkcji globalnej | Mały suwak inline | włączony · wyłączony · ładowanie (podczas zapisu) | Ukrycie awatara nie przerywa działania funkcji — dostęp pozostaje przez szybki wybór szyny nawigacji, skrót klawiszowy i polecenie języka naturalnego | Sekcja Always On Display (11.2, 11.3) | 2 | Sekcja „Always On Display” w selektorze sekcji `☰` |
| Przełącznik trybu obecności | Zakładki segmentowe, `.dn-zakladki--pigulki`, trzy warianty | Wybór trybu obecności funkcji: pełny, cichy, ukryty | Mały przełącznik segmentowy | pełny zaznaczony · cichy zaznaczony · ukryty zaznaczony | Wybór obowiązuje natychmiast na wszystkich urządzeniach Operatora | Sekcja Always On Display (11.2, 11.3) | 2 | Sekcja „Always On Display” w selektorze sekcji `☰` |
| Pole „Fraza wybudzająca” | Pole tekstowe jednowierszowe, `.dn-pole-kontrolka` | Ustanowienie globalnej frazy uruchamiającej tor głosowy funkcji | Pole formularza, waga średnia | domyślny · fokus · zapisywanie · ostrzeżenie (fraza zbyt krótka — zapis pozostaje do ponowienia) | Zapis natychmiastowy; fraza obowiązuje we wszystkich środowiskach i modułach | Sekcja Always On Display (11.2, 11.3) | 2 | Sekcja „Always On Display” w selektorze sekcji `☰` |
| Menu wyciszania i reguł wyzwalania `⋮` | Menu kebab w nagłówku sekcji | Zgrupowanie działań wyciszenia oraz wejścia do reguł wyzwalania sugestii | Mała ikona w nagłówku sekcji | zwinięte · rozwinięte · wyciszenie czynne | Wybór czasu i zakresu wyciszenia obowiązuje natychmiast; pozycja „Reguły wyzwalania” otwiera zestaw ustawień warstwy 4 | Sekcja Always On Display (11.2) | 3 | Wyzwalacz `⋮` w nagłówku sekcji |
| Tabela reguł wyzwalania sugestii | Tabela klas zdarzeń z przełącznikami i wartościami progów, `.dn-tabela` + `.dn-suwak` | Ustalenie, które klasy zdarzeń tworzą sugestie i przy jakich progach | Duża tabela pełnej szerokości panelu | ukryta (poza trybem administracyjnym) · wypełniona · wiersz w edycji | Zmiana klasy lub progu zapisywana natychmiast, bez zbiorczego zapisu | Sekcja Always On Display (11.2) | 4 | Tryb administracyjny, wyszukiwarka funkcji, skrót klawiszowy, polecenie języka naturalnego w Chat Window |
| Tabela kont modeli | Tabela danych, `.dn-tabela` | Prezentacja skonfigurowanych kont wiersza poleceń, w tym kont narzędzia Code CLI | Duża tabela pełnej szerokości panelu | ładowanie · wypełniona · pusta (`.dn-pusty-stan`) · wiersz w błędzie uwierzytelnienia | Menu akcji wiersza (ikona `wiecej`): „Ustaw jako aktywne”, „Testuj połączenie”, „Edytuj”, „Usuń” | Sekcja Konta modeli (8.2, 8.5, 8.6) | 2 | Sekcja „Konta modeli” w selektorze sekcji `☰` |
| Plakietka stanu konta modelu | Plakietka, `.dn-plakietka--sygnal` | Sygnalizacja, czy dane konto jest aktywne, skonfigurowane czy w błędzie | Mała pigułka w komórce tabeli | aktywne (`--sygnal`) · skonfigurowane (`--informacja`) · błąd uwierzytelnienia (`--blad`) | Brak interakcji bezpośredniej — stan wynikowy testu połączenia lub przełączenia | Sekcja Konta modeli (8.2, 8.5) | 2 | Widoczna w wierszu tabeli kont modeli |
| Przycisk „Dodaj konto modelu” | Przycisk główny akcji, `.dn-btn--sygnal` | Otwarcie formularza dodania nowego konta modelu | Duży przycisk, jedyny CTA sekcji | domyślny · ładowanie (po zatwierdzeniu formularza) | Otwiera okno nakładkowe opisane w rozdz. 8.3 | Sekcja Konta modeli (8.3, 8.5, 8.6) | 2 | Sekcja „Konta modeli”; polecenie języka naturalnego w Chat Window |
| Okno nakładkowe „Dodaj konto modelu” | Okno nakładkowe, `.dn-modal` (token `--dn-nakladka` na `::backdrop`) | Wprowadzenie danych nowego konta: nazwa, katalog profilu, dane dostępowe | Duży panel nakładkowy, wyśrodkowany, cztery pola | domyślny · ostrzeżenie pola (katalog zajęty, dane dostępowe puste) · test w toku · test powiódł się/nie powiódł się · zapisywanie | Przycisk „Zapisz konto” jest zawsze aktywny; zapis tworzy encję konta modelu z odwołaniem do danych dostępowych poza bazą, ostrzeżenie przy polu nie blokuje zatwierdzenia; okno nakładkowe zamyka się i pojawia dymek powiadomienia | Sekcja Konta modeli (8.3) | 3 | Rozwinięcie kontekstowe po użyciu działania „Dodaj konto modelu” |
| Przycisk „Testuj połączenie” | Przycisk drugorzędny, `.dn-btn--zarys` | Sprawdzenie poprawności danych dostępowych bez zamykania formularza | Mały przycisk tekstowy wewnątrz okna nakładkowego | domyślny · ładowanie · sukces (`.dn-plakietka--sukces`) · błąd (`.dn-plakietka--blad`) | Wywołuje próbne połączenie kanałem CLI; wynik wyświetla się obok przycisku | Okno nakładkowe dodania/edycji konta modelu (8.3) | 3 | Rozwinięcie kontekstowe wewnątrz okna nakładkowego konta modelu |
| Chat Window | Główne okno komunikacji Użytkownik ↔ Wykonawca, złożenie `.dn-karta--klikalna` `.dn-awatar` `.dn-pole-kontrolka` `.dn-btn-ikona` oraz bloków treści technicznej w kroju maszynowym `--dn-ff-mono` (kolumna) | Przyjmowanie poleceń w języku naturalnym, prezentacja strumienia odpowiedzi i wyników, zatwierdzanie i przerywanie działań, wyjaśnianie wyniku i kontekstu; podstawowy mechanizm sterowania wszystkimi procesami platformy | Lewa kolumna, stała, pełna wysokość obszaru roboczego | gotowe · oczekiwanie na odpowiedź Wykonawcy · strumieniowanie odpowiedzi · oczekiwanie na zatwierdzenie · błąd kanału | Polecenie Użytkownika uruchamia zadanie; zatwierdzenie i przerwanie sterują jego przebiegiem; wynik pozostaje w strumieniu rozmowy | Każdy kontekst pracy, w tym otwarte okno Ustawień (rozdz. 2.1, 9.1) | 1 | Widoczne bez interakcji |
| Execution Loop Window | Okno pętli wykonawczej Koordynator ↔ Wykonawca, złożenie `.dn-karta--klikalna` `.dn-tabela` `.dn-plakietka` `.dn-btn` (kolumna) | Prezentacja zlecenia i jego dekompozycji na zadania, kolejki i stanu zadań, wymiany komunikatów sterujących, wyników kontroli jakości i decyzji o ponowieniu oraz sterowanie przebiegiem pętli | Kolumna sąsiadująca z Chat Window, otwierana, pełna wysokość obszaru roboczego | zwinięte · pętla bezczynna · pętla w toku · wstrzymana · ponowienie zadania · przerwana · błąd | Wstrzymanie, wznowienie, przerwanie i korekta zlecenia zmieniają przebieg pętli; każde zdarzenie pętli trafia do strumienia Chat Window | Każdy kontekst pracy, w tym otwarte okno Ustawień (rozdz. 2.1, 9.2) | 2 | Wyzwalacz `▼` „Execution Loop” w kolumnie sąsiadującej; polecenie języka naturalnego w Chat Window |
| Pasek kontekstu | Zestaw znaczników kontekstowych, `.dn-plakietka` `.dn-plakietka--sygnal` | Prezentacja bieżącego środowiska, repozytorium, projektu, modelu i wykonawcy, na przykład `[Danaco Console] [Ubuntu] [Worktree] [Fable 5] [Ultra]` | Wąski pas znaczników w kolumnie obszaru roboczego, waga niska | gotowe · znacznik aktywny · selektor rozwinięty | Kliknięcie znacznika otwiera odpowiedni selektor; wybór zwija selektor samoczynnie | Każdy kontekst pracy (rozdz. 2.5) | 1 | Widoczny bez interakcji; kliknięcie znacznika otwiera selektor warstwy 2 |
| Selektor sekcji `☰` | Menu progresywne dziewięciu sekcji Ustawień | Wybór sekcji prezentowanej w panelu zawartości | Kolumna wąska, waga niska; w stanie spoczynku zwinięty do wyzwalacza `☰` | zwinięty · rozwinięty · pozycja wybrana | Wybór pozycji zamienia zawartość panelu zawartości i zwija menu | Całe okno (rozdz. 2.1, 2.2) | 2 | Wyzwalacz `☰`; skrót klawiszowy; wyszukiwarka funkcji |
| Menu akcji wiersza `⋮` | Menu kebab w wierszu tabeli | Zgrupowanie akcji wiersza tabeli urządzeń i tabeli kont modeli | Mała ikona w kolumnie akcji | zwinięte · rozwinięte · pozycja w stanie ładowania | Wybór pozycji wykonuje akcję i zwija menu | Sekcje Urządzenia i Konta modeli (rozdz. 6.1, 8.2) | 3 | Wyzwalacz `⋮`; menu kontekstowe wiersza |
| Wyszukiwarka funkcji | Pole wyszukiwania funkcji platformy | Odnalezienie i uruchomienie dowolnej funkcji okna, w tym funkcji warstwy 4 | Panel nakładkowy wyśrodkowany, wywoływany skrótem | ukryta · otwarta · wyniki wyszukiwania · brak wyników | Wybór wyniku uruchamia funkcję jednym kliknięciem | Całe okno (rozdz. 2.5) | 4 | Skrót klawiszowy; polecenie języka naturalnego w Chat Window |
| Tryb administracyjny | Stan okna udostępniający funkcje eksperckie i diagnostyczne | Konfiguracja warstw widoczności według roli, podgląd danych dostępowych, narzędzia diagnostyczne kanałów | Bez własnej powierzchni w stanie spoczynku; po włączeniu rozszerza panel zawartości | wyłączony · włączony · operacja diagnostyczna w toku | Włączenie odsłania sekcję „Warstwy widoczności” i pozycje warstwy 4 | Sekcje Konta modeli i Warstwy widoczności (rozdz. 8, 10) | 4 | Konfiguracja roli; polecenie języka naturalnego w Chat Window; skrót klawiszowy |
| Objaśnienie kontekstowe `[?]` | Dymek podpowiedzi, `.dn-tooltip` | Wyjaśnienie działania danego ustawienia i jego wpływu na aplikację | Mała ikona/znak przy etykiecie ustawienia | ukryty · widoczny (wskazanie kursorem lub fokus) | Wskazanie kursorem lub fokus klawiaturą wyświetla dymek z treścią; brak zmiany stanu ustawienia | Każda pozycja ustawienia w każdej sekcji (rozdz. 2.3) | 2 | Wskazanie kursorem lub fokus klawiatury na znaku `[?]` przy etykiecie ustawienia |

---

## 14. Stany — zestawienie zbiorcze

Rozdział 13 dokumentuje stany każdego elementu z osobna. Poniższe zestawienie porządkuje stany na poziomie okna i poszczególnych sekcji — sytuacje złożone, w których stan wynika z więcej niż jednego elementu jednocześnie.

### 14.1 Stany okna jako całości

| Stan okna | Wyzwalacz | Prezentacja |
|---|---|---|
| Ładowanie początkowe | Otwarcie okna Ustawień | Szkielet paneli z `.dn-spinner` w miejscu treści sekcji domyślnej |
| Gotowe | Dane Konta, Uwierzytelniania, Urządzeń i Kont modeli wczytane | Pełna zawartość sekcji domyślnej (Konto) |
| Błąd połączenia | Utrata kanału WebSocket podczas pracy w oknie | `.dn-toast--blad` na szczycie panelu zawartości: „Brak połączenia z serwerem — zmiany nie będą zapisywane do czasu przywrócenia łączności”; pola pozostają edytowalne lokalnie, zapis wstrzymany |
| Zapisywanie w toku | Dowolna zmiana wymagająca potwierdzenia serwera | `.dn-spinner` przy elemencie zmienianym; pozostałe elementy okna pozostają aktywne |
| Sterowanie poleceniem języka naturalnego | Polecenie Użytkownika skierowane do Wykonawcy w Chat Window zmienia pozycję Ustawień | Element zmieniany przechodzi w stan „ładowanie”, wynik potwierdzony w strumieniu Chat Window; Chat Window pozostaje w lewej kolumnie, panel zawartości pozostaje aktywny |
| Pętla wykonawcza w toku | Koordynator prowadzi pętlę zleconą z Chat Window | Execution Loop Window otwarte jako kolumna sąsiadująca, wskaźniki przebiegu pętli aktywne; praca w oknie Ustawień pozostaje możliwa bez przerywania pętli |

### 14.2 Stany sekcji z listami i tabelami

| Sekcja | Stan pustej listy | Stan błędu wczytania | Stan wypełniony |
|---|---|---|---|
| Urządzenia (6) | Nie występuje — co najmniej urządzenie bieżące jest zawsze obecne | `.dn-toast--blad`: „Nie udało się pobrać wykazu urządzeń” + przycisk „Ponów” | Tabela `.dn-tabela` z co najmniej jednym wierszem |
| Konta modeli (8) | `.dn-pusty-stan`: „Brak skonfigurowanych kont modeli” + CTA „Dodaj konto modelu” | `.dn-toast--blad`: „Nie udało się pobrać kont modeli” + przycisk „Ponów” | Tabela `.dn-tabela` z co najmniej jednym wierszem |
| Powiadomienia (7) | Nie dotyczy — klasy zdarzeń są stałym katalogiem platformy, nie listą dynamiczną | — | Tabela klas zdarzeń zawsze w pełni wypełniona |

### 14.3 Stany formularzy i okien nakładkowych

| Okno nakładkowe / formularz | Stan początkowy | Stan ostrzeżenia walidacji | Stan sukcesu |
|---|---|---|---|
| Zmiana hasła (3.4) | Pola nowego hasła puste (pole hasła bieżącego widoczne tylko, gdy ustawienie jest włączone — rozdz. 4.1); przycisk zatwierdzenia aktywny od otwarcia okna nakładkowego | „Hasło bieżące nieprawidłowe” / „Nowe hasła nie są identyczne” / „Uzupełnij nowe hasło” pod właściwym polem, `.dn-pole-blad` — ostrzeżenie, zapis można ponowić od razu | Zamknięcie okna nakładkowego + dymek powiadomienia „Hasło zmienione” |
| Dodanie konta modelu (8.3) | Cztery pola puste; przycisk zatwierdzenia aktywny od otwarcia okna nakładkowego | „Katalog profilu już wykorzystany przez inne konto” / „Dane dostępowe odrzucone przez dostawcę” / „Uzupełnij pola wymagane” pod właściwym polem, `.dn-pole-blad` — ostrzeżenie, zapis można ponowić od razu | Zamknięcie okna nakładkowego + dymek powiadomienia „Konto modelu dodane” + nowy wiersz w tabeli |
| Parowanie urządzenia (6.2) | Kod w postaci znakowej i kod graficzny wyświetlone, oczekiwanie | „Czas parowania wygasł” — kod wygasły, działanie „Wydaj nowy kod” | Zamknięcie okna nakładkowego + dymek powiadomienia „Sparowano: <nazwa urządzenia>” + nowy wiersz w tabeli urządzeń |

### 14.4 Stany przełączników (wzorzec wspólny)

Wzorzec obowiązuje jednolicie dla wszystkich przełączników `.dn-suwak` sekcji Uwierzytelnianie i Powiadomienia, zgodnie z wymogiem dokumentowania pięciu stanów przy każdym elemencie interaktywnym.

| Stan | Wygląd | Znaczenie |
|---|---|---|
| Domyślny/wyłączony | Tor jasny, wskaźnik po lewej | Ustawienie nieaktywne |
| Aktywny/włączony | Tor w gradiencie sygnałowym, wskaźnik po prawej | Ustawienie aktywne |
| Ostrzeżenie środowiskowe | Tor z obrysem `--ostrz`, `.dn-tooltip` z komunikatem | Warunek techniczny niespełniony (np. Windows Hello na urządzeniu bez czytnika biometrycznego) — przełącznik pozostaje klikalny; kliknięcie wyświetla komunikat wyjaśniający zamiast zmieniać stan |
| Ładowanie | Wskaźnik z `.dn-spinner` w miejscu przełącznika | Zapis zmiany w toku po stronie serwera |
| Błąd | Tor z obrysem `--blad`, `.dn-tooltip` z komunikatem przy wskazaniu kursorem | Zapis się nie powiódł — przełącznik wraca do stanu sprzed próby zmiany |

---

## 15. Komunikacja — polecenia i zdarzenia

Zmiany dokonane w oknie Ustawień przekazywane są tym samym kanałem WebSocket, który obsługuje pozostałe polecenia i zdarzenia platformy (Architektura, rozdz. 11), zgodnie z konwencją nazewniczą `obszar.czynność` przyjętą w całej platformie (Bezpieczeństwo i uwierzytelnianie, rozdz. 12; Model konfiguracji, rozdz. 8).

### 15.1 Polecenia i zdarzenia przejęte wprost

| Zakres | Polecenie (klient → serwer) | Zdarzenie (serwer → klient) | Źródło |
|---|---|---|---|
| Uwierzytelnianie — metody dodatkowe | `auth.method.add`, `auth.method.remove`, `auth.method.list` | `auth.changed` | Bezpieczeństwo i uwierzytelnianie, rozdz. 12 |
| Urządzenia | `device.list`, `device.revoke` | `device.changed` | Bezpieczeństwo i uwierzytelnianie, rozdz. 12 |
| Wygląd, język, powiadomienia (zakres Aplikacja) | `config.get`, `config.set` | `config.changed` | Model konfiguracji, rozdz. 8 |

### 15.2 Polecenia i zdarzenia zakresu Konta Operatora, kont modeli, okien komunikacji i warstw widoczności

Kontrakt komunikacji obejmuje polecenia właściwe Kontu Operatora edytowanemu z Ustawień oraz obszar `account` kontraktu, wspólny dla kont modeli i kont narzędzia wiersza poleceń (rozdz. 8). Ustawienia obu kanałów komunikacji operacyjnej (rozdz. 9.3) oraz przypisanie warstw widoczności do roli (rozdz. 10.1) nie mają dziś odrębnej komendy kontraktu — oznaczone `[DO DECYZJI OPERATORA]` w tabeli niżej.

| Polecenie / zdarzenie | Kierunek | Opis |
|---|---|---|
| `auth.profile.set` | klient → serwer | Zmiana loginu i/lub adresu e-mail uwierzytelniającego z poziomu Ustawień (rozdz. 3.2); odrębne od `auth.register` |
| `auth.password.reset` | klient → serwer | Zmiana hasła bramki ze znanym hasłem bieżącym — Operator już zalogowany (rozdz. 3.4); odrębne od `auth.reset`, właściwego wyłącznie odzyskiwaniu konta bez znanego hasła |
| `account.list` | klient → serwer | Pobranie wykazu skonfigurowanych kont modeli i kont Code CLI (rozdz. 8.2) |
| `account.add` | klient → serwer | Dodanie nowego konta modelu albo konta Code CLI wraz z danymi dostępowymi (rozdz. 8.3) |
| `account.test` | klient → serwer | Próbne sprawdzenie konta bez zmiany jego zapisu — dla konta Code CLI: obecność katalogu konfiguracji i programu wiersza poleceń w pakiecie serwera (rozdz. 8.3) |
| `account.default.set` | klient → serwer | Wskazanie konta domyślnego swojego rodzaju — dokładnie jedno konto domyślne na rodzaj (rozdz. 8.4) |
| `account.remove` | klient → serwer | Usunięcie konta; kanały powołujące się na nie tracą powiązanie, nie znikają |
| `account.changed` | serwer → klient | Rozgłoszenie zmiany konta modelu albo konta Code CLI do wszystkich urządzeń (rozdz. 8.4); poświadczeń zdarzenie nie niesie |
| **[DO DECYZJI OPERATORA]** — nazwa komendy nieustalona | — | Ustawienia Chat Window i Execution Loop Window: zapamiętane położenie paska rozdzielającego kolumny, gęstość strumienia, zakres wyświetlanych komunikatów sterujących, progi potwierdzeń (rozdz. 9.3). Obszar `window` kontraktu nie niesie komendy dedykowanej; rozstrzygnięcia wymaga, czy ustawienia zapisują się przez `config.session.set` (obszar `config`) jako pola `SessionConfig`, czy zyskują własny obszar |
| **[DO DECYZJI OPERATORA]** — nazwa komendy nieustalona | — | Przypisanie warstw widoczności do roli użytkownika (rozdz. 10.1). Kontrakt nie niesie komendy dedykowanej; mechanizm pokrewny co do formy, lecz o innym przedmiocie, to `panel.sections.get` i `panel.sections.set` (obszar `panel`) — zapisują układ sekcji panelu jednego okna, nie przypisanie warstwy 1–4 do roli |
| `config.set` | klient → serwer | Zapis ustawienia zakresu „Always On Display”: obecność awatara, tryb obecności, reguły wyzwalania, wyciszanie, tor głosowy (rozdz. 11.2; Model konfiguracji, rozdz. 5.17) |
| `config.changed` | serwer → klient | Rozgłoszenie zmiany ustawień zakresu „Always On Display” do wszystkich urządzeń Operatora (rozdz. 11.4) |
| `aod.status.get` | klient → serwer | Odczyt stanu aktywacji i widoczności funkcji globalnej po zmianie ustawienia ([Always On Display](../funkcje-globalne/always-on-display.md), rozdz. 9.5) |

### 15.3 Schemat przepływu — zapis pojedynczego ustawienia

```
Diagram 4 — Wzorzec zapisu ustawienia w oknie Ustawień (dowolna sekcja)

Operator zmienia wartość elementu (przełącznik, pole, wybór z listy)
        │
        ▼
Element przechodzi w stan „ładowanie” (rozdz. 14.4)
        │
        ▼
Klient wysyła polecenie właściwe zakresowi (rozdz. 15.1–15.2) kanałem WebSocket
        │
        ▼
Serwer zapisuje zmianę jako jedyne źródło prawdy (Architektura, rozdz. 15)
        │
        ├── sukces ──► zdarzenie *.changed rozgłoszone do wszystkich urządzeń Operatora
        │               element wraca do stanu „gotowe” z nową wartością
        │
        └── błąd ──► element wraca do wartości poprzedniej, `.dn-tooltip` z komunikatem błędu
```

---

## 16. Zgodność z zasadami nadrzędnymi platformy

| Zasada nadrzędna (Koncepcja platformy, rozdz. 14) | Realizacja w oknie Ustawień |
|---|---|
| 1. Pełna kompozycyjność | Wszystkie dziewięć sekcji dostępnych z tego samego selektora sekcji, każda jednym kliknięciem (zasada jednego kliknięcia); Operator może skonfigurować dowolny podzbiór metod uwierzytelniania, kanałów powiadomień i kont modeli w dowolnej kombinacji |
| 2. Pełna konfigurowalność (zasada centralna) | Żadne ustawienie sekcji 4–10 nie jest wymuszone; każde ma wartość domyślną wystarczającą do pracy bez interwencji Operatora (rozdz. 1.3) |
| 3. Jawność i konfigurowalność zależności | Relacja między kontem modelu (Ustawienia) a przypisaniem kanału do sesji lub roli (Konfiguracja) jest jawna i opisana wprost (rozdz. 8.1, Schemat 4), nie ukryta w mechanizmie wewnętrznym; przełączenie konta aktywnego jest zawsze widoczną akcją Operatora (rozdz. 8.4) |
| 4. Orkiestracja | Katalog kont modeli (rozdz. 8) udostępnia przypisywanie odrębnych kont modeli poszczególnym rolom środowiska MultitaskingAI; Execution Loop Window (rozdz. 9.2) prezentuje przebieg orkiestracji prowadzonej przez Koordynatora (Koncepcja platformy, rozdz. 14, zasada 4) |
| 5. Stopniowe ujawnianie funkcjonalności | Każdy element okna należy do jednej z czterech warstw widoczności; w stanie spoczynku widoczne są wyłącznie elementy warstwy 1 oraz zwinięte wyzwalacze warstw 2–3, a funkcje eksperckie warstwy 4 są osiągalne poleceniem języka naturalnego w Chat Window, skrótem klawiszowym, wyszukiwarką funkcji lub w trybie administracyjnym (rozdz. 2.5, rozdz. 10) |

---

## 17. Słowniczek pojęć

| Pojęcie | Znaczenie |
|---|---|
| **Ustawienia** | Okno poziomu aplikacji, odrębne od Konfiguracji, gromadzące ustawienia dotyczące Operatora i jego relacji z platformą jako całością: konto, uwierzytelnianie, wygląd, urządzenia, powiadomienia, konta modeli (rozdz. 1) |
| **Konfiguracja** | Okno odrębne od Ustawień, gromadzące dwanaście pozostałych zakresów Modelu konfiguracji — zachowanie platformy podczas pracy (rozdz. 1.1) |
| **Konto Operatora** | Jedyna encja Konto właściwa Operatorowi platformy; zarządzana w sekcji Konto niniejszego okna (rozdz. 3, Bezpieczeństwo i uwierzytelnianie, rozdz. 11.1) |
| **Metoda dodatkowa uwierzytelniania** | Jedna z trzech metod konfigurowalnych obok hasła: PIN, Windows Hello, adres e-mail jako samodzielna metoda logowania (rozdz. 4, Bezpieczeństwo i uwierzytelnianie, rozdz. 6) |
| **Konto modelu** | Nazwany, zapisany zestaw danych dostępowych kanału CLI (w szczególności konta narzędzia Code CLI) wraz z katalogiem profilu, zarządzany w sekcji Konta modeli (rozdz. 8) |
| **Katalog profilu** | Lokalny katalog przechowujący dane logowania i pliki konfiguracyjne jednego konta modelu, odrębny dla każdego skonfigurowanego konta (rozdz. 8.2) |
| **Konto aktywne** | Konto modelu stosowane domyślnie dla nowo tworzonych kanałów CLI, dopóki sesja lub rola nie wskaże innego wprost (rozdz. 8.2, 8.4) |
| **Okno globalne** | Okno operacyjne dostępne niezależnie od aktywnego środowiska i modułu; Ustawienia i Konfiguracja są oknami tej kategorii (rozdz. 1.2, Specyfikacja okien operacyjnych, rozdz. 8.1) |
| **Chat Window** | Główne okno komunikacji Użytkownik ↔ Wykonawca; centralny punkt pracy i podstawowy mechanizm sterowania wszystkimi procesami platformy; lewa kolumna obszaru roboczego, stała, pełna wysokość (rozdz. 9.1) |
| **Execution Loop Window** | Okno pętli wykonawczej prezentujące komunikację Koordynator ↔ Wykonawca; drugi kanał komunikacji operacyjnej, otwierany jako kolumna sąsiadująca z Chat Window (rozdz. 9.2) |
| **Koordynator** | Komponent orkiestrujący platformy: dekomponuje zlecenie, przydziela i nadzoruje zadania (rozdz. 9.2) |
| **Wykonawca** | AI, agent lub system wykonawczy realizujący zadania zlecone przez Użytkownika i nadzorowane przez Koordynatora (rozdz. 9) |
| **Warstwa widoczności** | Jedna z czterech warstw, do której należy każdy element interfejsu: zawsze widoczna, widoczna na żądanie, rozwinięcie kontekstowe, funkcja ekspercka (rozdz. 2.5, rozdz. 10) |
| **Objaśnienie kontekstowe `[?]`** | Element towarzyszący każdemu ustawieniu, opisujący jego działanie i wpływ na aplikację; nośnik: `.dn-tooltip` (rozdz. 2.3) |

---

## Załącznik A. Scenariusze użycia

| Scenariusz | Kroki | Rozdziały właściwe |
|---|---|---|
| Operator zmienia hasło po podejrzeniu jego ujawnienia | 1. Otwiera Ustawienia → Konto. 2. Klika „Zmień hasło”. 3. Podaje nowe hasło dwukrotnie (oraz hasło bieżące, jeśli wymóg jest włączony w Uwierzytelnianiu, rozdz. 4.1). 4. Po zapisie pozostałe urządzenia otrzymują zdarzenie `auth.changed` i mogą wymagać ponownego zalogowania przy najbliższym unieważnieniu tokenu | 3.4, 14.1 |
| Operator włącza logowanie kodem PIN na komputerze biurowym | 1. Otwiera Ustawienia → Uwierzytelnianie. 2. Włącza przełącznik PIN. 3. W formularzu „Ustaw PIN” podaje czterocyfrowy kod dwukrotnie. 4. Od tej chwili logowanie na tym urządzeniu oferuje PIN obok hasła | 4.1, 4.2 |
| Operator paruje nowy telefon służbowy | 1. Na komputerze już uwierzytelnionym otwiera Ustawienia → Urządzenia. 2. Klika „Sparuj urządzenie”. 3. Skanuje wyświetlony kod aplikacją mobilną na telefonie. 4. Nowy wiersz pojawia się w tabeli urządzeń, funkcja Mobile staje się dostępna na telefonie | 6.2, 6.4 |
| Operator traci telefon i odłącza go zdalnie | 1. Z dowolnego innego uwierzytelnionego urządzenia otwiera Ustawienia → Urządzenia. 2. Odnajduje wiersz zgubionego telefonu. 3. Klika „Odłącz” — token zostaje unieważniony od razu (potwierdzenie w oknie nakładkowym jest ustawieniem konfiguracyjnym sekcji Uwierzytelnianie). 4. Telefon przy najbliższej próbie połączenia wymaga ponownego logowania; odłączenie da się cofnąć z poziomu dymka powiadomienia przez kilka sekund po wykonaniu | 6.1, 13.2 |
| Operator wyłącza powiadomienia e-mail o zmianach konfiguracji, zachowując powiadomienia w aplikacji | 1. Otwiera Ustawienia → Powiadomienia. 2. W wierszu „Zmiany konfiguracji” odznacza kolumnę „e-mail”, pozostawiając kolumnę „w aplikacji” bez zmian | 7.2, 7.4 |
| Operator dodaje drugie konto narzędzia Code CLI dla zespołu prawnego i przełącza się między kontami | 1. Otwiera Ustawienia → Konta modeli. 2. Klika „Dodaj konto modelu”, podaje nazwę „Zespół prawny”, katalog profilu i dane dostępowe, testuje połączenie, zapisuje. 3. W dowolnym momencie klika „Ustaw jako aktywne” przy koncie „Zespół prawny”. 4. Nowe sesje i role bez jawnie przypisanego konta korzystają odtąd z tego konta; już uruchomione sesje zachowują konto poprzednie | 8.3, 8.4 |
| Operator zmienia ustawienie poleceniem języka naturalnego | 1. W Chat Window w lewej kolumnie wydaje Wykonawcy polecenie „włącz logowanie kodem PIN”. 2. Wykonawca wykonuje zmianę i potwierdza ją w strumieniu rozmowy. 3. Przełącznik PIN w sekcji Uwierzytelnianie przyjmuje stan włączony na wszystkich urządzeniach Operatora | 2.5, 9.1, 4.1 |
| Operator nadzoruje pętlę wykonawczą w trakcie pracy w Ustawieniach | 1. Zleca zadanie w Chat Window. 2. Otwiera Execution Loop Window wyzwalaczem `▼` w kolumnie sąsiadującej. 3. Śledzi dekompozycję zlecenia na zadania i wyniki kontroli jakości, wstrzymuje i wznawia pętlę. 4. Równolegle zmienia ustawienia w panelu zawartości — obie kolumny komunikacji pozostają widoczne | 9.2, 9.3 |
| Administrator ogranicza widoczność funkcji eksperckich użytkownikowi podstawowemu | 1. Włącza tryb administracyjny. 2. Otwiera sekcję „Warstwy widoczności”. 3. Dla roli „użytkownik podstawowy” ustawia najwyższą warstwę widoczną na 2. 4. Funkcje warstw 3–4 pozostają dla tej roli osiągalne poleceniem języka naturalnego w Chat Window i wyszukiwarką funkcji | 2.5, 10.1 |
| Operator pracuje z urządzenia w innym kraju i przełącza motyw na jasny mimo ustawienia systemowego „ciemny” | 1. Otwiera Ustawienia → Wygląd i język. 2. W przełączniku motywu wybiera „Jasny” zamiast „System”. 3. Ustawienie zapisuje się w warstwie globalnej i obowiązuje na wszystkich urządzeniach Operatora, nadpisując odczyt z systemu | 5.1, 5.3, 5.4 |

---

## Załącznik B. Szablony konfiguracji

Szablony redakcyjne poniżej porządkują pola nowych elementów wprowadzonych niniejszym dokumentem, w konwencji przyjętej przez pozostałe opracowania platformy (np. Bezpieczeństwo i uwierzytelnianie, rozdz. 6.4; Strona główna i nawigacja, rozdz. 6.3).

```
# Szablon — pozycja Konta Operatora (rozdz. 3)
konto_operatora:
  login:                          <nazwa logowania>
  adres_email_uwierzytelniajacy:  <adres e-mail>
  stan_weryfikacji_email:         <potwierdzony | niepotwierdzony>
  haslo:                          # nigdy nie odczytywane wprost — wyłącznie zmiana
    ostatnia_zmiana:              <znacznik czasu>
  awatar:                         <odwołanie do pliku obrazu | brak — inicjał>
  data_utworzenia:                <znacznik czasu rejestracji>
```

```
# Szablon — pozycja metody uwierzytelniania (rozdz. 4; rozszerza
# szablon "Metody uwierzytelniania" z Bezpieczeństwa i uwierzytelniania, rozdz. 6.4)
metoda_uwierzytelniania:
  nazwa:            <hasło | e-mail uwierzytelniający | PIN | Windows Hello>
  stan:              <aktywna | nieaktywna>
  domyslna:          <tak — dotyczy wyłącznie hasła | nie>
  wyłączalna:        tak            # każda metoda konfigurowalna, hasło włącznie
  warstwa_widocznosci: 2
  sposob_wywolania:   selektor sekcji ☰ → „Uwierzytelnianie”
  objasnienie [?]:   opis działania metody i jej wpływu na logowanie
```

```
# Szablon — konto modelu (rozdz. 8)
konto_modelu:
  nazwa:                <nazwa własna nadana przez Operatora>
  typ_kanalu:            CLI                       # zgodnie z Integracją modeli, rozdz. 3
  narzedzie:             <nazwa narzędzia Code CLI>
  katalog_profilu:       <ścieżka lokalna, odrębna dla każdego konta>
  dane_dostepowe:        → magazyn danych dostępowych poza bazą (rozdz. 8.1)
  stan:                  <aktywne | skonfigurowane | błąd uwierzytelnienia>
  ostatnio_uzyte:        <znacznik czasu ostatniego wywołania>
  aktywne_domyslnie:     <tak | nie>                # tylko jedno konto na raz
```

```
# Szablon — ustawienie okna komunikacji operacyjnej (rozdz. 9)
okno_komunikacji:
  kanal:                <Chat Window: Użytkownik ↔ Wykonawca |
                         Execution Loop Window: Koordynator ↔ Wykonawca>
  kolumna:               <lewa, stała, pełna wysokość | kolumna sąsiadująca, otwierana>
  szerokosc_kolumny:     <wartość systemu wizualnego | wartość Operatora>
  zakres_komunikatow:    <zlecenia i zadania | pełna wymiana>
  prog_ponowienia:       <liczba ponowień zadania>
  warstwa_widocznosci:   <1 | 2>
  sposob_wywolania:      <widoczne bez interakcji | wyzwalacz ▼ | polecenie języka
                          naturalnego w Chat Window>
```

```
# Szablon — przypisanie warstw widoczności do roli (rozdz. 10)
warstwy_widocznosci_roli:
  rola:                  <Operator | użytkownik podstawowy | użytkownik zaawansowany |
                          administrator>
  najwyzsza_warstwa:     <1 | 2 | 3 | 4>
  warstwa_1:             <elementy widoczne bez interakcji>
  warstwa_2:             <wyzwalacze ▼, ☰ i znaczniki kontekstowe>
  warstwa_3:             <akcje menu ⋮ i rozwinięć kontekstowych>
  warstwa_4:             <funkcje eksperckie i tryb administracyjny>
  kanaly_funkcji_ukrytych:
    jezyk_naturalny:     <wł. | wył.>
    skrot_klawiszowy:    <wł. | wył.>
    wyszukiwarka_funkcji: <wł. | wył.>
```

```
# Szablon — klasa zdarzeń (rozdz. 7)
klasa_zdarzen:
  nazwa:               <jedna z klas rozdz. 7.2, np. Zakończenie>
  zrodlo_zdarzenia:    <polecenie/zdarzenie kontraktu komunikacji źródłowego>
  kanaly:
    centrum_powiadomien: zawsze          # kanał podstawowy — rozdz. 7.3
    mobile:              <wł. | wył.>
    email:               <wł. | wył.>
  aktywna_gdy_przelacznik_glowny_wylaczony: nie       # wygaszona, nie skasowana
```

---

## Załącznik C. Mapa zgodności z dokumentami źródłowymi

| Rozdział niniejszego dokumentu | Dokument źródłowy | Rozdział źródła |
|---|---|---|
| 1. Miejsce Ustawień w architekturze | Koncepcja platformy; Model konfiguracji; Specyfikacja okien operacyjnych | 5.3–5.4; 3.1, 5.1; 8.1 |
| 2. Struktura, nawigacja i warstwy widoczności | Model konfiguracji; Koncepcja platformy | 3.2, 3.3; 12 |
| 3. Konto Operatora | Bezpieczeństwo i uwierzytelnianie | 3, 4.2, 11.1 |
| 4. Uwierzytelnianie | Bezpieczeństwo i uwierzytelnianie | 6, 8, 13 |
| 5. Wygląd: motyw i język | System wizualny; Model konfiguracji | 9; 5.1 |
| 6. Urządzenia | Bezpieczeństwo i uwierzytelnianie; [Mobile](../funkcje-globalne/mobile.md); Architektura | 3, 11.2, 13; 9; 13 |
| 7. Powiadomienia | Model konfiguracji; System wizualny | 5.1; 8.10 |
| 8. Zarządzanie kontami modeli | Integracja modeli; Bezpieczeństwo i uwierzytelnianie | 3, 5, 7, 8, 10; 10 |
| 9. Okna komunikacji operacyjnej | Koncepcja platformy; Architektura | 1, 11.8; 11 |
| 10. Warstwy widoczności według roli | Koncepcja platformy; System wizualny | 12; 10 |
| 11. Always On Display | funkcje-globalne/always-on-display.md; Model konfiguracji; Koncepcja platformy | 3, 5, 9.5, 10; 5.17; 12.2 |
| 12–14. Makieta, katalog elementów, stany | System wizualny; Specyfikacja okien operacyjnych | 8, 11.3; 2.4, 6.0 |
| 15. Komunikacja | Bezpieczeństwo i uwierzytelnianie; Model konfiguracji; Architektura | 12; 8; 11 |
| 16. Zgodność z zasadami nadrzędnymi | Koncepcja platformy | 12 |

---

*Koniec dokumentu. Okno Ustawień (poziom aplikacji) — Specyfikacja docelowa, wersja 2.0, 2026-08-20.*

---
*Danaco Console — Platforma AI Workspace OS · v2.0 · status Deweloperski*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz.*
*Warunki korzystania: [Licencja produktu](../LICENSE.md). Kontakt: support@danaco-group.pl*
