/**
 * Wiązanie panelu repozytorium sesji Studia z rdzeniem. Znacznik niesie
 * biblioteka Właściciela; ten plik wypełnia wykaz wersji odpowiedzią rdzenia,
 * powielając wzór wiersza zdjęty z treści przykładowej, i wiąże czynności
 * wersji z komendami kontraktu.
 */

import {
  Command,
  EventType,
  StudioAuthor,
  type StudioDocument,
  type StudioVersion,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';

/** Węzły panelu repozytorium; brak treści panelu znaczy, że karta nie stoi w dokumencie. */
interface WezlyRepozytorium {
  tresc: HTMLElement;
  licznik: HTMLElement | null;
  eksport: HTMLButtonElement | null;
}

/** Wzory zdjęte z wierszy przykładowych: kształt wiersza wersji i oznaczenie wersji kluczowej. */
interface WzoryWersji {
  wiersz: HTMLElement | null;
  plakietka: HTMLElement | null;
}

/**
 * Wiąże panel repozytorium sesji okna Studia. Zwraca prawdę, gdy znacznik
 * panelu stał i wiązanie zostało założone.
 */
export function zwiazRepozytorium(kanal: Kanal, idOkna: string): boolean {
  const znalezione = zbierzWezly();
  if (znalezione === null) return false;
  const wezly: WezlyRepozytorium = znalezione;

  const wiersze = [...wezly.tresc.querySelectorAll<HTMLElement>('.dn-wersja')];
  const wzory = zdejmijWzory(wiersze);
  // Wiersze przykładowe znikają, zanim padnie pierwsza odpowiedź rdzenia:
  // pusty wykaz jest uczciwy, wykaz z cudzą historią — nie.
  for (const wiersz of wiersze) wiersz.remove();

  let dokument: StudioDocument | null = null;

  wezly.tresc.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element) || dokument === null) return;
    const przycisk = cel.closest('.dn-wersja-akcje .dn-btn--zarys');
    const idWersji = przycisk?.closest<HTMLElement>('.dn-wersja')?.dataset.wersja;
    if (!(przycisk instanceof HTMLButtonElement) || idWersji === undefined) return;
    przycisk.disabled = true;
    void przywrocWersje(kanal, dokument.id, idWersji).then((przywrocony) => {
      przycisk.disabled = false;
      if (przywrocony === null) return;
      dokument = przywrocony;
      void wypelnijWykaz(kanal, wezly, wzory, przywrocony);
    });
  });

  const eksport = wezly.eksport;
  if (eksport !== null) {
    // Wydanie bez wskazania wersji bierze całą historię. Archiwum zostaje
    // w magazynie: panel nie ma węzła na jego odnośnik, więc o trwaniu
    // czynności mówi wyłącznie zablokowany przycisk.
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

  kanal.naZdarzenie(EventType.StudioDocumentChanged, (zmiana) => {
    if (dokument === null || zmiana.document.id !== dokument.id) return;
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

  return true;
}

/** Wskazuje węzły panelu; pustka znaczy, że karta repozytorium nie stoi w dokumencie. */
function zbierzWezly(): WezlyRepozytorium | null {
  const panel = document.querySelector('#panel-repo');
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

/** Zdejmuje z wierszy przykładowych wzór wiersza wraz z jego czynnością oraz wzór oznaczenia wersji kluczowej. */
function zdejmijWzory(wiersze: HTMLElement[]): WzoryWersji {
  const wiersz = sklonuj(wiersze[0] ?? null);
  if (wiersz !== null) ustawCzynnoscWiersza(wiersz, wiersze[1] ?? null);
  return { wiersz, plakietka: zdejmijWzorPlakietki(wiersze) };
}

/*
Wiersz zostaje przy jednej czynności — przywróceniu, które bierze samą wersję
wiersza. Podgląd wymaga formatu docelowego, menu wersji nazwy własnej,
a porównanie wskazania drugiej wersji; panel nie ma węzła, w którym Operator
poda którąkolwiek z tych wartości.
*/
function ustawCzynnoscWiersza(wzor: HTMLElement, drugi: HTMLElement | null): void {
  wzor.querySelector('.dn-wersja-akcje .dn-btn--duch')?.remove();
  wzor.querySelector('.dn-wersja-akcje .dn-btn-ikona--sm')?.remove();
  const przycisk = wzor.querySelector('.dn-wersja-akcje .dn-btn--zarys');
  if (przycisk === null) return;
  const napis = drugi?.querySelector('.dn-wersja-akcje .dn-btn--zarys')?.textContent?.trim() ?? '';
  if (napis === '') przycisk.remove();
  else przycisk.textContent = napis;
}

/** Zdejmuje opakowanie plakietki z wiersza przykładowego i zostawia w nim sam znak wersji kluczowej. */
function zdejmijWzorPlakietki(wiersze: HTMLElement[]): HTMLElement | null {
  for (const wiersz of wiersze) {
    const przykladowa = wiersz.querySelector('.dn-plakietka');
    const opakowanie = sklonuj(przykladowa?.parentElement ?? null);
    const plakietka = opakowanie?.querySelector('.dn-plakietka') ?? null;
    if (opakowanie === null || plakietka === null) continue;
    const znak = (przykladowa?.textContent ?? '').trim().split(' ')[0] ?? '';
    if (znak === '') continue;
    /* Nazwa własna wersji stoi już w tytule wiersza, więc plakietka niesie sam
       znak wersji kluczowej; znak pochodzi ze znacznika, nie z tego pliku. */
    plakietka.textContent = znak;
    return opakowanie;
  }
  return null;
}

/** Klon węzła wzorcowego, odporny na jego brak w znaczniku. */
function sklonuj(wezel: Element | null): HTMLElement | null {
  return wezel instanceof HTMLElement ? (wezel.cloneNode(true) as HTMLElement) : null;
}

/** Wczytuje historię wersji dokumentu i stawia z niej wykaz wraz z licznikiem w belce. */
async function wypelnijWykaz(
  kanal: Kanal,
  wezly: WezlyRepozytorium,
  wzory: WzoryWersji,
  dokument: StudioDocument,
): Promise<void> {
  const szereg = await wywolaj(kanal, Command.StudioVersionSeriesList, {
    documentId: dokument.id,
  });
  if (szereg.udany && szereg.wynik !== undefined) {
    opiszLicznik(wezly.licznik, szereg.wynik.operatorCount + szereg.wynik.autosaveCount);
    postawWiersze(wezly, wzory, szereg.wynik.versions, dokument);
    return;
  }
  /* Licznik w belce ma pokrycie wyłącznie w rozbiciu na szeregi: historia
     z `studio.repository.list` jest ucięta limitem, więc długość jej tablicy
     nie jest liczbą wersji dokumentu. */
  wezly.licznik?.remove();
  const wykaz = await wywolaj(kanal, Command.StudioRepositoryList, { documentId: dokument.id });
  if (!wykaz.udany || wykaz.wynik === undefined) return;
  postawWiersze(wezly, wzory, wykaz.wynik.versions, dokument);
}

/** Stawia po jednym wierszu na wersję, przed stopką panelu; wiersze poprzedniego wczytania schodzą. */
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

/** Zwraca klon wzoru wiersza opisany wersją rdzenia; wiersz niesie identyfikator wersji, bo po nim rozstrzyga się przywrócenie. */
function zbudujWiersz(
  wzor: HTMLElement,
  wzorPlakietki: HTMLElement | null,
  wersja: StudioVersion,
  biezaca: boolean,
): HTMLElement {
  const wiersz = wzor.cloneNode(true) as HTMLElement;
  wiersz.dataset.wersja = wersja.id;
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
  if (wersja.milestone === true) oznaczKluczowa(wiersz, wzorPlakietki);
  return wiersz;
}

/** Wpisuje w tytuł wiersza nazwę własną wersji; wersja bez nadanej nazwy zostaje bez tytułu, bo liczby porządkowej kontrakt nie niesie. */
function wpiszTytul(wiersz: HTMLElement, nazwa: string): void {
  const tytul = wiersz.querySelector('.st-panel-wiersz b');
  if (tytul === null) return;
  if (nazwa === '') tytul.remove();
  else tytul.textContent = nazwa;
}

/** Wstawia w wiersz oznaczenie wersji kluczowej, przed pasem czynności. */
function oznaczKluczowa(wiersz: HTMLElement, wzor: HTMLElement | null): void {
  if (wzor === null) return;
  const oznaczenie = wzor.cloneNode(true) as HTMLElement;
  const akcje = wiersz.querySelector('.dn-wersja-akcje');
  if (akcje === null) wiersz.appendChild(oznaczenie);
  else akcje.before(oznaczenie);
}

/** Godzina założenia wersji wraz z autorem; autor nieobowiązkowy w kontrakcie schodzi z opisu, a nie dostaje wartości domyślnej. */
function opiszZalozenie(wersja: StudioVersion): string {
  const godzina = new Date(wersja.createdAt).toLocaleTimeString('pl-PL', {
    hour: '2-digit',
    minute: '2-digit',
  });
  const autor = wersja.authorAgentName ?? nazwaAutora(wersja.author);
  return autor === '' ? godzina : `${godzina} · ${autor}`;
}

/** Nazwa autora wersji w brzmieniu znacznika Właściciela; słownictwo pochodzi z prototypu, nie z tego pliku. */
function nazwaAutora(autor: StudioAuthor | undefined): string {
  if (autor === StudioAuthor.Uzytkownik) return 'użytkownik';
  if (autor === StudioAuthor.Model) return 'model';
  return '';
}

/** Wpisuje w znacznik belki liczbę wersji dokumentu. */
function opiszLicznik(licznik: HTMLElement | null, liczba: number): void {
  if (licznik === null) return;
  licznik.textContent = `${liczba} ${odmianaWersji(liczba)}`;
}

/** Odmiana rzeczownika z licznika belki; liczebnik polski wymaga trzech form tego samego słowa. */
function odmianaWersji(liczba: number): string {
  if (liczba === 1) return 'wersja';
  const jednosci = liczba % 10;
  const dziesiatki = liczba % 100;
  const mnoga = jednosci >= 2 && jednosci <= 4 && (dziesiatki < 12 || dziesiatki > 14);
  return mnoga ? 'wersje' : 'wersji';
}

/** Otwiera w oknie dokument prowadzony w tej sesji; panel dokumentu nie zakłada, bo pokazuje historię pracy już prowadzonej. */
async function otworzDokument(kanal: Kanal, idOkna: string): Promise<StudioDocument | null> {
  const wynik = await wywolaj(kanal, Command.StudioDocumentOpen, { windowId: idOkna });
  if (!wynik.udany || wynik.wynik === undefined) return null;
  return wynik.wynik.document;
}

/** Przywraca wskazaną wersję dokumentu i zwraca dokument po przywróceniu. */
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
