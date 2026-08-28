import type { DesignPrompt } from '../../../../shared/contract';
import { przyciskBrakuDrogi } from './brak-drogi';
import { BRAKI } from './etykiety-designu';

/**
 * Historia poleceń Prompt Buildera wraz z porównaniem wersji promptu. Wykaz
 * obejmuje wyłącznie prompty wydane w tym oknie i ginie razem z nim, a
 * porównanie wymienia pola kontraktu `DesignPrompt`, którymi wersje się różnią.
 */
export interface HistoriaPromptow {
  element: HTMLElement;
  /** Dopisuje prompt wydany do rdzenia. */
  dopisz(prompt: DesignPrompt): void;
  /** Style użyte w tej sesji — podpowiedź biblioteki stylów. */
  style(): readonly string[];
}

export function utworzHistoriePromptow(
  naPrzywrocenie: (prompt: DesignPrompt) => void,
): HistoriaPromptow {
  const wpisy: DesignPrompt[] = [];

  const tytul = document.createElement('h4');
  tytul.className = 'md-panel__tytul';
  tytul.textContent = 'Historia poleceń i wersje promptu';

  const zasieg = document.createElement('p');
  zasieg.className = 'dn-pole-opis md-panel__opis';
  zasieg.textContent =
    'Wykaz obejmuje prompty wydane w tym oknie od jego otwarcia. Komenda historii promptów ' +
    'jest już w kontrakcie, ale czeka na uchwyt w rdzeniu — do tego czasu poza tym oknem ' +
    'promptów nie ma nigdzie i giną z zamknięciem karty.';

  const wykaz = document.createElement('ol');
  wykaz.className = 'md-historia';

  const roznica = document.createElement('pre');
  roznica.className = 'md-historia__roznica';
  roznica.hidden = true;

  const braki = document.createElement('div');
  braki.className = 'md-braki__rzad';
  braki.append(przyciskBrakuDrogi(BRAKI.szablonPromptu));

  const element = document.createElement('section');
  element.className = 'md-panel md-historia-panel';
  element.append(tytul, zasieg, wykaz, roznica, braki);

  function pokazRoznice(numer: number): void {
    const nowszy = wpisy[numer];
    const starszy = wpisy[numer + 1];
    if (nowszy === undefined || starszy === undefined) {
      roznica.hidden = false;
      roznica.textContent = 'To najstarszy prompt tej sesji — nie ma go z czym porównać.';
      return;
    }
    roznica.hidden = false;
    roznica.textContent = opiszRoznice(starszy, nowszy);
  }

  function odswiez(): void {
    wykaz.replaceChildren(
      ...wpisy.map((prompt, numer) => wierszHistorii(prompt, numer, naPrzywrocenie, pokazRoznice)),
    );
  }

  odswiez();

  return {
    element,

    dopisz(prompt) {
      wpisy.unshift(prompt);
      odswiez();
    },

    style() {
      const zebrane = new Set<string>();
      for (const prompt of wpisy) {
        const styl = (prompt.style ?? '').trim();
        if (styl !== '') zebrane.add(styl);
      }
      return [...zebrane].sort((a, b) => a.localeCompare(b, 'pl'));
    },
  };
}

/**
 * Składa jeden wpis historii: opis promptu wraz z parametrami oraz przyciski
 * przywrócenia promptu do pól okna i porównania go z wersją poprzednią.
 */
function wierszHistorii(
  prompt: DesignPrompt,
  numer: number,
  naPrzywrocenie: (prompt: DesignPrompt) => void,
  naPorownanie: (numer: number) => void,
): HTMLElement {
  const opis = document.createElement('span');
  opis.className = 'md-historia__opis';
  opis.textContent = `${prompt.subject === '' ? '(bez tematu)' : prompt.subject} · ${opiszParametry(prompt)}`;

  const przywroc = document.createElement('button');
  przywroc.type = 'button';
  przywroc.className = 'dn-btn dn-btn--duch dn-btn--sm';
  przywroc.textContent = 'Przywróć do pól';
  przywroc.addEventListener('click', () => naPrzywrocenie(prompt));

  const porownaj = document.createElement('button');
  porownaj.type = 'button';
  porownaj.className = 'dn-btn dn-btn--duch dn-btn--sm';
  porownaj.textContent = 'Porównaj z poprzednim';
  porownaj.addEventListener('click', () => naPorownanie(numer));

  const element = document.createElement('li');
  element.className = 'md-historia__wiersz';
  element.dataset['wersja'] = String(numer);
  element.append(opis, przywroc, porownaj);
  return element;
}

/**
 * Zestawia parametry promptu w jedno zdanie opisu: kreatywność i liczbę
 * wariantów zawsze, a ziarno oraz silnik wtedy, gdy prompt je niesie.
 */
function opiszParametry(prompt: DesignPrompt): string {
  const czesci = [`kreatywność ${prompt.creativity ?? 0}`, `warianty ${prompt.variants ?? 1}`];
  if (prompt.seed !== undefined) czesci.push(`ziarno ${prompt.seed}`);
  if (prompt.engine !== undefined) czesci.push(`silnik ${prompt.engine}`);
  return czesci.join(' · ');
}

/**
 * Wylicza pola, którymi dwa prompty się różnią, zestawiając wartość starszą
 * z nowszą; zgodność pole w pole daje zamiast wykazu osobne zdanie.
 */
function opiszRoznice(starszy: DesignPrompt, nowszy: DesignPrompt): string {
  const klucze = new Set<string>([...Object.keys(starszy), ...Object.keys(nowszy)]);
  const wiersze: string[] = [];
  // Pola promptu mają różne typy, więc porównanie idzie po zapisie tekstowym.
  const przedPola = starszy as unknown as Record<string, unknown>;
  const poPola = nowszy as unknown as Record<string, unknown>;
  for (const klucz of [...klucze].sort()) {
    const przed = String(przedPola[klucz] ?? '');
    const po = String(poPola[klucz] ?? '');
    if (przed !== po) wiersze.push(`${klucz}: „${przed}" → „${po}"`);
  }
  return wiersze.length === 0 ? 'Prompty są identyczne polem w pole.' : wiersze.join('\n');
}
