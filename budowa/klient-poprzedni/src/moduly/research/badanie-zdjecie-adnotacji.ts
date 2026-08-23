import type { ResearchAnnotation } from '../../../../shared/contract';

/**
 * Zapowiedź zdjęcia adnotacji — jedno ostrzeżenie przed czynnością
 * nieodwracalną.
 *
 * `research.annotation.remove` usuwa podświetlenie albo notatkę Operatora
 * i rdzeń nie ma komendy, która by je przywróciła: adnotacja niesie cytat,
 * komentarz i kotwicę pozycji, a po zdjęciu nie ma z czego ich odtworzyć.
 * Dlatego czynność mówi to PRZED wykonaniem i wymaga drugiego naciśnięcia tego
 * samego chwytu.
 *
 * Uzbrojenie jest wiedzą modułu, nie okna: adnotacja wskazana w Reading View
 * i chwyt naciśnięty w panelu akcji tego samego okna dotyczą jednego badania,
 * a moduł prowadzi jedno badanie na sesję (`stan-badania.ts`). Osobne
 * uzbrojenie per okno pozwoliłoby uzbroić w jednym oknie, a zdjąć w drugim —
 * czyli bez ostrzeżenia tam, gdzie Operator naciska.
 *
 * Zapowiedź nie jest okienkiem dialogowym i nie odbiera klikalności: pierwsze
 * naciśnięcie oddaje zdanie odmowy z powodem, drugie wykonuje. Adnotacja
 * wskazana inna niż uzbrojona uzbraja od nowa — Operator nigdy nie zdejmie
 * czegoś, o czym nie został ostrzeżony imiennie.
 */

/** Adnotacja uzbrojona do zdjęcia; puste znaczy „nic nie jest uzbrojone”. */
let uzbrojona = '';

/** Czy ta adnotacja została już zapowiedziana i czeka na drugie naciśnięcie. */
export function czyUzbrojona(idAdnotacji: string): boolean {
  return uzbrojona !== '' && uzbrojona === idAdnotacji;
}

/** Uzbraja adnotację i oddaje zdanie zapowiedzi dla okna. */
export function uzbrojZdjecie(adnotacja: ResearchAnnotation): string {
  uzbrojona = adnotacja.id;
  return (
    'ZDJĘCIE ADNOTACJI JEST NIEODWRACALNE — rdzeń nie ma komendy, która by ją przywróciła. ' +
    `Do zdjęcia stoi ${opisAdnotacji(adnotacja)}. Naciśnij „Zdejmij adnotację" po raz drugi, ` +
    'żeby ją zdjąć; wskazanie innej adnotacji zaczyna ostrzeżenie od nowa.'
  );
}

/** Rozbraja po wykonaniu albo po zmianie wskazania. */
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

/** Najnowsza adnotacja wykazu; `null`, gdy wykaz jest pusty. */
export function najnowszaAdnotacja(
  adnotacje: readonly ResearchAnnotation[],
): ResearchAnnotation | null {
  let najnowsza: ResearchAnnotation | null = null;
  for (const adnotacja of adnotacje) {
    if (najnowsza === null || adnotacja.createdAt > najnowsza.createdAt) najnowsza = adnotacja;
  }
  return najnowsza;
}

/** Skrót długiego napisu; granica jest w oknie, więc nie obcina po cichu. */
function skrot(tekst: string): string {
  return tekst.length <= 120 ? tekst : `${tekst.slice(0, 120)}… (${String(tekst.length)} znaków)`;
}
