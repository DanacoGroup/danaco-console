import type { WindowRole } from '../../../shared/contract';
import type { OpisOkna } from '../okno-komunikacji/opis-okna';
import { tytulGniazda } from './etykiety-ukladu';
import type { IdGniazda } from './identyfikatory';

/**
 * Funkcja wyprowadza opis okna dla pojedynczego gniazda układu jako kopię opisu wspólnego, z tytułem i rolą właściwymi temu gniazdu.
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
