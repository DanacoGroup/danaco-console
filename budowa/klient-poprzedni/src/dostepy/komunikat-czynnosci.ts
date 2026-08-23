import { ErrorCode } from '../../../shared/contract';
import type { Wynik } from '../protokol/kanal';

/**
 * Odpowiedź widoku na czynność Operatora — jedno miejsce w sekcji dostępów.
 *
 * Każde naciśnięcie daje odpowiedź. Nadanie dostępu, odebranie go, zmiana trybu
 * i sprawdzenie punktu kończą się zdaniem przy kontrolce — powodzenie mówi, co
 * się stało, niepowodzenie mówi, co odpowiedział rdzeń. Po milczeniu Operator
 * nie wie, czy model ma już dostęp, czy nie.
 */
export interface KomunikatCzynnosci {
  /** Element osadzany pod kontrolką. */
  element: HTMLElement;
  /** Pokazuje zdanie i oznacza je jako powodzenie albo niepowodzenie. */
  pokaz(tresc: string, powodzenie: boolean): void;
  /** Nanosi zdanie wprost z wyniku komendy. */
  zWyniku(wynik: Wynik<unknown>, powodzenieTresc: string): void;
  /** Zdejmuje zdanie — czynność przestała być aktualna. */
  wyczysc(): void;
}

export function utworzKomunikatCzynnosci(klasa = ''): KomunikatCzynnosci {
  const element = document.createElement('p');
  element.className = klasa === '' ? 'dd-komunikat' : `dd-komunikat ${klasa}`;
  element.hidden = true;

  function pokaz(tresc: string, powodzenie: boolean): void {
    element.textContent = tresc;
    element.hidden = tresc === '';
    element.dataset.powodzenie = String(powodzenie);
  }

  return {
    element,
    pokaz,

    zWyniku(wynik, powodzenieTresc) {
      pokaz(wynik.udany ? powodzenieTresc : zdanieBledu(wynik), wynik.udany);
    },

    wyczysc: () => pokaz('', true),
  };
}

/**
 * Odmowa własna widoku — czynność, której nie ma po co wysyłać do rdzenia.
 *
 * Kształt jest ten sam co wyniku komendy, więc widok nie potrzebuje drugiej
 * ścieżki obsługi: odmowa braku okna rozmowy wygląda dla niego jak odmowa
 * rdzenia i tak samo trafia do zdania przy kontrolce.
 */
export function odmowaWlasna(tresc: string): Wynik<never> {
  return {
    udany: false,
    blad: { code: ErrorCode.ValidationFailed, message: tresc, retryable: false },
  };
}

/**
 * Zdanie o niepowodzeniu. Rdzeń, który nie podał treści błędu, nie zostawia
 * pustego miejsca — Operator dostaje przynajmniej kod albo informację o jego
 * braku.
 */
export function zdanieBledu(wynik: Wynik<unknown>): string {
  const blad = wynik.blad;
  if (blad === undefined) return 'Rdzeń nie przyjął czynności i nie podał powodu.';
  const powod = blad.message !== '' ? blad.message : blad.code;
  return `Rdzeń nie przyjął czynności: ${powod}`;
}
