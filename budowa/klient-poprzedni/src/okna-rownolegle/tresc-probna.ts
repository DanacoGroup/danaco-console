import { WindowRole } from '../../../shared/contract';
import { wpisModelu, wpisOperatora, wpisSystemowy } from '../okno-komunikacji/wpis';
import type { UkladOkien } from './uklad-okien';
import { ID_GNIAZD } from './identyfikatory';

/**
 * Treść przykładowa strony podglądu układu.
 *
 * Wyłącznie na potrzeby podglądu wizualnego: strona nie łączy się z rdzeniem,
 * więc historia okien byłaby pusta, a rozkład sceny nieczytelny. Treść jest
 * jawnie oznaczona jako przykładowa i nie występuje w aplikacji — punkt
 * wejścia klienta jej nie importuje.
 */
export function zasiejTrescProbna(uklad: UkladOkien): void {
  const [pierwsze, drugie, trzecie] = ID_GNIAZD;

  const koordynator = uklad.fasadaOkna(pierwsze);
  if (koordynator !== null) {
    koordynator.dopisz(wpisSystemowy('Podgląd układu — treść przykładowa, bez połączenia z rdzeniem.'));
    koordynator.dopisz(
      wpisOperatora('Rozpisz przegląd umowy ramowej na kroki i prowadź wykonawcę.', 'Operator'),
    );
    koordynator.dopisz(
      wpisModelu(
        'Plan na dziesięć kroków. Krok 1: wykaz załączników. Krok 2: klauzule kar umownych. '
          + 'Po każdym kroku sprawdzam wynik i układam kolejne zlecenie.',
        nazwaPersony(WindowRole.Coordinator),
      ),
    );
  }

  const wykonawca = uklad.fasadaOkna(drugie);
  if (wykonawca !== null) {
    wykonawca.dopisz(wpisSystemowy('Podgląd układu — treść przykładowa, bez połączenia z rdzeniem.'));
    wykonawca.dopisz(
      wpisModelu(
        'Krok 1 zamknięty: wykaz załączników odczytany, dwa braki odnotowane. Raport przekazany koordynatorowi.',
        nazwaPersony(WindowRole.Executor),
      ),
    );
  }

  const drugiWykonawca = uklad.fasadaOkna(trzecie);
  if (drugiWykonawca !== null) {
    drugiWykonawca.dopisz(
      wpisSystemowy('Podgląd układu — treść przykładowa, bez połączenia z rdzeniem.'),
    );
    drugiWykonawca.dopisz(
      wpisModelu(
        'Oczekuję na zlecenie. Katalog roboczy i tryb uprawnień mam własne, niezależne od pozostałych okien.',
        nazwaPersony(WindowRole.Executor),
      ),
    );
  }
}

/** Persona wypowiedzi w oknie o danej roli. */
function nazwaPersony(rola: WindowRole): string {
  return rola === WindowRole.Coordinator ? 'Koordynator' : 'Wykonawca';
}
