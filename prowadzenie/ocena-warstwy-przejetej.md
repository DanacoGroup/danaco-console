# Ocena warstwy przejętej

Rozstrzygnięcie, co z logiki przejętego klienta przechodzi do nowego produktu bez
zmian, co wymaga pracy i czego w przejętym kodzie nie ma. Ocena nie buduje
niczego — wyznacza granicę, poza którą kolejne tereny nie powstają na przedmiot
już istniejący.

Odniesieniem jest model produktu z [rejestru decyzji](decyzje.md): pozycja 5
(aplikacja → okno robocze → karta → widok), pozycja 6 (kaskada zakresu na
czterech kondygnacjach, dwie osie) i pozycja 9 (wstążka narzędziowa i przybornik
karty).

## 1. Jak mierzono

Każda warstwa oceniona uruchomieniem. Poniższe przyrządy dały wyniki
przytaczane w dalszej części; wynik pusty jest wszędzie odróżniony od braku
pomiaru.

| Przyrząd | Wywołanie | Wynik zbiorczy |
|---|---|---|
| sprawdzenie typów całości | `npx tsc --noEmit` w `budowa/klient-poprzedni` | kod wyjścia 0, czas 9,2 s |
| sprawdzenie typów warstwy osobno | `npx tsc` z jawnym wykazem plików warstwy, flagi z `tsconfig.json`, `--types vite/client` | 0 błędów w każdej z dziewięciu warstw |
| sprawdziany | `npx vitest run` | 40 plików, 508 sprawdzianów, wszystkie zdane, 4,40 s |
| żywy rdzeń | `danaco-console -port 17781 -dane /tmp/rdzen-ocena-dane -wymog-logowania false` | `rdzeń gotowy: komend=1077`, `transport: nasłuch 127.0.0.1:17781 gniazdo=/ws` |
| próba przeciw rdzeniowi | kod warstw zbundlowany bez przeróbki (`esbuild`, `--platform=node`) i uruchomiony pod Node 22.23.2 przeciw gniazdu `/ws` | wyniki przy każdej warstwie w rozdziale 4 |
| próba w przeglądarce | ten sam kod zbundlowany pod przeglądarkę, serwowany z gniazda lokalnego, prowadzony przez Playwright/Chromium (`PLAYWRIGHT_BROWSERS_PATH=/opt/ms-playwright`) | wyniki w rozdziale 4.1 |

Środowisko pomiaru: Node v22.23.2, TypeScript 5.8.3 (wersja przypięta
w `budowa/klient-poprzedni/package.json`, nie wersja maszyny), Go 1.26.5,
Chromium `chromium-1234`.

Próba w przeglądarce potwierdza najpierw, że mierzy: odczytuje znacznik strony
(`sonda wczytana`) i obecność funkcji sondy (`sonda_osadzona = true`), zanim
sięgnie po jakikolwiek wynik. Bez tych dwóch potwierdzeń zwraca „POMIAR
NIEWYKONANY" zamiast liczby.

**Miara dotknięcia drzewa dokumentu.** Plik uznaje się za budujący dokument,
gdy niesie odwołanie do `document.*`, `HTMLElement` i pokrewnych, `querySelector`,
`getElementById`, `innerHTML`, `classList`, `appendChild`, `insertBefore`,
`customElements` albo `ShadowRoot`. Miara szersza — obejmująca `window.*`,
`localStorage`, `addEventListener` i zdarzenia wskaźnika — daje wyniki
zafałszowane: `protokol/uzgodnienie.ts` trafia pod nią wyłącznie za sprawą
nazwy komendy `window.create` w komentarzu, a `polaczenie/gniazdo.ts` — za
sprawą `addEventListener` wołanego na gnieździe WebSocket, nie na dokumencie.
Wszystkie liczby poniżej pochodzą z miary węższej.

## 2. Wynik zbiorczy

| Warstwa | Plików | Wierszy | Buduje dokument | Wierszy w plikach budujących | Wierszy poza nimi | Sprawdzianów | Rozstrzygnięcie |
|---|---|---|---|---|---|---|---|
| `polaczenie` | 12 | 1016 | 0 z 12 | 0 | 1016 | 23 | **przechodzi bez zmian** |
| `protokol` | 22 | 1026 | 0 z 22 | 0 | 1026 | 0 | **wymaga pracy** — jeden plik z dwudziestu dwóch |
| `uwierzytelnienie` | 8 | 1853 | 2 z 8 | 893 | 960 | 0 | **wymaga pracy** |
| `modele` | 26 | 2994 | 19 z 26 | 2381 | 613 | 0 | **wymaga pracy** |
| `sterowanie` | 33 | 3484 | 19 z 33 | 1960 | 1524 | 0 | **wymaga pracy** |
| `aplikacja` | 32 | 3494 | 20 z 32 | 2460 | 1034 | 0 | **wymaga pracy** — przebudowa złożenia |
| `punkty-izolacji` | 18 | 3614 | 13 z 18 | 3304 | 310 | 4 | **wymaga pracy** |
| `konfiguracja` | 29 | 3732 | 16 z 29 | 2308 | 1424 | 0 | **wymaga pracy** |
| `dostepy` | 25 | 2842 | 15 z 25 | 1725 | 1117 | 0 | **wymaga pracy** |
| razem | 205 | 24 055 | 104 z 205 | 15 031 | 9 024 | 27 | |

Żadna z dziewięciu warstw nie jest nieobecna w całości. Rzeczy nieistniejące
wylicza rozdział 6 — są to pojęcia modelu, nie katalogi.

## 3. Kompletność wykazu

**Wykaz dziewięciu warstw nie jest podziałem stu trzech tysięcy wierszy — jest
podzbiorem obejmującym 24 055 wierszy, czyli dziesiątą część klienta.**

Pomiar całego drzewa `budowa/klient-poprzedni/src`:

| Miara | Wynik |
|---|---|
| plików `.ts` | 1311 |
| wierszy | 242 267 |
| wierszy w plikach niebudujących dokumentu (miara węższa) | 85 142 |
| wierszy w plikach bez jakiegokolwiek odwołania do przeglądarki (miara szersza) | 67 964 |
| wierszy miary szerszej leżących w dziewięciu warstwach wykazu | 6 597 |
| wierszy miary węższej leżących w dziewięciu warstwach wykazu | 9 024 |

Liczby stu trzech tysięcy nie odtwarza żadna z dwóch miar. Największe skupisko
logiki niebudującej dokumentu leży poza wykazem — w `moduly/` (47 159 wierszy
miarą szerszą, w tym `studio` 8371, `design` 7045, `apps` 5485, `browser` 3462,
`roundtable` 3321). Poza `moduly/` największe skupiska to `rozmowa` (1933),
`okno-komunikacji` (1763), `strona-glowna` (1550), `powloka` (1501)
i `okna-rownolegle` (1492).

Rozstrzygnięcia dla tej masy niniejsza ocena nie podejmuje: jej przedmiotem jest
wykaz dziewięciu warstw. Rozbieżność między liczbą stu trzech tysięcy a każdym
możliwym pomiarem jest usterką źródła i wraca do Prowadzącego (rozdział 7).

Kolumna „plików dotyka DOM" w wykazie wyjściowym różni się od pomiaru w pięciu
warstwach na dziewięć: `sterowanie` 12 wobec zmierzonych 19, `aplikacja` 14
wobec 20, `modele` 17 wobec 19, `konfiguracja` 15 wobec 16, `punkty-izolacji`
12 wobec 13. Zgadzają się `polaczenie` (0), `protokol` (0), `uwierzytelnienie`
(2) i `dostepy` (15). Wykaz wyjściowy zaniża liczbę plików związanych z widokiem
i przez to zawyża masę logiki nadającej się do przeniesienia.

## 4. Rozstrzygnięcia warstwa po warstwie

### 4.1 `polaczenie` — przechodzi bez zmian

Transport ramek WebSocket, kolejka wychodząca, magistrala zdarzeń, polityka
ponawiania, ustalenie adresu rdzenia, dziennik komunikatów nierozpoznanych,
obserwator ogniska.

**Pomiar.** Sprawdzenie typów warstwy osobno: 0 błędów. Sprawdziany:
`src/polaczenie/adres-rdzenia.test.ts` i `src/polaczenie/warstwa-lacznosci.test.ts`
— 23 zdane, w tym „odkłada ramkę wpisaną po zerwaniu i wydaje ją po wznowieniu"
oraz „przechodzi stany w kolejności łączenie → połączony → ponawianie →
połączony". Żaden z dwunastu plików nie dotyka drzewa dokumentu.

**Próba przeciw żywemu rdzeniowi w przeglądarce.** Chromium, gniazdo
`ws://127.0.0.1:17783/ws`, rdzeń zatrzymywany i podnoszony w trakcie:

```
etap_1_po_polaczeniu  stany=[rozlaczony +25ms, laczenie +25ms, polaczony +27ms]
                      etapy=[powitanie:true, sesja:true, okno:true, gotowe:true]
etap_2_po_zatrzymaniu stanTransportu="ponawianie" (+3032ms), oczekujace=0
etap_2_ramek_w_kolejce_po_wyslaniu = 1
etap_3_po_wznowieniu  stanTransportu="polaczony" (+8248ms), oczekujace=0
                      odpowiedzZKolejki={"udany":true,"liczba":1}
                      bledyOkna=[]
```

Ramka wpisana przy rozłączeniu trafiła do kolejki, przeżyła zerwanie, została
wydana po wznowieniu i wróciła z odpowiedzią rdzenia. Okno nie zgłosiło ani
jednego błędu.

**Zależność od modelu okna.** Jedna: `obserwator-ogniska.ts` obsługuje zdarzenie
`session.focus.changed` i nazywa sesję „kartą sesji" (wiersze 12 i 25). Jest to
nazewnictwo, nie rozstrzygnięcie — obserwator przenosi identyfikatory
z kontraktu i nie zakłada, ile okien ani kart istnieje. Warstwa przechodzi.

**Do naprawy przy przeniesieniu, bez zmiany rozstrzygnięcia.**
`polaczenie/gniazdo.ts` wiersz 46 niesie
`gniazdo.addEventListener('error', () => this.gniazdo?.close())`. Pod
implementacją WebSocket Node 22 (undici) wywołanie `close()` na gnieździe
w stanie nawiązywania wyzwala kolejne zdarzenie `error`, a to kolejne
`close()` — pomiar kończy się `RangeError: Maximum call stack size exceeded`.
W Chromium ta sama ścieżka przebiega poprawnie, co wykazuje próba wyżej.
Sprawdziany warstwy tego nie łapią, bo stoją na atrapie gniazda. Skutek dotyczy
wyłącznie uruchomień poza przeglądarką; produkt stoi na WebView, więc nie jest
to przeszkoda w przejęciu, lecz zapora do postawienia, zanim warstwa zostanie
użyta w narzędziach.

### 4.2 `protokol` — wymaga pracy

Koperta kontraktu, korelacja żądanie–odpowiedź, ramka tekstowa, kanał
komunikatów, sprawdzian kształtu odpowiedzi, wywołanie jako obietnica, sesja,
uzgodnienie, opakowania osiemnastu komend `connection.*`, `session.*`,
`home.*`, `environment.*`, `workspace.*`, `window.*` i `module.*`.

**Pomiar.** Sprawdzenie typów warstwy osobno: 0 błędów. Żaden z dwudziestu
dwóch plików nie dotyka drzewa dokumentu. Sprawdzianów własnych warstwa nie ma;
`src/kontrakt.test.ts` (8 zdanych) pilnuje zgodności generatu ze źródłem prawdy
i równości stałej portu klienta z `PortDomyslny` rdzenia.

**Próba przeciw żywemu rdzeniowi.** Uzgodnienie domknięte w całości:

```
etapy_uzgodnienia = [powitanie:true, sesja:true, okno:true, gotowe:true]
identyfikator_sesji_od_rdzenia = "ses_88c93dc51c701359"
okno_otwarte_przez_rdzen = {"id":"okn_4e94124dc004608b","sessionId":"ses_88c93dc51c701359"}
session_list_udany = true, session_list_liczba_sesji = 1
komenda_nieistniejaca_udany = false, kod_bledu = "not_found"
dziennik_nieznanych_liczba = 0
```

Pozostałe drogi warstwy odpowiadają: `home.enter` (4 środowiska, 4 sesje,
4 obecności), `environment.enter` (środowisko, 9 modułów), `workspace.enter`
(sesja, moduł, okno), `module.list` (9 z 15), `session.rename`, `session.focus`,
`session.copy`, `session.archive`, `session.archive.list`. Odmowa wraca jako
`Wynik` z polem `blad`, nigdy jako wyjątek — sprawdzone komendą spoza kontraktu.

**Zależność od modelu okna — nazwana wprost.** `protokol/uzgodnienie.ts` zakłada
sesję i otwiera okno komunikacji **w chwili nawiązania połączenia**, jako drugi
i trzeci krok powitania. Decyzja 5 stanowi: „Moment powstania sesji jest jeden:
pierwsza wiadomość wysłana do modelu. Nie otwarcie okna, nie wejście do modułu,
nie otwarcie karty czatu". Uzgodnienie w obecnej postaci zakłada sesję rdzenia
każdemu, kto uruchomi aplikację, i to zanim Operator cokolwiek napisze. Jest to
rozstrzygnięcie warstwy, nie jej nazewnictwo, więc plik nie przechodzi.

**Druga usterka tego samego pliku.** `zalozSesje` ogłasza powodzenie etapu
`sesja` (`zglos('sesja', wynik)`) **przed** zapisaniem identyfikatora
(`kanal.sesja().ustaw(sesja.id)`). Słuchacz postępu czytający `sesja.id()`
w tym etapie dostaje napis pusty — pomiar w Chromium zwrócił `idSesji: ""` przy
`etapy=[…,"sesja:true",…]`.

**Zakres pracy:** `uzgodnienie.ts` (155 wierszy) — przeniesienie założenia sesji
z powitania na pierwszą wiadomość i odwrócenie kolejności ogłoszenia względem
zapisu identyfikatora. Pozostałe dwadzieścia jeden plików (871 wierszy)
przechodzi bez zmian.

### 4.3 `uwierzytelnienie` — wymaga pracy

Rozpoznanie stanu bramki, wejście hasłem albo kodem, rejestracja, odzyskanie,
zmiana hasła, przedłużenie tokenu, magazyn sesji bramki, tożsamość urządzenia,
ekran logowania.

**Pomiar.** Sprawdzenie typów warstwy osobno: 0 błędów. Sprawdzianów brak.
Dwa pliki budują dokument: `ekran-logowania.ts` (569 wierszy) i
`postac-bramki.ts` (324) — razem 893 wiersze widoku, który nie jest przejmowany.

**Próba przeciw żywemu rdzeniowi.** Wszystkie siedem komend warstwy istnieje
i odpowiada odmową merytoryczną, nie brakiem komendy:

```
auth.verify         conflict         "bramka: konta właściciela jeszcze nie ma — najpierw rejestracja"
auth.login          validation_failed "bramka: metoda wejścia  nie należy do kontraktu"
auth.register       validation_failed "bramka: rejestracja bez loginu"
auth.recover        validation_failed "bramka: odzyskanie konta bez podanego adresu"
auth.reset          validation_failed "bramka: ustawienie nowego hasła bez hasła"
auth.password.reset validation_failed "bramka: zmiana hasła wymaga hasła bieżącego i nowego"
auth.token.refresh  validation_failed "bramka: przedłużenie bez tokenu z połączenia"
```

Rdzeń uruchomiony z `-wymog-logowania false` nadal prowadzi bramkę: `auth.verify`
odpowiada, że konta właściciela nie ma. Wynik pusty jest tu odróżniony od braku
pomiaru — każda komenda zwróciła treść odmowy.

**Zależność od modelu okna.** Brak. Warstwa nie zna ani sesji pracy, ani okna
komunikacji; sesja bramki jest bytem osobnym.

**Zakres pracy.** Trzy rzeczy:

1. `sesja-bramki.ts` wiersze 47 i 56 oraz `tozsamosc-urzadzenia.ts` wiersz 34
   sięgają po `window.localStorage` i `window.sessionStorage`. Miejsce
   przechowania sesji bramki w powłoce natywnej jest rozstrzygnięciem do
   podjęcia, nie przeniesieniem.
2. `zrodlo-auth.ts` wiersz 301 woła `window.setTimeout` — wiąże warstwę
   z globalnym obiektem okna bez potrzeby.
3. Widok bramki (893 wiersze) powstaje na nowo.

Sześć plików nie buduje dokumentu (960 wierszy). Przechodzą bez zmian trzy —
374 wiersze: `odmowy-auth.ts` (198, zdania odmów dla Operatora),
`rozpoznanie-bramki.ts` (155, rozstrzygnięcie drogi wejścia) oraz `indeks.ts`
(21). Pozostałe trzy — `zrodlo-auth.ts` (397), `sesja-bramki.ts` (151),
`tozsamosc-urzadzenia.ts` (38) — czekają na rozstrzygnięcie miejsca
przechowania sesji bramki.

### 4.4 `modele` — wymaga pracy

Rejestr kont modeli, wybór konta domyślnego, warstwy tożsamości i ich dokumenty,
podgląd promptu złożonego, wykaz kategorii tożsamości.

**Pomiar.** Sprawdzenie typów warstwy osobno: 0 błędów. Sprawdzianów brak.
Dziewiętnaście z dwudziestu sześciu plików buduje dokument — 2381 wierszy
widoku wobec 613 wierszy poza nim.

**Próba przeciw żywemu rdzeniowi**, prowadzona przez własne źródła warstwy
(`utworzZrodloKont`, `utworzZrodloModeli`) bez przeróbki:

```
modele_lista_kont              = {"liczba":0,"niepowodzenia":[]}
modele_wykaz_modeli            = {"liczba":1,"pierwsze":["bielik"]}
identity.category.list         = {"udany":true,"categories":"tablica[13]"}
identity.effective.get         = {"udany":true,"mode":"ZASTAP","layers":"tablica[0]",
                                  "missingRequiredCategoryIds":"tablica[1]"}
identity.document.get          = {"udany":true,"documents":"tablica[0]"}
```

Wykaz kont pusty przy pustym wykazie niepowodzeń — to wynik pusty, nie brak
pomiaru; źródło warstwy rozróżnia oba przypadki wywołaniem `NaNiepowodzenie`.

**Zależność od modelu okna.** Żadna. Wyszukanie `sessionId` i `windowId`
w całym katalogu `modele/` nie daje ani jednego trafienia. Warstwa jest
niezależna od tego, jak produkt organizuje okna i karty.

**Zakres pracy:** wyłącznie widok — 2381 wierszy formularzy, paneli, wykazów,
zakładek i kontrolek. Do widoku należą przy tym dwa pliki, których nazwa tego nie
zapowiada: `zrodlo-tozsamosci.ts` (122) i `stan-tozsamosci.ts` (174) budują
dokument i nie przechodzą razem z resztą źródeł.

Przechodzi siedem plików niebudujących dokumentu — 613 wierszy:
`stan-kont.ts` (154), `zrodlo-kont.ts` (125), `warstwy-tozsamosci.ts` (80),
`zapisy-kont.ts` (82), `rodzaje-kont.ts` (67), `indeks.ts` (53),
`zrodlo-modeli.ts` (52). Brak sprawdzianów jest jedynym powodem, dla którego
„przechodzi bez zmian" nie zostaje tu wypowiedziane: przeniesienie odbywa się
bez sieci zabezpieczającej.

### 4.5 `sterowanie` — wymaga pracy

Komplet sterowania okna: tryb uprawnień, model główny i zapasowy, nakład
rozumowania, moduł, rola, środowisko wykonania, host, katalogi robocze, rejestr
kanałów modeli, rejestr agentów, wartość obowiązująca.

**Pomiar.** Sprawdzenie typów warstwy osobno: 0 błędów. Sprawdzianów brak.
Dziewiętnaście z trzydziestu trzech plików buduje dokument — 1960 wierszy widoku
wobec 1524 poza nim.

**Próba przeciw żywemu rdzeniowi**, prowadzona przez rejestry warstwy:

```
sterowanie_rejestr_kanalow = {"liczba":2}
sterowanie_rejestr_modulow = {"liczba":15,"pierwsze":["agents","apps","assistant",
                                                     "automations","browser","design"]}
sterowanie_rejestr_agentow = {"liczba":0}
sterowanie_window_update   = {"udany":true,"window":"obiekt"}
```

Odczyt wartości obowiązującej (`utworzOdczytObowiazujacej` → `config.get`
bez poziomu) domknięty na żywym rdzeniu:

```
ster_znane   = true
ster_naklad  = {"wartosc":"high","zasieg":"window","bytZasiegu":"okn_42e60925b211e6bd"}
ster_zdanie  = "high"
polityka_efektywna_liczba_kluczy = 64
```

**Zależność od modelu okna — nazwana wprost.** Cała warstwa traktuje `Window`
kontraktu jako miejsce nastawy: „Tryb jest parametrem okna", „Moduł jest
parametrem okna, nie sesji", „Zasięg wykonania jest parametrem okna". W modelu
decyzji 5 `Window` jest **kartą**, więc te nastawy siedzą na kondygnacji karty —
dokładnie tam, gdzie decyzja 6 je stawia („karta z interakcją modelu:
komponenty w oknie czatu"). Założenie warstwy nie sięga jej rozstrzygnięć;
rozjazd jest nazewniczy i mieści się w rozdziale 5.

**Zakres pracy.**

1. `tryb-uprawnien.ts` buduje listę wprost z `Object.values(PermissionMode)`.
   Kontrakt niesie **sześć** trybów, decyzja 6 nazywa **cztery**. Pomiar:
   ```
   tryby_w_kontrakcie              = ["manual","acceptEdits","plan","auto","dontAsk","bypassPermissions"]
   tryby_nadmiarowe_wobec_decyzji_6 = ["acceptEdits","dontAsk"]
   rdzen_przyjmuje_tryby            = wszystkie sześć zapisane i odczytane zgodnie
   ```
   Operator zobaczy sześć pozycji tam, gdzie decyzja przewiduje cztery. Rozjazd
   jest w kontrakcie, nie w kliencie (rozdział 7).
2. `atrapa-rdzenia.ts` (211 wierszy) nie ma ani jednego wołacza w całym drzewie,
   a jego nagłówek stanowi, że „nie ma ani jednego wołacza poza plikami
   `*.test.ts`" — plików sprawdzianów w tej warstwie nie ma wcale. Kod martwy
   wraz z nagłówkiem mówiącym nieprawdę.
3. Widok (1960 wierszy) powstaje na nowo.

Czternaście plików warstwy nie buduje dokumentu (1524 wiersze). Po odjęciu
kodu martwego przechodzi trzynaście — 1313 wierszy: `wartosc-obowiazujaca.ts`
(176), `etykiety-sterowania.ts` (164), `klucze-ustawien.ts` (147),
`zrodlo-kanalow.ts` (142), `zmiana-ustawienia.ts` (131), `adnotacje-wykonania.ts`
(116), `stan-sterowania.ts` (102), `rejestr-agentow.ts` (84),
`zmiana-okna.ts` (66), `rejestr-kanalow.ts` (66), `rejestr-modulow.ts` (55),
`komunikat-zmiany.ts` (37), `indeks.ts` (27).

`wartosc-obowiazujaca.ts` jest w tej warstwie rzeczą najcenniejszą: nazywa
wprost, czego komenda `config.get` nie obejmuje, i nie udaje, że obejmuje.

### 4.6 `aplikacja` — wymaga pracy, w części zasadniczej przebudowy

Złożenie aplikacji, router tras, katalog tras, przełącznik motywu, pasek
aplikacji, menu, scena sesji, przestrzeń modułu, wiązanie gniazda z oknem,
wskaźnik łączności, widoki trzech tras.

**Pomiar.** Sprawdzenie typów warstwy osobno: 0 błędów. Sprawdzianów brak.
Dwadzieścia z trzydziestu dwóch plików buduje dokument — 2460 wierszy widoku
wobec 1034 poza nim.

**Próba przeciw żywemu rdzeniowi.** Drogi warstwy odpowiadają:
`window.create` dla drugiego okna — powodzenie; `session.open` — powodzenie,
sesja wraz z dwoma oknami. `queue.action` i `context.transfer` odpowiadają
`not_found` na żądanie bez wskazania bytu, czyli istnieją i sprawdzają wejście.

**Zależność od modelu okna — nazwana wprost i sięgająca rozstrzygnięć.**

`aplikacja/trasy.ts` ustala **trzy trasy najwyższego rzędu**: `strona-glowna`
(Centrum dowodzenia), `srodowisko`, `pulpit` (Mission Control), z komentarzem
„Trasa otwierana przy uruchomieniu — Centrum dowodzenia, nie okno pracy".
Decyzja 5 stanowi wprost: „Nowe okno robocze otwiera się na centrum dowodzenia.
Centrum nie jest więc osobnym miejscem, do którego się wchodzi, lecz stanem
początkowym każdego okna roboczego". Router pokazuje w danej chwili jedną trasę
i trzyma po jednym egzemplarzu widoku na trasę; model decyzji 5 przewiduje wiele
okien roboczych otwartych naraz i przełącznik między nimi.

`aplikacja/scena-sesji.ts` stanowi w nagłówku: „Okno komunikacji to nie karta
sesji. **Karta w pasie powłoki jest sesją rdzenia** i rządzi się komendami
`session.*`; okno sceny jest bytem podrzędnym wobec sesji i rządzi się komendami
`window.*`. Liczbę okien ustawia wyłącznie przełącznik »Okna komunikacji: 1 2 3«".
Scena stawia od jednego do czterech okien obok siebie
(`okna-rownolegle/identyfikatory.ts`, `LICZBA_MAX = 4`), przy czym sufit czterech
jest umotywowany obsadą pętli wielomodelowej, nie widokiem. Decyzja 5 przewiduje
karty bez ograniczenia liczby, widok dzielony zawsze dwuczłonowy i okno boczne
jako trzecie miejsce wyświetlenia.

`aplikacja/polaczenie-z-rdzeniem.ts` składa **jedno** połączenie, **jedną**
sesję i **jedno** uzgodnienie na całą aplikację. Model decyzji 5 przewiduje
wiele okien roboczych naraz, każde z co najwyżej jedną sesją.

**Zakres pracy.** Dwanaście plików warstwy nie buduje dokumentu (1034 wiersze).
Z nich przechodzi dziewięć — 831 wierszy:

| Plik | Wierszy | Co niesie |
|---|---|---|
| `zrodlo-posuniec.ts` | 279 | rozstrzyganie, czy zdarzenie wyszło z tego połączenia |
| `rozstrzyganie-sprawcy.ts` | 145 | ta sama rzecz na poziomie pojedynczego zdarzenia |
| `zamiary-pulpitu.ts` | 99 | zamiary wnoszone z pulpitu |
| `ustawienia-okna-sledzone.ts` | 77 | odtworzenie „co zmieniono" z `window.changed` |
| `polaczenie-z-rdzeniem.ts` | 71 | złożenie transportu, kanału i uzgodnienia |
| `akcje-ustawien.ts` | 66 | czynności okna ustawień |
| `otwarcie-okna.ts` | 50 | `window.create` dla okna spoza uzgodnienia |
| `ulotnosc-okna.ts` | 25 | polityka pamięci rozmowy modułu |
| `arkusze-stylow.ts` | 19 | kolejność wczytania warstw wizualnych |

Nie przechodzą — mimo że dokumentu nie budują — `aplikacja.ts` (103),
`trasy.ts` (49) i `stan-pary-gniazda.ts` (51): pierwsze dwa niosą model trzech
tras, trzeci wiąże stan pętli z gniazdem sceny okien równoległych. Wraz
z widokiem (2460 wierszy) daje to 2663 wiersze powstające na nowo: złożenie,
router, katalog tras, przełącznik tras, scena sesji, przestrzeń modułu,
wiązanie gniazda i trzy widoki tras.

### 4.7 `punkty-izolacji` — wymaga pracy

Okno Punktów Izolacji: trzy przełączniki kontekstu, osiem zakresów technicznych,
profile, podgląd polityki obowiązującej, selektor poziomu zasięgu i warstwy.

**Pomiar.** Sprawdzenie typów warstwy osobno: 0 błędów. Sprawdziany:
`src/punkty-izolacji/okno-punktow-izolacji.test.ts` — 4 zdane, w tym „zasięg
i warstwa jadą w żądaniach, nie są zaszyte". Trzynaście z osiemnastu plików
buduje dokument — 3304 wiersze widoku wobec 310 poza nim; jest to warstwa
najbardziej związana z widokiem z całego wykazu.

**Próba przeciw żywemu rdzeniowi.** Przy poprawnym zaadresowaniu wszystkie
dwanaście komend `isolation.*` domyka pełny obieg:

```
izolacja_scope_list             = 9 poziomów: application global environment module
                                  modulePair project session role window
izolacja_context_get (okno / karta sesji / globalny / projekt) = switches: tablica[3]
izolacja_context_set  (poziom session, warstwa session)        = switches: tablica[3]
izolacja_technical_set(poziom session, warstwa session)        = switches: tablica[8]
izolacja_context_set  (poziom window,  warstwa default)        = switches: tablica[3]
izolacja_layer_set                                             = {"layer":"session"}
izolacja_policy_preview   = polityka z objaśnieniem każdego przełącznika
izolacja_profile_save / assign / load / list / delete          = wszystkie powodzenie
```

Warstwa `session` obowiązuje wyłącznie na poziomie `session`: żądanie
`{scope: window, layer: session}` rdzeń odrzuca zdaniem „warstwa karty sesji
obowiązuje na poziomie session, nie na poziomie window".

**Zależność od modelu okna.** Poziom `window` jest wśród poziomów, na których
izolację wolno zapisać, i zapis na nim przechodzi (pomiar wyżej). Decyzja 5
stanowi: „Izolacja ma dwa poziomy: okno robocze niesie izolację swojej sesji,
karta może nieść własną". Rdzeń to obsługuje. Zależność warstwy jest
nazewnicza — `katalog-izolacji.ts` nazywa `ConfigScope.Session` „Kartą sesji",
a `ConfigScope.Window` „Oknem komunikacji", czyli odwrotnie niż decyzja 5.

**Zakres pracy.**

1. `katalog-izolacji.ts`, `ZASIEGI_OD_NAJSZERSZEGO` — osiem pozycji z dziewięciu.
   Brakuje `ConfigScope.Application`, choć `isolation.scope.list` żywego rdzenia
   oddaje go w komplecie. Poziomu aplikacji nie da się wskazać w selektorze:
   ```
   katalog_etykiety_zasiegu                    = 9
   katalog_zasiegi_od_najszerszego             = 8
   katalog_poziomy_pominiete_w_liscie_wyboru   = ["application"]
   ```
2. `posijKomende` jest zdefiniowana **czterokrotnie**: raz jako funkcja
   udostępniona w `komenda.ts` (wiersz 22) i trzy razy prywatnie —
   `obszar-kontekst.ts` wiersz 109, `obszar-poziomy.ts` wiersz 81,
   `obszar-techniczny.ts` wiersz 234. Kopia w `obszar-kontekst.ts` jest przy tym
   nietypowana (`<T>` wraz z rzutami `as never`), więc znosi wiązanie żądania
   z odpowiedzią, które daje wersja udostępniona. Jedno pojęcie pod czterema
   osobnymi definicjami tej samej nazwy.
3. Widok (3304 wiersze) powstaje na nowo.

### 4.8 `konfiguracja` — wymaga pracy

Okno konfiguracji: katalog kategorii i definicji, pasek punktu widzenia, adres
ustawienia, rozstrzyganie wartości obowiązującej, łańcuch zapisów, obszary
konfiguracji sesji, kontrolki pól.

**Pomiar.** Sprawdzenie typów warstwy osobno: 0 błędów. Sprawdzianów brak.
Szesnaście z dwudziestu dziewięciu plików buduje dokument — 2308 wierszy widoku
wobec 1424 poza nim.

**Próba przeciw żywemu rdzeniowi**, prowadzona przez własne źródła warstwy:

```
konfiguracja_kategorie   = {"liczba":14}
konfiguracja_definicje   = {"liczba":64,"pierwsze":["kanal_modelu","harness.program_claude",
                            "tozsamosc.tryb_domyslny","tryb_uprawnien","katalogi_robocze"]}
konfiguracja_zapis_pod_adresem       = {"udany":true,"entry":"obiekt"}
konfiguracja_przywrocenie            = {"udany":true,"entries":"tablica[1]"}
konfiguracja_obszary_obowiazujace    = {"udany":true,"effective":"obiekt"}
config.session.set                   = {"udany":true,"config":"obiekt","storedAreas":"tablica[0]"}
```

**Kaskada zakresu z decyzji 6 nie jest zrealizowana — ani po stronie klienta,
ani po stronie rdzenia.** Jest to najważniejsze ustalenie całej oceny i jedyne,
które podważa rozstrzygnięcie warstwy, a nie jej widok.

Rachunek klienta, uruchomiony na wpisach zapisanych na trzech kondygnacjach
(`application`, `project`, `session`) przy punkcie widzenia ustawionym na kartę
(`window`):

```
kaskada_klienta_punkt_karta_wartosc   = "DOMYSLNA"
kaskada_klienta_punkt_karta_domyslna  = true
kaskada_klienta_dlugosc_lancucha      = 3
kaskada_klienta_punkt_sesja_wartosc   = "Z_SESJI"
kaskada_klienta_punkt_projekt_wartosc = "Z_PROJEKTU"
```

Wpisy istnieją i widnieją w łańcuchu, lecz w punkcie widzenia karty nie
obowiązuje żaden — obowiązuje wartość domyślna katalogu. `rozstrzygniecie.ts`,
funkcja `zasiegObowiazuje`, przepuszcza wyłącznie wpis globalny albo wpis
zapisany dokładnie na poziomie punktu widzenia i dokładnie dla jego bytu.
Dziedziczenia z kondygnacji szerszej nie ma. Plik mówi o tym wprost
w nagłówku: „Klient nie zna pełnej ścieżki bytów […], więc jej nie zgaduje".

Rachunek rdzenia, na żywym połączeniu, z zapisami dokładanymi po kolei:

```
zapis na poziomie global   → polityka dla okna: {"value":"low","scope":"global"}
zapis na poziomie session  → polityka dla okna: {"value":"low","scope":"global"}   ← karta sesji pominięta
zapis na poziomie window   → polityka dla okna: {"value":"high","scope":"window"}
```

Wartość zapisana na kondygnacji sesji **nie obowiązuje** karty tej sesji.
Powodem jest `server/internal/core/adapter_ustawienia.go`, funkcja
`kontekstZasiegu` (wiersz 75), która buduje kontekst rozstrzygania wypełniając
wyłącznie pole `Okno`, podczas gdy `server/internal/konfig/kontekst_zasiegu.go`
przewiduje siedem pól kondygnacji. Żądanie `ConfigGetRequest` nie ma pól, którymi
klient mógłby brakujące kondygnacje podać.

**Zależność od modelu okna.** `zasiegi.ts` nazywa `ConfigScope.Session` „kartą
sesji", a `ConfigScope.Window` „oknem komunikacji" — odwrotnie niż decyzja 5.

**Zakres pracy.**

1. Kaskada czterech kondygnacji z dziedziczeniem — do zbudowania. Rozstrzygnięcie,
   czy liczy ją rdzeń, czy klient, jest zgłoszeniem (rozdział 7).
2. `zasiegi.ts`, `ZASIEGI_OD_NAJWEZSZEGO` — osiem pozycji z dziewięciu; brakuje
   `ConfigScope.Application`. `pierwszenstwoZasiegu('application')` zwraca 8,
   czyli wartość poza wykazem, przypadkiem zbieżną z właściwym miejscem.
   ```
   liczba_poziomow_w_kontrakcie             = 9
   liczba_poziomow_w_zasiegi_ts             = 8
   poziomy_kontraktu_pominiete_w_zasiegi_ts = ["application"]
   ```
3. Widok (2308 wierszy) powstaje na nowo.

Trzynaście plików warstwy nie buduje dokumentu (1424 wiersze). Pracy wymagają
dwa z nich — `rozstrzygniecie.ts` (115) i `zasiegi.ts` (120); przechodzi
jedenaście, 1189 wierszy: `stan-obszarow-sesji.ts` (213),
`stan-konfiguracji.ts` (201), `zrodlo-wartosci.ts` (167), `obszary-sesji.ts`
(156), `zrodlo-katalogu.ts` (101), `zrodlo-obszarow-sesji.ts` (100),
`wybor-kontrolki.ts` (73), `adres-ustawienia.ts` (68), `indeks.ts` (54),
`widocznosc-pol.ts` (33), `dymek-objasnienia.ts` (23).

Przechodzi bez zmian rozróżnienie „nieustawione" od „ustawione", którego wymaga
decyzja 6: `Rozstrzygniecie` niesie pole `domyslna` i pole `zrodlo`, więc brak
zapisu jest odróżniony od zapisu o wartości pełnej.

### 4.9 `dostepy` — wymaga pracy

Sekcja dostępów: punkty dostępu (most MCP albo katalog lokalny), nadania okna,
tryb nadania, wybór korzeni, obszar katalogu roboczego.

**Pomiar.** Sprawdzenie typów warstwy osobno: 0 błędów. Sprawdzianów brak.
Piętnaście z dwudziestu pięciu plików buduje dokument — 1725 wierszy widoku
wobec 1117 poza nim.

**Próba przeciw żywemu rdzeniowi.** Pełny obieg poprowadzony przez własne źródła
warstwy (`utworzZrodloPunktow`, `utworzZrodloNadan`), bez przeróbki:

```
dostepy_lista_poczatkowa  = {"liczba":3,"niepowodzenia":0}
dostepy_dodaj_punkt       = {"udany":true,"point":"obiekt"}
dostepy_lista_po_dodaniu  = {"liczba":4,"rodzaje":["localDirectory","mcpBridge","mcpBridge","mcpBridge"],
                             "nazwy":["sonda-katalog","danaco-data","danaco-system","danaco-web"]}
dostepy_sprawdz_punkt     = {"udany":true,"status":"reachable","checkedAt":1787523483604,"roots":"tablica[1]"}
dostepy_nadaj             = {"udany":true,"grant":"obiekt","grants":"tablica[1]"}
dostepy_nadania_po        = {"liczba":1,"tryby":["write"],"korzenie":[["/tmp/rdzen-ocena-dane"]]}
dostepy_zmien_nadanie     = {"udany":true,"grant":"obiekt","grants":"tablica[1]"}
dostepy_odbierz_nadanie   = {"udany":true,"removed":true,"grants":"tablica[0]"}
dostepy_usun_punkt        = {"udany":true,"removed":true}
```

Warstwa rozróżnia brak wyniku od wyniku pustego wprost: `stany-odczytu.ts` niesie
cztery fazy odczytu (spoczynek, odczyt, błąd, gotowe) i osobne zdanie powodu, bo
— jak stanowi jej nagłówek — „pusty wykaz punktów znaczy co innego, gdy rdzeń
jeszcze nie odpowiedział, co innego, gdy odmówił, i co innego, gdy rejestr jest
pusty".

Sposób nadawania zakresu odpowiada decyzji 6: `AccessPoint` niesie `roots`,
a `AccessGrant` — podzbiór tych korzeni, przy czym pusty podzbiór znaczy komplet.
Zakres nadaje się wskazaniem miejsca, nie wyliczeniem uprawnień.

**Zależność od modelu okna — nazwana wprost i sięgająca rozstrzygnięć.**
Nadanie dostępu istnieje **wyłącznie na kondygnacji karty**.
`AccessGrantAddRequest` przyjmuje jedno pole bytu — `windowId`; kontrakt stanowi
przy polu `AccessGrant.windowId`: „nadanie zyje per okno, nie per sesja".
Wszystkie dziewięć komend `access.*` — `point.add`, `point.list`,
`point.update`, `point.remove`, `point.check`, `grant.add`, `grant.list`,
`grant.update`, `grant.remove` — zna wyłącznie punkt i okno.

Decyzja 6 wymaga osi dostępu na **czterech** kondygnacjach — aplikacji,
projektu, sesji i karty — z dziedziczeniem i z nadpisaniem w obie strony.
Trzech górnych kondygnacji tej osi nie ma ani w kontrakcie, ani w rdzeniu, ani
w kliencie.

**Zakres pracy.**

1. Oś dostępu na kondygnacjach aplikacji, projektu i sesji — do zbudowania,
   poczynając od kontraktu.
2. Widok (1725 wierszy) powstaje na nowo.

Przechodzi bez zmian dziesięć plików niebudujących dokumentu — 1117 wierszy:
`stan-dostepow.ts` (243), `zrodlo-punktow.ts` (150), `nazwy-dostepow.ts` (142),
`stan-katalogu-roboczego.ts` (125), `zrodlo-nadan.ts` (120),
`klucze-katalogu.ts` (100), `zapisy-nadan.ts` (88), `ostrzezenie-zapisu.ts` (72),
`indeks.ts` (57), `dialog-katalogu.ts` (20).

## 5. Zależność od modelu okna: nazwy są odwrócone

Przejęty klient i decyzja 5 nazywają te same dwa byty kontraktu odwrotnie.
Nie jest to różnica etykiety w jednym pliku — jest to konsekwentna wykładnia
przenikająca cały klient.

| Byt kontraktu | Jak nazywa przejęty klient | Jak nazywa decyzja 5 |
|---|---|---|
| `Session` | karta sesji, karta w pasie powłoki | **okno robocze** |
| `Window` | okno komunikacji, okno sceny | **karta** (czatu) |

Miejsca, w których wykładnia przejętego klienta stoi wprost w kodzie:

- `aplikacja/scena-sesji.ts`, nagłówek: „Karta w pasie powłoki jest sesją rdzenia".
- `sterowanie/model-karty-sesji.ts`, nagłówek: „Karta z czterema oknami to jedno
  żądanie zamiast czterech".
- `konfiguracja/zasiegi.ts`, `NAZWY_ZASIEGOW`: `Session` → „karta sesji",
  `Window` → „okno komunikacji".
- `punkty-izolacji/katalog-izolacji.ts`, `ETYKIETY_ZASIEGU`: tak samo.
- `polaczenie/obserwator-ogniska.ts`: „słuchacz jednej karty sesji".

Wykładnia ta ma pokrycie w samym kontrakcie, więc nie jest wymysłem klienta:
`shared/contract.ts` opisuje `ConfigScope.Session` jako „Karta sesji", a schemat
bazy rdzenia trzyma dla tego poziomu kod `karta_sesji`
(`WARTOSCI_BAZY_CONFIG_SCOPE`). Nazewnictwa kontraktu decyzja nie zmienia, więc
rozjazd jest trwały i wymaga jednego rozstrzygnięcia dla całej budowy: w kodzie
i w kontrakcie obowiązują nazwy `Session` i `Window`, w interfejsie
i dokumentacji — „okno robocze" i „karta". Zapisanie tego wprost jest
zgłoszeniem (rozdział 7); dopóki nie zapadnie, każdy teren czytający te pliki
będzie odczytywał je odwrotnie.

Cztery kondygnacje decyzji 6 odpowiadają czterem z dziewięciu poziomów
kontraktu:

| Kondygnacja decyzji 6 | Poziom kontraktu |
|---|---|
| aplikacja | `ConfigScope.Application` |
| projekt | `ConfigScope.Project` |
| sesja (okno robocze) | `ConfigScope.Session` |
| karta | `ConfigScope.Window` |

Pięć poziomów kontraktu decyzja nie nazywa: `global`, `environment`, `module`,
`modulePair`, `role`. Są one obecne w rdzeniu i rozstrzygają wartości — poziom
`global` jest wręcz jedynym, który w obecnym rachunku obowiązuje ponad punktem
widzenia (rozdział 4.8). Ich stosunek do czterech kondygnacji jest do
rozstrzygnięcia (rozdział 7).

## 6. Czego w przejętym kodzie nie ma

Poniższego nie ma ani w kliencie, ani w kontrakcie `budowa/shared/contract.json`.
Wyszukania przeprowadzone na całym drzewie klienta i na kontrakcie.

| Rzecz z modelu produktu | Stan | Gdzie powstanie |
|---|---|---|
| **układ paneli i widok dzielony** | brak. `splitView`, `paneLayout`, `sideWindow`, `tabStrip` — zero trafień w `contract.json`. W kliencie zero plików o nazwie widoku dzielonego. `StudioSplitOrientation` dotyczy podziału powierzchni dokumentu w module Studio i dopuszcza podział w pionie, którego decyzja 5 zakazuje | kontrakt (nowy typ układu w stanie sesji) oraz nowa warstwa widoku; decyzja 5 stanowi, że układ „wchodzi do modelu danych, nie tylko do warstwy widoku" |
| **okno boczne** | brak. Prawa strona ramy z własnym pasmem kart nie istnieje ani w kontrakcie, ani w kliencie | jw. |
| **karty inne niż czat** | brak. `Window` kontraktu nie ma pola rodzaju; opis typu brzmi „Okno komunikacji — byt posredni miedzy sesja a wiadomoscia". Edytor, przeglądarka, wykaz plików, podgląd zmian, terminal i tłumacz istnieją wyłącznie jako **moduły** (`moduly/`), nie jako karty osadzane w oknie roboczym | kontrakt (rodzaj karty w `Window`) oraz warstwa kart |
| **przełącznik okien roboczych** | brak. Wykaz okien roboczych jako byt odrębny od wykazu sesji nie istnieje; klient prowadzi jeden pas kart, którymi są sesje rdzenia | kontrakt oraz rama aplikacji |
| **projekt jako byt** | **brak w kontrakcie**. Nie ma typu `Project` ani żadnej komendy `project.*`. Istnieją wyłącznie `session.project.set` i `session.project.clear`, czyli przypisanie sesji do projektu, którego nie da się założyć. Tabela `projekt` istnieje w schemacie rdzenia — sprawdziany rdzenia zakładają wiersze wprost zapytaniem SQL. Pomiar: `session.project.set {projectId:"brak"}` → `not_found` „projekt brak nie istnieje" | kontrakt — pełny cykl życia projektu |
| **kaskada zakresu na czterech kondygnacjach** | brak. Rachunek klienta nie dziedziczy w ogóle, rachunek rdzenia obejmuje wyłącznie poziom okna i poziom globalny (rozdział 4.8) | kontrakt (pola kondygnacji w `ConfigGetRequest`) oraz rdzeń |
| **oś dostępu ponad kartą** | brak. Nadania istnieją wyłącznie per `windowId` (rozdział 4.9) | kontrakt oraz rdzeń |
| **karta incognito** | brak. Słowo `incognito` nie występuje ani w kontrakcie, ani w kliencie | kontrakt oraz warstwa kart |
| **gałąź robocza karty** | brak. `workspaceBranch`, `cardBranch` — zero trafień w kontrakcie | kontrakt oraz warstwa kart |
| **wstążka narzędziowa okna roboczego** (decyzja 9) | brak. `toolbar`, `ribbon` — zero trafień w kontrakcie. `moduly/studio/wstazka-pracy.ts` jest wstążką dokumentu w module Studio, nie wstążką karty | nowa warstwa widoku okna roboczego |
| **przybornik przeglądarki i inspektor elementów** (decyzja 9) | brak. `inspector` — zero trafień w kontrakcie. `moduly/design/inspektor-warstwy.ts` bada warstwę kompozycji projektowej, nie stronę w karcie przeglądarki | karta przeglądarki wraz z jej wstążką |
| **menu wskazanego elementu** (decyzja 9) | brak | jw. |

Nie brakuje natomiast: trybów upoważnienia (kontrakt niesie ich sześć, w tym
wszystkie cztery nazwane decyzją 6), izolacji na poziomie karty (`ConfigScope.Window`
jest wśród poziomów przyjmowanych przez `isolation.*`, co wykazuje pomiar
w rozdziale 4.7), zakresu nadawanego wskazaniem miejsca (`roots` punktu
i podzbiór korzeni nadania) oraz rozróżnienia „nieustawione" od „ustawione"
(`Rozstrzygniecie.domyslna`).

## 7. Usterki źródeł — do rozstrzygnięcia przez Prowadzącego

Poniższe stwierdzono pomiarem. Żadnej z tych rzeczy ocena nie uzupełnia
domysłem.

1. **Liczba stu trzech tysięcy wierszy nie ma pokrycia w pomiarze.** Klient liczy
   242 267 wierszy w 1311 plikach `.ts`; wierszy w plikach niebudujących dokumentu
   jest 85 142 miarą węższą albo 67 964 miarą szerszą. Dziewięć warstw wykazu
   niesie 24 055 wierszy, z czego 9 024 miarą węższą i 6 597 miarą szerszą.
   Skąd bierze się sto trzy tysiące — nie wiadomo; wykaz wymaga uzgodnienia
   z pomiarem.

2. **Kontrakt niesie sześć trybów upoważnienia, decyzja 6 nazywa cztery.**
   Nadmiarowe: `acceptEdits`, `dontAsk`. Rdzeń przyjmuje i zapisuje wszystkie
   sześć. Klient buduje listę wprost z wyliczenia kontraktu, więc Operator
   zobaczy sześć pozycji. Do rozstrzygnięcia, czy zwężamy kontrakt, czy
   rozszerzamy decyzję.

3. **Kontrakt niesie dziewięć poziomów zasięgu, decyzja 6 nazywa cztery
   kondygnacje.** Stosunek pięciu pozostałych — `global`, `environment`,
   `module`, `modulePair`, `role` — do kaskady jest nierozstrzygnięty. Rzecz nie
   jest teoretyczna: `global` jest dziś jedynym poziomem obowiązującym ponad
   punktem widzenia.

4. **Nazwy `Session` i `Window` są w przejętym kliencie odwrócone względem
   decyzji 5.** Rozstrzygnięcie jest jedno dla całej budowy i należy do rejestru
   decyzji, nie do wykonawcy terenu (rozdział 5).

5. **Kaskada zakresu nie działa między kondygnacjami — po żadnej ze stron.**
   Do rozstrzygnięcia, czy rachunek dziedziczenia prowadzi rdzeń (wtedy
   `ConfigGetRequest` potrzebuje pól projektu i sesji, a `kontekstZasiegu` —
   wypełnienia), czy klient (wtedy potrzebuje drogi do poznania ścieżki bytów,
   której — jak sam stanowi — nie zgaduje).

6. **Projektu nie da się założyć żadną komendą kontraktu.** Typu `Project` nie
   ma; komend `project.*` nie ma. Kondygnacja projektu z decyzji 6 i poziom
   projektu z decyzji 5 nie mają w kontrakcie na czym stanąć.

7. **Rdzeń liczy dziewięć poziomów, a mówi o ośmiu.**
   `server/internal/konfig/poziomy.go` wylicza dziewięć i tak nazywa je
   w komentarzu funkcji `Znany`, lecz zdania odmowy w
   `server/internal/core/adapter_modul_isolation.go` (wiersz 255) i
   `server/internal/core/adapter_konfiguracja_prowenancja.go` (wiersz 327) mówią
   „nie należy do **ośmiu** poziomów kontraktu". To samo zdanie niesie opis
   `isolation.policy.preview` w kontrakcie. Odmowa trafia do Operatora, więc
   liczba w niej ma znaczenie.

8. **Rdzeń wymaga bytu dla poziomu `application`, kontrakt stanowi, że bytu tam
   nie ma.** Kontrakt przy `ConfigScope.Application`: „zapisu na nim nie ma czym
   zawęzić, bo bytu nie ma (scopeId pusty)". Pomiar:
   `isolation.context.get {scope: application}` → `validation_failed`
   „poziom zasięgu application wymaga wskazania bytu — zapis bez niego nie
   doszedłby do rozstrzygania". Dwa źródła sprzeczne.

9. **Rdzeń przyjmuje żądania niepełne i wartości spoza wyliczeń, milcząc.**
   `access.point.add` bez pól obowiązkowych `kind`, `roots` i `defaultMode`
   odpowiada powodzeniem i zakłada punkt rodzaju `mcpBridge`.
   `window.update {permissionMode: "tryb-spoza-kontraktu"}` odpowiada
   powodzeniem i zostawia `manual`. Wynik wyglądający dobrze jest tu brakiem
   sprawdzenia podanym jako wynik.

10. **`polaczenie/gniazdo.ts` wiersz 46 wpada w nieskończoną rekurencję pod Node.**
    Szczegół w rozdziale 4.1. Pod Chromium ta ścieżka jest poprawna.

11. **`protokol/uzgodnienie.ts` zakłada sesję przy nawiązaniu połączenia**, co
    jest sprzeczne z decyzją 5, oraz **ogłasza powodzenie etapu `sesja` przed
    zapisaniem identyfikatora sesji**. Szczegół w rozdziale 4.2.

12. **`sterowanie/atrapa-rdzenia.ts` jest kodem martwym z nagłówkiem mówiącym
    nieprawdę** — 211 wierszy bez ani jednego wołacza, przy nagłówku
    stwierdzającym istnienie wołaczy w plikach sprawdzianów, których ta warstwa
    nie ma.

13. **`posijKomende` jest w `punkty-izolacji/` zdefiniowana czterokrotnie**, w tym
    raz w postaci nietypowanej, znoszącej wiązanie żądania z odpowiedzią, które
    daje kontrakt. Szczegół w rozdziale 4.7.

14. **`ConfigScope.Application` jest pominięty w obu wykazach poziomów po stronie
    klienta** — `konfiguracja/zasiegi.ts` (`ZASIEGI_OD_NAJWEZSZEGO`) oraz
    `punkty-izolacji/katalog-izolacji.ts` (`ZASIEGI_OD_NAJSZERSZEGO`) — mimo że
    oba pliki niosą dla niego etykietę, a rdzeń oddaje go w wykazie poziomów.

15. **Siedem z dziewięciu warstw nie ma ani jednego sprawdzianu.** Sprawdziany
    ma `polaczenie` (23) i `punkty-izolacji` (4). Przeniesienie warstwy bez
    sprawdzianów odbywa się bez sieci zabezpieczającej i to jest powód, dla
    którego zdania „przechodzi bez zmian" nie postawiono przy `modele`, mimo że
    warstwa ta nie ma żadnej zależności od modelu okna.

## 8. Czego nie dało się sprawdzić

1. **Zachowania przejętych warstw wewnątrz powłoki Tauri.** Pomiar prowadzono
   pod Chromium i pod Node. Ścieżka `adres-rdzenia.ts` obsługująca pakiet
   osadzony (pochodzenie `tauri.localhost`, polecenie `adres_rdzenia`) ma
   sprawdziany jednostkowe, lecz nie została uruchomiona w powłoce natywnej.
   Wymaga zbudowania powłoki, co wykracza poza teren.

2. **Znaczenia pustego zbioru nadań i pustego wykazu katalogów roboczych.**
   Nowe okno powstaje z `workingDirs: []`, a kontrakt niesie stan degradacji
   `outsideGrantedRoots` — „Katalog ustawiony lezy poza korzeniami nadan dostepu
   okna". Czy pusty zbiór nadań znaczy w rdzeniu zakres pełny (jak wymaga
   decyzja 6), czy zakres pusty — nie zmierzono. Rozstrzygnięcie wymaga
   doprowadzenia rdzenia do wykonania czynności modelu, czyli uruchomienia toru
   wykonawczego z prawdziwym kanałem modelu; rdzeń w tej próbie zamknął tor
   wykonawczy natychmiast po starcie („tor wykonawczy zakończony: wyczerpane
   wejście").

3. **Pełnej treści polityki efektywnej obszarów sesji.** `config.effective.get`
   odpowiada polami `config`, `origins`, `workingDirectory`, `resolvedAt`;
   zawężenie żądania do obszaru `permissions` nie wystawia tego obszaru na
   poziomie wierzchnim odpowiedzi. Które obszary rdzeń wypełnia naprawdę, a które
   zwraca puste — wymaga osobnego pomiaru z założonymi zapisami w każdym z
   osiemnastu obszarów.

4. **Zachowania warstw wobec rdzenia z włączonym wymogiem logowania.** Rdzeń
   uruchamiano z `-wymog-logowania false`, zgodnie ze wskazaniem terenu. Droga
   `rozpoznanie-bramki.ts` → `ekran-logowania.ts` → `sesja-bramki.ts` przy
   włączonym wymogu nie została przebyta.

5. **Rozstrzygnięcia dla logiki spoza wykazu.** Największe skupisko logiki
   niebudującej dokumentu leży w `moduly/` (47 159 wierszy). Przedmiotem oceny
   jest wykaz dziewięciu warstw, więc masa ta pozostaje nieoceniona i wymaga
   własnego terenu.
