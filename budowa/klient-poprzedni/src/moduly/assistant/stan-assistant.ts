import {
  AssistantActionStatus,
  ChangeKind,
  type AssistantAction,
  type AssistantActivityEntry,
  type Envelope,
} from '../../../../shared/contract';
import type { FazaOkna } from '../../komponenty/faza-okna';
import type { Kanal } from '../../protokol/kanal';
import {
  odczytajDziennik,
  odczytajZlecenia,
  usunZlecenie,
  utworzZapisModulu,
  wchlonObce,
  wchlonZlecenia,
  type ZapisModulu,
} from './zapis-modulu';
import { utworzZrodloAssistant, type ZrodloAssistant } from './zrodlo-assistant';
import { utworzZrodloKontekstow, type ZrodloKontekstow } from './zrodlo-kontekstow';
import { utworzZrodloMowy, type ZrodloMowy } from './zrodlo-mowy';
import { utworzZrodloSchowka, type ZrodloSchowka } from './zrodlo-schowka';
import { utworzZrodloZaplecza, type ZrodloZaplecza } from './zrodlo-zaplecza';

// Jedno źródło prawdy modułu Assistant: subskrypcja rdzenia i powiadamianie widoków.

/** Kod modułu w katalogu rdzenia; jedyne miejsce, w którym widok wiąże się z danymi modułu Assistant udostępnianymi przez kanał komunikacyjny. */
export const KOD_MODULU = 'assistant';

export interface StanAssistant {
  /** Źródło komend obszaru `assistant.*` — okna wołają je wprost. */
  zrodlo: ZrodloAssistant;
  /** Zaplecze: okna sesji i katalog akcji. */
  zaplecze: ZrodloZaplecza;
  /** Trzy źródła obok rdzenia: mowa, konteksty pamięci i zajętość okna oraz historia schowka. */
  mowa: ZrodloMowy;
  konteksty: ZrodloKontekstow;
  schowek: ZrodloSchowka;
  zlecenia(): readonly AssistantAction[];
  /** Zlecenia asystenta tej sesji założone poza oknem modułu, na przykład z nakładki ekranowej. */
  zleceniaObce(): readonly AssistantAction[];
  wpisy(): readonly AssistantActivityEntry[];
  /** Okno modułu przypisane przez rdzeń; pusty napis znaczy brak. */
  idOkna(): string;
  /** Karta sesji, dla której powłoka wczytała moduł; pusty napis znaczy brak. */
  idSesji(): string;
  /** Zlecenie zawężające dziennik; pusty napis znaczy całość zapisu. */
  wybrane(): string;
  /** Przełącza zawężenie dziennika na wskazane zlecenie i czyta je ponownie. */
  wybierz(idZlecenia: string): Promise<void>;
  faza(): FazaOkna;
  powod(): string;
  fazaDziennika(): FazaOkna;
  powodDziennika(): string;
  /** Czy rdzeń odpowiedział już w sprawie zleceń — nazywa pustkę wykazu. */
  pytanoOZlecenia(): boolean;
  /** Czy rdzeń odpowiedział już w sprawie dziennika. */
  pytanoODziennik(): boolean;
  /** Czy padło już pytanie o okno modułu, niezależnie od udzielonej przez rdzeń odpowiedzi. */
  pytanoOOkno(): boolean;
  /** Ustala okno modułu dla sesji; bez niego polecenia nie ma dokąd wysłać. */
  ustalOkno(idSesji: string): Promise<void>;
  odswiezZlecenia(): Promise<void>;
  odswiezDziennik(): Promise<void>;
  /** Wciąga zlecenia po własnym działaniu, bez czekania na zdarzenie. */
  wchlon(przyslane: readonly AssistantAction[]): void;
  obserwuj(sluchacz: () => void): () => void;
  rozlacz(): void;
}

/**
 * Czy zdarzenie dotyczy okna tego modułu.
 *
 * Bez znanego okna nie ma z czym porównywać, więc nie wpuszczamy niczego:
 * wykaz zleceń zbudowany przed odczytem okna byłby wykazem cudzym.
 */
function czyMojeOkno(okno: string, zlecenie: AssistantAction): boolean {
  return okno !== '' && zlecenie.windowId === okno;
}

/**
 * Czy zdarzenie dotyczy sesji, w której stoi ten klient.
 *
 * To jedyna granica, przez którą zlecenie asystenta nie przechodzi. Sesja
 * pochodzi z koperty, bo treść `assistant.action.changed` jej nie niesie —
 * niesie okno, a okno nie mówi, czyja to sesja.
 */
function czyMojaSesja(kanal: Kanal, koperta: Envelope): boolean {
  const moja = kanal.sesja().id();
  return moja !== '' && koperta.sessionId === moja;
}

/** Stany końcowe automatu zlecenia, po których osiągnięciu rdzeń ma już dziennik czynności dopisany i gotowy do odczytu. */
function zeszloZToru(zlecenie: AssistantAction): boolean {
  return (
    zlecenie.status === AssistantActionStatus.Done ||
    zlecenie.status === AssistantActionStatus.Failed ||
    zlecenie.status === AssistantActionStatus.Cancelled
  );
}

export function utworzStanAssistant(kanal: Kanal): StanAssistant {
  const zrodlo = utworzZrodloAssistant(kanal);
  const zaplecze = utworzZrodloZaplecza(kanal);
  const mowa = utworzZrodloMowy(kanal);
  const konteksty = utworzZrodloKontekstow(kanal);
  const schowek = utworzZrodloSchowka(kanal);
  const zapis: ZapisModulu = utworzZapisModulu();
  const sluchacze = new Set<() => void>();

  function oglos(): void {
    for (const sluchacz of [...sluchacze]) sluchacz();
  }

  /** Odczyt z ogłoszeniem fazy przed wywołaniem i po odpowiedzi rdzenia. */
  async function przezOdczyt(
    ustawFaze: () => void,
    odczyt: () => Promise<void>,
  ): Promise<void> {
    ustawFaze();
    oglos();
    await odczyt();
    oglos();
  }

  const odswiezZlecenia = (): Promise<void> =>
    przezOdczyt(
      () => {
        zapis.fazaZlecen = 'ladowanie';
        zapis.powodZlecen = '';
      },
      () => odczytajZlecenia(zapis, zrodlo),
    );

  const odswiezDziennik = (): Promise<void> =>
    przezOdczyt(
      () => {
        zapis.fazaDziennika = 'ladowanie';
        zapis.powodDziennika = '';
      },
      () => odczytajDziennik(zapis, zrodlo),
    );

  const odsubskrybuj = zrodlo.naZmianeZlecenia((tresc, koperta) => {
    // Cudza sesja odpada tu i tylko tu; pusta sesja kanału znaczy, że rdzeń jeszcze jej nie nadał.
    if (!czyMojaSesja(kanal, koperta)) return;
    const moje = czyMojeOkno(zapis.okno, tresc.action);
    if (tresc.change === ChangeKind.Deleted) usunZlecenie(zapis, tresc.action.id);
    else if (moje) wchlonZlecenia(zapis, [tresc.action]);
    else wchlonObce(zapis, [tresc.action]);
    oglos();
    // Dziennik czytamy dopiero po zejściu zlecenia z toru; zawężenie `actionId` pomija okno.
    if (!zeszloZToru(tresc.action)) return;
    if (moje || zapis.wybor === tresc.action.id) void odswiezDziennik();
  });

  return {
    zrodlo,
    zaplecze,
    mowa,
    konteksty,
    schowek,
    zlecenia: () => zapis.zlecenia,
    zleceniaObce: () => zapis.zleceniaObce,
    wpisy: () => zapis.dziennik,
    idOkna: () => zapis.okno,
    idSesji: () => zapis.sesja,
    wybrane: () => zapis.wybor,
    faza: () => zapis.fazaZlecen,
    powod: () => zapis.powodZlecen,
    fazaDziennika: () => zapis.fazaDziennika,
    powodDziennika: () => zapis.powodDziennika,
    pytanoOZlecenia: () => zapis.pytanoOZlecenia,
    pytanoODziennik: () => zapis.pytanoODziennik,
    pytanoOOkno: () => zapis.pytanoOOkno,

    async wybierz(idZlecenia) {
      zapis.wybor = zapis.wybor === idZlecenia ? '' : idZlecenia;
      await odswiezDziennik();
    },

    async ustalOkno(idSesji) {
      zapis.sesja = idSesji;
      const wynik = await zaplecze.okna(idSesji);
      // Bierzemy okno modułu, nie pierwsze lepsze: inne polecenie trafiłoby w cudzą historię.
      zapis.okno = wynik.wynik?.windows.find((wpis) => wpis.moduleId === KOD_MODULU)?.id ?? '';
      // Znacznik pada także po odmowie: pusty przydział po pytaniu znaczy odmowę, a nie brak pytania.
      zapis.pytanoOOkno = true;
      oglos();
    },

    odswiezZlecenia,
    odswiezDziennik,

    wchlon(przyslane) {
      wchlonZlecenia(zapis, przyslane);
      oglos();
    },

    obserwuj(sluchacz) {
      sluchacze.add(sluchacz);
      return () => sluchacze.delete(sluchacz);
    },

    rozlacz() {
      odsubskrybuj();
      sluchacze.clear();
    },
  };
}
