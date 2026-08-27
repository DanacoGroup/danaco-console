# Rejestr terenów

Żywy wykaz terenów. Prowadzi go wyłącznie Prowadzący budowę, na gałęzi `main`.
Teren bez wpisu w tym rejestrze nie jest otwarty, a praca na nim nie zostanie
przyjęta. Zasady podziału opisuje [ustrój budowy](ustroj-budowy.md).

## Tereny otwarte

### okno-przygotowania

Ostatni ekran drogi wejscia ma **slepy zaulek** i pasek, ktory nigdy nie dobiegnie
konca. Zmierzone przez sesje pomiarowa w przegladarce wobec zywego rdzenia.

| | |
|---|---|
| **Galaz** | `teren/okno-przygotowania` z `main` |
| **Wykaz plikow** | `budowa/klient/src/wejscie/ekrany/przygotowanie.ts`, `budowa/klient/src/wejscie/tresci.ts`, `budowa/klient/src/wejscie/przebieg.ts`, sprawdziany w `budowa/klient/src/wejscie/` |
| **Poza terenem** | `budowa/server/`, `budowa/shared/`, `budowa/desktop/`, `design/`, `prowadzenie/`, `budowa/klient/src/polaczenie`, `budowa/klient/src/protokol` |

**Pierwsza — przycisk prowadzi donikad.** „Pomin przywracanie sesji"
(`przygotowanie.ts:64`) ma `komunikat: 'pominiecie'` i **nie ma przypisanej
czynnosci**. Po klikniecu wystawia zdanie i nic wiecej sie nie dzieje; wykaz
etapow i pasek zostaja bez zmiany. To **jedyne wyjscie** z ostatniego ekranu poza
wylogowaniem.

**Druga — pasek staje na 40% na zawsze.** Trzy z pieciu etapow przygotowania nie
maja komendy w kontrakcie i zostaja w stanie oczekiwania. Zachowanie jest
udokumentowane w kodzie jako zamierzone, ale odbior jest taki, jakby aplikacja
sie zawiesila.

**Czego NIE robisz.** Nie dokladasz komend do kontraktu — kontrakt nalezy w tej
turze do innego terenu. Okno ma stac sie **uczciwe wobec tego, co potrafi**, a nie
udawac postep, ktorego nie ma.

**Kryteria odbioru.**

1. Przycisk pominiecia **prowadzi dalej** albo znika. Jesli prowadzi - wykazane
   przejsciem w przegladarce, z przytoczona odslona przed i po.
2. Ostatni ekran **nie zostawia Operatora bez wyjscia**: z kazdej odslony jest
   droga naprzod albo nazwany powod, dla ktorego jej nie ma.
3. Postep nie udaje. Albo dobiega konca, albo mowi wprost, na co czeka - decyzja
   Twoja, uzasadniona w raporcie.
4. Tekst widoczny dla uzytkownika **wylacznie** w katalogu tresci - sprawdzian
   `katalog-tresci.test.ts` przechodzi.
5. `tsc --noEmit` bez bledu; sprawdziany klienta zdane (26 przebiegu, 7 katalogu).
6. Wykazane **przejsciem w przegladarce** wobec zywego rdzenia, ze droga wejscia
   dalej przechodzi od konca do konca. Zrzuty ekranu w katalogu tymczasowym.

### witryna-pobierania

Wykaz wydan sprowadzono do dwoch postaci hybrydowych, tresc strony wokol niego
dalej obiecuje szesc.

| | |
|---|---|
| **Galaz** | `teren/witryna-pobierania` z `main` |
| **Wykaz plikow** | `budowa/witryna/tresc/` |
| **Poza terenem** | `budowa/witryna/wydania.json`, `budowa/server/`, `budowa/klient/`, `budowa/desktop/`, `budowa/shared/`, `design/`, `prowadzenie/` |

`strony.mjs:114,116,292,295,339` i `pobierz.mjs:331,427,435,445,475` niosa
rozdzial o wyborze miedzy hybryda a postacia natywna, pakiety Linuksa, instrukcje
sumy kontrolnej dla AppImage oraz zdanie o braku instalki natywnej dla Windows.
Wszystkie te postaci **zniosla pozycja 8** rejestru decyzji - przeczytaj ja.

**Kryteria odbioru.**

1. Zero wzmianek o postaci natywnej, pakietach Linuksa i AppImage - wykazane
   przeszukaniem **z sonda dodatnia** dowodzaca, ze wzorzec cokolwiek lapie.
2. Tresc strony zgadza sie z `wydania.json` co do liczby i nazw postaci -
   wykazane zestawieniem obu.
3. Instrukcja sumy kontrolnej dotyczy postaci, ktore naprawde sa wystawione.
4. Zdanie o tym, czego jeszcze nie ma, mowi prawde wobec pozycji 8 albo znika.
5. Jesli witryna daje sie zbudowac albo obejrzec - wykazane uruchomieniem.
   Jesli nie daje - napisane wprost, czego brakuje.


### zdolnosc-wyszukiwania

**Buduje dwie nowe zdolności produktu** na modelach, które stoją odłogiem.
Umocowanie: pozycja 17 rejestru decyzji. **Jedyny teren tej tury z prawem zmiany
kontraktu.**

| | |
|---|---|
| **Gałąź** | `teren/zdolnosc-wyszukiwania` z `main` |
| **Wykaz plików** | `budowa/shared/contract.json`, `budowa/server/internal/wiedza/`, `budowa/server/internal/core/adapter_modul_wiedza*.go` oraz adaptery obszaru `knowledge`, sprawdziany tych pakietów |
| **Poza terenem** | `budowa/klient/`, `budowa/desktop/`, `design/`, `prowadzenie/`, `internal/core/adapter_rozmowa_*`, `adapter_przejecie_sterowania.go`, `zgodnosc_kontraktu_test.go`, `blokady_skutek_test.go` (teren `usterki-rdzenia`), `internal/store/` migracje nastaw (teren `nastawy-wdrozenia`) |

**Przedmiot pierwszy — przesiew wyników.** `knowledge.search` dostaje pole
`rerank` wraz z liczbą kandydatów. Wyszukanie robi wtedy dwa przebiegi: kosinus
wektorów zbiera kandydatów, krzyżowy koder układa je ponownie. Model:
`/opt/danaco-modele/reranker` (bge-reranker-v2-m3, 2,2 GB, `safetensors`).

**Przedmiot drugi — oś obrazu.** Nowa komenda wyszukania obrazu zdaniem. Model:
`/opt/danaco-modele/clip` (CLIP ViT-L/14, 1,6 GB, `safetensors`).

**Zmiana kontraktu obwarowana.** Wolno **dołożyć**; nie wolno zmienić ani usunąć
niczego istniejącego — żadnej komendy, pola, wartości wyliczenia ani opisu.
Suma zastana: `b7d0436880d878e78576`.

**Kryteria odbioru.**

1. `knowledge.search` z przesiewem oddaje wynik **inaczej uszeregowany** niż bez
   przesiewu, na tej samej treści i tym samym zapytaniu — z przytoczonymi obiema
   odpowiedziami. Sam fakt odpowiedzi nie jest wykazaniem; przesiew ma **coś
   zmienić** i masz pokazać co.
2. Wyszukanie obrazu zdaniem oddaje **trafienie**, nie pustkę — z przytoczonym
   żądaniem, odpowiedzią i wskazaniem, który obraz wrócił i dlaczego.
3. Wagi z `/opt/danaco-modele` są **użyte** — wykazane tak, żeby widać było brak
   pobrania drugiej kopii; sonda dodatnia ma dowodzić, że twój instrument pobranie
   w ogóle wykryje.
4. Brak modelu albo biblioteki daje odmowę **nazywającą brak i drogę naprawy**,
   nie błąd wewnętrzny — wykazane sprawdzianem.
5. Kontrakt ruszony **wyłącznie dodaniami** — wykazane porównaniem z sumą zastaną,
   z wykazem tego, co przybyło. Generator daje wynik bajtowo powtarzalny w dwóch
   przebiegach; klient przechodzi `tsc --noEmit`.
6. Nowe komendy przechodzą **bramę kontraktu**: treść niepełna dostaje odmowę
   nazywającą brakujące pola, treść pełna przechodzi — obie przytoczone.
7. `gotestsum -- -count=1 ./...` — zero niepowodzeń wobec 2106 zdanych,
   17 pominiętych, zero niezdanych.

### odtwarzanie-twarzy

**Buduje zdolność, na którą kontrakt już czeka.** `image.upscale` ma pole
`faces: bool`; rdzeń odmawia dziś, bo silnika nie ma. Zmiana kontraktu
**niepotrzebna**.

| | |
|---|---|
| **Gałąź** | `teren/odtwarzanie-twarzy` z `main` |
| **Wykaz plików** | `budowa/server/internal/core/adapter_narzedzia_obraz_model_silniki.go` oraz pozostałe `adapter_narzedzia_obraz_*.go`, `zaleznosci_zewnetrzne.go`, sprawdziany tych pakietów |
| **Poza terenem** | `budowa/shared/`, `budowa/klient/`, `budowa/desktop/`, `design/`, `prowadzenie/`, `internal/wiedza/`, **cały obszar `design.*`** (zapora fotografii) |

**Stan zmierzony.** `/opt/danaco-modele/twarze` niesie `GFPGANv1.4.pth`
i `codeformer.pth`. Rdzeń szuka programu `gfpgan-ncnn-vulkan`, którego **na
maszynie nie ma**; stoi `realesrgan-ncnn-vulkan`, ale jego wydanie sieci
twarzowej nie niesie. Wagi `.pth` żądają stosu, którego rdzeń nie woła.

**Instalowanie dozwolone w tym terenie**, na tych samych zasadach co
`zaplecze-modeli`: wyjątek nazwany, wygasa z terenem, wszystko postawione trafia
do raportu wraz z wagą na dysku. Wybór drogi — wydanie `ncnn` niosące sieć
twarzową albo pomocnik pythonowy na stojących wagach — jest Twój i **uzasadniasz
go w raporcie**.

**Kryteria odbioru.**

1. `image.upscale` z `faces: true` oddaje obraz **różny** od tego samego
   powiększenia z `faces: false` — z przytoczonymi obiema odpowiedziami i miarą
   różnicy. Sam brak odmowy nie jest wykazaniem.
2. Wynik jest **prawdziwym obrazem** — format i rozmiar przytoczone.
3. Brak silnika dalej daje odmowę nazywającą brak i drogę naprawy.
4. Wagi z `/opt/danaco-modele/twarze` użyte albo **wprost napisane, dlaczego
   rdzeń sięga po inne**.
5. Zapora fotografii przechodzi; `gotestsum -- -count=1 ./...` — zero
   niepowodzeń. Kontrakt nietknięty — wykazane sumą.
6. Wykaz wszystkiego, co postawione na maszynie, wraz z wagą.

## Zgłoszenia oczekujące na teren

### Cala rodzina komend `control.*` nie jest zmontowana

`nowyAdapterPrzejeciaSterowania` (`internal/core/adapter_przejecie_sterowania.go`)
nie ma w drzewie **ani jednego wywolania** - przeszukanie po nazwie dalo sama
definicje, przy sondzie dodatniej na `nowyAdapterRozmowy`, ktory trafia
w `montaz_rozmowa.go:53`.

**Skutek:** `control.takeover`, `control.release` i stan sterujacego wchodza do
rejestru komend droga `zarejestrujPrzejecieSterowaniaNiewpiete`, czyli **jako
odmowa**. Rejestr steru nigdy nie dostaje wpisu. To wyjasnia, dlaczego funkcja
`sterZlecenia` byla martwa, i dlaczego jej dokumentacja opisywala uklad, ktorego
kontrakt nie zna.

**Montaz wymaga pol, ktorych kontrakt nie ma.** `shared.LoopState` nie niesie
`controller`, `takenOverAt` ani `takenOverBy` - zmierzone odczytem
`shared/contract.go:13104-13123`. Wpiecie rodziny jest wiec rozstrzygnieciem
o zakresie produktu, nie praca inzynierska: wymaga pozycji rejestru decyzji
i dolozenia pol do kontraktu.

**Druga martwa funkcja w tym samym pliku:** `zachowajStery` wraz z metoda
`rejestrSteru.zachowaj`. Zostawiona swiadomie - to mechanizm sprzatania po
zamknietych oknach, wzorowany na czynnym `rejestrBiegow.Zachowaj`; usuniecie
zabraloby zabezpieczenie w chwili, gdy adapter zostanie wpiety.


### Wagi mowy stoja poza katalogiem modeli

Teren `nastawy-wdrozenia` slusznie **nie zalozyl** trzeciej nastawy. Zgloszenie
podawalo `/opt/danaco-modele/mowa`, a takiego katalogu na maszynie nie ma -
sprawdzone przez wykonawce i potwierdzone przez Prowadzacego. **641 MB wag mowy**,
w tym `models--Systran--faster-whisper-small`, stoi w `~/.cache/huggingface`.

Nastawa wskazujaca katalog nieistniejacy sprowadzilaby silnik mowy do pobierania
od nowa, czyli **pogorszyla stan**. Domkniecie wymaga terenu obejmujacego zarazem
przeniesienie wag na dysku i `internal/mowa/ustawienia.go`, ktory lezal poza
wykazem tamtego terenu.

### Sprawdziany rdzenia sa na granicy domyslnego limitu czasu

Pakiet `internal/core` przekracza domyslne 10 minut `go test`: bieg bez
`-timeout` konczy sie zrzutem gorutyn po 600 s, nie wynikiem. Zmierzone:
`FAIL danacoconsole/server/internal/core 600.053s` przy 978 sprawdzianach.

**Skutek:** kto uruchomi sprawdziany dokladnie tak, jak podaje przekazanie,
dostanie niepowodzenie zamiast pomiaru. Obowiazujaca postac polecenia to odtad
`gotestsum -- -count=1 -timeout 40m ./...`. Rozbicie pakietu `core` albo
skrocenie najwolniejszych sprawdzianow jest osobna praca.

### Sciezka wag na wdrozeniu - obawa rozstrzygnieta pozycja 8

Wykonawca zglosil, ze `/opt/danaco-modele/embedder` jako wartosc domyslna jest
wlasnoscia tej maszyny, a jedyna postacia produktu jest hybryda Windows 11.

**Obawa znika.** Pozycja 8 stanowi, ze w hybrydzie **rdzen stoi na serwerze
Danaco**, a u Operatora staje samo okno. Rdzen nigdy nie biegnie na Windowsie,
wiec sciezka linuksowa jest sciezka serwera wdrozenia i jest poprawna. Nastawa
pozostaje bez rozgalezienia po systemie.


### Strona pobierania obiecuje warianty zniesione pozycja 8

Wykaz wydan sprowadzono do dwoch postaci hybrydowych, tresci strony wokol niego
nie. `budowa/witryna/tresc/strony.mjs:114,116,292,295,339` oraz
`tresc/pobierz.mjs:331,427,435,445,475` niosa rozdzial o wyborze miedzy hybryda
a postacia natywna, pakiety Linuksa, instrukcje sumy kontrolnej dla AppImage
i zdanie o braku instalki natywnej dla Windows. Dane wystawiaja dwie postaci,
a strona obiecuje szesc.

### Instalka wychodzi bez licencji i instrukcji

Skrypty instalek wymagaly z korzenia czterech dokumentow produktu i pakowaly je
do zasobow instalki. W korzeniu stoi dzis sam `CLAUDE.md`, wiec krok padal na
pierwszym pliku i teren `powloka-i-wydanie` go zdjal, zeby skrypty w ogole biegly.
Skutek: produkt koncowy nie niesie warunkow licencji ani instrukcji. Skrypty
wydania 1.0 przerywaly budowe przy ich braku, nazywajac to wprost: instalator bez
licencji i instrukcji nie jest produktem koncowym. Rozstrzygniecie nalezy do
Wlasciciela.

### Droga aktualizacji dla Linuksa zyje w powloce

`budowa/desktop/src-tauri/src/aktualizacja/droga.rs` niesie wariant AppImage,
sprawdzenie naglowka ELF i piec sprawdzianow dla platformy, ktorej pozycja 8 nie
przewiduje. Zostawione swiadomie: usuniecie zmienialo liczbe sprawdzianow, ktora
kryterium odbioru bralo za odniesienie.

### Katalog pakietu interfejsu ma dwie nazwy w drodze wdrozenia

`server/internal/konfiguracja/katalog_klienta.go` liczy domyslnie na `klient/dist`
i wprost ostrzega, ze dwie nazwy jednego katalogu wracaja odmowa przy pierwszym
uruchomieniu. Jednostka systemd w `budowa/packaging/` klade tymczasem sciezke
`client/dist`. Dziala wylacznie dlatego, ze jednostka podaje ja jawnie.
Ujednolicenie wymaga terenu obejmujacego `budowa/packaging/`.


Ustalenia z zamkniętych i biegnących terenów, które wykraczają poza ich zakres.
Każde zgłoszenie ma wskazany plik i wiersz. Zgłoszenie staje się terenem, gdy
Prowadzący je otworzy; do tego czasu jest wykazem, nie pracą.

### Ukończenie biegu może raz na jakiś czas skłamać — domknięcie poza terenem

Ustalenie terenu `warunek-ukonczenia-zadania`, wyprowadzone z odczytu dwóch
funkcji. `powodTury` (`adapter_rozmowa_petla.go:46`) nie widzi `zamkniecie.Blad`,
więc tura zamknięta zdarzeniem `result` z `is_error: true` przy sprawnym kanale
idzie do pętli jako `PowodWynik` i pętla ogłosi `completed`, choć wiadomość
dostaje stan `error`. Pętla używa jedynego sygnału, który dostaje, i nie dubluje
odczytu zamknięcia — drugi czytelnik byłby drugą prawdą. Domknięcie: jeden wiersz
w `adapter_rozmowa_petla.go` (przekazanie `zamkniecie` do `powodTury`) — plik
poza terenem, więc zgłoszenie.

Osobno, ta sama warstwa: pętla nie wie o **początku** tury wykonawcy (dowiaduje
się z pierwszego fragmentu), więc wykonawca zlecony tuż przed końcem tury
koordynatora jest przejściowo niewidoczny. Szczelina jest nieszkodliwa —
samopodjęcie biegu po ukończeniu i tak rusza kolejny obieg, stan kontrolki sam
się prostuje. Domknięcie do zera: `petla.ZTuraWBiegu(rozmowa.CzyTuraWBiegu)`,
jeden wiersz w `montaz_rozmowa.go` — plik poza terenem. Uwaga dla wpinającego:
w chwili `ZakonczTure` mapa `biegnace` wciąż zawiera okno koordynatora
(`zapomnijBieg` jest `defer`), więc pytać wolno tylko o wykonawców.

### Martwa funkcja `sterZlecenia` z komentarzem w nieistniejące miejsce

`core/adapter_przejecie_sterowania.go:146` — doc funkcji `sterZlecenia` mówi
„Woła ją przekład licznika obiegów w `core/stan_obiegu.go`", a `stan_obiegu.go`
jej nie woła i nikt inny też nie. Funkcja martwa, komentarz kieruje w nieistniejące
miejsce. Zauważone przy terenie `warunek-ukonczenia-zadania`, poza jego zakresem.

### Granica 15 s w uprzęży kontraktu jest ciasna dla warstwy skanera

Ustalenie zmierzone przez teren `odmowy-skanera` i potwierdzone przez kontrolę.
Generyczna uprząż zgodności kontraktu woła każdą komendę rdzenia z twardym
limitem 15 s: `zgodnosc_kontraktu_test.go:215` i `:254`,
`blokady_skutek_test.go:127`. Tymczasem czynność skanera na tej maszynie
dochodzi do ~14,8 s (`scanimage --format=png` ~7,4 s + wykaz `scanimage -L`
~7,4 s), a własne limity warstwy są znacznie wyższe — `granicaWykazuUrzadzen`
45 s, `granicaSkanowaniaUrzadzenia` 5 min (`urzadzenia_skaner.go:63,68`). Pod
obciążeniem maszyny generyczne 15 s ucina czynność przed jej własnym limitem
i odmowa spada na kod czasu (`internal_error`) zamiast rozpoznać brak
urządzenia (`not_found`). Na maszynie budowlanej bez żywego skanera komendy
odmawiają natychmiast, więc chwiejność ujawnia się dopiero tam, gdzie warstwa
realnie sonduje urządzenie — u Operatora na wolnej maszynie ta sama granica
może uciąć legalną odmowę „brak urządzenia".

Teren obejmuje uprząż zgodności kontraktu (pliki poza terenami zaplecza).
Rozstrzygnięcie progu — podnieść do wartości spójnej z limitami warstwy —
należy rozważyć wraz z tym, jak uprząż ma traktować komendy o własnych,
dłuższych limitach.

### Rejestr rozjechany na gałęzi centrum poprawek — rozstrzygnięcie przy scaleniu

Gałąź `teren/centrum-poprawki` niesie dwie rewizje dotykające tego pliku —
`6361df8` i `31e37de` — założone tam omyłkowo przez poprzedniego Prowadzącego,
bo drzewo główne stoi przełączone na tę gałąź. Ich treść przeniesiono na `main`
jako `a52fa05` i `ea0c673`, ale na gałęzi zostały. Od tamtej pory `main`
przerobił ten plik mocno: pięć terenów zamkniętych, zgłoszenia i reguły odbioru.

**Przy scaleniu gałęzi designu konflikt w tym pliku rozstrzyga się na rzecz
`main`.** Wersja z gałęzi cofnęłaby zamknięcia terenów i zgłoszenia z pomiarów,
a wyglądałoby to na zwykłe scalenie. Praca w `design/` scala się normalnie —
rozstrzygnięcie dotyczy wyłącznie `prowadzenie/rejestr-terenow.md`.

### Kontrakt ruszył dwa razy od przejęcia — suma bieżąca `b7d0436880d878e78576`

Do 27.08.2026 kontrakt stał nietknięty pod sumą `2cbb843d33f4531b05cd`. Dwie
zmiany, obie wyłącznie dodania, obie przez porównanie strukturalne potwierdzone
jako niezmieniające niczego istniejącego:

| Teren | Co dołożył | Suma po |
|---|---|---|
| `pomiar-stron` | trzy komendy pomiaru strony (1077 → 1080) | `334705bd88c2efc13779` |
| `warunek-ukonczenia-zadania` | czwarta wartość `LoopStopReason` — `completed` | `b7d0436880d878e78576` |

Sprawdzian nietykalności kontraktu odnosi się do sumy **bieżącej**
`b7d0436880d878e78576`, nie do sumy z pierwszego dnia; kolejny teren, który
kontrakt dokłada, podaje w rejestrze sumę zastaną i sumę po sobie.

### Wdrożenie musi założyć trzy nastawy, inaczej stojące wagi leżą odłogiem

Ustalenie terenu `zaplecze-modeli`, zmierzone na żywym rdzeniu. Silnik wiedzy
działa na wagach z `/opt/danaco-modele/embedder`, ale **dopiero po nastawach** —
wartości domyślne rdzenia wskazują co innego i wdrożenie pobierze drugi model
zamiast użyć stojących 4,3 GB.

| Nastawa | Wartość | Bez niej |
|---|---|---|
| `wiedza_model` | `BAAI/bge-m3` | rdzeń sięga po `mpnet` z migracji 115 |
| `wiedza_katalog_modeli` | `/opt/danaco-modele/embedder` | wagi pobierane na nowo do katalogu danych |
| `mowa_katalog_modeli` | `/opt/danaco-modele/mowa` | wagi mowy stoją w pamięci podręcznej konta, które uruchomiło rdzeń |

Wag mowy dotyczy osobne ustalenie: 464 MB modelu `faster-whisper-small` stoi
dziś w pamięci podręcznej pod katalogiem domowym, bo `mowa/ustawienia.go` przy
pustej nastawie zostawia miejsce bibliotece, podczas gdy `wiedza/pomocnik.go`
przy pustej nastawie **przypina** wagi do katalogu danych rdzenia — droga
wołania procesu nie dziedziczy środowiska, więc pomocnik nie zna nawet `HOME`.
Skutek: czyszczenie pamięci podręcznej kasuje działającą funkcję bez śladu
w produkcie, a wagi są przywiązane do konta uruchamiającego. Rozstrzygnięcie
przyjęte: przenieść wagi do `/opt/danaco-modele/mowa` i wskazać je nastawą.

Do rozstrzygnięcia zostaje jedno: czy zrównać obie rodziny w kodzie, żeby pusta
nastawa mowy znaczyła to samo co pusta nastawa wiedzy. To zmiana zachowania
domyślnego wraz z migracją — `migracja_075_mowa.sql` niesie dziś wartość pustą
zgodną ze stałą co do znaku, więc ruszenie samej stałej stworzyłoby dwie prawdy.

### Trzy modele bez komendy w kontrakcie

Ustalenie terenu `zaplecze-modeli`. Wagi stoją, komend nie ma — rejestr zabronił
ich dokładania, bo kontrakt należał w tej turze do innego terenu.

| Model | Waga | Propozycja obszaru |
|---|---|---|
| reranker (bge-reranker-v2-m3) | 2,2 GB | `knowledge` — przełącznik `rerank` w `knowledge.search` albo osobna `knowledge.rerank` |
| CLIP | 1,6 GB | `library` albo `knowledge` — wyszukiwanie obrazów po znaczeniu |
| twarze (GFPGAN, codeformer) | 692 MB | pole `faces` w `image.upscale` **już istnieje**; brakuje silnika: wydanie ncnn nie niesie sieci twarzowej, a wagi `.pth` żądają stosu torch, którego rdzeń nie woła. Odmowa nazywa dziś brak `gfpgan-ncnn-vulkan` |

Osobno rozstrzygnięte i zamknięte: **Kokoro nie jest wart deklaracji** — 54 głosy,
żadnego polskiego (`pf`/`pm` to portugalski), a piper z `pl_PL-darkman-medium`
działa i rdzeń go zna.

### Trzy narzędzia treści pisanej bez legalnego miejsca wpięcia

Ustalenia terenu `dokumenty-i-tekst`, każde zmierzone, żadne nienaprawialne
w granicach tamtego terenu.

| Rzecz | Zmierzone | Skutek |
|---|---|---|
| LanguageTool przy `translate.quality.check` | tryb `--bitext` (LT 6.6) na parze z jaskrawą rozbieżnością liczb (10:30→9:30, 250→500 EUR) oddaje **zero ustaleń**, a wyjścia maszynowego dla bitext program nie ma wcale | wpięcie dałoby instrument milczący przy realnej usterce; komenda mierzy wierność panelu wobec źródła i dziś nie ma czym |
| unpaper przy `research.source.ocr` | uchwyt `adapter_modul_badania_lektura.go:589` upuszcza pola `preprocess` i `languages`, a `DocumentTextExtractRequest` kontraktu tych pól nie ma | obróbka wstępna nieosiągalna z drogi badań; pole `preprocess` pozostaje martwe |
| Tika przy `library.metadata.get` | obsługa stoi w `adapter_modul_library_technika.go` — pliku terenu `obraz-i-diagramy`, gdzie exiftool wszedł w tę samą komendę | rozłączność terenów; do rozważenia, czy Tika ma tam co dołożyć wobec exiftoola |

### mermaid-cli nie ma legalnego miejsca wpięcia

Ustalenie terenu `obraz-i-diagramy`, zmierzone uruchomieniem. Obie komendy
z tabeli narzędzi są zamknięte: `design.diagram.render`
(`adapter_modul_design_wykresy.go:790`) leży za zaporą fotografii, która
zabrania **każdego** wołania procesu w plikach `adapter_modul_design*.go` —
plik próbny z samą wzmianką `zewnetrzne.Wolaj` wywrócił zaporę, choć zapora
nazw przepuszcza `mmdc`; nadto kontrakt tej komendy przyjmuje węzły i krawędzie,
nie ma pola na źródło Mermaid. `apps.architecture.export`
(`adapter_modul_aplikacje_architektura.go:427`) niesie w nagłówku
rozstrzygnięcie „żaden format nie woła programu z zewnątrz, więc eksport działa
na instalce, która niesie sam rdzeń".

Samo narzędzie działa: z przeglądarką stojącą w `/opt/ms-playwright` renderuje
źródło w kształcie eksportu rdzenia (SVG 12 539 B, PNG 265×278). Bez wskazania
przeglądarki `mmdc` odmawia — puppeteer żąda własnej kopii.

Do rozstrzygnięcia: czy rysowanie diagramu ma powstać poza obszarem Design
(nowa komenda i nowe miejsce), czy zostaje niezrobione. Rozstrzygnięcie dotyka
zapory niosącej rozstrzygnięcie Właściciela, więc nie jest samą robotą.

### Komunikat rdzenia radzi budowę usuniętymi skryptami

`budowa/server/internal/narzedzia/wpiecie.go:52` i `:111` — komentarz i treść
komunikatu błędu wskazują `scripts/wydanie.sh` i `scripts/pakowanie.sh`,
usunięte z repozytorium jako materiał instalek natywnych (pozycja 8 pkt 4).
Odmowa radzi Operatorowi drogę naprawy, której nie ma. Drobny teren rdzenia.

### Skrypty hybrydowe wskazują usunięte profile i nieistniejący katalog

Wszystkie trzy `budowa/scripts/instalka-hybryda-*.sh` wskazują profile
`tauri.hybryda-*.conf.json`, których nie ma (w `budowa/desktop/src-tauri/`
stoi wyłącznie `tauri.conf.json`), oraz `budowa/client/` zamiast
`budowa/klient/`; to samo `pakiet-serwera.sh` i `pokaz.sh`. Hybryda jest
jedyną postacią produktu, więc skrypty są żywe, a odmawiają na starcie.
Naturalny moment terenu: powstanie prawdziwej budowy `budowa/klient/dist`.
Osobno do rozstrzygnięcia: `instalka-hybryda-linux.sh` wobec pozycji 8
(jedyna platforma Windows 11), gdy `droga.rs` wciąż niesie tor AppImage.

### Kontrast metadanych w przedsionkach — naprawione

`design/zasoby/przedsionek.css` używał `--dn-tekst-3` (#787C85, szary-500) dla
metadanych 12 px w ośmiu klasach (`pd-kafel-opis`, `pd-sesja-meta`, `pd-kafel-meta`,
`pd-nadtytul`, `pd-motto`, `pd-strefa-opis`, `pd-listwa-meta`, `pd-filtr`) — łamiąc
regułę żetonu „wyłącznie ≥18,66 px" i dając ~28 węzłów poniżej progu WCAG na okno.
Wszystkie osiem przeniesiono na `--dn-tekst-2` (dark szary-400, light szary-600).
Zweryfikowane: 4 przedsionki × 2 motywy = 0 naruszeń axe. Przy okazji zdjęto
zastany `aria-selected="true"` z przycisku `.pd-sesja` (nieprawidłowy ARIA na
`<button>`) na `aria-current="true"` we wszystkich czterech oknach, z hakiem CSS
`.pd-sesja[aria-current='true']`. Decyzja Właściciela: naprawić teraz.

Pozostałe pliki z wpisu „Kontrast warstwy wspólnej" (`rama.css`, `stanowisko.css`,
`css/komponenty.css`, `panel-sesji.css`) nie były tu ruszane — to osobny zakres.

### Żeton mikro-odstępu `--dn-od-05` — wprowadzony przy domknięciu centrum

`zetony.css` dodaje `--dn-od-05: 2px` — półstopień siatki 4 px, między `--dn-od-0`
i `--dn-od-1`. Powód: kryterium maszynowe etapu 1 wymaga, by odstępy pochodziły
z żetonów, a najciaśniejszy realny odstęp (etykieta↔opis w kaflu) wynosił 2 px
bez pokrycia w skali. Zdjęto surowe `gap: 2px` z `okna/centrum-dowodzenia.css`
(2×) i `okna/studio.css` (4×). Decyzja delegowana przez Właściciela: „wariant
zgodny z profesjonalnym standardem" — skala odstępów zawiera używane wartości
jako żetony (wzór: Tailwind `0.5`, Material). Zmiana warstwy wspólnej odnotowana
tutaj zgodnie z ograniczeniem warstwy wspólnej z planu etapów.

### Droga wejścia — braki po stronie rdzenia, ujawnione pomiarem

Wszystkie pochodzą z terenu `droga-wejscia` i są zmierzone na żywym rdzeniu.
Żaden nie wstrzymał pracy; okno podłącza je dziś do odmowy nazywającej brak.

| Rzecz | Zmierzone | Skutek |
|---|---|---|
| brak drogi powtórnego potwierdzenia adresu | `auth.verify` z cudzą drogą → `not_authenticated`; powtórny `auth.register` → `conflict`; żadna z dziewięciu komend `auth.*` nie wydaje drogi drugi raz | Operator, do którego list nie dotarł, **nie ma wyjścia z okna** — dotyczy także instalacji **z pocztą**, nie tylko bez |
| trzy z pięciu etapów przygotowania środowiska bez komendy | prototyp wymienia pięć; droga wejścia obsługuje dwa. „Profil i uprawnienia", „Kanały modeli", „Magistrala kontekstu" nie mają przypisanej komendy | pasek postępu zatrzymuje się na 40% — wartość prawdziwa, nie ozdobna |
| odmowa logowania nie niesie długości zwłoki | osiem prób: 100, 350, 600, 1100, 2100, 4100, 5100, 5100 ms; `details` puste za każdym razem | okno mierzy zwłokę samo, czasem trwania próby poprzedniej |
| magazyn tokenu bramki między uruchomieniami | `connection.hello` przyjmuje token i oddaje `authenticated: true`; gdzie token mieszka, nie rozstrzyga ani prototyp, ani kontrakt, ani rejestr | odsłona „rozpoznano zaufane urządzenie" działa dopiero po wskazaniu magazynu — należy do powłoki |

### Adres serwera wdrożenia nie stoi w żadnym źródle

Teren `powloka-tauri` przeszukał `budowa/`, `design/`, `docs/` i `prowadzenie/`:
pozycja 8 nazwy serwera nie podaje, prototyp instalatora nie ma kroku adresu,
`wydania.json` niesie wyłącznie kanał pobrań. Powłoka bierze adres ze wskazania
Operatora albo ze zmiennej `DANACO_HOST_RDZENIA`. Czy adres ma być wkompilowany
przy składaniu instalki — i skąd wtedy pochodzi — jest rozstrzygnięciem
Właściciela.

### Zastane usterki powłoki, ujawnione przy przejęciu

| Rzecz | Stan |
|---|---|
| `aktualizacja/probne.rs` nie istnieje | `droga.rs:293` importuje `KatalogProbny` w siedmiu sprawdzianach; `cargo test` **nie kompiluje się od chwili przejęcia** |
| `wykonaj_aktualizacje` pobiera spod dowolnego adresu | `pobranie.rs` sprawdza wyłącznie sumę SHA-256, też podaną przez stronę; brak przywiązania do kanału z `wydania.json` |
| cztery z sześciu poleceń IPC bez odbiorcy | `stan_rdzenia`, `wskazanie_rdzenia`, `wybierz_katalog_roboczy`, `wykonaj_aktualizacje` — opisywały odbiorców w usuniętym `client/src/powloka/` |
| `app.security.csp: null` | strona bez polityki treści; ułożenie wymaga interfejsu, który dopiero powstaje |
| sześć skryptów w `budowa/scripts/` | wskazują usunięte profile nastaw, nieistniejący `budowa/client/` i wariant natywny zniesiony pozycją 8; pozycja 8 pkt 4 mówi, że mają leżeć **poza** repozytorium |

### Łańcuchy widoczne poza katalogiem treści w paczce designu

Właściwość, którą pozycja 12 nazywa wiążącą, jest w dwóch miejscach naruszona:
`design/zasoby/okna/wejscie/skladniki/pole-hasla.js:55` niesie zaszyte
`'Pokaż hasło'`, a `design/zasoby/okna/przeplyw-wejscia.js` powiela w kodzie
`SILA_PUSTE`, `SILA_OPISY` oraz dwa napisy o schowku, które stoją już
w katalogu. W kliencie wszystkie te napisy siedzą w katalogu treści.

### Kontrast warstwy wspólnej — gotowe do otwarcia

Żeton `--dn-tekst-3` (`design/zasoby/zetony/zetony.css` w. 222) niesie własną
regułę: wyłącznie metadane i tekst od 18,66 px półgrubego. Warstwa wspólna łamie
ją w około siedemdziesięciu miejscach: `rama.css` 24, `css/komponenty.css` 16,
`okna/centrum-dowodzenia.css` 14, `stanowisko.css` 8, `prototyp.css` 4,
`panel-sesji.css` 4. Żeton jest poprawny — wadliwe jest jego użycie. Dopóki to
stoi, każde okno korzystające z ramy niesie naruszenia wagi `serious` i żaden
teren nie domknie kryterium dostępności bez wyjątku.

| Plik i wiersz | Rzecz | Zmierzony kontrast |
|---|---|---|
| `rama.css:981` | `.dn-stan` — pasek stanu, 6 pozycji | 3,85 ciemny · 3,99 jasny |
| `stanowisko.css:148` | `.sta-kom-pole span`, zaszyte 10 px | 3,85 · 3,99 |
| `stanowisko.css:288` | `.sta-wpis-godzina`, zaszyte 10 px | 4,25 · 4,17 |
| `stanowisko.css:123` | `.sta-okno-znacznik` | 3,85 · 3,99 |
| `css/komponenty.css:226` | `.dn-pole-opis` przy 13 px | 4,17 · 4,25 |
| `panel-sesji.css:59, 139` · `okna/centrum-dowodzenia.css:258, 288` | tekst 12 px | poniżej progu |

Osobno, ta sama warstwa: `zetony.css:69` deklaruje przy `--dn-rama-tekst-3`
kontrast 4,74 : 1 na ramie. Zmierzone: 4,45 : 1, czyli poniżej progu 4,5.
Deklaracja w komentarzu jest nieprawdziwa.

### Mechanizmy warstwy prototypu — gotowe do otwarcia

| Plik i wiersz | Usterka | Waga |
|---|---|---|
| `prototyp.js:66–83` | `przelaczWidok` nadaje `aria-selected` elementom `<button>` i `<a>`, którym atrybut nie przysługuje | krytyczna |
| `prototyp.js:88–102` | wędrujący `tabindex` ustawiany tylko przy starcie, nieodświeżany po przełączeniu | poważna |
| `prototyp.js:300–341` | wstrzykiwany pasek prototypu stoi poza punktami orientacyjnymi | umiarkowana |
| `powloka.js` | montuje pełną powłokę bezwarunkowo; brak trybu „rama dopiero po uwierzytelnieniu" | poważna |
| `css/komponenty.css:1014` | `.dn-postep-wartosc` bez `display: block` — pasek postępu renderuje pusty tor | poważna |
| `css/komponenty.css:1481–1499` | `.dn-alert` bez gniazda ikony, bez części tytułu i treści, bez wariantów błędu i powodzenia | poważna |
| `rama.css:102–114` | `.dn-narzedzia-pas` na barwie gruntu roboczego — przyczyna źródłowa braku rozgraniczenia wstążek | poważna |
| `rama.css` — szyna | strefa środowisk z 38 pozycjami wypycha strefę szybkiego wyboru poza kadr | poważna |

### Wstążka narzędziowa do wyniesienia do warstwy wspólnej

Szkielet wstążki okna roboczego istnieje dziś **wyłącznie w arkuszu modułu
Studio**. Kontrakt wstążki obowiązuje wszystkie moduły — pozycja 9 rejestru
decyzji — więc przy drugim module rozjedzie się bez niczyjej złej woli.

Do wyniesienia: forma paska wraz z trzema strefami, stała wysokość, zwijanie
prawej grupy do menu nadmiaru, trwały stan wybrania narzędzia, nieruchomość przy
przewijaniu treści. Zmienna pozostaje wyłącznie zawartość stref, właściwa
rodzajowi karty.

### Komponent kart okna — wstrzymane

Wstrzymane do rozstrzygnięcia Właściciela w sprawie nazw dwóch pięter kart.
Otwarcie terenu na drugi moduł przed tym rozstrzygnięciem odtworzy pomieszanie
z biblioteki w każdym kolejnym prototypie.

| Plik i wiersz | Rzecz |
|---|---|
| `karty-okna.css` | arkusz nosi nazwę poprawną, a definiuje kontener `.dn-karty-sesji`, w którym stoją elementy `.dn-karta-widoku` — sprzeczność w jednym pliku |
| `karty-okna.js:27` | mechanizm kart zakotwiczony w `.dn-obszar-panel--glowny` z `rama.css`; w oknie roboczym komponent jest martwy |
| `karty-okna.css:140` | przycisk zamknięcia wewnątrz `div[role="tab"]` — zagnieżdżona interaktywność, przenosi się na każde okno używające komponentu |
| `.dn-izolacja-wskaznik` | stoi w ramie aplikacji, a izolacja jest cechą okna roboczego (`izolacja.css`) — etykietuje niewłaściwy poziom |

### Usterki zastane, ujawnione przy próbce — gotowe do otwarcia

Wszystkie sprzed terenu, żadna nie jest jego skutkiem.

| Miejsce | Rzecz |
|---|---|
| `design/INDEKS.html` | martwy odsyłacz `href="kontrakt systemu projektowego"` — jedyne 404 wśród 93 odsyłaczy strony |
| `design/zasoby/okna/studio.js` oraz `studio.html` w. 786 | zdublowana obsługa `[data-srod]`; klik w przycisk trybu wywołuje dwa komunikaty naraz. Zachowana bez zmiany, bo kryterium wymagało zachowania identycznego |
| plansze i indeks | metryki „N linii · M interakcji" oraz „N w. · M kB" rozjechane ze stanem plików także dla kart nietkniętych — `centrum-dowodzenia.html` podane jako 627 wierszy i 42 kB przy faktycznych 1023 wierszach i 305 kB |

Metryki wymagają rozstrzygnięcia Właściciela: albo są normatywne i dostają
definicję sposobu liczenia, albo znikają. Metody liczenia „interakcji" nie da
się odtworzyć z treści plików, więc dziś nikt nie jest w stanie ich utrzymać.

### Droga SANE bez odpowiednika `bladWarstwyWia` — gotowe do otwarcia

Brak urządzenia jest już nazwany. Brak samego programu `scanimage` dalej wychodzi
odmową arsenału bez wskazania drogi obejścia, podczas gdy `bladWarstwyWia`
(`urzadzenia_skaner.go`) dla tej samej sytuacji na Windowsie podaje
`studio.ingest.queue.add`. Ta sama asymetria dotyczy `wykazSkanerow`: gałąź
Windows przekłada odmowę, gałąź Linux oddaje ją surową. Operator na Linuksie bez
`sane-utils` nie dowie się, że materiał da się wnieść inną drogą.

### Sprawdzian katalogu akcji szuka nieistniejącego katalogu — naprawione

`katalog_akcji_test.go` szukał `../../../client/src/ikony/zrodla`, którego nie
ma w żadnej gałęzi. Rewizja `9ba1f90` wiąże sprawdzian ze źródłami ikon nowego
klienta — `sciezkaZrodelIkon` wskazuje `../../../klient/src/ikony/zrodla`.
Przechodzi w biegu odniesienia 2041 zdanych, zero niezdanych.

### Reguła odbioru wyprowadzona z pomiarów

`axe.run()` sam wywołuje dwa błędy 404 (`menu.css`, `ruch.css`), bo rozwiązuje
`@import` względem adresu dokumentu, a nie arkusza. Konsola przed wstrzyknięciem
axe jest pusta. Tych dwóch wpisów nie liczy się jako brudnej konsoli.

**Pełny bieg wymaga podniesionego limitu czasu.** Pakiet `internal/core`
przekroczył domyślną granicę `go test` (10 minut) po dołożeniu przez cztery
tereny około czterdziestu sprawdzianów wołających prawdziwe programy zewnętrzne.
Bieg pada wtedy paniką „test timed out after 10m0s" wskazującą **przypadkowy**
sprawdzian, który akurat biegł — wygląda to na usterkę tego sprawdzianu i nią
nie jest. Objaw uboczny myli podwójnie: `gotestsum` naliczy wtedy około 947
sprawdzianów zamiast ponad dwóch tysięcy, bo pakiet, który spanikował, nie
policzy swoich. Pełny bieg uruchamia się z `-timeout 30m`.

**Pełne biegi szereguje się zamkiem.** Maszyna ma 16 rdzeni; przy czterech
terenach naraz biegło na niej dziewięć procesów sprawdzianów i sprawdziany
z twardymi granicami czasu zaczęły się chwiać — każdy pełny bieg gubił **inny**,
a pojedynczo wszystkie przechodziły. Granice tego rodzaju stoją m.in.
w `zgodnosc_kontraktu_test.go`. Bieg pełny wykonuje się pod zamkiem:

```bash
cd budowa/server && flock /tmp/danaco-bieg-pelny.lock \
  gotestsum -- -count=1 -timeout 30m ./...
```

Sprawdzian, który padł w pełnym biegu, a przechodzi uruchomiony pojedynczo, jest
chwiejnością pod obciążeniem, nie usterką — i tak się go nazywa.

**Pełny bieg regresji jest bramką scalenia u Prowadzącego, nie bramką wyjścia
u wykonawcy.** Wykonawca dowodzi swojego przedmiotu biegiem **celowanym** pakietu,
który dotyka — jest szybki i odporny na obciążenie. Pełny bieg `./...` w drzewie
terenu mierzy drzewo **bez** pozostałych terenów biegnących równolegle, więc
rozstrzyga wyłącznie bieg na scalonym `main` po scaleniu, a ten wykonuje
Prowadzący raz, pod zamkiem, jako warunek bramki scalenia. Zaprzęganie czterech
terenów do czterech pełnych biegów naraz zapycha zamek na godziny i nic nie
rozstrzyga — każdy z nich mierzy inny, niepełny stan.

Czas dostępu do pliku nie dowodzi, że pliku nie czytano. Drzewo stoi na `ext4`
zamontowanym z `relatime`, gdzie jądro odświeża czas dostępu wyłącznie wtedy,
gdy poprzedni jest starszy od czasu zmiany albo starszy niż doba. Odczyt pliku,
którego czas dostępu jest już późniejszy od czasu zmiany, **nie zostawia
śladu** — niezmieniony czas dostępu jest tam brakiem pomiaru, nie dowodem
nietknięcia. Wyszło to przy sporze o sprawstwo pobrania wag mowy: wykonawca
podał czas dostępu jako dowód, a instrument z założenia milczał.

Pomiar w przeglądarce wymaga jawnego ustawienia `PLAYWRIGHT_BROWSERS_PATH` na
`/opt/ms-playwright` w poleceniu, a nie polegania na środowisku powłoki — powłoka
uruchomiona przed ustawieniem zmiennej jej nie widzi i pobiera przeglądarki
po raz drugi.

## Tereny zamknięte

| Nazwa | Gałąź | Rewizje | Kontrola |
|---|---|---|---|
| `usterki-rdzenia` | `teren/usterki-rdzenia` | trzy rewizje | weryfikacja Prowadzacego pomiarem: sprawdzian `TestTuraZamknietaBledemNieOglaszaUkonczenia` **padl na kodzie sprzed naprawy** z wlasciwym zdaniem i przeszedl po przywroceniu; zakres wylacznie `internal/core`; kontrakt nietkniety |
| `nastawy-wdrozenia` | `teren/nastawy-wdrozenia` | `678ff3e` | weryfikacja Prowadzacego pomiarem: 4 pliki w zakresie, migracja 115 nietknieta, nowa migracja 401; sonda dodatnia wykonawcy pokazala **206 polaczen i 1,1 GB pobrania przed naprawa wobec zera po niej** |
| `powloka-i-wydanie` | `teren/powloka-i-wydanie` | `13b947b`, `8d99ae4` | weryfikacja Prowadzacego pomiarem: wykaz wydan sprowadzony do dwoch postaci hybrydowych (0 trafien wzorca natywna/AppImage przy 10 kontrolnych); sprawdzian wiazacy adres padl po wprowadzonym rozjezdzie i przeszedl po cofnieciu; 25 sprawdzianow powloki zdanych |
| `brama-i-droga-wejscia` | `teren/brama-i-droga-wejscia` | `23a6b84` wyjątek powitania · `e41dd78` straże drogi bez poczty · `9a8ec0c` brak skanera | weryfikacja Prowadzącego pomiarem: bieg wymuszony `-count=1` 537 s — 2029 sprawdzianów, 1 niezdany wobec 4 zastanych; powitanie niepełne odpowiada wersją protokołu na żywym rdzeniu, `channel.add` z brakiem pola dalej odmawia; kontrakt nietknięty |
| `prototypy` | `teren/prototypy` | paczki instalatora i drogi wejścia | przyjęte przez Właściciela; weryfikacja Prowadzącego pomiarem: oba okna wczytują się bez błędu konsoli, zero łańcuchów widocznych poza katalogiem treści |
| `fundament-klienta` | `teren/fundament-klienta` | `d19bfeb` warstwa połączenia i protokołu | weryfikacja Prowadzącego pomiarem: kompilacja bez błędu, 17 sprawdzianów zdanych, rozmowa z żywym rdzeniem, generat bajtowo powtarzalny, zero dotknięć DOM, kontrakt nietknięty |
| `naprawy-rdzenia` | `teren/naprawy-rdzenia` | `060d5b7` naprawy i brama kontraktu | weryfikacja Prowadzącego pomiarem: 2022 zdane wobec 2003 zastanych, te same 4 niezdane, kontrakt nietknięty |
| `proba-prototypow` | `teren/prototypy` | `66e5ee0` przepływ wejścia · `61a5867` moduł Studio · `53bc3d5` odsyłacze | kontrola sesji nadzorującej wykonanie, weryfikacja Prowadzącego pomiarem |
| `pomiar-stron` | `teren/pomiar-stron` | `607a3d8` kontrakt · `4191c4b` trzy osie pomiaru · `6b189e7` sprawdziany | kontrola osobnej sesji: kontrakt porównany strukturalnie, nie liniowo — zero sekcji zmienionych, żaden istniejący element nietknięty bajtowo, dodane 3 komendy, 6 struktur, 3 wyliczenia; generator powtórzony dwa razy, sumy zgodne; oś obciążenia zmierzona niezależnie własnym klientem kanału — 102 666 żądań wobec 102 669 na liczniku serwera sprawdzianu; osiem sprawdzianów zdanych, żaden pominięty. Scalone `73e9aad` |
| `zaplecze-modeli` | `teren/zaplecze-modeli` | `090acd6` silnik osadzeń na wagach stojących | kontrola osobnej sesji: teren zwrócony za niepełny wykaz postawionego, po uzupełnieniu przyjęty; kontroler powtórzył pomiar osadzeń niezależnie i przy odciętej sieci (`HF_HUB_OFFLINE=1`) — wymiar 1024, normy 1,000000, cos(kot,kot) 0,8348 wobec 0,2863–0,3353 dla par odległych; wagi `/opt` czytane, nie pobierane po raz drugi; kontrakt nietknięty. Scalone `ba82cb6` |
| `obraz-i-diagramy` | `teren/obraz-i-diagramy` | `f7fab06` dogniecenie zapisu i metadane osadzone | kontrola osobnej sesji własnym biegiem: osiem zapór zdanych nietkniętych, 33 sprawdziany zdane w realnych czasach z przytoczonymi rozmiarami przed i po (PNG 3322→2456 B, JPEG 45399→39623 B, model barw paleta dowodzi wejścia pngquanta); sprawdziany maszyny bez programu wytwarzają ją naprawdę (`t.Setenv` na pusty katalog); kontrakt nietknięty; zero śladu mmdc w rdzeniu. Scalone `5b98585` |
| `dokumenty-i-tekst` | `teren/dokumenty-i-tekst` | `ea05603` sześć programów treści pisanej · `aa695e6` poprawka nazwy po zwrocie | kontrola osobnej sesji: teren zwrócony za martwą nazwę `zasiegSyntezy` w komentarzu, po poprawce przyjęty; własny bieg kontrolera 14 zdanych, zero pominiętych; sprawdziany mierzą skutek — PDF czytany drugą komendą i innym programem, korekta aż po treść panelu w bazie, OCR dwoma przebiegami, brak programu wywołany `t.Setenv`; wybory hunspell i typst zweryfikowane uruchomieniem; kontrakt nietknięty. Scalone `43d05f6` |
| `odmowy-skanera` | `teren/odmowy-skanera` | `1eaa58a` parytet odmów skanera Linux wobec Windows | kontrola osobnej sesji: `bladWarstwySane` wierne lustro `bladWarstwyWia` — ten sam kod `channel_unavailable`, ta sama droga obejścia `studio.ingest.queue.add`, pakiet czytany z braku nie zaszyty; sprawdzian parytetu woła obie warstwy i wymaga jednego kodu; brak wymuszony atrapą PATH z potwierdzeniem `zewnetrzne.Stoi=false`; bieg celowany 3 zdane; kontrakt nietknięty. Scalone `1dca877` |
| `sprawdziany-drogi-wejscia` | `teren/sprawdziany-drogi-wejscia` | `1393e74` sprawdzian zdjęcia znacznika bramki | weryfikacja Prowadzącego mutacją rdzenia: wyłączenie `zdejmijZnacznikBezPoczty` daje sprawdzian niezdany („znacznik przeżył swój powód"), przywrócenie — zdany; trzy sprawdziany rozjazdu z rdzeniem były już przerobione przez `brama-i-droga-wejscia`, wykonawca to zmierzył i nie tknął; bieg celowany osiem zdanych; kontrakt nietknięty, jeden plik terenu. Scalone `4e70ef1` |
| `odwolania-do-usunietych-skryptow` | `teren/odwolania-do-usunietych-skryptow` | `8451157` żywy skrypt pakietu w odmowie braku serwera narzędzi | **wykonawca urwał się na limicie sesji przed rewizją; pracę dokończył i zweryfikował Prowadzący sam, bez niezależnej kontroli — odstępstwo od rozdziału ról, wymuszone urwaniem, w trybie samodzielnym na polecenie Właściciela.** Weryfikacja: droga naprawy w odmowie istnieje — `scripts/pakiet-serwera.sh` żyje i buduje `danaco-narzedzia` z `server/cmd/danaco-narzedzia`, stawiając obok rdzenia (w. 72,83,101); zero odwołań do zniesionych skryptów w rdzeniu; próba mutacji: podmiana stałej na `wydanie.sh` daje obie straże niezdane, przywrócenie — zdane; kontrakt nietknięty. Scalone `8e504c3` |
| `warunek-ukonczenia-zadania` | `teren/warunek-ukonczenia-zadania` | `06c81a1` rdzeń odróżnia ukończenie od trzech zatrzymań | weryfikacja Prowadzącego (tryb samodzielny) dwiema mutacjami: usunięcie mapowania `completed` daje „rozróżnienia maszynowego nie ma", wyłączenie samopodjęcia daje „ukończenie zachowuje się jak bramka akceptacji" — obie przywrócone zdane; kontrakt tylko dodania (`334705…` → `b7d04368…`), generator bajtowo powtarzalny w dwóch przebiegach; wykonawca ustalił pomiarem, że koniec tury koordynatora dotąd ginął (Wybudzacz odrzucał okna niebędące wykonawcami); brak bramki akceptacji zrobiony mechanizmem samopodjęcia, nie deklaracją. Scalone `e677a65` |
| `aktualizacja-powloki-i-skrypty` | `teren/aktualizacja-powloki-i-skrypty` | `d1b6207` probne.rs · `54cb287` kanał pobrań · `880b41f` skrypty natywne | kontrola osobnej sesji własnym biegiem: `cargo test` 23 zdane, zero niezdanych (stan zastany: nie kompilował się); sha256 sześciu kopii w materiale zamkniętym tożsame; `ADRES_KANALU` zgodny znak w znak z `kanal.adres` wykazu wydań; scalone `adb2a85`, drzewo scalone tożsame z kontrolowanym. Uwaga trwała: kompilacja powłoki wymaga `budowa/klient/dist` (w `.gitignore`) i zmiennych CARGO/RUSTUP ze środowiska maszyny |

Teren `naprawy-rdzenia` scalony do `main`. Piąta usterka — wartość domyślna
`createVersion` — wróciła jako zgłoszenie, bo kontrakt jej nie ustala. Brama
kontraktu ujawniła, że powitanie kanału musi stać poza nią; rozstrzygnięcie
niesie pozycja 10 rejestru decyzji.

Wynik terenu `proba-prototypow`, przeniesionego na gałąź `teren/prototypy`: oba przedmioty wykonane, wszystkie kryteria spełnione. Zakresy trzech
rewizji rozłączne — sprawdzone. Drzewo czyste. Kryterium 7a zwraca zero trafień
w całym repozytorium, nie tylko w `design/`. Gałąź czeka na ocenę kierunku
przez Właściciela; **nie jest scalona** — próbka rozstrzyga kierunek, a nie
wnosi dorobek.

Usterki wykazu popełnione przez Prowadzącego, obie wychwycone przez sesję
nadzorującą: wpisanie katalogu `przeplyw/` zamiast nazw plików oraz wpisanie
nieistniejącego `design/README.md` (pliki README stoją wyłącznie
w podkatalogach). Wniosek na przyszłość: wykaz plików terenu sprawdza się
odczytem drzewa przed otwarciem, nie z pamięci.

## Wzór wpisu

Otwarcie terenu wymaga wypełnienia wszystkich pól. Pole puste blokuje otwarcie.

- **Nazwa** — rzeczownikowa, opisuje przedmiot pracy, bez oznaczeń literowych
  i numerycznych.
- **Gałąź bazowa** — gałąź, z której teren wyrasta i do której wraca.
- **Wykaz plików** — pełne ścieżki albo katalog. Wykaz nie może przecinać się
  z żadnym terenem otwartym.
- **Kryteria odbioru** — zdania sprawdzalne. Kryterium, którego nie da się
  sprawdzić uruchomieniem albo odczytem pliku, nie jest kryterium.
- **Wykonawca** — oznaczenie sesji prowadzącej teren.
