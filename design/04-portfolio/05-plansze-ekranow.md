# Plansze ekranów — galeria prezentacyjna

| | |
|---|---|
| **Produkt** | Danaco Console — AI Operating Environment |
| **Rodzaj** | Plansza portfolio design identity |
| **Wersja** | v2.0 · Status: Deweloperski |
| **Zakres** | Przeglądowa galeria wszystkich klikalnych prototypów okien platformy |
| **Źródło** | KANON.md rozdz. 7 (inwentarz okien) + katalog `05-okna/` (realne pliki) |

## 1. Przeznaczenie

Plansza jest punktem przeglądowym całego dorobku okien: prowadzi do każdego klikalnego prototypu,
pokazuje żywą miniaturę układu właściwą rodzinie okna i zestawia zbudowane ekrany z inwentarzem okien
z dokumentacji (KANON rozdz. 7) w formie uczciwej tabeli pokrycia. Wszystkie liczby — liczba linii i liczba
odwzorowanych interakcji — są odczytane z realnych plików, nie deklarowane.

## 2. Rodziny okien

| Rodzina | Zbudowane | Inwentarz | Pokrycie |
|---|---|---|---|
| Przepływ główny | 3 | 3 | pełne |
| Powłoki środowisk | 6 (4 powłoki + 2 uzupełniające MTAI) | 4 | pełne |
| Okna platformowe | 4 | 4 | pełne |
| Okna operacyjne modułów | 15 | 15 | pełne |
| **Razem** | **28 prototypów** | — | — |

Chat Window nie jest osobnym plikiem — występuje jako wspólny pas komunikacji rekonfigurowany
w kontekście każdego z 15 modułów, zgodnie z dokumentacją (`elementy-okien.md` rozdz. 3.6).

## 3. Statystyka dorobku

- 28 prototypów okien, wszystkie klikalne i dynamiczne, oba motywy równoprawne
- ~62 600 linii kodu prototypów
- ~3 100 odwzorowanych interakcji (przełączanie widoków, modale, listy wyboru, symulacje przejść stanów)

## 4. Decyzje projektowe

Miniatury-wireframe są zbudowane wyłącznie z żetonów systemu jako uproszczone szkielety proporcji powłoki —
nie są zrzutami ekranu, lecz schematami układu. Trzy warianty miniatury (przepływ z trzema strefami,
powłoka modułowa, nakładka modalna) odpowiadają trzem typom kompozycji ekranu platformy.

---
*Danaco Console — AI Operating Environment · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. — Dariusz Naharnowicz*
