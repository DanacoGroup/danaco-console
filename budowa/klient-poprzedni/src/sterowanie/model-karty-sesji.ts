import { Command } from '../../../shared/contract';
import type { Kanal } from '../protokol/kanal';
import { niepowodzenie, potwierdzenie, type OdbiorcaKomunikatu } from './komunikat-zmiany';
import type { StanSterowania } from './stan-sterowania';

/**
 * Nadanie modelu całej karcie sesji — zdolność komendy `model.channel.set`,
 * której `window.update` nie ma.
 *
 * `window.update` nadaje kanał jednemu oknu; `model.channel.set` przyjmuje
 * `sessionId` i przestawia wszystkie okna karty naraz (`kanalKartySesji`
 * w `adapter_modul_model.go`). Karta z czterema oknami to jedno żądanie zamiast
 * czterech, a więc i jedna okazja do niepowodzenia zamiast czterech.
 *
 * Żądanie niesie `agentId` obok `modelChannelId`, bo karta ma dostać ten sam
 * wybór, który stoi w oknie: kanał surowy albo agenta nałożonego na kanał.
 * Pominięcie agenta przestawiłoby pozostałe okna na model surowy.
 *
 * Przycisk jest zawsze klikalny. Karta jednookienna nie jest powodem do
 * wyszarzenia — nadanie jednemu oknu tą drogą daje ten sam skutek co przez okno.
 */

const ETYKIETA = 'Zastosuj do całej karty sesji';

/** Wybór okna przełożony na żądanie karty; kształt zgodny z `zlecenieWyboru`. */
export interface WyborModelu {
  modelChannelId?: string;
  agentId: string;
}

export interface ModelKartySesji {
  element: HTMLElement;
  /** Zapamiętuje wybór, który przycisk rozniesie na kartę. */
  ustawWybor(wybor: WyborModelu): void;
}

export function utworzModelKartySesji(
  stan: StanSterowania,
  kanal: Kanal,
  zglos: OdbiorcaKomunikatu,
): ModelKartySesji {
  let ostatni: WyborModelu | null = null;

  const element = document.createElement('button');
  element.type = 'button';
  element.className = 'dn-btn dn-btn--sm dn-btn--atrament';
  element.textContent = ETYKIETA;

  element.addEventListener('click', () => {
    const okno = stan.migawka().okno;
    const idSesji = okno.sessionId;
    if (idSesji === '') {
      // Okno bez karty nie jest awarią — jest stanem, w którym nie ma czego
      // rozesłać. Mówimy to wprost zamiast wysyłać żądanie bez adresata.
      zglos({ tresc: 'Okno nie należy do żadnej karty sesji — nie ma na co rozesłać modelu.', udany: false });
      return;
    }
    // Bez wcześniejszego wyboru rozsyłamy nastawę bieżącą okna, żeby model
    // ustawiony już w oknie nie wymagał wybrania go drugi raz.
    const wybor: WyborModelu = ostatni ?? {
      modelChannelId: okno.modelChannelId,
      agentId: okno.agentId ?? '',
    };
    const kanalModelu = wybor.modelChannelId ?? okno.modelChannelId;
    if (kanalModelu === '') {
      zglos({ tresc: 'Okno nie ma nadanego kanału modelu — nie ma czego rozesłać na kartę.', udany: false });
      return;
    }
    kanal.wyslij(
      Command.ModelChannelSet,
      { modelChannelId: kanalModelu, sessionId: idSesji, agentId: wybor.agentId },
      (wynik) => {
        if (!wynik.udany) {
          zglos(niepowodzenie(ETYKIETA, wynik.blad));
          return;
        }
        // Potwierdzenie nie podaje liczby przestawionych okien:
        // `ModelChannelSetResponse` niesie kanał, okno i kartę, a nie wykaz
        // okien objętych czynnością. Rdzeń rozgłasza je osobno zdarzeniem
        // `window.changed`, a każde okno przyjmuje swoje.
        zglos(potwierdzenie(ETYKIETA));
      },
    );
  });

  return {
    element,
    ustawWybor(wybor) {
      ostatni = wybor;
    },
  };
}
