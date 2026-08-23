# Ocena modułów przejętego klienta

Rozstrzygnięcie, co z katalogu `budowa/klient-poprzedni/src/moduly/` przechodzi
do nowego produktu bez zmian, co wymaga pracy i czego w przejętym kodzie nie ma.
Ocena nie buduje niczego — wyznacza granicę, poza którą kolejne tereny nie
powstają na przedmiot już istniejący.

Odniesieniem jest model produktu z [rejestru decyzji](decyzje.md): pozycja 5
(aplikacja → okno robocze → karta → widok; moduł jako praca z modelem na
narzędziu głównym; karta pomocnicza odrębna od modułu), pozycja 6 (kaskada
zakresu) i pozycja 9 (wstążka narzędziowa i przybornik karty). Miarą jest ta
sama miara, którą zastosowano w [ocenie warstwy przejętej](ocena-warstwy-przejetej.md).

## 1. Przedmiot i jego rozmiar

**Zlecenie nazywa przedmiot liczbą 47 159 wierszy. Katalog `moduly/` liczy
170 976 wierszy `.ts` w 759 plikach oraz 11 689 wierszy `.css` w 39 plikach.**
Liczba 47 159 nie jest rozmiarem katalogu — jest rozmiarem jego podzbioru
wyliczonego miarą szerszą z rozdziału 1 oceny warstwy przejętej („plik bez
jakiegokolwiek odwołania do przeglądarki"). Ocena obejmuje katalog w całości,
bo modułu nie da się rozstrzygnąć po jego podzbiorze.

Odtworzenie miary szerszej na tym drzewie daje **49 155 wierszy**, nie 47 159.
Zgadza się co do wiersza dla `apps` (5485), `browser` (3462), `design` (7045)
i `roundtable` (3321); rozjazd leży w `studio` — 8518 zmierzone wobec 8371
podanych. Rozbieżność wraca do Prowadzącego (rozdział 8).

Wykaz szesnastu modułów ze zlecenia zgadza się z drzewem co do nazwy i liczby.
Poza nim leży siedemnasty katalog — `wiedza` — którego zlecenie nie wymienia
(rozdział 5.17).

## 2. Jak mierzono

Każdy moduł oceniony uruchomieniem. Wynik pusty jest wszędzie odróżniony od
braku pomiaru.

| Przyrząd | Wywołanie | Wynik zbiorczy |
|---|---|---|
| żywy rdzeń | `danaco-console -port 17801 -dane /tmp/rdzen-moduly-dane -wymog-logowania false` | `rdzeń gotowy: komend=1077`, `transport: nasłuch 127.0.0.1:17801 gniazdo=/ws`, `zależności zewnętrzne: 30 z 30 obecnych` |
| wykaz komend rdzenia | `connection.hello` przez gniazdo `/ws`, odczyt pola `commands` | 1077 nazw, 66 przedrostków; zapisane i użyte jako miara pokrycia |
| sprawdzenie typów całości | `npx tsc --noEmit` w `budowa/klient-poprzedni` | kod wyjścia 0, czas 9,1 s |
| sprawdzenie typów modułu osobno | `npx tsc` z jawnym wykazem plików modułu i flagami z `tsconfig.json` | **0 błędów w każdym z szesnastu modułów** |
| sprawdziany | `npx vitest run src/moduly` | 30 plików, 404 sprawdziany, wszystkie zdane, 4,35 s |
| montaż modułu w przeglądarce | kod modułów zbundlowany bez przeróbki (`esbuild` 0.25.12, `--platform=browser`), serwowany z gniazda lokalnego, prowadzony przez Playwright/Chromium (`PLAYWRIGHT_BROWSERS_PATH=/opt/ms-playwright`), przeciw żywemu rdzeniowi | wyniki w rozdziale 3 |
| droga dziedzinowa modułu | źródła modułu (`zrodlo-*.ts`, katalogi czynności) wywołane wprost przeciw rdzeniowi | wyniki w rozdziale 5 |

Środowisko pomiaru: Node v22.23.2, TypeScript 5.8.3 (wersja przypięta
w `package.json` klienta, nie wersja maszyny), Go 1.26.5, Chromium
z `/opt/ms-playwright`.

**Potwierdzenie, że pomiar mierzy.** Sonda przeglądarkowa odczytuje najpierw
znacznik strony (`tytul === 'sonda wczytana'`) i obecność funkcji sondy
(`sonda_osadzona === true`); bez obu potwierdzeń zwraca „POMIAR NIEWYKONANY"
zamiast liczby. Każde wywołanie komendy ma własny licznik czasu i zwraca
`BRAK ODPOWIEDZI (Nms)` odrębnie od odmowy rdzenia i odrębnie od odpowiedzi
udanej o pustej treści.

**Punkt zaczepienia sondy.** Moduły wystawiają powłoce jeden kształt —
`OpisModulu { kod, utworzWidok(kanal) }` z `moduly/<nazwa>/indeks.ts`. Sonda
bierze go bez przeróbki, podstawia kanał nasłuchujący (przekazuje wywołania do
prawdziwego kanału i zapisuje nazwę komendy oraz rozstrzygnięcie), zakłada
oknu sesję i okno własne modułu komendą `window.create`, montuje widok
w dokumencie, woła `wczytaj(idSesji)` i po 1800 ms odczytuje stan.

**Miara dotknięcia drzewa dokumentu** jest miarą wąską z oceny warstwy
przejętej: `document.*`, `HTML*Element`, `querySelector`, `getElementById`,
`innerHTML`, `classList`, `appendChild`, `insertBefore`, `customElements`,
`ShadowRoot`.

**Podział na trzy klasy** (rozdział 6) rozdziela to, co miara wąska zlewa
w jedno:

- **dziedzina czysta** — plik bez drzewa dokumentu i bez kanału: dane,
  przeliczenia, katalogi pojęć;
- **dziedzina przy rdzeniu** — plik bez drzewa dokumentu, z kanałem: nazwy
  komend, kształty żądań, sprawdzian kształtu odpowiedzi;
- **obsługa interfejsu** — plik dotykający drzewa dokumentu.

## 3. Wynik zbiorczy

| Moduł | Plików | Wierszy | Dziedzina czysta | Dziedzina przy rdzeniu | Interfejs | Sprawdzianów | Komend | Okna: buduje / katalog rdzenia | Rozstrzygnięcie |
|---|---|---|---|---|---|---|---|---|---|
| `studio` | 133 | 42 954 | 5041 | 3851 | 29 100 | 264 | 228 | 8 / 8 | **wymaga pracy** |
| `design` | 66 | 15 370 | 1813 | 5548 | 7743 | 7 | 124 | 11 / 11 | **wymaga pracy** |
| `apps` | 46 | 10 754 | 1492 | 4329 | 4419 | 24 | 100 | 21 / 6 | **wymaga pracy** |
| `translate` | 49 | 9306 | 1028 | 1958 | 6320 | 0 | 73 | 7 / 7 | **wymaga pracy** |
| `automations` | 39 | 9161 | 2659 | 924 | 5095 | 31 | 64 | 6 / 6 | **wymaga pracy** |
| `terminal` | 30 | 9052 | 1748 | 625 | 6373 | 10 | 38 | 6 / 7 | **wymaga pracy** |
| `browser` | 50 | 8688 | 1945 | 2350 | 4393 | 0 | 58 | 5 / 6 | **wymaga pracy** |
| `roundtable` | 34 | 8543 | 2370 | 951 | 5222 | 0 | 46 | 6 / 7 | **wymaga pracy** |
| `agents` | 45 | 8398 | 630 | 1446 | 6171 | 5 | 53 | 6 / 6 | **wymaga pracy** |
| `library` | 48 | 8301 | 1372 | 1534 | 5164 | 9 | 60 | 5 / 6 | **wymaga pracy** |
| `multitasking` | 44 | 7815 | 1336 | 769 | 5710 | 0 | 41 | 4 / — | **wymaga pracy** |
| `research` | 50 | 7464 | 1650 | 2134 | 3361 | 13 | 105 | 7 / 6 | **wymaga pracy** |
| `assistant` | 43 | 6960 | 811 | 1496 | 4481 | 2 | 45 | 5 / 6 | **wymaga pracy** |
| `developer` | 23 | 6543 | 849 | 835 | 4541 | 12 | 51 | 16 / 5 | **wymaga pracy** |
| `diagnostics` | 25 | 5143 | 548 | 618 | 3401 | 21 | 24 | 12 / 6 | **wymaga pracy** |
| `workspace` | 27 | 5056 | 1989 | 800 | 2036 | 6 | 58 | 8 / 6 | **wymaga pracy** |
| razem | 752 | 169 508 | 27 281 | 30 168 | 103 530 | 404 | | | |

Trzy kolumny podziału nie sumują się do kolumny „Wierszy": różnicą są wiersze
sprawdzianów (8529 w 30 plikach), liczone osobno. Do sumy całego katalogu
brakuje jeszcze `wiedza` (526 wierszy) i czterech plików korzenia `moduly/`
(942 wiersze).

**Żaden moduł nie przechodzi bez zmian, i żaden nie jest nieobecny w całości.**
Powód jest jeden dla wszystkich szesnastu i nie leży w jakości kodu: **każdy
moduł jest złożeniem od czterech do dwudziestu jeden okien operacyjnych
ułożonych naraz w jednym obszarze**, a decyzja 5 stanowi, że funkcje wnosi się
do okna roboczego jako karty przełączane w paśmie, najwyżej dwie zestawione
w widoku dzielonym, plus okno boczne. Sześć paneli jeden pod drugim nie jest
pasmem kart. Rozstrzygnięcie „wymaga pracy" jest więc rozstrzygnięciem
o złożeniu modułu, a nie o jego dziedzinie — tę rozdział 6 wydziela osobno.

## 4. Montaż przeciw żywemu rdzeniowi

Wszystkie szesnaście modułów zamontowało się w Chromium przeciw rdzeniowi
na porcie 17801 **bez ani jednego wyjątku, bez ani jednego błędu strony
i bez ani jednego wpisu w konsoli błędów**. Uzgodnienie domknięte w całości:
`etapy = {powitanie:true, sesja:true, okno:true, gotowe:true}`.

| Moduł | Węzłów | Znaków tekstu | Komend wysłanych | ok | odmowa | brak odpowiedzi | Wyjątki |
|---|---|---|---|---|---|---|---|
| `studio` | 4086 | 61 138 | 20 | 19 | 1 | 0 | 0 |
| `design` | 2454 | 65 692 | 14 | 14 | 0 | 0 | 0 |
| `apps` | 1849 | 23 532 | 38 | 38 | 0 | 0 | 0 |
| `translate` | 1739 | 46 753 | 5 | 5 | 0 | 0 | 0 |
| `automations` | 695 | 9 061 | 10 | 10 | 0 | 0 | 0 |
| `terminal` | 870 | 15 357 | 9 | 9 | 0 | 0 | 0 |
| `browser` | 632 | 14 650 | 7 | 4 | 3 | 0 | 0 |
| `roundtable` | 609 | 21 234 | 3 | 3 | 0 | 0 | 0 |
| `agents` | 880 | 14 952 | 12 | 11 | 1 | 0 | 0 |
| `library` | 567 | 10 221 | 34 | 30 | 4 | 0 | 0 |
| `multitasking` | 633 | 11 208 | 8 | 8 | 0 | 0 | 0 |
| `research` | 973 | 40 631 | 2 | 2 | 0 | 0 | 0 |
| `assistant` | 801 | 30 312 | 18 | 17 | 1 | 0 | 0 |
| `developer` | 908 | 25 080 | 13 | 13 | 0 | 0 | 0 |
| `diagnostics` | 2088 | 72 620 | 4 | 4 | 0 | 0 | 0 |
| `workspace` | 407 | 6 093 | 5 | 3 | 2 | 0 | 0 |

Liczby węzłów i znaków są zależne od stanu rdzenia i nie są stałe między
przebiegami: `diagnostics` rysuje dziennik rdzenia, więc rośnie wraz z nim
(791 węzłów przy rdzeniu świeżym, 2088 po serii pomiarów), a `automations`
rysuje zapisane automatyki. Dwa pełne przebiegi sondy pod rząd dały zgodne
rozstrzygnięcie dla piętnastu modułów z szesnastu; rozbieżność w `automations`
to jedno wywołanie `config.effective.get` więcej, przy zerze błędów w obu
przebiegach. **Zero błędów, zero wyjątków i zero braków odpowiedzi powtórzyło
się w obu.**

Trzy moduły — `terminal`, `roundtable`, `developer` — montują się przez
`widokZOknaSesji` i **bez okna własnego w sesji oddają jedno zdanie zamiast
widoku**: „Ta sesja ma 1 okno/okien, ale żadne nie należy do modułu terminal —
moduł otworzy się razem ze swoim oknem". Jest to zachowanie zamierzone
i nazwane wprost w `moduly/rejestracja.ts`, nie usterka: stan pusty ma zdanie,
a nie ciszę. Po założeniu okna własnego komendą `window.create` wszystkie trzy
budują pełny widok (kolumna „Węzłów" wyżej pochodzi z przebiegu z oknem).

**Moduły sięgają przez stałe kontraktu (`Command.*`) po 976 różnych komend
i wszystkie 976 mają uchwyt w żywym rdzeniu.** Sprawdzone przecięciem wykazu
`Command.*` wyłuskanego z 759 plików katalogu z wykazem 1077 komend oddanym
przez `connection.hello`. Kontrakt niesie 1077 komend i rdzeń ma 1077 uchwytów —
rozjazdu nie ma; moduły wykorzystują 91 % powierzchni kontraktu. Nazwy spoza
kontraktu, których rdzeń nie ma, wymienia rozdział 5 przy właściwym module — są
to nazwy wpisane w kod jako powód niedziałania pozycji, a nie wywołania.

## 5. Rozstrzygnięcia moduł po module

### 5.1 `studio` — wymaga pracy

Praca z modelem w edytorze dokumentu: cyfryzacja i rozpoznanie pisma, redakcja
treści na kartce o nastawach wydania, warsztat PDF, znakowanie i komentarze,
blokady fragmentów, dziennik czynności, repozytorium wersji, wykaz operacji
kontekstowych, pętla wykonawcza.

**133 pliki, 42 954 wiersze, 264 sprawdziany w 16 plikach — najlepiej obłożony
sprawdzianami moduł katalogu.** Sprawdzenie typów osobno: 0 błędów.

**Cykl życia dokumentu przeciw żywemu rdzeniowi, drogą własną modułu**
(`zrodlo-studio.ts`, `zrodlo-kontroli-studio.ts`, `zrodlo-dokumentu-studio.ts`):

```
studio.document.open        ok / klucze=[document]   → studio-dok-49-bd1facbdfa7b3954
studio.document.save        ok / klucze=[document]        (dwa zapisy)
studio.repository.list      ok / klucze=[versions] versions.dlugosc=0
studio.version.series.list  ok / versions=0 operatorCount=0 autosaveCount=0
studio.journal.list         ok / klucze=[actions,total] actions.dlugosc=0
studio.autosave.get         ok / klucze=[settings]
studio.autosave.run         ok / saved=true, version=studio-wer-a8-7cf90fe5b23dd057
studio.repository.list      ok / versions.dlugosc=1   (po autozapisie)
studio.backup.list          ok / klucze=[backups,unsavedCount]
studio.markup.type.list     ok / markupTypes.dlugosc=3
studio.lock.add             ok / klucze=[lock,locks]
document.text.extract       ok / klucze=[text,usedOcr]
document.convert            ok / klucze=[asset,sizeBytes]
```

**Katalogi czynności modułu przeciw rdzeniowi.** `czynnosci-warsztatu.ts` (601
wierszy, 15 komend `studio.pdf.*` i `studio.security.*`) oraz
`czynnosci-redakcji.ts` (608 wierszy, 15 komend gałęzi, wydania, podglądu,
wsadu i wczytywania) opisują czynności danymi i składają żądanie funkcją
`zloz`. Sonda przeszła oba katalogi, składając żądanie funkcją modułu
i wysyłając je do rdzenia. **Wszystkie 30 komend trafiły na uchwyt rdzenia.**
Ani jedna odmowa nie brzmiała „rdzeń nie ma uchwytu komendy" — wszystkie były
sprawdzeniem dziedzinowym wartości podstawionych przez sondę:

```
studio.branch.list        ok / klucze=[branches] branches.dlugosc=0
studio.batch.run          ok / klucze=[runId,accepted,rejected]
studio.pdf.split          ODMOWA not_found — warsztat PDF: zasobu proba nie ma w magazynie
studio.security.sensitive.detect  ODMOWA validation_failed — kategorie nieznane:
                                  proba; rozpoznawane są: poczta,
                                  karta-platnicza, numer-ewidencyjny…
studio.ingest.device.scan ODMOWA internal_error — arsenał: SANE (scanimage) zakończył się
                          niepowodzeniem   (na maszynie pomiaru nie ma skanera)
```

**Zależność od modelu okna — nazwana wprost.** Jest jedna i jest środkiem
ciężkości pracy do wykonania.

`osadzenie-modulu.ts` (122 wiersze) wraz z `zrodlo-osadzenia.ts` (48 wierszy)
stawia nad wszystkimi oknami modułu **pole wyboru „Okno komunikacji sesji"**
z dymkiem: „Komendy modułu Studio działają na oknie komunikacji sesji: jego
identyfikator idzie w polu `windowId` komend `studio.document.open`
i `studio.contextual.op`". Moduł dostaje z powłoki wyłącznie identyfikator
sesji, więc sam czyta `window.list` i każe Operatorowi wskazać, w imieniu
którego okna pracuje.

Wobec decyzji 5 (`Session` = okno robocze, `Window` = karta czatu) znaczy to,
że Operator ma ręcznie wybrać kartę czatu, do której przypisze edytor. W nowym
modelu edytor i czat są dwiema kartami tego samego okna roboczego, a wiązanie
wynika z układu kart, nie z listy rozwijanej. Pole wyboru odpada wraz
z powodem swojego istnienia; zostaje potrzeba przekazania `windowId` — ta
zostaje, bo jest polem kontraktu.

**Złożenie do przebudowy.** `modul-studio.ts` (223 wiersze) składa osiem okien
operacyjnych — osadzenie, cyfryzację, pętlę wykonawczą, okno pracy z dokumentem,
Tools Panel, warsztat, redakcję i repozytorium — w sześć miejsc obszaru jeden
pod drugim (trzy okna wprost i trzy pasy po jednym lub dwóch oknach) i nadaje
elementowi etykietę „Moduł Studio — okna operacyjne". Sonda potwierdziła to
w drzewie dokumentu: osiem znaczników `data-okno`, 25 ram `section[aria-label]`,
4086 węzłów. Pod decyzją 5 te osiem funkcji to osiem kart pasma, z których dwie
da się zestawić w widok dzielony, a resztę odłożyć do okna bocznego.

**Rozjazd katalogu okien z rdzeniem.** Rdzeń wymienia dla Studia dziewięć kodów
okien operacyjnych: `chat-window`, `execution-loop-window`, `studio-editor`,
`diff-grep-panel`, `preview-window`, `session-repository`, `tools-panel`,
`ingest-ocr-panel`, `studio.document-workshop`. Moduł buduje osiem, z których
**pokrywa się jeden** (`studio.document-workshop`). Powód jest zapisany
w `modul-studio.ts` i jest świadomy: „Studio Editor, kanwa tekstowa, Preview
Window i Diff/Grep Panel pracowały nad TĄ SAMĄ treścią w czterech miejscach […]
Zeszły się w okno pracy z dokumentem". Cztery kody rdzenia mają w module jedno
okno o kodzie własnym `studio.praca-z-dokumentem`. Rozjazd jest usterką
katalogu rdzenia wobec modułu, nie usterką modułu — ale zostaje do
rozstrzygnięcia, bo nowy klient będzie czytał ten sam katalog.

**Co przechodzi.** 8892 wiersze w 45 plikach nie dotykają drzewa dokumentu:
5041 wierszy dziedziny czystej (`nastawy-strony.ts` 565 — nośniki druku,
marginesy, skala; `dziennik-czynnosci.ts` 284; `kategorie-operacji.ts` 176;
`wyszukiwanie-tekstu.ts` 133; `blokada-fragmentow.ts` 179; `kopie-zapasu.ts`
217; `ocena-redaktora.ts` 243) i 3851 wierszy dziedziny przy rdzeniu
(`zrodlo-kontroli-studio.ts` 872 — 44 typy odpowiedzi kontraktu; oba katalogi
czynności; jedenaście mniejszych źródeł). Ta masa jest kandydatem do przejęcia
niezależnie od tego, jak wygląda okno: wejściem są dane, wyjściem żądanie
kontraktu albo liczba.

**Usterki wykazane uruchomieniem** (rozdział 8 zbiera je z całości):

1. `studio.document.save` **nie zakłada wersji**. Trzy zapisy pod rząd zostawiły
   `studio.repository.list` z `versions: []`, a `studio.version.series.list`
   z zerami. Wersję zakłada dopiero `studio.autosave.run`. Kontrakt niesie
   `StudioDocumentSaveResponse.version?` jako pole opcjonalne; rdzeń go nie
   wypełnia. Moduł twierdzi przeciwnie w dwóch miejscach: `zrodlo-studio.ts`
   („`document.save` zapisuje treść i zakłada wersję") oraz zdanie wypisywane
   Operatorowi w `okno-pracy-z-dokumentem.ts` („Zapisz go — zapis zakłada
   pierwszą wersję sesji"). Zdanie jest nieprawdziwe wobec żywego rdzenia.
2. `studio.diff.compare` **bez wskazania dwóch wersji oddaje kopertę pustą**
   `{}` — ani `hunks`, ani `matches`, nawet jako tablice puste. Sprawdzone
   wzorcem obecnym w treści (`pattern: 'ma'` na treści „ala ma psa i kota")
   i wzorcem regularnym; obydwa razy `{}`. Z dwiema wersjami komenda działa
   w pełni: `hunks` z jednym fragmentem `changed` i `matches` z trafieniem.
   Wyszukiwanie wzorca po samej treści nie ma więc drogi do rdzenia.
3. `nastawy-strony.ts` wiersz 45 powołuje się na sprawdzian
   `nastawy-strony.test.ts`, którego w drzewie nie ma — ani w katalogu modułu,
   ani nigdzie w kliencie. Zdanie „Sprawdzian […] pilnuje, żeby wchłonięcie
   naprawdę nadpisywało wymiary wbudowane" nie ma pokrycia.
4. **Tools Panel dostaje z rdzenia trzy akcje, a żadna nie jest operacją
   redakcyjną.** `action.list` w zasięgu modułu (`scope: module`,
   `scopeId: 'studio'`, `enabledOnly`) oddaje dokładnie
   `studio.message.send`, `studio.message.stop`, `studio.message.list` — akcje
   okna komunikacji. `studio.operation.list`, wykaz własny Studia, oddaje na
   świeżym rdzeniu zero pozycji. Sam moduł nazywa ten stan wprost
   w `zrodlo-akcji-studio.ts` („Wierszy operacji kontekstowych Studia
   w katalogu nie ma; okno pokazuje to wprost, zamiast dorabiać je po stronie
   klienta") — jest to więc brak w katalogu rdzenia, nie w module.

   Droga wykonania operacji przy tym **działa**: `studio.contextual.op`
   z identyfikatorem z wykazu początkowego modułu
   (`kategorie-operacji.ts`, `OPERACJE_PASKA`) przechodzi sprawdzenie rdzenia
   i dochodzi do kanału modelu — `studio.korekta.ortografia`, `studio.styl.ton`
   i `studio.streszczenie.akapitowe` kończą się `channel_unavailable`
   z powodem „kanał lokalny-claude: injection: uruchomienie claude: exec:
   claude: executable file not found in $PATH", a nie `not_found`. Rdzeń
   przyjmuje identyfikator operacji od modułu i przekazuje go modelowi;
   niedostępny jest wyłącznie model.

### 5.2 `design` — wymaga pracy

Praca z modelem nad materiałem wizualnym: plansza, zasoby, podgląd, budowniczy
przywołań, sześć warsztatów (fotografia, wektor, druk, publikacja, makieta,
przeglądarka zasobów zewnętrznych).

**66 plików, 15 370 wierszy, 7 sprawdzianów.** Montaż: 2454 węzły, 65 692 znaki,
14 komend, wszystkie przyjęte. `design.asset.list` i `design.collection.list`
odpowiadają. Sprawdzenie typów: 0 błędów.

Moduł ma **najwyższy udział dziedziny przy rdzeniu w katalogu — 5548 wierszy
w 10 plikach** wobec 1813 wierszy dziedziny czystej. Buduje 11 okien
operacyjnych i pokrywa 10 z 11 kodów katalogu rdzenia (poza katalogiem stoi
`tokens-system-panel`, niezbudowany zostaje `execution-loop-window`).

Zależność od modelu okna: 72 wystąpienia `windowId`, żadnego „okna
komunikacji"; moduł wiąże się z oknem przez `window.list` i `window.state.get`.
Sześć nazw, których rdzeń nie ma jako komend — `design.mockup-workshop`,
`design.photo-workshop`, `design.print-workshop`, `design.publication-workshop`,
`design.stock-browser`, `design.vector-workshop` — to **kody okien
operacyjnych, nie komendy**; rdzeń wymienia je w `operationalWindowCodes`
modułu Design. Nie jest to brak.

### 5.3 `apps` — wymaga pracy

Budowanie produktów: budowniczy produktu, projektant architektury, przestrzeń
przodu i tyłu, panel wdrożenia, katalog aplikacji, zarządca zainstalowanych,
węzeł integracji, centrum uprawnień, konsola łącznika MCP, panel wydawcy.

**46 plików, 10 754 wiersze, 24 sprawdziany w 2 plikach.** Montaż: 1849 węzłów,
**38 komend wysłanych, wszystkie przyjęte** — najruchliwszy moduł katalogu przy
wczytaniu. Sprawdzenie typów: 0 błędów.

**Buduje 21 okien wobec 6 w katalogu rdzenia.** Dziesięć nadwyżkowych to kody
kart pomocniczych wspólne z `developer` i `diagnostics`: `podglad-bash`,
`terminal`, `terminal-tabs`, `output-console`, `process-monitor`,
`historia-rozmowy`, `pliki`, `przegladarka`, `artefakty`, `pliki-srodowiska`.
Nie są one kopiami w module — pochodzą ze wspólnego rejestru
`src/okna-pomocnicze/` (14 plików, 2206 wierszy) leżącego **poza terenem tej
oceny**. Jest to bezpośrednie odniesienie do decyzji 5: karta pomocnicza jest
w przejętym kliencie zbudowana i wspólna, tyle że wnoszona do modułu jako stały
panel, a nie jako karta pasma.

### 5.4 `translate` — wymaga pracy

Praca z modelem nad tłumaczeniem: panel źródła, panele tłumaczeń, glosariusz,
pamięć tłumaczeń, studio formatów, centrum jakości i przeglądu, warsztat.

**49 plików, 9306 wierszy, sprawdzianów własnych ZERO.** Montaż: 1739 węzłów,
46 753 znaki, 5 komend przyjętych. Sprawdzenie typów: 0 błędów. Buduje 7 okien;
6 pokrywa katalog rdzenia, `translation-workshop` stoi poza nim.

Droga dziedzinowa przeciw rdzeniowi (rdzeń niesie 64 komendy `translate.*`):

```
translate.source.set    ok / sourceLanguage="" segmentCount=2 panels=[]
translate.step.list     ok / steps=[{command:"translate.source.set",
                          name:"Ustal tekst źródłowy", kind:"translation", …}]
translate.glossary.list ok / terms=[] total=0
translate.memory.list   ok / entries=[] total=0
```

`translate.source.set` policzył segmenty prawdziwego zdania (2 z „Dzień dobry.
To jest zdanie próbne."), więc dziedzina tłumaczenia stoi w rdzeniu i moduł ma
do niej drogę. **Brak sprawdzianów przy 9306 wierszach jest usterką do
zgłoszenia**: nie ma czym wykazać, że przeniesienie czegokolwiek stąd niczego
nie zepsuło.

### 5.5 `automations` — wymaga pracy

Budowanie automatyk: wykaz gotowych pętli, budowniczy przepływu, orkiestrator,
harmonogram, zarządca kolejki, monitor wykonania.

**39 plików, 9161 wierszy, 31 sprawdzianów w 2 plikach.** Montaż: 695 węzłów,
10 komend przyjętych. Sprawdzenie typów: 0 błędów. Droga dziedzinowa działa
w całości przeciw żywemu rdzeniowi:

```
automation.workflow.save ok / workflow={id:"automat-cg-0bd2f2ac1b0ec417",
                              name:"sonda", enabled:true, version:1}
automation.workflow.list ok / workflows=[…jeden zapisany…]
automation.template.list ok / templates=[]
```

**Usterka źródła.** Rdzeń wymienia `automations` wśród piętnastu modułów
(`module.list` bez środowiska), ale z **pustym przypisaniem do środowisk**
(`environmentCodes: []`). Decyzja 5 stanowi: „Środowisko zawęża pole dostępnych
narzędzi do swoich modułów". Moduł nieprzypisany do żadnego środowiska nie ma
drogi przez okno startowe → środowisko → centrum dowodzenia; zostaje mu tylko
pasek szybkiego dostępu. Klient buduje dla niego pełny widok.

Drugi rozjazd: rodzaj modułu w rdzeniu to `kompozytor`, nie
`srodowisko_robocze` — a decyzja 5 mówi o oknie roboczym modułu wytwarzającego
mechanizm, że sesji nie niesie. Rozróżnienie istnieje w rdzeniu i klient go nie
używa.

### 5.6 `terminal` — wymaga pracy

Praca z modelem w powłoce: karty terminala, konsola wyjścia, monitor procesów,
zarządca sesji, harmonogram zadań, biblioteka skryptów.

**30 plików, 9052 wiersze, 10 sprawdzianów.** Montaż z oknem własnym: 870
węzłów, 15 357 znaków, 9 komend przyjętych (`terminal.session.list`,
`terminal.process.list`, `terminal.host.list`, `terminal.key.list`,
`terminal.tunnel.list`, `terminal.script.list`). Sprawdzenie typów: 0 błędów.

**Droga dziedzinowa domknięta przeciw żywemu rdzeniowi** — powłoka na maszynie
pomiaru, nie atrapa:

```
terminal.session.open  ok / session={id:"term-cu-76686dce10556a22", shell:"bash",
                         status:"running", workingDir:"…/sesje/ses_5f73cc4a76499c18"}
terminal.command.exec  ok / process={id:"tproc-cv-…", pid:3512662, parentPid:3494701,
                         command:"echo POMIAR-TERMINALA-OK && pwd", status:"running"}
terminal.output.read   ok / stdout="POMIAR-TERMINALA-OK\n\n/opt/…/sesje/ses_…\n
                         [proces finished, kod wyjścia 0]", exitCode=0, status:"finished"
```

**Usterka źródła — ta sama co przy `automations`:** rdzeń wymienia `terminal`
jako moduł z ośmioma oknami operacyjnymi i **pustym przypisaniem do środowisk**.

**Zależność od modelu okna — nazwana wprost i podwójna.** Moduł 30 razy nazywa
sesję powłoki „sesją", a 6 razy „kartą" (`karta`, `karty`). Rdzeń nazywa ją
kartą: `terminal.command.exec` odmawia zdaniem „moduł Terminal: wykonanie
polecenia wymaga wskazania karty", a kontrakt opisuje pole
`TerminalSessionOpenRequest.title` jako „Nazwa karty". Nakłada się to na
`Session` = okno robocze i `Window` = karta czatu z decyzji 5: w jednym module
słowo „karta" znaczy trzy różne rzeczy — kartę pasma okna roboczego, kartę
czatu i kartę powłoki. Ujednolicenie słownika jest pracą do wykonania przed
przeniesieniem.

Pięć nazw, których rdzeń nie ma — `terminal.dotfiles.sync`,
`terminal.knownhosts.get`, `terminal.notification.set`, `terminal.pane.split`,
`terminal.schedule` — to nazwy podane mechanizmowi `pokrycie-komend.ts` jako
powód niedziałania pozycji, nie wywołania. Mechanizm działa: rdzeń odpowiada na
nieznaną komendę `not_found — rdzeń nie ma uchwytu komendy X`, co sonda
potwierdziła.

### 5.7 `browser` — wymaga pracy

Praca z modelem na przeglądarce: okno przeglądarki, panel źródeł, panel
notatek, Automation Studio, panel przechwyceń i monitora, macierz izolacji.

**50 plików, 8688 wierszy, sprawdzianów własnych ZERO.** Montaż: 632 węzły,
7 komend, 4 przyjęte i 3 odmówione. Sprawdzenie typów: 0 błędów.

Trzy odmowy przy montażu **nie są usterką modułu** — wykazał to pomiar drogi
dziedzinowej. Komendy `browser.source.list`, `browser.note.list`
i `browser.snapshot.get` wymagają okna, po którym rdzeń wykonał już nawigację:

```
przed nawigacją:  browser.source.list  ODMOWA not_found — moduł Browser:
                    nie ma okna przeglądania o wskazaniu: okn_870e5398c90f1520
browser.navigate  ok / snapshot={id:"migawka-bb-…", url:"https://example.com",
                    title:"Example Domain", text:"Example…"}
po nawigacji:     browser.source.list   ok / sources=[]
                  browser.snapshot.get  ok / snapshot={…"Example Domain"…}
                  browser.note.list     ok / notes=[]
```

Rdzeń wyszedł do sieci i wrócił z tytułem strony. **Dziedzina przeglądarki
działa w całości.** Odmowa przed nawigacją jest stanem prawdziwym, nie brakiem.

Moduł buduje 5 okien z 6 katalogu rdzenia. Decyzja 9 wyznacza mu granicę wobec
karty pomocniczej: pełna konsola strony i oprzyrządowanie analityczne należą do
modułu, przybornik podstawowy do karty. Podział przedmiotów między
`moduly/browser/` a `src/okna-pomocnicze/` (kod `przegladarka`) tej granicy
dzisiaj nie odzwierciedla i jest pracą do wykonania.

**Brak sprawdzianów przy 8688 wierszach jest usterką do zgłoszenia.**

### 5.8 `roundtable` — wymaga pracy

Praca wielu modeli przy jednym stole: panele modeli, panel debaty, mapa
argumentów, centrum głosowania i oceny, panel moderatora, panel uzgodnienia.

**34 pliki, 8543 wiersze, sprawdzianów własnych ZERO.** Montaż z oknem własnym:
609 węzłów, 21 234 znaki, 3 komendy przyjęte. Sprawdzenie typów: 0 błędów.
Droga dziedzinowa:

```
roundtable.debate.get  ok / snapshot={windowId:"okn_…", participants:[],
                         turns:[], statements:[]}, total=0
roundtable.model.list  ok / participants=[]
```

**Zależność od modelu okna — najgłębsza w katalogu.** Decyzja 5 stanowi:
„Kilka kart czatu to kilka modeli pracujących w tej samej sesji. Podział widoku
pozwala postawić dwa czaty obok siebie […] To jest sedno produktu, nie funkcja
dodatkowa". Roundtable jest dokładnie tym mechanizmem, ale zbudowanym jako
sześć paneli w jednym obszarze modułu — a nie jako karty czatu okna roboczego
zestawione w widoku dzielonym. Moduł dwa razy odwołuje się do kodu
`chat-window` (jedyny moduł katalogu, który to robi). Przebudowa dotyczy tu
nie ozdoby, lecz przedmiotu: to samo, co Roundtable robi we własnym obszarze,
w nowym modelu robi pasmo kart okna roboczego.

**Brak sprawdzianów przy 8543 wierszach jest usterką do zgłoszenia.**

### 5.9 `agents` — wymaga pracy

Budowanie ekspertów: budowniczy agenta, konfiguracja modelu, centrum uprawnień,
zarządca umiejętności, zarządca łączników.

**45 plików, 8398 wierszy, 5 sprawdzianów.** Montaż: 880 węzłów, 12 komend,
11 przyjętych. Sprawdzenie typów: 0 błędów. Buduje 6 okien, 5 pokrywa katalog
rdzenia; szóste to `braki-modulu` — okno własne wykazu braków.

Jedyna odmowa jest usterką rdzenia, nie modułu:

```
tools.scope.list {agentId:"x"} ODMOWA not_found — profil asystenta domyślny nie istnieje
tools.scope.list {}            ODMOWA not_found — profil asystenta domyślny nie istnieje
```

Rdzeń świeżo założony nie ma domyślnego profilu asystenta, więc komenda odmawia
niezależnie od treści żądania. Moduł nie ma jak tego obejść.

**Najniższy udział dziedziny czystej w katalogu: 630 wierszy na 8398** — moduł
jest w przeważającej części formularzem.

### 5.10 `library` — wymaga pracy

Repozytorium wiedzy: przeglądarka repozytorium, podgląd pliku, panel wersji,
znaczniki i zbiory, panel danych opisowych i archiwum.

**48 plików, 8301 wierszy, 9 sprawdzianów.** Montaż: 567 węzłów, 34 komendy
wysłane, 30 przyjętych. Sprawdzenie typów: 0 błędów.

**Cztery odmowy to jedna usterka rdzenia, powtórzona czterokrotnie** — moduł
odczytuje panel czterema drogami i za każdym razem dostaje ten sam błąd:

```
library.stats.get {}  ODMOWA internal_error — dane: nie można policzyć zasobów
  repozytorium: sql: Scan error on column index 2, name "SUM(CASE WHEN suma_kon…"
```

Błąd jest odtwarzalny na żądaniu pustym, przy pustym repozytorium, poza
modułem. Jest to usterka rdzenia w zapytaniu SQL, nie usterka klienta.

`library.file.list`, `library.tag.list`, `library.retention.list`,
`library.rule.list`, `library.webhook.list`, `library.share.list`
i `library.suggestion.list` odpowiadają poprawnie.

Zależność od modelu okna: 8 wystąpień „okna komunikacji". Moduł montuje się bez
okna własnego, ale wypisuje wtedy zdanie: „Rdzeń oddał okna sesji (1), ale żadne
nie należy do modułu Library — moduł nie działa w imieniu żadnego z nich i nie
ma okna źródłowego przeniesienia kontekstu". Stan pusty ma zdanie.

### 5.11 `multitasking` — wymaga pracy

Orkiestracja: czat koordynatora, dwa czaty wykonawców, analizator wyników,
nadania ról, granice podagentów, monitor.

**44 pliki, 7815 wierszy, sprawdzianów własnych ZERO.** Montaż: 633 węzły,
8 komend przyjętych (`config.session.get`, `subagent.list`, `monitor.status`,
`role.list`). Sprawdzenie typów: 0 błędów.

**Moduł, który modułem nie jest.** Jego kod to `multitaskingai` — kod
**środowiska**, nie modułu. Potwierdza to rdzeń: `home.enter` oddaje środowisko
`{"id":"4","code":"multitaskingai","name":"MultitaskingAI",
"navigationKind":"orchestration"}` **bez pola `moduleCodes`**, a `module.list`
dla tego środowiska oddaje zero modułów. Klient wie o tym i mówi to wprost
w `aplikacja/rejestr-modulow.ts`: „Wpis MultitaskingAI niesie kod
`multitaskingai`, a to kod środowiska […] więc `opisModulu` tego wpisu
z nawigacji nie znajdzie". **Widok jest zbudowany i nieosiągalny drogą
nawigacji.**

Rdzeń dokłada drugi rozjazd: komendy `orchestration.*` (8 sztuk) odmawiają
zdaniem **„moduł Automations: komenda bez wskazania automatyki"**. Warstwa,
na której multitasking stoi, jest w rdzeniu przypisana Automations.

Zależność od modelu okna: 26 wystąpień „karty sesji" i 155 wystąpień rdzenia
„sesj" — moduł mówi językiem rdzenia, czyli językiem odwróconym wobec
decyzji 5.

### 5.12 `research` — wymaga pracy

Badanie: przestrzeń badania, panel odkryć, zarządca źródeł, widok czytania,
panel ustaleń, budowniczy raportu, panel wydania.

**50 plików, 7464 wiersze, 13 sprawdzianów.** Montaż: 973 węzły, 40 631 znaków,
2 komendy przyjęte. Sprawdzenie typów: 0 błędów. Droga dziedzinowa:

```
research.workspace.get  ok / scope="" stages=[] updatedAt=1787526234173
research.source.list    ok / sources=[] total=0
research.finding.list   ok / findings=[] total=0
```

**Dwadzieścia trzy nazwy `research.*` nie są komendami — są identyfikatorami
akcji okna** z `akcje-okien.ts`, a każda niesie pole `droga` rozstrzygające,
gdzie się wykonuje. Czternaście ma `droga: 'okno'` i wykonuje się w oknie, bez
udziału rdzenia. Dziewięć ma `droga: 'akcja'` i jedzie do rdzenia komendą
`window.action` — a katalog akcji rdzenia ich nie zna:

```
window.action research.discovery.filters  ODMOWA not_found — window.action:
                            akcja research.discovery.filters nie istnieje
                            w katalogu akcji
```

Ta sama odmowa wraca dla `research.discovery.providers`,
`research.export.content`, `research.export.formats`, `research.export.target`,
`research.finding.compare`, `research.reading.find`, `research.scope.kanban`
i `research.scope.timeline`. **Dziewięć czynności modułu Research nie ma
wiersza w katalogu akcji rdzenia** — brakuje wiersza katalogu, nie uchwytu
komendy, i moduł mówi to Operatorowi zamiast milczeć.

### 5.13 `assistant` — wymaga pracy

Asystent: konsola głosowa, monitor czynności, strumień aktywności, zarządca
pamięci i kontekstu, węzeł poleceń i narzędzi.

**43 pliki, 6960 wierszy, 2 sprawdziany — najsłabsze obłożenie wśród modułów,
które sprawdziany w ogóle mają: jeden na 3480 wierszy.** Montaż: 801 węzłów, 18 komend,
17 przyjętych. Sprawdzenie typów: 0 błędów.

Moduł mówi wprost, czego mu brakuje: „Rdzeń nie oddał okna modułu Assistant dla
tej sesji (`window.list`). `assistant.voice.command` wymaga pola `windowId`,
więc…". Odmowa `memory.list` przy montażu jest odmową dziedzinową rdzenia:
„moduł Workspace: karta sesji ses_… nie pracuje w żadnym projekcie".

**Zależność od modelu okna — 42 wystąpienia „karty sesji", najwięcej
w katalogu.** Moduł przejmuje słownik rdzenia, w którym `Session` jest „kartą
sesji" — czyli dokładnie odwrotnie niż w decyzji 5, gdzie kartą jest `Window`,
a `Session` jest oknem roboczym.

### 5.14 `developer` — wymaga pracy

Programowanie: edytor kodu, drzewo projektu, panel Git, wyjście budowy, warsztat
kodu, narzędzia dewelopera.

**23 pliki, 6543 wiersze, 12 sprawdzianów.** Montaż z oknem własnym: 908 węzłów,
13 komend przyjętych. Sprawdzenie typów: 0 błędów. Droga dziedzinowa działa na
prawdziwym systemie plików:

```
developer.tree.get {path:"/tmp"}  ok / root="/tmp" nodes=[{path:"/tmp",
                                    kind:"directory", updatedAt:…}, …]
developer.git.log  {path:"/tmp"}  ODMOWA not_found — moduł Developer: katalog
                                    roboczy okna nie jest repozytorium
```

Obie odpowiedzi są prawdziwe wobec stanu maszyny.

**Buduje 16 okien wobec 5 w katalogu rdzenia** — dziesięć nadwyżkowych to karty
pomocnicze ze wspólnego rejestru `src/okna-pomocnicze/` (jak w `apps`), jedenaste
to `warsztat-kodu`.

### 5.15 `diagnostics` — wymaga pracy

Diagnostyka: centrum diagnostyki, przeglądarka dzienników, panel błędów, panel
zaleceń, narzędzia obserwowalności.

**25 plików, 5143 wiersze, 21 sprawdzianów w 2 plikach.** Montaż: 2088 węzłów —
moduł rysuje dziennik rdzenia, więc liczba rośnie wraz z nim; przy rdzeniu
świeżym było ich 791. Cztery komendy, wszystkie przyjęte. Sprawdzenie typów:
0 błędów. Droga dziedzinowa czyta prawdziwe dzienniki rdzenia:

```
diagnostics.log.query  ok / entries=[{id:"log-b7-…", level:"info", source:"rdzeń",
                         message:"transport: rozłączenie pol-18 (kanał …"}, …]
diagnostics.error.list ok / errors=[{id:"blad-bv-…", fingerprint:"6463b55f…",
                         message:"moduł Developer: katalog roboczy okna nie jest…"}]
```

Drugi zapis jest odmową, którą sam pomiar wywołał minutę wcześniej — moduł
odczytał ślad własnego przebiegu sondy. **Buduje 12 okien wobec 6 w katalogu
rdzenia**; sześć nadwyżkowych to karty pomocnicze ze wspólnego rejestru.
Dziedziny ma 1166 wierszy na 5143 — 23 %, przy 21 % w Studiu najniższe wartości
w katalogu.

### 5.16 `workspace` — wymaga pracy

Projekt: pulpit projektu, panel instrukcji, pamięć kontekstu, biblioteka
projektu, zarządca agentów, węzeł planowania, wiki, współpraca.

**27 plików, 5056 wierszy, 6 sprawdzianów.** Montaż: 407 węzłów, 5 komend,
3 przyjęte. Sprawdzenie typów: 0 błędów.

Dwie odmowy są odmowami dziedzinowymi rdzenia, prawdziwymi wobec stanu:
`memory.list` → „moduł Workspace: karta sesji ses_… nie pracuje w żadnym
projekcie". Bez projektu nie ma pamięci do wyliczenia.

**Najwyższy udział dziedziny w katalogu: 2789 wierszy na 5056, czyli 55 %** —
odwrotnie niż wszędzie indziej. Moduł jest w przeważającej części modelem
projektu, a nie formularzem. To czyni go **najlepszym kandydatem do przejęcia
dziedziny w całości**.

Zależność od modelu okna: 27 wystąpień „karty sesji"; `windowId` nie występuje
ani razu — moduł pracuje na poziomie sesji i projektu, nie okna. Wobec decyzji 5
(`Session` = okno robocze) i decyzji 6 (kaskada: aplikacja → projekt → sesja →
karta) jest to poziom, który w nowym modelu zostaje bez zmian co do przedmiotu;
zmienia się nazwa.

### 5.17 `wiedza` — istnieje i jest nieosiągalny

Katalog spoza wykazu zlecenia: 3 pliki `.ts`, 526 wierszy, wyszukiwanie po
znaczeniu (`knowledge.*`). Sprawdzianów brak.

**Nie ma do niego drogi.** Jego własny nagłówek stanowi: „Okno wchodzi do sceny
przez złożenie modułu Poczta (`moduly/poczta/indeks.ts`), które ma w rejestrze
jeden wpis" — **katalog `moduly/poczta/` w drzewie nie istnieje**. Sprawdzenie
importów: żaden plik klienta nie importuje z `moduly/wiedza`. Katalog nie
wystawia `MODUL`, więc nie ma go w rejestrze modułów.

Rdzeń niesie 2 komendy `knowledge.*` i 10 komend `mail.*`, a `module.list` nie
wymienia ani modułu wiedzy, ani modułu poczty. **Jest to kod martwy** — 526
wierszy bez wejścia i bez sprawdzianu.

## 6. Logika dziedzinowa a obsługa interfejsu

Kryterium 4 zlecenia rozstrzygnięte pomiarem trzech klas z rozdziału 2.

| Klasa | Plików | Wierszy | Udział |
|---|---|---|---|
| dziedzina czysta — bez drzewa dokumentu i bez kanału | 186 | 27 281 | 16,1 % |
| dziedzina przy rdzeniu — bez drzewa dokumentu, z kanałem | 117 | 30 168 | 17,8 % |
| obsługa interfejsu — dotyka drzewa dokumentu | 419 | 103 530 | 61,1 % |
| sprawdziany | 30 | 8529 | 5,0 % |

**57 449 wierszy w 303 plikach — 34 % katalogu — nie dotyka drzewa dokumentu
i jest kandydatem do przejęcia niezależnie od tego, jak wygląda okno.** Jest to
o 8652 wiersze więcej, niż daje na tych samych szesnastu modułach miara szersza
(48 797), bo miara szersza odrzuca plik za samo wspomnienie `window.` albo
`addEventListener` — także wtedy, gdy chodzi o nasłuch na gnieździe, a nie na
dokumencie.

Podział wewnątrz tej masy ma znaczenie praktyczne:

- **Dziedzina czysta (27 281 wierszy)** przenosi się bez żadnego warunku: są to
  katalogi pojęć, przeliczenia i dane. Największe skupiska: `studio` 5041
  (nastawy strony, dziennik czynności, blokady, statystyka różnicy),
  `automations` 2659, `roundtable` 2370, `workspace` 1989, `browser` 1945,
  `terminal` 1748.
- **Dziedzina przy rdzeniu (30 168 wierszy)** przenosi się pod jednym warunkiem:
  wraz z nią przenosi się kontrakt, bo pliki te są niczym innym niż nazwami
  komend, kształtami żądań i sprawdzianem kształtu odpowiedzi. Kontrakt zostaje
  (rdzeń zostaje), więc warunek jest spełniony. Największe skupiska: `design`
  5548, `apps` 4329, `studio` 3851, `browser` 2350, `research` 2134.
- **Obsługa interfejsu (103 530 wierszy)** nie przenosi się w tej postaci
  w żadnym module, bo w każdym buduje okna operacyjne ułożone naraz w obszarze,
  a nie karty pasma.

Rozkład jest bardzo nierówny i rozstrzyga o kolejności przejmowania:

| Moduł | Udział dziedziny | Wierszy dziedziny | Wierszy modułu |
|---|---|---|---|
| `workspace` | 55,2 % | 2789 | 5056 |
| `apps` | 54,1 % | 5821 | 10754 |
| `research` | 50,7 % | 3784 | 7464 |
| `browser` | 49,4 % | 4295 | 8688 |
| `design` | 47,9 % | 7361 | 15370 |
| `automations` | 39,1 % | 3583 | 9161 |
| `roundtable` | 38,9 % | 3321 | 8543 |
| `library` | 35,0 % | 2906 | 8301 |
| `assistant` | 33,1 % | 2307 | 6960 |
| `translate` | 32,1 % | 2986 | 9306 |
| `multitasking` | 26,9 % | 2105 | 7815 |
| `terminal` | 26,2 % | 2373 | 9052 |
| `developer` | 25,7 % | 1684 | 6543 |
| `agents` | 24,7 % | 2076 | 8398 |
| `diagnostics` | 22,7 % | 1166 | 5143 |
| `studio` | 20,7 % | 8892 | 42954 |

Studio ma najwięcej dziedziny w liczbach bezwzględnych (8892 wiersze)
i najmniej w udziale (21 %) — jest zarazem największym zasobem i największą
przebudową interfejsu w katalogu.

## 7. Zależność od modelu okna

Kryterium 3 zlecenia. Zależność występuje w każdym module i ma trzy postacie.

**Postać pierwsza — słownik odwrócony.** Kontrakt rdzenia nazywa `Session`
„sesją", a w 58 miejscach „kartą sesji"; `Window` nazywa „oknem komunikacji".
Rdzeń mówi tak samo w treści odmów: „moduł Workspace: karta sesji ses_… nie
pracuje w żadnym projekcie". Decyzja 5 stanowi odwrotnie: `Session` = okno
robocze, `Window` = karta czatu. Moduły przejmują słownik rdzenia. Rozkład
wystąpień:

| Moduł | „okno komunikacji" | „karta sesji" |
|---|---|---|
| `assistant` | 4 | **42** |
| `workspace` | 1 | **27** |
| `multitasking` | 0 | **26** |
| `studio` | **13** | 0 |
| `agents` | **13** | 2 |
| `library` | **11** | 1 |
| `terminal` | 4 | 6 |
| `automations` | 2 | 4 |
| `translate` | 3 | 3 |
| pozostałe | 0–1 | 0–1 |

Nie jest to nazewnictwo do poprawienia zamianą wyrazów. Zdanie „komendy modułu
Studio działają na oknie komunikacji sesji" jest po zamianie na „karcie czatu
okna roboczego" nadal rozstrzygnięciem o modelu — mówi, że moduł pracuje
w imieniu jednej karty czatu wskazanej ręcznie.

**Postać druga — wskazanie okna jako czynność Operatora.** `studio` stawia nad
oknami pole wyboru „Okno komunikacji sesji" (`osadzenie-modulu.ts`).
`terminal`, `roundtable` i `developer` montują się dopiero, gdy w sesji stoi
okno o ich `moduleId` (`widokZOknaSesji` w `moduly/rejestracja.ts`), a
`multitasking` — gdy jest sesja (`widokZSesji`). W nowym modelu okno robocze
niesie karty, a przypisanie edytora do czatu wynika z układu kart, nie z listy
rozwijanej ani z obecności okna o właściwym `moduleId`.

**Postać trzecia — moduł jako obszar okien, nie jako pasmo kart.** Zmierzone
w drzewie dokumentu po montażu: liczba znaczników `data-okno` i liczba ram
`section[aria-label]`. `studio` 8 okien i 25 ram, `apps` 21 okien i 19 ram,
`developer` 16 i 13, `diagnostics` 12 i 9, `design` 11 i 5. Rdzeń przy tym
wymienia w `operationalWindowCodes` **kod `execution-loop-window` dla każdego
z piętnastu modułów, i ani jeden moduł nie buduje okna o tym kodzie** — Studio
buduje `studio.execution-loop`, reszta nie buduje nic.

## 8. Wykaz braków

Zgłoszenia do Prowadzącego. Każde wykazane uruchomieniem, z przytoczonym
wynikiem.

**Braki w źródłach zlecenia i katalogu rdzenia**

1. **Rozmiar przedmiotu podany w zleceniu nie odtwarza się.** Miara szersza daje
   na całym katalogu 49 155 wierszy wobec podanych 47 159 (48 797 na samych
   szesnastu modułach); zgadza się co do wiersza dla `apps`, `browser`, `design`
   i `roundtable`, rozjazd leży w `studio` — 8518 wobec 8371. Sam katalog
   `moduly/` liczy 170 976 wierszy `.ts` w 759 plikach.
2. **`terminal` i `automations` są w rdzeniu modułami bez przypisania do
   któregokolwiek środowiska** (`environmentCodes: []`). Klient buduje dla nich
   pełne widoki (9052 i 9161 wierszy). Decyzja 5 wiąże dostępność modułu ze
   środowiskiem; nie rozstrzyga, co znaczy moduł spoza wszystkich środowisk.
3. **`multitaskingai` jest w kliencie modułem, a w rdzeniu środowiskiem.** Rdzeń
   oddaje je jako środowisko `navigationKind: "orchestration"` bez `moduleCodes`;
   `module.list` dla niego oddaje zero modułów. Widok (7815 wierszy) jest
   nieosiągalny drogą nawigacji, co klient sam odnotowuje.
4. **Katalog okien operacyjnych rdzenia rozjeżdża się z modułami.** Kod
   `execution-loop-window` stoi przy każdym z piętnastu modułów i nie jest
   budowany przez żaden. Katalog Studia wymienia cztery kody
   (`studio-editor`, `diff-grep-panel`, `preview-window`, `session-repository`),
   które moduł świadomie scalił w jedno okno o kodzie własnym. Katalog jest
   niespójny także wewnętrznie: osiem kodów Studia bez przedrostka i jeden
   z przedrostkiem (`studio.document-workshop`), w Design cztery bez i sześć z.
5. **Katalog akcji rdzenia nie niesie operacji, które moduły wykonują.**
   `action.list` w zasięgu modułu Studio oddaje trzy akcje okna komunikacji
   (`studio.message.send`, `studio.message.stop`, `studio.message.list`), a ani
   jednej operacji redakcyjnej; `studio.operation.list` oddaje zero pozycji.
   Dziewięć akcji modułu Research idących drogą `window.action` wraca odmową
   `not_found — akcja … nie istnieje w katalogu akcji`. Brakuje wierszy
   katalogu, nie uchwytów komend — obie strony klienta mówią to wprost.
6. **Rdzeń podstawia `agents` pod pusty `moduleId`.** `window.create`
   z `moduleId: ''` oddaje okno o `moduleId: "agents"`. Rdzeń przyjmuje też
   `moduleId` spoza katalogu bez odmowy: `window.create` z
   `moduleId: 'nieistniejacy-modul'` kończy się `status=ok`. Opis okna klienta
   (`okno-komunikacji/opis-okna.ts`) zostawia moduł pusty świadomie — i dostaje
   cudzy moduł.

**Usterki rdzenia wykazane uruchomieniem**

7. **`library.stats.get` kończy się błędem wewnętrznym przy każdym wywołaniu.**
   `internal_error — dane: nie można policzyć zasobów repozytorium: sql: Scan
   error on column index 2, name "SUM(CASE WHEN suma_kon…"`. Odtwarzalne na
   żądaniu pustym, przy pustym repozytorium, poza modułem. Moduł Library woła tę
   komendę czterokrotnie przy jednym wczytaniu i czterokrotnie dostaje ten sam
   błąd.
8. **`studio.document.save` nie zakłada wersji.** Trzy zapisy zostawiły
   `studio.repository.list` z `versions: []`. Wersję zakłada dopiero
   `studio.autosave.run`. Kontrakt niesie `StudioDocumentSaveResponse.version?`.
9. **`studio.diff.compare` bez dwóch wersji oddaje kopertę pustą.** `{}` — ani
   `hunks`, ani `matches`, także jako tablice puste; sprawdzone wzorcem obecnym
   w treści i wzorcem regularnym. Z dwiema wersjami komenda działa w pełni.
   Wyszukiwanie wzorca po samej treści nie ma drogi do rdzenia.
10. **`tools.scope.list` odmawia niezależnie od treści żądania:** `not_found —
   profil asystenta domyślny nie istnieje`. Rdzeń świeżo założony nie ma
   domyślnego profilu; moduł Agents nie ma jak tego obejść.
11. **`apps.architecture.get` oddaje kopertę pustą `{}`** — wynik udany bez ani
    jednego pola.

**Usterki klienta wykazane uruchomieniem**

12. **`studio/nastawy-strony.ts` wiersz 45 powołuje się na sprawdzian
    `nastawy-strony.test.ts`, którego w drzewie nie ma.**
13. **Zdanie wypisywane Operatorowi w `studio/okno-pracy-z-dokumentem.ts` —
    „Zapisz go — zapis zakłada pierwszą wersję sesji" — jest nieprawdziwe wobec
    żywego rdzenia** (punkt 8).
14. **Cztery moduły nie mają ani jednego sprawdzianu: `browser` (8688 wierszy),
    `translate` (9306), `roundtable` (8543), `multitasking` (7815) — razem
    34 352 wiersze.** Nie ma czym wykazać, że przeniesienie czegokolwiek stamtąd
    niczego nie zepsuło. `assistant` ma 2 sprawdziany na 6960 wierszy.
15. **`moduly/wiedza/` (526 wierszy) jest kodem martwym**: nikt go nie
    importuje, nie wystawia `MODUL`, a jego własny nagłówek wskazuje na
    `moduly/poczta/indeks.ts`, którego nie ma.
16. **Dziewięć czynności modułu Research nie ma wiersza w katalogu akcji
    rdzenia** (`research.discovery.filters`, `research.discovery.providers`,
    `research.export.content`, `research.export.formats`,
    `research.export.target`, `research.finding.compare`,
    `research.reading.find`, `research.scope.kanban`,
    `research.scope.timeline`). Każda wraca z `window.action` odmową
    `not_found`. Pozostałe czternaście nazw `research.*` z tej rodziny
    wykonuje się w oknie i wiersza w rdzeniu nie potrzebuje — to samo pole
    `droga` je rozróżnia.

## 9. Czego nie dało się sprawdzić

1. **Zachowania modułów pod obciążeniem i w czasie.** Sonda mierzy montaż
   i pierwszy odczyt (1800 ms na moduł). Zachowanie po godzinie pracy, przy
   zerwaniu połączenia i po wznowieniu sesji nie było mierzone.
2. **Dróg zapisujących, które zmieniają stan trwały.** Poza Studio, Terminalem
   i Automations sonda uruchamiała odczyty. Komendy `*.set`, `*.remove`,
   `*.delete` pozostałych modułów nie były wywołane.
3. **Skutku czynności wychodzących do modelu.** Droga jest zmierzona do końca
   po stronie rdzenia, ale odpowiedź modelu nie: `studio.contextual.op`
   z prawidłową pozycją rejestru akcji dochodzi do warstwy kanału i kończy się
   `channel_unavailable — kanał lokalny-claude: injection: uruchomienie claude:
   exec: claude: executable file not found in $PATH`. Powitanie rdzenia niesie
   przy tym `gatewayConfigured=false`. Rdzeń ma w rejestrze dwa kanały modelu
   (`lokalny-claude`, `bielik-ollama`), ale żaden nie jest na tej maszynie
   wykonywalny. Wszystko, co po dojściu do modelu miałoby wrócić — sedno każdego
   modułu — pozostaje niezmierzone.
4. **`studio.ingest.device.scan`** dochodzi do rdzenia i tam kończy się
   `internal_error — arsenał: SANE (scanimage) zakończył się niepowodzeniem`.
   Na maszynie pomiaru nie ma skanera; nie da się rozstrzygnąć, czy droga jest
   dobra.
5. **Wyglądu.** Ocena mierzy drzewo dokumentu (liczba węzłów, znaczniki okien,
   ramy), nie obraz. Zrzutów nie wykonywano — wygląd należy do warstwy
   projektowej, nie do modułów.
6. **Granicy wobec `src/okna-pomocnicze/`** (14 plików, 2206 wierszy). Karty
   pomocnicze z decyzji 5 są tam zbudowane i wspólne, a moduły `apps`,
   `developer` i `diagnostics` wnoszą je do siebie jako stałe panele. Katalog
   leży poza terenem tej oceny i nie był oceniany.
7. **Arkuszy stylów modułów** (39 plików, 11 689 wierszy `.css`). Należą do
   warstwy projektowej, ocenionej osobno.
