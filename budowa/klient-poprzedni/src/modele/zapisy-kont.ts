import type {
  Account,
  AccountAddRequest,
  AccountDefaultSetRequest,
  AccountUpdateRequest,
} from '../../../shared/contract';
import type { Wynik } from '../protokol/kanal';
import type { ZrodloKont } from './zrodlo-kont';

/**
 * Cztery komendy zmieniające rejestr kont wraz z uzgodnieniem stanu widoku
 * z odpowiedzią rdzenia. Żadna ścieżka nie odrzuca obietnicy: niepowodzenie
 * wraca jako `Wynik` z błędem, aby formularz mógł podać treść odpowiedzi
 * rdzenia.
 */
export interface ZapisyKont {
  /** `account.add` — zakłada konto i czyni je czynnym. */
  dodaj(zadanie: AccountAddRequest): Promise<Wynik<unknown>>;
  /** `account.update` — zmienia konto wskazane w żądaniu. */
  zmien(zadanie: AccountUpdateRequest): Promise<Wynik<unknown>>;
  /** `account.remove` — usuwa konto i zdejmuje jego wybór. */
  usun(identyfikator: string): Promise<Wynik<unknown>>;
  /** `account.default.set` — wskazuje konto domyślne swojego rodzaju. */
  ustawDomyslne(zadanie: AccountDefaultSetRequest): Promise<Wynik<unknown>>;
}

/**
 * Wejścia zapisów: źródło komend rdzenia oraz cztery czynności stanu rejestru
 * — przyjęcie potwierdzonego konta, odłączenie konta, wybór konta czynnego
 * i ponowny odczyt wykazu.
 */
export interface ZaleznosciZapisow {
  zrodlo: ZrodloKont;
  /** Nanosi konto potwierdzone przez rdzeń na wykaz. */
  przyjmij(konto: Account): void;
  /** Zdejmuje konto z wykazu wraz z jego wyborem. */
  odlacz(identyfikator: string): void;
  /** Czyni konto czynnym — po założeniu nowego. */
  wybierz(identyfikator: string): void;
  /** Czyta wykaz od nowa, gdy odpowiedź nie wystarcza do uzgodnienia. */
  wczytaj(): Promise<void>;
}

export function utworzZapisyKont(zaleznosci: ZaleznosciZapisow): ZapisyKont {
  const { zrodlo, przyjmij, odlacz, wybierz, wczytaj } = zaleznosci;

  return {
    async dodaj(zadanie) {
      const wynik = await zrodlo.dodaj(zadanie);
      const konto = wynik.wynik?.account;
      if (wynik.udany && konto !== undefined) {
        przyjmij(konto);
        wybierz(konto.id);
      }
      return wynik;
    },

    async zmien(zadanie) {
      const wynik = await zrodlo.zmien(zadanie);
      const konto = wynik.wynik?.account;
      if (wynik.udany && konto !== undefined) przyjmij(konto);
      return wynik;
    },

    async usun(identyfikator) {
      const wynik = await zrodlo.usun({ accountId: identyfikator });
      if (wynik.udany) odlacz(identyfikator);
      return wynik;
    },

    async ustawDomyslne(zadanie) {
      const wynik = await zrodlo.ustawDomyslne(zadanie);
      // Oznaczenie domyślnego przestawia dwa wiersze naraz, a odpowiedź niesie
      // tylko jeden.
      if (wynik.udany) await wczytaj();
      return wynik;
    },
  };
}
