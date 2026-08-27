import type { ErrorInfo } from '../../../../shared/contract';

/**
 * Odczyt treści polecenia ocalałej w odmowie `design.asset.generate`: funkcja
 * wyjmuje gotową treść z pola `details.polecenie`, a pusty łańcuch znaczy, że
 * rdzeń treści nie podał albo że pole ma inny kształt.
 */
export function polecenieZOdmowy(blad: ErrorInfo | undefined): string {
  const szczegoly = blad?.details;
  if (typeof szczegoly !== 'object' || szczegoly === null) return '';
  const pole = (szczegoly as Record<string, unknown>)['polecenie'];
  return typeof pole === 'string' ? pole.trim() : '';
}

/**
 * Zdanie o tym, czym jest oddany tekst; pokazywane razem z nim, nigdy osobno.
 * Tekst polecenia obok przycisku generowania czytałby się jak wynik
 * generowania, a jest wynikiem składania promptu.
 */
export function zdanieOPoleceniu(polecenie: string): string {
  if (polecenie === '') return '';
  return (
    'Obrazu nie ma i nie będzie go w tej turze. Rdzeń złożył za to pola kreatora w gotową ' +
    `treść polecenia: „${polecenie}". To wynik składania promptu, nie generowania — ` +
    'tekst nadaje się do podania silnikowi graficznemu poza tym produktem.'
  );
}
