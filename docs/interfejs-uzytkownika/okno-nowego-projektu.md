# Danaco Console — Okno nowego projektu

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
| **Tytuł** | Okno nowego projektu |
| **Klasa dokumentu** | Specyfikacja docelowa |
| **Odbiorcy** | projektant — buduje formularz i jego zachowanie · deweloper — wiąże okno z kontraktem i modelem projektu |
| **Przeznaczenie** | Ustala kompletną budowę i przebieg okna zakładania projektu — okna Workspace środowiska otwieranego w trybie zakładania projektu — tak aby deweloper zbudował je bez rozstrzygania czegokolwiek samodzielnie. |
| **Zakres** | wywołanie okna z szyny nawigacji, rama i płótno dwukolumnowe, pięć sekcji formularza, panel podsumowania, żetony i komponenty, walidacje, stany kontrolek, skróty, punkty łamania, dostępność, komendy kontraktu wraz z lukami wykazanymi wprost, scenariusze, kryteria odbioru |
| **Poza zakresem** | praca w gotowym projekcie i Project Dashboard — [Moduł Workspace](../moduly/workspace.md); rama okna aplikacji (belka, szyna, wstążka, pasek stanu) jako byt osobny — [Rama okna](rama-okna.md); definicja encji projektu i jej pola — [Model danych](../architektura/model-danych.md); definicja agenta jako bytu — [Specyfikacja agentów](../specyfikacje/specyfikacja-agentow.md); przepływ „Utwórz z szablonu projektu” — brak odrębnego prototypu, poza zakresem niniejszego dokumentu |
| **Dokument nadrzędny** | [Elementy okien](elementy-okien.md) |
| **Dokumenty powiązane** | [Moduł Workspace](../moduly/workspace.md) · [Rama okna](rama-okna.md) · [Model danych](../architektura/model-danych.md) · [Okno historii sesji](okno-historii-sesji.md) · [Specyfikacja agentów](../specyfikacje/specyfikacja-agentow.md) · [Strona główna i nawigacja](strona-glowna-i-nawigacja.md) |
| **Prototypy odniesienia** | `design/05-okna/platformowe/nowy-projekt.html` |
| **Źródła normatywne** | `design/05-okna/platformowe/nowy-projekt.html` (rama, menu szyny, pięć sekcji formularza, panel podsumowania) · `budowa/shared/contract.json` (obszary `workspace`, `session`, `agent`) · `design/zasoby/zetony/zetony.css` · `design/zasoby/css/komponenty.css` |
| **Zasada nadrzędna** | Okno nowego projektu nie jest nakładką — jest oknem Workspace środowiska otwartym w trybie zakładania projektu. Wybór środowiska poprzedza wszystko inne, bo dokonuje się poza tym oknem: rozstrzyga go pozycja menu szyny nawigacji, zanim okno w ogóle się otworzy. |

---

## Spis treści

1. [Czym jest okno nowego projektu](#1-czym-jest-okno-nowego-projektu)
   - [1.1 Workspace środowiska w trybie zakładania projektu](#11-workspace-środowiska-w-trybie-zakładania-projektu)
   - [1.2 Wybór środowiska poprzedza okno, nie stoi w nim](#12-wybór-środowiska-poprzedza-okno-nie-stoi-w-nim)
   - [1.3 Relacja do encji WorkspaceProject](#13-relacja-do-encji-workspaceproject)
   - [1.4 Terminologia lokalna Workspace](#14-terminologia-lokalna-workspace)
2. [Umiejscowienie i sposób wywołania](#2-umiejscowienie-i-sposób-wywołania)
   - [2.1 Menu szyny „Nowy projekt”](#21-menu-szyny-nowy-projekt)
   - [2.2 Kontrast z „Nową sesją”](#22-kontrast-z-nową-sesją)
3. [Rama i płótno](#3-rama-i-płótno)
   - [3.1 Szkic pełnego okna](#31-szkic-pełnego-okna)
   - [3.2 Mapa stref płótna](#32-mapa-stref-płótna)
   - [3.3 Nagłówek płótna](#33-nagłówek-płótna)
4. [Pięć sekcji formularza](#4-pięć-sekcji-formularza)
   - [4.1 Sekcja 1 — Tożsamość projektu](#41-sekcja-1--tożsamość-projektu)
   - [4.2 Sekcja 2 — Instrukcje i model](#42-sekcja-2--instrukcje-i-model)
   - [4.3 Sekcja 3 — Zespół agentów](#43-sekcja-3--zespół-agentów)
   - [4.4 Sekcja 4 — Automatyzacje](#44-sekcja-4--automatyzacje)
   - [4.5 Sekcja 5 — Materiały wejściowe](#45-sekcja-5--materiały-wejściowe)
5. [Panel podsumowania](#5-panel-podsumowania)
6. [Żetony i komponenty](#6-żetony-i-komponenty)
7. [Walidacje i komunikaty](#7-walidacje-i-komunikaty)
8. [Stany kontrolek](#8-stany-kontrolek)
9. [Skróty klawiszowe](#9-skróty-klawiszowe)
10. [Punkty łamania](#10-punkty-łamania)
11. [Dostępność](#11-dostępność)
12. [Komendy kontraktu](#12-komendy-kontraktu)
   - [12.1 Struktury danych w postaci wymuszonej](#121-struktury-danych-w-postaci-wymuszonej)
   - [12.2 Wykaz komend](#122-wykaz-komend)
   - [12.3 Pola formularza bez pokrycia w komendach](#123-pola-formularza-bez-pokrycia-w-komendach)
   - [12.4 Zasięgi konfiguracji (ConfigScope)](#124-zasięgi-konfiguracji-configscope)
   - [12.5 Uproszczenie encji Agent w sekcji 3](#125-uproszczenie-encji-agent-w-sekcji-3)
   - [12.6 Weryfikacja unikalności nazwy](#126-weryfikacja-unikalności-nazwy)
   - [12.7 Od formularza do Project Dashboard](#127-od-formularza-do-project-dashboard)
13. [Etykiety interfejsu](#13-etykiety-interfejsu)
14. [Struktura zakładanego projektu](#14-struktura-zakładanego-projektu)
15. [Scenariusze](#15-scenariusze)
16. [Kryteria odbioru](#16-kryteria-odbioru)

---

## 1. Czym jest okno nowego projektu

### 1.1 Workspace środowiska w trybie zakładania projektu

Okno nowego projektu **nie jest oknem nakładkowym**. Prototyp
`design/05-okna/platformowe/nowy-projekt.html` deklaruje siebie atrybutem `data-prototyp-okno` jako
„Zakładanie projektu w Workspace środowiska” i niesie pełną ramę aplikacji — belkę tytułową, szynę
nawigacji, pasek edycji (wstążkę) i pasek stanu — dokładnie tę samą ramę, jaką niesie każde inne
okno platformy, opisaną w [Ramie okna](rama-okna.md). Belka tytułowa tego konkretnego zrzutu niesie
tytuł „Nowy projekt — WorkSpace › Workspace”: pierwszy człon to nazwa środowiska, drugi — moduł, w
którym okno stoi.

Komentarz źródłowy arkusza stylów lokalnych (prefiks `.np-*`) w prototypie stwierdza to wprost:
„Workspace środowiska w trybie zakładania projektu. Okno otwiera pozycja »Nowy projekt« strefy
pracy szyny nawigacji, po wskazaniu środowiska.” Innymi słowy: nie istnieje osobne „okno nowego
projektu” jako odrębny typ okna w rejestrze okien operacyjnych — istnieje moduł Workspace danego
środowiska, który w jednym ze swoich trybów pokazuje formularz zakładania projektu zamiast pulpitu
projektu.

Nagłówek treści okna (`.pt-okno-tytul`) niesie „Nowy projekt — WorkSpace”, z podpisem krojem
maszynowym (`.pt-mono`): „okno Workspace środowiska · wywoływane ze strefy pracy szyny nawigacji” — ten sam
podpis potwierdza status okna jako trybu modułu, nie osobnej klasy okna.

**Rozróżnienie nazw bliźniaczych: środowisko a moduł.** Belka tytułowa niesie „Nowy projekt —
WorkSpace › Workspace” — dwa człony oddzielone znakiem „›”, oba pochodzące od tego samego rdzenia
nazwy, lecz oznaczające dwa różne byty platformy. Pierwszy człon, „WorkSpace” (wielka litera S
wewnątrz wyrazu), jest nazwą **środowiska** — jednego z czterech wskazanych w menu szyny (rozdz.
2.1: TalkIn, WorkSpace, CodeStudio, MultitaskingAI). Drugi człon, „Workspace” (bez wielkiej litery
wewnątrz wyrazu), jest nazwą **modułu** — tego samego modułu, który niesie to okno w trybie
zakładania projektu (rozdz. 1.1). Zbieżność nazw nie jest przypadkiem redakcyjnym: środowisko
WorkSpace jest środowiskiem, dla którego moduł Workspace jest naturalnym miejscem pracy głównej —
podobnie jak środowisko TalkIn ma swój moduł wiodący Studio (rozdz. 3.4, wariant TalkIn), a
CodeStudio swój moduł programistyczny. Rozróżnienie wielkości litery „S” jest jedynym sygnałem
tekstowym odróżniającym oba byty w samej belce — czytelnik dokumentu i przyszły deweloper powinni
tej różnicy pilnować, bo dwa identycznie brzmiące, różnie zapisane słowa obok siebie w jednym
nagłówku są miejscem podatnym na pomyłkę redakcyjną w przyszłych zmianach tego okna.

### 1.2 Wybór środowiska poprzedza okno, nie stoi w nim

Poprzednia redakcja niniejszego dokumentu opisywała wybór środowiska jako selektor czterech kart
**wewnątrz** okna. Zweryfikowany prototyp tego nie potwierdza: selektor środowiska nie występuje w
treści okna (`.np-plotno`) w żadnej postaci. Wybór środowiska dokonuje się **przed** otwarciem okna,
w menu rozwijanym pozycji „Nowy projekt” szyny nawigacji (rozdz. 2.1) — dosłowny tekst opisu tej
pozycji menu w prototypie brzmi:

> „Po wskazaniu środowiska otwiera się od razu okno Workspace tego środowiska w trybie zakładania
> projektu — nazwa, katalog, instrukcje, agenci, automatyzacje i materiały wejściowe. Przedsionek
> środowiska nie jest po drodze.”

To zdanie jest zweryfikowanym spisem treści okna: sześć rzeczowników („nazwa, katalog, instrukcje,
agenci, automatyzacje i materiały wejściowe”) odpowiada dokładnie polom i sekcjom rozdz. 4 poniżej.
Zdanie „Przedsionek środowiska nie jest po drodze” jest zarazem kluczowym rozstrzygnięciem
nawigacyjnym: w odróżnieniu od otwarcia zwykłej nowej sesji (rozdz. 2.2), zakładanie projektu
pomija przedsionek środowiska i otwiera Workspace bezpośrednio.

Konsekwencja dla budowy dokumentu: rozdziały opisujące „wybór środowiska jako pole formularza”,
obecne w poprzedniej redakcji, są usunięte — środowisko jest kontekstem, w którym okno istnieje, nie
polem, które Operator wypełnia wewnątrz niego.

### 1.3 Relacja do encji WorkspaceProject

Projekt jest w Danaco Console jednostką organizującą pracę. Lid płótna (`.np-lid`, cytat dosłowny)
opisuje to wprost:

> „Projekt jest jednostką pracy środowiska: skupia katalog na dysku, instrukcję stałą, zespół
> agentów, automatyzacje i materiały wejściowe. Sesje otwierane później dziedziczą wszystko, co tu
> ustalisz — i mogą to uzupełniać, nie znosząc.”

Struktura `WorkspaceProject` z `budowa/shared/contract.json` niesie pola: `id`, `name`, `description`,
`status`, `createdAt`, `updatedAt`, `ownerNote`, `owner`. **Żadne z tych pól nie niesie identyfikatora
środowiska.** Jest to zweryfikowana rozbieżność między tym, co lid i formularz okna sugerują
(projekt jako byt „środowiska”), a tym, co struktura danych faktycznie przechowuje: powiązanie
projektu ze środowiskiem, jeśli istnieje, nie jest polem encji `WorkspaceProject`, lecz wynika
pośrednio z kontekstu okna/sesji, w którym projekt powstał. **[DO DECYZJI OPERATORA]** — czy
struktura `WorkspaceProject` wymaga dodania pola przechowującego środowisko, czy przynależność
środowiskowa ma pozostać relacją pośrednią, nieskładowaną wprost w encji.

Podobnie, pola widoczne w formularzu okna — katalog projektu, termin zamknięcia, jednostka pracy,
model wiodący, profil izolacji kontekstu — **nie mają odpowiednika** w strukturze `WorkspaceProject`
ani w polach żądania komendy `workspace.project.create` (`name`, `description`, `ownerNote`, `owner`,
`projectId` — rozdz. 12). Rozdział 12 wykazuje to pole po polu, bez pomijania żadnej rozbieżności
milczeniem.

### 1.4 Terminologia lokalna Workspace

Formularz i panel podsumowania przywołują pięć nazw własnych właściwych modułowi Workspace, które
nie mają osobnego opracowania w zestawie dokumentów przydzielonym temu redaktorowi. Zestawione tu
raz, w miejscu pierwszego pełnego użycia w niniejszym dokumencie, zgodnie z rozdz. 5.3 standardu
redakcyjnego zbioru — pełny opis każdego bytu należy do [Modułu Workspace](../moduly/workspace.md),
poza zakresem tego dokumentu.

| Nazwa własna | Występowanie w tym oknie | Rola (o ile wynika z kontekstu przywołania) |
|---|---|---|
| Project Dashboard | wzmiankowany pośrednio przez komendę `workspace.dashboard.get` (rozdz. 12.7) | widok stanu projektu po założeniu; to okno go nie pokazuje, wyłącznie poprzedza go w czasie |
| Project Library | podpowiedź sekcji 2 („pozycja źródłowa z Project Library”) i sekcja 5 („Trafiają do Project Library”) | miejsce docelowe plików wniesionych w sekcji 5 (rozdz. 4.5), zasilane komendą `workspace.library.upload` |
| Session Repository | opis automatyzacji „Zapis wersji po zaakceptowanej zmianie” (rozdz. 4.4) | miejsce, do którego trafiają wersje zapisywane automatycznie po akceptowanej zmianie w sesji projektu |
| Context Memory | lid płótna nie przywołuje tej nazwy wprost, lecz [Model danych](../architektura/model-danych.md) i inne opracowania zbioru wiążą ją z pamięcią projektu; pole `memoryEntryCount` struktury `WorkspaceDashboard` (rozdz. 12.7) jest jej odpowiednikiem liczbowym | pamięć trwała projektu, poza formularzem zakładania |
| Instructions Panel | opis pola `instructionSetCount` struktury `WorkspaceDashboard` w `contract.json` („licznik kafla Instructions Panel”) | widok wykazu zestawów instrukcji projektu w Project Dashboard; instrukcja stała sekcji 2 (rozdz. 4.2) jest jednym z zestawów, które ten kafel liczy |

Żadna z pięciu nazw nie jest tłumaczona na polski — rozdz. 5.2 standardu redakcyjnego zbioru
zachowuje nazwy własne modułów i widoków bez zmiany, tak jak nazwy środowisk (TalkIn, WorkSpace,
CodeStudio, MultitaskingAI) i modułów (Studio, Library, Browser). Rozróżnienie to obowiązuje
konsekwentnie w całym niniejszym dokumencie: gdziekolwiek pojawia się nazwa własna angielska
(„Project Dashboard”, „Project Library”, „Session Repository”, „Context Memory”, „Instructions
Panel”), towarzyszy jej — przy pierwszym użyciu w każdym rozdziale, w którym się pojawia — dość
kontekstu polskiego, by czytelnik nieznający wcześniejszych rozdziałów zrozumiał rolę bytu bez
konieczności przeszukiwania całego dokumentu wstecz.

---

## 2. Umiejscowienie i sposób wywołania

### 2.1 Menu szyny „Nowy projekt”

Jedyna zweryfikowana droga otwarcia okna: pozycja „Nowy projekt” w sekcji „Praca” szyny nawigacji
(`design/zasoby/rama.css`, opisana w [Ramie okna](rama-okna.md)), z etykietą wskazania kursorem:
„Założenie projektu — wskazanie środowiska otwiera od razu jego Workspace.” Kliknięcie rozwija menu
z nagłówkiem „Nowy projekt — w którym środowisku”, niosące sześć pozycji:

| Pozycja menu | Atrybut danych | Podpis (skrót etykiety) |
|---|---|---|
| TalkIn | `data-nowy-projekt-srodowisko="talkin"` | „Wiedza, komunikacja i praca z treścią” |
| WorkSpace | `data-nowy-projekt-srodowisko="workspace"` | „Projekty, procesy i produkty” |
| CodeStudio | `data-nowy-projekt-srodowisko="codestudio"` | „Programowanie, terminal i kontrola wersji” |
| MultitaskingAI | `data-nowy-projekt-srodowisko="multitaskingai"` | „Orkiestracja pracy ciągłej” |
| *(separator)* | — | — |
| Otwórz projekt istniejący | — | skrót `Ctrl+O` |
| Utwórz z szablonu projektu | — | brak zrzutu ekranu; poza zakresem niniejszego dokumentu (metryka, pole „Poza zakresem”) |

Wybór jednej z czterech pierwszych pozycji jest jedyną drogą dotarcia do formularza opisanego w
rozdz. 3–5 niniejszego dokumentu. Pozycja „Otwórz projekt istniejący” prowadzi do wykazu projektów
istniejących (poza zakresem — [Moduł Workspace](../moduly/workspace.md)), nie do formularza
zakładania. Skrót `Ctrl+O` tej pozycji jest tożsamy ze skrótem pozycji „Otwórz projekt…” w menu
aplikacji Plik, co potwierdza, że obie pozycje prowadzą do tego samego miejsca dwiema drogami.

Pozycja „Utwórz z szablonu projektu” nie ma zrzutu ekranu w żadnym z prototypów przeanalizowanych
dla niniejszego dokumentu, ale jej sama obecność w menu jest zweryfikowana — potwierdza, że
mechanizm szablonu projektu (przywoływany też przez przycisk „Zapisz jako szablon projektu” panelu
podsumowania, rozdz. 5) jest zamierzoną częścią przepływu zakładania projektu, nie tylko
funkcją zapisu jednokierunkową. Dwa fakty zweryfikowane niezależnie — istnienie pozycji odczytu
szablonu w menu szyny i istnienie przycisku zapisu szablonu w panelu — sugerują pętlę pełną: zapisz
formularz jako szablon (rozdz. 5) → wybierz „Utwórz z szablonu projektu” przy następnym zakładaniu
→ formularz sekcji 1–5 wypełnia się wartościami szablonu zamiast pustymi polami. Żaden z dwóch
punktów tej pętli nie ma jednak zweryfikowanego kontraktu (rozdz. 12.1 nie wykazuje struktury
niosącej szablon formularza) — cała pętla pozostaje poza zakresem normatywnym tego dokumentu,
odnotowana tu wyłącznie jako obserwacja spójności między dwoma miejscami interfejsu.

### 2.2 Kontrast z „Nową sesją”

Szyna nawigacji niesie w tej samej sekcji „Praca” pozycję siostrzaną „Nowa sesja”, o innym
zachowaniu — kontrast wart odnotowania, bo obie pozycje dotyczą zakładania czegoś nowego w
środowisku, lecz prowadzą różnymi drogami:

| Cecha | „Nowa sesja” | „Nowy projekt” |
|---|---|---|
| Etykieta wskazania kursorem | „Otwarcie nowej sesji — najpierw wskazanie środowiska, w którym ma powstać.” | „Założenie projektu — wskazanie środowiska otwiera od razu jego Workspace.” |
| Nagłówek menu | „Nowa sesja — w którym środowisku” | „Nowy projekt — w którym środowisku” |
| Krok po wskazaniu środowiska | otwiera **przedsionek** środowiska; moduł wiodący wybiera się kaflem w przedsionku | otwiera **od razu** okno Workspace w trybie zakładania projektu; przedsionek pomijany |
| Opis dosłowny (dymek podpowiedzi menu) | „Sesja otwiera przedsionek wskazanego środowiska. Moduł wiodący wybierasz kaflem w przedsionku — wtedy przestrzeń robocza otwiera się już z nim.” | „Po wskazaniu środowiska otwiera się od razu okno Workspace tego środowiska w trybie zakładania projektu — nazwa, katalog, instrukcje, agenci, automatyzacje i materiały wejściowe. Przedsionek środowiska nie jest po drodze.” |

Legenda: różnica jest istotna funkcjonalnie, nie tylko redakcyjnie — „Nowa sesja” zostawia Operatorowi
wybór modułu wiodącego w przedsionku; „Nowy projekt” zakłada z góry moduł Workspace jako miejsce
zakładania, bo projekt jest bytem Workspace niezależnie od tego, który moduł stanie się wiodący dla
pierwszej sesji otwartej w tym projekcie (rozdz. 5, nota panelu podsumowania).

---

## 3. Rama i płótno

### 3.1 Szkic pełnego okna

```
┌───────────────────────────────────────────────────────────────────────────────┐
│ ⧫ Danaco Console      Nowy projekt — WorkSpace › Workspace       [_][□][×]    │  belka — rama-okna.md
├───────────────────────────────────────────────────────────────────────────────┤
│☰│▤│⧗│      [Wstecz][Naprzód] [Centrum] [Odśwież]   [Szukaj……… Ctrl+K] [⋯]    │  pasek edycji — rama-okna.md
├──┬────────────────────────────────────────────────────────────────────────────┤
│  │  Nowy projekt — WorkSpace                                                  │  .pt-okno-naglowek
│sz│  okno Workspace środowiska · wywoływane ze strefy pracy szyny nawigacji    │
│y │ ┌────────────────────────────────────────────────────┐ ┌──────────────────┐│
│n │ │ WorkSpace · Workspace środowiska                   │ │  PODSUMOWANIE    ││  .np-panel
│a │ │ Nowy projekt                                       │ │ Środowisko       ││  (sticky)
│  │ │ Projekt jest jednostką pracy środowiska: …         │ │ WorkSpace        ││
│  │ │                                                    │ │ Nazwa            ││
│  │ │ ① Tożsamość projektu                               │ │ …                ││
│  │ │   Nazwa projektu · Katalog projektu · …            │ │ Katalog          ││
│  │ │                                                    │ │ …                ││
│  │ │ ② Instrukcje i model                               │ │ Model wiodący    ││
│  │ │   Instrukcja stała · Model wiodący · Izolacja      │ │ …                ││
│  │ │                                                    │ │ Izolacja kont.   ││
│  │ │ ③ Zespół agentów                                   │ │ …                ││
│  │ │   ☑ Redaktor  ☑ Analityk  ☐ Weryfikator  ☐ Korektor│ │ Agenci           ││
│  │ │                                                    │ │ 2 z 4            ││
│  │ │ ④ Automatyzacje                                    │ │ Automatyzacje    ││
│  │ │   ☑ Zapis wersji  ☑ Kopia dobowa  ☐ …  ☐ …         │ │ 2 czynne         ││
│  │ │                                                    │ │ Materiały        ││
│  │ │ ⑤ Materiały wejściowe                              │ │ 3 pliki·2,6 MB   ││
│  │ │   [ przeciągnij pliki albo wskaż w katalogu ]      │ ├──────────────────┤│
│  │ │   • zestawienie-kwartalne-Q3.xlsx      1,8 MB      │ │ [Załóż projekt   ││
│  │ │   • rejestr-umow-serwisowych.csv        640 kB     │ │  i otwórz sesję] ││
│  │ │   • wytyczne-redakcyjne-2026.docx       210 kB     │ │ [Załóż bez sesji]││
│  │ │                                                    │ │ [Zapisz szablon] ││
│  │ └────────────────────────────────────────────────────┘ └──────────────────┘│
├──┴────────────────────────────────────────────────────────────────────────────┤
│ Operator · support@danaco-group.pl  WorkSpace·Workspace   Sesje: 0   ⋯        │  pasek stanu — rama-okna.md
└───────────────────────────────────────────────────────────────────────────────┘
```

Legenda: cztery pasy ramy (belka, szyna, pasek edycji, pasek stanu) są opisane w
[Ramie okna](rama-okna.md) i nie są powtórzone tu poza szkicem — treść własna okna zaczyna się od
`.pt-okno-naglowek`. Kolumna prawa (`.np-panel`, 340 px) jest przyklejona (`position: sticky`) i
pozostaje widoczna przy przewijaniu pięciu sekcji lewej kolumny.

### 3.2 Mapa stref płótna

| Strefa | Klasa | Zawartość | Zachowanie |
|---|---|---|---|
| Płótno | `.np-plotno` | siatka dwukolumnowa `minmax(0, 1fr) 340px` | dzieli okno na kolumnę formularza i kolumnę podsumowania |
| Kolumna formularza | `.np-kolumna` (pierwsza) | nagłówek (rozdz. 3.3) + pięć sekcji (rozdz. 4) | przewijana, bez przyklejenia |
| Nagłówek płótna | `.np-glowa` | nadtytuł, tytuł, wprowadzenie | stały u góry kolumny formularza |
| Sekcja formularza | `.np-sekcja` | numer, tytuł, opis, pola | pięć wystąpień, numeracja `1`–`5` |
| Kolumna podsumowania | `.np-kolumna` (druga, `<aside>`) | panel podsumowania (rozdz. 5) | przyklejona (`position: sticky`, `top: var(--dn-od-6)`) powyżej progu 1240 px |

### 3.3 Nagłówek płótna

| Element | Klasa | Treść dosłowna |
|---|---|---|
| Nadtytuł | `.np-nadtytul` | „WorkSpace · Workspace środowiska” |
| Tytuł | `.np-tytul` | „Nowy projekt” |
| Wprowadzenie | `.np-lid` | cytat pełny w rozdz. 1.3 |

Nadtytuł niesie dwuczłonową wartość „{Środowisko} · {nazwa modułu}” — w zrzucie prototypu środowisko
WorkSpace, moduł Workspace. Dla pozostałych trzech środowisk (TalkIn, CodeStudio, MultitaskingAI)
nadtytuł przyjmuje analogiczną postać z nazwą właściwego środowiska — wartość nie ma odrębnego
zrzutu dla tych trzech i jest ustalona przez analogię, nie przez osobno zweryfikowany zrzut.
**[DO DECYZJI OPERATORA]** — dokładne brzmienie nadtytułu i modułu docelowego dla TalkIn, CodeStudio
i MultitaskingAI (czy moduł docelowy nosi zawsze nazwę „Workspace”, czy nazwę właściwą środowisku,
w postaci „Workspace TalkIn”).

---

## 4. Pięć sekcji formularza

Każda z pięciu sekcji (`.np-sekcja`) niesie strukturę stałą: numer w okrągłym znaczniku
(`.np-sekcja-nr`), tytuł (`.np-sekcja-tytul`), jednozdaniowy albo dwuzdaniowy opis
(`.np-sekcja-opis`) i pola właściwe. Numeracja sekcji (`1`–`5`) jest numeracją wizualną okna, nie
numeracją rozdziałów dokumentu — nie należy jej mylić z numeracją `## 4.1` poniżej.

### 4.1 Sekcja 1 — Tożsamość projektu

| Cecha | Wartość |
|---|---|
| Opis sekcji | „Nazwa, miejsce na dysku i ramy czasowe. Nazwa pojawi się w szynie sesji, w historii sesji i w pasku stanu.” |

| Pole | Typ kontrolki | Przykładowa wartość w zrzucie | Szerokość |
|---|---|---|---|
| Nazwa projektu | `.dn-pole-kontrolka` (`text`) | „Sprawozdawczość roczna 2026” | pełna szerokość sekcji (`.np-pole-szeroko`) |
| Katalog projektu | `.dn-pole-kontrolka` (`text`) | „/home/operator/danaco/sprawozdawczosc-2026” | pełna szerokość |
| *(podpowiedź)* | `.np-podpowiedz` | „Katalog jest jedynym miejscem, w którym projekt zapisuje wytwory. Poza nim sesja nie zapisze niczego bez osobnej zgody Operatora.” | pełna szerokość |
| Jednostka pracy | `.dn-pole-kontrolka` (`select`) | „projekt sprawozdawczy” (wybrana) — opcje: „projekt sprawozdawczy”, „projekt treściowy”, „projekt badawczy” | pół szerokości (`.np-para`, dwie kolumny) |
| Termin zamknięcia | `.dn-pole-kontrolka` (`text`) | „31 stycznia 2027” | pół szerokości |
| Opis projektu | `.dn-pole-kontrolka` (`textarea`, 3 wiersze) | „Roczne zestawienie wyników operacyjnych spółki wraz z załącznikami do sprawozdania zarządu. Materiał wejściowy pochodzi z rejestru umów oraz z zestawień kwartalnych.” | pełna szerokość |

Pole „Termin zamknięcia” niesie wartość tekstową swobodną w zrzucie („31 stycznia 2027”), nie
wartość kontrolki wyboru daty — **[DO DECYZJI OPERATORA]** czy pole docelowo pozostaje polem
tekstowym, czy przyjmuje kontrolkę wyboru daty z katalogu komponentów.

**Warianty dla pozostałych trzech środowisk.** Zrzut zweryfikowany w prototypie pokazuje wyłącznie
instancję środowiska WorkSpace (belka: „Nowy projekt — WorkSpace › Workspace”; nadtytuł płótna:
„WorkSpace · Workspace środowiska”). Menu szyny (rozdz. 2.1) oferuje trzy pozostałe środowiska —
TalkIn, CodeStudio, MultitaskingAI — jako równorzędne cele tej samej pozycji „Nowy projekt”, lecz
żadne z nich nie ma odrębnego zrzutu w prototypie. Poniższa tabela zbiera to, co jest zweryfikowane,
i oznacza to, co jest przeniesione przez analogię strukturalną (te same pięć sekcji, inne wartości
przykładowe dziedziny) bez potwierdzenia zrzutem.

| Środowisko | Nadtytuł płótna (analogia) | Jednostka pracy — przykład dziedzinowy | Status weryfikacji |
|---|---|---|---|
| WorkSpace | „WorkSpace · Workspace środowiska” | „projekt sprawozdawczy” | zweryfikowane zrzutem |
| TalkIn | „TalkIn · Workspace środowiska” | „projekt treściowy” | **[DO DECYZJI OPERATORA]** — analogia, bez zrzutu |
| CodeStudio | „CodeStudio · Workspace środowiska” | „projekt programistyczny” | **[DO DECYZJI OPERATORA]** — analogia, bez zrzutu |
| MultitaskingAI | „MultitaskingAI · Workspace środowiska” | „projekt orkiestracji” | **[DO DECYZJI OPERATORA]** — analogia, bez zrzutu; MultitaskingAI niesie ponadto panel orkiestracji zamiast listy modułów (`NavigationKind.orchestration`, [Strona główna i nawigacja](strona-glowna-i-nawigacja.md)), co może zmieniać treść sekcji 2 (pole „Model wiodący” traci sens dla panelu wieloagentowego) |

Struktura pięciu sekcji (rozdz. 4) i panelu podsumowania (rozdz. 5) jest przyjęta jako wspólna
wszystkim czterem środowiskom — nic w prototypie ani w kontrakcie nie sugeruje istnienia czterech
odrębnych układów formularza. Różni się wyłącznie treść przykładowa i, potencjalnie, dziedzinowe
brzmienie pola „Jednostka pracy” — żadna z tych różnic nie ma jednak potwierdzenia poza analogią
nazwaną wprost w tabeli powyżej.

### 4.2 Sekcja 2 — Instrukcje i model

| Cecha | Wartość |
|---|---|
| Opis sekcji | „Zasady obowiązujące każdą sesję projektu oraz zasięg, w jakim sesje mogą sięgać po dane.” |

| Pole | Typ kontrolki | Przykładowa wartość / opcje |
|---|---|---|
| Instrukcja stała projektu | `.dn-pole-kontrolka` (`textarea`, 5 wierszy) | „Stosuj rejestr formalny i liczby zapisuj słownie do dziesięciu. Każde twierdzenie liczbowe wiąż z pozycją źródłową z Project Library. Wnioski zapisuj osobnym akapitem, nigdy wewnątrz opisu danych.” |
| *(podpowiedź)* | `.np-podpowiedz` | „Instrukcja obowiązuje każdą sesję projektu i każdego agenta w nim pracującego. Sesja może ją uzupełnić, nie może jej znieść.” |
| Model wiodący | `.dn-pole-kontrolka` (`select`) | opcje: „Fable 5 · Ultra” (wybrana), „Fable 5 · Standard”, „Model lokalny” |
| Profil izolacji kontekstu | `.dn-pole-kontrolka` (`select`) | opcje: „Zamknięty — tylko katalog projektu” (wybrana), „Rozszerzony — katalog i biblioteka środowiska”, „Otwarty — z dostępem do sieci” |

Pole „Instrukcja stała projektu” odpowiada literalnie komendzie `workspace.instructions.set`
(rozdz. 12) — pole `content` tej komendy niesie dokładnie tę treść. Pola „Model wiodący” i „Profil
izolacji kontekstu” **nie mają odpowiednika** w strukturze `WorkspaceProject` ani w żadnej komendzie
obszaru `workspace` zweryfikowanej dla tego okna — rozdz. 12 odnotowuje to wprost jako lukę.

### 4.3 Sekcja 3 — Zespół agentów

| Cecha | Wartość |
|---|---|
| Opis sekcji | „Agenci przypisani do projektu. Wybór nie jest ostateczny — zestaw zmienia się w trakcie pracy.” |

| Agent | Zaznaczenie domyślne | Opis dosłowny |
|---|---|---|
| Redaktor | zaznaczony | „Redakcja tekstu w rejestrze formalnym; pilnuje instrukcji projektu.” |
| Analityk | zaznaczony | „Zestawienia liczbowe, kontrola zgodności danych ze źródłem.” |
| Weryfikator źródeł | niezaznaczony | „Sprawdza, czy każde twierdzenie ma pozycję w Project Library.” |
| Korektor | niezaznaczony | „Kontrola językowa i spójność terminologii w całym projekcie.” |

*(podpowiedź)* „Agenci pracują na tej samej instrukcji stałej. Zestaw zmienisz później w module
Agents — także w trakcie trwania projektu.”

Każdy wiersz jest polem wyboru (`.dn-check`) wewnątrz kafla `.np-wybor` w siatce
`repeat(auto-fit, minmax(220px, 1fr))` — liczba kolumn zależy od szerokości okna, nie jest stała.
Zaznaczenie agenta odpowiada literalnie komendzie `workspace.agent.assign` (rozdz. 12): dwóch
domyślnie zaznaczonych agentów (Redaktor, Analityk) odpowiada panelowi podsumowania „Agenci: 2 z 4”
(rozdz. 5).

### 4.4 Sekcja 4 — Automatyzacje

| Cecha | Wartość |
|---|---|
| Opis sekcji | „Czynności wykonywane bez polecenia. Każdą można wstrzymać w module Automations.” |

| Automatyzacja | Zaznaczenie domyślne | Opis dosłowny |
|---|---|---|
| Zapis wersji po zaakceptowanej zmianie | zaznaczona | „Każda przyjęta korekta tworzy wersję w Session Repository.” |
| Kopia dobowa do Library | zaznaczona | „Codziennie o 23:00 projekt trafia do repozytorium środowiska.” |
| Powiadomienie o zakończeniu zadania | niezaznaczona | „Wpis w pasku stanu oraz powiadomienie systemowe.” |
| Zestawienie tygodniowe postępu | niezaznaczona | „Skrót prac projektu wysyłany na adres konta w poniedziałki.” |

Dwie automatyzacje domyślnie zaznaczone odpowiadają panelowi podsumowania „Automatyzacje: 2 czynne”
(rozdz. 5). Żadna z czterech pozycji nie ma zweryfikowanego odpowiednika 1:1 w komendach obszaru
`automation` z `contract.json` (rozdz. 12) — obszar `automation` niesie model pełnych definicji
przepływu (kroki, harmonogram, szablony), nie prosty przełącznik projektowy. **[DO DECYZJI
OPERATORA]** — czy cztery pozycje tej sekcji są uproszczonym interfejsem nad `automation.template.apply`
(zastosowanie gotowego szablonu przepływu w kontekście projektu), czy osobnym mechanizmem
nieopisanym jeszcze w kontrakcie.

### 4.5 Sekcja 5 — Materiały wejściowe

| Cecha | Wartość |
|---|---|
| Opis sekcji | „Pliki wniesione do projektu na starcie. Trafiają do Project Library i są widoczne dla wszystkich sesji projektu.” |

| Element | Klasa | Treść dosłowna |
|---|---|---|
| Strefa upuszczenia | `.np-zrzutnia` | ikona pliku + „Przeciągnij pliki albo” + przycisk `.dn-btn--zarys.dn-btn--sm` „wskaż w katalogu” |
| Podpowiedź strefy | `.np-podpowiedz` | „Przyjmowane są dokumenty, arkusze, pliki tekstowe i obrazy — do 200 MB na plik.” |

Wykaz plików w zrzucie (`.np-plik`, wewnątrz `.np-pliki`):

| Plik | Rozmiar |
|---|---|
| zestawienie-kwartalne-Q3.xlsx | 1,8 MB |
| rejestr-umow-serwisowych.csv | 640 kB |
| wytyczne-redakcyjne-2026.docx | 210 kB |

Trzy pliki odpowiadają panelowi podsumowania „Materiały: 3 pliki · 2,6 MB” (rozdz. 5) — suma
1,8 MB + 0,64 MB + 0,21 MB ≈ 2,65 MB, zaokrąglona w podsumowaniu do „2,6 MB”. Każdy plik z tej listy
odpowiada jednemu wywołaniu komendy `workspace.library.upload` (rozdz. 12).

**Macierz zależności między pięcioma sekcjami.** Sekcje formularza nie są całkowicie niezależne —
trzy relacje warte odnotowania, choć żadna z nich nie jest wymuszona technicznie (zasada zero
blokad obowiązuje tu identycznie jak w rozdz. 7):

| Zależność | Kierunek | Skutek braku powiązania |
|---|---|---|
| Sekcja 2 (instrukcja) ↔ sekcja 3 (agenci) | instrukcja obowiązuje każdego agenta przypisanego w sekcji 3 | pusta instrukcja przy niepustym zestawie agentów nie blokuje założenia — agenci pracują bez instrukcji stałej, dziedzicząc wyłącznie własną definicję ([Specyfikacja agentów](../specyfikacje/specyfikacja-agentow.md)) |
| Sekcja 3 (agenci) ↔ sekcja 4 (automatyzacje) | brak zależności wprost — automatyzacje działają niezależnie od tego, którzy agenci są przypisani | żadna z czterech automatyzacji sekcji 4 nie odwołuje się do konkretnego agenta w opisie dosłownym (rozdz. 4.4) |
| Sekcja 1 (katalog) ↔ sekcja 5 (materiały) | pliki wniesione w sekcji 5 trafiają koncepcyjnie „do” projektu, którego katalog ustala sekcja 1 | pole `path` struktury `LibraryFile` (rozdz. 12.1) jest opcjonalne i względne wobec katalogu projektu — brak katalogu (rozdz. 12.3, pole bez pokrycia) nie ma zweryfikowanego wpływu na zapis plików biblioteki, ponieważ `workspace.library.upload` przyjmuje `projectId`, nie ścieżkę katalogu |

```
   Strefa upuszczenia i lista plików — sekcja 5
   ┌──────────────────────────────────────────────────────────────────┐
   │                                ⬆                                 │
   │           Przeciągnij pliki albo ( wskaż w katalogu )            │
   │  Przyjmowane są dokumenty, arkusze, pliki tekstowe i obrazy —    │
   │  do 200 MB na plik.                                              │
   ├──────────────────────────────────────────────────────────────────┤
   │ 📄 zestawienie-kwartalne-Q3.xlsx                        1,8 MB   │
   │ 📄 rejestr-umow-serwisowych.csv                         640 kB   │
   │ 📄 wytyczne-redakcyjne-2026.docx                        210 kB   │
   └──────────────────────────────────────────────────────────────────┘
```

---

## 5. Panel podsumowania

Panel `.np-panel` w kolumnie prawej (`<aside class="np-kolumna">`) jest przyklejony i widoczny
niezależnie od przewijania pięciu sekcji. Niesie trzy części: tytuł i osiem pozycji podsumowania,
trzy przyciski akcji, notę zamykającą.

| Pozycja podsumowania | Wartość w zrzucie | Sekcja źródłowa |
|---|---|---|
| Środowisko | „WorkSpace” | ustalone przez wejście z rozdz. 2.1, nie przez pole formularza |
| Nazwa | „Sprawozdawczość roczna 2026” | sekcja 1 |
| Katalog | „…/sprawozdawczosc-2026” (skrócone) | sekcja 1 |
| Model wiodący | „Fable 5 · Ultra” | sekcja 2 |
| Izolacja kontekstu | „Zamknięty” | sekcja 2 (wartość skrócona z pełnej etykiety opcji) |
| Agenci | „2 z 4” | sekcja 3 |
| Automatyzacje | „2 czynne” | sekcja 4 |
| Materiały | „3 pliki · 2,6 MB” | sekcja 5 |

Panel podsumowania **nie** niesie pozycji dla „Termin zamknięcia”, „Jednostka pracy” ani „Opis
projektu” — trzy pola sekcji 1 obecne w formularzu, nieobecne w podsumowaniu. Nie jest to
przeoczenie tego dokumentu: zrzut prototypu potwierdza dokładnie osiem pozycji wymienionych powyżej,
bez trzech pozostałych.

| Przycisk | Klasa | Rola |
|---|---|---|
| „Załóż projekt i otwórz sesję” | `.dn-btn--sygnal` (pełna szerokość, wyśrodkowany tekst) | akcja wiodąca — zakłada projekt i od razu otwiera pierwszą sesję w nim |
| „Załóż bez otwierania sesji” | `.dn-btn--zarys.dn-btn--sm` | zakłada projekt, zostawia go czekającym w Workspace środowiska |
| „Zapisz jako szablon projektu” | `.dn-btn--duch.dn-btn--sm` | zapisuje bieżące wypełnienie formularza jako szablon do ponownego użycia — bez zakładania projektu |

Nota zamykająca panel (`.np-nota`), cytat dosłowny:

> „Moduł wiodący pierwszej sesji wybierzesz kaflem w przedsionku środowiska. Projekt założony bez
> sesji czeka w Workspace środowiska i widnieje w historii sesji dopiero po pierwszym otwarciu.”

Nota rozstrzyga dwie rzeczy naraz: po pierwsze, że wybór modułu wiodącego pierwszej sesji **nie**
jest częścią tego formularza — dokonuje się później, w przedsionku, tym samym mechanizmem co przy
zwykłej nowej sesji (rozdz. 2.2); po drugie, że projekt założony przyciskiem „Załóż bez otwierania
sesji” nie pojawia się w [oknie historii sesji](okno-historii-sesji.md), dopóki żadna sesja nie
zostanie w nim otwarta — historia sesji rejestruje sesje, nie projekty jako takie.

**Reguły przeliczania ośmiu pozycji.** Panel podsumowania nie jest zrzutem statycznym niezależnym od
formularza — każda z ośmiu pozycji ma regułę przeliczenia jawną, wynikającą z pól sekcji 1–5
(rozdz. 4). Tabela poniżej podaje regułę dla każdej pozycji, nie tylko wartość przykładową.

| Pozycja | Reguła przeliczenia | Stan przy formularzu pustym (rozdz. 8) |
|---|---|---|
| Środowisko | wartość stała, ustalona przez wejście z rozdz. 2.1 — nie zmienia się w trakcie wypełniania formularza | nazwa środowiska wybranego przy wejściu |
| Nazwa | kopiuje dosłownie wartość pola „Nazwa projektu” (sekcja 1) | puste, dopóki pole nazwy jest puste |
| Katalog | kopiuje wartość pola „Katalog projektu” (sekcja 1), skróconą do końcowego segmentu ścieżki (`…/{ostatni-katalog}`) | puste |
| Model wiodący | kopiuje dosłownie wybraną opcję pola „Model wiodący” (sekcja 2) | pierwsza opcja listy — `select` nie ma stanu bez wyboru (rozdz. 8) |
| Izolacja kontekstu | kopiuje pierwszy człon wybranej opcji „Profil izolacji kontekstu” (sekcja 2), do myślnika — „Zamknięty — tylko katalog projektu” daje „Zamknięty” | pierwsza opcja listy, jak wyżej |
| Agenci | licznik `{zaznaczeni} z {wszyscy}` — liczba kafli `.np-wybor` zaznaczonych w sekcji 3, przez liczbę kafli dostępnych (cztery, rozdz. 4.3) | „0 z 4”, gdy żaden agent nie jest zaznaczony |
| Automatyzacje | licznik `{zaznaczone} czynne` — liczba kafli zaznaczonych w sekcji 4 (rozdz. 4.4) | „0 czynne” |
| Materiały | licznik `{n} pliki · {suma} MB` — liczba wierszy `.np-plik` i suma pól `sizeBytes` struktury `LibraryFile` (rozdz. 12.1), zaokrąglona do jednego miejsca po przecinku | „0 pliki”, bez człona wielkości — **[DO DECYZJI OPERATORA]**, dokładna forma gramatyczna licznika zerowego (np. „0 plików” zamiast „0 pliki”) nie ma zrzutu potwierdzającego |

Reguła pozycji „Katalog” i „Izolacja kontekstu” — obie skracają wartość pełną pola formularza do
postaci krótszej w panelu — jest jedynym miejscem panelu, w którym wartość widoczna w podsumowaniu
różni się dosłownie od wartości wpisanej w polu źródłowym. Pozostałych sześć pozycji kopiuje albo
liczy wartość źródłową bez skracania. Reguły przeliczenia opisane tabelą powyżej są odtworzone z
zachowania widocznego w jednym zrzucie statycznym — żaden przebieg dynamiczny (na przykład
odznaczenie agenta po uprzednim zaznaczeniu, i obserwacja, czy licznik cofa się natychmiast) nie ma
osobnego zrzutu potwierdzającego; reguły przyjęto jako najprostszą interpretację zgodną z jednym
stanem zaobserwowanym, spójną z ogólną zasadą aktualizacji na żywo (rozdz. 11).

---

## 6. Żetony i komponenty

Wykaz obejmuje wyłącznie żetony i klasy właściwe temu oknu (prefiks `.np-*` oraz komponenty
`.dn-*` użyte wewnątrz niego); żetony i komponenty ramy aplikacji (belka, szyna, wstążka, pasek
stanu) są przedmiotem [Ramy okna](rama-okna.md) i nie są powtórzone tu zgodnie z regułą jednego
źródła prawdy (standard redakcyjny zbioru, rozdz. 7.4).

| Miejsce zastosowania | Żeton `--dn-*` | Czego dotyczy |
|---|---|---|
| Siatka płótna | `--dn-od-6` | odstęp między kolumną formularza a kolumną podsumowania oraz dopełnienie płótna |
| Odstęp między sekcjami | `--dn-od-5` | rytm pionowy kolumny formularza |
| Karta sekcji | `--dn-obrys`, `--dn-powierzchnia`, `--dn-r-lg` | obrys, tło i zaokrąglenie `.np-sekcja` |
| Znacznik numeru sekcji | `--dn-obrys-mocny`, `--dn-r-pill`, `--dn-ff-mono` | `.np-sekcja-nr` |
| Nadtytuł i etykiety krojem maszynowym | `--dn-ff-mono`, `--dn-ls-mono-wersaliki`, `--dn-fs-xs` | `.np-nadtytul`, `.np-panel-tytul` |
| Tytuł płótna | `--dn-ff-naglowek`, `--dn-fw-polgruba`, `--dn-fs-2xl` | `.np-tytul` |
| Kafel wyboru (agent/automatyzacja) | `--dn-obrys`, `--dn-powierzchnia-2`, `--dn-r-md`, `--dn-obrys-mocny` (wskazanie kursorem) | `.np-wybor` |
| Strefa upuszczenia plików | obrys kreskowany `--dn-obrys-mocny`, `--dn-powierzchnia-2`, `--dn-wym-ikona-xl` | `.np-zrzutnia` |
| Wiersz pliku | `--dn-powierzchnia-2`, `--dn-r-sm`, `--dn-wym-ikona-sm` | `.np-plik` |
| Panel podsumowania | `--dn-obrys`, `--dn-panel`, `--dn-r-lg` | `.np-panel` — tło `--dn-panel` odróżnia panel od kart sekcji (`--dn-powierzchnia`) |
| Kreska działowa panelu | `--dn-obrys-subtelny` | `.np-kreska` |
| Przycisk wiodący panelu | `--dn-sygnal-wypelnienie`, `--dn-szary-0` | „Załóż projekt i otwórz sesję” (`.dn-btn--sygnal`) |
| Przyciski drugorzędne panelu | brak wypełnienia / `--dn-tekst-2` | „Załóż bez otwierania sesji” (`.dn-btn--zarys`), „Zapisz jako szablon projektu” (`.dn-btn--duch`) |
| Pole wyboru | `--dn-wym-check`, `--dn-atrament` | `.dn-check` w sekcjach 3 i 4 |
| Pierścień ogniska | `--dn-fokus`, `--dn-wym-fokus` | wszystkie kontrolki interaktywne |

**Zachowanie żetonu `--dn-panel` w obu motywach.** `--dn-panel` różni się od `--dn-powierzchnia`
(tło kart sekcji, rozdz. 4) w obu motywach, nie tylko w jednym: w motywie jasnym
`--dn-panel` = `--dn-szary-25` (`#FAFAFA`), o odcień ciemniejszy od bieli `--dn-powierzchnia`
(`--dn-szary-0`, `#FFFFFF`); w motywie ciemnym `--dn-panel` = `--dn-szary-850` (`#212121`), jaśniejszy
od `--dn-powierzchnia` w tym motywie (`--dn-szary-900`, `#181818`). Skutek jest spójny w obu
motywach: panel podsumowania jest zawsze subtelnie odróżnialny od kart sekcji, nie tylko obrysem —
w odróżnieniu od kolumny tożsamości okien wejściowych ([Okno instalatora](okno-instalatora.md)
rozdz. 3.2), gdzie odróżnienie opiera się na żetonie stałym w obu motywach, tutaj oba żetony
(`--dn-panel` i `--dn-powierzchnia`) zmieniają się z motywem, zachowując różnicę względną między
sobą.

| Żeton | Rola | Motyw jasny | Motyw ciemny |
|---|---|---|---|
| `--dn-panel` | tło panelu podsumowania | `--dn-szary-25` (`#FAFAFA`) | `--dn-szary-850` (`#212121`) |
| `--dn-powierzchnia` | tło kart sekcji formularza | `--dn-szary-0` (`#FFFFFF`) | `--dn-szary-900` (`#181818`) |

| Klasa `.np-*` | Rola | Uwaga |
|---|---|---|
| `.np-plotno` | siatka dwukolumnowa okna | `minmax(0, 1fr) 340px` |
| `.np-kolumna` | kolumna formularza albo kolumna podsumowania | druga kolumna jest `<aside>` |
| `.np-glowa`, `.np-nadtytul`, `.np-tytul`, `.np-lid` | nagłówek płótna | rozdz. 3.3 |
| `.np-sekcja`, `.np-sekcja-glowa`, `.np-sekcja-nr`, `.np-sekcja-tytul`, `.np-sekcja-opis` | struktura wspólna pięciu sekcji | rozdz. 4 |
| `.np-para` | dwie kolumny pól wewnątrz sekcji | sekcja 1 (jednostka pracy / termin), sekcja 2 (model / izolacja) |
| `.np-pole-szeroko` | pole zajmujące pełną szerokość `.np-para` | nazwa, katalog, opis |
| `.np-podpowiedz` | tekst pomocniczy pod polem albo sekcją | rozdz. 4.1, 4.2, 4.3, 4.5 |
| `.np-wybory`, `.np-wybor` | siatka kafli wyboru wielokrotnego | sekcje 3 i 4 |
| `.np-zrzutnia`, `.np-zrzutnia-tekst` | strefa upuszczenia plików | sekcja 5 |
| `.np-pliki`, `.np-plik`, `.np-plik-nazwa`, `.np-plik-meta` | wykaz plików wniesionych | sekcja 5 |
| `.np-panel`, `.np-panel-tytul` | panel podsumowania | rozdz. 5 |
| `.np-podsumowanie`, `.np-poz` | osiem pozycji podsumowania | rozdz. 5 |
| `.np-kreska` | separator poziomy wewnątrz panelu | rozdz. 5 |
| `.np-akcje` | trzy przyciski panelu | rozdz. 5 |
| `.np-nota` | tekst zamykający panel | rozdz. 5 |

---

## 7. Walidacje i komunikaty

Prototyp jest zrzutem statycznym, wypełnionym danymi przykładowymi — nie niesie zrzutu stanu błędu
ani komunikatu walidacji. Poniższa tabela ustala zachowanie zgodnie z zasadą zero blokad
(`komponenty.css`, w. 7: „żaden wariant nie odbiera klikalności — niegotowość komunikuje się opisem
albo komunikatem”), stosowaną konsekwentnie w całym zbiorze; treść komunikatów, gdzie nie ma
zrzutu potwierdzającego dosłowne brzmienie, jest oznaczona wprost.

Konsekwencja praktyczna tej zasady dla formularza o pięciu sekcjach jest większa niż dla formularza
jednopolowego: gdyby przycisk „Załóż projekt i otwórz sesję” (rozdz. 5) był wyłączany technicznie
do chwili spełnienia wszystkich warunków, Operator wypełniający sekcje w kolejności inne niż 1→5 —
na przykład zaczynający od wniesienia materiałów w sekcji 5 przed nadaniem nazwy w sekcji 1 — nie
miałby żadnego sygnału, dlaczego przycisk nie reaguje, dopóki nie cofnąłby się do sekcji 1. Zamiast
tego przycisk pozostaje klikalny od chwili otwarcia okna, a każda próba zatwierdzenia z brakującym
warunkiem odpowiada dokładnie tym brakiem — bez względu na to, którą z pięciu sekcji Operator
wypełnił jako ostatnią.

| Warunek | Zachowanie zgodne z zasadą zero blokad | Komunikat | Źródło |
|---|---|---|---|
| Pusta nazwa projektu, próba założenia | przycisk „Załóż projekt i otwórz sesję” pozostaje klikalny; pole nazwy przechodzi w stan błędu | „Nadaj projektowi nazwę” | **[DO DECYZJI OPERATORA]** — treść przyjęta analogicznie do wzorca innych formularzy zbioru |
| Katalog projektu bez prawa zapisu | pole katalogu przechodzi w stan błędu | „Brak prawa zapisu do wskazanego katalogu” | **[DO DECYZJI OPERATORA]**, jak wyżej |
| Nazwa projektu zajęta | pole nazwy przechodzi w stan błędu | „Projekt o tej nazwie już istnieje” | **[DO DECYZJI OPERATORA]** — zasięg unikalności (w obrębie środowiska, konta czy globalnie) nie ma zrzutu potwierdzającego |
| Plik przekraczający 200 MB (sekcja 5) | plik nie trafia do wykazu; komunikat przy strefie upuszczenia | „Plik przekracza limit 200 MB na plik” | wynika wprost z podpowiedzi sekcji 5 („do 200 MB na plik”), rozdz. 4.5 |
| Błąd zakładania po stronie rdzenia | okno pozostaje otwarte, formularz zachowuje wszystkie wypełnione pola | „Nie udało się założyć projektu — spróbuj ponownie” | **[DO DECYZJI OPERATORA]**, analogicznie do innych okien zakładających byty w zbiorze |

Zasada wspólna całemu zbiorowi: **okno nigdy nie gubi wprowadzonych danych przy błędzie zakładania.**
Nieudane założenie zostawia formularz — wraz z zaznaczeniami sekcji 3 i 4 oraz plikami sekcji 5 —
wypełniony i gotowy do ponowienia, bez konieczności odtwarzania wyborów od nowa.

---

## 8. Stany kontrolek

Wykaz w postaci ośmiu stanów wymaganych przez rozdz. 4.2 standardu redakcyjnego zbioru: spoczynek,
wskazanie kursorem, wciśnięcie, ognisko, nieaktywny, ładowanie, pusty, błąd. Kolumna „Nieaktywny”
niesie wartość „brak (zasada zero blokad)” tam, gdzie `komponenty.css` (w. 7) wyklucza stan
odbierający klikalność — nie jest to pominięcie kolumny, lecz jej rozstrzygnięcie wprost.

| Kontrolka | Spoczynek | Wskazanie kursorem | Wciśnięcie | Ognisko | Nieaktywny | Ładowanie | Pusty | Błąd |
|---|---|---|---|---|---|---|---|---|
| Pole tekstowe (`.dn-pole-kontrolka`) | obrys `--dn-obrys` | obrys mocniejszy | — | pierścień `--dn-fokus` | brak (zasada zero blokad) | — | tekst podpowiedzi, gdy pole bez wartości domyślnej (rozdz. 4 poniżej) | obrys `--dn-blad-obrys` + komunikat (rozdz. 7) |
| Pole `textarea` (opis, instrukcja) | obrys `--dn-obrys`, wysokość stała (3 albo 5 wierszy) | obrys mocniejszy | zmiana rozmiaru niedostępna — wysokość ustalona atrybutem `rows` | pierścień `--dn-fokus` | brak | — | tekst podpowiedzi | obrys `--dn-blad-obrys` |
| Lista rozwijana (`select`) | obrys `--dn-obrys` | obrys mocniejszy | rozwinięcie listy opcji | pierścień `--dn-fokus` | brak | — | pierwsza opcja wybrana domyślnie — lista `select` nie ma stanu „bez wyboru” w HTML | — |
| Kafel wyboru agenta/automatyzacji (`.np-wybor`) | obrys `--dn-obrys`, tło `--dn-powierzchnia-2` | obrys `--dn-obrys-mocny` | przełączenie zaznaczenia, przeliczenie panelu (rozdz. 5) | pierścień `--dn-fokus` na `.dn-check` wewnętrznym | brak | — | żaden kafel nie ma stanu pustego — wykaz czterech agentów/automatyzacji jest stały (rozdz. 4.3, 4.4) | — |
| Strefa upuszczenia (`.np-zrzutnia`) | obrys kreskowany `--dn-obrys-mocny` | podświetlenie tła podczas przeciągania pliku nad strefą | upuszczenie dodaje plik do wykazu | pierścień `--dn-fokus` na przycisku „wskaż w katalogu” | brak | pasek postępu wgrywania — **[DO DECYZJI OPERATORA]**, brak zrzutu tego stanu | wykaz plików pusty przed pierwszym wniesieniem — stan wyjściowy formularza | plik odrzucony przy przekroczeniu 200 MB (rozdz. 7) |
| Wiersz pliku (`.np-plik`) | tło `--dn-powierzchnia-2` | — | — | — | — | — | nie dotyczy — wiersz istnieje tylko dla pliku obecnego | — |
| Przycisk „Załóż projekt i otwórz sesję” | `.dn-btn--sygnal` | rozjaśnienie tła | przesunięcie `translateY(1px)` | pierścień `--dn-fokus` | brak (zasada zero blokad); brak nazwy daje komunikat, nie blokadę | wskaźnik `aria-busy="true"` w trakcie sekwencji wywołań (rozdz. 15, diagram) — **[DO DECYZJI OPERATORA]**, brak zrzutu | nie dotyczy | dymek błędu zakładania (rozdz. 7) |
| Przycisk „Załóż bez otwierania sesji” | `.dn-btn--zarys.dn-btn--sm` | tło `--dn-hover` | przesunięcie | pierścień `--dn-fokus` | brak | jak wyżej | nie dotyczy | jak wyżej |
| Przycisk „Zapisz jako szablon projektu” | `.dn-btn--duch.dn-btn--sm` | tło `--dn-hover` | przesunięcie | pierścień `--dn-fokus` | brak | — | nie dotyczy | — |
| Panel podsumowania (`.np-panel`) | osiem pozycji wypełnionych wartościami z formularza | nie dotyczy — panel nie jest kontrolką interaktywną poza trzema przyciskami | — | — | brak | — | wartości domyślne przy formularzu pustym (rozdz. 4 poniżej) | — |

**Stan pusty formularza — uwaga o zrzucie ilustracyjnym.** Zrzut zweryfikowany w prototypie pokazuje
formularz z wartościami przykładowymi już wypełnionymi (rozdz. 4: „Sprawozdawczość roczna 2026”,
pliki wniesione, agenci i automatyzacje zaznaczone) — jest to konwencja ilustracyjna powszechna w
prototypach zbioru, pokazująca układ z realistyczną treścią zamiast pustego formularza. **Stan
faktycznie pusty** (pola tekstowe bez wartości, żaden plik wniesiony, panel podsumowania z pozycjami
pustymi albo zerowymi) nie ma odrębnego zrzutu — kolumna „Pusty” powyżej ustala go przez wnioskowanie
z budowy pól HTML (`value` puste zamiast wypełnione, brak elementów `.np-plik`), nie przez zrzut
osobny. **[DO DECYZJI OPERATORA]** — dosłowna treść tekstu podpowiedzi każdego pola tekstowego w
stanie pustym nie ma potwierdzenia w prototypie, który pokazuje wyłącznie stan wypełniony; ta sama
uwaga dotyczy pola `textarea` w wierszu powyżej.

---

## 9. Skróty klawiszowe

| Skrót | Działanie | Zasięg | Kolizje |
|---|---|---|---|
| `Ctrl+O` | „Otwórz projekt istniejący” z menu szyny albo z menu Plik (rozdz. 2.1) | globalny | brak — `Ctrl+O` niesie to samo znaczenie w całej aplikacji, zgodnie z [Ramą okna](rama-okna.md) |
| `Tab` / `Shift+Tab` | przejście między polami formularza i kartami wyboru w kolejności sekcji 1 → 5, następnie panel podsumowania | okno | brak — kolejność ogniska okna nie przechwytuje przejścia do ramy aplikacji |
| `Spacja` | przełącza kafel wyboru (agent, automatyzacja) z ogniskiem | kafel z ogniskiem | kafel z ogniskiem przechwytuje `Spacja` przed przewinięciem widoku; poza kaflem `Spacja` przewija formularz |
| `Enter` w polu tekstowym jednowierszowym | nie zatwierdza formularza — pole tekstowe jednowierszowe nie zwalnia ogniska; zatwierdzenie wymaga aktywacji przycisku panelu | pole z ogniskiem | brak — okno nie ma skrótu zatwierdzenia formularza (akapit pod tabelą) |

Skróty ramy aplikacji (nawigacja, wstążka, `Alt+D` wyszukiwanie, `Ctrl+K` paleta poleceń) obowiązują bez zmian, zgodnie z
[Ramą okna](rama-okna.md) — okno nowego projektu jest oknem pełnym ramy, nie oknem wejściowym ani
nakładkowym, więc nie ma odrębnej puli skrótów ponad te cztery pozycje właściwe formularzowi.

Zbiór czterech skrótów jest zauważalnie węższy niż w oknach o porównywalnej liczbie pól — brak
skrótu zatwierdzającego formularz wprost z klawiatury (`Ctrl+Enter`) jest świadomym wynikiem
budowy okna: formularz niesie trzy przyciski o różnym skutku (rozdz. 5), więc jeden uniwersalny
skrót zatwierdzenia musiałby rozstrzygać, który z trzech uruchamia — rozstrzygnięcie, którego
prototyp nie podejmuje. **[DO DECYZJI OPERATORA]** — czy okno docelowo zyskuje skrót zatwierdzenia
przypisany jednoznacznie do przycisku wiodącego „Załóż projekt i otwórz sesję” (rozdz. 5), czy
zatwierdzenie pozostaje wyłącznie gestem myszy albo klawiatury przez `Tab` i aktywację przycisku.

---

## 10. Punkty łamania

| Próg szerokości | Zachowanie |
|---|---|
| > 1240 px | płótno dwukolumnowe: kolumna formularza elastyczna + panel podsumowania 340 px, przyklejony (`position: sticky`) |
| ≤ 1240 px | `.np-plotno` przechodzi na jedną kolumnę (`grid-template-columns: minmax(0, 1fr)`); panel podsumowania traci przyklejenie (`position: static`) i spada pod pięć sekcji formularza |

Próg 1240 px nie odpowiada dokładnie żadnemu żetonowi `--dn-bp-*` — najbliższy jest `--dn-bp-w3`
(1280 px). Wartość jest zapisana w arkuszu lokalnym `nowy-projekt.html` wprost, jako liczba
pikseli. Siatka kafli wyboru (`.np-wybory`, sekcje 3 i 4) korzysta z `repeat(auto-fit, minmax(220px, 1fr))`
niezależnie od progu 1240 px — liczba kolumn kafli zmienia się płynnie z szerokością dostępną, nie
skokowo na jednym progu nazwanym.

```
   Skutek przejścia progu 1240 px
   ┌───────────────────────────┬──────────────┐         ┌───────────────────────────┐
   │  KOLUMNA FORMULARZA       │  PODSUMOWANIE│   ──►   │     KOLUMNA FORMULARZA    │
   │  elastyczna               │  340 px      │ ≤1240px │                           │
   │                           │  (sticky)    │         ├───────────────────────────┤
   │                           │              │         │     PODSUMOWANIE          │
   └───────────────────────────┴──────────────┘         │     (static, pod spodem)  │
                    > 1240 px                           └───────────────────────────┘
```

Zejście panelu podsumowania pod pięć sekcji formularza poniżej progu 1240 px ma konsekwencję
praktyczną wartą odnotowania: na szerokościach wąskich Operator przewija całą treść pięciu sekcji,
zanim dotrze do trzech przycisków akcji (rozdz. 5) — inaczej niż powyżej progu, gdzie przyklejony
panel niesie te same przyciski widoczne bez przewijania niezależnie od postępu wypełniania
formularza. Próg 1240 px jest więc nie tylko progiem układu wizualnego, ale progiem różnicującym
liczbę gestów przewijania potrzebnych do zatwierdzenia formularza — od zera (panel zawsze widoczny)
do przewinięcia pełnej długości pięciu sekcji (panel na końcu strony). Punkt ten uzupełnia rozdz. 3.4
(mapa stref) o wymiar behawioralny, nie tylko geometryczny.

---

## 11. Dostępność

| Wymóg | Realizacja |
|---|---|
| Nawigacja klawiaturą | pełna; kolejność `Tab` przez pięć sekcji, następnie panel podsumowania (rozdz. 9) |
| Ognisko widoczne | pierścień `--dn-fokus` na każdej kontrolce interaktywnej |
| Etykiety pól | każde pole `.dn-pole` niesie `.dn-pole-etykieta` powiązaną z kontrolką przez zagnieżdżenie `<label>` |
| Kafle wyboru | `.np-wybor` jest elementem `<label>` obejmującym `.dn-check` i opis — kliknięcie w dowolnym miejscu kafla przełącza zaznaczenie |
| Strefa upuszczenia | przycisk „wskaż w katalogu” daje drogę klawiaturową równoważną przeciągnięciu pliku — upuszczenie nie jest jedyną drogą dodania pliku |
| Komunikaty walidacji | czytane przez `role="alert"` przy błędzie pola (rozdz. 7), zgodnie z konwencją zbioru |
| Panel podsumowania | aktualizuje się na żywo wraz z wypełnianiem formularza; wartości liczbowe („2 z 4”, „2 czynne”, „3 pliki · 2,6 MB”) są tekstem czytelnym przez czytnik ekranu, nie wyłącznie wskaźnikiem graficznym |
| Region aktywny dla przeliczeń panelu | zmiana wartości w `.np-podsumowanie` powinna nosić `aria-live="polite"`, tak aby czytnik ekranu ogłaszał przeliczenie bez przenoszenia ogniska — **[DO DECYZJI OPERATORA]**, atrybut nieobecny w znaczniku statycznym prototypu |
| Numer sekcji jako kontekst, nie jako jedyny identyfikator | `.np-sekcja-nr` ([]„1”–„5”) jest `aria-hidden="true"` w prototypie — numer jest ozdobą wizualną; tytuł sekcji (`.np-sekcja-tytul`) niesie znaczenie dla technologii wspomagającej, zgodnie z regułą „stan nigdy samym kolorem ani samym symbolem” |
| Odrzucenie pliku | komunikat o przekroczeniu limitu (rozdz. 7) powinien przenosić ognisko na strefę upuszczenia, tak aby komunikat błędu był natychmiast czytelny — **[DO DECYZJI OPERATORA]**, zachowanie ogniska przy odrzuceniu nie ma zrzutu potwierdzającego |
| Kontrast | zgodny z progami `budowa/shared/kontrasty-progi.json` w obu motywach |

---

## 12. Komendy kontraktu

Okno operuje na komendach obszaru `workspace` (projekt, instrukcje, agenci, biblioteka) i `session`
(otwarcie pierwszej sesji projektu). Nazwy i pola wprost z `budowa/shared/contract.json` — tabela
odróżnia jawnie pola formularza **z pokryciem** w kontrakcie od pól **bez pokrycia**, zamiast
milczeć o rozbieżności.

### 12.1 Struktury danych w postaci wymuszonej

Cztery struktury z `budowa/shared/contract.json` niosą byty, które okno zakłada albo zapisuje.
Zapisane w postaci wymuszonej — pole po polu, z typem i wymagalnością — zgodnie z rozdz. 4.2 wymogu
normatywnej szczegółowości, nie jako opis swobodny.

```json
{
  "struktura": "WorkspaceProject",
  "opis": "Projekt przestrzeni roboczej",
  "pola": {
    "id":          "string — wymagane",
    "name":        "string — wymagane",
    "description": "string — opcjonalne",
    "status":      "WorkspaceProjectStatus — wymagane (active | archived | paused)",
    "createdAt":   "int64 — wymagane, milisekundy epoki",
    "updatedAt":   "int64 — wymagane, milisekundy epoki",
    "ownerNote":   "string — opcjonalne",
    "owner":       "string — opcjonalne"
  }
}
```

```json
{
  "struktura": "WorkspaceInstructions",
  "opis": "Instrukcje systemowe projektu — nosi wartość pola „Instrukcja stała projektu”, sekcja 2",
  "pola": {
    "projectId":    "string — wymagane",
    "content":      "string — wymagane",
    "contentHash":  "string — opcjonalne",
    "scope":        "ConfigScope — wymagane (rozdz. 12.4)",
    "scopeId":      "string — opcjonalne",
    "updatedAt":    "int64 — wymagane, milisekundy epoki"
  }
}
```

```json
{
  "struktura": "WorkspaceAgentAssignment",
  "opis": "Przypisanie eksperta jako wykonawcy w projekcie — jeden wpis na agenta zaznaczonego w sekcji 3",
  "pola": {
    "projectId":       "string — wymagane",
    "agentId":         "string — wymagane",
    "role":            "string — opcjonalne",
    "defaultExecutor": "bool — opcjonalne",
    "assignedAt":      "int64 — wymagane, milisekundy epoki"
  }
}
```

```json
{
  "struktura": "LibraryFile",
  "opis": "Plik repozytorium wiedzy — jeden wpis na plik wniesiony w sekcji 5",
  "pola": {
    "id":              "string — wymagane",
    "name":            "string — wymagane",
    "path":            "string — opcjonalne",
    "mimeType":        "string — opcjonalne",
    "sizeBytes":       "int64 — opcjonalne",
    "sourceModuleId":  "string — opcjonalne",
    "projectId":       "string — opcjonalne",
    "tags":            "string[] — opcjonalne",
    "collectionIds":   "string[] — opcjonalne",
    "versionId":       "string — opcjonalne",
    "checksum":        "string — opcjonalne",
    "createdAt":       "int64 — wymagane, milisekundy epoki",
    "updatedAt":       "int64 — wymagane, milisekundy epoki"
  }
}
```

Pole `role` struktury `WorkspaceAgentAssignment` niesie rolę eksperta w projekcie jako tekst
swobodny — formularz sekcji 3 nie wypełnia go wprost (rozdz. 4.3 nie pokazuje pola roli obok
zaznaczenia), więc wartość `role` przy wywołaniu z tego okna pozostaje pusta albo przyjmuje nazwę
agenta jako wartość domyślną. **[DO DECYZJI OPERATORA]** — czy okno ma udostępnić pole roli osobno,
czy `role` ma pozostać nieużywane z poziomu tego formularza.

Piąta struktura, wynikowa komendy `session.create` wywoływanej przy przycisku „Załóż projekt i
otwórz sesję” (rozdz. 5, 12.2), zamyka sekwencję zapisu:

```json
{
  "struktura": "Session",
  "opis": "Sesja — wspólna dla plików, pamięci, projektu i agentów; wynik session.create wywołanej z projectId nowo założonego projektu",
  "pola": {
    "id":         "string — wymagane",
    "title":      "string — opcjonalne",
    "projectId":  "string — opcjonalne, tu wypełnione identyfikatorem projektu z workspace.project.create",
    "status":     "SessionStatus — wymagane (active | paused | finished | archived)",
    "windowIds":  "string[] — opcjonalne, puste bezpośrednio po session.create — okna zakłada się osobno",
    "createdAt":  "int64 — wymagane, milisekundy epoki",
    "updatedAt":  "int64 — wymagane, milisekundy epoki",
    "metadata":   "json — opcjonalne"
  }
}
```

Pole `windowIds` sesji założonej tą drogą jest puste bezpośrednio po wywołaniu `session.create` —
sesja istnieje, lecz nie niesie jeszcze żadnego okna komunikacji. Zgodnie z notą panelu (rozdz. 5),
moduł wiodący i pierwsze okno komunikacji Operator zakłada dopiero w przedsionku środowiska, poza
zakresem tego dokumentu — `session.create` wywołane z tego okna otwiera pustą kartę sesji powiązaną
z projektem, nie gotowe okno robocze.

### 12.2 Wykaz komend

| Komenda | Przeznaczenie | Pola żądania (z kontraktu) | Pola wyniku |
|---|---|---|---|
| `workspace.project.create` | zakłada projekt | `name` (wym), `description` (opc), `ownerNote` (opc), `owner` (opc), `projectId` (opc) | `project: WorkspaceProject` |
| `workspace.instructions.set` | zapisuje instrukcję stałą projektu (sekcja 2) | `projectId` (wym), `content` (wym), `scope` (opc), `scopeId` (opc) | `instructions: WorkspaceInstructions` |
| `workspace.agent.assign` | przypisuje każdego zaznaczonego agenta (sekcja 3) | `projectId` (wym), `agentId` (wym), `role` (opc), `defaultExecutor` (opc), `scope` (opc), `scopeId` (opc) | `assignment: WorkspaceAgentAssignment` |
| `workspace.library.upload` | wgrywa każdy plik wykazu (sekcja 5) | `projectId` (wym), `name` (wym), `contentBase64` (wym), `path` (opc) | `file: LibraryFile` |
| `session.create` | otwiera pierwszą sesję, gdy Operator wybierze „Załóż projekt i otwórz sesję” | `title` (opc), `projectId` (opc), `metadata` (opc) | `session: Session` |

### 12.3 Pola formularza bez pokrycia w komendach

| Pole formularza | Sekcja | Status weryfikacji |
|---|---|---|
| Katalog projektu | 1 | brak pola odpowiadającego w `WorkspaceProject` ani w `workspace.project.create` — **[DO DECYZJI OPERATORA]** |
| Jednostka pracy | 1 | brak pola odpowiadającego — **[DO DECYZJI OPERATORA]** |
| Termin zamknięcia | 1 | brak pola odpowiadającego — **[DO DECYZJI OPERATORA]** |
| Model wiodący | 2 | brak pola projektowego odpowiadającego; `Window.modelChannelId` istnieje na poziomie okna komunikacji, nie projektu — **[DO DECYZJI OPERATORA]**, czy wartość sekcji 2 ustawia domyślny kanał dla okien zakładanych później w tym projekcie |
| Profil izolacji kontekstu | 2 | brak pola projektowego odpowiadającego; `Window.executionEnv` istnieje na poziomie okna — **[DO DECYZJI OPERATORA]**, analogicznie do wiersza powyżej |
| Automatyzacje (cztery pozycje) | 4 | brak komendy 1:1; obszar `automation` niesie model pełnych definicji przepływu — rozdz. 4.4 |

Struktura żądania `workspace.project.create` nie niesie pola `environmentId` (rozdz. 1.3) — powiązanie
ze środowiskiem wybranym w rozdz. 2.1 nie jest przekazywane jako parametr tej komendy w kontrakcie
zweryfikowanym dla niniejszego dokumentu. Kolejność wywołań przy zatwierdzeniu formularza (przycisk
„Załóż projekt i otwórz sesję”) jest przedmiotem rozdz. 15, scenariusz A.

### 12.4 Zasięgi konfiguracji (ConfigScope)

Pola `scope` i `scopeId` komend `workspace.instructions.set` i `workspace.agent.assign` (rozdz. 12.2)
odwołują się do wyliczenia `ConfigScope` — dziewięć poziomów zasięgu, od najszerszego do
najwęższego, zweryfikowanych w `contract.json` wraz z opisem: „window jest najwęższy i wygrywa,
application najszerszy i ustępuje każdemu innemu”.

| Poziom | Wartość | Rola wobec tego okna |
|---|---|---|
| Aplikacja | `application` | najszerszy — nastawy programu jako całości, nie treści projektu; poza zasięgiem formularza |
| Globalny | `global` | poza zasięgiem formularza |
| Środowisko | `environment` | poza zasięgiem formularza — środowisko jest kontekstem wejścia (rozdz. 1.2), nie warstwą zapisu instrukcji |
| Moduł | `module` | poza zasięgiem formularza |
| Para modułów | `modulePair` | poza zasięgiem formularza |
| **Projekt** | `project` | **domyślny zasięg tego okna** — `workspace.instructions.set` bez pola `scope` przyjmuje ten poziom (opis komendy: „brak znaczy zasięg projektu”) |
| Karta sesji | `session` | zasięg węższy niż projektu; dostępny dopiero po otwarciu sesji, poza formularzem zakładania |
| Rola | `role` | poza zasięgiem formularza |
| Okno komunikacji | `window` | najwęższy — poza zasięgiem formularza, dotyczy pojedynczego okna otwartego wewnątrz sesji |

Instrukcja stała projektu (sekcja 2, rozdz. 4.2) jest więc zapisem na poziomie `project` — obowiązuje
każdą kartę sesji i każde okno komunikacji otwarte w tym projekcie, chyba że warstwa węższa (sesja
albo okno) nadpisze ją zapisem własnym, poza zakresem tego formularza.

### 12.5 Uproszczenie encji Agent w sekcji 3

Struktura `Agent` z `contract.json` niesie osiemnaście pól — tożsamość, model bazowy, warstwy
promptu, poziomy pamięci, uprawnienia, limit podagentów i więcej, opisane w pełni w
[Specyfikacji agentów](../specyfikacje/specyfikacja-agentow.md). Kafel `.np-wybor` sekcji 3
(rozdz. 4.3) pokazuje z tej struktury wyłącznie dwa pola: `name` (jako `<b>`, w brzmieniu „Redaktor”) i
`description` (jako `<small>`, w brzmieniu „Redakcja tekstu w rejestrze formalnym; pilnuje instrukcji
projektu.”). Żadne inne pole struktury `Agent` — model bazowy, poziomy pamięci, uprawnienia — nie
jest widoczne w tym formularzu; Operator zakładający projekt wybiera agentów po nazwie i opisie
zwięzłym, nie po ich pełnej konfiguracji. Pełna konfiguracja agenta pozostaje poza zakresem tego
okna, dostępna z poziomu modułu Agents (rozdz. 4.3, podpowiedź: „Zestaw zmienisz później w module
Agents”).

### 12.6 Weryfikacja unikalności nazwy

Komenda `workspace.project.list` (`query` opcjonalne, `includeArchived` opcjonalne,
`includeDeleted` opcjonalne, `limit` opcjonalne → `projects: WorkspaceProject[]`, `total: int`) jest
kandydatem naturalnym do sprawdzenia unikalności nazwy przed założeniem (rozdz. 7, komunikat
„Projekt o tej nazwie już istnieje”): wywołanie z `query` równym wpisanej nazwie i sprawdzenie, czy
`total` przekracza zero. Komenda nie niesie jednak pola zawężającego wykaz do jednego środowiska —
spójnie z brakiem pola `environmentId` w całej rodzinie komend `workspace.project.*` (rozdz. 1.3,
12.3) — więc zasięg sprawdzenia unikalności (w obrębie jednego środowiska czy globalnie, w obrębie
całego konta) pozostaje **[DO DECYZJI OPERATORA]**, niezależnie od tego, czy technicznie komenda
`workspace.project.list` jest tą właściwą do wywołania w tym celu.

### 12.7 Od formularza do Project Dashboard

Po założeniu projekt jest odczytywalny komendą `workspace.dashboard.get`, zwracającą strukturę
`WorkspaceDashboard` — zestawienie stanu dla widoku Project Dashboard, poza zakresem niniejszego
dokumentu ([Moduł Workspace](../moduly/workspace.md)). Tabela poniżej zestawia osiem pozycji panelu
podsumowania tego okna (rozdz. 5, wypełniane w trakcie zakładania) z polami `WorkspaceDashboard`
(wypełnianymi po założeniu, na bieżąco w trakcie życia projektu) — pokazuje, które pozycje panelu
mają żywy odpowiednik w dashboardzie, a które są jednorazowym podglądem właściwym wyłącznie
formularzowi zakładania.

| Pozycja panelu podsumowania (rozdz. 5) | Pole `WorkspaceDashboard` | Uwaga |
|---|---|---|
| Środowisko | — | brak pola środowiska w `WorkspaceDashboard`, spójnie z luką rozdz. 1.3 |
| Nazwa | `project.name` | pośrednio, przez pole zagnieżdżone `project: WorkspaceProject` |
| Katalog | — | brak pola — spójnie z luką rozdz. 12.3 |
| Model wiodący | — | brak pola — spójnie z luką rozdz. 12.3 |
| Izolacja kontekstu | — | brak pola — spójnie z luką rozdz. 12.3 |
| Agenci: „2 z 4” | `assignedAgentIds` (liczność) | dashboard niesie wykaz identyfikatorów, nie ułamek „X z Y” — mianownik „4” jest wyłącznie liczbą kafli sekcji 3 formularza zakładania, nie polem trwałym |
| Automatyzacje: „2 czynne” | — | brak pola dedykowanego automatyzacjom w `WorkspaceDashboard`; najbliższe pole `taskCount`/`taskDoneCount` dotyczy zadań, nie automatyzacji — **[DO DECYZJI OPERATORA]** |
| Materiały: „3 pliki · 2,6 MB” | `libraryFileCount`, `libraryUsedBytes` | odpowiednik żywy istnieje — dashboard utrzymuje te same dwie wielkości w czasie, formularz pokazuje ich wartość początkową |

`WorkspaceDashboard` niesie ponadto pola bez odpowiednika w panelu podsumowania tego okna:
`openSessionIds`, `memoryEntryCount`, `lastActivityAt`, `instructionSetCount`, `taskCount`,
`taskDoneCount`, `noteCount`, `memoryUsedBytes`, `memoryCapacityBytes` — wszystkie właściwe stanowi
projektu **po** rozpoczęciu pracy, nieistotne w chwili zakładania, gdy sesje, zadania i notatki
jeszcze nie istnieją.

---

## 13. Etykiety interfejsu

Pełny wykaz dosłownych brzmień tekstowych własnych okna (poza ramą aplikacji, dokumentowaną w
[Ramie okna](rama-okna.md)), zweryfikowanych w `design/05-okna/platformowe/nowy-projekt.html`. Wykaz
obejmuje wyłącznie treść stałą interfejsu — etykiety, tytuły, opisy i podpowiedzi — nie wartości
przykładowe wpisane w pola formularza (nazwa projektu, katalog, treść instrukcji), które są danymi
ilustracyjnymi, a nie tekstem interfejsu; te ostatnie są zestawione osobno przy każdym polu w
rozdz. 4.

| Element | Dosłowne brzmienie | Miejsce wystąpienia |
|---|---|---|
| Tytuł treści okna | „Nowy projekt — WorkSpace” | `.pt-okno-tytul` |
| Podpis pod tytułem | „okno Workspace środowiska · wywoływane ze strefy pracy szyny nawigacji” | `.pt-mono` |
| Nadtytuł płótna | „WorkSpace · Workspace środowiska” | `.np-nadtytul` |
| Tytuł płótna | „Nowy projekt” | `.np-tytul` |
| Wprowadzenie | cytat pełny w rozdz. 1.3 | `.np-lid` |
| Tytuł sekcji 1 | „Tożsamość projektu” | `.np-sekcja-tytul` |
| Opis sekcji 1 | „Nazwa, miejsce na dysku i ramy czasowe. Nazwa pojawi się w szynie sesji, w historii sesji i w pasku stanu.” | `.np-sekcja-opis` |
| Etykieta pola | „Nazwa projektu” | `.dn-pole-etykieta` |
| Etykieta pola | „Katalog projektu” | `.dn-pole-etykieta` |
| Podpowiedź katalogu | „Katalog jest jedynym miejscem, w którym projekt zapisuje wytwory. Poza nim sesja nie zapisze niczego bez osobnej zgody Operatora.” | `.np-podpowiedz` |
| Etykieta pola | „Jednostka pracy” | `.dn-pole-etykieta` |
| Etykieta pola | „Termin zamknięcia” | `.dn-pole-etykieta` |
| Etykieta pola | „Opis projektu” | `.dn-pole-etykieta` |
| Tytuł sekcji 2 | „Instrukcje i model” | `.np-sekcja-tytul` |
| Opis sekcji 2 | „Zasady obowiązujące każdą sesję projektu oraz zasięg, w jakim sesje mogą sięgać po dane.” | `.np-sekcja-opis` |
| Etykieta pola | „Instrukcja stała projektu” | `.dn-pole-etykieta` |
| Podpowiedź instrukcji | „Instrukcja obowiązuje każdą sesję projektu i każdego agenta w nim pracującego. Sesja może ją uzupełnić, nie może jej znieść.” | `.np-podpowiedz` |
| Etykieta pola | „Model wiodący” | `.dn-pole-etykieta` |
| Etykieta pola | „Profil izolacji kontekstu” | `.dn-pole-etykieta` |
| Tytuł sekcji 3 | „Zespół agentów” | `.np-sekcja-tytul` |
| Opis sekcji 3 | „Agenci przypisani do projektu. Wybór nie jest ostateczny — zestaw zmienia się w trakcie pracy.” | `.np-sekcja-opis` |
| Podpowiedź agentów | „Agenci pracują na tej samej instrukcji stałej. Zestaw zmienisz później w module Agents — także w trakcie trwania projektu.” | `.np-podpowiedz` |
| Tytuł sekcji 4 | „Automatyzacje” | `.np-sekcja-tytul` |
| Opis sekcji 4 | „Czynności wykonywane bez polecenia. Każdą można wstrzymać w module Automations.” | `.np-sekcja-opis` |
| Tytuł sekcji 5 | „Materiały wejściowe” | `.np-sekcja-tytul` |
| Opis sekcji 5 | „Pliki wniesione do projektu na starcie. Trafiają do Project Library i są widoczne dla wszystkich sesji projektu.” | `.np-sekcja-opis` |
| Tekst strefy upuszczenia | „Przeciągnij pliki albo” + przycisk „wskaż w katalogu” | `.np-zrzutnia-tekst` |
| Podpowiedź strefy | „Przyjmowane są dokumenty, arkusze, pliki tekstowe i obrazy — do 200 MB na plik.” | `.np-podpowiedz` |
| Tytuł panelu | „Podsumowanie” | `.np-panel-tytul` |
| Przycisk wiodący | „Załóż projekt i otwórz sesję” | `.dn-btn--sygnal` |
| Przycisk drugorzędny | „Załóż bez otwierania sesji” | `.dn-btn--zarys.dn-btn--sm` |
| Przycisk trzeciorzędny | „Zapisz jako szablon projektu” | `.dn-btn--duch.dn-btn--sm` |
| Nota panelu | cytat pełny w rozdz. 5 | `.np-nota` |
| Nagłówek menu szyny | „Nowy projekt — w którym środowisku” | `.dn-menu-naglowek` |
| Opis menu szyny | cytat pełny w rozdz. 1.2 | `.dn-menu-opis` |
| Etykieta wskazania kursorem pozycji szyny | „Założenie projektu — wskazanie środowiska otwiera od razu jego Workspace.” | `.dn-szyna-etyk` |

---

## 14. Struktura zakładanego projektu

Po zatwierdzeniu formularza projekt niesie zestaw powiązań wynikający z pięciu sekcji. Drzewo
poniżej jest strukturą logiczną — odzwierciedla powiązania encji w
[Modelu danych](../architektura/model-danych.md), nie układ katalogów na dysku, z wyjątkiem pola
„Katalog projektu”, które jest jedyną wartością wprost odnoszącą się do systemu plików.

Drzewo odpowiada dokładnie zestawowi przykładowemu z rozdz. 4–5 (projekt „Sprawozdawczość roczna
2026”, dwaj agenci zaznaczeni, trzy pliki wniesione) — nie jest strukturą uogólnioną z symbolami
zastępczymi, lecz odtworzeniem konkretnego przebiegu zweryfikowanego zrzutem, tak aby czytelnik mógł
zestawić każdą gałąź z konkretnym polem formularza opisanym wcześniej w dokumencie.

```
Projekt „Sprawozdawczość roczna 2026”
├── name:            „Sprawozdawczość roczna 2026”                (workspace.project.create)
├── description:     „Roczne zestawienie wyników operacyjnych…”   (workspace.project.create)
├── ownerNote:       {puste w zrzucie — nie wypełnione}
├── owner:           {puste w zrzucie — nie wypełnione}
├── instructions:    treść instrukcji stałej                      (workspace.instructions.set)
├── agentAssignments:
│   ├── Redaktor      (workspace.agent.assign)
│   └── Analityk      (workspace.agent.assign)
├── library:
│   ├── zestawienie-kwartalne-Q3.xlsx     1,8 MB   (workspace.library.upload)
│   ├── rejestr-umow-serwisowych.csv       640 kB   (workspace.library.upload)
│   └── wytyczne-redakcyjne-2026.docx      210 kB   (workspace.library.upload)
├── {katalog, jednostka pracy, termin, model wiodący, izolacja}  ← rozdz. 12, bez pokrycia w kontrakcie
└── pierwsza sesja:  otwierana warunkowo    (session.create, gdy „Załóż projekt i otwórz sesję”)
```

Legenda: gałęzie oznaczone nazwą komendy w nawiasie mają zweryfikowane pokrycie w
`budowa/shared/contract.json`; gałąź przedostatnia zbiera pola formularza bez takiego pokrycia,
wykazane osobno w rozdz. 12, nie pominięte milczeniem.

---

## 15. Scenariusze

Sekwencja poniżej rozwija krok 3 scenariusza A w postać diagramu — kolejność wywołań komend
kontraktu przy kliknięciu „Załóż projekt i otwórz sesję”, z zestawem przykładowym z rozdz. 4
(dwóch agentów, trzy pliki).

```
Operator         Okno Workspace                                Rdzeń
   │                    │                                        │
   │  „Załóż projekt    │                                        │
   │   i otwórz sesję”  │                                        │
   ├───────────────────►│                                        │
   │                    │  workspace.project.create              │
   │                    │  {name, description, ownerNote, owner} │
   │                    ├───────────────────────────────────────►│
   │                    │◄───────────────────────────────────────┤
   │                    │  project: WorkspaceProject             │
   │                    │                                        │
   │                    │  workspace.instructions.set            │
   │                    │  {projectId, content}                  │
   │                    ├───────────────────────────────────────►│
   │                    │◄───────────────────────────────────────┤
   │                    │                                        │
   │                    │  workspace.agent.assign  ×2            │
   │                    │  {projectId, agentId: Redaktor}        │
   │                    │  {projectId, agentId: Analityk}        │
   │                    ├───────────────────────────────────────►│
   │                    │◄───────────────────────────────────────┤
   │                    │                                        │
   │                    │  workspace.library.upload  ×3          │
   │                    │  {projectId, name, contentBase64}      │
   │                    ├───────────────────────────────────────►│
   │                    │◄───────────────────────────────────────┤
   │                    │                                        │
   │                    │  session.create                        │
   │                    │  {projectId}                           │
   │                    ├───────────────────────────────────────►│
   │                    │◄───────────────────────────────────────┤
   │                    │  session: Session                      │
   │  przedsionek       │                                        │
   │  środowiska        │                                        │
   │◄───────────────────┤                                        │
```

Legenda: siedem wywołań łącznie w zestawie przykładowym (jedno `workspace.project.create`, jedno
`workspace.instructions.set`, dwa `workspace.agent.assign`, trzy `workspace.library.upload`) plus
`session.create` warunkowe — ósme wywołanie sekwencji, gdy Operator wybrał „Załóż projekt i otwórz sesję”.
Automatyzacje sekcji 4 nie są narysowane — rozdz. 4.4 i 12.3 odnotowują brak zweryfikowanej komendy
dla tego pola.

**Scenariusz A — założenie projektu z otwarciem sesji (droga główna).**

1. Operator otwiera menu „Nowy projekt” szyny nawigacji i wskazuje środowisko WorkSpace
   (rozdz. 2.1). Otwiera się okno Workspace w trybie zakładania projektu, z pięcioma sekcjami
   pustymi albo z wartościami domyślnymi.
2. Operator wypełnia sekcję 1 (nazwa, katalog, jednostka pracy, termin, opis), sekcję 2 (instrukcja,
   model wiodący, izolacja), zaznacza dwóch agentów w sekcji 3, dwie automatyzacje w sekcji 4 i
   wnosi trzy pliki w sekcji 5 — panel podsumowania (rozdz. 5) przelicza się na bieżąco po każdej
   zmianie.
3. Operator klika „Załóż projekt i otwórz sesję”. Okno wywołuje kolejno: `workspace.project.create`
   (pola z sekcji 1 mające pokrycie), `workspace.instructions.set` (sekcja 2), `workspace.agent.assign`
   dwukrotnie (sekcja 3), `workspace.library.upload` trzykrotnie (sekcja 5), na końcu `session.create`
   z `projectId` nowo założonego projektu.
4. Po powodzeniu wszystkich siedmiu albo ośmiu wywołań otwiera się przedsionek środowiska
   WorkSpace — zgodnie z notą panelu (rozdz. 5), moduł wiodący pierwszej sesji Operator wybiera tam
   kaflem, nie w tym oknie — przedsionek środowiska jest opisany poza zakresem niniejszego
   dokumentu, w [Stronie głównej i nawigacji](strona-glowna-i-nawigacja.md) i w
   [Ramie okna](rama-okna.md).

**Scenariusz B — założenie bez otwierania sesji.**

1. Kroki 1–2 jak w scenariuszu A.
2. Operator klika „Załóż bez otwierania sesji”. Okno wywołuje te same komendy co w scenariuszu A, z
   wyjątkiem `session.create` — projekt powstaje, żadna sesja się nie otwiera.
3. Projekt czeka w Workspace środowiska WorkSpace. Zgodnie z notą panelu (rozdz. 5) nie pojawia się
   w [oknie historii sesji](okno-historii-sesji.md), dopóki Operator nie otworzy w nim pierwszej
   sesji — z poziomu Workspace środowiska, poza zakresem niniejszego dokumentu.

**Scenariusz C — zapisanie jako szablon, bez założenia projektu.**

1. Operator wypełnia formularz częściowo albo w całości, jak w krokach 1–2 scenariusza A.
2. Operator klika „Zapisz jako szablon projektu”. Żadna z pięciu komend rozdz. 12 nie jest
   wywoływana w tym scenariuszu w odniesieniu do bytu `WorkspaceProject` — zapisywany jest szablon,
   nie projekt. **[DO DECYZJI OPERATORA]** — komenda kontraktu obsługująca zapis szablonu formularza
   nie ma odpowiednika zweryfikowanego wśród komend obszaru `workspace` przeanalizowanych dla tego
   okna; pozycja menu szyny „Utwórz z szablonu projektu” (rozdz. 2.1) sugeruje istnienie takiego
   mechanizmu, lecz jego kontrakt pozostaje nieustalony w źródłach dostępnych temu dokumentowi.

**Scenariusz D — błąd zakładania po stronie rdzenia.**

1. Kroki 1–2 jak w scenariuszu A.
2. Operator klika „Załóż projekt i otwórz sesję”. Wywołanie `workspace.project.create` zwraca błąd
   (jeden z dziewięciu kodów globalnych kontraktu, w postaci `internal_error`; `budowa/shared/contract.json`
   zna dziś osiem, dziewiąty `command_not_understood` czeka na wdrożenie w kodzie).
3. Zgodnie z rozdz. 7 okno pozostaje otwarte, wszystkie pola pięciu sekcji zachowują wprowadzone
   wartości, komunikat „Nie udało się założyć projektu — spróbuj ponownie” pojawia się przy panelu
   podsumowania. Operator ponawia kliknięcie bez ponownego wypełniania formularza.

**Scenariusz E — projekt minimalny, tylko nazwa.**

1. Operator otwiera okno jak w kroku 1 scenariusza A, lecz wypełnia wyłącznie pole „Nazwa projektu”
   w sekcji 1 — katalog, jednostka pracy, termin i opis pozostają puste; sekcja 2 (instrukcja, model,
   izolacja) pozostaje pusta; żaden agent w sekcji 3 nie jest zaznaczony ponad domyślne dwa
   (Redaktor, Analityk — rozdz. 4.3, zaznaczenie domyślne niezależne od działania Operatora); żadna
   automatyzacja ponad domyślne dwie (rozdz. 4.4); sekcja 5 pozostaje bez plików.
2. Operator klika „Załóż bez otwierania sesji”. Zgodnie z rozdz. 12.2 wywoływane jest wyłącznie
   `workspace.project.create` z polem `name` wypełnionym i polami `description`, `ownerNote`,
   `owner` pustymi — trzy pozostałe komendy sekwencji (instructions.set, agent.assign,
   library.upload) uruchamiają się warunkowo, tylko gdy odpowiadająca im sekcja niesie treść albo
   zaznaczenie odbiegające od stanu wyjściowego; w tym scenariuszu `workspace.agent.assign` wywołuje
   się nadal dwukrotnie, bo dwaj agenci są zaznaczeni domyślnie, nie przez świadomy wybór Operatora
   ponad domyślność.
3. Projekt powstaje z minimalnym zestawem pól — potwierdza to, że żadne pole poza „Nazwa projektu”
   nie jest wymagane do założenia (rozdz. 7 nie wykazuje warunku wymagalności dla pozostałych pól
   sekcji 1–5).

**Scenariusz F — odrzucenie pliku przekraczającego limit w sekcji 5.**

1. Operator wykonuje kroki 1–2 scenariusza A, a w sekcji 5 usiłuje wnieść plik o rozmiarze
   250 MB — powyżej limitu „do 200 MB na plik” (rozdz. 4.5, podpowiedź strefy).
2. Zgodnie z rozdz. 7 strefa upuszczenia nie dodaje pliku do wykazu (`.np-pliki`); komunikat „Plik
   przekracza limit 200 MB na plik” pojawia się przy strefie `.np-zrzutnia`. Panel podsumowania nie
   przelicza pozycji „Materiały” — pozostaje przy stanie sprzed próby.
3. Operator wnosi plik mniejszy w miejsce odrzuconego; wykaz i podsumowanie przeliczają się
   normalnie, jak w kroku 2 scenariusza A. Odrzucenie pojedynczego pliku nie wpływa na pozostałe
   cztery sekcje formularza — stan sekcji 1–4 pozostaje nienaruszony przez błąd w sekcji 5.

---

## 16. Kryteria odbioru

| Warunek | Sposób sprawdzenia |
|---|---|
| Okno nie jest nakładką — niesie pełną ramę aplikacji | porównanie z [Ramą okna](rama-okna.md): obecność belki, szyny, wstążki i paska stanu w prototypie |
| Jedyną drogą wejścia jest menu szyny „Nowy projekt” z wyborem środowiska | otwarcie menu i sprawdzenie sześciu pozycji z rozdz. 2.1 |
| Wybór środowiska otwiera od razu Workspace, z pominięciem przedsionka | wybranie środowiska i obserwacja przejścia — kontrast z „Nową sesją” (rozdz. 2.2) |
| Pięć sekcji formularza obecnych w kolejności 1–5 z polami zgodnymi z rozdz. 4 | porównanie z `design/05-okna/platformowe/nowy-projekt.html` |
| Panel podsumowania niesie dokładnie osiem pozycji, przelicza się z formularzem | wypełnianie kolejnych sekcji i obserwacja panelu |
| Panel podsumowania jest przyklejony powyżej 1240 px i statyczny poniżej | zmiana szerokości okna wokół progu 1240 px |
| Trzy przyciski panelu wywołują różne zestawy komend zgodnie z rozdz. 12 i scenariuszami rozdz. 15 | wykonanie scenariuszy A, B, C i podgląd wywołań kontraktu |
| Błąd zakładania nie usuwa wprowadzonych danych | wymuszenie błędu (scenariusz D) i sprawdzenie zachowania formularza |
| Pola bez pokrycia kontraktowego są wykazane, nie ukryte | przegląd rozdz. 12 wobec `contract.json` |
| Materiał wejściowy powyżej 200 MB jest odrzucany z komunikatem, nie cichym pominięciem | próba wniesienia pliku przekraczającego limit (scenariusz F) |
| Projekt minimalny — wyłącznie nazwa — zakłada się bez błędu | wykonanie scenariusza E |
| Sekwencja wywołań przy „Załóż projekt i otwórz sesję” odpowiada diagramowi rozdz. 15 | przechwycenie wywołań kontraktu podczas zatwierdzenia i porównanie kolejności |
| Sesja założona `session.create` niesie puste `windowIds` bezpośrednio po założeniu | odczyt wyniku `session.create` zaraz po scenariuszu A, przed otwarciem przedsionka |
| Panel podsumowania nie niesie pozycji „Termin zamknięcia”, „Jednostka pracy” ani „Opis projektu” | porównanie ośmiu pozycji panelu (rozdz. 5) z jedenastoma polami formularza (rozdz. 4) |
| Wartości nadtytułu płótna dla TalkIn, CodeStudio i MultitaskingAI są rozstrzygnięte, nie pozostawione analogii | odczyt nadtytułu po wejściu z każdego z trzech środowisk (rozdz. 3.4 tabela wariantów) |
| Panel podsumowania stoi na żetonie `--dn-panel`, odróżnialnym od `--dn-powierzchnia` kart sekcji w obu motywach | porównanie koloru tła panelu i kart sekcji w motywie jasnym i ciemnym |
| Każda kontrolka osiągalna klawiaturą, ognisko widoczne | nawigacja `Tab` przez wszystkie pięć sekcji i panel |

---

*Koniec dokumentu. Okno nowego projektu — Specyfikacja docelowa, wersja 2.0, 2026-08-20.*

---
*Danaco Console — Platforma AI Workspace OS · v2.0 · status Deweloperski*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz.*
*Warunki korzystania: [Licencja produktu](../LICENSE.md). Kontakt: support@danaco-group.pl*
