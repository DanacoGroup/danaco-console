# Danaco Console — przewodnik wykonawcy

Aplikacja desktopowa (Tauri 2): rdzeń w Go, interfejs TypeScript/Vite, poczta
IMAP, WebSocket. Repozytorium jest w fazie pierwszego okna modułowego —
**kod aplikacji jeszcze nie powstaje**. Kolejność prac podaje
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

## Gałęzie

| Gałąź | Zawiera | Wolno pisać |
|---|---|---|
| `main` | prowadzenie budowy, narzędzia | tylko Prowadzący budowę |
| `teren/proba-prototypow` | przepływ wejścia i okno Studio | **wyłącznie Właściciel** — warsztat etapu 1 |
| `przebudowa/design`, `przebudowa/dokumentacja` | dorobek zastany | nikt — materiał do czytania |
| `teren/<nazwa>` | praca jednej sesji | sesja prowadząca teren |

Materiał na gałęziach przebudowy **nie jest źródłem prawdy** — pochodzi
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

## Weryfikacja

Nie deklarujesz, że praca jest gotowa — wykazujesz to. Przytaczasz rzeczywisty
wynik uruchomienia i wprost nazywasz to, czego sprawdzić się nie dało.

Wykaz narzędzi dostępnych na maszynie wraz z wersjami prowadzi
[prowadzenie/srodowisko-maszyny.md](prowadzenie/srodowisko-maszyny.md).
Wszystkie są zainstalowane — nie instaluj niczego. Brakujące narzędzie jest
zgłoszeniem do Prowadzącego.

Przed zamknięciem terenu wygaś procesy, które uruchomiłeś, i usuń pliki robocze.
Sesja, która zostawia po sobie działający proces, nie zamknęła pracy.
