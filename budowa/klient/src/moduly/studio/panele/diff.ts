/**
 * Diff/Grep Panel — porównanie dwóch wersji dokumentu Studio: fragmenty dodane,
 * usunięte i zmienione, wyszukiwanie wzorca w treści porównywanych wersji,
 * oraz zmiany odłożone przez śledzenie zmian wraz z decyzją Operatora.
 *
 * Panel nie zakłada okna ani dokumentu — czeka na `studio.document.changed`
 * dla okna wskazanego w zależnościach i dopiero wtedy woła rodziny komend
 * `studio.diff.*`, `studio.tracking.*` i `studio.annotation.*`. Zdarzenie
 * sprzed zamontowania panelu nie dociera do niego: rdzeń nie daje dziś drogi
 * do odpytania o dokument bieżący okna spoza tych trzech rodzin.
 */

import {
  ChangeKind,
  Command,
  DiffHunkKind,
  EventType,
  StudioChangeKind,
  type ErrorInfo,
  type StudioDiffHunk,
  type StudioTextMatch,
  type StudioTrackedChange,
} from '../../../../../shared/contract.ts';
import { wywolaj } from '../../../protokol/wywolanie.ts';
import { ikony } from '../ikony.ts';
import { el, tekst, zeZnacznika } from '../narzedzia.ts';
import { tresc as t } from './diff-tresci.ts';
import type { MontazPanelu } from './umowa.ts';
import { opisOdmowy } from '../odmowa.ts';
import { zLiczba } from '../liczebnik.ts';

/** Wynik porównania wersji, licząc trafienia wzorca jako część tej samej odpowiedzi. */
type ZasobPorownania =
  | { rodzaj: 'ladowanie' }
  | { rodzaj: 'odmowa'; blad?: ErrorInfo }
  | { rodzaj: 'dane'; hunks: StudioDiffHunk[]; trafienia: StudioTextMatch[] | null };

/** Wykaz zmian oczekujących na decyzję Operatora. */
type ZasobZmian =
  | { rodzaj: 'ladowanie' }
  | { rodzaj: 'odmowa'; blad?: ErrorInfo }
  | { rodzaj: 'dane'; zmiany: StudioTrackedChange[] };

const ZNAK_LUPA =
  '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><circle cx="11" cy="11" r="7"/><path d="m20 20-3.5-3.5"/></svg>';

function ikonaLupa(): SVGElement {
  return zeZnacznika(ZNAK_LUPA);
}

function ikonaDymek(): SVGElement {
  return zeZnacznika(ikony.dymek);
}

/** Znaczek odmowy jednakowy dla porównania, zmian i adnotacji: tytuł nazywa czynność, treść — zdanie dobrane po kodzie odmowy. */
function alertOdmowy(tytul: string, blad?: ErrorInfo): HTMLElement {
  return el('div', { klasa: 'dn-alert dn-alert--wstega dn-alert--blad', role: 'alert' }, [
    el('span', { klasa: 'dn-alert-tresc' }, [
      el('b', { tekst: tytul }),
      el('span', { tekst: opisOdmowy(blad, 'diff') }),
    ]),
  ]);
}

function wierszLadowania(napis: string): HTMLElement {
  return el('div', { klasa: 'dn-wykaz-modulu-poz' }, [
    el('span', { klasa: 'dn-kropka dn-kropka--sygnal dn-kropka--tetno', 'aria-hidden': 'true' }),
    napis,
  ]);
}

export const panelDiff: MontazPanelu = (wezel, zal) => {
  let zdjete = false;
  let dokumentId: string | null = null;
  /* Podbijane przy każdej zmianie dokumentu; odpowiedź spóźniona wobec
     dokumentu, który panel już porzucił, ląduje w nicości zamiast nadpisać
     stan nowszego dokumentu. */
  let pokolenie = 0;

  let wersjaBazowa = '';
  let wersjaCel = '';
  let wzorzec = '';
  let regex = false;

  let porownanie: ZasobPorownania | null = null;
  let zmianyStan: ZasobZmian | null = null;

  let adnotowanyHunk: number | null = null;
  let trescAdnotacji = '';
  let adnotacjaWTrakcie = false;
  const bladAdnotacji = new Map<number, ErrorInfo | undefined>();

  const decyzjeWTrakcie = new Set<string>();
  const bledyDecyzji = new Map<string, ErrorInfo | undefined>();

  function ustawDokument(id: string): void {
    if (id === dokumentId) return;
    dokumentId = id;
    pokolenie += 1;
    porownanie = null;
    zmianyStan = null;
    void uruchomPorownanie();
    void wczytajZmiany();
  }

  function usunDokument(): void {
    dokumentId = null;
    pokolenie += 1;
    porownanie = null;
    zmianyStan = null;
    odswiez();
  }

  async function uruchomPorownanie(): Promise<void> {
    if (dokumentId === null) return;
    const mojDokument = dokumentId;
    const mojePokolenie = pokolenie;
    porownanie = { rodzaj: 'ladowanie' };
    odswiez();

    const wzorzecPrzyciety = wzorzec.trim();
    const wynik = await wywolaj(zal.kanal, Command.StudioDiffCompare, {
      documentId: mojDokument,
      baseVersionId: wersjaBazowa.trim() || undefined,
      targetVersionId: wersjaCel.trim() || undefined,
      pattern: wzorzecPrzyciety === '' ? undefined : wzorzecPrzyciety,
      regex: wzorzecPrzyciety === '' ? undefined : regex,
    });
    if (zdjete || mojePokolenie !== pokolenie) return;

    porownanie =
      wynik.udany && wynik.wynik !== undefined
        ? {
            rodzaj: 'dane',
            hunks: wynik.wynik.hunks ?? [],
            trafienia: wzorzecPrzyciety === '' ? null : (wynik.wynik.matches ?? []),
          }
        : { rodzaj: 'odmowa', blad: wynik.blad };
    odswiez();
  }

  async function wczytajZmiany(): Promise<void> {
    if (dokumentId === null) return;
    const mojDokument = dokumentId;
    const mojePokolenie = pokolenie;
    zmianyStan = { rodzaj: 'ladowanie' };
    odswiez();

    const wynik = await wywolaj(zal.kanal, Command.StudioTrackingList, {
      documentId: mojDokument,
      pendingOnly: true,
    });
    if (zdjete || mojePokolenie !== pokolenie) return;

    zmianyStan =
      wynik.udany && wynik.wynik !== undefined
        ? { rodzaj: 'dane', zmiany: wynik.wynik.changes }
        : { rodzaj: 'odmowa', blad: wynik.blad };
    odswiez();
  }

  async function decyduj(changeId: string, przyjmij: boolean): Promise<void> {
    if (dokumentId === null) return;
    const mojePokolenie = pokolenie;
    decyzjeWTrakcie.add(changeId);
    bledyDecyzji.delete(changeId);
    odswiez();

    const wynik = await wywolaj(zal.kanal, Command.StudioTrackingDecide, {
      documentId: dokumentId,
      changeIds: [changeId],
      accept: przyjmij,
    });
    if (zdjete || mojePokolenie !== pokolenie) return;
    decyzjeWTrakcie.delete(changeId);

    if (!wynik.udany) {
      bledyDecyzji.set(changeId, wynik.blad);
      odswiez();
      return;
    }
    await wczytajZmiany();
  }

  async function dodajAdnotacje(hunkIndex: number): Promise<void> {
    const body = trescAdnotacji.trim();
    if (body === '' || dokumentId === null) return;
    const mojePokolenie = pokolenie;
    adnotacjaWTrakcie = true;
    bladAdnotacji.delete(hunkIndex);
    odswiez();

    const wynik = await wywolaj(zal.kanal, Command.StudioAnnotationAdd, {
      documentId: dokumentId,
      hunkIndex,
      baseVersionId: wersjaBazowa.trim() || undefined,
      targetVersionId: wersjaCel.trim() || undefined,
      body,
    });
    if (zdjete || mojePokolenie !== pokolenie) return;
    adnotacjaWTrakcie = false;

    if (!wynik.udany) {
      bladAdnotacji.set(hunkIndex, wynik.blad);
    } else {
      adnotowanyHunk = null;
      trescAdnotacji = '';
    }
    odswiez();
  }

  function wierszWersji(): HTMLElement {
    const poleBazowa = el('input', {
      klasa: 'dn-pole-kontrolka',
      type: 'text',
      'aria-label': t.porownanie.etykietaBazowa,
      placeholder: t.porownanie.zastepczaWersja,
      value: wersjaBazowa,
    }) as HTMLInputElement;
    poleBazowa.addEventListener('input', () => {
      wersjaBazowa = poleBazowa.value;
    });
    poleBazowa.addEventListener('keydown', (zdarzenie) => {
      if (zdarzenie.key === 'Enter') void uruchomPorownanie();
    });

    const poleCel = el('input', {
      klasa: 'dn-pole-kontrolka',
      type: 'text',
      'aria-label': t.porownanie.etykietaCel,
      placeholder: t.porownanie.zastepczaWersja,
      value: wersjaCel,
    }) as HTMLInputElement;
    poleCel.addEventListener('input', () => {
      wersjaCel = poleCel.value;
    });
    poleCel.addEventListener('keydown', (zdarzenie) => {
      if (zdarzenie.key === 'Enter') void uruchomPorownanie();
    });

    const przyciskPorownaj = el('button', {
      klasa: 'dn-btn dn-btn--zarys dn-btn--sm',
      type: 'button',
      tekst: t.porownanie.porownaj,
    });
    przyciskPorownaj.addEventListener('click', () => void uruchomPorownanie());

    return el('div', { klasa: 'st-panel-wiersz' }, [poleBazowa, '→', poleCel, przyciskPorownaj]);
  }

  function wierszGrep(): HTMLElement {
    const poleGrep = el('input', {
      type: 'search',
      'aria-label': t.grep.etykieta,
      placeholder: t.grep.zastepczaTresc,
      value: wzorzec,
    }) as HTMLInputElement;
    poleGrep.addEventListener('input', () => {
      wzorzec = poleGrep.value;
    });
    poleGrep.addEventListener('keydown', (zdarzenie) => {
      if (zdarzenie.key === 'Enter') void uruchomPorownanie();
    });
    const etykietaGrep = el('label', { klasa: 'dn-szukaj' }, [ikonaLupa(), poleGrep]);

    const poleRegex = el('input', { type: 'checkbox' }) as HTMLInputElement;
    poleRegex.checked = regex;
    poleRegex.addEventListener('change', () => {
      regex = poleRegex.checked;
    });
    const etykietaRegex = el('label', { klasa: 'dn-check-etyk' }, [poleRegex, t.grep.regex]);

    return el('div', { klasa: 'st-panel-wiersz' }, [etykietaGrep, etykietaRegex]);
  }

  function wierszKompozytoraAdnotacji(hunkIndex: number): HTMLElement {
    const pole = el('input', {
      klasa: 'dn-pole-kontrolka',
      type: 'text',
      'aria-label': t.adnotacja.etykieta,
      placeholder: t.adnotacja.zastepczaTresc,
      value: trescAdnotacji,
      disabled: adnotacjaWTrakcie,
    }) as HTMLInputElement;
    const przycisk = el('button', {
      klasa: 'dn-btn dn-btn--zarys dn-btn--sm',
      type: 'button',
      tekst: t.adnotacja.dodaj,
      disabled: adnotacjaWTrakcie || trescAdnotacji.trim() === '',
    }) as HTMLButtonElement;
    przycisk.addEventListener('click', () => void dodajAdnotacje(hunkIndex));

    /* Odświeżenie całego panelu przy każdym znaku gubiłoby ognisko pola — tylko
       stan przycisku reaguje na wpis, reszta widoku czeka na wynik zgłoszenia. */
    pole.addEventListener('input', () => {
      trescAdnotacji = pole.value;
      przycisk.disabled = trescAdnotacji.trim() === '';
    });

    const dzieci: (Node | string)[] = [el('div', { klasa: 'dn-pole-zestaw' }, [pole, przycisk])];
    if (bladAdnotacji.has(hunkIndex)) {
      dzieci.push(alertOdmowy(t.odmowa.adnotacja, bladAdnotacji.get(hunkIndex)));
    }
    return el('div', { klasa: 'st-panel-wiersz' }, dzieci);
  }

  function wierszHunka(hunk: StudioDiffHunk): HTMLElement[] {
    const tresc: (Node | string)[] = [];
    if (hunk.kind === DiffHunkKind.Added) {
      tresc.push(el('span', { klasa: 'dn-diff-dod', tekst: hunk.after ?? '' }));
    } else if (hunk.kind === DiffHunkKind.Removed) {
      tresc.push(el('span', { klasa: 'dn-diff-usu', tekst: hunk.before ?? '' }));
    } else if (hunk.kind === DiffHunkKind.Changed) {
      tresc.push(el('span', { klasa: 'dn-diff-usu', tekst: hunk.before ?? '' }));
      tresc.push(' ');
      tresc.push(el('span', { klasa: 'dn-diff-dod', tekst: hunk.after ?? '' }));
    } else {
      tresc.push(hunk.after ?? hunk.before ?? '');
    }

    const przyciskAdnotacji = el(
      'button',
      { klasa: 'dn-btn-ikona dn-btn-ikona--sm', type: 'button', 'aria-label': t.adnotacja.etykieta },
      [ikonaDymek()],
    );
    przyciskAdnotacji.addEventListener('click', () => {
      adnotowanyHunk = adnotowanyHunk === hunk.index ? null : hunk.index;
      trescAdnotacji = '';
      odswiez();
    });

    const wiersz = el('div', { klasa: 'st-panel-wiersz' }, [...tresc, przyciskAdnotacji]);
    return adnotowanyHunk === hunk.index ? [wiersz, wierszKompozytoraAdnotacji(hunk.index)] : [wiersz];
  }

  function statystykaPorownania(hunks: StudioDiffHunk[]): HTMLElement {
    let dodane = 0;
    let usuniete = 0;
    let zmienione = 0;
    for (const hunk of hunks) {
      if (hunk.kind === DiffHunkKind.Added) dodane += 1;
      else if (hunk.kind === DiffHunkKind.Removed) usuniete += 1;
      else if (hunk.kind === DiffHunkKind.Changed) zmienione += 1;
    }
    return el('div', {
      klasa: 'dn-nota',
      tekst:
        `${t.porownanie.statystykaEtykieta} ` +
        `${t.porownanie.znakDodania}${zLiczba(dodane, t.porownanie.jednostkaFragmenty)} · ` +
        `${t.porownanie.znakUsuniecia}${zLiczba(usuniete, t.porownanie.jednostkaFragmenty)} · ` +
        `${zLiczba(zmienione, t.porownanie.jednostkaZmian)}`,
    });
  }

  function widokTrafien(trafienia: StudioTextMatch[]): HTMLElement[] {
    if (trafienia.length === 0) {
      return [el('p', { klasa: 'dn-pusty-stan dn-pusty-stan--zwarty', tekst: t.grep.brakTrafien })];
    }
    return trafienia.map((trafienie) =>
      el('div', { klasa: 'st-panel-wiersz' }, [
        el('span', { klasa: 'dn-meta', tekst: String(trafienie.line) }),
        trafienie.text,
      ]),
    );
  }

  function widokPorownania(): HTMLElement[] {
    const glowa = [wierszWersji(), wierszGrep()];
    if (porownanie === null) return glowa;
    if (porownanie.rodzaj === 'ladowanie') return [...glowa, wierszLadowania(t.porownanie.ladowanie)];
    if (porownanie.rodzaj === 'odmowa') return [...glowa, alertOdmowy(t.odmowa.porownanie, porownanie.blad)];

    const wiersze: HTMLElement[] = [];
    if (porownanie.hunks.length === 0) {
      wiersze.push(el('p', { klasa: 'dn-pusty-stan dn-pusty-stan--zwarty', tekst: t.porownanie.brakRoznic }));
    } else {
      for (const hunk of porownanie.hunks) wiersze.push(...wierszHunka(hunk));
      wiersze.push(statystykaPorownania(porownanie.hunks));
    }
    if (porownanie.trafienia !== null) wiersze.push(...widokTrafien(porownanie.trafienia));
    return [...glowa, ...wiersze];
  }

  function wierszZmiany(zmiana: StudioTrackedChange): HTMLElement[] {
    const klasa =
      zmiana.kind === StudioChangeKind.Wstawienie
        ? 'dn-diff-dod'
        : zmiana.kind === StudioChangeKind.Usuniecie
          ? 'dn-diff-usu'
          : 'dn-diff-fmt';
    const trescZmiany = zmiana.kind === StudioChangeKind.Usuniecie ? zmiana.before : (zmiana.after ?? zmiana.before);

    const dzieci: (Node | string)[] = [];
    if (trescZmiany !== undefined && trescZmiany !== '') {
      dzieci.push(el('span', { klasa, tekst: trescZmiany }));
    }

    const wTrakcie = decyzjeWTrakcie.has(zmiana.id);
    const przyjmij = el('button', {
      klasa: 'dn-btn dn-btn--zarys dn-btn--sm',
      type: 'button',
      tekst: t.zmiany.akceptuj,
      disabled: wTrakcie,
    });
    przyjmij.addEventListener('click', () => void decyduj(zmiana.id, true));

    const odrzuc = el('button', {
      klasa: 'dn-btn dn-btn--duch dn-btn--sm',
      type: 'button',
      tekst: t.zmiany.odrzuc,
      disabled: wTrakcie,
    });
    odrzuc.addEventListener('click', () => void decyduj(zmiana.id, false));

    dzieci.push(przyjmij, odrzuc);
    const wiersz = el('div', { klasa: 'st-panel-wiersz' }, dzieci);

    return bledyDecyzji.has(zmiana.id)
      ? [wiersz, alertOdmowy(t.odmowa.decyzja, bledyDecyzji.get(zmiana.id))]
      : [wiersz];
  }

  function widokZmian(): HTMLElement[] {
    if (zmianyStan === null) return [];
    const naglowek = el('div', { klasa: 'dn-nota', tekst: t.zmiany.naglowek });

    if (zmianyStan.rodzaj === 'ladowanie') return [naglowek, wierszLadowania(t.zmiany.ladowanie)];
    if (zmianyStan.rodzaj === 'odmowa') return [naglowek, alertOdmowy(t.odmowa.zmiany, zmianyStan.blad)];
    if (zmianyStan.zmiany.length === 0) {
      return [naglowek, el('p', { klasa: 'dn-pusty-stan dn-pusty-stan--zwarty', tekst: t.zmiany.brak })];
    }
    return [naglowek, ...zmianyStan.zmiany.flatMap(wierszZmiany)];
  }

  function zawartosc(): HTMLElement[] {
    if (zal.idOkna === null) {
      return [alertOdmowy(tekst('dokumentOdmowa.okno'))];
    }
    if (dokumentId === null) {
      return [el('p', { klasa: 'dn-pusty-stan dn-pusty-stan--zwarty', tekst: t.brakDokumentu })];
    }
    return [...widokPorownania(), ...widokZmian()];
  }

  function odswiez(): void {
    wezel.replaceChildren(...zawartosc());
  }

  wezel.classList.add('st-panel-lista');
  odswiez();

  const odsubDokumentu =
    zal.idOkna === null
      ? null
      : zal.kanal.naZdarzenie(EventType.StudioDocumentChanged, (zdarzenie) => {
          if (zdarzenie.document.windowId !== zal.idOkna) return;
          if (zdarzenie.change === ChangeKind.Deleted) {
            if (zdarzenie.document.id === dokumentId) usunDokument();
            return;
          }
          ustawDokument(zdarzenie.document.id);
        });

  return {
    zdejmij() {
      zdjete = true;
      odsubDokumentu?.();
      wezel.replaceChildren();
    },
  };
};
