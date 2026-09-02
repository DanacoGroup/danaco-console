# Przekazanie prowadzenia budowy — stan na 2 września 2026

Dokument opisuje stan zastany, nie przebieg prac. Mówi, co stoi, co jest otwarte
i czego nie wolno ruszyć.

## 1. Czym rzecz stoi

Rdzeń działa. Interfejs prowadzi od wejścia przez przedsionek do Centrum i Studia;
pozostałe moduły stoją na zapowiedzi niegotowości. Praca leży na gałęzi `teren/naprawy-audytu`;
`main` i punkt cofnięcia `2881c425` są nietknięte.

| miara | wartość |
|---|---|
| komendy kontraktu | 1086 |
| komendy wołane przez klienta | 86 |
| tabele w bazie po przejeździe migracji | 418 |
| tabele niosące pracę Operatora | 373, w tym 172 korzenie |
| korzenie z granicą konta w zapytaniach | 47 najcięższych i 53 pliki rundy trzeciej |
| kroki migracji | do 485 |
| prototypy okien w `design/05-okna/` | 35; okien w produkcie 8 |
| droga wejścia | rejestracja, logowanie, metody, urządzenia, sesje bramki |

## 2. Gałęzie, dokumenty, materiał

| co | gdzie |
|---|---|
| gałąź prac | `teren/naprawy-audytu` |
| stan sprzed napraw | rewizja `2881c425` na `teren/centrum-poprawki` |
| rejestr rozstrzygnięć | `prowadzenie/decyzje.md` — 33 pozycje, wiążą bezwzględnie |
| audyt pierwszy | `~/robocze/przekazanie-2026-09-02/audyt-pierwszy.md` — 189 ustaleń, plan siedmiu etapów |
| audyt powtórny, część | `~/robocze/przekazanie-2026-09-02/audyt2-czesciowy.md` — 7 z 13 kategorii, 79 ustaleń |
| materiał granicy konta | `~/robocze/przekazanie-2026-09-02/` — rozpoznanie korzeni, pozycje otwarte, przebiegi |

Audyt powtórny stanął na 7 z 13 kategorii. **Nie ruszyły**: rdzeń Go, klient
TypeScript, warstwa projektowa, powłoka i instalator, poczta i droga wejścia,
braki produktowe. Skrypt przebiegu leży w materiale przekazania.

## 3. Granica konta — rzecz największa i niedokończona

**Miara.** Platforma prowadzi wiele kont od migracji 406. Granica konta objęła do
tej sesji cztery tabele: trzy tabele bramki i karty sesji. Rozpoznanie wszystkich
418 tabel wykazało **373 tabele niosące pracę Operatora**; z nich **172 to
korzenie** (wiersz należy do Operatora wprost), 178 to dzieci wiszące na
korzeniach kluczem obcym, 22 to łączniki.

**Mechanizm.** Wskazanie konta stoi w korzeniu; dziecko granicę dziedziczy drogą
po kluczu obcym. Warunek jest jeden na całą bazę:

- `dane.WarunekKonta` — zawężenie wiersza do konta żądania, jeden argument
  zapytania: `dane.KontoOperatora(ctx)`. Porównanie idzie przez `IS`, więc
  instalacja przed rejestracją widzi pracę zastaną.
- `dane.WskazanieKonta` — zapis konta w wierszu zakładanym.
- Wzór dla dziecka: `AND EXISTS (SELECT 1 FROM <korzeń> k WHERE k.id =
  <dziecko>.<klucz_obcy> AND <WarunekKonta z aliasem k>)`.

Konto wchodzi do kontekstu raz, w `core/rdzen.go` przed rozdziałem na
obsługiwacze, więc warstwa danych bierze je z kontekstu, nie argumentem.

**Kroki migracji tej sesji.** 484 dokłada `konto_id` 171 korzeniom (`kanal_modelu`
ma własny). 485 rozstrzyga zderzenie nazw — `kanal_modelu.konto_id` znaczyło
konto dostawcy i nosi teraz nazwę `konto_dostawcy_id` — oraz przebudowuje siedem
wskaźników jednoznacznych na kontowe.

**Co zrobione.** Zapytania 47 korzeni o najwyższej wadze — magazyny poświadczeń
i treści Operatora: `punkt_dostepu`, `terminal_klucz`, `terminal_host`,
`host_zdalny`, `konto`, `kanal_modelu`, `skrzynka_pocztowa`,
`sekret_rozszerzenia`, `poswiadczenie_automatyki`, `dokument_studio`,
`blok_wiadomosci`, `plik_biblioteki` i pozostałe. Zawężonych zapytań: 301
w pierwszym przejściu, 106 dziur domkniętych w drugim.

**Runda trzecia (2 września, wieczór).** Trzynaście terenów domknęło pozycje
otwarte 53 plików i zapisało konto w każdym `INSERT` do korzenia w tych plikach
(rewizja `a9ea0447`); sprawozdania wykonawców co do pliku i stałej:
`~/robocze/przekazanie-2026-09-02/runda-3/wyniki-terenow-01-13.txt`, instrukcja
dokończenia: `runda-3/STAN.md`. Kontrolerzy terenów NIE zdążyli — sesja stanęła na
limicie; rewizja `a9ea0447` przeszła drabinę szybką (budowa, vet, walidator), nie
przeszła sprawdzianów ani kontroli. Teren 14 (nastawy per konto przez
`konfig.Kontekst.KontoOperatora`) przerwany w połowie, stan zapisany rewizją `abc78ea8`
— rozstrzygacz z zerem konta czyta konto najstarsze, więc stan częściowy nie cofa
zachowania. Teren 15 (własność kanału przy użyciu, 11 miejsc w `core`, wykaz
`channel.list` po koncie) nie ruszył.

**Reguła Właściciela dla rund (2 września).** Agent dostaje pracę wartą co najmniej
500K tokenów wyniku i ma dać więcej wyniku, niż przeczyta; rundy dzielić na 3–5
szerokich terenów, nie na kilkanaście wąskich.

**Co otwarte — praca następnej sesji.**

1. **Kontrola rundy trzeciej i sprawdziany.** `go test -count=1 -timeout 30m` na
   `dane`, `konfig`, `zdalne`, `session`, `injection`, `transport`; jeden kontroler na
   rodzinę plików, nie na teren.
2. **Poza terenem z rundy trzeciej** (pełny wykaz w `wyniki-terenow-01-13.txt`):
   przekład `ErrKolizjaWiersza` na `conflict` w kilkunastu adapterach (`bladStudio`,
   `bladTlumaczenia`, `bladBiblioteki`, `bladDesignu`, `bladDebaty`, `bladBadan`,
   `bladAutomatyki`, `bladKomponentu`, `bladRozszerzenia`, `bladPrzegladarki`,
   `bladNadania`, `bladWskazania`); `sesje_historia.go` i `agenci.go` bez warunku
   konta; stan sesji odtwarzany z kontekstu montażu (`trwalosc_stanow.go`,
   `odtworzenie_stanu.go`, `stan_sesji_nadzor.go`) — po zawężeniu odtwarza wyłącznie
   konto najstarsze, do rozstrzygnięcia wg decyzji 34; tunele i obserwacje terminala
   na kontekście życia procesu (`adapter_modul_terminal_tunele.go`,
   `_obserwacje.go`); `przegladarka_okno.go` `sladOkna` bez warunku.
3. **Tabele bez drogi do konta**: `pamiec_tlumaczen` (bez `konto_id`, `panel_id`
   bywa NULL), `nastawa_pracy_studio` na poziomie okna, `przestrzen_badania`
   jednowierszowa na instalację (`CHECK(id = 1)`, migracja 049) — po zawężeniu konto
   młodsze dostaje kolizję; wymagają migracji ponad 485.
4. **125 korzeni pozostałych** oraz **77 zapisów do korzeni bez konta w 49 plikach**
   spoza rundy: `zapisy-korzeni-bez-konta-poza-runda-3.json`.
5. **176 więzów UNIQUE wewnątrz `CREATE TABLE`** — osobny etap z migracją.

## 4. Naprawy krytyczne — stan

| co | stan |
|---|---|
| sekret nawiązania wykluczał serwer narzędzi | naprawione; poświadczenie zastępuje sekret (decyzja 31) |
| ping zrywał bezczynne gniazdo serwera narzędzi | naprawione; gniazdo czytane biegiem odczytu, ponowienie na gnieździe zastanym |
| sumy migracji uzgadniane po numerze | naprawione; nazwa kroku sprawdzana przy każdym starcie (decyzja 32) |
| trzy produkcyjne maszyny w zaczynie każdej instalacji | naprawione krokiem 483 (decyzja 33) |
| znacznik „bramka bez poczty" jeden na instalację | naprawione; znacznik na konto, wpis zastany przejmowany raz |

Termin gniazda przed bramką, dołożony w trakcie prac, został **zdjęty**: zabijał
gniazdo serwera narzędzi i rejestrację z potwierdzeniem listem. Gniazdo
z poświadczeniem nie liczy się odtąd do granicy gniazd niezwiązanych; samo
ustalenie audytu o granicy 64 zostaje otwarte.

## 5. Ustalenia audytu powtórnego czekające na pracę

Poza granicą konta, z siedmiu zamkniętych kategorii (pełny wykaz:
`audyt2-ustalenia-7-kategorii.json`):

- **wysokie**: poświadczenia kanału w `option_env!` bez wznowienia budowy —
  instalator złożony po próbnym złożeniu wyniesie hasło
  (`instalator/src-tauri/build.rs`); pakietu serwera wdrożenia nadal nie ma,
  wykaz kieruje na plik z 18 sierpnia; `zloz.mjs` ślepy na pozycję za hasłem;
  sekret nawiązania nie jest przez powłokę ani wytwarzany, ani przekazywany
  rdzeniowi; rodzina `notification.*` bez wołającego; pokrycie klienta 5,34 %.
- **średnie**: `ZamknijKodem` daje 5 s na każde zamknięcie gniazda;
  `progress.changed` idzie do gniazd przed bramką; poświadczenie serwera
  narzędzi jedzie w wierszu poleceń procesu potomnego; wskaźnik PIN-u nie zna
  konta; `go test ./...` nie kończy się przez sprawdzian odtwarzania twarzy.

## 6. Rozstrzygnięcia czekające na Właściciela

Czynność, nie rozstrzygnięcie: zakup certyfikatu OV w Certum na dokumenty spółki
i konto SimplySign. Do zakupu wydania idą z `DANACO_PODPIS=pomijany`.

Pozycje audytu bez rozstrzygnięcia: los kanału `/wydania/` za hasłem; kształt
rodziny `project.*`; reguła wykazu narzędzi modelu; kolejność wpinania
podgrup `studio.*`; los siedmiu okien platformowych; zgoda na naprawę migracji
226, 269 i 378 nowym krokiem.

## 7. Czego nie wolno

- **Nie edytować migracji zastosowanych.** Strażnik sumy kontrolnej wywróci start
  każdej istniejącej bazy. Naprawa schematu idzie wyłącznie nowym krokiem
  o numerze wyższym niż 485.
- **Nie odtwarzać prototypów.** 27 arkuszy i 79 skryptów w `design/zasoby/` oraz
  35 prototypów okien w `design/05-okna/` są biblioteką do wpięcia. Kompozycji
  okna nie przekazuje się zleceniem.
- **Nie budować `klient-poprzedni`** w drzewie repozytorium — `emptyOutDir`
  kasuje pliki terenów roboczych.
- **Nie stawiać nowych testów, bramek i walidatorów.** Sprawdzian pisze się
  wtedy, gdy jest jedynym sposobem wykazania, że naprawa działa.
- **Nie pytać Właściciela o rzeczy rozstrzygalne ze źródeł.** Rozstrzygnięcia
  wiążące stoją w `prowadzenie/decyzje.md`.
- **Nie wystawiać niczego jako gotowego bez pokazania wyniku.** Zrzut z maszyny
  Windows z pliku pobranego z serwera, nie z dysku maszyny budującej.

## 8. Maszyny

**Maszyna budująca** — ta. Podgląd budowy wychodzi pod `http://51.75.62.180/`,
tylko port 80; Caddy kieruje na `127.0.0.1:17896`, gdzie stoi rdzeń podglądowy
z `/tmp/rdzen-podglad`, katalogiem danych `~/robocze/podglad-dane` i pakietem
interfejsu z `budowa/klient/dist`. **To proces, nie usługa** — ginie po restarcie
maszyny.

**danaco-system — 57.128.253.74.** Rdzeń wdrożenia (`danaco-console.service`,
port 17870, dane `/var/lib/danaco-console`), rdzeń wydania Studio
(`danaco-console-studio.service`, 17871), kanał pobrań `pobierz.danaco-group.pl`,
portal, serwer poczty Stalwart. Dostęp: klucz `~/.ssh/danaco_operator`, `root`.
Hasło skrzynki nadawczej stoi w `/etc/danaco-console/srodowisko` na tamtej
maszynie — nie w repozytorium.

**Uwaga wiążąca dla wdrożenia.** Baza `danaco-system` idzie inną linią numeracji:
jej kroki 406–480 noszą nazwy, których repozytorium nie zna. Po naprawie sum
kontrolnych rdzeń **odmówi na niej startu** z czytelnym komunikatem, zamiast paść
kilkadziesiąt kroków dalej na braku kolumny. Przed wdrożeniem rejestr tamtej bazy
trzeba doprowadzić do linii repozytorium; kopie sprzed ręcznych zmian leżą obok
jako `*.bak-2026-09-01-2340`.

**Maszyna Windows do sprawdzania.** Kontener `danaco-win` (dockurr/windows),
noVNC na `http://127.0.0.1:8006/`, przez `sudo docker`. Katalog wspólny
`/srv/win-danaco/wspolny` widziany jako `\\host.lan\Data`. Sterowanie Playwright
przez CDP, `executablePath: '/usr/bin/chromium-browser'`. **Pułapka:** `Alt+F4`
na pulpicie wywołuje zamknięcie systemu i kontener pada.

## 9. Jak sprawdzić, że wszystko stoi

```bash
cd ~/budowa/budowa && go build ./... && go vet ./server/... && gofmt -l server
cd ~/budowa/budowa && go test -count=1 ./server/internal/dane/... ./server/internal/store/... ./server/internal/transport/... ./server/internal/narzedzia/...
cd ~/budowa/budowa/klient && npx tsc --noEmit && npx vite build
cd ~/budowa/budowa/desktop/src-tauri && PATH=$HOME/.local/bin:$PATH cargo check --target x86_64-pc-windows-gnu
python3 ~/budowa/narzedzia/pokrycie-kontraktu.py
```

Pakiet `server/internal/core` uruchamiać wybiórczo przez `-run` albo z `-skip`
na sprawdzianach odtwarzania twarzy — sięgają silników zewnętrznych i nie kończą
się w dziesięć minut. Dziennik startu rdzenia ma mówić `komend=1086`,
`zależności zewnętrzne: 55 z 55 obecnych` i `klient=klient/dist`.

## 10. Pułapki tej maszyny

- Cel `x86_64-pc-windows-msvc` buduje się przez `cargo xwin`, a `cc-rs` szuka
  `llvm-lib` bez przyrostka wersji. W `~/.local/bin` stoją dowiązania do plików
  `*-21`. Bez `PATH=$HOME/.local/bin:$PATH` budowa pada.
- Powłoka główna składa się celem `gnu` i wymaga `WebView2Loader.dll` obok
  siebie; niesie ją instalka NSIS. Instalator kreatora składa się celem `msvc`
  i jest samodzielnym plikiem.
- Playwright chodzi wyłącznie z `executablePath: '/usr/bin/chromium-browser'`
  i `--no-sandbox`.
- Katalog tymczasowy sesji znika wraz z sesją. Co ma przetrwać, ląduje
  w `~/robocze/` albo w repozytorium.
