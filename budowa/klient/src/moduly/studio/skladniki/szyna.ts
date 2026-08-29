/**
 * Strefa 3 — szyna dokumentów sesji. Wykaz sesji bierze z rdzenia komendą
 * `session.list` i przeciąga go od nowa na zdarzenie `session.changed`
 * (sesja powstała) oraz `window.changed` (okno weszło do sesji) — bez tego
 * szyna stała pusta, bo okno modułu zakłada sesję już po złożeniu bryły.
 *
 * Dokumenty stoją pod grupą sesji bieżącej i pochodzą ze zdarzeń
 * `studio.document.changed` własnego okna: kontrakt nie niesie komendy
 * zwracającej wykaz dokumentów, więc pozostałe sesje nazywają tę niegotowość
 * zamiast pokazywać wymyślony spis. Menu porządkowania z prototypu z tego
 * samego powodu stoi nieczynne.
 */

import {
  ChangeKind,
  Command,
  EventType,
  type ErrorInfo,
  type Session,
  type StudioDocument,
} from '../../../../../shared/contract.ts';
import type { Kanal } from '../../../protokol/kanal.ts';
import { wywolaj } from '../../../protokol/wywolanie.ts';
import { ikony, type NazwaZnaku } from '../ikony.ts';
import { el, tekst, zeZnacznika } from '../narzedzia.ts';
import { opisOdmowy } from '../odmowa.ts';

/**
 * Zdania tej szyny, których katalog modułu jeszcze nie ma. Stoją osobnym
 * wykazem, nie wplecione w kod — `moduly/studio/tresci.ts` leży poza wykazem
 * plików tego zadania, więc przeniesienie ich tam jest ruchem samej treści.
 */
const TRESCI = {
  wczytywanie: 'Wczytywanie sesji…',
  brakSesji: 'Brak sesji na tym koncie.',
  sesjaBezNazwy: 'Sesja bez nazwy',
  sesjaBiezaca: 'Sesja bieżąca',
  pozostaleSesje: 'Pozostałe sesje',
  brakPrzejscia: 'Przejście do innej sesji z tej szyny nie jest dostępne w tej wersji.',
} as const;

export interface ZaleznosciSzyny {
  /** Kanał, którym szyna woła komendy rdzenia. */
  kanal: Kanal;
  /** Okno modułu — odczytywane przy wywołaniu, bo powstaje po montażu bryły. */
  idOkna(): string | null;
}

export interface ZamontowanaSzyna {
  wezel: HTMLElement;
  zdejmij(): void;
}

function znak(rysunek: NazwaZnaku): SVGElement {
  const wezel = zeZnacznika(ikony[rysunek]);
  wezel.setAttribute('aria-hidden', 'true');
  return wezel;
}

/** Nazwa sesji z rdzenia; sesja bez tytułu dostaje nazwany stan pusty, nie identyfikator. */
function nazwaSesji(sesja: Session): string {
  return sesja.title ?? TRESCI.sesjaBezNazwy;
}

function pustyStan(zdanie: string): HTMLElement {
  return el('p', { klasa: 'dn-pusty-stan dn-pusty-stan--zwarty', tekst: zdanie });
}

function grupa(nazwa: string): HTMLElement {
  return el('div', { klasa: 'st-szyna-grupa', tekst: nazwa });
}

export function szynaDokumentow(zaleznosci: ZaleznosciSzyny): ZamontowanaSzyna {
  const dokumenty = new Map<string, StudioDocument>();
  let sesje: Session[] = [];
  let wykazWczytany = false;
  let biezacy: string | null = null;
  let odmowa: ErrorInfo | undefined;
  let odmowaWykazu: ErrorInfo | undefined;
  let zdjete = false;

  const nowyDokument = el('button', {
    klasa: 'dn-btn dn-btn--zarys dn-btn--sm',
    type: 'button',
    tekst: tekst('szyna.nowyDokument'),
  }) as HTMLButtonElement;

  const naglowek = el('div', { klasa: 'st-szyna-naglowek' }, [
    nowyDokument,
    el(
      'button',
      {
        klasa: 'dn-btn-ikona dn-etykietka',
        type: 'button',
        disabled: 'disabled',
        'data-etykietka': tekst('szyna.brakPorzadkowania'),
        'aria-label': tekst('szyna.filtry'),
      },
      [znak('filtr')],
    ),
  ]);

  const lista = el('div', { klasa: 'st-szyna-lista' });

  function pozycjaDokumentu(dokument: StudioDocument): HTMLElement {
    const biezaca = dokument.id === biezacy;
    const wezelPozycji = el(
      'button',
      biezaca
        ? { klasa: 'pt-pozycja', type: 'button', 'aria-current': 'true' }
        : { klasa: 'pt-pozycja', type: 'button' },
      [
        ...(biezaca ? [el('span', { klasa: 'pt-tetno', 'aria-hidden': 'true' })] : []),
        el('span', { klasa: 'pt-pozycja-tytul', tekst: dokument.title ?? tekst('szyna.bezTytulu') }),
      ],
    );
    wezelPozycji.addEventListener('click', () => void otworz(dokument.id));
    return wezelPozycji;
  }

  /* Sesja obca stoi nieczynna: rdzeń zna `session.open`, ale przełożenie okna
     modułu na inną sesję jest czynnością panelu dokumentu, nie szyny. */
  function pozycjaSesji(sesja: Session): HTMLElement {
    return el(
      'button',
      { klasa: 'pt-pozycja', type: 'button', disabled: 'disabled', 'data-sesja-id': sesja.id },
      [el('span', { klasa: 'pt-pozycja-tytul', tekst: nazwaSesji(sesja) })],
    );
  }

  /** Sesja, do której należy okno modułu — rozpoznana po wykazie okien z rdzenia. */
  function sesjaOkna(): Session | undefined {
    const okno = zaleznosci.idOkna();
    if (okno === null) return undefined;
    return sesje.find((sesja) => (sesja.windowIds ?? []).includes(okno));
  }

  function odswiez(): void {
    const dzieci: HTMLElement[] = [];
    if (odmowa !== undefined) dzieci.push(pustyStan(opisOdmowy(odmowa, 'szyna')));

    if (!wykazWczytany) {
      lista.replaceChildren(...dzieci, pustyStan(TRESCI.wczytywanie));
      return;
    }
    if (odmowaWykazu !== undefined) {
      lista.replaceChildren(...dzieci, pustyStan(opisOdmowy(odmowaWykazu, 'szyna.wykazSesji')));
      return;
    }
    if (sesje.length === 0 && dokumenty.size === 0) {
      lista.replaceChildren(...dzieci, pustyStan(TRESCI.brakSesji));
      return;
    }

    const wlasna = sesjaOkna();
    dzieci.push(grupa(wlasna === undefined ? TRESCI.sesjaBiezaca : nazwaSesji(wlasna)));
    if (dokumenty.size === 0) dzieci.push(pustyStan(tekst('szyna.brakDokumentow')));
    else dzieci.push(...Array.from(dokumenty.values(), pozycjaDokumentu));

    const pozostale = sesje.filter((sesja) => sesja.id !== wlasna?.id);
    if (pozostale.length > 0) {
      dzieci.push(grupa(TRESCI.pozostaleSesje));
      dzieci.push(...pozostale.map(pozycjaSesji));
      dzieci.push(pustyStan(TRESCI.brakPrzejscia));
    }
    lista.replaceChildren(...dzieci);
  }

  async function wczytajSesje(): Promise<void> {
    const wynik = await wywolaj(zaleznosci.kanal, Command.SessionList, {});
    if (zdjete) return;
    wykazWczytany = true;
    odmowaWykazu = wynik.udany ? undefined : wynik.blad;
    if (wynik.udany) sesje = wynik.wynik?.sessions ?? [];
    odswiez();
  }

  async function zaloz(): Promise<void> {
    const okno = zaleznosci.idOkna();
    if (okno === null) return;
    nowyDokument.disabled = true;
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioDocumentCreate, {
      windowId: okno,
      title: tekst('szyna.nazwaNowego'),
    });
    if (zdjete) return;
    nowyDokument.disabled = false;
    odmowa = wynik.udany ? undefined : wynik.blad;
    odswiez();
  }

  async function otworz(idDokumentu: string): Promise<void> {
    const okno = zaleznosci.idOkna();
    if (okno === null) return;
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioDocumentOpen, {
      windowId: okno,
      documentId: idDokumentu,
    });
    if (zdjete) return;
    odmowa = wynik.udany ? undefined : wynik.blad;
    if (wynik.udany) biezacy = idDokumentu;
    odswiez();
  }

  nowyDokument.addEventListener('click', () => void zaloz());

  const odsubDokumenty = zaleznosci.kanal.naZdarzenie(EventType.StudioDocumentChanged, (zdarzenie) => {
    const okno = zaleznosci.idOkna();
    if (zdjete || okno === null || zdarzenie.document.windowId !== okno) return;
    if (zdarzenie.change === ChangeKind.Deleted) {
      dokumenty.delete(zdarzenie.document.id);
      if (biezacy === zdarzenie.document.id) biezacy = null;
    } else {
      dokumenty.set(zdarzenie.document.id, zdarzenie.document);
      biezacy = zdarzenie.document.id;
    }
    odswiez();
  });

  /* Sesja powstaje po złożeniu bryły, a okno wchodzi do niej jeszcze później —
     oba zdarzenia zmieniają wykaz, więc oba przeciągają go od nowa. */
  const odsubSesje = zaleznosci.kanal.naZdarzenie(EventType.SessionChanged, () => {
    if (!zdjete) void wczytajSesje();
  });
  const odsubOkna = zaleznosci.kanal.naZdarzenie(EventType.WindowChanged, () => {
    if (!zdjete) void wczytajSesje();
  });

  odswiez();
  void wczytajSesje();

  return {
    wezel: el('nav', { klasa: 'st-szyna', 'aria-label': tekst('szyna.etykieta') }, [naglowek, lista]),
    zdejmij() {
      zdjete = true;
      odsubDokumenty();
      odsubSesje();
      odsubOkna();
    },
  };
}
