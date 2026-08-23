import { Command, type ErrorInfo } from '../../../shared/contract';
import type { Kanal } from '../protokol/kanal';
import { NAPISY } from './etykiety-rozmowy';
import type { WpisRozmowy } from './wpis-rozmowy';

/**
 * Nadanie wiadomości i przerwanie tury — dwie czynności okna rozmowy, które
 * rozmawiają z rdzeniem komendami `message.send` i `message.stop`. Wysłanie
 * w trakcie odpowiedzi samo prosi o przerwanie, więc obie czynności są od
 * siebie zależne i stoją w jednym pliku.
 *
 * Warstwa nie zna widoku — dostaje wyłącznie haczyki: co ogłosić, jak zgłosić
 * błąd, jak przestawić stan. Dzięki temu ta sama logika obsługuje okno
 * komunikacji i podgląd.
 */

/** Zależności, których nadanie i przerwanie potrzebują od okna rozmowy. */
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
 * Nadaje wiadomość, a gdy okno prowadzi turę — poprzedza ją zatrzymaniem.
 *
 * Przerwanie jest jawne: `message.send` skierowany do okna, które odpowiada,
 * odmawia kodem `conflict` zamiast skasować odpowiedź w połowie zdania, więc
 * przerwanie jedzie osobną komendą. Para komend idzie sekwencyjnie, nie
 * równolegle — nadane naraz dotarłyby w kolejności niegwarantowanej,
 * a wysłanie, które wyprzedziło zatrzymanie, trafiłoby na okno nadal zajęte
 * i odmówiło.
 *
 * Nieudane zatrzymanie nie wstrzymuje wysłania: jeżeli tura zdążyła tymczasem
 * dobiec końca sama, wysłanie przejdzie, a jeżeli nie — odmówi rdzeń.
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

/** Samo nadanie wiadomości — bez rozstrzygania o przerwaniu. */
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
 * Przerywa turę. Każde z czterech rozstrzygnięć rdzenia — błąd komendy, brak
 * trwającej tury, brak wpisu bieżącego oraz przerwanie udane — kończy się
 * widocznym śladem w wątku albo na wskaźniku, żeby kliknięcie „Zatrzymaj"
 * nigdy nie wyglądało na martwe.
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

/** Opis błędu kontraktu dla Operatora. */
function opisBledu(blad: ErrorInfo | undefined): string {
  if (blad === undefined) return 'rdzeń nie podał przyczyny';
  return `${blad.message} (${blad.code})`;
}
