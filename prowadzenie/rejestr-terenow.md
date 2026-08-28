# Rejestr terenów

Żywy wykaz terenów. Prowadzi go wyłącznie Prowadzący budowę, na gałęzi `main`.
Teren bez wpisu w tym rejestrze nie jest otwarty, a praca na nim nie zostanie
przyjęta. Zasady podziału opisuje [ustrój budowy](ustroj-budowy.md).

## Tereny otwarte

### rama-wpieta-w-biblioteke

| | |
|---|---|
| **Galaz** | `teren/rama-wpieta` z `main` |
| **Drzewo** | `~/robocze/rama-wpieta` |
| **Wykaz plikow** | `budowa/klient/src/rama/` wraz z podkatalogami; `budowa/klient/src/aplikacja.ts`; `budowa/klient/index.html`; `budowa/klient/arkusze.css`; `budowa/klient/vite.config.*` jesli zajdzie potrzeba; `budowa/klient/package.json` |
| **Poza terenem** | `src/wejscie/`, `polaczenie/`, `protokol/`, `moduly/`; rdzen; kontrakt; **`design/` — CZYTASZ, NIE ZMIENIASZ**; `prowadzenie/` |

**Rzecz, ktora ten teren naprawia.** Rama zbudowana wczesniej NIE JEST wpieta w warstwe
projektowa — jest jej imitacja napisana od zera. Zmierzone przez prowadzenie:

| | Prototyp centrum dowodzenia | Klient dzis |
|---|---|---|
| arkuszy stylow | **20** | 8 |
| skryptow zachowania | **18** | **0** |

Wlasciciel zbudowal biblioteke wlasnie po to, zeby ja podpiac: 27 arkuszy i 79 skryptow
w `design/zasoby/`. Skladniki to funkcje domkniete dzialajace na ZNACZNIKU przez klasy
i atrybuty — `pasek-stanu.js` wystawia `window.dnPasekStanu` i czyta
`.dn-stan-tor[data-udzial]`. Kompozycja siedzi w znaczniku prototypu, zachowanie
w skryptach, wyglad w arkuszach. Klient ma wytwarzac TEN SAM ZNACZNIK i wciagac
TE SAME arkusze i skrypty, a nie pisac wlasnych klas.

**Wina lezy po stronie prowadzenia.** Poprzednie zlecenie mowilo „trzy strefy" zamiast
„wepnij biblioteke", a kryteria odbioru pytaly o wezly i lancuchy zamiast o zgodnosc
obrazu z prototypem. Wykonawca zrobil dokladnie to, o co go poproszono.

**Zrodlo prawdy: `design/05-okna/przeplyw/centrum-dowodzenia.html`.** Znacznik tego pliku
jest wzorcem struktury. Bierzesz z niego uklad, nazwy klas i atrybuty, a wypelniasz
DANYMI Z RDZENIA. Tresc przykladowa prototypu — nazwy sesji, kafle, liczby — NIE WCHODZI
do klienta; wchodzi tam, skad rdzen odda odpowiedz, albo stan pusty.

**Strefy, ktore maja stanac**, wszystkie sa w prototypie:
pasek okna ze sterowaniem oknem, lewa szyna ikon z OSOBNA IKONA NA MODUL, panel Sesje
i Projekty, pasmo kart, pasek narzedzi z wyszukiwaniem, obszar roboczy, pasek stanu
z tozsamoscia Operatora i licznikami.

**Kryteria odbioru — jedyne, ktore maja sens.**
1. Klient wciaga WSZYSTKIE 20 arkuszy i WSZYSTKIE 18 skryptow, ktorych uzywa prototyp.
   Policz je programem i przytocz liczbe.
2. ZRZUT EKRANU przy 2560x1440 zestawiony z prototypem. Robisz go SAM, ogladasz SAM
   i opisujesz roznice strefa po strefie. Zrzut jest dowodem; deklaracja nim nie jest.
3. Zadna ikona modulu nie powtarza sie miedzy modulami.
4. Zero tresci przykladowej z prototypu. Kazda wartosc pochodzi z rdzenia albo jest
   nazwanym stanem pustym.
5. `npm run typy` i `npm run budowanie` przechodza, a `dist/` niesie arkusze i skrypty
   biblioteki — sprawdz, ze wytworzona strona je naprawde laduje, nie tylko `index.html`.

**Pulapka do rozstrzygniecia:** skrypty biblioteki sa funkcjami domknietymi, nie modulami
— nie maja `export`. Vite ich nie zapakuje jak modulow. Rozstrzygnij, jak je wydac
(kopia do zasobow statycznych, wpis w `index.html`, alias w nastawach Vite) i NAZWIJ
wybor. Warstwy projektowej NIE ZMIENIASZ — jest wlasnoscia Wlasciciela.


### okno-studia-w-kliencie

| | |
|---|---|
| **Galaz** | `teren/okno-studia` z `teren/rama-aplikacji` — nie z `main`, bo okno wchodzi W RAME, a rama nie jest jeszcze scalona |
| **Drzewo** | `~/robocze/okno-studia` |
| **Wykaz plikow** | `budowa/klient/src/moduly/studio/` (nowy katalog) wraz z tresciami; `budowa/klient/src/rama/montaz.ts` i `skladniki/szyna.ts` wylacznie w zakresie otwierania okna modulu; `budowa/klient/arkusze.css` wylacznie w zakresie wpiecia arkusza okna; `budowa/klient/package.json` |
| **Poza terenem** | `src/wejscie/`, `polaczenie/`, `protokol/`; rdzen; kontrakt; `design/`; `prowadzenie/` |

**Przedmiot — brakujace ogniwo ciagu.** Wlasciciel zada pelnego ciagu od instalacji
do okna Studio, dzialajacego od konca do konca. Zmierzony stan ogniw:

| Ogniwo | Stan |
|---|---|
| instalator Windows | wytworzony, pakiet NSIS stoi |
| droga wejscia w kliencie | 8 plikow, dziala |
| rama aplikacji | 9 plikow, na galezi `teren/rama-aplikacji` |
| **okno Studio w kliencie** | **NIE ISTNIEJE — `src/moduly/` ma zero plikow** |

Do tego szyna nawigacji ramy dzis **nie otwiera okna modulu** — klikniecie zmienia
wylacznie tytul belki. Poprzedni wykonawca nazwal to wprost i slusznie zostawil,
bo mial to poza terenem. Teraz to jest przedmiot.

**Skad bierzesz ksztalt okna — i czego NIE ROBISZ.** Kompozycja okna Studio nalezy do
Wlasciciela i stoi w przyjetym prototypie `design/05-okna/moduly/studio.html`.
Czytasz go jako WZORZEC ukladu i nazw klas. **Nie proponujesz wlasnego ukladu, nie
dokladasz stref, nie zmieniasz kompozycji.** Gdy prototyp milczy o czyms, czego kod
potrzebuje, wybierasz rozwiazanie najblizsze temu, co juz stoi w `src/wejscie/`
i `src/rama/`, i nazywasz to w rozstrzygnieciach.

**Kryteria odbioru — wszystkie sprawdzalne uruchomieniem okna, nie sprawdzianem.**
1. Klikniecie pozycji Studio w szynie nawigacji OTWIERA okno Studia w obszarze roboczym
   ramy. Powrot do innego modulu zdejmuje je.
2. Okno pokazuje DANE Z RDZENIA, nie tresc wymyslona. Rdzen niesie komendy `studio.*` —
   przekroj od komendy do zapisu w SQLite zostal sprawdzony dzialaniem. Wybierasz komende,
   ktora daje sie pokazac na wejsciu (wykaz materialu albo wykaz urzadzen), wolasz ja
   przez warstwe `protokol/` tak, jak robi to droga wejscia, i pokazujesz odpowiedz.
3. Odmowa rdzenia jest POKAZANA Operatorowi, nie polkniete. Puste dane to nie to samo
   co odmowa i okno ma je rozrozniac.
4. Zero lancuchow widocznych dla Operatora poza katalogiem tresci — ten klient juz tego
   pilnuje i katalog `rama/` pokazuje wzorzec.
5. Barwy, odstepy i pismo wylacznie z zetonow warstwy projektowej.
6. `npm run typy` i `npm run budowanie` przechodza.


### wpiecie-analizy-typescriptu

| | |
|---|---|
| **Galaz** | `teren/wpiecie-analizy-typescriptu` z `main` |
| **Drzewo** | `~/robocze/wpiecie-analizy-typescriptu` |
| **Wykaz plikow** | `budowa/server/internal/core/adapter_modul_developer_jezyk.go`; `budowa/server/internal/core/adapter_modul_developer_narzedzia.go`; ich sprawdziany |
| **Poza terenem** | kontrakt i wytwory, migracje, klient, `design/`, `prowadzenie/` |

**Przedmiot — zdolnosc produktu, ktorej nie ma.** Analiza statyczna TypeScriptu
i JavaScriptu NIE JEST WPIETA, choc arsenal ja obiecuje. Wykaz zaleznosci wydawany przez
binarium rdzenia mowi o `eslint` wprost: „analiza statyczna kodu TypeScript i JavaScript".
Program stoi na maszynie. Pracy nie dostaje.

Zmierzone przez kontrole i potwierdzone przez prowadzenie:
- `narzedzieEslint` wystepuje w rdzeniu WYLACZNIE w wykazie warsztatu
  (`adapter_modul_developer_narzedzia.go:73`) i w wykazie zaleznosci. Jedyne uruchomienie
  to `eslint --version` z `developer.toolchain.check` — oslona na PATH pokazala dokladnie
  to i nic wiecej. Po takim kryterium „wpiety" bylby kazdy program na maszynie.
- `analizyRepozytorium` (`adapter_modul_developer_jezyk.go:499`) zna wylacznie
  `golangci-lint`, `ruff`, `stylelint` i `typos`. Pliki `.ts` i `.js` nie trafiaja nigdzie.
- `narzedzieStaticcheck` jest w tym samym polozeniu: tylko wykaz warsztatu, zero pracy.
- `serwerJezykaPliku` (`:471-477`) oddaje dla `.ts/.js` `narzedzieSerweraTypeScript`
  z druga wartoscia `false`, a komentarz `adapter_modul_developer_narzedzia.go:58` mowi
  wprost: „rdzen nie ma dzis czym go zapytac". Nawigacja po symbolach i przeksztalcenia
  w plikach TypeScriptu sa niedostepne.

**Skutek dla Operatora.** Modul Developer sprawdza kod Go i Pythona, a kodu, w ktorym
napisany jest wlasny klient produktu, nie sprawdza wcale.

**Kryteria odbioru.**
1. `developer.lint.get` na pliku `.ts` albo `.js` kieruje prace do `eslint` i oddaje
   jego uwagi. Dowodem jest URUCHOMIENIE przez montaz produkcyjny — oslona na PATH
   zapisujaca wywolanie albo odczyt uwag zwroconych przez program.
2. `developer.scan.run` obejmuje pliki `.ts` i `.js`, jesli tam jest co sprawdzic.
3. `staticcheck` dostaje prace albo ZNIKA z wykazu warsztatu. Program w wykazie, ktory
   nigdy nie dostaje pracy, klamie Operatorowi o zdolnosci produktu.
4. `serwerJezykaPliku` dla `.ts/.js`: albo rdzen ma czym go zapytac, albo odmowa jest
   NAZWANA i mowi Operatorowi, czego brakuje — nie cicha druga wartosc `false`.
   Rozstrzygniecie nalezy do wykonawcy i ma byc nazwane.
5. Sprawdzian mierzy SKUTEK i PADA przed naprawa. Sprawdzian na samo wywolanie programu
   jest brakiem sprawdzianu.
6. `go build ./...` oraz `go test -run 'Developer|Jezyk|Lint|Scan|Warsztat' ./internal/core/`.


### narzedzia-odblokowujace-arsenal

| | |
|---|---|
| **Galaz** | `teren/narzedzia-odblokowujace-arsenal` z `main` |
| **Drzewo** | `~/robocze/narzedzia-odblokowujace-arsenal` |
| **Wykaz plikow** | `budowa/shared/contract.json`; wytwory `budowa/shared/contract.go` i `budowa/shared/contract.ts` WYLACZNIE z generatora |
| **Poza terenem** | rdzen, klient, migracje, `design/`, `prowadzenie/` |

**Podstawa: rozstrzygniecie 20, etap drugi — z kolejnoscia zmieniona pomiarem.**
Etap drugi mial isc po wielkosci braku (Design 96 komend, Research 71, Studio 63).
Pomiar drogi od komendy do programu zmienil te kolejnosc: obszar z 96 niewystawionymi
komendami, ktory nie odblokowuje ZADNEGO programu, jest wart mniej niz jedna komenda
otwierajaca modelowi rozpoznanie pisma.

**Stan zmierzony 28.08.2026.** Arsenal niesie 55 programow; 50 jest osiagalnych z jakiejs
komendy. Komenda, ktora MODEL moze wywolac, osiagalne sa **23**. Pozostale **27 lezy za
komendami, ktorych nie ma w wykazie narzedzi modelu**. Model nie rozpozna dzis pisma,
nie uzyje przegladarki, nie sprawdzi ani nie sformatuje kodu, nie zbada dostepnosci
ani wydajnosci, nie sprawdzi pisowni i nie siegnie po skaner.

**Zadanie.** Dopisac do `narzedzia.pozycje` w `contract.json` DOKLADNIE ponizsze 23 pozycje,
ani jednej wiecej. Kazda odblokowuje wymienione przy niej programy:

1. `apps.performance.audit` — odblokowuje: lighthouse
2. `browser.accessibility.audit` — odblokowuje: pa11y
3. `browser.console.read` — odblokowuje: chromium-browser
4. `browser.device.emulate` — odblokowuje: chromium-browser
5. `browser.dom.inspect` — odblokowuje: chromium-browser
6. `browser.network.har` — odblokowuje: chromium-browser
7. `browser.screenshot.capture` — odblokowuje: chromium-browser
8. `browser.scroll` — odblokowuje: chromium-browser
9. `developer.api.load.run` — odblokowuje: autocannon
10. `developer.debug.session.start` — odblokowuje: dlv
11. `developer.format.run` — odblokowuje: gofmt, goimports, prettier
12. `developer.grep.replace` — odblokowuje: ast-grep
13. `developer.grep.search` — odblokowuje: ast-grep
14. `developer.lint.get` — odblokowuje: golangci-lint, ruff, stylelint
15. `developer.refactor.apply` — odblokowuje: gopls
16. `developer.scan.run` — odblokowuje: dupl, jscpd, semgrep, ruff
17. `developer.symbol.navigate` — odblokowuje: gopls
18. `developer.toolchain.check` — odblokowuje: gofmt, typos
19. `library.metadata.get` — odblokowuje: exiftool
20. `studio.ingest.device.list` — odblokowuje: scanimage
21. `studio.ingest.recognize` — odblokowuje: tesseract, unpaper
22. `terminal.script.lint` — odblokowuje: shellcheck, shfmt, pwsh, python3, ruff
23. `translate.proofread.run` — odblokowuje: hunspell, vale

Razem odblokowuje **27 programow** — wszystkie, ktore dzis stoja poza zasiegiem modelu.

**Jak pisac zdanie `zastosowanie` — nauka z rundy poprzedniej.** Kontrola obalila
32 zdania poprzedniego terenu na szesciu rzeczach i te same bledy zwroca i te prace:
- Zdanie ma mowic KIEDY siegnac, nie CO komenda robi. Model widzi opis komendy I Twoje
  zdanie w JEDNYM napisie — zdanie powtarzajace opis daje mu to samo dwa razy, a pole
  przeznaczone na wyzwalacz zostaje puste co do tresci. Zobacz sam:
  `grep -A2 'usage.summary.get' budowa/shared/contract.go`.
- Zdanie NIE powiela ksztaltu zadania: zadnych wartosci pol ani wyliczen, bo stoja obok.
- Zdanie NIE ZAWEZA skutku wbrew opisowi komendy. W rundzie poprzedniej zdanie mowilo
  „cala historie NIEPRZYPIETA", a komenda kasuje takze przypiete — model uwierzylby,
  ze wpisy przypiete sa bezpieczne.
- Wyzwalacz ma byc SPRAWDZALNY PRZEZ MODEL. Warunek, ktorego model nie widzi, nie jest
  wyzwalaczem.
- Zdanie nie gubi OSTRZEZENIA stojacego w opisie komendy.
- Pary bliskie sobie maja byc rozroznialne: `developer.lint.get` wobec `developer.scan.run`,
  `developer.grep.search` wobec `developer.grep.replace`, `browser.dom.inspect` wobec
  `browser.console.read`, `studio.ingest.recognize` wobec `studio.ingest.device.list`.

**Kryteria odbioru.**
1. Dopisane DOKLADNIE 23 pozycje z wykazu, ani jednej innej. Roznica zbiorow policzona
   programem, nie okiem.
2. Zadna nie wskazuje komendy nieistniejacej, nie powtarza komendy juz wystawionej,
   nie ma pustego `zastosowanie`.
3. Zdanie kazdej pozycji przechodzi szesc prob z akapitu wyzej. Kazda sprawdzona osobno.
4. `contract.go` i `contract.ts` powstaly GENERATOREM, nie recznie.
5. `bash narzedzia/drabina.sh szybka` przechodzi, w tym szczebel swiezosci wytworow.
6. `go test -run Kontrakt ./internal/core/` przechodzi.
7. Zadna inna sekcja `contract.json` nie ruszona.


### narzedzia-obszarow-zerowych

| | |
|---|---|
| **Galaz** | `teren/narzedzia-obszarow-zerowych` z `main` |
| **Drzewo** | `~/robocze/narzedzia-obszarow-zerowych` |
| **Wykaz plikow** | `budowa/shared/contract.json`; wytwory `budowa/shared/contract.go` i `budowa/shared/contract.ts` WYLACZNIE z generatora |
| **Poza terenem** | rdzen, klient, migracje, `design/`, `prowadzenie/` |

**Podstawa.** Rozstrzygniecie 20 rejestru decyzji. Kontrakt zapisuje kryterium
wystawiania komend kanalowi modelu; implementacja od niego odstaje o 716 komend,
a 13 obszarow nie ma ANI JEDNEGO narzedzia — model nie siegnie tam wcale.
Ten teren zamyka etap pierwszy: obszary zerowe.

**Zadanie.** Dopisac do `narzedzia.pozycje` w `contract.json` DOKLADNIE ponizsze
32 pozycje, ani jednej wiecej:

1. `mobile.status.get`
2. `mobile.process.list`
3. `mobile.process.control`
4. `device.list`
5. `retention.set`
6. `clipboard.list`
7. `clipboard.push`
8. `clipboard.pin`
9. `clipboard.delete`
10. `snippet.list`
11. `snippet.set`
12. `snippet.delete`
13. `launcher.hotkey.get`
14. `launcher.hotkey.set`
15. `provenance.call.list`
16. `provenance.call.get`
17. `provenance.call.replay`
18. `provenance.call.rate`
19. `provenance.trace.export`
20. `usage.summary.get`
21. `usage.report.build`
22. `alert.rule.save`
23. `alert.rule.list`
24. `alert.rule.remove`
25. `alert.trigger.list`
26. `alert.trigger.acknowledge`
27. `health.probe.save`
28. `health.probe.list`
29. `health.probe.remove`
30. `health.probe.run`
31. `health.result.list`
32. `health.uptime.get`

Kazda pozycja niesie dwa pola: `komenda` i `zastosowanie`. Zdanie `zastosowanie`
mowi modelowi, KIEDY siegnac po narzedzie, i powstaje z opisu komendy stojacego
w kontrakcie — nie z domyslu. Deklaracja NIE powiela ksztaltu zadania.

**Czego wystawic NIE WOLNO** — te komendy sa objete wyjatkiem zapisanym
w `narzedzia.opis` i ich pominiecie jest zgodne z zasada, nie brakiem:
`connection.hello`, wszystkie dziewiec `auth.*`, `device.revoke` (uniewaznienie
dostepu urzadzenia to zapis do punktu dostepu), `model.channel.set` (podmiana
kanalu obslugujacego okno to zmiana wlasnego dostepu modelu).
Wystawienie ktorejkolwiek z nich jest zwrotem terenu.

**RUNDA DRUGA — zdania zwrocone przez kontrole tresci 28.08.** Kontrola zgodnosci
przyjela w calosci: 372 pozycje, zero wiszacych, zero powtorzen, zero pustych, zadna
komenda objeta wyjatkiem nie wystawiona, wytwory bajt w bajt z generatora, 128 dopisan
i ZERO usuniec w kontrakcie. Liczby sa dobre. Wracaja ZDANIA.

| | Rzecz | Ile pozycji |
|---|---|---|
| A | Zdanie mowi, CO komenda robi, zamiast KIEDY po nia siegnac — opis komendy przepisany po slowie „Uzyj". W wytworzonym opisie modelu stoi wtedy dwa razy to samo w jednym napisie, np. `usage.summary.get`: „Zwraca zuzycie tokenow... Uzyj, aby poznac zuzycie tokenow...". `narzedzia.opis` zada czegos odwrotnego: zdania mowiacego, kiedy siegnac. | 15 z 32 |
| B | Powielenie ksztaltu zadania, zakazane wprost w `narzedzia.opis`: wartosci pol i wyliczenia wypisane w zdaniu, choc stoja obok w tej samej pozycji — `retention.set`, `mobile.process.list`, `provenance.call.list`, `usage.summary.get`, `clipboard.list`. | 5 |
| C | **GROZNE** — `clipboard.delete`: zdanie mowi „cala historie NIEPRZYPIETA", a opis komendy mowi „albo cala historie", przy czym pole `includePinned` pozwala skasowac takze przypiete. Model uwierzy, ze wpisy przypiete sa bezpieczne, i skasuje je nieodwracalnie. | 1 |
| D | Wyzwalacz niesprawdzalny przez model — `mobile.status.get` i `mobile.process.control` opisuja czynnosc urzadzenia („gdy urzadzenie pyta"), czyli warunek, ktorego model nie widzi. Zdania powstaly z opisu pola `deviceId`, nie z opisu komendy. | 2 |
| E | Zgubione ostrzezenia stojace w opisie komendy: `launcher.hotkey.set` gubi „skrot zajety przez inny program nie jest bledem zapisu"; `snippet.list` gubi „slownik nalezy do rdzenia, nie do jednego okna"; `provenance.call.replay` gubi powtorzenie na innym kanale i modelu; `provenance.call.get` podaje tresc promptu jako pewnik wbrew polu `includeContent`. | 4 |
| F | Nierozroznialnosc par bliskich sobie — `health.probe.list` wobec `health.result.list`, `alert.rule.list` wobec `alert.trigger.list`, `clipboard.list` wobec `snippet.list`. | — |

**Luka bramki, nie wina wykonawcy.** Szczebel dyscypliny w drabinie POMINAL SIE z powodem
„brak zmienionych plikow kodu", choc przybylo 256 wierszy w trzech plikach — walidator nie
obejrzal tej zmiany ani razu. Osobno: `budowa/klient/node_modules` nie istnieje, wiec
`contract.ts` nie zostal skompilowany przez nikogo, mimo ze osiem plikow klienta go wciaga.
Obie luki idą do zgloszen, nie do tego terenu.

**Kryteria odbioru.**
1. `narzedzia.pozycje` ma 372 pozycje (340 + 32). Ani jedna nie wskazuje komendy
   nieistniejacej, ani jedna nie powtarza komendy juz wystawionej, ani jedna nie
   ma pustego `zastosowanie`. Sprawdzone przeliczeniem, nie oszacowaniem.
2. Zadna z komend objetych wyjatkiem nie zostala wystawiona. Wypisz je i pokaz brak.
3. `contract.go` i `contract.ts` powstaly GENERATOREM (`budowa/shared/gen/`),
   nie recznie. Plik generowany edytowany recznie jest zwrotem terenu.
4. `bash narzedzia/drabina.sh szybka` przechodzi, w tym szczebel swiezosci wytworow.
5. Sprawdziany zgodnosci kontraktu w rdzeniu przechodza:
   `go test -run Kontrakt ./internal/core/`.
6. Zdanie `zastosowanie` kazdej nowej pozycji jest po polsku, bez ogonkow tam,
   gdzie sasiednie pozycje ich nie maja — sprawdz ksztalt pozycji juz stojacych
   i powtorz go, zamiast wprowadzac drugi styl.


### straznik-zapory-designu

| | |
|---|---|
| **Galaz** | `teren/straznik-zapory-designu` z `main` |
| **Drzewo** | `~/robocze/straznik-zapory-designu` |
| **Wykaz plikow** | `budowa/server/internal/core/zapora_fotografii_test.go` |
| **Poza terenem** | wszystkie pliki `adapter_modul_design*`, kontrakt, migracje, `prowadzenie/` |

**Przedmiot.** Usterka ZYWA NA `main`, znaleziona niezaleznie przez dwoch kontrolerow
zamknietego terenu mermaid. `TestObszarDesignNieWymieniaSilnikowObrazuSpozaInstalki`
(`zapora_fotografii_test.go:99-125`) **nie liczy przejrzanych plikow**,
ktory jego blizniak ma w wierszach 48, 62 i 91 (`sprawdzonych++` oraz
`if sprawdzonych == 0 { t.Fatal(...) }`).

Skutek zmierzony zmiana probna: w katalogu bez plikow o przedrostku `adapter_modul_design`
blizniak oblewa komunikatem „zapora nie znalazla ani jednego pliku", a ten test daje
`--- PASS (0.00s)`, nie sprawdziwszy niczego. Zmiana przedrostka nazw albo przeniesienie
obszaru Design do podkatalogu zdejmie zapore **bez jednego czerwonego sprawdzianu**.

To jest ten sam rodzaj usterki, ktory dzis wraca w budowie raz po raz: sprawdzian, ktory
przy braku materialu konczy sie bez bledu zamiast oblac. Wylaczona miara wyglada
wtedy jak wynik dobry.

**Kryteria odbioru.**
1. Test oblewa GLOSNO, gdy nie przejrzal ani jednego pliku. Dowodem jest zmiana probna:
   zmien filtr przedrostka w kopii poza drzewem i pokaz, ze test **pada**.
2. Przed naprawa ta sama mutacja daje `PASS` — przytocz oba biegi.
3. Przejrzyj CALY plik zapory i wypisz KAZDA petle oraz KAZDY wykaz, ktory moze
   przejsc na zerze przebiegow. Nie naprawiasz jednego miejsca, gdy sasiednie ma
   te sama dziure — ale i nie wychodzisz poza ten plik.
4. `go test -run 'Zapor|Design|Warsztat|Obszar' ./internal/core/` przechodzi.
   **Poprawka Prowadzacego 28.08:** poprzedni wybieg `-run 'Zapor|Design'` NIE WYBIERAL
   `TestWarsztatFotografiiNieWolaProcesu` — sprawdzone uruchomieniem, wybieral cztery inne
   sprawdziany, a glownej zapory nie dotykal. Bramka terenu przepuszczala stan, w ktorym
   pliki Design wolaja proces wlasna droga. Blad Prowadzacego, nie wykonawcy.

**RUNDA DRUGA — praca zwrocona przez obie kontrole 28.08.** Zalatana zostala JEDNA petla,
a trzy wykazy obok zostaly bez oslony — kazdy obalony uruchomieniem:
`silnikiSpozaInstalki` (`:12-20`, petla `:120`), `wolaniaProcesu` (`:23-28`, petla `:82`)
oraz `wolaniaWlasneProcesu` (`:31-35`, petla `:66`). Oproznienie ktoregokolwiek daje bieg
zdany. Zaden nie ma odpowiednika `if len(...) == 0 { t.Fatal(...) }`.

Do tego oslona zera nie wystarcza: pilnuje wylacznie przypadku, w ktorym zniknely WSZYSTKIE
pliki. Pokazane uruchomieniem — 40 z 41 plikow Design przemianowanych, `imagemagick`
i `potrace` wpisane do jednego z pozostalych, oba sprawdziany zdane. Miara ma trzymac
PODLOGE liczby plikow albo porownanie z rzeczywistym wykazem plikow obszaru, nie samo zero.


### domkniecie-obrobki-wstepnej

| | |
|---|---|
| **Galaz** | `teren/domkniecie-obrobki-wstepnej` z `teren/obrobka-wstepna-dokumentow`, scalona z `main` przy zalozeniu |
| **Drzewo** | `~/robocze/domkniecie-obrobki-wstepnej` |
| **Wykaz plikow** | `budowa/server/internal/core/adapter_narzedzia_dokument_tekst.go`; `budowa/server/internal/core/skutek_narzedzi_tresci_pisanej_test.go`; `budowa/server/internal/core/adapter_modul_badania_lektura.go`; `budowa/server/internal/core/adapter_modul_badania_lektura_test.go` |
| **Poza terenem** | kontrakt i jego wytwory, `adapter_modul_studio_cyfryzacja.go`, migracje, `prowadzenie/` |

**Przedmiot.** Teren `obrobka-wstepna-dokumentow` przyjeto: `preprocess` dziala na drodze
`document.text.extract`, sprawdzian mierzy roznice wyniku i pada przed naprawa. Kontrola
nazwala jednak szesc oslabien, a osobny werdykt — jedna usterke BLOKUJACA scalenie.
Ten teren domyka je wszystkie, zanim lancuch wejdzie do `main`.

**1. USTERKA BLOKUJACA — wielojezycznosc rozsypuje sie na dwoch jezykach.**
`adapter_modul_badania_lektura.go:99-101` sklada jezyki przez `strings.Join(jezyki, "+")`,
a `jezykRozpoznaniaDokumentu` (`adapter_narzedzia_dokument_tekst.go:342-353`) rozpoznaje
wylacznie nazwy pojedyncze i lancucha po `+` nie rozcina. `Languages: ["pl","en"]` daje
`-l pl+en`. Maszyna niesie wylacznie `eng`, `osd`, `pol`. Zadanie POPRAWNE WOBEC KONTRAKTU
konczy sie cicha awaria rozpoznania. Jeden jezyk przechodzi, wiec usterka czeka na drugi.

**2. Dwie prawdy o wierszu wywolania unpapera.** `argumentyObrobkiWstepnejDokumentu`
(`adapter_narzedzia_dokument_tekst.go:397-401`) powiela argumenty z `argumentyCzyszczenia`
(`adapter_modul_studio_cyfryzacja.go:454-468`) w wariancie „wszystko wlaczone". Komentarz
zapewnia „tak samo jak w drodze Studia" — dzis prawda, ale NIC TEGO NIE WIAZE. Zmiana
jednej strony rozejdzie sie po cichu.

**3. Sprawdzenie odmowy nie biegnie na zadnej maszynie drabiny.**
`skutek_narzedzi_tresci_pisanej_test.go:357-359` pomija sie sam wszedzie tam, gdzie
unpaper STOI — czyli na kazdej maszynie weryfikacji. W drabinie wypada SKIP, nie PASS.
Regresja w tresci odmowy przejdzie niezauwazona.

**4. Sprawdzian potrafi sam sie rozbroic.** `:329-334` niesie `t.Skip("pomiar
bezprzedmiotowy")` na wypadek, gdyby rozpoznanie odczytalo material pochylony bez obrobki.
Mocniejszy tesseract albo lagodniejszy skos zamieni miare w zielony SKIP. Rozbrojenie
bedzie wygladalo jak powodzenie — a to jest grozniejsze od bledu.

**5. Droga PDF nieprzemierzona z polem `preprocess`.** `:255-256, 301-302, 321` przekazuja
`obrobkaWstepna` przez `tekstZPdf` i `rozpoznajPismoWPdf`; kompiluje sie, ale zaden
sprawdzian ta galezia nie idzie. Zmierzona jest wylacznie droga obrazu.

**6. Odmowa zapada po pelnym renderze.** `zewnetrzne.Stoi(narzedzieCzyszczeniaSkanu)`
siedzi w `rozpoznajPismo` (`:355-364`), wiec na drodze PDF Operator bez unpapera placi
renderem wszystkich stron w 300 dpi, zanim zobaczy odmowe. Studio pyta o program PRZED praca.

**7. Sprawdzian, ktory klamie o sobie.** `adapter_modul_badania_lektura_test.go:29-32`
fabrykuje roznice wewnatrz atrapy (`if z.Preprocess != nil && *z.Preprocess { tekst =
"tekst po obrobce wstepnej" }`), a komentarz w wierszach 62-65 oglasza, ze „mierzy skutek".
Miare skutku na drodze produkcyjnej wnosi juz teren poprzedni. Ten sprawdzian albo znika,
albo przestaje twierdzic o sobie rzecz nieprawdziwa.

**RUNDA DRUGA — praca zwrocona przez kontrole awarii 28.08.** Kontroler miary przyjal,
kontroler awarii obalil. Punkty 1-7 wykonane, ale naprawa wniosla nowa brame, ktora sama
odmawia. Naprawiasz brame, nie przerabiasz calosci.

| | Rzecz | Gdzie |
|---|---|---|
| A | `osd` przechodzi weryfikacje jezyka, choc jest danymi ORIENTACJI PISMA, nie jezykiem. `language: "osd"` konczy sie odpowiedzia UDANA z trescia bedaca smieciami i `usedOcr: true` — odczyt zmyslony podany Operatorowi jako rozpoznanie. Wykaz z `--list-langs` NIE JEST wykazem jezykow rozpoznania. | `adapter_narzedzia_dokument_tekst.go:463-478` |
| B | Pusty odczyt `--list-langs` przy `err == nil` jest brany za prawde o maszynie: `wiersze[1:]` na wyniku jednoelementowym daje pusta mape, po czym KAZDY jezyk zostaje odrzucony. Brak dwoch oslon: wynik pusty ma dac odmowe zamiast isc dalej, a pierwszy wiersz ma byc sciety dopiero po sprawdzeniu, ze to naglowek. Prowadzenie NIE odtworzylo czestosci podanej przez kontrolera (40 przebiegow sekwencyjnych — wszystkie zdane), ale brak oslony jest w kodzie i zamienia kazde potkniecie w bledna odmowe na drodze produkcyjnej. | `:483-497` |
| C | Czlon pusty (`["pl",""]`) daje odmowe, ktora niczego nie nazywa: tresc brzmi „nie niesie danych jezykowych  — " i konczy sie „dla ". | `:470-475` |
| D | Przebieg odniesienia, gdy ODMOWI, po cichu opuszcza cala roznice przed/po — sprawdzian schodzi z pomiaru ROZNICY do pomiaru OBECNOSCI slowa. Ten sam gatunek co `t.Skip` z punktu 4, tyle ze bez sladu w wyjsciu drabiny. | `skutek_narzedzi_tresci_pisanej_test.go:445` i `:330` |
| E | Zaden sprawdzian nie przejezdza pelnego lancucha `research.source.ocr` z ladunkiem `Languages: ["pl","en"]` — a to jest ladunek nazwany w kryterium 1. Dzis mierzy sie forma juz poprawna `["pol","eng"]`, czyli dokladnie tak, jak usterka przeszla poprzednio. | `adapter_modul_badania_lektura_test.go:95` |
| F | Odmowa czlonu wewnatrz wykazu zlozonego (`"pol+xx"`) nie jest zmierzona zadnym uruchomieniem — tylko czlon pojedynczy. | `skutek_narzedzi_tresci_pisanej_test.go:379-395` |

**Kryterium dodatkowe rundy drugiej:** `go test -count=20` na sprawdzianach tego obszaru
przechodzi ZA KAZDYM przebiegiem. Sprawdzian oblewajacy przygodnie nie jest zdany.

**Kryteria odbioru.**
1. `Languages: ["pl","en"]` daje Tesseractowi wykaz, ktory ten rozumie, albo ODMOWE NAZWANA
   dla jezyka, ktorego maszyna nie niesie. Cicha awaria nie jest dopuszczalna. Sprawdzian
   mierzy to na drodze produkcyjnej i PADA przed naprawa.
2. Argumenty unpapera maja jedna prawde: albo jedno zrodlo, albo sprawdzian wiazacy obie
   strony, ktory pada przy rozejsciu.
3. Zaden sprawdzian tego obszaru nie przechodzi milczkiem tam, gdzie mial mierzyc.
   Rozstrzygniecie o drodze nalezy do wykonawcy i ma byc nazwane.
4. Droga PDF z polem `preprocess` jest zmierzona uruchomieniem, nie czytaniem.
5. Odmowa braku programu zapada PRZED kosztowna praca, tak jak w drodze Studia.
6. `go build ./...` oraz sprawdziany obszaru tresci pisanej przechodza.


### miara-skutku-twarzy

| | |
|---|---|
| **Galaz** | `teren/miara-skutku-twarzy` z `main` |
| **Drzewo** | `~/robocze/miara-skutku-twarzy` |
| **Wykaz plikow** | `budowa/server/internal/core/skutek_odtwarzania_twarzy_test.go` |
| **Poza terenem** | pomocnik pythonowy, adaptery, kontrakt, migracje, `prowadzenie/` |

**Przedmiot.** Teren `silnik-twarzy` przyjeto, ale jego kontroler nazwal trzy dziury
w samej mierze — sprawdzian przechodzi, a nie dowodzi tego, co obiecuje:

1. `skutek_odtwarzania_twarzy_test.go:56-65` — przy braku wag `t.Skipf` przepuszcza
   sprawdzian milczkiem. Na maszynie bez `/opt/danaco-modele/twarze` cala drabina
   swieci zielono, a kryterium nie jest zmierzone ANI RAZU. Bramka wydania, ktora
   sama siebie wylacza przy braku materialu, nie jest bramka.
2. `:114-141` — roznice liczy sie na CALYM obrazie. Pomocnik, ktory globalnie
   przyciemnilby albo rozmyl obraz, przeszedlby oba progi. Miara jest skutkiem
   NIEUMIEJSCOWIONYM: nie sprawdza, czy zmiana siedzi w wycinku twarzy.
3. `:96` — liczba „twarze poprawione" pochodzi z licznika wypisanego przez sam
   pomocnik (`adapter_narzedzia_obraz_pomocnik_twarzy.py:89`). Sprawdzian bada
   SAMOOPIS pomocnika, nie wielkosc zmierzona po stronie rdzenia.

**Kryteria odbioru.**
1. Brak wag nie daje cichego pominiecia: sprawdzian albo mierzy, albo GLOSNO oblewa.
   Rozstrzygniecie o drodze nalezy do wykonawcy i ma byc nazwane.
2. Roznica jest mierzona W WYCINKU TWARZY, a nie na calym obrazie. Dowodem jest
   przeciwsprawdzian: material zmieniony globalnie (przyciemniony albo rozmyty)
   NIE przechodzi progu.
3. Zadna liczba rozstrzygajaca o wyniku nie pochodzi z samoopisu pomocnika.
4. `DANACO_MODELE=/opt/danaco-modele go test -run Twarz ./internal/core/` przechodzi.


### obrobka-wstepna-w-porcie-dokumentow

| | |
|---|---|
| **Galaz** | `teren/obrobka-wstepna-dokumentow` z `teren/wpiecie-unpaper` — nie z `main`, bo praca zada pola `preprocess`, ktore wnosi tamten teren |
| **Drzewo** | `~/robocze/obrobka-wstepna-dokumentow` |
| **Wykaz plikow** | `budowa/server/internal/core/adapter_narzedzia_dokument_tekst.go` wraz ze sprawdzianami; `budowa/server/internal/core/skutek_narzedzi_tresci_pisanej_test.go` |
| **Poza terenem** | kontrakt i jego wytwory, `adapter_modul_badania_lektura.go`, `adapter_modul_studio_cyfryzacja.go`, migracje, `prowadzenie/` |

**Przedmiot.** Teren `wpiecie-unpaper` doprowadzil pola `preprocess` i `languages`
do granicy portu: `RozpoznajPismoZrodla` przekazuje je w `shared.DocumentTextExtractRequest`
do `Dokumenty.WyciagnijTekst`. Za ta granica pole jest jednak martwe —
`adapterNarzedziDokumentu.WyciagnijTekst` w ogole go nie czyta i nie wola programu
unpaper. Program jest dzis wolany wylacznie w module Studio, droga `studio.ingest.recognize`.

Skutek: obrobka wstepna pozostaje nieosiagalna z drogi `document.text.extract`,
a wiec i z `research.source.ocr`, mimo ze kontrakt i uchwyt ja niosa.

**Rzeczy zmierzone, ktore skracaja prace.** `adapterNarzedziDokumentu` ma juz metode
`wolaj` (`adapter_narzedzia_dokument.go:128`), wiec wolanie programu jest osiagalne
bez zmiany struktury adaptera. Funkcje `argumentyCzyszczenia`, `materialWPnm`,
`czyPnm` i `zapiszPpm` sa wolnymi funkcjami pakietu `core`, nie metodami adaptera
Studia — nadaja sie do ponownego uzycia bez przenoszenia kodu.

**Kryteria odbioru.**
1. `WyciagnijTekst` honoruje `z.Preprocess`: przy wartosci prawdziwej material
   przechodzi przez unpaper przed rozpoznaniem, przy falszywej i przy braku — nie.
2. Sprawdzian mierzy ROZNICE WYNIKU rozpoznania na materiale przekrzywionym,
   wzorem `TestObrobkaWstepnaProstujeSkosPrzedRozpoznaniem`, ktory robi to samo dla
   drogi Studia. Sprawdzian ma isc DROGA PRODUKCYJNA przez `montaz_porty.go`,
   nie przez atrape portu. Sprawdzian mierzacy samo wywolanie zamiast roznicy wyniku
   jest brakiem sprawdzianu.
3. Brak programu unpaper jest odmowa nazwana, nie cicha praca bez obrobki.
4. `go build ./...` oraz sprawdziany pakietu `core` dotyczace tresci pisanej przechodza.


### rama-aplikacji-w-kliencie

| | |
|---|---|
| **Galaz** | `teren/rama-aplikacji` z `main` |
| **Drzewo** | `~/robocze/rama-aplikacji` |
| **Wykaz plikow** | `budowa/klient/src/rama/` (nowy katalog) wraz ze sprawdzianami; `budowa/klient/src/aplikacja.ts` (nowy); `budowa/klient/index.html` i `budowa/klient/arkusze.css` wylacznie w zakresie wpiecia ramy; `budowa/klient/src/main.ts` (zniesiony na rzecz `aplikacja.ts` — dwa punkty wejscia otwieralyby dwa polaczenia i dwie maszyny stanu); `budowa/klient/package.json` wylacznie w zakresie dopisania sprawdzianow ramy do polecenia `testy` |
| **Poza terenem** | `budowa/klient/src/wejscie/`, `polaczenie/`, `protokol/`; caly `design/`; rdzen; `prowadzenie/` |

**Runda trzecia — wykaz zmian probnych jako kryterium odbioru.** Dwie niezalezne kontrole
zwrocily runde druga. Kazda z ponizszych zmian probnych przechodzi DZIS bez bledu; po naprawie
KAZDA ma oblewac sprawdzian. To jest miara odbioru, nie wskazowka.

| | Zmiana probna | Co pokazuje |
|---|---|---|
| M1 | wstawic z powrotem `przekazano = true` przed wolaniem przekazania | sprawdzian mierzy WLASNA KOPIE znacznika jednorazowosci (`przekazanie.test.ts:266-287`), nie kod produkcyjny; `aplikacja.ts` nie jest importowana przez zaden sprawdzian |
| M2 | zapisac `liczbaSesji: 0` na sztywno w `montaz.ts:40` i `:75` | nastawa niesie `sessions: []`, wiec sprawdzian porownuje 0 z 0 |
| M3 | `w.moduly.slice(0, 1)` w `skladniki/szyna.ts:43` | nastawa niesie JEDEN modul, wiec szyna z wykazu jest nieodrozznialna od szyny rysujacej pierwszy element |
| M4 | wyciac znak motywu i wezel `[data-stan-motyw]` | zaden sprawdzian po niego nie siega; `motywCiemny()` nie wykonuje sie ani razu |

Nastawy sprawdzianow musza NIESC MATERIAL, ktory rozroznia: co najmniej dwa moduly
i co najmniej jedna sesja. Nastawa pusta zamienia miare w porownanie zera z zerem.

**Cicha awaria, ktora runda druga PRZESUNELA, zamiast usunac.** `przekazanie.ts:41-49`:
`zdejmijOknoWejscia()` i `scenaWejscia.hidden = true` wykonuja sie PRZED `zamontujRame`,
a `miejsceRamy.hidden = false` dopiero po nim. Kazdy wyjatek montazu zostawia OBA wezly
ukryte — Operator dostaje bialy ekran. Wyjatek polyka `przebieg.ts:333-338`
(`try { sluchacz(stan) } catch { console.error }`), znacznik zostaje `false`, wiec proba
wraca przy kazdej zmianie stanu i rzuca od nowa. Zadnej odmowy nazwanej.
Pokazane uruchomieniem, nie czytaniem.

**Pozostale drogi bez odmowy:** `montaz.ts:46` i `skladniki/szyna.ts:43` robia
`w.moduly.map(...)` bez sprawdzenia — rdzen oddajacy `environment.enter` bez `modules`
albo z `null` daje to samo; `montaz.ts:66` zaklada `querySelector(...)!` zamiast odmowic;
`narzedzia.ts:50` rzuca wyjatek zamiast nazwac odmowe, i to w trakcie montazu.

**Rama nie subskrybuje przebiegu** (`montaz.ts`): liczba sesji i nazwa srodowiska sa
zamrozone na chwile montazu. Zerwanie polaczenia prowadzi przez `przebieg.ts:356-362`
cala droge wejscia od nowa, ale scena jest ukryta, a jej sluchacz odsubskrybowany.
Operator zostaje w ramie, ktora nie ma jak nic powiedziec, bez drogi powrotu.
Rozstrzygniecie o zakresie naprawy nalezy do wykonawcy i ma byc NAZWANE — jesli
subskrypcja wykracza poza ten teren, ma powstac zgloszenie, a nie milczenie.

`src/rama/dom-zastepczy.ts:1` — naglowek 355 znakow przy limicie 350.

**Do wykonania przy nastepnej rundzie, wina prowadzenia, nie wykonawcy.**
Nazwa wystawiona `utworzPrzekazanieJednorazowe` (`src/rama/przekazanie.ts:52`) zastapila
przenosnie, ktora wniosl do zlecenia Prowadzacy, a wykonawca ja stamtad wzial.
Rzecz jest znacznikiem jednorazowosci: `let przekazano = false` ustawiane na `true`
DOPIERO po powodzeniu, zeby odmowa nie zamykala drogi na kolejna zmiane stanu.
Nazwa ma opisywac te funkcje — `utworzPrzekazanieJednorazowe` albo blizsze temu,
co w kliencie juz stoi. Razem z nia schodza: komentarz w `:51`, wzmianka
w `aplikacja.ts:35` i siedem komunikatow sprawdzianow.

Zmiana NIE wchodzi w trakcie biegu kontroli — kontrolerzy uruchamiaja sprawdziany
na tym drzewie i podmiana nazw w locie zepsulaby ich pomiar.

**Przedmiot.** Droga wejscia konczy sie komenda `environment.enter` i okno zostaje
na ekranie na zawsze, bo nie ma dokad prowadzic. Klient bierzacy liczy 51 plikow
i wola piec komend z tysiaca osiemdziesieciu jeden. Teren stawia RAME APLIKACJI:
powierzchnie, w ktora droga wejscia przekazuje sterowanie po wejsciu do srodowiska.

Zrodlem ksztaltu jest przyjety prototyp `design/05-okna/przeplyw/centrum-dowodzenia.html`
wraz z przyjetym prototypem okna centrum dowodzenia. Rama obejmuje
trzy strefy prototypu: szyne nawigacji, belke tytulowa i pas stanu. Zawartosc okna
roboczego NIE nalezy do tego terenu.

**Kryteria odbioru.**
1. Po `environment.enter` droga wejscia oddaje sterowanie ramie, a scena wejscia
   schodzi. Sprawdzian mierzy przejscie, nie sama obecnosc funkcji.
2. Rama bierze barwy, odstepy i pismo WYLACZNIE z zetonow warstwy projektowej.
   Zero wartosci wpisanych liczbowo — wykazane przeszukaniem z sonda dodatnia.
3. Zadnego lancucha widocznego dla Operatora poza katalogiem tresci — tak samo,
   jak pilnuje tego droga wejscia swoim sprawdzianem.
4. `npm run typy`, `npm run testy` i `npm run budowanie` w `budowa/klient` przechodza.
5. Rama nie wola zadnej komendy kontraktu poza tymi, ktore juz sa w kliencie —
   okna modulowe to osobny teren.


### silnik-twarzy

| | |
|---|---|
| **Galaz** | `teren/silnik-twarzy` z `main` |
| **Drzewo** | `~/robocze/silnik-twarzy` |
| **Wykaz plikow** | `budowa/server/internal/core/adapter_narzedzia_obraz_model_twarze.go` wraz ze sprawdzianami; `budowa/server/internal/core/adapter_narzedzia_obraz_model_silniki.go`; pomocnik pythonowy `budowa/server/internal/core/adapter_narzedzia_obraz_pomocnik_twarzy.py` |
| **Poza terenem** | kontrakt, migracje, `internal/wiedza/`, `prowadzenie/` |

**Przedmiot.** Pole `faces` komendy `image.upscale` istnieje w kontrakcie, wagi
GFPGAN i CodeFormer stoja w `/opt/danaco-modele/twarze` (692 MB), a silnika nie ma:
wydanie ncnn nie niesie sieci twarzowej, wagi `.pth` zadaja stosu torch, ktorego
rdzen nie wola. Odmowa mowi dzis o braku `gfpgan-ncnn-vulkan`.

**Kryteria odbioru.**
1. `image.upscale` z `faces: true` oddaje obraz z poprawiona twarza, a nie odmowe.
2. Sprawdzian zdolnosci mierzy SKUTEK na obrazie, nie koperte odpowiedzi.
3. Brak wag jest odmowa nazwana, nie panika ani cicha praca bez poprawki.
4. `go build ./...` oraz sprawdziany pakietu `core` przechodza.

### wpiecie-unpaper

| | |
|---|---|
| **Galaz** | `teren/wpiecie-unpaper` z `main` |
| **Drzewo** | `~/robocze/wpiecie-unpaper` |
| **Wykaz plikow** | `budowa/shared/contract.json` wraz z wytworzonymi `contract.go` i `contract.ts`; `budowa/server/internal/core/adapter_modul_badania_lektura.go` wraz ze sprawdzianami |
| **Poza terenem** | migracje, `internal/wiedza/`, warstwa projektowa, `prowadzenie/` |

**Przedmiot.** Uchwyt `adapter_modul_badania_lektura.go:589` upuszcza pola
`preprocess` i `languages`, a `DocumentTextExtractRequest` kontraktu ich nie ma.
Obrobka wstepna unpaper jest nieosiagalna z drogi badan, a pole `preprocess`
pozostaje martwe.

**Kryteria odbioru.**
1. Kontrakt niesie pola, ktorych uchwyt zada; generatory `contract.go` i `contract.ts`
   przebudowane tym samym poleceniem, ktore repozytorium juz ma.
2. `research.source.ocr` z obrobka wstepna daje wynik rozny od wyniku bez niej —
   sprawdzian mierzy roznice, nie samo przejscie.
3. Suma kontrolna kontraktu odnotowana w raporcie; sprawdzian swiezosci generatorow przechodzi.

### trwalosc-wektorow-obrazu

| | |
|---|---|
| **Galaz** | `teren/trwalosc-wektorow-obrazu` z `main` |
| **Drzewo** | `~/robocze/trwalosc-wektorow-obrazu` |
| **Wykonawca** | sesja wykonawcza, kontrola osobna |
| **Wykaz plikow** | `budowa/server/internal/wiedza/obraz.go` wraz z jego sprawdzianami; nowa migracja `budowa/server/internal/store/migracja_*.sql`; `budowa/server/internal/core/adapter_modul_wiedza_obraz.go` wylacznie w zakresie nazwania obciecia wykazu |
| **Poza terenem** | wszystko inne, w szczegolnosci migracje juz zapisane, `internal/core/` poza wymienionym plikiem, katalog `prowadzenie/` |

**Przedmiot.** Wektor osi obrazu liczy sie przy kazdym zapytaniu: `Dopasuj` wola
pomocnika Pythona bez jednego odczytu z bazy, a kolumna `zakres` tabeli wskaznika
dopuszcza trzy wartosci i nie ma miejsca na wektor obrazu. Zmierzone: 17 s dla
trzech obrazow.

**Kryteria odbioru.** Sprawdzalne uruchomieniem, nie odczytem:

1. Drugie zapytanie o ten sam obraz NIE wola pomocnika zewnetrznego. Dowodzi tego
   sprawdzian liczacy wolania atrapa pomocnika; sprawdzian ma padac przed naprawa.
2. Zmiana nastawy `wiedza_model_obrazu` uniewaznia zapisany wektor. Wektor
   policzony innym modelem opisuje znaczenie w innej przestrzeni, wiec jego ponowne
   uzycie jest usterka ciezsza niz liczenie od nowa.
3. Zadna migracja juz zapisana nie zostaje zmieniona. Nowy stan wchodzi nowa
   migracja o kolejnym wolnym numerze.
4. Obciecie wykazu sufitem `GranicaObrazow` jest NAZWANE w odpowiedzi. Dzis pole
   `examined` mowi tylko, ile weszlo, wiec biblioteka o pieciu tysiacach obrazow
   jest nieodrozninalna od biblioteki o dwustu.
5. `go build ./...` oraz `DANACO_MODELE=/opt/danaco-modele go test -count=1 -timeout 20m
   ./internal/wiedza/ ./internal/store/` koncza sie powodzeniem.

**Rozstrzygniecia pozostawione wykonawcy**, do nazwania w raporcie: czy wektor
mieszka w nowej wartosci kolumny `zakres`, czy w osobnej tabeli — blizej tego, co
juz stoi dla wektorow tekstu.


**Orkiestracja czterech terenow komentarzy — rytm zapisu i wznowienie po
urwaniu.** Pisarz robi rewizje PO KAZDYM PLIKU (kod + przyrost docs razem);
dziennik trwaly stoi w `~/robocze/prowadzenie/komentarze/` (pomiar bazowy,
manifest porcji, dzienniki torow `<sektor>.jsonl` z liniami START/DONE).
Granica liczy komentarz glowny: odsylacze, wskazania i dyrektywy sa
z niej wylaczone (dopowiedzenie Wlasciciela w pozycji 18; wciela je
instrument). Granica jest schodkowa, nie proporcjonalna: 250 znakow na plik,
od dwoch tysiecy wierszy 250 za kazdy pelny tysiac (dopowiedzenie trzecie
pozycji 18; instrument wciela). Kazdy plik po pracy zaczyna sie naglowkiem: pelne zdanie
odpowiedzialnosci pliku (dopowiedzenie czwarte pozycji 18); plik niemy to
uchybienie zwracajace porcje. Komentarz jest kompletnym nosnikiem informacji: zakazane sa ODESLANIA
do opracowan (`*.md`, `docs/`, `prowadzenie/`, zwroty „Patrz", „Zob.",
„Uzasadnienie:", „szczegoly w"), dozwolone jest NAZYWANIE artefaktow
technicznych — bibliotek, arkuszy stylu, programow, migracji, plikow nastaw
i plikow zrodlowych (dopowiedzenia szoste i siodme; instrument liczy
odeslania jako POWOLANIA i plik z odeslaniem nie przechodzi); odsylacze i wskazania zastane zdejmuje sie
w toku prac; konwencja mapowania do docs zapisana raz w CLAUDE.md. Zakaz skrotow upychajacych tresc — wylacznie pelne zdania
i pelne slowa, w kodzie i w docs (dopowiedzenie drugie pozycji 18).
Stan jest odtwarzalny bez kontekstu zadnej sesji: plik ZROBIONY =
zacommitowany, W-GRANICY i tokeny tozsame z baza `285f4bf`; DO ZROBIENIA =
PONAD granica (pomiar `narzedzia/zrodlo-bez-komentarzy.go -gestosc`);
podejrzany po urwaniu = ostatni START bez DONE w dzienniku toru →
`git checkout -- <plik>` i plik wraca do puli. Wznowienie: pomiar
pozostalosci sektora → nowa orkiestracja na resztce manifestu.


### komentarz-studio-design

Przepisanie komentarzy do standardu zawodowego wraz z domknieciem granicy
z pozycji 18 rejestru decyzji. **Kodu nie zmieniamy.**

| | |
|---|---|
| **Galaz** | `teren/komentarz-studio-design` z `main` |
| **Wykaz plikow** | `internal/core/adapter_modul_studio_*.go`, `adapter_modul_design_*.go`, pakiety `internal/wiedza/`, `internal/mowa/`, `internal/poczta/`, `internal/zewnetrzne/` |
| **Skala** | 125 plikow, 854 tys. znakow komentarza (pomiar instrumentem, 27.08) |
| **Poza terenem** | wszystko inne, w tym `budowa/shared/`, `budowa/klient/`, `budowa/desktop/`, `design/`, `prowadzenie/` oraz pliki terenow `komentarz-moduly`, `komentarz-rdzen`, `komentarz-pakiety` |

**Przedmiot.** Komentarz stwierdza regule obowiazujaca — nie waży wariantow, nie
zwraca sie do czytelnika, nie prowadzi wykladu i nie relacjonuje przebiegu prac.
Uzasadnienie dluzsze niz zdanie idzie do `docs/`, a plik niesie zdanie i odsylacz.
Jezyk, terminologia i forma wedle standardu zawodowego, jednolite w calym zakresie.

**Kryteria odbioru.**

1. **Kod nietkniety** — dla kazdego zmienionego pliku strumien tokenow bez
   komentarzy (`narzedzia/zrodlo-bez-komentarzy.go`) bajtowo identyczny przed
   i po, a wykaz dyrektyw (`grep -n '^[[:space:]]*//go:'`) tozsamy. Binarium
   dowodem NIE jest — zmierzone 27.08: jedna dodana linia komentarza zmienia
   binarium, bo numery wierszy funkcji wchodza do tablic sladow stosu.
2. **Kazdy plik zakresu w granicy 250 znakow na 1000 wierszy** — wykazane
   pomiarem calego zakresu, z podaniem liczby plikow przed i po.
3. **Uzasadnienie niosace tresc nie zniklo** — przeniesione do `docs/`, a plik
   kodu odsyla do niego. Tresc, ktora byla narracja albo powtorzeniem
   oczywistosci, znika i **mowisz o tym wprost** w raporcie.
4. `gofmt -l` na zakresie zwraca pusto; `go vet` czysto.
5. `DANACO_MODELE=/opt/danaco-modele gotestsum -- -count=1 -timeout 40m ./...`
   — zero niepowodzen wobec stanu zastanego.
6. Kontrakt nietkniety — wykazane suma kontrolna.

### komentarz-moduly

Przepisanie komentarzy do standardu zawodowego wraz z domknieciem granicy
z pozycji 18 rejestru decyzji. **Kodu nie zmieniamy.**

| | |
|---|---|
| **Galaz** | `teren/komentarz-moduly` z `main` |
| **Wykaz plikow** | `internal/core/adapter_modul_*.go` **poza** `studio` i `design` |
| **Skala** | 262 plikow, 1072 tys. znakow komentarza (pomiar instrumentem, 27.08) |
| **Poza terenem** | wszystko inne, w tym `budowa/shared/`, `budowa/klient/`, `budowa/desktop/`, `design/`, `prowadzenie/` oraz pliki terenow `komentarz-studio-design`, `komentarz-rdzen`, `komentarz-pakiety` |

**Przedmiot.** Komentarz stwierdza regule obowiazujaca — nie waży wariantow, nie
zwraca sie do czytelnika, nie prowadzi wykladu i nie relacjonuje przebiegu prac.
Uzasadnienie dluzsze niz zdanie idzie do `docs/`, a plik niesie zdanie i odsylacz.
Jezyk, terminologia i forma wedle standardu zawodowego, jednolite w calym zakresie.

**Kryteria odbioru.**

1. **Kod nietkniety** — dla kazdego zmienionego pliku strumien tokenow bez
   komentarzy (`narzedzia/zrodlo-bez-komentarzy.go`) bajtowo identyczny przed
   i po, a wykaz dyrektyw (`grep -n '^[[:space:]]*//go:'`) tozsamy. Binarium
   dowodem NIE jest — zmierzone 27.08: jedna dodana linia komentarza zmienia
   binarium, bo numery wierszy funkcji wchodza do tablic sladow stosu.
2. **Kazdy plik zakresu w granicy 250 znakow na 1000 wierszy** — wykazane
   pomiarem calego zakresu, z podaniem liczby plikow przed i po.
3. **Uzasadnienie niosace tresc nie zniklo** — przeniesione do `docs/`, a plik
   kodu odsyla do niego. Tresc, ktora byla narracja albo powtorzeniem
   oczywistosci, znika i **mowisz o tym wprost** w raporcie.
4. `gofmt -l` na zakresie zwraca pusto; `go vet` czysto.
5. `DANACO_MODELE=/opt/danaco-modele gotestsum -- -count=1 -timeout 40m ./...`
   — zero niepowodzen wobec stanu zastanego.
6. Kontrakt nietkniety — wykazane suma kontrolna.

### komentarz-rdzen

Przepisanie komentarzy do standardu zawodowego wraz z domknieciem granicy
z pozycji 18 rejestru decyzji. **Kodu nie zmieniamy.**

| | |
|---|---|
| **Galaz** | `teren/komentarz-rdzen` z `main` |
| **Wykaz plikow** | `internal/core/*.go` **niezaczynajace sie** od `adapter_modul_` |
| **Skala** | 364 plikow, 1228 tys. znakow komentarza (pomiar instrumentem, 27.08) |
| **Poza terenem** | wszystko inne, w tym `budowa/shared/`, `budowa/klient/`, `budowa/desktop/`, `design/`, `prowadzenie/` oraz pliki terenow `komentarz-studio-design`, `komentarz-moduly`, `komentarz-pakiety` |

**Przedmiot.** Komentarz stwierdza regule obowiazujaca — nie waży wariantow, nie
zwraca sie do czytelnika, nie prowadzi wykladu i nie relacjonuje przebiegu prac.
Uzasadnienie dluzsze niz zdanie idzie do `docs/`, a plik niesie zdanie i odsylacz.
Jezyk, terminologia i forma wedle standardu zawodowego, jednolite w calym zakresie.

**Kryteria odbioru.**

1. **Kod nietkniety** — dla kazdego zmienionego pliku strumien tokenow bez
   komentarzy (`narzedzia/zrodlo-bez-komentarzy.go`) bajtowo identyczny przed
   i po, a wykaz dyrektyw (`grep -n '^[[:space:]]*//go:'`) tozsamy. Binarium
   dowodem NIE jest — zmierzone 27.08: jedna dodana linia komentarza zmienia
   binarium, bo numery wierszy funkcji wchodza do tablic sladow stosu.
2. **Kazdy plik zakresu w granicy 250 znakow na 1000 wierszy** — wykazane
   pomiarem calego zakresu, z podaniem liczby plikow przed i po.
3. **Uzasadnienie niosace tresc nie zniklo** — przeniesione do `docs/`, a plik
   kodu odsyla do niego. Tresc, ktora byla narracja albo powtorzeniem
   oczywistosci, znika i **mowisz o tym wprost** w raporcie.
4. `gofmt -l` na zakresie zwraca pusto; `go vet` czysto.
5. `DANACO_MODELE=/opt/danaco-modele gotestsum -- -count=1 -timeout 40m ./...`
   — zero niepowodzen wobec stanu zastanego.
6. Kontrakt nietkniety — wykazane suma kontrolna.

### komentarz-pakiety

Przepisanie komentarzy do standardu zawodowego wraz z domknieciem granicy
z pozycji 18 rejestru decyzji. **Kodu nie zmieniamy.**

| | |
|---|---|
| **Galaz** | `teren/komentarz-pakiety` z `main` |
| **Wykaz plikow** | `budowa/server/internal/` **poza** pakietami `core`, `wiedza`, `mowa`, `poczta`, `zewnetrzne` (te niesie teren studio-design), wraz z `budowa/server/cmd/` |
| **Skala** | 430 plikow, 1013 tys. znakow komentarza (pomiar instrumentem, 27.08) |
| **Poza terenem** | wszystko inne, w tym `budowa/shared/`, `budowa/klient/`, `budowa/desktop/`, `design/`, `prowadzenie/` oraz pliki terenow `komentarz-studio-design`, `komentarz-moduly`, `komentarz-rdzen` |

**Przedmiot.** Komentarz stwierdza regule obowiazujaca — nie waży wariantow, nie
zwraca sie do czytelnika, nie prowadzi wykladu i nie relacjonuje przebiegu prac.
Uzasadnienie dluzsze niz zdanie idzie do `docs/`, a plik niesie zdanie i odsylacz.
Jezyk, terminologia i forma wedle standardu zawodowego, jednolite w calym zakresie.

**Kryteria odbioru.**

1. **Kod nietkniety** — dla kazdego zmienionego pliku strumien tokenow bez
   komentarzy (`narzedzia/zrodlo-bez-komentarzy.go`) bajtowo identyczny przed
   i po, a wykaz dyrektyw (`grep -n '^[[:space:]]*//go:'`) tozsamy. Binarium
   dowodem NIE jest — zmierzone 27.08: jedna dodana linia komentarza zmienia
   binarium, bo numery wierszy funkcji wchodza do tablic sladow stosu.
2. **Kazdy plik zakresu w granicy 250 znakow na 1000 wierszy** — wykazane
   pomiarem calego zakresu, z podaniem liczby plikow przed i po.
3. **Uzasadnienie niosace tresc nie zniklo** — przeniesione do `docs/`, a plik
   kodu odsyla do niego. Tresc, ktora byla narracja albo powtorzeniem
   oczywistosci, znika i **mowisz o tym wprost** w raporcie.
4. `gofmt -l` na zakresie zwraca pusto; `go vet` czysto.
5. `DANACO_MODELE=/opt/danaco-modele gotestsum -- -count=1 -timeout 40m ./...`
   — zero niepowodzen wobec stanu zastanego.
6. Kontrakt nietkniety — wykazane suma kontrolna.


### uprzaz-komend-neuronowych

Generyczna uprzaz sprawdzianow urywa kazda komende liczaca modelem, wiec zadnej
z nich nie da sie zmierzyc jej droga. Dwa tereny obeszly to wlasna droga i oba
zglosily to jako obejscie.

| | |
|---|---|
| **Galaz** | `teren/uprzaz-komend-neuronowych` z `main` |
| **Wykaz plikow** | `budowa/server/internal/core/zgodnosc_kontraktu_test.go`, `budowa/server/internal/core/blokady_skutek_test.go` |
| **Poza terenem** | wszystko inne |

**Zmierzone.** Granica uprzezy wynosi dzis 60 s (`granicaWykazuUrzadzen` + 15 s).
`image.upscale` potrzebuje **90 s bez twarzy i 190 s z twarzami**, przesiew
wyszukiwania okolo **30 s** na wczytanie wag. Wartosc jednolita nie da sie
dobrac: 190 s razy ponad tysiac komend to bieg liczony w godzinach.

**Kryteria odbioru.**

1. Komenda liczaca modelem daje sie zmierzyc uprzezа — wykazane sprawdzianem,
   ktory **przechodzi** dla takiej komendy, a **zawodzi** przy granicy sprzed
   zmiany.
2. Bieg calego pakietu `internal/core` **nie wydluza sie** ponad to, co dzis —
   zmierzone przed i po, obie liczby przytoczone.
3. Rozroznienie komend nie jest wykazem imion wpisanym recznie, albo jest —
   i wtedy powod stoi w komentarzu, a wykaz ma jedno miejsce.
4. `gotestsum -- -count=1 -timeout 40m ./...` — zero niepowodzen.


### centrum-poprawki

| | |
|---|---|
| **Galaz** | `teren/centrum-poprawki` z `f8f05b1` |
| **Drzewo** | `~/budowa` — drzewo glowne, przelaczone na te galaz |
| **Wykonawca** | sesja designu |
| **Wykaz plikow** | `design/05-okna/przeplyw/centrum-dowodzenia.html`, `design/zasoby/okna/centrum-dowodzenia.css`, `design/zasoby/okna/centrum-dowodzenia.js`, `design/zasoby/okna/danaco-anim-3d.css` |
| **Warstwa wspolna** | **poza terenem** — zmiany w `rama.css` i `panel-sesji.css` wracaja zgloszeniem, patrz nizej |
| **Dokumentacja terenu** | `design/01-dokumentacja-md/centrum-dowodzenia-uzasadnienia.md` — dopisany do wykazu po fakcie |
| **Poza terenem** | `budowa/`, `docs/`, `prowadzenie/`, pliki `design/zasoby/` **inne niz wymienione wyzej** |

**Przedmiot.** Doprowadzenie okna centrum dowodzenia do postaci przyjetej przez
Wlasciciela. Kompozycji nie przekazuje sie zleceniem — etap 1 prowadzi Wlasciciel,
a wykonawca odpowiada na jego poprawki.

**Zmiana warstwy wspolnej wymaga odnotowania.** Plan etapow stanowi: „Warstwa
o zasiegu ogolnym — zetony, fundament, komponenty, rama, stanowisko, karty okna —
jest wspolna dla wszystkiego, co powstanie pozniej. **Zmiana w niej jest decyzja,
nie poprawka okna.**" Teren tknal `rama.css` i `panel-sesji.css`; skutek obejmuje
kazde okno, ktore z tej warstwy wyrosnie. Do raportu odbioru: **co i dlaczego**
zmienilo sie w obu plikach.

**Kryteria odbioru.** Pierwsze jest juz zapisane w planie etapow jako jedyne
kryterium etapu 1 sprawdzalne maszynowo; pozostale z niego wynikaja.

1. **Okno nie zawiera wartosci projektowych wpisanych liczbowo** — barwa, odstep
   i kroj pochodza z zetonow. Wykazane przeszukaniem, **z sonda dodatnia**
   dowodzaca, ze wzorzec cokolwiek lapie.
2. **Progi ukladu pochodza ze skali systemu** (`--dn-bp-w1` 640, `--dn-bp-w2` 960,
   `--dn-bp-w3` 1280) albo odstepstwo jest nazwane i uzasadnione w raporcie.
   Zmierzone dzis: `centrum-dowodzenia.css:184,188` niesie 860 px i 1180 px spoza
   skali.
3. **Zeton uzyty zgodnie ze swoim przeznaczeniem.** Zmierzone dzis:
   `centrum-dowodzenia.css:38` mierzy kreske karty biezacej zetonem sily hasla
   (`--dn-wym-sila`) zamiast wstegi aktywnosci karty (`--dn-wym-wstega`).
3a. **Barwa z palety, nie z surowych udzialow.** Zmierzone dzis:
   `centrum-dowodzenia.css:26` sklada obrys marki z `46%`/`13%` w `color-mix`
   poza paleta `zasoby/zetony/`.
4. **Rozstrzygniecie o zasiegu ogolnym mieszka w warstwie ogolnej.** Zmierzone
   dzis: `--cd-obrys-marki` — regula obowiazujaca kazde okno — zyje w pliku
   jednego okna zamiast w `zasoby/zetony/`.
5. **Granica modul–biblioteka zachowana.** Plik okna nie nadpisuje komponentow
   biblioteki (`.dn-narzedzia`, `.dn-obszar-panel`, `.dn-karta-widoku`,
   `.dn-nrz-btn`, `.dn-etykietka`) albo kazde nadpisanie ma podany powod.
6. **Zmiany w `rama.css` i `panel-sesji.css` opisane** — co, dlaczego i jaki
   skutek dla okien, ktore z tej warstwy wyrosna.
8. **Uklad mierzy to, co okno oglasza.** Zmierzone 27.08: platno deklaruje
   `container-type: inline-size` i wlasna nazwe **dwa razy**, a uklad przelacza
   **dziesiec zapytan `@media` i ani jedno `@container`**. Okno mierzy wiec
   szerokosc ekranu, choc samouczek zwezа samo platno — przy otwartym panelu
   uklad nie odpowie. Albo zapytania ida na `@container`, albo deklaracja
   platna znika. Wykazane przeszukaniem z sonda dodatnia.
9. **Dokument uzasadnien nie przeczy arkuszowi.** Zmierzone 27.08:
   `centrum-dowodzenia-uzasadnienia.md:13` twierdzi, ze punkty lamania mierza
   szerokosc platna przez `@container plotno` — arkusz tego nie robi. Dokument
   opisuje stan, ktorego nie ma.
10. **Dokument ma miejsce w serii albo odsylacz.** Stoi w
   `design/01-dokumentacja-md/` bez numeru w serii `01-`–`10-` i bez odsylacza
   z ktoregokolwiek z nich. Dokument, do ktorego nic nie prowadzi, nie zostanie
   przeczytany.
11. **Objetosc komentarza w granicach reguly.** Zmierzone 27.08:
   `centrum-dowodzenia.css` 24 882 znaki komentarza przy 1087 wierszach
   (dopuszczalne 271), `rama.css` 30 726 przy 1906 (dopuszczalne 476),
   `centrum-dowodzenia.js` 5 978 przy 492 (dopuszczalne 123). Komentarz opisuje
   **regule obowiazujaca**, nie przebieg poprawki ani stan poprzedni — ustroj
   zabrania kroniki w tresci.
12. **Wiersz znacznika daje sie przejrzec.** Zmierzone 27.08:
   `centrum-dowodzenia.html` niesie wiersze 23 252, 20 198 i 10 108 znakow bez
   lamania — nieczytelne w przegladzie i w roznicy rewizji.

7. Okno wczytuje sie **bez bledu konsoli** — wykazane uruchomieniem, z sonda
   dodatnia dowodzaca, ze odczyt konsoli lapie bledy.

**Uwaga o scaleniu.** Galaz wyrasta z `f8f05b1`, a `main` przeszedl od tamtej
pory dziesiec rewizji. Konflikt w `prowadzenie/rejestr-terenow.md` rozstrzyga sie
**na rzecz `main`** — wersja z galezi cofnelaby dorobek calej tury.


## Zgłoszenia oczekujące na teren

### Wielojęzyczność rozpoznania rozsypuje się na dwóch językach

Zgłoszenie dotyczy **gałęzi `teren/wpiecie-unpaper`, nie `main`** — na `main`
uchwyt nie przekazuje języków wcale, więc usterka przyjdzie dopiero ze scaleniem.

`adapter_modul_badania_lektura.go:99-101` skleja żądane języki w jeden łańcuch
`strings.Join(jezyki, "+")`, a normalizator `jezykRozpoznaniaDokumentu`
(`adapter_narzedzia_dokument_tekst.go:342-353`) rozpoznaje wyłącznie nazwy
pojedyncze i łańcucha po `+` nie rozcina. `Languages: ["pl","en"]` daje więc
`-l pl+en`. Maszyna niesie wyłącznie `eng`, `osd` i `pol` — zmierzone
w `/usr/share/tesseract-ocr/*/tessdata/`. Skutek: żądanie **poprawne wobec
kontraktu** kończy się cichą awarią rozpoznania.

Usterka ujawnia się dopiero przy dwóch językach; jeden przechodzi. Sprawdzian
wniesiony razem z polem używa wyłącznie form już poprawnych (`pol`, `eng`)
i dlatego jej nie łapie.

**Blokuje scalenie `teren/wpiecie-unpaper` do `main`.** Naprawa dotyka
`adapter_narzedzia_dokument_tekst.go`, który należy dziś do terenu
`obrobka-wstepna-w-porcie-dokumentow` — teren na tę usterkę otwiera się
dopiero po jego zamknięciu. Dwa tereny na jednym pliku to ta sama szkoda,
przed którą stoi ustrój.

### Sprawdzian pola `preprocess` mierzy atrapę, nie skutek

`adapter_modul_badania_lektura_test.go:29-32` (gałąź `teren/wpiecie-unpaper`)
fabrykuje różnicę wyniku wewnątrz atrapy portu: `if z.Preprocess != nil &&
*z.Preprocess { tekst = "tekst po obróbce wstępnej" }`. Sprawdzian mierzy
zatem dojście pola do granicy portu, a nie skutek obróbki wstępnej na obrazie.
Komentarz w wierszach 62-65 ogłasza, że „mierzy skutek" — wobec złożonego
rdzenia jest to nieprawda.

Miarę skutku na realnej drodze produkcyjnej wnosi teren
`obrobka-wstepna-w-porcie-dokumentow`. Po jego zamknięciu ten sprawdzian
albo znika, albo przestaje twierdzić o sobie rzecz nieprawdziwą.

### Teren `mermaid-za-zapora` nie wytworzył ani jednej rewizji

Przebieg z 28.08 zamknął się bez pracy: `git log main..teren/mermaid-za-zapora`
pusty, drzewo bez zmian, schowek pusty. Kontroler orzekł zwrot na wszystkich
trzech kryteriach. `design.diagram.render` nadal rysuje własnym płótnem
(`adapter_modul_design_wykresy.go:812`), a mermaid-cli stoi nietknięty.

Teren pozostaje otwarty i wraca do wykonania.


### Warstwa wspolna zmieniona bez wlasnego terenu

`design/zasoby/rama.css` (+13/-4) i `design/zasoby/panel-sesji.css` (+8/-2)
zostaly zmienione przy pracy nad jednym oknem. Plan etapow stanowi: warstwa
o zasiegu ogolnym „jest wspolna dla wszystkiego, co powstanie pozniej. **Zmiana
w niej jest decyzja, nie poprawka okna.**"

Skutek obejmuje **kazde okno**, ktore z tej warstwy wyrosnie — glif przycisku
okna 12→16 px i uniesienie przez `top`. Domkniecie wymaga wlasnego terenu
obejmujacego warstwe wspolna, z opisem, co i dlaczego sie zmienilo.


### Przesiew i os obrazu pobraly drugie kopie wag — 5,2 GB

Teren `nastawy-zdolnosci-i-mowy` zglosil to jako ryzyko. Prowadzacy zmierzyl,
ze **to juz sie stalo**:

| W pamieci podrecznej Huba | Na dysku |
|---|---|
| `models--BAAI--bge-reranker-v2-m3` **3,6 GB** | `/opt/danaco-modele/reranker` 2,2 GB |
| `models--openai--clip-vit-large-patch14` **1,6 GB** | `/opt/danaco-modele/clip` 1,6 GB |

**Przyczyna:** stale `katalogNiewskazany` w `internal/wiedza/ustawienia.go` sa
puste, a migracja 402 zalozyla wartosci domyslne **rowne im co do znaku** —
swiadomie, zeby nie stworzyc dwoch prawd. Pusta wartosc znaczy dla biblioteki
„uzyj pamieci podrecznej", wiec pobiera.

To ten sam stan, ktory dla osadzarki naprawila migracja 401, wskazujac
`/opt/danaco-modele/embedder`. Domkniecie wymaga terenu obejmujacego **zarazem**
stala w `internal/wiedza/` i wiersz katalogu w migracji — rozdzielenie ich
stworzyloby dwie prawdy.

**Domkniete 28.08.** Stala w `internal/wiedza/ustawienia.go` i migracja 404 zmienily sie
w jednej rewizji, a druga rewizja wciagnela oba klucze do petli `ustawienia()` w
`adapter_modul_wiedza.go` — bez tego migracja byla dla dzialajacego rdzenia bezskutkowa,
co wykazala kontrola. Sprawdzian odtwarzajacy wskazuje uszkodzony katalog wag nastawa
i zada odmowy zamiast cichego siegniecia po wagi wbudowane.
Po domknieciu: 5,2 GB w `~/.cache/huggingface` staje sie zbedne.

### Transkrypcja pyta siec o metadane przy kazdym przebiegu

`silnik.transkrybuj` (`budowa/pomocniki/transkrypcja/silnik.py`) wola
`WhisperModel(..., download_root=...)` **bez `local_files_only=True`**, choc
droga `--wersja` ten argument ma. Zmierzone: **13 polaczen poza DNS** do
huggingface.co przy zerowym pobraniu wag; z `HF_HUB_OFFLINE=1` przebieg jest
identyczny i bezsieciowy.

**Skutek:** maszyna bez sieci doklada opoznienie i ryzyko odmowy tam, gdzie
wszystko lezy na dysku.

**Domkniete 28.08.** `local_files_only=True` wchodzi takze do drogi transkrypcji.
Powstal `test_silnik.py` z atrapa biblioteki: pod `pytest` przechodzi, a po usunieciu
argumentu pada. Uwaga do drabiny weryfikacji: repozytorium nie ma skryptu drabiny, wiec
`pytest` nie biegnie nigdzie z urzedu — sprawdzian trzeba do niej wpiac, gdy powstanie.


### Warstwa `obowiazkowa-apt` wykazu zaleznosci niesie proze

`WarstwaZaleznosci` w `zaleznosci_wykaz_wydruk.go` traktuje `default` jako apt,
wiec do tej warstwy wpadaja podpowiedzi zdaniem: „srodowisko uruchomieniowe Javy
(default-jre) wraz z...", „pip install ruff" i trzy inne.

**Skutek zmierzony przez teren `arsenal-wdrozenia`:** rozbicie tego pola na
spacjach dawalo `apt-get install -y ... uruchomieniowe Javy (default-jre) ...`,
apt padal, a `set -e` zabijal przebieg **przed** warstwami Go, npm, snap i mowy.
Skrypt prowizjonowania serwera nie dochodzil do konca **od dawna**.

**Domkniete 28.08 przez prowadzenie, po dwoch zwrotach kontroli.** Rdzen dostal warstwe
`obowiazkowa-recznie`, a rozstrzyganie idzie po ksztalcie pola: do warstwy apt wchodzi
wylacznie pole zlozone z nazw pakietow dystrybucji oraz polecenie `pip install`, dla ktorego
skrypt ma osobna galaz. Pole, ktorego drugim czlonem jest `install`, jest poleceniem, nie
wykazem — to odroznia `cargo install typos-cli` od dwoch nazw pakietow.

Naprawa objela TRZY miejsca naraz, bo rozdzielenie ich dawalo regresje: rdzen, obie galezie
skryptu prowizjonowania (plan i postawienie) oraz odpis awaryjny w tym samym skrypcie, gdzie
szesc pozycji przeklasyfikowano. Zmierzone po naprawie: wiersz `apt-get install` niesie same
nazwy pakietow, a szesc pozycji recznych skrypt WYPISUJE — wczesniej proze gubil po cichu.
Sprawdziany odtwarzajace obie usterki padaja bez naprawy i przechodza z nia.

### `cargo install` klasyfikowany jako warsztat Go

Regula klasyfikacji pytala o podnapis `go install`, a ten stoi wewnatrz
`cargo install typos-cli`. **Domkniete 28.08 razem z pozycja wyzej** — dopasowanie bierze
przedrostek, a pozycja idzie do warstwy recznej, nie do apt.

### Pakiet `.deb` serwera nie niosl pomocnikow mowy — nieaktualne

**Zdjete 28.08 przegladem pomiarowym.** budowa/scripts/pakiet-serwera.sh:93-102 — blok kopiujący `pomocniki/` do drzewa pakietu istnieje: gdy `$BUDOWA/pomocniki` nie ma, skrypt woła `padnij "nie ma katalogu pomocników: ... — bez niego pakiet nie postawi mowy"` (funkcja `padnij` kończy `exit 1`, linie 24-27). Katalog `budowa/pomocniki` faktycznie istnieje w drzewie (zawiera `odprawa.mjs` i podkatalog `transkrypcja/`). Brak pomocników jest dziś odmową złożenia pakietu, dokładnie jak stwierdza zapis rejestru — usterka pierwotna (pakiet szedł bez pomocników po cichu) nie istnieje.

Rdzen szuka pomocnikow obok siebie (potwierdzone `--wykaz-mowy`, klucz
`rozpoznanie.pomocnik-szukano`), a `pakiet-serwera.sh` katalogu `pomocniki/`
do pakietu nie wkladal. **Warstwy mowy nie dalo sie postawic na serwerze nigdy.**
Teren to naprawil: brak `pomocniki/` jest teraz odmowa zlozenia pakietu.


### Dokumentacja zetonow rozjechana z arkuszem o caly stopien — nieaktualne

**Zdjete 28.08 przegladem pomiarowym.** grep -n -- "--dn-fs-xs|...|--dn-fs-xl" na design/zasoby/zetony/zetony.css:111-116 i design/01-dokumentacja-md/04-tokens.md:360-365 — obie strony niosą identyczny ciąg 12/13/14/15/17/21 px, żadna z sześciu par nie różni się już o 1 px. Rozjazd nazwany zgłoszeniem nie istnieje w bieżącym stanie plików.

Audyt wskazal jedna pozycje. Pomiar Prowadzacego pokazal, ze rozjazd jest
**systematyczny**: z dziewieciu zetonow rozmiaru pisma **szesc** rozni sie
miedzy `design/zasoby/zetony/zetony.css` a `design/01-dokumentacja-md/04-tokens.md`,
i wszystkie o **dokladnie 1 px** — arkusz jest wiekszy.

| Zeton | Arkusz | Dokumentacja |
|---|---|---|
| `--dn-fs-xs` | 12 px | 11 px |
| `--dn-fs-sm` | 13 px | 12 px |
| `--dn-fs-base` | 14 px | 13 px |
| `--dn-fs-md` | 15 px | 14 px |
| `--dn-fs-lg` | 17 px | 16 px |
| `--dn-fs-xl` | 21 px | 20 px |

To nie jest literowka w jednym wierszu, tylko **cala skala podniesiona o stopien**
w arkuszu i nieodnotowana w dokumencie. Skutek: kazdy, kto siegnie po wartosc do
`04-tokens.md` — sesja, audyt, wykonawca — dostanie liczbe o 1 px mniejsza od
prawdziwej i uzna zgodne uzycie zetonu za usterke.

**Obowiazuje arkusz.** `zasoby/zetony/zetony.css` jest zrodlem prawdy o wartosci
zetonu; dokumentacja jest jego opisem i to ona wymaga wyrownania. Do sprawdzenia
przy tej okazji, czy tak samo nie rozjechaly sie odstepy, promienie i cienie.


### Sprawdziany zdolnosci modelowych pomijaja sie bez `DANACO_MODELE`

Dwa sprawdziany, ktore **jako jedyne dowodza, ze przesiew i os obrazu dzialaja**,
pomijaja sie, gdy zmienna `DANACO_MODELE` nie wskazuje katalogu wag. Powod jest
uczciwy — bez 4 GB wag nie ma czym liczyc — a pominiecie nazywa wprost, co
ustawic.

**Skutek jest jednak dokladnie ta klasa bledu, przed ktora ostrzega ustroj:**
domyslny bieg konczy sie zielono, nie mierzac zdolnosci, ktore ta tura wniosla.
Zmierzone: bez zmiennej 2120 zdanych i **19 pominietych**; ze zmienna oba
sprawdziany przechodza (32,2 s i 19,8 s).

**Obowiazujaca postac polecenia** dla biegu, ktory ma cokolwiek orzec o tych
zdolnosciach:

```
DANACO_MODELE=/opt/danaco-modele gotestsum -- -count=1 -timeout 40m ./...
```

**Rozstrzygniete 28.08 przez prowadzenie: sprawdzian bierze wagi z katalogu
wdrozeniowego, gdy zmienna milczy.** Ustawianie zmiennej z urzedu ukrywaloby
zaleznosc, a wywracanie biegu psuloby maszyny, ktore wag nie maja. Gdy zmienna
nie wskazuje niczego, a `/opt/danaco-modele` istnieje, sprawdzian bierze wagi
stamtad; pominiecie zostaje wylacznie dla maszyny bez wag i nazywa oba braki.
Zmierzone po zmianie: bieg BEZ zmiennej trwa 86,5 s i daje ZERO pominiec —
wczesniej pomijal sie natychmiast.


### Droga wejscia nie ma dokad prowadzic

Po domknieciu przygotowania okno zostaje na ekranie **na zawsze**. `main.ts`
montuje wylacznie trzy okna wejscia, a okno pracy — Centrum dowodzenia — istnieje
jedynie jako prototyp w `design/05-okna/przeplyw/centrum-dowodzenia.html`.

Wykonawca **nie dopowiedzial** przejscia ani zdania o nim: okno konczy na stanie
prawdziwym (100%, komplet etapow zamkniety), a droga naprzod powstanie razem
z oknem, do ktorego ma prowadzic. To jest fala 3 i czeka na prototyp.

### Trzy etapy przygotowania wroca, gdy dostana komendy

Wykaz etapow skrocono z pieciu do dwoch — do tych, ktore dzialajacy uklad mierzy.
„Profil Operatora i uprawnienia", „Kanaly modeli i konektory" oraz „Magistrala
kontekstu i pamiec projektow" **nie byly etapami czekajacymi, tylko obietnica,
ktorej nikt nie wykona** — kontrakt nie ma dla nich komend.

Podstawa: pozycja 12 rejestru decyzji — prototyp nie rozstrzyga, „ile jest czego
— plikow, skladnikow, etapow", bo to mierzy dzialajacy uklad. Etapy wracaja do
wykazu w dniu, w ktorym kontrakt dostanie dla nich komendy.


### Granica uprzezy nie wystarcza komendom neuronowym

Teren `usterki-rdzenia` podniosl granice uprzezy zgodnosci kontraktu z 15 s do
60 s, wyprowadzajac ja z granicy warstwy skanera. Teren `odtwarzanie-twarzy`
zmierzyl, ze **problem jest szerszy**: `image.upscale` na procesorze potrzebuje
okolo **90 s bez twarzy i 190 s z twarzami**, a przesiew wyszukiwania wczytuje
wagi przez okolo 30 s.

**Skutek:** sprawdzianu skutku komendy neuronowej nie da sie napisac droga
generycznej uprzezy — kazdy taki przebieg zostanie urwany i zamelduje usterke
rdzenia tam, gdzie po prostu liczyl model. Oba tereny obeszly to wlasna droga
wykonania z osobna granica, i oba zglosily to jako obejscie, nie rozwiazanie.

**Rozstrzygniete 28.08 przez prowadzenie: granica jednolita, podniesiona do
dziesieciu minut.** Przeslanka o biegu liczonym w godzinach okazala sie falszywa
i zostala obalona pomiarem. Uprzaz wola `context.WithTimeout` i czeka na
zakonczenie komendy, wiec granica jest ograniczeniem gornym, a nie czasem
oczekiwania: komenda szybka wraca natychmiast. Bieg `TestKazdaKomendaZnosiPustyLadunek`
po podniesieniu granicy trwa 10,9 s, czyli tyle samo co przed nim.

Granica zalezna od rodzaju komendy zostala odrzucona z drugiego powodu: kontrakt
nie niesie znacznika komendy liczacej modelem, wiec kazdy taki wykaz bylby
dopowiedziany, a tego zapora przed dryfem zabrania.

Skutek: prywatna stala `granicaKomendyZWagami` zniesiona, oba tereny moga zdjac
swoje obejscia, a sprawdzian skutku komendy neuronowej pisze sie zwykla droga
uprzezy.


### Pakiet serwera nie stawia jeszcze pomocnika twarzy — nieaktualne

**Zdjete 28.08 przegladem pomiarowym.** budowa/scripts/arsenal-serwera.sh zawiera pełną implementację warstwy TWARZE: `planTwarzy()` (linie ok. 337-350) opisuje kroki, `postawTwarze()` (linie 369-395) buduje środowisko pythonowe, instaluje torch/torchvision/facexlib, zapisuje opakowanie `/usr/local/bin/danaco-twarze` i woła `pobierzWagiTwarzy`. `grep -n postawTwarze` pokazuje wywołanie w `trybPostaw` na linii 636, obok `postawMowe` (630) i `postawWiedze` (633) — a `case "$TRYB" in ... postaw) trybPostaw "$WYKAZ" ;;` (linia 739) potwierdza, że to realna gałąź trybu `postaw`, nie tylko `plan`. Skrypt dziś stawia pomocnika twarzy tak samo jak realesrgan-ncnn-vulkan i rembg.

Rdzen wola `danaco-twarze`, a `scripts/arsenal-serwera.sh` stawia dzis wylacznie
`realesrgan-ncnn-vulkan` i `rembg`. Na maszynie Operatora `image.upscale`
z `faces: true` **odmowi mimo poprawnego kodu**.

Wykaz zaleznosci niesie pelna podpowiedz instalacyjna — nazwe, program, skladniki
srodowiska i katalog wag — wiec odmowa mowi, co dociagnac. Wpisanie tego do
skryptu arsenalu jest osobnym terenem.

### Kafelkowanie Real-ESRGAN zostawia szwy w powiekszonym obrazie

Zauwazone przez Prowadzacego przy ogladzie dowodu terenu `odtwarzanie-twarzy`:
w powiekszonym obrazie widac prostokatne granice kafelkow. Sa w **obu** wynikach,
z `faces: false` i `faces: true`, wiec pochodza z powiekszenia, nie z przebiegu
twarzowego. Nie zmierzone liczbowo, tylko dostrzezone okiem — do sprawdzenia,
czy da sie je zniesc nakladaniem kafelkow.

### CodeFormer stoi na dysku i nie jest wolany

`/opt/danaco-modele/twarze/codeformer.pth` zostaje niewykorzystany. Powod podany
przez wykonawce: CodeFormer stoi na wlasnej architekturze (VQGAN + transformer),
ktorej wagi nie niosa, wiec wymagalby drugiego zestawu cudzego kodu obok tego,
ktorym GFPGAN juz liczy. Do rozstrzygniecia, czy wart jest tej pracy.


### Os obrazu rozumie polszczyzne slabo

CLIP ViT-L/14, wskazany pozycja 17, jest **modelem jednojezycznym angielskim**.
Pomiar terenu `zdolnosc-wyszukiwania` na trzech obrazach:

| Zdanie | Wynik |
|---|---|
| `a red circle` | trafienie pierwsze, margines wyrazny (0,250 wobec 0,179) |
| `niebieski kwadrat` | trafienie pierwsze, **margines znika** (0,149 wobec 0,139) |
| `zolty trojkat na czarnym tle` | trafienie **przegrywa** |

**Skutek:** obietnica wyszukania obrazu zdaniem jest dzis slabsza dla Operatora
piszacego po polsku niz po angielsku, a produkt jest polskojezyczny. Nastawa
`wiedza_model_obrazu` pozwala podmienic model na wydanie wielojezyczne **bez
zmiany kodu** — wybor modelu nalezy do Wlasciciela, bo poda go pozycja 17.

### Wektory obrazow licza sie przy kazdym zapytaniu

Os obrazu nie ma trwalosci wektorow: kolumna `zakres` tabeli wskaznika ma warunek
na trzy wartosci, a migracja nalezy do innego terenu. Zmierzone: 17 s dla trzech
obrazow. Przy bibliotece rzedu setek zapytanie zajmie odpowiednio wiecej, a sufit
`GranicaObrazow` = 200 **obetnie wykaz bez ostrzezenia w odpowiedzi** — pole
`examined` powie tylko, ile weszlo.

### Cztery nastawy zdolnosci wyszukiwania bez wiersza w katalogu ustawien — nieaktualne

**Zdjete 28.08 przegladem pomiarowym.** budowa/server/internal/store/migracja_402_nastawy_przesiewu_i_obrazu.sql wstawia wiersze `definicja_ustawienia` dla wszystkich czterech kluczy (wiedza_model_przesiewu, wiedza_katalog_przesiewu, wiedza_model_obrazu, wiedza_katalog_obrazu), wraz z definicja_ustawienia_zasieg (globalny) i definicja_ustawienia_os (platform); migracja_404_katalogi_przesiewu_i_obrazu_stojace.sql doklada domyslne sciezki katalogow. Migracje wchodza automatycznie przez `//go:embed migracja_*.sql` (zrodlo_migracji.go:16-17), bez recznego wpiecia. git log potwierdza commit 61abba24 "Cztery nastawy przesiewu i obrazu dostaja wiersz w katalogu" oraz 4cfe99b2 "Wskazuje stojace wagi przesiewu i osi obrazu zamiast pustej sciezki". `config.set` ma dzis wiersze do sprawdzenia dla wszystkich czterech kluczy.

`wiedza_model_przesiewu`, `wiedza_katalog_przesiewu`, `wiedza_model_obrazu`,
`wiedza_katalog_obrazu` maja klucze i wartosci domyslne w `wiedza/ustawienia.go`,
ale nie maja wierszy `definicja_ustawienia`. Zdolnosc dziala, bo rozstrzyganie
czyta zapis niezaleznie od katalogu — ale `config.set` odmawia klucza spoza
katalogu, wiec **Operator nie ustawi katalogu wag z okna konfiguracji**.
Do domkniecia: cztery wiersze w migracji nastaw, wzorem czterech z migracji 115.


### Materialy wizerunkowe moga niesc ta sama obietnice szesciu postaci

Teren `witryna-pobierania` sprowadzil tresc stron do dwoch postaci hybrydowych,
ale `budowa/witryna/portfolio/` oraz material `DO-WGRANIA-DANACO-GROUP` wskazany
w `wydania.json` jako strona wizerunkowa lezaly poza jego wykazem plikow.
Do sprawdzenia, czy nie powtarzaja zniesionych wariantow.

### Rozpoznanie systemu po stronie klienta moze miec galaz linuksowa

`budowa/witryna/tresc/pobierz.mjs` powoluje sie w komentarzach na
`client/src/aktualizacja/wykaz-wydan.ts` jako na druga strone tej samej reguly.
Teren zdjal galaz linuksowa ze skryptu rozpoznania na stronie, ale pliku klienta
nie tknal - lezal poza wykazem. Jesli klient dalej rozpoznaje Linuksa, strona
i baner w aplikacji mowia dwie rozne rzeczy o tej samej maszynie.


### Cala rodzina komend `control.*` nie jest zmontowana

`nowyAdapterPrzejeciaSterowania` (`internal/core/adapter_przejecie_sterowania.go`)
nie ma w drzewie **ani jednego wywolania** - przeszukanie po nazwie dalo sama
definicje, przy sondzie dodatniej na `nowyAdapterRozmowy`, ktory trafia
w `montaz_rozmowa.go:53`.

**Skutek:** `control.takeover`, `control.release` i stan sterujacego wchodza do
rejestru komend droga `zarejestrujPrzejecieSterowaniaNiewpiete`, czyli **jako
odmowa**. Rejestr steru nigdy nie dostaje wpisu. To wyjasnia, dlaczego funkcja
`sterZlecenia` byla martwa, i dlaczego jej dokumentacja opisywala uklad, ktorego
kontrakt nie zna.

**Montaz wymaga pol, ktorych kontrakt nie ma.** `shared.LoopState` nie niesie
`controller`, `takenOverAt` ani `takenOverBy` - zmierzone odczytem
`shared/contract.go:13104-13123`. Wpiecie rodziny jest wiec rozstrzygnieciem
o zakresie produktu, nie praca inzynierska: wymaga pozycji rejestru decyzji
i dolozenia pol do kontraktu.

**Druga martwa funkcja w tym samym pliku:** `zachowajStery` wraz z metoda
`rejestrSteru.zachowaj`. Zostawiona swiadomie - to mechanizm sprzatania po
zamknietych oknach, wzorowany na czynnym `rejestrBiegow.Zachowaj`; usuniecie
zabraloby zabezpieczenie w chwili, gdy adapter zostanie wpiety.


### Wagi mowy stoja poza katalogiem modeli — nieaktualne

**Zdjete 28.08 przegladem pomiarowym.** Katalog /opt/danaco-modele/mowa istnieje (utworzony 27.08, 464 MB) i zawiera models--Systran--faster-whisper-small; `budowa/server/internal/mowa/ustawienia.go:38` ma dziś `KatalogModeliDomyslny = "/opt/danaco-modele/mowa"` (zgodnie z opisem w zgłoszeniu 'katalog nieistniejący'), a `migracja_075_mowa.sql` niesie wiersz definicja_ustawienia dla `mowa_katalog_modeli`. `~/.cache/huggingface/hub/` (ls) już NIE zawiera `models--Systran--faster-whisper-small` — wagi zostały fizycznie przeniesione, a kod wskazuje ten sam katalog. Domknięcie opisane w zgłoszeniu ('teren obejmujący przeniesienie wag i ustawienia.go') zaszło.

Teren `nastawy-wdrozenia` slusznie **nie zalozyl** trzeciej nastawy. Zgloszenie
podawalo `/opt/danaco-modele/mowa`, a takiego katalogu na maszynie nie ma -
sprawdzone przez wykonawce i potwierdzone przez Prowadzacego. **641 MB wag mowy**,
w tym `models--Systran--faster-whisper-small`, stoi w `~/.cache/huggingface`.

Nastawa wskazujaca katalog nieistniejacy sprowadzilaby silnik mowy do pobierania
od nowa, czyli **pogorszyla stan**. Domkniecie wymaga terenu obejmujacego zarazem
przeniesienie wag na dysku i `internal/mowa/ustawienia.go`, ktory lezal poza
wykazem tamtego terenu.

### Sprawdziany rdzenia sa na granicy domyslnego limitu czasu

Pakiet `internal/core` przekracza domyslne 10 minut `go test`: bieg bez
`-timeout` konczy sie zrzutem gorutyn po 600 s, nie wynikiem. Zmierzone:
`FAIL danacoconsole/server/internal/core 600.053s` przy 978 sprawdzianach.

**Skutek:** kto uruchomi sprawdziany dokladnie tak, jak podaje przekazanie,
dostanie niepowodzenie zamiast pomiaru. Obowiazujaca postac polecenia to odtad
`gotestsum -- -count=1 -timeout 40m ./...`. Rozbicie pakietu `core` albo
skrocenie najwolniejszych sprawdzianow jest osobna praca.

### Sciezka wag na wdrozeniu - obawa rozstrzygnieta pozycja 8 — nieaktualne

**Zdjete 28.08 przegladem pomiarowym.** Rozstrzygnięcie stoi i nic go nie podważa. `prowadzenie/decyzje.md` sekcja '## 8. Model wdrożenia: wyłącznie hybryda' (linie 670-721) ma `Stan: obowiązuje`, wprost stwierdza że rdzeń biegnie na serwerze Danaco a u Operatora tylko cienka powłoka Windows (pkt 1a, 2, 5). Kod: `budowa/server/internal/wiedza/ustawienia.go:53` ma `KatalogModeliDomyslny = "/opt/danaco-modele/embedder"`, a katalog `/opt/danaco-modele/embedder` fizycznie stoi na tej maszynie (serwerze) z wagami (pytorch_model.bin 2,27 GB, ls -la). Ścieżka linuksowa jest więc poprawną ścieżką serwera wdrożenia, zgodnie z zapisem 'obawa znika'.

Wykonawca zglosil, ze `/opt/danaco-modele/embedder` jako wartosc domyslna jest
wlasnoscia tej maszyny, a jedyna postacia produktu jest hybryda Windows 11.

**Obawa znika.** Pozycja 8 stanowi, ze w hybrydzie **rdzen stoi na serwerze
Danaco**, a u Operatora staje samo okno. Rdzen nigdy nie biegnie na Windowsie,
wiec sciezka linuksowa jest sciezka serwera wdrozenia i jest poprawna. Nastawa
pozostaje bez rozgalezienia po systemie.


### Strona pobierania obiecuje warianty zniesione pozycja 8 — nieaktualne

**Zdjete 28.08 przegladem pomiarowym.** Wskazane wiersze `budowa/witryna/tresc/strony.mjs:114,116,292,295,339` opisują dziś wyłącznie dwie postacie hybrydowe: instalkę Windows (x64/ARM64) i osobny pakiet serwera wdrożenia — np. linia 292: „Windows 10, Linux i macOS nie są obsługiwane”, linia 339: instrukcja instalacji ogranicza się do wyboru x64/ARM64. `grep -n 'AppImage|Linux|linux|natywn' budowa/witryna/tresc/pobierz.mjs budowa/witryna/tresc/strony.mjs` daje jedno trafienie — właśnie zdanie o braku wsparcia Linuksa/macOS, żadnej wzmianki o AppImage, sumie kontrolnej AppImage ani natywnej instalce Windows. Próba dodatnia wzorca (samo `grep -n` bez filtra) potwierdza, że polecenie działa i po prostu nic więcej nie łapie. Treść strony jest dziś zgodna z dwoma wariantami z pozycji 8, obietnica sześciu wariantów zniknęła.

Wykaz wydan sprowadzono do dwoch postaci hybrydowych, tresci strony wokol niego
nie. `budowa/witryna/tresc/strony.mjs:114,116,292,295,339` oraz
`tresc/pobierz.mjs:331,427,435,445,475` niosa rozdzial o wyborze miedzy hybryda
a postacia natywna, pakiety Linuksa, instrukcje sumy kontrolnej dla AppImage
i zdanie o braku instalki natywnej dla Windows. Dane wystawiaja dwie postaci,
a strona obiecuje szesc.

### Instalka wychodzi bez licencji i instrukcji

**Sprzecznosc zrodel — do rozstrzygniecia przez Wlasciciela, 28.08.**
Skrypt `pakiet-serwera.sh:106-115` zada czterech dokumentow produktu z korzenia
i przerywa budowe, gdy ktoregos brak. Wszystkie cztery **istnialy i zostaly
swiadomie usuniete** rewizja `d28990b8` jako opracowania poprzedniego podejscia:
README.md (593 wiersze), LICENSE.md (1198), INSTALACJA-I-KONFIGURACJA.md (641),
INSTRUKCJA-UZYTKOWANIA.md (796). Sa odzyskiwalne co do znaku z `d28990b8^`.

Skutek dzisiejszy: pakietu serwera nie da sie zlozyc w ogole, a instalatory
Windows wychodza bez licencji. Stopki dokumentacji odsylaja przy tym do pliku
`LICENSE`, ktorego w drzewie nie ma — odsylacz jest martwy.

Prowadzenie NIE przywraca ich samo. Tresc licencji jest instrumentem prawnym,
a przywrocenie cofneloby decyzje o usunieciu materialu poprzedniego podejscia.
Rozstrzygniecia wymaga jedno: czy dokumenty wracaja z historii, powstaja na nowo
dla obecnego produktu, czy skrypt przestaje ich zadac.

Skrypty instalek wymagaly z korzenia czterech dokumentow produktu i pakowaly je
do zasobow instalki. W korzeniu stoi dzis sam `CLAUDE.md`, wiec krok padal na
pierwszym pliku i teren `powloka-i-wydanie` go zdjal, zeby skrypty w ogole biegly.
Skutek: produkt koncowy nie niesie warunkow licencji ani instrukcji. Skrypty
wydania 1.0 przerywaly budowe przy ich braku, nazywajac to wprost: instalator bez
licencji i instrukcji nie jest produktem koncowym. Rozstrzygniecie nalezy do
Wlasciciela.

### Droga aktualizacji dla Linuksa zyje w powloce

`budowa/desktop/src-tauri/src/aktualizacja/droga.rs` niesie wariant AppImage,
sprawdzenie naglowka ELF i piec sprawdzianow dla platformy, ktorej pozycja 8 nie
przewiduje. Zostawione swiadomie: usuniecie zmienialo liczbe sprawdzianow, ktora
kryterium odbioru bralo za odniesienie.

### Katalog pakietu interfejsu ma dwie nazwy w drodze wdrozenia

`server/internal/konfiguracja/katalog_klienta.go` liczy domyslnie na `klient/dist`
i wprost ostrzega, ze dwie nazwy jednego katalogu wracaja odmowa przy pierwszym
uruchomieniu. Jednostka systemd w `budowa/packaging/` klade tymczasem sciezke
`client/dist`. Dziala wylacznie dlatego, ze jednostka podaje ja jawnie.
Ujednolicenie wymaga terenu obejmujacego `budowa/packaging/`.


Ustalenia z zamkniętych i biegnących terenów, które wykraczają poza ich zakres.
Każde zgłoszenie ma wskazany plik i wiersz. Zgłoszenie staje się terenem, gdy
Prowadzący je otworzy; do tego czasu jest wykazem, nie pracą.

### Ukończenie biegu może raz na jakiś czas skłamać — domknięcie poza terenem

Ustalenie terenu `warunek-ukonczenia-zadania`, wyprowadzone z odczytu dwóch
funkcji. `powodTury` (`adapter_rozmowa_petla.go:46`) nie widzi `zamkniecie.Blad`,
więc tura zamknięta zdarzeniem `result` z `is_error: true` przy sprawnym kanale
idzie do pętli jako `PowodWynik` i pętla ogłosi `completed`, choć wiadomość
dostaje stan `error`. Pętla używa jedynego sygnału, który dostaje, i nie dubluje
odczytu zamknięcia — drugi czytelnik byłby drugą prawdą. Domknięcie: jeden wiersz
w `adapter_rozmowa_petla.go` (przekazanie `zamkniecie` do `powodTury`) — plik
poza terenem, więc zgłoszenie.

Osobno, ta sama warstwa: pętla nie wie o **początku** tury wykonawcy (dowiaduje
się z pierwszego fragmentu), więc wykonawca zlecony tuż przed końcem tury
koordynatora jest przejściowo niewidoczny. Szczelina jest nieszkodliwa —
samopodjęcie biegu po ukończeniu i tak rusza kolejny obieg, stan kontrolki sam
się prostuje. Domknięcie do zera: `petla.ZTuraWBiegu(rozmowa.CzyTuraWBiegu)`,
jeden wiersz w `montaz_rozmowa.go` — plik poza terenem. Uwaga dla wpinającego:
w chwili `ZakonczTure` mapa `biegnace` wciąż zawiera okno koordynatora
(`zapomnijBieg` jest `defer`), więc pytać wolno tylko o wykonawców.

### Martwa funkcja `sterZlecenia` z komentarzem w nieistniejące miejsce — nieaktualne

**Zdjete 28.08 przegladem pomiarowym.** grep -rni "sterZlecenia" w calym drzewie trafia wylacznie w typ `SterZlecenia` (adapter_przejecie_sterowania.go:34) i jego czynne uzycia (a.ster.ustaw/.stan/.zdejmij wolane wprost z metod Przejmij/Oddaj/Ster tego samego pliku) — funkcji o nazwie `sterZlecenia` (mala litera) nie ma. Komentarz odsylajacy do `stan_obiegu.go` tez nie istnieje: grep "przekład licznika obiegów|stan_obiegu" w plikach core/*.go daje zero trafien. git log -- adapter_przejecie_sterowania.go pokazuje commit 2b52fa57 "Usuwam martwa funkcje steru zlecenia" (27.08.2026, autor Dariusz Naharnowicz) usuwajacy dokladnie ta funkcje i komentarz (6 wstawien, 16 usuniec) — teren `warunek-ukonczenia-zadania` zgloszil usterke poza swoim zakresem, a nastepna rewizja ja zdjela.

`core/adapter_przejecie_sterowania.go:146` — doc funkcji `sterZlecenia` mówi
„Woła ją przekład licznika obiegów w `core/stan_obiegu.go`", a `stan_obiegu.go`
jej nie woła i nikt inny też nie. Funkcja martwa, komentarz kieruje w nieistniejące
miejsce. Zauważone przy terenie `warunek-ukonczenia-zadania`, poza jego zakresem.

### Granica 15 s w uprzęży kontraktu jest ciasna dla warstwy skanera — nieaktualne

**Zdjete 28.08 przegladem pomiarowym.** Stala `granicaKomendySprawdzianu` w budowa/server/internal/core/zgodnosc_kontraktu_test.go:26 wynosi dzis `10 * time.Minute`, nie 15 s, i jest uzywana konsekwentnie: zgodnosc_kontraktu_test.go:192,231, skutek_zdolnosci_wyszukiwania_test.go:129 oraz blokady_skutek_test.go:103 (to sa dokladnie pliki wskazane w zgloszeniu). Komentarz przy stalej (linie 17-25) wprost tlumaczy podniesienie: granica obejmuje "najwolniejsza zmierzona komende liczaca modelem na procesorze ... okolo stu dziewiedziesieciu sekund". grep "15 \* time.Second" w internal/core/*.go trafia tylko w cztery NIEZWIAZANE stale (adapter_kanaly_sprawdzenie.go, adapter_modul_extension_protokol.go, przegladarka_pobieranie.go, skutek_terminala_test.go) — zadna nie jest generyczna uprzeza zgodnosci kontraktu. git log potwierdza commity 1a3c537a "Granica komendy w uprzezy wychodzi z granicy warstwy" i 9cdec8aa "Ujednolica granice czasu uprzezy i znosi obejscie komend z wagami". Nowa granica 10 min jest daleko powyzej ~14,8 s czynnosci skanera, wiec opisany konflikt nie wystepuje.

Ustalenie zmierzone przez teren `odmowy-skanera` i potwierdzone przez kontrolę.
Generyczna uprząż zgodności kontraktu woła każdą komendę rdzenia z twardym
limitem 15 s: `zgodnosc_kontraktu_test.go:215` i `:254`,
`blokady_skutek_test.go:127`. Tymczasem czynność skanera na tej maszynie
dochodzi do ~14,8 s (`scanimage --format=png` ~7,4 s + wykaz `scanimage -L`
~7,4 s), a własne limity warstwy są znacznie wyższe — `granicaWykazuUrzadzen`
45 s, `granicaSkanowaniaUrzadzenia` 5 min (`urzadzenia_skaner.go:63,68`). Pod
obciążeniem maszyny generyczne 15 s ucina czynność przed jej własnym limitem
i odmowa spada na kod czasu (`internal_error`) zamiast rozpoznać brak
urządzenia (`not_found`). Na maszynie budowlanej bez żywego skanera komendy
odmawiają natychmiast, więc chwiejność ujawnia się dopiero tam, gdzie warstwa
realnie sonduje urządzenie — u Operatora na wolnej maszynie ta sama granica
może uciąć legalną odmowę „brak urządzenia".

Teren obejmuje uprząż zgodności kontraktu (pliki poza terenami zaplecza).
Rozstrzygnięcie progu — podnieść do wartości spójnej z limitami warstwy —
należy rozważyć wraz z tym, jak uprząż ma traktować komendy o własnych,
dłuższych limitach.

### Rejestr rozjechany na gałęzi centrum poprawek — rozstrzygnięcie przy scaleniu

Gałąź `teren/centrum-poprawki` niesie dwie rewizje dotykające tego pliku —
`6361df8` i `31e37de` — założone tam omyłkowo przez poprzedniego Prowadzącego,
bo drzewo główne stoi przełączone na tę gałąź. Ich treść przeniesiono na `main`
jako `a52fa05` i `ea0c673`, ale na gałęzi zostały. Od tamtej pory `main`
przerobił ten plik mocno: pięć terenów zamkniętych, zgłoszenia i reguły odbioru.

**Przy scaleniu gałęzi designu konflikt w tym pliku rozstrzyga się na rzecz
`main`.** Wersja z gałęzi cofnęłaby zamknięcia terenów i zgłoszenia z pomiarów,
a wyglądałoby to na zwykłe scalenie. Praca w `design/` scala się normalnie —
rozstrzygnięcie dotyczy wyłącznie `prowadzenie/rejestr-terenow.md`.

### Kontrakt ruszył dwa razy od przejęcia — suma bieżąca `b7d0436880d878e78576`

Do 27.08.2026 kontrakt stał nietknięty pod sumą `2cbb843d33f4531b05cd`. Dwie
zmiany, obie wyłącznie dodania, obie przez porównanie strukturalne potwierdzone
jako niezmieniające niczego istniejącego:

| Teren | Co dołożył | Suma po |
|---|---|---|
| `pomiar-stron` | trzy komendy pomiaru strony (1077 → 1080) | `334705bd88c2efc13779` |
| `warunek-ukonczenia-zadania` | czwarta wartość `LoopStopReason` — `completed` | `b7d0436880d878e78576` |

Sprawdzian nietykalności kontraktu odnosi się do sumy **bieżącej**
`b7d0436880d878e78576`, nie do sumy z pierwszego dnia; kolejny teren, który
kontrakt dokłada, podaje w rejestrze sumę zastaną i sumę po sobie.

### Wdrożenie musi założyć trzy nastawy, inaczej stojące wagi leżą odłogiem

Ustalenie terenu `zaplecze-modeli`, zmierzone na żywym rdzeniu. Silnik wiedzy
działa na wagach z `/opt/danaco-modele/embedder`, ale **dopiero po nastawach** —
wartości domyślne rdzenia wskazują co innego i wdrożenie pobierze drugi model
zamiast użyć stojących 4,3 GB.

| Nastawa | Wartość | Bez niej |
|---|---|---|
| `wiedza_model` | `BAAI/bge-m3` | rdzeń sięga po `mpnet` z migracji 115 |
| `wiedza_katalog_modeli` | `/opt/danaco-modele/embedder` | wagi pobierane na nowo do katalogu danych |
| `mowa_katalog_modeli` | `/opt/danaco-modele/mowa` | wagi mowy stoją w pamięci podręcznej konta, które uruchomiło rdzeń |

Wag mowy dotyczy osobne ustalenie: 464 MB modelu `faster-whisper-small` stoi
dziś w pamięci podręcznej pod katalogiem domowym, bo `mowa/ustawienia.go` przy
pustej nastawie zostawia miejsce bibliotece, podczas gdy `wiedza/pomocnik.go`
przy pustej nastawie **przypina** wagi do katalogu danych rdzenia — droga
wołania procesu nie dziedziczy środowiska, więc pomocnik nie zna nawet `HOME`.
Skutek: czyszczenie pamięci podręcznej kasuje działającą funkcję bez śladu
w produkcie, a wagi są przywiązane do konta uruchamiającego. Rozstrzygnięcie
przyjęte: przenieść wagi do `/opt/danaco-modele/mowa` i wskazać je nastawą.

**Domkniete 28.08 — sprawdzone, ze wszystkie trzy nastawy sa juz zalozone.**
Zgloszenie zdezaktualizowaly migracje 401 i 403 wraz z odpowiadajacymi im stalymi:

| Nastawa | Wartosc domyslna | Stala Go | Migracja |
|---|---|---|---|
| `wiedza_model` | `BAAI/bge-m3` | `wiedza.ModelDomyslny` | 401 |
| `wiedza_katalog_modeli` | `/opt/danaco-modele/embedder` | `wiedza.KatalogModeliDomyslny` | 401 |
| `mowa_katalog_modeli` | `/opt/danaco-modele/mowa` | `mowa.KatalogModeliDomyslny` | 403 |

Kazda para stala–migracja zgadza sie co do znaku, wiec dwoch prawd nie ma.
Pytanie o zrownanie rodzin rozstrzygnelo sie inaczej, niz je postawiono: pustej
nastawy mowy nie trzeba przedefiniowywac, bo wartosc domyslna przestala byc pusta.
Puste znaczy dalej „pamiec podreczna biblioteki" i jest to wybor Operatora, a nie
stan zastany po wdrozeniu.

### Trzy modele bez komendy w kontrakcie

Ustalenie terenu `zaplecze-modeli`. Wagi stoją, komend nie ma — rejestr zabronił
ich dokładania, bo kontrakt należał w tej turze do innego terenu.

| Model | Waga | Propozycja obszaru |
|---|---|---|
| reranker (bge-reranker-v2-m3) | 2,2 GB | `knowledge` — przełącznik `rerank` w `knowledge.search` albo osobna `knowledge.rerank` |
| CLIP | 1,6 GB | `library` albo `knowledge` — wyszukiwanie obrazów po znaczeniu |
| twarze (GFPGAN, codeformer) | 692 MB | pole `faces` w `image.upscale` **już istnieje**; brakuje silnika: wydanie ncnn nie niesie sieci twarzowej, a wagi `.pth` żądają stosu torch, którego rdzeń nie woła. Odmowa nazywa dziś brak `gfpgan-ncnn-vulkan` |

Osobno rozstrzygnięte i zamknięte: **Kokoro nie jest wart deklaracji** — 54 głosy,
żadnego polskiego (`pf`/`pm` to portugalski), a piper z `pl_PL-darkman-medium`
działa i rdzeń go zna.

### Trzy narzędzia treści pisanej bez legalnego miejsca wpięcia

Ustalenia terenu `dokumenty-i-tekst`, każde zmierzone, żadne nienaprawialne
w granicach tamtego terenu.

| Rzecz | Zmierzone | Skutek |
|---|---|---|
| LanguageTool przy `translate.quality.check` | tryb `--bitext` (LT 6.6) na parze z jaskrawą rozbieżnością liczb (10:30→9:30, 250→500 EUR) oddaje **zero ustaleń**, a wyjścia maszynowego dla bitext program nie ma wcale | wpięcie dałoby instrument milczący przy realnej usterce; komenda mierzy wierność panelu wobec źródła i dziś nie ma czym |
| unpaper przy `research.source.ocr` | uchwyt `adapter_modul_badania_lektura.go:589` upuszcza pola `preprocess` i `languages`, a `DocumentTextExtractRequest` kontraktu tych pól nie ma | obróbka wstępna nieosiągalna z drogi badań; pole `preprocess` pozostaje martwe |
| Tika przy `library.metadata.get` | obsługa stoi w `adapter_modul_library_technika.go` — pliku terenu `obraz-i-diagramy`, gdzie exiftool wszedł w tę samą komendę | rozłączność terenów; do rozważenia, czy Tika ma tam co dołożyć wobec exiftoola |

### mermaid-cli nie ma legalnego miejsca wpięcia

Ustalenie terenu `obraz-i-diagramy`, zmierzone uruchomieniem. Obie komendy
z tabeli narzędzi są zamknięte: `design.diagram.render`
(`adapter_modul_design_wykresy.go:790`) leży za zaporą fotografii, która
zabrania **każdego** wołania procesu w plikach `adapter_modul_design*.go` —
plik próbny z samą wzmianką `zewnetrzne.Wolaj` wywrócił zaporę, choć zapora
nazw przepuszcza `mmdc`; nadto kontrakt tej komendy przyjmuje węzły i krawędzie,
nie ma pola na źródło Mermaid. `apps.architecture.export`
(`adapter_modul_aplikacje_architektura.go:427`) niesie w nagłówku
rozstrzygnięcie „żaden format nie woła programu z zewnątrz, więc eksport działa
na instalce, która niesie sam rdzeń".

Samo narzędzie działa: z przeglądarką stojącą w `/opt/ms-playwright` renderuje
źródło w kształcie eksportu rdzenia (SVG 12 539 B, PNG 265×278). Bez wskazania
przeglądarki `mmdc` odmawia — puppeteer żąda własnej kopii.

Do rozstrzygnięcia: czy rysowanie diagramu ma powstać poza obszarem Design
(nowa komenda i nowe miejsce), czy zostaje niezrobione. Rozstrzygnięcie dotyka
zapory niosącej rozstrzygnięcie Właściciela, więc nie jest samą robotą.

### Komunikat rdzenia radzi budowę usuniętymi skryptami — nieaktualne

**Zdjęte 28.08 po sprawdzeniu.** Przeszukanie `wpiecie.go` nie znajduje odwołań
do `scripts/wydanie.sh` ani `scripts/pakowanie.sh`; usterkę naprawiono wcześniej,
a zapis ją przeżył.

### Skrypty hybrydowe wskazują usunięte profile i nieistniejący katalog — nieaktualne

**Zdjęte 28.08 po sprawdzeniu, każdy człon osobno.** Żaden skrypt nie wskazuje
profili `tauri.hybryda-*`. Skryptu `instalka-hybryda-linux.sh` nie ma, więc
sprzeczność z jedyną platformą Windows 11 rozstrzygnęła się sama; zostały dwa
skrypty Windows, x64 i ARM64.

Ścieżka `client` w `pakiet-serwera.sh` **nie jest literówką**: to układ wdrożenia
wewnątrz pakietu, wskazany jednostce systemd zmienną
`DANACO_KATALOG_KLIENTA=/opt/danaco-console/client/dist`, którą rdzeń czyta.
Wartość `klient/dist` w `konfiguracja/katalog_klienta.go` jest wyłącznie domyślną
dla uruchomienia z drzewa budowy.

Oba skrypty przeszły 28.08 przelotem do końca i złożyły instalatory obu
architektur — zapis o odmowie na starcie był nieprawdziwy.

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

### Droga wejścia — braki po stronie rdzenia, ujawnione pomiarem

Wszystkie pochodzą z terenu `droga-wejscia` i są zmierzone na żywym rdzeniu.
Żaden nie wstrzymał pracy; okno podłącza je dziś do odmowy nazywającej brak.

| Rzecz | Zmierzone | Skutek |
|---|---|---|
| brak drogi powtórnego potwierdzenia adresu | `auth.verify` z cudzą drogą → `not_authenticated`; powtórny `auth.register` → `conflict`; żadna z dziewięciu komend `auth.*` nie wydaje drogi drugi raz | Operator, do którego list nie dotarł, **nie ma wyjścia z okna** — dotyczy także instalacji **z pocztą**, nie tylko bez |
| trzy z pięciu etapów przygotowania środowiska bez komendy | prototyp wymienia pięć; droga wejścia obsługuje dwa. „Profil i uprawnienia", „Kanały modeli", „Magistrala kontekstu" nie mają przypisanej komendy | pasek postępu zatrzymuje się na 40% — wartość prawdziwa, nie ozdobna |
| odmowa logowania nie niesie długości zwłoki | osiem prób: 100, 350, 600, 1100, 2100, 4100, 5100, 5100 ms; `details` puste za każdym razem | okno mierzy zwłokę samo, czasem trwania próby poprzedniej |
| magazyn tokenu bramki między uruchomieniami | `connection.hello` przyjmuje token i oddaje `authenticated: true`; gdzie token mieszka, nie rozstrzyga ani prototyp, ani kontrakt, ani rejestr | odsłona „rozpoznano zaufane urządzenie" działa dopiero po wskazaniu magazynu — należy do powłoki |

### Adres serwera wdrożenia nie stoi w żadnym źródle

Teren `powloka-tauri` przeszukał `budowa/`, `design/`, `docs/` i `prowadzenie/`:
pozycja 8 nazwy serwera nie podaje, prototyp instalatora nie ma kroku adresu,
`wydania.json` niesie wyłącznie kanał pobrań. Powłoka bierze adres ze wskazania
Operatora albo ze zmiennej `DANACO_HOST_RDZENIA`. Czy adres ma być wkompilowany
przy składaniu instalki — i skąd wtedy pochodzi — jest rozstrzygnięciem
Właściciela.

### Zastane usterki powłoki, ujawnione przy przejęciu

| Rzecz | Stan |
|---|---|
| `aktualizacja/probne.rs` nie istnieje | `droga.rs:293` importuje `KatalogProbny` w siedmiu sprawdzianach; `cargo test` **nie kompiluje się od chwili przejęcia** |
| `wykonaj_aktualizacje` pobiera spod dowolnego adresu | `pobranie.rs` sprawdza wyłącznie sumę SHA-256, też podaną przez stronę; brak przywiązania do kanału z `wydania.json` |
| cztery z sześciu poleceń IPC bez odbiorcy | `stan_rdzenia`, `wskazanie_rdzenia`, `wybierz_katalog_roboczy`, `wykonaj_aktualizacje` — opisywały odbiorców w usuniętym `client/src/powloka/` |
| `app.security.csp: null` | strona bez polityki treści; ułożenie wymaga interfejsu, który dopiero powstaje |
| sześć skryptów w `budowa/scripts/` | wskazują usunięte profile nastaw, nieistniejący `budowa/client/` i wariant natywny zniesiony pozycją 8; pozycja 8 pkt 4 mówi, że mają leżeć **poza** repozytorium |

### Łańcuchy widoczne poza katalogiem treści w paczce designu

Właściwość, którą pozycja 12 nazywa wiążącą, jest w dwóch miejscach naruszona:
`design/zasoby/okna/wejscie/skladniki/pole-hasla.js:55` niesie zaszyte
`'Pokaż hasło'`, a `design/zasoby/okna/przeplyw-wejscia.js` powiela w kodzie
`SILA_PUSTE`, `SILA_OPISY` oraz dwa napisy o schowku, które stoją już
w katalogu. W kliencie wszystkie te napisy siedzą w katalogu treści.

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

### Droga SANE bez odpowiednika `bladWarstwyWia` — gotowe do otwarcia

Brak urządzenia jest już nazwany. Brak samego programu `scanimage` dalej wychodzi
odmową arsenału bez wskazania drogi obejścia, podczas gdy `bladWarstwyWia`
(`urzadzenia_skaner.go`) dla tej samej sytuacji na Windowsie podaje
`studio.ingest.queue.add`. Ta sama asymetria dotyczy `wykazSkanerow`: gałąź
Windows przekłada odmowę, gałąź Linux oddaje ją surową. Operator na Linuksie bez
`sane-utils` nie dowie się, że materiał da się wnieść inną drogą.

### Sprawdzian katalogu akcji szuka nieistniejącego katalogu — naprawione

`katalog_akcji_test.go` szukał `../../../client/src/ikony/zrodla`, którego nie
ma w żadnej gałęzi. Rewizja `9ba1f90` wiąże sprawdzian ze źródłami ikon nowego
klienta — `sciezkaZrodelIkon` wskazuje `../../../klient/src/ikony/zrodla`.
Przechodzi w biegu odniesienia 2041 zdanych, zero niezdanych.

### Reguła odbioru wyprowadzona z pomiarów

`axe.run()` sam wywołuje dwa błędy 404 (`menu.css`, `ruch.css`), bo rozwiązuje
`@import` względem adresu dokumentu, a nie arkusza. Konsola przed wstrzyknięciem
axe jest pusta. Tych dwóch wpisów nie liczy się jako brudnej konsoli.

**Pełny bieg wymaga podniesionego limitu czasu.** Pakiet `internal/core`
przekroczył domyślną granicę `go test` (10 minut) po dołożeniu przez cztery
tereny około czterdziestu sprawdzianów wołających prawdziwe programy zewnętrzne.
Bieg pada wtedy paniką „test timed out after 10m0s" wskazującą **przypadkowy**
sprawdzian, który akurat biegł — wygląda to na usterkę tego sprawdzianu i nią
nie jest. Objaw uboczny myli podwójnie: `gotestsum` naliczy wtedy około 947
sprawdzianów zamiast ponad dwóch tysięcy, bo pakiet, który spanikował, nie
policzy swoich. Pełny bieg uruchamia się z `-timeout 30m`.

**Pełne biegi szereguje się zamkiem.** Maszyna ma 16 rdzeni; przy czterech
terenach naraz biegło na niej dziewięć procesów sprawdzianów i sprawdziany
z twardymi granicami czasu zaczęły się chwiać — każdy pełny bieg gubił **inny**,
a pojedynczo wszystkie przechodziły. Granice tego rodzaju stoją m.in.
w `zgodnosc_kontraktu_test.go`. Bieg pełny wykonuje się pod zamkiem:

```bash
cd budowa/server && flock /tmp/danaco-bieg-pelny.lock \
  gotestsum -- -count=1 -timeout 30m ./...
```

Sprawdzian, który padł w pełnym biegu, a przechodzi uruchomiony pojedynczo, jest
chwiejnością pod obciążeniem, nie usterką — i tak się go nazywa.

**Pełny bieg regresji jest bramką scalenia u Prowadzącego, nie bramką wyjścia
u wykonawcy.** Wykonawca dowodzi swojego przedmiotu biegiem **celowanym** pakietu,
który dotyka — jest szybki i odporny na obciążenie. Pełny bieg `./...` w drzewie
terenu mierzy drzewo **bez** pozostałych terenów biegnących równolegle, więc
rozstrzyga wyłącznie bieg na scalonym `main` po scaleniu, a ten wykonuje
Prowadzący raz, pod zamkiem, jako warunek bramki scalenia. Zaprzęganie czterech
terenów do czterech pełnych biegów naraz zapycha zamek na godziny i nic nie
rozstrzyga — każdy z nich mierzy inny, niepełny stan.

Czas dostępu do pliku nie dowodzi, że pliku nie czytano. Drzewo stoi na `ext4`
zamontowanym z `relatime`, gdzie jądro odświeża czas dostępu wyłącznie wtedy,
gdy poprzedni jest starszy od czasu zmiany albo starszy niż doba. Odczyt pliku,
którego czas dostępu jest już późniejszy od czasu zmiany, **nie zostawia
śladu** — niezmieniony czas dostępu jest tam brakiem pomiaru, nie dowodem
nietknięcia. Wyszło to przy sporze o sprawstwo pobrania wag mowy: wykonawca
podał czas dostępu jako dowód, a instrument z założenia milczał.

Pomiar w przeglądarce wymaga jawnego ustawienia `PLAYWRIGHT_BROWSERS_PATH` na
`/opt/ms-playwright` w poleceniu, a nie polegania na środowisku powłoki — powłoka
uruchomiona przed ustawieniem zmiennej jej nie widzi i pobiera przeglądarki
po raz drugi.

## Tereny zamknięte

### mermaid-za-zapora — ZAMKNIETY JAKO BEZPRZEDMIOTOWY

**Rozstrzygniecie z 28.08.2026, oparte na pomiarze kontraktu.** Przeslanka terenu byla
falszywa. Teren zadal wpiecia mermaid-cli w `design.diagram.render`, tymczasem:

- `design.diagram.render` bierze **wezly i polaczenia** (`nodes`, `edges`), a nie tekst
  w notacji mermaid. Wyrysowanie ich wlasnym plotnem w Go jest realizacja poprawna,
  nie brakiem. `mmdc` przyjmuje zrodlo mermaid — czego ta komenda nie niesie.
- Kontrakt wymienia `mermaid` **dokladnie raz**: jako wartosc `AppExportFormat` komendy
  `app.architecture.export`. Te notacje rdzen juz wytwarza czystym Go
  (`adapter_modul_aplikacje_architektura.go`), bez wolania procesu.
- `mmdc` stoi na maszynie, ale produkt go nie potrzebuje. Program obecny nie jest
  zobowiazaniem do jego uzycia.

Dwa przebiegi wykonawcze przepadly na tym terenie, bo zadanie nie mialo legalnego
rozwiazania. Drugi wytworzyl zmiane SZKODLIWA — dopisal `"mermaid"` do wykazu
`silnikiSpozaInstalki`, ktory wylicza nazwy PROGRAMOW, przez co zabronil obszarowi
Design wymieniania formatu, ktory kontrakt dopuszcza. **Rewizja `c10f25a9` na galezi
`teren/mermaid-za-zapora` NIE WCHODZI do `main`.** Galaz ginie.

Zapora fotografii pozostaje nietknieta i nadal broni tego, co miala bronic.


| Nazwa | Gałąź | Rewizje | Kontrola |
|---|---|---|---|
| `nastawy-zdolnosci-i-mowy` | `teren/nastawy-zdolnosci-i-mowy` | `61abba2`, `97c6fc2` | weryfikacja Prowadzacego wlasnym pomiarem: cztery klucze przyjmowane przez `config.set`, zastanych migracji **nie ruszono** (same nowe 402 i 403), wagi mowy przeniesione — 464 MB w `/opt/danaco-modele/mowa`, `model.bin` rozwiazuje sie; sonda dodatnia wykonawcy: **486 MB pobrania przy pustym katalogu wobec zera przy pelnym** |
| `arsenal-wdrozenia` | `teren/arsenal-wdrozenia` | `541cdce` | weryfikacja Prowadzacego wlasnym pomiarem: odpis wykazu wyciagniety ze skryptu ma **55 pozycji, zgodnych z rdzeniem** co do warstwy, programu, pakietu i nazwy (bylo 30); `bash -n` i `shellcheck` czysto; zakres 2 pliki |
| `zdolnosc-wyszukiwania` | `teren/zdolnosc-wyszukiwania` | `5d02014` | weryfikacja Prowadzacego wlasnym pomiarem: kontrakt ruszony **wylacznie dodaniami** — stare pola `knowledge.search` sa prefiksem nowych, opis bez zmiany, zero usuniec; generator powtarzalny w dwoch przebiegach; oba sprawdziany zdolnosci **uruchomione z `DANACO_MODELE` i zdane** — przesiew zmienil kolejnosc, os obrazu trafila w kolo |
| `okno-przygotowania` | `teren/okno-przygotowania` | `fbb405d` | weryfikacja Prowadzacego wlasnym pomiarem: martwy przycisk **zniknal** (0 trafien przy 2 kontrolnych na nowa czynnosc), 27/27 sprawdzianow przebiegu, 7/7 katalogu tresci, `tsc` bez bledu; zrzuty obejrzane — postep dobiega 100%, odslona nieudana ma droge naprzod |
| `odtwarzanie-twarzy` | `teren/odtwarzanie-twarzy` | `9cc5ccf` | weryfikacja Prowadzacego wlasnym pomiarem: **RMSE 1301,37** miedzy wynikiem `faces:false` a `faces:true`, roznica zlokalizowana na twarzy; wycinek obejrzany — zeby, wargi i faktura skory wyraznie odtworzone; zakres 5 plikow w `internal/core`; kontrakt nietkniety; wagi `GFPGANv1.4.pth` wczytane `strict=True`, 285 kluczy |
| `witryna-pobierania` | `teren/witryna-pobierania` | `6ffe74f` | weryfikacja Prowadzacego wlasnym pomiarem: w tresci stron zostalo **jedno trafienie** wzorca zniesionych postaci i jest nim zdanie odmawiajace wprost z pozycji 8; sonda dodatnia tego samego wzorca daje 4 trafienia w `wydania.json` i 0 w `zloz.mjs`, wiec rozroznia; witryna sklada sie - 10 stron |
| `usterki-rdzenia` | `teren/usterki-rdzenia` | trzy rewizje | weryfikacja Prowadzacego pomiarem: sprawdzian `TestTuraZamknietaBledemNieOglaszaUkonczenia` **padl na kodzie sprzed naprawy** z wlasciwym zdaniem i przeszedl po przywroceniu; zakres wylacznie `internal/core`; kontrakt nietkniety |
| `nastawy-wdrozenia` | `teren/nastawy-wdrozenia` | `678ff3e` | weryfikacja Prowadzacego pomiarem: 4 pliki w zakresie, migracja 115 nietknieta, nowa migracja 401; sonda dodatnia wykonawcy pokazala **206 polaczen i 1,1 GB pobrania przed naprawa wobec zera po niej** |
| `powloka-i-wydanie` | `teren/powloka-i-wydanie` | `13b947b`, `8d99ae4` | weryfikacja Prowadzacego pomiarem: wykaz wydan sprowadzony do dwoch postaci hybrydowych (0 trafien wzorca natywna/AppImage przy 10 kontrolnych); sprawdzian wiazacy adres padl po wprowadzonym rozjezdzie i przeszedl po cofnieciu; 25 sprawdzianow powloki zdanych |
| `brama-i-droga-wejscia` | `teren/brama-i-droga-wejscia` | `23a6b84` wyjątek powitania · `e41dd78` straże drogi bez poczty · `9a8ec0c` brak skanera | weryfikacja Prowadzącego pomiarem: bieg wymuszony `-count=1` 537 s — 2029 sprawdzianów, 1 niezdany wobec 4 zastanych; powitanie niepełne odpowiada wersją protokołu na żywym rdzeniu, `channel.add` z brakiem pola dalej odmawia; kontrakt nietknięty |
| `prototypy` | `teren/prototypy` | paczki instalatora i drogi wejścia | przyjęte przez Właściciela; weryfikacja Prowadzącego pomiarem: oba okna wczytują się bez błędu konsoli, zero łańcuchów widocznych poza katalogiem treści |
| `fundament-klienta` | `teren/fundament-klienta` | `d19bfeb` warstwa połączenia i protokołu | weryfikacja Prowadzącego pomiarem: kompilacja bez błędu, 17 sprawdzianów zdanych, rozmowa z żywym rdzeniem, generat bajtowo powtarzalny, zero dotknięć DOM, kontrakt nietknięty |
| `naprawy-rdzenia` | `teren/naprawy-rdzenia` | `060d5b7` naprawy i brama kontraktu | weryfikacja Prowadzącego pomiarem: 2022 zdane wobec 2003 zastanych, te same 4 niezdane, kontrakt nietknięty |
| `proba-prototypow` | `teren/prototypy` | `66e5ee0` przepływ wejścia · `61a5867` moduł Studio · `53bc3d5` odsyłacze | kontrola sesji nadzorującej wykonanie, weryfikacja Prowadzącego pomiarem |
| `pomiar-stron` | `teren/pomiar-stron` | `607a3d8` kontrakt · `4191c4b` trzy osie pomiaru · `6b189e7` sprawdziany | kontrola osobnej sesji: kontrakt porównany strukturalnie, nie liniowo — zero sekcji zmienionych, żaden istniejący element nietknięty bajtowo, dodane 3 komendy, 6 struktur, 3 wyliczenia; generator powtórzony dwa razy, sumy zgodne; oś obciążenia zmierzona niezależnie własnym klientem kanału — 102 666 żądań wobec 102 669 na liczniku serwera sprawdzianu; osiem sprawdzianów zdanych, żaden pominięty. Scalone `73e9aad` |
| `zaplecze-modeli` | `teren/zaplecze-modeli` | `090acd6` silnik osadzeń na wagach stojących | kontrola osobnej sesji: teren zwrócony za niepełny wykaz postawionego, po uzupełnieniu przyjęty; kontroler powtórzył pomiar osadzeń niezależnie i przy odciętej sieci (`HF_HUB_OFFLINE=1`) — wymiar 1024, normy 1,000000, cos(kot,kot) 0,8348 wobec 0,2863–0,3353 dla par odległych; wagi `/opt` czytane, nie pobierane po raz drugi; kontrakt nietknięty. Scalone `ba82cb6` |
| `obraz-i-diagramy` | `teren/obraz-i-diagramy` | `f7fab06` dogniecenie zapisu i metadane osadzone | kontrola osobnej sesji własnym biegiem: osiem zapór zdanych nietkniętych, 33 sprawdziany zdane w realnych czasach z przytoczonymi rozmiarami przed i po (PNG 3322→2456 B, JPEG 45399→39623 B, model barw paleta dowodzi wejścia pngquanta); sprawdziany maszyny bez programu wytwarzają ją naprawdę (`t.Setenv` na pusty katalog); kontrakt nietknięty; zero śladu mmdc w rdzeniu. Scalone `5b98585` |
| `dokumenty-i-tekst` | `teren/dokumenty-i-tekst` | `ea05603` sześć programów treści pisanej · `aa695e6` poprawka nazwy po zwrocie | kontrola osobnej sesji: teren zwrócony za martwą nazwę `zasiegSyntezy` w komentarzu, po poprawce przyjęty; własny bieg kontrolera 14 zdanych, zero pominiętych; sprawdziany mierzą skutek — PDF czytany drugą komendą i innym programem, korekta aż po treść panelu w bazie, OCR dwoma przebiegami, brak programu wywołany `t.Setenv`; wybory hunspell i typst zweryfikowane uruchomieniem; kontrakt nietknięty. Scalone `43d05f6` |
| `odmowy-skanera` | `teren/odmowy-skanera` | `1eaa58a` parytet odmów skanera Linux wobec Windows | kontrola osobnej sesji: `bladWarstwySane` wierne lustro `bladWarstwyWia` — ten sam kod `channel_unavailable`, ta sama droga obejścia `studio.ingest.queue.add`, pakiet czytany z braku nie zaszyty; sprawdzian parytetu woła obie warstwy i wymaga jednego kodu; brak wymuszony atrapą PATH z potwierdzeniem `zewnetrzne.Stoi=false`; bieg celowany 3 zdane; kontrakt nietknięty. Scalone `1dca877` |
| `sprawdziany-drogi-wejscia` | `teren/sprawdziany-drogi-wejscia` | `1393e74` sprawdzian zdjęcia znacznika bramki | weryfikacja Prowadzącego mutacją rdzenia: wyłączenie `zdejmijZnacznikBezPoczty` daje sprawdzian niezdany („znacznik przeżył swój powód"), przywrócenie — zdany; trzy sprawdziany rozjazdu z rdzeniem były już przerobione przez `brama-i-droga-wejscia`, wykonawca to zmierzył i nie tknął; bieg celowany osiem zdanych; kontrakt nietknięty, jeden plik terenu. Scalone `4e70ef1` |
| `odwolania-do-usunietych-skryptow` | `teren/odwolania-do-usunietych-skryptow` | `8451157` żywy skrypt pakietu w odmowie braku serwera narzędzi | **wykonawca urwał się na limicie sesji przed rewizją; pracę dokończył i zweryfikował Prowadzący sam, bez niezależnej kontroli — odstępstwo od rozdziału ról, wymuszone urwaniem, w trybie samodzielnym na polecenie Właściciela.** Weryfikacja: droga naprawy w odmowie istnieje — `scripts/pakiet-serwera.sh` żyje i buduje `danaco-narzedzia` z `server/cmd/danaco-narzedzia`, stawiając obok rdzenia (w. 72,83,101); zero odwołań do zniesionych skryptów w rdzeniu; próba mutacji: podmiana stałej na `wydanie.sh` daje obie straże niezdane, przywrócenie — zdane; kontrakt nietknięty. Scalone `8e504c3` |
| `warunek-ukonczenia-zadania` | `teren/warunek-ukonczenia-zadania` | `06c81a1` rdzeń odróżnia ukończenie od trzech zatrzymań | weryfikacja Prowadzącego (tryb samodzielny) dwiema mutacjami: usunięcie mapowania `completed` daje „rozróżnienia maszynowego nie ma", wyłączenie samopodjęcia daje „ukończenie zachowuje się jak bramka akceptacji" — obie przywrócone zdane; kontrakt tylko dodania (`334705…` → `b7d04368…`), generator bajtowo powtarzalny w dwóch przebiegach; wykonawca ustalił pomiarem, że koniec tury koordynatora dotąd ginął (Wybudzacz odrzucał okna niebędące wykonawcami); brak bramki akceptacji zrobiony mechanizmem samopodjęcia, nie deklaracją. Scalone `e677a65` |
| `aktualizacja-powloki-i-skrypty` | `teren/aktualizacja-powloki-i-skrypty` | `d1b6207` probne.rs · `54cb287` kanał pobrań · `880b41f` skrypty natywne | kontrola osobnej sesji własnym biegiem: `cargo test` 23 zdane, zero niezdanych (stan zastany: nie kompilował się); sha256 sześciu kopii w materiale zamkniętym tożsame; `ADRES_KANALU` zgodny znak w znak z `kanal.adres` wykazu wydań; scalone `adb2a85`, drzewo scalone tożsame z kontrolowanym. Uwaga trwała: kompilacja powłoki wymaga `budowa/klient/dist` (w `.gitignore`) i zmiennych CARGO/RUSTUP ze środowiska maszyny |

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
