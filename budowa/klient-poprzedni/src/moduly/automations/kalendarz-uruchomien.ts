import type { AutomationExecution } from '../../../../shared/contract';
import { nastepneUruchomienia } from './nastepne-uruchomienia';

/**
 * Kalendarz uruchomień okna Scheduler — miesięczny i tygodniowy podgląd
 * terminów zaplanowanych oraz przebiegów już odbytych.
 *
 * Wykaz pięciu najbliższych terminów mówi „kiedy najbliżej"; kalendarz mówi
 * „jak gęsto" i pokazuje to razem z historią, dzięki czemu widać dzień, w
 * którym uruchomienie wypadało, a przebiegu nie było.
 *
 * Terminy zaplanowane liczy okno z wpisanej cykliczności (`nastepne-uruchomienia.ts`),
 * bo dotyczą reguły jeszcze niezapisanej. Terminem obowiązującym pozostaje
 * `AutomationSchedule.nextRunAt` liczone przez rdzeń — kalendarz jest podglądem
 * i mówi to wprost w swoim podpisie.
 *
 * Dni układają się według czasu miejscowego przeglądarki, bo tak Operator czyta
 * kalendarz; rachunek terminu idzie w UTC, tak samo jak w rdzeniu.
 */

/** Zakres podglądu kalendarza. */
export const ZAKRESY_KALENDARZA = {
  tydzien: 'tydzien',
  miesiac: 'miesiac',
} as const;

export type ZakresKalendarza = (typeof ZAKRESY_KALENDARZA)[keyof typeof ZAKRESY_KALENDARZA];

/** Zakresy w kolejności przełącznika wraz z ich nazwami na ekranie. */
export const NAZWY_ZAKRESOW: ReadonlyArray<[ZakresKalendarza, string]> = [
  [ZAKRESY_KALENDARZA.tydzien, 'tydzień'],
  [ZAKRESY_KALENDARZA.miesiac, 'miesiąc'],
];

/** Ile dni obejmuje każdy zakres; miesiąc liczymy czterema pełnymi tygodniami. */
const DNI_ZAKRESU: Readonly<Record<ZakresKalendarza, number>> = {
  [ZAKRESY_KALENDARZA.tydzien]: 7,
  [ZAKRESY_KALENDARZA.miesiac]: 28,
};

/** Ile terminów zaplanowanych wolno wyliczyć na potrzeby kalendarza. */
const GRANICA_TERMINOW = 60;

/** Jeden dzień kalendarza wraz z tym, co na nim stoi. */
export interface DzienKalendarza {
  /** Dzień w zapisie ISO (RRRR-MM-DD), w czasie miejscowym. */
  dzien: string;
  /** Liczba terminów zaplanowanych, które wypadają tego dnia. */
  zaplanowane: number;
  /** Liczba przebiegów rozpoczętych tego dnia. */
  odbyte: number;
}

/**
 * Składa dni kalendarza od dnia bieżącego wprzód.
 *
 * Funkcja jest czysta — sprawdzian czyta ją wprost, bez montażu okna.
 */
export function dniKalendarza(
  zapisCyklicznosci: string,
  przebiegi: readonly AutomationExecution[],
  zakres: ZakresKalendarza,
  od: Date = new Date(),
): DzienKalendarza[] {
  const ile = DNI_ZAKRESU[zakres];
  const dni: DzienKalendarza[] = [];
  const miejsca = new Map<string, DzienKalendarza>();
  for (let numer = 0; numer < ile; numer += 1) {
    const data = new Date(od.getFullYear(), od.getMonth(), od.getDate() + numer);
    const wpis: DzienKalendarza = { dzien: zapisDnia(data), zaplanowane: 0, odbyte: 0 };
    dni.push(wpis);
    miejsca.set(wpis.dzien, wpis);
  }

  for (const termin of nastepneUruchomienia(zapisCyklicznosci, od, GRANICA_TERMINOW)) {
    const wpis = miejsca.get(zapisDnia(termin));
    if (wpis !== undefined) wpis.zaplanowane += 1;
  }
  for (const przebieg of przebiegi) {
    const wpis = miejsca.get(zapisDnia(new Date(przebieg.startedAt)));
    if (wpis !== undefined) wpis.odbyte += 1;
  }
  return dni;
}

/**
 * Buduje siatkę kalendarza. Dzień pusty zostaje w siatce — luka między
 * uruchomieniami jest informacją, a siatka z wyciętymi dniami przestałaby być
 * kalendarzem.
 */
export function siatkaKalendarza(dni: readonly DzienKalendarza[]): HTMLElement {
  const siatka = document.createElement('div');
  siatka.className = 'da-kalendarz';
  siatka.setAttribute('role', 'list');
  siatka.setAttribute('aria-label', 'Kalendarz uruchomień automatyki');
  for (const wpis of dni) {
    const kafel = document.createElement('div');
    kafel.className = 'da-kalendarz__dzien';
    kafel.setAttribute('role', 'listitem');
    kafel.dataset['dzien'] = wpis.dzien;
    if (wpis.zaplanowane > 0) kafel.dataset['zaplanowane'] = String(wpis.zaplanowane);
    if (wpis.odbyte > 0) kafel.dataset['odbyte'] = String(wpis.odbyte);

    const data = document.createElement('span');
    data.className = 'da-kalendarz__data';
    data.textContent = wpis.dzien.slice(8);

    const opis = document.createElement('span');
    opis.className = 'da-kalendarz__opis';
    opis.textContent = opisDnia(wpis);

    kafel.title = `${wpis.dzien}: ${opisDnia(wpis)}`;
    kafel.append(data, opis);
    siatka.append(kafel);
  }
  return siatka;
}

/** Zdanie pod kalendarzem: skąd biorą się terminy i co jest terminem wiążącym. */
export function podpisKalendarza(dni: readonly DzienKalendarza[]): string {
  const zaplanowane = dni.reduce((suma, wpis) => suma + wpis.zaplanowane, 0);
  const odbyte = dni.reduce((suma, wpis) => suma + wpis.odbyte, 0);
  return (
    `Kalendarz obejmuje ${dni.length} dni: ${zaplanowane} terminów zaplanowanych, ` +
    `${odbyte} przebiegów odbytych. Terminy wylicza okno z wpisanej cykliczności; ` +
    'terminem obowiązującym pozostaje najbliższe uruchomienie oddane przez rdzeń po zapisie.'
  );
}

/** Zdanie w kaflu dnia; dzień bez niczego mówi to wprost, a nie zostaje pusty. */
function opisDnia(wpis: DzienKalendarza): string {
  if (wpis.zaplanowane === 0 && wpis.odbyte === 0) return 'bez uruchomień';
  const czesci: string[] = [];
  if (wpis.zaplanowane > 0) czesci.push(`zaplanowane ${wpis.zaplanowane}`);
  if (wpis.odbyte > 0) czesci.push(`odbyte ${wpis.odbyte}`);
  return czesci.join(' · ');
}

/** Dzień w zapisie ISO, w czasie miejscowym — kluczem siatki i podpisem kafla. */
function zapisDnia(data: Date): string {
  const miesiac = String(data.getMonth() + 1).padStart(2, '0');
  const dzien = String(data.getDate()).padStart(2, '0');
  return `${data.getFullYear()}-${miesiac}-${dzien}`;
}
