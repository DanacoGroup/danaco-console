# Danaco Console — Plansza portfolio 07: Typografia w działaniu

| | |
|---|---|
| **Produkt** | Danaco Console — AI Operating Environment |
| **Warstwa** | Portfolio Design Identity · plansza tematyczna 07 |
| **Wersja pakietu wizualnego** | v2.0 |
| **Status** | Deweloperski |
| **Data opracowania** | 2026-08-14 |
| **Odbiorcy** | zespół projektowy, zespół wdrożeniowy, recenzja właścicielska |
| **Zakres** | trzy role krojów · dziewięciostopniowa skala · zestawienia typograficzne · liczby tabelaryczne · wersaliki · diakrytyki polskie · gęstość |
| **Plik towarzyszący** | `07-typografia-w-dzialaniu.html` |

---

## 0. Po co ta plansza

Dokumentacja techniczna wymienia trzy kroje i dziewięć stopni. To jeszcze nie jest typografia —
to inwentarz. Plansza 07 pokazuje **ten sam inwentarz w miejscu pracy**: tam, gdzie stopień `xs`
naprawdę występuje (nagłówek kolumny tabeli), tam, gdzie `3xl` naprawdę występuje (tytuł karty
środowiska na Stronie głównej), tam, gdzie mono naprawdę występuje (identyfikator, ścieżka
repozytorium, znacznik czasu, wynik powłoki).

Zasada nadrzędna zespołu portfolio obowiązuje bez wyjątku: **każda próbka na planszy jest cytatem
z platformy**. Nie ma tu ani jednego „przykładowego nagłówka", ani jednej wymyślonej etykiety.
Gdy potrzebna jest treść wypełniająca, pochodzi ona z domeny produktu (nazwy modułów, nazwy okien
operacyjnych, ścieżki repozytoriów widoczne w prototypie Terminal, kod parowania widoczny
w prototypie Ustawień) i jest opisana jako przykładowa.

---

## 1. Trzy role, trzy kroje

Źródło wiążące: kontrakt systemu projektowego („Typografia — trzy role, trzy kroje"), kierunek systemu projektowego,
`zasoby/zetony/fonty.css`, `zasoby/zetony/zetony.css` sekcja 3.

| Rola | Krój | Żeton | Wagi w pakiecie | Gdzie występuje w platformie |
|---|---|---|---|---|
| Wejście do środowiska i tożsamość marki | **Space Grotesk** | `--dn-ff-naglowek` | 500 · 600 · 700 | logotyp `DANACO`, tytuły kart środowisk, nagłówki okien modalnych, tytuły sekcji Strony głównej |
| Praca | **IBM Plex Sans** | `--dn-ff-bazowa` | 400 · 500 · 600 · 700 | etykiety pól, treść wpisów rozmowy, opisy modułów, pozycje bocznej nawigacji, przyciski, tabele |
| Maszyna | **IBM Plex Mono** | `--dn-ff-mono` | 400 · 500 · 600 | identyfikatory, ścieżki repozytoriów, wynik powłoki, nazwy żetonów, znaczniki czasu, metryki |

### 1.1. Zdanie kierunkowe

> **Space Grotesk = wejście do środowiska i tożsamość marki; IBM Plex Sans = praca;
> IBM Plex Mono = maszyna.**

Para krojów w logotypie opowiada produkt: `DANACO` (Space Grotesk 700) to marka,
`CONSOLE` (IBM Plex Mono 500, rozstrzelone wersaliki) to maszyna. Ten sam podział przenosi się
na cały interfejs — użytkownik nie musi go znać, żeby z niego korzystać, ale rozpoznaje go
podświadomie: **to, co wygląda na mono, można skopiować i wkleić**.

### 1.2. Reguły przydziału (rozstrzygnięcia projektowe)

1. Space Grotesk nie występuje w tekście ciągłym. Wyłącznie tytuły, nagłówki i logotyp.
2. Space Grotesk nie występuje poniżej stopnia `lg` (16 px) — poniżej tego progu jego charakter
   ginie, a przewaga czytelności Plex Sans rośnie. Wyjątek: logotyp paska górnego, gdzie krój
   pełni funkcję znaku, nie tekstu.
3. Plex Mono nie występuje w zdaniach. Wyłącznie wartości, identyfikatory, ścieżki, kod, wynik.
4. Wartość, którą można skopiować, jest zapisana krojem mono — to sygnał afordancji.
5. Plex Sans jest krojem domyślnym całego dokumentu; pozostałe dwa są przydzielane świadomie.

---

## 2. Pełna skala — dziewięć stopni w kontekście

Źródło: `zetony.css` sekcja 3, kontrakt systemu projektowego. Baza **13 px** wynika z pokrętła
`GESTOSC_WIZUALNA = 8/10` — kokpit ma być zwarty.

| Żeton | Stopień | Krój typowy | Waga | Interlinia | Realny kontekst w platformie |
|---|---:|---|---:|---|---|
| `--dn-fs-xs` | 11 px | Plex Sans / Plex Mono | 500–600 | 1.25 | etykieta wersalikowa `ŚRODOWISKO`, nagłówek kolumny tabeli, licznik filtru `7 / 7` |
| `--dn-fs-sm` | 12 px | Plex Sans / Plex Mono | 400–500 | 1.45 | metadane wpisu, podtytuł okna, znacznik czasu, opis pola |
| `--dn-fs-base` | 13 px | Plex Sans | 400–500 | 1.45 | tekst interfejsu: przyciski, pozycje bocznej nawigacji, komórki tabeli |
| `--dn-fs-md` | 14 px | Plex Sans | 400 | 1.6 | treść wpisu rozmowy w Chat Window, akapit opisowy |
| `--dn-fs-lg` | 16 px | Space Grotesk / Plex Sans | 600 | 1.25 | nagłówek panelu (`Project Dashboard`, `Logs Viewer`) |
| `--dn-fs-xl` | 20 px | Space Grotesk | 600 | 1.25 | nagłówek okna modalnego, nagłówek okna operacyjnego |
| `--dn-fs-2xl` | 24 px | Space Grotesk | 600 | 1.25 | tytuł sekcji Strony głównej („Strefa 1 — Środowiska") |
| `--dn-fs-3xl` | 30 px | Space Grotesk | 600–700 | 1.25 | tytuł karty środowiska: `TalkIn`, `WorkSpace`, `CodeStudio`, `MultitaskingAI` |
| `--dn-fs-display` | 40 px | Space Grotesk | 700 | 1.25 | stopień ekspozycyjny — okno startowe, Always On Display, materiały marki |

### 2.1. Dlaczego dziewięć, a nie dwanaście

Skala jest krótka celowo. Kokpit ma **cztery poziomy hierarchii** (środowisko → moduł → okno →
wiersz danych) i każdemu odpowiada jeden stopień wiodący. Pozostałe pięć stopni obsługuje
metadane, etykiety i ekspozycję. Dwunastostopniowa skala wymusiłaby decyzje, których projekt nie
potrzebuje, i rozmyła granice między poziomami.

### 2.2. Skoki w skali

11 → 12 → 13 → 14 to skoki jednopunktowe: te cztery stopnie żyją obok siebie w jednym wierszu
tabeli albo w jednym wpisie rozmowy, więc różnica musi być subtelna. Powyżej 14 skoki rosną
(16 → 20 → 24 → 30 → 40), bo te stopnie nigdy nie sąsiadują — dzieli je zawsze warstwa układu.

---

## 3. Zestawienia typograficzne

Zestawienie to najmniejsza jednostka projektowa typografii: nie pojedynczy stopień, lecz **para
albo trójka stopni, które zawsze występują razem**.

| Zestawienie | Skład | Występowanie |
|---|---|---|
| Tytuł + podtytuł | `3xl` Space Grotesk 600 + `sm` Plex Sans 400 | karta środowiska (tytuł + motto) |
| Etykieta + wartość | `xs` Plex Sans 600 wersaliki 0,08em + `base` Plex Mono 500 | tabela danych `.dn-dane`, panel prowenancji |
| Nagłówek panelu + treść | `lg` Space Grotesk 600 + `base` Plex Sans 400 | każde okno operacyjne modułu |
| Wpis rozmowy | `sm` nadawca 600 + plakietka roli `xs` + `sm` mono godzina + `md` treść 400 | Chat Window — `.dn-wpis` |
| Wiersz tabeli | `xs` nagłówek kolumny wersaliki + `base` komórka + `sm` mono liczba | Queue Manager, Execution Monitor |
| Blok kodu | `sm` Plex Mono 400, numery wierszy `sm` mono tabular | Studio Editor, Code Editor |
| Pasek promptu | grot `sm` mono + `md` Plex Sans 400 | `.dn-prompt` — wspólny wszystkim modułom |

### 3.1. Reguła kontrastu w zestawieniu

W każdej parze zmienia się **co najmniej dwie cechy** (stopień + waga, stopień + krój, waga +
barwa tekstu). Zmiana samego stopnia daje hierarchię zbyt słabą przy gęstości zwartej; zmiana
samej barwy łamie zasadę „stan nigdy samym kolorem" przeniesioną na poziom typografii.

---

## 4. Liczby — `tabular-nums`

Źródło: kontrakt systemu projektowego („Liczby: `tabular-nums` w Plex Mono").

Liczba w kolumnie musi być **porównywalna wzrokiem**. Proporcjonalne cyfry sprawiają, że kolumna
identyfikatorów albo rozmiarów faluje, a operator traci możliwość odczytania rzędu wielkości bez
czytania cyfr. `font-variant-numeric: tabular-nums` wyrównuje szerokość każdej cyfry.

### 4.1. Formaty danych platformy

| Format | Przykład (z prototypów) | Krój | Wyrównanie |
|---|---|---|---|
| Kod parowania urządzenia | `DNC-4718-2K30` | mono 500 | do lewej |
| Znacznik czasu przebiegu | `2026-08-14 07:00` | mono 400 | do lewej |
| Czas trwania | `2,1 s` | mono 400 | do prawej |
| Rozmiar procesu | `210 MB` | mono 400 | do prawej |
| Licznik filtru | `7 / 7` | mono 500 | do prawej |
| Postęp | `72 %` | mono 500 | do prawej |
| Wersja pakietu | `v2.0` | mono 500 | do lewej |
| Suma kontrolna | `SHA-256` | mono 400 | do lewej |

### 4.2. Zasady

1. Każda liczba w tabeli, kolejce, monitorze i panelu metryk ma `tabular-nums`.
2. Liczby wyrównujemy do prawej wtedy i tylko wtedy, gdy są porównywalne między wierszami.
3. Identyfikatory (kody, wersje, znaczniki czasu w formacie ISO) wyrównujemy do lewej — czyta
   się je od początku, nie porównuje wielkością.
4. Separator dziesiętny — przecinek (polska konwencja). Separator tysięcy — spacja nierozdzielająca.
5. Jednostka jest częścią wartości i idzie tym samym krojem, oddzielona spacją nierozdzielającą.

---

## 5. Wersaliki

| Zastosowanie | Krój | Odstęp liter | Żeton |
|---|---|---|---|
| Etykieta sekcji, nagłówek kolumny | IBM Plex Sans 600 | 0,08em | `--dn-ls-wersaliki` |
| Etykieta techniczna, plakietka mono, `CONSOLE` w logotypie | IBM Plex Mono 500 | 0,14em | `--dn-ls-mono-wersaliki` |

Mono wymaga większego odstępu, bo jego wersaliki są już szerokie i gęste — bez rozstrzelenia
zlewają się w blok. Sans przy 0,08em zachowuje rytm słowa.

### 5.1. Gdzie wolno

- etykiety sekcji i grup pól,
- nagłówki kolumn tabel,
- drobne plakietki techniczne (nazwa żetonu, nazwa roli),
- `CONSOLE` w logotypie.

### 5.2. Gdzie nie wolno

- nazwy własne modułów i okien (`Studio Editor`, `Workflow Builder`) — dokumentacja zapisuje je
  w oryginalnej pisowni i zakaz parafrazy obejmuje także zmianę wielkości liter,
- nazwy środowisk (`TalkIn`, `MultitaskingAI`) — wielbłądzia pisownia jest częścią nazwy,
- treść wpisów rozmowy, opisy, komunikaty, teksty przycisków,
- wszystko powyżej stopnia `lg`.

### 5.3. Odstęp zerowy jako kontrprzykład

Ta sama etykieta bez rozstrzelenia czyta się jako skrót, nie jako etykieta. Rozstrzelenie
komunikuje: „to nie jest treść, to jest opis pola treści".

---

## 6. Diakrytyki polskie

Podzbiory `latin` + `latin-ext` są w pakiecie dla wszystkich trzech krojów i wszystkich wag
(`fonty.css` — zakresy `U+0100-024F` obejmują ą ć ę ł ń ó ś ź ż Ą Ć Ę Ł Ń Ó Ś Ź Ż).

Zdanie testowe kanoniczne: **ŁÓDŹ ŻÓŁĆ GĘŚLĄ JAŹŃ**.

Zestaw kontrolny małych liter: `ą ć ę ł ń ó ś ź ż`
Zestaw kontrolny wielkich liter: `Ą Ć Ę Ł Ń Ó Ś Ź Ż`

### 6.1. Punkty ryzyka

| Znak | Ryzyko | Rozstrzygnięcie |
|---|---|---|
| `Ł` / `ł` | przekreślenie kolidujące z sąsiadem w wersalikach mono | odstęp 0,14em rozwiązuje kolizję |
| `ą` / `ę` | ogonek wchodzący w interlinię przy 1.25 | interlinia ciasna zarezerwowana dla nagłówków w krótkich stopniach ≥ 16 px |
| `Ź` / `Ż` | kropka i akut mylone przy 11 px | stopień `xs` wyłącznie z wagą ≥ 500 |
| `Ń` | akut nad wąskim `N` w mono | akceptowalne we wszystkich stopniach pakietu |

### 6.2. Zasada wagi przy małych stopniach

Przy stopniach `xs` (11 px) i `sm` (12 px) waga 400 w Plex Mono osłabia znaki diakrytyczne.
Dlatego etykiety wersalikowe mono chodzą wagą 500, a nie 400.

---

## 7. Gęstość a typografia

Pokrętło `GESTOSC_WIZUALNA = 8/10` przekłada się na dwie konfiguracje typograficzne:

| Cecha | Zwarta (domyślna) | Przestronna |
|---|---|---|
| Baza tekstu interfejsu | 13 px | 14 px |
| Interlinia bazowa | 1.45 | 1.50 |
| Treść wpisu rozmowy | 14 px | 15 px odpowiednik przez skalowanie |
| Wysokość wiersza tabeli | 36 px | 44 px |
| Zastosowanie | biurko, pełny kokpit (`w3`, `w4`) | Mobile, `pointer: coarse`, dostępność |

Przestronna nie jest „wersją dla słabszego wzroku" — jest wariantem dla wskaźnika grubego
(dotyk), gdzie cele muszą urosnąć, a typografia idzie za nimi, żeby proporcja została zachowana.

### 7.1. Co się nie zmienia

Skala stopni ekspozycyjnych (`xl`, `2xl`, `3xl`, `display`) jest identyczna w obu gęstościach.
Zmienia się wyłącznie warstwa robocza: baza, interlinia i wysokość wiersza. Tytuł karty
środowiska ma 30 px zawsze — bo to element tożsamości, nie element pracy.

---

## 8. Decyzje projektowe planszy

1. **Zero abstrakcyjnych próbek.** Każdy okaz typograficzny to fragment realnego interfejsu —
   nagłówek panelu z nazwą okna operacyjnego z inwentarza, wpis rozmowy z klasami `.dn-wpis`,
   wiersz tabeli z licznikiem, ścieżka repozytorium przepisana z prototypu Terminal.
2. **Skala pokazana w kontekście, nie jako lista.** Dziewięć stopni występuje w dziewięciu
   miniaturach interfejsu; obok każdej podana jest nazwa żetonu, waga i interlinia — jako
   metadane okazu, nie jako jego treść.
3. **Porównania zawsze parami w jednym polu.** Tabelaryczne vs proporcjonalne cyfry, wersaliki
   z odstępem vs bez, gęstość zwarta vs przestronna — zawsze dwie kolumny w jednym kadrze,
   nigdy dwa oddzielne bloki, bo wtedy różnica znika.
4. **Interaktywny okaz zamiast statycznej specyfikacji.** Pole tekstowe z trzema suwakami
   (stopień, interlinia, odstęp liter) renderuje wpisany tekst równolegle we wszystkich trzech
   krojach. Domyślna treść pola to zdanie testowe diakrytyków.
5. **Przełącznik gęstości przebudowuje całą planszę**, nie tylko jedną demonstrację — bo gęstość
   jest właściwością systemu, nie właściwością komponentu.
6. **Ruch 3/10.** Jedyny ruch ciągły to tętno kropki sygnału przy okazie „praca w tle".
   Wszystkie pozostałe zmiany to mikroprzejścia 100–220 ms wywołane działaniem użytkownika.
7. **Zero blokad.** Żadna kontrolka planszy nie jest wyłączona. Suwaki mają pełny
   zakres; wyjście poza zakres skali systemowej jest komunikowane plakietką, nie blokadą.

---

## 9. Interakcje odwzorowane na planszy

| Interakcja | Mechanizm | Efekt |
|---|---|---|
| Przełączenie motywu | `data-przelacz-motyw` + `wspolne.js` | oba motywy równoprawne, pasek górny atramentowy w obu |
| Podświetlenie sekcji | `data-spis` + `IntersectionObserver` (`prototyp.js`) | aktywna pozycja spisu |
| Przełączenie roli kroju | `data-przelacz` / `data-grupa` | wymiana widoku `data-widok` |
| Wybór stopnia skali | `data-lista-wyboru` + `data-pozycja` + `data-pokaz` | podmiana miniatury kontekstu |
| Sortowanie tabeli metryk | `table[data-tabela]` + `th[data-sortuj]` | demonstracja wartości `tabular-nums` |
| Filtrowanie zestawu znaków | `input[data-filtruj]` + `data-filtr-pozycja` | zawężenie zestawu diakrytyków |
| Kopiowanie nazwy żetonu | `data-kopiuj` | powiadomienie `window.dnToast` |
| Komunikat zamiast blokady | `data-komunikat` | powiadomienie informacyjne |
| Symulacja renderowania | `data-symuluj` | przejście stanu „pracuje" → „gotowe" |
| Postęp pokrycia znaków | `.dn-postep[data-postep-do]` | animowany pasek |
| Okno modalne z kartą kroju | `data-otworz-modal` | `<dialog>` |
| Przełącznik gęstości | atrybut `data-gestosc` na `<html>` | przebudowa warstwy roboczej planszy |
| Tętno kropki | `.pt-tetno` | jedyny ruch ciągły, 2,4 s |

---

## 10. Źródła

| Zakres | Plik / rozdział |
|---|---|
| Trzy role, trzy kroje; skala; wagi; interlinie; odstępy | kontrakt systemu projektowego |
| Kontrakt kierunku, pokrętła | kontrakt systemu projektowego · kierunek systemu projektowego |
| Definicje żetonów typograficznych | `WYNIK/zasoby/zetony/zetony.css` sekcja 3 |
| Deklaracje `@font-face`, podzbiory latin-ext | `WYNIK/zasoby/zetony/fonty.css` |
| Klasy komponentów (`.dn-wpis`, `.dn-dane`, `.dn-tabela`, `.dn-plakietka`) | `WYNIK/zasoby/css/komponenty.css` · kontrakt systemu projektowego |
| Warstwa prototypu (`.pt-konsola`, `.pt-etykieta`, `.pt-mono`, `.pt-tetno`) | `WYNIK/zasoby/prototyp.css` |
| Nazwy modułów i okien operacyjnych | kontrakt systemu projektowego |
| Przykładowa ścieżka repozytorium i wynik powłoki | `WYNIK/05-okna/moduly/terminal.html` |
| Przykładowy kod parowania urządzenia | `WYNIK/05-okna/platformowe/ustawienia.html` · `mobile.html` |
| Zasada zero blokad | kontrakt systemu projektowego (zasada zero blokad) |
| Dostępność, `prefers-reduced-motion`, fokus | kontrakt systemu projektowego |

---

## 11. Lista kontrolna wykonania

- [x] Wszystkie nazwy własne zgodne z dokumentacją
- [x] Zero elementów bez pokrycia w dokumentacji
- [x] Zero wartości szesnastkowych wpisanych wprost w warstwie lokalnej
- [x] Zero `#000000`
- [x] Zero `disabled`
- [x] Zero emoji jako ikon — wyłącznie inline SVG 24×24, obrys 1,75, `currentColor`
- [x] Zero Lorem ipsum i zmyślonych metryk
- [x] Oba motywy działają; pasek górny atramentowy w obu
- [x] Fokus widoczny na każdej kontrolce
- [x] Panel „O tym opracowaniu" obecny
- [x] Polskie diakrytyki poprawne w całym pliku
- [x] Odnośnik powrotu do `../INDEKS.html`

---

*Danaco Console — AI Operating Environment · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o.*
