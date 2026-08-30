/**
 * Karta sesji, w której Operator teraz pracuje. Rdzeń prowadzi sesje trwale,
 * a okna modułów stają w ich wnętrzu — klient musi więc wiedzieć, która karta
 * jest bieżąca. Bez tego wskazania każde wejście w moduł zakładałoby kartę
 * nową, a praca Operatora nie miałaby dokąd wracać.
 */

import { Command, type Window } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { tozsamoscKlienta } from '../protokol/tozsamosc-klienta.ts';
import { wywolaj } from '../protokol/wywolanie.ts';

let biezaca = '';

/** Karta sesji bieżącej; pustka znaczy, że Operator żadnej jeszcze nie otworzył. */
export function sesjaBiezaca(): string {
  return biezaca;
}

/** Przenosi ognisko na wskazaną kartę sesji i czyni ją bieżącą. */
export async function wskazSesje(kanal: Kanal, idSesji: string): Promise<void> {
  if (idSesji === '') return;
  biezaca = idSesji;
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
export async function zapewnijSesje(kanal: Kanal, nazwa: string): Promise<string> {
  if (biezaca !== '') return biezaca;
  const wynik = await wywolaj(kanal, Command.SessionCreate, nazwa === '' ? {} : { title: nazwa });
  if (!wynik.udany || wynik.wynik === undefined) return '';
  await wskazSesje(kanal, wynik.wynik.session.id);
  return biezaca;
}

/** Zakłada nową kartę sesji i czyni ją bieżącą; zwraca jej identyfikator. */
export async function zalozSesje(kanal: Kanal, nazwa: string): Promise<string> {
  const wynik = await wywolaj(kanal, Command.SessionCreate, nazwa === '' ? {} : { title: nazwa });
  if (!wynik.udany || wynik.wynik === undefined) return '';
  await wskazSesje(kanal, wynik.wynik.session.id);
  return biezaca;
}

/**
 * Przejmuje kartę ogniskowaną przez rdzeń przy wejściu na stronę główną.
 * Rdzeń pamięta, gdzie Operator skończył — klient wraca tam, a nie do pustki.
 */
export async function przejmijOgnisko(kanal: Kanal): Promise<void> {
  const wynik = await wywolaj(kanal, Command.HomeEnter, { clientId: tozsamoscKlienta().id });
  if (!wynik.udany || wynik.wynik === undefined) return;
  const ogniskowana = wynik.wynik.focusedSessionId ?? '';
  if (ogniskowana !== '') biezaca = ogniskowana;
}
