import {
  ZDARZENIA_NIEZNANEJ,
  czyKomenda,
  czyZdarzenie,
  type EventType,
} from '../../../shared/contract';
import { utworzMagistrale, type Odsubskrybuj } from './magistrala-zdarzen';
import type { ZrodloZdarzen } from './zrodlo-zdarzen';

/** Pojedynczy komunikat, którego nie rozpoznał rdzeń albo klient, zapisany w dzienniku nierozpoznanych zdarzeń. */
export interface WpisNieznanego {
  /** Typ koperty, która przyszła: zdarzenie `*.unknown` albo typ spoza kontraktu. */
  typZdarzenia: string;
  /** Typ, którego rdzeń nie rozpoznał; pusty, gdy rdzeń go nie podał. */
  zadanyTyp: string;
  /** Identyfikator nierozpoznanego żądania; pusty, gdy rdzeń go nie podał. */
  idZadania: string;
  /** Powód nierozpoznania; pusty, gdy rdzeń go nie podał. */
  powod: string;
}

/**
 * Dziennik komunikatów nierozpoznanych jest bramą fail-open łączności: obejmuje wszystkie zdarzenia zapasowe kontraktu oraz koperty o typie spoza kontraktu, nie zrywając połączenia ani nie blokując sesji.
 */
export interface DziennikNieznanych {
  /** Liczba komunikatów nierozpoznanych od chwili założenia dziennika. */
  liczba(): number;
  /** Ostatni wpis albo `null`, gdy dziennik jest pusty. */
  ostatni(): WpisNieznanego | null;
  /** Subskrypcja wpisów — pozwala widokowi pokazać, że rdzeń czegoś nie rozpoznał. */
  naWpis(sluchacz: (wpis: WpisNieznanego) => void): Odsubskrybuj;
  /** Odłącza dziennik od źródła zdarzeń. */
  odlacz(): void;
}

export function zalozDziennikNieznanych(zrodlo: ZrodloZdarzen): DziennikNieznanych {
  const wpisy = utworzMagistrale<WpisNieznanego>();
  const odsubskrybowania: Odsubskrybuj[] = [];
  let licznik = 0;
  let ostatniWpis: WpisNieznanego | null = null;

  function zapisz(typZdarzenia: string, ladunek: unknown): void {
    const wpis: WpisNieznanego = { typZdarzenia, ...odczytajLadunek(ladunek) };
    licznik += 1;
    ostatniWpis = wpis;
    console.warn('[łączność] komunikat nierozpoznany', wpis);
    wpisy.oglos(wpis);
  }

  for (const rodzaj of rodzajeNieznanych()) {
    odsubskrybowania.push(
      zrodlo.naZdarzenie(rodzaj, (tresc, koperta) => zapisz(koperta.type, tresc)),
    );
  }

  odsubskrybowania.push(
    zrodlo.naDowolny((koperta) => {
      // Odpowiedź rozstrzyga korelacja żądania, a znany typ ma odbiorcę — dziennik zbiera resztę.
      if (koperta.status !== undefined) return;
      const typ: string = koperta.type;
      if (czyZdarzenie(typ) || czyKomenda(typ)) return;
      zapisz(typ, koperta.payload);
    }),
  );

  return {
    liczba: () => licznik,
    ostatni: () => ostatniWpis,
    naWpis: (sluchacz) => wpisy.subskrybuj(sluchacz),

    odlacz() {
      for (const odsubskrybuj of odsubskrybowania.splice(0)) {
        odsubskrybuj();
      }
    },
  };
}

/** Komplet zdarzeń *.unknown kontraktu, bez powtórzeń, wykorzystywany przy subskrypcji dziennika nierozpoznanych. */
function rodzajeNieznanych(): EventType[] {
  return [...new Set<EventType>(Object.values(ZDARZENIA_NIEZNANEJ))];
}

/** Odczyt ładunku UnknownCommandPayload odporny na jego brak i na inny kształt niż oczekiwany przez dziennik. */
function odczytajLadunek(ladunek: unknown): Omit<WpisNieznanego, 'typZdarzenia'> {
  if (typeof ladunek !== 'object' || ladunek === null) {
    return { zadanyTyp: '', idZadania: '', powod: '' };
  }
  const pola = ladunek as Record<string, unknown>;
  return {
    zadanyTyp: tekst(pola['requestedType']),
    idZadania: tekst(pola['requestId']),
    powod: tekst(pola['reason']),
  };
}

/** Napis albo pusty łańcuch — pole opcjonalne kontraktu bywa całkiem nieobecne w treści komunikatu od rdzenia. */
function tekst(wartosc: unknown): string {
  return typeof wartosc === 'string' ? wartosc : '';
}
