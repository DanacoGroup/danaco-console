/**
 * Zapowiedź zdjęcia adnotacji: jedno ostrzeżenie przed czynnością nieodwracalną.
 * Zdjęcie usuwa cytat, komentarz i kotwicę pozycji, a rdzeń nie ma komendy
 * przywracającej, więc czynność wymaga drugiego naciśnięcia tego samego chwytu.
 */
import type { ResearchAnnotation } from '../../../../shared/contract';

/**
 * Adnotacja uzbrojona do zdjęcia; puste znaczy, że nic nie jest uzbrojone.
 * Wartość jest jedna na moduł, ponieważ moduł prowadzi jedno badanie na sesję.
 */
let uzbrojona = '';

/**
 * Czy ta adnotacja została już zapowiedziana i czeka na drugie naciśnięcie.
 * Wskazanie innej adnotacji odpowiedzi twierdzącej nie da, więc ostrzeżenie
 * zaczyna się od nowa.
 */
export function czyUzbrojona(idAdnotacji: string): boolean {
  return uzbrojona !== '' && uzbrojona === idAdnotacji;
}

/**
 * Uzbraja adnotację i oddaje zdanie zapowiedzi dla okna. Zdanie nazywa adnotację
 * rodzajem, cytatem i czasem, żeby ostrzeżenie było imienne.
 */
export function uzbrojZdjecie(adnotacja: ResearchAnnotation): string {
  uzbrojona = adnotacja.id;
  return (
    'ZDJĘCIE ADNOTACJI JEST NIEODWRACALNE — rdzeń nie ma komendy, która by ją przywróciła. ' +
    `Do zdjęcia stoi ${opisAdnotacji(adnotacja)}. Naciśnij „Zdejmij adnotację" po raz drugi, ` +
    'żeby ją zdjąć; wskazanie innej adnotacji zaczyna ostrzeżenie od nowa.'
  );
}

/**
 * Rozbraja po wykonaniu albo po zmianie wskazania. Czyści jedyną wartość
 * uzbrojenia, więc kolejne naciśnięcie chwytu znowu zapowiada, zamiast zdejmować.
 */
export function rozbrojZdjecie(): void {
  uzbrojona = '';
}

/**
 * Adnotacja opisana tym, co Operator rozpozna: rodzajem, cytatem i czasem.
 *
 * Sam identyfikator nie mówi, co zniknie — a to jest jedyna rzecz, którą przed
 * czynnością nieodwracalną trzeba wiedzieć.
 */
export function opisAdnotacji(adnotacja: ResearchAnnotation): string {
  const cytat = (adnotacja.quote ?? '').trim();
  const komentarz = (adnotacja.comment ?? '').trim();
  const czesci = [`${adnotacja.kind} z ${new Date(adnotacja.createdAt).toLocaleString('pl-PL')}`];
  czesci.push(cytat === '' ? 'bez cytatu' : `cytat „${skrot(cytat)}"`);
  if (komentarz !== '') czesci.push(`komentarz „${skrot(komentarz)}"`);
  if (adnotacja.findingId !== undefined) {
    czesci.push('adnotacja ma ustalenie powstałe z niej — ustalenie zostaje, znika kotwica do fragmentu');
  }
  return czesci.join(', ');
}

/**
 * Najnowsza adnotacja wykazu; `null`, gdy wykaz jest pusty. Porównuje czas
 * utworzenia, więc kolejność wierszy oddana przez rdzeń nie ma znaczenia.
 */
export function najnowszaAdnotacja(
  adnotacje: readonly ResearchAnnotation[],
): ResearchAnnotation | null {
  let najnowsza: ResearchAnnotation | null = null;
  for (const adnotacja of adnotacje) {
    if (najnowsza === null || adnotacja.createdAt > najnowsza.createdAt) najnowsza = adnotacja;
  }
  return najnowsza;
}

/**
 * Skrót długiego napisu; granica jest w oknie, więc nie obcina po cichu: dopisuje
 * liczbę znaków napisu pełnego.
 */
function skrot(tekst: string): string {
  return tekst.length <= 120 ? tekst : `${tekst.slice(0, 120)}… (${String(tekst.length)} znaków)`;
}
