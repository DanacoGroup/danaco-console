/**
 * Pas sterowania okna komunikacji — strefa `.sta-sterowanie` ze źródła kształtu
 * (`design/05-okna/moduly/studio.html`, wnętrze `.sta-kom-dol`) wraz z trzema
 * strefami, które ten sam pas obsługuje danymi: nagłówkiem parametrów, pasem
 * kontekstu i wierszem monitora.
 *
 * Wszystko, co pas pokazuje, przychodzi z komend: `window.state.get` niesie
 * parametry okna, `channel.list` wykaz kanałów modelu, `config.session.get`
 * nakład rozumowania, a `speech.availability.get` stan silnika mowy. Zmiana
 * idzie tą samą drogą: `window.update` przestawia tryb uprawnień i kanał
 * modelu, `config.session.set` nakład rozumowania.
 *
 * Odmiana `sta-chip--diff` została pominięta świadomie: licznik zmian wiersza
 * kontekstu nie ma w kontrakcie żadnego źródła po stronie rozmowy, a chip
 * z wpisaną liczbą byłby treścią przykładową.
 *
 * Katalog treści stoi w tym samym pliku, nie w `tresci.ts` modułu — wzorem
 * `panele/*-tresci.ts`, żeby równolegli wykonawcy nie pisali jednego pliku.
 */

import {
  ChangeKind,
  Command,
  ConfigScope,
  EventType,
  ExecutionEnv,
  PermissionMode,
  ProgressStatus,
  ReasoningEffort,
  SessionConfigArea,
  WindowRole,
  type Channel,
  type Window,
} from '../../../../../shared/contract.ts';
import type { Kanal } from '../../../protokol/kanal.ts';
import { wywolaj } from '../../../protokol/wywolanie.ts';
import { el, zeZnacznika, type Dziecko } from '../narzedzia.ts';
import { opisOdmowy, zaloguj } from '../odmowa.ts';

/** Jedyne miejsce z tekstem widocznym dla Operatora w tej strefie. */
const tresc = {
  parametry: {
    srodowisko: 'Środowisko',
    model: 'Model',
    wysilek: 'Wysiłek',
    nieustalone: 'nieustalony',
    etykietaStanu: 'Stan okna',
  },

  kontekst: {
    etykieta: 'Kontekst',
    pusty: 'Żaden katalog roboczy nie jest przypięty do tego okna.',
    odepnij: 'Odepnij katalog roboczy',
  },

  monitor: {
    /* Wskaźnik trudności tekstu z prototypu nie ma w kontrakcie komendy,
       więc wiersz mówi to wprost zamiast pokazywać wymyśloną miarę. */
    bezMiary: 'Kontrola czytelności tekstu nie jest dostępna w tej wersji.',
    strumien: 'Odpowiedź powstaje…',
  },

  akcje: {
    etykieta: 'Zakres pracy okna',
  },

  pas: {
    etykieta: 'Sterowanie poleceniem',
    doKolejki: 'Do kolejki',
    nagrywaj: 'Nagrywaj',
    dodaj: 'Dodaj do polecenia',
    mikrofon: 'Mikrofon',
    wczytywanie: 'Wczytywanie…',
  },

  tryb: {
    tytul: 'Tryb uprawnień',
    nieustalony: 'Tryb nieustalony',
    odmowa: 'Trybu uprawnień nie udało się zmienić.',
  },

  model: {
    tytul: 'Model',
    nieustalony: 'Kanał nieustalony',
    brakKanalow: 'Rejestr kanałów modelu jest pusty.',
    odmowa: 'Kanału modelu nie udało się zmienić.',
    nieczynny: 'kanał nieczynny',
  },

  wysilek: {
    tytul: 'Wysiłek',
    szybko: 'Szybciej',
    dokladniej: 'Dokładniej',
    nota: 'Nakład rozumowania obowiązuje w tym oknie i przesłania ustawienie sesji.',
    odmowa: 'Nakładu rozumowania nie udało się zapisać.',
  },

  dodaj: {
    pliki: 'Dodaj pliki i zdjęcia',
    zLibrary: 'Dodaj dokument z Library',
    ukosnik: 'Polecenia ukośnikowe',
    /* Trzy czynności zostają widoczne, bo mają być; zdanie mówi, czego dziś
       nie zrobią, zamiast chować przycisk przed Operatorem. */
    niegotowe: 'Dołączanie plików, dokumentów Library i poleceń ukośnikowych wejdzie w kolejnym wydaniu.',
    konektory: 'Konektory',
    konektoryNiegotowe: 'Wykaz konektorów nie jest w tej wersji dostępny.',
    wtyczki: 'Wtyczki',
    wtyczkiNiegotowe: 'Wykaz wtyczek nie jest w tej wersji dostępny.',
  },

  mowa: {
    dyktowanie: 'Dyktowanie',
    odsluch: 'Odsłuch',
    wybudzanie: 'Wybudzanie frazą',
    nasluch: 'Nasłuch ciągły',
    silnik: 'Silnik',
    dostepne: 'dostępne',
    niedostepne: 'niedostępne',
    sprawdzanie: 'Sprawdzanie silnika mowy…',
    odmowa: 'Stanu silnika mowy nie udało się sprawdzić.',
    /* Nagrywanie z okna wymaga przechwycenia dźwięku, którego ta wersja
       programu nie prowadzi; sam przekład nagrania na tekst już stoi. */
    brakNagrywania: 'Nagrywanie głosu z okna wejdzie w kolejnym wydaniu — silnik mowy przekłada dziś tylko gotowe nagranie.',
  },

  stanOkna: {
    [ProgressStatus.Pending]: 'oczekuje',
    [ProgressStatus.Running]: 'pracuje',
    [ProgressStatus.Paused]: 'wstrzymane',
    [ProgressStatus.Stopped]: 'zatrzymane',
    [ProgressStatus.Done]: 'gotowe',
    [ProgressStatus.Failed]: 'zakończone błędem',
  } as Record<string, string>,

  srodowisko: {
    [ExecutionEnv.Local]: 'Ten komputer',
    [ExecutionEnv.Core]: 'Serwer programu',
    [ExecutionEnv.Remote]: 'Komputer zdalny',
  } as Record<string, string>,

  rolaOkna: {
    [WindowRole.Executor]: 'wykonawca',
    [WindowRole.Coordinator]: 'koordynator',
    [WindowRole.Standalone]: 'okno samodzielne',
  } as Record<string, string>,

  /** Nazwy trybów uprawnień po polsku wraz z tym, co każdy z nich znaczy w pracy. */
  trybNazwa: {
    [PermissionMode.Manual]: 'Ręczny',
    [PermissionMode.AcceptEdits]: 'Przyjmij zmiany',
    [PermissionMode.Plan]: 'Plan',
    [PermissionMode.Auto]: 'Auto',
    [PermissionMode.DontAsk]: 'Bez zapytań',
    [PermissionMode.BypassPermissions]: 'Z pominięciem zgód',
  } as Record<string, string>,

  trybOpis: {
    [PermissionMode.Manual]: 'pytanie o zgodę przed każdą zmianą',
    [PermissionMode.AcceptEdits]: 'zgoda na zmiany plików bez pytania',
    [PermissionMode.Plan]: 'praca planistyczna bez zmian w systemie',
    [PermissionMode.Auto]: 'o uprawnieniach rozstrzyga model',
    [PermissionMode.DontAsk]: 'bez zapytań, z zachowaniem ograniczeń',
    [PermissionMode.BypassPermissions]: 'pominięcie kontroli uprawnień',
  } as Record<string, string>,

  wysilekNazwa: {
    [ReasoningEffort.Low]: 'Szybki',
    [ReasoningEffort.Medium]: 'Pośredni',
    [ReasoningEffort.High]: 'Dokładny',
  } as Record<string, string>,
};

/** Rysunki strefy sterowania wzięte ze źródła kształtu; zestaw modułu ich nie niesie. */
const ZNAKI = {
  grot: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75"><path d="m6 9 6 6 6-6"/></svg>',
  ptaszek:
    '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><path d="m5 13 4 4L19 7"/></svg>',
  plus: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><path d="M12 5v14M5 12h14"/></svg>',
  mikrofon:
    '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="3" width="6" height="11" rx="3"/><path d="M5 11a7 7 0 0 0 14 0M12 18v3"/></svg>',
  nagranie: '<svg viewBox="0 0 24 24" fill="currentColor"><rect x="6" y="6" width="12" height="12" rx="2"/></svg>',
  komputer:
    '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="4" width="20" height="14" rx="2"/><path d="M8 20h8"/></svg>',
  teczka:
    '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><path d="M3 7a2 2 0 0 1 2-2h4l2 2h8a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/></svg>',
  wyslij:
    '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><path d="m18 15-6-6-6 6"/></svg>',
} as const;

/** Kolejność trybów uprawnień w dymku; ta sama, w której stoją w kontrakcie. */
const TRYBY: readonly PermissionMode[] = [
  PermissionMode.Manual,
  PermissionMode.AcceptEdits,
  PermissionMode.Plan,
  PermissionMode.Auto,
  PermissionMode.DontAsk,
  PermissionMode.BypassPermissions,
];

/** Nakłady rozumowania od najniższego; położenie suwaka to wprost pozycja w tym wykazie. */
const WYSILKI: readonly ReasoningEffort[] = [ReasoningEffort.Low, ReasoningEffort.Medium, ReasoningEffort.High];

export interface ZaleznosciSterowania {
  /** Kanał, którym pas woła komendy. */
  kanal: Kanal;
  /** Okno komunikacji; odczytywane przy wywołaniu, bo powstaje po montażu bryły. */
  idOkna(): string | null;
  /** Oddaje treść pola polecenia do kolejki — ta sama droga co przycisk formularza. */
  doKolejki(): void;
}

export interface ZamontowaneSterowanie {
  /** Nagłówek parametrów rozmowy — `.sta-kom-naglowek`. */
  parametry: HTMLElement;
  /** Pas przypiętych źródeł kontekstu — `.sta-kom-kontekst`. */
  kontekst: HTMLElement;
  /** Wiersz monitora nad polem polecenia — `.sta-kom-monitor`. */
  monitor: HTMLElement;
  /** Pas zakresu pracy nad polem polecenia — `.sta-kontekst-akcji`. */
  akcje: HTMLElement;
  /** Pas sterowania pod polem polecenia — `.sta-sterowanie`. */
  pas: HTMLElement;
  /** Wczytuje parametry okna; wołane, gdy okno komunikacji już stoi. */
  wczytaj(): void;
  /** Mówi monitorowi, czy w oknie rośnie odpowiedź. */
  strumien(trwa: boolean): void;
  zdejmij(): void;
}

function znak(rysunek: keyof typeof ZNAKI): SVGElement {
  const wezel = zeZnacznika(ZNAKI[rysunek]);
  wezel.setAttribute('aria-hidden', 'true');
  return wezel;
}

/**
 * Wiersz dymka. Wiersz z czynnością jest przyciskiem, choć prototyp rysuje go
 * blokiem: wybór trybu czy modelu musi być dosięgalny klawiaturą. Wiersz bez
 * czynności zostaje blokiem, bo martwy przycisk myli czytnik ekranu.
 */
function wierszDymka(etykieta: Dziecko[], klawisz: Dziecko, przyKliknieciu?: () => void): HTMLElement {
  const dzieci = [el('span', { klasa: 'sta-popover-etykieta' }, etykieta), klawisz];
  if (przyKliknieciu === undefined) return el('div', { klasa: 'sta-popover-wiersz' }, dzieci);
  const wiersz = el('button', { klasa: 'sta-popover-wiersz', type: 'button' }, dzieci);
  wiersz.addEventListener('click', przyKliknieciu);
  return wiersz;
}

export function sterowanieCzatu(zaleznosci: ZaleznosciSterowania): ZamontowaneSterowanie {
  let zdjete = false;
  let okno: Window | null = null;
  let stanProcesu: ProgressStatus | null = null;
  let kanaly: Channel[] = [];
  let wysilek: ReasoningEffort | null = null;
  let komunikat = '';
  let trwaStrumien = false;
  let wczytaneDlaOkna: string | null = null;

  /* Jeden dymek otwarty naraz: dwa nachodzące na siebie zasłaniałyby się
     wzajemnie, bo wszystkie wychodzą w górę z tego samego wiersza. */
  let otwarty: { przycisk: HTMLElement; dymek: HTMLElement } | null = null;

  const parametry = el('div', { klasa: 'sta-kom-naglowek' });
  const kontekst = el('div', { klasa: 'sta-kom-kontekst', 'aria-label': tresc.kontekst.etykieta });
  const monitor = el('div', { klasa: 'sta-kom-monitor', role: 'status' });
  const akcje = el('div', { klasa: 'sta-kontekst-akcji', 'aria-label': tresc.akcje.etykieta });
  const pas = el('div', { klasa: 'sta-sterowanie', 'aria-label': tresc.pas.etykieta });

  const trybNapis = el('span');
  const modelNapis = el('span');
  const wysilekNapis = el('span');
  const trybDymek = el('div', { klasa: 'sta-popover', id: 'pop-tryb-upr', hidden: true });
  const dodajDymek = el('div', { klasa: 'sta-popover', id: 'pop-plus', hidden: true });
  const mikrofonDymek = el('div', { klasa: 'sta-popover', id: 'pop-mik', hidden: true });
  const modelDymek = el('div', { klasa: 'sta-popover', id: 'pop-model', 'data-kotwica': 'prawo', hidden: true });
  const wysilekDymek = el('div', { klasa: 'sta-popover', id: 'pop-wysilek', 'data-kotwica': 'prawo', hidden: true });

  function powiadom(zdanie: string): void {
    komunikat = zdanie;
    odswiezMonitor();
  }

  function odswiezMonitor(): void {
    if (komunikat !== '') {
      monitor.replaceChildren(komunikat);
      return;
    }
    if (trwaStrumien) {
      monitor.replaceChildren(
        el('span', { klasa: 'pt-tetno', 'aria-hidden': 'true' }),
        ` ${tresc.monitor.strumien}`,
      );
      return;
    }
    monitor.replaceChildren(tresc.monitor.bezMiary);
  }

  function zamknijDymek(): void {
    if (otwarty === null) return;
    otwarty.dymek.hidden = true;
    otwarty.przycisk.setAttribute('aria-expanded', 'false');
    otwarty = null;
  }

  function przelaczDymek(przycisk: HTMLElement, dymek: HTMLElement): void {
    const bylOtwarty = otwarty?.dymek === dymek;
    zamknijDymek();
    if (bylOtwarty) return;
    dymek.hidden = false;
    przycisk.setAttribute('aria-expanded', 'true');
    otwarty = { przycisk, dymek };
  }

  /** Chip z dymkiem: przycisk i dymek stoją w jednym `.sta-nrz`, bo dymek kotwiczy się na nim. */
  function grupaDymka(przycisk: HTMLElement, dymek: HTMLElement): HTMLElement {
    przycisk.setAttribute('aria-expanded', 'false');
    przycisk.setAttribute('aria-controls', dymek.id);
    przycisk.addEventListener('click', () => przelaczDymek(przycisk, dymek));
    return el('div', { klasa: 'sta-nrz' }, [przycisk, dymek]);
  }

  const trybChip = el('button', { klasa: 'sta-chip sta-chip--tryb', type: 'button' }, [
    znak('ptaszek'),
    trybNapis,
    znak('grot'),
  ]);

  const dodajPrzycisk = el(
    'button',
    { klasa: 'sta-nrz-przycisk', type: 'button', 'aria-label': tresc.pas.dodaj, title: tresc.pas.dodaj },
    [znak('plus')],
  );

  const mikrofonPrzycisk = el(
    'button',
    { klasa: 'sta-nrz-przycisk', type: 'button', 'aria-label': tresc.pas.mikrofon, title: tresc.pas.mikrofon },
    [znak('mikrofon')],
  );

  const modelChip = el('button', { klasa: 'sta-chip sta-chip--model', type: 'button' }, [modelNapis, znak('grot')]);

  const wysilekChip = el('button', { klasa: 'sta-chip', type: 'button' }, [wysilekNapis, znak('grot')]);

  const suwak = el('input', {
    klasa: 'sta-suwak',
    type: 'range',
    min: 0,
    max: WYSILKI.length - 1,
    step: 1,
    'aria-label': tresc.wysilek.tytul,
  }) as HTMLInputElement;

  const nagrywanie = el(
    'button',
    { klasa: 'sta-rec', type: 'button', 'aria-label': tresc.pas.nagrywaj, title: tresc.pas.nagrywaj },
    [znak('nagranie')],
  );

  const doKolejki = el('button', { klasa: 'sta-chip--glowny sta-prompt-wyslij', type: 'button' }, [
    znak('wyslij'),
    tresc.pas.doKolejki,
  ]);
  doKolejki.addEventListener('click', () => {
    zamknijDymek();
    zaleznosci.doKolejki();
  });

  /* ── zmiany parametrów okna ─────────────────────────────────────────── */

  async function zmienOkno(zmiana: { permissionMode?: PermissionMode; modelChannelId?: string }, odmowa: string): Promise<void> {
    const idOkna = zaleznosci.idOkna();
    if (idOkna === null) return;
    zamknijDymek();
    const wynik = await wywolaj(zaleznosci.kanal, Command.WindowUpdate, { windowId: idOkna, ...zmiana });
    if (zdjete) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      powiadom(`${odmowa} ${opisOdmowy(wynik.blad, 'sterowanie.zmianaOkna')}`);
      return;
    }
    okno = wynik.wynik.window;
    powiadom('');
    odswiez();
  }

  async function zmienWysilek(nowy: ReasoningEffort): Promise<void> {
    const idOkna = zaleznosci.idOkna();
    if (idOkna === null) return;
    const wynik = await wywolaj(zaleznosci.kanal, Command.ConfigSessionSet, {
      areas: [SessionConfigArea.Model],
      config: { model: { reasoningEffort: nowy } },
      scope: ConfigScope.Window,
      scopeId: idOkna,
    });
    if (zdjete) return;
    if (!wynik.udany) {
      powiadom(`${tresc.wysilek.odmowa} ${opisOdmowy(wynik.blad, 'sterowanie.wysilek')}`);
      return;
    }
    wysilek = nowy;
    powiadom('');
    odswiez();
  }

  suwak.addEventListener('change', () => {
    const wybrany = WYSILKI[Number(suwak.value)];
    if (wybrany !== undefined) void zmienWysilek(wybrany);
  });

  nagrywanie.addEventListener('click', () => {
    zamknijDymek();
    powiadom(tresc.mowa.brakNagrywania);
  });

  /* ── rysowanie stref ────────────────────────────────────────────────── */

  /** Pole nagłówka: nazwa parametru i jego wartość, tak jak w źródle kształtu. */
  function poleParametru(nazwa: string, wartosc: string, zDanych: boolean): HTMLElement {
    return el('span', { klasa: 'sta-kom-pole' }, [
      el('span', { tekst: nazwa }),
      el('b', { klasa: zDanych ? 'dane' : null, tekst: wartosc }),
    ]);
  }

  function odswiezParametry(): void {
    const srodowisko = okno === null ? null : (tresc.srodowisko[okno.executionEnv] ?? okno.executionEnv);
    const kanal = kanaly.find((k) => k.id === okno?.modelChannelId);
    const nazwaKanalu = kanal?.name ?? null;
    const nazwaWysilku = wysilek === null ? null : (tresc.wysilekNazwa[wysilek] ?? wysilek);
    parametry.replaceChildren(
      poleParametru(tresc.parametry.srodowisko, srodowisko ?? tresc.parametry.nieustalone, srodowisko !== null),
      poleParametru(tresc.parametry.model, nazwaKanalu ?? tresc.parametry.nieustalone, nazwaKanalu !== null),
      poleParametru(tresc.parametry.wysilek, nazwaWysilku ?? tresc.parametry.nieustalone, nazwaWysilku !== null),
      el('span', { klasa: 'sta-kom-stan', 'aria-label': tresc.parametry.etykietaStanu }, [
        el('span', { klasa: 'dn-plakietka dn-plakietka--sygnal' }, [
          el('span', { klasa: 'pt-tetno', 'aria-hidden': 'true' }),
          ` ${stanProcesu === null ? tresc.parametry.nieustalone : (tresc.stanOkna[stanProcesu] ?? stanProcesu)}`,
        ]),
      ]),
    );
  }

  /** Odpina katalog roboczy, oddając rdzeniowi wykaz pomniejszony o wskazany. */
  async function odepnij(katalog: string): Promise<void> {
    const idOkna = zaleznosci.idOkna();
    if (idOkna === null || okno === null) return;
    const wynik = await wywolaj(zaleznosci.kanal, Command.WindowUpdate, {
      windowId: idOkna,
      workingDirs: okno.workingDirs.filter((k) => k !== katalog),
    });
    if (zdjete) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      zaloguj(wynik.blad, 'sterowanie.odepnij');
      powiadom(opisOdmowy(wynik.blad, 'sterowanie.odepnij'));
      return;
    }
    okno = wynik.wynik.window;
    odswiez();
  }

  function odswiezKontekst(): void {
    const katalogi = okno?.workingDirs ?? [];
    if (katalogi.length === 0) {
      kontekst.replaceChildren(el('span', { klasa: 'dn-meta', tekst: tresc.kontekst.pusty }));
      return;
    }
    kontekst.replaceChildren(
      ...katalogi.map((katalog) => {
        const odpinacz = el('button', {
          klasa: 'odepnij',
          type: 'button',
          'aria-label': `${tresc.kontekst.odepnij}: ${katalog}`,
          tekst: '×',
        });
        odpinacz.addEventListener('click', () => void odepnij(katalog));
        return el('span', { klasa: 'sta-zrodlo' }, [znak('teczka'), katalog, odpinacz]);
      }),
    );
  }

  function odswiezAkcje(): void {
    if (okno === null) {
      akcje.replaceChildren(el('span', { klasa: 'dn-meta', tekst: tresc.pas.wczytywanie }));
      return;
    }
    akcje.replaceChildren(
      el('span', { klasa: 'sta-chip' }, [
        znak('komputer'),
        tresc.srodowisko[okno.executionEnv] ?? okno.executionEnv,
      ]),
      el('span', { klasa: 'sta-chip', tekst: okno.moduleId }),
      el('span', { klasa: 'sta-chip', tekst: tresc.rolaOkna[okno.windowRole] ?? okno.windowRole }),
    );
  }

  function odswiezTryb(): void {
    const biezacy = okno?.permissionMode ?? null;
    trybNapis.textContent = biezacy === null ? tresc.tryb.nieustalony : (tresc.trybNazwa[biezacy] ?? biezacy);
    trybDymek.replaceChildren(
      el('div', { klasa: 'sta-popover-tytul', tekst: tresc.tryb.tytul }),
      ...TRYBY.map((tryb, numer) => {
        const nazwa = tresc.trybNazwa[tryb] ?? tryb;
        const etykieta: Dziecko[] =
          tryb === biezacy
            ? [el('b', { tekst: nazwa }), el('small', { tekst: tresc.trybOpis[tryb] ?? '' })]
            : [nazwa, el('small', { tekst: tresc.trybOpis[tryb] ?? '' })];
        return wierszDymka(etykieta, el('span', { klasa: 'pt-mono', tekst: String(numer + 1) }), () =>
          void zmienOkno({ permissionMode: tryb }, tresc.tryb.odmowa),
        );
      }),
    );
  }

  function odswiezModel(): void {
    const biezacy = kanaly.find((k) => k.id === okno?.modelChannelId) ?? null;
    modelNapis.textContent = biezacy === null ? tresc.model.nieustalony : biezacy.name;
    if (kanaly.length === 0) {
      modelDymek.replaceChildren(
        el('div', { klasa: 'sta-popover-tytul', tekst: tresc.model.tytul }),
        el('p', { klasa: 'dn-nota', tekst: tresc.model.brakKanalow }),
      );
      return;
    }
    modelDymek.replaceChildren(
      el('div', { klasa: 'sta-popover-tytul', tekst: tresc.model.tytul }),
      ...kanaly.map((kanal, numer) => {
        const opis = kanal.enabled ? (kanal.model ?? '') : tresc.model.nieczynny;
        const etykieta: Dziecko[] =
          kanal.id === biezacy?.id
            ? [el('b', { tekst: kanal.name }), opis !== '' && el('small', { tekst: opis })]
            : [kanal.name, opis !== '' && el('small', { tekst: opis })];
        return wierszDymka(etykieta, el('span', { klasa: 'pt-mono', tekst: String(numer + 1) }), () =>
          void zmienOkno({ modelChannelId: kanal.id }, tresc.model.odmowa),
        );
      }),
    );
  }

  function odswiezWysilek(): void {
    const nazwa = wysilek === null ? tresc.parametry.nieustalone : (tresc.wysilekNazwa[wysilek] ?? wysilek);
    wysilekNapis.textContent = nazwa;
    const pozycja = wysilek === null ? 0 : WYSILKI.indexOf(wysilek);
    suwak.value = String(pozycja === -1 ? 0 : pozycja);
    wysilekDymek.replaceChildren(
      el('div', { klasa: 'sta-popover-tytul', tekst: `${tresc.wysilek.tytul} — ${nazwa}` }),
      suwak,
      el('div', { klasa: 'sta-suwak-opis' }, [
        el('span', { tekst: tresc.wysilek.szybko }),
        el('span', { tekst: tresc.wysilek.dokladniej }),
      ]),
      el('p', { klasa: 'dn-nota', tekst: tresc.wysilek.nota }),
    );
  }

  function odswiez(): void {
    odswiezParametry();
    odswiezKontekst();
    odswiezAkcje();
    odswiezTryb();
    odswiezModel();
    odswiezWysilek();
  }

  /* ── wczytanie danych ───────────────────────────────────────────────── */

  async function wczytajStanOkna(idOkna: string): Promise<void> {
    const wynik = await wywolaj(zaleznosci.kanal, Command.WindowStateGet, { windowId: idOkna });
    if (zdjete) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      powiadom(opisOdmowy(wynik.blad, 'sterowanie.stanOkna'));
      return;
    }
    okno = wynik.wynik.window;
    stanProcesu = wynik.wynik.processStatus;
    odswiez();
  }

  async function wczytajKanaly(): Promise<void> {
    const wynik = await wywolaj(zaleznosci.kanal, Command.ChannelList, {});
    if (zdjete) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      zaloguj(wynik.blad, 'sterowanie.kanaly');
      return;
    }
    kanaly = wynik.wynik.channels;
    odswiezParametry();
    odswiezModel();
  }

  async function wczytajWysilek(idOkna: string): Promise<void> {
    const wynik = await wywolaj(zaleznosci.kanal, Command.ConfigSessionGet, {
      scope: ConfigScope.Window,
      scopeId: idOkna,
      areas: [SessionConfigArea.Model],
    });
    if (zdjete) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      zaloguj(wynik.blad, 'sterowanie.wysilek');
      return;
    }
    wysilek = wynik.wynik.config.model?.reasoningEffort ?? null;
    odswiezParametry();
    odswiezWysilek();
  }

  /** Znacznik obecności czynności mowy: ptaszek albo kreska, nigdy puste miejsce. */
  function znacznikMowy(czy: boolean | undefined): HTMLElement {
    return czy === true
      ? el('span', { klasa: 'ptaszek st-ptaszek-akt', tekst: '✓', 'aria-label': tresc.mowa.dostepne })
      : el('span', { klasa: 'pt-mono', tekst: '—', 'aria-label': tresc.mowa.niedostepne });
  }

  async function wczytajMowe(): Promise<void> {
    const wynik = await wywolaj(zaleznosci.kanal, Command.SpeechAvailabilityGet, {});
    if (zdjete) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      mikrofonDymek.replaceChildren(
        el('div', { klasa: 'sta-popover-tytul', tekst: tresc.pas.mikrofon }),
        el('p', { klasa: 'dn-nota', tekst: `${tresc.mowa.odmowa} ${opisOdmowy(wynik.blad, 'sterowanie.mowa')}` }),
      );
      return;
    }
    const stan = wynik.wynik;
    const wiersze: Node[] = [
      el('div', { klasa: 'sta-popover-tytul', tekst: tresc.pas.mikrofon }),
      wierszDymka([tresc.mowa.dyktowanie], znacznikMowy(stan.available)),
      wierszDymka([tresc.mowa.odsluch], znacznikMowy(stan.synthesisAvailable)),
      wierszDymka([tresc.mowa.wybudzanie], znacznikMowy(stan.wakeWordAvailable)),
      wierszDymka([tresc.mowa.nasluch], znacznikMowy(stan.listenAvailable)),
    ];
    if (stan.engine !== undefined) {
      wiersze.push(wierszDymka([tresc.mowa.silnik], el('span', { klasa: 'pt-mono', tekst: stan.engine })));
    }
    /* Powód niedostępności pisze rdzeń dla Operatora — kontrakt każe mu podać
       co, dlaczego i czym Operator to zmieni, więc idzie do dymka wprost. */
    if (stan.reason !== undefined) wiersze.push(el('p', { klasa: 'dn-nota', tekst: stan.reason }));
    mikrofonDymek.replaceChildren(...wiersze);
  }

  function wczytaj(): void {
    const idOkna = zaleznosci.idOkna();
    if (idOkna === null || wczytaneDlaOkna === idOkna) return;
    wczytaneDlaOkna = idOkna;
    void wczytajStanOkna(idOkna);
    void wczytajWysilek(idOkna);
  }

  /* ── złożenie pasa ──────────────────────────────────────────────────── */

  /* Trzy czynności dołączania zostają czynne wbrew brakowi komendy: naciśnięcie
     ma powiedzieć Operatorowi, czego program dziś nie zrobi. */
  function niegotowaCzynnosc(): void {
    zamknijDymek();
    powiadom(tresc.dodaj.niegotowe);
  }

  dodajDymek.replaceChildren(
    el('div', { klasa: 'sta-popover-tytul', tekst: tresc.pas.dodaj }),
    wierszDymka([tresc.dodaj.pliki], el('span', { klasa: 'pt-mono', tekst: 'Ctrl+U' }), niegotowaCzynnosc),
    wierszDymka([tresc.dodaj.zLibrary], null, niegotowaCzynnosc),
    wierszDymka([tresc.dodaj.ukosnik], null, niegotowaCzynnosc),
    el('p', { klasa: 'dn-nota', tekst: tresc.dodaj.niegotowe }),
    el('div', { klasa: 'sta-popover-tytul st-odsun-sekcja', tekst: tresc.dodaj.konektory }),
    el('p', { klasa: 'dn-nota', tekst: tresc.dodaj.konektoryNiegotowe }),
    el('div', { klasa: 'sta-popover-tytul st-odsun-sekcja', tekst: tresc.dodaj.wtyczki }),
    el('p', { klasa: 'dn-nota', tekst: tresc.dodaj.wtyczkiNiegotowe }),
  );

  mikrofonDymek.replaceChildren(
    el('div', { klasa: 'sta-popover-tytul', tekst: tresc.pas.mikrofon }),
    el('p', { klasa: 'dn-nota', tekst: tresc.mowa.sprawdzanie }),
  );

  pas.replaceChildren(
    grupaDymka(trybChip, trybDymek),
    grupaDymka(dodajPrzycisk, dodajDymek),
    grupaDymka(mikrofonPrzycisk, mikrofonDymek),
    el('span', { klasa: 'rozciag' }),
    grupaDymka(modelChip, modelDymek),
    grupaDymka(wysilekChip, wysilekDymek),
    nagrywanie,
    doKolejki,
  );

  /* Dymek zamyka się kliknięciem obok i klawiszem Escape — inaczej zostałby
     otwarty nad polem polecenia i zasłaniał je przy pisaniu. */
  const naKlik = (zdarzenie: MouseEvent): void => {
    if (otwarty === null) return;
    const cel = zdarzenie.target;
    if (cel instanceof Node && pas.contains(cel)) return;
    zamknijDymek();
  };
  const naKlawisz = (zdarzenie: KeyboardEvent): void => {
    if (zdarzenie.key === 'Escape') zamknijDymek();
  };
  document.addEventListener('click', naKlik);
  document.addEventListener('keydown', naKlawisz);

  const odsubskrybujOkno = zaleznosci.kanal.naZdarzenie(EventType.WindowChanged, (zdarzenie) => {
    if (zdjete || zdarzenie.window.id !== zaleznosci.idOkna()) return;
    if (zdarzenie.change === ChangeKind.Deleted) return;
    okno = zdarzenie.window;
    odswiez();
  });

  odswiez();
  odswiezMonitor();
  void wczytajKanaly();
  void wczytajMowe();
  wczytaj();

  return {
    parametry,
    kontekst,
    monitor,
    akcje,
    pas,
    wczytaj,
    strumien(trwa: boolean) {
      if (trwaStrumien === trwa) return;
      trwaStrumien = trwa;
      if (trwa) komunikat = '';
      odswiezMonitor();
    },
    zdejmij() {
      zdjete = true;
      odsubskrybujOkno();
      document.removeEventListener('click', naKlik);
      document.removeEventListener('keydown', naKlawisz);
    },
  };
}
