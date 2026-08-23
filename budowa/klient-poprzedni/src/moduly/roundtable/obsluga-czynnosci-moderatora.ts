import {
  RoundtableTurnStatus,
  type RoundtableModeratorDirectRequest,
  type RoundtableModeratorDirectResponse,
  type RoundtableParticipant,
  type RoundtableTurn,
} from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import type { Wynik } from '../../protokol/kanal';
import { nazwaUczestnika, type StanDebaty } from './stan-debaty';
import type { StanTresci } from './stany-okna';
import type { ZrodloRoundtable } from './zrodlo-roundtable';

/**
 * Obsługa odpowiedzi czynności moderatora — wydzielona z
 * `okno-moderator-panel.ts` wzorem `KontekstBudowania` z Build Output. Zwarta
 * odpowiedzialność: wysyłka `roundtable.moderator.direct`, zapis tury, warunkowy
 * zapis składu i złożenie zdania potwierdzenia.
 *
 * Potwierdzenie mówi, co zrobił rdzeń, a nie co wysłało okno: zdanie składa
 * funkcja biorąca odpowiedź, nie napis wnoszony razem z żądaniem. Gotowe zdanie
 * sukcesu nie odróżniłoby zamknięcia tury od potwierdzenia tury już zamkniętej,
 * a przy zmianie zagadnienia przemilczałoby, że rdzeń zamknął turę bieżącą
 * i otworzył nową.
 *
 * Tura zapisywana jest zawsze, skład wyłącznie wtedy, gdy rdzeń go przysłał:
 * `participants` jest w kontrakcie nieobowiązkowe, więc bez tego pola okno nie
 * ma dowodu, że zmiana składu zaszła.
 *
 * Odmowa nie jest jedynym niepowodzeniem. `sprawdzKsztalt` oddaje `udany: false`
 * także wtedy, gdy rdzeń odpowiedział, ale bez pola obowiązkowego — zdanie
 * „rdzeń odmówił” byłoby wtedy zarzutem, którego rdzeń nie orzekł. Dlatego
 * mówimy o czynności, która nie doszła do skutku, a powód niesie treść błędu
 * wprost z kontraktu.
 */

/** Zdanie potwierdzenia wraz z oceną: czy czynność naprawdę się odbyła. */
export interface Potwierdzenie {
  zdanie: string;
  udane: boolean;
}

/** Opis jednej czynności moderatora dla obsługi jej odpowiedzi. */
export interface OpcjeCzynnosci {
  /** Nazwa czynności w zdaniu o niepowodzeniu — wskazuje, co nie doszło do skutku. */
  czynnosc: string;
  /** Składa zdanie potwierdzenia z odpowiedzi rdzenia. */
  potwierdz(odpowiedz: RoundtableModeratorDirectResponse): Potwierdzenie;
  /** Zdanie doklejane do niepowodzenia, gdy okno wie o nim coś ponad powód rdzenia. */
  przyNiepowodzeniu?: string;
}

/** Kontekst wysyłki czynności moderatora i obsługi jej odpowiedzi. */
export interface KontekstCzynnosci {
  obsluz(zadanie: RoundtableModeratorDirectRequest, opcje: OpcjeCzynnosci): void;
}

export function utworzKontekstCzynnosci(
  zrodlo: ZrodloRoundtable,
  stan: StanDebaty,
  tresc: StanTresci,
  rysuj: () => void,
): KontekstCzynnosci {
  return {
    obsluz(zadanie, opcje) {
      tresc.ladowanie(`Wysyłanie czynności moderatora: ${opcje.czynnosc}…`);
      void zrodlo.moderuj(zadanie).then((wynik) => {
        obsluzOdpowiedzCzynnosci(wynik, opcje, stan, tresc, rysuj);
      });
    },
  };
}

function obsluzOdpowiedzCzynnosci(
  wynik: Wynik<RoundtableModeratorDirectResponse>,
  opcje: OpcjeCzynnosci,
  stan: StanDebaty,
  tresc: StanTresci,
  rysuj: () => void,
): void {
  if (!wynik.udany || wynik.wynik === undefined) {
    // Niepowodzenie czynności nie kasuje treści okna. `tresc.blad(...)` czyści
    // miejsce treści, więc nieudana czynność zabrałaby z ekranu turę i cały
    // skład — okno pokazywałoby pustą debatę tam, gdzie debata jest. Odmowa
    // idzie zatem pasem czynności (`data-udane='false'`, arkusz maluje go barwą
    // błędu), a wykaz zostaje przerysowany.
    rysuj();
    const uwaga = opcje.przyNiepowodzeniu === undefined ? '' : ` ${opcje.przyNiepowodzeniu}`;
    const zdanie = opisOdmowyBledu(
      `Czynność moderatora „${opcje.czynnosc}” nie doszła do skutku`,
      wynik.blad,
    );
    tresc.potwierdzenie(`${zdanie}.${uwaga}`, false);
    return;
  }
  stan.ustawTure(wynik.wynik.turn);
  const sklad = wynik.wynik.participants;
  if (sklad !== undefined) stan.ustawUczestnikow(sklad);
  rysuj();
  const potwierdzenie = opcje.potwierdz(wynik.wynik);
  tresc.potwierdzenie(potwierdzenie.zdanie, potwierdzenie.udane);
}

/**
 * Zamknięcie tury.
 *
 * Zdanie mówi, co rdzeń oddał, a nie co zrobił. Rdzeń nie odmawia zamknięcia
 * tury już zamkniętej — oddaje ją bez zmiany, więc okno nie ma czym odróżnić
 * zamknięcia od potwierdzenia stanu zastanego. Stan sprzed wywołania też nie
 * jest świadkiem pewnym, bo po zamknięciu potrafi jeszcze przyjść rozgłoszenie
 * z turą oznaczoną jako otwarta. Dlatego zdanie opisuje turę z odpowiedzi wraz
 * z chwilą zamknięcia, a czasownik „zamknął” pada wyłącznie tam, gdzie stan
 * sprzed wywołania mówił coś innego.
 */
export function potwierdzZamkniecie(
  odpowiedz: RoundtableModeratorDirectResponse,
  bylaZamknieta: boolean,
): Potwierdzenie {
  const tura = odpowiedz.turn;
  if (tura.status !== RoundtableTurnStatus.Closed) {
    return {
      zdanie:
        `Rdzeń przyjął wywołanie, ale oddał turę #${tura.index} w stanie „${tura.status}” — ` +
        'tura NIE jest zamknięta.',
      udane: false,
    };
  }
  const chwila =
    tura.closedAt === undefined || !Number.isFinite(tura.closedAt) || tura.closedAt <= 0
      ? 'rdzeń nie podał chwili zamknięcia'
      : `zamknięta ${new Date(tura.closedAt).toLocaleString('pl-PL')}`;
  if (bylaZamknieta) {
    return {
      zdanie: `Tura #${tura.index} była już zamknięta — rdzeń oddał ją bez zmiany (${chwila}).`,
      udane: false,
    };
  }
  return { zdanie: `Rdzeń oddał turę #${tura.index} jako zamkniętą (${chwila}).`, udane: true };
}

/**
 * Zmiana zagadnienia. Odpowiedź niesie turę po czynności — jeśli ma inny
 * identyfikator niż tura sprzed wywołania, rdzeń otworzył turę nową, a poprzednia
 * została zamknięta. Zdanie mówi to wprost, zamiast opisywać samo żądanie.
 */
export function potwierdzZagadnienie(
  odpowiedz: RoundtableModeratorDirectResponse,
  przed: RoundtableTurn | null,
  zadane: string,
): Potwierdzenie {
  const tura = odpowiedz.turn;
  const temat = tura.topic ?? '';
  if (temat !== zadane) {
    const oddane = temat === '' ? 'bez zagadnienia' : `„${temat}”`;
    return {
      zdanie:
        `Rdzeń przyjął wywołanie, ale oddał turę #${tura.index} ${oddane} zamiast żądanego ` +
        `„${zadane}” — zagadnienia nie zapisano tak, jak wysłano.`,
      udane: false,
    };
  }
  if (przed === null || przed.id !== tura.id) {
    return {
      zdanie:
        `Rdzeń otworzył turę #${tura.index} (${tura.status}) o zagadnieniu „${temat}”` +
        (przed === null ? '.' : `, a turę #${przed.index} zamknął.`),
      udane: true,
    };
  }
  return {
    zdanie: `Rdzeń zapisał zagadnienie „${temat}” w turze #${tura.index} (${tura.status}).`,
    udane: true,
  };
}

/**
 * Wyciszenie i jego zdjęcie. Rdzeń oddaje skład, więc okno czyta z odpowiedzi,
 * czy uczestnik jest wyciszony — nie powtarza tego, co samo wysłało.
 */
export function potwierdzWyciszenie(
  odpowiedz: RoundtableModeratorDirectResponse,
  idUczestnika: string,
  wyciszamy: boolean,
  stan: StanDebaty,
): Potwierdzenie {
  const czynnosc = wyciszamy ? 'wyciszenia' : 'zdjęcia wyciszenia';
  const sklad = odpowiedz.participants;
  if (sklad === undefined) {
    return {
      zdanie: `Rdzeń przyjął wywołanie, ale nie przysłał składu — okno nie potwierdza ${czynnosc}.`,
      udane: false,
    };
  }
  const po = sklad.find((uczestnik) => uczestnik.id === idUczestnika);
  if (po === undefined) {
    return {
      zdanie:
        `Rdzeń przyjął wywołanie, ale w oddanym składzie nie ma uczestnika ${idUczestnika} — ` +
        `okno nie potwierdza ${czynnosc}.`,
      udane: false,
    };
  }
  const nazwa = nazwaUczestnika(po, stan.opisKanalu(po.channelId));
  const wyciszony = po.muted === true;
  if (wyciszony !== wyciszamy) {
    return {
      zdanie:
        `Rdzeń przyjął wywołanie, ale w oddanym składzie „${nazwa}” jest ` +
        `${wyciszony ? 'wyciszony' : 'niewyciszony'} — ${czynnosc} nie odbyło się.`,
      udane: false,
    };
  }
  return {
    zdanie: `Rdzeń oddał skład, w którym „${nazwa}” jest ${wyciszony ? 'wyciszony' : 'niewyciszony'}.`,
    udane: true,
  };
}

/**
 * Kolejność głosu. Zdanie powstaje z pola `order` uczestników w odpowiedzi, nie
 * z listy, którą okno wysłało: gdyby rdzeń zapisał inny porządek, potwierdzenie
 * złożone z żądania byłoby niezgodne ze stanem rdzenia bez żadnego śladu.
 */
export function potwierdzKolejnosc(
  odpowiedz: RoundtableModeratorDirectResponse,
  zadana: readonly string[],
  stan: StanDebaty,
): Potwierdzenie {
  const sklad = odpowiedz.participants;
  if (sklad === undefined) {
    return {
      zdanie:
        'Rdzeń przyjął wywołanie, ale nie przysłał składu — okno nie potwierdza kolejności głosu.',
      udane: false,
    };
  }
  const oddana = [...sklad]
    .filter((uczestnik) => uczestnik.order !== undefined)
    .sort((pierwszy, drugi) => (pierwszy.order ?? 0) - (drugi.order ?? 0));
  if (oddana.length !== sklad.length) {
    return {
      zdanie:
        'Rdzeń oddał skład, w którym nie każdy uczestnik ma miejsce w kolejności głosu — ' +
        'okno nie potwierdza nowego porządku.',
      udane: false,
    };
  }
  const zgodna =
    oddana.length === zadana.length &&
    oddana.every((uczestnik, miejsce) => uczestnik.id === zadana[miejsce]);
  const wykaz = oddana.map((uczestnik, miejsce) => `#${miejsce + 1} ${podpis(uczestnik, stan)}`);
  if (!zgodna) {
    return {
      zdanie: `Rdzeń oddał kolejność inną niż wysłana: ${wykaz.join(', ')}.`,
      udane: false,
    };
  }
  return { zdanie: `Rdzeń zapisał kolejność głosu: ${wykaz.join(', ')}.`, udane: true };
}

function podpis(uczestnik: RoundtableParticipant, stan: StanDebaty): string {
  return nazwaUczestnika(uczestnik, stan.opisKanalu(uczestnik.channelId));
}
