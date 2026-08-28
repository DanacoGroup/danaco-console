import './panel-orkiestracji.css';

import { WindowRole } from '../../../../shared/contract';
import type { WidokModulu } from '../rejestracja';
import type { Kanal } from '../../protokol/kanal';
import { utworzZrodloSekcjiPaneli } from '../../powloka/zrodlo-sekcji-paneli';
import { KLUCZE_SEKCJI_ORKIESTRACJI } from '../../powloka/srodowiska';
import { zamontujMultitasking } from './indeks';
import { utworzSekcjeHarmonogramu } from './sekcja-harmonogram';
import { utworzSekcjeKolejek } from './sekcja-kolejki';
import { utworzSekcjeMonitora } from './sekcja-monitor';
import { utworzSekcjeOrkiestracji } from './sekcja-orkiestracja';
import { utworzSekcjeZespolow } from './sekcja-zespoly';
import { utworzZrodloBiegu } from './zrodlo-biegu';
import { utworzZrodloNadzoru } from './zrodlo-nadzoru';
import { utworzZrodloOkien } from './zrodlo-okien';
import { utworzZrodloZespolow } from './zrodlo-zespolow';

// Panel orkiestracji jest boczną nawigacją środowiska: nawiguje po rolach, kolejkach i orkiestracji.

/** Klucz sekcji panelu — ten sam, którym wskazuje ją boczna nawigacja całego środowiska orkiestracji zadań. */
export type KluczSekcji = (typeof KLUCZE_SEKCJI_ORKIESTRACJI)[number];

/** Czy klucz pozycji nawigacji jest jedną z sześciu sekcji panelu orkiestracji rozpoznawanych przez ten plik. */
export function czySekcjaOrkiestracji(klucz: string): klucz is KluczSekcji {
  return (KLUCZE_SEKCJI_ORKIESTRACJI as readonly string[]).includes(klucz);
}

/**
 * Widok jednej sekcji panelu orkiestracji oddaje kształt modułu, bo przestrzeń
 * robocza umie osadzić dokładnie taki kształt, a sekcja zajmuje to samo miejsce
 * na planszy i ma ten sam cykl życia.
 */
export function utworzWidokSekcji(klucz: KluczSekcji, kanal: Kanal): WidokModulu {
  if (klucz === 'role') return widokRol(kanal);
  return widokSterowania(klucz, kanal);
}

/**
 * Sekcja Role — scena czterech okien roboczych.
 *
 * Montaż jest odroczony do wczytania, bo okna ról przypisuje się do okien
 * karty sesji, a sesja bywa jeszcze nieuzgodniona w chwili wyboru pozycji.
 */
function widokRol(kanal: Kanal): WidokModulu {
  const gospodarz = document.createElement('div');
  gospodarz.className = 'dm-orkiestracja dm-orkiestracja--role';
  gospodarz.dataset['sekcja'] = 'role';
  let zamontowane: { zamknij(): void } | null = null;

  return {
    element: gospodarz,
    async wczytaj(idSesji: string) {
      if (zamontowane !== null || idSesji === '') return;
      zamontowane = zamontujMultitasking(gospodarz, kanal, idSesji);
    },
    zamknij: () => zamontowane?.zamknij(),
  };
}

/** Jedna z pięciu sekcji sterowania wraz z jej cyklem życia: element, odświeżenie, adres okna i rozłączenie. */
interface SekcjaSterowania {
  element: HTMLElement;
  odswiez(): void;
  ustawOkno(idOkna: string): void;
  rozlacz?(): void;
}

/**
 * Pięć sekcji sterowania dzieli wspólny montaż, bo różnią się wyłącznie
 * treścią, a źródła powstają dopiero przy wczytaniu, nie przy tworzeniu
 * widoku.
 */
function widokSterowania(klucz: KluczSekcji, kanal: Kanal): WidokModulu {
  const gospodarz = document.createElement('div');
  let sekcja: SekcjaSterowania | null = null;

  return {
    element: gospodarz,
    async wczytaj(idSesji: string) {
      if (sekcja !== null || idSesji === '') return;
      const zbudowana = zbuduj(klucz, kanal, () => idSesji);
      sekcja = zbudowana;
      gospodarz.replaceChildren(zbudowana.element);
      // Adres układu podsekcji ustala się przed pierwszym odczytem, inaczej sekcja czytałaby układ dwa razy.
      await ustawAdres(kanal, idSesji, zbudowana);
      zbudowana.odswiez();
    },
    zamknij: () => sekcja?.rozlacz?.(),
  };
}

/** Wytwórnia sekcji według klucza; źródła współdzielą jeden wspólny kanał połączenia z uruchomionym rdzeniem. */
function zbuduj(klucz: KluczSekcji, kanal: Kanal, sesja: () => string): SekcjaSterowania {
  const sekcje = utworzZrodloSekcjiPaneli(kanal);
  switch (klucz) {
    case 'zespoly':
      return utworzSekcjeZespolow({ zrodlo: utworzZrodloZespolow(kanal), sekcje });
    case 'kolejki':
      return utworzSekcjeKolejek({
        bieg: utworzZrodloBiegu(kanal),
        nadzor: utworzZrodloNadzoru(kanal),
        sekcje,
        sesja,
      });
    case 'orkiestracja':
      return utworzSekcjeOrkiestracji({ nadzor: utworzZrodloNadzoru(kanal), sekcje });
    case 'harmonogram':
      return utworzSekcjeHarmonogramu({ nadzor: utworzZrodloNadzoru(kanal), sekcje });
    default:
      return utworzSekcjeMonitora({
        bieg: utworzZrodloBiegu(kanal),
        nadzor: utworzZrodloNadzoru(kanal),
        okna: utworzZrodloOkien(kanal),
        sekcje,
        sesja,
      });
  }
}

/**
 * Wskazuje sekcji okno, do którego przypięty jest układ jej podsekcji:
 * koordynator ma pierwszeństwo, a sesja bez koordynatora oddaje okno
 * najstarsze.
 */
async function ustawAdres(kanal: Kanal, sesja: string, sekcja: SekcjaSterowania): Promise<void> {
  const wynik = await utworzZrodloOkien(kanal).okna({ sessionId: sesja });
  if (!wynik.udany || wynik.wynik === undefined) return;
  const wedlugCzasu = [...wynik.wynik].sort((a, b) => a.createdAt - b.createdAt);
  const koordynator = wedlugCzasu.find((okno) => okno.windowRole === WindowRole.Coordinator);
  sekcja.ustawOkno(koordynator?.id ?? wedlugCzasu[0]?.id ?? '');
}
