// Panel „File Preview" wobec podglądu i metadanych rdzenia. Prototyp ma trzy
// przyciski i jeden wiersz metadanych, więc naddatek podglądu rozstrzyga
// przycisk „Powiększ", a zapis tytułu — edycja tego wiersza.
import {
  Command,
  type LibraryFieldDefinition,
  type LibraryFile,
  type LibraryMetadata,
  type LibraryPreview,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { panel, pozycja, wartosc, zdejmijDzieci, type Kontekst } from './library-wspolne.ts';

const ZNAKI_DOPASOWANE = 2000;
const ZNAKI_POWIEKSZONE = 20000;

export function zwiazPodglad(
  kanal: Kanal,
  korzen: ParentNode,
  idOkna: string,
  kontekst: Kontekst,
  przy: AddEventListenerOptions,
): void {
  const trescMozeMoze = panel(korzen, 'panel-preview');
  if (trescMozeMoze === null) return;
  const trescMoze = trescMozeMoze;
  const tresc = trescMoze;
  const dzieci = [...tresc.children];
  const przyciski = dzieci[0] ?? null;
  const ramka = dzieci[1] ?? null;
  const opisMoze = tresc.querySelector('.lb-mono');
  if (opisMoze === null) return;
  const opis = opisMoze;
  if (ramka === null || !(opis instanceof HTMLElement)) return;

  let znakow = ZNAKI_DOPASOWANE;
  let pola: LibraryFieldDefinition[] = [];

  async function odswiez(): Promise<void> {
    const plik = kontekst.plik();
    if (plik === null) {
      ramka.textContent = 'Wybierz plik w wykazie, aby rdzeń podał jego podgląd.';
      opis.textContent = '';
      zdejmijSchemat(tresc);
      return;
    }
    const [odpowiedzPodgladu, odpowiedzMetadanych, odpowiedzSchematu] = await Promise.all([
      wywolaj(kanal, Command.LibraryFilePreview, { fileId: plik.id, page: 1, maxChars: znakow }),
      wywolaj(kanal, Command.LibraryMetadataGet, { fileId: plik.id, includeTechnical: true }),
      wywolaj(kanal, Command.LibrarySchemaGet, {}),
    ]);
    naniesPodglad(ramka, wartosc(odpowiedzPodgladu)?.preview ?? null);
    opis.textContent = opiszMetadane(plik, wartosc(odpowiedzMetadanych)?.metadata ?? null);
    pola = wartosc(odpowiedzSchematu)?.fields ?? [];
    naniesSchemat(tresc, pola);
  }

  przyciski?.addEventListener(
    'click',
    (zdarzenie) => {
      const cel = zdarzenie.target;
      if (!(cel instanceof Element)) return;
      const podpis = cel.closest('.dn-btn')?.textContent?.trim() ?? '';
      if (podpis === '') return;
      if (podpis.startsWith('Otwórz')) {
        void otworzWStudiu(kanal, idOkna, kontekst.plik());
        return;
      }
      znakow = podpis.startsWith('Powiększ') ? ZNAKI_POWIEKSZONE : ZNAKI_DOPASOWANE;
      void odswiez();
    },
    przy,
  );

  opis.addEventListener(
    'dblclick',
    () => {
      opis.contentEditable = 'true';
      opis.focus();
    },
    przy,
  );

  opis.addEventListener(
    'focusout',
    () => {
      if (opis.contentEditable !== 'true') return;
      opis.contentEditable = 'false';
      const tytul = (opis.textContent ?? '').split('·')[0].trim();
      void zapiszTytul(kanal, kontekst.plik(), tytul).then(odswiez);
    },
    przy,
  );

  tresc.addEventListener(
    'click',
    (zdarzenie) => {
      const cel = zdarzenie.target;
      if (!(cel instanceof Element)) return;
      const wierszMozeMoze = cel.closest<HTMLElement>('[data-pole-schematu]');
      if (wierszMozeMoze === null) return;
      const wierszMoze = wierszMozeMoze;
      const wiersz = wierszMoze;
      const pole = pola.find((kandydat) => kandydat.code === wiersz.dataset.poleSchematu);
      if (pole === undefined) return;
      void zapiszPole(kanal, pole, zdarzenie.shiftKey).then(odswiez);
    },
    przy,
  );

  kontekst.przyPliku(() => {
    void odswiez();
  });
  kontekst.zglosOdswiezenie(() => {
    void odswiez();
  });
  void odswiez();
}

function naniesPodglad(ramka: Element, podglad: LibraryPreview | null): void {
  if (podglad === null) {
    ramka.textContent = 'Rdzeń nie podał podglądu tego pliku.';
    return;
  }
  if (podglad.text !== undefined && podglad.text !== '') {
    ramka.textContent = podglad.text;
    return;
  }
  if (podglad.imageRef !== undefined && podglad.imageRef !== '') {
    zdejmijDzieci(ramka);
    const obraz = ramka.ownerDocument.createElement('img');
    obraz.src = podglad.imageRef;
    obraz.alt = '';
    ramka.appendChild(obraz);
    return;
  }
  ramka.textContent = `Rdzeń podał postać „${podglad.kind}" bez treści do pokazania.`;
}

function opiszMetadane(plik: LibraryFile, metadane: LibraryMetadata | null): string {
  const czesci = [
    metadane?.title ?? plik.name,
    plik.mimeType ?? '',
    metadane?.creator ?? '',
    plik.sourceModuleId ?? '',
    metadane?.language ?? '',
  ].filter((czesc) => czesc !== '');
  return czesci.join(' · ');
}

function naniesSchemat(tresc: HTMLElement, pola: LibraryFieldDefinition[]): void {
  zdejmijSchemat(tresc);
  for (const pole of pola) {
    const wiersz = pozycja(tresc.ownerDocument, pole.label, pole.required ? 'wymagane' : '');
    wiersz.dataset.poleSchematu = pole.code;
    tresc.appendChild(wiersz);
  }
}

function zdejmijSchemat(tresc: HTMLElement): void {
  for (const wiersz of tresc.querySelectorAll('[data-pole-schematu]')) wiersz.remove();
}

async function otworzWStudiu(
  kanal: Kanal,
  idOkna: string,
  plik: LibraryFile | null,
): Promise<void> {
  if (idOkna === '' || plik === null) return;
  await wywolaj(kanal, Command.StudioDocumentOpen, {
    windowId: idOkna,
    libraryFileId: plik.id,
  });
}

async function zapiszTytul(
  kanal: Kanal,
  plik: LibraryFile | null,
  tytul: string,
): Promise<void> {
  if (plik === null || tytul === '') return;
  await wywolaj(kanal, Command.LibraryMetadataSet, {
    fileId: plik.id,
    metadata: { fileId: plik.id, title: tytul },
    replace: false,
  });
}

async function zapiszPole(
  kanal: Kanal,
  pole: LibraryFieldDefinition,
  zdejmij: boolean,
): Promise<void> {
  await wywolaj(kanal, Command.LibrarySchemaSet, {
    field: zdejmij ? pole : { ...pole, required: !pole.required },
    remove: zdejmij,
  });
}
