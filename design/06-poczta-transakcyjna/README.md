# 06 · Poczta transakcyjna

Szablony wiadomości wychodzących z adresu `noreply@danaco-group.pl`
w ramach Danaco Console. Zgodne z księgą znaku, typografią i głosem
oraz zastosowaniami marki (rozdz. 4.6).

Każdy list ma dwie części: `text/html` i `text/plain`.
Wariant jasny jest podstawą, ciemny — zadeklarowanym nadpisaniem.

| Plik | Zawartość |
|---|---|
| `opracowanie-techniczne.md` | decyzje projektowe, żetony, ograniczenia klientów pocztowych, wariant ciemny, część tekstowa, zmienne, bezpieczeństwo, lista kontrolna |
| `html/list-01-aktywacja-konta.html` | kod aktywacji nowego konta Operatora |
| `html/list-02-uwierzytelnienie-logowania.html` | kod logowania (drugi składnik) |
| `html/list-03-autoresponder-brak-skrzynki.html` | odpowiedź na wiadomość nadesłaną na adres nieobsługiwany |
| `text/*.txt` | części `text/plain` odpowiadające trzem listom |
| `podglad.html` | podgląd trzech listów z danymi przykładowymi — otwórz w przeglądarce |
| `generator-podgladu.py` | wytwarza `podglad.html` z szablonów; uruchomienie: `python3 generator-podgladu.py` |

Budowa wiadomości: `multipart/alternative` (tekst przed HTML) wewnątrz
`multipart/related` ze znakami `cid:danaco-lockup` i `cid:danaco-lockup-dark`.
`data:` URI występuje wyłącznie w pliku podglądu.

Zmienne w składni `{{ }}` — spis w rozdz. 6 opracowania.
`{{original_subject}}` wymaga sanityzacji (rozdz. 6.5).

Podgląd bierze motyw z ustawienia systemu — przełącz motyw, aby zobaczyć
drugi wariant barwny.
