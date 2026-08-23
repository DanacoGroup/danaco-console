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

/**
 * Okno wyszukiwania po znaczeniu — pytanie w języku Operatora, odpowiedź
 * fragmentami jego własnych treści. Wyszukiwanie po słowach ma osobną drogę
 * (`library.file.search`); tutaj jedzie `knowledge.search` i `knowledge.index`.
 *
 * Źródło jest treścią wyniku, nie ozdobą przy nim, bo trafienie bez wskazania
 * źródła jest bezwartościowe. Dlatego:
 *  • każdy fragment niesie `source` w nagłówku pozycji, nie w dymku;
 *  • źródło jest klikalne — przyciskiem, który kładzie `sourceId` do schowka
 *    Operatora, bo to jedyne wskazanie, którym da się po źródło sięgnąć
 *    w innym oknie. Kontrakt nie niesie komendy „otwórz źródło fragmentu"
 *    (`KnowledgeHit` ma `source` i `sourceId`, i nic więcej), więc skok wprost
 *    do pliku biblioteki nie ma pokrycia i okno mówi to wprost;
 *  • fragment bez `sourceId` dostaje zdanie o braku wskazania zamiast martwego
 *    przycisku: rdzeń ma prawo źródła nie umieć wskazać.
 *
 * Trafność też jest widoczna: `score` jest w setnych (100 = najbliższy możliwy)
 * i stoi plakietką przy każdej pozycji. Bez niej wykaz posortowany malejąco
 * wygląda tak samo przy trafieniu bliskim i przy trafieniu przypadkowym.
 *
 * Przebudowa wskaźnika stoi w tym samym oknie, bo pusty wynik szukania ma
 * dokładnie dwie przyczyny: nie ma czego znaleźć albo wskaźnik nie został
 * zbudowany. Przycisk przebudowy obok pustego wyniku pozwala rozstrzygnąć to
 * jednym kliknięciem.
 *
 * Wygląd w całości z żetonów i biblioteki (`komponenty/`, `--dn-*`). Plik nie
 * zna ani jednej barwy.
 */
export interface OknoWyszukiwaniaZnaczenia {
  element: HTMLElement;
  /** Powtarza ostatnie pytanie; bez pytania nie woła rdzenia. */
  odswiez(): void;
  /**
   * Zwija ster zakresu wraz z jego nasłuchem dokumentu.
   *
   * `menu-drzewo` zakłada nasłuch `pointerdown` na dokumencie, dopóki jest
   * otwarte — okno zdjęte ze sceny przy rozwiniętym menu zostawiłoby ten
   * nasłuch wskazujący na element, którego w dokumencie już nie ma.
   */
  zamknij(): void;
}

/**
 * Wiersz ze sterem zamiast pola — podpis nad menu-drzewem.
 *
 * Osobno od `wiersz()` z biblioteki kontrolek, bo tamten buduje `<label>`,
 * a etykieta bez `for` wiąże się z pierwszym potomkiem dającym się etykietować
 * — czyli z uchwytem menu. Kliknięcie w podpis otwierałoby wtedy menu (podpis
 * nie jest sterem), a nazwa dostępna uchwytu wchodziłaby w spór z `aria-label`,
 * które mechanizm ustawia sam i które niesie bieżącą wartość nastawy.
 */
function wierszZakresu(etykieta: string, ster: HTMLElement): HTMLElement {
  const element = document.createElement('div');
  element.className = 'dwz-wiersz';

  const podpis = document.createElement('span');
  podpis.className = 'dn-pole-etykieta';
  podpis.textContent = etykieta;

  element.append(podpis, ster);
  return element;
}

/**
 * Zakresy wiedzy podane wykazem kontraktu, nie literałami okna, wraz ze
 * zdaniem mówiącym, co każdy z nich przeszukuje.
 *
 * Zdanie jest tu potrzebne bardziej niż przy większości nastaw: „wszystko"
 * i „pliki przestrzeni" brzmią podobnie, a przeszukują dwa różne zbiory —
 * i to od nich zależy, czy pusty wynik znaczy „nie ma", czy „nie tam".
 */
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

/** Zakres, od którego okno startuje — najszerszy, bo pytanie pada bez zawężeń. */
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

  /**
   * Zakres wskazany w oknie. Zmienna, a nie odczyt z kontrolki: menu-drzewo
   * jest sterem oddającym klucz przy wyborze — mechanizm z biblioteki nie
   * udaje `<input>` i nie ma być o wartość pytany.
   */
  let zakres: KnowledgeScope = ZAKRES_POCZATKOWY;
  const menuZakresu = utworzMenuDrzewo({
    nastawa: 'Zakres przeszukiwanej wiedzy',
    naWybor: (klucz) => {
      zakres = klucz as KnowledgeScope;
      przepiszDrzewoZakresu();
      // Zmiana zakresu po zadanym pytaniu powtarza je sama: Operator zawęża
      // zakres właśnie dlatego, że poprzednia odpowiedź go nie zadowoliła,
      // a drugie kliknięcie w „Szukaj" byłoby ruchem bez treści.
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

  /**
   * Ostatnie pytanie zadane, nie ostatnie wpisane. `odswiez()` powtarza to,
   * co naprawdę poszło do rdzenia — inaczej odświeżenie po zmianie treści pola
   * pokazałoby wynik jednego pytania pod nagłówkiem innego.
   */
  let ostatniePytanie = '';

  function szukajTeraz(): void {
    const fraza = pytanie.value.trim();
    if (fraza === '') {
      // Pytanie puste nie jedzie do rdzenia: `query` jest w kontrakcie polem
      // wymaganym, a odmowa walidacji powiedziałaby Operatorowi to samo, tylko
      // po podróży w obie strony i cudzym językiem.
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
      // Pustka bywa poprawna — i ma dwie przyczyny, których Operator z samego
      // pustego wykazu nie odróżni. Zdanie nazywa obie i wskazuje przycisk.
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
      // Wykaz przycięty granicą mówi to wprost: Operator, który widzi
      // dziesięć pozycji z czterdziestu i nie wie o tym, uzna, że reszty nie ma.
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
        // Po przebudowie powtarzamy ostatnie pytanie sami: Operator przebudował
        // wskaźnik właśnie dlatego, że poprzedni wynik go nie zadowolił.
        if (ostatniePytanie !== '') szukajTeraz();
      });
  }

  szukaj.addEventListener('click', szukajTeraz);
  przebuduj.addEventListener('click', przebudujTeraz);
  pytanie.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key === 'Enter') szukajTeraz();
  });

  tresc.pusto('Zadaj pytanie — okno nie pytało jeszcze rdzenia.');

  // Uchwyt niesie wartość od pierwszej chwili: pusty prostokąt w miejscu
  // nastawy nie mówi Operatorowi, czego okno w tej chwili przeszukuje.
  przepiszDrzewoZakresu();

  return {
    element: rama.element,
    odswiez() {
      // Bez zadanego pytania odświeżenie nie woła rdzenia: `knowledge.search`
      // wymaga `query`, więc wywołanie byłoby pewną odmową walidacji.
      if (ostatniePytanie === '') return;
      pytanie.value = ostatniePytanie;
      szukajTeraz();
    },
    zamknij() {
      menuZakresu.zwin();
    },
  };
}

/** Granica z pola liczbowego; puste i niedodatnie znaczy „granica rdzenia". */
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

/**
 * Uchwyt do źródła — przycisk kopiujący wskazanie, gdy rdzeń je podał.
 *
 * Kopia, a nie skok, bo kontrakt nie ma komendy otwierającej źródło fragmentu;
 * `sourceId` jest jedynym, czym da się po nie sięgnąć w oknie biblioteki albo
 * w rozmowie z modelem. Przycisk nazywa więc wprost to, co robi.
 */
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
        // Schowek bywa odmówiony przez przeglądarkę — wtedy podajemy wskazanie
        // wprost w potwierdzeniu, żeby Operator mógł je przepisać.
        potwierdz(`Schowek niedostępny. Wskazanie źródła: ${wskazanie}`, false),
      );
  });
  return uchwyt;
}

/** Zakres, z którego przyszedł fragment — biblioteka, historia, przestrzeń. */
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
  // Klasy dokładane jawnymi literałami, nie składane w napisie: kontrola
  // pokrycia klas CSS czyta napisy, a nazwa złożona w czasie działania jest
  // dla niej ślepa — reguła bez wołacza gnije wtedy niezauważona.
  element.className = 'dn-plakietka';
  element.classList.add(
    trafienie.score >= PROG_TRAFNOSCI_MOCNEJ ? 'dn-plakietka--sukces' : 'dn-plakietka--ostrzezenie',
  );
  element.textContent = `trafność ${trafienie.score}/100`;
  return element;
}

/**
 * Próg, powyżej którego trafność dostaje plakietkę sukcesu.
 *
 * Wartość jest umową okna, nie miarą rdzenia — kontrakt mówi tylko, że `score`
 * jest w setnych i że 100 to najbliższy możliwy. Próg stoi stałą nazwaną, żeby
 * liczba w kodzie nie udawała progu wziętego skądinąd.
 */
const PROG_TRAFNOSCI_MOCNEJ = 70;

/** Nazwy zakresów po polsku; klucze z kontraktu, nie literały okna. */
const OPIS_ZAKRESU: Record<string, string> = {
  [KnowledgeScope.Library]: 'biblioteka',
  [KnowledgeScope.History]: 'historia rozmów',
  [KnowledgeScope.Workspace]: 'pliki przestrzeni',
  [KnowledgeScope.All]: 'wszystko',
};
