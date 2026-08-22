# Źródła prawdy — czytaj przed pierwszą zmianą

Kształt produktu rozstrzygają **wyłącznie dwa katalogi**: `docs/` i `design/`.
Odpowiedzi nie ma nigdzie indziej — ani w kodzie, ani w komentarzach, ani
w notatkach po poprzednich wykonawcach. Czego nie ma w tych dwóch katalogach,
tego nie ma w produkcie.

## Czym jest każde źródło

| Źródło | Czym jest | Czym nie jest |
|---|---|---|
| `docs/` | opracowania funkcjonalne: komplet okien, funkcje, pola, czynności, stany | nie jest tekstem do wstawienia w okno |
| `design/` | prototypy HTML, opracowania systemu projektowego, żetony, marka | nie jest źródłem kopiowania treści |

**`docs/` mówi, co okno ma robić.** Zdanie opisujące zachowanie — na przykład
„niska pewność rozpoznania stawia plakietkę" — jest poleceniem dla wykonawcy:
zbuduj plakietkę. Nie jest napisem do wyświetlenia Operatorowi.

**Prototypy w `design/` pokazują wygląd i rozmieszczenie elementów.** Obowiązują
jako **minimum wytycznych**: układ, proporcje, hierarchia, gęstość. Prototypy są
okrojone z funkcji i komponentów, więc **nie odtwarza się ich jeden do jednego** —
wzorujesz się na wyglądzie, a pełny komplet funkcji bierzesz z `docs/`.

**Treść interfejsu pisze wykonawca.** Nie ma jej w żadnym pliku i nie wolno jej
skądkolwiek przepisywać. Etykieta działania to jedno do trzech słów. Nagłówek to
nazwa rzeczy, nie zdanie o niej. Stan pusty mówi, czego brak, i podaje jedno
wyjście. Komunikat błędu mówi, co się stało i co zrobić. W oknie nie stoi proza
o mechanice systemu, odsyłacz do rozdziału ani nazwa komendy kontraktu.

## Nietykalność

`docs/` i `design/` są **tylko do czytania**. Nie usuwaj, nie przenoś, nie
przemianowuj i nie redaguj niczego w środku. Gdy coś wygląda na nadmiarowe albo
sprzeczne — zgłoś Właścicielowi i czekaj na rozstrzygnięcie.

Opracowania produktu (`README.md`, `INSTALACJA-I-KONFIGURACJA.md`,
`INSTRUKCJA-UZYTKOWANIA.md`, `LICENSE.md`) mieszkają wyłącznie w `docs/`. Korzeń
repozytorium nie trzyma ich kopii — kopia rozjeżdża się z oryginałem i przestaje
być źródłem prawdy.

## Stan zbioru `docs/` — domknięty 2026-08-20

Zbiór jest **zamknięty i wiążący**: 51 opracowań  5 604 787 znaków. Od tej daty
`docs/` jest jedynym źródłem prawdy o tym  **czym platforma jest i jak ma być
zbudowana**. Rozstrzyga zakres  nazwy  wartości i zachowania; `design/` pokazuje
wygląd  ale nie rozstrzyga funkcji.

| Wejście do zbioru | Plik |
|---|---|
| Opis produktu i brama zbioru | `docs/README.md` |
| Katalog wszystkich opracowań | `docs/SPIS-OPRACOWAN.md` |
| Wzór redakcyjny  język i konwencje budowy | `docs/STANDARD-REDAKCYJNY-I-JEZYKOWY.md` |
| Warunki korzystania | `docs/LICENSE.md` |

Każde opracowanie niesie metrykę  spis treści i stopkę  co najmniej 30 % objętości
w formach wizualnych oraz wykazy normatywne: komendy kontraktu  etykiety  żetony 
komponenty  skróty  stany  punkty łamania i kryteria odbioru.

**Zasada zamkniętego zbioru.** Deweloper nie dopowiada niczego z własnej głowy.
Jeśli opracowanie nie podaje nazwy komendy  brzmienia etykiety  wartości żetonu
albo kodu błędu — jest to usterka opracowania do zgłoszenia  nie swoboda
wykonawcy. Miejsca nierozstrzygnięte są oznaczone wprost jako
**[DO DECYZJI OPERATORA]** i czekają na rozstrzygnięcie Właściciela.

**Wartości pochodzą wyłącznie ze źródeł normatywnych** — `budowa/shared/contract.json`
(komendy  zdarzenia  struktury  wyliczenia  kody błędów)  `design/zasoby/zetony/zetony.css`
(żetony)  arkusze `design/zasoby/css/` (komponenty)  prototypy `design/05-okna/`
(dosłowne brzmienia etykiet). Nazwa spoza tych źródeł nie wchodzi do dokumentacji
ani do kodu.

## Gdzie szukać warstwy wizualnej

Wartości wiążące stoją w pełnych opracowaniach, nie w streszczeniach:

| Czego szukasz | Opracowanie |
|---|---|
| żetony `--dn-*`, barwa, typografia, odstępy | `design/01-dokumentacja-md/04-tokens.md` |
| komponenty `.dn-*` i ich stany | `design/01-dokumentacja-md/06-components.md` |
| organizacja arkuszy i konwencje klas | `design/01-dokumentacja-md/05-styles-css.md` |
| układ okien i stanowiska | `design/01-dokumentacja-md/03-design-view.md` |
| ikony i widżety | `design/01-dokumentacja-md/07-icons-widgets.md` |
| ruch, dostępność, przekazanie | `design/01-dokumentacja-md/08-handoff-motion-dostepnosc.md` |
| znak, marka, ton | `design/01-dokumentacja-md/09-brand-system.md`, `design/03-marka/` |

## Praca wielu sesji naraz

Kilka sesji nie pisze równolegle do tego samego pliku — dzieli się **drzewem**,
a wyniki scala się po zamknięciu terenu.

- `main` — stan scalony i jedyna pozycja odniesienia.
- `teren/<nazwa>` — gałąź jednej sesji, na przykład `teren/uwierzytelnienie`,
  `teren/strona-glowna`, `teren/multitaskingai`.
- Terenów nie wolno nakładać: dwie gałęzie nie obejmują tego samego katalogu.
  Gdy zmiana wykracza poza teren, sesja ją zgłasza zamiast sięgać po cudze pliki.
- Scalenie do `main` następuje po zamknięciu terenu, nie w trakcie.

Commit obejmuje **wyłącznie pliki własnego terenu** — nigdy `git add -A`, gdy
w drzewie leży żywa praca innej sesji. Wpis w dzienniku: tryb oznajmujący, jedno
zdanie, do 70 znaków.

## Zakazy

- **Bez wymyślania.** Okno, nazwa, pole, metryka ani numeracja bez pokrycia
  w źródłach nie wchodzi do produktu.
- **Bez zmian zbiorczych.** Żadnego `sed -i` po katalogach ani `find -exec` —
  plik czytasz, zmieniasz celowaną zmianą.
- **Bez testów i bram.** Miarą pracy jest zgodność okna ze źródłami, nie zielony
  wynik. W drzewie zostają wyłącznie zapory pilnujące uruchamiania procesów
  i świeżości generatu kontraktu.
- **Bez kroniki w kodzie.** Komentarz mówi, dlaczego tak jest teraz — nigdy co
  było wcześniej. Jeden wpis na blok, do 70 znaków.
