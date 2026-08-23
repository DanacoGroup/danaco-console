import { przycisk } from '../../modele/kontrolki-formularza';

/**
 * Pasek etapów badania — wykaz kroków w kolejności, w jakiej zapisał je rdzeń.
 *
 * Etapy przychodzą polem `stages` odpowiedzi `research.workspace.set` i do tej
 * pory nigdzie się nie pokazywały: okno wiodące przyjmowało je w polu tekstowym
 * i odsyłało do rdzenia, a Operator nie widział, co rdzeń naprawdę zapisał.
 * Pasek pokazuje wykaz oddany, nie wpisany — dlatego stoi obok pola redakcji,
 * a nie zamiast niego.
 *
 * **Czego pasek NIE pokazuje i dlaczego.** Opracowanie modułu (rozdz. 3.3)
 * przewiduje przy etapie stan „ukończony · bieżący · nierozpoczęty" wraz
 * z checklistą. Kontrakt niesie etapy jako zwykły wykaz napisów — pola na stan
 * etapu nie ma żadnego. Znaczek ukończenia postawiony tutaj byłby danymi
 * zmyślonymi, więc pasek go nie stawia i mówi wprost, czego brakuje.
 *
 * Wskazanie etapu jest nastawą widoku i tylko nią: rdzeń nie ma gdzie zapisać,
 * na którym etapie stoi badanie, więc po ponownym wejściu do modułu wskazanie
 * zaczyna od zera. Okno tego nie ukrywa.
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
 * Jeden etap paska.
 *
 * Numer pozycji jest jedyną informacją porządkową, jaką kontrakt naprawdę
 * niesie — kolejnością w wykazie — więc to on stoi przy nazwie zamiast znaczka
 * stanu. Wskazanie niesie plakietkę z etykietą słowną obok wyróżnienia
 * graficznego: stan nigdy nie zależy od samej barwy.
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
