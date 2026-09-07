// Czynności plikowe okna Developer: drzewo, otwarcie i zapis pliku, wersje,
// szukanie z zamianą, formatowanie, uwagi narzędzia, przekształcenia i symbole.
import {
  Command,
  ContextualOpKind,
  RefactorKind,
  SymbolNavigationKind,
  TreeNodeKind,
} from '../../../shared/contract.ts';
import type {
  DeveloperDiagnostic,
  DeveloperFileVersion,
  DeveloperGrepMatch,
  DeveloperSymbol,
  DeveloperTreeNode,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import {
  czesciWpisu,
  potwierdzone,
  wpisPasa,
  wypelnijPanel,
  zalozPas,
  zwiazPas,
} from './developer-pas.ts';

const NAGLOWEK = 'Pliki projektu';
const ZNACZNIK = 'devplik';

const PRZYCISKI: ReadonlyArray<readonly [string, string]> = [
  ['drzewo', 'Odczytaj drzewo'],
  ['otworz', 'Otwórz plik'],
  ['zapisz', 'Zapisz plik'],
  ['zaloz', 'Załóż plik'],
  ['usun', 'Usuń plik'],
  ['przemianuj', 'Przemianuj plik'],
  ['przenies', 'Przenieś plik'],
  ['wersje', 'Wykaz wersji'],
  ['przywroc', 'Przywróć wersję'],
  ['szukaj', 'Szukaj w plikach'],
  ['zamien', 'Zamień w plikach'],
  ['sformatuj', 'Sformatuj plik'],
  ['uwagi', 'Odczytaj uwagi'],
  ['przeksztalc', 'Przemianuj symbol'],
  ['symbol', 'Znajdź określenie'],
  ['wyjasnij', 'Wyjaśnij plik'],
];

export function zwiazPlikiDevelopera(
  kanal: Kanal,
  korzen: Element,
  idOkna: () => string,
  przy: AddEventListenerOptions,
): void {
  zalozPas(korzen, 'panel-drzewo', ZNACZNIK, 'Ścieżka, wzorzec albo treść', PRZYCISKI);
  zwiazPas(korzen, ZNACZNIK, (czynnosc) => {
    void wykonaj(kanal, korzen, czynnosc, idOkna());
  }, przy);
}

function wpis(korzen: Element): string {
  return wpisPasa(korzen, ZNACZNIK);
}

function czesci(korzen: Element): string[] {
  return czesciWpisu(korzen, ZNACZNIK);
}

async function wykonaj(
  kanal: Kanal,
  korzen: Element,
  czynnosc: string,
  idOkna: string,
): Promise<void> {
  if (idOkna === '') {
    oglos(NAGLOWEK, 'Rdzeń nie dał okna projektu dla tej karty.', 'ostrzezenie');
    return;
  }
  if (czynnosc === 'drzewo') return odczytajDrzewo(kanal, korzen, idOkna);
  if (czynnosc === 'otworz') return otworzPlik(kanal, korzen, idOkna);
  if (czynnosc === 'zapisz') return zapiszPlik(kanal, korzen, idOkna);
  if (czynnosc === 'zaloz') return zalozPlik(kanal, korzen, idOkna);
  if (czynnosc === 'usun') return usunPlik(kanal, korzen, idOkna);
  if (czynnosc === 'przemianuj') return przemianujPlik(kanal, korzen, idOkna);
  if (czynnosc === 'przenies') return przeniesPlik(kanal, korzen, idOkna);
  if (czynnosc === 'wersje') return wykazWersji(kanal, korzen, idOkna);
  if (czynnosc === 'przywroc') return przywrocWersje(kanal, korzen, idOkna);
  if (czynnosc === 'szukaj') return szukajWPlikach(kanal, korzen, idOkna);
  if (czynnosc === 'zamien') return zamienWPlikach(kanal, korzen, idOkna);
  if (czynnosc === 'sformatuj') return sformatujPlik(kanal, korzen, idOkna);
  if (czynnosc === 'uwagi') return odczytajUwagi(kanal, korzen, idOkna);
  if (czynnosc === 'przeksztalc') return przemianujSymbol(kanal, korzen, idOkna);
  if (czynnosc === 'symbol') return znajdzOkreslenie(kanal, korzen, idOkna);
  if (czynnosc === 'wyjasnij') return wyjasnijPlik(kanal, korzen, idOkna);
  /* Czynność spoza obsłużonych odmawia zamiast milczeć: przycisk bez gałęzi
     wyglądałby jak działający. */
  oglos(NAGLOWEK, `Czynność „${czynnosc}" nie jest prowadzona przez to okno.`, 'ostrzezenie');
}

async function odczytajDrzewo(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const sciezka = wpis(korzen);
  const wynik = await wywolaj(kanal, Command.DeveloperTreeGet, {
    windowId: idOkna,
    ...(sciezka === '' ? {} : { path: sciezka }),
    depth: 2,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu drzewa.', 'ostrzezenie');
    return;
  }
  const wezly = wynik.wynik.nodes;
  wypelnijPanel(korzen, 'panel-drzewo', wezly.map((wezel: DeveloperTreeNode) =>
    `${wezel.kind === TreeNodeKind.Directory ? '▸' : '·'} ${wezel.path}`));
  if (wezly.length === 0) oglos(NAGLOWEK, `Katalog ${wynik.wynik.root} jest pusty.`);
}

async function otworzPlik(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const sciezka = wpis(korzen);
  if (sciezka === '') {
    oglos(NAGLOWEK, 'Otwarcie pliku potrzebuje ścieżki wpisanej w polu.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.DeveloperFileOpen, { windowId: idOkna, path: sciezka });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił otwarcia pliku.', 'ostrzezenie');
    return;
  }
  wypelnijPanel(korzen, 'panel-edytor', (wynik.wynik.file.content ?? '').split('\n'));
  oglos(NAGLOWEK, `Plik ${wynik.wynik.file.path} otwarty w edytorze.`);
}

/* Zapis zakłada wersję, bo bez niej nadpisanie pliku byłoby nie do cofnięcia
   z samego okna. */
async function zapiszPlik(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const [sciezka = '', tresc = ''] = czesci(korzen);
  if (sciezka === '' || tresc === '') {
    oglos(NAGLOWEK,
      'Zapis potrzebuje ścieżki i treści, na przykład „src/plik.ts | export const x = 1;".',
      'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.DeveloperFileSave, {
    windowId: idOkna,
    path: sciezka,
    content: tresc,
    createVersion: true,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu pliku.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Plik ${wynik.wynik.file.path} zapisany wraz z wersją.`);
}

async function zalozPlik(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const [sciezka = '', rodzaj = ''] = czesci(korzen);
  if (sciezka === '') {
    oglos(NAGLOWEK,
      'Założenie potrzebuje ścieżki, a po znaku pionowym słowa „katalog" dla katalogu.',
      'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.DeveloperFileCreate, {
    windowId: idOkna,
    path: sciezka,
    kind: rodzaj === 'katalog' ? TreeNodeKind.Directory : TreeNodeKind.File,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił założenia pozycji.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Pozycja ${wynik.wynik.node.path} założona.`);
}

/* Usunięcie jest nieodwracalne, więc pierwsze naciśnięcie uzbraja przycisk;
   znacznik cofnięcia od rdzenia okno nazywa, bo daje drogę odwrotu. */
async function usunPlik(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const sciezka = wpis(korzen);
  if (sciezka === '') {
    oglos(NAGLOWEK, 'Usunięcie potrzebuje ścieżki wpisanej w polu.', 'ostrzezenie');
    return;
  }
  if (!potwierdzone(korzen, ZNACZNIK, 'usun')) return;
  const wynik = await wywolaj(kanal, Command.DeveloperFileDelete, {
    windowId: idOkna,
    paths: [sciezka],
    recursive: true,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił usunięcia pozycji.', 'ostrzezenie');
    return;
  }
  const znacznikCofniecia = wynik.wynik.undoToken;
  oglos(NAGLOWEK, znacznikCofniecia === undefined
    ? `Usunięto ${wynik.wynik.deletedPaths.length} pozycji bez drogi odwrotu.`
    : `Usunięto ${wynik.wynik.deletedPaths.length} pozycji; znacznik cofnięcia ${znacznikCofniecia}.`);
}

async function przemianujPlik(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const [sciezka = '', nazwa = ''] = czesci(korzen);
  if (sciezka === '' || nazwa === '') {
    oglos(NAGLOWEK,
      'Przemianowanie potrzebuje ścieżki i nowej nazwy, na przykład „src/stary.ts | nowy.ts".',
      'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.DeveloperFileRename, {
    windowId: idOkna,
    path: sciezka,
    newName: nazwa,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił przemianowania.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Pozycja stoi teraz pod ${wynik.wynik.node.path}.`);
}

async function przeniesPlik(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const [sciezka = '', cel = ''] = czesci(korzen);
  if (sciezka === '' || cel === '') {
    oglos(NAGLOWEK,
      'Przeniesienie potrzebuje ścieżki i katalogu celu, na przykład „src/plik.ts | src/stare".',
      'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.DeveloperFileMove, {
    windowId: idOkna,
    paths: [sciezka],
    targetPath: cel,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił przeniesienia.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Przeniesiono ${wynik.wynik.nodes.length} pozycji do ${cel}.`);
}

async function wykazWersji(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const sciezka = wpis(korzen);
  if (sciezka === '') {
    oglos(NAGLOWEK, 'Wykaz wersji potrzebuje ścieżki wpisanej w polu.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.DeveloperFileVersionList, {
    windowId: idOkna,
    path: sciezka,
    limit: 50,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu wersji.', 'ostrzezenie');
    return;
  }
  const wersje = wynik.wynik.versions;
  wypelnijPanel(korzen, 'panel-artefakty', wersje.map((wersja: DeveloperFileVersion) =>
    `${wersja.id} · ${wersja.path}`));
  if (wersje.length === 0) oglos(NAGLOWEK, 'Ten plik nie ma jeszcze zapisanych wersji.');
}

/* Przywrócenie nadpisuje plik na dysku, więc pierwsze naciśnięcie uzbraja
   przycisk; wskazaniem jest oznaczenie wersji z wykazu. */
async function przywrocWersje(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const oznaczenie = wpis(korzen);
  if (oznaczenie === '') {
    oglos(NAGLOWEK, 'Przywrócenie potrzebuje oznaczenia wersji z wykazu.', 'ostrzezenie');
    return;
  }
  if (!potwierdzone(korzen, ZNACZNIK, 'przywroc')) return;
  const wynik = await wywolaj(kanal, Command.DeveloperFileVersionRestore, {
    windowId: idOkna,
    versionId: oznaczenie,
    writeToDisk: true,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił przywrócenia wersji.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Plik ${wynik.wynik.file.path} przywrócony z wersji ${oznaczenie}.`);
}

async function szukajWPlikach(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const wzorzec = wpis(korzen);
  if (wzorzec === '') {
    oglos(NAGLOWEK, 'Szukanie potrzebuje wzorca wpisanego w polu.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.DeveloperGrepSearch, {
    windowId: idOkna,
    pattern: wzorzec,
    limit: 200,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił szukania.', 'ostrzezenie');
    return;
  }
  const trafienia = wynik.wynik.matches;
  wypelnijPanel(korzen, 'panel-drzewo', trafienia.map((trafienie: DeveloperGrepMatch) =>
    `${trafienie.path}:${trafienie.line} ${trafienie.text.trim()}`));
  if (trafienia.length === 0) {
    oglos(NAGLOWEK, 'Żaden plik nie zawiera tego wzorca.');
    return;
  }
  if (wynik.wynik.truncated) {
    oglos(NAGLOWEK, 'Wykaz przycięty granicą wyszukiwania — trafień jest więcej.', 'ostrzezenie');
  }
}

/* Zamiana idzie najpierw jako podgląd: wykaz zmian wraca bez tknięcia plików,
   a dopiero drugie naciśnięcie zapisuje je na dysk. */
async function zamienWPlikach(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const [wzorzec = '', zamiennik = ''] = czesci(korzen);
  if (wzorzec === '') {
    oglos(NAGLOWEK,
      'Zamiana potrzebuje wzorca i treści zastępczej, na przykład „stary | nowy".',
      'ostrzezenie');
    return;
  }
  const podglad = !potwierdzone(korzen, ZNACZNIK, 'zamien');
  const wynik = await wywolaj(kanal, Command.DeveloperGrepReplace, {
    windowId: idOkna,
    pattern: wzorzec,
    replacement: zamiennik,
    preview: podglad,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zamiany.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, wynik.wynik.applied
    ? `Zamiana zapisana w ${wynik.wynik.changedPaths.length} plikach.`
    : `Podgląd: zamiana dotknie ${wynik.wynik.changedPaths.length} plików `
      + `w ${wynik.wynik.edits.length} miejscach. Naciśnij ponownie, aby zapisać.`);
}

async function sformatujPlik(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const sciezka = wpis(korzen);
  if (sciezka === '') {
    oglos(NAGLOWEK, 'Formatowanie potrzebuje ścieżki wpisanej w polu.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.DeveloperFormatRun, {
    windowId: idOkna,
    path: sciezka,
    writeToDisk: true,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił formatowania.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, wynik.wynik.changed
    ? `Plik ${sciezka} sformatowany narzędziem ${wynik.wynik.formatter ?? 'rdzenia'}.`
    : `Plik ${sciezka} był już sformatowany.`);
}

async function odczytajUwagi(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const sciezka = wpis(korzen);
  const wynik = await wywolaj(kanal, Command.DeveloperLintGet, {
    windowId: idOkna,
    ...(sciezka === '' ? {} : { paths: [sciezka] }),
    limit: 200,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu uwag.', 'ostrzezenie');
    return;
  }
  if (!wynik.wynik.linterAvailable) {
    oglos(NAGLOWEK, 'Narzędzie uwag nie stoi na maszynie rdzenia — nikt kodu nie sprawdził.',
      'ostrzezenie');
    return;
  }
  const uwagi = wynik.wynik.diagnostics;
  wypelnijPanel(korzen, 'panel-plan', uwagi.map((uwaga: DeveloperDiagnostic) =>
    `${uwaga.severity} ${uwaga.path}:${uwaga.line} ${uwaga.message}`));
  if (uwagi.length === 0) oglos(NAGLOWEK, 'Narzędzie nie ma uwag.');
}

/* Przemianowanie określenia idzie z podglądem, tak jak zamiana w plikach:
   pierwsze naciśnięcie pokazuje zmiany, drugie je zapisuje. */
async function przemianujSymbol(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const [sciezka = '', wiersz = '', kolumna = '', nazwa = ''] = czesci(korzen);
  const numerWiersza = Number.parseInt(wiersz, 10);
  const numerKolumny = Number.parseInt(kolumna, 10);
  if (sciezka === '' || !Number.isFinite(numerWiersza) || !Number.isFinite(numerKolumny)
    || nazwa === '') {
    oglos(NAGLOWEK,
      'Przemianowanie potrzebuje ścieżki, wiersza, kolumny i nowej nazwy, '
      + 'na przykład „src/plik.ts | 12 | 5 | nowaNazwa".',
      'ostrzezenie');
    return;
  }
  const podglad = !potwierdzone(korzen, ZNACZNIK, 'przeksztalc');
  const wynik = await wywolaj(kanal, Command.DeveloperRefactorApply, {
    windowId: idOkna,
    path: sciezka,
    line: numerWiersza,
    column: numerKolumny,
    kind: RefactorKind.Rename,
    newName: nazwa,
    preview: podglad,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił przemianowania określenia.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, wynik.wynik.applied
    ? `Określenie przemianowane w ${wynik.wynik.changedPaths?.length ?? 0} plikach.`
    : `Podgląd: przemianowanie dotknie ${wynik.wynik.edits.length} miejsc. `
      + 'Naciśnij ponownie, aby zapisać.');
}

async function znajdzOkreslenie(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const [sciezka = '', wiersz = '', kolumna = ''] = czesci(korzen);
  const numerWiersza = Number.parseInt(wiersz, 10);
  const numerKolumny = Number.parseInt(kolumna, 10);
  if (sciezka === '' || !Number.isFinite(numerWiersza) || !Number.isFinite(numerKolumny)) {
    oglos(NAGLOWEK,
      'Szukanie określenia potrzebuje ścieżki, wiersza i kolumny, na przykład „src/plik.ts | 12 | 5".',
      'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.DeveloperSymbolNavigate, {
    windowId: idOkna,
    path: sciezka,
    line: numerWiersza,
    column: numerKolumny,
    kind: SymbolNavigationKind.Definition,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił szukania określenia.', 'ostrzezenie');
    return;
  }
  if (!wynik.wynik.serverAvailable) {
    oglos(NAGLOWEK, 'Usługa języka nie stoi na maszynie rdzenia — określenia nikt nie wskazał.',
      'ostrzezenie');
    return;
  }
  const okreslenia = wynik.wynik.symbols;
  wypelnijPanel(korzen, 'panel-plan', okreslenia.map((okreslenie: DeveloperSymbol) =>
    `${okreslenie.name} · ${okreslenie.path}:${okreslenie.line}`));
  if (okreslenia.length === 0) oglos(NAGLOWEK, 'Usługa języka nie zna określenia w tym miejscu.');
}

async function wyjasnijPlik(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const sciezka = wpis(korzen);
  if (sciezka === '') {
    oglos(NAGLOWEK, 'Wyjaśnienie potrzebuje ścieżki wpisanej w polu.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.DeveloperContextualOp, {
    windowId: idOkna,
    operation: ContextualOpKind.Explain,
    path: sciezka,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wyjaśnienia.', 'ostrzezenie');
    return;
  }
  wypelnijPanel(korzen, 'panel-edytor', wynik.wynik.result.split('\n'));
  oglos(NAGLOWEK, `Wyjaśnienie pliku ${sciezka} stoi w edytorze.`);
}
