import {
  Command,
  ErrorCode,
  EventType,
  type RequestOf,
  type ResponseOf,
  type UnknownCommandPayload,
} from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../../protokol/kanal';

/**
 * Interfejs odpowiedzi badania rozszerza wynik wywołania o rozpoznanie odmowy nieznanej komendy, gdy rdzeń nie zna zadanego typu żądania.
 */
export interface OdpowiedzBadania<T> extends Wynik<T> {
  /** Typ, którego rdzeń nie zna; pusty, gdy odpowiedź nie jest odmową nieznanej. */
  nieznanyTyp?: string;
}

export interface NasluchOdmow {
  /** Wysyła komendę kontraktu; odmowa nieznanej rozstrzyga obietnicę tak samo jak odpowiedź. */
  wyslij<K extends Command>(
    komenda: K,
    zadanie: RequestOf<K>,
  ): Promise<OdpowiedzBadania<ResponseOf<K>>>;
  /** Odłącza subskrypcje zdarzeń odmowy. */
  rozlacz(): void;
}

/** Stała wylicza zdarzenia odmowy istotne dla okien modułu Research: odmowę własnego obszaru oraz odmowę obszaru okna. */
const ZDARZENIA_ODMOWY = [EventType.ResearchUnknown, EventType.WindowUnknown] as const;

export function utworzNasluchOdmow(kanal: Kanal): NasluchOdmow {
  const oczekujace = new Map<string, (tresc: UnknownCommandPayload) => void>();
  const odsubskrybowania: Odsubskrybuj[] = ZDARZENIA_ODMOWY.map((rodzaj) =>
    kanal.naZdarzenie(rodzaj, (tresc, koperta) => {
      const identyfikator = tresc?.requestId ?? koperta.id;
      const odbiorca = oczekujace.get(identyfikator);
      if (odbiorca === undefined) return;
      oczekujace.delete(identyfikator);
      odbiorca(tresc);
    }),
  );

  return {
    wyslij(komenda, zadanie) {
      return new Promise((rozstrzygnij) => {
        let rozstrzygnieto = false;
        const podaj = (odpowiedz: OdpowiedzBadania<ResponseOf<typeof komenda>>): void => {
          if (rozstrzygnieto) return;
          rozstrzygnieto = true;
          rozstrzygnij(odpowiedz);
        };
        const identyfikator = kanal.wyslij(komenda, zadanie, (wynik) => {
          oczekujace.delete(identyfikator);
          podaj(wynik);
        });
        oczekujace.set(identyfikator, (tresc) => podaj(odmowa(tresc, komenda)));
      });
    },

    rozlacz() {
      for (const odsubskrybuj of odsubskrybowania.splice(0)) odsubskrybuj();
      oczekujace.clear();
    },
  };
}

/**
 * Funkcja zamienia odmowę nieznanej komendy w zwykłe niepowodzenie wywołania, zachowując żądany typ w polu osobnym odpowiedzi.
 */
function odmowa<T>(tresc: UnknownCommandPayload, komenda: string): OdpowiedzBadania<T> {
  const zadany = tresc?.requestedType ?? komenda;
  const powod = (tresc?.reason ?? '').trim();
  return {
    udany: false,
    nieznanyTyp: zadany,
    blad: {
      code: ErrorCode.NotFound,
      message:
        powod === ''
          ? `Rdzeń nie zna komendy ${zadany}`
          : `Rdzeń nie zna komendy ${zadany} — ${powod}`,
      retryable: false,
    },
  };
}
