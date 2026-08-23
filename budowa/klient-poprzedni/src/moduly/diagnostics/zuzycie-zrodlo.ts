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
 * Rodzina `usage.*` widziana przez klienta — zestawienie zużycia
 * (`usage.summary.get`) i raport rozliczeniowy (`usage.report.build`).
 *
 * Zestawienie idzie po JEDNYM wymiarze naraz i tak stanowi kontrakt wprost:
 * „Jeden wymiar na żądanie: zestawienie po dwóch wymiarach naraz jest dwoma
 * żądaniami, nie jednym". Zakładka nie składa więc dwóch odpowiedzi w jedną
 * tabelę — pokazuje wymiar wybrany i mówi, który to.
 *
 * Raport bierze wymiary wykazem, bo jest jednym wytworem za okres, a nie
 * widokiem: puste znaczy wszystkie.
 *
 * Koszt niepełny jest tu osobnym pojęciem, nie zaokrągleniem. Kontrakt niesie
 * `priceCoverage` (udział wywołań objętych cennikiem) i `costWithoutPrice`
 * (wywołania pominięte w koszcie), a pole `cost` opisuje jako puste, gdy
 * cennika kanału nie ma — „puste znaczy brak cennika kanału, nie koszt
 * zerowy". Zakładka przenosi to rozróżnienie na ekran; złożenie kosztu z zer
 * dałoby liczbę wyglądającą na pomiar.
 *
 * Zdarzeń rodzina nie ma, więc nic tu nie nasłuchuje: zużycie odczytuje się na
 * żądanie.
 */

/** Wymiary grupowania wraz z nazwą dla Operatora — pełne nazwy, bez skrótów. */
export const WYMIARY_ZUZYCIA: readonly { kod: UsageDimension; nazwa: string }[] = [
  { kod: UsageDimension.Channel, nazwa: 'Kanał modelu' },
  { kod: UsageDimension.Provider, nazwa: 'Dostawca kanału' },
  { kod: UsageDimension.Account, nazwa: 'Konto dostępowe' },
  { kod: UsageDimension.Session, nazwa: 'Karta sesji' },
  { kod: UsageDimension.Project, nazwa: 'Projekt' },
  { kod: UsageDimension.Environment, nazwa: 'Środowisko' },
  { kod: UsageDimension.Window, nazwa: 'Okno komunikacji' },
];

/** Postaci raportu rozliczeniowego z wyliczenia kontraktu. */
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
        // Pusty wykaz jest poprawną odpowiedzią i sprawdzian tego nie myli
        // z odpowiedzią bez kształtu: wymagana jest sama tablica oraz granice
        // okresu, bo bez nich zakładka nie wie, o czym mówi.
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

/** Nazwa wymiaru dla Operatora; wymiar spoza wykazu wraca własnym kodem. */
export function nazwaWymiaru(kod: UsageDimension): string {
  return WYMIARY_ZUZYCIA.find((pozycja) => pozycja.kod === kod)?.nazwa ?? kod;
}

/** Rozszerzenie pliku raportu — postać OTLP jest JSON-em i tak się nazywa. */
export function rozszerzenieRaportu(postac: TelemetryFormat): string {
  return postac === TelemetryFormat.Csv ? 'csv' : postac === TelemetryFormat.Jsonl ? 'jsonl' : 'json';
}
