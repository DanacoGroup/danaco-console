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

/**
 * Panel orkiestracji — boczna nawigacja środowiska MultitaskingAI.
 *
 * Środowisko nie udostępnia modułów w bocznej nawigacji: nawiguje po tym, czym
 * się w nim steruje — po rolach, kolejkach i orkiestracji. Okien do modułu
 * nie przypinamy.
 *
 * Sekcja Role montuje istniejącą scenę czterech okien roboczych: Coordinator
 * Chat, dwa niezależne wystąpienia Executor Chat oraz Results Analyzer, wraz
 * z panelem Subagent Network w oknie wykonawcy. Executor 1 i 2 to jeden typ
 * okna w dwóch wystąpieniach, nie dwa osobne okna.
 *
 * Pozostałe cztery sekcje — Kolejki, Orkiestracja, Harmonogram i Monitor —
 * pracują na oknach modułu Automations, bo własnych okien nie mają. Każda
 * niesie podsekcję powiązania, która mówi, czy powiązanie jest skonfigurowane
 * i jak je założyć, zamiast rysować pustkę albo udawać własne okno.
 *
 * Sekcja nieznana temu plikowi nie gaśnie i nie znika: oddajemy `undefined`,
 * a przestrzeń robocza pokazuje stan pusty z nazwą sekcji. Brak ustawienia
 * układu podsekcji znaczy układ domyślny, nigdy niedostępność.
 *
 * Adresem układu jest okno: `panel.sections.*` adresuje parę (okno, panel),
 * a panel orkiestracji należy do środowiska. Adresem zostaje okno koordynatora,
 * a gdy obsady jeszcze nie ma — pierwsze okno karty sesji. Sesja bez ani
 * jednego okna zostawia układ miejscowy, o czym powłoka mówi przy pierwszej
 * zmianie.
 */

/** Klucz sekcji panelu — ten sam, którym wskazuje ją boczna nawigacja. */
export type KluczSekcji = (typeof KLUCZE_SEKCJI_ORKIESTRACJI)[number];

/** Czy klucz pozycji nawigacji jest jedną z sześciu sekcji panelu. */
export function czySekcjaOrkiestracji(klucz: string): klucz is KluczSekcji {
  return (KLUCZE_SEKCJI_ORKIESTRACJI as readonly string[]).includes(klucz);
}

/**
 * Widok jednej sekcji panelu orkiestracji.
 *
 * Oddaje `WidokModulu`, bo przestrzeń robocza umie osadzić dokładnie taki
 * kształt — sekcja nie jest modułem, ale zajmuje to samo miejsce na planszy
 * i ma ten sam cykl życia (element powstaje raz, treść dociąga się z sesją).
 * Drugiego kształtu widoku nie wprowadzamy.
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

/** Jedna z pięciu sekcji sterowania wraz z jej cyklem życia. */
interface SekcjaSterowania {
  element: HTMLElement;
  odswiez(): void;
  ustawOkno(idOkna: string): void;
  rozlacz?(): void;
}

/**
 * Pięć sekcji sterowania — wspólny montaż, bo różnią się wyłącznie treścią.
 *
 * Źródła powstają dopiero przy wczytaniu, a nie przy tworzeniu widoku: gdyby
 * powstawały wcześniej, subskrypcje zdarzeń stanęłyby dla sekcji, której
 * Operator jeszcze nie otworzył.
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
      // Adres układu podsekcji ustala się przed pierwszym odczytem: inaczej
      // sekcja przeczytałaby układ domyślny i zaraz po nim ten sam układ
      // z rdzenia — dwa odczyty o to samo.
      await ustawAdres(kanal, idSesji, zbudowana);
      zbudowana.odswiez();
    },
    zamknij: () => sekcja?.rozlacz?.(),
  };
}

/** Wytwórnia sekcji według klucza; źródła współdzielą jeden kanał. */
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
 * Wskazuje sekcji okno, do którego przypięty jest układ jej podsekcji.
 *
 * Koordynator ma pierwszeństwo, bo z jego okna steruje się całym środowiskiem.
 * Sesja bez koordynatora oddaje okno najstarsze — układ ma gdzie zamieszkać
 * także przed obsadzeniem ról. Sesja bez ani jednego okna zostawia adres pusty,
 * a powłoka mówi to przy pierwszej zmianie układu.
 */
async function ustawAdres(kanal: Kanal, sesja: string, sekcja: SekcjaSterowania): Promise<void> {
  const wynik = await utworzZrodloOkien(kanal).okna({ sessionId: sesja });
  if (!wynik.udany || wynik.wynik === undefined) return;
  const wedlugCzasu = [...wynik.wynik].sort((a, b) => a.createdAt - b.createdAt);
  const koordynator = wedlugCzasu.find((okno) => okno.windowRole === WindowRole.Coordinator);
  sekcja.ustawOkno(koordynator?.id ?? wedlugCzasu[0]?.id ?? '');
}
