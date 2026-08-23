import { ConfigScope } from '../../../../shared/contract';
import type { StanMultitaskingu } from './stan-multitaskingu';
import type { ZrodloOkien } from './zrodlo-okien';

/**
 * Cztery tryby współpracy wykonawców okien Executor Chat.
 *
 * Tryb rozstrzyga jedno przekazanie, a nie prowadzi biegu: mówi wyłącznie, który
 * wykonawca dostaje zlecenie i co do niego jedzie. Automatyczne podawanie wyniku
 * dalej po każdej turze byłoby drugim silnikiem pętli obok pętli sesji
 * w rdzeniu — moduł go nie buduje.
 *
 * Wartość utrwala `config.set` na poziomie okna koordynatora. Rdzeń nie
 * rejestruje `role.update`, więc profil roli po jego stronie się nie zmienia;
 * pokrycie kontraktu opisuje `braki-kontraktu.ts`.
 */

export const TrybWspolpracy = {
  /** Obaj wykonawcy dostają to samo zlecenie i pracują osobno. */
  Niezalezna: 'niezalezna',
  /** Zlecenie idzie do drugiego wykonawcy wraz z wynikiem pierwszego. */
  Przekazywanie: 'przekazywanie',
  /** Zlecenie idzie do tego wykonawcy, który poprzedniego nie dostał. */
  Naprzemienna: 'naprzemienna',
  /** Zlecenie wraca do wykonawcy poprzedniego wraz z jego własnym wynikiem. */
  Iteracyjna: 'iteracyjna',
} as const;
export type TrybWspolpracy = (typeof TrybWspolpracy)[keyof typeof TrybWspolpracy];

/** Klucz utrwalenia trybu na poziomie okna koordynatora. */
export const KLUCZ_TRYBU = 'multitasking.tryb_wspolpracy';

/** Nazwy trybów w selektorze okna wykonawcy. */
export const NAZWY_TRYBOW: ReadonlyArray<[TrybWspolpracy, string]> = [
  [TrybWspolpracy.Niezalezna, 'Praca niezależna'],
  [TrybWspolpracy.Przekazywanie, 'Przekazywanie wyników'],
  [TrybWspolpracy.Naprzemienna, 'Praca naprzemienna'],
  [TrybWspolpracy.Iteracyjna, 'Praca iteracyjna'],
];

/** Czy wartość jest jednym z czterech trybów; nieznana schodzi na niezależną. */
export function trybZWartosci(wartosc: unknown): TrybWspolpracy {
  return NAZWY_TRYBOW.some(([tryb]) => tryb === wartosc)
    ? (wartosc as TrybWspolpracy)
    : TrybWspolpracy.Niezalezna;
}

/** Jedno zlecenie skierowane do jednego okna wykonawcy. */
export interface Zlecenie {
  /** Okno wykonawcy, do którego idzie polecenie. */
  okno: string;
  /** Treść polecenia wraz z dołączonym wynikiem, gdy tryb tego wymaga. */
  tresc: string;
}

/** Stan potrzebny do rozdziału zlecenia między wykonawców. */
export interface RozdzialZlecenia {
  /** Okna wykonawców w kolejności Executor 1, Executor 2. */
  wykonawcy: readonly string[];
  /** Okno, które dostało zlecenie poprzednio; puste przed pierwszym. */
  poprzedni: string;
  /** Wynik ostatniej tury wykonawcy — treść dołączana w trybach zależnych. */
  wynik(okno: string): string;
}

/**
 * Rozdział zlecenia zgodny z trybem współpracy.
 *
 * Brak wykonawców daje pustą listę, nie błąd: okno wykonawcy może jeszcze nie
 * powstać, a zlecenie bez adresata nie jest usterką sceny.
 */
export function rozdziel(
  tryb: TrybWspolpracy,
  tresc: string,
  stan: RozdzialZlecenia,
): readonly Zlecenie[] {
  const [pierwszy, drugi] = stan.wykonawcy;
  if (pierwszy === undefined) return [];

  switch (tryb) {
    case TrybWspolpracy.Przekazywanie: {
      if (drugi === undefined) return [{ okno: pierwszy, tresc }];
      return [{ okno: drugi, tresc: zTrescia(tresc, stan.wynik(pierwszy), 'Executor 1') }];
    }
    case TrybWspolpracy.Naprzemienna: {
      const nastepny = stan.poprzedni === pierwszy ? (drugi ?? pierwszy) : pierwszy;
      return [{ okno: nastepny, tresc }];
    }
    case TrybWspolpracy.Iteracyjna: {
      const ten = stan.poprzedni === '' ? pierwszy : stan.poprzedni;
      return [{ okno: ten, tresc: zTrescia(tresc, stan.wynik(ten), 'poprzedniej iteracji') }];
    }
    default:
      return stan.wykonawcy.map((okno) => ({ okno, tresc }));
  }
}

/** Zlecenie wzbogacone o wynik poprzednika; pusty wynik nie dokłada nagłówka. */
function zTrescia(tresc: string, wynik: string, zrodlo: string): string {
  return wynik === '' ? tresc : `${tresc}\n\n--- Wynik ${zrodlo} ---\n${wynik}`;
}

/**
 * Zapis trybu współpracy na oknie koordynatora — jedno miejsce dla trzech okien.
 *
 * Selektor trybu stoi w oknie koordynatora i w obu oknach wykonawców, a zapis
 * jest jeden i ten sam, więc mieszka tutaj.
 *
 * Widok nie zostaje przy wartości, której rdzeń nie ma: `config.set` klucza
 * `multitasking.tryb_wspolpracy` kończy się odmową `validation_failed`, bo
 * klucz leży poza katalogiem ustawień. Odmowa cofa widok do wartości
 * poprzedniej, a zdanie mówi wprost, że zapisu nie ma.
 */
export async function zapiszTrybNaKoordynatorze(
  okna: ZrodloOkien,
  stan: StanMultitaskingu,
  nowy: TrybWspolpracy,
): Promise<{ zdanie: string; udane: boolean }> {
  const poprzedni = stan.tryb();
  const koordynator = stan.obsada().koordynator;
  if (koordynator === null) {
    stan.ustawTryb(nowy);
    return {
      zdanie: `Tryb „${nowy}" trzyma WYŁĄCZNIE ten widok — obsada nie ma koordynatora, na którego oknie rdzeń mógłby go zapisać.`,
      udane: false,
    };
  }
  stan.ustawTryb(nowy);
  const wynik = await okna.zapiszUstawienie({
    key: KLUCZ_TRYBU,
    value: nowy,
    scope: ConfigScope.Window,
    scopeId: koordynator.id,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    stan.ustawTryb(poprzedni);
    return {
      zdanie: `Rdzeń odmówił zapisu trybu, widok wraca do „${poprzedni}". Powód: ${wynik.blad?.message ?? 'rdzeń nie podał przyczyny'} (kod ${wynik.blad?.code ?? 'brak'}).`,
      udane: false,
    };
  }
  // Rdzeń oddaje zapisany wpis, więc to jego wartość mówi, co stoi
  // w konfiguracji okna koordynatora.
  const zapisany = trybZWartosci(wynik.wynik.value);
  stan.ustawTryb(zapisany);
  return {
    zdanie: `Rdzeń zapisał tryb współpracy „${zapisany}" na oknie koordynatora ${koordynator.id}.`,
    udane: true,
  };
}
