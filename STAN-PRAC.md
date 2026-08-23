# Stan prac — Danaco Console, doprowadzenie modułów do pełnej funkcjonalności

Zapis na wypadek urwania sesji. Data: 2026-08-17.

## Rachunek — stan bieżący

**Kontrakt liczy 1067 komend. Bez uchwytu nie zostaje ani jedna.**

Rejestr zmontowanego rdzenia odpowiada na każdą nazwę kontraktu i nie zna ani
jednej nazwy spoza niego. Na starcie tej roboty bez obsługi było **567 komend
z 876**.

### Jak to się mierzy — i czemu stara miara kłamała

Pomiar czytający źródła wyrażeniem `Zarejestruj\(shared\.(Command\w+)` jest
**wycofany**. Mylił się w obie strony:

- komendy wpinane przez parametr, a nie literałem, wypadały z pomiaru jako
  „bez obsługi" — tak przez cały czas wychodziło `action.list`, wpinane
  w `kompozycja.go` jako `zarejestrujAkcje(rejestr, p.Akcje, shared.CommandActionList)`,
  choć rdzeń odpowiada na nie od dawna;
- wywołanie `Zarejestruj` w gałęzi, do której montaż nigdy nie dochodzi, tamten
  pomiar liczył jako pokrycie.

Prawdę zna rejestr, bo to on rozstrzyga, czy komenda dostanie uchwyt, czy
odpowiedź `*.unknown`. Miarą jest więc sprawdzian na żywym rdzeniu — nie trzeba
go już przepisywać do dokumentu, bo pilnuje się sam:

```bash
cd /home/ubuntu/DanacoConsole/budowa && go test ./server/internal/core/ -run 'KomendaKontraktu|RejestrNieZna' -v
```

`server/internal/core/pokrycie_kontraktu_test.go` — dwa sprawdziany: każda
komenda kontraktu ma uchwyt, żaden uchwyt nie stoi pod nazwą spoza kontraktu.
Komenda dopisana do kontraktu bez obsługi wywraca ten sprawdzian od razu.

`provenance.call.replay` jest zarejestrowany świadomie inaczej niż reszta —
powód stoi w `server/internal/core/adapter_prowenancja.go`. Przed zmianą
przeczytać.

## Moduły zamknięte

Komplet komend w rejestrze i własne sprawdziany skutku wykonawcy:

**Studio** 63/63 · **Research** 76/76 · **Translate** 64/64 · **Browser** 47/47 ·
**Roundtable** 46/46 · **Library** 48/48 · **Developer** 50/50 ·
**Workspace** 40/40 · **Terminal** 27/27 · **Agents** 30/30 (+ orkiestracja 8/8,
narzędzia 3/3) · **Automations** 32/32 (+ kolejka 16/16, harmonogram 6/6) ·
**Apps** 41/41 (+ rozszerzenia 37/37) · **Diagnostics** 4/4

**Design** 29/29 · **Assistant** komplet, łącznie z `action.list` (katalog akcji
czyta rejestr `core/akcje_rejestr.go` nad tabelą `akcja` z migracji 008).

Tabela `akcja` jest przy tym PUSTA — migracja zakłada ją bez zaczynu. `action.list`
odpowiada więc wykazem pustym, a panele akcji i siatka szybkich akcji nazywają to
brakiem pozycji. To nie jest usterka rdzenia: wiersze katalogu są treścią
produktu i mają przyjść z opracowań, nie z kodu.

W Studiu warsztat PDF (9 komend) i bezpieczeństwo dokumentu (6) zbudowałem sam
na `pdfcpu` i `crypto` ze stdlib. W Developerze silnik repozytorium i 12 komend
też są moje.

## Co zbudowałem sam (nie oddawać nikomu bez powodu)

### `server/internal/repozytorium/` — pakiet samodzielny
`go-git` + `gotextdiff` + `regexp`. Stan, Roznica, Historia, Galezie, Konflikt,
RozstrzygnijKonflikt, Szukaj, Zamien, Zaloz, ZmienNazwe, Usun, Przenies.
12 sprawdzianów skutku, wszystkie zielone. Pakiet nie zależy od `dane`, więc
kompiluje się i mierzy nawet wtedy, gdy reszta drzewa jest w remoncie — i to się
wielokrotnie opłaciło.

### `server/internal/core/adapter_studio_pdf.go`
Dziewięć czynności rodziny `studio.pdf.*` na `pdfcpu`. Zero procesów potomnych.

### `server/internal/core/adapter_studio_bezpieczenstwo.go`
Sześć czynności `studio.security.*`. Redakcja WYCINA bloki tekstu z treści
strony (własna czytelnia składni treści PDF), dopiero potem kładzie czarne pole.
Podpis liczy skrót po TREŚCI STRON, nie po bajtach pliku — bajty zmieniają się
przy każdym zapisie, więc podpis liczony po nich nie dałby się zweryfikować.

### `server/internal/core/zapora_warsztatu_pdf_test.go` — ZAPORA, NIE OSŁABIAĆ
Pliki rodzin `studio.pdf.*` i `studio.security.*` nie mogą zawierać
`zewnetrzne.Wolaj` ani `exec.Command`; żaden plik rdzenia nie może wymieniać
`qpdf` ani `ghostscript`. Powód wpisany w treść pliku.

### `server/internal/core/adapter_modul_developer_git_odczyt.go` i `_wezly.go`
Dwanaście komend Developera. Odczyt repozytorium na `go-git`, bez programu `git`.

### `client/src/moduly/studio/` — okno „Warsztat dokumentu"
`okno-warsztatu-dokumentu.ts`, `czynnosci-warsztatu.ts` (katalog 15 czynności
opisanych danymi), `zrodlo-warsztatu-dokumentu.ts`, sprawdzian
`czynnosci-warsztatu.test.ts` (8 pozycji). Migracja **140** wnosi okno do
katalogu rdzenia.

## Rozstrzygnięcia obowiązujące wszystkich

1. **Zasada bezwzględna.** Żadna funkcja nie zależy od programu, którego instalka
   nie niesie. Arsenał stoi NA SERWERZE, u Operatora jest cienka instalka — samo
   okno. Gdy istnieje biblioteka Go — bierz bibliotekę: pdfcpu zamiast qpdf,
   go-git zamiast `git`, regexp zamiast ripgrep. PDF, kryptografia i podpis —
   wyłącznie biblioteki wkompilowane.
2. **Sprawdzian skutku, nie koperty.** `status: ok` z pustym wynikiem to wzorzec
   szkody, który w tym produkcie wystąpił (Design meldował wykaz zasobów, za
   którymi nie było ani jednego bajtu). Sprawdzian schodzi do bazy własnym
   zapytaniem SQL albo do pliku i mierzy niezależnie.
3. **Przedrostki nazw.** Pakiety `server/internal/dane` i `server/internal/core`
   mają jedną przestrzeń nazw, a pisze w nich dziewięciu wykonawców. Każda nazwa
   pomocnicza (`kolumny*`, `odczytaj*`, `lista*`, `usun*`, `baza*`, `zloz*`) ma
   nieść przedrostek obszaru. Dotyczy także plików `_test.go`. Dziewięć zatrzymań
   całego drzewa w tej turze miało tę jedną przyczynę.
4. **Bez rzutowania typu bez sprawdzenia w `Zloz`.** Rdzeń złożony bez portu ma
   powiedzieć, czego mu brakuje, i pracować dalej — tak jak `if r == nil { return }`
   w funkcjach `zarejestruj*`. Panika przy montażu kładzie cały produkt.

## Zakresy migracji (żeby nie było kolizji numerów)

140-149 moje · 150-159 Research · 160-169 Translate · 170-179 Browser ·
180-189 Library · 190-199 Roundtable · 200-215 Apps · 216-229 Workspace ·
230-245 Design · 246-259 Terminal · 260-275 Automations · 276-289 Agents ·
290-305 Assistant i platforma · 306-315 dokończenie Studia · 141-149 Developer

## Sprawy otwarte

- `apps.deployment.run` — sprawdzian `TestPakowaniePodpisIPublikacjaDajaPlikPodpisIPozycje`
  wypada niepomyślnie: okno `apps.product-builder/sprawdzian` nie istnieje.
  Należy do wykonawcy Apps.
- Silnik kontenerów **stoi na serwerze i tak zostaje** — rozstrzygnięcie
  Właściciela z 17.08.2026. Program jest potrzebny modułowi Developer, więc jest
  częścią aplikacji serwerowej. `developer.container.*`, `image.build`,
  `compose.up` na Docker SDK pracują sprawnie, nie odmawiają. Sprawa zamknięta.
- `studio.ingest.device.scan` — skaner stoi po stronie Operatora, nie serwera.
  Komenda ma odmawiać odmową nazwaną, nie udawać skanu.

## Jak sprawdzić stan po powrocie

```bash
cd /home/ubuntu/DanacoConsole/budowa
go build ./server/... && go vet ./server/internal/core/
go test ./server/internal/repozytorium/ -count=1
go test ./server/internal/core/ -run 'Zapora|Warsztat|Bezpieczen|Redakcja|Podpis|Szyfrowanie|Czyszczenie|Rozpoznanie|Pdf|Podzial|Przestawienie' -count=1
```

Ostatni pomiar: wszystkie 12 pozycji wyżej przechodzą, `repozytorium` przechodzi.

## Zagrożenie eksploatacyjne zgłoszone dwukrotnie

Dwaj wykonawcy niezależnie zgłosili, że **`go mod tidy` uruchomione przy
niekompilującym się drzewie wycina z `go.mod` zależności**, których nie widzi
w kodzie — zniknęły `pdfcpu`, `go-git`, `gotextdiff`, `x/image`, `docker`.
Obaj odtworzyli plik z lokalnej pamięci modułów.

Zasada na przyszłość: **`go mod tidy` wolno puścić dopiero wtedy, gdy
`go build ./...` przechodzi.** Przy pracy równoległej kilku wykonawców drzewo
bywa niekompilowalne przez kilkanaście minut i wtedy `tidy` kasuje cudze
zależności, a nie porządkuje własne.

Sprawdzenie po każdym takim biegu:

```bash
grep -cE "pdfcpu|go-git/go-git|gotextdiff|docker/docker|x/image|kin-openapi|pgx|go-sql-driver" go.mod
# ma dać 8
```


## Dobudowa kontraktu Designu — 17.08.2026, w toku

### Powód

Właściciel wskazał, że Design ma projektować grafikę użytkową: materiały
reklamowe, banery, grafikę wielkoformatową, biurową, sieciową i do druku.
Pomiar potwierdził dziurę: **opracowanie modułu wymienia 88 funkcji w dziesięciu
grupach, kontrakt niósł 29 komend**. Wykonawca zbudował 29 z 29 — czyli sto
procent tego, co kontrakt zamawiał. Brakowało samego kontraktu.

Grup nie było w ogóle: **C** (wektor), **D** (makiety interfejsu), **F** (kolor),
**G** (ikony i typografia), **I** (marketing i szablony), oraz — dołożona po
wskazaniu Właściciela — **K** (druk, wielki format, grafika biurowa).

### Co zostało wniesione do `shared/contract.json`

**+55 komend, +30 struktur, +23 wyliczenia.** Rodzina `design.*` urosła
z 29 do **84 komend**. Kontrakt: 876 → 931 komend.

Skrypt wnoszący leży w
`/tmp/claude-1000/-home-ubuntu/26f7cc8f-e2a2-47ab-b240-1640f7f43612/scratchpad/kontrakt-design.py`
wraz z kopią kontraktu sprzed zmiany (`contract.json.przed`). Skrypt pomija
pozycje o nazwach już obecnych, więc drugie uruchomienie niczego nie zmienia.

Nowe rodziny:

- `design.vector.*` — ścieżki Béziera, kształty, operacje logiczne, tekst na
  ścieżce i kontury, czyszczenie SVG, symbole z instancjami, wydanie
  wektorowe (SVG, PDF, EPS)
- `design.frame.*`, `design.layout.auto`, `design.constraint.set`,
  `design.component.*`, `design.prototype.*`, `design.mockup.*`,
  `design.grid.set` — makiety interfejsu
- `design.color.*` — palety z harmonii, ekstrakcja z obrazu, kontrast WCAG,
  symulacja wad widzenia barw, gradienty, konwersje przestrzeni, audyt
  dostępności całego zestawu żetonów
- `design.icon.*`, `design.favicon.build`, `design.font.*` — ikony i typografia
- `design.template.*`, `design.campaign.set.build`, `design.stock.*`,
  `design.product.mockup.render` — marketing i szablony
- `design.print.*`, `design.largeformat.tile`, `design.chart.render`,
  `design.diagram.render` — druk, wielki format, grafika biurowa

### Rozstrzygnięcia wpisane w treść kontraktu

- **`design.print.export` odmawia**, gdy kontrola przeddrukowa znajdzie wadę
  o wadze błędu. Plik nie do druku wydany jako gotowy do druku jest gorszy niż
  odmowa. Pominięcie kontroli jest jawnym wyborem Operatora i wraca w odpowiedzi
  polem `preflightSkipped`.
- **Kształt powstaje od razu jako węzły** ścieżki, więc da się go dalej edytować
  piórem — nie jest osobnym bytem, który trzeba potem zamieniać.
- **Zamiana tekstu w kontury jest nieodwracalna dla wyniku**, więc tekst
  źródłowy zostaje.
- **Bilans zamiast ciszy** w dziesięciu miejscach: `unreachableFrameIds`,
  `failedConcepts`, `gridWarnings`, `unmatchedNames`, `failedSizes`,
  `providersFailed`, `unplacedNodeIds`, `unrecognizedRegions`, `removedPathIds`,
  `iccProfiles`.
- **`design.chart.render` bierze serie danych, nie obraz** — wykres da się
  przerysować po zmianie liczb.

### Stan po wniesieniu

- `node shared/gen/generate.mjs` przepuszczony: `contract.ts` i `contract.go`
  odtworzone (931 komend, 138 zdarzeń, 228 narzędzi modelu).
- `go build ./shared/` czysty, `go test ./server/internal/store/` przechodzi
  (w tym `odwzorowanie_kontraktu_test.go`).
- `npx tsc --noEmit` w kliencie czysty, `kontrakt.test.ts` 8/8.

### Czego jeszcze NIE ma

**Ani jedna z 55 nowych komend nie ma uchwytu w rdzeniu.** Kontrakt jest
wniesiony i zwarty, budowa nie zaczęta. To jest następny krok: warstwa danych,
adaptery, migracje (wolny zakres: **316-345**), klient i sprawdziany skutku.

Biblioteki do wzięcia, wszystkie czysto Go, żadna nie wymaga programu
zewnętrznego: `tdewolff/canvas` (operacje logiczne na ścieżkach, wydanie PDF
i EPS), `lucasb-eyer/go-colorful` (harmonie, Lab, kontrast), `golang/freetype`
albo `go-text/typesetting` (kontury tekstu, glify), `ajstarks/svgo`
(serializacja SVG), `Kodeworks/golang-image-ico` (favicony), `pdfcpu` (PDF/X),
`golang.org/x/image/draw` (wyrys i kafle).


## Kto co buduje w tej chwili

| Wykonawca | Zakres | Komend | Migracje |
|---|---|---:|---|
| Design — wektor i makiety | `design.vector.*`, `design.frame.*`, `design.component.*`, `design.prototype.*`, `design.mockup.*`, `design.layout.auto`, `design.constraint.set`, `design.grid.set` | 24 | 316-325 |
| Design — kolor, ikony, typografia | `design.color.*`, `design.icon.*`, `design.favicon.build`, `design.font.*` | 15 | 326-335 |
| Design — druk, wielki format, marketing | `design.print.*`, `design.largeformat.tile`, `design.chart.render`, `design.diagram.render`, `design.template.*`, `design.campaign.set.build`, `design.stock.*`, `design.product.mockup.render` | 16 | 336-345 |
| Assistant i platforma | `action.*` — ostatnia komenda | 1 | 290-305 |

## Rzeczy do rozstrzygnięcia przez Właściciela

1. **Harmonogram zadań powłoki — czyj.** `docs/moduly/terminal.md` opisuje własne
   okno „Task & Schedule" z runnerem zadań powłoki i wyrażeniami cron, ale żadne
   opracowanie nie przypisuje sobie komend `queue.*` ani `schedule.*`. Kontrakt
   wiąże je z Automations i tak zostały zbudowane.
2. **Rozjazd kontraktu wobec opracowań w pozostałych modułach.** Design nie musi
   być jedynym miejscem, gdzie kontrakt jest węższy niż opracowanie. Katalogi
   funkcji: Library 136 wierszy, Diagnostics 136, Automations 135, Studio 132,
   Workspace 131, Browser 128. Kontrakty niosą od 27 do 84 komend na moduł.
   Część różnicy jest pozorna (jedna komenda obsługuje kilka pozycji katalogu,
   wiele funkcji jest czysto klienckich), część może być taka sama jak w Designie.
   **Pomiar nie został wykonany.**

## Defekty zastane, znalezione przy okazji i naprawione

- **Automatyka z harmonogramem nigdy się nie uruchamiała** — budzik zakładał
  kolejkę rodzaju `harmonogram`, a warunek tabeli z migracji 003 dopuszczał
  tylko `sesyjna` i `multitasking`. Naprawione migracją 269.
- **Odwołania do sekretów w krokach automatyki ginęły** — przychodziły żądaniem
  i nie były utrwalane; usunięcie klucza szło w ciszy.
- **Stan `paused` projektu odbijał się od warunku kolumny** — kontrakt go
  wymaga, tabela dopuszczała tylko `active` i `archived`. Migracja 226.
- **`design.asset.generate` nie zapisywał promptu** — `DesignAsset.PromptId`
  było polem, którego rdzeń nigdy nie wypełniał.
- **`sciezkaZasobu` w module Studio była zaślepką odmawiającą zawsze** — moduł
  nie miał wpiętego repozytorium zasobów, co blokowało też `ingest.queue.add`.
- **Panika przy montażu rdzenia** — rzutowanie portu bez sprawdzenia w `Zloz`.
  Trzy rzutowania przepisane na postać dwuwartościową.

## Sprawdzian niestabilny

`TestZasadaRetencjiNieKasujeWpisowWstecz` wypadł raz niepomyślnie, a w izolacji
i przy ponownym pełnym przebiegu przechodzi. Nie został zdiagnozowany.

# Studio — jedno okno pracy z dokumentem (scalenie czterech powierzchni)

Zlecenie: `ZLECENIE-STUDIO-JEDNO-OKNO.md` wraz z trzema uzupełnieniami
Właściciela (wzorzec pakietu biurowego, wstążka z panelem Redaktora i wierszem
polecenia, zakładka recenzji z dymkami komentarzy).

## Co scalone, co usunięte

Cztery powierzchnie tekstowe — **Studio Editor, kanwa tekstowa, Preview Window,
Diff/Grep Panel** — zeszły się w jedno okno `studio.praca-z-dokumentem`
(`client/src/moduly/studio/okno-pracy-z-dokumentem.ts`). Podgląd wydania
i różnica są trybami widoku tej samej treści, nie osobnymi oknami.

Usunięte pliki: `okno-studio-editor.ts`, `okno-kanwy-tekstowej.ts`,
`okno-diff-grep-panel.ts`, `pasek-narzedzi-tekstu.ts`, `pasek-zaznaczenia.ts`,
`decyzja-propozycji.ts`, `studio-kanwa.css`. Z `okno-preview-window.ts` zdjęty
wariant Studia; wariant Designu został nietknięty, bo plik jest wspólny.

Nowe pliki: `powierzchnia-dokumentu.ts`, `zapis-formatowany.ts`,
`nastawy-strony.ts`, `nastawy-wizualne.ts`, `wstazka-pracy.ts`,
`wiersz-polecenia.ts`, `panel-redaktora.ts`, `ocena-redaktora.ts`,
`dymki-komentarzy.ts`, `galeria-szablonow.ts`, `suwaki-koncepcyjne.ts`,
`zmiany-modelu.ts`, `karty-dokumentow.ts`, `czynnosci-pracy.ts`,
`zrodlo-pracy-studio.ts`, `studio-praca.css`.

## Wykaz czynności przed scaleniem (41)

**Studio Editor (18):** wybór źródła dokumentu · pole wskazania · wczytanie
(`studio.document.open`) wraz ze zdaniem o skutku · edycja treści · zapamiętanie
zaznaczenia · pasek zaznaczenia: cztery operacje · „Więcej" do Tools Panelu ·
wskaźnik zakresu operacji · szesnaście narzędzi znacznikowych · wskaźnik formatu
dokumentu · wydanie PDF · wydanie DOCX · brak drogi „zmień format dokumentu" ·
pasek statusu (zapis, liczniki, wersja) · znajdź/zamień z podglądem trafień ·
zapis i założenie wersji · przyjęcie wyniku · odrzucenie wyniku · wstawienie
wyniku w miejsce kursora.

**Kanwa tekstowa (5):** wyrys wyniku operacji · opis pochodzenia wyniku ·
poprawka propozycji przed decyzją · pasek powiązania z oknem rozmowy · trzy
stany własne (bez operacji, propozycja bez treści, odmowa rdzenia).

**Preview Window (6):** wybór formatu wydania · wyrys treści zaakceptowanej ·
eksport · przekazanie do Library · akceptacja wyniku · dwa stany własne.

**Diff/Grep Panel (12):** wersja odniesienia · wersja porównywana · wzorzec ·
przełącznik wyrażenia regularnego · porównanie i wyszukanie · filtr różnic ·
statystyka różnicy · wyrys fragmentów · wyrys trafień · adnotacja · przyjęcie
zmiany · odrzucenie zmiany · przyjęcie pary porównania z Session Repository.

## Wykaz czynności po scaleniu (41 przeniesionych + 34 nowe)

Wszystkie 41 pozycji stoi w oknie pracy. Przeniesienia, które zmieniły miejsce:

- narzędzia znacznikowe → wstążka, zakładka „Narzędzia główne" i „Wstawianie";
- wskaźnik formatu dokumentu → pasek statusu (`pasek-statusu.ts`);
- wskaźnik zakresu operacji → pas decyzji okna (`ms-praca__zakres`), wraz
  z trzema zdaniami rozróżniającymi, w tym o wyborze ręcznym z Tools Panelu;
- cztery operacje paska zaznaczenia i „Więcej" → wiersz polecenia przy kursorze;
- poprawka propozycji z kanwy → poprawia się TREŚĆ dokumentu, bo wynik modelu
  stoi już w niej jako zmiana oznaczona; bufor obok dokumentu zniknął razem
  z kanwą i nie ma następnika, bo nie ma czego buforować;
- wyrys treści zaakceptowanej z Preview Window → tryb widoku „Podgląd wydania",
  który czyta `trescZaakceptowana()` tak samo jak czytał go podgląd;
- przyjęcie/odrzucenie propozycji → `studio.proposal.decide` w rdzeniu, a przy
  propozycji bez odwołania w rdzeniu — decyzja po stronie okna, jak dawniej.

Nowe czynności (34): widok strony z marginesami, nagłówkiem, stopką, numeracją
i podziałem na kartki · skala widoku · wiele kartek obok siebie · nośnik
i orientacja · cztery marginesy · profile wydania (odczyt i zapis nastaw strony
w rdzeniu) · render paginacji rdzenia · styl nazwany bloku · krój · stopień ·
interlinia · wcięcie · odstęp akapitowy · wyrównanie · barwa · tryb źródłowy ·
podział strony jawny · zakładki dwóch dokumentów · pięć suwaków koncepcyjnych ·
wykaz operacji kontekstowych na wstążce · wiersz polecenia własnymi słowami ·
komentarz z wiersza polecenia · dymki komentarzy z wątkami i rozwiązaniem ·
skoki po komentarzach · śledzenie zmian · przyjęcie i odrzucenie wszystkich
zmian · skoki po zmianach · trzy tryby adiustacji · pasek zatwierdzenia zmiany
modelu („Gotowe"/„Cofnij") · decyzja wybiórcza o fragmentach propozycji ·
różnica wyglądu · galeria szablonów z wyszukiwaniem i kategoriami · panel
Redaktora (ocena, korekty, uściślenia, podobieństwa, statystyki) · grupa ochrony
prowadząca do Warsztatu dokumentu.

## Zmiany w rdzeniu

1. **`studio.contextual.op` pisze do dokumentu.** `odlozWynikModeluStudia`
   (`adapter_modul_studio_roznice.go`) wpisuje wynik w miejsce zakresu operacji,
   rejestruje go jako zmianę śledzoną autora `model` (`ZapiszZmianeSledzona`
   miało dotąd zero wołających) i zakłada wersję przez `zalozWersjeDokumentu`
   z autorem `model` i odwołaniem do propozycji. Komenda rozgłasza od teraz
   `studio.document.changed`.
2. **`studio.proposal.decide` honoruje `hunkIndexes`** — `zlozTrescZFragmentow`
   składa treść z fragmentów `policzFragmentyRoznicy`: wskazany bierze stronę
   propozycji, pominięty zostaje z dokumentu. Numer spoza rachunku to odmowa.
3. **Zakres zmiany śledzonej liczy się w ZNAKACH, nie w bajtach.**
   `zastosujZmianeSledzona` cięło treść po bajtach, więc decyzja o zmianie
   w tekście z polskimi literami wstawiała treść w środek znaku. Kontrakt mówi
   „w znakach" — naprawione.
4. **`params` operacji dojeżdża do modelu** — `trescOperacjiStudia` dokłada je do
   polecenia. Bez tego suwaki i wiersz polecenia przestawiałyby pole, którego
   model nie czyta.

## Potrzeby kontraktowe (kontraktu NIE ruszano)

1. **Styl nazwany i nastawy wizualne w dokumencie.** `StudioDocument` ma treść,
   tytuł, format i wersję; `studio.document.save` przyjmuje `documentId`,
   `content`, `title`, `createVersion`. Krój, stopień, interlinia, wcięcie,
   odstęp akapitowy, wyrównanie i barwa żyją więc przez sesję okna i giną z nią.
   Trwałe jest to, co jest składnią treści. Potrzeba: pole stylów dokumentu albo
   akapitu w `StudioDocument` i w żądaniu zapisu.
2. **Nastawy strony w dokumencie.** `StudioPageSetup` istnieje, ale wyłącznie
   przy PROFILU WYDANIA. Dokument nie ma własnych nastaw strony, więc kartka
   Operatora zapamiętuje się tylko przez nazwany profil.
3. **Wielkość ciągła operacji.** Nastawy jadą `params` (`json.RawMessage`) i to
   działa, ale kontrakt nie nazywa ani jednej wielkości, więc rdzeń nie może ich
   sprawdzić. Potrzeba: nazwane pola natężenia w żądaniu operacji.
4. **Usunięcie komentarza.** `studio.comment.*` niesie dodanie, wykaz
   i rozwiązanie; usunięcia nie ma. W oknie stoi jako brak nazwany.
5. **Pobranie bajtów zasobu do przeglądarki.** Wydanie, render podglądu
   i nakładki różnicy wyglądu oddają identyfikatory zasobów; komendy
   pobierającej ich treść nie ma, więc podgląd wydania rysuje kartki po stronie
   klienta, a liczby rdzenia stoją obok jego liczb.
6. **Znacznik „przejrzane" przy uściśleniu panelu Redaktora** — odhaczenie żyje
   przez sesję, bo kontrakt nie ma pola na taki stan.

## Czynności recenzji bez drogi w module Studio (brak nazwany, nie udawany)

- **tezaurus** — wymaga słownika języka; arsenał wkompilowany go nie niesie,
  a program spoza instalki jest zakazany;
- **czytanie na głos** — kanału mowy w module Studio kontrakt nie niesie;
- **sprawdzenie ułatwień dostępu dokumentu** — rdzeń nie ma takiego pomiaru;
- **tłumaczenie w miejscu** — należy do modułu Translate; Studio ma tylko
  operację kontekstową „tłumaczenie zaznaczenia" zlecaną modelowi;
- **pisownia i gramatyka w panelu Redaktora** — bez słownika i analizy
  składniowej nie ma pomiaru; wiersze stoją z wartością niepodaną i powodem.

## Kompilacja

`go build ./server/...` zielone. Sprawdzian typów klienta zielony —
uruchamiany jako `cd client && npx tsc --noEmit -p tsconfig.json`. Postać
`npx --prefix client tsc --noEmit` ze zlecenia NIE sprawdza projektu: wypisuje
pomoc kompilatora i wychodzi kodem 1, bo `--prefix` nie przestawia katalogu
roboczego, a bez `-p` kompilator nie widzi `tsconfig.json`. Warto poprawić
polecenie w zleceniach, żeby zielone nie było zielonym pozornym.

## Sprawdziany

- `server/internal/core/skutek_pracy_z_dokumentem_test.go` — pięć sprawdzianów
  skutku: przyjęcie wskazanych fragmentów zmienia wyłącznie je, fragment spoza
  rachunku odmawia, zakres zmiany liczy się w znakach (tekst z polskimi
  literami), zakres operacji schodzi na całość zamiast odmawiać, polecenie
  i nastawy dojeżdżają do modelu.
- `client/src/moduly/studio/praca-z-dokumentem.test.ts` — czternaście
  sprawdzianów, w tym **wykaz czynności pozycja po pozycji na zbudowanym
  widoku** (żadna z 41 nie zniknęła), odwracalność widoku formatowanego,
  paginacja z podziałem jawnym, oznaczenie zmiany modelu w miejscu, pomiary
  panelu Redaktora wraz z brakiem pomiaru pisowni, suwaki.

## Co zostało otwarte

- **Sprawdzian operacji kontekstowej od końca do końca.** Uprząż `core` nie
  stawia kanału modelu, więc droga „operacja → zmiana śledzona w dokumencie"
  jest zmierzona po częściach (rachunek zakresu, złożenie polecenia, decyzja
  o zmianie), a nie jednym przebiegiem.
- **Dwa dokumenty naraz to zakładki, nie dzielona powierzchnia.** Zlecenie
  dopuszczało jedno albo drugie; kartka A4 obok drugiej A4 zeszłaby do rozmiaru
  nieczytelnego.
- **Podgląd wydania rysuje kartki po stronie klienta.** Strony rdzenia
  (`studio.preview.render`) są liczone i ich liczba jest pokazana obok liczby
  klienta, ale obrazu stron okno nie pobiera — nie ma czym.
- **Trzy sprawdziany niepomyślne, wszystkie w cudzej robocie — w module Design.**
  `client/src/moduly/design/zrodlo-designu.test.ts` (komendy `design.photo.*` bez
  okna), `TestRejestrPokrywaKomendyKontraktu` (te same komendy bez uchwytu
  w rejestrze rdzenia) oraz `TestWarsztatDokumentuNieWymieniaZdjetychProgramow`
  — ten ostatni wskazuje `adapter_modul_design_druk_wydania.go`, który wymienia
  **ghostscript**, program zdjęty z rdzenia, bo instalka go nie niesie. Rodzina
  `design.photo.*` i ten plik weszły przy równoległej pracy innego wykonawcy;
  moduł Design ani kontrakt nie były tu ruszane, ale zaporę programów
  zewnętrznych warto mu zgłosić — łamie zasadę cienkiej instalki.
  Sprawdziany Studia: `go test` rodziny Studia oraz `vitest src/moduly/studio` —
  47 pomyślnych, zero niepomyślnych.
- **`go mod tidy` nie był potrzebny** — ani jednej nowej biblioteki.

---

# Design — domknięcie całego modułu (wykonawca Designu, 17.08.2026)

Zakres tej sekcji: WYŁĄCZNIE moduł Design. Plików Studia nie ruszano.
`shared/contract.json` był tu wnoszony **przed** ustaleniem jednego pisarza
kontraktu — od tej chwili nie jest już przeze mnie dotykany (zapis niżej,
w „Rozstrzygnięcia obowiązujące następnego wykonawcę").

## Krok 0 — drzewo się nie kompilowało

`adapter_modul_design_szablony_materialu.go:290` wołało `czyBrakWierszaDesignu`,
którego w drzewie nie ma. Jedna zmiana: → `czyBrakZasobuDesignu`
(`adapter_modul_design_etykiety.go:65`). Po niej `go build ./...` przechodzi.

## Co zbudowane, liczbami

| Rzecz | Liczba |
|---|---:|
| Komendy `design.*` w kontrakcie | 104 |
| Komendy `design.*` bez obsługi w rdzeniu | **0** |
| Komend dobudowanych w tej turze | **75** (55 z kontraktu zastanego + 20 nowych `design.photo.*`) |
| Metod dołożonych do portu `Design` | 75 + 2 odczyty dla szyny zdarzeń |
| Nowych plików rdzenia | 16 |
| Nowych plików warstwy danych | 1 (`dane/design_fotografia.go`) + strony szablonu w istniejącym |
| Nowych migracji | 4 — **339** (strony szablonu), **340** (okna warsztatów), **346** (łańcuch edycji), **347** (nastawy fotografii) |
| Nowych okien klienta | 6 warsztatów |
| Nowych plików klienta | 4 |
| Sprawdziany skutku | 8 w `skutek_fotografii_designu_test.go` + 2 zapory |

### Pliki rdzenia i czym liczą

| Plik | Komendy | Czym liczone |
|---|---:|---|
| `adapter_modul_design_wektor.go` + `_wektor_sciezki.go` | 10 — `design.vector.*` | `tdewolff/canvas` (operacje logiczne, PDF, EPS), rachunek własny na węzły i uchwyty, sklejenie dokumentu SVG |
| `adapter_modul_design_makiety.go` + `_makiety_zrzut.go` | 15 — `design.frame.*`, `design.component.*`, `design.prototype.*`, `design.mockup.*`, `design.layout.auto`, `design.constraint.set`, `design.grid.set` | rachunek własny na `dane/design_makiety.go`; segmentacja zrzutu — próg odstępstwa od tła + spójne bloki na kratce 8 px |
| `adapter_modul_design_kolor.go` + `_kolor_obraz.go` | 7 — `design.color.*` | gotowe `_barwy.go`, `lucasb-eyer/go-colorful` (HCL, Lab), k‑średnich w Lab (paleta z obrazu), macierze LMS Brettela–Viénota–Mollona (wady widzenia) |
| `adapter_modul_design_ikony.go` + `_ikony_katalog.go` + `_ikony_kroj.go` | 8 — `design.icon.*`, `design.favicon.build`, `design.font.*` | katalog **46 ikon wkompilowanych** w kod (siatka 24, obrys 2), `x/image/font/sfnt` (kontury glifów), `x/image/font/gofont` (8 krojów wkompilowanych), **własny zapis kroju TrueType** dla postaci `webfont` |
| `adapter_modul_design_marketing.go` + `_bazy_zdjeciowe.go` | 4 — `design.campaign.set.build`, `design.stock.*`, `design.product.mockup.render` | `x/image/draw` (skalowanie, homografia z obrotem i kryciem), 8 dostawców baz zdjęciowych po HTTP |
| `adapter_modul_design_druk.go` + `_druk_wydania.go` | 8 — `design.print.*`, `design.largeformat.tile` | `pdfcpu` (PDF, wielostronicowy), `x/image/tiff` (Deflate), **własny zapis EPS** z obrazem osadzonym operatorem `colorimage`, `x/image/draw` (kafle) |
| `adapter_modul_design_wykresy.go` | 2 — `design.chart.render`, `design.diagram.render` | `tdewolff/canvas` — jedno płótno, dwa wydania (SVG wydawcą `renderers/svg`, PNG rasteryzatorem `renderers/rasterizer`) |
| `adapter_modul_design_fotografia.go` + `_fotografia_rachunek.go` + `_fotografia_maski.go` + `_fotografia_wsad.go` | 20 — `design.photo.*` | wykaz niżej |
| `adapter_modul_design_szablony_materialu.go` | 3 — `design.template.*` (były napisane, niewpięte) + strony publikacji | `pdfcpu` |
| `adapter_modul_design_kroje.go` | — (warsztat wspólny) | `x/image/font/sfnt`, przegląd katalogów krojów serwera |

## Czym liczona fotografia

Wszystko **w procesie, czystym Go, wkompilowane w binarium**. Ani jednego
`exec.Command`, ani jednego `zewnetrzne.Wolaj` — pilnuje tego zapora
`server/internal/core/zapora_fotografii_test.go`.

- `disintegration/imaging` — Lanczos / dwuliniowy / najbliższy sąsiad, kadr,
  obrót, rozmycie, wyostrzenie, jasność, kontrast, gamma, nasycenie, krzywa
  sigmoidalna, skala szarości.
- `golang.org/x/image/draw` — przekształcenia afiniczne, kompozycja warstw.
- `golang.org/x/image` (webp, tiff, bmp) — odczyt formatów spoza biblioteki
  standardowej.
- `lucasb-eyer/go-colorful` — HSL i HCL przy przesunięciu odcienia i jasności
  percepcyjnej.
- **Rachunek własny**, bez biblioteki: homografia z czterech par punktów
  (eliminacja Gaussa z wyborem elementu głównego), dystorsja obiektywu
  (wielomian promieniowy, odwzorowanie odwrotne + próbkowanie dwuliniowe),
  auto‑poziomy z odcięciem 0,5 % z każdej strony histogramu, krzywe tonalne przez
  tablice przeglądowe 256‑elementowe, odcięcie tła (barwa tła z obwodu + rozrost
  obszaru od brzegów + wygładzenie samego kanału krycia), zaznaczenie obiektu
  (rozrost od punktu), domalowanie i leczenie (rozrost średniej od brzegu obszaru
  w głąb), rozszerzenie kadru (odbicie lustrzane brzegu), ziarno deterministyczne
  z położenia punktu, winieta kwadratowa od połowy promienia, tryby mieszania
  (multiply/screen/overlay), obrysowanie konturów (k‑średnich w Lab → mapa
  etykiet → spójne obszary → odcinki poziome scalane w pionie), odczyt EXIF
  (segment APP1, katalog IFD0, 15 rozpoznawanych pól).

### Cztery czynności o wariancie neuronowym — którą drogą naprawdę idą

`design.photo.upscale`, `design.photo.background.remove`,
`design.photo.inpaint`, `design.photo.expand`.

**Dziś idą RACHUNKIEM RDZENIA i odpowiedź mówi to wprost:** pole `computedBy`
niesie wartość `rachunekRdzenia` (wyliczenie `DesignPhotoComputeRoute`), a ta
sama wartość zapisuje się w łańcuchu edycji
(`czynnosc_fotografii_design.policzone_przez`), więc po tygodniu da się
powiedzieć, którą drogą powstał konkretny wariant.

Powód, dla którego droga neuronowa **jeszcze nie stoi**, jest twardy i zmierzony:
kanał obrazowy rdzenia (`models.AdapterObrazy`) przyjmuje **wyłącznie polecenie
tekstowe**. Ani `models.Zapytanie`, ani adapter kanału nie mają pola na obraz
**wejściowy**. Bez niego kanał wygenerowałby obraz **nowy**, a nie przetworzył
zdjęcia Operatora — i byłoby to podstawienie cudzej treści pod jego materiał.
Wskazany `channelId` jest natomiast **sprawdzany** tą samą drogą, co przy
`design.asset.generate` (kanał nieistniejący, nieczynny albo tekstowy → odmowa
z nazwaniem naprawy), żeby wskazanie nieistniejącego kanału nie przeszło w ciszy.

Gdy pole obrazu wejściowego dojdzie do warstwy modeli, zmienia się **jedna gałąź
w czterech metodach** — reszta drogi (wariant, łańcuch, `computedBy`) już stoi.

## Rozstrzygnięcia obowiązujące następnego wykonawcę

1. **`shared/contract.json` ma jednego pisarza i nie jest nim wykonawca Designu.**
   Dobudowa Designu weszła dwoma skryptami **idempotentnymi**, zostawiającymi
   kopię sprzed zmiany: `scripts/wnies-fotografie-do-kontraktu.mjs` (20 komend,
   5 wyliczeń, 7 struktur) i `scripts/wnies-publikacje-do-kontraktu.mjs`
   (`DesignTemplatePage`, `DesignPrintBinding`, pola `pages`, `pageNumber`,
   `templateId`, `pageOrder`, `binding`, `pageCount`). Kopie leżą jako
   `shared/contract.json.przed-fotografia` i `…przed-publikacjami`. Skrypty wolno
   puścić powtórnie — pomijają nazwy już obecne.
2. **Jednostką kompozycji przy druku jest MILIMETR.** Rozstrzygnięcie
   `adapter_modul_design_druk.go` i obowiązuje spójnie: kontrolę przeddrukową,
   wydanie i podział na kafle. Dzięki temu kompozycja 210×297 jest arkuszem A4.
   Wyrys ekranowy (`design.board.export`) liczy te same liczby jako **piksele** —
   tam nie ma nośnika, więc nie ma czego mierzyć w milimetrach.
3. **`design.print.export` ODMAWIA przy wadzie o wadze błędu.** Także przy
   publikacji wielostronicowej — wada na jednej stronie odmawia całego wydania,
   a numer strony wchodzi w zastrzeżenie. Pominięcie kontroli wraca polem
   `preflightSkipped`.
4. **Publikacje to szablon materiału o wielu stronach, nie rodzina komend.**
   Rodziny `design.publication.*` nie ma i nie ma jej zakładać. Szablon **bez**
   stron zostaje jednostronicowy i czyta warstwy tam, gdzie czytał (migracja 336)
   — baner nie ma stron.
5. **Obróbka zdjęcia zakłada WARIANT, oryginał zostaje.** `variantOfAssetId`
   wskazuje źródło, a obok wariantu powstaje ogniwo łańcucha (migracja 346)
   z nastawami. Wynik fotografii wychodzi **jako PNG** — bezstratnie i z kanałem
   krycia; dziesięć korekcji w JPEG to dziesięć kompresji.
6. **Kształt powstaje od razu jako węzły ścieżki.** Nie ma bytu „prostokąt" do
   późniejszej zamiany. Uchwyty węzła (`handleIn*`, `handleOut*`) są
   **odsunięciami** od węzła, nie współrzędnymi bezwzględnymi.
7. **Zamiana tekstu w kontury jest nieodwracalna dla wyniku**, więc tekst
   źródłowy zostaje w nazwie ścieżki (`sciezka_wektorowa_design.nazwa`).
8. **Bilans zamiast ciszy** — wypełnione są wszystkie dziesięć pól kontraktu
   przewidzianych na to (`unreachableFrameIds`, `failedConcepts`, `gridWarnings`,
   `unmatchedNames`, `failedSizes`, `providersFailed`, `unplacedNodeIds`,
   `unrecognizedRegions`, `removedPathIds`, `iccProfiles`) plus dołożone:
   `regionsSkipped`, `failedAssetIds`, `skippedLayerAssetIds`, `droppedRegions`,
   `appliedSteps`.
9. **Katalog ikon i kroje są WKOMPILOWANE.** 46 ikon w
   `adapter_modul_design_ikony_katalog.go` (siatka 24, obrys 2, zestaw
   `rdzen-24`) i 8 krojów rodziny Go. Wzór rdzenia **nie daje się nadpisać** —
   inaczej dwie instalacje miałyby dwa różne katalogi pod tą samą nazwą.
10. **Klucze baz zdjęciowych czyta się z sejfu rdzenia** pod bytem
    `design.stock.<dostawca>`. Moduł Design **wyłącznie czyta** — nigdy nie
    zapisuje. Brak klucza to `providersFailed`, nie odmowa całej komendy.
11. **Zasób z bazy zewnętrznej bez zapisanej licencji jest usterką.** Gdy zapis
    licencji padnie, wiersz zasobu jest **wycofywany**, żeby Assets Panel nie
    pokazał materiału bez prowenancji.
12. **Warsztaty klienta opisane są DANYMI.**
    `client/src/moduly/design/czynnosci-warsztatow-designu.ts` — 75 czynności
    w 6 grupach; okno buduje pola z katalogu
    (`okno-warsztatu-designu.ts`), a droga do rdzenia jest jedna
    (`zrodlo-warsztatow-designu.ts`). Kontrolki pól bierze się z warsztatu
    dokumentu Studia (`utworzKontrolke`) — **jeden przełącznik rodzajów pól na
    dwa moduły**, więc nowy rodzaj pola dokłada się tam raz.
13. **Zapora `zapora_fotografii_test.go` — nie osłabiać.** Pilnuje dwóch rzeczy:
    żaden plik `adapter_modul_design*.go` nie uruchamia procesu, i żaden nie
    wymienia nazw silników obrazu spoza instalki — także w komentarzu. Zakres jest
    obszarem Design, nie całym rdzeniem, i to jest **pomiar granicy
    odpowiedzialności**: te nazwy stoją dziś w nagłówkach modułu Library
    i narzędzi obrazowych modelu jako zapis decyzji ich autorów, a przepisanie
    cudzych plików jest osobną pracą i osobną zgodą.
14. **Nazwy pomocnicze noszą przedrostek obszaru** (`…Designu`, `…Design`) —
    także w `_test.go`. Jedna kolizja wyszła i została nazwana:
    `obrazZMagazynu` już istniało w `skutek_obrazu_wektor_test.go`, więc mój
    nazywa się `obrazZMagazynuFotografii`.

## Zakresy migracji zajęte w tej turze

- **339** — `strona_szablonu_materialu_design`, `warstwa_strony_szablonu_design`.
- **340** — sześć okien warsztatów Designu w katalogu okien.
- **346** — `czynnosc_fotografii_design` (łańcuch edycji).
- **347** — `nastawa_fotografii_design`.

Wolne w zakresie Designu zostają: **320–325, 327–335, 341–345, 348–360**.

## Defekty znalezione i naprawione przez sprawdziany skutku

Oba wyszły z pomiaru, nie z przeglądu kodu:

1. **`transparentShare` mówiło o obrazie pośrednim.** Udział punktów
   przezroczystych liczył się **przed** wygładzeniem krawędzi maski, a wygładzenie
   krycie zmienia — plik i odpowiedź rozjeżdżały się o 5 punktów procentowych.
   Pomiar przeniesiony za wygładzenie.
2. **Warstwy stron publikacji nie dały się odczytać.** Zapytanie wybierało 10
   kolumn, a wspólny odczytywacz warstwy żąda 12 (`sql: expected 10 destination
   arguments in Scan, not 12`). Zapytanie uzupełnione, bez drugiego odczytywacza.

## Sprawdziany

- `server/internal/core/skutek_fotografii_designu_test.go` — **8 sprawdzianów
  skutku**, każdy schodzi do magazynu i mierzy PLIK, nie kopertę: rozdzielczość
  po `upscale` (80×60 z 20×15), obecność **i treść** kanału krycia po
  `background.remove` (naroże 0, środek ≠ 0, udział przeliczony niezależnie),
  przesunięcie średniej jasności po `color.correct` (plus sprawdzenie, że
  **oryginał został nietknięty**), liczba i wymiary kafli po `largeformat.tile`
  (3×2 = 6, każdy z bajtami), **liczba stron w pliku PDF** po wydaniu publikacji
  (4 — liczone własnym rozbiorem, z rozpakowaniem strumieni obiektów przez
  `compress/zlib`, żeby nie mierzyć biblioteki jej samą sobą), odmowa wydania przy
  wadzie o wadze błędu wraz z powodzeniem przy jawnym pominięciu, kolejność
  i nastawy ogniw łańcucha edycji, metadane mierzone z pliku.
- `server/internal/core/zapora_fotografii_test.go` — 2 zapory (opisane wyżej).
- `client/src/moduly/design/zrodlo-designu.test.ts` — istniejący sprawdzian
  pokrycia **nauczony drugiej drogi** (warsztaty), bez osłabienia wymagania: nadal
  żąda drogi z okna do **każdej** komendy `design.*`.

Wyniki zmierzone: `go build ./...` zielony, `go vet ./server/internal/core/`
czysty, `go test` rodziny Designu pomyślny, `go test ./server/internal/store/`
pomyślny, `go test ./server/internal/dane/` pomyślny, `npx tsc --noEmit` czysty,
`npx vitest run` w kliencie **177/177**.
`grep -cE "pdfcpu|go-git/go-git|gotextdiff|docker/docker|x/image|kin-openapi|pgx|go-sql-driver" go.mod` → **8**.

## Co zostało OTWARTE — bez zaokrąglania w górę

1. **Droga neuronowa czterech czynności nie stoi.** Powód wyżej: warstwa modeli
   nie ma pola na obraz wejściowy. `computedBy` mówi o tym prawdę, ale funkcja
   jest dziś słabsza od zapowiedzi zlecenia.
2. **`design.mockup.import` nie CZYTA liter.** Rdzeń **wykrywa układ** obszarów
   i rozpoznaje, które są liniami tekstu (rzędy drobnych, regularnych bloków),
   i tak je znaczy. Odczytu treści napisów nie ma — wymaga biblioteki OCR, której
   w drzewie nie ma. **Zgłaszam brak, nie obszedłem go.** Pole `recognizeText`
   włącza to, co rdzeń naprawdę potrafi, a nagłówek pliku mówi to wprost.
3. **Smithsonian wymaga klucza, wbrew założeniu zlecenia.** Zlecenie stawiało go
   w grupie „bez klucza"; jego API (`api.si.edu/openaccess`) wymaga klucza
   `api.data.gov` — darmowego i natychmiastowego, ale wymaganego. Wszedł do grupy
   z kluczem i bez klucza wraca w `providersFailed`. Bez klucza pracują cztery:
   Openverse, Wikimedia Commons, Met Museum, NASA.
4. **Dostawcy baz zdjęciowych nie mają sprawdzianu od końca do końca.** Osiem
   dróg HTTP jest napisanych i zbudowanych, ale nie zmierzonych przelotem —
   uprząż `core` nie stawia serwera próbnego dla tych ośmiu kształtów odpowiedzi.
   Kształty pochodzą z dokumentacji dostawców, nie z pomiaru.
5. **Rdzeń nie rozkłada profili ICC i nie przelicza barw przez nie.** RGB → CMYK
   idzie przeliczeniem **wprost** i kontrola przeddrukowa mówi to jako
   zastrzeżenie rangi „informacja" przy każdym wydaniu w CMYK. Wykaz
   `iccProfiles` jest pomiarem katalogów tej maszyny — nie zapowiedzią.
6. **Zgodność z normą PDF/X nie jest weryfikowana.** Profil może jej żądać, rdzeń
   składa dokument biblioteką wkompilowaną i **mówi wprost**, że normy nie
   sprawdza (zastrzeżenie `norma-bez-weryfikacji`). Weryfikacja wymaga narzędzia
   spoza instalki.
7. **TIFF i EPS wychodzą jednostronicowo.** Wielostronicowy TIFF istnieje jako
   format, ale koder biblioteki go nie zapisuje; EPS jest jednostronicowy
   z natury. Publikacja wielostronicowa w tych formatach **odmawia** z nazwaniem
   powodu i wskazaniem PDF.
8. **`design.photo.vectorize` obrysowuje kontury prostokątnie.** Kontur idzie
   krawędziami punktów — schodkowy, ale **prawdziwy**: opisuje dokładnie te
   punkty, które do obszaru należą. Pole `smoothing` jest dziś przyjmowane
   i sprawdzane, ale rachunek go nie używa (scalanie odcinków w pionie jest
   bezwarunkowe). To brak, nie cisza — ale pole obiecuje więcej, niż daje.
9. **`design.photo.batch.apply` puszcza 8 z 20 czynności.** Wsad przyjmuje
   wyłącznie czynności działające na jednym obrazie i niepotrzebujące wskazania
   drugiego zasobu; kompozycja warstw, maska i zaznaczenie obiektu wymagają
   wskazań, których wsad nie ma skąd wziąć dla każdego zasobu osobno. Odmowa
   wymienia wykaz przyjmowanych.
10. **Krój ikonowy (`webfont`) zapisuje kontury spłaszczone do odcinków.** Glif
    niesie wyłącznie punkty na konturze, bez krzywych kwadratowych. Tolerancja
    spłaszczenia to jedna jednostka kroju (1/1024 firetu), więc różnicy nie widać
    — ale plik jest większy, niż byłby z krzywymi. **Krój nie został sprawdzony
    w przeglądarce ani w systemie** — nie ma czym na tej maszynie.
11. **Odszumienie jest rozmyciem Gaussa, nie filtrem bilateralnym.** Zachowuje
    krawędzie tylko o tyle, o ile pozwala rozmycie. Filtru bilateralnego
    biblioteka nie ma i rdzeń tego nie udaje.
12. **`design.mockup.generate` liczy układ regułą rdzenia.** Kanał modelu, gdy
    wskazany, rozpisuje **opis** na nazwy sekcji; samego układu nie liczy nigdy.
    Niepowodzenie kanału nie kończy komendy — liczy wtedy reguła.
13. **Sprawdziany warsztatu wektora, ikon, makiety i barwy to sprawdziany
    kompilacji, nie skutku.** Zmierzone przelotem są fotografia i druk. Wektor
    (operacje logiczne na węzłach), krój ikonowy (liczba glifów w pliku TTF)
    i schemat (liczba ścieżek w SVG) **zasługują na pomiar i go nie mają**.

## Potrzeby, których nie wolno mi było zrobić samemu

1. **Wykaz nośników druku powinien stanąć w miejscu wspólnym ze Studiem, wraz
   z uzupełnieniem o koperty C4, C5, C6.** Dziś stoi w
   `server/internal/core/adapter_modul_design_druk_wspolne.go` jako
   `nosnikiDruku` (17 pozycji: ISO A0–A6, ISO B1–B3, DL, wizytówka, Letter,
   Legal, Tabloid, dwie rolki wielkoformatowe) i jest **wiedzą rdzenia**, którą
   czytają trzy czynności Designu. Studio ma własne nastawy strony
   (`client/src/moduly/studio/nastawy-strony.ts`) — dwa wykazy rozmiarów
   w jednym produkcie rozjadą się przy pierwszej poprawce i Operator dostanie
   A3 o wymiarach A4. **Przeniesienia nie wykonałem, bo wymagałoby wejścia
   w pliki Studia, w których pracuje trzech wykonawców.** Brakujące koperty
   C4 (229×324 mm), C5 (162×229 mm) i C6 (114×162 mm) dokłada się wtedy raz,
   w miejscu wspólnym.
2. **Runtime ONNX z wagami modeli na serwerze** — poza zakresem wykonawcy,
   decyzja Właściciela tej samej wagi co silnik kontenerów. Bez niego
   punkt 1 sekcji „otwarte" zostaje otwarty.
3. **Zapora programów zewnętrznych dla całego rdzenia.** Moja pilnuje obszaru
   Design. Nazwy silników obrazu stoją dziś w nagłówkach modułu Library
   (`adapter_modul_library_odciski.go`) i narzędzi obrazowych modelu
   (`adapter_narzedzia_obraz*.go`, `kompozycja.go`, `handlers_narzedzia_obraz.go`).
   Rozszerzenie zapory na rdzeń wymaga przepisania tych nagłówków — cudzej pracy
   i osobnej zgody.
4. **Pole obrazu WEJŚCIOWEGO w `models.Zapytanie` i w adapterze kanału
   obrazowego.** Bez niego cztery czynności fotografii nie mają drogi
   neuronowej. Zmiana należy do warstwy modeli, nie do modułu Design.
5. **Biblioteka OCR** dla `design.mockup.import` (odczyt treści napisów) —
   dociągnięcie zależności tej wagi jest decyzją, której sam nie podejmuję.

# Studio, odcinek 3 — znakowanie, asystent, schowek, osadzenie (17.08.2026)

Zlecenie: `ZLECENIE-STUDIO-MODEL-DOKUMENTU.md` wraz z dwoma uzupełnieniami
dopisanymi w trakcie tej tury („operacje jako narzędzia ukryte" i „silnik mowy
idzie z pakietem serwera"). Podział: `PODZIAL-PRAC-STUDIO.md`, odcinek 3.

## Pliki odcinka

Nowe: `przybornik-znakowania.ts`, `przybornik-znakowania.css`,
`przybornik-wykaz.ts`, `przybornik-znaczniki.ts`, `przybornik-katalog.ts`,
`przybornik-plywak.ts`, `przybornik-uzycie.ts`, `przybornik-mowa.ts`,
`przybornik-zrodlo.ts`, `przybornik-zaplecze.ts`,
`przybornik-znakowania.test.ts`, `schowek-zrodlo.ts`, `schowek-historia.ts`,
`osadzenie-zrodel.ts`, `osadzenie-panel.ts`, `osadzenie-pochodzenie.ts`.

Przebudowane własne: `dymki-komentarzy.ts` (trzy byty marginesu),
`wiersz-polecenia.ts` (pływak, mikrofon, podpowiedzi), `kategorie-operacji.ts`
(wyszukiwanie, jeden wykaz), `okno-tools-panel.ts` (tryb do wyboru).

Tknięte pliki wspólne, bo bez nich odcinek nie byłby wpięty:
`okno-pracy-z-dokumentem.ts` (osadzenie paneli, czynności) i `modul-studio.ts`
(jedno zaplecze rodzin komend, nastawa trybu wspólna dwóm oknom). Oba nie
należały w podziale do żadnego odcinka; `wstazka-pracy.ts`,
`powierzchnia-dokumentu.ts`, `nastawy-strony.ts`, `karty-dokumentow.ts`
i `studio-praca.css` **nie były ruszane**.

## Co zbudowane

- **Przybornik znakowania przy krawędzi treści** — znakowanie fragmentu
  (komentarz, propozycja brzmienia, adnotacja), znacznik własny z nazwą i barwą,
  zakładka powrotu, wykaz znakowań z zawężeniem po rodzaju, autorze i stanie,
  skakanie po wykazie, legenda trzech bytów marginesu, jawna część „bez zaplecza
  w rdzeniu" z sześcioma brakami nazwanymi.
- **Trzy byty marginesu rozróżnione** — komentarz (obrys subtelny), propozycja
  brzmienia (obrys przerywany, „w treści jeszcze nie ma"), zmiana śledzona
  (obrys pełny, „jest już w treści"), każda z własną decyzją i własną komendą.
- **Narzędzia ukryte zamiast stałej kolumny** — pływak przy zaznaczeniu
  z czynnościami najczęstszymi liczonymi z użycia, uchwyt pełnego katalogu
  (grupy, opisy, szukanie, wędrówka strzałkami z `komponenty/menu-drzewo.ts`),
  suwaki wielkości ciągłych jako sterowanie ciągłe, przypinanie czynności.
  Pływak staje nad zaznaczeniem, a gdy nad nim brak miejsca — pod nim, wedle
  zmierzonej wysokości.
- **Stały panel operacji jako tryb do wyboru** — zwinięty do jednego wiersza,
  z pasa wiodącego schodzi selektorem stanu, nastawa trwała przez `config.set`
  na kluczu `studio_przybornik_tryb_operacji`, wspólna dla pływaka i panelu.
- **Operacje własne Operatora** w tym samym wykazie co fabryczne
  (`studio.operation.list`, `.save`, `.delete`), rozdzielone polem `builtin`;
  usunięcia fabrycznej klient nie uprzedza — odmowę oddaje rdzeń.
- **Mikrofon** w wierszu polecenia i **dyktafon do treści** — jeden rachunek
  (`speech.availability.get`, `speech.audio.upload`, `speech.transcribe`), trzy
  stany odpowiedzi rozróżnione, nagranie nie opuszcza maszyny rdzenia.
- **Schowek z historią** — wykaz do wyboru, przypinanie, dwie równorzędne drogi
  wklejenia (z postacią i jako czysty tekst), malarz formatów, odłożenie
  zaznaczenia, czyszczenie historii nieprzypiętej.
- **Osadzenie przeglądarki i Biblioteki** — panel obok treści albo na całej
  powierzchni, wniesienie wprost w miejsce kursora, wiersz pochodzenia w treści
  i wykaz pochodzeń sesji.
- **Adnotacja przestawiona z drogi generycznej na własną komendę** —
  `studio.annotation.add` i `.list` zamiast `window.action` z identyfikatorem
  `studio.diff.adnotacja`, którego katalog akcji nie ma.

## Potrzeby kontraktowe (kontraktu NIE ruszano)

1. **Znacznik własny fragmentu** — nazwa, barwa, autor, zakres, stan. Dziś żyje
   przez sesję okna; model go nie założy, bo komendy nie ma.
2. **Wyróżnienie fragmentu barwą** jako cecha postaci dokumentu. Okno oznacza
   dziś barwę PISMA akapitu i mówi wprost, że to nie dojeżdża do rdzenia.
3. **Zakładka w miejscu, przypis, odsyłacz, odwołanie wzajemne, wstawienie
   tabeli i pola** — aparat dokumentu i postać dokumentu.
4. **Pochodzenie fragmentu dokumentu** — skąd wniesiony (plik, wersja, adres).
   Dziś wiersz w treści plus wykaz sesji.
5. **Autor wpisu schowka** — `clipboard.*` nie rozróżnia wpisu modelu od wpisu
   Operatora; okno zapisuje `sourceWindowId`, bo tylko to pole niesie.
6. **Odsłuch treści („czytaj na głos") w panelu redaktora** — nie budowane,
   bo nie było zamówione. Droga istnieje: `speech.availability.get` niesie
   `synthesisAvailable`, a synteza stoi w `roundtable.speech.synthesize`
   i `translate.speech.synthesize`. Do rozstrzygnięcia przez Właściciela, czy
   Studio ma po nią sięgać, czy dostać własną komendę.

## Pakiet serwera

Silnik mowy jest wedle rozstrzygnięcia składnikiem pakietu serwera, ale
`scripts/instalka-windows.sh` i katalog `wydania` **nie wymieniają ani `piper`,
ani `espeak-ng`, ani `whisper`, ani Tesseractu, ani 7-Zipa**. Wykaz zależności
pakietu wymaga uzupełnienia — inaczej odmowa `speech.availability.get` będzie
prawdziwa u każdego Operatora.

## Kompilacja i sprawdziany

- `(cd client && npx tsc --noEmit -p tsconfig.json)` — czysto.
- `npx vitest run src/moduly/studio` — 19 nowych sprawdzianów odcinka
  pomyślnych; jeden sprawdzian cudzy (`praca-z-dokumentem.test.ts`, wykaz
  czynności) niepomyślny z powodu zmiany w `wstazka-pracy.ts` odcinka 2:
  przełącznik trybu źródłowego przestał nieść `data-czynnosc='tryb-zrodlowy'`.
- `go build ./...` nieuruchamiane — odcinek nie tknął ani jednego pliku Go.

---

# Odcinek 2 — Klient: powierzchnia, wstążka, widok

Zakres wykonany w `client/src/moduly/studio/`. Ani jednego pliku Go, ani jednej
zmiany w `shared/contract.json`, ani jednego wejścia w `server/internal/`
i w `client/src/moduly/design/`.

## Pliki własne odcinka

Zmienione: `nastawy-strony.ts`, `powierzchnia-dokumentu.ts`, `wstazka-pracy.ts`,
`karty-dokumentow.ts`, `studio-praca.css`, `okno-session-repository.ts`,
`wiersz-wersji.ts`, `filtr-historii.ts`, `okno-warsztatu-dokumentu.ts`,
`braki-cyfryzacji.ts`, `okno-ingest-ocr-panel.ts`.

Nowe: `linijka-podzialka.ts`, `linijka-pozioma.ts`, `linijka-pionowa.ts`,
`widok-nastawy-operatora.ts`, `widok-skali.ts`, `widok-ukladu-stron.ts`,
`widok-podzialu-powierzchni.ts`, `widok-pasek-widoku.ts`, `widok-druku.ts`,
`widok-i-linijka.test.ts`.

## Rozstrzygnięcia, które warto znać przed dalszą pracą

1. **Powierzchnia należy do dokumentu.** `.ms-praca__kolumny` ma domyślnie JEDNĄ
   kolumnę; dymki komentarzy i panel Redaktora stoją nakładką nad treścią. Stałe
   kolumny zostały jako `data-uklad-paneli='kolumny'` — tryb do wyboru Operatora,
   przestawiany z paska widoku. Historia wersji, wstążka PDF i narzędziownia
   cyfryzacji przestały być stałymi kolumnami: każda jest wyzwalaczem plus
   nakładką.
2. **Pasek widoku należy do powierzchni, nie do wstążki.** Skala, układ kartek,
   przewijanie, linijki, jednostka, granice marginesów, skok o kartkę, nastawy
   kartki i nastawy druku stoją w jednym miejscu — przy powierzchni. Wstążka bierze
   ten pasek gniazdem (`GniazdaWstazki.widok`, nieobowiązkowym), więc nastawy nie
   mają dwóch miejsc, które mogłyby się rozjechać.
3. **Nadpisania nastaw strony.** Okno pracy pcha całą `StronaPracy` przy każdej
   swojej zmianie. Powierzchnia trzyma osobno to, co Operator ustawił chwytem
   linijki albo paskiem, i przywraca to po pchnięciu okna — z wyjątkiem pól, które
   okno tym pchnięciem naprawdę zmieniło. Bez tego chwyt marginesu przeżywałby do
   pierwszego naciśnięcia czegokolwiek na wstążce.
4. **Skala idzie `zoom`, nie `transform: scale`.** Przekształcenie nie zmieniało
   miejsca w układzie, więc kartki nachodziły na siebie powyżej 100 %.
5. **Tryb źródłowy jest przełącznikiem**, nie czwartym przyciskiem trybu — zachował
   znacznik `data-czynnosc='tryb-zrodlowy'`, bo to ta sama czynność w innym
   kształcie (sprawdzian `praca-z-dokumentem.test.ts` przechodzi).

## Wykaz odbioru — czynność po czynności

**Wstążka:** osiem zakładek z nazwanymi grupami — jest. Zakładka „Widok"
przebudowana: tryby, przełącznik trybu źródłowego, gniazdo paska widoku, panele na
żądanie — jest. Zakładka kontekstowa PDF z pięcioma grupami (Strony, Nakładanie,
Treść, Bezpieczeństwo, Narzędziownia cyfryzacji), domyślnie ukryta, wchodząca
przyciskiem i samoczynnie przy dokumencie PDF — jest.

**Widok strony:** kartka o rozmiarze nośnika — jest; marginesy cztery osobno —
jest; oprawa i strona oprawy — jest; marginesy odbicia — jest; nastawy gotowe
marginesów (wąskie, normalne, szerokie, do oprawy) — jest; granica marginesu
w treści jako przełącznik — jest; podział na strony z podziałem jawnym — jest;
paginacja, nagłówek i stopka w miejscu pracy — jest; nagłówek i stopka OSOBNE dla
pierwszej strony i stron parzystych — **nie ma**: `StudioPageSetup` niesie po
jednym polu `header` i `footer`, brak pozycji kontraktu (zgłoszone niżej).

**Nośniki:** A0–A6, B1–B5, Letter, Legal, Tabloid, koperty DL, C4, C5, C6 — jest;
format własny w milimetrach — jest; orientacja pionowa i pozioma — jest; nastawy
osobne dla SEKCJI — **nie ma**: sekcji nie ma w kontrakcie ani w warstwie danych
klienta (odcinek 4). Nadruk koperty (adresat, nadawca, ich położenie) — **nie ma**,
brak pozycji kontraktu. Wykazu nośników nie skopiowano na sztywno: wbudowany jest
wykazem ZAPASOWYM, a `wchlonNosnikiRdzenia` nadpisuje go wykazem
`design.print.paper.list`, gdy ten dojedzie.

**Linijki:** pozioma i pionowa — jest; milimetry albo cale jako wybór Operatora —
jest; chwyty marginesów (cztery) — jest; wcięcie pierwszego wiersza, lewe i prawe
osobnymi znacznikami, z wcięciem pierwszego wiersza liczonym względem lewego —
jest; tabulatory zakładane naciśnięciem, cztery rodzaje w obiegu, znak wiodący,
zdejmowanie — jest; szerokości kolumn tabeli chwytem — jest (nastawa widoku:
składnia tabeli szerokości nie niesie i kontrakt nie ma jej gdzie zapisać);
położenie kursora i granice zaznaczenia na linijce — jest; przełącznik pokazania
i ukrycia — jest; wszystkie chwyty dostępne z klawiatury jako `slider` — jest.
Chwyt WYSOKOŚCI wiersza tabeli na linijce pionowej — świadomie **nie ma**: wysokość
bierze się z treści, a chwyt bez nastawy do przestawienia byłby pozorny.

**Skala:** procenty polem i suwakiem — jest; nastawy gotowe (do szerokości strony,
cała strona, do szerokości tekstu, 100 %, plus szereg 50–200 %) — jest; skala
pamiętana PRZY DOKUMENCIE — jest.

**Widok stron:** jedna kartka, wiele kartek w rzędzie (do ośmiu), rozkładówka ze
stroną pierwszą samą po prawej — jest; przewijanie ciągłe albo strona po stronie
(zaczepienie CSS plus skok o kartkę) — jest; licznik „kartka N z M" liczony
z przewinięcia — jest.

**Podgląd wydruku jako TRYB tego okna:** jest — pokazuje nośnik i orientację
naprawdę ustawione, pisanie wyłączone, render rdzenia (`studio.preview.render`)
wołany osobno przez okno pracy.

**Drukowanie (czynność Operatora):** droga druku z klienta — jest; wydruk bierze
kartki podglądu, nie surowy tekst — jest; zakres stron (wszystkie, bieżąca,
podany) — jest; skala druku osobna od skali widoku — jest; adiustacja na wydruku
albo tekst po zmianach — jest; szybkie drukowanie ostatnimi nastawami — jest;
nazwany brak PRZED próbą, gdy powłoka nie ma okna drukarki — jest. Liczba kopii
i druk dwustronny — **przenoszone jako życzenie do okna drukarki systemu, nie
narzucane przez stronę**; okno mówi to wprost.

**Dwa dokumenty — dwa równorzędne tryby:** zakładki z niezależnymi migawkami —
jest (było, zachowane); podział powierzchni pionowy i poziomy z przestawialną
granicą, dostępną z klawiatury, z polem czynnym i gniazdami stopek na własny pasek
statusu i wiersz polecenia każdego pola — **komponent gotowy i sprawdzony**; wybór
trybu jawny, odwracalny i pamiętany — jest. **Czego nie ma:** osadzenia drugiej
ŻYWEJ powierzchni w oknie pracy — wymaga zmiany w `okno-pracy-z-dokumentem.ts`
i drugiego stanu dokumentu, a ten plik należy do odcinka 3 i jest w trakcie prac.
Wywołanie do domknięcia: `podzial.ustawPola(powierzchniaA.element,
powierzchniaB.element)` plus drugi `StanStudio`.

**Historia wersji (Session Repository):** pod przyciskiem, jako nakładka — jest;
czas, autor, kropka stanu wraz ze słowem, etykieta własna — jest; Podgląd,
Przywróć, Porównaj — jest; „Przywróć" nie usuwa wersji nowszych i mówi to — jest;
menu pozycji (etykieta, eksport, odwołanie, usunięcie miejscowe i jawne) — jest;
filtr wedle autora i etykiety plus filtr wersji kluczowych — jest; zapisy
samoczynne osobnym szeregiem, domyślnie ukryte, z przełącznikiem — jest.
**Czego nie ma:** trzech wywołań do rdzenia — `studio.version.label.set`,
`studio.repository.export`, `studio.package.export`. Komendy w rdzeniu SĄ;
brakuje drogi po stronie klienta: `ZrodloStudio` niesie sześć komend i tych nie ma,
a `modul-studio.ts` (nie plik tego odcinka) nie podaje zaplecza. Okno przyjmuje
nieobowiązkowe `ZapleczeHistorii` i bez niego pokazuje brak NAZWANY. Domknięcie to
jedno wywołanie w `modul-studio.ts`.

**Narzędziownia cyfryzacji:** panel przestał być stałą kolumną — wchodzi
przyciskiem jako nakładka z kolejką, rozpoznaniem, czyszczeniem obrazu, podglądem
i przekazaniem do edytora — jest. **Czego nie ma:** przełożenia panelu z komendy
`document.text.extract` na rodzinę `studio.ingest.*` (osiem komend, wszystkie
w kontrakcie): wybór silnika, zestaw języków, próg pewności, poprawka rozpoznanych
słów przed przyjęciem, kolejka po stronie rdzenia, `ingest.url` wprost do
dokumentu, wykaz i odczyt urządzeń. `braki-cyfryzacji.ts` przestał kłamać, że tego
nie ma w kontrakcie — każda pozycja nazywa komendę, która czeka, i mówi, że brak
jest po stronie okna. Osadzenie narzędziowni w grupie wstążki PDF czeka na jedno
wywołanie: `utworzOknoWarsztatuDokumentu(stan, zrodlo, cyfryzacja.element)`.

## Braki pozycji kontraktu — do odcinka 1

Wszystkie dotyczą nastaw, które okno prowadzi przez sesję i traci przy zapisie:

1. `StudioPageSetup` bez: marginesu na oprawę, strony oprawy, marginesów odbicia,
   rozmiaru własnego w milimetrach, kolumn, znaku wodnego, formatu i punktu startu
   numeracji stron, jednostki miary.
2. Nagłówek i stopka **osobne dla sekcji, pierwszej strony i stron parzystych** —
   dziś po jednym polu na cały dokument.
3. **Sekcje** o własnych nastawach strony — bez nich „format osobno dla sekcji"
   nie ma czego adresować.
4. **Wcięcia akapitu, tabulatory (rodzaj i znak wiodący) oraz szerokości kolumn
   tabeli** — chwyty linijki działają, ale nie mają pola w kontrakcie, więc giną
   przy zapisie. To jest najboleśniejszy z tych braków, bo linijka bez trwałości
   wygląda na usterkę.
5. **Nastawy widoku dokumentu** (skala, układ kartek, przewijanie, jednostka,
   widoczność linijek) — dziś w `localStorage`; gdy kontrakt dostanie pozycję,
   `widok-nastawy-operatora.ts` przełoży się bez zmiany wołaczy, bo magazyn jest
   podawany, a nie brany na sztywno.
6. **Koperty C4, C5 i C6** w wykazie nośników rdzenia
   (`design.print.paper.list`) — Właściciel wymienia C5 i C4 wprost. Wykaz leży
   w pliku modułu Design, w który nie wolno wchodzić: należy przenieść go
   do miejsca wspólnego dla Studia i Designu i tam uzupełnić.
7. **Drukowanie** — świadomie bez komendy, decyzją Właściciela („póki co tylko
   funkcja dla Operatora"). Gdyby kiedyś miało dojść drukowanie po stronie serwera,
   potrzebne są trzy rzeczy: komenda wykazu drukarek widzianych przez maszynę,
   komenda wydruku przyjmująca dokument wraz z nastawami (zakres, kopie, dupleks,
   nośnik) oraz rozstrzygnięcie, CZYJA drukarka drukuje — serwerowa czy Operatora;
   bez tego trzeciego drukowanie z rdzenia wydrukuje pismo w serwerowni.

## Pominięcia świadome, nie przeoczenia

- **Gałęzie dokumentu** pominięte w warstwie okna decyzją Właściciela; komendy
  `studio.branch.create`, `.list`, `.merge` pracują dalej i nie są martwe.
- **Wyszukiwanie znaczeniowe** bez okna decyzją Właściciela;
  `studio.search.semantic` pracuje dalej.
- **Paczka redakcyjna przekazania** została pozycją menu bez osobnego warsztatu,
  wedle rozstrzygnięcia o jej opcjonalności.

Następny wykonawca niech nie dobudowuje tych trzech z własnej inicjatywy.

## Pakiet serwera

Potwierdzam ustalenie odcinka 3: wykaz zależności pakietu nie wymienia Tesseractu,
polskiego pakietu językowego, 7-Zipa ani silnika mowy. Wedle rozstrzygnięcia
Właściciela wszystkie trzy są składnikami pakietu serwera, więc wykaz wymaga
uzupełnienia. Treść odmów w `braki-cyfryzacji.ts` przestawiona na to znaczenie:
brak składnika jest **usterką wdrożenia serwera**, nie ograniczeniem produktu.

## Kompilacja i sprawdziany

- `(cd client && npx tsc --noEmit -p tsconfig.json)` — **czysto**.
- `npx vitest run` (całość klienta) — **265 sprawdzianów pomyślnych, 20 plików**,
  w tym 69 nowych tego odcinka. Sprawdzian cudzy zgłoszony wyżej przez odcinek 3
  (`praca-z-dokumentem.test.ts`, znacznik `tryb-zrodlowy`) — **naprawiony**:
  przełącznik trybu źródłowego niesie ten znacznik dalej.
- `go build ./...` — **czerwone, nie z tego odcinka**: `server/internal/dane`
  ma podwójne deklaracje między `studio_postac_*.go` i `studio_forma_dokumentu.go`
  / `studio_kontrola_pracy.go` / `studio_znakowanie_wykonawcy.go`
  (`StylNazwanyStudia`, `SekcjaDokumentuStudia`, `BlokadaFragmentuStudia`,
  `BlokadaSzablonuStudia`, `CzynnoscDokumentuStudia`, `KopiaZapasowaStudia`,
  `NastawaPracyStudia`, `ZnakowanieStudia`, `RodzajZnacznikaStudia`). To zderzenie
  dwóch innych odcinków w warstwie danych — do rozstrzygnięcia między nimi.

# Odcinek 4 — podjęcie po urwanej sesji (17.08.2026)

Sesja odcinka 4 urwała się w trakcie pisania, zostawiając `go build ./...`
czerwony w `internal/core`: napisane były wołacze, brakowało pomocników,
a trzy miejsca rozjeżdżały się z warstwą danych. Zderzenie deklaracji
w `internal/dane`, zgłoszone wyżej przez odcinek 3, w chwili podjęcia
**już nie występowało** — tamta warstwa buduje się czysto.

## Co zostało domknięte

**Pomocniki obszaru postaci** (`adapter_modul_studio_postac_pomocniki.go`):
`postacWskaznikLiczby64`, `postacSkladnica`, `postacNastawyDomyslne`,
`postacZlozSekcje`, `postacZlozObiekty`, `postacZlozAparat`, `postacZlozPola`,
`postacZlozBlokady`. Trzy rozstrzygnięcia w nich zawarte:

- Nastawa domyślna strony jedzie z `wejscieDomyslneNastawyStrony`, nie z drugiego
  wykazu — dwa wykazy domyślnych rozjechałyby dokument wczytany z założonym.
- Składacze czytają najpierw pole JSON, a POTEM nadpisują pola kolumnowe:
  kolumna jest prawdą, JSON niesie tylko to, na co kolumny nie ma.
- Blokady składa `blokadaZlozKontrakt` obszaru kontroli pracy. Drugie przełożenie
  tej samej tabeli byłoby drugą prawdą o jednym wierszu.

**Uzgodnienie z warstwą danych.** Blokady fragmentów i dziennik czynności nie
stoją w `RepozytoriumPostaciStudia`, tylko w `KontrolaPracyStudia` — oba wołania
przełożone na `a.kontrolaSkladnica()`. Pola wpisu dziennika nazwane tak, jak je
niesie `dane.CzynnoscDokumentuStudia` (`DokumentKod`, `StanPrzed`, `StanPo`,
`ZmianaSledzonaK`). Styl nazwany czyta i pisze `PostacZnakuJSON`
i `PostacAkapituJSON`.

**Dwie usterki, które kompilacja przepuszczała, a baza odrzuciłaby przy pierwszym
zapisie** — dlatego wypisane osobno, bo nie były błędami składni:

1. Stan wpisu dziennika był podawany jako `wniesiona`, a tabela zna wyłącznie
   `active` i `reverted` (migracja 364). Każda czynność na postaci wywracałaby się
   na ograniczeniu — po zapisaniu zmiany, czyli w najgorszym możliwym miejscu.
2. Rodzaj czynności dziennika był podawany rodzajem zmiany ŚLEDZONEJ
   (`wstawienie`, `usuniecie`, `formatowanie`), a dziennik ma własny słownik
   jedenastu wartości (`StudioActionKind`: `textEdit`, `formatChange`,
   `styleChange`, …). Oba słowniki są teraz **osobnymi polami** `postacZakoncz`,
   bo Operator cofa „zmianę stylu", nie „formatowanie". Jedenaście wołaczy podaje
   rodzaj czynności właściwy dla tego, co robią.

**Wpięcie do rejestru.** 21 komend miało metodę w rdzeniu i **nie było wpiętych**,
więc kontrakt widział je jako brak funkcji: `document.form.get`, `form.save`,
`text.get`, `text.edit`, dziewięć z rodziny `format.*`, cztery `style.*`
i trzy `lock.*`. Zdarzenia `studio.document.changed` te czynności nie rozgłaszają
— każda oddaje postać wprost w odpowiedzi, a przy pisaniu litera po literze
zdarzenie byłoby przeładowaniem okna na każde naciśnięcie klawisza. Wyjątkiem
jest `document.form.save`, który oddaje także dokument wraz z wersją.

## Stan sprawdzianów

- `go build ./...` — **czysto** (było: 27 błędów w trzech plikach).
- `go vet ./internal/core ./internal/dane` — czysto. `gofmt` — pliki tego
  podjęcia czyste.
- `go test ./internal/dane ./internal/store` — pomyślnie.
- `go test ./internal/core` — **jeden sprawdzian czerwony**:
  `TestRejestrPokrywaKomendyKontraktu`, komendy bez obsługiwacza **90 z 1067**
  (przed podjęciem 111). Żaden inny sprawdzian rdzenia nie stoi.

## Czego NIE MA — 90 komend bez metody w rdzeniu

To nie jest brak wpięcia, a brak napisanej czynności. Rozbicie po rodzinach:
`page` 9, `document` 8 (wejście, wydanie, kopia, założenie), `template` 7,
`markup` 7, `table` 6, `agents` 6, `list` 5, `symbol` 4, `object` 4,
`apparatus` 4, `section` 3, `model.changes` 3, `journal` 3, `field` 3,
`backup` 3, `autosave` 3, `view` 2, `version` 2, `insert` 2, `diff` 2,
`clipboard` 2, `ruler` 1, `provenance` 1.

Wedle podziału prac należą do odcinka 4 (`page`, `section`, `table`, `list`,
`object`, `apparatus`, `field`, `symbol`, `ruler` — 39), odcinka 5 (`template`,
`document.import/export`, `insert` — ok. 17) i odcinka 6 (`journal`, `autosave`,
`backup`, `markup`, `model.changes`, `diff.hunk`, `agents`, `clipboard`,
`provenance`, `version`, `view` — pozostałe). Pełny wykaz wypisuje sam
sprawdzian: `go test ./internal/core -run TestRejestrPokrywaKomendyKontraktu`.

# Studio — cztery odcinki oddane i wpięte (17.08.2026, wieczór)

Po podjęciu urwanej sesji ruszyły cztery odcinki naraz, na rozłącznych plikach,
z rejestrem komend trzymanym poza ich zasięgiem. Rejestr wpiął prowadzący —
jednym przebiegiem, po oddaniu wszystkich czterech.

## Miara odbioru

- **`studio.*` — 179 komend kontraktu, 179 z obsługą w rdzeniu.**
  `TestRejestrPokrywaKomendyKontraktu` **przechodzi**: zero komend całego
  kontraktu (1067) bez obsługiwacza. Przed podjęciem było 111.
- `go build ./...` czysto; `go test ./internal/core ./internal/dane
  ./internal/store ./internal/protocol` — **wszystko pomyślnie**.
- Klient: `npx tsc --noEmit` czysto, `npx vitest run` — **302 sprawdziany, 22
  pliki, wszystkie pomyślne**.
- Zapory nietknięte i zielone, w tym `zapora_warsztatu_pdf_test.go`.

## Co przyniósł który odcinek

**Strona, sekcje, listy, symbole (22).** Nastawy nośnika wraz z oprawą,
odbiciem, formatem własnym w milimetrach i kolumnami; nagłówek i stopka
w **trzech osobnych zasięgach** — zwykłe, pierwsza strona, strony parzyste;
numeracja stron z formatem, punktem startu i wznowieniem w sekcji; znak wodny;
podział strony i sekcji; nadruk koperty; tabulatory linijki wraz ze znakiem
wiodącym. Listy trzymają w akapicie kod listy i poziom, nie kopię nastaw, więc
zmiana definicji przestawia wszystkie miejsca. Tablica ~200 symboli w siedmiu
grupach; autozamiana czytana z bazy, nie z kodu.

**Tabele, obiekty, aparat, pola (17).** Tabela z policzonymi szerokościami,
sześć czynności na budowie, scalanie i podział komórek, powtarzanie wiersza
nagłówkowego, sortowanie całymi wierszami, zamiana tabeli w tekst w obie strony
z bilansem strat. Aparat: 13 rodzajów — spis treści, spisy ilustracji i tabel,
indeks, bibliografia z numeracją powołań, przypisy, ze znacznikiem nieświeżości.
Obiekty: rozmiar z proporcją, przycięcie, opływanie, warstwa, obrót.

**Wejście, wydanie, szablony (17).** Założenie dokumentu, wniesienie pliku,
PDF-u i obrazu, zapis pod nazwą, kopia, wydanie do txt, md, html, rtf, docx, odt
i pdf, wydanie wsadowe, warsztat szablonów wraz z polami do wypełnienia.
OOXML i ODF przez `archive/zip` i `encoding/xml`, PDF przez `pdfcpu` — **ani
jednego programu zewnętrznego**. Naprawione dwie rzeczy z zastanego kodu: docx
nie wypisywał nagłówka ani stopki do archiwum, a postać dokumentu miała **dwa
magazyny** — przełożona na jeden.

**Dziennik, znakowanie, agenci, schowek (34).** Cofanie pojedyncze i nie po
kolei, nakładające **różnicę drzew** byt po bycie, nie migawkę — praca naniesiona
po cofanej czynności zostaje. Zależność czynności odmawia i nazywa kod. Zmiany
modelu wraz z postacią, skakaniem i cofnięciem wybranych. Różnica dwóch wersji na
postaci i przeniesienie pojedynczego fragmentu. Autozapis osobnym szeregiem,
z kopią zakładaną PRZED zapisem. Znakowanie: propozycja na fragmencie
zablokowanym przechodzi — to jedyna droga wykonawcy. Zajęcia fragmentów wraz
z odmową nazywającą wykonawcę i czas. Naprawione zaszycie autora jako
`uzytkownik` w trzech miejscach `adapter_modul_studio_adnotacje.go`.

## Zapora blokad przestała być zapisem bez skutku

`zaporaBlokadStudia` stała napisana i **nie wołał jej nikt**, więc blokada
fragmentu nie zatrzymywała niczego. Wołanie stoi teraz w `kompozycja.go` po obu
wpięciach Studia — zapora owija to, co w rejestrze już jest, więc kolejność jest
częścią jej działania, nie szczegółem.

## Czego nie ma — do rozstrzygnięcia Właściciela

1. **Okno.** 115 komend `studio.*` nie jest wołanych z klienta. Rdzeń niesie
   funkcje, których Operator jeszcze nie dosięga — to jest teraz największy dług
   modułu i osobna robota.
2. **Wykaz nośników** stoi w pliku modułu Design i nie ma kopert C4, C5, C6.
   Właściciel wymienia C5 i C4 wprost. Do przeniesienia w miejsce wspólne.
3. **Paginacja** spisu treści liczy się silnikiem podglądu, który ma A4
   i marginesy zaszyte na stałe — numer strony zgadza się z podglądem, nie
   z prawdziwymi nastawami strony.
4. **Aparat dokumentu nie cofa się dziennikiem** — brak przekładu „element
   kontraktu → wiersz". Cofnięcie oddaje bilans nazywający brak, nie ciszę.
5. **Malarz formatów żyje w pamięci procesu**, choć migracja 368 ma na to tabelę
   z wygasaniem. Postać zabrana przepada przy przeładowaniu rdzenia.
6. **Tożsamość agenta nie jest stemplowana** na drogach postaci
   (`postacOdlozZmiane`), więc rozbicie zmian po wykonawcy pokaże je jako
   nienazwanego. Poprawka to jedno wywołanie w `postacZakoncz`.
7. **Porządek alfabetyczny** sortowania tabeli nie zna pisma polskiego (ł, ą).
8. **Postać przez schowek** — historia schowka platformy niesie sam tekst, więc
   wklejenie przejmuje postać miejsca, nie źródła; bilans mówi to wprost
   i kieruje do malarza formatów. Przeniesienie postaci wymaga kolumny
   w `wpis_schowka`.
9. **Brak pozycji kontraktu**: `StudioFragmentLock` bez `templateId`;
   `StudioParagraphFormat` bez wznowienia numeracji; `StudioActionKind` bez
   rodzaju „zmiana pola"; moduł przy `AodSuggestion` i `AodStatus`.

# Domknięcie okna — 98 % kontraktu osiągalne dla Operatora (17.08.2026, wieczór)

Po odcinkach rdzenia poszło sześć odcinków okna: trzy na Studio, jeden na
prowenancję i Research, jeden na dziewięć komend pojedynczych, jeden na kontrakt.
Montaż wykonał prowadzący, bo pliki złożenia są jedynym miejscem, w którym te
odcinki się spotykają.

## Miara

- **1077 komend kontraktu, 1061 z drogą z okna (98 %).** Przed tą turą: 923.
- Rdzeń: zero komend bez obsługiwacza (`TestRejestrPokrywaKomendyKontraktu`).
- Klient: `npx tsc --noEmit` czysto, `npx vitest run` — **503 sprawdziany
  w 40 plikach, wszystkie pomyślne** (przed turą 302).
- `go build ./...` czysto; sprawdziany rdzenia, danych, magazynu i protokołu
  pomyślne. Zapory nietknięte.

## Co dostał Operator

**Numeracja stron i marginesy — rozstrzygnięcie Właściciela wykonane.** Pięć
stylów numeracji, punkt startu od dowolnego numeru, wznowienie w sekcji,
„strona N z M", sześć umiejscowień. Marginesy: cztery osobno, oprawa, marginesy
odbicia, nastawy gotowe, format własny w milimetrach, kolumny — osobno dla każdej
sekcji. **Chwyty linijki dojeżdżają teraz do rdzenia**, więc nastawa przestała
ginąć przy zapisie.

**Postać dokumentu w oknie:** arkusz stylów nazwanych, 17 cech postaci znaku,
postać akapitu, malarz formatów, znajdź-i-zamień z postacią, listy wraz
z numeracją wielopoziomową, tablica symboli, tabulatory liczone różnicą wobec
wykazu. Nastawy widoku przełożone z zapisu przeglądarki na rdzeń **bez zmiany
ani jednego wołacza** — po to magazyn był podawany.

**Wstawienia:** tabele (siatka, budowa, sortowanie, zamiana w tekst), obiekty
(rozmiar z proporcją, opływanie, warstwa, obrót), aparat dokumentu (13 rodzajów,
odświeżanie, znacznik nieświeżości), pola. Wykres nazywa brak i kieruje do
modułu Design, zamiast udawać.

**Kontrola pracy:** dziennik czynności wraz z cofaniem i ponawianiem, przełącznik
zmian modelu, różnica dwóch wersji na postaci i przeniesienie fragmentu, kopie
zapasowe i autozapis osobnym szeregiem, znakowanie trwałe, blokady fragmentów,
zajęcia wykonawców, schowek dokumentu.

**Poza Studiem:** prowenancja wywołań modeli wraz z **oceną wywołania** i wydaniem
śladu; przestrzeń badania, załączniki źródeł i książka kodów; komponenty własne,
procesy z telefonu, materiał, zużycie i koszt, pakowanie archiwum.

**Kontrakt:** wyłączenia pamięci w sześciu zasięgach (całkiem, środowisko,
projekt, moduł, para modułów, karta sesji), wyciszenie nakładki jako byt rdzenia
wraz z rozgłoszeniem, moduł przy stanie i sugestii nakładki, nośnik sygnału klas
zdarzeń, trzy pozycje Studia. Sterowanie z konfiguracji **bez drugiej drogi
komend** — rodzina `config.*` niosła to zasięgiem ogólnym, brakowało wierszy
katalogu (migracja 376).

## Sześć zdań, które kłamały — poprawione

Wykaz braków w tym produkcie mylił się w OBIE strony. Poprawione zdania: zakładka
prowenancji („kontrakt nie niesie rodziny"), warstwa mobilna („nie jest wpięta"),
zużycie („kontrakt nie niesie ani pola tokenów"), dwa wpisy Biblioteki o materiale
i archiwum, wyróżnienie tekstu w oknie pracy („nie dojeżdża do rdzenia") oraz
zakładka dokumentu („zapisu kontrakt nie niesie" — zakładka jest elementem
aparatu i **od teraz zapisuje się w rdzeniu**).

## Zostało 16 komend bez drogi z okna — i wszystkie z powodu

- **`mail.*` (10)** — rodzina wymyślona przez wykonawcę; dostawa v2.0 nie ma
  modułu Poczty, ma pocztę jako konektor Asystenta i import EML w Bibliotece.
  Czeka na rozstrzygnięcie: wyciąć albo przenieść.
- **`memory.disable.list` i `.set` (2)** — wniesione dziś; sterowanie należy do
  zakresu 5.10 okna konfiguracji i jest następną robotą.
- **`memory.detach` (1)** — chwyt należy do okna Memory & Context Manager.
- **`aod.signal.report` i `.list` (2)** — sygnał zgłasza moduł wytwarzający
  zdarzenie, nie okno nakładki; do wpięcia u zgłaszających.
- **`advisor.consult` (1)** — czynność modelu w trakcie tury, nie Operatora.

# Arsenał serwera: pakiet wdrożeniowy stawia programy, na których stoi rdzeń

Model wdrożenia jest zamknięty — programy jadą WRAZ Z APLIKACJĄ NA SERWER,
u Operatora zostaje cienka instalka. Serwer nie miał czym ich postawić:
`scripts/instalka-windows.sh` buduje samą powłokę i ma taki zostać, a wykaz
30 zależności zewnętrznych istniał tylko jako diagnoza po fakcie — rdzeń mówił
Operatorowi, czego mu brak, i na tym się kończyło. Zbudowane:
`scripts/arsenal-serwera.sh` wraz z dwoma trybami wykazu w rdzeniu.

## Lista nie została przepisana — jest wyprowadzona z rejestru

Nagłówek `zaleznosci_zewnetrzne.go` ostrzega, że druga lista rozjedzie się
z pierwszą przy pierwszej zmianie pakietu. Dlatego skrypt **nie zna żadnej nazwy
pakietu**: woła rdzeń i konsumuje wynik.

- `danaco-console --wykaz-zaleznosci` — komplet 30 pozycji, po wierszu, pola
  rozdzielone tabulacją: `warstwa program pakiet stoi nazwa zakres`.
- `danaco-console --wykaz-mowy` — arsenał mowy spoza wykazu zależności.
- Oba w NOWYM pliku `server/internal/core/zaleznosci_wykaz_wydruk.go`; adaptery
  modułów nietknięte, deklaracje narzędzi czytane tam, gdzie stoją.

**Warstwę liczy rdzeń, nie powłoka.** Rozdział na obowiązkową i decyzyjną jest
rozstrzygnięciem, nie formatowaniem: gdyby robił go skrypt dopasowaniem napisów,
byłby drugą regułą obok deklaracji — nietypowaną i cichą przy pomyłce. Jest
jedną funkcją (`WarstwaZaleznosci`) z jednym sprawdzianem.

## Co pakiet serwera niesie — 17 + 5 + 2 + 2 + 2 + 2 = 30

- **obowiązkowa (apt), 17** — Tesseract WRAZ z `tesseract-ocr-pol`, 7-Zip,
  eSpeak NG, Pandoc, ffmpeg, ffprobe, LibreOffice, Chromium, **ImageMagick**,
  poppler ×2, OpenSSH, ShellCheck, shfmt, picocom, telnet, łańcuch Go. Skrypt
  dokłada `python3`, `python3-venv`, `python3-pip` — interpretera rdzeń nie woła
  jako narzędzia, więc w wykazie go nie ma, a bez niego mowa nie rusza.
- **warsztat Go, 5** — gopls, goimports, golangci-lint, staticcheck, Delve.
  Moduł Developer pracuje na serwerze, więc jego warsztat należy do serwera.
- **warsztat npm, 2** — Prettier, ESLint. **snap, 2** — kubectl, PowerShell.
- **kroki ręczne, 2** — Real-ESRGAN i rembg: wydania spoza repozytoriów
  dystrybucji, drukowane z treścią pola `Pakiet`, nie zgadywane.
- **decyzyjna, 2** — silnik kontenerów (obie deklaracje: warsztat Developera
  i moduł Terminal). WSTRZYMANY decyzją Właściciela, więc skrypt **nie stawia go
  milcząco**; wymaga `DANACO_SILNIK_KONTENEROW=tak`. Pilnuje tego sprawdzian.

## Mowa: wagi pobierane przy stawianiu, nie przy pierwszym mikrofonie

faster-whisper ściągał wagi przy pierwszym użyciu — pierwsze nagranie
u Operatora czekałoby na sieć. Prowizjonowanie pobiera je z góry, w rozmiarze
`mowa.ModelDomyslny` (dziś „small", **około 480 MB**; skrypt po pobraniu mierzy
katalog i wypisuje rozmiar zmierzony, żeby liczba była pomiarem, nie obietnicą).
Biblioteka idzie przez `pip install -r pomocniki/transkrypcja/wymagania.txt` —
nazwa i przypięta wersja stoją TAM, skrypt ich nie powtarza. Piper i jego głosy
`.onnx` zostają krokiem ręcznym z podanym miejscem arsenału i zmienną wskazania.

## Zmierzone, nie zapowiedziane

`go build ./...` czysto; osiem nowych sprawdzianów przechodzi; `gofmt`
i `shellcheck` bez uwag. Wykaz wypisał 30 z 30, `arsenal-serwera.sh sprawdz`
zgodnie z sondą startową. **Na tej maszynie nie postawiono niczego** — tryb
`postaw` bez roota odmawia i nic nie rusza. Otwarte: `postaw` nie był uruchomiony
na żadnym serwerze (droga apt/snap/venv niesprawdzona w boju); głosy pipera
i wagi Real-ESRGAN/rembg zostają ręczne; wag modelu mowy nie ma w katalogu
arsenału `/opt/danaco-arsenal/modele-mowy` — tu leżą w domyślnym cache'u.

# Studio, odcinek 6 — rdzeń: kontrola, historia, bezpieczeństwo pracy (18.08.2026)

Zakres: blokady fragmentów, odwracalny dziennik czynności, przełącznik
podświetlający zmiany wykonawców, różnica postaci dwóch wersji, autozapis
osobnym szeregiem, kopie zapasowe, znakowanie, zajęcia i spięcia wykonawców.

## Liczby

- **32 komendy wpięte** (wszystkie rodziny odcinka), obsługiwacz dla każdej;
  `TestRejestrPokrywaKomendyKontraktu` zielona.
- **Warstwa danych:** `dane/studio_kontrola_pracy.go` i
  `dane/studio_znakowanie_wykonawcy.go` nad tabelami migracji 363-366 i 370
  oraz kolumnami szeregu wersji z 367. **Migracji nie zakładałem żadnej** —
  wszystkie tabele przyszły z odcinka kontraktu.
- **42 funkcje sprawdzianowe w 8 plikach, 46 mierzonych przypadków**
  (`autor_na_drogach_test.go` to jedna funkcja o pięciu podprzypadkach):
  17 skutku czytających bazę osobnym połączeniem, 24 rachunku, 5 pokrycia
  zapisu autora. Wynik: **wszystkie zielone, 0 czerwonych**
  (`go test -run 'TestBlokada|TestDziennik|TestZmianyModelu|TestAutozapis|TestKopia|TestZnakowanie|TestAgenci'`).
- **Kompilacja** `go build ./...` zielona, `go vet` czysty, `gofmt` czysty,
  `grep` zależności `go.mod` daje 8. Kontraktu **nie tknąłem** — 1077 komend,
  `git diff shared/contract.json` puste.
- **Pełny przebieg `go test ./server/internal/core/` — ZIELONY, dwa razy:**
  506 s ze stemplowaniem autora i wywiedzioną zależnością czynności, oraz
  **495,6 s i 498,5 s po dołożeniu siatki śladu** (`sladWykonawcyStudia`) —
  dwa niezależne przebiegi na komplecie odcinka, oba EXIT=0. Wszystkie objęły zaporę
  `TestRejestrPokrywaKomendyKontraktu`. Przebiegów pośrednich, które wywróciły
  się na `signal: terminated`, nie liczę do wyniku — ubiłem je sam,
  porządkując procesy współbieżnych odcinków.

## Dwa rozstrzygnięcia, których nie da się cofnąć bez skutku dla produktu

### 1. Zapora blokad stoi w REJESTRZE, nie w obsługiwaczach

`zaporaBlokadStudia` owija w rejestrze każdą komendę `studio.*`, która może
zmienić dokument — wykazem **WYJĄTKÓW** (odczyty), nie wykazem objętych.
Komenda dopisana przez innego wykonawcę i niewpisana do wykazu jest domyślnie
**sprawdzana**; przy wykazie objętych byłaby cichą dziurą.

Trzy odpowiedzi zderzenia, nie dwie: całość w blokadzie → odmowa nazwana;
część → zmiana poza blokadą plus bilans dopisany do odpowiedzi (scalany
z bilansem obsługiwacza, nie nadpisujący go); poza blokadą → cisza.

Wpięcie idzie **po** `zarejestrujStudio` — owija to, co w rejestrze już stoi.
Przesunięcie go przed wpięcia zostawiłoby zaporę napisaną i nieaktywną.

### 2. Autor wykonawcy stemplowany w ładunku (`podpisemWykonawcy`)

**To była usterka i została naprawiona, nie zgłoszona.** Czynności postaci
rozstrzygają autora z pola żądania (`postacAutor`). Wykonawca, który pola
`author` nie podał, zapisywał się jako Operator — a wtedy jego zmiana **nie
odkładała się jako zmiana śledzona i nie dawała się podświetlić**. Zmierzone:
przed naprawą `zmiana_sledzona_studio` miała **zero** wierszy po edycji z
gniazda serwera narzędzi.

Naprawa pierwsza: gdy fakt gniazda mówi „wykonawca", wpięcie rejestru dopisuje
`author: model` do ładunku, zanim ładunek zobaczy obsługiwacz. Obejmuje
wszystkie komendy `studio.*`, także te, których jeszcze nikt nie napisał.

**Naprawa druga — siatka `sladWykonawcyStudia`.** Zmierzyłem pokrycie na PIĘCIU
drogach zamiast na jednej (`autor_na_drogach_test.go`) i pierwsza naprawa nie
starczyła: `studio.document.save` zawołane przez wykonawcę zmieniało treść
i dalej NIE odkładało ani zmiany śledzonej, ani wpisu dziennika — bo ta droga
nie prowadzi przez `postacZakoncz`. Siatka owija komendy Studia z zewnątrz
zapory, mierzy SKUTEK (czy treść jest inna) i dopisuje ślad tylko wtedy, gdy
obsługiwacz go nie odłożył; zakres bierze z różnicy treści liczonej w znakach,
a nie z „całego dokumentu", bo podświetlenie całego pisma nie pokazuje niczego.

Wynik pomiaru po obu naprawach — `zmiany_sledzone(model)` / `czynnosci(model)`:
`studio.text.edit` 1/1, `studio.format.character.set` 1/1,
`studio.format.paragraph.set` 1/1, `studio.document.save` 1/1 (**było 0/0**),
`studio.markup.add` 0/1 plus wiersz znakowania (znakowanie nie jest zmianą
treści, więc zmiany śledzonej nie zakłada — i tak ma być).

**Zasada tożsamości — jedna, konsekwentna:** wykonawcą jest ten, kogo wskazuje
SZERSZY z dwóch sygnałów (fakt gniazda `transport.Narzedzia()` albo podpis
żądania). Operatorem czynność jest wtedy i tylko wtedy, gdy MILCZĄ OBA. Stempel
idzie w jedną stronę — podnosi do wykonawcy, nigdy odwrotnie; `author:
uzytkownik` z gniazda narzędzi jest twierdzeniem modelu o sobie, nie faktem.

## Zależność czynności wywiedziona z zakresu

Tabela `zaleznosc_czynnosci_studio` istnieje, jest czytana i **nikt jej nie
wypełnia** — żaden odcinek nie woła `ZapiszZaleznoscCzynnosci`. Odmowa
nazywająca zależność nie mogła się więc nigdy odezwać. Dołożyłem drugie źródło:
czynność późniejsza, która ruszyła **ten sam fragment**, stoi na wcześniejszej
z natury rzeczy. Oba źródła działają razem; wywiedzione myli się w stronę
odmowy nazwanej, a nie cichego psucia dokumentu.

## Sprawdziany — co mierzą i czym

Wszystkie sprawdziany skutku czytają bazę **osobnym połączeniem** i pytają
o świat, nie o kopertę. Ręka modelu bierze się z **gniazda serwera narzędzi**,
a pola `author` sprawdziany celowo NIE podają — mierzą to, czego model nie może
o sobie zataić.

`blokady_skutek_test.go` (7): odmowa nazywająca fragment, blokadę i jej powód;
zamiana w całym dokumencie oddająca bilans z nazwą blokady przy treści
zablokowanej nietkniętej i zmianie poza blokadą wniesionej; Operator bez
przeszkód; zasięg `everyone` wiążący także Operatora; blokada przechodząca
przez przywrócenie wersji **i dalej pilnująca**; zdjęcie i założenie wyłącznie
przez Operatora, także wbrew podpisowi `author: uzytkownik` z gniazda narzędzi.

`dziennik_skutek_test.go` (6): cofnięcie ze ŚRODKA dziennika zostawiające pracę
wcześniejszą i późniejszą; wpis w stanie `reverted`, nie usunięty; zależność
odmawiająca i nazywająca czynność stojącą na cofanej, przy dokumencie
nietkniętym; cofnięcie obu razem przechodzące; licznik zmian modelu nieliczący
Operatora; filtr per wykonawca (dwóch agentów to nie jeden); cofnięcie
wszystkiego **z zachowaniem pracy Operatora** wraz z kopią zapasową; cofnięcie
wybranych nieruszające nieodhaczonych; autor zapisany na drodze bez pola
`author`.

`autozapis_skutek_test.go` (4): **uczciwość zapisu mierzona na zapisie
NIEUDANYM** — awaria wywołana wyzwalaczem odrzucającym zapis postaci przy
zdrowym odczycie (zabranie tabeli mierzyłoby awarię odczytu, czyli co innego);
praca zostaje w kopii ze wskaźnikiem zmian niezapisanych, `ostatni_zapis_nieudany`
stoi; szereg autozapisu 3 wersje wobec 1 wersji Operatora, zawężenie w obie
strony; przywrócenie kopii DO NOWEGO dokumentu nietykające pierwotnego; wykaz
zgłaszający kopie niezapisane.

### Sprawdziany zmierzone mutacją (nie tylko zielone)

- wyłączenie `zaporaBlokadStudia` → **3 z 7** sprawdzianów blokad wywraca się.
  Pozostałe 4 trzymają inne drogi (odmowy w moich obsługiwaczach oraz
  sprawdzenie blokad w warstwie postaci) — to obrona w głąb, nie martwy kod.
- `Saved: false` → `Saved: true` w gałęzi niepowodzenia autozapisu →
  sprawdzian uczciwości zapisu wywraca się z nazwaniem szkody.

## Co OTWARTE — bez zaokrąglania w górę

1. **`StudioVersion` w kontrakcie nie ma pola `series`.** Odróżnialność wersji
   autozapisu wychodzi dziś zawężeniem wykazu i dwoma licznikami
   (`operatorCount`, `autosaveCount`), a nie polem w pozycji. Wykonawca odcinka 2
   budujący przełącznik „pokaż także zapisy samoczynne" musi zawężać wykaz,
   bo po pozycji szeregu nie rozpozna. **Pole do dołożenia przez odcinek
   kontraktu** — sprawdzian mierzy stan rzeczywisty, nie zamierzony.
2. **Postać dokumentu ma dwa miejsca zapisu.** Migracja 361 daje
   `postac_dokumentu_studio` i tam pisze `postacZapisz`; `wejsciePostacDokumentu`
   (odcinek wejścia) trzyma postać jako ładunek w `profil_wydania_studio` pod
   kluczem. Kopia zapasowa czyta tę drugą drogę. Dwa magazyny jednego bytu
   rozjadą się przy pierwszej poprawce — **do uzgodnienia między odcinkami**,
   nie do rozstrzygnięcia jednostronnie.
3. **Uzgodnienie treści z blokadą idzie WIERSZAMI, nie znakami.** Gdy strony
   mają tyle samo wierszy w środku (zamiana w całym dokumencie — sedno
   wymagania), parowanie jest dokładne. Gdy liczba wierszy się różni, środek
   jest jedną całością: styka się z blokadą → nie wchodzi i wraca pominięciem.
   Zamiana wstawiająca znak nowego wiersza w zablokowanym sąsiedztwie zostanie
   więc pominięta szerzej, niż musi. Zgadywanie parowania byłoby gorsze —
   wstawiłoby treść w środek zablokowanego cytatu i nazwało to bilansem.
4. **Zapora dla czynności na całym dokumencie uzgadnia PO wykonaniu.** Zakresu
   przed wykonaniem nie ma (powstaje z rachunku obsługiwacza), więc zapora
   zakłada kopię, puszcza czynność i przywraca fragmenty zablokowane. Blokada
   nie jest naruszona na zewnątrz — uzgodnienie zamyka się w tym samym
   wywołaniu, a odpowiedź niesie stan uzgodniony wraz z bilansem. Ale to jest
   „przed dotknięciem treści" w słabszym znaczeniu niż dla czynności na
   fragmencie, gdzie obsługiwacz nie rusza wcale.
5. **Kres domknięcia tożsamości leży poza tym odcinkiem.** `Rodzaj=narzedzia`
   jedzie parametrem nawiązania gniazda, a `transport/tozsamosc.go` stanowi, że
   tożsamość połączenia jest OPISEM i nie ma być zaporą. Zasada „szerszy
   wygrywa" stoi po właściwej stronie tego zastrzeżenia — opis może czynność
   wyłącznie zawęzić w prawach, nigdy rozszerzyć. Ale wykonawca, który
   kontrolowałby OBA sygnały naraz (zataił parametr gniazda i nie podpisał
   żądania), byłby dla rdzenia nierozpoznawalny. **Domknięcie wymaga wiązania
   tożsamości serwera narzędzi przy jego starcie, nie w module Studio** — to
   sprawa warstwy transportu i decyzji Właściciela.
6. **Wyróżnienie barwą (`studio.markup.add` o rodzaju `highlight`) nie wchodzi
   do postaci dokumentu.** Znakowanie żyje w `znakowanie_studio`; pole `form`
   w odpowiedzi zostaje puste. Wyróżnienie jest więc widoczne wykazem znakowań,
   a nie jako cecha postaci — czyli nie wyjdzie w wydaniu do `docx`.
7. **Siatka śladu autora nie zna ZAKRESU tak dobrze jak obsługiwacz.** Ślad
   dopisany siatką ma zakres różnicy treści i rodzaj `wstawienie` albo
   `usuniecie` — nie odróżni zmiany stylu nazwanego od dopisania akapitu tak
   dokładnie, jak zrobiłby to obsługiwacz, który wie, co robił. Siatka jest
   siatką: lepiej ślad przybliżony niż zmiana modelu niewidzialna. Właściwym
   domknięciem jest odkładanie śladu przez samą czynność —
   `studio.document.save` powinien to robić u siebie, a wtedy siatka przestanie
   się dla niego odzywać sama (sprawdza, czy ślad już jest).
8. **Odpowiedź odcinkowi 7: zapora NIE MOŻE dziś czytać
   `kontrolaWykonawcaZKontekstu` i nie jest to zaniedbanie.** Nośnik kontekstowy
   wypełnia `podpisWykonawcyStudia`, wołane na KOŃCU `zarejestrujStudio`; zapora
   owija rejestr PÓŹNIEJ, z `kompozycja.go`, więc w czasie wykonania stoi
   NA ZEWNĄTRZ podpisu (kolejność: siatka śladu → zapora → podpis →
   obsługiwacz). W chwili, w której zapora rozstrzyga, kontekst jest jeszcze
   pusty — dlatego czyta ładunek, a nie nośnik. Odwrócenie kolejności owinięć
   dałoby zaporze nośnik, ale odebrałoby jej możliwość PRZEPISANIA ładunku
   (uzgodnienie treści z blokadą) przed podpisem. Podpis w ładunku, którym pętla
   wykonawcza się posłużyła, jest więc drogą właściwą, nie obejściem.
9. **Zależności czynności nikt nie zapisuje.** Wywiedzenie z zakresu (wyżej)
   jest siatką, nie zamiennikiem: powiązania niewidoczne z zakresu (wstawienie
   tabeli jako podstawa scalenia komórki) wymagają, żeby czynność je odłożyła.
   `ZapiszZaleznoscCzynnosci` czeka na wołających w odcinkach postaci i wejścia.

# Studio, odcinek 5 — rdzeń: wejście, wyjście, szablony (18.08.2026)

## Liczby

- **20 komend wpiętych** — wszystkie rodziny odcinka mają obsługiwacza
  w `adapter_modul_studio_uchwyty.go`: `studio.document.create`,
  `.import.file`, `.import.pdf`, `.image.import`, `.save.as`, `.copy`,
  `.export.format`, `.export.batch`; `studio.template.save/delete/field.set/
  field.list/import/export/fill`; `studio.insert.from.library`,
  `.from.web`; `studio.provenance.list`; `studio.clipboard.copy/paste`.
- **10 plików rdzenia, 12 684 wiersze** + **2 pliki sprawdzianów, 1 407 wierszy**.
- `go build ./...` **zielony**, `(cd client && npx tsc --noEmit -p tsconfig.json)`
  **zielony**, zapora `TestRejestrPokrywaKomendyKontraktu` **pomyślna**.
- **10 sprawdzianów skutku, wszystkie pomyślne** (5 w `_ooxml_test.go`,
  5 w `_wejscie_skutek_test.go`); pełny przebieg pakietu `core` **zielony**.
  Przed nimi nie było ani jednego pliku `*ooxml*_test.go` — obieg `.docx`,
  rzecz wskazana przez Właściciela wprost, nie był mierzony niczym.

## Czym to policzone

`archive/zip` + `encoding/xml` na OOXML i ODF, `pdfcpu` na PDF,
`golang.org/x/net/html` na HTML, `golang.org/x/text/encoding` na strony kodowe.
**Ani jednej biblioteki pomocniczej do formatów biurowych** — `unioffice`
odpada na licencji handlowej, a LibreOffice i `pandoc` byłyby procesem potomnym
tam, gdzie biblioteka Go wystarcza. Rozbiór XML idzie przez drzewo węzłów
(`ooxmlWezel`), nie przez automat stanów: postać akapitu OOXML stoi w `w:pPr`
przed fragmentami, postać znaku w `w:rPr` wewnątrz każdego z nich, a drzewo
pozwala czytać je tam, gdzie są. Ten sam rozbiór obsługuje ODF.

## Trzy usterki, które wykryły sprawdziany — i nic innego by ich nie wykryło

1. **`studio.document.save` NIE UŻYWAŁ pola `form`.** Kontrakt je miał, dziura
   z niego została. To jest ta jedna rzecz, od której zaczęło się całe zlecenie:
   postać przy każdym zapisie ginęła, więc model był wobec dokumentu ślepy.
   Zamknięte: `wejsciePrzyjmijPostacZapisu` utrwala postać tą samą drogą, którą
   jedzie wniesienie pliku (`wejscieUtrwalPostac`) — nie drugą. Kolejność jest
   tu treścią: postać idzie PRZED zapisem treści, bo utrwalenie postaci składa
   treść z bloków i wpisuje ją do wiersza, więc treść z żądania musi wejść po
   nim. Brak pola znaczy „bez zmiany postaci", nie „postać na zero" — zwykły
   zapis treści nie ma prawa zetrzeć arkusza stylów ani tabel, i sprawdzian
   mierzy także to.
2. **Wyróżnienia tła wychodziły do pliku.** Rozstrzygnięcie Właściciela mówi,
   że znakowanie żyje w sesji, a nie w piśmie wysyłanym na zewnątrz. Komentarze
   i propozycje do pliku nie wchodziły nigdy (składacz bierze postać, nie
   warstwę adnotacji), ale wyróżnienie jest CECHĄ POSTACI i samo by nie odpadło.
   `wydanieZdejmijZnakowanie` zdejmuje je z kopii postaci wydania — dokument
   Operatora zostaje ze swoim znakowaniem nietkniętym — i NAZYWA to w bilansie.
3. **Rozbiór ODF liczył scalenie poziome dwa razy.** OpenDocument zapisuje je
   podwójnie: raz jako `table:number-columns-spanned` komórki scalającej, raz
   jako `table:covered-table-cell` na każdej przykrytej kolumnie. Tabela
   trzykolumnowa wychodziła jako czterokolumnowa, a kolumna czwarta bez
   szerokości — czyli w dokumencie niewidoczna. Liczone jest teraz licznikiem
   komórek do pominięcia, nie granicą kolumny: po komórce scalającej numer
   kolumny stoi już ZA scaleniem, więc granica nie odróżnia przykrycia
   poziomego od pionowego.

## Sprawdziany — co mierzą i czym

Materiał wejściowy `.docx` i `.odt` jest **wpisany wprost w sprawdzianie**, a nie
składany własnym składaczem. Gdyby powstawał `wejscieZlozOoxml`, pomiar
mówiłby, że rdzeń czyta to, co sam napisał — i przechodziłby także wtedy, gdy
oba końce mylą się w ten sam sposób, a plik jest dla Worda nieczytelny.

Porównanie idzie **po odczycie, nie po bajtach**: ten sam dokument da się
zapisać na wiele poprawnych sposobów, więc bajt w bajt nie zgodzi się nigdy
i nie ma się zgodzić.

| Sprawdzian | Miara niezależna |
|---|---|
| `TestDocxPoOdczycieZachowujeStylSekcjeITabele` | cudzy `.docx` → edytor → wydanie `.docx` → **drugi dokument** z wydanego pliku; ta sama miara po obu odczytach: nazwa stylu akapitu, styl obecny w arkuszu, wytłuszczenie fragmentu, orientacja pozioma, marginesy 30/15 mm (niesymetryczne celowo — równe znaczyłyby nastawę domyślną), tabela 2×3, wiersz nagłówkowy powtarzany, scalenie dwóch kolumn, komórka wchłonięta oznaczona, szerokości kolumn niezerowe; składniki archiwum wyliczone po nazwach; „dzieło" wraca z literą „ł" |
| `TestOdtPoOdczycieZachowujeStylSekcjeITabele` | to samo na OpenDocument; dodatkowo `mimetype` sprawdzony jako składnik **pierwszy i nieskompresowany** po nagłówkach archiwum, bo tym OpenDocument rozpoznaje swoje pliki |
| `TestKonwersjaPdfOddajeBilansOdzyskania` | liczba stron w bilansie zestawiona z `api.PageCount` na pliku; strony z warstwą i bez niej muszą się zsumować do całości; wykaz nieodzyskanego niepusty; zdanie bilansu musi mówić, że odzyskanie jest **odtworzeniem**; treść ze strony trzeciej obecna w dokumencie |
| `TestKonwersjaPdfSamychSkanowKierujeNaRozpoznanie` | PDF bez warstwy tekstowej: `needsTextRecognition`, zero stron z warstwą, **pozycja kolejki rozpoznania odnaleziona w `studio.ingest.queue.list`** — droga dalsza wskazana, nie domyślona |
| `TestKopiaDokumentuJestOsobnymBytem` | zmiana treści i zdjęcie wytłuszczenia **w kopii**; oryginał czytany osobnym wywołaniem musi być nietknięty w obu warstwach |
| `TestKopiaBezHistoriiIZHistoriaSaJawnymWyborem` | `versionsCopied` zestawione z `studio.repository.list` na kopii — liczba z odpowiedzi skonfrontowana z repozytorium |
| `TestZapisPrzenosiPostacDokumentu` | postać podana polem `form`, odczytana **ponownie** przez `studio.document.form.get`: A5, orientacja pozioma, margines 33 mm, sekcja, tabela 2×2, szerokości kolumn, wiersz nagłówkowy, styl akapitu, wytłuszczenie; drugi zapis **bez** pola `form` nie ma prawa zetrzeć tabeli ani arkusza |
| `TestWydanieUbozszeOddajeWykazCechPominietych` | `txt` i `md`: wykaz niepusty, każda pozycja z powodem, wykaz **nazywa** przypis oraz — w `md` — scalone komórki; treść dokumentu w pliku obecna |
| `TestWydaniePdfDajePlikOWlasciwejLiczbieStron` | strony **policzone `pdfcpu` w pliku**; 160 akapitów musi dać ≥ 2 strony, a dokument czterokrotnie krótszy **mniej stron** — bez tej drugiej miary sprawdzian przechodziłby przy stałej liczbie stron |
| `TestPlikWyjsciowyJestCzystyZWarstwyZnakowania` | trzy miary naraz: brzmienia komentarza nie ma **w bajtach** pliku; wyróżnienia nie ma **po odczycie** wydanego pliku; `includeComments` jest **nazwane** w wykazie cech pominiętych; treść pisma nietknięta |

## Formaty — co w pełni, co częściowo

**Wczytanie w pełni** (postać przechodzi do edytora):

- `.docx`, `.dotx` — styl akapitu i znaku, arkusz stylów z dziedziczeniem,
  sekcje wraz z nastawami strony, nagłówki i stopki wedle zasięgu (strony
  zwykłe, pierwsza, parzyste), tabele ze scaleniami, siatką i wierszem
  nagłówkowym, listy, obrazy jako obiekty, odsyłacze, tabulatory, obramowania.
- `.odt`, `.ott` — to samo, wraz ze stylami automatycznymi i układem strony
  z `styles.xml`.
- tekst czysty, markdown — markdown rozpoznaje nagłówki, cytat, listy, tabele
  z wyrównaniem kolumn, blok kodu, pogrubienie, kursywę i kod w linii.
- HTML — rozbiorem `x/net/html`, czyli tak, jak czyta przeglądarka; postać
  znaku **płynie w dół drzewa**, więc `<b><i>` daje jedno i drugie naraz.
- Strony kodowe: znacznik kolejności bajtów, deklaracja w treści, sprawdzenie
  UTF-8, a na końcu **miara rozkładu bajtów** między windows-1250, ISO-8859-2
  i windows-1252. Rozpoznana nazwa wychodzi w bilansie, bo przy tej ostatniej
  drodze rdzeń **zgadywał**, a nie czytał, i ma to powiedzieć.

**Wczytanie częściowo — i czego brakuje:**

- **`.rtf`** — akapity, pogrubienie, kursywa, podkreślenie, przekreślenie,
  indeksy, stopień pisma, wyrównanie, znaki `\'hh` i `\uN`. **Nie ma: tabel
  ani obrazów.** Tabela RTF jest ciągiem akapitów z granicami komórek zapisanymi
  w rozkazach składu, więc jej odzyskanie byłoby odtworzeniem układu, nie
  odczytem struktury — wchodzi treścią i **wychodzi w bilansie** jako układ
  nierozpoznany. Barwa znaku (`\cf`) nie wchodzi: idzie numerem w tablicy barw,
  której ten rozbiór nie czyta, więc pole zostaje puste zamiast nieść barwę
  zmyśloną.
- **`.pdf`** — tekst i akapity z warstwy tekstowej, liczba stron, strony bez
  warstwy, obrazy policzone. **Tabele nie są rozpoznawane strukturalnie** — PDF
  nie niesie struktury tabeli wprost i bilans mówi to wprost. PDF ze samych
  skanów **nie jest udawany jako skonwertowany**: idzie na `studio.ingest.recognize`
  wraz z pozycją kolejki.
- **Bajty obrazów z archiwum** nie są przenoszone do magazynu zasobów przy
  wczytaniu: obiekt powstaje wraz z nazwą składnika i wymiarami, bajty zostają
  w archiwum. To brak nazwany, nie cisza.

**Wydanie w pełni:** `.docx`, `.odt` — postać zachowana (obieg zmierzony).

**Wydanie częściowo — i czego brakuje:**

- `txt` — postać schodzi w całości; wykaz wymienia postać znaku i akapitu,
  arkusz stylów, nastawy strony, siatkę tabel, obiekty, pola i **aparat
  rozbity po rodzajach** (przypisy, spisy, odsyłacze, podpisy).
- `md` — style na znaczniki, tabele na tabele markdown; wykaz wymienia nastawy
  strony, własne postacie stylów i **scalone komórki z liczbą**.
- `html` — style nazwane na arkusz; nastawy nośnika wychodzą regułą druku, nie
  układem strony; obrazy znacznikiem bez adresu, bo bajty leżą w magazynie.
- `rtf` — **wykonalne i zrobione**, nie zgłoszone jako brak; postać znaku
  i akapitu wchodzi wpisana przy akapitach, bez tablicy stylów; obrazy tekstem
  zastępczym; aparat treścią.
- `pdf` — przez `dokumentPdfZTekstu`: **jeden krój, bez krojów dokumentu**;
  tabele wierszami tekstu, obrazy tekstem zastępczym; bez wskazanego profilu
  paginacja i stopka są domyślne. Naprawa dla pisma wzorcowego wskazana
  w bilansie wprost: wydać do `docx` albo `odt`.

## Co OTWARTE — bez zaokrąglania w górę

1. **Wykaz nośników leży w pliku modułu Design.** `wejscieWymiaryNosnika`
   **czyta** `nosnikiDruku` z `adapter_modul_design_druk_wspolne.go` — nie
   kopiuje go i nie edytuje tamtego pliku, wedle polecenia. Wykaz należy
   przenieść do miejsca wspólnego dla Studia i Designu **i uzupełnić o koperty
   C4 (229×324), C5 (162×229) i C6 (114×162)**, które Właściciel wymienia
   wprost, a wykaz ich nie ma. Do czasu tej zmiany nośnik kopertowy inny niż DL
   nie ma wymiarów i nastawy strony schodzą na format własny.
2. **Treść płaska dokumentu nie niesie brzmienia komórek tabeli.**
   `postacTekstFormy` (odcinek postaci) składa treść z akapitów; tabela jest
   strukturą i jej brzmienie stoi w postaci. Sprawdziany mierzą je w komórkach,
   nie w treści — ale wyszukiwanie pełnotekstowe po treści dokumentu tabel
   **nie znajdzie**. Rozstrzygnięcie należy do odcinka postaci, nie do tego.
3. **Obramowanie jest jedno na cztery krawędzie.** Kontrakt niesie jedną
   odmianę i jedną grubość wraz z przełącznikami krawędzi, a OOXML i ODF opisują
   każdą krawędź osobno. Cztery krawędzie o różnej grubości przechodzą przez
   rdzeń spłaszczone do pierwszej, która grubość ma.
4. **Scalenie pionowe wraca jako rozpiętość dwóch wierszy**, niezależnie od
   tego, ile wierszy naprawdę objęło: `w:vMerge` nie niesie liczby, tylko
   znacznik wznowienia i znaczniki ciągłości, a policzenie ich wymaga przejścia
   wszystkich wierszy tabeli wstecz. Komórki przykryte są oznaczone poprawnie,
   więc tabela nie rozjeżdża się o kolumnę — ale liczba w `rowSpan` jest
   najmniejszą prawdziwą, nie dokładną.
5. **Bilans wniesienia nie liczy tabel nierozpoznanych w PDF** (`tablesMissed`
   zostaje zerem). Policzenie ich wymagałoby rozpoznania układu tabelarycznego,
   czyli tej samej pracy, której ten rachunek nie robi. Zero jest tu brakiem
   wiedzy podanym jako zero i to jest jedyna pozycja bilansu, która obiecuje
   więcej, niż daje.

---

# Studio, odcinek 7 — pętla wykonawcza: rozkład zlecenia na zadania i jej okno

Rozstrzygnięcie Właściciela, które określa cały ten odcinek: **pętla wykonawcza
i praca kilku agentów naraz to NARZĘDZIA NIEURUCHAMIANE NA STARCIE.** Domyślnie
wyłączone, włączane jawnym, odwracalnym ustawieniem Operatora. Wszystko poniżej
jest tej zasadzie podporządkowane, a nie dopasowane do niej po fakcie.

## Co wpięte — pięć komend, nie sześć

W kontrakcie stoi **pięć** komend rodziny `studio.plan.*`, nie sześć:
`create`, `get`, `task.update`, `run`, `stop`. Zgłoszenie mówiło o sześciu —
pomiar mówi pięć i pomiar wygrywa. Wszystkie pięć ma obsługiwacza
(`zarejestrujPetleWykonawczaStudia`, wołane z `kompozycja.go` jednym wierszem po
`zarejestrujStudio`).

Do tego okno pętli **wołaniem, nie budową drugi raz**: `studio.agents.*` (nastawy,
obsada fragmentów, spięcia — odcinek 6), `studio.chain.save/list/run`,
`studio.batch.run` i `studio.operation.list` — te cztery stały w rdzeniu zbudowane
i **nie miały czym być uruchomione**. To okno jest tym czymś.

## Trzy rozstrzygnięcia nośne

**1. Pętla jedzie REJESTREM, nie wywołaniem adaptera wprost.** Krótsza droga —
`p.studio.OperacjaKontekstowa(...)` — ominęłaby oba owinięcia stojące w rejestrze:
podpis wykonawcy i **zaporę blokad fragmentu**. Pętla omijająca zaporę byłaby
cichą dziurą w blokadzie: Operator oznaczyłby podstawę prawną jako nietykalną,
a pętla przepisałaby ją, bo „to nie klient wołał". Model woła komendy tą samą
drogą co klient — więc pętla też.

**Usterka wykryta przy tym i naprawiona:** zapora rozpoznaje wykonawcę z faktu
gniazda albo z podpisu w ładunku. Pętlę puszcza **Operator** przyciskiem, a pracę
wykonuje **model** — bez podpisu zapora brała robotę pętli za robotę Operatora
i przepuszczała ją przez fragment zablokowany przed wykonawcą. Pętla podpisuje
więc swoje wywołanie wewnętrzne (`petlaPodpiszLadunek`) nazwami pól, które zapora
czyta z każdego ładunku Studia. Podpis nie jedzie drutem — ładunek jest wywołaniem
wewnętrznym rdzenia — więc nie jest to nowe pole kontraktu, a oświadczenie pętli
o sobie. Sprawdzian mierzy skutek na treści: zablokowany fragment po przebiegu
został **dosłownie** taki, jaki był.

**2. Nastaw nie zakładamy po raz drugi.** Pętla czyta nastawy wywołaniem
obsługiwacza `studio.agents.settings.get` **z rejestru komend** — nie własnym
odczytem konfiguracji po odgadniętych kluczach. Nazwy kluczy należą do tego, kto
nastawy wystawia; odgadnięcie ich tutaj założyłoby drugi magazyn pod pozorem
odczytu. Brak obsługiwacza nastaw znaczy dla pętli to samo co nastawa wyłączona —
z odmową nazywającą właśnie ten brak.

**3. Okno bierze się z DOKUMENTU, nie z żądania.** Rodzina `studio.plan.*` okna
nie niesie i nie musi: dokument Studia zna swoje okno (`DokumentStudia.Okno`),
więc pętla bierze je stąd — tą samą drogą, którą wzięłaby je operacja kontekstowa.
Zgłaszany wcześniej „brak `windowId` w kontrakcie" **nie jest brakiem**.

## Rozkład zlecenia jest rachunkiem, nie ozdobą

Zlecenie „napisz pismo w tej sprawie na podstawie tych materiałów, sprawdź
terminologię i przygotuj wersję do druku" daje **cztery zadania**, nie jedno.
Trzy drogi po kolei, pierwsza która da zadania wygrywa: zadania podane wprost →
rozkład modelem (kontrakt: „brak znaczy rozkład ułożony przez model") → rozkład
własny rdzenia po czynnościach nazwanych w zleceniu. **Trzecia droga nie zawodzi
nigdy**: zlecenie nierozpoznane daje jedno zadanie niosące je w całości, bo
rozkład pusty byłby odpowiedzią udaną bez treści. Rozpoznanie idzie po
POŁOŻENIU słowa w zdaniu, nie po kolejności wykazu — kolejność czynności jest ta,
którą podał Operator. Zależności są łańcuchem: nie da się poprawić języka pisma,
którego jeszcze nie ma.

## Czym pętla wykonuje zadanie

- treść i postać → `studio.contextual.op`: wynik wchodzi jako **zmiana śledzona
  autora `model`**, więc praca pętli jest widoczna w podświetleniu zmian modelu
  i cofa się tak samo jak każda inna;
- `export` → `studio.document.export.format`, bo wydania nie wolno udawać
  operacją na treści;
- rodzina, której rdzeń nie wystawia → zadanie kończy się stanem `failed`
  z powodem **nazywającym** brak. Nigdy „gotowe" bez pracy.

Zatrzymanie ZATRZYMUJE: znacznik sprawdzany przed każdym zadaniem i po każdym
obiegu, zadanie w biegu wraca do czekania **wraz z powodem mówiącym, że zostało
PRZERWANE**. Nie znika i nie udaje domkniętego; to, co zrobiono, zostaje.

## Okno: powierzchnia należy do dokumentu

Opracowanie stawia Execution Loop Window jako stałą kolumnę obok okna rozmowy.
Rozstrzygnięcie Właściciela jest mocniejsze i ono obowiązuje: okno wchodzi
**nakładką na żądanie**, znacznikiem przebiegu `[⇄ pętla]`, i schodzi po zwinięciu.
Stała kolumna została — **wyłącznie jako tryb do wyboru** Operatora. Samoczynne
otwarcie przy zleceniu wielozadaniowym zostaje, bo wtedy okno jest naprawdę
używane. Zlecenie jednozadaniowe okna nie otwiera.

Cztery widoki, jedna kolejka: zadania ze stanami i wykonawcą, obsada wykonawców
na fragmentach, warsztat łańcucha, tryb wsadowy. Krok łańcucha i dokument wsadu
są **zadaniami tej samej kolejki** — kontrakt mówi wprost, że przebieg łańcucha
prowadzi pętla wykonawcza okna, więc druga maszyneria obok byłaby drugą prawdą.

**Odłożone brzmienie nie przepada.** `StudioAgentConflict.deferredText` stoi
w oknie wprost do przeczytania, z czynnością „Przyjmij brzmienie". Gdy rdzeń
odłożył je jako zmianę śledzoną albo propozycję na marginesie, widok **odsyła do
niej po identyfikatorze** zamiast trzymać drugą kopię tekstu.

**Wsad idzie dokument po dokumencie**, a nie jednym wywołaniem na czterdzieści.
Rachunek ten sam, komenda ta sama — różnica w tym, że po każdym dokumencie pętla
ma chwilę, w której może usłyszeć „przerwij". Jedno wywołanie na całość nie ma
gdzie się zatrzymać. Zatrzymanie nie wycofuje niczego.

## Sprawdziany skutku — siedem, mierzą dokument i stan, nie kopertę

1. **pętla wyłączona nastawą** — `started: false`, powód nazywa
   `executionLoopEnabled` ORAZ mówi, czym ją włączyć; **i kolejka jest nietknięta**
   (odczyt osobnym `studio.plan.get`). Dwie miary, bo pojedyncza dałaby się obejść;
2. **plan uruchomiony wykonuje albo nazywa brak** — żadne zadanie nie zostaje
   w czekaniu, każde nieudane ma powód, bilans niesie liczby;
3. **wydanie jako zadanie pętli** — zadanie `export` kończy `done`, bilans je liczy;
4. **zatrzymanie** — rozkład `stopped` z powodem Operatora, zadanie domknięte
   przed zatrzymaniem **zostaje** domknięte, żadne zadanie nie znika;
5. **blokada widoczna w kolejce** — zadanie nie kończy „gotowe", powód przy
   zadaniu nazywa blokadę, bilans ją nazywa, **a treść zablokowana została
   dosłownie** (miara ostateczna — reszta mówi o meldunkach, ta o dokumencie);
6. **wsad rozlicza KAŻDY dokument** — suma przyjętych i odrzuconych zgadza się
   z żądaniem, każde odrzucenie ma powód, dokument nieistniejący ma powód **inny**
   niż istniejące (jeden powód na wszystko byłby bilansem pozornym);
7. **rozkład rozpoznaje czynności** — rodzaje `research`, `draft`, `proofread`,
   `format` z jednego zdania, plus łańcuch zależności; oraz przejścia stanów
   zadania z odmową nazywającą oba stany.

## Zmierzone

`go build ./...` czysto. `(cd client && npx tsc --noEmit -p tsconfig.json)` czysto.
`go test ./server/internal/core/` — **ok, 508 s**, wraz z zaporą
`TestRejestrPokrywaKomendyKontraktu`. Klient: **40 plików, 508 sprawdzianów,
wszystkie przechodzą**. Siedem nowych sprawdzianów skutku pętli przechodzi.

## Co OTWARTE — bez zaokrąglania w górę

- **Rozkłady leżą w magazynie pamięciowym procesu, nie w bazie.** Warstwa danych
  nie ma tabeli rozkładu ani zadania, a zakładanie jej należy do odcinka
  fundamentu. Skutek nazwany, nie przemilczany: **rozkład nie przeżywa ponownego
  uruchomienia rdzenia.** Potrzebne tabele `plan_studio` i `zadanie_studio`.
- **Zadania treści wymagają kanału modelu okna.** Bez kanału kończą się `failed`
  z powodem nazywającym brak — uczciwie, ale to nie jest praca wykonana. Na
  stanowisku sprawdzianu kanału nie ma, więc miarą wykonania jest tam zadanie
  `export`.
- **Brak komendy „przyjmij odłożone brzmienie".** Okno wnosi je drogą rozkładu
  i pętli, żeby przyjęcie miało ślad — ale komendy własnej na to w kontrakcie nie
  ma. To jest brak wypisany, nie dziura zaklejona wywołaniem, którego nie ma.
- **Znacznik przebiegu stoi we własnym pasie modułu, nie w pasku kontekstu okna
  rozmowy** — opracowanie stawia go tam, a pliki paska kontekstu należą do
  odcinków 2 i 3. Przeniesienie to jeden wiersz po ich stronie.
- **`studio.chain.run` w rdzeniu nie wykonuje kroków** — waliduje i oddaje `runId`
  wraz z liczbą kroków. Zgodnie z kontraktem („przebieg prowadzi pętla wykonawcza
  okna") przebieg prowadzi okno, więc funkcja działa; ale gdyby ktoś zawołał
  `studio.chain.run` bez okna, nie stanie się nic i **nikt tego nie zgłosi**.
- **Zapora blokad nie czyta `kontrolaWykonawcaZKontekstu`**, choć odcinek 6 ten
  nośnik zbudował. Pętla obeszła to podpisem w ładunku; czystszą naprawą byłoby,
  żeby zapora sięgnęła po nośnik kontekstowy — to jest zmiana w pliku odcinka 6.
- **`PODZIAL-PRAC-STUDIO.md` i `ZLECENIE-STUDIO-MODEL-DOKUMENTU.md` już nie
  istnieją** w drzewie, więc dopisania odcinka do podziału prac nie było gdzie
  wykonać. Ten rozdział jest jedynym zapisem odcinka 7.

# Odcinek 4 — Rdzeń: postać dokumentu

Zakres wzięty w całości: styl znaku i akapitu, style nazwane wraz
z dziedziczeniem, nastawy strony i sekcje, nagłówki i stopki, formaty nośnika
i koperty, listy, tabele, obiekty, symbole, pola, aparat dokumentu, widok
i linijka.

## Wpięte — 59 komend, wszystkie z obsługą

`studio.document.form.get/save`, `studio.text.get/edit`,
`studio.format.character.set/get`, `studio.format.paragraph.set/get`,
`studio.format.clear`, `studio.format.case.set`,
`studio.format.painter.copy/apply`, `studio.format.similar.select`,
`studio.format.replace`, `studio.style.list/save/apply/delete`,
`studio.page.setup.get/set`, `studio.page.paper.list`,
`studio.page.envelope.set`, `studio.page.break.insert`,
`studio.page.headerfooter.set/get`, `studio.page.numbering.set`,
`studio.page.watermark.set`, `studio.section.save/list/delete`,
`studio.list.apply`, `studio.list.bullet.set`, `studio.list.numbering.set`,
`studio.list.restart`, `studio.list.level.indent`,
`studio.table.insert/structure.edit/format.set/sort/convert/list`,
`studio.object.insert/format.set/list/remove`, `studio.symbol.insert/list`,
`studio.symbol.autoreplace.set/list`, `studio.field.insert/refresh/list`,
`studio.apparatus.insert/list/remove/refresh`, `studio.view.get/set`,
`studio.ruler.tabstop.set`.

Każda jedzie tą samą drogą co narzędzie modelu — rejestrem, nie osobnym
wejściem — więc wpięcie obsługi jest równocześnie wpięciem narzędzia.

## Dwie usterki, które wykryły dopiero sprawdziany skutku

Obie były niewidoczne w odpowiedziach komend i obie kłamałyby Operatorowi.

1. **Styl AKAPITU nie niósł postaci znaku.** „Nagłówek poziomu 1" jest stylem
   akapitu, a mówi o stopniu pisma i pogrubieniu. Postać skuteczna liczyła się
   wyłącznie ze stylu ZNAKU, więc akapit ze stylem nagłówkowym dostawał stopień
   tekstu zasadniczego, a zmiana stylu nagłówkowego nie ruszała ani jednej
   litery. Naprawione trzema warstwami w `postacZnakSkutecznyWBloku`: styl
   akapitu, styl znaku, postać własna fragmentu — warstwa bliższa fragmentowi ma
   pierwszeństwo. Bez sprawdzianu numer 1 ta usterka przeszłaby cało: odpowiedź
   `style.save` meldowała skutek, którego nie było.

2. **Zmiana formatu nośnika zgłaszała, ale nie przeliczała.** Bilans nazywał
   tabelę szerszą niż nowy nośnik i na tym poprzestawał — dokument zostawał
   w stanie, którego nie da się wydrukować. Wymaganie Właściciela jest
   dwuczłonowe („przelicza układ, nie obcina treści; co się nie zmieściło, wraca
   bilansem"), więc `stronaBilansUkladu` teraz PRZELICZA: kolumny tabeli schodzą
   w tej samej proporcji, w jakiej stały, obraz z zachowaniem proporcji boków —
   i każde przeliczenie wchodzi do bilansu.

## Sprawdziany skutku — 85 pomyślnych, ani jednego niepomyślnego

Miara nigdy z odpowiedzi czynności: odpowiedź pisze ten sam kod, który zmieniał
dokument, więc potwierdzałaby samą siebie. Dwie drogi niezależne — DRUGIE
POŁĄCZENIE do pliku bazy własnym SQL oraz OSOBNE wywołanie komendy odczytu.

| Wymaganie Właściciela | Sprawdzian | Wynik |
|---|---|---|
| zmiana stylu nazwanego przestawia WSZYSTKIE miejsca użycia | `TestPostacStyluPrzestawiaWszystkieMiejscaUzycia` | pomyślny (wykrył usterkę 1) |
| postać na fragmencie zmienia WYŁĄCZNIE ten fragment | `TestPostacFragmentuZmieniaWylacznieFragment` | pomyślny |
| spis treści po odświeżeniu zgadza się z nagłówkami | `TestSpisTresciZgadzaSieZNaglowkami` | pomyślny |
| przypis przenumerowuje się po wstawieniu przypisu przed nim | `TestPrzypisPrzenumerowujeSiePoWstawieniuPrzedNim` | pomyślny |
| tabela po scaleniu ma policzone szerokości, nie zerowe | `TestTabelaPoScaleniuMaPoliczoneSzerokosci` | pomyślny |
| zmiana nośnika przelicza układ i oddaje bilans | `TestPostacZmianaNosnikaPrzeliczaUkladIOddajeBilans` | pomyślny (wykrył usterkę 2) |
| narzędzie modelu odkłada zmianę autora `model` | `TestPostacNarzedziaModeluOdkladaAutoraModel` | pomyślny |
| sortowanie tabeli zna pismo polskie | `TestSortowanieZnaPismoPolskie`, `TestSortowanieTabeliUkladaNazwiskaPoPolsku` | pomyślny |

Sprawdzian stylu mierzy jeszcze jedną rzecz, której nie ma na wykazie
Właściciela, a bez której punkt pierwszy byłby pozorny: zagląda do drzewa
w bazie i pilnuje, że fragmenty NIE niosą stopnia własnego. Gdyby stosowanie
stylu kopiowało jego postać na fragmenty, wszystkie miejsca użycia zmieniłyby
się przy zapisie stylu raz — i nigdy więcej.

Sprawdzian nośnika mierzy też, że pole `widthMm` wolno mieć puste dla nośnika
NAZWANEGO — tak stanowi kontrakt tego pola — a wymiar A5 czyta z wykazu
nośników, nie z drugiego wykazu Studia.

## Co OTWARTE — bez zaokrąglania w górę

1. **Kształt i ikona nie wołają modułu Design.** `studio.object.insert` zakłada
   obiekt we własnym rachunku Studia (rodzaj kształtu, geometria, wypełnienie,
   obrys, obrót, opływanie, tekst wewnątrz) i zapisuje `designNodeId`, gdy
   wołający go podał — ale sam `design.vector.shape.add` ani
   `design.icon.library.search` nie są wołane. Powód: `adapterStudia` nie ma
   portu Design, a jego dodanie wymaga zmiany konstruktora w
   `adapter_modul_studio.go` i w `kompozycja.go` — plikach, których ten odcinek
   nie posiada. Do wpięcia potrzebna jedna zależność w konstruktorze.
2. **`widthMm` i `heightMm` nastaw strony zostają puste dla nośnika nazwanego.**
   Zgodne z kontraktem, ale okno musi wtedy wołać `studio.page.paper.list`, żeby
   narysować kartkę. Wypełnianie ich zawsze byłoby wygodniejsze dla klienta
   i warto to rozważyć przy następnej poprawce kontraktu — nie robię tego sam,
   bo kontrakt należy do odcinka fundamentu.
3. **Cykl w arkuszu stylów przejętym z szablonu kończy się urwaniem łańcucha
   dziedziczenia po dziesięciu poziomach**, bez zgłoszenia Operatorowi. Kres
   jest konieczny (bez niego rdzeń wisi na jednym wywołaniu), ale milczenie
   o cyklu jest długiem: Operator dostaje styl liczony inaczej, niż zapisał,
   i nie wie dlaczego.
