// Czynności zasobów okna Developer: pojemniki i obrazy, połączenia z bazą
// wraz z zapytaniem, schematem i przeniesieniami, oraz zbiory żądań.
import { Command, ContainerActionKind, DataEngine } from '../../../shared/contract.ts';
import type { ContainerInfo, DataSchemaNode } from '../../../shared/contract.ts';
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

const NAGLOWEK = 'Zasoby projektu';
const ZNACZNIK = 'devzasob';

const PRZYCISKI: ReadonlyArray<readonly [string, string]> = [
  ['pojemnik', 'Zatrzymaj pojemnik'],
  ['zestaw', 'Podnieś zestaw'],
  ['obraz', 'Zbuduj obraz'],
  ['polaczenie', 'Zapisz połączenie'],
  ['zapytanie', 'Wykonaj zapytanie'],
  ['schemat', 'Odczytaj schemat'],
  ['przeniesienia', 'Wykonaj przeniesienia'],
  ['zadanie', 'Wyślij żądanie'],
  ['zbior', 'Zapisz zbiór żądań'],
  ['openapi', 'Wciągnij opis OpenAPI'],
  ['obciazenie', 'Zbadaj obciążenie'],
];

export function zwiazZasobyDevelopera(
  kanal: Kanal,
  korzen: Element,
  idOkna: () => string,
  odswiez: () => void,
  przy: AddEventListenerOptions,
): void {
  zalozPas(korzen, 'panel-terminal', ZNACZNIK, 'Nazwa, adres albo zapytanie', PRZYCISKI);
  zwiazPas(korzen, ZNACZNIK, (czynnosc) => {
    void wykonaj(kanal, korzen, czynnosc, idOkna(), odswiez);
  }, przy);
}

function wpis(korzen: Element): string {
  return wpisPasa(korzen, ZNACZNIK);
}

function czesci(korzen: Element): string[] {
  return czesciWpisu(korzen, ZNACZNIK);
}

/* Zapytanie, schemat i przeniesienia idą przez połączenie stojące w wykazie
   najwyżej: okno nie prowadzi wyboru połączenia, a wykaz podaje kolejność. */
async function pierwszePolaczenie(kanal: Kanal, idOkna: string): Promise<string> {
  const wykaz = await wywolaj(kanal, Command.DeveloperDataConnectionList, { windowId: idOkna });
  const cel = wykaz.wynik?.connections[0]?.id ?? '';
  if (cel === '') {
    oglos(NAGLOWEK, 'To okno nie ma jeszcze zapisanego połączenia z bazą.', 'ostrzezenie');
  }
  return cel;
}

async function wykonaj(
  kanal: Kanal,
  korzen: Element,
  czynnosc: string,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  if (idOkna === '') {
    oglos(NAGLOWEK, 'Rdzeń nie dał okna projektu dla tej karty.', 'ostrzezenie');
    return;
  }
  if (czynnosc === 'pojemnik') return zatrzymajPojemnik(kanal, korzen, idOkna, odswiez);
  if (czynnosc === 'zestaw') return podniesZestaw(kanal, korzen, idOkna);
  if (czynnosc === 'obraz') return zbudujObraz(kanal, korzen, idOkna);
  if (czynnosc === 'polaczenie') return zapiszPolaczenie(kanal, korzen, idOkna, odswiez);
  if (czynnosc === 'zapytanie') return wykonajZapytanie(kanal, korzen, idOkna);
  if (czynnosc === 'schemat') return odczytajSchemat(kanal, korzen, idOkna);
  if (czynnosc === 'przeniesienia') return wykonajPrzeniesienia(kanal, korzen, idOkna);
  if (czynnosc === 'zadanie') return wyslijZadanie(kanal, korzen, idOkna);
  if (czynnosc === 'zbior') return zapiszZbior(kanal, korzen, idOkna, odswiez);
  if (czynnosc === 'openapi') return wciagnijOpenapi(kanal, korzen, idOkna, odswiez);
  if (czynnosc === 'obciazenie') return zbadajObciazenie(kanal, korzen, idOkna);
  /* Czynność spoza obsłużonych odmawia zamiast milczeć: przycisk bez gałęzi
     wyglądałby jak działający. */
  oglos(NAGLOWEK, `Czynność „${czynnosc}" nie jest prowadzona przez to okno.`, 'ostrzezenie');
}

/* Zatrzymanie pojemnika przerywa usługę, którą on niesie, więc pierwsze
   naciśnięcie uzbraja przycisk; wskazaniem jest nazwa pojemnika z pola. */
async function zatrzymajPojemnik(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  const nazwa = wpis(korzen);
  const wykaz = await wywolaj(kanal, Command.DeveloperContainerList, { windowId: idOkna });
  if (wykaz.wynik?.engineAvailable === false) {
    oglos(NAGLOWEK, 'Rdzeń nie widzi silnika pojemników na swojej maszynie.', 'ostrzezenie');
    return;
  }
  const cel = wykaz.wynik?.containers.find((pojemnik: ContainerInfo) => pojemnik.name === nazwa)?.id
    ?? '';
  if (cel === '') {
    oglos(NAGLOWEK, 'Wykaz pojemników nie ma pozycji o tej nazwie.', 'ostrzezenie');
    return;
  }
  if (!potwierdzone(korzen, ZNACZNIK, 'pojemnik')) return;
  const wynik = await wywolaj(kanal, Command.DeveloperContainerAction, {
    containerId: cel,
    action: ContainerActionKind.Stop,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zatrzymania pojemnika.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Pojemnik ${nazwa} w stanie ${wynik.wynik.container.status}.`);
  odswiez();
}

async function podniesZestaw(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const plik = wpis(korzen);
  if (plik === '') {
    oglos(NAGLOWEK, 'Podniesienie zestawu potrzebuje ścieżki pliku zestawu.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.DeveloperComposeUp, { windowId: idOkna, file: plik });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił podniesienia zestawu.', 'ostrzezenie');
    return;
  }
  wypelnijPanel(korzen, 'panel-terminal', wynik.wynik.services.map((usluga: ContainerInfo) =>
    `${usluga.status} ${usluga.name}`));
  oglos(NAGLOWEK, `Zestaw podniesiony: ${wynik.wynik.services.length} usług.`);
}

/* Obraz nie idzie do składu wprost: wysłanie go tam byłoby wystawieniem pracy
   na zewnątrz, o czym rozstrzyga Operator osobnym poleceniem. */
async function zbudujObraz(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const [plik = '', oznaczenie = ''] = czesci(korzen);
  if (plik === '' || oznaczenie === '') {
    oglos(NAGLOWEK,
      'Budowa obrazu potrzebuje pliku i oznaczenia, na przykład „Dockerfile | sklep:1.0".',
      'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.DeveloperImageBuild, {
    windowId: idOkna,
    dockerfile: plik,
    tag: oznaczenie,
    push: false,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zbudowania obrazu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Obraz ${oznaczenie} zbudowany pod oznaczeniem ${wynik.wynik.imageId}.`);
}

/* Połączenie zakłada się tylko do odczytu: prawo zapisu do cudzej bazy ma być
   nadane wprost, nie wzięte z domyślnego ustawienia okna. */
async function zapiszPolaczenie(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  const [nazwa = '', baza = ''] = czesci(korzen);
  if (nazwa === '' || baza === '') {
    oglos(NAGLOWEK,
      'Połączenie potrzebuje nazwy i bazy, na przykład „miejscowa | dane.db".',
      'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.DeveloperDataConnectionSet, {
    windowId: idOkna,
    name: nazwa,
    engine: DataEngine.Sqlite,
    database: baza,
    readOnly: true,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu połączenia.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Połączenie „${wynik.wynik.connection.name}" zapisane tylko do odczytu.`);
  odswiez();
}

async function wykonajZapytanie(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const zapytanie = wpis(korzen);
  if (zapytanie === '') {
    oglos(NAGLOWEK, 'Zapytanie potrzebuje treści wpisanej w polu.', 'ostrzezenie');
    return;
  }
  const cel = await pierwszePolaczenie(kanal, idOkna);
  if (cel === '') return;
  const wynik = await wywolaj(kanal, Command.DeveloperDataQueryRun, {
    connectionId: cel,
    sql: zapytanie,
    limit: 200,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wykonania zapytania.', 'ostrzezenie');
    return;
  }
  const odpowiedz = wynik.wynik.result;
  wypelnijPanel(korzen, 'panel-zadania', [odpowiedz.columns.join(' · ')]);
  oglos(NAGLOWEK, odpowiedz.truncated
    ? `Zapytanie oddało ${odpowiedz.rowCount} wierszy, przycięte granicą odczytu.`
    : `Zapytanie oddało ${odpowiedz.rowCount} wierszy w ${odpowiedz.durationMs} ms.`,
  odpowiedz.truncated ? 'ostrzezenie' : 'informacja');
}

async function odczytajSchemat(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const cel = await pierwszePolaczenie(kanal, idOkna);
  if (cel === '') return;
  const wynik = await wywolaj(kanal, Command.DeveloperDataSchemaGet, { connectionId: cel });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu schematu.', 'ostrzezenie');
    return;
  }
  const wezly = wynik.wynik.nodes;
  wypelnijPanel(korzen, 'panel-zadania', wezly.map((wezel: DataSchemaNode) =>
    `${wezel.kind} ${wezel.name}`));
  if (wezly.length === 0) oglos(NAGLOWEK, 'Ta baza nie ma jeszcze żadnych tablic.');
}

/* Przeniesienia idą najpierw na sucho: wykaz zastosowanych i czekających
   wraca bez tknięcia bazy, a dopiero drugie naciśnięcie je wykonuje. */
async function wykonajPrzeniesienia(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const cel = await pierwszePolaczenie(kanal, idOkna);
  if (cel === '') return;
  const naSucho = !potwierdzone(korzen, ZNACZNIK, 'przeniesienia');
  const wynik = await wywolaj(kanal, Command.DeveloperDataMigrationRun, {
    connectionId: cel,
    windowId: idOkna,
    direction: 'up',
    dryRun: naSucho,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił przeniesień.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, naSucho
    ? `Na sucho: ${wynik.wynik.pending.length} przeniesień czeka. Naciśnij ponownie, aby wykonać.`
    : `Wykonano ${wynik.wynik.applied.length} przeniesień; czeka ${wynik.wynik.pending.length}.`);
}

/* Sposób i adres stoją w jednym polu rozdzielone spacją; bez sposobu okno
   bierze odczyt, bo tylko on niczego po drugiej stronie nie zmienia. */
async function wyslijZadanie(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const czesciWpisane = wpis(korzen).split(/\s+/);
  const sposob = czesciWpisane.length > 1 ? (czesciWpisane[0] ?? '') : 'GET';
  const adres = czesciWpisane.length > 1 ? (czesciWpisane[1] ?? '') : (czesciWpisane[0] ?? '');
  if (adres === '') {
    oglos(NAGLOWEK, 'Żądanie potrzebuje adresu, na przykład „GET https://przyklad.pl".',
      'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.DeveloperApiRequest, {
    windowId: idOkna,
    method: sposob.toUpperCase(),
    url: adres,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wysłania żądania.', 'ostrzezenie');
    return;
  }
  const odpowiedz = wynik.wynik.response;
  wypelnijPanel(korzen, 'panel-artefakty', (odpowiedz.body ?? '').split('\n'));
  oglos(NAGLOWEK, `Odpowiedź ${odpowiedz.status} w ${odpowiedz.durationMs} ms.`);
}

async function zapiszZbior(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  const nazwa = wpis(korzen);
  if (nazwa === '') {
    oglos(NAGLOWEK, 'Zbiór żądań potrzebuje nazwy wpisanej w polu.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.DeveloperApiCollectionSave, {
    windowId: idOkna,
    name: nazwa,
    requests: [],
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu zbioru.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Zbiór „${wynik.wynik.collection.name}" zapisany.`);
  odswiez();
}

async function wciagnijOpenapi(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  const wskazanie = wpis(korzen);
  if (wskazanie === '') {
    oglos(NAGLOWEK, 'Wciągnięcie opisu potrzebuje ścieżki albo adresu wpisanego w polu.',
      'ostrzezenie');
    return;
  }
  const przezSiec = wskazanie.startsWith('http://') || wskazanie.startsWith('https://');
  const wynik = await wywolaj(kanal, Command.DeveloperApiOpenapiImport, {
    windowId: idOkna,
    ...(przezSiec ? { url: wskazanie } : { path: wskazanie }),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wciągnięcia opisu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Zbiór „${wynik.wynik.collection.name}" wciągnięty: `
    + `${wynik.wynik.requestCount} żądań.`);
  odswiez();
}

/* Badanie obciążenia wysyła cudzemu serwerowi wiele żądań naraz, więc pierwsze
   naciśnięcie uzbraja przycisk. */
async function zbadajObciazenie(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const adres = wpis(korzen);
  if (adres === '') {
    oglos(NAGLOWEK, 'Badanie obciążenia potrzebuje adresu wpisanego w polu.', 'ostrzezenie');
    return;
  }
  if (!potwierdzone(korzen, ZNACZNIK, 'obciazenie')) return;
  const wynik = await wywolaj(kanal, Command.DeveloperApiLoadRun, {
    windowId: idOkna,
    url: adres,
    connections: 10,
    durationSeconds: 10,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił badania obciążenia.', 'ostrzezenie');
    return;
  }
  const bieg = wynik.wynik.run;
  oglos(NAGLOWEK, `${bieg.requestsTotal} żądań, ${bieg.requestsPerSecond} na sekundę, `
    + `poza zakresem 2xx: ${bieg.non2xx}, błędów ${bieg.errors}.`);
}
