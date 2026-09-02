// Panel „Versioning Panel" wobec wersji rdzenia. Prototyp nie ma pola wyboru
// pliku wersji, więc nową wersję niesie upuszczenie pliku z dysku na panel.
import { Command, type LibraryDiff, type LibraryVersion } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import {
  panel,
  pozycja,
  pustka,
  skrocony,
  wartosc,
  zdejmijDzieci,
  type Kontekst,
} from './library-wspolne.ts';

const GRANICA_WERSJI = 100;

export function zwiazWersje(
  kanal: Kanal,
  korzen: ParentNode,
  kontekst: Kontekst,
  przy: AddEventListenerOptions,
): void {
  const trescMozeMoze = panel(korzen, 'panel-versioning');
  if (trescMozeMoze === null) return;
  const trescMoze = trescMozeMoze;
  const tresc = trescMoze;
  const przyciski = pasekPrzyciskow(tresc);
  const wzorPorownania = tresc.querySelector('.dn-btn')?.textContent?.trim() ?? '';
  const wersje = new Map<HTMLElement, LibraryVersion>();
  let wybrana: LibraryVersion | null = null;
  let wykaz: LibraryVersion[] = [];

  async function odswiez(): Promise<void> {
    const plik = kontekst.plik();
    wersje.clear();
    wybrana = null;
    wykaz = [];
    if (plik === null) {
      pustka(tresc, 'Wybierz plik w wykazie, aby rdzeń podał jego wersje.');
      if (przyciski !== null) tresc.appendChild(przyciski);
      return;
    }
    const odpowiedz = await wywolaj(kanal, Command.LibraryVersionList, {
      fileId: plik.id,
      limit: GRANICA_WERSJI,
    });
    wykaz = wartosc(odpowiedz)?.versions ?? [];
    if (wykaz.length === 0) {
      pustka(tresc, 'Rdzeń nie zwrócił wersji tego pliku.');
      if (przyciski !== null) tresc.appendChild(przyciski);
      return;
    }
    zdejmijDzieci(tresc);
    for (const wersja of wykaz) {
      const wiersz = zbudujWersje(tresc.ownerDocument, wersja, wersja.id === plik.versionId);
      wersje.set(wiersz, wersja);
      tresc.appendChild(wiersz);
    }
    if (przyciski !== null) {
      opiszPorownanie(przyciski, wzorPorownania, wykaz);
      tresc.appendChild(przyciski);
    }
  }

  tresc.addEventListener(
    'click',
    (zdarzenie) => {
      const cel = zdarzenie.target;
      if (!(cel instanceof Element)) return;
      const przycisk = cel.closest('.dn-btn');
      if (przycisk !== null) {
        const podpis = przycisk.textContent?.trim() ?? '';
        if (podpis.startsWith('Przywróć')) {
          void przywrocWersje(kanal, kontekst, wybrana).then(odswiez);
          return;
        }
        void porownaj(kanal, kontekst, wykaz);
        return;
      }
      const wiersz = cel.closest('.dn-wykaz-modulu-poz');
      const wersja = wiersz instanceof HTMLElement ? wersje.get(wiersz) : undefined;
      if (wersja === undefined) return;
      for (const inny of wersje.keys()) inny.removeAttribute('aria-selected');
      wiersz?.setAttribute('aria-selected', 'true');
      wybrana = wersja;
    },
    przy,
  );

  tresc.addEventListener(
    'dragover',
    (zdarzenie) => {
      if (zdarzenie.dataTransfer?.types.includes('Files') === true) zdarzenie.preventDefault();
    },
    przy,
  );

  tresc.addEventListener(
    'drop',
    (zdarzenie) => {
      const plikZDysku = zdarzenie.dataTransfer?.files.item(0);
      if (plikZDysku === null || plikZDysku === undefined) return;
      zdarzenie.preventDefault();
      const czytnik = new FileReader();
      czytnik.addEventListener('load', () => {
        const dane = typeof czytnik.result === 'string' ? czytnik.result : '';
        void dolozWersje(kanal, kontekst, plikZDysku.name, dane.slice(dane.indexOf(',') + 1)).then(
          odswiez,
        );
      });
      czytnik.readAsDataURL(plikZDysku);
    },
    przy,
  );

  kontekst.przyPliku(() => {
    void odswiez();
  });
  kontekst.zglosOdswiezenie(() => {
    void odswiez();
  });
  void odswiez();
}

function pasekPrzyciskow(tresc: HTMLElement): HTMLElement | null {
  const przycisk = tresc.querySelector('.dn-btn');
  return przycisk?.parentElement ?? null;
}

function zbudujWersje(
  dokument: Document,
  wersja: LibraryVersion,
  biezaca: boolean,
): HTMLElement {
  const podpis = wersja.label ?? wersja.id;
  const opis = [wersja.author ?? '', skrocony(wersja.createdAt)]
    .filter((czesc) => czesc !== '')
    .join(' · ');
  const wiersz = pozycja(dokument, biezaca ? `${podpis} — bieżąca` : podpis, opis);
  const znak = dokument.createElement('span');
  znak.className = biezaca ? 'dn-kropka dn-kropka--sukces' : 'dn-kropka dn-kropka--neutralna';
  znak.setAttribute('aria-hidden', 'true');
  wiersz.prepend(znak);
  wiersz.tabIndex = 0;
  return wiersz;
}

function opiszPorownanie(
  przyciski: Element,
  wzor: string,
  wykaz: LibraryVersion[],
): void {
  const przycisk = przyciski.querySelector('.dn-btn');
  if (przycisk === null || wykaz.length < 2) return;
  const rdzen = wzor.split(' ')[0];
  const lewa = wykaz[1].label ?? wykaz[1].id;
  const prawa = wykaz[0].label ?? wykaz[0].id;
  przycisk.textContent = `${rdzen} ${lewa} ↔ ${prawa}`;
}

async function przywrocWersje(
  kanal: Kanal,
  kontekst: Kontekst,
  wersja: LibraryVersion | null,
): Promise<void> {
  const plik = kontekst.plik();
  if (plik === null || wersja === null) return;
  const odpowiedz = await wywolaj(kanal, Command.LibraryVersionRestore, {
    fileId: plik.id,
    versionId: wersja.id,
  });
  const wynikMozeMoze = wartosc(odpowiedz);
  if (wynikMozeMoze === null) return;
  const wynikMoze = wynikMozeMoze;
  const wynik = wynikMoze;
  kontekst.ustawPlik(wynik.file);
}

async function porownaj(
  kanal: Kanal,
  kontekst: Kontekst,
  wykaz: LibraryVersion[],
): Promise<void> {
  const plik = kontekst.plik();
  if (plik === null || wykaz.length < 2) return;
  const odpowiedz = await wywolaj(kanal, Command.LibraryDiffCompare, {
    leftFileId: plik.id,
    leftVersionId: wykaz[1].id,
    rightVersionId: wykaz[0].id,
    contextLines: 3,
  });
  const roznicaMozeMoze = wartosc(odpowiedz)?.diff ?? null;
  if (roznicaMozeMoze === null) return;
  const roznicaMoze = roznicaMozeMoze;
  const roznica = roznicaMoze;
  oglos('Library', opiszRoznice(roznica));
}

function opiszRoznice(roznica: LibraryDiff): string {
  if (roznica.identical) return `${roznica.leftLabel} i ${roznica.rightLabel} są identyczne.`;
  return `${roznica.leftLabel} ↔ ${roznica.rightLabel}: ${String(roznica.hunks.length)} zmian.`;
}

async function dolozWersje(
  kanal: Kanal,
  kontekst: Kontekst,
  nazwa: string,
  tresc64: string,
): Promise<void> {
  const plik = kontekst.plik();
  if (plik === null || tresc64 === '') return;
  const odpowiedz = await wywolaj(kanal, Command.LibraryVersionAdd, {
    fileId: plik.id,
    contentBase64: tresc64,
    label: nazwa,
  });
  const wynikMozeMoze = wartosc(odpowiedz);
  if (wynikMozeMoze === null) return;
  const wynikMoze = wynikMozeMoze;
  const wynik = wynikMoze;
  kontekst.ustawPlik(wynik.file);
}
