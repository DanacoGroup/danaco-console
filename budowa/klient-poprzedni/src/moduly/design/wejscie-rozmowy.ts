import { ChangeKind, type DesignAsset } from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import { slowoModelu } from './slowo-modelu';
import type { ZrodloDesignu } from './zrodlo-designu';

/**
 * Wejście rozmowy do modułu Design — którędy wynik pracy z Chat Window wchodzi
 * na okno robocze modułu.
 *
 * Rozstrzyga, czy zasób przyszedł spoza okien tego modułu, i układa zdanie o
 * tym. Kładzeniem warstwy zajmuje się Design Board, subskrypcją — złożenie
 * modułu (`modul-design.ts`).
 *
 * Miarą jest `windowId`, nie identyfikator zasobu:
 *   (1) `design.asset.generate` przepisuje `WindowId` żądania na pole `windowId`
 *       zasobu (`server/internal/core/adapter_modul_design.go`), więc zasób
 *       niesie okno, z którego padło zlecenie;
 *   (2) zdarzenie `design.asset.changed` rozgłaszane jest dla każdego zasobu
 *       odpowiedzi, zanim odpowiedź wróci do autora żądania
 *       (`adapter_modul_design_uchwyty.go`).
 * Z (2) wynika, że po identyfikatorze zasobu rozróżnić się nie da: w chwili
 * zdarzenia własnego generowania klient nie zna jeszcze id własnego zasobu.
 * Z (1) wynika miara bez wyścigu — zasób z oknem innym niż okno tego modułu
 * powstał poza oknami modułu, a jedyną drogą zlecenia Designowi spoza jego
 * okien jest rozmowa.
 *
 * Puste okno modułu miary nie psuje: Prompt Builder bez okna modułu odmawia
 * zlecenia (`okno-prompt-builder.ts`, warunek `stan.idOkna() === ''`), więc
 * zasób, który wtedy przyszedł, na pewno nie jest jego.
 *
 * Czego miara nie rozstrzyga: gdyby model wołał `design.asset.generate` z
 * identyfikatorem okna tego modułu, wynik rozmowy byłby nie do odróżnienia od
 * wyniku Prompt Buildera i na kanwę sam by nie wszedł — trafiłby do Assets
 * Panel. Z jakim `windowId` Chat Window woła komendy modułu, nie jest tu
 * zgadywane.
 */

/**
 * Co rozmowa zrobiła zasobowi.
 *
 * Trzy rodzaje zmiany kontraktu. `zasob-zdjety` istnieje mimo że rdzeń komendy
 * usuwającej dziś nie ma: `ChangeKind` dopuszcza `deleted`, a pominięcie tej
 * gałęzi kazałoby oknu położyć na kanwie zasób ogłoszony jako usunięty (ta sama
 * racja, co przy `usunZasob` w `zapis-designu.ts`).
 */
export type SkutekRozmowy = 'zasob-nowy' | 'zasob-zmieniony' | 'zasob-zdjety';

/**
 * Czy zasób przyszedł spoza okien tego modułu — czyli rozmową albo innym
 * połączeniem. Miarą jest okno zasobu, nie jego identyfikator (patrz nagłówek).
 */
export function czyZasobZRozmowy(zasob: DesignAsset, oknoModulu: string): boolean {
  return zasob.windowId !== oknoModulu;
}

/** Przekład rodzaju zmiany kontraktu na skutek, który okno robocze umie nazwać. */
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
 * Nasłuch wejścia rozmowy — druga subskrypcja tego samego zdarzenia, obok tej,
 * którą prowadzi stan modułu.
 *
 * Rozdzielenie jest celowe: stan wciąga każdy zasób do wykazu zarządcy, a to
 * nasłuchiwanie odpowiada wyłącznie na pytanie „czy to wynik rozmowy" i nie
 * zmienia wykazu. Zdarzenie zasobu z własnego okna modułu jest tu pomijane: to
 * skutek Prompt Buildera, który melduje go sam, w swoim oknie, wraz ze słowem
 * modelu.
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
 * Zdanie o skutku rozmowy — mówi, CO przyszło i czy naprawdę leży na kanwie.
 *
 * `naKanwie` nie jest zapowiedzią, tylko sprawozdaniem okna roboczego: zasób
 * już dołożony do kompozycji nie jest kładziony drugi raz i Operator ma to
 * usłyszeć zamiast patrzeć na licznik, który podskoczył bez powodu.
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
