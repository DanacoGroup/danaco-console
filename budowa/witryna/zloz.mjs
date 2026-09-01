// Złożenie witryny `danaco-console.pl` — dziewięć stron statycznych z jednego
// szkieletu i jednego wykazu wydań.
//
// PO CO GENERATOR, SKORO TO STRONY STATYCZNE. Bo inaczej nagłówek, stopka
// i nawigacja stałyby w dziewięciu kopiach, a dziewiąta kopia rozjeżdża się
// pierwszego dnia. Tu szkielet jest jeden, treść stron jest jedna, a chronologia
// wydań ma JEDNO źródło — `wydania.json`. Dopisanie wydania to dopisanie wiersza
// w tym pliku i przebieg tego polecenia; nikt nie edytuje HTML-a ręcznie.
//
// ŻADNEJ ZALEŻNOŚCI Z SIECI. Witryna ma stanąć na `danaco-web` bez `npm install`
// i bez łańcucha budowy — wyjściem jest katalog plików, które serwer po prostu
// oddaje. Krój pisma, style i skrypt idą w treści strony, nie z cudzego adresu.
import { mkdir, writeFile, readFile, rm, cp, stat } from 'node:fs/promises';
import { dirname, join } from 'node:path';
import { fileURLToPath } from 'node:url';

import { STRONY } from './tresc/strony.mjs';
import { szkielet, wczytajSzkielet } from './tresc/szkielet.mjs';
import { stronaPobierz, najnowszeWydanie } from './tresc/pobierz.mjs';
import { stronaPrototypy, katalogPrototypow } from './tresc/prototypy.mjs';

const KATALOG = dirname(fileURLToPath(import.meta.url));
const WYJSCIE = join(KATALOG, 'dist');

const wykaz = JSON.parse(await readFile(join(KATALOG, 'wydania.json'), 'utf8'));

// WYKAZ SPRAWDZA SIĘ PRZED ZAPISEM, NIE W TRAKCIE. Literówka w kluczu
// (`wydanie` zamiast `wydania`) dałaby witrynę GOTOWĄ DO WGRANIA, która
// twierdzi „nie ma jeszcze żadnego wydania" — bo `stronaPobierz` brak tablicy
// czyta jako wykaz pusty — a polecenie wywróciłoby się dopiero na liczeniu
// długości, po zapisaniu wszystkich stron i pliku `wydania.json`. Kanał wydań
// milczałby wtedy (aplikacja odczytuje `wydania` jako nie-tablicę i nie
// proponuje niczego). Odmowa nazywa BRAK i pada, zanim cokolwiek trafi na dysk.
if (!Array.isArray(wykaz.wydania)) {
  console.error(
    'witryna NIE złożona: plik wydania.json nie ma klucza „wydania" z tablicą wydań.\n' +
      '  odczytano klucze: ' + Object.keys(wykaz).join(', ') + '\n' +
      '  pusta chronologia zapisuje się jako "wydania": [] — brak klucza to literówka,\n' +
      '  a nie deklaracja braku wydań; katalog dist został nietknięty.',
  );
  process.exit(1);
}


// KROK ODBIORU — STRONA POWSTAJE NA ODPOWIEDŹ KANAŁU, NIE NA SŁOWO WYKAZU.
//
// Strona wystawia przycisk „Pobierz" pod adresem wziętym z wykazu, a adres
// wpisany do wykazu i adres, pod którym plik naprawdę leży, rozchodzą się
// bezgłośnie: pozycja przeniesiona pod ścieżkę zamkniętą uwierzytelnieniem
// zostawia stronę zapowiadającą pobranie bez pytania o hasło. Dlatego przed
// złożeniem idzie na każdą pozycję żądanie HEAD i strona NIE POWSTAJE, gdy
// odpowiedź jest inna niż umówiona.
//
// Pozycja z `chronione_haslem: true` odpowiada 401 i tak ma być — to jej
// deklarowany stan, a strona zapowiada przy niej pytanie o hasło. Każdy inny kod
// jest rozjazdem wykazu z kanałem. Rozmiar sprawdzany jest przy okazji, bo
// nagłówek `content-length` rozstrzyga, czy pod adresem leży ten plik, o którym
// wykaz mówi, czy inny o tej samej nazwie.
const CZAS_ODBIORU_MS = 15000;

/** Pozycje wykazu niosące adres pliku, wraz z nazwą miejsca — nazwa miejsca
 *  wchodzi do odmowy, żeby było wiadomo, którą pozycję poprawić. Pozycje
 *  `w_przygotowaniu` są pominięte: z założenia nie mają jeszcze pliku. */
function pozycjeZPlikiem(w) {
  const zebrane = [];
  for (const [klucz, wartosc] of Object.entries(w)) {
    if (klucz === 'pola' || klucz === 'w_przygotowaniu') continue;
    if (Array.isArray(wartosc)) {
      wartosc.forEach((p, i) => {
        if (p && typeof p === 'object' && p.plik) zebrane.push({ p, gdzie: `${klucz}[${i}]` });
      });
    } else if (wartosc && typeof wartosc === 'object' && wartosc.plik) {
      zebrane.push({ p: wartosc, gdzie: klucz });
    }
  }
  return zebrane;
}

async function odbiorKanalu(w) {
  // Kanał niewdrożony nie ma czego odpowiadać, a strona mówi o tym wprost
  // i nie wystawia ani jednego przycisku — nie ma więc czego odbierać.
  if (!w.kanal || w.kanal.wdrozony !== true) {
    console.log('odbiór pominięty: kanał wykazu nie jest wdrożony — strona powie o tym wprost');
    return;
  }
  const podstawa = typeof w.kanal.adres === 'string' ? w.kanal.adres : '';
  const zarzuty = [];
  for (const { p, gdzie } of pozycjeZPlikiem(w)) {
    let adres;
    try {
      adres = new URL(p.plik, podstawa).href;
    } catch {
      zarzuty.push(`${gdzie}: adres „${p.plik}" jest nieczytelny`);
      continue;
    }
    const dopuszczone = p.chronione_haslem === true ? [200, 401] : [200];
    let odpowiedz;
    try {
      odpowiedz = await fetch(adres, { method: 'HEAD', signal: AbortSignal.timeout(CZAS_ODBIORU_MS) });
    } catch (powod) {
      zarzuty.push(`${gdzie}: ${adres} — kanał nie odpowiedział (${powod.message})`);
      continue;
    }
    if (!dopuszczone.includes(odpowiedz.status)) {
      zarzuty.push(
        `${gdzie}: ${adres} — kod ${odpowiedz.status}, a umówione ${dopuszczone.join(' albo ')}` +
          (p.chronione_haslem === true ? ' (pozycja deklaruje ścieżkę za hasłem)' : ''),
      );
      continue;
    }
    const dlugosc = Number.parseInt(odpowiedz.headers.get('content-length') ?? '', 10);
    if (Number.isFinite(dlugosc) && Number.isFinite(p.rozmiarBajty) && dlugosc !== p.rozmiarBajty) {
      zarzuty.push(`${gdzie}: ${adres} — pod adresem leży ${dlugosc} B, a wykaz mówi o ${p.rozmiarBajty} B`);
      continue;
    }
    console.log(`= odbiór ${gdzie}: ${odpowiedz.status} (${p.nazwaPliku ?? adres})`);
  }
  if (zarzuty.length > 0) {
    console.error(
      'witryna NIE złożona: kanał odpowiada inaczej, niż mówi wykaz wydań.\n  ' +
        zarzuty.join('\n  ') +
        '\n  Strona wystawiłaby przycisk „Pobierz" pod adresem, który tego pliku nie odda;\n' +
        '  katalog dist został nietknięty.',
    );
    process.exit(1);
  }
}

await odbiorKanalu(wykaz);

await rm(WYJSCIE, { recursive: true, force: true });
await mkdir(WYJSCIE, { recursive: true });

// MATERIAŁ MARKI IDZIE DO WYJŚCIA, NIE JEST DOWIĄZANY. Żetony, kroje pisma, znak
// marki i zrzuty okien pochodzą z dostawy Właściciela (`design/`), która leży
// POZA kontrolą wersji — witryna dowiązująca do niej stanęłaby na maszynie,
// gdzie tej dostawy nie ma, i wyglądała jak dokument bez stylu. Dlatego pliki
// są skopiowane raz do `witryna/media/`, a stąd wchodzą do `dist/media/`
// przy każdym złożeniu; `dist` jest kasowany na wejściu, więc kopia musi być tu.
const MEDIA = join(KATALOG, 'media');
let media = false;
try {
  media = (await stat(MEDIA)).isDirectory();
} catch {
  media = false;
}
if (media) {
  await cp(MEDIA, join(WYJSCIE, 'media'), { recursive: true });
} else {
  // Brak materiału nie jest awarią składania, ale JEST brakiem — strona wyjdzie
  // bez znaku i bez krojów, i lepiej, żeby powiedziało to polecenie, niż żeby
  // ktoś zobaczył to dopiero w przeglądarce.
  console.warn(
    'UWAGA: katalog media/ nie istnieje — witryna złoży się BEZ znaku marki, krojów pisma\n' +
      '  i zrzutów okien. Materiał kopiuje się z design/ (dostawa Właściciela).',
  );
}

// PROTOTYPY OKIEN IDĄ DO WYJŚCIA W CAŁOŚCI, Z UKŁADEM KATALOGÓW.
//
// Makiety z `design/05-okna/` są interaktywne i wołają wspólne zasoby ścieżkami
// względnymi (`../../zasoby/css/fundament.css`). Skopiowanie samych plików HTML
// dałoby okna bez stylów, a przepisanie tych ścieżek byłoby przerabianiem
// dostawy. Dlatego do `dist/prototypy/` wchodzą OBA katalogi w swoim wzajemnym
// położeniu: `05-okna/` (rodziny okien) oraz `zasoby/` (styl, kroje, ikony,
// skrypt prototypu).
//
// CZEGO ŚWIADOMIE NIE KOPIUJEMY: `KANON.md`, `01-dokumentacja-*`, `03-marka/`
// (księgi znaku), `04-portfolio/` i `INDEKS.html`. To wewnętrzna dokumentacja
// projektowa producenta, a witryna jest publiczna — o jej wydaniu na świat
// rozstrzyga Właściciel, nie generator. Skutek uboczny do wiadomości: trzy
// prototypy samodzielne mają w nagłówku odnośnik powrotny do `INDEKS.html`,
// który tu nie odpowiada; jest to odnośnik wewnątrz dostawy, nie element okna.
const DOSTAWA = join(KATALOG, '..', '..', 'design');
const PROTOTYPY = join(WYJSCIE, 'prototypy');
const DO_SKOPIOWANIA = ['05-okna', 'zasoby'];

let rodziny = [];
try {
  await stat(DOSTAWA);
  await mkdir(PROTOTYPY, { recursive: true });
  for (const czlon of DO_SKOPIOWANIA) {
    try {
      await cp(join(DOSTAWA, czlon), join(PROTOTYPY, czlon), { recursive: true });
    } catch (powod) {
      // Brak jednego członu nie przerywa kopiowania pozostałych: pętla przerwana
      // na pierwszym braku zostawia okna bez arkuszy i skryptu, czyli strony
      // wyglądające na zepsute zamiast witryny bez jednego działu.
      console.warn(`UWAGA: członu prototypów „${czlon}" nie skopiowano (${powod.message}).`);
    }
  }
  rodziny = await katalogPrototypow(DOSTAWA);
} catch (powod) {
  // Brak dostawy nie wywraca składania — strona „Prototypy" powie wtedy wprost,
  // że okien nie ma, zamiast pokazywać rysunek zastępczy. Ale powiedzieć trzeba
  // tu, przy składaniu, a nie dopiero w przeglądarce.
  console.warn(`UWAGA: prototypów okien nie skopiowano (${powod.message}) — strona „Prototypy" powie o braku.`);
}

// Strony „Pobierz" i „Prototypy" powstają z materiału, nie z napisanej treści —
// dlatego wchodzą do zbioru dopiero tutaj, po odczycie wykazu i dostawy.
const strony = STRONY.map((s) => {
  if (s.adres === 'pobierz.html') return stronaPobierz(s, wykaz);
  if (s.adres === 'prototypy.html') return stronaPrototypy(s, rodziny);
  return s;
});

// Szkielet witryny grupy czytany jest RAZ, przed pętlą — nie dziesięć razy.
const szablonKodu = await wczytajSzkielet();

for (const strona of strony) {
  const html = szkielet(strona, strony, szablonKodu);
  await writeFile(join(WYJSCIE, strona.adres), html, 'utf8');
}

await writeFile(join(WYJSCIE, 'wydania.json'), JSON.stringify(wykaz, null, 2), 'utf8');

const liczbaWydan = wykaz.wydania.length;
const liczbaOkien = rodziny.reduce((suma, r) => suma + r.okna.length, 0);
console.log(`witryna złożona: ${strony.length} stron → ${WYJSCIE}`);
console.log(
  liczbaOkien === 0
    ? 'prototypów okien: 0 — strona „Prototypy" mówi o braku dostawy'
    : `prototypów okien: ${liczbaOkien} w ${rodziny.length} rodzinach (${rodziny.map((r) => `${r.nazwa}: ${r.okna.length}`).join(', ')})`,
);
console.log(
  liczbaWydan === 0
    ? 'wydań w wykazie: 0 — strona „Pobierz" mówi wprost, że nie ma czego pobrać'
    // Najnowsze WYLICZONE porównaniem wersji — tak jak na stronie i tak jak
    // w banerze aplikacji; pozycja zerowa wykazu bywa czym innym.
    : `wydań w wykazie: ${liczbaWydan}, najnowsze: ${najnowszeWydanie(wykaz.wydania).wersja} z ${najnowszeWydanie(wykaz.wydania).data}`,
);
