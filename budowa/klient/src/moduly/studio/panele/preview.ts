/**
 * Karta „Preview Window" okna roboczego Studia — podgląd wydania dokumentu:
 * układ i paginacja w formacie docelowym, nie sam tekst. Panel nie zna
 * dokumentu z góry — dostaje go reaktywnie zdarzeniem `studio.document.changed`,
 * którego opis w kontrakcie wprost nazywa ten panel jako odbiorcę.
 *
 * Trzy warstwy stanu rdzenia idą niezależnie: render podglądu
 * (`studio.preview.render`), wykaz profili wydania (`studio.export.profile.list`)
 * i ustawienia strony bieżącego dokumentu (`studio.page.setup.get`, wyłącznie
 * do odczytu — zmiana strony należy do innej karty). Bez okna modułu żadna
 * z tych komend nie ma na czym stanąć, więc panel wtedy nie próbuje żadnej.
 */

import {
  ChangeKind,
  Command,
  EventType,
  StudioExportFormat,
  type ErrorInfo,
  type StudioDocument,
  type StudioExportProfile,
  type StudioPageSetup,
  type StudioPreviewRenderResponse,
} from '../../../../../shared/contract.ts';
import { wywolaj } from '../../../protokol/wywolanie.ts';
import { el, tekst, zeZnacznika } from '../narzedzia.ts';
import { ikony } from '../ikony.ts';
import type { MontazPanelu } from './umowa.ts';
import { trescPodgladu } from './preview-tresci.ts';

type TrybWidoku = 'strona' | 'ciagly';

interface StanWidoku {
  format: string;
  tryb: TrybWidoku;
  powiekszenie: number;
  podzielony: boolean;
  profileId: string | undefined;
  strona: number;
}

type StanDokumentu = { rodzaj: 'brak' } | { rodzaj: 'obecny'; dokument: StudioDocument };

type StanRenderu =
  | { rodzaj: 'bezczynny' }
  | { rodzaj: 'renderowanie' }
  | { rodzaj: 'gotowy'; wynik: StudioPreviewRenderResponse }
  | { rodzaj: 'odmowa'; blad: ErrorInfo | undefined };

type StanProfili =
  | { rodzaj: 'wczytywanie' }
  | { rodzaj: 'wykaz'; profile: StudioExportProfile[] }
  | { rodzaj: 'odmowa'; blad: ErrorInfo | undefined };

type StanUkladu =
  | { rodzaj: 'brak' }
  | { rodzaj: 'znany'; pageSetup: StudioPageSetup }
  | { rodzaj: 'odmowa'; blad: ErrorInfo | undefined };

const POWIEKSZENIA: readonly number[] = [50, 75, 100, 125, 150, 200];
const FORMATY: readonly string[] = Object.values(StudioExportFormat);

export const montujPodgladWydania: MontazPanelu = (wezel, zaleznosci) => {
  let zdjete = false;
  let epokaRenderu = 0;
  let epokaUkladu = 0;

  const ui: StanWidoku = {
    format: StudioExportFormat.Pdf,
    tryb: 'strona',
    powiekszenie: 100,
    podzielony: false,
    profileId: undefined,
    strona: 1,
  };
  let dokument: StanDokumentu = { rodzaj: 'brak' };
  let render: StanRenderu = { rodzaj: 'bezczynny' };
  let profile: StanProfili = { rodzaj: 'wczytywanie' };
  let uklad: StanUkladu = { rodzaj: 'brak' };

  const tresc = el('div', { klasa: 'sta-okno-tresc' });
  const bryla = el(
    'section',
    { klasa: 'sta-okno', id: 'panel-preview', role: 'tabpanel', 'aria-labelledby': 'karta-preview', hidden: true },
    [
      el('header', { klasa: 'sta-okno-belka' }, [
        el('span', { klasa: 'sta-okno-tytul' }, [el('b', { tekst: tekst('karty.preview') })]),
      ]),
      tresc,
    ],
  );
  wezel.appendChild(bryla);

  let odsubDokument: (() => void) | null = null;

  if (zaleznosci.idOkna !== null) {
    const idOkna = zaleznosci.idOkna;
    void wczytajProfile();
    odsubDokument = zaleznosci.kanal.naZdarzenie(EventType.StudioDocumentChanged, (zdarzenie) => {
      if (zdjete || zdarzenie.document.windowId !== idOkna) return;
      if (zdarzenie.change === ChangeKind.Deleted) {
        dokument = { rodzaj: 'brak' };
        render = { rodzaj: 'bezczynny' };
        uklad = { rodzaj: 'brak' };
        odswiez();
        return;
      }
      dokument = { rodzaj: 'obecny', dokument: zdarzenie.document };
      ui.strona = 1;
      void wczytajUklad(zdarzenie.document.id);
      void wykonajRender(zdarzenie.document.id);
    });
  }

  odswiez();

  async function wczytajProfile(): Promise<void> {
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioExportProfileList, {});
    if (zdjete) return;
    profile = wynik.udany && wynik.wynik !== undefined
      ? { rodzaj: 'wykaz', profile: wynik.wynik.profiles }
      : { rodzaj: 'odmowa', blad: wynik.blad };
    odswiez();
  }

  async function wczytajUklad(documentId: string): Promise<void> {
    const numer = ++epokaUkladu;
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioPageSetupGet, { documentId });
    if (zdjete || numer !== epokaUkladu) return;
    uklad = wynik.udany && wynik.wynik !== undefined
      ? { rodzaj: 'znany', pageSetup: wynik.wynik.pageSetup }
      : { rodzaj: 'odmowa', blad: wynik.blad };
    odswiez();
  }

  async function wykonajRender(documentId: string): Promise<void> {
    const numer = ++epokaRenderu;
    render = { rodzaj: 'renderowanie' };
    odswiez();
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioPreviewRender, {
      documentId,
      format: ui.format,
      profileId: ui.profileId,
      windowId: zaleznosci.idOkna ?? undefined,
    });
    if (zdjete || numer !== epokaRenderu) return;
    render = wynik.udany && wynik.wynik !== undefined
      ? { rodzaj: 'gotowy', wynik: wynik.wynik }
      : { rodzaj: 'odmowa', blad: wynik.blad };
    ui.strona = 1;
    odswiez();
  }

  function ustawFormat(f: string): void {
    ui.format = f;
    if (dokument.rodzaj === 'obecny') void wykonajRender(dokument.dokument.id);
  }

  function ustawProfil(id: string): void {
    if (id === '') {
      ui.profileId = undefined;
    } else {
      ui.profileId = id;
      const wybrany = profile.rodzaj === 'wykaz' ? profile.profile.find((p) => p.id === id) : undefined;
      if (wybrany !== undefined) ui.format = wybrany.format;
    }
    if (dokument.rodzaj === 'obecny') void wykonajRender(dokument.dokument.id);
    else odswiez();
  }

  function ustawTryb(t: TrybWidoku): void {
    ui.tryb = t;
    odswiez();
  }

  function ustawPowiekszenie(p: number): void {
    ui.powiekszenie = p;
    odswiez();
  }

  function ustawPodzielony(v: boolean): void {
    ui.podzielony = v;
    odswiez();
  }

  function idzDoStrony(n: number): void {
    ui.strona = n;
    odswiez();
  }

  function odswiez(): void {
    tresc.replaceChildren(...zawartosc());
  }

  function zawartosc(): HTMLElement[] {
    if (zaleznosci.idOkna === null) return [alertProste(trescPodgladu.brakOkna)];
    return [pasekGora(), obszarPodgladu(), pasekDol()];
  }

  function pasekGora(): HTMLElement {
    const selFormat = el('select', {
      klasa: 'dn-btn dn-btn--duch dn-btn--sm',
      'aria-label': trescPodgladu.format,
    }) as HTMLSelectElement;
    for (const f of FORMATY) selFormat.appendChild(el('option', { value: f, tekst: f.toUpperCase() }));
    selFormat.value = ui.format;
    selFormat.addEventListener('change', () => ustawFormat(selFormat.value));

    const btnStrona = el('button', {
      klasa: `dn-btn dn-btn--sm ${ui.tryb === 'strona' ? 'dn-btn--zarys' : 'dn-btn--duch'}`,
      type: 'button',
      tekst: trescPodgladu.strona,
      'aria-pressed': String(ui.tryb === 'strona'),
    });
    btnStrona.addEventListener('click', () => ustawTryb('strona'));

    const btnCiagly = el('button', {
      klasa: `dn-btn dn-btn--sm ${ui.tryb === 'ciagly' ? 'dn-btn--zarys' : 'dn-btn--duch'}`,
      type: 'button',
      tekst: trescPodgladu.ciagly,
      'aria-pressed': String(ui.tryb === 'ciagly'),
    });
    btnCiagly.addEventListener('click', () => ustawTryb('ciagly'));

    const selZoom = el('select', {
      klasa: 'dn-btn dn-btn--duch dn-btn--sm',
      'aria-label': trescPodgladu.powiekszenie,
    }) as HTMLSelectElement;
    for (const p of POWIEKSZENIA) selZoom.appendChild(el('option', { value: p, tekst: `${p}%` }));
    selZoom.value = String(ui.powiekszenie);
    selZoom.addEventListener('change', () => ustawPowiekszenie(Number(selZoom.value)));

    return el('div', { klasa: 'st-panel-lista st-panel-lista--bez-dolu' }, [
      el('div', { klasa: 'st-panel-wiersz' }, [
        selFormat,
        btnStrona,
        btnCiagly,
        el('span', { klasa: 'dn-meta' }, [selZoom]),
      ]),
    ]);
  }

  function obszarPodgladu(): HTMLElement {
    if (dokument.rodzaj === 'brak') {
      return pustyStan(trescPodgladu.brakDokumentuTytul, trescPodgladu.brakDokumentuOpis);
    }
    if (render.rodzaj === 'bezczynny' || render.rodzaj === 'renderowanie') {
      return el('div', { klasa: 'st-podglad' }, [
        el('div', { klasa: 'dn-wykaz-modulu-poz' }, [
          el('span', { klasa: 'dn-kropka dn-kropka--sygnal dn-kropka--tetno', 'aria-hidden': 'true' }),
          trescPodgladu.renderowanie,
        ]),
      ]);
    }
    if (render.rodzaj === 'odmowa') {
      return el('div', { klasa: 'st-podglad' }, [alertBlad(trescPodgladu.odmowaRenderu, render.blad)]);
    }
    return widokGotowy(render.wynik);
  }

  function widokGotowy(wynik: StudioPreviewRenderResponse): HTMLElement {
    if (wynik.pages <= 0) return pustyStan(trescPodgladu.brakStronTytul, trescPodgladu.brakStronOpis);
    if (ui.strona > wynik.pages) ui.strona = wynik.pages;
    if (ui.strona < 1) ui.strona = 1;

    const indeksy =
      ui.tryb === 'ciagly'
        ? Array.from({ length: wynik.pages }, (_, i) => i + 1)
        : ui.podzielony && ui.strona < wynik.pages
          ? [ui.strona, ui.strona + 1]
          : [ui.strona];

    const kartki = el(
      'div',
      {
        klasa: ui.tryb === 'ciagly' ? 'st-panel-lista' : 'st-panel-wiersz',
        style: `transform: scale(${(ui.powiekszenie / 100).toFixed(2)}); transform-origin: top center;`,
      },
      indeksy.map((n) => kartka(n)),
    );

    const dzieci: HTMLElement[] = [...ukladNota(), kartki];
    if (ui.tryb === 'strona') dzieci.push(pasekStron(ui.strona, wynik.pages));

    return el('div', { klasa: 'st-podglad' }, dzieci);
  }

  function kartka(n: number): HTMLElement {
    return el('div', { klasa: 'dn-kartka' }, [el('p', { klasa: 'dn-kartka-numer', tekst: String(n) })]);
  }

  function ukladNota(): HTMLElement[] {
    if (uklad.rodzaj === 'znany') {
      const czesci = [uklad.pageSetup.pageSize, uklad.pageSetup.orientation].filter(
        (x): x is string => typeof x === 'string' && x.length > 0,
      );
      if (czesci.length === 0) return [];
      return [el('p', { klasa: 'dn-nota', tekst: czesci.join(' · ') })];
    }
    if (uklad.rodzaj === 'odmowa') return [alertBlad(trescPodgladu.odmowaUkladu, uklad.blad)];
    return [];
  }

  function pasekStron(strona: number, pages: number): HTMLElement {
    const naPoczatku = strona <= 1;
    const naKoncu = strona >= pages;

    const wstecz = el('button', {
      klasa: 'dn-btn dn-btn--duch dn-btn--sm',
      type: 'button',
      tekst: trescPodgladu.poprzednia,
      'aria-disabled': String(naPoczatku),
    });
    wstecz.addEventListener('click', () => {
      if (!naPoczatku) idzDoStrony(strona - 1);
    });

    const dalej = el('button', {
      klasa: 'dn-btn dn-btn--duch dn-btn--sm',
      type: 'button',
      tekst: trescPodgladu.nastepna,
      'aria-disabled': String(naKoncu),
    });
    dalej.addEventListener('click', () => {
      if (!naKoncu) idzDoStrony(strona + 1);
    });

    return el('div', { klasa: 'st-panel-wiersz' }, [
      wstecz,
      el('span', { klasa: 'dn-meta', tekst: `${strona} / ${pages}` }),
      dalej,
    ]);
  }

  function pasekDol(): HTMLElement {
    const wieleStron = render.rodzaj === 'gotowy' && render.wynik.pages > 1;
    const podzielNieczynny = dokument.rodzaj !== 'obecny' || ui.tryb !== 'strona' || !wieleStron;
    const btnPodziel = el('button', {
      klasa: `dn-btn dn-btn--sm ${ui.podzielony ? 'dn-btn--zarys' : 'dn-btn--duch'}`,
      type: 'button',
      tekst: trescPodgladu.podzielEkran,
      'aria-pressed': String(ui.podzielony),
      'aria-disabled': String(podzielNieczynny),
    });
    btnPodziel.addEventListener('click', () => {
      if (!podzielNieczynny) ustawPodzielony(!ui.podzielony);
    });

    const wierszProfilu: HTMLElement[] = [selektorProfilu()];
    if (profile.rodzaj === 'odmowa') {
      wierszProfilu.push(
        el('span', {
          klasa: 'dn-meta',
          tekst: `${trescPodgladu.profilOdmowa}: ${profile.blad?.message ?? trescPodgladu.brakOpisu}`,
        }),
      );
    } else if (profile.rodzaj === 'wykaz' && profile.profile.length === 0) {
      wierszProfilu.push(el('span', { klasa: 'dn-meta', tekst: trescPodgladu.profilBrak }));
    }

    const btnLibrary = el('button', {
      klasa: 'dn-btn dn-btn--duch dn-btn--sm',
      type: 'button',
      tekst: trescPodgladu.wyslijDoLibrary,
      'aria-disabled': 'true',
      title: trescPodgladu.wyslijNiegotowe,
    });

    return el('div', { klasa: 'st-panel-lista st-panel-lista--bez-gory' }, [
      el('div', { klasa: 'st-panel-wiersz' }, [btnPodziel, ...wierszProfilu, btnLibrary]),
    ]);
  }

  function selektorProfilu(): HTMLElement {
    const sel = el('select', {
      klasa: 'dn-btn dn-btn--atrament dn-btn--sm',
      'aria-label': trescPodgladu.eksportuj,
      disabled: profile.rodzaj !== 'wykaz' || profile.profile.length === 0,
    }) as HTMLSelectElement;

    if (profile.rodzaj === 'wczytywanie') {
      sel.appendChild(el('option', { value: '', tekst: trescPodgladu.profilWczytywanie }));
    } else {
      sel.appendChild(el('option', { value: '', tekst: trescPodgladu.profilDomyslny }));
      if (profile.rodzaj === 'wykaz') {
        for (const p of profile.profile) sel.appendChild(el('option', { value: p.id, tekst: p.name }));
      }
    }
    sel.value = ui.profileId ?? '';
    sel.addEventListener('change', () => ustawProfil(sel.value));
    return sel;
  }

  function pustyStan(tytul: string, opis: string): HTMLElement {
    const znak = zeZnacznika(ikony.pusto);
    znak.setAttribute('aria-hidden', 'true');
    return el('div', { klasa: 'st-podglad' }, [
      el('div', { klasa: 'dn-pusty-stan' }, [
        znak,
        el('span', { klasa: 'dn-pusty-stan-tytul', tekst: tytul }),
        el('span', { klasa: 'dn-pusty-stan-opis', tekst: opis }),
      ]),
    ]);
  }

  function alertBlad(tytul: string, blad: ErrorInfo | undefined): HTMLElement {
    return el('div', { klasa: 'dn-alert dn-alert--wstega dn-alert--blad', role: 'alert' }, [
      el('span', { klasa: 'dn-alert-tresc' }, [
        el('b', { tekst: tytul }),
        el('span', { tekst: blad?.message ?? trescPodgladu.brakOpisu }),
      ]),
    ]);
  }

  function alertProste(msg: string): HTMLElement {
    return el('div', { klasa: 'dn-alert dn-alert--wstega dn-alert--blad', role: 'alert' }, [
      el('span', { klasa: 'dn-alert-tresc' }, [el('b', { tekst: msg })]),
    ]);
  }

  return {
    zdejmij() {
      zdjete = true;
      odsubDokument?.();
      if (bryla.isConnected) bryla.remove();
    },
  };
};
