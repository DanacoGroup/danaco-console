import {
  Command,
  ImageAdjustKind,
  ImageTransformKind,
  type ArchiveUnpackResponse,
  type ImageAdjustRequest,
  type ImageTransformRequest,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyLiczba, czyObiekt, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Przygotowanie materiału wejściowego przed rozpoznaniem pisma — trzy komendy
 * spoza obszaru `studio`.
 *
 * Opracowanie modułu żąda przed rozpoznaniem prostowania skosu, odszumiania,
 * podniesienia kontrastu i przycięcia marginesów oraz wczytywania wsadowego
 * z archiwum. Obszar `studio` nie niesie żadnej z tych czynności, ale kontrakt
 * niesie je gdzie indziej i mają one uchwyt w rdzeniu: poprawka obrazu idzie
 * `image.adjust`, przekształcenie geometryczne `image.transform`,
 * a rozpakowanie wsadu `archive.unpack`. Panel woła je wprost, zamiast nazywać
 * brakiem czynność, którą rdzeń umie.
 *
 * Każda poprawka obrazu oddaje NOWY zasób magazynu, a źródło zostaje nietknięte.
 * Pozycja kolejki przestawia się więc na zasób wynikowy i od tej chwili
 * rozpoznanie pisma pracuje na obrazie po czyszczeniu.
 */

/** Poprawka obrazu wykonalna przed rozpoznaniem pisma. */
export interface PoprawkaObrazu {
  /** Kod pozycji; trafia do `data-poprawka` przycisku. */
  kod: string;
  /** Etykieta widoczna dla Operatora. */
  nazwa: string;
  /** Zdanie o tym, co poprawka robi z obrazem i po co przed rozpoznaniem. */
  opis: string;
  rodzaj: ImageAdjustKind;
}

/**
 * Cztery poprawki wskazane przez opracowanie i obecne w kontrakcie.
 *
 * Binaryzacji w wykazie nie ma i nie ma jej w kontrakcie — skala szarości wraz
 * z wyrównaniem poziomów jest tym, co rdzeń w tej drodze umie, i tak jest
 * nazwana. Zrównanie jej z binaryzacją byłoby obietnicą progowania, którego
 * nikt nie wykonuje.
 */
export const POPRAWKI_OBRAZU: readonly PoprawkaObrazu[] = [
  {
    kod: 'odszumienie',
    nazwa: 'Odszum',
    opis: 'Usuwa ziarno skanu; rozpoznanie pisma myli ziarno z interpunkcją.',
    rodzaj: ImageAdjustKind.Denoise,
  },
  {
    kod: 'kontrast',
    nazwa: 'Kontrast',
    opis: 'Podnosi kontrast; blady skan daje litery zlewające się z tłem.',
    rodzaj: ImageAdjustKind.Contrast,
  },
  {
    kod: 'szarosc',
    nazwa: 'Skala szarości',
    opis: 'Sprowadza obraz do szarości — barwa papieru nie niesie treści pisma.',
    rodzaj: ImageAdjustKind.Grayscale,
  },
  {
    kod: 'poziomy',
    nazwa: 'Wyrównaj poziomy',
    opis: 'Rozciąga poziomy jasności; zastępuje binaryzację, której kontrakt nie zna.',
    rodzaj: ImageAdjustKind.AutoLevels,
  },
];

/** Wskazanie materiału dla poprawki: zasób magazynu albo ścieżka rdzenia. */
export interface WskazanieObrazu {
  idZasobu: string;
  sciezka: string;
}

export interface ZrodloMaterialuStudio {
  /** Wykonuje poprawkę obrazu; oddaje identyfikator nowego zasobu. */
  popraw(
    idOkna: string,
    wskazanie: WskazanieObrazu,
    rodzaj: ImageAdjustKind,
    sila: number,
  ): Promise<Wynik<{ idZasobu: string }>>;
  /** Prostuje skos obrotem o zadany kąt; oddaje identyfikator nowego zasobu. */
  obroc(
    idOkna: string,
    wskazanie: WskazanieObrazu,
    stopnie: number,
  ): Promise<Wynik<{ idZasobu: string }>>;
  /** Rozpakowuje archiwum wsadu i oddaje ścieżki jego pozycji. */
  rozpakuj(idOkna: string, sciezka: string): Promise<Wynik<ArchiveUnpackResponse>>;
}

export function utworzZrodloMaterialuStudio(kanal: Kanal): ZrodloMaterialuStudio {
  /**
   * Dopisuje wskazanie materiału do żądania. Zasób magazynu ma pierwszeństwo
   * przed ścieżką: jego treść stoi już pod sumą kontrolną, więc rdzeń nie
   * wciąga jej po raz drugi.
   */
  function wskaz(
    zadanie: { assetId?: string; sourcePath?: string },
    wskazanie: WskazanieObrazu,
  ): void {
    if (wskazanie.idZasobu !== '') zadanie.assetId = wskazanie.idZasobu;
    else zadanie.sourcePath = wskazanie.sciezka;
  }

  return {
    async popraw(idOkna, wskazanie, rodzaj, sila) {
      const zadanie: ImageAdjustRequest = { operation: rodzaj };
      wskaz(zadanie, wskazanie);
      if (idOkna !== '') zadanie.windowId = idOkna;
      // Siła wchodzi wyłącznie podana: brak pola znaczy „weź wartość rozsądną
      // dla operacji", co jest rozstrzygnięciem rdzenia, a nie okna.
      if (sila > 0) zadanie.amount = sila;
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.ImageAdjust, zadanie),
        Command.ImageAdjust,
        (tresc) => czyObiekt(tresc.asset),
      );
      return przepiszZasob(wynik);
    },

    async obroc(idOkna, wskazanie, stopnie) {
      const zadanie: ImageTransformRequest = {
        operation: ImageTransformKind.Rotate,
        degrees: stopnie,
      };
      wskaz(zadanie, wskazanie);
      if (idOkna !== '') zadanie.windowId = idOkna;
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.ImageTransform, zadanie),
        Command.ImageTransform,
        (tresc) => czyObiekt(tresc.asset),
      );
      return przepiszZasob(wynik);
    },

    async rozpakuj(idOkna, sciezka) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.ArchiveUnpack, {
          sourcePath: sciezka,
          ...(idOkna === '' ? {} : { windowId: idOkna }),
        }),
        Command.ArchiveUnpack,
        (tresc) => czyLiczba(tresc.entries),
      );
    },
  };
}

/**
 * Sprowadza dwie odpowiedzi o różnym kształcie do jednej wartości, której
 * potrzebuje kolejka: identyfikatora zasobu wynikowego. Panel nie ogląda
 * pozostałych pól zasobu, więc ich przenoszenie tylko rozszerzałoby styk.
 */
function przepiszZasob(
  wynik: Wynik<{ asset: { id: string } }>,
): Wynik<{ idZasobu: string }> {
  if (!wynik.udany || wynik.wynik === undefined) {
    return wynik.blad === undefined ? { udany: false } : { udany: false, blad: wynik.blad };
  }
  return { udany: true, wynik: { idZasobu: wynik.wynik.asset.id } };
}
