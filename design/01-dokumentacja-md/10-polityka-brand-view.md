# Danaco Console — Polityka Brand View

| | |
|---|---|
| **Produkt** | Danaco Console · AI Operating Environment (warstwa wizualna v2.0) |
| **Dokument** | Główna polityka Brand View — dokument nadrzędny |
| **Numer** | 10 · Marka i identyfikacja |
| **Wersja** | v2.0 |
| **Status** | **Wiążący** — Deweloperski (obowiązuje od chwili wydania) |
| **Data wydania** | 2026-08-14 |
| **Producent** | Danaco Holding Group Sp. z o.o. |
| **Twórca** | Dariusz Naharnowicz |
| **Kontakt** | support@danaco-group.pl |
| **Odbiorcy** | projektanci interfejsu · programiści warstwy widoku · redaktorzy treści interfejsu · osoba odbierająca pracę |
| **Zakres** | wszystkie widoki platformy: okno startowe, okno rejestracji i logowania, Centrum dowodzenia, powłoki czterech środowisk, okno Konfiguracji, okno Ustawień, Always On Display, Mobile, Chat Window, okna operacyjne piętnastu modułów, okna robocze ról MultitaskingAI, nakładki (modal, toast, tooltip, Centrum poleceń) |
| **Poza zakresem** | materiały marketingowe, strona internetowa producenta, komunikacja handlowa, dokumentacja drukowana |
| **Nadrzędność** | kontrakt systemu projektowego · kierunek systemu projektowego · zasada zero blokad |
| **Podrzędne wobec tego dokumentu** | opracowania 01 (Design System), 02 (Design Architecture), 03 (Design View), 09 (System marki) w zakresie sposobu **objawiania się marki w widoku** |

---

> **Deklaracja mocy wiążącej.** Niniejszy dokument nie jest poradnikiem, zbiorem
> sugestii ani inspiracją. Jest **polityką**. Widok, który narusza którąkolwiek
> z pięciu reguł konstytutywnych rozdziału 3, **nie zostaje odebrany** — niezależnie
> od tego, jak dobrze działa i jak dobrze wygląda. Odstępstwo jest możliwe wyłącznie
> w trybie rozdziału 13 i musi zostać zapisane.

---

## Spis treści

1. [Cel, zakres i status polityki](#1-cel-zakres-i-status-polityki)
2. [Zasada nadrzędna brand view](#2-zasada-nadrzędna-brand-view)
3. [Pięć reguł konstytutywnych brand view](#3-pięć-reguł-konstytutywnych-brand-view)
4. [Obecność znaku per widok](#4-obecność-znaku-per-widok)
5. [Hierarchia obecności marki](#5-hierarchia-obecności-marki)
6. [Polityka barwy w widoku](#6-polityka-barwy-w-widoku)
7. [Polityka typografii w widoku](#7-polityka-typografii-w-widoku)
8. [Polityka ikonografii w widoku](#8-polityka-ikonografii-w-widoku)
9. [Polityka ruchu w widoku](#9-polityka-ruchu-w-widoku)
10. [Polityka języka interfejsu](#10-polityka-języka-interfejsu)
11. [Polityka treści przykładowych](#11-polityka-treści-przykładowych)
12. [Polityka dostępności jako element marki](#12-polityka-dostępności-jako-element-marki)
13. [Procedura odstępstwa](#13-procedura-odstępstwa)
14. [Procedura odbioru widoku](#14-procedura-odbioru-widoku)
15. [Katalog naruszeń](#15-katalog-naruszeń)
16. [Decyzje projektowe](#16-decyzje-projektowe)
- [Załącznik A — karta odbioru widoku (do wydruku)](#załącznik-a--karta-odbioru-widoku-do-wydruku)
- [Załącznik B — mapa źródeł polityki](#załącznik-b--mapa-źródeł-polityki)
- [Załącznik C — słownik terminów obowiązujących](#załącznik-c--słownik-terminów-obowiązujących)

---

## 1. Cel, zakres i status polityki

### 1.1. Czym jest brand view

**Brand view** to sposób, w jaki tożsamość marki Danaco Console jest **widoczna
i egzekwowalna w każdym pojedynczym widoku produktu** — od okna startowego, przez
Centrum dowodzenia, po okno operacyjne pojedynczego modułu.

Pojęcie obejmuje pięć warstw, uporządkowanych od najbardziej do najmniej trwałej:

```
┌───────────────────────────────────────────────────────────────────┐
│  BRAND VIEW — pięć warstw obecności marki w widoku                │
├───────────────────────────────────────────────────────────────────┤
│  1 · ZNAK        godło, jego forma, rozmiar, pozycja, wariant     │
│  2 · RAMA        atramentowy pasek górny — stały dom godła        │
│  3 · JĘZYK       barwa, typografia, ikonografia, gęstość, ruch    │
│  4 · SŁOWO       nazewnictwo, ton, komunikaty, terminologia       │
│  5 · POSTAWA     zero blokad, dostępność, uczciwość treści        │
└───────────────────────────────────────────────────────────────────┘
        ▲ najtrwalsze — zmiana wymaga decyzji Właściciela
        ▼ najbardziej operacyjne — zmiana w trybie rozdziału 13
```

Marka nie jest w tym produkcie warstwą naklejoną na interfejs. Marka **jest**
interfejsem: rama, którą Operator widzi przez cały czas pracy, jest tym samym
elementem, który niesie godło.

### 1.2. Cel polityki

| Cel | Skutek praktyczny |
|---|---|
| Ujednolicić objawianie się marki | Widok z modułu Terminal i widok z modułu Research są rozpoznawalne jako jeden produkt bez czytania napisów |
| Uczynić markę **sprawdzalną** | Zgodność z marką jest werdyktem z listy kontrolnej, nie opinią |
| Ochronić znak | Geometria sygnetu jest zamrożona; wolno budować zastosowania, nie wolno zmieniać krzywych |
| Zabezpieczyć powściągliwość | Budżet akcentu i budżet ruchu są liczbami, nie odczuciem |
| Uczynić odstępstwo kosztownym, ale możliwym | Tryb rozdziału 13 zapisuje każde odstępstwo z terminem przeglądu |

### 1.3. Kogo obowiązuje

| Rola | Zakres obowiązku |
|---|---|
| **Projektant interfejsu** | Pełny. Projekt widoku niezgodny z polityką nie przechodzi do realizacji. |
| **Programista warstwy widoku** | Pełny w zakresie rozdziałów 3, 6–9, 12. Odpowiada za żetony, ikony, ruch, dostępność. |
| **Redaktor treści interfejsu** | Rozdziały 10, 11. Odpowiada za język, terminologię, treści przykładowe. |
| **Osoba odbierająca pracę** | Rozdziały 14, 15. Prowadzi listę kontrolną, wydaje werdykt, prowadzi rejestr odstępstw. |
| **Zespół systemu i architektury** | Rozdziały 6–9 jako warunek brzegowy przy definiowaniu żetonów i komponentów. |
| **Zespół marki** | Właściciel dokumentu; rozdziały 2–5 są jego wyłączną domeną redakcyjną. |
| **Zespół okien i prototypów** | Pełny — każdy prototyp okna jest widokiem w rozumieniu tej polityki. |
| **Wykonawca zewnętrzny** | Pełny, przekazywany razem z pakietem `zasoby/`. |

### 1.4. Jak stosuje się politykę przy odbiorze pracy

```
   PROJEKT WIDOKU
        │
        ▼
   [1] Samosprawdzenie wykonawcy — załącznik A, 39 pozycji
        │  wynik zapisany w opisie przekazania
        ▼
   [2] Odbiór formalny — rozdział 14, ta sama lista, osoba odbierająca
        │
        ├─► wszystkie pozycje spełnione ─────────► ZGODNY  ── widok przyjęty
        │
        ├─► naruszenia wyłącznie drobne ────────► DO POPRAWY ── przyjęty warunkowo,
        │                                                       poprawka w bieżącej iteracji
        │
        └─► choć jedno naruszenie krytyczne ────► NIEZGODNY ── zwrot do wykonawcy
            albo ≥ 3 naruszenia poważne
```

Werdykt jest **binarny w skutkach**: „zgodny" i „do poprawy" pozwalają widok
zaimplementować, „niezgodny" — nie. Osoba odbierająca nie negocjuje reguł;
negocjacja odbywa się wyłącznie w trybie rozdziału 13.

### 1.5. Status i cykl życia dokumentu

| Zdarzenie | Skutek |
|---|---|
| Zmiana w kontrakcie systemu projektowego | Obowiązkowy przegląd rozdziałów 2–9 tej polityki |
| Zmiana geometrii znaku | Wyłącznie decyzją Właściciela; unieważnia rozdział 4 do czasu nowego wydania |
| Nowe okno w inwentarzu | Obowiązkowe uzupełnienie tabeli rozdziału 4 przed pierwszym projektem tego okna |
| Trzy odstępstwa tego samego typu | Sygnał do przeglądu reguły — reguła albo się zmienia, albo odstępstwa się cofa |
| Wydanie produkcyjne | Zmiana statusu z „Deweloperski" na „Produkcyjny" bez zmiany treści reguł |

---

## 2. Zasada nadrzędna brand view

> **Marka Danaco Console objawia się STRUKTURĄ i PRECYZJĄ, nie ozdobą.**
> **Jeden akcent. Jedna rama. Jeden ruch.**

### 2.1. Rozwinięcie

Kierunek projektowy określa produkt jako **instrument pomiarowy**: zwarta gęstość,
dane krojem mono, zero dekoracji bez funkcji. Instrument pomiarowy nie jest
rozpoznawalny dlatego, że ma logo na obudowie. Jest rozpoznawalny dlatego, że
**zawsze pokazuje to samo w ten sam sposób**: ta sama skala, ta sama igła, ta sama
rama, ten sam sposób sygnalizowania, że pomiar trwa.

Z tego wynika trójdzielna definicja marki w widoku:

| Filar | Nośnik w widoku | Co go łamie |
|---|---|---|
| **Jeden akcent** | błękit sygnałowy `--dn-sygnal-*`, użyty wyłącznie jako **wskazanie** | druga barwa marki; sygnał jako powierzchnia |
| **Jedna rama** | atramentowy pasek górny `--dn-rama`, stały w obu motywach | pasek przełączający się z motywem; drugi pasek |
| **Jeden ruch** | tętno kropki 2,4 s — jedyny ruch ciągły | animowane tło, pulsujące przyciski, karuzele |

### 2.2. Trzy zdania, które muszą być prawdziwe o każdym widoku

1. **Gdyby usunąć z widoku wszystkie barwy poza neutralnymi, widok nadal działa** —
   każdy stan ma ikonę albo etykietę, hierarchia wynika z masy i pozycji.
2. **Gdyby zatrzymać wszystkie animacje, widok nadal informuje** — ruch jest
   wzmocnieniem komunikatu o stanie systemu, nigdy jedynym jego nośnikiem.
3. **Gdyby zasłonić godło, widok nadal jest rozpoznawalny jako Danaco Console** —
   po ramie, po gęstości, po parze krojów, po sposobie sygnalizowania pracy.

Widok, o którym któreś z tych zdań jest nieprawdziwe, opiera markę na ozdobie.

### 2.3. Czego zasada nadrzędna zakazuje wprost

- Dekoracji bez funkcji: ozdobnych linii, tekstur, wzorów tła, ilustracji „dla klimatu".
- Powtarzania godła jako elementu dekoracyjnego w treści widoku.
- Barwy jako środka wyrazu estetycznego — barwa w tym produkcie **znaczy**.
- Ruchu, który nie niesie informacji o stanie systemu.
- Nadmiaru: dwóch pasków, dwóch akcentów, dwóch rodzin krojów w tej samej roli.

---

## 3. Pięć reguł konstytutywnych brand view

Pięć reguł poniższych jest **konstytutywnych**: nie wynikają z estetyki, lecz z kontraktu
kierunku projektowego i z zasady nadrzędnej platformy (zasada zero blokad). Naruszenie
którejkolwiek z nich jest naruszeniem krytycznym albo poważnym — katalog wag
w rozdziale 15.

---

### 3.1. Rama atramentowa

> **Pasek górny jest zawsze atramentowy — w obu motywach.
> To jedyny stały dom godła w całym produkcie.**

#### Co znaczy

- Pasek górny ma wysokość `--dn-wym-pasek` (48 px) i tło `--dn-rama`
  (`--dn-szary-925`, `#131313`) **niezależnie od `data-theme`**.
- Tekst na pasku: `--dn-rama-tekst` (podstawowy) i `--dn-rama-tekst-2` (drugorzędny).
- Obrys i najechanie: `--dn-rama-obrys`, `--dn-rama-hover`.
- Przyciski ikonowe na pasku używają modyfikatora `.dn-btn-ikona--na-ramie`.
- **Godło występuje na pasku i tylko tam** — w widokach, które mają pasek.
- Godło zachowuje barwy własne marki (grot `#F4F4F4`, kropka `#5C8CEC` — wariant
  na ciemnym tle) w **obu** motywach, bo dom godła jest w obu motywach ciemny.

Rama daje efekt stanowiska dowodzenia: przełączalna treść w środku, stały
horyzont dookoła. To jest pierwsza rzecz, po której produkt się rozpoznaje.

#### Jak się ją sprawdza

1. Otwórz widok w motywie ciemnym. Zapamiętaj jasność paska górnego.
2. Przełącz motyw na jasny przyciskiem `[data-przelacz-motyw]`.
3. Pasek górny **nie może zmienić jasności ani barwy**. Zmienia się wyłącznie
   obszar pod nim.
4. Sprawdź, że godło na pasku jest wersją na ciemnym tle w obu motywach.
5. Sprawdź, że w treści widoku nie ma drugiego wystąpienia godła (poza widokami
   rozdziału 4 oznaczonymi jako „głośne").

#### Co jest naruszeniem

| Naruszenie | Waga |
|---|---|
| Pasek górny przyjmuje jasne tło w motywie jasnym | **krytyczne** |
| Godło poza paskiem w widoku, który nie jest „głośny" wg rozdziału 4 | poważne |
| Drugi pasek nawigacyjny o tej samej randze co górny | poważne |
| Kontrolki na pasku bez modyfikatora `--na-ramie` (kontrast tekstu) | poważne |
| Godło w wersji na jasnym tle użyte na pasku | poważne |

#### Przykład zgodny

```
┌──────────────────────────────────────────────────────────────┐  ← --dn-rama
│ »». DANACO Console │ TALKIN         [szukaj]  [ksiezyc]  ●   │     #131313
├──────────────────────────────────────────────────────────────┤     w OBU motywach
│                                                              │
│   obszar przełączalny z motywem: --dn-tlo / --dn-powierzchnia│
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

#### Przykład niezgodny

```
MOTYW JASNY — pasek przyjął barwę tła strony
┌──────────────────────────────────────────────────────────────┐  ← #F4F4F4
│ »». DANACO Console │ TALKIN         [szukaj]  [slonce]  ●    │     ZŁE:
├──────────────────────────────────────────────────────────────┤     rama zniknęła,
│                                                              │     godło straciło dom,
│   obszar treści — nieodróżnialny od paska                    │     produkt stracił
│                                                              │     rozpoznawalność
└──────────────────────────────────────────────────────────────┘
```

---

### 3.2. Jeden akcent

> **Sygnał zajmuje co najwyżej 5% powierzchni ekranu.
> Nigdy nie jest tłem sekcji. Nigdy nie jest tłem przycisku głównego.**

#### Co znaczy

- Jedyną barwą akcentu jest **błękit sygnałowy** — rodzina `--dn-sygnal-100…800`.
- Sygnał ma cztery role i tylko cztery: **fokus**, **stan aktywny/wybrany**,
  **odnośnik**, **wskaźnik pracy w tle** (kropka, postęp).
- **Działanie główne to inwersja atramentu**, nie sygnał: przycisk główny jest
  czarny w motywie jasnym (`.dn-btn--atrament`) i biały w motywie ciemnym.
  Modyfikator `.dn-btn--sygnal` jest zarezerwowany dla działań o charakterze
  „wskazania", nie dla głównego wezwania na widoku.
- Powierzchnia sygnału liczona łącznie: wypełnienia, tła stanowe rodziny
  informacyjnej, obrysy sygnałowe, kropki, tory postępu, pierścienie fokusu.
- **Budżet: ≤ 5% powierzchni widoku.** Budżety szczegółowe per widok — rozdział 6.4.
- Stany (zieleń, bursztyn, czerwień) są **wyjątkiem**, nie akcentem — pojawiają się
  tylko wtedy, gdy system ma coś do zakomunikowania, i znikają, gdy nie ma.

#### Jak się ją sprawdza

**Metoda A — pomiar zgrubny (obowiązkowa).**
1. Zrzut widoku w rozdzielczości progu `w3` (1280 px szerokości).
2. Nałóż siatkę 20 × 20 pól (400 pól, każde = 0,25% powierzchni).
3. Policz pola, w których występuje jakikolwiek piksel z rodziny sygnału.
4. Wynik: `liczba pól × 0,25%`. Próg: **≤ 5%** (≤ 20 pól).

**Metoda B — inwentaryzacja elementów (dla sporów).**
Wypisz każdy element sygnałowy z jego wymiarem, zsumuj pola, podziel przez
powierzchnię widoku. Kalkulator w wersji HTML tej polityki liczy to automatycznie.

**Metoda C — próba usunięcia (jakościowa).**
Wyłącz w widoku wszystkie barwy sygnału. Jeżeli widok stracił czytelność
hierarchii, sygnał niesie hierarchię zamiast jej wskazywać — to naruszenie reguły „jeden akcent"
niezależnie od wyniku pomiaru.

#### Co jest naruszeniem

| Naruszenie | Waga |
|---|---|
| Sygnał jako tło sekcji, panelu, nagłówka strony | **krytyczne** |
| Sygnał jako tło przycisku głównego widoku | **krytyczne** |
| Druga barwa akcentu poza rodziną sygnału i rodzinami stanów | poważne |
| Budżet akcentu przekroczony (5–8%) | poważne |
| Budżet akcentu przekroczony rażąco (> 8%) | **krytyczne** |
| Gradient sygnałowy jako tło karty, przycisku, sekcji | poważne |
| Stan komunikowany barwą stanu bez zdarzenia w systemie („zawsze zielone") | poważne |

#### Przykład zgodny

```
Centrum dowodzenia — sygnał = 2,25% (9 pól z 400)

┌──────────────────────────────────────────────────────────┐
│ RAMA: godło · kropka aktywności ●                        │  ← ● sygnał
├──────────────────────────────────────────────────────────┤
│  STREFA 1 — cztery karty środowisk                       │
│  ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐             │
│  │ TalkIn │ │WorkSpc │ │CodeStd │ │Multitsk│             │
│  │ emblem │ │ emblem │ │ emblem │ │ emblem │             │
│  │   ●    │ │   ●    │ │   ●    │ │   ●    │             │  ← 4 kropki emblematów
│  └────────┘ └────────┘ └────────┘ └────────┘             │
│  ═══ obrys karty pod kursorem: --dn-sygnal-obrys ═══     │  ← 1 obrys
│  STREFA 2 — cztery kafle (bez sygnału w stanie spoczynku)│
│  STREFA 3 — listwa (bez sygnału w stanie spoczynku)      │
│  Przycisk główny: ATRAMENT (czarny/biały), NIE sygnał    │
└──────────────────────────────────────────────────────────┘
```

#### Przykład niezgodny

```
Centrum dowodzenia — sygnał ≈ 34%

┌──────────────────────────────────────────────────────────┐
│ RAMA                                                     │
├──────────────────────────────────────────────────────────┤
│▓▓▓▓▓▓▓▓ STREFA 1 — całe tło sekcji w błękicie ▓▓▓▓▓▓▓▓▓▓▓│  ← ZŁE: tło sekcji
│▓ ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐           ▓│
│▓ │        │ │        │ │        │ │        │           ▓│
│▓ └────────┘ └────────┘ └────────┘ └────────┘           ▓│
│▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓│
│  [▓▓ ROZPOCZNIJ ▓▓]  ← ZŁE: przycisk główny w sygnale    │
└──────────────────────────────────────────────────────────┘
Skutek: sygnał przestał wskazywać. Kropka pracy w tle jest
niewidoczna, bo tonie w tle tej samej rodziny.
```

---

### 3.3. Jeden ruch ciągły

> **W widoku biegnie dokładnie jeden ruch ciągły: tętno kropki sygnału, 2,4 s.
> Każdy inny ruch jest przejściem, nie animacją.**

#### Co znaczy

- **Ruch ciągły** = powtarzalny bez końca, dopóki trwa stan. W produkcie jest
  jeden taki ruch: `--dn-czas-tetno` = 2,4 s, klasa `.dn-kropka--tetno`.
- Tętno oznacza **jedną rzecz**: „tu biegnie praca". Nie oznacza „to jest ważne",
  „tu kliknij", „nowość".
- **Przejście** = jednorazowa zmiana stanu, 100–220 ms, `--dn-ease`:
  `--dn-czas-1` (0,1 s) mikroreakcje · `--dn-czas-2` (0,16 s) barwy i motyw ·
  `--dn-czas-3` (0,22 s) warstwy (modal, panel, toast).
- Spinner `.dn-spinner` (14 px) jest dopuszczalny **wyłącznie** wewnątrz kontrolki
  wykonującej operację i wyłącznie na czas jej trwania — jest przejściem
  o nieznanym z góry czasie, nie animacją dekoracyjną.
- W jednym widoku może biec kilka **instancji** tętna (np. trzy karty sesji
  z pracą w tle), ale to nadal **jeden rodzaj ruchu**. Zakaz dotyczy drugiego
  *rodzaju* ruchu ciągłego.
- `prefers-reduced-motion` jest obsłużone **globalnie w żetonach**: tętno zamiera,
  zastępuje je pierścień statyczny. Nie dubluj tej obsługi w komponencie.

#### Jak się ją sprawdza

1. Otwórz widok, nie dotykaj myszy przez 10 sekund.
2. Wypisz wszystko, co się rusza. Dopuszczalna lista: **tętno kropki**, spinner
   trwającej operacji, pasek postępu z realnym postępem.
3. Wszystko poza tą listą jest naruszeniem.
4. Zmierz czasy przejść w narzędziach przeglądarki — każde przejście musi mieścić
   się w 100–220 ms i używać żetonu `--dn-czas-*`.
5. Włącz `prefers-reduced-motion` w systemie. Tętno musi zamrzeć, widok musi
   pozostać w pełni czytelny.

#### Co jest naruszeniem

| Naruszenie | Waga |
|---|---|
| Drugi rodzaj ruchu ciągłego (pulsujący przycisk, animowane tło, karuzela) | poważne |
| Tętno na elemencie, w którym nie biegnie praca | poważne |
| Przejście dłuższe niż 220 ms poza `--dn-czas-3` | drobne |
| Czas przejścia wpisany wprost zamiast żetonu `--dn-czas-*` | drobne |
| Parallaks, scrollytelling, ruch sterowany przewijaniem w oknie roboczym | poważne |
| Ruch, który jest jedynym nośnikiem informacji o stanie | **krytyczne** |

#### Przykład zgodny

```
Pas kart sesji — trzy sesje, dwie z pracą w tle

┌ Studio · analiza repozytorium ────── ● ┐  ← tętno 2,4 s (praca biegnie)
├ Research · przegląd źródeł ────────── ● ┤  ← tętno 2,4 s (praca biegnie)
├ Library · katalog zasobów ───────────── ┤  ← brak kropki (praca stoi)
└────────────────────────────────────────┘
Jeden rodzaj ruchu, dwie instancje, jedno znaczenie.
Najechanie na kartę: tło 0,1 s. Wybór karty: wstęga 0,16 s. To przejścia.
```

#### Przykład niezgodny

```
┌ Studio · analiza repozytorium ────── ● ┐  ← tętno (poprawne)
├ Research · przegląd źródeł ──────── ✦✦ ┤  ← ZŁE: migające gwiazdki „nowość"
├ Library · katalog zasobów ─────────────┤
└────────────────────────────────────────┘
  [ ROZPOCZNIJ ]  ← ZŁE: przycisk pulsuje, żeby zwrócić uwagę
  ▓▓▓▓ tło sekcji przesuwa gradient w pętli ▓▓▓▓  ← ZŁE: ruch dekoracyjny

Skutek: tętno kropki przestało znaczyć „tu biegnie praca",
bo rusza się wszystko. Sygnaturowy element marki został zużyty.
```

---

### 3.4. Stan nigdy samym kolorem

> **Każdy stan komunikowany barwą ma obok ikonę albo etykietę.
> Barwa jest wzmocnieniem, nigdy jedynym nośnikiem.**

#### Co znaczy

- Cztery rodziny stanów: **sukces** (zieleń), **ostrzeżenie** (bursztyn),
  **błąd** (czerwień), **informacja** (rodzina sygnału — celowe scalenie).
- Każde wystąpienie stanu musi nieść **co najmniej dwa** z trzech nośników:
  barwa · ikona z zestawu · etykieta słowna. Barwa sama nie wystarcza nigdy.
- Dotyczy: plakietek `.dn-plakietka--*`, kropek `.dn-kropka--*`, toastów
  `.dn-toast--*`, kroków kolejki `.dn-krok--*`, wierszy tabel, wskaźników
  w kartach sesji, pól z błędem walidacji.
- Dotyczy również **nadawców w oknie komunikacji**: trzy klasy semantyczne
  (`--czlowiek`, `--inteligencja`, `--system`) + ikona + etykieta + plakietka roli.
  Zakaz dziewięciu barw tła dla dziewięciu nadawców.
- Kropka sygnału jest wyjątkiem pozornym: sama kropka nie komunikuje stanu
  „sukces/błąd", tylko fakt „praca biegnie" — i zawsze towarzyszy jej nazwa
  sesji albo etykieta procesu, która ten fakt nazywa.

#### Jak się ją sprawdza

1. **Próba desaturacji.** Zrzut widoku → filtr `grayscale(1)`. Każdy stan musi
   pozostać rozpoznawalny.
2. **Próba czytnika.** Przejdź widok czytnikiem ekranu. Każdy stan musi być
   odczytany słowem, nie tylko obecnością elementu.
3. **Inwentarz stanów.** Wypisz każdy element barwny stanowy, przy każdym zaznacz,
   który drugi nośnik występuje (ikona / etykieta). Puste pole = naruszenie.

#### Co jest naruszeniem

| Naruszenie | Waga |
|---|---|
| Stan komunikowany wyłącznie barwą | **krytyczne** |
| Barwa stanu bez odpowiadającego jej zdarzenia w systemie | poważne |
| Więcej niż trzy klasy semantyczne nadawców w oknie komunikacji | poważne |
| Ikona stanu spoza zestawu (w tym emoji) | **krytyczne** |
| Etykieta stanu w żargonie zamiast w nazwie stanu | drobne |

#### Przykład zgodny

```
Kolejka MultitaskingAI — cztery kroki

  ✓  ZADANIE 001 · przygotowanie kontekstu      POPRAWNY   ● zieleń
  ⟳  ZADANIE 002 · wykonanie przez Executor 1   PRACUJE    ● sygnał (tętno)
  △  ZADANIE 003 · walidacja wyniku             OSTRZEŻENIE ● bursztyn
  ✕  ZADANIE 004 · scalenie odpowiedzi          BŁĘDY      ● czerwień

Każdy wiersz: ikona + etykieta wersalikowa + barwa. Po desaturacji
wiersze nadal czytelne — niesie je ikona i słowo.
```

#### Przykład niezgodny

```
Kolejka MultitaskingAI — cztery kroki

  ● ZADANIE 001 · przygotowanie kontekstu
  ● ZADANIE 002 · wykonanie przez Executor 1
  ● ZADANIE 003 · walidacja wyniku
  ● ZADANIE 004 · scalenie odpowiedzi

ZŁE: różnicuje wyłącznie barwa kropki. Po desaturacji cztery
identyczne wiersze. Operator z deuteranopią nie odróżni
zieleni od bursztynu. Czytnik ekranu odczyta cztery razy to samo.
```

---

### 3.5. Zero blokad

> **Brak atrybutu `disabled` w całym produkcie.
> Zamiast bramy — komunikat. Zamiast wyszarzenia — opis obok.**

#### Co znaczy

Zasada wywodzi się z zasady zero blokad: *Danaco Console nie narzuca twardych blokad, bram
bezpieczeństwa ani wymuszonych zgód. Domyślne zachowanie systemu to wykonanie
polecenia.* W warstwie widoku oznacza to:

- **Przycisk jest zawsze klikalny i zawsze fokusowalny.** Atrybut `disabled`
  nie występuje w kodzie widoku.
- Jeżeli działanie w danym momencie nie może się wykonać — po naciśnięciu pojawia
  się **komunikat mówiący dlaczego i co zrobić dalej** (toast, opis pod kontrolką,
  modal wyjaśniający).
- Alternatywnie: **opis obok kontrolki** wyjaśnia warunek, zanim Operator naciśnie.
- Przycisk „Pomiń" w oknie logowania jest obecny i klikalny **w obu fazach wdrożenia**.
- Odliczanie przy „Wyślij ponownie" jest **wyłącznie informacyjne** — nie blokuje
  kliknięcia.
- Metoda uwierzytelniania wyłączona w konfiguracji **nie jest renderowana wcale**:
  „brak metody = mniej segmentów, nie zablokowany segment".
- Wyszarzenie jako komunikat wizualny jest zakazane — również wtedy, gdy element
  formalnie pozostaje klikalny (wyszarzenie kłamie o dostępności).

#### Jak się ją sprawdza

1. Wyszukaj w kodzie widoku ciąg `disabled` — wynik musi być pusty.
   Sprawdź też `aria-disabled`, `pointer-events: none`, `[inert]`.
2. Przejdź widok wyłącznie klawiaturą (`Tab`). Każda kontrolka musi być osiągalna.
3. Naciśnij każdą kontrolkę, która w danym stanie „nie powinna działać".
   Musi pojawić się komunikat z **następnym krokiem**, nie sam brak reakcji.
4. Sprawdź, czy w widoku nie ma elementów o obniżonej kryciu, które udają blokadę.

#### Co jest naruszeniem

| Naruszenie | Waga |
|---|---|
| Atrybut `disabled` na jakiejkolwiek kontrolce | **krytyczne** |
| `pointer-events: none` / `[inert]` użyte jako brama | **krytyczne** |
| Wyszarzenie kontrolki jako komunikat o niedostępności | poważne |
| Naciśnięcie bez reakcji (cisza zamiast komunikatu) | poważne |
| Komunikat, który opisuje przeszkodę, ale nie podaje następnego kroku | poważne |
| Zablokowany segment metody uwierzytelniania zamiast jego pominięcia | poważne |

#### Przykład zgodny

```
Okno rejestracji — konto niepotwierdzone

  ┌────────────────────────────────────────────────┐
  │  Kod potwierdzający                            │
  │  [ ______ ]                                    │
  │                                                │
  │  [ Wyślij ponownie ]  Kolejna wysyłka za 0:42  │  ← informacja, nie blokada
  │                        ▲ opis obok, kontrolka  │
  │                          pozostaje klikalna    │
  │  [ Potwierdź ]   [ Pomiń ]                     │  ← „Pomiń" zawsze obecny
  └────────────────────────────────────────────────┘

Naciśnięcie „Wyślij ponownie" przed upływem 0:42:
  → toast: „Poprzednia wiadomość została wysłana 18 s temu.
            Sprawdź skrzynkę; kolejną można wysłać za 42 s."
```

#### Przykład niezgodny

```
  ┌────────────────────────────────────────────────┐
  │  Kod potwierdzający                            │
  │  [ ______ ]                                    │
  │                                                │
  │  [ Wyślij ponownie ]  ← wyszarzony, disabled   │  ← ZŁE
  │                                                │
  │  [ Potwierdź ]  ← wyszarzony do czasu wpisania │  ← ZŁE
  └────────────────────────────────────────────────┘

Skutek: Operator nie wie, czy przycisk jest zepsuty, czy czeka.
Czytnik ekranu pomija kontrolkę. Nawigacja klawiaturą przeskakuje.
Produkt zaczyna decydować za Operatora — wbrew zasadzie zero blokad.
```

---

### 3.6. Reguły w jednej tabeli

| Reguła | Zdanie | Miara | Główne naruszenie |
|---|---|---|---|
| **1** | Rama atramentowa | pasek `--dn-rama` w obu motywach | pasek zmienia się z motywem |
| **2** | Jeden akcent | sygnał ≤ 5% powierzchni | sygnał jako tło sekcji lub przycisku głównego |
| **3** | Jeden ruch ciągły | tętno 2,4 s, reszta ≤ 220 ms | drugi rodzaj ruchu ciągłego |
| **4** | Stan nigdy samym kolorem | ≥ 2 nośniki na stan | barwa jako jedyny nośnik |
| **5** | Zero blokad | zero `disabled` | brama zamiast komunikatu |

---

## 4. Obecność znaku per widok

### 4.1. Zasady ogólne

1. **Godło występuje raz na widok.** Wyjątek: widoki „głośne" (rozdział 5),
   w których znak jest treścią, a nie oznaczeniem.
2. **Poniżej 24 px obowiązuje wariant uproszczony** (jeden grot + kropka).
   Powyżej — wariant pełny (dwa groty + kropka).
3. **Pole ochronne** = wysokość jednego grotu (44 jednostki siatki 96, czyli
   ok. 0,46 × wysokość znaku) z każdej strony. W polu ochronnym nie ma nic.
4. **Godło zachowuje barwy własne marki** — nie przełącza się z motywem.
   Wybór wariantu (`sygnet.svg` / `sygnet-na-ciemnym.svg`) zależy od jasności
   podłoża, nie od `data-theme`.
5. **Emblemat środowiska nie jest godłem.** Jest ikoną domenową na siatce 24 × 24,
   obrys 1,75, `currentColor`, z dokładnie jedną wypełnioną kropką sygnału.
6. Znak wklejamy **inline jako `<svg>`**, nie przez `<img>` — żeby działał
   z `file://` i reagował na `currentColor` tam, gdzie to zamierzone.

### 4.2. Tabela obecności — przepływ główny

| Widok | Forma znaku | Rozmiar | Pozycja | Wariant | Uzasadnienie |
|---|---|---|---|---|---|
| **Okno startowe (ładowania)** | logo pionowy (sygnet + logotyp) | sygnet 96 px, blok 159 × 172 | oś pionowa, ok. 42% wysokości od góry | na ciemnym / na jasnym wg motywu strony | Jedyny moment, w którym produkt nie pracuje. Marka mówi pełnym głosem, bo nie zabiera miejsca pracy. |
| **Okno startowe — stan „Łączenie"** | kropka sygnału pod logo | 6 px, tętno 2,4 s | pod logotypem, odstęp `--dn-od-6` | kropka `--dn-kropka` | Tętno mówi „system żyje" w chwili, gdy nie ma jeszcze żadnej treści. |
| **Okno startowe — „Błąd połączenia"** | logo pionowy bez kropki | jw. | jw. | jw. | Kropka znika, bo praca nie biegnie. Zastępuje ją plakietka błędu z ikoną i etykietą (reguła „stan nigdy samym kolorem"). |
| **Okno rejestracji i logowania** | logo pionowy | sygnet 64 px, blok skalowany do 106 × 115 | nad panelem formularza, wyśrodkowany | na ciemnym / na jasnym wg motywu | Wejście do produktu — drugi i ostatni moment „głośny" przed pracą. |
| **Okno rejestracji — segmenty metod** | brak znaku | — | — | — | Segmenty niosą ikony metod, nie znak. Godło raz na widok. |
| **Strona główna — Centrum dowodzenia** | sygnet pełny + logotyp w ramie | sygnet 26 px | pasek górny, skrajnie z lewej | na ciemnym (rama) | Marka wraca do roli oznaczenia. Treść widoku należy do czterech środowisk. |
| **Centrum dowodzenia — Strefa 1** | emblematy środowisk | 40 px | godło karty środowiska | `currentColor`, kropka sygnału | Emblemat identyfikuje środowisko, nie producenta. |
| **Centrum dowodzenia — Strefa 2** | ikony komponentów własnych | 20 px | kafel, lewa krawędź | `currentColor` | Komponenty własne nie mają własnych emblematów — używają ikon zestawu. |

### 4.3. Tabela obecności — powłoki środowisk

| Widok | Forma znaku | Rozmiar | Pozycja | Wariant | Uzasadnienie |
|---|---|---|---|---|---|
| **Powłoka TalkIn** | sygnet + logotyp | 26 px | pasek górny, lewa | na ciemnym | Rama jest domem godła; obszar roboczy należy do pracy. |
| **Powłoka TalkIn — emblemat środowiska** | `srodowisko-talkin.svg` | 20 px | pasek górny, po separatorze, przy nazwie środowiska | `currentColor` = `--dn-rama-tekst-2` | Operator musi wiedzieć „w jakim trybie pracuję" bez czytania. |
| **Powłoka WorkSpace** | sygnet + logotyp | 26 px | pasek górny, lewa | na ciemnym | jw. |
| **Powłoka WorkSpace — emblemat** | `srodowisko-workspace.svg` | 20 px | jw. | `currentColor` | jw. |
| **Powłoka CodeStudio** | sygnet + logotyp | 26 px | pasek górny, lewa | na ciemnym | jw. |
| **Powłoka CodeStudio — emblemat** | `srodowisko-codestudio.svg` | 20 px | jw. | `currentColor` | jw. |
| **Powłoka MultitaskingAI** | sygnet + logotyp | 26 px | pasek górny, lewa | na ciemnym | jw. |
| **Powłoka MultitaskingAI — emblemat** | `srodowisko-multitaskingai.svg` | 20 px | jw. | `currentColor` | Emblemat węzłów opowiada orkiestrację; panel orkiestracji zastępuje boczną nawigację, więc emblemat jest jedynym oznaczeniem trybu. |
| **Boczna nawigacja modułów** | brak znaku | — | — | — | Boczna nawigacja niesie ikony modułów. Drugie godło byłoby powtórzeniem. |
| **Panel orkiestracji (6 sekcji)** | brak znaku | — | — | — | jw. |
| **Pas kart sesji** | kropka sygnału | 6 px | prawa krawędź karty sesji | `--dn-kropka`, tętno gdy praca biegnie | Element sygnaturowy w roli operacyjnej — nie oznaczenie, lecz informacja. |

### 4.4. Tabela obecności — okna platformowe

| Widok | Forma znaku | Rozmiar | Pozycja | Wariant | Uzasadnienie |
|---|---|---|---|---|---|
| **Okno Konfiguracji** | sygnet + logotyp | 26 px | pasek górny, lewa | na ciemnym | Okno platformowe zachowuje ramę powłoki. |
| **Okno Konfiguracji — Panel prowenancji** | brak znaku | — | — | — | Panel prezentuje pochodzenie danych; znak producenta w tym miejscu wprowadzałby w błąd co do źródła. |
| **Okno Konfiguracji — modal „Pula kont Code CLI"** | brak znaku | — | — | — | Modal nie jest widokiem samodzielnym — dziedziczy oznaczenie z powłoki pod nakładką. |
| **Okno Ustawień** | sygnet + logotyp | 26 px | pasek górny, lewa | na ciemnym | jw. |
| **Okno Ustawień — sekcja „Wygląd"** | podgląd motywu z godłem | sygnet 20 px | wewnątrz miniatury podglądu motywu | uproszczony (< 24 px) | Podgląd pokazuje ramę z godłem, bo rama jest tym, co się zmienia — a raczej: co się **nie** zmienia. |
| **Always On Display** | sygnet uproszczony | 16 px | listwa tożsamości nakładki, górna krawędź | uproszczony, na ciemnym | AOD leży nad ramą (`--dn-z-aod` 1200), więc traci dom godła. Uproszczony sygnet przywraca oznaczenie przy minimalnym koszcie powierzchni. |
| **Always On Display — rdzeń awatara** | kropka sygnału (rdzeń) | zależny od stanu | centrum warstwy | `--dn-grad-sygnal` dopuszczony ilustracyjnie | Rdzeń jest elementem sygnaturowym w największej skali w całym produkcie — patrz odstępstwo stałe dla Always On Display (rozdz. 13.5). |
| **Mobile** | sygnet uproszczony | 20 px | pasek górny mobilny, lewa | uproszczony, na ciemnym | Szerokość paska na progu `w1` nie mieści logotypu; poniżej 24 px obowiązuje wariant uproszczony. |
| **Mobile — ekran główny** | logotyp bez sygnetu | wysokość 14 px | pasek górny, wyśrodkowany | na ciemnym | Wyjątek dopuszczony: sygnet po lewej + logotyp na środku dublowałyby znak; na ekranie głównym prowadzi logotyp. |

### 4.5. Tabela obecności — okno wspólne i okna operacyjne

| Widok | Forma znaku | Rozmiar | Pozycja | Wariant | Uzasadnienie |
|---|---|---|---|---|---|
| **Chat Window (okno komunikacji)** | brak znaku | — | — | — | Okno wspólne wszystkim modułom; oznaczenie dziedziczy z ramy powłoki. |
| **Chat Window — medalion nadawcy „inteligencja"** | awatar `--inteligencja` | 28 px | lewa krawędź wpisu | `.dn-awatar--inteligencja` | Awatar inteligencji niesie kropkę sygnału jako rdzeń; nie jest godłem producenta. |
| **Chat Window — nadawca aktywnie piszący** | kropka sygnału | 6 px | przy nadawcy | tętno 2,4 s | jw. rozdz. 4.3, pas kart sesji. |
| **Okna operacyjne — panelowe** (Tools Panel, Sources Manager, Findings Panel, Instructions Panel, Context Memory, Model Configuration, Skills Manager, Connectors Manager, Permissions Center, Consensus Panel, Moderator Panel, Notes Panel, Glossary Manager, Assets Panel, Versioning Panel, Recommendations Panel) | brak znaku | — | — | — | Marka mówi cicho: rama, gęstość, para krojów, kropka pracy. Znak w panelu odbierałby miejsce danym. |
| **Okna operacyjne — wiodące** (Studio Editor, Research Workspace, Library Explorer, Translation Panels, Browser Window, Voice Console, Model Panels, Project Dashboard, Workflow Builder, Design Board, Product Builder, Terminal Tabs, Code Editor, Diagnostics Center, Agent Builder) | brak znaku | — | — | — | jw. Wiodące okno modułu jest miejscem pracy, nie ekspozycji marki. |
| **Terminal Tabs / Output Console** | brak znaku, brak akcentu poza fokusem | — | — | — | Terytorium milczenia (rozdz. 5.3): treść terminala rządzi się własną semantyką barw. |
| **Code Editor / Build Output** | brak znaku | — | — | — | jw. Podświetlenie składni nie podlega budżetowi akcentu (rozdz. 6.5). |
| **Podgląd dokumentu / File Preview / Preview Window** | brak znaku | — | — | — | jw. Treść dokumentu Operatora nie jest powierzchnią marki. |
| **Okna robocze ról MultitaskingAI** (Executor 1/2/3, Coordinator, Validator, Executor Chat, Coordinator Chat, Results Analyzer, Subagent Network) | brak znaku | — | — | — | Rola ma plakietkę roli i emblemat środowiska w ramie; znak byłby trzecim oznaczeniem. |
| **Modal `.dn-modal`** | brak znaku | — | — | — | Nakładka nad widokiem, nie widok samodzielny. |
| **Toast `.dn-toast`** | ikona rodziny stanu | 12 px | lewa krawędź | ikona z zestawu | Nośnik stanu wg reguły „stan nigdy samym kolorem", nie nośnik marki. |
| **Tooltip `.dn-tooltip`** | brak znaku | — | — | — | jw. |
| **Centrum poleceń** | sygnet uproszczony | 16 px | lewa krawędź pola polecenia | uproszczony, na ciemnym | Warstwa `--dn-z-centrum-polecen` (1300) leży nad wszystkim, w tym nad ramą — traci dom godła jak AOD. |
| **Pusty stan `.dn-pusty-stan`** | ikona z zestawu | 24 px | nad tytułem pustego stanu | `currentColor` = `--dn-tekst-3` | Pusty stan to komunikat, nie okazja do ekspozycji znaku. |

### 4.6. Tabela obecności — konteksty systemowe (poza oknem)

| Kontekst | Forma znaku | Rozmiar | Wariant | Uzasadnienie |
|---|---|---|---|---|
| **Favicon** | sygnet uproszczony | 16 / 32 / 48 px | `favicon.svg` z `prefers-color-scheme` | Przy 16 px dwa groty zlewają się w plamę; uproszczony zachowuje czytelność „» ." |
| **Ikona aplikacji** | sygnet pełny na kaflu | 180 / 192 / 512 / 1024 px | kafel `--dn-szary-925`, znak na ciemnym | Kafel odtwarza ramę kokpitu: znak zawsze na atramencie. |
| **Ikona maskowalna** | sygnet pełny, pole bezpieczne 80% | 512 px | jw. | System operacyjny może przyciąć kafel; znak mieści się w polu bezpiecznym. |
| **Karta udostępnienia (OG)** | logo poziomy | wysokość 96 px | na ciemnym | Poza zakresem tej polityki co do kompozycji; forma znaku wiążąca. |

### 4.7. Diagram: gdzie znak jest, a gdzie go nie ma

```
                          ZNAK GŁOŚNY            ZNAK CICHY           ZNAK MILCZY
                     (znak jest treścią)     (znak oznacza)      (znak nieobecny)
                              │                     │                     │
  Okno startowe ──────────────●                     │                     │
  Logowanie / rejestracja ────●                     │                     │
  Centrum dowodzenia ─────────┼─────────────────────●                     │
  Powłoka środowiska ─────────┼─────────────────────●                     │
  Okno Konfiguracji ──────────┼─────────────────────●                     │
  Okno Ustawień ──────────────┼─────────────────────●                     │
  Always On Display ──────────┼─────────────────────●                     │
  Mobile ─────────────────────┼─────────────────────●                     │
  Centrum poleceń ────────────┼─────────────────────●                     │
  Chat Window ────────────────┼─────────────────────┼─────────────────────●
  Okna operacyjne ────────────┼─────────────────────┼─────────────────────●
  Terminal / Code Editor ─────┼─────────────────────┼─────────────────────●
  Podgląd dokumentu ──────────┼─────────────────────┼─────────────────────●
  Modal / toast / tooltip ────┼─────────────────────┼─────────────────────●
                              │                     │                     │
     rozmiar znaku:        64–96 px              16–26 px               0 px
     udział marki:         wysoki                 średni                 zerowy
     pierwszeństwo:        marka                  równowaga              praca
```

---

## 5. Hierarchia obecności marki

Marka nie ma jednej głośności. Ma trzy — i przypisanie widoku do jednej z nich
jest decyzją wiążącą, nie preferencją.

### 5.1. Marka głośna — dwa widoki

**Gdzie:** okno startowe, okno rejestracji i logowania.

**Dlaczego tylko tam:** to jedyne dwa momenty, w których Operator **nie pracuje**.
Nic mu nie zabieramy, pokazując znak w pełnej skali. Po zalogowaniu każdy piksel
poświęcony marce jest pikselem odebranym pracy.

| Cecha | Wartość |
|---|---|
| Forma znaku | logo pionowy (sygnet + logotyp) |
| Rozmiar sygnetu | 64–96 px |
| Udział marki w kompozycji | do 25% powierzchni widoku (znak + pole ochronne) |
| Budżet akcentu | ≤ 2% (kropka, pierścień postępu, fokus) |
| Ruch | tętno kropki, gdy trwa łączenie; poza tym cisza |
| Typografia | Space Grotesk dopuszczony w nagłówku widoku |

**Ograniczenie:** nawet w widoku głośnym marka nie ozdabia. Okno startowe nie ma
ilustracji, gradientu tła ani animowanego tła. Ma znak, kropkę i komunikat.

### 5.2. Marka cicha — powłoki i okna platformowe

**Gdzie:** Centrum dowodzenia, powłoki czterech środowisk, okno Konfiguracji,
okno Ustawień, Always On Display, Mobile, Centrum poleceń.

**Dlaczego:** to widoki, w których Operator **orientuje się** — wybiera środowisko,
przełącza sesję, zmienia ustawienie. Marka musi być obecna jako punkt odniesienia
(„jestem w Danaco Console, w środowisku TalkIn"), ale nie może konkurować z treścią.

| Cecha | Wartość |
|---|---|
| Forma znaku | sygnet + logotyp w ramie |
| Rozmiar sygnetu | 16–26 px |
| Udział marki w kompozycji | ≤ 4% powierzchni (pasek 48 px na 800 px wysokości = 6% powierzchni, z czego godło zajmuje ułamek) |
| Budżet akcentu | ≤ 3–4% wg tabeli 6.4 |
| Ruch | tętno wyłącznie tam, gdzie biegnie praca |
| Typografia | Space Grotesk w tytułach środowisk i nagłówkach sekcji; Plex Sans w treści |

**Zasada praktyczna:** w widoku cichym marka jest **ramą wokół pracy**.
Operator ma ją widzieć peryferyjnie i przestać zauważać po pięciu minutach.

### 5.3. Marka milczy — okna operacyjne i terytoria treści

**Gdzie:** wszystkie okna operacyjne modułów, Chat Window, Terminal Tabs,
Output Console, Code Editor, Build Output, Preview Window, File Preview,
podgląd dokumentu, okna robocze ról MultitaskingAI, modale, toasty, tooltipy.

**Dlaczego:** to widoki, w których Operator **pracuje**. Praca ma bezwzględne
pierwszeństwo. Znak w oknie operacyjnym nie dodaje niczego — Operator wie,
gdzie jest, bo ma nad sobą ramę.

| Cecha | Wartość |
|---|---|
| Forma znaku | brak |
| Nośniki marki | rama nad widokiem, gęstość, para krojów, kropka sygnału, ikonografia |
| Budżet akcentu | ≤ 5% (pułap ogólny), ≤ 2% w terytoriach treści |
| Ruch | tętno wyłącznie przy realnie biegnącej pracy |
| Typografia | Plex Sans (interfejs) + Plex Mono (dane); **Space Grotesk zakazany** |

#### 5.3.1. Terytoria milczenia bezwzględnego

Trzy obszary, w których marka nie objawia się **wcale** — nawet przez akcent:

| Terytorium | Reguła |
|---|---|
| **Treść dokumentu** (Preview Window, File Preview, Report Builder — podgląd) | Treść Operatora jest renderowana jej własną typografią i barwą. Sygnał wyłącznie jako pierścień fokusu. |
| **Terminal** (Terminal Tabs, Output Console) | Barwy wyjścia procesu pochodzą z procesu. Podświetlenia produktu nie ma. Tło terminala `--dn-powierzchnia-2`, tekst `--dn-ff-mono`. |
| **Edytor kodu** (Code Editor, Studio Editor — obszar kodu, Diff/Grep Panel) | Podświetlenie składni jest osobnym systemem barw i **nie wlicza się do budżetu akcentu** (rozdz. 6.5). Chrom edytora podlega polityce; treść nie. |

**Uzasadnienie:** produkt, który koloruje cudzą treść własnymi barwami marki,
kłamie o źródle. Instrument pomiarowy pokazuje pomiar, nie własne logo na wykresie.

### 5.4. Tabela przypisania — wszystkie widoki

| Widok | Głośność | Znak | Budżet akcentu |
|---|---|---|---|
| Okno startowe | głośna | logo pionowy 96 px | ≤ 1% |
| Okno rejestracji i logowania | głośna | logo pionowy 64 px | ≤ 2% |
| Centrum dowodzenia | cicha | sygnet 26 px w ramie | ≤ 3% |
| Powłoka TalkIn | cicha | sygnet 26 px + emblemat 20 px | ≤ 4% |
| Powłoka WorkSpace | cicha | sygnet 26 px + emblemat 20 px | ≤ 4% |
| Powłoka CodeStudio | cicha | sygnet 26 px + emblemat 20 px | ≤ 4% |
| Powłoka MultitaskingAI | cicha | sygnet 26 px + emblemat 20 px | ≤ 4% |
| Okno Konfiguracji | cicha | sygnet 26 px | ≤ 4% |
| Okno Ustawień | cicha | sygnet 26 px | ≤ 3% |
| Always On Display | cicha | sygnet uproszczony 16 px | ≤ 8% (odstępstwo stałe) |
| Mobile | cicha | sygnet uproszczony 20 px | ≤ 4% |
| Centrum poleceń | cicha | sygnet uproszczony 16 px | ≤ 3% |
| Chat Window | milczy | — | ≤ 3% |
| Okna operacyjne — wiodące | milczy | — | ≤ 5% |
| Okna operacyjne — panelowe | milczy | — | ≤ 4% |
| Okna robocze ról MultitaskingAI | milczy | — | ≤ 5% |
| Terminal Tabs / Output Console | milczy bezwzględnie | — | ≤ 2% |
| Code Editor / Studio Editor (kod) | milczy bezwzględnie | — | ≤ 2% + składnia poza budżetem |
| Preview Window / File Preview | milczy bezwzględnie | — | ≤ 2% |
| Modal / toast / tooltip | milczy | — | ≤ 5% powierzchni nakładki |

---

## 6. Polityka barwy w widoku

### 6.1. Trzy warstwy barwy

```
┌─────────────────────────────────────────────────────────────────┐
│ WARSTWA 1 · MONOCHROM — podłoże                                 │
│ --dn-tlo · --dn-powierzchnia · --dn-powierzchnia-2 · --dn-panel │
│ --dn-obrys · --dn-tekst · --dn-tekst-2 · --dn-tekst-3           │
│ Zajmuje ≥ 90% powierzchni widoku. Neutralne czysto neutralne.   │
├─────────────────────────────────────────────────────────────────┤
│ WARSTWA 2 · SYGNAŁ — wskazanie                                  │
│ --dn-sygnal · --dn-sygnal-wypelnienie · --dn-kropka · --dn-fokus│
│ ≤ 5% powierzchni. Cztery role: fokus, aktywność, odnośnik, praca│
├─────────────────────────────────────────────────────────────────┤
│ WARSTWA 3 · STANY — wyjątek                                     │
│ sukces · ostrzeżenie · błąd · informacja (= rodzina sygnału)     │
│ Pojawia się tylko wtedy, gdy system ma coś do zakomunikowania.   │
│ Zawsze z ikoną albo etykietą (reguła „stan nigdy samym kolorem"). Znika, gdy stan mija.        │
└─────────────────────────────────────────────────────────────────┘
```

### 6.2. Reguły bezwzględne barwy

1. **Zero wartości szesnastkowych wprost.** Każda barwa przez `var(--dn-*)`.
   Jedyny wyjątek: pliki SVG znaku, gdzie barwy są własnością marki.
2. **Komponent sięga po żetony semantyczne**, nigdy po prymitywy.
   `--dn-tekst`, nie `--dn-szary-900`.
3. **`#000000` jest zakazane** — jako tło, jako tekst, jako obrys, wszędzie.
   Skala kończy się na `--dn-szary-1000` (`#0A0A0A`).
4. **`#FFFFFF` wyłącznie jako powierzchnia karty w motywie jasnym i tekst na atramencie** —
   nigdy jako tło całej strony (tło jasne to `--dn-szary-50`).
5. **Działanie główne = inwersja atramentu**, nie sygnał.
6. **Gradient nigdy nie jest tłem przycisku, karty ani sekcji.** Zakres gradientów:
   awatary bez zdjęcia, rdzeń AOD, grafiki brandowe.
7. **Rozmycie wyłącznie w nakładce modala.** Glassmorfizm w panelach jest zakazany.
8. **Cień wyłącznie z kompletu `--dn-cien-1/2/3/lg`.** `--dn-cien-sygnal`
   zarezerwowany dla stanu fokusu i wskazania.

### 6.3. Cztery role sygnału — i nic poza nimi

| Rola | Żeton | Przykład wystąpienia | Zakaz |
|---|---|---|---|
| **Fokus** | `--dn-fokus`, `--dn-fokus-cien` | pierścień 2 px + odsunięcie 2 px | fokus w barwie stanu |
| **Stan aktywny / wybrany** | `--dn-sygnal-tlo`, `--dn-sygnal-obrys` | wybrana pozycja bocznej nawigacji, wybrany wiersz tabeli, wybrana zakładka | zaznaczenie wypełnieniem sygnałowym całego wiersza |
| **Odnośnik** | `--dn-sygnal` | odnośnik w treści, odnośnik w opisie pola | odnośnik w barwie stanu |
| **Wskaźnik pracy w tle** | `--dn-kropka`, `--dn-sygnal-wypelnienie` | kropka sesji, tor postępu, wstęga aktywności | kropka bez biegnącej pracy |

### 6.4. Budżet akcentu per widok

Budżet to **maksymalny udział powierzchni widoku** zajęty przez wszystkie elementy
rodziny sygnału łącznie (warstwa 2 + informacyjna część warstwy 3).

| Widok | Budżet | Uzasadnienie budżetu |
|---|---:|---|
| Okno startowe | **1%** | Jedna kropka, jeden pierścień. Nic więcej się nie dzieje. |
| Okno rejestracji i logowania | **2%** | Fokus pola + jeden odnośnik + ewentualny pasek postępu odzyskiwania konta. |
| Centrum dowodzenia | **3%** | Cztery kropki emblematów + obrys karty pod kursorem + fokus. |
| Powłoka środowiska (4 ×) | **4%** | Wybrana pozycja bocznej nawigacji + kropki kart sesji + fokus + wstęga aktywnej karty. |
| Okno Konfiguracji | **4%** | Wiele przełączników w stanie włączonym — każdy niesie wypełnienie sygnałowe. |
| Okno Ustawień | **3%** | Mniej przełączników niż Konfiguracja; brak paneli technicznych. |
| Chat Window | **3%** | Kropka nadawcy piszącego + fokus pola promptu + odnośniki w treści. |
| Okno operacyjne — wiodące | **5%** | Pułap ogólny. Okna wiodące mają najwięcej stanów jednocześnie. |
| Okno operacyjne — panelowe | **4%** | Panel to jedna lista i jedno zaznaczenie. |
| Okna robocze ról MultitaskingAI | **5%** | Kolejka z wieloma krokami „pracuje" jednocześnie. |
| Always On Display | **8%** | **Odstępstwo stałe dla Always On Display** — rdzeń awatara jest elementem sygnaturowym w największej skali. |
| Mobile | **4%** | Mniejsza powierzchnia, ta sama liczba wskazań — udział rośnie naturalnie. |
| Centrum poleceń | **3%** | Fokus pola + zaznaczenie wyniku. |
| Terminal / edytor kodu / podgląd dokumentu | **2%** | Terytorium milczenia. Wyłącznie fokus i zaznaczenie. |
| Modal / toast / tooltip | **5%** powierzchni nakładki | Nakładka jest mała; udział liczy się względem niej, nie względem ekranu. |

**Przekroczenie budżetu** o mniej niż 3 punkty procentowe = naruszenie poważne.
Przekroczenie o 3 punkty procentowe lub więcej = naruszenie krytyczne.

### 6.5. Co nie wlicza się do budżetu akcentu

| Element | Powód wyłączenia |
|---|---|
| Podświetlenie składni w edytorze kodu i terminalu | Osobny system barw; treść nie należy do marki (rozdz. 5.3.1) |
| Treść dokumentu Operatora w podglądzie | jw. |
| Obraz wgrany przez Operatora, miniatura pliku | jw. |
| Pierścień fokusu widoczny **w danej chwili** na jednym elemencie | Fokus jest przejściowy i pojedynczy; liczy się stan spoczynku widoku |
| Barwy stanów sukces/ostrzeżenie/błąd | Należą do warstwy 3 „wyjątek"; podlegają odrębnej regule 6.6 |

### 6.6. Budżet stanów

Stany są wyjątkiem, więc mają własny, ostrzejszy budżet:

| Reguła | Wartość |
|---|---|
| Łączna powierzchnia barw stanowych (zieleń + bursztyn + czerwień) | ≤ 3% widoku |
| Liczba jednocześnie widocznych rodzajów stanu | ≤ 3 z 4 |
| Stan „sukces" utrzymywany po zakończeniu zdarzenia | zakazany — sukces wygasa |
| Wiersz listy wypełniony barwą stanu na całej szerokości | dopuszczalny wyłącznie dla „błąd" |

### 6.7. Kontrast — warunek, nie cel

Każda para tekst/tło w widoku musi mieć pokrycie w `zasoby/zetony/kontrasty.json`
(33 zmierzone pary). Pary spoza tego kompletu wymagają pomiaru **przed** użyciem.

| Żeton | Minimalne zastosowanie |
|---|---|
| `--dn-tekst` | dowolny tekst |
| `--dn-tekst-2` | tekst drugorzędny ≥ 12 px |
| `--dn-tekst-3` | **wyłącznie** metadane i tekst ≥ 18,66 px półgruby |
| `--dn-sygnal` | odnośnik i tekst sygnałowy — nigdy jako tekst na `--dn-sygnal-tlo` w motywie jasnym bez pomiaru |

---

## 7. Polityka typografii w widoku

### 7.1. Trzy role, trzy kroje, zero mieszania

| Krój | Rola wyłączna | Żeton | Gdzie występuje |
|---|---|---|---|
| **Space Grotesk** 500–700 | **wejście do środowiska i tożsamość** | `--dn-ff-naglowek` | logotyp, tytuły kart środowisk, nagłówki widoków głośnych i cichych, tytuły sekcji Centrum dowodzenia, tytuły modali |
| **IBM Plex Sans** 400–700 | **praca** | `--dn-ff-bazowa` | cały interfejs, etykiety, opisy, treść wpisów rozmowy, komunikaty |
| **IBM Plex Mono** 400–600 | **maszyna** | `--dn-ff-mono` | dane, identyfikatory, liczby w kolumnach, ścieżki, nazwy modeli, terminal, etykiety wersalikowe techniczne |

> **Różnica krojów niesie znaczenie.** Space Grotesk = marka i wejście.
> Plex Sans = praca. Plex Mono = maszyna. Mieszanie ról niszczy tę informację.

### 7.2. Zakazy mieszania ról

| Zakaz | Waga | Powód |
|---|---|---|
| Space Grotesk w treści roboczej, w etykietach pól, w wierszach tabel | poważne | Krój tożsamości użyty do pracy przestaje znaczyć „wejście" |
| Space Grotesk w oknie operacyjnym (terytorium milczenia) | poważne | Marka ma tam milczeć |
| Plex Mono w prozie, w opisach, w komunikatach ciągłych | poważne | Mono w prozie czyta się wolniej i sygnalizuje „to są dane", gdy to nie są dane |
| Plex Sans w kolumnie liczb bez `tabular-nums` | drobne | Liczby nie ustawiają się w kolumnę |
| Czwarty krój w widoku | **krytyczne** | Trzy role, trzy kroje |
| Krój systemowy jako świadomy wybór | poważne | Zapas w `font-family` jest zabezpieczeniem, nie decyzją projektową |

### 7.3. Skala i jej stosowanie

| Żeton | Wartość | Zastosowanie wiążące |
|---|---:|---|
| `--dn-fs-xs` | 11 px | etykiety wersalikowe, nagłówki kolumn |
| `--dn-fs-sm` | 12 px | metadane, tekst pomocniczy |
| `--dn-fs-base` | **13 px** | **bazowy — tekst interfejsu** |
| `--dn-fs-md` | 14 px | wyróżniona treść, wpisy rozmowy |
| `--dn-fs-lg` | 16 px | nagłówki paneli |
| `--dn-fs-xl` | 20 px | nagłówki okien i modali |
| `--dn-fs-2xl` | 24 px | tytuły sekcji strony głównej |
| `--dn-fs-3xl` | 30 px | tytuły kart środowisk |
| `--dn-fs-display` | 40 px | największy stopień ekspozycyjny — **wyłącznie widoki głośne** |

**Reguła:** stopień `--dn-fs-display` nie występuje w widoku cichym ani milczącym.
Stopień powyżej `--dn-fs-xl` nie występuje w oknie operacyjnym.

### 7.4. Wersaliki, odstępy liter, liczby

| Zastosowanie | Reguła |
|---|---|
| Etykieta wersalikowa Plex Sans | `--dn-ls-wersaliki` = 0,08em, obowiązkowo |
| Etykieta wersalikowa Plex Mono | `--dn-ls-mono-wersaliki` = 0,14em, obowiązkowo |
| Nagłówek Space Grotesk ≥ 24 px | `--dn-ls-naglowek` = −0,01em |
| Liczby w kolumnie, identyfikatory, czasy | `--dn-ff-mono` + `tabular-nums`, obowiązkowo |
| Interlinia interfejsu | `--dn-lh-bazowy` = 1,45 |
| Interlinia treści ciągłej | `--dn-lh-luzny` = 1,6 |
| Interlinia nagłówków | `--dn-lh-ciasny` = 1,25 |

### 7.5. Diakrytyki

Test obowiązkowy dla każdego widoku: **`ŁÓDŹ ŻÓŁĆ GĘŚLĄ JAŹŃ`**.
Wszystkie trzy kroje mają pełny zakres `latin-ext`. Widok, w którym którykolwiek
diakrytyk renderuje się z zapasu systemowego, jest niezgodny.

---

## 8. Polityka ikonografii w widoku

### 8.1. Jeden zestaw

| Reguła | Wartość |
|---|---|
| Źródło | **Lucide** (licencja ISC) — `zasoby/ikony/svg/*.svg`, manifest `zasoby/ikony/manifest.json` |
| Siatka | 24 × 24 |
| Obrys | **1,75** — ujednolicony, bez wyjątków |
| Wypełnienie | `none` (wyjątek: kropka sygnału w emblemacie środowiska) |
| Barwa | `currentColor` |
| Zakończenia i łączenia | `round` |
| Nazewnictwo | polskie, opisowe (kompletna lista w kontrakcie systemu projektowego) |

**Ikon nie rysuje się ręcznie, jeżeli istnieją w bibliotece.** Dorysowuje się
wyłącznie brakujące pojęcia domenowe — emblematy czterech środowisk i rola
koordynatora — na tej samej siatce i tym samym obrysem.

### 8.2. Jeden obrys

Grubość obrysu 1,75 jest **elementem tożsamości marki**, nie ustawieniem
technicznym. To ta sama decyzja co „precyzja instrumentu": jedna kreska
o jednej grubości w całym produkcie.

```
ZGODNE — jeden obrys 1,75 w całym widoku
   ⌂     △     ◷     ▤     ⌗       wszystkie 1,75

NIEZGODNE — obrysy z różnych źródeł
   ⌂     △     ◷     ▤     ⌗
  1,0   2,4   1,75  3,2   1,5      widok „drga"; oko czyta różnicę
                                    grubości jako różnicę ważności
```

### 8.3. Jeden rozmiar per kontekst

| Kontekst | Rozmiar | Żeton |
|---|---:|---|
| Ikona w plakietce, w toaście | 12–14 px | `--dn-wym-ikona-sm` |
| Ikona w wierszu, w polu, w przycisku | 16 px | `--dn-wym-ikona` |
| Ikona w pozycji bocznej nawigacji, na pasku górnym | 20 px | `--dn-wym-ikona-lg` |
| Ikona w kaflu, w pustym stanie, emblemat środowiska w listwie | 24 px | `--dn-wym-ikona-xl` |
| Emblemat środowiska w karcie środowiska | 40 px | — |

**Reguła:** w jednym kontekście występuje **jeden** rozmiar. Boczna nawigacja
z ikonami 20 px nie ma pozycji z ikoną 16 px.

### 8.4. Zakazy ikonograficzne

| Zakaz | Waga |
|---|---|
| **Emoji jako ikona** | **krytyczne** |
| Ikona z innego zestawu (inna metryka, inny obrys) | poważne |
| Ikona wypełniona tam, gdzie zestaw jest obrysowy | poważne |
| Ikona przeskalowana niezgodnie ze skalą 14/16/20/24 | drobne |
| Ikona jako jedyny nośnik stanu bez etykiety i bez `aria-label` | **krytyczne** (reguła „stan nigdy samym kolorem" + dostępność) |
| Emblemat środowiska z więcej niż jedną wypełnioną kropką | poważne |
| Emblemat środowiska bez kropki sygnału | poważne |
| Ikona dekoracyjna bez `aria-hidden="true"` | drobne |

### 8.5. Emblematy środowisk — konstytucja

Cztery emblematy są **jedynymi znakami domenowymi** produktu poza godłem.
Ich konstytucja jest wiążąca:

| Cecha | Wartość |
|---|---|
| Siatka | 24 × 24 |
| Obrys | 1,75, `round` |
| Barwa | `currentColor` — emblemat przejmuje barwę kontekstu |
| Kropka sygnału | **dokładnie jedna**, `fill="currentColor"`, `stroke="none"` |
| TalkIn | ramka rozmowy z ogonkiem, dwie linie tekstu, kropka po prawej |
| WorkSpace | cztery kwadranty, kropka w kwadrancie prawym dolnym |
| CodeStudio | ramka okna, grot kodu, linia, kropka po prawej górnej |
| MultitaskingAI | trzy węzły wokół rdzenia, rdzeń wypełniony jako kropka |

Zmiana krzywych emblematu wymaga tego samego trybu co zmiana godła — decyzji
Właściciela, nie odstępstwa z rozdziału 13.

---

## 9. Polityka ruchu w widoku

### 9.1. Budżet ruchu: INTENSYWNOSC_RUCHU = 3/10

Pokrętło intensywności ruchu ma wartość **3 na 10** i jest to liczba wiążąca.
W praktyce oznacza:

```
  0 ─────────── 3 ─────────────────────────────────── 10
  │             │                                      │
  cisza      DANACO                              strona
  całkowita  CONSOLE                             marketingowa
             │
             ├─ mikroprzejścia 100–220 ms
             ├─ jeden ruch znaczący: tętno 2,4 s
             ├─ przejścia widoków: opacity + translateY(4px)
             └─ zero animacji dekoracyjnych
```

### 9.2. Co wolno animować

| Element | Czas | Żeton | Warunek |
|---|---|---|---|
| Tło przy najechaniu, wciśnięcie | 0,1 s | `--dn-czas-1` | zawsze |
| Barwa, obrys, przełączenie motywu | 0,16 s | `--dn-czas-2` | zawsze |
| Wejście/wyjście modala, panelu, toastu | 0,22 s | `--dn-czas-3` | zawsze |
| Przejście między widokami | ≤ 0,22 s | `--dn-czas-3` | `opacity` + `translateY(4px)` |
| Tętno kropki sygnału | 2,4 s, w pętli | `--dn-czas-tetno` | **wyłącznie** gdy biegnie praca |
| Tor postępu | zależny od postępu | — | wyłącznie przy realnym postępie |
| Spinner kontrolki | ciągły na czas operacji | — | wyłącznie wewnątrz kontrolki wykonującej operację |
| Obrót grota rozwinięcia | 0,16 s | `--dn-czas-2` | zawsze |
| Wstęga aktywnej karty sesji | 0,16 s | `--dn-czas-2` | zawsze |

### 9.3. Czego animować nie wolno

| Zakaz | Waga |
|---|---|
| Parallaks, scrollytelling, ruch sterowany przewijaniem | poważne |
| Animowane tło, przesuwający się gradient, tekstura w ruchu | poważne |
| Pulsujący przycisk, migający element „zwracający uwagę" | poważne |
| Karuzela, automatyczne przewijanie treści | poważne |
| Animacja wejścia elementów listy jeden po drugim (kaskada) | drobne |
| Przejście dłuższe niż 220 ms poza modalem | drobne |
| Animacja `width`/`height` zamiast `transform`/`opacity` | drobne |
| Ruch jako jedyny nośnik informacji o stanie | **krytyczne** |

### 9.4. Krzywa i wartości

- Jedna krzywa w całym produkcie: `--dn-ease: cubic-bezier(0.2, 0, 0, 1)`.
- Czas wpisany wprost (np. `transition: .3s`) zamiast żetonu = naruszenie drobne.
- Krzywa inna niż `--dn-ease` (np. `ease-in-out`) = naruszenie drobne.

### 9.5. Ograniczony ruch

`prefers-reduced-motion` jest obsłużone **globalnie w żetonach**. Komponent
nie dubluje tej obsługi. Przy włączonym ograniczeniu:

| Element | Zachowanie |
|---|---|
| Tętno kropki | zamiera; zastępuje je **pierścień statyczny** wokół kropki |
| Przejścia 100–220 ms | skracane do zera przez warstwę żetonów |
| Spinner | pozostaje (niesie informację o trwaniu operacji) |
| Tor postępu | pozostaje (niesie informację o postępie) |

**Sprawdzenie:** widok z włączonym `prefers-reduced-motion` musi nadal
komunikować wszystkie stany. Jeżeli po wyłączeniu ruchu znika informacja —
to naruszenie krytyczne reguły „jeden ruch ciągły".

---

## 10. Polityka języka interfejsu

### 10.1. Zasady ogólne

| Zasada | Reguła |
|---|---|
| **Język** | polski; `lang="pl"` w dokumencie |
| **Forma** | bezosobowa („Nie udało się połączyć") **albo** bezpośredni zwrot do Operatora („Wybierz środowisko"). Jedna forma na widok — bez mieszania. |
| **Nazwy własne** | dokładnie jak w dokumentacji, po angielsku: `Studio Editor`, `Workflow Builder`, `Chat Window`, `Coordinator`, `Subagent Network`. **Zakaz tłumaczenia i parafrazowania.** |
| **Identyfikatory techniczne** | po angielsku (klasy CSS, klucze konfiguracji, nazwy modeli) |
| **Opisy, etykiety, komunikaty** | po polsku |
| **Ton** | rzeczowy, bez emfazy, bez wykrzykników, bez humoru |
| **Długość** | etykieta ≤ 3 wyrazy · opis pola ≤ 1 zdanie · komunikat błędu ≤ 2 zdania |

### 10.2. Zakaz żargonu marketingowego

Interfejs instrumentu pomiarowego nie sprzedaje. Zakazane w warstwie widoku:

| Zakazane | Zamiast tego |
|---|---|
| „potężny", „rewolucyjny", „inteligentny", „magiczny", „bezproblemowy" | opis funkcji: „uruchamia zadanie w tle" |
| „Twoje AI", „Twój asystent" jako nazwa elementu | nazwa własna z dokumentacji: `Assistant`, `Agent Builder` |
| „Zaczynajmy!", „Świetnie!", „Ups!" | „Gotowe", „Nie udało się połączyć" |
| Wykrzykniki w komunikatach | kropka |
| Emoji w tekście interfejsu | brak |
| „Kliknij tutaj" | nazwa działania: „Otwórz Konfigurację" |
| Metafory z innej dziedziny („odpal", „wystrzel", „magia") | czasownik dosłowny: „uruchom", „wykonaj" |

### 10.3. Komunikaty błędu — konstrukcja obowiązkowa

Każdy komunikat błędu ma **dwie części**:

```
┌─────────────────────────────────────────────────────────────┐
│  [1] STAN — co się stało, w formie faktu                    │
│  [2] NASTĘPNY KROK — co Operator może zrobić teraz          │
└─────────────────────────────────────────────────────────────┘
```

| Zgodne | Niezgodne | Dlaczego |
|---|---|---|
| „Nie udało się połączyć z usługą modeli. Sprawdź konfigurację konta w oknie Ustawień." | „Wystąpił błąd." | brak stanu i brak kroku |
| „Kolejka zawiera 3 zadania w stanie błędu. Otwórz Execution Monitor, aby zobaczyć przyczyny." | „Ups! Coś poszło nie tak [emotikona]" | żargon, emoji, brak informacji |
| „Poprzednia wiadomość została wysłana 18 s temu. Kolejną można wysłać za 42 s." | (przycisk wyszarzony bez komunikatu) | narusza regułę „zero blokad" |
| „Repozytorium `danaco/console` nie jest podłączone. Dodaj je w Project Tree." | „Brak dostępu." | brak kroku, brzmi jak brama |

**Zakaz:** komunikat nie może opisywać przeszkody bez podania kroku.
Komunikat bez kroku jest bramą w przebraniu i narusza regułę „zero blokad".

### 10.4. Nazewnictwo elementów interfejsu

| Element | Nazwa obowiązująca | Nazwa zakazana |
|---|---|---|
| Użytkownik produktu | **Operator** | „użytkownik", „klient", „ty" jako nazwa roli |
| TalkIn / WorkSpace / CodeStudio / MultitaskingAI | **środowisko** | „tryb", „przestrzeń", „workspace" (kolizja z modułem `Workspace`) |
| Studio, Research, Library… | **moduł** | „sekcja", „zakładka", „aplikacja" |
| Studio Editor, Workflow Builder… | **okno operacyjne** | „widok", „ekran", „panel" (panel to podtyp) |
| Pozycja w pasie kart | **karta sesji** | „tab", „zakładka" (zakładka to `.dn-zakladka`) |
| Automations, Agents, Workspace, Assistant w Strefie 2 | **komponent własny** | „widżet", „skrót", „aplikacja" |
| Model w konfiguracji roli MultitaskingAI | **ekspert** | „model" (model to nazwa techniczna), „bot" |
| Executor / Coordinator / Validator | **rola** | „agent" (agent to moduł `Agents`), „węzeł" |
| Lista zadań silnika MultitaskingAI | **kolejka** | „lista", „stos", „pipeline" |

### 10.5. Wersaliki w interfejsie

| Zastosowanie | Reguła |
|---|---|
| Etykieta sekcji, nagłówek kolumny | wersaliki + `--dn-ls-wersaliki`, `--dn-fs-xs` |
| Etykieta stanu w kroku kolejki | wersaliki + `--dn-ls-mono-wersaliki`, Plex Mono |
| Nazwa środowiska w ramie | wersaliki + `--dn-ls-mono-wersaliki` |
| Treść zdaniowa | **nigdy wersaliki** |
| Nazwa własna okna (`Studio Editor`) | zapis jak w dokumentacji, bez zmiany wielkości liter |

---

## 11. Polityka treści przykładowych

### 11.1. Zasada

> **Każda treść przykładowa pochodzi z domeny produktu i jest oznaczona
> jako przykładowa.** Treść spoza domeny produktu nie występuje w żadnym widoku,
> prototypie ani opracowaniu.

### 11.2. Zakazy bezwzględne

| Zakazane | Waga |
|---|---|
| Lorem ipsum i jakikolwiek tekst zastępczy bez znaczenia | **krytyczne** |
| Zmyślone nazwiska osób (katalog anty-domyślnych, kontrakt systemu projektowego) | **krytyczne** |
| Zmyślone metryki („99,9% skuteczności", „3× szybciej") | **krytyczne** |
| Fikcyjne firmy i klienci | **krytyczne** |
| Zmyślone cytaty i opinie | **krytyczne** |
| Zdjęcia stockowe osób jako awatary | poważne |
| Dane wyglądające na prawdziwe, bez oznaczenia „przykładowe" | poważne |

### 11.3. Skąd bierze się treść przykładowa

Treść przykładowa budowana jest **ze świata produktu**, spójnie z modułem:

| Rodzaj danych | Źródło przykładu | Przykład zgodny |
|---|---|---|
| Nazwa sesji | moduł + wykonywane zadanie | `Studio · analiza repozytorium` |
| Identyfikator zadania w kolejce | schemat kolejki MultitaskingAI | `ZADANIE 002` |
| Ścieżka repozytorium | nazwa produktu / nazwa modułu | `danaco/console`, `danaco/pilot` |
| Nazwa modelu | nazwa techniczna z konfiguracji | pole konfiguracji, bez podawania konkretnego dostawcy jako rekomendacji |
| Nazwa automatyki | czasownik + rzeczownik z domeny | `Nocna synchronizacja biblioteki` |
| Nazwa projektu | moduł `Workspace` | `Projekt: dokumentacja wdrożeniowa` |
| Rola MultitaskingAI | nazwa własna z dokumentacji | `Executor 1`, `Coordinator`, `Validator` |
| Nazwa okna | nazwa własna z dokumentacji | `Diagnostics Center`, `Sources Manager` |
| Tożsamość nadawcy w rozmowie | rola, nie osoba | `Operator`, `Inteligencja`, `System`, plakietka roli narzędzia |

### 11.4. Jak oznacza się treść przykładową

| Kontekst | Sposób oznaczenia |
|---|---|
| Prototyp okna | plakietka `--informacja` „dane przykładowe" w nagłówku widoku **albo** zdanie w panelu „O tym opracowaniu" |
| Opracowanie HTML | sekcja `<details>` „O tym opracowaniu" wymienia, które dane są przykładowe |
| Opracowanie MD | zdanie w nagłówku rozdziału zawierającego przykłady |
| Zrzut w dokumentacji | podpis z adnotacją „dane przykładowe" |

### 11.5. Metryki i liczby

| Reguła |
|---|
| Liczba w interfejsie pochodzi z systemu albo jest jawnie przykładowa |
| Zakaz podawania liczb wydajnościowych produktu (czasy odpowiedzi, skuteczność, oszczędność) |
| Liczby przykładowe muszą być wiarygodne co do rzędu wielkości (kolejka na 12 zadań, nie na 12 000) |
| Liczby przykładowe nie powtarzają się identycznie w wielu widokach — sugerowałoby to dane systemowe |

---

## 12. Polityka dostępności jako element marki

### 12.1. Deklaracja

> **WCAG 2.1 AA nie jest w Danaco Console zgodnością formalną.
> Jest obietnicą marki: instrument pomiarowy, którego nie da się odczytać,
> nie jest instrumentem pomiarowym.**

Dostępność należy do warstwy 5 brand view („postawa"). Jest tym samym
zobowiązaniem co zero blokad: produkt nie decyduje za Operatora i nie wyklucza go.

### 12.2. Warunki wejściowe — nie do negocjacji

| Warunek | Reguła | Sprawdzenie |
|---|---|---|
| **Kontrast** | każda para tekst/tło zmierzona; komplet w `kontrasty.json` (33 pomiary) | pomiar narzędziem, nie na oko |
| **`--dn-tekst-3`** | wyłącznie metadane i tekst ≥ 18,66 px półgruby | przegląd wystąpień |
| **Stan nigdy samym kolorem** | reguła konstytutywna — zawsze ikona albo etykieta | próba desaturacji |
| **Fokus** | pierścień 2 px + odsunięcie 2 px, `--dn-fokus`, widoczny na każdej kontrolce | przejście `Tab` przez cały widok |
| **Klawiatura** | każda funkcja osiągalna bez myszy; kolejność `Tab` zgodna z układem wizualnym | przejście klawiaturą |
| **`prefers-reduced-motion`** | obsłużone globalnie w żetonach | przełączenie w systemie |
| **Cele dotykowe** | `pointer: coarse` → kontrolka 40 px, wiersz 44 px, check 20 px, przełącznik 44 × 24 | emulacja dotyku |
| **Role i `aria-*`** | semantyczne znaczniki; `aria-label` na kontrolkach ikonowych; `aria-live` na obszarach zmiennych | przegląd drzewa dostępności |
| **`lang="pl"`** | w każdym dokumencie | przegląd źródła |
| **Zero `disabled`** | zero blokad | wyszukanie w kodzie |

### 12.3. Dostępność a pięć reguł

| Reguła | Wkład w dostępność |
|---|---|
| Rama atramentowa | stały punkt orientacji; kontrast tekstu na ramie zmierzony raz i niezmienny |
| Jeden akcent | ≤ 5% sygnału oznacza, że fokus jest zawsze widoczny — nie tonie w błękitnym tle |
| Jeden ruch | brak ruchu dekoracyjnego = brak wyzwalaczy dla wrażliwości przedsionkowej |
| Stan nigdy kolorem | bezpośrednio realizuje kryterium 1.4.1 |
| Zero blokad | brak `disabled` = brak kontrolek pomijanych przez czytnik i klawiaturę |

**Wniosek:** przestrzeganie brand view realizuje znaczną część WCAG 2.1 AA
niejako przy okazji. To nie przypadek — reguły były tak formułowane.

### 12.4. Czego dostępność nie usprawiedliwia

| Argument | Odpowiedź |
|---|---|
| „Zwiększyłem kontrast, więc dodałem drugą barwę" | Kontrast osiąga się jasnością, nie liczbą barw. |
| „Dodałem animację, żeby zwrócić uwagę na błąd" | Uwagę zwraca ikona, etykieta i pozycja komunikatu. |
| „Wyszarzenie informuje, że nie można kliknąć" | Wyszarzenie obniża kontrast i narusza regułę „zero blokad". |
| „Powiększyłem tekst poza skalę, bo tak czytelniej" | Skala ma dziewięć stopni; czytelność zapewnia stopień bazowy 13 px i interlinia 1,45. |

---

## 13. Procedura odstępstwa

### 13.1. Zasada

Odstępstwo od polityki jest **możliwe, kosztowne i zapisane**. Nie ma odstępstw
milczących. Element niezgodny bez zapisanego odstępstwa jest naruszeniem —
niezależnie od tego, kto go zaprojektował.

### 13.2. Kto może odstąpić

| Zakres odstępstwa | Uprawniony | Forma |
|---|---|---|
| Geometria znaku, geometria emblematów, para krojów, barwa sygnału | **Właściciel** (Dariusz Naharnowicz) | decyzja, nie odstępstwo — powoduje nowe wydanie kontrakt systemu projektowego i tej polityki |
| Reguła konstytutywna | **Właściciel** | odstępstwo stałe, wpisane do rozdziału 13.5 |
| Budżet akcentu, budżet ruchu, budżet stanów | **Zespół marki** | odstępstwo terminowe, rejestr 13.4 |
| Obecność znaku w konkretnym widoku (rozdz. 4) | **Zespół marki** | odstępstwo terminowe |
| Rozmiar ikony, stopień typograficzny, żeton odstępu | **Osoba odbierająca pracę** | odstępstwo jednorazowe, adnotacja w karcie odbioru |
| Nazwa własna okna, modułu, roli | **nikt** | nazwy własne są cytatem z dokumentacji |
| Zero blokad | **nikt** | zasada nadrzędna platformy |
| Zakaz `#000000`, Lorem ipsum, emoji jako ikon | **nikt** | zakaz bezwzględny |

### 13.3. Tryb

```
  [1] WNIOSEK
      Kto: wykonawca widoku
      Co zawiera: widok · reguła · zakres odstępstwa · uzasadnienie
                  · rozważone alternatywy · proponowany termin przeglądu
        │
        ▼
  [2] OCENA
      Kto: uprawniony wg tabeli 13.2
      Kryteria: (a) czy alternatywa zgodna z polityką naprawdę nie istnieje
                (b) czy odstępstwo dotyczy jednego widoku, czy tworzy precedens
                (c) czy skutek dla rozpoznawalności marki jest odwracalny
        │
        ├─► ODMOWA ────► widok wraca do poprawy; powód zapisany
        │
        └─► ZGODA ─────► wpis do rejestru 13.4 z terminem przeglądu
                             │
                             ▼
  [3] PRZEGLĄD (obowiązkowy, w terminie z wpisu)
      ├─► odstępstwo przedłużone (nowy termin)
      ├─► odstępstwo cofnięte (widok poprawiany)
      └─► reguła zmieniona (odstępstwo staje się regułą — nowe wydanie polityki)
```

### 13.4. Rejestr odstępstw terminowych

Rejestr prowadzi osoba odbierająca pracę. Wpis zawiera obowiązkowo:

| Pole | Opis |
|---|---|
| Numer | `OT-nn` |
| Data | data zgody |
| Widok | pełna nazwa widoku wg inwentarza |
| Reguła / rozdział | np. `jeden akcent · budżet akcentu 6,5%` |
| Uzasadnienie | jedno zdanie, rzeczowo |
| Alternatywy rozważone | co sprawdzono i dlaczego nie działa |
| Uprawniony | kto wyraził zgodę |
| Termin przeglądu | data |
| Skutek przeglądu | przedłużone / cofnięte / reguła zmieniona |

**Stan na dzień wydania:** rejestr odstępstw terminowych jest pusty.

### 13.5. Odstępstwa stałe

| Nr | Zakres | Reguła | Treść | Uzasadnienie |
|---|---|---|---|---|
| **1** | Always On Display | jeden akcent · budżet akcentu | Budżet akcentu podniesiony z 5% do **8%** | Rdzeń awatara AOD jest elementem sygnaturowym marki w największej skali w produkcie. AOD nie jest widokiem roboczym — nie ma w nim treści, której sygnał mógłby przeszkodzić. Gradient sygnałowy dopuszczony ilustracyjnie zgodnie z kontraktem systemu projektowego (zakres gradientów: rdzeń AOD). |
| **2** | Mobile — ekran główny | Rozdz. 4 · jedno godło na widok | Logotyp bez sygnetu na środku paska, przy jednoczesnym braku sygnetu po lewej | Szerokość paska na progu `w1` nie mieści sygnetu i logotypu jednocześnie z kontrolkami. Wybór: logotyp niesie więcej informacji niż sam sygnet dla Operatora, który dopiero wchodzi do produktu na urządzeniu mobilnym. |
| **3** | Edytor kodu, terminal | jeden akcent · budżet akcentu | Podświetlenie składni poza budżetem akcentu | Barwy składni nie należą do marki (rozdz. 5.3.1). Wliczanie ich do budżetu wymuszałoby monochromatyczny edytor, co obniżyłoby użyteczność bez zysku dla marki. |

---

## 14. Procedura odbioru widoku

### 14.1. Kiedy przeprowadza się odbiór

| Moment | Zakres |
|---|---|
| Przekazanie projektu widoku do realizacji | pełna lista, 39 pozycji |
| Przekazanie zaimplementowanego widoku | pełna lista, 39 pozycji |
| Zmiana w widoku już odebranym | pozycje dotknięte zmianą + wszystkie pięć reguł konstytutywnych zawsze |
| Przegląd okresowy pakietu | próba losowa 5 widoków, pełna lista |

### 14.2. Tryb weryfikacji

| Krok | Czynność | Narzędzie |
|---|---|---|
| 1 | Otwarcie widoku w motywie ciemnym (domyślnym) | przeglądarka, `file://` |
| 2 | Przełączenie motywu `[data-przelacz-motyw]`, kontrola ramy | przeglądarka |
| 3 | Zrzut widoku w progu `w3` (1280 px) | zrzut ekranu |
| 4 | Pomiar budżetu akcentu metodą A (siatka 20 × 20) | kalkulator w wersji HTML polityki |
| 5 | Próba desaturacji — kontrola reguły „stan nigdy samym kolorem" | filtr `grayscale(1)` |
| 6 | Obserwacja 10 s bez interakcji — kontrola reguły „jeden ruch ciągły" | obserwacja |
| 7 | Przejście całego widoku klawiszem `Tab` — kontrola fokusu i reguły „zero blokad" | klawiatura |
| 8 | Wyszukanie `disabled`, `#000000`, wartości hex, `pointer-events: none` | wyszukiwanie w źródle |
| 9 | Kontrola treści: nazwy własne, brak żargonu, komunikaty z krokiem | odczyt |
| 10 | Włączenie `prefers-reduced-motion`, ponowna obserwacja | ustawienia systemu |
| 11 | Wypełnienie listy kontrolnej, werdykt, zapis | karta odbioru (załącznik A) |

### 14.3. Lista kontrolna brand view — 39 pozycji

#### Blok A · Reguły konstytutywne (8)

| # | Pozycja | Reguła | Waga naruszenia |
|---:|---|---|---|
| 1 | Pasek górny jest atramentowy w **obu** motywach | rama atramentowa | krytyczne |
| 2 | Godło występuje wyłącznie na pasku (albo zgodnie z rozdz. 4 dla widoków głośnych) | rama atramentowa | poważne |
| 3 | Godło zachowuje barwy własne marki i wariant na ciemnym tle | rama atramentowa | poważne |
| 4 | Sygnał zajmuje ≤ budżet z tabeli 6.4 | jeden akcent | poważne / krytyczne |
| 5 | Sygnał nie jest tłem sekcji ani tłem przycisku głównego | jeden akcent | krytyczne |
| 6 | W widoku biegnie co najwyżej jeden **rodzaj** ruchu ciągłego | jeden ruch ciągły | poważne |
| 7 | Każdy stan barwny ma ikonę albo etykietę | stan nigdy samym kolorem | krytyczne |
| 8 | W kodzie widoku nie ma `disabled`, `[inert]`, `pointer-events: none` jako bramy | zero blokad | krytyczne |

#### Blok B · Znak i emblematy (5)

| # | Pozycja | Rozdział | Waga |
|---:|---|---|---|
| 1 | Forma, rozmiar i pozycja znaku zgodne z tabelą rozdz. 4 | 4 | poważne |
| 2 | Poniżej 24 px użyto wariantu uproszczonego | 4.1 | poważne |
| 3 | Pole ochronne znaku zachowane (≈ 0,46 × wysokość znaku) | 4.1 | drobne |
| 4 | Geometria sygnetu nietknięta (groty i kropka co do współrzędnej) | 4.1 | krytyczne |
| 5 | Emblemat środowiska ma dokładnie jedną wypełnioną kropkę | 8.5 | poważne |

#### Blok C · Barwa (6)

| # | Pozycja | Rozdział | Waga |
|---:|---|---|---|
| 1 | Zero wartości szesnastkowych wprost (poza SVG znaku) | 6.2 | krytyczne |
| 2 | Zero `#000000` | 6.2 | krytyczne |
| 3 | Komponenty sięgają po żetony semantyczne, nie po prymitywy | 6.2 | poważne |
| 4 | Przycisk główny w inwersji atramentu, nie w sygnale | 6.2 | krytyczne |
| 5 | Gradient nie jest tłem przycisku, karty ani sekcji | 6.2 | poważne |
| 6 | Rozmycie wyłącznie w nakładce modala | 6.2 | poważne |

#### Blok D · Typografia (4)

| # | Pozycja | Rozdział | Waga |
|---:|---|---|---|
| 1 | Space Grotesk wyłącznie w roli tożsamości i wejścia | 7.2 | poważne |
| 2 | Plex Mono wyłącznie w danych, identyfikatorach, terminalu | 7.2 | poważne |
| 3 | Liczby w kolumnach z `tabular-nums` w Plex Mono | 7.4 | drobne |
| 4 | Diakrytyki `ŁÓDŹ ŻÓŁĆ GĘŚLĄ JAŹŃ` renderują się z kroju, nie z zapasu | 7.5 | poważne |

#### Blok E · Ikonografia (3)

| # | Pozycja | Rozdział | Waga |
|---:|---|---|---|
| 1 | Wszystkie ikony z jednego zestawu, obrys 1,75 | 8.1–8.2 | poważne |
| 2 | Jeden rozmiar ikony per kontekst (14/16/20/24) | 8.3 | drobne |
| 3 | Zero emoji jako ikon | 8.4 | krytyczne |

#### Blok F · Ruch (3)

| # | Pozycja | Rozdział | Waga |
|---:|---|---|---|
| 1 | Czasy przejść wyłącznie z żetonów `--dn-czas-*`, krzywa `--dn-ease` | 9.4 | drobne |
| 2 | Tętno wyłącznie tam, gdzie realnie biegnie praca | 9.2 | poważne |
| 3 | Przy `prefers-reduced-motion` widok nadal komunikuje wszystkie stany | 9.5 | krytyczne |

#### Blok G · Język i treść (5)

| # | Pozycja | Rozdział | Waga |
|---:|---|---|---|
| 1 | Nazwy własne okien, modułów, ról dokładnie jak w dokumentacji | 10.1 | poważne |
| 2 | Jedna forma zwracania się do Operatora w całym widoku | 10.1 | drobne |
| 3 | Zero żargonu marketingowego, zero wykrzykników, zero emoji w tekście | 10.2 | poważne |
| 4 | Każdy komunikat błędu zawiera stan **i** następny krok | 10.3 | poważne |
| 5 | Treści przykładowe z domeny produktu, oznaczone jako przykładowe; zero Lorem ipsum, zmyślonych nazwisk i metryk | 11 | krytyczne |

#### Blok H · Dostępność (5)

| # | Pozycja | Rozdział | Waga |
|---:|---|---|---|
| 1 | Fokus widoczny na każdej kontrolce (2 px + odsunięcie 2 px) | 12.2 | krytyczne |
| 2 | Cały widok osiągalny klawiaturą, kolejność `Tab` zgodna z układem | 12.2 | krytyczne |
| 3 | Pary tekst/tło mają pokrycie w `kontrasty.json` albo zmierzony pomiar | 12.2 | krytyczne |
| 4 | Kontrolki ikonowe mają `aria-label`; ikony dekoracyjne `aria-hidden` | 12.2 | poważne |
| 5 | `lang="pl"`, role semantyczne, `aria-live` na obszarach zmiennych | 12.2 | poważne |

**Razem: 39 pozycji w 8 blokach.**

### 14.4. Werdykt

| Werdykt | Warunek | Skutek |
|---|---|---|
| **ZGODNY** | 39 z 39 pozycji spełnionych | Widok przyjęty. Data i podpis w karcie odbioru. |
| **DO POPRAWY** | zero naruszeń krytycznych · ≤ 2 naruszenia poważne · dowolna liczba drobnych | Widok przyjęty warunkowo. Poprawki w bieżącej iteracji, ponowny odbiór wyłącznie pozycji naruszonych. |
| **NIEZGODNY** | ≥ 1 naruszenie krytyczne **albo** ≥ 3 naruszenia poważne | Zwrot do wykonawcy. Ponowny odbiór pełny, 39 pozycji. |

### 14.5. Zapis odbioru

Każdy odbiór zapisuje się w karcie (załącznik A) z polami: widok · data ·
osoba odbierająca · wynik per blok · lista naruszeń z wagami · werdykt ·
odstępstwa powołane · termin ponownego odbioru.

---

## 15. Katalog naruszeń

### 15.1. Wagi

| Waga | Definicja | Skutek |
|---|---|---|
| **Krytyczne** | narusza regułę konstytutywną w sposób nieodwracalny dla rozpoznawalności albo dostępności; łamie zakaz bezwzględny | werdykt NIEZGODNY, zwrot |
| **Poważne** | narusza politykę w sposób widoczny, ale odwracalny bez przeprojektowania widoku | trzy takie = NIEZGODNY |
| **Drobne** | odstępstwo od skali, żetonu albo konwencji, niewidoczne dla Operatora, widoczne w kodzie | DO POPRAWY |

### 15.2. Naruszenia krytyczne

| Nr | Naruszenie | Reguła | Naprawa |
|---|---|---|---|
| **1** | Pasek górny przełącza się z motywem | rama atramentowa | Ustaw tło paska na `--dn-rama`; usuń zależność od `data-theme` |
| **2** | Zmieniona geometria sygnetu (groty, kropka, proporcje, promień) | reguła „rama atramentowa" / rozdz. 4 | Przywróć plik z `zasoby/marka/logo/`; wklej inline bez modyfikacji ścieżek |
| **3** | Atrybut `disabled` na kontrolce | zero blokad | Usuń atrybut; dodaj komunikat po naciśnięciu albo opis obok |
| **4** | `pointer-events: none` albo `[inert]` użyte jako brama | zero blokad | jw. |
| **5** | Stan komunikowany wyłącznie barwą | stan nigdy samym kolorem | Dodaj ikonę z zestawu **i** etykietę słowną |
| **6** | Sygnał jako tło sekcji, panelu albo nagłówka | jeden akcent | Zamień na `--dn-powierzchnia` / `--dn-panel`; sygnał zostaw obrysowi i kropce |
| **7** | Sygnał jako tło przycisku głównego | jeden akcent | Zamień na `.dn-btn--atrament` |
| **8** | Budżet akcentu przekroczony o ≥ 3 punkty procentowe | jeden akcent | Zredukuj wypełnienia sygnałowe do obrysów i kropek |
| **9** | Ruch jako jedyny nośnik informacji o stanie | jeden ruch ciągły | Dodaj etykietę stanu; ruch zostaw jako wzmocnienie |
| **10** | Widok traci informację przy `prefers-reduced-motion` | reguła „jeden ruch ciągły" / rozdz. 12 | Przenieś informację do etykiety i ikony |
| **11** | `#000000` jako tło, tekst albo obrys | rozdz. 6.2 | Zamień na `--dn-szary-1000` albo właściwy żeton semantyczny |
| **12** | Wartość szesnastkowa wpisana wprost w CSS widoku | rozdz. 6.2 | Zamień na `var(--dn-*)`; brakujący żeton zgłoś zespołowi systemu i architektury |
| **13** | Emoji użyte jako ikona | rozdz. 8.4 | Zamień na ikonę z `zasoby/ikony/svg/` |
| **14** | Lorem ipsum, zmyślone nazwisko, zmyślona metryka, fikcyjna firma | rozdz. 11.2 | Zamień na treść z domeny produktu, oznacz jako przykładową |
| **15** | Brak widocznego fokusu na kontrolce | rozdz. 12.2 | Przywróć `:focus-visible` z `--dn-fokus`, 2 px + odsunięcie 2 px |
| **16** | Kontrast pary tekst/tło poniżej progu WCAG AA | rozdz. 12.2 | Zamień żeton na zmierzoną parę z `kontrasty.json` |
| **17** | Czwarty krój pisma w widoku | rozdz. 7.2 | Usuń; przypisz treść do jednej z trzech ról |
| **18** | Element osiągalny wyłącznie myszą | rozdz. 12.2 | Dodaj obsługę klawiatury i kolejność `Tab` |

### 15.3. Naruszenia poważne

| Nr | Naruszenie | Reguła | Naprawa |
|---|---|---|---|
| **1** | Godło poza paskiem w widoku niebędącym „głośnym" | rozdz. 4 | Usuń wystąpienie; oznaczenie dziedziczy z ramy |
| **2** | Godło w wersji na jasnym tle użyte na ramie | rama atramentowa | Zamień na `sygnet-na-ciemnym.svg` |
| **3** | Sygnet pełny użyty poniżej 24 px | rozdz. 4.1 | Zamień na `sygnet-uproszczony*.svg` |
| **4** | Druga barwa akcentu poza rodziną sygnału i rodzinami stanów | jeden akcent | Usuń; przypisz znaczenie do sygnału albo do stanu |
| **5** | Budżet akcentu przekroczony o < 3 punkty procentowe | jeden akcent | Zamień wypełnienia na obrysy |
| **6** | Gradient jako tło przycisku, karty albo sekcji | rozdz. 6.2 | Zamień na powierzchnię kryjącą |
| **7** | Rozmycie (glassmorfizm) poza nakładką modala | rozdz. 6.2 | Zamień na `--dn-powierzchnia` |
| **8** | Drugi rodzaj ruchu ciągłego | jeden ruch ciągły | Usuń; zostaw tętno |
| **9** | Tętno na elemencie bez biegnącej pracy | jeden ruch ciągły | Usuń klasę `--tetno`; kropka statyczna albo brak kropki |
| **10** | Parallaks, scrollytelling, karuzela w oknie roboczym | rozdz. 9.3 | Usuń |
| **11** | Space Grotesk w treści roboczej albo w oknie operacyjnym | rozdz. 7.2 | Zamień na `--dn-ff-bazowa` |
| **12** | Plex Mono w prozie i w komunikatach ciągłych | rozdz. 7.2 | Zamień na `--dn-ff-bazowa` |
| **13** | Ikony o różnej grubości obrysu w jednym widoku | rozdz. 8.2 | Ujednolić do 1,75 z zestawu |
| **14** | Ikona spoza zestawu | rozdz. 8.4 | Zamień na ikonę z manifestu; brak pojęcia zgłoś zespołowi systemu i architektury |
| **15** | Emblemat środowiska bez kropki albo z więcej niż jedną kropką | rozdz. 8.5 | Przywróć plik z `zasoby/marka/srodowiska/` |
| **16** | Komunikat błędu bez następnego kroku | rozdz. 10.3 | Dopisz zdanie z krokiem |
| **17** | Żargon marketingowy, wykrzyknik, emoji w tekście interfejsu | rozdz. 10.2 | Przeredaguj rzeczowo |
| **18** | Nazwa własna przetłumaczona albo sparafrazowana | rozdz. 10.1 | Przywróć zapis z dokumentacji |
| **19** | Wyszarzenie kontrolki jako komunikat o niedostępności | zero blokad | Usuń wyszarzenie; dodaj opis obok |
| **20** | Naciśnięcie bez reakcji zamiast komunikatu | zero blokad | Dodaj toast z przyczyną i krokiem |
| **21** | Dane wyglądające na prawdziwe bez oznaczenia „przykładowe" | rozdz. 11.2 | Dodaj oznaczenie |
| **22** | Kontrolka ikonowa bez `aria-label` | rozdz. 12.2 | Dodaj `aria-label` po polsku |
| **23** | Więcej niż trzy klasy semantyczne nadawców w oknie komunikacji | stan nigdy samym kolorem | Zredukuj do `--czlowiek / --inteligencja / --system` + plakietka roli |
| **24** | Kontrolki na pasku bez modyfikatora `--na-ramie` | rama atramentowa | Dodaj `.dn-btn-ikona--na-ramie` |
| **25** | Stopień `--dn-fs-display` w widoku cichym albo milczącym | rozdz. 7.3 | Zejdź do `--dn-fs-xl` lub niżej |

### 15.4. Naruszenia drobne

| Nr | Naruszenie | Rozdział | Naprawa |
|---|---|---|---|
| **1** | Odstęp spoza skali 4 px | kontrakt systemu projektowego | Zamień na `--dn-od-*` |
| **2** | Promień spoza skali `--dn-r-*` | kontrakt systemu projektowego | Zamień na żeton |
| **3** | Czas przejścia wpisany wprost | 9.4 | Zamień na `--dn-czas-*` |
| **4** | Krzywa inna niż `--dn-ease` | 9.4 | Zamień na `--dn-ease` |
| **5** | Etykieta wersalikowa bez rozstrzelenia | 7.4 | Dodaj `--dn-ls-wersaliki` / `--dn-ls-mono-wersaliki` |
| **6** | Liczby w kolumnie bez `tabular-nums` | 7.4 | Dodaj `--dn-ff-mono` + `tabular-nums` |
| **7** | Ikona w rozmiarze spoza 14/16/20/24 | 8.3 | Zamień na żeton `--dn-wym-ikona*` |
| **8** | Cień spoza kompletu `--dn-cien-*` | 6.2 | Zamień na żeton |
| **9** | Przejście dłuższe niż 220 ms poza modalem | 9.3 | Skróć do `--dn-czas-3` |
| **10** | Kaskadowe wejście elementów listy | 9.3 | Usuń opóźnienia |
| **11** | Animacja `width`/`height` zamiast `transform`/`opacity` | 9.3 | Przenieś na `transform` |
| **12** | Mieszanie formy bezosobowej i zwrotu do Operatora w jednym widoku | 10.1 | Ujednolić |
| **13** | Ikona dekoracyjna bez `aria-hidden="true"` | 8.4 | Dodaj atrybut |
| **14** | Pole ochronne znaku naruszone | 4.1 | Zwiększ odstęp |
| **15** | Etykieta stanu w żargonie zamiast w nazwie stanu | stan nigdy samym kolorem | Zamień na nazwę stanu |

### 15.5. Rozkład naruszeń — podsumowanie

```
  KRYTYCZNE  18 ████████████████████████████████████  → NIEZGODNY natychmiast
  POWAŻNE    25 ██████████████████████████████████████████████████  → 3 = NIEZGODNY
  DROBNE     15 ██████████████████████████████  → DO POPRAWY
             ──
  RAZEM      58 skatalogowanych naruszeń
```

---

## 16. Decyzje projektowe

Poniższe rozstrzygnięcia nie wynikają wprost z dokumentacji źródłowej.
Zostały podjęte w tym opracowaniu, zgodnie z kontraktem systemu projektowego,
i wymagają odnotowania.

| # | Decyzja | Alternatywy rozważone | Uzasadnienie wyboru |
|---:|---|---|---|
| **1** | Wprowadzenie pojęcia **brand view** jako pięciu warstw (znak · rama · język · słowo · postawa) | (a) brand view jako sam zbiór reguł znaku; (b) brand view jako rozdział księgi marki | Pięć warstw pozwala uporządkować trwałość: znak zmienia Właściciel, budżety zmienia zespół marki, żetony zmienia zespół systemu. Bez tej warstwowości procedura odstępstwa nie miałaby adresata. |
| **2** | Pięć reguł konstytutywnych zamiast dłuższej listy | (a) dziesięć reguł; (b) reguły per warstwa (barwa, typografia, ruch…) | Pięć reguł da się zapamiętać i wyrecytować przy odbiorze. Reguły per warstwa powielałyby dokumentację systemu; reguły konstytutywne dotyczą **sposobu objawiania się marki**, nie definicji żetonów. |
| **3** | Budżet akcentu **5%** jako pułap ogólny, z budżetami szczegółowymi 1–8% per widok | (a) jedna wartość dla wszystkich widoków; (b) brak liczby, ocena jakościowa | Jedna wartość byłaby zbyt luźna dla okna startowego i zbyt ciasna dla AOD. Ocena jakościowa nie daje się sprawdzić przy odbiorze — a polityka ma być sprawdzalna. Wartości szczegółowe wyprowadzone z liczby jednoczesnych wskazań w każdym widoku. |
| **4** | Metoda pomiaru budżetu: siatka **20 × 20 pól** (1 pole = 0,25%) | (a) pomiar pikselowy narzędziem; (b) inwentaryzacja elementów | Siatka daje powtarzalny wynik w minutę, bez narzędzi, przy `file://`. Pomiar pikselowy jest dokładniejszy, ale nie jest wykonalny podczas odbioru projektu (przed implementacją). Metoda B (inwentaryzacja) pozostaje jako rozstrzygająca w sporze. |
| **5** | Trzy poziomy głośności marki (**głośna / cicha / milcząca**) z przypisaniem każdego widoku | (a) dwa poziomy (marka jest / marki nie ma); (b) głośność jako decyzja projektanta per widok | Dwa poziomy nie oddają różnicy między powłoką (marka oznacza) a oknem operacyjnym (marka milczy, ale rama nad nim istnieje). Decyzja per widok prowadziłaby do niespójności między modułami. |
| **6** | Wprowadzenie **terytoriów milczenia bezwzględnego** (treść dokumentu, terminal, edytor kodu) | (a) objęcie tych obszarów zwykłym budżetem 5%; (b) brak rozróżnienia | Produkt, który koloruje cudzą treść własnymi barwami, kłamie o źródle. Wyłączenie podświetlenia składni z budżetu (odstępstwo stałe dla edytora kodu i terminala) jest tego konsekwencją — inaczej edytor musiałby być monochromatyczny. |
| **7** | Sygnet **26 px** na pasku 48 px jako rozmiar standardowy w widokach cichych | (a) 20 px; (b) 32 px | 26 px daje pole ochronne 11 px w pasku 48 px (48 − 26 = 22, po 11 z każdej strony) i zachowuje czytelność dwóch grotów. 20 px zbliżałoby do progu wariantu uproszczonego; 32 px nie mieści pola ochronnego. |
| **8** | **Trzy widoki nad ramą** (AOD, Centrum poleceń, Mobile) dostają sygnet uproszczony 16–20 px | (a) brak znaku; (b) sygnet pełny | Warstwy `--dn-z-aod` (1200) i `--dn-z-centrum-polecen` (1300) leżą nad paskiem, więc godło traci swój dom. Uproszczony sygnet przywraca oznaczenie przy koszcie powierzchni poniżej 0,1%. |
| **9** | Odstępstwo stałe dla Always On Display dla AOD (budżet 8%) zamiast wyjątku milczącego | (a) trzymanie 5% i zmniejszenie rdzenia; (b) brak budżetu dla AOD | Rdzeń AOD jest wprost wymieniony w kierunku systemu projektowego jako wystąpienie elementu sygnaturowego. Zapisanie tego jako odstępstwa stałego pokazuje procedurę rozdz. 13 w działaniu i chroni przed rozlewaniem się wyjątku na inne widoki. |
| **10** | Lista kontrolna ma **39 pozycji w 8 blokach** | (a) 25 pozycji minimum; (b) 60 pozycji szczegółowych | 25 nie pokrywa wszystkich pięciu reguł plus czterech polityk warstwowych plus dostępności. 60 pozycji nie zostałoby przeprowadzone w praktyce. 39 pozycji = ok. 25 minut odbioru jednego widoku. |
| **11** | Werdykt trójstopniowy z progiem „**3 poważne = niezgodny**" | (a) próg 5; (b) tylko krytyczne decydują | Trzy naruszenia poważne w jednym widoku oznaczają, że wykonawca nie czytał polityki — poprawianie punktowe jest wtedy droższe niż powtórzenie widoku. |
| **12** | **Katalog naruszeń** jako osobny rozdział uporządkowany wagą naruszenia, z kolumną „naprawa" | (a) naruszenia opisane przy regułach; (b) brak katalogu | Podział na wagi pozwala zapisać werdykt jako liczbę naruszeń w każdej z trzech kategorii, co czyni odbiór porównywalnym między widokami i w czasie. Kolumna „naprawa" skraca iterację — wykonawca dostaje rozwiązanie, nie tylko zarzut. |
| **13** | Terminologia obowiązująca (rozdz. 10.4, załącznik C) ze wskazaniem nazw **zakazanych** | (a) tylko lista nazw obowiązujących | Kolizje są realne: `Workspace` jest jednocześnie modułem i komponentem własnym, więc „workspace" jako synonim środowiska wprowadzałby w błąd. Podobnie `Agents` (moduł) vs. „agent" (potoczna nazwa roli). |
| **14** | Komunikat błędu jako konstrukcja **dwuczęściowa** (stan + następny krok) | (a) sam stan; (b) trójczęściowa (stan + przyczyna + krok) | Sam stan zamienia komunikat w bramę i narusza regułę „zero blokad". Trzecia część (przyczyna techniczna) należy do panelu diagnostycznego, nie do komunikatu w widoku roboczym. |
| **15** | Odstępstwo stałe dla ekranu głównego Mobile | (a) sygnet + skrócony logotyp; (b) sam sygnet | Pasek na progu `w1` mieści godło **albo** logotyp obok kontrolek. Wybrano logotyp, bo na urządzeniu mobilnym Operator częściej wchodzi do produktu „od zera" i nazwa niesie więcej niż znak. |

---

## Załącznik A — karta odbioru widoku (do wydruku)

```
╔══════════════════════════════════════════════════════════════════════════╗
║  DANACO CONSOLE · KARTA ODBIORU BRAND VIEW                     v2.0      ║
╠══════════════════════════════════════════════════════════════════════════╣
║  Widok ................................................................  ║
║  Głośność marki:   [ ] głośna   [ ] cicha   [ ] milczy                   ║
║  Budżet akcentu wg tabeli 6.4: ......... %   Pomiar: ......... %         ║
║  Data ..................  Osoba odbierająca ...........................  ║
╠══════════════════════════════════════════════════════════════════════════╣
║  BLOK · REGUŁY KONSTYTUTYWNE                                             ║
║  [ ] rama atramentowa w obu motywach                 krytyczne           ║
║  [ ] godło tylko na pasku (lub wg rozdz. 4)          poważne             ║
║  [ ] godło w barwach własnych, wariant ciemny        poważne             ║
║  [ ] budżet akcentu dotrzymany                       poważne             ║
║  [ ] sygnał nie jest tłem sekcji ani przycisku       krytyczne           ║
║  [ ] jeden rodzaj ruchu ciągłego                     poważne             ║
║  [ ] każdy stan z ikoną albo etykietą                krytyczne           ║
║  [ ] zero blokad w kodzie                            krytyczne           ║
║                                                                          ║
║  BLOK · ZNAK                  [ ] [ ] [ ] [ ] [ ]                        ║
║  BLOK · BARWA                 [ ] [ ] [ ] [ ] [ ] [ ]                    ║
║  BLOK · TYPOGRAFIA            [ ] [ ] [ ] [ ]                            ║
║  BLOK · IKONOGRAFIA           [ ] [ ] [ ]                                ║
║  BLOK · RUCH                  [ ] [ ] [ ]                                ║
║  BLOK · JĘZYK I TREŚĆ         [ ] [ ] [ ] [ ] [ ]                        ║
║  BLOK · DOSTĘPNOŚĆ            [ ] [ ] [ ] [ ] [ ]                        ║
╠══════════════════════════════════════════════════════════════════════════╣
║  SPEŁNIONE: ...... / 39                                                  ║
║  NARUSZENIA — krytyczne: ......  poważne: ......  drobne: ......         ║
║  Uwagi: ...............................................................  ║
║  Odstępstwa powołane: .................................................  ║
╠══════════════════════════════════════════════════════════════════════════╣
║  WERDYKT:   [ ] ZGODNY     [ ] DO POPRAWY     [ ] NIEZGODNY              ║
║  Termin ponownego odbioru: ............................................  ║
╚══════════════════════════════════════════════════════════════════════════╝
```

---

## Załącznik B — mapa źródeł polityki

| Rozdział polityki | Źródło wiążące | Miejsce |
|---|---|---|
| 2 · Zasada nadrzędna | kierunek systemu projektowego — deklaracja kierunku | § 1 |
| 3 · Rama atramentowa | kierunek systemu projektowego · kontrakt systemu projektowego (Rama kokpitu) | wprost |
| 3 · Jeden akcent | kierunek systemu projektowego · kontrakt systemu projektowego (Sygnał) | wprost; **budżet 5% — decyzja nr 3** |
| 3 · Jeden ruch | kierunek systemu projektowego · kontrakt systemu projektowego (Ruch) | wprost |
| 3 · Stan nie kolorem | kierunek systemu projektowego · kontrakt systemu projektowego | wprost |
| 3 · Zero blokad | zasada zero blokad · kontrakt systemu projektowego | wprost |
| 4 · Obecność znaku | kierunek systemu projektowego · kontrakt systemu projektowego (inwentarz okien) | forma znaku wprost; **przypisanie per widok — decyzje nr 7, 8** |
| 5 · Hierarchia obecności | — | **decyzja nr 5, 6** |
| 6 · Barwa | kontrakt systemu projektowego · `zasoby/zetony/zetony.css` | wprost; **budżety — decyzja nr 3** |
| 7 · Typografia | kierunek systemu projektowego · kontrakt systemu projektowego | wprost |
| 8 · Ikonografia | kierunek systemu projektowego · kontrakt systemu projektowego | wprost |
| 9 · Ruch | kierunek systemu projektowego · kontrakt systemu projektowego (wzorzec animacji) | wprost |
| 10 · Język | kontrakt systemu projektowego | wprost; **słownik terminów — decyzja nr 13** |
| 11 · Treści przykładowe | kontrakt systemu projektowego (anty-domyślne), § 10 pkt 3–4 | wprost |
| 12 · Dostępność | kontrakt systemu projektowego · `zasoby/zetony/kontrasty.json` | wprost |
| 13 · Odstępstwa | — | **decyzje nr 9, 15** |
| 14 · Odbiór | kontrakt systemu projektowego (lista sprawdzeń) rozszerzona | **decyzje nr 10, 11** |
| 15 · Katalog naruszeń | — | **decyzja nr 12** |

**Pliki znaku odczytane i zachowane co do współrzędnej:**
`zasoby/marka/logo/sygnet.svg` · `sygnet-na-ciemnym.svg` · `sygnet-uproszczony.svg` ·
`logo-poziomy.svg` · `logo-poziomy-na-ciemnym.svg` · `logotyp.svg` · `logo-pionowy.svg` ·
`zasoby/marka/srodowiska/srodowisko-{talkin,workspace,codestudio,multitaskingai}.svg` ·
`zasoby/marka/favicon/favicon.svg` · `zasoby/marka/ikona-aplikacji/ikona-aplikacji.svg`

**Uwaga o poprzednie wydanie przewodnika marki:** dokument opisuje warstwę wizualną v1.0
(granat + złoto, styl „kancelaryjno-nowoczesny"), uznaną przez Właściciela
za zastępczą i nieobowiązującą (kierunek systemu projektowego nagłówek „Zastępuje").
Jego wartości wizualne **nie obowiązują**; obowiązuje wyłącznie zawarta w nim
zasada oszczędności akcentu, przeniesiona tutaj jako reguła „jeden akcent" w nowej rodzinie barw.

---

## Załącznik C — słownik terminów obowiązujących

Terminy poniżej są **wiążące w interfejsie, w dokumentacji i w komunikacji
zespołów**. Kolumna „nie mówimy" wymienia nazwy zakazane.

| Termin | Znaczenie | Nie mówimy |
|---|---|---|
| **Operator** | Zawodowy użytkownik Danaco Console, zarządzający cyfrową organizacją złożoną z modeli, agentów, projektów i automatyzacji. Jedyna nazwa roli człowieka w produkcie. | użytkownik, klient, końcowy odbiorca |
| **Środowisko** | Jeden z czterech trybów pracy: TalkIn, WorkSpace, CodeStudio, MultitaskingAI. Odpowiada na pytanie „w jakim trybie pracuję?". | tryb, przestrzeń, workspace |
| **Moduł** | Jedna z piętnastu jednostek zadaniowych dostępnych w bocznej nawigacji środowiska. Odpowiada na pytanie „jakie zadanie wykonuję?". | sekcja, zakładka, aplikacja |
| **Okno operacyjne** | Konkretne narzędzie wewnątrz modułu (np. `Studio Editor`, `Workflow Builder`). Odpowiada na pytanie „jakim narzędziem realizuję?". | ekran, widok modułu |
| **Powłoka** | Stały układ środowiska: pasek górny, pas kart sesji, boczna nawigacja albo panel orkiestracji, obszar roboczy. | layout, szkielet, rama aplikacji |
| **Karta sesji** | Pozycja w pasie kart reprezentująca jedną sesję pracy; niesie kropkę sygnału, gdy w tle biegnie praca. | tab, zakładka |
| **Komponent własny** | Element Strefy 2 Centrum dowodzenia: `Automations`, `Agents`, `Workspace`, `Assistant`. Wpina się w sesję modułu. | widżet, skrót, dodatek |
| **Ekspert** | Model przypisany do roli w środowisku MultitaskingAI. | bot, silnik |
| **Rola** | Jednostka wykonawcza w MultitaskingAI: `Executor 1`, `Executor 2`, `Coordinator`, `Executor 3 / Validator`, `Subagent Network`. | agent, węzeł, worker |
| **Kolejka** | Lista zadań silnika MultitaskingAI z jedenastoma akcjami; każde zadanie ma stan komunikowany ikoną, etykietą i barwą. | lista zadań, stos, pipeline |
| **Chat Window** | Okno komunikacji wspólne wszystkim modułom, rekonfigurowane w kontekście każdego z nich. Nazwa własna — nie tłumaczymy. | czat, okno rozmowy |
| **Centrum dowodzenia** | Strona główna platformy; trzy strefy: środowiska, komponenty własne, ustawienia. | dashboard, pulpit, strona startowa |
| **Okno startowe** | Widok ładowania przed rejestracją i logowaniem; stany: Łączenie, Powrót z ważnym tokenem, Błąd połączenia. | splash, ekran powitalny |
| **Always On Display** | Warstwa centralna nad interfejsem z rdzeniem awatara i monitorem procesu. Skrót AOD dopuszczalny w dokumentacji, nie w interfejsie. | tryb zawsze aktywny, overlay |
| **Panel orkiestracji** | Sześciosekcyjny panel MultitaskingAI zastępujący boczną nawigację modułów. | nawigacja, menu ról |
| **Rama** | Atramentowy pasek górny, stały w obu motywach; jedyny stały dom godła. | header, nagłówek strony |
| **Godło** | Znak marki w użyciu jako oznaczenie: sygnet albo sygnet z logotypem. | logo (dopuszczalne potocznie, nie w dokumentacji) |
| **Sygnet** | Sam znak graficzny „»»." bez logotypu; warianty: pełny (dwa groty) i uproszczony (jeden grot). | ikona marki, symbol |
| **Logotyp** | `DANACO` (Space Grotesk 700) + `CONSOLE` (IBM Plex Mono 500, rozstrzelone wersaliki). | wordmark, napis |
| **Emblemat środowiska** | Znak domenowy jednego z czterech środowisk; siatka 24 × 24, obrys 1,75, dokładnie jedna wypełniona kropka. | ikona środowiska (dopuszczalne technicznie), logo środowiska |
| **Kropka sygnału** | Wypełniony punkt błękitu sygnałowego oznaczający „tu biegnie praca". Element sygnaturowy marki. | wskaźnik, dioda, badge |
| **Tętno** | Animacja kropki sygnału o okresie 2,4 s; jedyny ruch ciągły w produkcie. | puls, animacja, migotanie |
| **Sygnał** | Błękit `#3B6FE0` i jego rodzina; jedyna barwa akcentu. Cztery role: fokus, aktywność, odnośnik, praca w tle. | kolor akcentu, primary, brand color |
| **Żeton** | Zmienna projektowa `--dn-*`; jedyny dopuszczony nośnik wartości w CSS widoku. | token (dopuszczalne w kodzie), zmienna |
| **Brand view** | Sposób, w jaki tożsamość marki jest widoczna i egzekwowana w pojedynczym widoku produktu. Pięć warstw: znak, rama, język, słowo, postawa. | branding, identyfikacja wizualna |
| **Budżet akcentu** | Maksymalny udział powierzchni widoku zajęty przez rodzinę sygnału; wartości per widok w tabeli 6.4. | limit koloru |
| **Punkt izolacji** | Konfigurowalne miejsce ograniczenia dostępu; stan wyjściowy „wyłączone / pełny dostęp" zgodnie z zasadą zero blokad. | blokada, zabezpieczenie |
| **Panel prowenancji** | Element okna Konfiguracji prezentujący pochodzenie danych i decyzji. | audyt, historia |
| **Zero blokad** | Zasada nadrzędna platformy (zasada zero blokad): brak `disabled`, komunikat zamiast bramy. | soft lock, tryb bezpieczny |

---

*Danaco Console — AI Operating Environment · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. — Dariusz Naharnowicz*
