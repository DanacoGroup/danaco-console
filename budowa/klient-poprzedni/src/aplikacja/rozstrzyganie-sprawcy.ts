import type { Envelope, Window } from '../../../shared/contract';
import { odciskOkna } from './ustawienia-okna-sledzone';

/**
 * Rozstrzyganie sprawcy — odpowiedź na jedno pytanie: czy ta zmiana wyszła
 * z tego połączenia, czy z innego.
 *
 * Rdzeń rozgłasza każdą zmianę do wszystkich połączeń konta
 * (`transport/rozgloszenie.go`), nie pomijając nadawcy, a zdarzenie nie niesie
 * ani identyfikatora połączenia, ani autora: `WindowChangedEvent` ma okno,
 * `MessageChangedEvent` ma wiadomość. Kontrakt jest zamrożony, więc pola
 * sprawcy w nim nie przybędzie, a Operator ma odróżnić własny ruch od cudzego.
 *
 * Mechanizm odpowiada wyłącznie „to połączenie / nie to połączenie". Nie mówi
 * „asystent", bo drugie urządzenie Operatora wygląda stąd identycznie —
 * wołający mówi „spoza tego połączenia".
 *
 * Dwa rodzaje zapisu, bo dwa rodzaje pytania:
 *  • wiadomość powstaje raz i nigdy się nie zmienia, więc wystarczy
 *    identyfikator: wiadomość, której to połączenie nie dostało w odpowiedzi,
 *    napisał ktoś inny;
 *  • okno żyje długo i zmienia się wiele razy. Sam identyfikator odpowiadałby
 *    na pytanie „czy kiedykolwiek dotknąłem tego okna" zamiast „czy to
 *    przestawienie jest moje", więc zapisywany jest odcisk ustawień
 *    z odpowiedzi (`ustawienia-okna-sledzone.ts`).
 */
export interface RozstrzyganieSprawcy {
  /**
   * Zbiera dowody własności z KOPERTY — czynne wyłącznie dla odpowiedzi.
   *
   * @param naOknoZOdpowiedzi wołane dla okna niesionego przez odpowiedź; tędy
   *   wołający zasiewa swój własny stan wyjściowy, nie sięgając po drugą
   *   subskrypcję kanału.
   */
  zapamietajWlasne(koperta: Envelope, naOknoZOdpowiedzi?: (okno: Window) => void): void;
  /** Rozstrzyga po zwłoce, czy byt o tym identyfikatorze jest cudzy. */
  poZwloce(id: string, gdyCudze: () => void): void;
  /** Rozstrzyga po zwłoce, czy TA zmiana okna wyszła z innego połączenia. */
  poZwloceOkno(okno: Window, gdyCudze: () => void): void;
  /** Zdejmuje zegary rozstrzygnięć jeszcze niezapadłych. */
  rozlacz(): void;
}

/**
 * Ile czasu zdarzenie czeka na własną odpowiedź, zanim zostanie uznane za cudze.
 *
 * Zwłoka jest konieczna: rdzeń rozgłasza zdarzenie przed oddaniem odpowiedzi
 * (`handlers_window.go` — `e.okno(...)` przed `return`), więc zdarzenie
 * o własnej czynności zawsze wyprzedza własną odpowiedź. Bez zwłoki każde
 * kliknięcie Operatora meldowałoby się jako cudze.
 *
 * 700 ms pokrywa drogę powrotną własnej odpowiedzi z rdzenia lokalnego na
 * maszynie obciążonej. Krócej daje fałszywe „cudze" przy własnych
 * kliknięciach, dłużej rozjeżdża napis z tym, co widać na scenie. Zwłoka
 * dotyczy samego werdyktu: pokazanie okna i wpisu nie czeka na nic.
 */
const ZWLOKA_ROZSTRZYGNIECIA_MS = 700;

/**
 * Jak długo odcisk własnej odpowiedzi zaświadcza o własności.
 *
 * 1500 ms, czyli nieco ponad {@link ZWLOKA_ROZSTRZYGNIECIA_MS} — tyle dzieli
 * zdarzenie od odpowiedzi na tę samą komendę. Odcisk starszy nie jest już
 * świadkiem niczego, a szersze okno unieważnia cudze posunięcie: odcisk stanu
 * „okno w Studiu" zapisany przy starcie aplikacji zjadłby powrót asystenta do
 * Studia kilka sekund później i Operator zobaczyłby okno pracujące w innym
 * module, niż wskazuje kolumna.
 */
const PAMIEC_ODCISKU_MS = 1500;

export function utworzRozstrzyganieSprawcy(): RozstrzyganieSprawcy {
  /** Identyfikatory wiadomości, które TO połączenie dostało w odpowiedzi. */
  const wlasne = new Set<string>();
  /** Odciski stanów okna oddanych w odpowiedziach na własne komendy. */
  const odciskiWlasne = new Map<string, number>();
  const zegary = new Set<number>();

  /**
   * Czy odcisk pochodzi z własnej odpowiedzi sprzed chwili — i zużywa go.
   *
   * Odcisk opisuje stan, a nie czynność, więc ten sam stan potrafi wystąpić
   * dwa razy z różnych powodów (Operator przechodzi Studio → Workspace,
   * asystent wraca do Studia). Stąd dwa zabezpieczenia: wiek — patrz
   * {@link PAMIEC_ODCISKU_MS} — i zużycie. Jedna własna komenda rodzi dokładnie
   * jedno rozgłoszone zdarzenie, więc odcisk unieważnia jedno zdarzenie
   * i znika; drugie zdarzenie o tym samym stanie pochodzi z innej czynności
   * i nie ma już czym się wylegitymować.
   */
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
    // Odczyt jest świadomie tolerancyjny i płytki — kształt spoza kontraktu nie
    // psuje niczego poza tym jednym rozstrzygnięciem o sprawcy.
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
