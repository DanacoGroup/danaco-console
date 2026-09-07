// Sekcja modeli w oknie Ustawień: kanały modeli, którymi okna komunikacji
// rozmawiają z dostawcą — dołożenie, poprawa i zdjęcie kanału.
import { Command } from '../../../shared/contract.ts';
import type { Channel } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import {
  grupaKatalogu,
  komorka,
  komorkaCzynnosci,
  zdanieWiersza,
} from './ustawienia-katalog.ts';
import { zapytajWOknieModalnym } from './pytanie-modalne.ts';

const NAGLOWEK = 'Kanały modeli';

const RODZAJE: ReadonlyArray<readonly [string, string]> = [
  ['cli', 'Narzędzie wiersza poleceń'],
  ['api', 'Usługa sieciowa'],
];

/*
zwiazKanalyUstawien dokłada do sekcji modeli katalog kanałów. Konto mówi, czym
Operator płaci za rozmowę; kanał mówi, którędy ta rozmowa idzie i jakim modelem
— dlatego oba katalogi stoją w jednej sekcji, jeden pod drugim.
*/
export function zwiazKanalyUstawien(kanal: Kanal, korzen: Element): void {
  const sekcja = korzen.querySelector('#sekcja-modele');
  if (sekcja === null) return;
  sekcja.append(grupaKatalogu(sekcja.ownerDocument, {
    kod: 'us-kanaly',
    tytul: 'Kanały modeli',
    opis: 'Kanał wiąże okno komunikacji z modelem dostawcy. Okno bez czynnego '
      + 'kanału nie ma czym rozmawiać.',
    napisPrzycisku: 'Dołóż kanał',
    kolumny: ['Kanał', 'Rodzaj', 'Model', 'Stan', 'Akcje'],
  }));

  void wypelnij(kanal, sekcja);

  sekcja.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    if (cel.closest('#us-kanaly-dodaj') !== null) {
      zdarzenie.stopPropagation();
      void doloz(kanal, sekcja);
      return;
    }
    const wiersz = cel.closest<HTMLElement>('[data-kanal-modelu]');
    const czynnosc = cel.closest<HTMLElement>('[data-kanal-czynnosc]');
    if (wiersz === null || czynnosc === null) return;
    zdarzenie.stopPropagation();
    void wykonaj(kanal, sekcja, czynnosc.dataset.kanalCzynnosc ?? '',
      wiersz.dataset.kanalModelu ?? '');
  });
}

async function wypelnij(kanal: Kanal, sekcja: Element): Promise<void> {
  const cialo = sekcja.querySelector('#us-kanaly-cialo');
  if (cialo === null) return;
  const wynik = await wywolaj(kanal, Command.ChannelList, { enabledOnly: false });
  if (!wynik.udany || wynik.wynik === undefined) {
    cialo.replaceChildren(zdanieWiersza(cialo, 5,
      wynik.blad?.message ?? 'Wykaz kanałów nie doszedł.'));
    return;
  }
  const kanaly = wynik.wynik.channels;
  if (kanaly.length === 0) {
    cialo.replaceChildren(zdanieWiersza(cialo, 5, 'Katalog nie ma jeszcze żadnego kanału.'));
    return;
  }
  cialo.replaceChildren(...kanaly.map((pozycja) => wiersz(cialo, pozycja)));
}

function wiersz(cialo: Element, pozycja: Channel): HTMLElement {
  const dokument = cialo.ownerDocument;
  const wezel = dokument.createElement('tr');
  wezel.dataset.kanalModelu = pozycja.id;
  const rodzaj = RODZAJE.find((para) => para[0] === pozycja.kind);
  wezel.append(
    komorka(dokument, pozycja.name),
    komorka(dokument, rodzaj === undefined ? pozycja.kind : rodzaj[1]),
    komorka(dokument, pozycja.model ?? 'bez wskazania modelu', 'pt-mono'),
    plakietka(dokument, pozycja),
    komorkaCzynnosci(dokument, [
      ['popraw', 'Popraw'],
      ['sprawdz', 'Sprawdź'],
      ['zdejmij', 'Zdejmij'],
    ], 'data-kanal-czynnosc'),
  );
  return wezel;
}

function plakietka(dokument: Document, pozycja: Channel): HTMLElement {
  const wezel = dokument.createElement('td');
  const znak = dokument.createElement('span');
  znak.className = pozycja.enabled ? 'dn-plakietka dn-plakietka--sukces' : 'dn-plakietka';
  znak.textContent = pozycja.enabled ? 'czynny' : 'wyłączony';
  wezel.appendChild(znak);
  return wezel;
}

async function wykonaj(
  kanal: Kanal,
  sekcja: Element,
  czynnosc: string,
  identyfikator: string,
): Promise<void> {
  if (czynnosc === 'popraw') return popraw(kanal, sekcja, identyfikator);
  if (czynnosc === 'sprawdz') return sprawdz(kanal, identyfikator);
  if (czynnosc === 'zdejmij') return zdejmij(kanal, sekcja, identyfikator);
  /* Czynność spoza obsłużonych odmawia zamiast milczeć. */
  oglos(NAGLOWEK, `Czynność „${czynnosc}" nie jest prowadzona przez to okno.`, 'ostrzezenie');
}

async function pobierz(kanal: Kanal, identyfikator: string): Promise<Channel | undefined> {
  const wykaz = await wywolaj(kanal, Command.ChannelList, { enabledOnly: false });
  return wykaz.wynik?.channels.find((pozycja) => pozycja.id === identyfikator);
}

async function doloz(kanal: Kanal, sekcja: Element): Promise<void> {
  const wartosci = await zapytajWOknieModalnym({
    tytul: 'Nowy kanał modelu',
    opis: 'Rodzaj rozstrzyga, czym rdzeń woła model: narzędziem czy usługą sieciową.',
    pola: [
      { klucz: 'nazwa', etykieta: 'Nazwa kanału', wymagane: true },
      { klucz: 'rodzaj', etykieta: 'Rodzaj kanału', wybor: RODZAJE },
      { klucz: 'model', etykieta: 'Model' },
    ],
    wykonanie: 'Dołóż kanał',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.ChannelAdd, {
    name: wartosci.nazwa ?? '',
    kind: wartosci.rodzaj ?? 'cli',
    ...(wartosci.model === '' ? {} : { model: wartosci.model }),
    enabled: true,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił dołożenia kanału.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Kanał „${wartosci.nazwa ?? ''}" stoi w katalogu.`);
  await wypelnij(kanal, sekcja);
}

async function popraw(kanal: Kanal, sekcja: Element, identyfikator: string): Promise<void> {
  const pozycja = await pobierz(kanal, identyfikator);
  if (pozycja === undefined) {
    oglos(NAGLOWEK, 'Katalog nie ma już tego kanału.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWOknieModalnym({
    tytul: `Poprawa kanału „${pozycja.name}"`,
    pola: [
      { klucz: 'nazwa', etykieta: 'Nazwa kanału', wartosc: pozycja.name },
      { klucz: 'model', etykieta: 'Model', wartosc: pozycja.model ?? '' },
      {
        klucz: 'czynny',
        etykieta: 'Stan kanału',
        wybor: [['tak', 'Czynny'], ['nie', 'Wyłączony']],
        wartosc: pozycja.enabled ? 'tak' : 'nie',
      },
    ],
    wykonanie: 'Zapisz poprawki',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.ChannelUpdate, {
    channelId: identyfikator,
    name: wartosci.nazwa ?? '',
    model: wartosci.model ?? '',
    enabled: wartosci.czynny !== 'nie',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił poprawy kanału.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Kanał „${wartosci.nazwa ?? ''}" poprawiony.`);
  await wypelnij(kanal, sekcja);
}

/* Sprawdzenie mówi o dwóch rzeczach naraz: czy kanał odpowiada i czy rdzeń ma
   dla niego poświadczenie — kanał bez poświadczenia odpowie odmową. */
async function sprawdz(kanal: Kanal, identyfikator: string): Promise<void> {
  /* Sprawdzenie sięga do dostawcy i potrafi trwać, więc okno mówi o tym od
     razu — inaczej naciśnięcie wygląda na nieprzyjęte. */
  oglos(NAGLOWEK, 'Pytam kanał o odpowiedź…');
  const wynik = await wywolaj(kanal, Command.ChannelCheck, { channelId: identyfikator });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił sprawdzenia kanału.', 'ostrzezenie');
    return;
  }
  const poswiadczenie = await wywolaj(kanal, Command.ChannelCredentialStatus, {
    channelId: identyfikator,
  });
  const stanPoswiadczenia = poswiadczenie.wynik?.status.present === true
    ? 'poświadczenie stoi'
    : 'bez poświadczenia';
  const opoznienie = wynik.wynik.latencyMs;
  oglos(NAGLOWEK, (wynik.wynik.reachable ? 'Kanał odpowiada' : 'Kanał nie odpowiada')
    + (opoznienie === undefined ? '' : ` · ${String(opoznienie)} ms`)
    + ` · ${stanPoswiadczenia}`
    + (wynik.wynik.detail === undefined || wynik.wynik.detail === ''
      ? '' : ` · ${wynik.wynik.detail}`));
}

async function zdejmij(kanal: Kanal, sekcja: Element, identyfikator: string): Promise<void> {
  const pozycja = await pobierz(kanal, identyfikator);
  if (pozycja === undefined) {
    oglos(NAGLOWEK, 'Katalog nie ma już tego kanału.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWOknieModalnym({
    tytul: `Zdjęcie kanału „${pozycja.name}"`,
    pola: [],
    wykonanie: 'Zdejmij kanał',
    nieodwracalne: 'Okna, które rozmawiały tym kanałem, stracą drogę do modelu.',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.ChannelRemove, { channelId: identyfikator });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zdjęcia kanału.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Kanał „${pozycja.name}" zdjęty z katalogu.`);
  await wypelnij(kanal, sekcja);
}
