# 06 · Poczta transakcyjna

Siedem szablonów wiadomości wychodzących z adresu `noreply@danaco-group.pl`
w ramach Danaco Console, wraz z pakietem składającym wiadomość i runbookiem
wdrożenia. Zgodne z księgą znaku, typografią i głosem oraz zastosowaniami
marki (rozdz. 4.6).

Każdy list ma dwie części: `text/html` i `text/plain`.
Wariant jasny jest podstawą, ciemny — zadeklarowanym nadpisaniem.

## Dokumenty

| Plik | Zawartość |
|---|---|
| `opracowanie-techniczne.md` | decyzje projektowe, żetony, ograniczenia klientów pocztowych, wariant ciemny, część tekstowa, zmienne, bezpieczeństwo, pakiet rdzenia, lista kontrolna |
| `runbook-wdrozenia.md` | uruchomienie wysyłki na Stalwart 0.16.x — SPF, DKIM, DMARC, autoresponder, limity, dziennik, kontrola po wdrożeniu |

## Szablony

| Nr | Plik (bez rozszerzenia) | Okoliczność |
|---:|---|---|
| 1 | `list-01-aktywacja-konta` | kod aktywacji nowego konta |
| 2 | `list-02-uwierzytelnienie-logowania` | kod logowania |
| 3 | `list-03-autoresponder-brak-skrzynki` | odpowiedź na wiadomość nadesłaną na adres nieobsługiwany |
| 4 | `list-04-reset-hasla` | kod resetu hasła |
| 5 | `list-05-zmiana-adresu-potwierdzenie` | kod na **nowy** adres konta |
| 6 | `list-06-zmiana-adresu-powiadomienie` | ostrzeżenie na **dotychczasowy** adres konta |
| 7 | `list-07-przebieg-zakonczony` | przebieg automatyki zakończony albo zatrzymany |

Część HTML w `html/`, część tekstowa w `text/`. Listy 5 i 6 wychodzą razem,
na dwa różne adresy.

## Kod i podgląd

| Plik | Zawartość |
|---|---|
| `rdzen/mail/` | pakiet Go składający wiadomość — do przeniesienia do repozytorium rdzenia |
| `podglad.html` | podgląd siedmiu listów z danymi przykładowymi — otwórz w przeglądarce |
| `generator-podgladu.py` | wytwarza `podglad.html`; uruchomienie: `python3 generator-podgladu.py` |

Testy pakietu:

```bash
cd rdzen/mail && go test ./...
```

Testy czytają szablony wprost z katalogów `html/` i `text/`, więc rozjazd
między szablonem a kodem wychodzi przy `go test`, nie przy wysyłce.

## Budowa wiadomości

```
multipart/related
├── multipart/alternative
│   ├── text/plain        (pierwsza część — klient wybiera ostatnią, którą umie)
│   └── text/html
├── image/png  Content-ID: <danaco-lockup>
└── image/png  Content-ID: <danaco-lockup-dark>
```

`data:` URI występuje wyłącznie w pliku podglądu.

Zmienne w składni `{{ }}` — spis w rozdz. 6 opracowania.
`{{original_subject}}` i `{{run_name}}` pochodzą spoza rdzenia i wymagają
sanityzacji (rozdz. 6.6).

Podgląd bierze motyw z ustawienia systemu — przełącz motyw, aby zobaczyć
drugi wariant barwny.
