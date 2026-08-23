import { StudioAuthor, type StudioVersion } from '../../../../shared/contract';

/**
 * Filtr wykazu wersji historii — zawężenie po tym, co wersja NIESIE.
 *
 * ── Filtr po autorze jest już wykonalny ─────────────────────────────────────
 * `StudioVersion` niesie dziś `author` (`uzytkownik` albo `model`) i `milestone`,
 * więc rozdzielenie zmian Operatora od zmian modelu ma po czym przebiegać —
 * inaczej niż w turze, w której ten plik powstał. Pole jest **nieobowiązkowe**:
 * wersje założone przed jego wprowadzeniem autora nie niosą i takich wersji filtr
 * autora nie odsiewa po cichu, tylko trzyma je w wykazie z autorem pustym. Odsianie
 * ich byłoby ukryciem historii przed Operatorem.
 *
 * ── Autozapis idzie osobnym szeregiem ───────────────────────────────────────
 * Zapis samoczynny poznaje się po tym, że NIE jest wersją kluczową i nie ma
 * etykiety własnej — a taką nadaje wyłącznie Operator. Drugiego pojęcia okno nie
 * zakłada; rozróżnienie stoi na `studio.version.label.set`, jak stanowi zlecenie.
 *
 * Plik nie zna DOM: wejściem są wersje kontraktu, wyjściem wersje przefiltrowane.
 */

/** Wartość filtru wykazu wersji. */
export type FiltrHistorii =
  | 'wszystkie'
  | 'etykietowane'
  | 'zopisem'
  | 'kluczowe'
  | 'operator'
  | 'model';

/** Pozycje listy wyboru filtru wraz z ich znaczeniem. */
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

/** Zawęża wykaz wersji; `wszystkie` niczego nie odsiewa. */
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

/**
 * Czy wersja pochodzi z zapisu samoczynnego.
 *
 * Wersja kluczowa albo z etykietą własną NIE jest samoczynna — oba oznaczenia
 * nadaje Operator. Reszta wersji autora `model` też nie: zmiana modelu jest pracą
 * nad pismem, a nie zapisem w tle, i ukrycie jej byłoby ukryciem tego, co model
 * zrobił.
 */
export function czyZapisSamoczynny(wersja: StudioVersion): boolean {
  if (wersja.milestone === true) return false;
  if ((wersja.label ?? '') !== '') return false;
  if (wersja.author === StudioAuthor.Model) return false;
  return (wersja.summary ?? '') === '';
}

/**
 * Czy wersja pasuje do szukanej nazwy autora.
 *
 * Zapytanie puste przepuszcza wszystko. Wersja bez zapisanego autora przechodzi
 * wyłącznie przy zapytaniu pustym: przy szukaniu „model" nie wolno jej oddać, bo
 * nie wiadomo, czy nim jest, ale i nie wolno o niej zapomnieć — dlatego okno
 * pokazuje liczbę zawężenia obok wykazu.
 */
export function pasujeDoAutora(wersja: StudioVersion, szukane: string): boolean {
  const zapytanie = szukane.trim().toLowerCase();
  if (zapytanie === '') return true;
  if (wersja.author === undefined) return false;
  const nazwa = wersja.author === StudioAuthor.Model ? 'model' : 'operator uzytkownik';
  return nazwa.includes(zapytanie);
}

/** Zdanie o skutku filtru — pustka po zawężeniu to nie pustka historii. */
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
