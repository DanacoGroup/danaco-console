// Obszar roboczy i granice wykonawcy: zapis i zdjęcie przestrzeni kart,
// odczyt i nastawa granic, sterowanie pobieraniem, nagranie makra oraz
// dołożenie wytworu z treści przeglądanej strony.

import {
  BrowserArtifactKind,
  BrowserMacroAction,
  Command,
  type BrowserDownloadAction,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import {
  dolozPrzycisk,
  naBase64,
  odmowa,
  polePasa,
  potwierdzone,
  powiedz,
  zalozPas,
} from './browser-wspolne.ts';

interface PolaObszaru {
  nazwaObszaru: HTMLInputElement | null;
  krokiGranicy: HTMLInputElement | null;
  czasGranicy: HTMLInputElement | null;
}

export interface OdswiezeniaObszaru {
  odswiezPrzestrzenie: () => void;
  odswiezPobrania: () => void;
}

export function zwiazObszarIGranice(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  przy: AddEventListenerOptions,
  odswiezenia: OdswiezeniaObszaru,
): void {
  const pola = postawPas(korzen);
  let makro = '';

  korzen.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;

    if (cel.closest('[data-obszar-zapisz]') !== null) {
      void zapiszObszar(kanal, idOkna, pola.nazwaObszaru, odswiezenia.odswiezPrzestrzenie);
      return;
    }
    const usuniecie = cel.closest<HTMLElement>('[data-przestrzen-usun]');
    if (usuniecie !== null) {
      void usunObszar(kanal, usuniecie, odswiezenia.odswiezPrzestrzenie);
      return;
    }
    if (cel.closest('[data-granice-odczytaj]') !== null) {
      void odczytajGranice(kanal, idOkna, pola);
      return;
    }
    if (cel.closest('[data-granice-zapisz]') !== null) {
      void zapiszGranice(kanal, idOkna, pola);
      return;
    }
    const pobranie = cel.closest<HTMLElement>('[data-pobranie-czynnosc]');
    if (pobranie !== null) {
      void sterujPobraniem(kanal, pobranie, odswiezenia.odswiezPobrania);
      return;
    }
    const makropis = cel.closest<HTMLElement>('[data-makro]');
    if (makropis !== null) {
      void nagrajMakro(kanal, idOkna, makropis.dataset.makro === 'stop', makro)
        .then((kod) => { makro = kod; });
      return;
    }
    if (cel.closest('[data-wytwor-doloz]') !== null) {
      void dolozWytwor(kanal, idOkna);
    }
  }, przy);
}

function postawPas(korzen: Element): PolaObszaru {
  const szyna = korzen.querySelector('.dn-szyna-modulu-lista');
  if (szyna === null) return { nazwaObszaru: null, krokiGranicy: null, czasGranicy: null };
  const pas = zalozPas('Obszar roboczy i granice');
  const nazwaObszaru = polePasa(pas, 'Nazwa obszaru kart');
  dolozPrzycisk(pas, 'Zapisz obszar', 'obszarZapisz');
  const krokiGranicy = polePasa(pas, 'Górna liczba kroków');
  const czasGranicy = polePasa(pas, 'Górny czas w sekundach');
  dolozPrzycisk(pas, 'Odczytaj granice', 'graniceOdczytaj');
  dolozPrzycisk(pas, 'Zapisz granice', 'graniceZapisz');
  dolozPrzycisk(pas, 'Nagrywaj makro', 'makro', 'start');
  dolozPrzycisk(pas, 'Zakończ makro', 'makro', 'stop');
  dolozPrzycisk(pas, 'Dołóż wytwór', 'wytworDoloz');
  szyna.append(pas);
  return { nazwaObszaru, krokiGranicy, czasGranicy };
}

async function zapiszObszar(
  kanal: Kanal,
  idOkna: string,
  pole: HTMLInputElement | null,
  odswiez: () => void,
): Promise<void> {
  const name = (pole?.value ?? '').trim();
  if (name === '') {
    odmowa(undefined, 'Obszar wymaga nazwy — pole nazwy obszaru jest puste.');
    return;
  }
  const wynik = await wywolaj(kanal, Command.BrowserWorkspaceSave, { windowId: idOkna, name });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił zapisu obszaru kart.');
    return;
  }
  powiedz(`Obszar „${name}” obejmuje ${wynik.wynik?.workspace.tabCount ?? 0} kart.`);
  odswiez();
}

/* Usunięcie obszaru jest nieodwracalne: zestaw kart schodzi bez kopii, więc
   czynność wymaga drugiego naciśnięcia. */
async function usunObszar(
  kanal: Kanal,
  przycisk: HTMLElement,
  odswiez: () => void,
): Promise<void> {
  const workspaceId = przycisk.dataset.przestrzenUsun ?? '';
  if (workspaceId === '') return;
  if (!potwierdzone(przycisk, 'Usunięcie obszaru kasuje zestaw kart. Naciśnij drugi raz.')) return;
  const wynik = await wywolaj(kanal, Command.BrowserWorkspaceRemove, { workspaceId });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił usunięcia obszaru.');
    return;
  }
  powiedz(wynik.wynik?.removed === true
    ? 'Obszar usunięty; karty otwarte zostały nietknięte.'
    : 'Rdzeń przyjął żądanie, ale obszaru nie usunął.');
  odswiez();
}

async function odczytajGranice(kanal: Kanal, idOkna: string, pola: PolaObszaru): Promise<void> {
  const wynik = await wywolaj(kanal, Command.BrowserExecutorLimitsGet, { windowId: idOkna });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił odczytu granic wykonawcy.');
    return;
  }
  const granice = wynik.wynik?.limits;
  if (granice === undefined) return;
  if (pola.krokiGranicy !== null) pola.krokiGranicy.value = String(granice.maxSteps);
  if (pola.czasGranicy !== null) pola.czasGranicy.value = String(granice.maxDurationSeconds);
  powiedz(`Granice warstwy ${granice.scope}: ${granice.maxSteps} kroków, `
    + `${granice.maxDurationSeconds} s.`);
}

async function zapiszGranice(kanal: Kanal, idOkna: string, pola: PolaObszaru): Promise<void> {
  const kroki = liczbaZPola(pola.krokiGranicy);
  const czas = liczbaZPola(pola.czasGranicy);
  if (kroki === undefined && czas === undefined) {
    odmowa(undefined, 'Ani liczba kroków, ani czas nie są liczbą — nie ma czego zapisać.');
    return;
  }
  const wynik = await wywolaj(kanal, Command.BrowserExecutorLimitsSet, {
    windowId: idOkna,
    maxSteps: kroki,
    maxDurationSeconds: czas,
  });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił nastawy granic wykonawcy.');
    return;
  }
  const granice = wynik.wynik?.limits;
  powiedz(`Granice stoją: ${granice?.maxSteps ?? '?'} kroków, `
    + `${granice?.maxDurationSeconds ?? '?'} s.`);
}

function liczbaZPola(pole: HTMLInputElement | null): number | undefined {
  const wpis = (pole?.value ?? '').trim();
  if (wpis === '') return undefined;
  const liczba = Number.parseInt(wpis, 10);
  return Number.isFinite(liczba) && liczba > 0 ? liczba : undefined;
}

async function sterujPobraniem(
  kanal: Kanal,
  przycisk: HTMLElement,
  odswiez: () => void,
): Promise<void> {
  const downloadId = przycisk.closest<HTMLElement>('[data-pobranie]')?.dataset.pobranie ?? '';
  const czynnosc = przycisk.dataset.pobranieCzynnosc ?? '';
  if (downloadId === '' || czynnosc === '') return;
  const wynik = await wywolaj(kanal, Command.BrowserDownloadControl, {
    downloadId,
    action: czynnosc as BrowserDownloadAction,
  });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił sterowania pobieraniem.');
    return;
  }
  powiedz(`Pobranie stoi w stanie ${wynik.wynik?.download.status ?? '?'}.`);
  odswiez();
}

async function nagrajMakro(
  kanal: Kanal,
  idOkna: string,
  koniec: boolean,
  stojace: string,
): Promise<string> {
  const wynik = await wywolaj(kanal, Command.BrowserMacroRecord, {
    windowId: idOkna,
    action: koniec ? BrowserMacroAction.Stop : BrowserMacroAction.Start,
    macroId: koniec && stojace !== '' ? stojace : undefined,
  });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił nagrania makra przeglądania.');
    return stojace;
  }
  const makro = wynik.wynik?.macro;
  powiedz(wynik.wynik?.recording === true
    ? `Makro „${makro?.name ?? ''}” się nagrywa.`
    : `Makro „${makro?.name ?? ''}” zamknięte; kroków: ${makro?.steps?.length ?? 0}.`);
  return makro?.id ?? '';
}

/* Wytwór powstaje z treści migawki: rdzeń przyjmuje treść wprost, a okno nie
   ma innego materiału niż to, co sam pokazał na płótnie. */
async function dolozWytwor(kanal: Kanal, idOkna: string): Promise<void> {
  const migawka = await wywolaj(kanal, Command.BrowserSnapshotGet, {
    windowId: idOkna,
    includeHtml: false,
  });
  if (!migawka.udany) {
    odmowa(migawka.blad, 'Rdzeń odmówił odbitki strony, więc wytwór nie ma treści.');
    return;
  }
  const tresc = migawka.wynik?.snapshot.text ?? '';
  if (tresc === '') {
    odmowa(undefined, 'Odbitka strony przyszła bez treści — wytwór nie powstanie.');
    return;
  }
  const wynik = await wywolaj(kanal, Command.BrowserArtifactAdd, {
    windowId: idOkna,
    kind: BrowserArtifactKind.Extraction,
    title: migawka.wynik?.snapshot.title ?? migawka.wynik?.snapshot.url,
    contentBase64: naBase64(tresc),
    mimeType: 'text/plain',
    sourceUrl: migawka.wynik?.snapshot.url,
    snapshotId: migawka.wynik?.snapshot.id,
  });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił dołożenia wytworu.');
    return;
  }
  powiedz(`Wytwór „${wynik.wynik?.artifact.title ?? ''}” stoi w magazynie rdzenia.`);
}
