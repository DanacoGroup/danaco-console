import { KnowledgeScope, type KnowledgeHit } from '../../../../shared/contract';
import { utworzMenuDrzewo, type PozycjaMenu } from '../../komponenty/menu-drzewo';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import { utworzStanTresci } from '../../komponenty/stan-tresci';
import {
  poleLiczbowe,
  przyciskAkcji,
  wiersz,
} from '../../modele/kontrolki-formularza-braki';
import { utworzZrodloWiedzy } from './zrodlo-wiedzy';
import type { Kanal } from '../../protokol/kanal';

/** Okno wyszukiwania po znaczeniu zadaje pytanie w języku operatora i odpowiada fragmentami jego własnych treści wraz ze źródłem i trafnością. */
export interface OknoWyszukiwaniaZnaczenia {
  element: HTMLElement;
  /** Powtarza ostatnie pytanie; bez pytania nie woła rdzenia. */
  odswiez(): void;
  /** Zwinięcie odpina nasłuch dokumentu, żeby zamknięte okno nie wskazywało na usunięty element. */
  zamknij(): void;
}

/** Wiersz ze sterem zamiast pola tworzy podpis nad menu drzewem osobno od wiersza biblioteki, bo tamten wiąże etykietę z uchwytem menu. */
function wierszZakresu(etykieta: string, ster: HTMLElement): HTMLElement {
  const element = document.createElement('div');
  element.className = 'dwz-wiersz';

  const podpis = document.createElement('span');
  podpis.className = 'dn-pole-etykieta';
  podpis.textContent = etykieta;

  element.append(podpis, ster);
  return element;
}

/** Zakresy wiedzy pochodzą z wykazu kontraktu, nie z literałów okna, i niosą zdanie mówiące, co każdy z nich przeszukuje. */
const ZAKRESY: ReadonlyArray<readonly [KnowledgeScope, string, string]> = [
  [
    KnowledgeScope.All,
    'wszystko',
    'Biblioteka, historia rozmów i pliki przestrzeni naraz — najszersze pytanie, jakie da się zadać.',
  ],
  [KnowledgeScope.Library, 'biblioteka', 'Wyłącznie treści wniesione do biblioteki Operatora.'],
  [
    KnowledgeScope.History,
    'historia rozmów',
    'Wyłącznie to, co padło w rozmowach — ustalenia, o których pamięta zapis, a nie plik.',
  ],
  [
    KnowledgeScope.Workspace,
    'pliki przestrzeni',
    'Wyłącznie pliki przestrzeni roboczej tego okna; bez wskazania okna zakres nie ma czego przeszukać.',
  ],
];

/** Zakres, od którego okno startuje, jest najszerszy możliwy, bo pierwsze pytanie pada bez żadnych zawężeń. */
const ZAKRES_POCZATKOWY = KnowledgeScope.All;

export function utworzOknoWyszukiwaniaZnaczenia(
  kanal: Kanal,
  idOkna: string,
): OknoWyszukiwaniaZnaczenia {
  const zrodlo = utworzZrodloWiedzy(kanal);
  const rama = utworzRameOkna({
    tytul: 'Wyszukiwanie po znaczeniu',
    rola: 'wiodące',
    przeznaczenie:
      'Pytanie w języku Operatora; odpowiedzią są fragmenty jego własnych treści wraz ze źródłem i trafnością.',
    modul: 'Poczta i wiedza',
    przedrostek: 'dwz',
    ogniskowalne: true,
  });
  const tresc = utworzStanTresci('dwz');

  const pytanie = document.createElement('input');
  pytanie.type = 'search';
  pytanie.className = 'dn-pole-kontrolka';
  pytanie.placeholder = 'O co pytasz? np. „ustalenia z rozmowy o dostawcy"';
  pytanie.setAttribute('aria-label', 'Pytanie do wiedzy Operatora');

  // Zakres jest zmienną, nie odczytem z kontrolki, bo menu drzewa oddaje klucz przy wyborze.
  let zakres: KnowledgeScope = ZAKRES_POCZATKOWY;
  const menuZakresu = utworzMenuDrzewo({
    nastawa: 'Zakres przeszukiwanej wiedzy',
    naWybor: (klucz) => {
      zakres = klucz as KnowledgeScope;
      przepiszDrzewoZakresu();
      // Zmiana zakresu po zadanym pytaniu powtarza je sama, bo zawężenie bez nowego wyszukania nic nie da.
      if (ostatniePytanie !== '') szukajTeraz();
    },
  });

  function przepiszDrzewoZakresu(): void {
    const drzewo: PozycjaMenu[] = ZAKRESY.map(([wartosc, nazwa, opis]) => ({
      rodzaj: 'wybor',
      klucz: wartosc,
      nazwa,
      opis,
      wybrany: wartosc === zakres,
    }));
    menuZakresu.ustaw(OPIS_ZAKRESU[zakres] ?? zakres, drzewo);
  }

  const granica = poleLiczbowe('Górna granica liczby fragmentów', '20');

  const szukaj = przyciskAkcji('Szukaj po znaczeniu', 'dn-btn dn-btn--sygnal');
  const przebuduj = przyciskAkcji('Przebuduj wskaźnik', 'dn-btn dn-btn--zarys');

  rama.narzedzia.append(
    wiersz('Pytanie', pytanie, { klasa: 'dwz-wiersz dwz-wiersz--szeroki' }),
    wierszZakresu('Zakres', menuZakresu.element),
    wiersz('Najwyżej fragmentów', granica, {
      klasa: 'dwz-wiersz',
      objasnienie: 'Puste znaczy: granica rdzenia.',
    }),
  );
  rama.akcje.append(szukaj, przebuduj);
  rama.cialo.append(tresc.element);

  // Ostatnie pytanie zadane, nie ostatnie wpisane — odświeżenie powtarza to, co poszło do rdzenia.
  let ostatniePytanie = '';

  function szukajTeraz(): void {
    const fraza = pytanie.value.trim();
    if (fraza === '') {
      // Pytanie puste nie jedzie do rdzenia, bo pole jest w kontrakcie wymagane, a odmowa byłaby zbędna.
      tresc.pusto('Wpisz pytanie — wyszukiwanie po znaczeniu potrzebuje zdania, nie pustego pola.');
      return;
    }
    ostatniePytanie = fraza;
    tresc.ladowanie('Szukam po znaczeniu…');
    rama.ustawZnacznik('');
    void zrodlo
      .szukaj({
        query: fraza,
        scope: zakres,
        ...(idOkna === '' ? {} : { windowId: idOkna }),
        ...(granicaLiczbowa(granica) === undefined ? {} : { limit: granicaLiczbowa(granica) }),
      })
      .then((wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          tresc.blad(opisOdmowyBledu('Szukanie po znaczeniu nie udało się', wynik.blad), wynik.blad);
          rama.ustawZnacznik('odmowa', 'blad');
          return;
        }
        rysujTrafienia(wynik.wynik.wyniki, wynik.wynik.wszystkich, fraza);
      });
  }

  function rysujTrafienia(
    trafienia: readonly KnowledgeHit[],
    wszystkich: number,
    fraza: string,
  ): void {
    if (trafienia.length === 0) {
      // Pustka bywa poprawna z dwóch różnych przyczyn, których sam pusty wykaz nie odróżnia.
      tresc.pusto(
        `Wiedza Operatora nie ma fragmentu bliskiego pytaniu „${fraza}". Jeżeli treści przybyło od ostatniego razu, przebuduj wskaźnik — bez tego nowe pozycje są dla szukania niewidoczne.`,
      );
      rama.ustawZnacznik('0 fragmentów', 'ostrzezenie');
      return;
    }
    rama.ustawZnacznik(`${trafienia.length} z ${wszystkich}`, 'sukces');
    const miejsce = tresc.tresc();
    const wykaz = document.createElement('ol');
    wykaz.className = 'dwz-trafienia';
    wykaz.setAttribute('aria-label', `Fragmenty odnalezione dla pytania ${fraza}`);
    for (const trafienie of trafienia) wykaz.append(pozycjaTrafienia(trafienie, tresc.potwierdzenie));
    miejsce.append(wykaz);
    if (trafienia.length < wszystkich) {
      // Wykaz przycięty granicą mówi to wprost, żeby brak reszty nie wyglądał na brak wyniku.
      const przypis = document.createElement('p');
      przypis.className = 'dwz-przypis';
      przypis.textContent = `Pokazano ${trafienia.length} z ${wszystkich} fragmentów — podnieś granicę, żeby zobaczyć resztę.`;
      miejsce.append(przypis);
    }
  }

  function przebudujTeraz(): void {
    tresc.potwierdzenie('Przebudowa wskaźnika w toku…', true);
    przebuduj.disabled = true;
    void zrodlo
      .przebudujWskaznik({
        scope: zakres,
        ...(idOkna === '' ? {} : { windowId: idOkna }),
        rebuild: true,
      })
      .then((wynik) => {
        przebuduj.disabled = false;
        if (!wynik.udany || wynik.wynik === undefined) {
          tresc.potwierdzenie(
            opisOdmowyBledu('Przebudowa wskaźnika nie udała się', wynik.blad),
            false,
          );
          return;
        }
        const model = wynik.wynik.model === undefined ? '' : ` (model osadzeń: ${wynik.wynik.model})`;
        tresc.potwierdzenie(
          `Wskaźnik przebudowany: wniesiono ${wynik.wynik.wniesione}, we wskaźniku ${wynik.wynik.wszystkich} pozycji${model}.`,
          true,
        );
        // Po przebudowie wskaźnika ostatnie pytanie jest powtarzane, bo po to wskaźnik przebudowano.
        if (ostatniePytanie !== '') szukajTeraz();
      });
  }

  szukaj.addEventListener('click', szukajTeraz);
  przebuduj.addEventListener('click', przebudujTeraz);
  pytanie.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key === 'Enter') szukajTeraz();
  });

  tresc.pusto('Zadaj pytanie — okno nie pytało jeszcze rdzenia.');

  // Uchwyt niesie wartość od pierwszej chwili, żeby było widać, czego okno w tej chwili przeszukuje.
  przepiszDrzewoZakresu();

  return {
    element: rama.element,
    odswiez() {
      // Bez zadanego pytania odświeżenie nie woła rdzenia, bo pole zapytania jest w kontrakcie wymagane.
      if (ostatniePytanie === '') return;
      pytanie.value = ostatniePytanie;
      szukajTeraz();
    },
    zamknij() {
      menuZakresu.zwin();
    },
  };
}

/** Funkcja czyta granicę z pola liczbowego; wartość pusta albo niedodatnia znaczy granicę ustaloną przez rdzeń. */
function granicaLiczbowa(pole: HTMLInputElement): number | undefined {
  const liczba = Number.parseInt(pole.value, 10);
  return Number.isFinite(liczba) && liczba > 0 ? liczba : undefined;
}

/**
 * Pozycja wykazu trafień: źródło, trafność, treść fragmentu, uchwyt do źródła.
 *
 * Kolejność jest treścią: źródło stoi pierwsze, przed fragmentem, bo Operator
 * czyta wykaz wzrokiem po lewej krawędzi i ma wiedzieć, czyje to zdanie, zanim
 * je przeczyta.
 */
function pozycjaTrafienia(
  trafienie: KnowledgeHit,
  potwierdz: (zdanie: string, udane: boolean) => void,
): HTMLElement {
  const element = document.createElement('li');
  element.className = 'dwz-trafienie';

  const naglowek = document.createElement('div');
  naglowek.className = 'dwz-trafienie__naglowek';
  naglowek.append(uchwytZrodla(trafienie, potwierdz), plakietkaZakresu(trafienie), plakietkaTrafnosci(trafienie));

  const fragment = document.createElement('p');
  fragment.className = 'dwz-trafienie__tresc';
  fragment.textContent = trafienie.text;

  element.append(naglowek, fragment);
  return element;
}

/** Uchwyt do źródła jest przyciskiem kopiującym wskazanie źródła, gdy rdzeń je podał, bo kontrakt nie ma komendy skoku do źródła. */
function uchwytZrodla(
  trafienie: KnowledgeHit,
  potwierdz: (zdanie: string, udane: boolean) => void,
): HTMLElement {
  const wskazanie = (trafienie.sourceId ?? '').trim();
  if (wskazanie === '') {
    const bezWskazania = document.createElement('span');
    bezWskazania.className = 'dwz-trafienie__zrodlo dwz-trafienie__zrodlo--bez-wskazania';
    bezWskazania.textContent = trafienie.source;
    bezWskazania.title =
      'Rdzeń podał nazwę źródła, ale nie podał wskazania (sourceId) — nie ma czym po nie sięgnąć.';
    return bezWskazania;
  }
  const uchwyt = przyciskAkcji(trafienie.source, 'dn-btn dn-btn--duch dwz-trafienie__zrodlo');
  uchwyt.title = `Kopiuje wskazanie źródła: ${wskazanie}`;
  uchwyt.addEventListener('click', () => {
    void navigator.clipboard
      .writeText(wskazanie)
      .then(() => potwierdz(`Wskazanie źródła skopiowane: ${wskazanie}`, true))
      .catch(() =>
        // Schowek bywa odmówiony przez przeglądarkę, więc wskazanie trafia wprost do treści potwierdzenia.
        potwierdz(`Schowek niedostępny. Wskazanie źródła: ${wskazanie}`, false),
      );
  });
  return uchwyt;
}

/** Funkcja buduje plakietkę zakresu, z którego przyszedł fragment: biblioteka, historia albo przestrzeń. */
function plakietkaZakresu(trafienie: KnowledgeHit): HTMLElement {
  const element = document.createElement('span');
  element.className = 'dn-plakietka';
  element.textContent = OPIS_ZAKRESU[trafienie.scope] ?? trafienie.scope;
  return element;
}

/**
 * Trafność w setnych. Bez niej wykaz posortowany malejąco wygląda tak samo przy
 * trafieniu bliskim i przy przypadkowym — a rozstrzyga to Operator, nie wykaz.
 * Rdzeń ma prawo trafności nie podać; wtedy plakietka mówi o braku miary,
 * a nie zmyśla zera.
 */
function plakietkaTrafnosci(trafienie: KnowledgeHit): HTMLElement {
  const element = document.createElement('span');
  if (trafienie.score === undefined) {
    element.className = 'dn-plakietka';
    element.textContent = 'trafność nieznana';
    element.title = 'Rdzeń nie podał miary trafności tego fragmentu.';
    return element;
  }
  // Klasy są dokładane jawnymi literałami, bo kontrola pokrycia klas CSS czyta napisy, nie składa je.
  element.className = 'dn-plakietka';
  element.classList.add(
    trafienie.score >= PROG_TRAFNOSCI_MOCNEJ ? 'dn-plakietka--sukces' : 'dn-plakietka--ostrzezenie',
  );
  element.textContent = `trafność ${trafienie.score}/100`;
  return element;
}

/** Próg trafności mocnej jest umową okna, nie miarą rdzenia, i stoi stałą nazwaną, żeby liczba nie udawała wartości znikąd. */
const PROG_TRAFNOSCI_MOCNEJ = 70;

/** Nazwy zakresów podane po polsku odpowiadają kluczom zakresu z kontraktu, nie literałom wymyślonym w oknie. */
const OPIS_ZAKRESU: Record<string, string> = {
  [KnowledgeScope.Library]: 'biblioteka',
  [KnowledgeScope.History]: 'historia rozmów',
  [KnowledgeScope.Workspace]: 'pliki przestrzeni',
  [KnowledgeScope.All]: 'wszystko',
};
