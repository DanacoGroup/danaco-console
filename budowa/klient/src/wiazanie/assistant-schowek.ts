// Panel plików w oknie Assistant: schowek platformy, wycinki tekstu, skrót
// wywołania oraz zajętość okna rozmowy.
import { ClipboardEntryKind, Command } from '../../../shared/contract.ts';
import type { ClipboardEntry, TextSnippet } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { dolozCzynnosciPanelu, zapytajWSzufladzie } from './czynnosci-okna.ts';
import { oglos } from './ogloszenie.ts';

const NAGLOWEK = 'Schowek i wycinki';
const PANEL = 'panel-pliki';

const RODZAJE: ReadonlyArray<readonly [string, string]> = [
  [ClipboardEntryKind.Text, 'Tekst'],
  [ClipboardEntryKind.Image, 'Obraz'],
  [ClipboardEntryKind.File, 'Plik'],
];

interface Otoczenie {
  kanal: Kanal;
  korzen: Element;
  idOkna: () => string;
}

export function zwiazSchowekAsystenta(
  kanal: Kanal,
  korzen: Element,
  idOkna: () => string,
  dolacz: (zdejmij: () => void) => void,
  przy: AddEventListenerOptions,
): void {
  const otoczenie: Otoczenie = { kanal, korzen, idOkna };
  const zdejmij = dolozCzynnosciPanelu(korzen, PANEL, 'Schowek i wycinki', [
    {
      naglowek: 'Schowek platformy',
      pozycje: [
        { kod: 'schowek', nazwa: 'Wykaz schowka' },
        { kod: 'odloz', nazwa: 'Odłóż do schowka…' },
        { kod: 'przypnij', nazwa: 'Przypnij pozycję…' },
        { kod: 'wyrzuc', nazwa: 'Wyrzuć pozycję…' },
      ],
    },
    {
      naglowek: 'Wycinki tekstu',
      pozycje: [
        { kod: 'wycinki', nazwa: 'Wykaz wycinków' },
        { kod: 'wycinek-zapisz', nazwa: 'Zapisz wycinek…' },
        { kod: 'wycinek-zdejmij', nazwa: 'Zdejmij wycinek…' },
      ],
    },
    {
      naglowek: 'Skrót i zajętość',
      pozycje: [
        { kod: 'skrot', nazwa: 'Odczytaj skrót wywołania' },
        { kod: 'skrot-ustaw', nazwa: 'Ustaw skrót wywołania…' },
        { kod: 'zajetosc', nazwa: 'Zajętość okna rozmowy' },
      ],
    },
  ], (kod) => {
    void wykonaj(otoczenie, kod);
  }, przy);
  if (zdejmij !== null) dolacz(zdejmij);
}

function wypelnij(korzen: Element, wiersze: string[]): void {
  const cialo = korzen.querySelector(`#${PANEL} .sta-okno-tresc`);
  if (cialo === null) return;
  cialo.replaceChildren(...wiersze.map((tresc) => {
    const wiersz = cialo.ownerDocument.createElement('div');
    wiersz.className = 'dn-wykaz-modulu-poz';
    wiersz.textContent = tresc;
    return wiersz;
  }));
}

async function wykonaj(otoczenie: Otoczenie, kod: string): Promise<void> {
  if (kod === 'schowek') return wykazSchowka(otoczenie);
  if (kod === 'odloz') return odlozDoSchowka(otoczenie);
  if (kod === 'przypnij') return przypnijPozycje(otoczenie);
  if (kod === 'wyrzuc') return wyrzucPozycje(otoczenie);
  if (kod === 'wycinki') return wykazWycinkow(otoczenie);
  if (kod === 'wycinek-zapisz') return zapiszWycinek(otoczenie);
  if (kod === 'wycinek-zdejmij') return zdejmijWycinek(otoczenie);
  if (kod === 'skrot') return odczytajSkrot(otoczenie);
  if (kod === 'skrot-ustaw') return ustawSkrot(otoczenie);
  if (kod === 'zajetosc') return zajetoscOkna(otoczenie);
  /* Czynność spoza obsłużonych odmawia zamiast milczeć. */
  oglos(NAGLOWEK, `Czynność „${kod}" nie jest prowadzona przez to okno.`, 'ostrzezenie');
}

function opiszPozycje(pozycja: ClipboardEntry): string {
  const zapowiedz = (pozycja.preview ?? pozycja.content).slice(0, 60);
  return `${pozycja.pinned ? '📌 ' : ''}${zapowiedz} · ${String(pozycja.sizeBytes)} B`;
}

async function pozycjeSchowka(kanal: Kanal): Promise<ClipboardEntry[]> {
  const wynik = await wywolaj(kanal, Command.ClipboardList, { limit: 50 });
  return wynik.wynik?.entries ?? [];
}

async function wyborSchowka(kanal: Kanal): Promise<ReadonlyArray<readonly [string, string]>> {
  return (await pozycjeSchowka(kanal)).map((pozycja) =>
    [pozycja.id, opiszPozycje(pozycja)] as const);
}

async function wykazSchowka(otoczenie: Otoczenie): Promise<void> {
  const wynik = await wywolaj(otoczenie.kanal, Command.ClipboardList, { limit: 50 });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wykazu schowka.', 'ostrzezenie');
    return;
  }
  const pozycje = wynik.wynik.entries;
  wypelnij(otoczenie.korzen, pozycje.length === 0
    ? ['Schowek platformy jest pusty.']
    : pozycje.map(opiszPozycje));
  oglos(NAGLOWEK, `Pozycji w schowku: ${String(wynik.wynik.total)}.`);
}

/* Pozycja oznaczona jako wrażliwa nie wraca do okna w zapowiedzi: rdzeń trzyma
   jej treść, ale okno pokazuje samą metrykę. */
async function odlozDoSchowka(otoczenie: Otoczenie): Promise<void> {
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, PANEL, {
    tytul: 'Odłożenie do schowka',
    pola: [
      { klucz: 'tresc', etykieta: 'Treść', obszerne: true, wymagane: true },
      { klucz: 'rodzaj', etykieta: 'Rodzaj', wybor: RODZAJE },
      {
        klucz: 'wrazliwa',
        etykieta: 'Treść wrażliwa',
        wybor: [['nie', 'Nie'], ['tak', 'Tak']],
      },
    ],
    wykonanie: 'Odłóż',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.ClipboardPush, {
    content: wartosci.tresc ?? '',
    kind: (wartosci.rodzaj ?? ClipboardEntryKind.Text) as ClipboardEntryKind,
    sensitive: wartosci.wrazliwa === 'tak',
    ...(otoczenie.idOkna() === '' ? {} : { sourceWindowId: otoczenie.idOkna() }),
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odłożenia do schowka.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Treść leży w schowku platformy.');
  await wykazSchowka(otoczenie);
}

async function przypnijPozycje(otoczenie: Otoczenie): Promise<void> {
  const wybor = await wyborSchowka(otoczenie.kanal);
  if (wybor.length === 0) {
    oglos(NAGLOWEK, 'Schowek platformy jest pusty.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, PANEL, {
    tytul: 'Przypięcie pozycji',
    opis: 'Przypięta pozycja zostaje w schowku mimo sprzątania.',
    pola: [
      { klucz: 'pozycja', etykieta: 'Pozycja', wybor },
      { klucz: 'stan', etykieta: 'Przypięcie', wybor: [['tak', 'Przypnij'], ['nie', 'Odepnij']] },
    ],
    wykonanie: 'Zapisz',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.ClipboardPin, {
    entryId: wartosci.pozycja ?? '',
    pinned: wartosci.stan !== 'nie',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił przypięcia pozycji.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, wartosci.stan === 'nie' ? 'Pozycja odpięta.' : 'Pozycja przypięta.');
  await wykazSchowka(otoczenie);
}

async function wyrzucPozycje(otoczenie: Otoczenie): Promise<void> {
  const wybor = await wyborSchowka(otoczenie.kanal);
  if (wybor.length === 0) {
    oglos(NAGLOWEK, 'Schowek platformy jest pusty.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, PANEL, {
    tytul: 'Wyrzucenie ze schowka',
    pola: [{ klucz: 'pozycja', etykieta: 'Pozycja', wybor }],
    wykonanie: 'Wyrzuć',
    nieodwracalne: 'Treść znika ze schowka platformy bez możliwości odzyskania.',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.ClipboardDelete, {
    entryId: wartosci.pozycja ?? '',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wyrzucenia pozycji.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Pozycja wyrzucona ze schowka.');
  await wykazSchowka(otoczenie);
}

function opiszWycinek(wycinek: TextSnippet): string {
  return `${wycinek.shortcut} → ${wycinek.content.slice(0, 50)}`
    + (wycinek.enabled ? '' : ' · wyłączony');
}

async function wycinki(kanal: Kanal): Promise<TextSnippet[]> {
  const wynik = await wywolaj(kanal, Command.SnippetList, { limit: 50 });
  return wynik.wynik?.snippets ?? [];
}

async function wykazWycinkow(otoczenie: Otoczenie): Promise<void> {
  const wynik = await wywolaj(otoczenie.kanal, Command.SnippetList, { limit: 50 });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wykazu wycinków.', 'ostrzezenie');
    return;
  }
  const zbior = wynik.wynik.snippets;
  wypelnij(otoczenie.korzen, zbior.length === 0
    ? ['Nie ma jeszcze żadnego wycinka tekstu.']
    : zbior.map(opiszWycinek));
  oglos(NAGLOWEK, `Wycinków: ${String(wynik.wynik.total)}.`);
}

/* Skrót wycinka jest jego tożsamością przy pisaniu: to on rozwija się w treść,
   więc dwa wycinki nie mogą nosić tego samego skrótu. */
async function zapiszWycinek(otoczenie: Otoczenie): Promise<void> {
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, PANEL, {
    tytul: 'Zapis wycinka',
    pola: [
      { klucz: 'skrot', etykieta: 'Skrót', wymagane: true },
      { klucz: 'tresc', etykieta: 'Treść', obszerne: true, wymagane: true },
      { klucz: 'opis', etykieta: 'Opis' },
    ],
    wykonanie: 'Zapisz wycinek',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.SnippetSet, {
    shortcut: wartosci.skrot ?? '',
    content: wartosci.tresc ?? '',
    ...(wartosci.opis === '' ? {} : { description: wartosci.opis }),
    enabled: true,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu wycinka.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Wycinek „${wartosci.skrot ?? ''}" zapisany.`);
  await wykazWycinkow(otoczenie);
}

async function zdejmijWycinek(otoczenie: Otoczenie): Promise<void> {
  const zbior = await wycinki(otoczenie.kanal);
  if (zbior.length === 0) {
    oglos(NAGLOWEK, 'Nie ma jeszcze żadnego wycinka tekstu.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, PANEL, {
    tytul: 'Zdjęcie wycinka',
    pola: [{
      klucz: 'wycinek',
      etykieta: 'Wycinek',
      wybor: zbior.map((wycinek) => [wycinek.id, opiszWycinek(wycinek)] as const),
    }],
    wykonanie: 'Zdejmij wycinek',
    nieodwracalne: 'Skrót przestanie się rozwijać w treść.',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.SnippetDelete, {
    snippetId: wartosci.wycinek ?? '',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zdjęcia wycinka.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Wycinek zdjęty.');
  await wykazWycinkow(otoczenie);
}

async function odczytajSkrot(otoczenie: Otoczenie): Promise<void> {
  const wynik = await wywolaj(otoczenie.kanal, Command.LauncherHotkeyGet, {});
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu skrótu.', 'ostrzezenie');
    return;
  }
  const stan = wynik.wynik;
  const powod = stan.reason === undefined || stan.reason === '' ? '' : ` · ${stan.reason}`;
  wypelnij(otoczenie.korzen, [
    `skrót: ${stan.hotkey}`,
    `zarejestrowany: ${stan.registered ? 'tak' : 'nie'}`,
    `obsługiwany przez system: ${stan.supported ? 'tak' : 'nie'}${powod}`,
  ]);
  oglos(NAGLOWEK, `Skrót wywołania: ${stan.hotkey}.`);
}

async function ustawSkrot(otoczenie: Otoczenie): Promise<void> {
  const stojacy = await wywolaj(otoczenie.kanal, Command.LauncherHotkeyGet, {});
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, PANEL, {
    tytul: 'Skrót wywołania',
    opis: 'Zapis klawiszy rozdzielony znakiem plus, na przykład Ctrl+Alt+Spacja.',
    pola: [{
      klucz: 'skrot',
      etykieta: 'Skrót',
      wartosc: stojacy.wynik?.hotkey ?? '',
      wymagane: true,
    }],
    wykonanie: 'Ustaw skrót',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.LauncherHotkeySet, {
    hotkey: wartosci.skrot ?? '',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił ustawienia skrótu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Skrót wywołania to teraz ${wartosci.skrot ?? ''}.`);
  await odczytajSkrot(otoczenie);
}

/* Zajętość liczy rdzeń dzielnikiem modelu okna: bez okna nie ma czego mierzyć. */
async function zajetoscOkna(otoczenie: Otoczenie): Promise<void> {
  if (otoczenie.idOkna() === '') {
    oglos(NAGLOWEK, 'To okno nie stoi jeszcze przy sesji.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(otoczenie.kanal, Command.ContextUsageGet, {
    windowId: otoczenie.idOkna(),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił pomiaru zajętości.', 'ostrzezenie');
    return;
  }
  if (!wynik.wynik.available) {
    oglos(NAGLOWEK, wynik.wynik.reason ?? 'Rdzeń nie zmierzył zajętości tego okna.',
      'ostrzezenie');
    return;
  }
  const zajetosc = wynik.wynik.usage;
  wypelnij(otoczenie.korzen, [
    `zajęte: ${String(zajetosc.usedTokens)} z ${String(zajetosc.limitTokens)} żetonów`,
    `dzielnik: ${zajetosc.tokenizer}`,
    `historia: ${zajetosc.historyTokens === undefined
      ? 'bez wskazania' : String(zajetosc.historyTokens)}`,
  ]);
  oglos(NAGLOWEK, `Okno zajmuje ${String(zajetosc.usedTokens)} żetonów.`);
}
