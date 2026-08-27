/**
 * Rodzina komend `usage.*` widziana przez klienta: zestawienie zużycia oraz
 * raport rozliczeniowy. Zestawienie idzie po jednym wymiarze naraz, a raport
 * bierze wymiary wykazem, ponieważ jest jednym wytworem za okres.
 */
import {
  Command,
  TelemetryFormat,
  UsageDimension,
  type UsageReportBuildRequest,
  type UsageReportBuildResponse,
  type UsageSummaryGetRequest,
  type UsageSummaryGetResponse,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyLiczba, czyTablica, czyTekst, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Wymiary grupowania zużycia wraz z nazwą przeznaczoną dla Operatora. Nazwy są
 * pełne, bez skrótów, ponieważ wymiar czytany na ekranie ma mówić to samo, co
 * wymiar nazwany w kontrakcie.
 */
export const WYMIARY_ZUZYCIA: readonly { kod: UsageDimension; nazwa: string }[] = [
  { kod: UsageDimension.Channel, nazwa: 'Kanał modelu' },
  { kod: UsageDimension.Provider, nazwa: 'Dostawca kanału' },
  { kod: UsageDimension.Account, nazwa: 'Konto dostępowe' },
  { kod: UsageDimension.Session, nazwa: 'Karta sesji' },
  { kod: UsageDimension.Project, nazwa: 'Projekt' },
  { kod: UsageDimension.Environment, nazwa: 'Środowisko' },
  { kod: UsageDimension.Window, nazwa: 'Okno komunikacji' },
];

/**
 * Postaci raportu rozliczeniowego wzięte z wyliczenia kontraktu wraz z nazwami
 * pokazywanymi Operatorowi. Wykaz jest zamknięty, więc postać spoza niego nie
 * powstaje w oknie ani nie idzie do rdzenia.
 */
export const POSTACI_RAPORTU: readonly { kod: TelemetryFormat; nazwa: string }[] = [
  { kod: TelemetryFormat.Csv, nazwa: 'Wartości rozdzielone przecinkiem (CSV)' },
  { kod: TelemetryFormat.Json, nazwa: 'JSON platformy' },
  { kod: TelemetryFormat.Jsonl, nazwa: 'JSON Lines' },
  { kod: TelemetryFormat.Otlp, nazwa: 'OpenTelemetry Protocol w postaci JSON' },
];

export interface ZrodloZuzycia {
  /** `usage.summary.get` — zużycie zebrane w jednym wymiarze. */
  zestawienie(zadanie: UsageSummaryGetRequest): Promise<Wynik<UsageSummaryGetResponse>>;
  /** `usage.report.build` — raport rozliczeniowy za okres; wytworem jest treść. */
  raport(zadanie: UsageReportBuildRequest): Promise<Wynik<UsageReportBuildResponse>>;
}

export function utworzZrodloZuzycia(kanal: Kanal): ZrodloZuzycia {
  return {
    async zestawienie(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.UsageSummaryGet, zadanie),
        Command.UsageSummaryGet,
        // Pusty wykaz jest poprawną odpowiedzią; wymagane są granice okresu.
        (tresc) =>
          czyTablica(tresc.aggregates) && czyLiczba(tresc.fromTime) && czyLiczba(tresc.toTime),
      );
    },

    async raport(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.UsageReportBuild, zadanie),
        Command.UsageReportBuild,
        (tresc) => czyTekst(tresc.content) && czyLiczba(tresc.generatedAt),
      );
    },
  };
}

/**
 * Nazwa wymiaru pokazywana Operatorowi; wymiar spoza wykazu wraca własnym
 * kodem kontraktu, ponieważ nazwa zmyślona po stronie okna byłaby słowem, za
 * którym nie stoi żaden zapis rdzenia.
 */
export function nazwaWymiaru(kod: UsageDimension): string {
  return WYMIARY_ZUZYCIA.find((pozycja) => pozycja.kod === kod)?.nazwa ?? kod;
}

/**
 * Rozszerzenie pliku raportu dobrane do jego postaci. Postać zapisu śladów jest
 * zapisem JSON i tak też nazywa się jej plik, żeby Operator otwierał go
 * narzędziem właściwym dla treści.
 */
export function rozszerzenieRaportu(postac: TelemetryFormat): string {
  return postac === TelemetryFormat.Csv ? 'csv' : postac === TelemetryFormat.Jsonl ? 'jsonl' : 'json';
}
