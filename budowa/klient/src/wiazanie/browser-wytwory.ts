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
import { zglosUchwyt } from '../polaczenie/rozdzielacz-zdarzen.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';

export async function zwiazWytworyPrzegladania(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  przy: AddEventListenerOptions,
): Promise<void> {
  const szyna = korzen.querySelector('.dn-szyna-modulu-lista');
  if (szyna !== null) zdejmijPozycjeWzorcowe(szyna);

  await Promise.all([
    postawZakladki(kanal, szyna, idOkna),
    postawObserwacje(kanal, szyna, idOkna),
    postawKanaly(kanal, szyna, idOkna),
    postawPrzestrzenie(kanal, szyna),
  ]);
  void wypelnijPobrania(kanal, korzen, idOkna);

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

  zglosUchwyt(EventType.BrowserDownloadChanged, () => {
    void wypelnijPobrania(kanal, korzen, idOkna);
  });
  zglosUchwyt(EventType.BrowserMonitorChanged, () => {
    void postawObserwacje(kanal, szyna, idOkna);
  });
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
  postawGrupe(szyna, 'Zakładki', zakladki.map((zakladka) => ({
    tytul: zakladka.title ?? zakladka.url,
    cechy: { zakladkaZdejmij: zakladka.id },
  })));
}

async function postawObserwacje(kanal: Kanal, szyna: Element | null, idOkna: string): Promise<void> {
  if (szyna === null) return;
  const wynik = await wywolaj(kanal, Command.BrowserMonitorList, { windowId: idOkna });
  if (!wynik.udany) return;
  const obserwacje = (wynik.wynik as { monitors?: BrowserMonitor[] } | undefined)?.monitors ?? [];
  postawGrupe(szyna, 'Obserwacje stron', obserwacje.map((obserwacja) => ({
    tytul: obserwacja.url,
    cechy: { obserwacjaSprawdz: obserwacja.id },
  })));
}

async function postawKanaly(kanal: Kanal, szyna: Element | null, idOkna: string): Promise<void> {
  if (szyna === null) return;
  const wynik = await wywolaj(kanal, Command.BrowserFeedList, { windowId: idOkna });
  if (!wynik.udany) return;
  const kanaly = (wynik.wynik as { feeds?: BrowserFeed[] } | undefined)?.feeds ?? [];
  postawGrupe(szyna, 'Kanały treści', kanaly.map((pozycja) => ({
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
  postawGrupe(szyna, 'Przestrzenie kart', przestrzenie.map((przestrzen) => ({
    tytul: przestrzen.name,
    cechy: { przestrzenOtworz: przestrzen.id },
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
    wiersz.textContent = `${pobranie.fileName ?? pobranie.url} — ${pobranie.status}`;
    gniazdo.append(wiersz);
  }
}

interface PozycjaSzyny {
  tytul: string;
  cechy: Record<string, string>;
}

function postawGrupe(szyna: Element, etykieta: string, pozycje: PozycjaSzyny[]): void {
  if (pozycje.length === 0) return;
  const naglowek = document.createElement('div');
  naglowek.className = 'dn-etyk-mono';
  naglowek.textContent = etykieta;
  szyna.append(naglowek);
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
    szyna.append(przycisk);
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
    oglos('Przeglądanie', wynik.blad?.message ?? 'Zakładka nie została zdjęta.', 'blad');
    return;
  }
  void postawZakladki(kanal, szyna, idOkna);
}

async function sprawdzObserwacje(kanal: Kanal, id: string): Promise<void> {
  if (id === '') return;
  const wynik = await wywolaj(kanal, Command.BrowserMonitorCheck, { monitorId: id });
  if (!wynik.udany) {
    oglos('Przeglądanie', wynik.blad?.message ?? 'Obserwacja nie została sprawdzona.', 'blad');
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
    oglos('Przeglądanie', wynik.blad?.message ?? 'Kanał nie został zdjęty.', 'blad');
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
    oglos('Przeglądanie', wynik.blad?.message ?? 'Przestrzeń nie została otwarta.', 'blad');
  }
}
