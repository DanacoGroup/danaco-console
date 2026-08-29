import {
  ZDARZENIA_NIEZNANEJ,
  czyKomenda,
  czyZdarzenie,
  type EventType,
} from '../../../shared/contract.ts';
import { utworzMagistrale, type Odsubskrybuj } from './magistrala-zdarzen.ts';
import type { ZrodloZdarzen } from './zrodlo-zdarzen.ts';

/** Pojedynczy komunikat, którego nie rozpoznał rdzeń albo klient, niosący typ koperty, typ nierozpoznany, identyfikator żądania oraz powód nierozpoznania. */
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
 * Dziennik komunikatów nierozpoznanych stanowi bramę fail-open połączenia:
 * żaden z nich nie zrywa połączenia ani nie blokuje sesji, wpis idzie do
 * dziennika, a klient pracuje dalej.
 */
export interface DziennikNieznanych {
  /** Liczba komunikatów nierozpoznanych od chwili założenia dziennika. */
  liczba(): number;
  /** Ostatni wpis albo `null`, gdy dziennik jest pusty. */
  ostatni(): WpisNieznanego | null;
  /** Subskrypcja wpisów — pozwala warstwie wyższej pokazać, czego rdzeń nie rozpoznał. */
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
    console.warn('[połączenie] komunikat nierozpoznany', wpis);
    wpisy.oglos(wpis);
  }

  for (const rodzaj of rodzajeNieznanych()) {
    odsubskrybowania.push(
      zrodlo.naZdarzenie(rodzaj, (tresc, koperta) => zapisz(koperta.type, tresc)),
    );
  }

  odsubskrybowania.push(
    zrodlo.naDowolny((koperta) => {
      // Odpowiedź rozstrzyga korelacja żądania; typ znany kontraktowi ma odbiorcę, dziennik zbiera resztę.
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

/** Komplet zdarzeń `*.unknown` kontraktu, złożony z mapy zdarzeń zapasowych każdego obszaru nazwy, bez powtórzeń. */
function rodzajeNieznanych(): EventType[] {
  return [...new Set<EventType>(Object.values(ZDARZENIA_NIEZNANEJ))];
}

/** Odczyt ładunku `UnknownCommandPayload` odporny na jego brak i na inny kształt niż oczekiwany, zwracający wartości domyślne zamiast zgłaszać błąd. */
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

/** Napis albo pusty łańcuch — pole opcjonalne kontraktu bywa nieobecne, więc brak wartości zamienia się na pusty łańcuch zamiast wartości niezdefiniowanej. */
function tekst(wartosc: unknown): string {
  return typeof wartosc === 'string' ? wartosc : '';
}
