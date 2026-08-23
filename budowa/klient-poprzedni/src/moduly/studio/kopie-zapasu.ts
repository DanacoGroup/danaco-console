import {
  StudioBackupReason,
  StudioVersionSeries,
  type StudioAutosaveSettings,
  type StudioDocumentBackup,
  type StudioVersion,
} from '../../../../shared/contract';

/**
 * Autozapis, kopie zapasowe i szeregi wersji — słowa, którymi okno mówi
 * Operatorowi, czy jego praca jest bezpieczna.
 *
 * ── Uczciwość zapisu jest tu obowiązkowa ────────────────────────────────────
 * Wskaźnik „zapisano" pokazany, gdy zapis się nie udał, jest najgorszym możliwym
 * błędem tego modułu: Operator zamknie okno i straci pracę. Dlatego stan zapisu
 * ma tu **cztery** wartości, nie dwie — i „nieudany" jest jedną z nich, wraz
 * z powodem podanym przez rdzeń i wskazaniem, że treść leży w kopii zapasowej.
 *
 * ── Autozapis idzie OSOBNYM szeregiem ──────────────────────────────────────
 * Rozróżnienie szeregów nie jest tu wymyślone przez okno: rdzeń oddaje je
 * w `studio.version.series.list` wraz z liczbą wersji w każdym z nich. Okno
 * czyta to pole, a nie odgaduje autozapisu po braku etykiety.
 *
 * ── Kopia zapasowa jest zakładana PRZED zapisem ────────────────────────────
 * Kopia po zapisie nie chroni od niczego: awaria zapisu zostawiłaby dokument
 * uszkodzony bez kopii sprzed. Zdanie o tym stoi przy czynności, bo jest to
 * powód, dla którego kopie zapasowe są czymś innym niż historia wersji.
 *
 * Plik nie zna DOM ani rdzenia: wejściem są byty kontraktu, wyjściem napisy
 * i rozstrzygnięcia.
 */

/** Stan zapisu widoczny dla Operatora — cztery wartości, nie dwie. */
export type StanZapisu = 'zapisano' | 'zapisywanie' | 'niezapisane' | 'nieudany';

/** Nazwa powodu założenia kopii — pełna, widoczna dla Operatora. */
export const NAZWY_POWODOW_KOPII: Readonly<Record<StudioBackupReason, string>> = {
  [StudioBackupReason.Interval]: 'odstęp autozapisu',
  [StudioBackupReason.Event]: 'zdarzenie okna — odejście, zamknięcie, przełączenie',
  [StudioBackupReason.BeforeIrreversible]: 'przed czynnością nieodwracalną',
  [StudioBackupReason.Manual]: 'polecenie Operatora',
};

/** Zdanie o tym, dlaczego kopia idzie przed zapisem, a nie po nim. */
export const KOPIE_PRZED_ZAPISEM =
  'Kopia zapasowa zakładana jest PRZED zapisem, nie po nim: kopia po zapisie nie chroni od ' +
  'niczego, bo awaria zapisu zostawiłaby dokument uszkodzony bez stanu sprzed. Kopie idą też ' +
  'niezależnie od historii wersji, żeby przetrwały awarię procesu.';

/** Nazwa powodu kopii; powód nieznany oddaje swoją wartość, nie pustkę. */
export function kopieNazwaPowodu(powod: StudioBackupReason | string): string {
  return NAZWY_POWODOW_KOPII[powod as StudioBackupReason] ?? powod;
}

/**
 * Stan zapisu wyczytany z nastaw autozapisu.
 *
 * `zmianyNiezapisane` podaje okno, bo to ono wie, czy Operator pisał od
 * ostatniego zapisu — rdzeń tego nie widzi. Niepowodzenie ostatniego zapisu
 * wygrywa nad wszystkim: dopóki go nie naprawiono, „zapisano" jest fałszem.
 */
export function kopieStanZapisu(
  nastawy: StudioAutosaveSettings | null,
  zmianyNiezapisane: boolean,
  zapisWToku: boolean,
): StanZapisu {
  if (nastawy !== null && nastawy.lastSaveFailed === true) return 'nieudany';
  if (zapisWToku) return 'zapisywanie';
  if (zmianyNiezapisane) return 'niezapisane';
  return 'zapisano';
}

/** Zdanie o stanie zapisu wraz z czasem ostatniego udanego zapisu. */
export function kopieOpiszStanZapisu(
  stan: StanZapisu,
  nastawy: StudioAutosaveSettings | null,
): string {
  const czas =
    nastawy === null || nastawy.lastSaveAt === undefined
      ? 'rdzeń nie zgłosił jeszcze żadnego udanego zapisu'
      : `ostatni udany zapis: ${new Date(nastawy.lastSaveAt).toLocaleString('pl-PL')}`;
  if (stan === 'nieudany') {
    const powod =
      nastawy?.lastFailureReason === undefined || nastawy.lastFailureReason === ''
        ? 'rdzeń nie podał powodu'
        : nastawy.lastFailureReason;
    return (
      `ZAPIS SIĘ NIE UDAŁ — ${powod}. Treść leży w kopii zapasowej i da się ją przywrócić; ` +
      `nie zamykaj okna, dopóki zapis nie przejdzie. Wcześniej: ${czas}.`
    );
  }
  if (stan === 'zapisywanie') return `zapisywanie w toku · ${czas}`;
  if (stan === 'niezapisane') {
    return `niezapisane zmiany · ${czas}`;
  }
  return `zapisano · ${czas}`;
}

/** Zdanie o nastawach autozapisu — co dokładnie stoi, a nie „włączony". */
export function kopieOpiszNastawy(nastawy: StudioAutosaveSettings): string {
  if (!nastawy.enabled) {
    return (
      'Autozapis STOI WYŁĄCZONY. Zapis samoczynny jest jawnym, odwracalnym ustawieniem ' +
      'Operatora, więc okno go nie włącza samo.'
    );
  }
  const zdarzenia = [
    nastawy.onBlur === true ? 'odejście od okna' : '',
    nastawy.onClose === true ? 'zamknięcie dokumentu' : '',
    nastawy.onSwitch === true ? 'przełączenie dokumentu' : '',
  ].filter((pozycja) => pozycja !== '');
  const odstep =
    nastawy.intervalSeconds === undefined || nastawy.intervalSeconds <= 0
      ? 'bez odstępu — wyłącznie przy zdarzeniach'
      : `co ${nastawy.intervalSeconds} sekund`;
  const wygasanie = kopieOpiszWygasanie(nastawy);
  return (
    `Autozapis czynny ${odstep}. Zapis przy zdarzeniach: ` +
    `${zdarzenia.length === 0 ? 'przy żadnym' : zdarzenia.join(', ')}. ${wygasanie} ` +
    'Zapis samoczynny obejmuje treść I postać, i idzie osobnym szeregiem wersji.'
  );
}

/** Zdanie o zasadzie wygasania kopii — jawnym ustawieniu, nie zaszytej liczbie. */
export function kopieOpiszWygasanie(nastawy: StudioAutosaveSettings): string {
  const ile =
    nastawy.backupRetentionCount === undefined || nastawy.backupRetentionCount <= 0
      ? ''
      : `zachowywanych kopii: ${nastawy.backupRetentionCount}`;
  const godziny =
    nastawy.backupRetentionHours === undefined || nastawy.backupRetentionHours <= 0
      ? ''
      : `kopia wygasa po ${nastawy.backupRetentionHours} godzinach`;
  const czesci = [ile, godziny].filter((pozycja) => pozycja !== '');
  if (czesci.length === 0) {
    return 'Zasady wygasania kopii rdzeń nie zgłosił — kopie zostają, dopóki Operator jej nie ustawi.';
  }
  return `Wygasanie kopii: ${czesci.join(', ')}.`;
}

/** Zdanie o jednej kopii zapasowej wraz z czasem, rozmiarem i powodzeniem. */
export function kopieOpiszKopie(kopia: StudioDocumentBackup): string {
  const rozmiar =
    kopia.bytes === undefined ? 'rozmiaru rdzeń nie podał' : `${kopia.bytes} bajtów`;
  const niezapisane =
    kopia.unsavedChanges === true
      ? ' · NIESIE ZMIANY NIEZAPISANE — to po niej wraca się po nagłym zamknięciu'
      : '';
  const niepowodzenie = kopia.succeeded
    ? ''
    : ` · ZAPIS KOPII SIĘ NIE UDAŁ: ${kopia.failureReason ?? 'rdzeń nie podał powodu'}`;
  return (
    `${kopieNazwaPowodu(kopia.reason)} · ${rozmiar} · ` +
    `${new Date(kopia.createdAt).toLocaleString('pl-PL')}${niezapisane}${niepowodzenie}`
  );
}

/**
 * Zgłoszenie po nagłym zamknięciu — Studio mówi pierwsze, nie czeka na domysł.
 *
 * `null` znaczy „nie ma czego zgłaszać". Zgłoszenie dotyczy kopii najświeższej
 * niosącej zmiany niezapisane: pytanie o pięć kopii naraz byłoby pytaniem,
 * na które Operator nie ma jak odpowiedzieć.
 */
export function kopieZgloszeniePoZamknieciu(
  kopie: readonly StudioDocumentBackup[],
): { kopia: StudioDocumentBackup; zdanie: string } | null {
  const niezapisane = kopie
    .filter((kopia) => kopia.unsavedChanges === true && kopia.succeeded)
    .sort((pierwsza, druga) => druga.createdAt - pierwsza.createdAt);
  const najswiezsza = niezapisane[0];
  if (najswiezsza === undefined) return null;
  return {
    kopia: najswiezsza,
    zdanie:
      `Mam niezapisany dokument z ${new Date(najswiezsza.createdAt).toLocaleString('pl-PL')} ` +
      `(${kopieNazwaPowodu(najswiezsza.reason)}). Przywrócić? Przywrócenie DO NOWEGO DOKUMENTU ` +
      'nie rusza tego, co stoi w oknie — dopiero przywrócenie na miejsce je zamienia.' +
      (niezapisane.length > 1
        ? ` Kopii niosących zmiany niezapisane jest ${niezapisane.length}; ta jest najświeższa.`
        : ''),
  };
}

/** Zdanie o szeregach wersji — dwa szeregi, nie jedna lista z domysłem. */
export function kopieOpiszSzeregi(
  wersje: readonly StudioVersion[],
  liczbaOperatora: number,
  liczbaAutozapisu: number,
  idZalozycielskiej: string | undefined,
): string {
  const zalozycielska =
    idZalozycielskiej === undefined || idZalozycielskiej === ''
      ? 'Wersji założycielskiej rdzeń nie wskazał — dokument mógł jeszcze nie mieć pierwszego zapisu.'
      : `Wersja założycielska: ${idZalozycielskiej} — powrót do niej idzie jednym poleceniem, ` +
        'bez szukania jej w wykazie.';
  return (
    `Wersji w wykazie: ${wersje.length}. Szereg Operatora: ${liczbaOperatora}; ` +
    `szereg autozapisu: ${liczbaAutozapisu} — zapisy samoczynne NIE zaśmiecają historii ` +
    `Operatora, bo idą osobno. ${zalozycielska}`
  );
}

/**
 * Nazwa szeregu widoczna dla Operatora.
 *
 * Szereg jest tu szeregiem **zapytania**, nie cechą wiersza: `StudioVersion`
 * pola `series` NIE niesie, więc pojedyncza wersja nie mówi, z którego szeregu
 * pochodzi. Wykaz zawężony do jednego szeregu nazwać się daje, wykaz zbiorczy —
 * nie, i okno mówi to wprost, zamiast zgadywać po braku etykiety. Pozycja
 * kontraktu potrzebna do rozróżnienia wiersza stoi w sprawozdaniu.
 */
export function kopieNazwaSzeregu(szereg: StudioVersionSeries | undefined): string {
  if (szereg === StudioVersionSeries.Autosave) return 'zapis samoczynny';
  if (szereg === StudioVersionSeries.Operator) return 'wersja Operatora';
  return 'oba szeregi razem — wiersz sam szeregu nie niesie';
}
