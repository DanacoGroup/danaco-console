// Szukanie w panelu bocznym: zawęża widoczne wiersze sesji i projektów.
import { oglos } from './ogloszenie.ts';

const ZNACZNIK = '[data-etykietka="Szukaj w oknie"]';
const NAGLOWEK = 'Szukanie';

/* Rdzeń nie ma komendy szukania sesji, więc zakresem są wiersze już wczytane
   do panelu: zawężenie pokazuje mniej tego, co i tak stoi w oknie. */
export function zwiazSzukanieWOknie(): void {
  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const przycisk = cel.closest<HTMLElement>(ZNACZNIK);
    if (przycisk === null) return;
    zdarzenie.preventDefault();
    przelaczPole(przycisk);
  }, true);
}

// Puste pole zdejmuje zawężenie.
function przelaczPole(przycisk: HTMLElement): void {
  const pasek = przycisk.closest<HTMLElement>('.dn-obszar-pasek');
  if (pasek === null) return;
  const stojace = pasek.querySelector<HTMLInputElement>('[data-szukanie-w-oknie]');
  if (stojace !== null) {
    zawez(pasek, '');
    stojace.remove();
    return;
  }
  const pole = pasek.ownerDocument.createElement('input');
  pole.type = 'search';
  pole.className = 'dn-pole dn-pole--sm';
  pole.placeholder = 'Szukaj w oknie';
  pole.setAttribute('aria-label', 'Szukaj w oknie');
  pole.dataset.szukanieWOknie = '';
  pole.addEventListener('input', () => zawez(pasek, pole.value));
  pole.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key === 'Escape') {
      zawez(pasek, '');
      pole.remove();
    }
  });
  przycisk.insertAdjacentElement('afterend', pole);
  pole.focus();
}

/* Porównanie małymi literami, bo Operator nie pisze nazw tak, jak zapisał je
   rdzeń. Brak trafień jest nazwany, nie zostawiony jako pusty panel. */
function zawez(pasek: HTMLElement, szukane: string): void {
  const panel = pasek.parentElement;
  if (panel === null) return;
  const wzorzec = szukane.trim().toLowerCase();
  let trafienia = 0;
  const wiersze = panel.querySelectorAll<HTMLElement>('[data-id-sesji], [data-id-projektu]');
  for (const wiersz of wiersze) {
    const nazwa = (wiersz.textContent ?? '').toLowerCase();
    const pasuje = wzorzec === '' || nazwa.includes(wzorzec);
    wiersz.hidden = !pasuje;
    if (pasuje) trafienia += 1;
  }
  if (wzorzec !== '' && trafienia === 0 && wiersze.length > 0) {
    oglos(NAGLOWEK, 'Żaden wiersz w tym oknie nie niesie szukanego zapisu.', 'ostrzezenie');
  }
}
