import { EventType, type ChunkKind, type Envelope, type StreamChunkEvent } from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { Kanal } from '../../protokol/kanal';

/**
 * Źródło fragmentów wypowiedzi debaty — jedna subskrypcja `stream.chunk`
 * zawężona do okna debaty.
 *
 * Rdzeń nadaje wypowiedź uczestnika na żywo, fragment po fragmencie, zanim
 * utrwali ją i rozgłosi po raz drugi jako `roundtable.debate.changed` rodzaju
 * `updated`. Ten plik odbiera tamtą nadawaną treść i nic ponadto.
 *
 * Filtr idzie po oknie, bo tylko okno jest w zdarzeniu pewne: `stream.chunk`
 * niesie `windowId` i `messageId` (kontrakt, `StreamChunkEvent`), a rdzeń wpisuje
 * w `windowId` okno debaty. Okno oddziela więc głosy tej debaty od strumieni okna
 * rozmowy, podglądu w tle i każdego innego nadawcy wspólnej drogi.
 *
 * `messageId` zostaje surowy. Jeden strumień niesie pod tym polem dwie różne
 * wartości: fragmenty treści dostają identyfikator uczestnika, a fragment
 * domykający i fragment błędu — identyfikator wypowiedzi. Źródło nie rozstrzyga
 * tej dwoistości i nie zgaduje po przedrostku identyfikatora; przypisanie do
 * uczestnika robi `strumien-wypowiedzi.ts`, pytając o ten identyfikator stan
 * debaty, czyli jedyny byt znający i skład, i wykaz wypowiedzi tury.
 *
 * Źródło nie woła ani jednej komendy — komendy obszaru `roundtable.*` niesie
 * `zrodlo-roundtable.ts`. Nie gromadzi też treści: gromadzenie wymaga wiedzy
 * o składzie, a subskrypcja nie.
 */

/** Fragment wypowiedzi przyjęty z okna debaty. */
export interface FragmentDebaty {
  /**
   * `messageId` koperty, surowy. Bywa identyfikatorem uczestnika (fragmenty
   * treści) albo identyfikatorem wypowiedzi (fragment domykający i błąd).
   */
  identyfikator: string;
  /** Rodzaj fragmentu z kontraktu — tekst, tok rozumowania, błąd. */
  rodzaj: ChunkKind;
  /** Treść tekstowa; pusta, gdy fragment jej nie niósł. */
  tekst: string;
  /** `done` koperty — prawda w ostatnim fragmencie strumienia tego głosu. */
  ostatni: boolean;
}

export interface ZrodloStrumieniaDebaty {
  /**
   * Subskrypcja `stream.chunk` zawężona do okna debaty.
   *
   * Zwracana funkcja zdejmuje subskrypcję. Panel, który jej nie wywoła, zostawia
   * nasłuch żywy po zejściu ze sceny — dlatego `PanelPomocniczy` wymaga `zamknij()`.
   */
  naFragmentWypowiedzi(sluchacz: (fragment: FragmentDebaty) => void): Odsubskrybuj;
}

/**
 * @param okno odczyt okna debaty w chwili nadejścia fragmentu, nie w chwili
 *   subskrypcji: `StanDebaty.ustawOkno` przestawia moduł na inną debatę bez
 *   zakładania panelu od nowa, a filtr zapamiętany przy subskrypcji
 *   przepuszczałby wtedy fragmenty debaty poprzedniej.
 */
export function utworzZrodloStrumieniaDebaty(
  kanal: Kanal,
  okno: () => string,
): ZrodloStrumieniaDebaty {
  return {
    naFragmentWypowiedzi(sluchacz) {
      return kanal.naZdarzenie(EventType.StreamChunk, (tresc, koperta) => {
        const fragment = fragmentDlaOkna(tresc, koperta, okno());
        if (fragment === null) return;
        sluchacz(fragment);
      });
    },
  };
}

/**
 * Fragment przeznaczony dla okna debaty albo `null`.
 *
 * Okno puste odrzuca wszystko. Gospodarz, który nie dostał okna od rdzenia,
 * podaje `okno: ''`; przepuszczenie wtedy całego ruchu `stream.chunk` wsypałoby
 * do panelu debaty fragmenty okna rozmowy i podglądu w tle, czyli cudze zdania
 * podpisane uczestnikami tej debaty.
 */
function fragmentDlaOkna(
  tresc: StreamChunkEvent,
  koperta: Envelope,
  okno: string,
): FragmentDebaty | null {
  if (okno === '' || tresc.windowId !== okno) return null;
  return {
    identyfikator: tresc.messageId,
    rodzaj: tresc.kind,
    tekst: tresc.text ?? '',
    ostatni: koperta.done === true,
  };
}
