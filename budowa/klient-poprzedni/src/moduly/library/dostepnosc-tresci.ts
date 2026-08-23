import { ErrorCode, LibraryPreviewKind, type LibraryPreview } from '../../../../shared/contract';
import { opisOdmowyBledu, type PowodOdmowy } from '../../komponenty/odmowa';
import type { ZrodloBiblioteki } from './zrodlo-biblioteki';

/**
 * Droga odczytu treści pliku — jedyne miejsce, w którym moduł Library
 * rozstrzyga, czy repozytorium rdzenia ma treść wskazanego pliku.
 *
 * `library.file.upload` odpowiada powodzeniem także wtedy, gdy bajtów nie ma
 * gdzie zapisać, więc samo `status: "ok"` nie jest świadkiem wgrania. Treść
 * bajtów trzyma magazyn rdzenia (`adapter_modul_library_magazyn.go`).
 *
 * Kontrakt nie ma pola osiągalności treści, a metryka o niej nie mówi:
 * `LibraryFile` niesie `checksum` i `versionId`, lecz plik może mieć obie
 * wartości i nie mieć treści, a wersja czytana spod odwołania przychodzi bez
 * sumy kontrolnej w ogóle. Jedynym świadkiem jest odpowiedź rdzenia na
 * `library.file.preview` — i to ją moduł zapisuje, nigdy własny domysł.
 */

/**
 * Co rdzeń powiedział o treści pliku.
 *
 * `nieznana` jest stanem osobnym od `brak`: „nikt jeszcze nie pytał" to co innego
 * niż „rdzeń odmówił treści".
 *
 * `odmowa` jest osobna, bo rdzeń odmawia podglądu trzema kodami i tylko jeden
 * z nich orzeka o treści:
 *   — `validation_failed`: „plik … nie ma odwołania do treści — podgląd
 *     niemożliwy" — jedyna odpowiedź orzekająca o treści, i tylko po niej
 *     wychodzi `brak`;
 *   — `not_found`: „plik nie istnieje: …" — zdanie o pliku, nie o jego treści
 *     (tym samym kodem straż odmów nazywa komendę bez uchwytu w rdzeniu);
 *   — `internal_error`: „treść nie leży pod odwołaniem …", `retryable: true` —
 *     nośnik odmówił odczytu.
 * Klientowy sprawdzian kształtu (`protokol/ksztalt-odpowiedzi.ts`) sam wytwarza
 * `validation_failed`, gdy odpowiedź nie ma zapowiedzianego pola, więc i taka
 * odpowiedź trafia na `brak`; rozróżnienia nie ma po czym zrobić bez znaku
 * postawionego po stronie sprawdzianu, a ten plik leży poza modułem.
 *
 * `odwolanie` oddziela wskazanie miejsca treści od samej treści: podgląd pliku,
 * którego odwołanie wskazuje ścieżkę nieistniejącą, wraca powodzeniem tak samo
 * jak podgląd pliku leżącego w magazynie rdzenia. Świadkiem treści jest więc sama
 * treść odpowiedzi — podgląd tekstowy z polem `text`; bajtów spod odnośnika klient
 * nie ma czym pobrać, więc podgląd z samym odnośnikiem dostaje własny werdykt.
 */
export type WerdyktTresci = 'nieznana' | 'osiagalna' | 'brak' | 'odmowa' | 'odwolanie';

export interface StanTresci {
  werdykt: WerdyktTresci;
  /**
   * Co rdzeń powiedział: powód odmowy albo opis odpowiedzi, która treści nie
   * niosła. Pusty przy `nieznana` i przy `osiagalna`.
   */
  powod: string;
}

/** Stan pliku, o którego treść nikt jeszcze nie pytał — brak wiedzy, nie brak treści. */
export const TRESC_NIEZNANA: StanTresci = { werdykt: 'nieznana', powod: '' };

/**
 * Werdykt złożony z odpowiedzi rdzenia na podgląd.
 *
 * Wydzielony, bo werdykt powstaje w dwóch miejscach: przy pytaniu zadanym
 * wprost o treść (`zbadajTresc`) i przy zwykłym podglądzie okna File Preview,
 * które i tak tę odpowiedź dostaje — oba mają nazywać ten stan tak samo.
 *
 * Brak treści orzeka rdzeń: `brak` wychodzi stąd wyłącznie po
 * `validation_failed`, jedynej odmowie mówiącej o treści (wykaz kodów przy
 * `WerdyktTresci`). Każdą inną odmowę i każde milczenie moduł nazywa odmową
 * i powtarza jej powód.
 *
 * Osiągalność orzeka treść odpowiedzi, nie jej powodzenie: podgląd bierze się
 * z `preview`, a nie z samego `udany`, bo rdzeń odpowiada powodzeniem także
 * wtedy, gdy zamiast treści oddaje odnośnik do niej.
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
 * Rozbiór podglądu, który rdzeń oddał: co w nim jest, a co tylko wskazuje.
 *
 * Treścią jest wyłącznie podgląd tekstowy z polem `text` — jedyny rodzaj,
 * w którym kontrakt niesie bajty. Trzy pozostałe (`image`, `pdf`, `binary`)
 * dają najwyżej `imageRef`, a klient nie ma komendy, którą pobrałby cokolwiek
 * spod odnośnika (`tresc-podgladu.ts`). Zdanie o nich nazywa więc dokładnie to,
 * co przyszło, i nie dopowiada, czy pod odnośnikiem coś leży.
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
