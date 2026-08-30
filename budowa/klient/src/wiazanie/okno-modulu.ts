/**
 * Wiązanie okna modułu z rdzeniem. Znacznik niesie biblioteka Właściciela;
 * ten plik zakłada oknu sesję w rdzeniu i oddaje je wiązaniu szczegółowemu
 * modułu, a treść przykładową zdejmuje, zanim Operator ją zobaczy.
 */

import {
  Command,
  ExecutionEnv,
  PermissionMode,
  WindowRole,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';

/** Wiązania szczegółowe modułów, po kodzie rejestru rdzenia; moduł bez wpisu dostaje samo okno rdzenia. */
type WiazanieModulu = (kanal: Kanal, idOkna: string) => void;
const WIAZANIA = new Map<string, WiazanieModulu>();

/** Rejestruje wiązanie szczegółowe modułu. Woła to moduł wiązania przy wczytaniu, więc kolejność plików nie ma znaczenia. */
export function zglosWiazanieModulu(kod: string, wiazanie: WiazanieModulu): void {
  WIAZANIA.set(kod, wiazanie);
}

/**
 * Zakłada oknu modułu sesję i okno komunikacji w rdzeniu, po czym oddaje je
 * wiązaniu szczegółowemu. Nazwa środowiska wchodzi w nagłówek okna.
 */
export function zwiazOkno(kanal: Kanal, kodModulu: string, nazwaSrodowiska: string): void {
  opiszNaglowek(nazwaSrodowiska);
  void otworzOkno(kanal, kodModulu).then((idOkna) => {
    if (idOkna === '') return;
    WIAZANIA.get(kodModulu)?.(kanal, idOkna);
  });
}

/** Zakłada sesję i okno komunikacji modułu; pusty wynik znaczy, że rdzeń odmówił i wiązania szczegółowego nie ma po co wołać. */
async function otworzOkno(kanal: Kanal, kodModulu: string): Promise<string> {
  const sesja = await wywolaj(kanal, Command.SessionCreate, {});
  if (!sesja.udany || sesja.wynik === undefined) return '';

  const moduly = await wywolaj(kanal, Command.ModuleList, {});
  if (!moduly.udany || moduly.wynik === undefined) return '';
  const modul = moduly.wynik.modules.find((pozycja) => pozycja.code === kodModulu);
  if (modul === undefined) return '';

  const kanaly = await wywolaj(kanal, Command.ChannelList, { enabledOnly: true });
  const kanalModelu = kanaly.udany ? (kanaly.wynik?.channels[0]?.id ?? '') : '';
  if (kanalModelu === '') return '';

  const okno = await wywolaj(kanal, Command.WindowCreate, {
    sessionId: sesja.wynik.session.id,
    moduleId: modul.id,
    modelChannelId: kanalModelu,
    workingDirs: [],
    executionEnv: ExecutionEnv.Local,
    permissionMode: PermissionMode.Manual,
    windowRole: WindowRole.Standalone,
  });
  return okno.udany && okno.wynik !== undefined ? okno.wynik.window.id : '';
}

/** Wpisuje nazwę środowiska w nagłówek okna; pole bez wartości znika, bo podpis bez danych niczego nie mówi. */
function opiszNaglowek(nazwaSrodowiska: string): void {
  const pole = [...document.querySelectorAll('.sta-kom-naglowek .sta-kom-pole')].find(
    (kandydat) => kandydat.textContent?.startsWith('Środowisko') === true,
  );
  if (pole === undefined) return;
  if (nazwaSrodowiska === '') {
    pole.remove();
    return;
  }
  const dane = pole.querySelector('.dane');
  if (dane !== null) dane.textContent = nazwaSrodowiska;
}
