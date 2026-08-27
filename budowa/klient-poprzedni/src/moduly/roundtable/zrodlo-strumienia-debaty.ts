import { EventType, type ChunkKind, type Envelope, type StreamChunkEvent } from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { Kanal } from '../../protokol/kanal';

/** Fragment wypowiedzi debaty przyjęty z okna, niosący identyfikator koperty, rodzaj, treść tekstową i znacznik ostatniego fragmentu strumienia. */
export interface FragmentDebaty {
  /** Pole messageId koperty, surowe: identyfikator uczestnika albo wypowiedzi wedle rodzaju. */
  identyfikator: string;
  /** Rodzaj fragmentu z kontraktu — tekst, tok rozumowania, błąd. */
  rodzaj: ChunkKind;
  /** Treść tekstowa; pusta, gdy fragment jej nie niósł. */
  tekst: string;
  /** `done` koperty — prawda w ostatnim fragmencie strumienia tego głosu. */
  ostatni: boolean;
}

export interface ZrodloStrumieniaDebaty {
  /** Subskrypcja fragmentów okna debaty; zwrócona funkcja zdejmuje subskrypcję po zamknięciu panelu. */
  naFragmentWypowiedzi(sluchacz: (fragment: FragmentDebaty) => void): Odsubskrybuj;
}

/** Tworzy źródło fragmentów wypowiedzi debaty, filtrując zdarzenia strumienia rdzenia do okna odczytywanego w chwili nadejścia każdego fragmentu. */
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

/** Zwraca fragment przeznaczony dla okna debaty albo null, gdy okno jest puste lub fragment należy do innej debaty niż wskazana. */
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
