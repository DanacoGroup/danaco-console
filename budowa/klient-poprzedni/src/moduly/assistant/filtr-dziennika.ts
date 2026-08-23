import type { AssistantActivityEntry } from '../../../../shared/contract';
import { pole, wybor } from '../../modele/kontrolki-formularza';
import { POZYCJE_FILTRA_DZIENNIKA, type KodFiltraDziennika } from './etykiety-assistant';

/**
 * Wyszukiwanie i filtr rodzaju wpisu Activity Feed.
 *
 * Zawężenie liczy się w oknie. `assistant.activity.list` przyjmuje wyłącznie
 * `windowId`, `actionId` i `limit` (`adapter_modul_asystent_czynnosci.go`),
 * więc rdzeń nie ma czym zawęzić wykazu po treści ani po rodzaju wpisu —
 * odczyt na każde naciśnięcie klawisza byłby wywołaniem zwracającym za każdym
 * razem to samo.
 *
 * Wyszukiwanie jest pełnotekstowe i tylko takie. Wyszukiwania po znaczeniu
 * w dzienniku asystenta kontrakt nie ma: `knowledge.search` sięga biblioteki,
 * historii rozmów i plików przestrzeni roboczej (`KnowledgeScope`), a dziennik
 * asystenta nie jest żadnym z tych trzech zakresów. Brak nazywa okno wprost
 * (`braki-kontraktu.ts`), zamiast podstawiać dopasowanie po literach pod nazwę
 * „semantyczne".
 *
 * Porównanie idzie po zwinięciu wielkości liter właściwym dla polszczyzny
 * (`toLocaleLowerCase('pl')`), więc „Ł" znajduje „ł".
 */
export interface FiltrDziennika {
  /** Pasek osadzany w nagłówku okna: pole szukania i wybór rodzaju. */
  element: HTMLElement;
  /** Czy wpis przechodzi przez bieżące zawężenie. */
  przepusc(wpis: AssistantActivityEntry): boolean;
  /** Czy cokolwiek jest zawężone — rozstrzyga zdanie pustki okna. */
  zawezony(): boolean;
}

export function utworzFiltrDziennika(naZmiane: () => void): FiltrDziennika {
  const szukanie = pole('Szukaj w historii działań', 'szukaj w treści wpisów');
  szukanie.type = 'search';
  szukanie.addEventListener('input', naZmiane);

  const rodzaj = wybor('Filtr rodzaju wpisu', POZYCJE_FILTRA_DZIENNIKA);
  rodzaj.value = 'wszystkie';
  rodzaj.addEventListener('change', naZmiane);

  const element = document.createElement('div');
  element.className = 'ma-dziennik__filtr';
  element.append(szukanie, rodzaj);

  const fraza = (): string => szukanie.value.trim().toLocaleLowerCase('pl');
  const kod = (): KodFiltraDziennika => rodzaj.value as KodFiltraDziennika;

  return {
    element,
    przepusc(wpis) {
      const szukana = fraza();
      if (szukana !== '' && !wpis.content.toLocaleLowerCase('pl').includes(szukana)) return false;
      const wybrany = kod();
      return wybrany === 'wszystkie' || wpis.kind === wybrany;
    },
    zawezony: () => fraza() !== '' || kod() !== 'wszystkie',
  };
}
