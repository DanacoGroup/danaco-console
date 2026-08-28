/**
 * Przełożenie ładunku zdarzenia pulpitu na zdanie widoczne w pasku ostatniego
 * działania. Zdanie podaje nazwę komendy kontraktu, gdy zamiar ją ma; zamiar
 * bez odpowiednika w kontrakcie mówi o tym wprost.
 */

import { ETYKIETA_DZIALANIA_KOLEJKI } from './etykiety-pulpitu';
import type {
  WejscieDoSesji,
  ZamiarDecyzji,
  ZamiarKolejki,
  ZamiarUtworzenia,
} from './zdarzenia-pulpitu';

/**
 * Składa zdanie o wejściu do sesji: podaje tytuł sesji, jej środowisko oraz
 * nazwę komendy kontraktu, która wejście realizuje.
 */
export function opiszWejscie(wejscie: WejscieDoSesji): string {
  return `wejście do sesji „${wejscie.tytul}" (${wejscie.srodowisko}) — ${wejscie.komenda}`;
}

/**
 * Składa zdanie o działaniu na kolejce roli. Rodzaj zamiaru rozstrzyga treść:
 * działanie kolejki oraz przekazanie kontekstu podają komendę kontraktu,
 * natomiast podniesienie priorytetu odpowiednika w kontrakcie nie ma.
 */
export function opiszZamiarKolejki(zamiar: ZamiarKolejki): string {
  const rola = zamiar.rola;

  switch (zamiar.rodzaj) {
    case 'kolejka':
      return `${ETYKIETA_DZIALANIA_KOLEJKI[zamiar.dzialanie]} — kolejka ${rola} — ${zamiar.komenda}`;
    case 'przekazanie':
      return `Przekazanie kontekstu — kolejka ${rola} — ${zamiar.komenda}`;
    case 'priorytet':
      return `Podniesienie priorytetu — kolejka ${rola} — brak odpowiednika w kontrakcie, zgłoszone do Rejestru`;
  }
}

/**
 * Składa zdanie o utworzeniu bytu z etykiety zamiaru oraz nazwy jego rodzaju.
 * Ładunek zdarzenia nie niesie tu komendy kontraktu, więc zdanie jej nie podaje.
 */
export function opiszZamiarUtworzenia(zamiar: ZamiarUtworzenia): string {
  return `${zamiar.etykieta} — rodzaj „${zamiar.rodzaj}"`;
}

/**
 * Składa zdanie o wezwaniu do decyzji. Przy niezerowej liczbie wstrzymanych
 * przepływów podaje ją wraz z nazwą najstarszego, a przy zerowej stwierdza,
 * że żaden przepływ na decyzję nie czeka.
 */
export function opiszZamiarDecyzji(zamiar: ZamiarDecyzji): string {
  return zamiar.przeplywy > 0
    ? `rozstrzygnięcie ${zamiar.przeplywy} wstrzymanych przepływów, począwszy od „${zamiar.najstarszy}"`
    : 'przegląd przepływów — żaden nie czeka na decyzję';
}
