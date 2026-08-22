# Danaco Console — Funkcja globalna Mobile

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
| **Tytuł** | Funkcja globalna Mobile |
| **Klasa dokumentu** | Specyfikacja docelowa |
| **Odbiorcy** | deweloper · projektant |
| **Przeznaczenie** | Pełnozakresowy materiał źródłowy do projektowania i budowy interfejsu funkcji globalnej Mobile oraz do implementacji protokołu parowania urządzeń: co istnieje, gdzie leży, w jakiej formie i do czego służy — dla każdego ekranu, każdego elementu, każdego komunikatu i każdego przepływu |
| **Zakres** | przeznaczenie i kontekst funkcji, zakres funkcjonalny na urządzeniu przenośnym, układ interfejsu mobilnego, komplet okien i ekranów, oba kanały komunikacji operacyjnej, warstwy widoczności, praca bez połączenia, powiadomienia wypychane, protokół parowania urządzenia, punkty sterowania, stany i dane, scenariusze użycia |
| **Poza zakresem** | wnętrze modułów i środowisk osiągalnych z Mobile — opisują je [Środowisko MultitaskingAI](../srodowiska/multitaskingai.md) i opracowania katalogu `moduly/`; mechanizm uwierzytelniania Operatora jako taki — [Bezpieczeństwo i uwierzytelnianie](../architektura/bezpieczenstwo-i-uwierzytelnianie.md) |
| **Dokument nadrzędny** | [Koncepcja platformy](../architektura/koncepcja-platformy.md) |
| **Dokumenty powiązane** | [Bezpieczeństwo i uwierzytelnianie](../architektura/bezpieczenstwo-i-uwierzytelnianie.md) · [Kontrakty komunikacji](../architektura/kontrakty-komunikacji.md) · [Model danych](../architektura/model-danych.md) · [Strona główna i nawigacja](../interfejs-uzytkownika/strona-glowna-i-nawigacja.md) · [System wizualny](../interfejs-uzytkownika/system-wizualny.md) · [Model konfiguracji](../architektura/model-konfiguracji.md) · [Ustawienia](../interfejs-uzytkownika/ustawienia.md) · [Elementy okien](../interfejs-uzytkownika/elementy-okien.md) · [Przepływ okien](../interfejs-uzytkownika/przeplyw-okien.md) · [Środowisko MultitaskingAI](../srodowiska/multitaskingai.md) |
| **Prototypy odniesienia** | `design/05-okna/platformowe/mobile.html` |
| **Źródła normatywne** | `budowa/shared/contract.json` (obszar `mobile`) · `design/zasoby/zetony/zetony.css` · `design/zasoby/css/komponenty.css`, `design/zasoby/rama.css`, `design/zasoby/prototyp.css` · `design/05-okna/platformowe/mobile.html` |
| **Zasada nadrzędna** | Danaco Console nie narzuca twardych blokad, bram bezpieczeństwa ani wymuszonych zgód; domyślne zachowanie systemu to wykonanie polecenia; żaden element interfejsu nie traci klikalności — niekompletność sygnalizowana jest komunikatem kontekstowym, nie wyłączeniem kontrolki; izolacja techniczna, izolacja kontekstu i wszelkie ograniczenia uprawnień są ustawieniami konfiguracyjnymi Operatora ze stanem wyjściowym „pełny dostęp” |

---

## Spis treści

1. [Wprowadzenie](#wprowadzenie)
2. [Przeznaczenie i kontekst funkcji](#1-przeznaczenie-i-kontekst-funkcji)
   - [1.1 Przedmiot funkcji](#11-przedmiot-funkcji)
   - [1.2 Relacja do stanowiska stacjonarnego](#12-relacja-do-stanowiska-stacjonarnego)
   - [1.3 Relacja do środowisk platformy](#13-relacja-do-środowisk-platformy)
   - [1.4 Punkty wejścia funkcji](#14-punkty-wejścia-funkcji)
3. [Zakres funkcjonalny na urządzeniu przenośnym](#2-zakres-funkcjonalny-na-urządzeniu-przenośnym)
   - [2.1 Zakres dostępny na urządzeniu przenośnym](#21-zakres-dostępny-na-urządzeniu-przenośnym)
   - [2.2 Zakres realizowany na stanowisku stacjonarnym](#22-zakres-realizowany-na-stanowisku-stacjonarnym)
   - [2.3 Granica funkcji](#23-granica-funkcji)
4. [Układ interfejsu mobilnego](#3-układ-interfejsu-mobilnego)
   - [3.1 Zasada układu](#31-zasada-układu)
   - [3.2 Punkty przełamania](#32-punkty-przełamania)
   - [3.3 Zachowanie każdej kolumny przy zwężaniu](#33-zachowanie-każdej-kolumny-przy-zwężaniu)
   - [3.4 Przełącznik kolumn](#34-przełącznik-kolumn)
5. [Komplet okien i ekranów mobilnych](#4-komplet-okien-i-ekranów-mobilnych)
   - [4.0 Zestawienie ekranów](#40-zestawienie-ekranów)
   - [4.1 Ekran startowy](#41-ekran-startowy)
   - [4.2 Logowanie](#42-logowanie)
   - [4.3 Strona główna Mobile](#43-strona-główna-mobile)
   - [4.4 Okno komunikacji — Chat Window](#44-okno-komunikacji--chat-window)
   - [4.5 Okno pętli wykonawczej — Execution Loop Window](#45-okno-pętli-wykonawczej--execution-loop-window)
   - [4.6 Przegląd zadań i procesów](#46-przegląd-zadań-i-procesów)
   - [4.7 Powiadomienia](#47-powiadomienia)
   - [4.8 Ustawienia urządzenia](#48-ustawienia-urządzenia)
6. [Kanały komunikacji operacyjnej na urządzeniu przenośnym](#5-kanały-komunikacji-operacyjnej-na-urządzeniu-przenośnym)
   - [5.1 Chat Window — kanał Użytkownik ↔ Wykonawca](#51-chat-window--kanał-użytkownik--wykonawca)
   - [5.2 Execution Loop Window — kanał Koordynator ↔ Wykonawca](#52-execution-loop-window--kanał-koordynator--wykonawca)
   - [5.3 Punkty obowiązkowego zatwierdzenia](#53-punkty-obowiązkowego-zatwierdzenia)
7. [Warstwy widoczności na urządzeniu przenośnym](#6-warstwy-widoczności-na-urządzeniu-przenośnym)
   - [6.1 Zasada nadrzędna](#61-zasada-nadrzędna)
   - [6.2 Cztery warstwy w funkcji Mobile](#62-cztery-warstwy-w-funkcji-mobile)
   - [6.3 Mechanizmy ukrywania funkcjonalności](#63-mechanizmy-ukrywania-funkcjonalności)
   - [6.4 Zasada jednego kliknięcia](#64-zasada-jednego-kliknięcia)
8. [Praca bez połączenia](#7-praca-bez-połączenia)
   - [7.1 Zakres dostępny lokalnie](#71-zakres-dostępny-lokalnie)
   - [7.2 Kolejkowanie działań](#72-kolejkowanie-działań)
   - [7.3 Uzgadnianie stanu po odzyskaniu połączenia](#73-uzgadnianie-stanu-po-odzyskaniu-połączenia)
   - [7.4 Rozstrzyganie konfliktów](#74-rozstrzyganie-konfliktów)
9. [Powiadomienia wypychane](#8-powiadomienia-wypychane)
   - [8.1 Klasy zdarzeń](#81-klasy-zdarzeń)
   - [8.2 Treść powiadomienia](#82-treść-powiadomienia)
   - [8.3 Działania z poziomu powiadomienia](#83-działania-z-poziomu-powiadomienia)
   - [8.4 Relacja do ustawień powiadomień](#84-relacja-do-ustawień-powiadomień)
   - [8.5 Kanały dostarczenia](#85-kanały-dostarczenia)
10. [Protokół parowania urządzenia](#9-protokół-parowania-urządzenia)
   - [9.1 Przeznaczenie protokołu](#91-przeznaczenie-protokołu)
   - [9.2 Sekwencja kroków](#92-sekwencja-kroków)
   - [9.3 Diagram protokołu](#93-diagram-protokołu)
   - [9.4 Kod parowania — prezentacja i ważność](#94-kod-parowania--prezentacja-i-ważność)
   - [9.5 Wymiana i przechowywanie poświadczeń](#95-wymiana-i-przechowywanie-poświadczeń)
   - [9.6 Rejestr sparowanych urządzeń](#96-rejestr-sparowanych-urządzeń)
   - [9.7 Wygaśnięcie poświadczenia](#97-wygaśnięcie-poświadczenia)
   - [9.8 Wylogowanie zdalne i wylogowanie z urządzenia](#98-wylogowanie-zdalne-i-wylogowanie-z-urządzenia)
   - [9.9 Postępowanie przy utracie urządzenia](#99-postępowanie-przy-utracie-urządzenia)
   - [9.10 Tabela komunikatów protokołu](#910-tabela-komunikatów-protokołu)
   - [9.11 Warunki błędu](#911-warunki-błędu)
11. [Punkty sterowania z okna konfiguracji i z okna ustawień](#10-punkty-sterowania-z-okna-konfiguracji-i-z-okna-ustawień)
   - [10.1 Punkty sterowania w oknie Ustawień](#101-punkty-sterowania-w-oknie-ustawień)
   - [10.2 Punkty sterowania w oknie konfiguracji](#102-punkty-sterowania-w-oknie-konfiguracji)
   - [10.3 Punkty sterowania na urządzeniu przenośnym](#103-punkty-sterowania-na-urządzeniu-przenośnym)
12. [Stany, dane i powiązania](#11-stany-dane-i-powiązania)
   - [11.1 Stany funkcji](#111-stany-funkcji)
   - [11.2 Stany elementów](#112-stany-elementów)
   - [11.3 Encje modelu danych](#113-encje-modelu-danych)
   - [11.4 Powiązania między oknami](#114-powiązania-między-oknami)
   - [11.5 Powiązania z pozostałymi dokumentami zbioru](#115-powiązania-z-pozostałymi-dokumentami-zbioru)
13. [Scenariusze użycia](#12-scenariusze-użycia)
   - [12.1 Sparowanie nowego telefonu](#121-sparowanie-nowego-telefonu)
   - [12.2 Zatwierdzenie kroku nieodwracalnego poza stanowiskiem](#122-zatwierdzenie-kroku-nieodwracalnego-poza-stanowiskiem)
   - [12.3 Przerwanie zadania z urządzenia przenośnego](#123-przerwanie-zadania-z-urządzenia-przenośnego)
   - [12.4 Praca bez łączności i uzgodnienie po powrocie](#124-praca-bez-łączności-i-uzgodnienie-po-powrocie)
   - [12.5 Nadzór nad przebiegiem MultitaskingAI](#125-nadzór-nad-przebiegiem-multitaskingai)
   - [12.6 Utrata telefonu i odwołanie dostępu](#126-utrata-telefonu-i-odwołanie-dostępu)
14. [Załącznik A. Skróty, gesty i ikonografia](#załącznik-a-skróty-gesty-i-ikonografia)
   - [A.1. Gesty urządzenia przenośnego](#a1-gesty-urządzenia-przenośnego)
   - [A.2. Skróty klawiatury zewnętrznej](#a2-skróty-klawiatury-zewnętrznej)
   - [A.3. Ikonografia funkcji](#a3-ikonografia-funkcji)
15. [Załącznik B. Zbiorcza tabela komunikatów funkcji Mobile](#załącznik-b-zbiorcza-tabela-komunikatów-funkcji-mobile)
16. [Załącznik C. Mapa zgodności z dokumentami źródłowymi](#załącznik-c-mapa-zgodności-z-dokumentami-źródłowymi)
17. [Załącznik — pełny wykaz komend kontraktu obszaru `mobile`](#załącznik--pełny-wykaz-komend-kontraktu-obszaru-mobile)
   - [Obszar `mobile` — 3 komendy](#obszar-mobile--3-komendy)

---

## Wprowadzenie

Mobile jest jedną z dwóch funkcji globalnych platformy Danaco Console (Koncepcja platformy, rozdz. 12). Nie tworzy własnej, izolowanej przestrzeni roboczej, nie podlega mechanizmowi przełączania i przeładowania właściwemu środowiskom i modułom, nie zakłada nowych sesji. Jest globalnym trybem dostępu do platformy z urządzenia przenośnego — telefonu albo tabletu — obejmującym dostęp mobilny, monitoring procesów, zarządzanie zadaniami i nadzór operacyjny nad pracą Wykonawcy prowadzoną w tle.

Oba kanały komunikacji operacyjnej są na urządzeniu przenośnym elementami pierwszoplanowymi. **Chat Window** — kanał Użytkownik ↔ Wykonawca — jest centralnym punktem pracy i podstawowym mechanizmem sterowania wszystkimi procesami platformy; na urządzeniu przenośnym zajmuje lewą kolumnę obszaru roboczego. **Execution Loop Window** — kanał Koordynator ↔ Wykonawca — prowadzi pętlę wykonawczą, koordynuje zadania, nadzoruje realizację i kontroluje przebieg procesów; otwiera się jako kolumna sąsiadująca z Chat Window i pozwala zatwierdzać punkty decyzyjne oraz przerywać zadania z urządzenia przenośnego.

Dokument czyta się w pięciu warstwach nałożonych na siebie:

| Warstwa dokumentu | Co dostarcza | Rozdziały |
|---|---|---|
| Kontekst | Miejsce funkcji w platformie, relacja do stanowiska stacjonarnego i do środowisk, zakres funkcjonalny | 1–2 |
| Forma | Układ kolumnowy, punkty przełamania, komplet ekranów i katalogi elementów | 3–4 |
| Praca operacyjna | Oba kanały komunikacji, warstwy widoczności, praca bez połączenia, powiadomienia | 5–8 |
| Dostęp urządzenia | Pełny protokół parowania, rejestr urządzeń, poświadczenia, odwołanie dostępu | 9 |
| Sterowanie i dane | Punkty konfiguracji, stany, encje modelu danych, scenariusze, załączniki | 10–12, załączniki |

---

## 1. Przeznaczenie i kontekst funkcji

### 1.1. Przedmiot funkcji

Mobile odpowiada na potrzebę zachowania pełnej kontroli nad platformą poza stanowiskiem roboczym — nad procesami długotrwałymi, uruchomionymi w tle automatyzacjami i projektami środowiska MultitaskingAI, które wymagają zatwierdzenia, wstrzymania albo interwencji niezależnie od położenia Użytkownika (Koncepcja platformy, rozdz. 12.1).

| Możliwość | Zakres na urządzeniu przenośnym |
|---|---|
| Dostęp mobilny | Obsługa platformy z urządzenia przenośnego, niezależna od aktywnego środowiska |
| Monitoring procesów | Podgląd projektów, automatyzacji, agentów, aplikacji i sesji |
| Zarządzanie zadaniami | Uruchamianie, zatrzymywanie, zatwierdzanie, wstrzymywanie, wznawianie i modyfikowanie procesów |
| Nadzór operacyjny | Kontrola pracy Wykonawcy prowadzonej w tle, nadzór nad pętlą wykonawczą, zatwierdzanie punktów decyzyjnych |
| Sterowanie językiem naturalnym | Chat Window jako podstawowy mechanizm sterowania wszystkimi procesami platformy |

**Kluczowa idea:** zarządzanie platformą z dowolnego miejsca.

### 1.2. Relacja do stanowiska stacjonarnego

Serwer jest jedynym źródłem prawdy (Architektura, rozdz. 11; Kontrakty komunikacji, rozdz. 1). Urządzenie przenośne prezentuje obraz pochodzący z serwera i nie utrzymuje własnego, rozbieżnego stanu. Stanowisko stacjonarne i urządzenie przenośne są dwoma równorzędnymi klientami tego samego konta, różniącymi się wyłącznie powierzchnią ekranu i zestawem ujawnionych funkcji.

| Wymiar | Stanowisko stacjonarne | Urządzenie przenośne |
|---|---|---|
| Jedyne źródło prawdy | Serwer | Serwer |
| Rola klienta | Pełna przestrzeń robocza środowisk i modułów | Globalny tryb dostępu, nadzór i sterowanie |
| Chat Window | Lewa kolumna obszaru roboczego, pełna wysokość | Lewa kolumna obszaru roboczego, pełna wysokość |
| Execution Loop Window | Kolumna sąsiadująca, otwierana | Kolumna sąsiadująca, otwierana |
| Zakładanie sesji i kart | Tak | Nie — funkcja globalna nie tworzy nowych sesji |
| Uwierzytelnienie | Login i hasło, metody dodatkowe (Bezpieczeństwo i uwierzytelnianie, rozdz. 5–6) | Poświadczenie urządzenia wydane w parowaniu (rozdz. 9) |
| Zależność od aktywnego środowiska | Praca wewnątrz wybranego środowiska | Brak — dostęp niezależny od aktywnego środowiska |

Trwałość stanu procesów po stronie serwera jest bezpośrednią podstawą funkcji: zamknięcie widoku na stanowisku stacjonarnym nie przerywa pętli wykonawczej, a urządzenie przenośne podłącza się do stanu zastanego (Strona główna i nawigacja, rozdz. 8.4; Kontrakty komunikacji, rozdz. 2, krok rozłączenia).

### 1.3. Relacja do środowisk platformy

Mobile działa ponad czterema środowiskami i piętnastoma modułami. Nie przełącza środowiska i nie przeładowuje przestrzeni roboczej — prezentuje przekrój procesów prowadzonych we wszystkich środowiskach jednocześnie.

```
STANOWISKO STACJONARNE                      URZĄDZENIE PRZENOŚNE
    strona główna                                Mobile
        │                                          │
        ├── TalkIn ─────┐                          │  przekrój procesów
        ├── CodeStudio ─┤                          │  wszystkich środowisk
        ├── Workspace ──┤   procesy w tle          │  bez przełączania
        └── Multitasking┘   trwałe na serwerze ────┤  środowiska
                                                   │
PONAD CAŁĄ STRUKTURĄ                               ├── Chat Window
    Always On Display — globalny agent             ├── Execution Loop Window
    Mobile — nadzór i kontynuacja pracy            └── przegląd zadań i procesów
             poza stanowiskiem
```

| Środowisko | Co Mobile udostępnia z tego środowiska |
|---|---|
| [Środowisko TalkIn](../srodowiska/talkin.md) | Podgląd i sterowanie pętlą wykonawczą kart treści, zatwierdzanie wyników redakcji i przekładu |
| [Środowisko CodeStudio](../srodowiska/codestudio.md) | Nadzór nad przebiegami budowy, testów i zadań repozytorium, zatwierdzanie kroków nieodwracalnych |
| [Środowisko WorkSpace](../srodowiska/workspace.md) | Podgląd projektów, zasobów i zadań porządkowych, sterowanie automatyzacjami |
| [Środowisko MultitaskingAI](../srodowiska/multitaskingai.md) | Monitor procesu, przebiegi zespołów i ról, decyzje procesu, kolejki i orkiestracja |

### 1.4. Punkty wejścia funkcji

Mobile jest funkcją globalną — obecność w listwie ustawień strony głównej jest jednym z kilku równoważnych punktów dostępu, nie jedynym (Strona główna i nawigacja, rozdz. 9.2).

| Punkt wejścia | Miejsce | Warstwa | Sposób wywołania |
|---|---|---|---|
| Pozycja „Mobile” w listwie ustawień | Strefa 3 strony głównej ([Strona główna i nawigacja](../interfejs-uzytkownika/strona-glowna-i-nawigacja.md), rozdz. 3.4) | 1 | Kliknięcie pozycji listwy |
| Ikona „Mobile” w szybkim wyborze szyny nawigacji | Każde środowisko i moduł ([Elementy okien](../interfejs-uzytkownika/elementy-okien.md), rozdz. 5) | 2 | Szybki wybór szyny nawigacji albo menu kebab powłoki (⋮) |
| Uruchomienie aplikacji na urządzeniu sparowanym | Urządzenie przenośne | 1 | Otwarcie aplikacji |
| Powiadomienie wypychane | Urządzenie przenośne, poza aplikacją | 1 | Dotknięcie powiadomienia |
| Polecenie języka naturalnego | Chat Window na dowolnym urządzeniu | 4 | Polecenie „pokaż stan procesów na urządzeniu przenośnym” |

---

## 2. Zakres funkcjonalny na urządzeniu przenośnym

### 2.1. Zakres dostępny na urządzeniu przenośnym

| Obszar platformy | Zakres na urządzeniu przenośnym | Uzasadnienie zakresu |
|---|---|---|
| Chat Window | Pełny: polecenie w języku naturalnym, strumień odpowiedzi, zatwierdzanie, przerywanie, wyjaśnianie wyniku i kontekstu | Podstawowy mechanizm sterowania wszystkimi procesami platformy |
| Execution Loop Window | Pełny nadzór: dekompozycja zlecenia, kolejka i stan zadań, komunikaty sterujące, kontrola jakości, wstrzymanie, wznowienie, przerwanie, korekta zlecenia | Nadzór nad pętlą wykonawczą jest istotą funkcji |
| Przegląd zadań i procesów | Pełny: projekty, automatyzacje, agenci, aplikacje, sesje, przebiegi MultitaskingAI, kolejki | Monitoring procesów wskazany jako możliwość funkcji |
| Zarządzanie zadaniami | Pełne: `start`, `stop`, `approve`, `pause`, `resume`, `modify` | Zarządzanie zadaniami wskazane jako możliwość funkcji |
| Decyzje procesu MultitaskingAI | Pełne: zatwierdzanie punktów decyzyjnych, odrzucenie wyniku, decyzja o ponowieniu | Nadzór operacyjny nad pracą prowadzoną w tle |
| Powiadomienia | Pełne: odbiór, przegląd, działanie z poziomu powiadomienia | Powiadomienie jest wejściem do interwencji |
| Ustawienia urządzenia | Nazwa urządzenia, zakres powiadomień, tryb pracy bez połączenia, wylogowanie urządzenia | Zakres dotyczy tego urządzenia, nie platformy |
| Podgląd artefaktu | Odczyt treści artefaktu, wersji i metadanych | Podgląd wyniku niezbędny do decyzji o zatwierdzeniu |

### 2.2. Zakres realizowany na stanowisku stacjonarnym

| Obszar platformy | Miejsce realizacji | Dlaczego stanowisko stacjonarne |
|---|---|---|
| Zakładanie i zamykanie kart sesji, przełączanie środowisk i modułów | Powłoka środowiska ([Przepływ okien](../interfejs-uzytkownika/przeplyw-okien.md)) | Funkcja globalna nie tworzy przestrzeni roboczej ani nowych sesji (Koncepcja platformy, rozdz. 5.3) |
| Okna operacyjne modułów — edytory, kreatory, zarządcy, przeglądarki | Obszar roboczy modułu | Praca redakcyjna i konstrukcyjna wymaga powierzchni obszaru roboczego i wielu kolumn jednocześnie |
| Okno konfiguracji i wszystkie jego zakresy | [Model konfiguracji](../architektura/model-konfiguracji.md) | Konfiguracja obejmuje warstwy globalna → środowisko → projekt → sesja i macierze o gęstości nieprzenoszalnej na wąską kolumnę |
| Okno punktów izolacji | [Izolacja i konfigurowalność zależności](../architektura/izolacja-i-zaleznosci.md) | Macierz izolacji technicznej i kontekstu jest zestawem przełączników wymagającym pełnej szerokości obszaru roboczego |
| Panel orkiestracji MultitaskingAI w pełnej postaci | [Środowisko MultitaskingAI](../srodowiska/multitaskingai.md) | Panel prowadzi równolegle sekcje ról, kolejek i zależności; urządzenie przenośne prezentuje z niego Monitor procesu i decyzje |
| Rejestracja konta, zmiana hasła, metody dodatkowe uwierzytelniania | [Bezpieczeństwo i uwierzytelnianie](../architektura/bezpieczenstwo-i-uwierzytelnianie.md), rozdz. 4–6 | Konto właściciela zakładane jest raz, przy pierwszym uruchomieniu platformy |
| Inicjowanie parowania nowego urządzenia | [Okno Ustawień](../interfejs-uzytkownika/ustawienia.md), sekcja Urządzenia | Parowanie zakładane jest w obecności Operatora przy maszynie już uwierzytelnionej (rozdz. 9.2) |
| Definiowanie agentów, przepływów pracy, automatyzacji i rozszerzeń | Moduły Agents, Automations, [Rozszerzenia](../architektura/rozszerzenia.md) | Definiowanie jest pracą konstrukcyjną; urządzenie przenośne steruje uruchomionymi przebiegami |

### 2.3. Granica funkcji

| Warstwa | Za co odpowiada na urządzeniu przenośnym |
|---|---|
| Funkcja globalna Mobile | Dostęp niezależny od środowiska, oba kanały komunikacji operacyjnej, przegląd i sterowanie procesami, powiadomienia, ustawienia tego urządzenia |
| Serwer | Stan sesji, procesów, kolejek i pętli wykonawczych; nadawanie identyfikatorów; rozgłaszanie zmian do wszystkich urządzeń konta |
| Stanowisko stacjonarne | Zakładanie przestrzeni roboczych, praca w oknach operacyjnych modułów, konfiguracja platformy |

---

## 3. Układ interfejsu mobilnego

### 3.1. Zasada układu

Na urządzeniu przenośnym obowiązuje ten sam standard projektowy co na stanowisku stacjonarnym: **układ pionowy — podział lewa–prawa**. Wszystkie okna rozmieszczone są w kolumnach sąsiadujących poziomo, na pełną wysokość obszaru roboczego. Regulacji podlega wyłącznie szerokość kolumn. Zwężenie ekranu zwęża i zwija kolumny w kolejności od lewej do prawej, zachowując niezmienną kolejność i tożsamość każdej kolumny; żadna kolumna nie przechodzi pod inną i nie zmienia się w sekcję układaną nad inną sekcją.

Kanoniczne rozmieszczenie kolumn na urządzeniu przenośnym:

| Kolumna | Zawartość | Waga wizualna |
|---|---|---|
| Skrajna lewa, zwijalna | Kolumna nawigacji funkcji: strona główna Mobile, przegląd zadań i procesów, powiadomienia, ustawienia urządzenia | Kolumna ikonowa, pełna wysokość obszaru roboczego |
| Lewa, stała, pełna wysokość | **Chat Window** — główne okno komunikacji Użytkownik ↔ Wykonawca | Lewa kolumna, stała, pełna wysokość obszaru roboczego |
| Kolumna sąsiadująca (otwierana) | **Execution Loop Window** — okno pętli wykonawczej Koordynator ↔ Wykonawca | Kolumna sąsiadująca, otwierana, pełna wysokość obszaru roboczego |
| Prawa, dominująca | Obszar roboczy funkcji: przegląd zadań i procesów, powiadomienia, ustawienia urządzenia, podgląd artefaktu | Kolumna dominująca, największy udział powierzchni |
| Kolejne kolumny boczne | Okna pomocnicze — otwierane jako rozszerzenia boczne, po prawej stronie obszaru roboczego | Kolumna boczna, otwierana jako rozszerzenie boczne |

### 3.2. Punkty przełamania

Rozdział uzgadnia zachowanie kolumn z rozdziałem 4.4 dokumentu [System wizualny](../interfejs-uzytkownika/system-wizualny.md). Tamta tabela ustala zachowanie kolumn od 768 px wzwyż i wskazuje, że poniżej 768 px widoczna jest jedna kolumna naraz, wybierana przełącznikiem kolumn, w kolejności odpowiadającej kolejności kolumn od lewej do prawej. Niniejszy rozdział doprecyzowuje przedział poniżej 768 px na trzy punkty przełamania urządzenia przenośnego; wszystkie szerokości powyżej 768 px obowiązują bez zmian w brzmieniu Systemu wizualnego.

| Dostępna szerokość | Widoczne kolumny | Zachowanie kolumn |
|---|---|---|
| 1024 px i więcej (tablet w orientacji poziomej) | Nawigacja funkcji · Chat Window · obszar roboczy · panel pomocniczy | Wszystkie kolumny w szerokościach bazowych; Execution Loop Window otwierane jako kolumna sąsiadująca |
| 768–1024 px (tablet w orientacji pionowej) | Nawigacja funkcji (ikonowa) · Chat Window · obszar roboczy | Execution Loop Window zwinięte do wyzwalacza w nagłówku Chat Window; Chat Window i obszar roboczy zajmują po połowie szerokości |
| 600–768 px (tablet mały, telefon w orientacji poziomej) | Chat Window · obszar roboczy | Kolumna nawigacji funkcji zwinięta do wyzwalacza `☰` w nagłówku Chat Window; panel pomocniczy zwinięty do wyzwalacza ikonowego w kolumnie skrajnej prawej |
| 420–600 px (telefon w orientacji poziomej) | Jedna kolumna naraz, z podglądem krawędzi kolumny sąsiadującej | Kolumna aktywna zajmuje pełną szerokość pomniejszoną o pas podglądu 24 px kolumny następnej; przełącznik kolumn prowadzi po kolejności od lewej do prawej |
| Poniżej 420 px (telefon w orientacji pionowej) | Jedna kolumna naraz | Kolumna aktywna zajmuje pełną szerokość; przełącznik kolumn prowadzi po kolejności od lewej do prawej; kolumny zwinięte pozostają na swoich miejscach w kolejności |

Zwinięcie i rozwinięcie kolumny przebiega jako zmiana szerokości w czasie tokenu `--dn-czas-3` (System wizualny, rozdz. 3a.5, 4.4). Zwinięcie kolumny zwalnia zajmowaną powierzchnię w całości — kolumna nie pozostawia rezerwy miejsca ani śladu obrysu.

### 3.3. Zachowanie każdej kolumny przy zwężaniu

| Kolumna | Kolejność zwijania | Postać zwężona | Postać zwinięta | Przywołanie kolumny zwiniętej |
|---|---|---|---|---|
| Nawigacja funkcji | 1 — zwijana jako pierwsza | Kolumna ikonowa 56 px bez etykiet tekstowych | Wyzwalacz `☰` w nagłówku Chat Window | Kliknięcie `☰`; kolumna wysuwa się na czas wyboru i zwija samoczynnie po wyborze pozycji |
| Panel pomocniczy | 2 | Szerokość minimalna 280 px | Wyzwalacz ikonowy w skrajnej prawej krawędzi obszaru roboczego | Kliknięcie wyzwalacza ikonowego |
| Execution Loop Window | 3 | Szerokość minimalna 20% obszaru roboczego | Wyzwalacz „Execution Loop ▼” w nagłówku Chat Window | Kliknięcie wyzwalacza, przesunięcie poziome z krawędzi Chat Window albo polecenie w Chat Window |
| Chat Window | 4 | Szerokość minimalna 24% obszaru roboczego | Pozycja w przełączniku kolumn; kolumna zachowuje pierwszą pozycję w kolejności | Wybór w przełączniku kolumn, przesunięcie poziome albo skrót `Ctrl/Cmd + L` na klawiaturze zewnętrznej |
| Obszar roboczy funkcji | 5 — nigdy nie zwijany do zera | Pozostała szerokość, minimum 40% przy dwóch kolumnach | Pozycja w przełączniku kolumn | Wybór w przełączniku kolumn albo przesunięcie poziome |

### 3.4. Przełącznik kolumn

Przełącznik kolumn jest elementem nagłówka obszaru roboczego. Prezentuje kolumny w niezmiennej kolejności od lewej do prawej i zaznacza kolumnę aktywną. Zmiana kolumny nie zmienia stanu żadnej z kolumn pozostałych — kolumna niewidoczna zachowuje swoją zawartość, strumień i pozycję przewinięcia.

```
 ┌─ Przełącznik kolumn ─────────────────────────────
 │   ☰   │  Chat  │  Loop  │  Obszar  │  Panel  ⋮
 │  nawi │  ●     │        │  roboczy │
 └──────────────────────────────────────────────────
     1        2        3         4         5
     kolejność niezmienna, zgodna z układem lewa–prawa
```

---

## 4. Komplet okien i ekranów mobilnych

Funkcja udostępnia osiem ekranów. Kolumna „Warstwa” podaje warstwę widoczności (rozdz. 6), kolumna „Sposób wywołania” — mechanizm otwarcia. Makiety rysowane są w stanie spoczynku interfejsu: widoczne są wyłącznie elementy warstwy 1 oraz zwinięte wyzwalacze warstw 2–3.

### 4.0. Zestawienie ekranów

| Ekran | Rola | Waga wizualna | Warstwa | Sposób wywołania |
|---|---|---|---|---|
| Ekran startowy | Nawiązanie połączenia, rozpoznanie urządzenia, przejście do strony głównej Mobile albo do logowania | Pełna powierzchnia | 1 | Uruchomienie aplikacji na urządzeniu przenośnym |
| Logowanie | Uwierzytelnienie urządzenia bez ważnego poświadczenia; wejście do potwierdzenia parowania | Pełna powierzchnia | 1 | Brak ważnego poświadczenia urządzenia (rozdz. 9.10) |
| Strona główna Mobile | Przekrój stanu platformy: procesy wymagające decyzji, przebiegi w toku, skróty do kanałów komunikacji | Kolumna dominująca | 1 | Zakończenie ekranu startowego; pozycja nawigacji funkcji |
| Okno komunikacji (Chat Window) | Kanał Użytkownik ↔ Wykonawca — sterowanie wszystkimi procesami | Lewa kolumna, stała, pełna wysokość obszaru roboczego | 1 | Widoczne bez interakcji; pozycja przełącznika kolumn |
| Okno pętli wykonawczej (Execution Loop Window) | Kanał Koordynator ↔ Wykonawca — nadzór, punkty decyzyjne, przerywanie zadań | Kolumna sąsiadująca, otwierana | 2 | Wyzwalacz „Execution Loop ▼”, przesunięcie poziome, polecenie w Chat Window |
| Przegląd zadań i procesów | Projekty, automatyzacje, agenci, aplikacje, sesje, przebiegi, kolejki | Kolumna dominująca | 1 | Pozycja nawigacji funkcji |
| Powiadomienia | Lista zdarzeń wymagających uwagi wraz z działaniami | Kolumna dominująca | 1 | Pozycja nawigacji funkcji; dotknięcie powiadomienia wypychanego |
| Ustawienia urządzenia | Nazwa urządzenia, zakres powiadomień, praca bez połączenia, wylogowanie urządzenia | Kolumna dominująca | 2 | Pozycja nawigacji funkcji |

### 4.1. Ekran startowy

```
 ┌──────────────────────────────────┐
 │                                  │
 │            ◈ Danaco              │
 │            Console               │
 │                                  │
 │        ⇄  nawiązywanie           │
 │           połączenia             │
 │                                  │
 │        Urządzenie: Telefon       │
 │        Operatora                 │
 │                                  │
 │        ▮▮▮▮▮▯▯▯                  │
 │                                  │
 │        Pracuj bez połączenia  →  │
 │                                  │
 └──────────────────────────────────┘
```

| Element ekranu | Co to jest | Do czego służy | Stany | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|
| Godło platformy | Znak marki `ikony/logo-danaco.svg` | Identyfikacja platformy | Domyślny | 1 | Widoczne |
| Wskaźnik nawiązywania połączenia | Pasek postępu z etykietą stanu | Sygnał przebiegu kroków: kanał WebSocket, przedstawienie poświadczenia, powiązanie | Nawiązywanie · uwierzytelnianie · powiązanie · błąd | 1 | Widoczny |
| Nazwa rozpoznanego urządzenia | Etykieta z nazwą z rejestru urządzeń | Potwierdzenie, którym urządzeniem rejestru jest ten egzemplarz | Rozpoznane · nierozpoznane | 1 | Widoczna |
| Przejście „Pracuj bez połączenia” | Odnośnik do zakresu lokalnego | Wejście do zakresu dostępnego bez połączenia (rozdz. 7) | Domyślny · aktywny przy braku łączności | 1 | Widoczne; przy braku łączności podnoszone do pierwszej pozycji |
| Komunikat kontekstowy błędu | Wiersz tekstu pod wskaźnikiem | Wskazanie przyczyny: brak łączności, poświadczenie wygasłe, dostęp odwołany | Ukryty · widoczny | 1 | Pojawia się przy błędzie kroku |
| Działanie „Ponów połączenie” | Przycisk `.dn-btn` | Powtórzenie nawiązania połączenia | Domyślny · ładowanie | 2 | Widoczne po błędzie połączenia |
| Działanie „Sparuj to urządzenie” | Przycisk prowadzący do potwierdzenia parowania | Wejście urządzenia nierozpoznanego do protokołu parowania (rozdz. 9) | Domyślny · ładowanie | 2 | Widoczne, gdy urządzenie nie ma poświadczenia |
| Podgląd wersji kontraktu | Etykieta wersji uzgodnionej z serwerem | Diagnostyka zgodności klienta i serwera | Zgodna · niezgodna | 4 | Przytrzymanie godła platformy |

### 4.2. Logowanie

```
 ┌──────────────────────────────────┐
 │  ◈ Danaco Console                │
 ├──────────────────────────────────┤
 │  Logowanie urządzenia            │
 │                                  │
 │  Login                           │
 │  [                            ]  │
 │                                  │
 │  Hasło                           │
 │  [                            ]  │
 │                                  │
 │  [        Zaloguj            ]   │
 │                                  │
 │  Kod PIN  ·  Kod e-mail       ⋮  │
 │                                  │
 │  Mam kod parowania            →  │
 └──────────────────────────────────┘
```

| Element ekranu | Co to jest | Do czego służy | Stany | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|
| Pole „Login” | Pole `.dn-pole` z kontrolką `.dn-pole-kontrolka` | Podanie loginu konta właściciela | Puste · wypełnione · błąd walidacji | 1 | Widoczne |
| Pole „Hasło” | Pole `.dn-pole` z kontrolką `.dn-pole-kontrolka` w trybie maskowanym | Podanie hasła konta | Puste · wypełnione · błąd walidacji | 1 | Widoczne |
| Przycisk „Zaloguj” | Przycisk główny `.dn-btn--sygnal` | Uwierzytelnienie urządzenia i wydanie poświadczenia | Domyślny · ładowanie · błąd | 1 | Widoczny |
| Metoda „Kod PIN” | Pozycja metody dodatkowej | Szybkie potwierdzenie tożsamości na urządzeniu wcześniej uwierzytelnionym (Bezpieczeństwo, rozdz. 6.1) | Dostępna · niedostępna | 2 | Widoczna, gdy metoda jest aktywna dla konta |
| Metoda „Kod e-mail” | Pozycja metody dodatkowej | Logowanie kodem przesłanym na adres uwierzytelniający (Bezpieczeństwo, rozdz. 6.3) | Dostępna · niedostępna | 2 | Widoczna, gdy metoda jest aktywna dla konta |
| Przejście „Mam kod parowania” | Odnośnik do potwierdzenia parowania | Wejście do kroku 4 protokołu parowania (rozdz. 9.3) | Domyślny | 1 | Widoczne |
| Kebab ekranu `⋮` | Zbiorcze menu ekranu | Adres serwera, odzyskanie konta, diagnostyka połączenia | Domyślny · rozwinięty | 3 | Kliknięcie ikony ⋮ |
| Komunikat kontekstowy | Wiersz pod przyciskiem | Wskazanie przyczyny odrzucenia uwierzytelnienia | Ukryty · widoczny | 1 | Pojawia się po odrzuceniu |
| Diagnostyka połączenia | Widok adresu serwera, wersji kontraktu i przebiegu kroków połączenia | Rozpoznanie przyczyny niepowodzenia po stronie sieci | Zamknięta · otwarta | 4 | Pozycja kebaba ekranu albo polecenie języka naturalnego |

### 4.3. Strona główna Mobile

```
 ═══════════════════════════════════════════════════════════
  ☰ │ Chat Window            │ Strona główna Mobile      ⋮
    │ Użytkownik ↔ Wykonawca │
    │ [Danaco Console]       │  DECYZJE OCZEKUJĄCE     (2)
    │ [Multitasking] ▼       │   ⚑ Wdrożenie pakietu
    │                        │     Koordynator · 2 min
    │  Wykonawca             │   ⚑ Publikacja raportu
    │   wynik częściowy 62%  │     Koordynator · 14 min
    │   [Zatwierdź][Przerwij]│
    │                        │  PRZEBIEGI W TOKU       (3)
    │                        │   ⟳ Analiza źródeł  62%
    │                        │   ⟳ Budowa pakietu  30%
    │                        │   ⟳ Przekład opisu  81%
    │ ────────────────────── │
    │ pole polecenia    ⋮ ➤  │  SKRÓT: Execution Loop ▼
    │ Execution Loop ▼       │  ● połączony
 ═══════════════════════════════════════════════════════════
```

| Element ekranu | Co to jest | Do czego służy | Stany | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|
| Sekcja „Decyzje oczekujące” | Lista punktów decyzyjnych pętli wykonawczej i decyzji procesu | Wskazanie tego, co wstrzymuje pracę Wykonawcy | Pusta · zapełniona · zapełniona z pozycją przeterminowaną | 1 | Widoczna |
| Pozycja decyzji | Wiersz z opisem, źródłem i czasem oczekiwania | Wejście do zatwierdzenia albo odrzucenia | Oczekująca · zatwierdzona · odrzucona | 1 | Widoczna; kliknięcie otwiera Execution Loop Window na tym punkcie |
| Sekcja „Przebiegi w toku” | Lista przebiegów z wskaźnikiem zaawansowania | Orientacja w pracy prowadzonej w tle | Pusta · zapełniona | 1 | Widoczna |
| Pozycja przebiegu | Wiersz z nazwą, wskaźnikiem i środowiskiem pochodzenia | Wejście do szczegółów przebiegu | W toku · wstrzymany · błąd · zakończony | 1 | Widoczna; kliknięcie otwiera przegląd zadań i procesów na tym przebiegu |
| Wskaźnik połączenia | Kropka `.dn-kropka` z etykietą | Sygnał stanu kanału: połączony, wznawianie, bez połączenia | Połączony · wznawianie · bez połączenia | 1 | Widoczny |
| Skrót „Execution Loop ▼” | Zwinięty wyzwalacz okna pętli wykonawczej | Otwarcie kolumny nadzoru | Zwinięty · rozwinięty | 2 | Kliknięcie wyzwalacza, przesunięcie poziome albo polecenie w Chat Window |
| Licznik działań w kolejce lokalnej | Plakietka z liczbą działań oczekujących na wysłanie | Sygnał, że urządzenie pracuje bez połączenia (rozdz. 7) | Ukryty · widoczny | 1 | Widoczny, gdy kolejka lokalna nie jest pusta |
| Filtr przekroju | Element zbiorczy `Zakres ▼` | Zawężenie przekroju do środowiska, projektu albo rodzaju procesu | Domyślny · zawężony | 2 | Kliknięcie elementu zbiorczego |
| Kebab ekranu `⋮` | Zbiorcze menu ekranu | Odświeżenie przekroju, kolejność sekcji, zakres powiadomień, eksport zestawienia | Domyślny · rozwinięty | 3 | Kliknięcie ikony ⋮ |
| Wyszukiwarka funkcji | Pole wyszukiwania funkcji platformy | Dojście do funkcji warstw 2–4 jednym wpisaniem nazwy | Zamknięta · otwarta | 4 | Pozycja kebaba ekranu albo skrót klawiatury zewnętrznej |

### 4.4. Okno komunikacji — Chat Window

```
 ┌─ Chat Window · Użytkownik ↔ Wykonawca ─────────
 │  [Danaco Console] [Multitasking] [Fable 5]  ⋮
 ├────────────────────────────────────────────────
 │   Użytkownik
 │     wstrzymaj przebieg budowy pakietu
 │
 │   Wykonawca
 │     przebieg wstrzymany na zadaniu 3
 │     ◈ dziennik przebiegu
 │     [Wznów]  [Przerwij]
 │
 ├────────────────────────────────────────────────
 │   pole polecenia                       ⋮   ➤
 ├────────────────────────────────────────────────
 │   Execution Loop ▼    Koordynator ↔ Wykonawca
 └────────────────────────────────────────────────
```

| Element okna | Co to jest | Do czego służy | Stany | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|
| Nagłówek okna | Belka z nazwą kanału i znacznikami kontekstu | Identyfikacja kanału i bieżącego kontekstu | Domyślny | 1 | Widoczny |
| Znacznik platformy `[Danaco Console]` | Lekka etykieta paska kontekstu | Potwierdzenie zakresu globalnego rozmowy | Domyślny | 1 | Widoczny; kliknięcie otwiera selektor zakresu |
| Znacznik środowiska `[Multitasking]` | Lekka etykieta paska kontekstu | Wskazanie środowiska, którego dotyczy polecenie | Domyślny · zmieniony ręcznie | 1 | Widoczny; kliknięcie otwiera listę środowisk |
| Znacznik modelu `[Fable 5]` | Lekka etykieta paska kontekstu | Wskazanie Wykonawcy odpowiadającego w tym kanale | Domyślny · nadpisany | 1 | Widoczny; kliknięcie otwiera selektor modelu |
| Strumień rozmowy | Obszar zawartości z wypowiedziami Użytkownika i Wykonawcy | Prezentacja przebiegu pracy i wyników | Pusty · strumień aktywny · przewijalny · odtwarzany z pamięci lokalnej | 1 | Widoczny |
| Referencja zasobu `◈` | Odnośnik do artefaktu albo dziennika przebiegu | Podgląd wyniku przed decyzją | Domyślna · niedostępna bez połączenia | 1 | Widoczna w strumieniu; kliknięcie otwiera podgląd w kolumnie dominującej |
| Przycisk „Zatwierdź” | Kontrolka akceptacji kroku | Zatwierdzenie kroku oczekującego na decyzję | Domyślny · ładowanie · zatwierdzony · zakolejkowany lokalnie | 1 | Widoczny przy kroku oczekującym |
| Przycisk „Przerwij” | Kontrolka zatrzymania działania | Przerwanie operacji Wykonawcy w toku | Domyślny · ładowanie · zakolejkowany lokalnie | 1 | Widoczny podczas pracy Wykonawcy |
| Pole polecenia | Pole `.dn-pole` z kontrolką `.dn-pole-kontrolka` (wieloliniowa) i wejściem głosowym | Wydanie polecenia w języku naturalnym | Puste · wypełnione · wysyłanie · zakolejkowane lokalnie | 1 | Widoczne |
| Przycisk wyślij `➤` | Ikona wysłania treści pola | Przekazanie polecenia Wykonawcy | Domyślny · ładowanie | 1 | Widoczny; przy pustym polu wyświetla komunikat kontekstowy „Wpisz polecenie przed wysłaniem” |
| Dyktowanie polecenia | Wejście głosowe pola polecenia | Wydanie polecenia bez klawiatury ekranowej | Bezczynne · nasłuch · przetwarzanie | 2 | Kliknięcie ikony mikrofonu w polu polecenia |
| Kebab pola polecenia `⋮` | Zbiorcze menu pola | Załącznik, poziom wysiłku Wykonawcy, tryb odpowiedzi, kompresja kontekstu | Domyślny · rozwinięty | 3 | Kliknięcie ikony ⋮ |
| Nagłówek „Execution Loop ▼” | Zwinięty wyzwalacz okna pętli wykonawczej | Otwarcie kolumny sąsiadującej | Zwinięty · rozwinięty | 2 | Kliknięcie nagłówka, przesunięcie poziome albo polecenie w Chat Window |
| Kebab okna `⋮` | Zbiorcze menu okna | Zakres kontekstu, retencja historii, eksport rozmowy, szerokość kolumny | Domyślny · rozwinięty | 3 | Kliknięcie ikony ⋮ w nagłówku |
| Wywołanie funkcji warstwy 4 | Polecenie języka naturalnego uruchamiające funkcję ekspercką | Dostęp do trybu administracyjnego i narzędzi diagnostycznych | Rozpoznane · nierozpoznane | 4 | Treść polecenia w polu Chat Window |

### 4.5. Okno pętli wykonawczej — Execution Loop Window

```
 ┌─ Execution Loop Window · Koordynator ↔ Wykonawca ─
 │  zlecenie: „Budowa i publikacja pakietu”      ⋮
 ├───────────────────────────────────────────────────
 │  DEKOMPOZYCJA ZLECENIA
 │   1. Repozytorium · pobranie zmian   ✔ zakończone
 │   2. Budowa       · pakiet wydania   ✔ zakończone
 │   3. Publikacja   · wydanie pakietu  ⚑ decyzja
 │   4. Powiadomienie· raport wydania   ◌ oczekuje
 ├───────────────────────────────────────────────────
 │  KOMUNIKATY STERUJĄCE
 │   Koordynator → Wykonawca: zadanie 3, kontekst 2
 │   Wykonawca → Koordynator: krok nieodwracalny
 ├───────────────────────────────────────────────────
 │  KONTROLA JAKOŚCI
 │   zadanie 2 · wynik przyjęty
 ├───────────────────────────────────────────────────
 │  PRZEBIEG PĘTLI   ▮▮▮▮▮▯▯▯   2/4
 │  [Zatwierdź] [Wstrzymaj] [Wznów] [Przerwij]   ⋮
 └───────────────────────────────────────────────────
```

| Element okna | Co to jest | Do czego służy | Stany | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|
| Nagłówek zlecenia | Belka z treścią bieżącego zlecenia | Identyfikacja procesu prowadzonego przez pętlę | Domyślny · zlecenie skorygowane | 2 | Widoczny po otwarciu okna |
| Lista dekompozycji zlecenia | Wykaz zadań wywiedzionych ze zlecenia | Wgląd w podział pracy | Pusta · w toku · zakończona | 2 | Widoczna po otwarciu okna |
| Pozycja zadania | Wiersz z opisem, modułem docelowym i statusem | Śledzenie pojedynczego zadania pętli | Oczekuje · w toku · decyzja · zakończone · błąd · ponowione | 2 | Widoczna; kliknięcie otwiera szczegóły zadania |
| Znacznik punktu decyzyjnego `⚑` | Oznaczenie zadania wstrzymanego do decyzji Użytkownika | Wskazanie miejsca wymagającego zatwierdzenia | Oczekujący · rozstrzygnięty | 2 | Widoczny przy zadaniu wstrzymanym |
| Strumień komunikatów sterujących | Obszar wymiany Koordynator ↔ Wykonawca | Wgląd w treść koordynacji i nadzoru | Pusty · strumień aktywny | 2 | Widoczny po otwarciu okna |
| Panel kontroli jakości | Zestawienie ocen wyników zadań | Podstawa decyzji o przyjęciu albo ponowieniu | Brak ocen · wynik przyjęty · wynik odrzucony | 2 | Widoczny po otwarciu okna |
| Wskaźnik przebiegu pętli | Pasek postępu z licznikiem zadań | Orientacja w zaawansowaniu zlecenia | 0–100% · nieokreślony | 2 | Widoczny po otwarciu okna |
| Przycisk „Zatwierdź” | Kontrolka rozstrzygnięcia punktu decyzyjnego | Zatwierdzenie kroku wstrzymanego z urządzenia przenośnego | Domyślny · ładowanie · zakolejkowany lokalnie | 2 | Widoczny przy punkcie decyzyjnym |
| Przycisk „Wstrzymaj” | Kontrolka zatrzymania pętli | Wstrzymanie realizacji bez utraty stanu zadań | Domyślny · wstrzymane | 2 | Widoczny w pasku sterowania przebiegiem |
| Przycisk „Wznów” | Kontrolka podjęcia wstrzymanej pętli | Kontynuacja realizacji od zadania bieżącego | Domyślny · ładowanie | 2 | Widoczny w pasku sterowania przebiegiem |
| Przycisk „Przerwij” | Kontrolka zakończenia pętli | Przerwanie zadania i zamknięcie realizacji zlecenia | Domyślny · ładowanie · zakolejkowany lokalnie | 2 | Widoczny w pasku sterowania przebiegiem |
| Korekta zlecenia | Kontrolka zmiany treści zlecenia w toku | Zmiana zakresu pracy bez zakładania nowego zlecenia | Domyślna · edycja | 3 | Pozycja kebaba okna albo polecenie w Chat Window |
| Decyzja o ponowieniu | Kontrolka ponownego uruchomienia zadania z korektą | Naprawa wyniku odrzuconego w kontroli jakości | Domyślna · ładowanie | 3 | Kebab pozycji zadania albo panel kontroli jakości |
| Kebab okna `⋮` | Zbiorcze menu okna | Historia przebiegów, przypisanie Wykonawcy do zadań, szerokość kolumny, eksport przebiegu | Domyślny · rozwinięty | 3 | Kliknięcie ikony ⋮ w nagłówku |
| Dziennik audytu przebiegu | Pełny zapis zleceń, zadań, komunikatów i decyzji (Bezpieczeństwo, rozdz. 13.5) | Odtworzenie przebiegu decyzji po fakcie | Zamknięty · otwarty | 4 | Polecenie języka naturalnego albo wyszukiwarka funkcji |

### 4.6. Przegląd zadań i procesów

```
 ┌─ Przegląd zadań i procesów ────────────────────
 │  Zakres ▼   Rodzaj ▼   Stan ▼                ⋮
 ├────────────────────────────────────────────────
 │  ⟳  Budowa pakietu           CodeStudio   30%
 │     zadanie 3 z 4 · Koordynator
 │  ⚑  Publikacja raportu       Workspace  decyzja
 │     oczekuje 14 min
 │  ⟳  Analiza źródeł           TalkIn       62%
 │  ⏸  Nocna synchronizacja     Workspace  wstrzym.
 │  ✔  Przekład opisu           TalkIn     zakończ.
 ├────────────────────────────────────────────────
 │  ● połączony        ⟳ 3 w toku    ⚑ 2 decyzje
 └────────────────────────────────────────────────
```

| Element ekranu | Co to jest | Do czego służy | Stany | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|
| Lista procesów | Wykaz projektów, automatyzacji, agentów, aplikacji, sesji i przebiegów | Monitoring procesów niezależny od aktywnego środowiska | Pusta · zapełniona · filtrowana · odtworzona z pamięci lokalnej | 1 | Widoczna |
| Pozycja procesu | Wiersz z nazwą, środowiskiem pochodzenia, stanem i wskaźnikiem | Wejście do sterowania pojedynczym procesem | W toku · decyzja · wstrzymany · zakończony · błąd | 1 | Widoczna; kliknięcie otwiera szczegóły procesu |
| Wiersz podsumowania | Pas liczników: połączenie, procesy w toku, decyzje oczekujące | Orientacja bez przewijania listy | Domyślny | 1 | Widoczny |
| Filtr `Zakres ▼` | Element zbiorczy wyboru środowiska albo projektu | Zawężenie listy do jednego środowiska albo projektu | Domyślny · zawężony | 2 | Kliknięcie elementu zbiorczego |
| Filtr `Rodzaj ▼` | Element zbiorczy wyboru rodzaju procesu | Zawężenie do automatyzacji, agentów, sesji, przebiegów | Domyślny · zawężony | 2 | Kliknięcie elementu zbiorczego |
| Filtr `Stan ▼` | Element zbiorczy wyboru stanu | Zawężenie do procesów w toku, wstrzymanych, z decyzją, zakończonych | Domyślny · zawężony | 2 | Kliknięcie elementu zbiorczego |
| Zestaw akcji procesu `Operacje ▼` | Element zbiorczy akcji: uruchom, zatrzymaj, wstrzymaj, wznów, zatwierdź, modyfikuj | Zarządzanie zadaniem bez wchodzenia w szczegóły | Domyślny · rozwinięty · ładowanie | 3 | Kebab pozycji `⋮` albo przesunięcie poziome wiersza |
| Podgląd artefaktu wyniku | Widok treści, wersji i metadanych artefaktu | Ocena wyniku przed decyzją | Zamknięty · otwarty · niedostępny bez połączenia | 2 | Kliknięcie referencji wyniku w pozycji procesu |
| Kebab ekranu `⋮` | Zbiorcze menu ekranu | Kolejność listy, odświeżenie, zapis filtru jako domyślnego, eksport zestawienia | Domyślny · rozwinięty | 3 | Kliknięcie ikony ⋮ |
| Widok kolejek i zależności orkiestracji | Zestawienie kolejek, pozycji i zależności między zadaniami | Rozpoznanie przyczyny zastoju przebiegu | Zamknięty · otwarty | 4 | Polecenie języka naturalnego albo wyszukiwarka funkcji |

### 4.7. Powiadomienia

```
 ┌─ Powiadomienia ────────────────────────────────
 │  Klasa ▼                                     ⋮
 ├────────────────────────────────────────────────
 │  ⚑  Decyzja: publikacja raportu        14 min
 │     [Zatwierdź]  [Odrzuć]  [Otwórz pętlę]
 │  ✔  Zakończono: przekład opisu          31 min
 │     [Otwórz wynik]
 │  ⚠  Błąd zadania: budowa pakietu         1 godz
 │     [Ponów]  [Otwórz pętlę]
 │  ◈  Sparowano urządzenie: Tablet         2 godz
 │     [Otwórz rejestr urządzeń]
 ├────────────────────────────────────────────────
 │  Ustawienia powiadomień                     →
 └────────────────────────────────────────────────
```

| Element ekranu | Co to jest | Do czego służy | Stany | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|
| Lista powiadomień | Wykaz zdarzeń wymagających uwagi, uporządkowany od najnowszego | Wejście do interwencji bez szukania procesu | Pusta · zapełniona · filtrowana | 1 | Widoczna |
| Pozycja powiadomienia | Wiersz z klasą zdarzenia, treścią, czasem i zestawem działań | Rozstrzygnięcie zdarzenia z jednego miejsca | Nowa · przeczytana · rozstrzygnięta · wygasła | 1 | Widoczna |
| Działania pozycji | Przyciski właściwe klasie zdarzenia | Wykonanie działania bez otwierania kolejnego ekranu | Domyślne · ładowanie · zakolejkowane lokalnie | 1 | Widoczne w pozycji |
| Filtr `Klasa ▼` | Element zbiorczy wyboru klasy zdarzenia | Zawężenie listy do jednej klasy | Domyślny · zawężony | 2 | Kliknięcie elementu zbiorczego |
| Przejście „Ustawienia powiadomień” | Odnośnik do ustawień urządzenia | Zmiana zakresu klas dostarczanych na to urządzenie | Domyślny | 1 | Widoczne |
| Oznaczenie wszystkich jako przeczytane | Akcja zbiorcza listy | Uporządkowanie listy po serii zdarzeń | Domyślna · ładowanie | 3 | Pozycja kebaba ekranu |
| Kebab ekranu `⋮` | Zbiorcze menu ekranu | Oznaczenie jako przeczytane, retencja listy, zakres klas, eksport zestawienia | Domyślny · rozwinięty | 3 | Kliknięcie ikony ⋮ |
| Dziennik dostarczeń powiadomień | Zapis prób dostarczenia kanałem własnym i kanałem pośrednika | Diagnostyka braku powiadomienia na urządzeniu | Zamknięty · otwarty | 4 | Polecenie języka naturalnego albo wyszukiwarka funkcji |

### 4.8. Ustawienia urządzenia

```
 ┌─ Ustawienia urządzenia ────────────────────────
 │                                              ⋮
 │  URZĄDZENIE
 │   Nazwa:  Telefon Operatora            [Zmień]
 │   Typ:    telefon
 │   Sparowano: 2026-07-14
 │   Poświadczenie ważne do: 2026-10-12
 │
 │  POWIADOMIENIA
 │   Decyzje i zatwierdzenia            ( ● wł.)
 │   Zakończenie procesu                ( ● wł.)
 │   Błędy zadań                        ( ● wł.)
 │   Konto i dostęp                     ( ● wł.)
 │   Zmiany konfiguracji                ( ○ wył.)
 │
 │  PRACA BEZ POŁĄCZENIA
 │   Zakres lokalny                     ( ● wł.)
 │   Działania w kolejce: 0
 │
 │  [        Wyloguj to urządzenie        ]
 └────────────────────────────────────────────────
```

| Element ekranu | Co to jest | Do czego służy | Stany | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|
| Nazwa urządzenia | Pole tekstowe z akcją „Zmień” | Rozpoznanie urządzenia w rejestrze (rozdz. 9.6) | Domyślna · edycja · zapisano | 2 | Widoczna po otwarciu ekranu |
| Typ urządzenia | Etykieta wartości `telefon` albo `tablet` | Odwzorowanie pola `urzadzenie.typ` | Domyślny | 2 | Widoczny |
| Data parowania | Etykieta z datą wydania poświadczenia | Potwierdzenie momentu wejścia urządzenia do rejestru | Domyślna | 2 | Widoczna |
| Ważność poświadczenia | Etykieta z datą wygaśnięcia | Uprzedzenie o zbliżającym się wygaśnięciu (rozdz. 9.7) | Ważne · zbliża się wygaśnięcie · wygasłe | 2 | Widoczna |
| Przełączniki klas powiadomień | Zestaw kontrolek `.dn-suwak` | Wybór klas zdarzeń dostarczanych na to urządzenie | Włączony · wyłączony | 2 | Widoczne po otwarciu ekranu |
| Przełącznik zakresu lokalnego | Kontrolka `.dn-suwak` | Sterowanie utrzymywaniem zakresu dostępnego bez połączenia (rozdz. 7) | Włączony · wyłączony | 2 | Widoczny po otwarciu ekranu |
| Licznik działań w kolejce lokalnej | Etykieta z liczbą i wejściem do listy | Wgląd w działania oczekujące na wysłanie | Pusta · zapełniona · uzgadnianie | 2 | Widoczny po otwarciu ekranu |
| Lista działań kolejki lokalnej | Wykaz działań z czasem zapisu i stanem uzgodnienia | Rozstrzygnięcie konfliktu po odzyskaniu połączenia | Zamknięta · otwarta | 3 | Kliknięcie licznika działań |
| Przycisk „Wyloguj to urządzenie” | Przycisk główny akcji | Zakończenie dostępu tego urządzenia i usunięcie poświadczenia lokalnego (rozdz. 9.8) | Domyślny · ładowanie · potwierdzenie | 2 | Widoczny po otwarciu ekranu |
| Kebab ekranu `⋮` | Zbiorcze menu ekranu | Rejestr urządzeń konta, odświeżenie poświadczenia, diagnostyka połączenia, czyszczenie zakresu lokalnego | Domyślny · rozwinięty | 3 | Kliknięcie ikony ⋮ |
| Rejestr urządzeń konta | Wykaz wszystkich urządzeń konta z akcją odwołania dostępu | Odwołanie dostępu urządzenia utraconego (rozdz. 9.9) | Zamknięty · otwarty | 3 | Pozycja kebaba ekranu |
| Tryb administracyjny | Zestaw operacji administracyjnych funkcji | Odświeżenie poświadczenia, wymuszenie ponownego uzgodnienia stanu, podgląd wersji kontraktu | Nieaktywny · aktywny | 4 | Polecenie języka naturalnego, wyszukiwarka funkcji albo konfiguracja roli |

---

## 5. Kanały komunikacji operacyjnej na urządzeniu przenośnym

Oba kanały są na urządzeniu przenośnym elementami pierwszoplanowymi architektury, nie funkcjami dodatkowymi. Występują w tym samym miejscu układu co na stanowisku stacjonarnym: Chat Window w lewej kolumnie obszaru roboczego, Execution Loop Window w kolumnie z nim sąsiadującej.

### 5.1. Chat Window — kanał Użytkownik ↔ Wykonawca

Chat Window jest centralnym punktem pracy na urządzeniu przenośnym i podstawowym mechanizmem sterowania wszystkimi procesami realizowanymi przez platformę. Przyjmuje polecenia w języku naturalnym — wpisane albo podyktowane — prezentuje strumień odpowiedzi i wyników, umożliwia zatwierdzanie oraz przerywanie działań, a także wyjaśnianie wyniku i kontekstu. Polecenie wydane z urządzenia przenośnego uruchamia operacje we wszystkich środowiskach i modułach platformy dokładnie tak samo, jak polecenie wydane na stanowisku stacjonarnym.

| Zdolność kanału na urządzeniu przenośnym | Realizacja |
|---|---|
| Polecenie w języku naturalnym | Pole polecenia z wejściem tekstowym i głosowym; komunikat `message.send` (Kontrakty komunikacji, rozdz. 9.1) |
| Strumień odpowiedzi i wyników | Fragmenty strumienia z numeracją `seq` i znacznikiem `done` |
| Zatwierdzanie działania | Przycisk „Zatwierdź” przy kroku oczekującym na decyzję |
| Przerywanie działania | Przycisk „Przerwij” podczas pracy Wykonawcy |
| Wyjaśnianie wyniku i kontekstu | Polecenie „wyjaśnij ten wynik” kierowane do Wykonawcy w tym samym strumieniu |
| Sterowanie procesami wszystkich środowisk | Polecenie odnoszące się do procesu wskazanego nazwą albo referencją; zakres wskazuje znacznik środowiska w pasku kontekstu |
| Dostęp do funkcji warstwy 4 | Polecenie języka naturalnego jako podstawowa droga wywołania funkcji eksperckich na wąskim ekranie |

### 5.2. Execution Loop Window — kanał Koordynator ↔ Wykonawca

Execution Loop Window prezentuje komunikację między Koordynatorem — komponentem orkiestrującym platformy — a Wykonawcą. Na urządzeniu przenośnym odpowiada za nadzór nad pętlą wykonawczą prowadzoną po stronie serwera: prezentuje bieżące zlecenie i jego dekompozycję na zadania, kolejkę i stan zadań, wymianę komunikatów sterujących, wyniki kontroli jakości i decyzje o ponowieniu oraz wskaźniki przebiegu pętli, a także udostępnia sterowanie przebiegiem — wstrzymanie, wznowienie, przerwanie i korektę zlecenia. Okno otwiera się jako kolumna sąsiadująca z Chat Window.

| Zdolność nadzoru z urządzenia przenośnego | Realizacja |
|---|---|
| Podgląd dekompozycji zlecenia | Lista zadań z modułem docelowym, opisem i statusem |
| Zatwierdzanie punktu decyzyjnego | Przycisk „Zatwierdź” przy zadaniu oznaczonym znacznikiem `⚑`; komunikat `loop.control` z akcją zatwierdzenia |
| Odrzucenie wyniku i ponowienie zadania | Decyzja o ponowieniu z kebaba pozycji zadania albo z panelu kontroli jakości |
| Przerywanie zadania | Przycisk „Przerwij” w pasku sterowania przebiegiem |
| Wstrzymanie i wznowienie pętli | Przyciski „Wstrzymaj” i „Wznów”; stan zadań zachowany po stronie serwera |
| Korekta zlecenia w toku | Zmiana treści zlecenia bez zakładania nowego |
| Wgląd w komunikaty sterujące | Strumień wymiany Koordynator ↔ Wykonawca |

### 5.3. Punkty obowiązkowego zatwierdzenia

Punkty obowiązkowego zatwierdzenia przez Użytkownika (Bezpieczeństwo i uwierzytelnianie, rozdz. 13.4) są rozstrzygalne z urządzenia przenośnego w pełnym zakresie. Zatwierdzenie wydane na urządzeniu przenośnym jest tym samym zatwierdzeniem co wydane na stanowisku stacjonarnym — serwer rozgłasza rozstrzygnięcie do wszystkich urządzeń konta zdarzeniem właściwym zmienionemu obszarowi.

| Rodzaj punktu decyzyjnego | Prezentacja na urządzeniu przenośnym | Działania |
|---|---|---|
| Krok nieodwracalny zlecenia | Znacznik `⚑` w liście dekompozycji, powiadomienie klasy „Decyzje i zatwierdzenia” | Zatwierdź · Odrzuć · Otwórz pętlę |
| Negatywna kontrola jakości | Wpis panelu kontroli jakości z wynikiem odrzuconym | Ponów zadanie · Korekta zlecenia · Przerwij |
| Decyzja procesu MultitaskingAI | Pozycja sekcji „Decyzje oczekujące” strony głównej Mobile | Zatwierdź · Odrzuć · Otwórz Monitor procesu |
| Przekroczenie granicy autonomii Wykonawcy | Wpis strumienia Chat Window z prośbą o rozstrzygnięcie | Zatwierdź · Przerwij · Wyjaśnij |

---

## 6. Warstwy widoczności na urządzeniu przenośnym

### 6.1. Zasada nadrzędna

Jeżeli funkcja nie jest potrzebna do realizacji aktualnego zadania, nie jest widoczna. Interfejs urządzenia przenośnego ujawnia możliwości funkcji stopniowo — zależnie od kontekstu, roli użytkownika i wykonywanej czynności — zachowując maksymalną moc funkcjonalną przy minimalnej złożoności wizualnej. Warstwa 1 zajmuje ponad 80% powierzchni interfejsu; pozostałą część zajmują wyłącznie zwinięte wyzwalacze warstw 2–3. Ograniczona szerokość ekranu nie zmniejsza zakresu funkcji — zmienia wyłącznie proporcję między tym, co widoczne bez interakcji, a tym, co wywoływane jednym kliknięciem.

### 6.2. Cztery warstwy w funkcji Mobile

| Warstwa | Nazwa | Zawartość na urządzeniu przenośnym | Sposób wywołania |
|---|---|---|---|
| 1 | Zawsze widoczna | Chat Window, sekcja decyzji oczekujących, lista przebiegów w toku, lista procesów, lista powiadomień wraz z działaniami, wskaźnik połączenia, licznik kolejki lokalnej, przełącznik kolumn. Zajmuje ponad 80% powierzchni interfejsu | Widoczna bez interakcji |
| 2 | Widoczna na żądanie | Execution Loop Window, wybór środowiska, wybór modelu, filtry `Zakres ▼`, `Rodzaj ▼`, `Stan ▼`, ustawienia urządzenia, podgląd artefaktu, dyktowanie polecenia, panel pomocniczy | Znacznik kontekstowy, ikona, przełącznik, wyzwalacz `▼`; po użyciu element zwija się samoczynnie |
| 3 | Rozwinięcia kontekstowe | Zestaw akcji procesu `Operacje ▼`, kebab okna `⋮`, kebab pozycji `⋮`, menu nawigacji funkcji `☰`, rejestr urządzeń konta, lista działań kolejki lokalnej, korekta zlecenia, decyzja o ponowieniu, akcje zbiorcze listy powiadomień | Menu kebab (⋮), menu hamburger (☰), menu kontekstowe, panel popover, panel wysuwany, lista rozwijana |
| 4 | Funkcje eksperckie | Tryb administracyjny funkcji, dziennik audytu przebiegu, dziennik dostarczeń powiadomień, widok kolejek i zależności orkiestracji, diagnostyka połączenia, podgląd wersji kontraktu, odświeżenie poświadczenia | Polecenie języka naturalnego w Chat Window, wyszukiwarka funkcji, tryb administracyjny, konfiguracja roli. Użytkownik podstawowy nie widzi tych elementów |

### 6.3. Mechanizmy ukrywania funkcjonalności

- **Menu progresywne.** Wybór środowiska, modelu, zakresu przekroju i klasy powiadomień prezentowany jest jako jeden element zwinięty (`Zakres ▼`, `Rodzaj ▼`, `Stan ▼`, `Klasa ▼`); lista pozycji rozwija się po kliknięciu.
- **Panele wysuwane.** Podgląd artefaktu, lista działań kolejki lokalnej, rejestr urządzeń konta i szczegóły zadania są kolumnami bocznymi otwieranymi jako rozszerzenia boczne; po zamknięciu panel znika całkowicie z przestrzeni roboczej.
- **Grupowanie logiczne akcji.** Uruchomienie, zatrzymanie, wstrzymanie, wznowienie, zatwierdzenie i modyfikacja procesu występują jako jeden element zbiorczy `Operacje ▼`, którego rozwinięcie zawiera pełną listę akcji.
- **Znaczniki kontekstowe.** Zakres globalny, środowisko, projekt i model występują jako lekkie znaczniki w pasku kontekstu Chat Window, na przykład `[Danaco Console] [Multitasking] [Wydanie 4] [Fable 5]`; kliknięcie znacznika otwiera odpowiedni selektor.

### 6.4. Zasada jednego kliknięcia

Każda ukryta funkcja urządzenia przenośnego jest osiągalna jednym kliknięciem, jednym gestem, jednym skrótem klawiatury zewnętrznej albo jednym poleceniem języka naturalnego wydanym w Chat Window. Zagnieżdżanie funkcji głęboko w hierarchii menu jest zabronione. Na wąskim ekranie podstawową drogą dostępu do funkcji warstwy 4 jest polecenie języka naturalnego oraz wyszukiwarka funkcji.

| Funkcja ukryta | Droga jednego kliknięcia |
|---|---|
| Execution Loop Window | Wyzwalacz „Execution Loop ▼” w nagłówku Chat Window |
| Zestaw akcji procesu | Kebab `⋮` pozycji procesu albo przesunięcie poziome wiersza |
| Nawigacja funkcji | Wyzwalacz `☰` w nagłówku Chat Window |
| Rejestr urządzeń konta | Pozycja kebaba ekranu ustawień urządzenia |
| Dziennik audytu przebiegu | Polecenie „pokaż dziennik audytu tego przebiegu” |
| Tryb administracyjny funkcji | Polecenie „włącz tryb administracyjny” albo wyszukiwarka funkcji |

---

## 7. Praca bez połączenia

### 7.1. Zakres dostępny lokalnie

Serwer pozostaje jedynym źródłem prawdy. Zakres lokalny jest migawką odczytową ostatniego stanu potwierdzonego przez serwer, utrzymywaną na urządzeniu przenośnym po to, by brak łączności nie odbierał orientacji w stanie pracy.

| Zakres | Dostępność bez połączenia |
|---|---|
| Strumień Chat Window bieżącego kontekstu | Odczyt migawki ostatniej wymiany |
| Lista dekompozycji zlecenia i stan zadań | Odczyt migawki ostatniego stanu pętli |
| Lista procesów i przebiegów | Odczyt migawki z czasem pobrania |
| Lista powiadomień | Odczyt pełny; działania kolejkowane |
| Ustawienia urządzenia | Odczyt i zapis lokalny; zapis kolejkowany do uzgodnienia |
| Podgląd artefaktu | Niedostępny — treść artefaktu pobierana jest z serwera |
| Rejestr urządzeń konta | Niedostępny — rejestr odczytywany jest z serwera |
| Wydanie polecenia w Chat Window | Dostępne; polecenie trafia do kolejki działań lokalnych |
| Sterowanie procesem i zatwierdzanie decyzji | Dostępne; działanie trafia do kolejki działań lokalnych |

Każdy element prezentowany z migawki opatrzony jest znacznikiem czasu pobrania oraz plakietką „bez połączenia”. Żaden element nie traci klikalności — działanie wykonane bez łączności jest przyjmowane i kolejkowane, a jego stan sygnalizuje plakietka „zakolejkowane”.

### 7.2. Kolejkowanie działań

Działanie wykonane bez połączenia zapisywane jest w kolejce działań lokalnych (encja `dzialanie_lokalne`, [Model danych](../architektura/model-danych.md), rozdz. 3.6) wraz z czasem zapisu, rodzajem działania, identyfikatorem przedmiotu działania i wersją stanu, na którym Użytkownik podejmował decyzję.

```
Działanie Użytkownika bez łączności
        │
        ▼
Zapis do kolejki działań lokalnych  (rodzaj · przedmiot · wersja stanu · czas)
        │
        ▼
Plakietka „zakolejkowane” przy elemencie interfejsu
        │
        ▼
Odzyskanie połączenia ──► uzgadnianie stanu (rozdz. 7.3)
```

| Rodzaj działania | Kolejkowanie | Zachowanie po odzyskaniu połączenia |
|---|---|---|
| Polecenie w Chat Window | Tak | Wysłanie w kolejności zapisu; odpowiedź trafia do strumienia |
| Zatwierdzenie punktu decyzyjnego | Tak | Wysłanie z wersją stanu; serwer rozstrzyga zgodność |
| Przerwanie zadania | Tak | Wysłanie z wersją stanu; serwer rozstrzyga zgodność |
| Wstrzymanie i wznowienie pętli | Tak | Wysłanie z wersją stanu; serwer rozstrzyga zgodność |
| Zmiana ustawień urządzenia | Tak | Zapis ustawień urządzenia po stronie serwera |
| Oznaczenie powiadomienia jako przeczytanego | Tak | Zapis zbiorczy |
| Odczyt artefaktu, rejestru urządzeń, kolejek orkiestracji | Nie | Pobranie z serwera po odzyskaniu połączenia |

### 7.3. Uzgadnianie stanu po odzyskaniu połączenia

```
Urządzenie przenośne                 Serwer
      │  otwarcie kanału WebSocket     │
      │ ──────────────────────────────►│
      │  connection.hello              │
      │  (deviceId, poświadczenie)     │
      │ ──────────────────────────────►│
      │  response (authenticated)      │
      │ ◄──────────────────────────────│
      │  mobile.queue.sync             │
      │  (lista działań z wersjami)    │
      │ ──────────────────────────────►│
      │                                │  weryfikacja wersji
      │                                │  każdego działania
      │  response (przyjęte,           │
      │  odrzucone, do rozstrzygnięcia)│
      │ ◄──────────────────────────────│
      │  mobile.process.changed        │
      │  (stan po zastosowaniu)        │
      │ ◄──────────────────────────────│
```

| Krok uzgadniania | Przebieg |
|---|---|
| 1. Uwierzytelnienie urządzenia | Urządzenie przedstawia poświadczenie w komunikacie `connection.hello` (rozdz. 9.5) |
| 2. Przekazanie kolejki | Urządzenie wysyła `mobile.queue.sync` **[DO DECYZJI OPERATORA]** — nazwa nie występuje w `budowa/shared/contract.json`; do rozstrzygnięcia, czy komenda ma powstać z listą działań w kolejności zapisu, każde z wersją stanu |
| 3. Weryfikacja wersji | Serwer porównuje wersję stanu z wersją bieżącą przedmiotu działania |
| 4. Zastosowanie działań zgodnych | Działania o zgodnej wersji stosowane są w kolejności zapisu |
| 5. Rozstrzygnięcie działań niezgodnych | Działania o wersji nieaktualnej trafiają do rozstrzygnięcia (rozdz. 7.4) |
| 6. Zastąpienie migawki | Serwer przekazuje stan bieżący; migawka lokalna zostaje zastąpiona w całości |

**Pokrycie kontraktem.** `mobile.queue.sync` jest nazwą ustaloną w [Kontraktach komunikacji](../architektura/kontrakty-komunikacji.md) rozdz. 16.1 dla kanału Mobile; w `budowa/shared/contract.json` obszar `mobile` niesie obecnie wyłącznie `mobile.status.get`, `mobile.process.list`, `mobile.process.control` oraz zdarzenie `mobile.process.changed` — komenda uzgadniania kolejki działań lokalnych nie ma jeszcze odpowiednika w kontrakcie wdrożeniowym. **[DO DECYZJI OPERATORA]** — czy `mobile.queue.sync` wchodzi do `contract.json` jako nowa komenda obszaru `mobile`, czy uzgadnianie kolejki lokalnej odbywa się istniejącymi komendami obszarów `session` i `window` powtórzonymi w pętli po stronie klienta.

### 7.4. Rozstrzyganie konfliktów

Ponieważ serwer jest jedynym źródłem prawdy, nie występuje scalanie zmian po stronie klienta — każdy stan zatwierdzony przez serwer zastępuje poprzedni obraz na wszystkich urządzeniach (Kontrakty komunikacji, rozdz. 18). Konflikt oznacza wyłącznie to, że przedmiot działania zmienił stan między zapisem działania a jego wysłaniem.

| Rodzaj konfliktu | Rozstrzygnięcie | Prezentacja na urządzeniu |
|---|---|---|
| Punkt decyzyjny rozstrzygnięty w międzyczasie z innego urządzenia | Działanie zakolejkowane zostaje odrzucone jako bezprzedmiotowe | Wpis w strumieniu Chat Window: rozstrzygnięcie i jego źródło; plakietka „bezprzedmiotowe” |
| Zadanie zakończone przed przyjęciem przerwania | Przerwanie zostaje odrzucone; zadanie zachowuje wynik | Wpis w strumieniu z wynikiem zadania; plakietka „bezprzedmiotowe” |
| Pętla przerwana przed przyjęciem wznowienia | Wznowienie zostaje odrzucone | Komunikat kontekstowy przy pozycji przebiegu |
| Zlecenie skorygowane po zapisie polecenia | Polecenie zostaje przekazane Wykonawcy w kontekście zlecenia bieżącego | Wpis w strumieniu z aktualną treścią zlecenia |
| Ustawienie urządzenia zmienione z innego urządzenia | Zapis lokalny zostaje zastosowany jako późniejszy | Dymek powiadomienia potwierdzającego zapis |
| Poświadczenie odwołane w czasie pracy bez połączenia | Cała kolejka zostaje odrzucona, zakres lokalny usunięty | Ekran logowania z komunikatem „dostęp urządzenia odwołany” |

Działania odrzucone pozostają widoczne na liście działań kolejki lokalnej wraz z przyczyną odrzucenia do chwili potwierdzenia ich przez Użytkownika.

---

## 8. Powiadomienia wypychane

### 8.1. Klasy zdarzeń

Klasy odpowiadają kategoriom zdarzeń ustalonym w sekcji „Powiadomienia” okna Ustawień ([Okno Ustawień](../interfejs-uzytkownika/ustawienia.md), rozdz. 7.2). Funkcja nie wprowadza nowych zdarzeń — dostarcza na urządzenie przenośne zdarzenia zdefiniowane w kontraktach komunikacji poszczególnych obszarów platformy.

> **Zapis kluczy nastaw.** Nazwy w postaci `obszar.grupa.nastawa` użyte w tym rozdziale są
> **kluczami konfiguracji**, nie komendami kontraktu. Klucz wskazuje miejsce wartości w modelu
> konfiguracji; komendy kontraktu, którymi się go odczytuje i zapisuje, to `config.get` i
> `config.set` (obszar `config` w `budowa/shared/contract.json`).


| Klasa zdarzenia | Zdarzenie źródłowe | Domyślnie | Priorytet dostarczenia |
|---|---|---|---|
| Decyzje i zatwierdzenia | `message.approve` (kanał pierwszy), `loop.quality.result` z werdyktem `rejected` (kanał drugi) — [Kontrakty komunikacji](../architektura/kontrakty-komunikacji.md) rozdz. 9.1, 9.2 | aktywne | Najwyższy — dostarczane niezwłocznie |
| Zakończenie procesu | `loop.order.close`, `automation.execution.status`, `mobile.process.changed` | aktywne | Wysoki |
| Błędy zadań | `loop.task.result` ze statusem `failed`, `channel_unavailable` w kontekście zadania | aktywne | Wysoki |
| Zdarzenia MultitaskingAI | Monitor procesu, panel orkiestracji (Koncepcja platformy, rozdz. 13.8) | aktywne | Wysoki |
| Konto i dostęp | `auth.changed`, `device.changed`, `device.pair.completed` | aktywne | Wysoki |
| Zmiany konfiguracji | `config.changed`, `memory.changed`, `history.changed` | nieaktywne | Zwykły |

### 8.2. Treść powiadomienia

| Element treści | Zawartość |
|---|---|
| Klasa | Ikona i nazwa klasy zdarzenia |
| Tytuł | Nazwa przedmiotu zdarzenia: zlecenie, zadanie, automatyzacja, przebieg, urządzenie |
| Treść | Jedno zdanie opisujące zdarzenie i jego skutek dla pracy |
| Źródło | Środowisko albo projekt pochodzenia zdarzenia |
| Czas | Znacznik czasu wystąpienia zdarzenia |
| Działania | Zestaw działań właściwy klasie (rozdz. 8.3) |

Powiadomienie klasy „Decyzje i zatwierdzenia” nie zawiera treści przedmiotu decyzji, jeżeli kanał dostarczenia przebiega przez pośrednika — treść pobierana jest po otwarciu aplikacji (rozdz. 8.5).

### 8.3. Działania z poziomu powiadomienia

| Klasa zdarzenia | Działania dostępne bez otwierania aplikacji | Działania po otwarciu aplikacji |
|---|---|---|
| Decyzje i zatwierdzenia | Zatwierdź · Odrzuć | Otwórz Execution Loop Window na punkcie decyzyjnym |
| Zakończenie procesu | Oznacz jako przeczytane | Otwórz wynik · Otwórz przegląd procesów |
| Błędy zadań | Ponów zadanie | Otwórz Execution Loop Window · Otwórz dziennik przebiegu |
| Zdarzenia MultitaskingAI | Zatwierdź · Wstrzymaj | Otwórz Monitor procesu |
| Konto i dostęp | Oznacz jako przeczytane | Otwórz rejestr urządzeń konta |
| Zmiany konfiguracji | Oznacz jako przeczytane | Otwórz ustawienia urządzenia |

Działanie wykonane z poziomu powiadomienia bez łączności trafia do kolejki działań lokalnych na tych samych zasadach co działanie wykonane w aplikacji (rozdz. 7.2).

**Pokrycie kontraktem.** Ładunek `mobile.notification.action`, którym system operacyjny urządzenia przenośnego przekazuje działanie wykonane bezpośrednio na powiadomieniu wypychanym (bez otwierania aplikacji), jest nazwą ustaloną w [Kontraktach komunikacji](../architektura/kontrakty-komunikacji.md) rozdz. 16.1; w `budowa/shared/contract.json` obszar `mobile` nie niesie jeszcze tej komendy. **[DO DECYZJI OPERATORA]** — czy `mobile.notification.action` wchodzi do `contract.json` jako komenda obszaru `mobile`, czy działanie z poziomu powiadomienia przekłada się po stronie klienta wprost na komendę właściwą klasie zdarzenia (`notification.acknowledge`, `notification.resolve` albo komendę obszaru źródłowego zdarzenia, np. `automation.execution.status`).

### 8.4. Relacja do ustawień powiadomień

Zakres klas dostarczanych na urządzenie przenośne wynika z dwóch miejsc sterowania i jest ich iloczynem — klasa jest dostarczana wtedy, gdy jest aktywna w obu.

| Miejsce sterowania | Zakres | Zawartość |
|---|---|---|
| Okno Ustawień, sekcja „Powiadomienia” ([Okno Ustawień](../interfejs-uzytkownika/ustawienia.md), rozdz. 7) | Konto | Które kategorie zdarzeń i którym kanałem są dostarczane; kanał „Mobile — powiadomienie push” |
| Ekran ustawień urządzenia (rozdz. 4.8) | To urządzenie | Które klasy zdarzeń dostarczane są na ten konkretny egzemplarz |

### 8.5. Kanały dostarczenia

| Kanał | Przebieg | Zakres treści |
|---|---|---|
| Kanał własny | Dostarczenie otwartym kanałem WebSocket, gdy aplikacja pracuje na pierwszym planie | Pełna treść zdarzenia wraz z działaniami |
| Kanał pośrednika | Dostarczenie usługą powiadomień systemu urządzenia, gdy aplikacja nie pracuje na pierwszym planie | Klasa, tytuł i czas; treść przedmiotu pobierana po otwarciu aplikacji |

Powiadomienie dostarczone kanałem pośrednika i rozstrzygnięte z jego poziomu wywołuje nawiązanie połączenia, przedstawienie poświadczenia urządzenia i wysłanie działania; do chwili potwierdzenia przez serwer działanie pozostaje w kolejce działań lokalnych.

---

## 9. Protokół parowania urządzenia

### 9.1. Przeznaczenie protokołu

Parowanie jest drogą wejścia urządzenia przenośnego do funkcji globalnej Mobile. Zakładane jest w obecności Operatora, przy maszynie już uwierzytelnionej, i wydaje nowemu urządzeniu poświadczenie odrębne od poświadczeń pozostałych urządzeń — odwołanie jednego urządzenia nie odcina pozostałych. Protokół realizuje model właściciela „jedno konto, wiele urządzeń” ([Bezpieczeństwo i uwierzytelnianie](../architektura/bezpieczenstwo-i-uwierzytelnianie.md), rozdz. 3) i kończy się wpisem w rejestrze urządzeń konta.

| Właściwość | Ustalenie |
|---|---|
| Inicjator | Operator przy urządzeniu stacjonarnym z ważnym tokenem dostępu |
| Miejsce inicjowania | Okno Ustawień, sekcja „Urządzenia”, działanie „Sparuj urządzenie” ([Okno Ustawień](../interfejs-uzytkownika/ustawienia.md), rozdz. 6.2) |
| Przedmiot wydania | Poświadczenie urządzenia — odrębne dla każdego urządzenia, przechowywane poza bazą danych |
| Nośnik protokołu | Ten sam, jedyny kanał WebSocket co pozostała komunikacja platformy |
| Zakończenie pozytywne | Wpis urządzenia w rejestrze urządzeń konta i uaktywnienie funkcji Mobile na tym urządzeniu |
| Zakończenie negatywne | Unieważnienie kodu parowania bez zmiany w rejestrze |

### 9.2. Sekwencja kroków

| Krok | Miejsce | Przebieg |
|---|---|---|
| 1. Inicjowanie parowania | Stanowisko stacjonarne, okno Ustawień → Urządzenia | Operator wybiera działanie „Sparuj urządzenie”; klient wysyła `device.pair.start` |
| 2. Wydanie kodu parowania | Serwer | Serwer tworzy kod parowania o ograniczonej ważności, wiąże go z kontem i zwraca w odpowiedzi wraz z czasem wygaśnięcia |
| 3. Prezentacja kodu | Stanowisko stacjonarne | Okno nakładkowe prezentuje kod w postaci ciągu znaków oraz kodu graficznego; obok kodu widoczny jest licznik pozostałej ważności |
| 4. Podjęcie kodu | Urządzenie przenośne | Urządzenie odczytuje kod graficzny albo Operator wpisuje kod ręcznie; klient wysyła `device.pair.confirm` z kodem, identyfikatorem urządzenia, typem i nazwą wstępną |
| 5. Weryfikacja | Serwer | Serwer sprawdza istnienie kodu, jego ważność i to, że nie został wcześniej użyty |
| 6. Potwierdzenie przy maszynie | Stanowisko stacjonarne | Serwer wysyła zdarzenie z danymi urządzenia zgłaszającego; Operator potwierdza wpuszczenie urządzenia do rejestru |
| 7. Wydanie poświadczenia | Serwer → urządzenie przenośne | Serwer wydaje poświadczenie urządzenia, unieważnia kod parowania i zapisuje wpis w rejestrze urządzeń |
| 8. Przechowanie poświadczenia | Urządzenie przenośne | Poświadczenie zapisywane jest w magazynie danych dostępowych urządzenia; baza przechowuje wyłącznie odwołanie |
| 9. Zamknięcie protokołu | Oba urządzenia | Stanowisko stacjonarne prezentuje nowy wiersz w tabeli urządzeń i dymek powiadomienia potwierdzającego; urządzenie przenośne przechodzi do strony głównej Mobile |

### 9.3. Diagram protokołu

```
Stanowisko stacjonarne          Serwer              Urządzenie przenośne
   (token ważny)          (źródło prawdy)          (bez poświadczenia)
        │                        │                          │
        │ device.pair.start      │                          │   Krok 1
        │ ──────────────────────►│                          │
        │                        │  utworzenie kodu         │   Krok 2
        │  response              │  parowania i terminu     │
        │  (kod, wygasa)         │  ważności                │
        │ ◄──────────────────────│                          │
        │                        │                          │
   [ prezentacja kodu ]          │                          │   Krok 3
   [ licznik ważności ]          │                          │
        │                        │   device.pair.confirm    │   Krok 4
        │                        │   (kod, deviceId, typ,   │
        │                        │    nazwa wstępna)        │
        │                        │ ◄────────────────────────│
        │                        │                          │
        │                        │  weryfikacja kodu:       │   Krok 5
        │                        │  istnieje · ważny ·      │
        │                        │  nieużyty                │
        │  device.pair.request   │                          │   Krok 6
        │  (typ, nazwa, czas)    │                          │
        │ ◄──────────────────────│                          │
   [ potwierdzenie Operatora ]   │                          │
        │  device.pair.approve   │                          │
        │ ──────────────────────►│                          │
        │                        │  wydanie poświadczenia,  │   Krok 7
        │                        │  unieważnienie kodu,     │
        │                        │  wpis w rejestrze        │
        │                        │  device.pair.completed   │
        │                        │ ────────────────────────►│
        │                        │                          │
        │                        │  [ zapis poświadczenia   │   Krok 8
        │                        │    w magazynie ]         │
        │  device.changed        │                          │   Krok 9
        │ ◄──────────────────────│                          │
   [ nowy wiersz tabeli ]        │      [ strona główna Mobile ]
```

Legenda: trzy pionowe linie życia to — od lewej — stanowisko stacjonarne, serwer
i urządzenie przenośne; strzałka wskazuje kierunek komunikatu, a jej opis podaje
nazwę komunikatu kontraktu wraz z polami istotnymi dla kroku. Oznaczenia „Krok 1”
do „Krok 9” stoją poza prawą linią życia i wskazują numer kroku z tabeli rozdziału
9.2, do którego odnosi się wiersz albo blok wierszy stojący na jego wysokości.
Nawiasy kwadratowe oznaczają czynność wykonywaną w obrębie jednej linii życia,
bez wymiany komunikatu.

### 9.4. Kod parowania — prezentacja i ważność

| Właściwość kodu parowania | Ustalenie |
|---|---|
| Postać | Ciąg znaków do wpisania ręcznego oraz kod graficzny do odczytu kamerą urządzenia przenośnego |
| Zasięg | Jedno konto, jedno użycie |
| Ważność | Ograniczona czasem; licznik pozostałej ważności widoczny obok kodu |
| Unieważnienie przed użyciem | Zamknięcie okna nakładkowego, wybranie działania „Unieważnij kod”, wydanie nowego kodu, upływ ważności |
| Unieważnienie po użyciu | Natychmiastowe, w kroku 7 — kod nie podlega ponownemu użyciu |
| Przechowywanie | Poza bazą danych, w magazynie danych dostępowych; baza przechowuje odwołanie i termin ważności ([Model danych](../architektura/model-danych.md), rozdz. 3.5) |
| Prezentacja w interfejsie | Okno nakładkowe sekcji „Urządzenia” okna Ustawień; kod nigdy nie trafia do strumienia Chat Window ani do dziennika audytu w postaci jawnej |

### 9.5. Wymiana i przechowywanie poświadczeń

| Właściwość poświadczenia urządzenia | Ustalenie |
|---|---|
| Moment wydania | Krok 7 protokołu parowania; ponadto przy odświeżeniu poświadczenia (rozdz. 9.7) |
| Przypisanie | Do encji `urzadzenie`, nie do pojedynczego połączenia |
| Odrębność | Każde urządzenie ma własne poświadczenie; odwołanie jednego nie odcina pozostałych |
| Trwałość | Ważne między kolejnymi połączeniami tego samego urządzenia, aż do wygaśnięcia albo odwołania |
| Przechowywanie po stronie serwera | Poza bazą danych, w magazynie danych dostępowych; baza przechowuje odwołanie ([Model danych](../architektura/model-danych.md), rozdz. 1.5) |
| Przechowywanie po stronie urządzenia | W magazynie danych dostępowych urządzenia; nigdy w postaci jawnej w pamięci aplikacji ani w migawce zakresu lokalnego |
| Przedstawienie | Pole `token` komunikatu `connection.hello` przy każdym nawiązaniu połączenia (Kontrakty komunikacji, rozdz. 2, krok uwierzytelnienia) |
| Relacja do tokenu urządzenia stacjonarnego | Ten sam mechanizm — token urządzenia opisany w [Bezpieczeństwie i uwierzytelnianiu](../architektura/bezpieczenstwo-i-uwierzytelnianie.md), rozdz. 8.1; parowanie jest drogą jego wydania bez wprowadzania loginu i hasła na urządzeniu przenośnym |

### 9.6. Rejestr sparowanych urządzeń

Rejestr urządzeń konta jest wykazem encji `urzadzenie` powiązanych z jedynym kontem właściciela. Prezentowany jest w sekcji „Urządzenia” okna Ustawień na stanowisku stacjonarnym oraz — w postaci odczytowej z akcją odwołania — na ekranie ustawień urządzenia przenośnego.

| Pole rejestru | Zawartość | Źródło |
|---|---|---|
| Identyfikator urządzenia | Skrócony identyfikator techniczny | `urzadzenie.id` |
| Nazwa | Nazwa własna urządzenia, nadana przy parowaniu i zmienialna | `urzadzenie.nazwa` |
| Typ | `komputer` \| `telefon` \| `tablet` | `urzadzenie.typ` |
| Tryb Mobile | Czy urządzenie korzysta z funkcji globalnej Mobile | `urzadzenie.tryb_mobile_aktywny` |
| Odwołanie do poświadczenia | Wskazanie wpisu magazynu danych dostępowych | `poswiadczenie_urzadzenia.wartosc_odwolanie` |
| Data parowania | Data wydania poświadczenia | `poswiadczenie_urzadzenia.data_wydania` |
| Termin ważności poświadczenia | Data wygaśnięcia poświadczenia | `poswiadczenie_urzadzenia.data_waznosci` |
| Stan poświadczenia | `aktywne` \| `wygasle` \| `odwolane` | `poswiadczenie_urzadzenia.stan` |
| Ostatnie połączenie | Data ostatniej aktywności urządzenia | `urzadzenie.data_ostatniego_polaczenia` |
| Stan bieżący | Plakietka „to urządzenie” dla urządzenia bieżącego, „aktywne” dla pozostałych | Wyliczane z połączenia |

Nazwa urządzenia nadawana jest w kroku 4 protokołu jako nazwa wstępna zaproponowana przez urządzenie zgłaszające, potwierdzana przez Operatora w kroku 6 i zmienialna z dwóch miejsc: z sekcji „Urządzenia” okna Ustawień oraz z ekranu ustawień urządzenia przenośnego. Rozpoznanie urządzenia w interfejsie opiera się na parze nazwa i typ; identyfikator techniczny prezentowany jest w postaci skróconej jako rozstrzygnięcie przy nazwach powtarzalnych.

### 9.7. Wygaśnięcie poświadczenia

| Zdarzenie | Zachowanie platformy |
|---|---|
| Zbliżający się termin ważności | Ekran ustawień urządzenia prezentuje stan „zbliża się wygaśnięcie”; klient wysyła `auth.token.refresh` przy najbliższym połączeniu |
| Odświeżenie poświadczenia | Serwer wydaje poświadczenie o nowym terminie ważności, unieważnia poprzednie, aktualizuje wpis rejestru; przebieg nie wymaga udziału Operatora |
| Upływ terminu ważności | Poświadczenie przyjmuje stan `wygasle`; połączenie nie przechodzi etapu uwierzytelnienia |
| Praca po wygaśnięciu | Urządzenie przechodzi do ekranu logowania z komunikatem „poświadczenie urządzenia wygasło”; zakres lokalny pozostaje dostępny odczytowo, kolejka działań lokalnych zostaje zachowana do czasu ponownego uwierzytelnienia |
| Ponowne uzyskanie dostępu | Logowanie loginem i hasłem albo metodą dodatkową; ponowne parowanie nie jest wymagane, wpis rejestru zachowuje ciągłość |

### 9.8. Wylogowanie zdalne i wylogowanie z urządzenia

| Rodzaj wylogowania | Inicjator | Przebieg | Skutek |
|---|---|---|---|
| Wylogowanie z urządzenia | Użytkownik przy urządzeniu przenośnym | Działanie „Wyloguj to urządzenie” na ekranie ustawień urządzenia; komunikat `auth.logout` | Sesja bramki unieważniona, zakres lokalny i kolejka działań lokalnych usunięte, zdarzenie `auth.changed` (powód `sessionRevoked`) odświeża sekcję Urządzenia na pozostałych urządzeniach |
| Wylogowanie zdalne | Operator z dowolnego urządzenia konta | Akcja „Odłącz” w rejestrze urządzeń; komunikat `device.revoke` | Poświadczenie wskazanego urządzenia unieważnione; otwarte połączenie tego urządzenia zamykane niezwłocznie |
| Wylogowanie wszystkich urządzeń | Operator | Akcja zbiorcza rejestru urządzeń; odzyskanie konta ([Bezpieczeństwo i uwierzytelnianie](../architektura/bezpieczenstwo-i-uwierzytelnianie.md), rozdz. 7) | Unieważnienie wszystkich poświadczeń wydanych przed rozstrzygnięciem; ponowne logowanie na wszystkich urządzeniach |

Urządzenie, którego poświadczenie zostało odwołane, przy najbliższej próbie połączenia otrzymuje odmowę uwierzytelnienia, usuwa zakres lokalny i kolejkę działań lokalnych oraz przechodzi do ekranu logowania z komunikatem „dostęp urządzenia odwołany”.

### 9.9. Postępowanie przy utracie urządzenia

| Krok | Miejsce | Działanie |
|---|---|---|
| 1 | Dowolne urządzenie konta z ważnym poświadczeniem | Otwarcie rejestru urządzeń: sekcja „Urządzenia” okna Ustawień albo rejestr urządzeń konta na urządzeniu przenośnym |
| 2 | Rejestr urządzeń | Rozpoznanie urządzenia utraconego po nazwie, typie i dacie ostatniego połączenia |
| 3 | Rejestr urządzeń | Akcja „Odłącz” — unieważnienie poświadczenia urządzenia utraconego (`device.revoke`) |
| 4 | Serwer | Zamknięcie otwartych połączeń urządzenia, odrzucenie działań kolejki lokalnej nadesłanych po odwołaniu |
| 5 | Wszystkie urządzenia konta | Zdarzenie `device.changed` aktualizuje rejestr; powiadomienie klasy „Konto i dostęp” |
| 6 | Dziennik audytu | Zapis odwołania dostępu wraz z urządzeniem inicjującym i czasem |
| 7 | Konto | Przy podejrzeniu ujawnienia materiału uwierzytelniającego — odzyskanie konta, unieważniające wszystkie poświadczenia wydane wcześniej |

Utrata urządzenia nie wymaga zmiany hasła konta, jeżeli materiał uwierzytelniający konta nie został ujawniony — poświadczenie urządzenia jest odrębne i jego odwołanie zamyka dostęp tego jednego egzemplarza.

### 9.10. Tabela komunikatów protokołu

| Typ komunikatu | Kierunek | Opis |
|---|---|---|
| `device.pair.start` | polecenie | Inicjowanie parowania z sekcji „Urządzenia” okna Ustawień; odpowiedź niesie kod parowania i termin jego ważności. |
| `device.pair.cancel` | polecenie | Unieważnienie wydanego kodu parowania przed jego użyciem. |
| `device.pair.confirm` | polecenie | Zgłoszenie urządzenia przenośnego z kodem parowania, identyfikatorem urządzenia, typem i nazwą wstępną. |
| `device.pair.request` | zdarzenie | Zgłoszenie urządzenia przekazane do potwierdzenia na maszynie inicjującej parowanie. |
| `device.pair.approve` | polecenie | Potwierdzenie wpuszczenia urządzenia do rejestru przez Operatora. |
| `device.pair.reject` | polecenie | Odmowa wpuszczenia urządzenia; kod parowania zostaje unieważniony. |
| `device.pair.completed` | zdarzenie | Wydanie poświadczenia urządzeniu przenośnemu i zamknięcie protokołu. |
| `auth.token.refresh` | polecenie | Przedłużenie sesji bramki dla wskazanego urządzenia (`deviceId`); wydaje sesję (`AuthSession`) o nowym terminie ważności w miejsce sesji zbliżającej się do wygaśnięcia. |
| `device.list` | polecenie | Odczyt rejestru urządzeń konta. |
| `device.rename` | polecenie | Zmiana nazwy urządzenia w rejestrze. |
| `device.revoke` | polecenie | Odwołanie poświadczenia wskazanego urządzenia — wylogowanie zdalne. |
| `auth.logout` | polecenie | Unieważnienie sesji bramki tego urządzenia — pozycja „Wyloguj to urządzenie” ekranu ustawień urządzenia; skutek rozgłaszany zdarzeniem `auth.changed` (powód `sessionRevoked`). |
| `device.changed` | zdarzenie | Zmiana w rejestrze urządzeń, rozgłaszana do pozostałych urządzeń konta. |

Ładunek zgłoszenia urządzenia i odpowiedzi zamykającej protokół:

```json
// Urządzenie przenośne → Serwer, krok 4
{ "type": "device.pair.confirm", "id": "p-3", "payload": {
    "pairingCode": "••••-••••", "deviceId": "dev-9127",
    "deviceType": "telefon", "deviceName": "Telefon Operatora" },
  "timestamp": "2026-08-06T10:14:00Z" }

// Serwer → Urządzenie przenośne, krok 7
{ "type": "device.pair.completed", "id": "e-4410", "payload": {
    "deviceId": "dev-9127", "credentialRef": "•••",
    "credentialExpiresAt": "2026-11-04T10:14:00Z",
    "mobileEnabled": true },
  "timestamp": "2026-08-06T10:14:07Z" }
```

### 9.11. Warunki błędu

Struktura obiektu `error` jest wspólna dla całej platformy: `code`, `message`, `details`, `retryable` (Kontrakty komunikacji, rozdz. 19). Błąd protokołu nie przerywa połączenia WebSocket.

| Kod | Znaczenie | `retryable` | Zachowanie interfejsu |
|---|---|---|---|
| `pairing_code_invalid` | Podany kod parowania nie istnieje albo należy do innego konta. | Nie | Komunikat kontekstowy przy polu kodu na urządzeniu przenośnym; kod pozostaje do poprawienia |
| `pairing_code_expired` | Kod parowania utracił ważność przed użyciem. | Nie | Okno nakładkowe na stanowisku stacjonarnym prezentuje działanie „Wydaj nowy kod”; licznik ważności wyzerowany |
| `pairing_code_used` | Kod parowania został już użyty przez inne urządzenie. | Nie | Komunikat kontekstowy z odesłaniem do rejestru urządzeń |
| `pairing_rejected` | Operator odmówił wpuszczenia urządzenia w kroku 6. | Nie | Urządzenie przenośne wraca do ekranu logowania z komunikatem odmowy |
| `pairing_timeout` | Zgłoszenie nie zostało potwierdzone na maszynie inicjującej w czasie ważności kodu. | Tak | Okno nakładkowe proponuje wydanie nowego kodu; urządzenie przenośne prezentuje komunikat o braku potwierdzenia |
| `credential_expired` | Przedstawione poświadczenie urządzenia utraciło ważność. | Nie | Ekran logowania z komunikatem „poświadczenie urządzenia wygasło”; kolejka lokalna zachowana |
| `credential_revoked` | Poświadczenie urządzenia zostało odwołane. | Nie | Ekran logowania z komunikatem „dostęp urządzenia odwołany”; zakres lokalny i kolejka usunięte |
| `not_authenticated` | Połączenie nie zostało uwierzytelnione, a faza wdrożenia tego wymaga. | Nie | Ekran logowania |
| `validation_failed` | Ładunek polecenia nie spełnia wymagań pola lub typu. | Nie | Komunikat kontekstowy przy polu wskazanym w `details` |
| `conflict` | Operacja koliduje z bieżącym stanem rejestru — na przykład odwołanie poświadczenia już odwołanego. | Nie | Komunikat kontekstowy przy wierszu rejestru; rejestr odświeżany |
| `internal_error` | Błąd wewnętrzny serwera, niezwiązany z treścią polecenia. | Tak | Działanie „Ponów” przy komunikacie |

---

## 10. Punkty sterowania z okna konfiguracji i z okna ustawień

Model konfiguracji obejmuje warstwy **globalny → środowisko → projekt → sesja**. Brak jawnego ustawienia oznacza wartość domyślną i zachowanie kanoniczne. Wszystkie zakresy są sterowane przez Operatora bez blokad w interfejsie.

### 10.1. Punkty sterowania w oknie Ustawień

| Zakres | Co Operator ustawia | Miejsce | Warstwa konfiguracji |
|---|---|---|---|
| Urządzenia połączone | Wykaz urządzeń konta, nazwa urządzenia, odwołanie dostępu | [Okno Ustawień](../interfejs-uzytkownika/ustawienia.md), rozdz. 6.1 | globalna |
| Parowanie urządzenia | Wydanie kodu parowania, potwierdzenie zgłoszenia, unieważnienie kodu | [Okno Ustawień](../interfejs-uzytkownika/ustawienia.md), rozdz. 6.2 | globalna |
| Synchronizacja | Prezentacja stanu synchronizacji urządzeń konta | [Okno Ustawień](../interfejs-uzytkownika/ustawienia.md), rozdz. 6.3 | globalna |
| Powiadomienia — kategorie i kanał | Które kategorie zdarzeń i którym kanałem są dostarczane, w tym kanał powiadomienia wypychanego | [Okno Ustawień](../interfejs-uzytkownika/ustawienia.md), rozdz. 7 | globalna, sesja |
| Metody uwierzytelniania | Metody dostępne na ekranie logowania urządzenia przenośnego | [Okno Ustawień](../interfejs-uzytkownika/ustawienia.md), rozdz. 4 | globalna |
| Motyw i język | Motyw jasny albo ciemny oraz język interfejsu, wspólne dla wszystkich urządzeń konta | [Okno Ustawień](../interfejs-uzytkownika/ustawienia.md), rozdz. 5 | globalna |

### 10.2. Punkty sterowania w oknie konfiguracji

| Zakres | Co Operator ustawia | Warstwa konfiguracji |
|---|---|---|
| Chat Window | Zasięg kontekstu, retencja historii, próg ostrzeżenia o zapełnieniu kontekstu, domyślny poziom wysiłku Wykonawcy | globalna / środowisko / sesja |
| Execution Loop Window | Domyślna widoczność okna, zakres prezentowanego strumienia komunikatów sterujących, próg kontroli jakości uruchamiający ponowienie, limit ponowień zadania | globalna / środowisko / projekt |
| Punkty obowiązkowego zatwierdzenia | Zestaw działań wymagających zatwierdzenia przez Użytkownika i granice autonomii Wykonawcy | globalna / środowisko / projekt |
| Zakres funkcji Mobile | Które klasy procesów są prezentowane w przekroju urządzenia przenośnego; domyślny zakres filtra przekroju | globalna / środowisko |
| Praca bez połączenia | Utrzymywanie zakresu lokalnego, retencja migawki, retencja kolejki działań lokalnych | globalna |
| Warstwy widoczności według roli | Zestaw funkcji ujawnianych roli użytkownika na urządzeniu przenośnym, zakres wyszukiwarki funkcji, dostęp do trybu administracyjnego | globalna / rola |
| Poświadczenia urządzeń | Termin ważności poświadczenia, zachowanie przy odświeżeniu, ważność kodu parowania | globalna |
| Izolacja techniczna i kontekstu | Zasięg izolacji obejmujący sesje nadzorowane z urządzenia przenośnego; stan wyjściowy „pełny dostęp” | Okno punktów izolacji |

### 10.3. Punkty sterowania na urządzeniu przenośnym

| Zakres | Co Użytkownik ustawia | Miejsce |
|---|---|---|
| Nazwa urządzenia | Nazwa własna widoczna w rejestrze urządzeń | Ekran ustawień urządzenia (rozdz. 4.8) |
| Klasy powiadomień tego urządzenia | Które klasy zdarzeń są dostarczane na ten egzemplarz | Ekran ustawień urządzenia |
| Zakres lokalny | Utrzymywanie migawki i kolejki działań lokalnych na tym urządzeniu | Ekran ustawień urządzenia |
| Domyślny filtr przekroju | Zakres, rodzaj i stan procesów zapisane jako domyślne dla tego urządzenia | Kebab ekranu przeglądu zadań i procesów |
| Zakończenie dostępu | Wylogowanie tego urządzenia | Ekran ustawień urządzenia |

---

## 11. Stany, dane i powiązania

### 11.1. Stany funkcji

| Stan funkcji | Wyzwalacz | Prezentacja | Wyjście ze stanu |
|---|---|---|---|
| Nierozpoznana | Urządzenie bez poświadczenia | Ekran logowania z działaniem „Mam kod parowania” | Parowanie albo logowanie |
| Uwierzytelniona, połączona | Ważne poświadczenie, otwarty kanał | Wskaźnik `● połączony`; strumienie aktualizowane na żywo | Utrata łączności albo odwołanie dostępu |
| Uwierzytelniona, wznawianie | Przerwa w łączności | Wskaźnik `⇄ wznawianie`; interfejs sterowalny, działania kolejkowane | Wznowienie połączenia |
| Bez połączenia | Brak łączności | Wskaźnik `⇄ bez połączenia`; zakres lokalny odczytowy, licznik kolejki widoczny | Odzyskanie połączenia i uzgodnienie stanu |
| Uzgadnianie | Odzyskanie połączenia z niepustą kolejką | Plakietka „uzgadnianie” przy liczniku kolejki | Zamknięcie uzgadniania |
| Poświadczenie wygasłe | Upływ terminu ważności | Ekran logowania z komunikatem wygaśnięcia | Ponowne uwierzytelnienie |
| Dostęp odwołany | Odwołanie poświadczenia | Ekran logowania z komunikatem odwołania; zakres lokalny usunięty | Ponowne parowanie |

### 11.2. Stany elementów

| Element | Stany |
|---|---|
| Chat Window | Bezczynne · Wykonawca pracuje · oczekiwanie na zatwierdzenie · odtworzone z migawki · działanie zakolejkowane |
| Execution Loop Window | Zwinięte · rozwinięte, pętla bezczynna · rozwinięte, pętla w toku · rozwinięte, pętla wstrzymana · rozwinięte, punkt decyzyjny oczekuje |
| Pozycja procesu | W toku · decyzja · wstrzymany · zakończony · błąd · bezprzedmiotowy |
| Pozycja powiadomienia | Nowa · przeczytana · rozstrzygnięta · wygasła |
| Kolumna | Rozwinięta · zwężona do szerokości minimalnej · zwinięta do wyzwalacza · aktywna w przełączniku kolumn |
| Kolejka działań lokalnych | Pusta · zapełniona · uzgadniana · z pozycjami odrzuconymi |
| Poświadczenie urządzenia | Aktywne · zbliża się wygaśnięcie · wygasłe · odwołane |
| Kod parowania | Wydany · użyty · unieważniony · wygasły |

### 11.3. Encje modelu danych

Funkcja pracuje na encjach opisanych w [Modelu danych](../architektura/model-danych.md). Rozdział 18.1 tego dokumentu wskazuje odwzorowanie aspektów funkcji w encjach istniejących; poniższa tabela zestawia pełny zakres.

| Encja | Rozdział [Modelu danych](../architektura/model-danych.md) | Rola w funkcji Mobile |
|---|---|---|
| `konto` | 3.1 | Jedyne konto właściciela; właściciel wszystkich urządzeń |
| `urzadzenie` | 3.2 | Rejestr sparowanych urządzeń; pola `typ`, `nazwa`, `tryb_mobile_aktywny`, `data_ostatniego_polaczenia` |
| `polaczenie` | 3.3 | Połączenie WebSocket urządzenia przenośnego z serwerem |
| `poswiadczenie_urzadzenia` | 3.4 | Poświadczenie wydane w parowaniu: odwołanie do wartości, termin ważności, stan |
| `kod_parowania` | 3.5 | Kod wydany w kroku 2 protokołu: odwołanie, termin ważności, stan użycia |
| `dzialanie_lokalne` | 3.6 | Kolejka działań wykonanych bez połączenia wraz z wersją stanu i wynikiem uzgodnienia |
| `sesja`, `proces_sesji` | 8.2, 8.3 | Przedmiot monitoringu i sterowania z urządzenia przenośnego |
| `wiadomosc` | 8.4 | Strumień Chat Window prezentowany i uzupełniany z urządzenia przenośnego |
| `zadanie`, `kolejka`, `pozycja_kolejki` | 10.1, 10.3, 10.4 | Przedmiot zarządzania zadaniami: uruchomienie, zatrzymanie, wstrzymanie, wznowienie |
| `log_akcji_kolejki` | 10.5 | Zapis akcji wykonanych z urządzenia przenośnego |
| `przebieg_multitasking`, `decyzja_procesu` | 13.5, 13.6 | Przebiegi nadzorowane i punkty decyzyjne rozstrzygane z urządzenia przenośnego |
| `automatyka` | 6.2 | Automatyzacje uruchamiane i wstrzymywane z urządzenia przenośnego |
| `artefakt`, `wersja_artefaktu` | 16.1, 16.2 | Wynik podglądany przed decyzją o zatwierdzeniu |
| `ustawienie` | 17.1 | Ustawienia funkcji i ustawienia tego urządzenia |

### 11.4. Powiązania między oknami

| Okno źródłowe | Okno docelowe | Charakter powiązania |
|---|---|---|
| Chat Window | Execution Loop Window | Polecenie wydane w Chat Window uruchamia albo koryguje pętlę wykonawczą i zapełnia listę dekompozycji |
| Execution Loop Window | Chat Window | Wynik zadania i decyzja kontroli jakości trafiają do strumienia rozmowy jako wpis z referencją |
| Strona główna Mobile | Execution Loop Window | Kliknięcie pozycji decyzji otwiera okno pętli na wskazanym punkcie decyzyjnym |
| Strona główna Mobile | Przegląd zadań i procesów | Kliknięcie pozycji przebiegu otwiera przegląd zawężony do tego przebiegu |
| Powiadomienia | Execution Loop Window | Działanie „Otwórz pętlę” prowadzi do punktu decyzyjnego zdarzenia |
| Powiadomienia | Podgląd artefaktu | Działanie „Otwórz wynik” otwiera artefakt w kolumnie dominującej |
| Przegląd zadań i procesów | Podgląd artefaktu | Kliknięcie referencji wyniku otwiera podgląd jako rozszerzenie boczne |
| Ustawienia urządzenia | Rejestr urządzeń konta | Pozycja kebaba otwiera rejestr jako rozszerzenie boczne |
| Wszystkie ekrany | Kolejka działań lokalnych | Działanie wykonane bez łączności zapisuje pozycję kolejki i plakietkę „zakolejkowane” |
| Wszystkie ekrany | Dziennik audytu | Zlecenie, zadanie, komunikat sterujący i decyzja podlegają rejestracji (Bezpieczeństwo, rozdz. 13.5) |

### 11.5. Powiązania z pozostałymi dokumentami zbioru

| Dokument | Co wnosi do funkcji |
|---|---|
| [Koncepcja platformy](../architektura/koncepcja-platformy.md) | Definicja funkcji globalnej, tabela możliwości Mobile, macierz funkcji globalnych |
| [Bezpieczeństwo i uwierzytelnianie](../architektura/bezpieczenstwo-i-uwierzytelnianie.md) | Model jednego konta i wielu urządzeń, token urządzenia, magazyn danych dostępowych, punkty obowiązkowego zatwierdzenia, dziennik audytu |
| [Kontrakty komunikacji](../architektura/kontrakty-komunikacji.md) | Koperta komunikatu, cykl życia połączenia, komunikaty funkcji Mobile, synchronizacja wielourządzeniowa, obsługa błędów |
| [Model danych](../architektura/model-danych.md) | Encje, na których funkcja pracuje |
| [System wizualny](../interfejs-uzytkownika/system-wizualny.md) | Siatka kolumnowa, waga wizualna, punkty przełamania, tokeny, komponenty, ikonografia |
| [Strona główna i nawigacja](../interfejs-uzytkownika/strona-glowna-i-nawigacja.md) | Punkt wejścia funkcji w listwie ustawień i jej równoważne punkty dostępu |
| [Okno Ustawień](../interfejs-uzytkownika/ustawienia.md) | Sekcja urządzeń, parowanie, sekcja powiadomień |
| [Elementy okien](../interfejs-uzytkownika/elementy-okien.md) | Rama okna (szyna nawigacji, belka, wstążka, pasek stanu), karty sesji, wskaźnik stanu procesu w tle |
| [Przepływ okien](../interfejs-uzytkownika/przeplyw-okien.md) | Trwałość procesów w tle jako podstawa funkcji, aktywacja funkcji bez nowej przestrzeni roboczej |
| [Środowisko MultitaskingAI](../srodowiska/multitaskingai.md) | Monitor procesu, przebiegi, decyzje procesu nadzorowane z urządzenia przenośnego |

---

## 12. Scenariusze użycia

### 12.1. Sparowanie nowego telefonu

| Krok | Miejsce | Działanie |
|---|---|---|
| 1 | Stanowisko stacjonarne, okno Ustawień → Urządzenia | Operator wybiera działanie „Sparuj urządzenie”; okno nakładkowe prezentuje kod parowania i licznik ważności |
| 2 | Telefon | Aplikacja odczytuje kod graficzny; zgłoszenie trafia do serwera z typem `telefon` i nazwą wstępną |
| 3 | Stanowisko stacjonarne | Zdarzenie zgłoszenia prezentuje typ, nazwę i czas; Operator potwierdza wpuszczenie urządzenia |
| 4 | Serwer | Wydanie poświadczenia, unieważnienie kodu, wpis w rejestrze urządzeń |
| 5 | Stanowisko stacjonarne | Nowy wiersz w tabeli urządzeń i dymek powiadomienia „Sparowano: Telefon Operatora” |
| 6 | Telefon | Przejście do strony głównej Mobile; przekrój prezentuje decyzje oczekujące i przebiegi w toku |

### 12.2. Zatwierdzenie kroku nieodwracalnego poza stanowiskiem

| Krok | Miejsce | Działanie |
|---|---|---|
| 1 | Telefon, poza aplikacją | Powiadomienie klasy „Decyzje i zatwierdzenia”: publikacja pakietu wymaga zatwierdzenia |
| 2 | Telefon | Dotknięcie powiadomienia otwiera Execution Loop Window na zadaniu oznaczonym `⚑` |
| 3 | Execution Loop Window | Podgląd dekompozycji zlecenia i komunikatów sterujących; referencja wyniku otwiera podgląd artefaktu |
| 4 | Execution Loop Window | Użytkownik wybiera „Zatwierdź”; pętla podejmuje zadanie 3 |
| 5 | Chat Window | Wykonawca potwierdza publikację i przedstawia raport wydania |
| 6 | Stanowisko stacjonarne | Zdarzenie rozgłoszone przez serwer aktualizuje przebieg pętli w otwartym oknie |

### 12.3. Przerwanie zadania z urządzenia przenośnego

| Krok | Miejsce | Działanie |
|---|---|---|
| 1 | Przegląd zadań i procesów | Pozycja „Budowa pakietu” prezentuje stan błędu po nieudanym kroku |
| 2 | Pozycja procesu | Element zbiorczy `Operacje ▼` udostępnia akcje sterowania zadaniem |
| 3 | Execution Loop Window | Otwarcie okna pętli prezentuje komunikaty sterujące i wynik kontroli jakości |
| 4 | Execution Loop Window | Użytkownik przerywa zadanie; stan zadań pozostaje zachowany po stronie serwera |
| 5 | Chat Window | Polecenie „ponów budowę po korekcie zależności” koryguje zlecenie i podejmuje pętlę |

### 12.4. Praca bez łączności i uzgodnienie po powrocie

| Krok | Miejsce | Działanie |
|---|---|---|
| 1 | Strona główna Mobile | Wskaźnik zmienia się na `⇄ bez połączenia`; sekcje prezentują migawkę ze znacznikiem czasu |
| 2 | Powiadomienia | Użytkownik zatwierdza decyzję z poziomu listy; działanie trafia do kolejki działań lokalnych |
| 3 | Chat Window | Użytkownik wydaje polecenie wstrzymania automatyzacji nocnej; polecenie zostaje zakolejkowane |
| 4 | Ustawienia urządzenia | Licznik działań w kolejce prezentuje dwie pozycje wraz z czasem zapisu |
| 5 | Strona główna Mobile | Odzyskanie łączności uruchamia uzgadnianie: kolejka przekazywana jest z wersjami stanu |
| 6 | Chat Window | Zatwierdzenie zostaje przyjęte; polecenie wstrzymania odrzucone jako bezprzedmiotowe — automatyzacja zakończyła się wcześniej |
| 7 | Ustawienia urządzenia | Pozycja odrzucona pozostaje na liście z przyczyną do chwili potwierdzenia przez Użytkownika |

### 12.5. Nadzór nad przebiegiem MultitaskingAI

| Krok | Miejsce | Działanie |
|---|---|---|
| 1 | Strona główna Mobile | Sekcja „Przebiegi w toku” prezentuje przebieg zespołu środowiska MultitaskingAI |
| 2 | Przegląd zadań i procesów | Filtr `Zakres ▼` zawęża listę do środowiska MultitaskingAI |
| 3 | Execution Loop Window | Podgląd dekompozycji zlecenia, kolejki zadań i wyników kontroli jakości |
| 4 | Execution Loop Window | Zadanie odrzucone w kontroli jakości wraca do kolejki decyzją o ponowieniu |
| 5 | Chat Window | Wykonawca przedstawia wynik ponowienia; Użytkownik zatwierdza zamknięcie zlecenia |
| 6 | Powiadomienia | Zdarzenie zakończenia przebiegu zamyka serię powiadomień tego zlecenia |

### 12.6. Utrata telefonu i odwołanie dostępu

| Krok | Miejsce | Działanie |
|---|---|---|
| 1 | Tablet Operatora, ustawienia urządzenia | Otwarcie rejestru urządzeń konta z kebaba ekranu |
| 2 | Rejestr urządzeń konta | Rozpoznanie urządzenia utraconego po nazwie, typie i dacie ostatniego połączenia |
| 3 | Rejestr urządzeń konta | Akcja „Odłącz” unieważnia poświadczenie telefonu |
| 4 | Serwer | Zamknięcie otwartych połączeń telefonu; odrzucenie działań nadesłanych po odwołaniu |
| 5 | Wszystkie urządzenia konta | Zdarzenie `device.changed` aktualizuje rejestr; powiadomienie klasy „Konto i dostęp” |
| 6 | Stanowisko stacjonarne | Wpis dziennika audytu rejestruje odwołanie dostępu wraz z urządzeniem inicjującym |

---

## Załącznik A. Skróty, gesty i ikonografia

### A.1. Gesty urządzenia przenośnego

Gesty odpowiadają regulacji szerokości kolumn i przełączaniu kolumn — nigdy nie zmieniają układu z pionowego na inny.

| Gest | Działanie |
|---|---|
| Przesunięcie poziome w obszarze roboczym | Przejście do kolumny sąsiadującej w kolejności od lewej do prawej |
| Przesunięcie poziome od lewej krawędzi | Wysunięcie kolumny nawigacji funkcji |
| Przesunięcie poziome wiersza listy | Odsłonięcie zestawu akcji procesu `Operacje ▼` |
| Przeciągnięcie granicy między kolumnami | Regulacja szerokości kolumny |
| Przytrzymanie pozycji listy | Menu kontekstowe pozycji |
| Przeciągnięcie listy w dół | Odświeżenie przekroju z serwera |
| Przytrzymanie godła platformy | Podgląd wersji kontraktu i diagnostyka połączenia |

### A.2. Skróty klawiatury zewnętrznej

Mapa skrótów jest w całości konfigurowalna z edytora mapy skrótów (rozdz. 10.2). Konflikt przypisań sygnalizowany jest ostrzeżeniem, nie blokadą.

| Skrót | Działanie |
|---|---|
| `Ctrl/Cmd + K` | Wyszukiwarka funkcji |
| `Ctrl/Cmd + L` | Fokus na pole polecenia Chat Window |
| `Ctrl/Cmd + J` | Execution Loop Window |
| `Ctrl/Cmd + Enter` | Zatwierdzenie punktu decyzyjnego |
| `Ctrl/Cmd + .` | Przerwanie bieżącego działania Wykonawcy |
| `Ctrl/Cmd + Alt + ←` / `→` | Przejście do kolumny poprzedniej albo następnej |
| `Ctrl/Cmd + Shift + N` | Ekran powiadomień |
| `Ctrl/Cmd + Shift + Z` | Przegląd zadań i procesów |
| `Ctrl/Cmd + ,` | Ustawienia urządzenia |
| `?` | Ściągawka skrótów i gestów |

### A.3. Ikonografia funkcji

Ikony pochodzą z zestawu opisanego w [Systemie wizualnym](../interfejs-uzytkownika/system-wizualny.md), rozdz. 7; rysunek liniowy na siatce 24×24, jeden styl w całym zestawie.

| Znak | Znaczenie | Miejsce występowania |
|---|---|---|
| `◈` | Referencja zasobu, artefakt, dziennik przebiegu | Strumień Chat Window, pozycja procesu |
| `⚑` | Punkt decyzyjny oczekujący na zatwierdzenie | Lista dekompozycji, sekcja decyzji oczekujących, powiadomienia |
| `⟳` | Proces w toku, licznik procesów | Lista procesów, wiersz podsumowania |
| `⏸` | Proces wstrzymany | Lista procesów |
| `✔` | Zadanie albo przebieg zakończony | Lista dekompozycji, lista procesów, powiadomienia |
| `◌` | Zadanie oczekujące w kolejce | Lista dekompozycji |
| `⚠` | Błąd zadania | Lista procesów, powiadomienia |
| `⇄` | Stan kanału połączenia | Wskaźnik połączenia |
| `●` | Kropka stanu semantycznego | Wskaźnik połączenia, plakietki stanu |
| `▮▮▮▯` | Wskaźnik przebiegu pętli | Execution Loop Window, ekran startowy |
| `▼` | Zwinięty wyzwalacz warstwy 2 | Nagłówek Chat Window, filtry, elementy zbiorcze |
| `⋮` | Menu kebab warstwy 3 | Nagłówki okien, pozycje list |
| `☰` | Menu nawigacji funkcji warstwy 3 | Nagłówek Chat Window przy kolumnie zwiniętej |
| `➤` | Wysłanie polecenia | Pole polecenia Chat Window |
| `→` | Przejście do ekranu powiązanego | Ekran startowy, powiadomienia, ustawienia urządzenia |

---

## Załącznik B. Zbiorcza tabela komunikatów funkcji Mobile

| Typ komunikatu | Kierunek | Rozdział | Opis |
|---|---|---|---|
| `mobile.status.get` | polecenie | 1.4, 11.1 | Stan dostępu mobilnego dla bieżącego urządzenia. |
| `mobile.process.list` | polecenie | 4.6 | Podgląd projektów, automatyzacji, agentów, aplikacji i sesji — monitoring procesów niezależny od aktywnego środowiska. |
| `mobile.process.control` | polecenie | 4.6, 5.3 | Zarządzanie zadaniem: `start`, `stop`, `approve`, `pause`, `resume`, `modify`. |
| `mobile.process.changed` | zdarzenie | 7.3, 8.1 | Zmiana stanu monitorowanego procesu, w szczególności procesu długotrwałego uruchomionego przez środowisko MultitaskingAI. |
| `mobile.queue.sync` | polecenie | 7.3 | Przekazanie kolejki działań lokalnych wraz z wersjami stanu do uzgodnienia po odzyskaniu połączenia. Nazwa z [Kontraktów komunikacji](../architektura/kontrakty-komunikacji.md) rozdz. 16.1; brak w `contract.json` — [DO DECYZJI OPERATORA]. |
| `mobile.notification.action` | polecenie | 8.3 | Działanie wykonane z poziomu powiadomienia wypychanego. Nazwa z [Kontraktów komunikacji](../architektura/kontrakty-komunikacji.md) rozdz. 16.1; brak w `contract.json` — [DO DECYZJI OPERATORA]. |
| `message.send` | polecenie | 5.1 | Polecenie języka naturalnego w kanale Użytkownik ↔ Wykonawca. |
| `loop.control` | polecenie | 5.2 | Sterowanie przebiegiem pętli: zatwierdzenie, wstrzymanie, wznowienie, przerwanie, korekta zlecenia. |
| `device.pair.start` | polecenie | 9.10 | Inicjowanie parowania z okna Ustawień. |
| `device.pair.cancel` | polecenie | 9.10 | Unieważnienie kodu parowania przed użyciem. |
| `device.pair.confirm` | polecenie | 9.10 | Zgłoszenie urządzenia przenośnego z kodem parowania. |
| `device.pair.request` | zdarzenie | 9.10 | Zgłoszenie przekazane do potwierdzenia na maszynie inicjującej. |
| `device.pair.approve` | polecenie | 9.10 | Potwierdzenie wpuszczenia urządzenia do rejestru. |
| `device.pair.reject` | polecenie | 9.10 | Odmowa wpuszczenia urządzenia. |
| `device.pair.completed` | zdarzenie | 9.10 | Wydanie poświadczenia i zamknięcie protokołu. |
| `auth.token.refresh` | polecenie | 9.7 | Przedłużenie sesji bramki urządzenia; wydanie `AuthSession` o nowym terminie ważności. |
| `device.list` | polecenie | 9.6 | Odczyt rejestru urządzeń konta. |
| `device.rename` | polecenie | 9.6 | Zmiana nazwy urządzenia w rejestrze. |
| `device.revoke` | polecenie | 9.8 | Odwołanie poświadczenia wskazanego urządzenia. |
| `auth.logout` | polecenie | 9.8 | Unieważnienie sesji bramki tego urządzenia; skutek rozgłaszany zdarzeniem `auth.changed` (powód `sessionRevoked`). |
| `device.changed` | zdarzenie | 9.8, 9.9 | Zmiana w rejestrze urządzeń, rozgłaszana do pozostałych urządzeń konta. |

---

## Załącznik C. Mapa zgodności z dokumentami źródłowymi

| Rozdział niniejszego dokumentu | Dokument źródłowy | Rozdział źródła |
|---|---|---|
| 1. Przeznaczenie i kontekst | [Koncepcja platformy](../architektura/koncepcja-platformy.md); [Strona główna i nawigacja](../interfejs-uzytkownika/strona-glowna-i-nawigacja.md); [Przepływ okien](../interfejs-uzytkownika/przeplyw-okien.md) | 12.1, 5.3, załącznik C; 3.4, 9.2; 5.6, 8.4 |
| 2. Zakres funkcjonalny | [Koncepcja platformy](../architektura/koncepcja-platformy.md); [Model konfiguracji](../architektura/model-konfiguracji.md); [Izolacja i konfigurowalność zależności](../architektura/izolacja-i-zaleznosci.md) | 12.1, 5.3; 5.1; 3 |
| 3. Układ interfejsu mobilnego | [System wizualny](../interfejs-uzytkownika/system-wizualny.md) | 4.3, 4.4, 3a.5 |
| 4. Komplet okien i ekranów | [Elementy okien](../interfejs-uzytkownika/elementy-okien.md); [System wizualny](../interfejs-uzytkownika/system-wizualny.md); [Specyfikacja okien operacyjnych](../specyfikacje/specyfikacja-okien-operacyjnych.md) | 3–5; 8, 10; 2.4 |
| 5. Kanały komunikacji operacyjnej | [Koncepcja platformy](../architektura/koncepcja-platformy.md); [Kontrakty komunikacji](../architektura/kontrakty-komunikacji.md); [Bezpieczeństwo i uwierzytelnianie](../architektura/bezpieczenstwo-i-uwierzytelnianie.md) | 1, 11.8; 9.1, 9.2; 13.3, 13.4 |
| 6. Warstwy widoczności | [Koncepcja platformy](../architektura/koncepcja-platformy.md); [System wizualny](../interfejs-uzytkownika/system-wizualny.md); [Bezpieczeństwo i uwierzytelnianie](../architektura/bezpieczenstwo-i-uwierzytelnianie.md) | 12; 3a; 14 |
| 7. Praca bez połączenia | [Kontrakty komunikacji](../architektura/kontrakty-komunikacji.md); [Architektura techniczna](../architektura/architektura.md); [Model danych](../architektura/model-danych.md) | 2, 18; 7, 16; 3.6 |
| 8. Powiadomienia wypychane | [Okno Ustawień](../interfejs-uzytkownika/ustawienia.md); [Kontrakty komunikacji](../architektura/kontrakty-komunikacji.md) | 7; 9.2, 11.6 |
| 9. Protokół parowania | [Bezpieczeństwo i uwierzytelnianie](../architektura/bezpieczenstwo-i-uwierzytelnianie.md); [Kontrakty komunikacji](../architektura/kontrakty-komunikacji.md); [Okno Ustawień](../interfejs-uzytkownika/ustawienia.md); [Model danych](../architektura/model-danych.md) | 3, 8, 10.2, 12; 2, 5, 19; 6.2; 3.4, 3.5 |
| 10. Punkty sterowania | [Model konfiguracji](../architektura/model-konfiguracji.md); [Okno Ustawień](../interfejs-uzytkownika/ustawienia.md) | 4.1, 5.1; 4–7 |
| 11. Stany, dane i powiązania | [Model danych](../architektura/model-danych.md); [Bezpieczeństwo i uwierzytelnianie](../architektura/bezpieczenstwo-i-uwierzytelnianie.md) | 3, 8, 10, 13, 16, 17, 18.1; 13.5 |
| 12. Scenariusze użycia | [Środowisko MultitaskingAI](../srodowiska/multitaskingai.md); [Okno Ustawień](../interfejs-uzytkownika/ustawienia.md) | 4, 5; załącznik A |

---

## Załącznik — pełny wykaz komend kontraktu obszaru `mobile`

Wykaz wygenerowany wprost z `budowa/shared/contract.json` — jedynego źródła nazw komend, pól ładunków i zdarzeń. Każda pozycja jest wiążąca dla warstwy klienckiej i serwerowej; deweloper nie dodaje ani nie zmienia nazw poza tym wykazem. Komunikaty obszarów sąsiednich wykorzystywane przez funkcję — `message`, `device`, `auth`, `loop` — zestawia Załącznik B.

### Obszar `mobile` — 3 komendy

| Komenda | Przeznaczenie | Pola żądania | Pola wyniku |
|---|---|---|---|
| `mobile.status.get` | Zwraca stan platformy dla mobilnego centrum dowodzenia | `deviceId:string` (opc) | `status:MobileStatus` (wym) |
| `mobile.process.list` | Zwraca procesy platformy widoczne z urządzenia przenośnego | `deviceId:string` (opc)<br>`sessionId:string` (opc)<br>`status:ProgressStatus` (opc) | `processes:MobileProcess[]` (wym) |
| `mobile.process.control` | Steruje procesem platformy z urządzenia przenośnego | `processId:string` (wym)<br>`control:MobileProcessControl` (wym)<br>`deviceId:string` (opc) | `process:MobileProcess` (wym) |

**Zdarzenia obszaru `mobile` — 1:**

| Zdarzenie | Przeznaczenie | Pola |
|---|---|---|
| `mobile.process.changed` | Zmiana procesu platformy zgłoszona warstwie mobilnej | `change:ChangeKind` (wym)<br>`process:MobileProcess` (wym) |

Razem w wykazie: **3 komendy** z 1 obszaru kontraktu. Dwie dalsze nazwy przywoływane w niniejszym opracowaniu — `mobile.queue.sync` (rozdz. 7.3) oraz `mobile.notification.action` (rozdz. 8.3) — pochodzą z rozdziału 16.1 [Kontraktów komunikacji](../architektura/kontrakty-komunikacji.md) i nie mają jeszcze odpowiednika w `budowa/shared/contract.json`; obie pozostają oznaczone **[DO DECYZJI OPERATORA]** w miejscach ich opisu.

---

*Koniec dokumentu. Funkcja globalna Mobile — Specyfikacja docelowa, wersja 2.0, 2026-08-20.*

---
*Danaco Console — Platforma AI Workspace OS · v2.0 · status Deweloperski*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz.*
*Warunki korzystania: [Licencja produktu](../LICENSE.md). Kontakt: support@danaco-group.pl*
