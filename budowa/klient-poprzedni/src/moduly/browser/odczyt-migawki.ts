import { ErrorCode, type BrowserSnapshot } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import type { ZrodloBrowser } from './zrodlo-browser';

/**
 * Migawka strony wraz z drogą, którą moduł po nią sięga.
 *
 * Jedna odpowiedzialność: treść strony widziana w oknie i przez model, oraz
 * jeden odczyt `browser.snapshot.get`, który ją przynosi. Bez ani jednego
 * elementu widoku — okna pytają o wynik, nie o sposób (wzór: `wykazy-zebranego`).
 * Treść pochodzi z tabeli `migawka_strony`, którą wypełnia `browser.navigate`,
 * więc przeżywa przeładowanie powłoki i wymianę karty.
 *
 * Odmowa `not_found` nie jest błędem okna: rdzeń odpowiada nią, dopóki w oknie
 * nie odbyło się ani jedno przejście, co jest normalnym stanem okna świeżo
 * otwartego. Wtedy zostaje stan `nietknieta` — podgląd pusty, zdanie o powodzie,
 * bez stanu błędu. Każda inna odmowa mówi swoim powodem.
 *
 * Odczyt w toku jest osobnym stanem, nie odmianą pustki: „nie ma czego pokazać"
 * i „pytam" to dwie różne rady — jedna każe działać, druga czekać.
 */
export type StanMigawki = 'nietknieta' | 'odczyt' | 'gotowa' | 'blad';

export interface OdczytMigawki {
  /** Migawka widoczna oknom modułu; `null` znaczy „nie ma czego pokazać". */
  wartosc(): BrowserSnapshot | null;
  /** Stan ostatniego odczytu — rozstrzyga, czy okno pokazuje pustkę czy odmowę. */
  stan(): StanMigawki;
  /** Zdanie o ostatnim odczycie; puste, gdy migawka przyszła bez uwag. */
  powod(): string;
  /** Wchłania migawkę z własnego wywołania albo ze zdarzenia `browser.page.changed`. */
  wchlon(migawka: BrowserSnapshot): void;
  /**
   * Czyta migawkę okna z rdzenia. `zeZrodlem` steruje polem `includeHtml` —
   * źródło strony bywa ciężkie, więc wychodzi tylko wtedy, gdy ktoś o nie prosi.
   */
  zaciagnij(idOkna: string, zeZrodlem?: boolean): Promise<void>;
}

export function utworzOdczytMigawki(zrodlo: ZrodloBrowser, oglos: () => void): OdczytMigawki {
  let migawka: BrowserSnapshot | null = null;
  let stan: StanMigawki = 'nietknieta';
  let powod = 'Migawka strony nie została jeszcze odczytana z rdzenia.';

  /** Wpisuje rozstrzygnięcie odczytu i ogłasza je oknom modułu. */
  function przyjmij(nowyStan: StanMigawki, zdanie: string): void {
    stan = nowyStan;
    powod = zdanie;
    oglos();
  }

  return {
    wartosc: () => migawka,
    stan: () => stan,
    powod: () => powod,

    wchlon(nowa) {
      migawka = nowa;
      // Migawka świeżo wchłonięta unieważnia zdanie o poprzednim odczycie:
      // powód „okno nie ma jeszcze migawki" przestał być prawdą w tej chwili.
      przyjmij('gotowa', '');
    },

    async zaciagnij(idOkna, zeZrodlem = false) {
      if (idOkna === '') {
        przyjmij('nietknieta', 'Migawka nie ma o co pytać: okno przeglądarki nie jest ustalone.');
        return;
      }
      przyjmij('odczyt', 'Odczyt treści strony z rdzenia w toku…');
      const wynik = await zrodlo.migawka({ windowId: idOkna, includeHtml: zeZrodlem });
      if (wynik.udany && wynik.wynik !== undefined) {
        migawka = wynik.wynik.snapshot;
        przyjmij('gotowa', '');
        return;
      }
      // Migawki nie kasujemy przy odmowie: nieudany odczyt nie unieważnia
      // strony, którą Operator już czyta.
      if (wynik.blad?.code === ErrorCode.NotFound) {
        przyjmij(
          'nietknieta',
          'Rdzeń nie ma jeszcze migawki tego okna — treść strony pojawi się po pierwszym przejściu pod adres.',
        );
        return;
      }
      przyjmij(
        'blad',
        opisOdmowy('Odczyt migawki strony', wynik.blad?.code, wynik.blad?.message),
      );
    },
  };
}
