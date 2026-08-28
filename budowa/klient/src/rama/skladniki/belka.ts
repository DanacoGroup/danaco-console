/**
 * Strefa 1 — belka tytułowa. Znak i nazwa produktu po lewej, tytuł bieżącego
 * widoku pośrodku, sterowanie oknem po prawej. Przyciski sterowania stoją
 * w znaczniku tej strefy, jak w prototypie; czynność, którą wykonują na
 * oknie systemowym, nadaje powłoka Tauri — poza tym terenem.
 */

import { ikony } from '../ikony.ts';
import { el, tekst, zeZnacznika } from '../narzedzia.ts';

export interface WlasciwosciBelki {
  /** Nazwa środowiska, do którego przebieg wszedł — jedyny fragment tytułu, który nie stoi w katalogu treści. */
  srodowisko: string;
}

/** Belka wraz z odrębnym uchwytem do węzła tytułu — wywołujący zmienia tytuł bez ponownego odpytywania drzewa. */
export interface Belka {
  wezel: HTMLElement;
  tytul: HTMLElement;
}

export function belka(w: WlasciwosciBelki): Belka {
  const znak = zeZnacznika(ikony.godlo);
  znak.setAttribute('aria-hidden', 'true');

  const tytul = el('span', {
    klasa: 'dn-belka-tytul',
    tekst: `${tekst('belka.marka')} ${tekst('belka.separator')} ${w.srodowisko}`,
    'data-belka-tytul': true,
  });

  function przyciskOkna(rysunek: keyof typeof ikony, etykieta: string, dodatkowaKlasa?: string): HTMLElement {
    const wezelZnaku = zeZnacznika(ikony[rysunek]);
    wezelZnaku.setAttribute('aria-hidden', 'true');
    return el(
      'button',
      {
        klasa: dodatkowaKlasa ? `dn-belka-btn ${dodatkowaKlasa} dn-etykietka` : 'dn-belka-btn dn-etykietka',
        type: 'button',
        'data-etykietka': etykieta,
        'aria-label': etykieta,
      },
      [wezelZnaku],
    );
  }

  const sterowanieOknem = el('span', { klasa: 'dn-belka-okno' }, [
    przyciskOkna('minimalizuj', tekst('belka.minimalizuj')),
    przyciskOkna('maksymalizuj', tekst('belka.maksymalizuj')),
    przyciskOkna('zamknij', tekst('belka.zamknij'), 'dn-belka-btn--zamknij'),
  ]);

  const wezel = el('header', { klasa: 'dn-belka' }, [
    el('span', { klasa: 'dn-belka-marka' }, [
      znak,
      el('span', { klasa: 'dn-belka-nazwa', tekst: tekst('belka.marka') }),
    ]),
    tytul,
    sterowanieOknem,
  ]);

  return { wezel, tytul };
}
