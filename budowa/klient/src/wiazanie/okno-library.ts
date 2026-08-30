/**
 * Wiązanie wnętrza okna modułu Library z rdzeniem. Znacznik należy do
 * Właściciela — ten plik wypełnia stojące węzły odpowiedziami rodziny
 * library.*, powiela wzory zdjęte z treści przykładowej, a węzły, dla których
 * kontrakt nie ma pola, zdejmuje. Repozytorium jest wspólne dla konta, więc
 * żadna komenda rodziny nie bierze okna i wiązanie okna nie potrzebuje.
 */
import {
  Command,
  EventType,
  type LibraryCollection,
  type LibraryFile,
  type LibraryVersion,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import {
  data,
  opiszWezel,
  sklonuj,
  uzgodnijPrzelacznikiPaneli,
  wpiszTekst,
  zdejmijSterowanieWspolne,
  zdejmijTrescWspolna,
} from './okno-modulu.ts';

/** Panele prototypu, którym rodzina library.* nie oddaje ani jednego pola; schodzą razem z przełącznikami menu sesji. */
const PANELE_BEZ_POKRYCIA = [
  'panel-zadania',
  'panel-artefakty',
  'panel-pliki',
  'panel-kolejka',
  'panel-subagenci',
  'panel-terminal',
  'panel-przegladarka',
];

/** Wzory wierszy i pozycji zdjęte z treści przykładowej, zanim wykazy zostaną wyczyszczone. */
interface Wzory {
  grupa: HTMLElement | null;
  pozycja: HTMLElement | null;
  wiersz: HTMLElement | null;
  kolekcja: HTMLElement | null;
  etykieta: HTMLElement | null;
  wersja: HTMLElement | null;
}

/** Zasób wskazany w wykazie, jego wersja i kolekcja zawężająca; wszystkie wartości pochodzą z odpowiedzi rdzenia. */
interface Stan {
  plik: string;
  kolekcja: string;
  wersja: string;
  pliki: Map<string, LibraryFile>;
}

/**
 * Wiąże wnętrze okna modułu Library. Brak węzłów wnętrza kończy funkcję bez
 * żadnego działania — okna w ramie nie ma, więc nie ma czego wypełniać.
 */
export function zwiazLibrary(kanal: Kanal): void {
  const cialo = document.querySelector<HTMLElement>('.sta-cialo');
  if (cialo === null) return;

  const wzory = zdejmijWzory(cialo);
  zdejmijTrescPrzykladowa(cialo);
  zdejmijSterowanieBezPokrycia(cialo);
  wyczyscWskazanie(cialo);

  const stan: Stan = { plik: '', kolekcja: '', wersja: '', pliki: new Map() };

  /* Odwołanie stoi wartością, nie deklaracją: deklaracja wchodzi na wierzch
     zasięgu, więc kompilator bierze dla węzła wnętrza typ sprzed rozpoznania
     pustki. */
  /** Powtarza odczyt wykazu zasobów; wołane przy wejściu, po zawężeniu i po zmianie pliku w rdzeniu. */
  const odswiez = (): void => {
    void wypelnijWykaz(kanal, cialo, wzory, stan, frazaSzukania(cialo));
  };

  zwiazSzyne(kanal, cialo, wzory, stan, odswiez);
  zwiazSzukanie(cialo, odswiez);
  zwiazWykaz(kanal, cialo, wzory, stan);
  zwiazPrzywrocenie(kanal, cialo, wzory, stan, odswiez);
  zwiazZdarzeniaPlikow(kanal, odswiez);

  void wypelnijStanowisko(kanal, cialo, wzory, stan);
}

/** Zdejmuje wzory wierszy i pozycji; wzór traci oznaczenia stanu, których kontrakt zasobowi nie nadaje. */
function zdejmijWzory(cialo: HTMLElement): Wzory {
  const pozycja = sklonuj(cialo.querySelector('.dn-szyna-modulu-lista .pt-pozycja'));
  pozycja?.querySelector('.pt-tetno')?.remove();
  pozycja?.removeAttribute('aria-current');
  const wiersz = sklonuj(cialo.querySelector('.lb-tabela tbody tr'));
  // Kropka wiodąca i plakietka podejrzenia duplikatu nie mają pola w zasobie:
  // duplikaty rozstrzyga osobny przebieg, nie wykaz.
  wiersz?.querySelector('.dn-kropka')?.remove();
  // Kolumna autora schodzi także ze wzoru; wzór stoi przed czyszczeniem
  // wykazu, więc zdjęcie kolumny z tabeli samo go nie dosięga.
  wiersz?.querySelectorAll('td')[2]?.remove();
  const kolekcja = sklonuj(cialo.querySelector('#panel-tags .dn-wykaz-modulu-poz'));
  kolekcja?.querySelector('.pt-tetno')?.remove();
  kolekcja?.querySelector('.dn-kropka')?.remove();
  const wersja = sklonuj(cialo.querySelector('#panel-versioning .dn-wykaz-modulu-poz'));
  wersja?.querySelector('.dn-kropka')?.remove();
  return {
    grupa: sklonuj(cialo.querySelector('.dn-szyna-modulu-lista .dn-etyk-mono')),
    pozycja,
    wiersz,
    kolekcja,
    etykieta: sklonuj(cialo.querySelector('#panel-tags .dn-plakietka')),
    wersja,
  };
}

/**
 * Zdejmuje treść przykładową wnętrza. Wykaz rozmów i znaki kontekstu wiszą na
 * rodzinach spoza library.*, a nazwa repozytorium, przyrost dzienny i podpis
 * wersji nie mają pola w całym kontrakcie.
 */
function zdejmijTrescPrzykladowa(cialo: HTMLElement): void {
  zdejmijTrescWspolna(cialo);
  cialo.querySelector('.dn-szyna-modulu-lista')?.replaceChildren();
  for (const znak of cialo.querySelectorAll('.sta-kom-kontekst > .sta-zrodlo')) znak.remove();
  oznaczChipyCzynnosci(cialo);
  zdejmijKolumneAutora(cialo);
  cialo.querySelector('#panel-versioning .sta-okno-tresc > .dn-meta')?.remove();
}

/**
 * Opróżnia węzły opisujące zasób wskazany i kolekcję bieżącą. Wskazania przy
 * otwarciu okna nie ma, ale pola te kontrakt zna, więc węzły zostają puste
 * i czekają na wybór Operatora — schodzą wyłącznie węzły bez pokrycia.
 */
function wyczyscWskazanie(cialo: HTMLElement): void {
  opiszSciezke(cialo, '');
  for (const stojaca of cialo.querySelectorAll('#panel-versioning .dn-wykaz-modulu-poz')) {
    stojaca.remove();
  }
  for (const nazwa of ['#panel-preview', '#panel-versioning']) {
    const znacznik = cialo.querySelector<HTMLElement>(`${nazwa} .sta-okno-znacznik`);
    if (znacznik !== null) znacznik.textContent = '';
  }
  const podglad = cialo.querySelector<HTMLElement>('#panel-preview .sta-okno-tresc > div');
  if (podglad !== null) podglad.textContent = '';
  const opis = cialo.querySelector<HTMLElement>('#panel-preview .lb-mono');
  if (opis !== null) opis.textContent = '';
}

/**
 * Porządkuje pas czynności: znak zasięgu wykonania, nazwa repozytorium
 * i przyrost dzienny schodzą, a znaki liczby zasobów i kolekcji zostają
 * oznaczone do wypełnienia pulpitem stanu repozytorium.
 */
function oznaczChipyCzynnosci(cialo: HTMLElement): void {
  const chipy = [...cialo.querySelectorAll<HTMLElement>('.sta-kontekst-akcji > .sta-chip')];
  chipy[0]?.remove();
  chipy[1]?.remove();
  const pliki = chipy[2];
  const kolekcje = chipy[3];
  if (pliki !== undefined) pliki.dataset.pole = 'pliki';
  if (kolekcje !== undefined) kolekcje.dataset.pole = 'kolekcje';
  cialo.querySelector('.sta-kontekst-akcji .sta-chip--diff')?.remove();
  // Odłożenie zapytania jako wyszukiwania zapisanego nie ma komendy w rodzinie.
  cialo.querySelector('.sta-kontekst-akcji button.dn-btn')?.remove();
}

/** Zdejmuje z wykazu kolumnę autora wraz z jej komórkami; zasób repozytorium niesie moduł wytwórcy, autora — nie. */
function zdejmijKolumneAutora(cialo: HTMLElement): void {
  const tabela = cialo.querySelector('.lb-tabela');
  if (tabela === null) return;
  tabela.querySelectorAll('thead th')[2]?.remove();
  for (const wiersz of tabela.querySelectorAll('tbody tr')) wiersz.querySelectorAll('td')[2]?.remove();
}

/**
 * Zdejmuje sterowanie, którego rodzina library.* nie obsługuje: panele bez
 * pokrycia wraz z ich przełącznikami, wybór widoku wykazu, wciągnięcie pliku
 * bez jego treści, powiększenie podglądu i porównanie wersji bez miejsca na
 * wynik.
 */
function zdejmijSterowanieBezPokrycia(cialo: HTMLElement): void {
  uzgodnijPrzelacznikiPaneli(cialo, PANELE_BEZ_POKRYCIA);
  zdejmijSterowanieWspolne(cialo);
  cialo.querySelector('#panel-explorer .lb-zakladki')?.remove();
  // Założenie kolekcji bierze jej nazwę, a szyna pola nazwy nie niesie.
  cialo.querySelector('.dn-szyna-modulu-glowa .dn-btn')?.remove();
  // Wciągnięcie zasobu bierze jego treść, a pasek wykazu pola treści nie niesie.
  cialo.querySelector('#panel-explorer .lb-pasek .dn-btn')?.remove();
  cialo.querySelector('#panel-explorer .sta-okno-akcje .dn-btn-ikona')?.remove();
  const podglad = [...cialo.querySelectorAll('#panel-preview .sta-okno-tresc > div')];
  podglad[0]?.remove();
  cialo.querySelector('#panel-tags .dn-btn--zarys')?.remove();
  const wersje = [...cialo.querySelectorAll('#panel-versioning .dn-btn--zarys')];
  wersje[0]?.remove();
}

/** Wczytuje stan repozytorium: pulpit liczb, kolekcje szyny, słownik etykiet i wykaz zasobów. */
async function wypelnijStanowisko(
  kanal: Kanal,
  cialo: HTMLElement,
  wzory: Wzory,
  stan: Stan,
): Promise<void> {
  const kolekcje = await wywolaj(kanal, Command.LibraryCollectionList, {});
  const wykazKolekcji = kolekcje.udany ? (kolekcje.wynik?.collections ?? []) : [];
  const pulpit = await wywolaj(kanal, Command.LibraryStatsGet, {});
  opiszZnak(cialo, 'pliki', pulpit.udany ? String(pulpit.wynik?.stats.fileCount ?? '') : '');
  opiszZnak(cialo, 'kolekcje', String(wykazKolekcji.length));
  wypelnijSzyne(cialo, wzory, wykazKolekcji);
  wypelnijKolekcje(cialo, wzory, wykazKolekcji);
  await wypelnijEtykiety(kanal, cialo, wzory);
  await wypelnijWykaz(kanal, cialo, wzory, stan, '');
}

/** Nanosi wartość na znak pasa czynności rozpoznany po oznaczeniu. */
function opiszZnak(cialo: HTMLElement, pole: string, wartosc: string): void {
  opiszWezel(cialo.querySelector(`.sta-kontekst-akcji [data-pole="${pole}"]`), wartosc);
}

/** Stawia w szynie kolekcje repozytorium w dwóch grupach: zwykłych i wyznaczonych regułą. */
function wypelnijSzyne(cialo: HTMLElement, wzory: Wzory, kolekcje: LibraryCollection[]): void {
  const lista = cialo.querySelector<HTMLElement>('.dn-szyna-modulu-lista');
  if (lista === null || wzory.grupa === null || wzory.pozycja === null) return;
  const grupy: [string, LibraryCollection[]][] = [
    ['Kolekcje', kolekcje.filter((kolekcja) => (kolekcja.ruleId ?? '') === '')],
    ['Kolekcje inteligentne', kolekcje.filter((kolekcja) => (kolekcja.ruleId ?? '') !== '')],
  ];
  for (const [nazwa, zbior] of grupy) {
    if (zbior.length === 0) continue;
    const grupa = wzory.grupa.cloneNode(true) as HTMLElement;
    grupa.textContent = nazwa;
    lista.appendChild(grupa);
    for (const kolekcja of zbior) {
      const pozycja = wzory.pozycja.cloneNode(true) as HTMLElement;
      pozycja.dataset.kolekcja = kolekcja.id;
      const tytul = pozycja.querySelector('.pt-pozycja-tytul');
      if (tytul === null) continue;
      tytul.textContent = kolekcja.name;
      lista.appendChild(pozycja);
    }
  }
}

/** Wiąże wybór kolekcji w szynie z zawężeniem wykazu zasobów do tej kolekcji. */
function zwiazSzyne(
  kanal: Kanal,
  cialo: HTMLElement,
  wzory: Wzory,
  stan: Stan,
  odswiez: () => void,
): void {
  const lista = cialo.querySelector<HTMLElement>('.dn-szyna-modulu-lista');
  if (lista === null) return;
  lista.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const pozycja = cel.closest<HTMLElement>('.pt-pozycja');
    const kolekcja = pozycja?.dataset.kolekcja;
    if (pozycja === null || pozycja === undefined || kolekcja === undefined) return;
    for (const inna of lista.querySelectorAll('.pt-pozycja')) inna.removeAttribute('aria-current');
    pozycja.setAttribute('aria-current', 'true');
    stan.kolekcja = kolekcja;
    opiszSciezke(cialo, pozycja.querySelector('.pt-pozycja-tytul')?.textContent ?? '');
    void wypelnijWykaz(kanal, cialo, wzory, stan, frazaSzukania(cialo));
    odswiez();
  });
}

/** Nanosi nazwę kolekcji na ścieżkę wykazu i na znak okna komunikacji; wykaz całego repozytorium zostawia oba człony puste. */
function opiszSciezke(cialo: HTMLElement, nazwa: string): void {
  const czlon = cialo.querySelector<HTMLElement>('.pt-okruszki [aria-current]');
  if (czlon !== null) czlon.textContent = nazwa;
  opiszWezel(cialo.querySelector('#okno-czat-1 .sta-okno-belka > .sta-chip'), nazwa);
}

/** Stawia w panelu etykiet kolekcje wyznaczone regułą wraz z liczbą zasobów i nanosi liczbę kolekcji na znacznik panelu. */
function wypelnijKolekcje(cialo: HTMLElement, wzory: Wzory, kolekcje: LibraryCollection[]): void {
  const znacznik = cialo.querySelector<HTMLElement>('#panel-tags .sta-okno-znacznik');
  if (znacznik !== null) znacznik.textContent = String(kolekcje.length);
  const panel = cialo.querySelector<HTMLElement>('#panel-tags .sta-okno-tresc');
  if (panel === null || wzory.kolekcja === null) return;
  for (const stojaca of panel.querySelectorAll('.dn-wykaz-modulu-poz')) stojaca.remove();
  let kotwica = panel.querySelector('.pt-etykieta');
  for (const kolekcja of kolekcje.filter((pozycja) => (pozycja.ruleId ?? '') !== '')) {
    const wiersz = wzory.kolekcja.cloneNode(true) as HTMLElement;
    if (!wpiszTekst(wiersz, ` ${kolekcja.name} `)) continue;
    const meta = wiersz.querySelector('.dn-meta');
    if (meta !== null) meta.textContent = String(kolekcja.fileCount);
    kotwica?.after(wiersz);
    kotwica = wiersz;
  }
}

/** Stawia w słowniku etykiet po jednej plakietce na etykietę repozytorium wraz z licznikiem jej użycia. */
async function wypelnijEtykiety(kanal: Kanal, cialo: HTMLElement, wzory: Wzory): Promise<void> {
  const slownik = cialo.querySelector<HTMLElement>('#panel-tags .dn-plakietka')?.parentElement;
  if (slownik === null || slownik === undefined || wzory.etykieta === null) return;
  slownik.replaceChildren();
  const wynik = await wywolaj(kanal, Command.LibraryTagList, {});
  if (!wynik.udany || wynik.wynik === undefined) return;
  for (const etykieta of wynik.wynik.tags) {
    const plakietka = wzory.etykieta.cloneNode(true) as HTMLElement;
    plakietka.textContent = `${etykieta.name} · ${etykieta.fileCount}`;
    slownik.appendChild(plakietka);
  }
}

/** Fraza zawężania wpisana w pole wykazu; pole zdjęte albo puste daje wykaz całego repozytorium. */
function frazaSzukania(cialo: HTMLElement): string {
  return cialo.querySelector<HTMLInputElement>('#panel-explorer .dn-szukaj input')?.value.trim() ?? '';
}

/** Wiąże pole zawężania wykazu z wyszukiwaniem po treści i znaczeniu zasobu. */
function zwiazSzukanie(cialo: HTMLElement, odswiez: () => void): void {
  cialo
    .querySelector<HTMLInputElement>('#panel-explorer .dn-szukaj input')
    ?.addEventListener('input', odswiez);
}

/** Stawia w wykazie po jednym wierszu na zasób repozytorium zawężony kolekcją i frazą. */
async function wypelnijWykaz(
  kanal: Kanal,
  cialo: HTMLElement,
  wzory: Wzory,
  stan: Stan,
  fraza: string,
): Promise<void> {
  const tresc = cialo.querySelector<HTMLElement>('.lb-tabela tbody');
  if (tresc === null || wzory.wiersz === null) return;
  tresc.replaceChildren();
  const wynik = await wywolaj(kanal, Command.LibraryFileList, {
    query: fraza === '' ? undefined : fraza,
    collectionId: stan.kolekcja === '' ? undefined : stan.kolekcja,
  });
  stan.pliki.clear();
  if (!wynik.udany || wynik.wynik === undefined) return;
  for (const plik of wynik.wynik.files) {
    stan.pliki.set(plik.id, plik);
    const wiersz = wzory.wiersz.cloneNode(true) as HTMLElement;
    wiersz.dataset.plik = plik.id;
    opiszWierszPliku(wiersz, plik);
    tresc.appendChild(wiersz);
  }
}

/** Nanosi na wiersz wykazu nazwę zasobu, moduł wytwórcy, etykiety, czas zmiany i rozmiar; pole nieoddane zostawia komórkę pustą. */
function opiszWierszPliku(wiersz: HTMLElement, plik: LibraryFile): void {
  const komorki = [...wiersz.querySelectorAll<HTMLElement>('td')];
  if (komorki[0] !== undefined) komorki[0].textContent = plik.name;
  if (komorki[1] !== undefined) komorki[1].textContent = plik.sourceModuleId ?? '';
  const etykiety = komorki[2]?.querySelector('.dn-plakietka');
  const nazwy = plik.tags ?? [];
  if (etykiety !== null && etykiety !== undefined) {
    if (nazwy.length === 0) etykiety.remove();
    else etykiety.textContent = nazwy.join(', ');
  }
  if (komorki[3] !== undefined) komorki[3].textContent = data(plik.updatedAt);
  if (komorki[4] !== undefined) komorki[4].textContent = zapiszRozmiar(plik.sizeBytes);
}

/** Rzędy wielkości rozmiaru zasobu, od bajta w górę. */
const JEDNOSTKI_ROZMIARU = ['B', 'kB', 'MB', 'GB', 'TB'];

/** Rozmiar zasobu w zapisie znacznika: bajty rdzenia przeliczone na rząd tysięczny; brak pola daje pustkę. */
function zapiszRozmiar(bajty: number | undefined): string {
  if (bajty === undefined) return '';
  let wartosc = bajty;
  let rzad = 0;
  while (wartosc >= 1000 && rzad < JEDNOSTKI_ROZMIARU.length - 1) {
    wartosc /= 1000;
    rzad += 1;
  }
  return `${wartosc.toLocaleString('pl-PL', { maximumFractionDigits: 1 })} ${JEDNOSTKI_ROZMIARU[rzad]}`;
}

/** Wiąże wybór wiersza wykazu z podglądem zasobu i jego wykazem wersji. */
function zwiazWykaz(kanal: Kanal, cialo: HTMLElement, wzory: Wzory, stan: Stan): void {
  const tresc = cialo.querySelector<HTMLElement>('.lb-tabela tbody');
  if (tresc === null) return;
  tresc.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const wiersz = cel.closest<HTMLElement>('tr');
    const plik = wiersz?.dataset.plik;
    if (wiersz === null || wiersz === undefined || plik === undefined) return;
    for (const inny of tresc.querySelectorAll('tr')) inny.removeAttribute('aria-selected');
    wiersz.setAttribute('aria-selected', 'true');
    stan.plik = plik;
    void wypelnijPodglad(kanal, cialo, stan.pliki.get(plik));
    void wypelnijWersje(kanal, cialo, wzory, stan);
  });
}

/**
 * Wczytuje podgląd zasobu wraz z jego opisem. Opis składa się z rodzaju
 * treści, modułu wytwórcy i etykiet; karta sesji, przy której zasób powstał,
 * nie ma pola w zasobie repozytorium, więc do opisu nie wchodzi.
 */
async function wypelnijPodglad(
  kanal: Kanal,
  cialo: HTMLElement,
  plik: LibraryFile | undefined,
): Promise<void> {
  const znacznik = cialo.querySelector<HTMLElement>('#panel-preview .sta-okno-znacznik');
  if (znacznik !== null) znacznik.textContent = plik?.name ?? '';
  const pole = cialo.querySelector<HTMLElement>('#panel-preview .sta-okno-tresc > div');
  const opis = cialo.querySelector<HTMLElement>('#panel-preview .lb-mono');
  if (plik === undefined) return;
  const wynik = await wywolaj(kanal, Command.LibraryFilePreview, { fileId: plik.id });
  if (pole !== null) pole.textContent = wynik.udany ? (wynik.wynik?.preview.text ?? '') : '';
  if (opis === null) return;
  const czlony = [plik.mimeType ?? '', plik.sourceModuleId ?? '', (plik.tags ?? []).join(', ')];
  opis.textContent = czlony.filter((czlon) => czlon !== '').join(' · ');
}

/** Stawia w panelu wersji po jednym wierszu na wersję zasobu i nanosi ich liczbę na znacznik panelu. */
async function wypelnijWersje(
  kanal: Kanal,
  cialo: HTMLElement,
  wzory: Wzory,
  stan: Stan,
): Promise<void> {
  const panel = cialo.querySelector<HTMLElement>('#panel-versioning .sta-okno-tresc');
  if (panel === null || wzory.wersja === null) return;
  for (const stojaca of panel.querySelectorAll('.dn-wykaz-modulu-poz')) stojaca.remove();
  const wynik = await wywolaj(kanal, Command.LibraryVersionList, { fileId: stan.plik });
  const wersje = wynik.udany ? (wynik.wynik?.versions ?? []) : [];
  const znacznik = cialo.querySelector<HTMLElement>('#panel-versioning .sta-okno-znacznik');
  if (znacznik !== null) znacznik.textContent = String(wersje.length);
  const przyciski = panel.querySelector('div');
  for (const wersja of wersje) {
    const wiersz = postawWierszWersji(wzory.wersja, wersja);
    if (wiersz === null) continue;
    if (przyciski === null) panel.appendChild(wiersz);
    else przyciski.before(wiersz);
  }
}

/** Składa wiersz wersji z jej podpisu oraz autora i daty założenia; wersja bez podpisu wchodzi identyfikatorem rdzenia. */
function postawWierszWersji(wzor: HTMLElement, wersja: LibraryVersion): HTMLElement | null {
  const wiersz = wzor.cloneNode(true) as HTMLElement;
  wiersz.dataset.wersja = wersja.id;
  if (!wpiszTekst(wiersz, ` ${wersja.label ?? wersja.id} `)) return null;
  const meta = wiersz.querySelector('.dn-meta');
  if (meta === null) return wiersz;
  const autor = wersja.author ?? '';
  meta.textContent = autor === '' ? data(wersja.createdAt) : `${autor} · ${data(wersja.createdAt)}`;
  return wiersz;
}

/** Wiąże przywrócenie wersji: wybór wiersza wskazuje wersję, przycisk panelu odkłada ją jako obowiązującą. */
function zwiazPrzywrocenie(
  kanal: Kanal,
  cialo: HTMLElement,
  wzory: Wzory,
  stan: Stan,
  odswiez: () => void,
): void {
  const panel = cialo.querySelector<HTMLElement>('#panel-versioning .sta-okno-tresc');
  if (panel === null) return;
  panel.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const wiersz = cel.closest<HTMLElement>('.dn-wykaz-modulu-poz');
    const wersja = wiersz?.dataset.wersja;
    if (wiersz === null || wiersz === undefined || wersja === undefined) return;
    for (const inny of panel.querySelectorAll('.dn-wykaz-modulu-poz')) {
      inny.removeAttribute('aria-current');
    }
    wiersz.setAttribute('aria-current', 'true');
    stan.wersja = wersja;
  });
  panel.querySelector('.dn-btn--zarys')?.addEventListener('click', () => {
    if (stan.plik === '' || stan.wersja === '') return;
    void wywolaj(kanal, Command.LibraryVersionRestore, {
      fileId: stan.plik,
      versionId: stan.wersja,
    }).then(() => {
      void wypelnijWersje(kanal, cialo, wzory, stan);
      odswiez();
    });
  });
}

/** Nasłuchuje zmian zasobów repozytorium: wykaz bierze stan po zmianie, bo zasób może wejść z innego modułu. */
function zwiazZdarzeniaPlikow(kanal: Kanal, odswiez: () => void): void {
  kanal.naZdarzenie(EventType.LibraryFileChanged, () => {
    odswiez();
  });
}
