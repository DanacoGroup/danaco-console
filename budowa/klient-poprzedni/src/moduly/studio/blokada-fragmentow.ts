import {
  StudioAuthor,
  StudioLockScope,
  type StudioActionBalance,
  type StudioAgentSlot,
  type StudioFragmentLock,
} from '../../../../shared/contract';

/**
 * Blokada fragmentu i zajęcie fragmentu — dwie różne rzeczy, oba w rdzeniu.
 *
 * ── Blokada obowiązuje W RDZENIU, nie w oknie ───────────────────────────────
 * Sprawdzenie stoi na drodze każdej komendy zmieniającej dokument, po stronie
 * serwera, PRZED dotknięciem treści. Blokada pilnowana przez okno byłaby
 * pozorna: model woła komendy rdzenia tak samo jak klient, więc ominąłby ją bez
 * wysiłku. Okno pokazuje więc blokady i ich skutki, a nie wykonuje ich.
 *
 * ── Dwa byty, których nie wolno pomieszać ───────────────────────────────────
 *   — **blokada** (`studio.lock.*`) jest trwała i skierowana przeciw MODELOWI;
 *     zdejmuje ją wyłącznie Operator, a model, który uzna zmianę za konieczną,
 *     zakłada propozycję na marginesie i tyle;
 *   — **zajęcie** (`studio.agents.claim`) jest czasowe i skierowane przeciw
 *     DRUGIEMU WYKONAWCY, żeby dwóch agentów nie pisało po tym samym akapicie.
 *     Wygasa samo, a odmowa nazywa wykonawcę i czas.
 *
 * ── Bilans zamiast przemilczenia ────────────────────────────────────────────
 * Zmiana obejmująca blokadę częściowo wykonuje się POZA blokadą i oddaje bilans:
 * co przeszło, co pominięte i przez którą blokadę. Odmowa całości byłaby
 * nieproporcjonalna, ale przemilczenie pominięcia jest zakazane — dlatego zdanie
 * o bilansie nazywa blokady po nazwie, nie po liczbie.
 *
 * Plik nie zna DOM ani rdzenia.
 */

/** Zdanie o tym, gdzie blokada obowiązuje — jedno miejsce tej treści w oknie. */
export const BLOKADA_STOI_W_RDZENIU =
  'Blokada obowiązuje W RDZENIU, przed dotknięciem treści — nie w oknie. Czynność modelu godząca ' +
  'w zablokowany fragment wraca błędem nazywającym fragment i blokadę, a nie cichą bezczynnością.';

/** Zdanie o tym, kto blokadę zdejmuje. */
export const BLOKADA_ZDEJMUJE_OPERATOR =
  'Blokadę zdejmuje WYŁĄCZNIE Operator. Model jej nie zdejmuje i nie prosi o zdjęcie obejściem — ' +
  'jeśli uzna, że fragment wymaga zmiany, zakłada propozycję na marginesie.';

/** Nazwa zasięgu blokady widoczna dla Operatora. */
export function blokadaNazwaZasiegu(zasieg: StudioLockScope): string {
  return zasieg === StudioLockScope.Everyone
    ? 'model i Operator — ustawienie jawne, nie domyślne'
    : 'wyłącznie model — Operator zmienia fragment bez przeszkód';
}

/** Zdanie o jednej blokadzie: nazwa, powód, zakres, zasięg i pochodzenie. */
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

/** Zdanie o wykazie blokad; pusty wykaz też jest odpowiedzią. */
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

/** Zdanie o zajęciu fragmentu przez wykonawcę wraz z czasem wygaśnięcia. */
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

/** Czy zakres zaznaczenia nadaje się na blokadę albo zajęcie. */
export function blokadaZakresPoprawny(
  zakres: { poczatek: number; koniec: number } | null,
): boolean {
  return zakres !== null && zakres.koniec > zakres.poczatek;
}
