import { ustawMotyw, type Motyw } from '../motyw/motyw';
import { czyIdGniazda, ograniczLiczbe } from './identyfikatory';
import type { StanPary } from './stan-pary';
import type { UkladOkien } from './uklad-okien';

/**
 * Parametry adresu strony podglądu.
 *
 * Wyłącznie na potrzeby podglądu wizualnego: pozwalają otworzyć scenę od razu
 * w wybranym stanie, bez klikania. Każdy wariant ekranu okien równoległych —
 * jedno, dwa i trzy okna, oba motywy, chwila przekazania — daje się przywołać
 * samym adresem.
 *
 *   ?liczba=3&motyw=light&stan=kolejka-wstrzymana&przekazanie=okno-1:okno-2
 *   ?modul=automations&liczba=4   → scena bierze dwa okna, bo tyle prowadzi moduł
 */
export interface ParametryPodgladu {
  /**
   * Kod modułu sceny; pusty zostawia moduł nieustalony.
   *
   * Bez rdzenia to jedyna droga, żeby zobaczyć figurę rozmowy modułu —
   * w produkcie moduł przychodzi zdarzeniem `window.changed`.
   */
  modul: string;
  liczba: number;
  motyw: Motyw | null;
  stan: StanPary | null;
  przekazanie: { od: string; do_: string } | null;
}

/** Odczyt parametrów z adresu; brak parametru oznacza wartość domyślną. */
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

/** Zastosowanie pozostałych parametrów na gotowym układzie. */
export function zastosujParametry(uklad: UkladOkien, parametry: ParametryPodgladu): void {
  // Moduł przed liczbą: to on rozstrzyga, ile okien scena w ogóle otworzy.
  uklad.ustawModulSceny(parametry.modul);
  uklad.ustawLiczbe(parametry.liczba);
  if (parametry.stan !== null) uklad.ustawStanPary(parametry.stan);
  if (parametry.przekazanie !== null) {
    uklad.pokazPrzekazanie(parametry.przekazanie.od, parametry.przekazanie.do_);
  }
}

/** Nazwa motywu albo brak wskazania. */
function czyMotyw(wartosc: string | null): Motyw | null {
  return wartosc === 'light' || wartosc === 'dark' ? wartosc : null;
}

/** Nazwa stanu pary albo brak wskazania. */
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

/** Para gniazd zapisana jako `od:do` albo brak wskazania. */
function czyPrzekazanie(wartosc: string | null): { od: string; do_: string } | null {
  if (wartosc === null) return null;
  const [od, do_] = wartosc.split(':');
  if (od === undefined || do_ === undefined) return null;
  return czyIdGniazda(od) && czyIdGniazda(do_) ? { od, do_ } : null;
}
