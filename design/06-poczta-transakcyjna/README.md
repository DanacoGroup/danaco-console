# 06 · Poczta transakcyjna

Szablony wiadomości wychodzących z adresu `noreply@danaco-group.pl`
w ramach Danaco Console. Wariant jasny, zgodny z księgą znaku,
typografią i głosem oraz zastosowaniami marki (rozdz. 4.6).

| Plik | Zawartość |
|---|---|
| `opracowanie-techniczne.md` | decyzje projektowe, żetony, ograniczenia klientów pocztowych, zmienne szablonu, bezpieczeństwo, lista kontrolna |
| `html/list-01-aktywacja-konta.html` | kod aktywacji nowego konta Operatora |
| `html/list-02-uwierzytelnienie-logowania.html` | kod logowania (drugi składnik) |
| `html/list-03-autoresponder-brak-skrzynki.html` | odpowiedź na wiadomość nadesłaną na adres nieobsługiwany |
| `podglad.html` | podgląd trzech listów z danymi przykładowymi — otwórz w przeglądarce |
| `generator-podgladu.py` | wytwarza `podglad.html` z szablonów; uruchomienie: `python3 generator-podgladu.py` |

Szablony zawierają zmienne w składni `{{ }}` — spis w rozdz. 4 opracowania.
Znak w wysyłce produkcyjnej dołączany jako `cid:danaco-lockup`;
`data:` URI występuje wyłącznie w pliku podglądu.
