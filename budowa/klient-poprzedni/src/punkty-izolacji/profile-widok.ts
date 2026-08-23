import type { IsolationProfile } from '../../../shared/contract';
import {
  POZYCJE_KONTEKSTU,
  POZYCJE_TECHNICZNE,
  etykietaKontekstu,
  etykietaTechniczna,
} from './katalog-izolacji';

/**
 * Rysowanie profili izolacji: wykaz zapisanych i podgląd wczytanego.
 *
 * Przy każdym wierszu stoją trzy czynności i dwie z nich łatwo pomylić.
 * „Wczytaj" pobiera profil do podglądu i do formularza — żaden poziom zasięgu
 * po nim nie działa inaczej. „Przypisz do poziomu" dopiero wiąże profil
 * z wybranym poziomem i warstwą, i to po niej izolacja faktycznie się zmienia.
 * Obie stoją obok siebie, więc różnica jest wypisana słowami przy każdym
 * wierszu, a nie domyślana z kolejności przycisków.
 *
 * „Usuń" usuwa od razu — bez pytania „czy na pewno", bez wygaszania, bez
 * uprawnień. Profil izolacji rozstrzyga, co model widzi z sąsiedniego okna,
 * i zmienia się na żądanie; skutek usunięcia stoi wprost przy przycisku.
 *
 * Ten plik nie woła rdzenia. Buduje węzły i oddaje naciśnięcia wywołującemu
 * (`obszar-profile.ts`); nie zna ani jednej komendy, ani jednej barwy.
 */

/** Czynności podpięte pod przyciski wiersza; każdą wykonuje `obszar-profile.ts`. */
export interface CzynnosciProfilu {
  wczytaj(profil: IsolationProfile): void;
  przypisz(profil: IsolationProfile): void;
  usun(profil: IsolationProfile): void;
}

/** Wykaz zapisanych profili; jeden wiersz na profil, trzy czynności przy każdym. */
export function zbudujListeProfili(
  profile: readonly IsolationProfile[],
  czynnosci: CzynnosciProfilu,
  opisCelu: string,
): HTMLElement {
  const lista = document.createElement('ul');
  lista.className = 'pi-profile__lista';
  lista.setAttribute('aria-label', 'Zapisane profile izolacji');
  lista.append(...profile.map((profil) => zbudujWiersz(profil, czynnosci, opisCelu)));
  return lista;
}

function zbudujWiersz(
  profil: IsolationProfile,
  czynnosci: CzynnosciProfilu,
  opisCelu: string,
): HTMLElement {
  const tytul = document.createElement('strong');
  tytul.className = 'pi-profile__tytul';
  tytul.textContent = profil.name === '' ? '(profil bez nazwy)' : profil.name;

  const identyfikator = document.createElement('code');
  identyfikator.textContent = profil.id;

  const opis = document.createElement('span');
  opis.className = 'pi-profile__opis';
  opis.textContent =
    profil.description === undefined || profil.description === ''
      ? 'Bez opisu — rdzeń nie niesie przy tym profilu żadnego zdania.'
      : profil.description;

  const zestaw = document.createElement('span');
  zestaw.className = 'pi-profile__zestaw';
  zestaw.textContent = streszczenieZestawu(profil);

  const akcje = document.createElement('span');
  akcje.className = 'pi-profile__akcje';
  akcje.append(
    przycisk('Wczytaj', 'dn-btn dn-btn--zarys', 'Pobiera zestaw do podglądu i do formularza. NIE przypisuje go do żadnego poziomu — izolacja po tym nie zmienia się ani o krok.', () =>
      czynnosci.wczytaj(profil),
    ),
    przycisk('Przypisz do poziomu', 'dn-btn', `Wiąże ten profil z celem: ${opisCelu}. Dopiero ta czynność zmienia obowiązującą izolację.`, () =>
      czynnosci.przypisz(profil),
    ),
    przycisk('Usuń', 'dn-btn dn-btn--zarys', 'Usuwa profil od razu, bez pytania. Przypisania tego profilu przestają go wskazywać; zestaw możesz zapisać z powrotem formularzem poniżej.', () =>
      czynnosci.usun(profil),
    ),
  );

  const element = document.createElement('li');
  element.className = 'pi-profile__pozycja';
  element.append(tytul, identyfikator, opis, zestaw, akcje);
  return element;
}

function przycisk(etykieta: string, klasa: string, wyjasnienie: string, naKlik: () => void): HTMLButtonElement {
  const kontrolka = document.createElement('button');
  kontrolka.type = 'button';
  kontrolka.className = klasa;
  kontrolka.textContent = etykieta;
  kontrolka.title = wyjasnienie;
  kontrolka.setAttribute('aria-description', wyjasnienie);
  kontrolka.addEventListener('click', naKlik);
  return kontrolka;
}

/** Jednozdaniowe streszczenie zestawu: ile kluczy odciętych z trzech i z ośmiu. */
function streszczenieZestawu(profil: IsolationProfile): string {
  const odrebne = profil.contextSwitches.filter((p) => p.isolated).length;
  const wlaczone = profil.technicalSwitches.filter((p) => p.isolated).length;
  return (
    `Kontekst: ${odrebne} z ${POZYCJE_KONTEKSTU.length} odrębnych. ` +
    `Zakres techniczny: ${wlaczone} z ${POZYCJE_TECHNICZNE.length} włączonych.`
  );
}

/**
 * Podgląd profilu wczytanego (`isolation.profile.load`) — komplet jedenastu
 * kluczy z wartościami tego profilu i jawne zdanie, że to jest podgląd, a nie
 * stan obowiązujący.
 */
export function zbudujPodgladProfilu(profil: IsolationProfile): HTMLElement {
  const naglowek = document.createElement('h4');
  naglowek.className = 'pi-profile__naglowek';
  naglowek.textContent = `Wczytany profil: ${profil.name === '' ? '(bez nazwy)' : profil.name}`;

  const zastrzezenie = document.createElement('p');
  zastrzezenie.className = 'pi-profile__opis';
  zastrzezenie.textContent =
    'To jest ZAWARTOŚĆ profilu, nie stan obowiązujący. Żaden poziom zasięgu nie działa dziś według ' +
    'tego zestawu, dopóki nie naciśniesz „Przypisz do poziomu”. Wartości obowiązujące pokazuje ' +
    'zakładka „Polityka efektywna”.';

  const tabela = document.createElement('table');
  tabela.className = 'dn-tabela';
  tabela.setAttribute('aria-label', 'Jedenaście przełączników wczytanego profilu izolacji');

  const glowa = document.createElement('thead');
  const wierszGlowy = document.createElement('tr');
  for (const etykieta of ['Punkt izolacji', 'Wartość w profilu']) {
    const komorka = document.createElement('th');
    komorka.scope = 'col';
    komorka.textContent = etykieta;
    wierszGlowy.append(komorka);
  }
  glowa.append(wierszGlowy);

  const cialo = document.createElement('tbody');
  for (const pozycja of POZYCJE_KONTEKSTU) {
    const znaleziony = profil.contextSwitches.find((p) => p.kind === pozycja.kind);
    cialo.append(zbudujWierszPodgladu(pozycja.nazwa, znaleziony === undefined ? undefined : etykietaKontekstu(znaleziony.isolated)));
  }
  for (const pozycja of POZYCJE_TECHNICZNE) {
    const znaleziony = profil.technicalSwitches.find((p) => p.scope === pozycja.scope);
    cialo.append(zbudujWierszPodgladu(pozycja.nazwa, znaleziony === undefined ? undefined : etykietaTechniczna(znaleziony.isolated)));
  }

  tabela.append(glowa, cialo);

  const element = document.createElement('div');
  element.className = 'pi-profile__podglad';
  element.append(naglowek, zastrzezenie, tabela);
  return element;
}

/** Wiersz podglądu; brak klucza w profilu nazwany wprost, nie domalowany domyślną wartością. */
function zbudujWierszPodgladu(nazwa: string, wartosc: string | undefined): HTMLElement {
  const element = document.createElement('tr');

  const komorkaNazwy = document.createElement('td');
  komorkaNazwy.textContent = nazwa;

  const komorkaWartosci = document.createElement('td');
  const plakietka = document.createElement('span');
  if (wartosc === undefined) {
    plakietka.className = 'dn-plakietka dn-plakietka--ostrzezenie';
    plakietka.textContent = 'Profil nie niesie tego klucza';
    plakietka.title = 'Rdzeń nie zwrócił tego przełącznika w tym profilu — to nie jest wartość domyślna, tylko jej brak.';
  } else {
    plakietka.className = 'dn-plakietka dn-plakietka--informacja';
    plakietka.textContent = wartosc;
  }
  komorkaWartosci.append(plakietka);

  element.append(komorkaNazwy, komorkaWartosci);
  return element;
}
