/**
 * Diff/Grep Panel — strefa `section#panel-diff` z prototypu
 * (`design/05-okna/moduly/studio.html`). Porównanie dwóch wersji dokumentu
 * Studia w czterech zakresach — treść scalona, postać, wygląd i materiał
 * wejściowy — wyszukiwanie wzorca, adnotacje przy fragmentach, przeniesienie
 * fragmentu, zmiany odłożone przez śledzenie wraz z decyzją Operatora oraz
 * raport zmian odkładany w magazynie.
 *
 * Panel nie zakłada okna ani dokumentu — czeka na `studio.document.changed`
 * dla okna wskazanego w zależnościach i dopiero wtedy woła rodziny komend
 * `studio.diff.*`, `studio.tracking.*` i `studio.annotation.*`. Zdarzenie
 * sprzed zamontowania panelu nie dociera do niego: kontrakt nie daje dziś drogi
 * do odpytania o dokument bieżący okna spoza tych trzech rodzin.
 *
 * Wykaz wersji prowadzi rodzina `studio.repository.*` — poza tym panelem —
 * więc miejsce, w którym prototyp stawia listę wersji, niesie pole
 * identyfikatora i zdanie mówiące, skąd wykaz wziąć.
 */

import {
  ChangeKind,
  Command,
  DiffHunkKind,
  EventType,
  StudioAuthor,
  StudioChangeKind,
  StudioDocumentFormat,
  type ErrorInfo,
  type StudioAnnotation,
  type StudioDiffHunk,
  type StudioFormDiffEntry,
  type StudioTextMatch,
  type StudioTrackedChange,
  type StudioVisualDiffRegion,
} from '../../../../../shared/contract.ts';
import { wywolaj } from '../../../protokol/wywolanie.ts';
import { ikony } from '../ikony.ts';
import { el, tekst, zeZnacznika, type Dziecko } from '../narzedzia.ts';
import { tresc as t } from './diff-tresci.ts';
import type { MontazPanelu } from './umowa.ts';
import { opisOdmowy } from '../odmowa.ts';
import { zLiczba } from '../liczebnik.ts';

/** Zakres porównania — prototypowe „scalony ▾”, każdy zakres własną komendą rodziny `studio.diff.*`. */
type Tryb = 'scalony' | 'postac' | 'wyglad' | 'zrodlo';

/** Wynik porównania w zakresie wybranym przez Operatora; trafienia wzorca wracają tą samą odpowiedzią co fragmenty. */
type ZasobPorownania =
  | { rodzaj: 'ladowanie' }
  | { rodzaj: 'odmowa'; blad?: ErrorInfo }
  | { rodzaj: 'brakDanych'; zdanie: string }
  | { rodzaj: 'tekst'; hunks: StudioDiffHunk[]; trafienia: StudioTextMatch[] | null }
  | { rodzaj: 'postac'; wpisy: StudioFormDiffEntry[]; dodane: number; usuniete: number; zmienione: number }
  | { rodzaj: 'wyglad'; obszary: StudioVisualDiffRegion[]; nakladki: number }
  | { rodzaj: 'zrodlo'; hunks: StudioDiffHunk[]; odczytane: boolean };

/** Wykaz zmian oczekujących na decyzję Operatora. */
type ZasobZmian =
  | { rodzaj: 'ladowanie' }
  | { rodzaj: 'odmowa'; blad?: ErrorInfo }
  | { rodzaj: 'dane'; zmiany: StudioTrackedChange[] };

/** Adnotacje założone przy fragmentach bieżącego porównania. */
type ZasobAdnotacji =
  | { rodzaj: 'ladowanie' }
  | { rodzaj: 'odmowa'; blad?: ErrorInfo }
  | { rodzaj: 'dane'; adnotacje: StudioAnnotation[] };

/** Stan wydania raportu zmian. */
type ZasobRaportu =
  | { rodzaj: 'wydaje' }
  | { rodzaj: 'odmowa'; blad?: ErrorInfo }
  | { rodzaj: 'gotowy'; nazwa: string; bajty: number };

const TRYBY: readonly { kod: Tryb; nazwa: string }[] = [
  { kod: 'scalony', nazwa: t.tryb.scalony },
  { kod: 'postac', nazwa: t.tryb.postac },
  { kod: 'wyglad', nazwa: t.tryb.wyglad },
  { kod: 'zrodlo', nazwa: t.tryb.zrodlo },
];

const FORMATY_RAPORTU: readonly { kod: string; nazwa: string }[] = [
  { kod: StudioDocumentFormat.Pdf, nazwa: t.raport.formaty.pdf },
  { kod: StudioDocumentFormat.Docx, nazwa: t.raport.formaty.docx },
  { kod: StudioDocumentFormat.Txt, nazwa: t.raport.formaty.txt },
  { kod: StudioDocumentFormat.Markdown, nazwa: t.raport.formaty.markdown },
];

function znak(rysunek: keyof typeof ikony): SVGElement {
  const wezelZnaku = zeZnacznika(ikony[rysunek]);
  wezelZnaku.setAttribute('aria-hidden', 'true');
  return wezelZnaku;
}

/** Znaczek odmowy jednakowy dla każdej czynności panelu: tytuł nazywa czynność, treść — zdanie dobrane po kodzie odmowy. */
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

function pustyStan(zdanie: string): HTMLElement {
  return el('p', { klasa: 'dn-pusty-stan dn-pusty-stan--zwarty', tekst: zdanie });
}

function nota(zdanie: string): HTMLElement {
  return el('div', { klasa: 'dn-nota', tekst: zdanie });
}

/** Lista wyboru z prototypowego miejsca „▾”: wartość wskazuje Operator, nie zgaduje jej panel. */
function listaWyboru(
  etykieta: string,
  pozycje: readonly { kod: string; nazwa: string }[],
  wybrana: string,
  naZmiane: (kod: string) => void,
): HTMLSelectElement {
  const wybor = el('select', { klasa: 'dn-pole-kontrolka', 'aria-label': etykieta }) as HTMLSelectElement;
  for (const pozycja of pozycje) {
    wybor.appendChild(el('option', { value: pozycja.kod, tekst: pozycja.nazwa }));
  }
  wybor.value = wybrana;
  wybor.addEventListener('change', () => naZmiane(wybor.value));
  return wybor;
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
  let tryb: Tryb = 'scalony';
  let formatRaportu: string = StudioDocumentFormat.Pdf;

  let porownanie: ZasobPorownania | null = null;
  let zmianyStan: ZasobZmian | null = null;
  let adnotacjeStan: ZasobAdnotacji | null = null;
  let raportStan: ZasobRaportu | null = null;

  /* Rdzeń podaje stan śledzenia wyłącznie w odpowiedzi na jego przestawienie,
     więc przed pierwszym przestawieniem stan jest nieznany i pole wyboru stoi
     nierozstrzygnięte, zamiast udawać, że śledzenie jest wyłączone. */
  let sledzenie: boolean | null = null;
  let sledzenieWTrakcie = false;
  let bladSledzenia: ErrorInfo | undefined;

  let adnotowanyHunk: number | null = null;
  let trescAdnotacji = '';
  let adnotacjaWTrakcie = false;
  const bladAdnotacji = new Map<number, ErrorInfo | undefined>();

  let przenoszonyHunk: number | null = null;
  const bledyPrzeniesienia = new Map<number, ErrorInfo | undefined>();
  let pominietePrzeniesienia: number | null = null;

  const decyzjeWTrakcie = new Set<string>();
  const bledyDecyzji = new Map<string, ErrorInfo | undefined>();

  const trescPanelu = el('div', { klasa: 'sta-okno-tresc st-panel-lista' });

  function ustawDokument(id: string): void {
    if (id === dokumentId) return;
    dokumentId = id;
    pokolenie += 1;
    porownanie = null;
    zmianyStan = null;
    adnotacjeStan = null;
    raportStan = null;
    sledzenie = null;
    void uruchomPorownanie();
    void wczytajZmiany();
  }

  function usunDokument(): void {
    dokumentId = null;
    pokolenie += 1;
    porownanie = null;
    zmianyStan = null;
    adnotacjeStan = null;
    raportStan = null;
    odswiez();
  }

  /** Porównanie w zakresie treści scalonej wraz z wyszukaniem wzorca — `studio.diff.compare`. */
  async function porownajTresc(dokument: string): Promise<ZasobPorownania> {
    const wzorzecPrzyciety = wzorzec.trim();
    const wynik = await wywolaj(zal.kanal, Command.StudioDiffCompare, {
      documentId: dokument,
      baseVersionId: wersjaBazowa.trim() || undefined,
      targetVersionId: wersjaCel.trim() || undefined,
      pattern: wzorzecPrzyciety === '' ? undefined : wzorzecPrzyciety,
      regex: wzorzecPrzyciety === '' ? undefined : regex,
    });
    if (!wynik.udany || wynik.wynik === undefined) return { rodzaj: 'odmowa', blad: wynik.blad };
    return {
      rodzaj: 'tekst',
      hunks: wynik.wynik.hunks ?? [],
      trafienia: wzorzecPrzyciety === '' ? null : (wynik.wynik.matches ?? []),
    };
  }

  /** Porównanie postaci — `studio.diff.form.compare`; łapie zmianę kroju czy wcięcia, której różnica treści nie widzi. */
  async function porownajPostac(dokument: string): Promise<ZasobPorownania> {
    const wynik = await wywolaj(zal.kanal, Command.StudioDiffFormCompare, {
      documentId: dokument,
      baseVersionId: wersjaBazowa.trim() || undefined,
      targetVersionId: wersjaCel.trim() || undefined,
    });
    if (!wynik.udany || wynik.wynik === undefined) return { rodzaj: 'odmowa', blad: wynik.blad };
    return {
      rodzaj: 'postac',
      wpisy: wynik.wynik.entries,
      dodane: wynik.wynik.added,
      usuniete: wynik.wynik.removed,
      zmienione: wynik.wynik.changed,
    };
  }

  /** Porównanie wyglądu — `studio.diff.visual`; obie wersje są wymagane kontraktem, więc ich brak nazywa się zdaniem. */
  async function porownajWyglad(dokument: string): Promise<ZasobPorownania> {
    const baza = wersjaBazowa.trim();
    const cel = wersjaCel.trim();
    if (baza === '' || cel === '') return { rodzaj: 'brakDanych', zdanie: t.wyglad.wymaganeWersje };

    const wynik = await wywolaj(zal.kanal, Command.StudioDiffVisual, {
      documentId: dokument,
      baseVersionId: baza,
      targetVersionId: cel,
      windowId: zal.idOkna ?? undefined,
    });
    if (!wynik.udany || wynik.wynik === undefined) return { rodzaj: 'odmowa', blad: wynik.blad };
    return {
      rodzaj: 'wyglad',
      obszary: wynik.wynik.regions,
      nakladki: wynik.wynik.overlayAssetIds?.length ?? 0,
    };
  }

  /** Zestawienie z materiałem wejściowym — `studio.diff.source`; materiał wskazuje rdzeń z powiązania dokumentu. */
  async function porownajZeZrodlem(dokument: string): Promise<ZasobPorownania> {
    const wynik = await wywolaj(zal.kanal, Command.StudioDiffSource, {
      documentId: dokument,
      versionId: wersjaCel.trim() || undefined,
    });
    if (!wynik.udany || wynik.wynik === undefined) return { rodzaj: 'odmowa', blad: wynik.blad };
    return { rodzaj: 'zrodlo', hunks: wynik.wynik.hunks, odczytane: wynik.wynik.sourceResolved };
  }

  async function uruchomPorownanie(): Promise<void> {
    if (dokumentId === null) return;
    const mojDokument = dokumentId;
    const mojePokolenie = pokolenie;
    porownanie = { rodzaj: 'ladowanie' };
    bledyPrzeniesienia.clear();
    pominietePrzeniesienia = null;
    odswiez();

    const wynik =
      tryb === 'scalony'
        ? await porownajTresc(mojDokument)
        : tryb === 'postac'
          ? await porownajPostac(mojDokument)
          : tryb === 'wyglad'
            ? await porownajWyglad(mojDokument)
            : await porownajZeZrodlem(mojDokument);
    if (zdjete || mojePokolenie !== pokolenie) return;

    porownanie = wynik;
    odswiez();
    void wczytajAdnotacje();
  }

  /** Adnotacje bieżącego porównania — `studio.annotation.list`; wchodzą pod fragment, którego dotyczą. */
  async function wczytajAdnotacje(): Promise<void> {
    if (dokumentId === null) return;
    const mojDokument = dokumentId;
    const mojePokolenie = pokolenie;
    adnotacjeStan = { rodzaj: 'ladowanie' };
    odswiez();

    const wynik = await wywolaj(zal.kanal, Command.StudioAnnotationList, {
      documentId: mojDokument,
      baseVersionId: wersjaBazowa.trim() || undefined,
      targetVersionId: wersjaCel.trim() || undefined,
    });
    if (zdjete || mojePokolenie !== pokolenie) return;

    adnotacjeStan =
      wynik.udany && wynik.wynik !== undefined
        ? { rodzaj: 'dane', adnotacje: wynik.wynik.annotations }
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

  /** Przestawienie śledzenia zmian — `studio.tracking.set`; stan po zmianie bierze się z odpowiedzi, nie z pola wyboru. */
  async function przestawSledzenie(wlacz: boolean): Promise<void> {
    if (dokumentId === null) return;
    const mojePokolenie = pokolenie;
    sledzenieWTrakcie = true;
    bladSledzenia = undefined;
    odswiez();

    const wynik = await wywolaj(zal.kanal, Command.StudioTrackingSet, {
      documentId: dokumentId,
      enabled: wlacz,
    });
    if (zdjete || mojePokolenie !== pokolenie) return;
    sledzenieWTrakcie = false;

    if (wynik.udany && wynik.wynik !== undefined) {
      sledzenie = wynik.wynik.enabled;
      odswiez();
      await wczytajZmiany();
      return;
    }
    bladSledzenia = wynik.blad;
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
      odswiez();
      return;
    }
    adnotowanyHunk = null;
    trescAdnotacji = '';
    await wczytajAdnotacje();
  }

  /** Przeniesienie fragmentu z wersji odniesienia do stanu bieżącego — `studio.diff.hunk.apply`. */
  async function przeniesFragment(hunkIndex: number): Promise<void> {
    const zrodlo = wersjaBazowa.trim();
    if (dokumentId === null || zrodlo === '') return;
    const mojePokolenie = pokolenie;
    przenoszonyHunk = hunkIndex;
    bledyPrzeniesienia.delete(hunkIndex);
    pominietePrzeniesienia = null;
    odswiez();

    const wynik = await wywolaj(zal.kanal, Command.StudioDiffHunkApply, {
      documentId: dokumentId,
      sourceVersionId: zrodlo,
      hunkIndex,
      author: StudioAuthor.Uzytkownik,
    });
    if (zdjete || mojePokolenie !== pokolenie) return;
    przenoszonyHunk = null;

    if (!wynik.udany || wynik.wynik === undefined) {
      bledyPrzeniesienia.set(hunkIndex, wynik.blad);
      odswiez();
      return;
    }
    if (wynik.wynik.balance.skippedCount > 0) pominietePrzeniesienia = wynik.wynik.balance.skippedCount;
    await uruchomPorownanie();
    await wczytajZmiany();
  }

  /** Raport zmian jako osobny dokument w magazynie — `studio.diff.report.export`. */
  async function wydajRaport(): Promise<void> {
    const baza = wersjaBazowa.trim();
    if (dokumentId === null || baza === '') return;
    const mojePokolenie = pokolenie;
    raportStan = { rodzaj: 'wydaje' };
    odswiez();

    const wynik = await wywolaj(zal.kanal, Command.StudioDiffReportExport, {
      documentId: dokumentId,
      baseVersionId: baza,
      targetVersionId: wersjaCel.trim() || undefined,
      format: formatRaportu,
      includeAnnotations: true,
      windowId: zal.idOkna ?? undefined,
    });
    if (zdjete || mojePokolenie !== pokolenie) return;

    raportStan =
      wynik.udany && wynik.wynik !== undefined
        ? {
            rodzaj: 'gotowy',
            nazwa: wynik.wynik.asset.name ?? wynik.wynik.asset.id,
            bajty: wynik.wynik.sizeBytes,
          }
        : { rodzaj: 'odmowa', blad: wynik.blad };
    odswiez();
  }

  function poleWersji(etykieta: string, wartosc: string, ustaw: (nowa: string) => void): HTMLInputElement {
    const pole = el('input', {
      klasa: 'dn-pole-kontrolka',
      type: 'text',
      'aria-label': etykieta,
      placeholder: t.porownanie.zastepczaWersja,
      value: wartosc,
    }) as HTMLInputElement;
    pole.addEventListener('input', () => ustaw(pole.value));
    pole.addEventListener('keydown', (zdarzenie) => {
      if (zdarzenie.key === 'Enter') void uruchomPorownanie();
    });
    return pole;
  }

  function wierszWersji(): HTMLElement {
    const przyciskPorownaj = el('button', {
      klasa: 'dn-btn dn-btn--zarys dn-btn--sm',
      type: 'button',
      tekst: t.porownanie.porownaj,
    });
    przyciskPorownaj.addEventListener('click', () => void uruchomPorownanie());

    const wyborTrybu = listaWyboru(t.tryb.etykieta, TRYBY, tryb, (kod) => {
      tryb = kod as Tryb;
      void uruchomPorownanie();
    });

    return el('div', { klasa: 'st-panel-wiersz' }, [
      poleWersji(t.porownanie.etykietaBazowa, wersjaBazowa, (nowa) => {
        wersjaBazowa = nowa;
      }),
      '→',
      poleWersji(t.porownanie.etykietaCel, wersjaCel, (nowa) => {
        wersjaCel = nowa;
      }),
      przyciskPorownaj,
      el('span', { klasa: 'dn-meta' }, [wyborTrybu]),
    ]);
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
    const etykietaGrep = el('label', { klasa: 'dn-szukaj' }, [znak('szukaj'), poleGrep]);

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

    const dzieci: Dziecko[] = [el('div', { klasa: 'dn-pole-zestaw' }, [pole, przycisk])];
    if (bladAdnotacji.has(hunkIndex)) {
      dzieci.push(alertOdmowy(t.odmowa.adnotacja, bladAdnotacji.get(hunkIndex)));
    }
    return el('div', { klasa: 'st-panel-wiersz' }, dzieci);
  }

  /** Adnotacje przypisane do fragmentu; autor idzie nazwą roli, bo rdzeń podaje rodzaj autora, nie osobę. */
  function wierszeAdnotacji(hunkIndex: number): HTMLElement[] {
    if (adnotacjeStan === null || adnotacjeStan.rodzaj !== 'dane') return [];
    return adnotacjeStan.adnotacje
      .filter((adnotacja) => adnotacja.hunkIndex === hunkIndex)
      .map((adnotacja) =>
        el('div', { klasa: 'st-panel-wiersz' }, [
          znak('dymek'),
          adnotacja.body,
          el('span', {
            klasa: 'dn-meta',
            tekst: adnotacja.author === StudioAuthor.Model ? t.adnotacja.autorWykonawca : t.adnotacja.autorOperator,
          }),
        ]),
      );
  }

  /** Znaki różnic z prototypu: treść dopisana `dn-diff-dod`, skreślona `dn-diff-usu`. */
  function trescHunka(hunk: StudioDiffHunk): Dziecko[] {
    if (hunk.kind === DiffHunkKind.Added) {
      return [el('span', { klasa: 'dn-diff-dod', tekst: hunk.after ?? '' })];
    }
    if (hunk.kind === DiffHunkKind.Removed) {
      return [el('span', { klasa: 'dn-diff-usu', tekst: hunk.before ?? '' })];
    }
    if (hunk.kind === DiffHunkKind.Changed) {
      return [
        el('span', { klasa: 'dn-diff-usu', tekst: hunk.before ?? '' }),
        ' ',
        el('span', { klasa: 'dn-diff-dod', tekst: hunk.after ?? '' }),
      ];
    }
    return [hunk.after ?? hunk.before ?? ''];
  }

  function wierszHunka(hunk: StudioDiffHunk, zPrzeniesieniem: boolean): HTMLElement[] {
    const dzieci: Dziecko[] = trescHunka(hunk);

    if (zPrzeniesieniem) {
      const przenies = el('button', {
        klasa: 'dn-btn dn-btn--duch dn-btn--sm',
        type: 'button',
        tekst: przenoszonyHunk === hunk.index ? t.przeniesienie.wTrakcie : t.przeniesienie.etykieta,
        disabled: przenoszonyHunk !== null,
      });
      przenies.addEventListener('click', () => void przeniesFragment(hunk.index));
      dzieci.push(przenies);
    }

    const przyciskAdnotacji = el(
      'button',
      { klasa: 'dn-btn-ikona dn-btn-ikona--sm', type: 'button', 'aria-label': t.adnotacja.etykieta },
      [znak('dymek')],
    );
    przyciskAdnotacji.addEventListener('click', () => {
      adnotowanyHunk = adnotowanyHunk === hunk.index ? null : hunk.index;
      trescAdnotacji = '';
      odswiez();
    });
    dzieci.push(przyciskAdnotacji);

    const wiersze = [el('div', { klasa: 'st-panel-wiersz' }, dzieci), ...wierszeAdnotacji(hunk.index)];
    if (bledyPrzeniesienia.has(hunk.index)) {
      wiersze.push(alertOdmowy(t.odmowa.przeniesienie, bledyPrzeniesienia.get(hunk.index)));
    }
    if (adnotowanyHunk === hunk.index) wiersze.push(wierszKompozytoraAdnotacji(hunk.index));
    return wiersze;
  }

  function statystykaHunkow(hunks: StudioDiffHunk[]): HTMLElement {
    let dodane = 0;
    let usuniete = 0;
    let zmienione = 0;
    for (const hunk of hunks) {
      if (hunk.kind === DiffHunkKind.Added) dodane += 1;
      else if (hunk.kind === DiffHunkKind.Removed) usuniete += 1;
      else if (hunk.kind === DiffHunkKind.Changed) zmienione += 1;
    }
    return nota(
      `${t.porownanie.statystykaEtykieta} ` +
        `${t.porownanie.znakDodania}${zLiczba(dodane, t.porownanie.jednostkaFragmenty)} · ` +
        `${t.porownanie.znakUsuniecia}${zLiczba(usuniete, t.porownanie.jednostkaFragmenty)} · ` +
        `${zLiczba(zmienione, t.porownanie.jednostkaZmian)}`,
    );
  }

  function widokTrafien(trafienia: StudioTextMatch[]): HTMLElement[] {
    const naglowek = nota(t.grep.naglowek);
    if (trafienia.length === 0) return [naglowek, pustyStan(t.grep.brakTrafien)];
    return [
      naglowek,
      ...trafienia.map((trafienie) =>
        el('div', { klasa: 'st-panel-wiersz' }, [
          el('span', { klasa: 'dn-meta', tekst: String(trafienie.line) }),
          trafienie.text,
        ]),
      ),
    ];
  }

  /** Różnice postaci znakiem `dn-diff-fmt` — tym samym, którym prototyp oznacza zmianę formatowania. */
  function widokPostaci(zasob: Extract<ZasobPorownania, { rodzaj: 'postac' }>): HTMLElement[] {
    if (zasob.wpisy.length === 0) return [pustyStan(t.postac.brak)];
    const wiersze = zasob.wpisy.map((wpis) =>
      el('div', { klasa: 'st-panel-wiersz' }, [
        el('span', { klasa: 'dn-diff-fmt', tekst: wpis.detail }),
        el('span', { klasa: 'dn-meta', tekst: wpis.area }),
      ]),
    );
    wiersze.push(
      nota(
        `${t.postac.statystykaEtykieta} ` +
          `${t.porownanie.znakDodania}${zLiczba(zasob.dodane, t.postac.jednostkaCech)} · ` +
          `${t.porownanie.znakUsuniecia}${zLiczba(zasob.usuniete, t.postac.jednostkaCech)} · ` +
          `${zLiczba(zasob.zmienione, t.postac.jednostkaZmienionych)}`,
      ),
    );
    return wiersze;
  }

  function widokWygladu(zasob: Extract<ZasobPorownania, { rodzaj: 'wyglad' }>): HTMLElement[] {
    if (zasob.obszary.length === 0) return [pustyStan(t.wyglad.brak)];
    const wiersze = zasob.obszary.map((obszar: StudioVisualDiffRegion) =>
      el('div', { klasa: 'st-panel-wiersz' }, [
        el('span', { klasa: 'dn-diff-fmt', tekst: `${t.wyglad.obszar} ${obszar.page}` }),
        el('span', {
          klasa: 'dn-meta',
          tekst: `${t.wyglad.udzial} ${Math.round(obszar.changeRatio * 100)}%`,
        }),
      ]),
    );
    if (zasob.nakladki > 0) {
      wiersze.push(nota(`${t.wyglad.nakladki} ${zLiczba(zasob.nakladki, t.wyglad.jednostkaNakladek)}`));
    }
    return wiersze;
  }

  function widokPorownania(): HTMLElement[] {
    const glowa: HTMLElement[] = [wierszWersji(), nota(t.porownanie.skadWersje)];
    /* Wzorzec niesie wyłącznie porównanie treści scalonej; pozostałe zakresy nie
       mają go w kontrakcie, więc pole nie stoi tam bezczynnie, tylko ustępuje
       zdaniu mówiącemu, gdzie wzorzec działa. */
    glowa.push(tryb === 'scalony' ? wierszGrep() : nota(t.grep.tylkoTekst));

    if (porownanie === null) return glowa;
    if (porownanie.rodzaj === 'ladowanie') {
      const napis =
        tryb === 'postac'
          ? t.postac.ladowanie
          : tryb === 'wyglad'
            ? t.wyglad.ladowanie
            : tryb === 'zrodlo'
              ? t.zrodlo.ladowanie
              : t.porownanie.ladowanie;
      return [...glowa, wierszLadowania(napis)];
    }
    if (porownanie.rodzaj === 'brakDanych') return [...glowa, pustyStan(porownanie.zdanie)];
    if (porownanie.rodzaj === 'odmowa') {
      const tytul =
        tryb === 'postac'
          ? t.odmowa.postac
          : tryb === 'wyglad'
            ? t.odmowa.wyglad
            : tryb === 'zrodlo'
              ? t.odmowa.zrodlo
              : t.odmowa.porownanie;
      return [...glowa, alertOdmowy(tytul, porownanie.blad)];
    }
    if (porownanie.rodzaj === 'postac') return [...glowa, ...widokPostaci(porownanie)];
    if (porownanie.rodzaj === 'wyglad') return [...glowa, ...widokWygladu(porownanie)];

    const wiersze: HTMLElement[] = [];
    if (adnotacjeStan !== null && adnotacjeStan.rodzaj === 'odmowa') {
      wiersze.push(alertOdmowy(t.odmowa.adnotacje, adnotacjeStan.blad));
    }

    if (porownanie.rodzaj === 'zrodlo') {
      if (!porownanie.odczytane) wiersze.push(pustyStan(t.zrodlo.nieodczytany));
      if (porownanie.hunks.length === 0) wiersze.push(pustyStan(t.zrodlo.brak));
      else {
        for (const hunk of porownanie.hunks) wiersze.push(...wierszHunka(hunk, false));
        wiersze.push(statystykaHunkow(porownanie.hunks));
      }
      return [...glowa, ...wiersze];
    }

    /* Fragment bierze się z wersji odniesienia, więc bez jej wskazania
       przeniesienia nie ma czym zawołać — przycisk nie staje, a zdanie mówi,
       czego brakuje. */
    const zPrzeniesieniem = wersjaBazowa.trim() !== '';
    if (porownanie.hunks.length === 0) wiersze.push(pustyStan(t.porownanie.brakRoznic));
    else {
      for (const hunk of porownanie.hunks) wiersze.push(...wierszHunka(hunk, zPrzeniesieniem));
      if (!zPrzeniesieniem) wiersze.push(nota(t.przeniesienie.wymaganaWersja));
      if (pominietePrzeniesienia !== null) {
        wiersze.push(nota(`${t.przeniesienie.pominiete} ${pominietePrzeniesienia}`));
      }
      wiersze.push(statystykaHunkow(porownanie.hunks));
    }
    if (porownanie.trafienia !== null) wiersze.push(...widokTrafien(porownanie.trafienia));
    return [...glowa, ...wiersze];
  }

  function wierszSledzenia(): HTMLElement[] {
    const pole = el('input', { type: 'checkbox', disabled: sledzenieWTrakcie }) as HTMLInputElement;
    pole.checked = sledzenie === true;
    pole.indeterminate = sledzenie === null;
    pole.addEventListener('change', () => void przestawSledzenie(pole.checked));

    const wiersz = el('div', { klasa: 'st-panel-wiersz' }, [
      el('label', { klasa: 'dn-check-etyk' }, [pole, t.sledzenie.etykieta]),
    ]);
    const zdanie =
      sledzenie === null ? t.sledzenie.nieznane : sledzenie ? t.sledzenie.wlaczone : t.sledzenie.wylaczone;
    const wiersze = [wiersz, nota(zdanie)];
    if (bladSledzenia !== undefined) wiersze.push(alertOdmowy(t.odmowa.sledzenie, bladSledzenia));
    return wiersze;
  }

  function wierszZmiany(zmiana: StudioTrackedChange): HTMLElement[] {
    const klasa =
      zmiana.kind === StudioChangeKind.Wstawienie
        ? 'dn-diff-dod'
        : zmiana.kind === StudioChangeKind.Usuniecie
          ? 'dn-diff-usu'
          : 'dn-diff-fmt';
    const trescZmiany = zmiana.kind === StudioChangeKind.Usuniecie ? zmiana.before : (zmiana.after ?? zmiana.before);

    const dzieci: Dziecko[] = [];
    if (trescZmiany !== undefined && trescZmiany !== '') {
      dzieci.push(el('span', { klasa, tekst: trescZmiany }));
    }
    dzieci.push(
      el('span', {
        klasa: 'dn-meta',
        tekst: zmiana.author === StudioAuthor.Model ? t.zmiany.autorWykonawca : t.zmiany.autorOperator,
      }),
    );

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
    const glowa = [nota(t.zmiany.naglowek), ...wierszSledzenia()];
    if (zmianyStan === null) return glowa;
    if (zmianyStan.rodzaj === 'ladowanie') return [...glowa, wierszLadowania(t.zmiany.ladowanie)];
    if (zmianyStan.rodzaj === 'odmowa') return [...glowa, alertOdmowy(t.odmowa.zmiany, zmianyStan.blad)];
    if (zmianyStan.zmiany.length === 0) return [...glowa, pustyStan(t.zmiany.brak)];
    return [...glowa, ...zmianyStan.zmiany.flatMap(wierszZmiany)];
  }

  function widokRaportu(): HTMLElement[] {
    const bazaPusta = wersjaBazowa.trim() === '';
    const wydaj = el('button', {
      klasa: 'dn-btn dn-btn--atrament dn-btn--sm',
      type: 'button',
      tekst: raportStan !== null && raportStan.rodzaj === 'wydaje' ? t.raport.wTrakcie : t.raport.wydaj,
      disabled: bazaPusta || (raportStan !== null && raportStan.rodzaj === 'wydaje'),
    });
    wydaj.addEventListener('click', () => void wydajRaport());

    const wiersz = el('div', { klasa: 'st-panel-wiersz' }, [
      listaWyboru(t.raport.etykieta, FORMATY_RAPORTU, formatRaportu, (kod) => {
        formatRaportu = kod;
      }),
      wydaj,
    ]);

    const wiersze = [wiersz];
    if (bazaPusta) wiersze.push(nota(t.raport.wymaganaWersja));
    if (raportStan !== null && raportStan.rodzaj === 'odmowa') {
      wiersze.push(alertOdmowy(t.odmowa.raport, raportStan.blad));
    }
    if (raportStan !== null && raportStan.rodzaj === 'gotowy') {
      wiersze.push(
        nota(
          `${t.raport.gotowy} ${raportStan.nazwa} · ${zLiczba(raportStan.bajty, t.raport.jednostkaBajtow)}`,
        ),
      );
    }
    return wiersze;
  }

  function zawartosc(): HTMLElement[] {
    if (zal.idOkna === null) return [alertOdmowy(tekst('dokumentOdmowa.okno'))];
    if (dokumentId === null) return [pustyStan(t.brakDokumentu)];
    return [...widokPorownania(), ...widokZmian(), ...widokRaportu()];
  }

  function odswiez(): void {
    trescPanelu.replaceChildren(...zawartosc());
  }

  wezel.classList.add('sta-okno');
  wezel.replaceChildren(
    el('header', { klasa: 'sta-okno-belka' }, [
      el('span', { klasa: 'sta-okno-tytul' }, [znak('diff'), el('b', { tekst: t.panel.tytul })]),
    ]),
    trescPanelu,
  );
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
