/**
 * Karta sesji, w której Operator teraz pracuje. Rdzeń prowadzi sesje trwale,
 * a okna modułów stają w ich wnętrzu — klient musi więc wiedzieć, która karta
 * jest bieżąca. Bez tego wskazania każde wejście w moduł zakładałoby kartę
 * nową, a praca Operatora nie miałaby dokąd wracać. Identyfikator sesji stoi
 * w `protokol/sesja.ts` — ten sam, który niosą koperty wychodzące.
 */
import { Command, type Window } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { sesjaKlienta } from '../protokol/sesja.ts';
import { tozsamoscKlienta } from '../protokol/tozsamosc-klienta.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import { odtworzOknaRobocze } from './okna-robocze.ts';

/* Środowisko, przez które Operator wszedł do pracy. Sesja zakładana pierwszą
   wiadomością bierze je stąd, bo karta środowiska liczy po nim swoje sesje. */
let srodowiskoWejscia = '';

/** Zapamiętuje środowisko wejścia Operatora. */
export function wskazSrodowisko(kod: string): void {
  srodowiskoWejscia = kod;
}

/** Karta sesji bieżącej; pustka znaczy, że Operator żadnej jeszcze nie otworzył. */
export function sesjaBiezaca(): string {
  return sesjaKlienta().id();
}

/** Przenosi ognisko na wskazaną kartę sesji i czyni ją bieżącą. */
export async function wskazSesje(kanal: Kanal, idSesji: string): Promise<void> {
  if (idSesji === '') return;
  sesjaKlienta().ustaw(idSesji);
  await wywolaj(kanal, Command.SessionFocus, {
    sessionId: idSesji,
    clientId: tozsamoscKlienta().id,
  });
}

/**
 * Otwiera wskazaną kartę sesji wraz z jej oknami i czyni ją bieżącą.
 * Zwraca okna karty — po nich rozstrzyga się, czy jest do czego wracać.
 */
export async function otworzSesje(kanal: Kanal, idSesji: string): Promise<Window[]> {
  const wynik = await wywolaj(kanal, Command.SessionOpen, { sessionId: idSesji });
  if (!wynik.udany || wynik.wynik === undefined) return [];
  await wskazSesje(kanal, wynik.wynik.session.id);
  return wynik.wynik.windows ?? [];
}

/**
 * Karta sesji, w której ma stanąć okno: bieżąca, a gdy Operator żadnej nie
 * otworzył — założona pod nazwą wejścia, żeby wykaz nie zbierał kart bez nazwy.
 */
export async function zapewnijSesje(
  kanal: Kanal,
  nazwa: string,
  kodSrodowiska = '',
): Promise<string> {
  const biezaca = sesjaKlienta().id();
  if (biezaca !== '') return biezaca;
  const zadanie: { title?: string; environmentCode?: string } = {};
  if (nazwa !== '') zadanie.title = nazwa;
  const srodowisko = kodSrodowiska === '' ? srodowiskoWejscia : kodSrodowiska;
  if (srodowisko !== '') zadanie.environmentCode = srodowisko;
  const wynik = await wywolaj(kanal, Command.SessionCreate, zadanie);
  if (!wynik.udany || wynik.wynik === undefined) return '';
  await wskazSesje(kanal, wynik.wynik.session.id);
  return sesjaKlienta().id();
}

/**
 * Przejmuje przy wejściu na stronę główną stan, który rdzeń prowadzi za
 * Operatora: kartę ogniskowaną, sesje czynne konta i ich żywy odpis. Z sesji
 * i ich okien odtwarza się wykaz okien roboczych — klient nie zapisuje go
 * nigdzie sam, bo kontrakt nie ma na niego pola.
 */
export async function przejmijOgnisko(kanal: Kanal): Promise<void> {
  const wynik = await wywolaj(kanal, Command.HomeEnter, { clientId: tozsamoscKlienta().id });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos('Strona główna', wynik.blad?.message
      ?? 'Rdzeń odmówił wejścia na stronę główną — sesje i okna robocze zostają puste.', 'blad');
    return;
  }
  const ogniskowana = wynik.wynik.focusedSessionId ?? '';
  /* Ognisko idzie przez `session.focus`, nie samym wpisem: koperty wychodzące
     biorą sesję z `protokol/sesja.ts`, a rdzeń potwierdza ją zdarzeniem. */
  if (ogniskowana !== '') await wskazSesje(kanal, ogniskowana);
  const odtworzone = await odtworzOknaRobocze(
    kanal, wynik.wynik.sessions, wynik.wynik.presence ?? [], ogniskowana,
  );
  if (!odtworzone) {
    oglos('Okna robocze', 'Rdzeń nie podał okien komunikacji — przełącznik okien '
      + 'wstaje pusty, choć praca w sesjach mogła zostać.', 'ostrzezenie');
  }
}
