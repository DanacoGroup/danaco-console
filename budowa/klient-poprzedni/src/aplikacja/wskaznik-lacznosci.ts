import type { Transport } from '../polaczenie/gniazdo';
import type { StanPolaczenia } from '../polaczenie/stan-polaczenia';
import { stanRdzeniaZPowloki, type StanRdzenia } from '../powloka/most-rdzenia';

/** Wskaźnik łączności zamontowany w grupie akcji paska. */
export interface WskaznikLacznosci {
  /** Plakietka montowana na pasku. */
  element: HTMLElement;
  /** Odłącza subskrypcję transportu. */
  rozlacz(): void;
}

/** Nazwa stanu widoczna dla Operatora. */
const NAZWY: Readonly<Record<StanPolaczenia, string>> = {
  rozlaczony: 'Rozłączony',
  laczenie: 'Łączenie',
  polaczony: 'Połączony',
  ponawianie: 'Ponawianie',
};

/** Odmiana kropki stanu ze słownika (`komponenty/plakietka.css`); bez barw własnych. */
const KROPKI: Readonly<Record<StanPolaczenia, string>> = {
  rozlaczony: 'dn-kropka--blad',
  laczenie: 'dn-kropka--tetno',
  polaczony: 'dn-kropka--sukces',
  ponawianie: 'dn-kropka--ostrzezenie',
};

/**
 * Stany, w których łączność jest w toku — nośnikiem jest wtedy `.dn-spinner`.
 *
 * Dla łączenia z serwerem nośnikiem jest wskaźnik ładowania z etykietą obok.
 * Spinner zastępuje kropkę zamiast stawać przy niej: dwa ruchy naraz
 * w plakietce wielkości pigułki spierałyby się o uwagę, a znaczenie niesie
 * i tak etykieta, nie sam znak.
 */
const W_TOKU: ReadonlySet<StanPolaczenia> = new Set<StanPolaczenia>(['laczenie', 'ponawianie']);

/** Odstęp dobijania licznika kolejki po połączeniu. */
const KROK_ODSWIEZANIA_MS = 150;

/**
 * Stan łączności z rdzeniem pokazany na pasku.
 *
 * Jedna odpowiedzialność: przełożenie stanu transportu na plakietkę. Wskaźnik
 * niczego nie wyłącza i nie blokuje — stan jest informacją, a nie bramą; treść
 * wpisana przy rozłączeniu czeka w kolejce wychodzącej. Liczba ramek
 * oczekujących trafia do plakietki, bo przy zerze blokad przejrzystość jest
 * jedynym zabezpieczeniem.
 *
 * Transport ogłasza `polaczony` przed opróżnieniem kolejki, więc odczyt zrobiony
 * w chwili zmiany stanu zamarłby na wartości sprzed wysłania. Po połączeniu
 * licznik odświeża się cyklicznie, aż kolejka spadnie do zera — wtedy pętla
 * gaśnie. W stanach innych niż połączony nic z gniazda nie schodzi, więc pętla
 * nie jest potrzebna.
 *
 * Napis „Rozłączony" mówi, co widzi transport, i nic o powodzie. Powód zna
 * powłoka natywna: to ona stawia proces rdzenia, wie, czy nasłuch odpowiada,
 * i prowadzi dziennik uruchomienia. Jej zdanie (polecenie `stan_rdzenia`)
 * dopisuje się do podpowiedzi plakietki, gdy łączności nie ma; poza powłoką
 * natywną pytanie nie pada.
 */
export function utworzWskaznikLacznosci(transport: Transport): WskaznikLacznosci {
  const element = document.createElement('span');
  // Plakietka neutralna: stan dopowiadają kropka i etykieta obok, nigdy sama
  // barwa tła.
  element.className = 'dn-plakietka dn-lacznosc';

  const kropka = document.createElement('span');
  kropka.className = 'dn-kropka';

  const wskaznik = document.createElement('span');
  wskaznik.className = 'dn-spinner dn-lacznosc__wskaznik';
  wskaznik.setAttribute('aria-hidden', 'true');
  wskaznik.hidden = true;

  // Napis stanu nie nosi klasy: całe jego zachowanie — barwę, stopień pisma
  // i zakaz łamania wiersza — niesie już plakietka. Klasa bez ani jednej reguły
  // we własnym arkuszu byłaby hakiem żyjącym na aliasie zgodności.
  const opis = document.createElement('span');

  element.append(kropka, wskaznik, opis);

  // Zdanie powłoki natywnej o rdzeniu; puste, dopóki powłoka nie odpowie.
  let diagnoza = '';
  // Pytanie w locie — plakietka zmienia stan częściej niż powłoka odpowiada,
  // a jedno pytanie naraz wystarczy, żeby nie mnożyć wywołań IPC.
  let pytanie = false;
  // Plakietka zdjęta z paska nie przyjmuje spóźnionej odpowiedzi.
  let czynna = true;
  let biezacy: StanPolaczenia = transport.stan();
  // Uchwyt pętli dobijającej licznik kolejki; null, gdy pętla nie biegnie.
  let uchwytOdswiezania: ReturnType<typeof setInterval> | null = null;

  function opiszPodpowiedz(): void {
    const podstawa = `Łączność z rdzeniem: ${opis.textContent}`;
    element.title = diagnoza === '' ? podstawa : `${podstawa}\n${diagnoza}`;
  }

  /** Przerysowuje etykietę i podpowiedź na podstawie bieżącej liczby w kolejce. */
  function odswiezEtykiete(stan: StanPolaczenia): void {
    const oczekujace = transport.oczekujace();
    opis.textContent =
      oczekujace > 0 ? `${NAZWY[stan]} · w kolejce ${oczekujace}` : NAZWY[stan];
    opiszPodpowiedz();
  }

  function zatrzymajOdswiezanie(): void {
    if (uchwytOdswiezania !== null) {
      clearInterval(uchwytOdswiezania);
      uchwytOdswiezania = null;
    }
  }

  /**
   * Po połączeniu kolejka opróżnia się już po ogłoszeniu stanu, więc licznik
   * odczytany w chwili zmiany jest nieaktualny. Dopóki zostają ramki, dobijamy
   * odczyt cyklicznie; gdy kolejka spadnie do zera, pętla się gasi. W stanach
   * innych niż połączony nic z otwartego gniazda nie schodzi, więc pętla nie
   * jest potrzebna.
   */
  function zarzadzajOdswiezaniem(stan: StanPolaczenia): void {
    if (stan === 'polaczony' && transport.oczekujace() > 0) {
      if (uchwytOdswiezania === null) {
        uchwytOdswiezania = setInterval(() => {
          odswiezEtykiete(biezacy);
          if (transport.oczekujace() === 0) zatrzymajOdswiezanie();
        }, KROK_ODSWIEZANIA_MS);
      }
      return;
    }
    zatrzymajOdswiezanie();
  }

  /** Zdanie dla Operatora złożone z tego, co wie wyłącznie powłoka. */
  function zdanieORdzeniu(stan: StanRdzenia): string {
    const nasluch = stan.pracuje
      ? `Rdzeń odpowiada pod ${stan.adres}.`
      : `Rdzeń nie odpowiada pod ${stan.adres}.`;
    return `${nasluch} ${stan.opis} Dziennik: ${stan.dziennik}`;
  }

  function dopytajPowloke(): void {
    if (pytanie) return;
    pytanie = true;
    void stanRdzeniaZPowloki().then((stan) => {
      pytanie = false;
      // Odpowiedź spóźniona o odzyskanie łączności jest już nieprawdziwa —
      // wtedy milczy się zamiast tłumaczyć ciszę, której nie ma.
      if (!czynna || stan === null || biezacy === 'polaczony') return;
      diagnoza = zdanieORdzeniu(stan);
      element.dataset.rdzen = stan.pracuje ? 'nasluchuje' : 'milczy';
      opiszPodpowiedz();
    });
  }

  function pokaz(stan: StanPolaczenia): void {
    biezacy = stan;
    const wToku = W_TOKU.has(stan);
    kropka.className = `dn-kropka ${KROPKI[stan]}`;
    kropka.hidden = wToku;
    wskaznik.hidden = !wToku;
    element.dataset.stan = stan;

    // Łączność odzyskana unieważnia poprzednią diagnozę: powód ciszy przestał
    // istnieć, więc nie zostaje w podpowiedzi jako nieaktualna wymówka.
    if (stan === 'polaczony') {
      diagnoza = '';
      delete element.dataset.rdzen;
    }
    odswiezEtykiete(stan);
    zarzadzajOdswiezaniem(stan);
    if (stan !== 'polaczony') dopytajPowloke();
  }

  // Transport ogłasza subskrybentowi stan bieżący, więc plakietka nie czeka
  // na pierwszą zmianę; wywołanie poniżej domyka przypadek transportu, który
  // jeszcze nie ruszył.
  const odsubskrybuj = transport.naStan(pokaz);
  pokaz(transport.stan());

  return {
    element,
    rozlacz: () => {
      // Plakietka zdjęta z paska: ani spóźniona odpowiedź powłoki, ani pętla
      // licznika nie mają już czego opisywać.
      czynna = false;
      zatrzymajOdswiezanie();
      odsubskrybuj();
    },
  };
}
