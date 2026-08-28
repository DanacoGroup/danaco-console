import { ErrorCode, LibraryPreviewKind, type LibraryPreview } from '../../../../shared/contract';
import { opisOdmowyBledu, type PowodOdmowy } from '../../komponenty/odmowa';
import type { ZrodloBiblioteki } from './zrodlo-biblioteki';

// Droga odczytu treści pliku jest jedynym miejscem, w którym moduł rozstrzyga jej osiągalność.

/**
 * Werdykt o treści pliku odróżnia brak wiedzy od odmowy rdzenia i od wskazania
 * miejsca treści bez samej treści, każdy stan osobnym słowem.
 */
export type WerdyktTresci = 'nieznana' | 'osiagalna' | 'brak' | 'odmowa' | 'odwolanie';

export interface StanTresci {
  werdykt: WerdyktTresci;
  /** Co rdzeń powiedział: powód odmowy albo opis odpowiedzi bez treści; puste przy stanach bez odmowy. */
  powod: string;
}

/** Stan pliku, o którego treść nikt jeszcze nie pytał — brak wiedzy o treści, a nie stwierdzony brak samej treści. */
export const TRESC_NIEZNANA: StanTresci = { werdykt: 'nieznana', powod: '' };

/**
 * Werdykt złożony z odpowiedzi rdzenia na podgląd: treść brakuje wyłącznie po
 * odmowie walidacji, a osiągalność orzeka treść odpowiedzi, nie samo jej
 * powodzenie.
 */
export function werdyktZPodgladu(
  podglad: LibraryPreview | undefined,
  blad?: PowodOdmowy | null,
): StanTresci {
  if (podglad !== undefined) return werdyktZTresciPodgladu(podglad);
  if (blad?.code === ErrorCode.ValidationFailed) {
    return { werdykt: 'brak', powod: opisOdmowyBledu('Treść pliku', blad) };
  }
  return { werdykt: 'odmowa', powod: opisOdmowyBledu('Odczyt treści pliku', blad) };
}

/**
 * Rozbiór podglądu, który rdzeń oddał: treścią jest wyłącznie podgląd tekstowy,
 * pozostałe rodzaje dają najwyżej wskazanie miejsca, którego klient nie ma czym
 * pobrać.
 */
function werdyktZTresciPodgladu(podglad: LibraryPreview): StanTresci {
  if (podglad.kind === LibraryPreviewKind.Text && podglad.text !== undefined) {
    return { werdykt: 'osiagalna', powod: '' };
  }
  const odnosnik = (podglad.imageRef ?? '').trim();
  if (odnosnik !== '') {
    return {
      werdykt: 'odwolanie',
      powod:
        `rdzeń oddał podgląd rodzaju „${podglad.kind}" ze wskazaniem miejsca treści ` +
        `(${odnosnik}), a nie samą treść — czy leży ona pod tym wskazaniem, okno nie ma ` +
        'jak sprawdzić',
    };
  }
  return {
    werdykt: 'odwolanie',
    powod:
      `rdzeń oddał podgląd rodzaju „${podglad.kind}" bez treści i bez wskazania jej ` +
      'miejsca — kontrakt nie niesie w nim bajtów',
  };
}

/**
 * Pyta rdzeń o pierwszą stronę podglądu wyłącznie po to, żeby poznać werdykt.
 *
 * Innej drogi kontrakt nie daje: nie ma komendy pytającej o samą obecność
 * treści, więc świadkiem jest ta, która treść oddaje albo jej odmawia.
 */
export async function zbadajTresc(
  zrodlo: ZrodloBiblioteki,
  idPliku: string,
): Promise<StanTresci> {
  const wynik = await zrodlo.podglad(idPliku, 1);
  return werdyktZPodgladu(wynik.udany ? wynik.wynik?.preview : undefined, wynik.blad);
}
