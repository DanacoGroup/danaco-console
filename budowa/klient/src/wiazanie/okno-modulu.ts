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
import { zwiazDeveloper } from './okno-developer.ts';

/** Wiązania szczegółowe modułów, po kodzie rejestru rdzenia; moduł bez wpisu dostaje samo okno rdzenia. */
type WiazanieModulu = (kanal: Kanal, idOkna: string) => void;
const WIAZANIA = new Map<string, WiazanieModulu>([['developer', zwiazDeveloper]]);

/**
 * Czy moduł ma wiązanie wypełniające jego wnętrze odpowiedzią rdzenia. Wnętrze
 * bez wiązania pokazałoby treść przykładową prototypu jako pracę Operatora,
 * więc takiego okna wydanie nie stawia wcale.
 */
export function maWiazanie(kod: string): boolean {
  return WIAZANIA.has(kod);
}

/** Rejestruje wiązanie szczegółowe modułu. Woła to moduł wiązania przy wczytaniu, więc kolejność plików nie ma znaczenia. */
export function zglosWiazanieModulu(kod: string, wiazanie: WiazanieModulu): void {
  WIAZANIA.set(kod, wiazanie);
}

/**
 * Zakłada oknu modułu sesję i okno komunikacji w rdzeniu, po czym oddaje je
 * wiązaniu szczegółowemu. Nazwa środowiska wchodzi w nagłówek okna.
 */
export function zwiazOkno(kanal: Kanal, kodModulu: string, nazwaSrodowiska: string): void {
  void opiszNaglowek(kanal, nazwaSrodowiska);
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

/**
 * Opisuje nagłówek okna komunikacji środowiskiem wejścia i modelem kanału.
 * Pozostałe podpisy prototypu — wysiłek, wykonawca, pamięć, rola, format,
 * tura — schodzą: ich wartości niesie dopiero praca podjęta w oknie, a okno
 * dopiero co stanęło, więc podpis pokazywałby wartość wymyśloną.
 */
async function opiszNaglowek(kanal: Kanal, nazwaSrodowiska: string): Promise<void> {
  const naglowek = document.querySelector('.sta-kom-naglowek');
  if (naglowek === null) return;
  const wynik = await wywolaj(kanal, Command.ChannelList, { enabledOnly: true });
  const kanalModelu = wynik.udany ? wynik.wynik?.channels[0] : undefined;
  const wartosci: Record<string, string> = {
    'Środowisko': nazwaSrodowiska,
    Model: kanalModelu?.model ?? kanalModelu?.name ?? '',
  };
  for (const pole of [...naglowek.querySelectorAll('.sta-kom-pole')]) {
    const podpis = Object.keys(wartosci).find(
      (kandydat) => pole.textContent?.startsWith(kandydat) === true,
    );
    wpiszPole([pole], podpis ?? '', podpis === undefined ? '' : wartosci[podpis]);
  }
}

/** Wpisuje wartość w pole nagłówka rozpoznane po jego podpisie; wartość pusta zdejmuje całe pole, bo pole bez wartości niczego nie mówi. */
export function wpiszPole(pola: Element[], podpis: string, wartosc: string): void {
  const pole = podpis === ''
    ? pola[0]
    : pola.find((kandydat) => kandydat.textContent?.startsWith(podpis) === true);
  if (pole === undefined) return;
  if (wartosc === '') {
    pole.remove();
    return;
  }
  /* Wartość stoi w prototypie raz w `.dane`, raz w samym wyróżnieniu — pole
     środowiska nie niesie klasy danych, a wpisać się musi tak samo. */
  const dane = pole.querySelector('.dane') ?? pole.querySelector('b');
  if (dane !== null) dane.textContent = wartosc;
}
