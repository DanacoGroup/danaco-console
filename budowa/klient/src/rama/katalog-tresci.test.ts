/**
 * Rama aplikacji — katalog treści jest jedynym miejscem z tekstem. Sprawdzian
 * pilnuje, że cały tekst widoczny dla użytkownika stoi w jednym pliku, a poza
 * nim żaden łańcuch nie czyta się jak zdanie. Ten sam instrument co w drodze
 * wejścia, uruchomiony na własnym katalogu tego terenu.
 */

import { bieg, rowne, sprawdz } from '../sprawdzian.ts';

/* Odczyt plików bez typów środowiska: klient nie zaciąga deklaracji Node,
   a specyfikator spoza literału zostawia moduł nieopisanym. */
const nazwaModulu = 'node:fs';
const pliki = (await import(nazwaModulu)) as unknown as {
  readdirSync(
    sciezka: string,
    opcje: { withFileTypes: true },
  ): { name: string; isDirectory(): boolean }[];
  readFileSync(sciezka: string, kodowanie: string): string;
};

/** Katalog terenu; sprawdzian stoi w nim, więc odczyt plików idzie od własnego położenia pliku sprawdzianu. */
const KORZEN = new URL('.', import.meta.url).pathname;

/** Plik, który tekst niesie — jedyny plik w ramie, w którym wolno zapisać zdanie widoczne dla użytkownika. */
const KATALOG_TRESCI = 'tresci.ts';

/** Plik, który niesie rysunki; jego łańcuchy muszą być rysunkami, a nie zdaniami czytelnymi jak zwykły tekst. */
const ZESTAW_ZNAKOW = 'ikony.ts';

/** Przyrostek plików sprawdzianów — instrumentów mierzących ramę, a nie samej ramy, pomijanych przez tę regułę. */
const SPRAWDZIAN = '.test.ts';

/** Jedyny plik, który ten teren dołożył poza katalogiem ramy — bez niego straż nie sięga do korzenia `src/`. */
const PLIK_U_KORZENIA = '../aplikacja.ts';

/** Polskie znaki diakrytyczne — po ich obecności w łańcuchu poznaje się zdanie napisane po polsku wprost. */
const DIAKRYTYKI = /[ąćęłńóśźżĄĆĘŁŃÓŚŹŻ]/;

/**
 * Znak właściwy nazwie klasy albo selektorowi w arkuszu stylu; wyraz zdania
 * po polsku nigdy go nie niesie.
 */
const ZNAK_SELEKTORA = /[-.#[\]=,>+~*]/;

/**
 * Czy łańcuch ma kształt techniczny, czyli nie da się go przeczytać jak
 * zdania. Rozstrzyga próba na członach: w wykazie klas i w selektorze każdy
 * człon niesie znak selektora albo jest nazwą elementu, czyli ma najwyżej
 * dwie litery.
 */
function ksztaltTechniczny(lancuch: Lancuch, wZestawieZnakow: boolean): boolean {
  const tresc = lancuch.tresc;
  if (!tresc.includes(' ')) return true;
  if (tresc.trim().length === 0) return true;
  if (tresc.includes('${')) return true;
  if (wZestawieZnakow) return tresc.startsWith('<svg');
  return tresc
    .split(' ')
    .filter((czlon) => czlon.length > 0)
    .every((czlon) => ZNAK_SELEKTORA.test(czlon) || czlon.length <= 2);
}

/** Znalezisko: jeden łańcuch wraz z miejscem, w którym stoi — nazwą pliku i numerem jego wiersza w tym pliku. */
interface Lancuch {
  plik: string;
  wiersz: number;
  tresc: string;
  /** Czy łańcuch idzie do dziennika wywołaniem `console.*`. */
  doDziennika: boolean;
}

/** Wykaz plików terenu, wraz z podkatalogami, zebrany rekurencyjnie w kolejności odczytu katalogu na dysku. */
function wykazPlikow(katalog: string, przedrostek = ''): string[] {
  const znalezione: string[] = [];
  for (const wpis of pliki.readdirSync(katalog, { withFileTypes: true })) {
    const sciezka = `${przedrostek}${wpis.name}`;
    if (wpis.isDirectory()) znalezione.push(...wykazPlikow(`${katalog}${wpis.name}/`, `${sciezka}/`));
    else if (wpis.name.endsWith('.ts')) znalezione.push(sciezka);
  }
  return znalezione.sort();
}

/**
 * Wydobywa łańcuchy z pliku, pomijając komentarze i wyrażenia regularne, aby
 * pomiar mierzył wyłącznie tekst zdań, nie treść komentarzy ani znaki klas
 * wyrażeń.
 */
function lancuchyPliku(nazwa: string, tresc: string): Lancuch[] {
  const znalezione: Lancuch[] = [];
  let wiersz = 1;
  let i = 0;
  /* Ostatni znaczący znak rozstrzyga, czy ukośnik zaczyna wyrażenie regularne,
     czy jest dzieleniem. Po nawiasie zamykającym albo nazwie jest dzieleniem. */
  let ostatniZnaczacy = '';
  /* Początek wiersza bieżącego — po nim poznaje się wpis do dziennika. */
  let poczatekWiersza = 0;

  while (i < tresc.length) {
    const znak = tresc[i]!;
    const nastepny = tresc[i + 1];

    if (znak === '\n') {
      wiersz += 1;
      i += 1;
      poczatekWiersza = i;
      continue;
    }
    if (znak === '/' && nastepny === '/') {
      while (i < tresc.length && tresc[i] !== '\n') i += 1;
      continue;
    }
    if (znak === '/' && nastepny === '*') {
      i += 2;
      while (i < tresc.length && !(tresc[i] === '*' && tresc[i + 1] === '/')) {
        if (tresc[i] === '\n') wiersz += 1;
        i += 1;
      }
      i += 2;
      continue;
    }
    if (znak === '/' && ostatniZnaczacy !== '' && '=(,:[!&|?{};+-*%~^<>'.includes(ostatniZnaczacy)) {
      i += 1;
      let wKlasie = false;
      while (i < tresc.length) {
        const z = tresc[i]!;
        if (z === '\\') i += 2;
        else if (z === '[') {
          wKlasie = true;
          i += 1;
        } else if (z === ']') {
          wKlasie = false;
          i += 1;
        } else if (z === '/' && !wKlasie) {
          i += 1;
          break;
        } else if (z === '\n') break;
        else i += 1;
      }
      ostatniZnaczacy = '/';
      continue;
    }
    if (znak === "'" || znak === '"' || znak === '`') {
      const otwarcie = znak;
      const wiersOtwarcia = wiersz;
      let zebrane = '';
      i += 1;
      while (i < tresc.length) {
        const z = tresc[i]!;
        if (z === '\\') {
          zebrane += tresc[i + 1] ?? '';
          i += 2;
          continue;
        }
        if (z === otwarcie) {
          i += 1;
          break;
        }
        if (z === '\n') wiersz += 1;
        zebrane += z;
        i += 1;
      }
      znalezione.push({
        plik: nazwa,
        wiersz: wiersOtwarcia,
        tresc: zebrane,
        doDziennika: /console\.|new Error\(/.test(tresc.slice(poczatekWiersza, i)),
      });
      ostatniZnaczacy = otwarcie;
      continue;
    }
    if (znak.trim().length > 0) ostatniZnaczacy = znak;
    i += 1;
  }
  return znalezione;
}

const wykaz = [...wykazPlikow(KORZEN), PLIK_U_KORZENIA];
const plikiRamy = wykaz.filter((nazwa) => !nazwa.endsWith(SPRAWDZIAN));
const pominiete = wykaz.filter((nazwa) => nazwa.endsWith(SPRAWDZIAN));
const wszystkie: Lancuch[] = plikiRamy.flatMap((nazwa) =>
  lancuchyPliku(nazwa, pliki.readFileSync(`${KORZEN}${nazwa}`, 'utf8')),
);
const doDziennika = wszystkie.filter((l) => l.doDziennika);
const pozaKatalogiem = wszystkie.filter(
  (l) => l.plik !== KATALOG_TRESCI && !l.doDziennika,
);

function wypisz(nazwa: string, wartosc: unknown): void {
  console.log(`           ${nazwa}: ${JSON.stringify(wartosc)}`);
}

await bieg('rama aplikacji — katalog treści', {
  'instrument zmierzył to, co miał zmierzyć'() {
    // Zapora przed wynikiem pustym: zero trafień w katalogu nieistniejącym wygląda jak zero naruszeń.
    sprawdz(wykaz.length > 0, `nie znaleziono ani jednego pliku w ${KORZEN}`);
    sprawdz(plikiRamy.includes(KATALOG_TRESCI), `nie znaleziono katalogu ${KATALOG_TRESCI}`);
    sprawdz(plikiRamy.includes(ZESTAW_ZNAKOW), `nie znaleziono zestawu znaków ${ZESTAW_ZNAKOW}`);
    sprawdz(plikiRamy.includes(PLIK_U_KORZENIA), `nie znaleziono pliku u korzenia ${PLIK_U_KORZENIA}`);
    sprawdz(wszystkie.length > 0, 'nie znaleziono ani jednego łańcucha — wzorzec nic nie łapie');
    sprawdz(pozaKatalogiem.length > 0, 'poza katalogiem nie znaleziono ani jednego łańcucha');
    wypisz('plików ramy', plikiRamy.length);
    wypisz('sprawdzianów pominiętych', pominiete);
    wypisz('łańcuchów w plikach ramy', wszystkie.length);
    wypisz('łańcuchów poza katalogiem treści', pozaKatalogiem.length);
    wypisz('łańcuchów oddanych do dziennika', doDziennika.map((l) => `${l.plik}:${l.wiersz}`));
  },

  'instrument rozpoznaje zdanie po polsku, gdy je zobaczy'() {
    // Kontrola dodatnia: bez niej sprawdzian milczałby, gdyby wydobywanie łańcuchów przestało działać.
    const probka = lancuchyPliku('probka.ts', "const a = 'Środowisko nie zostało jeszcze rozpoznane.';\n");
    rowne(probka.length, 1, 'próbka dała jeden łańcuch');
    sprawdz(DIAKRYTYKI.test(probka[0]!.tresc), 'wzorzec nie rozpoznał polskiego zdania');
    const wKomentarzu = lancuchyPliku('probka.ts', "// 'Środowisko w komentarzu'\nconst a = 1;\n");
    rowne(wKomentarzu.length, 0, 'łańcuch w komentarzu został pominięty');
    const zWyrazeniem = lancuchyPliku('probka.ts', "const a = /[a-z']/.test(b);\nconst c = 'znak';\n");
    rowne(zWyrazeniem.length, 1, 'wyrażenie regularne nie połknęło reszty pliku');
  },

  'poza katalogiem treści nie ma ani jednego polskiego znaku w łańcuchu'() {
    const naruszenia = pozaKatalogiem.filter((l) => DIAKRYTYKI.test(l.tresc));
    for (const n of naruszenia) console.log(`           ${n.plik}:${n.wiersz} ${n.tresc}`);
    rowne(naruszenia.length, 0, 'łańcuchy z polskim znakiem poza katalogiem treści');
  },

  'poza katalogiem treści każdy łańcuch ma kształt techniczny'() {
    const naruszenia = pozaKatalogiem.filter(
      (l) => !ksztaltTechniczny(l, l.plik === ZESTAW_ZNAKOW),
    );
    for (const n of naruszenia) console.log(`           ${n.plik}:${n.wiersz} ${n.tresc}`);
    rowne(naruszenia.length, 0, 'łańcuchy o kształcie zdania poza katalogiem treści');
  },

  'instrument rozpoznaje zdanie, gdy stoi w składniku'() {
    // Kontrola dodatnia dla reguły kształtu: bez niej reguła mogłaby przepuszczać wszystko i milczeć.
    const zdanie: Lancuch = {
      plik: 'skladniki/probka.ts',
      wiersz: 1,
      tresc: 'Wybierz moduł z listy poniżej.',
      doDziennika: false,
    };
    sprawdz(!ksztaltTechniczny(zdanie, false), 'zdanie w składniku nie zostało rozpoznane');
    const klasy: Lancuch = {
      plik: 'skladniki/probka.ts',
      wiersz: 2,
      tresc: 'dn-szyna-poz dn-szyna-poz--modul',
      doDziennika: false,
    };
    sprawdz(ksztaltTechniczny(klasy, false), 'wykaz klas został wzięty za zdanie');
  },

  'wpis do dziennika nie przemyca treści dla użytkownika'() {
    // Odstępstwo jest wąskie: wpis diagnostyczny wolno pisać po polsku, ale nie zastępuje katalogu treści.
    sprawdz(
      doDziennika.length <= 2,
      `wpisów do dziennika przybyło: ${doDziennika.map((l) => `${l.plik}:${l.wiersz}`).join(', ')}`,
    );
    const bezNawiasu = doDziennika.filter((l) => !l.tresc.startsWith('['));
    for (const n of bezNawiasu) console.log(`           ${n.plik}:${n.wiersz} ${n.tresc}`);
    rowne(bezNawiasu.length, 0, 'wpisy do dziennika bez znacznika warstwy');
  },

  'każdy łańcuch zestawu znaków jest rysunkiem, nie zdaniem'() {
    const znaki = wszystkie.filter((l) => l.plik === ZESTAW_ZNAKOW);
    sprawdz(znaki.length > 0, 'zestaw znaków nie oddał ani jednego łańcucha');
    const naruszenia = znaki.filter((l) => !l.tresc.startsWith('<svg'));
    for (const n of naruszenia) console.log(`           ${n.plik}:${n.wiersz} ${n.tresc}`);
    rowne(naruszenia.length, 0, 'łańcuchy zestawu znaków, które nie są rysunkiem');
    wypisz('rysunków w zestawie znaków', znaki.length);
  },
});
