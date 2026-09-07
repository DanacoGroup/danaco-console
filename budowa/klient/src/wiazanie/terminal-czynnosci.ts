// Czynności okna Terminal rozłożone na trzy panele wedle tego, czego dotykają:
// karty powłoki, wyjście procesu i zasoby połączeń — gospodarze, klucze,
// skrypty i tunele.
import {
  Command,
  EventType,
  TerminalKeyType,
  TerminalScriptKind,
  TerminalShell,
  TerminalTunnelKind,
} from '../../../shared/contract.ts';
import type { TerminalOutputLine } from '../../../shared/contract.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import { zglosUchwyt } from '../polaczenie/rozdzielacz-zdarzen.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { dolozCzynnosciPanelu, zapytajWSzufladzie } from './czynnosci-okna.ts';
import { oglos } from './ogloszenie.ts';

const NAGLOWEK = 'Terminal';

const POWLOKI: ReadonlyArray<readonly [string, string]> = [
  [TerminalShell.Bash, 'Bash'],
  [TerminalShell.Powershell, 'PowerShell'],
  [TerminalShell.Cmd, 'Wiersz poleceń'],
  [TerminalShell.Node, 'Node'],
  [TerminalShell.Python, 'Python'],
];

interface Otoczenie {
  kanal: Kanal;
  korzen: Element;
  idOkna: () => string;
  odswiez: () => void;
  dolacz: (odsubskrybuj: Odsubskrybuj) => void;
}

/*
zwiazCzynnosciTerminala wiesza czynności na belkach trzech paneli okna.

Podział idzie za tym, czym panel jest: karty powłoki dostają czynności kart,
wyjście — czynności odczytu, a panel plików zasoby połączeń, bo tam Operator
patrzy na to, co leży poza samą kartą.
*/
export function zwiazCzynnosciTerminala(
  kanal: Kanal,
  korzen: Element,
  idOkna: () => string,
  odswiez: () => void,
  dolacz: (odsubskrybuj: Odsubskrybuj) => void,
  przy: AddEventListenerOptions,
): void {
  const otoczenie: Otoczenie = { kanal, korzen, idOkna, odswiez, dolacz };
  const zdejmij = [
    dolozCzynnosciPanelu(korzen, 'panel-tabs', 'Czynności kart powłoki', [
      {
        naglowek: 'Karta powłoki',
        pozycje: [
          { kod: 'karta-otworz', nazwa: 'Otwórz kartę powłoki…' },
          { kod: 'karta-obserwuj', nazwa: 'Obserwuj pliki karty…' },
        ],
      },
      {
        naglowek: 'Poza maszyną rdzenia',
        pozycje: [
          { kod: 'tunel-otworz', nazwa: 'Otwórz tunel przez gospodarza…' },
        ],
      },
    ], (kod) => {
      void wykonaj(otoczenie, kod, 'panel-tabs');
    }, przy),
    dolozCzynnosciPanelu(korzen, 'panel-output', 'Czynności wyjścia', [
      {
        naglowek: 'Wyjście procesu',
        pozycje: [
          { kod: 'wyjscie-odczytaj', nazwa: 'Odczytaj wyjście ostatniego procesu' },
          { kod: 'wyjscie-sledz', nazwa: 'Śledź wyjście wszystkich kart' },
        ],
      },
      {
        naglowek: 'Pliki karty',
        pozycje: [
          { kod: 'plik-odczytaj', nazwa: 'Odczytaj plik…' },
        ],
      },
    ], (kod) => {
      void wykonaj(otoczenie, kod, 'panel-output');
    }, przy),
    dolozCzynnosciPanelu(korzen, 'panel-pliki', 'Zasoby połączeń', [
      {
        naglowek: 'Książka gospodarzy',
        pozycje: [
          { kod: 'gospodarz-zapisz', nazwa: 'Zapisz gospodarza…' },
          { kod: 'gospodarz-usun', nazwa: 'Usuń gospodarza…' },
        ],
      },
      {
        naglowek: 'Klucze SSH',
        pozycje: [
          { kod: 'klucz-wytworz', nazwa: 'Wytwórz klucz…' },
          { kod: 'klucz-wnies', nazwa: 'Wnieś klucz z dysku rdzenia…' },
          { kod: 'klucz-usun', nazwa: 'Zdejmij klucz z wykazu…' },
        ],
      },
      {
        naglowek: 'Biblioteka skryptów',
        pozycje: [
          { kod: 'skrypt-zapisz', nazwa: 'Zapisz skrypt…' },
          { kod: 'skrypt-sprawdz', nazwa: 'Sprawdź treść skryptu…' },
          { kod: 'skrypt-usun', nazwa: 'Usuń skrypt…' },
        ],
      },
    ], (kod) => {
      void wykonaj(otoczenie, kod, 'panel-pliki');
    }, przy),
  ];
  for (const zdjecie of zdejmij) {
    if (zdjecie !== null) {
      dolacz(() => {
        zdjecie();
      });
    }
  }
}

async function wykonaj(otoczenie: Otoczenie, kod: string, panel: string): Promise<void> {
  if (otoczenie.idOkna() === '') {
    oglos(NAGLOWEK, 'Rdzeń nie dał okna terminala dla tej karty.', 'ostrzezenie');
    return;
  }
  if (kod === 'karta-otworz') return otworzKarte(otoczenie, panel);
  if (kod === 'karta-obserwuj') return obserwujPliki(otoczenie, panel);
  if (kod === 'tunel-otworz') return otworzTunel(otoczenie, panel);
  if (kod === 'wyjscie-odczytaj') return odczytajWyjscie(otoczenie);
  if (kod === 'wyjscie-sledz') return sledzWyjscie(otoczenie);
  if (kod === 'plik-odczytaj') return odczytajPlik(otoczenie, panel);
  if (kod === 'gospodarz-zapisz') return zapiszGospodarza(otoczenie, panel);
  if (kod === 'gospodarz-usun') return usunGospodarza(otoczenie, panel);
  if (kod === 'klucz-wytworz') return wytworzKlucz(otoczenie, panel);
  if (kod === 'klucz-wnies') return wniesKlucz(otoczenie, panel);
  if (kod === 'klucz-usun') return usunKlucz(otoczenie, panel);
  if (kod === 'skrypt-zapisz') return zapiszSkrypt(otoczenie, panel);
  if (kod === 'skrypt-sprawdz') return sprawdzSkrypt(otoczenie, panel);
  if (kod === 'skrypt-usun') return usunSkrypt(otoczenie, panel);
  /* Czynność spoza obsłużonych odmawia zamiast milczeć: pozycja bez gałęzi
     wyglądałaby jak działająca. */
  oglos(NAGLOWEK, `Czynność „${kod}" nie jest prowadzona przez to okno.`, 'ostrzezenie');
}

async function pierwszaKarta(kanal: Kanal, idOkna: string): Promise<string> {
  const wynik = await wywolaj(kanal, Command.TerminalSessionList, { windowId: idOkna });
  const karta = wynik.wynik?.sessions[0]?.id ?? '';
  if (karta === '') {
    oglos(NAGLOWEK, 'Okno nie ma jeszcze żadnej karty powłoki.', 'ostrzezenie');
  }
  return karta;
}

async function otworzKarte({ kanal, korzen, idOkna, odswiez }: Otoczenie, panel: string):
Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Otwarcie karty powłoki',
    pola: [
      { klucz: 'nazwa', etykieta: 'Nazwa karty', podpowiedz: 'Bez nazwy rdzeń nazwie ją sam' },
      { klucz: 'powloka', etykieta: 'Powłoka', wybor: POWLOKI },
      { klucz: 'katalog', etykieta: 'Katalog roboczy', podpowiedz: 'Domyślnie katalog okna' },
    ],
    wykonanie: 'Otwórz kartę',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.TerminalSessionOpen, {
    windowId: idOkna(),
    shell: (wartosci.powloka ?? TerminalShell.Bash) as TerminalShell,
    ...(wartosci.nazwa === '' ? {} : { title: wartosci.nazwa }),
    ...(wartosci.katalog === '' ? {} : { workingDir: wartosci.katalog }),
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił otwarcia karty powłoki.', 'ostrzezenie');
    return;
  }
  odswiez();
}

/* Wyjście czytane jest z procesu stojącego w wykazie najwyżej: okno nie
   prowadzi wskazania procesu, a wykaz podaje ich kolejność. */
async function odczytajWyjscie({ kanal, korzen, idOkna }: Otoczenie): Promise<void> {
  const procesy = await wywolaj(kanal, Command.TerminalProcessList, { windowId: idOkna() });
  const idProcesu = procesy.wynik?.processes[0]?.id ?? '';
  if (idProcesu === '') {
    oglos(NAGLOWEK, 'Żaden proces nie stoi w wykazie tego okna.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.TerminalOutputRead, { processId: idProcesu, tail: 60 });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu wyjścia.', 'ostrzezenie');
    return;
  }
  const wiersze = [wynik.wynik.stdout, wynik.wynik.stderr]
    .join('\n')
    .split('\n')
    .filter((wiersz) => wiersz.trim() !== '');
  wypelnijWyjscie(korzen, wiersze);
  if (wiersze.length === 0) {
    oglos(NAGLOWEK, 'Proces nie wypisał dotąd niczego.');
    return;
  }
  if (wynik.wynik.truncated) {
    oglos(NAGLOWEK, 'Wyjście przycięte granicą bufora — proces wypisał więcej.', 'ostrzezenie');
  }
}

function wypelnijWyjscie(korzen: Element, wiersze: string[]): void {
  const cialo = korzen.querySelector('#panel-output .sta-okno-tresc');
  if (cialo === null) return;
  cialo.replaceChildren(...wiersze.map((tekst) => wierszWyjscia(cialo, tekst)));
}

function wierszWyjscia(cialo: Element, tekst: string): HTMLElement {
  const wiersz = cialo.ownerDocument.createElement('div');
  wiersz.className = 'oc-wiersz';
  wiersz.textContent = tekst;
  return wiersz;
}

/*
sledzWyjscie zapisuje okno na zbiorcze wyjście wszystkich kart powłoki.

Zapis jest subskrypcją, nie jednorazowym odczytem: wiersze przychodzą potem
zdarzeniem `stream.chunk`. Uchwyt idzie do zwolnień karty, inaczej zostałby
po jej zamknięciu i dopisywał do zdjętego już drzewa.
*/
async function sledzWyjscie({ kanal, korzen, idOkna, dolacz }: Otoczenie): Promise<void> {
  const panel = korzen.querySelector('#panel-output');
  const cialo = panel?.querySelector('.sta-okno-tresc') ?? null;
  if (panel === null || panel === undefined || cialo === null) return;
  if (panel.getAttribute('data-strumien') === 'tak') {
    oglos(NAGLOWEK, 'Okno już śledzi zbiorcze wyjście kart.');
    return;
  }
  const okno = idOkna();
  const wynik = await wywolaj(kanal, Command.TerminalOutputStream, { windowId: okno, tail: 60 });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu na wyjście.', 'ostrzezenie');
    return;
  }
  if (!wynik.wynik.subscribed) {
    oglos(NAGLOWEK, 'Rdzeń oddał ogon historii, ale obserwacji nie założył.', 'ostrzezenie');
    return;
  }
  panel.setAttribute('data-strumien', 'tak');
  cialo.replaceChildren(...wynik.wynik.lines.map((wiersz: TerminalOutputLine) =>
    wierszWyjscia(cialo, wiersz.text)));
  dolacz(zglosUchwyt(EventType.StreamChunk, (tresc) => {
    if (tresc.windowId !== okno) return;
    const tekst = (tresc.text ?? '').replace(/\n$/, '');
    if (tekst === '') return;
    cialo.appendChild(wierszWyjscia(cialo, tekst));
  }));
  oglos(NAGLOWEK, 'Okno śledzi zbiorcze wyjście kart powłoki.');
}

async function odczytajPlik({ kanal, korzen, idOkna }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Odczyt pliku karty',
    opis: 'Plik czytany jest na maszynie karty powłoki i w jej katalogu roboczym.',
    pola: [
      {
        klucz: 'sciezka',
        etykieta: 'Ścieżka pliku',
        podpowiedz: 'Względem katalogu karty albo bezwzględna',
        wymagane: true,
      },
      { klucz: 'ogon', etykieta: 'Ostatnich wierszy', wartosc: '200' },
    ],
    wykonanie: 'Odczytaj plik',
  });
  if (wartosci === null) return;
  const idKarty = await pierwszaKarta(kanal, idOkna());
  if (idKarty === '') return;
  const ogon = Number.parseInt(wartosci.ogon ?? '', 10);
  const wynik = await wywolaj(kanal, Command.TerminalFileRead, {
    sessionId: idKarty,
    path: wartosci.sciezka ?? '',
    ...(Number.isFinite(ogon) ? { tail: ogon } : {}),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu pliku.', 'ostrzezenie');
    return;
  }
  wypelnijWyjscie(korzen, wynik.wynik.content.split('\n'));
  oglos(NAGLOWEK, wynik.wynik.truncated
    ? `Plik ${wynik.wynik.path} przycięty granicą odczytu.`
    : `Plik ${wynik.wynik.path} odczytany w całości.`);
}

/* Obserwacja bez polecenia po zmianie jest bezczynna, więc okno pyta o oba i
   bez polecenia czynności nie wysyła. */
async function obserwujPliki({ kanal, korzen, idOkna }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Obserwacja plików karty',
    pola: [
      { klucz: 'wzorzec', etykieta: 'Wzorzec ścieżek', podpowiedz: '*.go', wymagane: true },
      {
        klucz: 'polecenie',
        etykieta: 'Polecenie po zmianie',
        podpowiedz: 'go build ./...',
        wymagane: true,
      },
    ],
    wykonanie: 'Załóż obserwację',
  });
  if (wartosci === null) return;
  const idKarty = await pierwszaKarta(kanal, idOkna());
  if (idKarty === '') return;
  const wynik = await wywolaj(kanal, Command.TerminalWatchStart, {
    sessionId: idKarty,
    pattern: wartosci.wzorzec ?? '',
    command: wartosci.polecenie ?? '',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił założenia obserwacji.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Obserwacja stoi na wzorcu „${wartosci.wzorzec ?? ''}".`);
}

async function zapiszGospodarza({ kanal, korzen, odswiez }: Otoczenie, panel: string):
Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Wpis książki gospodarzy',
    pola: [
      { klucz: 'nazwa', etykieta: 'Nazwa wpisu', podpowiedz: 'serwer budowy', wymagane: true },
      { klucz: 'cel', etykieta: 'Adres celu', podpowiedz: 'ubuntu@10.0.0.1', wymagane: true },
      { klucz: 'port', etykieta: 'Port', podpowiedz: 'Domyślnie port protokołu' },
      { klucz: 'katalog', etykieta: 'Katalog roboczy' },
    ],
    wykonanie: 'Zapisz gospodarza',
  });
  if (wartosci === null) return;
  const port = Number.parseInt(wartosci.port ?? '', 10);
  const wynik = await wywolaj(kanal, Command.TerminalHostSave, {
    host: {
      id: '',
      name: wartosci.nazwa ?? '',
      target: wartosci.cel ?? '',
      ...(Number.isFinite(port) ? { port } : {}),
      ...(wartosci.katalog === '' ? {} : { workingDir: wartosci.katalog }),
      createdAt: 0,
      updatedAt: 0,
    },
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu gospodarza.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Gospodarz „${wartosci.nazwa ?? ''}" stoi w książce.`);
  odswiez();
}

async function usunGospodarza({ kanal, korzen, odswiez }: Otoczenie, panel: string): Promise<void> {
  const wykaz = await wywolaj(kanal, Command.TerminalHostList, {});
  const gospodarze = wykaz.wynik?.hosts ?? [];
  if (gospodarze.length === 0) {
    oglos(NAGLOWEK, 'Książka gospodarzy jest pusta.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Usunięcie gospodarza',
    pola: [{
      klucz: 'gospodarz',
      etykieta: 'Wpis książki',
      wybor: gospodarze.map((wpis) => [wpis.id, `${wpis.name} · ${wpis.target}`] as const),
    }],
    wykonanie: 'Usuń wpis',
    nieodwracalne: 'Wpis znika z książki; karty już otwarte do tego gospodarza zostają.',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.TerminalHostRemove, {
    hostId: wartosci.gospodarz ?? '',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił usunięcia gospodarza.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Gospodarz zdjęty z książki.');
  odswiez();
}

/* Materiał klucza prywatnego nigdy nie idzie do ogłoszenia ani do schowka:
   okno nazywa oznaczenie i odcisk, bo to wystarcza do rozpoznania klucza. */
async function wytworzKlucz({ kanal, korzen }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Wytworzenie klucza',
    opis: 'Klucz powstaje na maszynie rdzenia; okno pozna jego odcisk, nie treść.',
    pola: [
      { klucz: 'nazwa', etykieta: 'Nazwa klucza', wymagane: true },
      {
        klucz: 'rodzaj',
        etykieta: 'Rodzaj',
        wybor: [
          [TerminalKeyType.Ed25519, 'Ed25519'],
          [TerminalKeyType.Ecdsa, 'ECDSA'],
          [TerminalKeyType.Rsa, 'RSA'],
        ],
      },
      { klucz: 'komentarz', etykieta: 'Komentarz klucza publicznego' },
    ],
    wykonanie: 'Wytwórz klucz',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.TerminalKeyGenerate, {
    name: wartosci.nazwa ?? '',
    keyType: (wartosci.rodzaj ?? TerminalKeyType.Ed25519) as TerminalKeyType,
    ...(wartosci.komentarz === '' ? {} : { comment: wartosci.komentarz }),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wytworzenia klucza.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Klucz „${wartosci.nazwa ?? ''}" wytworzony; `
    + `odcisk ${wynik.wynik.key.fingerprint}.`);
}

/* Wniesienie bierze klucz leżący już na maszynie rdzenia, więc okno pyta o
   ścieżkę, a nie o treść klucza — treść nie przechodzi tędy nigdy. */
async function wniesKlucz({ kanal, korzen }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Wniesienie klucza',
    opis: 'Klucz musi leżeć na maszynie rdzenia; okno podaje jego miejsce, nie treść.',
    pola: [
      { klucz: 'nazwa', etykieta: 'Nazwa w wykazie', wymagane: true },
      {
        klucz: 'sciezka',
        etykieta: 'Ścieżka klucza prywatnego',
        podpowiedz: '/home/operator/.ssh/id_ed25519',
        wymagane: true,
      },
    ],
    wykonanie: 'Wnieś klucz',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.TerminalKeyImport, {
    name: wartosci.nazwa ?? '',
    path: wartosci.sciezka ?? '',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wniesienia klucza.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Klucz „${wartosci.nazwa ?? ''}" wciągnięty do wykazu.`);
}

async function usunKlucz({ kanal, korzen }: Otoczenie, panel: string): Promise<void> {
  const wykaz = await wywolaj(kanal, Command.TerminalKeyList, {});
  const klucze = wykaz.wynik?.keys ?? [];
  if (klucze.length === 0) {
    oglos(NAGLOWEK, 'Wykaz kluczy jest pusty.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Zdjęcie klucza z wykazu',
    pola: [
      {
        klucz: 'klucz',
        etykieta: 'Klucz',
        wybor: klucze.map((klucz) => [klucz.id, `${klucz.name} · ${klucz.fingerprint}`] as const),
      },
      {
        klucz: 'pliki',
        etykieta: 'Pliki klucza na dysku rdzenia',
        wybor: [['zostaw', 'Zostaw nietknięte'], ['skasuj', 'Skasuj razem z wpisem']],
      },
    ],
    wykonanie: 'Zdejmij klucz',
    nieodwracalne: 'Wpisy książki gospodarzy wskazujące ten klucz stracą wskazanie. '
      + 'Skasowanych plików klucza nie da się odzyskać.',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.TerminalKeyRemove, {
    keyId: wartosci.klucz ?? '',
    deleteFiles: wartosci.pliki === 'skasuj',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił usunięcia klucza.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, wartosci.pliki === 'skasuj'
    ? 'Klucz zdjęty z wykazu wraz z plikami.'
    : 'Klucz zdjęty z wykazu; pliki zostały na dysku rdzenia.');
}

async function zapiszSkrypt({ kanal, korzen, odswiez }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Pozycja biblioteki skryptów',
    pola: [
      { klucz: 'nazwa', etykieta: 'Nazwa', podpowiedz: 'budowa rdzenia', wymagane: true },
      { klucz: 'powloka', etykieta: 'Powłoka', wybor: POWLOKI },
      {
        klucz: 'tresc',
        etykieta: 'Treść',
        podpowiedz: 'go build ./...',
        obszerne: true,
        wymagane: true,
      },
      { klucz: 'skrot', etykieta: 'Wywołanie skrócone' },
    ],
    wykonanie: 'Zapisz w bibliotece',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.TerminalScriptSave, {
    script: {
      id: '',
      name: wartosci.nazwa ?? '',
      kind: TerminalScriptKind.Script,
      shell: (wartosci.powloka ?? TerminalShell.Bash) as TerminalShell,
      content: wartosci.tresc ?? '',
      ...(wartosci.skrot === '' ? {} : { alias: wartosci.skrot }),
      version: 0,
      createdAt: 0,
      updatedAt: 0,
    },
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu pozycji.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Pozycja „${wartosci.nazwa ?? ''}" stoi w bibliotece.`);
  odswiez();
}

async function sprawdzSkrypt({ kanal, korzen }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Sprawdzenie treści skryptu',
    pola: [
      { klucz: 'powloka', etykieta: 'Powłoka', wybor: POWLOKI },
      { klucz: 'tresc', etykieta: 'Treść', obszerne: true, wymagane: true },
    ],
    wykonanie: 'Sprawdź treść',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.TerminalScriptLint, {
    content: wartosci.tresc ?? '',
    shell: (wartosci.powloka ?? TerminalShell.Bash) as TerminalShell,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił sprawdzenia skryptu.', 'ostrzezenie');
    return;
  }
  if (!wynik.wynik.analyzerAvailable) {
    oglos(NAGLOWEK,
      `Narzędzie ${wynik.wynik.analyzer} nie stoi na maszynie rdzenia — treści nikt nie sprawdził.`,
      'ostrzezenie');
    return;
  }
  const uwagi = wynik.wynik.findings.length;
  oglos(NAGLOWEK, uwagi === 0
    ? `Narzędzie ${wynik.wynik.analyzer} nie ma uwag do treści.`
    : `Narzędzie ${wynik.wynik.analyzer} zgłasza ${uwagi} uwag.`);
}

async function usunSkrypt({ kanal, korzen, odswiez }: Otoczenie, panel: string): Promise<void> {
  const wykaz = await wywolaj(kanal, Command.TerminalScriptList, {});
  const skrypty = wykaz.wynik?.scripts ?? [];
  if (skrypty.length === 0) {
    oglos(NAGLOWEK, 'Biblioteka skryptów jest pusta.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Usunięcie pozycji biblioteki',
    pola: [{
      klucz: 'skrypt',
      etykieta: 'Pozycja',
      wybor: skrypty.map((skrypt) => [skrypt.id, `${skrypt.name} · ${skrypt.shell}`] as const),
    }],
    wykonanie: 'Usuń pozycję',
    nieodwracalne: 'Pozycja znika wraz ze wszystkimi swoimi wersjami.',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.TerminalScriptRemove, {
    scriptId: wartosci.skrypt ?? '',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił usunięcia pozycji.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Pozycja zdjęta z biblioteki.');
  odswiez();
}

/* Tunel idzie przez wpis książki gospodarzy: bez gospodarza rdzeń nie ma przez
   co go poprowadzić, więc okno pyta o wpis, a nie o adres wpisywany z ręki. */
async function otworzTunel({ kanal, korzen, idOkna }: Otoczenie, panel: string): Promise<void> {
  const wykaz = await wywolaj(kanal, Command.TerminalHostList, {});
  const gospodarze = wykaz.wynik?.hosts ?? [];
  if (gospodarze.length === 0) {
    oglos(NAGLOWEK, 'Książka gospodarzy jest pusta — tunel nie ma przez co iść.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Otwarcie tunelu',
    pola: [
      {
        klucz: 'gospodarz',
        etykieta: 'Gospodarz',
        wybor: gospodarze.map((wpis) => [wpis.id, `${wpis.name} · ${wpis.target}`] as const),
      },
      {
        klucz: 'rodzaj',
        etykieta: 'Rodzaj przekierowania',
        wybor: [
          [TerminalTunnelKind.Dynamic, 'Dynamiczne'],
          [TerminalTunnelKind.Local, 'Miejscowe'],
          [TerminalTunnelKind.Remote, 'Zdalne'],
        ],
      },
      { klucz: 'port', etykieta: 'Port po stronie rdzenia', podpowiedz: 'Wolny wybierze rdzeń' },
    ],
    wykonanie: 'Otwórz tunel',
  });
  if (wartosci === null) return;
  const port = Number.parseInt(wartosci.port ?? '', 10);
  const wynik = await wywolaj(kanal, Command.TerminalTunnelOpen, {
    windowId: idOkna(),
    kind: (wartosci.rodzaj ?? TerminalTunnelKind.Dynamic) as TerminalTunnelKind,
    hostId: wartosci.gospodarz ?? '',
    ...(Number.isFinite(port) ? { localPort: port } : {}),
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił otwarcia tunelu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Tunel otwarty.');
}
