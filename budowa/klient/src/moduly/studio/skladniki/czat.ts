/**
 * Strefa 4 — okno czatu. Belka, nagłówek, kontekst, historia i dół z polem
 * polecenia — kształt `.sta-czaty > .sta-kom` z prototypu. Rdzeń nie prowadzi
 * dziś rozmowy tego okna (osobny zakres prac), więc nagłówek, kontekst
 * i historia niosą nazwany stan pusty zamiast parametrów i treści zmyślonych.
 * Pole polecenia stoi gotowe do wpięcia, samo wysłanie na razie nie ma dokąd
 * pójść, więc przechwytuje `submit` bez skutku.
 */

import { ikony, type NazwaZnaku } from '../ikony.ts';
import { el, tekst, zeZnacznika } from '../narzedzia.ts';

function znak(rysunek: NazwaZnaku): SVGElement {
  const wezel = zeZnacznika(ikony[rysunek]);
  wezel.setAttribute('aria-hidden', 'true');
  return wezel;
}

function belka(): HTMLElement {
  return el('header', { klasa: 'sta-okno-belka' }, [
    el('span', { klasa: 'sta-okno-tytul' }, [znak('dymek'), el('b', { tekst: tekst('czat.tytul') })]),
  ]);
}

function dol(): HTMLElement {
  const pole = el('textarea', {
    klasa: 'sta-prompt-obszar',
    rows: 1,
    placeholder: tekst('czat.zastepczaTresc'),
    'aria-label': tekst('czat.etykietaTresci'),
  });
  const wyslijBtn = el(
    'button',
    { klasa: 'dn-btn dn-btn--atrament dn-btn--sm sta-prompt-wyslij', type: 'submit', 'aria-label': tekst('czat.wyslij') },
    [znak('wyslij')],
  );
  const formularz = el('form', { klasa: 'sta-prompt' }, [
    el('span', { klasa: 'sta-prompt-grot', 'aria-hidden': 'true', tekst: '»»' }),
    pole,
    wyslijBtn,
  ]);
  // Bez wpiętej rozmowy wysłanie nie ma dokąd pójść — przechwycone, żeby strona się nie przeładowała.
  formularz.addEventListener('submit', (zdarzenie) => zdarzenie.preventDefault());
  return el('div', { klasa: 'sta-kom-dol' }, [formularz]);
}

export function oknoCzatu(): HTMLElement {
  const kom = el(
    'section',
    { klasa: 'sta-okno sta-kom', id: 'okno-czat-1', 'data-nazwa': 'Chat Window', 'aria-label': tekst('czat.tytul') },
    [
      belka(),
      el('div', { klasa: 'sta-kom-naglowek' }, [el('span', { klasa: 'dn-meta', tekst: tekst('czat.brakParametrow') })]),
      el('div', { klasa: 'sta-kom-kontekst', 'aria-label': tekst('czat.etykietaKontekst') }, [
        el('span', { klasa: 'dn-meta', tekst: tekst('czat.brakKontekstu') }),
      ]),
      el('div', { klasa: 'sta-kom-historia' }, [
        el('p', { klasa: 'dn-pusty-stan dn-pusty-stan--zwarty', tekst: tekst('czat.brakWiadomosci') }),
      ]),
      dol(),
    ],
  );

  return el('div', { klasa: 'sta-czaty', 'data-liczba': '1' }, [kom]);
}
