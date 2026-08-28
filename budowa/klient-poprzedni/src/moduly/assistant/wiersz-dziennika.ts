/**
 * Zamiana wpisu dziennika w pozycję listy wraz z trzema czynnościami:
 * odtworzeniem przebiegu zlecenia, odsłuchem nagrania i wyróżnieniem wpisu.
 * Każda z nich ma drogę w kontrakcie i wykonuje pracę, a nie tylko wygląda.
 */

import type { AssistantActivityEntry } from '../../../../shared/contract';
import { zglosBrak } from './braki-kontraktu';
import { chwila, NAZWY_RODZAJOW } from './etykiety-assistant';

/**
 * Zamiar odtworzenia przebiegu zlecenia, do którego należy wpis. Wykonanie
 * należy do okna dziennika, bo to ono wie, gdzie przebieg pokazać; wiersz
 * wyłącznie zgłasza numer zlecenia.
 */
export type NaPrzebieg = (idZlecenia: string) => void;

/**
 * Czynności wpisu wykonywane przez okno: odsłuch nagrania i wyróżnienie.
 * Wiersz nie woła rdzenia sam — oddaje czynność oknu, które trzyma kanał
 * i odpowiada za odświeżenie wykazu po zapisie.
 */
export interface CzynnosciWpisu {
  /** `speech.audio.fetch` — pobranie bajtów nagrania spod odnośnika wpisu. */
  odsluch(odnosnik: string): void;
  /** `assistant.activity.flag` — wyróżnienie wpisu albo jego zdjęcie. */
  wyroznienie(idWpisu: string, wazny: boolean): void;
}

export function wierszDziennika(
  wpis: AssistantActivityEntry,
  wybrany: boolean,
  naPrzebieg: NaPrzebieg,
  czynnosciWpisu: CzynnosciWpisu,
): HTMLElement {
  const element = document.createElement('li');
  element.className = 'ma-dziennik__wpis';
  element.dataset['wpis'] = wpis.id;
  element.dataset['wybrany'] = wybrany ? 'tak' : 'nie';
  element.dataset['wazny'] = wpis.important === true ? 'tak' : 'nie';

  const czolo = document.createElement('p');
  czolo.className = 'ma-dziennik__czolo';
  czolo.textContent = `${NAZWY_RODZAJOW[wpis.kind]} · ${chwila(wpis.createdAt)}`;

  const tresc = document.createElement('p');
  tresc.className = 'ma-dziennik__tresc';
  tresc.textContent = wpis.content;

  element.append(czolo, tresc, czynnosci(wpis, naPrzebieg, czynnosciWpisu));
  return element;
}

/**
 * Trzy czynności wykazu: przebieg, odsłuch nagrania, wyróżnienie wpisu.
 * Przycisk odsłuchu przy wpisie bez odnośnika nagrania mówi wprost, że wpis
 * jest tekstowy, zamiast milczeć albo znikać.
 */
function czynnosci(
  wpis: AssistantActivityEntry,
  naPrzebieg: NaPrzebieg,
  czynnosciWpisu: CzynnosciWpisu,
): HTMLElement {
  const panel = document.createElement('div');
  panel.className = 'ma-dziennik__akcje';

  panel.append(
    guzik('Odtwórz przebieg', () => {
      if (wpis.actionId === undefined || wpis.actionId === '') {
        zglosBrak(
          'Odtworzenie przebiegu',
          'Ten wpis nie wskazuje zlecenia (pole actionId jest puste), więc nie ma czego odtworzyć.',
        );
        return;
      }
      naPrzebieg(wpis.actionId);
    }),
    guzik('Odsłuchaj nagranie', () => {
      if (wpis.audioRef === undefined || wpis.audioRef === '') {
        zglosBrak(
          'Odsłuch nagrania',
          'Ten wpis nie niesie odnośnika nagrania (pole audioRef jest puste) — ' +
            'polecenie zostało wydane tekstem, więc nie ma czego odsłuchać.',
        );
        return;
      }
      czynnosciWpisu.odsluch(wpis.audioRef);
    }),
    guzik(wpis.important === true ? 'Zdejmij wyróżnienie' : 'Oznacz jako ważne', () =>
      czynnosciWpisu.wyroznienie(wpis.id, wpis.important !== true),
    ),
  );
  return panel;
}

function guzik(etykieta: string, naKlik: () => void): HTMLButtonElement {
  const element = document.createElement('button');
  element.type = 'button';
  element.className = 'dn-btn dn-btn--sm dn-btn--duch';
  element.textContent = etykieta;
  element.addEventListener('click', naKlik);
  return element;
}
