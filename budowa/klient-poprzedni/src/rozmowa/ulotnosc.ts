import { profilModulu } from '../okno-komunikacji/rejestr-profilow';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';

/**
 * Ulotność rozmowy — moduł, którego czat nie ma pamięci sesyjnej.
 *
 * Źródłem jest profil modułu (`profil-modulu.ts`, pole `pamiecSesyjna`), nie
 * wykaz kodów prowadzony tutaj.
 *
 * Ulotność to dwa czyszczenia, nie jedno:
 *   (a) zamknięcie okna — rozmowa nie wraca, bo historia nie jest odtwarzana
 *       z rdzenia (`message.list` nie idzie w ogóle);
 *   (b) zmiana kontekstu roboczego — na przykład zmiana testowanego eksperta.
 *
 * Oba muszą być widoczne: czat, który po cichu gubi wątek, czyta się jak awaria
 * aplikacji. Dlatego polityka niesie gotowe zdania, nie samą wartość logiczną,
 * a rozmowa wypisuje je w oknie.
 *
 * Zasięg ulotności kończy się na kliencie:
 *   • `adapter_rozmowa.go` → `Wyslij` dopisuje każdą wiadomość Operatora
 *     i każdą odpowiedź modelu do dziennika rozmowy (`dziennik_rozmowy.go`),
 *     a ten pisze do tabeli `wiadomosc` (`dane/wiadomosci.go`). Zapis nie pyta
 *     o moduł okna.
 *   • Kontrakt nie ma komendy kasowania wiadomości: obszar `message.*` to
 *     `send`, `stop`, `list` (`shared/contract.ts`).
 *   • `window.close` zmienia stan okna; wierszy wiadomości nie usuwa.
 * „Bez pamięci sesyjnej" znaczy więc: klient nie odtwarza i nie pokazuje, a
 * rdzeń zapis trzyma. Okno mówi o tym Operatorowi wprost.
 */

/**
 * Kontekst roboczy rozmowy ulotnej — to, czego zmiana zaczyna nową rozmowę.
 *
 * Port jest ogólny: warstwa rozmowy nie zna ani modułu Agents, ani pojęcia
 * „ekspert". Wie tylko, że kontekst ma klucz, a zmiana klucza kończy
 * dotychczasową rozmowę i przychodzi ze zdaniem dla Operatora.
 */
export interface KontekstRoboczy {
  /** Klucz kontekstu; pusty znaczy „kontekstu jeszcze nie wskazano". */
  klucz(): string;
  /**
   * Subskrypcja zmiany kontekstu. Słuchacz dostaje gotowe zdanie mówiące, co
   * się zmieniło — nie sam klucz, bo z klucza nie da się złożyć zdania
   * prawdziwego bez wiedzy o dziedzinie.
   */
  naZmiane(sluchacz: (zdanie: string) => void): Odsubskrybuj;
}

/** Reguła pamięci rozmowy jednego okna wraz ze zdaniami dla Operatora. */
export interface PolitykaUlotnosci {
  /** Kod modułu, którego polityka dotyczy. */
  modul: string;
  /** Czy rozmowa przeżywa zamknięcie okna; fałsz znaczy „rozmowa ulotna". */
  pamiecSesyjna: boolean;
  /**
   * Zdania stanu — wypisywane przy otwarciu rozmowy ulotnej i powtarzane po
   * każdym czyszczeniu. Puste dla rozmowy z pamięcią: nie ma czego zapowiadać.
   */
  zapowiedz: readonly string[];
  /** Kontekst, którego zmiana czyści rozmowę; `null` = takiego kontekstu nie ma. */
  kontekst: KontekstRoboczy | null;
}

/**
 * Zdanie mówiące, dokąd sięga ulotność, a dokąd nie. Bez niego określenie
 * „czat roboczy" czytałoby się jako obietnica kasowania, którego rdzeń nie
 * wykonuje.
 */
export const ZDANIE_O_ZAPISIE_RDZENIA =
  'Zasięg: okno nie odtwarza i nie pokazuje tej rozmowy ponownie. Rdzeń zapisuje ' +
  'wiadomości w swojej bazie mimo to — kontrakt nie ma dziś komendy ich kasowania ' +
  '(obszar message.* to send, stop, list).';

/** Polityka rozmowy z pamięcią — stan zwykły, bez ani jednego zdania w oknie. */
export function politykaTrwala(modul: string): PolitykaUlotnosci {
  return { modul, pamiecSesyjna: true, zapowiedz: [], kontekst: null };
}

/**
 * Polityka okna pracującego we wskazanym module.
 *
 * Regułę bierze z profilu modułu, żeby nie prowadzić drugiej listy modułów bez
 * pamięci, która rozjechałaby się z profilem przy pierwszej jego zmianie.
 *
 * Kontekst roboczy podaje warstwa składająca, bo tylko ona zna moduły. Moduł
 * bez pamięci sesyjnej i bez podanego kontekstu jest wciąż ulotny — traci
 * rozmowę przy zamknięciu okna — mówi tylko o jednym czyszczeniu zamiast dwóch.
 */
export function politykaModulu(
  modul: string,
  kontekst: KontekstRoboczy | null = null,
): PolitykaUlotnosci {
  const profil = profilModulu(modul);
  if (profil.pamiecSesyjna) return politykaTrwala(modul);

  const zapowiedz = [
    `Moduł ${profil.nazwa}: czat roboczy BEZ pamięci sesyjnej. Ten wątek nie zostanie ` +
      'odtworzony po zamknięciu okna' +
      (kontekst === null ? '.' : ' ani po zmianie kontekstu testowania.'),
    ZDANIE_O_ZAPISIE_RDZENIA,
  ];

  return { modul, pamiecSesyjna: false, zapowiedz, kontekst };
}

/** Czy polityka każe rozmowie być ulotną. */
export function czyUlotna(polityka: PolitykaUlotnosci | null): boolean {
  return polityka !== null && !polityka.pamiecSesyjna;
}
