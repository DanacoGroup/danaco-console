import { PermissionMode, type Window } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import { ustawPozycje, poleWyboru } from '../../modele/kontrolki-formularza';
import { utworzStanOkna, type StanOkna } from './stan-okna';
import type { ZrodloZaplecza } from './zrodlo-zaplecza';

/**
 * Zasięg uprawnień w podziale na okna komunikacji pokazuje stan wszystkich
 * okien sesji naraz. Tryb okna niesie pole `Window.permissionMode`, a zmienia
 * je komenda `window.update`; każda odmowa cofa listę wyboru do trybu okna.
 */
export interface ZasiegOkien {
  element: HTMLElement;
  /** Odczytuje okna wskazanej sesji i pokazuje ich tryby. */
  wczytaj(idSesji: string): Promise<void>;
}

/**
 * Nazwy trybów uprawnień okna wraz ze zdaniem o znaczeniu każdego z nich.
 * Wykaz obejmuje komplet wartości wyliczenia `PermissionMode`, ponieważ lista
 * wyboru ma podać Operatorowi każdy tryb, który rdzeń przyjmie.
 */
const OPISY_TRYBOW: Record<PermissionMode, string> = {
  [PermissionMode.Manual]: 'pytanie o zgodę przed każdą zmianą',
  [PermissionMode.AcceptEdits]: 'automatyczna zgoda na zmiany plików',
  [PermissionMode.Plan]: 'praca planistyczna bez zmian w systemie',
  [PermissionMode.Auto]: 'decyzje o uprawnieniach podejmuje model',
  [PermissionMode.DontAsk]: 'bez zapytań, z zachowaniem ograniczeń',
  [PermissionMode.BypassPermissions]: 'pominięcie kontroli uprawnień',
};

/**
 * Etykieta okna komunikacji w wykazie zasięgu składa się z identyfikatora
 * modułu oraz tytułu okna, a przy tytule pustym — z identyfikatora okna. Okno
 * bez tytułu pozostaje wtedy rozpoznawalne zamiast zlewać się z pozostałymi.
 */
function etykietaOkna(okno: Window): string {
  const tytul = (okno.title ?? '').trim();
  return tytul === '' ? `${okno.moduleId} · ${okno.id}` : `${okno.moduleId} · ${tytul}`;
}

export function utworzZasiegOkien(zaplecze: ZrodloZaplecza): ZasiegOkien {
  const okno: StanOkna = utworzStanOkna();

  const wykaz = poleWyboru({ etykieta: 'Okno komunikacji sesji' }, []);
  const tryb = poleWyboru(
    { etykieta: 'Tryb uprawnień okna (--permission-mode)' },
    Object.values(PermissionMode).map((wartosc) => ({
      wartosc,
      etykieta: `${wartosc} — ${OPISY_TRYBOW[wartosc]}`,
    })),
  );

  const odpowiedz = document.createElement('p');
  odpowiedz.className = 'da-odpowiedz';
  odpowiedz.hidden = true;

  okno.tresc.append(wykaz.element, tryb.element, odpowiedz);

  const tytul = document.createElement('h4');
  tytul.className = 'da-panel__tytul';
  tytul.textContent = 'Zasięg per okno komunikacji';

  const element = document.createElement('section');
  element.className = 'da-panel da-zasieg';
  element.append(tytul, okno.element);

  let okna: Window[] = [];

  function powiedz(tresc: string, powodzenie: boolean): void {
    odpowiedz.textContent = tresc;
    odpowiedz.hidden = tresc === '';
    odpowiedz.dataset['powodzenie'] = String(powodzenie);
  }

  function pokazTrybWybranego(): void {
    const wybrane = okna.find((wpis) => wpis.id === wykaz.kontrolka.value);
    if (wybrane === undefined) return;
    tryb.kontrolka.value = wybrane.permissionMode;
  }

  async function zapiszTryb(): Promise<void> {
    const idOkna = wykaz.kontrolka.value;
    if (idOkna === '') {
      powiedz('Wskaż okno komunikacji — tryb uprawnień należy do okna, nie do sesji.', false);
      return;
    }
    powiedz('Zapis trybu uprawnień okna…', true);
    const wynik = await zaplecze.ustawTrybOkna(idOkna, tryb.kontrolka.value as PermissionMode);
    if (!wynik.udany || wynik.wynik === undefined) {
      pokazTrybWybranego();
      powiedz(opisOdmowy('Zapis trybu uprawnień', wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    const zmienione = wynik.wynik.window;
    okna = okna.map((wpis) => (wpis.id === zmienione.id ? zmienione : wpis));
    powiedz(`Okno ${etykietaOkna(zmienione)} pracuje w trybie ${zmienione.permissionMode}.`, true);
  }

  wykaz.kontrolka.addEventListener('change', pokazTrybWybranego);
  tryb.kontrolka.addEventListener('change', () => void zapiszTryb());

  return {
    element,

    async wczytaj(idSesji) {
      if (idSesji === '') {
        okno.puste('Sesja nie jest jeszcze otwarta — zasięg per okno nie ma czego pokazać.');
        return;
      }
      okno.ladowanie('Odczyt okien komunikacji sesji w toku…');
      const wynik = await zaplecze.okna(idSesji);
      if (!wynik.udany || wynik.wynik === undefined) {
        okno.blad(opisOdmowy('Odczyt okien sesji', wynik.blad?.code, wynik.blad?.message));
        return;
      }
      okna = wynik.wynik.windows;
      ustawPozycje(
        wykaz.kontrolka,
        okna.map((wpis) => ({ wartosc: wpis.id, etykieta: etykietaOkna(wpis) })),
      );
      if (okna.length === 0) {
        okno.puste('Sesja nie ma jeszcze okna komunikacji.');
        return;
      }
      pokazTrybWybranego();
      okno.gotowe();
    },
  };
}
