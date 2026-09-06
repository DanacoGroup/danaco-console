# Danaco Console — przewodnik wykonawcy

Aplikacja desktopowa (Tauri 2): rdzeń w Go, interfejs TypeScript/Vite, poczta
IMAP, WebSocket. Rdzeń wersji poprzedniej jest przejęty i działa; warstwa widoku
powstaje od nowa. Kolejność prac podaje
[plan etapów](prowadzenie/plan-etapow.md).

## Przed pierwszą zmianą

Przeczytaj [prowadzenie/ustroj-budowy.md](prowadzenie/ustroj-budowy.md).
W skrócie:

- Pracujesz w jednej roli: wykonawca terenu albo kontroler. Nie kontrolujesz
  własnej pracy.
- Pracujesz na jednym terenie, we własnym drzewie roboczym, wyłącznie na plikach
  wymienionych w [rejestrze terenów](prowadzenie/rejestr-terenow.md).
- Niczego nie dopowiadasz. Brak nazwy, wartości albo zachowania w źródle jest
  usterką do zgłoszenia, nie swobodą wykonawcy.
- Rozstrzygnięcie obowiązuje wyłącznie wtedy, gdy stoi
  w [rejestrze decyzji](prowadzenie/decyzje.md).

## Gdzie co leży

Jedno repozytorium: `~/budowa`. Poza nim nie ma drugiego drzewa z tą samą treścią.

| Ścieżka | Zawiera |
|---|---|
| `budowa/server/` | rdzeń w Go |
| `budowa/shared/` | `contract.json` — jedyne źródło prawdy typów — wraz z generatorem |
| `budowa/klient/` | nowy klient TypeScript; powstaje w terenie `fundament-klienta` |
| `budowa/desktop/` | powłoka Tauri |
| `design/` | system projektowy wraz z prototypami okien |
| `docs/` | dokumentacja projektowa |
| `prowadzenie/` | prowadzenie budowy; **znika przed wydaniem** |
| `narzedzia/` | skrypty budowy |

Klient wersji poprzedniej nie leży już w drzewie. Materiał do czytania i przeszczepu
wyjmuje się z historii: `git show klient-poprzedni-ostatni:budowa/klient-poprzedni/<ścieżka>`.

**Design stoi w `design/` tego repozytorium.** Nie ma dla niego osobnego drzewa
ani osobnej gałęzi; wcześniejszy warsztat `teren/prototypy` w `~/robocze/prototypy`
został zwinięty i wskazywanie go kieruje sesję w nieistniejące miejsce.
Prototypy okien etapu 1 leżą w `design/05-okna/`: `platformowe/instalator.html`,
`przeplyw/przeplyw-wejscia.html`, `przeplyw/centrum-dowodzenia.html`,
`srodowiska/talkin-przedsionek.html`, `moduly/studio.html`.

Materiał zabezpieczony poza gitem stoi w `~/robocze/material/` — kopie
bezpieczeństwa, punkt kontrolny wersji poprzedniej i opracowania zamkniętego
podejścia. **Nie jest źródłem prawdy i nie wchodzi do budowy.**

## Gałęzie

| Gałąź | Zawiera | Wolno pisać |
|---|---|---|
| `main` | całość budowy | tylko Prowadzący budowę |
| `teren/<nazwa>` | praca jednej sesji | sesja prowadząca teren |
| `refs/przeniesienie/*` | dorobek zamkniętego podejścia | nikt — materiał do czytania |

Materiał pod `refs/przeniesienie/` **nie jest źródłem prawdy** — pochodzi
z zamkniętego podejścia do budowy. Odwołuje się do usuniętego drzewa kodu;
odwołanie do nieistniejącego pliku nie jest wskazówką, jest pozostałością.

**Etap 1 prowadzi Właściciel własnoręcznie.** Kompozycji okna nie przekazuje się
zleceniem. Sesja nie proponuje układu okna Studio, nie otwiera na nie terenu
i nie zmienia go bez wyraźnego polecenia.

## Otwarcie terenu

```bash
bash narzedzia/nowy-teren.sh <nazwa> <galaz-bazowa>
```

Praca toczy się w `~/robocze/<nazwa>`. Zamknięcie terenu: zgłoszenie do
Prowadzącego, kontrola przez inną sesję, scalenie, `git worktree remove`.

## Zasady bezwzględne

- Rewizja obejmuje wyłącznie pliki własnego terenu — nigdy `git add -A`.
- Komunikat rewizji: tryb oznajmujący, jedno zdanie, do 70 znaków.
- Bez zmian zbiorczych: żadnego `sed -i` po katalogach ani `find -exec`.
- Bez narośli plikowych: nowa treść wchodzi edycją dokumentu kanonicznego.
  Notatki, podsumowania i plany robocze zostają w katalogu tymczasowym sesji.
- Bez kroniki: dokument i komentarz opisują stan obecny i jego uzasadnienie,
  nigdy przebieg prac.
- Bez wymyślonych oznaczeń: żadnych autorskich kodów, sygnatur i numeracji.
- Jedno pojęcie — jedna nazwa, w całym repozytorium.
- Komentarze: plik od 100 wierszy kodu niesie ich najwyżej 5% wierszy, komentarz
  do 350 znaków; sprawdza walidator pluginu `danaco-dyscyplina` (decyzja 35).

## Weryfikacja

Nie deklarujesz, że praca jest gotowa — wykazujesz to. Przytaczasz rzeczywisty
wynik uruchomienia i wprost nazywasz to, czego sprawdzić się nie dało.

Wykaz narzędzi dostępnych na maszynie wraz z wersjami prowadzi
[prowadzenie/srodowisko-maszyny.md](prowadzenie/srodowisko-maszyny.md).
Wszystkie są zainstalowane — nie instaluj niczego. Brakujące narzędzie jest
zgłoszeniem do Prowadzącego.

**Wyjątek udziela się terenowi wpisem w rejestrze terenów**, nazywa go wprost
i wygasa razem z tym terenem. Wszystko, co teren postawi, trafia do raportu
odbioru wraz z wagą na dysku i do `prowadzenie/srodowisko-maszyny.md`.

Przed zamknięciem terenu wygaś procesy, które uruchomiłeś, i usuń pliki robocze.
Sesja, która zostawia po sobie działający proces, nie zamknęła pracy.
