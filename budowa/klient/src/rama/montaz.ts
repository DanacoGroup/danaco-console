/**
 * Rama aplikacji — montaż. Składa trzy strefy z prototypu — szynę nawigacji,
 * belkę tytułową i pas stanu — i wstawia je w miejsce wskazane w dokumencie.
 * Ta sama zasada co w drodze wejścia: montaż jest jedynym miejscem, które
 * dotyka dokumentu.
 */

import type { Environment, Module, Session } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { zamontujOknoStudio, type OknoStudio } from '../moduly/studio/montaz.ts';
import { el, tekst } from './narzedzia.ts';
import { belka } from './skladniki/belka.ts';
import { szyna } from './skladniki/szyna.ts';
import { stan as pasStanu } from './skladniki/stan.ts';

export interface NastawyRamy {
  /** Miejsce w dokumencie, w które rama się wstawia. */
  miejsce: HTMLElement;
  /** Środowisko, do którego przebieg wszedł. */
  srodowisko: Environment;
  /** Moduły tego środowiska, w kolejności odebranej od rdzenia. */
  moduly: Module[];
  /** Karty sesji odtworzone przez rdzeń — ich liczba zasila pas stanu. */
  sesje: Session[];
  /** Kanał, którym okno modułu woła komendy rdzenia; brak — rama montuje się bez okna modułu. */
  kanal?: Kanal;
}

/** Kod modułu Studio w wykazie rdzenia — jedyny moduł z oknem roboczym dziś zmontowanym. */
const KOD_MODULU_STUDIO = 'studio';

/** Zapytanie o preferencję systemową motywu — bez odstępu po dwukropku, składnia CSS przyjmuje oba zapisy. */
const ZAPYTANIE_MOTYW_CIEMNY = '(prefers-color-scheme:dark)';

function motywCiemny(): boolean {
  const jawny = document.documentElement.dataset['theme'];
  if (jawny === 'dark') return true;
  if (jawny === 'light') return false;
  return globalThis.matchMedia?.(ZAPYTANIE_MOTYW_CIEMNY).matches === true;
}

export function zamontujRame(w: NastawyRamy): void {
  const glowna = el('main', { 'aria-label': tekst('glowna.etykieta') });
  const { wezel: belkaWezel, tytul: tytulWezel } = belka({ srodowisko: w.srodowisko.name });
  let stanWezel = pasStanu({
    srodowisko: w.srodowisko.name,
    liczbaSesji: w.sesje.length,
    motywCiemny: motywCiemny(),
  });

  const powloka = el('div', { klasa: 'sta-powloka' }, [
    el('div', { klasa: 'dn-rama-korpus' }, [
      szyna({ moduly: w.moduly }),
      el('div', { klasa: 'dn-rama-prawa' }, [belkaWezel, glowna]),
    ]),
    stanWezel,
  ]);

  w.miejsce.replaceChildren(powloka);

  /* Okno modułu czynne dziś wyłącznie dla Studio; zejście na inny moduł je
     zdejmuje, żeby obszar roboczy nie niósł treści modułu już opuszczonego. */
  let oknoModulu: OknoStudio | undefined;

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
    if (przycisk.dataset['modulKod'] === KOD_MODULU_STUDIO && w.kanal !== undefined) {
      oknoModulu = zamontujOknoStudio({ miejsce: glowna, kanal: w.kanal });
    } else {
      glowna.replaceChildren();
    }
  }
  w.miejsce.addEventListener('click', naKlikniecie);

  /* Pierwsza pozycja szyny startuje bieżąca — jeśli to Studio, okno wchodzi
     od razu, bez czekania na klik Operatora. */
  const pierwszy = w.moduly[0];
  if (pierwszy?.code === KOD_MODULU_STUDIO && w.kanal !== undefined) {
    oknoModulu = zamontujOknoStudio({ miejsce: glowna, kanal: w.kanal });
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
