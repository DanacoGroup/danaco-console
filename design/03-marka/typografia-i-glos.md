# Danaco Console — Typografia marki i głos marki

| | |
|---|---|
| **Produkt** | **Danaco Console** — AI Operating Environment (warstwa wizualna v2.0) |
| **Dokument** | Typografia marki i identyfikacja werbalna (głos marki) — opracowanie pełne |
| **Producent** | Danaco Holding Group Sp. z o.o. |
| **Twórca** | Dariusz Naharnowicz |
| **Kontakt** | support@danaco-group.pl |
| **Wersja** | v2.0 · Status: Deweloperski |
| **Data** | 2026-08-14 |
| **Zastępuje** | rozdz. 13 `ksiega-znaku.md` (skrót typograficzny) oraz rozdz. 4 przewodnika marki v1.0 (Cormorant / Inter) |
| **Odbiorcy** | projektanci interfejsu, deweloperzy front-end, redaktorzy tekstów interfejsu, autorzy dokumentacji użytkownika, dział wsparcia, wykonawcy materiałów firmowych |
| **Zakres** | trzy kroje marki i ich role, pełna skala typograficzna, zestawienia typograficzne, wersaliki, liczby, diakrytyki, wcielenie techniczne, katalog zakazów typograficznych; charakterystyka głosu marki, zasady pisania w interfejsie, słownik terminów obowiązujących, wzorce tekstów, mikrokopia, lista kontrolna redakcyjna |
| **Status wiążący** | **wiążący bezwarunkowo** — kroje, skala i słownik są zamknięte; odstępstwa wymagają pisemnej zgody Właściciela |
| **Dokument towarzyszący** | `typografia-i-glos.html` (wersja interaktywna z żywymi okazami i wyszukiwarką słownika) |
| **Źródła** | kontrakt systemu projektowego · kierunek systemu projektowego · `zasoby/zetony/zetony.css` · `zasoby/zetony/fonty.css` · `zasoby/css/fundament.css` · `zasoby/css/komponenty.css` · dokumentacja modułów `docs/` |

---

## Spis treści

**Część I — Typografia marki**

1. [Po co osobny dokument typograficzny](#1-po-co-osobny-dokument-typograficzny)
2. [Trzy kroje, trzy role — uzasadnienie doboru](#2-trzy-kroje-trzy-role--uzasadnienie-doboru)
3. [Space Grotesk — krój tożsamości](#3-space-grotesk--krój-tożsamości)
4. [IBM Plex Sans — krój pracy](#4-ibm-plex-sans--krój-pracy)
5. [IBM Plex Mono — krój maszyny](#5-ibm-plex-mono--krój-maszyny)
6. [Zasada nadrzędna i jej egzekwowanie](#6-zasada-nadrzędna-i-jej-egzekwowanie)
7. [Pełna skala typograficzna — dziewięć stopni](#7-pełna-skala-typograficzna--dziewięć-stopni)
8. [Zestawienia typograficzne — pary](#8-zestawienia-typograficzne--pary)
9. [Wersaliki](#9-wersaliki)
10. [Liczby](#10-liczby)
11. [Diakrytyki polskie](#11-diakrytyki-polskie)
12. [Wcielenie techniczne](#12-wcielenie-techniczne)
13. [Zakazy typograficzne](#13-zakazy-typograficzne)

**Część II — Głos marki**

14. [Charakterystyka głosu — pięć cech](#14-charakterystyka-głosu--pięć-cech)
15. [Zwrot do odbiorcy](#15-zwrot-do-odbiorcy)
16. [Zasady pisania w interfejsie](#16-zasady-pisania-w-interfejsie)
17. [Zasada zero blokad w języku](#17-zasada-zero-blokad-w-języku)
18. [Słownik terminów obowiązujących](#18-słownik-terminów-obowiązujących)
19. [Terminy zakazane i ich zamienniki](#19-terminy-zakazane-i-ich-zamienniki)
20. [Wzorce tekstów](#20-wzorce-tekstów)
21. [Mikrokopia](#21-mikrokopia)
22. [Lista kontrolna redakcyjna](#22-lista-kontrolna-redakcyjna)
23. [Decyzje projektowe](#23-decyzje-projektowe)

[Załącznik A · Karta typograficzna na jedną stronę](#załącznik-a--karta-typograficzna-na-jedną-stronę)
[Załącznik B · Ciąg kontrolny redakcyjny](#załącznik-b--ciąg-kontrolny-redakcyjny)
[Załącznik C · Mapa żetonów typograficznych](#załącznik-c--mapa-żetonów-typograficznych)

---

# CZĘŚĆ I — TYPOGRAFIA MARKI

## 1. Po co osobny dokument typograficzny

### 1.1 Typografia jako 95 % powierzchni produktu

Danaco Console jest kokpitem dowodzenia. W kokpicie nie ma zdjęć, ilustracji ani
dużych powierzchni barwnych — jest **tekst i liczba**. Po odjęciu paska górnego,
obrysów paneli i ikon 16-pikselowych, ponad dziewięćdziesiąt procent widocznej
powierzchni każdego okna operacyjnego stanowi **litera**. Typografia nie jest
w tym produkcie warstwą dekoracyjną nałożoną na interfejs — **jest interfejsem**.

Z tego wynika ranga niniejszego dokumentu. Błąd w doborze cienia kosztuje
estetykę; błąd w typografii kosztuje **czytelność identyfikatora zadania
w kolejce**, a więc czas Operatora i poprawność decyzji.

### 1.2 Co ten dokument rozstrzyga, a czego nie

| Rozstrzyga | Nie rozstrzyga |
|---|---|
| jakie trzy kroje obowiązują i dlaczego | geometrii znaku (→ `ksiega-znaku.md`, rozdz. 3–4) |
| dziewięciostopniową skalę i przypisanie stopni | barw tekstu i tła (→ kontrakt systemu projektowego) |
| gotowe zestawienia typograficzne dla komponentów | układu siatki i punktów łamania |
| reguły wersalików, liczb i diakrytyków | zestawu ikon (→ `emblematy-i-ikony.md`) |
| sposób osadzenia plików fontów | strategii marketingowej i pozycjonowania |
| **słownik terminów obowiązujących i wzorce tekstów interfejsu** | treści merytorycznej dokumentacji modułów |

### 1.3 Dlaczego typografia i głos w jednym dokumencie

Litera i słowo są jednym przedmiotem projektowym. Etykieta `WYWOŁANIA` złożona
wersalikami IBM Plex Mono z odstępem 0,14 em wygląda jak nagłówek kolumny danych,
ponieważ **jest** nagłówkiem kolumny danych — decyzja typograficzna i decyzja
redakcyjna zapadły w tym samym momencie. Rozdzielenie ich na dwa dokumenty
prowadzi do interfejsu, w którym forma i treść mówią co innego.

Zasada wiążąca: **każda reguła typograficzna z części I ma swój odpowiednik
redakcyjny w części II i odwrotnie.** Wzorzec „etykieta + wartość" z rozdz. 8.2
jest jednocześnie wzorcem redakcyjnym z rozdz. 20.5.

---

## 2. Trzy kroje, trzy role — uzasadnienie doboru

### 2.1 Kontrakt kierunku

kierunek systemu projektowego rozstrzyga:

> Trzy role, trzy kroje (wszystkie OFL, pełny `latin-ext`).
> **Space Grotesk = wejście do środowiska i tożsamość marki;
> Plex Sans = praca; Plex Mono = maszyna.**

Rozstrzygnięcie jest zamknięte. Poniższe rozdziały nie otwierają go ponownie —
tłumaczą je na poziom wykonawczy.

### 2.2 Zestawienie ról

| Rola | Krój | Wagi wgrywane | Żeton | Licencja | Udział powierzchni |
|---|---|---|---|---|---|
| Nagłówki, tytuły środowisk, logotyp | **Space Grotesk** | 500 · 600 · 700 | `--dn-ff-naglowek` | SIL OFL 1.1 | ok. 5 % |
| Interfejs i treść | **IBM Plex Sans** | 400 · 500 · 600 · 700 | `--dn-ff-bazowa` | SIL OFL 1.1 | ok. 65 % |
| Dane, identyfikatory, terminal | **IBM Plex Mono** | 400 · 500 · 600 | `--dn-ff-mono` | SIL OFL 1.1 | ok. 30 % |

Udział powierzchni jest szacunkiem wyprowadzonym z inwentarza okien
(kontrakt systemu projektowego) — nie pomiarem. Podano go, żeby uzmysłowić proporcję:
**krój tożsamości zajmuje najmniej miejsca i dlatego musi być najmocniejszy.**

### 2.3 Dlaczego trzy, a nie dwa albo pięć

**Dlaczego nie dwa.** Zestaw dwuelementowy (sans + mono) jest typowym rozwiązaniem
narzędzia deweloperskiego. Danaco Console nie jest narzędziem deweloperskim —
jest **środowiskiem pracy o czterech różnych trybach** (TalkIn, WorkSpace,
CodeStudio, MultitaskingAI). Moment wejścia do środowiska jest w tym produkcie
przeżyciem osobnym: Operator zmienia tryb myślenia, nie tylko widok. Ten moment
wymaga własnego kroju — takiego, którego nie widać nigdzie indziej w interfejsie.

**Dlaczego nie pięć.** Każdy dodatkowy krój to dodatkowe pliki `woff2`
(dziś: 22 pliki, 372 kB), dodatkowa decyzja przy każdym komponencie i dodatkowa
możliwość pomyłki. Trzy kroje dają **trzy jednoznaczne odpowiedzi na pytanie
„czym to jest?"** — a projektant nigdy nie stoi przed czwartą możliwością.

**Dlaczego akurat te trzy.** Space Grotesk i rodzina IBM Plex łączą trzy cechy
wymagane kontraktem jednocześnie: licencja OFL (osadzanie bez opłat i bez ryzyka
prawnego), pełny `latin-ext` (polskie diakrytyki narysowane, nie składane),
charakter techniczny bez ozdobności. Krojów spełniających wszystkie trzy warunki
jest niewiele; z tej puli wybrano parę o wspólnym rodowodzie inżynierskim
(Plex Sans + Plex Mono) i jeden krój kontrastujący (Space Grotesk).

### 2.4 Trzy kroje jako system znaczeń

```
   ┌──────────────────────────────────────────────────────────────┐
   │  SPACE GROTESK          „gdzie jestem"                       │
   │  ────────────────────────────────────────────────────        │
   │  DANACO · CodeStudio · Centrum dowodzenia · Konfiguracja      │
   │  geometryczny · osobowość · rzadki · duży stopień             │
   └───────────────────────────┬──────────────────────────────────┘
                               │  wejście
                               ▼
   ┌──────────────────────────────────────────────────────────────┐
   │  IBM PLEX SANS          „co robię"                           │
   │  ────────────────────────────────────────────────────        │
   │  Uruchom automatykę · Kolejka zadań · Zapisz zmianę          │
   │  humanistyczny · neutralny · wszechobecny · 13 px            │
   └───────────────────────────┬──────────────────────────────────┘
                               │  polecenie
                               ▼
   ┌──────────────────────────────────────────────────────────────┐
   │  IBM PLEX MONO          „co odpowiedział system"             │
   │  ────────────────────────────────────────────────────        │
   │  queue.action · SHA-256 · 00:12:04 · 1 284 wywołań           │
   │  stałoszerokościowy · tabular-nums · dane · dowód            │
   └──────────────────────────────────────────────────────────────┘
```

Kierunek strzałek jest kierunkiem pracy Operatora: wchodzi do środowiska
(Space Grotesk), wydaje polecenie (Plex Sans), odbiera wynik (Plex Mono).
Ten sam kierunek opowiada znak marki — dwa groty i kropka.

### 2.5 Zgodność z anty-domyślnymi

kierunek systemu projektowego blokuje odruch „Inter + slate-900 jako niezadeklarowana
baza". Konsekwencja typograficzna jest jednoznaczna: **Inter nie występuje
w produkcie w żadnej roli, także jako krój zapasowy.** Łańcuch zapasowy
(rozdz. 12.4) prowadzi do krojów systemowych, nie do Intera.

---

## 3. Space Grotesk — krój tożsamości

### 3.1 Charakterystyka

| Cecha | Wartość |
|---|---|
| Klasyfikacja | grotesk geometryczno-techniczny o proporcjach zbliżonych do neo-groteskowych |
| Pochodzenie rysunku | wyprowadzony z kroju wystawowego; charakter „technicznego szyldu" |
| Osobowość | wyraźna, ale bez ozdobności — pojedyncze znaki mają cechy rozpoznawalne (`a`, `g`, `y`, `Q`), całość pozostaje neutralna |
| Wersaliki | mocne, szerokie, o stabilnej podstawie — **główny powód wyboru** |
| Cyfry | proporcjonalne w rysunku, wystarczająco jednoznaczne w dużych stopniach |
| Licencja | **SIL Open Font License 1.1** |
| Wgrywane wagi | 500 · 600 · 700 (`latin` + `latin-ext`, sześć plików `woff2`, 76 kB) |

### 3.2 Dlaczego ten krój do tożsamości

Trzy powody, w kolejności wagi:

1. **Wersaliki.** Logotyp `DANACO` jest złożony wyłącznie wersalikami. Space
   Grotesk ma wersaliki o szerokim, stabilnym rysunku, które w wadze 700
   utrzymują masę także w stopniu 12 px na wizytówce. Krój humanistyczny
   traciłby w tej sytuacji strukturę.
2. **Odróżnialność od interfejsu.** Nagłówek `CodeStudio` musi być natychmiast
   rozpoznawalny jako **tytuł środowiska**, a nie jako etykieta pola. Kontrast
   rysunku między Space Grotesk a Plex Sans jest wystarczający, żeby ta różnica
   działała bez podnoszenia stopnia ani wagi.
3. **Brak ozdobności.** Kokpit nie znosi krojów o wyrazistej „osobowości
   projektanckiej". Space Grotesk mieści się w granicy: ma charakter, ale nie ma
   szeryfów, kontrastu grubości ani kaligraficznych zakończeń.

### 3.3 Wagi i ich zastosowanie

| Waga | Żeton | Zastosowanie | Stopnie |
|---:|---|---|---|
| **500** | `--dn-fw-srednia` | tytuły paneli w oknach operacyjnych, nagłówki sekcji drugiego rzędu | 16 · 20 px |
| **600** | `--dn-fw-polgruba` | **domyślna waga nagłówków** (`h1`–`h3` w `fundament.css`), nagłówki okien i modali, tytuły sekcji strony głównej | 16 · 20 · 24 px |
| **700** | `--dn-fw-gruba` | logotyp `DANACO`, tytuły kart środowisk, stopień ekspozycyjny | 30 · 40 px |

Waga 400 **nie jest wgrywana**. Space Grotesk w wadze regularnej traci charakter,
a jego rolą jest właśnie charakter. Jeżeli nagłówek wymaga lekkości — jest to
sygnał, że powinien być złożony IBM Plex Sans, nie że należy dograć czwartą wagę.

### 3.4 Katalog zastosowań — wyczerpujący

Space Grotesk występuje **wyłącznie** w poniższych miejscach:

| # | Miejsce | Stopień | Waga | Uwaga |
|---|---|---:|---:|---|
| 1 | Logotyp `DANACO` | zmienny | 700 | w plikach SVG zamieniony na krzywe |
| 2 | Tytuł karty środowiska (`.dn-karta-srodowiska-tytul`) | 30 px | 700 | TalkIn · WorkSpace · CodeStudio · MultitaskingAI |
| 3 | Nagłówek okna operacyjnego i modala | 20 px | 600 | np. `Workflow Builder`, `Panel prowenancji` |
| 4 | Tytuł sekcji strony głównej | 24 px | 600 | np. „Środowiska", „Komponenty własne" |
| 5 | Nagłówek panelu (`.dn-karta-tytul`) | 16 px | 500–600 | np. „Kolejka zadań", „Historia przebiegów" |
| 6 | Nagłówek strony w dokumentacji HTML | 24–40 px | 600–700 | dokumenty tego pakietu |
| 7 | Stopień ekspozycyjny w materiałach firmowych | 40 px | 700 | OG, banner, tło slajdu |

**Poza tą listą Space Grotesk nie występuje.** W szczególności: nie w treści
przycisku, nie w etykiecie pola, nie w komórce tabeli, nie w tekście wpisu
rozmowy, nie w tooltipie.

### 3.5 Zakazy właściwe krojowi

| Zakaz | Powód |
|---|---|
| **Zakaz stosowania poniżej 16 px** | rysunek geometryczny traci czytelność w małych stopniach; poniżej 16 px obowiązuje Plex Sans |
| **Zakaz w tekście ciągłym** | krój nie jest przeznaczony do dłuższych bloków; akapit złożony Space Grotesk męczy po trzech wierszach |
| **Zakaz wersalikowania nagłówków środowisk** | `TALKIN` zamiast `TalkIn` niszczy nazwę własną i łamie rozdz. 10 kontrakt systemu projektowego |
| **Zakaz odmiany pochyłej** | kursywa nie jest wgrywana; syntetyczna kursywa zablokowana przez `font-synthesis: none` |
| **Zakaz wagi 400 i 800+** | nie wgrywane; przeglądarka nie ma z czego syntezować |
| **Zakaz dodatniego odstępu liter** | dla stopni ≥ 24 px obowiązuje `--dn-ls-naglowek` = −0,01 em; rozstrzelanie rozbija masę wersalików |
| **Zakaz w danych liczbowych** | cyfry proporcjonalne rozjeżdżają kolumny — liczba należy do Plex Mono |

### 3.6 Licencja i jej konsekwencje

Space Grotesk jest objęty **SIL Open Font License 1.1**. Plik licencji:
`zasoby/fonty/LICENCJA-space-grotesk.txt`.

| Wolno | Nie wolno |
|---|---|
| osadzać w aplikacji desktopowej, mobilnej i na stronie | sprzedawać samych plików kroju |
| konwertować do `woff2` i tworzyć podzbiory | rozpowszechniać bez dołączonego tekstu licencji |
| stosować komercyjnie bez opłat i bez zgłoszenia | modyfikować rysunku **zachowując nazwę** (klauzula Reserved Font Name) |
| dołączać do materiałów przekazywanych podwykonawcom | zastrzegać kroju jako elementu znaku towarowego |

**Konsekwencja operacyjna:** każda dystrybucja pakietu (repozytorium, paczka dla
drukarni, materiał dla podwykonawcy) musi zawierać plik licencji. Pliki znaku
mają typografię zamienioną na krzywe — odbiorca nie potrzebuje kroju, żeby
poprawnie wydrukować logotyp.

---

## 4. IBM Plex Sans — krój pracy

### 4.1 Charakterystyka

| Cecha | Wartość |
|---|---|
| Klasyfikacja | grotesk humanistyczny o rysunku inżynierskim |
| Pochodzenie | krój korporacyjny zaprojektowany dla środowiska technologicznego; rodowód inżynierski jest cechą programową, nie skojarzeniem |
| Osobowość | celowo powściągliwa — krój „nie zabiera głosu"; rozpoznawalne pozostają `a`, `g`, `l`, `t` |
| Czytelność w małych stopniach | wysoka wysokość x, otwarte światła wewnętrzne, wyraźne zakończenia — krój utrzymuje czytelność w 11 px |
| Diakrytyki | **narysowane, nie składane** — pełny `latin-ext`, ogonki i kreski o proporcjach dopasowanych do rysunku liter |
| Licencja | **SIL Open Font License 1.1** |
| Wgrywane wagi | 400 · 500 · 600 · 700 (`latin` + `latin-ext`, osiem plików `woff2`, 159 kB) |

### 4.2 Dlaczego krój inżynierski

Gęstość wizualna produktu wynosi **8/10** (kierunek systemu projektowego). Wiersz
tabeli ma 36 px, kontrolka 32 px, stopień bazowy 13 px. W tej gęstości krój
o wyrazistej osobowości staje się szumem: każda litera walczy o uwagę
z sąsiednią liczbą i z ikoną 16-pikselową.

Krój inżynierski robi rzecz przeciwną — **znika**. Operator czyta treść, nie
literę. Jednocześnie rysunek inżynierski niesie właściwą konotację: to jest
narzędzie robocze, nie aplikacja rozrywkowa. Kontrakt kierunku nazywa to
„instrumentem pomiarowym"; Plex Sans jest typograficznym odpowiednikiem tej
metafory.

Drugi powód jest praktyczny: **Plex Sans i Plex Mono pochodzą z jednej rodziny.**
Mają wspólne proporcje, wspólną wysokość x i wspólną linię bazową. Etykieta
złożona Plex Sans i wartość złożona Plex Mono w tym samym wierszu wyrównują się
optycznie bez korekt — a to jest w tym produkcie sytuacja występująca setki razy
na jednym ekranie.

### 4.3 Jakość polskich diakrytyków

Sprawdzenie objęło cztery aspekty krytyczne dla języka polskiego:

| Aspekt | Ocena | Uwaga wykonawcza |
|---|---|---|
| Ogonek w `ą`, `ę`, `Ą`, `Ę` | dobra | ogonek osadzony na osi litery, nie doklejony; nie koliduje z wierszem niższym przy interlinii 1,45 |
| Kreska w `ł`, `Ł` | dobra | kreska przecina laskę pod kątem, zachowuje szerokość znaku; `ł` nie myli się z `t` |
| Akcent w `ó`, `ć`, `ń`, `ś`, `ź` | dobra | akcent krótki i stromy; w stopniu 11 px pozostaje odróżnialny od kropki `ż` |
| Kropka w `ż`, `Ż` | dostateczna | w stopniu 11 px różnica `ź`/`ż` wymaga uwagi — **stąd reguła: nazwy własne i identyfikatory nigdy w 11 px** |
| Kolizja akcentu z linią górną przy `--dn-lh-ciasny` (1,25) | brak | sprawdzone na ciągu `ŁÓDŹ ŻÓŁĆ` w stopniach 20 · 24 · 30 · 40 px |

**Reguła wynikowa [WIĄŻĄCA]:** stopień 11 px (`--dn-fs-xs`) jest zarezerwowany
dla etykiet wersalikowych i nagłówków kolumn — nigdy dla treści zawierającej
nazwy własne, nazwiska ani identyfikatory. Uzasadnienie leży w wierszu `ż`/`ź`
powyżej.

### 4.4 Wagi i ich zastosowanie

| Waga | Żeton | Zastosowanie | Przykład |
|---:|---|---|---|
| **400** | `--dn-fw-normalna` | tekst ciągły, treść wpisu rozmowy, opisy pól, treść tooltipa, opis komponentu | „Automatyka uruchamia się codziennie o 06:00." |
| **500** | `--dn-fw-srednia` | treść przycisku, pozycja bocznej nawigacji, etykieta zakładki, plakietka roli | `Uruchom automatykę` |
| **600** | `--dn-fw-polgruba` | etykieta pola, nagłówek kolumny tabeli, etykieta wersalikowa, tytuł kafla, wyróżnienie w tekście | `HARMONOGRAM` |
| **700** | `--dn-fw-gruba` | **rzadko** — wyłącznie wyróżnienie krytyczne w komunikacie błędu i liczba w karcie wskaźnika | `Przebieg przerwany` |

Waga 700 w Plex Sans jest **wagą wyjątku**. Jeżeli w jednym oknie występuje
więcej niż dwa razy — hierarchia została zbudowana źle i należy ją oprzeć na
stopniu albo na barwie tekstu (`--dn-tekst-2`, `--dn-tekst-3`), nie na wadze.

### 4.5 Katalog zastosowań

Plex Sans jest **krojem domyślnym całego interfejsu**: `body` w `fundament.css`
deklaruje `font-family: var(--dn-ff-bazowa)`. Wszystko, co nie zostało jawnie
przypisane do Space Grotesk (rozdz. 3.4) albo do Plex Mono (rozdz. 5.4), jest
złożone Plex Sans.

Miejsca wymagające jawnej deklaracji (bo dziedziczą po elemencie o innym kroju):

| Miejsce | Uwaga |
|---|---|
| opis pod tytułem karty środowiska | rodzic ma Space Grotesk — opis wymaga przywrócenia `--dn-ff-bazowa` |
| tekst w nagłówku modala obok tytułu | jw. |
| etykieta obok wartości w bloku danych | rodzic bywa `--dn-ff-mono` — etykieta wraca do Plex Sans |
| treść komunikatu w oknie terminalowym | otoczenie mono, komunikat systemowy po polsku — Plex Sans |

### 4.6 Zakazy właściwe krojowi

| Zakaz | Powód |
|---|---|
| **Zakaz w danych liczbowych porównywanych pionowo** | cyfry proporcjonalne — kolumny się rozjeżdżają |
| **Zakaz w identyfikatorach, ścieżkach, skrótach kryptograficznych** | brak rozróżnienia `0`/`O` i `1`/`l`/`I` |
| **Zakaz odmiany pochyłej** | kursywa nie jest wgrywana; wyróżnienie realizuje waga albo barwa |
| **Zakaz podkreślenia poza odnośnikiem** | podkreślenie jest w tym produkcie zarezerwowane dla odnośnika |
| **Zakaz wersalikowania tekstu dłuższego niż trzy wyrazy** | rozdz. 9.4 |
| **Zakaz stopni spoza skali** | wyłącznie dziewięć stopni z rozdz. 7 |

---

## 5. IBM Plex Mono — krój maszyny

### 5.1 Charakterystyka

| Cecha | Wartość |
|---|---|
| Klasyfikacja | krój o stałej szerokości znaku, rodowód wspólny z Plex Sans |
| Osobowość | techniczna, ale nie „terminalowa retro" — brak przerysowanych zakończeń i brak efektu maszyny do pisania |
| Rozróżnialność znaków mylonych | `0` przekreślone ukośnie · `1` z podstawą · `l` z ogonkiem · `I` z belkami — komplet czterech rozróżnień |
| Cyfry | **stałoszerokościowe z natury kroju**; dodatkowo `font-variant-numeric: tabular-nums` jako gwarancja |
| Wersaliki | wąskie, wymagają rozstrzelenia 0,14 em (`--dn-ls-mono-wersaliki`) |
| Licencja | **SIL Open Font License 1.1** (wspólna z Plex Sans) |
| Wgrywane wagi | 400 · 500 · 600 (`latin` + `latin-ext`, sześć plików `woff2`, 86 kB) |

### 5.2 DNA konsoli

Nazwa produktu brzmi **Danaco Console**. Człon `CONSOLE` w logotypie jest złożony
IBM Plex Mono 500 rozstrzelonymi wersalikami — **krój maszyny jest wpisany
w nazwę marki**. To nie jest ozdobnik: para krojów w logotypie opowiada produkt
(marka + maszyna), a ten sam podział rządzi całym interfejsem.

Konsekwencja: **Plex Mono nie jest krojem „technicznym dodatkiem", lecz drugim
krojem tożsamości.** Wszędzie, gdzie system pokazuje własną odpowiedź —
identyfikator, czas, hash, ścieżkę, kod wyjścia, liczbę wywołań — pokazuje ją
krojem, który jest w jego nazwie.

### 5.3 Wagi i ich zastosowanie

| Waga | Żeton | Zastosowanie | Przykład |
|---:|---|---|---|
| **400** | `--dn-fw-normalna` | dane w komórce tabeli, treść okna terminalowego, blok kodu, ścieżka, hash | `sha256:9f2c…a71b` |
| **500** | `--dn-fw-srednia` | `CONSOLE` w logotypie, etykieta wersalikowa mono, plakietka identyfikatora, wartość wyróżniona | `ZADANIE #0412` |
| **600** | `--dn-fw-polgruba` | **rzadko** — liczba kluczowa w karcie wskaźnika, wartość w stanie błędu | `3 nieudane` |

Waga 700 **nie jest wgrywana**. Mono w wadze grubej zlewa się w blok i przestaje
być czytelny w gęstości zwartej.

### 5.4 Katalog zastosowań — wyczerpujący

Plex Mono występuje **wyłącznie** tam, gdzie treść jest **wytworem maszyny albo
adresem maszynowym**:

| # | Kategoria | Przykłady z produktu |
|---|---|---|
| 1 | Identyfikatory | numer zadania w kolejce, identyfikator sesji, identyfikator przebiegu automatyki |
| 2 | Czas i data techniczna | `2026-08-14 06:00`, czas trwania `00:12:04`, znacznik czasu wpisu |
| 3 | Wielkości i liczniki | `1 284 wywołań`, `18,4 MB`, `12 / 15 zadań` |
| 4 | Skróty kryptograficzne | hash Konstytucji w Panelu prowenancji (SHA-256, postać skrócona) |
| 5 | Ścieżki i adresy | ścieżka repozytorium w module Developer, ścieżka pliku w Library Explorer |
| 6 | Nazwy techniczne poleceń | `queue.action`, `enqueue`, `dequeue`, `retry`, `route`, `branch` |
| 7 | Treść okna terminalowego | Terminal Tabs, Output Console, Build Output, Logs Viewer |
| 8 | Blok kodu i kod w wierszu | Code Editor, Studio Editor, Diff/Grep Panel |
| 9 | Nazwy modeli i kanałów | oznaczenia kanałów modelu w Model Configuration |
| 10 | Etykiety wersalikowe kolumn danych | `STATUS`, `CZAS`, `WYWOŁANIA`, `ZASIĘG` |
| 11 | `CONSOLE` w logotypie | jedyne wystąpienie z rozstrzeleniem 0,30 em |

### 5.5 `tabular-nums` — reguła i uzasadnienie

```css
/* Wszystkie dane liczbowe — bezwarunkowo */
.dn-dane,
.dn-tabela td.dn-dane,
.dn-kolejka .dn-krok-meta {
  font-family: var(--dn-ff-mono);
  font-variant-numeric: tabular-nums;
}
```

**Dlaczego mimo stałej szerokości kroju.** Plex Mono ma stałą szerokość znaku
z natury, więc `tabular-nums` jest technicznie nadmiarowe. Deklaracja pozostaje
z dwóch powodów: (1) zabezpiecza sytuację, w której krój nie zdążył się wczytać
i działa krój zapasowy o cyfrach proporcjonalnych; (2) jest **czytelną deklaracją
intencji** dla dewelopera przeglądającego arkusz — ten blok danych ma się
wyrównywać w pionie.

**Skutek praktyczny.** W Execution Monitor kolumna „Czas trwania" z wartościami
`00:00:07`, `00:12:04`, `01:44:19` układa się w prostokąt. Operator porównuje
wartości ruchem oka w pionie, bez czytania — dokładnie tak, jak odczytuje się
wskazania instrumentu.

### 5.6 Zakazy właściwe krojowi

| Zakaz | Powód |
|---|---|
| **Zakaz w tekście ciągłym po polsku** | stała szerokość rozbija rytm; akapit staje się nieczytelny po dwóch wierszach |
| **Zakaz w treści przycisku** | przycisk mówi ludzkim głosem — należy do Plex Sans |
| **Zakaz w nagłówku okna** | nagłówek to tożsamość — należy do Space Grotesk |
| **Zakaz wersalików bez rozstrzelenia** | bez `0,14 em` wersaliki mono zlepiają się w blok |
| **Zakaz wagi 700** | nie wgrywana |
| **Zakaz jako „krój techniczny na wszelki wypadek"** | mono oznacza „to powiedziała maszyna"; użyte do treści ludzkiej kłamie |

---

## 6. Zasada nadrzędna i jej egzekwowanie

### 6.1 Brzmienie zasady

> **Space Grotesk = wejście do środowiska i tożsamość marki.
> IBM Plex Sans = praca.
> IBM Plex Mono = maszyna.**

### 6.2 Test trzech pytań

Przy każdym nowym elemencie tekstowym projektant zadaje trzy pytania w kolejności:

```
   ┌─ 1. Czy ten napis odpowiada na pytanie „gdzie jestem"? ────────┐
   │      (nazwa środowiska, nazwa okna, tytuł sekcji, znak)        │
   │      TAK ──► SPACE GROTESK · ≥ 16 px · waga 500–700            │
   └───────────────────────────┬───────────────────────────────────┘
                               │ NIE
   ┌───────────────────────────▼───────────────────────────────────┐
   │  2. Czy tę treść wytworzyła maszyna albo jest to adres         │
   │     maszynowy? (identyfikator, czas, liczba, ścieżka, kod)     │
   │      TAK ──► IBM PLEX MONO · tabular-nums · waga 400–600       │
   └───────────────────────────┬───────────────────────────────────┘
                               │ NIE
   ┌───────────────────────────▼───────────────────────────────────┐
   │  3. Wszystko pozostałe ──► IBM PLEX SANS · waga 400–600        │
   └───────────────────────────────────────────────────────────────┘
```

Test jest rozstrzygający. Nie ma sytuacji, w której odpowiedź brzmi „zależy".

### 6.3 Przypadki graniczne — rozstrzygnięcia wiążące

| Przypadek | Rozstrzygnięcie | Uzasadnienie |
|---|---|---|
| Nazwa własna okna `Workflow Builder` w nagłówku | **Space Grotesk** | odpowiada na „gdzie jestem" |
| Ta sama nazwa w zdaniu opisowym w tooltipie | **Plex Sans** | zdanie jest treścią, nie tytułem |
| Nazwa automatyki nadana przez Operatora, np. „Raport tygodniowy" | **Plex Sans** | treść ludzka, nie wytwór maszyny |
| Identyfikator tej samej automatyki | **Plex Mono** | adres maszynowy |
| Liczba w zdaniu: „Kolejka zawiera 12 zadań." | **Plex Sans** | liczba wpleciona w zdanie należy do zdania |
| Ta sama liczba w kolumnie tabeli | **Plex Mono** | liczba porównywana pionowo |
| Nazwa akcji silnika kolejek `retry` w opisie | **Plex Mono** | nazwa polecenia technicznego |
| Etykieta przycisku wywołującego tę akcję: „Ponów zadanie" | **Plex Sans** | przycisk mówi ludzkim głosem |
| Nagłówek kolumny `STATUS` w tabeli danych | **Plex Mono** wersaliki 0,14 em | nagłówek kolumny danych należy do bloku danych |
| Etykieta pola formularza `HARMONOGRAM` | **Plex Sans** wersaliki 0,08 em | formularz jest pracą, nie danymi |
| Skrót `SHA-256` w zdaniu | **Plex Mono** | nazwa techniczna |
| Skrót `AOD` w zdaniu | **Plex Sans** | skrót nazwy funkcji produktu, nie nazwa techniczna |

### 6.4 Jak zasadę egzekwuje kod

Egzekwowanie odbywa się **przez żetony i klasy, nie przez dyscyplinę**:

1. `body` deklaruje `--dn-ff-bazowa` — domyślnie wszystko jest Plex Sans.
2. `fundament.css` przypisuje `h1, h2, h3, .dn-naglowek` → `--dn-ff-naglowek`.
   Projektant nie musi pamiętać o Space Grotesk w nagłówku — dostaje go z semantyki.
3. `komponenty.css` przypisuje `--dn-ff-mono` klasom `.dn-dane`, blokom kodu
   i elementom terminalowym. Wpisanie liczby do `.dn-dane` automatycznie daje mono
   z `tabular-nums`.
4. **Deklaracja `font-family` wprost w arkuszu okna jest sygnałem błędu.**
   Jeżeli okno potrzebuje własnej deklaracji — albo używa złego komponentu, albo
   odkryło brakującą klasę, którą należy dodać do `komponenty.css`.

### 6.5 Kontrola przy przeglądzie

| Pytanie kontrolne | Sposób sprawdzenia |
|---|---|
| Czy w oknie występuje Space Grotesk poniżej 16 px? | wyszukanie `--dn-ff-naglowek` w sąsiedztwie `--dn-fs-sm`, `--dn-fs-xs`, `--dn-fs-base` |
| Czy liczba w kolumnie jest złożona mono? | wizualnie: kolumna wyrównana w prostokąt |
| Czy w arkuszu okna występuje `font-family` wprost? | wyszukanie ciągu `font-family` poza `zasoby/` |
| Czy występuje krój spoza trzech? | wyszukanie nazw krojów w plikach okna |
| Czy występuje kursywa? | wyszukanie `font-style: italic` oraz `<em>` bez klasy |

---

## 7. Pełna skala typograficzna — dziewięć stopni

### 7.1 Konstrukcja skali

Skala jest **zwarta u dołu i rozstrzelona u góry**. Sześć stopni mieści się
w przedziale 11–20 px (obszar pracy), trzy pozostałe obsługują ekspozycję.
Taka konstrukcja wynika wprost z gęstości 8/10: w kokpicie różnica 13 → 14 px
niesie realną hierarchię, bo tekst zajmuje całą powierzchnię.

```
  11 ▏
  12 ▎
  13 ▍  ← BAZOWY (--dn-fs-base) · punkt odniesienia całego produktu
  14 ▍
  16 ▌
  20 ▋
  24 ▊
  30 ▉
  40 █
```

### 7.2 Dziewięć stopni — tabela wykonawcza

| Żeton | Stopień | Krój | Waga | Interlinia | Odstęp liter | Zastosowanie |
|---|---:|---|---:|---:|---:|---|
| `--dn-fs-xs` | **11 px** | Plex Sans / Plex Mono | 600 / 500 | 1,25 | 0,08 em / 0,14 em | etykiety wersalikowe, nagłówki kolumn tabel, plakietki drobne, oznaczenia zasięgu |
| `--dn-fs-sm` | **12 px** | Plex Sans / Plex Mono | 400 / 400 | 1,45 | 0 | metadane, godzina wpisu, opis pod polem, tekst pomocniczy, wartości w komórkach danych |
| `--dn-fs-base` | **13 px** | Plex Sans | 400–500 | 1,45 | 0 | **bazowy** — treść przycisku, etykieta pola, pozycja bocznej nawigacji, komórka tabeli, treść tooltipa |
| `--dn-fs-md` | **14 px** | Plex Sans | 400 | 1,45–1,60 | 0 | treść wpisu rozmowy, wyróżniona treść, tekst ciągły w dokumentacji |
| `--dn-fs-lg` | **16 px** | Space Grotesk | 500–600 | 1,25 | −0,01 em | nagłówki paneli, tytuł karty, tytuł sekcji w oknie konfiguracji |
| `--dn-fs-xl` | **20 px** | Space Grotesk | 600 | 1,25 | −0,01 em | nagłówki okien operacyjnych i modali |
| `--dn-fs-2xl` | **24 px** | Space Grotesk | 600 | 1,25 | −0,01 em | tytuły sekcji strony głównej (Centrum dowodzenia) |
| `--dn-fs-3xl` | **30 px** | Space Grotesk | 700 | 1,25 | −0,01 em | tytuły kart środowisk |
| `--dn-fs-display` | **40 px** | Space Grotesk | 700 | 1,25 | −0,01 em | stopień ekspozycyjny: okno startowe, materiały firmowe, okładki dokumentów |

### 7.3 Interlinia — trzy wartości, trzy zastosowania

| Żeton | Wartość | Kiedy | Dlaczego |
|---|---:|---|---|
| `--dn-lh-ciasny` | **1,25** | wszystkie nagłówki (Space Grotesk), etykiety wersalikowe, plakietki | nagłówek jest blokiem — luźna interlinia rozbija go na osobne wiersze |
| `--dn-lh-bazowy` | **1,45** | interfejs: przyciski, pola, tabele, listy, wpisy rozmowy | kompromis gęstości zwartej — jeszcze czytelne, już oszczędne |
| `--dn-lh-luzny` | **1,60** | tekst ciągły dłuższy niż pięć wierszy: dokumentacja, opisy modułów, treść pomocy | akapit potrzebuje powietrza; zastosowanie w oknie roboczym jest wyjątkiem |

**Reguła [WIĄŻĄCA]:** interlinia nie jest wartością dowolną. Trzy wartości
wyczerpują zbiór. Wartość `1,5` albo `normal` w arkuszu okna jest błędem.

### 7.4 Stopnie a dostępność

| Stopień | Ograniczenie barwy tekstu |
|---:|---|
| 11 px | wyłącznie `--dn-tekst` albo `--dn-tekst-2`; **`--dn-tekst-3` zakazany** |
| 12 px | `--dn-tekst`, `--dn-tekst-2`; `--dn-tekst-3` wyłącznie dla metadanych bez znaczenia operacyjnego |
| 13–14 px | wszystkie żetony tekstu zgodnie z `kontrasty.json` |
| ≥ 16 px | wszystkie; przy wadze ≥ 600 i stopniu ≥ 18,66 px obowiązuje próg WCAG dla dużego tekstu |

Podstawa: kontrakt systemu projektowego — `--dn-tekst-3` (#7C7C7C) wyłącznie metadane
i tekst ≥ 18,66 px półgruby.

### 7.5 Skala a punkty łamania

Skala **nie skaluje się płynnie**. Stopnie są stałe we wszystkich punktach
łamania — kokpit wymaga przewidywalności (`WARIANCJA_PROJEKTOWA` = 4/10).
Zmienia się wyłącznie liczba widocznych kolumn i szerokość paneli.

Wyjątek jedyny: w widoku mobilnym (`w1` = 640 px) stopień ekspozycyjny
`--dn-fs-display` (40 px) schodzi do `--dn-fs-3xl` (30 px), ponieważ 40 px
nie mieści nazwy `MultitaskingAI` w jednym wierszu na ekranie telefonu.

---

## 8. Zestawienia typograficzne — pary

Zestawienie typograficzne (para) jest **gotowym rozstrzygnięciem** dla powtarzalnej
sytuacji. Projektant nie składa pary od nowa — sięga po numer z tego rozdziału.

### 8.1 Para 1 · Tytuł + podtytuł

```
┌──────────────────────────────────────────────────┐
│  CodeStudio                                      │  30 px · Space Grotesk 700 · 1,25 · −0,01em
│  Środowisko wytwarzania oprogramowania           │  13 px · Plex Sans 400 · 1,45 · --dn-tekst-2
└──────────────────────────────────────────────────┘
   odstęp między wierszami pary: --dn-od-1 (4 px)
```

| Element | Krój | Stopień | Waga | Barwa | Interlinia |
|---|---|---:|---:|---|---:|
| Tytuł | Space Grotesk | 30 px | 700 | `--dn-tekst` | 1,25 |
| Podtytuł | IBM Plex Sans | 13 px | 400 | `--dn-tekst-2` | 1,45 |

**Zastosowanie:** karty środowisk na stronie głównej, okładki dokumentów,
nagłówek okna startowego.
**Wariant zmniejszony:** tytuł 20 px / 600, podtytuł 12 px — nagłówek okna
operacyjnego z dopiskiem modułu.

### 8.2 Para 2 · Etykieta + wartość

```
┌──────────────────────────────────────────────────┐
│  ZASIĘG KOLEJKI                                  │  11 px · Plex Sans 600 · WERSALIKI 0,08em · --dn-tekst-3
│  Lokalna                                         │  13 px · Plex Sans 500 · --dn-tekst
└──────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────┐
│  IDENTYFIKATOR PRZEBIEGU                         │  11 px · Plex Mono 500 · WERSALIKI 0,14em · --dn-tekst-3
│  0412-A                                          │  13 px · Plex Mono 400 · tabular-nums
└──────────────────────────────────────────────────┘
```

**Rozstrzygnięcie:** o kroju etykiety decyduje krój wartości. Wartość ludzka →
etykieta Plex Sans 0,08 em. Wartość maszynowa → etykieta Plex Mono 0,14 em.
Para nigdy nie miesza krojów etykiety i wartości.

**Zastosowanie:** Panel prowenancji, bloki metadanych w Execution Monitor,
karty wskaźników, podsumowania sesji.

### 8.3 Para 3 · Nagłówek panelu + treść

```
┌──────────────────────────────────────────────────┐
│  Kolejka zadań                        [12]       │  16 px · Space Grotesk 600 · 1,25
├──────────────────────────────────────────────────┤  obrys --dn-obrys
│  Zadania oczekujące na pobranie przez            │  13 px · Plex Sans 400 · 1,45
│  wykonawcę. Kolejność wynika z reguł             │  --dn-tekst-2
│  orkiestracji.                                   │
└──────────────────────────────────────────────────┘
   odstęp nagłówek → treść: --dn-od-3 (12 px)
```

| Element | Krój | Stopień | Waga | Barwa |
|---|---|---:|---:|---|
| Nagłówek panelu | Space Grotesk | 16 px | 600 | `--dn-tekst` |
| Licznik przy nagłówku | Plex Mono | 12 px | 500 | `--dn-tekst-2` |
| Treść | Plex Sans | 13 px | 400 | `--dn-tekst-2` |

**Zastosowanie:** wszystkie panele w oknach operacyjnych — Queue Manager,
Scheduler, Orchestrator, Instructions Panel, Findings Panel.

### 8.4 Para 4 · Wpis rozmowy

```
┌──────────────────────────────────────────────────────────┐
│ ◐  Coordinator            [ROLA]           14:32         │
│    ↑ 13px Sans 600        ↑ 11px Sans 600  ↑ 12px Mono   │
│      --dn-tekst             WERSALIKI        --dn-tekst-3│
│                             0,08em                        │
│                                                           │
│    Zadanie 0412 skierowane do Executor 2. Reguła          │  14 px · Plex Sans 400
│    warunkowa `condition` rozstrzygnęła na ścieżkę         │  interlinia 1,60
│    walidacji.                                             │  --dn-tekst
│         ↑ nazwa akcji: 13px Plex Mono 400, tło --dn-kod   │
└──────────────────────────────────────────────────────────┘
```

| Element | Krój | Stopień | Waga | Barwa |
|---|---|---:|---:|---|
| Nadawca | Plex Sans | 13 px | 600 | `--dn-tekst` |
| Plakietka roli | Plex Sans wersaliki | 11 px | 600 | wg klasy `.dn-plakietka--rola` |
| Godzina | Plex Mono | 12 px | 400 | `--dn-tekst-3` |
| Treść | Plex Sans | 14 px | 400 | `--dn-tekst` |
| Kod w wierszu | Plex Mono | 13 px | 400 | `--dn-tekst`, tło `--dn-r-xs` |

**Uwaga wykonawcza:** treść wpisu jest jedynym miejscem w oknie roboczym, gdzie
obowiązuje interlinia luźna (1,60) — wpis rozmowy jest tekstem ciągłym.

### 8.5 Para 5 · Dane w tabeli

```
╔═══════════════════╤════════════╤═════════════╤═══════════╗
║ ZADANIE           │ STATUS     │ CZAS        │ WYWOŁANIA ║  11px Mono 500 · 0,14em · --dn-tekst-3
╠═══════════════════╪════════════╪═════════════╪═══════════╣
║ 0412-A            │ ● pracuje  │ 00:12:04    │        18 ║  12px Mono 400 · tabular-nums
║ 0413-A            │ ● gotowe   │ 00:00:47    │         3 ║
║ 0414-B            │ ● wstrzym. │ 00:00:00    │         0 ║
╚═══════════════════╧════════════╧═════════════╧═══════════╝
   wiersz 36 px · liczby wyrównane do prawej · reszta do lewej
```

| Element | Krój | Stopień | Waga | Wyrównanie |
|---|---|---:|---:|---|
| Nagłówek kolumny | Plex Mono wersaliki | 11 px | 500 | jak kolumna |
| Identyfikator | Plex Mono | 12 px | 400 | do lewej |
| Status (etykieta) | Plex Sans | 12 px | 400 | do lewej, **zawsze z kropką lub ikoną** |
| Wartość czasu | Plex Mono | 12 px | 400 | do prawej, `tabular-nums` |
| Liczba | Plex Mono | 12 px | 400–600 | **do prawej**, `tabular-nums` |

**Reguła bezwzględna:** status nigdy nie jest wyrażony samym kolorem —
kropka `.dn-kropka` albo ikona towarzyszy etykiecie zawsze (kontrakt systemu projektowego).

### 8.6 Para 6 · Kod

```
┌──────────────────────────────────────────────────┐
│  queue.action                                    │  11 px · Mono 500 · WERSALIKI? NIE — nazwa polecenia
├──────────────────────────────────────────────────┤
│  enqueue    dodanie zadania do kolejki           │  13 px Mono 400 │ 13 px Sans 400
│  dequeue    pobranie kolejnego zadania           │
│  retry      ponowienie po niepowodzeniu          │
└──────────────────────────────────────────────────┘
```

| Element | Krój | Stopień | Waga | Interlinia |
|---|---|---:|---:|---:|
| Kod blokowy | Plex Mono | 13 px | 400 | 1,45 |
| Kod w wierszu tekstu | Plex Mono | 0,92 em rodzica | 400 | dziedziczy |
| Komentarz w kodzie | Plex Mono | jak kod | 400 | `--dn-tekst-3` |
| Opis obok kodu | Plex Sans | 13 px | 400 | 1,45 |

**Uwaga:** kod w wierszu ma stopień **względny** (0,92 em), nie bezwzględny —
inaczej w wierszu 14-pikselowym byłby za mały, a w 11-pikselowym za duży.
Współczynnik 0,92 wyrównuje optycznie wysokość x Plex Mono do Plex Sans.

### 8.7 Tablica zbiorcza par

| # | Para | Krój górny | Krój dolny | Odstęp | Główne zastosowanie |
|---|---|---|---|---|---|
| 1 | tytuł + podtytuł | Space Grotesk 700/30 | Plex Sans 400/13 | 4 px | karty środowisk |
| 2 | etykieta + wartość | Sans/Mono 600/11 wersaliki | Sans/Mono 400–500/13 | 4 px | Panel prowenancji, wskaźniki |
| 3 | nagłówek panelu + treść | Space Grotesk 600/16 | Plex Sans 400/13 | 12 px | wszystkie panele |
| 4 | wpis rozmowy | Plex Sans 600/13 | Plex Sans 400/14 | 8 px | Chat Window |
| 5 | dane w tabeli | Plex Mono 500/11 wersaliki | Plex Mono 400/12 | wiersz 36 px | tabele danych |
| 6 | kod | Plex Mono 500/11 | Plex Mono 400/13 | 8 px | Terminal, Code Editor |

---

## 9. Wersaliki

### 9.1 Kiedy wersaliki

Wersaliki występują w Danaco Console **wyłącznie** w pięciu sytuacjach:

| # | Sytuacja | Krój | Stopień | Odstęp liter |
|---|---|---|---:|---:|
| 1 | Etykieta nad wartością (para 2) | Plex Sans 600 | 11 px | **0,08 em** |
| 2 | Nagłówek kolumny w tabeli danych | Plex Mono 500 | 11 px | **0,14 em** |
| 3 | Plakietka roli i plakietka stanu | Plex Sans 600 | 11 px | **0,08 em** |
| 4 | `CONSOLE` w logotypie | Plex Mono 500 | zmienny | **0,30 em** |
| 5 | Oznaczenie zasięgu / warstwy w konfiguracji | Plex Sans 600 | 11 px | **0,08 em** |

**Poza tymi pięcioma sytuacjami wersaliki nie występują.**

### 9.2 Dlaczego dwa różne odstępy

Wersaliki Plex Sans mają rysunek szerszy i lepiej rozdzielony — 0,08 em
wystarcza, żeby ciąg `HARMONOGRAM` nie zlepił się w blok.

Wersaliki Plex Mono są **wąskie w stosunku do pola znaku** (stała szerokość
wymusza kompromis), przez co bez rozstrzelenia tworzą jednolity pas.
0,14 em przywraca im rytm.

0,30 em w `CONSOLE` jest wartością **wyłącznie logotypową** — nigdy nie pojawia
się w interfejsie. Jej zadaniem jest zrównanie długości członu `CONSOLE`
z członem `DANACO` w układzie pionowym logotypu.

### 9.3 Kompensacja odstępu na końcu ciągu

Odstęp liter dodaje przestrzeń **także po ostatnim znaku**, co przesuwa
wyrównanie do prawej o wartość odstępu. Przy wyrównaniu do prawej albo przy
wyśrodkowaniu wymagana jest kompensacja:

```css
.dn-etykieta--wersaliki {
  letter-spacing: var(--dn-ls-wersaliki);
  margin-right: calc(var(--dn-ls-wersaliki) * -1);
}
```

Przy wyrównaniu do lewej (przypadek domyślny) kompensacja jest zbędna.

### 9.4 Zakaz w treści ciągłej [BEZWZGLĘDNY]

**Wersaliki nie występują w treści ciągłej — nigdy.**

Powody, w kolejności wagi:

1. **Czytelność.** Tekst wersalikowy czyta się o ok. 15 % wolniej, ponieważ
   znika zarys wyrazu (wszystkie wyrazy mają kształt prostokąta).
2. **Diakrytyki.** `ŁÓDŹ ŻÓŁĆ GĘŚLĄ JAŹŃ` — akcenty nad wersalikami kolidują
   z wierszem wyższym przy interlinii 1,25, a przy 1,45 rozbijają rytm bloku.
3. **Głos marki.** Wersaliki w zdaniu czyta się jako krzyk. Głos marki jest
   rzeczowy (rozdz. 14.2) — nie podnosi tonu.
4. **Dostępność.** Czytniki ekranu bywają skłonne odczytywać ciągi wersalikowe
   litera po literze albo jako skrótowce.

**Granica praktyczna:** wersaliki obejmują **maksymalnie trzy wyrazy**.
`ZASIĘG KOLEJKI` — dobrze. `NAKŁADKA I WYWOŁANIE MODELU` — źle; to jest tytuł
obszaru, a więc należy do Space Grotesk w postaci `Nakładka i wywołanie modelu`.

### 9.5 `text-transform` czy zapis wersalikami

**Zawsze `text-transform: uppercase`, nigdy zapis wersalikami w treści.**

| Powód | Wyjaśnienie |
|---|---|
| Tłumaczenie | tekst źródłowy pozostaje w naturalnej formie i daje się przetłumaczyć |
| Czytnik ekranu | odczytuje treść źródłową, nie transformację wizualną |
| Wyszukiwanie | tekst da się znaleźć zwykłym zapytaniem |
| Zmiana decyzji | wycofanie wersalików to zmiana jednej reguły CSS, nie edycja treści |

**Wyjątek jedyny:** akronimy i skróty zapisywane wersalikami z natury
(`SHA-256`, `AOD`, `CLI`, `WCAG`, `OFL`) pozostają wersalikami w treści źródłowej.

---

## 10. Liczby

### 10.1 Zasada ogólna

| Kontekst liczby | Krój | Wyrównanie | `tabular-nums` |
|---|---|---|---|
| Liczba w zdaniu | Plex Sans | wg zdania | nie |
| Liczba w kolumnie tabeli | Plex Mono | **do prawej** | **tak** |
| Liczba w karcie wskaźnika | Plex Mono | do lewej | **tak** |
| Licznik przy etykiecie (np. `[12]`) | Plex Mono | wg etykiety | **tak** |
| Numer kroku w kolejce | Plex Mono | do prawej | **tak** |
| Wartość na pasku postępu | Plex Mono | do prawej | **tak** |

### 10.2 Formaty — tabela wiążąca

| Rodzaj | Format | Przykład | Uwagi |
|---|---|---|---|
| **Identyfikator zadania** | `NNNN-L` (cztery cyfry, myślnik, litera zasięgu) | `0412-A` | zera wiodące zachowane — szerokość stała |
| **Identyfikator sesji** | ciąg alfanumeryczny skrócony do 8 znaków + wielokropek | `k7f2b91c…` | pełną wartość udostępnia kopiowanie |
| **Skrót kryptograficzny** | `algorytm:pierwsze4…ostatnie4` | `sha256:9f2c…a71b` | postać skrócona z opcją skopiowania pełnej |
| **Data** | `RRRR-MM-DD` | `2026-08-14` | format porządkowalny leksykograficznie |
| **Godzina** | `GG:MM` | `06:00` | doba 24-godzinna, bez oznaczeń pory dnia |
| **Data i godzina** | `RRRR-MM-DD GG:MM` | `2026-08-14 06:00` | rozdzielone spacją, bez „T" |
| **Czas trwania** | `GG:MM:SS` | `00:12:04` | zawsze pełne trzy segmenty — kolumna równa |
| **Czas trwania krótki** | `M m S s` | `12 m 04 s` | wyłącznie w zdaniu, nigdy w kolumnie |
| **Rozmiar** | liczba + spacja nierozdzielająca + jednostka | `18,4 MB` | jednostki dziesiętne: kB, MB, GB |
| **Procent** | liczba + spacja nierozdzielająca + `%` | `84 %` | **spacja przed znakiem procentu — norma polska** |
| **Ułamek dziesiętny** | przecinek | `18,4` · `1,75` | **przecinek, nie kropka** |
| **Separator tysięcy** | spacja nierozdzielająca | `1 284` | **nie kropka, nie przecinek** |
| **Zakres** | półpauza bez spacji | `100–220 ms` | półpauza `–`, nie dywiz |
| **Postęp** | `wykonane / wszystkie` | `12 / 15` | spacje wokół ukośnika |
| **Liczba wywołań** | liczba + rzeczownik odmieniony | `1 wywołanie` · `3 wywołania` · `18 wywołań` | rozdz. 21.4 — odmiana przez liczebnik |

### 10.3 Wyrównanie w tabelach — schemat

```
╔════════════╤══════════════╤═══════════╤═════════════╗
║ ZADANIE    │ HARMONOGRAM  │  ROZMIAR  │  WYWOŁANIA  ║
╠════════════╪══════════════╪═══════════╪═════════════╣
║ 0412-A     │ codziennie   │   18,4 MB │       1 284 ║
║ 0413-A     │ pon., śr.    │    2,1 MB │          18 ║
║ 0414-B     │ co godzinę   │  104,7 MB │      12 067 ║
╚════════════╧══════════════╧═══════════╧═════════════╝
   ▲ do lewej   ▲ do lewej     ▲ do prawej  ▲ do prawej
   identyfikator  tekst          wielkość     licznik
```

**Reguła:** wszystko, co się porównuje, wyrównuje się do prawej.
Wszystko, co się czyta, wyrównuje się do lewej.
Nagłówek kolumny dziedziczy wyrównanie kolumny.

### 10.4 Liczba jako element stanu

Liczba nigdy nie stoi sama, gdy niesie stan. Zawsze towarzyszy jej jednostka
albo rzeczownik:

| Źle | Dobrze |
|---|---|
| `3` | `3 nieudane przebiegi` |
| `12/15` | `12 / 15 zadań ukończonych` |
| `84` | `84 %` |
| `18,4` | `18,4 MB` |

---

## 11. Diakrytyki polskie

### 11.1 Ciąg kontrolny

Każdy krój, każda waga i każdy stopień podlegają sprawdzeniu na ciągu:

```
ŁÓDŹ ŻÓŁĆ GĘŚLĄ JAŹŃ
łódź żółć gęślą jaźń
ĄĆĘŁŃÓŚŹŻ      ąćęłńóśźż
0123456789      0O 1lI
« » „ " – — · § № ° %
```

Ciąg został dobrany tak, by wyczerpać wszystkie dziewięć polskich znaków
diakrytycznych w obu wielkościach, plus znaki mylone i interpunkcję polską.

### 11.2 Wyniki sprawdzenia — trzy kroje, wszystkie wagi

| Krój / waga | `ą ę` ogonek | `ł Ł` kreska | `ó ć ń ś ź` akcent | `ż Ż` kropka | Ocena |
|---|---|---|---|---|---|
| Space Grotesk 500 | pełny | pełny | pełny | pełny | **komplet** |
| Space Grotesk 600 | pełny | pełny | pełny | pełny | **komplet** |
| Space Grotesk 700 | pełny | pełny | pełny | pełny | **komplet** |
| IBM Plex Sans 400 | pełny | pełny | pełny | pełny | **komplet** |
| IBM Plex Sans 500 | pełny | pełny | pełny | pełny | **komplet** |
| IBM Plex Sans 600 | pełny | pełny | pełny | pełny | **komplet** |
| IBM Plex Sans 700 | pełny | pełny | pełny | pełny | **komplet** |
| IBM Plex Mono 400 | pełny | pełny | pełny | pełny | **komplet** |
| IBM Plex Mono 500 | pełny | pełny | pełny | pełny | **komplet** |
| IBM Plex Mono 600 | pełny | pełny | pełny | pełny | **komplet** |

Wszystkie dziesięć kombinacji ma pełny zestaw `latin-ext` dostarczony jako
osobny plik `woff2`. Zestawienie plików: rozdz. 12.2.

### 11.3 Punkty ryzyka

| Ryzyko | Objaw | Zabezpieczenie |
|---|---|---|
| **Brak wpięcia podzbioru `latin-ext`** | polskie znaki w kroju systemowym — widoczna zmiana rysunku w środku wyrazu | wszystkie dwanaście deklaracji `latin-ext` obecne w `fonty.css`; kontrola: rozdz. 12.3 |
| **Kolizja akcentu z wierszem wyższym** | `Ó` w nagłówku dotyka wiersza nad nim | interlinia nagłówka nie schodzi poniżej 1,25 |
| **Mylenie `ź` i `ż` w 11 px** | błędny odczyt nazwy własnej | stopień 11 px zarezerwowany dla etykiet, nie dla nazw (rozdz. 4.3) |
| **Syntetyczne pogrubienie** | zniekształcone ogonki `ą`, `ę` | `font-synthesis: none` w `fundament.css` |
| **Wersaliki z akcentem w ciasnym bloku** | `ŻÓŁĆ` traci światło międzywierszowe | wersaliki wyłącznie jednowierszowe (rozdz. 9.4) |

### 11.4 Interpunkcja polska — reguły składu

| Znak | Zastosowanie | Przykład |
|---|---|---|
| Cudzysłów polski `„ "` | cytat, nazwa własna w cudzysłowie | „Delegacja" |
| Cudzysłów `« »` | cytat w cytacie | — |
| Półpauza `–` | zakres liczbowy, myślnik zdaniowy z odstępami | `100–220 ms` · „Kolejka — bez zmian." |
| Pauza `—` | wtrącenie w tekście ciągłym | jak w tym zdaniu |
| Dywiz `-` | wyłącznie w wyrazach złożonych i identyfikatorach | `0412-A` · `SHA-256` |
| Spacja nierozdzielająca | między liczbą a jednostką, po spójniku jednoliterowym | `18,4 MB` · `84 %` |
| Kropka po skrócie | zgodnie z normą | `pon.` · `godz.` · `np.` |
| Wielokropek `…` | skrócenie wartości, wielokropek w tekście | `k7f2b91c…` |

**Zakaz:** cudzysłów prosty `"` w treści interfejsu i dokumentacji.
Występuje wyłącznie w kodzie.

---

## 12. Wcielenie techniczne

### 12.1 Struktura plików

```
zasoby/
├── zetony/
│   ├── fonty.css               ← dwanaście par @font-face (22 deklaracje)
│   ├── zetony.css              ← --dn-ff-* , --dn-fs-* , --dn-fw-* , --dn-lh-* , --dn-ls-*
│   └── fonty/                  ← 22 pliki .woff2  (kopia robocza dla dokumentów)
└── fonty/
    ├── space-grotesk-latin-{500,600,700}-normal.woff2
    ├── space-grotesk-latin-ext-{500,600,700}-normal.woff2
    ├── ibm-plex-sans-latin-{400,500,600,700}-normal.woff2
    ├── ibm-plex-sans-latin-ext-{400,500,600,700}-normal.woff2
    ├── ibm-plex-mono-latin-{400,500,600}-normal.woff2
    ├── ibm-plex-mono-latin-ext-{400,500,600}-normal.woff2
    ├── LICENCJA-space-grotesk.txt
    └── LICENCJA-ibm-plex.txt
```

### 12.2 Inwentarz plików

| Krój | Wagi | Podzbiory | Pliki | Objętość |
|---|---|---|---:|---:|
| Space Grotesk | 500 · 600 · 700 | latin + latin-ext | 6 | ~76 kB |
| IBM Plex Sans | 400 · 500 · 600 · 700 | latin + latin-ext | 8 | ~159 kB |
| IBM Plex Mono | 400 · 500 · 600 | latin + latin-ext | 6 | ~86 kB |
| **Razem** | **10 wag** | — | **20** | **~321 kB** |

Do tego dwa pliki licencji (~9 kB). Objętość całkowita pakietu typograficznego
mieści się poniżej 340 kB — wartość akceptowalna dla aplikacji desktopowej
i dla pierwszego wczytania w przeglądarce.

### 12.3 Wzorzec `@font-face`

Każda para (krój, waga) wymaga **dwóch** deklaracji — `latin` i `latin-ext`.
Pominięcie drugiej powoduje, że polskie znaki wypadają na krój zapasowy.

```css
/* latin — znaki podstawowe */
@font-face {
  font-family: 'IBM Plex Sans';
  font-style: normal;
  font-weight: 400;
  font-display: swap;
  src: url('./fonty/ibm-plex-sans-latin-400-normal.woff2') format('woff2');
  unicode-range: U+0000-00FF, U+0131, U+0152-0153, U+02BB-02BC, U+02C6,
                 U+02DA, U+02DC, U+2000-206F, U+2074, U+20AC, U+2122,
                 U+2191, U+2193, U+2212, U+2215, U+FEFF, U+FFFD;
}

/* latin-ext — polskie diakrytyki: ą ć ę ł ń ó ś ź ż */
@font-face {
  font-family: 'IBM Plex Sans';
  font-style: normal;
  font-weight: 400;
  font-display: swap;
  src: url('./fonty/ibm-plex-sans-latin-ext-400-normal.woff2') format('woff2');
  unicode-range: U+0100-024F, U+0259, U+1E00-1EFF, U+2020, U+20A0-20AB,
                 U+20AD-20CF, U+2113, U+2C60-2C7F, U+A720-A7FF;
}
```

**Rola `unicode-range`:** przeglądarka pobiera plik `latin-ext` **wyłącznie
wtedy, gdy na stronie wystąpi znak z tego zakresu**. Dokument w całości
angielski nie pobierze ani jednego bajtu podzbioru rozszerzonego. Dokument polski
pobiera oba — i to jest sytuacja normalna w tym produkcie.

**Rola `font-display: swap`:** tekst jest widoczny natychmiast w kroju zapasowym
i podmienia się po wczytaniu. Alternatywa (`block`) daje niewidoczny tekst przez
do trzech sekund — w kokpicie niedopuszczalne.

### 12.4 Łańcuchy zapasowe

```css
--dn-ff-naglowek: 'Space Grotesk', 'IBM Plex Sans', 'Segoe UI', system-ui, sans-serif;
--dn-ff-bazowa:   'IBM Plex Sans', 'Segoe UI', system-ui, -apple-system, sans-serif;
--dn-ff-mono:     'IBM Plex Mono', 'Cascadia Mono', 'Consolas', monospace;
```

| Uwaga | Wyjaśnienie |
|---|---|
| Space Grotesk spada na **Plex Sans**, nie na krój systemowy | jeżeli jeden plik się nie wczyta, drugi zwykle tak — dokument zachowuje spójność rodziny |
| Inter **nie występuje** w żadnym łańcuchu | zablokowany jako anty-domyślny (kierunek systemu projektowego) |
| `system-ui` przed `sans-serif` | krój systemowy jest lepszym zapasem niż domyślny szeryfowy |
| Łańcuch mono zaczyna się od krojów o rozróżnialnym `0` | `Cascadia Mono` i `Consolas` mają przekreślone zero |

### 12.5 Wgrywanie wstępne

Dla okien uruchamianych jako pierwsze (okno startowe, okno logowania) wskazane
jest wgranie wstępne dwóch plików o największym udziale:

```html
<link rel="preload" as="font" type="font/woff2" crossorigin
      href="zasoby/zetony/fonty/ibm-plex-sans-latin-400-normal.woff2">
<link rel="preload" as="font" type="font/woff2" crossorigin
      href="zasoby/zetony/fonty/ibm-plex-sans-latin-ext-400-normal.woff2">
```

Wgrywanie wstępne **obejmuje wyłącznie te dwa pliki**. Wgranie wstępne całego
pakietu opóźnia pierwsze wyświetlenie zamiast je przyspieszyć.

### 12.6 Podzbiór dla materiałów firmowych

Materiały o zamkniętej treści (banner, tło slajdu, OG) nie wymagają pełnych
plików kroju. Obowiązuje jedna z dwóch dróg:

| Droga | Kiedy | Uwaga |
|---|---|---|
| **Krzywe** | grafika statyczna, plik przekazywany na zewnątrz | domyślna — pliki znaku mają typografię na krzywych |
| **Podzbiór znakowy** | dokument o treści zmiennej, ale ograniczonym alfabecie | podzbiór wykonywany narzędziem zachowującym `latin-ext` |

**Zakaz:** przekazywanie na zewnątrz plików `woff2` bez towarzyszącego pliku
licencji (rozdz. 3.6).

### 12.7 Reguły globalne w `fundament.css`

```css
body {
  font-family: var(--dn-ff-bazowa);
  font-size: var(--dn-fs-base);         /* 13 px */
  line-height: var(--dn-lh-bazowy);     /* 1,45 */
  font-synthesis: none;                 /* zakaz syntetycznej kursywy i pogrubienia */
  text-rendering: optimizeLegibility;
  -webkit-font-smoothing: antialiased;
}

h1, h2, h3, .dn-naglowek {
  font-family: var(--dn-ff-naglowek);
  font-weight: var(--dn-fw-polgruba);   /* 600 */
  letter-spacing: var(--dn-ls-naglowek);/* −0,01 em */
  line-height: var(--dn-lh-ciasny);     /* 1,25 */
}
```

`font-synthesis: none` jest **regułą bezpieczeństwa typograficznego**: blokuje
tworzenie przez przeglądarkę odmian, których nie wgrano. Bez niej `<em>` dałby
zniekształconą kursywę syntetyczną z uszkodzonymi ogonkami `ą` i `ę`.

---

## 13. Zakazy typograficzne

Lista zamknięta. Każda pozycja jest wiążąca bezwarunkowo.

| # | Zakaz | Uzasadnienie |
|---|---|---|
| 1 | **Krój spoza trzech** | zestaw jest zamknięty; czwarty krój wymaga zgody Właściciela |
| 2 | **Kursywa w jakimkolwiek kroju** | żadna odmiana pochyła nie jest wgrywana; `font-synthesis: none` blokuje syntezę |
| 3 | **Waga spoza wgrywanych** | Space Grotesk 400/800 · Plex Sans 300/800 · Plex Mono 700 — brak plików |
| 4 | **Stopień spoza dziewięciu** | skala jest zamknięta; wartość pośrednia (np. 15 px) jest błędem |
| 5 | **Interlinia spoza trzech** | `1,25` · `1,45` · `1,60`; `normal` jest błędem |
| 6 | **Space Grotesk poniżej 16 px** | traci czytelność |
| 7 | **Plex Mono w tekście ciągłym** | rozbija rytm czytania |
| 8 | **Wersaliki w treści ciągłej** | rozdz. 9.4 |
| 9 | **Wersaliki bez odstępu liter** | zlepiają się w blok |
| 10 | **Odstęp liter dodatni w nagłówku** | rozbija masę Space Grotesk |
| 11 | **Podkreślenie poza odnośnikiem** | podkreślenie jest zarezerwowane |
| 12 | **Przekreślenie jako wyróżnienie** | zarezerwowane dla wartości wycofanej w Diff/Grep Panel |
| 13 | **Cień tekstu** | zakaz bezwzględny — instrument pomiarowy nie ma cieni na literach |
| 14 | **Obrys tekstu (`-webkit-text-stroke`)** | jw. |
| 15 | **Gradient na tekście** | gradienty wyłącznie ilustracyjne (kontrakt systemu projektowego) |
| 16 | **Tekst na obrazie bez warstwy kryjącej** | nieprzewidywalny kontrast |
| 17 | **Justowanie obustronne** | rzeki w tekście przy wąskich kolumnach kokpitu |
| 18 | **Dzielenie wyrazów (`hyphens: auto`)** | dzieli identyfikatory i nazwy własne |
| 19 | **Liczba porównywana bez `tabular-nums`** | kolumna się rozjeżdża |
| 20 | **Cyfry proporcjonalne w tabeli** | jw. |
| 21 | **`font-family` wprost w arkuszu okna** | krój dostarcza żeton albo klasa komponentu |
| 22 | **Wartość stopnia wprost (`font-size: 13px`)** | wyłącznie `var(--dn-fs-*)` |
| 23 | **Zmiana kroju wewnątrz zdania bez powodu semantycznego** | zmiana kroju niesie znaczenie, nie ozdobę |
| 24 | **Wielkie litery w nazwach własnych modułów** | `TALKIN` zamiast `TalkIn` — łamie rozdz. 10 kontrakt systemu projektowego |
| 25 | **Emoji w treści interfejsu** | zakaz z kontraktu systemu projektowego |
| 26 | **Cudzysłów prosty `"` w treści** | obowiązuje cudzysłów polski `„ "` |
| 27 | **Kropka jako separator tysięcy** | obowiązuje spacja nierozdzielająca |
| 28 | **Kropka dziesiętna** | obowiązuje przecinek |

---

# CZĘŚĆ II — GŁOS MARKI

## 14. Charakterystyka głosu — pięć cech

### 14.1 Deklaracja

> **Danaco Console mówi: precyzyjnie · rzeczowo · sprawczo ·
> bez żargonu · bez przechwałek.**

Głos marki jest jeden — ten sam w etykiecie przycisku, w komunikacie błędu,
w dokumentacji i w wiadomości do wsparcia. Nie ma „tonu marketingowego"
i „tonu produktowego"; jest jeden głos w różnych długościach.

### 14.2 Pięć cech — opis i przykład

#### Cecha 1 · Precyzyjny

Nazywa dokładnie tę rzecz, o której mówi. Nie uogólnia, nie zaokrągla,
nie zastępuje nazwy własnej opisem.

| | |
|---|---|
| **Zasada** | Jeżeli produkt ma nazwę na daną rzecz — używamy tej nazwy. Jeżeli istnieje liczba — podajemy liczbę. |
| **Źle** | „Coś poszło nie tak podczas przetwarzania." |
| **Dobrze** | „Przebieg automatyki *Raport tygodniowy* zatrzymał się na kroku 3 z 5." |
| **Dlaczego** | Operator musi wiedzieć **co**, **gdzie** i **na którym etapie**. „Coś" i „gdzieś" nie pozwalają podjąć decyzji. |

#### Cecha 2 · Rzeczowy

Podaje stan rzeczy bez emocji, bez przepraszania i bez podnoszenia tonu.
Zdanie oznajmujące jest formą domyślną.

| | |
|---|---|
| **Zasada** | Zero wykrzykników. Zero słów wartościujących („niestety", „na szczęście", „świetnie"). Zero przeprosin za zachowanie systemu. |
| **Źle** | „Ups! Niestety nie udało się zapisać. Bardzo przepraszamy!" |
| **Dobrze** | „Zmiana nie została zapisana — połączenie z serwerem przerwane. Ponów zapis." |
| **Dlaczego** | Emocja w komunikacie systemu zabiera miejsce informacji i obniża zaufanie. Instrument pomiarowy nie przeprasza — pokazuje wskazanie. |

#### Cecha 3 · Sprawczy

Każdy komunikat kończy się czymś, co Operator może zrobić. Tekst nie zostawia
w stanie zawieszenia.

| | |
|---|---|
| **Zasada** | Komunikat = **stan + następny krok**. Jeżeli następnego kroku nie ma — podajemy, co system zrobi sam i kiedy. |
| **Źle** | „Model niedostępny." |
| **Dobrze** | „Kanał modelu nie odpowiada. Wybierz inny kanał w Model Configuration albo ponów za chwilę — kolejka wznowi się automatycznie." |
| **Dlaczego** | Zasada zero blokad w warstwie językowej: system nigdy nie zostawia Operatora bez drogi dalej. |

#### Cecha 4 · Bez żargonu

Nie używa słownictwa branżowego tam, gdzie istnieje słowo zwykłe. Terminy
techniczne występują wyłącznie wtedy, gdy **są nazwą własną mechanizmu produktu**.

| | |
|---|---|
| **Zasada** | Termin techniczny jest dopuszczalny, jeżeli jest w słowniku (rozdz. 18). Poza słownikiem — słowo zwykłe. |
| **Źle** | „Zainicjuj deployment pipeline'u w kontekście danego tenanta." |
| **Dobrze** | „Uruchom wdrożenie w Deployment Panel." |
| **Dlaczego** | Operator jest zawodowcem, ale nie jest deweloperem tego produktu. Żargon wewnętrzny wyklucza, nie uwiarygodnia. |

#### Cecha 5 · Bez przechwałek

Nie chwali własnego produktu, nie używa przymiotników reklamowych
i nie obiecuje wyników.

| | |
|---|---|
| **Zasada** | Zero słów: „inteligentny", „zaawansowany", „potężny", „rewolucyjny", „bez wysiłku", „automagicznie". Zero obietnic wydajności bez pomiaru. |
| **Źle** | „Nasz zaawansowany silnik orkiestracji inteligentnie przyspieszy Twoją pracę." |
| **Dobrze** | „Silnik kolejek udostępnia jedenaście akcji sterujących obiegiem zadań między rolami." |
| **Dlaczego** | Kokpit jest oceniany po wskazaniach, nie po opisie. Przechwałka w interfejsie narzędzia roboczego czyta się jako brak pewności. |

### 14.3 Czym głos NIE jest

| Nie jest | Wyjaśnienie |
|---|---|
| **oschły** | rzeczowość to nie chłód; „Zapisano." jest rzeczowe i uprzejme jednocześnie |
| **protekcjonalny** | nie tłumaczy Operatorowi rzeczy oczywistych i nie ostrzega przed każdym kliknięciem |
| **żartobliwy** | brak dowcipów, gier słownych i mrugnięć okiem — kokpit nie żartuje w trakcie przebiegu |
| **bezosobowy do granicy nieuprzejmości** | forma bezosobowa dotyczy opisów; akcje mówią wprost do Operatora |
| **anglojęzyczny** | interfejs jest po polsku; angielskie są wyłącznie nazwy własne okien i modułów |

### 14.4 Głos a trzy kroje

| Krój | Odpowiednik w głosie |
|---|---|
| Space Grotesk | **nazwa** — rzeczownik, nazwa własna, bez czasownika: `CodeStudio`, `Panel prowenancji` |
| IBM Plex Sans | **zdanie** — pełne zdanie oznajmujące albo rozkaźnik: „Uruchom automatykę." |
| IBM Plex Mono | **wartość** — sam fakt, bez zdania: `00:12:04`, `sha256:9f2c…a71b` |

---

## 15. Zwrot do odbiorcy

### 15.1 Kim jest odbiorca

**Operator.** Termin jest wiążący w całym produkcie i całej dokumentacji.
Nie „użytkownik", nie „klient", nie „Ty" jako etykieta roli.

| Cecha Operatora | Konsekwencja językowa |
|---|---|
| zawodowiec zarządzający cyfrową organizacją | nie tłumaczymy podstaw; tłumaczymy mechanizmy własne produktu |
| pracuje w kokpicie, nie „korzysta z aplikacji" | czasowniki operacyjne: uruchom, skieruj, wstrzymaj, wznów, zatwierdź |
| podejmuje decyzje na podstawie wskazań | podajemy stan i liczbę, nie ocenę |
| ma pełny dostęp do systemu (zero blokad) | nie mówimy „nie możesz"; mówimy, gdzie się to włącza |

### 15.2 Dwie formy — reguła podziału

| Kontekst | Forma | Przykład |
|---|---|---|
| **Opis** — co produkt robi, czym jest element, jak działa mechanizm | **bezosobowa** | „Kolejka lokalna obejmuje jeden proces orkiestracji." |
| **Akcja** — co Operator ma zrobić teraz | **bezpośredni zwrot, tryb rozkazujący** | „Wybierz zasięg kolejki." |
| **Stan po działaniu** | **bezosobowa, dokonana** | „Zmiana zapisana." |
| **Skutek działania Operatora** | **bezosobowa** | „Zadanie skierowane do Executor 2." |
| **Pytanie w potwierdzeniu** | **bezpośredni zwrot** | „Usunąć projekt *Analiza rynku Q3*?" |

### 15.3 Zakaz formy „my"

Produkt nie mówi o sobie „my". Nie ma w interfejsie zdania „Nie mogliśmy zapisać
zmiany" ani „Wysłaliśmy Ci wiadomość".

| Zamiast | Forma poprawna |
|---|---|
| „Nie mogliśmy zapisać zmiany." | „Zmiana nie została zapisana." |
| „Wysłaliśmy kod na Twój adres." | „Kod wysłany na adres podany przy rejestracji." |
| „Przygotowaliśmy dla Ciebie raport." | „Raport gotowy." |

**Wyjątek jedyny:** korespondencja z działu wsparcia podpisana imieniem i nazwiskiem
oraz dokumenty prawne producenta — tam „my" oznacza spółkę, nie system.

### 15.4 Zakaz zdrobnień i spieszczeń

| Zakazane | Poprawne |
|---|---|
| „chwileczkę", „momencik" | „Trwa przetwarzanie." |
| „okienko" | „okno" |
| „ustawionka" | „ustawienia" |
| „szybciutko" | — (usunąć) |

---

## 16. Zasady pisania w interfejsie

### 16.1 Przycisk — czasownik + dopełnienie

| Reguła | Wartość |
|---|---|
| Konstrukcja | **czasownik w trybie rozkazującym + dopełnienie** |
| Długość | 1–3 wyrazy; **maksymalnie 24 znaki** |
| Wielkość liter | pierwsza wielka, reszta mała (bez wersalików) |
| Kropka | **brak** |
| Krój | Plex Sans 500, stopień 13 px |

| Źle | Dobrze | Powód |
|---|---|---|
| „OK" | „Zapisz zmianę" | „OK" nie mówi, co się stanie |
| „Wyślij" | „Wyślij do walidacji" | brakuje dopełnienia |
| „Kliknij tutaj, aby uruchomić automatykę" | „Uruchom automatykę" | przycisk nie opisuje kliknięcia |
| „URUCHOM" | „Uruchom" | wersaliki zarezerwowane dla etykiet |
| „Uruchomienie procesu" | „Uruchom proces" | rzeczownik odsłowny zamiast czasownika |
| „Anuluj" (w modalu usuwania) | „Zachowaj projekt" | para przycisków ma nazywać dwa skutki, nie skutek i wycofanie |

**Reguła pary przycisków:** w modalu decyzyjnym oba przyciski nazywają **skutek**.
Nie „Usuń / Anuluj", lecz „Usuń projekt / Zachowaj projekt".

### 16.2 Etykieta — rzeczownik

| Reguła | Wartość |
|---|---|
| Konstrukcja | **rzeczownik albo grupa rzeczownikowa**, mianownik |
| Długość | 1–3 wyrazy; **maksymalnie 28 znaków** |
| Zakończenie | bez dwukropka, bez kropki |
| Wersaliki | tak, gdy etykieta stoi **nad** wartością (para 2); nie, gdy stoi **przed** polem |

| Źle | Dobrze |
|---|---|
| „Wpisz nazwę automatyki:" | „Nazwa automatyki" |
| „Jak często ma się uruchamiać?" | „Cykliczność" |
| „Wybierz zasięg" | „Zasięg kolejki" |
| „Model:" | „Model bazowy" |

### 16.3 Komunikat stanu

| Rodzaj stanu | Konstrukcja | Przykład |
|---|---|---|
| **Trwa** | czasownik niedokonany, forma bezosobowa + kropka sygnału | „Trwa pobieranie kontekstu." |
| **Gotowe** | imiesłów bierny dokonany | „Raport gotowy." · „Zmiana zapisana." |
| **Wstrzymane** | imiesłów + powód | „Kolejka wstrzymana przez Coordinatora." |
| **Oczekuje** | czasownik + na co | „Zadanie oczekuje na wolny slot wykonawcy." |
| **Puste** | rzeczownik + brak | rozdz. 20.3 |

**Reguła bezwzględna:** stan nigdy nie jest wyrażony samym kolorem.
Etykieta tekstowa albo ikona towarzyszy zawsze (kontrakt systemu projektowego).

### 16.4 Komunikat błędu — stan + następny krok

Konstrukcja obowiązkowa, trójczłonowa:

```
   [ CO SIĘ STAŁO ]  —  [ DLACZEGO, jeżeli wiadomo ]  ·  [ CO ZROBIĆ TERAZ ]
```

| Człon | Reguła |
|---|---|
| **Co się stało** | zdanie bezosobowe, forma dokonana, nazwa własna elementu |
| **Dlaczego** | wyłącznie jeżeli system zna przyczynę; **nigdy domysł** |
| **Co zrobić** | tryb rozkazujący, jeden konkretny krok albo nazwa okna, w którym się to ustawia |

**Zakaz bezwzględny: nigdy nie obwiniamy Operatora.**

| Źle | Dobrze |
|---|---|
| „Wprowadziłeś nieprawidłowy adres." | „Adres w formacie nierozpoznanym. Sprawdź zapis — oczekiwany format: `nazwa@domena.pl`." |
| „Błąd. Spróbuj ponownie." | „Zapis nieudany — serwer nie odpowiedział w czasie 30 s. Ponów zapis." |
| „Nie masz uprawnień do tej akcji." | rozdz. 17 |
| „Nieprawidłowe dane wejściowe." | „Pole *Cykliczność* wymaga wartości. Wybierz regułę czasową albo pozostaw automatykę uruchamianą ręcznie." |
| „Wystąpił błąd 500." | „Serwer zwrócił błąd wewnętrzny (500). Przebieg zapisany w Logs Viewer — identyfikator `0412-A`." |

### 16.5 Kolejność informacji

Zasada odwróconej piramidy — **najważniejsze najpierw**:

```
   1. co się stało / czym to jest        ← Operator czyta zawsze
   2. czego dotyczy (nazwa własna)       ← Operator czyta zwykle
   3. dlaczego / szczegół techniczny     ← Operator czyta czasem
   4. identyfikator do zgłoszenia        ← Operator kopiuje przy zgłoszeniu
```

Przykład zbudowany według piramidy:

> **Przebieg przerwany.** Automatyka *Raport tygodniowy*, krok 3 z 5.
> Kanał modelu nie odpowiedział w czasie 60 s.
> Identyfikator przebiegu: `0412-A`.

---

## 17. Zasada zero blokad w języku

### 17.1 Reguła

kontrakt systemu projektowego (zasada zero blokad): **Danaco Console nie narzuca twardych blokad.
Izolacja, uprawnienia i profile są wyłącznie opcjami konfigurowalnymi
ze stanem wyjściowym „wyłączone / pełny dostęp".**

W warstwie językowej oznacza to regułę:

> **Nie piszemy o braku uprawnienia. Opisujemy stan bieżący i wskazujemy
> miejsce, w którym Operator włącza albo wyłącza daną opcję.**

Komunikat nigdy nie mówi „nie możesz". Mówi „ta opcja jest teraz wyłączona
i włącza się tutaj" albo „ta akcja jest nieodwracalna — oto co się stanie".

### 17.2 Pary przed / po

| Przed (język blokady) | Po (język stanu i ścieżki) |
|---|---|
| „Brak uprawnień do tej operacji." | „Punkt izolacji *Sieć* jest włączony dla tej roli. Wyłącz go w oknie konfiguracji punktów izolacji, aby wykonać żądanie sieciowe." |
| „Funkcja niedostępna w Twoim planie." | „Kanał modelu nie jest jeszcze skonfigurowany. Dodaj konto w Ustawieniach → Konta modeli." |
| „Nie możesz usunąć tego projektu." | „Projekt *Analiza rynku Q3* zawiera trzy aktywne sesje. Usunięcie zamknie je wszystkie." |
| „Akcja zablokowana." | „Kolejka jest wstrzymana. Wznów kolejkę, aby zadanie zostało pobrane." |
| „Wymagane potwierdzenie administratora." | „Potwierdzenie przed wykonaniem jest włączone dla akcji nieodwracalnych. Wyłącz je w oknie konfiguracji, jeżeli ma się wykonywać od razu." |
| „Dostęp zabroniony." | „Ten agent ma zawężony zakres narzędzi. Rozszerz go w Permissions Center." |
| „Poczekaj 60 s przed ponownym wysłaniem." | „Kod wysłany. Kolejny można wysłać za 60 s — przycisk pozostaje aktywny." |
| „Najpierw uzupełnij wszystkie pola." | „Pole *Nazwa automatyki* jest puste. Automatyka zapisze się pod nazwą roboczą, jeżeli pozostawisz je puste." |

### 17.3 Trzy sytuacje i ich formy

| Sytuacja | Forma językowa | Wzorzec |
|---|---|---|
| **Opcja wyłączona przez Operatora** | opis stanu + miejsce zmiany | „*{Opcja}* jest teraz {stan}. {Czynność} w {nazwa okna}." |
| **Zasób nieskonfigurowany** | opis braku + czynność konfiguracyjna | „{Zasób} nie jest skonfigurowany. {Czynność} w {ścieżka}." |
| **Akcja nieodwracalna** | opis skutku + oba przyciski nazywające skutek | „{Akcja} spowoduje {skutek}." + „{Wykonaj} / {Zachowaj}" |

### 17.4 Czego nie robimy w interfejsie

| Zakaz | Powód |
|---|---|
| **Zakaz atrybutu `disabled`** | przycisk jest zawsze klikalny (kontrakt systemu projektowego) |
| **Zakaz wyszarzania jako komunikatu** | stan nigdy samym kolorem |
| **Zakaz słowa „zablokowane"** w odniesieniu do funkcji | funkcja nie jest zablokowana, tylko wyłączona przez Operatora |
| **Zakaz odliczania blokującego** | odliczanie jest informacyjne, przycisk działa |
| **Zakaz renderowania wyłączonej metody uwierzytelniania** | „brak metody = mniej segmentów, nie zablokowany segment" |

---

## 18. Słownik terminów obowiązujących

Terminy poniżej są **wiążące**. Odmiana podana w kolumnie „Formy" jest jedyną
dopuszczalną. Kolumna „Nie używać" wymienia formy zakazane wraz z powodem.

### 18.1 Podmiot i przestrzenie

| Termin PL | Termin EN / nazwa własna | Definicja | Formy odmiany | Nie używać |
|---|---|---|---|---|
| **Operator** | — | Odbiorca produktu: zawodowiec zarządzający cyfrową organizacją złożoną z modeli, agentów, projektów i procesów. Ma pełny dostęp do systemu. | Operator, Operatora, Operatorowi, Operatorem, Operatorze; lm. Operatorzy | „użytkownik" (zbyt ogólne), „klient" (relacja handlowa), „admin" (żargon), „user" |
| **środowisko** | TalkIn · WorkSpace · CodeStudio · MultitaskingAI | Najwyższy poziom hierarchii — tryb pracy. Odpowiada na pytanie „w jakim trybie pracuję?". Cztery, zamknięta lista. | środowisko, środowiska, środowisku, środowiskiem; lm. środowiska | „przestrzeń", „workspace" (kolizja z modułem Workspace), „obszar", „sekcja" |
| **moduł** | Studio · Research · Library · Translate · Browser · Assistant · Roundtable · Workspace · Automations · Design · Apps · Terminal · Developer · Diagnostics · Agents | Drugi poziom hierarchii — zadanie. „Jakie zadanie wykonuję?". Piętnaście, zamknięta lista. Nazwy własne po angielsku, nieodmienne. | moduł, modułu, modułowi, modułem, module; lm. moduły | „funkcja", „narzędzie", „aplikacja", „zakładka" |
| **okno operacyjne** | — | Trzeci poziom hierarchii — narzędzie. „Jakim narzędziem realizuję?". Nazwy własne okien po angielsku. | okno operacyjne, okna operacyjnego, oknie operacyjnym; lm. okna operacyjne | „ekran", „widok", „panel" (panel jest częścią okna), „strona" |
| **karta sesji** | — | Element pasa kart reprezentujący trwającą sesję pracy. Może nieść pulsującą kropkę sygnału = proces w tle. | karta sesji, karty sesji, kartę sesji; lm. karty sesji | „zakładka" (zakładka to element `.dn-zakladka` wewnątrz okna), „tab" |
| **komponent własny** | Automations · Agents · Workspace · Assistant | Wytwór zbudowany przez Operatora w strefie 2 strony głównej, wpinany w sesję dowolnego modułu. | komponent własny, komponentu własnego, komponentem własnym; lm. komponenty własne | „plugin", „dodatek", „rozszerzenie", „widget" |

### 18.2 Rola, zespół, orkiestracja

| Termin PL | Termin EN / nazwa własna | Definicja | Formy odmiany | Nie używać |
|---|---|---|---|---|
| **rola** | Executor 1 · Executor 2 · Coordinator · Executor 3 / Validator | Pozycja w zespole środowiska MultitaskingAI, z własnym oknem roboczym i przypisanym agentem albo modelem. | rola, roli, rolę, rolą; lm. role | „stanowisko", „bot", „instancja", „worker" |
| **ekspert** | — | Agent o zawężonej, specjalistycznej ekspertyzie zadaniowej, przypisywany do roli. | ekspert, eksperta, ekspertowi, ekspertem; lm. eksperci | „specjalista AI", „asystent" (Assistant to moduł), „ekspertka" jako wariant systemowy |
| **kolejka** | Queue Manager · silnik kolejek | Uporządkowany zbiór zadań o jednym z pięciu zasięgów: globalna, lokalna, dla modelu, dla agenta, dla projektu. | kolejka, kolejki, kolejce, kolejkę, kolejką; lm. kolejki | „lista zadań", „backlog", „pipeline" |
| **silnik kolejek** | `queue.action` | Mechanizm techniczny obsługujący jedenaście akcji: `enqueue`, `dequeue`, `delay`, `retry`, `pause`, `resume`, `split`, `merge`, `route`, `branch`, `condition`. | silnik kolejek, silnika kolejek, silnikiem kolejek | „kolejkowanie", „system kolejek", „broker" |
| **orkiestracja** | Orchestrator | Warstwa egzekwująca zależności między modelami, agentami, zadaniami, kolejkami, automatykami i projektami. | orkiestracja, orkiestracji, orkiestrację, orkiestracją | „zarządzanie przepływem", „workflow" (Workflow Builder to okno), „koordynacja" (Coordinator to rola) |
| **panel orkiestracji** | — | Boczna nawigacja środowiska MultitaskingAI — sześć sekcji: Zespoły, Role, Kolejki, Orkiestracja, Harmonogram i automatyki, Monitor procesu. Zastępuje listę modułów. | panel orkiestracji, panelu orkiestracji, panelem orkiestracji | „menu boczne", „sidebar", „nawigacja modułów" (to inna rzecz) |

### 18.3 Konfiguracja i wywołanie modelu

| Termin PL | Termin EN / nazwa własna | Definicja | Formy odmiany | Nie używać |
|---|---|---|---|---|
| **nakładka** | — | Warstwa programowana ponad modelem bazowym, składana z trzech poziomów: Konstytucja → Profil/Rola → Ekspertyza zadaniowa. Przekazywana modelowi przy każdym wywołaniu. | nakładka, nakładki, nakładce, nakładkę, nakładką | „prompt systemowy" jako nazwa warstwy (to nazwa pola), „wrapper", „preset" |
| **Konstytucja** | — | Najwyższa warstwa nakładki: trwałe, firmowe zasady pracy AI obowiązujące każde wywołanie modelu na platformie. Zasięg: globalny albo per środowisko. Wielka litera — nazwa własna mechanizmu. | Konstytucja, Konstytucji, Konstytucję, Konstytucją | „zasady globalne", „polityka", „reguły firmowe" |
| **prowenancja** | Panel prowenancji | Weryfikowalny zapis tego, co faktycznie zostało przekazane modelowi przy danym wywołaniu — po uwzględnieniu całego dziedziczenia warstw. | prowenancja, prowenancji, prowenancję, prowenancją | „historia", „log", „audyt", „debug" |
| **okno operacyjne** | — | patrz 18.1 | — | — |

### 18.4 Funkcje globalne

| Termin PL | Termin EN / nazwa własna | Definicja | Formy odmiany | Nie używać |
|---|---|---|---|---|
| **Always On Display** | AOD | Funkcja globalna: warstwa centralna z rdzeniem awatara i monitorem procesu, obecna nad całą platformą. Nazwa własna angielska, **nieodmienna**. Skrót `AOD` dopuszczalny po pierwszym rozwinięciu. | nieodmienne: „w Always On Display", „przez Always On Display" | „AOD" bez rozwinięcia przy pierwszym wystąpieniu, „wyświetlacz", „awatar" jako nazwa funkcji, tłumaczenie „Zawsze włączony ekran" |
| **Mobile** | — | Funkcja globalna dostępu mobilnego. Aktywacja nie tworzy nowej przestrzeni roboczej. Nazwa własna, nieodmienna. | nieodmienne: „funkcja Mobile", „w Mobile" | „aplikacja mobilna", „wersja mobilna", „tryb telefonu" |

### 18.5 Odmiana nazw własnych angielskich [REGUŁA]

Nazwy własne środowisk, modułów, okien i ról **są nieodmienne**. Odmienia się
rzeczownik polski stojący przed nimi.

| Źle | Dobrze |
|---|---|
| „w CodeStudiu" | „w środowisku CodeStudio" |
| „z Workflow Buildera" | „z okna Workflow Builder" |
| „Coordinatorowi" | „roli Coordinator" |
| „automatykę w Automationsach" | „automatykę w module Automations" |
| „Executora 2" | „wykonawcy Executor 2" albo „roli Executor 2" |

**Wyjątek dopuszczony:** w tekście technicznym o wysokiej gęstości (opis kroku
w kolejce, log) dopuszcza się skrót bez rzeczownika, jeżeli fraza pozostaje
w mianowniku: „skierowane do Executor 2".

### 18.6 Zapis nazw własnych — wielkość liter

| Kategoria | Zapis | Przykład |
|---|---|---|
| Środowiska | dokładnie jak w dokumentacji | `TalkIn` · `WorkSpace` · `CodeStudio` · `MultitaskingAI` |
| Moduły | pierwsza litera wielka | `Automations` · `Roundtable` · `Diagnostics` |
| Okna operacyjne | każdy człon wielką literą | `Workflow Builder` · `Execution Monitor` · `Chat Window` |
| Role | jak w dokumentacji, z cyfrą | `Executor 1` · `Executor 3 / Validator` |
| Mechanizmy własne PL | wielka litera tylko w nazwie własnej | `Konstytucja` · `Panel prowenancji` · `silnik kolejek` (mała) |
| Akcje techniczne | małe litery, krój mono | `enqueue` · `retry` · `queue.action` |

---

## 19. Terminy zakazane i ich zamienniki

### 19.1 Żargon techniczny

| Zakazane | Zamiennik | Powód |
|---|---|---|
| user, użytkownik | **Operator** | rola nazwana w produkcie |
| deploy, deployment | **wdrożenie** (albo nazwa własna `Deployment Panel`) | istnieje słowo polskie |
| pipeline | **proces** albo **przebieg** | zależnie od kontekstu |
| endpoint | **adres usługi** | — |
| feature | **funkcja** | — |
| bug | **błąd** | — |
| issue | **zgłoszenie** albo **usterka** | — |
| log | **zapis** (okno: `Logs Viewer`) | nazwa okna zostaje |
| debug | **diagnostyka** (moduł: `Diagnostics`) | — |
| token (w znaczeniu żetonu projektowego) | **żeton** | `--dn-*` to żetony |
| token (uwierzytelniający) | **token** | dopuszczalne — termin techniczny bez odpowiednika |
| prompt | **polecenie** (jeżeli treść Operatora) / **prompt systemowy** (nazwa pola) | rozróżnienie obowiązkowe |
| tab | **zakładka** albo **karta sesji** | dwa różne obiekty |
| dashboard | **pulpit** albo nazwa własna okna | `Project Dashboard` zostaje |
| tenant | **organizacja** | — |
| workspace (małą literą) | **przestrzeń robocza** albo nazwa modułu `Workspace` | rozróżnienie obowiązkowe |
| onboarding | **wdrożenie Operatora** | — |
| performance | **wydajność** | — |
| timeout | **przekroczenie czasu** | — |
| retry (w tekście PL) | **ponowienie** (akcja techniczna `retry` zostaje) | — |

### 19.2 Słownictwo reklamowe

| Zakazane | Zamiennik | Powód |
|---|---|---|
| inteligentny, smart | **opis mechanizmu** | „inteligentny" nie niesie informacji |
| zaawansowany | **opis możliwości** | j.w. |
| potężny, wydajny (bez liczby) | **liczba albo pomiar** | przechwałka bez pokrycia |
| rewolucyjny, przełomowy | — (usunąć) | — |
| bez wysiłku, automagicznie | **opis automatyzacji** | — |
| po prostu, wystarczy tylko | — (usunąć) | umniejsza trudność, brzmi protekcjonalnie |
| najlepszy, wiodący | — (usunąć) | twierdzenie bez pomiaru |
| błyskawicznie | **wartość czasu** | „w 1,2 s" zamiast „błyskawicznie" |
| pełna kontrola | **opis zakresu opcji** | — |

### 19.3 Słownictwo blokady

| Zakazane | Zamiennik | Powód |
|---|---|---|
| brak uprawnień | **opis stanu opcji + miejsce zmiany** | zasada zero blokad |
| nie możesz, nie wolno | **opis stanu** | j.w. |
| zablokowane, zabronione | **wyłączone** (o opcji Operatora) | j.w. |
| wymagane, obowiązkowe (o polu) | **opis skutku pozostawienia pustego** | zero blokad |
| niedostępne | **nieskonfigurowane** albo **wyłączone** | — |
| odmowa dostępu | **opis punktu izolacji** | — |

### 19.4 Słownictwo emocjonalne i obwiniające

| Zakazane | Zamiennik |
|---|---|
| Ups!, Ojej! | — (usunąć) |
| niestety, na szczęście | — (usunąć) |
| przepraszamy | — (usunąć; podać stan i krok) |
| wprowadziłeś błędne dane | **opis oczekiwanego formatu** |
| zapomniałeś | **opis stanu pola** |
| coś poszło nie tak | **nazwa elementu + stan** |
| spróbuj ponownie | **„Ponów {czynność}."** |

### 19.5 Zapożyczenia zbędne

| Zakazane | Zamiennik |
|---|---|
| dedykowany (w znaczeniu „przeznaczony") | **właściwy**, **przeznaczony do** |
| adresować problem | **rozwiązać**, **usunąć** |
| implementować | **wdrożyć**, **zbudować** |
| eksportować w znaczeniu „wysłać" | **wyeksportować** (jeżeli to Export Panel) albo **wysłać** |
| aplikować zmiany | **zastosować zmiany** |
| customizacja | **dostosowanie** |
| feedback | **informacja zwrotna** |

---

## 20. Wzorce tekstów

Każdy wzorzec podano w postaci szablonu ze zmiennymi `{…}` oraz trzech przykładów
zbudowanych z **realnych elementów produktu**. Przykłady oznaczone są jako
przykładowe — nie są danymi rzeczywistymi.

### 20.1 Tytuł okna

```
WZORZEC:  {Nazwa własna okna}
          — krój Space Grotesk 600 · 20 px · bez kropki · bez dopisku
```

| # | Moduł | Tytuł |
|---|---|---|
| 1 | Automations | `Workflow Builder` |
| 2 | Diagnostics | `Diagnostics Center` |
| 3 | Agents | `Permissions Center` |

**Wariant z kontekstem** (gdy okno występuje w wielu modułach):
`{Nazwa okna}` + separator `·` + `{Nazwa modułu}` — np. `Chat Window · Terminal`.
Człon kontekstowy: Plex Sans 400, 13 px, `--dn-tekst-3`.

### 20.2 Opis modułu

```
WZORZEC:  {Co robi — zdanie oznajmujące, bezosobowe, czas teraźniejszy}.
          {Dla kogo albo kiedy — zdanie drugie, opcjonalne}.
          — 2 zdania, maksymalnie 180 znaków, bez przymiotników oceniających
```

| # | Moduł | Opis |
|---|---|---|
| 1 | Automations | „Przekształca powtarzalne czynności w procesy działające bez stałego nadzoru. Gotowa automatyka wpina się jako komponent własny w sesję dowolnego modułu." |
| 2 | Roundtable | „Zestawia odpowiedzi kilku modeli w jednym widoku i prowadzi między nimi debatę. Rozstrzygnięcie zapada w Consensus Panel." |
| 3 | Diagnostics | „Zbiera zapisy przebiegów, błędy i wskazania wydajności platformy w jednym miejscu. Rekomendacje pojawiają się w Recommendations Panel." |

### 20.3 Pusty stan

```
WZORZEC:  {Rzeczownik — czego nie ma}.
          {Jedno zdanie: co zrobić, żeby było}.
          [ {Przycisk: czasownik + dopełnienie} ]
          — bez słów „jeszcze", „na razie", „pusto tu"
```

| # | Miejsce | Treść |
|---|---|---|
| 1 | Queue Manager, kolejka bez zadań | „Kolejka pusta." / „Zadania trafiają tutaj po wywołaniu akcji `enqueue` przez Coordinatora albo po wpięciu automatyki." / `[Dodaj zadanie]` |
| 2 | Library Explorer, biblioteka bez plików | „Biblioteka nie zawiera plików." / „Wgraj plik albo przenieś wynik pracy z okna Studio Editor." / `[Wgraj plik]` |
| 3 | Execution Monitor, brak przebiegów | „Brak zapisanych przebiegów." / „Pierwszy przebieg pojawi się po uruchomieniu automatyki albo po nadejściu terminu z harmonogramu." / `[Uruchom teraz]` |

### 20.4 Potwierdzenie

```
WZORZEC:  Nagłówek:  {Czasownik dokonany w bezokoliczniku} {nazwa obiektu}?
          Treść:     {Skutek — co dokładnie się stanie, z liczbą jeżeli jest}.
          Przyciski: [ {Wykonaj — czasownik + obiekt} ] [ {Zachowaj — czasownik + obiekt} ]
          — oba przyciski nazywają skutek; „Anuluj" jest zakazane
```

| # | Sytuacja | Treść |
|---|---|---|
| 1 | Usunięcie projektu (Project Dashboard) | „Usunąć projekt *Analiza rynku Q3*?" / „Projekt zawiera trzy aktywne sesje i dwie wpięte automatyki. Usunięcie zamknie sesje i odepnie automatyki." / `[Usuń projekt]` `[Zachowaj projekt]` |
| 2 | Wstrzymanie kolejki globalnej (Queue Manager) | „Wstrzymać kolejkę globalną?" / „Wstrzymanie zatrzyma pobieranie zadań przez wszystkie role. Zadania w trakcie wykonania dokończą się." / `[Wstrzymaj kolejkę]` `[Pozostaw aktywną]` |
| 3 | Nadpisanie Konstytucji (okno Konfiguracji) | „Zapisać nową treść Konstytucji?" / „Zmiana obejmie każde wywołanie modelu na platformie. Hash Konstytucji zostanie przeliczony." / `[Zapisz Konstytucję]` `[Wróć do edycji]` |

### 20.5 Ostrzeżenie

```
WZORZEC:  {Stan, który wymaga uwagi} — {skutek, jeżeli nic się nie zmieni}.
          {Czynność zapobiegawcza, tryb rozkazujący}.
          — plakietka .dn-plakietka--ostrzezenie + ikona; nigdy sam kolor
```

| # | Miejsce | Treść |
|---|---|---|
| 1 | Scheduler | „Harmonogram nakłada się z automatyką *Kopia biblioteki* — oba przebiegi wypadają o 06:00. Przesuń jeden z terminów albo połącz je w jeden proces." |
| 2 | Model Configuration | „Kanał modelu zbliża się do limitu wywołań: 1 284 z 1 500 w bieżącym oknie rozliczeniowym. Dodaj drugi kanał albo obniż częstotliwość w Scheduler." |
| 3 | Permissions Center | „Agent *Walidator raportów* ma dostęp do wszystkich narzędzi. Zawęź zakres w Permissions Center, jeżeli ma pracować wyłącznie na plikach biblioteki." |

### 20.6 Błąd

```
WZORZEC:  {Co się stało — bezosobowo, dokonanie} — {przyczyna, jeżeli znana}.
          {Co zrobić teraz — tryb rozkazujący}.
          [ {Identyfikator do zgłoszenia, krój mono} ]
```

| # | Miejsce | Treść |
|---|---|---|
| 1 | Execution Monitor | „Przebieg przerwany na kroku 3 z 5 — kanał modelu nie odpowiedział w czasie 60 s. Ponów przebieg albo wybierz inny kanał w Model Configuration. Identyfikator: `0412-A`" |
| 2 | Deployment Panel | „Wdrożenie nieukończone — budowanie zakończone kodem wyjścia `1`. Otwórz Build Output, aby zobaczyć zapis błędu." |
| 3 | okno logowania | „Logowanie nieudane — kod jednorazowy nie pasuje albo wygasł. Wpisz kod ponownie albo wyślij nowy." |

### 20.7 Powiadomienie

```
WZORZEC:  Tytuł: {Rzeczownik + stan — maksymalnie 40 znaków}
          Treść: {Jedno zdanie z nazwą własną i liczbą — maksymalnie 120 znaków}
          — toast .dn-toast; znika po 6 s; treść krytyczna nie znika sama
```

| # | Klasa | Treść |
|---|---|---|
| 1 | `--sukces` | „Automatyka uruchomiona" / „*Raport tygodniowy* — przebieg rozpoczęty o 06:00, pięć kroków w kolejce." |
| 2 | `--informacja` | „Zadanie skierowane" / „Zadanie `0412-A` przekazane do roli Executor 2 przez regułę warunkową." |
| 3 | `--blad` | „Przebieg przerwany" / „*Raport tygodniowy* zatrzymany na kroku 3 z 5. Zapis w Logs Viewer." |

### 20.8 Podpowiedź (tooltip)

```
WZORZEC:  {Co robi ten element — jedno zdanie, maksymalnie 90 znaków}.
          {Skrót klawiaturowy, jeżeli istnieje — krój mono}
          — bez kropki na końcu przy tekście jednozdaniowym krótszym niż 40 znaków
```

| # | Element | Treść |
|---|---|---|
| 1 | Przycisk ikonowy `odswiez` w Queue Manager | „Ponawia wybrane zadanie — akcja `retry`" |
| 2 | Przycisk ikonowy `tarcza` w oknie konfiguracji | „Otwiera punkty izolacji dla tej roli" |
| 3 | Kropka tętniąca na karcie sesji | „Proces działa w tle — sesja pozostaje żywa po zamknięciu okna" |

### 20.9 Objaśnienie kontekstowe `[?]`

```
WZORZEC:  {Czym jest pole — jedno zdanie}.
          {Co się dzieje przy braku wartości}.
          — maksymalnie 200 znaków, Plex Sans 400, 12 px
```

| # | Pole | Treść |
|---|---|---|
| 1 | Konstytucja (okno Konfiguracji) | „Trwałe zasady pracy AI obowiązujące każde wywołanie na platformie albo w tym środowisku. Brak treści oznacza pominięcie tej warstwy." |
| 2 | Zasięg kolejki (panel orkiestracji) | „Określa, których procesów dotyczy kolejka. Brak wskazania oznacza zasięg lokalny — jedna karta sesji." |
| 3 | Cykliczność (Scheduler) | „Reguła czasowa uruchamiania automatyki. Brak reguły oznacza uruchamianie wyłącznie ręczne." |

---

## 21. Mikrokopia

### 21.1 Długości maksymalne — tabela wiążąca

| Element | Maksimum | Optimum | Uwaga |
|---|---:|---:|---|
| Treść przycisku | 24 znaki | 12–18 | mieści się w kontrolce 32 px przy 13 px |
| Treść przycisku `--sm` | 16 znaków | 8–12 | — |
| Etykieta pola | 28 znaków | 12–20 | jeden wiersz |
| Etykieta wersalikowa | 20 znaków | 8–14 | odstęp 0,08 em wydłuża wizualnie o ok. 8 % |
| Nagłówek kolumny | 16 znaków | 6–12 | mono 0,14 em wydłuża o ok. 14 % |
| Pozycja bocznej nawigacji | 22 znaki | 8–16 | boczna ma 224 px |
| Tytuł okna | 40 znaków | 12–28 | Space Grotesk 20 px |
| Tytuł powiadomienia | 40 znaków | 16–30 | toast 280–420 px |
| Treść powiadomienia | 120 znaków | 60–100 | dwa wiersze |
| Tooltip | 90 znaków | 30–70 | jeden lub dwa wiersze |
| Objaśnienie `[?]` | 200 znaków | 80–150 | — |
| Opis pod polem | 120 znaków | 40–90 | — |
| Pusty stan — tytuł | 40 znaków | 14–28 | — |
| Pusty stan — opis | 160 znaków | 60–120 | — |
| Opis modułu | 180 znaków | 100–160 | dwa zdania |
| Komunikat błędu | 220 znaków | 90–170 | z identyfikatorem |

### 21.2 Łamanie wierszy

| Reguła | Zastosowanie |
|---|---|
| **Zakaz dzielenia wyrazów** | `hyphens: none` — dzielenie rozbija identyfikatory i nazwy własne |
| **Zakaz łamania nazwy własnej** | `Workflow Builder` nie łamie się między członami — `white-space: nowrap` na nazwie |
| **Zakaz łamania między liczbą a jednostką** | spacja nierozdzielająca: `18,4 MB` |
| **Zakaz łamania po spójniku jednoliterowym** | spacja nierozdzielająca po `i`, `w`, `z`, `o`, `a`, `u` |
| **Przycięcie zamiast łamania** | w komórce tabeli i na pozycji bocznej nawigacji: `text-overflow: ellipsis` + pełna treść w `title` |
| **Maksymalnie dwa wiersze** | w powiadomieniu, tooltipie i opisie pod polem — `-webkit-line-clamp: 2` |

### 21.3 Skracanie

| Sytuacja | Metoda | Przykład |
|---|---|---|
| Nazwa dłuższa niż pole | wielokropek na końcu | `Analiza rynku Q3 — wersja robo…` |
| Identyfikator długi | pierwsze 8 znaków + wielokropek | `k7f2b91c…` |
| Hash | `algorytm:4…4` | `sha256:9f2c…a71b` |
| Ścieżka długa | wielokropek **w środku**, człon końcowy zachowany | `~/projekty/…/raport.md` |
| Lista wartości | trzy pierwsze + licznik | `Studio, Research, Library +6` |

**Zasada:** skracamy **środek albo koniec**, nigdy początek. Operator rozpoznaje
element po początku ciągu.

### 21.4 Odmiana przez liczebnik

Język polski wymaga trzech form. Reguła wiążąca:

| Liczba | Forma | Przykład |
|---|---|---|
| 1 | mianownik lp. | `1 zadanie` · `1 przebieg` · `1 wywołanie` |
| 2, 3, 4 oraz końcówki 2–4 (poza 12–14) | mianownik lm. | `3 zadania` · `24 przebiegi` · `102 wywołania` |
| 0, 5–21 oraz pozostałe | dopełniacz lm. | `0 zadań` · `12 przebiegów` · `1 284 wywołań` |

```
   liczba % 10 == 1 && liczba % 100 != 11   →  forma 1   (zadanie)
   liczba % 10 ∈ {2,3,4} && liczba % 100 ∉ {12,13,14}  →  forma 2  (zadania)
   pozostałe                                →  forma 3   (zadań)
```

**Zakaz:** konstrukcje omijające odmianę typu `Zadania: 3` zamiast `3 zadania`
są dopuszczalne **wyłącznie** w nagłówku kolumny tabeli, gdzie liczba stoi
w osobnej komórce.

**Zakaz bezwzględny:** zapis `3 zadanie(a)` oraz `3 zadań/zadania`.

### 21.5 Skróty dopuszczone

| Skrót | Rozwinięcie | Kiedy |
|---|---|---|
| `AOD` | Always On Display | po pierwszym rozwinięciu w danym oknie |
| `CLI` | Command Line Interface | w kontekście puli kont Code CLI |
| `SHA-256` | — | nazwa algorytmu, zawsze wersalikami |
| `s`, `ms`, `min`, `godz.` | sekunda, milisekunda, minuta, godzina | przy wartości liczbowej |
| `kB`, `MB`, `GB` | — | przy wartości liczbowej |
| `pon.`, `wt.`, `śr.`, `czw.`, `pt.`, `sob.`, `niedz.` | dni tygodnia | w kolumnie harmonogramu |
| `np.`, `m.in.`, `itd.` | — | wyłącznie w dokumentacji, nie w interfejsie |

**Zakaz w interfejsie:** `itp.`, `ew.`, `tzn.`, `ok.` — pełne słowo albo
przeformułowanie.

---

## 22. Lista kontrolna redakcyjna

### 22.1 Przed oddaniem tekstu interfejsu

**Głos**
- [ ] Zero wykrzykników
- [ ] Zero słów oceniających („niestety", „świetnie", „na szczęście")
- [ ] Zero przeprosin za zachowanie systemu
- [ ] Zero formy „my" (poza korespondencją wsparcia)
- [ ] Zero zdrobnień
- [ ] Zero słownictwa reklamowego (rozdz. 19.2)
- [ ] Zero obwiniania Operatora
- [ ] Każdy komunikat kończy się czymś wykonalnym

**Terminologia**
- [ ] Odbiorca nazwany **Operator**, nie „użytkownik"
- [ ] Nazwy własne środowisk, modułów, okien i ról dokładnie jak w dokumentacji
- [ ] Nazwy własne angielskie **nieodmienione**
- [ ] Terminy ze słownika (rozdz. 18) użyte w formie słownikowej
- [ ] Zero terminów z listy zakazanych (rozdz. 19)
- [ ] `Konstytucja` wielką literą; `silnik kolejek` małą

**Zero blokad**
- [ ] Zero sformułowań „brak uprawnień", „nie możesz", „zablokowane"
- [ ] Każde ograniczenie opisane jako stan opcji + miejsce zmiany
- [ ] Pary przycisków nazywają dwa skutki, nie skutek i „Anuluj"
- [ ] Odliczanie opisane jako informacja, nie jako blokada

**Forma**
- [ ] Przyciski: czasownik + dopełnienie, bez kropki, ≤ 24 znaki
- [ ] Etykiety: rzeczownik, bez dwukropka, ≤ 28 znaków
- [ ] Błędy: stan + przyczyna (jeżeli znana) + krok
- [ ] Puste stany: czego nie ma + co zrobić + przycisk
- [ ] Długości mieszczą się w tabeli 21.1

**Liczby i skład**
- [ ] Przecinek dziesiętny, spacja jako separator tysięcy
- [ ] Spacja nierozdzielająca przed `%` i przed jednostką
- [ ] Odmiana przez liczebnik poprawna (rozdz. 21.4)
- [ ] Cudzysłów polski `„ "`, nie prosty `"`
- [ ] Półpauza w zakresach, dywiz wyłącznie w wyrazach złożonych
- [ ] Daty w formacie `RRRR-MM-DD`, godziny `GG:MM`

**Typografia**
- [ ] Krój zgodny z testem trzech pytań (rozdz. 6.2)
- [ ] Liczby porównywane pionowo w Plex Mono z `tabular-nums`
- [ ] Wersaliki wyłącznie w pięciu dozwolonych sytuacjach, ≤ 3 wyrazy
- [ ] Zero kursywy, zero podkreśleń poza odnośnikiem
- [ ] Stopnie i interlinie wyłącznie ze skali
- [ ] Polskie diakrytyki poprawne w całym pliku

### 22.2 Test szybki — trzy pytania do każdego zdania

```
   1. Czy Operator wie po tym zdaniu, CO SIĘ STAŁO?          → nie: dopisz stan
   2. Czy Operator wie, CO ZROBIĆ DALEJ?                     → nie: dopisz krok
   3. Czy da się to skrócić bez utraty żadnej z odpowiedzi?  → tak: skróć
```

### 22.3 Test głosu — czytanie na głos

Zdanie przechodzi, jeżeli po przeczytaniu na głos:

- nie brzmi jak reklama,
- nie brzmi jak przeprosiny,
- nie brzmi jak instrukcja dla dziecka,
- brzmi jak zdanie wypowiedziane przez zawodowca do zawodowca.

---

## 23. Decyzje projektowe

| # | Decyzja | Uzasadnienie | Alternatywa odrzucona |
|---|---|---|---|
| **1** | Typografia i głos w jednym dokumencie | litera i słowo są jednym przedmiotem projektowym; rozdzielenie prowadzi do rozjazdu formy i treści | dwa osobne dokumenty |
| **2** | Space Grotesk **nie schodzi poniżej 16 px** | rysunek geometryczny traci czytelność; Plex Sans przejmuje wszystko poniżej | dopuszczenie 13–14 px dla „małych nagłówków" |
| **3** | Waga 400 Space Grotesk **nie jest wgrywana** | krój ma nieść charakter; waga regularna go traci, a każdy plik kosztuje | komplet czterech wag |
| **4** | Waga 700 Plex Mono **nie jest wgrywana** | mono w wadze grubej zlewa się w blok przy gęstości zwartej | komplet czterech wag |
| **5** | O kroju **etykiety** decyduje krój **wartości** | para etykieta+wartość jest jednym blokiem semantycznym; mieszanie krojów w parze rozbija go | etykieta zawsze Plex Sans |
| **6** | Nagłówek kolumny danych = Plex Mono 0,14 em, etykieta formularza = Plex Sans 0,08 em | tabela danych należy do świata maszyny, formularz do świata pracy | jednolity krój etykiet wersalikowych |
| **7** | `tabular-nums` deklarowane mimo stałej szerokości Plex Mono | zabezpieczenie na czas przed wczytaniem kroju + czytelna deklaracja intencji | poleganie na naturze kroju |
| **8** | Stopień **11 px zarezerwowany** dla etykiet, nigdy dla nazw własnych | rozróżnienie `ź`/`ż` w 11 px wymaga uwagi; nazwa własna nie może być odczytana błędnie | dopuszczenie 11 px dla metadanych z nazwami |
| **9** | Skala **nie skaluje się płynnie** przy punktach łamania | kokpit wymaga przewidywalności (`WARIANCJA_PROJEKTOWA` 4/10) | typografia płynna `clamp()` |
| **10** | Wyjątek: `--dn-fs-display` → `--dn-fs-3xl` w widoku mobilnym | `MultitaskingAI` nie mieści się w 40 px na ekranie telefonu | łamanie nazwy środowiska |
| **11** | Kod w wierszu ma stopień **względny** 0,92 em | stopień bezwzględny byłby za mały w wierszu 14 px i za duży w 11 px | 12 px na sztywno |
| **12** | **Inter nieobecny także w łańcuchu zapasowym** | anty-domyślne kierunek systemu projektowego blokuje odruch, nie tylko deklarację główną | Inter jako pierwszy zapas |
| **13** | Space Grotesk spada na **Plex Sans**, nie na krój systemowy | zachowanie spójności rodziny przy częściowym niepowodzeniu wczytania | `system-ui` bezpośrednio |
| **14** | Odbiorca nazwany **Operator** w interfejsie, nie tylko w dokumentacji | nazwa roli jest częścią pozycjonowania produktu; „użytkownik" należy do innej klasy narzędzi | „użytkownik" w interfejsie, „Operator" w dokumentacji |
| **15** | Nazwy własne angielskie **nieodmienne** | odmiana `w CodeStudiu` niszczy nazwę i łamie kontrakt systemu projektowego | odmiana fonetyczna |
| **16** | Para przycisków w modalu **nazywa dwa skutki** | „Anuluj" nie jest skutkiem; Operator ma wybierać między dwoma stanami świata | „Usuń / Anuluj" |
| **17** | Zero komunikatów „brak uprawnień" — **opis stanu i ścieżki** | bezpośrednia konsekwencja zasady zero blokad w warstwie językowej | standardowe komunikaty odmowy |
| **18** | Zakaz formy „my" w interfejsie | system nie jest osobą; „my" wprowadza fikcyjny podmiot między Operatora a maszynę | „Nie mogliśmy zapisać…" |
| **19** | Wersaliki **wyłącznie przez `text-transform`** | tłumaczenie, czytniki ekranu, wyszukiwanie, odwracalność decyzji | zapis wersalikami w treści |
| **20** | Odmiana przez liczebnik **obowiązkowa**, zapis `3 zadanie(a)` zakazany | poprawność językowa jest częścią precyzji głosu | forma z nawiasem |
| **21** | Skracanie **środka albo końca**, nigdy początku | Operator rozpoznaje element po początku ciągu | skracanie od lewej |
| **22** | `font-synthesis: none` jako reguła globalna | blokuje syntetyczną kursywę i pogrubienie, które niszczą ogonki `ą`, `ę` | domyślne zachowanie przeglądarki |
| **23** | Wgrywanie wstępne obejmuje **dwa pliki**, nie cały pakiet | wgranie wstępne całości opóźnia pierwsze wyświetlenie | preload wszystkich 20 plików |
| **24** | Słownik podaje **formy odmiany**, nie tylko definicje | bez form odmiany słownik nie rozstrzyga sporu redakcyjnego | słownik definicyjny |

---

## Załącznik A · Karta typograficzna na jedną stronę

```
┌──────────────────────────────────────────────────────────────────────────┐
│  DANACO CONSOLE — KARTA TYPOGRAFICZNA                             v2.0   │
├──────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  TRZY KROJE                                                              │
│  Space Grotesk  500·600·700   nagłówki, tytuły środowisk, DANACO         │
│  IBM Plex Sans  400·500·600·700   interfejs i treść                      │
│  IBM Plex Mono  400·500·600   dane, identyfikatory, CONSOLE              │
│  Wszystkie SIL OFL 1.1 · latin + latin-ext · woff2 · ~321 kB             │
│                                                                          │
│  DZIEWIĘĆ STOPNI                                                         │
│  11 · 12 · [13] · 14 · 16 · 20 · 24 · 30 · 40 px      [13] = bazowy      │
│                                                                          │
│  INTERLINIA        1,25 nagłówki · 1,45 interfejs · 1,60 treść ciągła    │
│  ODSTĘP LITER      −0,01 em nagłówek · 0,08 em wersaliki Sans ·          │
│                     0,14 em wersaliki Mono · 0,30 em CONSOLE (logotyp)   │
│  LICZBY            tabular-nums · przecinek dziesiętny · spacja tysięcy  │
│                                                                          │
│  TEST TRZECH PYTAŃ                                                       │
│  1. „gdzie jestem?"  → Space Grotesk ≥ 16 px                             │
│  2. „to wytwór maszyny?" → IBM Plex Mono                                 │
│  3. wszystko inne    → IBM Plex Sans                                     │
│                                                                          │
│  DIAKRYTYKI        ŁÓDŹ ŻÓŁĆ GĘŚLĄ JAŹŃ · łódź żółć gęślą jaźń          │
│                                                                          │
│  GŁOS   precyzyjny · rzeczowy · sprawczy · bez żargonu · bez przechwałek │
│  ODBIORCA          Operator                                              │
│  BŁĄD              stan + przyczyna + następny krok · nigdy wina         │
│  BLOKADA           nie istnieje — opis stanu i miejsca zmiany            │
│                                                                          │
│  ZAKAZY  kursywa · krój spoza trzech · wersaliki w treści · cień tekstu  │
│          gradient na tekście · justowanie · dzielenie wyrazów ·          │
│          „brak uprawnień" · „Ups!" · „użytkownik" · disabled             │
└──────────────────────────────────────────────────────────────────────────┘
```

---

## Załącznik B · Ciąg kontrolny redakcyjny

Ciąg do wklejenia w dowolne okno w celu sprawdzenia kompletu reguł
typograficznych i redakcyjnych naraz:

```
ŁÓDŹ ŻÓŁĆ GĘŚLĄ JAŹŃ · łódź żółć gęślą jaźń · ĄĆĘŁŃÓŚŹŻ ąćęłńóśźż
0123456789 · 0O 1lI · 1 284 wywołań · 18,4 MB · 84 % · 100–220 ms
2026-08-14 06:00 · 00:12:04 · 0412-A · sha256:9f2c…a71b · queue.action
CodeStudio · Workflow Builder · Executor 3 / Validator · Konstytucja
ZASIĘG KOLEJKI · STATUS · WYWOŁANIA · DANACO C O N S O L E
1 zadanie · 3 zadania · 12 zadań · „Delegacja" · Operator
```

Sprawdzeniu podlega:

| Wiersz | Co sprawdza |
|---|---|
| 1 | komplet `latin-ext` we wszystkich krojach i wagach |
| 2 | rozróżnialność `0`/`O` i `1`/`l`/`I`, separator tysięcy, przecinek dziesiętny, spacja przed `%`, półpauza w zakresie |
| 3 | formaty daty, czasu trwania, identyfikatora, hasha, nazwy polecenia |
| 4 | zapis nazw własnych, brak odmiany, wielkie litery |
| 5 | wersaliki w dwóch odstępach, rozstrzelenie logotypowe |
| 6 | odmiana przez liczebnik, cudzysłów polski, nazwa odbiorcy |

---

## Załącznik C · Mapa żetonów typograficznych

| Żeton | Wartość | Plik źródłowy | Rozdział |
|---|---|---|---|
| `--dn-ff-naglowek` | `'Space Grotesk', 'IBM Plex Sans', 'Segoe UI', system-ui, sans-serif` | `zetony.css` | 3, 12.4 |
| `--dn-ff-bazowa` | `'IBM Plex Sans', 'Segoe UI', system-ui, -apple-system, sans-serif` | `zetony.css` | 4, 12.4 |
| `--dn-ff-mono` | `'IBM Plex Mono', 'Cascadia Mono', 'Consolas', monospace` | `zetony.css` | 5, 12.4 |
| `--dn-fs-xs` | `11px` | `zetony.css` | 7.2 |
| `--dn-fs-sm` | `12px` | `zetony.css` | 7.2 |
| `--dn-fs-base` | `13px` | `zetony.css` | 7.2 |
| `--dn-fs-md` | `14px` | `zetony.css` | 7.2 |
| `--dn-fs-lg` | `16px` | `zetony.css` | 7.2 |
| `--dn-fs-xl` | `20px` | `zetony.css` | 7.2 |
| `--dn-fs-2xl` | `24px` | `zetony.css` | 7.2 |
| `--dn-fs-3xl` | `30px` | `zetony.css` | 7.2 |
| `--dn-fs-display` | `40px` | `zetony.css` | 7.2 |
| `--dn-fw-normalna` | `400` | `zetony.css` | 4.4, 5.3 |
| `--dn-fw-srednia` | `500` | `zetony.css` | 3.3, 4.4, 5.3 |
| `--dn-fw-polgruba` | `600` | `zetony.css` | 3.3, 4.4, 5.3 |
| `--dn-fw-gruba` | `700` | `zetony.css` | 3.3, 4.4 |
| `--dn-lh-ciasny` | `1.25` | `zetony.css` | 7.3 |
| `--dn-lh-bazowy` | `1.45` | `zetony.css` | 7.3 |
| `--dn-lh-luzny` | `1.6` | `zetony.css` | 7.3 |
| `--dn-ls-naglowek` | `-0.01em` | `zetony.css` | 7.2 |
| `--dn-ls-wersaliki` | `0.08em` | `zetony.css` | 9.1 |
| `--dn-ls-mono-wersaliki` | `0.14em` | `zetony.css` | 9.1 |

---

*Danaco Console — AI Operating Environment · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. — Dariusz Naharnowicz*
