// Sekcja modeli w oknie Ustawień: katalog kont dostawców — wykaz z rdzenia,
// dołożenie, poprawa, wskazanie domyślnego i zdjęcie konta.
import { AccountKind, Command } from '../../../shared/contract.ts';
import type { Account } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import { zapytajWOknieModalnym } from './pytanie-modalne.ts';

const NAGLOWEK = 'Konta modeli';

const RODZAJE: ReadonlyArray<readonly [string, string]> = [
  [AccountKind.Cli, 'Narzędzie wiersza poleceń'],
  [AccountKind.Api, 'Usługa sieciowa'],
];

/*
zwiazKontaUstawien wypełnia sekcję modeli wykazem kont i wiąże jej czynności.

Czynność dołożenia stoi w narzędziach sekcji — tam, gdzie prototyp postawił
przycisk „Dodaj konto". Czynności pojedynczego konta stoją w jego wierszu,
w kolumnie akcji, bo dotyczą tego konta, a nie całej sekcji.
*/
export function zwiazKontaUstawien(kanal: Kanal, korzen: Element): void {
  const sekcja = korzen.querySelector('#sekcja-modele');
  if (sekcja === null) return;
  void wypelnijKonta(kanal, sekcja);

  sekcja.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    if (cel.closest('#us-dodaj-konto') !== null) {
      zdarzenie.stopPropagation();
      void dolozKonto(kanal, sekcja);
      return;
    }
    const wiersz = cel.closest<HTMLElement>('[data-konto]');
    const czynnosc = cel.closest<HTMLElement>('[data-konto-czynnosc]');
    if (wiersz === null || czynnosc === null) return;
    zdarzenie.stopPropagation();
    void wykonaj(kanal, sekcja, czynnosc.dataset.kontoCzynnosc ?? '', wiersz.dataset.konto ?? '');
  });
}

async function wypelnijKonta(kanal: Kanal, sekcja: Element): Promise<void> {
  const cialo = sekcja.querySelector('.us-tabela tbody');
  if (cialo === null) return;
  const wynik = await wywolaj(kanal, Command.AccountList, {});
  if (!wynik.udany || wynik.wynik === undefined) {
    cialo.replaceChildren(wierszZdania(cialo,
      wynik.blad?.message ?? 'Wykaz kont nie doszedł.'));
    return;
  }
  const konta = wynik.wynik.accounts;
  if (konta.length === 0) {
    cialo.replaceChildren(wierszZdania(cialo, 'Katalog nie ma jeszcze żadnego konta.'));
    return;
  }
  cialo.replaceChildren(...konta.map((konto) => wierszKonta(cialo, konto)));
}

function wierszZdania(cialo: Element, zdanie: string): HTMLElement {
  const wiersz = cialo.ownerDocument.createElement('tr');
  const komorka = cialo.ownerDocument.createElement('td');
  komorka.colSpan = 5;
  komorka.className = 'us-opis--drobny';
  komorka.textContent = zdanie;
  wiersz.appendChild(komorka);
  return wiersz;
}

function wierszKonta(cialo: Element, konto: Account): HTMLElement {
  const dokument = cialo.ownerDocument;
  const wiersz = dokument.createElement('tr');
  wiersz.dataset.konto = konto.id;
  wiersz.append(
    komorka(dokument, konto.name),
    komorka(dokument, konto.provider, 'pt-mono'),
    komorka(dokument, konto.defaultModel ?? 'bez wskazania modelu'),
    plakietka(dokument, konto),
    komorkaCzynnosci(dokument, konto),
  );
  return wiersz;
}

function komorka(dokument: Document, tresc: string, klasa = ''): HTMLElement {
  const wezel = dokument.createElement('td');
  if (klasa !== '') wezel.className = klasa;
  wezel.textContent = tresc;
  return wezel;
}

/* Stan konta niesie dwie rzeczy naraz: czy jest czynne i czy rdzeń ma dla niego
   poświadczenie — bez poświadczenia konto czynne i tak nie zadziała. */
function plakietka(dokument: Document, konto: Account): HTMLElement {
  const wezel = dokument.createElement('td');
  const znak = dokument.createElement('span');
  znak.className = konto.enabled && konto.hasCredential
    ? 'dn-plakietka dn-plakietka--sukces'
    : 'dn-plakietka';
  znak.textContent = konto.enabled
    ? (konto.hasCredential ? 'czynne' : 'czynne bez poświadczenia')
    : 'wyłączone';
  wezel.appendChild(znak);
  if (konto.isDefault) {
    const domyslne = dokument.createElement('span');
    domyslne.className = 'dn-plakietka dn-plakietka--sygnal';
    domyslne.textContent = 'domyślne';
    wezel.appendChild(domyslne);
  }
  return wezel;
}

function komorkaCzynnosci(dokument: Document, konto: Account): HTMLElement {
  const wezel = dokument.createElement('td');
  const gniazdo = dokument.createElement('span');
  gniazdo.className = 'us-akcje-komorka';
  gniazdo.append(
    przycisk(dokument, 'popraw', 'Popraw'),
    ...(konto.isDefault ? [] : [przycisk(dokument, 'domyslne', 'Ustaw domyślnym')]),
    przycisk(dokument, 'zdejmij', 'Zdejmij'),
  );
  wezel.appendChild(gniazdo);
  return wezel;
}

function przycisk(dokument: Document, czynnosc: string, napis: string): HTMLButtonElement {
  const wezel = dokument.createElement('button');
  wezel.type = 'button';
  wezel.className = 'dn-btn dn-btn--duch dn-btn--sm';
  wezel.dataset.kontoCzynnosc = czynnosc;
  wezel.textContent = napis;
  return wezel;
}

async function wykonaj(
  kanal: Kanal,
  sekcja: Element,
  czynnosc: string,
  idKonta: string,
): Promise<void> {
  if (czynnosc === 'popraw') return popraw(kanal, sekcja, idKonta);
  if (czynnosc === 'domyslne') return ustawDomyslne(kanal, sekcja, idKonta);
  if (czynnosc === 'zdejmij') return zdejmij(kanal, sekcja, idKonta);
  /* Czynność spoza obsłużonych odmawia zamiast milczeć: przycisk bez gałęzi
     wyglądałby jak działający. */
  oglos(NAGLOWEK, `Czynność „${czynnosc}" nie jest prowadzona przez to okno.`, 'ostrzezenie');
}

/* Poświadczenie idzie do sejfu rdzenia i nie wraca: okno nie pokazuje go
   nigdy, a puste pole zostawia poświadczenie stojące nietknięte. */
async function dolozKonto(kanal: Kanal, sekcja: Element): Promise<void> {
  const wartosci = await zapytajWOknieModalnym({
    tytul: 'Nowe konto modelu',
    opis: 'Poświadczenie idzie do sejfu rdzenia; okno nie pokaże go ponownie.',
    pola: [
      { klucz: 'nazwa', etykieta: 'Nazwa konta', wymagane: true },
      { klucz: 'rodzaj', etykieta: 'Rodzaj konta', wybor: RODZAJE },
      { klucz: 'dostawca', etykieta: 'Dostawca', podpowiedz: 'anthropic', wymagane: true },
      { klucz: 'model', etykieta: 'Model domyślny' },
      { klucz: 'katalog', etykieta: 'Katalog profilu' },
      { klucz: 'adres', etykieta: 'Adres usługi' },
      { klucz: 'poswiadczenie', etykieta: 'Poświadczenie' },
    ],
    wykonanie: 'Dołóż konto',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.AccountAdd, {
    name: wartosci.nazwa ?? '',
    kind: (wartosci.rodzaj ?? AccountKind.Cli) as AccountKind,
    provider: wartosci.dostawca ?? '',
    ...(wartosci.model === '' ? {} : { defaultModel: wartosci.model }),
    ...(wartosci.katalog === '' ? {} : { configDir: wartosci.katalog }),
    ...(wartosci.adres === '' ? {} : { baseUrl: wartosci.adres }),
    ...(wartosci.poswiadczenie === '' ? {} : { credential: wartosci.poswiadczenie }),
    enabled: true,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił dołożenia konta.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Konto „${wartosci.nazwa ?? ''}" stoi w katalogu.`);
  await wypelnijKonta(kanal, sekcja);
}

async function popraw(kanal: Kanal, sekcja: Element, idKonta: string): Promise<void> {
  const wykaz = await wywolaj(kanal, Command.AccountList, {});
  const konto = wykaz.wynik?.accounts.find((pozycja) => pozycja.id === idKonta);
  if (konto === undefined) {
    oglos(NAGLOWEK, 'Katalog nie ma już tego konta.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWOknieModalnym({
    tytul: `Poprawa konta „${konto.name}"`,
    opis: 'Puste pole poświadczenia zostawia poświadczenie stojące nietknięte.',
    pola: [
      { klucz: 'nazwa', etykieta: 'Nazwa konta', wartosc: konto.name },
      { klucz: 'dostawca', etykieta: 'Dostawca', wartosc: konto.provider },
      { klucz: 'model', etykieta: 'Model domyślny', wartosc: konto.defaultModel ?? '' },
      { klucz: 'katalog', etykieta: 'Katalog profilu', wartosc: konto.configDir ?? '' },
      { klucz: 'adres', etykieta: 'Adres usługi', wartosc: konto.baseUrl ?? '' },
      { klucz: 'poswiadczenie', etykieta: 'Nowe poświadczenie' },
      {
        klucz: 'czynne',
        etykieta: 'Stan konta',
        wybor: [['tak', 'Czynne'], ['nie', 'Wyłączone']],
        wartosc: konto.enabled ? 'tak' : 'nie',
      },
    ],
    wykonanie: 'Zapisz poprawki',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.AccountUpdate, {
    accountId: idKonta,
    name: wartosci.nazwa ?? '',
    provider: wartosci.dostawca ?? '',
    defaultModel: wartosci.model ?? '',
    configDir: wartosci.katalog ?? '',
    baseUrl: wartosci.adres ?? '',
    ...(wartosci.poswiadczenie === '' ? {} : { credential: wartosci.poswiadczenie }),
    enabled: wartosci.czynne !== 'nie',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił poprawy konta.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Konto „${wartosci.nazwa ?? ''}" poprawione.`);
  await wypelnijKonta(kanal, sekcja);
}

async function ustawDomyslne(kanal: Kanal, sekcja: Element, idKonta: string): Promise<void> {
  const wykaz = await wywolaj(kanal, Command.AccountList, {});
  const konto = wykaz.wynik?.accounts.find((pozycja) => pozycja.id === idKonta);
  if (konto === undefined) {
    oglos(NAGLOWEK, 'Katalog nie ma już tego konta.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.AccountDefaultSet, {
    accountId: idKonta,
    kind: konto.kind,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wskazania domyślnego konta.',
      'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Konto „${konto.name}" jest teraz domyślne dla swojego rodzaju.`);
  await wypelnijKonta(kanal, sekcja);
}

async function zdejmij(kanal: Kanal, sekcja: Element, idKonta: string): Promise<void> {
  const wykaz = await wywolaj(kanal, Command.AccountList, {});
  const konto = wykaz.wynik?.accounts.find((pozycja) => pozycja.id === idKonta);
  if (konto === undefined) {
    oglos(NAGLOWEK, 'Katalog nie ma już tego konta.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWOknieModalnym({
    tytul: `Zdjęcie konta „${konto.name}"`,
    pola: [],
    wykonanie: 'Zdejmij konto',
    nieodwracalne: 'Konto znika z katalogu wraz z poświadczeniem w sejfie rdzenia.',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.AccountRemove, { accountId: idKonta });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zdjęcia konta.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Konto „${konto.name}" zdjęte z katalogu.`);
  await wypelnijKonta(kanal, sekcja);
}
