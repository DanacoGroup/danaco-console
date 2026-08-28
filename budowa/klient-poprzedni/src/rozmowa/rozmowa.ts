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

/**
 * Ustawienia rozmowy, których w żaden sposób nie da się wyprowadzić z samego okna
 * komunikacji z rdzeniem.
 */
export interface OpcjeRozmowy {
  /** Tożsamość mówiącego po stronie modelu — nazwa kanału albo agenta. */
  persona?: string;
  /** Tożsamość Operatora widoczna w historii. */
  personaOperatora?: string;
  /** Rola okna w pętli koordynator–wykonawca. */
  rolaOkna?: WindowRole;
  // Polityka pamięci rozmowy w chwili złożenia; pominięta znaczy rozmowę z pamięcią.
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
  // Stawia w wątku zdanie warstwy automatycznej — wpis, którego nikt wprost nie
  // wypowiedział.
  zglosKomunikat(tresc: string): void;
  /** Subskrypcja wpisu założonego albo zmienionego. */
  naWpis(sluchacz: (wpis: WpisRozmowy) => void): Odsubskrybuj;
  /** Subskrypcja stanu wysyłania. */
  naStan(sluchacz: (stan: StanRozmowy) => void): Odsubskrybuj;
  /** Wpisy historii w kolejności powstania. */
  wpisy(): WpisRozmowy[];
  /** Bieżący stan wysyłania. */
  stan(): StanRozmowy;
  // Zdejmuje całą rozmowę i podaje powód, bo wyczyszczenie bez słowa wygląda jak awaria.
  wyczysc(powod: string): void;
  // Przestawia politykę pamięci, gdy okno zmieniło moduł roboczy, czyszcząc rozmowę wprost.
  ustawUlotnosc(polityka: PolitykaUlotnosci): void;
  /** Subskrypcja wyczyszczenia rozmowy — widok zdejmuje wszystkie pozycje. */
  naWyczyszczenie(sluchacz: (powod: string) => void): Odsubskrybuj;
  // Subskrypcja wypowiedzi, która weszła do okna spoza tego połączenia sieciowego.
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

  // Otoczenie, którego nadanie i przerwanie potrzebują od okna, składane raz dla obu
  // czynności.
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

  // Wysyła wiadomość Operatora; przerwaniem biegnącej tury zajmuje się osobna warstwa
  // nadania.
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

  // Niepowodzenie komendy zamyka turę i nic ponadto; treść błędu idzie bez doklejonego kodu.
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

  // Wyczyszczenie rozmowy ulotnej: najpierw powód, potem stan, bo kolejność jest tu treścią.
  function wyczysc(powod: string): void {
    historia.wyczysc();
    ustawStan({ wysyla: false, idWiadomosci: '', fragmenty: 0, ciszaSekundy: 0 });
    wyczyszczenia.oglos(powod);
    if (powod.trim().length > 0) zglosKomunikat(powod);
    for (const zdanie of polityka.zapowiedz) zglosKomunikat(zdanie);
  }

  // Podpięcie kontekstu roboczego bieżącej polityki; jego zmiana kończy rozmowę jak
  // zamknięcie.
  function podepnijKontekst(): void {
    odsubskrybujKontekst();
    odsubskrybujKontekst = () => {};
    const kontekst = polityka.kontekst;
    if (!czyUlotna(polityka) || kontekst === null) return;
    odsubskrybujKontekst = kontekst.naZmiane((zdanie) => wyczysc(zdanie));
  }

  // Przestawienie polityki po zmianie modułu okna, czyszczące wątek modułu poprzedniego.
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

  // Przełożenie wiadomości kontraktu na wpis wątku mieszka w osobnym pliku wiadomości z
  // rdzenia.
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

  /**
   * Historia okna wraca z rdzenia, nie ze strumienia, bo wiadomości leżą w bazie po
   * restarcie.
   */
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
        // Odpowiedź spóźniona o zmianę modułu nie wchodzi do wątku, ale pominięcie zostaje
        // nazwane.
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

/**
 * Zdanie zastępcze wypisywane w wątku rozmowy, gdy rdzeń domknął turę, nie podając żadnej
 * przyczyny.
 */
const BEZ_PRZYCZYNY = 'rdzeń nie podał przyczyny';
