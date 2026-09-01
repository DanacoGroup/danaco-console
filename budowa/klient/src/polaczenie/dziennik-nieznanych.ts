import {
  ZDARZENIA_NIEZNANEJ,
  czyKomenda,
  czyZdarzenie,
  type Envelope,
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

/** Dziennik komunikatów nierozpoznanych — brama fail-open: komunikat nie zrywa połączenia ani sesji, zostaje wpis, klient pracuje dalej. */
export interface DziennikNieznanych {
  /** Liczba komunikatów nierozpoznanych od chwili założenia dziennika. */
  liczba(): number;
  /** Ostatni wpis albo `null`, gdy dziennik jest pusty. */
  ostatni(): WpisNieznanego | null;
  /** Subskrypcja wpisów — pozwala warstwie wyższej pokazać, czego rdzeń nie rozpoznał. */
  naWpis(sluchacz: (wpis: WpisNieznanego) => void): Odsubskrybuj;
  /** Odnotowuje odpowiedź rdzenia, która przyszła po terminie korelacji i nie ma już odbiorcy. */
  odnotujSpozniona(koperta: Envelope): void;
  /** Odłącza dziennik od źródła zdarzeń. */
  odlacz(): void;
}

/** Komplet zdarzeń `*.unknown` kontraktu, złożony z mapy zdarzeń zapasowych każdego obszaru nazwy, bez powtórzeń. */
const NIEZNANE: ReadonlySet<string> = new Set<string>(Object.values(ZDARZENIA_NIEZNANEJ));

export function zalozDziennikNieznanych(zrodlo: ZrodloZdarzen): DziennikNieznanych {
  const wpisy = utworzMagistrale<WpisNieznanego>();
  const odsubskrybowania: Odsubskrybuj[] = [];
  let licznik = 0;
  let ostatniWpis: WpisNieznanego | null = null;

  function zapisz(wpis: WpisNieznanego): void {
    licznik += 1;
    ostatniWpis = wpis;
    wpisy.oglos(wpis);
  }

  function zapiszKoperte(typZdarzenia: string, ladunek: unknown): void {
    zapisz({ typZdarzenia, ...odczytajLadunek(ladunek) });
  }

  odsubskrybowania.push(
    zrodlo.naDowolny((koperta) => {
      const typ: string = koperta.type;
      // Odmowa nierozpoznania wraca ze statusem błędu, więc odsiew po statusie zamykałby bramę: rozstrzyga typ.
      if (NIEZNANE.has(typ)) {
        zapiszKoperte(typ, koperta.payload);
        return;
      }
      // Odpowiedź rozstrzyga korelacja żądania; typ znany kontraktowi ma odbiorcę, dziennik zbiera resztę.
      if (koperta.status !== undefined) return;
      if (czyZdarzenie(typ) || czyKomenda(typ)) return;
      zapiszKoperte(typ, koperta.payload);
    }),
  );

  return {
    liczba: () => licznik,
    ostatni: () => ostatniWpis,
    naWpis: (sluchacz) => wpisy.subskrybuj(sluchacz),

    odnotujSpozniona(koperta) {
      zapisz({
        typZdarzenia: koperta.type,
        zadanyTyp: koperta.type,
        idZadania: koperta.id,
        powod: 'odpowiedź po terminie korelacji',
      });
    },

    odlacz() {
      for (const odsubskrybuj of odsubskrybowania.splice(0)) {
        odsubskrybuj();
      }
    },
  };
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
