import type { KorzenAplikacji } from '../aplikacja/korzen-dokumentu';
import type { PolaczenieZRdzeniem } from '../aplikacja/polaczenie-z-rdzeniem';
import type { OpisOkna } from '../okno-komunikacji/opis-okna';
import { polaczPrzeplyw } from '../okno-komunikacji/przeplyw-komunikatow';
import { ID_GNIAZD } from './identyfikatory';
import { utworzUkladOkien, type UkladOkien } from './uklad-okien';

/**
 * Zamontowanie układu okien równoległych na scenie powłoki zastępuje montaż pojedynczego okna, a przepływ komunikatów podpina się do fasady gniazda, nie do widoku okna, więc zmiana roli nie zrywa łączności.
 */
export function zamontujUkladOkien(
  korzen: KorzenAplikacji,
  rdzen: PolaczenieZRdzeniem,
  opis: OpisOkna,
): UkladOkien {
  // Transport idzie do układu, bo łączność ma być widoczna, a układ rozsyła odczyt do wszystkich gniazd.
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

  // Kod okna wykonania idzie z uzgodnienia, nie z opisu — opis powstaje w kliencie, kod nadaje rdzeń.
  const zastane = rdzen.uzgodnienie.okno();
  if (zastane !== null) uklad.ustawOknoWykonania(ID_GNIAZD[0], zastane.id);
  rdzen.uzgodnienie.naOtwarcieOkna((okno) => {
    uklad.ustawOknoWykonania(ID_GNIAZD[0], okno.id);
  });

  korzen.scena.append(uklad.element);
  return uklad;
}
