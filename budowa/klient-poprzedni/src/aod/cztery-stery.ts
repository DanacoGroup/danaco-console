import type { Queue } from '../../../shared/contract';
import { opisOdmowyAod } from './odmowy-aod';
import { utworzAkapit, utworzPodtytul } from './pola-wykazu';
import type { DecyzjaCzekajaca } from './rozpoznanie-decyzji';
import type { ZrodloDecyzji } from './zrodlo-decyzji';

/**
 * Cztery stery decyzji w jednym pasie przy wpisie: zatwierdzenie kroku,
 * wstrzymanie, konfiguracja Koordynatora, przejęcie bezpośredniego sterowania.
 *
 * Ster ma przycisk tylko wtedy, gdy stoi za nim komenda kontraktu; inaczej
 * wyświetla zdanie nazywające granicę. Zatwierdzenia kroku kontrakt nie zna,
 * a przestawienie ogniska wymaga identyfikatora klienta z powitania połączenia,
 * którego okno nakładki nie otrzymuje. Żaden przycisk nie pyta o potwierdzenie
 * i żaden nie jest wyszarzany — rozstrzyga rdzeń, nakładka pokazuje odpowiedź.
 */

/** Tożsamość klienta potrzebna wyłącznie sterowi przejęcia. */
export interface TozsamoscDlaOgniska {
  /** Identyfikator klienta z powitania połączenia; nie wolno nadać go po raz drugi. */
  id: string;
}

export interface OpisCzterechSterow {
  zrodlo: ZrodloDecyzji;
  /** Krótkie potwierdzenie czynności na pasku okna. */
  zamelduj(zdanie: string, udane: boolean): void;
  /** Ponawia odczyt kolejki po czynności zmieniającej stan rdzenia. */
  poCzynnosci(): void;
  /** Otwiera Okno Konfiguracji klienta — ster 3 kończy się nim, nie samym `opened`. */
  otworzOknoKonfiguracji(): void;
  /** Tożsamość klienta z powitania; brak znaczy „ogniska nie ma czym przestawić". */
  klient?: TozsamoscDlaOgniska;
}

/** Kontekst jednej decyzji — to, co stery mają czym obsłużyć. */
export interface KontekstDecyzji {
  decyzja: DecyzjaCzekajaca;
  /** Kolejka dopasowana do decyzji; brak, gdy `queue.list` żadnej nie wskazał. */
  kolejka?: Queue;
}

/** Buduje pas czterech sterów dla jednej decyzji. */
export function utworzCzteryStery(
  opis: OpisCzterechSterow,
  kontekst: KontekstDecyzji,
): HTMLElement {
  const pas = document.createElement('div');
  pas.className = 'ao-stery';
  pas.append(
    sterZatwierdzenia(),
    sterWstrzymania(opis, kontekst),
    sterKonfiguracji(opis, kontekst),
    sterPrzejecia(opis, kontekst),
  );
  return pas;
}

/** Wspólna obudowa jednego steru: nazwa drogi, przyciski, granice. */
function utworzSter(nazwa: string): { element: HTMLElement; tresc: HTMLElement } {
  const element = document.createElement('div');
  element.className = 'ao-ster';

  const tresc = document.createElement('div');
  tresc.className = 'ao-ster__tresc';

  element.append(utworzPodtytul(nazwa), tresc);
  return { element, tresc };
}

/** Przycisk czynny — bez wyszarzania i bez pytania o potwierdzenie. */
function utworzPrzycisk(napis: string, czynnosc: () => void): HTMLButtonElement {
  const przycisk = document.createElement('button');
  przycisk.type = 'button';
  przycisk.className = 'dn-btn dn-btn--zarys ao-ster__przycisk';
  przycisk.textContent = napis;
  przycisk.addEventListener('click', czynnosc);
  return przycisk;
}

/** Ster zatwierdzenia kroku. Kontrakt nie zna tej komendy, więc ster nazywa brak. */
function sterZatwierdzenia(): HTMLElement {
  const { element, tresc } = utworzSter('1. Zatwierdź krok');
  tresc.append(
    utworzAkapit(
      'ao-granica',
      'Rdzeń nie ma dziś komendy zatwierdzającej krok: window.action z akcją „approve" ' +
        'odpowiada not_found, a kontrakt nie niesie ani typu pozycji kolejki, ani komendy ' +
        'czytającej pozycje — identyfikator pozycji wychodzi wyłącznie z window.handoff ' +
        'i nie da się go odczytać ponownie. Przycisk bez komendy byłby atrapą, więc go tu nie ma.',
    ),
  );
  return element;
}

/** Ster wstrzymania: kolejka, tura okna, tury sesji. */
function sterWstrzymania(opis: OpisCzterechSterow, kontekst: KontekstDecyzji): HTMLElement {
  const { element, tresc } = utworzSter('2. Wstrzymaj');
  const { decyzja, kolejka } = kontekst;

  if (kolejka !== undefined) {
    tresc.append(
      utworzPrzycisk(`Wstrzymaj kolejkę ${kolejka.id}`, () => {
        void wstrzymajKolejke(opis, kolejka.id);
      }),
    );
  } else {
    tresc.append(
      utworzAkapit(
        'ao-granica',
        'Do tej decyzji nie da się dopasować kolejki — queue.list żadnej nie wiąże ' +
          'z tym oknem ani z tą sesją. Wstrzymanie idzie wtedy przez okno albo przez sesję.',
      ),
    );
  }

  if (decyzja.idOkna !== undefined && decyzja.idOkna !== '') {
    const idOkna = decyzja.idOkna;
    tresc.append(
      utworzPrzycisk('Zatrzymaj turę okna', () => {
        void zatrzymajTureOkna(opis, idOkna);
      }),
    );
  }

  if (decyzja.idSesji !== undefined && decyzja.idSesji !== '') {
    const idSesji = decyzja.idSesji;
    tresc.append(
      utworzPrzycisk('Zatrzymaj tury sesji', () => {
        void zatrzymajTurySesji(opis, idSesji);
      }),
    );
  }

  tresc.append(
    utworzAkapit(
      'ao-granica',
      'Samego procesu telemetrii wstrzymać się nie da — rdzeń nie ma takiej drogi. ' +
        'Wariant mobilny też nie: mobile.process.control odpowiada „warstwa mobilna: ' +
        'port nawigacji nie niesie mobilnego centrum dowodzenia".',
    ),
  );

  return element;
}

/** Ster konfiguracji Koordynatora; kończy się otwartym oknem konfiguracji klienta. */
function sterKonfiguracji(opis: OpisCzterechSterow, kontekst: KontekstDecyzji): HTMLElement {
  const { element, tresc } = utworzSter('3. Zmodyfikuj konfigurację Koordynatora');

  // Okno Koordynatora, gdy bieg je wskazał; inaczej okno, w którym decyzja stoi.
  const idOkna = kontekst.decyzja.idOknaKoordynatora ?? kontekst.decyzja.idOkna;

  if (idOkna === undefined || idOkna === '') {
    tresc.append(
      utworzAkapit(
        'ao-granica',
        'Ta decyzja nie wskazuje żadnego okna, więc nie ma czego skonfigurować na ' +
          'zasięgu okna. Konfigurację szerszą otwiera listwa Ustawienia.',
      ),
    );
    return element;
  }

  const okno = idOkna;
  tresc.append(
    utworzPrzycisk(`Otwórz konfigurację okna ${okno}`, () => {
      void otworzKonfiguracje(opis, okno);
    }),
    utworzAkapit(
      'ao-granica',
      kontekst.decyzja.idOknaKoordynatora === undefined
        ? 'Bieg naprawczy nie wskazał okna koordynatora, więc konfiguracja otwiera się na ' +
            'oknie, w którym decyzja stoi — poziom window jest najwęższy i wygrywa z każdym szerszym.'
        : 'Okno wskazane przez bieg naprawczy jako koordynator; poziom window jest najwęższy ' +
            'i wygrywa z każdym szerszym zapisem.',
    ),
  );
  return element;
}

/** Ster przejęcia sterowania. Kontrakt nie ma jednej komendy przejęcia, tylko jego części. */
function sterPrzejecia(opis: OpisCzterechSterow, kontekst: KontekstDecyzji): HTMLElement {
  const { element, tresc } = utworzSter('4. Przejmij bezpośrednie sterowanie');
  const { decyzja } = kontekst;
  const idOkna = decyzja.idOkna;

  if (idOkna !== undefined && idOkna !== '') {
    tresc.append(
      utworzPrzycisk('Wyjmij okno z pętli koordynatora', () => {
        void wyjmijZPetli(opis, idOkna);
      }),
    );

    if (opis.klient !== undefined && decyzja.idSesji !== undefined && decyzja.idSesji !== '') {
      const idSesji = decyzja.idSesji;
      const idKlienta = opis.klient.id;
      tresc.append(
        utworzPrzycisk('Przestaw ognisko na to okno', () => {
          void przeniesOgnisko(opis, { sessionId: idSesji, clientId: idKlienta, windowId: idOkna });
        }),
      );
    } else {
      tresc.append(
        utworzAkapit(
          'ao-granica',
          'Ogniska nakładka dziś nie przestawia: session.focus żąda identyfikatora klienta ' +
            'z powitania połączenia, a okno nakładki dostaje sam kanał. Nadanie identyfikatora ' +
            'po raz drugi wskazałoby klienta, którego nie ma — komenda odpowiedziałaby pomyślnie, ' +
            'a ekran stałby w miejscu.',
        ),
      );
    }
  } else {
    tresc.append(
      utworzAkapit(
        'ao-granica',
        'Ta decyzja nie wskazuje okna, więc nie ma czego wyjąć z pętli ani na co przestawić ogniska.',
      ),
    );
  }

  tresc.append(
    utworzAkapit(
      'ao-granica',
      'Kontrakt nie ma jednej komendy pełnego przejęcia sterowania. Nakładka wyjmuje okno ' +
        'z pętli koordynatora i — gdy zna klienta — przestawia ognisko.',
    ),
  );

  return element;
}

// Czynności sterów. Każda melduje odpowiedź rdzenia, nie powtórzenie zamiaru operatora.

async function wstrzymajKolejke(opis: OpisCzterechSterow, queueId: string): Promise<void> {
  const wynik = await opis.zrodlo.wstrzymajKolejke(queueId);
  if (!wynik.udany || wynik.wynik === undefined) {
    opis.zamelduj(opisOdmowyAod('Wstrzymanie kolejki', 'wstrzymanie', wynik.blad), false);
    return;
  }
  opis.zamelduj(
    `Kolejka ${wynik.wynik.queue.id} po wstrzymaniu jest w stanie ${wynik.wynik.queue.status}.`,
    true,
  );
  opis.poCzynnosci();
}

async function zatrzymajTureOkna(opis: OpisCzterechSterow, windowId: string): Promise<void> {
  const wynik = await opis.zrodlo.zatrzymajTureOkna(windowId);
  if (!wynik.udany || wynik.wynik === undefined) {
    opis.zamelduj(opisOdmowyAod('Zatrzymanie tury okna', 'wstrzymanie', wynik.blad), false);
    return;
  }
  // `stopped:false` nie jest błędem — rdzeń mówi, że nie było czego zatrzymywać.
  opis.zamelduj(
    wynik.wynik.stopped
      ? `Tura okna ${windowId} zatrzymana (wiadomość ${wynik.wynik.messageId}).`
      : `W oknie ${windowId} nie biegła żadna tura — rdzeń nie miał czego zatrzymać.`,
    true,
  );
  opis.poCzynnosci();
}

async function zatrzymajTurySesji(opis: OpisCzterechSterow, sessionId: string): Promise<void> {
  const wynik = await opis.zrodlo.zatrzymajTurySesji(sessionId);
  if (!wynik.udany || wynik.wynik === undefined) {
    opis.zamelduj(opisOdmowyAod('Zatrzymanie tur sesji', 'wstrzymanie', wynik.blad), false);
    return;
  }
  const okna = wynik.wynik.stoppedWindowIds;
  opis.zamelduj(
    okna.length === 0
      ? `W sesji ${sessionId} nie biegła żadna tura — rdzeń nie miał czego zatrzymać.`
      : `Zatrzymano tury w oknach: ${okna.join(', ')}.`,
    true,
  );
  opis.poCzynnosci();
}

async function otworzKonfiguracje(opis: OpisCzterechSterow, windowId: string): Promise<void> {
  const wynik = await opis.zrodlo.otworzKonfiguracjeOkna(windowId);
  if (!wynik.udany || wynik.wynik === undefined) {
    opis.zamelduj(opisOdmowyAod('Otwarcie konfiguracji okna', 'konfiguracja', wynik.blad), false);
    return;
  }
  if (!wynik.wynik.opened) {
    opis.zamelduj(
      `Rdzeń nie otworzył konfiguracji okna ${windowId} — odpowiedział opened:false.`,
      false,
    );
    return;
  }
  // Rdzeń potwierdził zasięg; teraz Operator ma zobaczyć konfigurację na ekranie.
  opis.otworzOknoKonfiguracji();
  opis.zamelduj(
    `Konfiguracja otwarta na zasięgu ${wynik.wynik.scope ?? 'window'} dla okna ${windowId}.`,
    true,
  );
}

async function wyjmijZPetli(opis: OpisCzterechSterow, windowId: string): Promise<void> {
  const wynik = await opis.zrodlo.wyjmijZPetli(windowId);
  if (!wynik.udany || wynik.wynik === undefined) {
    opis.zamelduj(opisOdmowyAod('Wyjęcie okna z pętli', 'przejecie', wynik.blad), false);
    return;
  }
  opis.zamelduj(
    `Okno ${wynik.wynik.windowId} ma teraz rolę ${wynik.wynik.role} — jest poza pętlą koordynatora.`,
    true,
  );
  opis.poCzynnosci();
}

async function przeniesOgnisko(
  opis: OpisCzterechSterow,
  zadanie: { sessionId: string; clientId: string; windowId: string },
): Promise<void> {
  const wynik = await opis.zrodlo.przeniesOgnisko(zadanie);
  if (!wynik.udany || wynik.wynik === undefined) {
    opis.zamelduj(opisOdmowyAod('Przeniesienie ogniska', 'przejecie', wynik.blad), false);
    return;
  }
  opis.zamelduj(
    `Ognisko na sesji ${wynik.wynik.sessionId}, okno ${wynik.wynik.windowId ?? zadanie.windowId}.`,
    true,
  );
  opis.poCzynnosci();
}
