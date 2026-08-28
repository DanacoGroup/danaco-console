/**
 * Rama aplikacji — montaż. Składa strefy z prototypu — szynę nawigacji,
 * belkę tytułową, pasek narzędzi, panel Sesje i Projekty, pasmo kart, obszar
 * roboczy i pas stanu — i wstawia je w miejsce wskazane w dokumencie. Ta sama
 * zasada co w drodze wejścia: montaż jest jedynym miejscem, które dotyka
 * dokumentu.
 */

import type { Environment, Module, Session } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { zamontujOknoStudio, type OknoStudio } from '../moduly/studio/montaz.ts';
import { el, tekst } from './narzedzia.ts';
import { belka } from './skladniki/belka.ts';
import { szyna } from './skladniki/szyna.ts';
import { pasekNarzedzi } from './skladniki/pasek-narzedzi.ts';
import { panelSesje } from './skladniki/panel-sesje.ts';
import { pasmoKart } from './skladniki/pasmo-kart.ts';
import { stan as pasStanu } from './skladniki/stan.ts';

export interface NastawyRamy {
  /** Miejsce w dokumencie, w które rama się wstawia. */
  miejsce: HTMLElement;
  /** Środowisko, do którego przebieg wszedł. */
  srodowisko: Environment;
  /** Moduły tego środowiska, w kolejności odebranej od rdzenia. */
  moduly: Module[];
  /** Karty sesji odtworzone przez rdzeń — zasilają panel sesji, pasmo kart i pas stanu. */
  sesje: Session[];
  /** Kanał, którym okno modułu woła komendy rdzenia; brak — rama montuje się bez okna modułu. */
  kanal?: Kanal;
}

/** Kod modułu Studio w wykazie rdzenia — jedyny moduł z oknem roboczym dziś zmontowanym. */
const KOD_MODULU_STUDIO = 'studio';

/** Identyfikator obszaru roboczego — cel `aria-controls` kart pasma. */
const ID_OBSZARU_ROBOCZEGO = 'dn-obszar-glowna';

/** Zapytanie o preferencję systemową motywu — bez odstępu po dwukropku, składnia CSS przyjmuje oba zapisy. */
const ZAPYTANIE_MOTYW_CIEMNY = '(prefers-color-scheme:dark)';

function motywCiemny(): boolean {
  const jawny = document.documentElement.dataset['theme'];
  if (jawny === 'dark') return true;
  if (jawny === 'light') return false;
  return globalThis.matchMedia?.(ZAPYTANIE_MOTYW_CIEMNY).matches === true;
}

export function zamontujRame(w: NastawyRamy): void {
  const glowna = el('main', {
    id: ID_OBSZARU_ROBOCZEGO,
    klasa: 'dn-obszar-tresc',
    'aria-label': tekst('glowna.etykieta'),
    tekst: tekst('glowna.brakModulu'),
  });
  const { wezel: belkaWezel, tytul: tytulWezel } = belka({ srodowisko: w.srodowisko.name });
  let stanWezel = pasStanu({
    srodowisko: w.srodowisko.name,
    liczbaSesji: w.sesje.length,
    motywCiemny: motywCiemny(),
  });

  const obszarGlowny = el('section', { klasa: 'dn-obszar-panel dn-obszar-panel--glowny', 'aria-label': tekst('glowna.etykieta') }, [
    pasmoKart({ sesje: w.sesje, idTresci: ID_OBSZARU_ROBOCZEGO }),
    glowna,
  ]);

  const powloka = el('div', { klasa: 'sta-powloka' }, [
    el('div', { klasa: 'dn-rama-korpus' }, [
      szyna({ srodowisko: w.srodowisko, moduly: w.moduly }),
      el('div', { klasa: 'dn-rama-prawa' }, [
        belkaWezel,
        pasekNarzedzi(),
        el('div', { klasa: 'dn-obszar' }, [panelSesje({ sesje: w.sesje }), obszarGlowny]),
      ]),
    ]),
    stanWezel,
  ]);

  w.miejsce.replaceChildren(powloka);

  /* Okno modułu czynne dziś wyłącznie dla Studio; zejście na inny moduł je
     zdejmuje, żeby obszar roboczy nie niósł treści modułu już opuszczonego. */
  let oknoModulu: OknoStudio | undefined;

  /* Moduł Studia z wykazu rdzenia. Jego `id` — nie kod — idzie do `window.create`,
     bo okno rdzenia wiąże się z modułem po identyfikatorze. */
  const modulStudia = w.moduly.find((m) => m.code === KOD_MODULU_STUDIO);

  function naKlikniecie(zdarzenie: Event): void {
    const przycisk = (zdarzenie.target as Element | null)?.closest(
      '.dn-szyna-poz--modul',
    ) as HTMLElement | null;
    if (przycisk === null) return;
    for (const inny of w.miejsce.querySelectorAll('.dn-szyna-poz--modul')) {
      inny.removeAttribute('aria-current');
    }
    przycisk.setAttribute('aria-current', 'true');
    const nazwa = przycisk.dataset['modulNazwa'] ?? w.srodowisko.name;
    tytulWezel.textContent = `${tekst('belka.marka')} ${tekst('belka.separator')} ${nazwa}`;

    oknoModulu?.zdejmij();
    oknoModulu = undefined;
    if (przycisk.dataset['modulKod'] === KOD_MODULU_STUDIO && w.kanal !== undefined
      && modulStudia !== undefined) {
      oknoModulu = zamontujOknoStudio({ miejsce: glowna, kanal: w.kanal, sesje: w.sesje, modul: modulStudia });
    } else {
      glowna.textContent = tekst('glowna.brakModulu');
    }
  }
  w.miejsce.addEventListener('click', naKlikniecie);

  /* Pierwsza pozycja szyny startuje bieżąca — jeśli to Studio, okno wchodzi
     od razu, bez czekania na klik Operatora. */
  const pierwszy = w.moduly[0];
  if (pierwszy?.code === KOD_MODULU_STUDIO && w.kanal !== undefined && modulStudia !== undefined) {
    oknoModulu = zamontujOknoStudio({ miejsce: glowna, kanal: w.kanal, sesje: w.sesje, modul: modulStudia });
  }

  const zapytanieMotywu = globalThis.matchMedia?.(ZAPYTANIE_MOTYW_CIEMNY);
  function naZmianeMotywu(): void {
    const nowy = pasStanu({
      srodowisko: w.srodowisko.name,
      liczbaSesji: w.sesje.length,
      motywCiemny: motywCiemny(),
    });
    stanWezel.replaceWith(nowy);
    stanWezel = nowy;
  }
  zapytanieMotywu?.addEventListener('change', naZmianeMotywu);
}
