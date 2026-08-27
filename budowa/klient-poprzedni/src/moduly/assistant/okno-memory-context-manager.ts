import { utworzNaglowekOkna } from '../../komponenty/naglowek-okna';
import { utworzZakladkiSekcji } from '../../modele/zakladki-sekcji';
import { utworzPanelFaktow, type PanelFaktow } from './panel-faktow';
import { utworzPanelKontekstow } from './panel-kontekstow';
import { utworzPanelWiedzy } from './panel-wiedzy';
import { utworzPanelZestawow, type PanelZestawow } from './panel-zestawow';
import type { StanAssistant } from './stan-assistant';
import type { ZrodloPamieci } from './zrodlo-pamieci';

/** Kod okna operacyjnego modułu; rdzeń nie ma go dziś w katalogu okien. Wartość trafia do atrybutu danych sekcji okna i służy odszukaniu okna w drzewie dokumentu. */
export const KOD_OKNA = 'memory-context-manager';

/**
 * Okno zarządcy pamięci i kontekstów modułu Assistant. Składa pięć zakładek:
 * fakty, pamięć semantyczna, konteksty, baza wiedzy oraz zestawy i retencja.
 * Odpowiada wyłącznie za skład okna; wywołania kontraktu i fazę odczytu niesie
 * każdy obszar osobno.
 */
export interface OknoMemoryContextManager {
  element: HTMLElement;
  /** Odczyt pamięci widocznej w zasięgu karty sesji. */
  wczytaj(): Promise<void>;
}

export function utworzOknoMemoryContextManager(
  stan: StanAssistant,
  zrodlo: ZrodloPamieci,
): OknoMemoryContextManager {
  const fakty: PanelFaktow = utworzPanelFaktow(stan, zrodlo);
  const wiedza = utworzPanelWiedzy(stan, zrodlo);
  const konteksty = utworzPanelKontekstow(stan, zrodlo);
  const zestawy: PanelZestawow = utworzPanelZestawow(stan);

  const zakladki = utworzZakladkiSekcji([
    { kod: 'fakty', nazwa: 'Fakty', element: fakty.element },
    { kod: 'pamiec-semantyczna', nazwa: 'Pamięć semantyczna', element: wiedza.szukanie },
    { kod: 'konteksty', nazwa: 'Konteksty', element: konteksty.element },
    { kod: 'baza-wiedzy', nazwa: 'Baza wiedzy', element: wiedza.wskaznik },
    { kod: 'zestawy', nazwa: 'Zestawy i retencja', element: zestawy.element },
  ]);

  const element = document.createElement('section');
  element.className = 'ma-okno ma-okno--zarzadca';
  element.dataset['okno'] = KOD_OKNA;
  element.append(
    utworzNaglowekOkna({
      tytul: 'Memory & Context Manager',
      rola: 'zarządca · pamięć, konteksty i baza wiedzy asystenta',
    }),
    zakladki.element,
  );

  return {
    element,
    wczytaj: async () => {
      await Promise.all([fakty.wczytaj(), zestawy.wczytaj()]);
    },
  };
}

