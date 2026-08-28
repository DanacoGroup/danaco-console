import type {
  AutomationSchedule,
  AutomationWorkflow,
  BrowserNote,
  BrowserSnapshot,
  BrowserSource,
  ContextTransferResponse,
  Message,
} from '../../../../shared/contract';
import type { TrescNotatki } from './czynnosci-notatek';
import { nazwaRodzaju, type RodzajPrzechwycenia } from './material-sesji';

/**
 * Zdania o skutku czynności modułu Browser — budowane z odpowiedzi rdzenia,
 * nigdy z zamówienia, które okno wysłało. Jedna odpowiedzialność: przekład
 * tego, co wróciło z komendy, na zdanie dla operatora.
 */
export interface Skutek {
  /** Zdanie dla Operatora — wiersz odpowiedzi okna. */
  zdanie: string;
  /** Czy rdzeń zrobił to, co zamówiono; rozbieżność jest odmową, nie sukcesem. */
  udany: boolean;
}

/**
 * Strona pobrana przez rdzeń — adres brany z migawki, która wróciła.
 *
 * `co` nazywa rzecz otwieraną w bierniku („stronę", „źródło") i nie decyduje
 * o niczym więcej: obie drogi idą tą samą komendą `browser.navigate`.
 */
export function skutekPobrania(migawka: BrowserSnapshot, zamowiony: string, co: string): Skutek {
  if (migawka.url !== zamowiony) {
    return {
      zdanie:
        `Rdzeń otworzył ${co} pod adresem ${migawka.url}, a wysłano ${zamowiony}. ` +
        'Obowiązuje adres, który oddał rdzeń — nie ten, który stał w polu.',
      udany: false,
    };
  }
  const tytul = (migawka.title ?? '').trim();
  return {
    zdanie:
      `Rdzeń otworzył ${co} pod adresem ${migawka.url} i oddał migawkę ${migawka.id}` +
      `${tytul === '' ? ' (strona bez tytułu)' : ` strony „${tytul}"`}.`,
    udany: true,
  };
}

/**
 * Migawka zamówiona ze zrzutem ekranu — sprawdzana po polu `screenshotRef`,
 * nie po fladze, którą okno wysłało.
 */
export function skutekZrzutu(migawka: BrowserSnapshot): Skutek {
  const zrzut = (migawka.screenshotRef ?? '').trim();
  if (zrzut === '') {
    return {
      zdanie:
        'Rdzeń oddał migawkę BEZ odnośnika do zrzutu ekranu (pole screenshotRef puste) — ' +
        'zamówiony obraz strony nie przyszedł. Widoczna jest treść strony, nie jej zrzut.',
      udany: false,
    };
  }
  return { zdanie: `Migawka ze zrzutem ekranu pobrana z rdzenia (zrzut: ${zrzut}).`, udany: true };
}

/**
 * Pozycja materiału sesji — nazwana tym, co rdzeń w migawce naprawdę oddał,
 * nie tym, co przechwycenie zamówiło. Zdanie osobno mówi o zamówionym
 * zrzucie, którego w odpowiedzi nie było.
 */
export function skutekPrzechwycenia(
  migawka: BrowserSnapshot,
  rodzaj: RodzajPrzechwycenia,
): Skutek {
  const nazwa = nazwaRodzaju(rodzaj);
  if (rodzaj === 'zrzut') {
    return {
      zdanie: `Przechwycono ${nazwa} strony ${migawka.url} (migawka ${migawka.id}).`,
      udany: true,
    };
  }
  return {
    zdanie:
      `Przechwycono ${nazwa} ${migawka.url} (migawka ${migawka.id}) — rdzeń oddał ją BEZ ` +
      'odnośnika do zrzutu ekranu (pole screenshotRef puste), więc obrazu strony w pozycji nie ma.',
    udany: false,
  };
}

/**
 * Automatyka zapisana w rdzeniu — nazwa i liczba kroków brane z definicji,
 * którą rdzeń oddał, nie z formularza okna. Definicja przyjęta z krokami
 * odrzuconymi jest scenariuszem, który nic nie zrobi.
 */
export function skutekZapisuAutomatyki(
  automatyka: AutomationWorkflow,
  zamowiona: { nazwa: string; krokow: number },
): Skutek {
  const rozjazdy: string[] = [];
  if (automatyka.name !== zamowiona.nazwa) rozjazdy.push('nazwa');
  const krokow = (automatyka.steps ?? []).length;
  if (krokow !== zamowiona.krokow) rozjazdy.push(`liczba kroków (${krokow} zamiast ${zamowiona.krokow})`);
  if (rozjazdy.length > 0) {
    return {
      zdanie:
        `Rdzeń zapisał automatykę ${automatyka.id}, ale oddał ją inaczej niż zamówiono — ` +
        `rozjazd: ${rozjazdy.join(', ')}. Obowiązuje definicja rdzenia.`,
      udany: false,
    };
  }
  return {
    zdanie:
      `Automatyka „${automatyka.name}" zapisana w rdzeniu (${automatyka.id}), ` +
      `kroków: ${krokow}, stan: ${automatyka.enabled ? 'czynna' : 'wstrzymana'}.`,
    udany: true,
  };
}

/** Harmonogram zapisany w rdzeniu — cykliczność brana zawsze z wiersza zwrotnego, nie z pola formularza okna. */
export function skutekHarmonogramu(harmonogram: AutomationSchedule, zamowionyCron: string): Skutek {
  const cron = (harmonogram.cron ?? '').trim();
  if (cron !== zamowionyCron) {
    return {
      zdanie:
        `Rdzeń zapisał harmonogram ${harmonogram.id} z cyklicznością „${cron}", ` +
        `a wysłano „${zamowionyCron}". Obowiązuje wartość rdzenia.`,
      udany: false,
    };
  }
  return {
    zdanie:
      `Harmonogram ${harmonogram.id} zapisany: cykliczność „${cron}", ` +
      `${harmonogram.enabled ? 'obowiązuje' : 'nie obowiązuje'}.`,
    udany: true,
  };
}

/** Źródło zapisane w rdzeniu — adres brany z wiersza, który wrócił w odpowiedzi, nie z pola formularza. */
export function skutekZapisuZrodla(zrodlo: BrowserSource, zamowiony: string): Skutek {
  if (zrodlo.url !== zamowiony) {
    return {
      zdanie:
        `Rdzeń zapisał źródło pod adresem ${zrodlo.url}, a wysłano ${zamowiony}. ` +
        'W wykazie stoi wartość rdzenia — zapisu wysłanego adresu nie było.',
      udany: false,
    };
  }
  return { zdanie: `Źródło ${zrodlo.url} zapisane w rdzeniu (${zrodlo.id}).`, udany: true };
}

/** Notatka zapisana w rdzeniu — treść, powiązanie i cytat brane zawsze z wiersza zwrotnego, nie z formularza. */
export function skutekZapisuNotatki(notatka: BrowserNote, zamowione: TrescNotatki): Skutek {
  const rozjazdy: string[] = [];
  if (notatka.content !== zamowione.tresc) rozjazdy.push('treść notatki');
  if (zamowione.idZrodla !== '' && (notatka.sourceId ?? '') !== zamowione.idZrodla) {
    rozjazdy.push('powiązanie ze źródłem');
  }
  if (zamowione.cytat !== '' && (notatka.quote ?? '') !== zamowione.cytat) {
    rozjazdy.push('cytowany fragment');
  }
  if (rozjazdy.length > 0) {
    return {
      zdanie:
        `Rdzeń zapisał notatkę ${notatka.id}, ale oddał ją inaczej niż zamówiono — ` +
        `rozjazd: ${rozjazdy.join(', ')}. Obowiązuje wiersz rdzenia.`,
      udany: false,
    };
  }
  return { zdanie: `Notatka ${notatka.id} zapisana w rdzeniu${dopisekNotatki(notatka)}.`, udany: true };
}

/** Co notatka niesie po zapisie w rdzeniu — brane z wiersza zwrotnego odpowiedzi, nie z formularza okna. */
function dopisekNotatki(notatka: BrowserNote): string {
  const czesci: string[] = [];
  if ((notatka.sourceId ?? '') !== '') czesci.push(`ze źródłem ${notatka.sourceId ?? ''}`);
  if ((notatka.quote ?? '') !== '') czesci.push('z cytatem fragmentu strony');
  return czesci.length === 0 ? '' : ` ${czesci.join(' i ')}`;
}

/**
 * Pytanie zadane w oknie rozmowy — sprawdzane po wiadomości, którą rdzeń
 * założył, a nie po tym, że w ogóle odpowiedział. Wiadomość zapisana w cudzym
 * oknie znaczy, że pytanie poszło nie tam, gdzie operator je zadał.
 */
export function skutekPytania(wiadomosc: Message, zamowioneOkno: string): Skutek {
  if (wiadomosc.windowId !== zamowioneOkno) {
    return {
      zdanie:
        `Rdzeń zapisał pytanie w oknie ${wiadomosc.windowId}, a wysłano je do ${zamowioneOkno}. ` +
        'Odpowiedź modelu pojawi się w tamtym oknie, nie w tym.',
      udany: false,
    };
  }
  return {
    zdanie: `Pytanie o zaznaczony fragment przyjęte przez okno rozmowy (wiadomość ${wiadomosc.id}).`,
    udany: true,
  };
}

/**
 * Adnotacja spłaszczona do PNG i dołączona do rozmowy — oceniana po tym, co
 * rdzeń oddał w polu załączników, nie po tym, że okno obraz wysłało.
 */
export function skutekAdnotacji(wiadomosc: Message, zamowioneOkno: string, obraz: string): Skutek {
  if (wiadomosc.windowId !== zamowioneOkno) {
    return {
      zdanie:
        `Rdzeń zapisał adnotację w oknie ${wiadomosc.windowId}, a wysłano ją do ${zamowioneOkno}. ` +
        'Rysunek stoi w tamtej rozmowie, nie w tej.',
      udany: false,
    };
  }
  if (!(wiadomosc.attachments ?? []).includes(obraz)) {
    return {
      zdanie:
        `Rdzeń przyjął wiadomość ${wiadomosc.id}, ale oddał ją BEZ obrazu adnotacji ` +
        '(pole attachments nie niesie wysłanego PNG) — w rozmowie stoi sam opis, nie rysunek.',
      udany: false,
    };
  }
  return {
    zdanie: `Adnotacja dołączona do rozmowy jako obraz PNG (wiadomość ${wiadomosc.id}).`,
    udany: true,
  };
}

/**
 * Przeniesienie kompletu kontekstu — oceniane po znaczniku przeniesienia
 * i module okna, które wróciło w odpowiedzi, a nie po tym, co okno wysłało.
 */
export function skutekPrzekazania(
  odpowiedz: ContextTransferResponse,
  zamowionyModul: string,
  czynnosc: string,
): Skutek {
  if (!odpowiedz.transferred) {
    return {
      zdanie:
        `${czynnosc}: rdzeń przyjął żądanie, ale oddał je jako NIEPRZENIESIONE ` +
        `(transferred=false, okno ${odpowiedz.window.id}).`,
      udany: false,
    };
  }
  if (odpowiedz.window.moduleId !== zamowionyModul) {
    return {
      zdanie:
        `${czynnosc}: rdzeń oddał okno modułu ${odpowiedz.window.moduleId}, ` +
        `a przekazanie szło do ${zamowionyModul} (okno ${odpowiedz.window.id}).`,
      udany: false,
    };
  }
  return {
    zdanie: `${czynnosc} — rdzeń przeniósł komplet do modułu ${odpowiedz.window.moduleId} (okno ${odpowiedz.window.id}).`,
    udany: true,
  };
}
