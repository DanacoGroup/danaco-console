/**
 * Wiązanie panelu Pliki okna Studia z repozytorium rdzenia. Znacznik należy do
 * Właściciela — ten plik wypełnia stojące węzły, a wiersze wykazu powiela
 * z wzoru zdjętego z treści przykładowej.
 */
import { Command, EventType, type LibraryFile } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';

/** Węzły panelu Pliki, na których wiązanie pracuje; brak któregokolwiek znaczy, że panel nie stoi w dokumencie. */
interface WezlyPlikow {
  lista: HTMLElement;
  szukaj: HTMLInputElement;
}

/** Jednostki rozmiaru w brzmieniu znacznika prototypu, w rzędach tysięcznych. */
const JEDNOSTKI_ROZMIARU = ['B', 'kB', 'MB', 'GB', 'TB'];

/**
 * Wiąże panel Pliki okna Studia z wykazem repozytorium. Zwraca prawdę, gdy
 * znacznik panelu stał i wiązanie zostało założone.
 */
export function zwiazPliki(kanal: Kanal, idOkna: string): boolean {
  const znalezione = zbierzWezly();
  if (znalezione === null) return false;
  const wezly: WezlyPlikow = znalezione;

  const wzor = sklonuj(wezly.lista.querySelector('.st-panel-wiersz'));
  // Treść przykładowa znika, zanim padnie pierwsza odpowiedź rdzenia: pusty
  // wykaz jest uczciwy, wykaz z cudzymi plikami — nie.
  zdejmijWiersze(wezly.lista);

  const pliki = new Map<HTMLElement, string>();
  let idProjektu = '';
  let numerZadania = 0;

  /** Wczytuje wykaz zawężony frazą pola i wstawia go w miejsce wierszy poprzednich; odpowiedź spóźniona wobec kolejnej odpada. */
  async function odswiez(): Promise<void> {
    if (wzor === null) return;
    const numer = (numerZadania += 1);
    const wykaz = await pobierzPliki(kanal, idProjektu, wezly.szukaj.value.trim());
    if (numer !== numerZadania) return;
    zdejmijWiersze(wezly.lista);
    pliki.clear();
    for (const plik of wykaz) {
      const wiersz = zbudujWiersz(wzor, plik);
      pliki.set(wiersz, plik.id);
      wezly.lista.appendChild(wiersz);
    }
  }

  wezly.szukaj.addEventListener('input', () => {
    void odswiez();
  });

  wezly.lista.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const wiersz = cel.closest('.st-panel-wiersz');
    if (!(wiersz instanceof HTMLElement)) return;
    const idPliku = pliki.get(wiersz);
    if (idPliku === undefined) return;
    void otworzPlik(kanal, idOkna, idPliku).then((otwarty) => {
      if (!otwarty) return;
      for (const inny of wezly.lista.querySelectorAll('.st-panel-wiersz')) {
        inny.removeAttribute('aria-current');
      }
      wiersz.setAttribute('aria-current', 'true');
    });
  });

  kanal.naZdarzenie(EventType.LibraryFileChanged, () => {
    void odswiez();
  });

  void wskazProjekt(kanal, idOkna).then((znaleziony) => {
    idProjektu = znaleziony;
    void odswiez();
  });

  return true;
}

/** Odczytuje projekt sesji, do której należy okno; pustka znaczy sesję bez projektu, dla której biblioteka projektu nie ma czym się zawęzić. */
async function wskazProjekt(kanal: Kanal, idOkna: string): Promise<string> {
  if (idOkna === '') return '';
  const okna = await wywolaj(kanal, Command.WindowList, {});
  if (!okna.udany || okna.wynik === undefined) return '';
  const okno = okna.wynik.windows.find((pozycja) => pozycja.id === idOkna);
  if (okno === undefined) return '';
  const sesje = await wywolaj(kanal, Command.SessionList, {});
  if (!sesje.udany || sesje.wynik === undefined) return '';
  return sesje.wynik.sessions.find((pozycja) => pozycja.id === okno.sessionId)?.projectId ?? '';
}

/** Wykaz plików: z biblioteki projektu, gdy sesja okna ma projekt, w przeciwnym razie z repozytorium bez zawężenia projektem. */
async function pobierzPliki(
  kanal: Kanal,
  idProjektu: string,
  fraza: string,
): Promise<LibraryFile[]> {
  const szukana = fraza === '' ? undefined : fraza;
  if (idProjektu !== '') {
    const wykaz = await wywolaj(kanal, Command.WorkspaceLibraryList, {
      projectId: idProjektu,
      query: szukana,
    });
    return wykaz.udany && wykaz.wynik !== undefined ? wykaz.wynik.files : [];
  }
  const wykaz = await wywolaj(kanal, Command.LibraryFileList, { query: szukana });
  return wykaz.udany && wykaz.wynik !== undefined ? wykaz.wynik.files : [];
}

/** Wczytuje plik repozytorium do Studia; prawda znaczy, że rdzeń dokument oddał. */
async function otworzPlik(kanal: Kanal, idOkna: string, idPliku: string): Promise<boolean> {
  if (idOkna === '') return false;
  const wynik = await wywolaj(kanal, Command.StudioDocumentOpen, {
    windowId: idOkna,
    libraryFileId: idPliku,
  });
  return wynik.udany && wynik.wynik !== undefined;
}

/** Zwraca klon wzoru wiersza opisany nazwą pliku; miara przy nazwie zostaje wyłącznie dla pliku o znanym rozmiarze. */
function zbudujWiersz(wzor: HTMLElement, plik: LibraryFile): HTMLElement {
  const wiersz = wzor.cloneNode(true) as HTMLElement;
  // Nazwa stoi we wzorze jako tekst przed miarą, bez własnego węzła
  // elementowego, więc podmiana idzie po pierwszym węźle tekstowym.
  const nazwa = [...wiersz.childNodes].find((wezel) => wezel.nodeType === Node.TEXT_NODE);
  if (nazwa !== undefined) nazwa.nodeValue = `${plik.name} `;
  const miara = wiersz.querySelector('.dn-meta');
  if (miara === null) return wiersz;
  if (plik.sizeBytes === undefined) miara.remove();
  else miara.textContent = zapiszRozmiar(plik.sizeBytes);
  return wiersz;
}

/** Rozmiar pliku w zapisie znacznika: bajty rdzenia przeliczone na rząd tysięczny, z jednym miejscem po przecinku. */
function zapiszRozmiar(bajty: number): string {
  let wartosc = bajty;
  let rzad = 0;
  while (wartosc >= 1000 && rzad < JEDNOSTKI_ROZMIARU.length - 1) {
    wartosc /= 1000;
    rzad += 1;
  }
  return `${wartosc.toLocaleString('pl-PL', { maximumFractionDigits: 1 })} ${JEDNOSTKI_ROZMIARU[rzad]}`;
}

/** Zdejmuje wiersze wykazu, zostawiając pole zawężania, które pozycją pliku nie jest. */
function zdejmijWiersze(lista: HTMLElement): void {
  for (const wiersz of lista.querySelectorAll('.st-panel-wiersz')) wiersz.remove();
}

/** Wskazuje węzły panelu Pliki; pustka znaczy, że panel nie stoi w dokumencie. */
function zbierzWezly(): WezlyPlikow | null {
  const panel = document.querySelector('#panel-pliki');
  const lista = panel?.querySelector('.sta-okno-tresc.st-panel-lista');
  const szukaj = panel?.querySelector('input[type="search"]');
  if (!(lista instanceof HTMLElement) || !(szukaj instanceof HTMLInputElement)) return null;
  return { lista, szukaj };
}

/** Klon węzła wzorcowego, odporny na jego brak w znaczniku. */
function sklonuj(wezel: Element | null): HTMLElement | null {
  return wezel instanceof HTMLElement ? (wezel.cloneNode(true) as HTMLElement) : null;
}
