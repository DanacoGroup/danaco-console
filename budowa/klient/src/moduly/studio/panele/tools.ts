/**
 * Panel Tools Panel okna Studia. Operacje modelu na zaznaczeniu albo na
 * całym dokumencie okna: wykaz operacji — fabrycznych i własnych Operatora —
 * wczytany komendą `studio.operation.list`, przełącznik zasięgu i uruchomienie
 * wybranej operacji komendą `studio.contextual.op`.
 *
 * Dokument, na którym operacje działają, panel poznaje ze zdarzenia rdzenia
 * `studio.document.changed` ograniczonego do własnego okna — kontrakt nie
 * niesie komendy zwracającej wprost dokument otwarty w oknie, a to samo
 * zdarzenie dochodzi do Studio Editor przy każdym założeniu i zapisie
 * dokumentu. Zaznaczenie w edytorze jest stanem samej przeglądarki, nie bytem
 * rdzenia, więc rozmiar zakresu „Zaznaczenie” panel nazywa stanem pustym,
 * dopóki rdzeń nie zacznie go podawać osobną drogą.
 *
 * `studio.operation.save`, `studio.operation.delete` i `studio.search.semantic`
 * nie mają w źródle kształtu (`design/05-okna/moduly/studio.html`,
 * `section#panel-tools`) żadnej powierzchni sterowania — panel ich nie woła,
 * żeby nie dokładać układu, którego prototyp nie niesie.
 */

import {
  ChangeKind,
  Command,
  EventType,
  StudioOperationScope,
  type ErrorInfo,
  type StudioDocument,
  type StudioOperation,
} from '../../../../../shared/contract.ts';
import { wywolaj } from '../../../protokol/wywolanie.ts';
import { ikony } from '../ikony.ts';
import { el, zeZnacznika, type Dziecko } from '../narzedzia.ts';
import { tresciNarzedzi as T } from './tools-tresci.ts';
import type { MontazPanelu } from './umowa.ts';
import { opisOdmowy } from '../odmowa.ts';
import { zLiczba } from '../liczebnik.ts';

type StanBiegu =
  | { rodzaj: 'spoczynek' }
  | { rodzaj: 'uruchamia' }
  | { rodzaj: 'wynik'; tresc?: string }
  | { rodzaj: 'odmowa'; etykieta: string; blad?: ErrorInfo };

type Stan =
  | { rodzaj: 'brakOkna' }
  | { rodzaj: 'ladowanie' }
  | { rodzaj: 'odmowaListy'; blad?: ErrorInfo }
  | { rodzaj: 'gotowy' };

function znak(rysunek: keyof typeof ikony): SVGElement {
  const wezelZnaku = zeZnacznika(ikony[rysunek]);
  wezelZnaku.setAttribute('aria-hidden', 'true');
  return wezelZnaku;
}

function alert(etykieta: string, blad?: ErrorInfo): HTMLElement {
  return el('div', { klasa: 'dn-alert dn-alert--wstega dn-alert--blad', role: 'alert' }, [
    el('span', { klasa: 'dn-alert-tresc' }, [
      el('b', { tekst: etykieta }),
      el('span', { tekst: opisOdmowy(blad, 'tools') }),
    ]),
  ]);
}

function wierszPulsu(etykieta: string): HTMLElement {
  return el('div', { klasa: 'dn-wykaz-modulu-poz' }, [
    el('span', { klasa: 'dn-kropka dn-kropka--sygnal dn-kropka--tetno', 'aria-hidden': 'true' }),
    etykieta,
  ]);
}

export const montujPanelTools: MontazPanelu = (wezel, zaleznosci) => {
  let zdjete = false;
  let stan: Stan = zaleznosci.idOkna === null ? { rodzaj: 'brakOkna' } : { rodzaj: 'ladowanie' };

  let operacje: StudioOperation[] = [];
  let dokument: StudioDocument | null = null;
  let zakres: StudioOperationScope = StudioOperationScope.Selection;
  let wybrana: string | null = null;
  let bieg: StanBiegu = { rodzaj: 'spoczynek' };

  const tresc = el('div', { klasa: 'sta-okno-tresc st-panel-lista' });
  wezel.classList.add('sta-okno');
  wezel.replaceChildren(
    el('header', { klasa: 'sta-okno-belka' }, [
      el('span', { klasa: 'sta-okno-tytul' }, [znak('klucz'), el('b', { tekst: T.panel.tytul })]),
    ]),
    tresc,
  );

  function przelacznikZakresu(): HTMLElement {
    function przycisk(wartosc: StudioOperationScope, etykieta: string): HTMLElement {
      const aktywny = zakres === wartosc;
      const guzik = el('button', {
        klasa: aktywny ? 'dn-btn dn-btn--zarys dn-btn--sm' : 'dn-btn dn-btn--duch dn-btn--sm',
        type: 'button',
        'aria-pressed': aktywny ? 'true' : 'false',
        tekst: etykieta,
      });
      guzik.addEventListener('click', () => {
        if (zakres === wartosc) return;
        zakres = wartosc;
        odswiezGotowy();
      });
      return guzik;
    }
    return el(
      'div',
      { klasa: 'dn-zakladki dn-zakladki--pigulki', role: 'group', 'aria-label': T.zakres.etykietaPrzelacznika },
      [przycisk(StudioOperationScope.Selection, T.zakres.zaznaczenie), przycisk(StudioOperationScope.Document, T.zakres.dokument)],
    );
  }

  function etykietaRozmiaru(): HTMLElement {
    let opis: string;
    if (dokument === null) {
      opis = T.zakres.brakDokumentu;
    } else if (zakres === StudioOperationScope.Selection) {
      opis = T.zakres.brakZaznaczenia;
    } else {
      const trescDokumentu = dokument.content ?? '';
      const znaki = trescDokumentu.length;
      const slowa = trescDokumentu.split(/\s+/).filter((czlon) => czlon.length > 0).length;
      opis = `${T.zakres.etykietaRozmiaru} ${zLiczba(znaki, T.zakres.jednostkaZnaki)} · ${zLiczba(slowa, T.zakres.jednostkaSlowa)}`;
    }
    return el('div', { klasa: 'pt-etykieta', tekst: opis });
  }

  function wierszOperacji(op: StudioOperation): HTMLElement {
    const aktywna = wybrana === op.id;
    const dzieci: Dziecko[] = [
      el('span', { klasa: `dn-kropka ${aktywna ? 'dn-kropka--sukces' : 'dn-kropka--neutralna'}`, 'aria-hidden': 'true' }),
      op.name,
      op.builtin ? null : el('span', { klasa: 'dn-meta', tekst: T.operacje.wlasna }),
    ];
    /* `<div>`, jak w źródle kształtu — `.st-panel-wiersz` nie niesie resetu
       wyglądu natywnego przycisku, więc rolę i klawiaturę dokłada się wprost
       zamiast przez `<button>`, który wniósłby obramowanie przeglądarki. */
    const wiersz = el(
      'div',
      { klasa: 'st-panel-wiersz', role: 'button', tabindex: '0', 'aria-pressed': aktywna ? 'true' : 'false' },
      dzieci,
    );
    const przelacz = (): void => {
      wybrana = aktywna ? null : op.id;
      odswiezGotowy();
    };
    wiersz.addEventListener('click', przelacz);
    wiersz.addEventListener('keydown', (zdarzenie) => {
      if (zdarzenie.key !== 'Enter' && zdarzenie.key !== ' ') return;
      zdarzenie.preventDefault();
      przelacz();
    });
    return wiersz;
  }

  function listaOperacji(): HTMLElement[] {
    if (operacje.length === 0) {
      return [el('p', { klasa: 'dn-pusty-stan dn-pusty-stan--zwarty', tekst: T.operacje.brakWykazu })];
    }
    const wezly: HTMLElement[] = [];
    let ostatniaKategoria: string | null = null;
    for (const op of operacje) {
      if (op.category !== ostatniaKategoria) {
        wezly.push(el('div', { klasa: 'pt-etykieta', tekst: op.category }));
        ostatniaKategoria = op.category;
      }
      wezly.push(wierszOperacji(op));
    }
    return wezly;
  }

  function widokBiegu(): HTMLElement[] {
    switch (bieg.rodzaj) {
      case 'spoczynek':
        return [];
      case 'uruchamia':
        return [wierszPulsu(T.uruchom.wBiegu)];
      case 'wynik':
        return [el('div', { klasa: 'dn-nota' }, [bieg.tresc ?? T.uruchom.brakWyniku])];
      case 'odmowa':
        return [alert(bieg.etykieta, bieg.blad)];
    }
  }

  function przyciskUruchom(): HTMLElement {
    const zablokowany = dokument === null || wybrana === null || bieg.rodzaj === 'uruchamia';
    const przycisk = el('button', {
      klasa: 'dn-btn dn-btn--atrament st-odsun-sekcja',
      type: 'button',
      disabled: zablokowany,
      tekst: bieg.rodzaj === 'uruchamia' ? T.uruchom.wBiegu : T.uruchom.przycisk,
    });
    przycisk.addEventListener('click', () => void uruchomOperacje());
    return przycisk;
  }

  function widokGotowy(): HTMLElement[] {
    return [przelacznikZakresu(), etykietaRozmiaru(), ...listaOperacji(), ...widokBiegu(), przyciskUruchom()];
  }

  function zawartosc(s: Stan): HTMLElement[] {
    switch (s.rodzaj) {
      case 'brakOkna':
        return [alert(T.odmowa.brakOkna)];
      case 'ladowanie':
        return [wierszPulsu(T.ladowanie)];
      case 'odmowaListy':
        return [alert(T.odmowa.lista, s.blad)];
      case 'gotowy':
        return widokGotowy();
    }
  }

  function odswiez(s: Stan): void {
    stan = s;
    tresc.replaceChildren(...zawartosc(s));
  }

  function odswiezGotowy(): void {
    if (stan.rodzaj !== 'gotowy') return;
    tresc.replaceChildren(...zawartosc(stan));
  }

  async function zaladujOperacje(): Promise<void> {
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioOperationList, {});
    if (zdjete) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      odswiez({ rodzaj: 'odmowaListy', blad: wynik.blad });
      return;
    }
    operacje = wynik.wynik.operations;
    odswiez({ rodzaj: 'gotowy' });
  }

  async function uruchomOperacje(): Promise<void> {
    const idOkna = zaleznosci.idOkna;
    const dokumentBiezacy = dokument;
    const operacjaId = wybrana;
    if (idOkna === null || dokumentBiezacy === null || operacjaId === null) return;

    bieg = { rodzaj: 'uruchamia' };
    odswiezGotowy();

    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioContextualOp, {
      windowId: idOkna,
      documentId: dokumentBiezacy.id,
      actionId: operacjaId,
      scope: zakres,
    });
    if (zdjete) return;
    bieg = wynik.udany
      ? { rodzaj: 'wynik', tresc: wynik.wynik?.resultText }
      : { rodzaj: 'odmowa', etykieta: T.odmowa.uruchomienie, blad: wynik.blad };
    odswiezGotowy();
  }

  odswiez(stan);

  const odsubskrybujDokument = zaleznosci.kanal.naZdarzenie(EventType.StudioDocumentChanged, (zdarzenie) => {
    if (zdjete || zaleznosci.idOkna === null || zdarzenie.document.windowId !== zaleznosci.idOkna) return;
    if (zdarzenie.change === ChangeKind.Deleted) {
      if (dokument !== null && dokument.id === zdarzenie.document.id) dokument = null;
    } else {
      dokument = zdarzenie.document;
    }
    odswiezGotowy();
  });

  if (zaleznosci.idOkna !== null) void zaladujOperacje();

  return {
    zdejmij() {
      zdjete = true;
      odsubskrybujDokument();
      wezel.replaceChildren();
    },
  };
};
