import { profilModulu } from '../../okno-komunikacji/rejestr-profilow';
import { utworzMagistrale, type Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import { ZDANIE_O_ZAPISIE_RDZENIA, type KontekstRoboczy } from '../../rozmowa/ulotnosc';

/**
 * Kod modułu w katalogu rdzenia, ten sam, którym moduł opisuje się powłoce.
 * Pod nim stoi profil modułu niosący pole pamięci sesyjnej rozmowy.
 */
const KOD_MODULU = 'agents';

/**
 * Rozgłos testowanego eksperta: kontekst roboczy czatu testowego modułu Agents.
 * Niczego nie czyści — ogłasza wyłącznie zmianę testowanego eksperta wraz ze
 * zdaniem o niej, a wnioski wyciąga z tego polityka ulotności rozmowy.
 */
export interface RozglosTestowanego extends KontekstRoboczy {
  // Ustawia testowanego eksperta; ta sama para nie jest zmiana i nic nie oglasza.
  ustaw(idEksperta: string, nazwaEksperta: string): void;
  /** Nazwa testowanego eksperta; pusta, dopóki żaden nie jest wybrany. */
  nazwa(): string;
}

/**
 * Zdanie o zmianie testowanego eksperta.
 *
 * Zdanie nazywa poprzedniego eksperta, bo samo „rozmowa wyczyszczona" nie mówi
 * Operatorowi, co ją wyczyściło.
 */
function zdanieOZmianie(poprzedni: string, biezacy: string): string {
  return (
    `Rozmowa testowa wyczyszczona: testowany agent zmieniony z „${poprzedni}" na „${biezacy}". ` +
    'Nowy agent to nowy kontekst roboczy — poprzedni wątek do niego nie należy.'
  );
}

/**
 * Składa zdanie o zejściu testowanego eksperta, gdy moduł przestaje testować
 * kogokolwiek; zdanie nazywa poprzedniego eksperta i przypomina, że czat
 * roboczy dotyczy zawsze jednego wybranego eksperta.
 */
function zdanieOZejsciu(poprzedni: string): string {
  return (
    `Rozmowa testowa wyczyszczona: moduł nie testuje już agenta „${poprzedni}". ` +
    'Czat roboczy dotyczy zawsze jednego, wybranego eksperta.'
  );
}

export function utworzRozglosTestowanego(): RozglosTestowanego {
  const zmiany = utworzMagistrale<string>();
  let klucz = '';
  let nazwa = '';

  return {
    klucz: () => klucz,
    nazwa: () => nazwa,

    naZmiane: (sluchacz): Odsubskrybuj => zmiany.subskrybuj(sluchacz),

    ustaw(idEksperta, nazwaEksperta) {
      if (idEksperta === klucz) {
        // Sama nazwa nie jest zmiana kontekstu i rozmowy nie konczy.
        nazwa = nazwaEksperta;
        return;
      }
      const poprzedni = nazwa === '' ? klucz : nazwa;
      const bylWybor = klucz !== '';
      klucz = idEksperta;
      nazwa = nazwaEksperta;
      // Pierwszy wybor nie jest zmiana, bo modul wchodzi z wyborem pustym.
      if (!bylWybor) return;
      zmiany.oglos(
        idEksperta === ''
          ? zdanieOZejsciu(poprzedni)
          : zdanieOZmianie(poprzedni, nazwaEksperta === '' ? idEksperta : nazwaEksperta),
      );
    },
  };
}

/**
 * Jeden rozgłos testowanego eksperta na klienta. Byt jest wspólny, ponieważ
 * nadawanie i odbiór powstają w dwóch miejscach powłoki, których nikt nie
 * składa razem.
 */
export const testowanyAgent: RozglosTestowanego = utworzRozglosTestowanego();

/**
 * Składa zdanie dla okna kreatora o tym, czym jest czat tego modułu i czego po
 * nim nie oczekiwać. Treść bierze z pola pamięci sesyjnej w profilu modułu,
 * więc nadanie modułowi pamięci zmienia zdanie bez zmiany w tym pliku.
 */
export function zdanieOCzacieTestowym(): string {
  const profil = profilModulu(KOD_MODULU);
  if (profil.pamiecSesyjna) {
    return `Czat modułu ${profil.nazwa} prowadzi pamięć sesyjną — wątek przeżywa zamknięcie okna.`;
  }
  return (
    `Czat modułu ${profil.nazwa} jest środowiskiem TESTOWANIA wybranego eksperta i nie ma ` +
    'pamięci sesyjnej: rozmowa znika przy zamknięciu okna oraz przy zmianie testowanego ' +
    `agenta. ${ZDANIE_O_ZAPISIE_RDZENIA}`
  );
}
