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
// położeniu: `05-okna/` (rodziny okien), `_samodzielne/` (prototypy w jednym
// pliku) oraz `zasoby/` (styl, kroje, ikony, skrypt prototypu).
//
// CZEGO ŚWIADOMIE NIE KOPIUJEMY: `KANON.md`, `01-dokumentacja-*`, `03-marka/`
// (księgi znaku), `04-portfolio/` i `INDEKS.html`. To wewnętrzna dokumentacja
// projektowa producenta, a witryna jest publiczna — o jej wydaniu na świat
// rozstrzyga Właściciel, nie generator. Skutek uboczny do wiadomości: trzy
// prototypy samodzielne mają w nagłówku odnośnik powrotny do `INDEKS.html`,
// który tu nie odpowiada; jest to odnośnik wewnątrz dostawy, nie element okna.
const DOSTAWA = join(KATALOG, '..', '..', 'design');
const PROTOTYPY = join(WYJSCIE, 'prototypy');
const DO_SKOPIOWANIA = ['05-okna', '_samodzielne', 'zasoby'];

let rodziny = [];
try {
  await stat(DOSTAWA);
  await mkdir(PROTOTYPY, { recursive: true });
  for (const czlon of DO_SKOPIOWANIA) {
    await cp(join(DOSTAWA, czlon), join(PROTOTYPY, czlon), { recursive: true });
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
