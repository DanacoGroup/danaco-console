import { ClipboardEntryKind, type ClipboardEntry } from '../../../../shared/contract';
import { schowekOpiszWpis, schowekPodglad } from './schowek-zrodlo';

/**
 * Historia schowka jako wykaz do wyboru — nie samo „wklej ostatnie".
 *
 * ── Trzy wymagania Właściciela w jednym panelu ──────────────────────────────
 * Wykaz wpisów sprzed kilku ruchów, przypinanie wpisów, które mają zostać na
 * stałe, i wklejanie w dwóch postaciach: z zachowaniem postaci albo jako czysty
 * tekst. Ostatnie jest jawnym wyborem Operatora i oba warianty są równorzędne —
 * przycisk jest jeden dla każdego, a nie jeden domyślny z ukrytym drugim.
 *
 * ── Malarz formatów stoi tu, nie w schowku rdzenia ──────────────────────────
 * Malarz kopiuje POSTAĆ, nie treść, więc nie jest wpisem `clipboard.*`:
 * kontrakt schowka niesie `content` i rodzaj (tekst, obraz, ścieżka pliku),
 * a nie arkusz nastaw akapitu. Malarz trzyma więc pobraną postać w tym panelu
 * i oddaje ją oknu, które jedyne wie, czym postać akapitu jest. Gdy postać
 * dokumentu wejdzie do kontraktu (rodzina `studio.document.*` postaci —
 * odcinek 1 i 4), malarz przejdzie na komendę rdzenia i przetrwa zamknięcie
 * karty; dziś żyje przez sesję okna i panel mówi to wprost.
 */

/** Czynności panelu schowka zlecane oknu. */
export interface CzynnosciSchowka {
  /** Wkleja treść wpisu w miejsce kursora z zachowaniem postaci akapitu. */
  naWklejenieZPostacia(tresc: string): void;
  /** Wkleja treść wpisu jako czysty tekst — bez postaci akapitu. */
  naWklejenieCzyste(tresc: string): void;
  /** Odkłada zaznaczony fragment dokumentu do historii schowka. */
  naOdlozenie(): void;
  /** Przypina wpis albo zdejmuje przypięcie. */
  naPrzypiecie(idWpisu: string, przypiety: boolean): void;
  /** Usuwa wpis; puste wskazanie kasuje historię nieprzypiętą. */
  naUsuniecie(idWpisu: string): void;
  /** Odczytuje historię z rdzenia wedle frazy i zawężenia. */
  naOdczyt(fraza: string, tylkoPrzypiete: boolean): void;
  /** Pobiera postać akapitu, w którym stoi kursor; `null`, gdy kursora nie ma. */
  naPobraniePostaci(): string | null;
  /** Nakłada pobraną postać na akapit, w którym stoi kursor. */
  naNalozeniePostaci(): void;
}

/** Panel schowka wraz z jego odświeżeniem. */
export interface SchowekHistoria {
  element: HTMLElement;
  /** Pokazuje wpisy oddane przez rdzeń. */
  pokaz(wpisy: readonly ClipboardEntry[], wszystkich: number): void;
  /** Wypisuje odmowę rdzenia w miejscu wykazu — cisza jest zakazana. */
  odmowa(powod: string): void;
  /** Zapisuje pobraną postać malarza formatów i nazywa ją w panelu. */
  ustawPostacMalarza(opis: string | null): void;
  /** Czy malarz ma pobraną postać do nałożenia. */
  malarzGotowy(): boolean;
  przestawWidocznosc(): void;
  widoczny(): boolean;
}

export function utworzSchowekHistorie(czynnosci: CzynnosciSchowka): SchowekHistoria {
  let otwarty = false;
  let postacMalarza: string | null = null;

  const fraza = document.createElement('input');
  fraza.type = 'search';
  fraza.className = 'dn-pole-kontrolka';
  fraza.placeholder = 'zawężenie historii — dopasowanie w treści wpisu';
  fraza.setAttribute('aria-label', 'Zawężenie historii schowka');

  const tylkoPrzypiete = document.createElement('input');
  tylkoPrzypiete.type = 'checkbox';
  tylkoPrzypiete.className = 'dn-przelacznik';
  tylkoPrzypiete.dataset['czynnosc'] = 'tylko-przypiete';
  tylkoPrzypiete.setAttribute('aria-label', 'Tylko wpisy przypięte');

  const etykietaPrzypietych = document.createElement('label');
  etykietaPrzypietych.className = 'ms-schowek__zawezenie';
  etykietaPrzypietych.append(tylkoPrzypiete, document.createTextNode('tylko przypięte'));

  const odswiez = document.createElement('button');
  odswiez.type = 'button';
  odswiez.className = 'dn-btn dn-btn--sm dn-btn--zarys';
  odswiez.textContent = 'Odczytaj historię';
  odswiez.dataset['czynnosc'] = 'odczyt';
  odswiez.title = 'Idzie komendą clipboard.list — historia leży w rdzeniu, nie w karcie.';

  const odloz = document.createElement('button');
  odloz.type = 'button';
  odloz.className = 'dn-btn dn-btn--sm dn-btn--atrament';
  odloz.textContent = 'Odłóż zaznaczenie';
  odloz.dataset['czynnosc'] = 'odloz';
  odloz.title =
    'Idzie komendą clipboard.push. Powtórzenie treści identycznej nie mnoży wpisów — rdzeń ' +
    'podnosi wpis zastany na czoło wykazu i mówi to polem alreadyPresent.';

  const wyczysc = document.createElement('button');
  wyczysc.type = 'button';
  wyczysc.className = 'dn-btn dn-btn--sm dn-btn--duch';
  wyczysc.textContent = 'Wyczyść historię nieprzypiętą';
  wyczysc.dataset['czynnosc'] = 'wyczysc';
  wyczysc.title =
    'clipboard.delete bez wskazania wpisu kasuje historię nieprzypiętą; wpisy przypięte zostają.';

  const pobierzPostac = document.createElement('button');
  pobierzPostac.type = 'button';
  pobierzPostac.className = 'dn-btn dn-btn--sm dn-btn--zarys';
  pobierzPostac.textContent = 'Malarz: pobierz postać';
  pobierzPostac.dataset['czynnosc'] = 'malarz-pobierz';

  const nalozPostac = document.createElement('button');
  nalozPostac.type = 'button';
  nalozPostac.className = 'dn-btn dn-btn--sm dn-btn--sygnal';
  nalozPostac.textContent = 'Malarz: nałóż postać';
  nalozPostac.dataset['czynnosc'] = 'malarz-naloz';
  nalozPostac.disabled = true;

  const zdanieMalarza = document.createElement('p');
  zdanieMalarza.className = 'dn-pole-opis ms-schowek__malarz';

  const wykaz = document.createElement('ul');
  wykaz.className = 'ms-schowek__wykaz';

  const podsumowanie = document.createElement('p');
  podsumowanie.className = 'dn-pole-opis ms-schowek__podsumowanie';

  const pas = document.createElement('div');
  pas.className = 'ms-schowek__pas';
  pas.append(fraza, etykietaPrzypietych, odswiez, odloz, wyczysc);

  const pasMalarza = document.createElement('div');
  pasMalarza.className = 'ms-schowek__pas';
  pasMalarza.append(pobierzPostac, nalozPostac);

  const tytul = document.createElement('p');
  tytul.className = 'ms-schowek__tytul';
  tytul.textContent = 'Schowek — historia i przypięcia';

  const element = document.createElement('section');
  element.className = 'ms-schowek';
  element.hidden = true;
  element.setAttribute('aria-label', 'Schowek Operatora — historia wpisów');
  element.append(tytul, pas, podsumowanie, wykaz, pasMalarza, zdanieMalarza);

  function zlecOdczyt(): void {
    czynnosci.naOdczyt(fraza.value.trim(), tylkoPrzypiete.checked);
  }

  odswiez.addEventListener('click', () => zlecOdczyt());
  fraza.addEventListener('change', () => zlecOdczyt());
  tylkoPrzypiete.addEventListener('change', () => zlecOdczyt());
  odloz.addEventListener('click', () => czynnosci.naOdlozenie());
  wyczysc.addEventListener('click', () => czynnosci.naUsuniecie(''));

  pobierzPostac.addEventListener('click', () => {
    const opis = czynnosci.naPobraniePostaci();
    ustawPostacMalarza(opis);
  });
  nalozPostac.addEventListener('click', () => czynnosci.naNalozeniePostaci());

  function ustawPostacMalarza(opis: string | null): void {
    postacMalarza = opis;
    nalozPostac.disabled = opis === null;
    zdanieMalarza.textContent =
      opis === null
        ? 'Malarz formatów nie ma pobranej postaci. Postaw kursor w akapicie wzorcowym ' +
          'i naciśnij „Malarz: pobierz postać". Postać żyje przez tę sesję okna: kontrakt ' +
          'schowka niesie treść i rodzaj wpisu, nie arkusz nastaw akapitu.'
        : `Malarz trzyma postać: ${opis}. Nałożenie przestawi akapit, w którym stoi kursor — ` +
          'treści nie ruszy.';
  }

  ustawPostacMalarza(null);

  return {
    element,

    pokaz(wpisy, wszystkich) {
      if (wpisy.length === 0) {
        podsumowanie.textContent =
          tylkoPrzypiete.checked || fraza.value.trim() !== ''
            ? 'Rdzeń nie oddał ani jednego wpisu przy tym zawężeniu — historia nie jest pusta, ' +
              'zawężenie jest ostre. Zdejmij je, żeby zobaczyć całość.'
            : 'Historia schowka jest pusta. Zaznacz fragment dokumentu i naciśnij ' +
              '„Odłóż zaznaczenie" — wpis pojedzie komendą clipboard.push do rdzenia i przetrwa ' +
              'zamknięcie karty.';
        wykaz.replaceChildren();
        return;
      }
      podsumowanie.textContent =
        `Wpisów pokazanych ${wpisy.length} z ${wszystkich} spełniających zawężenie; ` +
        `przypiętych ${wpisy.filter((wpis) => wpis.pinned).length}.`;
      wykaz.replaceChildren(...wpisy.map((wpis) => wierszWpisu(wpis, czynnosci)));
    },

    odmowa(powod) {
      wykaz.replaceChildren();
      podsumowanie.textContent = powod;
    },

    ustawPostacMalarza,
    malarzGotowy: () => postacMalarza !== null,

    przestawWidocznosc() {
      otwarty = !otwarty;
      element.hidden = !otwarty;
      if (otwarty) zlecOdczyt();
    },

    widoczny: () => otwarty,
  };
}

/** Jeden wiersz historii wraz z dwiema drogami wklejenia i przypięciem. */
function wierszWpisu(wpis: ClipboardEntry, czynnosci: CzynnosciSchowka): HTMLElement {
  const glowa = document.createElement('p');
  glowa.className = 'ms-schowek__glowa';
  glowa.textContent = schowekOpiszWpis(wpis);

  const tresc = document.createElement('p');
  tresc.className = 'ms-schowek__tresc';
  tresc.textContent = schowekPodglad(wpis);

  const tekstowy = wpis.kind === ClipboardEntryKind.Text;

  const zPostacia = document.createElement('button');
  zPostacia.type = 'button';
  zPostacia.className = 'dn-btn dn-btn--sm dn-btn--sygnal';
  zPostacia.textContent = 'Wklej z postacią';
  zPostacia.dataset['czynnosc'] = 'wklej-postac';
  zPostacia.disabled = !tekstowy;
  zPostacia.title =
    'Treść wchodzi w miejsce kursora, a akapit zachowuje postać akapitu docelowego wraz ' +
    'z nastawami wizualnymi okna.';
  zPostacia.addEventListener('click', () => czynnosci.naWklejenieZPostacia(wpis.content));

  const czyste = document.createElement('button');
  czyste.type = 'button';
  czyste.className = 'dn-btn dn-btn--sm dn-btn--zarys';
  czyste.textContent = 'Wklej jako czysty tekst';
  czyste.dataset['czynnosc'] = 'wklej-czysto';
  czyste.disabled = !tekstowy;
  czyste.title =
    'Znaczniki postaci schodzą z treści przed wstawieniem: wchodzą same litery. Oba warianty ' +
    'są równorzędne — wybór należy do Operatora.';
  czyste.addEventListener('click', () => czynnosci.naWklejenieCzyste(wpis.content));

  const przypnij = document.createElement('button');
  przypnij.type = 'button';
  przypnij.className = 'dn-btn dn-btn--sm dn-btn--duch';
  przypnij.textContent = wpis.pinned ? 'Odepnij' : 'Przypnij';
  przypnij.dataset['czynnosc'] = 'przypnij';
  przypnij.title =
    'Wpis przypięty nie wygasa wraz z zasadą retencji historii — komenda clipboard.pin.';
  przypnij.addEventListener('click', () => czynnosci.naPrzypiecie(wpis.id, !wpis.pinned));

  const usun = document.createElement('button');
  usun.type = 'button';
  usun.className = 'dn-btn dn-btn--sm dn-btn--duch';
  usun.textContent = 'Usuń wpis';
  usun.dataset['czynnosc'] = 'usun';
  usun.addEventListener('click', () => czynnosci.naUsuniecie(wpis.id));

  const pasWpisu = document.createElement('div');
  pasWpisu.className = 'ms-schowek__pas';
  pasWpisu.append(zPostacia, czyste, przypnij, usun);

  const pozycja = document.createElement('li');
  pozycja.dataset['wpis'] = wpis.id;
  pozycja.dataset['przypiety'] = wpis.pinned ? 'tak' : 'nie';
  pozycja.dataset['rodzaj'] = wpis.kind;
  pozycja.append(glowa, tresc, pasWpisu);
  return pozycja;
}
