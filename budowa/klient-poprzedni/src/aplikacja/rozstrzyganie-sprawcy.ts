import type { Envelope, Window } from '../../../shared/contract';
import { odciskOkna } from './ustawienia-okna-sledzone';

/** Rozstrzyganie sprawcy — odpowiedź na jedno pytanie: czy ta zmiana wyszła z tego połączenia, czy z innego. */
export interface RozstrzyganieSprawcy {
  // Zbiera dowody własności z koperty odpowiedzi; parametr zasiewa własny stan wyjściowy.
  zapamietajWlasne(koperta: Envelope, naOknoZOdpowiedzi?: (okno: Window) => void): void;
  /** Rozstrzyga po zwłoce, czy byt o tym identyfikatorze jest cudzy. */
  poZwloce(id: string, gdyCudze: () => void): void;
  /** Rozstrzyga po zwłoce, czy TA zmiana okna wyszła z innego połączenia. */
  poZwloceOkno(okno: Window, gdyCudze: () => void): void;
  /** Zdejmuje zegary rozstrzygnięć jeszcze niezapadłych. */
  rozlacz(): void;
}

/** Ile czasu zdarzenie czeka na własną odpowiedź, zanim zostanie ostatecznie uznane za cudze — 700 milisekund. */
const ZWLOKA_ROZSTRZYGNIECIA_MS = 700;

/** Jak długo odcisk własnej odpowiedzi zaświadcza o własności — 1500 milisekund, nieco ponad zwłokę rozstrzygnięcia. */
const PAMIEC_ODCISKU_MS = 1500;

export function utworzRozstrzyganieSprawcy(): RozstrzyganieSprawcy {
  /** Identyfikatory wiadomości, które TO połączenie dostało w odpowiedzi. */
  const wlasne = new Set<string>();
  /** Odciski stanów okna oddanych w odpowiedziach na własne komendy. */
  const odciskiWlasne = new Map<string, number>();
  const zegary = new Set<number>();

  // Czy odcisk pochodzi z własnej odpowiedzi sprzed chwili i mieści się w pamięci odcisku — i zużywa go.
  function odciskSwiezy(odcisk: string): boolean {
    const chwila = odciskiWlasne.get(odcisk);
    if (chwila === undefined) return false;
    odciskiWlasne.delete(odcisk);
    return Date.now() - chwila <= PAMIEC_ODCISKU_MS;
  }

  /** Odkłada werdykt o zwłokę; `sprawdz` orzeka, czy czynność była własna. */
  function odlozWerdykt(sprawdz: () => boolean, gdyCudze: () => void): void {
    if (sprawdz()) return;
    const zegar = window.setTimeout(() => {
      zegary.delete(zegar);
      if (sprawdz()) return;
      gdyCudze();
    }, ZWLOKA_ROZSTRZYGNIECIA_MS);
    zegary.add(zegar);
  }

  return {
    // Odczyt jest świadomie tolerancyjny — kształt spoza kontraktu psuje tylko to jedno rozstrzygnięcie.
    zapamietajWlasne(koperta, naOknoZOdpowiedzi) {
      if (koperta.status === undefined) return;
      const tresc = koperta.payload;
      if (typeof tresc !== 'object' || tresc === null) return;
      const pola = tresc as Record<string, unknown>;

      const wiadomosc = pola['message'];
      if (typeof wiadomosc === 'object' && wiadomosc !== null) {
        const id = (wiadomosc as Record<string, unknown>)['id'];
        if (typeof id === 'string' && id.length > 0) wlasne.add(id);
      }

      const okno = pola['window'];
      if (typeof okno !== 'object' || okno === null) return;
      const kandydat = okno as Partial<Window>;
      if (typeof kandydat.id !== 'string' || kandydat.id.length === 0) return;
      const pelne = kandydat as Window;
      odciskiWlasne.set(odciskOkna(pelne), Date.now());
      naOknoZOdpowiedzi?.(pelne);
    },

    poZwloce(id, gdyCudze) {
      odlozWerdykt(() => wlasne.has(id), gdyCudze);
    },

    poZwloceOkno(okno, gdyCudze) {
      const odcisk = odciskOkna(okno);
      odlozWerdykt(() => odciskSwiezy(odcisk), gdyCudze);
    },

    rozlacz() {
      for (const zegar of zegary) window.clearTimeout(zegar);
      zegary.clear();
    },
  };
}
