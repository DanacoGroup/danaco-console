# Rejestr terenów

Żywy wykaz terenów. Prowadzi go wyłącznie Prowadzący budowę, na gałęzi `main`.
Teren bez wpisu w tym rejestrze nie jest otwarty, a praca na nim nie zostanie
przyjęta. Zasady podziału opisuje [ustrój budowy](ustroj-budowy.md).

## Tereny otwarte

Fala 1 — fundament klienta. Cztery tereny; pierwszy blokuje trzy pozostałe, bo
one biorą z niego typy. Żaden nie zależy od prototypów.

### typy-kontraktu · **blokujący**

| | |
|---|---|
| **Gałąź** | `teren/typy-kontraktu` z `main` |
| **Wykaz plików** | `budowa/klient/src/kontrakt.ts`, `budowa/narzedzia-kontraktu/` |
| **Poza terenem** | `budowa/shared/contract.json` — źródło, tylko do odczytu; `budowa/server/` |

**Przedmiot.** Wytworzenie typów TypeScript z `budowa/shared/contract.json` oraz
polecenie wykrywające rozjazd między typami Go a TypeScript.

**Kryteria odbioru.**

1. Wytworzenie jest powtarzalne: dwa przebiegi dają plik identyczny co do bajta.
2. Każda z 1077 komend, 542 struktur i 363 wyliczeń ma odpowiednik w typach.
3. Zmiana w `contract.json` bez ponownego wytworzenia jest wykrywana poleceniem
   kończącym się kodem różnym od zera.
4. `tsc --noEmit` przechodzi na wytworzonym pliku.
5. Źródło kontraktu nietknięte — wykazane `git status`.

### warstwa-polaczenia

| | |
|---|---|
| **Gałąź** | `teren/warstwa-polaczenia` z `teren/typy-kontraktu` po jego zamknięciu |
| **Wykaz plików** | `budowa/klient/src/polaczenie/` |

**Przedmiot.** Gniazdo, koperta kontraktu, korelacja żądanie–odpowiedź,
wznowienie po zerwaniu, przeciwciśnienie.

**Kryteria odbioru.**

1. Żądanie otrzymuje swoją odpowiedź także wtedy, gdy w locie jest wiele żądań —
   wykazane próbą z dwudziestoma naraz.
2. Zerwanie połączenia w trakcie strumienia i powrót nie gubi zdarzeń —
   wykazane próbą z przerwaniem i porównaniem wykazu odebranych zdarzeń.
3. Koperta niezgodna z kontraktem jest odrzucana przed wysłaniem, nie przez rdzeń.
4. Warstwa nie dotyka drzewa dokumentu — wykazane brakiem odwołań do `document`
   i `window` w plikach terenu.

### warstwa-stanu

| | |
|---|---|
| **Gałąź** | `teren/warstwa-stanu` z `teren/typy-kontraktu` po jego zamknięciu |
| **Wykaz plików** | `budowa/klient/src/stan/` |

**Przedmiot.** Odbiór 75 zdarzeń kontraktu, jedno źródło stanu, brak stanu
w warstwie widoku.

**Kryteria odbioru.**

1. Każde z 75 zdarzeń ma obsługę albo jawne pominięcie wraz z powodem —
   wykazane wykazem zestawionym z kontraktem.
2. Stan zmienia się wyłącznie przez zdarzenia; brak zapisu stanu z widoku.
3. Warstwa nie dotyka drzewa dokumentu.

### warstwa-zakresu

| | |
|---|---|
| **Gałąź** | `teren/warstwa-zakresu` z `teren/typy-kontraktu` po jego zamknięciu |
| **Wykaz plików** | `budowa/klient/src/zakres/` |

**Przedmiot.** Kaskada czterech kondygnacji — aplikacja, projekt, sesja, karta —
na dwóch osiach: upoważnienie i dostęp. Dziedziczenie i nadpisanie w obie strony.

**Kryteria odbioru.**

1. Rozstrzygnięcie zakresu zgodne z pozycją 6 rejestru decyzji — wykazane tabelą
   przypadków obejmującą dziedziczenie, zawężenie i rozszerzenie.
2. Stan „nieustawione" jest odróżniony od „ustawione na pełne" — wykazane
   przypadkiem, w którym zmiana ustawienia projektu dochodzi do karty.
3. Tryby upoważnienia noszą nazwy wiążące: `manual`, `auto`, `plan`,
   `bypass permissions`.
4. Warstwa nie dotyka drzewa dokumentu.

## Zgłoszenia oczekujące na teren

Ustalenia z zamkniętych i biegnących terenów, które wykraczają poza ich zakres.
Każde zgłoszenie ma wskazany plik i wiersz. Zgłoszenie staje się terenem, gdy
Prowadzący je otworzy; do tego czasu jest wykazem, nie pracą.

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
| `proba-prototypow` | `teren/proba-prototypow` | `66e5ee0` przepływ wejścia · `61a5867` moduł Studio · `53bc3d5` odsyłacze | kontrola sesji nadzorującej wykonanie, weryfikacja Prowadzącego pomiarem |

Wynik: oba przedmioty wykonane, wszystkie kryteria spełnione. Zakresy trzech
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
