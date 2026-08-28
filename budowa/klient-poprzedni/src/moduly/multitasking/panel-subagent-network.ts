import {
  ConfigScope,
  SessionConfigArea,
  type SessionConfigTools,
} from '../../../../shared/contract';
import { przyciskAkcji as przycisk, pozycjaWykazu, wykaz } from '../../modele/kontrolki-formularza';
import type { WykazKomendRdzenia } from './braki-kontraktu';
import { utworzFormularzPowolania } from './formularz-powolania';
import { PRZEDROSTEK } from './kontrolki';
import { utworzTrescPodagentow, type TrescPodagentow } from './panel-podagentow';
import type { ZrodloOkien } from './zrodlo-okien';
import type { ZrodloPodagentow } from './zrodlo-podagentow';

// Panel Subagent Network okna wykonawcy pokazuje definicje podagentów sesji i żywy wykaz do limitu.

/** Górny limit podagentów jednego wykonawcy, ustalony wprost przez ten plik klienta, niezależnie od kontraktu. */
export const LIMIT_PODAGENTOW = 15;

export interface PanelPodagentow {
  /** Odpina nasłuch żywego wykazu; woła je zamknięcie okna wykonawcy. */
  rozlacz(): void;
  element: HTMLElement;
  /** Odczytuje definicje podagentów z konfiguracji sesji oraz żywy wykaz. */
  odswiez(): void;
  // Żywy wykaz zależy od okna wykonawcy, definicje od sesji i przy zmianie obsady się nie zmieniają.
  odswiezZywych(): void;
}

export interface OpcjePanelu {
  zrodlo: ZrodloOkien;
  /** Wykaz komend rdzenia — stąd bierze się powód nieczynnych kontrolek. */
  komendy: WykazKomendRdzenia;
  /** Sesja, z której czytamy obszar `tools`. */
  sesja(): string;
  potwierdz(zdanie: string, udane: boolean): void;
  // Żywy wykaz powołanych jest sekcją zamkniętą; pola opcjonalne, bo wytwórnia bywa wołana bez źródeł.
  zywi?: {
    podagenci: ZrodloPodagentow;
    /** Okno tego wykonawcy; puste, dopóki obsada go nie założyła. */
    okno(): string;
  };
}

export function utworzPanelPodagentow(opcje: OpcjePanelu): PanelPodagentow {
  const lista = wykaz('Podagenci dostępni wykonawcy', 'dm-wykaz');

  const odczyt = przycisk('Odczytaj definicje', 'dn-btn dn-btn--sm');
  odczyt.addEventListener('click', () => {
    void odswiez();
  });

  const pasek = document.createElement('div');
  pasek.className = 'dm-podagenci__pasek';
  pasek.append(odczyt);

  const element = document.createElement('div');
  element.className = 'dm-podagenci';
  element.setAttribute('aria-label', 'Panel Subagent Network');
  element.append(pasek, lista);

  // Żywy wykaz powołanych wchodzi pod wykaz definicji, gdy wołający dał mu
  // źródła.
  let zywi: TrescPodagentow | null = null;
  if (opcje.zywi === undefined) {
    // Bez okna wykonawcy nie ma windowId, więc żadnej z trzech komend nie da się wywołać.
    pasek.append(
      opcje.komendy.przyciskNieczynny('Uruchom Subagent Network', 'subagent.spawn', 'subagent.list'),
      opcje.komendy.przyciskNieczynny('Scal wyniki (merge)', 'subagent.result.collect'),
      opcje.komendy.przyciskNieczynny('+ dodaj podagenta', 'subagent.spawn'),
    );
  } else {
    zywi = utworzTrescPodagentow({
      podagenci: opcje.zywi.podagenci,
      okno: opcje.zywi.okno,
      potwierdz: opcje.potwierdz,
    });
    element.append(zywi.element, sekcjaPowolania(opcje.zywi, opcje, () => zywi?.odswiez()));
  }

  // Odczyt udany też melduje, inaczej po naciśnięciu zostawałoby potwierdzenie innej czynności.
  async function odswiez(): Promise<void> {
    const sesja = opcje.sesja();
    if (sesja === '') {
      pokaz([], 'Sesja nieznana — nie ma gdzie szukać definicji podagentów.');
      opcje.potwierdz('Sesja nieznana — odczyt definicji podagentów nie ruszył.', false);
      return;
    }
    const wynik = await opcje.zrodlo.konfiguracjaSesji({
      scope: ConfigScope.Session,
      scopeId: sesja,
      areas: [SessionConfigArea.Tools],
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      pokaz([], `Odczyt konfiguracji sesji odmówiony (kod ${wynik.blad?.code ?? 'brak'}).`);
      opcje.potwierdz(
        `Rdzeń odmówił odczytu definicji podagentów. Powód: ${wynik.blad?.message ?? 'rdzeń nie podał przyczyny'} (kod ${wynik.blad?.code ?? 'brak'}).`,
        false,
      );
      return;
    }
    const nazwy = nazwyPodagentow(wynik.wynik.tools);
    pokaz(nazwy, 'Sesja nie ma zapisanych definicji podagentów.');
    opcje.potwierdz(
      nazwy.length === 0
        ? 'Rdzeń oddał konfigurację sesji bez ani jednej definicji podagenta.'
        : `Rdzeń oddał ${nazwy.length} definicji podagentów (limit na wykonawcę: ${LIMIT_PODAGENTOW}).`,
      true,
    );
  }

  /** Rysuje wykaz albo stan pusty; limit widoczny zawsze. */
  function pokaz(nazwy: readonly string[], zdaniePustki: string): void {
    lista.replaceChildren();
    if (nazwy.length === 0) {
      const pusto = document.createElement('li');
      pusto.className = 'dm-pozycja dm-pozycja--pusta';
      pusto.textContent = zdaniePustki;
      lista.append(pusto);
      return;
    }
    nazwy.slice(0, LIMIT_PODAGENTOW).forEach((nazwa, numer) => {
      const { element: wiersz } = pozycjaWykazu(
        nazwa,
        `karta ${numer + 1} z ${Math.min(nazwy.length, LIMIT_PODAGENTOW)} · limit ${LIMIT_PODAGENTOW}`,
        PRZEDROSTEK,
      );
      wiersz.dataset['podagent'] = nazwa;
      lista.append(wiersz);
    });
    if (nazwy.length > LIMIT_PODAGENTOW) {
      const nadmiar = document.createElement('li');
      nadmiar.className = 'dm-pozycja dm-pozycja--pusta';
      nadmiar.textContent = `Sesja definiuje ${nazwy.length} podagentów — ponad limit ${LIMIT_PODAGENTOW} na wykonawcę.`;
      lista.append(nadmiar);
    }
  }

  pokaz([], 'Definicje podagentów nieodczytane — użyj „Odczytaj definicje".');

  return {
    element,
    odswiez: () => {
      void odswiez();
      zywi?.odswiez();
    },
    odswiezZywych: () => {
      zywi?.odswiez();
    },
    rozlacz: () => {
      zywi?.rozlacz?.();
    },
  };
}

/**
 * Sekcja powołania podagenta — formularz `subagent.spawn` domyślnie zamknięty.
 *
 * Sekcja jest zamknięta, bo formularz o czterech polach rozpychałby okno
 * wykonawcy przy każdym otwarciu sceny, i oznaczona nagłówkiem mówiącym wprost,
 * co robi.
 */
function sekcjaPowolania(
  zywi: NonNullable<OpcjePanelu['zywi']>,
  opcje: OpcjePanelu,
  odswiezWykaz: () => void,
): HTMLElement {
  const formularz = utworzFormularzPowolania({
    podagenci: zywi.podagenci,
    okno: zywi.okno,
    potwierdz: opcje.potwierdz,
    poPowolaniu: odswiezWykaz,
  });

  const tytul = document.createElement('summary');
  tytul.className = 'dm-podagenci-zywi__tytul';
  tytul.textContent = 'Dodaj subagenta (subagent.spawn)';

  const element = document.createElement('details');
  element.className = 'dm-podagenci-zywi';
  element.append(tytul, formularz.element);
  return element;
}

/**
 * Nazwy podagentów z pola o kształcie nieokreślonym w kontrakcie: widok nie
 * zgaduje struktury, bierze nazwę, gdy pozycja ją ma, a w przeciwnym razie
 * zapis pozycji surowej.
 */
function nazwyPodagentow(tools?: SessionConfigTools): readonly string[] {
  const definicje = tools?.subagentDefinitions;
  if (!Array.isArray(definicje)) return [];
  return definicje.map((pozycja, numer) => {
    if (typeof pozycja === 'string') return pozycja;
    if (typeof pozycja === 'object' && pozycja !== null && 'name' in pozycja) {
      const nazwa = (pozycja as { name: unknown }).name;
      if (typeof nazwa === 'string' && nazwa !== '') return nazwa;
    }
    return `podagent ${numer + 1}`;
  });
}
