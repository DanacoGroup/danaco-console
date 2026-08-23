import {
  StudioRulerUnit,
  StudioScrollMode,
  StudioSplitOrientation,
  StudioSurfaceMode,
  StudioViewMode,
  type StudioViewSetRequest,
  type StudioViewSettings,
} from '../../../../shared/contract';
import type { ZrodloPostaciStudio } from './zrodlo-postaci-studio';
import {
  domyslneNastawyWidoku,
  type MagazynNastawWidoku,
  type NastawyOperatoraWidoku,
} from './widok-nastawy-operatora';

/**
 * Nastawy widoku prowadzone przez RDZEŃ, a nie przez magazyn przeglądarki.
 *
 * ── Dlaczego to jest przełożenie, a nie nowa droga ──────────────────────────
 * `widok-nastawy-operatora.ts` od początku brał magazyn PODANY, a nie sięgał po
 * `localStorage` z globalnej przestrzeni — właśnie po to, żeby zapis dał się
 * przełożyć na rdzeń bez zmiany ani jednego wołacza. Kontrakt ma dziś
 * `studio.view.get` i `studio.view.set`, więc ten plik jest tym przełożeniem:
 * podstawia się go tam, gdzie stał magazyn przeglądarki, i nastawy przestają
 * ginąć przy przesiadce na inne urządzenie.
 *
 * ── Dlaczego magazyn jest nadal natychmiastowy ──────────────────────────────
 * `MagazynNastawWidoku` jest z natury natychmiastowy: `getItem` musi oddać
 * wartość w chwili składania okna, a rdzeń odpowiada później. Zapis dwustopniowy
 * rozwiązuje to bez kłamstwa: pamięć podręczna trzyma nastawy zapisane w rdzeniu,
 * a `wczytaj()` ją napełnia. Do chwili odpowiedzi obowiązuje odbicie w magazynie
 * przeglądarki (gdy jest) albo nastawy domyślne — okno nie stoi wtedy w miejscu
 * i nie udaje, że zna wybór Operatora.
 *
 * ── Dlaczego odbicie w przeglądarce zostaje ─────────────────────────────────
 * Bo pierwsza klatka okna rysuje się przed odpowiedzią rdzenia. Odbicie nie jest
 * drugim źródłem prawdy: rdzeń nadpisuje je przy każdym odczycie, a zapis idzie
 * do obu naraz. Bez odbicia każde wejście do modułu zaczynałoby się skalą 100 %
 * i zakładkami, choćby Operator ustawił co innego.
 *
 * ── Uczciwość niepowodzenia ─────────────────────────────────────────────────
 * Nieudany zapis nastawy widoku NIE przerywa pracy — ale i nie milczy: idzie
 * zdaniem do wołacza (`naZdanie`), więc Operator wie, że wybór nie przeżyje
 * zamknięcia okna. Cisza byłaby tu obietnicą trwałości bez pokrycia.
 */

/** Klucz zapisu — ten sam, którym posługuje się `widok-nastawy-operatora.ts`. */
const KLUCZ_ZAPISU = 'dn.studio.widok';

/** Magazyn nastaw widoku oparty na rdzeniu. */
export interface MagazynWidokuRdzenia extends MagazynNastawWidoku {
  /**
   * Napełnia pamięć podręczną nastawami z rdzenia.
   *
   * Wywołuje się raz, przy wczytaniu okna. Odmowa zostawia nastawy zastane
   * i oddaje `false` — okno ma wtedy powiedzieć, że nastawy są miejscowe.
   */
  wczytaj(): Promise<boolean>;
}

/** Zaplecze magazynu: źródło komend, wskazanie okna i dokumentu, zdanie o skutku. */
export interface ZapleczeMagazynuWidoku {
  zrodlo: ZrodloPostaciStudio;
  idDokumentu(): string;
  idOkna(): string;
  /** Zdanie o nieudanym zapisie nastawy — brak trwałości jest brakiem nazwanym. */
  naZdanie(tresc: string, powodzenie: boolean): void;
  /** Odbicie na tym urządzeniu; `null` znaczy „bez odbicia". */
  odbicie?: MagazynNastawWidoku | null;
}

export function utworzMagazynWidokuRdzenia(
  zaplecze: ZapleczeMagazynuWidoku,
): MagazynWidokuRdzenia {
  const odbicie = zaplecze.odbicie === undefined ? magazynPrzegladarki() : zaplecze.odbicie;
  let pamiec: string | null = null;

  /** Pola wskazujące zasięg nastaw; oba opcjonalne w kontrakcie. */
  function zasieg(): { documentId?: string; windowId?: string } {
    const dokument = zaplecze.idDokumentu();
    const okno = zaplecze.idOkna();
    return {
      ...(dokument === '' ? {} : { documentId: dokument }),
      ...(okno === '' ? {} : { windowId: okno }),
    };
  }

  return {
    getItem(klucz) {
      if (klucz !== KLUCZ_ZAPISU) return null;
      if (pamiec !== null) return pamiec;
      return odbicie?.getItem(klucz) ?? null;
    },

    setItem(klucz, wartosc) {
      if (klucz !== KLUCZ_ZAPISU) return;
      pamiec = wartosc;
      try {
        odbicie?.setItem(klucz, wartosc);
      } catch {
        // Odbicie jest wygodą pierwszej klatki, nie trwałością — jego awaria
        // (tryb prywatny, osadzenie w ramce) niczego nie przerywa.
      }
      const nastawy = odczytajNastawy(wartosc);
      if (nastawy === null) return;
      void zapiszWRdzeniu(nastawy);
    },

    async wczytaj() {
      const wynik = await zaplecze.zrodlo.widok(zasieg());
      if (!wynik.udany || wynik.wynik === undefined) {
        zaplecze.naZdanie(
          'Nastawy widoku nie przyjechały z rdzenia' +
            (wynik.blad === undefined ? '' : `: ${wynik.blad.message}`) +
            '. Obowiązują nastawy z tego urządzenia i nie przeniosą się na inne.',
          false,
        );
        return false;
      }
      pamiec = JSON.stringify(zNastawRdzenia(wynik.wynik.settings));
      try {
        odbicie?.setItem(KLUCZ_ZAPISU, pamiec);
      } catch {
        // Jak wyżej: odbicie nie jest powodem, żeby cokolwiek przerywać.
      }
      return true;
    },
  };

  async function zapiszWRdzeniu(nastawy: NastawyOperatoraWidoku): Promise<void> {
    const wynik = await zaplecze.zrodlo.ustawWidok({ ...zasieg(), ...naNastawyRdzenia(nastawy) });
    if (wynik.udany) return;
    zaplecze.naZdanie(
      'Nastawa widoku NIE zapisała się w rdzeniu' +
        (wynik.blad === undefined ? '' : `: ${wynik.blad.message}`) +
        '. Działa na tym urządzeniu i tu zostanie — na innym Operator zobaczy nastawę wcześniejszą.',
      false,
    );
  }
}

/** Nastawy okna z zapisu; zapis uszkodzony oddaje `null`. */
function odczytajNastawy(zapis: string): NastawyOperatoraWidoku | null {
  try {
    const tresc = JSON.parse(zapis) as unknown;
    if (typeof tresc !== 'object' || tresc === null) return null;
    return { ...domyslneNastawyWidoku(), ...(tresc as Partial<NastawyOperatoraWidoku>) };
  } catch {
    return null;
  }
}

/**
 * Nastawy okna przełożone na kontrakt.
 *
 * Układ kartek okna („jedna", „obok", „rozkładówka") rozkłada się w kontrakcie na
 * DWA pola — liczbę stron w rzędzie i rozkładówkę — bo to dwie różne rzeczy:
 * rozkładówka ma stronę otwarcia po prawej, a dwie kartki obok siebie nie mają.
 */
export function naNastawyRdzenia(nastawy: NastawyOperatoraWidoku): StudioViewSetRequest {
  return {
    surfaceMode:
      nastawy.trybDokumentow === 'podzial' ? StudioSurfaceMode.Split : StudioSurfaceMode.Tabs,
    splitOrientation:
      nastawy.kierunekPodzialu === 'poziomy'
        ? StudioSplitOrientation.Horizontal
        : StudioSplitOrientation.Vertical,
    splitRatio: nastawy.udzialPodzialu,
    // Tryb widoku niesie tu wyłącznie przełącznik trybu źródłowego. Podglądu
    // wydruku i różnicy okno przestawia własną drogą (`ustawTryb`), a nadpisywanie
    // ich stąd zabierałoby Operatorowi tryb, w którym właśnie pracuje.
    viewMode: nastawy.trybZrodlowy ? StudioViewMode.Source : StudioViewMode.Edit,
    zoomPercent: nastawy.skala,
    rulersVisible: nastawy.linijkiWidoczne,
    rulerUnit: nastawy.jednostka === 'cal' ? StudioRulerUnit.Inch : StudioRulerUnit.Millimeter,
    marginGuides: nastawy.graniceMarginesow,
    pagesPerRow: nastawy.ukladKartek === 'obok' ? nastawy.kartekWRzedzie : 1,
    spreadView: nastawy.ukladKartek === 'rozkladowka',
    scrollMode:
      nastawy.przewijanie === 'strona-po-stronie'
        ? StudioScrollMode.Page
        : StudioScrollMode.Continuous,
  };
}

/** Nastawy kontraktu przełożone na nastawy okna; brak pola bierze wartość domyślną. */
export function zNastawRdzenia(nastawy: StudioViewSettings): NastawyOperatoraWidoku {
  const okno = domyslneNastawyWidoku();
  if (nastawy.surfaceMode !== undefined) {
    okno.trybDokumentow = nastawy.surfaceMode === StudioSurfaceMode.Split ? 'podzial' : 'zakladki';
  }
  if (nastawy.splitOrientation !== undefined) {
    okno.kierunekPodzialu =
      nastawy.splitOrientation === StudioSplitOrientation.Horizontal ? 'poziomy' : 'pionowy';
  }
  if (nastawy.splitRatio !== undefined) okno.udzialPodzialu = nastawy.splitRatio;
  if (nastawy.zoomPercent !== undefined) okno.skala = nastawy.zoomPercent;
  if (nastawy.rulersVisible !== undefined) okno.linijkiWidoczne = nastawy.rulersVisible;
  if (nastawy.rulerUnit !== undefined) {
    okno.jednostka = nastawy.rulerUnit === StudioRulerUnit.Inch ? 'cal' : 'mm';
  }
  if (nastawy.marginGuides !== undefined) okno.graniceMarginesow = nastawy.marginGuides;
  if (nastawy.scrollMode !== undefined) {
    okno.przewijanie =
      nastawy.scrollMode === StudioScrollMode.Page ? 'strona-po-stronie' : 'ciagle';
  }
  // Rozkładówka ma pierwszeństwo nad liczbą stron w rzędzie: jedno wyklucza
  // drugie, a rdzeń wolno mu oddać oba naraz.
  if (nastawy.spreadView === true) okno.ukladKartek = 'rozkladowka';
  else if (nastawy.pagesPerRow !== undefined && nastawy.pagesPerRow > 1) {
    okno.ukladKartek = 'obok';
    okno.kartekWRzedzie = nastawy.pagesPerRow;
  } else if (nastawy.pagesPerRow === 1) okno.ukladKartek = 'jedna';
  if (nastawy.viewMode !== undefined) okno.trybZrodlowy = nastawy.viewMode === StudioViewMode.Source;
  return okno;
}

/** Magazyn przeglądarki albo `null`, gdy niedostępny. */
function magazynPrzegladarki(): MagazynNastawWidoku | null {
  try {
    return globalThis.localStorage ?? null;
  } catch {
    return null;
  }
}
