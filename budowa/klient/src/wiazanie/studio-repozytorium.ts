// Panel repozytorium sesji Studia: wykaz wersji dokumentu wraz z przywróceniem
// wersji, oznaczeniem wersji kluczowej i wydaniem historii.

import {
  Command,
  EventType,
  StudioAuthor,
  type StudioDocument,
  type StudioVersion,
} from '../../../shared/contract.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import { nazwijOdmowe } from './studio-dokument.ts';

interface WezlyRepozytorium {
  tresc: HTMLElement;
  licznik: HTMLElement | null;
  eksport: HTMLButtonElement | null;
}

interface WzoryWersji {
  wiersz: HTMLElement | null;
  plakietka: HTMLElement | null;
}

let wzoryPanelu: WzoryWersji | null = null;

export function zdejmijTrescPrzykladowaRepozytorium(korzen: ParentNode): boolean {
  return przygotujPanel(korzen) !== null;
}

function przygotujPanel(korzen: ParentNode): { wezly: WezlyRepozytorium; wzory: WzoryWersji } | null {
  const znalezione = zbierzWezly(korzen);
  if (znalezione === null) return null;
  const wiersze = [...znalezione.tresc.querySelectorAll<HTMLElement>('.dn-wersja')];
  wzoryPanelu ??= zdejmijWzory(wiersze);
  for (const wiersz of wiersze) wiersz.remove();
  return { wezly: znalezione, wzory: wzoryPanelu };
}

export function zwiazRepozytorium(
  kanal: Kanal,
  idOkna: string,
  korzen: ParentNode,
): Odsubskrybuj | null {
  const przygotowany = przygotujPanel(korzen);
  if (przygotowany === null) return null;
  const wezly: WezlyRepozytorium = przygotowany.wezly;
  const wzory: WzoryWersji = przygotowany.wzory;

  let dokument: StudioDocument | null = null;

  wezly.tresc.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element) || dokument === null) return;
    const przycisk = cel.closest('.dn-wersja-akcje button');
    const wiersz = przycisk?.closest<HTMLElement>('.dn-wersja');
    const idWersji = wiersz?.dataset.wersja;
    if (!(przycisk instanceof HTMLButtonElement) || idWersji === undefined) return;
    const biezacy = dokument;
    przycisk.disabled = true;
    const czynnosc = przycisk.classList.contains('dn-btn-ikona--sm')
      ? oznaczKluczowa(kanal, idWersji, wiersz?.dataset.kluczowa !== 'true')
      : przywrocWersje(kanal, biezacy.id, idWersji);
    void czynnosc.then((przywrocony) => {
      przycisk.disabled = false;
      dokument = przywrocony ?? biezacy;
      void wypelnijWykaz(kanal, wezly, wzory, dokument);
    });
  });

  const eksport = wezly.eksport;
  if (eksport !== null) {
    // Archiwum zostaje w magazynie: panel nie ma węzła na jego odnośnik.
    eksport.addEventListener('click', () => {
      const biezacy = dokument;
      if (biezacy === null) return;
      eksport.disabled = true;
      void wywolaj(kanal, Command.StudioRepositoryExport, {
        documentId: biezacy.id,
        windowId: idOkna,
      }).then(() => {
        eksport.disabled = false;
      });
    });
  }

  // Panel wiąże się przed założeniem dokumentu, więc zdarzenie rozstrzyga oknem.
  const odlacz = kanal.naZdarzenie(EventType.StudioDocumentChanged, (zmiana) => {
    if (zmiana.document.windowId !== idOkna) return;
    dokument = zmiana.document;
    void wypelnijWykaz(kanal, wezly, wzory, zmiana.document);
  });

  void otworzDokument(kanal, idOkna).then((otwarty) => {
    dokument = otwarty;
    if (otwarty === null) {
      wezly.licznik?.remove();
      return;
    }
    void wypelnijWykaz(kanal, wezly, wzory, otwarty);
  });

  return odlacz;
}

function zbierzWezly(korzen: ParentNode): WezlyRepozytorium | null {
  const panel = korzen.querySelector('#panel-repo');
  const tresc = panel?.querySelector('.sta-okno-tresc');
  if (!(tresc instanceof HTMLElement)) return null;
  const licznik = panel?.querySelector('.sta-okno-znacznik');
  const eksport = tresc.querySelector('.dn-wersja-stopka .dn-btn');
  return {
    tresc,
    licznik: licznik instanceof HTMLElement ? licznik : null,
    eksport: eksport instanceof HTMLButtonElement ? eksport : null,
  };
}

function zdejmijWzory(wiersze: HTMLElement[]): WzoryWersji {
  const wiersz = sklonuj(wiersze[0] ?? null);
  if (wiersz !== null) ustawCzynnoscWiersza(wiersz, wiersze[1] ?? null);
  return { wiersz, plakietka: zdejmijWzorPlakietki(wiersze) };
}

// Podgląd wersji wymaga formatu docelowego, a porównanie wskazania drugiej
// wersji; panel nie ma węzła, w którym Operator poda którąkolwiek z tych wartości.
function ustawCzynnoscWiersza(wzor: HTMLElement, drugi: HTMLElement | null): void {
  wzor.querySelector('.dn-wersja-akcje .dn-btn--duch')?.remove();
  wzor.querySelector('.dn-wersja-akcje .dn-btn-ikona--sm')?.setAttribute(
    'aria-label', 'Wersja kluczowa',
  );
  const przycisk = wzor.querySelector('.dn-wersja-akcje .dn-btn--zarys');
  if (przycisk === null) return;
  const napis = drugi?.querySelector('.dn-wersja-akcje .dn-btn--zarys')?.textContent?.trim() ?? '';
  if (napis === '') przycisk.remove();
  else przycisk.textContent = napis;
}

// Nazwa własna wersji nie ma w panelu pola, więc czynność przestawia samo
// oznaczenie wersji kluczowej; dokument zostaje ten sam.
async function oznaczKluczowa(
  kanal: Kanal,
  idWersji: string,
  kluczowa: boolean,
): Promise<StudioDocument | null> {
  const wynik = await wywolaj(kanal, Command.StudioVersionLabelSet, {
    versionId: idWersji,
    milestone: kluczowa,
  });
  if (!wynik.udany) {
    oglos('Studio', nazwijOdmowe('Oznaczenie wersji kluczowej', wynik.blad), 'ostrzezenie');
  }
  return null;
}

function zdejmijWzorPlakietki(wiersze: HTMLElement[]): HTMLElement | null {
  for (const wiersz of wiersze) {
    const przykladowa = wiersz.querySelector('.dn-plakietka');
    const opakowanie = sklonuj(przykladowa?.parentElement ?? null);
    const plakietka = opakowanie?.querySelector('.dn-plakietka') ?? null;
    if (opakowanie === null || plakietka === null) continue;
    const znak = (przykladowa?.textContent ?? '').trim().split(' ')[0] ?? '';
    if (znak === '') continue;
    // Plakietka niesie sam znak wersji kluczowej, wzięty ze znacznika.
    plakietka.textContent = znak;
    return opakowanie;
  }
  return null;
}

function sklonuj(wezel: Element | null): HTMLElement | null {
  return wezel instanceof HTMLElement ? (wezel.cloneNode(true) as HTMLElement) : null;
}

async function wypelnijWykaz(
  kanal: Kanal,
  wezly: WezlyRepozytorium,
  wzory: WzoryWersji,
  dokument: StudioDocument,
): Promise<void> {
  const szereg = await wywolaj(kanal, Command.StudioVersionSeriesList, {
    documentId: dokument.id,
  });
  const wRozbiciu = szereg.udany ? szereg.wynik?.versions ?? [] : [];
  const liczba = szereg.udany && szereg.wynik !== undefined
    ? szereg.wynik.operatorCount + szereg.wynik.autosaveCount
    : 0;
  if (wRozbiciu.length > 0) {
    opiszLicznik(wezly.licznik, liczba);
    postawWiersze(wezly, wzory, wRozbiciu, dokument);
    return;
  }
  // Rozbicie na szeregi bywa puste, a wtedy wykaz historii jest jedynym źródłem wierszy.
  const wykaz = await wywolaj(kanal, Command.StudioRepositoryList, { documentId: dokument.id });
  if (!wykaz.udany || wykaz.wynik === undefined) {
    wezly.licznik?.remove();
    return;
  }
  opiszLicznik(wezly.licznik, liczba > 0 ? liczba : wykaz.wynik.versions.length);
  postawWiersze(wezly, wzory, wykaz.wynik.versions, dokument);
}

function postawWiersze(
  wezly: WezlyRepozytorium,
  wzory: WzoryWersji,
  wersje: StudioVersion[],
  dokument: StudioDocument,
): void {
  if (wzory.wiersz === null) return;
  for (const stojacy of wezly.tresc.querySelectorAll('.dn-wersja')) stojacy.remove();
  const stopka = wezly.tresc.querySelector('.dn-wersja-stopka');
  for (const wersja of wersje) {
    const biezaca = wersja.id === dokument.versionId;
    const wiersz = zbudujWiersz(wzory.wiersz, wzory.plakietka, wersja, biezaca);
    if (stopka === null) wezly.tresc.appendChild(wiersz);
    else stopka.before(wiersz);
  }
}

function zbudujWiersz(
  wzor: HTMLElement,
  wzorPlakietki: HTMLElement | null,
  wersja: StudioVersion,
  biezaca: boolean,
): HTMLElement {
  const wiersz = wzor.cloneNode(true) as HTMLElement;
  wiersz.dataset.wersja = wersja.id;
  wiersz.dataset.kluczowa = wersja.milestone === true ? 'true' : 'false';
  const kropka = wiersz.querySelector('.dn-kropka');
  if (kropka !== null) {
    kropka.classList.toggle('dn-kropka--sukces', biezaca);
    kropka.classList.toggle('dn-kropka--neutralna', !biezaca);
  }
  wpiszTytul(wiersz, wersja.label ?? '');
  const meta = wiersz.querySelector('.dn-meta');
  if (meta !== null) meta.textContent = opiszZalozenie(wersja);
  const nota = wiersz.querySelector('.dn-nota');
  if (nota !== null) {
    if (wersja.summary === undefined || wersja.summary === '') nota.remove();
    else nota.textContent = wersja.summary;
  }
  if (wersja.milestone === true) wstawZnakKluczowej(wiersz, wzorPlakietki);
  return wiersz;
}

function wpiszTytul(wiersz: HTMLElement, nazwa: string): void {
  const tytul = wiersz.querySelector('.st-panel-wiersz b');
  if (tytul === null) return;
  if (nazwa === '') tytul.remove();
  else tytul.textContent = nazwa;
}

function wstawZnakKluczowej(wiersz: HTMLElement, wzor: HTMLElement | null): void {
  if (wzor === null) return;
  const oznaczenie = wzor.cloneNode(true) as HTMLElement;
  const akcje = wiersz.querySelector('.dn-wersja-akcje');
  if (akcje === null) wiersz.appendChild(oznaczenie);
  else akcje.before(oznaczenie);
}

function opiszZalozenie(wersja: StudioVersion): string {
  const godzina = new Date(wersja.createdAt).toLocaleTimeString('pl-PL', {
    hour: '2-digit',
    minute: '2-digit',
  });
  const autor = wersja.authorAgentName ?? nazwaAutora(wersja.author);
  return autor === '' ? godzina : `${godzina} · ${autor}`;
}

function nazwaAutora(autor: StudioAuthor | undefined): string {
  if (autor === StudioAuthor.Uzytkownik) return 'użytkownik';
  if (autor === StudioAuthor.Model) return 'model';
  return '';
}

function opiszLicznik(licznik: HTMLElement | null, liczba: number): void {
  if (licznik === null) return;
  licznik.textContent = `${liczba} ${odmianaWersji(liczba)}`;
}

function odmianaWersji(liczba: number): string {
  if (liczba === 1) return 'wersja';
  const jednosci = liczba % 10;
  const dziesiatki = liczba % 100;
  const mnoga = jednosci >= 2 && jednosci <= 4 && (dziesiatki < 12 || dziesiatki > 14);
  return mnoga ? 'wersje' : 'wersji';
}

async function otworzDokument(kanal: Kanal, idOkna: string): Promise<StudioDocument | null> {
  const wynik = await wywolaj(kanal, Command.StudioDocumentOpen, { windowId: idOkna });
  if (!wynik.udany || wynik.wynik === undefined) return null;
  return wynik.wynik.document;
}

async function przywrocWersje(
  kanal: Kanal,
  idDokumentu: string,
  idWersji: string,
): Promise<StudioDocument | null> {
  const wynik = await wywolaj(kanal, Command.StudioRepositoryRestore, {
    documentId: idDokumentu,
    versionId: idWersji,
  });
  if (!wynik.udany || wynik.wynik === undefined) return null;
  return wynik.wynik.document;
}
