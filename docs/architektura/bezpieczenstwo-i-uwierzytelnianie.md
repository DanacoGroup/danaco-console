# Danaco Console — Bezpieczeństwo i uwierzytelnianie

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
| **Źródło** | Koncepcja platformy i architektura |

Dokument opisuje model bezpieczeństwa platformy Danaco Console: uwierzytelnianie jako mechanizm kontroli dostępu do platformy — rejestrację konta właściciela przy pierwszym uruchomieniu, logowanie, metody dodatkowe uwierzytelniania (PIN, Windows Hello, adres e-mail uwierzytelniający), odzyskiwanie konta, wydawanie tokenu dostępu i jego rolę w nawiązywaniu połączenia WebSocket, zasadę przechowywania danych dostępowych poza bazą danych, uwierzytelnianie i autoryzację w obu kanałach komunikacji operacyjnej — Chat Window (Użytkownik ↔ Wykonawca) i Execution Loop Window (Koordynator ↔ Wykonawca) — oraz warstwy widoczności jako mechanizm kontroli dostępu do funkcji platformy.

---

## Spis treści

- [Wprowadzenie](#wprowadzenie)
1. [Cel i zakres dokumentu](#1-cel-i-zakres-dokumentu)
2. [Zasada nadrzędna: uwierzytelnianie jako mechanizm kontroli dostępu do platformy](#2-zasada-nadrzędna-uwierzytelnianie-jako-mechanizm-kontroli-dostępu-do-platformy)
3. [Model właściciela: jedno konto, wiele urządzeń](#3-model-właściciela-jedno-konto-wiele-urządzeń)
4. [Rejestracja](#4-rejestracja)
5. [Logowanie](#5-logowanie)
6. [Metody dodatkowe uwierzytelniania](#6-metody-dodatkowe-uwierzytelniania)
7. [Odzyskiwanie konta](#7-odzyskiwanie-konta)
8. [Token dostępu i połączenie WebSocket](#8-token-dostępu-i-połączenie-websocket)
9. [Przechowywanie danych dostępowych poza bazą danych](#9-przechowywanie-danych-dostępowych-poza-bazą-danych)
10. [Model danych tożsamości i dostępu](#10-model-danych-tożsamości-i-dostępu)
11. [Komunikacja — polecenia i zdarzenia uwierzytelniania](#11-komunikacja--polecenia-i-zdarzenia-uwierzytelniania)
12. [Ustawienia uwierzytelniania w oknie konfiguracji](#12-ustawienia-uwierzytelniania-w-oknie-konfiguracji)
13. [Bezpieczeństwo kanałów komunikacji operacyjnej](#13-bezpieczeństwo-kanałów-komunikacji-operacyjnej)
14. [Warstwy widoczności jako mechanizm kontroli dostępu do funkcji](#14-warstwy-widoczności-jako-mechanizm-kontroli-dostępu-do-funkcji)
15. [Zgodność z zasadami nadrzędnymi platformy](#15-zgodność-z-zasadami-nadrzędnymi-platformy)
16. [Słowniczek pojęć](#16-słowniczek-pojęć)
- [Załącznik A. Schemat przepływu uwierzytelniania](#załącznik-a-schemat-przepływu-uwierzytelniania)
- [Załącznik B. Macierz metod uwierzytelniania](#załącznik-b-macierz-metod-uwierzytelniania)
- [Załącznik C. Scenariusze użycia](#załącznik-c-scenariusze-użycia)
- [Załącznik D. Szablony konfiguracji tożsamości i dostępu](#załącznik-d-szablony-konfiguracji-tożsamości-i-dostępu)

---

## Wprowadzenie

Model bezpieczeństwa platformy Danaco Console obejmuje trzy powiązane mechanizmy. Uwierzytelnianie rozstrzyga dostęp do platformy jako całości: urządzenie uzyskuje dostęp wyłącznie po przedstawieniu ważnego tokenu wydanego właścicielowi konta. Autoryzacja w kanałach komunikacji operacyjnej rozstrzyga, jakie działania wolno wykonać Koordynatorowi i Wykonawcy w imieniu uwierzytelnionego Użytkownika oraz w których punktach przebiegu wymagane jest zatwierdzenie przez Użytkownika. Warstwy widoczności rozstrzygają, które funkcje platformy są ujawniane danej roli użytkownika w interfejsie.

| Mechanizm | Rozstrzygane pytanie | Rozdział |
|---|---|---|
| Uwierzytelnianie | Czy urządzenie uzyskuje dostęp do platformy jako właściciel konta? | 2–9 |
| Model danych tożsamości i kontrakt komunikacji uwierzytelniania | Jak zapisywana i wymieniana jest tożsamość oraz stan dostępu? | 10–12 |
| Autoryzacja w kanałach komunikacji operacyjnej | Jakie działania wolno wykonać Koordynatorowi i Wykonawcy oraz co wymaga zatwierdzenia Użytkownika? | 13 |
| Warstwy widoczności | Które funkcje platformy są ujawniane danej roli użytkownika? | 14 |

Trzy mechanizmy działają kolejno: bez uwierzytelnienia nie powstaje sesja pracy, bez sesji pracy nie działa żaden z kanałów komunikacji operacyjnej, a zakres funkcji ujawnianych w interfejsie wynika z roli przypisanej uwierzytelnionemu użytkownikowi. Materiał uwierzytelniający — skrót hasła, materiał metod dodatkowych i tokeny dostępu — przechowywany jest poza bazą danych strukturalnych, w wydzielonym magazynie danych dostępowych (rozdz. 9). Każde zlecenie, zadanie, komunikat sterujący i decyzja o ponowieniu podlega rejestracji w dzienniku audytu (rozdz. 13.5).

---

## 1. Cel i zakres dokumentu

Celem dokumentu jest ustalenie pełnego modelu bezpieczeństwa platformy Danaco Console: uwierzytelniania jako mechanizmu kontroli dostępu do platformy (Architektura, rozdz. 15), autoryzacji w obu kanałach komunikacji operacyjnej oraz warstw widoczności ograniczających zestaw funkcji ujawnianych roli użytkownika. Zakres dokumentu obejmuje:

| Zakres dokumentu | Rozdział |
|---|---|
| Zasada nadrzędna uwierzytelniania jako mechanizmu kontroli dostępu do platformy oraz jej relacja do modelu jednego właściciela i wielu urządzeń | 2–3 |
| Rejestracja konta przy pierwszym uruchomieniu platformy, wraz z wymaganymi danymi i weryfikacją adresu e-mail | 4 |
| Logowanie loginem i hasłem | 5 |
| Trzy metody dodatkowe uwierzytelniania — PIN, Windows Hello, adres e-mail uwierzytelniający — oraz ich konfigurowalność z okna konfiguracji | 6 |
| Odzyskiwanie konta i ustawienie nowego hasła przy wykorzystaniu adresu e-mail uwierzytelniającego | 7 |
| Wydawanie tokenu dostępu po uwierzytelnieniu oraz jego rola w nawiązywaniu połączenia WebSocket | 8 |
| Zasada przechowywania danych dostępowych poza bazą danych | 9 |
| Model danych encji Konto i Urządzenie | 10 |
| Polecenia i zdarzenia komunikacji WebSocket właściwe uwierzytelnianiu | 11 |
| Miejsce ustawień uwierzytelniania w oknie konfiguracji | 12 |
| Uwierzytelnianie i autoryzacja w kanale Chat Window (Użytkownik ↔ Wykonawca) i w kanale Execution Loop Window (Koordynator ↔ Wykonawca), granice autonomii Wykonawcy, punkty obowiązkowego zatwierdzenia, dziennik audytu i izolacja danych między zleceniami | 13 |
| Warstwy widoczności jako mechanizm kontroli dostępu do funkcji, konfiguracja warstw według roli użytkownika, tryb administracyjny i zachowanie wyszukiwarki funkcji | 14 |

Dokument nie zawiera specyfikacji kryptograficznej ani wyboru konkretnych bibliotek i protokołów implementujących poszczególne metody uwierzytelniania. Rozstrzygnięcia implementacyjne podlegają zasadom architektonicznym ustalonym w rozdziale 1.3 dokumentu Architektura, w szczególności zasadzie piątej: konfigurowalność zamiast wartości wpisanych na stałe.

---

## 2. Zasada nadrzędna: uwierzytelnianie jako mechanizm kontroli dostępu do platformy

Dokument Architektura ustanawia wprost: „Uwierzytelnianie jest jedynym mechanizmem kontroli dostępu w systemie" (rozdz. 13). Zasada ta ma dla platformy dwie konsekwencje.

| Konsekwencja | Treść |
|---|---|
| 1. Jeden próg wejścia do platformy | Platforma nie bramkuje wejścia do środowisk, modułów ani funkcji globalnych odrębnymi listami kontroli dostępu. O dostępie do platformy rozstrzyga to, czy urządzenie przedstawiło ważny token dostępu wydany po uwierzytelnieniu (rozdz. 8). Nie istnieje pośredni stan dostępu ograniczonego zakresowo przez sam mechanizm uwierzytelniania. |
| 2. Odróżnienie od mechanizmów działających po uwierzytelnieniu | Mechanizmy odpowiadające na inne pytanie — autoryzacja w kanałach komunikacji operacyjnej (rozdz. 13), warstwy widoczności i konfiguracja roli (rozdz. 14), Permissions Center oraz izolacja techniczna „konto i token per sesja" — działają wewnątrz sesji już uwierzytelnionej i nie zastępują progu wejścia do platformy (zob. tabela rozgraniczenia poniżej). |

Po uzyskaniu dostępu użytkownik korzysta z zakresu platformy wyznaczonego przypisaną mu rolą — środowisk, modułów i funkcji globalnych ujawnianych zgodnie z warstwami widoczności (rozdz. 14) — przy zachowaniu zasady pełnej kompozycyjności (Koncepcja platformy, rozdz. 14, zasada 1) oraz modelu jednego właściciela opisanego w rozdziale 3.

**Rozgraniczenie mechanizmów.** Zasada odróżnia uwierzytelnianie od mechanizmów odpowiadających na inne pytanie.

| Mechanizm | Na jakie pytanie odpowiada | Czym jest | Czym nie jest | Źródło |
|---|---|---|---|---|
| Uwierzytelnianie | Czy urządzenie uzyskuje dostęp do platformy jako właściciel konta? | Mechanizm kontroli dostępu do platformy | — | Architektura, rozdz. 15 |
| Autoryzacja w kanałach komunikacji operacyjnej | Jakie działania wolno zlecić, wykonać i zatwierdzić w kanale Chat Window i w kanale Execution Loop Window? | Zestaw uprawnień Koordynatora, granic autonomii Wykonawcy i punktów zatwierdzenia przez Użytkownika wewnątrz sesji uwierzytelnionej | Progiem wejścia do platformy | rozdz. 13 |
| Warstwy widoczności i konfiguracja roli | Które funkcje platformy są ujawniane danej roli użytkownika w interfejsie? | Mechanizm kontroli dostępu do funkcji, oparty na czterech warstwach widoczności i konfiguracji warstw według roli | Progiem wejścia do platformy | rozdz. 14 |
| Permissions Center (moduł Agents) | Jaki zakres działań wolno wykonać konkretnemu agentowi jako wykonawcy zadań? | Konfiguracja uprawnień operacyjnych jednostki AI działającej w imieniu już uwierzytelnionego Użytkownika | Progiem wejścia do platformy | Koncepcja platformy, rozdz. 11.15 |
| Izolacja techniczna „konto i token per sesja" | Czy proces pojedynczej sesji otrzymuje własne, odrębne konto i token? | Zawężenie techniczne procesu pojedynczej sesji, sterowane ustawieniem konfiguracyjnym, domyślnie wyłączonym | Odrębnym mechanizmem kontroli dostępu do platformy | Architektura, rozdz. 12; Koncepcja platformy, rozdz. 6.4 |

Wymienione mechanizmy nie zastępują uwierzytelniania. Relację tokenu izolacji sesji do tokenu uwierzytelniania wyjaśnia rozdział 8.4.

---

## 3. Model właściciela: jedno konto, wiele urządzeń

Dokument Architektura ustala model wdrożenia platformy w rozdziale 2: „Jeden użytkownik (Operator). Wiele urządzeń: komputery, telefony, tablety. Instalacja pakietu klienta: jednorazowa na urządzenie." Uwierzytelnianie opisane w niniejszym dokumencie realizuje ten model wprost — platforma nie jest systemem wielokontowym, a rejestracja opisana w rozdziale 4 nie jest procesem otwierania kolejnych kont, lecz jednorazowym utworzeniem jedynego konta właściciela.

```
                    ┌─────────────────────────────────┐
                    │   KONTO — jedyne (właściciel)    │
                    │   Użytkownik  (rozdz. 10.1)      │
                    └────────────────┬────────────────┘
                                     │  jedno konto, wiele urządzeń
        ┌────────────────┬───────────┴───────────┬────────────────┐
        ▼                ▼                        ▼                ▼
 ┌────────────┐   ┌────────────┐          ┌────────────┐   ┌────────────┐
 │ Urządzenie │   │ Urządzenie │          │ Urządzenie │   │    …       │
 │ (komputer) │   │  (telefon) │          │  (tablet)  │   │ bez limitu │
 │  token A   │   │   token B  │          │   token C  │   │            │
 └────────────┘   └────────────┘          └────────────┘   └────────────┘
   Każde urządzenie ma własny token dostępu (rozdz. 8.1, 11.2)
```

Konsekwencją tego modelu jest rozróżnienie dwóch odrębnych czynności, które w systemie wielokontowym byłyby tożsame:

| Czynność | Kiedy występuje | Skutek |
|---|---|---|
| Rejestracja (rozdz. 4) | Wyłącznie przy pierwszym uruchomieniu platformy, gdy konto właściciela jeszcze nie istnieje | Utworzenie jedynej encji Konto (rozdz. 10.1) |
| Logowanie (rozdz. 5) | Przy każdym kolejnym uruchomieniu platformy oraz przy podłączaniu każdego kolejnego urządzenia do istniejącego konta | Uwierzytelnienie urządzenia wobec istniejącego konta i wydanie mu tokenu dostępu (rozdz. 8) |

Urządzenie przenośne uzyskuje własne poświadczenie w protokole parowania, opisanym w całości w dokumencie `funkcje-globalne/mobile.md` (rozdz. 9); poświadczenie to jest tokenem urządzenia w rozumieniu rozdziału 8.1, wydanym bez wprowadzania loginu i hasła na urządzeniu przenośnym.

Wykaz urządzeń powiązanych z kontem oraz zarządzanie ich dostępem jest ustawieniem okna konfiguracji (Model konfiguracji, rozdz. 5.1, pozycja „Urządzenia połączone"), opisanym szerzej w rozdziale 12. Platforma nie ogranicza domyślnie liczby urządzeń powiązanych z kontem.

---

## 4. Rejestracja

### 4.1. Moment rejestracji

Rejestracja odbywa się wyłącznie przy pierwszym uruchomieniu platformy (Architektura, rozdz. 15), zanim istnieje jakiekolwiek konto właściciela. Okno rejestracji jest pierwszym oknem prezentowanym urządzeniu łączącemu się z platformą, na której nie utworzono jeszcze konta. Poniższy schemat przedstawia pełny przebieg rejestracji.

```
Pierwsze uruchomienie platformy  (konto właściciela nie istnieje)
        │
        ▼
Okno rejestracji  ── pierwsze okno prezentowane urządzeniu (rozdz. 4.1)
        │
        ▼
Podanie danych wymaganych (rozdz. 4.2)
   login  ·  adres e-mail uwierzytelniający  ·  hasło
        │
        ▼
Weryfikacja adresu e-mail (rozdz. 4.3)
   konto w stanie NIEPOTWIERDZONYM do chwili potwierdzenia
        │
        ▼  potwierdzenie e-mail
Utworzenie jedynej encji KONTO (rozdz. 4.4, 11.1)
   metody aktywne domyślnie:   hasło  +  e-mail uwierzytelniający
   metody nieaktywne:          PIN, Windows Hello (do włączenia — rozdz. 6.4)
        │
        ▼
Wydanie tokenu dostępu urządzeniu (rozdz. 8.1)
```

### 4.2. Dane wymagane

| Pole | Opis | Rola w dalszym uwierzytelnianiu |
|---|---|---|
| Login | Nazwa właściciela konta wykorzystywana przy logowaniu | Wykorzystywany łącznie z hasłem przy logowaniu (rozdz. 5) |
| Adres e-mail uwierzytelniający | Adres poczty elektronicznej właściciela | Weryfikowany przy rejestracji (rozdz. 4.3); pełni dodatkowo funkcję metody logowania (rozdz. 6.3) oraz jedynej drogi odzyskania konta (rozdz. 7) |
| Hasło | Hasło dostępowe do konta | Wykorzystywane łącznie z loginem przy logowaniu (rozdz. 5); przechowywane poza bazą danych w postaci skrótu (rozdz. 9) |

### 4.3. Weryfikacja adresu e-mail

Rejestracja wymaga potwierdzenia adresu e-mail uwierzytelniającego (Architektura, rozdz. 15: „potwierdzenie przez e-mail"). Konto pozostaje w stanie niepotwierdzonym do chwili wykonania potwierdzenia; adres e-mail wykorzystywany do potwierdzenia jest tym samym adresem, który później pełni funkcję metody logowania (rozdz. 6.3) i odzyskiwania konta (rozdz. 7) — platforma nie wprowadza odrębnego, drugiego adresu do żadnego z tych celów.

### 4.4. Wynik rejestracji

Pozytywnie zakończona i potwierdzona rejestracja tworzy jedyną encję Konto (rozdz. 10.1) w konfiguracji metod ustalonej poniżej.

| Metoda uwierzytelniania | Stan po rejestracji | Podstawa |
|---|---|---|
| Hasło | Aktywna | Metoda bazowa ustalana przy rejestracji |
| Adres e-mail uwierzytelniający | Aktywna | Wartość domyślna ustawienia „Metody uwierzytelniania" (Model konfiguracji, rozdz. 5.1) |
| PIN | Nieaktywna | Wymaga świadomego włączenia przez Użytkownika (rozdz. 6.4) |
| Windows Hello | Nieaktywna | Wymaga świadomego włączenia przez Użytkownika (rozdz. 6.4) |

---

## 5. Logowanie

Logowanie odbywa się loginem i hasłem (Architektura, rozdz. 15) i jest czynnością wykonywaną przy każdym uruchomieniu platformy na urządzeniu, które nie posiada jeszcze ważnego tokenu dostępu, oraz przy podłączaniu kolejnego urządzenia do istniejącego konta właściciela (rozdz. 3).

```
Uruchomienie platformy na urządzeniu bez ważnego tokenu
        │
        ▼
Okno logowania ──► Użytkownik wybiera metodę uwierzytelnienia:
        │             • login + hasło  (metoda bazowa, zawsze dostępna)
        │             • aktywna metoda dodatkowa (rozdz. 6):
        │               PIN · Windows Hello · e-mail uwierzytelniający
        ▼
Weryfikacja przez serwer  (skrót hasła / materiał metody — rozdz. 9)
        │
        ├── pozytywna ──► wydanie tokenu dostępu (rozdz. 8.1)
        │                       │
        │                       ▼
        │                 nawiązanie / potwierdzenie połączenia WebSocket (rozdz. 8.2)
        │
        └── negatywna ──► ponowna próba  albo  odzyskiwanie konta (rozdz. 7)
```

| Krok | Opis |
|---|---|
| 1. Podanie danych | Użytkownik podaje login i hasło w oknie logowania |
| 2. Weryfikacja | Serwer porównuje przedstawione hasło ze skrótem hasła przechowywanym poza bazą danych (rozdz. 9) |
| 3. Wydanie tokenu | Po pozytywnej weryfikacji serwer wydaje urządzeniu token dostępu (rozdz. 8.1) |
| 4. Nawiązanie połączenia | Token wykorzystywany jest do nawiązania lub potwierdzenia połączenia WebSocket (rozdz. 8.2) |

Logowanie pozostaje dostępne również wtedy, gdy Użytkownik skonfigurował metody dodatkowe (rozdz. 6) — hasło jest metodą bazową, obecną zawsze, natomiast metody dodatkowe rozszerzają, a nie zastępują, sposób uwierzytelnienia.

---

## 6. Metody dodatkowe uwierzytelniania

Poza logowaniem loginem i hasłem platforma udostępnia trzy metody dodatkowe, konfigurowane z okna konfiguracji (Architektura, rozdz. 15). Żadna z metod dodatkowych nie jest aktywna z konieczności technicznej — zgodnie z zasadą pełnej konfigurowalności (Koncepcja platformy, rozdz. 14, zasada 2) każda z nich jest włączana i wyłączana świadomą decyzją Użytkownika.

| Metoda | Charakter | Przeznaczenie | Przechowywanie materiału (rozdz. 9) | Aktywna domyślnie |
|---|---|---|---|---|
| Hasło | Bazowa, obecna zawsze | Uwierzytelnienie pełne loginem i hasłem (rozdz. 5) | Skrót hasła poza bazą danych | Tak |
| PIN | Dodatkowa | Szybkie potwierdzenie tożsamości na urządzeniu już uwierzytelnionym pełnym loginem i hasłem | Materiał metody PIN poza bazą danych | Nie |
| Windows Hello | Dodatkowa, wiązana z konkretnym urządzeniem | Uwierzytelnienie oparte na mechanizmie systemowym Windows Hello | Materiał metody Windows Hello poza bazą danych | Nie |
| Adres e-mail uwierzytelniający | Dodatkowa (weryfikacja · logowanie · odzyskiwanie) | Samodzielna metoda logowania alternatywna wobec loginu i hasła (rozdz. 6.3) | Adres w encji Konto (rozdz. 10.1); nie jest materiałem tajnym | Tak |

### 6.1. PIN

Krótki kod numeryczny, alternatywny wobec hasła, przeznaczony do szybkiego potwierdzenia tożsamości na urządzeniu już wcześniej uwierzytelnionym pełnym loginem i hasłem. Materiał PIN przechowywany jest poza bazą danych, na zasadach opisanych w rozdziale 9.

### 6.2. Windows Hello

Metoda uwierzytelniania oparta na mechanizmie systemowym Windows Hello, wiążąca uwierzytelnienie z konkretnym urządzeniem Użytkownika. Materiał uwierzytelniający właściwy tej metodzie przechowywany jest poza bazą danych, analogicznie do pozostałych metod (rozdz. 9).

### 6.3. Adres e-mail uwierzytelniający jako metoda logowania

Ten sam adres e-mail, który przy rejestracji pełni funkcję weryfikacyjną (rozdz. 4.3), a przy odzyskiwaniu konta funkcję jedynej drogi dostępu (rozdz. 7), może zostać skonfigurowany jako samodzielna metoda logowania — alternatywna wobec podania loginu i hasła. Zgodnie z wartością domyślną ustawienia „Metody uwierzytelniania" (Model konfiguracji, rozdz. 5.1) metoda ta jest aktywna od chwili rejestracji.

### 6.4. Konfigurowalność zestawu metod

| Ustawienie | Warstwa | Wartość domyślna | Źródło |
|---|---|---|---|
| Metody uwierzytelniania | globalna | hasło i e-mail uwierzytelniający | Model konfiguracji, rozdz. 5.1 |

Szablon konfiguracji ustawienia:

```
Ustawienie: „Metody uwierzytelniania"
  Zakres:            Aplikacja (rozdz. 12)
  Warstwa:           globalna
  Źródło:            Model konfiguracji, rozdz. 5.1
  Wartość domyślna:  hasło + e-mail uwierzytelniający
  Metody włączalne:  hasło (stała) · e-mail uwierzytelniający · PIN · Windows Hello
  Objaśnienie [?]:   opis działania ustawienia i jego wpływu na aplikację
  Zasada:            brak twardych blokad — wyłączenie metody dodatkowej
                     nie ogranicza dostępu; pozostaje logowanie metodą
                     bazową (hasło) oraz metodą domyślnie aktywną
                     (e-mail uwierzytelniający)
```

Włączenie lub wyłączenie metody PIN albo Windows Hello nie wpływa na pozostałe aktywne metody — Użytkownik może w dowolnym momencie logować się dowolną z aktywnych metod, zgodnie z zasadą braku twardych blokad: brak skonfigurowania metody dodatkowej nie ogranicza dostępu, a jedynie pozostawia logowanie przy metodzie bazowej (hasło) oraz metodzie domyślnie aktywnej (e-mail uwierzytelniający).

---

## 7. Odzyskiwanie konta

Adres e-mail uwierzytelniający pełni funkcję odzyskiwania konta i ustawienia nowego hasła (Architektura, rozdz. 15, dosłownie). Jest to jedyna droga odzyskania dostępu do konta przewidziana przez platformę — spójna z tym, że adres e-mail jest jedynym elementem tożsamości Użytkownika niezależnym od konkretnego urządzenia lub materiału uwierzytelniającego przechowywanego lokalnie (PIN, Windows Hello).

```
Utrata dostępu do dotychczasowego materiału uwierzytelniającego
   (zapomniane hasło; brak urządzenia z aktywnym tokenem)
        │
        ▼
Okno logowania ──► funkcja odzyskiwania konta  (krok 1)
        │            Użytkownik podaje adres e-mail uwierzytelniający
        ▼
Serwer wysyła na wskazany adres drogę potwierdzenia tożsamości  (krok 2)
        │            analogicznie do weryfikacji przy rejestracji (rozdz. 4.3)
        ▼  potwierdzenie
Ustawienie nowego hasła  (krok 3)
        │            zastępuje dotychczasowy skrót hasła poza bazą danych (rozdz. 9)
        ▼
Unieważnienie aktywnych tokenów wydanych przed odzyskaniem  (krok 4)
        │            → ponowne logowanie na wszystkich powiązanych urządzeniach
        │            → zdarzenie device.changed rozgłoszone (rozdz. 11.2)
        ▼
Ta sama, jedyna encja KONTO (rozdz. 3, 11.1)
   zmienia się wyłącznie materiał uwierzytelniający — nowa encja nie powstaje
```

| Krok | Opis |
|---|---|
| 1. Inicjacja | Użytkownik wskazuje w oknie logowania funkcję odzyskiwania konta i podaje adres e-mail uwierzytelniający |
| 2. Weryfikacja adresu | Serwer wysyła na wskazany adres wiadomość z drogą potwierdzenia tożsamości (analogicznie do weryfikacji przy rejestracji, rozdz. 4.3) |
| 3. Ustawienie nowego hasła | Po potwierdzeniu Użytkownik ustawia nowe hasło, zastępujące dotychczasowy skrót hasła przechowywany poza bazą danych (rozdz. 9) |
| 4. Unieważnienie tokenów | Odzyskanie konta unieważnia aktywne tokeny dostępu wydane przed odzyskaniem, zgodnie z zarządzaniem dostępem urządzeń (rozdz. 3, Model konfiguracji rozdz. 5.1) |

Odzyskiwanie konta nie tworzy nowej encji Konto — działa na tej samej, jedynej encji właściciela (rozdz. 3, rozdz. 10.1), zmieniając wyłącznie materiał uwierzytelniający.

---

## 8. Token dostępu i połączenie WebSocket

### 8.1. Wydanie tokenu

Po uwierzytelnieniu — czy to przy rejestracji (rozdz. 4), czy przy logowaniu (rozdz. 5), czy po odzyskaniu konta (rozdz. 7) — urządzenie otrzymuje token dostępu wykorzystywany do nawiązania połączenia WebSocket (Architektura, rozdz. 15).

| Właściwość tokenu dostępu | Ustalenie |
|---|---|
| Moment wydania | Po uwierzytelnieniu: przy rejestracji (rozdz. 4), logowaniu (rozdz. 5) lub po odzyskaniu konta (rozdz. 7) |
| Przypisanie | Do encji Urządzenie (rozdz. 10.2), nie do pojedynczego połączenia |
| Trwałość | Ważny między kolejnymi połączeniami tego samego urządzenia, aż do unieważnienia |
| Unieważnienie | Odzyskanie konta (rozdz. 7, krok 4), unieważnienie z okna konfiguracji (rozdz. 3, rozdz. 12) lub odwołanie dostępu urządzenia przenośnego (`funkcje-globalne/mobile.md`, rozdz. 9.8) |
| Przechowywanie | Poza bazą danych, w magazynie danych dostępowych (rozdz. 9) |

### 8.2. Token w cyklu życia połączenia

Kontrakty komunikacji opisują cykl życia połączenia WebSocket w pięciu krokach (rozdz. 2), z których drugi jest właściwy uwierzytelnianiu.

| Krok | Opis | Relacja do niniejszego dokumentu |
|---|---|---|
| 1. Nawiązanie połączenia | Klient otwiera kanał WebSocket z serwerem | Połączenie wstępne, przed uwierzytelnieniem |
| 2. Uwierzytelnienie | Klient przedstawia token dostępu uzyskany po zalogowaniu; bez ważnego tokenu połączenie nie przechodzi do kroku „Powiązanie" | Token wydany zgodnie z rozdz. 8.1 |
| 3. Powiązanie | Klient zgłasza aktywne środowisko i moduł; serwer wiąże połączenie z sesją | Poza zakresem niniejszego dokumentu (Koncepcja platformy, rozdz. 8) |
| 4. Praca | Wymiana poleceń, zdarzeń i strumieni | Poza zakresem niniejszego dokumentu |
| 5. Rozłączenie | Zamknięcie kanału nie kończy sesji po stronie serwera | Token pozostaje ważny; ponowne połączenie przedstawia ten sam token |

### 8.3. Pierwsze połączenie a połączenia kolejne

Krok „Uwierzytelnienie" zakłada, że klient przedstawia token już uzyskany po zalogowaniu (rozdz. 8.2). Dla urządzenia łączącego się po raz pierwszy — przy rejestracji lub przy pierwszym logowaniu na nowym urządzeniu — token jeszcze nie istnieje. Rejestracja i logowanie odbywają się tym samym, jedynym kanałem WebSocket co pozostała komunikacja platformy (Architektura, rozdz. 11; Kontrakty komunikacji, rozdz. 1).

```
PIERWSZE POŁĄCZENIE URZĄDZENIA  (token jeszcze nie istnieje)
  Klient ──[1] nawiązanie połączenia WebSocket ─────────────────► Serwer
  Klient ──[2] krok „Uwierzytelnienie" BEZ tokenu ──────────────► Serwer
  Klient ──[ ] auth.register / auth.login (rozdz. 11) ──────────► Serwer
  Klient ◄──────── token dostępu w odpowiedzi (po weryfikacji) ── Serwer

POŁĄCZENIA KOLEJNE  (token uzyskany wcześniej)
  Klient ──[1] nawiązanie połączenia WebSocket ─────────────────► Serwer
  Klient ──[2] krok „Uwierzytelnienie" z tokenem ───────────────► Serwer
  Klient ◄──────── przejście do kroku „Powiązanie" (rozdz. 8.2) ─ Serwer
```

Pierwsze połączenie urządzenia przechodzi krok uwierzytelnienia bez tokenu, klient przesyła polecenie rejestracji albo logowania (rozdz. 11), a serwer, po pozytywnej weryfikacji, wydaje token dostępu w odpowiedzi. Każde kolejne połączenie tego samego lub innego już zarejestrowanego urządzenia przedstawia w kroku uwierzytelnienia token uzyskany wcześniej.

### 8.4. Relacja do izolacji technicznej „konto i token per sesja"

Token opisany w niniejszym rozdziale jest tokenem urządzenia — wydawanym raz na urządzenie i wykorzystywanym do uwierzytelnienia każdego jego połączenia z platformą (rozdz. 8.1). Jest to mechanizm odrębny od zakresu izolacji technicznej „konto i token per sesja", wymienionego w rozdziale 9 dokumentu Architektura oraz w macierzy izolacji technicznej okna konfiguracji punktów izolacji (Koncepcja platformy, rozdz. 6.4).

| Cecha | Token urządzenia (rozdz. 8.1) | „Konto i token per sesja" (izolacja techniczna) |
|---|---|---|
| Charakter | Mechanizm uwierzytelniania połączenia | Zawężenie techniczne procesu sesji, sterowane ustawieniem konfiguracyjnym, domyślnie wyłączonym |
| Zasięg | Jeden na urządzenie | Proces pojedynczej sesji |
| Domyślnie | Zawsze obecny po uwierzytelnieniu | Wyłączony |
| Konfiguracja | Nie dotyczy — element uwierzytelniania | Okno konfiguracji punktów izolacji (Koncepcja platformy, rozdz. 6.4), np. dla roli środowiska MultitaskingAI (rozdz. 4.5) |
| Relacja wzajemna | Nadrzędny — uwierzytelnia samo połączenie | Zagnieżdżony — dodatkowy środek izolacji procesu, niezależny od tokenu urządzenia |

Gdy Użytkownik włącza izolację „konto i token per sesja" dla wybranego zasięgu — na przykład dla pojedynczej roli środowiska MultitaskingAI (Koncepcja platformy, rozdz. 6.5) — proces sesji otrzymuje własne, odrębne konto i token wykorzystywane wewnętrznie przez ten proces, niezależnie od tokenu urządzenia, które nawiązało połączenie WebSocket. Token urządzenia pozostaje w takim wypadku mechanizmem nadrzędnym, uwierzytelniającym samo połączenie; token sesji jest zagnieżdżonym środkiem izolacji procesu, włączanym ustawieniem konfiguracyjnym punktów izolacji.

---

## 9. Przechowywanie danych dostępowych poza bazą danych

Model danych ustala dla danych dostępowych kanału modelu zasadę: „Dane dostępowe przechowywane są poza bazą, zgodnie z zasadami bezpieczeństwa. Baza przechowuje wyłącznie odwołanie" (rozdz. 8.1). Niniejszy dokument stosuje tę samą zasadę do danych dostępowych uwierzytelniania — materiału stanowiącego bezpośredni dowód tożsamości Użytkownika nie przechowuje się w bazie SQLite będącej ogólnym magazynem danych strukturalnych platformy (Architektura, rozdz. 10), lecz w wydzielonym magazynie danych dostępowych po stronie serwera, poza zasięgiem odczytu właściwym operacjom na encjach modelu danych (Model danych, rozdz. 13).

```
                         SERWER PLATFORMY
 ┌───────────────────────────────┐        ┌──────────────────────────────────┐
 │  Baza SQLite                  │        │  Magazyn danych dostępowych      │
 │  dane strukturalne            │        │  poza bazą (Model danych 8.1)    │
 │  (Architektura, rozdz. 10)     │        │                                  │
 │                               │        │   • skrót hasła                  │
 │  encja Konto (11.1):          │─odw.──►│   • materiał metody PIN          │
 │    login, e-mail,             │        │   • materiał Windows Hello       │
 │    aktywne metody, daty        │        │   • token dostępu urządzenia     │
 │                               │        │   • dane dostępowe kanału modelu │
 │  encja Urządzenie (11.2):     │─odw.──►│                                  │
 │    identyfikator, nazwa, data  │        │  Kontrakt dostępu — odrębny od   │
 │                               │        │  operacji na encjach modelu      │
 │  Baza przechowuje             │        │  (Architektura, rozdz. 1.3,      │
 │  WYŁĄCZNIE odwołania          │        │   zasada 4)                      │
 └───────────────────────────────┘        └──────────────────────────────────┘
```

| Rodzaj danych | Miejsce przechowywania | Uwaga |
|---|---|---|
| Login, adres e-mail uwierzytelniający, aktywne metody uwierzytelniania, daty | Baza SQLite (encja Konto, rozdz. 10.1) | Dane nieujawniające materiału uwierzytelniającego |
| Skrót hasła | Poza bazą danych, w magazynie danych dostępowych | Baza przechowuje wyłącznie odwołanie do wpisu magazynu |
| Materiał metody PIN | Poza bazą danych, w magazynie danych dostępowych | Aktywna wyłącznie, gdy metoda PIN jest włączona (rozdz. 6.1, 6.4) |
| Materiał metody Windows Hello | Poza bazą danych, w magazynie danych dostępowych | Aktywna wyłącznie, gdy metoda Windows Hello jest włączona (rozdz. 6.2, 6.4) |
| Token dostępu urządzenia | Poza bazą danych, w magazynie danych dostępowych | Odwołanie do urządzenia (rozdz. 10.2) pozostaje w bazie |
| Dane dostępowe kanału modelu | Poza bazą danych, w magazynie danych dostępowych | Zasada już ustalona w Modelu danych, rozdz. 8.1 |

Struktura pojedynczego wpisu magazynu danych dostępowych:

```
Wpis magazynu danych dostępowych
  Odwołanie:       identyfikator wpisu (przechowywany w bazie SQLite)
  Rodzaj:          skrót hasła | materiał PIN | materiał Windows Hello
                   | token dostępu urządzenia | dane dostępowe kanału modelu
  Materiał:        wartość poza bazą danych (nie replikowana do bazy SQLite)
  Właściciel:      odwołanie do encji Konto (11.1) lub Urządzenie (11.2)
  Dostęp:          wyłącznie przez kontrakt dostępu (rozdz. 1.3 Architektury, zasada 4)
```

Wydzielenie magazynu danych dostępowych poza bazę danych strukturalnych jest spójne z zasadą architektoniczną rozdzielenia przez kontrakty (Architektura, rozdz. 1.3, zasada 4): operacje na danych dostępowych realizowane są przez odrębny kontrakt dostępu, nie przez ogólny interfejs operacji na encjach modelu danych opisany w rozdziale 10 Architektury.

---

## 10. Model danych tożsamości i dostępu

Rozdział 2 dokumentu Model danych opisuje dwie encje właściwe tożsamości i dostępowi. Niniejszy rozdział przywołuje je w pełnym brzmieniu i doprecyzowuje atrybut „skrót hasła" zgodnie z zasadą przechowywania poza bazą ustaloną w rozdziale 9.

```
┌───────────────────────────────────┐          ┌───────────────────────────────────┐
│  KONTO  (rozdz. 10.1)             │ 1      N │  URZĄDZENIE  (rozdz. 10.2)        │
│  jedyne — właściciel (Użytkownik)   ├──────────┤  urządzenie połączone z platformą │
│                                   │ posiada  │                                   │
│  login                            │          │  identyfikator urządzenia         │
│  adres e-mail uwierzytelniający   │          │  nazwa                            │
│  → skrót hasła (poza bazą)        │          │  → token dostępu (poza bazą)      │
│  aktywne metody uwierzytelniania  │          │  data ostatniego połączenia       │
│  data utworzenia                  │          │                                   │
└───────────────────────────────────┘          └───────────────────────────────────┘
   Wykaz encji Urządzenie = ustawienie „Urządzenia połączone" (rozdz. 3, rozdz. 12)
```

### 10.1. Konto

Jedyne konto właściciela (Użytkownik), zgodnie z modelem opisanym w rozdziale 3.

| Atrybut | Opis | Przechowywanie |
|---|---|---|
| login | Nazwa właściciela konta wykorzystywana przy logowaniu | Baza SQLite |
| adres e-mail uwierzytelniający | Adres pełniący funkcje weryfikacji, logowania i odzyskiwania (rozdz. 4.3, 6.3, 7) | Baza SQLite |
| odwołanie do skrótu hasła | Wskazanie wpisu magazynu danych dostępowych (rozdz. 9) | Baza SQLite (odwołanie); skrót hasła poza bazą danych |
| aktywne metody uwierzytelniania | Które metody Użytkownik włączył: hasło, PIN, Windows Hello, e-mail uwierzytelniający | Baza SQLite (stan ustawienia, bez materiału — rozdz. 6.4, 13) |
| data utworzenia | Data utworzenia konta | Baza SQLite |

Szablon encji:

```
Encja: Konto  (jedyna — właściciel)
  login:                          nazwa właściciela (logowanie)
  adres e-mail uwierzytelniający: weryfikacja (4.3) · logowanie (6.3) · odzyskiwanie (7)
  odwołanie do skrótu hasła:      → magazyn danych dostępowych (rozdz. 9)
  aktywne metody uwierzytelniania: hasło | PIN | Windows Hello | e-mail uwierzytelniający
  data utworzenia:                data utworzenia konta
```

Atrybut „aktywne metody uwierzytelniania" jest zapisem stanu ustawienia „Metody uwierzytelniania" opisanego w rozdziale 6.4 i w rozdziale 12 — nie zawiera samego materiału uwierzytelniającego, wyłącznie informację o tym, które metody Użytkownik włączył.

### 10.2. Urządzenie

Urządzenie połączone z platformą.

| Atrybut | Opis | Przechowywanie |
|---|---|---|
| identyfikator urządzenia | Jednoznaczny identyfikator urządzenia | Baza SQLite |
| nazwa | Nazwa urządzenia | Baza SQLite |
| odwołanie do tokenu dostępu | Wskazanie wpisu magazynu danych dostępowych (rozdz. 9) | Baza SQLite (odwołanie); token dostępu poza bazą danych |
| data ostatniego połączenia | Data ostatniego połączenia urządzenia z platformą | Baza SQLite |

Szablon encji:

```
Encja: Urządzenie  (urządzenie połączone; wiele na jedno Konto)
  identyfikator urządzenia:    jednoznaczny identyfikator
  nazwa:                       nazwa urządzenia
  odwołanie do tokenu dostępu: → magazyn danych dostępowych (rozdz. 9)
  data ostatniego połączenia:  data ostatniego połączenia
```

Wykaz encji Urządzenie powiązanych z jedynym Kontem jest zawartością ustawienia „Urządzenia połączone" (rozdz. 3, rozdz. 12).

---

## 11. Komunikacja — polecenia i zdarzenia uwierzytelniania

Kontrakt komunikacji obejmuje odrębną kategorię poleceń i zdarzeń właściwą rejestracji, logowaniu, konfiguracji metod uwierzytelniania, odzyskiwaniu konta i zarządzaniu urządzeniami. Kategoria stosuje konwencję nazewniczą (`obszar.czynność`) i strukturę koperty komunikatu obowiązujące pozostałe obszary operacji platformy (Kontrakty komunikacji, rozdz. 3), a uwierzytelnienie zajmuje drugi krok cyklu życia połączenia (Kontrakty komunikacji, rozdz. 2).

### 11.1. Polecenia (klient do serwera)

| Polecenie | Opis |
|---|---|
| `auth.register` | Rejestracja jedynego konta właściciela — login, adres e-mail uwierzytelniający, hasło (rozdz. 4) |
| `auth.verify` | Potwierdzenie adresu e-mail uwierzytelniającego po rejestracji (rozdz. 4.3) |
| `auth.login` | Logowanie loginem i hasłem albo aktywną metodą dodatkową (rozdz. 5, rozdz. 6) |
| `auth.method.set` | Włączenie albo wyłączenie metody dodatkowej uwierzytelniania w oknie konfiguracji (rozdz. 6.4) |
| `auth.recover` | Inicjacja odzyskiwania konta adresem e-mail uwierzytelniającym (rozdz. 7) |
| `auth.reset` | Ustawienie nowego hasła po potwierdzeniu odzyskania konta (rozdz. 7) |
| `device.list` | Wykaz urządzeń połączonych z kontem (rozdz. 3) |
| `device.revoke` | Unieważnienie tokenu dostępu wskazanego urządzenia (rozdz. 3, rozdz. 7) |

### 11.2. Zdarzenia (serwer do klienta)

| Zdarzenie | Opis |
|---|---|
| `auth.changed` | Zmiana stanu uwierzytelnienia konta — dodanie lub usunięcie metody, potwierdzenie adresu e-mail, zmiana hasła |
| `device.changed` | Zmiana wykazu urządzeń połączonych z kontem, w tym unieważnienie tokenu |

Zgodnie z ogólną zasadą kontraktów komunikacji (Kontrakty komunikacji, rozdz. 1) każde z powyższych zdarzeń jest rozgłaszane przez serwer do wszystkich połączonych urządzeń — dzięki temu unieważnienie tokenu jednego urządzenia (rozdz. 3) albo zmiana zestawu aktywnych metod uwierzytelniania (rozdz. 6.4) jest natychmiast widoczna na pozostałych urządzeniach powiązanych z kontem. Polecenia i zdarzenia kanałów komunikacji operacyjnej, obejmujące zlecenia, zadania, komunikaty sterujące i decyzje o zatwierdzeniu, opisuje rozdział 13.6.

---

## 12. Ustawienia uwierzytelniania w oknie konfiguracji

Ustawienia uwierzytelniania należą do zakresu „Aplikacja" okna konfiguracji (Model konfiguracji, rozdz. 5.1) — zakresu obejmującego ustawienia platformy jako całości, niezależne od konkretnego środowiska, modułu lub sesji.

| Ustawienie | Opis | Warstwa | Wartość domyślna |
|---|---|---|---|
| Metody uwierzytelniania | Aktywne metody logowania: hasło, PIN, Windows Hello, e-mail uwierzytelniający (rozdz. 6) | globalna | hasło i e-mail uwierzytelniający |
| Urządzenia połączone | Wykaz urządzeń powiązanych z kontem oraz zarządzanie ich dostępem, w tym unieważnienie tokenu (rozdz. 3, rozdz. 8, rozdz. 11.1) | globalna | brak ograniczeń liczby urządzeń |

Szablon konfiguracji ustawienia „Urządzenia połączone":

```
Ustawienie: „Urządzenia połączone"
  Zakres:            Aplikacja (rozdz. 12)
  Warstwa:           globalna
  Zawartość:         wykaz encji Urządzenie powiązanych z Kontem (rozdz. 10.2)
  Pola pozycji:      identyfikator · nazwa · data ostatniego połączenia
  Operacje:          device.list (wykaz) · device.revoke (unieważnienie tokenu)
                     — rozdz. 11.1
  Wartość domyślna:  brak ograniczeń liczby urządzeń
  Objaśnienie [?]:   opis działania ustawienia i jego wpływu na aplikację
```

Każde z powyższych ustawień, zgodnie z mechanizmem objaśnień kontekstowych obowiązującym w całym oknie konfiguracji (Architektura, rozdz. 13; Model konfiguracji, rozdz. 3.3), jest opatrzone objaśnieniem `[?]` opisującym jego działanie i wpływ na aplikację. Ustawienia autoryzacji kanałów komunikacji operacyjnej opisuje rozdział 13.7, a ustawienie konfiguracji warstw widoczności według roli — rozdział 14.3.

---

## 13. Bezpieczeństwo kanałów komunikacji operacyjnej

Praca operacyjna platformy przebiega w dwóch kanałach komunikacji: Chat Window prowadzi wymianę Użytkownik ↔ Wykonawca, a Execution Loop Window prowadzi wymianę Koordynator ↔ Wykonawca. Oba kanały są elementami pierwszoplanowymi architektury i podlegają jednolitym regułom uwierzytelniania i autoryzacji. Kanały działają wyłącznie wewnątrz sesji uwierzytelnionej: token dostępu urządzenia (rozdz. 8.1) uwierzytelnia połączenie WebSocket przenoszące komunikację obu kanałów, a każde polecenie, zlecenie i komunikat sterujący jest wiązany z kontem właściciela, urządzeniem i identyfikatorem zlecenia.

Rozmieszczenie obu kanałów w układzie pionowym obszaru roboczego, wraz z punktami sterowania bezpieczeństwem:

```
 ═══════════════════════════════════════════════════════════════════════════
  Boczna    │ Chat Window          │ Obszar roboczy modułu    │ Panel
  nawigacja │ Użytkownik ↔         │                          │ pomocniczy
  modułów   │ Wykonawca            │  wynik zadania           │ (rozszerzenie
            │                      │  przedstawiony do        │  boczne)
            │ • polecenie          │  zatwierdzenia           │
            │ • zatwierdzenie      │                          │ • dziennik
            │ • przerwanie         │                          │   audytu
            │ ─────────────────    │                          │   (rozdz. 13.5)
            │ Execution Loop       │                          │
            │ Koordynator ↔        │                          │
            │ Wykonawca            │                          │
            │ • zlecenie i zadania │                          │
            │ • komunikaty ster.   │                          │
            │ • decyzje ponowienia │                          │
 ═══════════════════════════════════════════════════════════════════════════
```

### 13.1. Role i ich uprawnienia

| Rola | Charakter | Uprawnienia | Ograniczenia |
|---|---|---|---|
| Użytkownik | Osoba uwierzytelniona wobec konta właściciela (rozdz. 3) | Wydaje zlecenia w Chat Window, zatwierdza wyniki i działania nieodwracalne, przerywa przebieg, zmienia granice autonomii Wykonawcy i konfigurację warstw widoczności (rozdz. 14.3) | Uprawnienia obowiązują wyłącznie po uwierzytelnieniu urządzenia (rozdz. 8) |
| Koordynator | Komponent orkiestrujący platformy | Dekomponuje zlecenie Użytkownika na zadania, przydziela zadania Wykonawcom, nadzoruje realizację, ocenia wynik kontroli jakości, podejmuje decyzję o ponowieniu zadania, wstrzymuje i wznawia pętlę wykonawczą | Działa wyłącznie w granicach zlecenia wydanego przez Użytkownika; nie tworzy zleceń z własnej inicjatywy, nie rozszerza zakresu zlecenia, nie zatwierdza działań zastrzeżonych dla Użytkownika (rozdz. 13.4) |
| Wykonawca | AI, agent lub system wykonawczy | Realizuje przydzielone zadanie, korzysta z narzędzi i zasobów w zakresie wyznaczonym granicami autonomii (rozdz. 13.3), zgłasza wynik i stan zadania Koordynatorowi | Nie przydziela zadań innym Wykonawcom, nie sięga po dane innego zlecenia (rozdz. 13.6), nie wykonuje działań wymagających zatwierdzenia bez uzyskanej decyzji Użytkownika |

Uprawnienia Koordynatora do przydzielania zadań wynikają z zakresu zlecenia: Koordynator przydziela wyłącznie zadania będące dekompozycją bieżącego zlecenia, wyłącznie Wykonawcom uprawnionym do danego typu zadania i wyłącznie w obrębie zasobów objętych zleceniem. Zadanie wykraczające poza ten zakres nie jest przydzielane — Koordynator kieruje wniosek o rozszerzenie zlecenia do Użytkownika przez Chat Window.

### 13.2. Uwierzytelnianie i autoryzacja w kanale Użytkownik ↔ Wykonawca

| Element | Zasada |
|---|---|
| Uwierzytelnienie kanału | Chat Window działa na połączeniu WebSocket uwierzytelnionym tokenem urządzenia (rozdz. 8.2); bez ważnego tokenu okno nie przyjmuje poleceń |
| Tożsamość nadawcy | Każde polecenie wprowadzone w Chat Window jest wiązane z kontem właściciela i identyfikatorem urządzenia (rozdz. 10.2) |
| Zakres polecenia | Polecenie w języku naturalnym wykonywane jest w granicach roli przypisanej użytkownikowi; funkcje niedostępne roli nie są uruchamiane poleceniem języka naturalnego (rozdz. 14.4) |
| Zatwierdzanie i przerywanie | Chat Window jest miejscem, w którym Użytkownik zatwierdza działania wymagające zatwierdzenia (rozdz. 13.4) oraz przerywa przebieg |
| Prezentacja wyniku | Wynik zadania przedstawiany jest wraz z informacją o zakresie danych, z których powstał, i o zleceniu, w ramach którego został wytworzony |

### 13.3. Autoryzacja w kanale Koordynator ↔ Wykonawca i granice autonomii Wykonawcy

Execution Loop Window prowadzi pętlę wykonawczą: zlecenie i jego dekompozycję na zadania, kolejkę i stan zadań, wymianę komunikatów sterujących, wyniki kontroli jakości, decyzje o ponowieniu oraz sterowanie przebiegiem. Komunikacja tego kanału podlega tej samej sesji uwierzytelnionej co Chat Window i nie stanowi odrębnego wejścia do platformy.

Granice autonomii Wykonawcy wyznaczają, które działania Wykonawca realizuje samodzielnie, a które wymagają decyzji nadrzędnej.

| Zakres działania | Realizacja |
|---|---|
| Odczyt danych objętych zleceniem, analiza, przygotowanie wyniku roboczego | Wykonawca realizuje samodzielnie w granicach zadania |
| Ponowienie zadania po negatywnym wyniku kontroli jakości | Koordynator podejmuje decyzję o ponowieniu; liczba ponowień wynika z ustawienia konfiguracyjnego (rozdz. 13.7) |
| Zapis zmian trwałych, operacje na zasobach zewnętrznych, wysłanie treści poza platformę, operacje nieodwracalne | Wymaga zatwierdzenia przez Użytkownika (rozdz. 13.4) |
| Rozszerzenie zakresu zlecenia, sięgnięcie po zasób spoza zlecenia | Niedostępne dla Wykonawcy i dla Koordynatora; wymaga nowego lub zmienionego zlecenia Użytkownika |
| Zmiana konfiguracji platformy, uprawnień roli i warstw widoczności | Zastrzeżona dla Użytkownika w trybie administracyjnym (rozdz. 14.2) |

Przekroczenie granicy autonomii zatrzymuje zadanie w stanie oczekiwania na decyzję i odnotowuje zdarzenie w dzienniku audytu (rozdz. 13.5).

### 13.4. Punkty obowiązkowego zatwierdzenia przez Użytkownika

| Punkt zatwierdzenia | Moment | Kanał zatwierdzenia |
|---|---|---|
| Przyjęcie zlecenia do realizacji | Po dekompozycji zlecenia na zadania, przed uruchomieniem pętli wykonawczej | Chat Window |
| Działanie nieodwracalne | Przed zapisem zmian trwałych, operacją na zasobach zewnętrznych i wysłaniem treści poza platformę | Chat Window |
| Rozszerzenie zakresu zlecenia | Przed sięgnięciem po zasób nieobjęty zleceniem | Chat Window |
| Przekroczenie ustalonej liczby ponowień zadania | Po wyczerpaniu liczby ponowień wynikającej z ustawienia konfiguracyjnego | Chat Window, na podstawie stanu przedstawionego w Execution Loop Window |
| Zamknięcie zlecenia | Po zrealizowaniu wszystkich zadań i pozytywnej kontroli jakości | Chat Window |

Decyzja Użytkownika jest jawna: zatwierdzenie ma postać czynności wykonanej w Chat Window i zostaje odnotowane w dzienniku audytu wraz z jego przedmiotem i chwilą wykonania. Brak decyzji utrzymuje zadanie w stanie oczekiwania; platforma nie zastępuje decyzji Użytkownika zachowaniem domyślnym.

### 13.5. Dziennik audytu

Dziennik audytu rejestruje przebieg pracy obu kanałów. Zapis dziennika jest nieusuwalny z poziomu Koordynatora i Wykonawcy.

| Rejestrowany element | Zakres zapisu |
|---|---|
| Zlecenie | Identyfikator zlecenia, treść zlecenia, konto i urządzenie źródłowe, chwila wydania |
| Zadanie | Identyfikator zadania, przynależność do zlecenia, Wykonawca, stan, chwila przydzielenia i zakończenia |
| Komunikat sterujący | Nadawca i odbiorca (Koordynator, Wykonawca), rodzaj komunikatu, chwila wymiany |
| Decyzja o ponowieniu | Zadanie objęte decyzją, wynik kontroli jakości, numer ponowienia, chwila decyzji |
| Decyzja Użytkownika | Przedmiot zatwierdzenia lub odmowy, punkt zatwierdzenia (rozdz. 13.4), chwila decyzji |
| Sterowanie przebiegiem | Wstrzymanie, wznowienie, przerwanie i korekta zlecenia wraz z ich źródłem |
| Przekroczenie granicy autonomii | Zadanie, granica, sposób rozstrzygnięcia |

Wpisy dziennika przechowywane są z odwołaniem do konta i urządzenia, bez powielania materiału uwierzytelniającego, który pozostaje w magazynie danych dostępowych (rozdz. 9).

### 13.6. Izolacja danych między zleceniami

| Zasada | Treść |
|---|---|
| Zasięg danych zadania | Zadanie ma dostęp wyłącznie do danych zlecenia, którego jest dekompozycją; identyfikator zlecenia jest częścią każdego komunikatu obu kanałów |
| Rozdział kontekstów | Kontekst pracy jednego zlecenia nie jest przenoszony do innego zlecenia ani między Wykonawcami realizującymi zadania różnych zleceń |
| Wynik pośredni | Wyniki pośrednie zadania pozostają w zasięgu zlecenia i są usuwane wraz z jego zamknięciem, o ile nie zostały zatwierdzone jako wynik trwały (rozdz. 13.4) |
| Zlecenia równoległe | Zlecenia realizowane równolegle prowadzone są w rozdzielnych zasięgach; zasięg procesu zawęża dodatkowo izolacja techniczna „konto i token per sesja" (rozdz. 8.4) |
| Dziennik audytu | Wpisy dziennika są przypisane do zlecenia; wgląd w pełny dziennik przysługuje Użytkownikowi zgodnie z rolą (rozdz. 14) |

### 13.7. Ustawienia autoryzacji kanałów komunikacji operacyjnej

| Ustawienie | Opis | Warstwa | Wartość domyślna |
|---|---|---|---|
| Granice autonomii Wykonawcy | Zakres działań realizowanych samodzielnie oraz działań wymagających zatwierdzenia (rozdz. 13.3) | globalna, przesłanialna na poziomie środowiska i roli | działania nieodwracalne wymagają zatwierdzenia |
| Punkty zatwierdzenia | Zestaw punktów obowiązkowego zatwierdzenia przez Użytkownika (rozdz. 13.4) | globalna, przesłanialna na poziomie środowiska | pełny zestaw punktów rozdziału 13.4 |
| Liczba ponowień zadania | Liczba ponowień podejmowanych przez Koordynatora przed skierowaniem sprawy do Użytkownika (rozdz. 13.3) | globalna, przesłanialna na poziomie roli | wartość ustalona w konfiguracji pętli wykonawczej |
| Zasięg dziennika audytu | Zakres rejestrowanych elementów przebiegu (rozdz. 13.5) | globalna | pełny zakres tabeli rozdziału 13.5 |

Każde ustawienie opatrzone jest objaśnieniem kontekstowym `[?]` opisującym jego działanie i wpływ na aplikację (Architektura, rozdz. 13; Model konfiguracji, rozdz. 3.3). Wyłączenie punktu zatwierdzenia nie znosi rejestracji zdarzenia w dzienniku audytu.

---

## 14. Warstwy widoczności jako mechanizm kontroli dostępu do funkcji

Interfejs Danaco Console ujawnia funkcje stopniowo: funkcja niepotrzebna do realizacji aktualnego zadania nie jest widoczna. Każdy element interfejsu należy do jednej z czterech warstw widoczności. Warstwy pełnią, obok uwierzytelniania i autoryzacji kanałów operacyjnych, funkcję kontroli dostępu do funkcji platformy — zakres funkcji ujawnianych użytkownikowi wynika z jego roli.

### 14.1. Warstwy widoczności a uprawnienia roli

| Warstwa | Zawartość | Sposób dostępu | Zależność od roli |
|---|---|---|---|
| 1 — zawsze widoczna | Chat Window, aktywne okno wiodące, kontekst pracy, podstawowa nawigacja, wskaźniki stanu wykonania | Widoczna bez interakcji | Dostępna każdej roli |
| 2 — widoczna na żądanie | Wybór modelu, wykonawcy, środowiska i trybu pracy, poziom wysiłku, parametry przepływu pracy | Ikona, przycisk, przełącznik, znacznik kontekstowy | Zakres pozycji wynika z konfiguracji roli |
| 3 — rozwinięcia kontekstowe | Zestawy akcji, ustawienia szybkie, warianty operacji | Menu kebab (⋮), menu hamburger (☰), menu kontekstowe, panel popover, lista rozwijana | Zakres pozycji wynika z konfiguracji roli |
| 4 — funkcje eksperckie | Najbardziej zaawansowane operacje, tryby administracyjne, narzędzia diagnostyczne niskiego poziomu, zmiana uprawnień i konfiguracji warstw | Polecenie języka naturalnego w Chat Window, skrót klawiszowy, wyszukiwarka funkcji palety poleceń `Ctrl/Cmd + K`, tryb administracyjny, konfiguracja roli | Dostępne wyłącznie w trybie administracyjnym albo przy konfiguracji roli obejmującej warstwę 4; użytkownik podstawowy nie widzi tych elementów |

Warstwa 4 stanowi granicę uprawnień: operacje administracyjne, diagnostyka niskiego poziomu, zmiana granic autonomii Wykonawcy (rozdz. 13.7), zarządzanie urządzeniami (rozdz. 12) i zmiana konfiguracji warstw są osiągalne wyłącznie w trybie administracyjnym albo przy roli, której konfiguracja obejmuje warstwę 4.

### 14.2. Tryb administracyjny

Drogę wejścia w tryb administracyjny, makietę znacznika stanu i diagram przejść stanu opisuje `interfejs-uzytkownika/konfiguracja.md` (rozdz. 9.6) — miejsce rozstrzygające. Niniejsza sekcja podaje model uprawnień trybu.

| Element modelu uprawnień | Treść |
|---|---|
| Warunek widoczności pozycji wejścia | Rola obejmująca warstwę 4 |
| Warunek wejścia | Potwierdzenie uprawnień — uwierzytelnienie właściciela konta (rozdz. 5) |
| Zakres | Funkcje warstwy 4: operacje administracyjne, narzędzia diagnostyczne, konfiguracja warstw widoczności według roli, ustawienia autoryzacji kanałów operacyjnych (rozdz. 13.7) |
| Stan spoczynku interfejsu | Funkcje trybu administracyjnego pozostają niewidoczne do chwili wejścia w tryb |
| Wygaśnięcie | Tryb wygasa po okresie bezczynności właściwym poświadczeniu urządzenia (rozdz. 12) oraz przy rozłączeniu urządzenia |
| Rejestracja | Wejście w tryb administracyjny, wyjście z niego, jego wygaśnięcie oraz każda zmiana konfiguracji uprawnień odnotowywane są w dzienniku audytu (rozdz. 13.5) |

### 14.3. Konfiguracja warstw według roli użytkownika

Konfiguracja warstw widoczności według roli jest elementem modelu uprawnień platformy: rola rozstrzyga, które warstwy i które pozycje warstw są ujawniane w interfejsie.

| Ustawienie | Opis | Warstwa | Wartość domyślna |
|---|---|---|---|
| Warstwy widoczności roli | Przypisanie warstw 1–4 do roli użytkownika oraz zakres pozycji warstw 2–3 (rozdz. 14.1) | globalna, przesłanialna na poziomie środowiska | rola właściciela obejmuje warstwy 1–4; rola podstawowa obejmuje warstwy 1–3 |
| Tryb administracyjny | Dostępność trybu administracyjnego dla roli (rozdz. 14.2) | globalna | dostępny roli obejmującej warstwę 4 |

Szablon konfiguracji ustawienia:

```
Ustawienie: „Warstwy widoczności roli"
  Zakres:            Aplikacja (rozdz. 12)
  Warstwa:           globalna, przesłanialna na poziomie środowiska
  Zawartość:         przypisanie warstw 1–4 do roli oraz zakres pozycji warstw 2–3
  Wartość domyślna:  rola właściciela — warstwy 1–4
                     rola podstawowa  — warstwy 1–3
  Skutek:            funkcje warstwy nieprzypisanej roli nie są prezentowane
                     w interfejsie ani zwracane przez wyszukiwarkę funkcji
  Objaśnienie [?]:   opis działania ustawienia i jego wpływu na aplikację
```

### 14.4. Wyszukiwarka funkcji a uprawnienia roli

| Zasada | Treść |
|---|---|
| Zakres wyników | Wyszukiwarka funkcji zwraca wyłącznie funkcje dostępne roli użytkownika |
| Nieujawnianie funkcji niedostępnych | Funkcja poza zakresem roli nie występuje w wynikach, w podpowiedziach ani w komunikatach wyszukiwarki — jej istnienie nie jest ujawniane |
| Polecenie języka naturalnego | Polecenie wskazujące funkcję poza zakresem roli nie uruchamia jej; Chat Window przedstawia informację o braku uprawnienia, bez opisu funkcji |
| Skrót klawiszowy | Skrót przypisany funkcji poza zakresem roli pozostaje bez działania |
| Zasada jednego kliknięcia | Funkcja dostępna roli osiągalna jest jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego, niezależnie od warstwy, w której się znajduje |

---

## 15. Zgodność z zasadami nadrzędnymi platformy

Uwierzytelnianie opisane w niniejszym dokumencie podlega czterem zasadom nadrzędnym platformy ustanowionym w rozdziale 14 Koncepcji platformy.

| Zasada nadrzędna | Realizacja w uwierzytelnianiu |
|---|---|
| Pełna kompozycyjność | Uwierzytelnianie jest mechanizmem wspólnym, niezależnym od środowiska i modułu — obowiązuje jednolicie przed wejściem do dowolnego z czterech środowisk platformy. Oba kanały komunikacji operacyjnej podlegają jednolitym regułom autoryzacji w każdym środowisku i module (rozdz. 2, rozdz. 13). |
| Pełna konfigurowalność (zasada centralna) | Zestaw aktywnych metod dodatkowych (rozdz. 6.4), wykaz i zarządzanie urządzeniami (rozdz. 3, rozdz. 12), granice autonomii Wykonawcy i punkty zatwierdzenia (rozdz. 13.7) oraz konfiguracja warstw widoczności według roli (rozdz. 14.3) podlegają jawnej konfiguracji, bez twardych blokad. |
| Jawność i konfigurowalność zależności | Relacja między tokenem urządzenia a zakresem izolacji technicznej „konto i token per sesja" jest jawna i osobno konfigurowalna z okna konfiguracji punktów izolacji, nie ukryta w mechanizmie uwierzytelniania (rozdz. 8.4). |
| Rozszerzenie orkiestracji (cel) | Zakres izolacji technicznej „konto i token per sesja", przypisywalny do poziomu zasięgu „rola" (Koncepcja platformy, rozdz. 6.5), nadaje odrębny token poszczególnym rolom środowiska MultitaskingAI; autoryzacja kanału Koordynator ↔ Wykonawca (rozdz. 13) obejmuje przydzielanie zadań i nadzór nad ich realizacją. |

---

## 16. Słowniczek pojęć

| Pojęcie | Znaczenie |
|---|---|
| **Uwierzytelnianie** | Jedyny mechanizm kontroli dostępu do platformy Danaco Console, obejmujący rejestrację, logowanie, metody dodatkowe i odzyskiwanie konta (rozdz. 2, Architektura rozdz. 15). |
| **Konto** | Jedyna encja właściciela platformy (Użytkownik), tworzona przy rejestracji i wykorzystywana przy każdym logowaniu; nie istnieje więcej niż jedno konto (rozdz. 3, rozdz. 10.1). |
| **Urządzenie** | Encja reprezentująca pojedyncze urządzenie Użytkownika połączone z platformą, posiadająca własny token dostępu (rozdz. 3, rozdz. 10.2). |
| **Token dostępu** | Poświadczenie wydawane urządzeniu po uwierzytelnieniu, wykorzystywane do nawiązania i utrzymania połączenia WebSocket; przechowywane poza bazą danych (rozdz. 8, rozdz. 9). |
| **Adres e-mail uwierzytelniający** | Adres poczty elektronicznej podany przy rejestracji, pełniący trzy funkcje: weryfikację konta przy rejestracji, metodę logowania oraz jedyną drogę odzyskania konta (rozdz. 4.3, rozdz. 6.3, rozdz. 7). |
| **Metoda dodatkowa uwierzytelniania** | Jedna z trzech metod konfigurowalnych z okna konfiguracji obok logowania loginem i hasłem: PIN, Windows Hello, adres e-mail uwierzytelniający jako samodzielna metoda logowania (rozdz. 6). |
| **Chat Window** | Główne okno komunikacji w kanale Użytkownik ↔ Wykonawca, umieszczone w lewej kolumnie obszaru roboczego; miejsce wydawania poleceń, zatwierdzania działań i przerywania przebiegu (rozdz. 13.2). |
| **Execution Loop Window** | Okno pętli wykonawczej w kanale Koordynator ↔ Wykonawca, otwierane jako kolumna sąsiadująca z Chat Window; prezentuje zlecenie, zadania, komunikaty sterujące i decyzje o ponowieniu (rozdz. 13.3). |
| **Użytkownik** | Osoba uwierzytelniona wobec konta właściciela, zlecająca i zatwierdzająca działania platformy (rozdz. 13.1). |
| **Koordynator** | Komponent orkiestrujący platformy, dekomponujący zlecenie na zadania, przydzielający je Wykonawcom i nadzorujący ich realizację (rozdz. 13.1). |
| **Wykonawca** | AI, agent lub system wykonawczy realizujący przydzielone zadania w granicach autonomii ustalonych przez Użytkownika (rozdz. 13.1, rozdz. 13.3). |
| **Granice autonomii Wykonawcy** | Zakres działań realizowanych przez Wykonawcę samodzielnie oraz działań wymagających zatwierdzenia przez Użytkownika (rozdz. 13.3, rozdz. 13.7). |
| **Punkt zatwierdzenia** | Moment przebiegu, w którym realizacja zadania wymaga jawnej decyzji Użytkownika wykonanej w Chat Window (rozdz. 13.4). |
| **Dziennik audytu** | Rejestr zleceń, zadań, komunikatów sterujących, decyzji o ponowieniu i decyzji Użytkownika, prowadzony dla obu kanałów komunikacji operacyjnej (rozdz. 13.5). |
| **Warstwa widoczności** | Jedna z czterech warstw ujawniania funkcji interfejsu; przypisanie warstw do roli użytkownika jest elementem modelu uprawnień (rozdz. 14.1, rozdz. 14.3). |
| **Tryb administracyjny** | Stan interfejsu udostępniający funkcje warstwy 4, dostępny roli obejmującej tę warstwę; wejście następuje z palety poleceń `Ctrl/Cmd + K` po potwierdzeniu uprawnień i jest odnotowywane w dzienniku audytu (rozdz. 14.2; `interfejs-uzytkownika/konfiguracja.md`, rozdz. 9.6). |
| **Magazyn danych dostępowych** | Wydzielony magazyn po stronie serwera, przechowujący poza bazą danych SQLite materiał bezpośrednio dowodzący tożsamości: skrót hasła, materiał metod PIN i Windows Hello oraz tokeny dostępu; baza danych przechowuje wyłącznie odwołania (rozdz. 9). |
| **Odzyskiwanie konta** | Procedura ustawienia nowego hasła przy wykorzystaniu adresu e-mail uwierzytelniającego, dostępna w sytuacji utraty dostępu do dotychczasowego materiału uwierzytelniającego (rozdz. 7). |

---

## Załącznik A. Schemat przepływu uwierzytelniania

Poniższy schemat przedstawia w formie tekstowej pełny przepływ uwierzytelniania: rejestrację, logowanie, odzyskiwanie konta oraz wejście do sesji pracy obsługującej oba kanały komunikacji operacyjnej.

```
Uruchomienie platformy na urządzeniu
    │
    ▼
Okno rejestracji i logowania
    │
    ├── Konto nie istnieje
    │       │
    │       ▼
    │   REJESTRACJA (rozdz. 4)
    │   login · adres e-mail uwierzytelniający · hasło
    │       │
    │       ▼
    │   Weryfikacja adresu e-mail (rozdz. 4.3)
    │       │
    │       ▼
    │   Konto właściciela utworzone (rozdz. 10.1)
    │
    └── Konto istnieje
            │
            ▼
        LOGOWANIE (rozdz. 5)
        login + hasło  albo  metoda dodatkowa (rozdz. 6):
        PIN · Windows Hello · e-mail uwierzytelniający
            │
 ┌──────────┴──────────┐
 ▼                     ▼
Weryfikacja pozytywna   Utrata dostępu
 │                     │
 ▼                     ▼
Token dostępu          ODZYSKIWANIE KONTA (rozdz. 7)
wydany urządzeniu      adres e-mail uwierzytelniający
(rozdz. 8.1)                │
 │                     ▼
 │               Ustawienie nowego hasła
 │                     │
 │                     ▼
 │               Unieważnienie tokenów
 │               wydanych wcześniej
 │                     │
 └──────────┬──────────┘
            ▼
 Połączenie WebSocket (Kontrakty komunikacji, rozdz. 2)
 krok 2 „Uwierzytelnienie" — token przedstawiony
            │
            ▼
 krok 3 „Powiązanie" — połączenie związane z sesją
            │
            ▼
 SESJA PRACY
   Chat Window            — kanał Użytkownik ↔ Wykonawca      (rozdz. 13.2)
   Execution Loop Window  — kanał Koordynator ↔ Wykonawca     (rozdz. 13.3)
   zakres funkcji ujawniony zgodnie z warstwami roli          (rozdz. 14)
```

---

## Załącznik B. Macierz metod uwierzytelniania

| Metoda | Charakter | Gdzie konfigurowana | Aktywna domyślnie | Przechowywanie materiału (rozdz. 9) | Pełni funkcję odzyskiwania konta |
|---|---|---|---|---|---|
| Hasło | Bazowa, obecna zawsze | — (metoda podstawowa, ustalana przy rejestracji) | Tak | Skrót hasła poza bazą danych | Nie |
| Adres e-mail uwierzytelniający | Weryfikacyjna, dodatkowa metoda logowania, odzyskiwanie konta | Okno konfiguracji (rozdz. 6.4, rozdz. 12) | Tak | Adres w encji Konto (rozdz. 10.1); nie jest materiałem tajnym | Tak — jedyna metoda pełniąca tę funkcję |
| PIN | Dodatkowa | Okno konfiguracji (rozdz. 6.4, rozdz. 12) | Nie | Materiał metody PIN poza bazą danych | Nie |
| Windows Hello | Dodatkowa, wiązana z urządzeniem | Okno konfiguracji (rozdz. 6.4, rozdz. 12) | Nie | Materiał metody Windows Hello poza bazą danych | Nie |

---

## Załącznik C. Scenariusze użycia

Scenariusze przedstawiono jako ponumerowane mini-przepływy ilustrujące działanie mechanizmów opisanych w rozdziałach 4–14.

### C.1. Podłączenie kolejnego urządzenia

1. Użytkownik ma założone konto na komputerze stacjonarnym.
2. Instaluje pakiet klienta na telefonie (Architektura, rozdz. 2).
3. Przy pierwszym uruchomieniu na telefonie wyświetla się okno logowania — nie rejestracji, ponieważ konto już istnieje (rozdz. 3) — z prośbą o login i hasło albo o aktywną metodę dodatkową.
4. Po pozytywnym uwierzytelnieniu telefon otrzymuje własny token dostępu (rozdz. 8.1).
5. Telefon zostaje dopisany do wykazu urządzeń połączonych z kontem (rozdz. 12), widocznego i zarządzalnego z okna konfiguracji na dowolnym z powiązanych urządzeń.
6. Drogą równoważną wobec logowania loginem i hasłem jest parowanie telefonu z maszyny już uwierzytelnionej — pełny protokół opisuje `funkcje-globalne/mobile.md`, rozdz. 9.

### C.2. Utrata dostępu i odzyskanie konta

1. Użytkownik zapomina hasło i nie ma dostępu do żadnego urządzenia z aktywnym tokenem.
2. W oknie logowania wybiera funkcję odzyskiwania konta i podaje adres e-mail uwierzytelniający (rozdz. 7).
3. Po potwierdzeniu tożsamości drogą mailową ustawia nowe hasło.
4. Platforma unieważnia tokeny wydane przed odzyskaniem, wymagając ponownego logowania na wszystkich dotąd powiązanych urządzeniach.
5. Zdarzenie `device.changed` (rozdz. 11.2) informuje o tym każde urządzenie pozostające podłączone w chwili unieważnienia.

### C.3. Zlecenie wymagające zatwierdzenia działania nieodwracalnego

1. Użytkownik wydaje zlecenie w Chat Window; polecenie zostaje powiązane z kontem, urządzeniem i identyfikatorem zlecenia (rozdz. 13.2).
2. Koordynator dekomponuje zlecenie na zadania i przedstawia dekompozycję do przyjęcia; Użytkownik zatwierdza uruchomienie pętli wykonawczej (rozdz. 13.4).
3. Execution Loop Window prezentuje kolejkę zadań, stan każdego z nich i wymianę komunikatów sterujących (rozdz. 13.3).
4. Zadanie osiąga punkt zapisu zmian trwałych — działanie nieodwracalne — i zatrzymuje się w stanie oczekiwania na decyzję.
5. Użytkownik zatwierdza działanie w Chat Window; zatwierdzenie, wraz z jego przedmiotem i chwilą, zostaje odnotowane w dzienniku audytu (rozdz. 13.5).
6. Po pozytywnej kontroli jakości Koordynator przedstawia zlecenie do zamknięcia, a Użytkownik je zatwierdza.

### C.4. Negatywna kontrola jakości i decyzja o ponowieniu

1. Kontrola jakości zadania kończy się wynikiem negatywnym; wynik zostaje przedstawiony w Execution Loop Window.
2. Koordynator podejmuje decyzję o ponowieniu zadania w granicach liczby ponowień wynikającej z ustawienia konfiguracyjnego (rozdz. 13.7).
3. Każde ponowienie zostaje odnotowane w dzienniku audytu wraz z numerem ponowienia i wynikiem kontroli jakości (rozdz. 13.5).
4. Po wyczerpaniu liczby ponowień zadanie zatrzymuje się, a sprawa zostaje skierowana do Użytkownika przez Chat Window (rozdz. 13.4).

### C.5. Wywołanie funkcji warstwy 4 i zakres wyszukiwarki funkcji

1. Użytkownik o roli podstawowej otwiera wyszukiwarkę funkcji i wpisuje nazwę narzędzia diagnostycznego niskiego poziomu.
2. Wyszukiwarka nie zwraca żadnego wyniku ani podpowiedzi — funkcja poza zakresem roli nie jest ujawniana (rozdz. 14.4).
3. Użytkownik o roli obejmującej warstwę 4 wywołuje pozycję „Tryb administracyjny” z palety poleceń `Ctrl/Cmd + K`, potwierdza uprawnienia i wchodzi w tryb administracyjny (rozdz. 14.2).
4. Wejście w tryb administracyjny zostaje odnotowane w dzienniku audytu (rozdz. 13.5).
5. Po zakończeniu czynności funkcje warstwy 4 przestają być widoczne, a interfejs wraca do stanu spoczynku obejmującego warstwę 1 i zwinięte wyzwalacze warstw 2–3.

---

## Załącznik D. Szablony konfiguracji tożsamości i dostępu

Poniższe szablony porządkują pola konfigurowalnych struktur opisanych w dokumencie. Są to szablony redakcyjne wspierające implementację; wszystkie ustawienia podlegają zasadzie braku twardych blokad oraz „brak ustawienia = wartość domyślna" (Koncepcja platformy, rozdz. 6, rozdz. 14).

### D.1. Ustawienie „Metody uwierzytelniania"

| Pole | Opis | Wartości / przykład |
|---|---|---|
| Zakres | Zakres okna konfiguracji | Aplikacja (rozdz. 12) |
| Warstwa | Warstwa konfiguracji | globalna |
| Metody włączalne | Zestaw metod do włączenia lub wyłączenia | hasło (stała) \| e-mail uwierzytelniający \| PIN \| Windows Hello |
| Wartość domyślna | Zestaw aktywny po rejestracji | hasło + e-mail uwierzytelniający |
| Zasada | Zachowanie przy braku metody dodatkowej | brak twardych blokad — logowanie metodą bazową i domyślnie aktywną |
| Objaśnienie `[?]` | Objaśnienie kontekstowe ustawienia | opis działania i wpływu na aplikację (rozdz. 12) |

### D.2. Ustawienie „Urządzenia połączone"

| Pole | Opis | Wartości / przykład |
|---|---|---|
| Zakres | Zakres okna konfiguracji | Aplikacja (rozdz. 12) |
| Warstwa | Warstwa konfiguracji | globalna |
| Zawartość | Wykaz encji Urządzenie powiązanych z Kontem | pozycje: identyfikator · nazwa · data ostatniego połączenia (rozdz. 10.2) |
| Operacje | Operacje dostępne na wykazie | `device.list` \| `device.revoke` (rozdz. 11.1) |
| Wartość domyślna | Ograniczenie liczby urządzeń | brak ograniczeń |
| Objaśnienie `[?]` | Objaśnienie kontekstowe ustawienia | opis działania i wpływu na aplikację (rozdz. 12) |

### D.3. Encja Konto

| Pole | Opis | Przechowywanie |
|---|---|---|
| login | Nazwa właściciela konta | Baza SQLite |
| adres e-mail uwierzytelniający | Adres wielofunkcyjny (rozdz. 4.3, 6.3, 7) | Baza SQLite |
| odwołanie do skrótu hasła | Wskazanie wpisu magazynu danych dostępowych | Baza SQLite (odwołanie); skrót poza bazą |
| aktywne metody uwierzytelniania | Stan ustawienia „Metody uwierzytelniania" | Baza SQLite (bez materiału) |
| data utworzenia | Data utworzenia konta | Baza SQLite |

### D.4. Encja Urządzenie

| Pole | Opis | Przechowywanie |
|---|---|---|
| identyfikator urządzenia | Jednoznaczny identyfikator | Baza SQLite |
| nazwa | Nazwa urządzenia | Baza SQLite |
| odwołanie do tokenu dostępu | Wskazanie wpisu magazynu danych dostępowych | Baza SQLite (odwołanie); token poza bazą |
| data ostatniego połączenia | Data ostatniego połączenia | Baza SQLite |

### D.5. Wpis magazynu danych dostępowych

| Pole | Opis | Wartości / przykład |
|---|---|---|
| Odwołanie | Identyfikator wpisu przechowywany w bazie SQLite | odwołanie z encji Konto (11.1) lub Urządzenie (11.2) |
| Rodzaj | Rodzaj przechowywanego materiału | skrót hasła \| materiał PIN \| materiał Windows Hello \| token dostępu urządzenia \| dane dostępowe kanału modelu |
| Materiał | Wartość przechowywana poza bazą danych | nie replikowana do bazy SQLite |
| Dostęp | Droga dostępu do materiału | wyłącznie przez kontrakt dostępu (Architektura, rozdz. 1.3, zasada 4) |

---

*Koniec dokumentu. Danaco Console — Bezpieczeństwo i uwierzytelnianie, wersja 2.0.*

---
*Danaco Console — AI Workspace OS · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — [LICENSE](LICENSE). Kontakt: support@danaco-group.pl*
