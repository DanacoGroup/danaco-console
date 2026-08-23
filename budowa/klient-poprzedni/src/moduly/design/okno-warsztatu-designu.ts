import type { DesignAsset } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import { poleWyboru, przycisk, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import { opiszSkutek } from '../studio/czynnosci-warsztatu';
import { utworzKontrolke, type KontrolkaPola } from '../studio/okno-warsztatu-dokumentu';
import {
  czynnosciWarsztatu,
  type CzynnoscWarsztatuDesignu,
  type GrupaWarsztatu,
} from './czynnosci-warsztatow-designu';
import { naglowekOkna, utworzStanOkna, type StanOkna } from './stan-okna';
import type { StanDesignu } from './stan-designu';
import type { ZrodloWarsztatowDesignu } from './zrodlo-warsztatow-designu';

/**
 * Warsztat modułu Design — jedno okno na warsztat, jeden formularz na czynność.
 *
 * Pięć warsztatów (fotografia, wektor, druk, bazy zdjęciowe, publikacja) dzieli
 * ten sam budowniczy, bo dzielą to samo zadanie: wybór czynności przestawia pola,
 * pola pochodzą z katalogu, żądanie składa katalog, a odpowiedź jest opisana
 * SKUTKIEM. Pięć osobnych budowniczych byłoby pięcioma miejscami, w których
 * pomyłka w składaniu żądania mieszka osobno.
 *
 * Kontrolki bierze `utworzKontrolke` z warsztatu dokumentu Studio — ten sam
 * przełącznik rodzajów pól, bo katalogi opisują pola tymi samymi typami. Druga
 * kopia tego przełącznika rozjechałaby się z pierwszą przy pierwszym nowym
 * rodzaju pola.
 *
 * Materiał wchodzi z magazynu okna i wynik do niego wraca. Żadna czynność nie
 * zmienia materiału w miejscu — okno mówi to przy polu materiału, bo Operator ma
 * wiedzieć, że pomyłka nie kosztuje go zdjęcia źródłowego.
 *
 * Odpowiedź jest opisana liczbą, nie słowem „gotowe": nowy zasób, liczba stron,
 * udział punktów przezroczystych, liczba kafli, dostawcy, którzy nie
 * odpowiedzieli. Meldunek bez liczby nie odróżnia czynności wykonanej od
 * czynności przyjętej — a to jest wzorzec szkody, który ten moduł ma w historii.
 */
export interface OknoWarsztatuDesignu {
  element: HTMLElement;
  /** Wczytuje wykaz materiału z magazynu okna. */
  wczytaj(): Promise<void>;
  /** Przenosi ognisko do okna — woła to wyszukiwarka funkcji modułu. */
  przenieOgnisko(): void;
  odswiez(): void;
}

/** Nastawy jednego warsztatu: kod katalogu rdzenia, tytuł, rola i objaśnienie. */
export interface NastawyWarsztatuDesignu {
  kod: string;
  tytul: string;
  rola: string;
  grupa: GrupaWarsztatu;
  objasnienie: string;
}

export function utworzOknoWarsztatuDesignu(
  stan: StanDesignu,
  zrodlo: ZrodloWarsztatowDesignu,
  nastawy: NastawyWarsztatuDesignu,
): OknoWarsztatuDesignu {
  const czynnosci = czynnosciWarsztatu(nastawy.grupa);
  const okno: StanOkna = utworzStanOkna();

  const objasnienie = document.createElement('p');
  objasnienie.className = 'dn-pole-opis md-okno__objasnienie';
  objasnienie.textContent = nastawy.objasnienie;

  const wybor = poleWyboru(
    {
      etykieta: 'Czynność',
      opis: 'Wybór przestawia pola poniżej — każda czynność bierze inne wskazania.',
    },
    czynnosci.map((czynnosc) => ({ wartosc: czynnosc.komenda, etykieta: czynnosc.nazwa })),
  );

  const objasnienieCzynnosci = document.createElement('p');
  objasnienieCzynnosci.className = 'dn-pole-opis md-okno__wskaznik';

  const formularz = document.createElement('div');
  formularz.className = 'md-warsztat__pola';

  const wykonanie = przycisk('Wykonaj czynność', 'dn-btn dn-btn--sygnal');
  wykonanie.dataset['czynnosc'] = 'wykonaj';
  const odpowiedz = utworzWierszOdpowiedzi();

  const pasek = document.createElement('div');
  pasek.className = 'md-okno__pasek';
  pasek.append(wykonanie);

  okno.tresc.append(objasnienie, wybor.element, objasnienieCzynnosci, formularz, odpowiedz.element);

  const element = document.createElement('section');
  element.className = 'md-okno md-warsztat';
  element.dataset['okno'] = nastawy.kod;
  element.append(naglowekOkna(nastawy.tytul, nastawy.rola), pasek, okno.element);

  let materialy: DesignAsset[] = [];
  let kontrolki = new Map<string, KontrolkaPola>();

  /** Czynność wskazana na liście; pierwsza, gdy lista dopiero się zakłada. */
  function czynnoscBiezaca(): CzynnoscWarsztatuDesignu {
    const wybrana = czynnosci.find((czynnosc) => czynnosc.komenda === wybor.kontrolka.value);
    return wybrana ?? (czynnosci[0] as CzynnoscWarsztatuDesignu);
  }

  /** Przebudowuje formularz pod czynność wskazaną. */
  function przestawFormularz(): void {
    const czynnosc = czynnoscBiezaca();
    objasnienieCzynnosci.textContent = czynnosc.opis;
    kontrolki = new Map();
    const wiersze: HTMLElement[] = [];
    for (const pole of czynnosc.pola) {
      const kontrolka = utworzKontrolke(pole, materialy);
      kontrolki.set(pole.kod, kontrolka);
      wiersze.push(kontrolka.element);
    }
    formularz.replaceChildren(...wiersze);
    // Czynność bez ani jednego pola nie jest usterką — wykaz nastaw i wykaz
    // nośników biorą wszystko z okna. Zdanie mówi to wprost, bo pusty formularz
    // wygląda jak formularz, który się nie zbudował.
    if (czynnosc.pola.length === 0) {
      const zdanie = document.createElement('p');
      zdanie.className = 'dn-pole-opis';
      zdanie.textContent = 'Ta czynność nie bierze wskazań — pracuje na oknie modułu.';
      formularz.replaceChildren(zdanie);
    }
  }

  /** Zbiera wartości formularza w postaci, którą czyta katalog czynności. */
  function zbierzWartosci(): Record<string, string> {
    const wartosci: Record<string, string> = {};
    for (const [kod, kontrolka] of kontrolki) wartosci[kod] = kontrolka.odczytaj();
    return wartosci;
  }

  async function wykonaj(): Promise<void> {
    const czynnosc = czynnoscBiezaca();
    const zlozenie = czynnosc.zloz(zbierzWartosci(), {
      idOkna: stan.idOkna(),
      idDokumentu: stan.wybrany()?.id ?? null,
    });
    if ('odmowa' in zlozenie) {
      odpowiedz.pokaz(zlozenie.odmowa, false);
      return;
    }

    okno.ladowanie(`${czynnosc.nazwa} — czynność w toku…`);
    const wynik = await zrodlo.wykonaj(czynnosc.komenda, zlozenie.zadanie);
    if (!wynik.udany) {
      odpowiedz.pokaz(
        opisOdmowy(czynnosc.nazwa, wynik.blad?.code, wynik.blad?.message),
        false,
      );
      okno.gotowe();
      return;
    }
    odpowiedz.pokaz(`${czynnosc.nazwa}: ${opiszSkutekDesignu(wynik.wynik)}`, true);
    okno.gotowe();
    // Wynik bywa nowym zasobem magazynu, więc wykaz materiału zestarzał się
    // dokładnie w tej chwili. Odczyt idzie od razu, żeby następna czynność
    // widziała to, co przed chwilą powstało.
    await wczytaj();
  }

  async function wczytaj(): Promise<void> {
    const idOkna = stan.idOkna();
    if (idOkna === '') {
      okno.puste('Okno modułu nie jest jeszcze ustalone — warsztat nie ma skąd wziąć materiału.');
      return;
    }
    okno.ladowanie('Odczyt materiału warsztatu…');
    const wynik = await zrodlo.zasoby(idOkna);
    if (!wynik.udany || wynik.wynik === undefined) {
      // Odmowa odczytu NIE gasi okna: czynności pracują na wskazaniach wpisanych
      // wprost, a materiał z listy jest wygodą, nie warunkiem.
      okno.blad(opisOdmowy('Odczyt materiału', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    materialy = [...wynik.wynik];
    przestawFormularz();
    okno.gotowe();
  }

  function odswiez(): void {
    przestawFormularz();
  }

  wybor.kontrolka.addEventListener('change', () => przestawFormularz());
  wykonanie.addEventListener('click', () => void wykonaj());
  przestawFormularz();

  return {
    element,
    wczytaj,
    przenieOgnisko: () => wybor.kontrolka.focus(),
    odswiez,
  };
}

/**
 * Opis skutku czynności warsztatu Design.
 *
 * Bierze najpierw opis wspólny z warsztatu dokumentu (nowy zasób, liczba stron,
 * ubytek bajtów), a potem dokłada liczby, które mają znaczenie wyłącznie tutaj:
 * wymiary wyniku, drogę rachunku, udział punktów przezroczystych, liczbę kafli,
 * dostawców, którzy nie odpowiedzieli, bilanse pominięć.
 *
 * BILANS ZAMIAST CISZY także w oknie: pole, które rdzeń wypełnił powodem
 * niepowodzenia, ma być widoczne w meldunku. Meldunek „gotowe" nad odpowiedzią
 * z trzema nieudanymi rozmiarami byłby tym samym kłamstwem, przed którym broni
 * się rdzeń.
 */
export function opiszSkutekDesignu(odpowiedz: unknown): string {
  const wspolny = opiszSkutek(odpowiedz);
  if (typeof odpowiedz !== 'object' || odpowiedz === null) return wspolny;
  const tresc = odpowiedz as Record<string, unknown>;
  const czesci: string[] = [];

  if (wspolny !== 'Rdzeń przyjął czynność.') czesci.push(wspolny.replace(/\.$/, ''));

  if (typeof tresc['width'] === 'number' && typeof tresc['height'] === 'number') {
    czesci.push(`wymiary ${tresc['width']}×${tresc['height']} px`);
  }
  if (typeof tresc['computedBy'] === 'string') {
    czesci.push(
      tresc['computedBy'] === 'kanalModelu'
        ? 'policzone kanałem modelu'
        : 'policzone rachunkiem rdzenia',
    );
  }
  if (typeof tresc['pageCount'] === 'number') czesci.push(`stron w pliku: ${tresc['pageCount']}`);
  if (typeof tresc['transparentShare'] === 'number') {
    czesci.push(`punktów przezroczystych: ${Math.round(tresc['transparentShare'] * 100)}%`);
  }
  if (typeof tresc['coverage'] === 'number') {
    czesci.push(`zaznaczono: ${Math.round(tresc['coverage'] * 100)}%`);
  }
  if (typeof tresc['histogramShift'] === 'number') {
    czesci.push(`przesunięcie histogramu: ${tresc['histogramShift'].toFixed(1)}`);
  }
  if (typeof tresc['effectiveDpi'] === 'number') {
    czesci.push(`rozdzielczość skuteczna: ${tresc['effectiveDpi']} dpi`);
  }
  if (typeof tresc['rows'] === 'number' && typeof tresc['columns'] === 'number') {
    czesci.push(`kafli: ${tresc['rows']}×${tresc['columns']}`);
  }
  if (typeof tresc['pathCount'] === 'number') czesci.push(`ścieżek: ${tresc['pathCount']}`);
  if (typeof tresc['pathsAffected'] === 'number') {
    czesci.push(`ścieżek oczyszczonych: ${tresc['pathsAffected']}`);
  }
  if (typeof tresc['errors'] === 'number' && typeof tresc['warnings'] === 'number') {
    czesci.push(`wad o wadze błędu: ${tresc['errors']}, ostrzeżeń: ${tresc['warnings']}`);
  }
  if (typeof tresc['ready'] === 'boolean') {
    czesci.push(tresc['ready'] === true ? 'materiał gotowy do druku' : 'materiał NIE do druku');
  }
  if (typeof tresc['preflightSkipped'] === 'boolean' && tresc['preflightSkipped'] === true) {
    czesci.push('kontrola przeddrukowa POMINIĘTA na wyraźne żądanie');
  }
  if (typeof tresc['applied'] === 'number') czesci.push(`zasobów przetworzonych: ${tresc['applied']}`);
  if (typeof tresc['exported'] === 'number') czesci.push(`materiałów wydanych: ${tresc['exported']}`);
  if (typeof tresc['composited'] === 'number') czesci.push(`warstw złożonych: ${tresc['composited']}`);
  if (typeof tresc['regionsApplied'] === 'number') {
    czesci.push(`obszarów wyretuszowanych: ${tresc['regionsApplied']}`);
  }
  if (typeof tresc['total'] === 'number') czesci.push(`pozycji: ${tresc['total']}`);
  if (typeof tresc['removed'] === 'boolean') {
    czesci.push(tresc['removed'] === true ? 'usunięto' : 'nie było czego usunąć');
  }

  // Bilanse — pola, które rdzeń wypełnia POWODEM niepowodzenia.
  for (const [pole, opis] of [
    ['providersFailed', 'dostawcy bez odpowiedzi'],
    ['failedSizes', 'rozmiary nieudane'],
    ['failedAssetIds', 'zasoby nieudane'],
    ['failedConcepts', 'pojęcia bez ikony'],
    ['unmatchedNames', 'podstawienia niedopasowane'],
    ['gridWarnings', 'odstępstwa od siatki'],
    ['unreachableFrameIds', 'ramki nieosiągalne'],
    ['unplacedNodeIds', 'węzły nieumieszczone'],
    ['removedPathIds', 'ścieżki źródłowe usunięte'],
    ['skippedLayerAssetIds', 'warstwy pominięte'],
    ['regionsSkipped', 'obszary pominięte'],
    ['iccProfiles', 'profile ICC na serwerze'],
    ['appliedSteps', 'kroki, które weszły'],
  ] as const) {
    const wartosc = tresc[pole];
    if (Array.isArray(wartosc) && wartosc.length > 0) {
      czesci.push(`${opis}: ${wartosc.length}`);
    }
  }
  if (typeof tresc['unrecognizedRegions'] === 'number') {
    czesci.push(`obszarów nierozłożonych: ${tresc['unrecognizedRegions']}`);
  }
  if (typeof tresc['droppedRegions'] === 'number') {
    czesci.push(`obszarów odrzuconych: ${tresc['droppedRegions']}`);
  }

  return czesci.length === 0 ? 'Rdzeń przyjął czynność.' : `${czesci.join(', ')}.`;
}
