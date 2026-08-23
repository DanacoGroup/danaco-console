import {
  Command,
  ProcessInitiator,
  TerminalProcessStatus,
  type TerminalProcess,
} from '../../../../shared/contract';
import { pozycjaWykazu, przyciskAkcji } from '../../modele/kontrolki-formularza';

/**
 * Jedna pozycja wykazu Process Monitora wraz z jej czynnościami.
 *
 * Inicjator stoi w opisie pozycji, nie w podpowiedzi: różnica między procesem
 * uruchomionym przez Operatora a uruchomionym przez model jest różnicą
 * odpowiedzialności i ma być widoczna bez najeżdżania kursorem.
 *
 * Zakończenie ma dwa przyciski, bo ma dwa skutki: sygnał łagodny kończy sam
 * proces polecenia, wymuszony obejmuje całe drzewo potomstwa. Jeden przycisk
 * z przełącznikiem ukryłby tę różnicę.
 *
 * Żaden przycisk pozycji nie jest wygaszany. Stan procesu znany oknu jest kopią —
 * przychodzi wykazem `terminal.process.list` i zdarzeniami
 * `terminal.process.changed` — a kopia bywa nieświeża: zgubione zdarzenie albo
 * rozłączenie zostawia `finished` przy procesie, który biegnie, i blokada odbiera
 * wtedy jedyną drogę zakończenia. Wygaszenie w drugą stronę też nic nie daje:
 * zakończenie procesu już zakończonego nie ma skutku, a odpowiedź na taką próbę
 * należy do rdzenia, nie do okna. Stan jest więc informacją — stoi w opisie
 * pozycji i w podpowiedzi przycisku — a nie bramą.
 */
export interface CzynnosciProcesu {
  przypiety: boolean;
  /** Czy podgląd wyjścia tego procesu jest w tej chwili rozwinięty. */
  podgladOtwarty: boolean;
  zakoncz(wymuszony: boolean): void;
  /**
   * Wstrzymanie procesu albo jego wznowienie (`terminal.process.suspend`).
   *
   * Czynność jest odrębna od zakończenia, bo służy czemu innemu: oddaje procesor
   * BEZ utraty wykonanej pracy. Proces wstrzymany zostaje w stanie `running`,
   * bo taka jest prawda — ma PID, pamięć i otwarte pliki — a kontrakt stanu
   * „wstrzymany” nie ma.
   */
  wstrzymaj(wznowienie: boolean): void;
  uruchomPonownie(): void;
  przypnij(): void;
  eksportuj(): void;
  doKarty(): void;
  /** Odczyt `terminal.output.read` albo zwinięcie podglądu, gdy już stoi. */
  pokazWyjscie(): void;
}

export function pozycjaProcesu(proces: TerminalProcess, czynnosci: CzynnosciProcesu): HTMLElement {
  const czynny = proces.status === TerminalProcessStatus.Running;
  const { element, akcje } = pozycjaWykazu(proces.command, opisProcesu(proces), 'dt');
  element.dataset['stan'] = proces.status;
  element.dataset['inicjator'] = proces.initiator;
  element.dataset['przypiety'] = String(czynnosci.przypiety);

  const lagodny = przyciskAkcji('Zakończ', 'dn-btn dn-btn--zarys');
  lagodny.title = czynny
    ? 'Kończy sam proces polecenia; potomstwo zostaje.'
    : `Stan znany oknu to ${proces.status} — o tym, czy proces naprawdę biegnie, rozstrzygnie rdzeń.`;
  lagodny.addEventListener('click', () => czynnosci.zakoncz(false));

  const wymuszony = przyciskAkcji('Zakończ wymuszenie', 'dn-btn dn-btn--zarys');
  wymuszony.title = czynny
    ? 'Kończy proces wraz z całym drzewem potomstwa.'
    : `Stan znany oknu to ${proces.status} — o tym, czy proces naprawdę biegnie, rozstrzygnie rdzeń.`;
  wymuszony.addEventListener('click', () => czynnosci.zakoncz(true));

  const wstrzymaj = przyciskAkcji('Wstrzymaj', 'dn-btn dn-btn--zarys');
  wstrzymaj.title =
    'Zatrzymuje proces wraz z jego potomstwem, nie tracąc wykonanej pracy. Systemy uniksowe robią to ' +
    'sygnałem; Windows takiej czynności dla obcego procesu nie ma i odpowiedź mówi to wprost.';
  wstrzymaj.addEventListener('click', () => czynnosci.wstrzymaj(false));

  const wznow = przyciskAkcji('Wznów', 'dn-btn dn-btn--zarys');
  wznow.title = 'Podejmuje pracę procesu wstrzymanego w miejscu, w którym stanęła.';
  wznow.addEventListener('click', () => czynnosci.wstrzymaj(true));

  const ponownie = przyciskAkcji('Uruchom ponownie');
  ponownie.addEventListener('click', () => czynnosci.uruchomPonownie());

  const przypnij = przyciskAkcji(czynnosci.przypiety ? 'Odepnij' : 'Przypnij');
  przypnij.addEventListener('click', () => czynnosci.przypnij());

  const eksport = przyciskAkcji('Eksportuj wpis');
  eksport.addEventListener('click', () => czynnosci.eksportuj());

  const doKarty = przyciskAkcji('→ Karta źródłowa');
  doKarty.title =
    proces.sessionId === undefined || proces.sessionId === ''
      ? 'Proces nie wskazuje karty źródłowej.'
      : `Przestawia okno wiodące na kartę ${proces.sessionId}.`;
  doKarty.addEventListener('click', () => czynnosci.doKarty());

  // Podgląd wyjścia — jedyna droga do treści, którą rdzeń pamięta, a bufor okna
  // już nie (bufor umiera z połączeniem, rejestr rdzenia żyje cały jego bieg).
  // Bez `disabled` jak reszta: stan procesu znany oknu jest kopią i nie ma prawa
  // odbierać jedynej drogi do wyjścia.
  const wyjscie = przyciskAkcji(
    czynnosci.podgladOtwarty ? 'Ukryj wyjście' : 'Pokaż wyjście',
    'dn-btn dn-btn--atrament',
  );
  wyjscie.title = czynnosci.podgladOtwarty
    ? 'Zwija podgląd. Nic nie jedzie do rdzenia — treść znika tylko z widoku.'
    : czynny
      ? `Pyta rdzeń o wyjście tego procesu (${Command.TerminalOutputRead}) i chwilę czeka na jego domknięcie. Proces wciąż biegnący odda wyjście dotychczasowe.`
      : `Pyta rdzeń o wyjście tego procesu (${Command.TerminalOutputRead}) — także po rozłączeniu, gdy bufor okna jest już pusty.`;
  wyjscie.addEventListener('click', () => czynnosci.pokazWyjscie());

  akcje.append(wyjscie, wstrzymaj, wznow, lagodny, wymuszony, ponownie, przypnij, eksport, doKarty);
  return element;
}

/** Opis procesu: PID, rodzic, inicjator, stan, kod wyjścia i czasy. */
function opisProcesu(proces: TerminalProcess): string {
  const czesci = [
    `PID ${proces.pid ?? '—'}`,
    `rodzic ${proces.parentPid ?? '—'}`,
    proces.initiator === ProcessInitiator.Model ? 'inicjator: AI' : 'inicjator: Operator',
    `stan: ${proces.status}`,
    `kod wyjścia: ${proces.exitCode ?? '—'}`,
    `start: ${chwila(proces.startedAt)}`,
  ];
  if (proces.finishedAt !== undefined) czesci.push(`koniec: ${chwila(proces.finishedAt)}`);
  // Zużycia procesora i pamięci nie pokazujemy: rdzeń ich nie mierzy, a wpisanie
  // tu zera znaczyłoby „nic nie zużywa”, co byłoby nieprawdą.
  return czesci.join(' · ');
}

function chwila(znacznik: number): string {
  if (znacznik <= 0) return '—';
  return new Date(znacznik).toISOString().replace('T', ' ').slice(0, 19);
}
