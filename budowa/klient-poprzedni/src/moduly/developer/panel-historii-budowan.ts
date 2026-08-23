import { opisOdmowyBledu } from '../../komponenty/odmowa';
import { pole, przyciskAkcji, wiersz } from '../../modele/kontrolki-formularza-braki';
import type { StanDevelopera } from './stan-developer';
import { utworzStanTresci } from './stany-okna';
import { cialoZakladki, objasnienieZakladki, pasekZakladki } from './zakladki-okna';
import type { ZrodloWarsztatu } from './zrodlo-warsztatu';

/**
 * Historia i pomiary — trzecia część kolumny monitora modułu Developer.
 *
 * ── Po co osobna zakładka obok Build Output ────────────────────────────────
 * Build Output prowadzi przebieg BIEŻĄCY: log narasta zdarzeniem, a okno
 * pokazuje go na żywo. Historia mówi o przebiegach ZAKOŃCZONYCH i o tym, co po
 * nich zostało — wyniku testów i pokryciu kodu. To są dwa różne pytania i dwa
 * różne czasy: pierwsze „co się teraz dzieje”, drugie „co wyszło wtedy”.
 *
 * ── Log historii bywa przycięty i okno o tym mówi ──────────────────────────
 * W dzienniku przebiegu zostaje OGON logu, bo budowanie dużego projektu ma
 * dziesiątki tysięcy wierszy. Odpowiedź niesie znacznik przycięcia i okno go
 * pokazuje: „ostatnie 500 wierszy” to inne zdanie niż „tyle ich było”, a
 * Operator szukający wiersza z początku budowania ma wiedzieć, że go tu nie ma.
 *
 * ── Pusty pomiar to nie zero ───────────────────────────────────────────────
 * Przebieg, w którym nikt nie uruchamiał testów, nie ma wyników — i okno pisze
 * to wprost, zamiast pokazać „0 z 0 przeszło”. Zero przy zerze wygląda jak
 * powodzenie, a znaczy brak pomiaru.
 */
export interface PanelHistoriiBudowan {
  element: HTMLElement;
  odswiez(): void;
}

export function utworzPanelHistoriiBudowan(
  zrodlo: ZrodloWarsztatu,
  stan: StanDevelopera,
): PanelHistoriiBudowan {
  const tresc = utworzStanTresci();
  const przebieg = pole('Przebieg budowania', 'identyfikator przebiegu z historii');

  const historia = przyciskAkcji('Odczytaj historię przebiegów');
  const log = przyciskAkcji('Pokaż log przebiegu');
  const testy = przyciskAkcji('Pokaż wynik testów');
  const pokrycie = przyciskAkcji('Pokaż pokrycie kodu');

  historia.addEventListener('click', () => void odczytajHistorie());
  log.addEventListener('click', () => void pokazLog());
  testy.addEventListener('click', () => void pokazTesty());
  pokrycie.addEventListener('click', () => void pokazPokrycie());

  async function odczytajHistorie(): Promise<void> {
    tresc.ladowanie('Odczyt historii przebiegów…');
    const wynik = await zrodlo.wykazBudowan({ windowId: stan.okno() });
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(opisOdmowyBledu('Odczyt historii budowań', wynik.blad), wynik.blad);
      return;
    }
    if (wynik.wynik.builds.length === 0) {
      tresc.pusto('Okno nie ma jeszcze ani jednego przebiegu budowania w dzienniku.');
      return;
    }
    // Pierwszy przebieg wykazu wchodzi do pola: kolejne pytania dotyczą
    // najczęściej właśnie ostatniego biegu, a przepisywanie identyfikatora ręką
    // z wykazu do pola jest przepisywaniem, które okno może zrobić za Operatora.
    przebieg.value = wynik.wynik.builds[0]?.id ?? '';
    const miejsce = tresc.tresc();
    miejsce.replaceChildren(
      ...wynik.wynik.builds.map((pozycja) =>
        akapitHistorii(
          `${pozycja.id} — ${pozycja.task} — ${pozycja.status}` +
            `${pozycja.exitCode === undefined ? '' : ` (kod ${String(pozycja.exitCode)})`}`,
        ),
      ),
    );
  }

  async function pokazLog(): Promise<void> {
    const kod = przebieg.value.trim();
    if (kod === '') {
      tresc.blad('Odczyt logu wymaga wskazania przebiegu budowania.');
      return;
    }
    tresc.ladowanie('Odczyt logu…');
    const wynik = await zrodlo.logBudowania({ buildId: kod });
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(opisOdmowyBledu('Odczyt logu przebiegu', wynik.blad), wynik.blad);
      return;
    }
    if (wynik.wynik.lines.length === 0) {
      tresc.pusto(`Przebieg ${kod} nie zostawił ani jednego wiersza logu.`);
      return;
    }
    const miejsce = tresc.tresc();
    miejsce.replaceChildren(
      akapitHistorii(
        wynik.wynik.truncated
          ? 'Log jest przycięty — dziennik zachowuje ogon przebiegu, nie jego całość.'
          : 'Log w całości.',
      ),
      blokHistorii(wynik.wynik.lines.join('\n')),
    );
  }

  async function pokazTesty(): Promise<void> {
    const kod = przebieg.value.trim();
    if (kod === '') {
      tresc.blad('Odczyt wyniku testów wymaga wskazania przebiegu budowania.');
      return;
    }
    tresc.ladowanie('Odczyt wyniku testów…');
    const wynik = await zrodlo.wynikTestow({ buildId: kod });
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(opisOdmowyBledu('Odczyt wyniku testów', wynik.blad), wynik.blad);
      return;
    }
    const pomiar = wynik.wynik;
    if (pomiar.results.length === 0) {
      tresc.pusto(
        `Przebieg ${kod} nie ma wyniku testów — w tym biegu nikt ich nie uruchamiał. To nie jest ` +
          'to samo, co „wszystkie przeszły”.',
      );
      return;
    }
    const miejsce = tresc.tresc();
    miejsce.replaceChildren(
      akapitHistorii(
        `Przeszło ${String(pomiar.passed)} · nie przeszło ${String(pomiar.failed)} · ` +
          `pominięto ${String(pomiar.skipped)}` +
          (pomiar.durationMs === undefined ? '' : ` · czas ${String(pomiar.durationMs)} ms`),
      ),
      ...pomiar.results.map((test) =>
        akapitHistorii(
          `[${test.status}] ${test.name}` +
            `${test.suite === undefined ? '' : ` — ${test.suite}`}` +
            `${test.path === undefined ? '' : ` — ${test.path}`}` +
            `${test.line === undefined ? '' : `:${String(test.line)}`}` +
            `${test.message === undefined ? '' : ` — ${test.message}`}`,
        ),
      ),
    );
  }

  async function pokazPokrycie(): Promise<void> {
    const kod = przebieg.value.trim();
    if (kod === '') {
      tresc.blad('Odczyt pokrycia wymaga wskazania przebiegu budowania.');
      return;
    }
    tresc.ladowanie('Odczyt pokrycia…');
    const wynik = await zrodlo.pokrycie({ buildId: kod });
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(opisOdmowyBledu('Odczyt pokrycia kodu', wynik.blad), wynik.blad);
      return;
    }
    if (wynik.wynik.files.length === 0) {
      tresc.pusto(
        `Przebieg ${kod} nie zostawił pomiaru pokrycia — zadanie budowania nie zbierało go ` +
          '(brak profilu pokrycia w parametrach).',
      );
      return;
    }
    const miejsce = tresc.tresc();
    miejsce.replaceChildren(
      akapitHistorii(`Pokrycie zbiorcze: ${String(wynik.wynik.percent)}%`),
      ...wynik.wynik.files.map((plik) =>
        akapitHistorii(
          `${plik.path} — ${String(plik.percent)}% (${String(plik.covered)} z ${String(plik.statements)})` +
            `${plik.uncoveredLines === undefined || plik.uncoveredLines.length === 0 ? '' : ` · bez pokrycia: ${plik.uncoveredLines.join(', ')}`}`,
        ),
      ),
    );
  }

  const element = cialoZakladki(
    objasnienieZakladki(
      'Historia mówi o przebiegach zakończonych. Wynik testów i pokrycie powstają z rozbioru ' +
        'wyjścia W CHWILI biegu — w dzienniku zostaje sam ogon logu, więc pomiar odczytany ' +
        'godzinę później nie miałby z czego powstać.',
    ),
    wiersz('Przebieg budowania', przebieg, {
      klasa: 'dn-pole',
      objasnienie: 'Wypełnia się samo ostatnim przebiegiem po odczycie historii.',
    }),
    pasekZakladki(historia, log, testy, pokrycie),
    tresc.element,
  );

  return { element, odswiez: () => void odczytajHistorie() };
}

/** akapitHistorii składa jeden wiersz treści panelu. */
function akapitHistorii(zdanie: string): HTMLElement {
  const element = document.createElement('p');
  element.className = 'mdev-wiersz';
  element.textContent = zdanie;
  return element;
}

/** blokHistorii składa miejsce na treść wielowierszową — log przebiegu. */
function blokHistorii(zawartosc: string): HTMLElement {
  const element = document.createElement('pre');
  element.className = 'mdev-blok';
  element.textContent = zawartosc;
  return element;
}
