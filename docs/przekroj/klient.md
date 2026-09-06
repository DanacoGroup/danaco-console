# Danaco Console — przekrój pionowy: klient

Opracowanie opisuje klienta interfejsu w postaci, w jakiej działa w bieżącej
budowie. Obejmuje układ warstw, drogę wywołania, składanie wnętrz okien oraz
miary zgodności z kontraktem.

## Układ warstw

Klient jest programem TypeScript budowanym narzędziem Vite; źródła leżą
w `budowa/klient/src/`:

- `polaczenie/` — gniazdo WebSocket i jego otoczenie: adres rdzenia, stan
  połączenia, ponawianie, kolejka wychodząca, magistrala i rozdzielacz zdarzeń,
  dziennik komunikatów nieznanych;
- `protokol/` — koperta, korelacja odpowiedzi po identyfikatorze zadania,
  powitanie `connection.hello`, sesja i token sesji, tożsamość klienta,
  wywołanie komendy i kształt odpowiedzi;
- `wiazanie/` — wiązania okien: każda część interfejsu czyta stan z rdzenia
  i wpina czynności użytkownika w komendy kontraktu;
- `model/` — modele pomocnicze interfejsu;
- `aplikacja.ts` — złożenie całości.

Nazwy komend i zdarzeń pochodzą wyłącznie z wytworu generatora
`budowa/shared/contract.ts`.

## Droga wywołania

Wywołanie z okna przechodzi przez `protokol/wywolanie.ts`: otrzymuje kopertę
z identyfikatorem zadania i idzie do gniazda, a odpowiedź wraca po korelacji
identyfikatora; w czasie rozłączenia ramki czekają w kolejce wychodzącej
i wychodzą po odzyskaniu połączenia. Zdarzenia rdzenia trafiają z gniazda do
rozdzielacza i dalej magistralą do wiązań, które odświeżają widok. Komunikat
spoza kontraktu ląduje w dzienniku nieznanych zamiast przerywać pracę.

## Składanie wnętrz okien

Wnętrza okien pochodzą z prototypów warstwy projektowej. Budowanie
(`budowa/klient/vite.config.js`) wypełnia szablony `dn-tresc-*` treścią
prototypów z `design/05-okna/` i kopiuje skrypty biblioteki projektowej do
wyniku; klient dokłada do tych wnętrz wiązania danych i czynności.

## Weryfikacja

Klient przechodzi kontrolę typów `npm run typy` oraz budowanie
`npm run budowanie`. Pokrycie kontraktu mierzy drabina weryfikacji
`narzedzia/drabina.sh`: interfejs woła 1089 z 1089 komend kontraktu.
