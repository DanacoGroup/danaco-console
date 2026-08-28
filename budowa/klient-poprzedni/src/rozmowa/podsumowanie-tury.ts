import { liczba, listaTekstow, obiekt, prawda, tekst } from './odczyt-fragmentu';

/**
 * Podsumowanie zakończonej tury kanału głównego, dokładane przez rdzeń do ostatniego
 * fragmentu strumienia oznaczonego polem `done`.
 */
export interface PodsumowanieTury {
  /** Identyfikator rozmowy po stronie kanału; kolejne wywołanie go wznawia. */
  idSesjiKanalu: string;
  /** Podtyp zdarzenia kończącego: `success` albo `error_*`. */
  podtyp: string;
  /** Czy tura zakończyła się błędem. */
  bledna: boolean;
  /** Ostateczna treść odpowiedzi podana przez kanał. */
  tekst: string;
  /** Koszt tury w dolarach. */
  koszt: number;
  /** Liczba tur wewnętrznych kanału. */
  tury: number;
  /** Czas trwania tury w milisekundach. */
  czasMs: number;
  /** Konto, na którym tura faktycznie się wykonała. */
  konto: string;
  /** Typy linii przechwyconych w turze — przejrzystość, nie diagnostyka błędu. */
  typyZdarzen: string[];
  /** Liczba linii wyjścia, których nie dało się odczytać; nie przerywają tury. */
  linieNierozpoznane: number;
}

/**
 * Odczytuje podsumowanie tury z ładunku ostatniego fragmentu kończącego strumień
 * odpowiedzi rdzenia.
 */
export function odczytajPodsumowanie(dane: unknown): PodsumowanieTury | null {
  const zrodlo = obiekt(dane);
  if (zrodlo === null) return null;
  const podtyp = tekst(zrodlo, 'subtype');
  const sesja = tekst(zrodlo, 'cliSessionId');
  // Ładunek bez pola własnego podsumowania nie jest podsumowaniem, tylko domknięciem pustej
  // tury.
  if (podtyp === '' && sesja === '' && !('turns' in zrodlo)) return null;
  return {
    idSesjiKanalu: sesja,
    podtyp,
    bledna: prawda(zrodlo, 'isError'),
    tekst: tekst(zrodlo, 'text'),
    koszt: liczba(zrodlo, 'costUsd'),
    tury: liczba(zrodlo, 'turns'),
    czasMs: liczba(zrodlo, 'durationMs'),
    konto: tekst(zrodlo, 'account'),
    typyZdarzen: listaTekstow(zrodlo, 'eventTypes'),
    linieNierozpoznane: liczba(zrodlo, 'unparsedLines'),
  };
}

/**
 * Podsumowanie tury zapisane jedną linią tekstu, przeznaczone do wyświetlenia w stopce
 * wpisu rozmowy.
 */
export function opisPodsumowania(p: PodsumowanieTury): string {
  const czesci: string[] = [];
  if (p.podtyp.length > 0) czesci.push(p.podtyp);
  if (p.czasMs > 0) czesci.push(`${(p.czasMs / 1000).toFixed(1)} s`);
  if (p.tury > 0) czesci.push(`tur: ${p.tury}`);
  if (p.koszt > 0) czesci.push(`${p.koszt.toFixed(4)} USD`);
  if (p.konto.length > 0) czesci.push(`konto: ${p.konto}`);
  if (p.linieNierozpoznane > 0) czesci.push(`linii nierozpoznanych: ${p.linieNierozpoznane}`);
  return czesci.join(' · ');
}
