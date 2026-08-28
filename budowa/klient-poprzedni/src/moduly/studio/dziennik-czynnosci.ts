import {
  StudioActionKind,
  StudioActionState,
  StudioAuthor,
  type StudioActionBalance,
  type StudioDocumentAction,
  type StudioFormDiffEntry,
  type StudioModelChangeSummary,
} from '../../../../shared/contract';

/** Odwracalny dziennik czynności dokumentu opisuje słowami, co okno cofa i czego cofnąć nie może, przy czym cofnięcie nakłada różnicę drzew dokumentu, a nie migawkę stanu sprzed. */
export const DZIENNIK_ROZNICA_DRZEW =
  'Cofnięcie nakłada RÓŻNICĘ drzew dokumentu, a nie migawkę stanu sprzed: praca naniesiona po ' +
  'tej czynności ZOSTAJE. Cofasz tę jedną czynność, nie wszystko, co po niej.';

/** Nazwa rodzaju czynności widoczna dla operatora, pełna i bez kodu wewnętrznego, dla każdego z dwunastu rodzajów niesionych przez kontrakt. */
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

/** Nazwa rodzaju czynności odczytana ze słownika nazw; rodzaj nieznany słownikowi oddaje swoją surową wartość, nie pustkę. */
export function dziennikNazwaRodzaju(rodzaj: StudioActionKind | string): string {
  return NAZWY_RODZAJOW_CZYNNOSCI[rodzaj as StudioActionKind] ?? rodzaj;
}

/** Zdanie mówiące, kto wykonał czynność: operator albo wykonawca modelowy nazwany, gdy rdzeń podał jego imię i wersję. */
export function dziennikNazwaWykonawcy(wpis: StudioDocumentAction): string {
  if (wpis.author !== StudioAuthor.Model) return 'Operator';
  const nazwa = wpis.authorAgentName ?? '';
  if (nazwa !== '') {
    return wpis.authorAgentVersion === undefined || wpis.authorAgentVersion === ''
      ? `model — ${nazwa}`
      : `model — ${nazwa} (wersja ${wpis.authorAgentVersion})`;
  }
  // Rdzeń nie stempluje tożsamości wykonawcy na wszystkich drogach; brak jest nazwany, nie zasłonięty.
  return 'model — wykonawca nienazwany';
}

/** Zdanie opisujące wpis dziennika czynności: co się zmieniło, kto to zrobił, kiedy i którego miejsca dokumentu dotknęło. */
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

/** Sprawdza, czy tę czynność da się cofnąć samodzielnie, biorąc pod uwagę wykaz późniejszych czynności, które na niej stoją. */
export function dziennikCzyCofnieciePowstrzymane(wpis: StudioDocumentAction): boolean {
  return (wpis.blockedBy ?? []).length > 0;
}

/** Zdanie opisujące zależności wpisu dziennika wobec innych czynności; brak zależności też jest pełnoprawną odpowiedzią. */
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

/** Zdanie opisujące bilans czynności: ile weszło, ile stanęło i przez co, z pominięciami wypisanymi po jednym wraz z powodem. */
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

/** Zdanie opisujące jedno pominięcie z bilansu czynności; kształt pozycji bilansu bywa niepełny i pola bywają nieobecne. */
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

/** Zdanie opisujące odmowę cofnięcia czynności nazywające jej powód, zamiast ogólnikowego stwierdzenia, że się nie udało. */
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

/** Zdanie zestawiające zmiany modelu w dokumencie, wyświetlane przy przełączniku podświetlenia tych zmian. */
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

/** Zdanie opisujące cofnięcie zmian modelu wraz z liczbą zmian operatora zachowanych, bo to ona odróżnia je od przywrócenia wersji. */
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

/** Zdanie opisujące różnicę postaci między wersjami dokumentu; brak różnic też jest pełnoprawną odpowiedzią, nie pustką. */
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

/** Zdanie opisujące jedną różnicę postaci wraz z obszarem, rodzajem zmiany oraz stanem przed i po niej. */
export function dziennikOpiszRoznicePostaciWpis(wpis: StudioFormDiffEntry): string {
  const gdzie =
    wpis.rangeStart === undefined || wpis.rangeEnd === undefined
      ? 'cały dokument'
      : `znaki ${wpis.rangeStart}–${wpis.rangeEnd}`;
  const przed = wpis.before === undefined || wpis.before === '' ? 'brak' : wpis.before;
  const po = wpis.after === undefined || wpis.after === '' ? 'brak' : wpis.after;
  return `${wpis.area} · ${wpis.kind} · ${gdzie} · ${wpis.detail} · przed: ${przed} → po: ${po}`;
}

/** Sortuje wpisy dziennika czynności od najświeższego według kolejności nadanej przez rdzeń, bez jej odwracania. */
export function dziennikPoKolejnosci(
  wpisy: readonly StudioDocumentAction[],
): readonly StudioDocumentAction[] {
  return [...wpisy].sort((pierwszy, drugi) => drugi.sequence - pierwszy.sequence);
}
