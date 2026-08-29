/**
 * Strefa 4 — okno komunikacji. Kształt `.sta-czaty > .sta-kom` ze źródła
 * (`design/05-okna/moduly/studio.html`): belka, nagłówek parametrów, pas
 * kontekstu, historia wpisów, monitor i dół z polem polecenia.
 *
 * Rozmowę prowadzi rdzeń: `message.list` wczytuje wpisy okna, `message.send`
 * wysyła polecenie Operatora, a zdarzenie `message.changed` dokłada odpowiedzi
 * w miarę, jak powstają. Żaden wpis nie powstaje po stronie okna.
 *
 * Nagłówek parametrów i pas kontekstu stoją nazwanym stanem pustym: środowisko,
 * model i wysiłek rozstrzyga się dziś poza oknem komunikacji, a przypięte
 * źródła kontekstu nie mają w kontrakcie komendy, która by je oddawała.
 */

import {
  ChangeKind,
  Command,
  EventType,
  MessageRole,
  MessageStatus,
  type ErrorInfo,
  type Message,
} from '../../../../../shared/contract.ts';
import type { Kanal } from '../../../protokol/kanal.ts';
import { wywolaj } from '../../../protokol/wywolanie.ts';
import { ikony, type NazwaZnaku } from '../ikony.ts';
import { el, tekst, zeZnacznika, type Dziecko } from '../narzedzia.ts';
import { opisOdmowy, zaloguj } from '../odmowa.ts';

export interface ZaleznosciCzatu {
  /** Kanał, którym okno woła komendy rdzenia. */
  kanal: Kanal;
  /** Okno komunikacji — odczytywane przy wywołaniu, bo powstaje po montażu bryły. */
  idOkna(): string | null;
  /** Sesja, w której rozmowa się toczy. */
  idSesji(): string | null;
}

export interface ZamontowaneOknoCzatu {
  wezel: HTMLElement;
  zdejmij(): void;
}

/** Rysunek medalionu wpisu po roli nadawcy; rola spoza wykazu bierze rysunek warstwy systemowej. */
const ZNAK_ROLI: Record<string, NazwaZnaku> = {
  [MessageRole.User]: 'olowek',
  [MessageRole.Assistant]: 'dymek',
  [MessageRole.System]: 'plan',
  [MessageRole.Tool]: 'klucz',
};

/** Odmiana wpisu po roli nadawcy — trzy warianty, tak jak w źródle kształtu. */
const WARIANT_ROLI: Record<string, string> = {
  [MessageRole.User]: 'sta-wpis--czlowiek',
  [MessageRole.Assistant]: 'sta-wpis--inteligencja',
  [MessageRole.System]: 'sta-wpis--system',
  [MessageRole.Tool]: 'sta-wpis--system',
};

function znak(rysunek: NazwaZnaku): SVGElement {
  const wezel = zeZnacznika(ikony[rysunek]);
  wezel.setAttribute('aria-hidden', 'true');
  return wezel;
}

/** Godzina wpisu w postaci zegarowej; sekundy nie niosą tu nic dla czytającego. */
function godzina(znacznikCzasu: number): string {
  return new Date(znacznikCzasu).toLocaleTimeString('pl-PL', { hour: '2-digit', minute: '2-digit' });
}

export function oknoCzatu(zaleznosci: ZaleznosciCzatu): ZamontowaneOknoCzatu {
  let zdjete = false;
  let wpisy: Message[] = [];
  let stan: 'ladowanie' | 'gotowe' | 'odmowa' = 'ladowanie';
  let bladWczytania: ErrorInfo | undefined;
  let wysylka = false;

  const historia = el('div', { klasa: 'sta-kom-historia' });
  const pole = el('textarea', {
    klasa: 'sta-prompt-obszar',
    rows: 1,
    placeholder: tekst('czat.zastepczaTresc'),
    'aria-label': tekst('czat.etykietaTresci'),
  }) as HTMLTextAreaElement;

  const wyslijBtn = el(
    'button',
    {
      klasa: 'dn-btn dn-btn--atrament dn-btn--sm sta-prompt-wyslij',
      type: 'submit',
      'aria-label': tekst('czat.wyslij'),
    },
    [znak('wyslij')],
  ) as HTMLButtonElement;

  function wpis(m: Message): HTMLElement {
    const tresc: Dziecko[] = [m.content];
    /* Wpis w strumieniu niesie znacznik tętna: odpowiedź jeszcze rośnie,
       a wpis bez tego znaku wyglądałby na urwany w połowie zdania. */
    if (m.status === MessageStatus.Streaming) {
      tresc.push(el('span', { klasa: 'pt-tetno', 'aria-hidden': 'true' }));
    }
    const plakietki: Dziecko[] = [];
    if (m.status === MessageStatus.Error) {
      plakietki.push(el('span', { klasa: 'dn-plakietka', tekst: tekst('czat.stanBlad') }));
    } else if (m.status === MessageStatus.Stopped) {
      plakietki.push(el('span', { klasa: 'dn-plakietka', tekst: tekst('czat.stanZatrzymana') }));
    }
    return el('article', { klasa: `sta-wpis ${WARIANT_ROLI[m.role] ?? 'sta-wpis--system'}` }, [
      el('div', { klasa: 'sta-wpis-medalion', 'aria-hidden': 'true' }, [znak(ZNAK_ROLI[m.role] ?? 'plan')]),
      el('div', {}, [
        el('div', { klasa: 'sta-wpis-tozsamosc' }, [
          el('span', { klasa: 'sta-wpis-nadawca', tekst: tekst(`czat.role.${m.role}`) }),
          ...plakietki,
          el('span', { klasa: 'sta-wpis-godzina', tekst: godzina(m.createdAt) }),
        ]),
        el('div', { klasa: 'sta-wpis-tresc' }, tresc),
      ]),
    ]);
  }

  function odswiez(): void {
    if (stan === 'ladowanie') {
      historia.replaceChildren(
        el('div', { klasa: 'dn-wykaz-modulu-poz' }, [
          el('span', { klasa: 'dn-kropka dn-kropka--sygnal dn-kropka--tetno', 'aria-hidden': 'true' }),
          tekst('czat.wczytywanie'),
        ]),
      );
      return;
    }
    if (stan === 'odmowa') {
      historia.replaceChildren(
        el('div', { klasa: 'dn-alert dn-alert--wstega dn-alert--blad', role: 'alert' }, [
          el('span', { klasa: 'dn-alert-tresc' }, [
            el('b', { tekst: tekst('czat.odmowaWykazu') }),
            el('span', { tekst: opisOdmowy(bladWczytania, 'czat') }),
          ]),
        ]),
      );
      return;
    }
    if (wpisy.length === 0) {
      historia.replaceChildren(
        el('p', { klasa: 'dn-pusty-stan dn-pusty-stan--zwarty', tekst: tekst('czat.brakWiadomosci') }),
      );
      return;
    }
    historia.replaceChildren(...wpisy.map(wpis));
    /* Nowy wpis wchodzi na dole, więc historia zjeżdża za nim — inaczej
       odpowiedź powstawałaby poza polem widzenia Operatora. */
    historia.scrollTop = historia.scrollHeight;
  }

  /* Okno modułu powstaje wywołaniem do rdzenia PO montażu bryły, więc przy
     pierwszym wczytaniu jeszcze go nie ma. Wczytanie jednorazowe zostawiłoby
     rozmowę pustą na zawsze — stąd powtórka, gdy okno się pojawi. */
  let wczytaneDlaOkna: string | null = null;

  async function wczytaj(): Promise<void> {
    const okno = zaleznosci.idOkna();
    if (okno === null) {
      stan = 'gotowe';
      odswiez();
      return;
    }
    if (wczytaneDlaOkna === okno) return;
    wczytaneDlaOkna = okno;
    stan = 'ladowanie';
    odswiez();
    const wynik = await wywolaj(zaleznosci.kanal, Command.MessageList, { windowId: okno });
    if (zdjete) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      stan = 'odmowa';
      bladWczytania = wynik.blad;
    } else {
      stan = 'gotowe';
      wpisy = wynik.wynik.messages;
    }
    odswiez();
  }

  async function wyslij(): Promise<void> {
    const okno = zaleznosci.idOkna();
    const trescPolecenia = pole.value.trim();
    if (okno === null || trescPolecenia === '' || wysylka) return;
    wysylka = true;
    wyslijBtn.disabled = true;
    const wynik = await wywolaj(zaleznosci.kanal, Command.MessageSend, {
      windowId: okno,
      content: trescPolecenia,
    });
    if (zdjete) return;
    wysylka = false;
    wyslijBtn.disabled = false;
    if (!wynik.udany) {
      zaloguj(wynik.blad, 'czat.wyslij');
      stan = 'odmowa';
      bladWczytania = wynik.blad;
      odswiez();
      return;
    }
    /* Pole czyści się dopiero po przyjęciu polecenia przez rdzeń: wyczyszczone
       wcześniej zabrałoby Operatorowi tekst, którego rdzeń nie przyjął. */
    pole.value = '';
  }

  const formularz = el('form', { klasa: 'sta-prompt' }, [
    el('span', { klasa: 'sta-prompt-grot', 'aria-hidden': 'true', tekst: '»»' }),
    pole,
    wyslijBtn,
  ]);
  formularz.addEventListener('submit', (zdarzenie) => {
    zdarzenie.preventDefault();
    void wyslij();
  });

  const odsubskrybuj = zaleznosci.kanal.naZdarzenie(EventType.MessageChanged, (zdarzenie) => {
    const okno = zaleznosci.idOkna();
    if (zdjete || okno === null || zdarzenie.message.windowId !== okno) return;
    if (zdarzenie.change === ChangeKind.Deleted) {
      wpisy = wpisy.filter((m) => m.id !== zdarzenie.message.id);
    } else {
      const indeks = wpisy.findIndex((m) => m.id === zdarzenie.message.id);
      if (indeks === -1) wpisy = [...wpisy, zdarzenie.message];
      else wpisy = wpisy.map((m) => (m.id === zdarzenie.message.id ? zdarzenie.message : m));
    }
    stan = 'gotowe';
    odswiez();
  });

  /* Dokument okna zgłasza się zdarzeniem rdzenia — to samo, po którym panele
     poznają swój dokument. Jego nadejście znaczy, że okno już stoi, więc
     rozmowa ma czego szukać. */
  const odsubskrybujDokument = zaleznosci.kanal.naZdarzenie(EventType.StudioDocumentChanged, (zdarzenie) => {
    const okno = zaleznosci.idOkna();
    if (zdjete || okno === null || zdarzenie.document.windowId !== okno) return;
    void wczytaj();
  });

  odswiez();
  void wczytaj();

  const kom = el(
    'section',
    { klasa: 'sta-okno sta-kom', id: 'okno-czat-1', 'data-nazwa': 'Chat Window', 'aria-label': tekst('czat.tytul') },
    [
      el('header', { klasa: 'sta-okno-belka' }, [
        el('span', { klasa: 'sta-okno-tytul' }, [znak('dymek'), el('b', { tekst: tekst('czat.tytul') })]),
      ]),
      el('div', { klasa: 'sta-kom-naglowek' }, [
        el('span', { klasa: 'dn-meta', tekst: tekst('czat.brakParametrow') }),
      ]),
      el('div', { klasa: 'sta-kom-kontekst', 'aria-label': tekst('czat.etykietaKontekst') }, [
        el('span', { klasa: 'dn-meta', tekst: tekst('czat.brakKontekstu') }),
      ]),
      historia,
      el('div', { klasa: 'sta-kom-dol' }, [formularz]),
    ],
  );

  return {
    wezel: el('div', { klasa: 'sta-czaty', 'data-liczba': '1' }, [kom]),
    zdejmij() {
      zdjete = true;
      odsubskrybuj();
      odsubskrybujDokument();
    },
  };
}
