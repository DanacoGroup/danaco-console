/**
 * Panel Sesje i Projekty. Dwie zakładki nad jednym pasmem treści: karty
 * sesji odtworzone przez rdzeń przy wejściu do środowiska, i projekty —
 * rdzeń nie wystawia dziś osobnego wykazu projektów, więc ta zakładka stoi
 * nazwanym stanem pustym zamiast wartości zmyślonej po stronie klienta.
 */

import type { Session } from '../../../../shared/contract.ts';
import { SessionStatus } from '../../../../shared/contract.ts';
import { ikony } from '../ikony.ts';
import { el, tekst, zeZnacznika } from '../narzedzia.ts';

export interface WlasciwosciPanelu {
  /** Karty sesji otwarte w środowisku, w kolejności odebranej od rdzenia. */
  sesje: Session[];
}

/** Wartość `data-stan` znana arkuszowi `panel-sesji.css` — status spoza wykazu zostaje bez znacznika. */
const STAN_WIERSZA: Partial<Record<Session['status'], string>> = {
  [SessionStatus.Active]: 'praca',
  [SessionStatus.Paused]: 'praca',
  [SessionStatus.Finished]: 'zakonczone',
  [SessionStatus.Archived]: 'zakonczone',
};

function znak(rysunek: string): SVGElement {
  const wezel = zeZnacznika(rysunek);
  wezel.setAttribute('aria-hidden', 'true');
  return wezel;
}

function przyciskIkony(rysunek: keyof typeof ikony, etykieta: string, atrybuty?: Record<string, string>): HTMLElement {
  return el(
    'button',
    { klasa: 'dn-btn-ikona dn-etykietka', type: 'button', 'data-etykietka': etykieta, 'aria-label': etykieta, ...atrybuty },
    [znak(ikony[rysunek])],
  );
}

function wierszSesji(sesja: Session, biezaca: boolean): HTMLElement {
  const stan = STAN_WIERSZA[sesja.status];
  const nazwa = sesja.title ?? sesja.id;
  return el('div', { klasa: 'dn-nrz-menu sta-menu dn-panel-wiersz' }, [
    el(
      'button',
      {
        klasa: 'dn-obszar-pozycja dn-panel-poz',
        type: 'button',
        'data-sesja-id': sesja.id,
        'data-stan': stan ?? null,
        'aria-current': biezaca ? 'true' : null,
      },
      [
        el('span', { klasa: 'dn-stan-znacznik', 'data-stan': stan ?? null }),
        el('span', { klasa: 'dn-obszar-pozycja-tekst' }, [
          el('span', { klasa: 'dn-obszar-pozycja-nazwa', tekst: nazwa }),
        ]),
      ],
    ),
  ]);
}

function zakladka(id: string, rysunek: string, etykieta: string, wybrana: boolean): HTMLElement {
  return el(
    'div',
    {
      klasa: 'dn-karta-widoku',
      role: 'tab',
      tabindex: '0',
      'aria-selected': wybrana ? 'true' : 'false',
      'aria-controls': id,
    },
    [znak(rysunek), el('span', { klasa: 'dn-karta-widoku-nazwa', tekst: etykieta })],
  );
}

export function panelSesje(w: WlasciwosciPanelu): HTMLElement {
  const pasmoZakladek = el(
    'div',
    { klasa: 'dn-karty-pasmo', role: 'tablist', 'aria-label': tekst('panel.sekcje') },
    [
      zakladka('panel-sesje', ikony.karty, tekst('panel.sesje'), true),
      zakladka('panel-projekty', ikony.folder, tekst('panel.projekty'), false),
    ],
  );

  const paskiCzynnosci = el('div', { klasa: 'dn-obszar-pasek' }, [
    przyciskIkony('szukaj', tekst('panel.szukaj')),
    przyciskIkony('plus', tekst('panel.nowaSesja'), { 'data-okno-nowe': '' }),
    przyciskIkony('filtr', tekst('panel.filtry')),
  ]);

  const wykazSesji =
    w.sesje.length > 0
      ? el(
          'div',
          { klasa: 'dn-panel-wykaz', id: 'wykaz-sesji', 'data-sesje-widok': 'wszystkie' },
          w.sesje.map((sesja, indeks) => wierszSesji(sesja, indeks === 0)),
        )
      : el('p', { klasa: 'dn-pusty-stan dn-pusty-stan--zwarty', tekst: tekst('panel.brakSesji') });

  const panelSesjeWezel = el(
    'div',
    { klasa: 'dn-obszar-tresc', id: 'panel-sesje', role: 'tabpanel' },
    [wykazSesji],
  );

  const panelProjektyWezel = el(
    'div',
    { klasa: 'dn-obszar-tresc', id: 'panel-projekty', role: 'tabpanel', hidden: true },
    [el('p', { klasa: 'dn-pusty-stan dn-pusty-stan--zwarty', tekst: tekst('panel.brakProjektow') })],
  );

  return el(
    'aside',
    { klasa: 'dn-obszar-panel dn-obszar-panel--boczny', 'aria-label': tekst('panel.etykieta') },
    [pasmoZakladek, paskiCzynnosci, panelSesjeWezel, panelProjektyWezel],
  );
}
