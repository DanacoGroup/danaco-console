import {
  StudioIngestState,
  type StudioIngestItem,
  type StudioLayoutBlock,
  type StudioRecognizedWord,
} from '../../../../shared/contract';

/**
 * Kolejka wczytywania — od dziś kolejka RDZENIA, nie kliencka.
 *
 * ── Co się zmieniło i dlaczego ──────────────────────────────────────────────
 * Kolejka była prowadzona w oknie i ginęła z jego odświeżeniem, bo panel stał na
 * `document.text.extract`: jedno wydobycie na jedno wywołanie, bez pojęcia
 * pozycji. Rodzina `studio.ingest.*` niesie kolejkę po stronie rdzenia —
 * `queue.add` dokłada materiał, `queue.list` oddaje pozycje wraz ze stanem
 * i wynikiem. Ten plik przestał więc być kolejką, a stał się jej ODBICIEM:
 * pamięcią tego, co rdzeń ostatnio oddał, wraz z tym, czego rdzeń nie prowadzi.
 *
 * ── Czego rdzeń nie prowadzi i co zostaje po stronie okna ───────────────────
 *   — wskazanie pozycji, na której Operator pracuje (to nastawa widoku),
 *   — słowa rozpoznane i bloki układu odebrane z `ingest.recognize`: kontrakt
 *     oddaje je w odpowiedzi rozpoznania, a `queue.list` ich nie powtarza, więc
 *     bez odłożenia tutaj poprawianie słów nie miałoby na czym pracować,
 *   — znacznik wywołania w toku, żeby dwa naciśnięcia nie poszły naraz.
 *
 * Stan pozycji NIE jest już liczony w oknie: przychodzi z rdzenia jako
 * `StudioIngestState` i jest pięciowartościowy — oczekuje, przetwarzanie,
 * gotowa, ponowienie (pewność poniżej progu), odmowa. „Ponowienie" jest tu
 * stanem osobnym i tak ma być pokazane: pozycja nie jest ani gotowa, ani
 * odmówiona, tylko wraca do rozpoznania.
 */

/** Odbicie kolejki rdzenia wraz z tym, czego rdzeń nie prowadzi. */
export interface KolejkaCyfryzacji {
  /** Pozycje w kolejności oddanej przez rdzeń. */
  pozycje(): readonly StudioIngestItem[];
  /** Wstawia pozycje oddane przez rdzeń; zastępuje odbicie w całości. */
  ustawPozycje(pozycje: readonly StudioIngestItem[]): void;
  /** Wstawia albo podmienia jedną pozycję — po rozpoznaniu i po poprawce. */
  ustawPozycje1(pozycja: StudioIngestItem): void;
  /** Dokłada pozycje oddane przez rdzeń po dołożeniu materiału. */
  dolozPozycje(pozycje: readonly StudioIngestItem[]): void;
  pozycja(idPozycji: string): StudioIngestItem | null;
  /** Pierwsza pozycja oczekująca albo wracająca do ponowienia. */
  nastepnaDoRozpoznania(): StudioIngestItem | null;
  /** Pozycja wskazana przez Operatora; pusty napis znaczy „żadna". */
  wskazana(): string;
  ustawWskazana(idPozycji: string): void;
  /** Słowa rozpoznane pozycji wraz z pewnością — materiał poprawiania. */
  slowa(idPozycji: string): readonly StudioRecognizedWord[];
  /** Bloki układu pozycji — kolumny, tabele, nagłówki, stopki, przypisy. */
  uklad(idPozycji: string): readonly StudioLayoutBlock[];
  /** Odkłada wynik rozpoznania: słowa i układ, których `queue.list` nie powtarza. */
  ustawRozpoznanie(
    idPozycji: string,
    slowa: readonly StudioRecognizedWord[] | undefined,
    uklad: readonly StudioLayoutBlock[] | undefined,
  ): void;
  /** Czy na tej pozycji stoi wywołanie w toku. */
  wToku(idPozycji: string): boolean;
  ustawWToku(idPozycji: string, trwa: boolean): void;
  /** Liczby pozycji w rozbiciu na stany — materiał paska podsumowania. */
  bilans(): BilansKolejki;
}

/** Liczby pozycji w rozbiciu na pięć stanów kontraktu. */
export interface BilansKolejki {
  wszystkie: number;
  oczekujace: number;
  przetwarzane: number;
  gotowe: number;
  ponowienia: number;
  odmowy: number;
  /** Ile pozycji ma tekst gotowy do przyjęcia do edytora. */
  zTekstem: number;
}

export function utworzKolejkeCyfryzacji(): KolejkaCyfryzacji {
  let pozycje: StudioIngestItem[] = [];
  let wskazana = '';
  const slowa = new Map<string, readonly StudioRecognizedWord[]>();
  const uklady = new Map<string, readonly StudioLayoutBlock[]>();
  const wToku = new Set<string>();

  function znajdz(idPozycji: string): StudioIngestItem | null {
    return pozycje.find((pozycja) => pozycja.id === idPozycji) ?? null;
  }

  function ile(stan: StudioIngestState): number {
    return pozycje.filter((pozycja) => pozycja.state === stan).length;
  }

  return {
    pozycje: () => pozycje,

    ustawPozycje(nowe) {
      pozycje = [...nowe];
      // Wskazanie przeżywa odczyt, o ile pozycja nadal jest: odświeżenie kolejki
      // nie może przestawiać Operatora na inną pozycję niż ta, którą poprawia.
      if (wskazana !== '' && !pozycje.some((pozycja) => pozycja.id === wskazana)) wskazana = '';
    },

    ustawPozycje1(nowa) {
      const miejsce = pozycje.findIndex((pozycja) => pozycja.id === nowa.id);
      if (miejsce >= 0) pozycje[miejsce] = nowa;
      else pozycje.push(nowa);
    },

    dolozPozycje(nowe) {
      for (const pozycja of nowe) {
        const miejsce = pozycje.findIndex((istniejaca) => istniejaca.id === pozycja.id);
        if (miejsce >= 0) pozycje[miejsce] = pozycja;
        else pozycje.push(pozycja);
      }
      const pierwsza = nowe[0];
      if (wskazana === '' && pierwsza !== undefined) wskazana = pierwsza.id;
    },

    pozycja: (idPozycji) => znajdz(idPozycji),

    nastepnaDoRozpoznania: () =>
      pozycje.find(
        (pozycja) =>
          pozycja.state === StudioIngestState.Oczekuje ||
          pozycja.state === StudioIngestState.Ponowienie,
      ) ?? null,

    wskazana: () => wskazana,

    ustawWskazana(idPozycji) {
      wskazana = idPozycji;
    },

    slowa: (idPozycji) => slowa.get(idPozycji) ?? [],
    uklad: (idPozycji) => uklady.get(idPozycji) ?? [],

    ustawRozpoznanie(idPozycji, noweSlowa, nowyUklad) {
      // Brak słów w odpowiedzi NIE zeruje słów odłożonych: rozpoznanie bez
      // odtwarzania układu jest wynikiem poprawnym, a poprawki Operatora
      // wniesione wcześniej nie mają za co przepadać.
      if (noweSlowa !== undefined) slowa.set(idPozycji, [...noweSlowa]);
      if (nowyUklad !== undefined) uklady.set(idPozycji, [...nowyUklad]);
    },

    wToku: (idPozycji) => wToku.has(idPozycji),

    ustawWToku(idPozycji, trwa) {
      if (trwa) wToku.add(idPozycji);
      else wToku.delete(idPozycji);
    },

    bilans() {
      return {
        wszystkie: pozycje.length,
        oczekujace: ile(StudioIngestState.Oczekuje),
        przetwarzane: ile(StudioIngestState.Przetwarzanie),
        gotowe: ile(StudioIngestState.Gotowa),
        ponowienia: ile(StudioIngestState.Ponowienie),
        odmowy: ile(StudioIngestState.Odmowa),
        zTekstem: pozycje.filter((pozycja) => (pozycja.text ?? '') !== '').length,
      };
    },
  };
}
