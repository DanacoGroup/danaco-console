/**
 * Zdanie opisujące odmowę rdzenia — jedno dla całego interfejsu.
 *
 * Zdanie powstaje w jednym miejscu, bo kopie rozjeżdżają się po cichu: poprawka
 * w jednym module zostawiłaby pozostałe mówiące Operatorowi co innego o tej
 * samej odmowie.
 *
 * Kod błędu zostaje przy wiadomości z zamysłem: `validation_failed`
 * i `not_found` znaczą dla Operatora co innego, a treść wiadomości bywa dla obu
 * tym samym zdaniem. Brak wiadomości nie zostawia pustki — idzie wtedy sam kod,
 * a gdy nie ma i kodu, zdanie nazywa milczenie rdzenia wprost, żeby Operator
 * nie patrzył na pusty prostokąt.
 */

/**
 * Kształt powodu odmowy przyjmowany bez rozbierania na części.
 *
 * Luźniejszy od `ErrorInfo` z kontraktu (tam `code` i `message` są wymagane):
 * moduły dostają powód z pola `blad?` wyniku komendy, które bywa puste, a
 * pojedyncze okna trzymają własne, częściowe zapisy odmowy. `ErrorInfo`
 * pasuje tu bez rzutowania.
 */
export interface PowodOdmowy {
  code?: string;
  message?: string;
}

/**
 * Zdanie opisujące odmowę rdzenia.
 *
 * `nieznanyTyp` obsługuje przypadek osobny, otwarty dla każdego modułu: rdzeń
 * odpowiedział zdarzeniem „nie znam tej komendy”.
 * To nie jest awaria wykonania i Operator musi widzieć różnicę — nazwa typu
 * idzie na początek zdania i wyprzedza kod oraz wiadomość, bo rozstrzyga,
 * czy patrzy na usterkę, czy na czynność, której rdzeń jeszcze nie umie.
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
 * To samo zdanie z powodu wziętego w całości, bez rozbierania na kod i treść.
 *
 * Istnieje, żeby wywołanie po odmowie komendy nie powtarzało w każdym miejscu
 * `wynik.blad?.code, wynik.blad?.message` — para pól rozjeżdża się przy
 * przepisywaniu (łatwo podać kod z jednego wyniku, a wiadomość z drugiego).
 * Zdanie wychodzi identyczne jak z `opisOdmowy`.
 */
export function opisOdmowyBledu(
  czynnosc: string,
  blad?: PowodOdmowy | null,
  nieznanyTyp?: string,
): string {
  return opisOdmowy(czynnosc, blad?.code, blad?.message, nieznanyTyp);
}
