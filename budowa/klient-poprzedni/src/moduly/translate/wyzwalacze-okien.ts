import { przycisk } from '../../modele/kontrolki-formularza';
import { WARSTWY, oznaczWarstwe, type WarstwaWidocznosci } from './warstwy-translate';

/**
 * Wyzwalacze okien warstw drugiej i trzeciej to pas, z którego otwiera się to, co w stanie
 * spoczynku modułu jest zwinięte, ale nie zablokowane.
 */
export interface OpisWyzwalacza {
  /** Kod okna — ten sam, którym moduł zgłasza je w katalogu okien. */
  readonly kod: string;
  /** Nazwa własna okna. */
  readonly nazwa: string;
  /** Warstwa, na której okno stoi. */
  readonly warstwa: WarstwaWidocznosci;
  /** Znacznik wywołania z wykazu ikonografii. */
  readonly znacznik: string;
  /** Element okna, którego widocznością wyzwalacz steruje. */
  readonly okno: HTMLElement;
  /** Skrót klawiszowy okna, gdy opracowanie mu go nadaje. */
  readonly skrot?: string;
}

export interface WyzwalaczeOkien {
  /** Pas wyzwalaczy osadzany w module. */
  element: HTMLElement;
  /** Otwiera okno i przewija do niego — droga skrótu klawiszowego. */
  otworz(kod: string): void;
}

export function utworzWyzwalaczeOkien(opisy: readonly OpisWyzwalacza[]): WyzwalaczeOkien {
  const uchwyty = new Map<string, HTMLButtonElement>();

  const element = document.createElement('div');
  element.className = 'mt-wyzwalacze';
  element.setAttribute('aria-label', 'Wyzwalacze okien warstw drugiej i trzeciej');

  for (const opis of opisy) {
    const uchwyt = przycisk(etykieta(opis), 'dn-btn dn-btn--sm dn-btn--zarys mt-wyzwalacze__uchwyt');
    oznaczWarstwe(uchwyt, opis.warstwa);
    // Uchwyt znakuje się inaczej niż okno, którym steruje, inaczej dwa węzły odpowiadałyby zapytaniu.
    uchwyt.dataset['wyzwalacz'] = opis.kod;
    uchwyt.title = `${WARSTWY[opis.warstwa].nazwa} — ${WARSTWY[opis.warstwa].dostep}`;
    uchwyt.addEventListener('click', () => ustaw(opis, opis.okno.hidden));
    uchwyty.set(opis.kod, uchwyt);
    element.append(uchwyt);
    // Postać spoczynkowa: okno warstwy drugiej i trzeciej startuje zwinięte.
    ustaw(opis, false);
  }

  function ustaw(opis: OpisWyzwalacza, otwarte: boolean): void {
    opis.okno.hidden = !otwarte;
    const uchwyt = uchwyty.get(opis.kod);
    if (uchwyt === undefined) return;
    uchwyt.setAttribute('aria-expanded', String(otwarte));
    uchwyt.textContent = etykieta(opis, otwarte);
  }

  return {
    element,

    otworz(kod) {
      const opis = opisy.find((wpis) => wpis.kod === kod);
      if (opis === undefined) return;
      ustaw(opis, true);
      opis.okno.scrollIntoView({ block: 'nearest' });
    },
  };
}

/**
 * Etykieta uchwytu niesie czynność, którą naciśnięcie wykona: uchwyt zwinięty mówi Otwórz,
 * rozwinięty Zamknij, zamiast liczyć na sam stan otwarcia.
 */
function etykieta(opis: OpisWyzwalacza, otwarte = false): string {
  const czynnosc = otwarte ? 'Zamknij' : 'Otwórz';
  const skrot = opis.skrot === undefined ? '' : ` (${opis.skrot})`;
  return `${opis.znacznik} ${czynnosc}: ${opis.nazwa}${skrot}`;
}
