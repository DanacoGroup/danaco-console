import type {
  StudioAgentConflict,
  StudioAgentSlot,
} from '../../../../shared/contract';
import { opisWykonawcy, przycisk } from './zadania-wykaz';

/**
 * Praca kilku wykonawców nad jednym dokumentem — strona okna pętli.
 *
 * Dwa wykazy, obie prawdy o tym samym dokumencie:
 *
 * 1. **Obsada** (`studio.agents.slots.list`) — który wykonawca zajął który
 *    fragment, w jakim jest stanie i kiedy jego zajęcie wygasa. Operator ma
 *    widzieć, kto pracuje nad czym, zanim zacznie pisać w tym samym akapicie.
 *
 * 2. **Spięcia** (`studio.agents.conflicts.list`) — co się stało, gdy dwóch
 *    sięgnęło po ten sam fragment: czyja zmiana weszła, czyja została odłożona
 *    i wedle jakiej nastawy.
 *
 * ── Odłożone brzmienie NIE PRZEPADA ─────────────────────────────────────────
 * `StudioAgentConflict.deferredText` niesie brzmienie, które nie weszło. To jest
 * sedno tego widoku, nie ozdoba: praca wykonawcy odłożona bez pokazania jej
 * Operatorowi byłaby pracą wyrzuconą po cichu. Brzmienie stoi więc wprost do
 * przeczytania, wraz z czynnością „Przyjmij brzmienie", która wnosi je do
 * dokumentu decyzją Operatora. Gdzie rdzeń odłożył brzmienie jako propozycję
 * albo zmianę śledzoną, widok odsyła do niej po identyfikatorze — dwóch kopii
 * tego samego brzmienia nie zakłada.
 */

/** Czynności obsady sięgające poza ten widok. */
export interface CzynnosciObsady {
  /** Wnosi odłożone brzmienie do dokumentu decyzją Operatora. */
  przyjmijBrzmienie(spiecie: StudioAgentConflict): void;
  /** Otwiera fragment, o który poszło spięcie, w oknie pracy z dokumentem. */
  pokazFragment(od: number, do_: number): void;
}

export interface WidokObsady {
  element: HTMLElement;
  odswiez(obsada: readonly StudioAgentSlot[], spiecia: readonly StudioAgentConflict[]): void;
}

/** Nazwa stanu wykonawcy dla Operatora. */
const NAZWA_STANU_WYKONAWCY: Record<string, string> = {
  idle: 'bez zajętego fragmentu',
  working: 'pracuje',
  waiting: 'czeka na fragment',
  conflicted: 'zmiana odłożona po spięciu',
  stopped: 'odstawiony',
};

/** Nazwa nastawy, wedle której spięcie rozstrzygnięto. */
const NAZWA_NASTAWY: Record<string, string> = {
  refuse: 'odmowa drugiemu',
  queue: 'odłożenie zmiany drugiego',
  fragmentLock: 'blokada fragmentu na czas pracy',
};

export function utworzWidokObsady(czynnosci: CzynnosciObsady): WidokObsady {
  const naglowekObsady = document.createElement('h4');
  naglowekObsady.className = 'petla-podtytul';
  naglowekObsady.textContent = 'Kto pracuje nad którym fragmentem';

  const listaObsady = document.createElement('ul');
  listaObsady.className = 'petla-obsada';

  const pustkaObsady = document.createElement('p');
  pustkaObsady.className = 'dn-tekst-3';

  const naglowekSpiec = document.createElement('h4');
  naglowekSpiec.className = 'petla-podtytul';
  naglowekSpiec.textContent = 'Spięcia o ten sam fragment';

  const listaSpiec = document.createElement('ul');
  listaSpiec.className = 'petla-spiecia';

  const pustkaSpiec = document.createElement('p');
  pustkaSpiec.className = 'dn-tekst-3';

  const element = document.createElement('section');
  element.className = 'petla-agenci';
  element.append(
    naglowekObsady,
    pustkaObsady,
    listaObsady,
    naglowekSpiec,
    pustkaSpiec,
    listaSpiec,
  );

  function wierszObsady(zajecie: StudioAgentSlot): HTMLElement {
    const pozycja = document.createElement('li');
    pozycja.className = 'dn-karta petla-obsada__wiersz';
    pozycja.dataset['stan'] = zajecie.state;

    const kto = document.createElement('span');
    kto.className = 'petla-obsada__kto';
    kto.textContent = opisWykonawcy(zajecie.actor);

    const stan = document.createElement('span');
    stan.className = 'dn-plakietka';
    stan.textContent = NAZWA_STANU_WYKONAWCY[zajecie.state] ?? zajecie.state;

    const naglowek = document.createElement('div');
    naglowek.className = 'petla-obsada__naglowek';
    naglowek.append(kto, stan);
    pozycja.append(naglowek);

    if (zajecie.rangeStart !== undefined && zajecie.rangeEnd !== undefined) {
      const zakres = document.createElement('p');
      zakres.className = 'dn-tekst-3';
      zakres.textContent = `Fragment: znaki ${zajecie.rangeStart}–${zajecie.rangeEnd}`;
      pozycja.append(zakres);
      const od = zajecie.rangeStart;
      const do_ = zajecie.rangeEnd;
      pozycja.append(przycisk('Pokaż fragment', () => czynnosci.pokazFragment(od, do_)));
    }
    if (zajecie.taskId !== undefined && zajecie.taskId !== '') {
      const zadanie = document.createElement('p');
      zadanie.className = 'dn-tekst-3';
      zadanie.textContent = `W imieniu zadania: ${zajecie.taskId}`;
      pozycja.append(zadanie);
    }
    // Czas wygaśnięcia zajęcia jest treścią, nie szczegółem technicznym:
    // wykonawca ubity w pół pracy nie ma trzymać fragmentu na zawsze, a Operator
    // ma widzieć, do kiedy fragment jest zajęty.
    if (zajecie.expiresAt !== undefined) {
      const wygasa = document.createElement('p');
      wygasa.className = 'dn-tekst-3';
      wygasa.textContent = `Zajęcie wygasa: ${new Date(zajecie.expiresAt).toLocaleString('pl-PL')}`;
      pozycja.append(wygasa);
    }
    return pozycja;
  }

  function wierszSpiecia(spiecie: StudioAgentConflict): HTMLElement {
    const pozycja = document.createElement('li');
    pozycja.className = 'dn-karta petla-spiecie';
    pozycja.dataset['domkniete'] = spiecie.resolved === true ? 'tak' : 'nie';

    const naglowek = document.createElement('div');
    naglowek.className = 'petla-spiecie__naglowek';

    const zakres = document.createElement('span');
    zakres.className = 'petla-spiecie__zakres';
    zakres.textContent = `Znaki ${spiecie.rangeStart}–${spiecie.rangeEnd}`;

    const nastawa = document.createElement('span');
    nastawa.className = 'dn-plakietka';
    nastawa.textContent = NAZWA_NASTAWY[spiecie.policy] ?? spiecie.policy;

    naglowek.append(zakres, nastawa);
    pozycja.append(naglowek);

    const strony = document.createElement('p');
    strony.className = 'dn-tekst-3';
    strony.textContent =
      `Weszła zmiana: ${opisWykonawcy(spiecie.appliedActor)} · ` +
      `odłożona zmiana: ${opisWykonawcy(spiecie.deferredActor)}`;
    pozycja.append(strony);

    const powod = document.createElement('p');
    powod.className = 'dn-tekst-3 petla-spiecie__powod';
    powod.textContent = spiecie.reason;
    pozycja.append(powod);

    if (spiecie.deferredText !== undefined && spiecie.deferredText !== '') {
      const etykieta = document.createElement('p');
      etykieta.className = 'dn-tekst-3';
      etykieta.textContent = 'Brzmienie odłożone — nie przepadło:';
      const brzmienie = document.createElement('pre');
      brzmienie.className = 'dn-kod petla-spiecie__brzmienie';
      brzmienie.textContent = spiecie.deferredText;
      pozycja.append(etykieta, brzmienie);
      if (spiecie.resolved !== true) {
        pozycja.append(
          przycisk('Przyjmij brzmienie', () => czynnosci.przyjmijBrzmienie(spiecie)),
        );
      }
    } else if (spiecie.deferredChangeId !== undefined || spiecie.deferredMarkupId !== undefined) {
      // Rdzeń odłożył brzmienie jako zmianę śledzoną albo propozycję na
      // marginesie. Widok odsyła do niej po identyfikatorze zamiast wyświetlać
      // drugą kopię tekstu, którego prawdą jest tamten byt.
      const odeslanie = document.createElement('p');
      odeslanie.className = 'dn-tekst-3';
      const kod = spiecie.deferredChangeId ?? spiecie.deferredMarkupId ?? '';
      odeslanie.textContent =
        `Brzmienie odłożone leży jako ${
          spiecie.deferredChangeId !== undefined ? 'zmiana śledzona' : 'propozycja na marginesie'
        } ${kod} — rozstrzygasz je tam, gdzie stoi.`;
      pozycja.append(odeslanie);
    }
    return pozycja;
  }

  return {
    element,

    odswiez(obsada, spiecia) {
      listaObsady.replaceChildren();
      if (obsada.length === 0) {
        pustkaObsady.hidden = false;
        pustkaObsady.textContent =
          'Nad tym dokumentem nie pracuje w tej chwili ani jeden wykonawca. ' +
          'Praca kilku wykonawców naraz jest narzędziem nieuruchamianym na starcie — ' +
          'włącza ją nastawa Operatora, nie samo otwarcie okna.';
        listaObsady.hidden = true;
      } else {
        pustkaObsady.hidden = true;
        listaObsady.hidden = false;
        for (const zajecie of obsada) listaObsady.append(wierszObsady(zajecie));
      }

      listaSpiec.replaceChildren();
      if (spiecia.length === 0) {
        pustkaSpiec.hidden = false;
        pustkaSpiec.textContent =
          'Spięć nie było. Pusty wykaz znaczy tu dokładnie tyle: dwóch wykonawców ' +
          'nie sięgnęło po ten sam fragment — a nie że spięcia są przemilczane.';
        listaSpiec.hidden = true;
        return;
      }
      pustkaSpiec.hidden = true;
      listaSpiec.hidden = false;
      for (const spiecie of spiecia) listaSpiec.append(wierszSpiecia(spiecie));
    },
  };
}
