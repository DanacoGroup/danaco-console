/**
 * Słownik rodzajów konta pokazywany operatorowi. Rodzaj jest wyliczeniem
 * kontraktu, więc jego wartości pochodzą wyłącznie ze stałych `AccountKind`;
 * rodzaj spoza kontraktu pokazuje własny kod zamiast gasnąć.
 */
import { AccountKind } from '../../../shared/contract';



/**
 * Trzy rodzaje konta w kolejności wyświetlania w wykazie i w kontrolce wyboru:
 * konto API modelu, konto programu code CLI oraz konto SDK.
 */
export const RODZAJE_KONT: readonly AccountKind[] = [
  AccountKind.Api,
  AccountKind.Cli,
  AccountKind.Sdk,
];

/**
 * Nazwa rodzaju konta pokazywana operatorowi w formularzu, pełna dla każdej
 * z trzech wartości wyliczenia `AccountKind` niesionego przez kontrakt.
 */
export const NAZWY_RODZAJOW: Readonly<Record<AccountKind, string>> = {
  [AccountKind.Api]: 'konto API modelu',
  [AccountKind.Cli]: 'konto programu code CLI',
  [AccountKind.Sdk]: 'konto SDK',
};

/**
 * Nazwa krótka rodzaju konta, przeznaczona do plakietki w wykazie kont, gdzie
 * nazwa pełna nie mieści się w szerokości wiersza.
 */
export const NAZWY_KROTKIE: Readonly<Record<AccountKind, string>> = {
  [AccountKind.Api]: 'API',
  [AccountKind.Cli]: 'code CLI',
  [AccountKind.Sdk]: 'SDK',
};

/**
 * Nazwa rodzaju gotowa do wydruku; rodzaj spoza kontraktu pokazuje własny kod,
 * żeby było widać, co przyszło z rdzenia, zamiast pustego miejsca.
 */
export function nazwaRodzaju(rodzaj: string): string {
  return NAZWY_RODZAJOW[rodzaj as AccountKind] ?? rodzaj;
}

/**
 * Nazwa krótka gotowa do wydruku; rodzaj spoza kontraktu pokazuje własny kod,
 * ponieważ plakietka ma nazwać wartość oddaną przez rdzeń.
 */
export function nazwaKrotka(rodzaj: string): string {
  return NAZWY_KROTKIE[rodzaj as AccountKind] ?? rodzaj;
}

/**
 * Czy napis jest rodzajem znanym kontraktowi; sprawdzenie zawęża typ napisu do
 * wyliczenia `AccountKind` i strzeże odczytu obu słowników nazw.
 */
export function czyRodzajKonta(tekst: string): tekst is AccountKind {
  return (RODZAJE_KONT as readonly string[]).includes(tekst);
}

/**
 * Rodzaj odczytany z wartości kontrolki wyboru. Napis pusty znaczy „bez
 * ograniczenia"; napis spoza kontraktu również, bo ograniczenie wykazu do
 * wartości, której rdzeń nie zna, dałoby wykaz pusty bez powodu.
 */
export function rodzajZWyboru(wartosc: string): AccountKind | null {
  return czyRodzajKonta(wartosc) ? wartosc : null;
}

/**
 * Czy rodzaj konta korzysta z pola `configDir`. Katalog konfiguracji ma sens
 * wyłącznie dla programu code CLI — pozostałe rodzaje łączą się adresem punktu
 * końcowego.
 */
export function rodzajUzywaKataloguKonfiguracji(rodzaj: string): boolean {
  return rodzaj === AccountKind.Cli;
}
