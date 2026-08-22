# Danaco Console — Okno instrukcji użytkowania

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
| **Tytuł** | Okno instrukcji użytkowania — nakładka wywoływana z menu Pomoc |
| **Klasa dokumentu** | Specyfikacja docelowa |
| **Odbiorcy** | deweloper · projektant · redaktor dokumentacji |
| **Przeznaczenie** | Ustala pełny ciąg budowy okna nakładkowego „Instrukcja użytkowania”: układ, nawigację po rozdziałach, dosłowną treść rozdziałów 1–4 i zapowiedzi rozdziałów 5–9, makiety osadzone w treści, stopkę, żetony, komponenty i kryteria odbioru |
| **Zakres** | struktura nakładki, spis rozdziałów, treść rozdziałów 1–4 w pełnym brzmieniu, zapowiedzi rozdziałów 5–9 wraz z tabelą skrótów klawiszowych, makiety układu ramy i przestrzeni roboczej, stopka okna, warstwy widoczności, zachowanie na punktach łamania, relacja z dokumentem `../INSTRUKCJA-UZYTKOWANIA.md` |
| **Poza zakresem** | treść merytoryczna platformy poza streszczeniem niesionym przez instrukcję — [Rama okna aplikacji](rama-okna.md), [Przepływ okien](przeplyw-okien.md); pełna dokumentacja platformy otwierana klawiszem `F1` — poza zbiorem niniejszego katalogu; okno instalatora — [Okno instalatora](okno-instalatora.md) |
| **Dokument nadrzędny** | [Elementy okien przepływu głównego](elementy-okien.md) |
| **Dokumenty powiązane** | [Rama okna aplikacji](rama-okna.md) · [Przepływ okien](przeplyw-okien.md) · [Strona główna i nawigacja](strona-glowna-i-nawigacja.md) · [Katalog komponentów](katalog-komponentow.md) · [Okno instalatora](okno-instalatora.md) · [Instrukcja użytkowania](../INSTRUKCJA-UZYTKOWANIA.md) · [Instalacja i konfiguracja](../INSTALACJA-I-KONFIGURACJA.md) · [Standard redakcyjny i językowy](../STANDARD-REDAKCYJNY-I-JEZYKOWY.md) |
| **Prototypy odniesienia** | `design/05-okna/platformowe/instrukcja-uzytkowania.html` |
| **Źródła normatywne** | `design/05-okna/platformowe/instrukcja-uzytkowania.html` · `design/zasoby/zetony/zetony.css` · `design/zasoby/css/komponenty.css` · `design/zasoby/rama.css` · `design/zasoby/okna-modalne.css` |
| **Zasada nadrzędna** | Instrukcja otwiera się nad pracą, nie zamiast niej — zamyka się jednym naciśnięciem i przywraca dokładnie to miejsce, z którego została wywołana |

---

## Spis treści

1. [Czym jest to okno](#1-czym-jest-to-okno)
2. [Wywołanie z menu Pomoc](#2-wywołanie-z-menu-pomoc)
3. [Uwaga o nazewnictwie dwóch dokumentów](#3-uwaga-o-nazewnictwie-dwóch-dokumentów)
4. [Nakładka jako okno — rama nakładki](#4-nakładka-jako-okno--rama-nakładki)
5. [Układ okna — dwie kolumny](#5-układ-okna--dwie-kolumny)
6. [Spis rozdziałów — lewa kolumna](#6-spis-rozdziałów--lewa-kolumna)
   - [6.1 Wykaz dziewięciu pozycji](#61-wykaz-dziewięciu-pozycji)
   - [6.2 Stany pozycji spisu](#62-stany-pozycji-spisu)
   - [6.3 Nota zamykająca spis](#63-nota-zamykająca-spis)
7. [Treść rozdziału — prawa kolumna, jednostki typograficzne](#7-treść-rozdziału--prawa-kolumna-jednostki-typograficzne)
8. [Rozdział 1 instrukcji — Czym jest Danaco Console](#8-rozdział-1-instrukcji--czym-jest-danaco-console)
   - [8.1 Cztery środowiska (`.iu-def`)](#81-cztery-środowiska-iu-def)
   - [8.2 Moduły](#82-moduły)
   - [8.3 Sesje i projekty](#83-sesje-i-projekty)
9. [Rozdział 2 instrukcji — Pierwsze uruchomienie](#9-rozdział-2-instrukcji--pierwsze-uruchomienie)
   - [9.1 Treść dosłowna ośmiu etapów](#91-treść-dosłowna-ośmiu-etapów)
10. [Rozdział 3 instrukcji — Rama okna](#10-rozdział-3-instrukcji--rama-okna)
   - [10.1 Makieta ramy okna](#101-makieta-ramy-okna)
   - [10.2 Co gdzie stoi i dlaczego (`.iu-def`)](#102-co-gdzie-stoi-i-dlaczego-iu-def)
11. [Rozdział 4 instrukcji — Praca w środowisku](#11-rozdział-4-instrukcji--praca-w-środowisku)
   - [11.1 Przedsionek środowiska](#111-przedsionek-środowiska)
   - [11.2 Wybór modułu i układ przestrzeni roboczej — druga makieta](#112-wybór-modułu-i-układ-przestrzeni-roboczej--druga-makieta)
   - [11.3 Karty sesji](#113-karty-sesji)
   - [11.4 Trwałość pracy w tle](#114-trwałość-pracy-w-tle)
12. [Scenariusze użycia](#12-scenariusze-użycia)
   - [12.1 Pierwsze otwarcie w toku pracy](#121-pierwsze-otwarcie-w-toku-pracy)
   - [12.2 Sprawdzenie skrótu klawiszowego w trakcie pracy nad dokumentem](#122-sprawdzenie-skrótu-klawiszowego-w-trakcie-pracy-nad-dokumentem)
   - [12.3 Przeciąganie i zmiana rozmiaru nakładki](#123-przeciąganie-i-zmiana-rozmiaru-nakładki)
   - [12.4 Zawężenie okna aplikacji poniżej 900 px](#124-zawężenie-okna-aplikacji-poniżej-900-px)
   - [12.5 Operator instalujący produkt trafia na rozbieżność rozdz. 13.4](#125-operator-instalujący-produkt-trafia-na-rozbieżność-rozdz-134)
13. [Rozdziały 5–9 — zapowiedzi i tabela skrótów klawiszowych](#13-rozdziały-59--zapowiedzi-i-tabela-skrótów-klawiszowych)
   - [13.1 Rozdział 5 — Projekty](#131-rozdział-5--projekty)
   - [13.2 Rozdział 6 — Dostosowanie pasków](#132-rozdział-6--dostosowanie-pasków)
   - [13.3 Rozdział 7 — Skróty klawiszowe](#133-rozdział-7--skróty-klawiszowe)
   - [13.4 Rozdział 8 — Instrukcja instalacji](#134-rozdział-8--instrukcja-instalacji)
   - [13.5 Rozdział 9 — Wersja i licencja](#135-rozdział-9--wersja-i-licencja)
14. [Rozwinięcie rozdziału 5 — pięć grup ustaleń projektu](#14-rozwinięcie-rozdziału-5--pięć-grup-ustaleń-projektu)
15. [Rozwinięcie rozdziału 6 — cztery zakładki okna „Dostosuj”](#15-rozwinięcie-rozdziału-6--cztery-zakładki-okna-dostosuj)
16. [Budowa makiet osadzonych w treści](#16-budowa-makiet-osadzonych-w-treści)
17. [Stopka okna](#17-stopka-okna)
18. [Wyszukiwanie w treści, tryb czytania, druk i eksport, wersjonowanie — zakres nieobjęty prototypem](#18-wyszukiwanie-w-treści-tryb-czytania-druk-i-eksport-wersjonowanie--zakres-nieobjęty-prototypem)
19. [Warstwy widoczności i dostępność](#19-warstwy-widoczności-i-dostępność)
20. [Zachowanie na punktach łamania](#20-zachowanie-na-punktach-łamania)
21. [Skróty klawiszowe okna](#21-skróty-klawiszowe-okna)
22. [Żetony `--dn-*` i `--iu-*` użyte w widoku](#22-żetony---dn--i---iu--użyte-w-widoku)
23. [Komponenty `.dn-*` i `.iu-*` użyte w widoku](#23-komponenty-dn--i-iu--użyte-w-widoku)
24. [Etykiety interfejsu](#24-etykiety-interfejsu)
   - [24.1 Etykiety sekcji Pomoc przywołane w rozdz. 2](#241-etykiety-sekcji-pomoc-przywołane-w-rozdz-2)
25. [Komunikaty](#25-komunikaty)
26. [Komendy kontraktu — nie dotyczy](#26-komendy-kontraktu--nie-dotyczy)
27. [Relacja z `../INSTRUKCJA-UZYTKOWANIA.md` i kryteria odbioru](#27-relacja-z-instrukcja-uzytkowaniamd-i-kryteria-odbioru)
   - [27.1 Tabela zgodności rozdziałów](#271-tabela-zgodności-rozdziałów)
   - [27.2 Cztery okna platformowe tego przejścia redakcyjnego — zestawienie](#272-cztery-okna-platformowe-tego-przejścia-redakcyjnego--zestawienie)
   - [27.3 Kryteria odbioru](#273-kryteria-odbioru)

---

## 1. Czym jest to okno

Instrukcja użytkowania jest **oknem nakładkowym** (`<dialog class="dn-modal iu-okno">`) wywoływanym
z sekcji „Pomoc” menu aplikacji. Nie zajmuje całego okna platformy i nie zastępuje pracy w toku:
otwiera się nad bieżącym widokiem, karty sesji pracują dalej pod nią, a Operator wraca do dokładnie
tego samego miejsca jednym zamknięciem. Ciało nakładki jest dwukolumnowe — po lewej spis dziewięciu
rozdziałów z zaznaczoną pozycją bieżącą, po prawej treść rozdziału wskazanego. W rozdziałach o ramie
okna i o pracy w środowisku treść niesie dodatkowo **makiety układu** narysowane samym HTML-em
i żetonami CSS — schemat ramy okna oraz schemat przestrzeni roboczej — bez odwołania do zrzutu ekranu
i bez elementu graficznego spoza systemu projektowego.

Instrukcja **nie wprowadza żadnej funkcji, której nie ma w dokumentacji interfejsu produktu**. Jest streszczeniem:
każdy akapit rozdziału daje się wskazać w jednym z dokumentów źródłowych wymienionych w polu „Dokumenty
powiązane” metryki. Zmiana źródła pociąga aktualizację odpowiedniego rozdziału instrukcji — reguła
jednego źródła prawdy ([Standard redakcyjny i językowy](../STANDARD-REDAKCYJNY-I-JEZYKOWY.md) rozdz. 7.4) obowiązuje instrukcję tak samo
jak każde inne opracowanie zbioru.

```
┌──────────────────────────────────────────────────────────────────────────────┐
│  Trzy fakty, które rozstrzygają o kształcie tego okna                        │
├──────────────────────────────────────────────────────────────────────────────┤
│  1. Nakładka, nie pełny widok — praca w toku nie zostaje przerwana.          │
│  2. Dwie kolumny — spis rozdziałów stale widoczny, treść przewija się        │
│     niezależnie od spisu.                                                    │
│  3. Streszczenie, nie źródło — każde zdanie ma pokrycie w innym dokumencie   │
│     zbioru; instrukcja niczego nie ustala pierwsza.                          │
└──────────────────────────────────────────────────────────────────────────────┘
```

Prototyp deklaruje siebie w atrybutach `data-prototyp` znacznika `<body>`: `data-prototyp="OKNA
PLATFORMOWE · INSTRUKCJA UŻYTKOWANIA"` oraz `data-prototyp-okno="Instrukcja użytkowania — nakładka
wywoływana z menu Pomoc"`. Oba brzmienia potwierdzają dwa fakty rozstrzygnięte w niniejszym rozdziale:
przynależność do rodziny „okien platformowych” (a nie okien wejściowych ani okien modułów) oraz formę
nakładki, nie osobnego okna pełnoekranowego.

---

## 2. Wywołanie z menu Pomoc

Jedyną drogą otwarcia instrukcji jest **menu aplikacji** — hamburger w głowie szyny nawigacji
([Rama okna aplikacji](rama-okna.md) rozdz. 9) — sekcja **Pomoc**. Sekcja niesie dziewięć pozycji;
tabela poniżej podaje ich dosłowne brzmienie w kolejności występowania w znaczniku HTML.

| Pozycja | Rodzaj | Skrót | Skutek |
|---|---|---|---|
| Instrukcja użytkowania | otwiera okno | — | otwiera niniejszą nakładkę z rozdziałem 1 wybranym |
| Instrukcja instalacji | otwiera okno | — | otwiera [okno instalatora](okno-instalatora.md) |
| Dokumentacja platformy | czynność | `F1` | otwiera pełną dokumentację platformy — poza zbiorem niniejszego katalogu |
| Skróty klawiszowe | czynność | `Ctrl+/` | otwiera mapę skrótów klawiszowych aplikacji |
| Makiety i opracowania okien | czynność | — | otwiera zbiór prototypów `design/05-okna/` |
| Sprawdź dostępność aktualizacji | czynność | — | uruchamia sprawdzenie wydania bieżącego wobec wydań dostępnych |
| Zgłoś obserwację | czynność | — | otwiera formularz zgłoszenia obserwacji Operatora |
| Wersja i licencja | otwiera okno | — | otwiera okno „Wersja i licencja” |
| O programie | czynność | — | otwiera krótką informację o produkcie i wydawcy |

Znacznik `<button>` odpowiadający pierwszej pozycji niesie atrybut
`data-nawiguj="platformowe/instrukcja-uzytkowania.html"` — w środowisku prototypów jest to ścieżka
przejścia między plikami HTML; w aplikacji zbudowanej atrybut ten odpowiada otwarciu okna nakładkowego
o identyfikatorze `okno-instrukcji` bez przeładowania widoku. Ta sama pozycja powtarza się także jako
ostatnia w bloku „Pomoc” tuż przed rozdzielającą kreską `sta-menu-sep`, tym razem jako skrót do
rozdziału „Wersja i licencja” wewnątrz samej instrukcji (rozdz. 9 treści, rozdz. 13 niniejszego
dokumentu) — dwa różne wejścia do tego samego okna, różniące się wyłącznie rozdziałem otwieranym
domyślnie.

**Reguła doboru pozycji sekcji Pomoc.** Sekcja niesie wyłącznie czynności dotyczące aplikacji jako
całości: dokumentację, skróty, zgłoszenia, wersję. Czynność właściwa jednemu modułowi nie trafia tutaj —
ma własne miejsce w oknie roboczym tego modułu ([Rama okna aplikacji](rama-okna.md) rozdz. 9, „Reguła doboru”).

---

## 3. Uwaga o nazewnictwie dwóch dokumentów

Zbiór dokumentacji niesie **dwa różne dokumenty o zbliżonej nazwie**, których nie wolno pomylić:

| Plik | Położenie | Klasa dokumentu | Co niesie |
|---|---|---|---|
| `INSTRUKCJA-UZYTKOWANIA.md` | korzeń `docs/` | Stan wdrożenia | Instrukcja obsługi produktu w jego bieżącym, faktycznym stanie — dokument Właściciela, wielkimi literami w nazwie pliku |
| `interfejs-uzytkownika/okno-instrukcji.md` | niniejszy plik | Specyfikacja docelowa | Projekt okna, które **wewnątrz aplikacji** wyświetla treść pochodną od dokumentu powyższego |

Różnica nie jest tylko różnicą wielkości liter (zakazaną przez [Standard redakcyjny i językowy](../STANDARD-REDAKCYJNY-I-JEZYKOWY.md)
rozdz. 6 jako kolizję nazw na systemach nierozróżniających wielkości liter) — to **dwa odrębne
dokumenty o odrębnym przeznaczeniu**: jeden opisuje, co Operator może dziś zrobić z produktem, drugi
opisuje, jak zbudować okno, które o tym mówi. Dokument źródłowy `data-prototyp-zrodlo` prototypu HTML
niesie zapis `docs/interfejs-uzytkownika/instrukcja-uzytkowania.md` — nazwę wcześniejszą, sprzed
przyjęcia konwencji `okno-*.md` dla dokumentów tej rodziny ([Standard redakcyjny i językowy](../STANDARD-REDAKCYJNY-I-JEZYKOWY.md) rozdz.
7.3). Nazwa wiążąca dla niniejszego opracowania to `okno-instrukcji.md`; atrybut prototypu jest
świadectwem historycznym nazwy pliku i nie rozstrzyga o nazwie bieżącej.

**Konsekwencja dla redaktora tego opracowania.** Skoro treść dziewięciu rozdziałów jest streszczeniem
`../INSTRUKCJA-UZYTKOWANIA.md` (rozdz. 27.1), a nie ustaleniem pierwszym, każda różnica zauważona
między treścią dosłowną prototypu (rozdziały 1–4, w pełni rozwinięte w rozdz. 8–11 niniejszego
dokumentu) a treścią rzeczywistą dokumentu źródłowego jest sygnałem do przeglądu — nie automatycznie
błędem tego opracowania ani automatycznie błędem dokumentu źródłowego, lecz miejscem wymagającym
rozstrzygnięcia, które z dwóch brzmień jest aktualne. Rozdział 12.4 pokazuje przykład takiego miejsca
w praktyce: treść dosłowna prototypu (rozdział 8) i stan faktyczny wdrożenia rozeszły się, a niniejszy
dokument nazywa to rozejście wprost, zamiast wybierać jedną wersję milcząco.

---

## 4. Nakładka jako okno — rama nakładki

Instrukcja dziedziczy zachowanie wspólne wszystkim oknom nakładkowym platformy
([Rama okna aplikacji](rama-okna.md) rozdz. 21): przeciąganie, zmianę rozmiaru, minimalizację i maksymalizację. Zachowanie
dokłada `design/zasoby/okna-modalne.js` na dialogach oznaczonych `data-okno="1"` przez
przyciski `.dn-btn` w `.dn-modal-stopka` oraz wiersz `.dn-modal-naglowek`.

| Czynność | Uchwyt | Ograniczenie |
|---|---|---|
| Przeciąganie | nagłówek `.dn-modal-naglowek` | co najmniej 80 px nagłówka pozostaje w oknie |
| Zmiana rozmiaru | osiem uchwytów przy krawędziach i narożnikach | minimum 420 × 260 px |
| Minimalizacja | przycisk `--min` albo kliknięcie nagłówka | okno dokuje przy dolnej krawędzi od prawej |
| Maksymalizacja | przycisk `--max` albo dwuklik nagłówka | wypełnia okno z marginesem 16 px |

Klawiaturą: gdy nagłówek ma ognisko, strzałki przesuwają okno, `Shift`+strzałki zmieniają rozmiar (krok
16 px), `Escape` zamyka okno w całości niezależnie od tego, który rozdział jest otwarty. Stan (położenie,
rozmiar, zminimalizowanie) trzymają atrybuty `data-okno-*`. Przy `prefers-reduced-motion` przejścia
przeciągania i zmiany rozmiaru ustają — informacja o stanie okna nie ginie, zmienia się wyłącznie sposób
jej podania.

**Wymiar własny nakładki.** Poza wspólnym zachowaniem ramy nakładkowej, instrukcja ustala własny wymiar
bazowy niezależny od pozostałych okien nakładkowych platformy:

```
szerokość:    min(1180px, calc(100vw - var(--dn-od-8)))
wysokość maks.: min(88dvh, 820px)
```

Klasa `.dn-modal.iu-okno` układa wnętrze w kolumnę: nagłówek (`.dn-modal-naglowek`) i stopka
(`.dn-modal-stopka`) stoją nieruchomo (`flex: none`), a między nimi przewija się wyłącznie ciało
(`.dn-modal-cialo`) — inaczej stopka niosąca wersję, numer kompilacji i odnośnik do licencji wychodziłaby
poza wysokość okna przy dłuższym rozdziale.

---

## 5. Układ okna — dwie kolumny

Ciało nakładki (`.dn-modal-cialo`) mieści jeden kontener `.iu-korpus` o siatce dwukolumnowej: lewa
kolumna o stałej szerokości 272 px niesie spis rozdziałów, prawa — o szerokości elastycznej
(`minmax(0, 1fr)`) — treść rozdziału wskazanego. Obie kolumny przewijają się niezależnie
(`overflow-y: auto` na obu), więc długi spis i długi rozdział nie wymuszają wspólnego przewijania.

```
┌────────────────────────────────────────────────────────────────────────────────┐
│ Instrukcja użytkowania                                                    ✕    │  ← .dn-modal-naglowek
├────────────────────┬───────────────────────────────────────────────────────────┤
│ ROZDZIAŁY          │  Rozdział 1                                               │
│                    │  Czym jest Danaco Console                                 │
│ 1 Czym jest        │                                                           │
│   Danaco Console ► │  Danaco Console jest platformą klasy AI Workspace OS…     │
│ 2 Pierwsze         │                                                           │
│   uruchomienie     │  Cztery środowiska                                        │
│ 3 Rama okna        │  TalkIn        Środowisko rozmowy i pracy z wiedzą…       │
│ 4 Praca w          │  WorkSpace     Środowisko pracy nad treścią…              │
│   środowisku       │  CodeStudio    Środowisko wytwarzania oprogramowania…     │
│ 5 Projekty         │  MultitaskingAI  Środowisko orkiestracji pracy zespołu…   │
│ 6 Dostosowanie     │                                                           │
│   pasków           │  Moduły                                                   │
│ 7 Skróty           │  Platforma udostępnia piętnaście modułów…                 │
│   klawiszowe       │                                                           │
│ 8 Instrukcja       │  Sesje i projekty                                         │
│   instalacji       │  Sesja zawsze należy do środowiska…                       │
│ 9 Wersja i         │                                                           │
│   licencja         │                                                           │
│                    │                                                           │
│ Pełna dokumentacja:│                                                           │
│ F1 · mapa skrótów: │                                                           │
│ Ctrl+/             │                                                           │
├────────────────────┴───────────────────────────────────────────────────────────┤
│ Wersja 2.0 · kompilacja 2026.08-0000      Warunki licencji         Zamknij     │  ← .dn-modal-stopka
└────────────────────────────────────────────────────────────────────────────────┘
        272 px        ◄────────────────────  minmax(0, 1fr)  ───────────────────►
```

*Legenda:* `►` przy pozycji „1 Czym jest Danaco Console” oznacza rozdział bieżący (`aria-current="true"`),
odpowiednik kreski przy krawędzi `box-shadow: inset 2px 0 0 var(--dn-sygnal-500)` w rzeczywistym widoku.

Poniżej progu szerokości 900 px układ dwukolumnowy nie jest utrzymywany w tej postaci — zachowanie
opisuje rozdz. 20.

---

## 6. Spis rozdziałów — lewa kolumna

**Klasa kontenera:** `.iu-spis` (`<nav aria-label="Rozdziały instrukcji">`). Kolumna niesie nagłówek
`.iu-spis-glowa` z napisem **„Rozdziały”** krojem maszynowym w wersalikach, dziewięć przycisków pozycji
i notę zamykającą.

### 6.1 Wykaz dziewięciu pozycji

Wykaz poniżej podaje dosłowne brzmienie każdej pozycji spisu, w kolejności występowania w znaczniku.

| Nr | Nazwa pozycji (dosłownie) | `id` rozdziału w treści |
|---:|---|---|
| 1 | Czym jest Danaco Console | `iu-r1` |
| 2 | Pierwsze uruchomienie | `iu-r2` |
| 3 | Rama okna | `iu-r3` |
| 4 | Praca w środowisku | `iu-r4` |
| 5 | Projekty | *(zapowiedź, rozdz. 13 niniejszego dokumentu)* |
| 6 | Dostosowanie pasków | *(zapowiedź)* |
| 7 | Skróty klawiszowe | *(zapowiedź wraz z tabelą, rozdz. 13)* |
| 8 | Instrukcja instalacji | *(zapowiedź)* |
| 9 | Wersja i licencja | *(zapowiedź)* |

Każda pozycja jest przyciskiem `.iu-poz` o dwóch elementach wewnętrznych: `.iu-poz-nr` (numer
rozdziału, szerokość stała 16 px, krój maszynowy `--dn-ff-mono`) oraz `.iu-poz-nazwa` (nazwa, przycinana
wielokropkiem przy braku miejsca — `overflow: hidden; text-overflow: ellipsis; white-space: nowrap`).
Numer krojem maszynowym trzyma cyfry w jednej osi pionowej niezależnie od długości nazwy sąsiadującej —
lista czyta się jako spis numerowany, nie jako zbiór niezależnych przycisków o przypadkowej szerokości.

### 6.2 Stany pozycji spisu

| Stan | Wyzwalacz | Wygląd |
|---|---|---|
| Spoczynek | domyślny | tło przezroczyste, tekst `--dn-tekst-2`, bez obramowania |
| Wskazanie kursorem | `:hover` | tło `--dn-hover`, tekst `--dn-tekst` |
| Wciśnięcie | `:active` | dziedziczy zachowanie przycisku bazowego platformy — bez własnej reguły lokalnej |
| Ognisko klawiatury | `:focus-visible` | pierścień `--dn-fokus` zgodnie z regułą wspólną wszystkich kontrolek interaktywnych |
| Bieżący | `aria-current="true"` | tło `--dn-sygnal-tlo`, tekst `--dn-sygnal`, waga `--dn-fw-srednia`, kreska wewnętrzna `inset 2px 0 0 var(--dn-sygnal-500)` przy lewej krawędzi; numer rozdziału również w barwie `--dn-sygnal` |
| Nieaktywny | nie występuje | każda z dziewięciu pozycji jest zawsze klikalna — spis nie blokuje żadnego rozdziału |
| Ładowanie | nie występuje | treść rozdziału jest lokalna wobec dokumentu; zmiana rozdziału nie przechodzi przez sieć |
| Pusty | nie dotyczy | spis ma zawsze dokładnie dziewięć pozycji; liczba jest zamknięta, nie generowana z wykazu zewnętrznego |

**Dlaczego rozdział bieżący niesie kreskę oraz barwę, nie samą barwę.** Zasada powtórzona z ramy okna
([Rama okna aplikacji](rama-okna.md) rozdz. 3.2): wskazanie stanu nie może opierać się wyłącznie na kolorze — Operator
o obniżonej percepcji barw musi rozpoznać pozycję bieżącą po kształcie. Stąd kreska `inset 2px 0 0`
towarzysząca zawsze wypełnieniu sygnałowemu.

### 6.3 Nota zamykająca spis

Pod dziewiątą pozycją stoi akapit `.iu-spis-nota`:

> Pełna dokumentacja platformy otwiera się klawiszem **F1**, mapa skrótów — **Ctrl + /**.

Nota jest jedynym miejscem w spisie, w którym pojawia się odsyłacz do zasobów **poza** samą instrukcją —
przypomnienie, że instrukcja jest streszczeniem, a materiał źródłowy pełny stoi gdzie indziej (rozdz. 1).

---

## 7. Treść rozdziału — prawa kolumna, jednostki typograficzne

**Klasa kontenera:** `.iu-tresc`, wewnątrz — dokładnie jedna sekcja `.iu-rozdzial` widoczna naraz,
odpowiadająca pozycji wskazanej w spisie. Rozdział niesie stały zestaw jednostek typograficznych,
używanych w niezmiennej kolejności.

| Jednostka | Klasa | Rola |
|---|---|---|
| Nadtytuł | `.iu-nadtytul` | „Rozdział {n}” krojem maszynowym w wersalikach, barwa `--dn-tekst-3` |
| Tytuł | `.iu-tytul` (`<h3>`) | nazwa rozdziału, krój nagłówkowy, waga `--dn-fw-polgruba`, stopień `--dn-fs-xl` |
| Wiodące zdanie | `.iu-lid` | pierwszy akapit rozdziału, stopień `--dn-fs-base`, barwa `--dn-tekst-3`, do 90 znaków w wierszu |
| Podtytuł | `.iu-podtytul` (`<h4>`) | nagłówek podsekcji wewnątrz rozdziału |
| Akapit | `.iu-akapit` | treść bieżąca, stopień `--dn-fs-sm`, barwa `--dn-tekst-2`, wyróżnienia `<b>` w barwie `--dn-tekst` |
| Lista definicyjna | `.iu-def` (`<dl>`) | siatka dwukolumnowa 168 px + reszta — pojęcie i objaśnienie, użyta dla środowisk i dla elementów ramy |
| Lista wypunktowana | `.iu-lista` (`<ul>`) | wyliczenie z odstępem 4 px między pozycjami |
| Etapy przebiegu | `.iu-etapy` / `.iu-etap` | numerowany ciąg kroków z numerem liczonym przez CSS (`counter-increment`), użyty w rozdziale 2 |
| Makieta | `.iu-makieta` | schemat układu narysowany HTML-em, opisany w rozdz. 16 |
| Podpis makiety | `.iu-podpis` | zdanie „Na co patrzeć: …” pod każdą makietą, tłumaczące co widać na rysunku |
| Zapowiedź | `.iu-zapowiedz` / `.iu-zapowiedz-tytul` | blok streszczający rozdział jeszcze nierozwinięty w pełni (rozdziały 5–9) |
| Tabela | `.iu-tabela` | tabela skrótów klawiszowych, jedyna tabela HTML wewnątrz treści |
| Zapis klawisza | `.iu-klawisz` | fragment krojem maszynowym w obwódce, w postaci `Ctrl + F` |

Miara wiersza tekstu ciągłego (`.iu-akapit`, `.iu-lid`, `.iu-lista`) jest ograniczona do **90 znaków**
(`max-width: 90ch`) — wiersz dłuższy gubi powrót wzroku do początku wiersza następnego; makiety i tabela
mają miarę nieznacznie szerszą, do 92 znaków, bo niosą treść skanowaną punktowo, nie czytaną liniowo.

---

## 8. Rozdział 1 instrukcji — Czym jest Danaco Console

Rozdział otwiera się domyślnie przy każdym wywołaniu okna (`aria-current="true"` na pierwszej pozycji
spisu). Treść w pełnym, dosłownym brzmieniu źródła:

> Danaco Console jest platformą klasy AI Workspace OS: systemem operacyjnym dla sztucznej inteligencji,
> który łączy komunikację, zarządzanie wiedzą, tworzenie treści, projektowanie, automatyzacje procesów
> oraz rozwój oprogramowania w jednym oknie.

> Praca dzieli się na cztery poziomy: **środowisko** — najwyższy poziom organizacji, **moduł** —
> wyspecjalizowany obszar roboczy wewnątrz środowiska, **sesja** — pojedynczy tok pracy prowadzony we
> własnej karcie, oraz **projekt** — rama, w której moduły i sesje pracują nad wspólnym zamierzeniem.
> Obliczenia wykonuje serwer platformy; aplikacja na urządzeniu jest oknem do niego.

### 8.1 Cztery środowiska (`.iu-def`)

| Środowisko | Objaśnienie dosłowne |
|---|---|
| TalkIn | Środowisko rozmowy i pracy z wiedzą — komunikacja z modelami, badania, biblioteka materiałów. |
| WorkSpace | Środowisko pracy nad treścią i dokumentem — redakcja, projektowanie, zasoby projektu. |
| CodeStudio | Środowisko wytwarzania oprogramowania — kod, terminal, diagnostyka, aplikacje. |
| MultitaskingAI | Środowisko orkiestracji pracy zespołu ról; zamiast bocznej nawigacji modułów prowadzi je panel orkiestracji, a praca toczy się także wtedy, gdy żadne urządzenie Operatora nie jest podłączone. |

### 8.2 Moduły

> Platforma udostępnia piętnaście modułów: Assistant, Roundtable, Research, Library, Browser, Translate,
> Workspace, Studio, Design, Apps, Developer, Terminal, Diagnostics, Agents oraz Automations. Dostępność
> modułu w danym środowisku rozstrzyga macierz dostępności — środowisko pokazuje wyłącznie te, które do
> niego należą. Modułem wspólnym wszystkim środowiskom jest **Chat Window**: główne okno komunikacji
> Użytkownika z Wykonawcą, obecne w każdym module.

### 8.3 Sesje i projekty

> Sesja zawsze należy do środowiska — poza środowiskiem nie istnieje. Każda sesja działa jako **odrębny
> proces po stronie serwera**, z własnym katalogiem roboczym, i trwa niezależnie od tego, czy jej karta
> jest w tej chwili widoczna. Projekt gromadzi sesje wokół jednego zamierzenia: nadaje im katalog,
> instrukcję stałą, zespół agentów, automatyzacje i materiały wejściowe (rozdział 5 instrukcji).

Trzy jednostki niosące ten rozdział to: wiodące zdanie (`.iu-lid`), lista definicyjna dla czterech
środowisk (`.iu-def`) i dwa akapity zwykłe (`.iu-akapit`) dla modułów oraz sesji/projektów, każdy
poprzedzony podtytułem (`.iu-podtytul`) „Cztery środowiska”, „Moduły”, „Sesje i projekty”.

---

## 9. Rozdział 2 instrukcji — Pierwsze uruchomienie

Wiodące zdanie rozdziału:

> Droga od uruchomienia aplikacji do okna, w którym wykonuje się pracę, przebiega przez osiem etapów.
> Cztery pierwsze przechodzą samoczynnie albo wymagają jednej decyzji; cztery ostatnie są wyborem
> Operatora.

Osiem etapów jest wyliczeniem numerowanym przez CSS (`counter-increment: iu-etap` na każdym
`.iu-etap`) — numer nie jest wpisany ręcznie w treści, więc kolejność etapów jest kolejnością elementów
w dokumencie, nie liczbą przepisaną osobno.

```
┌───┐  Okno startowe
│ 1 │  Natywne okno aplikacji otwiera się i nawiązuje połączenie z serwerem
└───┘  platformy. Etap przejściowy — nie jest miejscem żadnej decyzji.
   │
┌───┐  Okno rejestracji i logowania
│ 2 │  Przy pierwszym uruchomieniu — rejestracja konta, przy kolejnych —
└───┘  logowanie. Przycisk „Pomiń” jest zawsze dostępny i klikalny.
   │
┌───┐  Przygotowanie środowiska pracy
│ 3 │  Klient odtwarza stan zapamiętany przy ostatnim zamknięciu: profil
└───┘  i uprawnienia, karty sesji trwające na serwerze, kanały modeli
       i konektory, magistralę kontekstu i pamięć projektów.
   │
┌───┐  Strona główna — Centrum dowodzenia
│ 4 │  Przedpokój przed strefą roboczą, w trzech strefach o malejącej wadze:
└───┘  karty czterech środowisk, kafle komponentów własnych, listwa ustawień.
   │
┌───┐  Wybór środowiska
│ 5 │  Kliknięcie jednej z czterech kart środowisk — jedyna droga wejścia
└───┘  do przestrzeni roboczej.
   │
┌───┐  Przedsionek środowiska
│ 6 │  Szyna sesji przy lewej krawędzi, kafle modułów i listwa działań
└───┘  środowiska. Przestrzeń robocza jeszcze się nie otwiera.
   │
┌───┐  Wybór modułu
│ 7 │  Kliknięcie kafla modułu w przedsionku. Przedsionek ustępuje
└───┘  przestrzeni roboczej: boczna nawigacja, pas kart sesji, Chat Window.
   │
┌───┐  Okno operacyjne
│ 8 │  Obszar roboczy wypełnia się zestawem okien właściwym modułowi.
└───┘  Tu wykonuje się pracę.
```

### 9.1 Treść dosłowna ośmiu etapów

| # | Nazwa etapu | Treść dosłowna |
|---:|---|---|
| 1 | Okno startowe | Natywne okno aplikacji otwiera się i nawiązuje połączenie z serwerem platformy. Etap przejściowy — nie jest miejscem żadnej decyzji. Jeżeli pakiet kliencki dostarczono bez wpisanego adresu serwera, okno prosi o jego podanie jednorazowo. |
| 2 | Okno rejestracji i logowania | Przy pierwszym uruchomieniu — rejestracja konta, przy kolejnych — logowanie. Przycisk „Pomiń” jest zawsze dostępny i klikalny; wymóg logowania Operator rozstrzyga sam w oknie konfiguracji. Urządzenie z ważnym tokenem przechodzi wprost do etapu trzeciego. |
| 3 | Przygotowanie środowiska pracy | Klient odtwarza stan zapamiętany przy ostatnim zamknięciu: profil i uprawnienia, karty sesji trwające na serwerze, kanały modeli i konektory, magistralę kontekstu i pamięć projektów. Widok pokazuje, który krok trwa. |
| 4 | Strona główna — Centrum dowodzenia | Przedpokój przed strefą roboczą, w trzech strefach o malejącej wadze: karty czterech środowisk, kafle komponentów własnych, listwa ustawień. Strona główna sama nie otwiera przestrzeni roboczej. |
| 5 | Wybór środowiska | Kliknięcie jednej z czterech kart środowisk. Jest to **jedyna droga wejścia** do przestrzeni roboczej — strefy druga i trzecia prowadzą do komponentów własnych oraz okien konfiguracji, nie do środowisk. |
| 6 | Przedsionek środowiska | Widok wejściowy środowiska: szyna sesji przy lewej krawędzi, kafle modułów i listwa działań środowiska. Przestrzeń robocza jeszcze się nie otwiera — żaden moduł nie jest wskazany, więc żaden nie zostaje otwarty. |
| 7 | Wybór modułu | Kliknięcie kafla modułu w przedsionku — albo pozycji w bocznej nawigacji, albo sekcji w panelu orkiestracji. Przedsionek ustępuje przestrzeni roboczej: pojawiają się boczna nawigacja, pas kart sesji i Chat Window; otwiera się pierwsza karta sesji. |
| 8 | Okno operacyjne | Obszar roboczy wypełnia się zestawem okien właściwym modułowi: Chat Window w lewej kolumnie na pełnej wysokości, Execution Loop Window w kolumnie sąsiedniej, obszar roboczy modułu w kolumnie dominującej, okna pomocnicze po prawej. Tu wykonuje się pracę. |

Zamykający akapit rozdziału:

> Z ósmego etapu prowadzą cztery drogi powrotne: zmiana modułu w tej samej karcie, otwarcie nowej karty
> przyciskiem „+”, powrót do Centrum dowodzenia (karta trwa w tle) oraz zamknięcie okna aplikacji —
> które również nie przerywa pracy sesji.

---

## 10. Rozdział 3 instrukcji — Rama okna

Wiodące zdanie:

> Rama okna to cztery pasy obecne w każdym oknie platformy, niezależnie od środowiska, modułu i widoku.
> Rama nie należy do żadnego środowiska ani modułu — należy do aplikacji, dlatego wygląda tak samo
> w oknie startowym, na stronie głównej, w przedsionku i w przestrzeni roboczej.

### 10.1 Makieta ramy okna

Rozdział niesie pierwszą z dwóch makiet osadzonych w instrukcji: `role="img"` z opisem dostępności
„Schemat ramy okna: szyna nawigacji przy lewej krawędzi, belka tytułowa, wstążka pozioma, treść widoku
i pasek stanu”.

```
┌──────┬──────────────────────────────────────────────────────────────────────────────┐
│ ☰   │ ◈ Danaco Console                  Raport końcowy — WorkSpace › Studio  ─ □ ✕ │
│praca │──────────────────────────────────────────────────────────────────────────────│
│ +    │ ⇤ │ ← → │ ⌂ ⟳ │ ⛶ 📋 │ [ wyszukiwanie ] │ karty sesji ▸ │ ⋯ │ Dostosuj       │
│ ⏱    │──────────────────────────────────────────────────────────────────────────────│
│ 🗀   │                                                                              │
│──────│   ┌─────────────────────────────────────────────────────────┐                │
│środow│   │ Treść widoku                                            │                │
│ ▣    │   │ przedsionek albo przestrzeń robocza — zależnie od etapu │                │
│ ▪▪   │   └─────────────────────────────────────────────────────────┘                │
│ ▣    │                                                                              │
│──────│                                                                              │
│szybki│                                                                              │
│ ◈    │                                                                              │
│ +    │                                                                              │
│──────│                                                                              │
│ ◉    │──────────────────────────────────────────────────────────────────────────────│
│      │ ◉ Operator │ WorkSpace · Studio │ Sesje: 7                  ◐ motyw │CPU│RAM │
└──────┴──────────────────────────────────────────────────────────────────────────────┘
```

Podpis makiety (`.iu-podpis`), dosłownie:

> Na co patrzeć: **szyna nawigacji** zaczyna się przy górnej krawędzi okna, nie pod belką — jej głowa
> tworzy kwadrat w narożniku, w którym stoi menu aplikacji. Szyna i belka niosą tę samą barwę i nie
> rozdziela ich kreska, więc rama czyta się jako jedna bryła otaczająca treść. Wstążka pozioma i belka
> tytułowa mieszczą się w kolumnie na prawo od szyny; pasek stanu biegnie przez całą szerokość okna.

### 10.2 Co gdzie stoi i dlaczego (`.iu-def`)

| Pas | Objaśnienie dosłowne |
|---|---|
| Belka tytułowa | Odpowiada na pytanie „w jakim programie jestem i jak steruję jego oknem?”. Niesie znak i nazwę produktu przy lewej krawędzi, tytuł bieżącego widoku pośrodku — w postaci „nazwa pracy w toku — miejsce, w którym się toczy” — oraz sterowanie oknem systemowym przy prawej. Zamknięcie okna nie przerywa sesji pracujących w tle. |
| Szyna nawigacji | Odpowiada na pytanie „gdzie chcę pracować?”. Stoi pionowo, ponieważ punktów wejścia jest kilkanaście, a pas poziomy zmieściłby je wyłącznie kosztem wyszukiwania i czynności. Trzy strefy: **praca** (nowa sesja, historia sesji, nowy projekt), **środowiska** (cztery środowiska, a pod środowiskiem rozwiniętym — jego moduły) i **szybki wybór** (komponenty własne i okna platformowe wraz z pozycją „Dodaj skrót”). Stopkę zamykają ustawienia, przełącznik motywu i profil użytkownika. |
| Wstążka pozioma | Odpowiada na pytanie „co zrobić z tym, co widzę?”. Układ jest przeglądarkowy: po lewej sterowanie i pole wyszukiwania, pośrodku pasmo kart otwartych sesji, po prawej czynności dotyczące pracy poza bieżącym widokiem, a na końcu — zawsze przyklejony do prawej krawędzi — przycisk „Dostosuj”. Nowa sesja, historia i nowy projekt na wstążce nie stoją: nie są czynnościami na treści, więc należą do strefy pracy szyny. |
| Pasek stanu | Odpowiada na pytanie „w jakim stanie jest praca i maszyna?”. Wyłącznie odczyt, bez czynności: Operator i konto, położenie (środowisko i moduł), liczba sesji czynnych, motyw bieżący oraz udział procesora i zajętość pamięci. Wartości nie są animowane — animowany licznik pokazywałby wielkości, których system nigdy nie zmierzył. |

Zamykający akapit rozdziału:

> Hierarchię pozycji szyny niesie **kształt, nie barwa**: środowisko jest większe i po rozwinięciu
> dostaje ramkę, moduł jest mniejszy i wcięty. Wypełnienie sygnałowe zarezerwowano dla modułu bieżącego —
> miejsca, w którym Operator faktycznie pracuje. Rozwinięte jest najwyżej jedno środowisko naraz. Każda
> pozycja ujawnia wskazaniem kursorem i ogniskiem klawiatury etykietę z nazwą i jednym zdaniem o przeznaczeniu.

Pełny opis ramy okna — wszystkie stany, wszystkie menu, wszystkie żetony — niesie
[Rama okna aplikacji](rama-okna.md); niniejszy rozdział jest jego streszczeniem do jednego ekranu
instrukcji.

---

## 11. Rozdział 4 instrukcji — Praca w środowisku

Wiodące zdanie:

> Wejście do środowiska prowadzi przez przedsionek, wybór modułu otwiera przestrzeń roboczą, a praca
> toczy się w kartach sesji, które trwają także wtedy, gdy nikt na nie nie patrzy.

### 11.1 Przedsionek środowiska

> Przedsionek jest widokiem wejściowym środowiska: przy lewej krawędzi stoi **szyna sesji** z wykazem
> sesji Operatora w tym środowisku, obok — płótno z nagłówkiem środowiska, kaflami modułów i listwą
> działań. Pasa kart sesji w przedsionku nie ma: karty należą do przestrzeni roboczej. Z listwy prowadzą
> te same drogi co ze strefy pracy szyny — „Nowy projekt” i „Pełna historia sesji” — z tą różnicą, że
> środowisko jest już znane z kontekstu, więc menu wyboru się nie pojawia.

### 11.2 Wybór modułu i układ przestrzeni roboczej — druga makieta

Druga makieta osadzona w instrukcji, `role="img"`, opis dostępności: „Schemat przestrzeni roboczej: szyna
nawigacji po lewej, wstążka z kartami sesji u góry, kolumny Chat Window, Execution Loop Window, obszar
roboczy modułu i panele pomocnicze”.

```
┌──────┬──────────────────────────────────────────────────────────────────────────────┐
│ ☰   │ ◈ Danaco Console                  Raport końcowy — WorkSpace › Studio  ─ □ ✕ │
│      │──────────────────────────────────────────────────────────────────────────────│
│      │ ⇤ │ ← → │ ⌂ ⟳ │ [ wyszukiwanie ] │ ▣ Raport końcowy │ ▫ Analiza │+│ Dostosuj │
│ ▣    │──────────────────────────────────────────────────────────────────────────────│
│ ▪▪▪  │  ┌───────────┬───────────┬─────────────────┬───────────────┐                 │
│ ▣    │  │Chat Window│ Execution │ Obszar roboczy  │ Panele        │                 │
│      │  │kolumna    │ Loop      │ modułu          │ pomocnicze    │                 │
│      │  │stała,     │ Window    │ ═══════════════ │ rozszerzenia  │                 │
│      │  │pełna wys. │ pętla     │ kolumna         │ boczne        │                 │
│ ◉    │  │Użytkownik │ wykonawcza│ dominująca      │ po prawej     │                 │
│      │  │↔ Wykonawca│ Koordyn.↔ │ edytor, lista   │               │                 │
│      │  │           │ Wykonawca │ albo monitor    │               │                 │
│      │  └───────────┴───────────┴─────────────────┴───────────────┘                 │
│      │──────────────────────────────────────────────────────────────────────────────│
│      │ ◉ Operator │ WorkSpace · Studio │ Sesje: 7                  ◐ motyw │CPU│RAM │
└──────┴──────────────────────────────────────────────────────────────────────────────┘
```

Podpis makiety, dosłownie:

> Na co patrzeć: kolumny sąsiadują **poziomo** i regulacji podlega wyłącznie ich szerokość — okna nie
> pływają i nie zachodzą na siebie. **Chat Window** zajmuje lewą kolumnę na pełnej wysokości w każdym
> module i przy zmianie modułu nie znika, tylko rekonfiguruje kontekst. Kolumna dominująca należy do
> modułu wybranego w szynie — stąd wskazanie sygnałowe. Zamknięcie Execution Loop Window albo panelu
> usuwa kolumnę z widoku, a zwolniona szerokość przypada kolumnom pozostałym; sama pętla wykonawcza
> biegnie dalej po stronie serwera.

### 11.3 Karty sesji

> Karty otwartych sesji zajmują pasmo wstążki poziomej, tak jak karty w oknie przeglądarki. Nowa karta
> powstaje przyciskiem „+” i prowadzi do wyboru modułu; przełączanie kart nie przerywa pracy żadnej
> z nich. Zmiana modułu w obrębie jednej karty przeładowuje wyłącznie tę kartę — pozostałe zostają
> nietknięte.

### 11.4 Trwałość pracy w tle

> Każda karta sesji jest odrębnym procesem na serwerze, dlatego praca trwa niezależnie od widoku
> i od połączenia:

| Sytuacja | Skutek dosłowny |
|---|---|
| Powrót do Centrum dowodzenia | karta znika z widoku, proces trwa dalej |
| Wybór innego środowiska | poprzednie środowisko i wszystkie jego karty pozostają czynne w tle i odtwarzają pełny stan przy powrocie |
| Zamknięcie okna aplikacji | kanał zostaje zamknięty, stan zapisany trwale, a procesy sesji trwają dalej po stronie serwera |
| Ponowne uruchomienie | karty odtwarzają się dokładnie w stanie sprzed rozłączenia: układ, historia i kontekst |

Zamykający akapit rozdziału:

> Ta trwałość jest podstawą funkcji globalnej Mobile — nadzoru nad procesami z urządzenia przenośnego —
> oraz pracy ciągłej środowiska MultitaskingAI, w którym role zespołu wykonują zadania również wtedy,
> gdy żadne urządzenie Operatora nie jest podłączone. Artefakt powstały w jednym module przenosi do
> innego **magistrala kontekstu**, bez eksportu i bez opuszczania karty sesji.

---

## 12. Scenariusze użycia

Cztery przebiegi obejmujące komplet zachowań okna, złożone z rozdziałów 2–10 w jeden ciąg zdarzeń.
Każdy przebieg jest niezależny od pozostałych — Operator może zamknąć nakładkę w dowolnym punkcie bez
skutku dla stanu aplikacji leżącej pod nią.

### 12.1 Pierwsze otwarcie w toku pracy

```
Operator                    Menu aplikacji              Nakładka instrukcji
   │                              │                              │
   │  otwiera menu aplikacji      │                              │
   ├─────────────────────────────►│                              │
   │                              │  najeżdża na „Pomoc”         │
   │                              │  → rozwija podmenu           │
   │  wskazuje „Instrukcja        │                              │
   │  użytkowania”                │                              │
   ├─────────────────────────────►│                              │
   │                              │  otwiera #okno-instrukcji    │
   │                              ├─────────────────────────────►│
   │                              │                              │  rozdział 1
   │                              │                              │  wybrany domyślnie
   │  czyta rozdział 1            │                              │
   │◄────────────────────────────────────────────────────────────┤
   │  naciska pozycję „3 Rama     │                              │
   │  okna” w spisie              │                              │
   ├────────────────────────────────────────────────────────────►│
   │                              │                              │  treść zamieniona,
   │                              │                              │  spis pozostaje
   │  naciska „Zamknij” w stopce  │                              │
   ├────────────────────────────────────────────────────────────►│
   │  wraca do karty sesji        │                              │
   │  dokładnie w stanie sprzed   │                              │
   │  otwarcia nakładki           │                              │
   │◄────────────────────────────────────────────────────────────┤
```

Karta sesji, z której Operator wywołał instrukcję, nie traci ogniska ani stanu przewijania — nakładka
stoi ponad nią przez cały czas otwarcia i nie modyfikuje jej zawartości.

### 12.2 Sprawdzenie skrótu klawiszowego w trakcie pracy nad dokumentem

Operator pracujący w module Studio nie pamięta kombinacji cofającej zmianę i sprawdza ją bez opuszczania
karty sesji:

1. Naciska hamburger szyny nawigacji → „Pomoc” → „Instrukcja użytkowania”.
2. Nakładka otwiera się z rozdziałem 1; Operator naciska pozycję „7 Skróty klawiszowe” w spisie.
3. Odnajduje w tabeli rozdz. 13.3 pozycję odpowiadającą szukanej czynności.
4. Naciska „Zamknij” — wraca do dokumentu w Studio z kombinacją zapamiętaną, bez utraty żadnej
   niezapisanej zmiany w edytorze pod spodem.

### 12.3 Przeciąganie i zmiana rozmiaru nakładki

Operator pracujący na dwóch monitorach przesuwa okno instrukcji na monitor pomocniczy, żeby czytać
rozdział 3 równolegle z pracą nad ramą okna w module Design:

1. Chwyta nagłówek `.dn-modal-naglowek` i przeciąga nakładkę na drugi ekran (rozdz. 4).
2. Chwyta uchwyt narożnika `--se` i powiększa nakładkę do maksymalnego wymiaru dopuszczalnego
   (`min(1180px, calc(100vw - var(--dn-od-8)))` szerokości, `min(88dvh, 820px)` wysokości).
3. Praca w oknie źródłowym (module Design) trwa przez cały czas równolegle — nakładka i widok pod nią są
   dwiema odrębnymi warstwami wizualnymi tego samego okna aplikacji, nie dwoma oknami systemowymi.

### 12.4 Zawężenie okna aplikacji poniżej 900 px

Operator zmniejsza okno aplikacji podczas pracy na wąskim ekranie:

1. Nakładka instrukcji jest już otwarta na rozdziale 4.
2. Szerokość okna aplikacji spada poniżej 900 px (rozdz. 20) — spis rozdziałów przechodzi z kolumny
   pionowej w pas poziomy nad treścią.
3. Wszystkie dziewięć pozycji spisu pozostaje dostępnych — żadna nie zostaje usunięta, zgodnie z regułą
   ustępowania wspólną całej ramie ([Rama okna aplikacji](rama-okna.md) rozdz. 17).
4. Treść rozdziału 4 pozostaje widoczna pod pasem spisu, bez przeładowania.

### 12.5 Operator instalujący produkt trafia na rozbieżność rozdz. 13.4

Operator, który dopiero zainstalował Danaco Console drogą opisaną w
[Oknie instalatora](okno-instalatora.md), otwiera instrukcję, żeby zrozumieć architekturę produktu
przed pierwszym uruchomieniem sesji roboczej:

1. Naciska hamburger szyny nawigacji → „Pomoc” → „Instrukcja użytkowania” (rozdz. 2).
2. Czyta rozdział 1 (domyślny) i przechodzi do rozdziału 8 „Instrukcja instalacji”, spodziewając się
   opisu tego, co właśnie wykonał.
3. Napotyka opis architektury serwerowej wieloużytkownikowej (rozdz. 13.4) — odmienny od
   jednostanowiskowej instalacji, którą faktycznie przeszedł, jeżeli instalacja przebiegła drogą
   opisaną w [Instalacja i konfiguracja](../INSTALACJA-I-KONFIGURACJA.md) (budowa ze źródeł, bez klasycznego pakietu
   instalacyjnego).
4. Rozbieżność, udokumentowana wprost w rozdz. 13.4 niniejszego opracowania, oznacza, że Operator w tym
   scenariuszu nie znajduje w treści instrukcji potwierdzenia własnego doświadczenia instalacyjnego —
   dokument nie ukrywa tego stanu rzeczy, lecz nazywa go i wskazuje dwa dokumenty, w których pełny
   obraz się składa: [Okno instalatora](okno-instalatora.md) dla stanu docelowego i
   [Instalacja i konfiguracja](../INSTALACJA-I-KONFIGURACJA.md) dla stanu faktycznego.

---

## 13. Rozdziały 5–9 — zapowiedzi i tabela skrótów klawiszowych

Nagłówek sekcji w treści nosi nadtytuł „Rozdziały 5–9” i tytuł „Dalsze rozdziały”, z wiodącym zdaniem:

> Rozdziały poniżej otwierają się z listy po lewej stronie okna. Każdy z nich jest samodzielny —
> kolejność czytania nie jest wymagana.

Pięć rozdziałów pozostałych stoi w treści jako bloki `.iu-zapowiedz` — krótkie streszczenie
zamiast pełnego rozwinięcia równego rozdziałom 1–4. Rozdział 7 wyróżnia się jedną różnicą: jego
zapowiedź poprzedza pełną tabelę HTML skrótów klawiszowych, jedyną tabelę w całej treści instrukcji.

### 13.1 Rozdział 5 — Projekty

> Rozdział opisuje zakładanie projektu — wskazanie środowiska, a następnie okno Workspace w trybie
> zakładania — oraz pięć grup ustaleń, które projekt zbiera: tożsamość wraz z katalogiem na dysku,
> instrukcję stałą i model wiodący, zespół agentów, automatyzacje i materiały wejściowe. Wyjaśnia,
> dlaczego katalog jest jedynym miejscem zapisu wytworów, dlaczego instrukcja stała obowiązuje każdą
> sesję i każdego agenta oraz w jaki sposób zestaw agentów i automatyzacji zmienia się później
> w modułach Agents i Automations.

Pełny projekt tego okna: [Okno nowego projektu](okno-nowego-projektu.md).

### 13.2 Rozdział 6 — Dostosowanie pasków

> Rozdział prowadzi przez okno „Dostosuj”, otwierane przyciskiem zamykającym wstążkę od prawej. Opisuje
> trzy zakładki — składniki pasków, karty sesji oraz wygląd i zachowanie — przenoszenie składników między
> wstążką poziomą, szyną pionową i zbiorem „poza paskami”, zasadę, że ten sam składnik nie stoi w dwóch
> miejscach naraz, oraz skład stały szyny, którego przenieść nie można. Ustawienie należy do konta
> i wędruje z nim na każde urządzenie.

Pełny projekt tego okna: [Rama okna aplikacji](rama-okna.md) rozdz. 7.

### 13.3 Rozdział 7 — Skróty klawiszowe

> Zestawienie najczęściej używanych kombinacji wraz z zasadami przypisywania własnych. Tabela poniżej
> jest jego częścią zasadniczą.

Tabela `.iu-tabela`, tytuł (`<caption>`) „Skróty klawiszowe — wybór”, kolumny **Kombinacja**, **Czynność**,
**Gdzie działa**. Wykaz pełny, w kolejności występowania w znaczniku:

| Kombinacja | Czynność | Gdzie działa |
|---|---|---|
| `Ctrl + K` | Paleta poleceń — kursor w polu wyszukiwania wstążki | całość aplikacji |
| `Ctrl + N` | Nowa sesja | całość aplikacji |
| `Ctrl + Shift + N` | Nowe okno aplikacji | całość aplikacji |
| `Ctrl + O` | Otwórz projekt | całość aplikacji |
| `Ctrl + S` | Zapisz sesję | przestrzeń robocza |
| `Ctrl + B` | Szyna sesji — pokaż albo ukryj | przedsionek i przestrzeń robocza |
| `Ctrl + Shift + B` | Szyna nawigacji — pokaż albo ukryj | całość aplikacji |
| `Ctrl + F` | Znajdź w widoku | bieżący widok |
| `Ctrl + H` | Znajdź i zamień | bieżący widok |
| `Ctrl + Shift + F` | Tryb skupienia | całość aplikacji |
| `F11` | Pełny ekran | całość aplikacji |
| `F5` | Odśwież widok | bieżący widok |
| `Ctrl + Tab` | Następna karta sesji | przestrzeń robocza |
| `Ctrl + Shift + Tab` | Poprzednia karta sesji | przestrzeń robocza |
| `Ctrl + Home` | Centrum dowodzenia | całość aplikacji |
| `Ctrl + E` | Przedsionek środowiska | przestrzeń robocza |
| `Ctrl + =` · `Ctrl + -` | Powiększ albo pomniejsz widok | całość aplikacji |
| `F1` | Dokumentacja platformy | całość aplikacji |
| `Ctrl + /` | Mapa skrótów klawiszowych | całość aplikacji |

Podpis pod tabelą (`.iu-podpis`), dosłownie:

> Skróty **Ctrl + F**, **F1** i **Ctrl + /** działają niezależnie od tego, czy odpowiadające im
> kontrolki są widoczne na wstążce, czy zwinięte pod rozwinięciem „Więcej”. Własne kombinacje przypisuje
> się w polach przypisania skrótu: pole pozostaje czynne również wtedy, gdy kombinacja jest już zajęta —
> Operator dostaje komunikat wskazujący czynność, która ją zajmuje, i wprowadza inną.

Zasady przypisywania własnych kombinacji, ze szczegółami stanów pola przypisania, niesie
[Rama okna aplikacji](rama-okna.md) rozdz. 16.

### 13.4 Rozdział 8 — Instrukcja instalacji

> Rozdział podaje wymagania urządzenia — komputer z systemem Windows albo Linux, urządzenie przenośne
> z Android, iOS albo iPadOS, stałe połączenie z serwerem, żadnego dodatkowego oprogramowania i żadnego
> układu GPU, ponieważ obliczenia wykonuje serwer — oraz przebieg instalacji pakietu klienckiego:
> instalator dla Windows, pakiet `deb` albo wariant `AppImage` dla Linuksa, wskazanie adresu serwera przy
> pierwszym uruchomieniu. Opisuje też aktualizacje: platforma aktualizuje się centralnie po stronie
> serwera, pakiet kliencki wymaga aktualizacji tylko wtedy, gdy zmienia się samo okno aplikacji, i nie
> usuwa konta, adresu ani tokenu urządzenia. Deinstalacja usuwa aplikację z urządzenia i nie usuwa żadnych
> danych — stan trwały platformy istnieje wyłącznie na serwerze; przed deinstalacją zaleca się odwołanie
> dostępu urządzenia w rejestrze urządzeń.

**Rozbieżność odnotowana wprost.** Treść rozdziału 8, wzięta dosłownie z prototypu, opisuje architekturę,
w której „obliczenia wykonuje serwer platformy” — model wieloużytkownikowy oprogramowania jako usługi, z centralną
aktualizacją. Dokument stanu faktycznego [Instalacja i konfiguracja](../INSTALACJA-I-KONFIGURACJA.md) (rozdz. 2) opisuje odmienny
stan wdrożenia: platformę jako aplikację jednego Operatora na jednej stacji Windows 11, z rdzeniem
lokalnym nasłuchującym na porcie 17870 w pętli zwrotnej, bez uwierzytelniania transportu w wersji v1.0.
**[DO DECYZJI OPERATORA]** — czy treść rozdziału 8 instrukcji ma opisywać architekturę docelową (serwer
wieloużytkownikowy, zgodną z niniejszym prototypem) czy stan faktyczny bieżącego wydania (stacja
jednoosobowa) — rozstrzygnięcie należy podjąć przed wdrożeniem tego rozdziału do produktu, inaczej
instrukcja obiecuje Operatorowi mechanizm, którego bieżące wydanie nie ma. Pełny projekt okna instalatora:
[Okno instalatora](okno-instalatora.md).

Ta sama rozbieżność, widziana z innej strony, jest przedmiotem [Okna instalatora](okno-instalatora.md)
rozdz. 1.3 „Stan faktyczny a stan docelowy”: ten dokument ustala, że Danaco Console w wersji v2.0 „nie
posiada klasycznego instalatora” ([Instalacja i konfiguracja](../INSTALACJA-I-KONFIGURACJA.md) rozdz. 5.1) i że konfiguracja
pakowania Windows (NSIS) przygotowana w `tauri.conf.json` pakuje — bez dalszych zmian — wyłącznie
powłokę, nie rdzeń (tamże, rozdz. 5.6). Rozdział 8 niniejszej instrukcji, opisujący „instalator dla
Windows, pakiet `deb` albo wariant `AppImage` dla Linuksa”, jest więc opisem tego samego stanu
docelowego, który [Okno instalatora](okno-instalatora.md) specyfikuje w pełni (sześć kroków, sześć
składników) — oba opracowania zgodnie traktują klasyczny instalator jako cel przyszły, nie stan
bieżący, choć żadne z nich osobno tego wprost nie rozstrzyga bez lektury drugiego. Rozstrzygnięcie
otwarte w obu miejscach jest tym samym rozstrzygnięciem: „Czy wprowadzić pełny instalator produktu
obejmujący rdzeń i czy go podpisywać” ([Instalacja i konfiguracja](../INSTALACJA-I-KONFIGURACJA.md) rozdz. 23.4, poz. 4).

### 13.5 Rozdział 9 — Wersja i licencja

> Rozdział podaje wydanie produktu i numer kompilacji pakietu klienckiego, warunki licencji Danaco
> Console wraz z zakresem uprawnień i ograniczeń oraz zestawienie składników obcych wykorzystanych
> w warstwie serwerowej i klienckiej — z nazwą, wersją, rodzajem licencji i rolą w produkcie. Wymienia
> także zależności zewnętrzne warunkujące działanie: konto u dostawcy modelu dla co najmniej jednego
> kanału, narzędzia i usługi właściwe wybranemu kanałowi, serwer platformy wraz z jego utrzymaniem oraz
> dostęp sieciowy urządzeń do serwera.

Wydanie i numer kompilacji cytowane w stopce okna (rozdz. 17) — „Wersja 2.0 · kompilacja 2026.08-0000” —
są tą samą parą wartości, którą rozdział 9 opisuje jako swoją treść: stopka jest skrótem stale widocznym,
rozdział 9 jego rozwinięciem.

---

## 14. Rozwinięcie rozdziału 5 — pięć grup ustaleń projektu

Zapowiedź rozdziału 5 (rozdz. 13.1) odsyła do okna zakładania projektu bez rozwijania treści pięciu
grup ustaleń, które to okno zbiera. Pełny wykaz, źródłowo z [Okna nowego projektu](okno-nowego-projektu.md) rozdz. 4 — dokument
przywoływany, nie powielany jako osobne ustalenie — dla ułatwienia lektury instrukcji w jednym miejscu:

| Grupa | Co ustala | Czego nie widać wprost |
|---|---|---|
| Tożsamość projektu | nazwa, katalog na dysku, jednostka pracy, termin, opis | Katalog jest **jedynym** miejscem zapisu wytworów; poza nim sesja nie zapisze niczego bez osobnej zgody Operatora. Nazwa pojawia się w szynie sesji, w historii sesji i w pasku stanu |
| Instrukcje i model | instrukcja stała, model wiodący, profil izolacji kontekstu | Instrukcja obowiązuje **każdą** sesję projektu i każdego agenta; sesja może ją uzupełnić, nie może jej znieść. Profil izolacji rozstrzyga, po co sesja może sięgnąć: katalog projektu, biblioteka środowiska, sieć |
| Zespół agentów | agenci przypisani do projektu | Agenci pracują na tej samej instrukcji stałej; zestaw zmienia się później w module Agents, także w trakcie trwania projektu |
| Automatyzacje | czynności wykonywane bez polecenia | Każdą można wstrzymać w module Automations; wstrzymanie nie usuwa jej z projektu |
| Materiały wejściowe | pliki wniesione na starcie | Trafiają do Project Library projektu i są widoczne dla wszystkich jego sesji, nie tylko dla pierwszej |

Okno kończy się trzema wyjściami: **założenie projektu wraz z otwarciem pierwszej sesji**, **założenie
bez sesji** (projekt czeka w Workspace środowiska i pojawia się w historii sesji dopiero po pierwszym
otwarciu) oraz **zapis jako szablon projektu** do ponownego użycia. Pełna specyfikacja okna, wraz ze
szkicami pól i komendami kontraktu, stoi w [Oknie nowego projektu](okno-nowego-projektu.md).

---

## 15. Rozwinięcie rozdziału 6 — cztery zakładki okna „Dostosuj”

Zapowiedź rozdziału 6 (rozdz. 13.2) nazywa trzy z czterech zakładek okna „Dostosuj”, bez rozwinięcia
zawartości którejkolwiek. Wykaz poniżej, źródłowo z [Ramy okna aplikacji](rama-okna.md) rozdz. 7.1–7.4, uzupełnia zapowiedź
o rozstrzygnięcia niesione przez każdą zakładkę.

**Zakładka „Składniki wstążki”** — trzy kolumny przenoszenia, każda z własnymi przyciskami:

| Kolumna | Kontrolki pod nią |
|---|---|
| Poza paskami | `Na wstążkę →` · `Na szynę →` |
| Wstążka pozioma | `← Odłóż` · `Na szynę →` · `↑` · `↓` |
| Szyna pionowa | `← Odłóż` · `Na wstążkę →` · `↑` · `↓` |

Cztery karty środowisk pracy oraz menu aplikacji, ustawienia i profil użytkownika należą do **stałego
składu szyny** — nie podlegają przenoszeniu ani usunięciu żadnym z powyższych przycisków.

**Zakładka „Karty sesji”:**

| Grupa | Rozstrzygnięcie |
|---|---|
| Zawartość karty | nazwa sesji i moduł · sama nazwa sesji · sama ikona modułu |
| Szerokość karty | stała 200 px · elastyczna |
| Nadmiar kart | przewijanie pasma · zwężanie kart · lista rozwijana |
| Oznaczenia karty | wskaźnik pracy w tle · przycisk zamknięcia · numer karty |

**Zakładka „Wygląd i zachowanie”:**

| Grupa | Rozstrzygnięcie |
|---|---|
| Postać kontrolek | same ikony · ikony z etykietami · napisy tylko w wybranych rodzinach |
| Forma wstążki | wstążka nakładana · wstążka na całą szerokość |
| Gęstość | zwarta 48 px · swobodna 56 px |
| Zachowanie | ukrycie w trybie skupienia · pole wyszukiwania na wstążce · jeden skład we wszystkich środowiskach |

**Zakładka „Personalizacja widoku”** (czwarta, nienazwana w zapowiedzi rozdziału 6, bo dodana do okna
„Dostosuj” później niż pozostałe trzy):

| Grupa | Rozstrzygnięcie |
|---|---|
| Kolorystyka pasów | atrament · grafit · powierzchnia · sygnał przygaszony |
| Kolor podświetlenia ikon | błękit sygnałowy · zieleń · bursztyn · neutralny |
| Rodzaj separacji pozycji | kreska między rodzinami · sam odstęp · tło rodziny |
| Akcja po kliknięciu pozycji | wykonaj od razu · podpowiedź, potem wykonaj · rozwiń menu wariantów |
| Odstępy i pozycjonowanie | pozwól wstawiać puste odstępy · wyrównaj rodziny do krawędzi (oba domyślnie włączone) |
| Grupowanie w jedną ikonę | zezwalaj na ikony zbiorcze · zwijaj rzadko używane pozycje szyny · pozwól usuwać pozycje szybkiego wyboru (wszystkie domyślnie włączone) |

Sześć grup w pełnym brzmieniu — z dosłownym opisem każdej opcji i notą łączącą grupę „Grupowanie w
jedną ikonę” z mechanizmem zakładki „Składniki wstążki” — stoi w
[Ramie okna aplikacji](rama-okna.md) rozdz. 7.4.

**Stopka okna „Dostosuj”:**

| Kontrolka | Czynność |
|---|---|
| `Przywróć układ domyślny` | Wraca skład wstążki i wszystkie rozstrzygnięcia do wartości fabrycznych |
| `Anuluj` | Zamyka okno bez zapisania zmian |
| `Zastosuj` | Zapisuje układ i zamyka okno; potwierdzenie niesie liczbę składników pozostałych na wstążce |

Ustawienie należy do konta i wędruje z nim na każde urządzenie — zapisane raz, obowiązuje niezależnie od
stacji, z której Operator się loguje.

---

## 16. Budowa makiet osadzonych w treści

Obie makiety instrukcji (rozdz. 10.1 i 11.2) są zbudowane z tych samych dwunastu klas lokalnych, bez
żadnego elementu graficznego spoza znaczników HTML i żetonów CSS. Zasada projektowa wprost z komentarza
źródła: *„Makieta jest rysunkiem poglądowym, nie odwzorowaniem widoku: pokazuje, co gdzie stoi, i celowo
nie niesie treści, żeby czytelnik patrzył na układ”*.

| Klasa | Rola w makiecie |
|---|---|
| `.iu-makieta` | kontener całości, obrys `--dn-obrys`, tło `--dn-powierzchnia-2`, zaokrąglenie `--dn-r-lg` |
| `.iu-makieta-rama` | siatka dwukolumnowa 56 px (szyna) + reszta, obrys `--dn-obrys-mocny` |
| `.iu-makieta-szyna` | kolumna szyny, tło `--dn-rama`, tekst `--dn-rama-tekst` |
| `.iu-makieta-kolumna` | kolumna prawa: belka, wstążka, płótno, pasek stanu, ułożone pionowo |
| `.iu-makieta-belka` | wiersz belki tytułowej wewnątrz makiety, tło `--dn-rama` |
| `.iu-makieta-wstazka` | wiersz wstążki poziomej, tło `--dn-panel`, kreska dolna `--dn-obrys-subtelny` |
| `.iu-makieta-plotno` | obszar treści makiety, tło `--dn-tlo`, wysokość minimalna 176 px |
| `.iu-makieta-stan` | wiersz paska stanu makiety, krój maszynowy, tło `--dn-panel` |
| `.iu-blok` | prostokąt z podpisem — jednostka budulcowa treści wewnątrz płótna; wariant `.iu-blok--wskazany` niesie wypełnienie sygnałowe i kreskę wewnętrzną dla elementu, na który pada uwaga podpisu |
| `.iu-znacznik` | kwadracik zastępujący ikonę pozycji szyny w makiecie; warianty `--srodowisko` (z obrysem) i `--modul` (mniejszy, wcięty) |
| `.iu-kreska-strefy` | krótka pozioma kreska 24 px rozdzielająca trzy strefy szyny wewnątrz makiety |
| `.iu-strefa-etyk` | podpis strefy szyny krojem maszynowym w wersalikach (widoczny wyłącznie w makiecie rozdz. 10.1 — makieta przestrzeni roboczej rozdz. 11.2 pomija etykiety stref, bo pokazuje wyłącznie środowisko już rozwinięte) |

**Dlaczego makieta jest rysunkiem, nie zrzutem ekranu.** Zrzut ekranu starzeje się z każdą zmianą
wizualną platformy i wymaga osobnego procesu aktualizacji poza redakcją tekstu. Makieta zbudowana
żetonami dziedziczy zmianę barwy i typografii automatycznie — zmienia się razem z resztą systemu
projektowego, bez interwencji redakcyjnej.

Dwanaście klas lokalnych wykazanych powyżej obsługuje dwie makiety różnej treści z tym samym szkieletem
strukturalnym (`.iu-makieta-rama` → `.iu-makieta-szyna` + `.iu-makieta-kolumna` →
`.iu-makieta-belka`/`.iu-makieta-wstazka`/`.iu-makieta-plotno`/`.iu-makieta-stan`) — różnica między
makietą ramy (rozdz. 10.1) a makietą przestrzeni roboczej (rozdz. 11.2) leży wyłącznie w zawartości
`.iu-makieta-plotno`: pierwsza pokazuje jeden blok treści ogólnej, druga cztery kolumny sąsiadujące.
Ten sam szkielet powtórzony dwukrotnie potwierdza zasadę jednego wzorca wizualnego dla obu makiet —
Operator czytający rozdział 4 rozpoznaje układ ramy z rozdziału 3 bez uczenia się nowej konwencji
graficznej.

---

## 17. Stopka okna

**Klasa:** `.dn-modal-stopka` (wspólna dla wszystkich okien nakładkowych) z zawartością własną instrukcji.

```
┌──────────────────────────────────────────────────────────────────────┐
│  Wersja 2.0 · kompilacja 2026.08-0000     Warunki licencji   Zamknij │
└──────────────────────────────────────────────────────────────────────┘
```

| Element | Klasa | Treść dosłowna | Zachowanie |
|---|---|---|---|
| Metryka wydania | `.iu-stopka-meta` | „Wersja 2.0 · kompilacja 2026.08-0000” | odczyt, `margin-right: auto` — przyklejona do lewej krawędzi stopki |
| Odnośnik licencji | `.iu-odnosnik` (`<a href="#warunki-licencji">`) | „Warunki licencji” | podkreślenie stałe, barwa `--dn-sygnal`, po wskazaniu kursorem `--dn-sygnal-mocny` |
| Przycisk zamknięcia | `.dn-btn.dn-btn--zarys` (`data-zamknij-modal`) | „Zamknij” | zamyka całą nakładkę, niezależnie od rozdziału otwartego |

Numer kompilacji „2026.08-0000” jest wartością przykładową prototypu (cztery zera końcowe) —
w wydaniu rzeczywistym pole niesie numer kompilacji faktycznej. Odnośnik „Warunki licencji” w prototypie
prowadzi do kotwicy `#warunki-licencji` wewnątrz tego samego pliku HTML; w opracowaniu docelowym pole
wskazuje [Licencję produktu](../LICENSE.md), nie kotwicę lokalną.

---

## 18. Wyszukiwanie w treści, tryb czytania, druk i eksport, wersjonowanie — zakres nieobjęty prototypem

Cztery mechanizmy bywają oczekiwane od okna pomocy tego rodzaju: wyszukiwanie pełnotekstowe w treści
instrukcji, tryb czytania oddzielony od trybu nawigacji, wydruk albo eksport treści oraz wersjonowanie
samej treści instrukcji (ślad zmian między wydaniami). Kontrola źródła — pełny przegląd znacznika
`instrukcja-uzytkowania.html`, w tym całego bloku modala `#okno-instrukcji` — nie ujawnia żadnego z tych
czterech mechanizmów: brak pola wyszukiwania wewnątrz nakładki, brak przełącznika trybu czytania, brak
przycisku druku albo eksportu, brak jakiegokolwiek znacznika wersji treści rozdziału.

| Mechanizm | Stan w prototypie | Najbliższy odpowiednik rzeczywisty |
|---|---|---|
| Wyszukiwanie w treści instrukcji | brak w znaczniku | menu „Znajdź w widoku” (`Ctrl+F`, [Rama okna aplikacji](rama-okna.md) rozdz. 10) przeszukuje **bieżący widok** — w tym otwartą nakładkę — ale nie przełącza rozdziałów za Operatora i nie przeszukuje ośmiu rozdziałów niewidocznych naraz |
| Tryb czytania | brak w znaczniku | dwukolumnowy układ (rozdz. 5) jest jedynym trybem; nie istnieje wariant „tylko treść” bez spisu |
| Druk i eksport | brak w znaczniku | brak odpowiednika; treść dostępna wyłącznie na ekranie, w oknie nakładkowym |
| Wersjonowanie treści | brak w znaczniku | metryka wydania w stopce (rozdz. 17) niesie wersję **pakietu**, nie wersję **treści instrukcji**; te dwie liczby nie muszą się zmieniać razem |

**[DO DECYZJI OPERATORA]** — czy okno instrukcji ma otrzymać własne pole wyszukiwania przeszukujące
wszystkie dziewięć rozdziałów naraz (z podświetleniem trafień i przeskokiem między nimi), własny
przełącznik trybu czytania, mechanizm druku albo eksportu treści do formatu PDF oraz osobny znacznik
wersji treści niezależny od wersji pakietu. Do chwili rozstrzygnięcia jedynym mechanizmem wyszukiwania
dostępnym Operatorowi wewnątrz tej nakładki pozostaje ogólne „Znajdź w widoku” wspólne całej platformie.

**Konsekwencja praktyczna braku wyszukiwania per-rozdział.** Ogólne „Znajdź w widoku” (`Ctrl+F`)
przeszukuje wyłącznie treść aktualnie widoczną w drzewie DOM — a treść niniejszej nakładki pokazuje
dokładnie jeden rozdział z dziewięciu naraz (rozdz. 7). Operator szukający frazy obecnej w rozdziale,
który akurat nie jest otwarty, nie znajdzie jej `Ctrl+F`, dopóki nie przełączy się na ten rozdział
ręcznie ze spisu — mechanizm ogólny nie zastępuje więc w pełni wyszukiwania przeszukującego wszystkie
rozdziały jednocześnie, o którym mowa w akapicie wyżej. Ta konkretna granica funkcjonalna jest
pierwszym praktycznym skutkiem braku mechanizmów wymienionych w tabeli powyżej, wartym odnotowania
osobno od samego faktu ich nieobecności w znaczniku.

---

## 19. Warstwy widoczności i dostępność

| Element | Warstwa | Uzasadnienie |
|---|---|---|
| Przycisk „Instrukcja użytkowania” w menu Pomoc | 2 | ujawniany dopiero po otwarciu menu aplikacji |
| Nakładka `#okno-instrukcji` w całości | 2 | ujawniana naciśnięciem pozycji menu; do tego czasu nieobecna w drzewie widocznym |
| Spis rozdziałów | 1 względem nakładki | widoczny bez dalszej interakcji, gdy nakładka jest otwarta |
| Treść rozdziału bieżącego | 1 względem nakładki | widoczna bez dalszej interakcji |
| Makiety osadzone (rozdz. 10.1, 11.2) | 1 względem rozdziału | widoczne bez dalszej interakcji wewnątrz rozdziałów 3 i 4 |

**Dostępność.** Wymagania poniżej obowiązują niezależnie od rozdziału otwartego w danej chwili.

| Wymaganie | Rozwiązanie |
|---|---|
| Etykieta okna | `aria-labelledby="iu-tytul"` na `<dialog>`, wskazujące nagłówek „Instrukcja użytkowania” |
| Nawigacja spisu | `<nav aria-label="Rozdziały instrukcji">`, pozycja bieżąca `aria-current="true"` |
| Makiety jako obraz opisowy | `role="img"` z `aria-label` opisującym układ — makieta nie jest odczytywana element po elemencie przez czytnik ekranu, tylko jako całość opisana jednym zdaniem |
| Zamknięcie | przycisk `.dn-modal-zamknij` z `aria-label="Zamknij okno"`, ikonowy krzyżyk `aria-hidden="true"` |
| Nawigacja klawiaturą | `Tab` przechodzi przez spis i treść w kolejności dokumentu; `Escape` zamyka całą nakładkę |
| Ognisko widoczne | pierścień `--dn-fokus`, wspólny dla wszystkich kontrolek interaktywnych platformy |
| Ograniczony ruch | `prefers-reduced-motion` wygasza przejścia przeciągania i zmiany rozmiaru ramy nakładki (rozdz. 4); zmiana tła pozycji spisu pozostaje, bez przejścia płynnego |

---

## 20. Zachowanie na punktach łamania

| Próg | Zachowanie |
|---|---|
| powyżej 900 px szerokości nakładki | układ dwukolumnowy pełny — spis 272 px, treść elastyczna (rozdz. 5) |
| poniżej 900 px szerokości nakładki | spis rozdziałów przechodzi w pas poziomy nad treścią; siatka `.iu-korpus` traci drugą kolumnę na rzecz układu jednokolumnowego pionowego |
| wysokość okna poniżej `88dvh` albo powyżej 820 px | wysokość nakładki nasyca się na `min(88dvh, 820px)` — ciało przewija się wewnątrz, nagłówek i stopka pozostają nieruchome (rozdz. 4) |
| szerokość okna aplikacji poniżej `calc(100vw - var(--dn-od-8))` mniejsza niż 1180 px | szerokość nakładki nasyca się na dostępną szerokość pomniejszoną o odstęp `--dn-od-8` z obu stron |

Reguła wspólna z ramą okna ([Rama okna aplikacji](rama-okna.md) rozdz. 17, „Reguła ustępowania”): żaden element nie znika
z systemu na progu węższym — spis rozdziałów zmienia wyłącznie orientację (z pionowej na poziomą), nie
znika żadna z dziewięciu pozycji, niezależnie od tego, jak wąskie stanie się okno aplikacji.

---

## 21. Skróty klawiszowe okna

| Kombinacja | Działanie | Zasięg | Kolizje |
|---|---|---|---|
| `Escape` | zamyka całą nakładkę, niezależnie od rozdziału otwartego | nakładka ma fokus | brak — `Escape` w innych oknach nakładkowych ma to samo znaczenie (rozdz. 4) |
| `Tab` / `Shift+Tab` | przechodzi między pozycjami spisu i elementami interaktywnymi treści | nakładka ma fokus | brak |
| `F1` | otwiera pełną dokumentację platformy — **nie** tę nakładkę | całość aplikacji | pozycja odrębna od tej, która otwiera niniejsze okno (rozdz. 2); Operator mylący `F1` z otwarciem instrukcji użytkowania trafia do innego zasobu |
| `Ctrl + /` | otwiera mapę skrótów klawiszowych — **nie** tę nakładkę | całość aplikacji | jak wyżej |
| strzałki, `Shift`+strzałki na nagłówku | przesuwanie i zmiana rozmiaru ramy nakładki | fokus na `.dn-modal-naglowek` | wspólne wszystkim oknom nakładkowym (rozdz. 4) |

Samo otwarcie okna **nie ma własnego skrótu klawiszowego** — jedyną drogą jest pozycja menu Pomoc
(rozdz. 2). Brak skrótu jest zamierzony: instrukcja jest zasobem sięganym rzadko, w przeciwieństwie do
`F1` (dokumentacja pełna) i `Ctrl+/` (mapa skrótów), które mają własne kombinacje jako narzędzia częstego
użytku, sięgane wielokrotnie w trakcie jednej sesji pracy.

---

## 22. Żetony `--dn-*` i `--iu-*` użyte w widoku

Wykaz żetonów pochodzących z `design/zasoby/zetony/zetony.css`, przywołanych w arkuszu lokalnym
`.iu-*`. Instrukcja **nie deklaruje żadnego żetonu własnego** — wyłącznie klasy CSS lokalne, wszystkie
wartości barw, wymiarów i typografii pochodzą z arkusza wspólnego.

| Rodzina | Żetony użyte | Zastosowanie w oknie |
|---|---|---|
| Kolor tekstu | `--dn-tekst`, `--dn-tekst-2`, `--dn-tekst-3` | trzy stopnie kontrastu tekstu: tytuły, treść bieżąca, metadane |
| Kolor sygnałowy | `--dn-sygnal`, `--dn-sygnal-tlo`, `--dn-sygnal-500`, `--dn-sygnal-obrys`, `--dn-sygnal-mocny` | rozdział bieżący w spisie, numer etapu w rozdz. 2, blok wskazany w makietach, odnośnik licencji |
| Powierzchnie | `--dn-powierzchnia`, `--dn-powierzchnia-2`, `--dn-panel`, `--dn-tlo` | tło spisu, tło makiety, tło etykiety pozycji, płótno makiety |
| Obrysy | `--dn-obrys`, `--dn-obrys-subtelny`, `--dn-obrys-mocny` | krawędzie makiety, kreska między spisem a treścią, kreska pod wstążką makiety |
| Rama (szyna/belka) | `--dn-rama`, `--dn-rama-tekst`, `--dn-rama-tekst-2`, `--dn-rama-hover`, `--dn-rama-obrys-mocny` | tło i tekst szyny oraz belki wewnątrz makiet |
| Interakcja | `--dn-hover`, `--dn-fokus` | wskazanie kursorem pozycji spisu, pierścień ogniska klawiatury |
| Typografia — krój | `--dn-ff-mono`, `--dn-ff-naglowek` | numer rozdziału i etykiety wersalikowe (mono), tytuły rozdziałów (nagłówkowy) |
| Typografia — stopień | `--dn-fs-xs`, `--dn-fs-sm`, `--dn-fs-base`, `--dn-fs-xl` | metadane, treść bieżąca, wiodące zdanie, tytuł rozdziału |
| Typografia — waga | `--dn-fw-normalna`, `--dn-fw-srednia`, `--dn-fw-polgruba` | tekst zwykły, wyróżnienia `<b>` i pozycja bieżąca, tytuły |
| Typografia — odstępy | `--dn-lh-luzny`, `--dn-lh-ciasny`, `--dn-ls-mono-wersaliki` | interlinia akapitów, interlinia bloków makiety, liternictwo etykiet wersalikowych |
| Zaokrąglenia | `--dn-r-sm`, `--dn-r-lg`, `--dn-r-xs`, `--dn-r-pill` | pozycja spisu, kontener makiety, blok makiety, znacznik etapu |
| Czas i easing | `--dn-czas-1`, `--dn-czas-2`, `--dn-ease` | przejście tła pozycji spisu, przejście etykiety podpowiedzi |

---

## 23. Komponenty `.dn-*` i `.iu-*` użyte w widoku

| Klasa | Pochodzenie | Rola |
|---|---|---|
| `.dn-modal`, `.dn-modal-naglowek`, `.dn-modal-tytul`, `.dn-modal-cialo`, `.dn-modal-stopka`, `.dn-modal-zamknij` | `design/zasoby/okna-modalne.css` | szkielet wspólny każdego okna nakładkowego platformy |
| `.dn-btn`, `.dn-btn--zarys`, `.dn-btn--sygnal` | `design/zasoby/css/komponenty.css` | przycisk „Zamknij” w stopce (`--zarys`) oraz przycisk „Otwórz instrukcję użytkowania” na tle za nakładką (`--sygnal`) |
| `.dn-btn-ikona` | `design/zasoby/css/komponenty.css` | przycisk zamknięcia w nagłówku nakładki (krzyżyk) |
| `.iu-okno` | lokalna, plik prototypu | wymiar własny nakładki (rozdz. 4) |
| `.iu-korpus`, `.iu-spis`, `.iu-spis-glowa`, `.iu-poz`, `.iu-poz-nr`, `.iu-poz-nazwa`, `.iu-spis-nota` | lokalne | układ dwukolumnowy i spis rozdziałów (rozdz. 5–6) |
| `.iu-tresc`, `.iu-rozdzial`, `.iu-nadtytul`, `.iu-tytul`, `.iu-podtytul`, `.iu-lid`, `.iu-akapit`, `.iu-lista`, `.iu-def` | lokalne | jednostki typograficzne treści (rozdz. 7) |
| `.iu-etapy`, `.iu-etap` | lokalne | ciąg ośmiu etapów rozdziału 2 (rozdz. 9) |
| `.iu-makieta`, `.iu-makieta-rama`, `.iu-makieta-szyna`, `.iu-makieta-kolumna`, `.iu-makieta-belka`, `.iu-makieta-wstazka`, `.iu-makieta-plotno`, `.iu-makieta-stan`, `.iu-blok`, `.iu-blok--wskazany`, `.iu-znacznik`, `.iu-znacznik--srodowisko`, `.iu-znacznik--modul`, `.iu-kreska-strefy`, `.iu-strefa-etyk`, `.iu-podpis` | lokalne | budowa dwóch makiet (rozdz. 13) |
| `.iu-zapowiedz`, `.iu-zapowiedz-tytul` | lokalne | pięć bloków zapowiedzi rozdziałów 5–9 (rozdz. 13) |
| `.iu-tabela`, `.iu-klawisz` | lokalne | tabela skrótów klawiszowych i zapis pojedynczej kombinacji (rozdz. 13.3) |
| `.iu-stopka-meta`, `.iu-odnosnik` | lokalne | metryka wydania i odnośnik licencji w stopce (rozdz. 17) |
| `.iu-tlo`, `.iu-tlo-karta`, `.iu-tlo-akcje`, `.iu-tlo-meta` | lokalne | karta wywołania na tle prototypu, poza samą nakładką — odpowiednik miejsca, z którego Operator otwiera okno w realnym widoku aplikacji |

Prefiks `iu-` (od „instrukcja użytkowania”) obejmuje wyłącznie klasy niewystępujące w arkuszach
wspólnych `komponenty.css`, `rama.css` i `okna-modalne.css` — zgodnie z regułą jednego źródła prawdy,
żadna z tych klas nie powiela nazwy ani przeznaczenia komponentu już istniejącego w tych trzech
arkuszach.

---

## 24. Etykiety interfejsu

Wykaz obowiązkowy zgodnie z [Standard redakcyjny i językowy](../STANDARD-REDAKCYJNY-I-JEZYKOWY.md) rozdz. 4.2 — element, dosłowne
brzmienie, miejsce wystąpienia. Wykaz obejmuje wyłącznie etykiety **własne** tego okna; etykiety menu
aplikacji poza sekcją Pomoc stoją w [Ramie okna aplikacji](rama-okna.md) rozdz. 9.

| Element | Brzmienie dosłowne | Miejsce wystąpienia |
|---|---|---|
| Tytuł nakładki | „Instrukcja użytkowania” | nagłówek `.dn-modal-tytul` |
| Nagłówek spisu | „Rozdziały” | `.iu-spis-glowa` |
| Pozycje spisu (9) | „Czym jest Danaco Console” · „Pierwsze uruchomienie” · „Rama okna” · „Praca w środowisku” · „Projekty” · „Dostosowanie pasków” · „Skróty klawiszowe” · „Instrukcja instalacji” · „Wersja i licencja” | `.iu-poz-nazwa`, rozdz. 6.1 |
| Nadtytuły rozdziałów (1–4) | „Rozdział 1” … „Rozdział 4” | `.iu-nadtytul` |
| Nadtytuł zbiorczy 5–9 | „Rozdziały 5–9” | `.iu-nadtytul` |
| Tytuł zbiorczy 5–9 | „Dalsze rozdziały” | `.iu-tytul` |
| Podtytuły rozdziału 1 | „Cztery środowiska” · „Moduły” · „Sesje i projekty” | `.iu-podtytul` |
| Podtytuł rozdziału 3 | „Co gdzie stoi i dlaczego” | `.iu-podtytul` |
| Tytuły zapowiedzi (5) | „5. Projekty” · „6. Dostosowanie pasków” · „7. Skróty klawiszowe” · „8. Instrukcja instalacji” · „9. Wersja i licencja” | `.iu-zapowiedz-tytul` |
| Tytuł tabeli skrótów | „Skróty klawiszowe — wybór” | `<caption>` w `.iu-tabela` |
| Nagłówki kolumn tabeli skrótów | „Kombinacja” · „Czynność” · „Gdzie działa” | `<th>` w `.iu-tabela` |
| Nota spisu | „Pełna dokumentacja platformy otwiera się klawiszem F1, mapa skrótów — Ctrl + /.” | `.iu-spis-nota` |
| Metryka wydania | „Wersja 2.0 · kompilacja 2026.08-0000” | `.iu-stopka-meta` |
| Odnośnik licencji | „Warunki licencji” | `.iu-odnosnik` |
| Przycisk zamknięcia stopki | „Zamknij” | `.dn-btn.dn-btn--zarys` |
| Przycisk zamknięcia nagłówka | ikona krzyżyka, `aria-label="Zamknij okno"` | `.dn-modal-zamknij` |
| Przycisk wywołania na tle prototypu | „Otwórz instrukcję użytkowania” | `.dn-btn.dn-btn--sygnal`, poza nakładką |
| Nota wywołania na tle prototypu | „nakładka otwiera się z rozdziałem 1 już wybranym” | `.iu-tlo-meta` |
| Nagłówek karty wywołania | „Instrukcja otwiera się nad pracą, nie zamiast niej” | `.iu-tytul` na tle prototypu, poza nakładką |
| Etykiety makiety ramy | „praca” · „środow.” · „szybki” | `.iu-strefa-etyk`, rozdz. 10.1 |
| Podpisy bloków makiety ramy | „Treść widoku” / „przedsionek albo przestrzeń robocza — zależnie od etapu” | `.iu-blok`, rozdz. 10.1 |
| Podpisy bloków makiety przestrzeni roboczej (4) | „Chat Window” · „Execution Loop Window” · „Obszar roboczy modułu” · „Panele pomocnicze” | `.iu-blok`, rozdz. 11.2 |

### 24.1 Etykiety sekcji Pomoc przywołane w rozdz. 2

Powtórzone tu dla kompletności wykazu — pełne brzmienie i skróty stoją w rozdz. 2 niniejszego dokumentu
oraz w [Ramie okna aplikacji](rama-okna.md) rozdz. 9: „Instrukcja użytkowania”, „Instrukcja instalacji”,
„Dokumentacja platformy”, „Skróty klawiszowe”, „Makiety i opracowania okien”, „Sprawdź dostępność
aktualizacji”, „Zgłoś obserwację”, „Wersja i licencja”, „O programie”.

---

## 25. Komunikaty

Wykaz obowiązkowy zgodnie z [Standard redakcyjny i językowy](../STANDARD-REDAKCYJNY-I-JEZYKOWY.md) rozdz. 4.2 — sytuacja, treść dosłowna,
rodzaj. Okno instrukcji jest widokiem statycznym (rozdz. 26): nie łączy się z rdzeniem, nie zgłasza
błędów walidacji, nie ma pól wprowadzania. Przegląd znacznika `#okno-instrukcji` w całości nie ujawnia
żadnego dymka powiadomienia, żadnego komunikatu ostrzeżenia ani błędu, żadnego stanu pustego wymagającego
własnego komunikatu — spis rozdziałów ma zawsze dokładnie dziewięć pozycji (rozdz. 6.2), więc stan pusty
nie występuje.

| Sytuacja | Treść dosłowna | Rodzaj |
|---|---|---|
| nie dotyczy — okno nie generuje komunikatów własnych | — | — |

Jedynym tekstem o charakterze zbliżonym do komunikatu jest nota spisu (rozdz. 24) i podpisy makiet
(„Na co patrzeć: …”, rozdz. 10.1, 11.2) — obie formy są treścią stałą, obecną niezależnie od działania
Operatora, nie reakcją na zdarzenie, więc nie spełniają definicji komunikatu w rozumieniu tego wykazu.
Jeżeli w toku rozstrzygnięcia luk rozdz. 18 (wyszukiwanie w treści) powstanie mechanizm zwracający
wynik pusty przy braku trafień, powstanie wówczas
pierwszy rzeczywisty komunikat tego okna — do tego czasu wykaz pozostaje pusty zasadnie, nie przez
przeoczenie.

---

## 26. Komendy kontraktu — nie dotyczy

Przegląd `budowa/shared/contract.json` (1119 komend w 68 obszarach) nie ujawnia obszaru właściwego
treści pomocy, dokumentacji ani wersji — nie istnieje rodzina komend `help.*`, `docs.*`, `manual.*` ani
`version.*`. Okno instrukcji jest widokiem w całości **statycznym wobec rdzenia platformy**: treść
dziewięciu rozdziałów jest częścią pakietu klienckiego, nie odpowiedzią serwera na zapytanie. Otwarcie
i zamknięcie okna, przełączanie rozdziałów oraz przewijanie treści są zdarzeniami lokalnymi warstwy
klienckiej, bez wywołania kontraktu.

Jedyny punkt styku z rdzeniem, pośredni i nieobowiązkowy, to pozycja „Sprawdź dostępność aktualizacji”
sekcji Pomoc (rozdz. 2) — czynność ta nie ma jednak własnej komendy w obszarach kontraktu udokumentowanych
na dzień redakcji niniejszego opracowania. **[DO DECYZJI OPERATORA]** — czy sprawdzenie dostępności
aktualizacji ma otrzymać własną komendę kontraktu (w przyszłym obszarze `application.*`), czy pozostać
mechanizmem czysto klienckim odpytującym zewnętrzny punkt dystrybucji pakietu.

Brak komend kontraktu odróżnia niniejsze okno od trzech pozostałych opracowań tego przejścia
redakcyjnego: [Okno instalatora](okno-instalatora.md) działa całkowicie poza połączeniem z rdzeniem
(przed jego nawiązaniem), natomiast [Okno historii sesji](okno-historii-sesji.md) i
[Okno nowego projektu](okno-nowego-projektu.md) wywołują odpowiednio obszary `session`/`history` oraz
`workspace`/`session` kontraktu w każdej czynności widocznej Operatorowi. Cztery opracowania tego
przejścia obejmują więc pełen zakres okien platformowych: od okna w całości statycznego, przez okno
przedstartowe, po dwa okna w pełni sprzężone z rdzeniem (rozdz. 27.2 zestawia wszystkie cztery).

**Weryfikacja wyczerpująca, nie próbkowa.** Twierdzenie o nieistnieniu obszarów `help`, `docs`,
`manual` i `version` nie opiera się na przeglądzie fragmentu kontraktu — pełny wykaz 68 nazw obszarów
z pola `obszary` dokumentu `contract.json`, odczytany programowo, nie zawiera żadnej z tych czterech
nazw ani żadnej nazwy pokrewnej semantycznie (`update`, `license`). Wykaz pełny sześćdziesięciu ośmiu
obszarów niesie [Architektura techniczna](../architektura/architektura.md) rozdz. 18; stoją w nim
`connection`, `home`, `environment`, `module`, `workspace`, `session`,
`window`, `message`, `history`, `agent`, `automation`, `studio`, `browser`, `research`, `translate`,
`mail`, `terminal`, `developer`, `design`, `mobile`, `settings`, `config`, `auth`, `account` — żaden
z nich nie niesie treści pomocy czy dokumentacji jako swojego przeznaczenia. To potwierdza, że
nieobecność nie jest luką w przeglądzie tego dokumentu, lecz faktem o całym kontrakcie: platforma w
wersji v2.0 nie modeluje treści pomocy jako bytu serwerowego w żadnym miejscu, nie tylko w tym oknie.

---

## 27. Relacja z `../INSTRUKCJA-UZYTKOWANIA.md` i kryteria odbioru

### 27.1 Tabela zgodności rozdziałów

| Rozdział okna (niniejszy dokument) | Rozdział źródłowy odpowiadający w `../INSTRUKCJA-UZYTKOWANIA.md` |
|---|---|
| 1 — Czym jest Danaco Console | rozdział wprowadzający produkt — środowiska, moduły, sesje, projekty |
| 2 — Pierwsze uruchomienie | rozdział opisujący przepływ pierwszego uruchomienia |
| 3 — Rama okna | rozdział opisujący ramę okna aplikacji |
| 4 — Praca w środowisku | rozdział opisujący przedsionek, przestrzeń roboczą i trwałość sesji |
| 5 — Projekty | rozdział opisujący zakładanie i strukturę projektu |
| 6 — Dostosowanie pasków | rozdział opisujący personalizację wstążki i szyny |
| 7 — Skróty klawiszowe | rozdział zestawiający skróty klawiszowe |
| 8 — Instrukcja instalacji | rozdział opisujący instalację, aktualizację i deinstalację — **z zastrzeżeniem rozbieżności odnotowanej w rozdz. 13.4** |
| 9 — Wersja i licencja | rozdział opisujący wydanie, licencję i składniki obce |

Każda zmiana rozdziału źródłowego w [Instrukcji użytkowania](../INSTRUKCJA-UZYTKOWANIA.md) pociąga
przegląd rozdziału odpowiadającego w treści niniejszego okna — reguła jednego źródła prawdy działa w
obu kierunkach: okno nie ustala treści pierwsze, ale musi pozostać z nią zgodne po każdej zmianie
źródła.

### 27.2 Cztery okna platformowe tego przejścia redakcyjnego — zestawienie

Rozdz. 26 wskazał już różnicę w sprzężeniu z rdzeniem między czterema opracowaniami tego przejścia
redakcyjnego. Tabela poniżej zbiera pełne zestawienie porównawcze — klasa okna, wywołanie, sprzężenie
z kontraktem, próg punktu łamania — w jednym miejscu, dla czytelnika przechodzącego między czterema
dokumentami jednej fali redakcyjnej.

| Cecha | [Okno instalatora](okno-instalatora.md) | Niniejsze okno (Instrukcja) | [Okno historii sesji](okno-historii-sesji.md) | [Okno nowego projektu](okno-nowego-projektu.md) |
|---|---|---|---|---|
| Klasa okna | okno wejściowe, bez ramy aplikacji | nakładka `.dn-modal` nad ramą pełną | nakładka `.dn-modal`, sekcja okna ustawień | okno pełne ramy, tryb modułu Workspace |
| Wywołanie | uruchomienie pliku pakietu, przed aplikacją | pozycja „Instrukcja użytkowania”, menu Pomoc (rozdz. 2) | trzy drogi — menu, szyna, przedsionek | menu szyny „Nowy projekt”, wybór środowiska |
| Sprzężenie z kontraktem | brak — przed nawiązaniem połączenia | brak — treść statyczna pakietu klienckiego (rozdz. 26) | pełne — obszary `session`, `history` | pełne — obszary `workspace`, `session` |
| Liczba komend eksponowanych | nie dotyczy | zero | sześć z dziewiętnastu komend obszaru `session`, plus `history.load` | pięć różnych obszaru `workspace` plus `session.create` |
| Próg punktu łamania własny | 1000 px (kolumna tożsamości) i 640 px | 900 px (spis rozdziałów) | 900 px (kolumny tabeli) | 1240 px (panel podsumowania) |
| Klasa dokumentu | Specyfikacja docelowa | Specyfikacja docelowa | Specyfikacja docelowa | Specyfikacja docelowa |

Cztery progi punktów łamania własnych — 1000 px, 900 px, 900 px, 1240 px — są czterema różnymi liczbami
nienazwanymi żadnym żetonem `--dn-bp-*` (zweryfikowane osobno w każdym z czterech dokumentów). Wzorzec
powtarzający się w całym przejściu redakcyjnym: okna platformowe zbioru ustalają progi lokalne wprost
w arkuszu stylów, niezależnie od skali `--dn-bp-w1`…`--dn-bp-w4` ramy aplikacji.

### 27.3 Kryteria odbioru

| Warunek sprawdzalny | Sposób sprawdzenia |
|---|---|
| Nakładka otwiera się wyłącznie z pozycji „Instrukcja użytkowania” sekcji Pomoc menu aplikacji | próba wywołania z każdego innego menu i skrótu kończy się niepowodzeniem |
| Nakładka otwiera się zawsze z rozdziałem 1 wybranym | stan `aria-current="true"` na pierwszej pozycji spisu przy każdym otwarciu |
| Spis rozdziałów niesie dokładnie dziewięć pozycji w kolejności z rozdz. 6.1 | porównanie z wykazem `.iu-poz` w znaczniku |
| Praca w toku pod nakładką nie zostaje przerwana | karta sesji aktywna przed otwarciem pozostaje aktywna po zamknięciu, bez zmiany stanu |
| Obie makiety (rozdz. 10.1, 11.2) renderują się bez elementu graficznego spoza HTML i żetonów | inspekcja znacznika — brak `<img>`, brak odwołania do pliku rastrowego |
| Tabela skrótów niesie dokładnie 19 wierszy z rozdz. 13.3 | porównanie z wykazem `<tr>` wewnątrz `.iu-tabela` |
| Stopka niesie metrykę wydania, odnośnik licencji i przycisk „Zamknij” w tej kolejności | inspekcja `.dn-modal-stopka` |
| `Escape` zamyka nakładkę z dowolnego rozdziału otwartego | próba na każdym z dziewięciu rozdziałów |
| Poniżej 900 px szerokości spis przechodzi w pas poziomy bez utraty żadnej z dziewięciu pozycji | próba na szerokości 899 px i 480 px |
| Żaden żeton barwy, wymiaru ani typografii nie jest wpisany na sztywno w arkuszu `.iu-*` | inspekcja arkusza lokalnego — wyłącznie `var(--dn-*)` |
| Rozbieżność rozdz. 13.4 (obliczenia serwera kontra stacja jednoosobowa) pozostaje jawnie oznaczona do chwili rozstrzygnięcia | obecność oznaczenia `[DO DECYZJI OPERATORA]` w treści rozdziału 8 wdrożonego produktu, dopóki Właściciel nie zdecyduje |
| Nakładka dziedziczy pełne zachowanie ramy nakładkowej (przeciąganie, zmiana rozmiaru, minimalizacja, maksymalizacja) opisane w rozdz. 4 | próba każdej z czterech czynności na nagłówku `.dn-modal-naglowek` i na ośmiu uchwytach krawędzi |
| Żadna z dziewięciu pozycji spisu nie prowadzi do treści pustej albo brakującej | otwarcie kolejno wszystkich dziewięciu pozycji i sprawdzenie obecności wiodącego zdania (`.iu-lid`) w każdej |
| Wykazy „Etykiety interfejsu” (rozdz. 24) i „Komunikaty” (rozdz. 25) pozostają zgodne ze znacznikiem źródłowym po każdej zmianie prototypu | porównanie wykazu z aktualną wersją `instrukcja-uzytkowania.html` przy każdej redakcji tego opracowania |
| Zero obszarów kontraktu właściwych treści pomocy (`help`, `docs`, `manual`, `version`) potwierdzone wobec pełnego wykazu 68 nazw | odczyt programowy pola `obszary` w `contract.json` (rozdz. 26) |
| Cztery progi punktów łamania własnych czterech okien tego przejścia (1000 px, 900 px, 900 px, 1240 px) pozostają udokumentowane jako nienazwane żadnym `--dn-bp-*` | porównanie rozdz. 27.2 z rozdziałami punktów łamania każdego z czterech dokumentów |
| Rozdział 9 treści („Wersja i licencja”) i stopka okna (rozdz. 17) niosą tę samą parę wartości wydania i kompilacji | porównanie treści dosłownej rozdz. 13.5 z `.iu-stopka-meta` |

---

*Koniec dokumentu. Okno instrukcji użytkowania — Specyfikacja docelowa, wersja 2.0, 2026-08-20.*

---
*Danaco Console — Platforma AI Workspace OS · v2.0 · status Deweloperski*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz.*
*Warunki korzystania: [Licencja produktu](../LICENSE.md). Kontakt: support@danaco-group.pl*
