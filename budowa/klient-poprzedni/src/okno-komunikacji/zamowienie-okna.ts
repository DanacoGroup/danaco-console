import type { ZamowienieOkna } from '../protokol/uzgodnienie';
import type { OpisOkna } from './opis-okna';

/**
 * Przekład opisu okna na treść komendy `window.create`.
 *
 * Nazwy pól należą do kontraktu, opis okna jest ich odpowiednikiem po stronie
 * widoku. Przekład leży tutaj, żeby warstwa protokołu nie znała widoku,
 * a widok nie budował koperty.
 */
export function zamowienieOkna(opis: OpisOkna): ZamowienieOkna {
  return {
    moduleId: opis.modul,
    modelChannelId: opis.kanalModelu,
    workingDirs: opis.katalogiRobocze,
    executionEnv: opis.srodowiskoWykonania,
    permissionMode: opis.trybUprawnien,
    windowRole: opis.rola,
    title: opis.tytul,
  };
}
