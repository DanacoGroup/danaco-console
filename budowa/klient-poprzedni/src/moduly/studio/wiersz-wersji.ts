import { StudioAuthor, type StudioVersion } from '../../../../shared/contract';

/** Interfejs WierszWersji niesie wiersz historii wersji dokumentu wraz z przyciskami: podgląd, przywrócenie, porównanie oraz przyciski menu pozycji. */
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

/** Interfejs PostacWiersza niesie postać wiersza wersji znaną oknu historii, a nie samej wersji dokumentu: czy jest bieżąca i czy pochodzi z zapisu samoczynnego. */
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

    // Oznacza wersję jako nieosiągalną z powodem odmowy rdzenia i wyłącza przywrócenie oraz podgląd.
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

/** Funkcja nazwaAutora zwraca nazwę autora wersji dokumentu; brak pola autora oddaje kreskę, a nie nazwę roli użytkownika. */
function nazwaAutora(autor: StudioAuthor | undefined): string {
  if (autor === StudioAuthor.Model) return 'model';
  if (autor === StudioAuthor.Uzytkownik) return 'Operator';
  return '—';
}
