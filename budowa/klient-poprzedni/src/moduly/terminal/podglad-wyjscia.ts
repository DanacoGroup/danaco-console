import {
  TerminalProcessStatus,
  type TerminalOutputReadResponse,
} from '../../../../shared/contract';

/**
 * Podgląd wyjścia procesu czyta treść przechowywaną w rdzeniu, osobno dla każdej pozycji wykazu.
 */

/**
 * Zdanie o stanie procesu w podglądzie; czasownik dobiera się ze stanu odpowiedzi rdzenia, nie z samego powodzenia odczytu.
 */
export function zdanieStanuWyjscia(odpowiedz: TerminalOutputReadResponse, idProcesu: string): string {
  const kod =
    odpowiedz.exitCode === undefined
      ? 'kodu wyjścia rdzeń nie podał'
      : `kod wyjścia ${odpowiedz.exitCode}`;
  switch (odpowiedz.status) {
    case TerminalProcessStatus.Running:
      // Proces biegnący nie ma kodu wyjścia, więc zdanie o jego braku wprowadzałoby w błąd.
      return `Proces ${idProcesu} wciąż biegnie — to, co niżej, jest wyjściem dotychczasowym, nie całością.`;
    case TerminalProcessStatus.Finished:
      return `Proces ${idProcesu} zakończył się — ${kod}.`;
    case TerminalProcessStatus.Failed:
      return `Proces ${idProcesu} zakończył się BŁĘDEM — ${kod}.`;
    case TerminalProcessStatus.Stopped:
      return `Proces ${idProcesu} został zatrzymany — ${kod}.`;
    default:
      // Stan spoza wyliczenia kontraktu: widok pokazuje go dosłownie, zamiast zgadywać znaczenie.
      return `Proces ${idProcesu} — rdzeń oddał stan „${String(odpowiedz.status)}", którego kontrakt nie wymienia (${kod}).`;
  }
}

/**
 * Zdanie o przycięciu treści w podglądzie wyjścia; pozostaje puste, gdy rdzeń oddał komplet strumienia.
 */
export function zdaniePrzyciecia(odpowiedz: TerminalOutputReadResponse): string {
  if (!odpowiedz.truncated) return '';
  const ile =
    odpowiedz.truncatedBytes === undefined
      ? 'Ilu bajtów rdzeń nie podał.'
      : `Odciętych bajtów: ${odpowiedz.truncatedBytes}.`;
  return `Rdzeń oddał FRAGMENT wyjścia, nie całość — granica rozmiaru albo wskazanie „ostatnie wiersze". ${ile}`;
}

/**
 * Węzeł podglądu jednego procesu.
 *
 * @param idProcesu proces, o którego wyjście pytano — stoi w każdym zdaniu, bo
 *   podgląd bywa otwarty przy kilku pozycjach wykazu naraz.
 */
export function podgladWyjscia(
  idProcesu: string,
  odpowiedz: TerminalOutputReadResponse,
): HTMLElement {
  const element = document.createElement('div');
  element.className = 'dt-podglad';
  element.dataset['proces'] = idProcesu;
  element.dataset['stan'] = odpowiedz.status;

  const stan = document.createElement('p');
  stan.className = 'dt-podglad__stan';
  stan.textContent = zdanieStanuWyjscia(odpowiedz, idProcesu);
  element.append(stan);

  const przyciecie = zdaniePrzyciecia(odpowiedz);
  if (przyciecie !== '') {
    const uwaga = document.createElement('p');
    uwaga.className = 'dt-podglad__przyciecie';
    uwaga.textContent = przyciecie;
    element.append(uwaga);
  }

  element.append(
    strumien('stdout', 'Wyjście zwykłe (stdout)', odpowiedz.stdout),
    strumien('stderr', 'Wyjście diagnostyczne (stderr)', odpowiedz.stderr),
  );
  return element;
}

/**
 * Jeden strumień wraz z podpisem.
 *
 * Rdzeń oddaje oba strumienie zawsze, a pusty `stderr` jest zwykłym przebiegiem
 * udanym, nie usterką odczytu — dlatego pustka dostaje własne zdanie zamiast
 * znikać razem z elementem.
 */
function strumien(rodzaj: 'stdout' | 'stderr', podpis: string, tresc: string): HTMLElement {
  const blok = document.createElement('section');
  blok.className = 'dt-podglad__strumien';
  blok.dataset['strumien'] = rodzaj;

  const naglowek = document.createElement('h4');
  naglowek.className = 'dt-podglad__podpis';
  naglowek.textContent = podpis;
  blok.append(naglowek);

  if (tresc === '') {
    const pusty = document.createElement('p');
    pusty.className = 'dt-podglad__pusty';
    pusty.textContent = `Ten strumień jest pusty — proces nie wypisał na ${rodzaj} ani jednego znaku.`;
    blok.append(pusty);
    return blok;
  }

  const konsola = document.createElement('pre');
  konsola.className = 'dt-konsola dt-podglad__tresc';
  konsola.setAttribute('aria-label', podpis);
  konsola.textContent = tresc;
  blok.append(konsola);
  return blok;
}
