// Panel „Library Explorer" i panel „Pliki" wobec repozytorium rdzenia.
// Zakładki paska prototypu rozdzielają źródło wykazu: trzy pierwsze biorą
// wykaz plików, czwarta dziennik czynności, bo tylko on niesie oś czasu.
import {
  Command,
  LibraryDuplicateKind,
  type LibraryAuditEntry,
  type LibraryFile,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import {
  dzien,
  panel,
  plakietka,
  pozycja,
  pustka,
  pustkaWiersza,
  rozmiar,
  wartosc,
  zdejmijDzieci,
  type Kontekst,
} from './library-wspolne.ts';

const KOLUMN = 6;
const GRANICA_WYKAZU = 200;

export function zwiazEksplorator(
  kanal: Kanal,
  korzen: ParentNode,
  kontekst: Kontekst,
  przy: AddEventListenerOptions,
): void {
  const tresc = panel(korzen, 'panel-explorer');
  const cialoMoze = tresc?.querySelector('.lb-tabela tbody') ?? null;
  const szukajMoze = tresc?.querySelector<HTMLInputElement>('input[type="search"]') ?? null;
  if (cialoMoze === null || szukajMoze === null) return;
  const cialoTabeli = cialoMoze;
  const szukaj = szukajMoze;
  const okruszki = tresc?.querySelector('.pt-okruszki') ?? null;
  const zakladki = [...(tresc?.querySelectorAll<HTMLElement>('.lb-zakladki [role="tab"]') ?? [])];
  const importuj = tresc?.querySelector('.lb-pasek .dn-btn') ?? null;
  if (cialoTabeli === null || !(szukaj instanceof HTMLInputElement)) return;

  const dokument = cialoTabeli.ownerDocument;
  const pliki = new Map<HTMLElement, LibraryFile>();
  let numerZadania = 0;

  const nazwaZakladki = (): string =>
    zakladki.find((karta) => karta.getAttribute('aria-selected') === 'true')?.textContent?.trim() ??
    'Lista';

  async function odswiez(): Promise<void> {
    const numer = (numerZadania += 1);
    const fraza = szukaj.value.trim();
    const zakladka = nazwaZakladki();
    if (zakladka === 'Oś czasu') {
      const dziennik = await pobierzDziennik(kanal);
      if (numer !== numerZadania) return;
      naniesDziennik(cialoTabeli, dziennik);
      return;
    }
    const wykaz = await pobierzPliki(kanal, fraza, kontekst.kolekcja());
    if (numer !== numerZadania) return;
    const widoczne =
      zakladka === 'Galeria'
        ? wykaz.filter((plik) => plik.mimeType?.startsWith('image/') === true)
        : wykaz;
    pliki.clear();
    if (widoczne.length === 0) {
      pustkaWiersza(cialoTabeli, KOLUMN, 'Rdzeń nie zwrócił plików dla tego zakresu.');
      return;
    }
    zdejmijDzieci(cialoTabeli);
    for (const plik of widoczne) {
      const wiersz = zbudujWiersz(dokument, plik);
      pliki.set(wiersz, plik);
      cialoTabeli.appendChild(wiersz);
    }
  }

  szukaj.addEventListener(
    'input',
    () => {
      void odswiez();
    },
    przy,
  );

  for (const karta of zakladki) {
    karta.addEventListener(
      'click',
      () => {
        for (const inna of zakladki) inna.setAttribute('aria-selected', String(inna === karta));
        void odswiez();
      },
      przy,
    );
  }

  cialoTabeli.addEventListener(
    'click',
    (zdarzenie) => {
      const cel = zdarzenie.target;
      if (!(cel instanceof Element)) return;
      const wiersz = cel.closest('tr');
      if (!(wiersz instanceof HTMLElement)) return;
      const plik = pliki.get(wiersz);
      if (plik === undefined) return;
      for (const inny of pliki.keys()) inny.removeAttribute('aria-selected');
      wiersz.setAttribute('aria-selected', 'true');
      kontekst.ustawPlik(plik);
    },
    przy,
  );

  cialoTabeli.addEventListener(
    'dragstart',
    (zdarzenie) => {
      const cel = zdarzenie.target;
      if (!(cel instanceof Element)) return;
      const wiersz = cel.closest('tr');
      const plik = wiersz instanceof HTMLElement ? pliki.get(wiersz) : undefined;
      if (plik === undefined) return;
      (zdarzenie as DragEvent).dataTransfer?.setData('text/dn-library-plik', plik.id);
    },
    przy,
  );

  for (const wiersz of cialoTabeli.querySelectorAll('tr')) wiersz.setAttribute('draggable', 'true');

  if (okruszki !== null) {
    okruszki.addEventListener(
      'dragover',
      (zdarzenie) => {
        if ((zdarzenie as DragEvent).dataTransfer?.types.includes('text/dn-library-plik') === true) {
          zdarzenie.preventDefault();
        }
      },
      przy,
    );
    okruszki.addEventListener(
      'drop',
      (zdarzenie) => {
        zdarzenie.preventDefault();
        const cel = zdarzenie.target;
        const idPliku = (zdarzenie as DragEvent).dataTransfer?.getData('text/dn-library-plik') ?? '';
        const sciezka = cel instanceof Element ? (cel.textContent?.trim() ?? '') : '';
        if (idPliku === '' || sciezka === '') return;
        void przeniesPlik(kanal, idPliku, sciezka).then(odswiez);
      },
      przy,
    );
  }

  if (importuj !== null) {
    importuj.addEventListener(
      'click',
      () => {
        wybierzZDysku(dokument, (nazwa, typ, tresc64) => {
          void wgrajPlik(kanal, nazwa, typ, tresc64, kontekst.kolekcja()).then(odswiez);
        });
      },
      przy,
    );
  }

  zwiazMenuSesji(kanal, korzen, kontekst, przy, odswiez);
  zwiazPanelPlikow(kanal, korzen, przy, odswiez);
  kontekst.przyKolekcji(() => {
    void odswiez();
  });
  kontekst.zglosOdswiezenie(() => {
    void odswiez();
  });
  void odswiezOkruszki(kanal, okruszki, kontekst);
  /* Bez pierwszego odczytu wykaz stoi przy wierszach z prototypu do chwili,
     w której Operator ruszy szukaniem albo zakładką. */
  void odswiez();
}

function zwiazMenuSesji(
  kanal: Kanal,
  korzen: ParentNode,
  kontekst: Kontekst,
  przy: AddEventListenerOptions,
  odswiez: () => Promise<void>,
): void {
  for (const pozycjaMenu of korzen.querySelectorAll<HTMLElement>(
    '#menu-sesja .sta-menu-poz[role="menuitem"]',
  )) {
    const podpis = pozycjaMenu.textContent?.trim() ?? '';
    if (podpis.startsWith('Znajdź duplikaty')) {
      pozycjaMenu.addEventListener(
        'click',
        () => {
          void skanujDuplikaty(kanal, kontekst.kolekcja());
        },
        przy,
      );
      continue;
    }
    if (podpis.startsWith('Archiwizuj kolekcję')) {
      pozycjaMenu.addEventListener(
        'click',
        () => {
          void archiwizuj(kanal, kontekst).then(odswiez);
        },
        przy,
      );
    }
  }
}

function zwiazPanelPlikow(
  kanal: Kanal,
  korzen: ParentNode,
  przy: AddEventListenerOptions,
  odswiezWykaz: () => Promise<void>,
): void {
  const trescMozeMoze = panel(korzen, 'panel-pliki');
  if (trescMozeMoze === null) return;
  const trescMoze = trescMozeMoze;
  const tresc = trescMoze;
  const archiwalne = new Map<HTMLElement, string>();

  async function odswiez(): Promise<void> {
    const odpowiedz = await wywolaj(kanal, Command.LibraryFileList, { limit: GRANICA_WYKAZU });
    const wynik = wartosc(odpowiedz);
    const zarchiwizowane = (wynik?.files ?? []).filter(
      (plik) => plik.path?.startsWith('archiwum/') === true,
    );
    archiwalne.clear();
    if (zarchiwizowane.length === 0) {
      pustka(tresc, 'Rdzeń nie zwrócił plików w archiwum.');
      return;
    }
    zdejmijDzieci(tresc);
    for (const plik of zarchiwizowane) {
      const wiersz = pozycja(tresc.ownerDocument, plik.name, rozmiar(plik.sizeBytes));
      archiwalne.set(wiersz, plik.id);
      wiersz.tabIndex = 0;
      tresc.appendChild(wiersz);
    }
  }

  tresc.addEventListener(
    'click',
    (zdarzenie) => {
      const cel = zdarzenie.target;
      if (!(cel instanceof Element)) return;
      const wiersz = cel.closest('.dn-wykaz-modulu-poz');
      const idPliku = wiersz instanceof HTMLElement ? archiwalne.get(wiersz) : undefined;
      if (idPliku === undefined) return;
      void przywroc(kanal, idPliku)
        .then(odswiez)
        .then(odswiezWykaz);
    },
    przy,
  );

  // Kontrakt żąda jawnego potwierdzenia usunięcia; gest niesie je klawiszem
  // Shift, bo prototyp nie ma w tym panelu ani okna pytania, ani przycisku.
  tresc.addEventListener(
    'keydown',
    (zdarzenie) => {
      if (zdarzenie.key !== 'Delete' || !zdarzenie.shiftKey) return;
      const cel = zdarzenie.target;
      if (!(cel instanceof HTMLElement)) return;
      const idPliku = archiwalne.get(cel);
      if (idPliku === undefined) return;
      zdarzenie.preventDefault();
      void usun(kanal, idPliku)
        .then(odswiez)
        .then(odswiezWykaz);
    },
    przy,
  );

  void odswiez();
}

async function odswiezOkruszki(
  kanal: Kanal,
  okruszki: Element | null,
  kontekst: Kontekst,
): Promise<void> {
  if (okruszki === null) return;
  const odpowiedz = await wywolaj(kanal, Command.LibraryCollectionList, {
    limit: GRANICA_WYKAZU,
  });
  const kolekcje = wartosc(odpowiedz)?.collections ?? [];
  const dokument = okruszki.ownerDocument;
  zdejmijDzieci(okruszki);
  const korzenSciezki = dokument.createElement('span');
  korzenSciezki.textContent = 'Repozytorium';
  okruszki.appendChild(korzenSciezki);
  const biezaca = kolekcje.find((kolekcja) => kolekcja.id === kontekst.kolekcja());
  if (biezaca === undefined) return;
  const liscie = dokument.createElement('span');
  liscie.setAttribute('aria-current', 'page');
  liscie.textContent = biezaca.name;
  okruszki.appendChild(liscie);
}

async function pobierzPliki(
  kanal: Kanal,
  fraza: string,
  idKolekcji: string,
): Promise<LibraryFile[]> {
  if (fraza !== '') {
    const znalezione = await wywolaj(kanal, Command.LibraryFileSearch, {
      query: fraza,
      limit: GRANICA_WYKAZU,
    });
    return wartosc(znalezione)?.files ?? [];
  }
  const wykaz = await wywolaj(kanal, Command.LibraryFileList, {
    ...(idKolekcji === '' ? {} : { collectionId: idKolekcji }),
    limit: GRANICA_WYKAZU,
  });
  return wartosc(wykaz)?.files ?? [];
}

async function pobierzDziennik(kanal: Kanal): Promise<LibraryAuditEntry[]> {
  const odpowiedz = await wywolaj(kanal, Command.LibraryAuditList, { limit: GRANICA_WYKAZU });
  return wartosc(odpowiedz)?.entries ?? [];
}

async function przeniesPlik(kanal: Kanal, idPliku: string, sciezka: string): Promise<void> {
  const odpowiedz = await wywolaj(kanal, Command.LibraryFileMove, {
    fileIds: [idPliku],
    path: sciezka,
  });
  const wynikMozeMoze = wartosc(odpowiedz);
  if (wynikMozeMoze === null) return;
  const wynikMoze = wynikMozeMoze;
  const wynik = wynikMoze;
  oglos('Library', `Przeniesiono plików: ${String(wynik.movedCount)}.`);
}

async function wgrajPlik(
  kanal: Kanal,
  nazwa: string,
  typ: string,
  tresc64: string,
  idKolekcji: string,
): Promise<void> {
  const odpowiedz = await wywolaj(kanal, Command.LibraryFileUpload, {
    name: nazwa,
    ...(typ === '' ? {} : { mimeType: typ }),
    contentBase64: tresc64,
  });
  const wgrany = wartosc(odpowiedz);
  if (wgrany === null || idKolekcji === '') return;
  await wywolaj(kanal, Command.LibraryCollectionAssign, {
    collectionId: idKolekcji,
    fileIds: [wgrany.file.id],
  });
}

async function skanujDuplikaty(kanal: Kanal, idKolekcji: string): Promise<void> {
  const odpowiedz = await wywolaj(kanal, Command.LibraryDuplicateScan, {
    ...(idKolekcji === '' ? {} : { collectionId: idKolekcji }),
    kinds: [LibraryDuplicateKind.Exact],
    limit: GRANICA_WYKAZU,
  });
  const wynikMozeMoze = wartosc(odpowiedz);
  if (wynikMozeMoze === null) return;
  const wynikMoze = wynikMozeMoze;
  const wynik = wynikMoze;
  oglos('Library', `Grup duplikatów: ${String(wynik.total)}.`);
}

async function archiwizuj(kanal: Kanal, kontekst: Kontekst): Promise<void> {
  const plikMozeMoze = kontekst.plik();
  if (plikMozeMoze === null) return;
  const plikMoze = plikMozeMoze;
  const plik = plikMoze;
  await wywolaj(kanal, Command.LibraryFileArchive, { fileIds: [plik.id] });
}

async function przywroc(kanal: Kanal, idPliku: string): Promise<void> {
  await wywolaj(kanal, Command.LibraryFileRestore, { fileIds: [idPliku] });
}

async function usun(kanal: Kanal, idPliku: string): Promise<void> {
  await wywolaj(kanal, Command.LibraryFileDelete, { fileIds: [idPliku], confirm: true });
}

function zbudujWiersz(dokument: Document, plik: LibraryFile): HTMLElement {
  const wiersz = dokument.createElement('tr');
  wiersz.setAttribute('draggable', 'true');
  wiersz.appendChild(komorka(dokument, plik.name));
  wiersz.appendChild(komorka(dokument, plik.sourceModuleId ?? ''));
  wiersz.appendChild(komorka(dokument, ''));
  const etykiety = dokument.createElement('td');
  for (const nazwa of plik.tags ?? []) etykiety.appendChild(plakietka(dokument, nazwa));
  wiersz.appendChild(etykiety);
  wiersz.appendChild(komorka(dokument, dzien(plik.updatedAt), 'lb-mono'));
  wiersz.appendChild(komorka(dokument, rozmiar(plik.sizeBytes), 'lb-mono'));
  return wiersz;
}

function naniesDziennik(cialo: Element, wpisy: LibraryAuditEntry[]): void {
  if (wpisy.length === 0) {
    pustkaWiersza(cialo, KOLUMN, 'Rdzeń nie zwrócił wpisów dziennika.');
    return;
  }
  const dokument = cialo.ownerDocument;
  zdejmijDzieci(cialo);
  for (const wpis of wpisy) {
    const wiersz = dokument.createElement('tr');
    wiersz.appendChild(komorka(dokument, wpis.detail ?? wpis.action));
    wiersz.appendChild(komorka(dokument, wpis.action));
    wiersz.appendChild(komorka(dokument, wpis.actor));
    wiersz.appendChild(komorka(dokument, ''));
    wiersz.appendChild(komorka(dokument, dzien(wpis.at), 'lb-mono'));
    wiersz.appendChild(komorka(dokument, '', 'lb-mono'));
    cialo.appendChild(wiersz);
  }
}

function komorka(dokument: Document, tresc: string, klasa = ''): HTMLElement {
  const wezel = dokument.createElement('td');
  if (klasa !== '') wezel.className = klasa;
  wezel.textContent = tresc;
  return wezel;
}

// Prototyp nie ma pola wyboru pliku, a kontrakt żąda treści w base64;
// pole wywołane z pamięci nie wchodzi do układu okna.
function wybierzZDysku(
  dokument: Document,
  odbierz: (nazwa: string, typ: string, tresc64: string) => void,
): void {
  const pole = dokument.createElement('input');
  pole.type = 'file';
  pole.multiple = true;
  pole.addEventListener('change', () => {
    for (const plik of pole.files ?? []) {
      const czytnik = new FileReader();
      czytnik.addEventListener('load', () => {
        const wynik = typeof czytnik.result === 'string' ? czytnik.result : '';
        odbierz(plik.name, plik.type, wynik.slice(wynik.indexOf(',') + 1));
      });
      /* Odczyt zawodzi, gdy plik zniknął, leży na nośniku odłączonym albo
         przeglądarka nie ma prawa go czytać. Bez tego zdania wskazanie takiego
         pliku kończy się ciszą i Operator nie wie, czy import w ogóle ruszył. */
      czytnik.addEventListener('error', () => {
        oglos('Library', `Pliku ${plik.name} nie udało się odczytać z dysku. `
          + 'Sprawdź, czy nadal tam leży i czy masz prawo go czytać.', 'blad');
      });
      czytnik.readAsDataURL(plik);
    }
  });
  pole.click();
}
