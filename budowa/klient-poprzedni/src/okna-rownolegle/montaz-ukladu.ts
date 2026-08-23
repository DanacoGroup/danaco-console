import type { KorzenAplikacji } from '../aplikacja/korzen-dokumentu';
import type { PolaczenieZRdzeniem } from '../aplikacja/polaczenie-z-rdzeniem';
import type { OpisOkna } from '../okno-komunikacji/opis-okna';
import { polaczPrzeplyw } from '../okno-komunikacji/przeplyw-komunikatow';
import { ID_GNIAZD } from './identyfikatory';
import { utworzUkladOkien, type UkladOkien } from './uklad-okien';

/**
 * Zamontowanie układu okien równoległych na scenie powłoki.
 *
 * Zastępuje montaż pojedynczego okna: scena dostaje układ, a układ trzyma
 * okna. Pierwsze gniazdo obejmuje okno uzgodnione z rdzeniem — to ono ma
 * identyfikator nadany komendą `window.create` i to na nim wisi przepływ
 * komunikatów. Pozostałe gniazda pracują z własnym opisem do czasu, aż rdzeń
 * otworzy dla nich osobne okna.
 *
 * Przepływ podpina się do **fasady** gniazda, nie do widoku okna. Fasada
 * przeżywa przebudowę widoku przy zmianie roli, więc zmiana roli nie zrywa
 * łączności ani nie wymaga ponownego podpięcia.
 */
export function zamontujUkladOkien(
  korzen: KorzenAplikacji,
  rdzen: PolaczenieZRdzeniem,
  opis: OpisOkna,
): UkladOkien {
  // Transport idzie do układu, bo łączność ma być widoczna w oknie. Wskaźnik
  // w pasku górnym nie sięga okna na pełnym ekranie ani dalszej kolumny sceny,
  // a polecenie wysłane po zerwanym łączu idzie w próżnię. Układ rozsyła odczyt
  // do wszystkich gniazd naraz — łącze jest jedno, więc i odpowiedź w każdym
  // nagłówku ma być ta sama.
  const uklad = utworzUkladOkien({
    podstawa: opis,
    kanal: rdzen.kanal,
    transport: rdzen.transport,
  });
  const pierwsze = uklad.fasadaOkna(ID_GNIAZD[0]);

  if (pierwsze !== null) {
    polaczPrzeplyw({
      okno: pierwsze,
      kanal: rdzen.kanal,
      transport: rdzen.transport,
      uzgodnienie: rdzen.uzgodnienie,
      opis,
    });
  }

  // Kod okna wykonania idzie z uzgodnienia, nie z opisu. Panele pomocnicze
  // wołają komendy żądające `windowId` okna otwartego, a opis okna takiego pola
  // nie ma i mieć nie może: opis powstaje w kliencie, a kod nadaje rdzeń. Okno
  // bywa już otwarte w chwili montażu (wtedy `okno()` je oddaje), a bywa, że
  // dopiero powstanie — stąd oba wejścia, nie jedno.
  const zastane = rdzen.uzgodnienie.okno();
  if (zastane !== null) uklad.ustawOknoWykonania(ID_GNIAZD[0], zastane.id);
  rdzen.uzgodnienie.naOtwarcieOkna((okno) => {
    uklad.ustawOknoWykonania(ID_GNIAZD[0], okno.id);
  });

  korzen.scena.append(uklad.element);
  return uklad;
}
