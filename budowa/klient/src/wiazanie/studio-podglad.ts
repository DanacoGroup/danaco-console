// Panel podglądu Studia: render oddaje strony jako zasoby magazynu, a kartka
// prototypu nie uniesie obrazu strony — panel pokazuje liczbę i numerację stron.
import {
  AssetContentDisposition,
  Command,
  StudioExportFormat,
  StudioScrollMode,
  StudioSurfaceMode,
  StudioZoomPreset,
  type StudioDocument,
  type StudioExportProfile,
  type StudioExportResult,
  type StudioPreviewRenderRequest,
  type StudioViewSetRequest,
  type StudioViewSettings,
} from '../../../shared/contract.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';

const KLASA_WCISNIETY = 'dn-btn--zarys';
const KLASA_SPOCZYNKU = 'dn-btn--duch';

interface WezlyPodgladu {
  podglad: HTMLElement;
  format: HTMLElement | null;
  strona: HTMLElement | null;
  ciagly: HTMLElement | null;
  skala: HTMLElement | null;
  podziel: HTMLElement | null;
  eksport: HTMLElement | null;
  biblioteka: HTMLElement | null;
}

interface StanPodgladu {
  dokument: StudioDocument | null;
  profil: StudioExportProfile | null;
  ustawienia: StudioViewSettings | null;
}

// Wzór kartki zdjęty raz: obszar stron opróżniony wzoru już nie oddaje.
let wzorKartki: HTMLElement | null = null;

// Woła się przy montażu okna: kartka prototypu niesie cudzy dokument.
export function zdejmijTrescPrzykladowaPodgladu(korzen: ParentNode): boolean {
  return przygotujPanel(korzen) !== null;
}

function przygotujPanel(korzen: ParentNode): { wezly: WezlyPodgladu; wzorStrony: HTMLElement } | null {
  const znalezione = zbierzWezly(korzen);
  if (znalezione === null) return null;
  wzorKartki ??= zdejmijWzorStrony(znalezione.podglad);
  if (wzorKartki === null) return null;
  znalezione.podglad.replaceChildren();
  return { wezly: znalezione, wzorStrony: wzorKartki };
}

export function zwiazPodglad(
  kanal: Kanal,
  idOkna: string,
  korzen: ParentNode,
): Odsubskrybuj | null {
  const przygotowany = przygotujPanel(korzen);
  if (przygotowany === null) return null;
  // Nasłuchy panelu schodzą razem z kartą, zdjęte sterownikiem przerwania.
  const sterowanie = new AbortController();
  const przy = { signal: sterowanie.signal };
  const wezly: WezlyPodgladu = przygotowany.wezly;
  const wzorStrony: HTMLElement = przygotowany.wzorStrony;

  const formaty = Object.values(StudioExportFormat);
  // Nastawa własna odpada ze zbioru skal: panel nie ma pola na jej wpisanie.
  const skale = Object.values(StudioZoomPreset).filter((pozycja) => pozycja !== StudioZoomPreset.Custom);

  let format: StudioExportFormat | null = null;
  let nastawaSkali: StudioZoomPreset | undefined;
  let idProfilu: string | undefined;
  let dokument: StudioDocument | null = null;
  let numerRenderu = 0;

  function przerysuj(): void {
    if (dokument === null || format === null) return;
    const numer = (numerRenderu += 1);
    void zbudujStrony(kanal, wzorStrony, {
      documentId: dokument.id,
      versionId: dokument.versionId,
      format,
      profileId: idProfilu,
      windowId: idOkna,
    }).then((strony) => {
      if (numer !== numerRenderu) return;
      wezly.podglad.replaceChildren(...strony);
    });
  }

  function przestaw(zmiana: StudioViewSetRequest): void {
    void ustawWidok(kanal, { ...zmiana, windowId: idOkna, documentId: dokument?.id }).then((ustawienia) => {
      // Odmowa rdzenia zostawia nastawę poprzednią w mocy.
      if (ustawienia !== null) nanies(wezly, ustawienia);
    });
  }

  wezly.format?.addEventListener('click', () => {
    if (format === null) return;
    format = nastepna(formaty, format);
    opiszFormat(wezly.format, format);
    przerysuj();
  }, przy);

  wezly.strona?.addEventListener('click', () => {
    przestaw({ scrollMode: StudioScrollMode.Page });
  }, przy);

  wezly.ciagly?.addEventListener('click', () => {
    przestaw({ scrollMode: StudioScrollMode.Continuous });
  }, przy);

  wezly.skala?.addEventListener('click', () => {
    nastawaSkali = nastepna(skale, nastawaSkali);
    przestaw({ zoomPreset: nastawaSkali });
  }, przy);

  wezly.podziel?.addEventListener('click', () => {
    przestaw({ surfaceMode: StudioSurfaceMode.Split });
  }, przy);

  wezly.eksport?.addEventListener('click', () => {
    if (dokument === null || format === null) return;
    void wydaj(kanal, dokument.id, format, idProfilu);
  }, przy);

  wezly.biblioteka?.addEventListener('click', () => {
    if (dokument === null || format === null) return;
    void przekazDoBiblioteki(kanal, dokument, format, idProfilu);
  }, przy);

  void wczytajStan(kanal, idOkna).then((stan) => {
    dokument = stan.dokument;
    idProfilu = stan.profil?.id;
    format = formatRenderu(formaty, stan.profil, stan.dokument);
    nastawaSkali = stan.ustawienia?.zoomPreset;
    if (format === null) zdejmijWydanie(wezly);
    else opiszFormat(wezly.format, format);
    nanies(wezly, stan.ustawienia);
    przerysuj();
  });

  return () => {
    sterowanie.abort();
  };
}

async function wczytajStan(kanal: Kanal, idOkna: string): Promise<StanPodgladu> {
  const otwarty = await wywolaj(kanal, Command.StudioDocumentOpen, { windowId: idOkna });
  const dokument = otwarty.udany && otwarty.wynik !== undefined ? otwarty.wynik.document : null;
  const profile = await wywolaj(kanal, Command.StudioExportProfileList, {});
  const profil = profile.udany && profile.wynik !== undefined ? (profile.wynik.profiles.at(0) ?? null) : null;
  const widok = await wywolaj(kanal, Command.StudioViewGet, { windowId: idOkna, documentId: dokument?.id });
  const ustawienia = widok.udany && widok.wynik !== undefined ? widok.wynik.settings : null;
  return { dokument, profil, ustawienia };
}

async function zbudujStrony(
  kanal: Kanal,
  wzor: HTMLElement,
  zadanie: StudioPreviewRenderRequest,
): Promise<HTMLElement[]> {
  const wynik = await wywolaj(kanal, Command.StudioPreviewRender, zadanie);
  if (!wynik.udany || wynik.wynik === undefined) return [];
  // Numeracja idzie porządkiem zasobów: odpowiedź renderu nie niesie numeru strony.
  const pierwsza = zadanie.pageFrom ?? 1;
  return wynik.wynik.pageAssetIds.map((_zasob, miejsce) => zbudujStrone(wzor, pierwsza + miejsce));
}

function zbudujStrone(wzor: HTMLElement, numer: number): HTMLElement {
  const kartka = wzor.cloneNode(true) as HTMLElement;
  const podpis = kartka.querySelector('.dn-kartka-numer');
  if (podpis !== null) podpis.textContent = `— ${numer} —`;
  return kartka;
}

async function ustawWidok(kanal: Kanal, zadanie: StudioViewSetRequest): Promise<StudioViewSettings | null> {
  const wynik = await wywolaj(kanal, Command.StudioViewSet, zadanie);
  return wynik.udany && wynik.wynik !== undefined ? wynik.wynik.settings : null;
}

async function wydaj(
  kanal: Kanal,
  idDokumentu: string,
  format: StudioExportFormat,
  idProfilu: string | undefined,
): Promise<StudioExportResult | null> {
  const wynik = await wywolaj(kanal, Command.StudioDocumentExportFormat, {
    documentId: idDokumentu,
    format,
    profileId: idProfilu,
  });
  return wynik.udany && wynik.wynik !== undefined ? wynik.wynik.result : null;
}

/* Bajty idą przez magazyn, bo ścieżka wydania jest ścieżką w systemie plików rdzenia. */
async function przekazDoBiblioteki(
  kanal: Kanal,
  dokument: StudioDocument,
  format: StudioExportFormat,
  idProfilu: string | undefined,
): Promise<void> {
  const wydanie = await wydaj(kanal, dokument.id, format, idProfilu);
  if (wydanie === null || wydanie.assetId === undefined) return;
  const nazwa = nazwaWydania(wydanie, dokument.title);
  if (nazwa === '') return;
  const tresc = await wywolaj(kanal, Command.DesignAssetContentGet, {
    assetId: wydanie.assetId,
    disposition: AssetContentDisposition.Inline,
  });
  if (!tresc.udany || tresc.wynik === undefined) return;
  const bajty = tresc.wynik.contentBase64;
  if (bajty === undefined) return;
  await wywolaj(kanal, Command.LibraryFileUpload, {
    name: nazwa,
    mimeType: tresc.wynik.mediaType,
    contentBase64: bajty,
  });
}

function nazwaWydania(wydanie: StudioExportResult, tytul: string | undefined): string {
  const spodSciezki = wydanie.path?.split(/[\\/]/u).pop();
  if (spodSciezki !== undefined && spodSciezki !== '') return spodSciezki;
  return tytul === undefined || tytul === '' ? '' : `${tytul}.${wydanie.format}`;
}

function nanies(wezly: WezlyPodgladu, ustawienia: StudioViewSettings | null): void {
  oznaczTryb(wezly.strona, ustawienia?.scrollMode === StudioScrollMode.Page);
  oznaczTryb(wezly.ciagly, ustawienia?.scrollMode === StudioScrollMode.Continuous);
  opiszSkale(wezly.skala, ustawienia?.zoomPercent);
}

function oznaczTryb(przycisk: HTMLElement | null, wcisniety: boolean): void {
  if (przycisk === null) return;
  przycisk.classList.toggle(KLASA_WCISNIETY, wcisniety);
  przycisk.classList.toggle(KLASA_SPOCZYNKU, !wcisniety);
}

/* Brak skali zdejmuje przycisk: procent zmyślony mówi o wydruku nieprawdę. */
function opiszSkale(przycisk: HTMLElement | null, procent: number | undefined): void {
  if (przycisk === null) return;
  if (procent === undefined) {
    przycisk.remove();
    return;
  }
  przycisk.textContent = `${procent}% ▾`;
}

function opiszFormat(przycisk: HTMLElement | null, format: StudioExportFormat): void {
  if (przycisk === null) return;
  przycisk.textContent = `Format: ${format.toUpperCase()} ▾`;
}

/* Format renderu z profilu wydania, a bez profilu z formatu dokumentu. */
function formatRenderu(
  formaty: StudioExportFormat[],
  profil: StudioExportProfile | null,
  dokument: StudioDocument | null,
): StudioExportFormat | null {
  const zProfilu = formaty.find((pozycja) => pozycja === profil?.format);
  return zProfilu ?? formaty.find((pozycja) => pozycja === dokument?.format) ?? null;
}

function zdejmijWydanie(wezly: WezlyPodgladu): void {
  wezly.format?.remove();
  wezly.eksport?.remove();
  wezly.biblioteka?.remove();
}

/* Kliknięcie obchodzi zbiór: prototyp nie niesie węzła listy rozwijanej. */
function nastepna<T>(zbior: readonly T[], czynna: T | undefined): T {
  const miejsce = czynna === undefined ? -1 : zbior.indexOf(czynna);
  return zbior[(miejsce + 1) % zbior.length];
}

/* Render oddaje wyłącznie zasoby stron, więc tekstu na kartce nie ma skąd wziąć. */
function zdejmijWzorStrony(podglad: HTMLElement): HTMLElement | null {
  const kartka = podglad.querySelector('.dn-kartka');
  if (!(kartka instanceof HTMLElement)) return null;
  kartka.querySelector('.dn-kartka-tytul')?.remove();
  for (const akapit of kartka.querySelectorAll('p:not(.dn-kartka-numer)')) akapit.remove();
  return kartka.cloneNode(true) as HTMLElement;
}

function zbierzWezly(korzen: ParentNode): WezlyPodgladu | null {
  const panel = korzen.querySelector('#panel-preview');
  if (panel === null) return null;
  const podglad = panel.querySelector('.st-podglad');
  if (!(podglad instanceof HTMLElement)) return null;
  const gorne = przyciski(panel, '.st-panel-lista--bez-dolu .st-panel-wiersz');
  const dolne = przyciski(panel, '.st-panel-lista--bez-gory .st-panel-wiersz');
  const skala = panel.querySelector('.st-panel-lista--bez-dolu .dn-meta button');
  return {
    podglad,
    format: gorne.at(0) ?? null,
    strona: gorne.at(1) ?? null,
    ciagly: gorne.at(2) ?? null,
    skala: skala instanceof HTMLElement ? skala : null,
    podziel: dolne.at(0) ?? null,
    eksport: dolne.at(1) ?? null,
    biblioteka: dolne.at(2) ?? null,
  };
}

function przyciski(panel: Element, wybor: string): HTMLElement[] {
  const wiersz = panel.querySelector(wybor);
  if (wiersz === null) return [];
  return [...wiersz.querySelectorAll(':scope > button')].filter(
    (kandydat): kandydat is HTMLElement => kandydat instanceof HTMLElement,
  );
}
