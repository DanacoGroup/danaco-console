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
 * Klucz zapisu nastaw widoku w magazynie zgodny z kluczem używanym przez moduł
 * nastaw operatora widoku, dla zgodności odczytu między obiema drogami zapisu.
 */
const KLUCZ_ZAPISU = 'dn.studio.widok';

/**
 * Magazyn nastaw widoku oparty na rdzeniu zamiast na magazynie przeglądarki;
 * zapis idzie dwustopniowo przez pamięć podręczną i odbicie lokalne.
 */
export interface MagazynWidokuRdzenia extends MagazynNastawWidoku {
  /** Napełnia pamięć podręczną nastawami z rdzenia; odmowa zostawia nastawy zastane. */
  wczytaj(): Promise<boolean>;
}

/**
 * Zaplecze magazynu nastaw widoku: źródło komend rdzenia, wskazanie okna
 * i dokumentu bieżącego oraz zdanie zwrotne o skutku zapisu nastawy.
 */
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

  /** Pola zasięgu nastaw: identyfikator dokumentu i okna; oba pola opcjonalne. */
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
        // Odbicie jest wygodą pierwszej klatki, awaria zapisu niczego nie przerywa.
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

/** Nastawy okna odtworzone z tekstu zapisu magazynu; zapis uszkodzony albo o złej postaci oddaje `null`. */
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
 * Nastawy okna przełożone na kontrakt; układ kartek rozkłada się na dwa pola
 * kontraktu — liczbę stron w rzędzie i rozkładówkę — bo to dwie różne rzeczy.
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
    // Tryb widoku niesie tu wyłącznie przełącznik trybu źródłowego, nie podgląd.
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

/** Nastawy kontraktu przełożone na nastawy okna operacyjnego Studio; brak pola bierze wartość domyślną modułu. */
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
  // Rozkładówka ma pierwszeństwo nad liczbą stron w rzędzie.
  if (nastawy.spreadView === true) okno.ukladKartek = 'rozkladowka';
  else if (nastawy.pagesPerRow !== undefined && nastawy.pagesPerRow > 1) {
    okno.ukladKartek = 'obok';
    okno.kartekWRzedzie = nastawy.pagesPerRow;
  } else if (nastawy.pagesPerRow === 1) okno.ukladKartek = 'jedna';
  if (nastawy.viewMode !== undefined) okno.trybZrodlowy = nastawy.viewMode === StudioViewMode.Source;
  return okno;
}

/** Magazyn przeglądarki dla odbicia nastaw widoku pierwszej klatki albo `null`, gdy magazyn przeglądarki jest niedostępny. */
function magazynPrzegladarki(): MagazynNastawWidoku | null {
  try {
    return globalThis.localStorage ?? null;
  } catch {
    return null;
  }
}
