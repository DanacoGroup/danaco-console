// Znak marki Danaco Console — sygnet (godło), logotyp i zestawienia.
// Znak czyta się »». : podwójny grot z kropką sygnału na linii bazowej.
//
// To jedyne znaki zestawu o barwach własnych — atrament i kropka sygnału są
// wpisane w plik, nie dziedziczone przez `currentColor`. Dlatego odmianę
// dobiera się do podłoża, nie do motywu: pasek kokpitu jest atramentowy
// w obu motywach, więc leży na nim odmiana „na ciemnym".
//
// Odmiana uproszczona (jeden grot) obowiązuje od 16 px w dół — poniżej tego
// boku drugi grot i prześwit między grotami przestają być czytelne.

import logoPionowyNaCiemnym from './zasoby-marki/logo-pionowy-na-ciemnym.svg?raw';
import logoPionowy from './zasoby-marki/logo-pionowy.svg?raw';
import logoPoziomyNaCiemnym from './zasoby-marki/logo-poziomy-na-ciemnym.svg?raw';
import logoPoziomy from './zasoby-marki/logo-poziomy.svg?raw';
import logotypNaCiemnym from './zasoby-marki/logotyp-na-ciemnym.svg?raw';
import logotyp from './zasoby-marki/logotyp.svg?raw';
import sygnetMonoBialy from './zasoby-marki/sygnet-mono-bialy.svg?raw';
import sygnetMonoCzarny from './zasoby-marki/sygnet-mono-czarny.svg?raw';
import sygnetNaCiemnym from './zasoby-marki/sygnet-na-ciemnym.svg?raw';
import sygnetUproszczonyNaCiemnym from './zasoby-marki/sygnet-uproszczony-na-ciemnym.svg?raw';
import sygnetUproszczony from './zasoby-marki/sygnet-uproszczony.svg?raw';
import sygnet from './zasoby-marki/sygnet.svg?raw';

/** Podłoże, na którym znak jest osadzany. */
export type PodlozeZnaku = 'jasne' | 'ciemne';

/** Odmiana znaku: samo godło albo godło zestawione z nazwą. */
export type OdmianaZnaku = 'sygnet' | 'logotyp' | 'poziomy' | 'pionowy';

/** Bok, od którego w dół obowiązuje sygnet uproszczony (jeden grot). */
export const PROG_SYGNETU_UPROSZCZONEGO = 16;

const ZNAKI: Readonly<Record<OdmianaZnaku, Readonly<Record<PodlozeZnaku, string>>>> = {
  sygnet: { jasne: sygnet, ciemne: sygnetNaCiemnym },
  logotyp: { jasne: logotyp, ciemne: logotypNaCiemnym },
  poziomy: { jasne: logoPoziomy, ciemne: logoPoziomyNaCiemnym },
  pionowy: { jasne: logoPionowy, ciemne: logoPionowyNaCiemnym },
};

const SYGNET_UPROSZCZONY: Readonly<Record<PodlozeZnaku, string>> = {
  jasne: sygnetUproszczony,
  ciemne: sygnetUproszczonyNaCiemnym,
};

const SYGNET_MONO: Readonly<Record<PodlozeZnaku, string>> = {
  jasne: sygnetMonoCzarny,
  ciemne: sygnetMonoBialy,
};

/** Zwraca źródło SVG wskazanej odmiany znaku dla danego podłoża. */
export function zrodloZnaku(odmiana: OdmianaZnaku, podloze: PodlozeZnaku): string {
  return ZNAKI[odmiana][podloze];
}

/**
 * Zwraca źródło SVG godła dobrane do boku renderowania: od 16 px w dół
 * odmiana uproszczona, powyżej pełny sygnet.
 */
export function zrodloGodla(rozmiar: number, podloze: PodlozeZnaku): string {
  return rozmiar <= PROG_SYGNETU_UPROSZCZONEGO
    ? SYGNET_UPROSZCZONY[podloze]
    : ZNAKI.sygnet[podloze];
}

/**
 * Zwraca źródło SVG godła w odmianie jednobarwnej — do tłoczenia, znaku
 * wodnego i podłoży, na których kropka sygnału łamie zasadę jednego akcentu.
 */
export function zrodloGodlaMono(podloze: PodlozeZnaku): string {
  return SYGNET_MONO[podloze];
}
