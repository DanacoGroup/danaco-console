# Przekazanie prowadzenia budowy — stan na 1 września 2026

Dokument jest dla sesji, która przejmuje prowadzenie budowy Danaco Console po sesji
z 30 sierpnia — 1 września 2026. Nie streszcza repozytorium; mówi, co zastajesz,
czego nie wolno ruszyć i gdzie leży praca w biegu.

## 1. Czym rzecz stoi

Rdzeń jest gotowy, interfejsu nie ma. Liczby z audytu z 1 września (pomiar, nie ocena):

| miara | wartość |
|---|---|
| komendy kontraktu | 1086 (przed naprawami 1083) |
| komendy z uchwytem w rdzeniu | wszystkie |
| komendy wołane przez klienta | 55 (przed naprawami 46) |
| rodziny bez ani jednego wołania | 52 z 67 |
| zdarzenia rdzenia z odbiorcą w kliencie | 5 z 64 przed naprawami; po naprawach rozdzielacz w `klient/src/wiazanie/zdarzenia.ts` |
| prototypy okien w `design/05-okna/` | 35 |
| okna obecne w produkcie | 8 |
| miejsca, w których rdzeń pyta „kto woła" | 12 przed naprawami |

Poprzedni klient (`budowa/klient-poprzedni/`) wołał 1061 komend. To materiał do czytania
i przeszczepu — **nie do rozwoju i nie do budowania**; jego budowanie w drzewie kasowało
pliki terenów roboczych (`emptyOutDir`).

## 2. Gałęzie i dokumenty

| co | gdzie |
|---|---|
| gałąź prac | `teren/naprawy-audytu` |
| stan sprzed napraw | rewizja `2881c425` na `teren/centrum-poprawki` |
| `main` | drzewo robocze `~/robocze/prowadzenie/main-roboczy`, scalone do rewizji `ed7be7c0` |
| audyt pierwszy | `~/robocze/audyt-danaco-console-2026-09-01.md`, wystawiony pod `http://51.75.62.180/wykaz-audyt.html` |
| materiał przekazania | `~/robocze/przekazanie-2026-09-01/` — audyt, sprawozdania terenów, trzy skrypty przebiegów |

Audyt pierwszy jest podstawą całej dalszej pracy: niesie 189 ustaleń z plikiem i wierszem,
recepty naprawcze, plan siedmiu etapów z miarami odbioru oraz wykaz dwunastu rozstrzygnięć
czekających na Właściciela. Zanim cokolwiek zrobisz — przeczytaj go w całości.

## 3. Maszyny

**Maszyna budująca** (ta, na której stoisz). Podgląd budowy wychodzi na świat pod
`http://51.75.62.180/`; tylko port 80. Pliki podglądu: `/srv/podglad/`, plik wystawiony
pod nazwą `wykaz-*.html` albo `wykaz-*.exe`. Port 80 obsługuje Caddy i kieruje na
`127.0.0.1:17896` — tam stoi rdzeń podglądowy uruchomiony z `/tmp/rdzen-podglad`
z katalogiem danych `~/robocze/podglad-dane` i pakietem interfejsu z `budowa/klient/dist`.
**To proces, nie usługa** — po restarcie maszyny trzeba go postawić ręcznie.

**danaco-system — 57.128.253.74.** Rdzeń wdrożenia (`danaco-console.service`, port 17870,
katalog `/opt/danaco/danaco-console`, dane `/var/lib/danaco-console`), rdzeń wydania Studio
(`danaco-console-studio.service`, port 17871), kanał pobrań `pobierz.danaco-group.pl` pod Caddy,
portal `/opt/danaco/danaco-console/portal`, serwer poczty Stalwart. Dostęp: klucz
`~/.ssh/danaco_operator`, użytkownik `root`. Hasło skrzynki nadawczej stoi w
`/etc/danaco-console/srodowisko` na tej maszynie — **nie w repozytorium**; rdzeń złożony
bez niego wchodzi w drogę bez poczty i rejestracja milczy.

Rdzeń szuka pakietu interfejsu w `klient/dist` (po polsku) względem katalogu roboczego.
Katalog `client/dist` z poprzedniego wdrożenia zostaje jako kopia — dziennik startu ma
mówić `klient=klient/dist`, a nie `klient=brak katalogu`.

**Maszyna Windows do sprawdzania.** Kontener `danaco-win` (dockurr/windows), noVNC na
`http://127.0.0.1:8006/`, dostęp przez `sudo docker`. Katalog wspólny
`/srv/win-danaco/wspolny` widziany w Windows jako `\\host.lan\Data` — to jedyna wygodna
droga wniesienia pliku do maszyny. Sterowanie: Playwright przez CDP,
`executablePath: '/usr/bin/chromium-browser'`, kliknięcia po współrzędnych na kanwie noVNC.
**Pułapka:** `Alt+F4` na pulpicie wywołuje okno zamknięcia systemu — kontener padnie
i trzeba go wystartować ponownie.

## 4. Co ta sesja zrobiła

**Instalator.** Nie działał, bo binarium celu `x86_64-pc-windows-gnu` niesie import
`WebView2Loader.dll`, a pobierany sam plik nie ma go obok siebie. Przebudowany celem
`x86_64-pc-windows-msvc` (`cargo xwin`, `WebView2LoaderStatic.lib`, `-C target-feature=+crt-static`) —
w imporcie zostały wyłącznie biblioteki systemu. Wykazane w Windows z pliku pobranego z serwera:
sześć kroków kreatora do końca.

**Okna bez ramy.** Instalator i aplikacja miały dwie belki naraz — systemową i własną.
`decorations(false)` w obu, belka warstwy projektowej związana z oknem powłoki
(`klient/src/wiazanie/belka-okna.ts`), okno powłoki wstaje ukryte i pokazuje się dopiero
z pierwszą klatką ekranu startowego.

**Animacja startowa.** Nie było jej z powodu usterki, nie opóźnienia: `ekran-startowy.js`
mierzył pole raz, przy montażu, i gdy arkusze nie zdążyły nadać mu wymiaru, kanwa zostawała
**1×1 piksela**. Dołożony obserwator wymiaru. Wykazane zrzutami z Windows.

**Kanał pobrań.** Składniki, po które sięga krok 5 instalatora, przeniesione pod otwartą
ścieżkę `pobierz.danaco-group.pl/pliki/` — kreator sześciu kroków nie ma pola na hasło,
więc plik za uwierzytelnieniem był dla niego nie do pobrania. Uwierzytelnienie zostaje
na `/wydania/*`.

**Wdrożenie.** Na danaco-system stała wersja rdzenia z 19 sierpnia. Podmieniona na bieżącą
wraz z serwerem narzędzi, pakietem interfejsu i pomocnikami; ten sam rdzeń wgrany do wydania
Studio. Arsenał uzupełniony z 50 do **55 z 55** (hunspell ze słownikami, Semgrep w osobnym
środowisku `/opt/danaco-arsenal/semgrep`, typos, typst, vale). Poprzednie binaria leżą obok
jako `danaco-console.bak-2026-08-30`.

**Audyt.** Trzynaście kategorii, 64 agenty, weryfikacja adwersaryjna na ustaleniach najcięższych.
189 ustaleń, 18 krytycznych (5 obalonych przez weryfikację).

**Naprawy.** Czternaście terenów o rozłącznych wykazach plików, praca równoległa.
Zapis: rewizja `955c03e0` — 94 pliki, +6490/−1423, drabina zdana.

**Scalenie i kontrola odbioru (sesja druga, 1 września).** Zgłoszenia „poza terenem"
domknięte w zakresach A–L, każdy własną rewizją (`db8d8eff` … `92e82cad`); kontrola
odbioru w pięciu zakresach; klucz sejfu rozstrzygnięty (26, 30); jedna maszyna i nazwa
wdrożenia (29); certyfikat podpisu (27); ARM64 poza etapem (28). Pełny bieg sprawdzianów
rdzenia: 662 przechodzą w 12 minut (limit domyślny 10 minut jest za krótki — bieg z `-timeout 20m`).

**Domknięcie (sesja trzecia, 1–2 września).** Zakres N — kanał klienta po zerwaniu
(`082ab081`): wykrywanie uśpienia skokiem zegara (próg 30 s = ping rdzenia 20 s + 10 s
czekania), nasłuch `online`/`visibilitychange`, sufit kolejki 256 ramek z powodem
porzucenia `przepelnienie`, odpowiedź po terminie korelacji odrzucana do dziennika.
Zakres M — 56 emisji zdarzeń rdzenia bez konta, trzy strumienie rozgłaszane do wszystkich
gniazd, wyścig przy zatrzymaniu tury. Skrypt `instalka-kreatora-win-x64.sh` przepisany
na cel `msvc` ze statycznym CRT i zaporą importów — poprzednia postać składała kreator
celem `gnu` i pakowała go w instalkę NSIS, czyli w to, co Właściciel odrzucił.

## 5. Praca w biegu — czym zaczyna następna sesja

Naprawy z audytu wykonało czternaście terenów pracujących równolegle na rozłącznych wykazach
plików (`955c03e0`); scalenie i kontrolę odbioru wykonała sesja druga (rozdział 4).
Rozdziały 5.1–5.3 zostają jako opis metody — audyt powtórny (5.3) wciąż nie był uruchomiony
i jest pierwszą rzeczą po wykazaniu wydania w Windows.

Materiał leży w `~/robocze/przekazanie-2026-09-01/`:

| plik | co niesie |
|---|---|
| `audyt-pierwszy.md` | audyt z 1 września: 189 ustaleń, recepty, plan siedmiu etapów |
| `sprawozdania-terenow.md` | co zrobił każdy teren, co pominął i czego nie mógł tknąć |
| `sprawozdania-terenow.json` | to samo w postaci do przetworzenia |
| `workflow-audyt-pierwszy.js` | przebieg audytu — do powtórzenia bez pisania od nowa |
| `workflow-naprawy.js` | przebieg napraw wraz z wykazami plików czternastu terenów |
| `workflow-audyt-powtorny.js` | przebieg audytu powtórnego, napisany i nieuruchomiony |

Bilans czternastu terenów: **145 napraw wykonanych, 37 pozycji pominiętych** (czekają
na rozstrzygnięcia Właściciela albo na pracę projektową) i **86 zgłoszeń „poza terenem"**.

### 5.0 Rejestracja: zderzenie terenów rozstrzygnięte 1 września

Teren trwałości zaszyfrował sejf poświadczeń kluczem ze zmiennej
`DANACO_KLUCZ_SEJFU`, a zakładanie konta zapisuje skrót hasła w sejfie — bez
zmiennej świeża instalacja nie zakładała żadnego konta. Rozstrzygnięcie 26
w `decyzje.md`: klucz ze zmiennej, a bez niej klucz własny rdzenia zakładany
w katalogu danych przy pierwszym użyciu. Sprawdzian
`TestKontoDrugieNieSiegaKontaPierwszego` przechodzi. Na serwerze wdrożenia
zmienna ma wskazać plik poza katalogiem danych — to wchodzi z wdrożeniem
(rozdział 5.1, jednostka systemd i `/etc/danaco-console/srodowisko`).

### 5.0a Wdrożenie z 1/2 września — co stoi i czego pilnować

Rdzeń po naprawach stoi na danaco-system w obu usługach (`danaco-console` 17870,
`danaco-console-studio` 17871): `komend=1086`, `55 z 55`, wymóg logowania włączony,
klucz sejfu ze zmiennej (`/etc/danaco-console/sejf.klucz`, dla Studia `sejf-studio.klucz`
przez drop-in `danaco-console-studio.service.d/bramka.conf`). Podgląd na 51.75.62.180
ma ten sam rdzeń. Wydanie 2.0.0 z 1 września leży w kanale pod `/pliki/2.0.0-2026-09-01/`
i jest wykazane w Windows z pliku pobranego z kanału: kreator (6 597 632 B, samodzielny,
cel `msvc`) → sześć kroków → założona powłoka (2 598 033 B, `https://console.danaco-group.pl:443`)
→ start z kroku 6 → `przyłączenie` w dzienniku danaco-system → po restarcie rdzenia aplikacja
wraca sama (`przyłączenie` sekundę po `serwer gotowy`).

**Dwie rzeczy, które o mało nie położyły wdrożenia — obie czekają na naprawę w kodzie:**

1. **Baza wdrożenia idzie inną linią numeracji migracji.** Powstała 18–19.08 rdzeniem
   wersji poprzedniej: jej wiersze 406–480 nazywają się `awatar_konta_wlasciciela`,
   `rejestr_urzadzen_konta`, `katalog_okna_ustawien`… i nie odpowiadają krokom
   repozytorium. Nowe kroki dołożone przy naprawach jako 408 i 409 zderzyły się z jej 408 —
   przemianowane na **481 i 482** (rewizja `7f5e5730`). Każdy następny krok ma numer
   wyższy od 482. Bazy podglądu i Studia, które dziś zdążyły przyjąć 408/409 pod nową
   treścią, cofnięto ręcznie (usunięte wiersze, kolumny `proby` i `klucz_hosta`, indeks);
   kopie `*.bak-2026-09-01-2345` leżą obok.
2. **`uzgodnijSumyKontrolne` w `store/migracje.go` przepisuje sumy po numerze, nie po treści.**
   Przy pierwszym starcie z `PRAGMA user_version = 0` każda zastosowana wersja dostaje sumę
   z bieżącego repozytorium — także wtedy, gdy krok o tym numerze ma zupełnie inną treść.
   Tak baza wdrożenia „miała" kroki 406/407 (`konto_id` w tabelach bramki i kart sesji),
   których DDL nigdy na niej nie wykonano; rdzeń padał na 481 (`no such column: konto_id`).
   DDL kroków 406 i 407 wykonano na tej bazie **ręcznie** (kopia `danaco-console.db.bak-2026-09-01-2340`
   w `/var/lib/danaco-console/`). Uzgodnienie sum ma porównywać nazwę kroku, nie sam numer,
   a rozjazd nazwy ma być odmową startu — to wchodzi do audytu powtórnego jako ustalenie.

Stan kanału: `zmierz.mjs --sprawdz` — rozjazdów 0; pozycja `serwer` (pakiet `.deb` z 18.08 za
hasłem) wciąż bez pliku na dysku i bez rozstrzygnięcia (audyt, rozdział 5, pozycja 4).

### 5.1 Scalenie

Osiemdziesiąt sześć zgłoszeń „poza terenem" to miejsca, których teren nie mógł tknąć, bo plik należał
do sąsiada pracującego w tej samej chwili. Teraz nikt nie pracuje równolegle, więc wolno tknąć
każdy plik. Wykaz stoi w `sprawozdania-terenow.md`, w blokach „Do domknięcia poza terenem".

Kolejność: najpierw doprowadzić drzewo do budowy, potem domknąć zgłoszenia, potem znów zbudować.

```bash
cd ~/budowa/budowa && go build ./... && go vet ./server/... && gofmt -l server
cd ~/budowa/budowa/klient && npx tsc --noEmit && npx vite build
cd ~/budowa/budowa/desktop/src-tauri && PATH=$HOME/.local/bin:$PATH cargo check --target x86_64-pc-windows-gnu
cd ~/budowa/budowa/instalator/src-tauri && PATH=$HOME/.local/bin:$PATH cargo check --target x86_64-pc-windows-gnu
cd ~/budowa/budowa && go test ./server/internal/dane/... ./server/internal/store/... ./server/internal/transport/...
```

Pakiet `server/internal/core` pomiń — ma sprawdzian sięgający po silniki zewnętrzne,
który nie kończy się w dziesięć minut. To osobne ustalenie audytu.

### 5.2 Kontrola odbioru

Sprawozdaniom terenów nie wolno wierzyć. Kontrola ma czytać kod po naprawie i mierzyć,
z założeniem, że naprawa jest pozorna, dopóki pomiar nie pokaże inaczej. Pięć zakresów:

1. **granica konta** — K1, W4, W5, W7: czy konto B nie sięga po dane i hasło konta A
2. **kanał i bramka** — W3, W6, W10, W26: czy bramki nie da się zdjąć z gniazda, czy unieważniona
   sesja traci dostęp natychmiast, czy rozgłoszenie trafia wyłącznie do konta wołającego
3. **model okna** — K2, W11–W15: czy sesja wraca z historią, czy dwie karty Studia żyją obok siebie
4. **Studio i atrapy** — K3, W1, W17, W18, W21, W22: czy strumień dochodzi, czy turę da się zatrzymać,
   czy w wydaniu nie została ani jedna atrapa
5. **instalator i wydanie** — K5, K6, K7, W24, W27: czy instalator naprawdę zakłada program

### 5.3 Audyt powtórny

Dopiero po kontroli. Skrypt jest napisany (`workflow-audyt-powtorny.js`) i ma cztery fazy:
odbiór wszystkich ustaleń pierwszego audytu jedno po drugim, trzynaście kategorii na nowo
z osobnym oznaczaniem **regresji**, weryfikacja adwersaryjna, synteza z bilansem.

Bilans ma odpowiedzieć na jedno pytanie: ile ustaleń pierwszego audytu naprawdę zniknęło,
a ile tylko przemalowano.

## 6. Co zostało

Plan siedmiu etapów stoi w audycie, rozdział 4, każdy z wykazem czynności plik po pliku,
miarą odbioru i kosztem. Kolejność wynika z zależności, nie z wagi:

1. fundament kanału — ponowienie powitania po zerwaniu, terminy korelacji, heartbeat, rozdzielacz zdarzeń
2. bramka, granica konta i brama kontraktu
3. model okna i sesji w kliencie
4. Studio jako okno pracujące
5. Centrum jako zarząd sesji
6. kontrakt, trwałość i sprawdzalność
7. wdrożenie i wydanie

Etapy 1 i 2 idą równolegle — pliki rozłączne. Etapów 3–5 nie da się zamknąć przed 1.
Etap 7 nie ma sensu przed 1, 2 i 6.

Po naprawach z tej sesji część czynności etapów 1, 2, 4, 5 i 7 jest wykonana; **stan każdego
ustalenia ma rozstrzygnąć audyt powtórny**, nie sprawozdania terenów.

## 7. Czego nie wolno

- **Nie edytować migracji 226, 269 i 378.** Niosą wadę (kaskadowe kasowanie przy przebudowie
  tabeli), ale strażnik sumy kontrolnej wywróci start każdej istniejącej bazy. Naprawa idzie
  wyłącznie nowym krokiem.
- **Nie odtwarzać prototypów.** 27 arkuszy i 79 skryptów w `design/zasoby/` oraz 35 prototypów
  okien w `design/05-okna/` są biblioteką do wpięcia. Kompozycji okna nie przekazuje się zleceniem
  i nie zmienia bez wyraźnego polecenia Właściciela.
- **Nie budować `klient-poprzedni`** w drzewie repozytorium.
- **Nie stawiać nowych testów, bramek i walidatorów** — reżim audytów cyklicznych został zniesiony;
  buduje się produkt i pokazuje wynik. Sprawdzian pisze się wtedy, gdy jest jedynym sposobem
  wykazania, że naprawa działa.
- **Nie pytać Właściciela o rzeczy rozstrzygalne ze źródeł.** Rozstrzygnięcia wiążące stoją
  w `prowadzenie/decyzje.md`; szczególnie rozstrzygnięcie 5 (model okna: aplikacja, okno robocze,
  karta, panel) i reguła, że sesja powstaje przy pierwszej wiadomości, nie przy wejściu w moduł.
- **Nie wystawiać niczego jako gotowego bez pokazania wyniku.** Zrzut z maszyny Windows
  z pliku pobranego z serwera, a nie z dysku maszyny budującej.

## 8. Rozstrzygnięcia czekające na Właściciela

Pełny wykaz stoi w audycie, rozdział 5 (12 pozycji). Trzy blokujące wydanie
zostały rozstrzygnięte 1 września z upoważnienia Właściciela („Rozstrzygaj za
mnie") i zapisane w `prowadzenie/decyzje.md`:

| pozycja audytu | rozstrzygnięcie | treść |
|---|---|---|
| 1. certyfikat podpisu kodu | 27 | OV od Certum, klucz w SimplySign, znacznik `http://time.certum.pl`; do zakupu wydania idą z `DANACO_PODPIS=pomijany` |
| 2. jedna maszyna i jeden port wdrożenia | 29 | `danaco-system` 57.128.253.74, rdzeń `127.0.0.1:17870` z wymogiem logowania, na świat `console.danaco-group.pl:443` przez Caddy; wydanie z `DANACO_HOST_WDROZENIA=console.danaco-group.pl`, `DANACO_PORT_WDROZENIA=443`, `DANACO_SCHEMAT_WDROZENIA=https` |
| 3. nazwa i certyfikat TLS dla rdzenia | 29 | nazwa i certyfikat ACME już stoją na Caddy; rdzeń nie potrzebuje `DANACO_TLS_*` |
| 5. wydanie ARM64 | 28 | nie wchodzi w ten etap; pozycja zostaje „w przygotowaniu" bez pliku |
| 8. domyślna wartość bramki na pętli zwrotnej | 14, 29 | domyślna zostaje; wdrożenie ustawia wymóg jawnie w jednostce systemd |
| 10. czym szyfrować sejf poświadczeń | 26, 30 | plik klucza ze zmiennej na serwerze, klucz własny rdzenia bez niej, na Windows DPAPI |

Wykaz, pakiet, skrypty składania i dokumenty doprowadzone do tych wartości
1 września (`wydania.json`, `danaco-console.service`, `srodowisko`,
`DEBIAN/control`, `kanal-wydan.md`, `INSTALACJA-I-KONFIGURACJA.md`).

**Zostaje po stronie Właściciela — czynność, nie rozstrzygnięcie:**

- zakup certyfikatu OV w Certum na dokumenty spółki i założenie konta
  SimplySign; po zakupie składanie idzie z `DANACO_PODPIS=wymagany`.

**Pozycje audytu bez rozstrzygnięcia:**

4. kanał `/wydania/` chroniony hasłem — utrzymać obok `/pliki/`, czy zwinąć;
6. kształt rodziny `project.*` w kontrakcie;
7. reguła wykazu narzędzi modelu (603 komendy poza wykazem bez podanej racji;
   rozstrzygnięcie 20 nazywa kierunek, nie treść reguły);
9. kolejność wpinania podgrup `studio.*` i projekt brakujących kontrolek;
11. los siedmiu okien platformowych;
12. zgoda na naprawę migracji nowym krokiem.

## 9. Jak sprawdzić, że wszystko stoi

```bash
cd ~/budowa/budowa && go build ./... && go vet ./server/... && gofmt -l server
cd ~/budowa/budowa/klient && npx tsc --noEmit && npx vite build
cd ~/budowa/budowa/desktop/src-tauri && PATH=$HOME/.local/bin:$PATH cargo check --target x86_64-pc-windows-gnu
python3 ~/budowa/narzedzia/pokrycie-kontraktu.py
ssh -i ~/.ssh/danaco_operator root@57.128.253.74 'journalctl -u danaco-console -n 5 --no-pager'
curl -s -o /dev/null -w '%{http_code}\n' https://console.danaco-group.pl/
```

Dziennik startu rdzenia ma mówić `komend=1086`, `zależności zewnętrzne: 55 z 55 obecnych`
i `klient=klient/dist`.

## 10. Pułapki tej maszyny

- Cel `x86_64-pc-windows-msvc` buduje się przez `cargo xwin`, a `cc-rs` szuka `llvm-lib`
  bez przyrostka wersji. W `~/.local/bin` stoją dowiązania `llvm-lib`, `llvm-rc`, `llvm-ar`,
  `llvm-dlltool`, `llvm-objcopy` do plików `*-21`. Bez `PATH=$HOME/.local/bin:$PATH` budowa pada.
- Powłoka główna składa się celem `gnu` i wymaga `WebView2Loader.dll` obok siebie; niesie ją
  instalka NSIS. Instalator kreatora składa się celem `msvc` i jest samodzielnym plikiem.
- Playwright chodzi wyłącznie z `executablePath: '/usr/bin/chromium-browser'` i `--no-sandbox`.
- Katalog tymczasowy sesji (`/tmp/claude-1000/...`) znika wraz z sesją. Co ma przetrwać,
  ląduje w `~/robocze/` albo w repozytorium.
