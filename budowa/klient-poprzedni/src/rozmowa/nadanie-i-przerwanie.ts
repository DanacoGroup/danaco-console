import { Command, type ErrorInfo } from '../../../shared/contract';
import type { Kanal } from '../protokol/kanal';
import { NAPISY } from './etykiety-rozmowy';
import type { WpisRozmowy } from './wpis-rozmowy';

/**
 * Plik łączy nadanie wiadomości i przerwanie tury, bo obie czynności zależą od siebie.
 */

/**
 * Zależności, których nadanie i przerwanie potrzebują od okna rozmowy, dostarczane przez
 * wywołującego.
 */
export interface OtoczenieNadania {
  kanal: Kanal;
  idOkna: string;
  /** Wpis tury bieżącej albo `undefined`, gdy żadna nie trwa. */
  biezaca: () => WpisRozmowy | undefined;
  /** Ogłasza zmieniony wpis widokowi. */
  oglos: (wpis: WpisRozmowy) => void;
  /** Sadza błąd komendy na wpisie oczekującym. */
  zglosBlad: (wpis: WpisRozmowy, blad: ErrorInfo | undefined) => void;
  /** Wpisuje komunikat systemowy do wątku. */
  zglosKomunikat: (tresc: string) => void;
  /** Przestawia wskaźnik wysyłania. */
  zakonczWysylanie: () => void;
}

/**
 * Nadaje wiadomość, a gdy okno prowadzi turę, poprzedza wysłanie jawnym przerwaniem
 * biegnącej odpowiedzi.
 */
export function nadajZPrzerwaniem(
  o: OtoczenieNadania,
  tresc: string,
  odpowiedz: WpisRozmowy,
): void {
  const przerwana = o.biezaca();
  if (przerwana === undefined || przerwana.domkniety) {
    nadaj(o, tresc, odpowiedz);
    return;
  }
  o.kanal.wyslij(
    Command.MessageStop,
    { windowId: o.idOkna, messageId: przerwana.idWiadomosci },
    () => {
      if (!przerwana.domkniety) {
        przerwana.stan = 'przerwany';
        przerwana.domkniety = true;
        o.oglos(przerwana);
      }
      nadaj(o, tresc, odpowiedz);
    },
  );
}

/**
 * Samo nadanie wiadomości do rdzenia komendą message.send, bez rozstrzygania o przerwaniu
 * biegnącej tury okna.
 */
function nadaj(o: OtoczenieNadania, tresc: string, odpowiedz: WpisRozmowy): void {
  o.kanal.wyslij(
    Command.MessageSend,
    { windowId: o.idOkna, content: tresc, stream: true },
    (wynik) => {
      if (wynik.udany) return;
      o.zglosBlad(odpowiedz, wynik.blad);
    },
  );
}

/**
 * Przerywa turę, a każde z czterech rozstrzygnięć rdzenia kończy się widocznym śladem w
 * wątku rozmowy albo na wskaźniku wysyłania.
 */
export function przerwijTure(o: OtoczenieNadania): void {
  const biezaca = o.biezaca();
  o.kanal.wyslij(
    Command.MessageStop,
    { windowId: o.idOkna, messageId: biezaca?.idWiadomosci },
    (wynik) => {
      if (!wynik.udany) {
        o.zglosKomunikat(`Zatrzymanie nieudane — ${opisBledu(wynik.blad)}`);
        return;
      }
      if (wynik.wynik?.stopped !== true) {
        o.zglosKomunikat(NAPISY.zatrzymanieBezTury);
        return;
      }
      if (biezaca === undefined) {
        o.zglosKomunikat(NAPISY.zatrzymanieBezWpisu);
        o.zakonczWysylanie();
        return;
      }
      biezaca.stan = 'przerwany';
      biezaca.domkniety = true;
      o.oglos(biezaca);
      o.zakonczWysylanie();
    },
  );
}

/**
 * Opis błędu kontraktu dla Operatora, złożony z treści komunikatu rdzenia oraz kodu błędu
 * zwróconego przez rdzeń.
 */
function opisBledu(blad: ErrorInfo | undefined): string {
  if (blad === undefined) return 'rdzeń nie podał przyczyny';
  return `${blad.message} (${blad.code})`;
}
