import { Command, type ErrorInfo, type Window, type WindowRole } from '../../../shared/contract';
import type { OpisOkna } from '../okno-komunikacji/opis-okna';
import { zamowienieOkna } from '../okno-komunikacji/zamowienie-okna';
import type { Kanal } from '../protokol/kanal';

/** Wynik zamówienia okna: rdzeń otworzył okno albo podał przyczynę odmowy. */
export interface OdpowiedzNaZamowienie {
  /** Okno otwarte przez rdzeń; `null`, gdy komenda się nie powiodła. */
  okno: Window | null;
  /** Przyczyna niepowodzenia, jeżeli rdzeń ją podał. */
  blad?: ErrorInfo;
}

/**
 * Zamówienie kolejnego okna komunikacji w rdzeniu.
 *
 * Jedna odpowiedzialność: wysłanie komendy `window.create` dla okna innego
 * niż uzgodnione przy nawiązaniu połączenia. Pierwsze okno otwiera
 * uzgodnienie (`protokol/uzgodnienie`); drugie i trzecie powstają dopiero
 * wtedy, gdy Operator wprowadza je na scenę, i przechodzą tą samą drogą
 * kontraktu — mają własny identyfikator, własną historię i własne ustawienia.
 *
 * Nazwa komendy i kształt żądania pochodzą z `shared/contract`;
 * przekładem opisu okna na treść żądania zajmuje się `zamowienieOkna`, więc
 * ten plik nie zna ani jednej nazwy pola kontraktu.
 *
 * Rola nadpisuje rolę z opisu, bo o roli okna na scenie rozstrzyga jego
 * miejsce w figurze koordynator–wykonawca, a nie konfiguracja
 * budowania.
 */
export function zamowOknoRdzenia(
  kanal: Kanal,
  opis: OpisOkna,
  rola: WindowRole,
  tytul: string,
  przyOdpowiedzi: (odpowiedz: OdpowiedzNaZamowienie) => void,
): void {
  kanal.wyslij(
    Command.WindowCreate,
    {
      ...zamowienieOkna(opis),
      sessionId: kanal.sesja().id(),
      windowRole: rola,
      title: tytul,
    },
    (wynik) => {
      przyOdpowiedzi({ okno: wynik.wynik?.window ?? null, blad: wynik.blad });
    },
  );
}
