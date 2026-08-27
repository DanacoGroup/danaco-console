# Danaco Console — Elementy okien przepływu głównego

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
| **Odbiorcy** | Designer (co, gdzie, w jakiej formie, do czego) · Deweloper (co zbudować) |
| **Autor koncepcji źródłowej** | Operator (Dariusz Naharnowicz, Danaco Group) |
| **Opracowanie** | Danaco Console — agent projektowy inwentarza elementów |
| **Źródła** | architektura/koncepcja-platformy.md (autorytatywna) · interfejs-uzytkownika/strona-glowna-i-nawigacja.md · specyfikacje/specyfikacja-okien-operacyjnych.md · interfejs-uzytkownika/system-wizualny.md · architektura/model-konfiguracji.md · architektura/architektura.md · architektura/bezpieczenstwo-i-uwierzytelnianie.md · `srodowiska/multitaskingai.md` · architektura/izolacja-i-zaleznosci.md · tokens.css · components.css |

Dokument wyszczególnia **wszystkie elementy** znajdujące się w każdym oknie przepływu głównego platformy Danaco Console, które **nie ma własnego agenta projektowego**: okno startowe/logowania, strona główna (trzy strefy) oraz powłoki czterech środowisk — TalkIn, WorkSpace, CodeStudio, MultitaskingAI (pasek górny, karty sesji, boczna nawigacja, obszar roboczy) wraz z obydwoma kanałami komunikacji operacyjnej — Chat Window (Użytkownik ↔ Wykonawca) i Execution Loop Window (Koordynator ↔ Wykonawca). Wszystkie okna rozmieszczone są w układzie pionowym, w kolumnach sąsiadujących poziomo. Dla każdego elementu podano: co to jest, do czego służy, formę i wagę wizualną, stany, zachowanie po interakcji, miejsce występowania, warstwę widoczności (1–4, rozdz. 7) oraz sposób wywołania — tak aby żaden element nie został zamieniony w projekcie na niewłaściwą formę (duży panel tam, gdzie miała być mała ikonka, i odwrotnie).

**Poza zakresem tego dokumentu** (mają własnego agenta projektowego albo są opisane w dokumentach modułów — rozdział 0.2): okno konfiguracji (pełny zakres trzynastu obszarów ustawień oraz okno punktów izolacji), moduł Agents wraz z jego sześcioma oknami operacyjnymi, piętnaście modułów platformy wraz z ich oknami operacyjnymi (Studio Editor, Code Editor, Workflow Builder itd.), wnętrze czterech okien roboczych ról środowiska MultitaskingAI, pełne widoki funkcji globalnych Mobile i Always On Display.

---

## Spis treści

- [0. Zakres dokumentu i mapa okien](#0-zakres-dokumentu-i-mapa-okien)
- [1. Okno startowe / logowania](#1-okno-startowe--logowania)
- [2. Strona główna — Centrum dowodzenia](#2-strona-główna--centrum-dowodzenia)
- [3. Powłoka środowisk modułowych — TalkIn / WorkSpace / CodeStudio](#3-powłoka-środowisk-modułowych--talkin--workspace--codestudio)
- [4. Powłoka środowiska MultitaskingAI](#4-powłoka-środowiska-multitaskingai)
- [5. Tabela zbiorcza dokumentu](#5-tabela-zbiorcza-dokumentu)
- [6. Stany wspólne i słowniczek](#6-stany-wspólne-i-słowniczek)
- [7. Warstwy widoczności elementów interfejsu](#7-warstwy-widoczności-elementów-interfejsu)
- [Załącznik A. Katalog ikon wykorzystanych w oknach tego dokumentu](#załącznik-a-katalog-ikon-wykorzystanych-w-oknach-tego-dokumentu)
- [Załącznik B. Mapa źródeł](#załącznik-b-mapa-źródeł)

---

## 0. Zakres dokumentu i mapa okien

### 0.1. Cztery okna/powłoki w zakresie

| Okno / powłoka | Zawartość dokumentowana tutaj | Rozdział |
|---|---|---|
| Okno startowe / logowania | Rejestracja, logowanie, metody dodatkowe, odzyskiwanie konta, faza budowy | 1 |
| Strona główna — Centrum dowodzenia | Strefa 1 (karty środowisk), Strefa 2 (kafle komponentów własnych), Strefa 3 (listwa ustawień) | 2 |
| Powłoka środowisk modułowych | Pasek górny, karty sesji, boczna nawigacja modułów, obszar roboczy (jako kontener), Chat Window, Execution Loop Window, panele pomocnicze — wspólne dla TalkIn, WorkSpace, CodeStudio | 3 |
| Powłoka środowiska MultitaskingAI | Pasek górny, karty sesji, panel orkiestracji (odpowiednik bocznej nawigacji), obszar roboczy (jako kontener), oba okna komunikacji | 4 |

### 0.2. Poza zakresem — gdzie szukać

| Element poza zakresem | Dlaczego poza zakresem | Gdzie opisany |
|---|---|---|
| Okno konfiguracji (13 zakresów ustawień + okno punktów izolacji) | Ma własnego agenta projektowego | Dokument okna Konfiguracji/Ustawień |
| Moduł Agents — 6 okien operacyjnych (Agent Builder, Model Configuration, Skills Manager, Connectors Manager, Permissions Center) | Ma własnego agenta projektowego | Dokument agenta modułu Agents |
| Piętnaście modułów platformy — ich okna operacyjne (Studio Editor, Code Editor, Workflow Builder, Research Workspace, Library Explorer, Design Board, Voice Console, Terminal Tabs, Diagnostics Center, Product Builder itd.) | Okna operacyjne modułów — dokumentowane per moduł | Dokumenty modułów (`moduly/`) |
| Wnętrze czterech okien roboczych ról MultitaskingAI (Executor 1, Executor 2, Coordinator, Executor 3 / Validator) — pola, zakładki, konfiguracja roli | Okno operacyjne środowiska — dokumentowane odrębnie | Dokument środowiska MultitaskingAI (`srodowiska/`) |
| Pełne widoki funkcji globalnych Mobile i Always On Display | Funkcje globalne o własnym zakresie, poza mechanizmem przełączania środowisk/modułów | `funkcje-globalne/mobile.md` |

Niniejszy dokument opisuje wyłącznie **punkt wejścia** do powyższych elementów (ikonę, kafel, pozycję listwy) — nie ich wnętrze.

**Schemat 1 — Miejsce okien tego dokumentu w hierarchii platformy**

```
OKNO STARTOWE / LOGOWANIA                                    ◄── rozdział 1 (ten dokument)
        │  uwierzytelnienie urządzenia — wymóg konfigurowalny przez Operatora (Pomiń, gdy wymóg wyłączony)
        ▼
STRONA GŁÓWNA — Centrum dowodzenia                            ◄── rozdział 2 (ten dokument)
├─ Strefa 1 · karty środowisk ──────► wybór środowiska
├─ Strefa 2 · kafle komponentów własnych ──► okno konfiguracji komponentu    [poza zakresem]
└─ Strefa 3 · listwa ustawień ──────► okno konfiguracji / Mobile / AOD      [poza zakresem]
        │
        ▼
POWŁOKA ŚRODOWISKA  (TalkIn · WorkSpace · CodeStudio · MultitaskingAI)      ◄── rozdziały 3–4 (ten dokument)
├─ pasek górny
├─ karty sesji
├─ boczna nawigacja modułów  albo  panel orkiestracji (MultitaskingAI)
└─ obszar roboczy  (kontener)
        │  wybór modułu / sekcji
        ▼
OKNA OPERACYJNE MODUŁU  (Chat Window ── ten dokument, rozdz. 3.5.1/4.5;         [wnętrze modułu:
         Execution Loop Window ── rozdz. 3.5.2;                                poza zakresem]
         panele pomocnicze ── rozdz. 3.5.3;
         pozostałe okna modułu ── poza zakresem)
```

---

## 1. Okno startowe / logowania

### 1.1. Miejsce w przepływie i model jednego konta

Danaco Console działa w modelu jednego właściciela (Operatora) i wielu urządzeń: jedno konto, dowolna liczba urządzeń połączonych (komputery, telefony, tablety), każde z własnym tokenem dostępu (Bezpieczeństwo i uwierzytelnianie, rozdz. 3). Z tego modelu wynikają dwa odrębne stany okna startowego, nigdy współwystępujące na tym samym urządzeniu w tej samej chwili:

**Schemat 2 — Rejestracja i logowanie wobec modelu „jedno konto, wiele urządzeń"**

```
                    ┌─────────────────────────────────┐
                    │   KONTO — jedyne (Operator)      │
                    └────────────────┬────────────────┘
                                     │
        ┌────────────────┬───────────┴───────────┬────────────────┐
        ▼                ▼                        ▼                ▼
 Urządzenie 1      Urządzenie 2             Urządzenie 3      Urządzenie …
 (pierwsze          (kolejne, np.            (kolejne)         (bez limitu)
  uruchomienie      telefon)
  w historii
  platformy)
        │                │                        │                │
        ▼                ▼                        ▼                ▼
   REJESTRACJA      LOGOWANIE                LOGOWANIE         LOGOWANIE
   (konto jeszcze    (konto już istnieje — Załącznik C.3 dokumentu Bezpieczeństwo)
   nie istnieje)
```

W praktyce niemal każde uruchomienie okna startowego pokazuje **stan logowania** — stan rejestracji występuje wyłącznie przy zupełnie pierwszym uruchomieniu platformy w jej historii, zanim istnieje jakiekolwiek konto (Bezpieczeństwo i uwierzytelnianie, rozdz. 4.1).

Dodatkowo istnieje wymiar niezależny od powyższego: **faza wdrożenia**. W fazie budowy (prace wykonawcze, przed migracją na serwer docelowy) okno ma pełny docelowy wygląd, a wymóg uwierzytelnienia jest domyślnie wyłączony — dostępny jest przycisk **Pomiń**. W fazie produkcyjnej (po migracji rdzenia na serwer i otwarciu portu) wymóg logowania jest domyślnie włączony, ale — zgodnie z zasadą nadrzędną „Zero blokerów" — pozostaje **konfigurowalny przez Operatora w oknie konfiguracji** (Bezpieczeństwo i uwierzytelnianie, rozdz. 8): to Operator decyduje, czy uwierzytelnienie jest wymagane, także w fazie produkcyjnej, i to jego ustawienie — nie sama faza wdrożenia — rozstrzyga o obecności przycisku **Pomiń**. Uwierzytelnianie pozostaje jedynym sankcjonowanym mechanizmem kontroli dostępu i jest zawsze dostępne jako granica bezpieczeństwa — zmienne jest wyłącznie jego egzekwowanie. Kod obsługi jest wspólny dla obu faz — różni je wyłącznie ustawienie wymogu logowania, jawnie widoczne i zmienialne przez Operatora w oknie konfiguracji.

### 1.2. Diagram pełnego przepływu

**Diagram 1 — Przepływ uwierzytelniania: rejestracja, logowanie, odzyskiwanie, faza budowy**

```
Uruchomienie platformy na urządzeniu bez ważnego tokenu
        │
        ▼
┌────────────────────────────────────────────────────────────────────────────┐
│                         OKNO STARTOWE / LOGOWANIA                          │
│                                                                            │
│  WYMÓG LOGOWANIA: WYŁĄCZONY          │  WYMÓG LOGOWANIA: WŁĄCZONY          │
│  (ustawienie Operatora, okno         │  (ustawienie Operatora, okno        │
│  konfiguracji — domyślnie w fazie    │  konfiguracji — domyślnie w fazie   │
│  budowy, zmienialne zawsze)          │  produkcyjnej, zmienialne zawsze)   │
│      │                               │      │                              │
│      ├─[Pomiń]→ wejście bez          │      ├── konto nie istnieje         │
│      │          uwierzytelnienia     │      │        │                     │
│      │                               │      │        ▼                     │
│      └─ formularz dostępny,          │      │   REJESTRACJA (1.4‑A)        │
│         nie blokuje dalszej pracy    │      │   login · e-mail · hasło     │
│                                      │      │        │                     │
│                                      │      │        ▼                     │
│                                      │      │   weryfikacja adresu e-mail  │
│                                      │      │        │                     │
│                                      │      │        ▼                     │
│                                      │      │   KONTO UTWORZONE            │
│                                      │      │                              │
│                                      │      └── konto istnieje             │
│                                      │               │                     │
│                                      │               ▼                     │
│                                      │         LOGOWANIE (1.4‑B)           │
│                                      │   login+hasło albo metoda dodatkowa │
│                                      │   (PIN · Windows Hello · e-mail)    │
│                                      │          │            │             │
│                                      │   pozytywne      utrata dostępu     │
│                                      │          │            │             │
│                                      │          ▼            ▼             │
│                                      │   token wydany   ODZYSKIWANIE       │
│                                      │   urządzeniu     KONTA (1.4‑D)      │
│                                      │                  e-mail uwierz.     │
│                                      │                  → nowe hasło       │
│                                      │                  → unieważnienie    │
│                                      │                    tokenów innych   │
│                                      │                    urządzeń         │
└──────────────────────────────────────┴─────────────────────────────────────┘
        │
        ▼
Połączenie WebSocket — krok „Powiązanie” z sesją → STRONA GŁÓWNA (rozdz. 2)
```

### 1.3. Elementy wspólne obu stanów (rejestracja i logowanie)

| Element | Co to jest | Do czego służy | Forma i waga |
|---|---|---|---|
| Rama okna | Kontener formularza, wyśrodkowany w oknie klienckim | Prezentacja formularza właściwego bieżącemu stanowi | Duży panel/karta (`.dn-modal`‑podobny, wyśrodkowany), tło powierzchni, na jasnym motywie kontrastujący z tłem aplikacji |
| Godło marki | `logo-danaco.svg` — tarcza granatowa z monogramem w złocie | Identyfikacja platformy przed uwierzytelnieniem | Ikona pełnokolorowa, nad formularzem, rozmiar średni-duży |
| Nazwa platformy | „Danaco Console” | Identyfikacja produktu | Krój nagłówkowy (Cormorant Garamond), pod godłem |

### 1.4. Makiety tekstowe

**Makieta 1 — Rejestracja (pierwsze uruchomienie platformy w jej historii)**

```
┌──────────────────────────────────────────────────────────┐
│                     ◈  Danaco Console                       │
│                                                            │
│              Utworzenie konta właściciela                 │
│                                                            │
│   Login                                                    │
│   [ ______________________________________________ ]      │
│                                                            │
│   Adres e-mail uwierzytelniający                           │
│   [ ______________________________________________ ]      │
│      pełni też funkcję logowania i odzyskiwania konta      │
│                                                            │
│   Hasło                                                     │
│   [ ______________________________________________ ]  👁   │
│                                                            │
│                [        Utwórz konto        ]              │
│                                                            │
│   ──────────────────────────────────────────────────       │
│   Masz już konto?  Zaloguj się                              │
└──────────────────────────────────────────────────────────┘

Po wysłaniu formularza (stan pośredni):
┌──────────────────────────────────────────────────────────┐
│   Konto utworzone — stan: NIEPOTWIERDZONE                  │
│   Wysłaliśmy wiadomość na podany adres e-mail.              │
│   Potwierdź adres, aby dokończyć rejestrację.               │
│                                                            │
│                [   Wyślij ponownie   ]                     │
└──────────────────────────────────────────────────────────┘
```

**Makieta 2 — Logowanie (stan typowy: konto już istnieje)**

```
┌──────────────────────────────────────────────────────────┐
│                     ◈  Danaco Console                        │
│                                                            │
│                        Logowanie                           │
│                                                            │
│   ( Login i hasło )( PIN )( Windows Hello )( E-mail )       │  ← przełącznik metody
│     aktywna zawsze   *      *                 aktywna       │    (tylko metody włączone
│                     wyłączona domyślnie                     │     w oknie konfiguracji)
│                                                            │
│   Login                                                     │
│   [ ______________________________________________ ]       │
│                                                            │
│   Hasło                                                      │
│   [ ______________________________________________ ]  👁    │
│                                                            │
│                [          Zaloguj          ]                │
│                                                            │
│   ──────────────────────────────────────────────────        │
│   Nie pamiętasz hasła?  Odzyskaj dostęp                      │
└──────────────────────────────────────────────────────────┘

Stan błędu (dane niepoprawne):
   Login
   [ ______________________________________________ ]
   Hasło
   [ ______________________________________________ ]  👁
   ⚠ Nieprawidłowy login lub hasło. Spróbuj ponownie.
                [          Zaloguj          ]
```

**Makieta 3 — Logowanie z wyłączonym wymogiem logowania (przycisk Pomiń obecny — domyślnie w fazie budowy, konfigurowalne przez Operatora w każdej fazie)**

```
┌──────────────────────────────────────────────────────────┐
│                     ◈  Danaco Console                        │
│                                                            │
│                        Logowanie                            │
│                                                            │
│   Login                                                      │
│   [ ______________________________________________ ]        │
│   Hasło                                                       │
│   [ ______________________________________________ ]  👁     │
│                                                            │
│                [          Zaloguj          ]                 │
│                                                            │
│   ──────────────────────────────────────────────────         │
│                    [        Pomiń        ]                    │
│         wejście do platformy bez uwierzytelnienia             │
│         (widoczny, gdy Operator wyłączy wymóg logowania —      │
│          domyślnie wyłączony w fazie budowy; Operator może      │
│          wyłączyć wymóg logowania także w fazie produkcyjnej —  │
│          rozdz. 8.1 dokumentu Bezpieczeństwo i uwierzytelnianie)│
└──────────────────────────────────────────────────────────┘
```

**Makieta 4 — Odzyskiwanie konta (dwa kroki)**

```
Krok 1 — żądanie odzyskania:
┌──────────────────────────────────────────────────────────┐
│              Odzyskiwanie konta                            │
│                                                            │
│   Adres e-mail uwierzytelniający                            │
│   [ ______________________________________________ ]        │
│                                                            │
│                [   Wyślij link odzyskiwania   ]              │
│   ──────────────────────────────────────────────────        │
│   Powrót do logowania                                        │
└──────────────────────────────────────────────────────────┘

Krok 2 — po potwierdzeniu tożsamości drogą mailową:
┌──────────────────────────────────────────────────────────┐
│              Ustawienie nowego hasła                        │
│                                                            │
│   Nowe hasło                                                 │
│   [ ______________________________________________ ]  👁     │
│                                                            │
│                [     Ustaw nowe hasło     ]                  │
│                                                            │
│   ℹ Po zmianie hasła wszystkie pozostałe urządzenia          │
│     zostaną wylogowane i będą wymagały ponownego              │
│     zalogowania (rozdz. 7 dokumentu Bezpieczeństwo).           │
└──────────────────────────────────────────────────────────┘
```

### 1.5. Tabela elementów — pełna

**Tabela 1 — Elementy okna startowego / logowania**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Pole „Login” | Pole tekstowe jednoliniowe | Wprowadzenie nazwy właściciela konta | `.dn-input` w `.dn-pole`, etykieta nad polem | Puste (placeholder) · wypełnione · fokus (obrys złoty) · błąd (`.dn-pole--blad`, obrys czerwony + komunikat pod polem) | Wpisywana wartość waliduje się na bieżąco (format); Enter przenosi fokus do pola Hasło | Rejestracja, Logowanie | 1 | Widoczne bez interakcji |
| Pole „Adres e-mail uwierzytelniający” | Pole tekstowe jednoliniowe | Weryfikacja konta, alternatywna metoda logowania, jedyna droga odzyskania konta (rozdz. 4.2, 6.3, 7 Bezpieczeństwo) | `.dn-input` w `.dn-pole` | Puste · wypełnione · fokus · błąd (format nieprawidłowy) | Format e-mail walidowany na bieżąco | Rejestracja, Odzyskiwanie konta (krok 1) | 1 | Widoczne bez interakcji |
| Pole „Hasło” | Pole hasła z przełącznikiem widoczności | Wprowadzenie/ustawienie hasła dostępowego | `.dn-input` typ hasło + ikona `oko` (pokaż/ukryj) | Puste · wypełnione · ukryte (domyślnie) · odkryte · fokus · błąd | Klik ikony `oko` przełącza maskowanie tekstu; hasło nigdy nie trafia do bazy danych w postaci jawnej (przechowywany wyłącznie skrót, poza bazą — rozdz. 10 Bezpieczeństwo) | Rejestracja, Logowanie (metoda bazowa), Odzyskiwanie konta (krok 2, jako „Nowe hasło”) | 1 | Widoczne bez interakcji |
| Przycisk „Utwórz konto” | Przycisk podstawowy formularza rejestracji | Wysłanie danych rejestracyjnych do serwera | `.dn-btn--zloty` (jedyny CTA widoku, pełna szerokość formularza) | Domyślny · hover · aktywny (wciśnięcie) · fokus (pierścień złoty) · ładowanie (`.dn-spinner` w miejscu etykiety podczas weryfikacji) | Przycisk zawsze klikalny; gdy wymagane pola nie są wypełnione, klik pokazuje ostrzeżenie przy brakującym polu zamiast wysyłki; w przeciwnym razie wysyła `auth.register` — sukces → stan „konto niepotwierdzone”; błąd serwera → komunikat pod właściwym polem | Rejestracja | 1 | Widoczny bez interakcji |
| Przycisk „Zaloguj” | Przycisk podstawowy formularza logowania | Wysłanie danych logowania do serwera | `.dn-btn--zloty`, pełna szerokość formularza | Domyślny · hover · aktywny · fokus · ładowanie | Przycisk zawsze klikalny; gdy pola są puste, klik pokazuje ostrzeżenie zamiast wysyłki; w przeciwnym razie wysyła `auth.login` — sukces → token wydany, przejście do strony głównej; błąd → komunikat błędu nad przyciskiem | Logowanie | 1 | Widoczny bez interakcji |
| Link „Masz już konto? Zaloguj się” | Odnośnik tekstowy | Przełączenie widoku formularza z rejestracji na logowanie w obrębie tego samego okna | Tekst-link, krój bazowy, kolor `--dn-akcent-txt` | Domyślny · hover (podkreślenie) · fokus | Klik podmienia zawartość okna na formularz logowania (bez przeładowania okna aplikacji) | Rejestracja | 1 | Widoczny bez interakcji |
| Link „Nie pamiętasz hasła? Odzyskaj dostęp” | Odnośnik tekstowy | Wejście w procedurę odzyskiwania konta | Tekst-link, krój bazowy | Domyślny · hover · fokus | Klik podmienia zawartość okna na formularz odzyskiwania (krok 1) | Logowanie | 1 | Widoczny bez interakcji |
| Link „Powrót do logowania” | Odnośnik tekstowy | Wyjście z procedury odzyskiwania bez jej dokończenia | Tekst-link, drugorzędny | Domyślny · hover · fokus | Klik wraca do formularza logowania | Odzyskiwanie konta | 1 | Widoczny bez interakcji |
| Przełącznik metody logowania | Grupa segmentowa (przełącznik gęsty) | Wybór aktywnej metody uwierzytelnienia spośród metod włączonych w oknie konfiguracji (rozdz. 6.4 Bezpieczeństwo) | `.dn-zakladki--pigulki` — segmenty: Login i hasło (zawsze) · PIN · Windows Hello · E-mail uwierzytelniający (tylko włączone) | Segment aktywny (tło `--dn-akcent-tlo`) · pozostałe segmenty nieaktywne · segment niedostępny (metoda wyłączona w konfiguracji — nie jest w ogóle renderowany, zgodnie z zasadą braku twardych blokad: brak metody = mniej segmentów, nie zablokowany segment) | Klik segmentu podmienia pola formularza pod przełącznikiem na właściwe wybranej metodzie | Logowanie (widoczny tylko, gdy co najmniej jedna metoda dodatkowa jest włączona) | 2 | Ujawniany, gdy w oknie konfiguracji włączona jest co najmniej jedna metoda dodatkowa; wybór segmentu |
| Przycisk „Utwórz konto” / komunikat „Wyślij ponownie” (weryfikacja e-mail) | Przycisk drugorzędny | Ponowne wysłanie wiadomości weryfikacyjnej, gdy pierwsza nie dotarła | `.dn-btn--zarys` | Domyślny · hover · aktywny · odliczanie informacyjne po wysłaniu (etykieta pokazuje czas do zalecanej ponownej próby — nie blokuje kliknięcia) | Przycisk pozostaje klikalny przez cały czas; klik zawsze ponawia wysyłkę wiadomości weryfikacyjnej — odliczanie ma charakter wyłącznie informacyjny, nie wstrzymuje akcji | Rejestracja (stan „konto niepotwierdzone”) | 2 | Ujawniany zdarzeniem — przejściem formularza w stan „konto niepotwierdzone” |
| Komunikat stanu „Konto niepotwierdzone” | Blok informacyjny | Poinformowanie, że rejestracja wymaga potwierdzenia adresu e-mail przed pierwszym logowaniem | `.dn-alert--info` | Widoczny do chwili potwierdzenia adresu | Zniknięcie po potwierdzeniu (odczytanym przez serwer) — okno przechodzi do stanu „konto utworzone” → wydanie tokenu | Rejestracja (stan pośredni) | 2 | Ujawniany zdarzeniem utworzenia konta |
| Przycisk „Pomiń” | Przycisk drugorzędny/duch | Wejście do platformy bez wypełniania formularza | `.dn-btn--duch` lub `.dn-btn--zarys`, pod przyciskiem głównym, oddzielony separatorem | Domyślny · hover · aktywny · fokus | Klik wchodzi bezpośrednio do strony głównej bez wydania tokenu uwierzytelniania (krok „Uwierzytelnienie” cyklu połączenia jest pomijany — rozdz. 8.1 Bezpieczeństwo) | Widoczny, gdy Operator wyłączy wymóg logowania w oknie konfiguracji — domyślnie wyłączony w fazie budowy; ta sama opcja dostępna Operatorowi również w fazie produkcyjnej (rozdz. 8 i 8.1 Bezpieczeństwo) | 2 | Ujawniany ustawieniem „wymóg logowania” w oknie konfiguracji |
| Pole „Nowe hasło” | Pole hasła | Ustawienie hasła zastępującego dotychczasowe po odzyskaniu konta | `.dn-input` typ hasło + ikona `oko` | Puste · wypełnione · fokus · błąd | Zapis zastępuje dotychczasowy skrót hasła (poza bazą danych); unieważnia tokeny pozostałych urządzeń | Odzyskiwanie konta (krok 2) | 1 | Widoczne bez interakcji |
| Przycisk „Wyślij link odzyskiwania” | Przycisk podstawowy | Zainicjowanie procedury odzyskiwania | `.dn-btn--zloty` | Domyślny · hover · aktywny · ładowanie | Przycisk zawsze klikalny; gdy pole e-mail jest puste, klik pokazuje ostrzeżenie zamiast wysyłki; w przeciwnym razie wysyła `auth.recover` — serwer wysyła wiadomość z drogą potwierdzenia tożsamości | Odzyskiwanie konta (krok 1) | 1 | Widoczny bez interakcji |
| Przycisk „Ustaw nowe hasło” | Przycisk podstawowy | Zatwierdzenie nowego hasła | `.dn-btn--zloty` | Domyślny · hover · aktywny · ładowanie | Przycisk zawsze klikalny; gdy pole nowego hasła jest puste lub nie spełnia reguł, klik pokazuje ostrzeżenie zamiast wysyłki; w przeciwnym razie wysyła `auth.reset` — sukces → powrót do formularza logowania z komunikatem potwierdzającym (`.dn-toast--sukces`) | Odzyskiwanie konta (krok 2) | 1 | Widoczny bez interakcji |
| Komunikat błędu logowania | Blok/tekst inline nad przyciskiem „Zaloguj” | Poinformowanie o nieprawidłowych danych | `.dn-alert--blad` lub tekst `--dn-blad-fg` z ikoną `blad` | Widoczny wyłącznie po nieudanej próbie | Znika przy kolejnej próbie wysłania formularza | Logowanie | 2 | Ujawniany zdarzeniem nieudanej próby logowania |
| Komunikat potwierdzający zmianę hasła | Toast | Potwierdzenie skutecznego odzyskania konta | `.dn-toast--sukces`, kreska lewa zielona | Pojawia się na 1 zdarzenie, znika automatycznie | Informuje dodatkowo o wylogowaniu pozostałych urządzeń | Odzyskiwanie konta (po kroku 2) | 2 | Ujawniany zdarzeniem zmiany hasła; znika samoczynnie |

### 1.6. Metody dodatkowe uwierzytelniania — tabela odniesienia

**Tabela 2 — Metody logowania dostępne w przełączniku (Bezpieczeństwo i uwierzytelnianie, rozdz. 6, Załącznik B)**

| Metoda | Aktywna domyślnie | Pola w formularzu | Pełni funkcję odzyskiwania |
|---|---|---|---|
| Login i hasło | Tak (bazowa, obecna zawsze) | Login · Hasło | Nie |
| Adres e-mail uwierzytelniający | Tak | Adres e-mail (jako identyfikator logowania) | Tak — jedyna metoda pełniąca tę funkcję |
| PIN | Nie — wymaga świadomego włączenia w oknie konfiguracji | Krótki kod numeryczny | Nie |
| Windows Hello | Nie — wymaga świadomego włączenia w oknie konfiguracji | Mechanizm systemowy (bez pola tekstowego) | Nie |

Włączanie i wyłączanie metod dodatkowych odbywa się w oknie konfiguracji (ustawienie „Metody uwierzytelniania”, zakres Aplikacja) — poza zakresem tego dokumentu. Zgodnie z zasadą braku twardych blokad, wyłączenie metody dodatkowej nie ogranicza dostępu: logowanie pozostaje możliwe metodą bazową (login i hasło) oraz metodą domyślnie aktywną (e-mail uwierzytelniający).

---

## 2. Strona główna — Centrum dowodzenia

### 2.1. Trzy ścieżki działania

Strona główna nie otwiera żadnej przestrzeni roboczej wprost — pełni funkcję przedpokoju przed strefą roboczą, z którego prowadzą trzy niezależne ścieżki (Strona główna i nawigacja, rozdz. 2).

**Schemat 3 — Trzy ścieżki z centrum dowodzenia**

```
Użytkownik po uwierzytelnieniu
    │
    ▼
STRONA GŁÓWNA — Centrum dowodzenia
    │
    ├─▶ Strefa 1 · wybór środowiska ────▶ przestrzeń robocza środowiska (rozdz. 3–4 tego dokumentu)
    ├─▶ Strefa 2 · komponent własny ────▶ okno konfiguracji komponentu       [poza zakresem]
    └─▶ Strefa 3 · ustawienia ──────────▶ okno konfiguracji / funkcja globalna [poza zakresem]
```

### 2.2. Układ całościowy

**Makieta 5 — Strona główna, układ całościowy**

```
╔═══════════════════════════════════════════════════════════════════════╗
║  ◈ Danaco Console                                                        ║
╠═══════════════════════════════════════════════════════════════════════╣
║                                                                         ║
║  STREFA 1 · KARTY ŚRODOWISK                    (waga główna — środek   ║
║                                                  ciężkości strony)      ║
║  ┌──────────┐   ┌───────────┐   ┌────────────┐   ┌────────────────┐   ║
║  │ ◈ godło  │   │ ◈ godło   │   │ ◈ godło    │   │ ◈ godło        │   ║
║  │ TalkIn   │   │ WorkSpace │   │ CodeStudio │   │ MultitaskingAI │   ║
║  │ opis     │   │ opis      │   │ opis       │   │ opis           │   ║
║  └──────────┘   └───────────┘   └────────────┘   └────────────────┘   ║
║                                                                         ║
║  STREFA 2 · KAFLE KOMPONENTÓW WŁASNYCH              (waga pośrednia)   ║
║  ┌─────────────┐  ┌────────┐  ┌───────────┐  ┌───────────┐            ║
║  │ ▧ Automations│  │▧ Agents│  │▧ Workspace│  │▧ Assistant│            ║
║  │ „Zbuduj      │  │„Skonfi-│  │„Załóż     │  │„Ustaw     │            ║
║  │  automatykę” │  │guruj   │  │ projekt”  │  │profil     │            ║
║  │              │  │agenta” │  │           │  │asystenta” │            ║
║  └─────────────┘  └────────┘  └───────────┘  └───────────┘            ║
║                                                                         ║
║  STREFA 3 · LISTWA USTAWIEŃ              (waga najniższa, stale        ║
║                                            dostępna)                    ║
║  [ ⚙ Okno konfiguracji ]   [ 📱 Mobile ]   [ 🔔 Always On Display ]     ║
║                                                                         ║
╚═══════════════════════════════════════════════════════════════════════╝
```

Trzy strefy różnią się formą prezentacji celowo — hierarchia skali (karty → kafle → listwa) komunikuje wagę, zanim użytkownik przeczyta którykolwiek napis (Koncepcja platformy, rozdz. 7.4).

### 2.3. Strefa 1 — karty środowisk

**Makieta 6 — Pojedyncza karta środowiska, trzy stany**

```
Stan spoczynkowy                    Stan hover / aktywny
┌──────────────────────┐            ┌══════════════════════┐
│                       │            ║ ▐▐▐ wstęga złota 3px ║
│       ◈  godło        │            ║        ◈  godło       ║
│                       │            ║                       ║
│      T a l k I n      │            ║      T a l k I n      ║
│                       │            ║                       ║
│  Wiedza, komunikacja  │            ║  Wiedza, komunikacja  ║
│  i praca z treścią    │            ║  i praca z treścią    ║
│                       │            ║  cień złoty, uniesienie║
└──────────────────────┘            ╚══════════════════════╝
 tło neutralne, bez złota            akcent złoty (wstęga górna
                                      3px + cień `--dn-cien-zloto`)
                                      uniesienie translateY(-2px)
```

**Tabela 3 — Elementy pojedynczej karty środowiska**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Karta środowiska (całość) | Duży kafel wejścia | Wejście do przestrzeni roboczej jednego z czterech środowisk | `.dn-karta--interaktywna --akcent`, jedna z czterech w siatce równej wagi — waga główna strony | Spoczynek (bez złota) · hover (wstęga górna 3px w `--dn-akcent`, cień `--dn-cien-zloto`, uniesienie −2px) · aktywny/fokus (pierścień złoty 2px, klawiatura) | Klik/Enter otwiera przestrzeń roboczą wybranego środowiska (krok 1 nawigacji, rozdz. 4.1 Strona główna) | Strefa 1 strony głównej — 4 wystąpienia | 1 | Widoczna bez interakcji |
| Godło środowiska | Symbol graficzny właściwy środowisku | Szybka identyfikacja wizualna środowiska | Ikona/maska `.dn-ikona--lg`, nad tytułem | Dziedziczy kolor z kontekstu karty (`currentColor`) | Brak własnej interakcji — element wewnątrz karty | W każdej z 4 kart | 1 | Widoczna bez interakcji |
| Tytuł karty | Nazwa środowiska | Jednoznaczna identyfikacja tekstowa | Krój nagłówkowy (Cormorant Garamond), rozmiar `2xl`/`3xl` | Stały | Brak własnej interakcji | W każdej z 4 kart | 1 | Widoczna bez interakcji |
| Opis trybu pracy | Jednozdaniowy opis | Doprecyzowanie, czym zajmuje się środowisko, przed wejściem | Krój bazowy, rozmiar `sm`/`base`, kolor drugorzędny | Stały | Brak własnej interakcji | W każdej z 4 kart | 1 | Widoczna bez interakcji |

**Tabela 4 — Treść czterech kart środowisk (kolejność źródłowa)**

| Kolejność | Środowisko | Opis trybu pracy na karcie | Prowadzi do |
|---|---|---|---|
| 1 | TalkIn | Wiedza, komunikacja i praca z treścią | Powłoka TalkIn (rozdz. 3) |
| 2 | WorkSpace | Produktywność, organizacja i realizacja projektów | Powłoka WorkSpace (rozdz. 3) |
| 3 | CodeStudio | Programowanie | Powłoka CodeStudio (rozdz. 3) |
| 4 | MultitaskingAI | Orkiestracja autonomicznej pracy ciągłej | Powłoka MultitaskingAI (rozdz. 4) |

### 2.4. Strefa 2 — kafle komponentów własnych

**Makieta 7 — Pojedynczy kafel komponentu własnego**

```
┌─────────────────────┐
│      ▧  ikona        │   ← skala mniejsza niż godło środowiska
│                      │
│    Automations       │   ← krój bazowy, semibold (NIE nagłówkowy)
│                      │
│  „Zbuduj automatykę” │   ← etykieta zorientowana na tworzenie
└─────────────────────┘
   bez wstęgi złotej i bez cienia złotego na spoczynku —
   wizualnie lżejszy niż karta środowiska (strefa 1)
```

**Tabela 5 — Elementy pojedynczego kafla komponentu własnego**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Kafel komponentu (całość) | Kafel średniej wagi | Wejście do okna konfiguracji jednego z czterech rodzajów komponentu własnego | `.dn-karta` (bez wariantu `--akcent`), mniejszy niż karta strefy 1, 4 w siatce — waga pośrednia | Spoczynek · hover (tło, bez wstęgi złotej) · aktywny/fokus (pierścień złoty) | Klik/Enter otwiera modal (`.dn-modal`) albo okno konfiguracji komponentu własnego — **poza zakresem tego dokumentu** | Strefa 2 strony głównej — 4 wystąpienia | 1 | Widoczny bez interakcji |
| Ikona komponentu | Symbol właściwy rodzajowi komponentu | Identyfikacja wizualna rodzaju wytworu | `.dn-ikona`, skala mniejsza niż godło środowiska | Dziedziczy kolor z kontekstu | Brak własnej interakcji | W każdym z 4 kafli | 1 | Widoczny bez interakcji |
| Nazwa komponentu | Nazwa modułu źródłowego komponentu | Identyfikacja tekstowa | Krój bazowy, waga `semibold`, NIE nagłówkowy — odróżnienie od strefy 1 | Stały | Brak własnej interakcji | W każdym z 4 kafli | 1 | Widoczny bez interakcji |
| Etykieta działania | Fraza zorientowana na tworzenie | Komunikacja, że kliknięcie tworzy nowy wytwór, a nie otwiera przestrzeń roboczą | Tekst drugorzędny pod nazwą | Stały | Brak własnej interakcji | W każdym z 4 kafli | 1 | Widoczny bez interakcji |

**Tabela 6 — Treść czterech kafli komponentów własnych**

| Kolejność | Komponent własny | Etykieta działania | Powstający wytwór | Okno docelowe |
|---|---|---|---|---|
| 1 | Automations | „Zbuduj automatykę” | Automatyka | Okno konfiguracji Automations — poza zakresem |
| 2 | Agents | „Skonfiguruj agenta” | Agent | Agent Builder — poza zakresem (własny agent) |
| 3 | Workspace | „Załóż projekt” | Projekt | Okno konfiguracji modułu Workspace — poza zakresem |
| 4 | Assistant | „Ustaw profil asystenta” | Profil asystenta | Okno konfiguracji Assistant — poza zakresem |

### 2.5. Strefa 3 — listwa ustawień

**Makieta 8 — Listwa ustawień**

```
┌──────────────────────────────────────────────────────────────────┐
│   [ ⚙  Okno konfiguracji ]    [ 📱  Mobile ]    [ 🔔  Always On   │
│                                                     Display ]      │
└──────────────────────────────────────────────────────────────────┘
  pasek poziomy zwarty · tło `--dn-powierzchnia-2` (nie tło marki) ·
  bez kroju nagłówkowego · złoto wyłącznie w fokusie/najechaniu
  pojedynczego elementu · najniższa waga wizualna trzech stref
```

**Tabela 7 — Elementy listwy ustawień**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Pozycja „Okno konfiguracji” | Przycisk z ikoną i etykietą tekstową | Otwarcie pełnego okna konfiguracji platformy (13 zakresów ustawień) | `.dn-btn-ikona` + etykieta, wzorzec lekki `.dn-pasek` | Domyślny · hover · aktywny · fokus | Klik otwiera okno konfiguracji — **poza zakresem tego dokumentu** (własny agent) | Strefa 3 strony głównej | 1 | Widoczna bez interakcji |
| Pozycja „Mobile” | Przycisk z ikoną i etykietą tekstową | Aktywacja funkcji globalnej Mobile (dostęp mobilny, monitoring procesów) | `.dn-btn-ikona` + etykieta | Domyślny · hover · aktywny · fokus | Klik aktywuje funkcję globalną — nie tworzy nowej przestrzeni roboczej, nie przeładowuje niczego (Koncepcja platformy, rozdz. 5.3) | Strefa 3 strony głównej — także dostępna z paska górnego każdej powłoki środowiska (rozdz. 3.2, 4.3) | 1 | Widoczna bez interakcji |
| Pozycja „Always On Display” | Przycisk z ikoną i etykietą tekstową | Aktywacja globalnego agenta towarzyszącego (pływający avatar) | `.dn-btn-ikona` + etykieta, ikona `dzwonek`/awatar | Domyślny · hover · aktywny · fokus | Klik przywołuje pływający avatar — dostępny wszędzie, nie tworzy przestrzeni roboczej | Strefa 3 strony głównej — także dostępna z paska górnego każdej powłoki środowiska (rozdz. 3.2, 4.3) | 1 | Widoczna bez interakcji |

### 2.6. Przepływ nawigacji ze strony głównej

**Diagram 2 — Cztery kroki nawigacji: strona główna → środowisko → moduł → okno → karty**

```
┌─────────────────────────────────────────────────────────────┐
│ STRONA GŁÓWNA — Centrum dowodzenia                            │
└─────────────────────────────────────────────────────────────┘
        │  KROK 1 — klik karty środowiska (strefa 1, rozdz. 2.3)
        ▼
┌─────────────────────────────────────────────────────────────┐
│ POWŁOKA ŚRODOWISKA   TalkIn · WorkSpace · CodeStudio ·        │
│                       MultitaskingAI  (rozdz. 3–4)             │
└─────────────────────────────────────────────────────────────┘
        │  KROK 2 — boczna nawigacja modułów (rozdz. 3.4)
        │           albo panel orkiestracji (rozdz. 4.4, MultitaskingAI)
        ▼
┌─────────────────────────────────────────────────────────────┐
│ MODUŁ  (Studio, Research, Developer …) — poza zakresem         │
│   albo ROLA / SEKCJA panelu orkiestracji — poza zakresem       │
└─────────────────────────────────────────────────────────────┘
        │  KROK 3 — otwarcie okien modułu w obszarze roboczym (rozdz. 3.5, 4.5)
        ▼
┌─────────────────────────────────────────────────────────────┐
│ OBSZAR ROBOCZY = Chat Window (ten dokument) + okna modułu       │
│                   (poza zakresem — dokumenty modułów)          │
└─────────────────────────────────────────────────────────────┘
        │  KROK 4 — praca odbywa się w karcie sesji bieżącej (rozdz. 3.3, 4.3)
        ▼
  [ Karta sesji A · Studio ]  [ Karta sesji B · Research ]  [ + nowa karta ]
```

---

## 3. Powłoka środowisk modułowych — TalkIn / WorkSpace / CodeStudio

Trzy środowiska — TalkIn, WorkSpace, CodeStudio — dzielą dokładnie tę samą mechanikę powłoki; różni je wyłącznie **zestaw modułów widocznych w bocznej nawigacji** (rozdz. 3.4). Z tego powodu opisano je łącznie w jednym rozdziale, z tabelą różnic w miejscu, gdzie faktycznie występują.

### 3.1. Anatomia powłoki

**Schemat 4 — Kolumny powłoki środowiska modułowego (układ pionowy, podział lewa–prawa)**

```
┌───────────────────────────────────────────────────────────────────────┐
│ PASEK GÓRNY                                                    (3.2)   │
│ godło+nazwa platformy │ dom │ szukaj │ [Ubuntu][Fable 5][Ultra] │ ⋮    │
├───────────────────────────────────────────────────────────────────────┤
│ KARTY SESJI                                                    (3.3)   │
│ [● Studio ✕] [ Research ✕] [+]                                        │
├──────────┬─────────────┬───────────────────────────────┬──────────────┤
│ BOCZNA   │ CHAT WINDOW │ OBSZAR ROBOCZY MODUŁU  (3.5)  │ PANEL        │
│ NAWIGACJA│ Użytkownik ↔│                               │ POMOCNICZY   │
│ MODUŁÓW  │ Wykonawca   │  okna operacyjne modułu       │ (rozszerzenie│
│ (3.4)    │   (3.5.1)   │  wybranego w bocznej          │  boczne,     │
│          │             │  nawigacji — poza zakresem    │  3.5.3)      │
│ • Studio │ ─────────── │  tego dokumentu               │              │
│ ▸ Research│ EXECUTION  │                               │              │
│ • Library│ LOOP WINDOW │                               │              │
│ • …      │ Koordynator │                               │              │
│          │ ↔ Wykonawca │                               │              │
│          │   (3.5.2)   │                               │              │
└──────────┴─────────────┴───────────────────────────────┴──────────────┘

Kolumny sąsiadują poziomo i zajmują pełną wysokość obszaru roboczego.
Regulacji podlega wyłącznie szerokość kolumn. Chat Window zajmuje lewą
kolumnę stałą, Execution Loop Window otwiera się jako kolumna sąsiadująca,
panele pomocnicze — jako rozszerzenia boczne po prawej stronie obszaru
roboczego (Specyfikacja okien operacyjnych, rozdz. 1.2).
```

**Makieta 9 — Powłoka środowiska TalkIn, przykład wypełniony (stan spoczynku interfejsu)**

```
╔═══════════════════════════════════════════════════════════════════════╗
║ ◈ Danaco Console  🏠  [ 🔍 Szukaj… ]  [Ubuntu][Fable 5][Ultra]   ⋮ ☀  ║
╠═══════════════════════════════════════════════════════════════════════╣
║ [●  Studio          ✕] [   Research         ✕] [ + ]                   ║
╠══════════╦═════════════╦══════════════════════════════════╦═══════════╣
║ TalkIn   ║ CHAT WINDOW ║ Studio Editor · Tools · Diff/Grep ║ Panel     ║
║──────────║ Użytkownik ↔║                                  ║ pomocniczy║
║   Studio ║ Wykonawca   ║   [ okna operacyjne modułu       ║ (zwinięty)║
║ ▸ Workspace║           ║     Studio — poza zakresem       ║           ║
║   Browser║ historia    ║     tego dokumentu ]             ║           ║
║   Research║ rozmowy    ║                                  ║           ║
║   Library║             ║                                  ║           ║
║   Translate║ [pole      ║                                  ║           ║
║   Roundtable║ poleceń…] ║                                  ║           ║
║   Assistant║ 📎 Operacje▼ ➤║                               ║           ║
║   Agents ║─────────────║                                  ║           ║
║          ║ EXECUTION   ║                                  ║           ║
║          ║ LOOP WINDOW ║                                  ║           ║
║          ║ Koordynator ║                                  ║           ║
║          ║ ↔ Wykonawca ║                                  ║           ║
╚══════════╩═════════════╩══════════════════════════════════╩═══════════╝
```

### 3.2. Pasek górny — rozstrzygnięcie źródłowe

System wizualny ustala wprost: „Pasek nawigacji stanowi wzorzec dla górnego paska okna środowiska (Koncepcja platformy, rozdz. 8), wspólnego dla wszystkich czterech środowisk” (System wizualny, rozdz. 8.7) — komponent `.dn-pasek`, wysokość `--dn-topbar-h` (56 px), tło marki (granat na motywie jasnym / grafit na motywie ciemnym, `tokens.css`).

Poniższa tabela łączy trzy rodzaje elementów: **(a)** wprost nazwane w komponencie `.dn-pasek` (blok marki, pole wyszukiwania, przyciski ikonowe), **(b)** funkcje, których źródła wprost stwierdzają dostępność „z poziomu każdego środowiska i modułu” niezależnie od konkretnego miejsca w interfejsie (Mobile, Always On Display — Strona główna i nawigacja, rozdz. 3.4, 9.2) i **(c)** uproszczone menu kontekstowe okna operacyjnego, którego istnienie i zakres ustala Model konfiguracji (rozdz. 3.1, 4.3) i Architektura (rozdz. 13) bez wskazania piksela — najbardziej spójnym miejscem jest pasek górny, wspólny dla całej powłoki. Pasek górny prezentuje wyłącznie elementy warstwy 1 oraz zwinięte wyzwalacze warstw 2–3 (znaczniki kontekstowe, menu kebab, element zbiorczy `Operacje ▼`) — rozdz. 7.

**Tabela 8 — Elementy paska górnego**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Blok marki/logotypu | `.dn-pasek-marka` — godło + nazwa „Danaco Console” | Identyfikacja platformy; orientacja, że użytkownik jest w powłoce środowiska (nie na stronie głównej) | Blok tekstowo-graficzny w lewym narożniku paska, krój nagłówkowy, fragment nazwy w złocie, mały dopisek wersalikowy pod spodem | Stały (bez interakcji własnej) | Brak — identyfikacja, nie nawigacja (nawigacja do strony głównej odbywa się przez odrębną, jawną ikonę „dom” — zasada świadomego, widocznego przejścia, Strona główna i nawigacja, rozdz. 8.2) | Pasek górny wszystkich 4 środowisk | 1 | Widoczny bez interakcji |
| Ikona „Strona główna” (dom) | Przycisk ikonowy | Świadomy powrót do Centrum dowodzenia | `.dn-btn-ikona`, ikona `dom`, na ciemnym tle paska | Domyślny · hover (podświetlenie tła) · aktywny (wciśnięcie) · fokus (pierścień złoty) | Klik zamyka powłokę bieżącego środowiska i pokazuje stronę główną; karty sesji bieżącego środowiska pozostają aktywne w tle i odtwarzają pełny stan przy powrocie (Strona główna i nawigacja, rozdz. 8.2–8.3) | Pasek górny wszystkich 4 środowisk | 1 | Widoczna bez interakcji |
| Pole wyszukiwania | `.dn-pasek-szukaj` — pole tekstowe z ikoną `szukaj` | Wyszukiwanie w obrębie platformy | Pole elastyczne do 520 px, ikona lupy po lewej wewnątrz pola, tło przyciemnione względem paska | Puste (placeholder) · wypełnione · fokus (obrys złoty) · wynik/brak wyniku | Wpisanie frazy inicjuje wyszukiwanie; Enter/wybór wyniku nawiguje do trafienia | Pasek górny wszystkich 4 środowisk | 1 | Widoczne bez interakcji |
| Ikona „Szybka konfiguracja sesji” | Przycisk ikonowy otwierający uproszczone menu kontekstowe | Szybka zmiana ustawień **warstwy sesji** bieżącej karty bez opuszczania okna operacyjnego (Model konfiguracji, rozdz. 3.1, 4.3; Architektura, rozdz. 13) | `.dn-btn-ikona` (ikona `ustawienia`), otwiera rozwijany panel/menu przy pasku | Domyślny · hover · aktywny · otwarte (podświetlone tło ikony + widoczny panel) | Klik otwiera podzbiór ustawień właściwy wyłącznie warstwie sesji (np. kanał modelu, motyw, retencja historii tej karty) + odnośnik „Otwórz pełne okno konfiguracji” prowadzący do okna poza zakresem tego dokumentu | Pasek górny wszystkich 4 środowisk; dotyczy zawsze aktywnej (frontowej) karty sesji | 3 | Klik ikony rozwija menu kontekstowe; wybór pozycji zwija je |
| Ikona „Mobile” | Przycisk ikonowy | Wejście w globalny tryb dostępu mobilnego / podgląd procesów w tle (`funkcje-globalne/mobile.md`) | `.dn-btn-ikona`, ikona urządzenia mobilnego | Domyślny · hover · aktywny · fokus | Klik aktywuje funkcję globalną Mobile — nie tworzy nowej przestrzeni roboczej (Koncepcja platformy, rozdz. 5.3) | Pasek górny wszystkich 4 środowisk — jeden z kilku równoważnych punktów dostępu (Strona główna i nawigacja, rozdz. 9.2) | 2 | Pozycja menu kebab paska górnego (⋮) — jedno kliknięcie |
| Ikona „Always On Display” | Przycisk ikonowy | Przywołanie/schowanie pływającego avatara agenta towarzyszącego | `.dn-btn-ikona`, ikona `dzwonek`/awatar (Załącznik C systemu wizualnego) | Domyślny · hover · aktywny · fokus | Klik przywołuje pływający avatar Always On Display, obecny ponad bieżącą przestrzenią roboczą | Pasek górny wszystkich 4 środowisk — jeden z kilku równoważnych punktów dostępu | 2 | Pozycja menu kebab paska górnego (⋮) — jedno kliknięcie |
| Przełącznik motywu jasny/ciemny | Mała ikona dwustanowa | Zmiana motywu wizualnego całej aplikacji | `.dn-btn-ikona`, ikony `slonce`/`ksiezyc` | Jasny (słońce) · ciemny (księżyc) · hover · fokus | Klik przełącza atrybut `data-theme` na elemencie głównym dokumentu; przejście animowane (`prefers-reduced-motion` respektowane); zapisywane jako ustawienie warstwy sesji lub domyślnej (System wizualny, rozdz. 9) | Pasek górny wszystkich 4 środowisk | 2 | Ikona dwustanowa w pasku górnym albo pozycja menu kebab (⋮) |

### 3.3. Karty sesji

**Makieta 10 — Pasek kart sesji**

```
┌────────────────────────┬────────────────────────┬───────┐
│  ●  Studio          ✕ │     Research         ✕ │   +   │
└────────────────────────┴────────────────────────┴───────┘
   ▲ karta aktywna            ▲ karta w tle            ▲ nowa karta
   ●  wskaźnik pracy w tle    ✕  kontrolka zamknięcia
```

**Makieta 11 — Anatomia pojedynczej karty sesji**

```
┌──────────────────────────────┐
│  ●   Studio               ✕ │
│  │      │                  │ │
│  │      │                  └── kontrolka zamknięcia karty
│  │      └── tytuł = nazwa aktualnie otwartego modułu
│  │          (aktualizuje się przy przeładowaniu — rozdz. 3.5)
│  └── wskaźnik stanu = sygnalizacja aktywnej pracy w tle
│      (powiązany z funkcją Mobile)
└──────────────────────────────┘
 Klik w dowolnym miejscu karty (poza „✕”) = przełączenie na tę sesję
```

**Tabela 9 — Elementy karty sesji**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Tytuł karty | Nazwa modułu aktualnie otwartego w tej karcie | Identyfikacja zawartości karty w pasku kart | Tekst krótki, krój bazowy, wewnątrz karty | Aktualizuje się przy przeładowaniu modułu (rozdz. 3.5) | Brak interakcji własnej — element informacyjny | Każda karta sesji | 1 | Widoczny bez interakcji |
| Wskaźnik stanu | Mała kropka/plakietka stanu | Sygnalizacja aktywnej pracy w tle (proces automatyki, aktywna rola MultitaskingAI itp.) | `.dn-kropka` (kolor semantyczny wg statusu) | Brak aktywności (niewidoczny) · aktywny (kropka widoczna, kolor wg statusu) | Powiązany z funkcją Mobile — monitoring procesów; klik otwiera podgląd statusu procesu | Każda karta sesji z aktywnym procesem w tle | 1 | Widoczny bez interakcji, gdy proces w tle jest aktywny |
| Kontrolka zamknięcia („✕”) | Mała ikona | Zamknięcie karty sesji | `.dn-btn-ikona` mały, ikona `zamknij`, w prawym rogu karty | Domyślny · hover · aktywny · fokus | Klik zamyka widok karty w bieżącym oknie; los procesu sesji leżącego u podstaw jest funkcją modelu sesji i procesów (Architektura, rozdz. 5), poza zakresem tego dokumentu | Każda karta sesji | 2 | Ujawniana najechaniem na kartę; klik zamyka kartę |
| Kontrolka nowej karty („+”) | Przycisk ikonowy | Otwarcie nowej, pustej karty sesji w bieżącym środowisku | `.dn-btn-ikona`, ikona `plus`, na końcu paska kart | Domyślny · hover · aktywny · fokus | Klik tworzy nową kartę z odrębnym układem, historią i kontekstem (stan wyjściowy izolacji, konfigurowalny); boczna nawigacja / panel orkiestracji aktywna, obszar roboczy w stanie pustym (rozdz. 3.5) | Koniec paska kart sesji, wszystkie 4 środowiska | 1 | Widoczna bez interakcji |
| Obszar przełączania karty | Cała powierzchnia karty poza kontrolką zamknięcia | Uczynienie karty aktywną | Cała karta jako obszar klikalny | Aktywna (podświetlona, na pierwszym planie) · w tle (przygaszona) · hover | Klik przełącza główny obszar roboczy na układ, historię i kontekst zapisane w tej karcie; boczna nawigacja/panel orkiestracji podświetla właściwy moduł/sekcję | Każda karta sesji | 1 | Widoczny bez interakcji |

**Diagram 3 — Cykl życia karty sesji**

```
                    kontrolka „+”
                         │
                         ▼
              ┌─────────────────────┐
              │  NOWA KARTA          │  obszar roboczy: stan pusty
              │  (pusta, oczekuje    │  (rozdz. 3.5) — oczekiwanie
              │   na wybór modułu)   │  na wybór modułu w bocznej
              └──────────┬──────────┘  nawigacji
                         │  wybór modułu w bocznej nawigacji
                         ▼
              ┌─────────────────────┐   klik innej karty
              │  KARTA AKTYWNA        │◄──────────────────┐
              │  (na pierwszym       │                    │
              │   planie)            │──────────────────► │
              └──────────┬──────────┘   klik tej karty    │
                         │                                 │
                         │  klik innej karty                │
                         ▼                                 │
              ┌─────────────────────┐                      │
              │  KARTA W TLE          │──────────────────────┘
              │  (stan zachowany:     │
              │   układ, historia,    │
              │   kontekst)           │
              └──────────┬──────────┘
                         │  kontrolka „✕”
                         ▼
              ┌─────────────────────┐
              │  KARTA ZAMKNIĘTA      │  znika z paska; los procesu
              │  (w widoku)           │  sesji — poza zakresem
              └─────────────────────┘  (Architektura, rozdz. 5)
```

### 3.4. Boczna nawigacja modułów

**Makieta 12 — Boczna nawigacja, środowisko TalkIn (przykład źródłowy)**

```
  ┌──────────────┐
  │ TalkIn       │   ← nagłówek: nazwa środowiska
  ├──────────────┤
  │   Studio     │
  │   Workspace  │
  │   Browser    │
  │ ▸ Research   │  ← stan wybrany (podświetlenie, tło `--dn-akcent-tlo`)
  │   Library    │
  │   Translate  │
  │   Roundtable │
  │   Assistant  │
  │   Agents     │
  └──────────────┘
```

**Tabela 10 — Elementy bocznej nawigacji modułów**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Nagłówek nawigacji | Nazwa bieżącego środowiska | Orientacja — potwierdzenie, w którym środowisku znajduje się użytkownik | Krój nagłówkowy, u góry panelu | Stały | Brak interakcji | Boczna nawigacja każdego z 3 środowisk modułowych | 1 | Widoczny bez interakcji |
| Pozycja modułu | Wiersz listy pionowej | Wybór modułu do otwarcia w bieżącej karcie sesji | `.dn-karta--pozycja`, tekst i ikona modułu | Domyślny · hover · **wybrany** (podświetlenie, tło `--dn-akcent-tlo`, tekst `--dn-akcent-txt`) · fokus | Klik przeładowuje obszar roboczy bieżącej karty do zestawu okien nowo wybranego modułu (rozdz. 3.5); podświetlenie przenosi się na klikniętą pozycję | Boczna nawigacja — liczba pozycji zależna od środowiska (Tabela 11) | 1 | Widoczna bez interakcji |
| Panel nawigacji (całość) | Stały, pionowy pas z lewej strony obszaru roboczego | Prezentacja listy modułów dostępnych w danym środowisku | Pas pionowy o stałej szerokości, wspólny dla wszystkich kart sesji tego środowiska (nie jest odrębny per karta) | Widoczny niezależnie od tego, który moduł jest otwarty | Wybór pozycji zmienia stan wybrany zgodnie z modułem otwartym w aktywnej (frontowej) karcie sesji | Wszystkie 3 środowiska modułowe | 1 | Widoczny bez interakcji; na wąskich ekranach rozwijany menu hamburger (☰) — warstwa 3 |

**Tabela 11 — Macierz modułów widocznych w bocznej nawigacji per środowisko (kolejność źródłowa)**

| Środowisko | Liczba pozycji | Lista modułów w kolejności wyświetlania |
|---|---|---|
| TalkIn | 9 | Studio · Workspace · Browser · Research · Library · Translate · Roundtable · Assistant · Agents |
| WorkSpace | 9 | Studio · Workspace · Browser · Research · Library · Roundtable · Design · Apps · Agents |
| CodeStudio | 8 | Workspace · Roundtable · Design · Terminal · Developer · Diagnostics · Apps · Agents |

Moduł Automations nie występuje w żadnej z trzech list — jest konfigurowany wyłącznie ze strony głównej (strefa 2), a gotowe automatyki trafiają do sesji jako komponenty własne (Koncepcja platformy, rozdz. 10). Treść każdego modułu (jego okna operacyjne) jest opisana w dokumencie tego modułu — poza zakresem niniejszego dokumentu.

### 3.5. Obszar roboczy — kontener

Ten podrozdział opisuje obszar roboczy **jako kontener powłoki** — jego anatomię ogólną, stany i jedyny element w nim obecny niezależnie od modułu: Chat Window. Wnętrze poszczególnych okien modułu (Studio Editor, Code Editor itd.) jest przedmiotem dokumentów modułów.

**Schemat 5 — Anatomia okna operacyjnego (wzorzec wspólny każdego okna w obszarze roboczym)**

```
 ┌─ Nagłówek okna ─────────────────────────────────────────────
 │   nazwa okna · przynależność do bieżącego modułu
 ├─────────────────────────────────────────────────────────────
 │   OBSZAR ZAWARTOŚCI
 │     treść właściwa oknu (edytor, lista, monitor, podgląd …)
 │     aktualizacja na żywo kanałem WebSocket — dotyczy okien
 │     monitorujących procesy oraz okien strumienia odpowiedzi
 ├─────────────────────────────────────────────────────────────
 │   ELEMENTY KONFIGURACJI
 │     każde ustawienie opatrzone objaśnieniem kontekstowym [?]
 ├─────────────────────────────────────────────────────────────
 │   STAN PROCESU SESJI (po stronie serwera)
 │     trwały — rozłączenie klienta nie zamyka okna
 └─────────────────────────────────────────────────────────────
```

**Makieta 13 — Obszar roboczy, stan pusty (nowa karta, przed wyborem modułu)**

```
┌───────────────────────────────────────────────────────────┐
│                                                             │
│                                                             │
│                    (ikona pustego stanu)                    │
│                                                             │
│         Wybierz moduł z bocznej nawigacji, aby              │
│                  rozpocząć pracę                             │
│                                                             │
│                                                             │
└───────────────────────────────────────────────────────────┘
   `.dn-pusty-stan` — analogicznie do nowej karty przeglądarki
   bez załadowanej strony (Strona główna i nawigacja, rozdz. 7.2)
```

**Makieta 14 — Obszar roboczy, stan wypełniony (moduł wybrany)**

```
┌─────────────┬────────────────────────────────────┬──────────────────┐
│ CHAT WINDOW │  okna operacyjne modułu wybranego  │ Panel pomocniczy │
│ Użytkownik ↔│  w bocznej nawigacji — liczba      │ (rozszerzenie    │
│ Wykonawca   │  i układ zależne od modułu —       │  boczne,         │
│ (3.5.1)     │  poza zakresem tego dokumentu,     │  zwinięty        │
│             │  patrz dokument modułu             │  w spoczynku)    │
│ historia    │                                    │                  │
│ rozmowy     │                                    │                  │
│             │                                    │                  │
│ [ pole      │                                    │                  │
│  poleceń… ] │                                    │                  │
│ 📎 Operacje▼ ➤│                                  │                  │
│─────────────│                                    │                  │
│ EXECUTION   │                                    │                  │
│ LOOP WINDOW │                                    │                  │
│ Koordynator │                                    │                  │
│ ↔ Wykonawca │                                    │                  │
│ (3.5.2)     │                                    │                  │
└─────────────┴────────────────────────────────────┴──────────────────┘
   lewa kolumna stała, pełna wysokość obszaru roboczego (Chat Window
   i Execution Loop Window) · kolumna dominująca: obszar roboczy modułu ·
   kolumny boczne: panele pomocnicze otwierane jako rozszerzenia boczne
```

**Tabela 12 — Elementy kontenera obszaru roboczego**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Obszar roboczy (kontener) | Główna, największa powierzchnia powłoki | Pomieszczenie zestawu okien operacyjnych bieżącego modułu | Kolumna dominująca układu, między kolumną komunikacji a kolumnami paneli pomocniczych, pod kartami sesji | Pusty (nowa karta) · przeładowanie (przejście między modułami) · wypełniony (moduł aktywny) | Wybór modułu w bocznej nawigacji zamienia całą zawartość kontenera (rozdz. 3.5.1, Diagram 4) | Każda karta sesji, wszystkie 3 środowiska modułowe | 1 | Widoczny bez interakcji |
| Stan pusty | Komunikat pustego widoku | Poinformowanie, że karta czeka na wybór modułu | `.dn-pusty-stan` (ikona + tytuł + opis), wyśrodkowany w kontenerze | Widoczny wyłącznie dla nowo otwartej karty bez wybranego modułu | Wybór dowolnej pozycji bocznej nawigacji zastępuje stan pusty zestawem okien modułu | Nowa karta sesji (kontrolka „+”) | 1 | Widoczny bez interakcji |
| Wskaźnik przeładowania | Krótki stan przejściowy między modułami | Sygnalizacja, że zestaw okien się wymienia | `.dn-spinner` lub przejście animowane (`--dn-czas-2`, respektuje `prefers-reduced-motion`) | Trwa krótko, automatycznie ustępuje stanowi wypełnionemu | Automatyczne — brak interakcji użytkownika | Przy każdej zmianie modułu w bocznej nawigacji | 1 | Widoczny bez interakcji |

#### 3.5.1. Chat Window — kanał Użytkownik ↔ Wykonawca

Chat Window jest głównym oknem komunikacji między Użytkownikiem a Wykonawcą (AI, agent lub system wykonawczy) i **jedynym oknem operacyjnym wspólnym dla wszystkich piętnastu modułów** — nie jest specyficzny dla żadnego z nich, więc jego elementy są w zakresie tego dokumentu (Specyfikacja okien operacyjnych, rozdz. 5). Stanowi centralny punkt pracy użytkownika i podstawowy mechanizm sterowania wszystkimi procesami realizowanymi przez platformę: przyjmuje polecenia w języku naturalnym, prezentuje strumień odpowiedzi i wyników, umożliwia zatwierdzanie oraz przerywanie działań, a także wyjaśnianie wyniku i kontekstu. Każde środowisko i każdy moduł udostępnia to okno w tym samym miejscu układu — **lewa kolumna, stała, pełna wysokość obszaru roboczego**. Chat Window należy w całości do warstwy 1 widoczności (rozdz. 7) i jest jednocześnie kanałem dostępu do funkcji warstwy 4: polecenie języka naturalnego wywołuje operacje niewidoczne w interfejsie.

**Makieta 15 — Chat Window (lewa kolumna, stan spoczynku)**

```
┌─────────────────────────────┐
│ CHAT WINDOW                  │
│ Użytkownik ↔ Wykonawca       │
├─────────────────────────────┤
│ historia rozmowy (przewijana)│
│ ┌─────────────────────────┐ │
│ │ Użytkownik: …            │ │
│ └─────────────────────────┘ │
│ ┌─────────────────────────┐ │
│ │ Wykonawca: … (strumień)  │ │
│ │                    ↩ ⋮   │ │
│ └─────────────────────────┘ │
├─────────────────────────────┤
│ [Fable 5] [Ultra] [Worktree] │ ← znaczniki kontekstowe (warstwa 2)
│ ┌─────────────────────────┐ │
│ │ pole poleceń…            │ │
│ └─────────────────────────┘ │
│ 📎   Operacje ▼          ➤   │ ← element zbiorczy (warstwa 3)
└─────────────────────────────┘
```

**Tabela 13 — Elementy Chat Window**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Historia rozmowy | Przewijana lista wiadomości | Prezentacja dotychczasowego przebiegu rozmowy w bieżącym module | Obszar przewijany, treść krojem bazowym; cytowania kodu `.dn-kod` krojem mono | Pusta (nowa karta/moduł) · z treścią · przewijana | Przewijanie w górę pokazuje starsze wiadomości; domyślnie odrębna per karta sesji, współdzielenie konfigurowalne (poza zakresem) | Chat Window w każdym module, wszystkie 4 środowiska | 1 | Widoczna bez interakcji |
| Pole wprowadzania poleceń | Pole tekstowe wieloliniowe | Wpisanie polecenia/pytania do AI | `.dn-input`/`.dn-textarea`, w stopce kolumny Chat Window, pełna szerokość kolumny | Puste (placeholder) · wypełnione · fokus | Pole pozostaje edytowalne również podczas generowania odpowiedzi (można przygotować kolejne polecenie w tym czasie); Enter/klik „Wyślij” wysyła treść — treść trafia do historii jako wiadomość użytkownika | Chat Window, wszystkie moduły | 1 | Widoczne bez interakcji |
| Przycisk „Wyślij” | Ikona strzałki/samolotu | Wysłanie wpisanego polecenia | `.dn-btn-ikona`, ikona `wyslij`, przy polu poleceń | Domyślny · hover · aktywny | Przycisk zawsze klikalny; przy pustym polu klik pokazuje krótkie ostrzeżenie zamiast wysyłki; w przeciwnym razie wysyła treść pola do AI i czyści pole | Chat Window, wszystkie moduły | 1 | Widoczny bez interakcji |
| Przycisk „Załącznik” | Ikona spinacza | Dołączenie pliku do polecenia | `.dn-btn-ikona`, ikona `spinacz` | Domyślny · hover · aktywny | Klik otwiera wybór pliku; załącznik pojawia się jako podgląd nad polem poleceń | Chat Window, wszystkie moduły | 1 | Widoczny bez interakcji |
| Kontrolka „Odpowiedz/cytuj” | Mała ikona przy wiadomości | Odniesienie się do konkretnej wcześniejszej wiadomości | `.dn-btn-ikona`, ikona `odpowiedz`, widoczna przy najechaniu na wiadomość | Ukryta domyślnie · widoczna na hover wiadomości | Klik wstawia cytat do pola poleceń | Chat Window, wszystkie moduły | 2 | Ujawniana najechaniem na wiadomość |
| Strumień odpowiedzi modelu | Obszar odpowiedzi generowanej na żywo | Prezentacja odpowiedzi AI w miarę jej powstawania | Blok wiadomości w historii, aktualizowany kanałem WebSocket | Generowanie (tekst narasta, wskaźnik aktywności) · zakończone | Automatyczne — odpowiedź dopisuje się do historii w czasie rzeczywistym | Chat Window, wszystkie moduły | 1 | Widoczny bez interakcji |
| Dostęp do operacji modułu | Mechanizm/przycisk kontekstowy | Wywołanie operacji AI właściwej bieżącemu modułowi (np. „zastosuj do zaznaczonego fragmentu” w Studio) | Zależny od modułu — konkretna forma poza zakresem tego dokumentu | Zależne od modułu | Zależne od modułu — patrz dokument modułu | Chat Window, treść operacji zmienna per moduł | 3 | Rozwinięcie elementu zbiorczego „Operacje ▼” |
| Element zbiorczy „Operacje ▼” | Jeden przycisk zwinięty zastępujący zestaw przycisków akcji | Grupowanie logiczne akcji Chat Window (zastosuj do zaznaczenia, przekaż do modułu, zapisz wynik, ponów, wyczyść kontekst) | Przycisk `.dn-btn--zarys` ze znacznikiem `▼`, przy polu poleceń | Zwinięty (spoczynek) · rozwinięty (lista akcji) · fokus | Klik rozwija pełną listę akcji; wybór akcji zwija element samoczynnie | Chat Window, wszystkie moduły | 3 | Klik przycisku „Operacje ▼” — jedno kliknięcie do każdej akcji |
| Menu kebab wiadomości (⋮) | Menu rozwijane pojedynczej wiadomości | Dostęp do wariantów operacji na wiadomości: kopiuj, cytuj, rozwiń kontekst, ponów odpowiedź, zgłoś wynik | Ikona `⋮` w narożniku bloku wiadomości | Ukryte w spoczynku · widoczne na hover/fokus wiadomości · otwarte | Klik otwiera menu; wybór pozycji zamyka menu | Chat Window, wszystkie moduły | 3 | Klik ikony ⋮ w narożniku wiadomości |
| Pasek znaczników kontekstowych | Lekkie znaczniki środowiska, repozytorium, projektu, modelu i wykonawcy | Prezentacja i zmiana kontekstu pracy bez otwierania okna konfiguracji, np. `[Danaco Console] [Ubuntu] [Worktree] [Fable 5] [Ultra]` | Rząd plakietek `.dn-plakietka` nad polem poleceń | Spoczynek (zwinięte znaczniki) · znacznik aktywny · selektor otwarty | Klik znacznika otwiera właściwy selektor (model, wykonawca, środowisko, tryb pracy, poziom wysiłku); po wyborze selektor zwija się samoczynnie | Chat Window, wszystkie moduły | 2 | Widoczne w spoczynku jako znaczniki zwinięte; klik znacznika otwiera selektor |
| Kontrolki zatwierdzenia i przerwania | Para przycisków sterujących przebiegiem | Zatwierdzenie działania zaproponowanego przez Wykonawcę oraz przerwanie działania w toku | `.dn-btn--zloty` (zatwierdź) i `.dn-btn--zarys` (przerwij), w bloku wiadomości wymagającej decyzji | Ukryte, gdy decyzja nie jest wymagana · widoczne przy oczekiwaniu na decyzję · widoczne w trakcie generowania (przerwij) | Klik „Zatwierdź” uruchamia działanie; klik „Przerwij” zatrzymuje strumień i zapisuje stan częściowy w historii | Chat Window, wszystkie moduły | 2 | Ujawniane zdarzeniem: oczekiwaniem na decyzję albo generowaniem odpowiedzi |
| Wyszukiwarka funkcji | Pole wywoływane skrótem klawiszowym | Dotarcie do dowolnej funkcji platformy, w tym funkcji warstwy 4, bez przeszukiwania menu | Panel popover wyśrodkowany nad obszarem roboczym, pole tekstowe z listą trafień | Zamknięta (spoczynek) · otwarta · z wynikami · bez wyników | Skrót klawiszowy otwiera panel; wpisanie frazy filtruje listę funkcji; wybór uruchamia funkcję i zamyka panel | Chat Window i cała powłoka, wszystkie 4 środowiska | 4 | Skrót klawiszowy; wpisana fraza filtruje listę funkcji |

**Diagram 4 — Przeładowanie obszaru roboczego przy zmianie modułu**

```
Boczna nawigacja: klik pozycji modułu B (poprzednio aktywny: moduł A)
        │
        ▼
Znikają okna operacyjne modułu A                      Chat Window
        │                                              NIE znika —
        ▼                                              rekonfiguruje
Pojawia się zestaw okien operacyjnych modułu B         kontekst
        │                                              (rozdz. 2.4
        ▼                                              Koncepcji)
Chat Window rekonfiguruje kontekst poleceń i dostępnych operacji
do nowego modułu B (historia/pamięć zależnie od ustawień izolacji)
        │
        ▼
Karta sesji: tytuł karty aktualizuje się na nazwę modułu B
Boczna nawigacja: stan wybrany przenosi się na pozycję modułu B

NIE ULEGA ZMIANIE: sama karta sesji, boczna nawigacja jako panel,
pozostałe otwarte karty, funkcje globalne Mobile / Always On Display
```

#### 3.5.2. Execution Loop Window — kanał Koordynator ↔ Wykonawca

Execution Loop Window jest drugim kanałem komunikacji operacyjnej platformy i elementem pierwszoplanowym architektury każdego środowiska. Prezentuje komunikację między Koordynatorem — komponentem orkiestrującym platformy — a Wykonawcą. Odpowiada za prowadzenie pętli wykonawczej, koordynację zadań, nadzór nad realizacją, orkiestrację działań oraz kontrolę realizacji procesów. Okno otwierane jest jako **kolumna sąsiadująca z Chat Window**, pod nim w tej samej lewej kolumnie układu, i zachowuje pełną wysokość przypisanej mu przestrzeni. Nagłówek okna należy do warstwy 1, jego rozwinięcia sterujące — do warstw 2 i 3 (rozdz. 7).

**Makieta 15a — Execution Loop Window**

```
┌─────────────────────────────┐
│ EXECUTION LOOP WINDOW        │
│ Koordynator ↔ Wykonawca      │
├─────────────────────────────┤
│ ZLECENIE: „Budowa modułu”    │
│ ├ zadanie 1  ● zakończone    │
│ ├ zadanie 2  ◐ w toku        │
│ └ zadanie 3  ○ w kolejce     │
├─────────────────────────────┤
│ Komunikaty sterujące          │
│  Koordynator → Wykonawca: …   │
│  Wykonawca → Koordynator: …   │
├─────────────────────────────┤
│ Kontrola jakości: ● zgodne    │
│ Ponowienia: 1                 │
├─────────────────────────────┤
│ Przebieg pętli: ▓▓▓▓▓░░░ 62 % │
│ Sterowanie ▼                  │ ← element zbiorczy (warstwa 3)
└─────────────────────────────┘
```

**Tabela 13a — Elementy Execution Loop Window**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Nagłówek okna pętli | Pasek tytułowy „Execution Loop Window — Koordynator ↔ Wykonawca” | Jednoznaczna identyfikacja drugiego kanału komunikacji | Nagłówek kolumny, krój bazowy semibold, kreska rozdzielająca od Chat Window | Stały | Brak interakcji własnej | Kolumna komunikacji, wszystkie 4 środowiska | 1 | Widoczny bez interakcji |
| Blok bieżącego zlecenia | Nazwa zlecenia wraz z jego dekompozycją na zadania | Prezentacja tego, co Koordynator zlecił Wykonawcy i jak zostało rozłożone na zadania | Lista zagnieżdżona `.dn-karta--pozycja` ze wskaźnikami stanu | Brak zlecenia · zlecenie aktywne · zlecenie zakończone | Klik zadania rozwija jego szczegóły w panelu popover | Kolumna komunikacji, wszystkie 4 środowiska | 1 | Widoczny bez interakcji |
| Kolejka i stan zadań | Lista zadań ze stanem wykonania | Podgląd kolejności realizacji i bieżącego stanu każdego zadania | `.dn-kropka` semantyczna + tekst zadania | W kolejce · w toku · zakończone · odrzucone | Klik zadania przenosi kontekst wymiany komunikatów do tego zadania | Kolumna komunikacji, wszystkie 4 środowiska | 1 | Widoczna bez interakcji |
| Strumień komunikatów sterujących | Przewijana wymiana komunikatów Koordynator ↔ Wykonawca | Podgląd poleceń sterujących, potwierdzeń i raportów Wykonawcy | Obszar przewijany, wiersze z oznaczeniem nadawcy | Pusty · z treścią · aktualizacja na żywo kanałem WebSocket | Przewijanie pokazuje wcześniejsze komunikaty; wiersz rozwija pełną treść | Kolumna komunikacji, wszystkie 4 środowiska | 1 | Widoczny bez interakcji |
| Wynik kontroli jakości | Blok statusu weryfikacji zadania | Prezentacja rozstrzygnięcia kontroli jakości i decyzji o ponowieniu | `.dn-plakietka--stan` + licznik ponowień | Zgodne · niezgodne (ponowienie) · oczekuje | Klik plakietki otwiera panel popover ze szczegółami niezgodności | Kolumna komunikacji, wszystkie 4 środowiska | 2 | Widoczny jako plakietka; klik otwiera panel popover ze szczegółami |
| Wskaźnik przebiegu pętli | Pasek postępu z wartością procentową | Sygnalizacja zaawansowania bieżącego przebiegu pętli wykonawczej | Pasek postępu w stopce kolumny | Bezczynny · w toku · wstrzymany · zakończony | Automatyczny — brak interakcji użytkownika | Kolumna komunikacji, wszystkie 4 środowiska | 1 | Widoczny bez interakcji |
| Element zbiorczy „Sterowanie ▼” | Jeden przycisk zwinięty grupujący akcje sterowania przebiegiem | Wstrzymanie, wznowienie, przerwanie oraz korekta zlecenia | Przycisk `.dn-btn--zarys` ze znacznikiem `▼`, stopka kolumny | Zwinięty · rozwinięty · fokus | Klik rozwija listę akcji sterujących; wybór akcji zwija element samoczynnie | Kolumna komunikacji, wszystkie 4 środowiska | 3 | Klik przycisku „Sterowanie ▼” |
| Kontrolka otwarcia i zamknięcia okna pętli | Przełącznik kolumny Execution Loop Window | Otwarcie kolumny pętli wykonawczej i jej zamknięcie | Ikona w nagłówku Chat Window oraz w pasku górnym | Zamknięte (kolumna nieobecna) · otwarte | Klik otwiera kolumnę sąsiadującą; ponowny klik zamyka ją, a kolumna znika całkowicie z przestrzeni roboczej | Nagłówek Chat Window i pasek górny, wszystkie 4 środowiska | 2 | Klik ikony w nagłówku Chat Window albo w pasku górnym |

#### 3.5.3. Panele pomocnicze — rozszerzenia boczne

Okna pomocnicze i panele otwierane są jako **kolumny boczne po prawej stronie obszaru roboczego**. Po zamknięciu panel znika całkowicie z przestrzeni roboczej, nie pozostawiając wyzwalacza innego niż ikona lub znacznik kontekstowy, z którego został wywołany.

**Tabela 13b — Elementy paneli pomocniczych**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Panel wysuwany | Kolumna boczna otwierana z prawej krawędzi obszaru roboczego | Pomieszczenie funkcji potrzebnych chwilowo: szczegółów elementu, historii wersji, listy narzędzi | Kolumna boczna o regulowanej szerokości, pełna wysokość obszaru roboczego | Zamknięty (nieobecny) · otwarty · szerokość regulowana | Klik wyzwalacza otwiera kolumnę; zamknięcie usuwa ją z przestrzeni roboczej bez pozostawiania śladu | Prawa strona obszaru roboczego, wszystkie 4 środowiska | 3 | Klik ikony lub znacznika kontekstowego |
| Panel popover | Mały panel zakotwiczony przy elemencie wywołującym | Prezentacja szczegółów i ustawień szybkich bez zmiany układu kolumn | Panel unoszony `.dn-popover`, szerokość dopasowana do treści | Zamknięty · otwarty · fokus wewnątrz | Klik poza panelem zamyka go; wybór pozycji zamyka panel i stosuje ustawienie | Przy znacznikach kontekstowych, plakietkach stanu i pozycjach list | 3 | Klik elementu kotwiczącego |
| Menu hamburger (☰) | Zwinięty wyzwalacz nawigacji | Otwarcie bocznej nawigacji modułów albo panelu orkiestracji, gdy szerokość okna nie mieści kolumny nawigacji | Ikona `menu` w lewym narożniku paska górnego | Zwinięte · rozwinięte | Klik rozwija kolumnę nawigacji jako panel wysuwany; wybór pozycji zwija ją samoczynnie | Pasek górny, wszystkie 4 środowiska | 3 | Klik ikony ☰ |
| Menu kebab powłoki (⋮) | Menu rozwijane paska górnego | Dostęp do zestawu akcji powłoki: konfiguracja sesji, motyw, Mobile, Always On Display, tryb administracyjny | Ikona `⋮` w prawej części paska górnego | Zwinięte · rozwinięte · pozycja podświetlona | Klik rozwija listę akcji; wybór pozycji zwija menu | Pasek górny, wszystkie 4 środowiska | 3 | Klik ikony ⋮ |


---

## 4. Powłoka środowiska MultitaskingAI

### 4.1. Różnice względem powłoki modułowej

**Tabela 14 — MultitaskingAI wobec powłoki środowisk modułowych**

| Cecha | Powłoka modułowa (TalkIn / WorkSpace / CodeStudio) | Powłoka MultitaskingAI |
|---|---|---|
| Co organizuje boczny panel | Zadania modułowe — lista modułów | Zespół modeli i agentów realizujących wspólny proces |
| Zawartość bocznego panelu | Lista modułów wg macierzy dostępności (rozdz. 3.4) | Panel orkiestracji — 6 sekcji stałych (rozdz. 4.4) |
| Po czym się nawiguje | Po modułach | Po rolach, kolejkach i orkiestracji |
| Co reprezentuje karta sesji | Jeden moduł otwarty w danej chwili | Jeden proces orkiestracji — stan ról, kolejek i orkiestracji widoczny przez panel (rozdz. 4.3) |
| Zawartość obszaru roboczego | Okna operacyjne wybranego modułu + Chat Window | Zawartość wybranej sekcji panelu orkiestracji (np. 4 karty ról w sekcji „Role”) |
| Pasek górny | Wzorzec wspólny `.dn-pasek` (rozdz. 3.2) | Ten sam wzorzec wspólny — bez zmian (rozdz. 4.2) |

### 4.2. Układ całościowy

**Makieta 16 — Powłoka środowiska MultitaskingAI, przykład wypełniony (sekcja „Role” wybrana)**

```
╔═══════════════════════════════════════════════════════════════════════╗
║ ◈ Danaco Console  🏠  [ 🔍 Szukaj… ]  [Ubuntu][Fable 5][Ultra]   ⋮ ☀  ║
╠═══════════════════════════════════════════════════════════════════════╣
║ [●  Proces „Budowa aplikacji” ✕] [ + ]                                 ║
╠════════════════╦═════════════╦════════════════════════╦═══════════════╣
║ MultitaskingAI ║ CHAT WINDOW ║ Executor 1 │ Executor 2 ║ Panel         ║
║ Panel          ║ Użytkownik ↔║ (👤 agent) │ (👤 agent) ║ pomocniczy    ║
║ orkiestracji   ║ Wykonawca   ║────────────┼────────────║ (zwinięty)    ║
║────────────────║ (per rola)  ║ Coordinator│ Executor 3 ║               ║
║   Zespoły      ║             ║ (👤 agent) │ / Validator║               ║
║ ▸ Role         ║ [pole       ║            │ (👤 agent) ║               ║
║   Kolejki      ║  poleceń…]  ║                         ║               ║
║   Orkiestracja ║ 📎 Operacje▼➤║ [ wnętrze okien        ║               ║
║   Harmonogram  ║─────────────║   roboczych ról —      ║               ║
║   i automatyki ║ EXECUTION   ║   poza zakresem tego   ║               ║
║   Monitor      ║ LOOP WINDOW ║   dokumentu ]          ║               ║
║   procesu      ║ Koordynator ║                         ║               ║
║                ║ ↔ Wykonawca ║                         ║               ║
╚════════════════╩═════════════╩════════════════════════╩═══════════════╝
```

W środowisku MultitaskingAI Execution Loop Window pełni rolę wiodącą: pętla wykonawcza Koordynator ↔ Wykonawca jest przedmiotem pracy całego środowiska, a sekcja „Monitor procesu” panelu orkiestracji prezentuje jej przebieg w skali wszystkich procesów (rozdz. 4.4).

### 4.3. Pasek górny i karty sesji — dostosowania

**Tabela 15 — Różnice elementów w MultitaskingAI**

| Element | Zachowanie w powłoce modułowej | Zachowanie w powłoce MultitaskingAI |
|---|---|---|
| Pasek górny | Wszystkie elementy z Tabeli 8 | Identyczny zestaw i mechanika — bez zmian (rozdz. 3.2) |
| Tytuł karty sesji | Nazwa modułu otwartego w karcie | Nazwa procesu/zespołu orkiestracji (np. nazwa zapisanej konfiguracji zespołu — rozdz. 4.4, sekcja „Zespoły”) |
| Wskaźnik stanu karty | Sygnalizacja procesu w tle danego modułu | Sygnalizacja przebiegu pętli orkiestracji — silniej powiązany z Monitorem procesu (rozdz. 4.4) i z funkcją Mobile |
| Kontrolka nowej karty „+” | Otwiera pustą kartę, oczekuje wyboru modułu | Otwiera pusty proces orkiestracji, oczekuje konfiguracji zespołu/ról w panelu orkiestracji |
| Poziom zasięgu izolacji karty | „Karta sesji” — poziom szósty z siedmiu (Koncepcja platformy, rozdz. 6.5) | Karta reprezentuje zasięg „karta sesji”, ale zawiera zagnieżdżony poziom „Rola” — najwęższy z siedmiu, z pierwszeństwem nad ustawieniem karty (rozdz. 6.5 Koncepcji; poza zakresem tego dokumentu — okno konfiguracji punktów izolacji) |

### 4.4. Panel orkiestracji

**Makieta 17 — Panel orkiestracji, boczna nawigacja MultitaskingAI**

```
  ┌────────────────────────────┐
  │ MultitaskingAI              │   ← nagłówek: nazwa środowiska
  │ Panel orkiestracji          │   ← podtytuł stały
  ├────────────────────────────┤
  │   Zespoły                   │
  │ ▸ Role                      │  ← stan wybrany
  │   Kolejki                   │
  │   Orkiestracja               │
  │   Harmonogram i automatyki   │
  │   Monitor procesu            │
  └────────────────────────────┘
    kolejność i widoczność sekcji — konfigurowalne (wartość domyślna
    jak powyżej; okno konfiguracji tej listy — poza zakresem)
```

**Tabela 16 — Elementy panelu orkiestracji (6 sekcji)**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Nagłówek panelu | „MultitaskingAI” + „Panel orkiestracji” | Orientacja — potwierdzenie środowiska i charakteru panelu | Krój nagłówkowy dla nazwy środowiska, podtytuł krojem bazowym | Stały | Brak interakcji | Boczny panel MultitaskingAI | 1 | Widoczny bez interakcji |
| Sekcja „Zespoły” | Pozycja listy | Otwarcie widoku zapisanych konfiguracji zespołu (presety ról, powiązań, kolejek) | `.dn-karta--pozycja`, ikona `gwiazdka`/`plik` | Domyślny · hover · wybrany · fokus | Klik pokazuje w obszarze roboczym listę zapisanych zespołów z akcjami zapisu/wczytania/duplikowania | Panel orkiestracji | 1 | Widoczna bez interakcji |
| Sekcja „Role” | Pozycja listy | Otwarcie czterech okien roboczych ról i przypisania agentów | `.dn-karta--pozycja`, ikona `uzytkownik` | Domyślny · hover · wybrany (domyślnie aktywna w przykładzie) · fokus | Klik pokazuje w obszarze roboczym 4 karty ról: Executor 1, Executor 2, Coordinator, Executor 3 / Validator (rozdz. 4.5) | Panel orkiestracji | 1 | Widoczna bez interakcji |
| Sekcja „Kolejki” | Pozycja listy | Otwarcie definicji kolejek (globalne, lokalne, modeli, agentów, projektów) i ich akcji | `.dn-karta--pozycja`, ikona `filtr`/`strzalka-prawo` | Domyślny · hover · wybrany · fokus | Klik pokazuje tabelę kolejek z akcjami enqueue/dequeue/retry jako przyciski ikonowe | Panel orkiestracji | 1 | Widoczna bez interakcji |
| Sekcja „Orkiestracja” | Pozycja listy | Otwarcie zależności między modelami, agentami, zadaniami, kolejkami, automatyzacjami, projektami | `.dn-karta--pozycja`, ikona `link-zewnetrzny` | Domyślny · hover · wybrany · fokus | Klik pokazuje listę zależności ze statusem zgodności (plakietka) | Panel orkiestracji | 1 | Widoczna bez interakcji |
| Sekcja „Harmonogram i automatyki” | Pozycja listy | Otwarcie harmonogramu pracy ciągłej i wpiętych automatyk z modułu Automations | `.dn-karta--pozycja`, ikona `zegar`/`kalendarz` | Domyślny · hover · wybrany · fokus | Klik pokazuje listę wpiętych automatyk i harmonogram — wzorzec Scheduler (moduł Automations, poza zakresem) | Panel orkiestracji | 1 | Widoczna bez interakcji |
| Sekcja „Monitor procesu” | Pozycja listy | Otwarcie podglądu przebiegu pętli, statusów przebiegów, hierarchii decyzji | `.dn-karta--pozycja`, ikona `oko` | Domyślny · hover · wybrany · fokus | Klik pokazuje tabelę statusów procesów z plakietkami stanu; miejsce nadzoru funkcji Always On Display | Panel orkiestracji | 1 | Widoczna bez interakcji |

### 4.5. Obszar roboczy MultitaskingAI

Analogicznie do rozdz. 3.5, obszar roboczy jest kontenerem, którego zawartość zależy od wybranej sekcji panelu orkiestracji. Poniżej opisano zachowanie na poziomie kontenera — wnętrze poszczególnych okien roboczych ról oraz managerów (Queue Manager, Orchestrator, Scheduler) jest przedmiotem dokumentu środowiska MultitaskingAI i dokumentu modułu Automations (poza zakresem).

**Makieta 18 — Obszar roboczy, sekcja „Role” wybrana**

```
┌───────────────────────────────────────────────────────────┐
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐            │
│  │ 👤 Executor 1│ │ 👤 Executor 2│ │ 👤 Coordinator│           │
│  │ agent: „X”   │ │ agent: „Y”   │ │ model: „Z”   │           │
│  │ ● aktywny     │ │ ○ nieaktywny │ │ ● aktywny     │           │
│  └─────────────┘ └─────────────┘ └─────────────┘            │
│  ┌─────────────┐                                             │
│  │ 👤 Executor 3│                                             │
│  │ / Validator  │                                             │
│  │ agent: „W”   │                                             │
│  └─────────────┘                                             │
│  [ wnętrze okna roboczego roli po kliknięciu karty —           │
│    poza zakresem tego dokumentu ]                             │
└───────────────────────────────────────────────────────────┘
```

**Makieta 19 — Obszar roboczy, sekcja „Monitor procesu” wybrana (przykład innej sekcji)**

```
┌───────────────────────────────────────────────────────────┐
│  Przebieg          Rola            Status        Czas       │
│  ──────────────────────────────────────────────────────────│
│  #128              Executor 1      ● sukces      12:04       │
│  #129              Coordinator     ◐ w toku       12:05       │
│  #130              Executor 2      ○ oczekuje     —          │
│                                                             │
│  [ szczegóły przebiegu, hierarchia decyzji —                 │
│    poza zakresem tego dokumentu ]                             │
└───────────────────────────────────────────────────────────┘
   `.dn-tabela` + `.dn-plakietka--stan` (wzorzec monitora, jak
   w oknach modułu Automations — Execution Monitor, Process Monitor)
```

**Tabela 17 — Elementy kontenera obszaru roboczego (poziom powłoki)**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Obszar roboczy (kontener) | Główna powierzchnia powłoki MultitaskingAI | Pomieszczenie zawartości wybranej sekcji panelu orkiestracji | Kolumna dominująca układu, między kolumną komunikacji a kolumnami paneli pomocniczych, pod kartami sesji | Pusty (nowy proces, brak skonfigurowanych ról) · przeładowanie (zmiana sekcji) · wypełniony | Klik sekcji panelu zamienia całą zawartość kontenera | Każda karta sesji środowiska MultitaskingAI | 1 | Widoczny bez interakcji |
| Karta roli (w sekcji „Role”) | Kafel reprezentujący jedną z czterech ról | Podgląd skrócony roli (przypisany agent/model, status) i wejście w jej okno robocze | `.dn-karta` z `.dn-awatar` (wariant `--kwadrat` dla agenta przypisanego, kołowy dla modelu bazowego) | Bez przypisanego wykonawcy · z przypisanym agentem/modelem · aktywny (praca w toku) · nieaktywny | Klik otwiera pełne okno robocze roli (Executor 1, Executor 2, Coordinator, Executor 3 / Validator) — **poza zakresem tego dokumentu** | Obszar roboczy, sekcja „Role” — dokładnie 4 karty | 1 | Widoczna bez interakcji |
| Rozwinięcie Subagent Network | Zagnieżdżona lista podagentów pod kartą wykonawcy | Podgląd do 15 podagentów uruchamianych przez rolę Executor | Zagnieżdżone pozycje `.dn-karta--pozycja` pod kartą roli, nie osobna sekcja panelu | Zwinięte · rozwinięte · pojedynczy podagent aktywny/zakończony | Klik roli z uruchomionym Subagent Network rozwija listę podagentów | Karta roli Executor 1 / Executor 2, gdy mechanizm aktywny | 3 | Klik karty roli z uruchomionym mechanizmem rozwija listę podagentów |
| Zawartość innych sekcji (Zespoły, Kolejki, Orkiestracja, Harmonogram, Monitor) | Widok listy/tabeli właściwy sekcji | Prezentacja danych organizacyjnych warstwy orkiestracji | `.dn-tabela` lub `.dn-karta--pozycja`, zależnie od sekcji (Tabela 16) | Pusty stan (brak konfiguracji) · z zawartością | Zależne od sekcji — szczegóły managerów (Queue Manager, Orchestrator, Scheduler, Execution Monitor) poza zakresem, patrz dokument modułu Automations | Obszar roboczy, sekcje inne niż „Role” | 1 | Widoczna bez interakcji |

**Diagram 5 — Przełączanie sekcji panelu orkiestracji → zawartość obszaru roboczego**

```
Panel orkiestracji: klik sekcji „Kolejki” (poprzednio: „Role”)
        │
        ▼
Znika zawartość sekcji „Role” (4 karty ról)
        │
        ▼
Pojawia się zawartość sekcji „Kolejki” (tabela kolejek + akcje)
        │
        ▼
Karta sesji: tytuł pozostaje nazwą procesu (nie zmienia się przy
zmianie sekcji — w odróżnieniu od zmiany modułu w środowiskach
modułowych, rozdz. 3.5.1, Diagram 4)
        │
        ▼
Panel orkiestracji: stan wybrany przenosi się na sekcję „Kolejki”
```

### 4.6. MultitaskingAI wobec środowisk modułowych — zestawienie strukturalne

**Schemat 6 — Zestawienie strukturalne**

```
ŚRODOWISKA MODUŁOWE                        MultitaskingAI
(TalkIn · WorkSpace · CodeStudio)
┌──────────────────────────┐               ┌──────────────────────────┐
│ pasek górny               │               │ pasek górny               │  ← identyczny
├──────────────────────────┤               ├──────────────────────────┤
│ karty sesji = moduły       │               │ karty sesji = procesy     │  ← inna semantyka
│                            │               │ orkiestracji              │     karty
├──────┬──────┬────────────┤               ├──────┬──────┬────────────┤
│boczna│ Chat │ obszar      │               │panel │ Chat │ obszar      │  ← inny panel
│nawig.│Window│ roboczy =   │               │orkie-│Window│ roboczy =   │     nawigacji,
│= li- │ ──── │ okna        │               │stracji│ ──── │ zawartość   │     ta sama
│sta   │Execu-│ modułu      │               │= 6   │Execu-│ wybranej    │     kolumnowa
│modu- │tion  │             │               │sekcji│tion  │ sekcji      │     pozycja
│łów   │Loop  │             │               │      │Loop  │             │     obu okien
│      │Window│             │               │      │Window│             │     komunikacji
└──────┴──────┴────────────┘               └──────┴──────┴────────────┘
```

## 5. Tabela zbiorcza dokumentu

**Tabela 18 — Wszystkie okna/powłoki i ich elementy główne**

| Okno / powłoka | Elementy główne (liczba pozycji w tabelach elementów) | Rozdział | Tabele | Makiety |
|---|---|---|---|---|
| Okno startowe / logowania | Pola formularza (4), przyciski (6), linki (3), przełącznik metody (1), komunikaty stanu (3) — 17 pozycji w Tabeli 1 | 1 | 2 | 4 |
| Strona główna — Strefa 1 (karty środowisk) | Karta jako całość + 3 podelementy (godło, tytuł, opis) — 4 pozycje w Tabeli 3, treść 4 kart w Tabeli 4 | 2.3 | 2 | 1 |
| Strona główna — Strefa 2 (kafle komponentów) | Kafel jako całość + 3 podelementy — 4 pozycje w Tabeli 5, treść 4 kafli w Tabeli 6 | 2.4 | 2 | 1 |
| Strona główna — Strefa 3 (listwa ustawień) | 3 pozycje listwy — Tabela 7 | 2.5 | 1 | 1 |
| Powłoka modułowa — pasek górny | 7 elementów — Tabela 8 | 3.2 | 1 | — |
| Powłoka modułowa — karty sesji | 5 elementów — Tabela 9 | 3.3 | 1 | 2 |
| Powłoka modułowa — boczna nawigacja | 3 elementy ogólne (Tabela 10) + macierz 3 środowisk × moduły (Tabela 11) | 3.4 | 2 | 1 |
| Powłoka modułowa — obszar roboczy (kontener) | 3 elementy kontenera — Tabela 12 | 3.5 | 1 | 2 |
| Powłoka modułowa — Chat Window (Użytkownik ↔ Wykonawca) | 12 elementów — Tabela 13 | 3.5.1 | 1 | 1 |
| Powłoka modułowa — Execution Loop Window (Koordynator ↔ Wykonawca) | 8 elementów — Tabela 13a | 3.5.2 | 1 | 1 |
| Powłoka modułowa — panele pomocnicze | 4 elementy — Tabela 13b | 3.5.3 | 1 | — |
| Powłoka MultitaskingAI — różnice i pasek górny/karty | Tabela różnic (Tabela 14) + tabela dostosowań (Tabela 15) | 4.1, 4.3 | 2 | — |
| Powłoka MultitaskingAI — panel orkiestracji | 6 sekcji + nagłówek — Tabela 16 | 4.4 | 1 | 1 |
| Powłoka MultitaskingAI — obszar roboczy | 4 elementy kontenera — Tabela 17 | 4.5 | 1 | 2 |
| Warstwy widoczności elementów interfejsu | 4 warstwy (Tabela 21) + 7 typów elementów warstw 2–4 (Tabela 22) | 7 | 2 | — |

---

## 6. Stany wspólne i słowniczek

### 6.1. Słownik stanów interakcji

**Tabela 19 — Stany wspólne wszystkich elementów interaktywnych (System wizualny, rozdz. 6)**

| Stan | Co oznacza wizualnie | Kiedy występuje |
|---|---|---|
| Domyślny (spoczynek) | Wygląd bazowy elementu, bez modyfikacji | Element widoczny, brak interakcji |
| Hover (najechanie) | Subtelna zmiana tła albo uniesienie 1–2 px dla elementów akcentowanych złotem | Kursor nad elementem interaktywnym |
| Aktywny (wciśnięcie) | Przesunięcie 1 px w dół (`translateY(1px)`) | Moment kliknięcia/naciśnięcia |
| Fokus | Pierścień złoty 2 px + odsunięcie 2 px | Nawigacja klawiaturą (`:focus-visible`) — obowiązkowy dla wszystkich kontrolek interaktywnych |
| Niegotowy / niekompletny (bez blokady interakcji) | Element pozostaje w pełni klikalny — bez wygaszenia, bez zmiany kursora, bez odebrania efektów najechania; niekompletność sygnalizowana opisowo obok elementu albo komunikatem/ostrzeżeniem po kliknięciu | Warunek wykonania niespełniony w bieżącym kontekście (np. puste pole wymagane) — klik wywołuje komunikat zamiast wykonania akcji; element nigdy nie traci klikalności |
| Zaznaczony / wybrany / aktywny (stan trwały) | Akcent złoty: podkreślenie zakładki, obrys karty, tło `--dn-akcent-tlo` | Pozycja bocznej nawigacji, zakładka, segment przełącznika aktualnie wybrane |
| Ładowanie | `.dn-spinner`, obrót ciągły (respektuje `prefers-reduced-motion`) | Oczekiwanie na odpowiedź serwera (logowanie, zapis, generowanie odpowiedzi AI) |
| Błąd | Kreska/obrys czerwony (`--dn-blad-fg/bg/br`), komunikat z ikoną `blad` | Walidacja nieudana, dane odrzucone przez serwer |

Zgodnie z zasadą nadrzędną „Zero blokerów” (przywołaną wprost w System wizualny, rozdz. 14) żaden element interaktywny nie traci klikalności: kontrolki pozostają zawsze aktywne, a niegotowość lub niekompletność danych wejściowych (np. puste pole wymagane) jest sygnalizowana opisowo obok elementu albo komunikatem/ostrzeżeniem po kliknięciu — nigdy przez odebranie interakcji.

### 6.2. Słowniczek pojęć własnych dokumentu

**Tabela 20 — Słowniczek**

| Pojęcie | Definicja | Rozdział |
|---|---|---|
| Okno startowe / logowania | Pierwsze okno prezentowane urządzeniu przy uruchomieniu platformy; łączy w jednym oknie formularz rejestracji, logowania i odzyskiwania konta, w formie zależnej od tego, czy konto istnieje i w jakiej fazie wdrożenia znajduje się platforma | 1 |
| Faza budowy / faza produkcyjna | Dwa stany wdrożenia różniące się domyślną wartością ustawienia „wymóg logowania” (wyłączony w fazie budowy, włączony w fazie produkcyjnej); ustawienie pozostaje konfigurowalne przez Operatora w każdej fazie, także w produkcyjnej — kod obsługi wspólny dla obu faz | 1.1 |
| Magistrala kontekstu | Mechanizm przekazywania artefaktu między modułami — kolumna boczna z referencjami do odłożonych artefaktów; opis rozstrzygający w `interfejs-uzytkownika/przeplyw-okien.md`, rozdz. 6a | 3.2 |
| Centrum powiadomień | Jedyny mechanizm powiadamiania platformy — kolumna boczna z rejestrem zdarzeń, wywoływana plakietką powiadomień w pasku kontekstu; karta komponentu w `interfejs-uzytkownika/katalog-komponentow.md`, rozdz. 11.6 | 3.2 |
| Paleta poleceń | Jedyna klawiaturowa nakładka platformy, wywoływana skrótem `Ctrl/Cmd + K`; karta komponentu w `interfejs-uzytkownika/katalog-komponentow.md`, rozdz. 15.8 | 3.2 |
| Centrum dowodzenia | Strona główna platformy — przedpokój przed strefą roboczą, trzy strefy: wybór środowiska, komponenty własne, ustawienia | 2 |
| Pasek górny | Poziomy pas na szczycie powłoki środowiska, wspólny wzorzec dla wszystkich czterech środowisk, niosący markę, wyszukiwanie, szybką konfigurację sesji i dostęp do funkcji globalnych | 3.2, 4.3 |
| Karta sesji | Poziomy element paska kart reprezentujący jedną, samodzielną przestrzeń roboczą (moduł — środowiska modułowe; proces orkiestracji — MultitaskingAI) | 3.3, 4.3 |
| Boczna nawigacja modułów | Stały pas pionowy z listą modułów dostępnych w środowisku, wspólny dla wszystkich kart sesji danego środowiska | 3.4 |
| Panel orkiestracji | Odpowiednik bocznej nawigacji w środowisku MultitaskingAI — 6 sekcji zamiast listy modułów | 4.4 |
| Obszar roboczy | Kontener główny powłoki, mieszczący zestaw okien operacyjnych wybranego modułu (albo zawartość wybranej sekcji panelu orkiestracji) oraz Chat Window | 3.5, 4.5 |
| Chat Window | Główne okno komunikacji Użytkownik ↔ Wykonawca; jedyne okno operacyjne wspólne wszystkim piętnastu modułom; główne okno komunikacji Użytkownik ↔ Wykonawca, umieszczone w lewej kolumnie obszaru roboczego, rekonfigurowany przy każdej zmianie modułu | 3.5.1 |
| Przeładowanie | Wymiana zestawu okien operacyjnych w obszarze roboczym bieżącej karty przy zmianie modułu/sekcji, bez zamykania samej karty | 3.5.1, 4.5 |
| Execution Loop Window | Okno pętli wykonawczej prezentujące komunikację Koordynator ↔ Wykonawca: zlecenie i jego dekompozycję, kolejkę i stan zadań, komunikaty sterujące, wyniki kontroli jakości oraz sterowanie przebiegiem; otwierane jako kolumna sąsiadująca z Chat Window | 3.5.2 |
| Koordynator | Komponent orkiestrujący platformy — dekomponuje zlecenie, przydziela i nadzoruje zadania | 3.5.2 |
| Wykonawca | AI, agent lub system wykonawczy realizujący zadania zlecone przez Użytkownika i nadzorowane przez Koordynatora | 3.5.1, 3.5.2 |
| Warstwa widoczności | Jedna z czterech warstw, do których należy każdy element interfejsu: zawsze widoczna, widoczna na żądanie, rozwinięcie kontekstowe, funkcja ekspercka | 7 |
| Znacznik kontekstowy | Lekka plakietka opisująca jeden wymiar kontekstu pracy (środowisko, repozytorium, projekt, model, wykonawca); klik otwiera właściwy selektor | 7.3, 7.5 |
| Element zbiorczy z rozwinięciem | Jeden przycisk zwinięty (`Operacje ▼`, `Sterowanie ▼`) zastępujący zestaw przycisków akcji | 7.3, 7.5 |
| Uproszczone menu kontekstowe | Skrócony dostęp z paska górnego/okna operacyjnego do ustawień warstwy sesji, bez opuszczania trwającej pracy | 3.2 |

---

## 7. Warstwy widoczności elementów interfejsu

### 7.1. Zasada nadrzędna

Zasadą nadrzędną interfejsu Danaco Console jest reguła: **jeżeli funkcja nie jest potrzebna do realizacji aktualnego zadania, nie jest widoczna**. Interfejs ujawnia możliwości systemu stopniowo — zależnie od kontekstu, roli użytkownika i wykonywanej czynności — zachowując maksymalną moc funkcjonalną przy minimalnej złożoności wizualnej. Złożoność platformy istnieje w architekturze i pozostaje niewidoczna w interfejsie do chwili wystąpienia potrzeby użycia danej funkcji. Liczba modułów, agentów, przepływów pracy, komponentów, narzędzi, paneli, ustawień i funkcji administracyjnych nie wpływa na postrzeganą prostotę interfejsu.

Każdy element interfejsu opisany w tym dokumencie należy do dokładnie jednej z czterech warstw widoczności. Tabele elementów w rozdziałach 1–4 podają dla każdego elementu jego warstwę oraz sposób wywołania.

### 7.2. Cztery warstwy widoczności

**Tabela 21 — Warstwy widoczności elementów interfejsu**

| Warstwa | Nazwa | Zawartość | Sposób dostępu |
|---|---|---|---|
| 1 | Zawsze widoczna | Chat Window, aktywne okno wiodące, kontekst pracy, podstawowa nawigacja, wskaźniki stanu wykonania. Zajmuje ponad 80% powierzchni interfejsu | Widoczna bez interakcji |
| 2 | Widoczna na żądanie | Wybór modelu, wybór wykonawcy, wybór środowiska, wybór trybu pracy, poziom wysiłku, parametry przepływu pracy | Ikona, przycisk, przełącznik, znacznik kontekstowy; po użyciu element zwija się samoczynnie |
| 3 | Rozwinięcia kontekstowe | Zestawy akcji, ustawienia szybkie, warianty operacji | Menu kebab (⋮), menu hamburger (☰), menu kontekstowe, panel popover, lista rozwijana |
| 4 | Funkcje eksperckie | Najbardziej zaawansowane operacje, tryby administracyjne, narzędzia diagnostyczne niskiego poziomu | Polecenie języka naturalnego w Chat Window, skrót klawiszowy, wyszukiwarka funkcji, tryb administracyjny, konfiguracja roli. Użytkownik podstawowy nie widzi tych elementów |

### 7.3. Mechanizmy ukrywania funkcjonalności

- **Menu progresywne.** Zbiór jednorodnych wyborów prezentowany jest jako jeden element zwinięty (`Agent ▼`), a lista pozycji rozwija się po kliknięciu.
- **Panele wysuwane.** Funkcjonalność umieszczana jest w panelach bocznych, panelach wysuwanych, oknach popover i panelach kontekstowych; po zamknięciu panel znika całkowicie z przestrzeni roboczej.
- **Grupowanie logiczne akcji.** Zamiast zestawu przycisków prezentowany jest jeden element zbiorczy (`Operacje ▼`), którego rozwinięcie zawiera pełną listę akcji.
- **Znaczniki kontekstowe.** Środowisko, repozytorium, projekt, model i wykonawca występują jako lekkie znaczniki w pasku kontekstu, na przykład `[Danaco Console] [Ubuntu] [Worktree] [Fable 5] [Ultra]`; kliknięcie znacznika otwiera odpowiedni selektor.

### 7.4. Zasada jednego kliknięcia

Każda ukryta funkcja jest osiągalna jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego. Zagnieżdżanie funkcji głęboko w hierarchii menu jest zabronione — ukrycie zmniejsza chaos wizualny, nie utrudnia dostępu.

### 7.5. Katalog typów elementów warstw 2–4

Poniższe typy są pełnoprawnymi typami elementów interfejsu Danaco Console i występują w tabelach elementów rozdziałów 1–4 na równi z przyciskami, polami i kartami.

**Tabela 22 — Typy elementów realizujących stopniowe ujawnianie funkcjonalności**

| Typ elementu | Co to jest | Do czego służy | Forma i waga | Stany | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|
| Element zbiorczy z rozwinięciem (`Operacje ▼`) | Jeden przycisk zwinięty zastępujący zestaw przycisków akcji | Grupowanie logiczne akcji jednorodnego przeznaczenia | Przycisk `.dn-btn--zarys` ze znacznikiem `▼`, waga pojedynczego przycisku | Zwinięty · rozwinięty · fokus | 3 | Klik elementu; wybór akcji zwija element samoczynnie |
| Menu kebab (⋮) | Menu rozwijane przypisane do pojedynczego elementu lub okna | Dostęp do wariantów operacji na konkretnym elemencie | Ikona `⋮` w narożniku elementu, waga ikony | Ukryte w spoczynku · widoczne na hover/fokus · otwarte | 3 | Klik ikony; klik poza menu zamyka je |
| Menu hamburger (☰) | Zwinięty wyzwalacz nawigacji | Otwarcie kolumny nawigacji, gdy szerokość okna jej nie mieści | Ikona `menu` w pasku górnym | Zwinięte · rozwinięte | 3 | Klik ikony; wybór pozycji zwija menu |
| Panel popover | Mały panel zakotwiczony przy elemencie wywołującym | Prezentacja szczegółów i ustawień szybkich bez zmiany układu kolumn | Panel unoszony, szerokość dopasowana do treści | Zamknięty · otwarty · fokus wewnątrz | 3 | Klik elementu kotwiczącego; klik poza panelem zamyka go |
| Panel wysuwany | Kolumna boczna otwierana z prawej krawędzi obszaru roboczego | Pomieszczenie funkcji potrzebnych chwilowo | Kolumna boczna o regulowanej szerokości, pełna wysokość obszaru roboczego | Zamknięty (nieobecny) · otwarty | 3 | Klik ikony lub znacznika; zamknięcie usuwa kolumnę z przestrzeni roboczej |
| Znacznik kontekstowy | Lekka plakietka opisująca jeden wymiar kontekstu pracy | Prezentacja i zmiana środowiska, repozytorium, projektu, modelu, wykonawcy, trybu pracy, poziomu wysiłku | Plakietka `.dn-plakietka` w pasku kontekstu | Spoczynek · aktywny · selektor otwarty | 2 | Klik znacznika otwiera selektor; po wyborze selektor zwija się samoczynnie |
| Wyszukiwarka funkcji | Pole wyszukiwania funkcji platformy | Dotarcie do dowolnej funkcji, w tym funkcji warstwy 4, bez przeszukiwania menu | Panel popover wyśrodkowany nad obszarem roboczym | Zamknięta · otwarta · z wynikami · bez wyników | 4 | Skrót klawiszowy; wybór trafienia uruchamia funkcję i zamyka panel |

### 7.6. Warstwy a makiety

Makiety ASCII w rozdziałach 1–4 przedstawiają interfejs w **stanie spoczynku**: widoczne są wyłącznie elementy warstwy 1 oraz zwinięte wyzwalacze warstw 2–3 (znaczniki kontekstowe, `▼`, `⋮`, `☰`). Elementy warstw 2–4 opisane są w tabelach elementów okien wraz z podaniem warstwy i sposobu wywołania.

---

## Załącznik A. Katalog ikon wykorzystanych w oknach tego dokumentu

**Tabela 23 — Ikony (System wizualny, Załącznik C, przefiltrowane do okien tego dokumentu)**

| Ikona | Zastosowanie w oknach tego dokumentu |
|---|---|
| `dom` | Powrót do strony głównej / Centrum dowodzenia — pasek górny (3.2, 4.3) |
| `ustawienia` | Otwarcie okna konfiguracji z listwy strony głównej (2.5); ikona szybkiej konfiguracji sesji w pasku górnym (3.2) |
| `szukaj` | Pole wyszukiwania w pasku górnym (3.2) |
| `plus` | Nowa karta sesji (3.3, 4.3) |
| `zamknij` | Zamknięcie karty sesji (3.3) |
| `wyslij` | Wysłanie wiadomości w Chat Window (3.5.1) |
| `spinacz` | Załącznik w Chat Window (3.5.1) |
| `odpowiedz` | Cytowanie wiadomości w Chat Window (3.5.1) |
| `oko` | Pokaż/ukryj hasło (1.5); sekcja „Monitor procesu” panelu orkiestracji (4.4) |
| `dzwonek` | Always On Display — pasek górny (3.2, 4.3), listwa strony głównej (2.5) |
| `uzytkownik` | Sekcja „Role” panelu orkiestracji (4.4) |
| `gwiazdka` | Sekcja „Zespoły” panelu orkiestracji (4.4) |
| `filtr` | Sekcja „Kolejki” panelu orkiestracji (4.4) |
| `link-zewnetrzny` | Sekcja „Orkiestracja” panelu orkiestracji (4.4) |
| `zegar` / `kalendarz` | Sekcja „Harmonogram i automatyki” panelu orkiestracji (4.4) |
| `slonce` / `ksiezyc` | Przełącznik motywu jasny/ciemny — pasek górny (3.2) |
| `blad` / `ostrzezenie` / `info` | Komunikaty stanu formularza logowania/rejestracji (1.5) |
| `⋮` (kebab) | Menu kebab wiadomości Chat Window (3.5.1) i menu kebab powłoki w pasku górnym (3.2, 3.5.3) |
| `menu` | Menu hamburger — otwarcie bocznej nawigacji / panelu orkiestracji na wąskich ekranach (3.4, 4.4, 7.3) |

Pełna galeria z podglądem SVG (47 ikon): `ikony/indeks.html`.

---

## Załącznik B. Mapa źródeł

**Tabela 24 — Dokument źródłowy → rozdziały wykorzystane w niniejszym opracowaniu**

| Dokument źródłowy | Rozdziały wykorzystane |
|---|---|
| architektura/koncepcja-platformy.md | 1–14 (w szczególności rozdz. 4–8, 10–12), Załącznik A |
| interfejs-uzytkownika/strona-glowna-i-nawigacja.md | 2–9 (w całości), Załączniki A–C |
| specyfikacje/specyfikacja-okien-operacyjnych.md | 1.2, 2.1–2.4, 3, 5 |
| interfejs-uzytkownika/system-wizualny.md | 6–9, 11, 13, 14, Załączniki A–C |
| architektura/model-konfiguracji.md | 3.1–3.3, 4.1–4.3 |
| architektura/architektura.md | 1–9 |
| architektura/bezpieczenstwo-i-uwierzytelnianie.md | 1–15 (w całości), Załączniki A–D |
| `srodowiska/multitaskingai.md` | 6 (Panel orkiestracji — boczna nawigacja) |
| architektura/izolacja-i-zaleznosci.md | 6 (poziomy zasięgu — kontekst dla karty sesji i roli) |
| tokens.css | Tokeny paska górnego (`--dn-topbar-h`, `--dn-pasek-tlo`) |
| components.css | `.dn-pasek`, `.dn-pasek-marka`, `.dn-pasek-szukaj`, `.dn-zakladki--pigulki`, `.dn-awatar` |

---

*Koniec dokumentu. Danaco Console — Elementy okien przepływu głównego, wersja 2.0.*

---
*Danaco Console — AI Workspace OS · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — [LICENSE](LICENSE). Kontakt: support@danaco-group.pl*
