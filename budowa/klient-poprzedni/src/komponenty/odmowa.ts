/**
 * Zdanie opisujące odmowę rdzenia powstaje w jednym miejscu dla całego
 * interfejsu, ponieważ kopie rozjeżdżają się po cichu: poprawka w jednym module
 * zostawiłaby pozostałe mówiące Operatorowi co innego o tej samej odmowie.
 */

/**
 * Kształt powodu odmowy przyjmowany bez rozbierania na części, luźniejszy
 * od kształtu `ErrorInfo` z kontraktu, w którym pola `code` oraz `message`
 * są wymagane; `ErrorInfo` pasuje tu bez rzutowania.
 */
export interface PowodOdmowy {
  code?: string;
  message?: string;
}

/**
 * Zdanie opisujące odmowę rdzenia: nazwa nieznanej komendy wyprzedza kod oraz
 * wiadomość, ponieważ rozstrzyga, czy Operator patrzy na usterkę wykonania,
 * czy na czynność, której rdzeń jeszcze nie umie.
 */
export function opisOdmowy(
  czynnosc: string,
  kod?: string,
  wiadomosc?: string,
  nieznanyTyp?: string,
): string {
  const typ = (nieznanyTyp ?? '').trim();
  if (typ !== '') {
    return `${czynnosc}: rdzeń nie zna komendy ${typ} — komenda jest w kontrakcie, uchwytu jeszcze nie ma`;
  }
  const powod = (wiadomosc ?? '').trim();
  const oznaczenie = (kod ?? '').trim();
  if (powod !== '' && oznaczenie !== '') return `${czynnosc}: ${powod} (${oznaczenie})`;
  if (powod !== '') return `${czynnosc}: ${powod}`;
  if (oznaczenie !== '') return `${czynnosc}: ${oznaczenie}`;
  return `${czynnosc}: rdzeń nie podał powodu`;
}

/**
 * To samo zdanie z powodu wziętego w całości, bez rozbierania na kod i treść;
 * wywołanie po odmowie komendy nie powtarza wtedy pary pól `code` i `message`,
 * a wynik wychodzi identyczny jak z funkcji `opisOdmowy`.
 */
export function opisOdmowyBledu(
  czynnosc: string,
  blad?: PowodOdmowy | null,
  nieznanyTyp?: string,
): string {
  return opisOdmowy(czynnosc, blad?.code, blad?.message, nieznanyTyp);
}
