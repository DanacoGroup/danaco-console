import {
  AssistantActionStatus,
  AssistantActivityKind,
  type AssistantAction,
} from '../../../shared/contract';
import type { Posuniecie, PracaAsystenta, ZrodloPosuniec } from '../aplikacja/zrodlo-posuniec';
import { opisOdmowyBledu } from '../komponenty/odmowa';
import { utworzZrodloAssistant, type ZrodloAssistant } from '../moduly/assistant/zrodlo-assistant';
import { utworzZrodloZaplecza, type ZrodloZaplecza } from '../moduly/assistant/zrodlo-zaplecza';
import type { Kanal } from '../protokol/kanal';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import { ZAPOWIEDZ_BRAKU_GLOSU } from './dostepnosc-mowy';
import { rozpoznajSprawce, zdanieOSprawcy, type ZnanySprawca } from './sprawca-zdarzenia';

/**
 * Stan pływającego dymka Asystenta — jedna rozmowa, jedno okno rdzenia.
 *
 * Droga do rdzenia jest cudza: dymek nie ma własnej warstwy wywołań, tylko
 * bierze `zrodlo-assistant.ts` i `zrodlo-zaplecza.ts` modułu Assistant. Własna
 * warstwa rozjechałaby się z modułem przy pierwszej zmianie kształtu odpowiedzi.
 *
 * Żadna ścieżka nie kończy się milczeniem. Każde miejsce, w którym rdzeń może
 * odmówić — ustalenie okna (`window.list`), wydanie polecenia
 * (`assistant.voice.command`), odczyt dziennika (`assistant.activity.list`)
 * oraz zlecenie zamknięte błędem lub anulowane — kończy się wypowiedzią
 * w historii, nie cichym `return`. Pusta historia czyta się jak „asystent nie
 * ma nic do powiedzenia", co jest nieprawdą.
 *
 * Okno rozróżnia sześć stanów, nie dwa: odmowa rdzenia i brak okna asystenta
 * prowadzą do różnych wniosków, więc nie mogą wyglądać tak samo.
 *
 * Odpowiedź asystenta przychodzi później niż odpowiedź komendy. Rdzeń
 * potwierdza samo przyjęcie zlecenia, a wpis rodzaju `result` dopisuje przy
 * domykaniu (`adapter_modul_asystent_wykonawca.go` → `domknijZlecenie`).
 * Dlatego zejście zlecenia z toru, ogłaszane zdarzeniem
 * `assistant.action.changed`, pociąga odczyt dziennika.
 *
 * Źródło posunięć jest jedno i wspólne z pasem dolnym (`aplikacja/
 * pas-posuniec.ts`). Druga subskrypcja byłaby drugim rozstrzyganiem sprawcy,
 * a `utworzRozstrzyganieSprawcy` zużywa odcisk okna przy weryfikacji: dwa
 * egzemplarze wydałyby dwa różne werdykty o jednym zdarzeniu. Gdy miejsce
 * montażu (`aplikacja/widok-srodowiska.ts`) źródła nie poda, dymek mówi
 * wprost, że posunięć nie dostaje.
 *
 * Sprawca idzie z koperty: `AssistantActionChangedEvent` niesie `actor`
 * i `actorClientId` wypełniane przez rdzeń (`core/sprawca.go`). Rozgłoszenie
 * idzie do całego konta (`transport/rozgloszenie.go`), więc zlecenie założone
 * poza tym dymkiem nie jest cudzą rozmową — wchodzi jako posunięcie (kto, co,
 * w jakim stanie), ale bez treści dziennika, bo treść odpowiedzi należy do
 * okna, w którym padło polecenie.
 */

/** Kod modułu, którego okna szukamy w rejestrze rdzenia. */
const KOD_MODULU = 'assistant';

/**
 * Kto się odezwał — rozstrzyga wygląd wiersza i jego rolę dostępności.
 *
 * `posuniecie` jest osobne celowo: to nie zdanie w rozmowie, lecz ruch
 * asystenta w aplikacji. Zlanie go z wypowiedziami asystenta zatarłoby
 * różnicę między odpowiedzią a czynnością.
 */
export type RodzajWypowiedzi = 'operator' | 'asystent' | 'rdzen' | 'odmowa' | 'posuniecie';

/** Jeden wiersz rozmowy dymka. */
export interface Wypowiedz {
  id: string;
  rodzaj: RodzajWypowiedzi;
  tresc: string;
  chwila: number;
}

/**
 * Stan okna modułu, do którego adresowane jest polecenie.
 *
 * `assistant.voice.command` wymaga pola `windowId`, więc bez okna polecenia
 * nie ma dokąd wysłać. Sześć stanów zamiast dwóch, bo każdy znaczy dla
 * Operatora co innego.
 */
export type StanOkna =
  | { rodzaj: 'niepytane' }
  | { rodzaj: 'bez-sesji' }
  | { rodzaj: 'ustalanie' }
  | { rodzaj: 'gotowe'; id: string }
  | { rodzaj: 'brak-okna' }
  | { rodzaj: 'odmowa'; zdanie: string };

/** Zdanie o stanie okna. Nigdy puste — pustka czyta się jak „nic tu nie ma". */
export function zdanieOStanieOkna(stan: StanOkna): string {
  switch (stan.rodzaj) {
    case 'niepytane':
      return 'Nie pytałem jeszcze rdzenia o okno Asystenta.';
    case 'bez-sesji':
      return 'Rdzeń nie założył jeszcze sesji. Okno Asystenta ustali się zaraz po uzgodnieniu.';
    case 'ustalanie':
      return 'Pytam rdzeń o okno Asystenta (window.list)…';
    case 'gotowe':
      return `Okno Asystenta: ${stan.id}. Polecenia jadą do niego.`;
    case 'brak-okna':
      return (
        'Rdzeń odpowiedział, ale nie ma w tej sesji okna modułu Assistant. ' +
        'assistant.voice.command wymaga pola windowId, więc polecenia nie ma dokąd wysłać. ' +
        'Otwórz moduł Assistant, żeby rdzeń założył jego okno.'
      );
    case 'odmowa':
      return stan.zdanie;
  }
}

/** Czy stan okna pozwala dziś wysłać polecenie. */
export function oknoGotowe(stan: StanOkna): stan is { rodzaj: 'gotowe'; id: string } {
  return stan.rodzaj === 'gotowe';
}

/**
 * Zdanie o posunięciu — składane tak samo, jak składa je pas dolny
 * (`aplikacja/pas-posuniec.ts`, funkcja `dopisz`).
 *
 * Oba widoki biorą posunięcie z jednego `ZrodloPosuniec`, więc muszą je też
 * jednakowo nazwać: ta sama rzecz w dwóch miejscach ma czytać się identycznie.
 *
 * Zdanie nie mówi „Asystent", bo `Posuniecie` niesie wyłącznie `pewnosc` —
 * `zrodlo-posuniec.ts` rozstrzyga sprawcę odciskiem okna zamiast czytać
 * `actor` z koperty.
 */
export function zdaniePosuniecia(posuniecie: Posuniecie): string {
  const zrodloRuchu =
    posuniecie.pewnosc === 'pewne' ? 'inne połączenie tego konta' : 'spoza tego połączenia';
  return `${posuniecie.opis} — ${zrodloRuchu}`;
}

/** Opcje warstwy asystenta — obie wymagają wpięcia w miejscu montażu. */
export interface OpcjeStanuDymka {
  /**
   * Wspólne źródło posunięć — TEN SAM egzemplarz, którym karmi się pas dolny.
   * Pominięte znaczy, że źródła nie podano, i dymek mówi to wprost.
   */
  posuniecia?: ZrodloPosuniec;
  /**
   * Identyfikator bieżącego połączenia (`uzgodnienie.klient.id`) — potrzebny,
   * by odróżnić Operatora przy tym ekranie od Operatora z innego urządzenia.
   * Pusty znaczy „nie znamy własnego" i wtedy tak też jest napisane.
   */
  idKlienta?: string;
}

export interface StanDymka {
  wypowiedzi(): readonly Wypowiedz[];
  okno(): StanOkna;
  /** Czy wywołanie rdzenia jest właśnie w drodze. */
  wToku(): boolean;
  /** Praca asystenta w toku albo `null`, gdy żadne zlecenie nie biegnie. */
  praca(): PracaAsystenta | null;
  /** Ile posunięć przybyło od ostatniego przeczytania (dymek bywa zwinięty). */
  nieprzeczytane(): number;
  /** Zeruje licznik nieprzeczytanych — wołane, gdy dymek jest odsłonięty. */
  oznaczPrzeczytane(): void;
  /** Pyta rdzeń o okno modułu; bezpieczne do powtórzenia. */
  ustalOkno(): Promise<void>;
  /** Wydaje polecenie. Każde zakończenie kończy się wypowiedzią w historii. */
  wyslij(tresc: string): Promise<void>;
  /**
   * Wnosi do historii odmowę, która nie przyszła z rdzenia — dziś jedyną taką
   * jest naciśnięcie mikrofonu. Dymek pokazuje ją także tutaj, nie tylko
   * dymkiem powiadomienia: powiadomienie znika po sekundach, a Operator wraca
   * po przebieg rozmowy do historii i ma tam znaleźć ślad każdej odmowy.
   */
  odnotujOdmowe(tresc: string): void;
  obserwuj(sluchacz: () => void): () => void;
  rozlacz(): void;
}

/** Nadaje wypowiedzi identyfikator — widok nie przerysowuje wierszy bez zmiany. */
let licznik = 0;

export function utworzStanDymka(kanal: Kanal, opcje: OpcjeStanuDymka = {}): StanDymka {
  const zrodlo: ZrodloAssistant = utworzZrodloAssistant(kanal);
  const zaplecze: ZrodloZaplecza = utworzZrodloZaplecza(kanal);
  const zrodloPosuniec = opcje.posuniecia;
  const idKlienta = opcje.idKlienta ?? '';

  const wypowiedzi: Wypowiedz[] = [];
  const sluchacze = new Set<() => void>();
  /** Zlecenia założone z tego dymka wraz z ostatnio pokazanym stanem. */
  const sledzone = new Map<string, AssistantActionStatus>();
  /** Zlecenia założone poza tym dymkiem — ten sam zapis, inny wniosek. */
  const obce = new Map<string, AssistantActionStatus>();
  /** Wpisy dziennika już wniesione do rozmowy — zapora przed powtórzeniem. */
  const wniesione = new Set<string>();

  let stanOkna: StanOkna = { rodzaj: 'niepytane' };
  let wysylka = false;
  let biezacaPraca: PracaAsystenta | null = null;
  let nieprzeczytanych = 0;
  const odsubskrybowania: Odsubskrybuj[] = [];

  function oglos(): void {
    for (const sluchacz of [...sluchacze]) sluchacz();
  }

  function powiedz(rodzaj: RodzajWypowiedzi, tresc: string): void {
    licznik += 1;
    wypowiedzi.push({ id: `w-${licznik}`, rodzaj, tresc, chwila: Date.now() });
    // Licznik rośnie wyłącznie od ruchów asystenta. Własne zdanie Operatora
    // i odpowiedź na nie nie są zaległością — Operator właśnie na nie patrzył.
    if (rodzaj === 'posuniecie') nieprzeczytanych += 1;
    oglos();
  }

  // Zdanie otwierające. Stoi w historii, a nie tylko w nagłówku, bo czytnik
  // ekranu czyta historię — a brak głosu jest tu rzeczą pierwszą do usłyszenia.
  powiedz('rdzen', ZAPOWIEDZ_BRAKU_GLOSU);

  // Drugie zdanie otwierające pada tylko wtedy, gdy miejsce montażu nie podało
  // źródła posunięć. Dymek bez tego źródła pokazuje samą rozmowę i ani jednego
  // ruchu asystenta — cisza wzięta za bezczynność byłaby myląca.
  if (zrodloPosuniec === undefined) {
    powiedz(
      'rdzen',
      'Ten dymek nie dostaje posunięć asystenta — miejsce montażu warstwy ' +
        '(aplikacja/widok-srodowiska.ts) nie podaje wspólnego źródła posunięć. ' +
        'Ruchy asystenta w aplikacji widać na pasie przy dolnej krawędzi; tutaj stoi ' +
        'wyłącznie rozmowa.',
    );
  }

  /** Odczyt dziennika po zejściu zlecenia z toru — wnosi odpowiedź asystenta. */
  async function wniesOdpowiedz(idZlecenia: string): Promise<void> {
    if (!oknoGotowe(stanOkna)) return;
    const wynik = await zrodlo.dziennik(stanOkna.id, idZlecenia);
    if (!wynik.udany || wynik.wynik === undefined) {
      // Zlecenie zeszło z toru, więc odpowiedź istnieje po stronie rdzenia.
      // Cichy powrót zostawiłby w historii samo zdanie Operatora i czytałby
      // się jak brak odpowiedzi asystenta.
      powiedz(
        'odmowa',
        `${opisOdmowyBledu('Odczyt odpowiedzi asystenta', wynik.blad)}. ` +
          `Zlecenie ${idZlecenia} rdzeń domknął — odpowiedź jest w dzienniku, ` +
          'ale ten odczyt jej nie przyniósł. Otwórz Activity Feed w module Assistant.',
      );
      return;
    }

    const odpowiedzi = wynik.wynik.entries
      .filter((wpis) => wpis.actionId === idZlecenia)
      .filter((wpis) => wpis.kind === AssistantActivityKind.Result)
      .filter((wpis) => !wniesione.has(wpis.id))
      .sort((a, b) => a.createdAt - b.createdAt);

    if (odpowiedzi.length === 0) {
      // Dziennik przyszedł, wpisu wyniku w nim nie ma. To też jest odpowiedź —
      // i też nie może wyglądać jak cisza.
      powiedz(
        'rdzen',
        `Rdzeń domknął zlecenie ${idZlecenia}, ale dziennik nie niesie dla niego ` +
          'wpisu rodzaju „wynik". Treści odpowiedzi nie ma czym pokazać.',
      );
      return;
    }

    for (const wpis of odpowiedzi) {
      wniesione.add(wpis.id);
      powiedz('asystent', wpis.content);
    }
  }

  /** Nanosi zmianę stanu zlecenia założonego z tego dymka. */
  function nanies(zlecenie: AssistantAction): void {
    const poprzedni = sledzone.get(zlecenie.id);
    if (poprzedni === zlecenie.status) return;
    sledzone.set(zlecenie.id, zlecenie.status);

    if (zlecenie.status === AssistantActionStatus.Failed) {
      const powod = (zlecenie.result ?? '').trim();
      powiedz(
        'odmowa',
        powod !== ''
          ? `Zlecenie ${zlecenie.id} zakończone błędem: ${powod}`
          : `Zlecenie ${zlecenie.id} zakończone błędem. Rdzeń nie podał powodu.`,
      );
      void wniesOdpowiedz(zlecenie.id);
      return;
    }

    if (zlecenie.status === AssistantActionStatus.Cancelled) {
      powiedz('rdzen', `Zlecenie ${zlecenie.id} anulowane — odpowiedzi nie będzie.`);
      return;
    }

    if (zlecenie.status === AssistantActionStatus.Done) {
      void wniesOdpowiedz(zlecenie.id);
      return;
    }

    if (zlecenie.status === AssistantActionStatus.Running) {
      powiedz('rdzen', `Zlecenie ${zlecenie.id} ruszyło — asystent nad nim pracuje.`);
      return;
    }

    oglos();
  }

  /**
   * Nanosi zlecenie założone poza tym dymkiem — z okna modułu Assistant,
   * z AOD, z telefonu albo ręką samego asystenta.
   *
   * Wnosi sam ruch: kto, jakie zlecenie i w jakim jest stanie. Treści dziennika
   * tu nie ma, bo odpowiedź asystenta na tamto polecenie należy do tamtego
   * okna — tutaj czytałaby się jako część tej rozmowy, którą nie jest.
   */
  function naniesObce(zlecenie: AssistantAction, sprawca: ZnanySprawca): void {
    if (obce.get(zlecenie.id) === zlecenie.status) return;
    obce.set(zlecenie.id, zlecenie.status);

    const tytul = (zlecenie.title ?? '').trim();
    const nazwa = tytul !== '' ? `„${tytul}"` : 'zlecenie bez nazwy';
    const etapy =
      (zlecenie.totalSteps ?? 0) > 0
        ? `, etap ${zlecenie.currentStep ?? 0} z ${zlecenie.totalSteps ?? 0}`
        : '';

    powiedz(
      'posuniecie',
      `${zdanieOSprawcy(sprawca)}: zlecenie ${zlecenie.id} ${nazwa} — stan ${zlecenie.status}${etapy}. ` +
        'Założone poza tym dymkiem; treść odpowiedzi stoi w oknie, w którym padło polecenie.',
    );
  }

  /** Praca asystenta odczytana z tego samego zdarzenia, którego słucha pas. */
  function naniesPrace(zlecenie: AssistantAction): void {
    const tytul = (zlecenie.title ?? '').trim();
    const nazwa = tytul !== '' ? tytul : 'zlecenie bez nazwy';
    if (zlecenie.status !== AssistantActionStatus.Running) {
      if (biezacaPraca?.tytul === nazwa) ustawPrace(null);
      return;
    }
    ustawPrace({
      tytul: nazwa,
      etap: zlecenie.currentStep ?? 0,
      etapow: zlecenie.totalSteps ?? 0,
    });
  }

  function ustawPrace(praca: PracaAsystenta | null): void {
    biezacaPraca = praca;
    oglos();
  }

  odsubskrybowania.push(
    zrodlo.naZmianeZlecenia((tresc) => {
      const zlecenie = tresc.action;
      if (zlecenie === undefined) return;

      // Sprawca wprost z koperty. Rdzeń wypełnia `actor` w tym zdarzeniu
      // (`core/sprawca.go`), więc nie ma tu czego wyprowadzać ze zwłoki
      // i odcisków — a gdy pola brak, mówimy „nie wiadomo", nie „Asystent".
      const sprawca = rozpoznajSprawce(tresc, idKlienta);

      // Wskaźnik pracy bierze się stąd tylko wtedy, gdy wspólnego źródła
      // posunięć nie podano. Podane źródło jest jedyną prawdą o pracy — pas
      // i favikon czytają wtedy ten sam stan.
      if (zrodloPosuniec === undefined) naniesPrace(zlecenie);

      if (sledzone.has(zlecenie.id)) {
        nanies(zlecenie);
        return;
      }
      naniesObce(zlecenie, sprawca);
    }),
  );

  // ── POSUNIĘCIA ZE WSPÓLNEGO ŹRÓDŁA ────────────────────────────────────────
  if (zrodloPosuniec !== undefined) {
    odsubskrybowania.push(
      zrodloPosuniec.naPosuniecie((posuniecie) => powiedz('posuniecie', zdaniePosuniecia(posuniecie))),
    );
    odsubskrybowania.push(zrodloPosuniec.naPrace((praca) => ustawPrace(praca)));
    biezacaPraca = zrodloPosuniec.praca();
  }

  return {
    wypowiedzi: () => wypowiedzi,
    okno: () => stanOkna,
    wToku: () => wysylka,
    praca: () => biezacaPraca,
    nieprzeczytane: () => nieprzeczytanych,

    oznaczPrzeczytane() {
      if (nieprzeczytanych === 0) return;
      nieprzeczytanych = 0;
      oglos();
    },

    async ustalOkno() {
      if (stanOkna.rodzaj === 'ustalanie' || stanOkna.rodzaj === 'gotowe') return;

      const idSesji = kanal.sesja().id();
      if (idSesji === '') {
        stanOkna = { rodzaj: 'bez-sesji' };
        oglos();
        return;
      }

      stanOkna = { rodzaj: 'ustalanie' };
      oglos();

      const wynik = await zaplecze.okna(idSesji);
      if (!wynik.udany || wynik.wynik === undefined) {
        // Odmowa `window.list` NIE jest brakiem okna. Gdyby oba stany zlały się
        // w jeden, Operator dostałby zdanie „otwórz moduł Assistant" w chwili,
        // gdy problemem jest zerwane połączenie z rdzeniem.
        stanOkna = {
          rodzaj: 'odmowa',
          zdanie:
            `${opisOdmowyBledu('Odczyt okien sesji', wynik.blad)}. ` +
            'Bez odpowiedzi rdzenia nie wiem, czy okno Asystenta istnieje — ' +
            'to nie znaczy, że go nie ma.',
        };
        oglos();
        return;
      }

      const okno = wynik.wynik.windows.find((wpis) => wpis.moduleId === KOD_MODULU);
      stanOkna = okno === undefined ? { rodzaj: 'brak-okna' } : { rodzaj: 'gotowe', id: okno.id };
      oglos();
    },

    async wyslij(tresc) {
      const polecenie = tresc.trim();
      if (polecenie === '') {
        powiedz('rdzen', 'Puste polecenie nie jedzie — rdzeń odmówiłby go tak samo.');
        return;
      }
      if (wysylka) {
        powiedz('rdzen', 'Poprzednie polecenie jest jeszcze w drodze. Poczekaj na odpowiedź rdzenia.');
        return;
      }

      if (!oknoGotowe(stanOkna)) {
        // Najpierw jedna próba ustalenia okna — Operator nie ma powodu wiedzieć,
        // że okno ustala się osobnym wywołaniem.
        await this.ustalOkno();
      }
      const cel = stanOkna;
      if (!oknoGotowe(cel)) {
        powiedz('operator', polecenie);
        powiedz(
          'odmowa',
          `Polecenia nie wysłałem. ${zdanieOStanieOkna(cel)}`,
        );
        return;
      }

      powiedz('operator', polecenie);
      wysylka = true;
      oglos();

      const wynik = await zrodlo.polecenie({
        idOkna: cel.id,
        transkrypcja: polecenie,
        profil: '',
        // `speak` zostaje fałszem: rdzeń pola nie odkłada (brak kolumny
        // w `migracja_050_asystent.sql`), a syntezy mowy nie ma. Prośba
        // o odczytanie odpowiedzi nic by nie zdziałała.
        czytaj: false,
      });

      wysylka = false;

      if (!wynik.udany || wynik.wynik === undefined) {
        powiedz('odmowa', opisOdmowyBledu('Wydanie polecenia', wynik.blad));
        return;
      }

      const zlecenie = wynik.wynik.action;
      sledzone.set(zlecenie.id, zlecenie.status);
      powiedz(
        'rdzen',
        `Rdzeń przyjął polecenie i założył zlecenie ${zlecenie.id} (stan: ${zlecenie.status}). ` +
          'Odpowiedź dojdzie tu, gdy zlecenie zejdzie z toru.',
      );

      // Zlecenie domknięte już w chwili odpowiedzi nie doczeka się zdarzenia —
      // zdarzenie poszło, zanim wpisaliśmy je do śledzonych.
      if (
        zlecenie.status === AssistantActionStatus.Done ||
        zlecenie.status === AssistantActionStatus.Failed
      ) {
        await wniesOdpowiedz(zlecenie.id);
      }
    },

    odnotujOdmowe(tresc) {
      powiedz('odmowa', tresc);
    },

    obserwuj(sluchacz) {
      sluchacze.add(sluchacz);
      return () => sluchacze.delete(sluchacz);
    },

    rozlacz() {
      for (const odsubskrybuj of odsubskrybowania.splice(0)) odsubskrybuj();
      sluchacze.clear();
    },
  };
}
