/**
 * Pasek narzędzi okna — pas `dn-narzedzia-pas` ze źródła kształtu
 * (`design/05-okna/moduly/studio.html`).
 *
 * Sześć rodzin czynności w kolejności prototypu: układ, nawigacja historii,
 * widok, treść i przechwytywanie, a po prawej kontekst i praca w tle oraz
 * widok i powiadomienia. Pas zamyka przycisk dostosowania wstążki.
 *
 * Czynności, których warstwa okien jeszcze nie prowadzi, stoją zapowiedziane:
 * przycisk zostaje i nazywa niegotowość zdaniem. Znikający przycisk zabrałby
 * Operatorowi wiedzę, że czynność w produkcie w ogóle jest.
 */

import { ikony } from '../ikony.ts';
import { el, tekst, zeZnacznika } from '../narzedzia.ts';

function znak(rysunek: string): SVGElement {
  const wezel = zeZnacznika(rysunek);
  wezel.setAttribute('aria-hidden', 'true');
  return wezel;
}

function przycisk(
  rysunek: keyof typeof ikony,
  etykieta: string,
  atrybuty?: Record<string, string>,
): HTMLElement {
  return el(
    'button',
    {
      klasa: 'dn-nrz-btn dn-etykietka',
      type: 'button',
      'data-etykietka': etykieta,
      'aria-label': etykieta,
      ...atrybuty,
    },
    [znak(ikony[rysunek])],
  );
}

/**
 * Przycisk czynności bez pokrycia w warstwie okien. Niesie `data-komunikat`,
 * który rama zamienia na zdanie dla Operatora — ta sama droga, którą idą
 * pozostałe zapowiedzi okna.
 */
function zapowiedziany(rysunek: keyof typeof ikony, etykieta: string): HTMLElement {
  return przycisk(rysunek, etykieta, {
    'data-komunikat': tekst('narzedzia.zapowiedziane'),
    'data-komunikat-tytul': etykieta,
  });
}

function grupa(dzieci: HTMLElement[], dodatkowaKlasa?: string): HTMLElement {
  return el(
    'div',
    { klasa: dodatkowaKlasa ? `dn-narzedzia-grupa ${dodatkowaKlasa}` : 'dn-narzedzia-grupa' },
    dzieci,
  );
}

export function pasekNarzedzi(): HTMLElement {
  const pole = el('div', { klasa: 'dn-szukaj dn-narzedzia-szukaj' }, [
    znak(ikony.szukaj),
    el('input', {
      type: 'search',
      placeholder: tekst('narzedzia.szukaj'),
      'aria-label': tekst('narzedzia.szukaj'),
    }),
  ]);

  const pas = el('div', { klasa: 'dn-narzedzia' }, [
    // Układ — jedyna czynność pasa działająca dziś w całości.
    grupa([
      przycisk('zwinPanel', tekst('narzedzia.zwinPanel'), {
        'aria-pressed': 'false',
        'data-zwin-szyne': '',
      }),
    ]),
    // Nawigacja historii.
    grupa([
      przycisk('strzalkaLewo', tekst('narzedzia.wstecz')),
      przycisk('strzalkaPrawo', tekst('narzedzia.naprzod')),
    ]),
    // Widok.
    grupa([
      zapowiedziany('godlo', tekst('narzedzia.centrum')),
      przycisk('odswiez', tekst('narzedzia.odswiez')),
    ]),
    // Treść i przechwytywanie.
    grupa([
      zapowiedziany('zrzut', tekst('narzedzia.zrzut')),
      zapowiedziany('schowek', tekst('narzedzia.schowek')),
    ]),
    pole,
    // Kontekst i praca w tle — prawa strona pasa.
    grupa(
      [
        zapowiedziany('kolejka', tekst('narzedzia.kolejka')),
        zapowiedziany('magistrala', tekst('narzedzia.magistrala')),
        zapowiedziany('izolacja', tekst('narzedzia.izolacja')),
      ],
      'dn-narzedzia-grupa--prawa',
    ),
    // Widok i powiadomienia.
    grupa(
      [
        zapowiedziany('powiadomienia', tekst('narzedzia.powiadomienia')),
        zapowiedziany('skupienie', tekst('narzedzia.skupienie')),
        zapowiedziany('pelnyEkran', tekst('narzedzia.pelnyEkran')),
      ],
      'dn-narzedzia-grupa--prawa',
    ),
    // Rozwinięcie kontrolek zwiniętych przy wąskim oknie.
    grupa([zapowiedziany('wiecej', tekst('narzedzia.wiecej'))], 'dn-narzedzia-grupa--prawa'),
    // Dostosowanie wstążki stoi zawsze na końcu pasa — tak stanowi prototyp.
    el(
      'button',
      {
        klasa: 'dn-narzedzia-dostosuj dn-etykietka',
        type: 'button',
        'data-etykietka': tekst('narzedzia.dostosuj'),
        'aria-label': tekst('narzedzia.dostosuj'),
        'data-komunikat': tekst('narzedzia.zapowiedziane'),
        'data-komunikat-tytul': tekst('narzedzia.dostosuj'),
      },
      [znak(ikony.ukladanka)],
    ),
  ]);

  return el('div', { klasa: 'dn-narzedzia-pas' }, [pas]);
}
