import type { AssistantAction, AssistantActivityEntry } from '../../../../shared/contract';
import type { FazaOkna } from '../../komponenty/faza-okna';
import { opisOdmowy } from '../../komponenty/odmowa';
import type { ZrodloAssistant } from './zrodlo-assistant';

/**
 * Zapis modułu Assistant niesie dane trzech okien wraz z przejściami między fazami odczytu,
 * zmieniany w miejscu, a nie odtwarzany.
 */
export interface ZapisModulu {
  zlecenia: AssistantAction[];
  /** Zlecenia tej samej sesji założone poza oknem modułu; bierzemy je wyłącznie ze zdarzenia rdzenia. */
  zleceniaObce: AssistantAction[];
  dziennik: AssistantActivityEntry[];
  /** Okno modułu przypisane przez rdzeń; pusty napis znaczy brak. */
  okno: string;
  /** Karta sesji, w której moduł został wczytany; okna zarządcze pytają rdzeń o byty przypisane do niej. */
  sesja: string;
  /** Zlecenie zawężające dziennik; pusty napis znaczy całość zapisu. */
  wybor: string;
  fazaZlecen: FazaOkna;
  powodZlecen: string;
  fazaDziennika: FazaOkna;
  powodDziennika: string;
  /** Czy rdzeń odezwał się już w sprawie zleceń; puste przed pytaniem i po odpowiedzi to różne sytuacje. */
  pytanoOZlecenia: boolean;
  /** To samo dla dziennika działań — czytany osobno, więc i znany osobno. */
  pytanoODziennik: boolean;
  /** To samo dla przydziału okna modułu; pusty napis zanim padło pytanie i gdy rdzeń okna nie oddał. */
  pytanoOOkno: boolean;
}

export function utworzZapisModulu(): ZapisModulu {
  return {
    zlecenia: [],
    zleceniaObce: [],
    dziennik: [],
    okno: '',
    sesja: '',
    wybor: '',
    fazaZlecen: 'puste',
    powodZlecen: '',
    fazaDziennika: 'puste',
    powodDziennika: '',
    pytanoOZlecenia: false,
    pytanoODziennik: false,
    pytanoOOkno: false,
  };
}

/** Wciąga zlecenia przysłane przez rdzeń po odczycie i po zdarzeniu; obraz starszy nie zdejmuje obrazu nowszego. */
export function wchlonZlecenia(zapis: ZapisModulu, przyslane: readonly AssistantAction[]): void {
  zapis.zlecenia = zlozWykaz(zapis.zlecenia, przyslane);
  // Zdarzenie rdzenia jest odpowiedzią tak samo jak odczyt: pusty wykaz po nim już coś znaczy.
  zapis.pytanoOZlecenia = true;
  przeliczFazeZlecen(zapis);
}

/** Wciąga zlecenie spoza okna modułu z tej samej sesji, nie ruszając znacznika pytania o zlecenia okna modułu. */
export function wchlonObce(zapis: ZapisModulu, przyslane: readonly AssistantAction[]): void {
  zapis.zleceniaObce = zlozWykaz(zapis.zleceniaObce, przyslane);
  przeliczFazeZlecen(zapis);
}

/** Zdejmuje zlecenie usunięte przez rdzeń wraz z zawężeniem dziennika i wyborem, jeśli dotyczyły tego zlecenia. */
export function usunZlecenie(zapis: ZapisModulu, idZlecenia: string): void {
  zapis.zlecenia = zapis.zlecenia.filter((wpis) => wpis.id !== idZlecenia);
  zapis.zleceniaObce = zapis.zleceniaObce.filter((wpis) => wpis.id !== idZlecenia);
  if (zapis.wybor === idZlecenia) zapis.wybor = '';
  przeliczFazeZlecen(zapis);
}

/** Wciąga przysłane zlecenia w istniejący wykaz, nie zdejmując obrazu nowszego niż ten, który już tam stoi. */
function zlozWykaz(
  wykaz: readonly AssistantAction[],
  przyslane: readonly AssistantAction[],
): AssistantAction[] {
  let wynik = [...wykaz];
  for (const zlecenie of przyslane) {
    const znane = wynik.find((wpis) => wpis.id === zlecenie.id);
    if (znane === undefined) wynik = [...wynik, zlecenie];
    else if (znane.updatedAt <= zlecenie.updatedAt) {
      wynik = wynik.map((wpis) => (wpis.id === zlecenie.id ? zlecenie : wpis));
    }
  }
  return wynik;
}

/** Faza wykazu zleceń liczona z obu wykazów naraz, bo widoczna tabela zleceń obcych zaprzecza fazie pustej. */
function przeliczFazeZlecen(zapis: ZapisModulu): void {
  const ile = zapis.zlecenia.length + zapis.zleceniaObce.length;
  zapis.fazaZlecen = ile === 0 ? 'puste' : 'gotowe';
}

/** Odczyt stanu zleceń, gdzie niepowodzenie nie rzuca wyjątkiem, tylko zostawia fazę błędu wraz z powodem. */
export async function odczytajZlecenia(zapis: ZapisModulu, zrodlo: ZrodloAssistant): Promise<void> {
  const wynik = await zrodlo.zlecenia(zapis.okno);
  zapis.pytanoOZlecenia = true;
  if (!wynik.udany || wynik.wynik === undefined) {
    zapis.fazaZlecen = 'blad';
    zapis.powodZlecen = opisOdmowy(
      'Odczyt zleceń asystenta',
      wynik.blad?.code,
      wynik.blad?.message,
    );
    return;
  }
  zapis.zlecenia = wynik.wynik.actions;
  // Rozłączność wykazów: zlecenie okna modułu wraca tu, inaczej stałoby w tabeli dwa razy.
  zapis.zleceniaObce = zapis.zleceniaObce.filter(
    (wpis) => wpis.windowId !== zapis.okno && !zapis.zlecenia.some((swoje) => swoje.id === wpis.id),
  );
  const znane = [...zapis.zlecenia, ...zapis.zleceniaObce];
  if (zapis.wybor !== '' && !znane.some((wpis) => wpis.id === zapis.wybor)) {
    zapis.wybor = '';
  }
  przeliczFazeZlecen(zapis);
  zapis.powodZlecen = '';
}

/** Odczyt dziennika; zawężenie bierze się z wybranego zlecenia, a niepowodzenie zostawia fazę błędu z powodem. */
export async function odczytajDziennik(zapis: ZapisModulu, zrodlo: ZrodloAssistant): Promise<void> {
  const wynik = await zrodlo.dziennik(zapis.okno, zapis.wybor);
  zapis.pytanoODziennik = true;
  if (!wynik.udany || wynik.wynik === undefined) {
    zapis.fazaDziennika = 'blad';
    zapis.powodDziennika = opisOdmowy(
      'Odczyt dziennika działań',
      wynik.blad?.code,
      wynik.blad?.message,
    );
    return;
  }
  zapis.dziennik = wynik.wynik.entries;
  zapis.fazaDziennika = zapis.dziennik.length === 0 ? 'puste' : 'gotowe';
  zapis.powodDziennika = '';
}
