/**
 * Wykazy Automation Studio: karty kroków scenariusza, wiersze automatyk oraz
 * trzy stany obowiązkowe wykazu. Plik odpowiada za postać wykazu na ekranie;
 * okno składa formularz i prowadzi rozmowę z rdzeniem, a tutaj jej wynik
 * staje się wierszem.
 */
import type {
  AutomationExecution,
  AutomationStep,
  AutomationWorkflow,
} from '../../../../shared/contract';
import { KLASA_PRZYCISKU, przyciskCzynnosci } from './przyciski-browser';
import type { StanOkna } from './stan-okna';

/**
 * Karty kroków scenariusza jako odczyt widoku tekstowego, a nie drugi zbiór
 * kroków. Każda karta niesie rodzaj kroku, jego komendę oraz nazwę, gdy krok
 * nazwę posiada.
 */
export function wierszeKrokow(kroki: readonly AutomationStep[]): HTMLElement[] {
  return kroki.map((krok) => {
    const element = document.createElement('li');
    element.className = 'mb-scenariusz__krok';
    const nazwa = (krok.name ?? '').trim();
    element.textContent =
      `${krok.kind}${krok.command === undefined ? '' : ` · ${krok.command}`}` +
      `${nazwa === '' ? '' : ` — ${nazwa}`}`;
    return element;
  });
}

/**
 * Jedna automatyka wykazu wraz ze zdaniem o jej przebiegach: nazwa, liczba
 * kroków, znacznik czynności oraz przycisk wczytujący automatykę do
 * formularza okna.
 */
export function wierszAutomatyki(
  automatyka: AutomationWorkflow,
  oPrzebiegach: string,
  wybrana: boolean,
  otworz: () => void,
): HTMLElement {
  const opis = document.createElement('span');
  opis.className = 'mb-automatyki__opis';
  opis.textContent =
    `${automatyka.name} · kroków ${(automatyka.steps ?? []).length} · ` +
    `${automatyka.enabled ? 'czynna' : 'wstrzymana'}` +
    `${oPrzebiegach === '' ? '' : ` · ${oPrzebiegach}`}`;

  const element = document.createElement('li');
  element.className = 'mb-automatyki__wiersz';
  element.dataset['wybrana'] = wybrana ? 'tak' : 'nie';
  element.append(opis, przyciskCzynnosci('Wczytaj do formularza', KLASA_PRZYCISKU.duch, otworz));
  return element;
}

/**
 * Zdanie o przebiegach automatyki: liczba przebiegów oraz stan najświeższego
 * z nich, wybranego po znaczniku rozpoczęcia. Pusty zbiór daje zdanie o braku
 * przebiegów.
 */
export function zdanieOPrzebiegach(przebiegi: readonly AutomationExecution[]): string {
  const najnowszy = [...przebiegi].sort((jeden, drugi) => drugi.startedAt - jeden.startedAt)[0];
  if (najnowszy === undefined) return 'bez przebiegów';
  return `przebiegów ${przebiegi.length}, ostatni: ${najnowszy.status}`;
}

/**
 * Stan wykazu widziany przez okno w chwili przerysowania: znacznik trwającego
 * odczytu, powód odmowy ostatniego odczytu oraz liczba wierszy stojących na
 * ekranie.
 */
export interface StanWykazuAutomatyk {
  wOdczycie: boolean;
  /** Powód odmowy ostatniego odczytu; pusty napis znaczy „bez uwag". */
  odmowa: string;
  pozycji: number;
}

/**
 * Trzy stany obowiązkowe wykazu automatyk: czekanie, odmowa, pustka.
 *
 * Odmowa jest błędem okna wyłącznie przy pustym wykazie — z automatykami na
 * ekranie wpisy zostają, a powód idzie zdaniem w wierszu odpowiedzi.
 */
export function nanieStanWykazu(okno: StanOkna, stan: StanWykazuAutomatyk): void {
  if (stan.wOdczycie) {
    okno.ladowanie('Odczyt automatyk z rdzenia w toku…');
    return;
  }
  if (stan.odmowa !== '' && stan.pozycji === 0) {
    okno.blad(stan.odmowa);
    return;
  }
  if (stan.pozycji === 0) {
    okno.pusteZOpisu();
    return;
  }
  okno.gotowe();
}
