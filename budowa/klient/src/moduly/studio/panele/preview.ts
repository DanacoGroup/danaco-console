/**
 * Karta „Preview Window" okna roboczego Studia — podgląd wydania dokumentu:
 * układ i paginacja w formacie docelowym, nie sam tekst. Panel nie zna
 * dokumentu z góry — dostaje go reaktywnie zdarzeniem `studio.document.changed`,
 * którego opis w umowie wprost nazywa ten panel jako odbiorcę.
 *
 * Render (`studio.preview.render`) oddaje STRONY JAKO ZASOBY, nie tekst. Sam
 * wykaz kodów zasobów nie jest jeszcze podglądem, więc panel dobiera treść
 * każdej strony komendą `design.asset.content.get` — umowa opisuje ją jako
 * czynność KAŻDEGO zasobu magazynu, nie tylko rodziny design. Dopiero obraz
 * strony jest tym, co rdzeń naprawdę oddał; pusta kartka byłaby atrapą.
 *
 * Nastawy strony idą w obie strony: odczyt `studio.page.setup.get`, zmiana
 * nośnika i orientacji `studio.page.setup.set`, wykaz formatów nośnika
 * `studio.page.paper.list`. Profile wydania to `studio.export.profile.list`
 * i `studio.export.profile.save` — zapis bierze format i nastawy strony
 * z tego, co widać, więc profil odkłada stan podglądu, nie wartości zmyślone.
 *
 * Bez okna modułu żadna z tych czynności nie ma na czym stanąć, więc panel
 * wtedy nie próbuje żadnej.
 */

import {
  AssetContentDisposition,
  ChangeKind,
  Command,
  EventType,
  StudioExportFormat,
  StudioPageOrientation,
  type ErrorInfo,
  type StudioDocument,
  type StudioExportProfile,
  type StudioPageSetup,
  type StudioPaperFormat,
  type StudioPreviewRenderResponse,
} from '../../../../../shared/contract.ts';
import { wywolaj } from '../../../protokol/wywolanie.ts';
import { el, zeZnacznika, type Dziecko } from '../narzedzia.ts';
import { ikony } from '../ikony.ts';
import type { MontazPanelu } from './umowa.ts';
import { trescPodgladu as T } from './preview-tresci.ts';
import { opisOdmowy } from '../odmowa.ts';

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
  | { rodzaj: 'zmieniany' }
  | { rodzaj: 'odmowa'; etykieta: string; blad: ErrorInfo | undefined };

type StanNosnikow =
  | { rodzaj: 'wczytywanie' }
  | { rodzaj: 'wykaz'; formaty: StudioPaperFormat[] }
  | { rodzaj: 'odmowa'; blad: ErrorInfo | undefined };

type StanObrazu =
  | { rodzaj: 'pobiera' }
  | { rodzaj: 'gotowy'; zrodlo: string }
  | { rodzaj: 'odmowa'; blad: ErrorInfo | undefined };

type StanZapisu =
  | { rodzaj: 'ukryty' }
  | { rodzaj: 'formularz'; ostrzezenie: string | undefined }
  | { rodzaj: 'biegnie' }
  | { rodzaj: 'odmowa'; blad: ErrorInfo | undefined };

const POWIEKSZENIA: readonly number[] = [50, 75, 100, 125, 150, 200];
const FORMATY: readonly string[] = Object.values(StudioExportFormat);

/* Szerokość kartki w arkuszu prototypu to 60ch; powiększenie mnoży tę miarę,
   zamiast skalować węzeł — skala zostawiłaby pasek przewijania bez pokrycia. */
const SZEROKOSC_KARTKI_CH = 60;

/* Znak zamknięcia stoi tutaj, a nie w `ikony.ts`: ten plik leży poza terenem
   zmiany. Rysunek przeniesiony z prototypu bez zmian. */
const ZNAK_ZAMKNIECIA =
  '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" ' +
  'stroke-linecap="round" stroke-linejoin="round"><path d="M18 6 6 18M6 6l12 12"/></svg>';

function znak(rysunek: keyof typeof ikony): SVGElement {
  const wezelZnaku = zeZnacznika(ikony[rysunek]);
  wezelZnaku.setAttribute('aria-hidden', 'true');
  return wezelZnaku;
}

function podstaw(wzorzec: string, wartosci: Record<string, string>): string {
  return wzorzec.replace(/\{(\w+)\}/g, (calosc, klucz: string) => wartosci[klucz] ?? calosc);
}

export const montujPodgladWydania: MontazPanelu = (wezel, zaleznosci) => {
  let zdjete = false;
  let epokaRenderu = 0;
  let epokaUkladu = 0;
  let odswiezanieZaplanowane = false;

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
  let nosniki: StanNosnikow = { rodzaj: 'wczytywanie' };
  let zapis: StanZapisu = { rodzaj: 'ukryty' };
  let nazwaProfilu = '';

  /* Treść stron trzymana pod kodem zasobu, żeby przewracanie stron i zmiana
     powiększenia nie wołały rdzenia po raz drugi o to samo. */
  const obrazy = new Map<string, StanObrazu>();

  /* Węzeł zaczepu jest już bryłą okna z identyfikatorem, rolą i ukryciem —
     panel dokłada do niego belkę i treść, zamiast zakładać drugą bryłę. */
  wezel.classList.add('sta-okno');
  const tresc = el('div', { klasa: 'sta-okno-tresc' });
  /* Przycisk zamknięcia w belce niesie ten sam znacznik co znak w paśmie kart,
     bo to pasmo prowadzi wykaz kart i ono zamyka kartę tego okna. */
  const zamknijKarte = el(
    'button',
    {
      klasa: 'dn-btn-ikona',
      type: 'button',
      'data-karta-zamknij': true,
      'aria-label': T.panel.zamknijKarte,
    },
    [zeZnacznika(ZNAK_ZAMKNIECIA)],
  );

  wezel.append(
    el('header', { klasa: 'sta-okno-belka' }, [
      el('span', { klasa: 'sta-okno-tytul' }, [znak('oko'), el('b', { tekst: T.panel.tytul })]),
      el('span', { klasa: 'sta-okno-akcje' }, [zamknijKarte]),
    ]),
    tresc,
  );

  let odsubDokument: (() => void) | null = null;

  if (zaleznosci.idOkna !== null) {
    const idOkna = zaleznosci.idOkna;
    void wczytajProfile();
    void wczytajNosniki();
    odsubDokument = zaleznosci.kanal.naZdarzenie(EventType.StudioDocumentChanged, (zdarzenie) => {
      if (zdjete || zdarzenie.document.windowId !== idOkna) return;
      if (zdarzenie.change === ChangeKind.Deleted) {
        dokument = { rodzaj: 'brak' };
        render = { rodzaj: 'bezczynny' };
        uklad = { rodzaj: 'brak' };
        obrazy.clear();
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

  async function wczytajNosniki(): Promise<void> {
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioPagePaperList, {});
    if (zdjete) return;
    nosniki = wynik.udany && wynik.wynik !== undefined
      ? { rodzaj: 'wykaz', formaty: wynik.wynik.papers }
      : { rodzaj: 'odmowa', blad: wynik.blad };
    odswiez();
  }

  async function wczytajUklad(documentId: string): Promise<void> {
    const numer = ++epokaUkladu;
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioPageSetupGet, { documentId });
    if (zdjete || numer !== epokaUkladu) return;
    uklad = wynik.udany && wynik.wynik !== undefined
      ? { rodzaj: 'znany', pageSetup: wynik.wynik.pageSetup }
      : { rodzaj: 'odmowa', etykieta: T.odmowaUkladu, blad: wynik.blad };
    odswiez();
  }

  /* Zmiana nastaw strony przelicza układ w rdzeniu, więc po niej podgląd musi
     powstać od nowa — inaczej pokazywałby stronę sprzed zmiany. */
  async function zmienUklad(zmiana: { paperName?: string; orientation?: StudioPageOrientation }): Promise<void> {
    if (dokument.rodzaj !== 'obecny') return;
    const documentId = dokument.dokument.id;
    const numer = ++epokaUkladu;
    uklad = { rodzaj: 'zmieniany' };
    odswiez();
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioPageSetupSet, { documentId, ...zmiana });
    if (zdjete || numer !== epokaUkladu) return;
    if (wynik.udany && wynik.wynik !== undefined) {
      uklad = { rodzaj: 'znany', pageSetup: wynik.wynik.pageSetup };
      void wykonajRender(documentId);
      return;
    }
    uklad = { rodzaj: 'odmowa', etykieta: T.odmowaZmianyUkladu, blad: wynik.blad };
    odswiez();
  }

  async function wykonajRender(documentId: string): Promise<void> {
    const numer = ++epokaRenderu;
    render = { rodzaj: 'renderowanie' };
    obrazy.clear();
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

  /* Zasób strony pobierany raz i tylko na żądanie widoku: w trybie ciągłym
     stron bywa wiele, a Operator ogląda je po kolei. */
  function zadajObraz(assetId: string): void {
    if (obrazy.has(assetId)) return;
    obrazy.set(assetId, { rodzaj: 'pobiera' });
    void (async () => {
      const wynik = await wywolaj(zaleznosci.kanal, Command.DesignAssetContentGet, {
        assetId,
        disposition: AssetContentDisposition.Inline,
      });
      if (zdjete) return;
      const oddane = wynik.wynik;
      const zapisTresci = oddane?.contentBase64 ?? '';
      obrazy.set(
        assetId,
        wynik.udany && oddane !== undefined && zapisTresci !== ''
          ? { rodzaj: 'gotowy', zrodlo: `data:${oddane.mediaType};base64,${zapisTresci}` }
          : { rodzaj: 'odmowa', blad: wynik.blad },
      );
      zaplanujOdswiez();
    })();
  }

  /* Strony wracają osobnymi odpowiedziami; przebudowa widoku po każdej z nich
     z osobna kosztowałaby tyle razy, ile stron. */
  function zaplanujOdswiez(): void {
    if (odswiezanieZaplanowane) return;
    odswiezanieZaplanowane = true;
    queueMicrotask(() => {
      odswiezanieZaplanowane = false;
      if (!zdjete) odswiez();
    });
  }

  async function zapiszProfil(): Promise<void> {
    const nazwa = nazwaProfilu.trim();
    if (nazwa === '') {
      zapis = { rodzaj: 'formularz', ostrzezenie: T.zapiszProfilBezNazwy };
      odswiez();
      return;
    }
    zapis = { rodzaj: 'biegnie' };
    odswiez();
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioExportProfileSave, {
      name: nazwa,
      format: ui.format,
      pageSetup: uklad.rodzaj === 'znany' ? uklad.pageSetup : undefined,
    });
    if (zdjete) return;
    if (wynik.udany && wynik.wynik !== undefined) {
      ui.profileId = wynik.wynik.profile.id;
      nazwaProfilu = '';
      zapis = { rodzaj: 'ukryty' };
      profile = { rodzaj: 'wczytywanie' };
      odswiez();
      void wczytajProfile();
      return;
    }
    zapis = { rodzaj: 'odmowa', blad: wynik.blad };
    odswiez();
  }

  function ustawFormat(f: string): void {
    ui.format = f;
    if (dokument.rodzaj === 'obecny') void wykonajRender(dokument.dokument.id);
    else odswiez();
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

  function odswiez(): void {
    tresc.replaceChildren(...zawartosc());
  }

  function zawartosc(): HTMLElement[] {
    if (zaleznosci.idOkna === null) return [alertProste(T.brakOkna)];
    return [pasekGora(), obszarPodgladu(), pasekDol()];
  }

  // ── Pasek górny: format wydania, tryb widoku, powiększenie, nastawy strony ──

  function pasekGora(): HTMLElement {
    const selFormat = el('select', {
      klasa: 'dn-btn dn-btn--duch dn-btn--sm',
      'aria-label': T.format,
    }) as HTMLSelectElement;
    for (const f of FORMATY) selFormat.appendChild(el('option', { value: f, tekst: f.toUpperCase() }));
    selFormat.value = ui.format;
    selFormat.addEventListener('change', () => ustawFormat(selFormat.value));

    const btnStrona = przelacznik(T.strona, ui.tryb === 'strona', () => {
      ui.tryb = 'strona';
      odswiez();
    });
    const btnCiagly = przelacznik(T.ciagly, ui.tryb === 'ciagly', () => {
      ui.tryb = 'ciagly';
      odswiez();
    });

    const selZoom = el('select', {
      klasa: 'dn-btn dn-btn--duch dn-btn--sm',
      'aria-label': T.powiekszenie,
    }) as HTMLSelectElement;
    for (const p of POWIEKSZENIA) selZoom.appendChild(el('option', { value: p, tekst: `${p}%` }));
    selZoom.value = String(ui.powiekszenie);
    selZoom.addEventListener('change', () => {
      ui.powiekszenie = Number(selZoom.value);
      odswiez();
    });

    return el('div', { klasa: 'st-panel-lista st-panel-lista--bez-dolu' }, [
      el('div', { klasa: 'st-panel-wiersz' }, [
        selFormat,
        btnStrona,
        btnCiagly,
        el('span', { klasa: 'dn-meta' }, [selZoom]),
      ]),
      wierszUkladu(),
    ]);
  }

  function przelacznik(etykieta: string, wlaczony: boolean, naKlik: () => void): HTMLElement {
    const przycisk = el('button', {
      klasa: `dn-btn dn-btn--sm ${wlaczony ? 'dn-btn--zarys' : 'dn-btn--duch'}`,
      type: 'button',
      tekst: etykieta,
      'aria-pressed': String(wlaczony),
    });
    przycisk.addEventListener('click', naKlik);
    return przycisk;
  }

  function wierszUkladu(): HTMLElement {
    const nieczynny = dokument.rodzaj !== 'obecny' || uklad.rodzaj === 'zmieniany';

    const selNosnik = el('select', {
      klasa: 'dn-btn dn-btn--duch dn-btn--sm',
      'aria-label': T.nosnik,
      disabled: nieczynny || nosniki.rodzaj !== 'wykaz',
    }) as HTMLSelectElement;
    if (nosniki.rodzaj === 'wczytywanie') {
      selNosnik.appendChild(el('option', { value: '', tekst: T.nosnikWczytywanie }));
    } else {
      selNosnik.appendChild(el('option', { value: '', tekst: T.nosnikBrak }));
      if (nosniki.rodzaj === 'wykaz') {
        for (const f of nosniki.formaty) {
          selNosnik.appendChild(
            el('option', { value: f.name, tekst: `${f.name} · ${f.widthMm}×${f.heightMm} mm` }),
          );
        }
      }
    }
    selNosnik.value = uklad.rodzaj === 'znany' ? uklad.pageSetup.pageSize ?? '' : '';
    selNosnik.addEventListener('change', () => {
      if (selNosnik.value !== '') void zmienUklad({ paperName: selNosnik.value });
    });

    const selOrientacja = el('select', {
      klasa: 'dn-btn dn-btn--duch dn-btn--sm',
      'aria-label': T.orientacja,
      disabled: nieczynny,
    }) as HTMLSelectElement;
    selOrientacja.append(
      el('option', { value: StudioPageOrientation.Pionowa, tekst: T.pionowa }),
      el('option', { value: StudioPageOrientation.Pozioma, tekst: T.pozioma }),
    );
    selOrientacja.value =
      uklad.rodzaj === 'znany' ? uklad.pageSetup.orientation ?? StudioPageOrientation.Pionowa : StudioPageOrientation.Pionowa;
    selOrientacja.addEventListener('change', () => {
      void zmienUklad({ orientation: selOrientacja.value as StudioPageOrientation });
    });

    const dodatki: Dziecko[] = [];
    if (uklad.rodzaj === 'zmieniany') dodatki.push(el('span', { klasa: 'dn-meta', tekst: T.ukladZmieniany }));
    if (uklad.rodzaj === 'odmowa') dodatki.push(notaOdmowy(uklad.etykieta, uklad.blad, 'preview.uklad'));
    if (nosniki.rodzaj === 'odmowa') dodatki.push(notaOdmowy(T.odmowaNosnikow, nosniki.blad, 'preview.nosniki'));

    return el('div', { klasa: 'st-panel-wiersz' }, [selNosnik, selOrientacja, ...dodatki]);
  }

  // ── Obszar kartek ───────────────────────────────────────────────────────

  function obszarPodgladu(): HTMLElement {
    if (dokument.rodzaj === 'brak') {
      return pustyStan(T.brakDokumentuTytul, T.brakDokumentuOpis);
    }
    if (render.rodzaj === 'bezczynny' || render.rodzaj === 'renderowanie') {
      return el('div', { klasa: 'st-podglad' }, [
        el('div', { klasa: 'dn-wykaz-modulu-poz' }, [
          el('span', { klasa: 'dn-kropka dn-kropka--sygnal dn-kropka--tetno', 'aria-hidden': 'true' }),
          T.renderowanie,
        ]),
      ]);
    }
    if (render.rodzaj === 'odmowa') {
      return el('div', { klasa: 'st-podglad' }, [alertBlad(T.odmowaRenderu, render.blad)]);
    }
    return widokGotowy(render.wynik);
  }

  function widokGotowy(wynik: StudioPreviewRenderResponse): HTMLElement {
    if (wynik.pages <= 0) return pustyStan(T.brakStronTytul, T.brakStronOpis);
    if (ui.strona > wynik.pages) ui.strona = wynik.pages;
    if (ui.strona < 1) ui.strona = 1;

    const indeksy =
      ui.tryb === 'ciagly'
        ? Array.from({ length: wynik.pages }, (_, i) => i + 1)
        : ui.podzielony && ui.strona < wynik.pages
          ? [ui.strona, ui.strona + 1]
          : [ui.strona];

    /* Pojedynczy arkusz leży wprost w obszarze podglądu — tak stoi w prototypie.
       Zawijarka wchodzi dopiero przy wielu arkuszach: obok siebie przy podziale
       ekranu, jeden pod drugim w trybie ciągłym. */
    const arkusze = indeksy.map((n) => kartka(n, wynik.pageAssetIds[n - 1]));
    const dzieci: Dziecko[] =
      arkusze.length === 1
        ? [...arkusze]
        : [el('div', { klasa: ui.tryb === 'ciagly' ? 'st-panel-lista' : 'st-panel-wiersz' }, arkusze)];
    if (ui.tryb === 'strona') dzieci.push(pasekStron(ui.strona, wynik.pages));

    return el('div', { klasa: 'st-podglad' }, dzieci);
  }

  /* Kartka niesie to, co rdzeń wyrysował: obraz strony z magazynu. Tytuł stoi
     tylko przy pierwszej kartce, żeby nie powtarzać go nad każdą stroną. */
  function kartka(numer: number, assetId: string | undefined): HTMLElement {
    const szerokosc = ((SZEROKOSC_KARTKI_CH * ui.powiekszenie) / 100).toFixed(1);
    const dzieci: Dziecko[] = [];

    if (numer === 1 && dokument.rodzaj === 'obecny') {
      const tytul = dokument.dokument.title ?? '';
      dzieci.push(el('div', { klasa: 'dn-kartka-tytul', tekst: tytul === '' ? T.bezTytulu : tytul }));
    }
    dzieci.push(trescStrony(numer, assetId));
    dzieci.push(el('p', { klasa: 'dn-kartka-numer', tekst: podstaw(T.numerStrony, { numer: String(numer) }) }));

    return el('div', { klasa: 'dn-kartka', style: `max-width: ${szerokosc}ch` }, dzieci);
  }

  function trescStrony(numer: number, assetId: string | undefined): HTMLElement {
    if (assetId === undefined) return el('p', { klasa: 'dn-nota', tekst: T.stronaBezObrazu });
    zadajObraz(assetId);
    const stan = obrazy.get(assetId);
    if (stan === undefined || stan.rodzaj === 'pobiera') {
      return el('p', { klasa: 'dn-nota', tekst: T.stronaWczytywanie });
    }
    if (stan.rodzaj === 'odmowa') return notaOdmowy(T.odmowaStrony, stan.blad, 'preview.strona');
    return el('img', {
      src: stan.zrodlo,
      alt: podstaw(T.numerStrony, { numer: String(numer) }),
      style: 'display: block; width: 100%; height: auto;',
    });
  }

  function pasekStron(strona: number, pages: number): HTMLElement {
    const naPoczatku = strona <= 1;
    const naKoncu = strona >= pages;

    const wstecz = el('button', {
      klasa: 'dn-btn dn-btn--duch dn-btn--sm',
      type: 'button',
      tekst: T.poprzednia,
      'aria-disabled': String(naPoczatku),
    });
    wstecz.addEventListener('click', () => {
      if (naPoczatku) return;
      ui.strona = strona - 1;
      odswiez();
    });

    const dalej = el('button', {
      klasa: 'dn-btn dn-btn--duch dn-btn--sm',
      type: 'button',
      tekst: T.nastepna,
      'aria-disabled': String(naKoncu),
    });
    dalej.addEventListener('click', () => {
      if (naKoncu) return;
      ui.strona = strona + 1;
      odswiez();
    });

    return el('div', { klasa: 'st-panel-wiersz' }, [
      wstecz,
      el('span', {
        klasa: 'dn-meta',
        tekst: podstaw(T.stronaZe, { numer: String(strona), ile: String(pages) }),
      }),
      dalej,
    ]);
  }

  // ── Pasek dolny: podział ekranu, profil wydania, wysyłka ────────────────

  function pasekDol(): HTMLElement {
    const wieleStron = render.rodzaj === 'gotowy' && render.wynik.pages > 1;
    const podzielNieczynny = dokument.rodzaj !== 'obecny' || ui.tryb !== 'strona' || !wieleStron;
    const btnPodziel = el('button', {
      klasa: `dn-btn dn-btn--sm ${ui.podzielony ? 'dn-btn--zarys' : 'dn-btn--duch'}`,
      type: 'button',
      tekst: T.podzielEkran,
      'aria-pressed': String(ui.podzielony),
      'aria-disabled': String(podzielNieczynny),
    });
    btnPodziel.addEventListener('click', () => {
      if (podzielNieczynny) return;
      ui.podzielony = !ui.podzielony;
      odswiez();
    });

    const wiersz: Dziecko[] = [btnPodziel, selektorProfilu()];
    if (profile.rodzaj === 'odmowa') wiersz.push(notaOdmowy(T.profilOdmowa, profile.blad, 'preview.profile'));
    else if (profile.rodzaj === 'wykaz' && profile.profile.length === 0) {
      wiersz.push(el('span', { klasa: 'dn-meta', tekst: T.profilBrak }));
    }

    const btnZapisz = el('button', {
      klasa: 'dn-btn dn-btn--duch dn-btn--sm',
      type: 'button',
      tekst: T.zapiszProfil,
      'aria-disabled': String(zapis.rodzaj === 'biegnie'),
    });
    btnZapisz.addEventListener('click', () => {
      if (zapis.rodzaj === 'biegnie') return;
      zapis = zapis.rodzaj === 'formularz' ? { rodzaj: 'ukryty' } : { rodzaj: 'formularz', ostrzezenie: undefined };
      odswiez();
    });
    wiersz.push(btnZapisz);

    /* Czynność bez pokrycia w umowie z rdzeniem: zostaje widoczna i mówi
       wprost, że nie jest gotowa, zamiast udawać przycisk bez skutku. */
    const btnLibrary = el('button', {
      klasa: 'dn-btn dn-btn--duch dn-btn--sm',
      type: 'button',
      tekst: T.wyslijDoLibrary,
      'aria-disabled': 'true',
      'aria-describedby': 'st-podglad-library-nota',
    });
    wiersz.push(btnLibrary);

    const dzieci: Dziecko[] = [el('div', { klasa: 'st-panel-wiersz' }, wiersz)];
    const formularz = wierszZapisuProfilu();
    if (formularz !== null) dzieci.push(formularz);
    dzieci.push(
      el('div', { klasa: 'st-panel-wiersz' }, [
        el('span', { klasa: 'dn-nota', id: 'st-podglad-library-nota', tekst: T.wyslijNiegotowe }),
      ]),
    );

    return el('div', { klasa: 'st-panel-lista st-panel-lista--bez-gory' }, dzieci);
  }

  function selektorProfilu(): HTMLElement {
    const sel = el('select', {
      klasa: 'dn-btn dn-btn--atrament dn-btn--sm',
      'aria-label': T.eksportuj,
      disabled: profile.rodzaj !== 'wykaz' || profile.profile.length === 0,
    }) as HTMLSelectElement;

    if (profile.rodzaj === 'wczytywanie') {
      sel.appendChild(el('option', { value: '', tekst: T.profilWczytywanie }));
    } else {
      sel.appendChild(el('option', { value: '', tekst: T.profilDomyslny }));
      if (profile.rodzaj === 'wykaz') {
        for (const p of profile.profile) sel.appendChild(el('option', { value: p.id, tekst: p.name }));
      }
    }
    sel.value = ui.profileId ?? '';
    sel.addEventListener('change', () => ustawProfil(sel.value));
    return sel;
  }

  function wierszZapisuProfilu(): HTMLElement | null {
    if (zapis.rodzaj === 'ukryty') return null;
    if (zapis.rodzaj === 'biegnie') {
      return el('div', { klasa: 'st-panel-wiersz' }, [
        el('span', { klasa: 'dn-meta', tekst: T.zapiszProfilBiegnie }),
      ]);
    }
    if (zapis.rodzaj === 'odmowa') {
      return el('div', { klasa: 'st-panel-wiersz' }, [
        notaOdmowy(T.zapiszProfilOdmowa, zapis.blad, 'preview.profil.zapis'),
      ]);
    }

    /* Nazwa żyje poza widokiem, bo panel przebudowuje treść przy każdej
       zmianie stanu — trzymana w węźle, znikałaby przy pierwszym odświeżeniu. */
    const pole = el('input', {
      klasa: 'dn-pole-kontrolka',
      type: 'text',
      value: nazwaProfilu,
      placeholder: T.zapiszProfilZnak,
      'aria-label': T.zapiszProfilNazwa,
    }) as HTMLInputElement;
    pole.addEventListener('input', () => {
      nazwaProfilu = pole.value;
    });

    const zatwierdz = el('button', {
      klasa: 'dn-btn dn-btn--atrament dn-btn--sm',
      type: 'button',
      tekst: T.zapiszProfilZatwierdz,
    });
    zatwierdz.addEventListener('click', () => void zapiszProfil());

    const porzuc = el('button', {
      klasa: 'dn-btn dn-btn--duch dn-btn--sm',
      type: 'button',
      tekst: T.zapiszProfilPorzuc,
    });
    porzuc.addEventListener('click', () => {
      zapis = { rodzaj: 'ukryty' };
      odswiez();
    });

    const dzieci: Dziecko[] = [pole, zatwierdz, porzuc];
    if (zapis.ostrzezenie !== undefined) {
      dzieci.push(el('span', { klasa: 'dn-nota', tekst: zapis.ostrzezenie }));
    }
    return el('div', { klasa: 'st-panel-wiersz' }, dzieci);
  }

  // ── Stany wspólne ───────────────────────────────────────────────────────

  function pustyStan(tytul: string, opis: string): HTMLElement {
    return el('div', { klasa: 'st-podglad' }, [
      el('div', { klasa: 'dn-pusty-stan' }, [
        znak('pusto'),
        el('span', { klasa: 'dn-pusty-stan-tytul', tekst: tytul }),
        el('span', { klasa: 'dn-pusty-stan-opis', tekst: opis }),
      ]),
    ]);
  }

  function alertBlad(tytul: string, blad: ErrorInfo | undefined): HTMLElement {
    return el('div', { klasa: 'dn-alert dn-alert--wstega dn-alert--blad', role: 'alert' }, [
      el('span', { klasa: 'dn-alert-tresc' }, [
        el('b', { tekst: tytul }),
        el('span', { tekst: opisOdmowy(blad, 'preview') }),
      ]),
    ]);
  }

  function alertProste(msg: string): HTMLElement {
    return el('div', { klasa: 'dn-alert dn-alert--wstega dn-alert--blad', role: 'alert' }, [
      el('span', { klasa: 'dn-alert-tresc' }, [el('b', { tekst: msg })]),
    ]);
  }

  function notaOdmowy(etykieta: string, blad: ErrorInfo | undefined, czynnosc: string): HTMLElement {
    return el('span', { klasa: 'dn-nota', role: 'alert', tekst: `${etykieta}. ${opisOdmowy(blad, czynnosc)}` });
  }

  return {
    zdejmij() {
      zdjete = true;
      odsubDokument?.();
      obrazy.clear();
      wezel.replaceChildren();
    },
  };
};
