import {
  ChangeKind,
  type ConfigEntry,
  type SettingCategory,
  type SettingDefinition,
} from '../../../shared/contract';
import type { Kanal, Wynik } from '../protokol/kanal';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import { adresPoczatkowy, adresWpisu, type AdresUstawienia } from './adres-ustawienia';
import { rozstrzygnij, type PunktWidzenia, type Rozstrzygniecie } from './rozstrzygniecie';
import { utworzZrodloKatalogu } from './zrodlo-katalogu';
import { utworzZrodloWartosci, wpisSpodAdresu } from './zrodlo-wartosci';

/**
 * Stan okna konfiguracji — jedno źródło prawdy dla nawigacji kategorii,
 * formularzy i wskaźników zasięgu.
 *
 * Katalog (kategorie i definicje) oraz wpisy konfiguracji trzymamy razem,
 * ponieważ pole formularza potrzebuje obu naraz: definicja mówi, jaką ma być
 * kontrolką, wpisy mówią, skąd bierze się jej wartość. Dwa równoległe stany
 * dałyby dwie prawdy o tej samej wartości.
 *
 * Żadna ścieżka nie zatrzymuje okna. Rdzeń, który nie odda katalogu, zostawia
 * wykazy puste; okno pokazuje wtedy komunikat i pozostaje otwarte.
 * Zapis, który się nie powiedzie, wraca jako `Wynik` z błędem — bez wyjątku.
 */
export type FazaOdczytu = 'spoczynek' | 'odczyt' | 'gotowe' | 'blad';

export interface StanKonfiguracji {
  /**
   * Faza odczytu katalogu i wpisów.
   *
   * Bez niej pusty katalog znaczy trzy rzeczy naraz: „jeszcze nie pytałem",
   * „pytam" i „rdzeń nie zna ani jednej kategorii". Każdej należy się inny
   * stan okna: nic, wskaźnik odczytu, stan pusty.
   */
  faza(): FazaOdczytu;
  /** Powód ostatniego niepowodzenia odczytu; pusty, gdy odczyt się powiódł. */
  powodNiepowodzenia(): string;
  /** Kategorie w kolejności wyświetlania. */
  kategorie(): readonly SettingCategory[];
  /** Pozycje katalogu należące do jednej kategorii. */
  definicjeKategorii(kategoria: string): SettingDefinition[];
  /** Pozycja katalogu o wskazanym kluczu; `null`, gdy katalog jej nie zna. */
  definicjaKlucza(klucz: string): SettingDefinition | null;
  /** Wartość obowiązująca klucza wraz z pochodzeniem i łańcuchem zapisów. */
  rozstrzygnij(definicja: SettingDefinition): Rozstrzygniecie;
  /** Punkt widzenia, z którego liczone jest dziedziczenie. */
  punkt(): PunktWidzenia;
  /** Zmienia punkt widzenia i ogłasza przeliczenie. */
  ustawPunkt(punkt: PunktWidzenia): void;
  /** Wczytuje katalog i wpisy z rdzenia. */
  odswiez(): Promise<void>;
  /** Zapisuje wartość pod adresem i uzgadnia stan z odpowiedzią rdzenia. */
  zapisz(klucz: string, wartosc: unknown, adres: AdresUstawienia): Promise<Wynik<unknown>>;
  /** Usuwa zapis spod adresu; wraca wartość odziedziczona albo domyślna. */
  przywroc(klucz: string, adres: AdresUstawienia): Promise<Wynik<unknown>>;
  /** Subskrypcja przeliczenia stanu. */
  naZmiane(sluchacz: () => void): void;
  /** Odłącza subskrypcję zdarzeń kanału. */
  rozlacz(): void;
}

export function utworzStanKonfiguracji(kanal: Kanal): StanKonfiguracji {
  // Powody zbierają się z trzech odczytów jednego odświeżenia (kategorie,
  // definicje, wpisy) — Operatorowi należy się zdanie o każdym, który nie
  // dojechał, a nie tylko o pierwszym.
  let powody: string[] = [];
  const zapiszPowod = (powod: string): void => void powody.push(powod);

  const katalog = utworzZrodloKatalogu(kanal, zapiszPowod);
  const wartosci = utworzZrodloWartosci(kanal, zapiszPowod);

  let kategorie: SettingCategory[] = [];
  let definicje: SettingDefinition[] = [];
  let wpisy: ConfigEntry[] = [];
  let punkt: PunktWidzenia = adresPoczatkowy();
  let faza: FazaOdczytu = 'spoczynek';

  // Każdy odczyt wpisów dostaje swój numer pokolenia; wolniejszy, wyprzedzony
  // zmianą punktu widzenia, nie nadpisuje świeższego.
  let pokolenie = 0;

  const sluchacze: Array<() => void> = [];
  const oglos = (): void => {
    for (const sluchacz of [...sluchacze]) sluchacz();
  };

  /**
   * Wykaz bez zapisu o wskazanym kluczu i adresie.
   *
   * Jeden przepis na dwie drogi: przywrócenie własne (przycisk „Przywróć" tego
   * okna) i przywrócenie cudze (zdarzenie `config.changed` z innego okna albo
   * urządzenia) zdejmują zapis dokładnie tak samo. Dwa przepisy rozeszłyby się
   * na wskaźniku pochodzenia: jeden egzemplarz stanu pokazywałby wartość
   * przywróconą, drugi wartość zapisaną i widmowy wiersz łańcucha dziedziczenia.
   */
  const bezZapisu = (klucz: string, adres: AdresUstawienia): ConfigEntry[] =>
    wpisy.filter((istniejacy) => !wpisSpodAdresu(istniejacy, klucz, adres));

  /**
   * Zmiana potwierdzona przez rdzeń: zapis zastępuje wpis o tym samym adresie,
   * usunięcie zdejmuje go z wykazu. Rodzaj pominięty znaczy zapis — tak wchodzi
   * odpowiedź na własną komendę `config.set`.
   */
  const przyjmij = (wpis: ConfigEntry, rodzaj: ChangeKind = ChangeKind.Updated): void => {
    const adres = adresWpisu(wpis);
    const pozostale = bezZapisu(wpis.key, adres);
    wpisy = rodzaj === ChangeKind.Deleted ? pozostale : [...pozostale, wpis];
    oglos();
  };

  const odsubskrybuj: Odsubskrybuj = wartosci.naZmiane(przyjmij);

  // Surowe wpisy łańcucha punktu widzenia. Zmiana punktu dociąga poziom, który
  // do tej pory nie był czytany; bez tego wartość spod okna, sesji czy projektu
  // nie ma pokrycia w stanie, a pochodzenie wskazuje poziom globalny.
  const przeczytajWpisy = async (): Promise<void> => {
    const moje = ++pokolenie;
    powody = [];
    faza = 'odczyt';
    oglos();
    const pobrane = await wartosci.wpisy(punkt);
    if (moje !== pokolenie) return;
    wpisy = pobrane;
    faza = powody.length > 0 ? 'blad' : 'gotowe';
    oglos();
  };

  return {
    faza: () => faza,

    powodNiepowodzenia: () => powody.join(' '),

    kategorie: () => kategorie,

    definicjeKategorii: (kategoria) =>
      definicje.filter((definicja) => definicja.categoryId === kategoria),

    definicjaKlucza: (klucz) =>
      definicje.find((definicja) => definicja.key === klucz) ?? null,

    rozstrzygnij: (definicja) => rozstrzygnij(definicja, wpisy, punkt),

    punkt: () => punkt,

    ustawPunkt(nowy) {
      punkt = nowy;
      // Nowy punkt liczy dziedziczenie z innego poziomu, więc jego surowe wpisy
      // trzeba dociągnąć; przeczytajWpisy ogłasza od razu (nowy punkt względem
      // wpisów już znanych) i ponownie, gdy poziom dojedzie.
      void przeczytajWpisy();
    },

    async odswiez() {
      const moje = ++pokolenie;
      powody = [];
      faza = 'odczyt';
      oglos();
      const [pobraneKategorie, pobraneDefinicje, pobraneWpisy] = await Promise.all([
        katalog.kategorie(),
        katalog.definicje(),
        wartosci.wpisy(punkt),
      ]);
      if (moje !== pokolenie) return;
      kategorie = pobraneKategorie;
      definicje = pobraneDefinicje;
      wpisy = pobraneWpisy;
      faza = powody.length > 0 ? 'blad' : 'gotowe';
      oglos();
    },

    async zapisz(klucz, wartosc, adres) {
      const wynik = await wartosci.zapisz(klucz, wartosc, adres);
      const wpis = wpisOdpowiedzi(wynik.wynik);
      if (wynik.udany && wpis !== null) przyjmij(wpis);
      else oglos();
      return wynik;
    },

    async przywroc(klucz, adres) {
      const wynik = await wartosci.przywroc(klucz, adres);
      if (wynik.udany) wpisy = bezZapisu(klucz, adres);
      oglos();
      return wynik;
    },

    naZmiane: (sluchacz) => void sluchacze.push(sluchacz),

    rozlacz: odsubskrybuj,
  };
}

/** Wpis z odpowiedzi `config.set`; kształt niespodziewany daje `null`. */
function wpisOdpowiedzi(tresc: unknown): ConfigEntry | null {
  if (typeof tresc !== 'object' || tresc === null) return null;
  const wpis = (tresc as { entry?: unknown }).entry;
  if (typeof wpis !== 'object' || wpis === null) return null;
  const kandydat = wpis as ConfigEntry;
  return typeof kandydat.key === 'string' ? kandydat : null;
}
