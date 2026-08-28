import { ChangeKind, type DesignAsset } from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import { slowoModelu } from './slowo-modelu';
import type { ZrodloDesignu } from './zrodlo-designu';

/**
 * Skutek, jaki rozmowa wywarła na zasobie modułu Design: zasób nowy, zmieniony
 * albo zdjęty. Trzy wartości odpowiadają trzem rodzajom zmiany wyliczenia
 * `ChangeKind`, którym rdzeń znakuje zdarzenie zasobu.
 */
export type SkutekRozmowy = 'zasob-nowy' | 'zasob-zmieniony' | 'zasob-zdjety';

/**
 * Rozstrzyga, czy zasób przyszedł spoza okien tego modułu, czyli rozmową albo
 * innym połączeniem. Miarą jest pole `windowId` zasobu zestawione z oknem
 * modułu, a nie identyfikator zasobu.
 */
export function czyZasobZRozmowy(zasob: DesignAsset, oknoModulu: string): boolean {
  return zasob.windowId !== oknoModulu;
}

/**
 * Przekłada rodzaj zmiany niesiony przez zdarzenie zasobu na skutek rozmowy,
 * który okno robocze umie nazwać; rodzaj nierozpoznany daje zasób zmieniony.
 */
export function skutekZmiany(zmiana: ChangeKind): SkutekRozmowy {
  switch (zmiana) {
    case ChangeKind.Created:
      return 'zasob-nowy';
    case ChangeKind.Deleted:
      return 'zasob-zdjety';
    default:
      return 'zasob-zmieniony';
  }
}

/**
 * Zakłada nasłuch wejścia rozmowy na zdarzeniu zmiany zasobu i woła wskazaną
 * funkcję ze skutkiem oraz zasobem. Zdarzenie zasobu z własnego okna modułu
 * pomija, bo o nim melduje Prompt Builder. Zwraca odsubskrybowanie.
 */
export function nasluchujWejsciaRozmowy(
  zrodlo: ZrodloDesignu,
  oknoModulu: () => string,
  naSkutek: (skutek: SkutekRozmowy, zasob: DesignAsset) => void,
): Odsubskrybuj {
  return zrodlo.naZmianeZasobu((tresc) => {
    if (!czyZasobZRozmowy(tresc.asset, oknoModulu())) return;
    naSkutek(skutekZmiany(tresc.change), tresc.asset);
  });
}

/**
 * Składa zdanie o skutku rozmowy: nazywa zasób, który przyszedł, i mówi, czy
 * leży już na kanwie. Znacznik obecności na kanwie jest sprawozdaniem okna
 * roboczego, a nie zapowiedzią złożenia warstwy.
 */
export function zdanieOSkutkuRozmowy(
  skutek: SkutekRozmowy,
  zasob: DesignAsset,
  naKanwie: boolean,
): string {
  const nazwa = nazwaZasobu(zasob);
  const skad = `okno rdzenia ${zasob.windowId === '' ? '(nie podane)' : zasob.windowId}`;
  switch (skutek) {
    case 'zasob-nowy':
      return naKanwie
        ? `zasób ${nazwa} z rozmowy (${skad}) leży warstwą na kanwie. Do rdzenia warstwa ` +
            'pojedzie dopiero zapisem kompozycji — samo wejście niczego w rdzeniu nie zmienia.'
        : `zasób ${nazwa} z rozmowy (${skad}) JUŻ ma warstwę na kanwie — drugiej okno nie ` +
            'dokłada, żeby ta sama praca nie liczyła się dwa razy.';
    case 'zasob-zmieniony':
      return `rozmowa ZMIENIŁA zasób ${nazwa} (${skad}) — u zasobu rdzeń zna dziś jeden zapis, ` +
        'etykiety. Kanwa nowej warstwy nie dostaje, bo to nie jest nowa praca.';
    default:
      return `rozmowa ZDJĘŁA zasób ${nazwa} (${skad}). Warstwa, jeżeli już leży na kanwie, ` +
        'zostaje i wskazuje byt, którego rdzeń nie zna — okno nie kasuje pracy Operatora bez ' +
        'jego wiedzy.';
  }
}

/**
 * Nazwa zasobu na tabliczce — słowo modelu albo identyfikator wraz z prawdą
 * o tym, że model nie oddał ani znaku (`slowo-modelu.ts`).
 */
function nazwaZasobu(zasob: DesignAsset): string {
  const slowo = slowoModelu(zasob);
  return slowo === '' ? `${zasob.id} (model nie oddał opisu)` : `„${slowo}"`;
}
