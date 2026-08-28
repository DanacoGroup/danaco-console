/**
 * Moduł Studio — panel dokumentu. Zamyka przekrój pionowy: zakłada sesję i okno
 * modułu w rdzeniu, zakłada dokument, oddaje jego treść do edycji i odsyła ją
 * komendą zapisu, po czym czyta dokument z powrotem.
 *
 * Okno modułu jest bytem rdzenia, nie widoku — bez niego żadna komenda Studia
 * dotykająca dokumentu nie ma gdzie stanąć, bo wszystkie wymagają `windowId`.
 */

import {
  Command,
  ExecutionEnv,
  PermissionMode,
  WindowRole,
  type ErrorInfo,
  type Module,
  type Session,
  type StudioDocument,
} from '../../../../shared/contract.ts';
import type { Kanal } from '../../protokol/kanal.ts';
import { wywolaj } from '../../protokol/wywolanie.ts';
import { el, tekst } from './narzedzia.ts';

export interface NastawyDokumentu {
  /** Kanał, którym panel woła komendy rdzenia. */
  kanal: Kanal;
  /** Karty sesji odtworzone przez rdzeń przy wejściu do środowiska. */
  sesje: Session[];
  /** Moduł, dla którego okno powstaje — jego `id` idzie do `window.create`. */
  modul: Module;
}

export interface PanelDokumentu {
  /** Węzeł panelu, gotowy do wstawienia w bryłę okna. */
  wezel: HTMLElement;
  /** Zdejmuje panel; wynik wywołania w locie ląduje w nicości. */
  zdejmij(): void;
}

type StanZapisu = 'spoczynek' | 'zapisuje' | 'zapisany';

type Stan =
  | { rodzaj: 'zakladanie' }
  | { rodzaj: 'dokument'; dokument: StudioDocument; zapis: StanZapisu }
  | { rodzaj: 'odmowa'; powod: string; blad?: ErrorInfo };

/* Katalogi robocze okna. Kontrakt wymaga pola, nie wymaga zawartości, a źródła
   nie podają katalogu należnego oknu modułu przed wskazaniem Operatora. Wykaz
   zostaje pusty do czasu rozstrzygnięcia: okno bez katalogu nie sięga plików,
   a granica obszaru liczy się właśnie wobec tego wykazu. */
const KATALOGI_ROBOCZE: string[] = [];

export function panelDokumentu(w: NastawyDokumentu): PanelDokumentu {
  let zdjete = false;
  let okno = '';

  const tresc = el('div', { klasa: 'sta-okno-tresc' });
  const wezel = el('section', { klasa: 'sta-okno', 'data-aktywne': 'tak' }, [
    el('header', { klasa: 'sta-okno-belka' }, [
      el('span', { klasa: 'sta-okno-tytul' }, [el('b', { tekst: tekst('dokument.tytul') })]),
    ]),
    tresc,
  ]);

  odswiez({ rodzaj: 'zakladanie' });
  void zaloz();

  /** Sesja pod okno modułu: pierwsza odtworzona przez rdzeń, a przy ich braku nowa. */
  async function sesja(): Promise<string | null> {
    const odtworzona = w.sesje[0];
    if (odtworzona !== undefined) return odtworzona.id;
    const wynik = await wywolaj(w.kanal, Command.SessionCreate, {});
    return wynik.udany ? (wynik.wynik?.session.id ?? null) : null;
  }

  /** Kanał modelu wymagany przez `window.create`: pierwszy czynny z wykazu rdzenia. */
  async function kanalModelu(): Promise<string | null> {
    const wynik = await wywolaj(w.kanal, Command.ChannelList, {});
    if (!wynik.udany) return null;
    const czynny = (wynik.wynik?.channels ?? []).find((k) => k.enabled);
    return czynny?.id ?? null;
  }

  async function zaloz(): Promise<void> {
    const idSesji = await sesja();
    if (zdjete) return;
    if (idSesji === null) return odswiez({ rodzaj: 'odmowa', powod: 'sesja' });

    const idKanalu = await kanalModelu();
    if (zdjete) return;
    if (idKanalu === null) return odswiez({ rodzaj: 'odmowa', powod: 'kanal' });

    const oknoWynik = await wywolaj(w.kanal, Command.WindowCreate, {
      sessionId: idSesji,
      moduleId: w.modul.id,
      modelChannelId: idKanalu,
      workingDirs: KATALOGI_ROBOCZE,
      executionEnv: ExecutionEnv.Local,
      permissionMode: PermissionMode.Manual,
      windowRole: WindowRole.Standalone,
    });
    if (zdjete) return;
    if (!oknoWynik.udany || oknoWynik.wynik === undefined) {
      return odswiez({ rodzaj: 'odmowa', powod: 'okno', blad: oknoWynik.blad });
    }
    okno = oknoWynik.wynik.window.id;

    const dokumentWynik = await wywolaj(w.kanal, Command.StudioDocumentCreate, {
      windowId: okno,
      title: tekst('dokument.nazwaNowego'),
    });
    if (zdjete) return;
    if (!dokumentWynik.udany || dokumentWynik.wynik === undefined) {
      return odswiez({ rodzaj: 'odmowa', powod: 'dokument', blad: dokumentWynik.blad });
    }
    odswiez({ rodzaj: 'dokument', dokument: dokumentWynik.wynik.document, zapis: 'spoczynek' });
  }

  /* Zapis i odczyt idą parą: dopiero treść wrócona z rdzenia dowodzi, że
     przeszła przez repozytorium, a nie została w polu edycji. */
  async function zapisz(dokument: StudioDocument, trescNowa: string): Promise<void> {
    odswiez({ rodzaj: 'dokument', dokument, zapis: 'zapisuje' });
    const zapis = await wywolaj(w.kanal, Command.StudioDocumentSave, {
      documentId: dokument.id,
      content: trescNowa,
      createVersion: true,
    });
    if (zdjete) return;
    if (!zapis.udany) return odswiez({ rodzaj: 'odmowa', powod: 'zapis', blad: zapis.blad });

    const odczyt = await wywolaj(w.kanal, Command.StudioDocumentOpen, {
      windowId: okno,
      documentId: dokument.id,
    });
    if (zdjete) return;
    if (!odczyt.udany || odczyt.wynik === undefined) {
      return odswiez({ rodzaj: 'odmowa', powod: 'zapis', blad: odczyt.blad });
    }
    odswiez({ rodzaj: 'dokument', dokument: odczyt.wynik.document, zapis: 'zapisany' });
  }

  function widokZakladania(): HTMLElement[] {
    return [
      el('div', { klasa: 'dn-wykaz-modulu-poz' }, [
        el('span', {
          klasa: 'dn-kropka dn-kropka--sygnal dn-kropka--tetno',
          'aria-hidden': 'true',
        }),
        tekst('dokument.zakladanie'),
      ]),
    ];
  }

  function widokOdmowy(powod: string, blad?: ErrorInfo): HTMLElement[] {
    return [
      el('div', { klasa: 'dn-alert dn-alert--wstega dn-alert--blad', role: 'alert' }, [
        el('span', { klasa: 'dn-alert-tresc' }, [
          el('b', { tekst: tekst(`dokumentOdmowa.${powod}`) }),
          el('span', { tekst: blad?.message ?? tekst('odmowa.brakOpisu') }),
        ]),
      ]),
    ];
  }

  function widokDokumentu(dokument: StudioDocument, zapis: StanZapisu): HTMLElement[] {
    const pole = el('textarea', {
      klasa: 'dn-pole dn-pole--wielowierszowe',
      'aria-label': tekst('dokument.etykietaTresci'),
      rows: 12,
    }) as HTMLTextAreaElement;
    pole.value = dokument.content ?? '';

    const przycisk = el('button', {
      klasa: 'dn-btn dn-btn--glowny',
      type: 'button',
      tekst: zapis === 'zapisuje' ? tekst('dokument.zapisywanie') : tekst('dokument.zapisz'),
      disabled: zapis === 'zapisuje',
    });
    przycisk.addEventListener('click', () => void zapisz(dokument, pole.value));

    const wersja = dokument.versionId ?? null;
    return [
      el('div', { klasa: 'dn-wykaz-modulu-poz' }, [
        el('b', { tekst: dokument.title ?? tekst('dokument.nazwaNowego') }),
        el('span', {
          klasa: 'dn-meta',
          tekst: wersja === null
            ? tekst('dokument.bezWersji')
            : `${tekst('dokument.wersja')} ${wersja}`,
        }),
      ]),
      pole,
      el('div', { klasa: 'dn-pas-dzialan' }, [
        przycisk,
        zapis === 'zapisany'
          ? el('span', { klasa: 'dn-meta', tekst: tekst('dokument.zapisany') })
          : null,
      ]),
    ];
  }

  function zawartosc(stan: Stan): HTMLElement[] {
    switch (stan.rodzaj) {
      case 'zakladanie':
        return widokZakladania();
      case 'odmowa':
        return widokOdmowy(stan.powod, stan.blad);
      case 'dokument':
        return widokDokumentu(stan.dokument, stan.zapis);
    }
  }

  function odswiez(stan: Stan): void {
    tresc.replaceChildren(...zawartosc(stan));
  }

  return {
    wezel,
    zdejmij() {
      zdjete = true;
      if (wezel.isConnected) wezel.remove();
    },
  };
}
