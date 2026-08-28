import { Command } from '../../../shared/contract';
import type { Kanal } from '../protokol/kanal';
import { niepowodzenie, potwierdzenie, type OdbiorcaKomunikatu } from './komunikat-zmiany';
import type { StanSterowania } from './stan-sterowania';

// Nadanie modelu całej karcie sesji: żądanie karty przestawia wszystkie okna karty jednym wywołaniem.

const ETYKIETA = 'Zastosuj do całej karty sesji';

/** Wybór okna przełożony na żądanie karty sesji; kształt zgodny z wynikiem funkcji zlecenieWyboru modelu. */
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
      // Okno bez karty nie jest awarią: nie ma czego rozesłać, więc mówimy to wprost.
      zglos({ tresc: 'Okno nie należy do żadnej karty sesji — nie ma na co rozesłać modelu.', udany: false });
      return;
    }
    // Bez wcześniejszego wyboru rozsyłamy nastawę bieżącą okna, bez ponownego jej wybierania.
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
        // Potwierdzenie nie podaje liczby okien; rdzeń rozgłasza je osobno zdarzeniem zmiany okna.
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
