// Poczta w oknie Assistant: konta, foldery i wiadomości, którymi asystent
// pracuje — panel zadań pokazuje podsumowanie poczty jako jedną ze swoich prac.
import { Command, MailProtocol } from '../../../shared/contract.ts';
import type { MailAccount, MailMessage } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { dolozCzynnosciPanelu, zapytajWSzufladzie } from './czynnosci-okna.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';

const NAGLOWEK = 'Poczta';
const PANEL = 'panel-zadania';

const PROTOKOLY: ReadonlyArray<readonly [string, string]> = [
  [MailProtocol.Imap, 'IMAP — działa z każdą skrzynką'],
  [MailProtocol.Jmap, 'JMAP — szybszy przy dużych skrzynkach'],
  [MailProtocol.Pop3, 'POP3 — tylko odbiór, bez folderów'],
];

interface Otoczenie {
  kanal: Kanal;
  korzen: Element;
}

export function zwiazPoczteAsystenta(
  kanal: Kanal,
  korzen: Element,
  dolacz: (zdejmij: () => void) => void,
  przy: AddEventListenerOptions,
): void {
  const otoczenie: Otoczenie = { kanal, korzen };
  const zdejmij = dolozCzynnosciPanelu(korzen, PANEL, 'Poczta', [
    {
      naglowek: 'Konta poczty',
      pozycje: [
        { kod: 'konta', nazwa: 'Wykaz kont' },
        { kod: 'odkryj', nazwa: 'Odkryj konta platformy' },
        { kod: 'dodaj', nazwa: 'Dodaj konto…' },
        { kod: 'usun', nazwa: 'Usuń konto…' },
      ],
    },
    {
      naglowek: 'Skrzynka',
      pozycje: [
        { kod: 'foldery', nazwa: 'Wykaz folderów' },
        { kod: 'wiadomosci', nazwa: 'Wiadomości folderu…' },
        { kod: 'wiadomosc', nazwa: 'Otwórz wiadomość…' },
        { kod: 'oznacz', nazwa: 'Oznacz wiadomość…' },
      ],
    },
    {
      naglowek: 'Pisanie',
      pozycje: [
        { kod: 'szkic', nazwa: 'Zapisz szkic…' },
        { kod: 'wyslij', nazwa: 'Wyślij wiadomość…' },
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
  if (kod === 'konta') return wykazKont(otoczenie);
  if (kod === 'odkryj') return odkryjKonta(otoczenie);
  if (kod === 'dodaj') return dodajKonto(otoczenie);
  if (kod === 'usun') return usunKonto(otoczenie);
  if (kod === 'foldery') return wykazFolderow(otoczenie);
  if (kod === 'wiadomosci') return wykazWiadomosci(otoczenie);
  if (kod === 'wiadomosc') return otworzWiadomosc(otoczenie);
  if (kod === 'oznacz') return oznaczWiadomosc(otoczenie);
  if (kod === 'szkic') return zapiszSzkic(otoczenie);
  if (kod === 'wyslij') return wyslijWiadomosc(otoczenie);
  oglos(NAGLOWEK, 'Ta czynność nie ma wiązania.', 'ostrzezenie');
}

async function pobierzKonta(otoczenie: Otoczenie): Promise<MailAccount[]> {
  const wynik = await wywolaj(otoczenie.kanal, Command.MailAccountList, {});
  return wynik.udany ? (wynik.wynik?.accounts ?? []) : [];
}

function zdanieKonta(konto: MailAccount): string {
  return `${konto.address} · ${konto.protocol}`
    + ` · ${konto.connected ? 'połączone' : 'bez połączenia'}`;
}

async function wykazKont(otoczenie: Otoczenie): Promise<void> {
  const wynik = await wywolaj(otoczenie.kanal, Command.MailAccountList, {});
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Wykaz kont nie doszedł.', 'ostrzezenie');
    return;
  }
  const konta = wynik.wynik.accounts;
  wypelnij(otoczenie.korzen, konta.length === 0
    ? ['Rdzeń nie prowadzi żadnego konta poczty.']
    : konta.map(zdanieKonta));
  oglos(NAGLOWEK, `Kont poczty: ${String(wynik.wynik.total)}.`);
}

/* Odkrywanie czyta konta, które platforma już zna, i niczego nie zakłada:
   dopisanie konta jest osobną czynnością, wykonywaną świadomie. */
async function odkryjKonta(otoczenie: Otoczenie): Promise<void> {
  const wynik = await wywolaj(otoczenie.kanal, Command.MailAccountDiscover, {});
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Odkrywanie kont nie doszło.',
      'ostrzezenie');
    return;
  }
  const konta = wynik.wynik.accounts;
  wypelnij(otoczenie.korzen, konta.length === 0
    ? ['Platforma nie podaje kont do odkrycia.']
    : konta.map(zdanieKonta));
  oglos(NAGLOWEK, `Kont odkrytych: ${String(wynik.wynik.total)}.`);
}

async function dodajKonto(otoczenie: Otoczenie): Promise<void> {
  const odpowiedz = await zapytajWSzufladzie(otoczenie.korzen, PANEL, {
    tytul: 'Dodanie konta poczty',
    pola: [
      { klucz: 'adres', etykieta: 'Adres', wymagane: true },
      { klucz: 'protokol', etykieta: 'Protokół', wybor: PROTOKOLY },
      { klucz: 'odbior', etykieta: 'Serwer odbioru' },
      { klucz: 'wysylka', etykieta: 'Serwer wysyłki' },
      { klucz: 'uzytkownik', etykieta: 'Użytkownik' },
      { klucz: 'nazwa', etykieta: 'Nazwa widoczna' },
    ],
    wykonanie: 'Dodaj konto',
  });
  if (odpowiedz === null) {
    oglos(NAGLOWEK, 'Dodanie konta odwołane; wykaz bez zmiany.');
    return;
  }
  const protokol = odpowiedz['protokol'] ?? '';
  const wynik = await wywolaj(otoczenie.kanal, Command.MailAccountAdd, {
    address: odpowiedz['adres'] ?? '',
    protocol: protokol === '' ? undefined : (protokol as typeof MailProtocol.Imap),
    incomingHost: pusteNaBrak(odpowiedz['odbior']),
    outgoingHost: pusteNaBrak(odpowiedz['wysylka']),
    username: pusteNaBrak(odpowiedz['uzytkownik']),
    displayName: pusteNaBrak(odpowiedz['nazwa']),
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił dodania konta.',
      'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Konto dodane: ${wynik.wynik?.account.address ?? 'bez adresu'}.`);
  await wykazKont(otoczenie);
}

/* Puste pole szuflady to brak wskazania, nie pusty ciąg: rdzeń ma wtedy
   dobrać wartość sam, zamiast zapisywać pustkę. */
function pusteNaBrak(wartosc: string | undefined): string | undefined {
  const tresc = (wartosc ?? '').trim();
  return tresc === '' ? undefined : tresc;
}

async function usunKonto(otoczenie: Otoczenie): Promise<void> {
  const konta = await pobierzKonta(otoczenie);
  if (konta.length === 0) {
    oglos(NAGLOWEK, 'Nie ma konta do usunięcia.');
    return;
  }
  const odpowiedz = await zapytajWSzufladzie(otoczenie.korzen, PANEL, {
    tytul: 'Usunięcie konta poczty',
    pola: [{
      klucz: 'konto',
      etykieta: 'Konto',
      wybor: konta.map((poz) => [poz.id, poz.address] as const),
      wymagane: true,
    }],
    wykonanie: 'Usuń konto',
    nieodwracalne: 'Konto znika z platformy wraz z jego ustawieniami połączenia;'
      + ' wiadomości na serwerze zostają nietknięte.',
  });
  if (odpowiedz === null) {
    oglos(NAGLOWEK, 'Usunięcie odwołane; konto zostaje.');
    return;
  }
  const wynik = await wywolaj(otoczenie.kanal, Command.MailAccountRemove, {
    accountId: odpowiedz['konto'] ?? '',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił usunięcia konta.',
      'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Konto usunięte z platformy.');
  await wykazKont(otoczenie);
}

async function wykazFolderow(otoczenie: Otoczenie): Promise<void> {
  const wynik = await wywolaj(otoczenie.kanal, Command.MailFolderList, {});
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Wykaz folderów nie doszedł.',
      'ostrzezenie');
    return;
  }
  const foldery = wynik.wynik.folders;
  wypelnij(otoczenie.korzen, foldery.length === 0
    ? ['Konto nie podaje folderów.']
    : foldery);
  oglos(NAGLOWEK, `Folderów: ${String(foldery.length)}.`);
}

function zdanieWiadomosci(wiadomosc: MailMessage): string {
  return `${wiadomosc.from} · ${wiadomosc.subject ?? 'bez tematu'}`
    + `${wiadomosc.unread === true ? ' · nieprzeczytana' : ''}`;
}

async function wykazWiadomosci(otoczenie: Otoczenie): Promise<void> {
  const odpowiedz = await zapytajWSzufladzie(otoczenie.korzen, PANEL, {
    tytul: 'Wiadomości folderu',
    pola: [
      { klucz: 'folder', etykieta: 'Folder', podpowiedz: 'puste znaczy domyślny' },
      { klucz: 'szukaj', etykieta: 'Szukany ciąg' },
    ],
    wykonanie: 'Pokaż wiadomości',
  });
  if (odpowiedz === null) {
    oglos(NAGLOWEK, 'Odczyt odwołany; panel bez zmiany.');
    return;
  }
  const wynik = await wywolaj(otoczenie.kanal, Command.MailMessageList, {
    folder: pusteNaBrak(odpowiedz['folder']),
    query: pusteNaBrak(odpowiedz['szukaj']),
    limit: 20,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Wiadomości nie doszły.', 'ostrzezenie');
    return;
  }
  const wiadomosci = wynik.wynik.messages;
  wypelnij(otoczenie.korzen, wiadomosci.length === 0
    ? ['Folder nie ma wiadomości pasujących do zapytania.']
    : wiadomosci.map(zdanieWiadomosci));
  oglos(NAGLOWEK, `Wiadomości: ${String(wynik.wynik.total)}.`);
}

async function otworzWiadomosc(otoczenie: Otoczenie): Promise<void> {
  const odpowiedz = await zapytajWSzufladzie(otoczenie.korzen, PANEL, {
    tytul: 'Otwarcie wiadomości',
    pola: [
      { klucz: 'wiadomosc', etykieta: 'Identyfikator wiadomości', wymagane: true },
    ],
    wykonanie: 'Otwórz',
  });
  if (odpowiedz === null) {
    oglos(NAGLOWEK, 'Otwarcie odwołane.');
    return;
  }
  const wynik = await wywolaj(otoczenie.kanal, Command.MailMessageGet, {
    messageId: odpowiedz['wiadomosc'] ?? '',
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Wiadomość nie doszła.', 'ostrzezenie');
    return;
  }
  const wiadomosc = wynik.wynik.message;
  wypelnij(otoczenie.korzen, [
    `od: ${wiadomosc.from}`,
    `temat: ${wiadomosc.subject ?? 'bez tematu'}`,
    `treść: ${String((wiadomosc.body ?? wiadomosc.preview ?? '').length)} znaków`,
    `załączników: ${String(wynik.wynik.attachmentAssetIds?.length ?? 0)}`,
  ]);
  oglos(NAGLOWEK, 'Wiadomość otwarta w panelu.');
}

async function oznaczWiadomosc(otoczenie: Otoczenie): Promise<void> {
  const odpowiedz = await zapytajWSzufladzie(otoczenie.korzen, PANEL, {
    tytul: 'Oznaczenie wiadomości',
    pola: [
      { klucz: 'wiadomosc', etykieta: 'Identyfikator wiadomości', wymagane: true },
      {
        klucz: 'stan',
        etykieta: 'Oznaczenie',
        wybor: [
          ['nieprzeczytana', 'jako nieprzeczytana'],
          ['przeczytana', 'jako przeczytana'],
          ['wyrozniona', 'wyróżniona'],
        ],
        wymagane: true,
      },
    ],
    wykonanie: 'Oznacz',
  });
  if (odpowiedz === null) {
    oglos(NAGLOWEK, 'Oznaczenie odwołane; wiadomość bez zmiany.');
    return;
  }
  const stan = odpowiedz['stan'] ?? '';
  const wynik = await wywolaj(otoczenie.kanal, Command.MailMessageFlag, {
    messageId: odpowiedz['wiadomosc'] ?? '',
    unread: stan === 'wyrozniona' ? undefined : stan === 'nieprzeczytana',
    flagged: stan === 'wyrozniona' ? true : undefined,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił oznaczenia.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Oznaczenie zapisane; wiadomość `
    + `${wynik.wynik?.message.unread === true ? 'nieprzeczytana' : 'przeczytana'}.`);
}

async function zapiszSzkic(otoczenie: Otoczenie): Promise<void> {
  const odpowiedz = await zapytajWSzufladzie(otoczenie.korzen, PANEL, {
    tytul: 'Zapis szkicu',
    pola: [
      { klucz: 'do', etykieta: 'Odbiorcy', podpowiedz: 'adresy po przecinku' },
      { klucz: 'temat', etykieta: 'Temat' },
      { klucz: 'tresc', etykieta: 'Treść', obszerne: true },
    ],
    wykonanie: 'Zapisz szkic',
  });
  if (odpowiedz === null) {
    oglos(NAGLOWEK, 'Zapis szkicu odwołany.');
    return;
  }
  const wynik = await wywolaj(otoczenie.kanal, Command.MailDraftSave, {
    to: adresy(odpowiedz['do']),
    subject: pusteNaBrak(odpowiedz['temat']),
    body: pusteNaBrak(odpowiedz['tresc']),
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu szkicu.',
      'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Szkic zapisany w folderze `
    + `${wynik.wynik?.draft.folder ?? 'bez wskazania'}.`);
}

function adresy(wartosc: string | undefined): string[] | undefined {
  const pozycje = (wartosc ?? '').split(',')
    .map((poz) => poz.trim())
    .filter((poz) => poz !== '');
  return pozycje.length === 0 ? undefined : pozycje;
}

/* Wysyłka jest nieodwracalna, bo wiadomość opuszcza platformę: dlatego
   szuflada żąda drugiego naciśnięcia i mówi wprost, co się stanie. */
async function wyslijWiadomosc(otoczenie: Otoczenie): Promise<void> {
  const odpowiedz = await zapytajWSzufladzie(otoczenie.korzen, PANEL, {
    tytul: 'Wysłanie wiadomości',
    pola: [
      { klucz: 'do', etykieta: 'Odbiorcy', podpowiedz: 'adresy po przecinku' },
      { klucz: 'temat', etykieta: 'Temat' },
      { klucz: 'tresc', etykieta: 'Treść', obszerne: true },
    ],
    wykonanie: 'Wyślij',
    nieodwracalne: 'Wiadomość opuszcza platformę i nie da się jej cofnąć'
      + ' ani po stronie rdzenia, ani po stronie odbiorcy.',
  });
  if (odpowiedz === null) {
    oglos(NAGLOWEK, 'Wysyłka odwołana; nic nie opuściło platformy.');
    return;
  }
  const wynik = await wywolaj(otoczenie.kanal, Command.MailSend, {
    to: adresy(odpowiedz['do']),
    subject: pusteNaBrak(odpowiedz['temat']),
    body: pusteNaBrak(odpowiedz['tresc']),
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wysyłki.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Wiadomość wysłana; rdzeń nadał jej `
    + `${wynik.wynik?.messageId ?? 'bez oznaczenia'}.`);
}
