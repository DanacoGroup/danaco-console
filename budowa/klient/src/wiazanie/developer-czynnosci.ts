// Czynności okna Developer rozłożone na panele: pliki w drzewie, kontrola
// i praca krokowa przy budowaniu, zasoby przy pojemnikach.
import {
  Command,
  ConflictResolutionKind,
  ContainerActionKind,
  ContextualOpKind,
  DataEngine,
  DebugStepKind,
  GitActionKind,
  RefactorKind,
  ScanKind,
  SymbolNavigationKind,
  TreeNodeKind,
} from '../../../shared/contract.ts';
import type {
  ContainerInfo,
  DataSchemaNode,
  DebugFrame,
  DeveloperCoverage,
  DeveloperDiagnostic,
  DeveloperFileVersion,
  DeveloperGrepMatch,
  DeveloperSymbol,
  DeveloperTestResult,
  DeveloperTreeNode,
  GitCommit,
  GitStatusEntry,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { dolozCzynnosciPanelu, zapytajWSzufladzie } from './czynnosci-okna.ts';
import { oglos } from './ogloszenie.ts';

const NAGLOWEK = 'Projekt';

/* Praca krokowa żyje w oknie, bo kontrakt nie ma wykazu jej przebiegów:
   oznaczenie bierze się z rozpoczęcia i służy krokom oraz wyliczeniu. */
const PRZEBIEGI = new Map<string, string>();

interface Otoczenie {
  kanal: Kanal;
  korzen: Element;
  idOkna: () => string;
  odswiez: () => void;
}

export function zwiazCzynnosciDevelopera(
  kanal: Kanal,
  korzen: Element,
  idOkna: () => string,
  odswiez: () => void,
  dolacz: (zdejmij: () => void) => void,
  przy: AddEventListenerOptions,
): void {
  const otoczenie: Otoczenie = { kanal, korzen, idOkna, odswiez };
  const zdejmowanie = [
    dolozCzynnosciPanelu(korzen, 'panel-drzewo', 'Czynności plików', [
      {
        naglowek: 'Drzewo i pliki',
        pozycje: [
          { kod: 'drzewo', nazwa: 'Odczytaj drzewo…' },
          { kod: 'otworz', nazwa: 'Otwórz plik…' },
          { kod: 'zapisz', nazwa: 'Zapisz plik…' },
          { kod: 'zaloz', nazwa: 'Załóż plik albo katalog…' },
          { kod: 'przemianuj', nazwa: 'Przemianuj…' },
          { kod: 'przenies', nazwa: 'Przenieś…' },
          { kod: 'usun', nazwa: 'Usuń…' },
        ],
      },
      {
        naglowek: 'Wersje pliku',
        pozycje: [
          { kod: 'wersje', nazwa: 'Wykaz wersji…' },
          { kod: 'przywroc', nazwa: 'Przywróć wersję…' },
        ],
      },
      {
        naglowek: 'Treść',
        pozycje: [
          { kod: 'szukaj', nazwa: 'Szukaj w plikach…' },
          { kod: 'zamien', nazwa: 'Zamień w plikach…' },
          { kod: 'sformatuj', nazwa: 'Sformatuj plik…' },
          { kod: 'uwagi', nazwa: 'Odczytaj uwagi narzędzia…' },
          { kod: 'przemianuj-symbol', nazwa: 'Przemianuj określenie…' },
          { kod: 'symbol', nazwa: 'Znajdź określenie…' },
          { kod: 'wyjasnij', nazwa: 'Wyjaśnij plik…' },
        ],
      },
    ], (kod) => {
      void wykonaj(otoczenie, kod, 'panel-drzewo');
    }, przy),
    dolozCzynnosciPanelu(korzen, 'panel-build', 'Czynności kontroli', [
      {
        naglowek: 'Budowanie',
        pozycje: [
          { kod: 'budowanie', nazwa: 'Uruchom budowanie…' },
          { kod: 'budowanie-dziennik', nazwa: 'Dziennik ostatniego budowania' },
          { kod: 'proby', nazwa: 'Wyniki prób' },
          { kod: 'pokrycie', nazwa: 'Pokrycie prób' },
        ],
      },
      {
        naglowek: 'Kontrola',
        pozycje: [
          { kod: 'przeglad', nazwa: 'Przegląd bezpieczeństwa…' },
          { kod: 'narzedzia', nazwa: 'Sprawdź narzędzia…' },
        ],
      },
      {
        naglowek: 'Praca krokowa',
        pozycje: [
          { kod: 'krokowa', nazwa: 'Zacznij pracę krokową…' },
          { kod: 'krok', nazwa: 'Krok dalej' },
          { kod: 'zatrzymanie', nazwa: 'Postaw zatrzymanie…' },
          { kod: 'wylicz', nazwa: 'Wylicz wyrażenie…' },
          { kod: 'zasieg', nazwa: 'Odczytaj zasięg' },
        ],
      },
    ], (kod) => {
      void wykonaj(otoczenie, kod, 'panel-build');
    }, przy),
    dolozCzynnosciPanelu(korzen, 'panel-git', 'Czynności repozytorium', [
      {
        naglowek: 'Stan i dzieje',
        pozycje: [
          { kod: 'git-stan', nazwa: 'Stan repozytorium' },
          { kod: 'git-roznica', nazwa: 'Różnica zmian' },
          { kod: 'git-dzieje', nazwa: 'Dzieje gałęzi…' },
        ],
      },
      {
        naglowek: 'Zmiany',
        pozycje: [
          { kod: 'git-zapisz', nazwa: 'Zapisz zmiany…' },
          { kod: 'git-spor', nazwa: 'Odczytaj spór…' },
          { kod: 'git-spor-rozstrzygnij', nazwa: 'Rozstrzygnij spór…' },
        ],
      },
    ], (kod) => {
      void wykonaj(otoczenie, kod, 'panel-git');
    }, przy),
    dolozCzynnosciPanelu(korzen, 'panel-terminal', 'Czynności zasobów', [
      {
        naglowek: 'Pojemniki',
        pozycje: [
          { kod: 'pojemnik', nazwa: 'Przestaw pojemnik…' },
          { kod: 'zestaw', nazwa: 'Podnieś zestaw…' },
          { kod: 'obraz', nazwa: 'Zbuduj obraz…' },
        ],
      },
      {
        naglowek: 'Bazy danych',
        pozycje: [
          { kod: 'polaczenie', nazwa: 'Zapisz połączenie…' },
          { kod: 'zapytanie', nazwa: 'Wykonaj zapytanie…' },
          { kod: 'schemat', nazwa: 'Odczytaj schemat' },
          { kod: 'przeniesienia', nazwa: 'Wykonaj przeniesienia…' },
        ],
      },
      {
        naglowek: 'Żądania sieciowe',
        pozycje: [
          { kod: 'zadanie', nazwa: 'Wyślij żądanie…' },
          { kod: 'zbior', nazwa: 'Zapisz zbiór żądań…' },
          { kod: 'openapi', nazwa: 'Wciągnij opis OpenAPI…' },
          { kod: 'obciazenie', nazwa: 'Zbadaj obciążenie…' },
        ],
      },
    ], (kod) => {
      void wykonaj(otoczenie, kod, 'panel-terminal');
    }, przy),
  ];
  for (const zdejmij of zdejmowanie) if (zdejmij !== null) dolacz(zdejmij);
}

function wypelnij(korzen: Element, panelKod: string, wiersze: string[]): void {
  const cialo = korzen.querySelector(`#${panelKod} .sta-okno-tresc`);
  if (cialo === null) return;
  cialo.replaceChildren(...wiersze.map((tresc) => {
    const wiersz = cialo.ownerDocument.createElement('div');
    wiersz.className = 'dn-wykaz-modulu-poz';
    wiersz.textContent = tresc;
    return wiersz;
  }));
}

async function wykonaj(otoczenie: Otoczenie, kod: string, panel: string): Promise<void> {
  if (otoczenie.idOkna() === '') {
    oglos(NAGLOWEK, 'Rdzeń nie dał okna projektu dla tej karty.', 'ostrzezenie');
    return;
  }
  if (kod === 'drzewo') return odczytajDrzewo(otoczenie, panel);
  if (kod === 'otworz') return otworzPlik(otoczenie, panel);
  if (kod === 'zapisz') return zapiszPlik(otoczenie, panel);
  if (kod === 'zaloz') return zalozPozycje(otoczenie, panel);
  if (kod === 'przemianuj') return przemianujPlik(otoczenie, panel);
  if (kod === 'przenies') return przeniesPlik(otoczenie, panel);
  if (kod === 'usun') return usunPozycje(otoczenie, panel);
  if (kod === 'wersje') return wykazWersji(otoczenie, panel);
  if (kod === 'przywroc') return przywrocWersje(otoczenie, panel);
  if (kod === 'szukaj') return szukajWPlikach(otoczenie, panel);
  if (kod === 'zamien') return zamienWPlikach(otoczenie, panel);
  if (kod === 'sformatuj') return sformatujPlik(otoczenie, panel);
  if (kod === 'uwagi') return odczytajUwagi(otoczenie, panel);
  if (kod === 'przemianuj-symbol') return przemianujSymbol(otoczenie, panel);
  if (kod === 'symbol') return znajdzOkreslenie(otoczenie, panel);
  if (kod === 'wyjasnij') return wyjasnijPlik(otoczenie, panel);
  if (kod === 'budowanie') return uruchomBudowanie(otoczenie, panel);
  if (kod === 'budowanie-dziennik') return dziennikBudowania(otoczenie);
  if (kod === 'proby') return wynikiProb(otoczenie);
  if (kod === 'pokrycie') return pokrycieProb(otoczenie);
  if (kod === 'przeglad') return przegladBezpieczenstwa(otoczenie, panel);
  if (kod === 'narzedzia') return sprawdzNarzedzia(otoczenie, panel);
  if (kod === 'krokowa') return zacznijPraceKrokowa(otoczenie, panel);
  if (kod === 'krok') return krokDalej(otoczenie);
  if (kod === 'zatrzymanie') return postawZatrzymanie(otoczenie, panel);
  if (kod === 'wylicz') return wyliczWyrazenie(otoczenie, panel);
  if (kod === 'zasieg') return odczytajZasieg(otoczenie);
  if (kod === 'git-stan') return stanRepozytorium(otoczenie);
  if (kod === 'git-roznica') return roznicaZmian(otoczenie);
  if (kod === 'git-dzieje') return dziejeGalezi(otoczenie, panel);
  if (kod === 'git-zapisz') return zapiszZmiany(otoczenie, panel);
  if (kod === 'git-spor') return odczytajSpor(otoczenie, panel);
  if (kod === 'git-spor-rozstrzygnij') return rozstrzygnijSpor(otoczenie, panel);
  if (kod === 'pojemnik') return przestawPojemnik(otoczenie, panel);
  if (kod === 'zestaw') return podniesZestaw(otoczenie, panel);
  if (kod === 'obraz') return zbudujObraz(otoczenie, panel);
  if (kod === 'polaczenie') return zapiszPolaczenie(otoczenie, panel);
  if (kod === 'zapytanie') return wykonajZapytanie(otoczenie, panel);
  if (kod === 'schemat') return odczytajSchemat(otoczenie);
  if (kod === 'przeniesienia') return wykonajPrzeniesienia(otoczenie, panel);
  if (kod === 'zadanie') return wyslijZadanie(otoczenie, panel);
  if (kod === 'zbior') return zapiszZbior(otoczenie, panel);
  if (kod === 'openapi') return wciagnijOpenapi(otoczenie, panel);
  if (kod === 'obciazenie') return zbadajObciazenie(otoczenie, panel);
  /* Czynność spoza obsłużonych odmawia zamiast milczeć: pozycja bez gałęzi
     wyglądałaby jak działająca. */
  oglos(NAGLOWEK, `Czynność „${kod}" nie jest prowadzona przez to okno.`, 'ostrzezenie');
}

async function odczytajDrzewo({ kanal, korzen, idOkna }: Otoczenie, panel: string):
Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Drzewo plików',
    pola: [
      { klucz: 'sciezka', etykieta: 'Katalog', podpowiedz: 'Puste bierze katalog roboczy' },
      { klucz: 'glebokosc', etykieta: 'Głębokość', wartosc: '2' },
    ],
    wykonanie: 'Odczytaj drzewo',
  });
  if (wartosci === null) return;
  const glebokosc = Number.parseInt(wartosci.glebokosc ?? '', 10);
  const wynik = await wywolaj(kanal, Command.DeveloperTreeGet, {
    windowId: idOkna(),
    ...(wartosci.sciezka === '' ? {} : { path: wartosci.sciezka }),
    depth: Number.isFinite(glebokosc) ? glebokosc : 2,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu drzewa.', 'ostrzezenie');
    return;
  }
  const wezly = wynik.wynik.nodes;
  wypelnij(korzen, 'panel-drzewo', wezly.map((wezel: DeveloperTreeNode) =>
    `${wezel.kind === TreeNodeKind.Directory ? '▸' : '·'} ${wezel.path}`));
  if (wezly.length === 0) oglos(NAGLOWEK, `Katalog ${wynik.wynik.root} jest pusty.`);
}

async function otworzPlik({ kanal, korzen, idOkna }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Otwarcie pliku',
    pola: [{ klucz: 'sciezka', etykieta: 'Ścieżka pliku', wymagane: true }],
    wykonanie: 'Otwórz plik',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.DeveloperFileOpen, {
    windowId: idOkna(),
    path: wartosci.sciezka ?? '',
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił otwarcia pliku.', 'ostrzezenie');
    return;
  }
  wypelnij(korzen, 'panel-edytor', (wynik.wynik.file.content ?? '').split('\n'));
  oglos(NAGLOWEK, `Plik ${wynik.wynik.file.path} otwarty w edytorze.`);
}

/* Zapis zakłada wersję, bo bez niej nadpisania nie da się cofnąć z okna. */
async function zapiszPlik({ kanal, korzen, idOkna }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Zapis pliku',
    opis: 'Zapis zakłada wersję, więc poprzednią treść da się przywrócić.',
    pola: [
      { klucz: 'sciezka', etykieta: 'Ścieżka pliku', wymagane: true },
      { klucz: 'tresc', etykieta: 'Treść', obszerne: true, wymagane: true },
    ],
    wykonanie: 'Zapisz plik',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.DeveloperFileSave, {
    windowId: idOkna(),
    path: wartosci.sciezka ?? '',
    content: wartosci.tresc ?? '',
    createVersion: true,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu pliku.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Plik ${wynik.wynik.file.path} zapisany wraz z wersją.`);
}

async function zalozPozycje({ kanal, korzen, idOkna }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Nowa pozycja drzewa',
    pola: [
      { klucz: 'sciezka', etykieta: 'Ścieżka', wymagane: true },
      {
        klucz: 'rodzaj',
        etykieta: 'Rodzaj',
        wybor: [[TreeNodeKind.File, 'Plik'], [TreeNodeKind.Directory, 'Katalog']],
      },
    ],
    wykonanie: 'Załóż pozycję',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.DeveloperFileCreate, {
    windowId: idOkna(),
    path: wartosci.sciezka ?? '',
    kind: (wartosci.rodzaj ?? TreeNodeKind.File) as TreeNodeKind,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił założenia pozycji.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Pozycja ${wynik.wynik.node.path} założona.`);
}

async function przemianujPlik({ kanal, korzen, idOkna }: Otoczenie, panel: string):
Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Przemianowanie pozycji',
    pola: [
      { klucz: 'sciezka', etykieta: 'Ścieżka pozycji', wymagane: true },
      { klucz: 'nazwa', etykieta: 'Nowa nazwa', wymagane: true },
    ],
    wykonanie: 'Przemianuj',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.DeveloperFileRename, {
    windowId: idOkna(),
    path: wartosci.sciezka ?? '',
    newName: wartosci.nazwa ?? '',
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił przemianowania.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Pozycja stoi teraz pod ${wynik.wynik.node.path}.`);
}

async function przeniesPlik({ kanal, korzen, idOkna }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Przeniesienie pozycji',
    pola: [
      { klucz: 'sciezka', etykieta: 'Ścieżka pozycji', wymagane: true },
      { klucz: 'cel', etykieta: 'Katalog celu', wymagane: true },
    ],
    wykonanie: 'Przenieś',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.DeveloperFileMove, {
    windowId: idOkna(),
    paths: [wartosci.sciezka ?? ''],
    targetPath: wartosci.cel ?? '',
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił przeniesienia.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Przeniesiono ${wynik.wynik.nodes.length} pozycji.`);
}

async function usunPozycje({ kanal, korzen, idOkna }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Usunięcie pozycji',
    pola: [{ klucz: 'sciezka', etykieta: 'Ścieżka pozycji', wymagane: true }],
    wykonanie: 'Usuń pozycję',
    nieodwracalne: 'Katalog schodzi wraz z zawartością. Rdzeń może oddać znacznik cofnięcia.',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.DeveloperFileDelete, {
    windowId: idOkna(),
    paths: [wartosci.sciezka ?? ''],
    recursive: true,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił usunięcia pozycji.', 'ostrzezenie');
    return;
  }
  const znacznik = wynik.wynik.undoToken;
  oglos(NAGLOWEK, znacznik === undefined
    ? `Usunięto ${wynik.wynik.deletedPaths.length} pozycji bez drogi odwrotu.`
    : `Usunięto ${wynik.wynik.deletedPaths.length} pozycji; znacznik cofnięcia ${znacznik}.`);
}

async function wykazWersji({ kanal, korzen, idOkna }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Wersje pliku',
    pola: [{ klucz: 'sciezka', etykieta: 'Ścieżka pliku', wymagane: true }],
    wykonanie: 'Odczytaj wersje',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.DeveloperFileVersionList, {
    windowId: idOkna(),
    path: wartosci.sciezka ?? '',
    limit: 50,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu wersji.', 'ostrzezenie');
    return;
  }
  const wersje = wynik.wynik.versions;
  wypelnij(korzen, 'panel-artefakty', wersje.map((wersja: DeveloperFileVersion) =>
    `${wersja.id} · ${wersja.path}`));
  if (wersje.length === 0) oglos(NAGLOWEK, 'Ten plik nie ma jeszcze zapisanych wersji.');
}

async function przywrocWersje({ kanal, korzen, idOkna }: Otoczenie, panel: string):
Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Przywrócenie wersji',
    pola: [{ klucz: 'wersja', etykieta: 'Oznaczenie wersji', wymagane: true }],
    wykonanie: 'Przywróć wersję',
    nieodwracalne: 'Treść pliku na dysku zostanie nadpisana treścią wskazanej wersji.',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.DeveloperFileVersionRestore, {
    windowId: idOkna(),
    versionId: wartosci.wersja ?? '',
    writeToDisk: true,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił przywrócenia wersji.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Plik ${wynik.wynik.file.path} przywrócony.`);
}

async function szukajWPlikach({ kanal, korzen, idOkna }: Otoczenie, panel: string):
Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Szukanie w plikach',
    pola: [
      { klucz: 'wzorzec', etykieta: 'Wzorzec', wymagane: true },
      {
        klucz: 'wyrazenie',
        etykieta: 'Wzorzec jest wyrażeniem regularnym',
        wybor: [['nie', 'Nie'], ['tak', 'Tak']],
      },
    ],
    wykonanie: 'Szukaj',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.DeveloperGrepSearch, {
    windowId: idOkna(),
    pattern: wartosci.wzorzec ?? '',
    regex: wartosci.wyrazenie === 'tak',
    limit: 200,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił szukania.', 'ostrzezenie');
    return;
  }
  const trafienia = wynik.wynik.matches;
  wypelnij(korzen, 'panel-drzewo', trafienia.map((trafienie: DeveloperGrepMatch) =>
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
   a zapis wymaga drugiego naciśnięcia. */
async function zamienWPlikach({ kanal, korzen, idOkna }: Otoczenie, panel: string):
Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Zamiana w plikach',
    opis: 'Rdzeń najpierw oddaje podgląd zmian; zapis wymaga potwierdzenia.',
    pola: [
      { klucz: 'wzorzec', etykieta: 'Wzorzec', wymagane: true },
      { klucz: 'zamiennik', etykieta: 'Treść zastępcza' },
    ],
    wykonanie: 'Zamień w plikach',
    nieodwracalne: 'Zapis zmieni treść plików na dysku.',
  });
  if (wartosci === null) return;
  const podglad = await wywolaj(kanal, Command.DeveloperGrepReplace, {
    windowId: idOkna(),
    pattern: wartosci.wzorzec ?? '',
    replacement: wartosci.zamiennik ?? '',
    preview: true,
  });
  if (!podglad.udany || podglad.wynik === undefined) {
    oglos(NAGLOWEK, podglad.blad?.message ?? 'Rdzeń odmówił zamiany.', 'ostrzezenie');
    return;
  }
  if (podglad.wynik.changedPaths.length === 0) {
    oglos(NAGLOWEK, 'Żaden plik nie zawiera tego wzorca.');
    return;
  }
  const wynik = await wywolaj(kanal, Command.DeveloperGrepReplace, {
    windowId: idOkna(),
    pattern: wartosci.wzorzec ?? '',
    replacement: wartosci.zamiennik ?? '',
    preview: false,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu zamiany.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Zamiana zapisana w ${wynik.wynik.changedPaths.length} plikach.`);
}

async function sformatujPlik({ kanal, korzen, idOkna }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Formatowanie pliku',
    pola: [{ klucz: 'sciezka', etykieta: 'Ścieżka pliku', wymagane: true }],
    wykonanie: 'Sformatuj plik',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.DeveloperFormatRun, {
    windowId: idOkna(),
    path: wartosci.sciezka ?? '',
    writeToDisk: true,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił formatowania.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, wynik.wynik.changed
    ? `Plik sformatowany narzędziem ${wynik.wynik.formatter ?? 'rdzenia'}.`
    : 'Plik był już sformatowany.');
}

async function odczytajUwagi({ kanal, korzen, idOkna }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Uwagi narzędzia',
    pola: [{ klucz: 'sciezka', etykieta: 'Ścieżka', podpowiedz: 'Puste bierze cały projekt' }],
    wykonanie: 'Odczytaj uwagi',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.DeveloperLintGet, {
    windowId: idOkna(),
    ...(wartosci.sciezka === '' ? {} : { paths: [wartosci.sciezka ?? ''] }),
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
  wypelnij(korzen, 'panel-plan', uwagi.map((uwaga: DeveloperDiagnostic) =>
    `${uwaga.severity} ${uwaga.path}:${uwaga.line} ${uwaga.message}`));
  if (uwagi.length === 0) oglos(NAGLOWEK, 'Narzędzie nie ma uwag.');
}

async function przemianujSymbol({ kanal, korzen, idOkna }: Otoczenie, panel: string):
Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Przemianowanie określenia',
    pola: [
      { klucz: 'sciezka', etykieta: 'Ścieżka pliku', wymagane: true },
      { klucz: 'wiersz', etykieta: 'Wiersz', wymagane: true },
      { klucz: 'kolumna', etykieta: 'Kolumna', wymagane: true },
      { klucz: 'nazwa', etykieta: 'Nowa nazwa', wymagane: true },
    ],
    wykonanie: 'Przemianuj określenie',
    nieodwracalne: 'Zmiana dotknie wszystkich miejsc użycia w projekcie.',
  });
  if (wartosci === null) return;
  const wiersz = Number.parseInt(wartosci.wiersz ?? '', 10);
  const kolumna = Number.parseInt(wartosci.kolumna ?? '', 10);
  if (!Number.isFinite(wiersz) || !Number.isFinite(kolumna)) {
    oglos(NAGLOWEK, 'Wiersz i kolumna muszą być liczbami.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.DeveloperRefactorApply, {
    windowId: idOkna(),
    path: wartosci.sciezka ?? '',
    line: wiersz,
    column: kolumna,
    kind: RefactorKind.Rename,
    newName: wartosci.nazwa ?? '',
    preview: false,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił przemianowania.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Określenie przemianowane w ${wynik.wynik.changedPaths?.length ?? 0} plikach.`);
}

async function znajdzOkreslenie({ kanal, korzen, idOkna }: Otoczenie, panel: string):
Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Szukanie określenia',
    pola: [
      { klucz: 'sciezka', etykieta: 'Ścieżka pliku', wymagane: true },
      { klucz: 'wiersz', etykieta: 'Wiersz', wymagane: true },
      { klucz: 'kolumna', etykieta: 'Kolumna', wymagane: true },
      {
        klucz: 'rodzaj',
        etykieta: 'Czego szukać',
        wybor: [
          [SymbolNavigationKind.Definition, 'Miejsca określenia'],
          [SymbolNavigationKind.References, 'Miejsc użycia'],
          [SymbolNavigationKind.Implementation, 'Wykonania'],
        ],
      },
    ],
    wykonanie: 'Szukaj',
  });
  if (wartosci === null) return;
  const wiersz = Number.parseInt(wartosci.wiersz ?? '', 10);
  const kolumna = Number.parseInt(wartosci.kolumna ?? '', 10);
  if (!Number.isFinite(wiersz) || !Number.isFinite(kolumna)) {
    oglos(NAGLOWEK, 'Wiersz i kolumna muszą być liczbami.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.DeveloperSymbolNavigate, {
    windowId: idOkna(),
    path: wartosci.sciezka ?? '',
    line: wiersz,
    column: kolumna,
    kind: (wartosci.rodzaj ?? SymbolNavigationKind.Definition) as SymbolNavigationKind,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił szukania określenia.', 'ostrzezenie');
    return;
  }
  if (!wynik.wynik.serverAvailable) {
    oglos(NAGLOWEK, 'Usługa języka nie stoi na maszynie rdzenia.', 'ostrzezenie');
    return;
  }
  const okreslenia = wynik.wynik.symbols;
  wypelnij(korzen, 'panel-plan', okreslenia.map((okreslenie: DeveloperSymbol) =>
    `${okreslenie.name} · ${okreslenie.path}:${okreslenie.line}`));
  if (okreslenia.length === 0) oglos(NAGLOWEK, 'Usługa języka nie zna nic w tym miejscu.');
}

async function wyjasnijPlik({ kanal, korzen, idOkna }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Wyjaśnienie pliku',
    pola: [{ klucz: 'sciezka', etykieta: 'Ścieżka pliku', wymagane: true }],
    wykonanie: 'Wyjaśnij plik',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.DeveloperContextualOp, {
    windowId: idOkna(),
    operation: ContextualOpKind.Explain,
    path: wartosci.sciezka ?? '',
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wyjaśnienia.', 'ostrzezenie');
    return;
  }
  wypelnij(korzen, 'panel-edytor', wynik.wynik.result.split('\n'));
  oglos(NAGLOWEK, 'Wyjaśnienie stoi w edytorze.');
}

/* Dziennik, wyniki prób i pokrycie idą z budowania stojącego w wykazie
   najwyżej: okno nie prowadzi wskazania budowania. */
async function ostatnieBudowanie(kanal: Kanal, idOkna: string): Promise<string> {
  const wykaz = await wywolaj(kanal, Command.DeveloperBuildList, { windowId: idOkna });
  const cel = wykaz.wynik?.builds[0]?.id ?? '';
  if (cel === '') oglos(NAGLOWEK, 'Żadne budowanie nie ruszyło w tym oknie.', 'ostrzezenie');
  return cel;
}

async function uruchomBudowanie({ kanal, korzen, idOkna, odswiez }: Otoczenie, panel: string):
Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Budowanie',
    pola: [
      { klucz: 'zlecenie', etykieta: 'Zlecenie', podpowiedz: 'build', wymagane: true },
      { klucz: 'argumenty', etykieta: 'Argumenty', podpowiedz: 'Rozdzielone spacją' },
    ],
    wykonanie: 'Uruchom budowanie',
  });
  if (wartosci === null) return;
  const argumenty = (wartosci.argumenty ?? '').split(/\s+/).filter((slowo) => slowo !== '');
  const wynik = await wywolaj(kanal, Command.DeveloperBuildRun, {
    windowId: idOkna(),
    task: wartosci.zlecenie ?? '',
    ...(argumenty.length === 0 ? {} : { arguments: argumenty }),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił uruchomienia budowania.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Budowanie w stanie ${wynik.wynik.build.status}.`);
  odswiez();
}

async function dziennikBudowania({ kanal, korzen, idOkna }: Otoczenie): Promise<void> {
  const cel = await ostatnieBudowanie(kanal, idOkna());
  if (cel === '') return;
  const wynik = await wywolaj(kanal, Command.DeveloperBuildLogGet, { buildId: cel, tail: 200 });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu dziennika.', 'ostrzezenie');
    return;
  }
  wypelnij(korzen, 'panel-build', wynik.wynik.lines);
  if (wynik.wynik.truncated) {
    oglos(NAGLOWEK, 'Dziennik przycięty granicą bufora — wierszy jest więcej.', 'ostrzezenie');
  }
}

async function wynikiProb({ kanal, korzen, idOkna }: Otoczenie): Promise<void> {
  const cel = await ostatnieBudowanie(kanal, idOkna());
  if (cel === '') return;
  const wynik = await wywolaj(kanal, Command.DeveloperTestResultGet, { buildId: cel });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu wyników prób.', 'ostrzezenie');
    return;
  }
  wypelnij(korzen, 'panel-zadania', wynik.wynik.results.map((proba: DeveloperTestResult) =>
    `${proba.status} ${proba.name}`));
  oglos(NAGLOWEK, `Próby: ${wynik.wynik.passed} zdanych, ${wynik.wynik.failed} niezdanych, `
    + `${wynik.wynik.skipped} pominiętych.`,
  wynik.wynik.failed === 0 ? 'informacja' : 'ostrzezenie');
}

async function pokrycieProb({ kanal, korzen, idOkna }: Otoczenie): Promise<void> {
  const cel = await ostatnieBudowanie(kanal, idOkna());
  if (cel === '') return;
  const wynik = await wywolaj(kanal, Command.DeveloperCoverageGet, { buildId: cel });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu pokrycia.', 'ostrzezenie');
    return;
  }
  wypelnij(korzen, 'panel-zadania', wynik.wynik.files.map((plik: DeveloperCoverage) =>
    `${plik.percent}% ${plik.path}`));
  oglos(NAGLOWEK, `Pokrycie prób: ${wynik.wynik.percent}%.`);
}

async function przegladBezpieczenstwa({ kanal, korzen, idOkna }: Otoczenie, panel: string):
Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Przegląd bezpieczeństwa',
    pola: [{
      klucz: 'zakres',
      etykieta: 'Zakres',
      wybor: [
        ['wszystko', 'Zależności, tajemnice, kod i licencje'],
        [ScanKind.Dependencies, 'Same zależności'],
        [ScanKind.Secrets, 'Same tajemnice'],
        [ScanKind.Code, 'Sam kod'],
        [ScanKind.Licenses, 'Same licencje'],
      ],
    }],
    wykonanie: 'Uruchom przegląd',
  });
  if (wartosci === null) return;
  const zakres = wartosci.zakres === 'wszystko' || wartosci.zakres === undefined
    ? [ScanKind.Dependencies, ScanKind.Secrets, ScanKind.Code, ScanKind.Licenses]
    : [wartosci.zakres as ScanKind];
  const wynik = await wywolaj(kanal, Command.DeveloperScanRun, {
    windowId: idOkna(),
    kinds: zakres,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił przeglądu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Przegląd w stanie ${wynik.wynik.scan.status}; `
    + `ustaleń ${wynik.wynik.scan.findingCount ?? 0}.`);
}

async function sprawdzNarzedzia({ kanal, korzen }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Narzędzia maszyny rdzenia',
    pola: [{
      klucz: 'programy',
      etykieta: 'Programy',
      wartosc: 'git go node python3',
      wymagane: true,
    }],
    wykonanie: 'Sprawdź narzędzia',
  });
  if (wartosci === null) return;
  const programy = (wartosci.programy ?? '').split(/\s+/).filter((nazwa) => nazwa !== '');
  const wynik = await wywolaj(kanal, Command.DeveloperToolchainCheck, { programs: programy });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił sprawdzenia narzędzi.', 'ostrzezenie');
    return;
  }
  wypelnij(korzen, 'panel-terminal', wynik.wynik.programs.map((program) =>
    `${program.present ? '✓' : '✗'} ${program.program} ${program.version ?? ''}`.trim()));
  const brakujace = wynik.wynik.programs.filter((program) => !program.present).length;
  if (brakujace > 0) {
    oglos(NAGLOWEK, `${brakujace} narzędzi nie stoi na maszynie rdzenia.`, 'ostrzezenie');
  }
}

async function zacznijPraceKrokowa({ kanal, korzen, idOkna }: Otoczenie, panel: string):
Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Praca krokowa',
    pola: [
      { klucz: 'program', etykieta: 'Program', podpowiedz: 'Puste bierze ustawienie okna' },
      { klucz: 'argumenty', etykieta: 'Argumenty' },
    ],
    wykonanie: 'Zacznij pracę krokową',
  });
  if (wartosci === null) return;
  const argumenty = (wartosci.argumenty ?? '').split(/\s+/).filter((slowo) => slowo !== '');
  const wynik = await wywolaj(kanal, Command.DeveloperDebugSessionStart, {
    windowId: idOkna(),
    ...(wartosci.program === '' ? {} : { program: wartosci.program }),
    ...(argumenty.length === 0 ? {} : { arguments: argumenty }),
    stopOnEntry: true,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił rozpoczęcia pracy krokowej.',
      'ostrzezenie');
    return;
  }
  PRZEBIEGI.set(idOkna(), wynik.wynik.session.id);
  oglos(NAGLOWEK, `Praca krokowa stoi na ${wynik.wynik.session.adapter}, `
    + `stan ${wynik.wynik.session.status}.`);
}

function przebieg(idOkna: string): string {
  const zapamietany = PRZEBIEGI.get(idOkna) ?? '';
  if (zapamietany === '') {
    oglos(NAGLOWEK, 'Okno nie prowadzi pracy krokowej — zacznij ją najpierw.', 'ostrzezenie');
  }
  return zapamietany;
}

async function krokDalej({ kanal, idOkna }: Otoczenie): Promise<void> {
  const cel = przebieg(idOkna());
  if (cel === '') return;
  const wynik = await wywolaj(kanal, Command.DeveloperDebugSessionControl, {
    sessionId: cel,
    step: DebugStepKind.StepOver,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił kroku.', 'ostrzezenie');
    return;
  }
  const stanowisko = wynik.wynik.session;
  oglos(NAGLOWEK, `Praca krokowa w stanie ${stanowisko.status}`
    + `${stanowisko.stoppedReason === undefined ? '' : `: ${stanowisko.stoppedReason}`}.`);
}

async function postawZatrzymanie({ kanal, korzen, idOkna }: Otoczenie, panel: string):
Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Zatrzymanie pracy krokowej',
    pola: [
      { klucz: 'sciezka', etykieta: 'Ścieżka pliku', wymagane: true },
      { klucz: 'wiersz', etykieta: 'Wiersz', wymagane: true },
      { klucz: 'warunek', etykieta: 'Warunek', podpowiedz: 'Puste zatrzymuje zawsze' },
    ],
    wykonanie: 'Postaw zatrzymanie',
  });
  if (wartosci === null) return;
  const wiersz = Number.parseInt(wartosci.wiersz ?? '', 10);
  if (!Number.isFinite(wiersz)) {
    oglos(NAGLOWEK, 'Wiersz musi być liczbą.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.DeveloperBreakpointSet, {
    windowId: idOkna(),
    path: wartosci.sciezka ?? '',
    line: wiersz,
    ...(wartosci.warunek === '' ? {} : { condition: wartosci.warunek }),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił postawienia zatrzymania.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Zatrzymań w oknie: ${wynik.wynik.breakpoints.length}.`);
}

/* Wyliczenie idzie w ramce stojącej najwyżej: okno nie prowadzi wyboru ramki,
   a zasięg podaje ich kolejność od miejsca zatrzymania. */
async function wyliczWyrazenie({ kanal, korzen, idOkna }: Otoczenie, panel: string):
Promise<void> {
  const cel = przebieg(idOkna());
  if (cel === '') return;
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Wyliczenie wyrażenia',
    pola: [{ klucz: 'wyrazenie', etykieta: 'Wyrażenie', wymagane: true }],
    wykonanie: 'Wylicz',
  });
  if (wartosci === null) return;
  const zasieg = await wywolaj(kanal, Command.DeveloperDebugScopeGet, { sessionId: cel });
  const ramka = zasieg.wynik?.frames[0]?.id ?? '';
  if (ramka === '') {
    oglos(NAGLOWEK, 'Praca krokowa nie stoi w żadnej ramce.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.DeveloperDebugEvaluate, {
    sessionId: cel,
    frameId: ramka,
    expression: wartosci.wyrazenie ?? '',
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wyliczenia.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `${wartosci.wyrazenie ?? ''} = ${wynik.wynik.value}`);
}

async function odczytajZasieg({ kanal, korzen, idOkna }: Otoczenie): Promise<void> {
  const cel = przebieg(idOkna());
  if (cel === '') return;
  const wynik = await wywolaj(kanal, Command.DeveloperDebugScopeGet, { sessionId: cel });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu zasięgu.', 'ostrzezenie');
    return;
  }
  wypelnij(korzen, 'panel-kolejka', [
    ...wynik.wynik.frames.map((ramka: DebugFrame) =>
      `ramka ${ramka.name} · ${ramka.path ?? ''}:${ramka.line ?? 0}`),
    ...wynik.wynik.variables.map((zmienna) => `${zmienna.name} = ${zmienna.value}`),
  ]);
  if (wynik.wynik.frames.length === 0) oglos(NAGLOWEK, 'Praca krokowa nie stoi w żadnej ramce.');
}

async function stanRepozytorium({ kanal, korzen, idOkna }: Otoczenie): Promise<void> {
  const wynik = await wywolaj(kanal, Command.DeveloperGitStatus, { windowId: idOkna() });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu stanu.', 'ostrzezenie');
    return;
  }
  const stan = wynik.wynik.status;
  if (!stan.isRepository) {
    oglos(NAGLOWEK, 'Katalog roboczy tego okna nie jest repozytorium.', 'ostrzezenie');
    return;
  }
  wypelnij(korzen, 'panel-git', stan.entries.map((pozycja: GitStatusEntry) =>
    `${pozycja.index}/${pozycja.worktree} ${pozycja.path}`));
  oglos(NAGLOWEK, `Gałąź ${stan.branch ?? 'bez nazwy'}: ${stan.entries.length} zmian, `
    + `${stan.hasConflicts ? 'są spory' : 'bez sporów'}.`);
}

async function roznicaZmian({ kanal, korzen, idOkna }: Otoczenie): Promise<void> {
  const wynik = await wywolaj(kanal, Command.DeveloperGitDiff, { windowId: idOkna() });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu różnicy.', 'ostrzezenie');
    return;
  }
  const fragmenty = wynik.wynik.hunks;
  wypelnij(korzen, 'panel-git', fragmenty.map((fragment) =>
    `${fragment.path} · ${fragment.header ?? ''}`));
  if (fragmenty.length === 0) oglos(NAGLOWEK, 'Katalog roboczy nie ma niezapisanych zmian.');
}

async function dziejeGalezi({ kanal, korzen, idOkna }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Dzieje gałęzi',
    pola: [
      { klucz: 'galaz', etykieta: 'Gałąź', podpowiedz: 'Puste bierze gałąź bieżącą' },
      { klucz: 'granica', etykieta: 'Najwyżej zapisów', wartosc: '50' },
    ],
    wykonanie: 'Odczytaj dzieje',
  });
  if (wartosci === null) return;
  const granica = Number.parseInt(wartosci.granica ?? '', 10);
  const wynik = await wywolaj(kanal, Command.DeveloperGitLog, {
    windowId: idOkna(),
    ...(wartosci.galaz === '' ? {} : { branch: wartosci.galaz }),
    limit: Number.isFinite(granica) ? granica : 50,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu dziejów.', 'ostrzezenie');
    return;
  }
  const zapisy = wynik.wynik.commits;
  wypelnij(korzen, 'panel-git', zapisy.map((zapis: GitCommit) =>
    `${zapis.shortId} ${zapis.author} · ${zapis.message.split('\n')[0] ?? ''}`));
  if (zapisy.length === 0) oglos(NAGLOWEK, 'Ta gałąź nie ma jeszcze zapisów.');
}

/* Zapis zmian bierze wszystkie pozycje stanu i zamyka je jednym opisem: okno
   nie prowadzi wyboru pozycji, więc dzielenie ich byłoby zgadywaniem. */
async function zapiszZmiany({ kanal, korzen, idOkna, odswiez }: Otoczenie, panel: string):
Promise<void> {
  const stan = await wywolaj(kanal, Command.DeveloperGitStatus, { windowId: idOkna() });
  const sciezki = (stan.wynik?.status.entries ?? []).map((pozycja) => pozycja.path);
  if (sciezki.length === 0) {
    oglos(NAGLOWEK, 'Nie ma zmian do zapisania.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Zapis zmian',
    opis: `Zapis obejmie wszystkie ${sciezki.length} zmienionych ścieżek.`,
    pola: [{ klucz: 'opis', etykieta: 'Opis zmiany', obszerne: true, wymagane: true }],
    wykonanie: 'Zapisz zmiany',
  });
  if (wartosci === null) return;
  const przygotowanie = await wywolaj(kanal, Command.DeveloperGitAction, {
    windowId: idOkna(),
    action: GitActionKind.Stage,
    paths: sciezki,
  });
  if (!przygotowanie.udany) {
    oglos(NAGLOWEK, przygotowanie.blad?.message ?? 'Rdzeń odmówił przygotowania zmian.',
      'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.DeveloperGitAction, {
    windowId: idOkna(),
    action: GitActionKind.Commit,
    message: wartosci.opis ?? '',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu zmian.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Zapisano ${sciezki.length} zmienionych ścieżek.`);
  odswiez();
}

async function odczytajSpor({ kanal, korzen, idOkna }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Spór w pliku',
    pola: [{ klucz: 'sciezka', etykieta: 'Ścieżka pliku', wymagane: true }],
    wykonanie: 'Odczytaj spór',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.DeveloperGitConflictGet, {
    windowId: idOkna(),
    path: wartosci.sciezka ?? '',
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu sporu.', 'ostrzezenie');
    return;
  }
  const spor = wynik.wynik.conflict;
  wypelnij(korzen, 'panel-edytor', [
    '— strona bieżąca —',
    ...spor.current.split('\n'),
    '— strona przychodząca —',
    ...spor.incoming.split('\n'),
  ]);
  oglos(NAGLOWEK, `Spór w pliku ${spor.path} stoi w edytorze.`);
}

async function rozstrzygnijSpor({ kanal, korzen, idOkna }: Otoczenie, panel: string):
Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Rozstrzygnięcie sporu',
    pola: [
      { klucz: 'sciezka', etykieta: 'Ścieżka pliku', wymagane: true },
      {
        klucz: 'strona',
        etykieta: 'Która strona zostaje',
        wybor: [
          [ConflictResolutionKind.TakeCurrent, 'Bieżąca'],
          [ConflictResolutionKind.TakeIncoming, 'Przychodząca'],
          [ConflictResolutionKind.TakeBoth, 'Obie'],
        ],
      },
    ],
    wykonanie: 'Rozstrzygnij spór',
    nieodwracalne: 'Treść pliku zostanie zastąpiona wskazaną stroną.',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.DeveloperGitConflictResolve, {
    windowId: idOkna(),
    path: wartosci.sciezka ?? '',
    resolution: (wartosci.strona ?? ConflictResolutionKind.TakeCurrent) as ConflictResolutionKind,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił rozstrzygnięcia sporu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, wynik.wynik.resolved
    ? `Spór rozstrzygnięty; zostaje ${wynik.wynik.remainingPaths.length} spornych plików.`
    : 'Rdzeń nie uznał sporu za rozstrzygnięty.',
  wynik.wynik.resolved ? 'informacja' : 'ostrzezenie');
}

async function przestawPojemnik({ kanal, korzen, idOkna, odswiez }: Otoczenie, panel: string):
Promise<void> {
  const wykaz = await wywolaj(kanal, Command.DeveloperContainerList, { windowId: idOkna() });
  if (wykaz.wynik?.engineAvailable === false) {
    oglos(NAGLOWEK, 'Rdzeń nie widzi na swojej maszynie niczego, co prowadzi pojemniki.',
      'ostrzezenie');
    return;
  }
  const pojemniki = wykaz.wynik?.containers ?? [];
  if (pojemniki.length === 0) {
    oglos(NAGLOWEK, 'Żaden pojemnik nie stoi dla tego okna.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Pojemnik',
    pola: [
      {
        klucz: 'pojemnik',
        etykieta: 'Pojemnik',
        wybor: pojemniki.map((p: ContainerInfo) => [p.id, `${p.name} · ${p.status}`] as const),
      },
      {
        klucz: 'czynnosc',
        etykieta: 'Czynność',
        wybor: [
          [ContainerActionKind.Start, 'Uruchom'],
          [ContainerActionKind.Stop, 'Zatrzymaj'],
          [ContainerActionKind.Restart, 'Uruchom ponownie'],
          [ContainerActionKind.Logs, 'Odczytaj dziennik'],
        ],
      },
    ],
    wykonanie: 'Wykonaj na pojemniku',
    nieodwracalne: 'Zatrzymanie przerywa usługę, którą pojemnik niesie.',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.DeveloperContainerAction, {
    containerId: wartosci.pojemnik ?? '',
    action: (wartosci.czynnosc ?? ContainerActionKind.Logs) as ContainerActionKind,
    tail: 100,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił czynności na pojemniku.', 'ostrzezenie');
    return;
  }
  if (wartosci.czynnosc === ContainerActionKind.Logs) {
    wypelnij(korzen, 'panel-terminal', (wynik.wynik.output ?? '').split('\n'));
  }
  oglos(NAGLOWEK, `Pojemnik ${wynik.wynik.container.name} w stanie `
    + `${wynik.wynik.container.status}.`);
  odswiez();
}

async function podniesZestaw({ kanal, korzen, idOkna }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Zestaw pojemników',
    pola: [
      { klucz: 'plik', etykieta: 'Plik zestawu', podpowiedz: 'compose.yaml', wymagane: true },
      { klucz: 'uslugi', etykieta: 'Usługi', podpowiedz: 'Puste bierze wszystkie' },
    ],
    wykonanie: 'Podnieś zestaw',
  });
  if (wartosci === null) return;
  const uslugi = (wartosci.uslugi ?? '').split(/\s+/).filter((nazwa) => nazwa !== '');
  const wynik = await wywolaj(kanal, Command.DeveloperComposeUp, {
    windowId: idOkna(),
    file: wartosci.plik ?? '',
    ...(uslugi.length === 0 ? {} : { services: uslugi }),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił podniesienia zestawu.', 'ostrzezenie');
    return;
  }
  wypelnij(korzen, 'panel-terminal', wynik.wynik.services.map((usluga: ContainerInfo) =>
    `${usluga.status} ${usluga.name}`));
  oglos(NAGLOWEK, `Zestaw podniesiony: ${wynik.wynik.services.length} usług.`);
}

/* Obraz nie idzie do składu: wysłanie go tam byłoby wystawieniem pracy na
   zewnątrz, o czym rozstrzyga Operator osobnym poleceniem. */
async function zbudujObraz({ kanal, korzen, idOkna }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Budowa obrazu',
    opis: 'Obraz zostaje na maszynie rdzenia; okno go nigdzie nie wysyła.',
    pola: [
      { klucz: 'plik', etykieta: 'Plik opisu', podpowiedz: 'Dockerfile', wymagane: true },
      { klucz: 'oznaczenie', etykieta: 'Oznaczenie obrazu', podpowiedz: 'sklep:1.0', wymagane: true },
    ],
    wykonanie: 'Zbuduj obraz',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.DeveloperImageBuild, {
    windowId: idOkna(),
    dockerfile: wartosci.plik ?? '',
    tag: wartosci.oznaczenie ?? '',
    push: false,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zbudowania obrazu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Obraz zbudowany pod oznaczeniem ${wynik.wynik.imageId}.`);
}

async function zapiszPolaczenie({ kanal, korzen, idOkna, odswiez }: Otoczenie, panel: string):
Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Połączenie z bazą',
    pola: [
      { klucz: 'nazwa', etykieta: 'Nazwa połączenia', wymagane: true },
      {
        klucz: 'rodzajBazy',
        etykieta: 'Rodzaj bazy',
        wybor: [
          [DataEngine.Sqlite, 'SQLite'],
          [DataEngine.Postgres, 'PostgreSQL'],
          [DataEngine.Mysql, 'MySQL'],
        ],
      },
      { klucz: 'baza', etykieta: 'Baza', podpowiedz: 'dane.db albo nazwa bazy', wymagane: true },
      { klucz: 'host', etykieta: 'Host', podpowiedz: 'Puste dla bazy w pliku' },
      {
        klucz: 'zapis',
        etykieta: 'Prawo zapisu',
        wybor: [['odczyt', 'Tylko odczyt'], ['zapis', 'Odczyt i zapis']],
      },
    ],
    wykonanie: 'Zapisz połączenie',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.DeveloperDataConnectionSet, {
    windowId: idOkna(),
    name: wartosci.nazwa ?? '',
    engine: (wartosci.rodzajBazy ?? DataEngine.Sqlite) as DataEngine,
    database: wartosci.baza ?? '',
    ...(wartosci.host === '' ? {} : { host: wartosci.host }),
    readOnly: wartosci.zapis !== 'zapis',
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu połączenia.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Połączenie „${wynik.wynik.connection.name}" zapisane`
    + `${wynik.wynik.connection.readOnly ? ' tylko do odczytu' : ''}.`);
  odswiez();
}

/* Zapytanie, schemat i przeniesienia idą przez połączenie stojące w wykazie
   najwyżej: okno nie prowadzi wyboru połączenia. */
async function pierwszePolaczenie(kanal: Kanal, idOkna: string): Promise<string> {
  const wykaz = await wywolaj(kanal, Command.DeveloperDataConnectionList, { windowId: idOkna });
  const cel = wykaz.wynik?.connections[0]?.id ?? '';
  if (cel === '') {
    oglos(NAGLOWEK, 'To okno nie ma jeszcze zapisanego połączenia z bazą.', 'ostrzezenie');
  }
  return cel;
}

async function wykonajZapytanie({ kanal, korzen, idOkna }: Otoczenie, panel: string):
Promise<void> {
  const cel = await pierwszePolaczenie(kanal, idOkna());
  if (cel === '') return;
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Zapytanie do bazy',
    pola: [
      { klucz: 'zapytanie', etykieta: 'Treść zapytania', obszerne: true, wymagane: true },
      { klucz: 'granica', etykieta: 'Najwyżej wierszy', wartosc: '200' },
    ],
    wykonanie: 'Wykonaj zapytanie',
  });
  if (wartosci === null) return;
  const granica = Number.parseInt(wartosci.granica ?? '', 10);
  const wynik = await wywolaj(kanal, Command.DeveloperDataQueryRun, {
    connectionId: cel,
    sql: wartosci.zapytanie ?? '',
    limit: Number.isFinite(granica) ? granica : 200,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wykonania zapytania.', 'ostrzezenie');
    return;
  }
  const odpowiedz = wynik.wynik.result;
  wypelnij(korzen, 'panel-zadania', [odpowiedz.columns.join(' · ')]);
  oglos(NAGLOWEK, odpowiedz.truncated
    ? `Zapytanie oddało ${odpowiedz.rowCount} wierszy, przycięte granicą odczytu.`
    : `Zapytanie oddało ${odpowiedz.rowCount} wierszy w ${odpowiedz.durationMs} ms.`,
  odpowiedz.truncated ? 'ostrzezenie' : 'informacja');
}

async function odczytajSchemat({ kanal, korzen, idOkna }: Otoczenie): Promise<void> {
  const cel = await pierwszePolaczenie(kanal, idOkna());
  if (cel === '') return;
  const wynik = await wywolaj(kanal, Command.DeveloperDataSchemaGet, { connectionId: cel });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu schematu.', 'ostrzezenie');
    return;
  }
  const wezly = wynik.wynik.nodes;
  wypelnij(korzen, 'panel-zadania', wezly.map((wezel: DataSchemaNode) =>
    `${wezel.kind} ${wezel.name}`));
  if (wezly.length === 0) oglos(NAGLOWEK, 'Ta baza nie ma jeszcze żadnych tablic.');
}

async function wykonajPrzeniesienia({ kanal, korzen, idOkna }: Otoczenie, panel: string):
Promise<void> {
  const cel = await pierwszePolaczenie(kanal, idOkna());
  if (cel === '') return;
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Przeniesienia bazy',
    pola: [{
      klucz: 'kierunek',
      etykieta: 'Kierunek',
      wybor: [['up', 'Naprzód'], ['down', 'Wstecz']],
    }],
    wykonanie: 'Wykonaj przeniesienia',
    nieodwracalne: 'Przeniesienia zmieniają układ tablic bazy.',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.DeveloperDataMigrationRun, {
    connectionId: cel,
    windowId: idOkna(),
    direction: wartosci.kierunek ?? 'up',
    dryRun: false,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił przeniesień.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Wykonano ${wynik.wynik.applied.length} przeniesień; `
    + `czeka ${wynik.wynik.pending.length}.`);
}

async function wyslijZadanie({ kanal, korzen, idOkna }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Żądanie sieciowe',
    pola: [
      {
        klucz: 'sposob',
        etykieta: 'Sposób',
        wybor: [['GET', 'GET'], ['POST', 'POST'], ['PUT', 'PUT'], ['PATCH', 'PATCH'],
          ['DELETE', 'DELETE']],
      },
      { klucz: 'adres', etykieta: 'Adres', podpowiedz: 'https://przyklad.pl', wymagane: true },
      { klucz: 'tresc', etykieta: 'Treść żądania', obszerne: true },
    ],
    wykonanie: 'Wyślij żądanie',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.DeveloperApiRequest, {
    windowId: idOkna(),
    method: wartosci.sposob ?? 'GET',
    url: wartosci.adres ?? '',
    ...(wartosci.tresc === '' ? {} : { body: wartosci.tresc }),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wysłania żądania.', 'ostrzezenie');
    return;
  }
  const odpowiedz = wynik.wynik.response;
  wypelnij(korzen, 'panel-artefakty', (odpowiedz.body ?? '').split('\n'));
  oglos(NAGLOWEK, `Odpowiedź ${odpowiedz.status} w ${odpowiedz.durationMs} ms.`);
}

async function zapiszZbior({ kanal, korzen, idOkna, odswiez }: Otoczenie, panel: string):
Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Zbiór żądań',
    pola: [{ klucz: 'nazwa', etykieta: 'Nazwa zbioru', wymagane: true }],
    wykonanie: 'Zapisz zbiór',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.DeveloperApiCollectionSave, {
    windowId: idOkna(),
    name: wartosci.nazwa ?? '',
    requests: [],
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu zbioru.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Zbiór „${wynik.wynik.collection.name}" zapisany.`);
  odswiez();
}

async function wciagnijOpenapi({ kanal, korzen, idOkna, odswiez }: Otoczenie, panel: string):
Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Opis OpenAPI',
    pola: [
      { klucz: 'wskazanie', etykieta: 'Ścieżka albo adres opisu', wymagane: true },
      { klucz: 'nazwa', etykieta: 'Nazwa zbioru', podpowiedz: 'Puste bierze nazwę z opisu' },
    ],
    wykonanie: 'Wciągnij opis',
  });
  if (wartosci === null) return;
  const wskazanie = wartosci.wskazanie ?? '';
  const przezSiec = wskazanie.startsWith('http://') || wskazanie.startsWith('https://');
  const wynik = await wywolaj(kanal, Command.DeveloperApiOpenapiImport, {
    windowId: idOkna(),
    ...(przezSiec ? { url: wskazanie } : { path: wskazanie }),
    ...(wartosci.nazwa === '' ? {} : { collectionName: wartosci.nazwa }),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wciągnięcia opisu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Zbiór „${wynik.wynik.collection.name}" wciągnięty: `
    + `${wynik.wynik.requestCount} żądań.`);
  odswiez();
}

async function zbadajObciazenie({ kanal, korzen, idOkna }: Otoczenie, panel: string):
Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Badanie obciążenia',
    pola: [
      { klucz: 'adres', etykieta: 'Adres', wymagane: true },
      { klucz: 'polaczenia', etykieta: 'Równoczesnych połączeń', wartosc: '10' },
      { klucz: 'sekundy', etykieta: 'Czas w sekundach', wartosc: '10' },
    ],
    wykonanie: 'Zbadaj obciążenie',
    nieodwracalne: 'Badanie wyśle wskazanemu serwerowi wiele żądań naraz.',
  });
  if (wartosci === null) return;
  const polaczenia = Number.parseInt(wartosci.polaczenia ?? '', 10);
  const sekundy = Number.parseInt(wartosci.sekundy ?? '', 10);
  const wynik = await wywolaj(kanal, Command.DeveloperApiLoadRun, {
    windowId: idOkna(),
    url: wartosci.adres ?? '',
    connections: Number.isFinite(polaczenia) ? polaczenia : 10,
    durationSeconds: Number.isFinite(sekundy) ? sekundy : 10,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił badania obciążenia.', 'ostrzezenie');
    return;
  }
  const bieg = wynik.wynik.run;
  oglos(NAGLOWEK, `${bieg.requestsTotal} żądań, ${bieg.requestsPerSecond} na sekundę, `
    + `poza zakresem 2xx: ${bieg.non2xx}, błędów ${bieg.errors}.`);
}
