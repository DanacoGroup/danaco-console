import {
  poleLogiczne,
  poleTekstowe,
  poleWielowierszowe,
  type PoleFormularza,
} from '../../modele/kontrolki-formularza';
import { OBJASNIENIA, PODPOWIEDZ_JEZYKOW } from './etykiety-translate';
import { dopnijDymek, podepnijPodpowiedz } from './kontrolki-translate';

/**
 * Pola formularza panelu źródła wraz z dymkami objaśnień, budujące wyłącznie
 * kształt formularza, bez wiązania zdarzeń należącego do okna.
 */
export interface PolaZrodlaFormularza {
  tekst: PoleFormularza<HTMLTextAreaElement>;
  jezyk: PoleFormularza<HTMLInputElement>;
  podpowiedzi: HTMLDataListElement;
  ponowna: PoleFormularza<HTMLInputElement>;
  /** Wskaźnik zaznaczenia — fragment jako przedmiot operacji. */
  zaznaczenie: HTMLElement;
}

export function zbudujPolaZrodla(): PolaZrodlaFormularza {
  const tekst = poleWielowierszowe(
    {
      etykieta: 'Tekst źródłowy',
      podpowiedz: 'wpisz albo wklej treść do przetłumaczenia',
      opis: 'Zapis tekstu uruchamia tłumaczenie we wszystkich panelach języków naraz.',
    },
    8,
  );

  const jezyk = poleTekstowe({
    etykieta: 'Język źródłowy',
    podpowiedz: 'puste = rozpoznanie przez rdzeń',
  });
  dopnijDymek(jezyk.element, OBJASNIENIA.jezykZrodlowy);
  const podpowiedzi = podepnijPodpowiedz(jezyk.kontrolka, 'mt-jezyki-zrodlowe', PODPOWIEDZ_JEZYKOW);

  const ponowna = poleLogiczne({ etykieta: 'Podziel na segmenty ponownie przy zapisie' });
  dopnijDymek(ponowna.element, OBJASNIENIA.ponownaSegmentacja);

  const zaznaczenie = document.createElement('p');
  zaznaczenie.className = 'mt-zrodlo__zaznaczenie';

  return { tekst, jezyk, podpowiedzi, ponowna, zaznaczenie };
}
