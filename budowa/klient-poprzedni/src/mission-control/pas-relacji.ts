import {
  ETYKIETA_BRAKU_ZRODLA,
  ETYKIETA_RODZAJU_RELACJI,
  STRZALKA_RELACJI,
} from './etykiety-pulpitu';
import type { RelacjaSesji } from './model-danych';

/**
 * Pas relacji pod matrycą sesji: powiązania między sesjami, także sesjami
 * z różnych środowisk. Kolumny matrycy pokazują sesje obok siebie, a dopiero pas
 * pokazuje, że procesy biegną razem.
 */
export interface PasRelacji {
  element: HTMLElement;
  odswiez(relacje: RelacjaSesji[] | null): void;
}

/**
 * Buduje pas relacji wraz z jego wykazem i zwraca odświeżenie. Odświeżenie
 * przyjmuje wykaz powiązań albo `null`, więc pas obsługuje zarówno brak powiązań,
 * jak i brak odczytu, bez dwóch osobnych dróg wywołania.
 */
export function utworzPasRelacji(relacje: RelacjaSesji[] | null): PasRelacji {
  const element = document.createElement('div');
  element.className = 'mc-relacje';

  const etykieta = document.createElement('span');
  etykieta.className = 'mc-relacje__etykieta';
  etykieta.textContent = 'Relacje';

  const wykaz = document.createElement('ul');
  wykaz.className = 'mc-relacje__wykaz';
  wykaz.setAttribute('aria-label', 'Powiązania między sesjami');

  element.append(etykieta, wykaz);

  const odswiez = (dane: RelacjaSesji[] | null): void => {
    if (dane === null) {
      wykaz.replaceChildren(brakZrodla());
      return;
    }
    wykaz.replaceChildren(...(dane.length > 0 ? dane.map(wiersz) : [brakRelacji()]));
  };
  odswiez(relacje);

  return { element, odswiez };
}

/**
 * Jedno powiązanie: źródło, znak kierunku, cel i przedmiot wymiany. Znak jest
 * podwójny — strzałka oraz słowo w podpowiedzi i w treści czytanej — żeby kierunek
 * nie zależał od samego kształtu.
 */
function wiersz(relacja: RelacjaSesji): HTMLLIElement {
  const element = document.createElement('li');
  element.className = 'mc-relacja';
  element.dataset.rodzaj = relacja.rodzaj;
  element.title = `${ETYKIETA_RODZAJU_RELACJI[relacja.rodzaj]} — ${relacja.przedmiot}`;

  const zrodlo = document.createElement('span');
  zrodlo.className = 'mc-relacja__wezel';
  zrodlo.textContent = relacja.zrodloNazwa;
  zrodlo.dataset.sesja = relacja.zrodloId;

  const znak = document.createElement('span');
  znak.className = 'mc-relacja__znak';
  znak.textContent = STRZALKA_RELACJI[relacja.rodzaj];
  znak.setAttribute('aria-label', ETYKIETA_RODZAJU_RELACJI[relacja.rodzaj]);
  znak.setAttribute('role', 'img');

  const cel = document.createElement('span');
  cel.className = 'mc-relacja__wezel';
  cel.textContent = relacja.celNazwa;
  cel.dataset.sesja = relacja.celId;

  const przedmiot = document.createElement('span');
  przedmiot.className = 'mc-relacja__przedmiot';
  przedmiot.textContent = relacja.przedmiot;

  element.append(zrodlo, znak, cel, przedmiot);
  return element;
}

/**
 * Pas nie znika przy braku powiązań — mówi wprost, że żadna sesja nie jest
 * powiązana z inną. Zniknięcie pasa byłoby nierozróżnialne od braku odczytu,
 * a to dwa różne stany.
 */
function brakRelacji(): HTMLLIElement {
  const element = document.createElement('li');
  element.className = 'mc-relacja mc-relacja--brak';
  element.textContent = 'Żadna sesja nie jest powiązana z inną.';
  return element;
}

/**
 * Pas nie znika też przy braku źródła — mówi, że odczytu dziś nie ma. Podpowiedź
 * nazywa przyczynę: kontrakt nie niesie komendy pytającej o powiązania między
 * sesjami, a brakujący odczyt jest zgłoszony.
 */
function brakZrodla(): HTMLLIElement {
  const element = document.createElement('li');
  element.className = 'mc-relacja mc-relacja--brak';
  element.textContent = `Powiązania sesji — ${ETYKIETA_BRAKU_ZRODLA}.`;
  element.title =
    'Kontrakt nie niesie dziś odczytu powiązań między sesjami. Brakujący odczyt zgłoszony.';
  return element;
}
