import { StudioAuthor, type StudioVersion } from '../../../../shared/contract';

/**
 * Wiersz historii wersji — czas, autor, kropka stanu, etykieta i czynności.
 *
 * ── Kropka stanu, nie sama nazwa ────────────────────────────────────────────
 * Wersja bieżąca jest wyróżniona kropką pełną, archiwalne mają kropkę pustą.
 * Stan nigdy samą barwą: kropce towarzyszy słowo, bo barwa bez słowa nie mówi
 * nic Operatorowi, który jej nie rozróżnia.
 *
 * ── Trzy czynności i menu pozycji ───────────────────────────────────────────
 * Podgląd, Przywróć, Porównaj stoją wprost, bo po nie sięga się najczęściej.
 * Etykieta, eksport, odwołanie i schowanie miejscowe siedzą w menu pozycji
 * rozwijanym `⋯` — wiersz z siedmioma przyciskami byłby nieczytelny, a te cztery
 * są czynnościami rzadkimi.
 *
 * ── Czego tu nie ma ─────────────────────────────────────────────────────────
 * Czynności „Rozgałęź" — decyzją Właściciela gałęzie są poza zakresem tej tury.
 * Komendy `studio.branch.*` pracują w rdzeniu dalej.
 *
 * Plik nie woła rdzenia: oddaje przyciski, a rozmowę prowadzi okno historii.
 */

/** Wiersz wersji wraz z jego przyciskami. */
export interface WierszWersji {
  element: HTMLElement;
  podejrzyj: HTMLButtonElement;
  przywroc: HTMLButtonElement;
  porownaj: HTMLButtonElement;
  /** Menu pozycji: nazwa własna wersji. */
  etykieta: HTMLButtonElement;
  eksportuj: HTMLButtonElement;
  odwolanie: HTMLButtonElement;
  /** Schowanie pozycji miejscowo — nie usunięcie wersji w rdzeniu. */
  usun: HTMLButtonElement;
  oznaczNiedostepna(powod: string): void;
}

/** Postać wiersza znana oknu, a nie samej wersji. */
export interface PostacWiersza {
  /** Czy ta wersja jest wersją bieżącą dokumentu. */
  biezaca: boolean;
  /** Czy wersja pochodzi z zapisu samoczynnego. */
  zapisSamoczynny: boolean;
}

export function utworzWierszWersji(
  wersja: StudioVersion,
  postac: PostacWiersza = { biezaca: false, zapisSamoczynny: false },
): WierszWersji {
  const element = document.createElement('tr');
  element.dataset['wersja'] = wersja.id;
  element.dataset['dostepnosc'] = 'osiagalna';
  element.dataset['biezaca'] = postac.biezaca ? 'tak' : 'nie';
  element.dataset['samoczynna'] = postac.zapisSamoczynny ? 'tak' : 'nie';

  const nazwa = document.createElement('td');
  nazwa.textContent =
    wersja.label === undefined || wersja.label === '' ? wersja.id : wersja.label;
  nazwa.title = `Identyfikator wersji: ${wersja.id}`;
  if (wersja.milestone === true) {
    const kluczowa = document.createElement('span');
    kluczowa.className = 'dn-plakietka';
    kluczowa.textContent = 'kluczowa';
    kluczowa.title =
      'Wersja oznaczona jako kluczowa (studio.version.label.set). Zapisy samoczynne nigdy nie ' +
      'są kluczowe, więc to oznaczenie odróżnia szereg Operatora od szeregu autozapisu.';
    nazwa.append(' ', kluczowa);
  }

  const opis = document.createElement('td');
  opis.textContent = wersja.summary === undefined || wersja.summary === '' ? '—' : wersja.summary;

  const chwila = document.createElement('td');
  chwila.className = 'ms-repozytorium__chwila';
  chwila.textContent = new Date(wersja.createdAt).toLocaleString('pl');

  const autor = document.createElement('td');
  autor.textContent = nazwaAutora(wersja.author);
  autor.title =
    wersja.author === undefined
      ? 'Wersja założona przed wprowadzeniem pola autora — rdzeń go dla niej nie zna. Puste ' +
        'miejsce jest tu uczciwsze niż podstawienie Operatora.'
      : `Autor wersji wedle pola StudioVersion.author: ${wersja.author}.`;

  const stanKomorka = document.createElement('td');
  const kropka = document.createElement('span');
  kropka.className = 'ms-repozytorium__kropka';
  kropka.dataset['stan'] = postac.biezaca ? 'biezaca' : 'archiwalna';
  kropka.textContent = postac.biezaca ? '●' : '○';
  kropka.setAttribute('aria-hidden', 'true');

  const plakietka = document.createElement('span');
  plakietka.className = `dn-plakietka ${postac.biezaca ? 'dn-plakietka--sukces' : ''}`.trim();
  plakietka.textContent = postac.biezaca
    ? 'bieżąca'
    : postac.zapisSamoczynny
      ? 'zapis samoczynny'
      : 'archiwalna';
  stanKomorka.append(kropka, ' ', plakietka);

  const podejrzyj = przyciskWiersza('Podgląd', 'podejrzyj',
    'Pokazuje treść tej wersji jako różnicę wobec bieżącej, bez jej przywracania.');
  const przywroc = przyciskWiersza('Przywróć', 'przywroc',
    'Przywraca treść tej wersji i zakłada stan bieżący. Wersje nowsze POZOSTAJĄ w historii, ' +
    'więc przywrócenie samo da się cofnąć.');
  const porownaj = przyciskWiersza('Porównaj', 'porownaj',
    'Wskazuje tę wersję jako odniesienie różnicy wobec bieżącej.');

  const etykieta = przyciskWiersza('Etykieta', 'etykieta',
    'Nadaje wersji nazwę własną i oznacza ją jako kluczową — studio.version.label.set.');
  const eksportuj = przyciskWiersza('Eksportuj', 'eksportuj',
    'Wydaje tę wersję z historii jako archiwum — studio.repository.export.');
  const odwolanie = przyciskWiersza('Odwołanie', 'odwolanie',
    'Pokazuje odwołanie do wersji: identyfikator, dokument i czas — do wklejenia w pismo.');
  const usun = przyciskWiersza('Usuń z wykazu', 'usun',
    'Schowanie MIEJSCOWE, w tym oknie. Wersja zostaje w rdzeniu — historia nic nie usuwa bez ' +
    'decyzji Operatora, a komendy usuwającej wersję kontrakt nie niesie.');

  const menu = document.createElement('details');
  menu.className = 'ms-repozytorium__menu';
  const znak = document.createElement('summary');
  znak.textContent = '⋯';
  znak.title = 'Menu pozycji: etykieta, eksport, odwołanie, schowanie miejscowe.';
  const pas = document.createElement('div');
  pas.className = 'ms-repozytorium__menu-pas';
  pas.append(etykieta, eksportuj, odwolanie, usun);
  menu.append(znak, pas);

  const czynnosci = document.createElement('td');
  czynnosci.className = 'ms-repozytorium__czynnosci';
  czynnosci.append(podejrzyj, przywroc, porownaj, menu);

  element.append(nazwa, opis, chwila, autor, stanKomorka, czynnosci);

  return {
    element,
    podejrzyj,
    przywroc,
    porownaj,
    etykieta,
    eksportuj,
    odwolanie,
    usun,

    /**
     * Oznacza wersję jako nieosiągalną wraz z powodem odmowy rdzenia.
     *
     * Przyciski schodzą, bo powtórzenie tej samej odmowy niczego nie zmieni;
     * powód zostaje przy wierszu, żeby Operator wiedział, czego dotyczyła.
     */
    oznaczNiedostepna(powod) {
      element.dataset['dostepnosc'] = 'niedostepna';
      plakietka.className = 'dn-plakietka dn-plakietka--blad';
      plakietka.textContent = 'nieosiągalna';
      plakietka.title = powod;
      kropka.dataset['stan'] = 'niedostepna';
      przywroc.disabled = true;
      podejrzyj.disabled = true;
    },
  };
}

function przyciskWiersza(nazwa: string, kod: string, objasnienie: string): HTMLButtonElement {
  const element = document.createElement('button');
  element.type = 'button';
  element.className = 'dn-btn dn-btn--sm dn-btn--duch';
  element.textContent = nazwa;
  element.dataset['czynnosc'] = kod;
  element.title = objasnienie;
  element.setAttribute('aria-description', objasnienie);
  return element;
}

/** Nazwa autora wersji; brak pola oddaje kreskę, a nie Operatora. */
function nazwaAutora(autor: StudioAuthor | undefined): string {
  if (autor === StudioAuthor.Model) return 'model';
  if (autor === StudioAuthor.Uzytkownik) return 'Operator';
  return '—';
}
