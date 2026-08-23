import type { EdytorKrokow } from './edytor-krokow';
import type { StanTresci } from './stany-okna';

/**
 * Wymiana definicji automatyki z plikiem — pozycja „Import/eksport definicji”
 * panelu akcji Workflow Buildera.
 *
 * Wymiana obywa się bez komendy kontraktu: definicja jest już w oknie, a plik
 * leży po stronie klienta. Import wypełnia edytor; do rdzenia trafia dopiero
 * zapis, tą samą komendą co każda inna zmiana.
 */

/** Wczytanie definicji z pliku — import bez pytania rdzenia, bo plik jest lokalny. */
export function wczytajPlik(edytor: EdytorKrokow, tresc: StanTresci): void {
  const wybor = document.createElement('input');
  wybor.type = 'file';
  wybor.accept = 'application/json,.json';
  wybor.addEventListener('change', () => {
    const plik = wybor.files?.[0];
    if (plik === undefined) return;
    void plik.text().then((zawartosc) => {
      try {
        const definicja = JSON.parse(zawartosc) as { steps?: unknown };
        if (!Array.isArray(definicja.steps)) {
          tresc.blad('Plik nie zawiera wykazu kroków — import wstrzymany.');
          return;
        }
        edytor.wczytaj(definicja.steps);
        tresc.potwierdzenie('Kroki wczytane z pliku. Zapis wyśle je do rdzenia.', true);
      } catch {
        tresc.blad('Plik nie jest poprawnym zapisem JSON — import wstrzymany.');
      }
    }).catch(() => {
      tresc.blad('Nie udało się odczytać pliku — import wstrzymany.');
    });
  });
  wybor.click();
}
