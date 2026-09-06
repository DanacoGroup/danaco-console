# Danaco Console

Aplikacja hybrydowa do pracy z modelami AI: rdzeń w Go, interfejs w TypeScripcie,
powłoka desktopowa na Tauri 2. Wydawana jako instalator dla Windows i Linux.

## Budowa repozytorium

| Katalog | Zawartość |
| --- | --- |
| `budowa/server` | Rdzeń w Go — logika, magazyn SQLite, migracje schematu |
| `budowa/klient` | Interfejs w TypeScripcie (Vite) |
| `budowa/desktop` | Powłoka Tauri 2 pakująca rdzeń i interfejs |
| `budowa/instalator` | Instalator wydania |
| `budowa/shared` | Kontrakt komunikacji i wytworzone z niego wiązania |
| `budowa/witryna` | Witryna produktu |
| `design` | Warstwa projektowa: marka, okna, zasoby |
| `docs` | Dokumentacja: architektura, moduły, instalacja, licencja |
| `narzedzia` | Drabina weryfikacji i walidatory dyscypliny |

## Kontrakt jako źródło prawdy

Nazwy komend, zdarzeń, kodów błędów i kształt koperty żyją wyłącznie
w `budowa/shared/contract.json`. Wiązania dla obu stron powstają z niego
generatorem — ręczna zmiana wytworu jest błędem:

```bash
node budowa/shared/gen/generate.mjs
```

Generator wytwarza `contract.go`, `contract.ts` oraz rejestr komend klienta.
Świeżości tych wytworów pilnuje sprawdzian `budowa/shared/swiezosc_generatu_test.go`
oraz szczebel drabiny weryfikacji.

## Weryfikacja

```bash
bash narzedzia/drabina.sh szybka
```

Tryb `szybka` obejmuje format, budowę, dyscyplinę kodu i świeżość kontraktu;
tryb `pelna` dokłada sprawdziany rdzenia i klienta. Drabinę uruchamia też hook
`pre-commit` z katalogu `.githooks`, włączany poleceniem:

```bash
git config core.hooksPath .githooks
```

## Licencja

Oprogramowanie własnościowe. Warunki korzystania określa [docs/LICENSE.md](docs/LICENSE.md).
Publikacja kodu w tym repozytorium nie przenosi żadnych praw ani nie udziela
licencji na użycie, zwielokrotnianie czy tworzenie opracowań.

---

© Danaco Holding Group. Wszelkie prawa zastrzeżone.
