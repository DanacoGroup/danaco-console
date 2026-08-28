import { ustawMotyw, type Motyw } from '../motyw/motyw';
import { czyIdGniazda, ograniczLiczbe } from './identyfikatory';
import type { StanPary } from './stan-pary';
import type { UkladOkien } from './uklad-okien';

/**
 * Parametry adresu strony podglądu pozwalają otworzyć scenę od razu w wybranym stanie bez klikania, przywołując samym adresem dowolny wariant ekranu okien równoległych — liczbę okien, motyw, stan pary czy chwilę przekazania.
 */
export interface ParametryPodgladu {
  /** Kod modułu sceny; pusty zostawia moduł nieustalony, bo to jedyna droga do figury bez rdzenia. */
  modul: string;
  liczba: number;
  motyw: Motyw | null;
  stan: StanPary | null;
  przekazanie: { od: string; do_: string } | null;
}

/** Odczyt parametrów z adresu strony podglądu; brak parametru w adresie oznacza przyjęcie wartości domyślnej. */
export function odczytajParametry(adres: string): ParametryPodgladu {
  const pytanie = new URL(adres).searchParams;
  return {
    modul: pytanie.get('modul') ?? '',
    liczba: ograniczLiczbe(Number(pytanie.get('liczba') ?? 2)),
    motyw: czyMotyw(pytanie.get('motyw')),
    stan: czyStan(pytanie.get('stan')),
    przekazanie: czyPrzekazanie(pytanie.get('przekazanie')),
  };
}

/**
 * Wybór motywu z adresu.
 *
 * Stosowany przed zbudowaniem sceny: przełączenie motywu ma przejście barwne,
 * więc ustawienie go po pierwszym rysowaniu dawałoby na obrazie ekranu barwy
 * uchwycone w połowie przejścia.
 */
export function zastosujMotywPodgladu(parametry: ParametryPodgladu): void {
  if (parametry.motyw !== null) ustawMotyw(parametry.motyw);
}

/** Zastosowanie pozostałych parametrów podglądu na już gotowym, zmontowanym układzie okien równoległych. */
export function zastosujParametry(uklad: UkladOkien, parametry: ParametryPodgladu): void {
  // Moduł przed liczbą: to on rozstrzyga, ile okien scena w ogóle otworzy.
  uklad.ustawModulSceny(parametry.modul);
  uklad.ustawLiczbe(parametry.liczba);
  if (parametry.stan !== null) uklad.ustawStanPary(parametry.stan);
  if (parametry.przekazanie !== null) {
    uklad.pokazPrzekazanie(parametry.przekazanie.od, parametry.przekazanie.do_);
  }
}

/** Nazwa motywu odczytana z parametru adresu albo brak wskazania, gdy wartość jest zupełnie nierozpoznawalna. */
function czyMotyw(wartosc: string | null): Motyw | null {
  return wartosc === 'light' || wartosc === 'dark' ? wartosc : null;
}

/** Nazwa stanu pary odczytana z parametru adresu albo brak wskazania, gdy wartość jest nierozpoznawalna. */
function czyStan(wartosc: string | null): StanPary | null {
  const znane: StanPary[] = [
    'brak-pary',
    'gotowa',
    'wykonawca-pracuje',
    'koordynator-wybudzony',
    'kolejka-wstrzymana',
    'przekazanie',
  ];
  return znane.find((stan) => stan === wartosc) ?? null;
}

/** Para gniazd zapisana w adresie jako źródło i cel rozdzielone dwukropkiem, albo brak jakiegokolwiek wskazania. */
function czyPrzekazanie(wartosc: string | null): { od: string; do_: string } | null {
  if (wartosc === null) return null;
  const [od, do_] = wartosc.split(':');
  if (od === undefined || do_ === undefined) return null;
  return czyIdGniazda(od) && czyIdGniazda(do_) ? { od, do_ } : null;
}
