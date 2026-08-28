import type { AutomationExecution } from '../../../../shared/contract';
import { nastepneUruchomienia } from './nastepne-uruchomienia';

/** Kalendarz uruchomień okna Scheduler pokazuje podgląd terminów zaplanowanych i przebiegów odbytych. */

/** Zakres podglądu kalendarza, tydzień albo miesiąc, wybierany przełącznikiem nad siatką dni kalendarza. */
export const ZAKRESY_KALENDARZA = {
  tydzien: 'tydzien',
  miesiac: 'miesiac',
} as const;

export type ZakresKalendarza = (typeof ZAKRESY_KALENDARZA)[keyof typeof ZAKRESY_KALENDARZA];

/** Zakresy w kolejności przełącznika wraz z ich nazwami widocznymi Operatorowi na ekranie tego kalendarza. */
export const NAZWY_ZAKRESOW: ReadonlyArray<[ZakresKalendarza, string]> = [
  [ZAKRESY_KALENDARZA.tydzien, 'tydzień'],
  [ZAKRESY_KALENDARZA.miesiac, 'miesiąc'],
];

/** Ile dni obejmuje każdy zakres podglądu; miesiąc liczy się czterema pełnymi tygodniami, nie datami kalendarza. */
const DNI_ZAKRESU: Readonly<Record<ZakresKalendarza, number>> = {
  [ZAKRESY_KALENDARZA.tydzien]: 7,
  [ZAKRESY_KALENDARZA.miesiac]: 28,
};

/** Ile terminów zaplanowanych wolno wyliczyć na potrzeby kalendarza, zanim wyliczanie zostanie odcięte. */
const GRANICA_TERMINOW = 60;

/** Jeden dzień kalendarza wraz z tym, co na nim stoi: zapis daty, liczba terminów i liczba przebiegów dnia. */
export interface DzienKalendarza {
  /** Dzień w zapisie ISO (RRRR-MM-DD), w czasie miejscowym. */
  dzien: string;
  /** Liczba terminów zaplanowanych, które wypadają tego dnia. */
  zaplanowane: number;
  /** Liczba przebiegów rozpoczętych tego dnia. */
  odbyte: number;
}

/** Składa dni kalendarza od dnia bieżącego wprzód; funkcja jest czysta, sprawdzian czyta ją wprost bez okna. */
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

/** Buduje siatkę kalendarza; dzień pusty zostaje w siatce, bo luka między uruchomieniami jest informacją. */
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

/** Zdanie pod kalendarzem nazywa, skąd biorą się terminy i co jest terminem wiążącym dla samego Operatora. */
export function podpisKalendarza(dni: readonly DzienKalendarza[]): string {
  const zaplanowane = dni.reduce((suma, wpis) => suma + wpis.zaplanowane, 0);
  const odbyte = dni.reduce((suma, wpis) => suma + wpis.odbyte, 0);
  return (
    `Kalendarz obejmuje ${dni.length} dni: ${zaplanowane} terminów zaplanowanych, ` +
    `${odbyte} przebiegów odbytych. Terminy wylicza okno z wpisanej cykliczności; ` +
    'terminem obowiązującym pozostaje najbliższe uruchomienie oddane przez rdzeń po zapisie.'
  );
}

/** Zdanie w kaflu dnia; dzień bez niczego mówi to wprost, a nie zostaje pusty bez żadnego zdania widocznego. */
function opisDnia(wpis: DzienKalendarza): string {
  if (wpis.zaplanowane === 0 && wpis.odbyte === 0) return 'bez uruchomień';
  const czesci: string[] = [];
  if (wpis.zaplanowane > 0) czesci.push(`zaplanowane ${wpis.zaplanowane}`);
  if (wpis.odbyte > 0) czesci.push(`odbyte ${wpis.odbyte}`);
  return czesci.join(' · ');
}

/** Dzień w zapisie ISO, w czasie miejscowym, jest kluczem siatki i podpisem kafla widocznym dla Operatora. */
function zapisDnia(data: Date): string {
  const miesiac = String(data.getMonth() + 1).padStart(2, '0');
  const dzien = String(data.getDate()).padStart(2, '0');
  return `${data.getFullYear()}-${miesiac}-${dzien}`;
}
