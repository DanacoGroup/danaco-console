import { WindowStatus, type AutomationStep } from '../../../../shared/contract';
import { czyObiekt, czyTablica } from '../../protokol/ksztalt-odpowiedzi';
import { KOD_MODULU } from './kody-okien';
import { OBSZAR_KONTEKSTU_ROZMOWY, type ZrodloAutomations } from './zrodlo-automations';

/**
 * Scenariusz przeniesiony z innego modułu za pomocą komendy context.transfer,
 * odczytany z konfiguracji sesji okna modułu docelowego, gotowy do zapisu
 * jako automatyka.
 */
export interface PrzekazanyScenariusz {
  /** Okno modułu, do którego rdzeń przeniósł komplet. */
  idOkna: string;
  /** Nazwa scenariusza podana przez moduł przekazujący. */
  nazwa: string;
  /** Opis scenariusza; pusty, gdy moduł przekazujący go nie podał. */
  opis: string;
  /** Czy scenariusz ma być czynny po zapisie. */
  czynny: boolean;
  /** Kroki scenariusza w kolejności wykonania. */
  kroki: AutomationStep[];
  /** Polecenie wyjściowe kompletu — zdanie, którym moduł opisał przekazanie. */
  polecenie: string;
}

/**
 * Skutek odczytu przekazań: powodzenie niesie wykaz scenariuszy i liczbę
 * sprawdzonych okien, niepowodzenie niesie zdanie z przyczyną braku wykazu.
 */
export type SkutekOdczytuPrzekazan =
  | { odczytane: true; scenariusze: PrzekazanyScenariusz[]; oknaSprawdzone: number }
  | { odczytane: false; zdanie: string };

/**
 * Odczytuje scenariusze przeniesione do modułu.
 *
 * Okna zamknięte pomijamy: przekazanie do okna zamkniętego jest przekazaniem
 * odbytym i zakończonym, a wykaz ma pokazywać to, co czeka na decyzję.
 */
export async function odczytajPrzekazania(
  zrodlo: ZrodloAutomations,
): Promise<SkutekOdczytuPrzekazan> {
  const okna = await zrodlo.oknaKomunikacji({ status: WindowStatus.Open });
  if (!okna.udany || okna.wynik === undefined) {
    return {
      odczytane: false,
      zdanie:
        'Rdzeń nie oddał wykazu okien komunikacji, więc nie wiadomo, czy któryś moduł przekazał tu scenariusz. ' +
        'Przekazanie idzie komendą context.transfer i zakłada okno tego modułu — bez wykazu okien nie ma do czego zajrzeć.',
    };
  }

  const nasze = okna.wynik.filter((okno) => okno.moduleId === KOD_MODULU);
  const scenariusze: PrzekazanyScenariusz[] = [];
  for (const okno of nasze) {
    const komplet = await zrodlo.kontekstOkna({
      windowId: okno.id,
      areas: [OBSZAR_KONTEKSTU_ROZMOWY],
    });
    if (!komplet.udany || komplet.wynik === undefined) continue;
    const przeniesiony = komplet.wynik.config.conversationContext?.transferredContext;
    if (przeniesiony === undefined) continue;
    const scenariusz = scenariuszZKompletu(okno.id, przeniesiony.prompt ?? '', przeniesiony.executionParams);
    if (scenariusz !== null) scenariusze.push(scenariusz);
  }
  return { odczytane: true, scenariusze, oknaSprawdzone: nasze.length };
}

/**
 * Wyjmuje scenariusz z parametrów wykonania przeniesionego kompletu. Komplet
 * bez nazwy albo bez ani jednego kroku zwraca `null`, ponieważ nie da się
 * z niego zapisać automatyki.
 */
export function scenariuszZKompletu(
  idOkna: string,
  polecenie: string,
  parametry: unknown,
): PrzekazanyScenariusz | null {
  if (!czyObiekt(parametry)) return null;
  const nazwa = typeof parametry['name'] === 'string' ? parametry['name'].trim() : '';
  if (nazwa === '') return null;
  const surowe = parametry['steps'];
  if (!czyTablica(surowe)) return null;
  const kroki = surowe.filter(czyKrok);
  if (kroki.length === 0) return null;
  return {
    idOkna,
    nazwa,
    opis: typeof parametry['description'] === 'string' ? parametry['description'].trim() : '',
    czynny: parametry['enabled'] !== false,
    kroki,
    polecenie,
  };
}

/**
 * Czy pozycja przeniesionego wykazu jest krokiem automatyki. Sprawdza pola
 * obowiązkowe kontraktu `id` i `kind`; pozostałe pola pozostają
 * nieobowiązkowe.
 */
function czyKrok(pozycja: unknown): pozycja is AutomationStep {
  return (
    czyObiekt(pozycja) &&
    typeof pozycja['id'] === 'string' &&
    pozycja['id'].trim() !== '' &&
    typeof pozycja['kind'] === 'string'
  );
}

/**
 * Zdanie o wyniku odczytu przekazań, przeznaczone do wyświetlenia nad wykazem
 * w oknie modułu: informuje o braku przekazań, liczbie sprawdzonych okien
 * albo liczbie oczekujących decyzji.
 */
export function zdanieOPrzekazaniach(skutek: SkutekOdczytuPrzekazan): string {
  if (!skutek.odczytane) return skutek.zdanie;
  if (skutek.oknaSprawdzone === 0) {
    return 'Żaden moduł nie przekazał tu jeszcze scenariusza — rdzeń nie ma ani jednego otwartego okna tego modułu.';
  }
  if (skutek.scenariusze.length === 0) {
    return (
      `Sprawdzono ${skutek.oknaSprawdzone} otwartych okien tego modułu; żadne nie niesie przeniesionego ` +
      'scenariusza z nazwą i krokami.'
    );
  }
  return (
    `Przekazania oczekujące na decyzję: ${skutek.scenariusze.length}. ` +
    'Zapis do magazynu automatyk jest osobną czynnością — przekazanie samo niczego nie zapisuje.'
  );
}
