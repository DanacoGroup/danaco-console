# Kreator instalacji — dokumentacja wykonawcza

Okno instalatora Danaco Console złożone ze składników. Ten dokument opisuje,
z czego okno jest zbudowane i jak to przejąć — żeby wykonawca aplikacji nie
musiał niczego wyławiać z pliku podglądu.

Podgląd `05-okna/platformowe/instalator.html` **nie jest produktem**. Jest
montażem: wpina warstwę wspólną i wskazuje miejsce, w którym okno ma stanąć.
Nie ma w nim ani jednego łańcucha widocznego dla użytkownika, ani jednej reguły
wyglądu, ani jednej linii skryptu.

---

## Co gdzie leży

| Warstwa | Plik | Czym jest |
|---|---|---|
| wartości wizualne | `zasoby/zetony/zetony.css` | barwy, odstępy, wymiary, kroje |
| fundament | `zasoby/css/fundament.css` | zerowanie, warstwa dokumentu |
| biblioteka wyglądu | `zasoby/css/komponenty.css` | klasy `.dn-*` wspólne wszystkim oknom |
| rodzina okna | `zasoby/kreator.css` | klasy `.dn-kreator-*` — układ kreatora |
| treść | `zasoby/okna/instalator/tresci.js` | **wszystkie** łańcuchy widoczne dla użytkownika |
| umowa licencyjna | `zasoby/tresci/licencja.js` | wygenerowana z `licencja-2.1.md` |
| narzędzia | `zasoby/okna/instalator/narzedzia.js` | budowa węzła, sięgnięcie po łańcuch |
| ikony | `zasoby/okna/instalator/ikony.js` | zestaw SVG używanych w oknie |
| składniki | `zasoby/okna/instalator/skladniki/` | 19 plików, jeden składnik na plik |
| ekrany | `zasoby/okna/instalator/ekrany/` | 6 plików, jeden krok na plik |
| mapa stanów | `zasoby/okna/instalator/stany.js` | co zmienia każdy stan warunkowy |
| montaż | `zasoby/okna/instalator/montaz.js` | złożenie okna z części |
| przepływ | `zasoby/okna/instalator.js` | mechanika: przejścia, zapis, dialogi |

Kolejność wpięcia arkuszy jest wiążąca — opisuje ją `zasoby/ARCHITEKTURA.md`.

---

## Zasada podziału

Trzy rzeczy trzymane są osobno, bo zmieniają się niezależnie i zmienia je kto
inny:

- **treść** — `tresci.js`, zmienia tłumacz i redaktor;
- **wygląd** — arkusze `.css`, zmienia projektant;
- **zachowanie** — `.js`, zmienia wykonawca.

Żaden plik składnika nie zawiera tekstu widocznego dla użytkownika. Składnik
przyjmuje **klucz katalogu**, nie napis:

```js
K.skladniki.naglowekBloku({ klucz: 'krok4.skroty.naglowek' })
```

Brak klucza w katalogu nie jest sytuacją do obsłużenia po cichu — do okna
wchodzi `⟨krok4.skroty.naglowek⟩`, żeby usterka rzucała się w oczy.

---

## Katalog treści

`tresci.js` ma budowę hierarchiczną — klucz odczytuje się jak ścieżkę:

```
krok3.warianty.x64.opis   →   tresci.krok3.warianty.x64.opis
```

Miejsca podstawień stoją w nawiasach klamrowych i **nie podlegają tłumaczeniu**:

```json
"pasekEtapu": "Etap {numer} z {ile} · {nazwa}"
```

```js
podstaw(tekst('krok5.pasekEtapu'), { numer: 2, ile: 4, nazwa: '…' })
```

### Dlaczego skrypt, a nie `.json`

Okno bywa otwierane wprost z dysku, dwuklikiem, a przeglądarka blokuje wtedy
pobieranie plików towarzyszących — `fetch` pliku `.json` kończy się błędem
i okno nie staje. Katalog wpina się więc znacznikiem `<script src>`, a plik
przypisuje gotowy obiekt:

```js
window.DanacoKreator.tresci = { … };
```

Zawartość pozostaje czystym JSON-em — klucze w cudzysłowach, bez przecinka po
ostatniej pozycji — więc narzędzia tłumaczy czytają ją tak samo jak wcześniej.
Ta sama zasada dotyczy umowy licencyjnej: `zasoby/tresci/licencja.js` wpisuje
fragment do rejestru `window.DanacoTresci`, skąd bierze go składnik `dokument`.
Wynika stąd wymóg kolejności — **pliki treści wpina się przed montażem**.

---

## Kontrakt składnika

Każdy składnik to funkcja `K.skladniki.<nazwa>(wlasciwosci)` zwracająca węzeł
DOM. Nie sięga do dokumentu, nie zna innych składników, nie trzyma stanu.
Wykaz właściwości stoi w nagłówku pliku składnika — poniżej skrót.

| Składnik | Właściwości |
|---|---|
| `nawigacjaKrokow` | `klucz`, `biezacy`, `naKrok` |
| `naglowekEkranu` | `nadtytul`, `tytul`, `podtytul`, `znak`, `dane` |
| `naglowekBloku` | `klucz` |
| `blokDanych` | `naglowek`, `wiersze`, `dane` |
| `poleWyboru` | `etykieta`, `opis`, `id`, `zaznaczone`, `dane`, `stan` |
| `poleSciezki` | `etykieta`, `wartosc`, `opis`, `id`, `zmien`, `naZmiane`, `stan` |
| `kartaOpcji` | `nazwa`, `opis`, `wartosc`, `grupa`, `plakietka`, `zaznaczona`, `fokus` |
| `baner` | `rodzaj`, `ikona`, `glowa`, `tresc`, `dane`, `ukryty` |
| `listaEtapow` | `naglowek`, `etapy`, `stany`, `poczatkowe` |
| `pasekPostepu` | `etykieta`, `opisPaska`, `wartosc`, `stan` |
| `pasekSzczegolow` | `dane` |
| `pasDzialan` | `poboczna`, `glowna`, `dodatkowa` |
| `frazaNawigacyjna` | `klucz`, `dane` |
| `tekstCiagly` | `klucze`, `dane` |
| `dokument` | `zrodlo`, `etykieta`, `rosnace` |
| `wierszMiary` | `wzor`, `dane`, `wytluszcz` |
| `blokBledu` | `szczegoly`, `rada`, `dane`, `ukryty` |
| `odsylaczPomocy` | `klucz`, `dane`, `naOtwarcie` |
| `oknoDialogowe` | `nazwa`, `dane`, `rola`, `tytul`, `tresc`, `cialo`, `czynnosci` |

Wartość `null` albo `false` pomija atrybut — dzięki temu warianty składnika
pisze się warunkiem, a nie rozgałęzieniem.

### Uchwyty

Mechanika nie szuka elementów po klasie wyglądu — klasa może się zmienić przy
przemalowaniu okna. Szuka po atrybucie `data-*`, który jest **kontraktem** —
52 uchwyty, m.in. `data-ekran`, `data-krok`, `data-tytul`, `data-wskazowka`,
`data-plakietka`, `data-niezgodnosc`, `data-zgoda-licencja`, `data-postep-tor`,
`data-etap`, `data-szczegoly`, `data-krok-dalej`, `data-krok-wstecz`.
Pełny wykaz daje:

```bash
grep -ohP '(data-\K[a-z-]+(?=[\]="]))|(dataset\.\K[a-zA-Z]+)' zasoby/okna/instalator.js | sort -u
```

Zmiana uchwytu jest zmianą umowy między składnikiem a przepływem — nazwy
uchwytów zmienia się w składniku i w przepływie jednocześnie.

---

## Ekrany i stany

Ekran rysuje się **raz**. Stan warunkowy nie jest rozgałęzieniem w kodzie
ekranu — jest daną, którą się na gotowym ekranie ustawia. Wszystko, co dany
stan zmienia, stoi w `stany.js`, a wartości są kluczami katalogu, nie tekstem.

| Ekran | Stany |
|---|---|
| 3 · wersja | `x64`, `arm`, `brak` |
| 5 · zapis | `przebieg`, `wycofywanie`, `blad` |
| 6 · wynik | `gotowe`, `ostrzezenia` |

Stany osiągalne adresem, do których przepływ sam nie doprowadzi:

```
?procesor=arm|brak       wykrycie procesora (krok 3)
?stan=blad|wycofywanie   odsłona zapisu (krok 5)
?wynik=ostrzezenia       wynik instalacji (krok 6)
```

---

## Bryła — animowana grafika kroku

`zasoby/bryla.js` jest składnikiem niezależnym od instalatora. Zakłada się na
każdy węzeł `.dn-bryla` z atrybutem `data-bryla`; barwy czerpie wyłącznie
z żetonów `--dn-bryla-*`, więc podąża za motywem i nie ma własnego tła.

Sterowanie przez `wezel.bryla`: `postep`, `pokaz`, `pracuje`, `domknij`,
`etap`, `etykiety`, `tempo`, `wstrzymaj`, `wznow`, `zdejmij`.

---

## Przejęcie okna

1. Skopiuj `zasoby/zetony/`, `zasoby/css/`, `zasoby/kreator.css`.
2. Skopiuj `zasoby/okna/instalator/` w całości i `zasoby/okna/instalator.js`.
3. Skopiuj `zasoby/bryla.js` i `zasoby/tresci/licencja.js`.
4. W swoim dokumencie umieść `<div data-kreator></div>`; montaż sam wykryje
   ten węzeł i złoży okno.
5. Wepnij pliki treści przed składnikami — kolejność wpięć podaje plik
   podglądu, który jest wzorem montażu.

Nic się nie dociąga w czasie działania: okno działa tak samo z serwera i wprost
z dysku. Sprawdzone porównaniem zrzutów — 5 kroków × 2 motywy, **0 pikseli
różnicy** między jednym a drugim.

Zdarzenie `kreator-gotowy` na dokumencie oznacza, że okno stoi i można podpiąć
własną mechanikę. `zasoby/okna/instalator.js` robi dokładnie to.

---

## Sprawdzian gotowości

Usuń z pliku okna wszystkie `<style>`, wszystkie `style=` i wszystkie
`<script>` bez `src`. Jeżeli okno wygląda i działa tak samo — jest zbudowane
ze składników, a nie zapakowane w jeden plik.

Stan obecny podglądu: **0 `<style>` · 0 `style=` · 0 `<script>` bez `src`**,
104 wiersze samego montażu.
