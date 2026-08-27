# Centrum dowodzenia — uzasadnienia reguł

Zapis powodów, dla których reguły okna są takie, a nie inne. Arkusze i mechanizm
niosą sam kod; uzasadnienie stoi tutaj, bo limit gęstości komentarzy w kodzie
wynosi 250 znaków na 1000 wierszy i nie pomieściłby tych treści.

Wiąże się z terenem `centrum-poprawki` w [rejestrze terenów](../../prowadzenie/rejestr-terenow.md).

## Arkusz `zasoby/okna/centrum-dowodzenia.css`

### `:root`

DANACO CONSOLE — CENTRUM DOWODZENIA: ZŁOŻENIE OKNA ---------------------------------------------------------------------------- Kompetencja: układ karty Centrum dowodzenia — nagłówek, wiersz pierwszego wejścia oraz cztery strefy wyboru (środowiska, komponenty aplikacji, komponenty własne Operatora, listwa ustawień) i panel samouczka. Wygląd komponentów pochodzi z arkuszy kompetencyjnych; tutaj stoi wyłącznie ich złożenie w tym oknie. Wymaga: `zetony/zetony.css`, `zetony/ruch.css`, `css/komponenty.css`, `stany.css`, `rama.css`, `karty-okna.css`, `panel-sesji.css`, `pasek-okna.css`, `izolacja.css`, `okno-robocze.css`. Mechanizm: `okna/centrum-dowodzenia.js`. Punkty łamania mierzą szerokość płótna (`@container plotno`), nie szerokość ekranu: panel samouczka zwęża okno główne o 320 px.

### `.cd-strona`

Włos marki ──────────────────────────────────────────────────────────── Obrys wstążek i okien niesie barwę znaku zmieszaną z obrysem podstawowym. Udział barwy znaku jest dwustopniowy: wyższy dla krawędzi, niższy dla tła karty bieżącej, gdzie barwa leży pod tekstem i musi zostawić mu kontrast. Mieszanie w przestrzeni oklab daje w obu motywach równy stopień odejścia od szarości, czego mieszanie kanałowe nie zapewnia. Wartości udziału stoją tu, a nie w `zasoby/zetony/`, bo rozstrzygnięcie obejmuje jedno okno; przy rozszerzeniu na pozostałe przenieść do palety.

### `.cd-strona .dn-karta-widoku[aria-selected='true']`

Bez mieszania barw tło karty bieżącej równa się tłu wstążki, więc nośnikiem wyróżnienia jest kreska sygnałowa u dołu karty. Miarę niesie żeton wstęgi aktywności karty, nie żeton siły hasła.

### `.cd-strona .dn-narzedzia`

Wstążka główna nie miała obrysu — dostaje włos marki, wewnątrz obecnej bryły, więc jej wysokość nie drgnie.

### `.cd-strona .dn-obszar-panel`

Trzy okna obszaru roboczego: boczne, robocze i samouczek. Sama barwa — skrót `border` przywracałby krawędzie, które warstwa wspólna zeruje celowo: dół każdego panelu (wchodzi pod pasek stanu), lewą okna bocznego i prawą samouczka (stykają się z ramą aplikacji).

### `.dn-obszar-panel--boczny[data-strona='prawo']`

Położenie panelu bocznego ──────────────────────────────────────────── Panel sesji stoi domyślnie po lewej. Menu panelu przenosi go na prawą stronę okna — obszar roboczy jest zginany w rzędzie, więc kolejność wystarczy zmienić porządkiem, bez przenoszenia znacznika.

### `.cd-tresc:not([hidden])`

Płótno karty. Warunek `:not([hidden])` jest konieczny: `display` zadeklarowany przez autora wygrywa z `[hidden]` z arkusza przeglądarki, a wtedy karta odłożona zostawałaby na ekranie razem z kartą bieżącą.

### `.cd-tresc > *`

Wejście stref w kolejności ważności. `backwards`, nie `both`: animacja wypełniająca w przód pozostaje czynna i czyni ze strefy kontekst nakładania, co przycina panele rozwijane w jej wnętrzu.

### `.cd-marka-lockup`

Lockup poziomy 280 px wobec minimum 120 px z księgi znaku; poniżej tej szerokości pola schodzi razem z płótnem, zamiast wypychać treść poza okno.

### `.cd-marka-lockup .kropka:not(.kropka--konsola)`

Kropka sygnału tętni nieprzerwanie — wspólnym pulsem platformy, 2,4 s. Oba znaki okna, lockup płótna i godło belki, biją tym samym rytmem.

### `.cd-start`

Belka pierwszego wejścia nie jest osobną powierzchnią — to zdanie na stronie, nie blok do oddzielenia. Bez obrysu i tła wyściółka pozioma stawiałaby ją z wcięciem wobec znaczników sekcji poniżej, więc schodzi do zera.

### `.cd-start`

Belka samouczka i listwa ustawień domykają się na siatce środowisk, nie na krawędzi płótna: kończą się równo z drugą kartą przy czterech kolumnach i z pierwszą przy dwóch. Obie miary daje jedno wyrażenie — dwie kolumny wraz z przerwą między nimi mierzą dokładnie tyle, co jedna kolumna z układu dwukolumnowego, czyli połowę siatki pomniejszoną o pół odstępu. `min-width` chroni listwę: jej pięć pozycji potrzebuje więcej miejsca, niż zostaje przy dwóch kolumnach, a treść nie ma prawa wyjść poza obrys.

### `.cd-listwa`

Rząd kontrolek czyta się jako jeden zestaw, więc pole listwy nie zwęża się poniżej połowy siatki z własnej woli. Gdy jednak okno zejdzie tak nisko, że pięć kontrolek nie mieści się w tej połowie, rząd łamie się na dwa wiersze. Podłoga na szerokości treści byłaby tu błędem: `min-width` bije `max-width`, więc poniżej 480 px listwa wychodziła poza pole treści o ponad 140 px.

### `.cd-listwa-poz, .cd-listwa-sep, .cd-listwa-menu`

Pozycja i przegroda zachowują swoją miarę — bez tego rząd nie łamie się, tylko ściska kontrolki do nieczytelnych kikutów.

### `@media (max-width: 860px)`

Jedna kolumna — karta zajmuje całą szerokość, więc belka też. Próg 860 px nie należy do skali `--dn-bp-*`; belka musi łamać się dokładnie tam, gdzie łamie się siatka środowisk, a ta stoi na progach 1180 i 860 px sprzed tej pracy. Przeniesienie siatki na skalę systemową rozstrzyga się osobno — do tego czasu belka idzie za siatką, nie za skalą.

### `.cd-start`

Zdanie zajmuje tyle, ile mierzy; przycisk staje bezpośrednio za nim, w odstępie spacji, a zamknięcie odsuwa się do prawej krawędzi wiersza. Odstęp wiersza schodzi do miary spacji, a znak odzyskuje swój przez własny margines.

### `.cd-sekcja-tytul::before`

Waga strefy zapisana atrybutem `data-waga`; nośnikiem jest kreska sygnałowa przed tytułem, w trzech stopniach krycia.

### `.cd-siatka-srodowisk`

Cztery środowiska stoją w jednej linii — komplet wyboru czyta się jako jeden rząd, a nie jako trzy plus sierota. Karta zwęża się razem z płótnem; nie ubywa kolumn, ubywa szerokości karty.

### `—`

Karta jest kolumną: wiersz meta z przyciskiem wejścia stoi przy dolnej krawędzi, więc cztery przyciski w rzędzie leżą na jednej wysokości.

### `.cd-karta-akcent`

Karta w rzędzie czterech ma węższe pole niż karta w rzędzie trzech: wyściółka i stopień tytułu schodzą o krok, żeby najdłuższa nazwa środowiska (MultitaskingAI) mieściła się w jednym wierszu przy 1366 px.

### `.cd-karta-akcent .cd-metryka-czlon`

Wiersz metryki łamał się sam z siebie i tylko na dwóch kartach z czterech: przy 1366 i 1440 px pole tekstu ma 200 px, a „9 modułów · rozmowa i wiedza" potrzebuje 224 px i „panel orkiestracji · 6 sekcji" 232 px. Zdanie opisowe startowało przez to na dwóch różnych wysokościach i rząd czterech kart czytał się jako niewyrównany. Jednowierszowa metryka jest nieosiągalna: najmniejszy żeton stopnia (--dn-fs-xs, 12 px) daje 203 px w Chromium i 208,8 px w Firefoksie — nadal ponad pole, a zejście z wyściółki karty do 8 px zostawiałoby margines poniżej piksela i pękało przy innych metrykach kroju. Dlatego metryka jest łamana celowo i jednakowo na wszystkich czterech kartach: człon ilościowy w wierszu pierwszym, człon opisowy w drugim. Podział niesie markup, nie zawijanie, więc wypada identycznie w obu silnikach i przy każdej z trzech szerokości okna — a zdanie opisowe zaczyna się wszędzie na tej samej wysokości. Treść metryki bez zmian: separator „·" zostaje i otwiera wiersz drugi.

### `.cd-karta-akcent::before`

Wstęga sygnałowa na górnej krawędzi kończyła się przed łukiem narożnika: pseudoelement kotwiczy się do pola wyściółki, a przycięcie liczy promień pomniejszony o obrys, więc jej koniec wystawał poza łuk i czytał się jak kreska urwana w rogu. Stan karty niesie sam obrys — na czterech krawędziach o jednej geometrii.

### `.cd-siatka-kafli`

Kafel niesie treść, nie pustkę: znak po lewej, nazwa i zdanie o przeznaczeniu po prawej, w jednym wierszu. Układ kolumnowy z osobnym wierszem czynności rozciągał kafel na dwukrotność potrzebnej wysokości i zostawiał w nim puste pole. Licznik zapisanych zniknął z narożnika — zapisane komponenty stoją w strefie własnej Operatora, niżej.

### `padding: var(--dn-od-6);`

Kafel komponentu stoi o połowę wyżej niż stopień biblioteczny; nośnikiem wzrostu jest wyściółka, bo treść ma zostać wyśrodkowana w polu.

### `.cd-kafel .dn-kafel-ikona`

Znak komponentu ma wielkość emblematu środowiska — 32 px, rozstrzygnięcie Właściciela z 2026-08-27. Obie strefy strony głównej niosą znaki tej samej miary; różnicę stref oddaje karta, nie rozmiar znaku.

### `.cd-kafel-tresc .dn-kafel-etykieta`

Nazwa komponentu stoi w stopniu tytułu karty środowiska — obie strefy czyta się jednym stopniem. Opis zostaje na 13 px.

### `.cd-strona .cd-kafel-tresc .dn-kafel-etykieta`

Znak modułu niesie barwę znaku środowiska — obie strefy strony głównej mają ten sam stopień kontrastu. Nazwa modułu schodzi z półgrubej na zwykłą: krój treści (IBM Plex Sans) czyta się przy tej samej wadze ciężej niż krój nagłówkowy karty środowiska (Space Grotesk).

### `.cd-strona`

Grot marki jako wejście ─────────────────────────────────────────────── Wejście do środowiska i do modułu niesie pojedynczy grot ze znaku firmowego wraz z kropką sygnału — ten sam kształt, którym otwiera się lockup. Miara grota modułu jest zapisana jako ułamek miary grota środowiska, żeby zależność między strefami trzymała się przy każdej zmianie tej pierwszej.

### `.cd-strona .cd-wejdz--grot`

Karta środowiska: wejście niesie sam grot, bez obrysu i tła. Kontrolką zostaje, bo karta nie jest przyciskiem — bez tego nie dałoby się wejść z klawiatury — ale nie nosi już chromu przycisku bibliotecznego.

### `margin: 0 var(--dn-od-2) var(--dn-od-2) auto;`

Lewy margines samoczynny dosuwa grot do prawej krawędzi wiersza metryki; skrót `margin` go zerował i grot wracał pod tekst.

### `.cd-strona .cd-wejdz--grot:focus-visible`

Grot nie ma obrysu przycisku, na którym obwódka fokusu mogłaby usiąść — siada na samym znaku, parą żetonów fokusu z biblioteki.

### `.cd-kafel-grot`

Kafel modułu jest przyciskiem sam w sobie, więc grot stoi w nim jako znak, nie kontrolka. W spoczynku znak i kropka są bezbarwne — barwę marki bierze dopiero pod wskaźnikiem; to odróżnia moduł od środowiska, gdzie kropka pali się zawsze. Reguła na `.kropka` bije wpisany w znacznik atrybut `fill`, bo dowolna reguła arkusza ma pierwszeństwo przed atrybutem prezentacji.

### `.cd-dodaj-pas`

Przycisk zakładający komponent własny stoi między strefą aplikacji a strefą Operatora — to on jest przejściem między jedną a drugą.

### `.cd-sekcja-menu`

Menu sekcji stoi w linii jej tytułu, wyrównane do środka wiersza — `align-items: baseline` nagłówka sekcji podciągało ikonę ponad tekst.

### `.cd-wlasny-menu`

Wyśrodkowanie rozciągnięciem, nie przesunięciem — `transform` zabierał etykiecie przycisku kotwicę w oknie widoku. Tak samo w wierszu wykazu.

### `.cd-wlasne-pusto`

Stan pusty strefy własnej: pokazuje się dopiero wtedy, gdy Operator nie ma ani jednego komponentu. Warunek liczy sam arkusz — bez skryptu.

### `.cd-listwa-poz:focus-visible`

Pierścień fokusu rysuje się do wewnątrz, na czterech krawędziach. Poprzednio stan fokusu niósł sam `--dn-cien-sygnal` (krycie 0,18–0,22) — obrys o kontraście 1,2 : 1 wobec listwy, praktycznie niewidoczny, a przy krawędzi listwy jeszcze przycinany.

### `—`

Liczba kolumn stref 1 i 2 wynika z miejsca na płótnie; przy otwartym panelu samouczka płótno jest o 320 px węższe, więc progi przesuwają się o tę samą wartość. Karta środowiska nie schodzi poniżej ~260 px — niżej tytuł środowiska łamie się w połowie wyrazu. Liczba powtórzeń w `repeat()` musi być literałem: `var()` w tym miejscu jest nieprawidłowe i deklaracja przepada.

### `@media (max-width: 1180px)`

Cztery kolumny to komplet wyboru, więc trzymają się do progu, poniżej którego karta przestaje mieścić nazwę środowiska. Panel samouczka zwęża płótno o 320 px, więc jego otwarcie przesuwa progi o tę samą wartość.

### `.sta-menu-poz > .ptaszek`

Menu tego okna ─────────────────────────────────────────────────────── Pozycja z zaznaczeniem ma ikonę po lewej, a znacznik stanu przy prawej krawędzi. Kolejność w dokumencie jest teraz zgodna z kolejnością widzianą, więc czytnik ekranu podaje nazwę przed znacznikiem, nie odwrotnie.

### `.cd-tresc--modul`

PŁÓTNO KARTY MODUŁOWEJ ──────────────────────────────────────────────── Karta modułu wskazuje własne płótno: wybranie karty zdejmuje pulpit Centrum. Prototyp niesie tu nagłówek warsztatu i odnośnik do prototypu modułu — układ warsztatu należy do tamtego okna, nie do Centrum dowodzenia.

### `.cd-strona .dn-karta-widoku`

KARTY PEŁNE — kwadratowe rogi (zakres: centrum, klasa `cd-strona`) -------------------------------------------------------------------------- Zaokrąglenie przy grubym akcencie `border-left` (wstęga sygnałowa) i cienkim obrysie reszty dawało skośny szew na rogu — widoczny artefakt. Pełna karta zdejmuje go u źródła, nie ruszając warstwy wspólnej w pozostałych oknach.

### `—`

Karty pełne zdejmują „klapki" łączące zaokrągloną zakładkę z panelem — pseudo-elementy ::before/::after zostawały jako ciemne prostokąty (widmo na karcie). Przy kwadratowej zakładce są zbędne. Zakres: centrum.

### `.cd-strona .dn-panel-drzewo .dn-obszar-pozycja`

Sesje w drzewie Projektów wyrównane do nazwy projektu (lewa linia) — wiersz niósł 12 px wcięcia, przez co bullet stał w prawo od nagłówka grupy.

### `.cd-strona .dn-karty-pasmo`

Jeden rozmiar ikon w pasie kart. Kontrolki (zamknij, „+", powłoki, narzędzia) szły na 20 px (--kart-ikona-sterowania), a zakładki i grot na 16 px — stąd rozjazd. Nadpisanie zmiennej sprowadza wszystkie do 16 px. Zakres: centrum.

### `.cd-strona .dn-karty-pasmo svg`

Jeden rozmiar ikony w całym oknie — 18 px, jak na wstążce głównej. Kontrolki pasa kart i wstążek okien szły na 16 px, przez co ikony zamknięcia, grotu i plakietki były drobniejsze od pozostałych.

### `.cd-strona .dn-karty-pasmo .dn-nrz-menu:hover`

Podwójne podświetlenie w pasie kart: tło hover dostawał i przycisk, i jego opakowanie (span .dn-nrz-menu / .dn-karty-wykaz / .dn-karty-powloki / .dn-karty-dodaj). Dwie warstwy rgba(255,255,255,.06) sumowały się do ~.12. Podświetla się wyłącznie przycisk. Zakres: centrum.

### `.cd-strona .dn-obszar-pasek`

WSTĄŻKA OKNA — pasek narzędzi każdego okna (lewego i roboczego) -------------------------------------------------------------------------- Lewe okno (sesje i projekty) oraz okno robocze są osobnymi oknami, więc każde ma własną wstążkę — tak jak wstążka główna aplikacji: zaokrąglona, w barwie panelu, z ikonami właściwymi dla tego okna.

### `margin: var(--dn-od-2);`

Odsunięcie ze wszystkich stron — bez marginesu górnego wstążka stykała się z pasem kart i oba zlewały się w jeden ciemny blok bez kształtu.

### `border: 1px solid var(--cd-obrys-marki);`

Wstążka nie jest osobną powierzchnią — ma tło okna, a wyodrębnia ją wyłącznie cienka ramka. Barwa `--dn-obrys` w motywie ciemnym równa się barwie wnętrza okna (#272A30), więc ramka na niej byłaby niewidoczna.

### `.cd-strona .dn-obszar-pasek::-webkit-scrollbar`

Wstążka nie jest polem przewijania: pole przewijania zamyka dymki przycisków w pasie własnej wysokości i daje wstążce pionowy pasek przesuwania. Komplet przycisków lewego okna zajmuje 238 px i mieści się w domyślnej szerokości okna (280 px).

### `overflow-x: auto;`

Komplet przycisków wstążki jest szerszy niż wstążka przy najmniejszej szerokości panelu (238 px zawartości wobec 201 px pola), więc pas przewija się w poziomie zamiast ucinać ostatni przycisk. Sam pasek przesuwania jest schowany — przewijanie idzie kółkiem i gestem. Dymki tego nie dotyczy: stoją w oknie widoku, poza polem przycinania wstążki.

### `.cd-strona .dn-obszar-pasek .dn-btn-ikona`

Przycisk wstążki okna jest węższy niż na wstążce głównej — okno boczne ma 280 px szerokości i przy 40 px na przycisk ikony nie mieściły się w rzędzie.

### `.cd-strona .dn-obszar-panel--boczny .dn-obszar-tresc`

Pozycje sesji i projektów zaczynają się na wysokości liter znaku DANACO. Pudełko `.cd-marka-lockup` stoi na 217 px, ale rysunek liter zaczyna się dopiero 33 px niżej — równanie do pudełka stawiało wykazy nad znakiem.

### `.cd-uchwyt`

UCHWYT ZMIANY SZEROKOŚCI OKIEN BOCZNYCH -------------------------------------------------------------------------- Lewe okno i okno samouczka mają regulowaną szerokość. Uchwyt jest wąski wizualnie, ale ma szerszą strefę chwytu, żeby dało się go złapać myszą.

### `.cd-strona .dn-obszar-panel--boczny`

ZWINIĘCIE LEWEGO OKNA I PODGLĄD PO NAJECHANIU NA KRAWĘDŹ -------------------------------------------------------------------------- Zwinięte okno nie znika z układu (jak przy `hidden`), tylko traci szerokość — dzięki temu zostaje krawędź, po której najechaniu okno rozwija się samo na czas podglądu i zwija z powrotem po zejściu kursorem.

### `.cd-strona .cd-kafel .dn-kafel-ikona`

Znak komponentu stoi bez płytki i bez ramki — samo godło. `.dn-kafel-ikona` z warstwy wspólnej niesie tło `--dn-powierzchnia` i obrys `--dn-obrys`; oba są tu zdejmowane. Pojemnik zachowuje 32 px z `.cd-kafel .dn-kafel-ikona` wyżej, żeby odstęp od brzegu kafla i od tekstu pozostał nietknięty, a znak wypełnia go w całości.

### `.cd-strona .dn-karty-lista ~ .dn-karty-dodaj::after`

Kreska przy „+” równo z krawędzią ostatniej karty. Kontrolka „+” odsuwa się od zbioru kart o `--kart-odstep-sterowania`, a jej kreska cofa się o tę samą wartość — czyli zatrzymuje się w połowie odstępu, nie na krawędzi karty. Cofnięcie o pełny odstęp stawia ją dokładnie tam, gdzie kończy się ostatnia karta.

### `.cd-strona .dn-karta-widoku[aria-selected='true']::before`

Przegrody pasa kart mają w tym oknie jedną formę i jedną barwę — w każdym pasie: w oknie roboczym, w oknie bocznym i w panelu samouczka. Barwa: `--dn-obrys` w motywie ciemnym ma wartość #272A30, czyli dokładnie tę samą co tło karty bieżącej i wnętrze okna. Kreska na niej ma kontrast 1,00 : 1 i jest niewidoczna z definicji, nie z powodu ustawień; przegrody idą więc na `--dn-obrys-mocny`. Forma: włos kreski wysokości `--kart-przegroda-wys`, pośrodku wysokości pasa. Krawędzie karty bieżącej szły z warstwy wspólnej na pełną wysokość karty i były jedynymi kreskami pasa innego rodzaju niż pozostałe — dwie formy przegrody w jednym rzędzie czytają się jako usterka, nie jako wyróżnienie.

### `.cd-strona .dn-karty-pasmo > .dn-karty-powloki::before`

Menu otwartych kart i menu okien roboczych to dwie osobne kontrolki, stykające się bokami na początku pasa. Każda inna szczelina pasa niesie przegrodę; ta jedna jej nie miała i obie kontrolki zlewały się w jeden blok.

### `.cd-strona .dn-karty-pasmo > .dn-karty-dodaj::after`

Kontrolka zwinięcia w pasie samouczka jest przyciskiem z etykietą najechania, a etykieta i przegroda idą na tym samym `::after`. Wygrywa przegroda, ale zostają po etykiecie `opacity: 0`, ramka, wyściółka i cień — kreska wychodzi niewidoczna i osiemnastokrotnie za szeroka. Reguła zdejmuje te pozostałości.

### `.cd-strona .dn-obszar-panel--boczny`

Krawędzie okien stoją w pionie krawędzi wstążki głównej: okna nie noszą marginesów ujemnych, a uchwyt szerokości okna ukrytego nie zajmuje miejsca.

### `.cd-strona .sta-menu-wiersz > .dn-btn-ikona`

Wykaz kart grupy jest jedynym miejscem, z którego kartę zgrupowaną można zamknąć bez wchodzenia w nią, więc jego wiersz niesie dwie czynności obok siebie: wybór karty i jej zamknięcie (`.sta-menu-wiersz`, biblioteka). Biblioteka ujawnia czynność dodatkową dopiero najechaniem — tutaj zamknięcie stoi w wierszu na stałe, bo ukryte byłoby tą samą nieobecnością, którą wiersz ma usunąć. Barwa metadanych oddaje pierwszeństwo nazwie karty; przy najechaniu i przy fokusie znak przechodzi w pełną barwę tekstu.

### `.cd-strona .sta-menu-wiersz > .dn-btn-ikona > svg`

Jeden rozmiar ikony w całym oknie — 18 px, jak zamknięcie karty w pasie kart. Panel menu przenosi się na czas otwarcia do warstwy okna, więc wychodzi spod reguły pasa kart wyżej i bierze ten rozmiar tutaj.

### `—`

DYMEK — kotwiczony do prostokąta wyzwalacza -------------------------------------------------------------------------- Dymek warstwy wspólnej stoi bezwzględnie względem wyzwalacza, więc należy do pola przycinania panelu, a panele okien przycinają zawartość. Przy skrajnym przycisku połowa napisu zostawała za krawędzią okna, a sam dymek podnosił zasięg przewijania wstążki z 38 do 61 px i dawał jej poziomy pasek. Umocowanie do okna widoku wyjmuje dymek z obu tych pól. Położenia nie da się przy tym wziąć z układu — `position: fixed` z `inset: auto` czyta miejsce spoczynkowe, które w każdym pojemniku wypada inaczej, przez co dymek odjeżdżał od ikony. Współrzędne podaje więc mechanizm okna w `--cd-dymek-x` i `--cd-dymek-y`, liczone z prostokąta wyzwalacza i przyciśnięte do widoku.

### `.cd-strona [data-etykietka]::after`

Umocowanie do okna widoku obowiązuje na KAŻDYM dymku, także spoczywającym: dymek pozycjonowany względem wyzwalacza wchodzi do pola przewijania swojego pojemnika i sam z siebie daje mu poziomy pasek, choć jest niewidoczny. Spoczywający stoi poza widokiem, bo i tak ma zerowe krycie. Współrzędne dostaje wyłącznie wyzwalacz bieżący. Zmienne własne dziedziczą się w dół drzewa, więc bez tego zawężenia dymek dziecka brałby współrzędne przodka i lądował po drugiej stronie okna.

### `.cd-strona [data-etykietka][data-dymek='tak']::after`

`right` musi zostać zdjęte razem z `left`. Warstwa wspólna przypina dymki pasa sterowania i ostatniej grupy narzędzi prawą krawędzią; przy umocowaniu do widoku taki dymek dostawał obie krawędzie naraz i rozciągał się od ikony po brzeg okna — 1726 px zamiast własnej szerokości.

### `.cd-strona .dn-narzedzia > .dn-narzedzia-grupa .dn-etykietka[data-dymek='tak']::after`

Skrajne grupy narzędzi przypinają dymek krawędzią grupy — `left: 0` w grupie pierwszej, `right: 0` w ostatniej — selektorem o swoistości wyższej niż kotwica, bo z `:first-child` i `:last-child`. Dymek dostawał przez to współrzędną z mechanizmu i krawędź z biblioteki naraz: pierwszy przycisk wstążki stawiał etykietę przy brzegu ekranu, pod szyną nawigacji.

### `.cd-strona [data-etykietka]:not([data-dymek='tak']):hover::after`

Dymek bez współrzędnych nie ma prawa się pokazać. Warstwa wspólna ujawnia go samym najechaniem, a współrzędne nadaje mechanizm okna — między jednym a drugim jest szczelina: gdy wyzwalacz pojawia się pod nieruchomym wskaźnikiem albo wjeżdża pod niego przewijaniem, najechanie zachodzi bez zdarzenia wskaźnika. Dymek stawał wtedy w miejscu spoczynkowym, które w pasie sterowania wypada przy prawej krawędzi okna, a gdzie indziej poza ekranem — stąd etykiety raz odległe od ikony, raz nieobecne.

### `.cd-strona .cd-siatka-srodowisk > .dn-karta-srodowiska`

Warstwa animacji karty środowiska ──────────────────────────────────── Godło zostaje na swoim miejscu; animacja jest tłem karty, nie elementem układu. Pole animacji zajmuje prawe 40 % szerokości karty na pełnej jej wysokości, a sama animacja stoi w tym polu pośrodku — w pionie i w poziomie. Rozmiar animacji liczy się od szerokości pola (`cqw`), więc karta zwęża się razem z płótnem, a animacja nadal mieści się w polu. Karta przycina zawartość do własnego łuku narożnika: przy `overflow: visible` animacja wylewałaby się poza narożniki i poza obrys karty.

### `opacity: 0.11;`

Animacja jest tłem pod tekstem: przygaszenie trzyma kontrast tytułu, opisu i etykiety „Wejdź" ponad progiem WCAG AA także w najjaśniejszej klatce pętli. Wartość wynika z pomiaru najjaśniejszego piksela pod tekstem.

### `.cd-strona .dn-obszar-pasek .dn-pasek-okna`

Narzędzia karty bieżącej na wstążce okna ───────────────────────────── Pod pasmem kart stoi jeden rząd narzędzi. Osobnego pasa modułu nad wstążką nie ma: narzędzia karty bieżącej weszły do wstążki okna roboczego i noszą jej miarę — przycisk 32 px, znak 18 px, kreska rozdzielająca 18 px. Klasa `dn-pasek-okna` została na grupie wyłącznie jako zaczep mechanizmu (`pasek-okna.js`): to on gasi narzędzia na karcie Centrum dowodzenia, która jest pulpitem okna, a nie modułem, i zwija nadmiar pod „…" przy zwężeniu. Deklaracje pasa — własne tło, kreska domykająca, wysokość i wyściółka — są tu wygaszone, bo grupa nie jest już pasem.

### `flex: 1 1 auto;`

Grupa bierze wolne miejsce wstążki: mechanizm zwijania mierzy zapas względem jej prawej krawędzi, więc bez rozciągnięcia zwinąłby wszystko.

### `.cd-strona .dn-obszar-pasek .dn-pasek-okna-odstep`

Przycisk „…" stoi zaraz za ostatnim narzędziem, nie przy prawej krawędzi wstążki — inaczej rozdzielałby narzędzia od menu widoku.

### `.cd-strona .dn-karta-srodowiska:hover .cd-karta-anim`

Kafel pod wskaźnikiem ujawnia animację. Metryka i zdanie opisowe biegną po niej, więc samo rozjaśnienie tła zabrałoby im czytelność — razem z animacją wzmacnia się dlatego barwa pisma, z tekstu drugorzędnego na pierwszorzędny. Obie wartości dobrane pomiarem kontrastu na najjaśniejszych klatkach pętli: przy 0.5 i piśmie pierwszorzędnym najsłabszy kontrast opisu to 4,74 : 1.

### `.cd-panel-stopka`

Stopka lewego okna ──────────────────────────────────────────────────── Dwie pozycje poza wykazem sesji i projektów: sprawdzenie wydania na serwerze oraz odsyłacz do strony wydań, otwierany w przeglądarce. Stopka stoi przy dolnej krawędzi okna, poza polem przewijania obu zakładek.

## Mechanizm `zasoby/okna/centrum-dowodzenia.js`

### `(function ()`

══════════════════════════════════════════════════════════════════════════ SKRYPT LOKALNY — Centrum dowodzenia v2. Rozszerza zachowania powłoki (menu, motyw, TRYBY) obsługiwane przez wspolne.js / prototyp.js / stanowisko.js. ══════════════════════════════════════════════════════════════════════════

### `qq('[data-przelacz-samouczek]').forEach(function (b)`

Stan samouczka niosą dwie kontrolki: przycisk pasa narzędzi (`aria-pressed`) i pozycja menu widoku (`aria-checked`). Przełącznik panelu ustawia tylko pierwszą z nich, więc drugą dostrajamy tu — inaczej ptaszek w menu kłamie.

### `—`

SAMOUCZEK · panel po prawej ───────────────────────────────────────── Przełączanie panelu obsługuje `zasoby/okna/zakladki-paneli.js` — wspólnie dla wyzwalacza w pasie narzędzi i dla przycisku „Otwórz samouczek" w wierszu pierwszego wejścia (oba noszą `data-przelacz-samouczek`). Tu nie dublujemy nasłuchu: dwa nasłuchy na tym samym kliknięciu zamykały panel w tej samej chwili, w której go otwierały.

### `if (e.target.closest('.cd-karta-menu'))`

Menu operacji leży na karcie, ale nie jest wejściem do środowiska — bez tego warunku każde otwarcie menu otwierało też środowisko.

### `var aod = q('#cd-aod-przelacz');`

Stan włączenia niesie sam przycisk (`aria-pressed`) — napis obok niego powtarzał tę samą informację słowem i był jedyną treścią w listwie, która nie była nazwą czynności.

### `(function ()`

══════════════════════════════════════════════════════════════════════════ OKNA BOCZNE — regulowana szerokość, zwijanie, podgląd po najechaniu ══════════════════════════════════════════════════════════════════════════

### `function zbudujUchwyt(panel, strona, zmienna, minSzer, maxSzer)`

Uchwyty zmiany szerokości ───────────────────────────────────────── Uchwyt jest elementem układu, nie nakładką — dzięki temu nie zasłania treści i sam trzyma się między oknami przy każdej szerokości okna.

### `u.setAttribute('aria-valuemin', String(minSzer));`

Ogniskowalny separator jest kontrolką o wartości — bez zakresu i wartości bieżącej czytnik ekranu nie ma czego odczytać.

### `D.documentElement.style.setProperty(zmienna, w + 'px');`

Szerokość idzie żetonem, nie stylem w linii: styl w linii bije regułę zwinięcia i okno nie dawało się zwinąć po zmianie szerokości.

### `if (lewy)`

Zwijanie lewego okna ────────────────────────────────────────────── Warstwa wspólna zwija panel atrybutem `hidden`, czyli usuwa go z układu — wtedy nie ma po czym najechać, żeby go podejrzeć. Tutaj zwinięcie odbiera szerokość, a okno zostaje w układzie. Przechwytujemy zdarzenie w fazie przechwytywania, zanim dojdzie do obsługi wspólnej.

### `var krawedz = D.createElement('div');`

Podgląd po najechaniu na lewą krawędź ─────────────────────────── Strefa czuła stoi przy krawędzi ekranu i działa wyłącznie wtedy, gdy okno jest zwinięte. Rozwinięcie na podgląd znika, gdy kursor opuści zarówno okno, jak i strefę.

### `(function ()`

══════════════════════════════════════════════════════════════════════════ WEJŚCIE DO PRZEDSIONKA ŚRODOWISKA Obsługa wspólna pokazywała komunikat, ale nie przechodziła do przedsionka — kafel środowiska nie prowadził donikąd. Nazwy plików są pisane małymi literami, a nazwa środowiska w znaczniku wielkimi, stąd `toLowerCase()`. ══════════════════════════════════════════════════════════════════════════

### `(function ()`

Kotwica dymka. Dymek jest umocowany do okna widoku, żeby wyjść poza pole przycinania paneli; współrzędne muszą więc pochodzić z prostokąta wyzwalacza, bo miejsce spoczynkowe pseudoelementu wypada w każdym pojemniku inaczej. Dymek staje pod wyzwalaczem, wyśrodkowany; przy krawędzi okna jest do niej przyciskany, a przy dolnej krawędzi przeskakuje nad wyzwalacz. Zmienne własne dziedziczą się w dół drzewa, więc same nie wystarczą: dymek dziecka brałby współrzędne przodka, który był pod wskaźnikiem wcześniej. Nośnikiem jest dlatego znacznik `data-dymek` na wyzwalaczu bieżącym — poza nim obowiązuje umocowanie biblioteczne.

### `function miara(el)`

Szerokości pseudoelementu nie da się odczytać wprost, ale da się złożyć: napis mierzony krojem dymka plus jego wyściółka i ramka. Krój bierze się z `::after`, nie z wyzwalacza — przycisk niesie własny, większy stopień pisma, przez który dymek wychodził szerszy, niż jest, i docisk do krawędzi odsuwał go od ikony nawet tam, gdzie mieścił się bez przesunięcia.

### `var lewaZapamietana = null;`

Szyna stoi nieruchomo, a jej miara jest potrzebna przy każdym najechaniu — wystarczy zmierzyć ją raz i odświeżyć, gdy okno zmieni rozmiar.

### `var lewaStrefa = strefaLewa();`

Lewa strefa bezpieczna omija pionową szynę nawigacji: dymek pierwszego przycisku wstążki dosuwał się do krawędzi okna i wchodził na szynę, która stoi wyżej w układzie. Granicę odsuwamy za prawą krawędź szyny.

### `var wysokosc = m.wysokosc;`

Przy dolnej krawędzi okna dymek nie ma dokąd opaść — staje nad wyzwalaczem. Wysokość bierze się z `::after`, bo stopień pisma dymka zmienia się razem z żetonem.

### `document.addEventListener('pointermove', function (e)`

Domknięcie szczeliny między najechaniem a zdarzeniem wskaźnika. Wyzwalacz, który pojawia się pod nieruchomym wskaźnikiem albo wjeżdża pod niego przewijaniem, dostaje stan najechania bez `pointerover`. Ruch wskaźnika nadrabia współrzędne, a warunek na znaczniku sprawia, że pełny rachunek wykonuje się raz na wyzwalacz, nie przy każdym ruchu.
