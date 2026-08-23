import { ProgressStatus, type ProgressChangedEvent } from '../../../../shared/contract';

/**
 * Pasek postępu generowania wraz z miniaturą w budowie.
 *
 * Nośnikiem stanu jest zdarzenie `progress.changed` z rdzenia. Pasek rusza
 * wyłącznie wtedy, gdy rdzeń przyśle postęp procesu tego okna.
 *
 * `design.asset.generate` postępu nie zgłasza — rdzeń nie ma czym wytworzyć
 * obrazu (`slowo-modelu.ts`) — więc okno kreatora wycisza pasek zamiast
 * zostawiać go na „czekam".
 *
 * Pasek bez oczekiwanego procesu milczy: `progress.changed` jedzie z każdego
 * strumienia rozmowy w rdzeniu (`core/telemetria_strumienia.go`), nie tylko
 * z pracy tego okna, więc przyjmowanie każdego zdarzenia pokazywałoby tu
 * postęp pracy, której to okno nie zlecało.
 *
 * Wygląd bierze z biblioteki (`dn-postep`) — plik nie zna ani jednej barwy.
 */
export interface PasekPostepu {
  element: HTMLElement;
  /** Przyjmuje zdarzenie postępu; obce okno i obcy proces są pomijane. */
  przyjmij(tresc: ProgressChangedEvent): void;
  /** Zapowiada proces, którego postępu okno oczekuje. */
  oczekuj(idProcesu: string, idOkna: string): void;
  /** Zdejmuje oczekiwanie i chowa pasek. */
  wycisz(): void;
  /** Stopień ukończenia ostatnio przyjęty; -1 znaczy „nic nie przyszło". */
  stopien(): number;
}

export function utworzPasekPostepu(): PasekPostepu {
  let oczekiwanyProces = '';
  let oczekiwaneOkno = '';
  let ostatni = -1;

  const wartosc = document.createElement('div');
  wartosc.className = 'dn-postep-wartosc';
  wartosc.style.width = '0%';

  const tor = document.createElement('div');
  tor.className = 'dn-postep-tor';
  tor.append(wartosc);

  const etykieta = document.createElement('span');
  etykieta.className = 'dn-postep-etykieta';
  etykieta.textContent = 'bez procesu';

  const pasek = document.createElement('div');
  pasek.className = 'dn-postep md-postep__pasek';
  pasek.setAttribute('role', 'progressbar');
  pasek.setAttribute('aria-valuemin', '0');
  pasek.setAttribute('aria-valuemax', '100');
  pasek.append(tor, etykieta);

  // Miniatura w budowie: prostokąt zastępczy rosnący wraz z postępem. Nie jest
  // podglądem wyniku — wyniku jeszcze nie ma — tylko widocznym śladem tego,
  // że rdzeń pracuje nad zasobem tego okna.
  const miniatura = document.createElement('div');
  miniatura.className = 'md-postep__miniatura';
  miniatura.dataset['stan'] = 'spoczynek';

  const element = document.createElement('div');
  element.className = 'md-postep';
  element.hidden = true;
  element.append(pasek, miniatura);

  /**
   * Czy to postęp tego okna. Bez oczekiwania odpowiedź brzmi „nie". Ze
   * wskazanym procesem liczy się zgodność procesu; dopóki rdzeń procesu nie
   * nazwał, wystarczy zgodność okna.
   */
  function czyNaszPostep(tresc: ProgressChangedEvent): boolean {
    if (oczekiwanyProces === '' && oczekiwaneOkno === '') return false;
    if (oczekiwanyProces !== '') return tresc.processId === oczekiwanyProces;
    return (tresc.windowId ?? '') === oczekiwaneOkno;
  }

  function pokaz(procent: number, opis: string, stan: string): void {
    ostatni = procent;
    element.hidden = false;
    wartosc.style.width = `${Math.max(0, Math.min(100, procent))}%`;
    pasek.setAttribute('aria-valuenow', String(procent));
    etykieta.textContent = opis;
    miniatura.dataset['stan'] = stan;
  }

  return {
    element,

    przyjmij(tresc) {
      if (!czyNaszPostep(tresc)) return;
      const etap = tresc.stepLabel ?? `etap ${tresc.currentStep} z ${tresc.totalSteps}`;
      pokaz(tresc.percent, `${tresc.percent}% · ${etap}`, tresc.status);
      if (tresc.status === ProgressStatus.Done) miniatura.dataset['stan'] = ProgressStatus.Done;
    },

    oczekuj(idProcesu, idOkna) {
      oczekiwanyProces = idProcesu;
      oczekiwaneOkno = idOkna;
      pokaz(0, 'zlecenie przyjęte — rdzeń nie przysłał jeszcze postępu', 'oczekiwanie');
    },

    wycisz() {
      oczekiwanyProces = '';
      oczekiwaneOkno = '';
      ostatni = -1;
      element.hidden = true;
      miniatura.dataset['stan'] = 'spoczynek';
    },

    stopien: () => ostatni,
  };
}
