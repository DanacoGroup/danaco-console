import { EventType, MessageRole, type ErrorInfo, type WindowRole } from '../../../shared/contract';
import { utworzMagistrale, type Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal } from '../protokol/kanal';
import { wczytajHistorie } from './historia-z-rdzenia';
import { kluczMiejscowy, utworzHistorieTur } from './historia-tur';
import { rozpoznajNadawce, RodzajNadawcy } from './nadawca';
import type { StanRozmowy } from './stan-rozmowy';
import { czyUlotna, politykaTrwala, type PolitykaUlotnosci } from './ulotnosc';
import { wpisAutomatyzacji, wpisOperatora, type WpisRozmowy } from './wpis-rozmowy';
import {
  nadajZPrzerwaniem,
  przerwijTure,
  type OtoczenieNadania,
} from './nadanie-i-przerwanie';
import { przyjmijWiadomosc } from './wiadomosc-z-rdzenia';
import { utworzZegarCiszy } from './zegar-ciszy';
import { dolaczFragment } from './zlozenie-tury';

/** Ustawienia rozmowy, których nie da się wyprowadzić z samego okna. */
export interface OpcjeRozmowy {
  /** Tożsamość mówiącego po stronie modelu — nazwa kanału albo agenta. */
  persona?: string;
  /** Tożsamość Operatora widoczna w historii. */
  personaOperatora?: string;
  /** Rola okna w pętli koordynator–wykonawca. */
  rolaOkna?: WindowRole;
  /**
   * Polityka pamięci rozmowy w chwili złożenia (`ulotnosc.ts`).
   *
   * Pominięta znaczy „rozmowa z pamięcią" — stan dotychczasowy. Podana
   * z `pamiecSesyjna: false` wstrzymuje odtworzenie historii z rdzenia,
   * bo w module bez pamięci sesyjnej `message.list` przywróciłby dokładnie ten
   * wątek, który miał zniknąć razem z oknem.
   */
  ulotnosc?: PolitykaUlotnosci;
}

/**
 * Rozmowa jednego okna komunikacji.
 *
 * Warstwa nie zna widoku: wystawia wpisy i stan, a rysowaniem zajmuje się
 * `widok-rozmowy`. Nie zna też transportu — nazwy komend i zdarzeń bierze
 * wyłącznie z `shared/contract`.
 */
export interface Rozmowa {
  /** Wysyła wypowiedź Operatora komendą `message.send`. */
  wyslij(tekst: string): void;
  /** Przerywa bieżącą turę komendą `message.stop`; czynne zawsze. */
  przerwij(): void;
  /**
   * Stawia w wątku zdanie warstwy automatycznej — wpis, którego nikt nie mówił.
   *
   * Tą samą drogą idą powody wyczyszczenia rozmowy, odmowa `message.list`
   * i zatrzymanie tury bez tury. Wejście jest wystawione na zewnątrz, bo
   * dołożenie modelowi narzędzia spoza okna też musi być widoczne, a jedynym
   * miejscem, w którym Operator patrzy, jest wątek rozmowy.
   *
   * To nie jest druga droga wypowiedzi: zdanie nie idzie do rdzenia, nie
   * otwiera tury i nie ma nadawcy — jest wpisem automatyzacji, tak samo jak
   * zdania o ulotności.
   */
  zglosKomunikat(tresc: string): void;
  /** Subskrypcja wpisu założonego albo zmienionego. */
  naWpis(sluchacz: (wpis: WpisRozmowy) => void): Odsubskrybuj;
  /** Subskrypcja stanu wysyłania. */
  naStan(sluchacz: (stan: StanRozmowy) => void): Odsubskrybuj;
  /** Wpisy historii w kolejności powstania. */
  wpisy(): WpisRozmowy[];
  /** Bieżący stan wysyłania. */
  stan(): StanRozmowy;
  /**
   * Zdejmuje całą rozmowę i podaje powód.
   *
   * Powód jest obowiązkowy: wyczyszczenie bez słowa jest dla Operatora
   * nieodróżnialne od awarii, w której aplikacja zgubiła jego pracę.
   */
  wyczysc(powod: string): void;
  /**
   * Przestawia politykę pamięci — okno zmieniło moduł (`workspace.enter`).
   *
   * Wejście w moduł bez pamięci sesyjnej czyści rozmowę i mówi o tym wprost:
   * wątek modułu poprzedniego nie ma prawa zostać na oczach Operatora w oknie,
   * które właśnie ogłosiło, że pamięci nie prowadzi.
   */
  ustawUlotnosc(polityka: PolitykaUlotnosci): void;
  /** Subskrypcja wyczyszczenia rozmowy — widok zdejmuje wszystkie pozycje. */
  naWyczyszczenie(sluchacz: (powod: string) => void): Odsubskrybuj;
  /**
   * Subskrypcja wypowiedzi, która weszła do okna spoza tego połączenia.
   *
   * Nośnikiem obrazu „moja klawiatura pisze za mnie" jest widok, nie ta
   * warstwa: rozmowa mówi wyłącznie, co przyszło i skąd, a widok rozstrzyga,
   * czy pokazać to w polu wypowiedzi i jak długo. Sygnał idzie osobno od wpisu,
   * bo „ktoś napisał to za mnie" jest zdarzeniem innego rodzaju niż „jest nowa
   * pozycja w wątku" — i tylko pierwsze ma prawo ruszyć pole Operatora.
   */
  naWypowiedzZZewnatrz(sluchacz: (tresc: string) => void): Odsubskrybuj;
  /** Odłącza subskrypcje kanału. */
  rozlacz(): void;
}

export function utworzRozmowe(kanal: Kanal, idOkna: string, opcje: OpcjeRozmowy = {}): Rozmowa {
  const wpisy = utworzMagistrale<WpisRozmowy>();
  const stany = utworzMagistrale<StanRozmowy>();
  const wyczyszczenia = utworzMagistrale<string>();
  const wypowiedziZZewnatrz = utworzMagistrale<string>();
  const historia = utworzHistorieTur();
  const odsubskrybowania: Odsubskrybuj[] = [];

  let polityka: PolitykaUlotnosci = opcje.ulotnosc ?? politykaTrwala('');
  /** Odwołanie subskrypcji kontekstu roboczego bieżącej polityki. */
  let odsubskrybujKontekst: Odsubskrybuj = () => {};

  const persona = opcje.persona ?? 'Model';
  const personaOperatora = opcje.personaOperatora ?? 'Operator';
  let rolaOkna: WindowRole | null = opcje.rolaOkna ?? null;
  let stan: StanRozmowy = { wysyla: false, idWiadomosci: '', fragmenty: 0, ciszaSekundy: 0 };

  /** Chwila ostatniego znaku życia tury: wysłania komendy albo fragmentu. */
  let znakZycia = 0;
  const zegarCiszy = utworzZegarCiszy({
    stan: () => stan,
    znakZycia: () => znakZycia,
    ustawStan,
  });

  /** Ogłasza wpis i stan po każdej zmianie — widok odświeża jedną pozycję. */
  function oglos(wpis: WpisRozmowy): void {
    wpisy.oglos(wpis);
  }

  function ustawStan(zmiana: Partial<StanRozmowy>): void {
    stan = { ...stan, ...zmiana };
    stany.oglos(stan);
    if (stan.wysyla) zegarCiszy.uruchom();
    else zegarCiszy.zatrzymaj();
  }

  /** Nadawca wypowiedzi modelu przy bieżącej roli okna. */
  function nadawcaModelu(): RodzajNadawcy {
    return rozpoznajNadawce(MessageRole.Assistant, rolaOkna);
  }

  /**
   * Otoczenie, którego nadanie i przerwanie potrzebują od okna
   * (`nadanie-i-przerwanie.ts`). Składa się raz — obie czynności widzą ten sam
   * wątek, tę samą historię i ten sam wskaźnik wysyłania.
   */
  function otoczenieNadania(): OtoczenieNadania {
    return {
      kanal,
      idOkna,
      biezaca: () => historia.biezaca(),
      oglos,
      zglosBlad,
      zglosKomunikat,
      zakonczWysylanie: () => ustawStan({ wysyla: false, ciszaSekundy: 0 }),
    };
  }

  /**
   * Wysyła wiadomość Operatora.
   *
   * Wpisy powstają od razu, przed jakąkolwiek komendą — interfejs nie ma
   * blokad. Przerwaniem tury biegnącej zajmuje się warstwa nadania: rdzeń nie
   * przerywa jej przy okazji wysłania, więc przerwanie jedzie osobną komendą
   * (`nadanie-i-przerwanie.ts`).
   */
  function wyslij(tekst: string): void {
    const tresc = tekst.trim();
    if (tresc.length === 0) return;

    oglos(historia.dodaj(wpisOperatora(kluczMiejscowy('operator'), tresc, personaOperatora)));
    const odpowiedz = historia.oczekujaca(nadawcaModelu(), persona);
    oglos(odpowiedz);
    znakZycia = Date.now();
    ustawStan({ wysyla: true, idWiadomosci: '', fragmenty: 0, ciszaSekundy: 0 });

    nadajZPrzerwaniem(otoczenieNadania(), tresc, odpowiedz);
  }

  /** Przerwanie tury przyciskiem Operatora (`nadanie-i-przerwanie.ts`). */
  function przerwij(): void {
    przerwijTure(otoczenieNadania());
  }

  /**
   * Niepowodzenie komendy zamyka turę i nic ponadto.
   *
   * Treść błędu idzie tu sama, bez doklejonego kodu: kod jedzie osobnym polem
   * i dokleja go widok wpisu. Sklejenie obu w tym miejscu pokazałoby kod dwa
   * razy — „…nie istnieje: okn_… (not_found) (not_found)".
   */
  function zglosBlad(wpis: WpisRozmowy, blad: ErrorInfo | undefined): void {
    wpis.bledy.push({
      kod: blad?.code ?? '',
      tresc: blad?.message ?? BEZ_PRZYCZYNY,
      ponawialny: blad?.retryable ?? true,
    });
    wpis.stan = 'bledny';
    wpis.domkniety = true;
    oglos(wpis);
    ustawStan({ wysyla: false, ciszaSekundy: 0 });
  }

  /** Komunikat warstwy automatycznej w historii. */
  function zglosKomunikat(tresc: string): void {
    oglos(historia.dodaj(wpisAutomatyzacji(kluczMiejscowy('automat'), tresc)));
  }

  /**
   * Wyczyszczenie rozmowy ulotnej — najpierw powód, potem stan.
   *
   * Kolejność jest tu treścią. Widok najpierw zdejmuje wszystkie pozycje, więc
   * gdyby na tym poprzestać, Operator zostałby z pustą listą i napisem „Rozmowa
   * jeszcze się nie zaczęła" — nieodróżnialnym od okna dopiero co otwartego
   * i od okna, które zgubiło pracę. Dlatego zaraz po wyczyszczeniu wraca zdanie
   * o powodzie, a po nim zapowiedź polityki: co zniknęło, dlaczego i dokąd ta
   * ulotność sięga.
   */
  function wyczysc(powod: string): void {
    historia.wyczysc();
    ustawStan({ wysyla: false, idWiadomosci: '', fragmenty: 0, ciszaSekundy: 0 });
    wyczyszczenia.oglos(powod);
    if (powod.trim().length > 0) zglosKomunikat(powod);
    for (const zdanie of polityka.zapowiedz) zglosKomunikat(zdanie);
  }

  /**
   * Podpięcie kontekstu roboczego bieżącej polityki.
   *
   * Zmiana kontekstu — dla Agents zmiana testowanego eksperta — kończy rozmowę
   * dokładnie tak samo, jak zamknięcie okna. Zdanie przychodzi z kontekstu,
   * bo tylko on wie, co się zmieniło; warstwa rozmowy nie zna dziedziny.
   */
  function podepnijKontekst(): void {
    odsubskrybujKontekst();
    odsubskrybujKontekst = () => {};
    const kontekst = polityka.kontekst;
    if (!czyUlotna(polityka) || kontekst === null) return;
    odsubskrybujKontekst = kontekst.naZmiane((zdanie) => wyczysc(zdanie));
  }

  /**
   * Przestawienie polityki po zmianie modułu okna.
   *
   * Wejście w moduł ulotny czyści to, co zostało po module poprzednim. Wyjście
   * z modułu ulotnego czyści tak samo — wątek testowy nie ma prawa pojechać za
   * Operatorem do modułu, który prowadzi pracę na serio.
   */
  function ustawUlotnosc(nowa: PolitykaUlotnosci): void {
    const bylaUlotna = czyUlotna(polityka);
    const bedzieUlotna = czyUlotna(nowa);
    const tenSamModul = polityka.modul === nowa.modul;
    polityka = nowa;
    podepnijKontekst();
    if (tenSamModul || (!bylaUlotna && !bedzieUlotna)) return;
    wyczysc(
      bedzieUlotna
        ? 'Okno weszło w moduł bez pamięci sesyjnej — dotychczasowy wątek zdjęty z okna.'
        : 'Wątek czatu roboczego zdjęty — okno wyszło z modułu bez pamięci sesyjnej.',
    );
  }

  odsubskrybowania.push(
    kanal.naZdarzenie(EventType.StreamChunk, (fragment, koperta) => {
      if (fragment.windowId !== idOkna) return;
      const wpis = historia.dlaWiadomosci(fragment.messageId, nadawcaModelu(), persona);
      dolaczFragment(wpis, fragment, {
        numer: koperta.seq ?? wpis.ostatniNumer + 1,
        ostatni: koperta.done === true,
      });
      oglos(wpis);
      znakZycia = Date.now();
      ustawStan({
        wysyla: !wpis.domkniety,
        idWiadomosci: wpis.idWiadomosci,
        fragmenty: wpis.fragmenty,
        ciszaSekundy: 0,
      });
    }),
  );

  // Przełożenie wiadomości kontraktu na wpis wątku mieszka w osobnym pliku
  // (`wiadomosc-z-rdzenia.ts`) wraz z uzasadnieniem, dlaczego okno pokazuje
  // także wypowiedzi roli `user` i jak broni się przed ich podwojeniem.
  odsubskrybowania.push(
    kanal.naZdarzenie(EventType.MessageChanged, ({ message }) => {
      if (message.windowId !== idOkna) return;
      przyjmijWiadomosc(
        {
          historia,
          rolaOkna: () => rolaOkna,
          persona,
          oglos,
          ustawStan,
          zglosWypowiedzZZewnatrz: (tresc: string) => wypowiedziZZewnatrz.oglos(tresc),
        },
        message,
      );
    }),
  );

  odsubskrybowania.push(
    kanal.naZdarzenie(EventType.WindowChanged, ({ window }) => {
      if (window.id !== idOkna) return;
      rolaOkna = window.windowRole ?? rolaOkna;
    }),
  );

  // Historia okna wraca z rdzenia, nie ze strumienia: wiadomości leżą w bazie
  // i przeżywają restart rdzenia, więc odświeżenie okna albo powrót do sesji
  // mają odtworzyć wątek.
  //
  // Odmowa `message.list` nie może wyglądać jak pusta rozmowa. Bez drugiej
  // gałęzi Operator widziałby ten sam pusty stan „Rozmowa jeszcze się nie
  // zaczęła" niezależnie od tego, czy rozmowa naprawdę jest pusta, czy rdzeń
  // po prostu nie oddał wiadomości. Niepowodzenie idzie więc do historii jako
  // wpis automatyzacji — tą samą drogą, którą `zglosKomunikat` niesie
  // przerwanie bez tury — i na ekranie zostaje ślad zamiast ciszy.
  //
  // Rozmowa ulotna historii nie odtwarza. Rdzeń zapisuje wiadomości niezależnie
  // od modułu (`adapter_rozmowa.go` → `dziennik_rozmowy.go` → tabela
  // `wiadomosc`), więc gdyby okno modułu bez pamięci sesyjnej wołało
  // `message.list` tak jak każde inne, wątek znikający po zamknięciu okna
  // wracałby przy każdym otwarciu w komplecie. Okno mówi o tym wprost zdaniami
  // zapowiedzi, zamiast po cichu pokazywać pustkę.
  if (czyUlotna(polityka)) {
    for (const zdanie of polityka.zapowiedz) zglosKomunikat(zdanie);
  } else {
    wczytajHistorie(
      kanal,
      idOkna,
      rolaOkna,
      persona,
      personaOperatora,
      (odtworzone) => {
        // Odpowiedź spóźniona o zmianę modułu: `message.list` poszło, gdy okno
        // stało jeszcze w module z pamięcią, a zanim wróciło, `workspace.enter`
        // przestawił je na moduł bez pamięci sesyjnej. Wstawienie tych wpisów
        // przywróciłoby wątek dopiero co zdjęty na oczach Operatora, więc wpisy
        // nie wchodzą — ale pominięcie zostaje nazwane, bo milczące porzucenie
        // odpowiedzi rdzenia jest nieodróżnialne od jej zgubienia.
        if (czyUlotna(polityka)) {
          if (odtworzone.length > 0) {
            zglosKomunikat(
              `Historia z rdzenia (${odtworzone.length} wpisów) nie została wstawiona: ` +
                'zanim odpowiedź wróciła, okno weszło w moduł bez pamięci sesyjnej.',
            );
          }
          return;
        }
        for (const wpis of odtworzone) oglos(historia.dodaj(wpis));
      },
      (powod) => zglosKomunikat(powod),
    );
  }

  podepnijKontekst();

  return {
    wyslij,
    przerwij,
    zglosKomunikat,
    naWpis: (sluchacz) => wpisy.subskrybuj(sluchacz),
    naStan: (sluchacz) => stany.subskrybuj(sluchacz),
    naWyczyszczenie: (sluchacz) => wyczyszczenia.subskrybuj(sluchacz),
    naWypowiedzZZewnatrz: (sluchacz) => wypowiedziZZewnatrz.subskrybuj(sluchacz),
    wpisy: () => historia.wpisy(),
    stan: () => stan,
    wyczysc,
    ustawUlotnosc,
    rozlacz: () => {
      zegarCiszy.zatrzymaj();
      odsubskrybujKontekst();
      odsubskrybujKontekst = () => {};
      for (const odsubskrybuj of odsubskrybowania.splice(0)) odsubskrybuj();
    },
  };
}

/** Zdanie zastępcze, gdy rdzeń domknął turę, nie podając przyczyny. */
const BEZ_PRZYCZYNY = 'rdzeń nie podał przyczyny';
