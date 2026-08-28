/**
 * Strumień wykonawcy widziany przez koordynatora po stronie widoku: odkładanie
 * fragmentów przychodzących zdarzeniem `stream.chunk`, żeby okno koordynatora
 * miało co narysować po przewinięciu.
 */
import { ChunkKind, type StreamChunkEvent } from '../../../../shared/contract';

/**
 * Jeden wpis strumienia jednego wykonawcy: okno pochodzenia, wiadomość tury,
 * rodzaj fragmentu wprost z kontraktu, treść tekstowa oraz chwila przyjęcia
 * liczona w milisekundach epoki.
 */
export interface WpisStrumienia {
  /** Okno wykonawcy, z którego wpis przyszedł. */
  okno: string;
  /** Wiadomość, w ramach której biegnie tura. */
  wiadomosc: string;
  /** Rodzaj fragmentu wprost z kontraktu. */
  rodzaj: ChunkKind;
  /** Treść tekstowa; dla fragmentów nietekstowych zapis skrócony. */
  tresc: string;
  /** Chwila przyjęcia wpisu w milisekundach epoki. */
  chwila: number;
}

/**
 * Ile ostatnich wpisów jednego wykonawcy zostaje jawnych. Starsze wypadają
 * z pamięci okna, ponieważ koordynator czyta ostatni odcinek pracy, a nie całą
 * historię tury.
 */
export const POJEMNOSC_STRUMIENIA = 200;

export interface StrumienWykonawcy {
  /** Dopisuje fragment; zwraca prawdę, gdy widok ma się przerysować. */
  dopisz(tresc: StreamChunkEvent, chwila: number): boolean;
  /** Wpisy wykonawcy w kolejności przyjęcia. */
  wpisy(okno: string): readonly WpisStrumienia[];
  /** Ostatnia treść tekstowa wykonawcy — nośnik przekazania wyniku. */
  wynik(okno: string): string;
  /** Liczba wywołań narzędzi wykonawcy — miara pracy wykonanej poza tekstem. */
  liczbaWywolan(okno: string): number;
  /** Porzuca ślad wykonawcy; wywoływane przy przepięciu obsady. */
  zapomnij(okno: string): void;
}

export function utworzStrumienWykonawcy(pojemnosc = POJEMNOSC_STRUMIENIA): StrumienWykonawcy {
  const wpisy = new Map<string, WpisStrumienia[]>();
  const tekst = new Map<string, string>();
  // Tekst narasta w obrębie jednej tury; nowa wiadomość zaczyna wynik od nowa.
  const tury = new Map<string, string>();

  function biezace(okno: string): WpisStrumienia[] {
    const gotowe = wpisy.get(okno);
    if (gotowe !== undefined) return gotowe;
    const nowe: WpisStrumienia[] = [];
    wpisy.set(okno, nowe);
    return nowe;
  }

  return {
    dopisz(tresc, chwila) {
      const zapis = zapisFragmentu(tresc);
      if (zapis === '') return false;
      const lista = biezace(tresc.windowId);
      lista.push({
        okno: tresc.windowId,
        wiadomosc: tresc.messageId,
        rodzaj: tresc.kind,
        tresc: zapis,
        chwila,
      });
      if (lista.length > pojemnosc) lista.splice(0, lista.length - pojemnosc);
      if (tresc.kind === ChunkKind.Text) {
        const tazSamaTura = tury.get(tresc.windowId) === tresc.messageId;
        tury.set(tresc.windowId, tresc.messageId);
        tekst.set(tresc.windowId, (tazSamaTura ? (tekst.get(tresc.windowId) ?? '') : '') + zapis);
      }
      return true;
    },

    wpisy: (okno) => wpisy.get(okno) ?? [],

    wynik: (okno) => tekst.get(okno) ?? '',

    liczbaWywolan: (okno) =>
      (wpisy.get(okno) ?? []).filter((wpis) => wpis.rodzaj === ChunkKind.ToolUse).length,

    zapomnij(okno) {
      wpisy.delete(okno);
      tekst.delete(okno);
      tury.delete(okno);
    },
  };
}

/**
 * Zapis fragmentu w postaci nadającej się do pokazania.
 *
 * Fragment nietekstowy niesie treść w polu `data` o kształcie zależnym od
 * rodzaju. Widok nie zgaduje jego struktury i pokazuje zapis JSON.
 */
function zapisFragmentu(tresc: StreamChunkEvent): string {
  if (tresc.text !== undefined && tresc.text !== '') return tresc.text;
  if (tresc.data === undefined || tresc.data === null) return '';
  try {
    return JSON.stringify(tresc.data);
  } catch {
    return String(tresc.data);
  }
}
