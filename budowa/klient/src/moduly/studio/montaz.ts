/**
 * Okno modułu Studio — montaż. Woła `studio.ingest.device.list` przez
 * `wywolaj` nad `Kanal`, jak droga wejścia; wynik rozstrzyga trzy stany —
 * lista, pustka, odmowa. Bryła powtarza `.st-okno-robocze`/`.sta-okno`
 * z prototypu; wykaz urządzeń idzie klasą biblioteki `.dn-wykaz-modulu`.
 */

import {
  Command,
  StudioInputDeviceKind,
  type ErrorInfo,
  type Module,
  type Session,
  type StudioInputDevice,
} from '../../../../shared/contract.ts';
import type { Kanal } from '../../protokol/kanal.ts';
import { wywolaj } from '../../protokol/wywolanie.ts';
import { ikony, type NazwaZnaku } from './ikony.ts';
import { panelDokumentu } from './dokument.ts';
import { el, tekst, zeZnacznika } from './narzedzia.ts';

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

type StanPanelu =
  | { rodzaj: 'ladowanie' }
  | { rodzaj: 'lista'; urzadzenia: StudioInputDevice[] }
  | { rodzaj: 'pusto' }
  | { rodzaj: 'odmowa'; blad?: ErrorInfo };

function znak(nazwa: NazwaZnaku): SVGElement {
  const rysunek = zeZnacznika(ikony[nazwa]);
  rysunek.setAttribute('aria-hidden', 'true');
  return rysunek;
}

/** Znak zależny od rodzaju urządzenia — kamera dostaje własny rysunek, reszta rysunek ogólny. */
function znakUrzadzenia(rodzaj: StudioInputDevice['kind']): SVGElement {
  return znak(rodzaj === StudioInputDeviceKind.Kamera ? 'kamera' : 'urzadzenie');
}

/** Nota wiersza: rozdzielczości i podajnik — wyłącznie to, co rdzeń podał, nic dopisanego. */
function notaUrzadzenia(u: StudioInputDevice): string | null {
  const czesci: string[] = [];
  if (u.resolutions !== undefined && u.resolutions.length > 0) {
    czesci.push(`${u.resolutions.join(' · ')} ${tekst('urzadzenie.dpi')}`);
  }
  if (u.hasFeeder === true) czesci.push(tekst('urzadzenie.podajnik'));
  return czesci.length > 0 ? czesci.join(' · ') : null;
}

function wierszUrzadzenia(u: StudioInputDevice): HTMLElement {
  const nota = notaUrzadzenia(u);
  return el('div', { klasa: 'dn-wykaz-modulu-poz' }, [
    znakUrzadzenia(u.kind),
    el('b', { tekst: u.name }),
    nota === null ? null : el('span', { klasa: 'dn-meta', tekst: nota }),
  ]);
}

function panelLadowania(): HTMLElement[] {
  return [
    el('div', { klasa: 'dn-wykaz-modulu-poz' }, [
      el('span', {
        klasa: 'dn-kropka dn-kropka--sygnal dn-kropka--tetno',
        'aria-hidden': 'true',
      }),
      tekst('panel.ladowanie'),
    ]),
  ];
}

function panelPusty(): HTMLElement[] {
  return [
    el('div', { klasa: 'dn-pusty-stan' }, [
      znak('pusto'),
      el('span', { klasa: 'dn-pusty-stan-tytul', tekst: tekst('pusto.tytul') }),
      el('span', { klasa: 'dn-pusty-stan-opis', tekst: tekst('pusto.opis') }),
    ]),
  ];
}

function panelOdmowy(blad: ErrorInfo | undefined): HTMLElement[] {
  return [
    el('div', { klasa: 'dn-alert dn-alert--wstega dn-alert--blad', role: 'alert' }, [
      el('span', { klasa: 'dn-alert-znak', 'aria-hidden': 'true' }, [znak('ostrzezenie')]),
      el('span', { klasa: 'dn-alert-tresc' }, [
        el('b', { tekst: tekst('odmowa.glowa') }),
        el('span', { tekst: blad?.message ?? tekst('odmowa.brakOpisu') }),
      ]),
    ]),
  ];
}

function panelListy(urzadzenia: StudioInputDevice[]): HTMLElement[] {
  return urzadzenia.map(wierszUrzadzenia);
}

function zawartoscPanelu(stan: StanPanelu): HTMLElement[] {
  switch (stan.rodzaj) {
    case 'ladowanie':
      return panelLadowania();
    case 'pusto':
      return panelPusty();
    case 'odmowa':
      return panelOdmowy(stan.blad);
    case 'lista':
      return panelListy(stan.urzadzenia);
  }
}

export function zamontujOknoStudio(w: NastawyOknaStudio): OknoStudio {
  let zdjete = false;

  const tresc = el('div', { klasa: 'sta-okno-tresc dn-wykaz-modulu' });
  const panel = el('section', { klasa: 'sta-okno', 'data-aktywne': 'tak' }, [
    el('header', { klasa: 'sta-okno-belka' }, [
      el('span', { klasa: 'sta-okno-tytul' }, [znak('urzadzenie'), el('b', { tekst: tekst('panel.tytul') })]),
    ]),
    tresc,
  ]);
  const dokument = panelDokumentu({ kanal: w.kanal, sesje: w.sesje, modul: w.modul });
  const obszar = el('div', { klasa: 'sta-obszar', 'data-czaty': 'ukryte' }, [
    el('div', { klasa: 'sta-robocza', 'data-liczba': '2' }, [dokument.wezel, panel]),
  ]);
  const bryla = el('section', { klasa: 'st-okno-robocze', 'aria-label': tekst('okno.etykieta') }, [obszar]);

  w.miejsce.replaceChildren(bryla);
  odswiez({ rodzaj: 'ladowanie' });
  void zaladuj();

  async function zaladuj(): Promise<void> {
    const wynik = await wywolaj(w.kanal, Command.StudioIngestDeviceList, {});
    // Odmontowane w trakcie oczekiwania na odpowiedź — wynik nie ma już gdzie wylądować.
    if (zdjete) return;
    if (!wynik.udany) {
      odswiez({ rodzaj: 'odmowa', blad: wynik.blad });
      return;
    }
    const urzadzenia = wynik.wynik?.devices ?? [];
    odswiez(urzadzenia.length === 0 ? { rodzaj: 'pusto' } : { rodzaj: 'lista', urzadzenia });
  }

  function odswiez(stan: StanPanelu): void {
    tresc.replaceChildren(...zawartoscPanelu(stan));
  }

  return {
    zdejmij() {
      zdjete = true;
      dokument.zdejmij();
      if (bryla.isConnected) bryla.remove();
    },
  };
}
