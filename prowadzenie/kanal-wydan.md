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

## Rdzeń wdrożenia — czeka na rozstrzygnięcie Właściciela

Pola `serwer.rdzen_wdrozenia.adres` i `.port` w wykazie stoją puste, bo trzy
zapisy w drzewie mówią dziś trzy różne rzeczy:

- jednostka systemd pakietu serwera stawia rdzeń na porcie **17870**, na
  wszystkich interfejsach;
- opis w `DEBIAN/control` zapowiada ten sam port 17870 (zgodność pilnuje zapora
  w `scripts/pakiet-serwera.sh`);
- wydana powłoka 2.0.0 celuje w **51.75.62.180:80**, a pod tym adresem port 80
  prowadzi do procesu podglądu maszyny budującej, nie do rdzenia wdrożenia.

Rozstrzygnięcie ma nazwać jedną maszynę i jeden port — osobno rozstrzyga się, czy
powłoka sięga rdzenia wprost, czy przez serwer stron. Obowiązuje naraz w trzech
miejscach: jednostce systemd pakietu, opisie w `DEBIAN/control` oraz zmiennych
`DANACO_HOST_WDROZENIA` i `DANACO_PORT_WDROZENIA` podawanych przy składaniu
instalki powłoki. Do tego czasu skrypty składania odmawiają budowy bez tych
zmiennych — wartości nie ma w drzewie i żaden skrypt jej nie dopowiada.

## Porządek składania wydania

1. `npm --prefix budowa/klient run budowanie` — powłoka bierze zbudowany interfejs.
2. `DANACO_HOST_WDROZENIA=… DANACO_PORT_WDROZENIA=… budowa/scripts/instalka-hybryda-win-x64.sh`
   (odpowiednio `-arm.sh` dla ARM64) — powłoka programu.
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
tam, gdzie `katalog_plikow` wykazu każe szukać plików do pomiaru. Katalog
wskazuje się z zewnątrz zmienną `DANACO_KATALOG_WYDANIA` — złożenie próbne nie
ma odkładać pliku między wydania.

## Podpis Authenticode

Producent nie ma dziś certyfikatu podpisywania kodu, więc każdy publikowany plik
Windows ma pusty katalog Security i niesie w wykazie `podpisany: false`; strona
„Pobierz” pisze o tym przy każdej pozycji Windows.

Skrypty składania mierzą katalog Security gotowego pliku i rozstrzygają zmienną
`DANACO_PODPIS`:

- `pomijany` (przyjmowane bez podania) — brak podpisu jest wypisany w pomiarze
  wyniku i wydanie wychodzi;
- `wymagany` — plik z pustym katalogiem Security nie zostaje odłożony do wydania.

Po zdobyciu certyfikatu: polecenie podpisujące wchodzi do `bundle.windows` obu
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
- Powłoki ARM64 pod numerem 2.0.0 nie ma — pozycja stoi w `w_przygotowaniu`
  bez adresu pliku, a kreator na maszynie ARM odmawia nazwanym powodem
  (`wydanie-nieopublikowane`). Czy i kiedy ARM64 wychodzi pod 2.0.0,
  rozstrzyga Właściciel; do tego czasu pozycji nie wolno wskazać na plik 1.0.0
  spod `/wydania/`, bo kreator nie ma pola na hasło i pobranie kończy się
  odmową 401.
