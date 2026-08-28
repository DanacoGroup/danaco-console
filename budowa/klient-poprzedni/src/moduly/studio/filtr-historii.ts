import { StudioAuthor, type StudioVersion } from '../../../../shared/contract';

/** Wartość filtru wykazu wersji historii dokumentu, zawężającego wykaz po tym, co wersja niesie: etykietę, opis, kluczowość albo autora. */
export type FiltrHistorii =
  | 'wszystkie'
  | 'etykietowane'
  | 'zopisem'
  | 'kluczowe'
  | 'operator'
  | 'model';

/** Pozycje listy wyboru filtru wykazu wersji dokumentu wraz z ich etykietami widocznymi w kontrolce wyboru. */
export const POZYCJE_FILTRU: readonly { wartosc: FiltrHistorii; etykieta: string }[] = [
  { wartosc: 'wszystkie', etykieta: 'Wszystkie wersje' },
  { wartosc: 'etykietowane', etykieta: 'Tylko wersje z etykietą' },
  { wartosc: 'zopisem', etykieta: 'Tylko wersje z opisem zmiany' },
  { wartosc: 'kluczowe', etykieta: 'Tylko wersje kluczowe' },
  { wartosc: 'operator', etykieta: 'Tylko zmiany Operatora' },
  { wartosc: 'model', etykieta: 'Tylko zmiany modelu' },
];

/**
 * Zdanie o granicy filtru autora.
 *
 * Stoi przy kontrolce, bo granica jest niewidoczna: Operator, który zawęża wykaz
 * do zmian modelu, ma wiedzieć, że wersje bez zapisanego autora nie znikają
 * w milczeniu, a ich autora rdzeń po prostu nie zna.
 */
export const BRAK_FILTRU_AUTORA =
  'Filtr autora czyta pole StudioVersion.author (uzytkownik albo model). Pole jest ' +
  'nieobowiązkowe: wersje założone przed jego wprowadzeniem autora nie niosą i zostają ' +
  'w wykazie z autorem pustym, zamiast wypadać z niego po cichu. Zawężenie po nazwie autora ' +
  'stoi osobnym polem obok.';

/** Zawęża wykaz wersji dokumentu wedle wybranego filtru; wartość wszystkie niczego nie odsiewa i oddaje cały wykaz. */
export function przefiltrujHistorie(
  wersje: readonly StudioVersion[],
  filtr: FiltrHistorii,
): readonly StudioVersion[] {
  if (filtr === 'etykietowane') {
    return wersje.filter((wersja) => (wersja.label ?? '') !== '');
  }
  if (filtr === 'zopisem') {
    return wersje.filter((wersja) => (wersja.summary ?? '') !== '');
  }
  if (filtr === 'kluczowe') {
    return wersje.filter((wersja) => wersja.milestone === true);
  }
  if (filtr === 'operator') {
    return wersje.filter((wersja) => wersja.author !== StudioAuthor.Model);
  }
  if (filtr === 'model') {
    return wersje.filter((wersja) => wersja.author === StudioAuthor.Model);
  }
  return wersje;
}

/** Sprawdza, czy wersja pochodzi z zapisu samoczynnego: nie jest kluczowa, nie ma etykiety własnej i nie jest autorstwa modelu. */
export function czyZapisSamoczynny(wersja: StudioVersion): boolean {
  if (wersja.milestone === true) return false;
  if ((wersja.label ?? '') !== '') return false;
  if (wersja.author === StudioAuthor.Model) return false;
  return (wersja.summary ?? '') === '';
}

/** Sprawdza, czy wersja pasuje do szukanej nazwy autora; zapytanie puste przepuszcza wszystkie wersje bez rozróżniania. */
export function pasujeDoAutora(wersja: StudioVersion, szukane: string): boolean {
  const zapytanie = szukane.trim().toLowerCase();
  if (zapytanie === '') return true;
  if (wersja.author === undefined) return false;
  const nazwa = wersja.author === StudioAuthor.Model ? 'model' : 'operator uzytkownik';
  return nazwa.includes(zapytanie);
}

/** Zdanie opisujące skutek filtru wykazu wersji; pustka po zawężeniu nie oznacza pustej historii dokumentu. */
export function opiszZawezenie(
  wszystkie: number,
  widoczne: number,
  filtr: FiltrHistorii,
): string {
  if (filtr === 'wszystkie') return `Wersji w historii: ${wszystkie}.`;
  if (widoczne === 0) {
    return (
      `Wersji w historii: ${wszystkie}, ale żadna nie spełnia wybranego zawężenia. ` +
      'Historia nie jest pusta — pusty jest wynik filtru.'
    );
  }
  return `Wersji w historii: ${wszystkie}; zawężenie pokazuje ${widoczne}.`;
}
