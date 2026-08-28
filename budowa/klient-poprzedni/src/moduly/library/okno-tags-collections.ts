import { utworzDymekObjasnienia } from '../../komponenty/dymek';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import {
  pobierzPlik,
  poleTekstowe,
  przycisk,
  utworzWierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import { opisOdmowy } from '../../komponenty/odmowa';
import { mapaKolekcji } from './eksporty-biblioteki';
import { BEZ_KOMENDY_ETYKIETY } from './etykiety-biblioteki';
import { utworzFormularzKolekcji } from './formularz-kolekcji';
import { utworzPanelAkcji } from './panel-akcji';
import type { StanBiblioteki } from './stan-biblioteki';
import { utworzStanOkna } from './stan-okna';
import { nadajEtykieteZbiorczo } from './zapisy-zbiorcze';
import type { ZrodloOtoczenia } from './zrodlo-otoczenia';

/**
 * Tags & Collections to okno zarządca modułu biblioteki: nadaje etykiety, tworzy kolekcje
 * i przypisuje zasoby, korzystając ze słownika etykiet i wykazu kolekcji pobranych z rdzenia.
 */
export interface OknoEtykiet {
  element: HTMLElement;
  odswiez(): void;
}

export function utworzOknoEtykiet(stan: StanBiblioteki, otoczenie: ZrodloOtoczenia): OknoEtykiet {
  const rama = utworzRameOkna({
    kod: 'tags-collections',
    tytul: 'Tags & Collections',
    rola: 'zarządca',
    modul: 'Library',
    dodatkiNaglowka: [
      utworzDymekObjasnienia(
        'Etykiety i kolekcje definiowane swobodnie przez Operatora. Etykieta idzie do plików ' +
          'zaznaczonych w Library Explorer; przypisanie zasobu nie ma tu własnego wejścia.',
        { powloka: 'ml-dymek', znak: 'ml-dymek__znak' },
      ),
    ],
  });
  const okno = utworzStanOkna();
  const odpowiedz = utworzWierszOdpowiedzi();
  const kolekcje = utworzFormularzKolekcji(stan);
  const panel = utworzPanelAkcji(otoczenie, stan, BEZ_KOMENDY_ETYKIETY);

  const etykieta = poleTekstowe({
    etykieta: 'Nowa etykieta',
    podpowiedz: 'np. umowy',
    opis: 'Etykieta dochodzi do etykiet już nadanych plikom zaznaczonym w Explorerze.',
  });

  // Wskaźnik odczytu stoi obok pola, nigdy zamiast kontrolki: pole pozostaje
  // edytowalne w czasie zapisu.
  const wskaznik = document.createElement('span');
  wskaznik.className = 'dn-spinner ml-etykiety__wskaznik';
  wskaznik.setAttribute('role', 'status');
  wskaznik.setAttribute('aria-label', 'Trwa zapis etykiety');
  wskaznik.hidden = true;

  const nadaj = przycisk('Nadaj etykietę', 'dn-btn dn-btn--sm dn-btn--atrament');
  nadaj.dataset['czynnosc'] = 'etykieta';
  nadaj.addEventListener('click', () => void nadajEtykiete());

  async function nadajEtykiete(): Promise<void> {
    const zaznaczone = stan.zaznaczone();
    const pliki = stan.pliki().filter((plik) => zaznaczone.includes(plik.id));
    wskaznik.hidden = false;
    odpowiedz.pokaz('Nadawanie etykiety…', true);
    const wynik = await nadajEtykieteZbiorczo(stan, etykieta.kontrolka.value, pliki);
    wskaznik.hidden = true;
    odpowiedz.pokaz(wynik.tresc, wynik.powodzenie);
  }

  /** Liczba plików wykazu noszących wskazaną etykietę — statystyka jej użycia. */
  function uzycieEtykiety(kod: string): number {
    return stan.pliki().filter((plik) => (plik.tags ?? []).includes(kod)).length;
  }

  const wywiezMape = przycisk('Mapa kolekcji (Markdown)', 'dn-btn dn-btn--sm dn-btn--duch');
  wywiezMape.dataset['czynnosc'] = 'mapa-kolekcji';
  wywiezMape.addEventListener('click', () => {
    const kolekcje = stan.kolekcje();
    if (stan.pliki().length === 0 && kolekcje.length === 0) {
      odpowiedz.pokaz('Wykaz plików i wykaz kolekcji są puste — mapa nie miałaby treści.', false);
      return;
    }
    pobierzPlik('mapa-kolekcji.md', mapaKolekcji(stan.pliki(), kolekcje), 'text/markdown');
    odpowiedz.pokaz(
      `Mapa złożona z wykazu pokazanego w oknie: kolekcji ${kolekcje.length}, ` +
        `plików ${stan.pliki().length}. Kolekcje wypisane kodami — nazwy własnej ` +
        'kontrakt nie oddaje poza chwilą założenia.',
      true,
    );
  });

  const wywiezTezaurus = przycisk('Tezaurus SKOS (Turtle)', 'dn-btn dn-btn--sm dn-btn--duch');
  wywiezTezaurus.dataset['czynnosc'] = 'tezaurus-skos';
  wywiezTezaurus.addEventListener('click', () => void wywiezTezaurusZRdzenia());

  async function wywiezTezaurusZRdzenia(): Promise<void> {
    odpowiedz.pokaz('Rdzeń składa tezaurus wraz z relacjami…', true);
    const wynik = await stan.zrodlo.wywiezTezaurus('turtle');
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowy('Wywóz tezaurusa', wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    pobierzPlik('tezaurus.ttl', wynik.wynik.content, 'text/turtle');
    odpowiedz.pokaz(
      `Tezaurus z rdzenia: pojęć ${wynik.wynik.conceptCount}, relacji ` +
        `${wynik.wynik.relationCount}, zapis ${wynik.wynik.format}.`,
      true,
    );
  }

  // ── Zarząd słownikiem etykiet ─────────────────────────────────────────────
  const etykietaZrodlowa = poleTekstowe({
    etykieta: 'Etykieta słownika',
    podpowiedz: 'nazwa etykiety już istniejącej',
    opis: 'Wejście czynności słownikowych: zmiany nazwy, łączenia i usunięcia.',
  });
  const etykietaDocelowa = poleTekstowe({
    etykieta: 'Nowa nazwa albo etykieta docelowa',
    podpowiedz: 'np. klient-iks',
    opis: 'Przy zmianie nazwy — nazwa po zmianie; przy łączeniu — etykieta, która zostaje.',
  });

  /** Wykaz słownika z rdzenia wraz z licznikami użycia po całym repozytorium. */
  const slownik = document.createElement('ul');
  slownik.className = 'ml-etykiety__slownik';

  async function odczytajSlownik(): Promise<void> {
    const wynik = await stan.zrodlo.slownikEtykiet();
    if (!wynik.udany || wynik.wynik === undefined) {
      slownik.replaceChildren();
      return;
    }
    slownik.replaceChildren(
      ...wynik.wynik.tags.map((wpis) =>
        pozycja(
          wpis.color === undefined ? wpis.name : `${wpis.name} (${wpis.color})`,
          'slownik',
          wpis.fileCount,
        ),
      ),
    );
  }

  const zmienNazwe = przycisk('Zmień nazwę etykiety', 'dn-btn dn-btn--sm dn-btn--zarys');
  zmienNazwe.dataset['czynnosc'] = 'etykieta-zmiana';
  zmienNazwe.addEventListener('click', () => void zmienNazweEtykiety());

  async function zmienNazweEtykiety(): Promise<void> {
    odpowiedz.pokaz('Zmiana nazwy etykiety w całym repozytorium…', true);
    const wynik = await stan.zrodlo.zmienEtykiete(
      etykietaZrodlowa.kontrolka.value.trim(),
      etykietaDocelowa.kontrolka.value.trim() === ''
        ? undefined
        : etykietaDocelowa.kontrolka.value.trim(),
    );
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowy('Zmiana etykiety', wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    odpowiedz.pokaz(
      `Etykieta po zmianie: ${wynik.wynik.tag.name}. Zasobów dotkniętych: ` +
        `${wynik.wynik.affectedFiles}.`,
      true,
    );
    await odczytajSlownik();
    await stan.odczytaj('', '');
  }

  const polacz = przycisk('Połącz etykiety', 'dn-btn dn-btn--sm dn-btn--zarys');
  polacz.dataset['czynnosc'] = 'etykieta-laczenie';
  polacz.addEventListener('click', () => void polaczEtykiety());

  async function polaczEtykiety(): Promise<void> {
    odpowiedz.pokaz('Łączenie etykiet duplikujących się…', true);
    const wynik = await stan.zrodlo.polaczEtykiety(
      [etykietaZrodlowa.kontrolka.value.trim()],
      etykietaDocelowa.kontrolka.value.trim(),
    );
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowy('Łączenie etykiet', wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    odpowiedz.pokaz(
      `Wchłonięto: ${wynik.wynik.mergedNames.join(', ') || 'nic'} → ${wynik.wynik.tag.name}. ` +
        `Zasobów dotkniętych: ${wynik.wynik.affectedFiles}.`,
      true,
    );
    await odczytajSlownik();
    await stan.odczytaj('', '');
  }

  /** Etykieta wskazana do usunięcia, potwierdzana drugim naciśnięciem. */
  let potwierdzanaEtykieta = '';

  const usun = przycisk('Usuń etykietę', 'dn-btn dn-btn--sm dn-btn--zarys');
  usun.dataset['czynnosc'] = 'etykieta-usuniecie';
  usun.addEventListener('click', () => void usunEtykiete());

  async function usunEtykiete(): Promise<void> {
    const nazwa = etykietaZrodlowa.kontrolka.value.trim();
    const potwierdzone = potwierdzanaEtykieta === nazwa && nazwa !== '';
    const wynik = await stan.zrodlo.usunEtykiete(nazwa, potwierdzone);
    if (!wynik.udany || wynik.wynik === undefined) {
      potwierdzanaEtykieta = nazwa;
      odpowiedz.pokaz(opisOdmowy('Usunięcie etykiety', wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    potwierdzanaEtykieta = '';
    odpowiedz.pokaz(
      `Etykieta ${nazwa} usunięta ze słownika i z ${wynik.wynik.affectedFiles} zasobów.`,
      true,
    );
    await odczytajSlownik();
    await stan.odczytaj('', '');
  }

  const relacjaNadrzedna = przycisk('Ustaw relację nadrzędną', 'dn-btn dn-btn--sm dn-btn--zarys');
  relacjaNadrzedna.dataset['czynnosc'] = 'tezaurus-relacja';
  relacjaNadrzedna.addEventListener('click', () => void ustawRelacje());

  async function ustawRelacje(): Promise<void> {
    odpowiedz.pokaz('Ustanawianie relacji tezaurusa…', true);
    const wynik = await stan.zrodlo.ustawRelacje(
      etykietaZrodlowa.kontrolka.value.trim(),
      etykietaDocelowa.kontrolka.value.trim(),
      'broader',
      false,
    );
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowy('Relacja tezaurusa', wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    odpowiedz.pokaz(
      `Relacja ${etykietaZrodlowa.kontrolka.value} → ${etykietaDocelowa.kontrolka.value} ` +
        `(nadrzędna) ${wynik.wynik.relationSet ? 'stoi' : 'zdjęta'}.`,
      true,
    );
  }

  const paskaSlownika = document.createElement('div');
  paskaSlownika.className = 'ml-etykiety__pasek';
  paskaSlownika.append(zmienNazwe, polacz, usun, relacjaNadrzedna);

  const zbiorEtykiet = document.createElement('ul');
  zbiorEtykiet.className = 'ml-etykiety__zbior';

  const drzewoKolekcji = document.createElement('ul');
  drzewoKolekcji.className = 'ml-kolekcje__drzewo';

  const pasek = document.createElement('div');
  pasek.className = 'ml-etykiety__pasek';
  pasek.append(nadaj, wskaznik);

  okno.tresc.append(
    etykieta.element,
    pasek,
    zbiorEtykiet,
    drzewoKolekcji,
    kolekcje.element,
    etykietaZrodlowa.element,
    etykietaDocelowa.element,
    paskaSlownika,
    slownik,
  );
  okno.pierwszaAkcja('Nadaj pierwszą etykietę', () => etykieta.kontrolka.focus());
  rama.pasek.append(wywiezMape, wywiezTezaurus);
  rama.cialo.append(okno.element, odpowiedz.element, panel.element);

  /** Stany obowiązkowe okna czytane z fazy wykazu plików — własnego odczytu okno mieć nie może. */
  function ustawStan(): void {
    if (stan.faza() === 'odczyt') {
      okno.ladowanie(
        'Rdzeń odczytuje wykaz plików — etykiety i kolekcje składają się z ich pól.',
      );
      return;
    }
    if (stan.faza() === 'blad') {
      okno.blad(stan.powod());
      return;
    }
    if (stan.faza() === 'spoczynek') {
      okno.puste(
        'Wykaz plików nieodczytany',
        'Etykiety i kolekcje biorą się ze zbioru plików, a ten nie przyszedł jeszcze ' +
          'z rdzenia. Odśwież wykaz w Library Explorer.',
      );
      return;
    }
    if (stan.etykiety().length === 0 && stan.kolekcje().length === 0) {
      okno.puste(
        'Bez etykiet i kolekcji',
        'Rdzeń nie zwrócił ani jednej etykiety i ani jednej kolekcji. Zaznacz pliki ' +
          'w Library Explorer i nadaj im etykietę wpisaną wyżej albo załóż kolekcję ' +
          'formularzem niżej.',
      );
      return;
    }
    okno.gotowe();
  }

  return {
    element: rama.element,

    odswiez() {
      kolekcje.odswiez();
      void odczytajSlownik();
      zbiorEtykiet.replaceChildren(
        ...stan.etykiety().map((kod) => pozycja(kod, 'etykieta', uzycieEtykiety(kod))),
      );
      drzewoKolekcji.replaceChildren(
        ...stan.kolekcje().map((kod) => pozycja(kod, 'kolekcja', uzycieKolekcji(kod))),
      );
      ustawStan();
    },
  };

  /** Liczba plików wykazu należących do wskazanej kolekcji. */
  function uzycieKolekcji(kod: string): number {
    return stan.pliki().filter((plik) => (plik.collectionIds ?? []).includes(kod)).length;
  }
}

/** Jedna pozycja zbioru etykiet albo drzewa kolekcji wraz z licznikiem użycia — liczba plików obecnych w wykazie, nie w całym repozytorium. */
function pozycja(kod: string, rodzaj: string, wWykazie: number): HTMLElement {
  const plakietka = document.createElement('span');
  plakietka.className = 'dn-plakietka';
  plakietka.textContent = kod;

  const licznik = document.createElement('span');
  licznik.className = 'ml-etykiety__licznik';
  licznik.textContent = String(wWykazie);
  licznik.title = `Plików z tą pozycją w bieżącym wykazie: ${wWykazie}`;

  const element = document.createElement('li');
  element.className = 'ml-etykiety__pozycja';
  element.dataset['rodzaj'] = rodzaj;
  element.append(plakietka, licznik);
  return element;
}
