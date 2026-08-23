import {
  StudioActionKind,
  StudioActionState,
  StudioAuthor,
  type StudioActionBalance,
  type StudioDocumentAction,
  type StudioFormDiffEntry,
  type StudioModelChangeSummary,
} from '../../../../shared/contract';

/**
 * Odwracalny dziennik czynności dokumentu — słowa, którymi okno mówi
 * Operatorowi, co cofa i czego cofnąć nie może.
 *
 * ── Trzy rzeczy, których nie wolno zgubić ───────────────────────────────────
 * 1. **Cofnięcie nakłada RÓŻNICĘ drzew, nie migawkę.** Praca naniesiona po
 *    cofanej czynności zostaje. Operator, który tego nie wie, nie odważy się
 *    cofnąć czynności ze środka dziennika — a to jest cała wartość tej rodziny.
 *    Zdanie stoi więc przy każdej czynności, nie w pomocy.
 * 2. **Czynność będąca podstawą późniejszej ODMAWIA i nazywa te czynności.**
 *    Rdzeń oddaje je w `blockedBy` wpisu oraz w bilansie odpowiedzi. Okno
 *    pokazuje powód, a nie samo „nie udało się" — inaczej Operator zostaje
 *    z zamkniętą drogą i bez wskazania, co ją zamknęło.
 * 3. **Wykonawca jest nazwany.** Wpis niesie `authorAgentName`, więc dwóch
 *    agentów nie pokazuje się jako jeden „model".
 *
 * Plik nie zna DOM ani rdzenia: wejściem są byty kontraktu, wyjściem napisy
 * i rozstrzygnięcia. Dzięki temu sprawdza się bez stawiania okna.
 */

/** Zdanie o różnicy drzew — jedno miejsce tej treści w całym oknie. */
export const DZIENNIK_ROZNICA_DRZEW =
  'Cofnięcie nakłada RÓŻNICĘ drzew dokumentu, a nie migawkę stanu sprzed: praca naniesiona po ' +
  'tej czynności ZOSTAJE. Cofasz tę jedną czynność, nie wszystko, co po niej.';

/** Nazwa rodzaju czynności widoczna dla Operatora — pełna, bez kodu. */
export const NAZWY_RODZAJOW_CZYNNOSCI: Readonly<Record<StudioActionKind, string>> = {
  [StudioActionKind.TextEdit]: 'zmiana treści',
  [StudioActionKind.FormatChange]: 'zmiana postaci',
  [StudioActionKind.StyleChange]: 'zmiana stylu nazwanego',
  [StudioActionKind.PageChange]: 'zmiana nastaw strony albo sekcji',
  [StudioActionKind.ListChange]: 'zmiana listy',
  [StudioActionKind.TableChange]: 'zmiana tabeli',
  [StudioActionKind.ObjectChange]: 'zmiana obiektu osadzonego',
  [StudioActionKind.ApparatusChange]: 'zmiana aparatu dokumentu',
  [StudioActionKind.MarkupChange]: 'zmiana znakowania',
  [StudioActionKind.ProposalAccept]: 'przyjęcie propozycji modelu',
  [StudioActionKind.ImportChange]: 'wniesienie treści z pliku albo źródła',
  [StudioActionKind.FieldChange]: 'zmiana pola dokumentu',
};

/** Nazwa rodzaju czynności; rodzaj nieznany oddaje swoją wartość, nie pustkę. */
export function dziennikNazwaRodzaju(rodzaj: StudioActionKind | string): string {
  return NAZWY_RODZAJOW_CZYNNOSCI[rodzaj as StudioActionKind] ?? rodzaj;
}

/** Kto czynność wykonał — wykonawca nazwany, gdy rdzeń podał jego nazwę. */
export function dziennikNazwaWykonawcy(wpis: StudioDocumentAction): string {
  if (wpis.author !== StudioAuthor.Model) return 'Operator';
  const nazwa = wpis.authorAgentName ?? '';
  if (nazwa !== '') {
    return wpis.authorAgentVersion === undefined || wpis.authorAgentVersion === ''
      ? `model — ${nazwa}`
      : `model — ${nazwa} (wersja ${wpis.authorAgentVersion})`;
  }
  // Rdzeń nie stempluje tożsamości wykonawcy na wszystkich drogach postaci.
  // Podstawienie tu nazwy własnej byłoby wskazaniem wykonawcy, którego rdzeń
  // nie zapisał — brak jest nazwany, a nie zasłonięty.
  return 'model — wykonawca nienazwany';
}

/** Zdanie o wpisie dziennika: co, kto, kiedy i czego dotknęło. */
export function dziennikOpiszWpis(wpis: StudioDocumentAction): string {
  const gdzie =
    wpis.rangeStart === undefined || wpis.rangeEnd === undefined
      ? 'cały dokument'
      : `znaki ${wpis.rangeStart}–${wpis.rangeEnd}`;
  const stan = wpis.state === StudioActionState.Reverted ? 'cofnięta' : 'stoi w dokumencie';
  return (
    `${dziennikNazwaRodzaju(wpis.kind)} · ${dziennikNazwaWykonawcy(wpis)} · ${gdzie} · ` +
    `${stan} · ${new Date(wpis.createdAt).toLocaleString('pl-PL')}`
  );
}

/**
 * Czy tę czynność da się cofnąć samodzielnie.
 *
 * Czynność już cofnięta nie ma czego cofać — dla niej właściwe jest ponowienie.
 * Czynność, na której stoją późniejsze, cofnąć się nie da i okno mówi to
 * ZAWCZASU, a nie po odmowie: przycisk pozostaje jednak czynny, bo prawdę
 * rozstrzyga rdzeń, a wykaz `blockedBy` mógł się zmienić po ostatnim odczycie.
 */
export function dziennikCzyCofnieciePowstrzymane(wpis: StudioDocumentAction): boolean {
  return (wpis.blockedBy ?? []).length > 0;
}

/** Zdanie o zależnościach wpisu; brak zależności też jest odpowiedzią. */
export function dziennikOpiszZaleznosci(wpis: StudioDocumentAction): string {
  const stoiNa = wpis.dependsOn ?? [];
  const stojaNaNiej = wpis.blockedBy ?? [];
  if (stoiNa.length === 0 && stojaNaNiej.length === 0) {
    return 'Ta czynność nie ma zależności: nic na niej nie stoi i sama nie stoi na niczym.';
  }
  const czesci: string[] = [];
  if (stojaNaNiej.length > 0) {
    czesci.push(
      `Na tej czynności stoją późniejsze (${stojaNaNiej.length}): ${stojaNaNiej.join(', ')}. ` +
        'Cofnięcie jej samej rdzeń ODMÓWI i nazwie te czynności — najpierw cofnij je.',
    );
  }
  if (stoiNa.length > 0) {
    czesci.push(`Ta czynność stoi na wcześniejszych (${stoiNa.length}): ${stoiNa.join(', ')}.`);
  }
  return czesci.join(' ');
}

/**
 * Zdanie o bilansie czynności — co przeszło, co stanęło i przez co.
 *
 * Bilans jest tu jedyną drogą, którą Operator dowiaduje się o pominięciu.
 * Przemilczenie pominięcia jest w tym module zakazane wprost, więc pominięcia
 * są wypisane po jednym, wraz z powodem podanym przez rdzeń.
 */
export function dziennikOpiszBilans(bilans: StudioActionBalance): string {
  const czesci: string[] = [
    `Weszło w ${bilans.applied} miejscach, stanęło w ${bilans.skippedCount}.`,
  ];
  if (bilans.note !== undefined && bilans.note !== '') czesci.push(bilans.note);
  for (const pominiete of bilans.skipped ?? []) {
    czesci.push(`Pominięte: ${dziennikOpiszPominiecie(pominiete)}`);
  }
  if (bilans.deferredCount !== undefined && bilans.deferredCount > 0) {
    czesci.push(
      `Zmian odłożonych, żeby nie nadpisać cudzej pracy: ${bilans.deferredCount}. ` +
        'Odłożone brzmienie NIE przepadło — stoi w wykazie spięć wykonawców.',
    );
  }
  for (const spiecie of bilans.conflicts ?? []) {
    czesci.push(
      `Spięcie o znaki ${spiecie.rangeStart}–${spiecie.rangeEnd}: ${spiecie.reason}`,
    );
  }
  return czesci.join(' ');
}

/** Zdanie o jednym pominięciu; kształt pozycji bilansu bywa niepełny. */
function dziennikOpiszPominiecie(pominiete: {
  rangeStart?: number;
  rangeEnd?: number;
  reason?: string;
  detail?: string;
  lockId?: string;
  lockName?: string;
}): string {
  const gdzie =
    pominiete.rangeStart === undefined || pominiete.rangeEnd === undefined
      ? 'miejsce niewskazane'
      : `znaki ${pominiete.rangeStart}–${pominiete.rangeEnd}`;
  const przezCo =
    pominiete.lockName !== undefined && pominiete.lockName !== ''
      ? ` — zatrzymała je blokada „${pominiete.lockName}"`
      : pominiete.lockId !== undefined && pominiete.lockId !== ''
        ? ` — zatrzymała je blokada ${pominiete.lockId}`
        : '';
  const powod =
    pominiete.reason === undefined || pominiete.reason === '' ? '' : ` (${pominiete.reason})`;
  const czego =
    pominiete.detail === undefined || pominiete.detail === '' ? '' : ` — ${pominiete.detail}`;
  return `${gdzie}${przezCo}${czego}${powod}`;
}

/**
 * Zdanie o odmowie cofnięcia — powód, nie „nie udało się".
 *
 * Wywołanie cofnięcia bywa UDANE, a mimo to nic nie cofa: rdzeń oddaje wykaz
 * cofniętych pusty i bilans z pominięciem nazywającym zależność. Okno musi to
 * odróżnić od powodzenia, bo inaczej pokazałoby „cofnięto" po czynności, która
 * nie zeszła.
 */
export function dziennikOpiszCofniecie(
  cofniete: readonly string[],
  bilans: StudioActionBalance,
): { udane: boolean; zdanie: string } {
  if (cofniete.length === 0) {
    return {
      udane: false,
      zdanie:
        'Rdzeń nie cofnął ani jednej czynności i podał powód: ' +
        `${dziennikOpiszBilans(bilans)} ${DZIENNIK_ROZNICA_DRZEW}`,
    };
  }
  return {
    udane: true,
    zdanie:
      `Cofnięto czynności: ${cofniete.length}. ${dziennikOpiszBilans(bilans)} ` +
      DZIENNIK_ROZNICA_DRZEW,
  };
}

/** Zestawienie zmian modelu jako zdanie przy przełączniku podświetlenia. */
export function dziennikOpiszZmianyModelu(zestawienie: StudioModelChangeSummary): string {
  if (zestawienie.total === 0) {
    return 'Model nie wniósł do tego dokumentu ani jednej zmiany — ani treści, ani postaci.';
  }
  const czesci = [
    `Zmian modelu w dokumencie: ${zestawienie.total} ` +
      `(treści ${zestawienie.contentChanges}, postaci ${zestawienie.formChanges}); ` +
      `czeka na rozstrzygnięcie ${zestawienie.openChanges}.`,
  ];
  const rozbicie = zestawienie.byAgent ?? [];
  if (rozbicie.length > 0) {
    czesci.push(
      'Wedle wykonawcy: ' +
        rozbicie
          .map(
            (pozycja) =>
              `${pozycja.actor.agentName ?? pozycja.actor.agentId ?? 'wykonawca nienazwany'} — ` +
              `${pozycja.total}`,
          )
          .join('; ') +
        '.',
    );
  }
  return czesci.join(' ');
}

/**
 * Zdanie o cofnięciu zmian modelu.
 *
 * Liczba zmian Operatora ZACHOWANYCH jest tu treścią równie ważną jak liczba
 * cofniętych: to ona odróżnia cofnięcie pracy modelu od przywrócenia wersji
 * sprzed, które skasowałoby też pracę Operatora.
 */
export function dziennikOpiszCofniecieModelu(
  cofnietych: number,
  zachowanychOperatora: number,
  bilans: StudioActionBalance,
  idKopii: string | undefined,
): string {
  const kopia =
    idKopii === undefined || idKopii === ''
      ? 'Kopii zapasowej rdzeń nie zgłosił — to jest brak do sprawdzenia, bo cofnięcie miało ją założyć.'
      : `Kopia zapasowa założona PRZED cofnięciem: ${idKopii}.`;
  return (
    `Cofnięto zmian modelu: ${cofnietych}. Zmian Operatora naniesionych w tym czasie ` +
    `ZACHOWANO: ${zachowanychOperatora} — to nie jest przywrócenie wersji sprzed. ` +
    `${dziennikOpiszBilans(bilans)} ${kopia}`
  );
}

/** Zdanie o różnicy postaci; zero cech też jest odpowiedzią, nie pustką. */
export function dziennikOpiszRoznicePostaci(
  wpisy: readonly StudioFormDiffEntry[],
  liczby: { added: number; removed: number; changed: number },
): string {
  if (wpisy.length === 0) {
    return (
      'Postać obu wersji jest ta sama: rdzeń nie znalazł ani jednej różnicy cechy. ' +
      'To odpowiedź, nie brak wyniku — różnica treści jest liczona osobno.'
    );
  }
  return (
    `Różnic postaci: ${wpisy.length} — cech doszło ${liczby.added}, odpadło ${liczby.removed}, ` +
    `zmieniło się ${liczby.changed}.`
  );
}

/** Zdanie o jednej różnicy postaci wraz z obszarem i stanem przed i po. */
export function dziennikOpiszRoznicePostaciWpis(wpis: StudioFormDiffEntry): string {
  const gdzie =
    wpis.rangeStart === undefined || wpis.rangeEnd === undefined
      ? 'cały dokument'
      : `znaki ${wpis.rangeStart}–${wpis.rangeEnd}`;
  const przed = wpis.before === undefined || wpis.before === '' ? 'brak' : wpis.before;
  const po = wpis.after === undefined || wpis.after === '' ? 'brak' : wpis.after;
  return `${wpis.area} · ${wpis.kind} · ${gdzie} · ${wpis.detail} · przed: ${przed} → po: ${po}`;
}

/** Wpisy dziennika od najświeższego; rdzeń tak je oddaje, okno tego nie odwraca. */
export function dziennikPoKolejnosci(
  wpisy: readonly StudioDocumentAction[],
): readonly StudioDocumentAction[] {
  return [...wpisy].sort((pierwszy, drugi) => drugi.sequence - pierwszy.sequence);
}
