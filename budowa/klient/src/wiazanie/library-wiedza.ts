// Panel podglądu w oknie Library: wyszukiwanie znaczeniowe po zbiorach oraz
// obróbka pojedynczego pliku — nagrania, dokumentu i archiwum.
import { Command, KnowledgeScope, MediaOperationKind } from '../../../shared/contract.ts';
import type { KnowledgeHit } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { dolozCzynnosciPanelu, zapytajWSzufladzie } from './czynnosci-okna.ts';
import { oglos } from './ogloszenie.ts';

const NAGLOWEK = 'Wiedza i pliki';
const PANEL = 'panel-preview';

const ZBIORY: ReadonlyArray<readonly [string, string]> = [
  [KnowledgeScope.All, 'Wszystko'],
  [KnowledgeScope.Library, 'Biblioteka'],
  [KnowledgeScope.History, 'Historia rozmów'],
  [KnowledgeScope.Workspace, 'Przestrzeń pracy'],
  [KnowledgeScope.ProjectMemory, 'Pamięć projektu'],
  [KnowledgeScope.AssistantActivity, 'Czynności asystenta'],
];

const OBROBKI: ReadonlyArray<readonly [string, string]> = [
  [MediaOperationKind.Convert, 'Przełóż na inną postać'],
  [MediaOperationKind.Trim, 'Wytnij odcinek'],
  [MediaOperationKind.ExtractAudio, 'Wyjmij ścieżkę dźwięku'],
  [MediaOperationKind.Resize, 'Przeskaluj obraz'],
  [MediaOperationKind.Frame, 'Wyjmij klatkę'],
];

interface Otoczenie {
  kanal: Kanal;
  korzen: Element;
  idOkna: () => string;
}

export function zwiazWiedzeBiblioteki(
  kanal: Kanal,
  korzen: Element,
  idOkna: () => string,
  dolacz: (zdejmij: () => void) => void,
  przy: AddEventListenerOptions,
): void {
  const otoczenie: Otoczenie = { kanal, korzen, idOkna };
  const zdejmij = dolozCzynnosciPanelu(korzen, PANEL, 'Wiedza i pliki', [
    {
      naglowek: 'Wyszukiwanie znaczeniowe',
      pozycje: [
        { kod: 'szukaj', nazwa: 'Szukaj w zbiorach…' },
        { kod: 'szukaj-obrazy', nazwa: 'Szukaj wśród obrazów…' },
        { kod: 'przelicz', nazwa: 'Przelicz wskaźnik zbioru…' },
      ],
    },
    {
      naglowek: 'Nagrania',
      pozycje: [
        { kod: 'obejrzyj', nazwa: 'Obejrzyj nagranie…' },
        { kod: 'przerob', nazwa: 'Przerób nagranie…' },
      ],
    },
    {
      naglowek: 'Dokumenty i archiwa',
      pozycje: [
        { kod: 'przeloz', nazwa: 'Przełóż dokument…' },
        { kod: 'wyjmij-tekst', nazwa: 'Wyjmij tekst z dokumentu…' },
        { kod: 'spakuj', nazwa: 'Spakuj do archiwum…' },
        { kod: 'rozpakuj', nazwa: 'Rozpakuj archiwum…' },
      ],
    },
  ], (kod) => {
    void wykonaj(otoczenie, kod);
  }, przy);
  if (zdejmij !== null) dolacz(zdejmij);
}

function wypelnij(korzen: Element, wiersze: string[]): void {
  const cialo = korzen.querySelector(`#${PANEL} .sta-okno-tresc`);
  if (cialo === null) return;
  cialo.replaceChildren(...wiersze.map((tresc) => {
    const wiersz = cialo.ownerDocument.createElement('div');
    wiersz.className = 'dn-wykaz-modulu-poz';
    wiersz.textContent = tresc;
    return wiersz;
  }));
}

async function wykonaj(otoczenie: Otoczenie, kod: string): Promise<void> {
  if (kod === 'szukaj') return szukajWZbiorach(otoczenie);
  if (kod === 'szukaj-obrazy') return szukajObrazow(otoczenie);
  if (kod === 'przelicz') return przeliczWskaznik(otoczenie);
  if (kod === 'obejrzyj') return obejrzyjNagranie(otoczenie);
  if (kod === 'przerob') return przerobNagranie(otoczenie);
  if (kod === 'przeloz') return przelozDokument(otoczenie);
  if (kod === 'wyjmij-tekst') return wyjmijTekst(otoczenie);
  if (kod === 'spakuj') return spakujArchiwum(otoczenie);
  if (kod === 'rozpakuj') return rozpakujArchiwum(otoczenie);
  /* Czynność spoza obsłużonych odmawia zamiast milczeć. */
  oglos(NAGLOWEK, `Czynność „${kod}" nie jest prowadzona przez to okno.`, 'ostrzezenie');
}

function opiszTrafienie(trafienie: KnowledgeHit): string {
  const ocena = trafienie.score === undefined ? '' : ` · ${trafienie.score.toFixed(2)}`;
  return `${trafienie.source}${ocena} · ${trafienie.text.slice(0, 60)}`;
}

/* Wyszukiwanie idzie po znaczeniu, nie po dosłownym słowie: rdzeń porównuje
   zapytanie ze wskaźnikiem zbioru, więc pusty wskaźnik nie znajdzie niczego. */
async function szukajWZbiorach(otoczenie: Otoczenie): Promise<void> {
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, PANEL, {
    tytul: 'Wyszukiwanie znaczeniowe',
    pola: [
      { klucz: 'zapytanie', etykieta: 'Czego szukasz', obszerne: true, wymagane: true },
      { klucz: 'zbior', etykieta: 'Gdzie szukać', wybor: ZBIORY },
    ],
    wykonanie: 'Szukaj',
  });
  if (wartosci === null) return;
  oglos(NAGLOWEK, 'Szukam w zbiorach…');
  const wynik = await wywolaj(otoczenie.kanal, Command.KnowledgeSearch, {
    query: wartosci.zapytanie ?? '',
    scope: (wartosci.zbior ?? KnowledgeScope.All) as KnowledgeScope,
    limit: 20,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wyszukiwania.', 'ostrzezenie');
    return;
  }
  const trafienia = wynik.wynik.results;
  wypelnij(otoczenie.korzen, trafienia.length === 0
    ? ['Nic nie odpowiada temu zapytaniu.']
    : trafienia.map(opiszTrafienie));
  oglos(NAGLOWEK, `Trafień: ${String(wynik.wynik.total)}.`);
}

async function szukajObrazow(otoczenie: Otoczenie): Promise<void> {
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, PANEL, {
    tytul: 'Wyszukiwanie wśród obrazów',
    opis: 'Rdzeń porównuje opis z treścią obrazów, nie z nazwą pliku.',
    pola: [{ klucz: 'zapytanie', etykieta: 'Co ma być na obrazie', obszerne: true, wymagane: true }],
    wykonanie: 'Szukaj',
  });
  if (wartosci === null) return;
  oglos(NAGLOWEK, 'Przeglądam obrazy…');
  const wynik = await wywolaj(otoczenie.kanal, Command.KnowledgeImageSearch, {
    query: wartosci.zapytanie ?? '',
    limit: 20,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił przeglądu obrazów.', 'ostrzezenie');
    return;
  }
  const trafienia = wynik.wynik.results;
  wypelnij(otoczenie.korzen, trafienia.length === 0
    ? ['Żaden obraz nie odpowiada temu opisowi.']
    : trafienia.map((trafienie) => JSON.stringify(trafienie).slice(0, 90)));
  oglos(NAGLOWEK, `Obrazów pasujących: ${String(wynik.wynik.total)}.`);
}

async function przeliczWskaznik(otoczenie: Otoczenie): Promise<void> {
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, PANEL, {
    tytul: 'Przeliczenie wskaźnika',
    opis: 'Przeliczenie od nowa liczy cały zbiór, nie tylko to, co doszło.',
    pola: [
      { klucz: 'zbior', etykieta: 'Zbiór', wybor: ZBIORY },
      { klucz: 'od-nowa', etykieta: 'Licz od nowa', wybor: [['nie', 'Nie'], ['tak', 'Tak']] },
    ],
    wykonanie: 'Przelicz',
  });
  if (wartosci === null) return;
  oglos(NAGLOWEK, 'Przeliczam wskaźnik zbioru…');
  const wynik = await wywolaj(otoczenie.kanal, Command.KnowledgeIndex, {
    scope: (wartosci.zbior ?? KnowledgeScope.Library) as KnowledgeScope,
    rebuild: wartosci['od-nowa'] === 'tak',
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił przeliczenia wskaźnika.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Policzono ${String(wynik.wynik.indexed)} z ${String(wynik.wynik.total)}`
    + (wynik.wynik.model === undefined ? '.' : ` · ${wynik.wynik.model}`));
}

async function wskazPlik(
  otoczenie: Otoczenie,
  tytul: string,
  wykonanie: string,
  dodatkowe: Parameters<typeof zapytajWSzufladzie>[2]['pola'] = [],
  opis?: string,
): Promise<Record<string, string> | null> {
  return zapytajWSzufladzie(otoczenie.korzen, PANEL, {
    tytul,
    ...(opis === undefined ? {} : { opis }),
    pola: [{ klucz: 'sciezka', etykieta: 'Ścieżka pliku', wymagane: true }, ...dodatkowe],
    wykonanie,
  });
}

async function obejrzyjNagranie(otoczenie: Otoczenie): Promise<void> {
  const wartosci = await wskazPlik(otoczenie, 'Obejrzenie nagrania', 'Obejrzyj');
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.MediaInspect, {
    sourcePath: wartosci.sciezka ?? '',
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił obejrzenia nagrania.', 'ostrzezenie');
    return;
  }
  const opis = wynik.wynik;
  wypelnij(otoczenie.korzen, [
    `postać: ${opis.format}`,
    `czas: ${String(opis.durationMs)} ms`,
    `wielkość: ${String(opis.sizeBytes)} B`,
    `ścieżki: ${opis.streams.slice(0, 120)}`,
  ]);
  oglos(NAGLOWEK, `Nagranie ${opis.format} · ${String(opis.durationMs)} ms.`);
}

async function przerobNagranie(otoczenie: Otoczenie): Promise<void> {
  const wartosci = await wskazPlik(otoczenie, 'Obróbka nagrania', 'Przerób', [
    { klucz: 'obrobka', etykieta: 'Co zrobić', wybor: OBROBKI },
    { klucz: 'postac', etykieta: 'Postać wyniku' },
  ]);
  if (wartosci === null) return;
  oglos(NAGLOWEK, 'Zlecam obróbkę nagrania…');
  const wynik = await wywolaj(otoczenie.kanal, Command.MediaTranscode, {
    sourcePath: wartosci.sciezka ?? '',
    operation: (wartosci.obrobka ?? MediaOperationKind.Convert) as MediaOperationKind,
    ...(wartosci.postac === '' ? {} : { format: wartosci.postac }),
    ...(otoczenie.idOkna() === '' ? {} : { windowId: otoczenie.idOkna() }),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił obróbki nagrania.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Wynik obróbki stoi w katalogu jako „${wynik.wynik.asset.name}".`);
}

async function przelozDokument(otoczenie: Otoczenie): Promise<void> {
  const wartosci = await wskazPlik(otoczenie, 'Przełożenie dokumentu', 'Przełóż', [
    { klucz: 'postac', etykieta: 'Postać wyniku', podpowiedz: 'pdf', wymagane: true },
  ]);
  if (wartosci === null) return;
  oglos(NAGLOWEK, 'Zlecam przełożenie dokumentu…');
  const wynik = await wywolaj(otoczenie.kanal, Command.DocumentConvert, {
    sourcePath: wartosci.sciezka ?? '',
    toFormat: wartosci.postac ?? '',
    ...(otoczenie.idOkna() === '' ? {} : { windowId: otoczenie.idOkna() }),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił przełożenia dokumentu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Dokument przełożony: „${wynik.wynik.asset.name}"`
    + ` · ${String(wynik.wynik.sizeBytes)} B.`);
}

/* Rozpoznanie pisma rdzeń włącza sam, gdy dokument nie niesie warstwy tekstu;
   wymuszenie każe mu przeczytać obraz nawet wtedy, gdy tekst już jest. */
async function wyjmijTekst(otoczenie: Otoczenie): Promise<void> {
  const wartosci = await wskazPlik(otoczenie, 'Wyjęcie tekstu', 'Wyjmij tekst', [
    {
      klucz: 'rozpoznanie',
      etykieta: 'Wymuś rozpoznanie pisma',
      wybor: [['nie', 'Nie'], ['tak', 'Tak']],
    },
  ]);
  if (wartosci === null) return;
  oglos(NAGLOWEK, 'Czytam dokument…');
  const wynik = await wywolaj(otoczenie.kanal, Command.DocumentTextExtract, {
    sourcePath: wartosci.sciezka ?? '',
    forceOcr: wartosci.rozpoznanie === 'tak',
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu dokumentu.', 'ostrzezenie');
    return;
  }
  const odczyt = wynik.wynik;
  wypelnij(otoczenie.korzen, [
    `stron: ${odczyt.pages === undefined ? 'bez wskazania' : String(odczyt.pages)}`,
    `rozpoznanie pisma: ${odczyt.usedOcr ? 'użyte' : 'niepotrzebne'}`,
    odczyt.text.slice(0, 200),
  ]);
  oglos(NAGLOWEK, `Odczytano ${String(odczyt.text.length)} znaków.`);
}

async function spakujArchiwum(otoczenie: Otoczenie): Promise<void> {
  const wartosci = await wskazPlik(otoczenie, 'Spakowanie do archiwum', 'Spakuj', [
    { klucz: 'postac', etykieta: 'Postać archiwum', podpowiedz: 'zip' },
    { klucz: 'nazwa', etykieta: 'Nazwa archiwum' },
  ]);
  if (wartosci === null) return;
  oglos(NAGLOWEK, 'Pakuję…');
  const wynik = await wywolaj(otoczenie.kanal, Command.ArchivePack, {
    sourcePath: wartosci.sciezka ?? '',
    ...(wartosci.postac === '' ? {} : { format: wartosci.postac }),
    ...(wartosci.nazwa === '' ? {} : { name: wartosci.nazwa }),
    ...(otoczenie.idOkna() === '' ? {} : { windowId: otoczenie.idOkna() }),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił spakowania.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Archiwum „${wynik.wynik.asset.name}" niesie`
    + ` ${String(wynik.wynik.entries)} pozycji.`);
}

async function rozpakujArchiwum(otoczenie: Otoczenie): Promise<void> {
  const wartosci = await wskazPlik(otoczenie, 'Rozpakowanie archiwum', 'Rozpakuj', [
    { klucz: 'cel', etykieta: 'Katalog docelowy' },
  ], 'Pusty katalog docelowy zostawia wybór rdzeniowi.');
  if (wartosci === null) return;
  oglos(NAGLOWEK, 'Rozpakowuję…');
  const wynik = await wywolaj(otoczenie.kanal, Command.ArchiveUnpack, {
    sourcePath: wartosci.sciezka ?? '',
    ...(wartosci.cel === '' ? {} : { targetPath: wartosci.cel }),
    ...(otoczenie.idOkna() === '' ? {} : { windowId: otoczenie.idOkna() }),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił rozpakowania.', 'ostrzezenie');
    return;
  }
  const sciezki = wynik.wynik.paths ?? [];
  wypelnij(otoczenie.korzen, sciezki.length === 0
    ? [`rozpakowano ${String(wynik.wynik.entries)} pozycji`]
    : sciezki.slice(0, 20));
  oglos(NAGLOWEK, `Rozpakowano ${String(wynik.wynik.entries)} pozycji.`);
}
