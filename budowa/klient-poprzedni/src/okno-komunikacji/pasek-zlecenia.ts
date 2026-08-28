import './pasek-zlecenia.css';

import { EventType } from '../../../shared/contract';
import { utworzZrodloRozszerzen } from '../moduly/agents/zrodlo-rozszerzen';
import type { Kanal } from '../protokol/kanal';
import { utworzRejestrAgentow } from '../sterowanie/rejestr-agentow';
import type { RejestrKanalow } from '../sterowanie/rejestr-kanalow';
import { utworzDyktowanie, type Dyktowanie } from './dyktowanie/dyktowanie';
import { utworzSterKatalogow } from './ster-katalogow';
import { utworzSterMikrofonu } from './ster-mikrofonu';
import { utworzSterModelu } from './ster-modelu';
import { utworzSterNakladu } from './ster-nakladu';
import type { SterPaska } from './ster-nastawy';
import { utworzSterRozszerzen } from './ster-rozszerzen';
import { utworzSterUprawnien } from './ster-uprawnien';
import type { ZrodloZlecenia } from './zrodlo-zlecenia';

// Pasek zlecenia to powierzchnia sterowania kopertą zlecenia pod czatem, ustawiona w jednym rzędzie.

/** Zależności paska są wąskie i wstrzykiwane z zewnątrz, obejmując kanał do rdzenia, nastawy okna, katalog kanałów modelu oraz opcjonalne stery pomocnicze. */
export interface ZaleznosciPaska {
  /** Droga do rdzenia — wyłącznie po to, by założyć rejestr ekspertów. */
  kanal: Kanal;
  /** Odczyt i zapis nastaw okna; jeden na pasek. */
  zrodlo: ZrodloZlecenia;
  /** Wspólny katalog kanałów modelu — jeden na klienta. */
  rejestrKanalow: RejestrKanalow;
  /** Gotowy przełącznik urządzenia; pominięty = paska bez tego steru. */
  sterowanieSrodowiska?: HTMLElement;
  /** Otwiera kolumnę sterowania okna — stopka steru katalogów. */
  otworzSterowanie(): void;
  /** Dokłada rozpoznany tekst do pola wypowiedzi; pominięcie znaczy pasek bez mikrofonu. */
  wstawTekst?(tekst: string): void;
  /** Otwiera katalog rozszerzeń; pominięta = menu „+" bez stopki. */
  otworzKatalogRozszerzen?(): void;
}

export interface PasekZlecenia {
  /** Element stawiany w rzędzie akcji pod polem wypowiedzi. */
  element: HTMLElement;
  /** Przerysowuje wszystkie stery z bieżącej migawki. */
  odswiez(): void;
  /** Zdejmuje nasłuchy sterów i zwalnia mikrofon, by subskrypcje nie przeżyły zejścia okna. */
  rozlacz(): void;
}

export function utworzPasekZlecenia(zaleznosci: ZaleznosciPaska): PasekZlecenia {
  const { kanal, zrodlo, rejestrKanalow } = zaleznosci;

  const element = document.createElement('div');
  element.className = 'dc-pasek-zlecenia';
  // Rola group ma znaczenie: czytnik ekranu zapowiada rząd uchwytów, zanim przeczyta same wartości.
  element.setAttribute('role', 'group');
  element.setAttribute('aria-label', 'Koperta zlecenia');

  const model = utworzSterModelu(
    {
      migawka: () => {
        const okno = zrodlo.migawka().okno;
        return okno.agentId === undefined
          ? { modelChannelId: okno.modelChannelId }
          : { modelChannelId: okno.modelChannelId, agentId: okno.agentId };
      },
      zastosuj: (nazwa, zmiana) => zrodlo.zastosuj(nazwa, zmiana),
    },
    rejestrKanalow,
    utworzRejestrAgentow(kanal),
  );

  const naklad = utworzSterNakladu({
    migawka: () => ({ naklad: zrodlo.migawka().ustawienia.nakladRozumowania }),
    zapisz: (nazwa, klucz, wartosc) => zrodlo.zapisz(nazwa, klucz, wartosc),
  });

  const uprawnienia = utworzSterUprawnien({
    migawka: () => ({ tryb: zrodlo.migawka().okno.permissionMode }),
    zastosuj: (nazwa, zmiana) => zrodlo.zastosuj(nazwa, zmiana),
  });

  const katalogi = utworzSterKatalogow({
    migawka: () => ({ katalogi: zrodlo.migawka().okno.workingDirs }),
    zastosuj: (nazwa, zmiana) => zrodlo.zastosuj(nazwa, zmiana),
    otworzSterowanie: () => zaleznosci.otworzSterowanie(),
  });

  const rozszerzenia = utworzSterRozszerzen({
    zrodlo: utworzZrodloRozszerzen(kanal),
    // Katalog rozszerzeń rozgłasza zmiany zdarzeniem rdzenia, pasek nie odpytuje go ani nie zgaduje.
    naZmiane: (sluchacz) => kanal.naZdarzenie(EventType.ExtensionChanged, () => sluchacz()),
    ...(zaleznosci.otworzKatalogRozszerzen === undefined
      ? {}
      : { otworzKatalog: zaleznosci.otworzKatalogRozszerzen }),
  });

  // Kolejność sterów niesie treść: gdzie, na czym, czym, jak pytać o zgodę, a na końcu jak podać treść.
  if (zaleznosci.sterowanieSrodowiska !== undefined) {
    element.append(zaleznosci.sterowanieSrodowiska);
  }
  element.append(
    katalogi.element,
    model.element,
    naklad.element,
    uprawnienia.element,
    rozszerzenia.element,
  );

  // Gniazdo mikrofonu stoi od razu, ster wchodzi później — bez miejsca wskoczyłby w złe miejsce rzędu.
  const gniazdoMikrofonu = document.createElement('span');
  gniazdoMikrofonu.className = 'dc-pasek-zlecenia__gniazdo';
  element.append(gniazdoMikrofonu);

  /** Silnik dyktowania zakładany tylko wtedy, gdy wynik ma dokąd trafić. */
  let dyktowanie: Dyktowanie | null = null;
  let mikrofon: (SterPaska & { rozlacz(): void }) | null = null;

  const wstawTekst = zaleznosci.wstawTekst;
  if (wstawTekst !== undefined) {
    dyktowanie = utworzDyktowanie(kanal);
    void utworzSterMikrofonu({ dyktowanie, wstawTekst }).then((ster) => {
      // Wartość `null` oznacza brak silnika mowy — pasek zostaje krótszy zamiast pokazywać ikonę.
      if (ster === null) return;
      mikrofon = ster;
      gniazdoMikrofonu.append(ster.element);
    });
  }

  /** Przerysowanie idzie z jednej subskrypcji migawki; rozszerzenia i mikrofon budzą się osobno. */
  function odswiez(): void {
    for (const ster of [katalogi, model, naklad, uprawnienia]) ster.odswiez();
  }

  // Subskrypcję migawki zdejmuje `zrodlo-zlecenia.ts` — pasek zdejmuje wyłącznie to, co sam założył.
  zrodlo.naZmiane(odswiez);

  return {
    element,
    odswiez,

    rozlacz() {
      rozszerzenia.rozlacz();
      mikrofon?.rozlacz();
      // Fasada dyktowania zdejmowana jest zawsze, nawet bez stera — nasłuch zakłada się przy jej utworzeniu.
      dyktowanie?.rozlacz();
    },
  };
}
