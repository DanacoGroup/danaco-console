import type { WindowRole } from '../../../shared/contract';
import type { OpisOkna } from '../okno-komunikacji/opis-okna';
import { tytulGniazda } from './etykiety-ukladu';
import type { IdGniazda } from './identyfikatory';

/**
 * Opis okna dla jednego gniazda układu.
 *
 * Jedna odpowiedzialność: wyprowadzenie opisu okna gniazda z opisu wspólnego
 * dla sceny. Okno jest bytem pośrednim między sesją a wiadomością —
 * projekt, katalogi i środowisko wykonania mogą być wspólne na starcie, ale
 * tytuł i rola należą już do gniazda i różnią się między oknami.
 *
 * Opis jest kopią, nie wskazaniem na wspólny obiekt. Zmiana ustawienia jednego
 * okna nie może przeciekać do drugiego: okno jest najwęższym poziomem zasięgu
 * i wygrywa z pozostałymi.
 */
export function opisGniazda(
  podstawa: OpisOkna,
  id: IdGniazda,
  rola: WindowRole,
): OpisOkna {
  return {
    ...podstawa,
    katalogiRobocze: [...podstawa.katalogiRobocze],
    tytul: tytulGniazda(id),
    rola,
  };
}
