# Okna wejścia — dokumentacja wykonawcza

Dwa okna, które użytkownik widzi przed wejściem do programu: **uruchomienie**
(nawiązanie połączenia) i **dostęp do konta** (logowanie, zakładanie konta,
odzyskiwanie dostępu). Ten dokument opisuje, z czego są zbudowane i jak to
przejąć.

Podgląd `05-okna/przeplyw/przeplyw-wejscia.html` **nie jest produktem**. Jest
rusztowaniem, na którym ogląda się okna: niesie przybornik etapów, dziennik
komunikatów i symulację uzgodnienia protokołu. Same okna wstawia w dwa miejsca:

```html
<div data-wejscie-okno="uruchomienie"></div>
<div data-wejscie-okno="dostep"></div>
```

W pliku podglądu nie ma ani jednej reguły wyglądu, ani jednej linii skryptu,
ani jednego łańcucha widocznego dla użytkownika.

---

## Co gdzie leży

| Warstwa | Plik | Czym jest |
|---|---|---|
| wartości wizualne | `zasoby/zetony/zetony.css` | barwy, odstępy, wymiary, kroje |
| fundament | `zasoby/css/fundament.css` | zerowanie, warstwa dokumentu |
| biblioteka wyglądu | `zasoby/css/komponenty.css` | klasy `.dn-*` wspólne wszystkim oknom |
| rodzina okna | `zasoby/wejscie.css` | klasy `.we-*` i `.au-*` — układ okien wejścia |
| narzędzia | `zasoby/narzedzia-okien.js` | budowa węzła, sięgnięcie po łańcuch — wspólne dla okien |
| treść | `zasoby/okna/wejscie/tresci.js` | **wszystkie** łańcuchy widoczne dla użytkownika |
| znaki | `zasoby/okna/wejscie/ikony.js` | 22 rysunki używane w oknach |
| składniki | `zasoby/okna/wejscie/skladniki/` | 15 plików, jeden składnik na plik |
| ekrany | `zasoby/okna/wejscie/ekrany/` | `uruchomienie.js` (3 odsłony), `dostep.js` (7 odsłon) |
| montaż | `zasoby/okna/wejscie/montaz.js` | złożenie okien i wstawienie w miejsce montażu |
| mechanika okna | `zasoby/okna/przeplyw-wejscia.js` | animacja startowa, przebieg łączenia, sprawdzanie danych |
| oprawa podglądu | `zasoby/okna/przeplyw-wejscia-podglad.js` | rusztowanie — **nie należy do produktu** |
| animacje | `zasoby/ekran-startowy.js`, `zasoby/powloki.js` | składniki biblioteki, każdy z własnym wykazem |

Kolejność wpięcia arkuszy jest wiążąca — opisuje ją `zasoby/ARCHITEKTURA.md`.

---

## Zasada podziału

Trzy rzeczy trzymane są osobno, bo zmieniają się niezależnie i zmienia je kto
inny: **treść** (`tresci.js` — tłumacz i redaktor), **wygląd** (arkusze `.css`
— projektant), **zachowanie** (`.js` — wykonawca).

Żaden plik składnika nie zawiera tekstu widocznego dla użytkownika. Składnik
przyjmuje **klucz katalogu**, nie napis:

```js
S.naglowekEkranu(N, { tytul: 'dostep.logowanie.tytul', lid: 'dostep.logowanie.lid' })
```

Brak klucza nie jest sytuacją do obsłużenia po cichu — do okna wchodzi
`⟨dostep.logowanie.tytul⟩`, żeby usterka rzucała się w oczy.

### Katalog treści

`tresci.js` ma budowę hierarchiczną — klucz odczytuje się jak ścieżkę:

```
dostep.odzyskiwanie.haslo.nowe  →  tresci.dostep.odzyskiwanie.haslo.nowe
```

Miejsca podstawień stoją w nawiasach klamrowych i **nie podlegają tłumaczeniu**:

```json
"lid": "Na adres {adres} został wysłany sześciocyfrowy kod. Zachowuje ważność przez {minuty} minut.",
"lidDane": { "adres": "operator@danaco-group.pl", "minuty": 10 }
```

Skrypt, a nie plik `.json`, bo okno bywa otwierane wprost z dysku, a przeglądarka
blokuje wtedy pobieranie plików towarzyszących. Zawartość pozostaje czystym
JSON-em, więc narzędzia tłumaczy czytają ją tak samo.

---

## Kontrakt składnika

Składnik to funkcja `DanacoWejscie.skladniki.<nazwa>(N, wlasciwosci)`, gdzie `N`
to narzędzia związane z katalogiem. Zwraca węzeł DOM albo **wykaz węzłów**, gdy
jego części są rodzeństwem w kolumnie panelu. Nie sięga do dokumentu, nie zna
innych składników, nie trzyma stanu.

| Składnik | Właściwości |
|---|---|
| `belkaOkna` | `tytul` |
| `kolumnaTozsamosci` | `odslona` (`uruchamianie` \| `dostep`) |
| `notaWydawcy` | — |
| `naglowekEkranu` | `nadtytul`, `tytul`, `lid`, `lidWezly`, `poTytule`, `dane` |
| `listaEtapow` | `stany`, `miary` |
| `baner` | `rodzaj`, `ikona`, `glowa`, `tresc`, `dane`, `id` |
| `frazaNawigacyjna` | `klucz` albo `pytanie` + `czynnosc` + `cel` + `grupa` |
| `pasDzialan` | `widok`, `grupa`, `aktywny`, `czynnosci` |
| `zakladkiPigulki` | `wybrana` |
| `poleTekstowe` | `etykieta`, `id`, `typ`, `uzupelnij`, `opis`, `bledne`, `opisuje` |
| `poleHasla` | `etykieta`, `id`, `uzupelnij`, `opis`, `bledne`, `opisuje` |
| `miernikSily` | `dla` (identyfikator pola hasła) |
| `poleSesji` | — |
| `metodyLogowania` | `dostepne` — zwraca **wykaz** węzłów |
| `poleKodu` | `odliczanie`, `czynnosc` (`wklej` \| `ponow`) — zwraca **wykaz** |
| `krokiOdzyskiwania` | `biezacy` (1..3) |

Wartość `null` albo `false` pomija atrybut — dzięki temu warianty składnika
pisze się warunkiem, a nie rozgałęzieniem.

### Uchwyty

Mechanika nie szuka elementów po klasie wyglądu — klasa może się zmienić przy
przemalowaniu okna. Szuka po atrybucie `data-*`, który jest **kontraktem**.
Szesnaście uchwytów; najważniejsze:

| Uchwyt | Kto go stawia | Co znaczy |
|---|---|---|
| `data-wejscie-okno` | podgląd | miejsce, w którym ma stanąć okno |
| `data-widok`, `data-widok-aktywny`, `data-grupa-widoku` | ekrany i pasy | odsłona przełączana razem ze swoim pasem |
| `data-po-polaczeniu` | ekran uruchomienia | etap, do którego okno przechodzi samo |
| `data-uruchomienie-scena`, `data-uruchomienie` | podgląd | scena i warstwa animacji startowej |
| `data-ekran-startowy`, `data-powloki`, `data-bryla` | okno | pola składników animowanych |
| `data-sila-dla` | `miernikSily` | pole hasła, które ten miernik ocenia |
| `data-odsloniecie` | `poleHasla` | pole, które ta kontrolka odsłania |
| `data-wklej-kod` | `poleKodu` | kontrolka wklejenia kodu ze schowka |
| `data-wartosc` | znacznik | miara postępu w procentach |

Pełny wykaz daje:

```bash
grep -ohP '(data-\K[a-z-]+(?=[\]="]))|(dataset\.\K[a-zA-Z]+)' zasoby/okna/przeplyw-wejscia.js | sort -u
```

Zmiana uchwytu jest zmianą umowy — nazwę zmienia się w składniku i w mechanice
jednocześnie.

---

## Odsłony

**Okno uruchomienia** (grupa `wariant`) — trzy odsłony jednego okna:

| `data-widok` | Kiedy | Czynność główna |
|---|---|---|
| `w-laczenie` | łączenie z serwerem | brak — przechodzi samo |
| `w-token` | urządzenie rozpoznane | brak — przechodzi samo |
| `w-blad` | serwer nie odpowiada | „Spróbuj ponownie" |

**Okno dostępu** (grupa `stan`) — siedem odsłon:

`logowanie` · `logowanie-blad` · `rejestracja` · `kod` ·
`odzyskiwanie-adres` · `odzyskiwanie-kod` · `odzyskiwanie-haslo`

Zakładki niosą wyłącznie dwie drogi równorzędne — logowanie i rejestrację.
Odzyskiwanie dostępu nie jest trzecią drogą, tylko wyjściem z logowania, więc
w zakładkach zostaje zaznaczone logowanie.

Odsłony osiągalne adresem, do których przepływ sam nie doprowadzi:

```
?rejestracja=login-zajety | email-bledny | haslo-slabe | hasla-rozne
```

---

## Co robi mechanika okna

`zasoby/okna/przeplyw-wejscia.js` rusza dopiero na zdarzenie `wejscie-gotowe`,
zgłaszane przez montaż — wcześniej nie ma czego obsługiwać.

1. **Uruchomienie.** Animacja startowa gra przed oknem; okno jest w tym czasie
   ukryte (`visibility`), więc nie widzi go ani oko, ani czytnik ekranu.
   Po zdarzeniu `ekran-startowy-koniec` okno staje i rusza przebieg łączenia.
2. **Przebieg łączenia.** Cztery etapy po 900 ms, potem przejście do etapu
   wskazanego przez `data-po-polaczeniu`. Bez czynności użytkownika — nie ma tu
   czego wybierać.
3. **Sprawdzanie danych przed wysłaniem.** Pięć formularzy, jedna reguła
   wymagań hasła wspólna z miernikiem siły. Rozpoznania: brak danych, brak
   loginu, brak hasła, brak adresu, login zajęty, adres nieprawidłowy, hasło
   poniżej wymagań, hasła niezgodne. Wszystkie pola puste to **jedna** sprawa,
   nie kilka. Napisy stoją w katalogu; w mechanice zostaje sam kontrakt z polami.
4. **Pola kodu.** Wiązane grupami — każdy zestaw rządzi się sam, kursor nie
   przeskakuje między oknami. Przyjmują wyłącznie cyfry.
5. **Odsłonięcie hasła.** Zmienia typ pola i własną etykietę, więc czytnik
   ekranu wie, w którym stanie stoi.

---

## Przejęcie okien

1. Skopiuj `zasoby/zetony/`, `zasoby/css/`, `zasoby/wejscie.css`.
2. Skopiuj `zasoby/narzedzia-okien.js` i `zasoby/okna/wejscie/` w całości.
3. Skopiuj `zasoby/okna/przeplyw-wejscia.js` oraz animacje:
   `zasoby/ekran-startowy.js`, `zasoby/powloki.js`.
4. **Nie kopiuj** `zasoby/okna/przeplyw-wejscia-podglad.js` — to rusztowanie
   podglądu, nie produkt.
5. W swoim dokumencie umieść `<div data-wejscie-okno="uruchomienie"></div>`
   albo `<div data-wejscie-okno="dostep"></div>`; montaż sam wykryje ten węzeł.
6. Wepnij pliki treści **przed** składnikami — kolejność wpięć podaje plik
   podglądu, który jest wzorem montażu.

Nic się nie dociąga w czasie działania: okna działają tak samo z serwera i wprost
z dysku.

---

## Sprawdzian gotowości

Usuń z pliku okna wszystkie `<style>`, wszystkie `style=` i wszystkie `<script>`
bez `src`. Jeżeli okno wygląda i działa tak samo — jest zbudowane ze składników,
a nie zapakowane w jeden plik.

Stan obecny podglądu: **0 `<style>` · 0 `style=` · 0 `<script>` bez `src`**.
