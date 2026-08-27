import { przycisk } from '../../modele/kontrolki-formularza';

/**
 * Panel etapów badania wymienia kroki w kolejności, w jakiej zapisał je rdzeń, obok pola
 * redakcji zakresu, nie zamiast niego.
 */
export interface PasekEtapow {
  element: HTMLElement;
  /** Wymienia wykaz etapów na oddany przez rdzeń. */
  odswiez(etapy: readonly string[]): void;
  /** Etap wskazany przez Operatora; pusty napis znaczy „nie wskazano". */
  wskazany(): string;
}

export function utworzPasekEtapow(naWskazanie: (etap: string) => void): PasekEtapow {
  let wskazanyEtap = '';

  const wykaz = document.createElement('ol');
  wykaz.className = 'mr-etapy__wykaz';

  const zastrzezenie = document.createElement('p');
  zastrzezenie.className = 'dn-pole-opis mr-etapy__zastrzezenie';
  zastrzezenie.textContent =
    'Etapy pochodzą z odpowiedzi rdzenia na zapis zakresu. Stanu etapu — ukończony, bieżący, ' +
    'nierozpoczęty — ani checklisty kontrakt nie ma czym przenieść, więc pasek ich nie pokazuje. ' +
    'Wskazanie etapu jest nastawą tego widoku i nie przeżywa ponownego wejścia do modułu.';

  const element = document.createElement('div');
  element.className = 'mr-etapy';
  element.setAttribute('aria-label', 'Etapy badania');

  const naglowek = document.createElement('p');
  naglowek.className = 'dn-pole-etykieta';
  naglowek.textContent = 'Etapy badania';

  element.append(naglowek, wykaz, zastrzezenie);

  return {
    element,
    wskazany: () => wskazanyEtap,

    odswiez(etapy) {
      if (!etapy.includes(wskazanyEtap)) wskazanyEtap = '';
      if (etapy.length === 0) {
        const pusto = document.createElement('li');
        pusto.className = 'mr-etapy__pusto';
        pusto.textContent =
          'Rdzeń nie zapisał ani jednego etapu — wpisz je w polu etapów i zapisz zakres.';
        wykaz.replaceChildren(pusto);
        return;
      }
      wykaz.replaceChildren(
        ...etapy.map((etap, kolejnosc) =>
          pozycjaEtapu(etap, kolejnosc + 1, etap === wskazanyEtap, (wybrany) => {
            wskazanyEtap = wybrany;
            naWskazanie(wybrany);
          }),
        ),
      );
    },
  };
}

/**
 * Jeden etap paska przedstawiony numerem pozycji zamiast znaczka stanu oraz plakietką ze słowną
 * etykietą przy wskazaniu, nigdy samą barwą.
 */
function pozycjaEtapu(
  etap: string,
  numer: number,
  wskazany: boolean,
  naWybor: (etap: string) => void,
): HTMLElement {
  const element = document.createElement('li');
  element.className = 'mr-etapy__pozycja';
  element.dataset['etap'] = etap;
  if (wskazany) element.dataset['wskazany'] = 'tak';

  const kontrolka = przycisk(`${String(numer)}. ${etap}`, 'dn-btn dn-btn--sm dn-btn--duch');
  kontrolka.addEventListener('click', () => naWybor(etap));

  element.append(kontrolka);
  if (wskazany) {
    const znacznik = document.createElement('span');
    znacznik.className = 'dn-plakietka dn-plakietka--informacja';
    znacznik.textContent = 'wskazany';
    element.append(znacznik);
  }
  return element;
}
