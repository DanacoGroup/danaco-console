import type { AssistantActivityEntry } from '../../../../shared/contract';
import { zglosBrak } from './braki-kontraktu';
import { chwila, NAZWY_RODZAJOW } from './etykiety-assistant';

/**
 * Zamiana wpisu dziennika w pozycję listy wraz z trzema czynnościami.
 *
 * Wszystkie trzy mają dziś drogę w kontrakcie i wszystkie trzy naprawdę coś
 * robią: „Odtwórz przebieg" oddaje zamiar oknu, „Odsłuchaj nagranie" pobiera
 * bajty spod odnośnika wpisu (`speech.audio.fetch`) i odtwarza je w karcie,
 * a „Oznacz jako ważne" zapisuje wyróżnienie w rdzeniu
 * (`assistant.activity.flag`) — znacznik przeżywa odświeżenie wykazu, bo ma
 * gdzie zamieszkać.
 *
 * Wpis bez odnośnika nagrania nie jest brakiem produktu, tylko wpisem
 * tekstowym: przycisk odsłuchu mówi to wprost, zamiast milczeć.
 */

/** Zamiar odtworzenia przebiegu zlecenia, do którego należy wpis. */
export type NaPrzebieg = (idZlecenia: string) => void;

/** Czynności wpisu wykonywane przez okno: odsłuch nagrania i wyróżnienie. */
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

/** Trzy czynności wykazu: przebieg, odsłuch nagrania, wyróżnienie wpisu. */
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
