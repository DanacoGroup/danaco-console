import type { StudioActionBalance, StudioSkippedItem } from '../../../../shared/contract';
import type { PoleFormularza } from '../../modele/kontrolki-formularza';

/**
 * Drobne kontrolki i zdania wspólne panelom postaci dokumentu.
 *
 * ── Dlaczego pole liczbowe jest własne ──────────────────────────────────────
 * `kontrolki-formularza.ts` niesie pole tekstowe, listę wyboru i przełącznik.
 * Nastawy postaci są jednak w większości liczbami z granicami: margines
 * w milimetrach, stopień pisma w punktach, punkt startu numeracji, krycie znaku
 * wodnego. Pole tekstowe przyjęłoby „dwadzieścia" i wysłało to do rdzenia.
 *
 * ── Dlaczego pole puste znaczy „nie ruszaj" ─────────────────────────────────
 * To jest zasada CAŁEGO tego odcinka. Rodzina `studio.page.*` i `studio.format.*`
 * ma pola opcjonalne w znaczeniu „tej cechy nie zmieniam". Margines zerowy jest
 * nastawą, którą Operator może wybrać świadomie, więc „podano zero" i „nie
 * podano" nie mogą znaczyć tego samego — inaczej każde naciśnięcie przycisku
 * zerowałoby wszystko, czego Operator nie wpisał.
 */

/** Granice i krok pola liczbowego. */
export interface GranicePola {
  dolna?: number;
  gorna?: number;
  krok?: number;
  /** Wartość pokazywana na starcie; brak zostawia pole puste. */
  wartosc?: number;
}

/** Pole liczbowe wraz z etykietą i zdaniem wyjaśniającym. */
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
 * Liczba z pola albo `undefined`, gdy pole jest puste.
 *
 * `undefined` jedzie do rdzenia jako brak pola, czyli „tej cechy nie zmieniam".
 * Wartość nieliczbowa daje to samo — przeglądarka nie wpuści jej do pola typu
 * liczbowego, a gdyby wpuściła, cisza jest bezpieczniejsza niż `NaN` w żądaniu.
 */
export function liczbaPola(kontrolka: HTMLInputElement): number | undefined {
  const surowa = kontrolka.value.trim();
  if (surowa === '') return undefined;
  const wartosc = Number(surowa);
  return Number.isFinite(wartosc) ? wartosc : undefined;
}

/** Napis z pola albo `undefined`, gdy pole jest puste. */
export function tekstPola(kontrolka: HTMLInputElement | HTMLTextAreaElement): string | undefined {
  const surowy = kontrolka.value;
  return surowy === '' ? undefined : surowy;
}

/**
 * Wartość listy wyboru albo `undefined` dla pozycji „bez zmiany".
 *
 * Pozycja o wartości pustej stoi w każdej liście tego odcinka jako pierwsza i
 * znaczy dosłownie „nie ruszaj tej cechy". Bez niej lista wyboru zawsze coś
 * narzucałaby, bo `select` nie ma stanu „nic nie wybrano".
 */
export function wyborPola(kontrolka: HTMLSelectElement): string | undefined {
  return kontrolka.value === '' ? undefined : kontrolka.value;
}

/** Pozycja „bez zmiany" — pierwsza w każdej liście nastaw postaci. */
export const BEZ_ZMIANY = { wartosc: '', etykieta: '— bez zmiany —' };

/**
 * Zdanie o bilansie czynności — co zmienione, co pominięte i przez co.
 *
 * Bilans jest sedno uczciwości tego odcinka: zamiana w całym dokumencie, która
 * trafiła w blokadę, wykonuje się POZA blokadą i musi powiedzieć, którą.
 * Przemilczenie pominięcia jest tu zakazane, a `applied: 0` bez słowa byłoby
 * najgorszą możliwą odpowiedzią — Operator myślałby, że czynność się wykonała.
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

/** Jedno pominięcie nazwane wraz z blokadą, która je wywołała. */
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
