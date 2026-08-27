import {
  StudioAuthor,
  StudioLockScope,
  type StudioActionBalance,
  type StudioAgentSlot,
  type StudioFragmentLock,
} from '../../../../shared/contract';

/** Zdanie wyświetlane oknu o tym, że blokada obowiązuje w rdzeniu przed dotknięciem treści dokumentu, jedyne miejsce niosące tę treść w interfejsie. */
export const BLOKADA_STOI_W_RDZENIU =
  'Blokada obowiązuje W RDZENIU, przed dotknięciem treści — nie w oknie. Czynność modelu godząca ' +
  'w zablokowany fragment wraca błędem nazywającym fragment i blokadę, a nie cichą bezczynnością.';

/** Zdanie wyświetlane oknu o tym, że blokadę fragmentu zdejmuje wyłącznie operator, a model może jedynie założyć propozycję na marginesie. */
export const BLOKADA_ZDEJMUJE_OPERATOR =
  'Blokadę zdejmuje WYŁĄCZNIE Operator. Model jej nie zdejmuje i nie prosi o zdjęcie obejściem — ' +
  'jeśli uzna, że fragment wymaga zmiany, zakłada propozycję na marginesie.';

/** Nazwa zasięgu blokady widoczna dla operatora, rozróżniająca ustawienie jawne obejmujące oboje od domyślnego skierowanego wyłącznie przeciw modelowi. */
export function blokadaNazwaZasiegu(zasieg: StudioLockScope): string {
  return zasieg === StudioLockScope.Everyone
    ? 'model i Operator — ustawienie jawne, nie domyślne'
    : 'wyłącznie model — Operator zmienia fragment bez przeszkód';
}

/** Zdanie opisujące jedną blokadę fragmentu: jej nazwę, powód założenia, zakres znaków, zasięg działania i pochodzenie z szablonu albo ręczne. */
export function blokadaOpisz(blokada: StudioFragmentLock): string {
  const powod =
    blokada.reason === undefined || blokada.reason === ''
      ? 'powodu Operator nie podał'
      : blokada.reason;
  const zSzablonu =
    blokada.fromTemplate === true
      ? blokada.templateId === undefined || blokada.templateId === ''
        ? ' · blokada wzorcowa z szablonu — KTÓREGO, rdzeń nie podał'
        : ` · blokada wzorcowa z szablonu ${blokada.templateId}`
      : '';
  const kto =
    blokada.createdBy === StudioAuthor.Model
      ? ' · założona przez model'
      : blokada.createdBy === StudioAuthor.Uzytkownik
        ? ' · założona przez Operatora'
        : '';
  const czas =
    blokada.createdAt === undefined
      ? ''
      : ` · ${new Date(blokada.createdAt).toLocaleString('pl-PL')}`;
  return (
    `„${blokada.name}" · znaki ${blokada.rangeStart}–${blokada.rangeEnd} · ${powod} · ` +
    `dotyczy: ${blokadaNazwaZasiegu(blokada.scope)}${zSzablonu}${kto}${czas}`
  );
}

/** Zdanie opisujące wykaz blokad dokumentu, w którym pusty wykaz też jest pełnoprawną odpowiedzią, a nie brakiem danych. */
export function blokadaOpiszWykaz(blokady: readonly StudioFragmentLock[]): string {
  if (blokady.length === 0) {
    return (
      'Dokument nie ma ani jednej blokady. Pusty wykaz znaczy tu „model może zmieniać wszystko", ' +
      'a nie „blokady są wyłączone".'
    );
  }
  const przeciwWszystkim = blokady.filter(
    (blokada) => blokada.scope === StudioLockScope.Everyone,
  ).length;
  const zSzablonu = blokady.filter((blokada) => blokada.fromTemplate === true).length;
  return (
    `Blokad: ${blokady.length}, z tego obejmujących także Operatora: ${przeciwWszystkim}, ` +
    `wzorcowych z szablonu: ${zSzablonu}. Blokady przechodzą przez wersje — przywrócenie ` +
    'wcześniejszej wersji ich nie gubi.'
  );
}

/**
 * Bilans zmiany, która trafiła w blokadę częściowo.
 *
 * `null` znaczy „blokada niczego nie zatrzymała". Zdanie wypisuje blokady po
 * nazwie, bo Operator ma wiedzieć, KTÓRA blokada stanęła na drodze, a nie ile
 * ich było.
 */
export function blokadaBilansPominiec(bilans: StudioActionBalance): string | null {
  const pominiete = (bilans.skipped ?? []).filter(
    (pozycja) =>
      (pozycja.lockId !== undefined && pozycja.lockId !== '') ||
      (pozycja.lockName !== undefined && pozycja.lockName !== ''),
  );
  if (pominiete.length === 0) return null;
  const nazwy = pominiete.map(
    (pozycja) =>
      `${pozycja.lockName === undefined || pozycja.lockName === '' ? pozycja.lockId : `„${pozycja.lockName}"`}` +
      (pozycja.rangeStart === undefined || pozycja.rangeEnd === undefined
        ? ''
        : ` (znaki ${pozycja.rangeStart}–${pozycja.rangeEnd})`),
  );
  return (
    `Zmiana weszła POZA blokadami w ${bilans.applied} miejscach, a w ${pominiete.length} stanęła. ` +
    `Zatrzymały ją: ${nazwy.join(', ')}. Pominięcie nie jest przemilczane — to jest jego bilans.`
  );
}

/** Zdanie opisujące zajęcie fragmentu przez wykonawcę: kto je trzyma, jaki zakres znaków obejmuje, w jakim jest stanie i kiedy zajęcie wygasa. */
export function blokadaOpiszZajecie(zajecie: StudioAgentSlot): string {
  const kto =
    zajecie.actor.agentName ?? zajecie.actor.agentId ?? 'wykonawca nienazwany';
  const zakres =
    zajecie.rangeStart === undefined || zajecie.rangeEnd === undefined
      ? 'bez wskazanego fragmentu'
      : `znaki ${zajecie.rangeStart}–${zajecie.rangeEnd}`;
  const wygasa =
    zajecie.expiresAt === undefined
      ? 'czasu wygaśnięcia rdzeń nie podał'
      : `zajęcie wygasa ${new Date(zajecie.expiresAt).toLocaleString('pl-PL')}`;
  return `${kto} · ${zakres} · stan: ${zajecie.state} · ${wygasa}`;
}

/**
 * Zdanie o odmowie zajęcia — z NAZWANIEM wykonawcy i czasu.
 *
 * Odmowa zajęcia wraca odpowiedzią udaną (`claimed: false`), więc bez tego
 * zdania Operator zobaczyłby powodzenie tam, gdzie fragmentu nie zajęto.
 */
export function blokadaOpiszOdmoweZajecia(
  trzymajacy: StudioAgentSlot | undefined,
  powod: string | undefined,
): string {
  if (trzymajacy === undefined) {
    return (
      'Fragmentu NIE zajęto. ' +
      (powod === undefined || powod === ''
        ? 'Rdzeń nie podał powodu ani wykonawcy trzymającego fragment — to jest brak do zgłoszenia, ' +
          'bo odmowa ma nazywać wykonawcę i czas.'
        : powod)
    );
  }
  const kto = trzymajacy.actor.agentName ?? trzymajacy.actor.agentId ?? 'wykonawca nienazwany';
  const doKiedy =
    trzymajacy.expiresAt === undefined
      ? 'bez podanego czasu wygaśnięcia'
      : `do ${new Date(trzymajacy.expiresAt).toLocaleString('pl-PL')}`;
  return (
    `Fragmentu NIE zajęto: trzyma go ${kto} ${doKiedy}` +
    (trzymajacy.rangeStart === undefined || trzymajacy.rangeEnd === undefined
      ? ''
      : ` na znakach ${trzymajacy.rangeStart}–${trzymajacy.rangeEnd}`) +
    (powod === undefined || powod === '' ? '.' : `. ${powod}`)
  );
}

/** Sprawdza, czy zakres zaznaczenia nadaje się na podstawę blokady albo zajęcia fragmentu, czyli czy koniec jest większy od początku. */
export function blokadaZakresPoprawny(
  zakres: { poczatek: number; koniec: number } | null,
): boolean {
  return zakres !== null && zakres.koniec > zakres.poczatek;
}
