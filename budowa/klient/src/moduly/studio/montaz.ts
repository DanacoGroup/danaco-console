/**
 * Okno modułu Studio — montaż. Powtarza kształt `rama/montaz.ts`: składa
 * strefy z prototypu (`design/05-okna/moduly/studio.html`) — pasmo siedmiu
 * kart, wstążkę okna roboczego, szynę dokumentów sesji, okno czatu i zaczep
 * paneli — i wstawia je w miejsce wskazane przez ramę. Karta „Studio Editor"
 * niesie jedyny dziś rzeczywisty przekrój do rdzenia (`dokument.ts`); pozostałe
 * sześć kart wchodzą osobnym zakresem prac.
 *
 * Pas stanu (strefa 6 z wpisu terenu) montuje `rama/skladniki/stan.ts` —
 * stoi raz dla całej powłoki aplikacji, nie osobno na moduł.
 */

import type { Module, Session } from '../../../../shared/contract.ts';
import type { Kanal } from '../../protokol/kanal.ts';
import { panelDokumentu } from './dokument.ts';
import { el, tekst } from './narzedzia.ts';
import { KARTA_EDITOR, karty } from './skladniki/definicje.ts';
import { pasmoKart } from './skladniki/pasmo.ts';
import { wstazkaOkna } from './skladniki/wstazka.ts';
import { szynaDokumentow } from './skladniki/szyna.ts';
import { oknoCzatu } from './skladniki/czat.ts';
import { panelNiegotowy } from './skladniki/panel-niegotowy.ts';
import type { MontazPanelu, ZamontowanyPanel } from './panele/umowa.ts';
import { montujPanelTools } from './panele/tools.ts';
import { panelDiff } from './panele/diff.ts';
import { montujPanelRepo } from './panele/repo.ts';
import { montujPodgladWydania } from './panele/preview.ts';
import { montujPliki } from './panele/pliki.ts';
import { montujPanelPlanu } from './panele/plan.ts';

/* Wykaz montaży paneli kartą. Karta bez wpisu zostaje przy stanie „panel jeszcze
   nie powstał" — wpis dopisany tutaj jest jedynym miejscem, w którym panel wchodzi
   do okna, więc dołożenie kolejnego nie dotyka żadnego z pozostałych plików. */
const montazePaneli: Record<string, MontazPanelu> = {
  tools: montujPanelTools,
  diff: panelDiff,
  repo: montujPanelRepo,
  preview: montujPodgladWydania,
  pliki: montujPliki,
  plan: montujPanelPlanu,
};

export interface NastawyOknaStudio {
  /** Miejsce w dokumencie, w które okno się wstawia — `main` ramy aplikacji. */
  miejsce: HTMLElement;
  /** Kanał, którym okno woła komendy rdzenia. */
  kanal: Kanal;
  /** Karty sesji odtworzone przez rdzeń — pod nimi staje okno modułu. */
  sesje: Session[];
  /** Moduł, którego okno powstaje w rdzeniu. */
  modul: Module;
}

export interface OknoStudio {
  /** Zdejmuje okno z miejsca montażu; wynik wywołania w locie ląduje w nicości. */
  zdejmij(): void;
}

/**
 * Wiąże klik na karcie pasma z widocznością panelu w zaczepie: oznacza kartę
 * wskazaną jako bieżącą i odsłania wyłącznie panel jej odpowiadający — ten
 * sam mechanizm co `karty-okna.js` dla pasma aplikacji, tu własny, bo ten
 * skrypt czyta wyłącznie pasmo poziomu aplikacji, nie zagnieżdżone pasmo
 * okna roboczego.
 */
function wirujKarty(bryla: HTMLElement, domontuj: (kod: string) => void): (kod: string) => void {
  const wykazKart = Array.from(bryla.querySelectorAll<HTMLElement>('.dn-karta[data-karta]'));
  const wykazPaneli = Array.from(bryla.querySelectorAll<HTMLElement>('.sta-robocza > [role="tabpanel"]'));
  const lista = bryla.querySelector('.st-karty');

  function pokaz(kodKarty: string): void {
    for (const inna of wykazKart) {
      const biezaca = inna.dataset['karta'] === kodKarty;
      inna.setAttribute('aria-selected', biezaca ? 'true' : 'false');
      inna.tabIndex = biezaca ? 0 : -1;
    }
    domontuj(kodKarty);
    for (const panel of wykazPaneli) {
      panel.hidden = panel.id !== `panel-${kodKarty}`;
    }
  }

  lista?.addEventListener('click', (zdarzenie) => {
    const wskazana = (zdarzenie.target as Element | null)?.closest('.dn-karta[data-karta]') as HTMLElement | null;
    const kodKarty = wskazana?.dataset['karta'];
    if (kodKarty !== undefined) pokaz(kodKarty);
  });

  return pokaz;
}

export function zamontujOknoStudio(w: NastawyOknaStudio): OknoStudio {
  const dokument = panelDokumentu({ kanal: w.kanal, sesje: w.sesje, modul: w.modul });
  /* Panel wchodzi dopiero przy pierwszym wejściu na jego kartę. Wcześniej nie ma
     po co: okno modułu i sesja powstają wywołaniem do rdzenia już po montażu bryły,
     a panel bez `windowId` nie ma czym zawołać ani jednej komendy Studia. */
  let pokazKarte: ((kod: string) => void) | undefined;
  const zamontowane: ZamontowanyPanel[] = [];
  const czekajace = new Map<string, HTMLElement>();
  const panele = karty
    .filter((k) => k.kod !== KARTA_EDITOR)
    .map((k) => {
      const wezel = panelNiegotowy(k);
      if (montazePaneli[k.kod]) czekajace.set(k.kod, wezel);
      return wezel;
    });

  function domontuj(kod: string): void {
    const wezel = czekajace.get(kod);
    const montuj = montazePaneli[kod];
    if (wezel === undefined || montuj === undefined) return;
    czekajace.delete(kod);
    wezel.replaceChildren();
    zamontowane.push(montuj(wezel, {
      kanal: w.kanal,
      idOkna: dokument.idOkna(),
      idSesji: dokument.idSesji(),
    }));
  }
  const robocza = el('div', { klasa: 'sta-robocza' }, [dokument.wezel, ...panele]);
  const obszar = el('div', { klasa: 'sta-obszar', 'data-robocza': 'widoczna' }, [oknoCzatu(), robocza]);
  const cialo = el('div', { klasa: 'sta-cialo st-cialo' }, [szynaDokumentow(), obszar]);

  /* `<section>`, nie `<main>` jak w prototypie: bryła wchodzi wewnątrz
     `<main id="dn-obszar-glowna">` ramy, a dokument nie niesie dwóch `<main>`. */
  const bryla = el('section', { klasa: 'st-okno-robocze', 'aria-label': tekst('okno.etykieta') }, [
    pasmoKart(),
    /* Karta wskazana z wstążki: bryła jeszcze nie stoi, więc przejście woła
       przez uchwyt uzupełniony zaraz po jej złożeniu. */
    wstazkaOkna({ sesje: w.sesje, naKarte: (kod) => pokazKarte?.(kod) }),
    cialo,
  ]);

  w.miejsce.replaceChildren(bryla);
  pokazKarte = wirujKarty(bryla, domontuj);

  return {
    zdejmij() {
      dokument.zdejmij();
      if (bryla.isConnected) bryla.remove();
    },
  };
}
