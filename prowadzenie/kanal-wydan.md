# Kanał wydań

Opis kanału, którym Danaco Console wychodzi do pobierającego: maszyny, ścieżki,
poświadczenia i porządek składania wydania.

**Dlaczego nie w `budowa/witryna/wydania.json`.** Wykaz wydań jest publiczny —
leży na stronie i jest wkompilowany w kreator instalacji oraz w powłokę, więc
czyta się go poleceniem `strings` na pobranym pliku. Nazwa konta, ścieżka do
konfiguracji serwera stron i adresy maszyn nie mają tam czego szukać. Wykaz
niesie wyłącznie to, czym posługują się strona i kreator: adres kanału, adresy
plików, sumy, rozmiary i stan każdej pozycji.

## Maszyny

| Co | Gdzie | Czym postawione |
|---|---|---|
| kanał pobrań `pobierz.danaco-group.pl` | 57.128.253.74 (IPv6 `2001:41d0:601:1100::2485`) | Caddy; TLS na 443 z certyfikatu ACME, przekierowanie z 80 robi Caddy sam |
| rdzeń wdrożenia `console.danaco-group.pl` | 57.128.253.74 (ta sama maszyna, `danaco-system`) | rdzeń na `127.0.0.1:17870` z jednostki systemd pakietu; na świat `:443` przez Caddy, TLS z ACME (rozstrzygnięcie 29) |
| katalog wydawany kanału | `/opt/danaco/danaco-console/portal` | — |
| pliki wydań | `/opt/danaco/danaco-console/wydania/` | — |
| strona wizerunkowa `danaco-group.pl` | 137.74.41.149 | nginx |

Nazwy `danaco-system` i `danaco-web` nie są nazwami DNS — to etykiety wewnętrzne
z migracji punktów dostępu i w DNS nie rozwiązują się.

## Dwie ścieżki kanału

| Ścieżka | Kto po nią sięga | Zamknięcie |
|---|---|---|
| `/pliki/<wersja-data>/` | kreator instalacji w kroku 5 oraz strona „Pobierz” | otwarta — kreator sześciu kroków nie ma pola na hasło |
| `/wydania/<wersja-data>/` | Właściciel i administrator | uwierzytelnienie podstawowe, użytkownik `wlasciciel`, hasło jako skrót bcrypt w `/etc/caddy/Caddyfile` |

Pozycja wykazu leżąca pod `/wydania/` musi nieść `chronione_haslem: true` —
strona zapowiada wtedy pytanie o hasło przy tej jednej pozycji, a krok odbioru
w `zloz.mjs` przyjmuje dla niej odpowiedź 401. Pozycja spod `/pliki/` niesie
`chronione_haslem: false` i musi odpowiadać kodem 200.

## Sekcja `kanal` wykazu

Wykaz niesie w `kanal` wyłącznie cztery pola czytane maszynowo; opis kanału
stoi tu, nie w wykazie. Kreator instalacji czyta z tej sekcji samo `adres`
(`instalator/src-tauri/src/pobranie.rs`) i pod nim szuka bieżącego
`wydania.json`, więc pole musi zostać pod tą nazwą.

| Pole | Kto czyta | Znaczenie |
|---|---|---|
| `adres` | kreator, `zloz.mjs`, strona „Pobierz” | adres kanału; względem niego rozwijane są adresy w `plik` |
| `wdrozony` | `zloz.mjs`, strona „Pobierz” | `true` — pliki leżą pod adresami i strona wystawia przycisk; krok odbioru w `zloz.mjs` pyta każdą pozycję żądaniem HEAD i odmawia złożenia przy innej odpowiedzi. `false` — strona pokazuje adres jako tekst, bez przycisku |
| `tymczasowy` | strona „Pobierz” | `true` dokłada nad kartami zdanie, że kanał stoi na maszynie budującej pod certyfikatem własnym i adres nie jest docelowy |
| `czego_brakuje` | strona „Pobierz” | zdanie o brakującym elemencie kanału, drukowane przy `wdrozony: false` |

Certyfikat kanału wystawia i odnawia Caddy (ACME, Let's Encrypt); przeglądarka
nie pokazuje ostrzeżenia, a `https://pobierz.danaco-group.pl/` odpowiada 200.
Do zrobienia zostaje strona wizerunkowa produktu w portfolio grupy — bez
pobierania plików, z jednym odnośnikiem do kanału.

## Rdzeń wdrożenia — rozstrzygnięcie 29

Pola `serwer.rdzen_wdrozenia.adres`, `.port` i `.schemat` w wykazie niosą
wartość z rozstrzygnięcia 29 (`prowadzenie/decyzje.md`):

| co | wartość |
|---|---|
| maszyna wdrożenia | `danaco-system`, 57.128.253.74 — ta sama, na której stoi kanał pobrań |
| nasłuch rdzenia | `127.0.0.1:17870` — jednostka systemd pakietu (`DANACO_ADRES=127.0.0.1`, `DANACO_PORT=17870`, `DANACO_WYMOG_LOGOWANIA=true`) |
| nazwa i port na świat | `console.danaco-group.pl:443` przez Caddy, TLS z ACME; TLS kończy się na Caddy, rdzeń nie potrzebuje `DANACO_TLS_*` |
| wskazanie w wydaniu | `DANACO_HOST_WDROZENIA=console.danaco-group.pl`, `DANACO_PORT_WDROZENIA=443`, `DANACO_SCHEMAT_WDROZENIA=https` |
| podgląd budowy | 51.75.62.180:80 — wyłącznie podgląd, nigdy cel wydania |

Wartość obowiązuje naraz w trzech miejscach: jednostce systemd pakietu
(`packaging/drzewo/lib/systemd/system/danaco-console.service`), opisie
w `DEBIAN/control` (zgodność portu pilnuje `scripts/pakiet-serwera.sh`) oraz
trzech zmiennych podawanych przy składaniu instalki powłoki. Skrypty składania
nadal odmawiają budowy bez tych zmiennych — wartość stoi w wykazie i w tym
opisie, a nie w skrypcie, żeby instalka złożona z inną wartością nie wyglądała
jak wydanie. Wydana powłoka 2.0.0 z 30 sierpnia celuje jeszcze
w 51.75.62.180:80 i wymaga ponownego złożenia z wartościami wyżej.

## Porządek składania wydania

1. `npm --prefix budowa/klient run budowanie` — powłoka bierze zbudowany interfejs.
2. `DANACO_HOST_WDROZENIA=console.danaco-group.pl DANACO_PORT_WDROZENIA=443 DANACO_SCHEMAT_WDROZENIA=https budowa/scripts/instalka-hybryda-win-x64.sh`
   — powłoka programu (`-arm.sh` dla ARM64 — poza tym etapem, rozstrzygnięcie 28).
3. `budowa/scripts/instalka-kreatora-win-x64.sh` — kreator instalacji; bierze
   adres i sumę powłoki z wykazu wydań, więc składa się PO kroku 2 i po pomiarze.
4. `budowa/scripts/pakiet-serwera.sh` — pakiet serwera wdrożenia.
5. Wgranie plików pod `/pliki/<wersja-data>/` na maszynie kanału.
6. `node budowa/witryna/zmierz.mjs` — wpisuje rozmiar i sumę Z PLIKU na dysku.
7. `node budowa/witryna/zmierz.mjs --sprawdz` — musi powiedzieć „rozjazdów 0,
   plików brakuje 0”; inaczej wykazu nie wolno publikować.
8. `node budowa/witryna/zloz.mjs` — krok odbioru pyta kanał żądaniem HEAD
   o każdą pozycję i odmawia złożenia strony, gdy kod jest inny niż umówiony.

Wszystkie skrypty odkładają wynik do `budowa/wydania/<wersja>-<data>/`, czyli
tam, gdzie `katalog_plikow` wykazu każe szukać plików do pomiaru; data to dzień
składania. Inny katalog wskazuje się z zewnątrz zmienną `DANACO_KATALOG_WYDANIA`
(przyjmują ją `instalka-hybryda-win-{x64,arm}.sh`, `instalka-kreatora-win-x64.sh`
i `pakiet-serwera.sh`) — złożenie próbne nie ma odkładać pliku między wydania.

## Podpis Authenticode

Rozstrzygnięcie 27: certyfikat OV od Certum, klucz w usłudze SimplySign,
podpis przez klienta SimplySign (PKCS#11, `osslsigncode`), znacznik czasu
`http://time.certum.pl`. Zakup wymaga dokumentów spółki i jest jedyną czynnością
po stronie Właściciela. Do zakupu każdy publikowany plik Windows ma pusty
katalog Security i niesie w wykazie `podpisany: false`; strona „Pobierz” pisze
o tym przy każdej pozycji Windows.

Skrypty składania mierzą katalog Security gotowego pliku i rozstrzygają zmienną
`DANACO_PODPIS`:

- `pomijany` (przyjmowane bez podania; do zakupu certyfikatu podawane jawnie) —
  brak podpisu jest wypisany w pomiarze wyniku i wydanie wychodzi;
- `wymagany` — plik z pustym katalogiem Security nie zostaje odłożony do wydania.

Po zakupie certyfikatu: polecenie podpisujące wchodzi do `bundle.windows` obu
profili powłoki i profilu kreatora, składanie idzie z `DANACO_PODPIS=wymagany`,
a pozycje wykazu dostają `podpisany: true`.

## Czego dziś brakuje

- Pliku kreatora i pakietu serwera nie ma w `budowa/wydania/2.0.0-2026-08-30/`,
  więc `zmierz.mjs --sprawdz` melduje „plików brakuje 2” i odmawia. Sumy tych
  dwóch pozycji są przepisane, nie zmierzone z pliku leżącego w drzewie.
- Wystawiony plik kreatora nie pochodzi ze skryptu: pod nazwą instalatora leży
  surowa binarka, nie pakiet NSIS. `instalka-kreatora-win-x64.sh` składa pakiet
  właściwy — wymiana pliku w kanale wymaga wgrania.
- Pakiet serwera w kanale jest wydaniem 1.0.0 z 18 sierpnia i leży pod ścieżką
  za hasłem; pakiet 2.0.0 trzeba złożyć i wgrać pod `/pliki/`.
- Powłoki ARM64 pod numerem 2.0.0 nie ma i w tym etapie nie powstaje
  (rozstrzygnięcie 28) — pozycja stoi w `w_przygotowaniu` bez adresu pliku.
  Kreator czyta wykaz najpierw z kanału (`pobranie.rs`, `wykaz_biezacy`), więc
  wykaz WYSTAWIONY dziś na kanale, w którym pozycja ARM64 wciąż leży
  w `wydania` i wskazuje plik 1.0.0 spod `/wydania/`, daje na maszynie ARM
  odmowę `odpowiedz-serwera` (401 — kreator nie ma pola na hasło). Po
  opublikowaniu świeżego wykazu z tego drzewa (`zloz.mjs`) pozycja ARM64 znika
  z `wydania` i odmowa staje się `wydanie-nieopublikowane`, nazwana z wykazu,
  nie z odpowiedzi serwera. Pozycji nie wolno wskazać na plik 1.0.0.
