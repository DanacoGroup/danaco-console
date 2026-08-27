# Danaco Console — architektura arkuszy i skryptów

Każdy plik niesie jedną kompetencję i jest jej jedynym właścicielem. Nowe okno
prototypu składa się z wpięcia potrzebnych plików w podanej kolejności — nie
z pisania stylów od nowa i nie z nadpisywania wartości z plików niższych warstw.

## Kolejność wpięcia

| # | Plik | Kompetencja | Zależy od |
|---|---|---|---|
| 1 | `zetony/fonty.css` | kroje pisma | — |
| 2 | `zetony/zetony.css` | żetony wartości wizualnych (barwy, pismo, przestrzeń, wymiary, warstwy) | wpina `zetony/ruch.css` |
| 3 | `zetony/ruch.css` | czasy, krzywa, wartości uwydatnienia, klatki kluczowe, `prefers-reduced-motion` | żetony |
| 4 | `css/fundament.css` | zerowanie, typografia bazowa, siatka | żetony |
| 5 | `css/komponenty.css` | biblioteka komponentów (przyciski, pola, karty, kafle, modale, plakietki); wymiar znaku w przycisku należy do przycisku | żetony, ruch; wpina `menu.css` |
| 5a | `menu.css` | panel menu rozwijanego: bryła, widełki szerokości, przewijanie, kotwica prawa | żetony |
| 6 | `stany.css` | jeden kod stanu: praca · reakcja · błąd · zakończone | żetony |
| 7 | `rama.css` | obudowa okna: belka, pas narzędzi, szyna nawigacji, obszary i panele, pasek stanu | żetony, komponenty |
| 8 | `karty-okna.css` | pasmo kart okna roboczego, anatomia karty, karta przypięta, karta grupy, sterowanie pasma | żetony, stany |
| 9 | `panel-sesji.css` | zawartość panelu bocznego: płaski wykaz sesji z sortowaniem i zawężeniem, drzewo projektów, nagłówek karty | rama, stany |
| 10 | `pasek-okna.css` | pasek narzędzi modułu wewnątrz okna roboczego | żetony, komponenty |
| 11 | `izolacja.css` | wskaźniki incognito, gałęzi, zakresu widzenia i uprawnień okna | żetony |
| 12 | `okno-robocze.css` | maksymalizacja okna roboczego | rama |
| 12a | `okna/<okno>.css` po nim | płótno karty modułowej i pulpit Centrum | wszystkie wyżej |
| 13 | `stanowisko.css` | wnętrze stanowiska: okna operacyjne, pas komunikacji, żetony polecenia | żetony, komponenty |
| 14 | `okna-modalne.css` | okna nakładkowe | komponenty |
| 15 | `przedsionek.css` | przedsionek środowiska | żetony, komponenty |
| 15a | `kreator.css` | okno kreatora: trzy pasma (belka · korpus · pas działań), szyna kroków, płótno treści | żetony, komponenty |
| 16 | `prototyp.css` | oprawa prototypu: plakietka, pasek dolny, nota o oknie | żetony |
| 17 | `okna/<okno>.css` | złożenie jednego okna — wyłącznie układ, bez definicji komponentów | wszystkie wyżej |

Okno, które nie ma nic własnego poza znacznikiem, **nie zakłada pliku z wiersza 17**.
Instalator jest takim oknem: wpina żetony, fundament, komponenty, `kreator.css`
i `prototyp.css`, a całą jego mechanikę niesie `okna/instalator.js`.

## Skrypty

| Plik | Kompetencja |
|---|---|
| `wspolne.js` | motyw, ogłoszenia dla czytnika ekranu, wspólne narzędzia |
| `prototyp.js` | mechanika `data-menu`, dymki, kopiowanie, oprawa prototypu |
| `okna-modalne.js` | otwieranie i zamykanie okien nakładkowych |
| `rama.js` | belka, pas narzędzi, nawigacja historii |
| `stanowisko.js` | wnętrze stanowiska |
| `stany.js` | stan zbiorczy pozycji nadrzędnej (`data-stan-zbiorczy`) |
| `pasek-stanu.js` | miary paska stanu (`data-udzial`) |
| `karty-okna.js` | karty, grupy, menu powłok, grot, „+", ustawienia widoku |
| `pasek-okna.js` | obecność paska narzędzi i zwijanie nadmiaru pod „…" |
| `panel-sesji.js` | grupowanie wykazu, drzewo, czynności wiersza, akcje zbiorcze |
| `okno-robocze.js` | maksymalizacja, skrót, Esc, przełączniki izolacji |
| `bryla.js` | bryła warstw platformy — scena przestrzenna w płótnie, barwy z żetonów |
| `okna/<okno>.js` | złożenie jednego okna |

## Zasady wiążące

1. **Jeden właściciel.** Klasa ma definicję w jednym pliku. Reguła w pliku okna,
   która zmienia wartość zadeklarowaną w arkuszu kompetencyjnym, jest błędem —
   wartość albo trafia do arkusza kompetencyjnego, albo dostaje własny wariant
   (klasę, atrybut `data-*`).
2. **Zależność zamiast kolejności.** Reguła nie może opierać się na tym, że
   arkusz wpina się później. Warianty rozstrzyga selektor stanu
   (`:has(> .dn-rama-korpus)`, `[data-*]`), nie pozycja w `<head>`.
3. **Tylko żetony.** Barwy, odstępy i stopnie pisma pochodzą z żetonów.
   Geometria zatwierdzona pomiarem wzorca (pasmo kart) stoi w jednym bloku
   zmiennych komponentu w arkuszu właściciela.
4. **Wygląd w arkuszu, mechanizm w skrypcie.** Plik okna nie zawiera `<style>`,
   `<script>` bez `src` ani wstawek `style="…"`.
5. **`[hidden]` rozstrzyga fundament.** `fundament.css` niesie
   `[hidden] { display: none !important }`, więc atrybut `hidden` wygrywa z każdym
   `display` autora — ukrycie deklarujemy tym atrybutem, a reguły układu ustawiają
   `display` tylko dla stanu widocznego. Odsłonięcie wbrew `[hidden]` wymaga
   własnego atrybutu lub klasy wariantu, nie samego `display`.
6. **Animacja wejścia z `backwards`.** `both` pozostawia animację czynną, a taki
   element jest kontekstem nakładania i blokiem odniesienia dla potomków
   `position: fixed`.
7. **Ikony z zestawu.** `ikony/svg/` + `manifest.json`; brak ikony zgłasza się
   Właścicielowi, nie rysuje własnej.

## Kontrakt znaczników

| Atrybut | Właściciel | Znaczenie |
|---|---|---|
| `data-stan` | `stany.css` | praca · reakcja · blad · zakonczone · brak |
| `data-tetno` | `stany.css` | tętno kropki pracy w strefie bieżącej |
| `data-stan-zbiorczy` | `stany.js` | selektor znacznika, który przyjmuje stan najwyższy |
| `data-karta`, `data-karta-rodzaj`, `data-modul`, `data-grupa` | `karty-okna.*` | tożsamość karty pasma |
| `data-grupowanie`, `data-gestosc` | `karty-okna.*` | ustawienia widoku pasma |
| `data-okno`, `data-okno-max`, `data-izolacja` | `okno-robocze.*` | stan okna roboczego |
| `data-projekt`, `data-srodowisko` | `panel-sesji.*` | wymiary grupowania wykazu |
| `data-narzedzie`, `data-nadmiar`, `data-lustro` | `pasek-okna.*` | narzędzia modułu i ich zwijanie |
| `data-udzial` | `pasek-stanu.js` | udział miary paska stanu |

## Zależności ustalone pomiarem

| Zależność | Powód |
|---|---|
| `css/komponenty.css` wpina `menu.css` | menu jest komponentem biblioteki; 33 prototypy wpinają samą bibliotekę i bez tego zostałyby bez panelu menu |
| `rama.css` trzyma zapasowe położenie menu ramy | widoki bez `menu.js` muszą mieć poprawne położenie menu szyny, stopki, grupy funkcji globalnych i sterowania panelu |
| karty okna głównego obsługuje wyłącznie `karty-okna.js` | dwa właściciele przełączania kart gasiły płótno karty bieżącej, gdy dwie karty wskazywały ten sam panel |
| `.cd-tresc` i karty panelu deklarują `display` wyłącznie z warunkiem `:not([hidden])` | `[hidden]` z `fundament.css` niesie `!important`, więc `display` autora go nie przebije; regułę układu pisze się z `:not([hidden])`, żeby nie stosowała się do karty odłożonej |
| wymiar znaku deklaruje przycisk (`.dn-btn`, `.dn-btn-ikona`) i pozycja listwy (`.cd-ik`) | `svg` bez zadeklarowanego wymiaru rośnie do wysokości pojemnika |
| zwijanie nadmiaru paska modułu prowadzi pomiar, nie suma szerokości | odstępy i rozdzielacze zmieniają się razem ze składem paska |
