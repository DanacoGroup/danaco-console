// Czynności wytwarzania w oknie Apps: wyrób i jego schemat, obszar roboczy
// warstw, plan etapów i kamieni, architektura, motyw, podgląd i sondy.
import {
  AppArchitectureTemplate,
  AppDeployEnvironment,
  AppEndpointMethod,
  AppExportFormat,
  AppWorkspaceLayer,
  Command,
  type AppSchemaTable,
  type AppTimelineEntry,
  type AppValidationIssue,
  type DeveloperFile,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';

const NAGLOWEK = 'Aplikacje';

const PRZYCISKI: ReadonlyArray<readonly [string, string]> = [
  ['wyrob-zapisz', 'Zapisz wyrób'],
  ['wyrob-odczytaj', 'Odczytaj wyrób'],
  ['schemat', 'Odczytaj schemat'],
  ['obszar', 'Odczytaj warstwę widoku'],
  ['obszar-zapisz', 'Zapisz plik warstwy'],
  ['dzieje', 'Odczytaj dzieje'],
  ['etap', 'Dołóż etap'],
  ['kamien', 'Dołóż kamień'],
  ['kamien-usun', 'Usuń kamień'],
  ['architektura', 'Załóż architekturę'],
  ['architektura-odczytaj', 'Odczytaj architekturę'],
  ['architektura-sprawdz', 'Sprawdź architekturę'],
  ['architektura-wydaj', 'Wydaj architekturę'],
  ['przypis', 'Dołóż przypis'],
  ['motyw-odczytaj', 'Odczytaj motyw'],
  ['motyw-ustaw', 'Ustaw motyw'],
  ['podglad', 'Uruchom podgląd'],
  ['podglad-zatrzymaj', 'Zatrzymaj podgląd'],
  ['wydajnosc', 'Zbadaj wydajność'],
  ['sonda', 'Sonduj punkt styku'],
];

/*
zwiazWytwarzanieAplikacji stawia pas nad panelem budowniczego.

Pas stoi nad ciałem panelu, bo ciało jest wymieniane przy każdym odświeżeniu
wykazu; pas postawiony w nim znikałby razem z pozycjami.
*/
export function zwiazWytwarzanieAplikacji(
  kanal: Kanal,
  korzen: Element,
  idOkna: () => string,
  odswiez: () => void,
  przy: AddEventListenerOptions,
): void {
  postawPas(korzen);
  korzen.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const czynnosc = cel.closest<HTMLElement>('[data-apps]')?.dataset.apps;
    if (czynnosc === undefined) return;
    zdarzenie.stopPropagation();
    void wykonaj(kanal, korzen, czynnosc, idOkna(), odswiez);
  }, przy);
}

function postawPas(korzen: Element): void {
  const panel = korzen.querySelector('#panel-builder');
  const cialo = panel?.querySelector('.sta-okno-tresc');
  if (panel === null || panel === undefined || cialo === null || cialo === undefined) return;
  if (panel.querySelector('[data-apps-wpis]') !== null) return;
  const pas = korzen.ownerDocument.createElement('div');
  pas.className = 'dn-pas-dzialan';
  const pole = korzen.ownerDocument.createElement('input');
  pole.type = 'text';
  pole.className = 'dn-pole dn-pole--sm';
  pole.placeholder = 'Nazwa, ścieżka albo treść czynności';
  pole.setAttribute('aria-label', 'Nazwa, ścieżka albo treść czynności okna aplikacji');
  pole.dataset.appsWpis = '';
  pas.appendChild(pole);
  for (const [czynnosc, etykieta] of PRZYCISKI) pas.appendChild(przycisk(korzen, czynnosc, etykieta));
  panel.insertBefore(pas, cialo);
}

function przycisk(korzen: Element, czynnosc: string, etykieta: string): HTMLButtonElement {
  const wezel = korzen.ownerDocument.createElement('button');
  wezel.type = 'button';
  wezel.className = 'dn-btn dn-btn--duch dn-btn--sm';
  wezel.dataset.apps = czynnosc;
  wezel.textContent = etykieta;
  return wezel;
}

function wpis(korzen: Element, znacznik: string): string {
  const pole = korzen.querySelector<HTMLInputElement>(`[${znacznik}]`);
  return (pole?.value ?? '').trim();
}

function pole(korzen: Element): string {
  return wpis(korzen, 'data-apps-wpis');
}

function czesci(korzen: Element): string[] {
  return pole(korzen).split('|').map((czesc) => czesc.trim());
}

/* Drugie naciśnięcie potwierdza czynność nieodwracalną; napis wraca sam, więc
   uzbrojenie nie zostaje na przycisku na stałe. */
export function potwierdzone(korzen: Element, wybor: string, czynnosc: string): boolean {
  const guzik = korzen.querySelector<HTMLElement>(`[${wybor}="${czynnosc}"]`);
  if (guzik === null) return true;
  if (guzik.dataset.uzbrojone === 'tak') {
    delete guzik.dataset.uzbrojone;
    return true;
  }
  const napis = guzik.textContent ?? '';
  guzik.dataset.uzbrojone = 'tak';
  guzik.textContent = 'Naciśnij ponownie';
  globalThis.setTimeout(() => {
    delete guzik.dataset.uzbrojone;
    guzik.textContent = napis;
  }, 5000);
  return false;
}

export function wypelnij(korzen: Element, panel: string, wiersze: string[]): void {
  const cialo = korzen.querySelector(`#${panel} .sta-okno-tresc`);
  if (cialo === null) return;
  cialo.replaceChildren(...wiersze.map((tresc) => {
    const wiersz = cialo.ownerDocument.createElement('div');
    wiersz.className = 'dn-wykaz-modulu-poz';
    wiersz.textContent = tresc;
    return wiersz;
  }));
}

async function wykonaj(
  kanal: Kanal,
  korzen: Element,
  czynnosc: string,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  if (idOkna === '') {
    oglos(NAGLOWEK, 'Rdzeń nie dał okna aplikacji dla tej karty.', 'ostrzezenie');
    return;
  }
  if (czynnosc === 'wyrob-zapisz') return zapiszWyrob(kanal, korzen, idOkna, odswiez);
  if (czynnosc === 'wyrob-odczytaj') return odczytajWyrob(kanal, idOkna);
  if (czynnosc === 'schemat') return odczytajSchemat(kanal, korzen, idOkna);
  if (czynnosc === 'obszar') return odczytajWarstwe(kanal, korzen, idOkna);
  if (czynnosc === 'obszar-zapisz') return zapiszPlikWarstwy(kanal, korzen, idOkna);
  if (czynnosc === 'dzieje') return odczytajDzieje(kanal, korzen, idOkna);
  if (czynnosc === 'etap') return dolozEtap(kanal, korzen, idOkna, odswiez);
  if (czynnosc === 'kamien') return dolozKamien(kanal, korzen, idOkna, odswiez);
  if (czynnosc === 'kamien-usun') return usunKamien(kanal, korzen, idOkna, odswiez);
  if (czynnosc === 'architektura') return zalozArchitekture(kanal, korzen, idOkna, odswiez);
  if (czynnosc === 'architektura-odczytaj') return odczytajArchitekture(kanal, korzen, idOkna);
  if (czynnosc === 'architektura-sprawdz') return sprawdzArchitekture(kanal, korzen, idOkna);
  if (czynnosc === 'architektura-wydaj') return wydajArchitekture(kanal, idOkna);
  if (czynnosc === 'przypis') return dolozPrzypis(kanal, korzen, idOkna);
  if (czynnosc === 'motyw-odczytaj') return odczytajMotyw(kanal, idOkna);
  if (czynnosc === 'motyw-ustaw') return ustawMotyw(kanal, korzen, idOkna);
  if (czynnosc === 'podglad') return uruchomPodglad(kanal, korzen, idOkna);
  if (czynnosc === 'podglad-zatrzymaj') return zatrzymajPodglad(kanal, idOkna);
  if (czynnosc === 'wydajnosc') return zbadajWydajnosc(kanal, korzen, idOkna);
  if (czynnosc === 'sonda') return sondujPunktStyku(kanal, korzen, idOkna);
  /* Czynność spoza obsłużonych odmawia zamiast milczeć: przycisk bez gałęzi
     wyglądałby jak działający. */
  oglos(NAGLOWEK, `Czynność „${czynnosc}" nie jest prowadzona przez to okno.`, 'ostrzezenie');
}

async function zapiszWyrob(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  const [nazwa = '', opis = ''] = czesci(korzen);
  if (nazwa === '') {
    oglos(NAGLOWEK, 'Wyrób potrzebuje nazwy, na przykład „Sklep | witryna zamówień".', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.AppsProductSave, {
    windowId: idOkna,
    name: nazwa,
    ...(opis === '' ? {} : { description: opis }),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu wyrobu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Wyrób „${wynik.wynik.product.name}" zapisany.`);
  odswiez();
}

async function odczytajWyrob(kanal: Kanal, idOkna: string): Promise<void> {
  const wynik = await wywolaj(kanal, Command.AppsProductGet, { windowId: idOkna });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu wyrobu.', 'ostrzezenie');
    return;
  }
  const wyrob = wynik.wynik.product;
  oglos(NAGLOWEK, wyrob === undefined
    ? 'To okno nie ma jeszcze zapisanego wyrobu.'
    : `Wyrób „${wyrob.name}": ${wyrob.description ?? 'bez opisu'}.`);
}

/* Schemat idzie do panelu warstwy rdzenia, bo opisuje jej tablice. */
async function odczytajSchemat(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const wynik = await wywolaj(kanal, Command.AppsSchemaGet, { windowId: idOkna });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu schematu.', 'ostrzezenie');
    return;
  }
  const tablice = wynik.wynik.schema?.tables ?? [];
  if (tablice.length === 0) {
    oglos(NAGLOWEK, 'Rdzeń nie zna tablic tej aplikacji.', 'ostrzezenie');
    return;
  }
  wypelnij(korzen, 'panel-backend', tablice.map((tablica: AppSchemaTable) =>
    `${tablica.name} · ${tablica.columns?.length ?? 0} kolumn`));
}

async function odczytajWarstwe(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const wynik = await wywolaj(kanal, Command.AppsWorkspaceList, {
    windowId: idOkna,
    layer: AppWorkspaceLayer.Frontend,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu warstwy widoku.', 'ostrzezenie');
    return;
  }
  const pliki = wynik.wynik.files;
  wypelnij(korzen, 'panel-frontend', pliki.map((plik: DeveloperFile) => plik.path));
  if (pliki.length === 0) oglos(NAGLOWEK, 'Warstwa widoku nie ma jeszcze plików.');
}

async function zapiszPlikWarstwy(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const [sciezka = '', tresc = ''] = czesci(korzen);
  if (sciezka === '' || tresc === '') {
    oglos(NAGLOWEK,
      'Plik warstwy potrzebuje ścieżki i treści, na przykład „src/widok.ts | export const widok = 1;".',
      'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.AppsWorkspaceUpdate, {
    windowId: idOkna,
    layer: AppWorkspaceLayer.Frontend,
    path: sciezka,
    content: tresc,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu pliku warstwy.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Plik ${wynik.wynik.file.path} zapisany w warstwie widoku.`);
}

async function odczytajDzieje(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const wynik = await wywolaj(kanal, Command.AppsTimelineList, { windowId: idOkna, limit: 100 });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu dziejów.', 'ostrzezenie');
    return;
  }
  const wpisy = wynik.wynik.entries;
  wypelnij(korzen, 'panel-plan', wpisy.map((wpis: AppTimelineEntry) => `${wpis.kind} · ${wpis.summary}`));
  if (wpisy.length === 0) oglos(NAGLOWEK, 'Dzieje tego okna są jeszcze puste.');
}

async function dolozEtap(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  const nazwa = pole(korzen);
  if (nazwa === '') {
    oglos(NAGLOWEK, 'Etap potrzebuje nazwy wpisanej w polu.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.AppsStageSave, { windowId: idOkna, name: nazwa });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu etapu.', 'ostrzezenie');
    return;
  }
  odswiez();
}

async function dolozKamien(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  const nazwa = pole(korzen);
  if (nazwa === '') {
    oglos(NAGLOWEK, 'Kamień milowy potrzebuje nazwy wpisanej w polu.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.AppsMilestoneSave, { windowId: idOkna, name: nazwa });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu kamienia.', 'ostrzezenie');
    return;
  }
  odswiez();
}

/* Kamień wskazuje nazwa z pola: okno nie prowadzi wyboru wiersza, a wykaz
   kamieni stoi w panelu zadań tylko do odczytu. */
async function usunKamien(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  const nazwa = pole(korzen);
  const wykaz = await wywolaj(kanal, Command.AppsMilestoneList, { windowId: idOkna });
  const cel = wykaz.wynik?.milestones.find((kamien) => kamien.name === nazwa)?.id ?? '';
  if (cel === '') {
    oglos(NAGLOWEK, 'Wykaz kamieni nie ma pozycji o tej nazwie.', 'ostrzezenie');
    return;
  }
  if (!potwierdzone(korzen, 'data-apps', 'kamien-usun')) return;
  const wynik = await wywolaj(kanal, Command.AppsMilestoneDelete, {
    windowId: idOkna,
    milestoneId: cel,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił usunięcia kamienia.', 'ostrzezenie');
    return;
  }
  odswiez();
}

/* Wzorzec architektury bierze się z pola; bez wskazania staje monolit, bo to
   układ, w którym aplikacja zaczyna się najczęściej. */
async function zalozArchitekture(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  const [nazwa = '', wzorzec = ''] = czesci(korzen);
  if (nazwa === '') {
    oglos(NAGLOWEK,
      'Architektura potrzebuje nazwy, a po znaku pionowym wzorca: monolith, microservices albo serverless.',
      'ostrzezenie');
    return;
  }
  const wzorce: Record<string, AppArchitectureTemplate> = {
    monolith: AppArchitectureTemplate.Monolith,
    microservices: AppArchitectureTemplate.Microservices,
    serverless: AppArchitectureTemplate.Serverless,
  };
  const wynik = await wywolaj(kanal, Command.AppsArchitectureDefine, {
    windowId: idOkna,
    name: nazwa,
    template: wzorce[wzorzec] ?? AppArchitectureTemplate.Monolith,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił założenia architektury.', 'ostrzezenie');
    return;
  }
  odswiez();
}

async function odczytajArchitekture(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const wynik = await wywolaj(kanal, Command.AppsArchitectureGet, { windowId: idOkna });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu architektury.', 'ostrzezenie');
    return;
  }
  const architektura = wynik.wynik.architecture;
  if (architektura === undefined) {
    oglos(NAGLOWEK, 'To okno nie ma jeszcze założonej architektury.', 'ostrzezenie');
    return;
  }
  wypelnij(korzen, 'panel-architektura',
    (architektura.components ?? []).map((skladnik) => `${skladnik.name} · ${skladnik.kind}`));
  oglos(NAGLOWEK, `Architektura „${architektura.name ?? 'bez nazwy'}" na wzorcu ${architektura.template}.`);
}

async function sprawdzArchitekture(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const wynik = await wywolaj(kanal, Command.AppsArchitectureValidate, { windowId: idOkna });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił sprawdzenia architektury.', 'ostrzezenie');
    return;
  }
  const uwagi = wynik.wynik.issues;
  if (uwagi.length === 0) {
    oglos(NAGLOWEK, 'Architektura bez uwag.');
    return;
  }
  wypelnij(korzen, 'panel-architektura',
    uwagi.map((uwaga: AppValidationIssue) => `${uwaga.severity} · ${uwaga.message}`));
  oglos(NAGLOWEK, `Architektura ma ${uwagi.length} uwag.`, 'ostrzezenie');
}

/* Wydanie wraca oznaczeniem wytworu magazynu, więc okno nazywa oznaczenie
   zamiast udawać zapis pliku, którego nie ma gdzie zostawić. */
async function wydajArchitekture(kanal: Kanal, idOkna: string): Promise<void> {
  const wynik = await wywolaj(kanal, Command.AppsArchitectureExport, {
    windowId: idOkna,
    format: AppExportFormat.Mermaid,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wydania architektury.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Architektura wydana jako wytwór ${wynik.wynik.artifactRef}.`);
}

async function dolozPrzypis(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const tresc = pole(korzen);
  if (tresc === '') {
    oglos(NAGLOWEK, 'Przypis potrzebuje treści wpisanej w polu.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.AppsArchitectureAnnotationSave, {
    windowId: idOkna,
    text: tresc,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu przypisu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Przypis dołożony do architektury.');
}

async function odczytajMotyw(kanal: Kanal, idOkna: string): Promise<void> {
  const wynik = await wywolaj(kanal, Command.AppsThemeGet, { windowId: idOkna });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu motywu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, wynik.wynik.theme === undefined
    ? 'To okno nie ma jeszcze ustawionego motywu.'
    : JSON.stringify(wynik.wynik.theme).slice(0, 300));
}

/* Motyw jest dla rdzenia treścią nieprzejrzystą, więc pole podaje go zapisem
   JSON; zapis niepoprawny odmawia zamiast wysyłać rdzeniowi tekst. */
async function ustawMotyw(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const zapis = pole(korzen);
  if (zapis === '') {
    oglos(NAGLOWEK, 'Motyw potrzebuje zapisu JSON wpisanego w polu.', 'ostrzezenie');
    return;
  }
  let motyw: unknown;
  try {
    motyw = JSON.parse(zapis);
  } catch {
    oglos(NAGLOWEK, 'Wpisana treść nie jest poprawnym zapisem JSON.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.AppsThemeSet, { windowId: idOkna, theme: motyw });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił ustawienia motywu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Motyw ustawiony.');
}

async function uruchomPodglad(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const wynik = await wywolaj(kanal, Command.AppsPreviewStart, {
    windowId: idOkna,
    layer: AppWorkspaceLayer.Frontend,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił uruchomienia podglądu.', 'ostrzezenie');
    return;
  }
  const poleWpisu = korzen.querySelector<HTMLInputElement>('[data-apps-wpis]');
  if (poleWpisu !== null) poleWpisu.value = wynik.wynik.previewUrl;
  oglos(NAGLOWEK, `Podgląd stoi pod ${wynik.wynik.previewUrl}; adres wpisany w pole.`);
}

async function zatrzymajPodglad(kanal: Kanal, idOkna: string): Promise<void> {
  const wynik = await wywolaj(kanal, Command.AppsPreviewStop, { windowId: idOkna });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zatrzymania podglądu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Podgląd zatrzymany.');
}

async function zbadajWydajnosc(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const adres = pole(korzen);
  if (adres === '') {
    oglos(NAGLOWEK, 'Badanie wydajności potrzebuje adresu wpisanego w polu.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.AppsPerformanceAudit, { windowId: idOkna, url: adres });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił badania wydajności.', 'ostrzezenie');
    return;
  }
  const badanie = wynik.wynik.audit;
  oglos(NAGLOWEK,
    `Wydajność ${badanie.url}: ocena ${badanie.performanceScore}, miar ${badanie.metrics.length}.`);
}

/* Metoda i ścieżka stoją w jednym polu rozdzielone spacją; metoda nieznana
   odmawia, bo wysłanie sondy inną metodą byłoby decyzją za Operatora. */
async function sondujPunktStyku(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const [oznaczenie = '', sciezka = ''] = pole(korzen).split(/\s+/);
  const metody: Record<string, AppEndpointMethod> = {
    get: AppEndpointMethod.Get,
    post: AppEndpointMethod.Post,
    put: AppEndpointMethod.Put,
    patch: AppEndpointMethod.Patch,
    delete: AppEndpointMethod.Delete,
    graphql: AppEndpointMethod.Graphql,
  };
  const metoda = metody[oznaczenie.toLowerCase()];
  if (metoda === undefined || sciezka === '') {
    oglos(NAGLOWEK, 'Sonda potrzebuje metody i ścieżki, na przykład „GET /zdrowie".', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.AppsEndpointProbe, {
    windowId: idOkna,
    method: metoda,
    path: sciezka,
    environment: AppDeployEnvironment.Dev,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wysłania sondy.', 'ostrzezenie');
    return;
  }
  const odpowiedz = wynik.wynik.result;
  oglos(NAGLOWEK, `Sonda ${oznaczenie.toUpperCase()} ${sciezka}: kod ${odpowiedz.statusCode}, `
    + `czas ${odpowiedz.durationMs} ms.`);
}
