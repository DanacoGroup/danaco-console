import type { StudioDiffCompareRequest } from '../../../../shared/contract';
import { poleLogiczne, poleTekstowe } from '../../modele/kontrolki-formularza';
import type { ParaPorownania, StanStudio } from './stan-studio';

/**
 * Pola wejściowe Diff/Grep Panelu wraz ze złożeniem żądania.
 *
 * Wydzielone od widoku: panel odpowiada za układ i stany, ten plik za to, jak
 * cztery kontrolki składają się w treść żądania `studio.diff.compare`. Kształt
 * żądania można dzięki temu zmienić bez dotykania układu okna.
 *
 * Puste pole wersji porównywanej nie znaczy „brak danych": kontrakt czyta brak
 * `targetVersionId` jako zgodę na porównanie z propozycją zmiany.
 *
 * Wersja bieżąca wchodzi do żądania wprost, bo rdzeń jej nie podstawia sam.
 * `studio.diff.compare` przeszukuje treść wskazanej strony porównania, więc
 * żądanie z samym `pattern` wraca puste nawet wtedy, gdy dokument wzorzec
 * zawiera. Dokument niesie `versionId`, więc okno podstawia go za puste pole
 * odniesienia i podpowiedź przy kontrolce odpowiada temu, co się dzieje.
 */
export interface PolaRoznicy {
  /** Kontrolki w kolejności osadzenia w oknie. */
  elementy: readonly HTMLElement[];
  /** Treść żądania porównania; `null`, gdy nie ma dokumentu czynnego. */
  zadanie(stan: StanStudio): StudioDiffCompareRequest | null;
  /**
   * Wpisuje parę wersji wskazaną w Session Repository.
   *
   * Wpis idzie do kontrolek, a nie obok nich: Operator ma zobaczyć, co zostanie
   * porównane, i móc to jeszcze zmienić przed naciśnięciem. Para wskazana
   * z zewnątrz jest podpowiedzią, nie rozkazem.
   */
  przyjmijPare(para: ParaPorownania): void;
}

export function utworzPolaRoznicy(): PolaRoznicy {
  const wersjaOdniesienia = poleTekstowe({
    etykieta: 'Wersja odniesienia (baseVersionId)',
    podpowiedz: 'puste = wersja bieżąca dokumentu',
    opis:
      'Pole puste podstawia wersję bieżącą dokumentu. Dokument jeszcze niezapisany wersji nie ma, ' +
      'a rdzeń przeszukuje treść WERSJI — wtedy wyszukiwanie nie ma czego przeszukać i mówi o tym wprost.',
  });
  const wersjaPorownywana = poleTekstowe({
    etykieta: 'Wersja porównywana (targetVersionId)',
    podpowiedz: 'puste = propozycja z Tools Panel',
  });
  const wzorzec = poleTekstowe({
    etykieta: 'Wzorzec wyszukiwania',
    podpowiedz: 'fraza albo wyrażenie regularne',
    opis: 'Wzorzec i porównanie wersji można podać razem — rdzeń odda trafienia i fragmenty w jednej odpowiedzi.',
  });
  const regularne = poleLogiczne({ etykieta: 'Wzorzec jest wyrażeniem regularnym' });

  return {
    elementy: [
      wersjaOdniesienia.element,
      wersjaPorownywana.element,
      wzorzec.element,
      regularne.element,
    ],

    przyjmijPare(para) {
      wersjaOdniesienia.kontrolka.value = para.odniesienie;
      wersjaPorownywana.kontrolka.value = para.porownywana;
    },

    zadanie(stan) {
      const dokument = stan.dokument();
      if (dokument === null) return null;
      const tresc: StudioDiffCompareRequest = { documentId: dokument.id };
      const odniesienie = wersjaOdniesienia.kontrolka.value.trim();
      const porownywana = wersjaPorownywana.kontrolka.value.trim();
      const szukana = wzorzec.kontrolka.value.trim();
      if (odniesienie !== '') tresc.baseVersionId = odniesienie;
      else if (dokument.versionId !== undefined && dokument.versionId !== '') {
        tresc.baseVersionId = dokument.versionId;
      }
      if (porownywana !== '') tresc.targetVersionId = porownywana;
      const propozycja = stan.propozycja();
      if (porownywana === '' && propozycja !== null && propozycja.idPropozycji !== '') {
        tresc.proposalId = propozycja.idPropozycji;
      }
      if (szukana !== '') {
        tresc.pattern = szukana;
        tresc.regex = regularne.kontrolka.checked;
      }
      return tresc;
    },
  };
}
