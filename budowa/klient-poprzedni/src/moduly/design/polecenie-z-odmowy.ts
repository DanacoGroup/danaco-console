import type { ErrorInfo } from '../../../../shared/contract';

/**
 * Odczyt treści polecenia ocalałej w odmowie generowania.
 *
 * Jedyne zadanie pliku: wyjęcie gotowej treści polecenia z `error.details`
 * odmowy `design.asset.generate`.
 *
 * Komenda kończy się błędem, gdy w rejestrze brakuje kanału obrazowego, kanał
 * jest nieczynny albo wskazany kanał okazuje się tekstowy, brakuje poświadczenia,
 * odpowiedź nie niesie obrazu lub bajtów nie da się pobrać. Każda taka odmowa
 * niesie gotową treść polecenia w `details.polecenie`, bo rdzeń składa prompt
 * w całości, zanim cokolwiek wyśle. Okno pokazuje ten tekst, żeby praca nad
 * promptem nie ginęła wraz z odmową i dała się użyć poza produktem.
 *
 * Kształt `details` jest sprawdzany, nie zakładany: kontrakt opisuje to pole
 * jako `unknown`, więc rzutowanie na własny typ byłoby obietnicą bez pokrycia.
 * Brak pola albo pole innego kształtu daje pusty wynik, nie wyjątek.
 */

/** Gotowa treść polecenia z odmowy; pusty łańcuch znaczy: rdzeń jej nie podał. */
export function polecenieZOdmowy(blad: ErrorInfo | undefined): string {
  const szczegoly = blad?.details;
  if (typeof szczegoly !== 'object' || szczegoly === null) return '';
  const pole = (szczegoly as Record<string, unknown>)['polecenie'];
  return typeof pole === 'string' ? pole.trim() : '';
}

/**
 * Zdanie o tym, czym jest oddany tekst — pokazywane razem z nim, nigdy osobno.
 *
 * Sam tekst polecenia obok przycisku generowania czytałby się jak wynik
 * generowania, a jest wynikiem składania promptu; zdanie mówi wprost, że obrazu
 * nie ma i skąd tekst pochodzi.
 */
export function zdanieOPoleceniu(polecenie: string): string {
  if (polecenie === '') return '';
  return (
    'Obrazu nie ma i nie będzie go w tej turze. Rdzeń złożył za to pola kreatora w gotową ' +
    `treść polecenia: „${polecenie}". To wynik składania promptu, nie generowania — ` +
    'tekst nadaje się do podania silnikowi graficznemu poza tym produktem.'
  );
}
