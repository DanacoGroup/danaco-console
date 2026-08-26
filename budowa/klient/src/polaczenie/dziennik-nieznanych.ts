import {
  ZDARZENIA_NIEZNANEJ,
  czyKomenda,
  czyZdarzenie,
  type EventType,
} from '../../../shared/contract.ts';
import { utworzMagistrale, type Odsubskrybuj } from './magistrala-zdarzen.ts';
import type { ZrodloZdarzen } from './zrodlo-zdarzen.ts';

/** Pojedynczy komunikat, którego nie rozpoznał rdzeń albo klient. */
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
 * Dziennik komunikatów nierozpoznanych — brama fail-open połączenia.
 *
 * Rdzeń odpowiada na nieznaną komendę zdarzeniem `*.unknown` właściwym dla
 * obszaru nazwy; każdy obszar kontraktu ma własne zdarzenie zapasowe. Dziennik
 * obejmuje je wszystkie, sięgając po komplet z mapy kontraktu
 * `ZDARZENIA_NIEZNANEJ`, nigdy po literał nazwy: dopisanie obszaru
 * w `contract.json` rozszerza dziennik samo, bez zmiany tego pliku.
 *
 * Osobno przechwytujemy koperty o typie spoza kontraktu — takie, których nie
 * zna ani wykaz komend, ani wykaz zdarzeń. Powstają, gdy rdzeń wyprzedził
 * klienta wersją albo gdy ramka była nieczytelna. Żaden z tych przypadków nie
 * zrywa połączenia i nie blokuje sesji: wpis idzie do dziennika, klient
 * pracuje dalej.
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
      // Odpowiedź rozstrzyga korelacja żądania, a typ znany kontraktowi ma
      // swojego odbiorcę — dziennik zbiera wyłącznie resztę.
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

/** Komplet zdarzeń `*.unknown` kontraktu, bez powtórzeń. */
function rodzajeNieznanych(): EventType[] {
  return [...new Set<EventType>(Object.values(ZDARZENIA_NIEZNANEJ))];
}

/** Odczyt ładunku `UnknownCommandPayload` odporny na jego brak i na inny kształt. */
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

/** Napis albo pusty łańcuch — pole opcjonalne kontraktu bywa nieobecne. */
function tekst(wartosc: unknown): string {
  return typeof wartosc === 'string' ? wartosc : '';
}
