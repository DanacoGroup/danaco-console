import {
  Command,
  EventType,
  type DeviceChangedEvent,
  type DeviceListResponse,
  type DeviceRevokeResponse,
} from '../../../shared/contract';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../protokol/kanal';
import { czyLogiczna, czyTablica, sprawdzKsztalt } from '../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../protokol/wywolanie';

/**
 * Czynności na urządzeniach powiązanych z kontem — rodzina `device.*`.
 *
 * Rodzina jest dwuczłonowa i tyle samo niesie to źródło: `device.list` oddaje
 * wykaz, `device.revoke` unieważnia token jednego urządzenia. Przemianowania
 * ani odłączenia urządzenia od konta kontrakt nie zna, więc nie ma tu czego
 * dołożyć.
 *
 * Rozdział względem `zrodlo-bramki.ts` przebiega po przedmiocie, nie po oknie:
 * tamto źródło prowadzi metody wejścia i sesję bramki (`auth.*`), to prowadzi
 * urządzenia. Obie rodziny spotykają się w polu `deviceId`, ale odpowiadają za
 * co innego — zdjęcie metody PIN zostawia token urządzenia nietknięty,
 * a unieważnienie tokenu nie zdejmuje metod.
 *
 * Wykaz po zmianie przychodzi zdarzeniem, nie z odpowiedzi: `device.revoke`
 * oddaje wyłącznie `revoked`, a pełny wykaz rdzeń rozgłasza `device.changed` do
 * wszystkich połączonych urządzeń. Dzięki temu unieważnienie wykonane na jednej
 * maszynie widać natychmiast na pozostałych ekranach — i dlatego sekcja nie
 * odpytuje rdzenia po każdej czynności.
 */
export interface ZrodloUrzadzen {
  /**
   * `device.list` — wykaz urządzeń powiązanych z kontem właściciela.
   *
   * Jedyna droga do wykazu przy wejściu do sekcji: zanim zajdzie jakakolwiek
   * zmiana, nie ma czego rozgłosić zdarzeniem.
   */
  wykaz(): Promise<Wynik<DeviceListResponse>>;
  /**
   * `device.revoke` — unieważnienie tokenu dostępu wskazanego urządzenia.
   *
   * Rdzeń nie broni unieważnienia własnego tokenu; orzeczenie o skutku należy
   * do niego, a ostrzeżenie do ekranu. Wynik niesie samo `revoked`, bo wykaz po
   * zmianie idzie zdarzeniem.
   */
  uniewaznij(idUrzadzenia: string): Promise<Wynik<DeviceRevokeResponse>>;
  /**
   * Subskrypcja `device.changed` — jedyna droga, którą wykaz nadąża za
   * czynnością wykonaną poza tym oknem.
   *
   * Zdarzenie niesie pełny wykaz po zmianie oraz — opcjonalnie — urządzenie,
   * którego zmiana dotyczy. Pole `deviceId` bywa puste (tak przychodzi zmiana
   * hasła, która unieważnia tokeny hurtem) i sekcja nie zgaduje za rdzeń, czego
   * dotyczyła.
   */
  naZmianeUrzadzen(sluchacz: (zmiana: DeviceChangedEvent) => void): Odsubskrybuj;
}

export function utworzZrodloUrzadzen(kanal: Kanal): ZrodloUrzadzen {
  return {
    async wykaz() {
      // Kształt sprawdzamy, bo rzutowanie kanału jest obietnicą kompilatora,
      // nie rdzenia: brak pola `devices` trafiłby do widoku jako `undefined`
      // w miejscu, w którym kontrakt obiecuje tablicę.
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DeviceList, {}),
        Command.DeviceList,
        (tresc) => czyTablica(tresc.devices),
      );
    },

    async uniewaznij(idUrzadzenia) {
      // Brak `revoked` bez sprawdzianu wyglądałby jak czynność udana — a to
      // jedyne pole, którym rdzeń orzeka, czy token faktycznie zszedł.
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DeviceRevoke, { deviceId: idUrzadzenia }),
        Command.DeviceRevoke,
        (tresc) => czyLogiczna(tresc.revoked),
      );
    },

    naZmianeUrzadzen(sluchacz) {
      return kanal.naZdarzenie(EventType.DeviceChanged, (tresc) => sluchacz(tresc));
    },
  };
}
