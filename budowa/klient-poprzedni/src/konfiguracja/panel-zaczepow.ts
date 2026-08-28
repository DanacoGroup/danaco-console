import type { HookBinding, SessionConfigHooks } from '../../../shared/contract';
import {
  pole,
  przelacznik,
  przyciskAkcji,
} from '../modele/kontrolki-formularza-braki';
import type { StanObszarowSesji } from './stan-obszarow-sesji';

/**
 * Panel zaczepów cyklu życia kanału — powierzchnia redagowania obszaru
 * `hooks`. Zaczep bez zdarzenia albo bez polecenia odrzuca rdzeń. Zapis idzie
 * adresem `config.session.set` z obszarem `hooks` pod punktem widzenia okna.
 */
export interface PanelZaczepow {
  element: HTMLElement;
  /** Nanosi stan rdzenia, o ile wykaz nie jest właśnie redagowany. */
  odswiez(): void;
}

/** Punkty cyklu życia obsługiwane przez program modelu — podpowiedź w polu wpisu, nie lista zamknięta wartości. */
const ZNANE_ZDARZENIA =
  'UserPromptSubmit · PreToolUse · PostToolUse · Stop · SessionStart · ' +
  'SessionEnd · Notification · SubagentStop · PreCompact';

/** Jeden redagowany wiersz zaczepu: zdarzenie, zawężenie, polecenie oraz przełącznik stanu czynnego zaczepu. */
interface WierszZaczepu {
  element: HTMLElement;
  odczyt(): HookBinding;
}

export function utworzPanelZaczepow(stan: StanObszarowSesji): PanelZaczepow {
  const wykaz = document.createElement('ul');
  wykaz.className = 'dk-zaczepy__wykaz';
  wykaz.setAttribute('aria-label', 'Zaczepy cyklu życia kanału');

  const wiersze: WierszZaczepu[] = [];
  let redagowane = false;

  const obszarCzynny = przelacznik('Zaczepy są wywoływane');
  obszarCzynny.checked = true;

  const odpowiedz = document.createElement('p');
  odpowiedz.className = 'dk-zaczepy__odpowiedz';
  odpowiedz.hidden = true;

  const dodaj = przyciskAkcji('Dodaj zaczep', 'dn-btn dn-btn--zarys');
  dodaj.addEventListener('click', () => {
    redagowane = true;
    dolozWiersz({ event: '', command: '', enabled: true });
  });

  const zapisz = przyciskAkcji('Zapisz zaczepy na tym poziomie', 'dn-btn dn-btn--atrament');
  zapisz.addEventListener('click', () => void wykonajZapis());

  const element = document.createElement('section');
  element.className = 'dk-zaczepy';
  element.append(zlozNaglowek(obszarCzynny), wykaz, zlozStopke(dodaj, zapisz, odpowiedz));

  function dolozWiersz(zaczep: HookBinding): void {
    const wiersz = zlozWierszZaczepu(zaczep, () => {
      redagowane = true;
      const indeks = wiersze.indexOf(wiersz);
      if (indeks >= 0) wiersze.splice(indeks, 1);
      wiersz.element.remove();
    });
    wiersz.element.addEventListener('input', () => {
      redagowane = true;
    });
    wiersze.push(wiersz);
    wykaz.append(wiersz.element);
  }

  function pokazZeStanu(): void {
    const zaczepy = stan.obowiazujaca()?.config.hooks;
    wiersze.length = 0;
    wykaz.replaceChildren();
    obszarCzynny.checked = zaczepy?.enabled !== false;
    for (const zaczep of zaczepy?.hooks ?? []) dolozWiersz(zaczep);
  }

  async function wykonajZapis(): Promise<void> {
    const tresc: SessionConfigHooks = {
      enabled: obszarCzynny.checked,
      hooks: wiersze.map((wiersz) => wiersz.odczyt()),
    };
    const wynik = await stan.utrwalZaczepy(tresc);
    odpowiedz.hidden = false;
    if (!wynik.udany) {
      odpowiedz.textContent = wynik.blad?.message ?? 'Rdzeń odmówił bez podania powodu.';
      odpowiedz.dataset['powodzenie'] = 'false';
      return;
    }
    redagowane = false;
    odpowiedz.textContent =
      `Zaczepy zapisane (${tresc.hooks?.length ?? 0}). ` +
      'Obowiązują od następnej tury okna objętego tym poziomem.';
    odpowiedz.dataset['powodzenie'] = 'true';
    void stan.odswiez();
  }

  function odswiez(): void {
    // Wykaz w redakcji nie jest nadpisywany odczytem — odświeżenie skasowałoby
    // niedokończony wiersz.
    if (redagowane) return;
    pokazZeStanu();
  }

  odswiez();
  return { element, odswiez };
}

function zlozNaglowek(obszarCzynny: HTMLInputElement): HTMLElement {
  const tytul = document.createElement('h3');
  tytul.className = 'dk-zaczepy__tytul';
  tytul.textContent = 'Zaczepy cyklu życia kanału';

  const opis = document.createElement('p');
  opis.className = 'dn-pole-opis';
  opis.textContent =
    'Polecenia wykonywane w punktach cyklu życia procesu modelu. Konfiguracja ' +
    'jedzie do programu sekcją hooks pliku ustawień, a zdarzenia zaczepów ' +
    'wracają strumieniem do dziennika zdarzeń i diagnostyki. ' +
    `Znane punkty: ${ZNANE_ZDARZENIA}.`;

  const element = document.createElement('header');
  element.className = 'dk-zaczepy__naglowek';
  element.append(tytul, opis, obszarCzynny.closest('label') ?? obszarCzynny);
  return element;
}

function zlozStopke(...czesci: HTMLElement[]): HTMLElement {
  const element = document.createElement('div');
  element.className = 'dk-zaczepy__stopka';
  element.append(...czesci);
  return element;
}

/** Wiersz redagowania jednego zaczepu: zdarzenie, zawężenie, polecenie, czynność oraz przycisk usunięcia wiersza. */
function zlozWierszZaczepu(zaczep: HookBinding, usunMnie: () => void): WierszZaczepu {
  const zdarzenie = pole('Punkt cyklu życia', 'np. UserPromptSubmit');
  zdarzenie.value = zaczep.event;

  const zawezenie = pole('Zawężenie (matcher)', 'puste znaczy każde wystąpienie');
  zawezenie.value = zaczep.matcher ?? '';

  const polecenie = pole('Polecenie zaczepu', 'polecenie powłoki');
  polecenie.value = zaczep.command;

  const czynny = przelacznik('Zaczep czynny');
  czynny.checked = zaczep.enabled;

  const usun = przyciskAkcji('Usuń', 'dn-btn dn-btn--zarys dk-zaczepy__usun');
  usun.addEventListener('click', usunMnie);

  const element = document.createElement('li');
  element.className = 'dk-zaczepy__wiersz';
  element.append(
    otul(zdarzenie, 'Punkt cyklu życia'),
    otul(zawezenie, 'Zawężenie'),
    otul(polecenie, 'Polecenie'),
    czynny.closest('label') ?? czynny,
    usun,
  );

  return {
    element,
    odczyt: () => {
      const wpis: HookBinding = {
        event: zdarzenie.value.trim(),
        command: polecenie.value.trim(),
        enabled: czynny.checked,
      };
      const matcher = zawezenie.value.trim();
      if (matcher !== '') wpis.matcher = matcher;
      return wpis;
    },
  };
}

/** Pole z podpisem — podpis mówi, czym pole jest, zanim użytkownik zdąży wpisać do niego treść ustawienia. */
function otul(kontrolka: HTMLInputElement, podpis: string): HTMLElement {
  const etykieta = document.createElement('label');
  etykieta.className = 'dk-zaczepy__pole';
  const napis = document.createElement('span');
  napis.className = 'dn-pole-opis';
  napis.textContent = podpis;
  etykieta.append(napis, kontrolka);
  return etykieta;
}
