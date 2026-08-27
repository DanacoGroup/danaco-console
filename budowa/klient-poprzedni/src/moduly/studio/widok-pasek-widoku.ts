import {
  NASTAWY_MARGINESOW,
  nosniki,
  POWOD_NIETRWALOSCI_STRONY,
  type JednostkaMiary,
} from './nastawy-strony';
import { przyciskBezKomendy } from '../../modele/kontrolki-formularza';
import { NASTAWY_SKALI, SKALE_GOTOWE } from './widok-skali';
import {
  czyDrukDostepny,
  czytajNastawyDruku,
  opiszNastawyDruku,
  POWOD_BRAKU_DRUKU,
  type AdiustacjaDruku,
  type NastawyDruku,
  type ZakresDruku,
} from './widok-druku';
import {
  SKALA_DOLNA,
  SKALA_GORNA,
  type KierunekPodzialu,
  type NastawyOperatoraWidoku,
  type TrybDwochDokumentow,
  type TrybPrzewijania,
  type UkladKartek,
} from './widok-nastawy-operatora';

/**
 * Pasek widoku pływa przy krawędzi powierzchni i chowa nastawy rzadsze
 * w nakładce, aby stała kolumna nie zabierała miejsca kartce. Czynności
 * zebrane niżej pasek zleca powierzchni, bo to ona przechowuje stan widoku
 * i skalę.
 */
export interface CzynnosciPaskaWidoku {
  /** Skala wpisana albo z suwaka, w procentach. */
  naSkale(procent: number): void;
  /** Nastawa gotowa skali — rachunek należy do powierzchni, bo ona zna swój widok. */
  naNastaweSkali(kod: string): void;
  naUklad(uklad: UkladKartek, kartekWRzedzie: number): void;
  naPrzewijanie(tryb: TrybPrzewijania): void;
  naLinijki(widoczne: boolean): void;
  naJednostke(jednostka: JednostkaMiary): void;
  naGraniceMarginesow(widoczne: boolean): void;
  /** Tryb źródłowy ze znacznikami — przełącznik, nie tryb domyślny. */
  naTrybZrodlowy(wlaczony: boolean): void;
  /** Podgląd wydruku jako tryb tego okna. */
  naPodgladWydruku(wlaczony: boolean): void;
  /** Skok do wskazanej kartki. */
  naKartke(numer: number): void;
  /** Zakładki albo podział powierzchni — dwa tryby równorzędne. */
  naTrybDokumentow(tryb: TrybDwochDokumentow): void;
  naKierunekPodzialu(kierunek: KierunekPodzialu): void;
  /** Panele jako nakładki albo jako stałe kolumny — wybór Operatora. */
  naUkladPaneli(uklad: 'nakladka' | 'kolumny'): void;
  /* ── Nastawy kartki ─────────────────────────────────────────────────────── */

  /** Nośnik z wykazu — oznaczenie ISO albo nazwa koperty. */
  naNosnik(oznaczenie: string): void;
  /** Nośnik własny podany wymiarami w milimetrach. */
  naNosnikWlasny(szerokoscMm: number, wysokoscMm: number): void;
  naOrientacje(pozioma: boolean): void;
  /** Nastawa gotowa marginesów — wąskie, normalne, szerokie, do oprawy. */
  naNastaweMarginesow(kod: string): void;
  naMargines(ktory: 'gora' | 'dol' | 'lewy' | 'prawy', milimetry: number): void;
  /** Margines na oprawę w milimetrach. */
  naOprawe(milimetry: number): void;
  naStroneOprawy(strona: 'wewnatrz' | 'gora'): void;
  /** Marginesy odbicia dla druku dwustronnego. */
  naOdbicia(wlaczone: boolean): void;
  /** Drukowanie wedle nastaw okna druku. */
  naDruk(nastawy: NastawyDruku): void;
  /** Szybkie drukowanie: ostatnie nastawy, bez okna nastaw. */
  naSzybkiDruk(): void;
}

/**
 * Pasek widoku wraz z jego sterowaniem: element gotowy do osadzenia w oknie,
 * odświeżenie kontrolek po zmianie nastaw powierzchni oraz zamknięcie nakładki,
 * gdy ognisko wraca do dokumentu.
 */
export interface PasekWidoku {
  element: HTMLElement;
  /** Przestawia kontrolki wedle nastaw obowiązujących. */
  odswiez(
    nastawy: NastawyOperatoraWidoku,
    podglad: boolean,
    liczbaStron: number,
    kartkaBiezaca: number,
  ): void;
  /** Zamyka nakładkę nastaw — wołane, gdy powierzchnia bierze ognisko. */
  zwin(): void;
}

export function utworzPasekWidoku(
  nastawy: NastawyOperatoraWidoku,
  czynnosci: CzynnosciPaskaWidoku,
  ukladPaneliPoczatkowy: 'nakladka' | 'kolumny' = 'nakladka',
): PasekWidoku {
  /* ── Rząd stały: skala, kartki, przełączniki ─────────────────────────────── */

  const suwakSkali = document.createElement('input');
  suwakSkali.type = 'range';
  suwakSkali.className = 'ms-widok__suwak';
  suwakSkali.min = String(SKALA_DOLNA);
  suwakSkali.max = String(SKALA_GORNA);
  suwakSkali.step = '5';
  suwakSkali.value = String(nastawy.skala);
  suwakSkali.setAttribute('aria-label', 'Skala widoku w procentach');
  suwakSkali.addEventListener('input', () => czynnosci.naSkale(Number(suwakSkali.value)));

  const poleSkali = document.createElement('input');
  poleSkali.type = 'number';
  poleSkali.className = 'dn-pole-kontrolka ms-widok__skala';
  poleSkali.min = String(SKALA_DOLNA);
  poleSkali.max = String(SKALA_GORNA);
  poleSkali.value = String(nastawy.skala);
  poleSkali.setAttribute('aria-label', 'Skala widoku wpisana w procentach');
  poleSkali.addEventListener('change', () => czynnosci.naSkale(Number(poleSkali.value)));

  const nastawySkali = document.createElement('select');
  nastawySkali.className = 'dn-pole-kontrolka ms-widok__nastawy-skali';
  nastawySkali.setAttribute('aria-label', 'Nastawy gotowe skali widoku');
  const bezNastawy = document.createElement('option');
  bezNastawy.value = '';
  bezNastawy.textContent = 'skala własna';
  nastawySkali.append(bezNastawy);
  for (const nastawa of NASTAWY_SKALI) {
    const opcja = document.createElement('option');
    opcja.value = nastawa.kod;
    opcja.textContent = nastawa.nazwa;
    opcja.title = nastawa.opis;
    nastawySkali.append(opcja);
  }
  for (const procent of SKALE_GOTOWE) {
    const opcja = document.createElement('option');
    opcja.value = `procent-${procent}`;
    opcja.textContent = `${procent} %`;
    nastawySkali.append(opcja);
  }
  nastawySkali.addEventListener('change', () => {
    const wybrana = nastawySkali.value;
    if (wybrana === '') return;
    if (wybrana.startsWith('procent-')) {
      czynnosci.naSkale(Number(wybrana.slice('procent-'.length)));
      return;
    }
    czynnosci.naNastaweSkali(wybrana);
  });

  const kartkaWstecz = przycisk('◂ kartka', 'kartka-wstecz', () => skok(-1),
    'Przewija do kartki poprzedniej. Przy przewijaniu strona po stronie jest to jedyny sposób ' +
    'przejścia — i taki jest sens tej nastawy.');
  const kartkaDalej = przycisk('kartka ▸', 'kartka-dalej', () => skok(1),
    'Przewija do kartki następnej.');

  const licznikKartek = document.createElement('span');
  licznikKartek.className = 'dn-plakietka ms-widok__licznik';

  let kartkaBiezaca = 1;
  let stronOgolem = 1;

  function skok(kierunek: number): void {
    czynnosci.naKartke(Math.min(Math.max(kartkaBiezaca + kierunek, 1), Math.max(1, stronOgolem)));
  }

  const linijki = przelacznik('Linijki', nastawy.linijkiWidoczne, (wlaczone) =>
    czynnosci.naLinijki(wlaczone),
    'Pokazuje albo ukrywa linijkę poziomą i pionową wraz z ich chwytami.');

  const granice = przelacznik('Granice marginesów', nastawy.graniceMarginesow, (wlaczone) =>
    czynnosci.naGraniceMarginesow(wlaczone),
    'Rysuje w treści kreskę tam, gdzie kończy się pole pisania — granica jest widoczna, ' +
    'a nie domyślana z położenia liter.');

  const zrodlowy = przelacznik('Tryb źródłowy', nastawy.trybZrodlowy, (wlaczony) =>
    czynnosci.naTrybZrodlowy(wlaczony),
    'Pokazuje treść ze znacznikami, tak jak jedzie do rdzenia. Jest to PRZEŁĄCZNIK obok widoku ' +
    'formatowanego, a nie widok domyślny: pismo pisze się na kartce, nie w znacznikach.');

  const podgladWydruku = przelacznik('Podgląd wydruku', false, (wlaczony) =>
    czynnosci.naPodgladWydruku(wlaczony),
    'Tryb tego samego okna, nie osobne okno: kartka w nośniku i orientacji NAPRAWDĘ ustawionych, ' +
    'wraz z paginacją, nagłówkiem, stopką i numeracją. Pisanie w tym trybie jest wyłączone.');

  // Szybkie drukowanie stoi w rzędzie stałym; bez sterownika druku jest przyciskiem
  // z nazwanym powodem.
  const szybkiDruk = czyDrukDostepny()
    ? przycisk('Szybkie drukowanie', 'szybki-druk', () => czynnosci.naSzybkiDruk(),
        'Drukuje od razu, ostatnimi nastawami druku, bez okna nastaw.')
    : przyciskBezKomendy('Szybkie drukowanie', POWOD_BRAKU_DRUKU);

  const rozwin = document.createElement('button');
  rozwin.type = 'button';
  rozwin.className = 'dn-btn dn-btn--sm dn-btn--zarys ms-widok__rozwin';
  rozwin.textContent = 'Nastawy widoku ▾';
  rozwin.dataset['czynnosc'] = 'nastawy-widoku';
  rozwin.setAttribute('aria-expanded', 'false');
  rozwin.title =
    'Rozwija nakładkę z nastawami rzadszymi: układ kartek, przewijanie, jednostka linijki, ' +
    'tryb dwóch dokumentów, układ paneli. Nakładka stoi nad treścią i schodzi Escapem — ' +
    'powierzchnia należy do dokumentu, nie do paneli.';
  rozwin.addEventListener('click', () => przestawNakladke(!nakladkaOtwarta()));

  const rzad = document.createElement('div');
  rzad.className = 'ms-widok__rzad';
  rzad.append(
    obudowa('Skala', poleSkali),
    suwakSkali,
    nastawySkali,
    kartkaWstecz,
    licznikKartek,
    kartkaDalej,
    linijki.element,
    granice.element,
    zrodlowy.element,
    podgladWydruku.element,
    szybkiDruk,
    rozwin,
  );

  /* ── Nakładka nastaw rzadszych ───────────────────────────────────────────── */

  const uklad = wybor(
    'Układ kartek',
    [
      { wartosc: 'jedna', etykieta: 'Jedna kartka w rzędzie' },
      { wartosc: 'obok', etykieta: 'Kartki obok siebie' },
      { wartosc: 'rozkladowka', etykieta: 'Rozkładówka jak w książce' },
    ],
    (wartosc) => czynnosci.naUklad(wartosc as UkladKartek, Number(wRzedzie.kontrolka.value)),
  );
  uklad.kontrolka.value = nastawy.ukladKartek;

  const wRzedzie = liczba('Kartek w rzędzie', nastawy.kartekWRzedzie, 1, 8, (wartosc) =>
    czynnosci.naUklad(uklad.kontrolka.value as UkladKartek, wartosc),
  );

  const przewijanie = wybor(
    'Przewijanie',
    [
      { wartosc: 'ciagle', etykieta: 'Ciągłe' },
      { wartosc: 'strona-po-stronie', etykieta: 'Strona po stronie' },
    ],
    (wartosc) => czynnosci.naPrzewijanie(wartosc as TrybPrzewijania),
  );
  przewijanie.kontrolka.value = nastawy.przewijanie;

  const jednostka = wybor(
    'Jednostka linijki',
    [
      { wartosc: 'mm', etykieta: 'Milimetry' },
      { wartosc: 'cal', etykieta: 'Cale' },
    ],
    (wartosc) => czynnosci.naJednostke(wartosc as JednostkaMiary),
  );
  jednostka.kontrolka.value = nastawy.jednostka;

  const trybDokumentow = wybor(
    'Dwa dokumenty',
    [
      { wartosc: 'zakladki', etykieta: 'Zakładki — jeden na całej powierzchni' },
      { wartosc: 'podzial', etykieta: 'Podział powierzchni — oba naraz' },
    ],
    (wartosc) => czynnosci.naTrybDokumentow(wartosc as TrybDwochDokumentow),
  );
  trybDokumentow.kontrolka.value = nastawy.trybDokumentow;

  const kierunek = wybor(
    'Kierunek podziału',
    [
      { wartosc: 'pionowy', etykieta: 'Pionowy — obok siebie' },
      { wartosc: 'poziomy', etykieta: 'Poziomy — jeden pod drugim' },
    ],
    (wartosc) => czynnosci.naKierunekPodzialu(wartosc as KierunekPodzialu),
  );
  kierunek.kontrolka.value = nastawy.kierunekPodzialu;

  const ukladPaneli = wybor(
    'Panele okna',
    [
      { wartosc: 'nakladka', etykieta: 'Na żądanie, nakładką nad treścią' },
      { wartosc: 'kolumny', etykieta: 'Stałe kolumny obok treści' },
    ],
    (wartosc) => czynnosci.naUkladPaneli(wartosc as 'nakladka' | 'kolumny'),
  );
  ukladPaneli.kontrolka.value = ukladPaneliPoczatkowy;

  const oPanelach = document.createElement('p');
  oPanelach.className = 'dn-pole-opis';
  oPanelach.textContent =
    'Powierzchnia należy do dokumentu: dymki komentarzy, panel Redaktora i historia wersji ' +
    'wchodzą nakładką na żądanie i schodzą, gdy nie są używane. Stałe kolumny zostają jako ' +
    'tryb do wyboru — Operator z szerokim ekranem może je chcieć — ale nie jako postać domyślna.';

  /* ── Nastawy kartki: nośnik, orientacja, marginesy, oprawa ───────────────── */

  const nosnik = wybor(
    'Nośnik',
    nosniki().map((pozycja) => ({
      wartosc: pozycja.oznaczenie,
      etykieta:
        `${pozycja.oznaczenie} (${pozycja.szerokoscMm}×${pozycja.wysokoscMm} mm)` +
        `${pozycja.rodzaj === 'koperta' ? ' — koperta' : ''}`,
    })),
    (wartosc) => czynnosci.naNosnik(wartosc),
  );

  const szerokoscWlasna = liczba('Szerokość własna (mm)', 210, 1, 2000, () => undefined);
  const wysokoscWlasna = liczba('Wysokość własna (mm)', 297, 1, 2000, () => undefined);
  const nadajWlasny = przycisk(
    'Ustaw format własny',
    'nosnik-wlasny',
    () =>
      czynnosci.naNosnikWlasny(
        Number(szerokoscWlasna.kontrolka.value),
        Number(wysokoscWlasna.kontrolka.value),
      ),
    'Format podany wymiarami — Operator nie jest uwięziony w wykazie nośników. Wymiar niedodatni ' +
      'jest odrzucany: kartka bez wymiaru nie jest nastawą, jest usterką.',
  );

  const orientacja = wybor(
    'Orientacja',
    [
      { wartosc: 'pionowa', etykieta: 'Pionowa' },
      { wartosc: 'pozioma', etykieta: 'Pozioma' },
    ],
    (wartosc) => czynnosci.naOrientacje(wartosc === 'pozioma'),
  );

  const nastawaMarginesow = wybor(
    'Nastawy marginesów',
    [
      { wartosc: '', etykieta: 'własne — cztery pola poniżej' },
      ...NASTAWY_MARGINESOW.map((nastawa) => ({
        wartosc: nastawa.kod,
        etykieta: nastawa.nazwa,
      })),
    ],
    (wartosc) => {
      if (wartosc === '') return;
      czynnosci.naNastaweMarginesow(wartosc);
    },
  );
  for (const opcja of Array.from(nastawaMarginesow.kontrolka.options)) {
    opcja.title = NASTAWY_MARGINESOW.find((nastawa) => nastawa.kod === opcja.value)?.opis ?? '';
  }

  const marginesy = new Map<'gora' | 'dol' | 'lewy' | 'prawy', HTMLInputElement>();
  const rzadMarginesow = document.createElement('div');
  rzadMarginesow.className = 'ms-widok__rzad';
  for (const strona of ['gora', 'dol', 'lewy', 'prawy'] as const) {
    const kontrolka = liczba(`Margines ${strona} (mm)`, 20, 0, 200, (wartosc) =>
      czynnosci.naMargines(strona, wartosc),
    );
    marginesy.set(strona, kontrolka.kontrolka);
    rzadMarginesow.append(kontrolka.element);
  }

  const oprawa = liczba('Oprawa (mm)', 0, 0, 100, (wartosc) => czynnosci.naOprawe(wartosc));
  const stronaOprawy = wybor(
    'Strona oprawy',
    [
      { wartosc: 'wewnatrz', etykieta: 'Wewnątrz (przy rowku)' },
      { wartosc: 'gora', etykieta: 'U góry' },
    ],
    (wartosc) => czynnosci.naStroneOprawy(wartosc === 'gora' ? 'gora' : 'wewnatrz'),
  );

  const odbicia = przelacznik('Marginesy odbicia', false, (wlaczone) =>
    czynnosci.naOdbicia(wlaczone),
    'Druk dwustronny: margines „lewy" staje się WEWNĘTRZNYM — na stronie nieparzystej stoi po ' +
      'lewej, na parzystej po prawej. Bez tego pas oprawy wypadałby raz w rowku, raz na krawędzi.',
  );

  const oTrwalosci = document.createElement('p');
  oTrwalosci.className = 'dn-pole-opis';
  oTrwalosci.textContent = POWOD_NIETRWALOSCI_STRONY;

  /* ── Druk: nastawy, drukowanie i szybkie drukowanie ──────────────────────── */

  let nastawyDruku = czytajNastawyDruku();

  const zakresDruku = wybor(
    'Zakres stron',
    [
      { wartosc: 'wszystkie', etykieta: 'Wszystkie strony' },
      { wartosc: 'biezaca', etykieta: 'Strona bieżąca' },
      { wartosc: 'podany', etykieta: 'Zakres podany' },
    ],
    (wartosc) => {
      nastawyDruku = { ...nastawyDruku, zakres: wartosc as ZakresDruku };
      odswiezDruk();
    },
  );
  zakresDruku.kontrolka.value = nastawyDruku.zakres;

  const odStrony = liczba('Od strony', nastawyDruku.odStrony, 1, 9999, (wartosc) => {
    nastawyDruku = { ...nastawyDruku, odStrony: wartosc };
    odswiezDruk();
  });
  const doStrony = liczba('Do strony', nastawyDruku.doStrony, 1, 9999, (wartosc) => {
    nastawyDruku = { ...nastawyDruku, doStrony: wartosc };
    odswiezDruk();
  });

  const kopie = liczba('Kopie', nastawyDruku.kopie, 1, 99, (wartosc) => {
    nastawyDruku = { ...nastawyDruku, kopie: wartosc };
    odswiezDruk();
  });

  const dwustronny = przelacznik('Druk dwustronny', nastawyDruku.dwustronny, (wlaczony) => {
    nastawyDruku = { ...nastawyDruku, dwustronny: wlaczony };
    odswiezDruk();
  },
  'Życzenie przenoszone do okna drukarki — zatwierdza je drukarka, nie strona. Marginesy odbicia ' +
    'ustawia się osobno, wyżej, bo dotyczą postaci kartki, a nie samego wydruku.');

  const skalaDruku = liczba('Skala druku (%)', nastawyDruku.skala, SKALA_DOLNA, SKALA_GORNA, (wartosc) => {
    nastawyDruku = { ...nastawyDruku, skala: wartosc };
    odswiezDruk();
  });

  const adiustacjaDruku = wybor(
    'Adiustacja na wydruku',
    [
      { wartosc: 'po-zmianach', etykieta: 'Tekst po zmianach — bez adiustacji' },
      { wartosc: 'z-adiustacja', etykieta: 'Z adiustacją — zmiany i komentarze widoczne' },
    ],
    (wartosc) => {
      nastawyDruku = { ...nastawyDruku, adiustacja: wartosc as AdiustacjaDruku };
      odswiezDruk();
    },
  );
  adiustacjaDruku.kontrolka.value = nastawyDruku.adiustacja;

  const podsumowanieDruku = document.createElement('p');
  podsumowanieDruku.className = 'dn-pole-opis';

  const pasDruku = document.createElement('div');
  pasDruku.className = 'ms-widok__rzad';
  if (czyDrukDostepny()) {
    pasDruku.append(
      przycisk('Drukuj', 'drukuj', () => czynnosci.naDruk(nastawyDruku),
        'Drukuje TO, co pokazuje podgląd wydruku: nośnik, orientację, marginesy, paginację, ' +
          'nagłówek i stopkę. Wydruk różniący się od podglądu byłby usterką.'),
    );
  } else {
    pasDruku.append(przyciskBezKomendy('Drukuj', POWOD_BRAKU_DRUKU));
  }

  const grupaDruku = document.createElement('div');
  grupaDruku.className = 'ms-widok__druk';
  grupaDruku.append(
    zakresDruku.element,
    odStrony.element,
    doStrony.element,
    kopie.element,
    dwustronny.element,
    skalaDruku.element,
    adiustacjaDruku.element,
    podsumowanieDruku,
    pasDruku,
  );

  /** Ile stron i która bieżąca — do podsumowania nastaw druku. */
  function odswiezDruk(): void {
    const podane = nastawyDruku.zakres === 'podany';
    odStrony.element.hidden = !podane;
    doStrony.element.hidden = !podane;
    podsumowanieDruku.textContent = opiszNastawyDruku(nastawyDruku, stronOgolem, kartkaBiezaca);
  }

  const nakladka = document.createElement('div');
  nakladka.className = 'ms-widok__nakladka';
  nakladka.hidden = true;
  nakladka.setAttribute('role', 'group');
  nakladka.setAttribute('aria-label', 'Nastawy widoku powierzchni');
  nakladka.append(
    uklad.element,
    wRzedzie.element,
    przewijanie.element,
    jednostka.element,
    trybDokumentow.element,
    kierunek.element,
    ukladPaneli.element,
    oPanelach,
    nosnik.element,
    orientacja.element,
    szerokoscWlasna.element,
    wysokoscWlasna.element,
    nadajWlasny,
    nastawaMarginesow.element,
    rzadMarginesow,
    oprawa.element,
    stronaOprawy.element,
    odbicia.element,
    oTrwalosci,
    grupaDruku,
  );

  const element = document.createElement('div');
  element.className = 'ms-widok';
  element.append(rzad, nakladka);

  function nakladkaOtwarta(): boolean {
    return !nakladka.hidden;
  }

  function przestawNakladke(otwarta: boolean): void {
    nakladka.hidden = !otwarta;
    rozwin.setAttribute('aria-expanded', otwarta ? 'true' : 'false');
    rozwin.textContent = otwarta ? 'Nastawy widoku ▴' : 'Nastawy widoku ▾';
  }

  nakladka.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key === 'Escape') przestawNakladke(false);
  });

  return {
    element,

    odswiez(nowe, podglad, liczbaStron, kartka) {
      stronOgolem = Math.max(1, liczbaStron);
      kartkaBiezaca = Math.min(Math.max(1, kartka), stronOgolem);
      if (poleSkali.value !== String(nowe.skala)) poleSkali.value = String(nowe.skala);
      if (suwakSkali.value !== String(nowe.skala)) suwakSkali.value = String(nowe.skala);
      licznikKartek.textContent = `kartka ${kartkaBiezaca} z ${stronOgolem}`;
      linijki.ustaw(nowe.linijkiWidoczne);
      granice.ustaw(nowe.graniceMarginesow);
      zrodlowy.ustaw(nowe.trybZrodlowy);
      podgladWydruku.ustaw(podglad);
      uklad.kontrolka.value = nowe.ukladKartek;
      wRzedzie.kontrolka.value = String(nowe.kartekWRzedzie);
      wRzedzie.element.hidden = nowe.ukladKartek !== 'obok';
      przewijanie.kontrolka.value = nowe.przewijanie;
      jednostka.kontrolka.value = nowe.jednostka;
      trybDokumentow.kontrolka.value = nowe.trybDokumentow;
      kierunek.element.hidden = nowe.trybDokumentow !== 'podzial';
      kierunek.kontrolka.value = nowe.kierunekPodzialu;
      // Tryb źródłowy i podgląd wydruku wykluczają się — wyłączony przełącznik
      // mówi to wprost.
      zrodlowy.kontrolka.disabled = podglad;
      podgladWydruku.kontrolka.disabled = nowe.trybZrodlowy;
      odswiezDruk();
    },

    zwin: () => przestawNakladke(false),
  };
}

/* ── Kawałki wspólne paska ────────────────────────────────────────────────── */

function przycisk(
  nazwa: string,
  kod: string,
  czynnosc: () => void,
  objasnienie: string,
): HTMLButtonElement {
  const element = document.createElement('button');
  element.type = 'button';
  element.className = 'dn-btn dn-btn--sm dn-btn--duch';
  element.textContent = nazwa;
  element.dataset['czynnosc'] = kod;
  element.title = objasnienie;
  element.setAttribute('aria-description', objasnienie);
  element.addEventListener('click', czynnosc);
  return element;
}

/**
 * Przełącznik paska wraz z etykietą i objaśnieniem widocznym jako podpowiedź
 * oraz jako opis dostępności kontrolki dla czytnika ekranu.
 */
function przelacznik(
  etykieta: string,
  wlaczony: boolean,
  naZmiane: (wlaczony: boolean) => void,
  objasnienie: string,
): { element: HTMLElement; kontrolka: HTMLInputElement; ustaw(wlaczony: boolean): void } {
  const kontrolka = document.createElement('input');
  kontrolka.type = 'checkbox';
  kontrolka.className = 'dn-przelacznik';
  kontrolka.checked = wlaczony;
  kontrolka.title = objasnienie;
  kontrolka.setAttribute('aria-description', objasnienie);
  kontrolka.addEventListener('change', () => naZmiane(kontrolka.checked));

  const napis = document.createElement('span');
  napis.className = 'ms-widok__etykieta';
  napis.textContent = etykieta;

  const element = document.createElement('label');
  element.className = 'ms-widok__przelacznik';
  element.append(kontrolka, napis);
  return {
    element,
    kontrolka,
    ustaw(nowy) {
      if (kontrolka.checked !== nowy) kontrolka.checked = nowy;
    },
  };
}

function wybor(
  etykieta: string,
  pozycje: readonly { wartosc: string; etykieta: string }[],
  naZmiane: (wartosc: string) => void,
): { element: HTMLElement; kontrolka: HTMLSelectElement } {
  const kontrolka = document.createElement('select');
  kontrolka.className = 'dn-pole-kontrolka';
  kontrolka.setAttribute('aria-label', etykieta);
  for (const pozycja of pozycje) {
    const opcja = document.createElement('option');
    opcja.value = pozycja.wartosc;
    opcja.textContent = pozycja.etykieta;
    kontrolka.append(opcja);
  }
  kontrolka.addEventListener('change', () => naZmiane(kontrolka.value));
  return { element: obudowa(etykieta, kontrolka), kontrolka };
}

function liczba(
  etykieta: string,
  domyslna: number,
  dol: number,
  gora: number,
  naZmiane: (wartosc: number) => void,
): { element: HTMLElement; kontrolka: HTMLInputElement } {
  const kontrolka = document.createElement('input');
  kontrolka.type = 'number';
  kontrolka.className = 'dn-pole-kontrolka ms-widok__liczba';
  kontrolka.value = String(domyslna);
  kontrolka.min = String(dol);
  kontrolka.max = String(gora);
  kontrolka.setAttribute('aria-label', etykieta);
  kontrolka.addEventListener('change', () => naZmiane(Number(kontrolka.value)));
  return { element: obudowa(etykieta, kontrolka), kontrolka };
}

function obudowa(etykieta: string, kontrolka: HTMLElement): HTMLElement {
  const napis = document.createElement('span');
  napis.className = 'ms-widok__etykieta';
  napis.textContent = etykieta;

  const element = document.createElement('label');
  element.className = 'ms-widok__pole';
  element.append(napis, kontrolka);
  return element;
}
