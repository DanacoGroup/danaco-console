import type { Agent, AgentVersion } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import { przycisk } from '../../modele/kontrolki-formularza';
import { utworzStanOkna, type StanOkna } from './stan-okna';
import { utworzZrodloWersjiEksperta, zdanieOWersji } from './zrodlo-wersji-eksperta';
import { utworzZrodloZakresuEksperta, type ZrodloZakresuEksperta } from './zrodlo-zakresu-eksperta';
import type { ZrodloWersjiEksperta } from './zrodlo-wersji-eksperta';
import type { Kanal } from '../../protokol/kanal';

/**
 * Panel „Historia wersji” Agent Buildera pokazuje przegląd wersji tożsamości
 * eksperta zapisanych trwale w rdzeniu oraz umożliwia ich przywrócenie.
 */
export interface PanelHistorii {
  element: HTMLElement;
  /** Nanosi eksperta czynnego i zleca odczyt jego historii. */
  ustaw(ekspert: Agent | null): void;
}

/** Zależności panelu przekazywane z zewnątrz: po przywróceniu wersji biblioteka musi zobaczyć nową wersję eksperta. */
export interface OpcjeHistorii {
  /** Wywoływane po przywróceniu wersji — moduł odświeża wtedy bibliotekę. */
  naPrzywroceniu(ekspert: Agent): void;
}

export function utworzPanelHistorii(kanal: Kanal, opcje: OpcjeHistorii): PanelHistorii {
  const zrodlo: ZrodloWersjiEksperta = utworzZrodloWersjiEksperta(kanal);
  const zakres: ZrodloZakresuEksperta = utworzZrodloZakresuEksperta(kanal);
  const okno: StanOkna = utworzStanOkna();

  const tytul = document.createElement('h4');
  tytul.className = 'da-panel__tytul';
  tytul.textContent = 'Historia wersji';

  const nota = document.createElement('p');
  nota.className = 'dn-pole-opis';
  nota.textContent =
    'Wykaz trwały, odczytany z rdzenia — obejmuje także wersje zapisane przed ' +
    'uruchomieniem klienta i na innym urządzeniu konta. Przywrócenie wstawia wskazaną ' +
    'wersję JAKO KOLEJNĄ: historia nie zostaje skrócona.';

  const lista = document.createElement('ul');
  lista.className = 'da-historia';
  okno.tresc.append(lista);

  const odpowiedz = document.createElement('p');
  odpowiedz.className = 'da-odpowiedz';
  odpowiedz.hidden = true;

  const element = document.createElement('section');
  element.className = 'da-panel da-historia__panel';
  element.append(tytul, nota, okno.element, odpowiedz);

  /** Ekspert, którego historia stoi w panelu — po nim poznajemy zmianę wyboru. */
  let pokazany = '';

  function powiedz(tresc: string, powodzenie: boolean): void {
    odpowiedz.textContent = tresc;
    odpowiedz.hidden = tresc === '';
    odpowiedz.dataset['powodzenie'] = String(powodzenie);
  }

  /** Odczyt historii jednego eksperta; odpowiedź przedawniona nie nadpisuje odpowiedzi świeższej. */
  async function wczytaj(idEksperta: string): Promise<void> {
    okno.ladowanie('Odczyt historii wersji w toku…');
    lista.replaceChildren();
    const wynik = await zrodlo.wersje(idEksperta);
    if (pokazany !== idEksperta) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      okno.blad(opisOdmowy('Odczyt historii wersji', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    pokaz(idEksperta, wynik.wynik.versions);
  }

  function pokaz(idEksperta: string, wersje: readonly AgentVersion[]): void {
    if (wersje.length === 0) {
      lista.replaceChildren();
      okno.puste('Rdzeń nie ma jeszcze ani jednej wersji tożsamości tego eksperta.');
      return;
    }
    const [najnowsza] = wersje;
    lista.replaceChildren(
      ...wersje.map((wersja) =>
        wierszWersji(
          wersja,
          wersja.id === najnowsza?.id,
          () => void przywroc(idEksperta, wersja),
          (miejsce) => void pokazTresc(idEksperta, wersja, miejsce),
        ),
      ),
    );
    okno.gotowe();
  }

  async function przywroc(idEksperta: string, wersja: AgentVersion): Promise<void> {
    powiedz(`Przywracanie wersji „${wersja.label ?? wersja.id}"…`, true);
    const wynik = await zrodlo.przywrocWersje(idEksperta, wersja.id);
    if (!wynik.udany || wynik.wynik === undefined) {
      powiedz(opisOdmowy('Przywrócenie wersji', wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    powiedz(
      `Wersja „${wersja.label ?? wersja.id}" wróciła jako wersja ` +
        `${wynik.wynik.agent.version ?? '?'} — historia nie została skrócona.`,
      true,
    );
    opcje.naPrzywroceniu(wynik.wynik.agent);
    await wczytaj(idEksperta);
  }

  /** Podgląd tożsamości utrwalonej w wersji rozwija się w miejscu, pod wierszem, i nie zastępuje wykazu. */
  async function pokazTresc(
    idEksperta: string,
    wersja: AgentVersion,
    miejsce: HTMLElement,
  ): Promise<void> {
    miejsce.hidden = false;
    miejsce.textContent = 'Odczyt treści wersji…';
    const wynik = await zakres.wersja(idEksperta, wersja.id);
    if (!wynik.udany || wynik.wynik === undefined) {
      miejsce.textContent = opisOdmowy(
        'Odczyt treści wersji',
        wynik.blad?.code,
        wynik.blad?.message,
      );
      return;
    }
    const migawka = wynik.wynik.snapshot;
    const wiersze = [`nazwa: ${migawka.name}`];
    if (migawka.displayName !== undefined && migawka.displayName !== '') {
      wiersze.push(`imię własne: ${migawka.displayName}`);
    }
    if (migawka.description !== undefined && migawka.description !== '') {
      wiersze.push(`przeznaczenie: ${migawka.description}`);
    }
    wiersze.push(`widoczność: ${migawka.visibility}`);
    wiersze.push(
      `poziomy pamięci: ${
        migawka.memoryLevels.length === 0 ? '— pamięć wyłączona' : migawka.memoryLevels.join(', ')
      }`,
    );
    if (migawka.systemPrompt !== undefined && migawka.systemPrompt !== '') {
      wiersze.push('instrukcje systemowe:', migawka.systemPrompt);
    }
    miejsce.textContent = wiersze.join('\n');
  }

  return {
    element,

    ustaw(ekspert) {
      const kod = ekspert?.id ?? '';
      if (kod === pokazany) return;
      pokazany = kod;
      powiedz('', true);
      if (ekspert === null) {
        lista.replaceChildren();
        okno.puste('Wybierz eksperta w bibliotece, aby zobaczyć jego historię wersji.');
        return;
      }
      void wczytaj(kod);
    },
  };
}

/**
 * Jeden wiersz wykazu wersji. Wersja najnowsza nie dostaje przycisku
 * przywrócenia, ponieważ przywrócenie jej do samej siebie wydłużyłoby
 * historię bez treści.
 */
function wierszWersji(
  wersja: AgentVersion,
  najnowsza: boolean,
  naPrzywrocenie: () => void,
  naPodglad: (miejsce: HTMLElement) => void,
): HTMLElement {
  const opis = document.createElement('span');
  opis.className = 'da-historia__opis';
  opis.textContent = zdanieOWersji(wersja);

  const znacznik = document.createElement('span');
  znacznik.className = 'dn-plakietka da-historia__znacznik';
  znacznik.textContent = najnowsza ? 'bieżąca' : 'archiwalna';

  const element = document.createElement('li');
  element.className = 'da-historia__wiersz';
  element.dataset['wersja'] = wersja.id;
  element.dataset['biezaca'] = String(najnowsza);
  const tresc = document.createElement('pre');
  tresc.className = 'dn-kod da-historia__tresc';
  tresc.hidden = true;

  const podglad = przycisk('Pokaż treść wersji', 'dn-btn dn-btn--sm dn-btn--zarys');
  podglad.addEventListener('click', () => naPodglad(tresc));

  element.append(opis, znacznik, podglad);
  if (!najnowsza) {
    const kontrolka = przycisk('Przywróć tę wersję', 'dn-btn dn-btn--sm dn-btn--zarys');
    kontrolka.addEventListener('click', naPrzywrocenie);
    element.append(kontrolka);
  }
  element.append(tresc);
  return element;
}
