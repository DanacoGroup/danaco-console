import { utworzKafelKomponentu } from './kafel-komponentu';
import { POZYCJE_KOMPONENTOW, type PozycjaKomponentu } from './pozycje-komponentow';
import { utworzStrefeZwijana } from './strefa-zwijana';
import { utworzSygnalWyboru, type SluchaczWyboru } from './sygnal-wyboru';

/**
 * Strefa druga — kafle komponentów własnych.
 *
 * Jedna odpowiedzialność: siatka kafli komponentów.
 *
 * Siatka niesie dwa rodzaje kafli. Pierwsze pochodzą z modułów, które rdzeń
 * oznaczył jako nastawiane na Stronie głównej (`Module.configuredOnHome`); do
 * pierwszej odpowiedzi stoi tam czwórka zastana z kontraktu. Za nimi stoją
 * kafle personalizowane — po jednym na komponent zbudowany i nazwany przez
 * Operatora. Żadna z dwóch liczb nie jest z góry znana, więc siatka przyjmuje
 * oba wykazy osobno i przerysowuje się, gdy rdzeń odpowie.
 *
 * Waga wizualna niższa niż strefy pierwszej: mniejsza powierzchnia, mniejszy
 * promień, mniejsza ikona, etykieta krojem bazowym półgrubym. Bez wstęgi
 * i bez cienia sygnału w spoczynku.
 */

/** Nazwa strefy z opracowania (rozdz. 3.1) i z makiety — bez parafrazy. */
const ETYKIETA = 'Strefa 2 · Kafle komponentów własnych';
const WYJASNIENIE =
  'Kliknięcie kafla nie otwiera przestrzeni roboczej — otwiera okno konfiguracji, w którym powstaje nazwany wytwór.';


export interface StrefaKomponentow {
  /** Element `<details>` — strefa przywoływana, nie rysowana z urzędu. */
  element: HTMLDetailsElement;
  /** Miejsce na formularz zakładania — wpina je warstwa, która ma kanał. */
  przybornik: HTMLElement;
  naWybor(sluchacz: SluchaczWyboru<PozycjaKomponentu>): void;
  /**
   * Ustawia kafle modułów nastawianych na Stronie głównej — wykazem z rdzenia.
   *
   * Wykaz pusty nie zdejmuje kafli: zostaje czwórka zastana. Strefa bez ani
   * jednego wejścia byłaby regresem widocznym na ekranie, a pusta odpowiedź
   * bywa też odpowiedzią bazy bez wykonanych migracji.
   */
  ustawRodzaje(pozycje: readonly PozycjaKomponentu[]): void;
  /**
   * Ustawia kafle komponentów zbudowanych przez Operatora. Wykaz pusty jest
   * poprawnym stanem: Operator, który nic jeszcze nie zbudował, widzi same
   * kafle modułów.
   */
  ustawPersonalizowane(pozycje: readonly PozycjaKomponentu[]): void;
}

export function utworzStrefeKomponentow(): StrefaKomponentow {
  const sygnal = utworzSygnalWyboru<PozycjaKomponentu>();

  // Strefa zwijalna, ale rozwinięta domyślnie: jest jedną z dróg wejścia
  // w pracę, więc zwinięcie z urzędu ukrywałoby ją przed Operatorem.
  // Zwinięcie wykonane ręcznie zostaje zapamiętane pod kluczem strefy.
  const strefa = utworzStrefeZwijana({
    etykieta: ETYKIETA,
    wyjasnienie: WYJASNIENIE,
    klucz: 'strona.komponenty',
    domyslnieRozwiniete: true,
  });
  const element = strefa.element;
  element.classList.add('dn-strona__strefa--komponenty');
  element.setAttribute('aria-label', ETYKIETA);

  // Jedna siatka na oba rodzaje kafli: personalizowane stoją obok rodzajów,
  // w tym samym rzędzie i tej samej wielkości. Druga siatka z własnym
  // podpisem wydłużałaby stronę o dwa pasy i spychała sesje w tle oraz
  // archiwum poniżej pierwszego ekranu. Rodzaj kafla rozpoznaje się po jego
  // treści (wezwanie „Zbuduj…" wobec nazwy własnej), nie po osobnej sekcji.
  const siatka = document.createElement('div');
  siatka.className = 'dn-strona__siatka dn-strona__siatka--komponenty';

  const przybornik = document.createElement('div');
  przybornik.className = 'dn-strona__przybornik';

  function kafel(pozycja: PozycjaKomponentu): HTMLElement {
    return utworzKafelKomponentu(pozycja, (wybrany) => sygnal.nadaj(wybrany)).element;
  }

  // Czwórka zastana stoi od pierwszej klatki, żeby strefa nie była pusta, zanim
  // rdzeń odpowie na `module.list`. Po odpowiedzi rozstrzyga rdzeń.
  let rodzaje: readonly PozycjaKomponentu[] = POZYCJE_KOMPONENTOW;
  let personalizowane: readonly PozycjaKomponentu[] = [];

  function przerysuj(): void {
    siatka.replaceChildren(...rodzaje.map(kafel), ...personalizowane.map(kafel));
    strefa.ustawDopisek(String(rodzaje.length + personalizowane.length));
  }

  przerysuj();

  // Przybornik stoi pod siatką, nie w wierszu nagłówka: utworzenie komponentu
  // jest czynnością wykonywaną po obejrzeniu tego, co już stoi, a nad kaflami
  // konkurowałoby o uwagę z etykietą strefy.
  strefa.cialo.append(siatka, przybornik);

  return {
    element,
    przybornik,
    naWybor: sygnal.sluchaj,

    ustawRodzaje(pozycje) {
      if (pozycje.length === 0) return;
      rodzaje = pozycje;
      przerysuj();
    },

    ustawPersonalizowane(pozycje) {
      personalizowane = pozycje;
      przerysuj();
    },
  };
}
