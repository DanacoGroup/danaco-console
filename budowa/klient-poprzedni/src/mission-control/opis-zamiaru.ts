import { ETYKIETA_DZIALANIA_KOLEJKI } from './etykiety-pulpitu';
import type {
  WejscieDoSesji,
  ZamiarDecyzji,
  ZamiarKolejki,
  ZamiarUtworzenia,
} from './zdarzenia-pulpitu';

/**
 * Przełożenie ładunku zdarzenia pulpitu na zdanie widoczne w pasku ostatniego
 * działania.
 *
 * Zdanie podaje nazwę komendy kontraktu, gdy zamiar ją ma; zamiar bez
 * odpowiednika w kontrakcie („podnieś priorytet") mówi o tym wprost.
 */

/** Zdanie o wejściu do sesji. */
export function opiszWejscie(wejscie: WejscieDoSesji): string {
  return `wejście do sesji „${wejscie.tytul}" (${wejscie.srodowisko}) — ${wejscie.komenda}`;
}

/** Zdanie o działaniu na kolejce. */
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

/** Zdanie o utworzeniu bytu. */
export function opiszZamiarUtworzenia(zamiar: ZamiarUtworzenia): string {
  return `${zamiar.etykieta} — rodzaj „${zamiar.rodzaj}"`;
}

/** Zdanie o wezwaniu do decyzji. */
export function opiszZamiarDecyzji(zamiar: ZamiarDecyzji): string {
  return zamiar.przeplywy > 0
    ? `rozstrzygnięcie ${zamiar.przeplywy} wstrzymanych przepływów, począwszy od „${zamiar.najstarszy}"`
    : 'przegląd przepływów — żaden nie czeka na decyzję';
}
