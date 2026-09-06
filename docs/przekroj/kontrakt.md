# Danaco Console — przekrój pionowy: kontrakt

Opracowanie opisuje warstwę kontraktu w postaci, w jakiej działa w bieżącej
budowie. Obejmuje źródło prawdy typów, kopertę komunikatu, zawartość kontraktu,
kody błędów oraz wytwory generatora wraz z ich kontrolą.

## Źródło prawdy

Jedynym źródłem prawdy nazw komend, zdarzeń, kodów błędów i kształtu koperty
jest plik `budowa/shared/contract.json`. Kontrakt nosi oznaczenie protokołu
`1.0` i wskazuje transport WebSocket z treścią w formacie JSON. Nazwy komend
i zdarzeń zapisane są po angielsku w notacji `obszar.zasob.akcja`
z separatorem kropki. Powielenie literału nazwy po którejkolwiek stronie —
w rdzeniu albo w kliencie — jest błędem: obie strony importują wyłącznie
wytwory generatora.

## Koperta komunikatu

Jedna koperta obsługuje zadanie, odpowiedź i fragment strumienia. Pola wspólne:

- `type` — typ komunikatu w notacji `obszar.zasob.akcja`, typ `MessageType`,
  wymagane;
- `id` — identyfikator zadania; odpowiedź i fragmenty strumienia powtarzają
  identyfikator zadania, przez co klient wiąże je z wywołaniem, wymagane;
- `sessionId` — sesja, której dotyczy komunikat; puste dla `connection.hello`;
- `payload` — treść właściwa, której kształt wyznacza typ komunikatu;
- `timestamp` — czas nadania w milisekundach epoki, typ `int64`, wymagane.

Odpowiedź niesie ponadto pole `status` typu `EnvelopeStatus`, a przy statusie
`error` — pole `error` typu `ErrorInfo`.

## Zawartość kontraktu

Kontrakt wymienia 69 obszarów, 1089 komend, 75 zdarzeń, 367 wyliczeń,
549 struktur oraz 466 pozycji narzędzi z przedrostkiem `danaco`. Wyliczenie
mające odpowiednik w schemacie SQLite niesie pole `kolumnaBazy` w zapisie
`tabela.kolumna` oraz pole `baza` przy każdej wartości; kontrakt zapisuje nazwy
po angielsku, a odwzorowanie wiąże je z modelem danych.

## Kody błędów

Kontrakt zna osiem kodów błędów. Kody `validation_failed`, `not_found`,
`not_authenticated`, `permission_denied` i `conflict` oznaczają odmowę bez
ponawiania. Kody `channel_unavailable`, `rate_limited` i `internal_error`
dopuszczają ponowienie wywołania.

## Wytwory generatora i ich kontrola

Generator w `budowa/shared/gen` wytwarza z kontraktu dwa pliki powiązań:
`budowa/shared/contract.go` dla rdzenia oraz `budowa/shared/contract.ts` dla
klienta. Świeżość wytworów wobec `contract.json` mierzy sprawdzian
`budowa/shared/swiezosc_generatu_test.go`; tę samą miarę przechodzi drabina
weryfikacji `narzedzia/drabina.sh`, która dodatkowo liczy pokrycie kontraktu
przez interfejs. Bieżący wynik obu miar: wytwory świeże, interfejs woła
1089 z 1089 komend kontraktu.
