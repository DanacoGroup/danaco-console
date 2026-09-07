// Zbieranie treści z przeglądanej strony: zakładka, lista do czytania, źródła
// wraz z grupami oraz subskrypcja kanału. Adres bierze się z paska adresu.

import { Command, type BrowserReadingItem, type BrowserSource } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import {
  adresWymagany,
  dolozPrzycisk,
  odmowa,
  pasNadCialem,
  polePasa,
  powiedz,
  zalozPas,
} from './browser-wspolne.ts';
import { postawGrupeSzyny, szynaModulu } from './browser-wytwory.ts';

export interface OdswiezeniaZbierania {
  odswiezZakladki: () => void;
  odswiezKanaly: () => void;
  odswiezZrodla: () => void;
}

export function zwiazZbieranieTresci(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  przy: AddEventListenerOptions,
  odswiezenia: OdswiezeniaZbierania,
): void {
  const szyna = szynaModulu(korzen);
  postawPasZbierania(korzen);
  const poleGrupy = postawPasZrodel(korzen);
  const odswiezCzytanie = (): void => void postawCzytanie(kanal, szyna, idOkna);
  odswiezCzytanie();
  void postawGrupyZrodel(kanal, szyna, idOkna);

  korzen.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;

    if (cel.closest('[data-zakladka-dodaj]') !== null) {
      void dodajZakladke(kanal, korzen, idOkna, odswiezenia.odswiezZakladki);
      return;
    }
    if (cel.closest('[data-czytanie-dodaj]') !== null) {
      void dodajDoCzytania(kanal, korzen, idOkna, odswiezCzytanie);
      return;
    }
    const zdjecie = cel.closest<HTMLElement>('[data-czytanie-zdejmij]');
    if (zdjecie !== null) {
      void zdejmijZCzytania(kanal, zdjecie.dataset.czytanieZdejmij ?? '', odswiezCzytanie);
      return;
    }
    if (cel.closest('[data-kanal-subskrybuj]') !== null) {
      void subskrybujKanal(kanal, korzen, idOkna, odswiezenia.odswiezKanaly);
      return;
    }
    if (cel.closest('[data-zrodlo-wnies]') !== null) {
      void wniesZrodlo(kanal, korzen, idOkna, odswiezenia.odswiezZrodla);
      return;
    }
    const doGrupy = cel.closest<HTMLElement>('[data-zrodlo-do-grupy]');
    if (doGrupy !== null) {
      void wstawZrodloDoGrupy(kanal, szyna, idOkna, doGrupy.dataset.zrodloDoGrupy ?? '',
        (poleGrupy?.value ?? '').trim());
    }
  }, przy);
}

/** Dokłada w wierszu źródła drogę do grupy; zdjęcie źródła stoi już obok. */
export function dolozCzynnosciZrodla(wiersz: Element, zrodlo: BrowserSource): void {
  dolozPrzycisk(wiersz, 'Do grupy', 'zrodloDoGrupy', zrodlo.id);
}

function postawPasZbierania(korzen: Element): void {
  const panel = korzen.querySelector('#panel-browser');
  const plotno = panel?.querySelector('.brw-plotno') ?? null;
  if (panel === null || plotno === null) return;
  const pas = zalozPas('Zbieranie treści');
  dolozPrzycisk(pas, 'Do zakładek', 'zakladkaDodaj');
  dolozPrzycisk(pas, 'Do czytania', 'czytanieDodaj');
  dolozPrzycisk(pas, 'Wnieś źródło', 'zrodloWnies');
  dolozPrzycisk(pas, 'Subskrybuj kanał', 'kanalSubskrybuj');
  panel.insertBefore(pas, plotno);
}

function postawPasZrodel(korzen: Element): HTMLInputElement | null {
  const cialo = korzen.querySelector('#panel-sources .sta-okno-tresc');
  const pas = pasNadCialem(cialo, 'Grupy źródeł');
  return polePasa(pas, 'Nazwa grupy źródeł');
}

async function dodajZakladke(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  const url = adresWymagany(korzen);
  if (url === '') return;
  const wynik = await wywolaj(kanal, Command.BrowserBookmarkAdd, { windowId: idOkna, url });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił założenia zakładki.');
    return;
  }
  powiedz(`Zakładka „${wynik.wynik?.bookmark.title ?? url}” stoi w wykazie.`);
  odswiez();
}

async function dodajDoCzytania(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  const url = adresWymagany(korzen);
  if (url === '') return;
  const wynik = await wywolaj(kanal, Command.BrowserReadlistAdd, { windowId: idOkna, url });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił dołożenia pozycji do czytania.');
    return;
  }
  powiedz(`Pozycja „${wynik.wynik?.item.title ?? url}” czeka na przeczytanie.`);
  odswiez();
}

/* Zdjęcie z listy do czytania oznacza pozycję jako przeczytaną — pozycja nie
   znika bez śladu, a wykaz pyta wyłącznie o zaległe. */
async function zdejmijZCzytania(kanal: Kanal, id: string, odswiez: () => void): Promise<void> {
  if (id === '') return;
  const wynik = await wywolaj(kanal, Command.BrowserReadlistRemove, {
    itemId: id,
    markRead: true,
  });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił zdjęcia pozycji z listy do czytania.');
    return;
  }
  powiedz('Pozycja zeszła z listy do czytania.');
  odswiez();
}

async function postawCzytanie(kanal: Kanal, szyna: Element | null, idOkna: string): Promise<void> {
  if (szyna === null) return;
  const wynik = await wywolaj(kanal, Command.BrowserReadlistList, {
    windowId: idOkna,
    pendingOnly: true,
  });
  if (!wynik.udany) return;
  const pozycje = (wynik.wynik as { items?: BrowserReadingItem[] } | undefined)?.items ?? [];
  postawGrupeSzyny(szyna, 'czytanie', 'Do czytania', pozycje.map((pozycja) => ({
    tytul: pozycja.title ?? pozycja.url,
    cechy: { czytanieOtworz: pozycja.id },
    dzialania: [{ etykieta: 'Zdejmij', cecha: 'czytanieZdejmij', wartosc: pozycja.id }],
  })));
}

async function subskrybujKanal(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  const url = adresWymagany(korzen);
  if (url === '') return;
  const wynik = await wywolaj(kanal, Command.BrowserFeedSubscribe, { windowId: idOkna, url });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił subskrypcji kanału.');
    return;
  }
  powiedz(`Kanał „${wynik.wynik?.feed.title ?? url}” jest subskrybowany.`);
  odswiez();
}

async function wniesZrodlo(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  const url = adresWymagany(korzen);
  if (url === '') return;
  const wynik = await wywolaj(kanal, Command.BrowserSourceAdd, { windowId: idOkna, url });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił wniesienia źródła.');
    return;
  }
  powiedz(`Źródło „${wynik.wynik?.source.title ?? url}” stoi w wykazie.`);
  odswiez();
}

/* Kontrakt nie zna dokładania źródła do grupy: zapis idzie całą grupą, więc
   skład bierze się z wykazu grup i dochodzi do niego źródło wskazane. */
async function wstawZrodloDoGrupy(
  kanal: Kanal,
  szyna: Element | null,
  idOkna: string,
  idZrodla: string,
  nazwa: string,
): Promise<void> {
  if (idZrodla === '') return;
  if (nazwa === '') {
    odmowa(undefined, 'Grupa źródeł wymaga nazwy — pole nazwy grupy jest puste.');
    return;
  }
  const wykaz = await wywolaj(kanal, Command.BrowserSourceGroupList, { windowId: idOkna });
  const grupa = wykaz.wynik?.groups.find((pozycja) => pozycja.name === nazwa);
  const sklad = new Set(grupa?.sourceIds ?? []);
  sklad.add(idZrodla);
  const wynik = await wywolaj(kanal, Command.BrowserSourceGroupSet, {
    windowId: idOkna,
    groupId: grupa?.id,
    name: nazwa,
    sourceIds: [...sklad],
  });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił zapisu grupy źródeł.');
    return;
  }
  powiedz(`Źródło stoi w grupie „${wynik.wynik?.group.name ?? nazwa}”.`);
  void postawGrupyZrodel(kanal, szyna, idOkna);
}

async function postawGrupyZrodel(
  kanal: Kanal,
  szyna: Element | null,
  idOkna: string,
): Promise<void> {
  if (szyna === null) return;
  const wynik = await wywolaj(kanal, Command.BrowserSourceGroupList, { windowId: idOkna });
  if (!wynik.udany) return;
  const grupy = wynik.wynik?.groups ?? [];
  postawGrupeSzyny(szyna, 'grupy-zrodel', 'Grupy źródeł', grupy.map((grupa) => ({
    tytul: `${grupa.name} — źródeł: ${grupa.sourceIds?.length ?? 0}`,
    cechy: { grupaZrodel: grupa.id },
  })));
}
