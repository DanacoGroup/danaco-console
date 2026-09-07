// Wiązanie karty modułu Browser z rdzeniem: karty przeglądania, źródła, notatki
// i artefakty. Węzły pochodzą ze znacznika Właściciela, treść z rdzenia.

import {
  Command,
  EventType,
  type BrowserNote,
  type BrowserSource,
  type BrowserTab,
} from '../../../shared/contract.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import { zglosUchwyt } from '../polaczenie/rozdzielacz-zdarzen.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { odmowa } from './browser-wspolne.ts';
import { dolozCzynnosciNotatki, zwiazNotatkiPrzegladania } from './browser-notatki.ts';
import { zwiazObserwacjeStron } from './browser-obserwacje.ts';
import { zwiazObszarIGranice } from './browser-obszar.ts';
import { dolozCzynnosciKarty, zwiazProwadzenieKarty } from './browser-prowadzenie.ts';
import { zwiazWgladTechniczny } from './browser-wglad.ts';
import { dolozCzynnosciZrodla, zwiazZbieranieTresci } from './browser-zbieranie.ts';
import { opiszNaglowek, zapewnijOknoModulu, zdejmijTrescWspolna } from './okno-modulu.ts';
import { zwiazKatalogModulu } from './katalog-modulu.ts';
import { zwiazWytworyPrzegladania } from './browser-wytwory.ts';

interface WiazaniePrzegladarki {
  korzen: Element;
  odlaczenia: Odsubskrybuj[];
}

const WIAZANIA = new Map<string, WiazaniePrzegladarki>();

export function zwiazPrzegladarke(
  kanal: Kanal,
  nazwaSrodowiska: string,
  idOknaStojacego: string,
  wskazanieKorzenia: Element | string,
): boolean {
  let idOkna = idOknaStojacego;
  const korzen = korzenKarty(wskazanieKorzenia);
  if (korzen === null) return false;
  const idKarty = korzen.getAttribute('data-karta') ?? '';
  if (idKarty === '') return false;
  if (WIAZANIA.get(idKarty)?.korzen === korzen) return false;
  zwolnijPrzegladarke(idKarty);

  const odlaczenia: Odsubskrybuj[] = [];
  const sterowanie = new AbortController();
  const przy = { signal: sterowanie.signal };
  odlaczenia.push(() => {
    sterowanie.abort();
  });
  WIAZANIA.set(idKarty, { korzen, odlaczenia });

  zdejmijTrescPrzykladowa(korzen);
  void opiszNaglowek(kanal, nazwaSrodowiska, korzen);

  /* Wykazy kart, źródeł i notatek rdzeń odmawia bez wskazania okna, a karta
     świeża okna jeszcze nie ma — okno powstaje więc tu, przed pytaniami. */
  void (async (): Promise<void> => {
    idOkna = await zapewnijOknoModulu(kanal, idKarty, 'browser', 'Browser', idOkna);
    if (idOkna === '') return;
    const odswiezKarty = (): void => void wypelnijKarty(kanal, korzen, idOkna);
    const odswiezZrodla = (): void => void wypelnijZrodla(kanal, korzen, idOkna);
    const odswiezNotatki = (): void => void wypelnijNotatki(kanal, korzen, idOkna);
    odswiezKarty();
    odswiezZrodla();
    odswiezNotatki();
    const wytwory = await zwiazWytworyPrzegladania(kanal, korzen, idOkna, przy);
    zwiazProwadzenieKarty(kanal, korzen, idOkna, przy, odswiezKarty);
    zwiazZbieranieTresci(kanal, korzen, idOkna, przy, {
      odswiezZakladki: wytwory.odswiezZakladki,
      odswiezKanaly: wytwory.odswiezKanaly,
      odswiezZrodla,
    });
    zwiazNotatkiPrzegladania(kanal, korzen, idOkna, przy, odswiezNotatki);
    zwiazObserwacjeStron(kanal, korzen, idOkna, przy, wytwory.odswiezObserwacje);
    zwiazWgladTechniczny(kanal, korzen, idOkna, przy);
    zwiazObszarIGranice(kanal, korzen, idOkna, przy, {
      odswiezPrzestrzenie: wytwory.odswiezPrzestrzenie,
      odswiezPobrania: wytwory.odswiezPobrania,
    });
    dopiszOdlaczenia(idKarty, korzen, wytwory.odlaczenia);
  })();

  korzen.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;

    const zamkniecie = cel.closest<HTMLElement>('[data-karta-zamknij]');
    if (zamkniecie !== null) {
      void zamknijKarte(kanal, korzen, idOkna, zamkniecie.dataset.kartaZamknij ?? '');
      return;
    }
    const zdjecie = cel.closest<HTMLElement>('[data-zrodlo-zdejmij]');
    if (zdjecie !== null) {
      void zdejmijZrodlo(kanal, korzen, idOkna, zdjecie.dataset.zrodloZdejmij ?? '');
    }
  }, przy);

  const adres = korzen.querySelector<HTMLInputElement>('[data-adres-przegladania]');
  adres?.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key !== 'Enter' || zdarzenie.shiftKey) return;
    zdarzenie.preventDefault();
    void otworzAdres(kanal, korzen, idOkna, adres.value);
  }, przy);

  odlaczenia.push(zglosUchwyt(EventType.BrowserTabChanged, () => {
    void wypelnijKarty(kanal, korzen, idOkna);
  }));
  odlaczenia.push(zglosUchwyt(EventType.BrowserPageChanged, () => {
    void wypelnijZrodla(kanal, korzen, idOkna);
  }));

  const katalog = zwiazKatalogModulu(kanal, idOkna, korzen, 'browser', 'Przeglądarka');
  if (katalog !== null) odlaczenia.push(katalog);

  return true;
}

export function zwolnijPrzegladarke(idKarty: string): void {
  const wiazanie = WIAZANIA.get(idKarty);
  if (wiazanie === undefined) return;
  for (const odlacz of wiazanie.odlaczenia) odlacz();
  WIAZANIA.delete(idKarty);
}

/* Wiązania podrzędne powstają po oczekiwaniu na okno, więc karta mogła w tym
   czasie zejść: subskrypcja spóźniona schodzi od razu zamiast zostać w tle. */
function dopiszOdlaczenia(idKarty: string, korzen: Element, nowe: Odsubskrybuj[]): void {
  const wiazanie = WIAZANIA.get(idKarty);
  if (wiazanie === undefined || wiazanie.korzen !== korzen) {
    for (const odlacz of nowe) odlacz();
    return;
  }
  wiazanie.odlaczenia.push(...nowe);
}

function korzenKarty(wskazanie: Element | string): Element | null {
  if (typeof wskazanie !== 'string') return wskazanie;
  return document.querySelector(`.cd-tresc--modul[data-karta="${wskazanie}"]`);
}

/* Znacznik niesie karty, źródła i notatki wpisane wprost. Wypełniacz zostaje
   w produkcie na zawsze, więc schodzi przed pierwszym pytaniem rdzenia. */
function zdejmijTrescPrzykladowa(korzen: Element): void {
  for (const wybor of ['.br-karta', '.br-zrodlo', '.br-notatka', '.br-artefakt', '.brw-tab']) {
    for (const wezel of korzen.querySelectorAll(wybor)) wezel.remove();
  }
  // Rozmowa, żetony kontekstu i znacznik pracy też są wpisane wprost w prototyp.
  if (korzen instanceof HTMLElement) zdejmijTrescWspolna(korzen);
  /* Pasek adresu niesie adres wzorcowy; czynności biorą adres wprost z pola,
     więc wpis zostawiony pracowałby na stronie, której nikt nie otwierał. */
  const adres = korzen.querySelector<HTMLInputElement>('[data-adres-przegladania]');
  if (adres !== null) adres.value = '';
  /* Wskaźnik obecności twierdzi, że model patrzy na stronę. Rdzeń takiego
     stanu nie podaje, więc plakietka mówiłaby to bez pokrycia. */
  for (const wezel of korzen.querySelectorAll('[title="Wskaźnik obecności AI"]')) {
    wezel.remove();
  }
  /* Pliki i zadania w tle to wyposażenie wspólne okna, którego przeglądanie nie
     wypełnia — panel zadań niesie w prototypie licznik wzięty znikąd. Oba stoją
     ukryte, więc przejście po oknach ich nie widzi; zdanie wchodzi mimo to. */
  for (const identyfikator of ['panel-pliki', 'panel-zadania']) {
    const cialo = korzen.querySelector(`#${identyfikator} .sta-okno-tresc`);
    if (cialo === null) continue;
    postawStanPusty(cialo, 'Rdzeń nie podaje tego wykazu dla okna przeglądania.');
  }
}

async function wypelnijKarty(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const gniazdo = korzen.querySelector('#panel-browser .sta-okno-tresc');
  if (gniazdo === null) return;
  const wynik = await wywolaj(kanal, Command.BrowserTabList, { windowId: idOkna });
  if (!wynik.udany) {
    postawStanPusty(gniazdo, 'Wykaz kart przeglądania nie doszedł.');
    return;
  }
  const karty = (wynik.wynik as { tabs?: BrowserTab[] } | undefined)?.tabs ?? [];
  gniazdo.replaceChildren();
  if (karty.length === 0) {
    postawStanPusty(gniazdo, 'Żadna karta przeglądania nie stoi.');
    return;
  }
  for (const karta of karty) {
    const wiersz = document.createElement('div');
    wiersz.className = 'br-karta';
    const tytul = document.createElement('span');
    tytul.textContent = karta.title ?? karta.url ?? '';
    wiersz.append(tytul);
    dolozCzynnosciKarty(wiersz, karta);
    const zamknij = document.createElement('button');
    zamknij.className = 'dn-btn dn-btn--zarys dn-btn--sm';
    zamknij.type = 'button';
    zamknij.dataset.kartaZamknij = karta.id;
    zamknij.textContent = 'Zamknij';
    wiersz.append(zamknij);
    gniazdo.append(wiersz);
  }
}

async function wypelnijZrodla(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const gniazdo = korzen.querySelector('#panel-sources .sta-okno-tresc');
  if (gniazdo === null) return;
  const wynik = await wywolaj(kanal, Command.BrowserSourceList, { windowId: idOkna });
  if (!wynik.udany) {
    postawStanPusty(gniazdo, 'Wykaz źródeł nie doszedł.');
    return;
  }
  const zrodla = (wynik.wynik as { sources?: BrowserSource[] } | undefined)?.sources ?? [];
  gniazdo.replaceChildren();
  if (zrodla.length === 0) {
    postawStanPusty(gniazdo, 'Żadne źródło nie zostało wniesione.');
    return;
  }
  for (const zrodlo of zrodla) {
    const wiersz = document.createElement('div');
    wiersz.className = 'br-zrodlo';
    const tytul = document.createElement('span');
    tytul.textContent = zrodlo.title ?? zrodlo.url;
    wiersz.append(tytul);
    dolozCzynnosciZrodla(wiersz, zrodlo);
    const zdejmij = document.createElement('button');
    zdejmij.className = 'dn-btn dn-btn--zarys dn-btn--sm';
    zdejmij.type = 'button';
    zdejmij.dataset.zrodloZdejmij = zrodlo.id;
    zdejmij.textContent = 'Zdejmij';
    wiersz.append(zdejmij);
    gniazdo.append(wiersz);
  }
}

async function wypelnijNotatki(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const gniazdo = korzen.querySelector('#panel-notes .sta-okno-tresc');
  if (gniazdo === null) return;
  const wynik = await wywolaj(kanal, Command.BrowserNoteList, { windowId: idOkna });
  if (!wynik.udany) {
    postawStanPusty(gniazdo, 'Wykaz notatek nie doszedł.');
    return;
  }
  const notatki = (wynik.wynik as { notes?: BrowserNote[] } | undefined)?.notes ?? [];
  gniazdo.replaceChildren();
  if (notatki.length === 0) {
    postawStanPusty(gniazdo, 'Żadna notatka nie powstała.');
    return;
  }
  for (const notatka of notatki) {
    const wiersz = document.createElement('div');
    wiersz.className = 'br-notatka';
    dolozCzynnosciNotatki(wiersz, notatka);
    gniazdo.append(wiersz);
  }
}

async function otworzAdres(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  adres: string,
): Promise<void> {
  const url = adres.trim();
  if (url === '') return;
  const wynik = await wywolaj(kanal, Command.BrowserTabOpen, { windowId: idOkna, url });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Karta nie została otwarta.');
    return;
  }
  void wypelnijKarty(kanal, korzen, idOkna);
}

async function zamknijKarte(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  idKarty: string,
): Promise<void> {
  if (idKarty === '') return;
  const wynik = await wywolaj(kanal, Command.BrowserTabClose, { tabId: idKarty });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Karta nie została zamknięta.');
    return;
  }
  void wypelnijKarty(kanal, korzen, idOkna);
}

async function zdejmijZrodlo(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  idZrodla: string,
): Promise<void> {
  if (idZrodla === '') return;
  const wynik = await wywolaj(kanal, Command.BrowserSourceRemove, { sourceId: idZrodla });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Źródło nie zostało zdjęte.');
    return;
  }
  void wypelnijZrodla(kanal, korzen, idOkna);
}

function postawStanPusty(gniazdo: Element, zdanie: string): void {
  gniazdo.replaceChildren();
  const wezel = document.createElement('p');
  wezel.className = 'dn-pusty';
  wezel.textContent = zdanie;
  gniazdo.append(wezel);
}
