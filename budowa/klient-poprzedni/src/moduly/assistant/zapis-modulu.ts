import type { AssistantAction, AssistantActivityEntry } from '../../../../shared/contract';
import type { FazaOkna } from '../../komponenty/faza-okna';
import { opisOdmowy } from '../../komponenty/odmowa';
import type { ZrodloAssistant } from './zrodlo-assistant';

/**
 * Zapis modułu Assistant — dane trzech okien wraz z przejściami między fazami
 * odczytu.
 *
 * Plik odpowiada wyłącznie za to, co moduł wie i jak ta wiedza zmienia się po
 * odpowiedzi rdzenia. Powiadamianie widoków i subskrypcja zdarzeń mieszkają
 * osobno, w `stan-assistant.ts`; rozdzielenie pozwala sprawdzić przejścia bez
 * budowania widoku.
 *
 * Zapis jest zmieniany w miejscu, a nie odtwarzany: trzy okna trzymają na niego
 * jedno odwołanie i patrzą na te same zlecenia.
 */
export interface ZapisModulu {
  zlecenia: AssistantAction[];
  /**
   * Zlecenia asystenta tej samej sesji założone poza oknem modułu.
   *
   * Wykaz stoi osobno, bo `odczytajZlecenia` czyta `assistant.action.status`
   * po `windowId` i podmienia nim cały wykaz — po takim odczycie zlecenie
   * z innego okna zniknęłoby bez śladu, choć wciąż biegnie. Zbiory są rozłączne
   * z definicji (`windowId` równy oknu modułu albo różny), a rozdział pozwala
   * jednemu z nich być podmienianym, a drugiemu trwać.
   *
   * Wpisy biorą się wyłącznie ze zdarzenia `assistant.action.changed`, które
   * rdzeń rozgłasza do wszystkich połączeń konta
   * (`transport/rozgloszenie.go`). Najczęstszą drogą jest `aod.voice.command`:
   * nakładka dobiera okno sama (`adapter_modul_aod.go`, `oknoZadania` — okno
   * wskazane wprost, okno karty sesji, ostatni punkt pracy) i nie sprawdza
   * przy tym `moduleId`, więc zlecenie zwykle siada na oknie rozmowy, nie na
   * oknie asystenta.
   *
   * Wykaz nie obejmuje zleceń spoza okna założonych, zanim moduł został
   * otwarty: kontrakt nie ma odczytu zleceń po sesji —
   * `assistant.action.status` czyta po `windowId` albo po `actionId`
   * (`adapter_modul_asystent_czynnosci.go`) — więc pierwszym momentem, w którym
   * klient dowiaduje się o takim zleceniu, jest jego pierwsza zmiana stanu.
   */
  zleceniaObce: AssistantAction[];
  dziennik: AssistantActivityEntry[];
  /** Okno modułu przypisane przez rdzeń; pusty napis znaczy brak. */
  okno: string;
  /**
   * Karta sesji, w której moduł został wczytany.
   *
   * Okna zarządcze pytają rdzeń o byty przypisane do karty sesji, nie do okna:
   * `memory.list` i `memory.toggle` biorą `sessionId`, a `session.tool.*`
   * czynią to samo dla dołożeń narzędzi. Kanał niesie własną sesję
   * (`Kanal.sesja()`), ale jest to sesja połączenia, nie karta, dla której
   * powłoka wczytała moduł — sięganie po nią zamieniłoby pamięć jednej karty
   * w pamięć innej.
   */
  sesja: string;
  /** Zlecenie zawężające dziennik; pusty napis znaczy całość zapisu. */
  wybor: string;
  fazaZlecen: FazaOkna;
  powodZlecen: string;
  fazaDziennika: FazaOkna;
  powodDziennika: string;
  /**
   * Czy rdzeń odezwał się już w sprawie zleceń.
   *
   * Pusty wykaz przed pytaniem i pusty wykaz po odpowiedzi to dwie różne
   * sytuacje, a faza okna ich nie rozróżnia — obie są `puste`. Rozróżnienie
   * należy do źródła danych i mieszka tutaj; okno bierze z niego wyłącznie
   * zdanie, którym pustkę nazywa.
   */
  pytanoOZlecenia: boolean;
  /** To samo dla dziennika działań — czytany osobno, więc i znany osobno. */
  pytanoODziennik: boolean;
  /**
   * To samo dla przydziału okna modułu (`window.list`).
   *
   * Pole `okno` niesie pusty napis w dwóch różnych sytuacjach: zanim padło
   * pytanie i wtedy, gdy rdzeń okna modułu nie oddał. Dla Voice Console jest to
   * różnica między „zaraz zacznę" a „nie mam dokąd wysłać polecenia" — czyli
   * między stanem pustym a stanem błędu. Bez tego znacznika okno musiałoby ją
   * zgadywać z upływu czasu.
   */
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

/**
 * Wciąga zlecenia przysłane przez rdzeń — po odczycie i po zdarzeniu.
 *
 * Obraz starszy nie zdejmuje nowszego. Odpowiedź na własną komendę i zdarzenie
 * o tym samym zleceniu jadą dwoma niezależnymi biegami rdzenia
 * (`transport/petla_odbioru.go`), więc odpowiedź z chwili wcześniejszej może
 * przyjść po zdarzeniu z chwili późniejszej. Wpuszczenie jej zamieniłoby
 * zlecenie już wykonane z powrotem w „w kolejce". O tym, który obraz jest
 * późniejszy, mówi znacznik `updatedAt`.
 */
export function wchlonZlecenia(zapis: ZapisModulu, przyslane: readonly AssistantAction[]): void {
  zapis.zlecenia = zlozWykaz(zapis.zlecenia, przyslane);
  // Zdarzenie rdzenia jest odpowiedzią tak samo jak odczyt: po nim wykaz pusty
  // znaczy już „rdzeń nie prowadzi zlecenia", a nie „jeszcze nie pytałem".
  zapis.pytanoOZlecenia = true;
  przeliczFazeZlecen(zapis);
}

/**
 * Wciąga zlecenie spoza okna modułu — z tej samej sesji.
 *
 * Znacznik `pytanoOZlecenia` zostaje nietknięty: zdarzenie o cudzym oknie nie
 * jest odpowiedzią rdzenia w sprawie zleceń okna modułu, więc po nim wciąż nie
 * wiadomo, czy okno modułu prowadzi coś swojego. Postawienie znacznika
 * zamieniłoby „jeszcze nie pytałem" w „rdzeń nie prowadzi ani jednego" bez
 * odczytu, który by za tym zdaniem stał.
 *
 * Zlecenie okna modułu przysłane tą drogą — bo zdarzenie przyszło, zanim
 * `ustalOkno` poznało okno — nie zostaje w wykazie obcym: `odczytajZlecenia`
 * czyści z niego wszystko, co należy już do okna modułu.
 */
export function wchlonObce(zapis: ZapisModulu, przyslane: readonly AssistantAction[]): void {
  zapis.zleceniaObce = zlozWykaz(zapis.zleceniaObce, przyslane);
  przeliczFazeZlecen(zapis);
}

/** Zdejmuje zlecenie usunięte przez rdzeń wraz z zawężeniem dziennika. */
export function usunZlecenie(zapis: ZapisModulu, idZlecenia: string): void {
  zapis.zlecenia = zapis.zlecenia.filter((wpis) => wpis.id !== idZlecenia);
  zapis.zleceniaObce = zapis.zleceniaObce.filter((wpis) => wpis.id !== idZlecenia);
  if (zapis.wybor === idZlecenia) zapis.wybor = '';
  przeliczFazeZlecen(zapis);
}

/**
 * Wciąga przysłane zlecenia w istniejący wykaz, nie zdejmując obrazu nowszego.
 *
 * Wydzielone z `wchlonZlecenia`, bo ta sama reguła pierwszeństwa obowiązuje oba
 * wykazy — a rozjazd między nimi znaczyłby, że zlecenie spoza okna cofa się
 * w czasie, gdy zlecenie okna modułu tego nie robi.
 */
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

/**
 * Faza wykazu zleceń liczona z obu wykazów.
 *
 * Okno stojące w fazie pustej nad widoczną tabelą zleceń spoza okna mówiłoby
 * dwie rzeczy naraz — a zdanie pustki („rdzeń nie prowadzi dziś ani jednego")
 * byłoby wtedy po prostu nieprawdą.
 */
function przeliczFazeZlecen(zapis: ZapisModulu): void {
  const ile = zapis.zlecenia.length + zapis.zleceniaObce.length;
  zapis.fazaZlecen = ile === 0 ? 'puste' : 'gotowe';
}

/**
 * Odczyt stanu zleceń.
 *
 * Niepowodzenie nie rzuca wyjątkiem: zostawia fazę `blad` wraz z powodem, żeby
 * pusty wykaz po odmowie rdzenia i pusty wykaz na świeżej instalacji pozostały
 * dwoma różnymi stanami.
 */
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
  // Rozłączność wykazów: zlecenie okna modułu, które wpadło do wykazu obcego
  // zdarzeniem przed poznaniem okna, wraca tu na swoje miejsce — inaczej ten sam
  // wiersz stałby w tabeli dwa razy.
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

/** Odczyt dziennika; zawężenie bierze się z wybranego zlecenia. */
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
