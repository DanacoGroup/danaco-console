// Wytwory przeglądania: zakładki, lista czytania, obserwacje stron, kanały
// treści, pobrania i przestrzenie kart. Pozycje pochodzą z wykazów rdzenia.

import {
  Command,
  EventType,
  type BrowserBookmark,
  type BrowserDownload,
  type BrowserFeed,
  type BrowserMonitor,
  type BrowserWorkspace,
} from '../../../shared/contract.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import { zglosUchwyt } from '../polaczenie/rozdzielacz-zdarzen.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { dolozPrzycisk, odmowa } from './browser-wspolne.ts';

export interface DzialanieSzyny {
  etykieta: string;
  cecha: string;
  wartosc: string;
}

export interface PozycjaSzyny {
  tytul: string;
  cechy: Record<string, string>;
  dzialania?: DzialanieSzyny[];
}

export interface WytworyPrzegladania {
  odlaczenia: Odsubskrybuj[];
  odswiezZakladki: () => void;
  odswiezObserwacje: () => void;
  odswiezKanaly: () => void;
  odswiezPrzestrzenie: () => void;
  odswiezPobrania: () => void;
}

export function szynaModulu(korzen: Element): Element | null {
  return korzen.querySelector('.dn-szyna-modulu-lista');
}

export async function zwiazWytworyPrzegladania(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  przy: AddEventListenerOptions,
): Promise<WytworyPrzegladania> {
  const szyna = szynaModulu(korzen);
  if (szyna !== null) zdejmijPozycjeWzorcowe(szyna);

  const odswiezZakladki = (): void => void postawZakladki(kanal, szyna, idOkna);
  const odswiezObserwacje = (): void => void postawObserwacje(kanal, szyna, idOkna);
  const odswiezKanaly = (): void => void postawKanaly(kanal, szyna, idOkna);
  const odswiezPrzestrzenie = (): void => void postawPrzestrzenie(kanal, szyna);
  const odswiezPobrania = (): void => void wypelnijPobrania(kanal, korzen, idOkna);

  await Promise.all([
    postawZakladki(kanal, szyna, idOkna),
    postawObserwacje(kanal, szyna, idOkna),
    postawKanaly(kanal, szyna, idOkna),
    postawPrzestrzenie(kanal, szyna),
  ]);
  odswiezPobrania();

  korzen.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;

    const zakladka = cel.closest<HTMLElement>('[data-zakladka-zdejmij]');
    if (zakladka !== null) {
      void zdejmijZakladke(kanal, szyna, idOkna, zakladka.dataset.zakladkaZdejmij ?? '');
      return;
    }
    const obserwacja = cel.closest<HTMLElement>('[data-obserwacja-sprawdz]');
    if (obserwacja !== null) {
      void sprawdzObserwacje(kanal, obserwacja.dataset.obserwacjaSprawdz ?? '');
      return;
    }
    const kanalTresci = cel.closest<HTMLElement>('[data-kanal-zdejmij]');
    if (kanalTresci !== null) {
      void zdejmijKanal(kanal, szyna, idOkna, kanalTresci.dataset.kanalZdejmij ?? '');
      return;
    }
    const przestrzen = cel.closest<HTMLElement>('[data-przestrzen-otworz]');
    if (przestrzen !== null) {
      void otworzPrzestrzen(kanal, idOkna, przestrzen.dataset.przestrzenOtworz ?? '');
    }
  }, przy);

  const odlaczenia = [
    zglosUchwyt(EventType.BrowserDownloadChanged, odswiezPobrania),
    zglosUchwyt(EventType.BrowserMonitorChanged, odswiezObserwacje),
  ];

  return {
    odlaczenia,
    odswiezZakladki,
    odswiezObserwacje,
    odswiezKanaly,
    odswiezPrzestrzenie,
    odswiezPobrania,
  };
}

/* Znacznik niesie nazwy zakładek i obserwacji wpisane wprost; pozycje wzorcowe
   schodzą, zanim rdzeń poda własne. */
function zdejmijPozycjeWzorcowe(szyna: Element): void {
  for (const pozycja of szyna.querySelectorAll('.pt-pozycja')) pozycja.remove();
}

async function postawZakladki(kanal: Kanal, szyna: Element | null, idOkna: string): Promise<void> {
  if (szyna === null) return;
  const wynik = await wywolaj(kanal, Command.BrowserBookmarkList, { windowId: idOkna });
  if (!wynik.udany) return;
  const zakladki = (wynik.wynik as { bookmarks?: BrowserBookmark[] } | undefined)?.bookmarks ?? [];
  postawGrupeSzyny(szyna, 'zakladki', 'Zakładki', zakladki.map((zakladka) => ({
    tytul: zakladka.title ?? zakladka.url,
    cechy: { zakladkaZdejmij: zakladka.id },
  })));
}

async function postawObserwacje(kanal: Kanal, szyna: Element | null, idOkna: string): Promise<void> {
  if (szyna === null) return;
  const wynik = await wywolaj(kanal, Command.BrowserMonitorList, { windowId: idOkna });
  if (!wynik.udany) return;
  const obserwacje = (wynik.wynik as { monitors?: BrowserMonitor[] } | undefined)?.monitors ?? [];
  postawGrupeSzyny(szyna, 'obserwacje', 'Obserwacje stron', obserwacje.map((obserwacja) => ({
    tytul: `${obserwacja.url} — ${obserwacja.status}`,
    cechy: { obserwacjaSprawdz: obserwacja.id },
    dzialania: [{ etykieta: 'Zdejmij', cecha: 'obserwacjaZdejmij', wartosc: obserwacja.id }],
  })));
}

async function postawKanaly(kanal: Kanal, szyna: Element | null, idOkna: string): Promise<void> {
  if (szyna === null) return;
  const wynik = await wywolaj(kanal, Command.BrowserFeedList, { windowId: idOkna });
  if (!wynik.udany) return;
  const kanaly = (wynik.wynik as { feeds?: BrowserFeed[] } | undefined)?.feeds ?? [];
  postawGrupeSzyny(szyna, 'kanaly', 'Kanały treści', kanaly.map((pozycja) => ({
    tytul: pozycja.title ?? pozycja.url,
    cechy: { kanalZdejmij: pozycja.id },
  })));
}

async function postawPrzestrzenie(kanal: Kanal, szyna: Element | null): Promise<void> {
  if (szyna === null) return;
  const wynik = await wywolaj(kanal, Command.BrowserWorkspaceList, {});
  if (!wynik.udany) return;
  const przestrzenie =
    (wynik.wynik as { workspaces?: BrowserWorkspace[] } | undefined)?.workspaces ?? [];
  postawGrupeSzyny(szyna, 'przestrzenie', 'Przestrzenie kart', przestrzenie.map((przestrzen) => ({
    tytul: `${przestrzen.name} — kart: ${przestrzen.tabCount}`,
    cechy: { przestrzenOtworz: przestrzen.id },
    dzialania: [{ etykieta: 'Usuń', cecha: 'przestrzenUsun', wartosc: przestrzen.id }],
  })));
}

async function wypelnijPobrania(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const gniazdo = korzen.querySelector('#panel-artefakty .sta-okno-tresc');
  if (gniazdo === null) return;
  const wynik = await wywolaj(kanal, Command.BrowserDownloadList, { windowId: idOkna });
  if (!wynik.udany) return;
  const pobrania = (wynik.wynik as { downloads?: BrowserDownload[] } | undefined)?.downloads ?? [];
  gniazdo.replaceChildren();
  if (pobrania.length === 0) {
    const pusty = document.createElement('p');
    pusty.className = 'dn-pusty';
    pusty.textContent = 'Żaden plik nie został pobrany.';
    gniazdo.append(pusty);
    return;
  }
  for (const pobranie of pobrania) {
    const wiersz = document.createElement('div');
    wiersz.className = 'br-artefakt';
    wiersz.dataset.pobranie = pobranie.id;
    const opis = document.createElement('span');
    opis.textContent = `${pobranie.fileName ?? pobranie.url} — ${pobranie.status}`;
    wiersz.append(opis);
    const rzad = document.createElement('div');
    rzad.className = 'sta-chip-rzad';
    for (const [etykieta, czynnosc] of CZYNNOSCI_POBRANIA) {
      dolozPrzycisk(rzad, etykieta, 'pobranieCzynnosc', czynnosc);
    }
    wiersz.append(rzad);
    gniazdo.append(wiersz);
  }
}

const CZYNNOSCI_POBRANIA: readonly (readonly [string, string])[] = [
  ['Wstrzymaj', 'pause'],
  ['Wznów', 'resume'],
  ['Przerwij', 'cancel'],
  ['Ponów', 'retry'],
  ['Usuń', 'remove'],
];

/* Grupa stoi we własnym pojemniku pod kluczem: odświeżenie jednej grupy
   wymienia jej zawartość zamiast dokładać drugi nagłówek obok pierwszego. */
export function postawGrupeSzyny(
  szyna: Element,
  klucz: string,
  etykieta: string,
  pozycje: PozycjaSzyny[],
): void {
  let pojemnik = szyna.querySelector<HTMLElement>(`[data-grupa-szyny="${klucz}"]`);
  if (pojemnik === null) {
    pojemnik = document.createElement('div');
    pojemnik.dataset.grupaSzyny = klucz;
    szyna.append(pojemnik);
  }
  pojemnik.replaceChildren();
  if (pozycje.length === 0) return;
  const naglowek = document.createElement('div');
  naglowek.className = 'dn-etyk-mono';
  naglowek.textContent = etykieta;
  pojemnik.append(naglowek);
  for (const pozycja of pozycje) {
    const przycisk = document.createElement('button');
    przycisk.className = 'pt-pozycja';
    przycisk.type = 'button';
    for (const [nazwa, wartosc] of Object.entries(pozycja.cechy)) {
      przycisk.dataset[nazwa] = wartosc;
    }
    const tytul = document.createElement('span');
    tytul.className = 'pt-pozycja-tytul';
    tytul.textContent = pozycja.tytul;
    przycisk.append(tytul);
    pojemnik.append(przycisk);
    if (pozycja.dzialania === undefined || pozycja.dzialania.length === 0) continue;
    const rzad = document.createElement('div');
    rzad.className = 'sta-chip-rzad';
    for (const dzialanie of pozycja.dzialania) {
      dolozPrzycisk(rzad, dzialanie.etykieta, dzialanie.cecha, dzialanie.wartosc);
    }
    pojemnik.append(rzad);
  }
}

async function zdejmijZakladke(
  kanal: Kanal,
  szyna: Element | null,
  idOkna: string,
  id: string,
): Promise<void> {
  if (id === '') return;
  const wynik = await wywolaj(kanal, Command.BrowserBookmarkRemove, { bookmarkId: id });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Zakładka nie została zdjęta.');
    return;
  }
  void postawZakladki(kanal, szyna, idOkna);
}

async function sprawdzObserwacje(kanal: Kanal, id: string): Promise<void> {
  if (id === '') return;
  const wynik = await wywolaj(kanal, Command.BrowserMonitorCheck, { monitorId: id });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Obserwacja nie została sprawdzona.');
  }
}

async function zdejmijKanal(
  kanal: Kanal,
  szyna: Element | null,
  idOkna: string,
  id: string,
): Promise<void> {
  if (id === '') return;
  const wynik = await wywolaj(kanal, Command.BrowserFeedRemove, { feedId: id });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Kanał nie został zdjęty.');
    return;
  }
  void postawKanaly(kanal, szyna, idOkna);
}

async function otworzPrzestrzen(kanal: Kanal, idOkna: string, id: string): Promise<void> {
  if (id === '') return;
  const wynik = await wywolaj(kanal, Command.BrowserWorkspaceOpen, {
    workspaceId: id,
    windowId: idOkna,
  });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Przestrzeń nie została otwarta.');
  }
}
