/**
 * Wiązanie panelu podglądu Studia z rdzeniem. Znacznik należy do Właściciela —
 * ten plik wypełnia stojące węzły, a kartki stron powiela z wzoru zdjętego
 * z treści przykładowej. Render oddaje strony jako zasoby magazynu, a kartka
 * prototypu nie niesie węzła zdolnego unieść obraz strony, więc podgląd
 * pokazuje liczbę i numerację stron, nie ich układ.
 */
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
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';

/** Klasy przycisku w znaczniku prototypu: wciśnięty niesie zarys, spoczywający postać duchową. */
const KLASA_WCISNIETY = 'dn-btn--zarys';
const KLASA_SPOCZYNKU = 'dn-btn--duch';

/** Węzły panelu podglądu, na których wiązanie pracuje; przyciski są opcjonalne, bo ich brak nie odbiera panelowi obszaru stron. */
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

/** Stan wejściowy panelu: dokument okna, profil wydania i nastawy widoku — czym rdzeń opisuje podgląd przed pierwszym renderem. */
interface StanPodgladu {
  dokument: StudioDocument | null;
  profil: StudioExportProfile | null;
  ustawienia: StudioViewSettings | null;
}

/**
 * Wiąże panel podglądu okna Studia z rdzeniem. Zwraca prawdę, gdy znacznik
 * panelu stał i wiązanie zostało założone.
 */
export function zwiazPodglad(kanal: Kanal, idOkna: string): boolean {
  const znalezione = zbierzWezly();
  if (znalezione === null) return false;
  const wezly: WezlyPodgladu = znalezione;
  const kartka = zdejmijWzorStrony(wezly.podglad);
  if (kartka === null) return false;
  const wzorStrony: HTMLElement = kartka;

  // Kartka przykładowa znika, zanim padnie pierwsza odpowiedź rdzenia: pusty
  // obszar jest uczciwy, kartka z cudzym dokumentem — nie.
  wezly.podglad.replaceChildren();

  const formaty = Object.values(StudioExportFormat);
  // Nastawa własna odpada ze zbioru skal: znaczy wartość wpisaną przez
  // Operatora, a panel nie ma pola, w które dałoby się ją wpisać.
  const skale = Object.values(StudioZoomPreset).filter((pozycja) => pozycja !== StudioZoomPreset.Custom);

  let format: StudioExportFormat | null = null;
  let nastawaSkali: StudioZoomPreset | undefined;
  let idProfilu: string | undefined;
  let dokument: StudioDocument | null = null;
  let numerRenderu = 0;

  /** Woła render w formacie czynnym i wymienia kartki; odpowiedź spóźniona wobec kolejnej odpada. */
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

  /** Przestawia nastawy widoku i nanosi na przyciski nastawy zwrócone przez rdzeń. */
  function przestaw(zmiana: StudioViewSetRequest): void {
    void ustawWidok(kanal, { ...zmiana, windowId: idOkna, documentId: dokument?.id }).then((ustawienia) => {
      // Odmowa rdzenia zostawia przyciski nietknięte: nastawa poprzednia obowiązuje dalej.
      if (ustawienia !== null) nanies(wezly, ustawienia);
    });
  }

  wezly.format?.addEventListener('click', () => {
    if (format === null) return;
    format = nastepna(formaty, format);
    opiszFormat(wezly.format, format);
    przerysuj();
  });

  wezly.strona?.addEventListener('click', () => {
    przestaw({ scrollMode: StudioScrollMode.Page });
  });

  wezly.ciagly?.addEventListener('click', () => {
    przestaw({ scrollMode: StudioScrollMode.Continuous });
  });

  wezly.skala?.addEventListener('click', () => {
    nastawaSkali = nastepna(skale, nastawaSkali);
    przestaw({ zoomPreset: nastawaSkali });
  });

  wezly.podziel?.addEventListener('click', () => {
    przestaw({ surfaceMode: StudioSurfaceMode.Split });
  });

  wezly.eksport?.addEventListener('click', () => {
    if (dokument === null || format === null) return;
    void wydaj(kanal, dokument.id, format, idProfilu);
  });

  wezly.biblioteka?.addEventListener('click', () => {
    if (dokument === null || format === null) return;
    void przekazDoBiblioteki(kanal, dokument, format, idProfilu);
  });

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

  return true;
}

/** Odczytuje dokument otwarty w oknie, pierwszy profil wydania i nastawy widoku; odmowa rdzenia oddaje pustkę na każdym z trzech pól. */
async function wczytajStan(kanal: Kanal, idOkna: string): Promise<StanPodgladu> {
  const otwarty = await wywolaj(kanal, Command.StudioDocumentOpen, { windowId: idOkna });
  const dokument = otwarty.udany && otwarty.wynik !== undefined ? otwarty.wynik.document : null;
  const profile = await wywolaj(kanal, Command.StudioExportProfileList, {});
  const profil = profile.udany && profile.wynik !== undefined ? (profile.wynik.profiles.at(0) ?? null) : null;
  const widok = await wywolaj(kanal, Command.StudioViewGet, { windowId: idOkna, documentId: dokument?.id });
  const ustawienia = widok.udany && widok.wynik !== undefined ? widok.wynik.settings : null;
  return { dokument, profil, ustawienia };
}

/** Woła render podglądu i zwraca po kartce na każdy zasób strony; odmowa rdzenia oddaje wykaz pusty. */
async function zbudujStrony(
  kanal: Kanal,
  wzor: HTMLElement,
  zadanie: StudioPreviewRenderRequest,
): Promise<HTMLElement[]> {
  const wynik = await wywolaj(kanal, Command.StudioPreviewRender, zadanie);
  if (!wynik.udany || wynik.wynik === undefined) return [];
  // Numeracja idzie porządkiem zasobów przesuniętym o pierwszą stronę żądania:
  // odpowiedź renderu nie niesie numeru strony osobnym polem.
  const pierwsza = zadanie.pageFrom ?? 1;
  return wynik.wynik.pageAssetIds.map((_zasob, miejsce) => zbudujStrone(wzor, pierwsza + miejsce));
}

/** Klon kartki opisany numerem strony. */
function zbudujStrone(wzor: HTMLElement, numer: number): HTMLElement {
  const kartka = wzor.cloneNode(true) as HTMLElement;
  const podpis = kartka.querySelector('.dn-kartka-numer');
  if (podpis !== null) podpis.textContent = `— ${numer} —`;
  return kartka;
}

/** Przestawia nastawy widoku okna i oddaje nastawy po zmianie; odmowa rdzenia oddaje pustkę i zostawia przyciski w stanie poprzednim. */
async function ustawWidok(kanal: Kanal, zadanie: StudioViewSetRequest): Promise<StudioViewSettings | null> {
  const wynik = await wywolaj(kanal, Command.StudioViewSet, zadanie);
  return wynik.udany && wynik.wynik !== undefined ? wynik.wynik.settings : null;
}

/** Wydaje dokument do formatu czynnego; profil idzie w żądaniu, bo wydanie PDF z paginacją i stopką wymaga ustawień strony. */
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

/** Przekazuje wydanie do Biblioteki dwoma krokami kontraktu: wydanie dokumentu, potem wgranie bajtów wziętych z magazynu. Bajty idą przez magazyn, bo ścieżka wydania jest ścieżką w systemie plików rdzenia i klient jej nie odczyta. */
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

/** Nazwa pliku wydania: nazwa spod ścieżki rdzenia, a gdy rdzeń miejsca nie wskazał — tytuł dokumentu z rozszerzeniem formatu. Dokument bez tytułu nie ma z czego wziąć nazwy i do Biblioteki nie idzie. */
function nazwaWydania(wydanie: StudioExportResult, tytul: string | undefined): string {
  const spodSciezki = wydanie.path?.split(/[\\/]/u).pop();
  if (spodSciezki !== undefined && spodSciezki !== '') return spodSciezki;
  return tytul === undefined || tytul === '' ? '' : `${tytul}.${wydanie.format}`;
}

/** Nanosi nastawy widoku na przyciski: tryb przewijania jako wciśnięcie, skalę jako podpis. */
function nanies(wezly: WezlyPodgladu, ustawienia: StudioViewSettings | null): void {
  oznaczTryb(wezly.strona, ustawienia?.scrollMode === StudioScrollMode.Page);
  oznaczTryb(wezly.ciagly, ustawienia?.scrollMode === StudioScrollMode.Continuous);
  opiszSkale(wezly.skala, ustawienia?.zoomPercent);
}

/** Nadaje przyciskowi postać wciśniętą albo spoczywającą w brzmieniu klas prototypu. */
function oznaczTryb(przycisk: HTMLElement | null, wcisniety: boolean): void {
  if (przycisk === null) return;
  przycisk.classList.toggle(KLASA_WCISNIETY, wcisniety);
  przycisk.classList.toggle(KLASA_SPOCZYNKU, !wcisniety);
}

/** Wpisuje skalę widoku w podpis przycisku; brak skali w nastawach zdejmuje przycisk, bo procent zmyślony mówi o wydruku nieprawdę. */
function opiszSkale(przycisk: HTMLElement | null, procent: number | undefined): void {
  if (przycisk === null) return;
  if (procent === undefined) {
    przycisk.remove();
    return;
  }
  przycisk.textContent = `${procent}% ▾`;
}

/** Wpisuje format czynny w podpis przycisku wyboru formatu, w brzmieniu podpisu prototypu. */
function opiszFormat(przycisk: HTMLElement | null, format: StudioExportFormat): void {
  if (przycisk === null) return;
  przycisk.textContent = `Format: ${format.toUpperCase()} ▾`;
}

/**
 * Format renderu wzięty z odpowiedzi rdzenia: z profilu wydania, a bez profilu
 * z formatu samego dokumentu. Pustka znaczy dokument, którego formatu zbiór
 * wydania kontraktu nie niesie.
 */
function formatRenderu(
  formaty: StudioExportFormat[],
  profil: StudioExportProfile | null,
  dokument: StudioDocument | null,
): StudioExportFormat | null {
  const zProfilu = formaty.find((pozycja) => pozycja === profil?.format);
  return zProfilu ?? formaty.find((pozycja) => pozycja === dokument?.format) ?? null;
}

/** Zdejmuje przyciski wydania; bez formatu wziętego z rdzenia panel nie ma czym renderować ani co wydać. */
function zdejmijWydanie(wezly: WezlyPodgladu): void {
  wezly.format?.remove();
  wezly.eksport?.remove();
  wezly.biblioteka?.remove();
}

/** Kolejna pozycja zbioru kontraktu po pozycji czynnej; pozycja spoza zbioru bierze pozycję pierwszą. Kliknięcie obchodzi zbiór, bo prototyp nie niesie węzła listy rozwijanej, a wiązanie węzłów nie buduje. */
function nastepna<T>(zbior: readonly T[], czynna: T | undefined): T {
  const miejsce = czynna === undefined ? -1 : zbior.indexOf(czynna);
  return zbior[(miejsce + 1) % zbior.length];
}

/** Zdejmuje wzór kartki z treści przykładowej wraz z tytułem i akapitem: render oddaje wyłącznie zasoby stron, więc tekstu na kartce nie ma z czego wziąć. */
function zdejmijWzorStrony(podglad: HTMLElement): HTMLElement | null {
  const kartka = podglad.querySelector('.dn-kartka');
  if (!(kartka instanceof HTMLElement)) return null;
  kartka.querySelector('.dn-kartka-tytul')?.remove();
  for (const akapit of kartka.querySelectorAll('p:not(.dn-kartka-numer)')) akapit.remove();
  return kartka.cloneNode(true) as HTMLElement;
}

/** Wskazuje węzły panelu podglądu; pustka znaczy, że panel nie stoi w dokumencie. */
function zbierzWezly(): WezlyPodgladu | null {
  const panel = document.querySelector('#panel-preview');
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

/** Przyciski stojące wprost w wierszu panelu, w porządku znacznika. */
function przyciski(panel: Element, wybor: string): HTMLElement[] {
  const wiersz = panel.querySelector(wybor);
  if (wiersz === null) return [];
  return [...wiersz.querySelectorAll(':scope > button')].filter(
    (kandydat): kandydat is HTMLElement => kandydat instanceof HTMLElement,
  );
}
