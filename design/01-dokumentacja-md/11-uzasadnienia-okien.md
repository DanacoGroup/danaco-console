# Uzasadnienia reguł okien

Powody, dla których reguły okien są takie, a nie inne. Arkusze i mechanizmy
niosą sam kod: limit gęstości komentarzy wynosi 250 znaków na 1000 wierszy
i nie pomieściłby tych treści. Zapis prowadzony jest per okno, w kolejności
wprowadzania okien do budowy.

Powiązane: [05 — style CSS](05-styles-css.md), [03 — widok](03-design-view.md),
[rejestr terenów](../../prowadzenie/rejestr-terenow.md).

---

## Centrum dowodzenia

Teren `centrum-poprawki`. Pliki: `05-okna/przeplyw/centrum-dowodzenia.html`,
`zasoby/okna/centrum-dowodzenia.css`, `zasoby/okna/centrum-dowodzenia.js`.

Punkty łamania układu mierzą **szerokość płótna** (`@container plotno`), nie
szerokość ekranu. Płótno przełącza siatki stref na dwie kolumny poniżej 820 px
i na jedną poniżej 500 px; belka pierwszego wejścia i listwa ustawień idą tym
samym progiem. Otwarcie panelu samouczka zwęża płótno i przełącza układ samo,
bez osobnego warunku.

Wyjątkiem jest wyściółka samego `.cd-tresc` — element nie może odpytywać
pojemnika, który sam deklaruje, więc ta jedna reguła zostaje przy progu ekranu.

### Arkusz

**`:root`**

DANACO CONSOLE — CENTRUM DOWODZENIA: ZŁOŻENIE OKNA ---------------------------------------------------------------------------- Kompetencja: układ karty Centrum dowodzenia — nagłówek, wiersz pierwszego wejścia oraz cztery strefy wyboru (środowiska, komponenty aplikacji, komponenty własne Operatora, listwa ustawień) i panel samouczka. Wygląd komponentów pochodzi z arkuszy kompetencyjnych; tutaj stoi wyłącznie ich złożenie w tym oknie.

Wymaga: `zetony/zetony.css`, `zetony/ruch.css`, `css/komponenty.css`, `stany.css`, `rama.css`, `karty-okna.css`, `panel-sesji.css`, `pasek-okna.css`, `izolacja.css`, `okno-robocze.css`. Mechanizm: `okna/centrum-dowodzenia.js`.

Punkty łamania mierzą szerokość płótna (`@container plotno`), nie szerokość ekranu: panel samouczka zwęża okno główne o 320 px.

**`.cd-strona`**

── Włos marki Obrys wstążek i okien niesie barwę znaku zmieszaną z obrysem podstawowym. Udział barwy znaku jest dwustopniowy: wyższy dla krawędzi, niższy dla tła karty bieżącej, gdzie barwa leży pod tekstem i musi zostawić mu kontrast. Mieszanie w przestrzeni oklab daje w obu motywach równy stopień odejścia od szarości, czego mieszanie kanałowe nie zapewnia. Wartości udziału stoją tu, a nie w `zasoby/zetony/`, bo rozstrzygnięcie obejmuje jedno okno; przy rozszerzeniu na pozostałe przenieść do palety.

**`.cd-tresc:not([hidden])`**

Płótno karty. Warunek `:not([hidden])` jest konieczny: `display` zadeklarowany przez autora wygrywa z `[hidden]` z arkusza przeglądarki, a wtedy karta odłożona zostawałaby na ekranie razem z kartą bieżącą.

**`.cd-tresc > *`**

Wejście stref w kolejności ważności. `backwards`, nie `both`: animacja wypełniająca w przód pozostaje czynna i czyni ze strefy kontekst nakładania, co przycina panele rozwijane w jej wnętrzu.

**`.cd-start`**

Belka samouczka i listwa ustawień domykają się na siatce środowisk, nie na krawędzi płótna: kończą się równo z drugą kartą przy czterech kolumnach i z pierwszą przy dwóch. Obie miary daje jedno wyrażenie — dwie kolumny wraz z przerwą między nimi mierzą dokładnie tyle, co jedna kolumna z układu dwukolumnowego, czyli połowę siatki pomniejszoną o pół odstępu. `min-width` chroni listwę: jej pięć pozycji potrzebuje więcej miejsca, niż zostaje przy dwóch kolumnach, a treść nie ma prawa wyjść poza obrys.

**`.cd-listwa`**

Rząd kontrolek czyta się jako jeden zestaw, więc pole listwy nie zwęża się poniżej połowy siatki z własnej woli. Gdy jednak okno zejdzie tak nisko, że pięć kontrolek nie mieści się w tej połowie, rząd łamie się na dwa wiersze. Podłoga na szerokości treści byłaby tu błędem: `min-width` bije `max-width`, więc poniżej 480 px listwa wychodziła poza pole treści o ponad 140 px.

**`@media (max-width: 860px)`**

Jedna kolumna — karta zajmuje całą szerokość, więc belka też. Próg 860 px nie należy do skali `--dn-bp-*`; belka musi łamać się dokładnie tam, gdzie łamie się siatka środowisk, a ta stoi na progach 1180 i 860 px sprzed tej pracy. Przeniesienie siatki na skalę systemową rozstrzyga się osobno — do tego czasu belka idzie za siatką, nie za skalą.

**`.cd-sekcja-tytul::before`**

Waga strefy zapisana atrybutem `data-waga`; nośnikiem jest kreska sygnałowa przed tytułem, w trzech stopniach krycia.

**`.cd-karta-akcent::before`**

Wstęga sygnałowa na górnej krawędzi kończyła się przed łukiem narożnika: pseudoelement kotwiczy się do pola wyściółki, a przycięcie liczy promień pomniejszony o obrys, więc jej koniec wystawał poza łuk i czytał się jak kreska urwana w rogu. Stan karty niesie sam obrys — na czterech krawędziach o jednej geometrii.

**`.cd-strona .cd-kafel-tresc .dn-kafel-etykieta`**

Znak modułu niesie barwę znaku środowiska — obie strefy strony głównej mają ten sam stopień kontrastu. Nazwa modułu schodzi z półgrubej na zwykłą: krój treści (IBM Plex Sans) czyta się przy tej samej wadze ciężej niż krój nagłówkowy karty środowiska (Space Grotesk).

**`.cd-wlasny-menu`**

Wyśrodkowanie rozciągnięciem, nie przesunięciem — `transform` zabierał etykiecie przycisku kotwicę w oknie widoku. Tak samo w wierszu wykazu.

**`.cd-listwa-poz:focus-visible`**

Pierścień fokusu rysuje się do wewnątrz, na czterech krawędziach. Poprzednio stan fokusu niósł sam `--dn-cien-sygnal` (krycie 0,18–0,22) — obrys o kontraście 1,2 : 1 wobec listwy, praktycznie niewidoczny, a przy krawędzi listwy jeszcze przycinany.

**`@media (max-width: 1180px)`**

Cztery kolumny to komplet wyboru, więc trzymają się do progu, poniżej którego karta przestaje mieścić nazwę środowiska. Panel samouczka zwęża płótno o 320 px, więc jego otwarcie przesuwa progi o tę samą wartość.

**`—`**

Karty pełne zdejmują „klapki" łączące zaokrągloną zakładkę z panelem — pseudo-elementy ::before/::after zostawały jako ciemne prostokąty (widmo na karcie). Przy kwadratowej zakładce są zbędne. Zakres: centrum.

**`—`**

(reguła zdjęta — ::after niesie prawą krawędź karty bieżącej).

**`.cd-strona .dn-karty-pasmo`**

Jeden rozmiar ikon w pasie kart. Kontrolki (zamknij, „+", powłoki, narzędzia) szły na 20 px (--kart-ikona-sterowania), a zakładki i grot na 16 px — stąd rozjazd. Nadpisanie zmiennej sprowadza wszystkie do 16 px. Zakres: centrum.

**`overflow-x: auto;`**

Komplet przycisków wstążki jest szerszy niż wstążka przy najmniejszej szerokości panelu (238 px zawartości wobec 201 px pola), więc pas przewija się w poziomie zamiast ucinać ostatni przycisk. Sam pasek przesuwania jest schowany — przewijanie idzie kółkiem i gestem. Dymki tego nie dotyczy: stoją w oknie widoku, poza polem przycinania wstążki.

**`.cd-strona .cd-kafel .dn-kafel-ikona`**

Znak komponentu stoi bez płytki i bez ramki — samo godło. `.dn-kafel-ikona` z warstwy wspólnej niesie tło `--dn-powierzchnia` i obrys `--dn-obrys`; oba są tu zdejmowane. Pojemnik zachowuje 32 px z `.cd-kafel .dn-kafel-ikona` wyżej, żeby odstęp od brzegu kafla i od tekstu pozostał nietknięty, a znak wypełnia go w całości.

**`.cd-strona .dn-karty-lista ~ .dn-karty-dodaj::after`**

Kreska przy „+” równo z krawędzią ostatniej karty. Kontrolka „+” odsuwa się od zbioru kart o `--kart-odstep-sterowania`, a jej kreska cofa się o tę samą wartość — czyli zatrzymuje się w połowie odstępu, nie na krawędzi karty. Cofnięcie o pełny odstęp stawia ją dokładnie tam, gdzie kończy się ostatnia karta.

**`.cd-strona .dn-karta-widoku[aria-selected='true']::before`**

Przegrody pasa kart mają w tym oknie jedną formę i jedną barwę — w każdym pasie: w oknie roboczym, w oknie bocznym i w panelu samouczka.

Barwa: `--dn-obrys` w motywie ciemnym ma wartość #272A30, czyli dokładnie tę samą co tło karty bieżącej i wnętrze okna. Kreska na niej ma kontrast 1,00 : 1 i jest niewidoczna z definicji, nie z powodu ustawień; przegrody idą więc na `--dn-obrys-mocny`.

Forma: włos kreski wysokości `--kart-przegroda-wys`, pośrodku wysokości pasa. Krawędzie karty bieżącej szły z warstwy wspólnej na pełną wysokość karty i były jedynymi kreskami pasa innego rodzaju niż pozostałe — dwie formy przegrody w jednym rzędzie czytają się jako usterka, nie jako wyróżnienie.

**`.cd-strona .dn-karty-pasmo > .dn-karty-powloki::before`**

Menu otwartych kart i menu okien roboczych to dwie osobne kontrolki, stykające się bokami na początku pasa. Każda inna szczelina pasa niesie przegrodę; ta jedna jej nie miała i obie kontrolki zlewały się w jeden blok.

**`.cd-strona .dn-karty-pasmo > .dn-karty-dodaj::after`**

Kontrolka zwinięcia w pasie samouczka jest przyciskiem z etykietą najechania, a etykieta i przegroda idą na tym samym `::after`. Wygrywa przegroda, ale zostają po etykiecie `opacity: 0`, ramka, wyściółka i cień — kreska wychodzi niewidoczna i osiemnastokrotnie za szeroka. Reguła zdejmuje te pozostałości.

**`—`**

DYMEK — kotwiczony do prostokąta wyzwalacza -------------------------------------------------------------------------- Dymek warstwy wspólnej stoi bezwzględnie względem wyzwalacza, więc należy do pola przycinania panelu, a panele okien przycinają zawartość. Przy skrajnym przycisku połowa napisu zostawała za krawędzią okna, a sam dymek podnosił zasięg przewijania wstążki z 38 do 61 px i dawał jej poziomy pasek.

Umocowanie do okna widoku wyjmuje dymek z obu tych pól. Położenia nie da się przy tym wziąć z układu — `position: fixed` z `inset: auto` czyta miejsce spoczynkowe, które w każdym pojemniku wypada inaczej, przez co dymek odjeżdżał od ikony. Współrzędne podaje więc mechanizm okna w `--cd-dymek-x` i `--cd-dymek-y`, liczone z prostokąta wyzwalacza i przyciśnięte do widoku.

**`.cd-strona [data-etykietka]::after`**

Umocowanie do okna widoku obowiązuje na KAŻDYM dymku, także spoczywającym: dymek pozycjonowany względem wyzwalacza wchodzi do pola przewijania swojego pojemnika i sam z siebie daje mu poziomy pasek, choć jest niewidoczny. Spoczywający stoi poza widokiem, bo i tak ma zerowe krycie.

Współrzędne dostaje wyłącznie wyzwalacz bieżący. Zmienne własne dziedziczą się w dół drzewa, więc bez tego zawężenia dymek dziecka brałby współrzędne przodka i lądował po drugiej stronie okna.

**`.cd-strona [data-etykietka][data-dymek='tak']::after`**

`right` musi zostać zdjęte razem z `left`. Warstwa wspólna przypina dymki pasa sterowania i ostatniej grupy narzędzi prawą krawędzią; przy umocowaniu do widoku taki dymek dostawał obie krawędzie naraz i rozciągał się od ikony po brzeg okna — 1726 px zamiast własnej szerokości.

**`.cd-strona .dn-narzedzia > .dn-narzedzia-grupa .dn-etykietka[data-dymek='tak']::after`**

Skrajne grupy narzędzi przypinają dymek krawędzią grupy — `left: 0` w grupie pierwszej, `right: 0` w ostatniej — selektorem o swoistości wyższej niż kotwica, bo z `:first-child` i `:last-child`. Dymek dostawał przez to współrzędną z mechanizmu i krawędź z biblioteki naraz: pierwszy przycisk wstążki stawiał etykietę przy brzegu ekranu, pod szyną nawigacji.

**`.cd-strona [data-etykietka]:not([data-dymek='tak']):hover::after`**

Dymek bez współrzędnych nie ma prawa się pokazać. Warstwa wspólna ujawnia go samym najechaniem, a współrzędne nadaje mechanizm okna — między jednym a drugim jest szczelina: gdy wyzwalacz pojawia się pod nieruchomym wskaźnikiem albo wjeżdża pod niego przewijaniem, najechanie zachodzi bez zdarzenia wskaźnika. Dymek stawał wtedy w miejscu spoczynkowym, które w pasie sterowania wypada przy prawej krawędzi okna, a gdzie indziej poza ekranem — stąd etykiety raz odległe od ikony, raz nieobecne.

**`.cd-strona .cd-siatka-srodowisk > .dn-karta-srodowiska`**

── Warstwa animacji karty środowiska Godło zostaje na swoim miejscu; animacja jest tłem karty, nie elementem układu. Pole animacji zajmuje prawe 40 % szerokości karty na pełnej jej wysokości, a sama animacja stoi w tym polu pośrodku — w pionie i w poziomie. Rozmiar animacji liczy się od szerokości pola (`cqw`), więc karta zwęża się razem z płótnem, a animacja nadal mieści się w polu. Karta przycina zawartość do własnego łuku narożnika: przy `overflow: visible` animacja wylewałaby się poza narożniki i poza obrys karty.

**`opacity: 0.11;`**

Animacja jest tłem pod tekstem: przygaszenie trzyma kontrast tytułu, opisu i etykiety „Wejdź" ponad progiem WCAG AA także w najjaśniejszej klatce pętli. Wartość wynika z pomiaru najjaśniejszego piksela pod tekstem.

**`.cd-strona .dn-karta-srodowiska:hover .cd-karta-anim`**

Kafel pod wskaźnikiem ujawnia animację. Metryka i zdanie opisowe biegną po niej, więc samo rozjaśnienie tła zabrałoby im czytelność — razem z animacją wzmacnia się dlatego barwa pisma, z tekstu drugorzędnego na pierwszorzędny. Obie wartości dobrane pomiarem kontrastu na najjaśniejszych klatkach pętli: przy 0.5 i piśmie pierwszorzędnym najsłabszy kontrast opisu to 4,74 : 1.

**`(function ()`**

Kotwica dymka. Dymek jest umocowany do okna widoku, żeby wyjść poza pole przycinania paneli; współrzędne muszą więc pochodzić z prostokąta wyzwalacza, bo miejsce spoczynkowe pseudoelementu wypada w każdym pojemniku inaczej. Dymek staje pod wyzwalaczem, wyśrodkowany; przy krawędzi okna jest do niej przyciskany, a przy dolnej krawędzi przeskakuje nad wyzwalacz.

Zmienne własne dziedziczą się w dół drzewa, więc same nie wystarczą: dymek dziecka brałby współrzędne przodka, który był pod wskaźnikiem wcześniej. Nośnikiem jest dlatego znacznik `data-dymek` na wyzwalaczu bieżącym — poza nim obowiązuje umocowanie biblioteczne.

**`function miara(el)`**

Szerokości pseudoelementu nie da się odczytać wprost, ale da się złożyć: napis mierzony krojem dymka plus jego wyściółka i ramka. Krój bierze się z `::after`, nie z wyzwalacza — przycisk niesie własny, większy stopień pisma, przez który dymek wychodził szerszy, niż jest, i docisk do krawędzi odsuwał go od ikony nawet tam, gdzie mieścił się bez przesunięcia.

**`var wysokosc = m.wysokosc;`**

Przy dolnej krawędzi okna dymek nie ma dokąd opaść — staje nad wyzwalaczem. Wysokość bierze się z `::after`, bo stopień pisma dymka zmienia się razem z żetonem.
