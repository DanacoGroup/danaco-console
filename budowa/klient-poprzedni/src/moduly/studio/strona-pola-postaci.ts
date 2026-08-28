import type { StudioActionBalance, StudioSkippedItem } from '../../../../shared/contract';
import type { PoleFormularza } from '../../modele/kontrolki-formularza';

/**
 * Granice i krok pola liczbowego postaci dokumentu; wartość początkowa nieobecna
 * zostawia pole puste, co w rodzinie komend postaci znaczy brak zmiany cechy.
 */
export interface GranicePola {
  dolna?: number;
  gorna?: number;
  krok?: number;
  /** Wartość pokazywana na starcie; brak zostawia pole puste. */
  wartosc?: number;
}

/** Pole liczbowe formularza postaci wraz z etykietą, granicami wartości, krokiem i zdaniem wyjaśniającym pod polem. */
export function poleLiczbowe(
  opis: { etykieta: string; opis?: string; podpowiedz?: string },
  granice: GranicePola = {},
): PoleFormularza<HTMLInputElement> {
  const kontrolka = document.createElement('input');
  kontrolka.type = 'number';
  kontrolka.className = 'dn-pole-kontrolka';
  kontrolka.id = `dm-postac-${licznik()}`;
  if (granice.dolna !== undefined) kontrolka.min = String(granice.dolna);
  if (granice.gorna !== undefined) kontrolka.max = String(granice.gorna);
  kontrolka.step = String(granice.krok ?? 1);
  if (granice.wartosc !== undefined) kontrolka.value = String(granice.wartosc);
  if (opis.podpowiedz !== undefined) kontrolka.placeholder = opis.podpowiedz;

  const etykieta = document.createElement('label');
  etykieta.className = 'dn-pole-etykieta';
  etykieta.htmlFor = kontrolka.id;
  etykieta.textContent = opis.etykieta;

  const wiersz = document.createElement('div');
  wiersz.className = 'dn-pole dm-pole';
  wiersz.append(etykieta, kontrolka);
  if (opis.opis !== undefined) {
    const zdanie = document.createElement('p');
    zdanie.className = 'dn-pole-opis';
    zdanie.textContent = opis.opis;
    wiersz.append(zdanie);
  }
  return { element: wiersz, kontrolka };
}

let numer = 0;
function licznik(): number {
  numer += 1;
  return numer;
}

/**
 * Liczba z pola albo `undefined`, gdy pole jest puste albo jego treść nie jest
 * liczbą; wartość nieokreślona jedzie do rdzenia jako brak zmiany tej cechy.
 */
export function liczbaPola(kontrolka: HTMLInputElement): number | undefined {
  const surowa = kontrolka.value.trim();
  if (surowa === '') return undefined;
  const wartosc = Number(surowa);
  return Number.isFinite(wartosc) ? wartosc : undefined;
}

/** Napis z pola tekstowego albo obszaru wielowierszowego formularza postaci; `undefined`, gdy pole jest puste. */
export function tekstPola(kontrolka: HTMLInputElement | HTMLTextAreaElement): string | undefined {
  const surowy = kontrolka.value;
  return surowy === '' ? undefined : surowy;
}

/**
 * Wartość listy wyboru postaci albo `undefined` dla pozycji oznaczającej brak
 * zmiany cechy; pozycja o wartości pustej stoi w każdej liście jako pierwsza.
 */
export function wyborPola(kontrolka: HTMLSelectElement): string | undefined {
  return kontrolka.value === '' ? undefined : kontrolka.value;
}

/** Pozycja „bez zmiany" — pierwsza pozycja w każdej liście wyboru nastaw postaci tego edytowanego dokumentu. */
export const BEZ_ZMIANY = { wartosc: '', etykieta: '— bez zmiany —' };

/**
 * Zdanie o bilansie czynności masowej: liczba miejsc zmienionych, pominiętych,
 * odłożonych i spiętych, wraz z wykazem powodów każdego pominięcia z osobna.
 */
export function opiszBilans(bilans: StudioActionBalance): string {
  const czesci: string[] = [`miejsc zmienionych ${bilans.applied}`];
  if (bilans.skippedCount > 0) czesci.push(`pominiętych ${bilans.skippedCount}`);
  if (bilans.deferredCount !== undefined && bilans.deferredCount > 0) {
    czesci.push(`odłożonych, żeby nie nadpisać cudzej pracy: ${bilans.deferredCount}`);
  }
  if (bilans.conflicts !== undefined && bilans.conflicts.length > 0) {
    czesci.push(`spięć o ten sam fragment ${bilans.conflicts.length}`);
  }
  const powody = (bilans.skipped ?? []).map(opiszPominiecie).filter((zdanie) => zdanie !== '');
  const nota = bilans.note === undefined || bilans.note === '' ? '' : ` ${bilans.note}`;
  const wykaz = powody.length === 0 ? '' : ` Pominięcia: ${powody.join('; ')}.`;
  return `${czesci.join(', ')}.${nota}${wykaz}`;
}

/** Jedno pominięcie czynności masowej nazwane wraz z zakresem znaków dokumentu i blokadą, która je wywołała. */
function opiszPominiecie(pozycja: StudioSkippedItem): string {
  const zakres =
    pozycja.rangeStart === undefined || pozycja.rangeEnd === undefined
      ? ''
      : ` (znaki ${pozycja.rangeStart}–${pozycja.rangeEnd})`;
  const blokada =
    pozycja.lockName !== undefined && pozycja.lockName !== ''
      ? ` — blokada „${pozycja.lockName}"`
      : pozycja.lockId !== undefined && pozycja.lockId !== ''
        ? ` — blokada ${pozycja.lockId}`
        : '';
  const szczegol =
    pozycja.detail === undefined || pozycja.detail === '' ? '' : `: ${pozycja.detail}`;
  return `${pozycja.reason}${szczegol}${zakres}${blokada}`;
}
