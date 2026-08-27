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
 * Stan pływającego dymka Asystenta: jedna rozmowa w jednym oknie rdzenia, obsługiwana
 * przez moduł Assistant; kod modułu wskazuje jego wpis w rejestrze okien.
 */
const KOD_MODULU = 'assistant';

/**
 * Kto się odezwał — rozstrzyga wygląd wiersza i jego rolę dostępności.
 *
 * `posuniecie` jest osobne celowo: to nie zdanie w rozmowie, lecz ruch
 * asystenta w aplikacji. Zlanie go z wypowiedziami asystenta zatarłoby
 * różnicę między odpowiedzią a czynnością.
 */
export type RodzajWypowiedzi = 'operator' | 'asystent' | 'rdzen' | 'odmowa' | 'posuniecie';

/**
 * Jeden wiersz rozmowy dymka: niesie rodzaj wypowiedzi, treść, chwilę wystąpienia oraz
 * identyfikator wykluczający ponowne wyrenderowanie tego samego wiersza.
 */
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

/**
 * Zdanie o stanie okna dla czytnika ekranu: nigdy puste, bo pusty komunikat
 * czytałby się jak brak odpowiedzi rdzenia.
 */
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

/**
 * Czy stan okna pozwala dziś wysłać polecenie: prawda wyłącznie dla stanu
 * `gotowe`, niosącego identyfikator okna modułu Assistant.
 */
export function oknoGotowe(stan: StanOkna): stan is { rodzaj: 'gotowe'; id: string } {
  return stan.rodzaj === 'gotowe';
}

/**
 * Zdanie o posunięciu — składane tak samo, jak w pasie dolnym, bo oba widoki biorą
 * je z jednego `ZrodloPosuniec` i muszą nazwać tę samą rzecz identycznie.
 */
export function zdaniePosuniecia(posuniecie: Posuniecie): string {
  const zrodloRuchu =
    posuniecie.pewnosc === 'pewne' ? 'inne połączenie tego konta' : 'spoza tego połączenia';
  return `${posuniecie.opis} — ${zrodloRuchu}`;
}

/**
 * Opcje warstwy asystenta: obie wymagają wpięcia w miejscu montażu widoku środowiska,
 * inaczej dymek działa bez posunięć i bez tożsamości połączenia.
 */
export interface OpcjeStanuDymka {
  /** Wspólne źródło posunięć — ten sam egzemplarz, którym karmi się pas dolny. */
  posuniecia?: ZrodloPosuniec;
  /** Identyfikator bieżącego połączenia — odróżnia Operatora przy tym ekranie od innych urządzeń. */
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
  /** Wnosi do historii odmowę, która nie przyszła z rdzenia — dziś jedyną taką jest mikrofon. */
  odnotujOdmowe(tresc: string): void;
  obserwuj(sluchacz: () => void): () => void;
  rozlacz(): void;
}

/**
 * Nadaje wypowiedzi identyfikator: widok nie przerysowuje wierszy rozmowy bez zmiany
 * tego identyfikatora.
 */
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
    // Licznik rośnie tylko od ruchów asystenta — zdanie Operatora nie jest zaległością do przeczytania.
    if (rodzaj === 'posuniecie') nieprzeczytanych += 1;
    oglos();
  }

  // Zdanie otwierające stoi w historii, bo czytnik ekranu czyta historię, nie nagłówek.
  powiedz('rdzen', ZAPOWIEDZ_BRAKU_GLOSU);

  // Drugie zdanie pada, gdy miejsce montażu nie podało źródła posunięć — cisza myli się z bezczynnością.
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
      // Zlecenie zeszło z toru — cichy powrót zostawiłby w historii samo zdanie Operatora, bez odpowiedzi.
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
      // Dziennik przyszedł bez wpisu wyniku — to też jest odpowiedź, nie cisza.
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

  /** Nanosi zlecenie założone poza tym dymkiem — wnosi sam ruch, bez treści dziennika tamtego okna. */
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

      // Sprawca wprost z koperty — rdzeń wypełnia actor w tym zdarzeniu, brak pola znaczy nie wiadomo.
      const sprawca = rozpoznajSprawce(tresc, idKlienta);

      // Wskaźnik pracy bierze się stąd tylko, gdy wspólnego źródła posunięć nie podano.
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
        // Odmowa window.list nie jest brakiem okna — zlanie stanów myliłoby zerwane połączenie z brakiem.
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
        // Najpierw jedna próba ustalenia okna — Operator nie musi wiedzieć, że to osobne wywołanie.
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
        // speak zostaje fałszem: rdzeń pola nie odkłada, a syntezy mowy w tej wersji nie ma.
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

      // Zlecenie domknięte już przy odpowiedzi nie doczeka się zdarzenia — ono poszło wcześniej.
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
