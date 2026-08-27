/**
 * Odczyt sprawcy z pól `actor` oraz `actorClientId` koperty zdarzenia. Odczyt,
 * nie domysł: gdy rdzeń nazwał sprawcę, warstwa mówi jego nazwą, a gdy pola nie
 * ma, mówi wprost, że nie wiadomo.
 */
import { ActorKind } from '../../../shared/contract';

/**
 * Co warstwa asystenta wie o ręce, która wykonała czynność. Wykaz rozdziela
 * przypadki nazwane przez rdzeń od przypadku, w którym koperta zdarzenia pola
 * sprawcy nie niosła.
 */
export type ZnanySprawca =
  /** Rdzeń nazwał: asystent działający za Operatora. */
  | { rodzaj: 'asystent' }
  /** Rdzeń nazwał: Operator — z tego samego klienta, co ten dymek. */
  | { rodzaj: 'operator-tutaj' }
  /** Rdzeń nazwał: Operator — ale z innego klienta (drugie okno, telefon, AOD). */
  | { rodzaj: 'operator-gdzie-indziej' }
  /** Rdzeń nazwał: Operator, a my nie znamy własnego identyfikatora klienta. */
  | { rodzaj: 'operator-nieznane-urzadzenie' }
  /** Rdzeń nazwał: model w oknie roboczym, narzędziem przez MCP. */
  | { rodzaj: 'model' }
  /** Rdzeń nazwał sam siebie: przemiatanie, harmonogram, odtworzenie stanu. */
  | { rodzaj: 'rdzen' }
  /** Koperta pola nie niosła — rdzeń nie potrafił rozstrzygnąć. */
  | { rodzaj: 'nieznany' };

/**
 * Kształt wspólny siedmiu zdarzeniom kontraktu, które sprawcę już niosą: pole
 * rodzaju sprawcy oraz identyfikator klienta, z którego czynność wyszła. Oba
 * pola są opcjonalne.
 */
export interface KopertaZeSprawca {
  actor?: ActorKind;
  actorClientId?: string;
}

/**
 * Rozstrzyga sprawcę wyłącznie z treści zdarzenia.
 *
 * @param idKlienta identyfikator tego połączenia. Pusty napis znaczy, że własny
 *   identyfikator nie jest znany, więc Operatora nie przypisuje się ani temu
 *   urządzeniu, ani innemu.
 */
export function rozpoznajSprawce(
  tresc: KopertaZeSprawca | undefined,
  idKlienta: string,
): ZnanySprawca {
  const actor = tresc?.actor;
  if (actor === undefined) return { rodzaj: 'nieznany' };

  switch (actor) {
    case ActorKind.Assistant:
      return { rodzaj: 'asystent' };
    case ActorKind.Model:
      return { rodzaj: 'model' };
    case ActorKind.Core:
      return { rodzaj: 'rdzen' };
    case ActorKind.Operator: {
      const czyj = tresc?.actorClientId ?? '';
      if (idKlienta === '' || czyj === '') return { rodzaj: 'operator-nieznane-urzadzenie' };
      return czyj === idKlienta
        ? { rodzaj: 'operator-tutaj' }
        : { rodzaj: 'operator-gdzie-indziej' };
    }
    default:
      // Wartość spoza wykazu `ActorKind` nie jest błędem klienta, lecz brakiem
      // rozstrzygnięcia.
      return { rodzaj: 'nieznany' };
  }
}

/**
 * Krótka nazwa ręki, która wykonała czynność, przeznaczona do etykiety wiersza
 * oraz do czytnika ekranu; każdemu rozpoznanemu rodzajowi sprawcy odpowiada
 * dokładnie jedno zdanie.
 */
export function nazwijSprawce(sprawca: ZnanySprawca): string {
  switch (sprawca.rodzaj) {
    case 'asystent':
      return 'Asystent';
    case 'operator-tutaj':
      return 'Operator — z tego okna';
    case 'operator-gdzie-indziej':
      return 'Operator z innego urządzenia';
    case 'operator-nieznane-urzadzenie':
      return 'Operator — nie wiadomo, z którego urządzenia';
    case 'model':
      return 'Model w oknie roboczym';
    case 'rdzen':
      return 'Rdzeń';
    case 'nieznany':
      return 'Nie wiadomo, czyja ręka';
  }
}

/**
 * Zdanie dopisywane do opisu posunięcia. Dla sprawcy nieznanego mówi wprost, że
 * koperta pola nie niosła, żeby czytelnik nie wziął tego za niepewność
 * interfejsu.
 */
export function zdanieOSprawcy(sprawca: ZnanySprawca): string {
  if (sprawca.rodzaj === 'nieznany') {
    return 'nie wiadomo, czyja ręka — to zdarzenie nie niesie pola sprawcy';
  }
  return nazwijSprawce(sprawca);
}

/**
 * Rozstrzyga, czy czynność wykonał asystent działający za Operatora; pytanie
 * pada osobno, ponieważ dymek asystenta znakuje własne posunięcia inaczej niż
 * posunięcia cudze.
 */
export function czyAsystent(sprawca: ZnanySprawca): boolean {
  return sprawca.rodzaj === 'asystent';
}
