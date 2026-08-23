import {
  poleLogiczne,
  poleTekstowe,
  poleWielowierszowe,
  type PoleFormularza,
} from '../../modele/kontrolki-formularza';
import { OBJASNIENIA, PODPOWIEDZ_JEZYKOW } from './etykiety-translate';
import { dopnijDymek, podepnijPodpowiedz } from './kontrolki-translate';

/**
 * Pola formularza Source Panel wraz z dymkami objaśnień [?].
 *
 * Wydzielone z okna, bo skład pola — etykieta, podpowiedź, dymek, wykaz
 * podpowiadanych wartości — to inna odpowiedzialność niż to, co się dzieje po
 * naciśnięciu przycisku. Okno składa całość i wiąże zdarzenia; ten plik zna
 * wyłącznie kształt formularza.
 *
 * Dymek objaśnienia dostaje każde pole zmieniające treść żądania — język
 * źródłowy i ponowna segmentacja.
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
