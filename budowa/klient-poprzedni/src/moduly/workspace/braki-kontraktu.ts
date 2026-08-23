import { Command, KOMENDY, PROTOCOL_VERSION } from '../../../../shared/contract';
import { przyciskBezKomendy } from '../../modele/kontrolki-formularza';
import type { Kanal } from '../../protokol/kanal';
import { czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { tozsamoscKlienta } from '../../protokol/tozsamosc-klienta';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Powody nieczynnych kontrolek modułu Workspace — składane z wykazu komend
 * uruchomionego rdzenia i z wykazu komend kontraktu.
 *
 * Powód bierze się z powitania: `connection.hello` oddaje w polu `commands`
 * wykaz komend zarejestrowanych przez rdzeń, a nie wykaz z kontraktu. Moduł
 * pyta o niego raz przy montażu i z odpowiedzi układa zdanie każdej nieczynnej
 * kontrolki, więc rozróżnia brak po stronie rdzenia od braku w kontrakcie.
 * Gdy rdzeń komendę zarejestruje, zdanie zmienia się samo (wzór:
 * `moduly/katalog-okien.ts`).
 *
 * Stanów jest więcej niż dwa: dopóki rdzeń nie odpowiedział, kontrolka nie
 * orzeka o braku, tylko mówi, że pytanie jest w drodze; odmowa powitania też
 * nie staje się orzeczeniem o braku.
 *
 * Kontrolka bez pokrycia nie znika i nie udaje, że działa — zostaje widoczna,
 * klikalna i niesie powód wprost (`przyciskBezKomendy`).
 */

/** Komendy kontraktu jako zbiór — po nim poznajemy, po czyjej stronie jest brak. */
const W_KONTRAKCIE: ReadonlySet<string> = new Set<string>(KOMENDY);

/** Wykaz komend rdzenia wraz z kontrolkami, które z niego biorą swoje zdanie. */
export interface WykazBrakow {
  /** Pyta rdzeń o wykaz komend i przerysowuje powody wszystkich kontrolek. */
  odczytaj(): Promise<void>;
  /**
   * Kontrolka bez pokrycia; powód dopisuje się sam po odpowiedzi rdzenia.
   *
   * @param etykieta napis na przycisku.
   * @param czynnosc czego kontrolka miała dokonać — wchodzi do zdania powodu.
   * @param komendy komendy, które by tego dokonały; brak nazw znaczy, że okno
   *   nie zna żadnej.
   */
  przyciskBraku(
    etykieta: string,
    czynnosc: string,
    ...komendy: readonly string[]
  ): HTMLButtonElement;
}

export function utworzWykazBrakow(kanal: Kanal): WykazBrakow {
  /** Wykaz z rdzenia; `null`, dopóki powitanie nie wróciło. */
  let komendyRdzenia: ReadonlySet<string> | null = null;
  /** Treść odmowy powitania; pusta, dopóki nic się nie zepsuło. */
  let odmowa = '';
  /** Kontrolki, które trzeba przerysować po odpowiedzi rdzenia. */
  const zalezne: Array<() => void> = [];

  /**
   * Zdanie o pokryciu jednej komendy. Rozróżnia cztery stany: odmowę powitania,
   * pytanie w drodze, komendę zarejestrowaną przez rdzeń oraz brak — osobno po
   * stronie rdzenia i osobno w kontrakcie.
   */
  function zdanieOKomendzie(nazwa: string): string {
    if (odmowa !== '') {
      return `Rdzeń nie oddał wykazu komend (${odmowa}) — o pokryciu komendy ${nazwa} nic nie wiadomo.`;
    }
    if (komendyRdzenia === null) {
      return `Sprawdzanie u rdzenia, czy obsługuje komendę ${nazwa}…`;
    }
    if (komendyRdzenia.has(nazwa)) {
      return `Rdzeń rejestruje komendę ${nazwa}, ale to okno jej nie wywołuje — czynność niezbudowana.`;
    }
    return W_KONTRAKCIE.has(nazwa)
      ? `Rdzeń tej wersji nie rejestruje komendy ${nazwa}, choć kontrakt ją niesie — ` +
          'brak jest po stronie rdzenia (zmierzone powitaniem connection.hello).'
      : `Ani kontrakt, ani rdzeń nie znają komendy ${nazwa} — czynności nie ma czym wykonać.`;
  }

  /** Powód nieczynności kontrolki opartej na jednej komendzie albo na kilku naraz. */
  function powod(czynnosc: string, komendy: readonly string[]): string {
    if (komendy.length === 0) {
      return `${czynnosc}: okno nie zna komendy, która by tego dokonała.`;
    }
    return `${czynnosc}. ${komendy.map(zdanieOKomendzie).join(' ')}`;
  }

  return {
    async odczytaj() {
      const klient = tozsamoscKlienta();
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.ConnectionHello, {
          clientId: klient.id,
          clientVersion: klient.wersja,
          protocolVersion: PROTOCOL_VERSION,
        }),
        Command.ConnectionHello,
        (tresc) => czyTablica(tresc.commands),
      );
      if (!wynik.udany || wynik.wynik === undefined) {
        odmowa = `kod ${wynik.blad?.code ?? 'brak'}`;
        komendyRdzenia = null;
      } else {
        odmowa = '';
        komendyRdzenia = new Set(wynik.wynik.commands ?? []);
      }
      for (const przerysuj of zalezne) przerysuj();
    },

    przyciskBraku(etykieta, czynnosc, ...komendy) {
      let kontrolka = przyciskBezKomendy(etykieta, powod(czynnosc, komendy));
      // Powód idzie trzema drogami naraz — `title`, `aria-description` i dymek
      // po naciśnięciu — a dymek `przyciskBezKomendy` domyka powód w chwili
      // budowy. Podmiana samych atrybutów zostawiłaby dymek ze zdaniem sprzed
      // odpowiedzi rdzenia, więc kontrolka powstaje na nowo i wchodzi na miejsce
      // poprzedniej; jednego źródła kontrolki (`kontrolki-formularza`) to nie rusza.
      zalezne.push(() => {
        const nowa = przyciskBezKomendy(etykieta, powod(czynnosc, komendy));
        if (kontrolka.parentNode === null) {
          // Kontrolka poza drzewem: trzyma ją jeszcze wywołujący, więc
          // podmiana węzła nic by nie dała — zostają dwie drogi z trzech.
          kontrolka.title = nowa.title;
          kontrolka.setAttribute('aria-description', nowa.title);
          return;
        }
        kontrolka.replaceWith(nowa);
        kontrolka = nowa;
      });
      return kontrolka;
    },
  };
}
