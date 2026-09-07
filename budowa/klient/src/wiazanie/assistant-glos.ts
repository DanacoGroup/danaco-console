// Głos w oknie Assistant: gotowość toru mowy, nasłuch, słowo budzące oraz
// praca na nagraniu — wniesienie, odczyt i spisanie.
import { Command, ListenMode } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { dolozCzynnosciPanelu, zapytajWSzufladzie } from './czynnosci-okna.ts';
import { oglos } from './ogloszenie.ts';

const NAGLOWEK = 'Głos';

/* Nagranie wniesione żyje w oknie: rdzeń oddaje jego oznaczenie, a odczyt i
   spisanie potrzebują tego samego oznaczenia w następnej czynności. */
const NAGRANIA = new Map<string, string>();

const TRYBY: ReadonlyArray<readonly [string, string]> = [
  [ListenMode.PushToTalk, 'Na przycisk'],
  [ListenMode.WakeWord, 'Po słowie budzącym'],
  [ListenMode.Continuous, 'Ciągły'],
];

interface Otoczenie {
  kanal: Kanal;
  korzen: Element;
  idOkna: () => string;
}

export function zwiazGlosAsystenta(
  kanal: Kanal,
  korzen: Element,
  idOkna: () => string,
  dolacz: (zdejmij: () => void) => void,
  przy: AddEventListenerOptions,
): void {
  const otoczenie: Otoczenie = { kanal, korzen, idOkna };
  const zdejmij = dolozCzynnosciPanelu(korzen, 'panel-voice', 'Czynności głosu', [
    {
      naglowek: 'Tor mowy',
      pozycje: [
        { kod: 'gotowosc', nazwa: 'Sprawdź gotowość toru' },
        { kod: 'budzenie', nazwa: 'Odczytaj słowo budzące' },
        { kod: 'budzenie-ustaw', nazwa: 'Ustaw słowo budzące…' },
      ],
    },
    {
      naglowek: 'Nasłuch',
      pozycje: [
        { kod: 'nasluch', nazwa: 'Zacznij nasłuch…' },
        { kod: 'nasluch-koniec', nazwa: 'Zakończ nasłuch' },
      ],
    },
    {
      naglowek: 'Nagranie',
      pozycje: [
        { kod: 'wnies', nazwa: 'Wnieś nagranie…' },
        { kod: 'odczytaj', nazwa: 'Odczytaj nagranie' },
        { kod: 'spisz', nazwa: 'Spisz nagranie…' },
      ],
    },
  ], (kod) => {
    void wykonaj(otoczenie, kod, 'panel-voice');
  }, przy);
  if (zdejmij !== null) dolacz(zdejmij);
}

function wypelnij(korzen: Element, wiersze: string[]): void {
  const cialo = korzen.querySelector('#panel-voice .sta-okno-tresc');
  if (cialo === null) return;
  cialo.replaceChildren(...wiersze.map((tresc) => {
    const wiersz = cialo.ownerDocument.createElement('div');
    wiersz.className = 'dn-wykaz-modulu-poz';
    wiersz.textContent = tresc;
    return wiersz;
  }));
}

async function wykonaj(otoczenie: Otoczenie, kod: string, panel: string): Promise<void> {
  if (kod === 'gotowosc') return sprawdzGotowosc(otoczenie);
  if (kod === 'budzenie') return odczytajBudzenie(otoczenie);
  if (kod === 'budzenie-ustaw') return ustawBudzenie(otoczenie, panel);
  if (kod === 'nasluch') return zacznijNasluch(otoczenie, panel);
  if (kod === 'nasluch-koniec') return zakonczNasluch(otoczenie);
  if (kod === 'wnies') return wniesNagranie(otoczenie, panel);
  if (kod === 'odczytaj') return odczytajNagranie(otoczenie);
  if (kod === 'spisz') return spiszNagranie(otoczenie, panel);
  /* Czynność spoza obsłużonych odmawia zamiast milczeć: pozycja bez gałęzi
     wyglądałaby jak działająca. */
  oglos(NAGLOWEK, `Czynność „${kod}" nie jest prowadzona przez to okno.`, 'ostrzezenie');
}

/* Gotowość mówi o każdej drodze z osobna: tor może spisywać nagrania, a nie
   umieć nasłuchiwać, i odwrotnie. */
async function sprawdzGotowosc({ kanal, korzen }: Otoczenie): Promise<void> {
  const wynik = await wywolaj(kanal, Command.SpeechAvailabilityGet, {});
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił sprawdzenia toru mowy.', 'ostrzezenie');
    return;
  }
  const stan = wynik.wynik;
  wypelnij(korzen, [
    `tor mowy: ${stan.available ? 'gotowy' : `niegotowy — ${stan.reason ?? 'bez wyjaśnienia'}`}`,
    `narzędzie: ${stan.engine ?? 'nieznane'} · model: ${stan.model ?? 'nieznany'}`,
    `nasłuch: ${stan.listenAvailable === true ? 'dostępny' : 'niedostępny'}`,
    `słowo budzące: ${stan.wakeWordAvailable === true ? 'dostępne' : 'niedostępne'}`,
    `wnoszenie nagrań: ${stan.uploadAvailable === true ? 'dostępne' : 'niedostępne'}`,
    `synteza mowy: ${stan.synthesisAvailable === true
      ? 'dostępna'
      : `niedostępna — ${stan.synthesisReason ?? 'bez wyjaśnienia'}`}`,
  ]);
}

async function odczytajBudzenie({ kanal, korzen }: Otoczenie): Promise<void> {
  const wynik = await wywolaj(kanal, Command.SpeechWakeGet, {});
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu słowa budzącego.', 'ostrzezenie');
    return;
  }
  const nastawa = wynik.wynik.config;
  wypelnij(korzen, [
    `słowo budzące: ${nastawa.phrase}`,
    `tryb nasłuchu: ${nastawa.mode}`,
    `próg głosu: ${nastawa.vadThreshold}`,
    `tłumienie szumu: ${nastawa.noiseSuppression ? 'włączone' : 'wyłączone'}`,
  ]);
  if (!wynik.wynik.available) {
    oglos(NAGLOWEK,
      `Słowo budzące nie działa: ${wynik.wynik.reason ?? 'bez wyjaśnienia'}.`, 'ostrzezenie');
  }
}

async function ustawBudzenie({ kanal, korzen }: Otoczenie, panel: string): Promise<void> {
  const stojace = await wywolaj(kanal, Command.SpeechWakeGet, {});
  const nastawa = stojace.wynik?.config;
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Słowo budzące',
    pola: [
      { klucz: 'slowo', etykieta: 'Słowo budzące', wartosc: nastawa?.phrase ?? '', wymagane: true },
      { klucz: 'tryb', etykieta: 'Tryb nasłuchu', wybor: TRYBY, wartosc: nastawa?.mode },
      {
        klucz: 'prog',
        etykieta: 'Próg głosu',
        wartosc: nastawa === undefined ? '0.5' : String(nastawa.vadThreshold),
      },
      {
        klucz: 'szum',
        etykieta: 'Tłumienie szumu',
        wybor: [['tak', 'Włączone'], ['nie', 'Wyłączone']],
        wartosc: nastawa?.noiseSuppression === false ? 'nie' : 'tak',
      },
    ],
    wykonanie: 'Ustaw słowo budzące',
  });
  if (wartosci === null) return;
  const prog = Number.parseFloat(wartosci.prog ?? '');
  const wynik = await wywolaj(kanal, Command.SpeechWakeSet, {
    phrase: wartosci.slowo ?? '',
    mode: (wartosci.tryb ?? ListenMode.WakeWord) as ListenMode,
    ...(Number.isFinite(prog) ? { vadThreshold: prog } : {}),
    noiseSuppression: wartosci.szum !== 'nie',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił ustawienia słowa budzącego.',
      'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Słowo budzące to teraz „${wartosci.slowo ?? ''}".`);
  await odczytajBudzenie({ kanal, korzen, idOkna: () => '' });
}

async function zacznijNasluch({ kanal, korzen, idOkna }: Otoczenie, panel: string):
Promise<void> {
  if (idOkna() === '') {
    oglos(NAGLOWEK, 'Rdzeń nie dał okna asystenta dla tej karty.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Nasłuch mowy',
    pola: [
      { klucz: 'tryb', etykieta: 'Tryb nasłuchu', wybor: TRYBY },
      { klucz: 'jezyk', etykieta: 'Język', podpowiedz: 'Puste bierze ustawienie rdzenia' },
    ],
    wykonanie: 'Zacznij nasłuch',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.SpeechListenStart, {
    windowId: idOkna(),
    mode: (wartosci.tryb ?? ListenMode.PushToTalk) as ListenMode,
    ...(wartosci.jezyk === '' ? {} : { language: wartosci.jezyk }),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił nasłuchu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, wynik.wynik.listening
    ? 'Nasłuch stoi; mowa przyjdzie zdarzeniem.'
    : `Rdzeń nie zaczął nasłuchu: ${wynik.wynik.reason ?? 'bez wyjaśnienia'}.`,
  wynik.wynik.listening ? 'informacja' : 'ostrzezenie');
}

async function zakonczNasluch({ kanal, idOkna }: Otoczenie): Promise<void> {
  if (idOkna() === '') {
    oglos(NAGLOWEK, 'Rdzeń nie dał okna asystenta dla tej karty.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.SpeechListenStop, { windowId: idOkna() });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zakończenia nasłuchu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Nasłuch zakończony.');
}

/* Nagranie wnosi się zapisem base64, bo tego oczekuje rdzeń; okno nie ma
   dostępu do plików maszyny Operatora inaczej niż przez wklejoną treść. */
async function wniesNagranie({ kanal, korzen, idOkna }: Otoczenie, panel: string):
Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Wniesienie nagrania',
    opis: 'Treść nagrania podaje się zapisem base64; rdzeń odda jego oznaczenie.',
    pola: [
      { klucz: 'tresc', etykieta: 'Nagranie w zapisie base64', obszerne: true, wymagane: true },
      { klucz: 'rodzaj', etykieta: 'Rodzaj treści', wartosc: 'audio/wav', wymagane: true },
      {
        klucz: 'zachowaj',
        etykieta: 'Nagranie w magazynie',
        wybor: [['nie', 'Zdejmij po spisaniu'], ['tak', 'Zachowaj']],
      },
    ],
    wykonanie: 'Wnieś nagranie',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.SpeechAudioUpload, {
    audio: wartosci.tresc ?? '',
    contentType: wartosci.rodzaj ?? 'audio/wav',
    ...(idOkna() === '' ? {} : { windowId: idOkna() }),
    retain: wartosci.zachowaj === 'tak',
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił przyjęcia nagrania.', 'ostrzezenie');
    return;
  }
  NAGRANIA.set(idOkna(), wynik.wynik.audioRef);
  oglos(NAGLOWEK, `Nagranie przyjęte jako ${wynik.wynik.audioRef}; `
    + `${wynik.wynik.sizeBytes} bajtów.`);
}

function nagranie(idOkna: string): string {
  const zapamietane = NAGRANIA.get(idOkna) ?? '';
  if (zapamietane === '') {
    oglos(NAGLOWEK, 'Okno nie zna nagrania — wnieś je najpierw.', 'ostrzezenie');
  }
  return zapamietane;
}

/* Odczyt nazywa rozmiar i czas trwania, a nie oddaje treści: nagranie w zapisie
   base64 nie jest niczym, co da się przeczytać w oknie. */
async function odczytajNagranie({ kanal, korzen, idOkna }: Otoczenie): Promise<void> {
  const cel = nagranie(idOkna());
  if (cel === '') return;
  const wynik = await wywolaj(kanal, Command.SpeechAudioFetch, { audioRef: cel });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu nagrania.', 'ostrzezenie');
    return;
  }
  wypelnij(korzen, [
    `oznaczenie: ${cel}`,
    `rodzaj treści: ${wynik.wynik.contentType}`,
    `rozmiar: ${wynik.wynik.sizeBytes} bajtów`,
    `czas trwania: ${wynik.wynik.durationMs ?? 'nieznany'} ms`,
  ]);
}

async function spiszNagranie({ kanal, korzen, idOkna }: Otoczenie, panel: string):
Promise<void> {
  const cel = nagranie(idOkna());
  if (cel === '') return;
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Spisanie nagrania',
    pola: [
      { klucz: 'jezyk', etykieta: 'Język', podpowiedz: 'Puste zostawia rozpoznanie rdzeniowi' },
      { klucz: 'model', etykieta: 'Model spisujący', podpowiedz: 'Puste bierze model domyślny' },
    ],
    wykonanie: 'Spisz nagranie',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.SpeechTranscribe, {
    audioRef: cel,
    ...(wartosci.jezyk === '' ? {} : { language: wartosci.jezyk }),
    ...(wartosci.model === '' ? {} : { model: wartosci.model }),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił spisania nagrania.', 'ostrzezenie');
    return;
  }
  if (!wynik.wynik.processed) {
    oglos(NAGLOWEK, 'Rdzeń nie spisał nagrania.', 'ostrzezenie');
    return;
  }
  wypelnij(korzen, wynik.wynik.transcript.split('\n'));
  oglos(NAGLOWEK, `Spisano ${wynik.wynik.characters} znaków modelem ${wynik.wynik.model}`
    + `${wynik.wynik.confidence === undefined
      ? '' : `, pewność ${wynik.wynik.confidence}`}.`);
}
