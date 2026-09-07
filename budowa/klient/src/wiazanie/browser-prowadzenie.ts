// Prowadzenie karty przeglądania: wejście pod adres, przewinięcie strony oraz
// zmiana i grupowanie kart. Panel wystawiał dotąd sam wykaz kart z zamknięciem.

import { Command, type BrowserSnapshot, type BrowserTab } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import {
  adresWymagany,
  dolozPrzycisk,
  odmowa,
  polePasa,
  powiedz,
  zalozPas,
} from './browser-wspolne.ts';

const PRZEWINIECIE_PIKSELE = 900;

export function zwiazProwadzenieKarty(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  przy: AddEventListenerOptions,
  odswiezKarty: () => void,
): void {
  const panel = korzen.querySelector('#panel-browser');
  const plotno = panel?.querySelector('.brw-plotno') ?? null;
  const poleGrupy = postawPas(panel, plotno);

  korzen.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;

    if (cel.closest('[data-wejdz-pod-adres]') !== null) {
      void wejdzPodAdres(kanal, korzen, idOkna, odswiezKarty);
      return;
    }
    const przewin = cel.closest<HTMLElement>('[data-przewin-strone]');
    if (przewin !== null) {
      void przewinStrone(kanal, idOkna, przewin.dataset.przewinStrone === 'koniec');
      return;
    }
    const przypniecie = cel.closest<HTMLElement>('[data-karta-przypnij]');
    if (przypniecie !== null) {
      void zmienKarte(kanal, przypniecie.dataset.kartaPrzypnij ?? '',
        { pinned: przypniecie.dataset.stanPrzypiecia !== 'tak' }, odswiezKarty);
      return;
    }
    const uspienie = cel.closest<HTMLElement>('[data-karta-uspij]');
    if (uspienie !== null) {
      void zmienKarte(kanal, uspienie.dataset.kartaUspij ?? '', { suspended: true }, odswiezKarty);
      return;
    }
    const grupowanie = cel.closest<HTMLElement>('[data-karta-do-grupy]');
    if (grupowanie !== null) {
      void wstawKarteDoGrupy(kanal, idOkna, grupowanie.dataset.kartaDoGrupy ?? '',
        (poleGrupy?.value ?? '').trim(), odswiezKarty);
    }
  }, przy);
}

/** Dokłada w wierszu karty drogi, których wykaz sam nie niesie. */
export function dolozCzynnosciKarty(wiersz: Element, karta: BrowserTab): void {
  const przypnij = dolozPrzycisk(wiersz, karta.pinned === true ? 'Odepnij' : 'Przypnij',
    'kartaPrzypnij', karta.id);
  if (przypnij !== null) przypnij.dataset.stanPrzypiecia = karta.pinned === true ? 'tak' : 'nie';
  dolozPrzycisk(wiersz, 'Uśpij', 'kartaUspij', karta.id);
  dolozPrzycisk(wiersz, 'Do grupy', 'kartaDoGrupy', karta.id);
}

function postawPas(panel: Element | null, plotno: Element | null): HTMLInputElement | null {
  if (panel === null || plotno === null) return null;
  const pas = zalozPas('Prowadzenie karty');
  dolozPrzycisk(pas, 'Wejdź pod adres', 'wejdzPodAdres');
  dolozPrzycisk(pas, 'Przewiń', 'przewinStrone', 'krok');
  dolozPrzycisk(pas, 'Przewiń na koniec', 'przewinStrone', 'koniec');
  const pole = polePasa(pas, 'Nazwa grupy kart');
  panel.insertBefore(pas, plotno);
  return pole;
}

/* Wejście pod adres idzie kartą czynną, a nie nową: komenda otwarcia karty
   stoi już pod klawiszem Enter w pasku adresu. */
async function wejdzPodAdres(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  odswiezKarty: () => void,
): Promise<void> {
  const url = adresWymagany(korzen);
  if (url === '') return;
  const wynik = await wywolaj(kanal, Command.BrowserNavigate, { windowId: idOkna, url });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił wejścia pod adres.');
    return;
  }
  powiedz(opiszMigawke(wynik.wynik?.snapshot));
  odswiezKarty();
}

async function przewinStrone(kanal: Kanal, idOkna: string, doKonca: boolean): Promise<void> {
  const zadanie = doKonca
    ? { windowId: idOkna, toEnd: true }
    : { windowId: idOkna, deltaY: PRZEWINIECIE_PIKSELE };
  const wynik = await wywolaj(kanal, Command.BrowserScroll, zadanie);
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił przewinięcia strony.');
    return;
  }
  powiedz(opiszMigawke(wynik.wynik?.snapshot));
}

async function zmienKarte(
  kanal: Kanal,
  idKarty: string,
  zmiana: { pinned?: boolean; suspended?: boolean },
  odswiezKarty: () => void,
): Promise<void> {
  if (idKarty === '') return;
  const wynik = await wywolaj(kanal, Command.BrowserTabUpdate, { tabId: idKarty, ...zmiana });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił zmiany karty.');
    return;
  }
  const karta = wynik.wynik?.tab;
  powiedz(`Karta „${karta?.title ?? karta?.url ?? idKarty}” stoi w stanie ${karta?.state ?? '?'}.`);
  odswiezKarty();
}

/* Kontrakt nie zna dokładania karty do grupy: zapis idzie całą grupą, więc
   wykaz kart podaje skład grupy, do którego dochodzi karta wskazana. */
async function wstawKarteDoGrupy(
  kanal: Kanal,
  idOkna: string,
  idKarty: string,
  nazwa: string,
  odswiezKarty: () => void,
): Promise<void> {
  if (idKarty === '') return;
  if (nazwa === '') {
    odmowa(undefined, 'Grupa kart wymaga nazwy — pole nazwy grupy jest puste.');
    return;
  }
  const sklad = await skladGrupy(kanal, idOkna, nazwa, idKarty);
  const wynik = await wywolaj(kanal, Command.BrowserTabGroupSet, {
    windowId: idOkna,
    groupId: sklad.groupId,
    name: nazwa,
    tabIds: sklad.tabIds,
  });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił zapisu grupy kart.');
    return;
  }
  powiedz(`Karta stoi w grupie „${wynik.wynik?.group.name ?? nazwa}”.`);
  odswiezKarty();
}

async function skladGrupy(
  kanal: Kanal,
  idOkna: string,
  nazwa: string,
  idKarty: string,
): Promise<{ groupId?: string; tabIds: string[] }> {
  const wykaz = await wywolaj(kanal, Command.BrowserTabList, { windowId: idOkna });
  const grupa = wykaz.wynik?.groups?.find((pozycja) => pozycja.name === nazwa);
  const karty = new Set(grupa?.tabIds ?? []);
  karty.add(idKarty);
  return { groupId: grupa?.id, tabIds: [...karty] };
}

function opiszMigawke(migawka: BrowserSnapshot | undefined): string {
  if (migawka === undefined) return 'Rdzeń przyjął czynność, ale migawki strony nie oddał.';
  return `Stoi strona „${migawka.title ?? migawka.url}”.`;
}
