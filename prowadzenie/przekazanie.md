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
| skrypt audytu powtórnego | `/tmp/claude-1000/-home-ubuntu/80e91b12-.../scratchpad/audyt-powtorny.js` — katalog tymczasowy sesji znika, skrypt trzeba odtworzyć albo napisać na nowo |

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
Zapis roboczy: rewizja `955c03e0` — 94 pliki, +6490/−1423, drabina zdana.
Po niej idzie scalenie (jeden agent na całym drzewie domykający zgłoszenia „poza terenem")
i pięciu kontrolerów odbioru.

## 5. Co zostało

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

## 6. Czego nie wolno

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

## 7. Rozstrzygnięcia czekające na Właściciela

Pełny wykaz stoi w audycie, rozdział 5. Trzy blokują wydanie:

1. **Certyfikat podpisywania kodu** — OV czy EV, kto trzyma klucz, serwer znacznika czasu.
   Dziś każdy pobierający dostaje „Nieznany wydawca".
2. **Jedna maszyna i jeden port wdrożenia.** Trzy niezgodne zapisy: wykaz mówi 57.128.253.74,
   wydanie celuje w 51.75.62.180:80 (a tam odpowiada rdzeń podglądowy), jednostka systemd
   stawia 17870.
3. **Nazwa i certyfikat TLS dla rdzenia wdrożenia** — bez nazwy `wss://` nie ma jak postawić,
   gołe IP nie zestawia TLS.

Dalej m.in.: kształt rodziny `project.*`, reguła wykazu narzędzi modelu (603 komendy poza
wykazem bez podanej racji), wydanie ARM64, domyślna wartość bramki na pętli zwrotnej,
czym szyfrować sejf poświadczeń, los siedmiu okien platformowych.

## 8. Jak sprawdzić, że wszystko stoi

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

## 9. Pułapki tej maszyny

- Cel `x86_64-pc-windows-msvc` buduje się przez `cargo xwin`, a `cc-rs` szuka `llvm-lib`
  bez przyrostka wersji. W `~/.local/bin` stoją dowiązania `llvm-lib`, `llvm-rc`, `llvm-ar`,
  `llvm-dlltool`, `llvm-objcopy` do plików `*-21`. Bez `PATH=$HOME/.local/bin:$PATH` budowa pada.
- Powłoka główna składa się celem `gnu` i wymaga `WebView2Loader.dll` obok siebie; niesie ją
  instalka NSIS. Instalator kreatora składa się celem `msvc` i jest samodzielnym plikiem.
- Playwright chodzi wyłącznie z `executablePath: '/usr/bin/chromium-browser'` i `--no-sandbox`.
- Katalog tymczasowy sesji (`/tmp/claude-1000/...`) znika wraz z sesją. Co ma przetrwać,
  ląduje w `~/robocze/` albo w repozytorium.
