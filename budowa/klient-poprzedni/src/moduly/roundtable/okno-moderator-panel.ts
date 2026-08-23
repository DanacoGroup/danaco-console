import {
  Command,
  ModeratorAction,
  RoundtableTurnStatus,
  type RoundtableParticipant,
} from '../../../../shared/contract';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import {
  pole,
  poleTresci,
  przyciskAkcji as przycisk,
  przyciskBezKomendy,
  wiersz,
} from '../../modele/kontrolki-formularza-braki';
import { powodBrakuCzynnosci, powodBrakuObslugi } from './braki-kontraktu';
import {
  potwierdzKolejnosc,
  potwierdzWyciszenie,
  potwierdzZagadnienie,
  potwierdzZamkniecie,
  utworzKontekstCzynnosci,
} from './obsluga-czynnosci-moderatora';
import { rysujTuraISklad, uczestnicyWKolejnosci } from './sklad-moderatora';
import type { StanDebaty } from './stan-debaty';
import { utworzStanTresci } from './stany-okna';
import { czynnosciModeracji } from './czynnosci-arsenalu';
import { utworzWykazFunkcji, utworzZestawAkcji } from './warstwy-modulu';
import type { ZrodloArsenaluRoundtable } from './zrodlo-arsenalu';
import { ODSTEP_ZEGARA_MS, zdanieZegara } from './zegar-tury';
import type { ZrodloRoundtable } from './zrodlo-roundtable';

/**
 * Moderator Panel — okno zarządcy modułu Roundtable: ukierunkowanie dyskusji,
 * zamknięcie tury, wybór kolejnego zagadnienia i kolejność głosu; aktywne
 * w toku całej sesji Roundtable.
 *
 * Wyliczenie `ModeratorAction` niesie dziś dziewięć wartości; okno wysyła pięć
 * starszych — `Direct`, `CloseTurn`, `NextTopic`, `Mute`, `Unmute`. Czterech
 * nowszych (`SetSpeakingOrder`, `InjectCounterVoice`, `OpenSideThread`,
 * `MergeSideThread`) jeszcze nie wysyła: obsługi nie zbudowano. Kolejność głosu
 * jedzie zatem nadal razem z `action: Direct`, samodzielnie albo obok `message`
 * (`RoundtableModeratorDirectRequest.speakingOrder?: string[]`), choć wartość
 * `SetSpeakingOrder` jest już w wyliczeniu i przeniesienie na nią kolejności
 * usunęłoby niepewność opisaną przy `przesunKolejnosc`.
 *
 * `participants` w odpowiedzi jest polem nieobowiązkowym, więc okno musi umieć
 * jego brak i wtedy nie potwierdza zmiany składu. Kiedy skład przyjdzie, okno
 * wpisuje go do `stan.ustawUczestnikow(...)` — to jedyna droga, którą Model
 * Panels i Consensus Panel dowiadują się o składzie.
 *
 * Skład i kanały idą wyłącznie ze stanu (`stan.uczestnicy()`,
 * `stan.opisKanalu`) — tego samego rejestru, którym jedzie okno rozmowy. Okno
 * nie woła `channel.list` na własną rękę.
 */
export interface OknoModeratorPanel {
  element: HTMLElement;
  odswiez(): void;
  /** Zamyka nasłuch `stan.naZmiane(...)`. */
  zamknij(): void;
}

export function utworzOknoModeratorPanel(
  zrodlo: ZrodloRoundtable,
  arsenal: ZrodloArsenaluRoundtable,
  stan: StanDebaty,
  idOkna: string,
): OknoModeratorPanel {
  stan.ustawOkno(idOkna);

  const rama = utworzRameOkna({
    tytul: 'Moderator Panel',
    rola: 'zarządca',
    kod: 'moderator-panel',
    przeznaczenie:
      'Ukierunkowanie dyskusji, zamknięcie tury, wybór kolejnego zagadnienia i kolejność głosu — aktywny w toku całej sesji Roundtable.',
    przedrostek: 'dr',
  });
  const tresc = utworzStanTresci();
  // Czynności moderacji wołają komendy obszaru wprost: zapis szablonu, wykaz
  // szablonów, wydanie transkryptu, odsłuch i odczyt stanu debaty z rdzenia.
  const powierzchnia = zlozPowierzchnieModeratora(
    rama,
    tresc.element,
    czynnosciModeracji(
      arsenal,
      stan,
      (zdanie, powodzenie) => tresc.potwierdzenie(zdanie, powodzenie),
      () => 'Szablon moderacji tej debaty',
      () => 'markdown',
    ),
  );

  /** Tura, o której mówi zdanie w pasie czynności — po niej poznajemy zmianę. */
  let turaZdania = '';

  /** Odświeżenie samego zegara — bez przerysowania składu. */
  function odswiezZegar(): void {
    powierzchnia.zegar.textContent = zdanieZegara(stan.definicjaTury(), Date.now());
  }

  /**
   * Zdanie o czynności należy do tury, w której padło.
   *
   * Pas czynności trzyma treść do następnego zapisu, więc bez tego odmowa
   * dotycząca jednej tury wisiałaby nad następną. Wraz ze zmianą tury zdanie
   * schodzi (ten sam wzór, co czyszczenie odpowiedzi przy zmianie eksperta
   * w `agents/okno-skills-manager.ts`).
   */
  function rysuj(): void {
    const biezaca = stan.definicjaTury()?.id ?? '';
    if (biezaca !== turaZdania) {
      turaZdania = biezaca;
      tresc.potwierdzenie('', false);
    }
    odswiezZegar();
    rysujTuraISklad(stan, tresc, { przelaczWyciszenie, przesunKolejnosc });
  }

  /**
   * Odmowa czynności idzie pasem czynności, nie stanem błędu okna.
   *
   * `tresc.blad(...)` czyści miejsce treści, więc zabrałby z ekranu turę i cały
   * skład — okno pokazywałoby pustą debatę tam, gdzie debata jest. Stan błędu
   * zostaje przy nieudanym odczycie, po którym treści naprawdę nie ma. Tutaj
   * treść zostaje przerysowana, a powód odmowy stoi obok niej
   * w `data-udane='false'` (arkusz maluje go barwą błędu).
   */
  function odmow(zdanie: string): void {
    rysuj();
    tresc.potwierdzenie(zdanie, false);
  }

  // Obsługa odpowiedzi czynności (zapis tury, warunkowy zapis składu, dobór
  // zdania potwierdzenia) jest wydzielona do `utworzKontekstCzynnosci`
  // — zwarta odpowiedzialność, która nie musi siedzieć w wytwórni.
  const kontekst = utworzKontekstCzynnosci(zrodlo, stan, tresc, rysuj);

  function skieruj(): void {
    const wiadomosc = powierzchnia.poleInterwencji.value.trim();
    if (wiadomosc === '') {
      odmow('Wpisz treść interwencji moderującej, zanim ją wyślesz.');
      return;
    }
    kontekst.obsluz(
      { windowId: idOkna, action: ModeratorAction.Direct, message: wiadomosc },
      {
        czynnosc: 'interwencja moderująca',
        potwierdz: (odpowiedz) => ({
          zdanie:
            `Rdzeń przyjął interwencję do tury #${odpowiedz.turn.index} ` +
            `(${odpowiedz.turn.status}). Wypowiedzi uczestników przychodzą osobno, zdarzeniem.`,
          udane: true,
        }),
      },
    );
    powierzchnia.poleInterwencji.value = '';
  }

  function zamknijTure(): void {
    const idTury = stan.tura();
    if (idTury === '') {
      odmow('Debata nie ma otwartej tury — nie ma czego zamykać.');
      return;
    }
    // Stan tury sprzed wywołania jest jedynym świadkiem tego, czy zamknięcie
    // cokolwiek zmieniło: rdzeń oddaje turę już zamkniętą bez odmowy.
    const przed = stan.definicjaTury();
    const bylaZamknieta = przed !== null && przed.status === RoundtableTurnStatus.Closed;
    kontekst.obsluz(
      { windowId: idOkna, action: ModeratorAction.CloseTurn, turnId: idTury },
      {
        czynnosc: 'zamknięcie tury',
        potwierdz: (odpowiedz) => potwierdzZamkniecie(odpowiedz, bylaZamknieta),
      },
    );
  }

  function nastepneZagadnienie(): void {
    const temat = powierzchnia.poleZagadnienia.value.trim();
    if (temat === '') {
      odmow('Wpisz zagadnienie kolejnej tury, zanim je ustawisz.');
      return;
    }
    const przed = stan.definicjaTury();
    kontekst.obsluz(
      { windowId: idOkna, action: ModeratorAction.NextTopic, topic: temat },
      {
        czynnosc: 'zmiana zagadnienia tury',
        potwierdz: (odpowiedz) => potwierdzZagadnienie(odpowiedz, przed, temat),
      },
    );
    powierzchnia.poleZagadnienia.value = '';
  }

  function przelaczWyciszenie(uczestnik: RoundtableParticipant): void {
    const wyciszamy = uczestnik.muted !== true;
    kontekst.obsluz(
      {
        windowId: idOkna,
        action: wyciszamy ? ModeratorAction.Mute : ModeratorAction.Unmute,
        participantId: uczestnik.id,
      },
      {
        czynnosc: wyciszamy ? 'wyciszenie uczestnika' : 'zdjęcie wyciszenia uczestnika',
        potwierdz: (odpowiedz) => potwierdzWyciszenie(odpowiedz, uczestnik.id, wyciszamy, stan),
      },
    );
  }

  /**
   * Przesunięcie w kolejności głosu.
   *
   * Przyciski skrajnych wierszy są klikalne zawsze, a ruch niewykonalny dostaje
   * zdanie w pasie czynności — obok wykazu, który zostaje na ekranie.
   * Wygaszenie przycisku milczałoby o powodzie.
   */
  function przesunKolejnosc(idUczestnika: string, kierunek: -1 | 1): void {
    const biezaca = uczestnicyWKolejnosci(stan).map((uczestnik) => uczestnik.id);
    const indeks = biezaca.indexOf(idUczestnika);
    if (indeks === -1) {
      odmow('Tego uczestnika nie ma już w składzie znanym oknu — odśwież skład czynnością moderatora.');
      return;
    }
    const nastepny = indeks + kierunek;
    if (nastepny < 0 || nastepny >= biezaca.length) {
      odmow(
        kierunek === -1
          ? 'Ten uczestnik jest pierwszy w kolejności głosu — wyżej nie ma go dokąd przesunąć.'
          : 'Ten uczestnik jest ostatni w kolejności głosu — niżej nie ma go dokąd przesunąć.',
      );
      return;
    }
    const przestawiona = [...biezaca];
    const [wyjety] = przestawiona.splice(indeks, 1);
    przestawiona.splice(nastepny, 0, wyjety!);
    kontekst.obsluz(
      { windowId: idOkna, action: ModeratorAction.Direct, speakingOrder: przestawiona },
      {
        czynnosc: 'zmiana kolejności głosu',
        potwierdz: (odpowiedz) => potwierdzKolejnosc(odpowiedz, przestawiona, stan),
        // Okno wysyła kolejność obok czynności `Direct`, więc po odmowie nie
        // wie, czego rdzeń odmówił: interwencji czy kolejności. Wyliczenie ma
        // dziś wartość `SetSpeakingOrder`, która tę niepewność usuwa — okno
        // jeszcze jej nie wysyła, bo obsługi nie zbudowano.
        przyNiepowodzeniu:
          'Okno wysyła kolejność głosu obok czynności moderatora, więc NIE WIE, czy rdzeń kolejność zapisał. ' +
          'Wyliczenie ModeratorAction ma już osobną wartość setSpeakingOrder; okno jeszcze jej nie wysyła.',
      },
    );
  }

  podepnijAkcjeModeratora(powierzchnia, { skieruj, zamknijTure, nastepneZagadnienie });

  // Jeden punkt wejścia dla pierwszego rysowania jest `odswiez()`, wołany raz
  // przez złożenie modułu — wytwórnia okna nie czyta stanu sama.
  const odsubskrybuj = stan.naZmiane(() => rysuj());

  // Zegar chodzi własnym rytmem, bo mierzy czas, a nie zmianę stanu: bez niego
  // licznik stałby nieruchomo między przyrostami debaty. Odświeża sam napis
  // zegara, nie cały skład — przerysowanie wykazu co sekundę kasowałoby wpisaną
  // treść interwencji i przenosiło ognisko.
  const zegarBiegnie = window.setInterval(odswiezZegar, ODSTEP_ZEGARA_MS);
  odswiezZegar();

  return {
    element: rama.element,
    odswiez: rysuj,
    zamknij() {
      window.clearInterval(zegarBiegnie);
      odsubskrybuj();
    },
  };
}

/** Kontrolki panelu akcji i ciała ramy. */
interface PowierzchniaModeratora {
  poleInterwencji: HTMLTextAreaElement;
  poleZagadnienia: HTMLInputElement;
  wyslijInterwencje: HTMLButtonElement;
  zamknijTure: HTMLButtonElement;
  ustawZagadnienie: HTMLButtonElement;
  /** Napis zegara tury — odświeżany własnym rytmem, poza rysowaniem składu. */
  zegar: HTMLElement;
}

/**
 * Składa pasek akcji ramy z czterema czynnościami prawdziwymi kontraktu
 * (interwencja, zamknięcie tury, zmiana zagadnienia; kolejność głosu i
 * wyciszenie siedzą przy wierszu uczestnika, bo potrzebują jego tożsamości)
 * oraz z pozycjami, których obsługi jeszcze nie zbudowano.
 *
 * Wydzielone z wytwórni okna — fragment czysty, bez domknięcia na
 * stanie modułu.
 */
function zlozAkcjeModeratora(gospodarz: HTMLElement): {
  wyslijInterwencje: HTMLButtonElement;
  zamknijTure: HTMLButtonElement;
  ustawZagadnienie: HTMLButtonElement;
} {
  const wyslijInterwencje = przycisk('Wyślij interwencję', 'dn-btn dn-btn--atrament');
  const zamknijTure = przycisk('Zamknij turę teraz', 'dn-btn');
  const ustawZagadnienie = przycisk('Ustaw zagadnienie', 'dn-btn');

  gospodarz.append(
    wyslijInterwencje,
    zamknijTure,
    ustawZagadnienie,
    przyciskBezKomendy(
      '+ Nowy wątek poboczny',
      powodBrakuObslugi(
        Command.RoundtableModeratorDirect,
        'Wątek poboczny otwierają i scalają wartości openSideThread i mergeSideThread wyliczenia ModeratorAction, a turę nadrzędną wskazuje pole parentTurnId tury. Okno wysyła dziś pięć wartości starszych i tych dwóch jeszcze nie wysyła.',
      ),
    ),
    przyciskBezKomendy(
      'Zapisz jako szablon',
      powodBrakuObslugi(
        Command.RoundtableModerationTemplateSave,
        'Format, liczbę tur i kolejność głosu zapisuje jako nazwany szablon ta komenda, a szablony zapisane oddaje roundtable.moderation.template.list.',
      ),
    ),
    przyciskBezKomendy(
      'Liczba tur',
      powodBrakuCzynnosci(
        'turnLimit jest polem żądania roundtable.debate.start, nie moderator.direct — to pole innego okna (Debate Panel).',
      ),
    ),
  );
  return { wyslijInterwencje, zamknijTure, ustawZagadnienie };
}

/** Zestaw akcji warstwy trzeciej — czynności moderacji jeszcze niezbudowane. */
function zlozZestawModeratora(czynnosci: HTMLButtonElement[]): HTMLElement {
  return utworzZestawAkcji('Zestaw operacji moderacyjnych', czynnosci);
}

/** Składa formularz interwencji i zagadnienia oraz ciało ramy. */
function zlozPowierzchnieModeratora(
  rama: { akcje: HTMLElement; narzedzia: HTMLElement; cialo: HTMLElement },
  stanTresci: HTMLElement,
  czynnosci: HTMLButtonElement[],
): PowierzchniaModeratora {
  const { wyslijInterwencje, zamknijTure, ustawZagadnienie } = zlozAkcjeModeratora(rama.akcje);
  const poleInterwencji = poleTresci(
    'Treść interwencji moderującej',
    3,
    'np. Wróćmy do pytania głównego…',
  );
  const poleZagadnienia = pole('Zagadnienie kolejnej tury', 'np. Skutki uboczne wariantu B');
  const zegar = document.createElement('p');
  zegar.className = 'dr-zegar';
  zegar.setAttribute('role', 'status');

  rama.narzedzia.append(zegar);
  rama.cialo.append(
    wiersz('Interwencja moderująca', poleInterwencji, { klasa: 'dr-wiersz' }),
    wiersz('Kolejne zagadnienie', poleZagadnienia, { klasa: 'dr-wiersz' }),
    zlozZestawModeratora(czynnosci),
    stanTresci,
    utworzWykazFunkcji('moderator-panel'),
  );
  return {
    poleInterwencji,
    poleZagadnienia,
    wyslijInterwencje,
    zamknijTure,
    ustawZagadnienie,
    zegar,
  };
}

/** Podpina pasek akcji do czynności okna. */
function podepnijAkcjeModeratora(
  powierzchnia: PowierzchniaModeratora,
  obsluga: { skieruj: () => void; zamknijTure: () => void; nastepneZagadnienie: () => void },
): void {
  powierzchnia.wyslijInterwencje.addEventListener('click', obsluga.skieruj);
  powierzchnia.zamknijTure.addEventListener('click', obsluga.zamknijTure);
  powierzchnia.ustawZagadnienie.addEventListener('click', obsluga.nastepneZagadnienie);
}
