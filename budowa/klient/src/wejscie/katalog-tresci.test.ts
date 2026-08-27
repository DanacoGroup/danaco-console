/**
 * DROGA WEJŚCIA — katalog treści jest jedynym miejscem z tekstem.
 *
 * Sprawdzian pilnuje właściwości, którą prototyp przyniósł i którą teren ma
 * przenieść: cały tekst widoczny dla użytkownika stoi w jednym pliku, a poza
 * nim nie ma ani jednego łańcucha, który dałoby się przeczytać jako zdanie.
 * Łańcuch dopisany poza `tresci.ts` jest usterką i ma tu upaść.
 *
 * Reguła obejmuje PLIKI OKNA — składniki, ekrany, montaż, przebieg, narzędzia
 * i zestaw znaków. Nie obejmuje sprawdzianów: opis sprawdzianu jest zdaniem dla
 * tego, kto czyta wynik uruchomienia, i do okna nie trafia nigdy. Wykaz plików
 * pominiętych jest wypisywany, więc pominięcie nie da się rozrosnąć po cichu.
 *
 * Reguła jest dwuczłonowa i cała maszynowa:
 *
 *   1. Poza katalogiem żaden łańcuch nie niesie polskiego znaku diakrytycznego.
 *   2. Poza katalogiem każdy łańcuch ma KSZTAŁT TECHNICZNY: jest nazwą bez
 *      odstępu, wykazem klas albo selektorem, wzorem z podstawieniem, samym
 *      odstępem rozdzielającym węzły tekstowe, rysunkiem `<svg` w zestawie
 *      znaków albo wpisem diagnostycznym. Zdanie żadnego z tych kształtów nie
 *      ma i tu upada.
 *
 * Odstępstwo jest jedno i wąskie: wpis diagnostyczny — łańcuch oddany do
 * dziennika wywołaniem `console.*` albo niesiony wyjątkiem `new Error`. To
 * zdanie dla wykonawcy, nie dla Operatora; okno go nie pokazuje. Wpisy są
 * liczone i wypisywane, więc odstępstwo nie rozrośnie się po cichu.
 *
 * Czego sprawdzian NIE wychwyci: pojedynczego słowa bez odstępu i bez polskiego
 * znaku, na przykład nazwy własnej wpisanej wprost w składnik. Granica jest
 * nazwana wprost, bo instrument, który udaje szczelność, jest gorszy od
 * instrumentu o znanym zasięgu.
 *
 * Pomiar odróżnia brak wyniku od wyniku pustego: sprawdzian upada, gdy nie
 * znalazł plików albo nie znalazł ani jednego łańcucha. Zero trafień
 * w katalogu, którego nie ma, nie jest wynikiem.
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

/** Katalog terenu; sprawdzian stoi w nim, więc idzie od własnego położenia. */
const KORZEN = new URL('.', import.meta.url).pathname;

/** Plik, który tekst NOSI — jedyny. */
const KATALOG_TRESCI = 'tresci.ts';

/** Plik, który niesie rysunki; jego łańcuchy muszą być rysunkami. */
const ZESTAW_ZNAKOW = 'ikony.ts';

/** Przyrostek plików sprawdzianów — instrumentów, nie okna. */
const SPRAWDZIAN = '.test.ts';

/** Polskie znaki diakrytyczne — po nich poznaje się zdanie po polsku. */
const DIAKRYTYKI = /[ąćęłńóśźżĄĆĘŁŃÓŚŹŻ]/;

/**
 * Znak właściwy nazwie klasy albo selektorowi. Wyraz zdania go nie niesie.
 */
const ZNAK_SELEKTORA = /[-.#[\]=,>+~*]/;

/**
 * Czy łańcuch ma kształt techniczny, czyli nie da się go przeczytać jak zdania.
 *
 * Rozstrzyga próba na członach: w wykazie klas i w selektorze KAŻDY człon
 * rozdzielony odstępem niesie znak selektora albo jest nazwą elementu, czyli
 * ma najwyżej dwie litery. W zdaniu członów takich nie ma — wyraz polski
 * dłuższy niż dwie litery i bez znaku selektora przepada na tej próbie,
 * choćby całe zdanie składało się ze znaków, których używa arkusz.
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

/** Znalezisko: jeden łańcuch wraz z miejscem, w którym stoi. */
interface Lancuch {
  plik: string;
  wiersz: number;
  tresc: string;
  /** Czy łańcuch idzie do dziennika wywołaniem `console.*`. */
  doDziennika: boolean;
}

/** Wykaz plików terenu, wraz z podkatalogami, w kolejności odczytu. */
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
 * Wydobywa łańcuchy z pliku, pomijając komentarze i wyrażenia regularne.
 *
 * Bez pomijania komentarzy pomiar mierzyłby polszczyznę komentarzy zamiast
 * łańcuchów, a bez pomijania wyrażeń regularnych ukośnik klasy znaków
 * wyglądałby jak początek komentarza i zjadał resztę pliku.
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

const wykaz = wykazPlikow(KORZEN);
const plikiOkna = wykaz.filter((nazwa) => !nazwa.endsWith(SPRAWDZIAN));
const pominiete = wykaz.filter((nazwa) => nazwa.endsWith(SPRAWDZIAN));
const wszystkie: Lancuch[] = plikiOkna.flatMap((nazwa) =>
  lancuchyPliku(nazwa, pliki.readFileSync(`${KORZEN}${nazwa}`, 'utf8')),
);
const doDziennika = wszystkie.filter((l) => l.doDziennika);
const pozaKatalogiem = wszystkie.filter(
  (l) => l.plik !== KATALOG_TRESCI && !l.doDziennika,
);

function wypisz(nazwa: string, wartosc: unknown): void {
  console.log(`           ${nazwa}: ${JSON.stringify(wartosc)}`);
}

await bieg('droga wejścia — katalog treści', {
  'instrument zmierzył to, co miał zmierzyć'() {
    // Zapora przed wynikiem pustym: zero trafień w katalogu, którego nie ma,
    // wygląda tak samo jak zero naruszeń.
    sprawdz(wykaz.length > 0, `nie znaleziono ani jednego pliku w ${KORZEN}`);
    sprawdz(plikiOkna.includes(KATALOG_TRESCI), `nie znaleziono katalogu ${KATALOG_TRESCI}`);
    sprawdz(plikiOkna.includes(ZESTAW_ZNAKOW), `nie znaleziono zestawu znaków ${ZESTAW_ZNAKOW}`);
    sprawdz(wszystkie.length > 0, 'nie znaleziono ani jednego łańcucha — wzorzec nic nie łapie');
    sprawdz(pozaKatalogiem.length > 0, 'poza katalogiem nie znaleziono ani jednego łańcucha');
    wypisz('plików okna', plikiOkna.length);
    wypisz('sprawdzianów pominiętych', pominiete);
    wypisz('łańcuchów w plikach okna', wszystkie.length);
    wypisz('łańcuchów poza katalogiem treści', pozaKatalogiem.length);
    wypisz('łańcuchów oddanych do dziennika', doDziennika.map((l) => `${l.plik}:${l.wiersz}`));
  },

  'instrument rozpoznaje zdanie po polsku, gdy je zobaczy'() {
    // Kontrola dodatnia: bez niej sprawdzian milczałby także wtedy, gdyby
    // wydobywanie łańcuchów przestało cokolwiek wydobywać.
    const probka = lancuchyPliku('probka.ts', "const a = 'Hasło nie spełnia wymagań.';\n");
    rowne(probka.length, 1, 'próbka dała jeden łańcuch');
    sprawdz(DIAKRYTYKI.test(probka[0]!.tresc), 'wzorzec nie rozpoznał polskiego zdania');
    const wKomentarzu = lancuchyPliku('probka.ts', "// 'Hasło w komentarzu'\nconst a = 1;\n");
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
    // Kontrola dodatnia dla reguły kształtu: bez niej reguła mogłaby przepuszczać
    // wszystko i milczeć tak samo jak reguła spełniona.
    const zdanie: Lancuch = {
      plik: 'skladniki/probka.ts',
      wiersz: 1,
      tresc: 'Podaj login albo adres e-mail.',
      doDziennika: false,
    };
    sprawdz(!ksztaltTechniczny(zdanie, false), 'zdanie w składniku nie zostało rozpoznane');
    const klasy: Lancuch = {
      plik: 'skladniki/probka.ts',
      wiersz: 2,
      tresc: 'dn-btn dn-btn--sygnal',
      doDziennika: false,
    };
    sprawdz(ksztaltTechniczny(klasy, false), 'wykaz klas został wzięty za zdanie');
  },

  'wpis do dziennika nie przemyca treści dla użytkownika'() {
    // Odstępstwo jest wąskie i policzone: wpis diagnostyczny wolno pisać po
    // polsku, ale nie wolno nim zastąpić katalogu treści.
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
