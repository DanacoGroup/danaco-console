/**
 * Strefa 3 — szyna dokumentów sesji. Kontrakt nie niesie komendy zwracającej
 * wykaz dokumentów sesji, więc szyna nie wczytuje listy: buduje ją ze zdarzeń
 * `studio.document.changed` ograniczonych do własnego okna — te same, z których
 * dokument poznaje Tools Panel. Pozycja stoi w szynie dlatego, że rdzeń o niej
 * powiedział, nie dlatego, że klient ją wymyślił.
 *
 * Grup projektowych z prototypu (`.st-szyna-grupa`) szyna nie stawia: dokument
 * w kontrakcie nie niesie projektu, a nagłówek grupy bez przynależności byłby
 * podziałem zmyślonym. Menu filtrowania i sortowania z tego samego powodu nie
 * ma pozycji — porządkowałoby po polach, których rdzeń nie podaje.
 */

import {
  ChangeKind,
  Command,
  EventType,
  type ErrorInfo,
  type StudioDocument,
} from '../../../../../shared/contract.ts';
import type { Kanal } from '../../../protokol/kanal.ts';
import { wywolaj } from '../../../protokol/wywolanie.ts';
import { ikony, type NazwaZnaku } from '../ikony.ts';
import { el, tekst, zeZnacznika } from '../narzedzia.ts';
import { opisOdmowy } from '../odmowa.ts';

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

export function szynaDokumentow(zaleznosci: ZaleznosciSzyny): ZamontowanaSzyna {
  const dokumenty = new Map<string, StudioDocument>();
  let biezacy: string | null = null;
  let odmowa: ErrorInfo | undefined;
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

  function pozycja(dokument: StudioDocument): HTMLElement {
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

  function odswiez(): void {
    if (odmowa !== undefined) {
      lista.replaceChildren(
        el('p', { klasa: 'dn-pusty-stan dn-pusty-stan--zwarty', tekst: opisOdmowy(odmowa, 'szyna') }),
      );
      return;
    }
    if (dokumenty.size === 0) {
      lista.replaceChildren(
        el('p', { klasa: 'dn-pusty-stan dn-pusty-stan--zwarty', tekst: tekst('szyna.brakDokumentow') }),
      );
      return;
    }
    lista.replaceChildren(...Array.from(dokumenty.values(), pozycja));
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

  const odsubskrybuj = zaleznosci.kanal.naZdarzenie(EventType.StudioDocumentChanged, (zdarzenie) => {
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

  odswiez();

  return {
    wezel: el('nav', { klasa: 'st-szyna', 'aria-label': tekst('szyna.etykieta') }, [naglowek, lista]),
    zdejmij() {
      zdjete = true;
      odsubskrybuj();
    },
  };
}
