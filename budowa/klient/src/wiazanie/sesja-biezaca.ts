/**
 * Sesja okna roboczego, w którym Operator teraz pracuje. Okno robocze niesie
 * najwyżej jedną sesję i z niego bierze się identyfikator nadawany kopertom
 * (`protokol/sesja.ts`): przełączenie okna przestawia sesję klienta, a okno
 * bez sesji zostawia ją pustą, aż pierwsza wiadomość ją założy.
 */
import { Command, type Session, type Window } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { sesjaKlienta } from '../protokol/sesja.ts';
import { tozsamoscKlienta } from '../protokol/tozsamosc-klienta.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import {
  odtworzOknaRobocze,
  oknoBiezace,
  oknoKarty,
  przypiszSesjeOkna,
  zdejmijSesjeOkien,
  type OknoRobocze,
} from './okna-robocze.ts';

/* Środowisko, przez które Operator wszedł do pracy. Sesja zakładana pierwszą
   wiadomością bierze je stąd, bo karta środowiska liczy po nim swoje sesje. */
let srodowiskoWejscia = '';

/** Zapamiętuje środowisko wejścia Operatora. */
export function wskazSrodowisko(kod: string): void {
  srodowiskoWejscia = kod;
}

/** Sesja okna roboczego bieżącego; pustka znaczy okno, w którym sesja jeszcze nie powstała. */
export function sesjaBiezaca(): string {
  return sesjaKlienta().id();
}

/** Przenosi ognisko rdzenia na wskazaną sesję i czyni ją sesją klienta. */
export async function wskazSesje(kanal: Kanal, idSesji: string): Promise<void> {
  if (idSesji === '') return;
  sesjaKlienta().ustaw(idSesji);
  await wywolaj(kanal, Command.SessionFocus, {
    sessionId: idSesji,
    clientId: tozsamoscKlienta().id,
  });
}

/**
 * Uzgadnia sesję klienta z oknem roboczym, na które Operator przeszedł: sesja
 * okna idzie w ognisko rdzenia, a okno bez sesji zeruje sesję klienta —
 * kontrakt nie ma komendy zdejmującej ognisko, więc pustka stoi tylko po
 * stronie klienta. Sesja już stojąca w kliencie nie idzie do rdzenia drugi raz.
 */
export async function uzgodnijSesjeOkna(kanal: Kanal, okno: OknoRobocze = oknoBiezace()): Promise<void> {
  if (okno.idSesji === '') {
    sesjaKlienta().ustaw('');
    return;
  }
  if (sesjaKlienta().id() === okno.idSesji) return;
  await wskazSesje(kanal, okno.idSesji);
}

/**
 * Otwiera wskazaną sesję wraz z jej oknami i czyni ją sesją klienta; wołający
 * stawia dla niej okno robocze. Pustka znaczy odmowę rdzenia.
 */
export async function otworzSesje(
  kanal: Kanal,
  idSesji: string,
): Promise<{ sesja: Session; okna: Window[] } | null> {
  const wynik = await wywolaj(kanal, Command.SessionOpen, { sessionId: idSesji });
  if (!wynik.udany || wynik.wynik === undefined) return null;
  await wskazSesje(kanal, wynik.wynik.session.id);
  return { sesja: wynik.wynik.session, okna: wynik.wynik.windows ?? [] };
}

/**
 * Sesja okna roboczego, w którym stoi wskazana karta: stojąca, a gdy okno
 * sesji jeszcze nie ma — założona pod nazwą wejścia, żeby wykaz nie zbierał
 * sesji bez nazwy. Pustka znaczy kartę już zdjętą albo odmowę rdzenia. Sesja
 * idzie po karcie, nie po oknie bieżącym: Operator mógł przejść na inne okno,
 * zanim rdzeń odpowiedział.
 */
export async function zapewnijSesje(
  kanal: Kanal,
  idKarty: string,
  nazwa: string,
): Promise<string> {
  const okno = oknoKarty(idKarty);
  if (okno === undefined) return '';
  if (okno.idSesji !== '') {
    await uzgodnijSesjeOkna(kanal, okno);
    return okno.idSesji;
  }
  const zadanie: { title?: string; environmentCode?: string } = {};
  if (nazwa !== '') zadanie.title = nazwa;
  if (srodowiskoWejscia !== '') zadanie.environmentCode = srodowiskoWejscia;
  const wynik = await wywolaj(kanal, Command.SessionCreate, zadanie);
  if (!wynik.udany || wynik.wynik === undefined) return '';
  const idSesji = wynik.wynik.session.id;
  /* Okno mogło dostać sesję z innej karty, gdy rdzeń odpowiadał; wtedy sesja
     założona tu zostaje osierocona w rejestrze i praca idzie w tę stojącą. */
  if (!przypiszSesjeOkna(okno.id, idSesji, nazwa)) return okno.idSesji;
  await wskazSesje(kanal, idSesji);
  return idSesji;
}

/** Zdejmuje sesję usuniętą w rdzeniu z okien roboczych; sesja klienta, gdy była tą sesją, pustoszeje. */
export function zapomnijSesje(idSesji: string): void {
  zdejmijSesjeOkien(idSesji);
  if (sesjaKlienta().id() === idSesji) sesjaKlienta().ustaw('');
}

/**
 * Przejmuje przy wejściu na stronę główną stan, który rdzeń prowadzi za
 * Operatora: sesję ogniskowaną, sesje czynne konta i ich żywy odpis. Z sesji
 * i ich okien odtwarza się wykaz okien roboczych — klient nie zapisuje go
 * nigdzie sam, bo kontrakt nie ma na niego pola. Po odtworzeniu sesja klienta
 * idzie za oknem bieżącym.
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
  await uzgodnijSesjeOkna(kanal);
}
