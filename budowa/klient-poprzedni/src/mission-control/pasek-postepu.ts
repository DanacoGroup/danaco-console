import { ProgressStatus } from '../../../shared/contract';

/**
 * Pasek postępu procesu: stopień ukończenia w postaci paska i napisu
 * „etap N/M". Napis jest obowiązkowy, bo stan nie bywa sygnalizowany samym
 * kolorem.
 *
 * Element niesie `role="progressbar"` wraz z `aria-valuenow`, więc czytnik ekranu
 * podaje tę samą wartość, którą widzi oko.
 */

/** Klasa modyfikująca barwę paska według stanu procesu. */
const KLASA_STANU: Readonly<Record<ProgressStatus, string>> = {
  [ProgressStatus.Pending]: 'mc-postep--oczekuje',
  [ProgressStatus.Running]: 'mc-postep--biegnie',
  [ProgressStatus.Paused]: 'mc-postep--wstrzymany',
  [ProgressStatus.Stopped]: 'mc-postep--zatrzymany',
  [ProgressStatus.Done]: 'mc-postep--zakonczony',
  [ProgressStatus.Failed]: 'mc-postep--bledny',
};

export interface OpisPostepu {
  etapBiezacy: number;
  etapowRazem: number;
  status: ProgressStatus;
  /** Nazwa procesu, do odczytu przez technologie wspomagające. */
  nazwa: string;
}

/** Oblicza stopień ukończenia w procentach; liczba etapów 0 znaczy „nieznana". */
export function stopienUkonczenia(etapBiezacy: number, etapowRazem: number): number {
  if (etapowRazem <= 0) {
    return 0;
  }
  const udzial = (etapBiezacy / etapowRazem) * 100;
  return Math.max(0, Math.min(100, Math.round(udzial)));
}

/** Buduje pasek postępu wraz z napisem „etap N/M". */
export function utworzPasekPostepu(opis: OpisPostepu): HTMLElement {
  const procent = stopienUkonczenia(opis.etapBiezacy, opis.etapowRazem);

  const owijka = document.createElement('div');
  owijka.className = 'mc-postep';

  const tor = document.createElement('div');
  tor.className = `mc-postep__tor ${KLASA_STANU[opis.status]}`;
  tor.setAttribute('role', 'progressbar');
  tor.setAttribute('aria-valuemin', '0');
  tor.setAttribute('aria-valuemax', String(opis.etapowRazem));
  tor.setAttribute('aria-valuenow', String(opis.etapBiezacy));
  tor.setAttribute(
    'aria-label',
    opis.etapowRazem > 0
      ? `${opis.nazwa} — etap ${opis.etapBiezacy} z ${opis.etapowRazem}`
      : `${opis.nazwa} — etap ${opis.etapBiezacy}, liczba etapów nieznana`,
  );

  const wypelnienie = document.createElement('span');
  wypelnienie.className = 'mc-postep__wypelnienie';
  wypelnienie.style.width = `${procent}%`;
  tor.append(wypelnienie);

  const napis = document.createElement('span');
  napis.className = 'mc-postep__napis';
  // Kontrakt: `totalSteps` równe 0 znaczy „liczba etapów nieznana".
  napis.textContent =
    opis.etapowRazem > 0 ? `etap ${opis.etapBiezacy}/${opis.etapowRazem}` : `etap ${opis.etapBiezacy}/?`;

  owijka.append(tor, napis);
  return owijka;
}
