import { ActorKind } from '../../../shared/contract';

/**
 * Odczyt sprawcy z pól `actor` / `actorClientId` koperty zdarzenia.
 *
 * Odczyt, nie domysł: gdy rdzeń nazwał sprawcę, warstwa mówi jego nazwą; gdy
 * pola nie ma, mówi „nie wiadomo". Napis nie bywa pewniejszy niż dowód, na
 * którym stoi.
 *
 * Osobne od `aplikacja/rozstrzyganie-sprawcy.ts`, który rozwiązuje inne
 * zadanie — „czy zrobiło to bieżące połączenie" dla zdarzeń bez pola `actor`.
 */

/** Co warstwa asystenta wie o ręce, która wykonała czynność. */
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

/** Kształt wspólny siedmiu zdarzeniom, które sprawcę już niosą. */
export interface KopertaZeSprawca {
  actor?: ActorKind;
  actorClientId?: string;
}

/**
 * Rozstrzyga sprawcę WYŁĄCZNIE z treści zdarzenia.
 *
 * @param idKlienta identyfikator TEGO połączenia (`uzgodnienie.klient.id`).
 *   Pusty napis znaczy „nie znamy własnego" — wtedy Operatora nie przypisujemy
 *   ani temu urządzeniu, ani innemu, bo nie ma czym porównać.
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
      // Wartość spoza wykazu `ActorKind` nie jest błędem klienta. Warstwa nie
      // ma dla niej nazwy, więc traktuje ją jak brak rozstrzygnięcia.
      return { rodzaj: 'nieznany' };
  }
}

/** Krótka nazwa ręki — do etykiety wiersza i do czytnika ekranu. */
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
 * Zdanie dopisywane do opisu posunięcia.
 *
 * Dla `nieznany` mówi wprost, że pola nie było — inaczej czytelnik wziąłby
 * „nie wiadomo" za niepewność interfejsu, podczas gdy jest to brak w kopercie
 * zdarzenia.
 */
export function zdanieOSprawcy(sprawca: ZnanySprawca): string {
  if (sprawca.rodzaj === 'nieznany') {
    return 'nie wiadomo, czyja ręka — to zdarzenie nie niesie pola sprawcy';
  }
  return nazwijSprawce(sprawca);
}

/** Czy czynność wykonał asystent. */
export function czyAsystent(sprawca: ZnanySprawca): boolean {
  return sprawca.rodzaj === 'asystent';
}
