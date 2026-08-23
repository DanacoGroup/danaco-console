import { utworzNaglowekOkna } from '../../komponenty/naglowek-okna';
import { utworzZakladkiSekcji } from '../../modele/zakladki-sekcji';
import { utworzPanelFaktow, type PanelFaktow } from './panel-faktow';
import { utworzPanelKontekstow } from './panel-kontekstow';
import { utworzPanelWiedzy } from './panel-wiedzy';
import { utworzPanelZestawow, type PanelZestawow } from './panel-zestawow';
import type { StanAssistant } from './stan-assistant';
import type { ZrodloPamieci } from './zrodlo-pamieci';

/** Kod okna operacyjnego modułu; rdzeń nie ma go dziś w katalogu okien. */
export const KOD_OKNA = 'memory-context-manager';

/**
 * Memory & Context Manager — okno zarządcy pamięci i kontekstów modułu
 * Assistant.
 *
 * Okno zamyka obszar pamięci modułu i realizuje zasadę jawności: pamięć
 * asystenta jest w całości widoczna, edytowalna i usuwalna przez Operatora.
 * Cztery zakładki odpowiadają czterem obszarom opracowania modułu — fakty,
 * pamięć semantyczna, konteksty i baza wiedzy.
 *
 * Plik odpowiada wyłącznie za skład okna. Wywołania kontraktu mieszkają
 * w `zrodlo-pamieci.ts`, a każda zakładka ma własny plik obszaru: ustalenia
 * w `panel-faktow.ts`, wskaźnik znaczenia w `panel-wiedzy.ts`, poziomy pamięci
 * karty sesji w `panel-kontekstow.ts`.
 *
 * Fazy odczytu nie ma na poziomie okna, tylko w obszarach. Cztery zakładki
 * czytają cztery różne byty i wołają je w różnych chwilach — jedna faza dla
 * całego okna kazałaby odmowie odczytu pamięci przesłonić zakładkę wskaźnika,
 * która o tej odmowie nic nie wie.
 *
 * Piąta zakładka — „Zestawy i retencja" — zamyka trzy obszary, które kontrakt
 * niesie w całości: nazwane konteksty pamięci (`memory.context.*`), zasady
 * retencji i wygaszania (`memory.retention.*`) oraz miernik okna kontekstu
 * (`context.usage.get`). Miernik pokazuje liczbę policzoną tokenizatorem
 * rdzenia wraz z nazwą słownika, którym policzono — a gdy pomiar jest
 * niewykonalny, sam powód.
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

