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
import { panelSesje } from './skladniki/panel-sesje.ts';
import { wypelnijPowloke } from './skladniki/wypelnienie-powloki.ts';
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
  /*
  Konto, którym Operator wszedł. Wskazanie idzie z drogi wejścia, nie z rdzenia:
  `AuthSession` niesie token, termin, urządzenie i metodę — konta nie. Puste
  znaczy wejście drogą, która konta nie podaje; pasek stanu nazywa je wtedy
  nierozpoznanym, zamiast zmyślać.
  */
  kontoOperatora?: string;
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
  const { tytul: tytulWezel } = belka({ srodowisko: w.srodowisko.name });
  let stanWezel = pasStanu({
    srodowisko: w.srodowisko.name,
    liczbaSesji: w.sesje.length,
    motywCiemny: motywCiemny(),
  });

  /* Pasmo kart sesji stoi tylko wtedy, gdy obszar roboczy nie niesie okna
     modułu. Źródło kształtu (`design/05-okna/moduly/studio.html`, w. 346) ma
     w widoku modułu jedno pasmo — to, które okno robocze przynosi ze sobą;
     drugie nad nim czytałoby się jak okno wstawione w okno. */
  const pasmoSesji = pasmoKart({ sesje: w.sesje, idTresci: ID_OBSZARU_ROBOCZEGO });
  const obszarGlowny = el('section', { klasa: 'dn-obszar-panel dn-obszar-panel--glowny', 'aria-label': tekst('glowna.etykieta') }, [
    pasmoSesji,
    glowna,
  ]);

  /*
  Powłoka pochodzi z biblioteki prototypu (`zasoby/powloka.js`), nie z tego
  pliku: szynę nawigacji, belkę tytułową, pas narzędzi i pasek stanu stawia ona,
  wypełniając gniazdo `#dn-powloka-montaz` przy wczytaniu dokumentu.

  Rama wstawia tu wyłącznie obszar modułu — w `.dn-rama-prawa`, pod pasem
  narzędzi, dokładnie tam, gdzie treść okna stoi w źródle kształtu.
  */
  const ramaPrawa = w.miejsce.querySelector('.dn-rama-prawa');
  if (ramaPrawa === null) {
    /* Biblioteka nie zamontowała powłoki — bez niej okno nie ma w co wejść.
       Zdanie idzie do konsoli przeglądarki, nie do Operatora: to usterka
       wydania, nie stan, który on mógłby naprawić. */
    console.error('powłoka nie została zamontowana: brak .dn-rama-prawa w gnieździe ramy');
    return;
  }

  /*
  Okno robocze pochodzi z prototypu: `vite.config.js` wstrzykuje blok
  `main.st-okno-robocze` ze źródła kształtu w szablon `#dn-tresc-okna`, a powłoka
  wkleja go w tym miejscu. Rama go NIE buduje — dokłada tylko to, czego prototyp
  z natury nie niesie, bo pochodzi z rdzenia.

  Obszar własny staje wyłącznie wtedy, gdy prototypowego okna nie ma — inaczej
  Operator zobaczyłby dwa okna robocze jedno pod drugim.
  */
  const oknoZPrototypu = ramaPrawa.querySelector('.st-okno-robocze');
  if (oknoZPrototypu === null) {
    ramaPrawa.append(el('div', { klasa: 'dn-obszar' }, [
      panelSesje({ sesje: w.sesje }),
      obszarGlowny,
    ]));
  }

  /* Powłoka przychodzi z treścią przykładową prototypu — nazwą innego okna
     w belce i zmyślonymi miarami maszyny w pasku stanu. Wypełnienie zastępuje
     je wartościami rdzenia, a pozycje bez źródła zdejmuje. */
  wypelnijPowloke(w.miejsce, {
    tytul: `${tekst('belka.tytul')} — ${w.srodowisko.name}`,
    srodowisko: w.srodowisko.name,
    liczbaSesji: w.sesje.length,
    motywCiemny: motywCiemny(),
    operator: w.kontoOperatora,
  });

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
    pasmoSesji.hidden = oknoModulu !== undefined;
  }
  w.miejsce.addEventListener('click', naKlikniecie);

  /* Pierwsza pozycja szyny startuje bieżąca — jeśli to Studio, okno wchodzi
     od razu, bez czekania na klik Operatora. */
  const pierwszy = w.moduly[0];
  if (pierwszy?.code === KOD_MODULU_STUDIO && w.kanal !== undefined && modulStudia !== undefined) {
    oknoModulu = zamontujOknoStudio({ miejsce: glowna, kanal: w.kanal, sesje: w.sesje, modul: modulStudia });
  }
  pasmoSesji.hidden = oknoModulu !== undefined;

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
