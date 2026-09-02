/**
 * Panel biblioteki projektu WorkSpace: wykaz plików, znaczniki, wgranie pliku,
 * wersje, wydobycie tekstu i scalanie duplikatów. Wywołania przeniesione
 * z poprzedniego klienta (`moduly/workspace/biblioteka-projektu.ts`
 * i `zrodlo-sasiadow.ts`): `workspace.library.list`,
 * `workspace.library.text.extract`, `workspace.library.duplicate.list`,
 * `workspace.library.duplicate.merge`, `library.file.upload`
 * i `library.version.list`.
 */

import {
  Command,
  type LibraryFile,
  type WorkspaceDuplicateGroup,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import { data } from './okno-modulu.ts';
import { odmien } from './workspace-projekt.ts';
import {
  cialoPanelu,
  kopiaWzoru,
  niegotowe,
  plakietkaPanelu,
  przyjmij,
  wpisz,
  type WezlyWorkspace,
} from './workspace-wezly.ts';

export interface StanBiblioteki {
  odswiez: () => Promise<void>;
}

export function zwiazBiblioteke(
  wezly: WezlyWorkspace,
  kanal: Kanal,
  idProjektu: () => string,
  przy: AddEventListenerOptions,
): StanBiblioteki {
  const panel = wezly.biblioteka;
  const cialo = cialoPanelu(panel);
  const pasek = cialo?.querySelector<HTMLElement>('.rt-mp-dol') ?? null;
  const wzorZnacznika = kopiaWzoru(pasek?.querySelector<HTMLElement>('.dn-plakietka') ?? null);
  const wzorWiersza = kopiaWzoru(cialo?.querySelector<HTMLElement>('.dn-wykaz-modulu-poz') ?? null);
  const gniazdo = document.createElement('div');
  gniazdo.className = 'dn-wykaz-modulu';
  let znacznik = '';

  // Żetony widoku prototypu przełączają układ, którego rdzeń nie opisuje.
  for (const zeton of pasek?.querySelectorAll('.sta-chip') ?? []) zeton.remove();
  for (const stary of cialo?.querySelectorAll('.dn-wykaz-modulu-poz') ?? []) stary.remove();
  cialo?.appendChild(gniazdo);

  const odswiez = async (): Promise<void> => {
    const projekt = idProjektu();
    if (projekt === '') return;
    const odpowiedz = await pobierzPliki(kanal, projekt, znacznik);
    const pliki = odpowiedz?.files ?? [];
    wpisz(plakietkaPanelu(panel), odmien(odpowiedz?.total ?? pliki.length, ['plik', 'pliki', 'plików']));
    opiszZnaczniki(pasek, wzorZnacznika, pliki, znacznik, (wybrany) => {
      znacznik = wybrany;
      void odswiez();
    }, przy);
    gniazdo.replaceChildren();
    if (pliki.length === 0) {
      niegotowe(gniazdo, 'Biblioteka tego projektu nie ma jeszcze plików.');
    } else {
      for (const plik of pliki) {
        const wiersz = wierszPliku(wzorWiersza, plik);
        if (wiersz !== null) gniazdo.appendChild(wiersz);
      }
    }
    await dolozDuplikaty(kanal, projekt, gniazdo, wzorWiersza);
  };

  gniazdo.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const wierszMoze = cel.closest<HTMLElement>('[data-plik], [data-duplikat]');
    if (wierszMoze === null) return;
    const wiersz = wierszMoze;
    const grupa = wiersz.dataset.duplikat ?? '';
    void (grupa === ''
      ? wydobadzTekst(kanal, idProjektu(), wiersz)
      : scalDuplikaty(kanal, idProjektu(), wiersz, odswiez));
  }, przy);

  zwiazWgranie(panel, kanal, idProjektu, odswiez, przy);
  return { odswiez };
}

/** Wgranie idzie treścią w base64 — poprzedni klient podawał ją tak samo. */
function zwiazWgranie(
  panel: HTMLElement | null,
  kanal: Kanal,
  idProjektu: () => string,
  odswiez: () => Promise<void>,
  przy: AddEventListenerOptions,
): void {
  const przycisk = panel?.querySelector<HTMLElement>(
    '.sta-okno-belka .dn-btn-ikona:not([data-okno-zamknij])',
  );
  if (przycisk === null || przycisk === undefined) return;
  const wybor = document.createElement('input');
  wybor.type = 'file';
  wybor.hidden = true;
  panel?.appendChild(wybor);
  przycisk.addEventListener('click', () => wybor.click(), przy);
  wybor.addEventListener('change', () => {
    const plik = wybor.files?.[0];
    if (plik !== undefined) void wgrajPlik(kanal, idProjektu(), plik, odswiez);
    wybor.value = '';
  }, przy);
}

async function wgrajPlik(
  kanal: Kanal,
  idProjektu: string,
  plik: File,
  odswiez: () => Promise<void>,
): Promise<void> {
  if (idProjektu === '') return;
  const bajty = new Uint8Array(await plik.arrayBuffer());
  let napis = '';
  for (const bajt of bajty) napis += String.fromCharCode(bajt);
  const wynik = await wywolaj(kanal, Command.LibraryFileUpload, {
    name: plik.name,
    projectId: idProjektu,
    sourceModuleId: 'workspace',
    mimeType: plik.type,
    contentBase64: globalThis.btoa(napis),
  });
  if (przyjmij('Biblioteka projektu', wynik) === null) return;
  await odswiez();
}

function opiszZnaczniki(
  pasek: HTMLElement | null,
  wzor: HTMLElement | null,
  pliki: LibraryFile[],
  wybrany: string,
  wybierz: (znacznik: string) => void,
  przy: AddEventListenerOptions,
): void {
  if (pasek === null || wzor === null) return;
  for (const stary of pasek.querySelectorAll('.dn-plakietka')) stary.remove();
  const znaczniki = new Set<string>(wybrany === '' ? [] : [wybrany]);
  for (const plik of pliki) for (const cecha of plik.tags ?? []) znaczniki.add(cecha);
  for (const cecha of [...znaczniki].sort()) {
    const zeton = kopiaWzoru(wzor);
    if (zeton === null) continue;
    zeton.textContent = cecha;
    zeton.setAttribute('aria-pressed', String(cecha === wybrany));
    zeton.addEventListener('click', () => wybierz(cecha === wybrany ? '' : cecha), przy);
    pasek.appendChild(zeton);
  }
}

function wierszPliku(wzor: HTMLElement | null, plik: LibraryFile): HTMLElement | null {
  const wierszMoze = kopiaWzoru(wzor);
  if (wierszMoze === null) return null;
  const wiersz = wierszMoze;
  const meta = wiersz.querySelector<HTMLElement>('.dn-meta');
  wiersz.replaceChildren();
  wiersz.append(`${plik.name} `);
  if (meta !== null) {
    meta.textContent = [rozmiar(plik.sizeBytes), data(plik.updatedAt)]
      .filter((czesc) => czesc !== '')
      .join(' · ');
    wiersz.appendChild(meta);
  }
  wiersz.dataset.plik = plik.id;
  return wiersz;
}

function rozmiar(bajty: number | undefined): string {
  if (bajty === undefined) return '';
  if (bajty < 1024) return `${bajty} B`;
  if (bajty < 1024 * 1024) return `${Math.round(bajty / 1024)} kB`;
  return `${(bajty / (1024 * 1024)).toFixed(1)} MB`;
}

async function pobierzPliki(
  kanal: Kanal,
  idProjektu: string,
  znacznik: string,
): Promise<{ files: LibraryFile[]; total?: number } | null> {
  const wynik = await wywolaj(kanal, Command.WorkspaceLibraryList, {
    projectId: idProjektu,
    limit: 100,
    ...(znacznik === '' ? {} : { tag: znacznik }),
  });
  return przyjmij('Biblioteka projektu', wynik) ?? null;
}

/** Grupa duplikatów wchodzi w ten sam wykaz; kliknięcie zleca scalenie. */
async function dolozDuplikaty(
  kanal: Kanal,
  idProjektu: string,
  gniazdo: HTMLElement,
  wzor: HTMLElement | null,
): Promise<void> {
  const wynik = await wywolaj(kanal, Command.WorkspaceLibraryDuplicateList, {
    projectId: idProjektu,
  });
  if (!wynik.udany || wynik.wynik === undefined) return;
  for (const grupa of wynik.wynik.groups) {
    const wiersz = wierszDuplikatu(wzor, grupa);
    if (wiersz !== null) gniazdo.appendChild(wiersz);
  }
}

function wierszDuplikatu(
  wzor: HTMLElement | null,
  grupa: WorkspaceDuplicateGroup,
): HTMLElement | null {
  const wierszMoze = kopiaWzoru(wzor);
  if (wierszMoze === null) return null;
  const wiersz = wierszMoze;
  const meta = wiersz.querySelector<HTMLElement>('.dn-meta');
  const zachowaj = grupa.suggestedKeepFileId ?? grupa.files[0]?.id ?? '';
  wiersz.replaceChildren();
  wiersz.append(`${grupa.files[0]?.name ?? 'plik'} `);
  if (meta !== null) {
    meta.textContent = `duplikaty: ${grupa.files.length} · scal`;
    wiersz.appendChild(meta);
  }
  wiersz.dataset.duplikat = grupa.checksum;
  wiersz.dataset.zachowaj = zachowaj;
  wiersz.dataset.scalane = grupa.files
    .map((plik) => plik.id)
    .filter((identyfikator) => identyfikator !== zachowaj)
    .join(',');
  return wiersz;
}

async function wydobadzTekst(kanal: Kanal, idProjektu: string, wiersz: HTMLElement): Promise<void> {
  const idPliku = wiersz.dataset.plik ?? '';
  if (idProjektu === '' || idPliku === '') return;
  const wynik = await wywolaj(kanal, Command.WorkspaceLibraryTextExtract, {
    projectId: idProjektu,
    fileId: idPliku,
  });
  const wydobycie = przyjmij('Wydobycie tekstu', wynik)?.extraction;
  if (wydobycie === undefined) return;
  const wersje = await wywolaj(kanal, Command.LibraryVersionList, { fileId: idPliku });
  const czesci = [
    odmien(wydobycie.characterCount, ['znak', 'znaki', 'znaków']),
    wydobycie.method,
  ];
  const liczbaWersji = wersje.wynik?.versions.length;
  if (liczbaWersji !== undefined) czesci.push(`wersje: ${liczbaWersji}`);
  wpisz(wiersz.querySelector('.dn-meta'), czesci.join(' · '));
}

async function scalDuplikaty(
  kanal: Kanal,
  idProjektu: string,
  wiersz: HTMLElement,
  odswiez: () => Promise<void>,
): Promise<void> {
  const zachowaj = wiersz.dataset.zachowaj ?? '';
  const scalane = (wiersz.dataset.scalane ?? '').split(',').filter((cel) => cel !== '');
  if (idProjektu === '' || zachowaj === '' || scalane.length === 0) return;
  const wynik = await wywolaj(kanal, Command.WorkspaceLibraryDuplicateMerge, {
    projectId: idProjektu,
    keepFileId: zachowaj,
    mergedFileIds: scalane,
  });
  const odpowiedzMoze = przyjmij('Scalenie duplikatów', wynik);
  if (odpowiedzMoze === null) return;
  const odpowiedz = odpowiedzMoze;
  oglos('Biblioteka projektu', `Scalono ${odpowiedz.merged} plików.`);
  await odswiez();
}
