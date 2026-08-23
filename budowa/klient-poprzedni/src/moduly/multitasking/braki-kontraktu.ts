import { Command, PROTOCOL_VERSION } from '../../../../shared/contract';
import { przyciskBezKomendy } from '../../modele/kontrolki-formularza';
import type { Kanal } from '../../protokol/kanal';
import { czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { tozsamoscKlienta } from '../../protokol/tozsamosc-klienta';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Komendy, których okna ról MultitaskingAI potrzebują, wraz z odpowiedzią
 * uruchomionego rdzenia na pytanie, czy je obsługuje.
 *
 * Zdanie o pokryciu bierze się z rdzenia, nie ze stałej: powitanie
 * `connection.hello` oddaje w polu `commands` wykaz komend zarejestrowanych po
 * montażu — obsługiwanych naprawdę, a nie tylko wypisanych w kontrakcie
 * (`handlers_connection.go`). Moduł pyta o niego raz przy montażu i układa
 * z odpowiedzi zdanie każdej nieczynnej kontrolki, więc gdy rdzeń domknie
 * kolejną komendę, zdanie zmienia się samo (wzór: `moduly/katalog-okien.ts`).
 *
 * Stany są trzy, nie dwa: dopóki rdzeń nie odpowiedział, kontrolka nie orzeka
 * o braku, a odmowa powitania też nie jest orzeczeniem braku.
 *
 * Kontrolka bez pokrycia nie znika i nie udaje, że działa — zostaje widoczna,
 * nieczynna i niesie powód wprost, żeby brak pozostał widoczny w oknie.
 */

/** Komenda wymagana przez okno roli wraz z jej przeznaczeniem. */
export interface KomendaRoli {
  /** Nazwa komendy w kontrakcie. */
  komenda: string;
  /** Co ta komenda miała umożliwić. */
  przeznaczenie: string;
}

export const WYMAGANE: readonly KomendaRoli[] = [
  { komenda: 'role.assign', przeznaczenie: 'nadanie roli oknu wraz z profilem wcielenia' },
  { komenda: 'role.update', przeznaczenie: 'zmiana wcielenia i trybu współpracy po stronie rdzenia' },
  { komenda: 'subagent.spawn', przeznaczenie: 'uruchomienie podagenta wykonawcy (do 15)' },
  { komenda: 'subagent.list', przeznaczenie: 'wykaz podagentów biegnących pod wykonawcą' },
  { komenda: 'subagent.result.collect', przeznaczenie: 'scalenie wyników podagentów' },
  { komenda: 'orchestration.dependency.set', przeznaczenie: 'zależności DAG między etapami planu' },
  { komenda: 'monitor.subscribe', przeznaczenie: 'subskrypcja przebiegu ról bez odpytywania' },
  { komenda: 'monitor.status', przeznaczenie: 'zbiorczy stan ról i werdykt oceny wyniku' },
];

/** Nazwa komendy wraz z przeznaczeniem; nieznana zostaje samą nazwą. */
function opisKomendy(nazwa: string): string {
  const znana = WYMAGANE.find((pozycja) => pozycja.komenda === nazwa);
  return znana === undefined ? nazwa : `${znana.komenda} (${znana.przeznaczenie})`;
}

/**
 * Wykaz komend rdzenia wraz z kontrolkami, które z niego biorą swoje zdanie.
 *
 * Jeden na scenę: cztery okna pytają rdzeń raz, nie cztery razy.
 */
export interface WykazKomendRdzenia {
  /** Pyta rdzeń o wykaz komend i przerysowuje powody wszystkich kontrolek. */
  odczytaj(): Promise<void>;
  /** Kontrolka bez pokrycia — powód dopisuje się po odpowiedzi rdzenia. */
  przyciskNieczynny(etykieta: string, ...komendy: readonly string[]): HTMLButtonElement;
  /** Wykaz komend niepokrytych jako element informacyjny okna. */
  wykazBrakow(komendy: readonly string[]): HTMLElement;
}

export function utworzWykazKomendRdzenia(kanal: Kanal): WykazKomendRdzenia {
  /** Wykaz z rdzenia; pusty, dopóki powitanie nie wróciło. */
  let komendyRdzenia: ReadonlySet<string> | null = null;
  /** Treść odmowy powitania; pusta, dopóki nic się nie zepsuło. */
  let odmowa = '';
  /** Kontrolki i wykazy, które trzeba przerysować po odpowiedzi rdzenia. */
  const zalezne: Array<() => void> = [];

  /**
   * Zdanie o pokryciu jednej komendy — trzy stany, bo trzeci naprawdę istnieje.
   *
   * Rdzeń, którego nie zapytano albo który odmówił odpowiedzi, nie orzekł
   * o braku niczego. Okno nie ma prawa zamienić tej ciszy w oskarżenie.
   */
  function zdanieOKomendzie(nazwa: string): string {
    if (odmowa !== '') {
      return `Rdzeń nie oddał wykazu komend (${odmowa}) — o pokryciu komendy ${opisKomendy(nazwa)} nic nie wiadomo.`;
    }
    if (komendyRdzenia === null) {
      return `Sprawdzanie u rdzenia, czy obsługuje komendę ${opisKomendy(nazwa)}…`;
    }
    if (!komendyRdzenia.has(nazwa)) {
      return `Rdzeń tej wersji nie rejestruje komendy ${opisKomendy(nazwa)} — zmierzone powitaniem connection.hello. Kontrakt ją niesie, więc brak jest po stronie rdzenia.`;
    }
    return `Rdzeń rejestruje komendę ${opisKomendy(nazwa)}, ale to okno jeszcze jej nie wywołuje — kontrolka niezbudowana.`;
  }

  /** Powód nieczynności kontrolki opartej na kilku komendach naraz. */
  function powod(komendy: readonly string[]): string {
    return komendy.map(zdanieOKomendzie).join(' ');
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

    przyciskNieczynny(etykieta, ...komendy) {
      const kontrolka = przyciskBezKomendy(etykieta, powod(komendy));
      zalezne.push(() => {
        const zdanie = powod(komendy);
        kontrolka.title = zdanie;
        kontrolka.setAttribute('aria-description', zdanie);
      });
      return kontrolka;
    },

    wykazBrakow(komendy) {
      const element = document.createElement('details');
      element.className = 'dm-braki';

      const naglowek = document.createElement('summary');
      element.append(naglowek);

      const lista = document.createElement('ul');
      lista.className = 'dm-braki__lista';
      element.append(lista);

      function przerysuj(): void {
        // Nagłówek nie liczy braków przed odpowiedzią rdzenia — licznik przy
        // nieznanym wykazie byłby orzeczeniem z ciszy.
        const niepokryte =
          komendyRdzenia === null
            ? null
            : komendy.filter((nazwa) => !komendyRdzenia!.has(nazwa)).length;
        naglowek.textContent =
          odmowa !== ''
            ? 'Pokrycie komend w rdzeniu — nieustalone'
            : niepokryte === null
              ? 'Pokrycie komend w rdzeniu — odczyt w toku…'
              : `Komendy tego okna nierejestrowane przez rdzeń (${niepokryte} z ${komendy.length})`;
        lista.replaceChildren(
          ...komendy.map((nazwa) => {
            const wiersz = document.createElement('li');
            wiersz.dataset['komenda'] = nazwa;
            wiersz.textContent = zdanieOKomendzie(nazwa);
            return wiersz;
          }),
        );
      }

      zalezne.push(przerysuj);
      przerysuj();
      return element;
    },
  };
}
