import type { Channel } from '../../../../shared/contract';

/**
 * Rozpoznanie kanału obrazowego po stronie klienta: czy danym kanałem da się
 * wygenerować obraz.
 *
 * Żądanie `design.asset.generate` niesie pole `channelId`, a rdzeń sprawdza
 * rodzaj wskazanego kanału zanim cokolwiek wyśle — kanał tekstowy odmawia kodem
 * `validation_failed` (`core/adapter_modul_design_kanal.go`). Wykaz silników
 * pokazuje więc kanały obrazowe jako wybieralne, a pozostałe wymienia z nazwy
 * i z powodem, zamiast je ukrywać.
 *
 * Kolejność rozpoznania jest przepisana z `models/definicja.go`,
 * `KluczAdaptera()`: parametr `adapter` ma pierwszeństwo, a rodzaj wiersza jest
 * wartością zapasową. Do kontraktu obie strony jadą tym samym wierszem
 * (`core/adapter_kanaly.go`: `Kind: k.RodzajKanalu`, `Config: ParametryJSON`),
 * więc `Channel.kind` to `Rodzaj`, a `Channel.config.adapter` to `Parametr`.
 * Kluczem adaptera obrazowego jest napis `obrazy` (`models/kanal.go`,
 * `AdapterObrazy`).
 *
 * Rozpoznanie może rozejść się z rdzeniem: klient czyta `config` jako treść
 * nieokreśloną (`config?: unknown` kontraktu), więc kanał o parametrach
 * zapisanych inaczej niż obiektem trafi między nieobrazowe. Ostateczną kontrolę
 * wykonuje rdzeń — dlatego wskazanie kanału spoza wykazu obrazowego nie jest
 * w oknie blokowane, tylko opisane.
 */

/** Klucz adaptera kanału obrazowego — `models/kanal.go`, `AdapterObrazy`. */
export const ADAPTER_OBRAZY = 'obrazy';

/**
 * Klucz adaptera kanału — parametr `adapter`, a w jego braku rodzaj wiersza.
 *
 * Pusty wynik znaczy „kanał nie mówi, czym jest": ani parametru, ani rodzaju.
 * To inny stan niż „kanał tekstowy" i wykaz mówi o nim osobno.
 */
export function kluczAdaptera(kanal: Channel): string {
  const parametr = parametrKonfiguracji(kanal.config, 'adapter');
  if (parametr !== '') return parametr;
  return kanal.kind.trim();
}

/** Czy kanał oddaje bajty obrazu — jedyny rodzaj, którym generowanie przejdzie. */
export function czyKanalObrazowy(kanal: Channel): boolean {
  return kluczAdaptera(kanal) === ADAPTER_OBRAZY;
}

/** Kanały obrazowe rejestru w kolejności, w jakiej oddał je rdzeń. */
export function kanalyObrazowe(kanaly: readonly Channel[]): readonly Channel[] {
  return kanaly.filter(czyKanalObrazowy);
}

/**
 * Dlaczego kanał nie nadaje się na silnik obrazów — zdanie dla wykazu.
 *
 * Puste znaczy „nadaje się". Zdanie powtarza powód, którym odmówiłby rdzeń, więc
 * powód jest znany bez zlecania generowania.
 */
export function powodNieprzydatnosci(kanal: Channel): string {
  if (czyKanalObrazowy(kanal)) return '';
  const klucz = kluczAdaptera(kanal);
  if (klucz === '') {
    return 'kanał nie podaje ani parametru adapter, ani rodzaju — rdzeń nie ma po czym poznać, czy odda obraz';
  }
  return `adapter „${klucz}" oddaje fragmenty tekstu, nie bajty obrazu`;
}

/**
 * Odczyt parametru kanału z pola `config` kontraktu.
 *
 * Pole jest w kontrakcie treścią nieokreśloną (`config?: unknown`), bo niesie
 * parametry dowolnego adaptera. Wszystko, co nie jest obiektem z napisem pod tym
 * kluczem, daje pustkę.
 */
function parametrKonfiguracji(config: unknown, klucz: string): string {
  if (config === null || typeof config !== 'object') return '';
  const wartosc = (config as Record<string, unknown>)[klucz];
  return typeof wartosc === 'string' ? wartosc.trim() : '';
}
