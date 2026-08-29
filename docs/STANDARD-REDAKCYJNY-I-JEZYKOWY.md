# Danaco Console — Standard redakcyjny i językowy

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
| **Tytuł** | Standard redakcyjny i językowy zbioru dokumentacji oraz konwencje dalszej budowy |
| **Klasa dokumentu** | Nawigacja |
| **Odbiorcy** | redaktor dokumentacji · projektant · deweloper |
| **Przeznaczenie** | Ustala jeden obowiązujący kształt każdego opracowania zbioru, jeden słownik terminologiczny oraz konwencje nazewnicze wiążące w dalszej budowie produktu. Na jego podstawie powstaje i jest przyjmowane każde opracowanie. |
| **Zakres** | budowa dokumentu, wymóg form wizualnych, wymóg normatywnej szczegółowości, język i typografia, konwencje nazewnicze bytów kodu, tryb postępowania przy braku ustalenia |
| **Poza zakresem** | treść merytoryczna produktu — nosi ją zbiór opracowań wskazany w [Spisie opracowań](SPIS-OPRACOWAN.md) |
| **Dokument nadrzędny** | [Spis opracowań](SPIS-OPRACOWAN.md) |
| **Dokumenty powiązane** | [Opis produktu](README.md) · [Licencja produktu](LICENSE.md) |
| **Prototypy odniesienia** | `design/05-okna/` — komplet 38 prototypów |
| **Źródła normatywne** | `budowa/shared/contract.json` · `design/zasoby/zetony/zetony.css` · `design/zasoby/css/komponenty.css` · `design/zasoby/rama.css` · `design/zasoby/prototyp.css` · `design/zasoby/ikony/manifest.json` · `budowa/shared/kontrasty-progi.json` |
| **Zasada nadrzędna** | Opracowanie niesie zamknięty zbiór informacji. Jeśli dokument czegoś nie rozstrzyga, wykonawca nie rozstrzyga tego sam — dokument oznacza brak i stawia pytanie. |

---

## Spis treści

1. [Do czego służy ten standard](#1-do-czego-służy-ten-standard)
   - [1.1 Trzy części standardu](#11-trzy-części-standardu)
   - [1.2 Moc obowiązująca](#12-moc-obowiązująca)
2. [Budowa opracowania](#2-budowa-opracowania)
   - [2.1 Nagłówek i metryka](#21-nagłówek-i-metryka)
   - [2.2 Klasy dokumentów](#22-klasy-dokumentów)
   - [2.3 Spis treści](#23-spis-treści)
   - [2.4 Korpus dokumentu](#24-korpus-dokumentu)
   - [2.5 Stopka](#25-stopka)
   - [2.6 Oznaczenia stanu wdrożenia](#26-oznaczenia-stanu-wdrożenia)
3. [Wymóg form wizualnych](#3-wymóg-form-wizualnych)
   - [3.1 Miara i próg](#31-miara-i-próg)
   - [3.2 Repertuar form obowiązkowych](#32-repertuar-form-obowiązkowych)
   - [3.3 Jakość diagramów](#33-jakość-diagramów)
4. [Wymóg normatywnej szczegółowości](#4-wymóg-normatywnej-szczegółowości)
   - [4.1 Źródła normatywne](#41-źródła-normatywne)
   - [4.2 Wykazy obowiązkowe](#42-wykazy-obowiązkowe)
   - [4.3 Wykluczenie luk wykończeniowych](#43-wykluczenie-luk-wykończeniowych)
5. [Język](#5-język)
   - [5.1 Słownik wiążący](#51-słownik-wiążący)
   - [5.2 Terminy zachowane](#52-terminy-zachowane)
   - [5.3 Zasady polszczyzny](#53-zasady-polszczyzny)
   - [5.4 Typografia](#54-typografia)
6. [Odsyłacze i nawigacja](#6-odsyłacze-i-nawigacja)
7. [Konwencje budowy obowiązujące na przyszłość](#7-konwencje-budowy-obowiązujące-na-przyszłość)
   - [7.1 Zakaz kodów wymyślonych](#71-zakaz-kodów-wymyślonych)
   - [7.2 Zakaz określeń abstrakcyjnych](#72-zakaz-określeń-abstrakcyjnych)
   - [7.3 Nazewnictwo bytów kodu](#73-nazewnictwo-bytów-kodu)
   - [7.4 Reguła jednego źródła prawdy](#74-reguła-jednego-źródła-prawdy)
   - [7.5 Reguła zamkniętego zbioru](#75-reguła-zamkniętego-zbioru)
   - [7.6 Tryb postępowania przy braku ustalenia](#76-tryb-postępowania-przy-braku-ustalenia)
8. [Kontrola zgodności](#8-kontrola-zgodności)
9. [Progi objętości](#9-progi-objętości)
10. [Wzorzec w pełni wypełniony](#10-wzorzec-w-pełni-wypełniony)
11. [Najczęstsze usterki redakcyjne](#11-najczęstsze-usterki-redakcyjne)
12. [Polecenia kontrolne dla redaktora](#12-polecenia-kontrolne-dla-redaktora)

---

## 1. Do czego służy ten standard

Zbiór dokumentacji Danaco Console liczy kilkadziesiąt opracowań pisanych w różnym czasie
i przez różne ręce. Bez jednego wzorca każde kolejne przejście redakcyjne poprawia to samo
od nowa: raz metrykę, raz stopkę, raz termin zapisany na trzy sposoby. Standard zamyka tę
pętlę — ustala kształt dokumentu, słownik i konwencje nazewnicze raz, wiążąco i sprawdzalnie.

### 1.1 Trzy części standardu

```
┌─────────────────────────────────────────────────────────────────────┐
│  STANDARD REDAKCYJNY I JĘZYKOWY                                     │
├──────────────────┬──────────────────────┬───────────────────────────┤
│  CZĘŚĆ PIERWSZA  │  CZĘŚĆ DRUGA         │  CZĘŚĆ TRZECIA            │
│  Redakcja        │  Język               │  Konwencje budowy         │
├──────────────────┼──────────────────────┼───────────────────────────┤
│  rozdz. 2 · 3    │  rozdz. 5            │  rozdz. 7                 │
│  rozdz. 4 · 6    │                      │                           │
├──────────────────┼──────────────────────┼───────────────────────────┤
│  jak dokument    │  jakimi słowami      │  jak nazywać byty,        │
│  ma być zbudowany│  ma być napisany     │  które dopiero powstaną   │
└──────────────────┴──────────────────────┴───────────────────────────┘
```

Legenda: część pierwsza dotyczy formy opracowania, druga jego języka, trzecia — tego, co
z opracowania przechodzi do kodu i do przyszłych dokumentów.

### 1.2 Moc obowiązująca

| Adresat | Zakres związania |
|---|---|
| Redaktor dokumentacji | całość standardu; opracowanie niezgodne ze standardem nie jest przyjmowane |
| Projektant | rozdziały 2–6 przy opisie interfejsu; rozdział 7 przy nazywaniu komponentów i żetonów |
| Deweloper | rozdział 7 w całości; rozdziały 4 i 5 przy nazywaniu komend, zdarzeń i etykiet |

---

## 2. Budowa opracowania

Każde opracowanie zbioru składa się z sześciu części w stałej kolejności:

```
┌───────────────────────────┐
│  1 · Tytuł (nagłówek H1)  │
├───────────────────────────┤
│  2 · Metryka produktowa   │  tabela stała, osiem wierszy
├───────────────────────────┤
│  3 · Metryka dokumentu    │  tabela zmienna, jedenaście wierszy
├───────────────────────────┤
│  4 · Spis treści          │  odsyłacze do wszystkich rozdziałów
├───────────────────────────┤
│  5 · Korpus               │  rozdziały numerowane, oddzielane ───
├───────────────────────────┤
│  6 · Stopka               │  trzy bloki, treść dosłowna
└───────────────────────────┘
```

### 2.1 Nagłówek i metryka

Tytuł dokumentu zapisujemy jako `# Danaco Console — {Tytuł opracowania}`. Bezpośrednio pod
nim stoją dwie tabele.

**Metryka produktowa** — wartości stałe dla całego zbioru, zmienia się wyłącznie data:

```markdown
| | |
|---|---|
| **Produkt** | Danaco Console |
| **Rodzaj** | Platforma AI Workspace OS |
| **Opis** | Platforma jest wielośrodowiskowym systemem operacyjnym dla sztucznej inteligencji, integrującym komunikację, zarządzanie wiedzą, tworzenie treści, projektowanie, automatyzacje procesów oraz rozwój oprogramowania. |
| **Producent** | Danaco Holding Group Sp. z o.o. |
| **Twórca** | Dariusz Naharnowicz |
| **Wersja** | v2.0 |
| **Status** | Deweloperski |
| **Data** | {RRRR-MM-DD} |
```

Wiersz **Opis** niesie brzmienie dosłowne podane wyżej i jest jednakowy w całym zbiorze —
jak pozostałe wiersze tej tabeli, nie podlega przeredagowaniu w pojedynczym opracowaniu.

**Metryka dokumentu** — wartości zmienne, poprzedzona wierszem
`**Informacje szczegółowe dokumentu:**`:

| Pole | Obowiązkowe | Treść |
|---|---|---|
| **Tytuł** | tak | pełny tytuł opracowania, bez przedrostka nazwy produktu |
| **Klasa dokumentu** | tak | jedna z czterech wartości z rozdz. 2.2 |
| **Odbiorcy** | tak | wskazani imiennie: deweloper, projektant, Operator — nie „wszyscy zainteresowani” |
| **Przeznaczenie** | tak | jedno zdanie: do czego dokument służy i co na jego podstawie powstaje |
| **Zakres** | tak | co dokument obejmuje |
| **Poza zakresem** | tak | co świadomie pominięto i gdzie tego szukać |
| **Dokument nadrzędny** | tak | odsyłacz do opracowania wyżej w hierarchii |
| **Dokumenty powiązane** | tak | odsyłacze rozdzielone `·`, powiązanie dwukierunkowe |
| **Prototypy odniesienia** | gdy dotyczy | ścieżki do plików w `design/05-okna/` |
| **Źródła normatywne** | tak | ścieżki źródeł, z których pochodzą wartości w dokumencie |
| **Zasada nadrzędna** | tak | jedno zdanie rozstrzygające, gdy treść dopuszcza dwie lektury |

Pole „Poza zakresem” nie jest ozdobnikiem: bez niego czytelnik nie wie, czy czegoś nie ma,
bo tego nie przewidziano, czy dlatego, że opisano to gdzie indziej.

### 2.2 Klasy dokumentów

| Klasa | Co opisuje | Przykłady |
|---|---|---|
| **Specyfikacja docelowa** | stan, który ma powstać — zakres budowy | opracowania w `architektura/`, `moduly/`, `interfejs-uzytkownika/`, `srodowiska/`, `specyfikacje/`, `funkcje-globalne/` |
| **Stan wdrożenia** | stan faktyczny kodu w bieżącym wydaniu | `README.md`, `INSTRUKCJA-UZYTKOWANIA.md`, `INSTALACJA-I-KONFIGURACJA.md` |
| **Akt prawny** | warunki korzystania z produktu | `LICENSE.md` |
| **Nawigacja** | porządek zbioru | `SPIS-OPRACOWAN.md`, niniejszy standard |

Klasa rozstrzyga o dwóch rzeczach: czy w dokumencie występują oznaczenia stanu wdrożenia
(rozdz. 2.6) oraz jaki repertuar form wizualnych obowiązuje (rozdz. 3.2).

### 2.3 Spis treści

```markdown
## Spis treści

1. [{Rozdział}](#1-rozdział)
   - [1.1 {Podrozdział}](#11-podrozdział)
   - [1.2 {Podrozdział}](#12-podrozdział)
2. [{Rozdział}](#2-rozdział)
```

| Zasada | Postać wiążąca |
|---|---|
| Nagłówek sekcji | zawsze `## Spis treści` |
| Poziom pierwszy | numeracja `1.`, `2.`, bez wcięcia |
| Poziom drugi | wcięcie trzema spacjami, myślnik `- `, numer `1.1` |
| Poziom trzeci | dopuszczalny w opracowaniach powyżej 1500 wierszy |
| Położenie numeru | wewnątrz odsyłacza: `[1.1 Tytuł](#…)`, nigdy `1.1 [Tytuł](#…)` |
| Kotwice | małe litery, polskie znaki zachowane, spacja → `-`, kropki i backticki usunięte, półpauza `—` daje podwójny `--` |
| Kompletność | spis obejmuje wszystkie nagłówki `##` i `###` korpusu |

Spisu treści nie umieszczamy w spisie treści. Rozdziałów „Stopka” i „Metryka” również.

### 2.4 Korpus dokumentu

Rozdziały numerujemy `## 1.`, podrozdziały `### 1.1`. Numeracja jest ciągła i nie ma luk.
Rozdziały oddzielamy poziomą linią `---`. Rozdział pierwszy zawsze odpowiada na pytanie
„czym jest opisywany byt”; rozdział ostatni zawsze zawiera kryteria odbioru (rozdz. 4.2).

Zakaz numeracji innej niż numeracja rozdziałów — patrz rozdz. 7.1.

### 2.5 Stopka

Trzy ostatnie bloki każdego pliku, treść dosłowna:

```markdown
---

*Koniec dokumentu. {Tytuł opracowania} — {klasa dokumentu}, wersja 2.0, {RRRR-MM-DD}.*

---
*Danaco Console — Platforma AI Workspace OS · v2.0 · status Deweloperski*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz.*
*Warunki korzystania: [Licencja produktu]({ścieżka}). Kontakt: support@danaco-group.pl*
```

Ścieżka do licencji jest względna wobec położenia pliku:

| Położenie pliku | Ścieżka w stopce |
|---|---|
| korzeń `docs/` | `LICENSE.md` |
| podkatalog `docs/{katalog}/` | `../LICENSE.md` |

Zapis `[Licencja produktu](../LICENSE.md)` w pliku leżącym w korzeniu `docs/` jest błędny —
wskazuje poza katalog zbioru, gdzie nie ma żadnego pliku licencji. Tak samo błędny jest
zapis `[Licencja produktu](LICENSE.md)` w pliku podkatalogu.

### 2.6 Oznaczenia stanu wdrożenia

Oznaczenia stanu **nie występują w opracowaniach klasy Specyfikacja docelowa**. Opracowanie
projektowe opisuje to, co ma powstać; znakowanie w nim, co już działa, myli zakres budowy
z postępem budowy i dezaktualizuje dokument przy każdym wydaniu. Postęp prac jest
przedmiotem raportu wykonawczego.

W opracowaniach klasy **Stan wdrożenia** obowiązuje jeden zestaw sześcioelementowy:

| Oznaczenie | Znaczenie |
|---|---|
| **[DZIAŁA]** | Funkcja domknięta od interfejsu albo od komendy kontraktu do skutku; potwierdzona audytem. |
| **[DZIAŁA CZĘŚCIOWO]** | Mechanizm działa w ograniczonym zakresie albo z udokumentowanym zastrzeżeniem. |
| **[NIEZINTEGROWANE]** | Kod istnieje i jest poprawny, lecz nie ma konsumenta — funkcja nie jest osiągalna dla Operatora. |
| **[ATRAPA]** | Element widoczny w interfejsie, za którym nie stoi realizacja. |
| **[BRAK]** | Przewidziane koncepcją, w bieżącej wersji niezaimplementowane. |
| **[DO DECYZJI OPERATORA]** | Kwestia nierozstrzygnięta; dokument jej nie przesądza. |

W opracowaniach klasy Specyfikacja docelowa z całego repertuaru zostaje **wyłącznie**
oznaczenie `[DO DECYZJI OPERATORA]` — jest niezbędne, bo rozdz. 4 zakazuje zgadywania.

---

## 3. Wymóg form wizualnych

### 3.1 Miara i próg

Każde opracowanie ma **co najmniej 30 % wierszy w formach innych niż ciągła proza**.

```
Objętość dokumentu
├── proza ciągła ............................ najwyżej 70 %
└── formy wizualne .......................... co najmniej 30 %
    ├── tabele
    ├── diagramy ASCII
    ├── bloki kodu i konfiguracji
    ├── listy definicyjne i wyliczeniowe
    ├── drzewa katalogów
    ├── maszyny stanów
    └── przebiegi decyzyjne i sekwencje
```

Miarą jest udział wierszy należących do form wizualnych w sumie wierszy pliku, liczony
z pominięciem metryki, spisu treści i stopki.

### 3.2 Repertuar form obowiązkowych

| Katalog | Formy wymagane w każdym opracowaniu |
|---|---|
| `architektura/` | diagram warstw · diagram sekwencji dla każdego przebiegu · schemat encji i relacji · drzewo katalogów · tabela kontraktów |
| `specyfikacje/` | tabela pełnego wykazu · maszyna stanów · macierz uprawnień |
| `interfejs-uzytkownika/` | szkic układu okna · mapa stref · tabela stanów kontrolki · przebieg nawigacji · tabela punktów łamania |
| `moduly/` | szkic okna głównego · mapa okien operacyjnych · tabela komend kontraktu · przepływ pracy jako diagram · tabela danych wykorzystywanych przez model |
| `srodowiska/` | mapa modułów środowiska · szkic przedsionka · przebieg wejścia do środowiska |
| `funkcje-globalne/` | szkic widoku · maszyna stanów · tabela wyzwalaczy |

### 3.3 Jakość diagramów

| Zasada | Wymóg |
|---|---|
| Znaki ramek | `─ │ ┌ ┐ └ ┘ ├ ┤ ┬ ┴ ┼ ► ▼ ◄ ▲` |
| Szerokość | do 100 znaków |
| Podpis | jedno zdanie nad diagramem, mówiące co diagram pokazuje |
| Legenda | pod diagramem, gdy diagram używa skrótów albo znaków o umownym znaczeniu |
| Nazewnictwo | każdy element nazwany dokładnie tak jak w prozie i w kodzie |
| Kompletność | diagram nie wprowadza nazwy, której nie ma w treści |
| Celowość | diagram objaśnia to, czego proza nie oddaje równie zwięźle; diagramy ozdobne są zakazane |

Przykład diagramu spełniającego wymogi — przebieg wejścia do środowiska:

```
Operator                Powłoka                           Rdzeń
   │                       │                                 │
   │  wybór środowiska     │                                 │
   ├──────────────────────►│                                 │
   │                       │  environment.enter              │
   │                       ├────────────────────────────────►│
   │                       │                                 │
   │                       │  environment · modules          │
   │                       │  sessions · focusedSessionId    │
   │                       │◄────────────────────────────────┤
   │  przedsionek          │                                 │
   │◄──────────────────────┤                                 │
   │                       │                                 │
```

Legenda: `environment.enter` — nazwa komendy kontraktu; `environment`, `modules`,
`sessions`, `focusedSessionId` — pola wyniku tej komendy.

---

## 4. Wymóg normatywnej szczegółowości

Opracowanie niesie **pełny ciąg budowy**. Deweloper nie ma prawa niczego dopowiadać.
Jeśli dokument nie podaje nazwy komendy, brzmienia etykiety, wartości żetonu albo kodu
błędu — jest to usterka opracowania, nie swoboda wykonawcy.

### 4.1 Źródła normatywne

Wartości pochodzą wyłącznie z poniższych źródeł. Wymyślanie wartości jest zakazane.

| Źródło | Zawartość | Co z niego bierzemy |
|---|---|---|
| `budowa/shared/contract.json` | 1119 komend · 68 obszarów · 552 struktury · 377 wyliczeń · 78 zdarzeń · 8 kodów błędów | nazwy komend i zdarzeń, pola żądania i wyniku, typy, wymagalność, wartości wyliczeń, kody błędów |
| `design/zasoby/zetony/zetony.css` | 184 żetony `--dn-*` | kolory, wymiary, odstępy, typografia, cienie, czasy, warstwy, punkty łamania |
| `design/zasoby/css/komponenty.css`, `rama.css`, `prototyp.css` | 199 klas `.dn-*` | nazwy komponentów, modyfikatorów i stanów |
| `design/zasoby/ikony/manifest.json` | wykaz ikon | nazwy ikon przypisane do kontrolek |
| `budowa/shared/kontrasty-progi.json` | progi kontrastu | wymagania dostępności |
| `design/05-okna/` | 38 prototypów | dosłowne brzmienia etykiet, przycisków, komunikatów |
| `budowa/client/src/` | kod powłoki | rzeczywiste identyfikatory i etykiety |

**Reguła weryfikowalności.** Każda wartość liczbowa, nazwa i etykieta w opracowaniu musi
dać się wskazać w jednym z powyższych źródeł. Wartość bez pokrycia oznaczamy jawnie
`[DO DECYZJI OPERATORA]`.

### 4.2 Wykazy obowiązkowe

Każde opracowanie klasy Specyfikacja docelowa zawiera komplet poniższych wykazów. Wykaz
niedotyczący opisywanego bytu zapisujemy z adnotacją „nie dotyczy” i uzasadnieniem — nie
pomijamy go milczeniem.

| Wykaz | Zawartość kolumn |
|---|---|
| **Komendy kontraktu** | nazwa komendy · obszar · pola żądania (nazwa, typ, wymagalność) · pola wyniku · zdarzenia zwrotne · kody błędów |
| **Etykiety interfejsu** | element · dosłowne brzmienie · miejsce wystąpienia |
| **Komunikaty** | sytuacja · dosłowna treść · rodzaj (potwierdzenie, ostrzeżenie, błąd, stan pusty) |
| **Żetony** | miejsce zastosowania · żeton `--dn-*` · czego dotyczy |
| **Komponenty** | klasa `.dn-*` · rola · modyfikatory · stany |
| **Skróty klawiszowe** | kombinacja · działanie · zasięg · kolizje |
| **Stany kontrolek** | kontrolka · spoczynek · wskazanie kursorem · wciśnięcie · ognisko · nieaktywny · ładowanie · pusty · błąd |
| **Punkty łamania** | próg szerokości · zachowanie układu |
| **Kryteria odbioru** | warunek sprawdzalny · sposób sprawdzenia |

**Fragmenty kodu jako wartość wymuszona.** Tam, gdzie konwencja projektowa nie dopuszcza
wariantu — struktura ładunku, sygnatura komendy, deklaracja żetonu, znacznik dostępności —
podajemy gotowy fragment, nie opis. Przykład zapisu komendy w postaci wymuszonej:

```json
{
  "typ": "environment.enter",
  "zadanie": {
    "environmentId": "string   — wymagane",
    "clientId":      "string   — wymagane",
    "sessionId":     "string   — opcjonalne; puste otwiera kartę pustą"
  }
}
```

### 4.3 Wykluczenie luk wykończeniowych

Sformułowania zakazane, ponieważ otwierają interpretację wykonawcy:

| Zakazane | Dlaczego | Czym zastąpić |
|---|---|---|
| „odpowiedni”, „stosowny”, „właściwy” | nie mówi, który | wskazać wartość |
| „w razie potrzeby”, „opcjonalnie” | nie mówi, kiedy | podać warunek |
| „na przykład”, „i tym podobne”, „między innymi” | zbiór otwarty | wyliczyć komplet |
| „można”, „powinno się rozważyć”, „warto” | nie zobowiązuje | rozstrzygnąć wprost |
| „typowo”, „zwykle”, „zazwyczaj” | opisuje zwyczaj, nie regułę | podać regułę |
| „i inne”, „itd.”, „itp.” | ukrywa braki | wyliczyć komplet |

Każde takie miejsce zastępujemy wartością rozstrzygniętą albo oznaczeniem
`[DO DECYZJI OPERATORA]` z pytaniem postawionym wprost.

---

## 5. Język

### 5.1 Słownik wiążący

Zamiana obejmuje **wyłącznie prozę**. Nazwy w blokach kodu, selektorach, ścieżkach, polach
kontraktu, identyfikatorach bazy danych oraz nazwy własne okien produktu pozostają
nietknięte.

| Termin porzucony | Postać wiążąca |
|---|---|
| hover | wskazanie kursorem |
| preview | podgląd |
| workflow | przepływ pracy |
| backend | warstwa serwerowa (rdzeń) |
| frontend | warstwa kliencka (powłoka) |
| toast | dymek powiadomienia |
| placeholder | tekst podpowiedzi |
| sandbox | piaskownica izolacyjna |
| checkbox | pole wyboru |
| tooltip | dymek podpowiedzi |
| layout | układ |
| toggle | przełącznik |
| mockup | makieta |
| modal | okno nakładkowe |
| dashboard | wskaźniki; typologia okna zbierającego wskaźniki brzmi „Okno monitorów i wskaźników” |
| klik | kliknięcie |
| feature | funkcja |
| update | aktualizacja |
| setup | konfiguracja początkowa |
| deploy, deployment | wdrożenie |
| hardcode, zero-hardcode | wartość wpisana na sztywno; bez wartości wpisanych na sztywno |
| Accepted (status decyzji) | Przyjęta |

**Zakaz odmieniania anglicyzmów po polsku.** Formy typu „checkboxa”, „toastu”,
„placeholderów”, „mockupów” są niedopuszczalne bez wyjątku — także tam, gdzie termin
angielski zostaje zachowany zgodnie z rozdz. 5.2.

**Wyjątek terminologiczny: `snippet`.** Termin `snippet` nie podlega zamianie na postać
polską — jak nazwy własne okien produktu, nazywa byt produktu, nie jego opis. Wyjątek
zapisany jest w wykazie terminów zachowanych rozdz. 5.2 i obejmuje wyłącznie postać
nieodmienną; zakaz odmiany z akapitu powyżej pozostaje wiążący także dla niego.

**Zakaz glosowania wariantowego.** Termin ma jedną postać w całym zbiorze. Zapisy
„Najechanie (hover)”, „Hover (najechanie)” i „stan hover” obok siebie są usterką nawet
wtedy, gdy każdy z osobna jest zrozumiały.

### 5.2 Terminy zachowane

| Termin | Powód zachowania |
|---|---|
| `commit`, `merge`, `branch`, `rebase` | terminy narzędzia Git, przyjęte w polskiej literaturze technicznej |
| Nazwy własne okien produktu | „Build Output”, „Metrics & Performance”, „Dev Tools”, „QA & Review Center” — nazwy elementów interfejsu, nie opisy |
| `snippet` | Wycinek kodu zapisany do ponownego wstawienia; termin przyjęty w polskiej praktyce inżynierskiej i nazywający byt produktu, nie jego opis. Postać nieodmienna, jak nazwy własne okien — zapisy „snippetu”, „snippetów”, „snippetami” pozostają zakazane rozdz. 5.1 |
| Nazwy modułów i środowisk | Assistant, Browser, Studio, TalkIn, WorkSpace, CodeStudio, MultitaskingAI |
| Identyfikatory techniczne | nazwy komend, pól kontraktu, kolumn bazy danych, klas i żetonów |

### 5.3 Zasady polszczyzny

| Zasada | Wymóg |
|---|---|
| Osoba | wyłącznie forma bezosobowa albo trzecia osoba; zakaz „dodałem”, „zrobiliśmy”, „zobaczysz” |
| Tryb | oznajmujący; zakaz trybu potocznego i zwrotów do czytelnika po imieniu |
| Rejestr | zawodowy; zakaz „po prostu”, „chyba”, „raczej”, „w sumie”, „jakoś” |
| Pleonazmy | zakaz konstrukcji typu „przepływ pracy pracy”, „okno okna dialogowego” |
| Glosy | termin objaśniamy raz, w miejscu pierwszego użycia w zbiorze, nie w każdym pliku |
| Zdania | jedna myśl na zdanie; zdanie powyżej czterdziestu wyrazów dzielimy |

### 5.4 Typografia

| Element | Postać wiążąca | Postać błędna |
|---|---|---|
| Myślnik zdaniowy | półpauza `—` | dywiz `-`, pauza `−` |
| Cudzysłów | `„treść”` | `"treść"`, `„treść"`, `“treść”` |
| Separator pozycji | `·` | `,` w wykazach metryki |
| Wyróżnienie | `**pogrubienie**` na wprowadzenie akapitu i na oznaczenia | podkreślenie, wersaliki |
| Kursywa | wyłącznie w stopce | w korpusie |
| Separator rozdziałów | `---` w osobnym wierszu | linia z dywizów |
| Łamanie wierszy | twarde, 78–100 znaków, poza tabelami | brak łamania |
| Odstępy | pojedyncza spacja | podwójna spacja, spacja przed przecinkiem |

---

## 6. Odsyłacze i nawigacja

| Zasada | Wymóg |
|---|---|
| Postać | odsyłacz markdown ``[Tytuł](../katalog/plik.md)``; zapis w backtickach nie jest odsyłaczem |
| Ścieżka | zawsze z katalogiem, także przy plikach o nazwie powtórzonej w zbiorze |
| Kotwica rozdziału | `[…](model-danych.md#…)` |
| Dwukierunkowość | jeśli opracowanie A wskazuje B, opracowanie B wymienia A w polu „Dokumenty powiązane” |
| Cel istniejący | odsyłacz do pliku nieistniejącego jest usterką; zakaz przywoływania opracowań planowanych |
| Liczby | zamiast odsyłać do wykazu, podajemy liczbę i wykaz wprost ze źródła normatywnego |

Nazwy plików kolidujące wielkością liter są zakazane: `INSTRUKCJA-UZYTKOWANIA.md`
i `instrukcja-uzytkowania.md` to jeden plik na Windows i na macOS.

---

## 7. Konwencje budowy obowiązujące na przyszłość

Rozdział wiąże nie tylko dokumentację, lecz także dalszą budowę produktu. Jego celem jest
zamknięcie źródła problemu, a nie tylko jego skutków.

### 7.1 Zakaz kodów wymyślonych

**Żadnych oznaczeń typu `ADL`, `D-01`, `P-3`, `L-P-7`, rejestrów decyzji ani numeracji
ustaleń.** Odwołujemy się nazwą rzeczy, nie kodem.

| Zapis zakazany | Zapis wiążący |
|---|---|
| „zgodnie z ADL-14” | „zgodnie z zasadą jednego źródła prawdy” |
| „patrz P-3” | „patrz [Model konfiguracji](architektura/model-konfiguracji.md)” |
| „decyzja D-07 przesądza” | „warstwowość konfiguracji przesądza” |
| „rejestr decyzji, poz. 22” | nazwa rozstrzygnięcia i miejsce jego opisu |

Numeracja dopuszczalna jest **wyłącznie** jako numeracja rozdziałów wewnątrz jednego
dokumentu. Kod, który istnieje tylko w głowie autora i wymaga tabeli rozwinięć, nie jest
oznaczeniem — jest szyfrem.

### 7.2 Zakaz określeń abstrakcyjnych

Nazwa opisuje rzecz, a nie jej rolę w narracji autora.

| Określenie ogólnikowe | Czego brakuje | Zapis wiążący |
|---|---|---|
| „mechanizm” | co to jest | „silnik kolejek zadań” |
| „warstwa pomocnicza” | czemu służy | „warstwa adapterów dostawcy modelu” |
| „element” | który | „kafel modułu na stronie głównej” |
| „odpowiedni komponent” | który | „`.dn-btn--sygnal`” |
| „system” bez dopowiedzenia | który | „rdzeń platformy” albo „powłoka kliencka” |

### 7.3 Nazewnictwo bytów kodu

Reguły spisane z tego, co w repozytorium już obowiązuje. Nowy byt otrzymuje nazwę zgodną
z rodziną, do której należy.

> **Zapis kluczy nastaw.** Nazwy w postaci `obszar.grupa.nastawa` użyte w tym rozdziale są
> **kluczami konfiguracji**, nie komendami kontraktu. Klucz wskazuje miejsce wartości w modelu
> konfiguracji; komendy kontraktu, którymi się go odczytuje i zapisuje, to `config.get` i
> `config.set` (obszar `config` w `budowa/shared/contract.json`).


| Rodzina | Reguła zapisu | Przykład wzorcowy |
|---|---|---|
| Pliki i katalogi dokumentacji | małe litery, myślnik jako separator, bez polskich znaków, nazwa rzeczownikowa | `model-konfiguracji.md` |
| Pliki źródłowe powłoki | jak wyżej, nazwa po polsku | `budowa/client/src/powloka/pasek-gorny.ts` |
| Żetony projektowe | `--dn-{rodzina}-{wariant}`; rodziny istniejące: `wym`, `sygnal`, `cien`, `szary`, `tekst`, `ostrzezenie`, `informacja`, `z`, `sukces`, `od`, `blad`, `obrys`, `fs`, `lh`, `ff`, `fw`, `atrament`, `czas`, `rama`, `wstazka`, `r`, `powierzchnia`, `fokus`, `tlo`, `nakladka`, `kropka`, `bp` | `--dn-sygnal-obrys`, `--dn-fs-base`, `--dn-bp-w2` |
| Klasy komponentów | `.dn-{komponent}` · element `.dn-{komponent}-{element}` · modyfikator `.dn-{komponent}--{wariant}` | `.dn-btn`, `.dn-btn--sygnal`, `.dn-awatar-stan` |
| Komendy kontraktu | `{obszar}.{czynność}` albo `{obszar}.{podobszar}.{czynność}`, notacja wielbłądzia w członach złożonych; nazwa obszaru wyłącznie z wykazu 68 obszarów w `budowa/shared/contract.json` — dziesięć z nich: `connection`, `home`, `environment`, `module`, `workspace`, `session`, `window`, `message`, `config`, `settings` | `environment.enter`, `speech.availability.get` |
| Pola ładunku kontraktu | notacja wielbłądzia, angielski, nazwa rzeczownikowa | `focusedSessionId`, `protocolVersion` |
| Kolumny bazy danych | małe litery, podkreślenie jako separator, bez polskich znaków | `srodowisko_id`, `zrodlo_typ` |
| Zdarzenia | `{obszar}.{rzecz}.{co się stało}` w czasie przeszłym | `session.message.appended` |

Nazwa nowego bytu nie powstaje w oderwaniu od tych rodzin. Jeśli byt nie mieści się
w żadnej, jest to sygnał, że rodzina wymaga rozszerzenia — rozszerzenie odnotowujemy tutaj,
a nie zakładamy własnej konwencji obok.

### 7.4 Reguła jednego źródła prawdy

Wartość istnieje w jednym miejscu i jest przywoływana, nigdy powielana.

```
                    ┌──────────────────────────┐
                    │   ŹRÓDŁO NORMATYWNE      │
                    │  contract.json           │
                    │  zetony.css              │
                    │  komponenty.css          │
                    └────────────┬─────────────┘
                                 │  przywołanie
              ┌──────────────────┼──────────────────┐
              ▼                  ▼                  ▼
       ┌────────────┐     ┌────────────┐     ┌────────────┐
       │ opracowanie│     │  prototyp  │     │    kod     │
       └────────────┘     └────────────┘     └────────────┘

       Zakazane: powielenie wartości w którymkolwiek z trzech miejsc.
```

| Zakazane | Wiążące |
|---|---|
| wartość szesnastkowa koloru w opisie wyglądu | nazwa żetonu `--dn-*` |
| wartość pikselowa wymiaru | nazwa żetonu z rodziny `wym` albo `od` |
| nazwa komendy przepisana z pamięci | nazwa odczytana z `contract.json` |
| własna lista pól ładunku | lista z `contract.json` |

### 7.5 Reguła zamkniętego zbioru

Nowa funkcja wchodzi do produktu **razem z opracowaniem**, które niesie komplet:

```
┌─────────────────────────────────────────────────────────────┐
│  Funkcja gotowa do budowy — warunki łączne                  │
├─────────────────────────────────────────────────────────────┤
│  ✓ nazwa komendy kontraktu i pola ładunku                   │
│  ✓ dosłowne brzmienie wszystkich etykiet i komunikatów      │
│  ✓ żetony i klasy komponentów użyte w widoku                │
│  ✓ komplet stanów kontrolek                                 │
│  ✓ zachowanie na wszystkich punktach łamania                │
│  ✓ skróty klawiszowe i ich kolizje                          │
│  ✓ kryteria odbioru — warunki sprawdzalne                   │
└─────────────────────────────────────────────────────────────┘
```

Brak któregokolwiek z siedmiu warunków oznacza, że funkcja nie jest gotowa do budowy.
Rozpoczęcie budowy mimo braku jest odstępstwem od standardu, a nie decyzją wykonawcy.

### 7.6 Tryb postępowania przy braku ustalenia

```
        czy wartość jest w źródle normatywnym?
                        │
            ┌───────────┴───────────┐
           tak                     nie
            │                       │
            ▼                       ▼
     przywołaj ją        czy Właściciel ją rozstrzygnął?
                                    │
                        ┌───────────┴───────────┐
                       tak                     nie
                        │                       │
                        ▼                       ▼
              zapisz i wskaż        [DO DECYZJI OPERATORA]
              podstawę              + pytanie postawione wprost
```

Zakazane jest przemycenie własnej propozycji wykonawcy jako faktu. Propozycja jest
dopuszczalna wyłącznie jako propozycja: oznaczona, opatrzona uzasadnieniem i skierowana
do rozstrzygnięcia.

---

## 8. Kontrola zgodności

Opracowanie jest przyjmowane, gdy wszystkie poniższe kontrole zwracają zero uchybień.

| Kontrola | Sprawdzany warunek |
|---|---|
| **Metryka** | obie tabele obecne, komplet pól obowiązkowych wypełniony; metryka produktowa liczy osiem wierszy, metryka dokumentu jedenaście pól |
| **Spis treści** | obecny, pozycje zgodne z rzeczywistymi nagłówkami, kotwice rozwiązywalne |
| **Stopka** | trzy bloki, treść dosłowna, ścieżka do licencji rozwiązywalna |
| **Odsyłacze** | wszystkie rozwiązywalne, wszystkie w postaci markdown, powiązania dwukierunkowe |
| **Formy wizualne** | udział nie mniejszy niż 30 % |
| **Wykazy obowiązkowe** | komplet z rozdz. 4.2 obecny albo opatrzony uzasadnioną adnotacją „nie dotyczy” |
| **Normatywność** | każda nazwa komendy obecna w `contract.json`; każdy żeton w `zetony.css`; każda klasa w arkuszach `design/zasoby/` |
| **Luki wykończeniowe** | zero sformułowań z rozdz. 4.3 bez towarzyszącego oznaczenia braku |
| **Terminologia** | zero terminów porzuconych z rozdz. 5.1 w prozie |
| **Typografia** | zero cudzysłowów prostych poza kodem, zero podwójnych spacji, zero dywizów w funkcji myślnika |
| **Kody wymyślone** | zero oznaczeń zakazanych rozdz. 7.1 |
| **Objętość** | patrz progi znakowe rozdz. 9 |

---

## 9. Progi objętości

Objętość jest miarą pomocniczą — świadczy, że wykazy z rozdz. 4.2 i formy wizualne
z rozdz. 3 rzeczywiście wypełniono, a nie że dokument jest długi dla samej długości.
Miarą jest liczba znaków pliku źródłowego, licząc znaki znaczników Markdown.

| Klasa pliku | Próg | Uzasadnienie |
|---|---|---|
| Cztery opracowania Właściciela w korzeniu (`README.md`, `INSTRUKCJA-UZYTKOWANIA.md`, `INSTALACJA-I-KONFIGURACJA.md`, `LICENSE.md`) | **objętość nie może ulec zmniejszeniu** względem stanu zastanego | dokumenty niosą opis stanu faktycznego kodu i akt prawny — redakcja poprawia formę, nie skraca treści |
| Opracowania merytoryczne bardziej rozbudowane — katalogi `architektura/`, `moduly/`, `interfejs-uzytkownika/`, `srodowiska/`, `specyfikacje/`, `funkcje-globalne/` | **co najmniej 85 000 znaków** | katalogi niosące pełny ciąg budowy (rozdz. 4) — komplet wykazów przy tej objętości merytorycznej osiąga tę wielkość naturalnie, bez dopełniaczy |
| Pozostałe opracowania własne zbioru (nawigacja, standard) | **co najmniej 45 000 znaków** | dokumenty porządkujące zbiór, nie opisujące pojedynczego bytu platformy w pełnym wykazie |

**Zakaz osiągania progu przez dopełniacz.** Wiersz powtórzony, akapit rozwodniony,
zdanie bez nowej informacji nie podnoszą jakości dokumentu, nawet jeśli podnoszą
licznik znaków. Objętość jest skutkiem kompletności wykazów i diagramów (rozdz. 3–4),
nie odwrotnie. Kontrola przyjęcia dokumentu sprawdza oba warunki niezależnie: próg
znakowy **i** komplet wykazów obowiązkowych — spełnienie jednego bez drugiego nie
domyka pracy redakcyjnej.

---

## 10. Wzorzec w pełni wypełniony

Poniższy szkielet pokazuje minimalny komplet elementów opracowania klasy Specyfikacja
docelowa, z przykładowymi wartościami wziętymi z rzeczywistych źródeł normatywnych —
nie jest to dokument istniejący w zbiorze, wyłącznie ilustracja złożenia elementów
z rozdziałów 2–4 w jedną całość.

```markdown
# Danaco Console — Moduł Automations

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
| **Tytuł** | Moduł Automations |
| **Klasa dokumentu** | Specyfikacja docelowa |
| **Odbiorcy** | deweloper · projektant |
| **Przeznaczenie** | Ustala zakres funkcjonalny i zachowanie modułu Automations — automatyk, harmonogramów i przebiegów |
| **Zakres** | okna operacyjne modułu, komendy obszaru `automation`, żetony i komponenty widoku, stany kontrolek |
| **Poza zakresem** | silnik harmonogramu jako mechanizm rdzenia — [Architektura techniczna](../architektura/architektura.md) rozdz. 7 |
| **Dokument nadrzędny** | [Specyfikacja modułów](../specyfikacje/specyfikacja-modulow.md) |
| **Dokumenty powiązane** | [Model danych](../architektura/model-danych.md) · [Katalog komponentów](../interfejs-uzytkownika/katalog-komponentow.md) |
| **Prototypy odniesienia** | `design/05-okna/moduly/automations.html` |
| **Źródła normatywne** | `budowa/shared/contract.json` (obszar `automation`) · `design/zasoby/zetony/zetony.css` |
| **Zasada nadrzędna** | Automatyka bez właściciela nie uruchamia się — każdy przebieg ma przypisaną odpowiedzialność |

---

## Spis treści

1. [Czym jest moduł Automations](#1-czym-jest-moduł-automations)
2. [Okna operacyjne](#2-okna-operacyjne)
3. [Komendy kontraktu](#3-komendy-kontraktu)

---

## 1. Czym jest moduł Automations

{proza — jeden do trzech akapitów}

---

## 2. Okna operacyjne

┌───────────────────────────────────────────────────┐
│  Automations — okno główne                        │
├──────────────────┬────────────────────────────────┤
│  Harmonogramy    │  Karta harmonogramu            │
│  (lewa kolumna)  │  (prawa kolumna, szczegóły)    │
└──────────────────┴────────────────────────────────┘

## 3. Komendy kontraktu

| Komenda | Pola żądania | Pola wyniku | Zdarzenia |
|---|---|---|---|
| `automation.list` **[DO DECYZJI OPERATORA]** — nazwa nie występuje w `budowa/shared/contract.json`; do rozstrzygnięcia, czy komenda ma powstać | `environmentId?: string` | `automations: Automation[]` | — |
| `automation.run` | `automationId: string` | `runId: string` | `automation.run.started` |
```

Elementy, których w tym szkicu **nie wolno pominąć** w dokumencie docelowym: wykaz
etykiet interfejsu, wykaz żetonów, wykaz komponentów, wykaz skrótów klawiszowych,
wykaz stanów kontrolek, wykaz punktów łamania, kryteria odbioru oraz stopka pełna
(rozdz. 2.5). Szkic je pomija wyłącznie dla zwięzłości ilustracji.

---

## 11. Najczęstsze usterki redakcyjne

Wykaz zebrany z audytu zbioru poprzedzającego niniejszy standard — każda pozycja
wystąpiła co najmniej kilkanaście razy w zbiorze przed redakcją. Wykaz służy jako
lista kontrolna przy przeglądzie każdego opracowania. Trzy ostatnie pozycje pochodzą
z fali redakcyjnej katalogu `specyfikacje/` — metryka zredukowana do jednego pola,
nagłówek sekcji zmieniony w tabeli zbiorczej bez zmiany w każdym z jego wystąpień
korpusu oraz opracowanie bez odsyłacza przychodzącego — i dokumentują usterki
odnalezione dopiero przy przeglądzie obejmującym cały zbiór naraz, nie pojedynczy plik.

| Usterka | Zapis błędny | Zapis wiążący |
|---|---|---|
| Cudzysłów niesparowany | `„Deweloperski"` | `„Deweloperski”` |
| Odsyłacz w backtickach zamiast markdown | `` `architektura/model-danych.md` `` | ``[Model danych](../architektura/model-danych.md)`` |
| Termin glosowany wariantowo w jednym pliku | „Hover (najechanie)” w jednym akapicie, „najechanie (hover)” w drugim | jedna postać w całym pliku: „wskazanie kursorem” |
| Numer rozdziału poza odsyłaczem | `1.1 [Tytuł](#…)` | `[1.1 Tytuł](#…)` |
| Odsyłacz do rozdziału, którego nie ma | `rozdz. 4.5` w pliku mającym rozdziały 4.1–4.4 | wskazanie rzeczywistego rozdziału po weryfikacji w pliku docelowym |
| Kod wymyślony jako odwołanie | „zgodnie z ADL-14” | nazwa rozstrzygnięcia wprost, bez kodu |
| Odmiana anglicyzmu po polsku | „checkboxa”, „toastu”, „mockupów” | „pola wyboru”, „dymka powiadomienia”, „makiet” |
| Wartość wpisana na sztywno zamiast żetonu | „przycisk w kolorze `#2F6FED`” | „przycisk w kolorze `--dn-sygnal`” |
| Sformułowanie otwierające interpretację | „ustawia odpowiedni znacznik” | „ustawia znacznik `sesja_aktywna`” albo `[DO DECYZJI OPERATORA]` |
| Pierwsza osoba liczby pojedynczej | „dodałem konto, a dalej odmawia” | forma bezosobowa: „po dodaniu konta system nadal odmawia” |
| Stopka niepełna albo brak stopki | plik kończy się ostatnim zdaniem korpusu | trzy bloki stopki z rozdz. 2.5 |
| Odsyłacz do nieistniejącego pliku | `budowa/docs/kontrakt.md` (katalog nieistniejący) | odsyłacz do rzeczywistego źródła `budowa/shared/contract.json` albo liczba podana wprost |
| Powielenie wykazu zamiast odsyłacza do źródła | własna kopia listy 68 obszarów kontraktu w kilku plikach | odsyłacz do `[Architektura techniczna](../architektura/architektura.md)` rozdz. 18, jedynego miejsca wykazu pełnego |
| Nazwa pliku kolidująca wielkością liter | `Instrukcja.md` obok `instrukcja.md` | jedna z dwóch nazw zmieniona na rzeczowo odrębną |
| Metryka niepełna — jedno pole zamiast jedenastu | `**Informacje szczegółowe dokumentu:**` z samym polem „Źródło” | komplet jedenastu pól rozdz. 2.1: Tytuł, Klasa dokumentu, Odbiorcy, Przeznaczenie, Zakres, Poza zakresem, Dokument nadrzędny, Dokumenty powiązane, Prototypy odniesienia (gdy dotyczy), Źródła normatywne, Zasada nadrzędna |
| Nagłówek sekcji przemianowany niespójnie w obrębie pliku | „Przebieg pracy modułu” w spisie wykazów rozdz. 4.2, „Workflow modułu” w każdym z piętnastu wystąpień korpusu | zamiana wszystkich wystąpień jednocześnie — jedna nazwa sekcji w całym pliku, nie tylko w tabeli zbiorczej |
| Opracowanie osierocone — bez odsyłacza przychodzącego | plik istnieje w katalogu tematycznym, lecz żadne inne opracowanie ani Spis opracowań do niego nie odsyła | dodanie odsyłacza z opracowania nadrzędnego wskazanego w polu „Dokument nadrzędny” własnej metryki pliku |

**Przykład: naprawa metryki niepełnej.** Poniższy fragment pokazuje usterkę spotykaną
najczęściej w opracowaniach redagowanych przed ustanowieniem niniejszego standardu —
drugą tabelę metryki zredukowaną do jednego, niewiążącego pola — oraz jej naprawę do
kompletu jedenastu pól rozdziału 2.1.

Zapis błędny:

```markdown
**Informacje szczegółowe dokumentu:**

| | |
|---|---|
| **Źródło** | Koncepcja platformy i architektura |
```

Zapis naprawiony — komplet pól, każde wypełnione treścią właściwą opisywanemu opracowaniu,
nie przepisaną z sąsiedniego pliku:

```markdown
**Informacje szczegółowe dokumentu:**

| | |
|---|---|
| **Tytuł** | {pełny tytuł opracowania} |
| **Klasa dokumentu** | Specyfikacja docelowa |
| **Odbiorcy** | deweloper · projektant |
| **Przeznaczenie** | {jedno zdanie: do czego dokument służy} |
| **Zakres** | {co dokument obejmuje} |
| **Poza zakresem** | {co świadomie pominięto i gdzie tego szukać} |
| **Dokument nadrzędny** | [Koncepcja platformy](../architektura/koncepcja-platformy.md) |
| **Dokumenty powiązane** | {odsyłacze rozdzielone `·`} |
| **Źródła normatywne** | {ścieżki źródeł} |
| **Zasada nadrzędna** | {jedno zdanie rozstrzygające} |
```

Naprawa nie ogranicza się do dopisania brakujących wierszy pustą treścią — pole
wypełnione frazą „do uzupełnienia” albo pozostawione puste jest tą samą usterką w innej
postaci. Każde z jedenastu pól niesie treść rzeczywiście prawdziwą dla redagowanego
pliku, ustaloną lekturą jego korpusu, nie skopiowaną z wzorca ani z sąsiedniego
opracowania tego samego katalogu.

---

## 12. Polecenia kontrolne dla redaktora

Rozdział podaje gotowe polecenia powłoki sprawdzające zgodność jednego pliku albo
całego zbioru z rozdziałami 2–8, bez potrzeby ręcznego przeglądu. Polecenia uruchamia
się z katalogu `docs/`.

**Cudzysłów prosty poza kodem** — zero wyników oznacza zgodność:

```bash
awk '!/^```/{f=!f} f{next} /"/{print FILENAME":"FNR}' **/*.md
```

**Odsyłacz w backtickach zamiast markdown** — każdy wynik wymaga zamiany na `[Tytuł](ścieżka)`:

```bash
grep -rnoE '`[a-z0-9/_-]+\.md`' --include='*.md' .
```

**Termin porzucony w prozie** (rozdz. 5.1) — dopasowanie poza blokiem kodu:

```bash
grep -rniE '\b(hover|preview|workflow|toast|placeholder|checkbox|tooltip|toggle|mockup|dashboard)\b' \
  --include='*.md' . | grep -v '^\s*```'
```

Wynik dla `dashboard` wymaga odsiania trzech dopuszczalnych wystąpień: nazwy własnej okna
„Project Dashboard” (rozdz. 5.2), nazwy komendy `workspace.dashboard.get` oraz pola
`dashboard:WorkspaceDashboard` — identyfikatory techniczne zamianie nie podlegają.

```bash
grep -rniE '\bdashboard\b' --include='*.md' . | grep -vE 'Project Dashboard|workspace\.dashboard|WorkspaceDashboard'
```

**Kod wymyślony** (rozdz. 7.1) — musi zwrócić zero:

```bash
grep -rnE '\bADL[- ]?[0-9]|\b[DPL]-[0-9]+\b' --include='*.md' .
```

**Brak stopki** — pliki, których trzy ostatnie niepuste wiersze nie zaczynają się
od `*Danaco Console`:

```bash
for f in $(find . -name '*.md'); do
  tail -5 "$f" | grep -q '^\*Danaco Console' || echo "brak stopki: $f"
done
```

**Objętość poniżej progu** (rozdz. 9):

```bash
for f in $(find . -name '*.md'); do
  n=$(wc -c < "$f"); echo "$n $f"
done | sort -n | head -20
```

**Odsyłacz martwy** — ścieżka niewskazująca istniejącego pliku:

```python
import re, glob, os
links = re.compile(r'\[([^\]]*)\]\(([^)#\s]+)(#[^)]*)?\)')
for f in glob.glob('**/*.md', recursive=True):
    txt = re.sub(r'```.*?```', '', open(f, encoding='utf-8').read(), flags=re.S)
    for m in links.finditer(txt):
        u = m.group(2)
        if u.startswith(('http', 'mailto:')):
            continue
        p = os.path.normpath(os.path.join(os.path.dirname(f), u))
        if not os.path.exists(p):
            print(f, u)
```

**Opracowanie osierocone** — plik `.md` katalogu `docs/`, do którego nie odsyła żaden
odsyłacz markdown z żadnego innego pliku zbioru poza Spisem opracowań; wynik jest listą
kandydatów do sprawdzenia, nie automatycznym wyrokiem — plik wskazany jako `[Spis
opracowań](../SPIS-OPRACOWAN.md)` w polu „Dokument nadrzędny” własnej metryki, lecz bez
odsyłacza przychodzącego z żadnego opracowania tematycznego, wymaga dodania takiego
odsyłacza w opracowaniu nadrzędnym:

```python
import re, glob, os

files = glob.glob('**/*.md', recursive=True)
incoming = {f: 0 for f in files}
link = re.compile(r'\[[^\]]*\]\(([^)#\s]+)(?:#[^)]*)?\)')
for f in files:
    body = re.sub(r'```.*?```', '', open(f, encoding='utf-8').read(), flags=re.S)
    for m in link.finditer(body):
        u = m.group(1)
        if u.startswith(('http', 'mailto:')):
            continue
        target = os.path.normpath(os.path.join(os.path.dirname(f), u))
        if target in incoming and target != f:
            incoming[target] += 1
for f, n in sorted(incoming.items()):
    if n == 0 and f not in ('SPIS-OPRACOWAN.md',):
        print('osierocone:', f)
```

Trzy pliki są wyłączone spod tej kontroli z natury swojej roli, nie przez wyjątek
zapisany w skrypcie: `README.md` jest punktem wejścia czytelnika z zewnątrz zbioru i nie
wymaga odsyłacza przychodzącego z wnętrza `docs/`; `SPIS-OPRACOWAN.md` jest korzeniem
nawigacyjnym całego zbioru z definicji rozdziału 2.2; `LICENSE.md` jest aktem prawnym
przywoływanym ze stopki każdego pliku, co dla skryptu liczy się jako odsyłacz
przychodzący i w praktyce nigdy nie oznacza go jako osieroconego. Wynik dla pozostałych
czterdziestu ośmiu plików jest rozstrzygający — zero wystąpień oznacza zgodność, każde
inne oznacza plik do naprawy przez dodanie odsyłacza w opracowaniu nadrzędnym.

Żadne z powyższych poleceń nie zastępuje przeglądu merytorycznego — wykrywają wyłącznie
usterki mechaniczne. Kompletność wykazów normatywnych (rozdz. 4.2) i jakość diagramów
(rozdz. 3.3) sprawdza się wyłącznie lekturą.

**Kolejność stosowania.** Polecenia mechaniczne uruchamia się przed przeglądem
merytorycznym, nie po nim — usterka typograficzna albo odsyłacz martwy odwraca uwagę
recenzenta od treści i utrudnia ocenę kompletności wykazów. Redaktor przechodzący
opracowanie stosuje kolejność: odsyłacze i cudzysłowy (mechaniczne) → rama redakcyjna
(rozdz. 2) → terminologia (rozdz. 5) → normatywność (rozdz. 4) → formy wizualne
(rozdz. 3) → objętość (rozdz. 9). Odwrócenie kolejności — dociąganie objętości
przed uzupełnieniem wykazów — grozi dopełniaczem zakazanym w rozdz. 9.

Polecenia mechaniczne uruchamia się ponownie po każdej turze redakcji, aż do zerowego
wyniku wszystkich ośmiu kontroli — dopiero wtedy opracowanie jest gotowe do przeglądu
merytorycznego przez drugą osobę. Kontrola opracowania osieroconego różni się od
pozostałych siedmiu tym, że nie sprawdza pojedynczego pliku, lecz cały zbiór naraz —
uruchamiana jest po zakończeniu fali redakcyjnej obejmującej wiele opracowań, nie przy
każdym pojedynczym pliku z osobna.

**Moc wiążąca standardu.** Niniejszy standard jest wiążący dla wszystkich opracowań
katalogu `docs/` powstałych po dniu 2026-08-20 oraz dla redakcji opracowań powstałych
wcześniej, bez wyjątku co do klasy dokumentu.

---

*Koniec dokumentu. Standard redakcyjny i językowy — Nawigacja, wersja 2.0, 2026-08-20.*

---
*Danaco Console — Platforma AI Workspace OS · v2.0 · status Deweloperski*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz.*
*Warunki korzystania: [Licencja produktu](LICENSE.md). Kontakt: support@danaco-group.pl*
