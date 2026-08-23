# Danaco Console — system projektowy

| Pole | Treść |
|---|---|
| **Produkt** | Danaco Console — AI Operating Environment (warstwa wizualna v2.0) |
| **Rodzaj** | Opracowanie merytoryczno-techniczne — fundament systemu projektowego (A1) |
| **Producent** | Danaco Holding Group Sp. z o.o. |
| **Twórca** | Dariusz Naharnowicz |
| **Wersja** | v2.0 |
| **Status** | Deweloperski |
| **Data** | 2026-08-14 (pakiet wizualny źródłowy: 2026-08-11) |
| **Odbiorcy** | projektanci interfejsu, deweloperzy warstwy klienckiej, autorzy makiet i prototypów, osoby prowadzące odbiór wizualny, redaktorzy tekstów interfejsu |
| **Zakres** | definicja systemu projektowego: kontrakt kierunku, filozofia, zasady nienegocjowalne, anatomia trójwarstwowa, katalog anty-domyślnych, zasada jednego akcentu, rama kokpitu, element sygnaturowy, gęstości, zarządzanie systemem, słownik pojęć, mapa zależności |
| **Czego NIE zawiera** | pełnych kart poszczególnych komponentów (osobne opracowanie katalogu komponentów), tabeli 33 pomiarów kontrastu w rozbiciu na pary (plik `zasoby/zetony/kontrasty.json` i opracowanie barw), księgi znaku i konstrukcji godła (opracowanie marki), makiet okien operacyjnych (katalog `05-okna/`), specyfikacji protokołu i warstwy funkcjonalnej Danaco Pilot |

---

## Spis treści

1. [Czym jest system projektowy Danaco Console i po co powstał](#1-czym-jest-system-projektowy-danaco-console-i-po-co-powstał)
2. [Filozofia „instrumentu pomiarowego”](#2-filozofia-instrumentu-pomiarowego)
3. [Pięć zasad nienegocjowalnych](#3-pięć-zasad-nienegocjowalnych)
4. [Anatomia systemu — trzy warstwy](#4-anatomia-systemu--trzy-warstwy)
5. [Katalog anty-domyślnych — zablokowane odruchy](#5-katalog-anty-domyślnych--zablokowane-odruchy)
6. [Zasada jednego akcentu](#6-zasada-jednego-akcentu)
7. [Rama kokpitu jako decyzja systemowa](#7-rama-kokpitu-jako-decyzja-systemowa)
8. [Element sygnaturowy — kropka sygnału](#8-element-sygnaturowy--kropka-sygnału)
9. [Gęstość: zwarta · przestronna · dotyk](#9-gęstość-zwarta--przestronna--dotyk)
10. [Zarządzanie systemem (governance)](#10-zarządzanie-systemem-governance)
11. [Słownik pojęć systemu (PL/EN)](#11-słownik-pojęć-systemu-plen)
12. [Mapa zależności systemu](#12-mapa-zależności-systemu)
13. [Decyzje projektowe](#13-decyzje-projektowe)

---

## 1. Czym jest system projektowy Danaco Console i po co powstał

### 1.1. Definicja robocza

System projektowy Danaco Console to **zamknięty zbiór wartości, reguł i komponentów**, z którego wyprowadza się każdy widok platformy. Nie jest zestawem sugestii ani biblioteką inspiracji. Jest **kontraktem wykonawczym**: element, którego nie da się wyprowadzić z tego zbioru, nie wchodzi do interfejsu.

System obejmuje pięć warstw materialnych:

| Warstwa | Nośnik | Rola |
|---|---|---|
| Wartości | `zasoby/zetony/zetony.css`, `zetony.json` | jedyne źródło prawdy wartości wizualnych |
| Pomiar | `zasoby/zetony/kontrasty.json` (33 pomiary) | dowód zgodności z WCAG 2.1 AA — pomiar, nie deklaracja |
| Podłoże | `zasoby/css/fundament.css` | wyzerowanie, typografia bazowa, fokus globalny, drobne wzorce tekstowe |
| Komponenty | `zasoby/css/komponenty.css` (1207 linii, klasy `.dn-*`) | biblioteka kontrolek ze wszystkimi stanami |
| Znak i ikony | `zasoby/marka/`, `zasoby/ikony/` (manifest + SVG) | tożsamość i alfabet piktograficzny |

### 1.2. Powód powstania

Pakiet v2.0 **zastępuje w całości** dotychczasową, zastępczą warstwę wizualną (granat + złoto, tarcza z wagą, kroje Cormorant/Inter), uznaną przez Właściciela za nieobowiązującą. Powód zastąpienia nie był estetyczny, lecz funkcjonalny: poprzednia warstwa opisywała **wizerunek instytucji**, podczas gdy produkt jest **stanowiskiem pracy**. Kokpit dowodzenia i sygnet z wagą to dwa różne obiecania.

Drugim powodem był rozjazd źródeł. Dokumentacja v1.0 opisywała dwie równoległe rodziny plików o tej samej dacie wersji i różnych wartościach (skala typografii `11·12·14·15·17·21·26·34·44` px wobec `12·13·15·16·18·22·28·36·46` px; paleta ciemna granatowa wobec grafitowej). System v2.0 zamyka tę dwoistość: **`zetony.css` i `zetony.json` są jedynym źródłem prawdy**, a każda wartość opisana prozą ustępuje wartości zapisanej w żetonie.

### 1.3. Kontrakt kierunku

Zdanie źródłowe, z którego wyprowadza się każda decyzja wizualna (`KIERUNEK.md`, rozdz. 1):

> Platforma operacyjna (**kokpit dowodzenia**) dla zawodowego **Operatora** zarządzającego cyfrową organizacją, w języku **monochromatycznej precyzji** — odcienie bieli w motywie jasnym, odcienie czerni w motywie ciemnym, **jeden chłodny sygnał** — z ciążeniem ku własnemu systemowi „**instrumentu pomiarowego**”: zwarta gęstość, dane krojem mono, zero dekoracji bez funkcji.

Zdanie to jest **kontraktem**, nie opisem. Każda decyzja wizualna musi dać się z niego wyprowadzić w jednym kroku; jeżeli wyprowadzenie wymaga dwóch kroków interpretacji, decyzja jest podejrzana.

### 1.4. Trzy pokrętła z uzasadnieniem

Kierunek jest sparametryzowany trzema pokrętłami. Wartości nie są preferencją stylistyczną — każda ma konsekwencję mierzalną w interfejsie.

| Pokrętło | Wartość | Uzasadnienie źródłowe | Konsekwencja mierzalna |
|---|---|---|---|
| `WARIANCJA_PROJEKTOWA` | **4 / 10** | Kokpit wymaga przewidywalności. Asymetria wyłącznie tam, gdzie niesie hierarchię: strefy Centrum dowodzenia, relacja koordynator–wykonawca. Zero ozdobnego chaosu. | jedna siatka 12 kolumn, przerwa 24 px, treść maks. 1200 px; jeden system narożników (3 · 6 · 8 · 10 · 14 · pill); jeden zestaw cieni per motyw |
| `INTENSYWNOSC_RUCHU` | **3 / 10** | Aplikacja robocza, nie strona marketingowa. Mikroprzejścia plus jeden ruch znaczący. | czasy 0,10 / 0,16 / 0,22 s, jedna krzywa `cubic-bezier(0.2, 0, 0, 1)`; jeden ruch ciągły: tętno 2,4 s |
| `GESTOSC_WIZUALNA` | **8 / 10** | Decyzja Właściciela: gęstość zwarta („kokpit”). Operator ma widzieć stan wielu procesów naraz, nie przewijać. | kontrolka 32 px, wiersz 36 px, pasek 48 px, panele 12–16 px, stopień bazowy 13 px |

**Reguła sprzężenia pokręteł.** Pokrętła nie są niezależne. Wysoka gęstość (8/10) wymusza niską wariancję (4/10) — w gęstym układzie każda niespodzianka kompozycyjna kosztuje czytelność. Niska wariancja z kolei uzasadnia niską intensywność ruchu (3/10): skoro układ jest przewidywalny, ruch nie musi go tłumaczyć, ma wyłącznie meldować stan.

### 1.5. Cztery dowody, które interfejs składa bez słów

Kierunek przekłada się na cztery twierdzenia, które użytkownik ma odczytać z samego układu, zanim przeczyta jakikolwiek tekst (`ksiega-marki.html`, rozdz. 1):

| Dowód | Nośnik wizualny |
|---|---|
| To nie jest czat | rama kokpitu, pas kart sesji, okna operacyjne i pas komunikacji zamiast jednej rozmowy |
| Praca biegnie w tle | kropka sygnału tętni na kartach sesji i przy nadawcach |
| Modele pracują dla siebie nawzajem | para okien koordynator–wykonawca ze wspólną szyną przekazania |
| Wszystko jest konfigurowalne per okno | panel ustawień przy każdym oknie komunikacji |

---

## 2. Filozofia „instrumentu pomiarowego”

### 2.1. Metafora nadrzędna

Instrument pomiarowy — miernik laboratoryjny, konsola realizatorska, pulpit dyspozytorski — ma trzy cechy, które system przejmuje dosłownie:

1. **Nie ozdabia odczytu.** Każdy piksel, który nie niesie informacji, obniża wiarygodność odczytu.
2. **Ma jeden kolor sygnalizacyjny.** Instrument jest szary, czarny albo biały; kolor pojawia się wyłącznie tam, gdzie coś się dzieje.
3. **Jest gęsty, bo pomiar jest zbiorowy.** Operator patrzy na wiele wskaźników naraz. Rozrzedzenie układu to nie „oddech”, to utrata pola widzenia.

### 2.2. Monochromatyczna precyzja — trzy konsekwencje

| Konsekwencja | Zapis w systemie |
|---|---|
| Neutralne są **czysto neutralne** — równe składowe RGB, bez podbarwienia slate ani ciepłego beżu | skala 18 kroków od `#FFFFFF` do `#0A0A0A`, każdy krok o równych składowych |
| Oba motywy są **definiowane osobno**, nie wywodzone przez odwrócenie | bloki `:root[data-theme='light']` i `:root[data-theme='dark']` w `zetony.css` — dwa niezależne komplety żetonów semantycznych |
| Działanie główne to **inwersja atramentu**, nie druga barwa | `--dn-atrament` = `szary-900` w motywie jasnym, `szary-100` w ciemnym; przycisk główny jest czarny na jasnym i biały na ciemnym |

Monochromatyczna inwersja jest **najsilniejszym możliwym akcentem, jaki nie wprowadza drugiej barwy**. To jest sedno decyzji: czerń działa, sygnał wskazuje.

### 2.3. Dane krojem mono

Trzy kroje niosą trzy różne komunikaty. Różnica kroju jest nośnikiem znaczenia, nie ozdobą.

| Rola | Krój | Żeton | Komunikat |
|---|---|---|---|
| Nagłówki, tytuły środowisk, logotyp | Space Grotesk 500–700 | `--dn-ff-naglowek` | **wejście do środowiska i tożsamość marki** |
| Interfejs i treść | IBM Plex Sans 400–700 | `--dn-ff-bazowa` | **praca** |
| Dane, identyfikatory, terminal | IBM Plex Mono 400–600 | `--dn-ff-mono` | **maszyna** |

Wszystkie liczby renderowane krojem mono używają `tabular-nums` (`.dn-dane`, `.dn-liczba` w `fundament.css`) — kolumny liczb muszą się zgadzać co do piksela, bo są odczytem, nie tekstem.

Test diakrytyków obowiązujący dla każdego kroju: `ŁÓDŹ ŻÓŁĆ GĘŚLĄ JAŹŃ`.

### 2.4. Zero dekoracji bez funkcji — lista rozstrzygnięć

| Element | Rozstrzygnięcie systemu | Uzasadnienie |
|---|---|---|
| Gradienty | istnieją **dwa** (`--dn-grad-atrament`, `--dn-grad-sygnal`) i służą **wyłącznie ilustracji**: awatarom bez zdjęcia, rdzeniowi Always On Display, grafice brandowej | gradient nie niesie stanu; jako tło przycisku, karty czy sekcji byłby czystą dekoracją |
| Cienie | dwa komplety tonowane per motyw, jeden kierunek światła (z góry); `--dn-cien-sygnal` **zarezerwowany** dla pierścienia pola aktywnego i poświaty kropki | cień oznacza wyniesienie warstwy — nie „ładność” |
| Rozmycie (glassmorfizm) | powierzchnie kryjące; rozmycie **wyłącznie** w nakładce modala (`::backdrop`, `blur(2px)`) | rozmycie ma odciąć uwagę od warstwy pod spodem — to jedyna jego funkcja |
| Promienie narożników | jeden system: 3 · 6 · 8 · 10 · 14 · pill | małe promienie to język precyzji; duże promienie to język konsumencki |
| Animacje wejścia | `opacity` + `translateY(4–8 px)`, maks. 220 ms; brak parallaxu, brak scrollytellingu, brak sprężyn | ruch potwierdza, nie ozdabia — nic nie wjeżdża „z ekranu”, nic nie skacze |

### 2.5. Głos interfejsu jako część systemu

System projektowy obejmuje też warstwę słowną, bo tekst zajmuje w kokpicie więcej powierzchni niż jakikolwiek komponent. Głos Danaco Console to **rzeczowy meldunek, nie gawęda**.

| Reguła | Zapis |
|---|---|
| Działania | tryb rozkazujący: „Uruchom pętlę”, „Zapisz plan”, „Dodaj katalog” |
| Stany | meldunek: „krok 3 z 9 · obieg 1” |
| Błąd | co się stało **plus** co zrobić; nigdy sam kod, nigdy sama emocja |
| Liczby | po polsku: przecinek dziesiętny, spacja tysięcy, jednostka po liczbie |
| Terminologia | stała, bez synonimów — jedno pojęcie, jedno słowo (rozdz. 11) |
| Zakazy | wykrzykniki, „Ups!”, infantylizowanie Operatora |

---

## 3. Pięć zasad nienegocjowalnych

Zasady zapisane są w nagłówku `zetony.css` (linie 13–14) i powtórzone maszynowo w `zetony.json` (`_meta.zasady`). Nienegocjowalność oznacza: **zmiana wymaga decyzji Właściciela i nowego numeru wersji pakietu**, nie decyzji projektanta w toku prac.

### 3.1. Zestawienie

| # | Zasada | Zapis techniczny | Test odbioru |
|---|---|---|---|
| 1 | **Jednostka 4 px** | `--dn-od-0…16` = 0 / 4 / 8 / 12 / 16 / 20 / 24 / 32 / 40 / 48 / 64 px | żaden odstęp w arkuszu nie jest liczbą spoza skali |
| 2 | **Oba motywy równoprawne** | bloki `light` i `dark` definiowane osobno; `color-scheme` ustawiany per motyw; `prefers-color-scheme` honorowany przy braku jawnego wyboru | każdy widok sprawdzony w obu motywach; żaden nie jest „wersją drugą” |
| 3 | **Stan nigdy samym kolorem** | każdy komponent stanu niesie ikonę albo etykietę: `.dn-plakietka`, `.dn-krok`, `.dn-toast`, `.dn-kropka` | wyłączenie barw (symulacja monochromii) nie odbiera odczytu stanu |
| 4 | **Zero blokad (ADL-017)** | brak atrybutu `disabled`; stany przez ARIA: `aria-pressed`, `aria-busy`, `aria-invalid`, `aria-selected`, `aria-current` | żaden przycisk nie jest wyszarzoną bramą; niegotowość komunikuje opis obok albo komunikat po naciśnięciu |
| 5 | **WCAG 2.1 AA na wejściu** | 33 zmierzone pary w `kontrasty.json`; wartość bez przekroczonego progu nie wchodzi do żetonów | pomiar, nie deklaracja — przebieg `kontrasty.json` jest bramą każdej fali wdrożenia |

**Zasada szósta zapisana w `zetony.json`** (rozszerzenie, nie odrębna reguła): zakaz `#000000`; `#FFFFFF` wyłącznie jako powierzchnia kart i tekst na atramencie. Traktujemy ją jako konsekwencję zasady 5 — czysta czerń zabija głębię i powoduje halację na OLED, a biel jako tło całej strony podnosi jasność ponad próg komfortu przy pracy wielogodzinnej.

### 3.2. Zasada 1 — jednostka 4 px w praktyce

| Wielokrotność | Żeton | Typowe zastosowanie |
|---|---|---|
| 1× | `--dn-od-1` = 4 px | odstęp etykieta–kontrolka, gap paska ikon |
| 2× | `--dn-od-2` = 8 px | gap wewnątrz przycisku, odstęp między plakietkami |
| 3× | `--dn-od-3` = 12 px | rytm wewnętrzny paneli (`--dn-odstep-panel`) |
| 4× | `--dn-od-4` = 16 px | wyściółka pozioma kart i przycisków |
| 5× | `--dn-od-5` = 20 px | wyściółka modala, wyściółka Always On Display |
| 6× | `--dn-od-6` = 24 px | rytm między sekcjami (`--dn-odstep-sekcji`), wyściółka karty środowiska |
| 8× | `--dn-od-8` = 32 px | odsunięcie treści dokumentowej |
| 10× / 12× / 16× | 40 / 48 / 64 px | pusty stan, marginesy stron ekspozycyjnych |

Wymiary kontrolek również są wielokrotnościami: 32 (8×), 36 (9×), 40 (10×), 44 (11×), 48 (12×), 56 (14×), 224 (56×), 320 (80×).

### 3.3. Zasada 4 — cztery sytuacje niegotowości

Zero blokad nie oznacza „system pozwala na wszystko”. Oznacza: **system nie udaje, że przycisk nie istnieje**. Cztery sytuacje i ich obsługa:

| Sytuacja | Reguła | Przykład z dokumentacji |
|---|---|---|
| Element zależy od danych w tym samym widoku | zawsze klikalny; brak danych sygnalizowany komunikatem przy właściwym polu po kliknięciu | przycisk „Zaloguj się” przy niewypełnionym formularzu |
| Element zależy od zakończenia poprzedniego etapu | zawsze klikalny; kliknięcie wyświetla komunikat o brakującym warunku | **Deployment Panel** modułu **Apps** przed zakończeniem **Frontend Workspace** / **Backend Workspace** |
| Element zależy od konfiguracji Operatora | obecny i w pełni dostępny; brak konfiguracji = wartość domyślna | przełącznik macierzy izolacji; powiązanie międzymodułowe |
| Element dotyczy progu uwierzytelnienia | przycisk „Pomiń” klikalny w obu fazach wdrożenia; skuteczność zależy od ustawienia „Wymóg logowania” | okno rejestracji i logowania |

Dodatkowe rozstrzygnięcia ADL-017 wiążące dla makiet:
- Odliczanie przy „Wyślij ponownie” ma charakter **wyłącznie informacyjny** — nie blokuje kliknięcia.
- Metoda uwierzytelniania wyłączona w konfiguracji **nie jest renderowana w ogóle** — „brak metody = mniej segmentów, nie zablokowany segment”.
- Jedyny sankcjonowany wyjątek to pole w widoku tylko do odczytu (`readonly`), gdzie element fizycznie nie ma czego przyjąć.

---

## 4. Anatomia systemu — trzy warstwy

### 4.1. Schemat warstw

```
┌─────────────────────────────────────────────────────────────────────────┐
│  WARSTWA 1 · PRYMITYWY            zetony.css, sekcja 1                  │
│  --dn-szary-0…1000 (18 kroków) · --dn-sygnal-100…800 (8 kroków)         │
│  --dn-zielen/bursztyn/czerwien-100/200/400/700                          │
│  ZAKAZ UŻYCIA WPROST W KOMPONENTACH                                     │
└──────────────────────────────┬──────────────────────────────────────────┘
                               │  mapowanie per motyw (dwa osobne bloki)
                               ▼
┌─────────────────────────────────────────────────────────────────────────┐
│  WARSTWA 2 · SEMANTYCZNE          zetony.css, sekcje 10–11              │
│  --dn-tlo · --dn-powierzchnia · --dn-powierzchnia-2 · --dn-panel        │
│  --dn-tekst · --dn-tekst-2 · --dn-tekst-3 · --dn-tekst-inv             │
│  --dn-atrament(-hover, -tekst) · --dn-sygnal(-mocny, -wypelnienie, …)   │
│  --dn-obrys(-subtelny, -mocny) · --dn-fokus · --dn-kropka               │
│  --dn-sukces/ostrzezenie/blad/informacja -tekst/-tlo/-obrys             │
│  --dn-cien-1/2/3/lg/sygnal · --dn-hover · --dn-wcisniecie · --dn-nakladka│
└──────────────────────────────┬──────────────────────────────────────────┘
                               │  jedyne wejście dozwolone dla komponentu
                               ▼
┌─────────────────────────────────────────────────────────────────────────┐
│  WARSTWA 3 · KOMPONENTOWE         komponenty.css, klasy .dn-*           │
│  .dn-btn · .dn-pole · .dn-karta · .dn-tabela · .dn-wpis · .dn-krok …    │
│  Każdy komponent: spoczynek · najechanie · naciśnięcie · fokus ·        │
│  wybrany · ładowanie · błąd                                             │
└─────────────────────────────────────────────────────────────────────────┘

   ↕ POPRZECZNIE (niezależne od motywu, jedna definicja dla całości):
     typografia · przestrzeń · promienie · ruch · wymiary · siatka ·
     warstwy z-index · gradienty · rama kokpitu
```

### 4.2. Reguła „komponent nie sięga po prymityw”

**Zapis reguły:** komponent sięga wyłącznie po żetony semantyczne (`--dn-tlo`, `--dn-tekst`, `--dn-obrys`…), nigdy po prymitywy (`--dn-szary-500`).

**Dlaczego.** Prymityw jest wartością, semantyczny żeton jest **rolą**. Komponent, który sięga po `--dn-szary-500`, przestaje działać w drugim motywie — bo „szary 500” znaczy co innego na białym papierze niż na czarnym. Komponent, który sięga po `--dn-tekst-3`, działa w obu motywach bez jednej linii kodu więcej, ponieważ rola „metadane” jest zdefiniowana osobno dla każdego motywu.

**Wyjątki jawne i zamknięte.** W bibliotece komponentów istnieją trzy odwołania do prymitywu — każde uzasadnione i policzone:

| Miejsce | Odwołanie | Uzasadnienie |
|---|---|---|
| `.dn-btn--sygnal` (kolor treści) | `var(--dn-szary-0)` | biel na wypełnieniu sygnałowym jest **stała w obu motywach** — pary zmierzone: 6,40:1 (jasny) i 4,63:1 (ciemny). Żeton `--dn-tekst-inv` przełącza się z motywem i zepsułby ten kontrast. |
| `.dn-przelacznik:checked::before` (suwak) | `var(--dn-szary-0)` | ta sama przesłanka — suwak leży na wypełnieniu sygnałowym w obu motywach |
| `.dn-awatar` / `.dn-aod-rdzen` (treść na gradiencie) | `var(--dn-szary-100)` / `var(--dn-szary-0)` | gradienty są **ilustracyjne i niezależne od motywu**; tekst na nich musi być stały |

Wszystkie trzy wyjątki mają wspólny mianownik: **powierzchnia pod tekstem nie przełącza się z motywem**, więc tekst też nie może.

### 4.3. Reguła kolejności ładowania

Kolejność arkuszy jest częścią kontraktu, nie szczegółem technicznym:

```
1. fonty.css        @font-face (Space Grotesk, IBM Plex Sans, IBM Plex Mono)
2. zetony.css       prymitywy → semantyczne → poprzeczne → dotyk → gęstość → ruch
3. fundament.css    wyzerowanie, podłoże, typografia bazowa, fokus globalny
4. komponenty.css   biblioteka .dn-*
5. <style> lokalny  wyłącznie układ konkretnego widoku, wyłącznie var(--dn-*)
```

Odwrócenie kolejności 2 ↔ 3 psuje fokus globalny; odwrócenie 3 ↔ 4 psuje nagłówki komponentów. Arkusz lokalny widoku **nie definiuje wartości** — komponuje wyłącznie układ z żetonów już istniejących.

### 4.4. Model stanów komponentu

Każdy komponent interaktywny ma siedem stanów. Model jest wspólny, nie per komponent.

| Stan | Nośnik | Reguła wizualna |
|---|---|---|
| Spoczynek | wygląd bazowy | powierzchnia neutralna, bez sygnału |
| Najechanie | `:hover` | `--dn-hover` jako tło albo uniesienie 1–2 px dla elementów o dużej masie |
| Naciśnięcie | `:active` | `translateY(1px)`, czas `--dn-czas-1` |
| Fokus | `:focus-visible` | pierścień `--dn-fokus` 2 px + odsunięcie 2 px — wyłącznie z klawiatury |
| Wybrany (trwały) | `aria-pressed` / `aria-selected` / `aria-current` | tło `--dn-sygnal-tlo`, obrys `--dn-sygnal-obrys`, tekst `--dn-sygnal` |
| Ładowanie | `aria-busy="true"` | wskaźnik obok treści; **przycisk pozostaje klikalny** |
| Błąd | `aria-invalid="true"` | obrys `--dn-blad-obrys`, tekst `--dn-blad-tekst` **plus** ikona lub etykieta |

Konsekwencja zasady 3 i 4 łącznie: **stan jest zawsze parą** (barwa + znak) i **nigdy nie odbiera interakcji**.

---

## 5. Katalog anty-domyślnych — zablokowane odruchy

Katalog jest listą odruchów, które w projektach interfejsów AI pojawiają się automatycznie i które w tym systemie są **bezwzględnie zakazane**. Każdy zakaz ma zamiennik i uzasadnienie wyprowadzone z kontraktu kierunku.

| # | Zablokowany odruch | Zamiast tego | Uzasadnienie z kontraktu |
|---|---|---|---|
| 1 | Fioletowy gradient „AI” | monochrom + jeden błękit sygnałowy | kontrakt mówi „monochromatyczna precyzja, jeden chłodny sygnał”; gradient AI jest znakiem kategorii, nie produktu — nie odróżnia kokpitu od czatu |
| 2 | Inter + slate-900 jako niezadeklarowana baza | IBM Plex Sans + czysta neutralna skala, zadeklarowane wprost | slate jest podbarwiony niebieskim; kontrakt żąda „odcieni bieli i czerni” dosłownie, o równych składowych RGB |
| 3 | Trzy równe karty funkcji | strefy o malejącej masie (Centrum dowodzenia), asymetria koordynator–wykonawca | równa masa oznacza równą wagę; w kokpicie wagi nie są równe — środowisko waży więcej niż komponent własny, komponent więcej niż ustawienie |
| 4 | Wyszarzone przyciski jako bramy | **zero blokad** — komunikat po naciśnięciu albo opis obok | ADL-017: domyślne zachowanie systemu to wykonanie polecenia; wyszarzenie jest odmową bez uzasadnienia |
| 5 | Stan samym kolorem | **zawsze** ikona albo etykieta | zasada 3; ok. 8% populacji męskiej ma zaburzenie rozróżniania barw, a kokpit czyta się także w monochromatycznym zrzucie i przy wysokim kontraście systemowym |
| 6 | Emoji jako ikony | wyłącznie SVG z zestawu (Lucide, obrys 1,75, siatka 24×24) | emoji renderuje się inaczej na każdej platformie, nie przyjmuje `currentColor` i nie skaluje się do 14 px bez utraty czytelności |
| 7 | Dziewięć kolorów tła dla dziewięciu nadawców | trzy klasy semantyczne + ikona + etykieta + plakietka roli | dziewięć teł rozbiłoby spójność i złamało zasadę „stan nigdy samym kolorem”; rola różnicuje się znakiem, nie barwą |
| 8 | Glassmorfizm wszędzie | powierzchnie kryjące; rozmycie **tylko** w nakładce modala | rozmycie obniża kontrast tekstu poniżej progu i zmienia go zależnie od treści pod spodem — nie da się go zmierzyć raz |
| 9 | `#000000` / `#FFFFFF` jako tło strony | `#0F0F0F` / `#F4F4F4` | czysta czerń zabija głębię i powoduje halację na OLED; czysta biel jako pełne tło podnosi jasność ponad próg komfortu przy pracy wielogodzinnej |
| 10 | Cień „domyślny szary wszędzie” | dwa komplety cieni tonowanych per motyw; `--dn-cien-sygnal` zarezerwowany | cień szary na ciemnym tle jest niewidoczny albo brudny; komplet ciemny jest głębszy i o innej krzywej |
| 11 | Lorem ipsum, „Jan Kowalski”, zmyślone metryki | treści operacyjne z domeny produktu, oznaczone jako przykładowe | makieta z fikcyjną treścią nie weryfikuje długości etykiet, zawijania ani gęstości; zmyślone metryki wchodzą potem do prezentacji jako fakty |

### 5.1. Rozszerzenie katalogu — odruchy odnotowane w toku prac

Poniższe pozycje nie występują w źródłowym katalogu `KIERUNEK.md`, ale wynikają wprost z zasad nienegocjowalnych. Odnotowane jako rozstrzygnięcia ponad źródła (rozdz. 13).

| Odruch | Zamiast tego | Zasada źródłowa |
|---|---|---|
| Wartość szesnastkowa wpisana wprost w arkuszu widoku | `var(--dn-*)` zawsze | zasada „jedno źródło prawdy wartości” |
| Odstęp spoza skali (np. 10 px, 18 px, 30 px) | najbliższa wielokrotność 4 px | zasada 1 |
| Usunięcie pierścienia fokusu „bo brzydki” | pierścień nieusuwalny bez równoważnego wskaźnika | zasada 5 |
| Animacja dekoracyjna „żeby żyło” | jeden ruch znaczący na widok — tętno pracy w tle | `INTENSYWNOSC_RUCHU` 3/10 |
| Gradient jako tło przycisku, karty albo sekcji | powierzchnia kryjąca | zakres gradientów: wyłącznie ilustracja |
| Nowa barwa dla nowego modułu | ta sama paleta, różnicowanie ikoną i etykietą | zasada jednego akcentu (rozdz. 6) |

---

## 6. Zasada jednego akcentu

### 6.1. Treść zasady

> **Sygnał zajmuje najwyżej 5% powierzchni ekranu. Sygnał nigdy nie jest tłem sekcji. Sygnał nigdy nie jest przyciskiem głównym.**

Neutralne zajmują pozostałe ok. 95% powierzchni: tła, powierzchnie, tekst, obrysy.

### 6.2. Podział ról barwy

| Rodzina | Zakres | Udział powierzchni | Ograniczenia |
|---|---|---|---|
| Neutralne (18 kroków, `#FFFFFF` → `#0A0A0A`) | tła, powierzchnie, tekst, obrysy | ok. 95% | zakaz `#000000`; `#FFFFFF` wyłącznie jako powierzchnia kart i tekst na atramencie |
| Sygnał (`#3B6FE0` i rodzina, 8 kroków) | fokus · wybór · odnośniki · postęp · kropka pracy | **≤ 5%** | nigdy tło sekcji, nigdy przycisk główny |
| Stany (zieleń · bursztyn · czerwień) | wyłącznie meldunek stanu | punktowo | zawsze z ikoną albo etykietą; informacja = rodzina sygnału |

### 6.3. Dlaczego działanie główne to inwersja atramentu, a nie sygnał

To rozstrzygnięcie jest osią całego systemu barwy. Cztery argumenty:

1. **Rozdzielenie ról.** Gdyby przycisk główny był sygnałowy, sygnał znaczyłby jednocześnie „tu kliknij” i „tu biegnie praca”. Dwa znaczenia jednej barwy to brak znaczenia. Rozdział jest zapisany zdaniem: **czerń działa, sygnał wskazuje**.
2. **Budżet 5%.** Przyciski główne występują w każdym panelu. Gdyby były sygnałowe, sam ten jeden komponent zjadłby cały budżet akcentu i kropka sygnału przestałaby być widoczna jako wyjątek.
3. **Siła kontrastu.** Inwersja atramentu daje 17,76:1 (biel na `szary-900`) i 16,23:1 (`szary-950` na `szary-100`). Wypełnienie sygnałowe daje 6,40:1 i 4,63:1. Inwersja jest po prostu **mocniejszym** akcentem.
4. **Neutralność motywu.** Inwersja definiuje się sama w obu motywach (czarny ↔ biały). Wypełnienie sygnałowe wymaga dwóch różnych kroków rodziny (`sygnal-600` w jasnym, `sygnal-500` w ciemnym), żeby zachować próg.

**Kiedy zatem przycisk sygnałowy?** `.dn-btn--sygnal` jest zarezerwowany wyłącznie dla **działania systemowego typu „uruchom / zatwierdź plan”** — i **jeden na widok**. Nie zastępuje atramentu jako domyślnego przycisku głównego.

### 6.4. Rejestr wystąpień sygnału

Pełna lista miejsc, w których sygnał ma prawo wystąpić. Lista jest zamknięta.

| Wystąpienie | Klasa / żeton | Powierzchnia |
|---|---|---|
| Pierścień fokusu | `--dn-fokus`, `:focus-visible` | 2 px obwodu jednej kontrolki naraz |
| Kropka sygnału | `.dn-kropka`, `--dn-kropka` | 6 × 6 px na wystąpienie |
| Wstęga aktywnej karty sesji | `.dn-karta-sesji[aria-selected]::before` | 2 px wysokości |
| Podkreślenie aktywnej zakładki | `.dn-zakladka[aria-selected]::after` | 2 px wysokości |
| Znacznik pozycji bocznej nawigacji | `.dn-boczna-pozycja[aria-current]::before` | 2 × 16 px |
| Tło pozycji wybranej | `--dn-sygnal-tlo` | wąskie pasy, nie sekcje |
| Wypełnienie paska postępu | `.dn-postep-wartosc` | 4 px wysokości |
| Suwak włączony | `.dn-przelacznik:checked` | 36 × 20 px |
| Odnośnik | `a { color: var(--dn-sygnal) }` | tekst, nie tło |
| Krawędź wpisu klasy „inteligencja” | `.dn-wpis--inteligencja` | 2 px lewej krawędzi |
| Przycisk systemowy „uruchom” | `.dn-btn--sygnal` | jeden na widok |
| Rdzeń awatara inteligencji / AOD | `--dn-grad-sygnal` | element ilustracyjny 28–36 px |

### 6.5. Test budżetu

Sposób weryfikacji zasady na makiecie: zrzut ekranu, izolacja pikseli należących do rodziny sygnału, podzielenie przez powierzchnię całkowitą. Wynik powyżej 5% oznacza błąd projektowy — najczęściej: sygnałowe tło sekcji, sygnałowy przycisk główny albo zbyt wiele elementów jednocześnie w stanie „wybrany”.

---

## 7. Rama kokpitu jako decyzja systemowa

### 7.1. Treść decyzji

**Pasek górny jest zawsze atramentowy (`--dn-rama` = `szary-925`, `#131313`) — w obu motywach.** To jedyny element interfejsu, który nie przełącza się razem z motywem (poza godłem, które zachowuje barwy własne).

### 7.2. Komplet żetonów ramy

| Żeton | Wartość | Rola |
|---|---|---|
| `--dn-rama` | `szary-925` | tło paska górnego, tło dymka objaśnienia |
| `--dn-rama-tekst` | `szary-100` | treść pierwszorzędna na ramie |
| `--dn-rama-tekst-2` | `szary-400` | treść drugorzędna, ikony w spoczynku |
| `--dn-rama-hover` | `rgba(255,255,255,.08)` | najechanie na ramie |
| `--dn-rama-obrys` | `rgba(255,255,255,.10)` | krawędzie kontrolek na ramie |

Żetony ramy są zdefiniowane **poza blokami motywów**, w sekcji poprzecznej `zetony.css` (sekcja 2) — to zapis techniczny decyzji o niezmienności.

### 7.3. Cztery uzasadnienia

| Uzasadnienie | Treść |
|---|---|
| **Stanowisko dowodzenia** | ciemna rama wokół przełączalnej treści daje efekt pulpitu: treść jest tym, co się zmienia, rama jest tym, co stoi |
| **Stały dom dla stałych elementów** | godło, wyszukiwarka, wskaźniki i przełącznik motywu mają jedno miejsce o niezmiennym kontraście — Operator nie uczy się ich położenia dwa razy |
| **Kotwica orientacyjna przy przełączeniu motywu** | przełączenie motywu zmienia całą treść; niezmienna rama utrzymuje ciągłość percepcyjną i skraca czas ponownej orientacji |
| **Jeden komplet kontrastów zamiast dwóch** | kontrolki na ramie (`.dn-btn-ikona--na-ramie`, `.dn-pasek-szukaj`) mają jeden zestaw wartości mierzonych raz, nie dwa |

### 7.4. Zasięg ramy

```
┌──────────────────────────────────────────────────────────────────┐
│  ▓▓▓  PASEK GÓRNY 48 px  — --dn-rama, STAŁY W OBU MOTYWACH  ▓▓▓  │  z-index 100
├──────────────────────────────────────────────────────────────────┤
│  PAS KART SESJI 36 px — --dn-powierzchnia-2, PRZEŁĄCZA SIĘ       │
├────────────────┬─────────────────────────────────────────────────┤
│  BOCZNA        │  OBSZAR ROBOCZY                                 │
│  NAWIGACJA     │  --dn-tlo / --dn-powierzchnia                   │
│  224 px        │  PRZEŁĄCZA SIĘ Z MOTYWEM                        │
│  --dn-panel    │                                                 │
│  PRZEŁĄCZA SIĘ ├─────────────────────────────────────────────────┤
│                │  PAS KOMUNIKACJI 320 px — PRZEŁĄCZA SIĘ         │
└────────────────┴─────────────────────────────────────────────────┘
```

Rama to **wyłącznie pasek górny**. Boczna nawigacja, pas kart sesji i pas komunikacji przełączają się normalnie — nie są ramą, są treścią.

### 7.5. Jeden element pochodny

Dymek objaśnienia (`.dn-tooltip-tresc`) używa `--dn-rama` jako tła i `--dn-rama-tekst` jako treści. To celowe: dymek jest warstwą **systemową**, nie treściową — należy do ramy pojęciowo, mimo że pojawia się nad treścią.

---

## 8. Element sygnaturowy — kropka sygnału

### 8.1. Definicja

**Kropka sygnału** to wypełniony punkt błękitu sygnałowego o średnicy **6 px** (`--dn-wym-kropka`), oznaczający jedno zdanie: **„tu biegnie praca”**.

Jest to element, po którym platformę się zapamięta — jedna rzecz, nie zestaw rzeczy. Wywodzi się z metafory nadrzędnej znaku: grot „❯” to miejsce wydania polecenia, kropka to dowód, że maszyna pracuje.

### 8.2. Pięć miejsc wystąpienia

| # | Miejsce | Forma | Znaczenie |
|---|---|---|---|
| 1 | **Godło** (sygnet „Delegacja”) | kropka zamykająca podwójny grot — znak czyta się „»».” | zlecenie przekazane od koordynatora do wykonawcy; kropka to praca, która właśnie ruszyła |
| 2 | **Emblematy czterech środowisk** (TalkIn · WorkSpace · CodeStudio · MultitaskingAI) | dokładnie **jedna** wypełniona kropka w każdym emblemacie | wspólne DNA czterech środowisk bez wprowadzania czterech barw |
| 3 | **Karty sesji** | kropka pulsująca (`.dn-kropka--tetno`) | proces biegnie w tle także wtedy, gdy Operator patrzy gdzie indziej |
| 4 | **Okno komunikacji** | kropka przy nadawcy aktywnie piszącym (`.dn-wpis--pracuje`) | jednostka inteligencji generuje odpowiedź teraz |
| 5 | **Always On Display** | rdzeń awatara (`.dn-aod-rdzen` z pierścieniem tętna) | platforma pracuje nawet przy zwiniętym oknie |

### 8.3. Tętno 2,4 s — jedyny ruch ciągły

| Parametr | Wartość | Zapis |
|---|---|---|
| Czas cyklu | **2,4 s** | `--dn-czas-tetno` |
| Krzywa | `cubic-bezier(0.2, 0, 0, 1)` | `--dn-ease` — ta sama dla całego interfejsu |
| Mechanika | pierścień `box-shadow` 0 → 5 px, zanik do przezroczystości w 40% cyklu | `@keyframes dn-tetno` |
| Wariant ograniczonego ruchu | **statyczny pierścień** 2 px zamiast animacji | `@media (prefers-reduced-motion: reduce)` |

**Dlaczego 2,4 s.** Rytm ma być wolniejszy od spoczynkowego tętna człowieka (ok. 60–80/min, czyli 0,75–1,0 s), żeby nie wywoływać poczucia pośpiechu, i szybszy od progu, przy którym ruch przestaje być odbierany jako powiązany (ok. 4 s). 2,4 s czyta się jako **spokojna praca w toku**, nie jako alarm.

**Dlaczego jeden ruch ciągły.** Przy `INTENSYWNOSC_RUCHU` = 3/10 każdy dodatkowy ruch ciągły odbiera znaczenie temu jednemu. Jeśli w kokpicie pulsuje pięć rzeczy, żadna nie pulsuje. Reguła operacyjna: **jeden ruch znaczący na widok**.

### 8.4. Warianty stanu kropki

Kropka ma warianty stanu, ale nadal obowiązuje zasada 3 — nigdy nie występuje sama.

| Wariant | Klasa | Barwa | Wymóg towarzyszący |
|---|---|---|---|
| Sygnał (praca) | `.dn-kropka` | `--dn-kropka` | etykieta lub kontekst karty sesji |
| Sukces | `.dn-kropka--sukces` | `--dn-sukces-tekst` | ikona `ptaszek` albo etykieta |
| Ostrzeżenie | `.dn-kropka--ostrzezenie` | `--dn-ostrzezenie-tekst` | ikona `ostrzezenie` albo etykieta |
| Błąd | `.dn-kropka--blad` | `--dn-blad-tekst` | ikona `blad` albo etykieta |
| Neutralna | `.dn-kropka--neutralna` | `--dn-tekst-3` | etykieta |

### 8.5. Czego kropka nie oznacza

| Nie oznacza | Właściwy komponent |
|---|---|
| „nieprzeczytane” | plakietka z liczbą (`.dn-plakietka`) |
| „nowość” / „polecane” | etykieta tekstowa |
| „element wybrany” | `aria-selected` + tło `--dn-sygnal-tlo` |
| „element zaznaczony na liście” | `.dn-check` |
| ozdoba przy nagłówku | brak — usunąć |

---

## 9. Gęstość: zwarta · przestronna · dotyk

### 9.1. Trzy tryby i ich status

| Tryb | Wywołanie | Status | Przeznaczenie |
|---|---|---|---|
| **Zwarta** | domyślny (brak atrybutu) | **obowiązujący** | biurko, kokpit, `GESTOSC_WIZUALNA` 8/10 |
| **Przestronna** | `data-gestosc="przestronna"` na `<html>` | **przygotowany, domyślnie nieaktywny** | wybór Operatora; ekrany prezentacyjne; praca długodystansowa |
| **Dotyk** | `@media (pointer: coarse)` — automatycznie | **aktywny warunkowo** | urządzenia dotykowe; nakłada się na tryb bieżący, nie zastępuje go |

### 9.2. Tabela porównawcza wymiarów

| Żeton | Zwarta (domyślna) | Przestronna | Dotyk (`pointer: coarse`) |
|---|---:|---:|---:|
| `--dn-fs-base` | 13 px | **14 px** | 13 px (bez zmiany) |
| `--dn-lh-bazowy` | 1,45 | **1,5** | 1,45 (bez zmiany) |
| `--dn-wym-kontrolka` | 32 px | **40 px** | **40 px** |
| `--dn-wym-ikonowy` | 32 px | **40 px** | **40 px** |
| `--dn-wym-wiersz` | 36 px | **44 px** | **44 px** |
| `--dn-wym-pasek` | 48 px | **56 px** | 48 px (bez zmiany) |
| `--dn-wym-pas-kart` | 36 px | **40 px** | 36 px (bez zmiany) |
| `--dn-wym-check` | 16 px | 16 px | **20 px** |
| `--dn-wym-przelacznik-szer` | 36 px | 36 px | **44 px** |
| `--dn-wym-przelacznik-wys` | 20 px | 20 px | **24 px** |
| `--dn-odstep-panel` | 12 px (`od-3`) | **20 px** (`od-5`) | 12 px (bez zmiany) |
| `--dn-odstep-sekcji` | 24 px (`od-6`) | **32 px** (`od-8`) | 24 px (bez zmiany) |
| `--dn-wym-boczna` | 224 px | 224 px | 224 px |
| `--dn-wym-pas-komunikacji` | 320 px | 320 px | 320 px |

### 9.3. Reguła kompozycji trybów

Dotyk i przestronna **nie kolidują** — nakładają się bez konfliktu wartości:

```
  ZWARTA (bazowa, :root)
     │
     ├── + DOTYK (pointer: coarse)      → podnosi cele dotykowe: kontrolka, wiersz,
     │                                     check, przełącznik
     │
     └── + PRZESTRONNA (data-gestosc)   → podnosi rytm i stopień: fs-base, interlinia,
                                           pasek, pas kart, odstępy paneli i sekcji

  ZWARTA + DOTYK + PRZESTRONNA jednocześnie:
     kontrolka 40 · wiersz 44 · check 20 · przełącznik 44×24 ·
     pasek 56 · pas kart 40 · fs-base 14 · odstępy 20/32
     → brak sprzeczności: obie warstwy podają dla kontrolki i wiersza te same
       wartości (40 / 44), a pozostałe żetony ustawia tylko jedna z nich
```

Wartości kontrolki (40 px) i wiersza (44 px) są **celowo identyczne** w obu warstwach — to nie zbieg okoliczności, lecz warunek bezkolizyjności. Gdyby przestronna podawała np. 44 px dla kontrolki, wynik zależałby od specyficzności selektorów, a nie od decyzji projektowej.

### 9.4. Co gęstość zmienia, a czego nie zmienia

| Zmienia | Nie zmienia |
|---|---|
| wysokości kontrolek i wierszy | barwy i kontrasty |
| stopień bazowy i interlinię | promienie narożników |
| rytm wewnętrzny paneli i sekcji | siatkę 12 kolumn i punkty łamania |
| wysokość paska i pasa kart | wymiary kropki, wstęgi, spinnera, pierścienia fokusu |
| — | szerokość bocznej nawigacji (224 px) i pasa komunikacji (320 px) |
| — | reguły komponentów (żaden selektor nie ma wariantu „dla gęstości”) |

**Reguła nadrzędna:** gęstość zmienia się **żetonem, nie wyjątkiem**. W bibliotece komponentów nie istnieje ani jeden selektor warunkowany gęstością — cała zmiana zachodzi w warstwie wartości.

### 9.5. Punkty łamania a gęstość

Punkty łamania są niezależne od gęstości i dotyczą **układu**, nie wymiarów kontrolek.

| Żeton | Wartość | Zachowanie układu |
|---|---|---|
| `--dn-bp-w1` | 640 px | telefon poziomo — widok mobilny |
| `--dn-bp-w2` | 960 px | tablet — boczna nawigacja zwija się do ikon |
| `--dn-bp-w3` | 1280 px | biurko — pełny kokpit |
| `--dn-bp-w4` | 1600 px | szerokie biurko — dwa okna komunikacji obok siebie |

---

## 10. Zarządzanie systemem (governance)

### 10.1. Czego nie wolno zmienić po cichu

Lista zamknięta, przeniesiona z `HANDOFF.md` rozdz. 4. Zmiana którejkolwiek pozycji wymaga decyzji Właściciela i nowego numeru wersji pakietu.

| # | Pozycja | Nośnik |
|---|---|---|
| 1 | wartości żetonów | `zasoby/zetony/zetony.css`, `zetony.json` |
| 2 | geometria znaku i emblematów (krzywe, nie fonty) | `zasoby/marka/logo/*.svg`, `zasoby/marka/srodowiska/*.svg` |
| 3 | zasada jednego akcentu (sygnał ≤ 5% ekranu, nigdy tła sekcji) | rozdz. 6 |
| 4 | pierścień fokusu i zachowanie `prefers-reduced-motion` | `fundament.css`, `zetony.css` sekcja 14 |
| 5 | gęstość zwarta jako domyślna | `zetony.css` sekcja 6 |

Do listy dołączamy — jako konsekwencje zasad nienegocjowalnych — pozycje szóstą i siódmą:

| # | Pozycja | Nośnik |
|---|---|---|
| 6 | pasek górny atramentowy w obu motywach | `zetony.css` sekcja 2 |
| 7 | zakaz `disabled` (ADL-017) | `komponenty.css`, nagłówek pliku |

### 10.2. Co wolno zmienić i w jakim trybie

| Zakres zmiany | Tryb | Kto rozstrzyga | Skutek wersyjny |
|---|---|---|---|
| Dodanie nowego **modyfikatora** do istniejącej klasy `.dn-*` | swobodny, w ramach istniejących żetonów | projektant komponentu | wersja pomocnicza |
| Dodanie nowej **klasy komponentowej** | procedura z rozdz. 10.4 | prowadzący system projektowy | wersja pomocnicza |
| Dodanie nowej **ikony domenowej** | rysunek na siatce 24×24, obrys 1,75, wpis do `manifest.json` | prowadzący system projektowy | wersja pomocnicza |
| Zmiana **wartości żetonu semantycznego** | wymaga ponownego pomiaru kontrastów i aktualizacji `kontrasty.json` | Właściciel | wersja główna |
| Dodanie **prymitywu** (nowy krok skali) | wymaga uzasadnienia luki w istniejącej skali | Właściciel | wersja główna |
| Dodanie **nowej rodziny barwnej** | zakazane bez zmiany kontraktu kierunku | Właściciel | nowy kontrakt kierunku |
| Zmiana **kroju pisma** | zakazana bez zmiany kontraktu kierunku | Właściciel | nowy kontrakt kierunku |

### 10.3. Kto rozstrzyga — trzy role

| Rola | Zakres decyzji | Ograniczenie |
|---|---|---|
| **Właściciel** (Dariusz Naharnowicz) | kontrakt kierunku, pokrętła, zasady nienegocjowalne, wartości żetonów, znak | jedyny podmiot uprawniony do zmiany listy z rozdz. 10.1 |
| **Prowadzący system projektowy** | katalog komponentów, katalog ikon, wzorce okien, redakcja opracowań | działa wyłącznie w granicach żetonów i zasad |
| **Wykonawca widoku** (projektant / deweloper) | kompozycja widoku z istniejących komponentów, teksty operacyjne | nie tworzy wartości; brak komponentu zgłasza procedurą 10.4 |

**Zasada rozstrzygania sporów:** przy rozbieżności między opisem prozą a plikiem CSS **rozstrzyga plik CSS** jako jedyne źródło prawdy o klasach i wartościach. Zasada przeniesiona z README katalogu komponentów (zasada 4) i rozszerzona na cały pakiet v2.0.

### 10.4. Procedura dodania komponentu

```
KROK 1 · SPRAWDŹ, CZY KOMPONENT JUŻ ISTNIEJE
   ├─ przegląd komponenty.css (119 tokenów klasowych .dn-*)
   ├─ przegląd katalogu komponentów (27 komponentów w 7 grupach)
   └─ jeśli istnieje rola, ale brakuje wariantu → to MODYFIKATOR, nie komponent

KROK 2 · UZASADNIJ ROLĘ
   ├─ jaka rola w interfejsie nie ma dziś nośnika?
   ├─ w którym oknie z inwentarza występuje?
   └─ brak odpowiedzi na drugie pytanie = komponent nie wchodzi

KROK 3 · ZBUDUJ WYŁĄCZNIE Z ŻETONÓW SEMANTYCZNYCH
   ├─ zero wartości zaszytych, zero prymitywów
   ├─ nazwa klasy po polsku wg konwencji .dn-<rzeczownik>
   └─ modyfikatory: .dn-<rzecz>--<przymiotnik>, elementy: .dn-<rzecz>-<część>

KROK 4 · SKOMPLETUJ SIEDEM STANÓW
   spoczynek · najechanie · naciśnięcie · fokus · wybrany · ładowanie · błąd
   ├─ stany przez ARIA, nigdy przez disabled
   └─ stan zawsze parą: barwa + ikona/etykieta

KROK 5 · ZMIERZ
   ├─ kontrast każdej nowej pary tekst/tło → dopisz do kontrasty.json
   ├─ sprawdź w OBU motywach
   └─ sprawdź w gęstości zwartej i przestronnej oraz przy pointer: coarse

KROK 6 · UDOKUMENTUJ
   ├─ karta komponentu w katalogu (rola, warianty, stany, wymiary, użycie)
   ├─ żywa demonstracja w galerii komponentów
   └─ wpis do listy klas

KROK 7 · ODBIÓR
   └─ brama: zgodność z listą kontrolną (rozdz. 10.6) w obu motywach
```

### 10.5. Wersjonowanie

| Poziom | Zapis | Wyzwalacz | Konsekwencja dla wdrożeń |
|---|---|---|---|
| **Kontrakt kierunku** | zmiana zdania kierunkowego | decyzja Właściciela o nowym języku wizualnym | pakiet zastępowany w całości (tak jak v2.0 zastąpił poprzedni) |
| **Wersja główna** (`v2.0` → `v3.0`) | zmiana wartości żetonów, prymitywów, krojów, znaku | decyzja Właściciela | wymagany ponowny odbiór wszystkich okien w obu motywach |
| **Wersja pomocnicza** (`v2.0` → `v2.1`) | nowe komponenty, nowe modyfikatory, nowe ikony, poprawki dokumentacji | prowadzący system | wdrożenie przyrostowe; istniejące okna bez zmian |

**Zasada jednoczesności plików.** Zmiana wartości musi trafić **jednocześnie** do `zetony.css`, `zetony.json` i — gdy dotyczy pary tekst/tło — do `kontrasty.json`. Rozjazd między tymi plikami jest przyczyną, dla której poprzednia generacja systemu wymagała zastąpienia (rozdz. 1.2). Nie powtarzamy go.

### 10.6. Fale wdrożenia i bramy odbioru

| Fala | Zakres | Brama weryfikacyjna |
|---|---|---|
| 1. Żetony | wartości, motywy, wymiary, warstwy, rama, gradienty | aplikacja rusza na nowej palecie; przebieg `kontrasty.json` (pomiar, nie deklaracja) |
| 2. Komponenty | `fundament.css` + biblioteka `.dn-*` | galeria komponentów jako wzorzec odbioru wizualnego, w obu motywach |
| 3. Ikony i marka | zestaw SVG, manifest, favicon, ikony aplikacji | zgodność nazw z `manifest.json`; literówka w nazwie zatrzymuje kompilację |
| 4. Okna | okno komunikacji + powłoki środowisk wg makiet | porównanie z makietami w obu motywach |

### 10.7. Lista kontrolna odbioru pliku

- [ ] wszystkie nazwy własne zgodne z dokumentacją (zero parafraz)
- [ ] zero elementów bez pokrycia w dokumentacji
- [ ] zero wartości szesnastkowych wpisanych wprost — tylko `var(--dn-*)`
- [ ] zero `#000000`
- [ ] zero `disabled`
- [ ] zero emoji jako ikon
- [ ] zero Lorem ipsum, zmyślonych nazwisk i metryk
- [ ] oba motywy działają i wyglądają poprawnie
- [ ] pasek górny atramentowy w obu motywach
- [ ] stan nigdy samym kolorem (ikona lub etykieta towarzyszy)
- [ ] fokus widoczny na każdej kontrolce
- [ ] interakcje działają po otwarciu z `file://`
- [ ] panel „O tym opracowaniu” obecny
- [ ] polskie diakrytyki poprawne w całym pliku

---

## 11. Słownik pojęć systemu (PL/EN)

### 11.1. Pojęcia systemu projektowego

| PL | EN | Definicja obowiązująca |
|---|---|---|
| żeton | design token | nazwana wartość wizualna; jedyny sposób zapisu wartości w systemie |
| prymityw | primitive token | surowa wartość skali (`--dn-szary-500`); zakaz użycia w komponentach |
| żeton semantyczny | semantic token | rola per motyw (`--dn-tekst-2`); jedyne wejście dozwolone dla komponentu |
| atrament | ink | barwa działania głównego: `szary-900` w motywie jasnym, `szary-100` w ciemnym |
| sygnał | signal | jedyna barwa akcentu (`#3B6FE0` i rodzina); wskazuje, nie działa |
| kropka sygnału | signal dot | element sygnaturowy 6 px oznaczający „tu biegnie praca” |
| tętno | pulse | animacja kropki 2,4 s; jedyny ruch ciągły interfejsu |
| rama kokpitu | cockpit frame | pasek górny atramentowy, stały w obu motywach |
| gęstość zwarta | compact density | tryb domyślny: kontrolka 32 px, wiersz 36 px |
| gęstość przestronna | comfortable density | tryb przygotowany: `data-gestosc="przestronna"` |
| anty-domyślne | anti-defaults | katalog zablokowanych odruchów projektowych |
| zero blokad | no hard gates (ADL-017) | zasada: system nie odbiera klikalności |
| wstęga | ribbon | pasek 2 px oznaczający element aktywny (karta sesji, karta środowiska) |
| plakietka | badge | etykieta stanu lub roli; zawsze z ikoną albo tekstem |
| medalion | medallion | kwadratowy nośnik ikony nadawcy we wpisie okna komunikacji |
| powierzchnia | surface | tło karty lub okna (`--dn-powierzchnia`) |
| podłoże | canvas | tło strony (`--dn-tlo`) |
| nakładka | overlay / backdrop | przyciemnienie pod modalem |
| grot | chevron | znak „❯” — sygnatura wejścia polecenia |

### 11.2. Pojęcia platformy

Nazwy własne modułów, okien i ról pozostają **w brzmieniu angielskim** zgodnie z dokumentacją — zakaz tłumaczenia i parafrazowania.

| PL (pojęcie) | EN (nazwa własna, nieprzekładalna) | Uwaga |
|---|---|---|
| Operator | Operator | zawodowy użytkownik platformy |
| środowisko | — | TalkIn · WorkSpace · CodeStudio · MultitaskingAI |
| moduł | — | Studio · Research · Library · Translate · Browser · Assistant · Roundtable · Workspace · Automations · Design · Apps · Terminal · Developer · Diagnostics · Agents |
| okno operacyjne | — | np. **Studio Editor**, **Workflow Builder**, **Agent Builder**, **Diagnostics Center** |
| okno komunikacji | **Chat Window** | pas wspólny wszystkim modułom |
| karta sesji | session tab | element pasa kart sesji |
| koordynator | **Coordinator** | rola MultitaskingAI |
| wykonawca | **Executor 1** / **Executor 2** | role MultitaskingAI |
| walidator | **Executor 3 / Validator** | rola MultitaskingAI |
| sieć podagentów | **Subagent Network** | piąta pozycja zespołu |
| panel orkiestracji | — | 6 sekcji, zastępuje boczną nawigację w MultitaskingAI |
| centrum dowodzenia | — | strona główna, trzy strefy |
| Always On Display | **Always On Display** | funkcja globalna, warstwa `z-index` 1200 |
| centrum poleceń | **Command Center** | najwyższa warstwa, `z-index` 1300 |

### 11.3. Terminologia stała tekstów interfejsu

Jedno pojęcie, jedno słowo — synonimy zabronione: **Operator · środowisko · moduł · sesja · okno komunikacji · zlecenie · pętla · koordynator / wykonawca / walidator · praca w tle · punkt decyzji**.

---

## 12. Mapa zależności systemu

### 12.1. Diagram pełny

```
        ┌───────────────────────────────────────────────────────────┐
        │  KONTRAKT KIERUNKU  (KIERUNEK.md)                         │
        │  monochromatyczna precyzja · instrument pomiarowy         │
        │  pokrętła: wariancja 4 · ruch 3 · gęstość 8               │
        └────────────────────────┬──────────────────────────────────┘
                                 │ wyprowadzenie
                                 ▼
   ┌─────────────────────────────────────────────────────────────────────┐
   │  ŻETONY                                                             │
   │  zetony.css ── zetony.json ── kontrasty.json (33 pomiary)          │
   │  ┌──────────┐   ┌──────────────┐   ┌───────────────────────────┐   │
   │  │prymitywy │──►│ semantyczne  │   │ poprzeczne:               │   │
   │  │18+8+12   │   │ jasny/ciemny │   │ typografia · przestrzeń · │   │
   │  └──────────┘   └──────────────┘   │ ruch · wymiary · siatka · │   │
   │                                     │ warstwy · rama · gradienty│   │
   │                                     └───────────────────────────┘   │
   └───────────┬─────────────────────────────────────────────────────────┘
               │
               ▼
   ┌───────────────────────────┐        ┌────────────────────────────────┐
   │  FUNDAMENT                │        │  MARKA I IKONY                 │
   │  fundament.css            │        │  marka/logo · marka/srodowiska │
   │  wyzerowanie · podłoże ·  │◄──────►│  ikony/manifest.json (82) ·    │
   │  typografia · fokus ·     │        │  ikony/svg/*.svg               │
   │  paski · wzorce tekstowe  │        │  Lucide ISC · 24×24 · 1,75     │
   └───────────┬───────────────┘        └───────────────┬────────────────┘
               │                                        │
               └──────────────┬─────────────────────────┘
                              ▼
   ┌─────────────────────────────────────────────────────────────────────┐
   │  KOMPONENTY  ── komponenty.css (1207 linii · 119 tokenów .dn-*)     │
   │  przyciski · pola · wybór · nawigacja · dane · nakładki ·           │
   │  tożsamość · komunikacja                                            │
   │  każdy komponent: 7 stanów · zero disabled · stan nigdy samym       │
   │  kolorem                                                            │
   └────────────────────────────┬────────────────────────────────────────┘
                                │
                                ▼
   ┌─────────────────────────────────────────────────────────────────────┐
   │  OKNA                                                               │
   │  przepływ:      Okno startowe → Rejestracja i logowanie →           │
   │                 Centrum dowodzenia                                  │
   │  powłoki:       TalkIn · WorkSpace · CodeStudio · MultitaskingAI    │
   │  platformowe:   Okno Konfiguracji · Okno Ustawień ·                 │
   │                 Always On Display · Mobile                          │
   │  wspólne:       Chat Window                                         │
   │  operacyjne:    15 modułów × okna właściwe                          │
   │  role:          Coordinator · Executor 1 · Executor 2 ·             │
   │                 Executor 3 / Validator · Subagent Network           │
   └────────────────────────────┬────────────────────────────────────────┘
                                │
                                ▼
   ┌─────────────────────────────────────────────────────────────────────┐
   │  ŚRODOWISKA — trójstopniowa hierarchia                              │
   │  Środowisko („w jakim trybie pracuję?")                             │
   │      └─► Moduł („jakie zadanie wykonuję?")                          │
   │              └─► Okno operacyjne („jakim narzędziem realizuję?")    │
   └─────────────────────────────────────────────────────────────────────┘
```

### 12.2. Kierunek zależności — reguła jednokierunkowości

```
   żetony ──► fundament ──► komponenty ──► okna ──► środowiska
      ▲                                                  │
      └──────────────  ZAKAZ  ◄──────────────────────────┘
              (okno nie definiuje wartości;
               środowisko nie definiuje komponentu)
```

| Poziom | Wolno | Nie wolno |
|---|---|---|
| Środowisko | wybierać okna, ustawiać kolejność modułów | definiować własnych barw ani komponentów |
| Okno | komponować z komponentów, definiować układ lokalny | definiować wartości; wprowadzać własnych klas globalnych |
| Komponent | sięgać po żetony semantyczne | sięgać po prymitywy (poza trzema wyjątkami z rozdz. 4.2) |
| Fundament | ustawiać podłoże i zachowania globalne | definiować wygląd komponentów |
| Żetony | definiować wartości | odwoływać się do czegokolwiek wyżej |

### 12.3. Macierz wpływu zmiany

| Zmiana w… | Wymusza przegląd… | Koszt |
|---|---|---|
| prymitywie | wszystkich żetonów semantycznych, `kontrasty.json`, wszystkich okien | **bardzo wysoki** |
| żetonie semantycznym | komponentów używających roli, `kontrasty.json`, okien w obu motywach | **wysoki** |
| żetonie poprzecznym (wymiar, ruch) | wszystkich komponentów; okna zwykle bez zmian | średni |
| komponencie | okien, w których występuje | średni |
| oknie | wyłącznie tego okna | niski |
| tekście interfejsu | tego widoku; sprawdzenie długości etykiet | niski |

---

## 13. Decyzje projektowe

Sekcja zawiera rozstrzygnięcia podjęte **ponad źródła** — tam, gdzie dokumentacja nie rozstrzygała formy. Każde rozstrzygnięcie jest zgodne z katalogiem komponentów i zasadami nienegocjowalnymi.

| # | Zagadnienie | Stan w źródłach | Rozstrzygnięcie | Uzasadnienie |
|---|---|---|---|---|
| **D-01** | Uzasadnienia poszczególnych zakazów w katalogu anty-domyślnych | `KIERUNEK.md` i `KANON.md` podają **parę** „odruch → zamiast tego”, bez uzasadnień | Dopisano kolumnę uzasadnień (rozdz. 5) | Katalog bez uzasadnień jest listą zakazów do obejścia; z uzasadnieniami staje się narzędziem decyzyjnym w sytuacjach nieprzewidzianych |
| **D-02** | Rozszerzenie katalogu anty-domyślnych o 6 pozycji (rozdz. 5.1) | brak w źródłach | Dodano jako **rozszerzenie oznaczone jawnie**, wyprowadzone z zasad nienegocjowalnych | Pozycje te były łamane najczęściej w toku prac (wartości hex wprost, odstępy spoza skali, usuwanie fokusu); zapis czyni je sprawdzalnymi |
| **D-03** | Trzy wyjątki „komponent sięga po prymityw” | `komponenty.css` zawiera trzy takie odwołania, bez komentarza wyjaśniającego | Zinwentaryzowano i uzasadniono wspólnym mianownikiem: powierzchnia pod tekstem nie przełącza się z motywem (rozdz. 4.2) | Bez zapisu wyjątki wyglądają na błąd i zostałyby „naprawione”, psując kontrast na wypełnieniu sygnałowym |
| **D-04** | Status trybu „dotyk” względem gęstości | `zetony.css` definiuje `pointer: coarse` i `data-gestosc` w osobnych sekcjach, bez opisu relacji | Ustalono, że dotyk **nakłada się** na tryb bieżący, nie zastępuje go; udokumentowano bezkolizyjność (rozdz. 9.3) | Wartości 40 px i 44 px są identyczne w obu warstwach, więc złożenie jest deterministyczne niezależnie od specyficzności selektorów |
| **D-05** | Uzasadnienie czasu tętna 2,4 s | źródła podają wartość, nie uzasadnienie | Dopisano uzasadnienie percepcyjne: wolniej niż tętno spoczynkowe, szybciej niż próg utraty powiązania ruchu (rozdz. 8.3) | Bez uzasadnienia wartość wygląda na dowolną i jest pierwszą kandydatką do „drobnej korekty” |
| **D-06** | Rejestr wystąpień sygnału | źródła podają zasadę ≤ 5% i listę ról, bez zamkniętego rejestru miejsc | Zbudowano **zamkniętą listę 12 wystąpień** (rozdz. 6.4) na podstawie faktycznych selektorów w `komponenty.css` i `fundament.css` | Zasada procentowa bez rejestru jest niesprawdzalna; rejestr pozwala odbierać widok bez pomiaru pikseli |
| **D-07** | Lista „czego kropka nie oznacza” | brak w źródłach | Dopisano (rozdz. 8.5) z odesłaniem do właściwego komponentu dla każdego przypadku | Element sygnaturowy rozmywa się przez nadużycie; zapis negatywny chroni skuteczniej niż pozytywny |
| **D-08** | Procedura dodania komponentu (7 kroków) | źródła opisują zasady komponentów i fale wdrożenia, nie opisują procedury rozszerzania | Zbudowano procedurę z istniejących elementów: katalog + żetony semantyczne + 7 stanów + pomiar + dokumentacja + brama odbioru (rozdz. 10.4) | Procedura nie wprowadza nowych reguł — porządkuje istniejące w kolejność wykonawczą |
| **D-09** | Model wersjonowania (kontrakt / główna / pomocnicza) | źródła podają wersję v2.0 i fakt zastąpienia v1.0, bez modelu | Zaproponowano trójstopniowy model z wyzwalaczami i konsekwencjami (rozdz. 10.5) | Trzy poziomy odpowiadają trzem warstwom anatomii: kierunek → wartości → komponenty |
| **D-10** | Podział ról decyzyjnych | źródła wskazują decyzje Właściciela punktowo (gęstość, emblematy, paleta) | Uogólniono do trzech ról z zakresami i ograniczeniami (rozdz. 10.3) | Zakresy odtworzono z faktycznych decyzji odnotowanych w `KIERUNEK.md` i README pakietu; nie wprowadzono nowych uprawnień |
| **D-11** | Macierz wpływu zmiany | brak w źródłach | Zbudowano na podstawie kierunku zależności z rozdz. 12.2 | Jest konsekwencją arytmetyczną jednokierunkowości zależności, nie nową regułą |
| **D-12** | Rozstrzygnięcie sporu proza ↔ CSS dla całego pakietu | zasada zapisana w README **katalogu komponentów** (zasada 4), ograniczona do klas | Rozszerzono na cały pakiet v2.0: przy rozbieżności rozstrzyga plik CSS | Rozjazd źródeł był przyczyną zastąpienia poprzedniej generacji (rozdz. 1.2); rozszerzenie zasady zamyka tę klasę problemów |
| **D-13** | Sprzężenie pokręteł | źródła podają trzy wartości niezależnie | Opisano regułę sprzężenia: gęstość 8 wymusza wariancję 4, wariancja 4 uzasadnia ruch 3 (rozdz. 1.4) | Bez opisu sprzężenia pokrętła wyglądają na trzy niezależne suwaki i kuszą do zmiany pojedynczej |
| **D-14** | Godło jako jedyne odstępstwo od reguły żetonów w SVG | `HANDOFF.md`: „godło zachowuje barwy własne w obu motywach” | W opracowaniach HTML sygnet rysowany inline z `currentColor` dla grotów i `var(--dn-sygnal-400)` dla kropki | Grot dziedziczy barwę ramy (stałą), kropka zachowuje tożsamość znaku; oba zapisy pozostają odwołaniem `var(--dn-*)`, nie wartością szesnastkową |

### 13.1. Kwestie pozostawione otwarte

| Kwestia | Stan | Kto rozstrzyga |
|---|---|---|
| Zakres widoku Mobile (L-P-31) | otwarta w źródłach | Właściciel |
| Układ okien per moduł (M-4) | makiety przyjmują układ „okno modułu + pas komunikacji na dole” | Właściciel |
| Domyślne włączenie gęstości przestronnej jako wyboru w Oknie Ustawień | żetony gotowe, przełącznik nie jest zadeklarowany w inwentarzu okien | prowadzący system + Właściciel |
| Zachowanie pasa kart sesji przy przepełnieniu (przewijanie poziome vs zwężanie kart) | „przedmiot fazy wykonawczej” | prowadzący system |
| Powrót na stronę główną: godło czy odrębna ikona `dom` | rozbieżność między dwoma dokumentami źródłowymi | prowadzący system |

---

*Danaco Console — AI Operating Environment · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. — Dariusz Naharnowicz*
