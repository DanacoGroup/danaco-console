import { profilModulu } from '../../okno-komunikacji/rejestr-profilow';
import { utworzMagistrale, type Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import { ZDANIE_O_ZAPISIE_RDZENIA, type KontekstRoboczy } from '../../rozmowa/ulotnosc';

/** Kod modułu w katalogu rdzenia — ten sam, którym moduł opisuje się powłoce. */
const KOD_MODULU = 'agents';

/**
 * Testowany ekspert — kontekst roboczy czatu testowego modułu Agents.
 *
 * Czat tego modułu jest środowiskiem testowania wybranego eksperta i nie ma
 * pamięci sesyjnej: rozmowa znika przy zamknięciu okna oraz przy zmianie
 * testowanego agenta.
 *
 * Rozgłos jest jednym bytem po stronie klienta, a nie polem modułu.
 * `aplikacja/przestrzen-modulu.ts` stawia widok modułu obok sceny sesji, a nie
 * zamiast niej — okno rozmowy jest oknem wiodącym każdego modułu i zostaje na
 * planszy. Czat modułu Agents to więc to samo okno komunikacji ze sceny,
 * przestawione na moduł `agents`, a nie okno wewnątrz modułu. Wybór eksperta
 * żyje po drugiej stronie planszy, w `stan-agentow.ts`.
 *
 * Plik niczego nie czyści. Ogłasza wyłącznie, że testowany ekspert się zmienił,
 * i podaje zdanie o tej zmianie. Co z tym zrobić, rozstrzyga polityka ulotności
 * rozmowy (`rozmowa/ulotnosc.ts`) — to ona zna profil modułu i pole
 * `pamiecSesyjna`.
 */
export interface RozglosTestowanego extends KontekstRoboczy {
  /**
   * Ustawia testowanego eksperta. Wywołuje wyłącznie `stan-agentow` —
   * jedyny właściciel wyboru w module.
   *
   * Wywołanie tą samą parą nie jest zmianą i nie ogłasza niczego: `oglos`
   * stanu agentów biegnie przy każdym odświeżeniu wykazu, przy każdej fazie
   * odczytu i przy każdym wchłonięciu zmiany z rdzenia. Gdyby każde z nich
   * liczyło się jako zmiana testowanego eksperta, rozmowa Operatora znikałaby
   * po zapisaniu instrukcji agenta, którego właśnie testuje.
   */
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

/** Zdanie o zmianie, gdy testowany ekspert znika (usunięty albo odznaczony). */
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
        // Sama nazwa mogła się zmienić — Operator przemianował testowanego
        // eksperta. To nie jest zmiana kontekstu i rozmowy nie kończy.
        nazwa = nazwaEksperta;
        return;
      }
      const poprzedni = nazwa === '' ? klucz : nazwa;
      const bylWybor = klucz !== '';
      klucz = idEksperta;
      nazwa = nazwaEksperta;
      // Pierwszy wybór nie jest zmianą. Moduł wchodzi z pustym wyborem i sam
      // wskazuje pierwszego eksperta z biblioteki po odczycie z rdzenia;
      // ogłoszenie tego jako „rozmowa wyczyszczona" kasowałoby zapowiedź
      // polityki, którą okno wypisało chwilę wcześniej, i meldowało utratę
      // wątku, którego jeszcze nie było.
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
 * Jeden rozgłos na klienta.
 *
 * Byt wspólny, bo jego dwie strony powstają w dwóch miejscach powłoki, których
 * nikt nie składa razem: `stan-agentow` przy budowie widoku modułu
 * (`przestrzen-modulu`), a odbiór — przy wiązaniu gniazda sceny z oknem rdzenia
 * (`aplikacja/wiazanie-gniazda`). Przekazanie „z rąk do rąk" wymagałoby
 * przeciągnięcia uchwytu przez rejestr modułów i przez scenę sesji, czyli
 * przez dwie warstwy, których ta rzecz nie dotyczy.
 */
export const testowanyAgent: RozglosTestowanego = utworzRozglosTestowanego();

/**
 * Zdanie dla okna kreatora: czym jest czat tego modułu i czego po nim nie
 * oczekiwać.
 *
 * Stoi w module, a nie tylko w oknie rozmowy, bo Operator przełączający eksperta
 * w bibliotece ma przeczytać, co się stanie z rozmową, zanim się to stanie.
 *
 * Treść czyta pole `pamiecSesyjna` z profilu modułu, więc nadanie modułowi
 * pamięci sesyjnej zmienia tę notę bez zmiany w tym pliku.
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
