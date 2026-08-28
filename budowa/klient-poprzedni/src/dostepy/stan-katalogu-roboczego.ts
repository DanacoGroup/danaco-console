import type { ConfigEntry, SettingDefinition } from '../../../shared/contract';
import { adresPoczatkowy, adresWpisu, tenSamAdres } from '../konfiguracja/adres-ustawienia';
import { rozstrzygnij, type Rozstrzygniecie } from '../konfiguracja/rozstrzygniecie';
import { utworzZrodloKatalogu } from '../konfiguracja/zrodlo-katalogu';
import { utworzZrodloWartosci } from '../konfiguracja/zrodlo-wartosci';
import type { Kanal, Wynik } from '../protokol/kanal';
import { definicjaZastepcza, KLUCZE_KATALOGU } from './klucze-katalogu';

/**
 * Stan obszaru katalogu roboczego. Obszar nie ma własnej drogi do rdzenia:
 * bierze definicje z katalogu ustawień i wartości komendą `config.get`,
 * a rachunek dziedziczenia oddaje modułowi `konfiguracja/rozstrzygniecie`.
 */
export interface StanKataloguRoboczego {
  /** Definicja pozycji katalogu; brak w katalogu daje definicję zastępczą. */
  definicja(klucz: string): SettingDefinition | null;
  /** Wartość obowiązująca wraz z pochodzeniem; `null` bez definicji. */
  rozstrzygnij(klucz: string): Rozstrzygniecie | null;
  /** Wczytuje definicje i wpisy z rdzenia. */
  odswiez(): Promise<void>;
  /** Zapisuje wartość klucza na poziomie globalnym, na osi platformy. */
  zapisz(klucz: string, wartosc: unknown): Promise<Wynik<unknown>>;
  /** Zdejmuje zapis globalny; wraca wartość domyślna rdzenia. */
  przywroc(klucz: string): Promise<Wynik<unknown>>;
  /** Subskrypcja przeliczenia. */
  naZmiane(sluchacz: () => void): void;
  /** Odłącza subskrypcję `config.changed`. */
  rozlacz(): void;
}

export function utworzStanKataloguRoboczego(kanal: Kanal): StanKataloguRoboczego {
  const katalog = utworzZrodloKatalogu(kanal);
  const wartosci = utworzZrodloWartosci(kanal);
  const adres = adresPoczatkowy();

  let definicje: SettingDefinition[] = [];
  let wpisy: ConfigEntry[] = [];

  const sluchacze: Array<() => void> = [];
  const oglos = (): void => {
    for (const sluchacz of [...sluchacze]) sluchacz();
  };

  /** Wpis potwierdzony przez rdzeń zastępuje zapis o tym samym adresie. */
  const przyjmij = (wpis: ConfigEntry): void => {
    if (!KLUCZE_KATALOGU.includes(wpis.key)) return;
    const adresWpisanego = adresWpisu(wpis);
    wpisy = [
      ...wpisy.filter(
        (istniejacy) =>
          istniejacy.key !== wpis.key || !tenSamAdres(adresWpisu(istniejacy), adresWpisanego),
      ),
      wpis,
    ];
    oglos();
  };

  const odsubskrybuj = wartosci.naZmiane(przyjmij);

  function definicjaKlucza(klucz: string): SettingDefinition | null {
    return definicje.find((pozycja) => pozycja.key === klucz) ?? definicjaZastepcza(klucz);
  }

  return {
    definicja: definicjaKlucza,

    rozstrzygnij(klucz) {
      const definicja = definicjaKlucza(klucz);
      return definicja === null ? null : rozstrzygnij(definicja, wpisy, adres);
    },

    async odswiez() {
      const [pobraneDefinicje, pobraneWpisy] = await Promise.all([
        katalog.definicje(),
        wartosci.wpisy(),
      ]);
      definicje = pobraneDefinicje.filter((pozycja) => KLUCZE_KATALOGU.includes(pozycja.key));
      wpisy = pobraneWpisy.filter((wpis) => KLUCZE_KATALOGU.includes(wpis.key));
      oglos();
    },

    async zapisz(klucz, wartosc) {
      const wynik = await wartosci.zapisz(klucz, wartosc, adres);
      const wpis = wpisOdpowiedzi(wynik.wynik);
      if (wynik.udany && wpis !== null) przyjmij(wpis);
      else oglos();
      return wynik;
    },

    async przywroc(klucz) {
      const wynik = await wartosci.przywroc(klucz, adres);
      if (wynik.udany) {
        wpisy = wpisy.filter(
          (istniejacy) =>
            istniejacy.key !== klucz || !tenSamAdres(adresWpisu(istniejacy), adres),
        );
      }
      oglos();
      return wynik;
    },

    naZmiane: (sluchacz) => void sluchacze.push(sluchacz),

    rozlacz: odsubskrybuj,
  };
}

/**
 * Wpis z odpowiedzi `config.set`; kształt niespodziewany daje `null`, dzięki
 * czemu stan nie przyjmuje wartości, której rdzeń nie potwierdził.
 */
function wpisOdpowiedzi(tresc: unknown): ConfigEntry | null {
  if (typeof tresc !== 'object' || tresc === null) return null;
  const wpis = (tresc as { entry?: unknown }).entry;
  if (typeof wpis !== 'object' || wpis === null) return null;
  const kandydat = wpis as ConfigEntry;
  return typeof kandydat.key === 'string' ? kandydat : null;
}
