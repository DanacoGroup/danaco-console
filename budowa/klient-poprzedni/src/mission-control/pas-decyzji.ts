import { elementIkony } from '../ikony/ikony';
import { ETYKIETA_BRAKU_ZRODLA } from './etykiety-pulpitu';
import type { PasDecyzji } from './model-danych';
import type { ZamiarDecyzji } from './zdarzenia-pulpitu';

/**
 * Pas eskalacji koordynatora ma jedną odpowiedzialność: pokazuje, że pętla
 * zatrzymała się na człowieku, i podaje drogę wyjścia. Zdarzenie nieoczekiwane
 * wstrzymuje kolejkę, a praca czeka na decyzję użytkownika.
 */
export interface PasDecyzjiWidok {
  element: HTMLElement;
  odswiez(decyzje: PasDecyzji | null): void;
}

/**
 * Buduje pas decyzji wraz z wezwaniem „Podejmij decyzję": znak ostrzegawczy,
 * nagłówek z liczbą przepływów, szczegół z czasem oczekiwania oraz przycisk
 * nadający zamiar rozstrzygnięcia.
 */
export function utworzPasDecyzji(
  decyzje: PasDecyzji | null,
  nadaj: (zamiar: ZamiarDecyzji) => void,
): PasDecyzjiWidok {
  const element = document.createElement('div');
  element.className = 'mc-decyzje';
  element.setAttribute('role', 'status');

  const znak = document.createElement('span');
  znak.className = 'mc-decyzje__znak';
  znak.append(elementIkony('ostrzezenie', { rozmiar: 24, etykieta: 'Wymagana decyzja' }));

  const tresc = document.createElement('div');
  tresc.className = 'mc-decyzje__tresc';

  const naglowek = document.createElement('p');
  naglowek.className = 'mc-decyzje__naglowek';

  const szczegol = document.createElement('p');
  szczegol.className = 'mc-decyzje__szczegol';

  tresc.append(naglowek, szczegol);

  const wezwanie = document.createElement('button');
  wezwanie.type = 'button';
  wezwanie.className = 'dn-btn dn-btn--atrament mc-decyzje__wezwanie';
  wezwanie.textContent = 'Podejmij decyzję';

  let biezace = decyzje;
  wezwanie.addEventListener('click', () => {
    if (biezace === null) return;
    nadaj({ przeplywy: biezace.przeplywyDoDecyzji, najstarszy: biezace.najstarszyPrzeplyw });
  });

  element.append(znak, tresc, wezwanie);

  const odswiez = (dane: PasDecyzji | null): void => {
    biezace = dane;
    if (dane === null) {
      naglowek.textContent = `Przepływy do decyzji — ${ETYKIETA_BRAKU_ZRODLA}`;
      szczegol.textContent =
        'Kontrakt nie niesie dziś odczytu przepływów wstrzymanych do decyzji — brakujący odczyt zgłoszony.';
      element.dataset.oczekuje = 'nie';
      wezwanie.hidden = true;
      return;
    }
    wezwanie.hidden = false;
    const liczba = dane.przeplywyDoDecyzji;
    naglowek.textContent =
      liczba > 0
        ? `${liczba} ${odmianaPrzeplywu(liczba)}`
        : 'Żaden przepływ nie czeka na decyzję';
    szczegol.textContent =
      liczba > 0
        ? `Najdłużej czeka: ${dane.najstarszyPrzeplyw} — ${dane.czekaMinut} min. Kolejka wstrzymana do rozstrzygnięcia.`
        : 'Pętla koordynator–wykonawca biegnie bez wstrzymania.';
    element.dataset.oczekuje = liczba > 0 ? 'tak' : 'nie';
  };
  odswiez(decyzje);

  return { element, odswiez };
}

/**
 * Odmiana orzeczenia i rzeczownika przez liczbę:
 * „1 przepływ wymaga decyzji", „3 przepływy wymagają decyzji",
 * „5 przepływów wymaga decyzji".
 */
function odmianaPrzeplywu(liczba: number): string {
  if (liczba === 1) {
    return 'przepływ wymaga decyzji';
  }
  const dziesiatki = liczba % 100;
  const jednosci = liczba % 10;
  const mnoga = jednosci >= 2 && jednosci <= 4 && !(dziesiatki >= 12 && dziesiatki <= 14);
  return mnoga ? 'przepływy wymagają decyzji' : 'przepływów wymaga decyzji';
}
