# Rejestr terenów

Żywy wykaz terenów. Prowadzi go wyłącznie Prowadzący budowę, na gałęzi `main`.
Teren bez wpisu w tym rejestrze nie jest otwarty, a praca na nim nie zostanie
przyjęta. Zasady podziału opisuje [ustrój budowy](ustroj-budowy.md).

## Tereny otwarte

### droga-wejscia

Widok drogi wejścia na gotowym fundamencie klienta. Prototyp jest przyjęty
(pozycja 12 rejestru decyzji), rdzeń przechodzi całą drogę, warstwy połączenia
i protokołu stoją. Nic tego terenu nie blokuje.

| | |
|---|---|
| **Gałąź** | `teren/droga-wejscia` z `main` |
| **Wykaz plików** | `budowa/klient/src/wejscie/` — katalog powstaje w tym terenie |
| **Do czytania, bez zapisu** | `design/zasoby/okna/wejscie/`, `design/05-okna/przeplyw/przeplyw-wejscia.html`, `budowa/shared/contract.json`, `budowa/klient/src/polaczenie`, `budowa/klient/src/protokol` |
| **Poza terenem** | `budowa/server/`, `budowa/desktop/`, `budowa/shared/`, `design/`, `prowadzenie/`, `budowa/klient/src/polaczenie`, `budowa/klient/src/protokol` |

**Przedmiot.** Trzy etapy prototypu w działającym kliencie: łączenie z rdzeniem,
rejestracja i logowanie wraz z odzyskaniem konta, przygotowanie środowiska.
Odsłony niesie prototyp — nawiązywanie połączenia, przywracanie sesji, błąd
połączenia, logowanie, wstrzymanie po pięciu próbach, założenie konta Operatora,
potwierdzenie adresu, odzyskanie dostępu, ustawienie nowego hasła.

**Kompozycji nie wymyślasz — czytasz ją z prototypu.** Pozycja 12 rejestru
decyzji rozstrzyga, co prototyp wiąże, a co jest parametrem inżynierskim.

**Kryteria odbioru.**

1. Każdy z trzech etapów rozmawia z **żywym rdzeniem** — z przytoczoną
   odpowiedzią rdzenia dla `connection.hello`, `auth.register`, `auth.login`
   i `environment.enter`.
2. Rejestracja obsługuje **obie gałęzie pozycji 11**: przy nadajniku
   `pendingVerification: true` prowadzi do potwierdzenia listem, przy braku
   nadajnika kończy się wejściem hasłem i **nazywa niepotwierdzony adres**.
3. Tekst widoczny dla użytkownika stoi w jednym katalogu treści — poza nim zero
   łańcuchów. Wykazane pomiarem, tak jak w prototypie.
4. `tsc --noEmit` bez błędu, sprawdziany zdane — z przytoczonym wynikiem.
5. Warstwy połączenia i protokołu **nietknięte** — wykazane `git diff`.
6. Zero dotknięć DOM poza katalogiem tego terenu.
7. Odsłony błędu i wstrzymania są osiągalne w sprawdzianie, nie tylko opisane.
8. Rzecz wymagająca rozstrzygnięcia Właściciela wraca jako zgłoszenie.

### powloka-tauri

Powłoka okna dla Windows 11 w dwóch architekturach. Wariant natywny zniesiony
(pozycja 8), więc powłoka niesie interfejs i **nie niesie rdzenia**.

| | |
|---|---|
| **Gałąź** | `teren/powloka-tauri` z `main` |
| **Wykaz plików** | `budowa/desktop/` |
| **Do czytania, bez zapisu** | `budowa/klient/`, `budowa/shared/contract.json`, `design/05-okna/platformowe/instalator.html` |
| **Poza terenem** | `budowa/server/`, `budowa/klient/`, `budowa/shared/`, `design/`, `prowadzenie/` |

**Przedmiot.** Doprowadzić powłokę do postaci budowalnej na oba cele Windows 11,
zgodnej z hybrydą: okno, cykl życia, wskazanie serwera, brak rdzenia w pakiecie.

**Kryteria odbioru.**

1. Powłoka buduje się na `x86_64-pc-windows-gnu` **oraz** `aarch64-pc-windows-msvc`
   — z przytoczonym wynikiem obu przebiegów.
2. Pakiet **nie zawiera** rdzenia ani serwera narzędzi — wykazane wykazem
   zawartości pakietu, nie deklaracją.
3. Adres rdzenia pochodzi z nastawy budowania; **budowa wydaniowa nie niesie
   adresu serwera rozwojowego** — wykazane przeszukaniem gotowego pliku.
4. Powłoka nie stawia i nie wygasza rdzenia — w hybrydzie rdzeń stoi na serwerze.
5. Uprawnienia powłoki odmawiają domyślnie; każde nadane ma podany powód.
6. Rzecz wymagająca rozstrzygnięcia Właściciela wraca jako zgłoszenie.

## Zgłoszenia oczekujące na teren

Ustalenia z zamkniętych i biegnących terenów, które wykraczają poza ich zakres.
Każde zgłoszenie ma wskazany plik i wiersz. Zgłoszenie staje się terenem, gdy
Prowadzący je otworzy; do tego czasu jest wykazem, nie pracą.

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

### Sprawdziany drogi wejścia rozjechane z rdzeniem — gotowe do otwarcia

Trzy sprawdziany zawodzą od chwili przejęcia rdzenia. Pomiar wykazał, że nie są
usterką rdzenia — opisują zamiar porzucony i nikt ich za zmianą nie poprawił.

| Sprawdzian | Czego żąda | Co rdzeń robi |
|---|---|---|
| `TestRejestracjaBezKontaNadawczegoOdmawiaINieZakladaKonta` | odmowy i zera wierszy przy braku nadajnika | zakłada konto, stawia znacznik, wpuszcza hasłem — pozycja 11 rejestru decyzji |
| `TestNieudaneNadanieListuCofaRejestracje` | cofnięcia rejestracji przy nadajniku nieosiągalnym | do zmierzenia w terenie |
| `TestSkanowanieZUrzadzeniaOdmawiaNazwanie` | odmowy nazywającej brak | do zmierzenia w terenie |

Osobno: **droga bez poczty nie ma ani jednego sprawdzianu własnego**. Zachowanie
rozstrzygnięte pozycją 11 stoi dziś wyłącznie na komentarzu i na pomiarze
jednorazowym — pierwsza zmiana w bramce zniesie je bez niczyjej wiedzy.

Teren ma dla każdego z trzech sprawdzianów rozstrzygnąć pomiarem, czy zawodzi
sprawdzian, czy rdzeń, i poprawić tę stronę, która się myli — a nie tę, którą
łatwiej. Do tego założyć sprawdziany drogi bez poczty: rejestracja, wejście
hasłem, potwierdzenie adresu po ustawieniu nadajnika, zdjęcie znacznika.

### Droga SANE bez odpowiednika `bladWarstwyWia` — gotowe do otwarcia

Brak urządzenia jest już nazwany. Brak samego programu `scanimage` dalej wychodzi
odmową arsenału bez wskazania drogi obejścia, podczas gdy `bladWarstwyWia`
(`urzadzenia_skaner.go`) dla tej samej sytuacji na Windowsie podaje
`studio.ingest.queue.add`. Ta sama asymetria dotyczy `wykazSkanerow`: gałąź
Windows przekłada odmowę, gałąź Linux oddaje ją surową. Operator na Linuksie bez
`sane-utils` nie dowie się, że materiał da się wnieść inną drogą.

### Sprawdzian katalogu akcji szuka nieistniejącego katalogu — gotowe do otwarcia

`budowa/server/internal/store/katalog_akcji_test.go:177` szuka
`../../../client/src/ikony/zrodla`. Katalog klienta nazywa się `budowa/klient`,
a `budowa/client` nie istnieje w żadnej gałęzi. To jedyne niepowodzenie
pozostałe w całym module. Do rozstrzygnięcia wraz z pierwszym terenem widoku,
bo dotyczy źródeł ikon nowego klienta.

### Reguła odbioru wyprowadzona z pomiarów

`axe.run()` sam wywołuje dwa błędy 404 (`menu.css`, `ruch.css`), bo rozwiązuje
`@import` względem adresu dokumentu, a nie arkusza. Konsola przed wstrzyknięciem
axe jest pusta. Tych dwóch wpisów nie liczy się jako brudnej konsoli.

Pomiar w przeglądarce wymaga jawnego ustawienia `PLAYWRIGHT_BROWSERS_PATH` na
`/opt/ms-playwright` w poleceniu, a nie polegania na środowisku powłoki — powłoka
uruchomiona przed ustawieniem zmiennej jej nie widzi i pobiera przeglądarki
po raz drugi.

## Tereny zamknięte

| Nazwa | Gałąź | Rewizje | Kontrola |
|---|---|---|---|
| `brama-i-droga-wejscia` | `teren/brama-i-droga-wejscia` | `23a6b84` wyjątek powitania · `e41dd78` straże drogi bez poczty · `9a8ec0c` brak skanera | weryfikacja Prowadzącego pomiarem: bieg wymuszony `-count=1` 537 s — 2029 sprawdzianów, 1 niezdany wobec 4 zastanych; powitanie niepełne odpowiada wersją protokołu na żywym rdzeniu, `channel.add` z brakiem pola dalej odmawia; kontrakt nietknięty |
| `prototypy` | `teren/prototypy` | paczki instalatora i drogi wejścia | przyjęte przez Właściciela; weryfikacja Prowadzącego pomiarem: oba okna wczytują się bez błędu konsoli, zero łańcuchów widocznych poza katalogiem treści |
| `brama-i-droga-wejscia` | `teren/brama-i-droga-wejscia` | `23a6b84` · `e41dd78` · `9a8ec0c` | weryfikacja Prowadzącego pomiarem: 2029 sprawdzianów, 1 niepowodzenie zastane spoza terenu wobec 4 zastanych; wyjątek bramy w jednym miejscu; kontrakt nietknięty |
| `fundament-klienta` | `teren/fundament-klienta` | `d19bfeb` warstwa połączenia i protokołu | weryfikacja Prowadzącego pomiarem: kompilacja bez błędu, 17 sprawdzianów zdanych, rozmowa z żywym rdzeniem, generat bajtowo powtarzalny, zero dotknięć DOM, kontrakt nietknięty |
| `naprawy-rdzenia` | `teren/naprawy-rdzenia` | `060d5b7` naprawy i brama kontraktu | weryfikacja Prowadzącego pomiarem: 2022 zdane wobec 2003 zastanych, te same 4 niezdane, kontrakt nietknięty |
| `proba-prototypow` | `teren/prototypy` | `66e5ee0` przepływ wejścia · `61a5867` moduł Studio · `53bc3d5` odsyłacze | kontrola sesji nadzorującej wykonanie, weryfikacja Prowadzącego pomiarem |

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
